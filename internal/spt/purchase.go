package spt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Money is an amount in its currency's smallest unit.
type Money struct {
	Minor    int64  `json:"minor"`
	Currency string `json:"currency"`
}

// Seller is the shop as Stripe names it: its Stripe profile and business
// name, read from Stripe's verified checkout or profile, never from a page.
type Seller struct {
	Profile string `json:"profile"`
	Name    string `json:"name"`
	// Verified says the seller came from Stripe's own answer. A seller from
	// anything else (page text, the taf's words) is never paid.
	Verified bool `json:"verified"`
}

// Line is one line of the cart.
type Line struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Amount   Money  `json:"amount"`
}

// Checkout is the shop's checkout, as read from the shop's own checkout
// session: the seller, the lines and the total. The taf names none of it.
type Checkout struct {
	Seller Seller `json:"seller"`
	Lines  []Line `json:"lines"`
	Total  Money  `json:"total"`
}

// Key is the checkout's mark: the seller, the lines and the total, hashed.
// Two approvals for the same cart carry the same key, so the cart is never
// paid twice.
func (c Checkout) Key() string {
	b, _ := json.Marshal(struct {
		Seller string `json:"seller"`
		Lines  []Line `json:"lines"`
		Total  Money  `json:"total"`
	}{c.Seller.Profile, c.Lines, Money{c.Total.Minor, strings.ToUpper(c.Total.Currency)}})
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Approval is what the owner approved with the phone's key and KIFF's
// permit carries: the purchase is paid only on these exact fields.
type Approval struct {
	Tenant        string
	Taf           string
	Purchase      string
	Permit        string
	SellerProfile string
	Amount        Money
	ExpiresAt     time.Time
}

// Card is the payment method the account saved with Stripe: its id and what
// may be shown of it, bound to one account.
type Card struct {
	Tenant        string
	PaymentMethod string
	Brand         string
	Last4         string
}

// Store keeps what the Issuer needs between calls.
type Store interface {
	// CardFor is the account's saved card, if any.
	CardFor(ctx context.Context, tenant string) (Card, bool, error)
	// MarkCheckout marks the checkout key for the account and reports false
	// when it was marked in the last 24 hours: a purchase of that cart
	// already ran, or its outcome isn't known.
	MarkCheckout(ctx context.Context, tenant, key string, now time.Time) (bool, error)
}

// Purchase is one paid-for purchase. Its token id stays in planque.
type Purchase struct {
	ID       string `json:"id"`
	Tenant   string `json:"tenant"`
	Taf      string `json:"taf"`
	Seller   Seller `json:"seller"`
	Amount   Money  `json:"amount"`
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	Limits   Limits `json:"limits"`
	State    string `json:"state"`
	Checkout string `json:"checkout"`
	Token    string `json:"-"`
}

// Purchase states.
const (
	StateIssued   = "issued"
	StatePaid     = "paid"
	StateNotPaid  = "not_paid"
	StateIncident = "incident"
)

// MaxLife is the longest a token lives: the window in which a halt can race
// a payment.
const MaxLife = 10 * time.Minute

// The Issuer's refusals.
var (
	ErrOff              = errors.New("spt: paying with a token is off")
	ErrTestTenant       = errors.New("spt: a proof test account never pays")
	ErrSellerUnverified = errors.New("spt: the seller isn't Stripe's word")
	ErrMismatch         = errors.New("spt: the approval isn't this checkout's seller and total")
	ErrCurrency         = errors.New("spt: the currency isn't one Stripe's decimals are known for")
	ErrAmount           = errors.New("spt: the amount must be more than zero")
	ErrExpiry           = errors.New("spt: the approval's expiry must be in the next 10 minutes")
	ErrNoCard           = errors.New("spt: the account has no saved card")
	ErrWrongTenant      = errors.New("spt: the card belongs to another account")
	ErrRepeat           = errors.New("spt: this cart was paid for, or its outcome isn't known, in the last 24 hours")
)

// Issuer issues a token for one approved purchase.
type Issuer struct {
	// On turns it on; off by default.
	On     bool
	Client *Client
	Store  Store
	// Decimals is Stripe's decimals for a currency, and whether they're
	// known (planque's stripeDecimals): an unknown one is never paid.
	Decimals func(currency string) (int, bool)
	// TestTenants are the proof test accounts, which never pay.
	TestTenants map[string]bool
	Now         func() time.Time
}

func (i *Issuer) now() time.Time {
	if i.Now != nil {
		return i.Now()
	}
	return time.Now()
}

// Pay issues the token for a, after every check, and returns the purchase
// with the token kept inside it. Handing the token to the shop's checkout is
// the caller's next step, through planque's guarded routes only.
func (i *Issuer) Pay(ctx context.Context, a Approval, c Checkout) (Purchase, error) {
	switch {
	case !i.On:
		return Purchase{}, ErrOff
	case i.TestTenants[a.Tenant]:
		return Purchase{}, ErrTestTenant
	case !c.Seller.Verified || c.Seller.Profile == "":
		return Purchase{}, ErrSellerUnverified
	case a.SellerProfile != c.Seller.Profile || a.Amount.Minor != c.Total.Minor || !strings.EqualFold(a.Amount.Currency, c.Total.Currency):
		return Purchase{}, ErrMismatch
	case a.Amount.Minor <= 0:
		return Purchase{}, ErrAmount
	}
	if _, ok := i.Decimals(a.Amount.Currency); !ok {
		return Purchase{}, ErrCurrency
	}
	now := i.now()
	if !a.ExpiresAt.After(now) || a.ExpiresAt.Sub(now) > MaxLife {
		return Purchase{}, ErrExpiry
	}
	card, ok, err := i.Store.CardFor(ctx, a.Tenant)
	if err != nil {
		return Purchase{}, err
	}
	if !ok {
		return Purchase{}, ErrNoCard
	}
	if card.Tenant != a.Tenant {
		return Purchase{}, ErrWrongTenant
	}
	key := c.Key()
	if fresh, err := i.Store.MarkCheckout(ctx, a.Tenant, key, now); err != nil {
		return Purchase{}, err
	} else if !fresh {
		return Purchase{}, ErrRepeat
	}
	t, err := i.Client.Issue(ctx, Issue{
		PaymentMethod:  card.PaymentMethod,
		SellerProfile:  c.Seller.Profile,
		Currency:       a.Amount.Currency,
		MaxAmount:      a.Amount.Minor,
		ExpiresAt:      a.ExpiresAt,
		IdempotencyKey: a.Permit,
		Purchase:       a.Purchase,
	})
	if err != nil {
		return Purchase{}, err
	}
	return Purchase{ID: a.Purchase, Tenant: a.Tenant, Taf: a.Taf, Seller: c.Seller, Amount: a.Amount,
		Brand: card.Brand, Last4: card.Last4, Limits: t.Limits, State: StateIssued, Checkout: key, Token: t.ID}, nil
}

// Settle reads the token's state from Stripe, which decides: a webhook only
// says when to read. Consumed is paid; expired or revoked is not paid; a
// token resolved to another processor, which this flow never asks for, is an
// incident; active stays issued.
func (i *Issuer) Settle(ctx context.Context, p Purchase) (string, error) {
	t, err := i.Client.Get(ctx, p.Token)
	if err != nil {
		return p.State, err
	}
	if t.Limits != p.Limits {
		return StateIncident, nil
	}
	switch {
	case t.Status == "active" || t.Status == "requires_action":
		return StateIssued, nil
	case t.Status == "deactivated" && t.DeactivatedReason == "consumed":
		return StatePaid, nil
	case t.Status == "deactivated" && (t.DeactivatedReason == "expired" || t.DeactivatedReason == "revoked"):
		return StateNotPaid, nil
	}
	return StateIncident, nil
}

// Stop revokes the tokens of purchases still issued: a halt, a Card change
// or the taf's deletion. A token used before its revoke reads paid at its
// next Settle.
func (i *Issuer) Stop(ctx context.Context, ps []Purchase) error {
	var errs []error
	for _, p := range ps {
		if p.State == StateIssued && p.Token != "" {
			if err := i.Client.Revoke(ctx, p.Token); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// String is the purchase without its token id, so a log never carries it.
func (p Purchase) String() string {
	return "purchase{" + p.ID + " " + p.State + " " + p.Seller.Name + " " + strconv.FormatInt(p.Amount.Minor, 10) + " " + p.Amount.Currency + "}"
}

// GoString is String, so %#v masks the token too.
func (p Purchase) GoString() string { return p.String() }
