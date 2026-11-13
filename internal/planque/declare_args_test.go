package planque

import (
	"strings"
	"testing"
)

// What KIFF scopes a Card on (resource_argument) and counts (amount_argument)
// must be one of the operation's own arguments, of its type, when the tool
// isn't narrowed: checked at the request, with planque's words, before
// anything is declared. Not, as in the Eval, declared as is and refused by
// KIFF as a 502 when the key comes.
func TestAKeyToolsResourceAndAmountAreTheOperationsOwn(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	tool := func(extra map[string]any) []map[string]any {
		t := map[string]any{"tool": "createRefund", "title": "Refund a charge"}
		for k, x := range extra {
			t[k] = x
		}
		return []map[string]any{t}
	}
	for name, c := range map[string]struct {
		extra map[string]any
		why   string
	}{
		"not an argument":     {map[string]any{"resource_argument": "order_id"}, "resource_argument order_id must be one of createRefund's string arguments (charge)"},
		"an integer resource": {map[string]any{"resource_argument": "amount"}, "resource_argument amount must be one of createRefund's string arguments (charge)"},
		"a string amount":     {map[string]any{"amount_argument": "charge", "amount_unit": "EUR"}, "amount_argument charge must be one of createRefund's integer arguments (amount)"},
	} {
		st, out := e.askAPIGrant(v, tool(c.extra), nil)
		if st != 400 || !strings.Contains(mustJSONString(out), c.why) {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
	// The operation's own: declared.
	st, out := e.askAPIGrant(v, tool(map[string]any{"resource_argument": "charge", "amount_argument": "amount", "amount_unit": "EUR"}), nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("its own arguments: %d %v", st, out)
	}
}

// describe answers each operation's arguments, so Coach points a tool's
// resource and amount at ones the request takes, and pins them, before any
// grant: what it answers is what the check above holds a request to.
func TestDescribeGivesEachOperationsArguments(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	api := map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/api/operations", map[string]any{"api": api, "q": "refund charge"}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("%d %v", st, out)
	}
	var args []any
	for _, o := range out["operations"].([]any) {
		if m := o.(map[string]any); m["name"] == "createRefund" {
			args = m["arguments"].([]any)
		}
	}
	got := mustJSONString(args)
	want := `[{"in":"path","name":"charge","required":true,"type":"string"},{"in":"body","name":"amount","required":true,"type":"integer"}]`
	if got != want {
		t.Fatalf("createRefund's arguments:\n%s\nwant\n%s", got, want)
	}
}
