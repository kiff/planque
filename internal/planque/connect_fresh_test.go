package planque

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// siblingOf saves a copy of g as another taf's connection on the same
// provider account, as a held sign-in's second connection is.
func siblingOf(t *testing.T, e *env, g *Grant, id, agent string) *Grant {
	t.Helper()
	raw, err := json.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	var c Grant
	if err := json.Unmarshal(raw, &c); err != nil {
		t.Fatal(err)
	}
	c.ID, c.Request.Grant, c.Request.AgentName, c.ver = id, id, agent, 0
	if err := e.s.saveGrant(context.Background(), &c); err != nil {
		t.Fatal(err)
	}
	return &c
}

func connectFresh(t *testing.T, e *env, f *connectorFixture, g *Grant) (int, map[string]any) {
	t.Helper()
	e.at(time.Minute) // a new signature, later, on both clocks
	e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/connect-fresh", reapprove(g), e.phone, "")
}

func idsOf(v any) []string {
	var out []string
	for _, x := range v.([]any) {
		out = append(out, x.(string))
	}
	return out
}

// One fresh sign-in brings back every signed-out connection that was on the
// dead provider account, each with its own tools and Card as they were; a
// connection on another account is left alone; the same confirmation twice
// moves nothing again.
// Category: authority
// Attack: a fresh sign-in that reaches a connection on another provider account, or is replayed
// Refused: only the signed-out connections on the old account move; a replayed confirm is 409
func TestAFreshSignInBringsBackEveryConnectionOnTheOldAccount(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct, f.noName = true, true
	first := connectConnector(t, e, f)
	if first.Connector.AccountName != "" {
		t.Fatalf("named %q", first.Connector.AccountName)
	}
	declared, _ := e.kiff.Connection(tenant, "reply_ticket")
	second := siblingOf(t, e, first, "grt_second", "Win-back offers")
	other := siblingOf(t, e, first, "grt_other", "Another account's taf")
	other.Connector.Account, other.Request.Connector.Account = "another_account", "another_account"
	if err := e.s.saveGrant(ctx, other); err != nil {
		t.Fatal(err)
	}
	// A signed-out connection on the same account of another toolkit is left
	// alone too: only the same app comes back with the sign-in.
	otherApp := siblingOf(t, e, first, "grt_otherapp", "Another app's taf")
	otherApp.Request.Connector.Toolkit = "another_toolkit"
	if err := e.s.saveGrant(ctx, otherApp); err != nil {
		t.Fatal(err)
	}
	signOut401(t, e, f, first)
	// The other account's connection is signed out too, by its own account.
	other, _ = e.s.loadGrant(ctx, tenant, other.ID)
	other.Status = StatusNeedsReconnect
	if err := e.s.saveGrant(ctx, other); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{first.ID, second.ID} {
		if g, _ := e.s.loadGrant(ctx, tenant, id); g.Status != StatusNeedsReconnect {
			t.Fatalf("%s after the 401: %s", id, g.Status)
		}
	}
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	_, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, nil, "k_setup")
	if view["reconnect"] != "fresh" {
		t.Fatalf("a connection with no name kept reads reconnect %v", view["reconnect"])
	}
	// Sign in again, the old way, can't work for it.
	f.noName = false
	path := "/v1/accounts/" + tenant + "/grants/" + first.ID
	if st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, ""); st != 409 || out["sign_in_url"] != nil {
		t.Fatalf("sign in again with no name kept: %d %v", st, out)
	}
	st, out := connectFresh(t, e, f, g)
	if st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("connect fresh: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&status=ACTIVE")
	st, out = e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if got := idsOf(out["grants"]); !slices.Equal(got, []string{first.ID, second.ID}) {
		t.Fatalf("the connections brought back: %v", got)
	}
	for _, id := range []string{first.ID, second.ID} {
		g, _ := e.s.loadGrant(ctx, tenant, id)
		if g.Status != StatusConnected || g.Connector.Account == first.Connector.Account || g.Request.Connector.Account != g.Connector.Account ||
			len(g.Tools) != len(first.Tools) || g.Connector.AccountName == "" || g.Reconnect || g.FreshFrom != "" {
			t.Fatalf("%s after the fresh sign-in: %s, account %q (was %q), bound %q, tools %d, name %q", id, g.Status, g.Connector.Account, first.Connector.Account, g.Request.Connector.Account, len(g.Tools), g.Connector.AccountName)
		}
	}
	a, _ := e.s.loadGrant(ctx, tenant, first.ID)
	b, _ := e.s.loadGrant(ctx, tenant, second.ID)
	if a.Connector.Account != b.Connector.Account {
		t.Fatalf("the two connections are on different sign-ins: %q, %q", a.Connector.Account, b.Connector.Account)
	}
	if c, ok := e.kiff.Connection(tenant, "reply_ticket"); !ok || c.Credential != declared.Credential {
		t.Fatal("its registration in KIFF changed: the taf's Card would no longer cover it")
	}
	if o, _ := e.s.loadGrant(ctx, tenant, otherApp.ID); o.Status != StatusNeedsReconnect || o.Connector.Account != first.Connector.Account {
		t.Fatalf("a connection of another toolkit moved: %s, %q", o.Status, o.Connector.Account)
	}
	if o, _ := e.s.loadGrant(ctx, tenant, other.ID); o.Status != StatusNeedsReconnect || o.Connector.Account != "another_account" {
		t.Fatalf("a connection on another account moved: %s, %q", o.Status, o.Connector.Account)
	}
	// It works: a call goes through on the new sign-in, for the second too.
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_back"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: b.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_back", Args: args})
	if got := mustJSONString(e.s.call(ctx, b, "reply_ticket", args, permit)); strings.Contains(got, `"isError":true`) {
		t.Fatalf("a call after the fresh sign-in: %s", got)
	}
	// The same confirmation again moves nothing.
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 409 {
		t.Fatalf("a replayed confirm: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); g.Connector.Account != a.Connector.Account || g.Status != StatusConnected {
		t.Fatalf("after a replay: %s, %q", g.Status, g.Connector.Account)
	}
}

// A connection that has its account's name signs in again on that account,
// never fresh: a sign-in to another account is never swapped under the Card
// the owner approved for the first.
// Category: authority
// Attack: connect-fresh on a connection whose account is known, to move its tools to another account
// Refused: 409, nothing opened, the connection still signed out, and its row says "same"
func TestAConnectionWithItsAccountsNameIsNeverSignedInFresh(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct = true
	first := connectConnector(t, e, f)
	if first.Connector.AccountName == "" {
		t.Fatal("no name kept")
	}
	signOut401(t, e, f, first)
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	if st, out := connectFresh(t, e, f, g); st != 409 || out["sign_in_url"] != nil || !strings.Contains(mustJSONString(out), "account is known") {
		t.Fatalf("connect fresh with the name kept: %d %v", st, out)
	}
	if g, _ = e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusNeedsReconnect || g.Reconnect || g.FreshFrom != "" {
		t.Fatalf("after the refusal: %s, reconnect %v, fresh from %q", g.Status, g.Reconnect, g.FreshFrom)
	}
	if _, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, nil, "k_setup"); view["reconnect"] != "same" {
		t.Fatalf("a connection with its name kept reads reconnect %v", view["reconnect"])
	}
	// And a connected one has none.
	e2, f2 := newConnectorFixture(t)
	c := connectConnector(t, e2, f2)
	if _, view := e2.do("GET", "/v1/accounts/"+tenant+"/grants/"+c.ID, nil, nil, "k_setup"); view["reconnect"] != nil {
		t.Fatalf("a connected connection reads reconnect %v", view["reconnect"])
	}
}
