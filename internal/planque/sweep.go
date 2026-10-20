package planque

import (
	"context"
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

func dueSK(at time.Time, tenant, id string) string {
	return at.UTC().Format(time.RFC3339) + "#" + tenant + "#" + id
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

// markDue indexes g's deadline. It's kept a week past it, by which time the
// sweep has long removed it.
func (s *Server) markDue(ctx context.Context, g *Grant) {
	at := g.deadline()
	if at.IsZero() {
		return
	}
	it := store.Item{PK: pkDue, SK: dueSK(at, g.Tenant, g.ID), Doc: []byte("{}"), TTL: at.Add(7 * 24 * time.Hour).Unix()}
	if err := s.KV.Put(ctx, it, store.Cond{}); err != nil {
		s.log().WarnContext(ctx, "grant deadline not indexed", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
}

// Sweep expires every grant whose deadline has passed: a sign-in nobody
// confirmed, and a sign-in whose tools nobody approved, are revoked on time
// even if nobody reads them again. It returns how many it ended. Each entry
// is done with once its deadline passes: a grant that moved on in the
// meantime is left as it is.
func (s *Server) Sweep(ctx context.Context) (int, error) {
	// Only what's due: anyone may call the sweep, so a call reads no
	// further than its work ("~" sorts after any tenant or grant id).
	now, ended := s.now(), 0
	items, err := s.KV.QueryTo(ctx, pkDue, dueSK(now, "~", "~"))
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		at, tenant, id, ok := parseDue(it.SK)
		if ok && at.After(now) {
			break // in order of deadline: the rest are later
		}
		if ok {
			if g, err := s.loadGrant(ctx, tenant, id); err == nil {
				before := g.Status
				s.expire(ctx, g)
				if g.Status != before {
					ended++
				}
			}
		}
		_, _, _ = s.KV.Delete(ctx, pkDue, it.SK)
	}
	return ended, nil
}

// sweep is the scheduled call (EventBridge, hourly). Anyone may call it: it
// only does on time what the next read of each grant would do anyway.
func (s *Server) sweep(w http.ResponseWriter, r *http.Request) {
	n, err := s.Sweep(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if n > 0 {
		s.log().InfoContext(r.Context(), "sweep ended grants", slog.Int("ended", n))
	}
	writeJSON(w, http.StatusOK, map[string]int{"ended": n})
}
