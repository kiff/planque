package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
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
