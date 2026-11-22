package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/stepkey"
)

// The extension pairs from its own origin with no host permission for
// planque, so Chrome asks before each call (22 Nov, the owner's Chrome: a 405
// on the preflight, and nothing reached planque). Only an extension's origin
// is let read the answers: never `*`, never with credentials.
func TestOnlyAnExtensionsOriginReadsTheExtensionRoutes(t *testing.T) {
	const ext = "chrome-extension://abcdefghijklmnopabcdefghijklmnop"
	e := chromeEnv(t)
	h := e.s.Handler()
	call := func(method, path, origin, body string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		if method == http.MethodOptions {
			r.Header.Set("Access-Control-Request-Method", "POST")
			r.Header.Set("Access-Control-Request-Headers", "content-type")
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	noCredentials := func(name string, w *httptest.ResponseRecorder) {
		t.Helper()
		if v := w.Header().Get("Access-Control-Allow-Credentials"); v != "" {
			t.Errorf("%s: Allow-Credentials %q", name, v)
		}
		if !strings.Contains(strings.Join(w.Header().Values("Vary"), ","), "Origin") {
			t.Errorf("%s: no Vary: Origin", name)
		}
	}

	// The preflight, for both routes: 204, the extension's origin echoed, GET and POST with a JSON body.
	for _, path := range []string{"/v1/extension/pairings/2222", "/v1/extension/pairings/2222?key=x"} {
		w := call(http.MethodOptions, path, ext, "")
		if w.Code != http.StatusNoContent || w.Header().Get("Access-Control-Allow-Origin") != ext ||
			w.Header().Get("Access-Control-Allow-Methods") != "GET, POST" || w.Header().Get("Access-Control-Allow-Headers") != "content-type" {
			t.Fatalf("preflight %s: %d %v", path, w.Code, w.Header())
		}
		noCredentials("preflight", w)
	}
	// The calls themselves, a refusal included, say the same origin.
	post := call(http.MethodPost, "/v1/extension/pairings/2222", ext, `{}`)
	if post.Header().Get("Access-Control-Allow-Origin") != ext || post.Code == http.StatusMethodNotAllowed {
		t.Fatalf("post: %d %v", post.Code, post.Header())
	}
	noCredentials("post", post)
	get := call(http.MethodGet, "/v1/extension/pairings/2222?key=x", ext, "")
	if get.Header().Get("Access-Control-Allow-Origin") != ext {
		t.Fatalf("collect: %d %v", get.Code, get.Header())
	}
	// CORS relaxes nothing: a body with no proof is refused with the extension's origin as without one, and so is a
	// collect with the wrong key.
	if bare := call(http.MethodPost, "/v1/extension/pairings/2222", "", `{}`); post.Code != bare.Code || post.Code < 400 {
		t.Fatalf("a post with no proof: %d with the origin, %d without", post.Code, bare.Code)
	}
	if bare := call(http.MethodGet, "/v1/extension/pairings/2222?key=x", "", ""); get.Code != bare.Code || get.Code != http.StatusNotFound {
		t.Fatalf("a collect with the wrong key: %d with the origin, %d without", get.Code, bare.Code)
	}
	// The preflight writes nothing and runs no handler: it answers before any count, so a burst of them takes no try.
	for range chromePostsPerAddress + 1 {
		call(http.MethodOptions, "/v1/extension/pairings/2222", ext, "")
	}
	if w := call(http.MethodPost, "/v1/extension/pairings/2222", ext, `{}`); w.Code == http.StatusTooManyRequests {
		t.Fatal("preflights took the post's tries")
	}

	// Any other origin, or none, gets no CORS answer: a web page, another scheme, an id not Chrome's shape, `null`.
	for _, origin := range []string{"https://evil.example", "moz-extension://abcdefghijklmnopabcdefghijklmnop",
		"chrome-extension://abcdefghijklmnopabcdefghijklmnoz", "chrome-extension://abcdefghijklmnopabcdefghijklmnop.evil.example",
		"chrome-extension://short", "null", ""} {
		for _, w := range []*httptest.ResponseRecorder{call(http.MethodOptions, "/v1/extension/pairings/2222", origin, ""),
			call(http.MethodPost, "/v1/extension/pairings/2222", origin, `{}`)} {
			if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
				t.Errorf("origin %q: Allow-Origin %q", origin, v)
			}
			if v := w.Header().Get("Access-Control-Allow-Methods"); v != "" {
				t.Errorf("origin %q: Allow-Methods %q", origin, v)
			}
		}
	}

	// Only the extension's routes: the phone's own routes answer no preflight.
	if w := call(http.MethodOptions, "/v1/accounts/tnt/chromes/pairings", ext, ""); w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("the phone's route answered a preflight: %v", w.Header())
	}
}

// The extension unpairs itself from its side panel (kiff/tafs sw.js
// unpairHere): POST /v1/extension/unpair {chrome_id, at, sig}, signed by the
// Chrome's own key. Only that key's fresh signature unpairs it; a wrong key,
// another body, an old time or a replay is refused, and the answer is the
// same 404 as for a Chrome that doesn't exist.
func TestTheExtensionUnpairsItselfWithItsOwnKeyOnly(t *testing.T) {
	const ext = "chrome-extension://abcdefghijklmnopabcdefghijklmnop"
	ctx := context.Background()
	r := newRelayRig(t)
	h := r.e.s.Handler()
	now := r.e.s.now()
	other := stepkey.Local{Key: func() *ecdsa.PrivateKey { k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader); return k }()}
	signed := func(key stepkey.Local, chrome string, at int64, signedAt int64) string {
		sig, err := stepkey.SignBody(ctx, key, labelUnpair, map[string]any{"chrome_id": chrome, "at": signedAt})
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(map[string]any{"chrome_id": chrome, "at": at, "sig": sig})
		return string(b)
	}
	post := func(body string) *httptest.ResponseRecorder {
		q := httptest.NewRequest(http.MethodPost, "/v1/extension/unpair", strings.NewReader(body))
		q.Header.Set("Origin", ext)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, q)
		return w
	}
	paired := func() bool {
		_, err := getDoc(ctx, r.e.s.KV, tenantPK(tenant), skChrome(r.chrome), &Chrome{})
		return err == nil
	}
	ms := now.UnixMilli()

	for name, c := range map[string]struct {
		body string
		want int
	}{
		"another key":           {signed(other, r.chrome, ms, ms), http.StatusNotFound},
		"another body signed":   {signed(r.ext, r.chrome, ms, ms+1), http.StatusNotFound},
		"an old time":           {signed(r.ext, r.chrome, now.Add(-6*time.Minute).UnixMilli(), now.Add(-6*time.Minute).UnixMilli()), http.StatusBadRequest},
		"a time ahead":          {signed(r.ext, r.chrome, now.Add(6*time.Minute).UnixMilli(), now.Add(6*time.Minute).UnixMilli()), http.StatusBadRequest},
		"an unknown Chrome":     {signed(r.ext, "chr_unknown0001", ms, ms), http.StatusNotFound},
		"no signature":          {`{"chrome_id":"` + r.chrome + `","at":` + fmt.Sprint(ms) + `,"sig":""}`, http.StatusBadRequest},
		"another field":         {`{"chrome_id":"` + r.chrome + `","at":` + fmt.Sprint(ms) + `,"sig":"x","tenant":"t"}`, http.StatusBadRequest},
		"a malformed chrome id": {signed(r.ext, "../x", ms, ms), http.StatusBadRequest},
	} {
		if w := post(c.body); w.Code != c.want {
			t.Errorf("%s: %d %s, want %d", name, w.Code, w.Body.String(), c.want)
		}
		if !paired() {
			t.Fatalf("%s unpaired the Chrome", name)
		}
	}

	// Its own key, now: unpaired, its key gone, a connect told 4410; the answer readable by the extension.
	good := signed(r.ext, r.chrome, ms, ms)
	w := post(good)
	if w.Code != http.StatusOK || w.Header().Get("Access-Control-Allow-Origin") != ext {
		t.Fatalf("its own unpair: %d %s %v", w.Code, w.Body.String(), w.Header())
	}
	if paired() {
		t.Fatal("the Chrome is still paired")
	}
	var ix chromeIndex
	if _, err := getDoc(ctx, r.e.s.KV, pkChromeIndex(r.chrome), "-", &ix); err != nil || !ix.Unpaired {
		t.Fatalf("a connect isn't told it was unpaired: %+v %v", ix, err)
	}
	// The same request again: nothing left to unpair.
	if w := post(good); w.Code != http.StatusNotFound {
		t.Fatalf("a replay: %d %s", w.Code, w.Body.String())
	}
}
