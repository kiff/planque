package planque

import (
	"cmp"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

// noFollow is a browser that shows each answer instead of following it.
var noFollow = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func visit(t *testing.T, link string) (*http.Response, string) {
	t.Helper()
	resp, err := noFollow.Get(link)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, string(b)
}

// signInAt follows a sign-in link through the vendor's page to planque's
// callback, and returns the callback's own answer, unfollowed, and its URL.
func (e *env) signInAt(t *testing.T, link string, lang string) (*http.Response, string, string) {
	t.Helper()
	callback := ""
	c := &http.Client{CheckRedirect: func(r *http.Request, _ []*http.Request) error {
		if strings.HasPrefix(r.URL.String(), e.app.URL) {
			return http.ErrUseLastResponse
		}
		if strings.HasPrefix(r.URL.String(), e.url) {
			callback = r.URL.String()
		}
		return nil
	}}
	req, _ := http.NewRequest("GET", link, nil)
	if lang != "" {
		req.Header.Set("Accept-Language", lang)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if callback == "" {
		t.Fatalf("the sign-in never reached planque's callback: %d", resp.StatusCode)
	}
	return resp, string(b), callback
}

var openLink = regexp.MustCompile(`<a id="open" href="([^"]+)">`)

// A finished sign-in goes back to the app the way the approving device's
// browser can take it: an iPhone's or iPad's by a redirect its browser
// session catches, as before, byte for byte; an Android phone's by a page
// with an Open tafs link that also tries it once at load (28 Nov: after the
// owner confirmed in their bank's own app, Chrome dropped the redirect, four
// times); a computer's, or one whose platform KIFF didn't name, by the same
// page without the try, since a desktop browser drops a redirect to an app
// with no click behind it. Both
// are never cached nor sent onward, and the link works only once: a second
// visit shows no link (17 Nov: Windows' Chrome dropped the redirect).
func TestASignInReturnsToTheAppTheWayItsDeviceCanTakeIt(t *testing.T) {
	for _, tc := range []struct {
		platform string
		redirect bool
		// tries: the page also follows its own link once at load (Android,
		// whose browser drops a redirect to an app with no tap behind it).
		tries bool
	}{{"ios", true, false}, {"ipados", true, false}, {"android", false, true}, {"windows", false, false}, {"macos", false, false}, {"linux", false, false}, {"", false, false}} {
		t.Run(cmp.Or(tc.platform, "unnamed"), func(t *testing.T) {
			e := newEnv(t)
			e.enrol()
			dev := e.device("dev_" + cmp.Or(tc.platform, "unnamed"))
			d := e.kiff.Devices[dev.DeviceID]
			d.Platform = tc.platform
			e.kiff.Devices[dev.DeviceID] = d
			id, digest := e.request()
			st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, dev, "")
			if st != 200 {
				t.Fatalf("approve: %d %v", st, out)
			}
			resp, page, link := e.signInAt(t, out["sign_in_url"].(string), "")
			if got := resp.Header.Get("Cache-Control"); got != "no-store" {
				t.Errorf("Cache-Control %q", got)
			}
			if got := resp.Header.Get("Referrer-Policy"); got != "no-referrer" {
				t.Errorf("Referrer-Policy %q", got)
			}
			var back string
			if tc.redirect {
				back = resp.Header.Get("Location")
				if resp.StatusCode != http.StatusFound || !strings.HasPrefix(back, e.app.URL+"/signed-in?") {
					t.Fatalf("a phone's sign-in: %d, Location %q", resp.StatusCode, back)
				}
			} else {
				if resp.StatusCode != 200 || resp.Header.Get("Location") != "" {
					t.Fatalf("a computer's sign-in: %d, Location %q", resp.StatusCode, resp.Header.Get("Location"))
				}
				m := openLink.FindStringSubmatch(page)
				if m == nil || !strings.Contains(page, ">Open tafs</a>") {
					t.Fatalf("no Open tafs link:\n%s", page)
				}
				back = html.UnescapeString(m[1])
				if !strings.HasPrefix(back, e.app.URL+"/signed-in?") {
					t.Fatalf("the link goes to %q", back)
				}
				csp := resp.Header.Get("Content-Security-Policy")
				if !strings.Contains(csp, "default-src 'none'") || !strings.Contains(csp, "frame-ancestors 'none'") || resp.Header.Get("X-Frame-Options") != "DENY" {
					t.Errorf("CSP %q, frames %q", csp, resp.Header.Get("X-Frame-Options"))
				}
				if tc.tries != strings.Contains(page, "<script>"+openScript+"</script>") || tc.tries != strings.Contains(csp, "'"+openScriptHash+"'") {
					t.Errorf("the load-time try: want %v; CSP %q\n%s", tc.tries, csp, page)
				}
			}
			for _, k := range []string{"tenant=", "grant=" + id, "confirm="} {
				if !strings.Contains(back, k) {
					t.Errorf("the return link has no %s: %q", k, back)
				}
			}
			// The person's click (or the phone's session) reaches the app,
			// and the approving device confirms with the code.
			if _, body := visit(t, back); body != "app opened" {
				t.Fatalf("the app: %q", body)
			}
			if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": e.app.code()}, dev, ""); st != 200 || out["status"] != StatusConnected {
				t.Fatalf("confirm: %d %v", st, out)
			}
			// A second visit: no link, and no code.
			resp, page = visit(t, link)
			if resp.Header.Get("Location") != "" || strings.Contains(page, e.app.URL) || strings.Contains(page, "confirm=") {
				t.Fatalf("a second visit hands the app's link again: %d %q\n%s", resp.StatusCode, resp.Header.Get("Location"), page)
			}
		})
	}
}

// Signing in again on a connected grant returns the same way, by the
// device that approved the new sign-in.
func TestANewSignInReturnsToTheAppTheWayItsDeviceCanTakeIt(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	pc := e.device("dev_pc")
	d := e.kiff.Devices["dev_pc"]
	d.Platform = "windows"
	e.kiff.Devices["dev_pc"] = d
	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": g.Digest}, pc, "")
	if st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	resp, page, _ := e.signInAt(t, out["sign_in_url"].(string), "")
	if resp.StatusCode != 200 || resp.Header.Get("Referrer-Policy") != "no-referrer" || resp.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("a computer's new sign-in: %d %v", resp.StatusCode, resp.Header)
	}
	if m := openLink.FindStringSubmatch(page); m == nil || !strings.Contains(html.UnescapeString(m[1]), "grant="+g.ID) {
		t.Fatalf("no Open tafs link for the new sign-in:\n%s", page)
	}
}

// The page is in the person's language, with the app's name.
func TestTheOpenTafsPageSpeaksThePersonsLanguage(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	pc := e.device("dev_pc")
	d := e.kiff.Devices["dev_pc"]
	d.Platform = "windows"
	e.kiff.Devices["dev_pc"] = d
	id, digest := e.request()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, pc, "")
	_, page, _ := e.signInAt(t, out["sign_in_url"].(string), "es-ES,es;q=0.9")
	if !strings.Contains(page, `lang="es"`) || !strings.Contains(page, ">Abrir tafs</a>") {
		t.Fatalf("not in Spanish:\n%s", page)
	}
}
