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
	"time"

	"github.com/kiff/planque/internal/browser"
)

// healthSite is a site whose /account is the fixed page: shown to a signed-in
// visitor, sent to /login otherwise (by redirect, or by script at /account-js).
// /trap tries, on load, to post, beacon and submit a form. writes counts every
// request it received that wasn't a GET or HEAD.
func healthSite(t *testing.T) (*httptest.Server, *atomic.Int64, *atomic.Int64, *atomic.Int64) {
	t.Helper()
	var writes, workers, sockets atomic.Int64
	mux := http.NewServeMux()
	signedIn := func(r *http.Request) bool {
		c, err := r.Cookie("session")
		return err == nil && (c.Value == "ok" || c.Value == "rotated")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			writes.Add(1)
		}
		// A socket's opening reaches the site as a GET with an Upgrade: counted
		// as a write, since its messages would act.
		if r.Header.Get("Upgrade") != "" {
			sockets.Add(1)
			http.Error(w, "no sockets here", http.StatusBadRequest)
			return
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
try{const s=new WebSocket("ws://"+location.host+"/sock");s.onopen=()=>s.send("act")}catch(e){}
try{const f=document.createElement("iframe");document.body.appendChild(f);const W=f.contentWindow.WebSocket;const s=new W("ws://"+location.host+"/sock");s.onopen=()=>s.send("act")}catch(e){}
try{const f=document.createElement("iframe");document.body.appendChild(f);new f.contentWindow.Worker("/worker.js")}catch(e){}
try{const f=document.createElement("iframe");document.body.appendChild(f);new f.contentWindow.Worker(URL.createObjectURL(new Blob(['fetch("/write",{method:"POST"}).catch(()=>{});try{new WebSocket("ws://'+location.host+'/sock")}catch(e){}'])))}catch(e){}
new Worker("/worker.js");navigator.serviceWorker&&navigator.serviceWorker.register("/sw.js").catch(()=>{});
setTimeout(()=>document.getElementById("f").submit(),300)</script></body></html>`)
		case "/worker.js", "/sw.js":
			workers.Add(1)
			w.Header().Set("Content-Type", "text/javascript")
			fmt.Fprint(w, `fetch("/write",{method:"POST",body:"from a worker"}).catch(()=>{});try{new WebSocket("ws://"+location.host+"/sock")}catch(e){}`)
		default:
			http.NotFound(w, r)
		}
	})
	s := httptest.NewServer(mux)
	t.Cleanup(s.Close)
	return s, &writes, &workers, &sockets
}

// signedInState is a browser state holding the site's sign-in cookie.
func signedInState(t *testing.T, s *httptest.Server) []byte {
	t.Helper()
	u, _ := url.Parse(s.URL)
	raw, _ := json.Marshal(map[string]any{"cookies": []map[string]any{{"name": "session", "value": "ok", "domain": u.Hostname(), "path": "/"}}})
	return raw
}

// A health check reads the verdict: signed in where the fixed page shows the
// signed-in words, signed out where it's redirected to a page asking for a
// password, and unknown where it can't tell. It runs no page script, and no
// request that isn't a GET or HEAD leaves the browser: a page that tries to
// post, beacon, open sockets, start workers and submit on load reaches the
// site with none of them.
func TestAHealthCheckReadsTheVerdictAndNeverWrites(t *testing.T) {
	w := worker(t)
	s, writes, workers, sockets := healthSite(t)
	host := func() string { u, _ := url.Parse(s.URL); return u.Hostname() }()
	ctx := context.Background()
	check := func(path string, state []byte) browser.HealthResult {
		t.Helper()
		res, err := w.Health(ctx, browser.HealthCheck{Grant: "grt_test", State: state, Site: host, URL: s.URL + path, Shows: "Your account"})
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
	// A site that sends a signed-out visitor away only by script: the check
	// runs no page script, so it can't follow, and without the signed-in words
	// it reads unknown, never signed in. Signed in, the server's own page shows them.
	if r := check("/account-js", nil); r.Verdict != browser.HealthUnknown {
		t.Fatalf("sent away by script, which a check never runs: %+v", r)
	}
	if r := check("/account-js", signedInState(t, s)); r.Verdict != browser.HealthSignedIn {
		t.Fatalf("signed in, by script-free proof: %+v", r)
	}
	r = check("/trap", signedInState(t, s))
	// No worker started (a check has none: a socket in a worker isn't blocked
	// by address), and no socket reached the site, the one borrowed from a
	// blank frame included.
	if workers.Load() != 0 {
		t.Fatalf("a worker started in a check: %d", workers.Load())
	}
	if writes.Load() != 0 || sockets.Load() != 0 {
		t.Fatalf("%d requests that write and %d sockets reached the site", writes.Load(), sockets.Load())
	}
	// With no page script, the trap's script never runs at all: nothing is
	// even tried. The control (TestTheHealthTrapWorksOutsideACheck) shows it
	// would, and TestAHealthChecksLowerLayersHoldWithoutItsScript that the
	// layers under the script hold by themselves.
	_ = r
	if w.Open() != 0 {
		t.Fatalf("%d sessions left open", w.Open())
	}
}

// The trap is real: in an ordinary session (not a check), the same page's
// socket reaches the site and its worker starts. So the check's zeros are
// the check's.
func TestTheHealthTrapWorksOutsideACheck(t *testing.T) {
	w := worker(t)
	s, _, workers, sockets := healthSite(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, browser.Start{Grant: "grt_test"})
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	if _, err := sess.Open(ctx, s.URL+"/trap"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 30 && (sockets.Load() == 0 || workers.Load() == 0); i++ {
		time.Sleep(100 * time.Millisecond)
	}
	if sockets.Load() == 0 || workers.Load() == 0 {
		t.Fatalf("outside a check, sockets %d and workers %d: the check's test would prove nothing", sockets.Load(), workers.Load())
	}
}

// The layers under the document script hold by themselves: with the script
// left out (so the page keeps WebSocket and Worker), the trap's socket is
// still blocked at the network layer, and its workers are closed at the
// target before they run. Nothing that would write reaches the site.
func TestAHealthChecksLowerLayersHoldWithoutItsScript(t *testing.T) {
	w := worker(t)
	w.healthNoDocScript = true
	s, writes, workers, sockets := healthSite(t)
	host := func() string { u, _ := url.Parse(s.URL); return u.Hostname() }()
	if _, err := w.Health(context.Background(), browser.HealthCheck{Grant: "grt_test", State: signedInState(t, s), Site: host, URL: s.URL + "/trap", Shows: "Your account"}); err != nil {
		t.Fatal(err)
	}
	if writes.Load() != 0 || sockets.Load() != 0 || workers.Load() != 0 {
		t.Fatalf("without the script: %d writes, %d sockets, %d workers reached the site", writes.Load(), sockets.Load(), workers.Load())
	}
}
