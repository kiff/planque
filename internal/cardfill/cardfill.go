// Package cardfill holds the payment-only worker's guard rails (Thune v1,
// batch B2, path (b)): what it may reach, what it may return, how a card
// number is found anywhere it must not be, and the mode it runs in. The
// worker itself fills a card into one checkout on a fresh machine and is
// destroyed after; nothing here types a card.
//
// Path (b) runs in test mode only until the card-data (PCI) decision. In
// test mode only Stripe's published test card numbers may be filled, so a
// real card can't be typed even if one reached the worker.
package cardfill

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Mode is how the worker may run: Off (the default) or Test. There is no
// live mode in code: adding one is a reviewed change, after the PCI
// decision.
type Mode string

const (
	Off  Mode = ""
	Test Mode = "test"
)

// ModeFrom reads PLANQUE_CARD_FILL: unset is Off, "test" is Test, and
// anything else is refused, so a typo never turns the worker on in another
// mode.
func ModeFrom(v string) (Mode, error) {
	switch strings.TrimSpace(v) {
	case "", "off":
		return Off, nil
	case "test":
		return Test, nil
	}
	return Off, fmt.Errorf("PLANQUE_CARD_FILL is %q: only test (or unset) is accepted; there is no live card fill", v)
}

// testCards are Stripe's published test card numbers (docs.stripe.com/testing):
// the only numbers the worker fills in test mode.
var testCards = map[string]bool{
	"4242424242424242": true, // Visa
	"4000056655665556": true, // Visa (debit)
	"5555555555554444": true, // Mastercard
	"2223003122003222": true, // Mastercard (2-series)
	"5200828282828210": true, // Mastercard (debit)
	"378282246310005":  true, // American Express
	"4000002500003155": true, // Visa, 3-D Secure required (the decline-a-challenge test)
	"4000000000000002": true, // Visa, declined
}

// ErrNotATestCard refuses a number in test mode that isn't one of Stripe's
// test cards. The number itself is never in the error.
var ErrNotATestCard = errors.New("card fill runs in test mode: only Stripe's published test card numbers may be filled")

// Fillable says whether number may be filled in mode: never when Off, and
// in Test only one of Stripe's test cards.
func Fillable(m Mode, number string) error {
	if m != Test {
		return errors.New("card fill is off")
	}
	if !testCards[digitsOf(number)] {
		return ErrNotATestCard
	}
	return nil
}

// paymentHosts are the payment services a checkout's card form talks to,
// beside the merchant's own site: the published list, kept small. A
// journey reaches the merchant's site and these, nothing else. 3-D Secure
// hosts aren't on it: a challenge on an issued card declines the purchase
// (the threat model's T9), it's never answered.
var paymentHosts = []string{
	"js.stripe.com", "m.stripe.network", "api.stripe.com", "checkout.stripe.com",
}

// PaymentHosts is a copy of the published list.
func PaymentHosts() []string { return append([]string(nil), paymentHosts...) }

var hostName = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)

// Egress is the sites one payment's browser may reach: the merchant's own
// site (a host name, as the case names it) and the published payment
// hosts. A site that isn't a plain host name is refused.
func Egress(merchantSite string) ([]string, error) {
	site := strings.ToLower(strings.TrimSpace(merchantSite))
	if !hostName.MatchString(site) {
		return nil, fmt.Errorf("the merchant's site %q isn't a host name", merchantSite)
	}
	return append([]string{site}, paymentHosts...), nil
}

// Hit is a place the canary found a card number where none may be.
type Hit struct {
	Where string // what was scanned: "returned text", "log", "capture"…
	What  string // "card number" or "a card's last eight"; never the digits
}

var digitRun = regexp.MustCompile(`(?:\d[ .-]?){12,18}\d`)

// Canary scans text, named where, for any Luhn-valid run of 13 to 19
// digits (spaces, dashes or dots allowed, as a card number is written) and for
// the last eight digits of any of cards. It returns what it found, never
// the digits.
func Canary(where, text string, cards ...string) []Hit {
	var hits []Hit
	for _, run := range digitRun.FindAllString(text, -1) {
		d := digitsOf(run)
		if len(d) >= 13 && len(d) <= 19 && luhn(d) {
			hits = append(hits, Hit{Where: where, What: "card number"})
		}
	}
	flat := digitsOf(text)
	for _, c := range cards {
		if d := digitsOf(c); len(d) >= 8 && strings.Contains(flat, d[len(d)-8:]) {
			hits = append(hits, Hit{Where: where, What: "a card's last eight"})
		}
	}
	return hits
}

// Result is all the worker returns to planque for one payment: the order
// reference and the confirmation's text, masked, and how it ended. Never a
// field's value.
type Result struct {
	Outcome      string `json:"outcome"` // "ordered", "not_ordered" or "unknown"
	OrderRef     string `json:"order_ref,omitempty"`
	Confirmation string `json:"confirmation,omitempty"`
}

// ErrCardInResult refuses a result the canary found a card number in: the
// result is dropped whole, and the outcome becomes unknown for
// reconciliation to settle.
var ErrCardInResult = errors.New("the worker's result carried a card number: it was dropped, and the outcome is unknown until reconciled")

// Clean is r as planque may keep it: its outcome one of the three, and
// nothing in it the canary finds, for any of cards. A result that fails is
// replaced by an unknown outcome with nothing else.
func (r Result) Clean(cards ...string) (Result, error) {
	switch r.Outcome {
	case "ordered", "not_ordered", "unknown":
	default:
		return Result{Outcome: "unknown"}, fmt.Errorf("the worker's outcome %q isn't ordered, not_ordered or unknown", r.Outcome)
	}
	if len(Canary("order reference", r.OrderRef, cards...))+len(Canary("confirmation", r.Confirmation, cards...)) > 0 {
		return Result{Outcome: "unknown"}, ErrCardInResult
	}
	return r, nil
}

func digitsOf(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

// luhn says digits pass the Luhn check, as every card number does.
func luhn(digits string) bool {
	sum, double := 0, false
	for i := len(digits) - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
