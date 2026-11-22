package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// A connection whose Chrome is gone, its record deleted before unpairing
// ended its connections ("Your Chrome" 3), is never stuck: the owner removes it
// as any other (the phone's device key, nothing asked of the Chrome), or
// switches it to a Chrome still paired ("Your Chrome" 2), and nothing is left
// listed on the Chrome that's gone.
func TestAConnectionWhoseChromeIsGoneIsRemovedOrSwitched(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ctx := context.Background()
	ask := func(agent string) string {
		t.Helper()
		st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": agent,
			"tools": []any{map[string]any{"tool": browseOpen}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("ask %s: %d %v", agent, st, out)
		}
		id := out["grant"].(map[string]any)["id"].(string)
		g, _ := c.s.loadGrant(ctx, tenant, id)
		g.ChromeSitesSig, g.Status = testSitesSig, StatusConnected
		if err := c.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
		return id
	}
	removed, switched := ask("agt_remove"), ask("agt_switch")
	home := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	if err := putDoc(ctx, c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(home.ID)}, home, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	// The Chrome both are on goes, as it could before unpairing ended them.
	if _, ok, err := c.s.KV.Delete(ctx, tenantPK(tenant), skChrome(c.r.chrome)); err != nil || !ok {
		t.Fatalf("the Chrome's record: %v %v", ok, err)
	}

	if st, out := c.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+removed, nil, c.phone, ""); st != 200 {
		t.Fatalf("remove: %d %v", st, out)
	}
	if g, _ := c.s.loadGrant(ctx, tenant, removed); g.Status != StatusDisconnected {
		t.Fatalf("removed: %s", g.Status)
	}

	newSig := strings.Replace(testSitesSig, "s1gned", "s1GNED", 1)
	if st, out := c.do("PUT", "/v1/accounts/"+tenant+"/grants/"+switched+"/chrome", map[string]any{"chrome_id": home.ID, "chrome_sites_sig": newSig}, c.phone, ""); st != 200 {
		t.Fatalf("switch: %d %v", st, out)
	}
	if g, _ := c.s.loadGrant(ctx, tenant, switched); g.Request.Browser.Chrome != home.ID || g.ChromeSitesSig != newSig || g.Status != StatusConnected {
		t.Fatalf("switched: %+v sig=%q %s", g.Request.Browser, g.ChromeSitesSig, g.Status)
	}
}
