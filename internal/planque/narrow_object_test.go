package planque

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// rowValues is a sheet row's values as an object argument: the one column the
// taf writes, named.
var rowValues = json.RawMessage(`{"type":"object","properties":{"Introduction":{"type":"string","maxLength":2000}},"required":["Introduction"]}`)

func updateRow() ToolRequest {
	return ToolRequest{Tool: "updateRow", Name: "update_row", Arguments: map[string]Argument{
		"row":    {Path: "row", Type: "integer", Required: true},
		"values": {Path: "values", Type: "object", Required: true, Schema: rowValues},
	}}
}

// A vendor argument that is an object (a sheet row's values, column → text)
// can be declared as one, naming each field the taf may set (the Eval's
// intros: Coach declared values as an object and planque refused it). A call
// goes to the vendor as the taf sent it, and the taf sees the fields named.
func TestAnObjectArgumentNamesItsFields(t *testing.T) {
	tool := updateRow()
	if msg := checkNarrowing(tool); msg != nil {
		t.Fatal(msg)
	}
	got, err := tool.vendorArgs(map[string]any{"row": json.Number("3"), "values": map[string]any{"Introduction": "Hello Dana"}})
	if want := `{"row":3,"values":{"Introduction":"Hello Dana"}}`; err != nil || string(got) != want {
		t.Fatalf("%s, %v; want %s", got, err, want)
	}
	var schema struct {
		Properties map[string]map[string]any `json:"properties"`
	}
	_ = json.Unmarshal(tool.schema(), &schema)
	v := schema.Properties["values"]
	if v["type"] != "object" || v["properties"] == nil || v["additionalProperties"] != false {
		t.Fatalf("the taf's schema for values: %v", v)
	}
	// The owner reads which fields the taf sets.
	if text := narrowingText(tool); !strings.Contains(text, "values (Introduction)") {
		t.Fatalf("the approval says %q", text)
	}
}

// A call that sets a field the object doesn't name (another column), a field
// of the wrong kind, leaves out a required one, or sends something that isn't
// an object is refused before the vendor gets it.
func TestAnObjectArgumentTakesOnlyItsNamedFields(t *testing.T) {
	tool := updateRow()
	for name, values := range map[string]any{
		"another column":    map[string]any{"Introduction": "Hi", "Status": "done"},
		"only another":      map[string]any{"Status": "done"},
		"the wrong kind":    map[string]any{"Introduction": 12},
		"required left out": map[string]any{},
		"not an object":     "Hi",
		"a list":            []any{"Hi"},
		"too long":          map[string]any{"Introduction": strings.Repeat("a", 2001)},
	} {
		if _, err := tool.vendorArgs(map[string]any{"row": json.Number("3"), "values": values}); err == nil {
			t.Errorf("%s: taken", name)
		}
	}
}

// An object argument's schema names its fields, each a plain value or an
// array as an array argument's schema is: never open, never another object,
// never null, a type list or alternatives (kiff-qa's negatives).
func TestAnObjectArgumentsSchemaIsChecked(t *testing.T) {
	for name, c := range map[string]struct {
		schema string
		want   string
	}{
		"open fields":        {`{"type":"object","properties":{"a":{"type":"string"}},"additionalProperties":true}`, "only the fields it names"},
		"open to a schema":   {`{"type":"object","properties":{"a":{"type":"string"}},"additionalProperties":{"type":"string"}}`, "JSON Schema of"},
		"no fields":          {`{"type":"object","properties":{}}`, "names its fields"},
		"fields unnamed":     {`{"type":"object"}`, "names its fields"},
		"no schema":          {``, "JSON Schema of"},
		"not an object":      {`{"type":"array","items":{"type":"string"}}`, "is an object"},
		"a nested object":    {`{"type":"object","properties":{"a":{"type":"object","properties":{"b":{"type":"string"}}}}}`, "never an object"},
		"a null field":       {`{"type":"object","properties":{"a":{"type":"null"}}}`, "type null isn't one"},
		"a type list":        {`{"type":"object","properties":{"a":{"type":["string","null"]}}}`, "JSON Schema of"},
		"alternatives":       {`{"type":"object","properties":{"a":{"oneOf":[{"type":"string"},{"type":"integer"}]}}}`, "JSON Schema of"},
		"a field untyped":    {`{"type":"object","properties":{"a":{}}}`, "type  isn't one"},
		"an array unchecked": {`{"type":"object","properties":{"a":{"type":"array"}}}`, "a: "},
		"required not named": {`{"type":"object","properties":{"a":{"type":"string"}},"required":["b"]}`, "required field b"},
		"a keyword":          {`{"type":"object","properties":{"a":{"type":"string","pattern":"x"}}}`, "JSON Schema of"},
	} {
		tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: "object", Schema: json.RawMessage(c.schema)}}}
		normaliseArguments(tool)
		msg := checkNarrowing(tool)
		if msg == nil || msg.Code != refArgumentSchema || !strings.Contains(msg.Msg, c.want) {
			t.Errorf("%s: %v, want %q", name, msg, c.want)
		}
	}
	// A field may be an array with its own reviewed schema, and an optional
	// field written as anyOf it and null is the field alone.
	ok := json.RawMessage(`{"type":"object","properties":{"Tags":{"type":"array","items":{"type":"string"}},"Note":{"anyOf":[{"type":"string"},{"type":"null"}]},"Count":{"type":"integer"},"Done":{"type":"boolean"}}}`)
	tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: "object", Schema: ok}}}
	normaliseArguments(tool)
	if msg := checkNarrowing(tool); msg != nil {
		t.Fatal(msg)
	}
	if _, err := tool.vendorArgs(map[string]any{"values": map[string]any{"Tags": []any{"a"}, "Note": "n"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := tool.vendorArgs(map[string]any{"values": map[string]any{"Tags": []any{1}}}); err == nil {
		t.Fatal("an array field took an item of the wrong kind")
	}
}

// An object argument can't share its path with a pin or another argument
// inside it, whichever comes first: the taf would set a field the pin fixes.
func TestAnObjectArgumentDoesntOverlap(t *testing.T) {
	inside := Argument{Path: "values.Status", Type: "string"}
	for range 20 { // map order varies: each order is met
		tool := updateRow()
		tool.Arguments["status"] = inside
		if msg := checkNarrowing(tool); msg == nil || msg.Code != refArgumentOverlap {
			t.Fatalf("an argument inside the object: %v", msg)
		}
	}
	tool := updateRow()
	tool.Fixed = map[string]json.RawMessage{"values.Status": json.RawMessage(`"done"`)}
	if msg := checkNarrowing(tool); msg == nil || msg.Code != refArgumentOverlap {
		t.Fatalf("a pin inside the object: %v", msg)
	}
}

// KIFF counts and scopes a Card on a plain value: an object argument is never
// the amount or the resource.
func TestAnObjectArgumentIsNeverTheAmountOrResource(t *testing.T) {
	tool := updateRow()
	tool.AmountArg, tool.AmountUnit = "values", "EUR"
	if msg := checkNarrowing(tool); msg == nil || msg.Code != refAmountArgument {
		t.Fatalf("an object as the amount: %v", msg)
	}
	tool = updateRow()
	tool.ResourceArg = "values"
	if msg := checkNarrowing(tool); msg == nil || msg.Code != refResourceArgument {
		t.Fatalf("an object as the resource: %v", msg)
	}
}

// A type planque doesn't take is refused with both ways to give an object,
// so a retry works whichever the taf picks.
func TestAnArgumentTypeRefusalShowsBothForms(t *testing.T) {
	for _, typ := range []string{"map", "", "null"} {
		tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: typ}}}
		msg := checkNarrowing(tool)
		if msg == nil || msg.Code != refArgumentType {
			t.Fatalf("type %q: %v", typ, msg)
		}
		for _, want := range []string{`"type": "object"`, `"properties"`, `values.Status`} {
			if !strings.Contains(msg.Want, want) || !strings.Contains(msg.Msg, want) {
				t.Errorf("type %q: the refusal doesn't show %s: %q", typ, want, msg.Want)
			}
		}
	}
}

// The intros sheet's row update, as its API describes it, takes values as an
// object: the request names it so, and the API's check takes it.
func TestAnObjectArgumentFitsTheAPIsObject(t *testing.T) {
	spec := []byte(`{"openapi": "3.0.3", "info": {"title": "Sheets", "version": "1"}, "servers": [{"url": "https://sheets.example/v1"}],
 "security": [{"k": []}], "components": {"securitySchemes": {"k": {"type": "apiKey", "in": "header", "name": "X-Sheets-Key"}},
  "schemas": {"Values": {"type": "object", "required": ["values"],
   "properties": {"values": {"type": "object", "additionalProperties": {"type": "string"}}}}}},
 "paths": {"/sheets/{sheet}/rows/{row}": {"patch": {"operationId": "updateRow",
   "parameters": [{"name": "sheet", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "row", "in": "path", "required": true, "schema": {"type": "integer"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Values"}}}}}}}}`)
	snap, err := openapi.Parse(spec, openapi.Auth{In: "header", Name: "X-Sheets-Key"})
	if err != nil {
		t.Fatal(err)
	}
	tool := updateRow()
	tool.Fixed = map[string]json.RawMessage{"sheet": json.RawMessage(`"sh_eval"`)}
	if msg := apiToolsCheck(snap, []ToolRequest{tool}); msg != nil {
		t.Fatal(msg)
	}
	// Declared as a string where the API has an object, it's still refused.
	tool.Arguments["values"] = Argument{Path: "values", Type: "string", Required: true}
	if msg := apiToolsCheck(snap, []ToolRequest{tool}); msg == nil || msg.Code != refArgumentType {
		t.Fatalf("a string for the API's object: %v", msg)
	}
}
