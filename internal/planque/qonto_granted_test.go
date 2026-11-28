package planque

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

// What a bank's token answer granted beyond the reviewed list: only for a
// sign-in on Qonto's host, any name not on the list, once each; an answer
// that names no scope grants what was asked.
func TestWhatABankGrantedBeyondTheList(t *testing.T) {
	asked := strings.Join(qontoScopes, " ")
	for _, c := range []struct {
		server, granted string
		want            []string
	}{
		{qontoServer, asked, nil},
		{qontoServer, "", nil},
		{qontoServer, "offline_access organization.read", nil},
		{qontoServer, asked + " request_transfers.write card.write request_transfers.write", []string{"request_transfers.write", "card.write"}},
		{"https://MCP.qonto.com/mcp/", "organization.read payment_link.write", []string{"payment_link.write"}},
		{"https://mcp.stripe.com", "anything at_all", nil},
	} {
		if got := widerThanAsked(c.server, c.granted); !slices.Equal(got, c.want) {
			t.Errorf("%s %q: %v, want %v", c.server, c.granted, got, c.want)
		}
	}
}

// signInExpectingFailure follows a sign-in link that planque refuses at
// its callback: the page says it didn't connect, and the app isn't opened.
func (e *env) signInExpectingFailure(link string) {
	e.t.Helper()
	client := e.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(link)
	if err != nil {
		e.t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode == 200 || string(page) == "app opened" {
		e.t.Fatalf("the refused sign-in reached the app: %d %s", resp.StatusCode, page)
	}
}

// asQonto makes the grant id a sign-in to Qonto's server, as a direct
// Qonto entry is, while the fake vendor answers it.
func (e *env) asQonto(id string) {
	e.t.Helper()
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	g.Request.Server = qontoServer
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		e.t.Fatal(err)
	}
}

// A sign-in Qonto answers with more than planque asked for is refused: the
// token is never kept, nothing connects, and the owner reads why, naming
// what was granted. One granted the list is kept with its scope names.
func TestABankSignInGrantedMoreIsRefused(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	e.asQonto(id)
	e.vendor.Scope = strings.Join(qontoScopes, " ") + " request_transfers.write"
	e.signInExpectingFailure(out["sign_in_url"].(string))
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "granted more than tafs asked for (request_transfers.write)") ||
		!strings.Contains(g.Reason, "Remove the connection in") {
		t.Fatalf("grant: %s %q", g.Status, g.Reason)
	}
	if _, _, err := e.s.openToken(context.Background(), g); !errors.Is(err, errNotFound) {
		t.Fatalf("a token was kept: %v", err)
	}
	if len(g.GrantedScopes) != 0 || len(e.kiff.Connections(tenant)) != 0 {
		t.Fatalf("granted %v, connections %d", g.GrantedScopes, len(e.kiff.Connections(tenant)))
	}

	// A sign-in granted no more than asked: kept, with its scope names on the
	// grant (names only), for any vendor.
	e2 := newEnv(t)
	e2.enrol()
	id2, digest2 := e2.request()
	st, out = e2.do("POST", "/v1/accounts/"+tenant+"/grants/"+id2+"/approve", map[string]string{"digest": digest2}, e2.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	e2.vendor.Scope = "offline_access organization.read"
	code := e2.signIn(out["sign_in_url"].(string))
	if st, out := e2.do("POST", "/v1/accounts/"+tenant+"/grants/"+id2+"/confirm", map[string]string{"confirm": code}, e2.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	_, view := e2.do("GET", "/v1/accounts/"+tenant+"/grants/"+id2, nil, nil, "k_setup")
	if got, _ := view["granted_scopes"].([]any); len(got) != 2 || got[0] != "offline_access" || got[1] != "organization.read" {
		t.Fatalf("granted_scopes: %v", view["granted_scopes"])
	}
}

// A refresh that grants Qonto's sign-in more than planque asked for ends
// it: the stored token is dropped, the grant waits for the owner to sign in
// again with the cause scope_wider, and no call goes out with it.
func TestABankRefreshGrantedMoreSignsOut(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.asQonto(g.ID)
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	ctx := context.Background()
	e.vendor.Expire()
	tok, it, _ := e.s.openToken(ctx, g)
	tok.Expiry = time.Now().Add(-time.Second)
	_ = e.s.storeToken(ctx, g, tok, &it.Ver)
	e.vendor.Scope = "organization.read sepa_direct_debit.write"
	if _, _, err := e.s.accessToken(ctx, g, 0); err == nil {
		t.Fatal("a refresh granting more was used")
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if now.Status != StatusNeedsReconnect || now.SignedOutCause != CauseScopeWider || !strings.Contains(now.Reason, "sepa_direct_debit.write") {
		t.Fatalf("grant: %s %q %q", now.Status, now.SignedOutCause, now.Reason)
	}
	if _, _, err := e.s.openToken(ctx, now); !errors.Is(err, errNotFound) {
		t.Fatalf("the token was kept: %v", err)
	}
}
