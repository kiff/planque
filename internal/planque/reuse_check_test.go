package planque

import (
	"context"
	"strings"
	"testing"
)

// A sign-in to reuse that the vendor refused (signed out: a password
// change, a revoked app) isn't offered as "no second sign-in": the request
// signs in anew to the same server on its first card, says why, and no
// grant fails after the owner approves.
func TestAReuseOfARefusedSignInSignsInAnew(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t)
	first := e.connect()
	e.vendor.SignOut()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "list_charges", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse of a refused sign-in: %d %v", st, out)
	}
	text, _ := out["text"].(string)
	if strings.Contains(text, "already signed in") || !strings.Contains(text, "You will sign in at") {
		t.Fatalf("the first card: %q", text)
	}
	id := out["id"].(string)
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil || g.Status != StatusPending || g.TokenOf != "" || g.SignedIn || g.Request.Reuse != "" {
		t.Fatalf("the new grant: %v status %s tokenOf %q signedIn %v reuse %q", err, g.Status, g.TokenOf, g.SignedIn, g.Request.Reuse)
	}
	if !strings.Contains(g.Events[len(g.Events)-1].What, "was refused") {
		t.Fatalf("its record: %+v", g.Events)
	}
	// The refused one says so, as any refused refresh does.
	if old, _ := e.s.loadGrant(ctx, tenant, first.ID); old.Status != StatusNeedsReconnect {
		t.Fatalf("the refused sign-in's grant: %s", old.Status)
	}
	// Approving it gives the vendor's sign-in link, as any first sign-in.
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
	if st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("approve: %d %v", st, out)
	}
}

// A live sign-in is reused as before: the check reads the server once and
// costs the owner nothing.
func TestALiveSignInIsStillReused(t *testing.T) {
	e := newEnv(t)
	first := e.connect()
	lists := e.vendor.Lists
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "list_charges", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 || !strings.Contains(out["text"].(string), "already signed in") {
		t.Fatalf("reuse: %d %v", st, out)
	}
	if e.vendor.Lists != lists+1 {
		t.Fatalf("the check read the server %d times", e.vendor.Lists-lists)
	}
}
