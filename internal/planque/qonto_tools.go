package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
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
	// The bank accounts, with their balances and IBANs where the member's own
	// Qonto access shows them. get_organization gives the business's identity
	// only, and names this read for balances; Qonto's live list (28 Nov) marks
	// it read-only, with only page and per_page as input.
	"list_bank_accounts",
	"list_client_invoices", "get_client_invoice", "list_credit_notes", "get_credit_note", "list_supplier_invoices", "get_supplier_invoice", "list_products",
	// Writes: a client invoice, always as a draft (qontoDraftOnly), and a
	// client added or updated.
	"create_client_invoice", "create_client", "update_client",
}

// servesAll says planque passes the vendor's whole live list through for
// tenant's grants on server: a test account's (PLANQUE_MANAGED_TEST_TENANTS),
// once PLANQUE_LIVE_LIST=test is set, which waits for KIFF's Card to decide
// each call's class there (step 9 of the policy move, brought forward for
// test accounts). Every other account, and Qonto's host until then, keeps
// qontoServed. Off any server planque filters, every tool is served anyway.
func (s *Server) servesAll(tenant, server string) bool {
	return s.LiveListTest && s.ManagedTestTenants[tenant] && onHost(server, qontoHost)
}

// served says planque serves tool on server for tenant's grants.
func (s *Server) served(tenant, server, tool string) bool {
	return s.servesAll(tenant, server) || qontoServedOn(server, tool)
}

// servedFor is servedOnly for tenant: the whole list where servesAll.
func (s *Server) servedFor(tenant, server string, listed []mcpTool) []mcpTool {
	if s.servesAll(tenant, server) {
		return listed
	}
	return servedOnly(server, listed)
}

// servedRule names the rule tenant's tools on server are filtered by: the
// whole live list ("live") where servesAll; on Qonto's MCP host, the tools in
// qontoServed, as a short digest; on any other server, none, since every tool
// there is served. A grant's stored catalog keeps it (Grant.CatalogRule), so a
// change to what's served reaches grants listed before it.
func (s *Server) servedRule(tenant, server string) string {
	if s.servesAll(tenant, server) {
		return "live"
	}
	if !onHost(server, qontoHost) {
		return ""
	}
	sum := sha256.Sum256([]byte(strings.Join(qontoServed, ",")))
	return "qonto:" + hex.EncodeToString(sum[:6])
}

// qontoServedOn says whether planque serves tool on server: on Qonto's MCP
// host only the tools in qontoServed; on any other server, any tool (each
// still under the owner's approval and KIFF's permit).
func qontoServedOn(server, tool string) bool {
	return !onHost(server, qontoHost) || slices.Contains(qontoServed, tool)
}

// servedOnly drops from a server's list the tools planque never serves there,
// so no catalog, no live catalog and no connection ever shows or grants one.
// A served tool's description is Qonto's own: planque adds nothing to it.
func servedOnly(server string, listed []mcpTool) []mcpTool {
	if !onHost(server, qontoHost) {
		return listed
	}
	return slices.DeleteFunc(slices.Clone(listed), func(t mcpTool) bool { return !slices.Contains(qontoServed, t.Name) })
}

// qontoNotServed refuses a Qonto tool planque doesn't serve on its sign-in.
func qontoNotServed(tool, field string) *toolRefusal {
	return &toolRefusal{Code: refToolNotServed, Tool: tool, Field: field, Want: "one of " + strings.Join(qontoServed, ", "),
		Msg: tool + " isn't served on Qonto's sign-in: planque serves only the business's identity (get_organization), its bank accounts with their balances and IBANs where the member's Qonto access shows them (list_bank_accounts), the team, transactions and attachments, clients, invoices and credit notes, supplier invoices and products, drafting a client invoice, and adding or updating a client, by these tools: " + strings.Join(qontoServed, ", ") + ". Everything else stays with the owner in Qonto; leave it out"}
}

// qontoToolsCheck refuses, before the owner sees anything, a requested tool
// planque doesn't serve on Qonto's sign-in.
func (s *Server) toolsCheck(tenant, server string, tools []ToolRequest, field string) *toolRefusal {
	for _, t := range tools {
		if !s.served(tenant, server, t.Tool) {
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
func (s *Server) callServed(ctx context.Context, tenant, server, access, tool string, args json.RawMessage) (json.RawMessage, error) {
	if !s.served(tenant, server, tool) {
		return nil, errNotServed
	}
	return s.MCP.CallTool(ctx, server, access, tool, args)
}

// qontoDraftOnly keeps a Qonto client invoice a draft, whatever the taf set:
// Qonto's create takes a status (draft, or unpaid, which finalizes it) and,
// for e-invoicing, whether to report it, which sends it on creation with no
// recall. It fails closed:
//   - without the tool's input schema as Qonto listed it (kept beside a
//     narrowed tool's own), or without a top-level status in it to pin,
//     nothing is sent;
//   - a top-level argument the schema doesn't name is refused, so a new
//     argument Qonto adds isn't passed on before it's reviewed;
//   - at any depth, a status other than draft, a report_einvoicing other
//     than false, or a key that would send, finalize or submit is refused;
//   - a status or report_einvoicing the schema declares below the top level
//     is refused, since it can't be pinned when the taf leaves it out;
//   - a status or report_einvoicing the schema takes is pinned to draft and
//     false, and returned as pinned, for the operation's record;
//   - the fields Qonto takes from Italian organisations only (qontoItalianOnly)
//     follow the organisation's country, read at sign-in (Grant.Country): for
//     another country, a call that sets payment_reporting or
//     stamp_duty_amount is refused, and report_einvoicing is left out, not
//     pinned, as Qonto asks of a non-Italian organisation. For Italy, or a
//     country planque couldn't read, report_einvoicing is pinned to false:
//     left out, an Italian organisation's invoice goes to SdI on creation.
//
// Other tools and servers pass as they are.
func qontoDraftOnly(server string, tool GrantedTool, args json.RawMessage, country string) (json.RawMessage, map[string]any, string) {
	if !onHost(server, qontoHost) || tool.Tool != "create_client_invoice" {
		return args, nil, ""
	}
	// Qonto's own schema: a narrowed tool's InputSchema is planque's, with
	// planque's names, so its vendor schema is kept beside it.
	vendor := tool.InputSchema
	if tool.narrowed() {
		vendor = tool.VendorInputSchema
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if json.Unmarshal(vendor, &schema) != nil || len(schema.Properties) == 0 {
		return nil, nil, "planque has no record of the arguments Qonto's create_client_invoice takes, so it can't keep the invoice a draft; connect Qonto again"
	}
	// The status is what keeps it a draft: without one to pin, Qonto's
	// default decides, so nothing is sent until that's reviewed.
	if _, ok := schema.Properties["status"]; !ok {
		return nil, nil, "Qonto's create_client_invoice lists no status planque can pin, so it can't keep the invoice a draft; it isn't sent until that's reviewed"
	}
	// A status or report_einvoicing the schema declares below the top level
	// can't be pinned when the taf leaves it out, and Qonto's default could
	// finalize or report the invoice: nothing is sent until it's reviewed.
	for k, v := range schema.Properties {
		if nested := nestedDraftField(v); nested != "" {
			return nil, nil, "Qonto's create_client_invoice takes " + nested + " inside " + k + ", which planque can't keep a draft; it isn't sent until that's reviewed"
		}
	}
	var in map[string]any
	dec := json.NewDecoder(bytes.NewReader(args))
	dec.UseNumber()
	if dec.Decode(&in) != nil || in == nil {
		return nil, nil, "the arguments are not a JSON object"
	}
	for k := range in {
		if _, ok := schema.Properties[k]; !ok {
			return nil, nil, k + " isn't an argument Qonto's create_client_invoice lists, so it isn't passed on"
		}
	}
	if why := draftOnlyTree(in); why != "" {
		return nil, nil, why
	}
	// Another country's organisation sends none of the Italian-only fields.
	elsewhere := country != "" && country != "IT"
	if elsewhere {
		for _, k := range qontoItalianOnly {
			if v, ok := in[k]; ok && v != nil && k != "report_einvoicing" {
				return nil, nil, k + " is for Italian organisations only, and this Qonto organisation's country is " + country + ", so it isn't sent; leave it out"
			}
		}
		delete(in, "report_einvoicing")
	}
	pinned := map[string]any{}
	if _, ok := schema.Properties["status"]; ok {
		in["status"], pinned["status"] = "draft", "draft"
	}
	if _, ok := schema.Properties["report_einvoicing"]; ok && !elsewhere {
		in["report_einvoicing"], pinned["report_einvoicing"] = false, false
	}
	out, err := json.Marshal(in)
	if err != nil {
		return nil, nil, "the arguments couldn't be kept as a draft"
	}
	if len(pinned) == 0 {
		pinned = nil
	}
	return out, pinned, ""
}

// qontoItalianOnly are the arguments of Qonto's create_client_invoice its
// own descriptions say only an Italian organisation includes ("Non-Italian
// organizations should not include this property").
var qontoItalianOnly = []string{"report_einvoicing", "payment_reporting", "stamp_duty_amount"}

// qontoOrganizationTool is Qonto's read of the signed-in organisation.
const qontoOrganizationTool = "get_organization"

// signedInCountry is the organisation's country a Qonto sign-in is for,
// read with Qonto's own get_organization, once, by planque itself, right
// after the owner signed in, with the new token: never a caller's word.
// "" when it can't be read or names no country, and the draft pin then
// stays as for an Italian organisation. The read is on the grant's record.
func (s *Server) signedInCountry(ctx context.Context, g *Grant, access string) string {
	if s.MCP == nil {
		return ""
	}
	res, err := s.callServed(ctx, g.Tenant, g.Request.Server, access, qontoOrganizationTool, json.RawMessage(`{}`))
	if err != nil {
		s.log().WarnContext(ctx, "the sign-in's organisation couldn't be read", slog.String("grant", g.ID), slog.String("err", err.Error()))
		g.event(s.now(), "couldn't read the organisation's country")
		return ""
	}
	country := organizationCountry(res)
	if country == "" {
		g.event(s.now(), "read the organisation; it named no country")
		return ""
	}
	g.event(s.now(), "read the organisation's country: "+country)
	return country
}

// organizationCountry reads get_organization's answer, its
// structuredContent or else the same JSON as its text: the organisation's
// legal_country, at the top or under organization, as two letters. "" for
// an error, another shape, or anything that isn't a country code.
func organizationCountry(res json.RawMessage) string {
	type org struct {
		Country      string `json:"legal_country"`
		Organization *struct {
			Country string `json:"legal_country"`
		} `json:"organization"`
	}
	var out struct {
		Structured *org `json:"structuredContent"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if json.Unmarshal(res, &out) != nil || out.IsError {
		return ""
	}
	found := []*org{out.Structured}
	for _, c := range out.Content {
		var t org
		if c.Type == "text" && json.Unmarshal([]byte(c.Text), &t) == nil {
			found = append(found, &t)
		}
	}
	for _, o := range found {
		if o == nil {
			continue
		}
		c := o.Country
		if c == "" && o.Organization != nil {
			c = o.Organization.Country
		}
		if c = strings.ToUpper(strings.TrimSpace(c)); len(c) == 2 && c[0] >= 'A' && c[0] <= 'Z' && c[1] >= 'A' && c[1] <= 'Z' {
			return c
		}
	}
	return ""
}

// draftOnlyTree refuses, anywhere in v, what would make a Qonto client
// invoice more than a draft.
func draftOnlyTree(v any) string {
	switch t := v.(type) {
	case map[string]any:
		for k, child := range t {
			switch key := strings.ToLower(k); {
			case key == "status":
				if child != "draft" {
					return "a Qonto client invoice is created only as a draft: finalizing it is the owner's, in Qonto"
				}
			case key == "report_einvoicing":
				if child != false {
					return "a Qonto client invoice isn't reported for e-invoicing on creation: sending it is the owner's, in Qonto"
				}
			case strings.Contains(key, "send"), strings.Contains(key, "finali"), strings.Contains(key, "submit"), strings.Contains(key, "einvoic"), strings.Contains(key, "e_invoic"):
				return k + " would send or finalize the invoice: a Qonto client invoice is created only as a draft, which the owner sends"
			}
			if why := draftOnlyTree(child); why != "" {
				return why
			}
		}
	case []any:
		for _, child := range t {
			if why := draftOnlyTree(child); why != "" {
				return why
			}
		}
	}
	return ""
}

// nestedDraftField names a status or report_einvoicing property anywhere in
// a schema node, or "" for none.
func nestedDraftField(node json.RawMessage) string {
	var v any
	if json.Unmarshal(node, &v) != nil {
		return ""
	}
	var walk func(any) string
	walk = func(v any) string {
		switch t := v.(type) {
		case map[string]any:
			if props, ok := t["properties"].(map[string]any); ok {
				for k := range props {
					if key := strings.ToLower(k); key == "status" || key == "report_einvoicing" {
						return k
					}
				}
			}
			for _, child := range t {
				if found := walk(child); found != "" {
					return found
				}
			}
		case []any:
			for _, child := range t {
				if found := walk(child); found != "" {
					return found
				}
			}
		}
		return ""
	}
	return walk(v)
}
