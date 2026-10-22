package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// The owner's halt (kiff-cloud#1061): before planque sends anything to a
// vendor for an agent (a relay call, a Setup read), it asks KIFF, with the
// account's connector key, whether the agent or its account is halted. It
// refuses while either is, and refuses a permit issued at or before the
// latest halt even once resumed, whatever the permit's lifetime. When KIFF
// can't say, it refuses: a halt is never assumed lifted.

// haltRefusal is why nothing may be sent for agent of tenant now, or "".
// issued is the permit's issue time; zero for a send with no permit.
func (s *Server) haltRefusal(ctx context.Context, tenant, agent string, issued time.Time) string {
	key, err := s.connectorKey(ctx, tenant)
	if err != nil {
		return "planque couldn't check whether the owner halted this agent"
	}
	h, err := s.KIFF.HaltOf(ctx, key, agent)
	if err != nil {
		s.log().WarnContext(ctx, "halt unknown; refusing", slog.String("tenant", tenant), slog.String("agent", agent), slog.String("err", err.Error()))
		return "planque couldn't check whether the owner halted this agent"
	}
	if h.Halted {
		return "the owner halted this agent, or the whole account"
	}
	if h.HaltedAt != nil && !issued.IsZero() && !issued.After(*h.HaltedAt) {
		return "this permit was issued before the owner's halt"
	}
	return ""
}

// confirmHalt is KIFF asking planque to confirm a halt: planque reads it back
// from KIFF the way its dispatch does, and answers seen only if KIFF says the
// agent (or the account) is halted from at least that time. An account with
// no enrolment here has nothing planque could send for it.
func (s *Server) confirmHalt(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	var body struct {
		Agent    string    `json:"agent"`
		HaltedAt time.Time `json:"halted_at"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body) != nil || body.HaltedAt.IsZero() {
		writeErr(w, http.StatusBadRequest, "the body is {agent, halted_at}")
		return
	}
	key, err := s.connectorKey(r.Context(), tenant)
	if errors.Is(err, errNotFound) {
		writeJSON(w, http.StatusOK, map[string]any{"seen": true, "enrolled": false})
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	h, err := s.KIFF.HaltOf(r.Context(), key, body.Agent)
	if err != nil {
		writeErr(w, http.StatusBadGateway, "planque couldn't read the halt from KIFF")
		return
	}
	seen := h.Halted && h.HaltedAt != nil && !h.HaltedAt.Before(body.HaltedAt)
	writeJSON(w, http.StatusOK, map[string]any{"seen": seen, "enrolled": true})
}
