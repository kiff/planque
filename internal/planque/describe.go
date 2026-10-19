package planque

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/store"
)

// What a describe call may read, at most (kiff/tafs#89).
const (
	// MaxDescribe is the most bytes of a vendor's answer planque returns.
	MaxDescribe = 32 << 10
	// MaxDescribes is how many describe calls one grant may make.
	MaxDescribes = 10
)

// describe lets the agent that asked for a grant read what the owner's
// sign-in covers, while it looks (kiff/tafs#89): it calls one vendor tool
// that the vendor marks read-only and that takes no arguments. Such a tool
// describes the session itself (Stripe's list_available_accounts_or_orgs
// lists the accounts and their modes), and with no arguments it can't be
// aimed at the owner's data. Nothing else is called, and nothing is given
// to the taf. The agent then pins what the vendor reported, instead of
// guessing.
func (s *Server) describe(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		Tool string `json:"tool"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if me.AgentID == "" || me.AgentID != g.Request.AskedBy {
		writeErr(w, http.StatusForbidden, "only the agent that asked for this grant may read what its sign-in covers")
		return
	}
	s.expire(r.Context(), g)
	// Only what the owner approved: a grant listed after connecting first,
	// on the text that says its sign-in may be read. Not a connected grant,
	// whose approval named its tools, nor one approved on an older text.
	if g.Status != StatusListed {
		writeErr(w, http.StatusConflict, "this grant is "+g.Status+": read what the sign-in covers while connecting first, before the tools")
		return
	}
	if !g.Describable {
		writeErr(w, http.StatusConflict, "this grant was approved on a text that didn't say its sign-in would be read; connect again to read it")
		return
	}
	i := slices.IndexFunc(g.Catalog, func(t CatalogTool) bool { return t.Name == req.Tool })
	if i < 0 {
		writeErr(w, http.StatusNotFound, "the vendor lists no tool "+req.Tool)
		return
	}
	if t := g.Catalog[i]; !t.ReadOnly || !takesNoArguments(t.InputSchema) {
		writeErr(w, http.StatusForbidden, "only a tool the vendor marks read-only and that takes no arguments describes the sign-in; "+req.Tool+" isn't one")
		return
	}
	n := 0
	for _, ev := range g.Events {
		if strings.HasPrefix(ev.What, "read ") {
			n++
		}
	}
	if n >= MaxDescribes {
		writeErr(w, http.StatusTooManyRequests, "this grant has read what its sign-in covers enough times")
		return
	}

	// Recorded first: a read the record doesn't show never happens, and two
	// racing reads can't both pass the bound (the write is conditional).
	g.event(s.now(), "read "+req.Tool+" at "+hostOf(g.Request.Server)+", to see what the sign-in covers")
	if err := s.saveGrant(r.Context(), g); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, "the grant changed while this read was checked; read again")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	access, ver, err := s.accessToken(r.Context(), g, 0)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var res json.RawMessage
	for attempt := 0; ; attempt++ {
		res, err = s.MCP.CallTool(r.Context(), g.Request.Server, access, req.Tool, json.RawMessage(`{}`))
		if !errors.Is(err, mcp.ErrUnauthorized) || attempt > 0 {
			break
		}
		// The vendor refused the token: refresh it once and try again.
		if access, ver, err = s.accessToken(r.Context(), g, ver); err != nil {
			s.fail(w, r, err)
			return
		}
	}
	var rpc *mcp.RPCError
	switch {
	case errors.As(err, &rpc):
		writeErr(w, http.StatusBadGateway, hostOf(g.Request.Server)+" returned an error: "+rpc.Message)
		return
	case err != nil:
		writeErr(w, http.StatusBadGateway, hostOf(g.Request.Server)+" could not be reached, or refused the sign-in")
		return
	}
	truncated := len(res) > MaxDescribe
	if truncated {
		res = nil
	}
	out := map[string]any{"grant": g.ID, "tool": req.Tool, "result": res}
	if truncated {
		out["truncated"] = true
	}
	writeJSON(w, http.StatusOK, out)
}

// takesNoArguments reports whether a tool's input schema asks for nothing:
// no properties, and nothing required.
func takesNoArguments(schema json.RawMessage) bool {
	if len(schema) == 0 {
		return true
	}
	var s struct {
		Properties map[string]json.RawMessage `json:"properties"`
		Required   []string                   `json:"required"`
	}
	if err := json.Unmarshal(schema, &s); err != nil {
		return false
	}
	return len(s.Properties) == 0 && len(s.Required) == 0
}
