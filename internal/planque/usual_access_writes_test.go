package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// 20 Nov, Google Calendar: a sign-in made for a read was reused for a taf that
// creates events. planque said "no second sign-in", the first write came back
// 403 "insufficient authentication scopes", and the owner signed in again
// anyway. A sign-in on the app's usual access (the provider couldn't narrow it,
// or reported its auth config's scopes, not the vendor's grant) isn't trusted
// for a tool that writes: a new sign-in is asked for up front. A read it
// covers is still reused.
// Category: authority
// Attack: a held usual-access sign-in reused for a tool that writes, then refused by the vendor
// Refused: not reused for the write; the read is
func TestAUsualAccessSignInIsntReusedForAToolThatWrites(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.search, f.searchDescribed, f.granted = true, true, "https://www.googleapis.com/auth/spreadsheets"
	f.replyScopes = []string{"https://www.googleapis.com/auth/spreadsheets"}
	f.searchScopes = []string{"https://www.googleapis.com/auth/spreadsheets.readonly"}
	first := connectConnector(t, e, f)
	if !first.Request.Connector.UsualAccess {
		t.Fatalf("the held sign-in isn't usual access: %+v", first.Request.Connector)
	}
	entry := "/v1/accounts/" + tenant + "/catalog/connector:composio:desk_demo/grants"

	st, out := e.do("POST", entry, map[string]any{"agent": "agt_reply", "tools": []string{"DESK_DEMO_REPLY"}}, nil, "k_setup")
	if v, _ := out["grant"].(map[string]any); st != 201 || v["status"] != StatusPending || v["signed_in"] == true {
		t.Fatalf("a write on a usual-access sign-in: %d %v", st, out)
	}
	st, out = e.do("POST", entry, map[string]any{"agent": "agt_search", "tools": []string{"DESK_DEMO_SEARCH"}}, nil, "k_setup")
	if v, _ := out["grant"].(map[string]any); st != 201 || v["status"] != StatusListed || v["signed_in"] != true {
		t.Fatalf("a read on a usual-access sign-in: %d %v", st, out)
	}
}

// A vendor's 403 for too little access ends the account's reuse for a tool that
// writes, on every connection of that account, even when planque had trusted
// its scopes (a narrowed sign-in): the scopes are unknown from then on.
// Category: authority
// Attack: a sign-in recorded as covering a write, refused by the vendor for too little access, reused for the next write
// Refused: the account is marked, its scopes unknown, and a new sign-in is asked for
func TestAVendorRefusalForTooLittleAccessEndsReuseForWrites(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.granted = "https://www.googleapis.com/auth/spreadsheets"
	f.replyScopes = []string{"https://www.googleapis.com/auth/spreadsheets"}
	g := connectConnector(t, e, f)
	ctx := context.Background()
	// A sign-in planque trusts: narrowed to its tools, not the app's usual access.
	held, _ := e.s.loadGrant(ctx, tenant, g.ID)
	held.Request.Connector.UsualAccess = false
	if err := e.s.saveGrant(ctx, held); err != nil {
		t.Fatal(err)
	}
	entry := "/v1/accounts/" + tenant + "/catalog/connector:composio:desk_demo/grants"
	st, out := e.do("POST", entry, map[string]any{"agent": "agt_first", "tools": []string{"DESK_DEMO_REPLY"}}, nil, "k_setup")
	if v, _ := out["grant"].(map[string]any); st != 201 || v["status"] != StatusListed || v["signed_in"] != true {
		t.Fatalf("the trusted sign-in isn't reused before the refusal: %d %v", st, out)
	}

	f.executeBody = `{"data":{"http_error":"Forbidden","status_code":403},"successful":false,"error":"Request had insufficient authentication scopes."}`
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_scopes"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: held.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_scopes", Args: args})
	if got := mustJSONString(e.s.call(ctx, held, "reply_ticket", args, permit)); !strings.Contains(got, "insufficient authentication scopes") || !strings.Contains(got, `"isError":true`) {
		t.Fatalf("the refused call: %s", got)
	}
	after, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if !after.Connector.ScopesRefused || grantedScopes(after) != nil {
		t.Fatalf("the account isn't marked: refused %v, scopes %v", after.Connector.ScopesRefused, grantedScopes(after))
	}
	st, out = e.do("POST", entry, map[string]any{"agent": "agt_second", "tools": []string{"DESK_DEMO_REPLY"}}, nil, "k_setup")
	if v, _ := out["grant"].(map[string]any); st != 201 || v["status"] != StatusPending || v["signed_in"] == true {
		t.Fatalf("a write on the refused account: %d %v", st, out)
	}
}

// Only a 403 about too little access marks the account: another refusal, or a
// 403 for something else, leaves reuse as it was.
// Category: authority
// Attack: an unrelated vendor refusal ending the reuse of a good sign-in
// Refused: nothing is marked
func TestOnlyAScopesRefusalMarksTheAccount(t *testing.T) {
	for _, reason := range []string{"Forbidden", "No such ticket: 'ticket_demo'", "Rate limit exceeded"} {
		if insufficientScopes(reason) {
			t.Errorf("%q read as too little access", reason)
		}
	}
	for _, reason := range []string{"Request had insufficient authentication scopes.", "Insufficient Permission", "insufficientPermissions"} {
		if !insufficientScopes(reason) {
			t.Errorf("%q not read as too little access", reason)
		}
	}
}
