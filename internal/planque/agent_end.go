package planque

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
)

// capManageTafs is what KIFF calls managing the account's tafs: only a
// person who holds it ends a taf's connections.
const capManageTafs = "cards.manage"

// endedGrant is one connection ending a taf's agent ended (or, previewed,
// would end).
type endedGrant struct {
	Grant  string `json:"grant"`
	Vendor string `json:"vendor"`
	// Was is the status it had: a waiting or failed request ends too.
	Was string `json:"was"`
}

// keptSignIn is a sign-in or key one of the ended connections used that
// stays, because another taf's live connection uses it too.
type keptSignIn struct {
	Vendor string   `json:"vendor"`
	Why    string   `json:"why"`
	UsedBy []string `json:"used_by"`
	// Grants are the live connections that still hold the sign-in, and
	// Agents their agents: an app offers to remove the sign-in too only when
	// none of those agents is a taf the owner still has, by disconnecting
	// those grants (DELETE …/grants/{id}, signed). The last one gone, the
	// sign-in is revoked as any is.
	Grants []string `json:"grants"`
	Agents []string `json:"agents"`
}

// appName is g's app as people name it, for words drawn as they are: the
// reviewed catalog entry's vendor ("Qonto", not the entry's title), else the
// name the vendor or the managed app gave, else its server's host, else its
// API's.
func (s *Server) appName(g *Grant) string {
	if e, ok := s.CatalogEntries[s.appEntry(g)]; ok && e.Vendor != "" {
		return e.Vendor
	}
	if g.VendorName != "" {
		return g.VendorName
	}
	if h := hostOf(g.Request.Server); h != "" {
		return h
	}
	if g.Request.API != nil {
		if h := hostOf(g.Request.API.Base); h != "" {
			return h
		}
		return hostOf(g.Request.API.URL)
	}
	return ""
}

// endAgent ends every connection of a taf's agent when the taf is removed,
// whatever its status: connected, waiting for the owner, or failed. A
// connected one's tools leave KIFF and its sign-in or key is revoked,
// unless another taf's live connection uses it, which keeps it (and the
// answer says so). It only takes away, so it needs no device signature:
// the caller is the person's own KIFF session, holding cards.manage on
// this account. No key ends anyone's connections, not even an admin's.
//
// It is idempotent: a second call finds nothing left to end, and an agent
// with no connections is 200 with nothing ended, so a caller tells "no
// such route" (404) from "nothing to end". ?preview=1 ends nothing and
// answers what the call would end and keep, for the confirmation before
// the fingerprint. Every call is logged: who, which agent, what it ended.
//
// planque doesn't know tafs: the caller ends an agent only when its taf is
// removed, or when the agent has no taf.
func (s *Server) endAgent(w http.ResponseWriter, r *http.Request) {
	tenant, agent := r.PathValue("tenant"), r.PathValue("agent")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if me.APIKey || me.AgentID != "" || !slices.Contains(me.Capabilities, capManageTafs) {
		writeErr(w, http.StatusForbidden, "only a person who manages this account's tafs ends a taf's connections, with their own session")
		return
	}
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if len(strings.TrimSpace(string(body))) > 0 {
		var none struct{}
		if err := decodeStrict(body, &none); err != nil {
			s.fail(w, r, err)
			return
		}
	}
	preview := r.URL.Query().Get("preview") == "1"
	ctx := r.Context()

	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var mine, others []*Grant
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil {
			continue
		}
		g.ver = it.Ver
		if g.Request.Agent == agent {
			if g.Status != StatusDisconnected {
				mine = append(mine, &g)
			}
		} else if g.Status != StatusDisconnected && g.Status != StatusFailed {
			others = append(others, &g)
		}
	}

	// What ends and what stays, worked out before anything changes, so the
	// preview and the call answer the same.
	ended := []endedGrant{}
	kept := []keptSignIn{}
	for _, g := range mine {
		ended = append(ended, endedGrant{Grant: g.ID, Vendor: s.appName(g), Was: g.Status})
		if g.Status == StatusFailed {
			continue // it holds no sign-in in use
		}
		var usedBy, grants, agents []string
		for _, o := range others {
			if sameSignIn(g, o) {
				name := o.Request.AgentName
				if name == "" {
					name = "another taf"
				}
				if !slices.Contains(usedBy, name) {
					usedBy = append(usedBy, name)
				}
				grants = append(grants, o.ID)
				if !slices.Contains(agents, o.Request.Agent) {
					agents = append(agents, o.Request.Agent)
				}
			}
		}
		if len(usedBy) > 0 {
			kept = append(kept, keptSignIn{Vendor: s.appName(g), Why: "another taf's connection uses the same sign-in", UsedBy: usedBy, Grants: grants, Agents: agents})
		}
	}
	out := map[string]any{"agent": agent, "ended": ended, "kept": kept, "preview": preview}
	if preview {
		s.log().InfoContext(ctx, "agent end previewed", slog.String("tenant", tenant), slog.String("agent", agent),
			slog.String("by", me.Subject), slog.Int("ended", len(ended)), slog.Int("kept", len(kept)))
		writeJSON(w, http.StatusOK, out)
		return
	}

	// One at a time, each saved before the next: a sign-in two of this
	// agent's connections share is kept by the first and revoked with the
	// last.
	var key string
	var keyErr error
	if len(mine) > 0 {
		key, keyErr = s.connectorKey(ctx, tenant)
	}
	var unfinished []string
	for _, g := range mine {
		var problems []string
		if len(g.inKIFF()) > 0 {
			if keyErr != nil {
				unfinished = append(unfinished, g.ID)
				continue
			}
			left := false
			for _, t := range g.inKIFF() {
				if err := s.KIFF.Disconnect(ctx, key, t.Name); err != nil {
					problems = append(problems, "KIFF kept "+t.Name+": "+err.Error())
					left = true
				}
			}
			if left {
				// Its tools are still in KIFF: it stays as it is, for the
				// next call to end.
				s.log().WarnContext(ctx, "agent end: KIFF kept a tool", slog.String("grant", g.ID), slog.String("problems", strings.Join(problems, "; ")))
				unfinished = append(unfinished, g.ID)
				continue
			}
		}
		was := g.Status
		if was != StatusFailed {
			problems = append(problems, s.revoke(ctx, g)...)
		}
		g.Status, g.Transport, g.Transports, g.Connected, g.Declared = StatusDisconnected, "", nil, nil, false
		g.Pending, g.PendingTransport = nil, ""
		g.HeldUntil, g.ConfirmHash, g.ConfirmBy, g.Renewal = time.Time{}, "", time.Time{}, nil
		what := "ended with its taf: removed by a person who manages the account's tafs"
		switch {
		case was == StatusFailed:
			what += "; it had failed, nothing was in use"
		case g.keptFor > 0:
			what += "; the sign-in stays, another taf's connection uses it"
		case g.revokedShared:
			what += "; the sign-in it shared was revoked, no other connection used it"
		case g.TokenOf != "":
			what += "; the sign-in it used stays with the connection it came from"
		case g.Request.Kind == KindAPIKey:
			what += "; key deleted"
		case g.SignedIn || was == StatusConnected:
			what += "; token deleted"
		}
		if g.Request.API != nil {
			s.dropSnapshot(ctx, g.Tenant, g.ID)
		}
		if len(problems) > 0 {
			what += "; " + strings.Join(problems, "; ")
		}
		g.event(s.now(), what)
		if err := s.saveGrant(ctx, g); err != nil {
			s.log().WarnContext(ctx, "agent end: a grant wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
			unfinished = append(unfinished, g.ID)
		}
	}
	s.log().InfoContext(ctx, "agent ended", slog.String("tenant", tenant), slog.String("agent", agent),
		slog.String("by", me.Subject), slog.Int("ended", len(ended)-len(unfinished)), slog.Int("kept", len(kept)),
		slog.Int("unfinished", len(unfinished)))
	if len(unfinished) > 0 {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "some connections could not be ended yet; calling again goes on",
			"unfinished": unfinished, "agent": agent})
		return
	}
	writeJSON(w, http.StatusOK, out)
}
