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
}
