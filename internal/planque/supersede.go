package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

// supersedePending ends the account's earlier requests that still wait for the
// owner when g, a newer request from the same agent for the same app, is made.
// Each time Coach asks again (a refusal, a new answer) the old request would
// otherwise stay "waiting for your approval" until its seven days ran out, and
// the person would read one waiting row per try. An earlier request that holds
// a sign-in, a connection, or another taf's sign-in is left alone; so is one
// from another agent or for another app. The ended request keeps its record,
// as disconnected, with tools declared for it taken out of KIFF first.
func (s *Server) supersedePending(ctx context.Context, g *Grant) {
	if g.Request.Agent == "" || g.Request.CatalogEntry == "" {
		return
	}
	items, err := s.KV.Query(ctx, tenantPK(g.Tenant), "GRANT#")
	if err != nil {
		s.log().WarnContext(ctx, "earlier requests weren't read to end them", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return
	}
	for _, it := range items {
		var old Grant
		if json.Unmarshal(it.Doc, &old) != nil || old.ID == "" || old.ID == g.ID || skGrant(old.ID) != it.SK {
			continue
		}
		old.ver = it.Ver
		if !old.waitsForApproval() || old.TokenOf != "" || old.Request.Agent != g.Request.Agent || old.Request.CatalogEntry != g.Request.CatalogEntry {
			continue
		}
		if err := s.undeclare(ctx, &old); err != nil {
			s.log().WarnContext(ctx, "an earlier request's declared tools stayed in KIFF", slog.String("grant", old.ID), slog.String("err", err.Error()))
			continue
		}
		old.HeldUntil = time.Time{}
		from := old.Status
		old.Status, old.Reason = StatusDisconnected, ""
		old.event(s.now(), "replaced by a newer request for the same app from the same agent: "+g.ID)
		if err := s.saveGrant(ctx, &old); err != nil {
			s.log().WarnContext(ctx, "an earlier request wasn't ended", slog.String("grant", old.ID), slog.String("err", err.Error()))
			continue
		}
		s.logEnded(ctx, &old, from, "a newer request, "+g.ID)
	}
}

// logEnded says a grant ended: which, from which status to which, and who or
// what ended it, so a row that went from waiting to gone can be explained.
func (s *Server) logEnded(ctx context.Context, g *Grant, from, by string) {
	s.log().InfoContext(ctx, "grant ended", slog.String("grant", g.ID), slog.String("from", from), slog.String("to", g.Status),
		slog.String("by", by), slog.String("agent", g.Request.Agent), slog.String("entry", g.Request.CatalogEntry))
}
