package planque

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/store"
)

// An account that connects and disconnects a lot ends up holding hundreds of
// finished grants and their operation records (the Eval's held 960 grants and
// 8,004 operations), and every read of its grants pays for them. With the
// signature of someone who may delete the account, the finished ones go, with
// their own operation records, and nothing live: no live grant, no live
// grant's operations, not the enrolment.
func TestDeletingFinishedGrantsKeepsTheLiveAndTheEnrolment(t *testing.T) {
	e := newEnv(t)
	live := e.connect()
	ctx := context.Background()
	acct := "/v1/accounts/" + tenant + "/finished-grants"

	save := func(id, status string, signedIn bool, created time.Time) {
		g := &Grant{ID: id, Tenant: tenant, Status: status, SignedIn: signedIn,
			Request: Request{Tenant: tenant, CatalogEntry: "gmail-key@1", CreatedAt: created.UTC().Format(time.RFC3339)}}
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	op := func(grant, name string) {
		it := store.Item{PK: tenantPK(tenant), SK: skOp(grant, name), Ver: 1}
		if err := putDoc(ctx, e.kv, it, map[string]string{"state": "done"}, store.Cond{NotExists: true}); err != nil {
			t.Fatal(err)
		}
	}
	old, recent := time.Now().Add(-45*24*time.Hour), time.Now().Add(-2*time.Hour)
	for i := range 4 {
		save(fmt.Sprintf("grt_done_%d", i), StatusDisconnected, true, recent)
		op(fmt.Sprintf("grt_done_%d", i), "op1")
		op(fmt.Sprintf("grt_done_%d", i), "op2")
	}
	save("grt_failed", StatusFailed, false, recent)
	op("grt_failed", "op1")
	// Finished by rule: asked for over a month ago, never signed in, no token.
	save("grt_stale", StatusApproved, false, old)
	op("grt_stale", "op1")
	// Not finished: asked for hours ago and waiting for its sign-in, and one
	// signed in a month ago and still held.
	save("grt_waiting", StatusApproved, false, recent)
	save("grt_held", StatusApproved, true, old)
	op(live.ID, "keep1")
	op("grt_waiting", "keep2")

	count := func() (grants, ops int) {
		all, err := e.kv.Query(ctx, tenantPK(tenant), "")
		if err != nil {
			t.Fatal(err)
		}
		for _, it := range all {
			switch {
			case strings.HasPrefix(it.SK, "GRANT#"):
				grants++
			case strings.HasPrefix(it.SK, "OP#"):
				ops++
			}
		}
		return
	}
	grantsBefore, opsBefore := count()

	// Not without a signature of someone who may delete the account.
	if st, _ := e.do("DELETE", acct, nil, nil, "k_setup"); st != 401 {
		t.Fatalf("an agent's key: %d, want 401", st)
	}
	editor := e.device("dev_editor")
	ed := e.kiff.Devices["dev_editor"]
	ed.Capabilities = []string{"holds.answer", "tools.connect"}
	e.kiff.Devices["dev_editor"] = ed
	if st, _ := e.do("DELETE", acct, nil, editor, ""); st != 403 {
		t.Fatalf("an editor: %d, want 403", st)
	}
	if g, o := count(); g != grantsBefore || o != opsBefore {
		t.Fatalf("a refused request deleted records: %d grants and %d operations, were %d and %d", g, o, grantsBefore, opsBefore)
	}

	st, out := e.do("DELETE", acct, nil, e.phone, "")
	if st != 200 {
		t.Fatalf("delete: %d %v", st, out)
	}
	// Six finished grants (4 disconnected, 1 failed, 1 stale) with 4*2+1+1 operations.
	if out["grants"] != float64(6) || out["operations"] != float64(10) ||
		out["grants_before"] != float64(grantsBefore) || out["grants_after"] != float64(grantsBefore-6) {
		t.Fatalf("the counts: %v (before %d grants)", out, grantsBefore)
	}
	grantsAfter, opsAfter := count()
	if grantsAfter != grantsBefore-6 || opsAfter != opsBefore-10 {
		t.Fatalf("kept %d grants and %d operations, want %d and %d", grantsAfter, opsAfter, grantsBefore-6, opsBefore-10)
	}
	for _, id := range []string{"grt_done_0", "grt_done_3", "grt_failed", "grt_stale"} {
		if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), skGrant(id)); ok {
			t.Errorf("%s was kept", id)
		}
	}
	// Everything not finished stays: the live grant (and still works), the
	// waiting one, the held one, their operations, and the enrolment.
	for _, id := range []string{live.ID, "grt_waiting", "grt_held"} {
		if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), skGrant(id)); !ok {
			t.Errorf("%s was deleted", id)
		}
	}
	for _, k := range []string{skOp(live.ID, "keep1"), skOp("grt_waiting", "keep2")} {
		if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), k); !ok {
			t.Errorf("%s was deleted", k)
		}
	}
	if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), skEnrolment); !ok {
		t.Fatal("the enrolment was deleted")
	}
	if code := e.relayCallRaw(live); code == 401 {
		t.Fatalf("the live grant stopped working: %d", code)
	}
	// Asked again, there is nothing left to delete.
	if st, out := e.do("DELETE", acct, nil, e.phone, ""); st != 200 || out["grants"] != float64(0) || out["operations"] != float64(0) {
		t.Fatalf("again: %d %v", st, out)
	}
}
