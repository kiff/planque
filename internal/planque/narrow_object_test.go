package planque

import (
	"encoding/json"
	"fmt"
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
		"open, no fields":    {`{"type":"object","additionalProperties":{"type":"string"}}`, "open object can't be checked: an object names each field the taf writes"},
		"open, fields empty": {`{"type":"object","properties":{},"additionalProperties":true}`, "open object can't be checked: an object names each field the taf writes"},
		"no fields":          {`{"type":"object","properties":{}}`, "an object names each field the taf writes"},
		"fields unnamed":     {`{"type":"object"}`, `in properties: {"type": "object", "properties": {"Status"`},
		"no schema":          {``, "it has no schema: an object names each field the taf writes"},
		"not an object":      {`{"type":"array","items":{"type":"string"}}`, "is an object"},
		"a nested object":    {`{"type":"object","properties":{"a":{"type":"object","properties":{"b":{"type":"string"}}}}}`, "never an object"},
		"a null field":       {`{"type":"object","properties":{"a":{"type":"null"}}}`, "type null isn't one"},
		"a type list":        {`{"type":"object","properties":{"a":{"type":["string","null"]}}}`, "JSON Schema of"},
		"alternatives":       {`{"type":"object","properties":{"a":{"oneOf":[{"type":"string"},{"type":"integer"}]}}}`, "JSON Schema of"},
		"a field untyped":    {`{"type":"object","properties":{"a":{}}}`, "type  isn't one"},
		"an array unchecked": {`{"type":"object","properties":{"a":{"type":"array"}}}`, "a: "},
		"required not named": {`{"type":"object","properties":{"a":{"type":"string"}},"required":["b"]}`, "required field b"},
		"a keyword":          {`{"type":"object","properties":{"a":{"type":"string","pattern":"x"}}}`, "JSON Schema of"},
		"a dot in a name":    {`{"type":"object","properties":{"A.B":{"type":"string"}}}`, "no dot"},
		"an empty name":      {`{"type":"object","properties":{"":{"type":"string"}}}`, "no dot"},
		"a name too long":    {`{"type":"object","properties":{"` + strings.Repeat("a", 65) + `":{"type":"string"}}}`, "1 to 64"},
		"an enum on it":      {`{"type":"object","properties":{"a":{"type":"string"}},"enum":["x"]}`, "description only"},
		"items on it":        {`{"type":"object","properties":{"a":{"type":"string"}},"items":{"type":"string"}}`, "description only"},
		"a length on it":     {`{"type":"object","properties":{"a":{"type":"string"}},"maxItems":2}`, "description only"},
		"too many fields":    {manyFields(51), "at most 50 fields"},
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
	// A sheet's column with a space in its name is a field; 50 fields are taken.
	for _, schema := range []string{`{"type":"object","properties":{"First name":{"type":"string"}}}`, manyFields(50)} {
		tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: "object", Schema: json.RawMessage(schema)}}}
		if msg := checkNarrowing(tool); msg != nil {
			t.Fatal(msg)
		}
	}
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

// manyFields is an object schema naming n string fields.
func manyFields(n int) string {
	fields := make([]string, n)
	for i := range fields {
		fields[i] = fmt.Sprintf(`"f%d":{"type":"string"}`, i)
	}
	return `{"type":"object","properties":{` + strings.Join(fields, ",") + `}}`
}

// An object that names its fields but is left open (additionalProperties true,
// or a schema such as the sheet's "a text per column", which a taf copies from
// the vendor: the Eval's intros) is stored closed, never refused. Closing only
// narrows it: what is stored, shown and signed says additionalProperties
// false, and a call with a field it doesn't name is still refused.
func TestAnOpenObjectNamingItsFieldsIsStoredClosed(t *testing.T) {
	for name, schema := range map[string]string{
		"open to a schema": `{"type":"object","properties":{"status":{"type":"string"}},"additionalProperties":{"type":"string"}}`,
		"open":             `{"type":"object","properties":{"status":{"type":"string"}},"additionalProperties":true}`,
		"left out":         `{"type":"object","properties":{"status":{"type":"string"}}}`,
		"closed":           `{"type":"object","properties":{"status":{"type":"string"}},"additionalProperties":false}`,
	} {
		tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: "object", Required: true, Schema: json.RawMessage(schema)}}}
		normaliseArguments(tool)
		if msg := checkNarrowing(tool); msg != nil {
			t.Fatalf("%s: %v", name, msg)
		}
		var stored map[string]any
		_ = json.Unmarshal(tool.Arguments["values"].Schema, &stored)
		if stored["additionalProperties"] != false {
			t.Errorf("%s: stored as %s", name, tool.Arguments["values"].Schema)
		}
		if _, err := tool.vendorArgs(map[string]any{"values": map[string]any{"status": "done", "Email": "x@example.com"}}); err == nil {
			t.Errorf("%s: a field it doesn't name was taken", name)
		}
		if _, err := tool.vendorArgs(map[string]any{"values": map[string]any{"status": "done"}}); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// A refused argument's declaration is logged by its keys and types alone: its
// fields' names and types, an enum's count, never a description, a value or an
// enum's values.
func TestARefusedArgumentsShapeHasNoValues(t *testing.T) {
	schema := `{"type":"object","description":"SECRET-DESC","properties":{"status":{"type":"string","enum":["SECRET-ENUM"],"pattern":"SECRET-PAT"}}}`
	tool := ToolRequest{Tool: "updateRow", Arguments: map[string]Argument{"values": {Path: "values", Type: "object", Schema: json.RawMessage(schema)}}}
	normaliseArguments(tool)
	msg := checkNarrowing(tool)
	if msg == nil || msg.Code != refArgumentSchema {
		t.Fatalf("%v", msg)
	}
	if !strings.Contains(msg.Shape, "type=object") || !strings.Contains(msg.Shape, "status:{enum:1,pattern,type:string}") {
		t.Fatalf("shape %q", msg.Shape)
	}
	if strings.Contains(msg.Shape, "SECRET") {
		t.Fatalf("a value is in the shape: %q", msg.Shape)
	}
	if r := checkNarrowing(ToolRequest{Tool: "x", Arguments: map[string]Argument{"n": {Path: "n", Type: "map"}}}); r == nil || r.Shape != "type=map" {
		t.Fatalf("a type refusal's shape: %+v", r)
	}
}
