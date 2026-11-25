package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// Re-pairing replaces (kiff/tafs docs/chrome.md, "Your Chrome" 1): the confirm
// names the Chrome it replaces and carries the connections that move, each
// with the sites the owner approved; those move to the new Chrome, the old
// Chrome goes, and every connection it didn't carry ends with it. The
// extension learns what was carried, inside what the phone signed for it.
// Category: authority
// Attack: carrying a connection that isn't on the Chrome replaced, carrying it with other sites than approved, carrying without a replace
// Refused: each, before anything is kept; a carried connection's step says its sites are vouched for by the pairing
func TestRePairingMovesTheCarriedConnectionsAndEndsTheRest(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ask := func(name string) string {
		t.Helper()
		st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_" + strings.ToLower(strings.ReplaceAll(name, " ", "_")), "agent_name": name,
			"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("ask %s: %d %v", name, st, out)
		}
		id := out["grant"].(map[string]any)["id"].(string)
		return id
	}
	// approved: the owner's phone signed its sites and it's live, as an approve leaves it.
	approved := func(id string, edit func(g *Grant)) string {
		t.Helper()
		g, _ := c.s.loadGrant(context.Background(), tenant, id)
		g.ChromeSitesSig, g.Status = testSitesSig, StatusConnected
		if edit != nil {
			edit(g)
		}
		if err := c.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
		return id
	}
	keep, drop := approved(ask("Wiki titles"), nil), approved(ask("Wiki links"), nil)
	old := c.r.chrome
	// Connections a replace must never carry.
	asked := approved(ask("Only asked"), func(g *Grant) { g.Status = StatusPending })
	ended := approved(ask("Ended"), func(g *Grant) { g.Status = StatusDisconnected })
	hosted := approved(ask("Hosted"), func(g *Grant) { g.Request.Browser.Driver = "" })
	elsewhere := approved(ask("Elsewhere"), func(g *Grant) { g.Request.Browser.Chrome = "chr_homelaptop01" })
	unsigned := approved(ask("Unsigned"), func(g *Grant) { g.ChromeSitesSig = "" })
	carry := func(ids ...string) []any {
		out := []any{}
		for _, id := range ids {
			out = append(out, map[string]any{"grant": id, "sites": []string{"wikipedia.org"}})
		}
		return out
	}
	// A Chrome paired to another account, never this one's to replace.
	theirs := Chrome{ID: "chr_otheraccount1", Tenant: "t_other", Name: "Theirs", Key: "k"}
	if err := putDoc(context.Background(), c.s.KV, store.Item{PK: tenantPK(theirs.Tenant), SK: skChrome(theirs.ID)}, theirs, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	// More than 50, each distinct: refused for the count (400), before any is read.
	many := []string{}
	for i := 0; i < chromeCarriedMax+1; i++ {
		many = append(many, "grt_many"+strings.Repeat("x", i%10)+string(rune('a'+i%26)))
	}

	acct := "/v1/accounts/" + tenant
	open := func() (string, string) {
		t.Helper()
		st, out := c.do("POST", acct+"/chromes/pairings", map[string]any{}, c.phone, "")
		if st != 201 {
			t.Fatalf("open: %d %v", st, out)
		}
		pairing, lookup := out["pairing"].(string), out["lookup"].(string)
		extKey, extSum := p256Key(t)
		if st, out := extensionPosts(c.env, lookup, extKey, releasedManifest); st != 202 {
			t.Fatalf("post: %d %v", st, out)
		}
		return pairing, extSum + "|" + extKey + "|" + lookup
	}
	confirm := func(pairing, sum string, extra map[string]any) (int, map[string]any) {
		body := map[string]any{"chrome_key_sha256": strings.Split(sum, "|")[0], "name": "Work Mac 2", "phone_key": strings.Split(sum, "|")[1],
			"phone_proof": b64(32), "phone_sig": b64(64)}
		for k, v := range extra {
			body[k] = v
		}
		return c.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", body, c.phone, "")
	}

	// Refused before anything is kept: carried with no replace, a connection
	// not on the Chrome replaced, or other sites than it was approved with.
	pairing, sum := open()
	for what, extra := range map[string]map[string]any{
		"carried with no replace": {"carried": []any{map[string]any{"grant": keep, "sites": []string{"wikipedia.org"}}}},
		"other sites":             {"replaces": old, "carried": []any{map[string]any{"grant": keep, "sites": []string{"wikipedia.org", "bank.example"}}}},
		"not a connection":        {"replaces": old, "carried": []any{map[string]any{"grant": "grt_nope", "sites": []string{"wikipedia.org"}}}},
		"only asked for":          {"replaces": old, "carried": carry(asked)},
		"already ended":           {"replaces": old, "carried": carry(ended)},
		"not the extension's":     {"replaces": old, "carried": carry(hosted)},
		"on another Chrome":       {"replaces": old, "carried": carry(elsewhere)},
		"sites never signed":      {"replaces": old, "carried": carry(unsigned)},
		"the same one twice":      {"replaces": old, "carried": carry(keep, keep)},
		"more than 50":            {"replaces": old, "carried": carry(many...)},
		"no such Chrome":          {"replaces": "chr_nosuchchrome0", "carried": carry(keep)},
		"another account's":       {"replaces": theirs.ID, "carried": carry(keep)},
		"half the one-touch":      {"tenant": tenant},
		"another account":         {"tenant": "t_other", "signed_at": "2026-11-22T20:00:00Z", "kiff_request_sha256": strings.Repeat("a", 64)},
	} {
		st, out := confirm(pairing, sum, extra)
		if what == "more than 50" && st != 400 || st != 400 && st != 409 {
			t.Errorf("%s: %d %v", what, st, out)
		}
	}
	if g, _ := c.s.loadGrant(context.Background(), tenant, keep); g.Request.Browser.Chrome != old {
		t.Fatalf("a refused confirm moved a connection")
	}

	// The replace: one connection carried, the other not.
	st, out := confirm(pairing, sum, map[string]any{"replaces": old, "carried": []any{map[string]any{"grant": keep, "sites": []string{"wikipedia.org"}}},
		"tenant": tenant, "signed_at": "2026-11-22T20:00:00Z", "kiff_request_sha256": strings.Repeat("b", 64)})
	if st != 200 || len(out["moved"].([]any)) != 1 || out["moved"].([]any)[0] != keep {
		t.Fatalf("confirm with a replace: %d %v", st, out)
	}
	fresh := out["chrome_id"].(string)
	moved, _ := c.s.loadGrant(context.Background(), tenant, keep)
	if b := moved.Request.Browser; b.Chrome != fresh || b.ChromeName != "Work Mac 2" || moved.ChromeSitesSig != "" || moved.ChromeSitesCarried != fresh {
		t.Fatalf("the carried connection: %+v sig=%q carried=%q", b, moved.ChromeSitesSig, moved.ChromeSitesCarried)
	}
	if ended, _ := c.s.loadGrant(context.Background(), tenant, drop); ended.Status != StatusDisconnected {
		t.Errorf("the connection not carried is %s, want it ended with the old Chrome", ended.Status)
	}
	if st, out := c.do("GET", acct+"/chromes", nil, c.phone, ""); st != 200 || len(out["chromes"].([]any)) != 1 {
		t.Errorf("the old Chrome is still listed: %d %v", st, out)
	}
	// The extension collects what was carried and the one-touch fields, as the phone sent them.
	lookup := strings.Split(sum, "|")[2]
	st, out = c.do("GET", "/v1/extension/pairings/"+lookup+"?key="+strings.Split(sum, "|")[1], nil, nil, "")
	if st != 200 || out["replaces"] != old || out["tenant"] != tenant || out["kiff_request_sha256"] != strings.Repeat("b", 64) || len(out["carried"].([]any)) != 1 {
		t.Fatalf("collect: %d %v", st, out)
	}
	// A step of the carried connection says the pairing vouches for its sites.
	sess, code, msg := c.s.startChrome(context.Background(), moved, "op_1", "agt_taf", "")
	if sess == nil || !sess.signed.Carried || sess.signed.Sig != "" {
		t.Fatalf("the carried connection's session: %v %s %s", sess, code, msg)
	}
}
