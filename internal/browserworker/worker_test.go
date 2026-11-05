package browserworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// shop is the test shop in HTML, for the real browser: the same pages the
// fake worker scripts, served on loopback. It never touches a merchant.
func shop(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	page := func(title, body string) string {
		return "<!doctype html><html><head><title>" + title + "</title></head><body>" + body + "</body></html>"
	}
	signed := func(r *http.Request) string {
		if c, err := r.Cookie("session"); err == nil && c.Value == "signed-in" {
			return "<p>Signed in as ana</p>"
		}
		return ""
	}
	mux.HandleFunc("/item/pen", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Blue pen", signed(r)+`<h1>Blue pen</h1><p>1.25 EUR each. In stock.</p>
<form method=post action=/cart><label>Quantity <input name=qty type=number value=1></label><button>Add to cart</button></form>
<a href=/cart>Cart</a> <a href=/login>Sign in</a>`))
	})
	mux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		qty := "2"
		if r.Method == http.MethodPost {
			_ = r.ParseForm()
			qty = r.Form.Get("qty")
			http.Redirect(w, r, "/cart?qty="+qty, http.StatusSeeOther)
			return
		}
		if q := r.URL.Query().Get("qty"); q != "" {
			qty = q
		}
		fmt.Fprint(w, page("Your cart", `<p>Blue pen × `+qty+`: 2.50 EUR</p><a href=/checkout>Checkout</a>`))
	})
	mux.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Checkout", `<h1>Checkout</h1><table><tr><td>Blue pen × 2</td><td>2.50 EUR</td></tr><tr><td>Subtotal</td><td>2.50 EUR</td></tr><tr><td>Total</td><td>2.50 EUR</td></tr></table>
<form method=post action=/place>
<label>Email <input name=email type=email></label>
<label>Card number <input name=cardnumber autocomplete=cc-number value="4242 4242 4242 4242"></label>
<label>Security code <input name=cvc autocomplete=cc-csc value=123></label>
<label>Coupon <input name=coupon></label>
<button formaction=/checkout>Apply coupon</button>
<button>Place order</button></form>`))
	})
	mux.HandleFunc("/place", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/done", http.StatusSeeOther)
	})
	mux.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Thank you", `<p>Order 1001 placed. Total 2.50 EUR.</p>`))
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_ = r.ParseForm()
			if r.Form.Get("pass") == "hunter2" {
				http.SetCookie(w, &http.Cookie{Name: "session", Value: "signed-in", Path: "/"})
			}
			http.Redirect(w, r, "/item/pen", http.StatusSeeOther)
			return
		}
		fmt.Fprint(w, page("Sign in", `<form method=post action=/login><label>Email <input name=user type=email></label>
<label>Password <input name=pass type=password value=hunter2></label><button>Sign in</button></form>`))
	})
	mux.HandleFunc("/einklick", func(w http.ResponseWriter, r *http.Request) {
		// A one-click buy with no words on it: an icon in a form that posts
		// to /buy. Not a checkout page: no total, no card field.
		fmt.Fprint(w, page("Blauer Stift", `<h1>Blauer Stift</h1><p>1,25 € pro Stück.</p>
<form method=post action=/buy><button type=submit><img src="data:image/gif;base64,R0lGODlhAQABAAAAACw=" alt=""></button></form>
<form method=post action=/merken><button type=submit>Merken</button></form>`))
	})
	mux.HandleFunc("/commande", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Commande", `<h1>Votre commande</h1><p>Stylo bleu × 2</p><p>Sous-total : 2,50 €</p><p>Montant total : 2,50 €</p>
<form method=post action=/commande/valider><label>Courriel <input name=email type=email></label><button>Commander</button></form>`))
	})
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=secret.txt")
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprint(w, "a download")
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func worker(t *testing.T) *Local {
	t.Helper()
	if _, err := chromePath(""); err != nil {
		t.Skip(err)
	}
	return &Local{Options{AllowLocal: true, NavigateWait: 20 * time.Second}}
}

// startFor is a session on the test shop, limited to sites when given.
func startFor(s *httptest.Server, sites ...string) browser.Start {
	return browser.Start{Grant: "grt_test", Sites: sites}
}

// find is the field or control with that label.
func find[T browser.Field | browser.Control](items []T, label string) *T {
	for i := range items {
		var l string
		switch v := any(items[i]).(type) {
		case browser.Field:
			l = v.Label
		case browser.Control:
			l = v.Label
		}
		if l == label {
			return &items[i]
		}
	}
	return nil
}

func mustLink(t *testing.T, p browser.Page, label string) string {
	t.Helper()
	for _, l := range p.Links {
		if l.Label == label {
			return l.Ref
		}
	}
	t.Fatalf("no link %q: %+v", label, p.Links)
	return ""
}

func stepFill(ref, value string) browser.Step { return browser.Step{Fill: ref, Value: value} }
func stepClick(ref string) browser.Step       { return browser.Step{Click: ref} }

// 1. The shop as the browser reads it: refs, kinds from autocomplete,
// secrets masked at the source, the checkout's total, a screenshot, and
// the cookies as the state.
func TestWorkerReadsTheShopAndMasksSecrets(t *testing.T) {
	w, s := worker(t), shop(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	p, err := sess.Open(ctx, s.URL+"/item/pen")
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Blue pen" || !strings.Contains(p.Text, "1.25 EUR") || len(p.Fields) != 1 || p.Fields[0].Kind != "number" || p.Fields[0].Label != "Quantity" {
		t.Fatalf("item page: %+v", p)
	}
	add := find(p.Controls, "Add to cart")
	if add == nil || !add.Submits || add.PlacesOrder {
		t.Fatalf("add to cart: %+v", p.Controls)
	}
	p, err = sess.Do(ctx, stepFill(p.Fields[0].Ref, "2"))
	if err != nil {
		t.Fatal(err)
	}
	if p.Fields[0].Value != "2" {
		t.Errorf("after the fill: %+v", p.Fields[0])
	}
	p, err = sess.Do(ctx, stepClick(add.Ref))
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Your cart" || !strings.Contains(p.Text, "× 2") || !strings.HasSuffix(p.URL, "/cart?qty=2") {
		t.Fatalf("after add to cart: %+v", p)
	}
	p, err = sess.Do(ctx, stepClick(mustLink(t, p, "Checkout")))
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Checkout" || p.Checkout == nil || p.Checkout.Total != "2.50" || p.Checkout.Currency != "EUR" {
		t.Fatalf("checkout: %+v", p)
	}
	kinds := map[string]string{}
	for _, f := range p.Fields {
		kinds[f.Label] = f.Kind
		if (f.Kind == "card" || f.Kind == "code") && (!f.Masked || f.Value != "") {
			t.Errorf("secret shown: %+v", f)
		}
	}
	// A security code is card data: kind card, as secret as the number.
	if kinds["Card number"] != "card" || kinds["Security code"] != "card" || kinds["Email"] != "email" || kinds["Coupon"] != "text" {
		t.Errorf("kinds: %v", kinds)
	}
	raw, _ := json.Marshal(p)
	for _, secret := range []string{"4242", "123"} {
		if strings.Contains(string(raw), secret) {
			t.Errorf("the page carries %q: %s", secret, raw)
		}
	}
	place, coupon := find(p.Controls, "Place order"), find(p.Controls, "Apply coupon")
	if place == nil || !place.Submits || !place.PlacesOrder {
		t.Errorf("place order: %+v", place)
	}
	if coupon == nil || !coupon.Submits {
		t.Errorf("apply coupon: %+v", coupon)
	}
	png, err := sess.Screenshot(ctx)
	if err != nil || !bytes.HasPrefix(png, []byte("\x89PNG")) {
		t.Errorf("screenshot: %v %d bytes", err, len(png))
	}
	// The worker itself places the order when asked to click: the guard
	// is planque's. The confirmation page reads as one.
	p, err = sess.Do(ctx, stepClick(place.Ref))
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Thank you" || !strings.Contains(p.Text, "Order 1001") {
		t.Errorf("after placing: %+v", p)
	}
	p, err = sess.Open(ctx, s.URL+"/login")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ = json.Marshal(p)
	if strings.Contains(string(raw), "hunter2") || find(p.Fields, "Password") == nil || find(p.Fields, "Password").Kind != "password" {
		t.Errorf("login page: %s", raw)
	}
	if find(p.Controls, "Sign in").PlacesOrder {
		t.Errorf("sign in places an order: %+v", p.Controls)
	}
	st, err := sess.Close(ctx)
	if err != nil || !strings.Contains(string(st), `"cookies"`) {
		t.Errorf("state: %v %s", err, st)
	}
}

// 2. Placing an order is a fact read from the page, in no language: an
// unlabelled icon button posting to /buy on a page that isn't a checkout,
// and "Commander" on a French checkout whose total is "2,50 €".
// (agentgabo's ask on kiff/planque#88.)
func TestWorkerReadsPlacesOrderWithoutWords(t *testing.T) {
	w, s := worker(t), shop(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	p, err := sess.Open(ctx, s.URL+"/einklick")
	if err != nil {
		t.Fatal(err)
	}
	if p.Checkout != nil {
		t.Errorf("a product page read as a checkout: %+v", p.Checkout)
	}
	var icon, merken *browser.Control
	for i := range p.Controls {
		switch p.Controls[i].Label {
		case "":
			icon = &p.Controls[i]
		case "Merken":
			merken = &p.Controls[i]
		}
	}
	if icon == nil || !icon.PlacesOrder || !icon.Submits {
		t.Errorf("the icon button: %+v", p.Controls)
	}
	if merken == nil || merken.PlacesOrder {
		t.Errorf("merken: %+v", merken)
	}
	p, err = sess.Open(ctx, s.URL+"/commande")
	if err != nil {
		t.Fatal(err)
	}
	if p.Checkout == nil || p.Checkout.Total != "2.50" || p.Checkout.Currency != "EUR" {
		t.Errorf("the French checkout: %+v", p.Checkout)
	}
	if c := find(p.Controls, "Commander"); c == nil || !c.PlacesOrder {
		t.Errorf("commander: %+v", p.Controls)
	}
}

// 3. A sign-in made in one session is the state the next starts with.
func TestWorkerKeepsASignInAsTheState(t *testing.T) {
	w, s := worker(t), shop(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	p, err := sess.Open(ctx, s.URL+"/login")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sess.Do(ctx, stepFill(find(p.Fields, "Email").Ref, "ana@example.com")); err != nil {
		t.Fatal(err)
	}
	// The password is already in the form (the owner typed it in the live
	// view); the worker's test only submits.
	p, err = sess.Do(ctx, stepClick(find(p.Controls, "Sign in").Ref))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(p.Text, "Signed in as ana") {
		t.Fatalf("after sign in: %+v", p)
	}
	st, err := sess.Close(ctx)
	if err != nil || !strings.Contains(string(st), "signed-in") {
		t.Fatalf("state: %v %s", err, st)
	}
	start := startFor(s)
	start.State = st
	again, err := w.Start(ctx, start)
	if err != nil {
		t.Fatal(err)
	}
	defer again.Close(ctx)
	p, err = again.Open(ctx, s.URL+"/item/pen")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(p.Text, "Signed in as ana") {
		t.Errorf("the sign-in didn't carry over: %+v", p)
	}
	// Without the state, a fresh session is signed out.
	fresh, err := w.Start(ctx, startFor(s))
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close(ctx)
	if p, err := fresh.Open(ctx, s.URL+"/item/pen"); err != nil || strings.Contains(p.Text, "Signed in") {
		t.Errorf("a fresh session is signed in: %v %+v", err, p)
	}
}

// 4. The egress holds inside the browser: a listed name that resolves to
// a private address is never fetched, a site off the list isn't, and a
// download is denied.
func TestWorkerEgressHoldsInsideTheBrowser(t *testing.T) {
	w, s := worker(t), shop(t)
	var asked []string
	w.Resolve = func(_ context.Context, host string) ([]netip.Addr, error) {
		asked = append(asked, host)
		return []netip.Addr{netip.MustParseAddr("10.0.0.7")}, nil
	}
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(s, "shop.example"))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	if _, err := sess.Open(ctx, "https://shop.example/item/pen"); err == nil || !strings.Contains(err.Error(), "egress") {
		t.Errorf("a listed name resolving to 10.x: %v", err)
	}
	if _, err := sess.Open(ctx, "https://other.example/"); err == nil || !strings.Contains(err.Error(), "egress") {
		t.Errorf("an unlisted site: %v", err)
	}
	// The browser retries a refused tunnel, so the counts are floors; the
	// unlisted site is never resolved at all.
	for _, h := range asked {
		if h != "shop.example" {
			t.Errorf("resolved %q: the unlisted site must not be resolved", h)
		}
	}
	if len(asked) == 0 {
		t.Error("the listed name wasn't resolved by the proxy")
	}
	if n := sess.(*session).proxy.Refused; n < 2 {
		t.Errorf("refused %d, want at least 2", n)
	}
	// A download: denied, the page stays where it was.
	if _, err := sess.Open(ctx, s.URL+"/item/pen"); err != nil {
		t.Fatal(err)
	}
	if p2, err := sess.Open(ctx, s.URL+"/file"); err == nil && strings.Contains(p2.Text, "a download") {
		t.Errorf("the download was shown: %+v", p2)
	}
}
