package planque

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
)

// The owner's retest (kiff/tafs#138, 22 Oct): improving a taf took five Face
// IDs and five Stripe sign-ins, because a connected grant couldn't gain a
// tool and each new one needed a grant of its own. A connected grant now
// gains tools: one approval, no sign-in, and its tools work meanwhile.

// connectedRefunds is a grant, connected first, whose one tool is create_refund.
func (e *env) connectedRefunds() string {
	e.t.Helper()
	id := e.listedGrant()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	refund := map[string]any{"tool": "create_refund", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("set tools: %d %v", st, out)
	}
	if st, out := e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	return id
}

func (e *env) refund() string {
	e.t.Helper()
	res, err := e.kiff.Call(context.Background(), "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 300})
	if err != nil {
		return err.Error()
	}
	return res.Decision
}

func TestAConnectedGrantGainsTools(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	list := map[string]any{"tool": "list_charges", "name": "list_charges", "read_only": true}

	// The asking agent alone adds, and only what the grant doesn't have.
	if st, _ := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_taf"); st != 403 {
		t.Fatalf("another agent adding: %d", st)
	}
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "create_refund", "name": "create_refund"}}}, nil, "k_setup"); st != 409 ||
		!strings.Contains(out["error"].(string), "send only the tools to add") {
		t.Fatalf("adding a tool it has: %d %v", st, out)
	}
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending || out["signed_in"] != true || out["held_until"] == nil ||
		out["status_words"] != "Connected · new tools waiting for your approval" {
		t.Fatalf("add: %d %v", st, out)
	}
	text := out["text"].(string)
	if !strings.HasPrefix(text, "This adds to a connection the taf already has: ") || !strings.Contains(text, "keeps working while you decide") ||
		!strings.Contains(text, "already signed in") || strings.Contains(text, "You will sign in") || !strings.Contains(text, "Refund a charge") {
		t.Fatalf("the owner's text:\n%s", text)
	}
	// Meanwhile the taf's refund still works, and the new tool isn't in KIFF.
	if d := e.refund(); d != "allowed" {
		t.Fatalf("a refund while the addition waits: %s", d)
	}
	if _, ok := e.kiff.Connection(tenant, "list_charges"); ok {
		t.Fatal("list_charges is in KIFF before the owner approved it")
	}

	// One approval connects it: no sign-in link, both tools in KIFF, and the
	// refund still works on its first credential.
	signIns := e.vendor.Live
	a0, _ := signIns()
	st, out = e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] != nil {
		t.Fatalf("approve the addition: %d %v", st, out)
	}
	if a, _ := signIns(); a != a0 {
		t.Fatalf("approving the addition signed in again: %d live tokens, were %d", a, a0)
	}
	if _, ok := e.kiff.Connection(tenant, "list_charges"); !ok {
		t.Fatal("list_charges isn't in KIFF")
	}
	if tools, _ := out["tools"].([]any); len(tools) != 2 {
		t.Fatalf("tools: %v", out["tools"])
	}
	if d := e.refund(); d != "allowed" {
		t.Fatalf("a refund after the addition: %s", d)
	}
	res, err := e.kiff.Call(context.Background(), "k_taf", "list_charges", map[string]any{})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("the added tool through KIFF: %+v %v", res, err)
	}
}

// An addition nobody approves goes after a day; the grant stays connected
// with what it had, and its sign-in stays.
func TestAnAdditionNobodyApprovesGoesBack(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 200 {
		t.Fatalf("add: %d %v", st, out)
	}
	e.s.Now = func() time.Time { return time.Now().Add(holdTTL + time.Minute) }
	st, out := e.do("GET", acct, nil, nil, "k_setup")
	if st != 200 || out["status"] != StatusConnected || out["signed_in"] != true {
		t.Fatalf("after a day: %d %v", st, out)
	}
	if tools, _ := out["tools"].([]any); len(tools) != 1 {
		t.Fatalf("tools after a day: %v", out["tools"])
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the sign-in was revoked with the addition")
	}
	if d := e.refund(); d != "allowed" {
		t.Fatalf("a refund after the addition went: %s", d)
	}
}

// When KIFF won't register an added tool, the addition is undone and the
// grant keeps working with what it had.
func TestAFailedAdditionKeepsTheTools(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("add: %d %v", st, out)
	}
	// Taken in KIFF after it was named, so its registration is refused.
	e.kiff.Add(fakes.Connection{Connection: kiff.Connection{Name: "list_charges", Mode: "sealed"}, Tenant: tenant})
	if st, out := e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 502 {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusConnected || len(g.Tools) != 1 || len(g.Request.Tools) != 1 || g.Connected != nil || len(g.Transports) != 0 {
		t.Fatalf("after a failed addition: %s, tools %d, asked %d, connected %v, transports %v", g.Status, len(g.Tools), len(g.Request.Tools), g.Connected, g.Transports)
	}
	// The record keeps the owner's approval of what is connected: its body
	// names the grant's digest, not the addition's (the owner's export).
	if g.Approval == nil || !strings.Contains(g.Approval.Body, g.Digest) {
		t.Fatalf("the approval after a failed addition is of another digest: %+v, grant digest %s", g.Approval, g.Digest)
	}
	if d := e.refund(); d != "allowed" {
		t.Fatalf("a refund after a failed addition: %s", d)
	}
}

// Setup, improving a taf, finds the connections it asked for that taf.
func TestSetupFindsTheTafsConnections(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/agents/agt_taf/grants", nil, nil, "k_setup")
	gs, _ := out["grants"].([]any)
	if st != 200 || len(gs) != 1 {
		t.Fatalf("Setup's list: %d %v", st, out)
	}
	g := gs[0].(map[string]any)
	if g["id"] != id || g["status"] != StatusConnected || g["server"] != e.vendor.MCPURL() || len(g["tools"].([]any)) != 1 || g["text"] != nil || g["approval"] != nil {
		t.Fatalf("the grant as Setup reads it: %v", g)
	}
	// Another agent sees none of them, and another account's key is refused.
	if _, out := e.do("GET", "/v1/accounts/"+tenant+"/agents/agt_taf/grants", nil, nil, "k_taf"); len(out["grants"].([]any)) != 0 {
		t.Fatalf("another agent's list: %v", out)
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/agents/agt_taf/grants", nil, nil, "k_other"); st != 403 {
		t.Fatalf("another account's key: %d", st)
	}
}

// Setup that lost a taf's agent id finds the connections it asked for, each
// with its agent, and none another agent asked for (kiff/planque#49).
func TestSetupFindsItsConnectionsWithoutTheTafsAgent(t *testing.T) {
	e := newEnv(t)
	id := e.connectedRefunds()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_setup")
	gs, _ := out["grants"].([]any)
	if st != 200 || len(gs) != 1 {
		t.Fatalf("Setup's list: %d %v", st, out)
	}
	g := gs[0].(map[string]any)
	if g["id"] != id || g["agent"] != "agt_taf" || g["status"] != StatusConnected || g["created"] == "" || g["text"] != nil || g["approval"] != nil {
		t.Fatalf("the grant as Setup reads it: %v", g)
	}
	if _, out := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_taf"); len(out["grants"].([]any)) != 0 {
		t.Fatalf("another agent's list: %v", out)
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_other"); st != 403 {
		t.Fatalf("another account's key: %d", st)
	}
}

// A connection that holds only Setup's own reads gives the taf nothing yet: its first
// tools are never called an addition to what it has (kiff/tafs#236, item 4).
func TestTheTafsFirstToolsAfterSetupsReadsArentAnAddition(t *testing.T) {
	added := []ToolRequest{{Tool: "create_refund", Name: "create_refund_x", Title: "Refund a payment"}}
	reads := []ToolRequest{{Tool: "list_accounts", Name: "list_accounts_x", ReadOnly: true, SetupReads: true}}
	if got := additionText(added, reads, nil); !strings.HasPrefix(got, "These are the taf's first tools on this connection: Refund a payment.") ||
		strings.Contains(got, "already has") {
		t.Fatalf("after Setup's reads:\n%s", got)
	}
	taf := append(reads, ToolRequest{Tool: "list_charges", Name: "list_charges_x", ReadOnly: true})
	if got := additionText(added, taf, nil); !strings.HasPrefix(got, "This adds to a connection the taf already has: Refund a payment.") {
		t.Fatalf("after a tool of the taf's own:\n%s", got)
	}
}
