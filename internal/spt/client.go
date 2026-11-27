// Package spt is the agent side of a purchase paid with a Stripe shared
// payment token: planque issues a token for one approved purchase, limited
// to the approved amount, currency, shop and a short expiry; reads its state;
// and revokes it. The taf never holds the token, and never names the amount
// or the shop: both come from the shop's verified checkout (purchase.go).
//
// It is a prototype: nothing in planque calls it yet, its Issuer is off by
// default, and its Client takes only a Stripe test key.
package spt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client calls Stripe's shared payment token routes on the agent's account.
type Client struct {
	HTTP *http.Client
	// Base is Stripe's API address; empty is https://api.stripe.com.
	Base string
	// Version is the Stripe-Version header the preview needs.
	Version string
	// Key gives the secret key at each call. Only a test key is taken.
	Key func(context.Context) (string, error)
}

// ErrLiveKey is a key that isn't a Stripe test key: the prototype never
// touches live mode.
var ErrLiveKey = errors.New("spt: only a Stripe test key (sk_test_ or rk_test_) is taken")

// ErrLiveAnswer is an answer from live mode, which a test key can't give:
// it's refused, never trusted.
var ErrLiveAnswer = errors.New("spt: Stripe answered in live mode")

func testKey(k string) bool {
	return strings.HasPrefix(k, "sk_test_") || strings.HasPrefix(k, "rk_test_")
}

// Limits are a token's usage limits, as Stripe keeps them.
type Limits struct {
	Currency  string `json:"currency"`
	MaxAmount int64  `json:"max_amount"`
	ExpiresAt int64  `json:"expires_at"`
}

// Token is an issued token as Stripe answers it. Its id is a bearer
// credential for the shop it was issued to: it never goes to the taf, a log,
// the Record or the receipt, so it prints masked.
type Token struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	DeactivatedReason string `json:"deactivated_reason"`
	Livemode          bool   `json:"livemode"`
	Limits            Limits `json:"usage_limits"`
	Seller            struct {
		Profile string `json:"network_business_profile"`
	} `json:"seller_details"`
}

// String is the token without its id.
func (t Token) String() string {
	return fmt.Sprintf("token{status:%s reason:%s limits:%d %s until %d}", t.Status, t.DeactivatedReason, t.Limits.MaxAmount, t.Limits.Currency, t.Limits.ExpiresAt)
}

// GoString is String, so %#v masks the id too.
func (t Token) GoString() string { return t.String() }

// Issue is one token to issue.
type Issue struct {
	PaymentMethod string
	SellerProfile string
	Currency      string // ISO 4217, any case
	MaxAmount     int64  // in the currency's smallest unit
	ExpiresAt     time.Time
	// IdempotencyKey is the KIFF permit's id: a retried call never issues a
	// second token.
	IdempotencyKey string
	// Purchase is our purchase id, shared with the shop as metadata.
	Purchase string
}

// Issue issues a token, and checks that Stripe's answer carries exactly the
// limits asked: an answer with any other limit is revoked and refused.
func (c *Client) Issue(ctx context.Context, in Issue) (Token, error) {
	form := url.Values{}
	form.Set("payment_method", in.PaymentMethod)
	form.Set("seller_details[network_business_profile]", in.SellerProfile)
	form.Set("usage_limits[currency]", strings.ToLower(in.Currency))
	form.Set("usage_limits[max_amount]", strconv.FormatInt(in.MaxAmount, 10))
	form.Set("usage_limits[expires_at]", strconv.FormatInt(in.ExpiresAt.Unix(), 10))
	form.Set("shared_metadata[purchase]", in.Purchase)
	var t Token
	if err := c.call(ctx, "POST", "/v1/shared_payment/issued_tokens", form, in.IdempotencyKey, &t); err != nil {
		return Token{}, err
	}
	want := Limits{Currency: strings.ToLower(in.Currency), MaxAmount: in.MaxAmount, ExpiresAt: in.ExpiresAt.Unix()}
	// The seller too: a token issued to any other profile than the one asked
	// is revoked, and an answer that doesn't say which is refused the same.
	if t.Limits != want || t.Seller.Profile != in.SellerProfile || t.Status != "active" || t.ID == "" {
		if t.ID != "" {
			_ = c.Revoke(ctx, t.ID)
		}
		return Token{}, fmt.Errorf("spt: the issued token isn't the one asked (%s)", t)
	}
	return t, nil
}

// Get reads a token's state from Stripe.
func (c *Client) Get(ctx context.Context, id string) (Token, error) {
	var t Token
	err := c.call(ctx, "GET", "/v1/shared_payment/issued_tokens/"+url.PathEscape(id), nil, "", &t)
	return t, err
}

// Revoke stops a token from being used.
func (c *Client) Revoke(ctx context.Context, id string) error {
	var t Token
	return c.call(ctx, "POST", "/v1/shared_payment/issued_tokens/"+url.PathEscape(id)+"/revoke", url.Values{}, "", &t)
}

func (c *Client) call(ctx context.Context, method, path string, form url.Values, idem string, out *Token) error {
	key, err := c.Key(ctx)
	if err != nil {
		return err
	}
	if !testKey(key) {
		return ErrLiveKey
	}
	base := c.Base
	if base == "" {
		base = "https://api.stripe.com"
	}
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, base+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	if c.Version != "" {
		req.Header.Set("Stripe-Version", c.Version)
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	hc := c.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		// A transport error carries the whole URL, and a token's id is in the
		// path of a read or a revoke: say the failure without it.
		var ue *url.Error
		if errors.As(err, &ue) {
			return fmt.Errorf("spt: Stripe unreachable (%s %s): %w", method, routeOf(path), errorOnly{ue.Err})
		}
		return fmt.Errorf("spt: Stripe unreachable (%s %s)", method, routeOf(path))
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	if err != nil {
		return err
	}
	if res.StatusCode/100 != 2 {
		var e struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		return fmt.Errorf("spt: Stripe answered %d %s", res.StatusCode, e.Error.Code)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("spt: Stripe's answer doesn't read: %w", err)
	}
	if out.Livemode {
		return ErrLiveAnswer
	}
	return nil
}

// routeOf is path with a token's id taken out.
func routeOf(path string) string {
	const prefix = "/v1/shared_payment/issued_tokens/"
	rest, ok := strings.CutPrefix(path, prefix)
	if !ok {
		return path
	}
	if _, tail, found := strings.Cut(rest, "/"); found {
		return prefix + "{id}/" + tail
	}
	return prefix + "{id}"
}

// errorOnly wraps a transport's own error, so errors.Is still finds a
// timeout or a cancel, without the URL a *url.Error would print.
type errorOnly struct{ err error }

func (e errorOnly) Error() string { return e.err.Error() }
func (e errorOnly) Unwrap() error { return e.err }
