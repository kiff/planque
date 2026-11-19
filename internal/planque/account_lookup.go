package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/store"
)

// A managed sign-in the provider has lost (deleted on its side) is learned
// when something looks it up, and every call does. A connection nobody calls
// (its taf's schedule paused after repeated failures) is never looked
// up, so it kept reading "Connected" after the provider lost it. The owner's
// read of the account's grants looks each connected managed account up once,
// and a gone one signs in again, as it would at its next call.
const (
	// accountLookupEvery is how long a provider account is left alone after
	// a lookup: one lookup per account per hour, however often the owner's
	// app reads.
	accountLookupEvery = time.Hour
	// accountLookupMax is the most accounts one read looks up, so a read of
	// an account with many connections never waits on all of them; the rest
	// are looked up by the next reads (the claims of these keep them out).
	accountLookupMax = 8
)

// accountLookupWithin is how long a read waits for its lookups. A lookup that
// isn't answered by then, or fails, changes nothing. A variable for the test.
var accountLookupWithin = 3 * time.Second

func skAccountLookup(account string) string { return "LOOKUP#" + account }

// lookUpAccounts is the owner's read of tenant's grants checking their
// provider accounts, read only: it asks the provider for each connected
// managed account it hasn't asked about in the last hour, and only an answer
// of "this account is gone or isn't active" (connector.ErrAccount, as at a
// call) signs the connections on it out. An answer that fails, times out or
// isn't about the account leaves every grant as it was. Best effort, never
// an error for the read.
func (s *Server) lookUpAccounts(ctx context.Context, tenant string) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return
	}
	var due []*Grant
	taken := map[string]bool{}
	for _, it := range items {
		if len(due) >= accountLookupMax {
			break
		}
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.Status != StatusConnected || g.Connector == nil || g.Connector.Account == "" ||
			g.Request.Connector == nil || g.Reconnect || g.renewing() || g.Connector.Revocation != nil && g.Connector.Revocation.Disabled {
			continue
		}
		if _, err := s.provider(&g); err != nil || taken[g.Connector.Account] {
			continue
		}
		taken[g.Connector.Account] = true
		// The claim: this account isn't looked up again for an hour, by this
		// planque or another's. A claim that can't be made is a lookup not made.
		claim := store.Item{PK: tenantPK(tenant), SK: skAccountLookup(g.Connector.Account), Ver: 1, Doc: []byte("{}"), TTL: s.now().Add(accountLookupEvery).Unix()}
		if err := s.KV.Put(ctx, claim, store.Cond{NotExists: true}); err != nil {
			if !errors.Is(err, store.ErrCondition) {
				s.log().WarnContext(ctx, "account lookup: the claim wasn't made", slog.String("grant", g.ID), slog.String("err", err.Error()))
			}
			continue
		}
		due = append(due, &g)
	}
	if len(due) == 0 {
		return
	}
	lctx, cancel := context.WithTimeout(ctx, accountLookupWithin)
	defer cancel()
	gone := make([]bool, len(due))
	var wg sync.WaitGroup
	for i, g := range due {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gone[i] = errors.Is(s.checkConnector(lctx, g), connector.ErrAccount)
		}()
	}
	wg.Wait()
	for i, g := range due {
		if gone[i] {
			s.log().InfoContext(ctx, "account lookup: the provider no longer has the account", slog.String("grant", g.ID))
			s.signInRefused(ctx, g, vendorName(g)+" no longer accepts this sign-in")
		}
	}
}
