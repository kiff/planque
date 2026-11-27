package planque

import "slices"

// qontoServed is every tool planque serves on Qonto's own sign-in (the
// `qonto` entry, Qonto's MCP server): its reads, drafting a client invoice,
// and adding or updating a client. Nothing else is served, whatever the
// token's scopes would let Qonto run: client_invoice.write also lets
// Qonto's tools finalize, cancel, send (by email, or as an e-invoice that
// can't be recalled), mark paid and delete, and client.write deletes a
// client. Those stay with the owner in Qonto. It fails closed: a tool Qonto
// adds later, or renames, isn't served until it's reviewed and named here.
//
// The names are Qonto's MCP tool reference (docs.qonto.com/mcp/tools, read
// 27 Nov 2026). Its exports and uploads (export_transactions,
// export_transaction_attachments, request_attachment_upload,
// upload_attachment) write files or create records, so they're left out.
var qontoServed = []string{
	// Reads.
	"get_organization", "list_user_organizations", "list_memberships", "get_authenticated_membership", "list_teams", "list_labels",
	"list_cards",
	"list_transactions", "get_transaction", "list_statements", "get_statement", "get_attachment", "list_transaction_attachments",
	"list_cash_flow_categories",
	"list_clients", "get_client",
	"list_client_invoices", "get_client_invoice", "list_reminders", "list_credit_notes", "get_credit_note",
	"list_quotes", "get_quote", "list_supplier_invoices", "get_supplier_invoice", "list_products",
	"list_payment_links", "get_payment_link", "list_requests",
	"get_subscription", "get_qonto_public_pricing",
	// Writes: a client invoice as a draft (Qonto's default; finalizing and
	// sending are other tools), and a client added or updated.
	"create_client_invoice", "create_client", "update_client",
}

// qontoServedOn says whether planque serves tool on server: on Qonto's MCP
// host only the tools in qontoServed; on any other server, any tool (each
// still under the owner's approval and KIFF's permit).
func qontoServedOn(server, tool string) bool {
	return !onHost(server, qontoHost) || slices.Contains(qontoServed, tool)
}

// servedOnly drops from a server's list the tools planque never serves there,
// so no catalog, no live catalog and no connection ever shows or grants one.
func servedOnly(server string, listed []mcpTool) []mcpTool {
	if !onHost(server, qontoHost) {
		return listed
	}
	return slices.DeleteFunc(slices.Clone(listed), func(t mcpTool) bool { return !slices.Contains(qontoServed, t.Name) })
}

// qontoNotServed refuses a Qonto tool planque doesn't serve on its sign-in.
func qontoNotServed(tool, field string) *toolRefusal {
	return &toolRefusal{Code: refToolNotServed, Tool: tool, Field: field, Want: "one of Qonto's reads, create_client_invoice (a draft), create_client or update_client",
		Msg: tool + " isn't served on Qonto's sign-in: planque serves its reads, drafting a client invoice, and adding or updating a client. Finalizing, cancelling, sending, marking paid and deleting stay with the owner in Qonto; leave it out"}
}

// qontoToolsCheck refuses, before the owner sees anything, a requested tool
// planque doesn't serve on Qonto's sign-in.
func qontoToolsCheck(server string, tools []ToolRequest, field string) *toolRefusal {
	for _, t := range tools {
		if !qontoServedOn(server, t.Tool) {
			return qontoNotServed(t.Tool, field)
		}
	}
	return nil
}
