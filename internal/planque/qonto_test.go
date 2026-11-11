package planque

import (
	"context"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/oauth"
)

// Every scope Qonto's server lists (its scopes_supported, read 11 Nov 2026).
var qontoListed = []string{"offline_access", "attachment.read", "attachment.write", "bank_account.write", "card.read", "card.write",
	"cash_flow_category.read", "cash_flow_category.write", "client.read", "client.write", "client_invoice.write", "client_invoices.read",
	"einvoicing.read", "insurance_contract.read", "insurance_contract.write", "membership.read", "membership.write", "organization.read",
	"payment_link.read", "payment_link.write", "product.read", "product.write", "request_cards.write", "request_review.write",
	"request_transfers.write", "sepa_direct_debit.read", "sepa_direct_debit.write", "subscription.read", "supplier_invoice.read",
	"supplier_invoice.write", "team.read", "team.write"}

// What a sign-in to a bank never asks for: anything that moves money, or
// touches cards, accounts, direct debits, the team or its members.
var neverAskedOfABank = []string{"transfer", "card", "bank_account", "sepa", "team", "membership"}

// A sign-in to Qonto asks for exactly the reviewed scopes, never all its
// server lists, and nothing that moves money or changes who may.
func TestQontoSignInAsksForTheReviewedScopesOnly(t *testing.T) {
	// It fails closed: whatever path Qonto's metadata names, and whatever it
	// names when the server connected is Qonto's, the list is the reviewed one.
	for _, c := range []struct{ server, resource string }{
		{qontoServer, "https://mcp.qonto.com"}, {qontoServer, "https://mcp.qonto.com/"}, {qontoServer, "https://MCP.qonto.com"},
		{qontoServer, "https://mcp.qonto.com/mcp"}, {"https://mcp.qonto.com/mcp/", "https://mcp.qonto.com/mcp"},
		{qontoServer, "https://auth.qonto.example"}, {"https://tools.example/mcp", "https://mcp.qonto.com/mcp"},
	} {
		resource := c.server + " → " + c.resource
		meta := oauth.Metadata{Resource: c.resource, Scopes: slices.Clone(qontoListed)}
		reviewScopes(c.server, &meta)
		if !slices.Equal(meta.Scopes, qontoScopes) {
			t.Fatalf("%s: asked for %v", resource, meta.Scopes)
		}
		u, err := url.Parse(oauth.AuthURL(meta, "client", "https://planque.example/oauth/callback", "state", "challenge"))
		if err != nil {
			t.Fatal(err)
		}
		asked := strings.Fields(u.Query().Get("scope"))
		if !slices.Equal(asked, qontoScopes) {
			t.Fatalf("%s: the sign-in page asks for %v", resource, asked)
		}
		for _, scope := range asked {
			for _, never := range neverAskedOfABank {
				if strings.Contains(scope, never) {
					t.Fatalf("Qonto's sign-in asks for %q", scope)
				}
			}
		}
	}
	// Every reviewed scope is one Qonto lists: a typo would ask for nothing.
	for _, scope := range qontoScopes {
		if !slices.Contains(qontoListed, scope) {
			t.Fatalf("%q isn't a scope Qonto lists", scope)
		}
	}
	// Another server's sign-in is as it was: all its server lists.
	for _, server := range []string{"https://mcp.qonto.com.example/mcp", "https://qonto.example/mcp", "https://mcp.stripe.com"} {
		other := oauth.Metadata{Resource: server, Scopes: []string{"a", "b"}}
		reviewScopes(server, &other)
		if !slices.Equal(other.Scopes, []string{"a", "b"}) {
			t.Fatalf("%s: scopes changed: %v", server, other.Scopes)
		}
	}
}

// Qonto is planque's own catalog entry, found by its name, and asked for
// like Stripe's: a direct sign-in to Qonto's MCP server.
func TestTheCatalogOffersQontoDirectly(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	base := "/v1/accounts/" + tenant + "/catalog"
	st, out := e.do("GET", base+"?q=qonto", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	entries := out["entries"].([]any)
	if len(entries) != 1 || entries[0].(map[string]any)["entry"] != QontoEntry || entries[0].(map[string]any)["state"] != "available" {
		t.Fatalf("entries: %v", entries)
	}
	// Asked, it discovers Qonto's sign-in from Qonto's own server, as
	// Stripe's entry does Stripe's, and keeps only the reviewed scopes.
	st, out = e.do("POST", base+"/"+QontoEntry+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("Qonto asked: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil || g.Request.Server != qontoServer || g.Request.CatalogEntry != QontoEntry || g.Request.Connector != nil {
		t.Fatalf("the grant: %v %+v", err, g.Request)
	}
	if g.Vendor.Resource != "https://mcp.qonto.com" || !slices.Equal(g.Vendor.Scopes, qontoScopes) {
		t.Fatalf("the sign-in: resource %q, scopes %v", g.Vendor.Resource, g.Vendor.Scopes)
	}
}

// One way to connect a bank: the managed connector's Qonto is neither
// listed, nor found in a tool search, nor connected.
func TestTheAppShowsOneQonto(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	e.s.Connectors = map[string]connector.Provider{"fake": fakeCatalogProvider{}}
	base := "/v1/accounts/" + tenant + "/catalog"
	st, out := e.do("GET", base, nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	qontos := []string{}
	for _, v := range out["entries"].([]any) {
		if v := v.(map[string]any); v["title"] == "Qonto" {
			qontos = append(qontos, v["entry"].(string))
		}
	}
	if len(qontos) != 1 || qontos[0] != QontoEntry {
		t.Fatalf("Qonto entries: %v", qontos)
	}
	st, out = e.do("GET", base+"/tools?q=invoices", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("search: %d %v", st, out)
	}
	for _, v := range out["tools"].([]any) {
		if v.(map[string]any)["app"] == qontoToolkit {
			t.Fatalf("the managed Qonto found: %v", v)
		}
	}
	if st, out := e.do("POST", base+"/connector:fake:"+qontoToolkit+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 409 || !strings.Contains(out["error"].(string), QontoEntry) {
		t.Fatalf("the managed Qonto asked: %d %v", st, out)
	}
}
