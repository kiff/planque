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
	tool := ToolRequest{Tool: "stripe_api_read", Name: "list_disputes",
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
