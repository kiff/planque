package planque

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
)

func ucpWith(handlers string) []byte {
	return []byte(`{"ucp":{"version":"2026-01-23","payment_handlers":{"com.stripe.payments":[` + handlers + `]}}}`)
}

const ucpTest = `{"id":"stripe_payments","config":{"environment":"testmode","business_profile":"profile_test_61TU90nIeGjU7NNVXA6TU90m7ISQWsBxpcx9lASWWXTk"}}`

// The seller's profile is read from the checkout page's own host, test mode
// only, never two that differ, and answered with the host it was read on.
func TestTheSellersStripeProfileComesFromItsOwnSite(t *testing.T) {
	cases := map[string]struct {
		doc  []byte
		want string
	}{
		"test mode":       {ucpWith(ucpTest), "profile_test_61TU90nIeGjU7NNVXA6TU90m7ISQWsBxpcx9lASWWXTk"},
		"live mode":       {ucpWith(`{"config":{"environment":"livemode","business_profile":"profile_live_1"}}`), ""},
		"two that differ": {ucpWith(ucpTest + `,{"config":{"environment":"testmode","business_profile":"profile_test_other"}}`), ""},
		"the same twice":  {ucpWith(ucpTest + "," + ucpTest), "profile_test_61TU90nIeGjU7NNVXA6TU90m7ISQWsBxpcx9lASWWXTk"},
		"not a profile":   {ucpWith(`{"config":{"environment":"testmode","business_profile":"acct_1"}}`), ""},
		"not UCP":         {[]byte(`<html>`), ""},
	}
	for name, c := range cases {
		if got := stripeProfileOf(c.doc); got != c.want {
			t.Errorf("%s: %q", name, got)
		}
	}

	e := newBrowserEnv(t)
	g := e.connect()
	var asked []string
	doc := ucpWith(ucpTest)
	e.s.UCPGet = func(_ context.Context, host string) ([]byte, error) {
		asked = append(asked, host)
		if doc == nil {
			return nil, errors.New("none")
		}
		return doc, nil
	}
	// Off by default: nothing is fetched, nothing answered.
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["stripe_profile"] != nil || len(asked) != 0 {
		t.Fatalf("off: %v, asked %v", sc, asked)
	}
	e.s.StripeProfiles = true
	out, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout})
	if sc["stripe_profile"] != "profile_test_61TU90nIeGjU7NNVXA6TU90m7ISQWsBxpcx9lASWWXTk" || sc["stripe_profile_from"] != "shop.example" || len(asked) != 1 || asked[0] != "shop.example" {
		t.Fatalf("on: %s, asked %v", out, asked)
	}
	// The site names another profile before the pay: nothing is sent.
	doc = ucpWith(`{"config":{"environment":"testmode","business_profile":"profile_test_other"}}`)
	out, sc = e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
	if code(sc) != browser.CodeCheckoutChanged || len(e.b.Orders) != 0 || !strings.Contains(out, "Stripe profile") {
		t.Fatalf("a changed profile: %s", out)
	}
}
