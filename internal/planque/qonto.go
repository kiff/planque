package planque

import (
	"context"
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

// widerThanAsked is what a bank's token answer granted beyond what planque
// asked for (the scopes on the grant's sign-in, g.Vendor.Scopes), for a
// sign-in on Qonto's host: granted equals requested, so a token never holds
// more than planque asked, whatever the vendor's consent page showed. An
// answer that names no scope grants what was asked (RFC 6749 §5.1), so it's
// nothing wider. Another vendor: nil, until it's checked the same way.
func widerThanAsked(g *Grant, granted string) []string {
	if !onHost(g.Request.Server, qontoHost) && !onHost(g.Vendor.Resource, qontoHost) {
		return nil
	}
	var extra []string
	for _, name := range scopeNames(granted) {
		if !slices.Contains(g.Vendor.Scopes, name) {
			extra = append(extra, name)
		}
	}
	return extra
}

// fullScopesFor reports whether tenant's sign-ins ask for every scope a
// server lists, Qonto's included (FullScopes).
func (s *Server) fullScopesFor(tenant string) bool {
	return s.FullScopes[tenant] || s.FullScopes["*"]
}

// scopeNames are a token answer's scope names, space-separated as OAuth
// writes them, once each; nil for none.
func scopeNames(scope string) []string {
	var names []string
	for _, name := range strings.Fields(scope) {
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	return names
}

// widerWhy is the reason a wider grant was refused, for the owner and KIFF:
// the vendor, the names it granted beyond the list, and where to end it too.
func widerWhy(g *Grant, extra []string) string {
	return vendorName(g) + " granted more than was asked for (" + strings.Join(extra, ", ") +
		"), so the sign-in was refused and its access wasn't kept. Remove the connection in " + vendorName(g) + "'s settings too."
}

// revokeTokenPair revokes a token answer's refresh and access tokens where
// the vendor offers revocation, best effort: the token is never kept either
// way. Qonto's metadata names no revocation endpoint, so there the owner
// removes it in Qonto (widerWhy says so).
func (s *Server) revokeTokenPair(ctx context.Context, g *Grant, tok oauth.Token) {
	for _, t := range []struct{ tok, hint string }{{tok.Refresh, "refresh_token"}, {tok.Access, "access_token"}} {
		if t.tok != "" {
			_ = s.revokeToken(ctx, g, t.tok, t.hint)
		}
	}
}
