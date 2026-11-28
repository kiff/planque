package planque

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/store"
)

// The sign-in health check (browser batch 2): to warn the owner before a
// saved sign-in in the hosted browser stops working, planque opens one fixed
// page of the site with it, read-only (browser.HealthCheck), at most once a
// day per connection and site, on demand only, and keeps only the verdict,
// when it checked, and the sign-in's own end as the site set it. A signed-out
// verdict is for the owner's own step ("Your turn"); planque never signs in.

//go:embed health_pages.json
var healthPagesJSON []byte

// HealthPage is one site's reviewed fixed page: its https URL, and words
// only a signed-in visitor sees on it, the proof a check reads signed in by.
type HealthPage struct {
	URL   string `json:"url"`
	Shows string `json:"shows"`
}

// ReviewedHealthPages is the reviewed list of each site's fixed page
// (health_pages.json): a site not on it has no check.
func ReviewedHealthPages() map[string]HealthPage {
	var f struct {
		Pages map[string]HealthPage `json:"pages"`
	}
	if err := json.Unmarshal(healthPagesJSON, &f); err != nil {
		panic("health_pages.json: " + err.Error())
	}
	return f.Pages
}

const (
	// healthEvery: a site's check runs at most once in this, per connection.
	healthEvery = 24 * time.Hour
	// healthBackoff: after a check that couldn't tell, none before this.
	healthBackoff = 6 * time.Hour
	// healthLease bounds the browser session slot a check holds.
	healthLease = time.Minute
)

func skHealth(grant, site string) string { return "HEALTH#" + grant + "#" + site }

// healthRecord is all a check keeps: the verdict and when, and after one
// that couldn't tell, when it may run again. Never anything of the page.
type healthRecord struct {
	Verdict   string    `json:"verdict"`
	CheckedAt time.Time `json:"checked_at"`
	RetryAt   time.Time `json:"retry_at,omitzero"`
	// StateVer is the saved sign-in's version the verdict is about: a new
	// sign-in (the owner's, after a signed-out verdict) makes a check due
	// again at once, whatever the day's limit.
	StateVer int64 `json:"state_ver,omitempty"`
}

// healthSite is one site's answer: its verdict and when it was checked,
// and when its saved sign-in ends as the site set it (the latest expiry of
// the site's cookies; none when they're all session cookies). Why, when
// there's no check to give.
type healthSite struct {
	Site      string     `json:"site"`
	Verdict   string     `json:"verdict"`
	CheckedAt *time.Time `json:"checked_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Why       string     `json:"why,omitempty"`
}

// grantHealth answers POST (run what's due) and GET (read what's kept)
// /v1/accounts/{tenant}/grants/{id}/health, with the account's key, for a
// hosted-browser connection on an account the hosted browser is offered to.
func (s *Server) grantHealth(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	switch {
	case g.Request.Kind != KindBrowser || chromeGrant(g) || g.Request.Browser == nil:
		writeErr(w, http.StatusBadRequest, "only a hosted-browser connection has sign-in checks")
		return
	case !s.browserOnFor(tenant):
		writeErr(w, http.StatusForbidden, "planque's hosted browser isn't offered to this account")
		return
	}
	run := r.Method == http.MethodPost
	// A connection's sign-in health is for the taf the connection is for, or
	// the agent that asked for it (Setup, which tafs cloud's reads go
	// through): another taf's key of the account neither runs a check (it
	// takes the connection's browser slot and may reseal its sign-in) nor
	// reads whether that sign-in holds.
	if me.AgentID != g.Request.Agent && me.AgentID != g.Request.AskedBy {
		writeErr(w, http.StatusForbidden, "only the taf this connection is for, or the agent that asked for it, reads or runs its sign-in check")
		return
	}
	out := make([]healthSite, 0, len(g.Request.Browser.Sites))
	for _, site := range g.Request.Browser.Sites {
		out = append(out, s.siteHealth(r.Context(), g, site, run))
	}
	writeJSON(w, http.StatusOK, map[string]any{"grant": g.ID, "sites": out})
}

// siteHealth is one site's answer: the kept verdict, or, when run and one is
// due, a new check.
func (s *Server) siteHealth(ctx context.Context, g *Grant, site string, run bool) healthSite {
	out := healthSite{Site: site, Verdict: browser.HealthUnknown}
	page := s.HealthPages[site]
	if page.URL == "" || strings.TrimSpace(page.Shows) == "" {
		out.Why = "no reviewed page for this site"
		return out
	}
	var rec healthRecord
	_, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skHealth(g.ID, site), &rec)
	if err != nil && !errors.Is(err, errNotFound) {
		out.Why = "the kept check couldn't be read"
		return out
	}
	now := s.now()
	kept := err == nil
	state, from, serr := s.openBrowserState(ctx, g)
	due := !kept || (rec.RetryAt.IsZero() && !now.Before(rec.CheckedAt.Add(healthEvery))) || (!rec.RetryAt.IsZero() && !now.Before(rec.RetryAt)) ||
		(serr == nil && from != rec.StateVer)
	if run && due && serr == nil {
		rec = s.runHealth(ctx, g, site, page, state, from)
		kept = true
		if st, _, err := s.openBrowserState(ctx, g); err == nil {
			state = st
		}
	}
	if kept {
		out.Verdict, out.CheckedAt = rec.Verdict, &rec.CheckedAt
	}
	if exp, ok := cookieEnd(state, site); ok {
		out.ExpiresAt = &exp
	}
	return out
}

// runHealth runs one check and keeps its verdict. With no saved sign-in
// there's nothing to open: signed out. A check that can't run or can't tell
// keeps "unknown" and waits healthBackoff. A signed-in check's cookies are
// sealed as a session's are, only over the state it started from.
func (s *Server) runHealth(ctx context.Context, g *Grant, site string, page HealthPage, state []byte, from int64) healthRecord {
	now := s.now().UTC()
	rec := healthRecord{Verdict: browser.HealthUnknown, CheckedAt: now}
	checker, ok := s.Browser.(browser.HealthChecker)
	switch {
	case len(state) == 0:
		rec.Verdict = browser.HealthSignedOut
	case !ok:
		rec.RetryAt = now.Add(healthBackoff)
	default:
		lease, _, err := s.takeBrowserLease(ctx, g.Tenant, now.Add(healthLease))
		if err != nil {
			// The account's sessions are busy, or unreadable: nothing kept, so
			// the next ask runs it.
			return healthRecord{Verdict: browser.HealthUnknown, CheckedAt: now}
		}
		res, err := checker.Health(ctx, browser.HealthCheck{Grant: g.ID, State: state, Site: site, URL: page.URL, Shows: page.Shows})
		s.dropBrowserLease(context.WithoutCancel(ctx), lease)
		switch {
		case err != nil:
			rec.RetryAt = now.Add(healthBackoff)
		case res.Verdict == browser.HealthSignedIn || res.Verdict == browser.HealthSignedOut:
			rec.Verdict = res.Verdict
		default:
			rec.RetryAt = now.Add(healthBackoff)
		}
		if err == nil && res.Verdict == browser.HealthSignedIn && len(res.State) > 0 {
			s.keepBrowserState(ctx, g, res.State, &from)
			// The verdict is about the sign-in as resealed, so the reseal
			// doesn't make the next ask due again.
			if _, ver, err := s.openBrowserState(ctx, g); err == nil {
				from = ver
			}
		}
		if err == nil && res.Blocked > 0 {
			s.log().InfoContext(ctx, "sign-in check failed requests that would write", slog.String("grant", g.ID), slog.String("site", site), slog.Int("blocked", res.Blocked))
		}
	}
	rec.StateVer = from
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skHealth(g.ID, site), TTL: now.Add(healthEvery + healthBackoff + time.Hour).Unix()}, rec, store.Cond{}); err != nil {
		s.log().WarnContext(ctx, "sign-in check not kept", slog.String("grant", g.ID), slog.String("site", site), slog.String("err", err.Error()))
	}
	s.log().InfoContext(ctx, "sign-in checked", slog.String("grant", g.ID), slog.String("site", site), slog.String("page_host", hostOfURL(page.URL)), slog.String("verdict", rec.Verdict))
	return rec
}

// cookieEnd is when site's saved sign-in ends as the site set it: the
// latest expiry among the state's cookies for that site. None when they're
// all session cookies, or there are none. Never a cookie's name or value.
func cookieEnd(state []byte, site string) (time.Time, bool) {
	var st struct {
		Cookies []struct {
			Domain  string  `json:"domain"`
			Expires float64 `json:"expires"`
		} `json:"cookies"`
	}
	if len(state) == 0 || json.Unmarshal(state, &st) != nil {
		return time.Time{}, false
	}
	site = strings.ToLower(site)
	var end float64
	for _, c := range st.Cookies {
		d := strings.ToLower(strings.TrimPrefix(c.Domain, "."))
		if c.Expires <= 0 || d == "" || (site != d && !strings.HasSuffix(site, "."+d)) {
			continue
		}
		end = max(end, c.Expires)
	}
	if end == 0 {
		return time.Time{}, false
	}
	return time.Unix(int64(end), 0).UTC(), true
}

func hostOfURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
