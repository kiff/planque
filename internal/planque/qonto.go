package planque

import (
	"net/url"
	"slices"
	"strings"

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

// qontoHost is the host of Qonto's MCP server and of the resource its
// metadata names.
const qontoHost = "mcp.qonto.com"

// qontoToolkit is the managed connector's Qonto, which the direct entry
// replaces: one way to connect a bank.
const qontoToolkit = "qonto"

// reviewScopes keeps a sign-in to a server planque knows to its reviewed
// scopes. It fails closed: the server being connected, or any resource on
// Qonto's host, whatever path the metadata names, gets the reviewed list,
// never all the server lists.
func reviewScopes(server string, meta *oauth.Metadata) {
	if onHost(server, qontoHost) || onHost(meta.Resource, qontoHost) {
		meta.Scopes = slices.Clone(qontoScopes)
	}
}

// onHost reports whether raw is a URL on host, whatever its scheme, port or path.
func onHost(raw, host string) bool {
	u, err := url.Parse(strings.TrimSpace(raw))
	return err == nil && strings.EqualFold(u.Hostname(), host)
}

// directEntryFor is planque's own entry that replaces the managed
// connector's toolkit, or "" for a toolkit offered as it is.
func directEntryFor(toolkit string) string {
	switch toolkit {
	case stripeToolkit:
		return StripeEntry
	case qontoToolkit:
		return QontoEntry
	}
	return ""
}
