package planque

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/store"
)

// pkKept holds one item per grant that ended while its sign-in was kept for
// another grant (a token it shared, or a provider's connected account), so
// the sweep revokes that sign-in once no live grant uses it. Two grants on
// one sign-in disconnected at once each see the other as live and keep it;
// without this, nothing would revoke it afterwards.
const pkKept = "KEPT"

func keptSK(tenant, id string) string { return tenant + "#" + id }

// markKept indexes g, which ended keeping its sign-in for another grant.
// The entry is harmless if g's own save then fails: the sweep acts only on a
// grant that has ended.
func (s *Server) markKept(ctx context.Context, g *Grant) error {
	it := store.Item{PK: pkKept, SK: keptSK(g.Tenant, g.ID), Doc: []byte("{}")}
	if err := s.KV.Put(ctx, it, store.Cond{}); err != nil {
		return fmt.Errorf("index the sign-in kept for another connection: %w", err)
	}
	return nil
}

// SweepKept revokes each kept sign-in no live grant uses any more, at the
// vendor or the connection provider, then drops it. A sign-in a live grant
// uses is never touched, and its entry waits for the next sweep. It returns
// how many sign-ins it revoked, and the first error; a failure keeps the
// entry for the next sweep.
func (s *Server) SweepKept(ctx context.Context) (int, error) {
	items, err := s.KV.Query(ctx, pkKept, "")
	if err != nil {
		return 0, err
	}
	revoked := 0
	var first error
	for _, it := range items {
		tenant, id, ok := strings.Cut(it.SK, "#")
		if ok {
			done, did, err := s.settleKept(ctx, tenant, id)
			if err != nil {
				s.log().WarnContext(ctx, "sweep: a kept sign-in waits for the next sweep", slog.String("grant", id), slog.String("err", err.Error()))
				first = cmpErr(first, err)
				continue
			}
			if !done {
				continue
			}
			if did {
				revoked++
			}
		}
		if _, _, err := s.KV.Delete(ctx, pkKept, it.SK); err != nil {
			first = cmpErr(first, err)
		}
	}
	return revoked, first
}

// settleKept revokes grant id's kept sign-in if no live grant uses it. done
// says the entry can go: the sign-in is revoked, or there is nothing left
// for this entry to do; revoked says this call revoked it.
func (s *Server) settleKept(ctx context.Context, tenant, id string) (done, revoked bool, err error) {
	g, err := s.loadGrant(ctx, tenant, id)
	if errors.Is(err, errNotFound) {
		return true, false, nil
	}
	if err != nil {
		return false, false, err
	}
	// A grant that hasn't ended (its disconnect wasn't saved) uses its own
	// sign-in: its own end handles it.
	if g.Status != StatusDisconnected && g.Status != StatusFailed {
		return true, false, nil
	}
	users, err := s.signInUsers(ctx, g)
	if err != nil {
		return false, false, err
	}
	if len(users) > 0 {
		return false, false, nil // still in use: never touched
	}
	// A provider account another ended grant already had revoked is done:
	// asked once.
	if g.Request.Connector != nil && g.Connector != nil {
		if receipt, err := s.accountRevoked(ctx, g); err != nil {
			return false, false, err
		} else if receipt != nil {
			g.Connector.Revocation = receipt
			if err := s.saveGrant(ctx, g); err != nil {
				return false, false, err
			}
			return true, false, nil
		}
	}
	// revokeNow reads the users again, then revokes at the vendor or the
	// provider and drops the token, or errs without dropping a live one.
	if err := s.revokeNow(ctx, g); err != nil {
		return false, false, err
	}
	g.event(s.now(), "the sign-in it kept for another connection was revoked: no connection uses it now")
	if err := s.saveGrant(ctx, g); err != nil {
		return false, false, err
	}
	return true, true, nil
}

// accountRevoked is the provider's confirmed revocation of g's connected
// account, if another grant on it already has one; nil when none does.
func (s *Server) accountRevoked(ctx context.Context, g *Grant) (*connector.Revocation, error) {
	items, err := s.KV.Query(ctx, tenantPK(g.Tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		var o Grant
		if json.Unmarshal(it.Doc, &o) != nil || o.ID == g.ID || o.Connector == nil || o.Connector.Revocation == nil {
			continue
		}
		if o.Connector.Revocation.Disabled && sameSignIn(g, &o) {
			r := *o.Connector.Revocation
			return &r, nil
		}
	}
	return nil, nil
}
