package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// The owner's halt reaches planque's dispatch (kiff-cloud#1061): while the
// agent or its account is halted, nothing is sent, whatever permit the call
// carries; a permit issued before the halt stays refused once resumed; and
// when KIFF can't say, nothing is sent.
func TestAHaltStopsTheRelayAtDispatch(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	args := func(op string) map[string]any { return map[string]any{"charge": "ch_9", "amount": 100, IdempotencyArg: op} }
	mint := func(op string, at time.Time) string {
		b, _ := json.Marshal(args(op))
		return e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Sub: "agt_taf", Op: op, Args: b, IssuedAt: at})
	}
	sent := e.vendor.CallCount()

	early := mint("op_early", time.Now().Add(-10*time.Second))
	e.kiff.Halt(tenant, "agt_taf", time.Now().Add(-5*time.Second))
	if out := e.relayCall(g, "create_refund", args("op_now"), mint("op_now", time.Now())); !strings.Contains(out, "Refused by planque: the owner halted this agent") {
		t.Fatalf("a call while halted: %s", out)
	}
	e.kiff.Resume(tenant, "agt_taf")
	if out := e.relayCall(g, "create_refund", args("op_early"), early); !strings.Contains(out, "issued before the owner's halt") {
		t.Fatalf("a permit from before the halt, after the resume: %s", out)
	}
	if e.vendor.CallCount() != sent {
		t.Fatalf("the vendor was called %d times during and after the halt", e.vendor.CallCount()-sent)
	}
	// A new permit after the resume goes through.
	if out := e.relayCall(g, "create_refund", args("op_after"), mint("op_after", time.Now().Add(time.Second))); strings.Contains(out, "Refused") || e.vendor.CallCount() != sent+1 {
		t.Fatalf("a call after the resume: %s", out)
	}
	// Halt all holds every agent.
	e.kiff.Halt(tenant, "", time.Now())
	if res, _ := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_9", "amount": 100}); !strings.Contains(string(res.Result), "the owner halted this agent, or the whole account") {
		t.Fatalf("a call while the account is halted: %+v", res)
	}
	e.kiff.Resume(tenant, "")
	// KIFF can't say: nothing is sent.
	e.kiff.HaltsDown(true)
	if out := e.relayCall(g, "create_refund", args("op_down"), mint("op_down", time.Now().Add(2*time.Second))); !strings.Contains(out, "couldn't check whether the owner halted") {
		t.Fatalf("a call with KIFF's halt unreadable: %s", out)
	}
	if e.vendor.CallCount() != sent+1 {
		t.Fatal("the vendor was called while the halt couldn't be read")
	}
}

// planque confirms a halt only once it reads it back from KIFF, from at
// least the time KIFF asked about.
func TestConfirmingAHaltReadsItBackFromKIFF(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	at := time.Now().UTC().Truncate(time.Microsecond)
	confirm := func(agent string, when time.Time) map[string]any {
		t.Helper()
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/halts/confirm", map[string]any{"agent": agent, "halted_at": when.Format(time.RFC3339Nano)}, nil, "")
		if st != 200 {
			t.Fatalf("confirm: %d %v", st, out)
		}
		return out
	}
	if out := confirm("agt_taf", at); out["seen"] != false {
		t.Fatalf("confirmed a halt KIFF doesn't have: %v", out)
	}
	e.kiff.Halt(tenant, "agt_taf", at)
	if out := confirm("agt_taf", at); out["seen"] != true {
		t.Fatalf("didn't confirm the halt: %v", out)
	}
	if out := confirm("agt_taf", at.Add(time.Minute)); out["seen"] != false {
		t.Fatalf("confirmed a later halt than KIFF's: %v", out)
	}
	if st, out := e.do("POST", "/v1/accounts/tnt_none/halts/confirm", map[string]any{"agent": "a", "halted_at": at.Format(time.RFC3339Nano)}, nil, ""); st != 200 || out["seen"] != true || out["enrolled"] != false {
		t.Fatalf("an account with nothing here: %d %v", st, out)
	}
}
