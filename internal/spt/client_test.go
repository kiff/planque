package spt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func client(t *testing.T, key string, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &Client{HTTP: srv.Client(), Base: srv.URL, Version: "2026-04-22.preview", Key: func(context.Context) (string, error) { return key, nil }}
}

var in = Issue{PaymentMethod: "pm_1", SellerProfile: "profile_test_1", Currency: "EUR", MaxAmount: 1999, ExpiresAt: time.Unix(1798761600, 0), IdempotencyKey: "k", Purchase: "case_7"}

// Only a test key; never a live answer.
func TestOnlyTestMode(t *testing.T) {
	c := client(t, "sk_live_x", func(w http.ResponseWriter, r *http.Request) { t.Error("Stripe was called with a live key") })
	if _, err := c.Issue(context.Background(), in); !errors.Is(err, ErrLiveKey) {
		t.Errorf("live key: %v", err)
	}
	c = client(t, "rk_test_x", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"spt_1","status":"active","livemode":true}`))
	})
	if _, err := c.Get(context.Background(), "spt_1"); !errors.Is(err, ErrLiveAnswer) {
		t.Errorf("live answer: %v", err)
	}
}

// A token whose limits or seller aren't the ones asked is revoked and
// refused; the asked one is taken, with the preview version.
func TestAnIssuedTokenMustBeTheOneAsked(t *testing.T) {
	var revoked []string
	answer := func(seller string, max int64) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/revoke") {
				revoked = append(revoked, r.URL.Path)
				_, _ = w.Write([]byte(`{"id":"spt_1","status":"deactivated"}`))
				return
			}
			if r.Header.Get("Stripe-Version") != "2026-04-22.preview" {
				t.Errorf("version: %q", r.Header.Get("Stripe-Version"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "spt_1", "status": "active",
				"usage_limits":   map[string]any{"currency": "eur", "max_amount": max, "expires_at": 1798761600},
				"seller_details": map[string]any{"network_business_profile": seller}})
		}
	}
	if _, err := client(t, "rk_test_x", answer("profile_test_other", 1999)).Issue(context.Background(), in); err == nil || len(revoked) != 1 {
		t.Errorf("another seller: %v, revoked %v", err, revoked)
	}
	if _, err := client(t, "rk_test_x", answer("profile_test_1", 5000)).Issue(context.Background(), in); err == nil || len(revoked) != 2 {
		t.Errorf("another cap: %v, revoked %v", err, revoked)
	}
	tok, err := client(t, "rk_test_x", answer("profile_test_1", 1999)).Issue(context.Background(), in)
	if err != nil || tok.ID != "spt_1" {
		t.Fatalf("the asked one: %+v %v", tok, err)
	}
	if strings.Contains(tok.String(), "spt_1") || strings.Contains(strings.ReplaceAll(tok.GoString(), " ", ""), "spt_1") {
		t.Errorf("the token prints its id: %s", tok)
	}
}

// A failure never says the token's id.
func TestAFailureNeverSaysTheID(t *testing.T) {
	c := &Client{Base: "http://127.0.0.1:1", Key: func(context.Context) (string, error) { return "rk_test_x", nil }}
	err := c.Revoke(context.Background(), "spt_secret_1")
	if err == nil || strings.Contains(err.Error(), "spt_secret_1") || !strings.Contains(err.Error(), "{id}") {
		t.Errorf("revoke failure: %v", err)
	}
}
