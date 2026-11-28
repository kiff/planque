package planque

import (
	"context"
	"testing"
	"time"
)

// A grant registered before KIFF read its tools' app, hint and pins is updated in KIFF in place, once: every tool,
// then the grant is stamped, and the next call or read sends nothing. Registering again is never tried, since KIFF
// refuses a name already connected and removing it would drop it from under the owner's Cards.
func TestAnOlderRegistrationIsUpdatedInPlaceOnce(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	e.s.refreshRegistration(ctx, g)
	if e.kiff.Patched != len(g.Tools) {
		t.Fatalf("updated %d tools, want %d", e.kiff.Patched, len(g.Tools))
	}
	var rec registrationRecord
	if _, err := getDoc(ctx, e.s.KV, tenantPK(tenant), skRegistration(g.ID), &rec); err != nil || rec.Version != kiffRegistration {
		t.Fatalf("not stamped: %+v %v", rec, err)
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if now.ver != g.ver {
		t.Fatalf("the grant's own record was written: version %d, was %d", now.ver, g.ver)
	}
	e.s.refreshRegistration(ctx, now)
	if e.kiff.Patched != len(g.Tools) {
		t.Fatalf("updated again: %d", e.kiff.Patched)
	}
	if got := len(e.kiff.Connections(tenant)); got != len(g.Tools) {
		t.Fatalf("connections %d after the update, want %d: nothing removed or added", got, len(g.Tools))
	}
}

// KIFF refusing an update for good (it never lowers a tool's class) is final for the tool: one PATCH, stamped, and no
// later call sends another.
func TestARefusedUpdateIsSentOnce(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.kiff.RefusePatch = 409
	for range 3 {
		e.s.refreshRegistration(context.Background(), g)
	}
	if e.kiff.PatchTries != len(g.Tools) {
		t.Fatalf("%d PATCHes for %d tools over three calls, want one each", e.kiff.PatchTries, len(g.Tools))
	}
	var rec registrationRecord
	if _, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skRegistration(g.ID), &rec); err != nil || rec.Version != kiffRegistration {
		t.Fatalf("not stamped after a final refusal: %+v, %v", rec, err)
	}
}

// A transient failure (KIFF down) leaves the older version and waits registrationRetry: calls within it send
// nothing, and the first after it tries again.
func TestATransientFailureWaitsBeforeTryingAgain(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.kiff.RefusePatch = 503
	e.s.refreshRegistration(context.Background(), g)
	e.s.refreshRegistration(context.Background(), g)
	if e.kiff.PatchTries != len(g.Tools) {
		t.Fatalf("%d PATCHes within the wait, want %d", e.kiff.PatchTries, len(g.Tools))
	}
	e.kiff.RefusePatch = 0
	later := time.Now().Add(registrationRetry + time.Minute)
	e.s.Now = func() time.Time { return later }
	e.s.refreshRegistration(context.Background(), g)
	var rec registrationRecord
	if _, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skRegistration(g.ID), &rec); err != nil || rec.Version != kiffRegistration || e.kiff.Patched != len(g.Tools) {
		t.Fatalf("after the wait: %+v, %v, %d updated", rec, err, e.kiff.Patched)
	}
}

// A tool KIFF says isn't the relay's to update (403) is said once: the grant is stamped, so no call asks again.
func TestAnUpdateThatIsntOursIsNotAskedAgain(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.kiff.RefusePatch = 403
	e.s.refreshRegistration(context.Background(), g)
	var rec registrationRecord
	if _, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skRegistration(g.ID), &rec); err != nil || rec.Version != kiffRegistration {
		t.Fatalf("not stamped after 403: %+v, %v", rec, err)
	}
}
