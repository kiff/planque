package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// A browser grant asked before sites were required and shared domains
// refused is held to those checks at each use: with no site, or a shared
// domain among its sites, a call opens nothing (sites_unsafe) and the live
// sign-in view doesn't start. A grant whose sites pass works as before.
func TestAGrantFromBeforeIsHeldToTheSiteChecksAtEachUse(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, browseOpen, map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("a grant whose sites pass: %v", sc)
	}
	for name, sites := range map[string][]string{
		"no site":           nil,
		"a shared domain":   {"shop.example", "github.io"},
		"a registry suffix": {"co.uk"},
	} {
		g.Request.Browser.Sites = sites
		if err := e.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
		sessions := e.b.Sessions
		out, sc := e.call(g, browseOpen, map[string]any{"url": fakes.ShopItem})
		if code(sc) != CodeSitesUnsafe || sc["sent"] != false || e.b.Sessions != sessions || !strings.Contains(out, "ask for") {
			t.Errorf("%s, a call: %s %v", name, out, sc)
		}
		// Each ask is signed anew (a signed ask is never made twice).
		st, view := e.startView(g.ID, map[string]any{"url": "https://shop.example/?case=" + strings.ReplaceAll(name, " ", "-")})
		if st != 409 || view["code"] != CodeSitesUnsafe || e.b.Sessions != sessions {
			t.Errorf("%s, the sign-in view: %d %v", name, st, view)
		}
	}
}
