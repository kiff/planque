package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
)

// An app taken off the published list closes the connections made while it was
// on: their calls and the tools added to them are refused ("Nothing was sent"),
// where before only a new connection was. Put back, the same calls go through,
// and the test accounts' own try-anything exemption stays as it was.
func TestAnAppTakenOffTheListClosesItsConnections(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	call := func(op string) string {
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: op, Args: args})
		return mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	}
	closeIt := func() {
		e.s.ManagedApps, e.s.ManagedTestTenants = map[string]bool{}, map[string]bool{}
	}
	openIt := func() { e.s.ManagedApps = map[string]bool{"deskdemo": true} }

	closeIt()
	before := f.executions
	if got := call("operation_closed"); !strings.Contains(got, "connection is closed") || f.executions != before {
		t.Fatalf("a call on a closed app's connection: %s (executions %d, was %d)", got, f.executions, before)
	}

	// The tools added to a connected grant are refused at their approval.
	f.search = true
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID
	search := connector.Tool{Toolkit: "desk_demo", Slug: "DESK_DEMO_SEARCH"}.Alias()
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "search_tickets", Title: "Search tickets", Kind: "write", Effect: profile.EffectSends, Shape: "single"}}}
	openIt()
	if st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup"); st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: search, Title: "Search tickets", Name: "search_tickets"}}, "profile": p}, nil, "k_setup"); st != 200 {
		t.Fatalf("add tools: %d %v", st, out)
	}
	pending, _ := e.s.loadGrant(ctx, tenant, g.ID)
	ph, _ := profile.Digest(pending.Request.Profile)
	closeIt()
	st, out := e.do("POST", path+"/approve", map[string]any{"digest": pending.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 403 || out["error"] != errAppClosed {
		t.Fatalf("approving tools added to a closed app's connection: %d %v", st, out)
	}

	// On the list again, the same call goes through: the switch is the list.
	openIt()
	before = f.executions
	if got := call("operation_open"); strings.Contains(got, "connection is closed") || f.executions != before+1 {
		t.Fatalf("a call with the app listed: %s (executions %d, was %d)", got, f.executions, before)
	}
}

// An account moved to an app's own client keeps the managed connection it made
// before: the per-account choice is about what's offered next, so a call on the
// earlier connection goes through while a new managed request is refused.
func TestAnOwnClientSwitchDoesNotCloseAnEarlierManagedConnection(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	e.s.ManagedApps, e.s.ManagedTestTenants = map[string]bool{"deskdemo": true}, map[string]bool{}
	e.s.CatalogEntries = map[string]CatalogEntry{"desk@1": {Application: "desk_demo", Vendor: "Support desk", Title: "Support desk"}}
	e.s.OwnClientTenants = map[string]bool{tenant: true}
	if e.s.managedAllowed(tenant, "desk_demo") {
		t.Fatal("a new managed connection is still offered")
	}
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"operation_own"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "operation_own", Args: args})
	before := f.executions
	if got := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); strings.Contains(got, "connection is closed") || f.executions != before+1 {
		t.Fatalf("a call on the earlier managed connection: %s (executions %d, was %d)", got, f.executions, before)
	}
}
