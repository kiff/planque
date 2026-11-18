package planque

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/store"
)

// An account that connects and disconnects a lot ends up holding hundreds of
// finished grants and their operation records (the Eval's held 960 grants and
// 8,004 operations), and every read of its grants pays for them. With the
// signature of someone who may delete the account, the finished ones go, with
// their own records, and nothing else: no live grant, no waiting or ended-by-
// the-vendor one, not the enrolment, and no finished grant that still holds a
// token or a provider access not confirmed disabled (that one is skipped and
// counted, since its record must never go while a secret of it remains).
func TestDeletingFinishedGrantsKeepsEverythingElse(t *testing.T) {
	e := newEnv(t)
	live := e.connect()
	ctx := context.Background()
	acct := "/v1/accounts/" + tenant + "/finished-grants"
	recent, old := time.Now().Add(-2*time.Hour), time.Now().Add(-45*24*time.Hour)

	save := func(g *Grant) {
		g.Tenant = tenant
		g.Request.Tenant, g.Request.CatalogEntry = tenant, "gmail-key@1"
		if g.Request.CreatedAt == "" {
			g.Request.CreatedAt = recent.UTC().Format(time.RFC3339)
		}
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	row := func(sk string) {
		if err := putDoc(ctx, e.kv, store.Item{PK: tenantPK(tenant), SK: sk, Ver: 1}, map[string]string{"x": "y"}, store.Cond{NotExists: true}); err != nil {
			t.Fatal(err)
		}
	}
	for i := range 4 {
		id := fmt.Sprintf("grt_done_%d", i)
		save(&Grant{ID: id, Status: StatusDisconnected, SignedIn: true})
		row(skOp(id, "op1"))
		row(skOp(id, "op2"))
	}
	row(skSnapshot("grt_done_0", 1))
	save(&Grant{ID: "grt_failed", Status: StatusFailed})
	row(skOp("grt_failed", "op1"))
	// A provider access confirmed disabled: its record may go.
	save(&Grant{ID: "grt_revoked", Status: StatusDisconnected, SignedIn: true,
		Connector: &ConnectorState{Account: "acct_1", Revocation: &connector.Revocation{Disabled: true}}})
	row(skOp("grt_revoked", "op1"))
	// Skipped: a finished grant that still has a token row, and one whose
	// provider access isn't confirmed disabled.
	save(&Grant{ID: "grt_token", Status: StatusFailed})
	row(skToken("grt_token"))
	row(skOp("grt_token", "op1"))
	save(&Grant{ID: "grt_unrevoked", Status: StatusDisconnected, SignedIn: true, Connector: &ConnectorState{Account: "acct_2"}})
	row(skOp("grt_unrevoked", "op1"))
	// Kept: waiting for its sign-in (recent), signed in a month ago and held,
	// ended by the vendor (needs_reconnect), and listed with triggers.
	save(&Grant{ID: "grt_waiting", Status: StatusApproved})
	save(&Grant{ID: "grt_held", Status: StatusApproved, SignedIn: true,
		Request: Request{CreatedAt: old.UTC().Format(time.RFC3339)}})
	save(&Grant{ID: "grt_reconnect", Status: StatusNeedsReconnect, SignedIn: true})
	save(&Grant{ID: "grt_triggers", Status: StatusListed, SignedIn: true, TriggerIDs: []string{"trg_1"}})
	row(skOp(live.ID, "keep1"))
	row(skOp("grt_waiting", "keep2"))
	// The due index: a finished grant's stale row goes, a live one's stays.
	due := func(id string) string { return dueSK(time.Now().Add(time.Hour), tenant, id) }
	for _, id := range []string{"grt_done_0", "grt_waiting"} {
		if err := e.kv.Put(ctx, store.Item{PK: pkDue, SK: due(id), Doc: []byte("{}")}, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}

	count := func(prefix string) int {
		all, err := e.kv.Query(ctx, tenantPK(tenant), prefix)
		if err != nil {
			t.Fatal(err)
		}
		return len(all)
	}
	grantsBefore, opsBefore, snapsBefore := count("GRANT#"), count("OP#"), count("SNAPSHOT#")

	// Not without the signature of someone who may delete the account.
	if st, _ := e.do("DELETE", acct, nil, nil, "k_setup"); st != 401 {
		t.Fatalf("an agent's key: %d, want 401", st)
	}
	// Nor another account's key.
	if st, _ := e.do("DELETE", acct, nil, nil, "k_other"); st != 401 && st != 403 {
		t.Fatalf("another account's key: %d, want 401 or 403", st)
	}
	editor := e.device("dev_editor")
	ed := e.kiff.Devices["dev_editor"]
	ed.Capabilities = []string{"holds.answer", "tools.connect"}
	e.kiff.Devices["dev_editor"] = ed
	if st, _ := e.do("DELETE", acct, nil, editor, ""); st != 403 {
		t.Fatalf("an editor: %d, want 403", st)
	}
	if count("GRANT#") != grantsBefore || count("OP#") != opsBefore {
		t.Fatal("a refused request deleted records")
	}

	st, out := e.do("DELETE", acct, nil, e.phone, "")
	if st != 200 {
		t.Fatalf("delete: %d %v", st, out)
	}
	// 4 disconnected + 1 failed + 1 revoked connector: 6 grants, 4*2+1+1 = 10 operations.
	if out["grants"] != float64(6) || out["operations"] != float64(10) || out["skipped"] != float64(2) || out["more"] != false ||
		out["grants_before"] != float64(grantsBefore) || out["grants_after"] != float64(grantsBefore-6) {
		t.Fatalf("the counts: %v (before %d grants)", out, grantsBefore)
	}
	if count("GRANT#") != grantsBefore-6 || count("OP#") != opsBefore-10 || count("SNAPSHOT#") != snapsBefore-1 {
		t.Fatalf("kept %d grants, %d operations, %d snapshots", count("GRANT#"), count("OP#"), count("SNAPSHOT#"))
	}
	has := func(sk string) bool {
		_, ok, err := e.kv.Get(ctx, tenantPK(tenant), sk)
		if err != nil {
			t.Fatal(err)
		}
		return ok
	}
	for _, id := range []string{"grt_done_0", "grt_done_3", "grt_failed", "grt_revoked"} {
		if has(skGrant(id)) || has(skLive(id)) {
			t.Errorf("%s was kept", id)
		}
	}
	// Everything else stays: the live grant (and it still works), the
	// waiting, held, reconnecting and triggered ones, the skipped ones and
	// their token and operations, the live ones' operations, the enrolment.
	for _, k := range []string{skGrant(live.ID), skGrant("grt_waiting"), skGrant("grt_held"), skGrant("grt_reconnect"), skGrant("grt_triggers"),
		skGrant("grt_token"), skToken("grt_token"), skOp("grt_token", "op1"), skGrant("grt_unrevoked"), skOp("grt_unrevoked", "op1"),
		skOp(live.ID, "keep1"), skOp("grt_waiting", "keep2"), skEnrolment} {
		if !has(k) {
			t.Errorf("%s was deleted", k)
		}
	}
	if g, err := e.s.loadGrant(ctx, tenant, live.ID); err != nil || g.Status != StatusConnected {
		t.Fatalf("the live grant changed: %v %v", err, g)
	}
	// The due index: the finished grant's row is gone, the waiting one's stays.
	if _, ok, _ := e.kv.Get(ctx, pkDue, due("grt_done_0")); ok {
		t.Error("a deleted grant's due row was kept")
	}
	if _, ok, _ := e.kv.Get(ctx, pkDue, due("grt_waiting")); !ok {
		t.Error("a kept grant's due row was deleted")
	}
	// One account event says who signed and the counts.
	events, _ := e.kv.Query(ctx, tenantPK(tenant), "ACCOUNT_EVENT#")
	if len(events) != 1 || !strings.Contains(string(events[0].Doc), `"grants":6`) || !strings.Contains(string(events[0].Doc), `"skipped":2`) ||
		!strings.Contains(string(events[0].Doc), `"by":"`) || strings.Contains(string(events[0].Doc), `"by":""`) {
		t.Fatalf("the account event: %v", events)
	}
	// Asked again, only the skipped ones are left, and they are skipped again.
	if st, out := e.do("DELETE", acct+"?again=1", nil, e.phone, ""); st != 200 || out["grants"] != float64(0) || out["skipped"] != float64(2) {
		t.Fatalf("again: %d %v", st, out)
	}
}

// One call deletes at most finishedPerCall grants and says more: true, so the
// next call carries on: 960 grants and 8,004 operations are too much for one
// function's time.
func TestDeletingFinishedGrantsStopsAtTheCapAndSaysMore(t *testing.T) {
	e := newEnv(t)
	e.connect()
	ctx := context.Background()
	acct := "/v1/accounts/" + tenant + "/finished-grants"
	old := finishedPerCall
	finishedPerCall = 4
	t.Cleanup(func() { finishedPerCall = old })
	for i := range 10 {
		g := &Grant{ID: fmt.Sprintf("grt_done_%d", i), Tenant: tenant, Status: StatusDisconnected, SignedIn: true,
			Request: Request{Tenant: tenant, CatalogEntry: "gmail-key@1"}}
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	for i, want := range []struct {
		grants float64
		more   bool
	}{{4, true}, {4, true}, {2, false}, {0, false}} {
		st, out := e.do("DELETE", fmt.Sprintf("%s?call=%d", acct, i+1), nil, e.phone, "")
		if st != 200 || out["grants"] != want.grants || out["more"] != want.more {
			t.Fatalf("call %d: %d %v, want %v grants and more=%v", i+1, st, out, want.grants, want.more)
		}
	}
}
