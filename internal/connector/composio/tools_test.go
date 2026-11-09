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
