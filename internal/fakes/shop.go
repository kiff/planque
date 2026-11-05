package fakes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
}

// NewShop serves the shop on loopback, plain http; NewTLSShop over TLS
// with a test certificate, for a browser reaching it as a site.
func NewShop() *Shop {
	s := &Shop{Total: "2.50", Currency: "EUR"}
	s.Server = httptest.NewServer(s.handler())
	return s
}

func NewTLSShop() *Shop {
	s := &Shop{Total: "2.50", Currency: "EUR"}
	s.Server = httptest.NewTLSServer(s.handler())
	return s
}

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
