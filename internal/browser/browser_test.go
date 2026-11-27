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
		{"a plain link on a payment page", Step{Click: "l1"}, CodePaymentStep},
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
	// whatever language, labelled or not, and so is every link.
	fr := Page{Checkout: &Checkout{Total: "2.50", Currency: "EUR"},
		Controls: []Control{{Ref: "cmd", Label: "Commander", Submits: true}, {Ref: "icon", Label: ""}, {Ref: "zb", Label: "Jetzt kaufen"}, {Ref: "cp", Label: "Code promo", Submits: true}},
		Links:    []Link{{Ref: "cgv", Label: "CGV", URL: "https://shop.example/cgv"}}}
	for _, ref := range []string{"cmd", "icon", "zb", "cp"} {
		if code, _ := Refuse(fr, Step{Click: ref}); code != CodePaymentStep {
			t.Errorf("checkout control %s: %s", ref, code)
		}
	}
	if code, why := Refuse(fr, Step{Click: "cgv"}); code != CodePaymentStep {
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

// A card number shown in a page's text, or echoed in a plain field, is
// masked to its last four digits (#86: a secret in a page's text); other
// numbers stay.
func TestMaskTextHidesCardNumbersInText(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"Paying with card 4242 4242 4242 4242 ending soon", "Paying with card •••• 4242 ending soon"},
		{"card 4242424242424242.", "card •••• 4242."},
		{"card 4242-4242-4242-4242", "card •••• 4242"},
		{"Amex 3782 822463 10005", "Amex •••• 0005"},
		{"Order 1001, total 2.50 EUR, ref 123456789012", "Order 1001, total 2.50 EUR, ref 123456789012"}, // fails Luhn
		{"Phone +33 6 12 34 56 78", "Phone +33 6 12 34 56 78"},
		{"IBAN-like 1234567890123456", "IBAN-like 1234567890123456"}, // fails Luhn
		{"", ""},
	} {
		if got := MaskText(c.in); got != c.want {
			t.Errorf("%q: %q, want %q", c.in, got, c.want)
		}
	}
	p := Mask(Page{Title: "Receipt 4242 4242 4242 4242", Text: "Card 4242 4242 4242 4242 charged",
		Fields: []Field{{Ref: "f1", Label: "Note", Kind: "text", Value: "paid with 4242424242424242"}}}, 0)
	if strings.Contains(mustJSON(p), "4242 4242") || strings.Contains(mustJSON(p), "4242424242424242") {
		t.Errorf("the card number shows: %s", mustJSON(p))
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

// A link can run the page's script instead of going anywhere, and a plain
// one can place an order (a GET to /place-order): one whose words pay is
// refused anywhere, and on a checkout or payment page no link is followed,
// whatever its address or words. Off them, a link is followed.
func TestALinkThatPaysOrIsOnACheckoutIsRefused(t *testing.T) {
	checkout := checkoutPage()
	checkout.Links = append(checkout.Links,
		Link{Ref: "l2", Label: "Continue", URL: "https://shop.example/checkout#"},
		Link{Ref: "l3", Label: "Next step", URL: "https://shop.example/checkout/?step=2"},
		Link{Ref: "l4", Label: "Place order", URL: "https://shop.example/thanks"},
		Link{Ref: "l5", Label: "Continue", URL: "https://shop.example/place-order"})
	product := Page{URL: "https://shop.example/p/1", Title: "Mug",
		Links: []Link{{Ref: "l1", Label: "Buy now", URL: "https://shop.example/p/1#"}, {Ref: "l2", Label: "Reviews", URL: "https://shop.example/p/1#reviews"},
			{Ref: "l3", Label: "Kaufen", URL: "https://shop.example/cart"}, {Ref: "l4", Label: "Shipping", URL: "https://shop.example/shipping"},
			{Ref: "l5", Label: "Continue", URL: "https://shop.example/place-order"}}}
	for _, c := range []struct {
		name string
		page Page
		ref  string
		code string
	}{
		{"a script link on a checkout", checkout, "l2", CodePaymentStep},
		{"the same path with a query, on a checkout", checkout, "l3", CodePaymentStep},
		{"a link whose words pay, to another page", checkout, "l4", CodePaymentStep},
		{"a checkout's link to another page", checkout, "l1", CodePaymentStep},
		{"a neutral link to another path, on a checkout", checkout, "l5", CodePaymentStep},
		{"the same link off a checkout or payment page", product, "l5", ""},
		{"a link whose words pay, on a product page", product, "l1", CodePaymentStep},
		{"a link whose words pay in another language", product, "l3", CodePaymentStep},
		{"a link within a page that doesn't pay", product, "l2", ""},
		{"a plain link", product, "l4", ""},
	} {
		if code, why := Refuse(c.page, Step{Click: c.ref}); code != c.code {
			t.Errorf("%s: %q (%s), want %q", c.name, code, why, c.code)
		}
	}
}

// A select whose options change under the same ref and label is another
// read: a step picks from the list it was built on.
func TestReadCoversASelectsOptions(t *testing.T) {
	p := Page{URL: "https://shop.example/ship", Fields: []Field{{Ref: "f1", Label: "Shipping", Kind: "select", Options: []string{"Standard", "Express"}}}}
	q := p
	q.Fields = []Field{{Ref: "f1", Label: "Shipping", Kind: "select", Options: []string{"Standard", "Express (€30)"}}}
	if Read(p) == Read(q) {
		t.Fatal("a select's new options kept its read")
	}
	r := p
	r.Fields = []Field{{Ref: "f1", Label: "Shipping", Kind: "select", Options: []string{"Standard", "Express"}, Value: "Express"}}
	if Read(p) != Read(r) {
		t.Fatal("a picked value moved the read")
	}
}

// The owner's Chrome's other steps: each well formed only as itself, and refused as the extension's guard refuses it,
// so planque refuses first.
func TestTheChromesOtherStepsAreCheckedAndRefusedAsAClickOrAFillIs(t *testing.T) {
	for _, c := range []struct {
		step Step
		ok   bool
	}{
		{Step{Tick: "f1"}, true},
		{Step{Press: "f1", Key: "Enter"}, true},
		{Step{Press: "f1", Key: "Tab"}, false},
		{Step{Press: "f1"}, false},
		{Step{Fill: "f1", Key: "Enter"}, false},
		{Step{Scroll: "down"}, true},
		{Step{Scroll: "l1"}, true},
		{Step{Back: true}, true},
		{Step{Back: true, Click: "c1"}, false},
		{Step{Fill: "f1", Value: "Pari", Keys: true}, true},
		{Step{Click: "c1", Keys: true}, false},
	} {
		if err := c.step.Check(); (err == nil) != c.ok {
			t.Errorf("%+v: %v", c.step, err)
		}
		if c.ok && c.step.Fill == "" && !c.step.ChromeOnly() {
			t.Errorf("%+v isn't the owner's Chrome's only", c.step)
		}
	}
	p := checkoutPage()
	// No card field and no paying control, so it isn't a payment page; its link "Buy now" still pays.
	plain := Page{URL: "https://shop.example/", Title: "Shop",
		Fields:   []Field{{Ref: "f1", Label: "Search", Kind: "search"}, {Ref: "f2", Label: "Password", Kind: "password"}},
		Controls: []Control{{Ref: "c1", Label: "Search"}},
		Links:    []Link{{Ref: "l1", Label: "Lamp", URL: "https://shop.example/lamp"}, {Ref: "l2", Label: "Buy now", URL: "https://shop.example/buy"}}}
	paying := Page{URL: "https://shop.example/lamp", Title: "Lamp",
		Fields:   []Field{{Ref: "f3", Label: "Card", Kind: "card"}},
		Controls: []Control{{Ref: "c2", Label: "Buy now"}}}
	for _, c := range []struct {
		name string
		page Page
		step Step
		code string
	}{
		{"Enter on a search field", plain, Step{Press: "f1", Key: "Enter"}, ""},
		{"ArrowDown on a search field", plain, Step{Press: "f1", Key: "ArrowDown"}, ""},
		{"a key on a password", plain, Step{Press: "f2", Key: "Escape"}, CodeSecretField},
		{"a key on a card field", paying, Step{Press: "f3", Key: "ArrowUp"}, CodePaymentStep},
		{"Enter on a control that pays", paying, Step{Press: "c2", Key: "Enter"}, CodePaymentStep},
		{"Enter on a link whose words pay", plain, Step{Press: "l2", Key: "Enter"}, CodePaymentStep},
		{"Enter on a plain link", plain, Step{Press: "l1", Key: "Enter"}, ""},
		{"Escape on a control that pays", paying, Step{Press: "c2", Key: "Escape"}, ""},
		{"Enter on a field of a checkout", p, Step{Press: "f1", Key: "Enter"}, CodePaymentStep},
		{"Enter on any control of a checkout", p, Step{Press: "c3", Key: "Enter"}, CodePaymentStep},
		{"a key on a ref that isn't there", plain, Step{Press: "f9", Key: "Escape"}, CodeRefUnknown},
		{"ticking a box", plain, Step{Tick: "f1"}, ""},
		{"ticking a password", plain, Step{Tick: "f2"}, CodeSecretField},
		{"ticking what isn't there", plain, Step{Tick: "f9"}, CodeRefUnknown},
		{"scrolling to a link", plain, Step{Scroll: "l1"}, ""},
		{"scrolling to what isn't there", plain, Step{Scroll: "x9"}, CodeRefUnknown},
		{"scrolling the page", p, Step{Scroll: "up"}, ""},
		{"back", p, Step{Back: true}, ""},
		{"keys into a card field", paying, Step{Fill: "f3", Value: "4", Keys: true}, CodePaymentStep},
		{"keys into a password", plain, Step{Fill: "f2", Value: "x", Keys: true}, CodeSecretField},
	} {
		if code, why := Refuse(c.page, c.step); code != c.code {
			t.Errorf("%s: %q (%s), want %q", c.name, code, why, c.code)
		}
	}
}

// Only what the Chrome may say it cut is kept, each once.
func TestMaskKeepsOnlyKnownCuts(t *testing.T) {
	got := Mask(Page{Truncated: []string{"links", "text", "links", "<<<page-words:x>>>", "options"}}, 0).Truncated
	if strings.Join(got, ",") != "links,text,options" {
		t.Fatalf("truncated: %v", got)
	}
}

// A step names its ref as the read does: the label (spaces folded, case
// kept), submits: true exactly for a control that submits a form, and to: the
// host a link goes to. Those are what a held call shows the owner, so a step
// that names its ref otherwise, a lying taf's or a mistaken one's, runs
// nothing (red-team, 27 Nov: a held browse_act showed its steps as refs the
// owner couldn't read).
func TestMatchesHoldsAStepToThePagesWordsForItsRef(t *testing.T) {
	p := Page{URL: "https://shop.example/form", Title: "Order",
		Fields:   []Field{{Ref: "f1", Label: "Customer name", Kind: "text"}},
		Controls: []Control{{Ref: "c1", Label: "Submit order", Submits: true}, {Ref: "c2", Label: "Save draft"}, {Ref: "c3", Label: " "}},
		Links:    []Link{{Ref: "l1", Label: "Partner", URL: "https://Other.example/x"}, {Ref: "l2", URL: "https://shop.example/y"}}}
	for _, c := range []struct {
		name string
		step Step
		ok   bool
	}{
		{"a fill with its label", Step{Fill: "f1", Value: "Test", Label: "Customer name"}, true},
		{"spaces folded", Step{Fill: "f1", Value: "Test", Label: "  Customer\tname "}, true},
		{"no label", Step{Fill: "f1", Value: "Test"}, false},
		{"another field's words", Step{Fill: "f1", Value: "Test", Label: "Notes"}, false},
		{"case kept", Step{Fill: "f1", Value: "Test", Label: "customer name"}, false},
		{"a submit named as one", Step{Click: "c1", Label: "Submit order", Submits: true}, true},
		{"a submit called a draft", Step{Click: "c1", Label: "Save draft", Submits: true}, false},
		{"a submit without submits", Step{Click: "c1", Label: "Submit order"}, false},
		{"a plain control said to submit", Step{Click: "c2", Label: "Save draft", Submits: true}, false},
		{"a plain control", Step{Click: "c2", Label: "Save draft"}, true},
		{"a link with its host", Step{Click: "l1", Label: "Partner", To: "other.example"}, true},
		{"a link without its host", Step{Click: "l1", Label: "Partner"}, false},
		{"a link to another host", Step{Click: "l1", Label: "Partner", To: "shop.example"}, false},
		{"submits on a field", Step{Fill: "f1", Value: "x", Label: "Customer name", Submits: true}, false},
		{"Enter on a submit names it", Step{Press: "c1", Key: "Enter", Label: "Submit order", Submits: true}, true},
		{"Enter in a field not said to submit", Step{Press: "f1", Key: "Enter", Label: "Customer name"}, false},
		{"Enter in a field said to submit", Step{Press: "f1", Key: "Enter", Label: "Customer name", Submits: true}, true},
		{"another key in a field said to submit", Step{Press: "f1", Key: "Escape", Label: "Customer name", Submits: true}, false},
		{"a line break typed key by key, not said to submit", Step{Fill: "f1", Value: "Ana\nx", Keys: true, Label: "Customer name"}, false},
		{"a line break typed key by key, said to submit", Step{Fill: "f1", Value: "Ana\r", Keys: true, Label: "Customer name", Submits: true}, true},
		{"a line break set at once submits nothing", Step{Fill: "f1", Value: "Ana\nx", Label: "Customer name"}, true},
		{"a control with no words", Step{Click: "c3", Label: " "}, false},
		{"a link with no words", Step{Click: "l2", Label: "", To: "shop.example"}, false},
		{"a scroll to a ref names nothing", Step{Scroll: "c1"}, true},
		{"back names nothing", Step{Back: true}, true},
		{"a ref the page doesn't have is Refuse's", Step{Click: "zz"}, true},
	} {
		code, why := Matches(p, c.step)
		if (code == "") != c.ok || (!c.ok && code != CodeLabelMismatch) {
			t.Errorf("%s: %q %q", c.name, code, why)
		}
	}
	// The refusal names the page's words, so the taf can build the step again.
	if _, why := Matches(p, Step{Fill: "f1", Value: "Test", Label: "Notes"}); !strings.Contains(why, `"Customer name"`) {
		t.Errorf("the words: %q", why)
	}
}
