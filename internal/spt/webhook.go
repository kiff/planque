package spt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Event is a verified webhook event about an issued token. It only says
// which token to read again (Settle): its own fields decide nothing.
type Event struct {
	ID    string
	Type  string
	Token string
}

// ErrSignature is a webhook whose Stripe signature doesn't check, or is too
// old: it changes nothing.
var ErrSignature = errors.New("spt: the webhook's Stripe signature doesn't check")

// Tolerance is how old a signed webhook may be.
const Tolerance = 5 * time.Minute

// Verify checks a webhook's Stripe-Signature header (t=…,v1=…, an
// HMAC-SHA256 of "t.payload" with the endpoint's secret) and its age, and
// reads the token it names. Only the token events are taken.
func Verify(payload []byte, header, secret string, now time.Time) (Event, error) {
	var ts int64
	var sigs []string
	for _, part := range strings.Split(header, ",") {
		k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		switch k {
		case "t":
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return Event{}, ErrSignature
			}
			ts = n
		case "v1":
			sigs = append(sigs, v)
		}
	}
	if ts == 0 || len(sigs) == 0 || secret == "" {
		return Event{}, ErrSignature
	}
	if age := now.Sub(time.Unix(ts, 0)); age > Tolerance || age < -Tolerance {
		return Event{}, ErrSignature
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(payload)
	want := mac.Sum(nil)
	good := false
	for _, s := range sigs {
		if got, err := hex.DecodeString(s); err == nil && hmac.Equal(got, want) {
			good = true
		}
	}
	if !good {
		return Event{}, ErrSignature
	}
	var e struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID     string `json:"id"`
				Object string `json:"object"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &e); err != nil {
		return Event{}, err
	}
	if !strings.HasPrefix(e.Type, "shared_payment.issued_token.") || e.Data.Object.ID == "" {
		return Event{}, errors.New("spt: not an issued token's event")
	}
	return Event{ID: e.ID, Type: e.Type, Token: e.Data.Object.ID}, nil
}
