package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// What a key is, for a proof test account: a key to Stripe's API in test
// mode is taken, any other key to Stripe's API isn't, and a live Stripe
// secret or restricted key isn't, for any vendor.
func TestWhichStripeKeysAreLive(t *testing.T) {
	stripe := &Grant{Request: Request{API: &APIRequest{Base: "https://api.stripe.com/v1"}}}
	other := &Grant{Request: Request{API: &APIRequest{Base: "https://api.example.com"}}}
	for _, c := range []struct {
		g    *Grant
		key  string
		live bool
	}{
		{stripe, "sk_test_abc", false}, {stripe, "rk_test_abc", false},
		{stripe, "sk_live_abc", true}, {stripe, "rk_live_abc", true}, {stripe, "pk_test_abc", true}, {stripe, "abc", true},
		{other, "sk_live_abc", true}, {other, "rk_live_abc", true}, {other, "sk_test_abc", false},
		{other, "a-vendor-key", false}, {&Grant{}, "a-vendor-key", false},
	} {
		if got := stripeLiveKey(c.g, c.key); got != c.live {
			t.Errorf("%s for %+v: live %v, want %v", c.key, c.g.Request.API, got, c.live)
		}
	}
}

// On a proof test account, a live Stripe key is refused with nothing kept,
// as given and as a replacement; a test key is taken. Another account takes
// the same live key as before.
func TestAProofTestAccountTakesNoLiveStripeKey(t *testing.T) {
	for _, testAccount := range []bool{true, false} {
		e := newEnv(t)
		v := newRESTVendor(t) // its key is sk_test_secret_123
		if testAccount {
			e.s.ProofTestTenants = map[string]bool{tenant: true}
		}
		e.enrol()
		st, out := e.askAPIGrant(v, []map[string]any{}, nil)
		if st != 201 {
			t.Fatalf("ask: %d %v", st, out)
		}
		id, digest := out["id"].(string), out["digest"].(string)
		for _, live := range []string{"sk_live_abc123", "rk_live_abc123"} {
			st, out := e.giveKey(id, digest, live)
			if !testAccount {
				if st != 200 || out["status"] != StatusListed {
					t.Fatalf("another account, %s: %d %v", live, st, out)
				}
				break
			}
			if st != 403 || out["code"] != CodeTestTenant || !strings.Contains(out["error"].(string), "test mode") {
				t.Fatalf("a test account, %s: %d %v", live, st, out)
			}
			g, err := e.s.loadGrant(context.Background(), tenant, id)
			if err != nil || g.Status != StatusPending || g.SignedIn {
				t.Fatalf("a refused key changed the grant: %+v %v", g, err)
			}
			if _, _, err := e.s.openToken(context.Background(), g); err == nil {
				t.Fatalf("a refused key was kept")
			}
		}
		if !testAccount {
			continue
		}
		if st, out := e.giveKey(id, digest, v.key); st != 200 || out["status"] != StatusListed {
			t.Fatalf("a test key on a test account: %d %v", st, out)
		}
	}
}

// On a proof test account, browse_pay is refused for every agent before
// anything starts, even with the owner's approval on its permit; another
// account's taf pays as before.
func TestAProofTestAccountsTafsPayNothing(t *testing.T) {
	for _, testAccount := range []bool{true, false} {
		e := newBrowserEnv(t)
		g := e.connect()
		if testAccount {
			e.s.ProofTestTenants = map[string]bool{tenant: true}
		}
		if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); code(sc) != "" {
			t.Fatalf("checkout: %v", sc)
		}
		sessions := e.b.Sessions
		out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
		if testAccount {
			if code(sc) != CodeTestTenant || sc["sent"] != false || sc["refused_by"] != "connection" || len(e.b.Orders) != 0 ||
				e.b.Sessions != sessions || !strings.Contains(out, "test account") {
				t.Fatalf("a test account's pay: %s %v", out, sc)
			}
			// Another agent on the same account: refused too.
			g.Request.Agent = "agt_other"
			if _, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != CodeTestTenant || len(e.b.Orders) != 0 {
				t.Fatalf("another agent's pay: %v", sc)
			}
			continue
		}
		if code(sc) != "" || len(e.b.Orders) != 1 {
			t.Fatalf("another account's pay: %s %v", out, sc)
		}
	}
}
