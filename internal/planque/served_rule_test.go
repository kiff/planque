package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
)

// Step 9 for the test accounts: once PLANQUE_LIVE_LIST=test is set, a test
// account's grants on Qonto's sign-in get Qonto's whole live list, and KIFF's
// Card decides each call. Every other account keeps qontoServed, and so does a
// test account until it's set.
func TestTheTestAccountsGetTheWholeLiveListOnceSwitchedOn(t *testing.T) {
	s := &Server{ManagedTestTenants: map[string]bool{"tnt_test": true}}
	listed := []mcpTool{{Name: "list_clients"}, {Name: "send_client_invoice"}, {Name: "create_transfer_request"}}
	check := func(tenant string, want bool) {
		t.Helper()
		if got := s.served(tenant, qontoServer, "send_client_invoice"); got != want {
			t.Fatalf("%s, live list %v: send_client_invoice served %v", tenant, s.LiveListTest, got)
		}
		if got := len(s.servedFor(tenant, qontoServer, listed)); (got == len(listed)) != want {
			t.Fatalf("%s, live list %v: %d of %d listed", tenant, s.LiveListTest, got, len(listed))
		}
		if r := s.toolsCheck(tenant, qontoServer, []ToolRequest{{Tool: "send_client_invoice"}}, "tools"); (r == nil) != want {
			t.Fatalf("%s, live list %v: refusal %+v", tenant, s.LiveListTest, r)
		}
	}
	check("tnt_test", false)
	check("tnt_any", false)
	s.LiveListTest = true
	check("tnt_test", true)
	check("tnt_any", false)
	if r := s.servedRule("tnt_test", qontoServer); !strings.HasPrefix(r, "live:") {
		t.Fatalf("a test account's rule: %q", r)
	}
	if r := s.servedRule("tnt_any", qontoServer); !strings.HasPrefix(r, "qonto:") {
		t.Fatalf("another account's rule on Qonto: %q", r)
	}
	if r := s.servedRule("tnt_test", "https://mcp.vendor.example/mcp"); r != "" {
		t.Fatalf("a server planque doesn't filter: %q", r)
	}
}

// A grant listed before a change to what planque serves keeps the catalog it
// was listed with: the Qonto re-test found list_bank_accounts missing from a
// sign-in made an hour before it was served. A catalog read that finds an
// earlier rule lists the tools again with the grant's own sign-in, keeps the
// new catalog and answers it; the tools the owner approved don't change. A
// listing that fails is refused, never answered stale.
func TestACatalogKeptUnderAnEarlierRuleIsListedAgainOnItsNextRead(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	ctx := context.Background()
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	all := len(g.Catalog)
	if all < 2 {
		t.Fatalf("the fake vendor lists %d tools", all)
	}
	approved := mustJSONString(g.Tools)
	// As kept under an earlier rule: a smaller list, and its stamp.
	g.Catalog, g.CatalogRule = g.Catalog[:1], "qonto:000000000000"
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id+"/catalog", nil, nil, "k_setup")
	if st != 200 || len(out["tools"].([]any)) != all {
		t.Fatalf("a catalog under an earlier rule: %d %v", st, out)
	}
	g, err = e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Catalog) != all || g.CatalogRule != e.s.servedRule(tenant, g.Request.Server) {
		t.Fatalf("kept %d tools under %q", len(g.Catalog), g.CatalogRule)
	}
	if mustJSONString(g.Tools) != approved {
		t.Fatalf("the approved tools changed: %s, were %s", mustJSONString(g.Tools), approved)
	}

	// The vendor can't be reached: refused, never the stale list.
	g.Catalog, g.CatalogRule = g.Catalog[:1], "qonto:000000000000"
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	e.vendor.Close()
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id+"/catalog", nil, nil, "k_setup"); st != 503 || out["tools"] != nil {
		t.Fatalf("a re-listing that failed: %d %v", st, out)
	}
}

// qontoAtFake sends what planque asks of Qonto's MCP host to the test's fake
// vendor, so a grant on Qonto's server name lists the tools the fake lists.
type qontoAtFake struct{ fake *url.URL }

func (q qontoAtFake) RoundTrip(r *http.Request) (*http.Response, error) {
	if onHost(r.URL.String(), qontoHost) {
		r = r.Clone(r.Context())
		r.URL.Scheme, r.URL.Host, r.Host = q.fake.Scheme, q.fake.Host, q.fake.Host
	}
	return http.DefaultTransport.RoundTrip(r)
}

// A listed grant on Qonto's own host, with Qonto's tool names, as listed
// before a change to what planque serves.
func qontoListedGrant(t *testing.T, e *env) *Grant {
	t.Helper()
	e.vendor.ChangeTools = func([]map[string]any) []map[string]any {
		tool := func(name string, readOnly bool) map[string]any {
			return map[string]any{"name": name, "description": name, "inputSchema": map[string]any{"type": "object"},
				"annotations": map[string]any{"readOnlyHint": readOnly}}
		}
		return []map[string]any{tool("list_clients", true), tool("list_bank_accounts", true), tool("send_client_invoice", false)}
	}
	fake, err := url.Parse(e.vendor.MCPURL())
	if err != nil {
		t.Fatal(err)
	}
	e.s.MCP.HTTP = &http.Client{Transport: qontoAtFake{fake}}
	id := e.listedGrant()
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	g.Request.Server = qontoServer
	g.Catalog, g.CatalogRule = []CatalogTool{{Name: "list_clients"}}, "qonto:000000000000"
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	return g
}

func catalogNames(t *testing.T, e *env, id string) []string {
	t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	var names []string
	for _, x := range out["tools"].([]any) {
		names = append(names, x.(map[string]any)["name"].(string))
	}
	slices.Sort(names)
	return names
}

// On Qonto's own host: a sign-in listed before list_bank_accounts was served
// offers it on its next catalog read, and still never send_client_invoice;
// with the live list switched on for a test account, it offers all three,
// stamped "live:…" (the whole list, less the bank's own approvals).
func TestAQontoCatalogIsListedAgainUnderTheRuleOfTheDay(t *testing.T) {
	e := newEnv(t)
	g := qontoListedGrant(t, e)
	if got := catalogNames(t, e, g.ID); !slices.Equal(got, []string{"list_bank_accounts", "list_clients"}) {
		t.Fatalf("under qontoServed: %v", got)
	}
	if g, _ = e.s.loadGrant(context.Background(), tenant, g.ID); !strings.HasPrefix(g.CatalogRule, "qonto:") || g.CatalogRule == "qonto:000000000000" {
		t.Fatalf("kept under %q", g.CatalogRule)
	}
	e.s.ManagedTestTenants, e.s.LiveListTest = map[string]bool{tenant: true}, true
	if got := catalogNames(t, e, g.ID); !slices.Equal(got, []string{"list_bank_accounts", "list_clients", "send_client_invoice"}) {
		t.Fatalf("the live list for a test account: %v", got)
	}
	if g, _ = e.s.loadGrant(context.Background(), tenant, g.ID); !strings.HasPrefix(g.CatalogRule, "live:") {
		t.Fatalf("kept under %q", g.CatalogRule)
	}
}

// Add-tools chooses from the catalog of the day too: a stale catalog is
// listed again before the tools asked for are checked against it.
func TestAddToolsChoosesFromTheCatalogOfTheDay(t *testing.T) {
	e := newEnv(t)
	g := qontoListedGrant(t, e)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/tools",
		map[string]any{"tools": []any{map[string]any{"tool": "list_bank_accounts", "title": "List bank accounts", "read_only": true}}}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("a tool served since the grant was listed: %d %v", st, out)
	}
}

// A signed-out grant whose catalog is stale answers as a signed-out read
// would, never 503 "try again", which wouldn't help.
func TestAStaleCatalogOnASignedOutGrantSaysItsSignedOut(t *testing.T) {
	e := newEnv(t)
	g := qontoListedGrant(t, e)
	g.Reconnect = true
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/catalog", nil, nil, "k_setup"); st == 200 || st == 503 {
		t.Fatalf("a signed-out grant's stale catalog: %d %v", st, out)
	}
}

// Turned off again: a tool approved under the live list stays approved but
// is no longer served, so nothing is sent to it.
func TestTurningTheLiveListOffStopsWhatOnlyItServed(t *testing.T) {
	e := newEnv(t)
	e.s.ManagedTestTenants, e.s.LiveListTest = map[string]bool{tenant: true}, true
	if !e.s.served(tenant, qontoServer, "send_client_invoice") {
		t.Fatal("not served with the live list on")
	}
	e.s.LiveListTest = false
	if e.s.served(tenant, qontoServer, "send_client_invoice") {
		t.Fatal("still served with the live list off")
	}
	if _, err := e.s.callServed(context.Background(), tenant, qontoServer, "token", "send_client_invoice", json.RawMessage(`{}`)); !errors.Is(err, errNotServed) {
		t.Fatalf("sent with the live list off: %v", err)
	}
}
