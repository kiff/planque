package planque

import (
	"net/http"
	"net/url"
	"strings"
)

// A proof test account (PLANQUE_PROOF_TEST_TENANTS: the same ids as KIFF's
// KIFF_PROOF_TEST_TENANTS, in the tenants:<id>[,<id>] form of
// PLANQUE_MANAGED_TEST_TENANTS, set at deploy and never kept in this
// repository) is run by a robot test user. It never moves real money:
// planque takes no live Stripe key for it, and pays nothing for any of its
// tafs (browse_pay). Unset, no account is one, and nothing changes.

// CodeTestTenant is what's refused on a proof test account: a live Stripe
// key, and every browse_pay.
const CodeTestTenant = "test_tenant"

// stripeTestKeyPrefixes are a Stripe key in test mode: a secret key or a
// restricted one.
var stripeTestKeyPrefixes = []string{"sk_test_", "rk_test_"}

// stripeLiveKey reports that key, given for g, is a Stripe key that isn't in
// test mode: any key to Stripe's API (api.stripe.com) that doesn't start
// sk_test_ or rk_test_, and a Stripe secret or restricted key (sk_, rk_) in
// live mode given for any vendor.
func stripeLiveKey(g *Grant, key string) bool {
	for _, p := range stripeTestKeyPrefixes {
		if strings.HasPrefix(key, p) {
			return false
		}
	}
	if strings.HasPrefix(key, "sk_") || strings.HasPrefix(key, "rk_") {
		return true
	}
	return stripeAPI(g)
}

// stripeAPI reports that g is an API-key grant to Stripe's own API.
func stripeAPI(g *Grant) bool {
	if g == nil || g.Request.API == nil {
		return false
	}
	u, err := url.Parse(g.Request.API.Base)
	return err == nil && strings.EqualFold(u.Hostname(), "api.stripe.com")
}

// testTenantKeyRefusal is why key isn't taken for g, on a proof test account,
// or "".
func (s *Server) testTenantKeyRefusal(g *Grant, key string) string {
	if !s.ProofTestTenants[g.Tenant] || !stripeLiveKey(g, key) {
		return ""
	}
	return "this is a test account: it takes only a Stripe key in test mode (sk_test_ or rk_test_), never a live one. Nothing was kept"
}

// On a proof test account, Stripe connects only by an API-key grant to
// Stripe's own API, whose key testTenantKeyRefusal reads. Its other two ways
// carry a key planque never sees: the managed Stripe app (its key is typed on
// the connection platform's page) and a sign-in to Stripe's MCP server
// (mcp.stripe.com, whose mode is chosen per call). Both are refused there,
// before anything is kept or started.
const stripeOnlyByTestKey = "this is a test account: Stripe connects here only with an API key in test mode (sk_test_ or rk_test_), never a managed connection or a sign-in. Nothing was kept"

// testTenantStripeManaged reports that tenant is a proof test account and
// toolkit is Stripe's managed app.
func (s *Server) testTenantStripeManaged(tenant, toolkit string) bool {
	return s.ProofTestTenants[tenant] && provedKey(toolkit) == provedKey(stripeToolkit)
}

// testTenantStripeSignIn reports that tenant is a proof test account and a
// grant of kind to server is a sign-in to Stripe's MCP server.
func (s *Server) testTenantStripeSignIn(tenant, kind, server string) bool {
	return s.ProofTestTenants[tenant] && kind != KindAPIKey && kind != KindBrowser && sameServer(server, stripeServer)
}

// managedOK writes why a managed connection to toolkit isn't made for
// tenant, and reports whether it may be: a proof test account's Stripe
// (test_tenant), or an app the published list doesn't name.
func (s *Server) managedOK(w http.ResponseWriter, tenant, toolkit string) bool {
	if s.testTenantStripeManaged(tenant, toolkit) {
		writeTestTenant(w, stripeOnlyByTestKey)
		return false
	}
	if !s.managedAllowed(tenant, toolkit) {
		writeErr(w, http.StatusForbidden, errNotDisclosed)
		return false
	}
	return true
}

func writeTestTenant(w http.ResponseWriter, why string) {
	writeJSON(w, http.StatusForbidden, map[string]string{"error": why, "code": CodeTestTenant})
}
