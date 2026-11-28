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
}

// refreshRegistration updates, in KIFF, a connected grant's tools whose
// registration predates kiffRegistration: each tool's app, destructive hint
// and pins, in place (KIFF refuses to register a name already connected, and
// removing it would drop it from under the owner's Cards and held calls).
// The stamp (skRegistration) is set only once every tool is updated or gone. It never stops
// the call or read it runs in: a failure is logged, and tried again next
// time; a refusal (KIFF never lowers a tool's class) is logged and kept, so
// it shows on the next sweep.
func (s *Server) refreshRegistration(ctx context.Context, g *Grant) {
	if g.Status != StatusConnected || g.registering() || len(g.Tools) == 0 || s.KIFF == nil {
		return
	}
	var rec registrationRecord
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skRegistration(g.ID), &rec); err == nil && rec.Version == kiffRegistration {
		return
	}
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
		case errors.As(err, &a) && (a.Status == http.StatusConflict || a.Status == http.StatusUnprocessableEntity):
			s.log().WarnContext(ctx, "registration update refused", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("why", a.Msg))
			done = false
		default:
			s.log().InfoContext(ctx, "registration not updated yet", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("err", err.Error()))
			done = false
		}
	}
	if !done {
		return
	}
	it := store.Item{PK: tenantPK(g.Tenant), SK: skRegistration(g.ID)}
	if err := putDoc(ctx, s.KV, it, registrationRecord{Version: kiffRegistration, At: s.now().UTC().Format(time.RFC3339)}, store.Cond{}); err != nil {
		// Updated in KIFF, not stamped here: the next call updates again, to the same values.
		s.log().InfoContext(ctx, "registration updated, not stamped", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
}
