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
	e.tafName = "Refund requests"
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
	// The owner can see where the tafs went: each re-pointed connection reads
	// the new account's name, and the taf that uses it.
	for id, taf := range map[string]string{first.ID: "Refund requests", second.ID: "Win-back offers"} {
		_, v := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
		if v["account_name"] != "Demo Shop" || !slices.Contains(idsOf(v["used_by"]), taf) {
			t.Fatalf("%s after the fresh sign-in reads account_name %v, used_by %v", id, v["account_name"], v["used_by"])
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

// Only the siblings that can come back do: signed out, on the old account, with
// no name kept, of the same app, and with no sign-in of their own under way.
// Each of these four stays as it was, and the one that fits moves.
// Category: authority
// Attack: a fresh sign-in that reaches a connection with its account known, one with its own sign-in under way, one connected, or one of another app
// Refused: each stays as it was; only the signed-out unnamed sibling of the same app comes back
func TestAFreshSignInLeavesEverySiblingThatCantComeBack(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct, f.noName = true, true
	first := connectConnector(t, e, f)
	fits := siblingOf(t, e, first, "grt_fits", "Win-back offers")
	named := siblingOf(t, e, first, "grt_named", "Named sibling")
	underway := siblingOf(t, e, first, "grt_underway", "Sibling signing in")
	connected := siblingOf(t, e, first, "grt_connected", "Connected sibling")
	otherProvider := siblingOf(t, e, first, "grt_otherprov", "Another provider")
	signOut401(t, e, f, first)
	change := func(id string, edit func(g *Grant)) {
		g, err := e.s.loadGrant(ctx, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
		edit(g)
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	change(named.ID, func(g *Grant) { g.Connector.AccountName = "Named Co" })
	change(underway.ID, func(g *Grant) { g.Reconnect = true })
	change(connected.ID, func(g *Grant) { g.Status = StatusConnected })
	change(otherProvider.ID, func(g *Grant) { g.Request.Connector.Provider = "another_provider" })
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	path := "/v1/accounts/" + tenant + "/grants/" + first.ID
	f.noName = false
	if st, out := connectFresh(t, e, f, g); st != 200 {
		t.Fatalf("connect fresh: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&status=ACTIVE")
	st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, "")
	if st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if got := idsOf(out["grants"]); !slices.Equal(got, []string{first.ID, fits.ID}) {
		t.Fatalf("the connections brought back: %v", got)
	}
	for _, c := range []struct {
		id   string
		want string
	}{{named.ID, StatusNeedsReconnect}, {underway.ID, StatusNeedsReconnect}, {connected.ID, StatusConnected}, {otherProvider.ID, StatusNeedsReconnect}} {
		o, _ := e.s.loadGrant(ctx, tenant, c.id)
		if o.Status != c.want || o.Connector.Account != first.Connector.Account || o.Request.Connector.Account != first.Request.Connector.Account {
			t.Fatalf("%s moved: %s (want %s), account %q (was %q)", c.id, o.Status, c.want, o.Connector.Account, first.Connector.Account)
		}
	}
	if o, _ := e.s.loadGrant(ctx, tenant, underway.ID); !o.Reconnect {
		t.Fatal("a sibling's own sign-in under way was ended")
	}
}

// A fresh sign-in that is never confirmed ends with the connection signed out
// and no memory of the old account: a later sign-in on it is held to an account
// like any other, never skipping the comparison.
// Category: authority
// Attack: a stale fresh marker left on a connection after its fresh sign-in expired, so a later sign-in skips the account comparison
// Refused: the expired sign-in clears the marker, and the connection reads signed out
func TestAnExpiredFreshSignInLeavesNoFreshMarker(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct, f.noName = true, true
	first := connectConnector(t, e, f)
	signOut401(t, e, f, first)
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	f.noName = false
	if st, out := connectFresh(t, e, f, g); st != 200 {
		t.Fatalf("connect fresh: %d %v", st, out)
	}
	e.signIn(f.callback + "&status=ACTIVE") // the sign-in comes back, and is never confirmed
	g, _ = e.s.loadGrant(ctx, tenant, first.ID)
	if g.Status != StatusSignedIn || g.FreshFrom == "" {
		t.Fatalf("before expiry: %s, fresh from %q", g.Status, g.FreshFrom)
	}
	e.at(time.Hour)
	if _, err := e.s.expire(ctx, g); err != nil {
		t.Fatal(err)
	}
	g, _ = e.s.loadGrant(ctx, tenant, first.ID)
	if g.Status != StatusNeedsReconnect || g.Reconnect || g.FreshFrom != "" {
		t.Fatalf("after expiry: %s, reconnect %v, fresh from %q", g.Status, g.Reconnect, g.FreshFrom)
	}
}
