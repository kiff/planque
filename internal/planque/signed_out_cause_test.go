package planque

import (
	"context"
	"strings"
	"testing"
	"time"
)

// A signed-out connection says why, from the cause planque kept, never as
// "expired" unless the sign-in really ran out (a Drive sign-in revoked by
// another connection's removal read "expired").
//
// Category: honesty
// Attack: a revoke worded as an expiry, and a removed shared sign-in worded as the vendor's doing
// Refused: only no_refresh_token reads expired; the vendor's refusal reads "by {vendor}", a removed shared sign-in and an unknown cause read plain
func TestASignedOutConnectionSaysItsCause(t *testing.T) {
	cases := []struct {
		name, cause string
		byVendor    bool
		expired     bool
	}{
		{"the vendor refused", CauseRefreshRefused, true, false},
		{"it really ran out", CauseNoRefreshToken, true, true},
		{"its shared sign-in was removed", CauseSharedSignInGone, false, false},
		{"a cause planque doesn't know", "some_future_cause", false, false},
		{"signed out before the cause was kept", "", false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e, f := newConnectorFixture(t)
			g := connectConnector(t, e, f)
			ctx := context.Background()
			if err := e.s.signOut(ctx, g, "why", c.cause); err == nil {
				t.Fatal("a signed-out grant's call should report it")
			}
			now, _ := e.s.loadGrant(ctx, tenant, g.ID)
			vendor := vendorName(now)
			if now.Status != StatusNeedsReconnect || now.SignedOutCause != c.cause {
				t.Fatalf("grant: %s, cause %q, want %q", now.Status, now.SignedOutCause, c.cause)
			}
			_, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup")
			words, _ := view["status_words"].(string)
			if c.byVendor != strings.HasPrefix(words, "Signed out by "+vendor) {
				t.Fatalf("status words %q, by vendor: %v", words, c.byVendor)
			}
			if !c.byVendor && words != "Signed out · sign in again" {
				t.Fatalf("a cause with no vendor to name reads %q", words)
			}
			if got, _ := view["signed_out_cause"].(string); got != c.cause {
				t.Fatalf("the view's cause %q, want %q", got, c.cause)
			}
			h, _ := health(listRow(t, e, g.ID, "Europe/Paris"))["words"].(string)
			if c.expired != strings.Contains(h, "expired") {
				t.Fatalf("row words %q, expired: %v", h, c.expired)
			}
			if strings.Contains(strings.ToLower(words), "expired") {
				t.Fatalf("the status line says expired: %q", words)
			}
		})
	}
}

// A connection that isn't signed out sends no cause, and one signed in again
// forgets it.
func TestTheCauseGoesWithTheSignOut(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	if c := signedOutCauseOf(g); c != "" {
		t.Fatalf("a connected grant has cause %q", c)
	}
	_ = e.s.signOut(ctx, g, "why", CauseRefreshRefused)
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if signedOutCauseOf(now) != CauseRefreshRefused {
		t.Fatalf("signed out: cause %q", signedOutCauseOf(now))
	}
	now.Status, now.SignedOutAt, now.SignedOutCause = StatusConnected, time.Time{}, ""
	if signedOutCauseOf(now) != "" {
		t.Fatalf("connected again: cause %q", signedOutCauseOf(now))
	}
}
