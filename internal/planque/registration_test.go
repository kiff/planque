package planque

import (
	"context"
	"testing"
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

// KIFF refusing an update (it never lowers a tool's class) leaves the grant unstamped, so the refusal shows again;
// the call it ran in isn't stopped.
func TestARefusedUpdateLeavesTheGrantUnstamped(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.kiff.RefusePatch = 409
	e.s.refreshRegistration(context.Background(), g)
	var rec registrationRecord
	if _, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skRegistration(g.ID), &rec); err == nil || e.kiff.Patched != 0 {
		t.Fatalf("stamped after a refusal: %+v, %d updated", rec, e.kiff.Patched)
	}
}
