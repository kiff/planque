package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/store"
)

// countingKV counts the items every Query returns: what a read of the table
// would be billed for, whatever the caller then keeps.
type countingKV struct {
	store.KV
	read int
	puts int
}

func (c *countingKV) Put(ctx context.Context, it store.Item, cond store.Cond) error {
	c.puts++
	return c.KV.Put(ctx, it, cond)
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
		// Made to fool a reader that took them: the same app and a sign-in
		// held, ids sorting after the live grant's, some on an entry no live
		// grant has.
		entry := "gmail-key@1"
		if i%3 == 0 {
			entry = "gone-key@1"
		}
		g := &Grant{ID: fmt.Sprintf("grt_zz_done_%d", i), Tenant: tenant, Status: status, SignedIn: true,
			Request: Request{CatalogEntry: entry, Tenant: tenant,
				Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
			Connector: &ConnectorState{Account: "acct_done"}}
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
		// "Live" is every status but finished: a connected grant asking again
		// (pending with tools connected) and a listed one with triggers read
		// connected; a plain pending one reads as nothing.
		reapproval := &Grant{ID: "grt_reapproval", Tenant: tenant, Status: StatusPending, SignedIn: true,
			Connected: []ToolRequest{{Tool: "t", Name: "t"}},
			Request:   Request{CatalogEntry: "third-key@1", Tenant: tenant}}
		triggers := &Grant{ID: "grt_triggers", Tenant: tenant, Status: StatusListed, SignedIn: true,
			TriggerIDs: []string{"trg_1"}, Request: Request{CatalogEntry: "fourth-key@1", Tenant: tenant}}
		plain := &Grant{ID: "grt_plain", Tenant: tenant, Status: StatusPending,
			Request: Request{CatalogEntry: "fifth-key@1", Tenant: tenant}}
		for _, g := range []*Grant{live, reconnect, reapproval, triggers, plain} {
			if err := s.saveGrant(t.Context(), g); err != nil {
				t.Fatal(err)
			}
		}
		finishedGrants(t, s, finished)
		// Steady state: the account's older grants are indexed once, then read.
		if err := s.ensureLive(t.Context(), tenant); err != nil {
			t.Fatal(err)
		}
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
		if states["gmail-key@1"].State != "connected" || states["typeform-key@1"].State != "needs_reconnect" ||
			states["third-key@1"].State != "connected" || states["fourth-key@1"].State != "connected" || len(states) != 4 {
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

// A grant that becomes finished leaves the catalog and the held sign-in at
// once, and one that becomes connected again (Reconnect) is back in them: the
// moves go through saveGrant, the one place a grant is written.
func TestAGrantMovingToFinishedOrBackIsReadAtOnce(t *testing.T) {
	s := &Server{KV: &store.Memory{}}
	g := &Grant{ID: "grt_live", Tenant: tenant, Status: StatusConnected, SignedIn: true,
		Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant,
			Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
		Connector: &ConnectorState{Account: "acct_1"}}
	state := func() (string, *Grant) {
		st, err := s.entryStates(t.Context(), tenant)
		if err != nil {
			t.Fatal(err)
		}
		h, err := s.heldSignIn(t.Context(), tenant, "composio", "gmail", nil)
		if err != nil {
			t.Fatal(err)
		}
		return st["gmail-key@1"].State, h
	}
	move := func(status string) {
		g.Status = status
		if err := s.saveGrant(t.Context(), g); err != nil {
			t.Fatal(err)
		}
	}
	move(StatusConnected)
	if st, h := state(); st != "connected" || h == nil {
		t.Fatalf("connected: %q %v", st, h)
	}
	move(StatusDisconnected)
	if st, h := state(); st != "" || h != nil {
		t.Fatalf("disconnected must leave both at once: %q %v", st, h)
	}
	move(StatusConnected)
	if st, h := state(); st != "connected" || h == nil || h.ID != "grt_live" {
		t.Fatalf("connected again must be back: %q %v", st, h)
	}
	move(StatusFailed)
	if st, h := state(); st != "" || h != nil {
		t.Fatalf("failed must leave both at once: %q %v", st, h)
	}
}

// The two real paths in: Reconnect (needs_reconnect back to connected) is read
// as connected at once, and a pending re-approval that ends leaves the catalog
// at once.
func TestReconnectAndAnEndedReapprovalAreReadAtOnce(t *testing.T) {
	s := &Server{KV: &store.Memory{}}
	g := &Grant{ID: "grt_1", Tenant: tenant, Status: StatusNeedsReconnect, SignedIn: true,
		Request: Request{CatalogEntry: "typeform-key@1", Tenant: tenant}}
	stateOf := func(entry string) string {
		st, err := s.entryStates(t.Context(), tenant)
		if err != nil {
			t.Fatal(err)
		}
		return st[entry].State
	}
	save := func() {
		if err := s.saveGrant(t.Context(), g); err != nil {
			t.Fatal(err)
		}
	}
	save()
	if got := stateOf("typeform-key@1"); got != "needs_reconnect" {
		t.Fatalf("needs_reconnect: %q", got)
	}
	g.Status = StatusConnected
	save()
	if got := stateOf("typeform-key@1"); got != "connected" {
		t.Fatalf("Reconnect must read connected at once: %q", got)
	}
	r := &Grant{ID: "grt_2", Tenant: tenant, Status: StatusPending, SignedIn: true,
		Connected: []ToolRequest{{Tool: "t", Name: "t"}}, Request: Request{CatalogEntry: "third-key@1", Tenant: tenant}}
	if err := s.saveGrant(t.Context(), r); err != nil {
		t.Fatal(err)
	}
	if got := stateOf("third-key@1"); got != "connected" {
		t.Fatalf("a re-approval reads connected: %q", got)
	}
	r.Status = StatusDisconnected
	if err := s.saveGrant(t.Context(), r); err != nil {
		t.Fatal(err)
	}
	if got := stateOf("third-key@1"); got != "" {
		t.Fatalf("an ended re-approval must leave the catalog at once: %q", got)
	}
}

// A grant written the old way (no live row, straight to the table) is found
// once the account is indexed; a second pass writes nothing; a finished grant
// gets no live row; and a row a newer save wrote is never overwritten.
func TestGrantsMadeBeforeLiveRowsAreIndexedOnce(t *testing.T) {
	kv := &countingKV{KV: &store.Memory{}}
	s := &Server{KV: kv}
	old := func(id, status, entry string) *Grant {
		g := &Grant{ID: id, Tenant: tenant, Status: status, SignedIn: true, Request: Request{CatalogEntry: entry, Tenant: tenant}}
		if err := putDoc(t.Context(), kv, store.Item{PK: tenantPK(tenant), SK: skGrant(id), Ver: 1}, g, store.Cond{NotExists: true}); err != nil {
			t.Fatal(err)
		}
		return g
	}
	old("grt_old_live", StatusConnected, "gmail-key@1")
	old("grt_old_done", StatusDisconnected, "typeform-key@1")
	old("grt_old_failed", StatusFailed, "third-key@1")
	newer := old("grt_newer", StatusConnected, "fourth-key@1")
	// A save in between wrote the grant's newer state and its live row (the
	// grant first, as saveGrant does), before the backfill reached it.
	newer.Status = StatusNeedsReconnect
	for _, sk := range []string{skGrant(newer.ID), skLive(newer.ID)} {
		if err := putDoc(t.Context(), kv, store.Item{PK: tenantPK(tenant), SK: sk, Ver: 2}, newer, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}

	st, err := s.entryStates(t.Context(), tenant)
	if err != nil {
		t.Fatal(err)
	}
	if st["gmail-key@1"].State != "connected" || st["fourth-key@1"].State != "needs_reconnect" || len(st) != 2 {
		t.Fatalf("after the backfill: %+v", st)
	}
	live, err := kv.KV.Query(t.Context(), tenantPK(tenant), "LIVE#")
	if err != nil {
		t.Fatal(err)
	}
	if len(live) != 2 {
		t.Fatalf("live rows: %d, want 2 (a finished grant gets none)", len(live))
	}
	before := kv.puts
	if err := s.ensureLive(t.Context(), tenant); err != nil {
		t.Fatal(err)
	}
	if kv.puts != before {
		t.Fatalf("a second pass wrote %d items", kv.puts-before)
	}
}

// A live row can outlive its grant: a refresh saved after a revoke, a write
// that failed halfway, the backfill racing a save. Each leaves a live row that
// says connected while the grant itself is finished, gone, or newer. None may
// read as held or as connected, and the stale row is gone or mended after.
func TestAStaleLiveRowIsNeverRead(t *testing.T) {
	kv := &countingKV{KV: &store.Memory{}}
	s := &Server{KV: kv}
	grant := func(id, status string) *Grant {
		return &Grant{ID: id, Tenant: tenant, Status: status, SignedIn: true,
			Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant,
				Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
			Connector: &ConnectorState{Account: "acct_1"}}
	}
	liveRow := func(g *Grant, ver int64) {
		if err := putDoc(t.Context(), kv, store.Item{PK: tenantPK(tenant), SK: skLive(g.ID), Ver: ver}, g, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	rowKept := func(id string) bool {
		_, ok, err := kv.Get(t.Context(), tenantPK(tenant), skLive(id))
		if err != nil {
			t.Fatal(err)
		}
		return ok
	}
	reads := func() (string, *Grant) {
		st, err := s.entryStates(t.Context(), tenant)
		if err != nil {
			t.Fatal(err)
		}
		h, err := s.heldSignIn(t.Context(), tenant, "composio", "gmail", nil)
		if err != nil {
			t.Fatal(err)
		}
		return st["gmail-key@1"].State, h
	}

	// 1. The grant was disconnected, and a stale connected copy came back.
	g := grant("grt_revoked", StatusConnected)
	if err := s.saveGrant(t.Context(), g); err != nil {
		t.Fatal(err)
	}
	g.Status = StatusDisconnected
	if err := s.saveGrant(t.Context(), g); err != nil {
		t.Fatal(err)
	}
	stale := grant("grt_revoked", StatusConnected)
	liveRow(stale, 1)
	if st, h := reads(); st != "" || h != nil {
		t.Fatalf("a revoked grant's stale row was read: %q %v", st, h)
	}
	if rowKept("grt_revoked") {
		t.Fatal("the stale row was left")
	}

	// 2. The grant's own record is gone.
	liveRow(grant("grt_gone", StatusConnected), 3)
	if st, h := reads(); st != "" || h != nil {
		t.Fatalf("a row with no grant was read: %q %v", st, h)
	}
	if rowKept("grt_gone") {
		t.Fatal("the orphan row was left")
	}

	// 3. The row is older than a grant that still counts: it reads as the
	// grant says, and the row is mended.
	cur := grant("grt_cur", StatusConnected)
	if err := s.saveGrant(t.Context(), cur); err != nil {
		t.Fatal(err)
	}
	old := grant("grt_cur", StatusNeedsReconnect)
	liveRow(old, 0)
	if st, h := reads(); st != "connected" || h == nil || h.ID != "grt_cur" {
		t.Fatalf("an older row must read as its grant: %q %v", st, h)
	}
	row, _, _ := kv.Get(t.Context(), tenantPK(tenant), skLive("grt_cur"))
	var mended Grant
	if err := json.Unmarshal(row.Doc, &mended); err != nil || mended.Status != StatusConnected {
		t.Fatalf("the row wasn't mended: %v %q", err, mended.Status)
	}
}

// The tool names a new connection can't reuse read the live grants only:
// the same answer, and the same number of items read, at 5 finished grants
// and at 200 (two test accounts held 341 and 151).
func TestNamesInUseReadsOnlyLiveGrants(t *testing.T) {
	read := func(finished int) (reads int, names map[string]string) {
		kv := &countingKV{KV: &store.Memory{}}
		s := &Server{KV: kv}
		connected := &Grant{ID: "grt_live", Tenant: tenant, Status: StatusConnected, SignedIn: true,
			Tools: []GrantedTool{{ToolRequest: ToolRequest{Name: "gmail_send"}}},
			Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant, Server: "https://mcp.gmail.example/",
				Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
			Connector: &ConnectorState{Account: "acct_1"}}
		sibling := &Grant{ID: "grt_sibling", Tenant: tenant, Status: StatusConnected, SignedIn: true,
			Request: Request{CatalogEntry: "gmail-key@1", Tenant: tenant,
				Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}},
			Connector: &ConnectorState{Account: "acct_1"}}
		for _, g := range []*Grant{connected, sibling} {
			if err := s.saveGrant(t.Context(), g); err != nil {
				t.Fatal(err)
			}
		}
		finishedGrants(t, s, finished)
		if err := s.ensureLive(t.Context(), tenant); err != nil {
			t.Fatal(err)
		}
		kv.read = 0
		n, err := s.namesInUse(t.Context(), tenant, "")
		if err != nil {
			t.Fatal(err)
		}
		return kv.read, n
	}
	smallReads, smallNames := read(5)
	bigReads, bigNames := read(200)
	for name, got := range map[string]map[string]string{"5 finished": smallNames, "200 finished": bigNames} {
		if len(got) != 1 || got["gmail_send"] == "" {
			t.Errorf("%s: names in use changed: %v", name, got)
		}
	}
	if bigReads != smallReads {
		t.Errorf("items read grew with finished grants: %d at 5, %d at 200", smallReads, bigReads)
	}
}

// tafs cloud's read of an account's connections and Setup's read of what it
// asked for, the two routes read on most requests, read the live grants
// only: the same answers, and the same number of items read, at 5 finished
// grants and at 200.
func TestTheAccountsConnectionsAndAskedGrantsReadOnlyLiveGrants(t *testing.T) {
	read := func(finished int) (reads int, conns, asked int) {
		e := newEnv(t)
		e.kiff.Keys["k_person"] = kiff.Me{TenantID: tenant}
		e.connect()
		finishedGrants(t, e.s, finished)
		if err := e.s.ensureLive(t.Context(), tenant); err != nil {
			t.Fatal(err)
		}
		kv := &countingKV{KV: e.s.KV}
		e.s.KV = kv
		st, cs := e.connections("k_person")
		if st != 200 {
			t.Fatalf("connections: %d", st)
		}
		return kv.read, len(cs), len(e.asked())
	}
	smallReads, smallConns, smallAsked := read(5)
	bigReads, bigConns, bigAsked := read(200)
	if smallConns != 1 || bigConns != 1 || smallAsked != 1 || bigAsked != 1 {
		t.Errorf("the connected grant alone is listed: connections %d/%d, asked %d/%d", smallConns, bigConns, smallAsked, bigAsked)
	}
	if bigReads != smallReads {
		t.Errorf("items read grew with finished grants: %d at 5, %d at 200", smallReads, bigReads)
	}
}
