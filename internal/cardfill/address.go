package cardfill

import (
	"errors"
	"strings"
	"unicode"
)

// Card is the one card a payment fills: an issued card for that one case,
// or, in test mode, one of Stripe's test cards. It reaches the worker
// alone; planque's API and KIFF never hold it. Its String never shows the
// number or the code.
type Card struct {
	Number   string
	ExpMonth int
	ExpYear  int // four digits
	CVC      string
	Name     string
}

// String is the card as a log may show it: no digit at all, so a Card
// printed by mistake leaks nothing.
func (Card) String() string { return "card(hidden)" }

// GoString hides the card from %#v as String does from %v.
func (Card) GoString() string { return "card(hidden)" }

// Address is the delivery address the case signs (contracts/purchase.json,
// case.delivery).
type Address struct {
	Name       string `json:"name"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Region     string `json:"region"`
	Country    string `json:"country"`
}

// ErrAddressNotShown refuses a checkout whose delivery address isn't the
// case's: nothing is filled or clicked.
var ErrAddressNotShown = errors.New("the checkout's delivery address isn't the one the case signed: nothing was filled or paid")

// AddressShown says the checkout delivers to a, as the page shows it. The
// page's address fields (by their autocomplete token: address-line1,
// postal-code, address-level2, country) are compared when the page has
// them; otherwise the page's text must hold the first line, the postal
// code and the city. Either way the comparison folds case, accents and
// spacing, and an empty line1, postal code or city in the case is refused:
// a case without an address delivers nowhere.
func AddressShown(a Address, fields map[string]string, text string) error {
	if fold(a.Line1) == "" || fold(a.PostalCode) == "" || fold(a.City) == "" {
		return ErrAddressNotShown
	}
	if len(fields) > 0 {
		want := map[string]string{"address-line1": a.Line1, "postal-code": a.PostalCode, "address-level2": a.City}
		for tok, v := range want {
			got, ok := fields[tok]
			if !ok || fold(got) != fold(v) {
				return ErrAddressNotShown
			}
		}
		if c, ok := fields["country"]; ok && a.Country != "" && fold(c) != fold(a.Country) {
			return ErrAddressNotShown
		}
		return nil
	}
	page := fold(text)
	for _, v := range []string{a.Line1, a.PostalCode, a.City} {
		if !strings.Contains(page, fold(v)) {
			return ErrAddressNotShown
		}
	}
	return nil
}

// unaccented maps the accented Latin letters of European addresses to
// their base letter (standard library only, so no Unicode normalizer).
var unaccented = func() map[rune]rune {
	m := map[rune]rune{}
	for base, accented := range map[rune]string{
		'a': "àáâãäåāăą", 'c': "çćĉċč", 'd': "ďđ", 'e': "èéêëēĕėęě", 'g': "ĝğġģ",
		'h': "ĥħ", 'i': "ìíîïĩīĭįı", 'j': "ĵ", 'k': "ķ", 'l': "ĺļľŀł", 'n': "ñńņňŉ",
		'o': "òóôõöøōŏő", 'r': "ŕŗř", 's': "śŝşšș", 't': "ţťŧț", 'u': "ùúûüũūŭůűų",
		'w': "ŵ", 'y': "ýÿŷ", 'z': "źżž",
	} {
		for _, r := range accented {
			m[r] = base
		}
	}
	return m
}()

// fold is s compared loosely: lower case, accents dropped, runs of space
// and punctuation as one space.
func fold(s string) string {
	var b strings.Builder
	space := true
	for _, r := range strings.ToLower(s) {
		if base, ok := unaccented[r]; ok {
			r = base
		}
		switch {
		case unicode.Is(unicode.Mn, r):
			continue
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			space = false
		default:
			if !space {
				b.WriteByte(' ')
				space = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}
