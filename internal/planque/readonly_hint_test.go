package planque

import (
	"encoding/json"
	"testing"
)

// A vendor's own server marks a tool read-only: it's granted as a read only
// when nothing about it says it can change things.
func TestAVendorsReadOnlyMarkIsKeptOnlyForAToolThatOnlyReads(t *testing.T) {
	ro := json.RawMessage(`{"readOnlyHint":true}`)
	for _, c := range []struct {
		tool mcpTool
		read bool
	}{
		{mcpTool{Name: "list_refunds", Description: "Lists refunds.", InputSchema: json.RawMessage(`{"type":"object","properties":{"limit":{"type":"integer"}}}`), Annotations: ro}, true},
		{mcpTool{Name: "graphql", Description: "Runs a GraphQL document against the store.", InputSchema: json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}}}`), Annotations: ro}, false},
		{mcpTool{Name: "orders", Description: "Lists orders, or executes a custom GraphQL document that may perform mutations.", InputSchema: json.RawMessage(`{"type":"object"}`), Annotations: ro}, false},
		{mcpTool{Name: "list_refunds", Description: "Lists refunds.", InputSchema: json.RawMessage(`{"type":"object"}`)}, false},
	} {
		if got := readOnlyHint(c.tool); got != c.read {
			t.Errorf("%s (%q): read-only %v, want %v", c.tool.Name, c.tool.Description, got, c.read)
		}
	}
}
