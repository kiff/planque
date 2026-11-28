package planque

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// signInLink approves grant id on phone and returns the vendor sign-in link.
func (e *env) approveFor(id, digest string) string {
	e.t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 200 {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	return out["sign_in_url"].(string)
}

// On Android, the sign-in's end is a page that opens tafs with a button
// (and tries the app's link once at load), never a bare redirect: Chrome
// drops a redirect to an app with no tap behind it, as after the owner
// confirmed in their bank's own app (seen live: the code never reached tafs).
// The button carries the one-use code for this grant; only the page's own
// script runs; it can't be framed or cached; loading the callback again
// shows no code; and the code confirms nothing without the device's
// signature.
func TestAnAndroidSignInEndsOnAButtonThatOpensTafs(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	d := e.kiff.Devices["dev_phone"]
	d.Platform = "android"
	e.kiff.Devices["dev_phone"] = d
	id, digest := e.request()
	link := e.approveFor(id, digest)
	client := e.client
	if client == nil {
		client = http.DefaultClient
	}
	var last *url.URL
	c := *client
	c.CheckRedirect = func(r *http.Request, via []*http.Request) error { last = r.URL; return nil }
	resp, err := c.Get(link)
	if err != nil {
		t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || string(page) == "app opened" {
		t.Fatalf("the end was a redirect to the app, not a page: %d %s", resp.StatusCode, page)
	}
	h := resp.Header
	csp := h.Get("Content-Security-Policy")
	scripts := csp[strings.Index(csp, "script-src"):]
	if !strings.Contains(scripts, "'"+openScriptHash+"'") || strings.Contains(scripts, "unsafe") || !strings.Contains(csp, "frame-ancestors 'none'") ||
		h.Get("Cache-Control") != "no-store" || h.Get("Referrer-Policy") != "no-referrer" || h.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("headers: %v", h)
	}
	if !strings.Contains(string(page), "<script>"+openScript+"</script>") || strings.Contains(string(page), "http://") && strings.Contains(string(page), "<link") {
		t.Fatalf("the page: %s", page)
	}
	m := regexp.MustCompile(`<a id="open" href="([^"]+)">`).FindStringSubmatch(string(page))
	if m == nil {
		t.Fatalf("no button: %s", page)
	}
	u, err := url.Parse(strings.ReplaceAll(m[1], "&amp;", "&"))
	if err != nil || u.Query().Get("grant") != id || u.Query().Get("confirm") == "" {
		t.Fatalf("the button's link: %q", m[1])
	}
	code := u.Query().Get("confirm")
	// Loading the callback again: the vendor's code is spent, so no code.
	if last != nil {
		again, err := client.Get(last.String())
		if err == nil {
			body, _ := io.ReadAll(again.Body)
			again.Body.Close()
			if strings.Contains(string(body), code) {
				t.Fatal("a second load of the callback showed the code")
			}
		}
	}
	// Without the device's signature, the code confirms nothing.
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, nil, ""); st == 200 {
		t.Fatal("confirmed with no device signature")
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm with the button's code: %d %v", st, out)
	}
}

// signedInGrant approves a request on e.phone and finishes its sign-in at
// the fake vendor, returning the grant (signed_in) and the code that went
// to the browser.
func (e *env) signedInGrant() (*Grant, string) {
	e.t.Helper()
	id, digest := e.request()
	code := e.signIn(e.approveFor(id, digest))
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil || g.Status != StatusSignedIn {
		e.t.Fatalf("signed in: %v %+v", err, g)
	}
	return g, code
}

// approveSigned approves g again on p, signing the words it shows (digest
// v2), or the request alone (v1) when words is false.
func (e *env) approveAgain(g *Grant, p *fakes.Phone, words bool) (int, map[string]any) {
	e.t.Helper()
	sum := g.Digest
	if words {
		var err error
		if sum, err = g.digestIn("en"); err != nil {
			e.t.Fatal(err)
		}
	}
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": sum}, p, "")
}

// The device that approved approves again while the sign-in waits for its
// confirm (its code never reached it), signing words that name the account
// the sign-in is in: the sign-in is kept, and a new code comes in the signed
// answer; the old code no longer confirms, the new one does.
func TestTheApprovingDeviceResumesASignInThatNamesItsAccount(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	g, old := e.signedInGrant()
	if g.Account == "" || !strings.Contains(g.Text, g.Account) {
		t.Fatalf("the words don't name the account: %q / %q", g.Account, g.Text)
	}
	st, out := e.approveAgain(g, e.phone, true)
	code, _ := out["confirm"].(string)
	if st != 200 || out["status"] != StatusSignedIn || out["account"] != g.Account || code == "" || code == old || out["sign_in_url"] != nil {
		t.Fatalf("resumed: %d %v", st, out)
	}
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/confirm", map[string]string{"confirm": old}, e.phone, ""); st == 200 {
		t.Fatal("the old code still confirmed")
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("the new code: %d %v", st, out)
	}
}

// No resume without the account in what was signed: a sign-in with no
// account, an approval that signed the request alone (no words), or another
// device all get today's revoke and a new link.
func TestASignInIsResumedOnlyWithItsAccountSignedByTheSameDevice(t *testing.T) {
	t.Run("no account", func(t *testing.T) {
		e := newEnv(t)
		e.enrol()
		e.vendor.UnnamedAccount = true
		g, _ := e.signedInGrant()
		if g.Account != "" {
			t.Skipf("the fake named an account anyway: %q", g.Account)
		}
		st, out := e.approveAgain(g, e.phone, true)
		if st != 200 || out["confirm"] != nil || out["sign_in_url"] == nil {
			t.Fatalf("no account: %d %v", st, out)
		}
	})
	t.Run("words not signed", func(t *testing.T) {
		e := newEnv(t)
		e.enrol()
		g, _ := e.signedInGrant()
		// The same body as the first approval: past the second in which an
		// identical signed request is a replay, as a phone's later tap is.
		time.Sleep(1100 * time.Millisecond)
		st, out := e.approveAgain(g, e.phone, false)
		if st != 200 || out["confirm"] != nil || out["sign_in_url"] == nil {
			t.Fatalf("request alone: %d %v", st, out)
		}
	})
	t.Run("another device", func(t *testing.T) {
		e := newEnv(t)
		e.enrol()
		g, _ := e.signedInGrant()
		st, out := e.approveAgain(g, e.device("dev_ipad"), true)
		if st != 200 || out["confirm"] != nil || out["sign_in_url"] == nil {
			t.Fatalf("another device: %d %v", st, out)
		}
	})
}
