package planque

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/permit"
)

// For an account that enforces pins, each way a pin can fail to hold refuses
// the call, by path only, never a value: pinned and the permit signed none,
// pinned and not signed, signed with another value, the body not carrying
// it, and signed but not pinned here. Pins that hold, or none on either
// side, refuse nothing.
func TestAPinThatDoesNotHoldIsRefused(t *testing.T) {
	want := []kiff.Pin{{Path: "status", Value: "draft"}}
	body := json.RawMessage(`{"status":"draft"}`)
	signed := func(pins map[string]any) permit.Permit {
		return permit.Permit{Op: "op_1", PinsSigned: true, Pins: pins}
	}
	for _, c := range []struct {
		name string
		want []kiff.Pin
		p    permit.Permit
		body json.RawMessage
		why  string
	}{
		{"pinned, none signed", want, permit.Permit{Op: "op_1"}, body, "weren't signed by KIFF"},
		{"pinned, not signed", want, signed(map[string]any{}), body, "status: pinned here, not signed"},
		{"signed otherwise", want, signed(map[string]any{"status": "final"}), body, "status: signed with another value"},
		{"sent otherwise", want, signed(map[string]any{"status": "draft"}), json.RawMessage(`{"status":"final"}`), "status: the body sent doesn't carry the pinned value"},
		{"signed, not pinned", nil, signed(map[string]any{"stripe_context": "acct_9"}), body, "stripe_context: signed, not pinned here"},
		{"they hold", want, signed(map[string]any{"status": "draft"}), body, ""},
		{"none either side", nil, permit.Permit{Op: "op_1"}, body, ""},
	} {
		got := pinRefusal(c.want, c.p, c.body)
		if (c.why == "") != (got == "") || !strings.Contains(got, c.why) {
			t.Errorf("%s: %q, want %q", c.name, got, c.why)
		}
		for _, v := range []string{"draft", "final", "acct_9", "op_1"} {
			if strings.Contains(got, v) {
				t.Errorf("%s: a value in %q", c.name, got)
			}
		}
	}
	s := &Server{EnforcePins: map[string]bool{"t_on": true}}
	if !s.enforcesPins("t_on") || s.enforcesPins("t_off") || (&Server{}).enforcesPins("t_on") || !(&Server{EnforcePins: map[string]bool{"*": true}}).enforcesPins("any") {
		t.Fatal("pins enforced for the wrong accounts")
	}
}

// Through the relay: an account that enforces pins has a call whose permit
// signs a pin planque doesn't set refused with pin_mismatch, and nothing
// reaches the vendor; the same call for an account that doesn't enforce is
// only logged and goes through, and a call with no pins on either side goes
// through for both.
func TestTheRelayRefusesAPinDifferenceForAnAccountThatEnforces(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	args := func(op string) map[string]any {
		return map[string]any{"charge": "ch_1", "amount": 1, IdempotencyArg: op}
	}
	mint := func(op string, pins map[string]any) string {
		return e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Sub: "agt_taf", Op: op,
			Args: mustJSON(args(op)), Pins: pins})
	}
	stray := map[string]any{"metadata.kiff_op": "op_x"}

	sent := e.vendor.CallCount()
	if out := e.relayCall(g, "create_refund", args("op_log"), mint("op_log", stray)); strings.Contains(out, CodePinMismatch) || e.vendor.CallCount() != sent+1 {
		t.Fatalf("not enforced: %s", out)
	}
	e.s.EnforcePins = map[string]bool{tenant: true}
	sent = e.vendor.CallCount()
	out := e.relayCall(g, "create_refund", args("op_refused"), mint("op_refused", stray))
	if !strings.Contains(out, `"code":"`+CodePinMismatch+`"`) || !strings.Contains(out, `"sent":false`) || !strings.Contains(out, "metadata.kiff_op: signed, not pinned here") || e.vendor.CallCount() != sent {
		t.Fatalf("enforced: %s, vendor calls %d", out, e.vendor.CallCount()-sent)
	}
	if out := e.relayCall(g, "create_refund", args("op_plain"), mint("op_plain", nil)); strings.Contains(out, CodePinMismatch) || e.vendor.CallCount() != sent+1 {
		t.Fatalf("no pins either side: %s", out)
	}
}
