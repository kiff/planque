package browserworker

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// noServiceWorker takes from a health check's documents what its request
// interception can't see: registering a service worker, sockets (WebSocket,
// WebTransport), whose messages never pause, and workers (Worker,
// SharedWorker). It's a second layer: a check runs no page script at all
// (Emulation.setScriptExecutionDisabled), and closes any worker target before
// it runs. A URL block doesn't stop a socket (measured), so none is relied on. WebRTC is already gone (noWebRTC). Socket addresses are also
// blocked at the network layer (socketBlock), frames and workers included.
const noServiceWorker = `(() => {
  const refuse = (what) => function () { throw new Error(what + " isn't allowed in a sign-in check"); };
  try { if (navigator.serviceWorker) Object.defineProperty(navigator.serviceWorker, "register", {value: () => Promise.reject(new Error("no service workers in a sign-in check")), configurable: false}); } catch (e) {}
  for (const name of ["WebSocket", "WebTransport", "Worker", "SharedWorker"]) {
    try { Object.defineProperty(window, name, {value: refuse(name), configurable: false, writable: false}); } catch (e) {}
  }
})();`

// healthSettle is how long a health check waits, after the page loads, for
// a site that sends a signed-out visitor to its sign-in page by script.
const healthSettle = 1500 * time.Millisecond

// healthTimeout bounds one check, the browser's start included.
const healthTimeout = 20 * time.Second

// Health runs one sign-in health check (browser.HealthCheck): a session
// held to the one site, every request that isn't a GET or HEAD failed before
// it leaves the browser, one navigation to the fixed page, and the verdict
// read from where it ended. It takes no input and runs no step. Nothing of
// the page is kept; its cookies after come back as the state, as a session's
// do, since many sites rotate the sign-in on each visit and dropping the new
// one would sign the owner out by planque's own check.
func (l *Local) Health(ctx context.Context, check browser.HealthCheck) (browser.HealthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, healthTimeout)
	defer cancel()
	s, err := l.start(ctx, browser.Start{Grant: check.Grant, State: check.State, Sites: []string{check.Site}}, true)
	if err != nil {
		return browser.HealthResult{}, err
	}
	res := browser.HealthResult{Verdict: browser.HealthUnknown}
	defer func() {
		if !s.closed {
			s.end()
		}
	}()
	p, err := s.Open(ctx, check.URL)
	if err == nil {
		// A site may send a signed-out visitor on by script after load.
		select {
		case <-time.After(healthSettle):
		case <-ctx.Done():
		}
		if s.settle(ctx) == nil {
			if again, err := s.page(ctx); err == nil {
				p = again
			}
		}
		res.Verdict = browser.HealthVerdict(check.URL, check.Shows, p)
	}
	if res.Verdict == browser.HealthSignedIn {
		// Only a check that found the sign-in holding gives its cookies
		// back: a signed-out or unclear one never replaces a kept sign-in.
		if st, err := s.Close(ctx); err == nil {
			res.State = st
		}
	}
	res.Blocked = int(s.blocked.Load())
	if res.Blocked > 0 {
		s.log().Info("health check failed requests that would write", slog.String("grant", check.Grant), slog.Int("blocked", res.Blocked))
	}
	return res, nil
}

// fetchPaused is the part of Fetch.requestPaused a health check reads.
type fetchPaused struct {
	RequestID string `json:"requestId"`
	Request   struct {
		Method string `json:"method"`
	} `json:"request"`
}

// readOnly answers each request a health session's targets pause: a GET or
// HEAD goes on; anything else is failed, never sent, and counted.
func (s *session) readOnly(m cdpMessage) {
	if m.Method != "Fetch.requestPaused" {
		return
	}
	var p fetchPaused
	if json.Unmarshal(m.Params, &p) != nil || p.RequestID == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if readOnlyMethod(p.Request.Method) {
		_ = s.c.call(ctx, m.SessionID, "Fetch.continueRequest", map[string]any{"requestId": p.RequestID}, nil)
		return
	}
	s.blocked.Add(1)
	_ = s.c.call(ctx, m.SessionID, "Fetch.failRequest", map[string]any{"requestId": p.RequestID, "errorReason": "BlockedByClient"}, nil)
}

// readOnlyMethod: the methods a health check lets through.
func readOnlyMethod(method string) bool { return method == "GET" || method == "HEAD" }

// readOnlyTarget puts a child target (a frame, a worker) of a health session
// under the same rule before it runs. One that can't be is closed, never
// run: no request of a health check goes out unchecked.
func (s *session) readOnlyTarget(ctx context.Context, p attachedTarget) bool {
	// A worker, however it was made (a constructor borrowed from a blank
	// frame included), is closed before it runs: it attaches paused, and a
	// socket in a worker isn't blocked by address. A check needs none.
	if strings.Contains(p.TargetInfo.Type, "worker") {
		s.log().Info("a health check's worker closed before it ran", slog.String("type", p.TargetInfo.Type))
		if err := s.c.call(ctx, "", "Target.closeTarget", map[string]any{"targetId": p.TargetInfo.TargetID}, nil); err != nil {
			_ = s.c.call(ctx, "", "Target.detachFromTarget", map[string]any{"sessionId": p.SessionID}, nil)
		}
		return false
	}
	err := s.c.call(ctx, p.SessionID, "Fetch.enable", map[string]any{"patterns": []map[string]string{{"urlPattern": "*", "requestStage": "Request"}}}, nil)
	if err == nil && (p.TargetInfo.Type == "iframe" || p.TargetInfo.Type == "page") {
		err = s.c.call(ctx, p.SessionID, "Emulation.setScriptExecutionDisabled", map[string]any{"value": true}, nil)
		if err == nil && !s.noDocScript {
			err = s.c.call(ctx, p.SessionID, "Page.addScriptToEvaluateOnNewDocument", map[string]any{"source": noServiceWorker}, nil)
		}
	}
	if err == nil {
		return true
	}
	s.log().Warn("a health check's child target couldn't be held to reads: closed, never run", slog.String("target", p.TargetInfo.TargetID), slog.String("err", err.Error()))
	if err := s.c.call(ctx, "", "Target.closeTarget", map[string]any{"targetId": p.TargetInfo.TargetID}, nil); err != nil {
		_ = s.c.call(ctx, "", "Target.detachFromTarget", map[string]any{"sessionId": p.SessionID}, nil)
	}
	return false
}
