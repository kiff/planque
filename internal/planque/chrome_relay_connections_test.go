package planque

import (
	"context"
	"testing"
)

// The owner's Chrome learns which of its connections are live (the side
// panel offers those tafs to talk to): in the welcome, and again whenever one
// is approved, ends or moves, posted to its open relay. For display only, never
// authority: a step still carries its own signed sites.
// Category: authority
// Attack: a list that names a connection that ended, another Chrome's, or one posted to a connection that isn't past its welcome
// Refused: each; the list is the Chromes list's own rule, per Chrome
func TestTheOwnersChromeHearsItsLiveConnections(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	st, out := c.do("POST", "/v1/accounts/"+tenant+"/catalog/browser@1/grants", map[string]any{"agent": "agt_wiki", "agent_name": "Wiki titles",
		"tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseRead}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)

	// Before the welcome, nothing is posted: there is no live connection yet.
	c.s.tellChromes(context.Background(), tenant, c.r.chrome)
	if n := len(c.r.poster.posts["c1"]); n != 0 {
		t.Fatalf("posted before any connection: %d", n)
	}

	// The welcome names it, with its agent and sites.
	c.r.connect(t, "c1")
	w := c.r.poster.last("c1")
	list, _ := w["connections"].([]any)
	if len(list) != 1 {
		t.Fatalf("the welcome's connections: %v", w["connections"])
	}
	one := list[0].(map[string]any)
	if one["grant"] != id || one["agent"] != "agt_wiki" || one["name"] != "Wiki titles" {
		t.Fatalf("the connection: %v", one)
	}
	if sites, _ := one["sites"].([]any); len(sites) != 1 || sites[0] != "wikipedia.org" {
		t.Fatalf("its sites: %v", one["sites"])
	}

	// It ends: the Chrome hears an empty list, never null.
	g, _ := c.s.loadGrant(context.Background(), tenant, id)
	if err := c.s.endGrant(context.Background(), g, "", "the test"); err != nil {
		t.Fatal(err)
	}
	m := c.r.poster.last("c1")
	if m["type"] != "connections" {
		t.Fatalf("after the end: %v", m)
	}
	if l, ok := m["connections"].([]any); !ok || len(l) != 0 {
		t.Fatalf("an ended connection is still listed: %v", m["connections"])
	}

	// Another Chrome's connection is never in this Chrome's list.
	if got := c.s.connectionsOn(context.Background(), tenant, "chr_otherchrome01"); len(got) != 0 {
		t.Fatalf("another Chrome's list: %v", got)
	}
}
