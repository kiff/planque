package kiff

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A refusal's short reason travels with its error, so a log can tell an
// approval that expired from a profile that changed; a long or spaced reason
// (free text) is left out.
func TestARefusalKeepsItsShortReason(t *testing.T) {
	for _, c := range []struct{ body, want string }{
		{`{"error":"approval_refused","reason":"profile_not_approved"}`, "KIFF answered 403: approval_refused (profile_not_approved)"},
		{`{"error":"approval_refused"}`, "KIFF answered 403: approval_refused"},
		{`{"error":"approval_refused","reason":"the device said something long"}`, "KIFF answered 403: approval_refused"},
		{`{"error":"approval_refused","reason":"` + strings.Repeat("a", 65) + `"}`, "KIFF answered 403: approval_refused"},
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(c.body))
		}))
		cl := &Client{Gateway: srv.URL}
		_, err := cl.do(context.Background(), http.MethodPost, srv.URL, "k", map[string]string{}, nil)
		srv.Close()
		if err == nil || err.Error() != c.want {
			t.Errorf("%s: %v, want %q", c.body, err, c.want)
		}
	}
}
