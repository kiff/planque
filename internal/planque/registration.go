package planque

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/store"
)

// kiffRegistration names the relay-owned fields planque registers each tool
// with: its app (K4) and its vendor's destructive hint, and its vendor pins
// (step 3). A grant whose tools were registered before them carries an
// older stamp, and is updated in place (refreshRegistration).
const kiffRegistration = "app-pins-1"

// skRegistration keeps, beside a grant, the version its tools were last
// updated with in KIFF. It's its own item, never the grant's record, so the
// update never races a registration that is saving the grant (KIFF probes the
// relay while a grant's tools are registered).
func skRegistration(grant string) string { return "KIFFREG#" + grant + "#" }

type registrationRecord struct {
	Version string `json:"version"`
	At      string `json:"at"`
	// Retry: after a transient failure (KIFF unreachable, a 5xx), no update is tried again before this time, so a
	// slow or down KIFF costs a call at most one try in registrationRetry, never one per call.
	Retry string `json:"retry,omitempty"`
}

// registrationRetry is how long a grant waits after a transient failure before its tools are tried again.
const registrationRetry = 15 * time.Minute

// registrationTimeout bounds one refresh, every tool together: it runs beside a call, so a slow KIFF never holds the
// call longer than this, and the call's own cancel never cuts a PATCH in half.
const registrationTimeout = 3 * time.Second

// refreshRegistration updates, in KIFF, a connected grant's tools whose
// registration predates kiffRegistration: each tool's app, destructive hint
// and pins, in place (KIFF refuses to register a name already connected, and
// removing it would drop it from under the owner's Cards and held calls).
// The stamp (skRegistration) is set once every tool is updated, gone, or
// refused for good (403, 409, 422: logged once). A transient failure keeps
// the older version and waits registrationRetry before trying again. It never
// stops the call or read it runs in, and holds it at most registrationTimeout.
func (s *Server) refreshRegistration(ctx context.Context, g *Grant) {
	if g.Status != StatusConnected || g.registering() || len(g.Tools) == 0 || s.KIFF == nil {
		return
	}
	var rec registrationRecord
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skRegistration(g.ID), &rec); err == nil {
		if rec.Version == kiffRegistration {
			return
		}
		if retry, err := time.Parse(time.RFC3339, rec.Retry); err == nil && s.now().Before(retry) {
			return
		}
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), registrationTimeout)
	defer cancel()
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		s.log().WarnContext(ctx, "registration not updated: no connector key", slog.String("grant", g.ID))
		return
	}
	app, done := s.appSlug(g), true
	for _, t := range g.Tools {
		err := s.KIFF.UpdateConnection(ctx, key, t.Name, kiff.Update{App: app, DestructiveHint: t.DestructiveHint, VendorPins: toolPins(g, t)})
		var a *kiff.Answered
		switch {
		case err == nil:
		case errors.As(err, &a) && a.Status == http.StatusNotFound:
			// The tool isn't connected in KIFF any more: nothing to update.
		case errors.As(err, &a) && a.Status == http.StatusForbidden:
			// KIFF says the tool isn't this relay's to update (an owner's own connection by that name): that won't
			// change, so it's said once and not asked again on every call.
			s.log().WarnContext(ctx, "registration update not ours", slog.String("grant", g.ID), slog.String("tool", t.Name))
		case errors.As(err, &a) && (a.Status == http.StatusConflict || a.Status == http.StatusUnprocessableEntity):
			// KIFF refuses what it always will (never lowering a class, an app already set, a hint going back to
			// false): final for this tool. Said once, in the log a sweep reads, and never asked again on a call.
			s.log().WarnContext(ctx, "registration update refused", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("why", a.Msg))
		default:
			s.log().InfoContext(ctx, "registration not updated yet", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("err", err.Error()))
			done = false
		}
	}
	it := store.Item{PK: tenantPK(g.Tenant), SK: skRegistration(g.ID)}
	now := s.now().UTC()
	next := registrationRecord{Version: kiffRegistration, At: now.Format(time.RFC3339)}
	if !done {
		// A transient failure: the older version stays, and the next try waits registrationRetry.
		next = registrationRecord{Version: rec.Version, At: now.Format(time.RFC3339), Retry: now.Add(registrationRetry).Format(time.RFC3339)}
	}
	if err := putDoc(ctx, s.KV, it, next, store.Cond{}); err != nil {
		// Updated in KIFF, not stamped here: the next call updates again, to the same values.
		s.log().InfoContext(ctx, "registration updated, not stamped", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
}
