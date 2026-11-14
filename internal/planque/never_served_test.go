package planque

import (
	"context"
	"encoding/json"
	"maps"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
)

// neverServe makes the fixture app's tools named never served, beside the
// shipped list, until the test ends. Stripe's own list can't be exercised
// end to end here: a managed Stripe connection isn't offered (directEntryFor).
func neverServe(t *testing.T, slugs ...string) {
	t.Helper()
	was := connector.NeverServedTools
	t.Cleanup(func() { connector.NeverServedTools = was })
	list := maps.Clone(was)
	list["desk_demo"] = slugs
	connector.NeverServedTools = list
}

// A managed app's never-served tool doesn't come back from the tools
// search, nor its catalog, and asking for it, by the narrow sign-in's tools
// or by setting a grant's tools, is refused with tool_not_served. The app's
// other tools, and other apps', are served as before.
func TestAManagedAppsNeverServedToolIsNeitherFoundNorAskedFor(t *testing.T) {
	e, f := newConnectorFixture(t)
	base := "/v1/accounts/" + tenant + "/catalog"

	neverServe(t, "DESK_DEMO_REPLY")
	st, out := e.do("GET", base+"/tools?q=reply", nil, nil, "k_setup")
	if tools := out["tools"].([]any); st != 200 || len(tools) != 1 || tools[0].(map[string]any)["slug"] != "OWN_DEMO_SEND" || out["skipped"].(float64) != 2 {
		t.Fatalf("search: %d %v", st, out)
	}

	neverServe(t, "DESK_DEMO_SEARCH")
	// Negative: the search finds the app's tool not on the list.
	if st, out := e.do("GET", base+"/tools?q=reply", nil, nil, "k_setup"); st != 200 || len(out["tools"].([]any)) != 2 {
		t.Fatalf("search, served: %d %v", st, out)
	}
	ask := func(tools ...string) (int, map[string]any) {
		return e.do("POST", base+"/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "new_sign_in": true, "tools": tools}, nil, "k_setup")
	}
	st, out = ask("DESK_DEMO_REPLY", "desk_demo_search")
	if st != 400 || out["code"] != refToolNotServed || out["tool"] != "DESK_DEMO_SEARCH" || out["field"] != "tools" || !strings.Contains(out["error"].(string), "never served") {
		t.Fatalf("asked with a never-served tool: %d %v", st, out)
	}
	// Negative: the app's other tool is asked for.
	if st, out := ask("DESK_DEMO_REPLY"); st != 201 {
		t.Fatalf("asked with a served tool: %d %v", st, out)
	}

	f.search = true
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	tools := catalog["tools"].([]any)
	if st != 200 || len(tools) != 1 || tools[0].(map[string]any)["title"] != "Reply to a ticket" || catalog["skipped"].(float64) != 1 || len(catalog["unusable"].([]any)) != 0 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	hidden := connector.Tool{Toolkit: "desk_demo", Slug: "DESK_DEMO_SEARCH"}.Alias()
	st, out = e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: hidden, Title: "Search tickets", Name: "search_tickets", ReadOnly: true}}}, nil, "k_setup")
	if st != 400 || out["code"] != refToolNotServed || out["tool"] != "DESK_DEMO_SEARCH" || out["field"] != "tool" {
		t.Fatalf("set a never-served tool: %d %v", st, out)
	}
}

// A connection that holds a tool from before planque stopped serving it
// never sends a call to it, and no longer lists it; until then it runs.
func TestACallToANeverServedToolIsNotSent(t *testing.T) {
	e, f := newConnectorFixture(t)
	neverServe(t, "DESK_DEMO_SEARCH")
	g := connectConnector(t, e, f)
	ctx := context.Background()
	call := func(op string) string {
		t.Helper()
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: op, Args: args})
		return mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	}
	// Negative: another tool of the app on the list leaves this one served.
	if got := call("op_served"); !strings.Contains(got, "reply_demo") || f.executions != 1 {
		t.Fatalf("a served tool: %s", got)
	}
	if listed, _, err := e.s.connectorTools(ctx, g); err != nil || len(listed) != 1 {
		t.Fatalf("listed: %v %v", listed, err)
	}

	neverServe(t, "DESK_DEMO_REPLY")
	got := call("op_never")
	if !strings.Contains(got, `"code":"`+CodeToolNotServed+`"`) || !strings.Contains(got, `"sent":false`) || !strings.Contains(got, "Not sent: planque never serves DESK_DEMO_REPLY") || f.executions != 1 {
		t.Fatalf("a never-served tool: %s, executions %d", got, f.executions)
	}
	if listed, _, err := e.s.connectorTools(ctx, g); err != nil || len(listed) != 0 {
		t.Fatalf("still listed: %v %v", listed, err)
	}
}

// A connection that holds a tool from before planque stopped serving it
// takes no new tool until that one is removed: adding any is refused with
// tool_not_served, naming the held tool, and the connection stays as it was.
// Category: money guard
// Attack: adding a served tool to re-bind a held never-served one
// Refused: 400 tool_not_served; nothing bound
func TestAddingToAConnectionHoldingANeverServedToolIsRefused(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	g := connectConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID
	f.search = true
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	search := connector.Tool{Toolkit: "desk_demo", Slug: "DESK_DEMO_SEARCH"}.Alias()
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "search_tickets", Title: "Search tickets", Kind: "write", Effect: profile.EffectSends, Shape: "single"}}}
	add := func() (int, map[string]any) {
		return e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: search, Title: "Search tickets", Name: "search_tickets"}}, "profile": p}, nil, "k_setup")
	}

	neverServe(t, "DESK_DEMO_REPLY")
	st, out := add()
	if st != 400 || out["code"] != refToolNotServed || out["tool"] != "DESK_DEMO_REPLY" || out["field"] != "tool" {
		t.Fatalf("added beside a held never-served tool: %d %v", st, out)
	}
	if held, _ := e.s.loadGrant(ctx, tenant, g.ID); held.Status != StatusConnected || len(held.Request.Tools) != 1 || len(held.Request.Connector.Tools) != 1 {
		t.Fatalf("changed by a refused addition: %s %v", held.Status, held.Request.Tools)
	}

	// Negative: with the held tool served again, the same addition is taken.
	neverServe(t)
	if st, out := add(); st != 200 || out["status"] != StatusPending {
		t.Fatalf("added beside served tools: %d %v", st, out)
	}
}
