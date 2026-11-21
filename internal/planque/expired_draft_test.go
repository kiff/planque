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
	// A sweep finds nothing left to do.
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 0 {
		t.Fatalf("a later sweep: %d %v", n, err)
	}
}
