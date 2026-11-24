package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// sendRunState is the extension's signed run_state for run, sent with key.
func (c *chromeRig) sendRunState(key stepkey.Local, run, state string, at time.Time) map[string]any {
	c.t.Helper()
	body := map[string]any{"run": run, "state": state, "at": at.UTC().Format(time.RFC3339Nano)}
	sig, err := stepkey.SignBody(context.Background(), key, labelRunState, body)
	if err != nil {
		c.t.Fatal(err)
	}
	msg := map[string]any{"type": "run_state", "sig": sig}
	for k, v := range body {
		msg[k] = v
	}
	c.r.event("$default", "c1", msg)
	return c.r.poster.last("c1")
}

// readRun is the grant's run and its state, as the taf's runtime reads it with key.
func (c *chromeRig) readRun(g *Grant, key string) (int, map[string]any) {
	c.t.Helper()
	return c.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/chrome/run", nil, nil, key)
}

// Pause in the owner's Chrome (tafs docs/chrome.md, the floating bar): the
// refusal of a paused step names its run; the run's Chrome says its state in
// a signed run_state, kept only for a run planque started for that Chrome,
// never after stopped, never from another key, never from far ahead; and the
// taf's own key reads the state, nothing else does.
func TestAPausedRunSaysSoAndItsChromeSaysWhenItGoesOn(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); sc["read"] == nil {
		t.Fatalf("open: %v", sc)
	}
	runs, _ := c.runs()
	run := runs[0]
	// Its Chrome hasn't said a state yet (an extension from before run_state): unreported, never working,
	// so the runtime doesn't take a pause it can't follow for a Resume.
	if st, out := c.readRun(g, "k_taf"); st != 200 || out["run"] != run || out["state"] != runUnreported {
		t.Fatalf("a run with no state said: %d %v", st, out)
	}
	if m := c.sendRunState(c.r.ext, run, runUnreported, c.s.now()); m["type"] != "run_state_refused" || m["code"] != "bad_state" {
		t.Fatalf("unreported sent by a Chrome: %v", m)
	}

	// The person pauses: the next step is answered paused_by_owner, which names the run.
	start := c.s.now()
	if m := c.sendRunState(c.r.ext, run, runPaused, start); m["type"] != "run_state_kept" || m["run"] != run {
		t.Fatalf("paused: %v", m)
	}
	answer := c.answer
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return false, nil, "paused_by_owner" }
	_, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"})
	if code(sc) != "paused_by_owner" || sc["run"] != run || sc["grant"] != g.ID {
		t.Fatalf("a paused step: %v", sc)
	}
	if st, out := c.readRun(g, "k_taf"); st != 200 || out["run"] != run || out["state"] != runPaused {
		t.Fatalf("a paused run: %d %v", st, out)
	}

	// Only the grant's own taf reads it: not Setup, not another account's key.
	for _, key := range []string{"k_setup", "k_other", ""} {
		if st, _ := c.readRun(g, key); st == 200 {
			t.Errorf("%q read the run", key)
		}
	}

	// Refused: another key, a run planque never started for this Chrome, a time far ahead or before the run.
	ok, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if m := c.sendRunState(stepkey.Local{Key: ok}, run, runWorking, start.Add(time.Second)); m["type"] != "run_state_refused" || m["code"] != "unsigned" {
		t.Fatalf("another key: %v", m)
	}
	if m := c.sendRunState(c.r.ext, "run_nobody", runWorking, start.Add(time.Second)); m["type"] != "run_state_refused" || m["code"] != "unknown_run" {
		t.Fatalf("an unknown run: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, runWorking, c.s.now().Add(time.Hour)); m["type"] != "run_state_refused" || m["code"] != "bad_at" {
		t.Fatalf("far ahead: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, runWorking, start.Add(-time.Hour)); m["type"] != "run_state_refused" || m["code"] != "bad_at" {
		t.Fatalf("before the run: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, "dancing", start.Add(time.Second)); m["type"] != "run_state_refused" || m["code"] != "bad_state" {
		t.Fatalf("a state that isn't one: %v", m)
	}

	// Resume: the state goes on, and a resend of the older pause changes nothing.
	if m := c.sendRunState(c.r.ext, run, runWorking, start.Add(2*time.Second)); m["type"] != "run_state_kept" {
		t.Fatalf("resumed: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, runPaused, start); m["type"] != "run_state_kept" {
		t.Fatalf("the older pause again: %v", m)
	}
	if _, out := c.readRun(g, "k_taf"); out["state"] != runWorking {
		t.Fatalf("after Resume: %v", out)
	}
	c.answer = answer

	// Stop is final: a late Resume is refused, whatever its time, and a resend of Stop is kept.
	if m := c.sendRunState(c.r.ext, run, runStopped, start.Add(3*time.Second)); m["type"] != "run_state_kept" {
		t.Fatalf("stopped: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, runWorking, start.Add(4*time.Second)); m["type"] != "run_state_refused" || m["code"] != "run_stopped" {
		t.Fatalf("a Resume after Stop: %v", m)
	}
	if m := c.sendRunState(c.r.ext, run, runStopped, start.Add(5*time.Second)); m["type"] != "run_state_kept" {
		t.Fatalf("Stop again: %v", m)
	}
	if _, out := c.readRun(g, "k_taf"); out["state"] != runStopped {
		t.Fatalf("after Stop: %v", out)
	}
}

// A paused run stays the grant's past its idle time, so the step sent again on
// Resume goes on in the same tab; one that isn't paused starts a new run, as before.
func TestAPausedRunOutlivesItsIdleTime(t *testing.T) {
	for _, paused := range []bool{true, false} {
		c, g := tabRig(t)
		if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); sc["read"] == nil {
			t.Fatalf("open: %v", sc)
		}
		runs, _ := c.runs()
		if paused {
			if m := c.sendRunState(c.r.ext, runs[0], runPaused, c.s.now()); m["type"] != "run_state_kept" {
				t.Fatalf("paused: %v", m)
			}
		}
		// The run's last call ended longer ago than its idle time.
		ctx := context.Background()
		var doc chromeRunDoc
		it, err := getDoc(ctx, c.s.KV, tenantPK(tenant), skChromeRun(g), &doc)
		if err != nil {
			t.Fatal(err)
		}
		doc.At = doc.At.Add(-chromeRunIdle - time.Minute)
		ver := it.Ver
		if err := putDoc(ctx, c.s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: ver + 1, TTL: it.TTL}, doc, store.Cond{Ver: &ver}); err != nil {
			t.Fatal(err)
		}
		if _, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"}); sc["read"] == nil {
			t.Fatalf("read: %v", sc)
		}
		after, _ := c.runs()
		if same := after[len(after)-1] == runs[0]; same != paused {
			t.Errorf("paused %v: runs %v", paused, after)
		}
	}
}
