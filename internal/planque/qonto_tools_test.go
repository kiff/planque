package planque

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	// In Qonto's live list stored on a sign-in (28 Nov): the bank accounts, read-only.
	"list_bank_accounts",
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
	// Reads the privacy policies don't name: cards and their holders, team requests, payment links,
	// quotes (Qonto's consent page says this connection can't read them), statements and the rest.
	"list_cards", "list_requests", "list_payment_links", "get_payment_link", "list_quotes", "get_quote",
	"list_statements", "get_statement", "list_reminders", "list_labels", "list_cash_flow_categories", "get_subscription", "get_qonto_public_pricing",
	"list_user_organizations",
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

	if r := (&Server{}).toolsCheck("tnt_any", qontoServer, []ToolRequest{{Tool: "list_clients"}, {Tool: "send_client_invoice"}}, "tools"); r == nil ||
		r.Code != refToolNotServed || r.Tool != "send_client_invoice" || r.Field != "tools" || !strings.Contains(r.Msg, "Everything else stays with the owner in Qonto") {
		t.Fatalf("refusal: %+v", r)
	}
	if r := (&Server{}).toolsCheck("tnt_any", qontoServer, []ToolRequest{{Tool: "list_clients"}, {Tool: "create_client_invoice"}, {Tool: "update_client"}}, "tool"); r != nil {
		t.Fatalf("served tools refused: %+v", r)
	}
	// The bank accounts and their balances are served by list_bank_accounts,
	// which Qonto's live list marks read-only (28 Nov). A taf that guessed
	// another name reads, in the refusal, that list_bank_accounts gives them,
	// never get_organization, which gives the business's identity only.
	if r := (&Server{}).toolsCheck("tnt_any", qontoServer, []ToolRequest{{Tool: "list_bank_accounts"}}, "tools"); r != nil {
		t.Fatalf("list_bank_accounts refused: %+v", r)
	}
	r := (&Server{}).toolsCheck("tnt_any", qontoServer, []ToolRequest{{Tool: "get_account_balance"}}, "tools")
	if r == nil || !strings.Contains(r.Msg, "(list_bank_accounts)") || !strings.Contains(r.Msg, "identity (get_organization)") ||
		strings.Contains(r.Msg, "balances and IBANs (get_organization)") || strings.Contains(r.Msg, "aren't among them") {
		t.Fatalf("the refusal doesn't name list_bank_accounts for balances: %+v", r)
	}
	for _, name := range qontoServed {
		if !strings.Contains(r.Msg, name) || !strings.Contains(r.Want, name) {
			t.Errorf("the refusal doesn't list %s: %+v", name, r)
		}
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
	if !strings.Contains(got, `"code":"`+CodeToolNotServed+`"`) || !strings.Contains(got, `"sent":false`) || !strings.Contains(got, "Not sent: planque doesn't serve send_client_invoice on Qonto's sign-in. Everything outside its list stays with the owner in Qonto.") {
		t.Fatalf("a never-served Qonto tool: %s", got)
	}
	if n := e.vendor.CallCount(); n != 0 {
		t.Fatalf("the vendor ran %d calls", n)
	}
	var listed []string
	for _, tool := range (&Server{}).relayTools(g) {
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

// Every tools/call to a vendor's MCP server goes through callServed: a new
// path that calls s.MCP.CallTool directly would skip the Qonto filter.
func TestEveryVendorCallGoesThroughCallServed(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || f == "qonto_tools.go" {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "MCP.CallTool(") {
			t.Fatalf("%s calls MCP.CallTool directly: use callServed", f)
		}
	}
	e := newEnv(t)
	if _, err := e.s.callServed(context.Background(), tenant, qontoServer, "token", "send_client_invoice", json.RawMessage(`{}`)); !errors.Is(err, errNotServed) {
		t.Fatalf("callServed sent a tool planque doesn't serve: %v", err)
	}
}

// A Qonto client invoice is created as a draft whatever the taf sets: the
// status and e-invoicing report the tool takes are pinned and returned as
// pinned; other values, at any depth, anything that would send or finalize,
// an argument the schema doesn't name, and a tool with no schema on record
// are refused; other tools and servers pass.
func TestAQontoInvoiceIsCreatedOnlyAsADraft(t *testing.T) {
	takes := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice"},
		InputSchema: json.RawMessage(`{"type":"object","properties":{"client_id":{"type":"string"},"status":{"type":"string"},"report_einvoicing":{"type":"boolean"},
			"items":{"type":"array"},"settings":{"type":"object"}}}`)}
	pinned := func(args string) (map[string]any, map[string]any) {
		t.Helper()
		out, pins, why := qontoDraftOnly(qontoServer, takes, json.RawMessage(args))
		if why != "" {
			t.Fatalf("%s refused: %s", args, why)
		}
		var m map[string]any
		_ = json.Unmarshal(out, &m)
		return m, pins
	}
	m, pins := pinned(`{"client_id":"c1","items":[{"title":"Audit","quantity":1}]}`)
	if m["status"] != "draft" || m["report_einvoicing"] != false || m["client_id"] != "c1" || pins["status"] != "draft" || pins["report_einvoicing"] != false {
		t.Fatalf("left out: %v, pinned %v", m, pins)
	}
	if m, _ := pinned(`{"client_id":"c1","status":"draft","report_einvoicing":false,"settings":{"status":"draft"}}`); m["status"] != "draft" || m["report_einvoicing"] != false {
		t.Fatalf("set to the draft values: %v", m)
	}
	for _, args := range []string{
		`{"status":"unpaid"}`, `{"status":"finalized"}`, `{"report_einvoicing":true}`, `{"report_einvoicing":"true"}`,
		// At any depth, in an object or an array.
		`{"settings":{"status":"unpaid"}}`, `{"settings":{"Status":"unpaid"}}`, `{"settings":{"report_einvoicing":true}}`,
		`{"settings":{"send_email":true}}`, `{"items":[{"title":"x","finalize":true}]}`, `{"settings":{"nested":{"submit_einvoice":true}}}`,
		// Top-level arguments the schema doesn't name, sending or not.
		`{"send_email":true}`, `{"finalize":true}`, `{"einvoicing":true}`, `{"memo":"hello"}`,
		`[]`, `nope`,
	} {
		if _, _, why := qontoDraftOnly(qontoServer, takes, json.RawMessage(args)); why == "" {
			t.Fatalf("%s was let through", args)
		}
	}
	// Without the schema Qonto listed, nothing is sent: Qonto's default could finalize it.
	for _, schema := range []string{``, `not json`, `{"type":"object"}`, `{"type":"object","properties":{}}`} {
		noSchema := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice"}, InputSchema: json.RawMessage(schema)}
		if _, _, why := qontoDraftOnly(qontoServer, noSchema, json.RawMessage(`{}`)); !strings.Contains(why, "can't keep the invoice a draft") {
			t.Fatalf("schema %q: %q", schema, why)
		}
	}
	// A status or report_einvoicing nested in the schema can't be pinned when left out: refused, whatever the arguments.
	for _, schema := range []string{
		`{"type":"object","properties":{"invoice":{"type":"object","properties":{"status":{"type":"string"}}}}}`,
		`{"type":"object","properties":{"invoice":{"type":"object","properties":{"report_einvoicing":{"type":"boolean"}}}}}`,
		`{"type":"object","properties":{"items":{"type":"array","items":{"type":"object","properties":{"Status":{"type":"string"}}}}}}`,
	} {
		nested := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice"}, InputSchema: json.RawMessage(schema)}
		if _, _, why := qontoDraftOnly(qontoServer, nested, json.RawMessage(`{}`)); why == "" || !strings.Contains(why, "draft") {
			t.Fatalf("nested schema %s: %q", schema, why)
		}
	}
	// A schema with no top-level status to pin leaves Qonto's default to decide: refused.
	plain := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice"}, InputSchema: json.RawMessage(`{"type":"object","properties":{"client_id":{"type":"string"}}}`)}
	if _, _, why := qontoDraftOnly(qontoServer, plain, json.RawMessage(`{"client_id":"c1"}`)); !strings.Contains(why, "lists no status") {
		t.Fatalf("a tool without a status: %q", why)
	}
	// A tool that takes a status but no report_einvoicing gets only the status: nothing it would reject is added.
	statusOnly := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice"}, InputSchema: json.RawMessage(`{"type":"object","properties":{"client_id":{"type":"string"},"status":{"type":"string"}}}`)}
	if out, pins, why := qontoDraftOnly(qontoServer, statusOnly, json.RawMessage(`{"client_id":"c1"}`)); why != "" || string(out) != `{"client_id":"c1","status":"draft"}` || len(pins) != 1 {
		t.Fatalf("a tool with a status only: %s %v %s", out, pins, why)
	}
	// A narrowed tool (a Card fixing the client) is checked against Qonto's own schema, kept beside planque's narrowed one.
	narrowed := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice", Fixed: map[string]json.RawMessage{"client_id": json.RawMessage(`"c1"`)}},
		InputSchema: json.RawMessage(`{"type":"object","properties":{"lines":{"type":"array"}}}`), VendorInputSchema: takes.InputSchema}
	if !narrowed.narrowed() {
		t.Fatal("the test tool isn't narrowed")
	}
	if out, pins, why := qontoDraftOnly(qontoServer, narrowed, json.RawMessage(`{"client_id":"c1","items":[]}`)); why != "" || pins["status"] != "draft" || !strings.Contains(string(out), `"status":"draft"`) {
		t.Fatalf("a narrowed tool: %s %v %s", out, pins, why)
	}
	// Narrowed, from before Qonto's schema was kept: refused, never checked against planque's names.
	narrowed.VendorInputSchema = nil
	if _, _, why := qontoDraftOnly(qontoServer, narrowed, json.RawMessage(`{"client_id":"c1"}`)); !strings.Contains(why, "can't keep the invoice a draft") {
		t.Fatalf("a narrowed tool without Qonto's schema: %q", why)
	}
	// Negative: another Qonto tool, and the same name on another server, pass as they are.
	for _, c := range []struct {
		server string
		tool   GrantedTool
	}{{qontoServer, GrantedTool{ToolRequest: ToolRequest{Tool: "create_client"}}}, {"https://mcp.vendor.example/mcp", takes}} {
		if out, _, why := qontoDraftOnly(c.server, c.tool, json.RawMessage(`{"status":"unpaid"}`)); why != "" || string(out) != `{"status":"unpaid"}` {
			t.Fatalf("%s %s: %s %s", c.server, c.tool.Tool, out, why)
		}
	}
}

// A served tool's description is Qonto's own. planque once added to
// get_organization's that it gives each bank account's balance and IBAN; Qonto
// says it gives the organization only, and planque doesn't serve the bank
// accounts. Nothing is added, on Qonto's host or any other.
func TestQontoDescriptionsAreQontosOwn(t *testing.T) {
	listed := []mcpTool{{Name: "get_organization", Description: "Retrieve the organization."}, {Name: "list_transactions", Description: "List transactions."}}
	got := servedOnly(qontoServer, listed)
	if len(got) != 2 || got[0].Description != "Retrieve the organization." || got[1].Description != "List transactions." {
		t.Fatalf("a served description changed: %+v", got)
	}
}
