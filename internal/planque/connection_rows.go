package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	// The owner's zone (X-Time-Zone) names any IANA zone, wherever planque runs.
	_ "time/tzdata"

	"github.com/kiff/planque/internal/store"
	"github.com/kiff/planque/internal/words"
)

// A connection as the owner's list shows it, the same in every app: which
// tafs use it, its health in plain words, and only the actions that can work
// for it. An app draws these as sent and keeps no rule of its own, so it
// never offers a button that answers with a state ("this connection is
// listed"), and never shows a provider's id.

// connectionHealth is a connection's health for people. State is one of the
// health states below; OK is true only while it works; Words says it in the
// language asked, its time in the owner's zone (X-Time-Zone, else UTC). At is
// the time Words names, when it names one: the last use, or the sign-out.
type connectionHealth struct {
	State string    `json:"state"`
	OK    bool      `json:"ok"`
	Words string    `json:"words"`
	At    time.Time `json:"at,omitzero"`
}

// The health states.
const (
	HealthWorking    = "working"
	HealthPaused     = "paused"
	HealthSignedOut  = "signed_out"
	HealthUnfinished = "unfinished"
	HealthWaiting    = "waiting"
	HealthEnded      = "ended"
)

// connectionAction is one thing the owner can do with a connection, that
// works for it now. How is the call it takes, which the apps already make:
// "approve" (…/approve, signed: a new sign-in on this same connection),
// "key" (…/key, signed, with the key) or "delete" (DELETE …/grants/{id},
// signed). Stops, on remove, says what stops, for the app's one more ask.
type connectionAction struct {
	ID    string `json:"id"`
	Words string `json:"words"`
	How   string `json:"how"`
	Stops string `json:"stops,omitempty"`
}

// The actions.
const (
	ActionSignInAgain      = "sign_in_again"
	ActionRemove           = "remove"
	ActionRemoveUnfinished = "remove_unfinished"
)

// usedDoc is when a connection last did what a taf asked, kept beside the
// grant (never in it, so a call never races an approval's write).
type usedDoc struct {
	At time.Time `json:"at"`
}

func skUsed(id string) string { return "USED#" + id }

// noteUsed keeps when g last answered a taf's call without an error. Best
// effort: a write that fails only leaves the last use older.
func (s *Server) noteUsed(ctx context.Context, g *Grant, res any) {
	if g == nil || isErrorResult(res) {
		return
	}
	_ = putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skUsed(g.ID)}, usedDoc{At: s.now().UTC()}, store.Cond{})
}

// isErrorResult says a tool result is an error (MCP's isError).
func isErrorResult(res any) bool {
	if m, ok := res.(map[string]any); ok {
		b, _ := m["isError"].(bool)
		return b
	}
	b, err := json.Marshal(res)
	if err != nil {
		return true
	}
	var r struct {
		IsError bool `json:"isError"`
	}
	return json.Unmarshal(b, &r) != nil || r.IsError
}

// lastUses is when each of tenant's connections was last used, by grant id.
func (s *Server) lastUses(ctx context.Context, tenant string) map[string]time.Time {
	out := map[string]time.Time{}
	items, err := s.KV.Query(ctx, tenantPK(tenant), "USED#")
	if err != nil {
		return out
	}
	for _, it := range items {
		var u usedDoc
		if json.Unmarshal(it.Doc, &u) == nil && !u.At.IsZero() {
			out[strings.TrimPrefix(it.SK, "USED#")] = u.At
		}
	}
	return out
}

// ownerZone is the zone the owner's app names (X-Time-Zone, an IANA name),
// and whether it named one; UTC when not.
func ownerZone(r *http.Request) (*time.Location, bool) {
	if r != nil {
		if z := strings.TrimSpace(r.Header.Get("X-Time-Zone")); z != "" && len(z) <= 64 {
			if loc, err := time.LoadLocation(z); err == nil {
				return loc, true
			}
		}
	}
	return time.UTC, false
}

// rowContext is what a row's words depend on: the language asked, the
// owner's zone, now, and the connection's last use when it was looked up.
type rowContext struct {
	lang  string
	zone  *time.Location
	named bool
	now   time.Time
	used  time.Time
	known bool
}

func (rc rowContext) when(t time.Time) string {
	c := words.Connection(rc.lang)
	s := c.When(t.In(rc.zone), rc.now.In(rc.zone))
	if !rc.named && !strings.ContainsAny(s, " .") {
		s += " UTC"
	}
	return s
}

func (rc rowContext) day(t time.Time) string {
	return words.Connection(rc.lang).Day(t.In(rc.zone))
}

// usedBy is the tafs that use g, by their names: its own taf, while g gives
// it tools (connected, paused, signed out, doors on, or tools added to a
// connected grant). A connection with no tools, or ended, is used by none.
func usedBy(g *Grant) []string {
	out := []string{}
	if !inUse(g) {
		return out
	}
	if n := label(g.Request.AgentName); n != "" {
		out = append(out, n)
	}
	return out
}

func inUse(g *Grant) bool {
	switch g.Status {
	case StatusConnected, StatusNeedsReconnect:
		return true
	case StatusListed:
		return len(g.TriggerIDs) > 0
	case StatusPending:
		return len(g.Connected) > 0
	}
	return false
}

// unfinished: a connection signed in, or on its way, that no taf uses yet.
func unfinished(g *Grant) bool {
	return g.Status == StatusListed && len(g.TriggerIDs) == 0
}

func healthOf(g *Grant, rc rowContext) connectionHealth {
	c := words.Connection(rc.lang)
	vendor := vendorName(g)
	switch {
	case g.Status == StatusDisconnected || g.Status == StatusFailed:
		return connectionHealth{State: HealthEnded, Words: statusIn(words.Approve(rc.lang), g)}
	case g.Paused != nil:
		return connectionHealth{State: HealthPaused, Words: statusIn(words.Approve(rc.lang), g)}
	case g.Status == StatusNeedsReconnect:
		if !g.SignedOutAt.IsZero() {
			return connectionHealth{State: HealthSignedOut, Words: fmt.Sprintf(c.ExpiredOnF, vendor, rc.day(g.SignedOutAt)), At: g.SignedOutAt}
		}
		return connectionHealth{State: HealthSignedOut, Words: fmt.Sprintf(c.ExpiredF, vendor)}
	case unfinished(g):
		if t, err := time.Parse(time.RFC3339, g.Request.CreatedAt); err == nil {
			return connectionHealth{State: HealthUnfinished, Words: fmt.Sprintf(c.UnfinishedOnF, rc.day(t))}
		}
		return connectionHealth{State: HealthUnfinished, Words: c.Unfinished}
	case inUse(g):
		switch {
		case !rc.known:
			return connectionHealth{State: HealthWorking, OK: true, Words: c.Working}
		case rc.used.IsZero():
			return connectionHealth{State: HealthWorking, OK: true, Words: c.WorkingNotYet}
		}
		return connectionHealth{State: HealthWorking, OK: true, Words: fmt.Sprintf(c.WorkingLastF, rc.when(rc.used)), At: rc.used}
	}
	// On its way: waiting for an approval, a sign-in or its confirmation.
	return connectionHealth{State: HealthWaiting, Words: statusIn(words.Approve(rc.lang), g)}
}

// signInAgainHow is how a signed-out g signs in again on itself, or "" when
// it can't: a managed app only when its account's name was kept (a new
// sign-in is confirmed only on the same account), a key connection with its
// key, a vendor's own sign-in with a new one; never one that uses another
// connection's sign-in, which is signed in again there.
func signInAgainHow(g *Grant) string {
	if g.Status != StatusNeedsReconnect {
		return ""
	}
	switch {
	case g.Request.Connector != nil:
		if g.Connector != nil && g.Connector.AccountName != "" {
			return "approve"
		}
		return ""
	case g.TokenOf != "":
		return ""
	case g.Request.Kind == KindAPIKey:
		return "key"
	case g.Request.Kind == KindBrowser:
		return ""
	}
	return "approve"
}

func actionsOf(g *Grant, rc rowContext, users []string) []connectionAction {
	c := words.Connection(rc.lang)
	out := []connectionAction{}
	switch {
	case g.Status == StatusDisconnected:
		return out
	case g.Status == StatusFailed || !inUse(g):
		// Unfinished, failed, or on its way: nothing uses it, and removing it
		// is all that can be done here.
		return append(out, connectionAction{ID: ActionRemoveUnfinished, Words: c.RemoveUnfinished, How: "delete", Stops: c.NothingStops})
	}
	if how := signInAgainHow(g); how != "" {
		out = append(out, connectionAction{ID: ActionSignInAgain, Words: c.SignInAgain, How: how})
	}
	stops := c.NothingStops
	if len(users) > 0 {
		lines := make([]string, 0, len(users))
		for _, u := range users {
			lines = append(lines, fmt.Sprintf(c.StopsF, u, vendorName(g)))
		}
		stops = strings.Join(lines, " ")
	}
	return append(out, connectionAction{ID: ActionRemove, Words: c.Remove, How: "delete", Stops: stops})
}

// fillRow sets v's used_by, health and actions for g, as r's app reads them.
// used and known are g's last use and whether it was looked up (the list
// does; a single grant's answer says "Working" with no time).
func (s *Server) fillRow(r *http.Request, g *Grant, v *grantView, used time.Time, known bool) {
	zone, named := ownerZone(r)
	lang := "en"
	if r != nil {
		lang = words.Pick("", r.Header.Get("Accept-Language"))
	}
	rc := rowContext{lang: lang, zone: zone, named: named, now: s.now(), used: used, known: known}
	v.UsedBy = usedBy(g)
	h := healthOf(g, rc)
	v.Health = &h
	v.Actions = actionsOf(g, rc, v.UsedBy)
}
