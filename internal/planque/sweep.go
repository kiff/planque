package planque

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kiff/planque/internal/store"
)

// pkDue holds one item per grant that has a deadline, keyed by when
// (dueSK), so the scheduled sweep finds what's due without reading every
// account. A grant missing from it is still expired on its next read.
const pkDue = "DUE"

// dueSK keys a deadline by its whole second, rounded up, so an entry is
// never due before its deadline (kiff/planque#23 C). Entries written before
// were rounded down; the sweep reads either.
func dueSK(at time.Time, tenant, id string) string {
	return ceilSecond(at).Format(time.RFC3339) + "#" + tenant + "#" + id
}

func ceilSecond(at time.Time) time.Time {
	at = at.UTC()
	if t := at.Truncate(time.Second); t.Before(at) {
		return t.Add(time.Second)
	}
	return at
}

func parseDue(sk string) (at time.Time, tenant, id string, ok bool) {
	parts := strings.SplitN(sk, "#", 3)
	if len(parts) != 3 {
		return time.Time{}, "", "", false
	}
	at, err := time.Parse(time.RFC3339, parts[0])
	return at, parts[1], parts[2], err == nil
}

// deadline is when g must be confirmed (a sign-in) or its tools approved (a
// held sign-in), as expire reads them; zero when it has none.
func (g *Grant) deadline() time.Time {
	if g.Status == StatusSignedIn && !g.ConfirmBy.IsZero() {
		return g.ConfirmBy
	}
	if g.SignedIn && (g.Status == StatusListed || g.Status == StatusPending) && !g.HeldUntil.IsZero() {
		return g.HeldUntil
	}
	return time.Time{}
}

// overdue says g's deadline passed. Such a grant acts on nothing, whether
// or not its expiry could be saved yet: the sweep or the next read finishes
// revoking it.
func (s *Server) overdue(g *Grant) bool {
	d := g.deadline()
	return !d.IsZero() && s.now().After(d)
}

// markDue indexes g's deadline, before g is saved with it: a grant is never
// stored with a deadline the sweep can't find (kiff/planque#23 B). An entry
// whose save then fails is harmless, since the sweep checks the grant
// itself. It's kept a week past its deadline, by which time the sweep has
// long removed it.
func (s *Server) markDue(ctx context.Context, g *Grant) error {
	at := g.deadline()
	if at.IsZero() {
		return nil
	}
	it := store.Item{PK: pkDue, SK: dueSK(at, g.Tenant, g.ID), Doc: []byte("{}"), TTL: at.Add(7 * 24 * time.Hour).Unix()}
	if err := s.KV.Put(ctx, it, store.Cond{}); err != nil {
		return fmt.Errorf("index the grant's deadline: %w", err)
	}
	return nil
}

// Sweep expires every grant whose deadline has passed: a sign-in nobody
// confirmed, and a sign-in whose tools nobody approved, are revoked on time
// even if nobody reads them again. It returns how many it ended, and the
// first error. An entry goes only once its work is done (kiff/planque#23 A):
// the grant is gone, it moved on to another deadline, or its expiry is
// saved with its tokens revoked. Any failure keeps the entry for the next
// sweep.
func (s *Server) Sweep(ctx context.Context) (int, error) {
	// Only what's due: anyone may call the sweep, so a call reads no
	// further than its work ("~" sorts after any tenant or grant id).
	now, ended := s.now(), 0
	items, err := s.KV.QueryTo(ctx, pkDue, now.UTC().Truncate(time.Second).Format(time.RFC3339)+"#~#~")
	if err != nil {
		return 0, err
	}
	var first error
	for _, it := range items {
		at, tenant, id, ok := parseDue(it.SK)
		if ok && at.After(now) {
			break // in order of deadline: the rest are later
		}
		if ok {
			done, n, err := s.settleDue(ctx, at, tenant, id)
			ended += n
			if err != nil {
				s.log().WarnContext(ctx, "sweep: a grant's expiry waits for the next sweep", slog.String("grant", id), slog.String("err", err.Error()))
				first = cmpErr(first, err)
				continue
			}
			if !done {
				continue
			}
		}
		if _, _, err := s.KV.Delete(ctx, pkDue, it.SK); err != nil {
			first = cmpErr(first, err)
		}
	}
	return ended, first
}

func cmpErr(first, err error) error {
	if first != nil {
		return first
	}
	return err
}

// settleDue does what the due entry at `at` asks of grant id: done says the
// entry can go, ended whether this ended the grant.
func (s *Server) settleDue(ctx context.Context, at time.Time, tenant, id string) (done bool, ended int, err error) {
	g, err := s.loadGrant(ctx, tenant, id)
	if errors.Is(err, errNotFound) {
		return true, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	d := g.deadline()
	// Moved on: no deadline now, or another one, which has its own entry.
	// An entry is d's whole second, rounded up (or, written before, down).
	if d.IsZero() || d.UTC().Truncate(time.Second).After(at) || ceilSecond(d).Before(at) {
		return true, 0, nil
	}
	if !s.now().After(d) {
		return false, 0, nil
	}
	ok, err := s.expire(ctx, g)
	if err != nil {
		return false, 0, err
	}
	if !ok {
		return false, 0, nil
	}
	return true, 1, nil
}

// sweep is the scheduled call (EventBridge, hourly). Anyone may call it: it
// only does on time what the next read of each grant would do anyway.
func (s *Server) sweep(w http.ResponseWriter, r *http.Request) {
	n, err := s.Sweep(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	// Sign-outs KIFF didn't take when they happened, sent again.
	if told, err := s.tellAgain(r.Context()); err != nil {
		s.log().WarnContext(r.Context(), "sweep: sign-outs not told again", slog.String("err", err.Error()))
	} else if told > 0 {
		s.log().InfoContext(r.Context(), "sweep told sign-outs to KIFF", slog.Int("told", told))
	}
	if n > 0 {
		s.log().InfoContext(r.Context(), "sweep ended grants", slog.Int("ended", n))
	}
	writeJSON(w, http.StatusOK, map[string]int{"ended": n})
}
