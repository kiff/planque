package browserworker

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

// fakeCaller is a browser that answers each method as told and records
// what was asked of it.
type fakeCaller struct {
	fail  map[string]error
	calls []string
}

func (f *fakeCaller) call(_ context.Context, _ string, method string, _ any, _ any) error {
	f.calls = append(f.calls, method)
	return f.fail[method]
}

// A frame that can't be armed against WebRTC is closed, never released:
// no runIfWaitingForDebugger reaches it (agentgabo's condition on
// kiff/planque#89 before the worker is deployed). One that is armed is
// released; a worker, which has no Page domain, is released without it.
func TestArmTargetFailsClosed(t *testing.T) {
	frame := attachedTarget{SessionID: "s1", WaitingForDebugger: true}
	frame.TargetInfo.TargetID, frame.TargetInfo.Type = "t1", "iframe"
	for _, c := range []struct {
		name string
		fail map[string]error
	}{
		{"Page.enable fails", map[string]error{"Page.enable": errors.New("not allowed")}},
		{"the script fails", map[string]error{"Page.addScriptToEvaluateOnNewDocument": errors.New("not allowed")}},
	} {
		f := &fakeCaller{fail: c.fail}
		if armTarget(context.Background(), f, frame, slog.Default()) {
			t.Errorf("%s: the frame was released", c.name)
		}
		for _, m := range f.calls {
			if m == "Runtime.runIfWaitingForDebugger" || m == "Target.setAutoAttach" {
				t.Errorf("%s: %s was called: %v", c.name, m, f.calls)
			}
		}
		if f.calls[len(f.calls)-1] != "Target.closeTarget" {
			t.Errorf("%s: the frame wasn't closed: %v", c.name, f.calls)
		}
	}
	// Closing fails too: the frame is detached, still never released.
	f := &fakeCaller{fail: map[string]error{"Page.enable": errors.New("no"), "Target.closeTarget": errors.New("gone")}}
	if armTarget(context.Background(), f, frame, slog.Default()) || f.calls[len(f.calls)-1] != "Target.detachFromTarget" {
		t.Errorf("close failing: %v", f.calls)
	}
	// Armed: released, in order.
	ok := &fakeCaller{}
	if !armTarget(context.Background(), ok, frame, slog.Default()) {
		t.Error("an armed frame wasn't released")
	}
	want := []string{"Page.enable", "Page.addScriptToEvaluateOnNewDocument", "Target.setAutoAttach", "Runtime.runIfWaitingForDebugger"}
	if len(ok.calls) != len(want) {
		t.Fatalf("calls: %v", ok.calls)
	}
	for i := range want {
		if ok.calls[i] != want[i] {
			t.Errorf("call %d: %s, want %s", i, ok.calls[i], want[i])
		}
	}
	// A worker: no Page domain asked, released.
	worker := attachedTarget{SessionID: "s2", WaitingForDebugger: true}
	worker.TargetInfo.TargetID, worker.TargetInfo.Type = "t2", "service_worker"
	w := &fakeCaller{}
	if !armTarget(context.Background(), w, worker, slog.Default()) || len(w.calls) != 2 || w.calls[0] != "Target.setAutoAttach" {
		t.Errorf("a worker: %v", w.calls)
	}
}
