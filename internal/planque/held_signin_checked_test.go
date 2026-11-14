package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
)

// askEntry asks for the managed app by its entry for a new taf, as Coach
// does, without new_sign_in: a held sign-in is reused when it's fit to be.
func askEntry(t *testing.T, e *env, agent string) map[string]any {
	t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": agent, "agent_name": "Refund requests"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	return out["grant"].(map[string]any)
}

// notReused says the request asks for a sign-in of its own: not the held
// one, and never "no second sign-in".
func notReused(t *testing.T, v map[string]any, held *Grant) {
	t.Helper()
	text, _ := v["text"].(string)
	if v["signed_in"] == true || v["status"] != StatusPending || strings.Contains(text, "no second sign-in") || strings.Contains(text, "already have there") {
		t.Fatalf("the held sign-in was offered: %v", v)
	}
	if strings.Contains(text, held.Connector.Account) {
		t.Fatalf("the provider's id reached the owner: %s", text)
	}
}

// A held sign-in the provider no longer calls active (expired, revoked) isn't
// reused: the new taf's connection asks for a sign-in, and the held one is
// marked signed out, so Connect apps says "Sign in again", never "Connected".
// Category: authority
// Attack: an expired held sign-in reused, with "no second sign-in" said over a sign-in the vendor refuses
// Refused: not reused; the owner is asked to sign in, and the held connection reads signed out
func TestAnExpiredHeldSignInIsntReused(t *testing.T) {
	e, f := newConnectorFixture(t)
	first := connectConnector(t, e, f)
	f.accountStatus = "EXPIRED"
	notReused(t, askEntry(t, e, "agt_refunds"), first)
	g, _ := e.s.loadGrant(context.Background(), tenant, first.ID)
	if g.Status != StatusNeedsReconnect || !strings.Contains(g.Reason, "no longer accepts this sign-in") {
		t.Fatalf("the held connection: %s, %q", g.Status, g.Reason)
	}
	_, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, nil, "k_setup")
	if words, _ := view["status_words"].(string); words == "Connected" || strings.Contains(words, "account_demo") {
		t.Fatalf("Connect apps would read: %q", words)
	}
}

// The vendor refuses the sign-in itself (401) while the provider still calls
// it active: every connection on that sign-in is marked signed out, and a new
// taf isn't offered it.
// Category: authority
// Attack: a sign-in Gmail answers 401 to, behind an ACTIVE status, reused for the next taf
// Refused: every connection on it signs in again, and the next request asks for a sign-in
func TestAVendorsRefusedSignInSignsOutEveryConnectionOnIt(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	first := connectConnector(t, e, f)
	f.executeBody = `{"data":{"status_code":401,"message":"Invalid Credentials"},"successful":false,"error":"Request had invalid authentication credentials."}`
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_401"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: first.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_401", Args: args})
	got := mustJSONString(e.s.call(ctx, first, "reply_ticket", args, permit))
	if !strings.Contains(got, "401") {
		t.Fatalf("the taf reads the vendor's refusal: %s", got)
	}
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	if g.Status != StatusNeedsReconnect {
		t.Fatalf("after the vendor's 401: %s", g.Status)
	}
	notReused(t, askEntry(t, e, "agt_refunds"), first)
}

// A refusal that isn't of the sign-in (the vendor's 403 on one item) leaves
// the connection as it is, and a held sign-in the provider can't be asked
// about now isn't reused but isn't marked either: it's asked again next time.
func TestOnlyARefusedSignInSignsOut(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	first := connectConnector(t, e, f)
	f.executeBody = `{"data":{"status_code":403,"message":"Forbidden"},"successful":false,"error":"Not allowed on this ticket."}`
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_403"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: first.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_403", Args: args})
	e.s.call(ctx, first, "reply_ticket", args, permit)
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusConnected {
		t.Fatalf("after a 403: %s", g.Status)
	}
	f.accountDown = true
	notReused(t, askEntry(t, e, "agt_refunds"), first)
	if g, _ := e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusConnected {
		t.Fatalf("after the provider didn't answer: %s", g.Status)
	}
	// Healthy again: reused, with no new sign-in.
	f.accountDown = false
	if v := askEntry(t, e, "agt_third"); v["signed_in"] != true || v["status"] != StatusListed {
		t.Fatalf("a healthy held sign-in wasn't reused: %v", v)
	}
}

// reapprove is the approve body the phone sends for g as it stands: its
// digest and, for its profile, the digests the owner read.
func reapprove(g *Grant) map[string]any {
	ph, _ := profile.Digest(g.Request.Profile)
	return map[string]any{"digest": g.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}}
}

// signOut401 has the vendor answer a call through g with 401, as Gmail did:
// the connection is marked signed out.
func signOut401(t *testing.T, e *env, f *connectorFixture, g *Grant) {
	t.Helper()
	f.executeBody = `{"data":{"status_code":401,"message":"Invalid Credentials"},"successful":false,"error":"Request had invalid authentication credentials. Expected OAuth 2 access token"}`
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_out"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_out", Args: args})
	got := mustJSONString(e.s.call(context.Background(), g, "reply_ticket", args, permit))
	// Plain words for the taf, with where to sign in again: never the vendor's own body.
	if strings.Contains(got, "Expected OAuth 2") || strings.Contains(got, "invalid authentication") || !strings.Contains(got, "Connect again") {
		t.Fatalf("what the taf reads: %s", got)
	}
	f.executeBody = ""
}

// Sign in again on a signed-out managed connection: one approval, the
// provider's sign-in, the device's confirm, and the same grant is connected
// again, with its tools and its registration in KIFF as they were (so the
// taf's Card stays), with no new Card and no Coach round.
func TestASignedOutManagedConnectionSignsInAgainOnItself(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct = true
	first := connectConnector(t, e, f)
	declared, _ := e.kiff.Connection(tenant, "reply_ticket")
	signOut401(t, e, f, first)
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	if g.Status != StatusNeedsReconnect {
		t.Fatalf("after the 401: %s", g.Status)
	}
	path := "/v1/accounts/" + tenant + "/grants/" + first.ID
	e.at(time.Minute) // a new signature, later, on both clocks
	e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
	st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, "")
	if st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&status=ACTIVE")
	st, out = e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(ctx, tenant, first.ID)
	if g.ID != first.ID || len(g.Tools) != len(first.Tools) || g.Request.Connector.Account != g.Connector.Account || g.Connector.Account == first.Connector.Account {
		t.Fatalf("the grant after reconnecting: tools %d→%d, bound %q, provider account %q (was %q)", len(first.Tools), len(g.Tools), g.Request.Connector.Account, g.Connector.Account, first.Connector.Account)
	}
	if c, ok := e.kiff.Connection(tenant, "reply_ticket"); !ok || c.Credential != declared.Credential {
		t.Fatal("its registration in KIFF changed: the taf's Card would no longer cover it")
	}
	// It works again: a call goes through on the new sign-in.
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_back"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_back", Args: args})
	if got := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); strings.Contains(got, `"isError":true`) {
		t.Fatalf("a call after reconnecting: %s", got)
	}
}

// A reconnect's sign-in to another vendor account connects nothing: the
// tools and Card were approved for the first account, never swapped under it.
// Category: authority
// Attack: Sign in again with another account, so the taf acts on an account the owner never approved it for
// Refused: 409 at the confirm, that sign-in revoked, the connection still signed out
func TestAReconnectToAnotherAccountConnectsNothing(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct = true
	first := connectConnector(t, e, f)
	signOut401(t, e, f, first)
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	path := "/v1/accounts/" + tenant + "/grants/" + first.ID
	e.at(time.Minute) // a new signature, later, on both clocks
	e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
	if st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, ""); st != 200 {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	f.accountName = "Another Shop"
	code := e.signIn(f.callback + "&status=ACTIVE")
	revoked := len(f.revoked)
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 409 || !strings.Contains(mustJSONString(out), "another account") {
		t.Fatalf("confirm on another account: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(ctx, tenant, first.ID)
	if g.Status != StatusNeedsReconnect || g.Request.Connector.Account != first.Connector.Account || len(f.revoked) != revoked+1 {
		t.Fatalf("after another account: %s, bound %q, revoked %d→%d", g.Status, g.Request.Connector.Account, revoked, len(f.revoked))
	}
}

// A reconnect fails closed when the account can't be compared: the provider
// shows no vendor id, only the account's name, so a new sign-in that comes
// back with no name is refused like another account's, and a connection
// whose account never had a name isn't offered a new sign-in at all.
// Category: authority
// Attack: Sign in again with another account the provider shows no name for, so the name check is skipped
// Refused: 409 at the confirm, that sign-in revoked, the connection still signed out; or 409 at approve, nothing opened
func TestAReconnectWhoseAccountCantBeComparedConnectsNothing(t *testing.T) {
	t.Run("the new sign-in has no name", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		ctx := context.Background()
		f.distinct = true
		first := connectConnector(t, e, f)
		signOut401(t, e, f, first)
		g, _ := e.s.loadGrant(ctx, tenant, first.ID)
		path := "/v1/accounts/" + tenant + "/grants/" + first.ID
		e.at(time.Minute)
		e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
		if st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, ""); st != 200 {
			t.Fatalf("sign in again: %d %v", st, out)
		}
		f.noName = true
		code := e.signIn(f.callback + "&status=ACTIVE")
		revoked := len(f.revoked)
		if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 409 || !strings.Contains(mustJSONString(out), "can't be shown") {
			t.Fatalf("confirm with no name: %d %v", st, out)
		}
		g, _ = e.s.loadGrant(ctx, tenant, first.ID)
		if g.Status != StatusNeedsReconnect || g.Request.Connector.Account != first.Connector.Account || len(f.revoked) != revoked+1 {
			t.Fatalf("after no name: %s, bound %q, revoked %d→%d", g.Status, g.Request.Connector.Account, revoked, len(f.revoked))
		}
	})
	t.Run("the connection never had a name", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		ctx := context.Background()
		f.distinct = true
		f.noName = true
		first := connectConnector(t, e, f)
		if first.Connector.AccountName != "" {
			t.Fatalf("named %q", first.Connector.AccountName)
		}
		signOut401(t, e, f, first)
		g, _ := e.s.loadGrant(ctx, tenant, first.ID)
		path := "/v1/accounts/" + tenant + "/grants/" + first.ID
		e.at(time.Minute)
		e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
		f.noName = false
		if st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, ""); st != 409 || out["sign_in_url"] != nil || !strings.Contains(mustJSONString(out), "remove access and connect") {
			t.Fatalf("sign in again with no name kept: %d %v", st, out)
		}
		if g, _ = e.s.loadGrant(ctx, tenant, first.ID); g.Status != StatusNeedsReconnect || g.Reconnect {
			t.Fatalf("after approve: %s, reconnect %v", g.Status, g.Reconnect)
		}
	})
}
