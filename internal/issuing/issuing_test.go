package issuing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeStripe answers the card calls planque makes, and keeps every request
// it was sent so a test can read what went to Stripe.
type fakeStripe struct {
	mu    sync.Mutex
	sent  []sentReq
	cards map[string]map[string]any
	n     int
}

type sentReq struct {
	Method, Path, Query, Idem string
	Form                      url.Values
}

func newFake(t *testing.T) (*fakeStripe, *Client) {
	t.Helper()
	f := &fakeStripe{cards: map[string]map[string]any{}}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	cl, err := New("rk_test_fake", "ich_test")
	if err != nil {
		t.Fatal(err)
	}
	cl.Base, cl.HTTP = srv.URL, srv.Client()
	return f, cl
}

func (f *fakeStripe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	form, _ := url.ParseQuery(string(b))
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sent = append(f.sent, sentReq{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Get("Idempotency-Key"), form})
	if r.Header.Get("Authorization") != "Bearer rk_test_fake" {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
		return
	}
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/v1/issuing/cards":
		f.n++
		id := fmt.Sprintf("ic_test_%04d", f.n)
		card := map[string]any{
			"id": id, "status": form.Get("status"), "last4": "4242", "currency": form.Get("currency"),
			"metadata": map[string]string{"case": form.Get("metadata[case]"), "case_version": form.Get("metadata[case_version]"), "permit": form.Get("metadata[permit]")},
			"spending_controls": map[string]any{"spending_limits": []map[string]any{
				{"amount": json.Number(form.Get("spending_controls[spending_limits][0][amount]")), "interval": form.Get("spending_controls[spending_limits][0][interval]")},
			}},
			// What Stripe would add on expand=number: the client never asks,
			// and must never decode it even if it came.
			"number": "4242424242424242", "cvc": "123",
		}
		f.cards[id] = card
		_ = json.NewEncoder(w).Encode(card)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/v1/issuing/cards/"):
		id := strings.TrimPrefix(r.URL.Path, "/v1/issuing/cards/")
		c := f.cards[id]
		if c == nil {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"error":{"message":"No such card"}}`))
			return
		}
		// Stripe refuses to change a canceled card.
		if c["status"] == "canceled" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":{"message":"This card has been canceled and can no longer be updated."}}`))
			return
		}
		c["status"] = form.Get("status")
		_ = json.NewEncoder(w).Encode(c)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/issuing/cards/"):
		c := f.cards[strings.TrimPrefix(r.URL.Path, "/v1/issuing/cards/")]
		if c == nil {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"error":{"message":"No such card"}}`))
			return
		}
		_ = json.NewEncoder(w).Encode(c)
	case r.Method == http.MethodGet && r.URL.Path == "/v1/issuing/cards":
		var ids []string
		for id, c := range f.cards {
			if c["status"] == "active" {
				ids = append(ids, id)
			}
		}
		slices.Sort(ids)
		if after := r.URL.Query().Get("starting_after"); after != "" {
			i := slices.Index(ids, after)
			ids = ids[i+1:]
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		more := len(ids) > limit
		if more {
			ids = ids[:limit]
		}
		data := []map[string]any{}
		for _, id := range ids {
			data = append(data, f.cards[id])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "has_more": more})
	default:
		w.WriteHeader(404)
	}
}

func approved(now time.Time) Case {
	return Case{ID: "case_1", Version: 2, State: "approved", MaxAmount: 4250, Currency: "eur",
		Categories: []string{"book_stores"}, Expires: now.Add(10 * time.Minute)}
}

// A card's spending is the case's, copied: one authorization at most the
// case's maximum, and never more over the card's life; only the case's
// merchant categories; the case and permit in its metadata; and one card per
// case version however often it's asked.
func TestACardIsLimitedToExactlyItsCase(t *testing.T) {
	f, cl := newFake(t)
	now := time.Now()
	card, err := cl.Create(context.Background(), approved(now), "jti_1", now)
	if err != nil {
		t.Fatal(err)
	}
	if card.Case != "case_1" || card.Version != 2 || card.Permit != "jti_1" || card.Max != 4250 || card.Last4 != "4242" {
		t.Fatalf("card: %+v", card)
	}
	got := f.sent[0].Form
	want := map[string]string{
		"type": "virtual", "status": "active", "cardholder": "ich_test", "currency": "eur",
		"spending_controls[spending_limits][0][amount]": "4250", "spending_controls[spending_limits][0][interval]": "per_authorization",
		"spending_controls[spending_limits][1][amount]": "4250", "spending_controls[spending_limits][1][interval]": "all_time",
		"spending_controls[allowed_categories][0]":        "book_stores",
		"lifecycle_controls[cancel_after][payment_count]": "1",
		"metadata[case]": "case_1", "metadata[case_version]": "2", "metadata[permit]": "jti_1",
	}
	for k, v := range want {
		if got.Get(k) != v {
			t.Errorf("%s = %q, want %q", k, got.Get(k), v)
		}
	}
	for k := range got {
		if strings.Contains(k, "expand") || strings.Contains(k, "number") {
			t.Errorf("the request asks for %s", k)
		}
	}
	if f.sent[0].Idem != "thune-case_1-v2" {
		t.Fatalf("idempotency key: %q", f.sent[0].Idem)
	}
}

// A card is never made for a case that isn't approved, has ended, runs past
// the window, names no category, or has no amount; nothing is sent.
func TestNoCardForACaseThatIsNotReady(t *testing.T) {
	f, cl := newFake(t)
	now := time.Now()
	for name, change := range map[string]func(*Case){
		"pending":           func(c *Case) { c.State = "pending" },
		"ended":             func(c *Case) { c.Expires = now.Add(-time.Second) },
		"past the window":   func(c *Case) { c.Expires = now.Add(MaxWindow + time.Minute) },
		"cash":              func(c *Case) { c.Categories = []string{"book_stores", "automated_cash_disburse"} },
		"gambling":          func(c *Case) { c.Categories = []string{"betting_casino_gambling"} },
		"no amount":         func(c *Case) { c.MaxAmount = 0 },
		"upper case":        func(c *Case) { c.Currency = "EUR" },
		"no version":        func(c *Case) { c.Version = 0 },
		"odd category name": func(c *Case) { c.Categories = []string{"books & more"} },
	} {
		c := approved(now)
		change(&c)
		if _, err := cl.Create(context.Background(), c, "jti", now); !errors.Is(err, ErrNotReady) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if len(f.sent) != 0 {
		t.Fatalf("%d requests went to Stripe", len(f.sent))
	}
}

// Only a test key: a live or unknown key is refused before any client exists.
func TestOnlyATestKey(t *testing.T) {
	for _, k := range []string{"rk_live_x", "sk_live_x", "pk_test_x", "", "key"} {
		if _, err := New(k, "ich"); !errors.Is(err, ErrLiveKey) {
			t.Errorf("%q: %v", k, err)
		}
	}
	if _, err := New("rk_test_x", ""); err == nil {
		t.Error("no cardholder accepted")
	}
}

// A card's number and code never reach planque's types: even when Stripe's
// answer carries them, nothing decoded holds them.
func TestACardsNumberIsNeverDecoded(t *testing.T) {
	_, cl := newFake(t)
	now := time.Now()
	card, err := cl.Create(context.Background(), approved(now), "jti", now)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(card)
	if strings.Contains(string(b), "4242424242424242") || strings.Contains(string(b), "123\"") {
		t.Fatalf("a card's number or code was kept: %s", b)
	}
	list, err := cl.Active(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = json.Marshal(list)
	if strings.Contains(string(b), "4242424242424242") {
		t.Fatalf("the list kept a number: %s", b)
	}
}

// Cancel ends a card; the active list no longer has it.
func TestACanceledCardIsNoLongerActive(t *testing.T) {
	_, cl := newFake(t)
	now := time.Now()
	card, _ := cl.Create(context.Background(), approved(now), "jti", now)
	if err := cl.Cancel(context.Background(), card.ID); err != nil {
		t.Fatal(err)
	}
	list, _ := cl.Active(context.Background(), 10)
	if len(list) != 0 {
		t.Fatalf("still active: %+v", list)
	}
}

// Stripe takes allowed or blocked categories, never both. An allowed list
// blocks every other category, so it carries no block list; a case with no
// categories gets the high-risk block list, and is never open to cash,
// money transfer or gambling.
func TestHighRiskCategoriesAreNeverOpen(t *testing.T) {
	f, cl := newFake(t)
	now := time.Now()
	if _, err := cl.Create(context.Background(), approved(now), "jti", now); err != nil {
		t.Fatal(err)
	}
	for k := range f.sent[0].Form {
		if strings.HasPrefix(k, "spending_controls[blocked_categories]") {
			t.Errorf("an allowed list also sent %s", k)
		}
	}
	c := approved(now)
	c.ID, c.Categories = "case_2", nil
	if _, err := cl.Create(context.Background(), c, "jti", now); err != nil {
		t.Fatal(err)
	}
	form := f.sent[1].Form
	var blocked []string
	for i := 0; form.Get(fmt.Sprintf("spending_controls[blocked_categories][%d]", i)) != ""; i++ {
		blocked = append(blocked, form.Get(fmt.Sprintf("spending_controls[blocked_categories][%d]", i)))
	}
	for _, cat := range []string{"automated_cash_disburse", "manual_cash_disburse", "wires_money_orders", "non_fi_money_orders", "financial_institutions", "betting_casino_gambling"} {
		if !slices.Contains(blocked, cat) {
			t.Errorf("%s isn't blocked: %v", cat, blocked)
		}
	}
	if form.Get("spending_controls[allowed_categories][0]") != "" || form.Get("lifecycle_controls[cancel_after][payment_count]") != "1" {
		t.Errorf("a case with no categories: %v", form)
	}
}

// A second cancel of a canceled card counts as done; a card that doesn't
// exist is still an error.
func TestCancelingTwiceIsDone(t *testing.T) {
	_, cl := newFake(t)
	now := time.Now()
	card, _ := cl.Create(context.Background(), approved(now), "jti", now)
	if err := cl.Cancel(context.Background(), card.ID); err != nil {
		t.Fatal(err)
	}
	if err := cl.Cancel(context.Background(), card.ID); err != nil {
		t.Errorf("a second cancel: %v", err)
	}
	if err := cl.Cancel(context.Background(), "ic_none"); err == nil {
		t.Error("a card that doesn't exist was canceled")
	}
}

// The sweep's list is every active card, page after page, not the first
// page alone.
func TestActiveReadsEveryPage(t *testing.T) {
	f, cl := newFake(t)
	now := time.Now()
	for i := range 7 {
		c := approved(now)
		c.ID = fmt.Sprintf("case_%d", i)
		if _, err := cl.Create(context.Background(), c, "jti", now); err != nil {
			t.Fatal(err)
		}
	}
	before := len(f.sent)
	list, err := cl.Active(context.Background(), 3)
	if err != nil || len(list) != 7 {
		t.Fatalf("%d cards, %v", len(list), err)
	}
	if pages := len(f.sent) - before; pages != 3 {
		t.Errorf("%d pages read, want 3", pages)
	}
	seen := map[string]bool{}
	for _, c := range list {
		if seen[c.ID] {
			t.Errorf("%s twice", c.ID)
		}
		seen[c.ID] = true
	}
}
