package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

func TestBrowserOfferFrom(t *testing.T) {
	for _, c := range []struct {
		offer, tenants string
		want           BrowserOffer
	}{
		{"", "", BrowserOffer{}},
		{"off", "tenants:tnt_a", BrowserOffer{}},
		{"on", "", BrowserOffer{All: true}},
		{"test", "tenants:tnt_a, tnt_b", BrowserOffer{Tenants: map[string]bool{"tnt_a": true, "tnt_b": true}}},
	} {
		got, err := BrowserOfferFrom(c.offer, c.tenants)
		if err != nil || got.All != c.want.All || len(got.Tenants) != len(c.want.Tenants) {
			t.Fatalf("%q %q: %+v %v", c.offer, c.tenants, got, err)
		}
		for id := range c.want.Tenants {
			if !got.Tenants[id] {
				t.Fatalf("%q %q: %s missing from %+v", c.offer, c.tenants, id, got)
			}
		}
	}
	for _, c := range [][2]string{{"test", ""}, {"test", "tnt_a"}, {"test", "tenants:not a tenant"}, {"yes", ""}, {"ON", ""}} {
		if _, err := BrowserOfferFrom(c[0], c[1]); err == nil {
			t.Fatalf("%q %q was taken", c[0], c[1])
		}
	}
}

// A configured worker offers the hosted browser to no account by itself:
// only to the accounts PLANQUE_BROWSER_OFFER names. An account off the list
// sees no "Use websites", its ask is refused, and a grant it already has
// opens nothing: no call, no sign-in view.
func TestTheHostedBrowserIsOfferedOnlyToTheAccountsNamed(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect() // while this account is on the list
	listed := func() bool {
		st, cat := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("catalog: %d %v", st, cat)
		}
		return strings.Contains(mustJSONString(cat), `"entry":"browser"`)
	}
	if !listed() {
		t.Fatal("the account on the list doesn't see Use websites")
	}
	for _, c := range []struct {
		name  string
		offer BrowserOffer
		view  string // each sign-in ask its own signed body
	}{
		{"unset", BrowserOffer{}, fakes.ShopLogin},
		{"another account only", BrowserOffer{Tenants: map[string]bool{"tnt_other": true}}, fakes.ShopItem},
	} {
		name := c.name
		e.s.BrowserOffer = c.offer
		if listed() {
			t.Errorf("%s: the catalog lists Use websites", name)
		}
		sites := map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}}
		if st, out := e.ask(sites); st != 403 || out["code"] != CodeBrowserOff {
			t.Errorf("%s: the ask: %d %v", name, st, out)
		}
		if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/browser/grants", map[string]any{"agent": "agt_taf", "browser": map[string]any{"sites": []string{"shop.example"}}}, nil, "k_setup"); st != 403 || out["code"] != CodeBrowserOff {
			t.Errorf("%s: the catalog's ask: %d %v", name, st, out)
		}
		sessions := e.b.Sessions
		for _, tool := range []string{"browse_open", "browse_act"} {
			args := map[string]any{"url": fakes.ShopItem}
			if tool == "browse_act" {
				args["read"], args["steps"] = e.b.ReadOf(fakes.ShopItem), []map[string]any{{"click": "add"}}
			}
			if out, sc := e.call(g, tool, args); code(sc) != CodeBrowserOff || sc["sent"] != false {
				t.Errorf("%s: %s on a grant from before: %s %v", name, tool, out, sc)
			}
		}
		if st, out := e.startView(g.ID, map[string]any{"url": c.view}); st != 409 || out["code"] != CodeBrowserOff {
			t.Errorf("%s: the sign-in view: %d %v", name, st, out)
		}
		if e.b.Sessions != sessions {
			t.Errorf("%s: a browser was started: %d sessions, before %d", name, e.b.Sessions, sessions)
		}
	}
	// "on": every account, this one included.
	e.s.BrowserOffer = BrowserOffer{All: true}
	if !listed() {
		t.Error("on: the catalog doesn't list Use websites")
	}
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Errorf("on: a call: %v", sc)
	}
}

// The account's connections say whether the hosted browser is offered to it,
// and whether it pays, so Coach drafts no payment rules for a browser that
// can't pay. Nothing about it for an account it isn't offered to.
func TestTheConnectionsSayWhetherTheHostedBrowserPays(t *testing.T) {
	e := newBrowserEnv(t)
	read := func() map[string]any {
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/connections", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("connections: %d %v", st, out)
		}
		return out
	}
	e.s.BrowserPay = false
	if b, _ := read()["browser"].(map[string]any); b == nil || b["pay"] != false {
		t.Errorf("paying off: %v", read()["browser"])
	}
	e.s.BrowserPay = true
	if b, _ := read()["browser"].(map[string]any); b == nil || b["pay"] != true {
		t.Errorf("paying on: %v", read()["browser"])
	}
	e.s.BrowserOffer = BrowserOffer{Tenants: map[string]bool{"tnt_other": true}}
	if b, ok := read()["browser"]; ok {
		t.Errorf("not offered to this account: %v", b)
	}
	e.s.BrowserOffer, e.s.Browser = BrowserOffer{All: true}, nil
	if b, ok := read()["browser"]; ok {
		t.Errorf("no worker: %v", b)
	}
}
