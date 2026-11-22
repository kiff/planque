package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/store"
)

// Coach asks for a browser connection as tafs's a-browser-taf says: the
// catalog's "browser@1", the sites and the tools, with no driver and no Chrome
// (it sees the paired Chromes by name only). On 22 Nov 2026 planque answered
// "no such catalog entry" to that id, and would then have answered "no
// browser" before reading the request. With no hosted browser and the owner's
// Chrome offered, it's a grant in the owner's Chrome: planque names the one
// paired, before the words and the digest are made, so the name the phone
// shows, the id it signs and what the owner approves are one Chrome.
func TestCoachsBrowserAskIsAGrantInTheOwnersPairedChrome(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil // production has no hosted browser
	ask := func(entry string, b map[string]any) (int, map[string]any) {
		return c.do("POST", "/v1/accounts/"+tenant+"/catalog/"+entry+"/grants", map[string]any{"agent": "agt_taf", "agent_name": "Wiki titles",
			"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": b}, nil, "k_setup")
	}
	for _, entry := range []string{"browser@1", "browser"} {
		st, out := ask(entry, map[string]any{"sites": []string{"wikipedia.org"}})
		if st != 201 {
			t.Fatalf("%s: %d %v", entry, st, out)
		}
		g, err := c.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
		if err != nil {
			t.Fatal(err)
		}
		b := g.Request.Browser
		if b == nil || b.Driver != DriverExtension || b.Chrome != c.r.chrome || b.ChromeName != "Work Mac" {
			t.Fatalf("%s: the grant's browser %+v, want the owner's paired Chrome", entry, b)
		}
		// The digest the phone checks is over the request with that Chrome in it,
		// and the words name it.
		if d, err := digest(g.Request); err != nil || d != g.Digest {
			t.Errorf("%s: the digest isn't the request's own: %v", entry, err)
		}
		if !strings.Contains(g.Text, `your own Chrome "Work Mac"`) {
			t.Errorf("%s: the words don't name the Chrome: %s", entry, g.Text)
		}
		// The phone approves it with its signature over the sites, as any
		// grant in the owner's Chrome.
		hash, _ := profile.Digest(g.Request.Profile)
		approved := ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}
		if st, out := c.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": approved,
			"chrome_sites_sig": testSitesSig}, c.phone, ""); st != 200 {
			t.Errorf("%s: approve: %d %v", entry, st, out)
		}
	}

	// Still refused: no sites (the owner's Chrome reaches only the sites named),
	// another catalog id, and the owner's Chrome no longer offered.
	if st, out := ask("browser@1", map[string]any{}); st != 400 {
		t.Errorf("no sites: %d %v", st, out)
	}
	if st, out := ask("browser@2", map[string]any{"sites": []string{"wikipedia.org"}}); st != 404 {
		t.Errorf("browser@2: %d %v", st, out)
	}
	c.s.ChromeOffer = ChromeOffer{}
	if st, out := ask("browser@1", map[string]any{"sites": []string{"wikipedia.org"}}); st != 503 {
		t.Errorf("not offered, no hosted browser: %d %v", st, out)
	}
}

// With several paired Chromes planque never guesses: the ask names one by the
// name the owner gave it, or it's refused with the names to choose from.
func TestWithSeveralChromesTheAskNamesOne(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	other := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	if err := putDoc(context.Background(), c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(other.ID)}, other, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	ask := func(b map[string]any) (int, map[string]any) {
		b["sites"] = []string{"wikipedia.org"}
		return c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_taf",
			"tools": []any{map[string]any{"tool": browseOpen}}, "browser": b}, nil, "k_setup")
	}
	st, out := ask(map[string]any{})
	if st != 400 || !strings.Contains(mustJSONString(out), "Work Mac") || !strings.Contains(mustJSONString(out), "Home laptop") {
		t.Fatalf("several, none named: %d %v", st, out)
	}
	if st, out := ask(map[string]any{"chrome_name": "Office"}); st != 400 {
		t.Errorf("a name no Chrome has: %d %v", st, out)
	}
	st, out = ask(map[string]any{"chrome_name": "home laptop"})
	if st != 201 {
		t.Fatalf("named: %d %v", st, out)
	}
	g, _ := c.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
	if g.Request.Browser.Chrome != other.ID || g.Request.Browser.ChromeName != "Home laptop" {
		t.Errorf("named: %+v", g.Request.Browser)
	}
}
