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
	"slices"
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
	// the Card's allowed ones. None: the card is blocked from the high-risk
	// categories only (HighRisk), and the live authorization's exact shop
	// match is the fence. A high-risk category is never allowed.
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

// HighRisk are the merchant categories no purchase card may ever be used
// at, whatever the Card allows: cash, quasi-cash and money transfer
// (crypto purchases are coded as non-financial money orders), and
// gambling. A card that leaks can't become cash. Stripe's names
// (docs.stripe.com/api/issuing/cards/create, spending_controls).
var HighRisk = []string{
	"automated_cash_disburse", "manual_cash_disburse",
	"financial_institutions", "non_fi_money_orders", "non_fi_stored_value_card_purchase_load", "wires_money_orders",
	"betting_casino_gambling", "government_licensed_online_casions_online_gambling_us_region_only",
	"government_licensed_horse_dog_racing_us_region_only", "government_owned_lotteries_non_us_region", "government_owned_lotteries_us_region_only",
}

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
	case !c.Expires.After(now):
		return fmt.Errorf("%w: it has ended", ErrNotReady)
	case c.Expires.Sub(now) > MaxWindow:
		return fmt.Errorf("%w: it runs longer than %s", ErrNotReady, MaxWindow)
	}
	for _, cat := range c.Categories {
		if !categoryName.MatchString(cat) {
			return fmt.Errorf("%w: %q is not a merchant category name", ErrNotReady, cat)
		}
		if slices.Contains(HighRisk, cat) {
			return fmt.Errorf("%w: it allows %s, and a purchase card never opens cash, money transfer or gambling", ErrNotReady, cat)
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
	// Stripe takes allowed or blocked categories, never both: an allowed
	// list already blocks every other category, high-risk ones included
	// (Check refuses them in it); with none, the high-risk ones are blocked.
	if len(c.Categories) > 0 {
		for i, cat := range c.Categories {
			f.Set(fmt.Sprintf("spending_controls[allowed_categories][%d]", i), cat)
		}
	} else {
		for i, cat := range HighRisk {
			f.Set(fmt.Sprintf("spending_controls[blocked_categories][%d]", i), cat)
		}
	}
	// Stripe cancels the card itself after its first payment (a non-zero
	// authorization), whatever planque does next (T2.2).
	f.Set("lifecycle_controls[cancel_after][payment_count]", "1")
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

// Cancel cancels card id for good. A card already canceled is done: when
// Stripe refuses the change, the card is read, and a canceled one counts
// as canceled.
func (cl *Client) Cancel(ctx context.Context, id string) error {
	f := url.Values{}
	f.Set("status", "canceled")
	var raw stripeCard
	err := cl.do(ctx, http.MethodPost, "/v1/issuing/cards/"+url.PathEscape(id), f, "", &raw)
	if err == nil {
		return nil
	}
	if card, gerr := cl.Get(ctx, id); gerr == nil && card.Status == "canceled" {
		return nil
	}
	return err
}

// Get reads card id: its status, case and limit, never its number.
func (cl *Client) Get(ctx context.Context, id string) (Card, error) {
	var raw stripeCard
	if err := cl.do(ctx, http.MethodGet, "/v1/issuing/cards/"+url.PathEscape(id), nil, "", &raw); err != nil {
		return Card{}, err
	}
	return raw.card(), nil
}

// maxPages bounds how many pages Active reads.
const maxPages = 50

// ErrTooManyCards: more active cards than Active reads; the sweep reports
// it rather than treat a partial list as the whole.
var ErrTooManyCards = errors.New("issuing: more active cards than one sweep reads")

// Active lists every active card of the cardholder, newest first, pages of
// pageSize (at most 100) followed to the end: what the sweep checks
// against planque's records. The cards read before an error are returned
// with it.
func (cl *Client) Active(ctx context.Context, pageSize int) ([]Card, error) {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}
	var out []Card
	after := ""
	for range maxPages {
		q := url.Values{}
		q.Set("status", "active")
		q.Set("cardholder", cl.Cardholder)
		q.Set("limit", strconv.Itoa(pageSize))
		if after != "" {
			q.Set("starting_after", after)
		}
		var page struct {
			Data    []stripeCard `json:"data"`
			HasMore bool         `json:"has_more"`
		}
		if err := cl.do(ctx, http.MethodGet, "/v1/issuing/cards?"+q.Encode(), nil, "", &page); err != nil {
			return out, err
		}
		for _, r := range page.Data {
			out = append(out, r.card())
		}
		if !page.HasMore || len(page.Data) == 0 {
			return out, nil
		}
		after = page.Data[len(page.Data)-1].ID
	}
	return out, ErrTooManyCards
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

// PreviewVersion is the Stripe-Version of the private preview that makes a
// PaymentMethod from an Issuing card (and issues a shared payment token).
const PreviewVersion = "2026-04-22.preview"

// PaymentMethod makes a card PaymentMethod (pm_…) from Issuing card id, in
// the same account, without the card's number or code ever being read: what
// a shared payment token is issued from (path a). It needs the private
// preview.
func (cl *Client) PaymentMethod(ctx context.Context, id string) (string, error) {
	f := url.Values{}
	f.Set("type", "card")
	f.Set("card[issuing_card]", id)
	var pm struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Livemode bool   `json:"livemode"`
	}
	if err := cl.doWith(ctx, http.MethodPost, "/v1/payment_methods", f, "thune-pm-"+id, PreviewVersion, &pm); err != nil {
		return "", err
	}
	if pm.Livemode {
		return "", ErrLiveKey
	}
	if !strings.HasPrefix(pm.ID, "pm_") || pm.Type != "card" {
		return "", errors.New("issuing: Stripe's answer isn't a card PaymentMethod")
	}
	return pm.ID, nil
}

// do sends one request. Stripe's error message is passed on, never the
// request's form (which carries nothing secret, but stays out of logs).
func (cl *Client) do(ctx context.Context, method, path string, form url.Values, idem string, out any) error {
	return cl.doWith(ctx, method, path, form, idem, "", out)
}

// doWith is do with a Stripe-Version, when given.
func (cl *Client) doWith(ctx context.Context, method, path string, form url.Values, idem, version string, out any) error {
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
	if version != "" {
		req.Header.Set("Stripe-Version", version)
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
