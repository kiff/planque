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

// What a bank's token answer granted beyond what planque asked for: only
// for a sign-in on Qonto's host, any name the sign-in didn't ask for, once
// each, whether it asked for the reviewed list or for all Qonto lists; an
// answer that names no scope grants what was asked.
func TestWhatABankGrantedBeyondWhatWasAsked(t *testing.T) {
	grant := func(server string, asked []string) *Grant {
		g := &Grant{}
		g.Request.Server, g.Vendor.Scopes = server, asked
		return g
	}
	reviewed, all := strings.Join(qontoScopes, " "), strings.Join(qontoListed, " ")
	for _, c := range []struct {
		g       *Grant
		granted string
		want    []string
	}{
		{grant(qontoServer, qontoScopes), reviewed, nil},
		{grant(qontoServer, qontoScopes), "", nil},
		{grant(qontoServer, qontoScopes), "offline_access organization.read", nil},
		{grant(qontoServer, qontoScopes), reviewed + " request_transfers.write card.write request_transfers.write", []string{"request_transfers.write", "card.write"}},
		{grant("https://MCP.qonto.com/mcp/", qontoScopes), "organization.read payment_link.write", []string{"payment_link.write"}},
		// Asked for all Qonto lists: all of it is as asked, and a name it doesn't list isn't.
		{grant(qontoServer, qontoListed), all, nil},
		{grant(qontoServer, qontoListed), all + " admin.everything", []string{"admin.everything"}},
		{grant("https://mcp.stripe.com", []string{"a"}), "anything at_all", nil},
	} {
		if got := widerThanAsked(c.g, c.granted); !slices.Equal(got, c.want) {
			t.Errorf("%s asked %d, granted %q: %v, want %v", c.g.Request.Server, len(c.g.Vendor.Scopes), c.granted, got, c.want)
		}
	}
}

// An account on PLANQUE_FULL_SCOPES asks Qonto for every scope it lists, as
// any server is asked; another account, and every account with it unset,
// keeps the reviewed list.
func TestFullScopesAskQontoForEverything(t *testing.T) {
	s := &Server{FullScopes: map[string]bool{"t_full": true}}
	if !s.fullScopesFor("t_full") || s.fullScopesFor("t_other") || (&Server{}).fullScopesFor("t_full") {
		t.Fatal("full scopes for the wrong accounts")
	}
	if !(&Server{FullScopes: map[string]bool{"*": true}}).fullScopesFor("anyone") {
		t.Fatal("all doesn't cover every account")
	}
	// Its grants on Qonto's host get Qonto's whole live tool list, for the
	// Card to decide; another account keeps the reviewed tools, and another
	// host is unchanged.
	if !s.servesAll("t_full", qontoServer) || s.servesAll("t_other", qontoServer) || s.servesAll("t_full", "https://mcp.stripe.com") {
		t.Fatal("the whole live list for the wrong accounts or hosts")
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

// asQonto makes the grant id a sign-in to Qonto's server that asked for the
// reviewed list, as a direct Qonto entry does, while the fake vendor answers it.
func (e *env) asQonto(id string) {
	e.t.Helper()
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	g.Request.Server, g.Vendor.Scopes = qontoServer, slices.Clone(qontoScopes)
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

// The stored token is checked at refresh too, not only the refresh answer:
// a token granted more at sign-in (before the check) is caught at its next
// refresh even when the refresh answer names no scope.
func TestABankRefreshChecksTheStoredToken(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.asQonto(g.ID)
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	ctx := context.Background()
	e.vendor.Expire()
	tok, it, _ := e.s.openToken(ctx, g)
	tok.Expiry, tok.Scope = time.Now().Add(-time.Second), "organization.read card.write"
	_ = e.s.storeToken(ctx, g, tok, &it.Ver)
	e.vendor.Scope = "" // the refresh answer names none
	if _, _, err := e.s.accessToken(ctx, g, 0); err == nil {
		t.Fatal("a stored token holding more was refreshed and used")
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if now.Status != StatusNeedsReconnect || now.SignedOutCause != CauseScopeWider || !strings.Contains(now.Reason, "card.write") {
		t.Fatalf("grant: %s %q %q", now.Status, now.SignedOutCause, now.Reason)
	}
}

// The bank's own approvals are never served, on any account, the whole
// live list included: tafs prepares a transfer or card request, and Qonto's
// own approval finishes it, whatever a Card opens.
func TestQontoBankApprovalsAreNeverServed(t *testing.T) {
	s := &Server{FullScopes: map[string]bool{"t_full": true}, LiveListTest: true, ManagedTestTenants: map[string]bool{"t_test": true}}
	listed := []mcpTool{{Name: "list_transactions"}, {Name: "approve_request"}, {Name: "decline_request"}, {Name: "create_multi_transfer_request"}}
	for _, tenant := range []string{"t_full", "t_test", "t_other"} {
		for _, tool := range qontoBankApproval {
			if s.served(tenant, qontoServer, tool) {
				t.Fatalf("%s: %s served", tenant, tool)
			}
		}
		for _, got := range s.servedFor(tenant, qontoServer, listed) {
			if slices.Contains(qontoBankApproval, got.Name) {
				t.Fatalf("%s: %s listed", tenant, got.Name)
			}
		}
	}
	// The whole list keeps preparing a request, for Qonto's approval to finish.
	if got := s.servedFor("t_full", qontoServer, listed); len(got) != 2 || got[1].Name != "create_multi_transfer_request" {
		t.Fatalf("the whole list: %v", got)
	}
	// Off Qonto's host, a tool of that name is another vendor's, served as any.
	if !s.served("t_other", "https://tools.example/mcp", "approve_request") {
		t.Fatal("another vendor's approve_request refused")
	}
	// Each name is one Qonto has: a typo would leave the real one served.
	for _, name := range qontoBankApproval {
		if !slices.Contains(qontoReference, name) {
			t.Fatalf("%s isn't in Qonto's reference", name)
		}
	}
}
