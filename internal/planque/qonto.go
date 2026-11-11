package planque

import (
	"slices"

	"github.com/kiff/planque/internal/oauth"
)

// QontoEntry is Qonto in the vendor catalog: planque's own entry for a
// direct sign-in to Qonto's MCP server, a business bank's, where the owner
// signs in on Qonto's page with Qonto's own approval step.
const QontoEntry = "qonto"

// qontoServer is Qonto's own MCP server. Its metadata names its origin,
// https://mcp.qonto.com, as the resource (oauth.coversServer).
const qontoServer = "https://mcp.qonto.com/mcp"

// qontoVendor and qontoIcon are how the catalog shows Qonto's entry.
const qontoVendor = "Qonto"
const qontoIcon = "https://logos.composio.dev/api/qonto"

// qontoScopes is all a Qonto sign-in asks for, reviewed: reading the
// business, its clients, its invoices and their e-invoicing state, and
// writing clients and client invoices, for a draft invoice the owner sends
// themselves. Qonto's server lists more (scopes_supported), among them
// transfers, cards, bank accounts, direct debits, team and members, and
// planque otherwise asks for all a server lists: a token for a bank never
// carries what moves money or changes who may. A Card limits each call too;
// this limits what the token could do at all.
var qontoScopes = []string{
	"offline_access",
	"organization.read",
	"client.read",
	"client.write",
	"client_invoices.read",
	"client_invoice.write",
	"supplier_invoice.read",
	"einvoicing.read",
	"product.read",
	"attachment.read",
}

// reviewScopes keeps a sign-in to a server planque knows to its reviewed
// scopes. Matched on the resource the server's own metadata names, so a
// grant that names Qonto's server another way gets the same list.
func reviewScopes(meta *oauth.Metadata) {
	if sameServer(meta.Resource, "https://mcp.qonto.com") {
		meta.Scopes = slices.Clone(qontoScopes)
	}
}
