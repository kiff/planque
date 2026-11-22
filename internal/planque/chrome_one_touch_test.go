package planque

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// One fingerprint pairs a Chrome (kiff/tafs docs/chrome.md, "Pairing"): a
// phone whose key signs once per prompt signs the pairing's confirm here and
// KIFF's view-only device for the extension (POST /v1/me/pairings) together,
// under one kiff.requests.v1 signature by its approval key. planque checks its
// own request with KIFF's hash as given; a changed body, a device key's
// combined signature or a request paired with itself confirms nothing.
// Category: authority
// Attack: the confirm's body changed after signing, the device key's combined signature, a request paired with itself
// Refused: each; the approval key's confirm goes through once, and the Chrome is kept
func TestOneFingerprintConfirmsAChromeBesideKIFFsPairing(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	approval := fakes.NewPhone()
	approval.DeviceID = "dev_phone"
	d := e.kiff.Devices["dev_phone"]
	d.Approval = approval
	e.kiff.Devices["dev_phone"] = d
	if st, _ := e.enrolAs(phone); st != 201 {
		t.Fatal("enrol")
	}
	e.s.ChromeBuilds = map[string]ChromeBuild{releasedManifest: {ManifestSHA256: releasedManifest, Version: "1.0.0"}}
	e.s.ChromeOffer = ChromeOffer{All: true}

	acct := "/v1/accounts/" + tenant
	st, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, phone, "")
	if st != 201 {
		t.Fatalf("open: %d %v", st, out)
	}
	pairing, lookup := out["pairing"].(string), out["lookup"].(string)
	extKey, extSum := p256Key(t)
	if st, out := extensionPosts(e, lookup, extKey, releasedManifest); st != 202 {
		t.Fatalf("the extension's key: %d %v", st, out)
	}
	phoneKey, _ := p256Key(t)
	confirm := map[string]any{"chrome_key_sha256": extSum, "name": "Work Mac", "phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)}
	body, _ := json.Marshal(confirm)
	path := acct + "/chromes/pairings/" + pairing + "/confirm"

	// KIFF's request the phone signs beside it: its kiff.request.v1 message,
	// at the same signed_at, as KIFF rebuilds it.
	at := time.Now().UTC().Format(time.RFC3339)
	approval.Now = func() time.Time { t, _ := time.Parse(time.RFC3339, at); return t }
	kiffBody := []byte(`{"offer":"…the extension's KIFF offer…","approval_pub":""}`)
	kiffSum := sha256.Sum256(kiffBody)
	kiffMsg := strings.Join([]string{"kiff.request.v1", tenant, "POST", "/v1/me/pairings", hex.EncodeToString(kiffSum[:]), at}, "\n")
	other := sha256.Sum256([]byte(kiffMsg))
	with := hex.EncodeToString(other[:])

	var answered string
	send := func(signer *fakes.Phone, signed, sent []byte) int {
		t.Helper()
		req, _ := http.NewRequest("POST", e.url+path, bytes.NewReader(sent))
		signer.Sign(req, tenant, signed)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		answered = string(b)
		return resp.StatusCode
	}

	// The device key's combined signature: refused, as on every route.
	phone.With = with
	if st := send(phone, body, body); st != http.StatusForbidden {
		t.Fatalf("the device key's combined signature: %d", st)
	}
	phone.With = ""
	// The body changed after the one fingerprint (another name): refused.
	approval.With = with
	renamed, _ := json.Marshal(map[string]any{"chrome_key_sha256": extSum, "name": "Not mine", "phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)})
	if st := send(approval, body, renamed); st != http.StatusForbidden && st != http.StatusUnauthorized {
		t.Fatalf("a body other than the one signed: %d", st)
	}
	// Paired with itself: refused.
	approval.With, approval.WithItself = "", true
	if st := send(approval, body, body); st != http.StatusUnauthorized {
		t.Fatalf("a request paired with itself: %d", st)
	}
	approval.WithItself = false
	if st, out := e.do("GET", acct+"/chromes", nil, phone, ""); st != 200 || len(out["chromes"].([]any)) != 0 {
		t.Fatalf("a Chrome kept by a refused confirm: %d %v", st, out)
	}

	// The approval key's, beside KIFF's request: the Chrome is confirmed,
	// and the extension collects the phone's key.
	approval.With = with
	if st := send(approval, body, body); st != http.StatusOK {
		t.Fatalf("the one fingerprint's confirm: %d", st)
	}
	st, out = e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+extKey, nil, nil, "")
	if st != 200 || out["phone_key"] != phoneKey || out["name"] != "Work Mac" {
		t.Fatalf("collect: %d %v", st, out)
	}
	// The same signed confirm again: refused as a replay, before the
	// pairing's own state is read (qa-kiff: "already confirmed" alone would
	// hide a missing replay check).
	if st := send(approval, body, body); st != http.StatusConflict || !strings.Contains(answered, "signature_replayed") {
		t.Fatalf("the combined signature used twice: %d %s", st, answered)
	}
}
