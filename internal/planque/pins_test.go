package planque

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/permit"
)

// What planque sets in a vendor's body after the permit is registered with KIFF as the tool's pins, read from
// the same steps the relay runs: Qonto's draft pins by the organisation's country, a Stripe tool's account, and
// none for any other tool.
func TestEachToolRegistersThePinsTheRelaySets(t *testing.T) {
	invoice := GrantedTool{ToolRequest: ToolRequest{Tool: "create_client_invoice", Name: "create_client_invoice"},
		InputSchema: json.RawMessage(`{"type":"object","properties":{"client_id":{"type":"string"},"status":{"type":"string"},"report_einvoicing":{"type":"boolean"}}}`)}
	qonto := func(country string) *Grant { return &Grant{Request: Request{Server: qontoServer}, Country: country} }
	cases := []struct {
		name string
		g    *Grant
		t    GrantedTool
		want []kiff.Pin
	}{
		{"Qonto, country unread", qonto(""), invoice, []kiff.Pin{{Path: "report_einvoicing", Value: false}, {Path: "status", Value: "draft"}}},
		{"Qonto, Italy", qonto("IT"), invoice, []kiff.Pin{{Path: "report_einvoicing", Value: false}, {Path: "status", Value: "draft"}}},
		{"Qonto, France", qonto("FR"), invoice, []kiff.Pin{{Path: "status", Value: "draft"}}},
		{"a Stripe read in one account", &Grant{Request: Request{Server: stripeServer}, Account: "acct_1"},
			GrantedTool{ToolRequest: ToolRequest{Tool: "stripe_api_read", Name: "list_disputes"}}, []kiff.Pin{{Path: "stripe_context", Value: "acct_1"}}},
		{"a Stripe read with no single account", &Grant{Request: Request{Server: stripeServer}},
			GrantedTool{ToolRequest: ToolRequest{Tool: "stripe_api_read", Name: "list_disputes"}}, nil},
		{"another tool", &Grant{Request: Request{Server: "https://mcp.vendor.example/mcp"}},
			GrantedTool{ToolRequest: ToolRequest{Tool: "do_it", Name: "do_it"}}, nil},
	}
	for _, c := range cases {
		if got := toolPins(c.g, c.t); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: pins %#v, want %#v", c.name, got, c.want)
		}
	}
}

// At call time the relay compares what it pins with what KIFF signed and with the body it sends, both ways, and
// names each difference by its path only, never a value.
func TestThePinsAreComparedBothWaysByPathOnly(t *testing.T) {
	want := []kiff.Pin{{Path: "status", Value: "draft"}, {Path: "metadata.kiff_op", Marker: true}}
	body := json.RawMessage(`{"status":"draft","metadata":{"kiff_op":"op_1"}}`)
	agree := permit.Permit{Op: "op_1", PinsSigned: true, Pins: map[string]any{"status": "draft", "metadata.kiff_op": "op_1"}}
	if d := pinDiffs(want, agree, body); len(d) != 0 {
		t.Fatalf("agreeing pins: %v", d)
	}
	cases := []struct {
		name string
		p    permit.Permit
		body json.RawMessage
		want string
	}{
		{"not signed", permit.Permit{Op: "op_1", PinsSigned: true, Pins: map[string]any{"metadata.kiff_op": "op_1"}}, body, "status: pinned here, not signed"},
		{"signed otherwise", permit.Permit{Op: "op_1", PinsSigned: true, Pins: map[string]any{"status": "final", "metadata.kiff_op": "op_1"}}, body, "status: signed with another value"},
		{"sent otherwise", agree, json.RawMessage(`{"status":"final","metadata":{"kiff_op":"op_1"}}`), "status: the body sent doesn't carry the pinned value"},
		{"a marker for another op", agree, json.RawMessage(`{"status":"draft","metadata":{"kiff_op":"op_2"}}`), "metadata.kiff_op: the body sent doesn't carry the pinned value"},
		{"signed, not pinned", permit.Permit{Op: "op_1", PinsSigned: true, Pins: map[string]any{"status": "draft", "metadata.kiff_op": "op_1", "stripe_context": "acct_9"}}, body, "stripe_context: signed, not pinned here"},
	}
	for _, c := range cases {
		d := pinDiffs(want, c.p, c.body)
		if len(d) != 1 || d[0] != c.want {
			t.Errorf("%s: %v, want [%s]", c.name, d, c.want)
		}
		for _, line := range d {
			if strings.Contains(line, "draft") || strings.Contains(line, "final") || strings.Contains(line, "acct_") || strings.Contains(line, "op_") {
				t.Errorf("%s: a value in the log line %q", c.name, line)
			}
		}
	}
}

// A tool with no pins registers as before: no vendor_pins field, which an older gateway would refuse.
func TestARegistrationWithNoPinsSendsNone(t *testing.T) {
	g := &Grant{Request: Request{Server: "https://mcp.vendor.example/mcp"}, Audience: "https://relay.example/g1"}
	b, _ := json.Marshal(kiffConnection(g, GrantedTool{ToolRequest: ToolRequest{Tool: "do_it", Name: "do_it"}}, "cred", ""))
	if strings.Contains(string(b), "vendor_pins") {
		t.Fatalf("sent vendor_pins with none: %s", b)
	}
}
