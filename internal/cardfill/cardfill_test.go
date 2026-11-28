package cardfill

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestTheModeIsOffOrTestOnly(t *testing.T) {
	for v, want := range map[string]Mode{"": Off, "off": Off, "test": Test, " test ": Test} {
		if m, err := ModeFrom(v); err != nil || m != want {
			t.Errorf("%q: %v %v", v, m, err)
		}
	}
	for _, v := range []string{"live", "on", "prod", "TEST"} {
		if _, err := ModeFrom(v); err == nil {
			t.Errorf("%q accepted", v)
		}
	}
}

// In test mode only Stripe's test cards are fillable; when off, none is;
// and a refusal never carries the number.
func TestOnlyStripesTestCardsAreFilled(t *testing.T) {
	if err := Fillable(Test, "4242 4242 4242 4242"); err != nil {
		t.Fatal(err)
	}
	if err := Fillable(Off, "4242424242424242"); err == nil {
		t.Fatal("filled while off")
	}
	err := Fillable(Test, "4111111111111111")
	if !errors.Is(err, ErrNotATestCard) || strings.Contains(err.Error(), "4111") {
		t.Fatalf("a non-test card: %v", err)
	}
}

func TestEgressIsTheMerchantAndThePaymentHostsOnly(t *testing.T) {
	got, err := Egress("Papeterie-Lumiere.example")
	if err != nil || got[0] != "papeterie-lumiere.example" || len(got) != 1+len(PaymentHosts()) {
		t.Fatalf("%v %v", got, err)
	}
	for _, bad := range []string{"", "https://shop.example", "10.0.0.1", "*.example", "shop.example/pay", "localhost"} {
		if _, err := Egress(bad); err == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}

// The canary finds a card number however it's written, and a card's last
// eight, and never reports the digits; an order number or a phone number
// doesn't trip it.
func TestTheCanaryFindsCardNumbersAndNothingElse(t *testing.T) {
	for _, text := range []string{"paid with 4242424242424242", "card 4242 4242 4242 4242 ok", "4242-4242-4242-4242"} {
		hits := Canary("log", text)
		if len(hits) == 0 || hits[0].What != "card number" || strings.Contains(hits[0].What+hits[0].Where, "4242") {
			t.Errorf("%q: %+v", text, hits)
		}
	}
	if hits := Canary("log", "ref ...42424242 end", "4242424242424242"); len(hits) != 1 || hits[0].What != "a card's last eight" {
		t.Errorf("last eight: %+v", hits)
	}
	for _, clean := range []string{"Order #10234567", "Call +33 1 23 45 67 89", "Total 84.20 EUR", "card ending 4242"} {
		if hits := Canary("log", clean, "4242424242424242"); len(hits) != 0 {
			t.Errorf("%q tripped: %+v", clean, hits)
		}
	}
}

// A result with a card number in it is dropped whole and becomes unknown;
// an outcome that isn't one of the three is unknown too.
func TestAResultWithACardInItIsDropped(t *testing.T) {
	ok := Result{Outcome: "ordered", OrderRef: "A-10234", Confirmation: "Thank you, order A-10234"}
	if got, err := ok.Clean("4242424242424242"); err != nil || got != ok {
		t.Fatalf("%+v %v", got, err)
	}
	bad := Result{Outcome: "ordered", OrderRef: "A-1", Confirmation: "Paid with 4242 4242 4242 4242"}
	if got, err := bad.Clean(); !errors.Is(err, ErrCardInResult) || got != (Result{Outcome: "unknown"}) {
		t.Fatalf("%+v %v", got, err)
	}
	if got, err := (Result{Outcome: "done"}).Clean(); err == nil || got.Outcome != "unknown" {
		t.Fatalf("%+v %v", got, err)
	}
}

// The delivery address: the page's address fields when it has them, its
// text otherwise, folding case, accents and spacing; an empty case refused.
func TestAddressShown(t *testing.T) {
	a := Address{Name: "Acme SAS", Line1: "12 rue des Lilas", PostalCode: "69003", City: "Lyon", Country: "FR"}
	ok := map[string]struct {
		fields map[string]string
		text   string
	}{
		"text":               {nil, "Deliver to: Acme SAS, 12 RUE DES LILAS, 69003 Lyon, France"},
		"text, spaced":       {nil, "12,  rue des Lilas\n69003 — Lyon"},
		"fields":             {map[string]string{"address-line1": "12 rue des Lilas", "postal-code": "69003", "address-level2": "LYON", "country": "fr"}, ""},
		"fields, no country": {map[string]string{"address-line1": "12 rue des lilas", "postal-code": "69003", "address-level2": "Lyon"}, ""},
	}
	for name, c := range ok {
		if err := AddressShown(a, c.fields, c.text); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	accented := Address{Line1: "3 place de l'Hôtel-de-Ville", PostalCode: "38000", City: "Grenoble"}
	if err := AddressShown(accented, nil, "3 Place de l Hotel de Ville, 38000 GRENOBLE"); err != nil {
		t.Errorf("accents and punctuation: %v", err)
	}
	bad := map[string]struct {
		a      Address
		fields map[string]string
		text   string
	}{
		"another city":           {a, nil, "12 rue des Lilas, 69003 Villeurbanne"},
		"another postal code":    {a, nil, "12 rue des Lilas, 69007 Lyon"},
		"a field differs":        {a, map[string]string{"address-line1": "1 avenue Foch", "postal-code": "69003", "address-level2": "Lyon"}, "12 rue des Lilas, 69003 Lyon"},
		"a field missing":        {a, map[string]string{"postal-code": "69003", "address-level2": "Lyon"}, "12 rue des Lilas, 69003 Lyon"},
		"another country":        {a, map[string]string{"address-line1": "12 rue des Lilas", "postal-code": "69003", "address-level2": "Lyon", "country": "DE"}, ""},
		"a case with no address": {Address{}, nil, "anything"},
	}
	for name, c := range bad {
		if err := AddressShown(c.a, c.fields, c.text); !errors.Is(err, ErrAddressNotShown) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// A Card printed by mistake shows no digit.
func TestCardPrintsNoDigit(t *testing.T) {
	c := Card{Number: "4242424242424242", ExpMonth: 12, ExpYear: 2034, CVC: "123"}
	for _, s := range []string{fmt.Sprint(c), fmt.Sprintf("%v %+v %#v %s", c, c, c, c)} {
		if strings.ContainsAny(s, "0123456789") {
			t.Errorf("printed: %q", s)
		}
	}
}

// A card number written with dots is found too.
func TestTheCanaryFindsDottedNumbers(t *testing.T) {
	if hits := Canary("log", "paid with 4242.4242.4242.4242 today"); len(hits) != 1 || hits[0].What != "card number" {
		t.Errorf("dotted: %+v", hits)
	}
	if hits := Canary("log", "total 1.234.56 EUR, order 2026.11.28"); len(hits) != 0 {
		t.Errorf("not a card: %+v", hits)
	}
}
