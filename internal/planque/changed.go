package planque

import (
	"context"
	"log/slog"
	"time"

	"github.com/kiff/planque/internal/jcs"
)

// definitionsReadEvery is how often, at most, planque reads a direct MCP
// server's tools again before a call, to notice one changed since the owner
// approved it (kiff/planque#71). Read on use, not on a schedule.
const definitionsReadEvery = 10 * time.Minute

// Pause is why a grant sends nothing: a tool the owner approved changed at
// the vendor, or is gone. Every call on the grant is refused before anything
// is sent. A pause is never lifted on the same grant: not by a new sign-in
// (a reconnect reads the tools again without the owner reviewing them), nor
// by adding tools (that keeps the approved ones as they were). The owner
// connects again, approving the current definitions on a new grant.
type Pause struct {
	Tool string    `json:"tool"`
	At   time.Time `json:"at"`
	Why  string    `json:"why"`
}

// pausedRefusal refuses a call on a paused grant; "" when it isn't paused.
func pausedRefusal(g *Grant) (code, why string) {
	if g.Paused == nil {
		return "", ""
	}
	return CodeOperationChanged, g.Paused.Why + "; disconnect it and connect again to approve its current definition"
}

// vendorSchemaSHA256 is the canonical hash of a tool's input schema as the
// vendor lists it: what the owner approved, compared later.
func vendorSchemaSHA256(schema []byte) string {
	if len(schema) == 0 {
		schema = []byte("{}")
	}
	sum, err := jcs.SHA256(schema)
	if err != nil {
		return ""
	}
	return sum
}

// checkDefinitions reads a direct MCP server's tools, at most every
// definitionsReadEvery, and pauses g if one it was granted is gone or its
// input schema changed since the owner approved it. A list it can't read
// leaves the call to its usual path, which says what happened (a vendor
// that signed the grant out, one that can't be reached); the next call
// reads again. Only a list read and found changed pauses. A tool approved
// before planque kept its schema's hash gets one on the first read.
func (s *Server) checkDefinitions(ctx context.Context, g *Grant) (code, why string) {
	// Only a direct MCP server lists tools to read again: a connector checks
	// its pinned definition at each call, and an API or browser grant sends
	// only what planque itself defines.
	if g.Request.Connector != nil || g.Request.API != nil || g.Request.Kind == KindBrowser || len(g.Tools) == 0 {
		return "", ""
	}
	now := s.now()
	if !g.DefinitionsReadAt.IsZero() && now.Sub(g.DefinitionsReadAt) < definitionsReadEvery {
		return "", ""
	}
	access, _, err := s.accessToken(ctx, g, 0)
	if err != nil {
		return "", ""
	}
	listed, _, err := s.MCP.ListTools(ctx, g.Request.Server, access)
	if err != nil {
		s.log().InfoContext(ctx, "definitions not read; read again at the next call", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return "", ""
	}
	current := make(map[string][]byte, len(listed))
	for _, t := range listed {
		current[t.Name] = t.InputSchema
	}
	for i, t := range g.Tools {
		schema, ok := current[t.Tool]
		if !ok {
			return s.pauseGrant(ctx, g, t.Name, vendorName(g)+" no longer lists "+t.Tool)
		}
		sum := vendorSchemaSHA256(schema)
		switch {
		case sum == "":
			// A schema that isn't JSON: changed from what was approved, which was.
			return s.pauseGrant(ctx, g, t.Name, vendorName(g)+" changed "+t.Tool+" since you approved it")
		case t.VendorSchemaSHA256 == "":
			g.Tools[i].VendorSchemaSHA256 = sum
		case t.VendorSchemaSHA256 != sum:
			return s.pauseGrant(ctx, g, t.Name, vendorName(g)+" changed "+t.Tool+" since you approved it")
		}
	}
	g.DefinitionsReadAt = now
	if err := s.saveGrant(ctx, g); err != nil {
		// Checked and unchanged: the call goes on; the next one reads again.
		s.log().WarnContext(ctx, "definitions check not kept", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	return "", ""
}

// pauseGrant pauses g for tool, says why on the grant and in the log, and
// refuses the call that noticed it.
func (s *Server) pauseGrant(ctx context.Context, g *Grant, tool, why string) (code, reason string) {
	now := s.now()
	g.Paused = &Pause{Tool: tool, At: now.UTC(), Why: why}
	g.event(now, "paused: "+why+"; nothing is sent until the owner connects again and approves its current definition")
	if err := s.saveGrant(ctx, g); err != nil {
		s.log().ErrorContext(ctx, "pause not kept", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	s.log().WarnContext(ctx, "grant paused: an approved tool changed at the vendor", slog.String("grant", g.ID), slog.String("tool", tool))
	return pausedRefusal(g)
}
