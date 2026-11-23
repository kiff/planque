package planque

import (
	"context"
	"strings"
	"testing"
	"time"

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

	// A connection the owner didn't approve, or that ended, or that isn't on
	// the owner's own Chrome, can't be switched; nor can any while the Chrome
	// isn't offered to the account. Each refused, and nothing changes.
	back := map[string]any{"chrome_id": c.r.chrome, "chrome_sites_sig": testSitesSig}
	// Each refusal is its own signed request (a second apart: the same request
	// within a second is a replay), so the answer is the check's own.
	refused := func(what, want string) {
		t.Helper()
		time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
		if st, out := c.do("PUT", path, back, c.phone, ""); st != 409 || !strings.Contains(mustJSONString(out), want) {
			t.Errorf("%s: %d %v, want 409 %q", what, st, out, want)
		}
	}
	const notApproved, notOwnChrome = "only a connection the owner approved", "doesn't use the owner's own Chrome"
	for what, tc := range map[string]struct {
		edit func(g *Grant)
		want string
	}{
		"ended":                      {func(g *Grant) { g.Status = StatusDisconnected }, notApproved},
		"failed":                     {func(g *Grant) { g.Status = StatusFailed }, notApproved},
		"only asked for (pending)":   {func(g *Grant) { g.Status = StatusPending }, notApproved},
		"not the owner's own Chrome": {func(g *Grant) { g.Request.Browser.Driver = "" }, notOwnChrome},
		"not a browser's (a key)":    {func(g *Grant) { g.Request.Browser = nil; g.Request.Kind = KindAPIKey }, notOwnChrome},
	} {
		g, _ := c.s.loadGrant(context.Background(), tenant, id)
		was := *g
		wasBrowser := *g.Request.Browser
		tc.edit(g)
		if err := c.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
		refused(what, tc.want)
		after, _ := c.s.loadGrant(context.Background(), tenant, id)
		after.Status, after.Request.Kind, after.Request.Browser = was.Status, was.Request.Kind, &wasBrowser
		if err := c.s.saveGrant(context.Background(), after); err != nil {
			t.Fatal(err)
		}
	}
	c.s.ChromeOffer = ChromeOffer{}
	refused("the owner's Chrome not offered", "isn't offered to this account")
	if g, _ := c.s.loadGrant(context.Background(), tenant, id); g.Request.Browser.Chrome != home.ID || g.ChromeSitesSig != newSig {
		t.Fatalf("a refused switch changed the connection: %+v", g.Request.Browser)
	}
}

// A connection in the owner's own Chrome offers "Use another Chrome" exactly
// where the PUT would move it ("Your Chrome" 2): the action and the switch
// ask one rule (chromeSwitchRefusal), so no app keeps a status list, and a
// button is never offered that answers 409.
// Category: authority
// Attack: an action offered for a connection the switch refuses (pending, failed, ended, not the owner's Chrome, Chrome not offered, no other Chrome)
// Refused: each shows no action and the PUT refuses it; where the action shows, the PUT moves it
func TestTheUseAnotherChromeActionShowsExactlyWhereTheMoveIsAccepted(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_wiki", "agent_name": "Wiki titles",
		"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	home := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	if err := putDoc(context.Background(), c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(home.ID)}, home, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	path := "/v1/accounts/" + tenant + "/grants/" + id + "/chrome"
	// set puts the connection back on the rig's Chrome, as edit leaves it.
	set := func(edit func(g *Grant)) {
		t.Helper()
		g, _ := c.s.loadGrant(context.Background(), tenant, id)
		g.Request.Browser.Driver, g.Request.Browser.Chrome, g.ChromeSitesSig, g.Status = DriverExtension, c.r.chrome, testSitesSig, StatusConnected
		edit(g)
		if err := c.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
	}
	n := 0
	// check reads the row, then asks the PUT to move it to home: the action
	// shows if and only if the move is accepted.
	check := func(what string, wantMove bool) {
		t.Helper()
		row := listRow(t, c.env, id, "")
		a := action(row, ActionUseAnotherChrome)
		// Each switch its own signed request: a different signature each time, never a replay.
		n++
		sig := testSitesSig[:len(testSitesSig)-1] + string(rune('a'+n))
		st, out := c.do("PUT", path, map[string]any{"chrome_id": home.ID, "chrome_sites_sig": sig}, c.phone, "")
		moved := st == 200
		if (a != nil) != moved || moved != wantMove {
			t.Errorf("%s: action %v, PUT %d %v; want a move: %v", what, a, st, out, wantMove)
		}
		if a != nil && (a["words"] != "Use another Chrome" || a["how"] != "chrome") {
			t.Errorf("%s: the action %v", what, a)
		}
	}
	for _, tc := range []struct {
		status string
		move   bool
	}{
		// signed_in isn't here: one with no sign-in time is expired as it's
		// read, so both say no (they agree; its expiry is tested on its own).
		{StatusConnected, true}, {StatusApproved, true}, {StatusListed, true}, {StatusNeedsReconnect, true},
		{StatusPending, false}, {StatusFailed, false},
	} {
		set(func(g *Grant) { g.Status = tc.status })
		check(tc.status, tc.move)
	}
	// An ended connection isn't in the owner's list at all, so it offers
	// nothing; the PUT refuses it.
	set(func(g *Grant) { g.Status = StatusDisconnected })
	n++
	if st, out := c.do("PUT", path, map[string]any{"chrome_id": home.ID, "chrome_sites_sig": testSitesSig[:len(testSitesSig)-1] + string(rune('a'+n))}, c.phone, ""); st != 409 {
		t.Errorf("ended: PUT %d %v, want 409", st, out)
	}
	set(func(g *Grant) { g.Request.Browser.Driver = "" })
	check("not the owner's own Chrome", false)
	set(func(g *Grant) {})
	c.s.ChromeOffer = ChromeOffer{}
	check("the owner's Chrome not offered", false)
	c.s.ChromeOffer = ChromeOffer{All: true}
	if _, _, err := c.s.KV.Delete(context.Background(), tenantPK(tenant), skChrome(home.ID)); err != nil {
		t.Fatal(err)
	}
	set(func(g *Grant) {})
	check("no other Chrome paired", false)
}
