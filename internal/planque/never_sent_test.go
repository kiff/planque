package planque

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// The operations Stripe's broad write is never narrowed to, and a call never
// sends, are those it lists, each an operation in Stripe's published API.
func TestStripesNeverSentOperations(t *testing.T) {
	want := []string{"PostInvoicesInvoiceVoid", "PostInvoicesInvoicePay", "PostInvoicesInvoiceMarkUncollectible", "DeleteInvoicesInvoice",
		"PostInvoicesInvoiceFinalize", "PostInvoicesInvoiceAttachPayment", "PostCreditNotesIdVoid", "PostBillingCreditGrantsIdVoid"}
	if got := describers[0].Never; !slices.Equal(got, want) {
		t.Fatalf("never: %v", got)
	}
	for _, op := range append(want, " postinvoicesinvoicevoid ", "POSTINVOICESINVOICEPAY") {
		if !neverSent("stripe_api_write", op) {
			t.Errorf("%q is sent", op)
		}
	}
	// Negative: refunds, the send, reads, and another vendor's tool of the same operation.
	for tool, op := range map[string]string{"stripe_api_write": "PostRefunds", "stripe_api_read": "PostInvoicesInvoiceVoid", "other_api_write": "PostInvoicesInvoicePay"} {
		if neverSent(tool, op) {
			t.Errorf("%s %s is never sent", tool, op)
		}
	}
	if neverSent("stripe_api_write", "PostInvoicesInvoiceSend") {
		t.Error("the send is never sent")
	}
}

// A broad write narrowed to an operation that moves money or closes a bill for good is refused
// with tool_not_served before anything is asked or read: on a new grant, as
// an addition to a connected one, in any spelling. Narrowing to a refund, a
// read of invoices, or the describing tool goes on as before.
func TestABroadWriteIsNeverNarrowedToAnOperationThatCantBeUndone(t *testing.T) {
	arg := func(path, typ string) map[string]any { return map[string]any{"path": path, "type": typ} }
	write := func(op string) map[string]any {
		return map[string]any{"tool": "stripe_api_write", "name": "chase_ef56", "title": "Chase an invoice",
			"fixed_arguments": map[string]any{"stripe_api_operation_id": op},
			"arguments":       map[string]any{"invoice": arg("parameters.invoice", "string")}}
	}
	refused := func(what string, st int, out map[string]any, op string) {
		t.Helper()
		msg, _ := out["error"].(string)
		if st != 400 || out["code"] != refToolNotServed || out["tool"] != "chase_ef56" || out["field"] != "stripe_api_operation_id" ||
			!strings.Contains(msg, "planque never narrows stripe_api_write to "+strings.TrimSpace(op)) || !strings.Contains(msg, "moves money or closes a bill for good") {
			t.Errorf("%s %q: %d %v", what, op, st, out)
		}
	}

	// A new grant: refused for the operation, not sent on to ask for the describing tool first.
	e := newEnv(t)
	e.enrol()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{write("PostInvoicesInvoiceVoid")}}, nil, "k_setup")
	refused("a new grant", st, out, "PostInvoicesInvoiceVoid")

	// An addition, once the describing tool is approved: each operation, and spellings of one.
	e = newEnv(t)
	id := e.detailsGrant()
	at := "/v1/accounts/" + tenant + "/grants/" + id + "/tools"
	before := e.vendor.CallCount()
	for _, op := range append(slices.Clone(describers[0].Never), " postinvoicesinvoicevoid", "DELETEINVOICESINVOICE") {
		st, out := e.do("POST", at, map[string]any{"tools": []any{write(op)}}, nil, "k_setup")
		refused("an addition", st, out, op)
	}
	if e.vendor.CallCount() != before {
		t.Fatalf("a refused narrowing read the vendor: %d calls, was %d", e.vendor.CallCount(), before)
	}
	_, g := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
	if b := string(mustJSON(g["events"])); strings.Contains(b, "Invoice") {
		t.Fatalf("a refused narrowing is on the record: %s", b)
	}

	// Negative: a refund still narrows, and is asked for.
	e = newEnv(t)
	refund := map[string]any{"tool": "stripe_api_write", "name": "refund_cd34", "title": "Refund a payment", "amount_argument": "amount",
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "PostRefunds"},
		"arguments":       map[string]any{"charge": arg("parameters.charge", "string"), "amount": arg("parameters.amount", "integer")}}
	if _, out := e.addToDetails(refund); out["status"] != StatusPending {
		t.Fatalf("a refund: %v", out)
	}
	// Negative: reads of invoices narrow, open ones pinned, with the same operation ids a write may not take.
	e = newEnv(t)
	id = e.listedGrant()
	reads := []any{
		map[string]any{"tool": "stripe_api_read", "name": "list_invoices", "title": "List overdue invoices", "read_only": true,
			"fixed_arguments": map[string]any{"stripe_api_operation_id": "GetInvoices", "parameters.status": "open"}},
		map[string]any{"tool": "stripe_api_read", "name": "get_invoice", "title": "Read an invoice", "read_only": true,
			"fixed_arguments": map[string]any{"stripe_api_operation_id": "GetInvoicesInvoice"},
			"arguments":       map[string]any{"invoice": arg("parameters.invoice", "string")}},
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": reads}, nil, "k_setup"); st != 200 {
		t.Fatalf("reads: %d %v", st, out)
	}
}

// A call through a broad write that names an operation that moves money or closes a bill for good
// is never sent: on a broad tool, where the taf sets the operation, and on a
// narrowed one a grant from before pins to it. Other operations are sent.
func TestACallToAnOperationThatCantBeUndoneIsNotSent(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	broad := map[string]any{"tool": "stripe_api_write", "name": "stripe_write", "title": "Write to Stripe"}
	id, out := e.addToDetails(broad)
	e.approveAddition(id, out)
	e.kiff.Cards["k_taf"] = map[string]int64{"stripe_write": 0}
	ctx := context.Background()
	call := func(op string) string {
		t.Helper()
		res, err := e.kiff.Call(ctx, "k_taf", "stripe_write", map[string]any{"stripe_api_operation_id": op, "parameters": map[string]any{"invoice": "in_1"}})
		if err != nil || res.Decision != "allowed" {
			t.Fatalf("%s: %+v %v", op, res, err)
		}
		return string(res.Result)
	}
	for _, op := range []string{"PostInvoicesInvoiceVoid", "postinvoicesinvoicepay", " PostCreditNotesIdVoid "} {
		before := e.vendor.CallCount()
		if res := call(op); !strings.Contains(res, "Not sent: planque never sends "+strings.TrimSpace(op)+" through stripe_api_write") || !strings.Contains(res, CodeToolNotServed) {
			t.Errorf("%q: %s", op, res)
		}
		if e.vendor.CallCount() != before {
			t.Errorf("%q reached the vendor", op)
		}
	}
	// Negative: another operation through the same broad tool is sent.
	before := e.vendor.CallCount()
	if res := call("PostInvoicesInvoiceSend"); strings.Contains(res, "Not sent") || e.vendor.CallCount() != before+1 {
		t.Fatalf("the send: %s (%d calls, was %d)", res, e.vendor.CallCount(), before)
	}

	// Fails closed: an operation that isn't one string is never sent either.
	for _, op := range []any{[]any{"PostInvoicesInvoiceVoid"}, map[string]any{"id": "PostInvoicesInvoicePay"}, 42, nil} {
		before := e.vendor.CallCount()
		res, err := e.kiff.Call(ctx, "k_taf", "stripe_write", map[string]any{"stripe_api_operation_id": op})
		if err != nil || !strings.Contains(string(res.Result), "Not sent: planque never sends a call through stripe_api_write whose stripe_api_operation_id isn't one operation's id") ||
			!strings.Contains(string(res.Result), CodeToolNotServed) || e.vendor.CallCount() != before {
			t.Errorf("operation %v: %s (%d calls, was %d)", op, res.Result, e.vendor.CallCount(), before)
		}
	}
	// Unit: arguments that aren't an object, and no operation at all (left to the vendor).
	if why := neverSentCall("stripe_api_write", json.RawMessage(`["x"]`)); !strings.Contains(why, "whose arguments aren't an object") {
		t.Errorf("not an object: %q", why)
	}
	if why := neverSentCall("stripe_api_write", json.RawMessage(`{"parameters":{}}`)); why != "" {
		t.Errorf("no operation: %q", why)
	}
	if why := neverSentCall("stripe_api_read", json.RawMessage(`["x"]`)); why != "" {
		t.Errorf("a read: %q", why)
	}

	// A narrowed write pinned, before this check, to an operation that moves money or closes a bill for good.
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(g.Tools, func(gt GrantedTool) bool { return gt.Name == "stripe_write" })
	g.Tools[i].Fixed = map[string]json.RawMessage{"stripe_api_operation_id": json.RawMessage(`"PostInvoicesInvoiceMarkUncollectible"`)}
	// Its fields checked then: Stripe's own server describes the operation, so the field check lets it by.
	g.Tools[i].FieldsChecked = true
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	before = e.vendor.CallCount()
	res, err := e.kiff.Call(ctx, "k_taf", "stripe_write", map[string]any{})
	if err != nil || !strings.Contains(string(res.Result), "Not sent: planque never sends PostInvoicesInvoiceMarkUncollectible") || e.vendor.CallCount() != before {
		t.Fatalf("pinned from before: %+v %v", res, err)
	}
}
