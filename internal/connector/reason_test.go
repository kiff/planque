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
