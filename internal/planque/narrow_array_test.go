package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// lines is an invoice's lines as a reviewed schema: each line's fields named,
// with no discount among them.
var lines = json.RawMessage(`{"type":"array","maxItems":2,"items":{"type":"object","properties":{
	"title":{"type":"string","maxLength":40},"quantity":{"type":"string"},"vat_rate":{"type":"string"},
	"unit_price":{"type":"object","properties":{"value":{"type":"string"},"currency":{"type":"string","enum":["EUR"]}},"required":["value","currency"]}},
	"required":["title","quantity","unit_price","vat_rate"]}}`)

func draftInvoice() ToolRequest {
	return ToolRequest{Tool: "create_client_invoice", Name: "qonto_draft_invoice", Title: "Draft an invoice",
		Fixed: map[string]json.RawMessage{"status": json.RawMessage(`"draft"`)},
		Arguments: map[string]Argument{
			"quote_id": {Path: "quote_id", Type: "string", Required: true},
			"items":    {Path: "items", Type: "array", Required: true, Schema: lines},
		}}
}

// An array argument goes to the vendor as the agent sent it, beside the fixed
// arguments, and the agent sees the reviewed schema for it.
func TestAnArrayArgumentIsSentWithItsFixedArguments(t *testing.T) {
	tool := draftInvoice()
	if msg := checkNarrowing(tool); msg != nil {
		t.Fatal(msg)
	}
	line := map[string]any{"title": "Pose", "quantity": "2", "vat_rate": "0.2", "unit_price": map[string]any{"value": "120.00", "currency": "EUR"}}
	got, err := tool.vendorArgs(map[string]any{"quote_id": "q_1", "items": []any{line}})
	want := `{"items":[{"quantity":"2","title":"Pose","unit_price":{"currency":"EUR","value":"120.00"},"vat_rate":"0.2"}],"quote_id":"q_1","status":"draft"}`
	if err != nil || string(got) != want {
		t.Fatalf("%s, %v; want %s", got, err, want)
	}
	// The agent can't set the fixed status: it isn't one of its arguments.
	if _, err := tool.vendorArgs(map[string]any{"quote_id": "q_1", "items": []any{line}, "status": "unpaid"}); err == nil {
		t.Fatal("the agent set the fixed status")
	}
	var schema struct {
		Properties map[string]map[string]any `json:"properties"`
	}
	_ = json.Unmarshal(tool.schema(), &schema)
	if schema.Properties["items"]["type"] != "array" || schema.Properties["items"]["items"] == nil {
		t.Fatalf("the agent's schema for items: %v", schema.Properties["items"])
	}
}

// A line with a field the schema doesn't name (a discount), a field of the
// wrong kind, a value outside its list, or a required field missing is refused
// before anything reaches the vendor.
func TestAnArrayArgumentTakesOnlyItsReviewedFields(t *testing.T) {
	tool := draftInvoice()
	ok := func() map[string]any {
		return map[string]any{"title": "Pose", "quantity": "2", "vat_rate": "0.2", "unit_price": map[string]any{"value": "120.00", "currency": "EUR"}}
	}
	for name, c := range map[string]struct {
		items any
		want  string
	}{
		"a discount": {[]any{func() map[string]any {
			l := ok()
			l["discount"] = map[string]any{"type": "percentage", "value": "10"}
			return l
		}()}, "items[0].discount is not part of this grant"},
		"not an array":      {ok(), "items must be an array"},
		"a number for text": {[]any{func() map[string]any { l := ok(); l["quantity"] = json.Number("2"); return l }()}, "items[0].quantity must be a string"},
		"another currency": {[]any{func() map[string]any {
			l := ok()
			l["unit_price"] = map[string]any{"value": "1", "currency": "USD"}
			return l
		}()}, "items[0].unit_price.currency must be one of EUR"},
		"no price":             {[]any{func() map[string]any { l := ok(); delete(l, "unit_price"); return l }()}, "items[0].unit_price is required"},
		"a title too long":     {[]any{func() map[string]any { l := ok(); l["title"] = strings.Repeat("a", 41); return l }()}, "items[0].title takes at most 40 characters"},
		"too many lines":       {[]any{ok(), ok(), ok()}, "items takes at most 2 items"},
		"a line not an object": {[]any{"Pose"}, "items[0] must be an object"},
	} {
		_, err := tool.vendorArgs(map[string]any{"quote_id": "q_1", "items": c.items})
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v, want %q", name, err, c.want)
		}
	}
}

// Only an array argument has a schema, and it says only what planque checks.
func TestAnArrayArgumentsSchemaIsChecked(t *testing.T) {
	for name, c := range map[string]struct {
		arg  Argument
		want string
	}{
		"no schema":            {Argument{Path: "items", Type: "array"}, "schema"},
		"not an array":         {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"object","properties":{"a":{"type":"string"}}}`)}, "is an array"},
		"items not objects":    {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"array","items":{"type":"string"}}`)}, "items are objects"},
		"open fields":          {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"array","items":{"type":"object","properties":{"a":{"type":"string"}},"additionalProperties":true}}`)}, "only the fields it names"},
		"an object unnamed":    {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"array","items":{"type":"object"}}`)}, "names its fields"},
		"a keyword unchecked":  {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"array","items":{"type":"object","properties":{"a":{"type":"string","pattern":"x"}}}}`)}, "JSON Schema of"},
		"required not a field": {Argument{Path: "items", Type: "array", Schema: json.RawMessage(`{"type":"array","items":{"type":"object","properties":{"a":{"type":"string"}},"required":["b"]}}`)}, "required field b"},
		"a scalar with schema": {Argument{Path: "n", Type: "string", Schema: lines}, "only an array argument"},
	} {
		tool := ToolRequest{Tool: "create", Fixed: map[string]json.RawMessage{"status": json.RawMessage(`"draft"`)}, Arguments: map[string]Argument{"items": c.arg}}
		if msg := checkNarrowing(tool); msg == nil || !strings.Contains(msg.Msg, c.want) {
			t.Errorf("%s: %q, want %q", name, msg, c.want)
		}
	}
}
