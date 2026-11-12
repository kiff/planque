package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/fakes"
)

// One fingerprint for two requests (kiff.requests.v1, kiff/tafs
// docs/secure-connection.md): the approval key signs this request's message
// beside another's, once. planque checks its own request with the other's
// hash as given, only from the approval key, and takes the signature once.
// Category: authority
// Attack: a combined signature from the device key (a window's tap), a request paired with itself, a combined signature used twice
// Refused: each; the approval key's goes through once, and its record says which message was signed
func TestARequestSignedTogetherWithAnother(t *testing.T) {
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
	var signed string
	other := sha256.Sum256([]byte("KIFF's confirm of the Card's draft"))
	with := hex.EncodeToString(other[:])

	// The device key's combined signature: refused, nothing approved.
	id, digest := e.request()
	approve := "/v1/accounts/" + tenant + "/grants/" + id + "/approve"
	phone.With = with
	if st, out := e.do("POST", approve, map[string]string{"digest": digest}, phone, ""); st != http.StatusForbidden {
		t.Fatalf("the device key's combined signature: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Approval != nil {
		t.Fatalf("approved by the device key: %+v", g.Approval)
	}

	// The approval key's: through, and the record holds the combined message.
	approval.With = with
	if st, out := e.do("POST", approve, map[string]string{"digest": digest}, approval, ""); st != http.StatusOK {
		t.Fatalf("the approval key's combined signature: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Approval == nil || g.Approval.Key != "critical" || !strings.HasPrefix(g.Approval.Message, "kiff.requests.v1\n"+tenant+"\n") || !strings.Contains(g.Approval.Message, with) {
		t.Fatalf("the record: %+v", g.Approval)
	}

	// A request paired with itself: refused, nothing approved.
	id2, digest2 := e.request()
	approval.With, approval.WithItself = "", true
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id2+"/approve", map[string]string{"digest": digest2}, approval, ""); st != http.StatusUnauthorized {
		t.Fatalf("a request paired with itself: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id2); g.Approval != nil {
		t.Fatalf("approved paired with itself: %+v", g.Approval)
	}
	approval.WithItself = false

	// The key route, as the sheet uses it: the key and its hash in the signed body, given once.
	v := newRESTVendor(t)
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	keyID, keyDigest := out["id"].(string), out["digest"].(string)
	sum := sha256.Sum256([]byte(v.key))
	b, _ := json.Marshal(map[string]string{"digest": keyDigest, "key_sha256": hex.EncodeToString(sum[:])})
	approval.With = with
	send := func(header string) int {
		t.Helper()
		req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tenant+"/grants/"+keyID+"/key", bytes.NewReader(b))
		if header == "" {
			approval.Sign(req, tenant, b)
		} else {
			req.Header.Set(device.Header, header)
		}
		req.Header.Set(KeyHeader, v.key)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if header == "" {
			signed = req.Header.Get(device.Header)
		}
		return resp.StatusCode
	}
	if st := send(""); st != http.StatusOK {
		t.Fatalf("the key, signed together with the Card's confirm: %d", st)
	}
	kg, _ := e.s.loadGrant(context.Background(), tenant, keyID)
	if kg.Approval == nil || kg.Approval.Key != "critical" || !strings.HasPrefix(kg.Approval.Message, "kiff.requests.v1\n") || kg.Status != StatusListed {
		t.Fatalf("the key's record: %+v %s", kg.Approval, kg.Status)
	}
	// The same signed request again: refused, it was used.
	if st := send(signed); st != http.StatusConflict {
		t.Fatalf("the combined signature used twice: %d", st)
	}
}
