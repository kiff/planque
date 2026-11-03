package connector

import (
	"strings"
	"testing"
	"time"
)

func TestWebhookAuthenticatesExactBytesAndRejectsReplayWindow(t *testing.T) {
	now := time.Unix(1800000000, 0)
	body := []byte(`{"tenant":"untrusted"}`)
	ts := "1800000000"
	sig := SignWebhook("synthetic-secret", "event_one", ts, body)
	if err := VerifyWebhook("synthetic-secret", "event_one", ts, sig, body, now); err != nil {
		t.Fatal(err)
	}
	if err := VerifyWebhook("synthetic-secret", "event_one", ts, "v1,Zm9yZ2Vk "+sig, body, now); err != nil {
		t.Fatal("rotation signature rejected")
	}
	for _, scenario := range []string{"missing", "bad_secret", "changed_id", "changed_timestamp", "changed_body", "expired", "future", "wrong_version", "noncanonical_timestamp"} {
		secret, id, timestamp, signature, raw, at := "synthetic-secret", "event_one", ts, sig, body, now
		switch scenario {
		case "missing":
			secret = ""
		case "bad_secret":
			secret = "another"
		case "changed_id":
			id = "event_other"
		case "changed_timestamp":
			timestamp = "1800000001"
		case "changed_body":
			raw = []byte(`{"tenant":"another"}`)
		case "expired":
			at = now.Add(301 * time.Second)
		case "future":
			at = now.Add(-301 * time.Second)
		case "wrong_version":
			signature = strings.Replace(sig, "v1,", "v2,", 1)
		case "noncanonical_timestamp":
			timestamp = "01800000000"
		}
		if VerifyWebhook(secret, id, timestamp, signature, raw, at) == nil {
			t.Fatalf("accepted %s", scenario)
		}
	}
}
