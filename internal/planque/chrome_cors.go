package planque

import (
	"net/http"
	"regexp"
)

// The tafs Chrome extension calls planque's /v1/extension routes from its
// service worker with no host permission for planque (kiff/tafs
// docs/chrome.md: none is asked at install), so Chrome sends each as a CORS
// request, and a JSON body or a key in the query first as a preflight. These
// routes carry no cookie and no bearer: each step proves itself (the code's
// HMAC proof, the signed bodies, and the released builds' manifest hashes).
// So a browser is let read their answers from an extension's own origin
// only, never `*` and never with credentials. The origin's scheme and shape
// are the gate, as no Web Store id is fixed yet; the proofs stay the checks.

// extensionOrigin is a Chrome extension's origin: the scheme, and an id of 32
// letters a to p, as Chrome makes them.
var extensionOrigin = regexp.MustCompile(`^chrome-extension://[a-p]{32}$`)

// allowExtension says a browser may read the answer, when the request comes
// from an extension's own origin: its origin echoed, never `*`, and no
// credentials. It reports whether it did.
func allowExtension(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Add("Vary", "Origin")
	origin := r.Header.Get("Origin")
	if !extensionOrigin.MatchString(origin) {
		return false
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	return true
}

// extensionCORS lets the extension read h's answers, refusals included.
func extensionCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowExtension(w, r)
		h(w, r)
	}
}

// extensionPreflight answers Chrome's preflight for /v1/extension: GET and
// POST with a JSON body, from an extension's origin. Anything else gets no
// CORS headers, so the browser stops it.
func extensionPreflight(w http.ResponseWriter, r *http.Request) {
	if allowExtension(w, r) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Access-Control-Max-Age", "600")
	}
	w.WriteHeader(http.StatusNoContent)
}
