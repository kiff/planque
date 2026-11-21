package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
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

// Who still needs a vendor account, asked of the store as the removal asks
// it: connections of one tenant, of the same provider and family of apps, on
// one account or one that names none; never another tenant's, another
// family's, a different named account's, or a disconnected or failed one
// (qa-kiff's review).
func TestWhoStillUsesAVendorAccountIsReadForTheTenantAndTheFamily(t *testing.T) {
	e := newInProcessEnv(t)
	ctx := context.Background()
	put := func(id, tnt, toolkit, name, status string) *Grant {
		g := &Grant{ID: id, Tenant: tnt, Status: status,
			Request:   Request{Server: "https://backend.composio.dev", Connector: &ConnectorRequest{Provider: "composio", Toolkit: toolkit}},
			Connector: &ConnectorState{Account: "acct_" + id, AccountName: name}}
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
		return g
	}
	users := func(g *Grant) int {
		t.Helper()
		ids, err := e.s.signInUsers(ctx, g)
		if err != nil {
			t.Fatal(err)
		}
		return len(ids)
	}
	// Google, one tenant: Drive removed, Docs on the same account, Gmail on another named account.
	drive := put("grt_drive", tenant, "googledrive", "carolina@kiff.dev", StatusConnected)
	put("grt_docs", tenant, "googledocs", "Carolina@kiff.dev", StatusConnected)
	put("grt_gmail_other", tenant, "gmail", "someone@else.example", StatusConnected)
	if n := users(drive); n != 1 {
		t.Fatalf("Google: %d users of the account, want 1 (Docs, same name in another case)", n)
	}
	// A disconnected or failed one uses nothing; another tenant's never counts; an unnamed one may.
	put("grt_sheets_gone", tenant, "googlesheets", "carolina@kiff.dev", StatusDisconnected)
	put("grt_calendar_failed", tenant, "googlecalendar", "carolina@kiff.dev", StatusFailed)
	put("grt_docs_other_tenant", "tnt_other", "googledocs", "carolina@kiff.dev", StatusConnected)
	if n := users(drive); n != 1 {
		t.Fatalf("Google: %d users after a disconnected, a failed and another tenant's, want 1", n)
	}
	put("grt_sheets_unnamed", tenant, "googlesheets", "", StatusConnected)
	if n := users(drive); n != 2 {
		t.Fatalf("Google: %d users with an unnamed one live, want 2", n)
	}
	// Microsoft, in the same store: Outlook removed, OneDrive on the same account, Drive (Google) never.
	outlook := put("grt_outlook", tenant, "outlook", "carolona@kiff.dev", StatusConnected)
	put("grt_onedrive", tenant, "one_drive", "CAROLONA@kiff.dev", StatusConnected)
	if n := users(outlook); n != 1 {
		t.Fatalf("Microsoft: %d users of the account, want 1 (OneDrive); Google's apps never count", n)
	}
	// The last one on its account: nothing else of its family on the account.
	lone := put("grt_typeform", tenant, "typeform", "", StatusConnected)
	if n := users(lone); n != 0 {
		t.Fatalf("an app of its own family with no other connection: %d users, want 0", n)
	}
}

// A removed connection whose sign-in is kept for another connection serves
// no call at all: the call is refused before anything reaches the provider,
// even with a valid permit, and its relay credential is gone (legal-i18n's
// question on the retention text: the kept sign-in is "unused").
// Category: authority
// Attack: a removed taf's connection, its sign-in kept for another's, still sends with the owner's account
// Refused: "Not sent", nothing reaches the provider, no relay credential
func TestARemovedConnectionWhoseSignInIsKeptServesNoCalls(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a, b := twoOwnSignIns(t, e, f)
	remove(t, e, a)
	if f.revocations != 0 {
		t.Fatalf("setup: the first removal revoked: %d", f.revocations)
	}
	g, _ := e.s.loadGrant(ctx, tenant, a.ID)
	if g.Status != StatusDisconnected || g.Transport != "" || len(g.Transports) != 0 || g.Connected != nil {
		t.Fatalf("the removed connection still holds a way in: %s %q %v %v", g.Status, g.Transport, g.Transports, g.Connected)
	}
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_removed"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_removed", Args: args})
	executions := f.executions
	got := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if !strings.Contains(got, `"code":"`+CodeConnectionDown+`"`) || !strings.Contains(got, `"sent":false`) || f.executions != executions {
		t.Fatalf("a call on the removed connection: %s (provider executions %d→%d)", got, executions, f.executions)
	}
	// The other connection is untouched, and its provider account was not revoked.
	if live, _ := e.s.loadGrant(ctx, tenant, b.ID); live.Status != StatusConnected || f.revocations != 0 {
		t.Fatalf("the other connection after the first's removal: %s, %d revocations", live.Status, f.revocations)
	}
}

// Category: authority
// Attack: a live connection whose LIVE# row was never written (keepLive is its own write) is missed, so removing its sibling revokes the sign-in both use
// Refused: the check before a revoke reads every grant, so the sign-in is kept while the other connection lives
func TestALiveConnectionMissingFromTheLiveIndexKeepsTheSignIn(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a := connectConnector(t, e, f)
	b := connectConnector(t, e, f) // the same sign-in: the fixture gives both one connected account
	if err := e.s.ensureLive(ctx, tenant); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := e.s.KV.Delete(ctx, tenantPK(tenant), skLive(b.ID)); err != nil || !ok {
		t.Fatalf("setup: drop b's LIVE# row: %v %v", ok, err)
	}
	remove(t, e, a)
	if f.revocations != 0 {
		t.Fatalf("removing a revoked the sign-in the live b still uses: %d revocations", f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("b after a's removal: %s", g.Status)
	}
	if n, err := e.s.SweepKept(ctx); err != nil || n != 0 || f.revocations != 0 {
		t.Fatalf("sweep with b live: %d, %v, %d revocations", n, err, f.revocations)
	}
}
