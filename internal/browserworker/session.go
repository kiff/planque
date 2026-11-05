package browserworker

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/browser"
)

//go:embed page.js
var pageJS string

//go:embed step.js
var stepJS string

// noWebRTC takes WebRTC away from a document: a page's script can't open
// a peer connection, so no UDP leaves the browser for any address.
const noWebRTC = `(() => {
	for (const name of ['RTCPeerConnection', 'webkitRTCPeerConnection', 'RTCDataChannel', 'RTCIceCandidate', 'RTCSessionDescription', 'RTCDtlsTransport', 'RTCIceTransport', 'RTCSctpTransport']) {
		try { Object.defineProperty(window, name, {value: undefined, writable: false, configurable: false}); } catch (e) {}
	}
	try { Object.defineProperty(navigator, 'mediaDevices', {value: undefined, writable: false, configurable: false}); } catch (e) {}
})();`

// Options is how the worker runs its browsers.
type Options struct {
	// Chrome is the browser's path; empty, chromePath looks.
	Chrome string
	// AllowLocal lets tests run against a shop on loopback, and Map sends
	// a host:port there ("shop.example:443" to the shop's TLS server), the
	// browser taking the test certificate.
	AllowLocal bool
	Map        map[string]string
	// Resolve is the proxy's resolver; nil, the system's.
	Resolve func(ctx context.Context, host string) ([]netip.Addr, error)
	Logger  *slog.Logger
	// NavigateWait bounds how long a page may take to load.
	NavigateWait time.Duration
}

// Local is the worker in-process: planque's browser.Worker for tests, and
// what Server runs behind its API. Each session is a Chromium of its own
// with a proxy of its own, gone when the session closes. At most
// MaxSessions run at once (measured: about 400 MB for the first, 270 to
// 300 MB for each more, so three fit 2 GB with headroom); beyond it a
// start is refused with ErrBusy.
type Local struct {
	Options
	// MaxSessions bounds the sessions open at once; 0 is 3.
	MaxSessions int

	mu   sync.Mutex
	open map[string]*session
	seq  int
}

// ErrBusy: the worker runs its most sessions already.
var ErrBusy = errors.New("the browser worker is busy: every session is in use")

func (l *Local) limit() int {
	if l.MaxSessions > 0 {
		return l.MaxSessions
	}
	return 3
}

// Resume is the open session with that id.
func (l *Local) Resume(_ context.Context, id string) (browser.Session, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	s, ok := l.open[id]
	if !ok {
		return nil, errors.New("no such session")
	}
	return s, nil
}

// Open is how many sessions run now.
func (l *Local) Open() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.open)
}

func (l *Local) forget(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.open, id)
}

// Start begins a session: a proxy bound to the sites, a fresh browser
// behind it, downloads denied, and the kept cookies restored.
func (l *Local) Start(ctx context.Context, start browser.Start) (browser.Session, error) {
	path, err := chromePath(l.Chrome)
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	if l.open == nil {
		l.open = map[string]*session{}
	}
	if len(l.open) >= l.limit() {
		l.mu.Unlock()
		return nil, ErrBusy
	}
	l.seq++
	id := fmt.Sprintf("s%d_%s", l.seq, rand.Text())
	l.open[id] = nil // the slot is taken while the browser starts
	l.mu.Unlock()
	p := &Proxy{Sites: start.Sites, AllowLocal: l.AllowLocal, Map: l.Map, Logger: l.Logger}
	if l.Resolve != nil {
		p.Resolve = l.Resolve
	}
	addr, err := p.Start()
	if err != nil {
		l.forget(id)
		return nil, err
	}
	c, err := launch(context.Background(), path, addr, l.AllowLocal)
	if err != nil {
		p.Close()
		l.forget(id)
		return nil, err
	}
	s := &session{c: c, proxy: p, wait: l.NavigateWait, logger: l.Logger, id: id, owner: l}
	if s.wait == 0 {
		s.wait = 20 * time.Second
	}
	if err := s.attach(ctx, start.State); err != nil {
		s.end()
		return nil, err
	}
	l.mu.Lock()
	l.open[id] = s
	l.mu.Unlock()
	return s, nil
}

type session struct {
	id      string
	owner   *Local
	c       *chrome
	proxy   *Proxy
	target  string
	sid     string
	wait    time.Duration
	closed  bool
	grantID string
	logger  *slog.Logger
}

// state is a session's storage state: its cookies, as the browser gives
// them. Nothing else of a site is kept.
type state struct {
	Cookies []json.RawMessage `json:"cookies"`
}

func (s *session) attach(ctx context.Context, kept []byte) error {
	if err := s.c.call(ctx, "", "Browser.setDownloadBehavior", map[string]any{"behavior": "deny"}, nil); err != nil {
		return err
	}
	var created struct {
		TargetID string `json:"targetId"`
	}
	if err := s.c.call(ctx, "", "Target.createTarget", map[string]any{"url": "about:blank"}, &created); err != nil {
		return err
	}
	s.target = created.TargetID
	var attached struct {
		SessionID string `json:"sessionId"`
	}
	if err := s.c.call(ctx, "", "Target.attachToTarget", map[string]any{"targetId": s.target, "flatten": true}, &attached); err != nil {
		return err
	}
	s.sid = attached.SessionID
	for _, m := range []string{"Page.enable", "Network.enable", "Runtime.enable"} {
		if err := s.c.call(ctx, s.sid, m, nil, nil); err != nil {
			return err
		}
	}
	if err := s.c.call(ctx, s.sid, "Page.setDownloadBehavior", map[string]any{"behavior": "deny"}, nil); err != nil {
		return err
	}
	// WebRTC's UDP goes around an HTTP proxy, and the browser's policy flag
	// alone still lets STUN checks out (measured): every document the
	// session loads, frames included, has WebRTC taken away before its
	// own script runs, as a property it can't put back.
	if err := s.c.call(ctx, s.sid, "Page.addScriptToEvaluateOnNewDocument", map[string]any{"source": noWebRTC}, nil); err != nil {
		return err
	}
	// A cross-site frame is a target of its own under site isolation, and
	// the page's script doesn't run in it. Every child target (frames,
	// frames in frames, workers) is attached with its debugger paused, gets
	// the same script, and only then runs: no frame's script sees WebRTC.
	s.c.handle(s.armChild)
	if err := s.c.call(ctx, s.sid, "Target.setAutoAttach", map[string]any{"autoAttach": true, "waitForDebuggerOnStart": true, "flatten": true}, nil); err != nil {
		return err
	}
	if len(kept) > 0 {
		var st state
		if err := json.Unmarshal(kept, &st); err != nil {
			return errors.New("the kept state isn't the browser's")
		}
		if len(st.Cookies) > 0 {
			cookies := make([]map[string]any, 0, len(st.Cookies))
			for _, raw := range st.Cookies {
				var ck map[string]any
				if json.Unmarshal(raw, &ck) != nil {
					continue
				}
				// What setCookies takes: the cookie's own fields, not the
				// browser's bookkeeping of it.
				for _, k := range []string{"size", "session", "priority", "sameParty", "sourceScheme", "sourcePort", "partitionKey", "partitionKeyOpaque"} {
					delete(ck, k)
				}
				cookies = append(cookies, ck)
			}
			if err := s.c.call(ctx, s.sid, "Network.setCookies", map[string]any{"cookies": cookies}, nil); err != nil {
				return fmt.Errorf("restore the sign-in: %w", err)
			}
		}
	}
	return nil
}

// cdpCaller is what armTarget drives: the browser, or a fake in a test.
type cdpCaller interface {
	call(ctx context.Context, session, method string, params any, out any) error
}

// attachedTarget is a Target.attachedToTarget event's part that matters.
type attachedTarget struct {
	SessionID  string `json:"sessionId"`
	TargetInfo struct {
		TargetID string `json:"targetId"`
		Type     string `json:"type"`
	} `json:"targetInfo"`
	WaitingForDebugger bool `json:"waitingForDebugger"`
}

// armChild takes WebRTC away from a target the browser attached (a
// cross-site frame, a worker) before its script runs, sets the same
// auto-attach on it for its own children, and lets it run.
func (s *session) armChild(m cdpMessage) {
	if m.Method != "Target.attachedToTarget" {
		return
	}
	var p attachedTarget
	if json.Unmarshal(m.Params, &p) != nil || p.SessionID == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	armTarget(ctx, s.c, p, s.log())
}

// armTarget arms one attached target and releases it, or, when a frame
// or page can't be armed, closes it instead: a target that would run with
// WebRTC never runs (fail closed). A worker has no Page domain; the script
// is for documents, and takes only once the domain is enabled on that
// target. The auto-attach applies to every kind.
func armTarget(ctx context.Context, c cdpCaller, p attachedTarget, logger *slog.Logger) (released bool) {
	if p.TargetInfo.Type == "iframe" || p.TargetInfo.Type == "page" {
		err := c.call(ctx, p.SessionID, "Page.enable", nil, nil)
		if err == nil {
			err = c.call(ctx, p.SessionID, "Page.addScriptToEvaluateOnNewDocument", map[string]any{"source": noWebRTC}, nil)
		}
		if err != nil {
			logger.Warn("a frame couldn't be armed against WebRTC: closed, never run", slog.String("target", p.TargetInfo.TargetID), slog.String("err", err.Error()))
			if err := c.call(ctx, "", "Target.closeTarget", map[string]any{"targetId": p.TargetInfo.TargetID}, nil); err != nil {
				_ = c.call(ctx, "", "Target.detachFromTarget", map[string]any{"sessionId": p.SessionID}, nil)
			}
			return false
		}
	}
	_ = c.call(ctx, p.SessionID, "Target.setAutoAttach", map[string]any{"autoAttach": true, "waitForDebuggerOnStart": true, "flatten": true}, nil)
	if p.WaitingForDebugger {
		_ = c.call(ctx, p.SessionID, "Runtime.runIfWaitingForDebugger", nil, nil)
	}
	return true
}

func (s *session) log() *slog.Logger {
	if s.logger != nil {
		return s.logger
	}
	return slog.Default()
}

// Open navigates to u and returns the page once loaded.
func (s *session) Open(ctx context.Context, u string) (browser.Page, error) {
	if s.closed {
		return browser.Page{}, errors.New("the session is closed")
	}
	var nav struct {
		ErrorText string `json:"errorText"`
	}
	ctx, cancel := context.WithTimeout(ctx, s.wait)
	defer cancel()
	if err := s.c.call(ctx, s.sid, "Page.navigate", map[string]any{"url": u}, &nav); err != nil {
		return browser.Page{}, err
	}
	if nav.ErrorText != "" {
		return browser.Page{}, errors.New("the page couldn't be opened: " + navError(nav.ErrorText))
	}
	if err := s.settle(ctx); err != nil {
		return browser.Page{}, err
	}
	return s.page(ctx)
}

// navError is a browser's network error in words.
func navError(code string) string {
	switch {
	case strings.Contains(code, "TUNNEL_CONNECTION_FAILED"), strings.Contains(code, "PROXY_CONNECTION_FAILED"):
		return "refused by the browser's egress (not a public site, or not one of the connection's sites)"
	case strings.Contains(code, "NAME_NOT_RESOLVED"):
		return "the site couldn't be resolved"
	case strings.Contains(code, "ABORTED"):
		return "the page was abandoned"
	}
	return code
}

// settle waits for the document to be complete, retrying while the frame
// is navigating.
func (s *session) settle(ctx context.Context) error {
	deadline := time.Now().Add(s.wait)
	for {
		var ev struct {
			Result struct {
				Value string `json:"value"`
			} `json:"result"`
		}
		err := s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": "document.readyState", "returnByValue": true}, &ev)
		if err == nil && ev.Result.Value == "complete" {
			return nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return err
			}
			return errors.New("the page didn't finish loading")
		}
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// page reads the page's facts with page.js.
func (s *session) page(ctx context.Context) (browser.Page, error) {
	var ev struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
		Exception *struct {
			Text string `json:"text"`
		} `json:"exceptionDetails"`
	}
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": pageJS, "returnByValue": true}, &ev)
		if err == nil && ev.Exception == nil {
			break
		}
		// The frame moved on under the script: read it again once settled.
		if err := s.settle(ctx); err != nil {
			return browser.Page{}, err
		}
	}
	if err != nil {
		return browser.Page{}, err
	}
	if ev.Exception != nil {
		return browser.Page{}, errors.New("the page couldn't be read: " + ev.Exception.Text)
	}
	var p browser.Page
	if err := json.Unmarshal([]byte(ev.Result.Value), &p); err != nil {
		return browser.Page{}, errors.New("the page couldn't be read")
	}
	// Masked at the source; masked again here, so no value of a secret
	// field leaves the worker whatever the script did.
	return browser.Mask(p, 0), nil
}

// Do runs one step on the open page and returns the page after.
func (s *session) Do(ctx context.Context, step browser.Step) (browser.Page, error) {
	if s.closed {
		return browser.Page{}, errors.New("the session is closed")
	}
	if err := step.Check(); err != nil {
		return browser.Page{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, s.wait)
	defer cancel()
	arg, _ := json.Marshal(step)
	var ev struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
		Exception *struct {
			Text string `json:"text"`
		} `json:"exceptionDetails"`
	}
	before := s.c.eventsSeen()
	if err := s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": "(" + strings.TrimSpace(stepJS) + ")(" + string(arg) + ")", "returnByValue": true, "awaitPromise": true}, &ev); err != nil {
		return browser.Page{}, err
	}
	if ev.Exception != nil {
		return browser.Page{}, errors.New("the step failed: " + ev.Exception.Text)
	}
	switch ev.Result.Value {
	case "ok":
	case "ref_unknown":
		return browser.Page{}, errors.New("the page has no element " + step.Ref())
	default:
		return browser.Page{}, errors.New("the step failed: " + ev.Result.Value)
	}
	// A click may start a navigation: give it a moment to begin, then
	// wait for the document to be complete.
	if step.Click != "" {
		waitCtx, cancelWait := context.WithTimeout(ctx, 400*time.Millisecond)
		_, _, _ = s.c.waitEvent(waitCtx, s.sid, "Page.frameStartedNavigating", before)
		cancelWait()
	}
	if err := s.settle(ctx); err != nil {
		return browser.Page{}, err
	}
	return s.page(ctx)
}

// Screenshot is the page as shown, PNG.
func (s *session) Screenshot(ctx context.Context) ([]byte, error) {
	if s.closed {
		return nil, errors.New("the session is closed")
	}
	var shot struct {
		Data string `json:"data"`
	}
	if err := s.c.call(ctx, s.sid, "Page.captureScreenshot", map[string]any{"format": "png"}, &shot); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(shot.Data)
}

// Close returns the session's cookies as its state and ends the browser.
func (s *session) Close(ctx context.Context) ([]byte, error) {
	if s.closed {
		return nil, errors.New("the session is closed")
	}
	var got struct {
		Cookies []json.RawMessage `json:"cookies"`
	}
	err := s.c.call(ctx, s.sid, "Network.getAllCookies", nil, &got)
	s.end()
	if err != nil {
		return nil, err
	}
	raw, _ := json.Marshal(state{Cookies: got.Cookies})
	return raw, nil
}

func (s *session) end() {
	s.closed = true
	s.c.stop()
	s.proxy.Close()
	if s.owner != nil {
		s.owner.forget(s.id)
	}
}

func (s *session) ID() string { return s.id }

// Frame is the page as shown, JPEG, for the live view: smaller than a
// screenshot, taken often.
func (s *session) Frame(ctx context.Context) ([]byte, error) {
	if s.closed {
		return nil, errors.New("the session is closed")
	}
	var shot struct {
		Data string `json:"data"`
	}
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		// A frame during a navigation fails; the next moment it doesn't.
		if err = s.c.call(ctx, s.sid, "Page.captureScreenshot", map[string]any{"format": "jpeg", "quality": 60}, &shot); err == nil {
			break
		}
		select {
		case <-time.After(200 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(shot.Data)
}

// Input is the owner's act in the live view, sent to the browser as the
// device events it would get from a person: nothing of the page is read.
func (s *session) Input(ctx context.Context, in browser.Input) error {
	if s.closed {
		return errors.New("the session is closed")
	}
	if err := in.Check(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	switch in.Kind {
	case "tap":
		for _, kind := range []string{"mousePressed", "mouseReleased"} {
			if err := s.c.call(ctx, s.sid, "Input.dispatchMouseEvent", map[string]any{"type": kind, "x": in.X, "y": in.Y, "button": "left", "clickCount": 1}, nil); err != nil {
				return err
			}
		}
	case "type":
		return s.c.call(ctx, s.sid, "Input.insertText", map[string]any{"text": in.Text}, nil)
	case "key":
		code := map[string]int{"Enter": 13, "Backspace": 8, "Tab": 9}[in.Key]
		for _, kind := range []string{"keyDown", "keyUp"} {
			params := map[string]any{"type": kind, "key": in.Key, "code": in.Key, "windowsVirtualKeyCode": code, "nativeVirtualKeyCode": code}
			if in.Key == "Enter" && kind == "keyDown" {
				params["text"] = "\r"
			}
			if err := s.c.call(ctx, s.sid, "Input.dispatchKeyEvent", params, nil); err != nil {
				return err
			}
		}
	case "scroll":
		return s.c.call(ctx, s.sid, "Input.dispatchMouseEvent", map[string]any{"type": "mouseWheel", "x": in.X, "y": in.Y, "deltaX": 0, "deltaY": in.DeltaY}, nil)
	case "go":
		var nav struct {
			ErrorText string `json:"errorText"`
		}
		if err := s.c.call(ctx, s.sid, "Page.navigate", map[string]any{"url": in.URL}, &nav); err != nil {
			return err
		}
		if nav.ErrorText != "" {
			return errors.New("the page couldn't be opened: " + navError(nav.ErrorText))
		}
		// Settled before the next frame, so the frame is the page asked for.
		return s.settle(ctx)
	}
	return nil
}
