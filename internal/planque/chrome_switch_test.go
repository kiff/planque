package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// Switching a browser connection to another Chrome ("Your Chrome" 2): the
// owner's phone signs it, with the connection's sites signed for the new
// Chrome; the sites and the rest stay, and the step goes to the new Chrome.
// Category: authority
// Attack: a switch to another account's Chrome, to a Chrome that doesn't exist, with no sites signature, on a connection that ended, unsigned
// Refused: each; the signed switch moves it once, with the new signature
func TestABrowserConnectionSwitchesToAnotherChrome(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_wiki", "agent_name": "Wiki titles",
		"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	g, _ := c.s.loadGrant(context.Background(), tenant, id)
	g.ChromeSitesSig, g.Status = testSitesSig, StatusConnected
	if err := c.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	put := func(doc Chrome, pk string) {
		if err := putDoc(context.Background(), c.s.KV, store.Item{PK: pk, SK: skChrome(doc.ID)}, doc, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	home := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	put(home, tenantPK(tenant))
	theirs := Chrome{ID: "chr_otheraccount1", Tenant: "t_other", Name: "Theirs", Key: "k"}
	put(theirs, tenantPK(theirs.Tenant))
	newSig := strings.Replace(testSitesSig, "s1gned", "s1GNED", 1)
	path := "/v1/accounts/" + tenant + "/grants/" + id + "/chrome"

	for what, tc := range map[string]struct {
		body  map[string]any
		phone bool
		want  int
	}{
		"unsigned":                   {map[string]any{"chrome_id": home.ID, "chrome_sites_sig": newSig}, false, 401},
		"no sites signature":         {map[string]any{"chrome_id": home.ID}, true, 400},
		"another account's Chrome":   {map[string]any{"chrome_id": theirs.ID, "chrome_sites_sig": newSig}, true, 404},
		"a Chrome that isn't there":  {map[string]any{"chrome_id": "chr_nosuchchrome0", "chrome_sites_sig": newSig}, true, 404},
		"the Chrome it's already on": {map[string]any{"chrome_id": c.r.chrome, "chrome_sites_sig": newSig}, true, 409},
	} {
		phone := c.phone
		if !tc.phone {
			phone = nil
		}
		if st, out := c.do("PUT", path, tc.body, phone, ""); st != tc.want {
			t.Errorf("%s: %d %v, want %d", what, st, out, tc.want)
		}
	}
	if g, _ := c.s.loadGrant(context.Background(), tenant, id); g.Request.Browser.Chrome != c.r.chrome || g.ChromeSitesSig != testSitesSig {
		t.Fatalf("a refused switch changed the connection: %+v", g.Request.Browser)
	}

	// The signed switch: the new Chrome, its name, the new signature; the sites stay.
	if st, out := c.do("PUT", path, map[string]any{"chrome_id": home.ID, "chrome_sites_sig": newSig}, c.phone, ""); st != 200 {
		t.Fatalf("switch: %d %v", st, out)
	}
	g, _ = c.s.loadGrant(context.Background(), tenant, id)
	if b := g.Request.Browser; b.Chrome != home.ID || b.ChromeName != "Home laptop" || g.ChromeSitesSig != newSig || strings.Join(b.Sites, ",") != "wikipedia.org" {
		t.Fatalf("after the switch: %+v sig=%q", b, g.ChromeSitesSig)
	}
	if last := g.Events[len(g.Events)-1]; !strings.Contains(last.What, `switched to the Chrome "Home laptop"`) {
		t.Errorf("the event: %+v", last)
	}

	// A connection that ended can't be switched.
	g.Status = StatusDisconnected
	if err := c.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	if st, out := c.do("PUT", path, map[string]any{"chrome_id": c.r.chrome, "chrome_sites_sig": testSitesSig}, c.phone, ""); st != 409 {
		t.Errorf("an ended connection: %d %v", st, out)
	}
}
