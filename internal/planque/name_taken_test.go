package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
)

// kiff/tafs#138: a tool name already taken, by another of the account's
// connections or by any tool in KIFF, is never the owner's problem. KIFF
// would refuse it only at connect, after the owner approved and signed in
// (the owner's iPad, 21 Oct). A tool's name is KIFF's, not the owner's, so
// planque gives it a free one when the agent sets the tools, for any vendor,
// and points the unit lookups that named it there. The grant's own names
// stay its own.
func TestATakenToolNameIsRenamedFirst(t *testing.T) {
	e := newEnv(t)
	first := e.listedGrant()
	// A second connection on the same, already enrolled, phone.
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	second := out["id"].(string)
	if st, out := e.approveAndSignIn(second, out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		t.Fatalf("the second connection: %d %v", st, out)
	}
	at := func(id string) string { return "/v1/accounts/" + tenant + "/grants/" + id + "/tools" }
	names := func(out map[string]any) []string {
		var ns []string
		req, _ := out["request"].(map[string]any)
		tools, _ := req["tools"].([]any)
		for _, tool := range tools {
			ns = append(ns, tool.(map[string]any)["name"].(string))
		}
		return ns
	}
	read := map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "name": "read_payment", "read_only": true}}}
	if st, out := e.do("POST", at(first), read, nil, "k_setup"); st != 200 || names(out)[0] != "read_payment" {
		t.Fatalf("the first connection's read: %d %v", st, out)
	}
	if st, out := e.do("POST", at(first), read, nil, "k_setup"); st != 200 || names(out)[0] != "read_payment" {
		t.Fatalf("its own name again: %d %v", st, out)
	}
	// Another connection asks for the same name: it gets a free one.
	st, out := e.do("POST", at(second), read, nil, "k_setup")
	if st != 200 || names(out)[0] == "read_payment" || !strings.HasPrefix(names(out)[0], "read_payment_") {
		t.Fatalf("the same name on another connection: %d %v", st, out)
	}
	// A tool connected in KIFF outside planque holds its name too, and a
	// lookup that named the taken tool follows it to its new name.
	e.kiff.Add(fakes.Connection{Connection: kiff.Connection{Name: "list_stripe_accounts"}, Tenant: tenant})
	withLookup := map[string]any{"tools": []any{
		map[string]any{"tool": "list_charges", "name": "list_stripe_accounts", "read_only": true},
		map[string]any{"tool": "create_refund", "name": "make_refund", "amount_argument": "amount", "amount_unit": "lookup", "resource_argument": "charge",
			"unit_lookup": map[string]any{"tool": "list_stripe_accounts", "argument": "id", "path": "currency"}},
	}}
	st, out = e.do("POST", at(second), withLookup, nil, "k_setup")
	if st != 200 {
		t.Fatalf("a name KIFF holds: %d %v", st, out)
	}
	got := names(out)
	if got[0] == "list_stripe_accounts" || !strings.HasPrefix(got[0], "list_stripe_accounts_") || got[1] != "make_refund" {
		t.Fatalf("names: %v", got)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, second)
	if err != nil {
		t.Fatal(err)
	}
	if l := g.Request.Tools[1].UnitLookup; l == nil || l.Tool != got[0] {
		t.Fatalf("the lookup names %+v, want %s", l, got[0])
	}
}

// Free names: the grant's suffix, then a number, never one taken or asked
// for in the same request.
func TestRenameTakenPicksAFreeName(t *testing.T) {
	tools := []ToolRequest{{Tool: "a", Name: "read"}, {Tool: "b", Name: "read_ab12"}}
	renameTaken(tools, "grt_ab12cd", map[string]string{"read": "x", "read_ab12": "y"})
	if tools[0].Name != "read_ab122" || tools[1].Name != "read_ab12_ab12" {
		t.Fatalf("renamed: %+v", tools)
	}
}
