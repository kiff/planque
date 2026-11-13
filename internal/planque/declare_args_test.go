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
	// listCharges takes the account to list in, as a vendor's reads do: what a taf's read pins.
	v.desc = strings.Replace(restDesc, `"parameters": [{"name": "limit"`, `"parameters": [{"name": "account", "in": "query", "schema": {"type": "string"}}, {"name": "limit"`, 1)
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
		"a nested resource":   {map[string]any{"resource_argument": "charge.id"}, "resource_argument charge.id is not a valid name"},
		"a pin it doesn't take": {map[string]any{"fixed_arguments": map[string]any{"account": "a_1"}, "arguments": map[string]any{"charge": map[string]any{"path": "charge", "type": "string"}}},
			"fixed argument account isn't one of createRefund's arguments (amount, charge)"},
		"a path it doesn't take": {map[string]any{"fixed_arguments": map[string]any{"charge": "ch_1"}, "arguments": map[string]any{"cents": map[string]any{"path": "cents", "type": "integer"}}},
			"argument cents's path cents isn't one of createRefund's arguments (amount, charge)"},
		"an argument of another type": {map[string]any{"fixed_arguments": map[string]any{"charge": "ch_1"}, "arguments": map[string]any{"amount": map[string]any{"path": "amount", "type": "string"}}},
			"argument amount's type is string, but createRefund's amount is integer"},
	} {
		st, out := e.askAPIGrant(v, tool(c.extra), nil)
		if st != 400 || !strings.Contains(mustJSONString(out), c.why) {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("a refused request reached KIFF: %+v", c)
	}
	// The operation's own: declared, what the write acts on with it, and a
	// read's pin kept on the grant the key then serves.
	tools := append(tool(map[string]any{"name": "create_refund", "resource_argument": "charge", "amount_argument": "amount", "amount_unit": "EUR"}),
		map[string]any{"tool": "listCharges", "name": "list_charges", "title": "List charges", "read_only": true, "fixed_arguments": map[string]any{"account": "acct_1"},
			"arguments": map[string]any{"limit": map[string]any{"path": "limit", "type": "integer"}}})
	st, out := e.askAPIGrant(v, tools, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("its own arguments: %d %v", st, out)
	}
	if c, ok := e.kiff.Connection(tenant, "create_refund"); !ok || c.ResourceArg != "charge" || c.AmountArg != "amount" {
		t.Fatalf("declared: %+v", c)
	}
	if !strings.Contains(mustJSONString(out), `"fixed_arguments":{"account":"acct_1"}`) {
		t.Fatalf("the read's pin isn't on the grant: %v", out)
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
