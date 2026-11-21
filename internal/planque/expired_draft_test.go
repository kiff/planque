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
// sign-in nobody will approve a taf for: a week past its request the
// sign-in is revoked, and the grant has the deadline the sweep finds it by.
// Before the week it stays, a day included: an app may be connected one
// evening and its taf set up the next.
func TestADraftThatNeverAskedForToolsEndsAfterAWeek(t *testing.T) {
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
	if d := g.deadline(); !d.Equal(asked.Add(pendingTTL)) {
		t.Fatalf("the draft's deadline is %v, want the request's time plus a week (%v)", d, asked.Add(pendingTTL))
	}

	// A day and a minute on: the next evening's taf is still welcome.
	e.s.Now = func() time.Time { return asked.Add(holdTTL + time.Minute) }
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 0 {
		t.Fatalf("the sweep after a day: ended %d, %v", n, err)
	}
	e.s.Now = func() time.Time { return asked.Add(pendingTTL - time.Hour) }
	if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("a week less an hour on: %d %v", st, out["status"])
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the sign-in was revoked before the week passed")
	}

	e.s.Now = func() time.Time { return asked.Add(pendingTTL + time.Minute) }
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 1 {
		t.Fatalf("the sweep after a week: ended %d, %v", n, err)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "no tools were asked for") || len(g.Tools) != 0 || g.Transport != "" {
		t.Fatalf("the grant after a week: %s %q tools %v", g.Status, g.Reason, g.Tools)
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

// A taf from a catalog entry that lists its tools has its Setup reads marked
// on its own grant (markEntrySetupReads), and the request's own setup_reads
// are cleared. A read-only taf whose every tool is marked is a live taf, not
// a draft: neither rule ends its sign-in, however old (kiff/planque#322
// review). Modelled as that grant is stored: the connect-first draft with
// Request.SetupReads cleared.
func TestAReadOnlyTafWhoseEveryToolIsMarkedSetupReadsIsNotADraft(t *testing.T) {
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
	ctx := context.Background()
	g, _ := e.s.loadGrant(ctx, tenant, id)
	g.Request.SetupReads = nil // as an entry's own marking leaves the grant
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatalf("save: %v", err)
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	if len(g.Tools) == 0 || !g.Tools[0].SetupReads {
		t.Fatalf("the tools aren't marked: %+v", g.Tools)
	}
	asked, _ := time.Parse(time.RFC3339, g.Request.CreatedAt)
	if !g.deadline().IsZero() {
		t.Fatalf("a live taf has a draft's deadline: %v", g.deadline())
	}
	acct := "/v1/accounts/" + tenant + "/grants/" + id

	// Rule 2: thirty days on, still connected, sign-in live, the sweep ends nothing.
	e.s.Now = func() time.Time { return asked.Add(30 * 24 * time.Hour) }
	if n, err := e.s.Sweep(ctx); err != nil || n != 0 {
		t.Fatalf("the sweep a month on: ended %d, %v", n, err)
	}
	if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("a month on: %d %v", st, out["status"])
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("a live taf's sign-in was revoked")
	}

	// Rule 1: an addition nobody approves goes after its day, and drops only the addition.
	e.s.Now = func() time.Time { return time.Now() }
	refund := map[string]any{"tool": "create_refund", "title": "Refund a charge", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup"); st != 200 || out["declared"] != true {
		t.Fatalf("the addition: %d %v", st, out)
	}
	e.s.Now = func() time.Time { return time.Now().Add(holdTTL + time.Minute) }
	if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("a day after the addition: %d %v", st, out["status"])
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	if len(g.Tools) != 1 || len(g.Pending) != 0 || g.Tools[0].Tool != "list_charges" {
		t.Fatalf("what the grant holds after the addition went: tools %+v pending %v", g.Tools, g.Pending)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the addition's expiry revoked the taf's sign-in")
	}
	if _, ok := e.kiff.Connection(tenant, "create_refund"); ok {
		t.Fatal("the unapproved addition is still in KIFF")
	}
}
