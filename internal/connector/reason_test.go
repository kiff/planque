package connector

import (
	"strings"
	"testing"
)

// A refusal's reason reaches the owner and the taf, so it is one bounded line
// that never carries a secret, whatever the vendor or provider echoed
// (kiff/planque: a refused call said only a fixed line, so nobody could tell
// a test/live mismatch from a broken tool).
func TestReasonIsOneSafeBoundedLine(t *testing.T) {
	for _, c := range []struct{ raw, want string }{
		{`{"message":"No such invoice: 'in_1'"}`, "No such invoice: 'in_1'"},
		{`{"error":{"message":"Invalid API Key provided: sk_live_abc123"}}`, "Invalid API Key provided: [redacted]"},
		{`{"error":"forbidden"}`, "forbidden"},
		{`"No such customer"`, "No such customer"},
		{"plain\ntext\x00 with\tcontrols", "plain text with controls"},
		{`{"detail":"Authorization: Bearer abc.def-ghi rejected"}`, "Authorization: [redacted] rejected"},
		{`{"message":"the project key synthetic-key is wrong"}`, "the project key [redacted] is wrong"},
		// Key formats with a fixed prefix, and a key the vendor quoted after naming it: planque never holds the
		// person's own key for an app, so the shape and the label are what catch it.
		{`{"message":"token xoxb-1234567890-abcdefghij revoked"}`, "token [redacted] revoked"},
		{`{"message":"bad credentials for ghp_aaaaBBBBccccDDDDeeeeFFFF1234"}`, "bad credentials for [redacted]"},
		{`{"message":"AKIAABCDEFGHIJKLMNOP is not valid"}`, "[redacted] is not valid"},
		{`{"message":"Invalid API key provided: 'a1b2c3d4e5f6a7b8c9d0e1f2'"}`, "Invalid API key provided: '[redacted]'"},
		{`{"error":"api_key=0123456789abcdef0123 rejected"}`, "api_key=[redacted] rejected"},
		{`{"error":"Invalid token: 9f8e7d6c5b4a39281706f5e4"}`, "Invalid token: [redacted]"},
		// What isn't a key stays: short words after a label, and ids a person needs to read.
		{`{"message":"Invalid API key provided"}`, "Invalid API key provided"},
		{`{"message":"The token is expired"}`, "The token is expired"},
		{`{"message":"No such customer: 'cus_ABC123def456GHI'; a similar object exists in live mode, but a test mode key was used to make this request."}`, "No such customer: 'cus_ABC123def456GHI'; a similar object exists in live mode, but a test mode key was used to make this request."},
		{`{"data":[]}`, ""},
		{``, ""},
	} {
		if got := Reason([]byte(c.raw), "synthetic-key"); got != c.want {
			t.Errorf("Reason(%q) = %q, want %q", c.raw, got, c.want)
		}
	}
	long := Reason([]byte(`"` + strings.Repeat("x", 1000) + `"`))
	if r := []rune(long); len(r) != maxReason || !strings.HasSuffix(long, "…") {
		t.Fatalf("a long reason is %d characters: %q…", len(r), long[:20])
	}
}
