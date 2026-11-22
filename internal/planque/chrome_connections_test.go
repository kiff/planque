package planque

import (
	"context"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// The Chromes list says each Chrome's live connections ("Your Chrome" 3 and 4):
// the taf's name and its sites, decided here once so no app writes the rule
// again. One the owner has yet to approve is listed (unpairing ends it too);
// an ended or failed one isn't, nor another Chrome's, nor one not driven by
// the owner's Chrome, and a Chrome with none lists an empty list.
func TestTheChromesListSaysEachChromesLiveConnections(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ask := func(agent, name string, edit func(g *Grant)) string {
		t.Helper()
		st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": agent, "agent_name": name,
			"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("ask %s: %d %v", name, st, out)
		}
		id := out["grant"].(map[string]any)["id"].(string)
		if edit != nil {
			g, _ := c.s.loadGrant(context.Background(), tenant, id)
			edit(g)
			if err := c.s.saveGrant(context.Background(), g); err != nil {
				t.Fatal(err)
			}
		}
		return id
	}
	live := ask("agt_wiki", "Wiki lookup", nil)
	ask("agt_ended", "Ended", func(g *Grant) { g.Status = StatusDisconnected })
	ask("agt_failed", "Failed", func(g *Grant) { g.Status = StatusFailed })
	ask("agt_hosted", "Hosted", func(g *Grant) { g.Request.Browser.Driver = "" })
	pending := ask("agt_asked", "", func(g *Grant) { g.Status = StatusPending })
	home := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	other := ask("agt_home", "On the laptop", func(g *Grant) { g.Request.Browser.Chrome = home.ID })
	spare := Chrome{ID: "chr_sparelaptop1", Tenant: tenant, Name: "Spare", Key: "k"}
	for _, doc := range []Chrome{home, spare} {
		if err := putDoc(context.Background(), c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(doc.ID)}, doc, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}

	st, out := c.do("GET", "/v1/accounts/"+tenant+"/chromes", nil, c.phone, "")
	if st != 200 {
		t.Fatalf("list: %d %v", st, out)
	}
	by := map[string][]any{}
	for _, row := range out["chromes"].([]any) {
		m := row.(map[string]any)
		by[m["chrome_id"].(string)] = m["connections"].([]any)
	}
	work := by[c.r.chrome]
	if len(work) != 2 || work[0].(map[string]any)["grant"] == work[1].(map[string]any)["grant"] {
		t.Fatalf("the paired Chrome's connections: %v", work)
	}
	if work[0].(map[string]any)["grant"] == pending {
		work[0], work[1] = work[1], work[0]
	}
	// With no name given, the agent's id names it.
	if work[1].(map[string]any)["grant"] != pending || work[1].(map[string]any)["name"] != "agt_asked" || work[0].(map[string]any)["grant"] != live || work[0].(map[string]any)["name"] != "Wiki lookup" ||
		work[0].(map[string]any)["sites"].([]any)[0] != "wikipedia.org" {
		t.Fatalf("the paired Chrome's connections: %v", work)
	}
	if len(by[home.ID]) != 1 || by[home.ID][0].(map[string]any)["grant"] != other {
		t.Fatalf("the laptop's connections: %v", by[home.ID])
	}
	if none, ok := by[spare.ID]; !ok || none == nil || len(none) != 0 {
		t.Fatalf("a Chrome with none: %v (listed %v)", none, ok)
	}
}
