package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/stepkey"
)

// A purchase case from a browse_pay (Thune v1): before KIFF decides a pay,
// its gateway asks planque to describe it, and builds the case only from
// what planque answers, never from the taf's arguments. The describe is the
// pay itself with the gateway's header KIFF-Checkout-Describe naming the
// pay's operation, under KIFF's permit for "<op>:describe". planque answers
// the checkout it read and kept at browse_checkout (lines, total, currency,
// shipping, tax, merchant, delivery address, and the seller's Stripe profile
// with the host it was read from), signed with its describe key under its own
// label. No browser starts and nothing reaches the site; the taf is told
// nothing. At the pay itself planque still checks the page against the same
// record, so what runs is what was described.

// CheckoutDescribeHeader asks a browse_pay to describe its checkout: the
// pay's operation. A taf can't set it: the gateway sends only its own.
const CheckoutDescribeHeader = "KIFF-Checkout-Describe"

// labelCheckout is what planque signs a checkout's block under: not the
// sealed fill's label, so neither block passes for the other.
const labelCheckout = "planque/checkout/v1"

var describeOp = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)

type checkoutHeaderKey struct{}

// withCheckoutHeader keeps a call's checkout-describe header, or "bad"
// when it doesn't read (the call is then refused).
func withCheckoutHeader(ctx context.Context, h http.Header) context.Context {
	v := strings.TrimSpace(h.Get(CheckoutDescribeHeader))
	if v == "" {
		return ctx
	}
	if !describeOp.MatchString(v) {
		v = "\x00bad"
	}
	return context.WithValue(ctx, checkoutHeaderKey{}, v)
}

func checkoutHeaderFrom(ctx context.Context) string {
	v, _ := ctx.Value(checkoutHeaderKey{}).(string)
	return v
}

// checkoutBlock is what planque signs for a pay's checkout, as KIFF builds
// a purchase case from it (FromCheckout).
type checkoutBlock struct {
	V                 string            `json:"v"`
	Tenant            string            `json:"tenant"`
	Op                string            `json:"op"`
	Grant             string            `json:"grant"`
	Site              string            `json:"site"`
	URL               string            `json:"url"`
	Ref               string            `json:"ref"`
	Lines             []browser.Line    `json:"lines"`
	Total             string            `json:"total"`
	Currency          string            `json:"currency"`
	Shipping          string            `json:"shipping,omitempty"`
	Tax               string            `json:"tax,omitempty"`
	Merchant          *browser.Merchant `json:"merchant,omitempty"`
	Delivery          *browser.Address  `json:"delivery,omitempty"`
	StripeProfile     string            `json:"stripe_profile,omitempty"`
	StripeProfileFrom string            `json:"stripe_profile_from,omitempty"`
	ReadAt            time.Time         `json:"read_at"`
}

// describeCheckout answers the gateway's describe of the pay a: the kept
// checkout read for a.URL, signed. A pay with no kept read, one past its
// time, or one whose earlier pay ended unknown, gets no block.
func (s *Server) describeCheckout(ctx context.Context, g *Grant, p permit.Permit, a browserArgs, op string) (string, any) {
	switch {
	case op == "\x00bad":
		return OpNotSent, browserRefusal(CodeBadArguments, "the checkout describe header doesn't read", false)
	case a.CheckoutDescribe != op:
		return OpNotSent, browserRefusal(CodeBadArguments, "a pay's describe carries checkout_describe, the pay's op, as its header does", false)
	case p.Op != op+":describe":
		return OpNotSent, browserRefusal(CodePermitInvalid, "the describe's permit isn't for this pay's describe", false)
	case s.ChromeSigner == nil:
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "planque can't sign a checkout's description here", false)
	case a.Ref == "" || len(a.Steps) > 0:
		return OpNotSent, browserRefusal(CodeBadArguments, "a pay's describe names its url and the control that pays, and nothing else", false)
	}
	var rec checkoutRecord
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skCheckout(g.ID, a.URL), &rec); errors.Is(err, errNotFound) {
		return OpNotSent, browserRefusal(browser.CodeTotalUnread, "no total was read on this checkout: read it with "+checkoutName(g)+" first", false)
	} else if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read the checkout's record", false)
	}
	switch {
	case !s.now().Before(rec.ReadAt.Add(checkoutKept)):
		return OpNotSent, browserRefusal(browser.CodeTotalUnread, "the checkout was read too long ago: read it again first", false)
	case !rec.UnknownAt.IsZero():
		return OpNotSent, browserRefusal(browser.CodeOutcomeUnknown, "an earlier pay on this checkout ended without a known outcome, so it isn't described again", false)
	}
	co := rec.Checkout
	b := checkoutBlock{V: labelCheckout, Tenant: g.Tenant, Op: op, Grant: g.ID, Site: rec.Site, URL: a.URL, Ref: a.Ref,
		Lines: co.Lines, Total: co.Total, Currency: co.Currency, Shipping: co.Shipping, Tax: co.Tax, Merchant: co.Merchant,
		Delivery: co.Delivery, StripeProfile: co.StripeProfile, StripeProfileFrom: co.StripeProfileFrom, ReadAt: rec.ReadAt.UTC()}
	if b.Lines == nil {
		b.Lines = []browser.Line{}
	}
	sig, err := stepkey.SignBody(ctx, s.ChromeSigner, labelCheckout, b)
	if err != nil {
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "planque could not sign the checkout's description", false)
	}
	s.log().InfoContext(ctx, "checkout described", slog.String("grant", g.ID), slog.String("op", op), slog.String("site", rec.Site))
	block, _ := json.Marshal(b)
	res := toolResult("Described for the owner's purchase.", false)
	res["structuredContent"] = map[string]any{"browser_block": map[string]any{"block": json.RawMessage(block), "sig": sig}}
	return OpDone, res
}
