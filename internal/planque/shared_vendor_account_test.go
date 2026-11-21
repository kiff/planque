package planque

import (
	"context"
	"testing"
)

// 21 Nov: removing a taf's Drive connection signed another taf's Drive
// connection out, though each was its own sign-in with its own connected
// account at the provider: revoking at Google ends the user's whole grant to
// the OAuth client, so every other token of that Google account died too.
// A removal now revokes at the vendor only when no other live connection of
// the same family of apps may be on the same vendor account.

func twoOwnSignIns(t *testing.T, e *env, f *connectorFixture) (a, b *Grant) {
	t.Helper()
	a = connectConnector(t, e, f)
	b = connectConnector(t, e, f)
	// Each its own connected account at the provider, as two sign-ins made
	// apart are (the fixture gives both the same one).
	b.Connector.Account = "account_demo_2"
	if err := e.s.saveGrant(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	if a.Connector.Account == b.Connector.Account {
		t.Fatal("setup: the two connections share a connected account")
	}
	return a, b
}

func remove(t *testing.T, e *env, g *Grant) {
	t.Helper()
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("remove %s: %d %v", g.ID, st, out)
	}
}

// Category: authority
// Attack: removing one taf's connection revokes the vendor account at the provider, signing out the owner's other taf on the same account
// Refused: no revocation while the other connection is live; the last removal revokes
func TestRemovingOneOfTwoConnectionsOnTheSameVendorAccountSignsNoOneOut(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a, b := twoOwnSignIns(t, e, f) // both on "Demo Shop", each its own connected account
	remove(t, e, a)
	if f.revocations != 0 {
		t.Fatalf("removing one connection revoked the vendor account the other is on: %d revocations", f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("the other connection after the first's removal: %s", g.Status)
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 || f.revocations != 0 {
		t.Fatalf("sweep with the other live: %d, %v, %d revocations", n, err, f.revocations)
	}
	remove(t, e, b)
	if f.revocations != 1 {
		t.Fatalf("the last connection removed: %d revocations, want 1", f.revocations)
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n != 1 || f.revocations != 2 {
		t.Fatalf("the kept first connection is revoked once none is live: %d, %v, %d revocations", n, err, f.revocations)
	}
}

// Category: authority
// Attack: a connection that can't say which account it is on is revoked under another that may be on it
// Refused: an unnamed connection is taken as possibly the other's, so nothing is revoked while the other lives
func TestAnUnnamedConnectionIsTakenAsPossiblyOnTheSameAccount(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a, b := twoOwnSignIns(t, e, f)
	b.Connector.AccountName = ""
	if err := e.s.saveGrant(ctx, b); err != nil {
		t.Fatal(err)
	}
	remove(t, e, a)
	if f.revocations != 0 {
		t.Fatalf("revoked under a connection that names no account: %d", f.revocations)
	}
}

// Two connections named as different accounts don't share: removing one
// revokes it, as before (the owner's "remove access" still ends the access).
func TestConnectionsOnDifferentNamedAccountsAreRevokedOnRemoval(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a, b := twoOwnSignIns(t, e, f)
	b.Connector.AccountName = "Another Shop"
	if err := e.s.saveGrant(ctx, b); err != nil {
		t.Fatal(err)
	}
	remove(t, e, a)
	if f.revocations != 1 {
		t.Fatalf("a connection on a different account was kept: %d revocations", f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("the other account's connection: %s", g.Status)
	}
}

// The families: Google's apps share a vendor account across apps, Microsoft's
// likewise, and an app outside both is its own family.
func TestTheAppsThatShareAVendorAccountAreOneFamily(t *testing.T) {
	grant := func(toolkit, name string) *Grant {
		return &Grant{Request: Request{Connector: &ConnectorRequest{Provider: "composio", Toolkit: toolkit}}, Connector: &ConnectorState{AccountName: name}}
	}
	for _, c := range []struct {
		a, b *Grant
		want bool
	}{
		{grant("googledrive", "x@y.z"), grant("googledocs", "X@Y.Z"), true},
		{grant("gmail", "x@y.z"), grant("googlecalendar", ""), true},
		{grant("outlook", ""), grant("one_drive", ""), true},
		{grant("googledrive", "x@y.z"), grant("googledocs", "other@y.z"), false},
		{grant("googledrive", ""), grant("outlook", ""), false},
		{grant("typeform", ""), grant("shopify", ""), false},
		{grant("typeform", ""), grant("typeform", ""), true},
	} {
		if got := mayShareVendorAccount(c.a, c.b); got != c.want {
			t.Errorf("%s (%q) with %s (%q): %v, want %v", c.a.Request.Connector.Toolkit, c.a.Connector.AccountName, c.b.Request.Connector.Toolkit, c.b.Connector.AccountName, got, c.want)
		}
	}
	other := grant("gmail", "")
	other.Request.Connector.Provider = "another"
	if mayShareVendorAccount(grant("gmail", ""), other) {
		t.Error("two providers' connections are never one account")
	}
}
