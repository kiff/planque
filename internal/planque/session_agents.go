package planque

import (
	"encoding/json"
	"net/http"
	"regexp"
	"slices"
)

// A request's session is the chat that asked for it, as the agents runtime's
// runner names it: an id, never words. planque keeps it beside the grant
// (Grant.AskedIn), outside the request, so it's in no request served and no
// digest signed, and clears it when the grant ends.
var askedInShape = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// askedInWant is the refusal of a session that isn't an id.
const askedInWant = "session is the chat's id as the runner gives it: 1 to 128 letters, digits, - or _"

// askedInOK says whether s is no session at all, or one shaped as an id.
func askedInOK(s string) bool { return s == "" || askedInShape.MatchString(s) }

// sessionAgent is one connection a chat asked for, as a chat's delete reads
// it: the agent it was asked under, where it stands, and the day it was
// asked. No grant id, token, tool or argument.
type sessionAgent struct {
	Agent  string `json:"agent"`
	Status string `json:"status"`
	Since  string `json:"since,omitempty"`
}

// sessionAgents answers which agents a chat asked connections for
// (GET /v1/accounts/{tenant}/sessions/{session}/agents): the grants whose
// request carried that session and that haven't ended or failed, for any of the
// account's keys or sessions, as the connections list is. tafs cloud reads it
// when the chat is deleted and ends, through the end route and with the
// person's own session, those no live taf runs. planque doesn't know tafs, so
// it says which, and decides nothing.
func (s *Server) sessionAgents(w http.ResponseWriter, r *http.Request) {
	tenant, session := r.PathValue("tenant"), r.PathValue("session")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	if session == "" || !askedInOK(session) {
		writeErr(w, http.StatusBadRequest, askedInWant)
		return
	}
	items, err := s.KV.Query(r.Context(), tenantPK(tenant), "GRANT#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	out := []sessionAgent{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.Tenant != tenant || g.AskedIn != session || g.Request.Agent == "" {
			continue
		}
		if g.Status == StatusDisconnected || g.Status == StatusFailed {
			continue
		}
		if slices.ContainsFunc(out, func(a sessionAgent) bool { return a.Agent == g.Request.Agent }) {
			continue
		}
		out = append(out, sessionAgent{Agent: g.Request.Agent, Status: g.Status, Since: sinceOf(&g)})
	}
	writeJSON(w, http.StatusOK, map[string]any{"agents": out})
}
