package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// A refused pin names the tool and the pin, says which part is wrong, and
// shows the form planque wants, so Coach can fix the request by itself.
func TestARefusedPinNamesWhatToFix(t *testing.T) {
	for name, c := range map[string]struct {
		pin   string
		value string
		want  []string
	}{
		"a space in the path":     {"post author", `"urn:li:person:1"`, []string{`linkedin_post: fixed argument "post author" isn't a dotted path`, `{"author": "urn:li:person:…"}`}},
		"a path too deep":         {"a.b.c.d.e", `"x"`, []string{`fixed argument "a.b.c.d.e" isn't a dotted path`, "up to 4 names"}},
		"a value that isn't JSON": {"author", `urn:li:person:1`, []string{`linkedin_post: fixed argument "author"'s value isn't JSON`, `{"author": "urn:li:person:…"}`}},
		"an id beyond 2^53 - 1":   {"administration_id", `504060379197867402`, []string{`linkedin_post: fixed argument "administration_id" holds a whole number beyond 2^53 - 1`, `write it as text, in quotes`, `{"administration_id": "123456789012345678"}`}},
		"one inside an object":    {"filter", `{"admin": 504060379197867402}`, []string{`fixed argument "filter" holds a whole number beyond 2^53 - 1`}},
	} {
		tool := ToolRequest{Tool: "CREATE_POST", Name: "linkedin_post", Title: "Post", Fixed: map[string]json.RawMessage{c.pin: json.RawMessage(c.value)}}
		r := checkNarrowing(tool)
		if r == nil {
			t.Errorf("%s: accepted", name)
			continue
		}
		if r.Code != refPinShape || r.Tool != "linkedin_post" || r.Field != c.pin {
			t.Errorf("%s: code %q, tool %q, field %q", name, r.Code, r.Tool, r.Field)
		}
		for _, w := range c.want {
			if !strings.Contains(r.Msg, w) {
				t.Errorf("%s: %q lacks %q", name, r.Msg, w)
			}
		}
	}
	ok := ToolRequest{Tool: "CREATE_POST", Name: "linkedin_post", Fixed: map[string]json.RawMessage{"author": json.RawMessage(`"urn:li:person:1"`)}}
	if msg := checkNarrowing(ok); msg != nil {
		t.Errorf("a good pin was refused: %s", msg)
	}
	// The same id as text, and a whole number within 2^53 - 1, are both fine.
	for _, v := range []string{`"504060379197867402"`, `9007199254740991`} {
		pin := ToolRequest{Tool: "CREATE_POST", Name: "linkedin_post", Fixed: map[string]json.RawMessage{"administration_id": json.RawMessage(v)}}
		if msg := checkNarrowing(pin); msg != nil {
			t.Errorf("%s was refused: %s", v, msg.Msg)
		}
	}
}
