package planque

import (
	"context"
	"testing"
	"time"
)

// Ending an agent's grants when its taf is removed: a route that only takes
// away, POST /v1/accounts/{tenant}/agents/{agent}/end, called with the
// account's Setup key. Every grant of that agent ends, whatever its status;
// a sign-in another agent's grant still uses stays; it is idempotent; and
// another account's key ends nothing. These are the negatives for it: what
// must not survive a removed taf, and what must not be taken from another.

func (e *env) endAgent(tenantID, agent, key string) (int, map[string]any) {
	e.t.Helper()
	return e.do("POST", "/v1/accounts/"+tenantID+"/agents/"+agent+"/end", map[string]any{}, nil, key)
}

// requestFor asks, as Setup, for a pending grant for agent.
func (e *env) requestFor(agent string) string {
	e.t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": agent,
		"tools": []map[string]any{{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		e.t.Fatalf("request for %s: %d %v", agent, st, out)
	}
	return out["id"].(string)
}

// live is every status in which a grant can still be approved, signed in
// or used.
var live = map[string]bool{StatusPending: true, StatusApproved: true, StatusSignedIn: true, StatusListed: true,
	StatusConnected: true, StatusNeedsReconnect: true}

func TestEndingARemovedTafsAgentEndsEveryGrantItHas(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	// A failed one: a sign-in whose tools nobody approved within a day.
	failed := e.listedGrant()
	e.later(holdTTL + time.Minute)
	if s := e.status(failed); s != StatusFailed {
		t.Fatalf("setup: %s, want failed", s)
	}
	pending, _ := e.request()
	connected := e.connect()
	other := e.requestFor("agt_other")

	st, out := e.endAgent(tenant, "agt_taf", "k_setup")
	if st != 200 {
		t.Fatalf("end: %d %v", st, out)
	}
	for name, id := range map[string]string{"pending": pending, "failed": failed, "connected": connected.ID} {
		if s := e.status(id); live[s] {
			t.Errorf("the removed taf's %s grant is still %s", name, s)
		}
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Errorf("tokens still live at the vendor for a removed taf: %d, %d", a, r)
	}
	if code := e.relayCallRaw(connected); code == 200 {
		t.Error("the removed taf's connection still relays a call")
	}
	if s := e.status(other); s != StatusPending {
		t.Errorf("another agent's grant: %s, want it left pending", s)
	}

	// Again: nothing more to end, nothing changed, no error.
	if st, out := e.endAgent(tenant, "agt_taf", "k_setup"); st != 200 {
		t.Fatalf("end again: %d %v", st, out)
	}
	if s := e.status(other); s != StatusPending {
		t.Errorf("another agent's grant after a second end: %s", s)
	}
}

// A provider's sign-in that another taf's connection uses stays signed in
// when one agent ends; it is revoked when the last agent using it ends.
func TestEndingAnAgentKeepsASignInAnotherTafUses(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a := connectConnector(t, e, f)
	b := connectConnector(t, e, f)
	b.Request.Agent = "agt_refunds"
	if b.Connector.Account != a.Connector.Account {
		b.Connector.Account = a.Connector.Account
	}
	if err := e.s.saveGrant(ctx, b); err != nil {
		t.Fatal(err)
	}
	if st, out := e.endAgent(tenant, a.Request.Agent, "k_setup"); st != 200 {
		t.Fatalf("end the first agent: %d %v", st, out)
	}
	if f.revocations != 0 {
		t.Fatalf("ending one taf revoked the sign-in another uses: %d revocations", f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("the other taf's connection: %s, want connected", g.Status)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, a.ID); live[g.Status] {
		t.Fatalf("the ended agent's connection is still %s", g.Status)
	}
	if st, out := e.endAgent(tenant, "agt_refunds", "k_setup"); st != 200 {
		t.Fatalf("end the second agent: %d %v", st, out)
	}
	if f.revocations != 1 {
		t.Fatalf("the last taf on the sign-in ended without revoking it: %d revocations", f.revocations)
	}
}

// Another account's Setup key ends nothing here, through its own account's
// path or this one's; and a request with no key ends nothing.
func TestAnotherAccountsKeyEndsNoGrant(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	pending, _ := e.request()
	connected := e.connect()
	for _, c := range []struct{ path, key string }{
		{tenant, "k_other"},
		{"tnt_other", "k_setup"},
		{tenant, ""},
	} {
		if st, out := e.endAgent(c.path, "agt_taf", c.key); st == 200 {
			t.Errorf("end on %s with %q: %d %v", c.path, c.key, st, out)
		}
	}
	if s := e.status(pending); s != StatusPending {
		t.Errorf("the pending grant after another account's end: %s", s)
	}
	if s := e.status(connected.ID); s != StatusConnected {
		t.Errorf("the connected grant after another account's end: %s", s)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Error("another account's end revoked a token here")
	}
}
