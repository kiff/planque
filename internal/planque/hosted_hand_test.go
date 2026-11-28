package planque

import (
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// newHandEnv is the browser env with the hosted hand-over on, and a grant
// on shop.example that has its two tools.
func newHandEnv(t *testing.T) (*browserEnv, *Grant) {
	t.Helper()
	e := newBrowserEnv(t)
	e.s.HostedHand = true
	g := e.connect()
	if !hasTool(g, browseHand) || !hasTool(g, browseHandResult) {
		t.Fatalf("the hosted grant's tools: %+v", g.Tools)
	}
	return e, g
}

func hasTool(g *Grant, tool string) bool {
	for _, t := range g.Tools {
		if t.Tool == tool {
			return true
		}
	}
	return false
}

// handView is the phone's signed ask for a hand-over's view.
func (e *browserEnv) handView(g *Grant, hand, what string) (int, map[string]any) {
	e.t.Helper()
	// Each ask is signed anew (a signed request is never replayed): the
	// body differs by a nonce the routes don't read.
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/hands/"+hand+"/"+what, map[string]any{"nonce": randomSuffix()}, e.phone, "")
}

// 1. "Your turn" in the hosted browser: the taf hands a step over, the
// session stays open for the phone's signed view, the taf waits, the person
// does it in the view, and Done seals what the site kept, which the taf's
// next call starts with; the answer is read once.
func TestHostedHandOverThroughTheLiveView(t *testing.T) {
	e, g := newHandEnv(t)
	out, sc := e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here, then tap Done"})
	hand, _ := sc["hand"].(string)
	if code(sc) != "" || sc["state"] != "handed" || !strings.HasPrefix(hand, "hh_") || sc["site"] != "shop.example" || e.b.Sessions != 1 {
		t.Fatalf("hand: %s", out)
	}
	// The taf waits: its calls are refused while the step is the person's.
	if out, sc := e.call(g, browseOpen, map[string]any{"url": fakes.ShopItem}); code(sc) != CodeSigningIn || !strings.Contains(out, hand) || e.b.Sessions != 1 {
		t.Fatalf("a call while handed: %s", out)
	}
	if _, sc := e.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handWaiting {
		t.Fatalf("result while waiting: %v", sc)
	}
	// Unsigned, or another hand, is refused.
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/hands/"+hand+"/view", nil, nil, "k_setup"); st != 401 {
		t.Fatalf("unsigned view: %d", st)
	}
	if st, _ := e.handView(g, "hh_other", "view"); st != 410 {
		t.Fatalf("another hand: %d", st)
	}
	st, v := e.handView(g, hand, "view")
	viewURL, _ := v["view"].(string)
	if st != 200 || !strings.HasPrefix(viewURL, e.url+"/browser/view/") || v["site"] != "shop.example" || v["words"] != "Sign in here, then tap Done" || v["expires"] == "" {
		t.Fatalf("view: %d %v", st, v)
	}
	// A second ask makes the first URL useless.
	st2, v2 := e.handView(g, hand, "view")
	if st2 != 200 {
		t.Fatalf("second view: %d %v", st2, v2)
	}
	if st, _, _ := e.view("POST", viewURL+"/frame", nil); st != 404 {
		t.Fatalf("the first URL after a second ask: %d", st)
	}
	viewURL = v2["view"].(string)
	// The page names the site, planque's, and the taf's words in quotes.
	st, _, page := e.view("GET", viewURL, nil)
	if st != 200 || !strings.Contains(string(page), "<title>shop.example</title>") || !strings.Contains(string(page), "“Sign in here, then tap Done”") {
		t.Fatalf("page: %d %s", st, page)
	}
	for _, in := range []map[string]any{{"kind": "type", "text": "hunter2"}, {"kind": "key", "key": "Enter"}} {
		if st, _, raw := e.view("POST", viewURL+"/input", in); st != 200 {
			t.Fatalf("input %v: %d %s", in, st, raw)
		}
	}
	if st, _, raw := e.view("POST", viewURL+"/done", nil); st != 200 {
		t.Fatalf("done: %d %s", st, raw)
	}
	if _, sc := e.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handDone {
		t.Fatalf("result: %v", sc)
	}
	if out, sc := e.call(g, browseHandResult, map[string]any{"hand": hand}); code(sc) != CodeBadArguments {
		t.Fatalf("read twice: %s", out)
	}
	// The next call starts with what the site kept.
	if out, sc := e.call(g, browseOpen, map[string]any{"url": fakes.ShopItem}); code(sc) != "" || !strings.Contains(string(e.b.States[len(e.b.States)-1]), "signed-in") {
		t.Fatalf("the next call: %s, states %q", out, e.b.States)
	}
}

// 2. Never on a page that pays: at the hand-over, nothing is kept open; in
// the view, a navigation to a checkout ends it before a typed digit passes,
// seals nothing, and the taf reads it as put off.
func TestHostedHandOverNeverOnAPaymentPage(t *testing.T) {
	e, g := newHandEnv(t)
	if out, sc := e.call(g, browseHand, map[string]any{"url": fakes.ShopCheckout, "words": "Type your card here"}); code(sc) != "payment_step" || sc["hand"] != nil {
		t.Fatalf("hand on a checkout: %s", out)
	}
	if out, sc := e.call(g, browseOpen, map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("nothing held after the refusal: %s", out)
	}

	_, sc := e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here"})
	hand := sc["hand"].(string)
	_, v := e.handView(g, hand, "view")
	viewURL := v["view"].(string)
	_, verBefore, _ := e.s.openBrowserState(t.Context(), g)
	// The person goes to the checkout on the grant's own site.
	if st, _, raw := e.view("POST", viewURL+"/input", map[string]any{"kind": "go", "url": fakes.ShopCheckout}); st != 200 {
		t.Fatalf("go: %d %s", st, raw)
	}
	inputsBefore := len(e.b.Inputs)
	st, _, raw := e.view("POST", viewURL+"/input", map[string]any{"kind": "type", "text": "4242 4242 4242 4242"})
	if st != 404 || !strings.Contains(string(raw), "takes a payment") {
		t.Fatalf("typing on the checkout: %d %s", st, raw)
	}
	for _, in := range e.b.Inputs[inputsBefore:] {
		if strings.ContainsAny(in.Text, "0123456789") {
			t.Fatalf("a digit reached the checkout: %+v", in)
		}
	}
	if st, _, _ := e.view("POST", viewURL+"/frame", nil); st != 404 {
		t.Fatalf("the view after it ended: %d", st)
	}
	if _, sc := e.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handDeclined {
		t.Fatalf("result: %v", sc)
	}
	if _, verAfter, _ := e.s.openBrowserState(t.Context(), g); verAfter != verBefore {
		t.Fatalf("something was sealed: version %d, was %d", verAfter, verBefore)
	}
}

// 3. The app's Not now ends the session and keeps nothing; a view left to
// lapse reads as expired.
func TestHostedHandOverNotNowAndLapse(t *testing.T) {
	e, g := newHandEnv(t)
	_, sc := e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here"})
	hand := sc["hand"].(string)
	if st, out := e.handView(g, hand, "decline"); st != 200 {
		t.Fatalf("not now: %d %v", st, out)
	}
	if _, sc := e.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handDeclined {
		t.Fatalf("result: %v", sc)
	}
	if st, _ := e.handView(g, hand, "view"); st != 410 {
		t.Fatalf("a view after Not now: %d", st)
	}

	_, sc = e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here"})
	hand = sc["hand"].(string)
	later := time.Now().Add(signInViewTTL + time.Minute)
	e.s.Now = func() time.Time { return later }
	h, err := e.s.readHostedHand(t.Context(), g, hand)
	e.s.Now = nil
	if err != nil || h.Outcome != handExpired {
		t.Fatalf("lapsed: %+v %v", h, err)
	}
}

// 4. Off, a hosted grant has no hand-over tools, and asking for one is
// refused.
func TestHostedHandOverIsOffByDefault(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if hasTool(g, browseHand) || hasTool(g, browseHandResult) {
		t.Fatalf("tools with the hand-over off: %+v", g.Tools)
	}
	st, out := e.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}}, "tools": []any{map[string]any{"tool": browseHand}}})
	if st != 400 || !strings.Contains(out["error"].(string), "isn't offered in the hosted browser") {
		t.Fatalf("asking for it: %d %v", st, out)
	}
}

// tafs cloud lists the hand-overs still waiting, with the account's key:
// the grant, its agent, the hand, the site, the words and the times, no
// URL or token; done, put off and lapsed ones aren't listed.
func TestWaitingHandOversAreListed(t *testing.T) {
	e, g := newHandEnv(t)
	list := func() []any {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/hands", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("list: %d %v", st, out)
		}
		hands, _ := out["hands"].([]any)
		return hands
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/hands", nil, nil, ""); st == 200 {
		t.Fatal("listed with no key")
	}
	if got := list(); len(got) != 0 {
		t.Fatalf("none yet: %v", got)
	}
	_, sc := e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here"})
	hand := sc["hand"].(string)
	got := list()
	if len(got) != 1 {
		t.Fatalf("one waiting: %v", got)
	}
	h := got[0].(map[string]any)
	if h["grant"] != g.ID || h["agent"] != g.Request.Agent || h["hand"] != hand || h["site"] != "shop.example" || h["words"] != "Sign in here" || h["since"] == "" || h["expires"] == "" || h["view"] != nil {
		t.Fatalf("the hand: %v", h)
	}
	if st, _ := e.handView(g, hand, "decline"); st != 200 {
		t.Fatal("decline")
	}
	if got := list(); len(got) != 0 {
		t.Fatalf("after Not now: %v", got)
	}
	e.call(g, browseHand, map[string]any{"url": fakes.ShopLogin, "words": "Sign in here"})
	later := time.Now().Add(signInViewTTL + time.Minute)
	e.s.Now = func() time.Time { return later }
	defer func() { e.s.Now = nil }()
	if got := list(); len(got) != 0 {
		t.Fatalf("lapsed: %v", got)
	}
}
