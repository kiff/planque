package browserworker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// healthSite is a site whose /account is the fixed page: shown to a signed-in
// visitor, sent to /login otherwise (by redirect, or by script at /account-js).
// /trap tries, on load, to post, beacon and submit a form. writes counts every
// request it received that wasn't a GET or HEAD.
func healthSite(t *testing.T) (*httptest.Server, *atomic.Int64, *atomic.Int64) {
	t.Helper()
	var writes, workers atomic.Int64
	mux := http.NewServeMux()
	signedIn := func(r *http.Request) bool {
		c, err := r.Cookie("session")
		return err == nil && (c.Value == "ok" || c.Value == "rotated")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			writes.Add(1)
		}
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/account":
			if !signedIn(r) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			// The site rotates the sign-in on each visit, as many do.
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "rotated", Path: "/"})
			fmt.Fprint(w, `<html><body><h1>Your account</h1><p>Ana</p></body></html>`)
		case "/account-js":
			if !signedIn(r) {
				fmt.Fprint(w, `<html><body><script>location.href="/login"</script></body></html>`)
				return
			}
			fmt.Fprint(w, `<html><body><h1>Your account</h1></body></html>`)
		case "/login":
			fmt.Fprint(w, `<html><body><form method="post" action="/login"><label>Email <input name="e" type="email"></label><label>Password <input name="p" type="password"></label><button>Sign in</button></form></body></html>`)
		case "/trap":
			fmt.Fprint(w, `<html><body><h1>Your account</h1><form id="f" method="post" action="/write"><input name="x" value="1"></form>
<script>fetch("/write",{method:"POST",body:"x"}).catch(()=>{});navigator.sendBeacon("/write","x");
new Worker("/worker.js");navigator.serviceWorker&&navigator.serviceWorker.register("/sw.js").catch(()=>{});
setTimeout(()=>document.getElementById("f").submit(),300)</script></body></html>`)
		case "/worker.js", "/sw.js":
			workers.Add(1)
			w.Header().Set("Content-Type", "text/javascript")
			fmt.Fprint(w, `fetch("/write",{method:"POST",body:"from a worker"}).catch(()=>{});`)
		default:
			http.NotFound(w, r)
		}
	})
	s := httptest.NewServer(mux)
	t.Cleanup(s.Close)
	return s, &writes, &workers
}

// signedInState is a browser state holding the site's sign-in cookie.
func signedInState(t *testing.T, s *httptest.Server) []byte {
	t.Helper()
	u, _ := url.Parse(s.URL)
	raw, _ := json.Marshal(map[string]any{"cookies": []map[string]any{{"name": "session", "value": "ok", "domain": u.Hostname(), "path": "/"}}})
	return raw
}

// A health check reads the verdict, signed in where the fixed page shows,
// signed out where it's sent to a page asking for a password, whether by
// redirect or by script, and never lets a request that isn't a GET or HEAD
// leave the browser: a page that tries to post, beacon and submit on load
// reaches the site with none of them.
func TestAHealthCheckReadsTheVerdictAndNeverWrites(t *testing.T) {
	w := worker(t)
	s, writes, workers := healthSite(t)
	host := func() string { u, _ := url.Parse(s.URL); return u.Hostname() }()
	ctx := context.Background()
	check := func(path string, state []byte) browser.HealthResult {
		t.Helper()
		res, err := w.Health(ctx, browser.HealthCheck{Grant: "grt_test", State: state, Site: host, URL: s.URL + path})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	r := check("/account", signedInState(t, s))
	if r.Verdict != browser.HealthSignedIn {
		t.Fatalf("signed in: %+v", r)
	}
	// The cookie the visit rotated comes back, to be sealed as a session's is.
	if !strings.Contains(string(r.State), "rotated") {
		t.Fatalf("the rotated sign-in wasn't given back: %s", r.State)
	}
	if r := check("/account", nil); r.Verdict != browser.HealthSignedOut || r.State != nil {
		t.Fatalf("redirected to sign in: %+v", r)
	}
	if r := check("/account-js", nil); r.Verdict != browser.HealthSignedOut {
		t.Fatalf("sent to sign in by script: %+v", r)
	}
	r = check("/trap", signedInState(t, s))
	// The dedicated worker ran (its script was fetched), and its post too was failed.
	if workers.Load() == 0 {
		t.Fatal("the trap's worker never started: the test proves nothing about workers")
	}
	if writes.Load() != 0 {
		t.Fatalf("%d requests that write reached the site", writes.Load())
	}
	if r.Blocked == 0 {
		t.Fatalf("the trap's writes weren't failed: %+v", r)
	}
	if w.Open() != 0 {
		t.Fatalf("%d sessions left open", w.Open())
	}
}
