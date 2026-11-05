package planque

import (
	"context"
	"strings"
	"testing"
)

// keptEntries are the sweep's index of kept sign-ins.
func keptEntries(t *testing.T, e *env) int {
	t.Helper()
	items, err := e.kv.Query(context.Background(), pkKept, "")
	if err != nil {
		t.Fatal(err)
	}
	return len(items)
}

// endWithoutRevoking ends g as a disconnect that saw another grant on its
// sign-in as live would: saved disconnected, its sign-in kept and indexed.
// It is what the second of two disconnects at once leaves.
func endWithoutRevoking(t *testing.T, e *env, id string) {
	t.Helper()
	ctx := context.Background()
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	g.Status, g.Transport, g.Transports = StatusDisconnected, "", nil
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	if err := e.s.markKept(ctx, g); err != nil {
		t.Fatal(err)
	}
}

// Two tafs on one sign-in, both disconnected at once: each saw the other as
// live and kept the sign-in. The sweep revokes it at the vendor and drops
// it once no live grant uses it, and never touches one a live grant uses.
func TestTheSweepRevokesASignInNoLiveGrantUses(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t)
	first := e.connect()
	e.kiff.Cards["k_taf"]["list_charges_b"] = 0
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "list_charges", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	// The first disconnected: kept for the second, and indexed.
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, e.phone, ""); st != 200 {
		t.Fatal("disconnect the first")
	}
	if keptEntries(t, e) != 1 {
		t.Fatalf("%d kept entries, want the first's", keptEntries(t, e))
	}
	// The second is live: the sweep never touches the sign-in, and the
	// entry waits.
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 {
		t.Fatalf("sweep with the second live: %d, %v", n, err)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the sweep revoked a sign-in a live grant uses")
	}
	if !mustHave(mustQuery(t, e), skToken(first.ID)) || keptEntries(t, e) != 1 {
		t.Fatal("the sweep dropped a sign-in a live grant uses, or its entry")
	}
	if res, _ := e.kiff.Call(ctx, "k_taf", "list_charges_b", map[string]any{}); res.Decision != "allowed" {
		t.Fatalf("the second after the sweep: %+v", res)
	} else if text, isErr := resultText(t, res.Result); isErr {
		t.Fatalf("the second after the sweep: %s", text)
	}
	// The second ended at the same moment, seeing the first as live: the
	// sign-in is still live at the vendor with no grant on it.
	endWithoutRevoking(t, e, id)
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("setup: the race left nothing live")
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n == 0 {
		t.Fatalf("sweep with no live grant on it: %d, %v", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the sign-in outlived every connection on it: %d access, %d refresh", a, r)
	}
	if mustHave(mustQuery(t, e), skToken(first.ID)) || keptEntries(t, e) != 0 {
		t.Fatal("the token or an entry outlived the sweep")
	}
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); !strings.Contains(g.Events[len(g.Events)-1].What, "no connection uses it now") {
		t.Fatalf("the first's record: %+v", g.Events[len(g.Events)-1])
	}
	// Nothing left to do.
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 {
		t.Fatalf("a second sweep: %d, %v", n, err)
	}
}

// The same for a provider's connected account: kept while another taf's
// connection uses it, revoked at the provider by the sweep once none does.
func TestTheSweepRevokesAProviderAccountNoLiveGrantUses(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a := connectConnector(t, e, f)
	b := connectConnector(t, e, f)
	if b.Connector.Account != a.Connector.Account {
		b.Connector.Account = a.Connector.Account
		if err := e.s.saveGrant(ctx, b); err != nil {
			t.Fatal(err)
		}
	}
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+a.ID, nil, e.phone, ""); st != 200 {
		t.Fatal("disconnect the first")
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 || f.revocations != 0 {
		t.Fatalf("sweep with the second live: %d, %v, %d revocations", n, err, f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("the other taf's connection: %s", g.Status)
	}
	endWithoutRevoking(t, e, b.ID)
	if n, err := e.s.SweepKept(ctx); err != nil || n != 1 {
		t.Fatalf("sweep with no live grant on it: %d, %v", n, err)
	}
	// Revoked at the provider once, though both grants kept it.
	if f.revocations != 1 {
		t.Fatalf("the provider account: %d revocations, want 1", f.revocations)
	}
	if keptEntries(t, e) != 0 {
		t.Fatal("an entry outlived the sweep")
	}
}

// A grant whose disconnect was indexed but not saved hasn't ended: the
// sweep leaves its sign-in to its own end.
func TestTheSweepLeavesALiveGrantsOwnSignIn(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t)
	g := e.connect()
	if err := e.s.markKept(ctx, g); err != nil {
		t.Fatal(err)
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 {
		t.Fatalf("%d, %v", n, err)
	}
	if a, _ := e.vendor.Live(); a == 0 || !mustHave(mustQuery(t, e), skToken(g.ID)) {
		t.Fatal("the sweep revoked a live grant's own sign-in")
	}
	if keptEntries(t, e) != 0 {
		t.Fatal("the entry for a live grant stayed")
	}
}
