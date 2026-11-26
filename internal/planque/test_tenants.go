package planque

import (
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
