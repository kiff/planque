package planque

import (
	"bytes"
	"context"
	"log/slog"
	"net/url"
	"strings"
	"testing"
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
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil || g.Status != StatusPending {
		t.Fatalf("grant: %v %+v", err, g)
	}
}
