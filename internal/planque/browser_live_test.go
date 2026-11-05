package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/browserworker"
	"github.com/kiff/planque/internal/fakes"
)

// The browser end to end, through planque's guard and the real browser
// (kiff/planque#86, step 3): the test shop is served over TLS on loopback
// and reached as https://shop.example through the worker's proxy, which
// the test maps there. Skipped where no Chromium is found.
func newLiveBrowserEnv(t *testing.T) (*browserEnv, *fakes.Shop) {
	t.Helper()
	if _, err := browserworker.ChromePath(""); err != nil {
		t.Skip(err)
	}
	shop := fakes.NewTLSShop()
	t.Cleanup(shop.Close)
	e := newBrowserEnv(t)
	e.s.Browser = &browserworker.Local{Options: browserworker.Options{AllowLocal: true, NavigateWait: 20 * time.Second,
		Map: map[string]string{"shop.example:443": strings.TrimPrefix(shop.URL, "https://")},
		Resolve: func(_ context.Context, host string) ([]netip.Addr, error) {
			return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
		}}}
	return e, shop
}

const (
	liveItem     = "https://shop.example/item/pen"
	liveCheckout = "https://shop.example/checkout"
	liveLogin    = "https://shop.example/login"
	liveOther    = "https://other.example/"
)

// Look up an item, put it in the cart, read the checkout, pay once with
// the permit; a replay answers as the first did; a changed total after
// the read sends nothing; a card number never reaches the agent.
func TestBrowserLiveLookUpApproveAndPayOnce(t *testing.T) {
	e, shop := newLiveBrowserEnv(t)
	g := e.connect()

	out, sc := e.call(g, "browse_open", map[string]any{"url": liveItem})
	if code(sc) != "" || sc["site"] != "shop.example" {
		t.Fatalf("open: %s", out)
	}
	page := sc["page"].(map[string]any)
	qty, add := refByLabel(page, "fields", "Quantity"), refByLabel(page, "controls", "Add to cart")
	out, sc = e.call(g, "browse_act", map[string]any{"url": liveItem, "steps": []map[string]any{{"fill": qty, "value": "2"}, {"click": add}}})
	if code(sc) != "" || sc["steps_run"] != float64(2) || !strings.Contains(out, "Your cart") {
		t.Fatalf("add to cart: %s", out)
	}

	// The checkout: its card fields masked, its total read, its one line.
	out, sc = e.call(g, "browse_open", map[string]any{"url": liveCheckout})
	if code(sc) != "" || strings.Contains(out, "4242") || strings.Contains(out, `"value": "123"`) {
		t.Fatalf("checkout page: %s", out)
	}
	page = sc["page"].(map[string]any)
	place := refByLabel(page, "controls", "Place order")
	// A click on the checkout through browse_act is refused before it runs.
	if out, sc := e.call(g, "browse_act", map[string]any{"url": liveCheckout, "steps": []map[string]any{{"click": place}}}); code(sc) != browser.CodePaymentStep || len(shop.Orders) != 0 {
		t.Fatalf("placing the order through browse_act: %s (orders %v)", out, shop.Orders)
	}
	out, sc = e.call(g, "browse_checkout", map[string]any{"url": liveCheckout})
	if code(sc) != "" || sc["total"] != "2.50" || sc["currency"] != "EUR" {
		t.Fatalf("checkout read: %s", out)
	}
	var lines []browser.Line
	raw, _ := json.Marshal(sc["lines"])
	if err := json.Unmarshal(raw, &lines); err != nil || len(lines) != 1 || lines[0] != (browser.Line{Item: "Blue pen", Quantity: 2, Price: "2.50"}) {
		t.Errorf("lines: %s", raw)
	}

	// The shop changes its total before the owner approves: nothing is paid.
	shop.Total = "3.00"
	if out, sc := e.call(g, "browse_pay", map[string]any{"url": liveCheckout, "ref": place}); code(sc) != browser.CodeCheckoutChanged || len(shop.Orders) != 0 {
		t.Fatalf("pay after a change: %s (orders %v)", out, shop.Orders)
	}
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": liveCheckout}); sc["total"] != "3.00" {
		t.Fatalf("second read: %v", sc)
	}
	// Paid once with the permit KIFF issued; the replay sends nothing.
	op := "op_live_pay"
	args := mustJSON(map[string]any{"url": liveCheckout, "ref": place, IdempotencyArg: op})
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args})
	first, fsc := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
	if code(fsc) != "" || len(shop.Orders) != 1 || !strings.Contains(first, `"total": "3.00"`) || !strings.Contains(first, "Order 1001 placed") {
		t.Fatalf("pay: %s (orders %v)", first, shop.Orders)
	}
	if again, _ := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit)); again != first || len(shop.Orders) != 1 {
		t.Fatalf("replay: %s (orders %v)", again, shop.Orders)
	}
	if refs := e.ev.Refs(tenant, g.ID); len(refs) != 2 {
		t.Errorf("screenshots kept: %v", refs)
	}
	for _, ref := range e.ev.Refs(tenant, g.ID) {
		if png := e.ev.Kept[ref]; !strings.HasPrefix(string(png), "\x89PNG") {
			t.Errorf("%s isn't a PNG", ref)
		}
	}
}

// Secrets and sites hold inside the real browser: a password field's
// value never leaves the worker, typing into it is refused, and a site
// off the list is refused by planque before any browser starts.
func TestBrowserLiveKeepsSecretsAndSites(t *testing.T) {
	e, _ := newLiveBrowserEnv(t)
	g := e.connect()
	out, sc := e.call(g, "browse_open", map[string]any{"url": liveLogin})
	if code(sc) != "" || strings.Contains(out, "hunter2") || !strings.Contains(out, `"masked": true`) {
		t.Fatalf("login page: %s", out)
	}
	pass := refByLabel(sc["page"].(map[string]any), "fields", "Password")
	if _, sc := e.call(g, "browse_act", map[string]any{"url": liveLogin, "steps": []map[string]any{{"fill": pass, "value": "x"}}}); code(sc) != browser.CodeSecretField {
		t.Errorf("typing a password: %v", sc)
	}
	if _, sc := e.call(g, "browse_open", map[string]any{"url": liveOther}); code(sc) != browser.CodeSiteNotListed {
		t.Errorf("another site: %v", sc)
	}
}

// refByLabel is the ref of the field or control with that label in a
// page as the agent reads it.
func refByLabel(page map[string]any, kind, label string) string {
	items, _ := page[kind].([]any)
	for _, it := range items {
		m, _ := it.(map[string]any)
		if m["label"] == label {
			return m["ref"].(string)
		}
	}
	return ""
}

// The owner signs in to the real shop through the view planque serves:
// Tab to the email field, type, Tab to the password (the shop's test page
// fills it), Enter. Done seals the cookies, and the taf's next call is
// signed in. No page is read for an agent along the way.
func TestBrowserLiveSignInViewSignsTheShopIn(t *testing.T) {
	e, _ := newLiveBrowserEnv(t)
	g := e.connect()
	st, out := e.startView(g.ID, map[string]any{"url": liveLogin})
	if st != 200 {
		t.Fatalf("start: %d %v", st, out)
	}
	viewURL := out["view"].(string)
	st, hdr, frame := e.view("POST", viewURL+"/frame", nil)
	if st != 200 || hdr.Get("Content-Type") != "image/jpeg" || !bytes.HasPrefix(frame, []byte{0xff, 0xd8}) {
		t.Fatalf("frame: %d %s %d bytes: %s", st, hdr.Get("Content-Type"), len(frame), frame)
	}
	for _, in := range []map[string]any{
		{"kind": "tap", "x": 640, "y": 450},
		{"kind": "key", "key": "Tab"},
		{"kind": "type", "text": "ana@example.com"},
		{"kind": "key", "key": "Tab"},
		{"kind": "key", "key": "Enter"},
	} {
		if st, _, raw := e.view("POST", viewURL+"/input", in); st != 200 {
			t.Fatalf("input %v: %d %s", in, st, raw)
		}
	}
	time.Sleep(700 * time.Millisecond)
	if st, _, raw := e.view("POST", viewURL+"/done", nil); st != 200 {
		t.Fatalf("done: %d %s", st, raw)
	}
	out2, sc := e.call(g, "browse_open", map[string]any{"url": liveItem})
	if code(sc) != "" || !strings.Contains(out2, "Signed in as ana") {
		t.Fatalf("the taf after the owner's sign-in: %s", out2)
	}
}
