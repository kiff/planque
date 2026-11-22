package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/store"
)

// heldMoneybird is a connected Moneybird sign-in of tenant's on account.
func heldMoneybird(id, account, created string) *Grant {
	return &Grant{ID: id, Tenant: tenant, Status: StatusConnected, SignedIn: true,
		Request:   Request{CatalogEntry: "moneybird@1", Tenant: tenant, CreatedAt: created, Connector: &ConnectorRequest{Provider: "composio", Toolkit: "moneybird"}},
		Connector: &ConnectorState{Account: account}}
}

// A tenant holds sign-ins to two accounts at the same app, one of them an
// account no taf is meant for: planque never picks one for a new taf, the
// latest or any other. Only a grant the taf names is reused, and a sign-in
// marked never to be used isn't, by any grant on its account.
// Category: authority
// Attack: a taf reusing the latest sign-in, an account it isn't meant for
// Refused: no reuse between two accounts unnamed; a never-use account isn't reused, named or not
func TestTwoSignInsAtAnAppAreNeverPickedBetween(t *testing.T) {
	ctx := t.Context()
	s := &Server{KV: &store.Memory{}}
	test := heldMoneybird("grt_testbooks0000001", "ca_test", "2026-11-18T10:00:00Z")
	real := heldMoneybird("grt_realbooks0000001", "ca_real", "2026-11-19T10:00:00Z")
	for _, g := range []*Grant{test, real} {
		if err := s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	if h, err := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, ""); err != nil || h != nil {
		t.Fatalf("two accounts, none named: reused %v (%v)", h, err)
	}
	for _, named := range []*Grant{test, real} {
		if h, err := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, named.ID); err != nil || h == nil || h.ID != named.ID {
			t.Fatalf("named %s: reused %v (%v)", named.ID, h, err)
		}
	}
	if h, _ := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, "grt_nosuchgrant00001"); h != nil {
		t.Fatalf("a grant named that isn't held: reused %v", h)
	}

	// The other account marked never to be used: the test sign-in is the only
	// candidate, and is reused.
	s.NeverUse = map[string]bool{real.ID: true}
	if h, err := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, ""); err != nil || h == nil || h.ID != test.ID {
		t.Fatalf("never-use marked: reused %v (%v)", h, err)
	}
	if h, _ := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, real.ID); h != nil {
		t.Fatalf("a never-use grant named: reused %v", h)
	}
	// Another grant on that account is refused too: the mark is
	// the account's, not only its grant's.
	again := heldMoneybird("grt_realbooks0000002", "ca_real", "2026-11-20T10:00:00Z")
	if err := s.saveGrant(ctx, again); err != nil {
		t.Fatal(err)
	}
	if h, _ := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, again.ID); h != nil {
		t.Fatalf("another grant on a never-use account, named: reused %v", h)
	}
	if h, err := s.heldSignIn(ctx, tenant, "composio", "moneybird", nil, ""); err != nil || h == nil || h.ID != test.ID {
		t.Fatalf("another grant on a never-use account counted: reused %v (%v)", h, err)
	}
}

// A sign-in marked never to be used serves no call, through its own grant
// or another on the same account, though each is connected and permitted.
// Category: authority
// Attack: a call through a connection to an account no taf may use
// Refused: not sent, and the vendor is never reached
func TestASignInMarkedNeverToBeUsedServesNoCall(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"operation_never"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "operation_never", Args: args})
	e.s.NeverUse = map[string]bool{g.ID: true}
	res := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if f.executions != 0 || !strings.Contains(res, "marked never to be used") {
		t.Fatalf("a never-use grant's call: %s executions=%d", res, f.executions)
	}
	// Marked by another grant on the same account: refused the same.
	other := heldMoneybird("grt_sameaccount00001", g.Connector.Account, "2026-11-18T10:00:00Z")
	if err := e.s.saveGrant(ctx, other); err != nil {
		t.Fatal(err)
	}
	e.s.NeverUse = map[string]bool{other.ID: true}
	args = json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"operation_never2"}`)
	permit = e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "operation_never2", Args: args})
	res = mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if f.executions != 0 || !strings.Contains(res, "marked never to be used") {
		t.Fatalf("a call on a never-use account: %s executions=%d", res, f.executions)
	}
	// Unmarked, the same call is sent.
	e.s.NeverUse = nil
	args = json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"operation_never3"}`)
	permit = e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "operation_never3", Args: args})
	if res := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); f.executions != 1 {
		t.Fatalf("unmarked: %s executions=%d", res, f.executions)
	}
}

func TestNeverUseFrom(t *testing.T) {
	got, err := NeverUseFrom("grants:grt_realbooks0000001, grt_realbooks0000002")
	if err != nil || len(got) != 2 || !got["grt_realbooks0000002"] {
		t.Fatalf("%v %v", got, err)
	}
	if got, err := NeverUseFrom(""); err != nil || len(got) != 0 {
		t.Fatalf("empty: %v %v", got, err)
	}
	for _, bad := range []string{"grt_realbooks0000001", "grants:", "grants:ca_real", "grants:grt_short"} {
		if _, err := NeverUseFrom(bad); err == nil {
			t.Fatalf("%q was taken", bad)
		}
	}
}
