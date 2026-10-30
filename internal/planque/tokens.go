package planque

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

type mcpTool = mcp.Tool

// refreshEarly refreshes a token this long before it expires.
const refreshEarly = time.Minute

// refreshLease is how long one refresher holds the token record.
const refreshLease = 20 * time.Second

func tokenBinding(g *Grant) map[string]string {
	return map[string]string{"tenant": g.Tenant, "grant": g.ID, "kind": "token"}
}

// storeToken seals tok for g. With prev, only if the record is still at
// that version; without, only if there is none.
func (s *Server) storeToken(ctx context.Context, g *Grant, tok oauth.Token, prev *int64) error {
	b, _ := json.Marshal(tok)
	sealed, err := s.Sealer.Seal(ctx, b, tokenBinding(g))
	clear(b)
	if err != nil {
		return err
	}
	it := store.Item{PK: tenantPK(g.Tenant), SK: skToken(g.ID), Ver: 1}
	c := store.Cond{NotExists: true}
	if prev != nil {
		it.Ver, c = *prev+1, store.Cond{Ver: prev}
	}
	return putDoc(ctx, s.KV, it, tokenRecord{Sealed: sealed, Expiry: tok.Expiry}, c)
}

func (s *Server) openToken(ctx context.Context, g *Grant) (oauth.Token, store.Item, error) {
	var rec tokenRecord
	it, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skToken(g.ID), &rec)
	if err != nil {
		return oauth.Token{}, it, err
	}
	b, err := s.Sealer.Open(ctx, rec.Sealed, tokenBinding(g))
	if err != nil {
		return oauth.Token{}, it, err
	}
	var tok oauth.Token
	err = json.Unmarshal(b, &tok)
	clear(b)
	tok.Expiry = rec.Expiry // the record's, which storeToken writes with the token
	return tok, it, err
}

func (s *Server) dropToken(ctx context.Context, g *Grant) {
	_, _, _ = s.KV.Delete(ctx, tenantPK(g.Tenant), skToken(g.ID))
}

var errNoToken = errors.New("planque holds no token for this grant")

// accessToken returns a usable access token for g and the version of
// the record it came from. stale, when not 0, is a version the vendor
// refused: the token is refreshed unless someone already replaced it.
func (s *Server) accessToken(ctx context.Context, g *Grant, stale int64) (string, int64, error) {
	if g.signedOut() {
		return "", 0, signedOutErr(g)
	}
	tok, it, err := s.openToken(ctx, g)
	if errors.Is(err, errNotFound) {
		return "", 0, errNoToken
	}
	if err != nil {
		return "", 0, err
	}
	fresh := tok.Expiry.IsZero() || s.now().Add(refreshEarly).Before(tok.Expiry)
	if (stale == 0 && fresh) || (stale != 0 && it.Ver != stale) {
		return tok.Access, it.Ver, nil
	}
	return s.refresh(ctx, g, it.Ver)
}

// refresh replaces the token at version ver, once. In this process one
// call refreshes per grant; across processes the record's lease decides,
// and the others wait for the version after ver. A refresh token the
// vendor rotates is kept; one it does not send again stays.
func (s *Server) refresh(ctx context.Context, g *Grant, ver int64) (string, int64, error) {
	mu, _ := s.refreshing.LoadOrStore(g.ID, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	tok, it, err := s.openToken(ctx, g)
	if err != nil {
		return "", 0, err
	}
	if it.Ver != ver {
		return tok.Access, it.Ver, nil // refreshed while we waited
	}
	now := s.now()
	leased := it
	leased.Lease = now.Add(refreshLease).UnixMilli()
	err = s.KV.Put(ctx, leased, store.Cond{Ver: &ver, LeaseFreeAt: now.UnixMilli()})
	if errors.Is(err, store.ErrCondition) {
		return s.waitRefresh(ctx, g, ver)
	}
	if err != nil {
		return "", 0, err
	}
	if tok.Refresh == "" {
		s.releaseLease(ctx, it, ver)
		return "", 0, s.signOut(ctx, g, "its token expired and the vendor gave no refresh token")
	}
	next, err := s.OAuth.Refresh(ctx, g.Vendor, g.ClientID, tok.Refresh)
	if err != nil {
		s.releaseLease(ctx, it, ver)
		if errors.Is(err, oauth.ErrInvalidGrant) {
			return "", 0, s.signOut(ctx, g, "the vendor refused its refresh (invalid_grant)")
		}
		return "", 0, err
	}
	if next.Refresh == "" {
		next.Refresh = tok.Refresh
	}
	if next.Account == "" {
		next.Account = tok.Account
	}
	if tok.Account != "" && next.Account != tok.Account {
		s.releaseLease(ctx, it, ver)
		return "", 0, errors.New("the refreshed token is for another vendor account")
	}
	if err := s.storeToken(ctx, g, next, &ver); err != nil {
		return "", 0, err
	}
	return next.Access, ver + 1, nil
}

// releaseLease clears the lease so another call may try the refresh.
func (s *Server) releaseLease(ctx context.Context, it store.Item, ver int64) {
	it.Lease = 0
	_ = s.KV.Put(ctx, it, store.Cond{Ver: &ver})
}

func (s *Server) waitRefresh(ctx context.Context, g *Grant, ver int64) (string, int64, error) {
	wait := s.RefreshWait
	if wait <= 0 {
		wait = refreshLease
	}
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return "", 0, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		tok, it, err := s.openToken(ctx, g)
		if err != nil {
			return "", 0, err
		}
		if it.Ver != ver {
			return tok.Access, it.Ver, nil
		}
	}
	return "", 0, errors.New("another refresh of the vendor token did not finish in time")
}

// signedOut says the vendor ended g's sign-in: it waits for the owner to
// sign in again, or did and that sign-in isn't confirmed yet.
func (g *Grant) signedOut() bool {
	return g.Status == StatusNeedsReconnect || g.Reconnect
}

// signedOutErr is what a call through a signed-out grant answers: the
// taf says it can't be done now, and the owner reads why in tafs.
func signedOutErr(g *Grant) error {
	return fmt.Errorf("%s ended the sign-in planque holds; the owner signs in to %s again in tafs", vendorName(g), vendorName(g))
}

// signOut marks a connected grant as signed out by its vendor, so its
// calls are refused saying why and the owner's next sign-in really signs
// in, rather than finding the grant connected. It reads the grant again,
// so a sign-in someone else just confirmed isn't undone.
func (s *Server) signOut(ctx context.Context, g *Grant, why string) error {
	now, err := s.loadGrant(ctx, g.Tenant, g.ID)
	if err == nil && now.Status == StatusConnected {
		now.Status, now.Reason = StatusNeedsReconnect, why
		now.event(s.now(), "signed out by "+vendorName(now)+": "+why+"; waiting for the owner to sign in again")
		if err = s.saveGrant(ctx, now); err == nil {
			*g = *now
		}
	}
	if err != nil && !errors.Is(err, store.ErrCondition) {
		s.log().WarnContext(ctx, "signed-out grant not marked", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	s.log().InfoContext(ctx, "grant signed out by its vendor", slog.String("grant", g.ID), slog.String("why", why))
	return signedOutErr(g)
}
