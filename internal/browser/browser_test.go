package browser

import (
	"encoding/json"
	"strings"
	"testing"
)

// A checkout page as a worker reads it: a card field, a security code,
// and a control that places the order.
func checkoutPage() Page {
	return Page{URL: "https://shop.example/checkout", Title: "Checkout",
		Fields: []Field{
			{Ref: "f1", Label: "Email", Kind: "email", Value: "ana@example.com"},
			{Ref: "f2", Label: "Card number", Kind: "card", Value: "4242 4242 4242 4242"},
			{Ref: "f3", Label: "Security code", Kind: "code", Value: "123"},
			{Ref: "f4", Label: "Password", Kind: "password", Value: "hunter2"},
		},
		Controls: []Control{{Ref: "c1", Label: "Place order", Submits: true}, {Ref: "c2", Label: "Apply coupon", Submits: true}, {Ref: "c3", Label: "Back"}},
		Links:    []Link{{Ref: "l1", Label: "Terms", URL: "https://shop.example/terms"}}}
}

// Every step that would pay, or type a secret, is refused with its code;
// the page's other steps run.
func TestRefuseKeepsMoneyAndSecretsOutOfSteps(t *testing.T) {
	p := checkoutPage()
	for _, c := range []struct {
		name string
		step Step
		code string
	}{
		{"placing the order", Step{Click: "c1"}, CodePaymentStep},
		{"submitting on a payment page", Step{Click: "c2"}, CodePaymentStep},
		{"typing a card number", Step{Fill: "f2", Value: "4000"}, CodePaymentStep},
		{"typing a security code", Step{Fill: "f3", Value: "999"}, CodeSecretField},
		{"typing a password", Step{Fill: "f4", Value: "x"}, CodeSecretField},
		{"selecting in a password", Step{Select: "f4", Option: "x"}, CodeSecretField},
		{"a field that isn't there", Step{Fill: "f9", Value: "x"}, CodeRefUnknown},
		{"a control that isn't there", Step{Click: "c9"}, CodeRefUnknown},
		{"typing an email", Step{Fill: "f1", Value: "b@example.com"}, ""},
		{"a plain click on a payment page", Step{Click: "c3"}, CodePaymentStep},
		{"a link", Step{Click: "l1"}, ""},
	} {
		code, why := Refuse(p, c.step)
		if code != c.code {
			t.Errorf("%s: %q (%s), want %q", c.name, code, why, c.code)
		}
		if code != "" && why == "" {
			t.Errorf("%s: a refusal with no words", c.name)
		}
	}
	// On a page that doesn't pay, a submit runs.
	form := Page{Fields: []Field{{Ref: "f1", Label: "Name", Kind: "text"}}, Controls: []Control{{Ref: "c1", Label: "Search", Submits: true}}}
	if code, why := Refuse(form, Step{Click: "c1"}); code != "" {
		t.Errorf("a search submit: %s %s", code, why)
	}
	// A control whose words pay is refused even off a payment page.
	buy := Page{Controls: []Control{{Ref: "c1", Label: "Buy now"}}}
	if code, _ := Refuse(buy, Step{Click: "c1"}); code != CodePaymentStep {
		t.Errorf("buy now: %s", code)
	}
	// On a checkout, every control is refused, whatever it says and in
	// whatever language, labelled or not; its links are followed.
	fr := Page{Checkout: &Checkout{Total: "2.50", Currency: "EUR"},
		Controls: []Control{{Ref: "cmd", Label: "Commander", Submits: true}, {Ref: "icon", Label: ""}, {Ref: "zb", Label: "Jetzt kaufen"}, {Ref: "cp", Label: "Code promo", Submits: true}},
		Links:    []Link{{Ref: "cgv", Label: "CGV", URL: "https://shop.example/cgv"}}}
	for _, ref := range []string{"cmd", "icon", "zb", "cp"} {
		if code, _ := Refuse(fr, Step{Click: ref}); code != CodePaymentStep {
			t.Errorf("checkout control %s: %s", ref, code)
		}
	}
	if code, why := Refuse(fr, Step{Click: "cgv"}); code != "" {
		t.Errorf("a checkout's link: %s %s", code, why)
	}
	// A one-click buy on a product page, in any of the languages.
	for _, label := range []string{"Kaufen", "Comprar", "Acheter", "Acquista ora", "Koop nu", "Comprar agora", "Commander", "Bestellen", "Pagar"} {
		one := Page{Controls: []Control{{Ref: "c1", Label: label}}}
		if code, _ := Refuse(one, Step{Click: "c1"}); code != CodePaymentStep {
			t.Errorf("one-click %q: %s", label, code)
		}
	}
	// The worker's fact, with no words at all.
	fact := Page{Controls: []Control{{Ref: "c1", Label: "", PlacesOrder: true}}}
	if code, _ := Refuse(fact, Step{Click: "c1"}); code != CodePaymentStep {
		t.Errorf("places_order fact: %s", code)
	}
}

func TestPaymentPages(t *testing.T) {
	if !Payment(checkoutPage()) {
		t.Error("a checkout with a card field isn't a payment page")
	}
	if !Payment(Page{PaymentFrame: true}) {
		t.Error("a payment frame isn't a payment page")
	}
	if !Payment(Page{Controls: []Control{{Ref: "c", Label: "Confirm and pay"}}}) {
		t.Error("confirm and pay isn't a payment page")
	}
	if Payment(Page{Fields: []Field{{Ref: "f", Label: "Password", Kind: "password"}}, Controls: []Control{{Ref: "c", Label: "Sign in", Submits: true}}}) {
		t.Error("a sign-in page is taken for a payment page")
	}
	for _, label := range []string{"Pay now", "Place your order", "Buy", "Complete purchase", "Subscribe", "PAY"} {
		if !Pays(Control{Label: label}) {
			t.Errorf("%q doesn't pay", label)
		}
	}
	for _, label := range []string{"Zur Kasse gehen", "Panier", "Suchen", "Añadir al carrito", "Aggiungi al carrello", "Weiter"} {
		if Pays(Control{Label: label}) {
			t.Errorf("%q pays", label)
		}
	}
	for _, label := range []string{"Search", "Add to cart", "Payment methods", "Sign in", "Repay history"} {
		if Pays(Control{Label: label}) {
			t.Errorf("%q pays", label)
		}
	}
}

// A secret's value never reaches an agent, whatever the worker sent.
func TestMaskHidesSecretsAndCutsText(t *testing.T) {
	m := Mask(checkoutPage(), 0)
	for _, f := range m.Fields {
		if Secret(f) && (f.Value != "" || !f.Masked) {
			t.Errorf("%s: %+v", f.Label, f)
		}
		if !Secret(f) && f.Value == "" {
			t.Errorf("%s was masked", f.Label)
		}
	}
	for _, s := range []string{"4242", "123", "hunter2"} {
		if b := mustJSON(m); strings.Contains(b, s) {
			t.Errorf("the masked page carries %q", s)
		}
	}
	long := Mask(Page{Text: strings.Repeat("a", 100)}, 10)
	if long.Text != strings.Repeat("a", 10)+"\n(cut)" {
		t.Errorf("text wasn't cut: %q", long.Text)
	}
}

func TestAllowedSitesAndAddresses(t *testing.T) {
	sites := []string{"shop.example", "docs.example.org"}
	for _, c := range []struct{ url, code string }{
		{"https://shop.example/item/1", ""},
		{"https://www.shop.example/", ""},
		{"https://docs.example.org/page", ""},
		{"https://other.example/", CodeSiteNotListed},
		{"https://shop.example.evil/", CodeSiteNotListed},
		{"https://notshop.example/", CodeSiteNotListed},
		{"http://shop.example/", CodeBadURL},
		{"https://user:pw@shop.example/", CodeBadURL},
		{"file:///etc/passwd", CodeBadURL},
		{"", CodeBadURL},
		{"https://127.0.0.1/", CodeAddressPrivate},
		{"https://10.1.2.3/", CodeAddressPrivate},
		{"https://8.8.8.8/", CodeAddressPrivate},
		{"https://[::1]/", CodeAddressPrivate},
		{"https://localhost/", CodeAddressPrivate},
		{"https://planque/", CodeAddressPrivate},
		{"https://db.internal/", CodeAddressPrivate},
		{"https://printer.local/", CodeAddressPrivate},
	} {
		if code, why := Allowed(c.url, sites); code != c.code {
			t.Errorf("%s: %q (%s), want %q", c.url, code, why, c.code)
		}
	}
	// With no sites listed, any public site; local names and addresses still never.
	if code, _ := Allowed("https://anything.example/", nil); code != "" {
		t.Errorf("no list: %s", code)
	}
	if code, _ := Allowed("https://169.254.169.254/latest/", nil); code != CodeAddressPrivate {
		t.Errorf("no list, metadata address: %s", code)
	}
}

func TestCheckSite(t *testing.T) {
	for _, ok := range []string{"shop.example", "a-b.co.uk", "x1.example.org"} {
		if err := CheckSite(ok); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	for _, bad := range []string{"https://shop.example", "shop.example/path", "shop.example:443", "Shop.Example", "localhost", "intranet", "10.0.0.1", "db.internal", "-a.example", ""} {
		if err := CheckSite(bad); err == nil {
			t.Errorf("%q was taken for a site", bad)
		}
	}
}

func TestStepCheck(t *testing.T) {
	for _, ok := range []Step{{Fill: "f", Value: "x"}, {Fill: "f"}, {Select: "f", Option: "o"}, {Click: "c"}} {
		if err := ok.Check(); err != nil {
			t.Errorf("%+v: %v", ok, err)
		}
	}
	for _, bad := range []Step{{}, {Fill: "f", Click: "c"}, {Click: "c", Value: "x"}, {Select: "f"}, {Fill: "f", Value: strings.Repeat("x", 4097)}} {
		if err := bad.Check(); err == nil {
			t.Errorf("%+v passed", bad)
		}
	}
}

func TestCheckoutSame(t *testing.T) {
	a := Checkout{Lines: []Line{{Item: "Pen", Quantity: 2, Price: "1.25"}}, Total: "2.50", Currency: "EUR"}
	b := a
	b.Lines = []Line{{Item: "Pen", Quantity: 2, Price: "1.25"}}
	if !a.Same(b) {
		t.Error("the same checkout differs")
	}
	for _, o := range []Checkout{
		{Lines: a.Lines, Total: "2.51", Currency: "EUR"},
		{Lines: a.Lines, Total: "2.50", Currency: "USD"},
		{Lines: []Line{{Item: "Pen", Quantity: 3, Price: "1.25"}}, Total: "2.50", Currency: "EUR"},
		{Lines: nil, Total: "2.50", Currency: "EUR"},
	} {
		if a.Same(o) {
			t.Errorf("%+v taken as the same", o)
		}
	}
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
