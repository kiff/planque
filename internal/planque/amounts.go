package planque

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strings"
)

// An API answers money in its currency's smallest unit ("amount": 700,
// "currency": "eur" is €7.00), and an agent reading 700 has to divide in its
// head: once, a taf told its owner "a €700 refund" for 700 cents. So for a
// vendor whose own API says every amount beside a currency is in that
// currency's smallest unit, planque writes each such amount in the currency's
// own unit beside it, by that vendor's decimals: "amount_major": "7.00 EUR".
// The agent says that, and never works the figure out. Nothing is guessed:
// a vendor that doesn't say it (an amount can be major there, "7.00" or
// 7.5), a currency whose decimals aren't known, or a field already there gets
// nothing, and such an answer is passed on as it was.

// Decimals is a vendor's number of decimals for a currency in its API, and
// whether it's known.
type Decimals func(currency string) (int, bool)

// minorAmountVendors are the vendors whose own API says every integer amount
// beside a currency is in that currency's smallest unit, by their catalog
// entry.
var minorAmountVendors = map[string]Decimals{
	// Stripe: "All API requests expect amounts to be provided in a currency's
	// smallest unit" (its docs, Currencies › Zero-decimal currencies).
	StripeEntry: stripeDecimals,
}

// Stripe's own decimals: its zero-decimal and three-decimal lists, and the
// currencies it treats as special cases (Icelandic króna, Hungarian forint,
// New Taiwan dollar, Ugandan shilling), which get no major amount rather than
// a guess. Every other currency on ISO 4217's two-decimal list is two.
var (
	stripeZero    = []string{"BIF", "CLP", "DJF", "GNF", "JPY", "KMF", "KRW", "MGA", "PYG", "RWF", "VND", "VUV", "XAF", "XOF", "XPF"}
	stripeThree   = []string{"BHD", "JOD", "KWD", "OMR", "TND"}
	stripeSpecial = []string{"ISK", "HUF", "TWD", "UGX"}
)

func stripeDecimals(currency string) (int, bool) {
	code := strings.ToUpper(strings.TrimSpace(currency))
	switch {
	case has(stripeSpecial, code):
		return 0, false
	case has(stripeZero, code):
		return 0, true
	case has(stripeThree, code):
		return 3, true
	case has(isoHundredths, code):
		return 2, true
	}
	return 0, false
}

// isoHundredths are the active ISO 4217 currencies at 2 decimals.
var isoHundredths = strings.Fields(`AED AFN ALL AMD ANG AOA ARS AUD AWG AZN BAM BBD BDT BGN BMD BND BOB BOV BRL
BSD BTN BWP BYN BZD CAD CDF CHE CHF CHW CNY COP COU CRC CUP CVE CZK DKK DOP DZD EGP ERN ETB EUR FJD FKP GBP
GEL GHS GIP GMD GTQ GYD HKD HNL HTG HUF IDR ILS INR IRR JMD KES KGS KHR KPW KYD KZT LAK LBP LKR LRD LSL MAD
MDL MKD MMK MNT MOP MRU MUR MVR MWK MXN MXV MYR MZN NAD NGN NIO NOK NPR NZD PAB PEN PGK PHP PKR PLN QAR
RON RSD RUB SAR SBD SCR SDG SEK SGD SHP SLE SOS SRD SSP STN SVC SYP SZL THB TJS TMT TOP TRY TTD TWD TZS UAH
USD USN UYU UZS VES WST XCD YER ZAR ZMW ZWG`)

func has(list []string, code string) bool {
	for _, c := range list {
		if c == code {
			return true
		}
	}
	return false
}

// MajorAmount is minor (in the currency's smallest unit) in its own unit, at
// decimals, with the currency's code: "7.00 EUR", "700 JPY", "0.700 KWD".
func MajorAmount(minor *big.Int, currency string, decimals int) string {
	code := strings.ToUpper(strings.TrimSpace(currency))
	sign, digits := "", new(big.Int).Abs(minor).String()
	if minor.Sign() < 0 {
		sign = "-"
	}
	if decimals == 0 {
		return sign + digits + " " + code
	}
	for len(digits) <= decimals {
		digits = "0" + digits
	}
	return sign + digits[:len(digits)-decimals] + "." + digits[len(digits)-decimals:] + " " + code
}

// isAmountField is whether an object's field holds an amount, in a vendor
// that says its amounts are minor: "amount", a name with "amount" as one of
// its words (amount_refunded, unit_amount), or a total (total, subtotal).
// Never one planque wrote itself.
func isAmountField(name string) bool {
	n := strings.ToLower(name)
	if strings.HasSuffix(n, "_major") {
		return false
	}
	switch n {
	case "amount", "total", "subtotal":
		return true
	}
	for _, w := range strings.FieldsFunc(n, func(r rune) bool { return r == '_' || r == '-' }) {
		if w == "amount" {
			return true
		}
	}
	return false
}

// MaxAmountsDepth bounds how deep an answer is read for amounts.
const MaxAmountsDepth = 32

// WithMajorAmounts is an answer's text with each integer amount beside a
// currency also in its own unit (<field>_major), at any depth, by the
// vendor's decimals. Text that isn't one JSON object or array, or has no such
// amount, comes back as it was.
func WithMajorAmounts(text string, decimals Decimals) string {
	t := bytes.TrimSpace([]byte(text))
	if len(t) == 0 || (t[0] != '{' && t[0] != '[') {
		return text
	}
	dec := json.NewDecoder(bytes.NewReader(t))
	dec.UseNumber()
	var v any
	if dec.Decode(&v) != nil || dec.More() {
		return text
	}
	if !addMajor(v, decimals, 0) {
		return text
	}
	out, err := json.Marshal(v)
	if err != nil {
		return text
	}
	return string(out)
}

// addMajor writes <field>_major beside each amount of an object with a known
// currency, in v and below, replacing one already there that differs, and
// takes one away beside an amount whose currency it can't place; it says
// whether it changed anything.
func addMajor(v any, decimals Decimals, depth int) bool {
	if depth > MaxAmountsDepth {
		return false
	}
	added := false
	switch x := v.(type) {
	case map[string]any:
		if cur, ok := x["currency"].(string); ok {
			places, known := decimals(cur)
			for k, raw := range x {
				n, isNum := raw.(json.Number)
				if !isNum || !isAmountField(k) || strings.ContainsAny(string(n), ".eE") {
					continue
				}
				minor, ok := new(big.Int).SetString(string(n), 10)
				if !ok {
					continue
				}
				// planque's own figure wins: a <field>_major the answer
				// already carries beside the amount and its currency is
				// replaced, and where planque can't place the currency
				// (Stripe's special cases, an unknown one) it's taken
				// away, so the taf never reads a figure planque didn't
				// work out from the amount itself.
				had, has := x[k+"_major"]
				if !known {
					if has {
						delete(x, k+"_major")
						added = true
					}
					continue
				}
				major := MajorAmount(minor, cur, places)
				if s, ok := had.(string); ok && s == major {
					continue
				}
				x[k+"_major"] = major
				added = true
			}
		}
		for _, child := range x {
			added = addMajor(child, decimals, depth+1) || added
		}
	case []any:
		for _, child := range x {
			added = addMajor(child, decimals, depth+1) || added
		}
	}
	return added
}

// withMajorAmounts is a done call's answer with its amounts in their own unit
// too, in each text block and the structured content, for a vendor that says
// its amounts are minor; any other answer is returned as it was.
func withMajorAmounts(g *Grant, res any) any {
	decimals, ok := minorAmountVendors[g.Request.CatalogEntry]
	if !ok {
		return res
	}
	var m map[string]any
	raw, err := json.Marshal(res)
	if err != nil || decodeNumbers(raw, &m) != nil || m == nil {
		return res
	}
	changed := false
	if blocks, ok := m["content"].([]any); ok {
		for _, b := range blocks {
			block, ok := b.(map[string]any)
			if !ok || block["type"] != "text" {
				continue
			}
			text, _ := block["text"].(string)
			if out := WithMajorAmounts(text, decimals); out != text {
				block["text"] = out
				changed = true
			}
		}
	}
	if sc, ok := m["structuredContent"]; ok && addMajor(sc, decimals, 0) {
		changed = true
	}
	if !changed {
		return res
	}
	return m
}
