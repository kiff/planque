package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"testing"

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
}
