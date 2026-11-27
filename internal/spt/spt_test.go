package spt

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeStripe answers the issued-token routes as Stripe's docs describe them,
// and records what it was asked.
type fakeStripe struct {
	mu      sync.Mutex
	tokens  map[string]*Token
	issued  []url.Values
	idem    []string
	revoked []string
	n       int
	// lie, when set, changes the limits of the next issued token's answer.
	lie func(*Token)
	// live answers in live mode.
	live bool
}

func (f *fakeStripe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer sk_test_") {
		w.WriteHeader(401)
		return
	}
	body, _ := io.ReadAll(r.Body)
	form, _ := url.ParseQuery(string(body))
	path := strings.TrimPrefix(r.URL.Path, "/v1/shared_payment/issued_tokens")
	answer := func(t *Token) {
		c := *t
		c.Livemode = f.live
		_ = json.NewEncoder(w).Encode(c)
	}
	switch {
	case r.Method == "POST" && path == "":
		f.issued = append(f.issued, form)
		f.idem = append(f.idem, r.Header.Get("Idempotency-Key"))
		f.n++
		max, _ := strconv.ParseInt(form.Get("usage_limits[max_amount]"), 10, 64)
		exp, _ := strconv.ParseInt(form.Get("usage_limits[expires_at]"), 10, 64)
		t := &Token{ID: fmt.Sprintf("spt_test%d", f.n), Status: "active",
			Limits: Limits{Currency: form.Get("usage_limits[currency]"), MaxAmount: max, ExpiresAt: exp}}
		if f.lie != nil {
			f.lie(t)
			f.lie = nil
		}
		f.tokens[t.ID] = t
		answer(t)
	case r.Method == "POST" && strings.HasSuffix(path, "/revoke"):
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/revoke")
		f.revoked = append(f.revoked, id)
		if t := f.tokens[id]; t != nil {
			t.Status, t.DeactivatedReason = "deactivated", "revoked"
			answer(t)
			return
		}
		w.WriteHeader(404)
	case r.Method == "GET":
		if t := f.tokens[strings.TrimPrefix(path, "/")]; t != nil {
			answer(t)
			return
		}
		w.WriteHeader(404)
	default:
		w.WriteHeader(400)
	}
}

type memStore struct {
	cards  map[string]Card
	marked map[string]time.Time
}

func (m *memStore) CardFor(_ context.Context, tenant string) (Card, bool, error) {
	c, ok := m.cards[tenant]
	return c, ok, nil
}

func (m *memStore) MarkCheckout(_ context.Context, tenant, key string, now time.Time) (bool, error) {
	k := tenant + "|" + key
	if at, ok := m.marked[k]; ok && now.Sub(at) < 24*time.Hour {
		return false, nil
	}
	m.marked[k] = now
	return true, nil
}

var now0 = time.Date(2026, 11, 27, 1, 0, 0, 0, time.UTC)

func decimals(c string) (int, bool) {
	switch strings.ToUpper(c) {
	case "EUR", "USD":
		return 2, true
	case "JPY":
		return 0, true
	}
	return 0, false // ISK, HUF, TWD, UGX and anything unknown
}

func setup(t *testing.T, key string) (*Issuer, *fakeStripe, *memStore) {
	t.Helper()
	f := &fakeStripe{tokens: map[string]*Token{}}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	st := &memStore{cards: map[string]Card{
		"tnt_a": {Tenant: "tnt_a", PaymentMethod: "pm_a", Brand: "visa", Last4: "4242"},
		// A store bug that hands tenant B tenant A's card must still be refused.
		"tnt_b": {Tenant: "tnt_a", PaymentMethod: "pm_a", Brand: "visa", Last4: "4242"},
	}, marked: map[string]time.Time{}}
	c := &Client{HTTP: srv.Client(), Base: srv.URL, Version: "2026-09-30.preview",
		Key: func(context.Context) (string, error) { return key, nil }}
	return &Issuer{On: true, Client: c, Store: st, Decimals: decimals, TestTenants: map[string]bool{"tnt_proof": true},
		Now: func() time.Time { return now0 }}, f, st
}

func shop() Checkout {
	return Checkout{Seller: Seller{Profile: "profile_test_shop", Name: "Test Shop", Verified: true},
		Lines: []Line{{Name: "Paper, A4", Quantity: 2, Amount: Money{1250, "EUR"}}}, Total: Money{2500, "EUR"}}
}

func approval() Approval {
	return Approval{Tenant: "tnt_a", Taf: "agt_buyer", Purchase: "pur_1", Permit: "prm_1",
		SellerProfile: "profile_test_shop", Amount: Money{2500, "EUR"}, ExpiresAt: now0.Add(10 * time.Minute)}
}

// A purchase approved for exactly this checkout issues one token, limited to
// the approved amount, currency, shop and expiry, with the permit as its
// idempotency key; settling reads Stripe's state.
func TestAnApprovedPurchaseIssuesOneTokenOnItsExactLimits(t *testing.T) {
	i, f, _ := setup(t, "sk_test_x")
	p, err := i.Pay(context.Background(), approval(), shop())
	if err != nil {
		t.Fatal(err)
	}
	if len(f.issued) != 1 || f.idem[0] != "prm_1" {
		t.Fatalf("issued %d, idempotency %v", len(f.issued), f.idem)
	}
	got := f.issued[0]
	for k, want := range map[string]string{
		"payment_method": "pm_a", "seller_details[network_business_profile]": "profile_test_shop",
		"usage_limits[currency]": "eur", "usage_limits[max_amount]": "2500",
		"usage_limits[expires_at]": strconv.FormatInt(now0.Add(10*time.Minute).Unix(), 10), "shared_metadata[purchase]": "pur_1",
	} {
		if got.Get(k) != want {
			t.Errorf("%s = %q, want %q", k, got.Get(k), want)
		}
	}
	if p.State != StateIssued || p.Brand != "visa" || p.Last4 != "4242" || p.Token == "" {
		t.Fatalf("purchase: %#v", p)
	}
	if s, err := i.Settle(context.Background(), p); err != nil || s != StateIssued {
		t.Fatalf("an active token: %s %v", s, err)
	}
	f.tokens[p.Token].Status, f.tokens[p.Token].DeactivatedReason = "deactivated", "consumed"
	if s, _ := i.Settle(context.Background(), p); s != StatePaid {
		t.Fatalf("a consumed token: %s", s)
	}
}

// Every check refuses before anything reaches Stripe.
func TestEachRefusalReachesNoStripe(t *testing.T) {
	cases := []struct {
		name string
		edit func(*Issuer, *Approval, *Checkout)
		want error
	}{
		{"off", func(i *Issuer, _ *Approval, _ *Checkout) { i.On = false }, ErrOff},
		{"a proof test account", func(_ *Issuer, a *Approval, _ *Checkout) { a.Tenant = "tnt_proof" }, ErrTestTenant},
		{"a seller read from a page", func(_ *Issuer, _ *Approval, c *Checkout) { c.Seller.Verified = false }, ErrSellerUnverified},
		{"another shop than approved", func(_ *Issuer, a *Approval, _ *Checkout) { a.SellerProfile = "profile_test_other" }, ErrMismatch},
		{"a total changed after the approval", func(_ *Issuer, _ *Approval, c *Checkout) { c.Total.Minor = 2600 }, ErrMismatch},
		{"another currency", func(_ *Issuer, _ *Approval, c *Checkout) { c.Total.Currency = "USD" }, ErrMismatch},
		{"a special currency", func(_ *Issuer, a *Approval, c *Checkout) { a.Amount.Currency, c.Total.Currency = "HUF", "HUF" }, ErrCurrency},
		{"a zero amount", func(_ *Issuer, a *Approval, c *Checkout) { a.Amount.Minor, c.Total.Minor = 0, 0 }, ErrAmount},
		{"an expiry past 10 minutes", func(_ *Issuer, a *Approval, _ *Checkout) { a.ExpiresAt = now0.Add(11 * time.Minute) }, ErrExpiry},
		{"an expiry passed", func(_ *Issuer, a *Approval, _ *Checkout) { a.ExpiresAt = now0 }, ErrExpiry},
		{"no saved card", func(_ *Issuer, a *Approval, _ *Checkout) { a.Tenant = "tnt_none" }, ErrNoCard},
		{"another account's card", func(_ *Issuer, a *Approval, _ *Checkout) { a.Tenant = "tnt_b" }, ErrWrongTenant},
	}
	for _, c := range cases {
		i, f, _ := setup(t, "sk_test_x")
		a, co := approval(), shop()
		c.edit(i, &a, &co)
		if _, err := i.Pay(context.Background(), a, co); !errors.Is(err, c.want) {
			t.Errorf("%s: %v, want %v", c.name, err, c.want)
		}
		if f.n != 0 {
			t.Errorf("%s: %d tokens issued", c.name, f.n)
		}
	}
}

// Two approvals for the same cart pay it once: the second, a new permit
// included, is refused for 24 hours.
func TestTheSameCartIsNeverPaidTwice(t *testing.T) {
	i, f, _ := setup(t, "sk_test_x")
	if _, err := i.Pay(context.Background(), approval(), shop()); err != nil {
		t.Fatal(err)
	}
	a := approval()
	a.Permit, a.Purchase = "prm_2", "pur_2"
	if _, err := i.Pay(context.Background(), a, shop()); !errors.Is(err, ErrRepeat) {
		t.Fatalf("a second purchase of the cart: %v", err)
	}
	if f.n != 1 {
		t.Fatalf("%d tokens issued", f.n)
	}
	i.Now = func() time.Time { return now0.Add(25 * time.Hour) }
	a.ExpiresAt = now0.Add(25*time.Hour + 5*time.Minute)
	if _, err := i.Pay(context.Background(), a, shop()); err != nil {
		t.Fatalf("the same cart a day later: %v", err)
	}
}

// Only a test key is taken, and a live answer is refused.
func TestOnlyTestModeIsTouched(t *testing.T) {
	for _, key := range []string{"sk_live_x", "rk_live_x", "pk_test_x", ""} {
		i, f, _ := setup(t, key)
		if _, err := i.Pay(context.Background(), approval(), shop()); !errors.Is(err, ErrLiveKey) {
			t.Errorf("%q: %v", key, err)
		}
		if f.n != 0 {
			t.Errorf("%q reached Stripe", key)
		}
	}
	i, f, _ := setup(t, "sk_test_x")
	f.live = true
	if _, err := i.Pay(context.Background(), approval(), shop()); !errors.Is(err, ErrLiveAnswer) {
		t.Errorf("a live answer: %v", err)
	}
}

// A token whose limits aren't the ones asked is revoked and refused.
func TestATokenWithOtherLimitsIsRevoked(t *testing.T) {
	i, f, _ := setup(t, "sk_test_x")
	f.lie = func(t *Token) { t.Limits.MaxAmount = 250000 }
	if _, err := i.Pay(context.Background(), approval(), shop()); err == nil {
		t.Fatal("a token for 100 times the amount was taken")
	}
	if len(f.revoked) != 1 {
		t.Fatalf("revoked: %v", f.revoked)
	}
}

// Settling follows Stripe's state: expired or revoked is not paid; limits
// that changed, or a token resolved to another processor, is an incident.
func TestSettleFollowsStripesState(t *testing.T) {
	i, f, _ := setup(t, "sk_test_x")
	p, err := i.Pay(context.Background(), approval(), shop())
	if err != nil {
		t.Fatal(err)
	}
	tok := f.tokens[p.Token]
	for _, c := range []struct{ status, reason, want string }{
		{"deactivated", "expired", StateNotPaid},
		{"deactivated", "revoked", StateNotPaid},
		{"deactivated", "resolved", StateIncident},
		{"deactivated", "consumed", StatePaid},
	} {
		tok.Status, tok.DeactivatedReason = c.status, c.reason
		if s, _ := i.Settle(context.Background(), p); s != c.want {
			t.Errorf("%s/%s: %s, want %s", c.status, c.reason, s, c.want)
		}
	}
	tok.Limits.MaxAmount++
	if s, _ := i.Settle(context.Background(), p); s != StateIncident {
		t.Errorf("changed limits: %s", s)
	}
}

// A halt revokes every token still issued, and only those.
func TestStopRevokesTheIssuedTokens(t *testing.T) {
	i, f, _ := setup(t, "sk_test_x")
	p, err := i.Pay(context.Background(), approval(), shop())
	if err != nil {
		t.Fatal(err)
	}
	paid := Purchase{ID: "pur_0", State: StatePaid, Token: "spt_old"}
	if err := i.Stop(context.Background(), []Purchase{p, paid}); err != nil {
		t.Fatal(err)
	}
	if len(f.revoked) != 1 || f.revoked[0] != p.Token {
		t.Fatalf("revoked: %v", f.revoked)
	}
	if s, _ := i.Settle(context.Background(), p); s != StateNotPaid {
		t.Fatalf("after the stop: %s", s)
	}
}

// The token id never prints: not with %v, %+v or %#v, not in a log line,
// not in the purchase's JSON.
func TestTheTokenIDNeverPrints(t *testing.T) {
	i, _, _ := setup(t, "sk_test_x")
	p, err := i.Pay(context.Background(), approval(), shop())
	if err != nil {
		t.Fatal(err)
	}
	tok, _ := i.Client.Get(context.Background(), p.Token)
	var buf bytes.Buffer
	slog.New(slog.NewTextHandler(&buf, nil)).Info("paid", "purchase", p, "token", tok)
	j, _ := json.Marshal(p)
	out := fmt.Sprintf("%v %+v %#v %v %+v %#v %s %s", p, p, p, tok, tok, tok, buf.String(), j)
	if strings.Contains(out, "spt_") {
		t.Fatalf("a token id printed: %s", out)
	}
}

func sign(payload []byte, secret string, ts int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(payload)
	return "t=" + strconv.FormatInt(ts, 10) + ",v1=" + hex.EncodeToString(mac.Sum(nil))
}

// A webhook counts only with Stripe's signature, fresh, for a token event;
// it only names the token to read again.
func TestAWebhookCountsOnlySignedAndFresh(t *testing.T) {
	payload := []byte(`{"id":"evt_1","type":"shared_payment.issued_token.used","data":{"object":{"id":"spt_test1","object":"shared_payment.issued_token"}}}`)
	ts := now0.Unix()
	e, err := Verify(payload, sign(payload, "whsec_x", ts), "whsec_x", now0)
	if err != nil || e.Token != "spt_test1" || e.Type != "shared_payment.issued_token.used" {
		t.Fatalf("a signed event: %+v %v", e, err)
	}
	for name, c := range map[string]struct {
		payload []byte
		header  string
		at      time.Time
	}{
		"another secret":  {payload, sign(payload, "whsec_y", ts), now0},
		"a changed body":  {bytes.Replace(payload, []byte("used"), []byte("deactivated"), 1), sign(payload, "whsec_x", ts), now0},
		"a replay later":  {payload, sign(payload, "whsec_x", ts), now0.Add(6 * time.Minute)},
		"no signature":    {payload, "t=" + strconv.FormatInt(ts, 10), now0},
		"no header":       {payload, "", now0},
		"not a token one": {[]byte(`{"id":"evt_2","type":"charge.refunded","data":{"object":{"id":"ch_1"}}}`), sign([]byte(`{"id":"evt_2","type":"charge.refunded","data":{"object":{"id":"ch_1"}}}`), "whsec_x", ts), now0},
	} {
		if _, err := Verify(c.payload, c.header, "whsec_x", c.at); err == nil {
			t.Errorf("%s was taken", name)
		}
	}
}
