package planque

import (
	"context"
	"testing"
)

// A newer request that KIFF turns down leaves the earlier one waiting: the
// owner may be about to approve it. One that stands ends the earlier one and
// takes the tools declared for it out of KIFF. One that reuses another grant's
// sign-in isn't ended by a newer ask.
func TestEndingAnEarlierRequestFollowsTheNewOneStanding(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	// A key grant declares its tools in KIFF at its request; naming its app lets a newer ask find it.
	e.s.ServerApps = map[string]string{v.URL: "connector:composio:refunds"}
	ctx := context.Background()
	ask := func() (int, map[string]any) { return e.askAPIGrant(v, apiTools, nil) }
	load := func(out map[string]any) *Grant {
		g, err := e.s.loadGrant(ctx, tenant, out["id"].(string))
		if err != nil {
			t.Fatal(err)
		}
		return g
	}
	held := func(g *Grant) int {
		n := 0
		for _, tool := range g.Tools {
			if _, ok := e.kiff.Connection(tenant, tool.Name); ok {
				n++
			}
		}
		return n
	}

	st, out := ask()
	if st != 201 {
		t.Fatalf("first: %d %v", st, out)
	}
	first := load(out)
	if !first.Declared || held(first) == 0 {
		t.Fatalf("the first request's tools aren't declared: %v", first.Declared)
	}

	// KIFF turns the newer one down: the earlier one is still waiting, its tools still in KIFF.
	e.kiff.RefuseConnectAfter = len(e.kiff.Connections(tenant))
	if st, out := ask(); st != 400 {
		t.Fatalf("a refused request: %d %v", st, out)
	}
	e.kiff.RefuseConnectAfter = 0
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusPending || held(g) != len(g.Tools) {
		t.Fatalf("the earlier request after a refused one: %s, %d of %d tools in KIFF", g.Status, held(g), len(g.Tools))
	}

	// One that holds another grant's sign-in isn't ended by a newer ask.
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	g.TokenOf = "grt_elsewhere"
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	st, out = ask()
	if st != 201 {
		t.Fatalf("second: %d %v", st, out)
	}
	second := load(out)
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusPending {
		t.Fatalf("a request holding another's sign-in was ended: %s", g.Status)
	}
	g, _ = e.s.loadGrant(ctx, tenant, first.ID)
	g.TokenOf = ""
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}

	// A newer one that stands ends the earlier ones, and their declared tools leave KIFF first.
	st, out = ask()
	if st != 201 {
		t.Fatalf("third: %d %v", st, out)
	}
	third := load(out)
	for _, old := range []*Grant{first, second} {
		g, _ := e.s.loadGrant(ctx, tenant, old.ID)
		if g.Status != StatusDisconnected || g.Declared || held(old) != 0 {
			t.Errorf("an earlier request: %s, declared %v, %d tools still in KIFF", g.Status, g.Declared, held(old))
		}
	}
	if g, _ := e.s.loadGrant(ctx, tenant, third.ID); g.Status != StatusPending || held(third) != len(third.Tools) || len(third.Tools) == 0 {
		t.Errorf("the newest: %s, %d of %d tools in KIFF", g.Status, held(third), len(third.Tools))
	}
}
