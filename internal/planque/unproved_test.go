package planque

import (
	"testing"
)

// A managed connection says it's unproved while no recorded run has proved
// its app on the managed route, so tafs cloud holds every change of it for
// the owner. A proof by the managed route lifts it; a proof by tafs's own
// client only ("own") doesn't. A grant of any other kind never carries it.
func TestAManagedGrantSaysWhenItsAppIsUnproved(t *testing.T) {
	e, _ := gated(t, []string{"Desk Demo"}, nil)
	unproved := func() bool {
		t.Helper()
		st, out := askManaged(e)
		if st != 201 {
			t.Fatalf("ask: %d %v", st, out)
		}
		u, _ := out["grant"].(map[string]any)["unproved"].(bool)
		return u
	}
	e.s.Proved = map[string]Proof{}
	if !unproved() {
		t.Error("an app nobody proved reads proved")
	}
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-20", How: "oauth", Route: "own"}}
	if !unproved() {
		t.Error("a proof through our own client lifted the managed route's hold")
	}
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-20", How: "oauth", Route: "managed"}}
	if unproved() {
		t.Error("an app proved on the managed route reads unproved")
	}
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-20", How: "key"}}
	if !unproved() {
		t.Error("a key proof lifted a sign-in's hold")
	}

	// A browser grant, never.
	b := newBrowserEnv(t)
	b.s.Proved = map[string]Proof{}
	if g := b.connect(); b.s.view(nil, g).Unproved {
		t.Error("a browser grant reads unproved")
	}
}
