package planque

import (
	"strings"
	"testing"
)

// A key tool whose amount is a number (a decimal) is refused when it's asked
// for: KIFF counts a whole number of minor units and refuses it, so declared
// it would only fail at the owner's key. Nothing is asked and nothing is
// left in KIFF.
func TestAKeyToolsAmountKIFFWouldRefuseIsRefusedAtTheRequest(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	v.mu.Lock()
	v.desc = strings.Replace(restDesc, `"properties": {"amount": {"type": "integer"}}`, `"properties": {"amount": {"type": "number"}}`, 1)
	v.mu.Unlock()
	st, out := e.askAPIGrant(v, []map[string]any{{"tool": "createRefund", "name": "create_refund", "title": "Refund",
		"amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2, "resource_argument": "charge"}}, nil)
	if st != 400 || !strings.Contains(mustJSONString(out), "amount_argument amount") || !strings.Contains(mustJSONString(out), "is a number") {
		t.Fatalf("a number amount was taken: %d %v", st, out)
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("left in KIFF: %v", c)
	}
}

// A tool KIFF refuses when it's declared (a 4xx)
// fails the request with KIFF's words, instead of going on undeclared to
// fail at the owner's key. The tools declared before it are taken back.
func TestADeclareKIFFRefusesFailsTheRequest(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	e.kiff.RefuseConnect = map[string]string{"create_refund": "resource_argument must be a string argument of the tool"}
	tools := []map[string]any{apiTools[1], apiTools[0]} // list_charges first, so one is declared before the refusal
	st, out := e.askAPIGrant(v, tools, nil)
	if st != 400 || !strings.Contains(mustJSONString(out), "KIFF won't take create_refund") || !strings.Contains(mustJSONString(out), "resource_argument must be") {
		t.Fatalf("a request KIFF refused was taken: %d %v", st, out)
	}
	if got := e.declared("list_charges"); len(got) != 0 {
		t.Fatalf("left in KIFF: %v", got)
	}
	id, _ := out["grant"].(string)
	if id == "" || e.statusOf(id) != StatusFailed {
		t.Fatalf("the refused request isn't failed: %v", out)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if st, give := e.giveKey(id, g.Digest, "k_vendor"); st < 400 || len(e.kiff.Connections(tenant)) != 0 {
		t.Fatalf("a key was taken for a refused request: %d %v", st, give)
	}
}

// KIFF answering 5xx isn't a refusal: asking again may work, so the request
// goes on undeclared, as before, and its tools are declared at the key.
func TestADeclareKIFFFailsLeavesTheRequestUndeclared(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	e.kiff.ConnectDown = true
	st, out := e.askAPIGrant(v, apiTools, nil)
	if st != 201 || out["status"] != StatusPending || out["declared"] == true {
		t.Fatalf("KIFF down: %d %v, want pending and undeclared", st, out)
	}
	if !strings.Contains(mustJSONString(out), "weren't declared to KIFF before the key") {
		t.Fatalf("no word of it: %v", out)
	}
}
