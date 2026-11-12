package planque

import (
	"encoding/json"
	"testing"
)

// A narrowed read called with no arguments (a schedule's check) still sends
// the object its lifted arguments go in: Stripe's stripe_api_read with {} got
// no parameters object and answered with an error, so the check never read.
// Set arguments land in it as before, and the fixed ones stay.
func TestANarrowedReadWithNoArgumentsSendsItsParametersObject(t *testing.T) {
	tool := ToolRequest{Tool: "stripe_api_read", Name: "list_disputes", ReadOnly: true,
		Fixed:     map[string]json.RawMessage{"stripe_api_operation_id": json.RawMessage(`"GetDisputes"`), "livemode": json.RawMessage(`false`)},
		Arguments: map[string]Argument{"limit": {Path: "parameters.limit", Type: "integer"}, "starting_after": {Path: "parameters.starting_after", Type: "string"}}}
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{}, `{"livemode":false,"parameters":{},"stripe_api_operation_id":"GetDisputes"}`},
		{map[string]any{"limit": json.Number("100")}, `{"livemode":false,"parameters":{"limit":100},"stripe_api_operation_id":"GetDisputes"}`},
	} {
		got, err := tool.vendorArgs(c.args)
		if err != nil || string(got) != c.want {
			t.Fatalf("%v: %s, %v; want %s", c.args, got, err, c.want)
		}
	}
}

// The object is only made when absent: a fixed value in it, or a fixed object
// at its key, is still sent as approved when the agent sets nothing.
func TestTheParametersObjectNeverDropsAFixedValue(t *testing.T) {
	lifted := map[string]Argument{"limit": {Path: "parameters.limit", Type: "integer"}}
	for _, c := range []struct {
		fixed map[string]json.RawMessage
		want  string
	}{
		{map[string]json.RawMessage{"parameters.expand": json.RawMessage(`["data.evidence"]`)}, `{"parameters":{"expand":["data.evidence"]}}`},
		{map[string]json.RawMessage{"parameters": json.RawMessage(`{"status":"open"}`)}, `{"parameters":{"status":"open"}}`},
	} {
		tool := ToolRequest{Tool: "stripe_api_read", Title: "List", Name: "list", ReadOnly: true, Fixed: c.fixed, Arguments: lifted}
		got, err := tool.vendorArgs(map[string]any{})
		if err != nil || string(got) != c.want {
			t.Fatalf("fixed %s: %s, %v; want %s", c.fixed, got, err, c.want)
		}
	}
	// A value already at the path's key is left as it is.
	m := map[string]any{"parameters": "kept"}
	parents(m, []string{"parameters", "inner"})
	if m["parameters"] != "kept" {
		t.Fatalf("parents replaced a value: %v", m)
	}
}

// A write never gets an empty object its agent didn't fill: a refund's
// unset human_confirmation stays out of the call.
func TestAWriteGetsNoEmptyObject(t *testing.T) {
	tool := ToolRequest{Tool: "stripe_api_write", Name: "refund",
		Arguments: map[string]Argument{"amount": {Path: "parameters.amount", Type: "integer"}, "approval_token": {Path: "human_confirmation.approval_token", Type: "string"}}}
	got, err := tool.vendorArgs(map[string]any{"amount": json.Number("300")})
	if err != nil || string(got) != `{"parameters":{"amount":300}}` {
		t.Fatalf("%s, %v", got, err)
	}
}
