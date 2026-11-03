package connector

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

var ErrSignature = errors.New("the event signature is invalid or expired")

// VerifyWebhook authenticates raw bytes before any routing field is trusted.
// The project secret is used as UTF-8, as documented by the provider. Several
// v1 signatures may be supplied during a signing-secret rotation.
func VerifyWebhook(secret, id, timestamp, signature string, body []byte, now time.Time) error {
	if secret == "" || id == "" || len(id) > 256 || len(signature) > 2048 || len(body) > 1<<20 {
		return ErrSignature
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || strconv.FormatInt(ts, 10) != timestamp {
		return ErrSignature
	}
	at := time.Unix(ts, 0)
	if at.Before(now.Add(-5*time.Minute)) || at.After(now.Add(5*time.Minute)) {
		return ErrSignature
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(id + "." + timestamp + "."))
	mac.Write(body)
	for _, candidate := range strings.Fields(signature) {
		encoded, ok := strings.CutPrefix(candidate, "v1,")
		if !ok {
			continue
		}
		got, err := base64.StdEncoding.DecodeString(encoded)
		if err == nil && hmac.Equal(got, mac.Sum(nil)) {
			return nil
		}
	}
	return ErrSignature
}
func SignWebhook(secret, id, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(id + "." + timestamp + "."))
	mac.Write(body)
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Event is the V3 trigger envelope. Tenant, agent and recipient are derived
// from planque's approved subscription, never supplied by this envelope.
type Event struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Metadata struct {
		Trigger string `json:"trigger_id"`
		Slug    string `json:"trigger_slug"`
		Account string `json:"connected_account_id"`
		User    string `json:"user_id"`
	} `json:"metadata"`
	Data      json.RawMessage `json:"data"`
	Timestamp string          `json:"timestamp"`
}

func ParseEvent(raw []byte) (Event, error) {
	var e Event
	if len(raw) > 1<<20 || json.Unmarshal(raw, &e) != nil || e.ID == "" || len(e.ID) > 256 || e.Type != "composio.trigger.message" || e.Metadata.Trigger == "" || e.Metadata.Slug == "" || e.Metadata.Account == "" || e.Metadata.User == "" || !json.Valid(e.Data) {
		return Event{}, ErrDefinition
	}
	return e, nil
}
