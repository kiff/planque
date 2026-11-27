package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
)

// qontoServed is every tool planque serves on Qonto's own sign-in (the
// `qonto` entry, Qonto's MCP server): the reads of what the privacy
// policies name for a connected bank (the business, its accounts, IBANs,
// transactions and team; clients; client and supplier invoices, with their
// credit notes; products; attachments), drafting a client invoice, and adding or updating a client.
// Nothing else is served, whatever the token's scopes would let Qonto run:
// client_invoice.write also lets Qonto's tools finalize, cancel, send (by
// email, or as an e-invoice that can't be recalled), mark paid and delete,
// and client.write deletes a client. Those stay with the owner in Qonto.
// Reads the policies don't name (cards and their holders, team requests,
// payment links, quotes, statements, reminders, labels,
// cash-flow categories, the subscription) aren't served either. It fails
// closed: a tool Qonto adds or renames isn't served until it's reviewed and
// named here.
//
// The names are Qonto's MCP tool reference (docs.qonto.com/mcp/tools, read
// 27 Nov 2026).
var qontoServed = []string{
	// Reads.
	"get_organization", "get_authenticated_membership", "list_memberships", "list_teams",
	"list_transactions", "get_transaction", "get_attachment", "list_transaction_attachments",
	"list_clients", "get_client",
	"list_client_invoices", "get_client_invoice", "list_credit_notes", "get_credit_note", "list_supplier_invoices", "get_supplier_invoice", "list_products",
	// Writes: a client invoice, always as a draft (qontoDraftOnly), and a
	// client added or updated.
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

// errNotServed is a vendor call planque refused to send: the tool isn't
// served on that server.
var errNotServed = errors.New("planque doesn't serve this tool on Qonto's sign-in, so it wasn't sent")

// callServed is every tools/call planque sends to a vendor's MCP server, the
// relay's, a describing tool's and a sign-in's account read alike: a tool it
// doesn't serve on that server is never sent, whoever asks.
func (s *Server) callServed(ctx context.Context, server, access, tool string, args json.RawMessage) (json.RawMessage, error) {
	if !qontoServedOn(server, tool) {
		return nil, errNotServed
	}
	return s.MCP.CallTool(ctx, server, access, tool, args)
}

// qontoDraftOnly keeps a Qonto client invoice a draft, whatever the taf set:
// Qonto's create takes a status (draft, or unpaid, which finalizes it) and,
// for e-invoicing, whether to report it, which sends it on creation with no
// recall. A field the tool takes is pinned to draft and to not reporting; a
// value other than those, or any argument that would send, finalize or
// submit, is refused. Other tools and servers pass as they are.
func qontoDraftOnly(server string, tool GrantedTool, args json.RawMessage) (json.RawMessage, string) {
	if !onHost(server, qontoHost) || tool.Tool != "create_client_invoice" {
		return args, ""
	}
	var in map[string]any
	dec := json.NewDecoder(bytes.NewReader(args))
	dec.UseNumber()
	if dec.Decode(&in) != nil || in == nil {
		return nil, "the arguments are not a JSON object"
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	_ = json.Unmarshal(tool.InputSchema, &schema)
	for k, v := range in {
		switch key := strings.ToLower(k); {
		case key == "status":
			if v != "draft" {
				return nil, "a Qonto client invoice is created only as a draft: finalizing it is the owner's, in Qonto"
			}
		case key == "report_einvoicing":
			if v != false {
				return nil, "a Qonto client invoice isn't reported for e-invoicing on creation: sending it is the owner's, in Qonto"
			}
		case strings.Contains(key, "send"), strings.Contains(key, "finali"), strings.Contains(key, "submit"), strings.Contains(key, "einvoic"), strings.Contains(key, "e_invoic"):
			return nil, k + " would send or finalize the invoice: a Qonto client invoice is created only as a draft, which the owner sends"
		}
	}
	if _, ok := schema.Properties["status"]; ok {
		in["status"] = "draft"
	}
	if _, ok := schema.Properties["report_einvoicing"]; ok {
		in["report_einvoicing"] = false
	}
	out, err := json.Marshal(in)
	if err != nil {
		return nil, "the arguments couldn't be kept as a draft"
	}
	return out, ""
}
