package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// draftTool narrows an email draft to its cc, with the schema an agent
// copied from the provider's catalog.
func draftTool(cc string) ToolRequest {
	return ToolRequest{Tool: "create_email_draft", Name: "draft_email", Title: "Draft an email",
		Fixed: map[string]json.RawMessage{"user_id": json.RawMessage(`"me"`)},
		Arguments: map[string]Argument{
			"to": {Path: "recipient_email", Type: "string", Required: true},
			"cc": {Path: "cc", Type: "array", Schema: json.RawMessage(cc)},
		}}
}

// A provider's generated schema for an optional list of addresses, with its
// title, default and examples, or written as anyOf it and null, is taken: its
// annotations are dropped, and the request carries, declares and signs the
// schema planque checks, without them.
func TestAProvidersAnnotationsAreDroppedFromAnArgumentsSchema(t *testing.T) {
	clean := `{"description":"Cc addresses.","items":{"type":"string"},"type":"array"}`
	for name, cc := range map[string]string{
		"annotated":        `{"type":"array","items":{"type":"string","format":"email","title":"Address"},"title":"Cc","default":[],"examples":[["a@b.com"]],"description":"Cc addresses."}`,
		"anyOf null":       `{"anyOf":[{"type":"array","items":{"type":"string"}},{"type":"null"}],"default":null,"title":"Cc","description":"Cc addresses."}`,
		"oneOf null first": `{"oneOf":[{"type":"null"},{"type":"array","items":{"type":"string"},"description":"inner"}],"nullable":true,"description":"Cc addresses."}`,
	} {
		tools := []ToolRequest{draftTool(cc)}
		if r := checkTools("agt_1", tools); r != nil {
			t.Errorf("%s: refused: %s", name, r.Msg)
			continue
		}
		if got := string(tools[0].Arguments["cc"].Schema); got != clean {
			t.Errorf("%s: stored %s, want %s", name, got, clean)
		}
		// What is declared to KIFF, and the digest the owner signs, are the
		// clean schema's.
		if s := string(tools[0].schema()); strings.Contains(s, "title") || strings.Contains(s, "default") || strings.Contains(s, "anyOf") {
			t.Errorf("%s: declared %s", name, s)
		}
		got, err1 := digest(Request{Tools: tools})
		want, err2 := digest(Request{Tools: []ToolRequest{draftTool(clean)}})
		if err1 != nil || err2 != nil || got != want {
			t.Errorf("%s: digest %s, want the clean schema's %s (%v %v)", name, got, want, err1, err2)
		}
		// A call is still checked against it: a list of addresses, nothing else.
		if _, err := tools[0].vendorArgs(map[string]any{"to": "x@y.z", "cc": []any{"a@b.com"}}); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		for _, bad := range []any{nil, "a@b.com", []any{json.Number("1")}} {
			if _, err := tools[0].vendorArgs(map[string]any{"to": "x@y.z", "cc": bad}); err == nil {
				t.Errorf("%s: cc %v was sent", name, bad)
			}
		}
	}
}

// NEGATIVE: annotations go at every level, but what restricts a value stays:
// an object's field list, its required fields, enum, maxLength and maxItems.
// A field named like an annotation ("title") is a field, never dropped.
func TestDroppingAnnotationsKeepsWhatRestricts(t *testing.T) {
	in := `{"type":"array","title":"Lines","maxItems":2,"items":{"type":"object","title":"Line","examples":[{}],"required":["title"],
		"properties":{"title":{"type":"string","maxLength":40,"default":"x"},
		"kind":{"anyOf":[{"type":"string","enum":["a","b"]},{"type":"null"}],"title":"Kind"}}}}`
	tool := draftTool(in)
	if r := checkTools("agt_1", []ToolRequest{tool}); r != nil {
		t.Fatal(r.Msg)
	}
	got := string(tool.Arguments["cc"].Schema)
	want := `{"items":{"properties":{"kind":{"enum":["a","b"],"type":"string"},"title":{"maxLength":40,"type":"string"}},"required":["title"],"type":"object"},"maxItems":2,"type":"array"}`
	if got != want {
		t.Fatalf("stored %s\nwant %s", got, want)
	}
	for name, c := range map[string]struct {
		v    any
		want string
	}{
		"too many":       {[]any{map[string]any{"title": "a"}, map[string]any{"title": "b"}, map[string]any{"title": "c"}}, "at most 2 items"},
		"too long":       {[]any{map[string]any{"title": strings.Repeat("a", 41)}}, "at most 40 characters"},
		"outside enum":   {[]any{map[string]any{"title": "a", "kind": "c"}}, "must be one of a, b"},
		"null kind":      {[]any{map[string]any{"title": "a", "kind": nil}}, "must be a string"},
		"unnamed field":  {[]any{map[string]any{"title": "a", "note": "x"}}, "not part of this grant"},
		"required title": {[]any{map[string]any{"kind": "a"}}, "is required"},
	} {
		if _, err := tool.vendorArgs(map[string]any{"to": "x@y.z", "cc": c.v}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v, want %q", name, err, c.want)
		}
	}
}

// NEGATIVE: a keyword that restricts a value planque doesn't check, or one it
// can't reason about, is still refused with argument_schema, never dropped:
// dropping it would let through what the owner was shown as refused.
func TestARestrictingKeywordIsStillRefused(t *testing.T) {
	for name, cc := range map[string]string{
		"pattern":          `{"type":"array","items":{"type":"string","pattern":"@example\\.com$"},"title":"Cc"}`,
		"minimum":          `{"type":"array","items":{"type":"integer","minimum":1}}`,
		"minItems":         `{"type":"array","minItems":1,"items":{"type":"string"}}`,
		"minLength":        `{"type":"array","items":{"type":"string","minLength":3}}`,
		"const":            `{"type":"array","items":{"type":"string","const":"a"}}`,
		"$ref":             `{"type":"array","items":{"$ref":"#/$defs/Address"}}`,
		"anyOf of two":     `{"anyOf":[{"type":"array","items":{"type":"string"}},{"type":"array","items":{"type":"integer"}}]}`,
		"anyOf of three":   `{"anyOf":[{"type":"array","items":{"type":"string"}},{"type":"null"},{"type":"string"}]}`,
		"anyOf of nulls":   `{"anyOf":[{"type":"null"},{"type":"null"}]}`,
		"allOf":            `{"allOf":[{"type":"array","items":{"type":"string"}},{"type":"null"}]}`,
		"anyOf overlaps":   `{"type":"array","items":{"type":"string"},"anyOf":[{"type":"array","items":{"type":"integer"}},{"type":"null"}]}`,
		"a null with more": `{"anyOf":[{"type":"array","items":{"type":"string"}},{"type":"null","maxLength":1}]}`,
	} {
		tool := draftTool(cc)
		r := checkNarrowing(func() ToolRequest { normaliseArguments(tool); return tool }())
		if r == nil || r.Code != refArgumentSchema || r.Field != "cc" {
			t.Errorf("%s: %+v, want argument_schema on cc", name, r)
		}
	}
}

// A schema with nothing to drop is kept byte for byte, so a grant asked
// before reads, compares and signs as it did.
func TestACleanSchemaIsKeptAsItCame(t *testing.T) {
	tool := draftInvoice()
	normaliseArguments(tool)
	if string(tool.Arguments["items"].Schema) != string(lines) {
		t.Fatalf("rewritten: %s", tool.Arguments["items"].Schema)
	}
}
