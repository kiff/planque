package composio

import (
	"encoding/json"
	"testing"
)

// The provider's read-only tag makes a tool read-only only when nothing
// about it says it can change things: a GraphQL tool tagged read-only that
// "may perform mutations" comes out a write.
func TestAReadOnlyTagIsKeptOnlyForAToolThatOnlyReads(t *testing.T) {
	for _, c := range []struct {
		raw  string
		read bool
	}{
		{`{"slug":"SHOPIFY_GRAPH_QL_ORDERS","name":"Orders","toolkit":{"slug":"shopify"},"version":"20300101_00",` +
			`"description":"Lists, retrieves, counts, and searches Shopify orders using GraphQL, or executes a custom GraphQL document that may perform mutations.",` +
			`"input_parameters":{"type":"object","properties":{"query":{"type":"string"}}},"tags":["readOnlyHint"]}`, false},
		{`{"slug":"DESK_LIST_TICKETS","name":"List tickets","toolkit":{"slug":"desk"},"version":"20300101_00",` +
			`"description":"Lists tickets.","input_parameters":{"type":"object","properties":{"limit":{"type":"integer"}}},"tags":["readOnlyHint"]}`, true},
		{`{"slug":"DESK_REPLY","name":"Reply","toolkit":{"slug":"desk"},"version":"20300101_00",` +
			`"description":"Replies to a ticket.","input_parameters":{"type":"object","properties":{}},"tags":[]}`, false},
	} {
		var w wireTool
		if err := json.Unmarshal([]byte(c.raw), &w); err != nil {
			t.Fatal(err)
		}
		tl, err := w.tool()
		if err != nil {
			t.Fatalf("%s: %v", w.Slug, err)
		}
		if tl.ReadOnly != c.read {
			t.Errorf("%s: read-only %v, want %v", w.Slug, tl.ReadOnly, c.read)
		}
	}
}

// A tool's output schema describes what a call relays: the provider's
// {data, error, successful} answer by its data alone, with its definitions
// kept so a reference still resolves; any other schema as given.
func TestAToolsOutputIsWhatACallRelays(t *testing.T) {
	tool := func(output string) map[string]any {
		t.Helper()
		var w wireTool
		raw := `{"slug":"GMAIL_FETCH_EMAILS","name":"Fetch emails","toolkit":{"slug":"gmail"},"version":"20300101_00",` +
			`"description":"Fetches emails.","input_parameters":{"type":"object","properties":{}},"tags":["readOnlyHint"],"output_parameters":` + output + `}`
		if err := json.Unmarshal([]byte(raw), &w); err != nil {
			t.Fatal(err)
		}
		tl, err := w.tool()
		if err != nil {
			t.Fatal(err)
		}
		var out map[string]any
		if err := json.Unmarshal(tl.Output, &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	// Wrapped: data's own schema, its messages at the top.
	got := tool(`{"type":"object","properties":{"data":{"type":"object","properties":{"messages":{"type":"array"}}},"error":{"type":"string"},"successful":{"type":"boolean"}}}`)
	if props, _ := got["properties"].(map[string]any); props["messages"] == nil || props["data"] != nil {
		t.Fatalf("wrapped: %v", got)
	}
	// Wrapped by reference: data's reference, with the definitions it points into.
	got = tool(`{"type":"object","properties":{"data":{"$ref":"#/$defs/D"},"error":{"type":"string"},"successful":{"type":"boolean"}},"$defs":{"D":{"type":"object","properties":{"messages":{"type":"array"}}}}}`)
	if got["$ref"] != "#/$defs/D" || got["$defs"] == nil {
		t.Fatalf("wrapped by reference: %v", got)
	}
	// Not the exact wrapper: a data field among others stays as given.
	for _, as := range []string{
		`{"type":"object","properties":{"data":{"type":"object"},"next":{"type":"string"}}}`,
		`{"type":"object","properties":{"data":{"type":"object"},"error":{"type":"string"},"successful":{"type":"boolean"},"log_id":{"type":"string"}}}`,
		`{"type":"object","properties":{"data":true,"error":{"type":"string"},"successful":{"type":"boolean"}}}`,
	} {
		got := tool(as)
		if props, _ := got["properties"].(map[string]any); props["data"] == nil {
			t.Fatalf("not the wrapper, changed: %s → %v", as, got)
		}
	}
}
