package planque

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kiff/planque/internal/kiff"
)

// Ending an agent's grants when its taf is removed: a route that only takes
// away, POST /v1/accounts/{tenant}/agents/{agent}/end, called with the
// person's own KIFF session, which must manage the account's tafs
// (cards.manage). Every grant of that agent ends, whatever its status; a
// sign-in another agent's grant still uses stays; it is idempotent; with
// ?preview=1 it ends nothing and answers what it would end and keep. No
// service key ends anyone's connections, and another account's session ends
// nothing. These are the negatives: what must not survive a removed taf,
// and what must not be taken from anyone else.

// people adds the sessions the route is called with: a person who manages
// the account's tafs, one who only reads, and one who manages another
// account's.
func (e *env) people() {
	e.kiff.Keys["s_owner"] = kiff.Me{TenantID: tenant, Capabilities: []string{"read", "tools.connect", "cards.manage"}}
	e.kiff.Keys["s_viewer"] = kiff.Me{TenantID: tenant, Capabilities: []string{"read"}}
	e.kiff.Keys["s_other_owner"] = kiff.Me{TenantID: "tnt_other", Capabilities: []string{"read", "tools.connect", "cards.manage"}}
}

func (e *env) endAgent(tenantID, agent, key, query string) (int, map[string]any) {
	e.t.Helper()
	return e.do("POST", "/v1/accounts/"+tenantID+"/agents/"+agent+"/end"+query, map[string]any{}, nil, key)
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

// removedTaf gives agt_taf a failed, a pending and a connected grant, and
// another agent a pending one.
func (e *env) removedTaf() (failed, pending string, connected *Grant, other string) {
	e.t.Helper()
	e.enrol()
	e.people()
	// A failed one: a sign-in whose tools nobody approved within a day.
	failed = e.listedGrant()
	e.later(holdTTL + time.Minute)
	if s := e.status(failed); s != StatusFailed {
		e.t.Fatalf("setup: %s, want failed", s)
	}
	pending, _ = e.request()
	connected = e.connect()
	other = e.requestFor("agt_other")
	return
}

func TestEndingARemovedTafsAgentEndsEveryGrantItHas(t *testing.T) {
	e := newEnv(t)
	failed, pending, connected, other := e.removedTaf()

	st, out := e.endAgent(tenant, "agt_taf", "s_owner", "")
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
	if st, out := e.endAgent(tenant, "agt_taf", "s_owner", ""); st != 200 {
		t.Fatalf("end again: %d %v", st, out)
	}
	if s := e.status(other); s != StatusPending {
		t.Errorf("another agent's grant after a second end: %s", s)
	}
}

// The preview ends nothing, revokes nothing, and answers what the call
// that follows then does.
func TestThePreviewOfEndingAnAgentEndsNothing(t *testing.T) {
	e := newEnv(t)
	failed, pending, connected, _ := e.removedTaf()
	before := map[string]string{}
	for _, id := range []string{failed, pending, connected.ID} {
		before[id] = e.status(id)
	}
	liveA, liveR := e.vendor.Live()

	st, preview := e.endAgent(tenant, "agt_taf", "s_owner", "?preview=1")
	if st != 200 {
		t.Fatalf("preview: %d %v", st, preview)
	}
	for id, was := range before {
		if s := e.status(id); s != was {
			t.Errorf("the preview changed %s from %s to %s", id, was, s)
		}
	}
	if a, r := e.vendor.Live(); a != liveA || r != liveR {
		t.Errorf("the preview revoked tokens: %d, %d (was %d, %d)", a, r, liveA, liveR)
	}
	if code := e.relayCallRaw(connected); code != 200 {
		t.Errorf("after a preview the connection relays %d, want it still working", code)
	}

	st, done := e.endAgent(tenant, "agt_taf", "s_owner", "")
	if st != 200 {
		t.Fatalf("end: %d %v", st, done)
	}
	list := func(v any) []any { l, _ := v.([]any); return l }
	if len(list(preview["ended"])) == 0 || len(list(preview["ended"])) != len(list(done["ended"])) {
		t.Errorf("the preview said it would end %v; the call ended %v", preview["ended"], done["ended"])
	}
	p, _ := json.Marshal(preview["kept"])
	d, _ := json.Marshal(done["kept"])
	if string(p) != string(d) {
		t.Errorf("the preview said it would keep %s; the call kept %s", p, d)
	}
}

// A provider's sign-in that another taf's connection uses stays signed in
// when one agent ends, and the preview says so; it is revoked when the
// last agent using it ends.
func TestEndingAnAgentKeepsASignInAnotherTafUses(t *testing.T) {
	e, f := newConnectorFixture(t)
	e.people()
	ctx := context.Background()
	a := connectConnector(t, e, f)
	b := connectConnector(t, e, f)
	b.Request.Agent = "agt_refunds"
	b.Connector.Account = a.Connector.Account
	if err := e.s.saveGrant(ctx, b); err != nil {
		t.Fatal(err)
	}
	if st, out := e.endAgent(tenant, a.Request.Agent, "s_owner", "?preview=1"); st != 200 || len(out["kept"].([]any)) == 0 {
		t.Fatalf("the preview doesn't say the shared sign-in stays: %d %v", st, out)
	}
	if st, out := e.endAgent(tenant, a.Request.Agent, "s_owner", ""); st != 200 {
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
	if st, out := e.endAgent(tenant, "agt_refunds", "s_owner", ""); st != 200 {
		t.Fatalf("end the second agent: %d %v", st, out)
	}
	if f.revocations != 1 {
		t.Fatalf("the last taf on the sign-in ended without revoking it: %d revocations", f.revocations)
	}
}

// Only a person who manages this account's tafs ends an agent: no service
// key (Setup's, a taf's own, an admin key), no person who only reads, no
// one from another account by either path, and no call without a key.
// Each refusal ends nothing, preview or not.
func TestOnlyThePersonWhoManagesTheTafsEndsAnAgent(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.people()
	pending, _ := e.request()
	connected := e.connect()
	for _, c := range []struct{ name, path, key string }{
		{"Setup's key", tenant, "k_setup"},
		{"the taf's own key", tenant, "k_taf"},
		{"an admin API key", tenant, "k_admin"},
		{"a person who only reads", tenant, "s_viewer"},
		{"another account's owner, on this account", tenant, "s_other_owner"},
		{"this account's owner, on another account's path", "tnt_other", "s_owner"},
		{"another account's key", tenant, "k_other"},
		{"no key", tenant, ""},
	} {
		for _, q := range []string{"", "?preview=1"} {
			if st, out := e.endAgent(c.path, "agt_taf", c.key, q); st == 200 {
				t.Errorf("%s%s ended an agent: %d %v", c.name, q, st, out)
			}
		}
	}
	if s := e.status(pending); s != StatusPending {
		t.Errorf("the pending grant after the refused calls: %s", s)
	}
	if s := e.status(connected.ID); s != StatusConnected {
		t.Errorf("the connected grant after the refused calls: %s", s)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Error("a refused call revoked a token")
	}
}
