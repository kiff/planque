package planque

import (
	"context"
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
	if r := s.servedRule("tnt_test", qontoServer); r != "live" {
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
