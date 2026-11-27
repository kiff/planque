package planque

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// Every tool in Qonto's MCP tool reference (docs.qonto.com/mcp/tools, read
// 27 Nov 2026), so a name planque serves is one Qonto has: a typo would
// serve nothing, and a write named by mistake would show here.
var qontoReference = []string{
	"get_organization", "list_user_organizations", "list_memberships", "get_authenticated_membership", "create_membership", "list_teams", "create_team", "list_labels",
	"list_cards", "create_card", "update_card", "change_card_status",
	"list_transactions", "get_transaction", "export_transactions", "list_statements", "get_statement", "get_attachment", "list_transaction_attachments",
	"remove_transaction_attachment", "export_transaction_attachments", "request_attachment_upload", "upload_attachment",
	"list_cash_flow_categories", "create_cash_flow_category", "modify_transaction_cash_flow_category",
	"list_clients", "get_client", "create_client", "update_client", "delete_client",
	"list_client_invoices", "get_client_invoice", "create_client_invoice", "update_client_invoice", "delete_client_invoice", "change_client_invoice_status",
	"mark_client_invoice_as_paid", "send_client_invoice", "list_reminders", "create_reminder", "delete_reminder", "list_credit_notes", "get_credit_note",
	"create_credit_note", "list_quotes", "get_quote", "create_quote", "update_quote", "delete_quote", "send_quote", "list_supplier_invoices",
	"get_supplier_invoice", "change_supplier_invoice_status", "list_products", "create_product",
	"list_payment_links", "create_payment_link", "get_payment_link",
	"list_requests", "create_card_request", "create_multi_transfer_request", "approve_request", "decline_request",
	"get_subscription", "get_qonto_public_pricing",
}

// What Qonto's sign-in never serves, by name: finalizing or cancelling,
// sending (an e-invoice can't be recalled), marking paid, deleting, and
// every other write but the three reviewed ones; and a tool Qonto might add.
var qontoNeverServed = []string{
	"change_client_invoice_status", "send_client_invoice", "mark_client_invoice_as_paid", "delete_client_invoice", "delete_client",
	"update_client_invoice", "change_supplier_invoice_status", "send_quote", "delete_quote", "create_quote", "update_quote", "create_credit_note",
	"create_reminder", "delete_reminder", "create_membership", "create_team", "create_card", "update_card", "change_card_status",
	"create_card_request", "create_multi_transfer_request", "approve_request", "decline_request", "create_payment_link",
	"remove_transaction_attachment", "export_transactions", "export_transaction_attachments", "request_attachment_upload", "upload_attachment",
	"create_cash_flow_category", "modify_transaction_cash_flow_category", "create_product",
	"pay_supplier_invoice", "create_transfer", // not in the reference: a tool Qonto adds isn't served
}

func TestQontoSignInServesOnlyItsReviewedTools(t *testing.T) {
	for _, name := range qontoServed {
		if !slices.Contains(qontoReference, name) {
			t.Fatalf("%s is served but isn't in Qonto's reference", name)
		}
	}
	for _, name := range qontoNeverServed {
		if slices.Contains(qontoServed, name) {
			t.Fatalf("%s is served", name)
		}
	}
	// Its only writes are the three reviewed ones: every other served tool is a read.
	for _, name := range qontoServed {
		if !strings.HasPrefix(name, "get_") && !strings.HasPrefix(name, "list_") && !slices.Contains([]string{"create_client_invoice", "create_client", "update_client"}, name) {
			t.Fatalf("%s is served and is neither a read nor a reviewed write", name)
		}
	}

	listed := make([]mcpTool, 0, len(qontoReference)+2)
	for _, name := range append(slices.Clone(qontoReference), "pay_supplier_invoice", "create_transfer") {
		listed = append(listed, mcpTool{Name: name})
	}
	names := func(tools []mcpTool) []string {
		out := make([]string, 0, len(tools))
		for _, t := range tools {
			out = append(out, t.Name)
		}
		return out
	}
	// On Qonto's host, whatever path or case, the list is the served tools; the server's own list is untouched.
	for _, server := range []string{qontoServer, "https://MCP.qonto.com/mcp", "https://mcp.qonto.com/mcp/", "https://mcp.qonto.com:443/other"} {
		got, want := names(servedOnly(server, listed)), slices.Clone(qontoServed)
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("%s: served %v", server, got)
		}
	}
	if len(listed) != len(qontoReference)+2 {
		t.Fatalf("servedOnly changed the list it was given")
	}
	// Negative: any other server keeps its whole list, sends included.
	for _, server := range []string{"https://mcp.stripe.com", "https://mcp.vendor.example/mcp", "https://qonto.com.example/mcp"} {
		if got := servedOnly(server, listed); len(got) != len(listed) {
			t.Fatalf("%s: dropped tools: %d of %d", server, len(got), len(listed))
		}
		if !qontoServedOn(server, "send_client_invoice") {
			t.Fatalf("%s: a tool refused off Qonto's host", server)
		}
	}

	if r := qontoToolsCheck(qontoServer, []ToolRequest{{Tool: "list_clients"}, {Tool: "send_client_invoice"}}, "tools"); r == nil ||
		r.Code != refToolNotServed || r.Tool != "send_client_invoice" || r.Field != "tools" || !strings.Contains(r.Msg, "stay with the owner in Qonto") {
		t.Fatalf("refusal: %+v", r)
	}
	if r := qontoToolsCheck(qontoServer, []ToolRequest{{Tool: "list_clients"}, {Tool: "create_client_invoice"}, {Tool: "update_client"}}, "tool"); r != nil {
		t.Fatalf("served tools refused: %+v", r)
	}
}

// Asked for on Qonto's MCP server, a tool planque doesn't serve is refused
// before anything is kept, so the owner never sees it to approve.
func TestAQontoSendIsRefusedBeforeTheOwnerSeesIt(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	ask := func(server string, tools ...string) (int, map[string]any) {
		asked := make([]map[string]any, 0, len(tools))
		for _, tool := range tools {
			asked = append(asked, map[string]any{"tool": tool, "name": tool, "title": tool})
		}
		return e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"agent": "agt_taf", "server": server, "tools": asked}, nil, "k_setup")
	}
	for _, never := range []string{"send_client_invoice", "change_client_invoice_status", "mark_client_invoice_as_paid", "delete_client_invoice", "delete_client"} {
		st, out := ask(qontoServer, "create_client_invoice", never)
		if st != 400 || out["code"] != refToolNotServed || out["tool"] != never || out["field"] != "tools" {
			t.Fatalf("%s: %d %v", never, st, out)
		}
	}
	// Negative: the reviewed draft and client tools pass this check.
	if st, out := ask(qontoServer, "create_client_invoice", "create_client"); out["code"] == refToolNotServed {
		t.Fatalf("served tools refused: %d %v", st, out)
	}
}

// A Qonto connection holding a tool from before planque stopped serving it
// never sends a call to it, and doesn't list it to KIFF; its served tools
// stay listed.
func TestACallToAQontoToolPlanqueDoesntServeIsNotSent(t *testing.T) {
	e := newEnv(t)
	g := &Grant{ID: "grt_qonto", Tenant: tenant, Status: StatusConnected, Audience: "aud_qonto",
		Request: Request{Server: qontoServer},
		Tools: []GrantedTool{
			{ToolRequest: ToolRequest{Name: "send_invoice", Tool: "send_client_invoice"}},
			{ToolRequest: ToolRequest{Name: "clients", Tool: "list_clients", ReadOnly: true}},
		}}
	got := mustJSONString(e.s.call(context.Background(), g, "send_invoice", json.RawMessage(`{}`), ""))
	if !strings.Contains(got, `"code":"`+CodeToolNotServed+`"`) || !strings.Contains(got, `"sent":false`) || !strings.Contains(got, "Not sent: planque doesn't serve send_client_invoice") {
		t.Fatalf("a never-served Qonto tool: %s", got)
	}
	if n := e.vendor.CallCount(); n != 0 {
		t.Fatalf("the vendor ran %d calls", n)
	}
	var listed []string
	for _, tool := range relayTools(g) {
		listed = append(listed, tool.Name)
	}
	if !slices.Equal(listed, []string{"clients"}) {
		t.Fatalf("listed to KIFF: %v", listed)
	}
	// Negative: the served tool goes on to the permit check, so it isn't refused as not served.
	if got := mustJSONString(e.s.call(context.Background(), g, "clients", json.RawMessage(`{}`), "")); strings.Contains(got, CodeToolNotServed) {
		t.Fatalf("a served Qonto tool: %s", got)
	}
}
