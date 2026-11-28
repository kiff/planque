package fakes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// The test shop in HTML (kiff/planque#86): the pages the scripted Browser
// fakes, served for the real browser worker. planque's end to end runs on
// it and never touches a merchant. Totals and the currency can be changed
// between calls, so a pay after a change is tested.
type Shop struct {
	*httptest.Server
	// Total and Currency are what the checkout shows now.
	Total, Currency string
	// Orders are the orders placed: one per "Place order".
	Orders []string
	// Ship is the delivery address the pay checkout (/pay) shows: its
	// first line, postal code and city. Paid is the last four digits of
	// each card the pay checkout was paid with, as the merchant receives
	// them; PayCalls counts the pay checkout's orders.
	Ship     [3]string
	Paid     []string
	PayCalls int
	// Exp is each paid order's expiry as received: "MM / YY", or the
	// month and year fields joined by "/".
	Exp []string
}

// NewShop serves the shop on loopback, plain http; NewTLSShop over TLS
// with a test certificate, for a browser reaching it as a site.
func NewShop() *Shop {
	s := &Shop{Total: "2.50", Currency: "EUR", Ship: defaultShip}
	s.Server = httptest.NewServer(s.handler())
	return s
}

func NewTLSShop() *Shop {
	s := &Shop{Total: "2.50", Currency: "EUR", Ship: defaultShip}
	s.Server = httptest.NewTLSServer(s.handler())
	return s
}

// defaultShip is the pay checkout's delivery address: the contract's test
// case's.
var defaultShip = [3]string{"12 rue des Lilas", "69003", "Lyon"}

func (s *Shop) handler() http.Handler {
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
		fmt.Fprint(w, page("Checkout", `<h1>Checkout</h1><table class=order>
<tr class=line><td>Blue pen × 2</td><td>2.50 EUR</td></tr>
<tr><td>Subtotal</td><td>2.50 EUR</td></tr>
<tr><td>Shipping</td><td>0.00 EUR</td></tr>
<tr><td>Total</td><td>`+s.Total+` `+s.Currency+`</td></tr></table>
<form method=post action=/place>
<label>Email <input name=email type=email></label>
<label>Card number <input name=cardnumber autocomplete=cc-number value="4242 4242 4242 4242"></label>
<label>Security code <input name=cvc autocomplete=cc-csc value=123></label>
<label>Coupon <input name=coupon></label>
<button formaction=/checkout>Apply coupon</button>
<button>Place order</button></form>`))
	})
	mux.HandleFunc("/place", func(w http.ResponseWriter, r *http.Request) {
		s.Orders = append(s.Orders, "1001")
		http.Redirect(w, r, "/done", http.StatusSeeOther)
	})
	mux.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Thank you", `<p>Order 1001 placed. Total `+s.Total+` `+s.Currency+`.</p>`))
	})
	// The pay checkout, for the payment-only worker: the delivery address
	// shown (as text, or as the address fields with ?fields=1), empty card
	// fields by their autocomplete tokens, and one control that places the
	// order. ?two=1 adds a second one.
	mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		ship := `<p class=ship>Deliver to: Acme SAS, ` + s.Ship[0] + `, ` + s.Ship[1] + ` ` + s.Ship[2] + `, France</p>`
		if r.URL.Query().Get("fields") == "1" {
			ship = `<fieldset><legend>Delivery</legend>
<label>Address <input autocomplete="shipping address-line1" value="` + s.Ship[0] + `"></label>
<label>Postal code <input autocomplete="shipping postal-code" value="` + s.Ship[1] + `"></label>
<label>City <input autocomplete="shipping address-level2" value="` + s.Ship[2] + `"></label>
<label>Country <select autocomplete="shipping country"><option value=FR selected>France</option><option value=DE>Germany</option></select></label></fieldset>`
		}
		second := ""
		if r.URL.Query().Get("two") == "1" {
			second = `<button formaction=/pay/place>Buy now</button>`
		}
		// ?end= picks the page after the order: "failed" (a payment that
		// didn't go through) or "plain" (no confirmation, no reference).
		action := "/pay/place"
		if end := r.URL.Query().Get("end"); end != "" {
			action += "?end=" + end
		}
		card := `<label>Card number <input name=cardnumber autocomplete=cc-number inputmode=numeric></label>
<label>Expiry <input name=exp autocomplete=cc-exp></label>
<label>Security code <input name=cvc autocomplete=cc-csc></label>`
		switch {
		case r.URL.Query().Get("split") == "1":
			card = `<label>Card number <input name=cardnumber autocomplete=cc-number inputmode=numeric></label>
<label>Month <input name=expm autocomplete=cc-exp-month maxlength=2></label>
<label>Year <input name=expy autocomplete=cc-exp-year maxlength=2 placeholder=YY></label>
<label>Security code <input name=cvc autocomplete=cc-csc></label>`
		case r.URL.Query().Get("frame") == "1":
			card = `<iframe name=__privateStripeFrame title="Secure card payment input frame" src="https://js.stripe.com/v3/elements-inner-card.html" width=400 height=60></iframe>`
		}
		fmt.Fprint(w, page("Checkout", `<script type="application/ld+json">{"@context":"https://schema.org","@type":"OnlineStore","name":"Test Pens","address":{"@type":"PostalAddress","addressCountry":"FR"}}</script>
<h1>Checkout</h1><table class=order>
<tr class=line><td>Blue pen × 2</td><td>2.50 EUR</td></tr>
<tr><td>Shipping</td><td>0.00 EUR</td></tr>
<tr><td>VAT</td><td>0.42 EUR</td></tr>
<tr><td>Total</td><td>`+s.Total+` `+s.Currency+`</td></tr></table>`+ship+`
<form method=post action=`+action+`>
`+card+`
<label>Name on card <input name=ccname autocomplete=cc-name></label>
<input type=hidden name=trap autocomplete=cc-number>
<button>Place order</button>`+second+`</form>`))
	})
	mux.HandleFunc("/pay/place", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		num := strings.Map(func(c rune) rune {
			if c >= '0' && c <= '9' {
				return c
			}
			return -1
		}, r.Form.Get("cardnumber"))
		exp := r.Form.Get("exp")
		if exp == "" && r.Form.Get("expm") != "" {
			exp = r.Form.Get("expm") + "/" + r.Form.Get("expy")
		}
		if len(num) < 13 || exp == "" || r.Form.Get("cvc") == "" {
			fmt.Fprint(w, page("Checkout", `<p>Your card was declined.</p><form method=post action=/pay/place><label>Card number <input name=cardnumber autocomplete=cc-number></label><button>Place order</button></form>`))
			return
		}
		s.PayCalls++
		s.Paid = append(s.Paid, num[len(num)-4:])
		s.Exp = append(s.Exp, exp)
		switch r.URL.Query().Get("end") {
		case "failed":
			http.Redirect(w, r, "/pay/failed", http.StatusSeeOther)
		case "plain":
			http.Redirect(w, r, "/pay/next", http.StatusSeeOther)
		case "unprocessed", "refonly":
			http.Redirect(w, r, "/pay/"+r.URL.Query().Get("end"), http.StatusSeeOther)
		default:
			http.Redirect(w, r, "/pay/done", http.StatusSeeOther)
		}
	})
	mux.HandleFunc("/pay/failed", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Payment", fmt.Sprintf(`<p>Your payment didn't go through. Order A-%d was not completed.</p>`, 1000+s.PayCalls)))
	})
	mux.HandleFunc("/pay/unprocessed", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Order", fmt.Sprintf(`<p>Order A-%d couldn't be processed.</p>`, 1000+s.PayCalls)))
	})
	mux.HandleFunc("/pay/refonly", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Order", fmt.Sprintf(`<p>Order A-%d</p>`, 1000+s.PayCalls)))
	})
	mux.HandleFunc("/pay/next", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Next steps", `<p>We'll email you shortly.</p>`))
	})
	mux.HandleFunc("/pay/done", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Thank you", fmt.Sprintf(`<p>Thank you. Order A-%d placed. Total %s %s.</p>`, 1000+s.PayCalls, s.Total, s.Currency)))
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
		fmt.Fprint(w, page("Commande", `<h1>Votre commande</h1><ul class=panier><li>Stylo bleu × 2 — 2,50 €</li></ul><p>Sous-total : 2,50 €</p><p>Montant total : 2,50 €</p>
<form method=post action=/commande/valider><label>Courriel <input name=email type=email></label><button>Commander</button></form>`))
	})
	mux.HandleFunc("/frame", func(w http.ResponseWriter, r *http.Request) {
		// A page embedding a frame from wherever src says: another site, to
		// a browser, when its host differs.
		fmt.Fprint(w, page("Framed", `<h1>Framed</h1><iframe id=f src="`+r.URL.Query().Get("src")+`" width=400 height=200></iframe>`))
	})
	mux.HandleFunc("/inner", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page("Inner", `<p>Inner frame</p>`))
	})
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=secret.txt")
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprint(w, "a download")
	})
	return mux
}
