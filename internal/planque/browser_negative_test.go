package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
)

// connectAs is connect for another taf's agent: a second browser grant on
// the same account, limited to the same site.
func (e *browserEnv) connectAs(agent string) *Grant {
	e.t.Helper()
	st, out := e.ask(map[string]any{"agent": agent, "browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}})
	if st != 201 {
		e.t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	hash, _ := profile.Digest(g.Request.Profile)
	approved := ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}
	if st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": approved}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	g, err = e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	return g
}

// named is the name g's tool goes by in KIFF: a second grant's names are
// suffixed so they don't collide with the first's.
func named(g *Grant, tool string) string {
	for _, t := range g.Tools {
		if t.Tool == tool {
			return t.Name
		}
	}
	return tool
}

// Two browser grants on one account, on the same site, share nothing: one
// grant's sign-in never starts the other's session, and a checkout one
// grant read can't be paid through the other.
func TestBrowserGrantsShareNoSignInAndNoCheckout(t *testing.T) {
	e := newBrowserEnv(t)
	a := e.connect()
	b := e.connectAs("agt_other_taf")
	if a.ID == b.ID {
		t.Fatal("one grant for two tafs")
	}
	if _, sc := e.call(a, "browse_act", map[string]any{"url": fakes.ShopLogin, "steps": []map[string]any{{"fill": "user", "value": "ana@example.com"}, {"click": "signin"}}}); code(sc) != "" {
		t.Fatalf("sign in on a: %v", sc)
	}
	if _, sc := e.call(b, named(b, "browse_open"), map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("open on b: %v", sc)
	}
	if last := e.b.States[len(e.b.States)-1]; strings.Contains(string(last), "signed-in") {
		t.Fatalf("b's session started with a's sign-in: %s", last)
	}
	if _, sc := e.call(a, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
		t.Fatalf("checkout on a: %v", sc)
	}
	sessions := e.b.Sessions
	if out, sc := e.call(b, named(b, "browse_pay"), map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeTotalUnread || len(e.b.Orders) != 0 || e.b.Sessions != sessions {
		t.Fatalf("b paid a's checkout: %s (orders %v)", out, e.b.Orders)
	}
}

// After disconnecting, the relay answers nothing for the grant, even to a
// call with a permit KIFF minted before, and no browser starts.
func TestBrowserDisconnectedGrantStartsNothing(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	c, ok := e.kiff.Connection(tenant, "browse_open")
	if !ok {
		t.Fatal("browse_open isn't connected in KIFF")
	}
	// The same path before: the gateway's call reaches the browser.
	before := mustJSON(map[string]any{"url": fakes.ShopItem, IdempotencyArg: "op_before"})
	if res, err := e.kiff.Forward(context.Background(), c, before, e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_open", Sub: g.Request.Agent, Op: "op_before", Args: before})); err != nil || e.b.Sessions != 1 {
		t.Fatalf("a call before disconnect: %s %v, sessions %d", res, err, e.b.Sessions)
	}
	args := mustJSON(map[string]any{"url": fakes.ShopItem, IdempotencyArg: "op_after"})
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_open", Sub: g.Request.Agent, Op: "op_after", Args: args})
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	res, err := e.kiff.Forward(context.Background(), c, args, permit)
	if err == nil && !strings.Contains(string(res), "Refused") {
		t.Fatalf("a call after disconnect: %s", res)
	}
	if e.b.Sessions != 1 {
		t.Fatalf("a browser started after disconnect: %d sessions, want the 1 from before", e.b.Sessions)
	}
}

// A grant with no list of sites still never opens an address: the forms
// a browser reads as an IPv4 address (a short or hexadecimal or octal
// one, or one written with another script's full stop) are refused like
// a dotted one, before a browser starts.
func TestBrowserRefusesEveryFormOfAnAddress(t *testing.T) {
	e := newBrowserEnv(t)
	st, out := e.ask(map[string]any{"agent": "agt_any", "browser": map[string]any{"currency": "EUR"}})
	if st != 201 {
		t.Skipf("a grant without sites isn't asked for here: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, out["id"].(string))
	hash, _ := profile.Digest(g.Request.Profile)
	if st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	for _, u := range []string{"https://127.1/", "https://0x7f.1/", "https://0177.0.0.1/", "https://1.2.3.0x4/", "https://127.0.0。1/", "https://169.254.169．254/"} {
		out, sc := e.call(g, "browse_open", map[string]any{"url": u})
		if code(sc) != browser.CodeAddressPrivate || sc["sent"] != false {
			t.Errorf("%s: %s", u, out)
		}
	}
	if e.b.Sessions != 0 {
		t.Fatalf("an address started a browser: %d sessions", e.b.Sessions)
	}
}
