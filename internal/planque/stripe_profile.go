package planque

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/egress"
)

// The seller's Stripe profile, for a purchase paid by shared payment token
// (Thune v1, path a): read from the merchant's own Stripe-backed
// declaration on the checkout's own site, never from the taf. A seller that
// takes Stripe's agentic payments publishes its UCP profile at
// https://<site>/.well-known/ucp, with the com.stripe.payments handler's
// config naming its business_profile and environment. planque reads it from
// the host the checkout page is on, over https, following no redirect,
// and reports the profile with that host (stripe_profile_from), so KIFF signs
// it into the case only when the host is the merchant's site. A profile in
// live mode isn't taken while tokens are issued in test mode; two that
// differ give none. On with PLANQUE_STRIPE_PROFILE=on.

// stripeProfileID is a Stripe business network profile's id.
var stripeProfileID = regexp.MustCompile(`^profile_[A-Za-z0-9_]{1,200}$`)

// ucpProfileMax bounds the profile read.
const ucpProfileMax = 64 << 10

// ucpGet reads https://host/.well-known/ucp: public addresses only, no
// redirect, at most ucpProfileMax. Tests replace it (Server.UCPGet).
func ucpGet(ctx context.Context, host string) ([]byte, error) {
	c := egress.Client(5 * time.Second)
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+host+"/.well-known/ucp", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("no UCP profile at the site")
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, ucpProfileMax+1))
	if err != nil || len(b) > ucpProfileMax {
		return nil, errors.New("the site's UCP profile is too large")
	}
	return b, nil
}

// stripeProfileOf is the test-mode Stripe profile a UCP profile names in its
// com.stripe.payments handlers, or "": none, one in live mode, or two that
// differ.
func stripeProfileOf(raw []byte) string {
	var doc struct {
		UCP struct {
			Handlers map[string][]struct {
				Config struct {
					Environment     string `json:"environment"`
					BusinessProfile string `json:"business_profile"`
				} `json:"config"`
			} `json:"payment_handlers"`
		} `json:"ucp"`
	}
	if json.Unmarshal(raw, &doc) != nil {
		return ""
	}
	found := ""
	for _, h := range doc.UCP.Handlers["com.stripe.payments"] {
		p := h.Config.BusinessProfile
		if h.Config.Environment != "testmode" || !stripeProfileID.MatchString(p) {
			continue
		}
		if found != "" && found != p {
			return ""
		}
		found = p
	}
	return found
}

// readStripeProfile is the seller's Stripe profile for the checkout at
// pageURL, and the host it was read on; both empty when there's none or
// the read is off.
func (s *Server) readStripeProfile(ctx context.Context, pageURL string) (profile, from string) {
	if !s.StripeProfiles {
		return "", ""
	}
	u, err := url.Parse(pageURL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.Port() != "" {
		return "", ""
	}
	host := strings.ToLower(u.Hostname())
	get := s.UCPGet
	if get == nil {
		get = ucpGet
	}
	raw, err := get(ctx, host)
	if err != nil {
		return "", ""
	}
	if p := stripeProfileOf(raw); p != "" {
		return p, host
	}
	return "", ""
}

// withStripeProfile sets co's seller profile, read for the page at pageURL.
func (s *Server) withStripeProfile(ctx context.Context, co *browser.Checkout, pageURL string) {
	co.StripeProfile, co.StripeProfileFrom = s.readStripeProfile(ctx, pageURL)
}
