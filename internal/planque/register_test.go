package planque

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/oauth"
)

// A vendor that doesn't allow planque's client (Qonto allow-lists the
// redirect URIs of the MCP clients it lets register) refuses it at
// registration. planque says so in one log line, with the vendor's host,
// status and error code, and answers in words with a code the apps can
// read, instead of a bare 502 the apps can only call "not on our side".
// Nothing is connected, and the grant waits as it was.
func TestAVendorThatRefusesPlanquesClientIsLoggedAndSaidSo(t *testing.T) {
	e := newEnv(t)
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))
	e.vendor.Redirects = []string{"https://another-client.example/oauth/callback"}
	e.vendor.RefusalWords = "redirect_uri not in the allow-list for this client: write to the vendor's partners desk"
	e.enrol()
	id, digest := e.request()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 502 || out["code"] != "vendor_refused_client" {
		t.Fatalf("approve: %d %v", st, out)
	}
	u, _ := url.Parse(e.vendor.Issuer())
	if msg, _ := out["error"].(string); !strings.Contains(msg, u.Host+" doesn't allow tafs to sign in yet") || !strings.Contains(msg, "nothing was connected") {
		t.Fatalf("words: %q", msg)
	}
	line := logs.String()
	if !strings.Contains(line, `"msg":"client registration refused"`) || !strings.Contains(line, `"status":400`) ||
		!strings.Contains(line, `"code":"invalid_redirect_uri"`) || !strings.Contains(line, `"grant":"`+id+`"`) {
		t.Fatalf("log: %s", line)
	}
	// Never the vendor's own words: not to the owner, not in the log.
	if strings.Contains(line, "allow-list") || strings.Contains(fmt.Sprint(out), "allow-list") {
		t.Fatalf("the vendor's description was repeated: %s %v", line, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil || g.Status != StatusPending {
		t.Fatalf("grant: %v %+v", err, g)
	}
}

// A busy vendor (429), or one that timed out (408), hasn't refused the
// client: the answer is the retryable 502, never "it doesn't allow tafs".
func TestABusyVendorIsNoRefusal(t *testing.T) {
	for _, status := range []int{408, 429} {
		e := newEnv(t)
		e.vendor.RegisterStatus = status
		e.enrol()
		id, digest := e.request()
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
		if st != 502 || out["code"] == "vendor_refused_client" || strings.Contains(fmt.Sprint(out["error"]), "doesn't allow tafs") {
			t.Fatalf("%d: approve: %d %v", status, st, out)
		}
	}
	for status, refused := range map[int]bool{400: true, 401: true, 403: true, 408: false, 429: false, 500: false, 503: false} {
		if got := (&oauth.RegistrationError{Status: status}).Refused(); got != refused {
			t.Errorf("%d: refused %v, want %v", status, got, refused)
		}
	}
}
