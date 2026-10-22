package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/permit"
)

// The owner's halt reaches planque's dispatch (kiff-cloud#1061): while the
// agent or its account is halted, nothing is sent, whatever permit the call
// carries; a permit issued before the halt stays refused once resumed; and
// when KIFF can't say, nothing is sent.
func TestAHaltStopsTheRelayAtDispatch(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	args := func(op string) map[string]any {
		return map[string]any{"charge": "ch_9", "amount": 100, IdempotencyArg: op}
	}
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
// least the time KIFF asked about, and only when KIFF asks: signed with its
// permit key, for that halt, still live.
func TestConfirmingAHaltReadsItBackFromKIFF(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	at := time.Now().UTC().Truncate(time.Microsecond)
	ask := func(tnt, agent string, when time.Time, token string) (int, map[string]any) {
		t.Helper()
		b, _ := json.Marshal(map[string]any{"agent": agent, "halted_at": when.Format(time.RFC3339Nano)})
		req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tnt+"/halts/confirm", bytes.NewReader(b))
		if token != "" {
			req.Header.Set(permit.HaltConfirmHeader, token)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var out map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&out)
		return resp.StatusCode, out
	}
	confirm := func(agent string, when time.Time) map[string]any {
		t.Helper()
		st, out := ask(tenant, agent, when, e.kiff.AskHaltConfirm(tenant, agent, when, 0))
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
	if st, out := ask("tnt_none", "a", at, e.kiff.AskHaltConfirm("tnt_none", "a", at, 0)); st != 200 || out["seen"] != true || out["enrolled"] != false {
		t.Fatalf("an account with nothing here: %d %v", st, out)
	}
	// Only KIFF's live ask for that halt is answered; nothing says whether
	// the account is enrolled here.
	stranger := fakes.NewKIFF()
	stranger.Issuer = e.kiff.Issuer
	for name, c := range map[string]struct {
		tnt, agent string
		token      string
	}{
		"unsigned":              {tenant, "agt_taf", ""},
		"signed by another key": {tenant, "agt_taf", stranger.AskHaltConfirm(tenant, "agt_taf", at, 0)},
		"for another agent":     {tenant, "agt_taf", e.kiff.AskHaltConfirm(tenant, "agt_other", at, 0)},
		"for another account":   {tenant, "agt_taf", e.kiff.AskHaltConfirm("tnt_none", "agt_taf", at, 0)},
		"for another time":      {tenant, "agt_taf", e.kiff.AskHaltConfirm(tenant, "agt_taf", at.Add(time.Second), 0)},
		"lives too long":        {tenant, "agt_taf", e.kiff.AskHaltConfirm(tenant, "agt_taf", at, 10*time.Minute)},
		"a permit":              {tenant, "agt_taf", e.kiff.Mint(fakes.Permit{Audience: "x", Tenant: tenant, Tool: "t", Sub: "agt_taf", Op: "op", Args: []byte("{}")})},
		"unknown account":       {"tnt_none", "a", ""},
	} {
		if st, out := ask(c.tnt, c.agent, at, c.token); st != http.StatusUnauthorized || out["seen"] != nil || out["enrolled"] != nil {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
}

// A trial's permit (sub agent~trial) is refused under its agent's halt, and
// after the resume when issued before it (kiff-cloud#1061, apps/api maps a
// trial to its agent).
func TestATrialPermitIsHeldByItsAgentsHalt(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	args := func(op string) map[string]any {
		return map[string]any{"charge": "ch_9", "amount": 100, IdempotencyArg: op}
	}
	mint := func(op string, at time.Time) string {
		b, _ := json.Marshal(args(op))
		return e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Sub: "agt_taf~trial", Op: op, Args: b, IssuedAt: at})
	}
	sent := e.vendor.CallCount()
	early := mint("op_trial_early", time.Now().Add(-10*time.Second))
	e.kiff.Halt(tenant, "agt_taf", time.Now().Add(-5*time.Second))
	if out := e.relayCall(g, "create_refund", args("op_trial_now"), mint("op_trial_now", time.Now())); !strings.Contains(out, "the owner halted this agent") {
		t.Fatalf("a trial's call while its agent is halted: %s", out)
	}
	e.kiff.Resume(tenant, "agt_taf")
	if out := e.relayCall(g, "create_refund", args("op_trial_early"), early); !strings.Contains(out, "issued before the owner's halt") {
		t.Fatalf("a trial permit from before the halt, after the resume: %s", out)
	}
	if e.vendor.CallCount() != sent {
		t.Fatalf("the vendor was called %d times for the halted agent's trial", e.vendor.CallCount()-sent)
	}
}
