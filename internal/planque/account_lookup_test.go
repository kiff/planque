package planque

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var reads atomic.Int64

// readGrants is the owner's app reading the account's grants, signed by the enrolled phone, d after the test began
// (on planque's clock, the table's and the phone's, each read signed a second later than the last).
func readGrants(t *testing.T, e *env, d time.Duration) {
	t.Helper()
	e.at(d)
	e.kv.Now = func() time.Time { return e.s.now() }
	e.phone.Now = func() time.Time { return time.Now().Add(d + time.Duration(reads.Add(1))*time.Second) }
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants", nil, e.phone, ""); st != 200 {
		t.Fatalf("the owner's read: %d %v", st, out)
	}
}

func statusOf(t *testing.T, e *env, g *Grant) *Grant {
	t.Helper()
	got, err := e.s.loadGrant(context.Background(), tenant, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// A connection nobody calls reads "Connected" after the provider lost its account (an Outlook taf whose schedule
// paused, 17 Nov). The owner's read of the grants looks the account up, and a gone one asks to sign in again.
// Category: reliability
// Attack: a deleted provider account kept reading "Connected" to the owner, so a taf was never asked to sign in again
// Refused: the owner's read ends the sign-in as a call would: the connection is signed out, with Connect again
func TestTheOwnersReadEndsAConnectionTheProviderLost(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	f.accountGone = true
	readGrants(t, e, 0)
	got := statusOf(t, e, g)
	if got.Status != StatusNeedsReconnect || !strings.Contains(got.Reason, "no longer accepts this sign-in") {
		t.Fatalf("after the owner's read: %s, %q", got.Status, got.Reason)
	}
}

// A healthy account is left connected, and it is looked up once an hour however often the app reads.
func TestAnAccountIsLookedUpOncePerHour(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	base := f.accountReads
	readGrants(t, e, 0)
	if f.accountReads != base+1 || statusOf(t, e, g).Status != StatusConnected {
		t.Fatalf("the first read: %d lookups, %s", f.accountReads-base, statusOf(t, e, g).Status)
	}
	for range 5 {
		readGrants(t, e, 0)
	}
	if f.accountReads != base+1 {
		t.Fatalf("five more reads looked up %d more times", f.accountReads-base-1)
	}
	// Within the hour a lost account isn't learned from a read; the next hour's read does.
	f.accountGone = true
	readGrants(t, e, 30*time.Minute)
	if statusOf(t, e, g).Status != StatusConnected || f.accountReads != base+1 {
		t.Fatalf("at 30 minutes: %s, %d lookups", statusOf(t, e, g).Status, f.accountReads-base)
	}
	readGrants(t, e, 61*time.Minute)
	if got := statusOf(t, e, g); got.Status != StatusNeedsReconnect || f.accountReads != base+2 {
		t.Fatalf("at 61 minutes: %s, %d lookups", got.Status, f.accountReads-base)
	}
}

// A lookup the provider doesn't answer (503) or answers too late changes nothing, and the read still answers.
// Category: reliability
// Attack: a provider outage on the owner's read signs every connection out
// Refused: a lookup that fails or times out leaves the grant connected
func TestALookupThatFailsOrTimesOutChangesNothing(t *testing.T) {
	t.Run("the provider is down", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		f.accountDown = true
		readGrants(t, e, 0)
		if got := statusOf(t, e, g); got.Status != StatusConnected {
			t.Fatalf("after a 503: %s", got.Status)
		}
	})
	t.Run("the provider is slow", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		was := accountLookupWithin
		accountLookupWithin = 50 * time.Millisecond
		t.Cleanup(func() { accountLookupWithin = was })
		f.accountGone, f.accountDelay = true, time.Minute
		start := time.Now()
		readGrants(t, e, 0)
		if time.Since(start) > 5*time.Second {
			t.Fatalf("the read waited %s", time.Since(start))
		}
		if got := statusOf(t, e, g); got.Status != StatusConnected {
			t.Fatalf("after a lookup that timed out: %s", got.Status)
		}
	})
}

// A read looks nothing up for a connection that isn't connected: a signed-out one has no account to ask about.
func TestAReadLooksUpOnlyConnectedAccounts(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	f.accountGone = true
	readGrants(t, e, 0)
	if statusOf(t, e, g).Status != StatusNeedsReconnect {
		t.Fatal("not signed out")
	}
	before := f.accountReads
	readGrants(t, e, 2*time.Hour)
	if f.accountReads != before {
		t.Fatalf("a signed-out connection was looked up %d times", f.accountReads-before)
	}
}
