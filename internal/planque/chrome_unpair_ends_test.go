package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// Unpairing a Chrome ends its connections (kiff/tafs docs/chrome.md, "Your
// Chrome" 3): none is left on a Chrome that's gone, and each says why. A
// connection on the account's other Chrome, and one that isn't a browser's,
// stays as it was.
// Category: authority
// Attack: a connection left pointing at an unpaired Chrome, reachable again by a Chrome paired later under the same id
// Refused: none is left: each of the Chrome's connections is disconnected with the unpair
func TestUnpairingAChromeEndsItsConnections(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ask := func(name string) string {
		t.Helper()
		st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_" + strings.ToLower(strings.ReplaceAll(name, " ", "_")), "agent_name": name,
			"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("ask %s: %d %v", name, st, out)
		}
		return out["grant"].(map[string]any)["id"].(string)
	}
	mine, kept := ask("Wiki titles"), ask("Wiki links")
	// The second one is on the account's other Chrome.
	other := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	if err := putDoc(context.Background(), c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(other.ID)}, other, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	g, _ := c.s.loadGrant(context.Background(), tenant, kept)
	g.Request.Browser.Chrome, g.Request.Browser.ChromeName = other.ID, other.Name
	if err := c.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}

	if st, out := c.do("DELETE", "/v1/accounts/"+tenant+"/chromes/"+c.r.chrome, nil, c.phone, ""); st != 200 {
		t.Fatalf("unpair: %d %v", st, out)
	}
	ended, _ := c.s.loadGrant(context.Background(), tenant, mine)
	if ended.Status != StatusDisconnected {
		t.Fatalf("the unpaired Chrome's connection is %s, want %s", ended.Status, StatusDisconnected)
	}
	if last := ended.Events[len(ended.Events)-1]; !strings.Contains(last.What, "its Chrome was unpaired") {
		t.Errorf("it doesn't say why it ended: %+v", last)
	}
	if still, _ := c.s.loadGrant(context.Background(), tenant, kept); still.Status == StatusDisconnected {
		t.Errorf("the other Chrome's connection ended too")
	}
}
