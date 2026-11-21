package planque

import (
	"context"
	"strings"
	"testing"
)

// A phone test on 21 Nov 2026: Coach reused a connected Stripe grant, was
// told "this grant has no catalog", and named the write tool from memory
// (stripe_api_execute, which Stripe's server no longer lists). The connection
// failed only after the owner approved it. A connected MCP grant's catalog is
// now its server's list, read with its own sign-in, and nothing is kept.
func TestAConnectedServersCatalogIsItsListNow(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID
	st, out := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog of a connected grant: %d %v", st, out)
	}
	names := map[string]map[string]any{}
	for _, x := range out["tools"].([]any) {
		c := x.(map[string]any)
		names[c["name"].(string)] = c
	}
	for _, want := range []string{"create_refund", "list_charges", "stripe_api_write", "stripe_api_details"} {
		if names[want] == nil {
			t.Fatalf("%s not in the catalog: %v", want, out["tools"])
		}
	}
	if names["stripe_api_execute"] != nil {
		t.Fatal("a tool the server doesn't list is in the catalog")
	}
	if names["list_charges"]["read_only"] != true || names["create_refund"]["read_only"] == true {
		t.Fatalf("read-only marks: %v %v", names["list_charges"], names["create_refund"])
	}
	if names["create_refund"]["input_schema"] == nil {
		t.Fatal("a tool's input schema isn't given")
	}
	after, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if len(after.Catalog) != 0 || after.Status != StatusConnected || len(after.Tools) != len(g.Tools) {
		t.Fatalf("reading the catalog changed the grant: %s, %d kept, %d tools", after.Status, len(after.Catalog), len(after.Tools))
	}
	// Another agent's key may not read it, as before.
	if st, _ := e.do("GET", path+"/catalog", nil, nil, "k_other"); st != 403 {
		t.Fatalf("another agent read the catalog: %d", st)
	}
}

// The same phone test's request: reusing that sign-in, a tool named from
// memory is refused when it's asked, with the names the server lists, before
// the owner sees an approval. A name the server lists still goes on.
func TestAReusedServersToolItDoesntListIsRefusedWhenAsked(t *testing.T) {
	e := newEnv(t)
	first := e.connect()
	ask := func(tool string) (int, map[string]any) {
		return e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_second",
			"tools": []any{map[string]any{"tool": tool, "name": "list_payments", "title": "List payments", "read_only": true}}}, nil, "k_setup")
	}
	st, out := ask("stripe_api_execute")
	if st != 400 || out["code"] != refNotListed || out["field"] != "tool" || !strings.Contains(out["error"].(string), "has no tool named stripe_api_execute now; it lists: ") ||
		!strings.Contains(out["error"].(string), "stripe_api_read") || !strings.Contains(out["error"].(string), "Read grant "+first.ID+"'s catalog") {
		t.Fatalf("an unlisted tool: %d %v", st, out)
	}
	if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "GRANT#"); len(items) != 1 {
		t.Fatalf("a refused request made a grant: %d grants", len(items))
	}
	if st, out := ask("list_charges"); st != 201 {
		t.Fatalf("a listed tool: %d %v", st, out)
	}
}
