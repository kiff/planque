// Package issuing makes one single-use virtual card for one purchase case
// KIFF approved, through Stripe Issuing, and cancels it. It is the
// connection only: whether to buy is KIFF's, and planque calls this only
// with KIFF's signature on an approved case (thune_cards.go).
//
// What it never does: read a card's number or code (no reveal, no
// expand=number), answer an authorization, or use a live key. A card's
// spending controls are copied from the case and never widened.
package issuing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MaxWindow is the longest a case may run from the card's creation: a card
// for a case that ends later is refused (the threat model's T2.2).
const MaxWindow = 15 * time.Minute

// Case is the part of a purchase case a card is made from, as KIFF signed
// it. Amounts are in the currency's minor unit.
type Case struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	// State must be "approved": the owner's fingerprint, or KIFF's
	// within-limits decision, on this version.
	State string `json:"state"`
	// MaxAmount is the most the card may be authorized for: the total plus
	// the signed price rule.
	MaxAmount int64  `json:"max_amount"`
	Currency  string `json:"currency"`
	// Categories are Stripe's merchant categories the card may be used at:
	// the case's merchant's.
	Categories []string  `json:"categories"`
	Expires    time.Time `json:"expires_at"`
}

// Card is a card as planque keeps and returns it: never its number or code.
type Card struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Last4    string `json:"last4"`
	Case     string `json:"case"`
	Version  int    `json:"case_version"`
	Permit   string `json:"permit"`
	Currency string `json:"currency"`
	Max      int64  `json:"max_amount"`
}

// Refusals. None sends anything to Stripe.
var (
	ErrLiveKey  = errors.New("issuing runs in Stripe test mode only: the key must be a test key (rk_test_ or sk_test_)")
	ErrNotReady = errors.New("the purchase case is not an approved one this card can be made for")
)

var (
	currencyCode = regexp.MustCompile(`^[a-z]{3}$`)
	categoryName = regexp.MustCompile(`^[a-z0-9_]{1,64}$`)
	caseID       = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
)

// Check says why c isn't a case a card can be made for at now, or nil.
func (c Case) Check(now time.Time) error {
	switch {
	case !caseID.MatchString(c.ID) || c.Version < 1:
		return fmt.Errorf("%w: it names no case id and version", ErrNotReady)
	case c.State != "approved":
		return fmt.Errorf("%w: its state is %q, not approved", ErrNotReady, c.State)
	case c.MaxAmount <= 0:
		return fmt.Errorf("%w: its maximum amount is not above zero", ErrNotReady)
	case !currencyCode.MatchString(c.Currency):
		return fmt.Errorf("%w: its currency is not a three-letter code in lower case", ErrNotReady)
	case len(c.Categories) == 0:
		return fmt.Errorf("%w: it names no merchant category, and a card is never left open to every one", ErrNotReady)
	case !c.Expires.After(now):
		return fmt.Errorf("%w: it has ended", ErrNotReady)
	case c.Expires.Sub(now) > MaxWindow:
		return fmt.Errorf("%w: it runs longer than %s", ErrNotReady, MaxWindow)
	}
	for _, cat := range c.Categories {
		if !categoryName.MatchString(cat) {
			return fmt.Errorf("%w: %q is not a merchant category name", ErrNotReady, cat)
		}
	}
	return nil
}

// Client is Stripe Issuing, in test mode, with a restricted key that can
// create, read and cancel cards: nothing that reveals a number.
type Client struct {
	Base       string // https://api.stripe.com, or a test server
	Key        string
	Cardholder string // the test cardholder every card is issued to
	HTTP       *http.Client
}

// New is a Client for key, refusing any key that isn't in test mode.
func New(key, cardholder string) (*Client, error) {
	if !strings.HasPrefix(key, "rk_test_") && !strings.HasPrefix(key, "sk_test_") {
		return nil, ErrLiveKey
	}
	if cardholder == "" {
		return nil, errors.New("issuing: a test cardholder is required")
	}
	return &Client{Base: "https://api.stripe.com", Key: key, Cardholder: cardholder, HTTP: &http.Client{Timeout: 20 * time.Second}}, nil
}

// form is c's card request: a virtual card, active, for the cardholder,
// with its spending limited to exactly the case, and the case and permit in
// its metadata so any card can be traced to what allowed it.
func (cl *Client) form(c Case, permit string) url.Values {
	f := url.Values{}
	f.Set("type", "virtual")
	f.Set("status", "active")
	f.Set("cardholder", cl.Cardholder)
	f.Set("currency", c.Currency)
	max := strconv.FormatInt(c.MaxAmount, 10)
	// The most one authorization may be, and the most the card may ever
	// spend: the same amount, so a second authorization finds no room.
	f.Set("spending_controls[spending_limits][0][amount]", max)
	f.Set("spending_controls[spending_limits][0][interval]", "per_authorization")
	f.Set("spending_controls[spending_limits][1][amount]", max)
	f.Set("spending_controls[spending_limits][1][interval]", "all_time")
	for i, cat := range c.Categories {
		f.Set(fmt.Sprintf("spending_controls[allowed_categories][%d]", i), cat)
	}
	f.Set("metadata[case]", c.ID)
	f.Set("metadata[case_version]", strconv.Itoa(c.Version))
	f.Set("metadata[permit]", permit)
	return f
}

// Create makes the card for c, under permit (the KIFF permit's id). The
// idempotency key is the case and version, so a retry never makes a second
// card for the same decision.
func (cl *Client) Create(ctx context.Context, c Case, permit string, now time.Time) (Card, error) {
	if err := c.Check(now); err != nil {
		return Card{}, err
	}
	var raw stripeCard
	if err := cl.do(ctx, http.MethodPost, "/v1/issuing/cards", cl.form(c, permit), "thune-"+c.ID+"-v"+strconv.Itoa(c.Version), &raw); err != nil {
		return Card{}, err
	}
	return raw.card(), nil
}

// Cancel cancels card id for good. A card already canceled is done.
func (cl *Client) Cancel(ctx context.Context, id string) error {
	f := url.Values{}
	f.Set("status", "canceled")
	var raw stripeCard
	return cl.do(ctx, http.MethodPost, "/v1/issuing/cards/"+url.PathEscape(id), f, "", &raw)
}

// Active lists the cardholder's active cards, newest first, up to limit
// (at most 100): what the sweep checks against planque's records.
func (cl *Client) Active(ctx context.Context, limit int) ([]Card, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	q := url.Values{}
	q.Set("status", "active")
	q.Set("cardholder", cl.Cardholder)
	q.Set("limit", strconv.Itoa(limit))
	var page struct {
		Data []stripeCard `json:"data"`
	}
	if err := cl.do(ctx, http.MethodGet, "/v1/issuing/cards?"+q.Encode(), nil, "", &page); err != nil {
		return nil, err
	}
	out := make([]Card, 0, len(page.Data))
	for _, r := range page.Data {
		out = append(out, r.card())
	}
	return out, nil
}

// stripeCard is the fields of Stripe's card object planque reads. It has no
// field for the number or code, so neither can be decoded by mistake.
type stripeCard struct {
	ID               string            `json:"id"`
	Status           string            `json:"status"`
	Last4            string            `json:"last4"`
	Currency         string            `json:"currency"`
	Metadata         map[string]string `json:"metadata"`
	SpendingControls struct {
		SpendingLimits []struct {
			Amount   int64  `json:"amount"`
			Interval string `json:"interval"`
		} `json:"spending_limits"`
	} `json:"spending_controls"`
}

func (r stripeCard) card() Card {
	c := Card{ID: r.ID, Status: r.Status, Last4: r.Last4, Currency: r.Currency, Case: r.Metadata["case"], Permit: r.Metadata["permit"]}
	c.Version, _ = strconv.Atoi(r.Metadata["case_version"])
	for _, l := range r.SpendingControls.SpendingLimits {
		if l.Interval == "per_authorization" {
			c.Max = l.Amount
		}
	}
	return c
}

// do sends one request. Stripe's error message is passed on, never the
// request's form (which carries nothing secret, but stays out of logs).
func (cl *Client) do(ctx context.Context, method, path string, form url.Values, idem string, out any) error {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(cl.Base, "/")+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cl.Key)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	resp, err := cl.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("issuing: Stripe could not be reached: %w", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		var e struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(b, &e)
		return fmt.Errorf("issuing: Stripe answered %d: %s", resp.StatusCode, e.Error.Message)
	}
	return json.Unmarshal(b, out)
}
