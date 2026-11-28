package planque

import (
	"context"
	"testing"
	"time"
)

// An hour after its sign-in, a connected grant's access token has expired
// at the vendor (Qonto's lasts less: 28 Nov, the owner's Version 2 step 3
// asked for a second sign-in after "no second sign-in"). Adding a tool
// lists the vendor's tools with a refreshed token, as a call would: the
// addition is declared while it waits, and one approval connects it, with
// no sign-in link.
func TestAnAdditionAfterTheAccessTokenExpiredNeedsNoSignIn(t *testing.T) {
	for _, c := range []struct {
		name string
		// stale: planque's record still says the token is good (the vendor
		// ended it early), so the vendor's refusal triggers the refresh;
		// else the record says it expired, and planque refreshes first.
		stale bool
	}{{"the vendor refuses a token planque thinks good", true}, {"the token's own time passed", false}} {
		t.Run(c.name, func(t *testing.T) {
			e := newEnv(t)
			id := e.connectedRefunds()
			acct := "/v1/accounts/" + tenant + "/grants/" + id
			e.vendor.Expire()
			if !c.stale {
				ctx := context.Background()
				g, _ := e.s.loadGrant(ctx, tenant, id)
				tok, it, err := e.s.openToken(ctx, g)
				if err != nil {
					t.Fatal(err)
				}
				tok.Expiry = time.Now().Add(-time.Second)
				_ = e.s.storeToken(ctx, g, tok, &it.Ver)
			}
			list := map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}
			st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup")
			if st != 200 || out["status"] != StatusPending || out["declared"] != true || out["sign_in"] != "none" {
				t.Fatalf("add after the token expired: %d %v", st, out)
			}
			live, _ := e.vendor.Live()
			st, out = e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
			if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] != nil {
				t.Fatalf("approve the addition: %d %v", st, out)
			}
			if now, _ := e.vendor.Live(); now != live {
				t.Fatalf("a new sign-in was made: %d live tokens, were %d", now, live)
			}
			if tools, _ := out["tools"].([]any); len(tools) != 2 {
				t.Fatalf("tools: %v", out["tools"])
			}
		})
	}
}

// The same when the addition wasn't declared while it waited (KIFF down
// then): the approval lists with a refreshed token and connects it.
func TestAnUndeclaredAdditionApprovedAfterTheTokenExpiredConnects(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	e.kiff.ConnectDown = true
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 200 || out["declared"] == true {
		t.Fatalf("add: %d %v", st, out)
	}
	e.kiff.ConnectDown = false
	e.vendor.Expire()
	st, out = e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] != nil {
		t.Fatalf("approve: %d %v", st, out)
	}
}

// A refresh the vendor refuses signs the grant out before any tap: the
// addition goes, the grant keeps what it had and waits signed out with its
// reason, and what the owner reads says a sign-in is needed (sign_in:
// "needed"), so approving it then gives a sign-in link the words said.
func TestAnAdditionWhoseRefreshIsRefusedSignsOutBeforeTheTap(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	e.vendor.SignOut()
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 200 || out["declared"] == true {
		t.Fatalf("add with a refused refresh: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusNeedsReconnect || g.SignedOutCause != CauseRefreshRefused || len(g.Pending) != 0 || g.Connected != nil || len(g.Request.Tools) != 1 {
		t.Fatalf("after a refused refresh: %s %q, pending %v, connected %v, asked %d", g.Status, g.SignedOutCause, g.Pending, g.Connected, len(g.Request.Tools))
	}
	_, view := e.do("GET", acct, nil, nil, "k_setup")
	if view["sign_in"] != "needed" {
		t.Fatalf("the view: %v", view["sign_in"])
	}
	// Approved as it now reads: a sign-in link, which its words said. (The
	// same digest as its first approval: past the second a replay is in.)
	time.Sleep(1100 * time.Millisecond)
	if st, out := e.do("POST", acct+"/approve", map[string]string{"digest": g.Digest}, e.phone, ""); st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("approve the signed-out grant: %d %v", st, out)
	}
}

// sign_in is "needed" for a grant not signed in yet, "none" for one whose
// held sign-in connects tools on the approval, and absent for a key's or a
// browser's grant.
func TestTheGrantSaysWhetherApprovingItTakesASignIn(t *testing.T) {
	for _, c := range []struct {
		g    Grant
		want string
	}{
		{Grant{Status: StatusPending}, "needed"},
		{Grant{Status: StatusApproved}, "needed"},
		{Grant{Status: StatusPending, SignedIn: true}, "none"},
		{Grant{Status: StatusConnected, SignedIn: true}, "none"},
		{Grant{Status: StatusNeedsReconnect, SignedIn: true}, "needed"},
		{Grant{Status: StatusFailed}, ""},
		{Grant{Status: StatusPending, Request: Request{Kind: KindBrowser}}, ""},
		{Grant{Status: StatusPending, Request: Request{Kind: KindAPIKey}}, ""},
	} {
		if got := signInNeed(&c.g); got != c.want {
			t.Errorf("%s signed_in=%v kind=%q: %q, want %q", c.g.Status, c.g.SignedIn, c.g.Request.Kind, got, c.want)
		}
	}
}
