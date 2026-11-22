package planque

import (
	"context"
	"testing"
)

// A phone re-pairing a Chrome carries its connections with their sites
// ("Your Chrome" 1), read from its own signature, never from planque's word:
// the grants read gives each connection in the owner's Chrome its proof as
// kept, the phone's signature over {chrome_id, grant, sites} or the Chrome a
// re-pairing carried it to. A connection not driven by the owner's Chrome
// shows neither, whatever is stored.
func TestTheGrantsReadGivesEachOwnersChromeConnectionItsSitesProof(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ask := func(agent string, edit func(g *Grant)) string {
		t.Helper()
		st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": agent,
			"tools": []any{map[string]any{"tool": browseOpen}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("ask %s: %d %v", agent, st, out)
		}
		id := out["grant"].(map[string]any)["id"].(string)
		g, _ := c.s.loadGrant(context.Background(), tenant, id)
		edit(g)
		if err := c.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
		return id
	}
	signed := ask("agt_signed", func(g *Grant) { g.ChromeSitesSig, g.Status = testSitesSig, StatusConnected })
	carried := ask("agt_carried", func(g *Grant) { g.ChromeSitesCarried, g.Status = c.r.chrome, StatusConnected })
	hosted := ask("agt_hosted", func(g *Grant) {
		g.Request.Browser.Driver, g.ChromeSitesSig, g.ChromeSitesCarried = "", testSitesSig, c.r.chrome
	})

	st, out := c.do("GET", "/v1/accounts/"+tenant+"/grants", nil, c.phone, "")
	if st != 200 {
		t.Fatalf("grants: %d %v", st, out)
	}
	seen := map[string]map[string]any{}
	for _, row := range out["grants"].([]any) {
		m := row.(map[string]any)
		seen[m["id"].(string)] = m
	}
	if g := seen[signed]; g == nil || g["chrome_sites_sig"] != testSitesSig || g["chrome_sites_carried"] != nil {
		t.Errorf("a signed connection: %v", g)
	}
	if g := seen[carried]; g == nil || g["chrome_sites_carried"] != c.r.chrome || g["chrome_sites_sig"] != nil {
		t.Errorf("a carried connection: %v", g)
	}
	if g := seen[hosted]; g == nil || g["chrome_sites_sig"] != nil || g["chrome_sites_carried"] != nil {
		t.Errorf("a hosted browser's: %v", g)
	}
}
