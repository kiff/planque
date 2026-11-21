package planque

import (
	"context"
	"strings"
	"testing"
	"time"
)

// A draft nobody approves keeps no sign-in. Its sign-in was approved for
// Setup's reads alone; the tools the taf then asked for were never approved
// and went after their day, but "what it had stays connected" left the
// sign-in, and the reads, connected for ever: a test Stripe sign-in stayed
// live four days after its draft was dropped.
func TestAnExpiredDraftThatOnlyHadSetupsReadsLeavesNoConnection(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	if st, out := e.enrolAs(phone); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	e.enrolled, e.phone = true, phone
	id, digest := e.askWithSetupReads("list_charges")
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	read := g.Tools[0].Name
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	refund := map[string]any{"tool": "create_refund", "title": "Refund a charge", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup"); st != 200 || out["declared"] != true {
		t.Fatalf("the taf's tools after the reads: %d %v", st, out)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("no live sign-in before the day passed")
	}

	e.s.Now = func() time.Time { return time.Now().Add(holdTTL + time.Minute) }
	st, out := e.do("GET", acct, nil, nil, "k_setup")
	if st != 200 || out["status"] == StatusConnected || out["status"] == StatusPending {
		t.Fatalf("after a day: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "not approved within a day") {
		t.Fatalf("the grant after a day: %s %q", g.Status, g.Reason)
	}
	if len(g.Tools) != 0 || len(g.Connected) != 0 || len(g.Pending) != 0 || g.Transport != "" || g.PendingTransport != "" {
		t.Fatalf("what the ended grant still holds: tools %v, connected %v, pending %v", g.Tools, g.Connected, g.Pending)
	}
	if a, _ := e.vendor.Live(); a != 0 {
		t.Fatalf("the sign-in is still live at the vendor: %d tokens", a)
	}
	for _, name := range []string{read, "create_refund"} {
		if _, ok := e.kiff.Connection(tenant, name); ok {
			t.Fatalf("%s is still in KIFF", name)
		}
	}
	// A call on the ended grant's read is not served.
	if res, err := e.kiff.Call(context.Background(), "k_taf", "list_charges", map[string]any{}); err == nil && res.Decision == "allowed" {
		if text, isErr := resultText(t, res.Result); !isErr {
			t.Fatalf("a call on the ended grant was served: %s", text)
		}
	}
	// A sweep finds nothing left to do.
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 0 {
		t.Fatalf("a later sweep: %d %v", n, err)
	}
}

// A draft that connected Setup's reads and never asked for tools holds a
// sign-in nobody will approve a taf for: a day past its request the sign-in
// is revoked, and the grant has the deadline the sweep finds it by. Before a
// day it stays, as Setup asks for the tools within minutes of the reads.
func TestADraftThatNeverAskedForToolsEndsAfterADay(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	if st, out := e.enrolAs(phone); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	e.enrolled, e.phone = true, phone
	id, digest := e.askWithSetupReads("list_charges")
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	read := g.Tools[0].Name
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	asked, err := time.Parse(time.RFC3339, g.Request.CreatedAt)
	if err != nil {
		t.Fatalf("the request's time: %v", err)
	}
	if d := g.deadline(); !d.Equal(asked.Add(holdTTL)) {
		t.Fatalf("the draft's deadline is %v, want the request's time plus a day (%v)", d, asked.Add(holdTTL))
	}

	e.s.Now = func() time.Time { return asked.Add(holdTTL - time.Hour) }
	if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("23 hours on: %d %v", st, out["status"])
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the sign-in was revoked before the day passed")
	}

	e.s.Now = func() time.Time { return asked.Add(holdTTL + time.Minute) }
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 1 {
		t.Fatalf("the sweep after a day: ended %d, %v", n, err)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "no tools were asked for") || len(g.Tools) != 0 || g.Transport != "" {
		t.Fatalf("the grant after a day: %s %q tools %v", g.Status, g.Reason, g.Tools)
	}
	if a, _ := e.vendor.Live(); a != 0 {
		t.Fatalf("the sign-in is still live at the vendor: %d tokens", a)
	}
	if _, ok := e.kiff.Connection(tenant, read); ok {
		t.Fatal("the read is still in KIFF")
	}
}

// A connection with a tool the owner approved is never taken for a draft,
// whatever its age.
func TestAConnectionWithAnApprovedToolIsNotADraft(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.abandonedDraft(time.Now().Add(30 * 24 * time.Hour)) {
		t.Fatal("a connection with an approved tool was taken for an abandoned draft")
	}
	e.s.Now = func() time.Time { return time.Now().Add(30 * 24 * time.Hour) }
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("a month on: %d %v", st, out["status"])
	}
}
