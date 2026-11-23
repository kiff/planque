package planque

import (
	"strings"
	"testing"
	"time"
)

// tafs#1746's extension sends profile_id: planque took it as an unknown
// field and refused the pairing.
// The extension's offer carries its profile id and whether it checks a
// one-touch phone_sig ("Your Chrome" 1 and 5): planque keeps both as given,
// the phone reads them, and the paired Chrome lists its profile id, so a
// later re-pairing can pre-select it. A malformed id is refused, and an
// extension that sends neither pairs as before.
func TestTheOfferSaysItsProfileAndOneTouchToThePhone(t *testing.T) {
	e := chromeEnv(t)
	acct := "/v1/accounts/" + tenant
	open := func() (string, string) {
		// The phone signs each request anew: the same one in the same second is a replay.
		time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
		st, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, "")
		if st != 201 {
			t.Fatalf("open: %d %v", st, out)
		}
		return out["pairing"].(string), out["lookup"].(string)
	}
	post := func(lookup, key string, extra map[string]any) (int, map[string]any) {
		body := map[string]any{"key": key, "version": "1.0.0", "manifest_hash": releasedManifest, "browser": "Chrome 131", "os": "macOS 15", "proof": b64(32)}
		for k, v := range extra {
			body[k] = v
		}
		return e.do("POST", "/v1/extension/pairings/"+lookup, body, nil, "")
	}
	profile := strings.Repeat("ab", 16)

	_, lookup := open()
	key, _ := p256Key(t)
	for what, extra := range map[string]map[string]any{
		"a malformed profile id":     {"profile_id": "not-hex"},
		"one_touch with no proof":    {"one_touch": true},
		"a proof with no flag":       {"one_touch_proof": b64(32)},
		"a proof that isn't an HMAC": {"one_touch": true, "one_touch_proof": b64(8)},
	} {
		if st, out := post(lookup, key, extra); st != 400 {
			t.Fatalf("%s: %d %v", what, st, out)
		}
	}

	pairing, lookup := open()
	key, sum := p256Key(t)
	if st, out := post(lookup, key, map[string]any{"profile_id": profile, "one_touch": true, "one_touch_proof": b64(32)}); st != 202 {
		t.Fatalf("post: %d %v", st, out)
	}
	st, out := e.do("GET", acct+"/chromes/pairings/"+pairing, nil, e.phone, "")
	if st != 200 || out["profile_id"] != profile || out["one_touch"] != true || out["one_touch_proof"] != b64(32) {
		t.Fatalf("the phone's read: %d %v", st, out)
	}
	if st, out := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", map[string]any{"chrome_key_sha256": sum, "name": "Work Mac",
		"phone_key": key, "phone_proof": b64(32), "phone_sig": b64(64)}, e.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	st, out = e.do("GET", acct+"/chromes", nil, e.phone, "")
	if list := out["chromes"].([]any); st != 200 || len(list) != 1 || list[0].(map[string]any)["profile_id"] != profile {
		t.Fatalf("the list: %d %v", st, out)
	}

	// An older extension: neither field, and the read says neither.
	pairing, lookup = open()
	key, _ = p256Key(t)
	if st, out := post(lookup, key, nil); st != 202 {
		t.Fatalf("an older extension's post: %d %v", st, out)
	}
	if _, out := e.do("GET", acct+"/chromes/pairings/"+pairing, nil, e.phone, ""); out["profile_id"] != nil || out["one_touch"] != nil || out["one_touch_proof"] != nil {
		t.Fatalf("an older extension's read: %v", out)
	}
}
