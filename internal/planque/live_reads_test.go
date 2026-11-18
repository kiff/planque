package planque

import (
	"context"
	"fmt"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// countingKV counts the items every Query returns: what a read of the table
// would be billed for, whatever the caller then keeps.
type countingKV struct {
	store.KV
	read int
}

func (c *countingKV) Query(ctx context.Context, pk, prefix string) ([]store.Item, error) {
	items, err := c.KV.Query(ctx, pk, prefix)
	c.read += len(items)
	return items, err
}

// finishedGrants adds n grants that no longer count (disconnected, failed),
// the way an account that has run many connections for a while holds them.
func finishedGrants(t *testing.T, s *Server, n int) {
	t.Helper()
	for i := range n {
		status := StatusDisconnected
		if i%2 == 1 {
			status = StatusFailed
		}
		g := &Grant{ID: fmt.Sprintf("grt_done_%d", i), Tenant: tenant, Status: status, SignedIn: true,
			Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant}}
		if err := s.saveGrant(t.Context(), g); err != nil {
			t.Fatal(err)
		}
	}
}

// What the catalog and a held sign-in read must stay what they read today
// (a connected grant counts, a revoked, failed or expired one never does),
// and what they read must not grow with the grants an account has finished
// with: a store that counts items read sees the same number at 5 finished
// grants and at 200 (the Eval tenant held 960 and every catalog read took
// all of them).
func TestTheCatalogAndAHeldSignInReadOnlyLiveGrants(t *testing.T) {
	read := func(finished int) (reads int, states map[string]catalogEntryView, held *Grant) {
		kv := &countingKV{KV: &store.Memory{}}
		s := &Server{KV: kv}
		live := &Grant{ID: "grt_live", Tenant: tenant, Status: StatusConnected, SignedIn: true,
			Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant,
				Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
			Connector: &ConnectorState{Account: "acct_1"}}
		reconnect := &Grant{ID: "grt_reconnect", Tenant: tenant, Status: StatusNeedsReconnect, SignedIn: true,
			Request: Request{CatalogEntry: "typeform-key@1", Tenant: tenant}}
		for _, g := range []*Grant{live, reconnect} {
			if err := s.saveGrant(t.Context(), g); err != nil {
				t.Fatal(err)
			}
		}
		finishedGrants(t, s, finished)
		kv.read = 0
		st, err := s.entryStates(t.Context(), tenant)
		if err != nil {
			t.Fatal(err)
		}
		h, err := s.heldSignIn(t.Context(), tenant, "composio", "gmail", nil)
		if err != nil {
			t.Fatal(err)
		}
		return kv.read, st, h
	}
	smallReads, smallStates, smallHeld := read(5)
	bigReads, bigStates, bigHeld := read(200)

	for name, states := range map[string]map[string]catalogEntryView{"5 finished": smallStates, "200 finished": bigStates} {
		if states["gmail-key@1"].State != "connected" || states["typeform-key@1"].State != "needs_reconnect" || len(states) != 2 {
			t.Errorf("%s: the catalog's states changed: %+v", name, states)
		}
	}
	for name, h := range map[string]*Grant{"5 finished": smallHeld, "200 finished": bigHeld} {
		if h == nil || h.ID != "grt_live" {
			t.Errorf("%s: the held sign-in must be the connected grant, got %v", name, h)
		}
	}
	if bigReads != smallReads {
		t.Errorf("items read grew with finished grants: %d at 5, %d at 200", smallReads, bigReads)
	}
}
