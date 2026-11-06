package planque

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"github.com/kiff/planque/internal/appconfig"
)

// telegramEntry is how an app knows a Telegram bot grant: it has no catalog
// entry of its own, and planque's description URL isn't a name for people.
const telegramEntry = "api:telegram"

// entryOf is the catalog entry recorded for g's app: the one its request
// names, which the owner signed, else the one planque recorded on the grant
// for a grant made before the request carried it (Grant.Entry).
func entryOf(g *Grant) string {
	if g.Request.CatalogEntry != "" {
		return g.Request.CatalogEntry
	}
	return g.Entry
}

// appEntry is the catalog entry for g's app, as the apps and the catalog
// read it: the one recorded, else the one deriveEntry finds. Empty when
// it's unclear.
func (s *Server) appEntry(g *Grant) string {
	if e := entryOf(g); e != "" {
		return e
	}
	return s.deriveEntry(g)
}

// deriveEntry is the catalog entry for g's app when it's certain:
//   - a managed app's: its provider and toolkit;
//   - the browser's, and a Telegram bot's: theirs;
//   - any other: the one reviewed catalog entry for the same server, reached
//     the same way (an MCP sign-in, or a REST API by key or registered app).
//
// It's empty when no entry matches, or more than one does: never guessed.
func (s *Server) deriveEntry(g *Grant) string {
	r := g.Request
	switch {
	case r.Connector != nil:
		if r.Connector.Provider == "" || r.Connector.Toolkit == "" {
			return ""
		}
		return "connector:" + r.Connector.Provider + ":" + r.Connector.Toolkit
	case r.Kind == KindBrowser:
		return BrowserEntry
	case isTelegram(g):
		return telegramEntry
	}
	family := requestFamily(r)
	if family == "" {
		return ""
	}
	match := ""
	for id, e := range s.CatalogEntries {
		ef, server := s.entryReach(e)
		if ef != family || !sameServer(server, r.Server) {
			continue
		}
		if match != "" {
			return "" // two entries for the same server: unclear
		}
		match = id
	}
	return match
}

// Families of reach: an MCP server's sign-in, or a REST API (by key or by a
// registered app). Two grants of one family to one server are one app.
const (
	familyMCP  = "mcp"
	familyREST = "rest"
)

func requestFamily(r Request) string {
	switch r.Kind {
	case "", appconfig.MCP:
		return familyMCP
	case KindAPIKey, appconfig.REST:
		return familyREST
	}
	return ""
}

// entryReach is how a reviewed catalog entry reaches its vendor: its family
// and its server, as a grant asked from it would record them.
func (s *Server) entryReach(e CatalogEntry) (family, server string) {
	switch {
	case e.Application != "":
		app, ok := s.Applications[e.Application]
		if !ok {
			return "", ""
		}
		switch app.Kind {
		case appconfig.MCP:
			return familyMCP, app.Server
		case appconfig.REST:
			return familyREST, app.Server
		}
		return "", ""
	case e.Kind == KindAPIKey:
		if e.API == nil {
			return "", ""
		}
		return familyREST, e.API.Base
	case e.Kind == "":
		return familyMCP, e.Server
	}
	return "", ""
}

// sameServer: two addresses of one server, whatever the case of the scheme
// and host, the default port, or a trailing slash.
func sameServer(a, b string) bool {
	ua, errA := url.Parse(strings.TrimSpace(a))
	ub, errB := url.Parse(strings.TrimSpace(b))
	if errA != nil || errB != nil || ua.Host == "" || ub.Host == "" || ua.User != nil || ub.User != nil {
		return false
	}
	return strings.EqualFold(ua.Scheme, ub.Scheme) && hostKey(ua) == hostKey(ub) &&
		strings.TrimSuffix(ua.EscapedPath(), "/") == strings.TrimSuffix(ub.EscapedPath(), "/") &&
		ua.RawQuery == ub.RawQuery
}

func hostKey(u *url.URL) string {
	host, port := strings.ToLower(u.Hostname()), u.Port()
	if port == "" || port == "443" && strings.EqualFold(u.Scheme, "https") || port == "80" && strings.EqualFold(u.Scheme, "http") {
		return host
	}
	return host + ":" + port
}

// recordEntry is the backfill for a grant made before its request named its
// app: planque records the entry on the grant itself, outside the request the
// owner signed, which stays byte for byte as signed. It does nothing for a
// grant that has its entry already, or whose entry is unclear. A grant
// someone else saved meanwhile is left for the next listing.
func (s *Server) recordEntry(ctx context.Context, g *Grant) {
	if entryOf(g) != "" {
		return
	}
	entry := s.deriveEntry(g)
	if entry == "" {
		return
	}
	saved := *g
	saved.Entry = entry
	if err := s.saveGrant(ctx, &saved); err != nil {
		s.log().InfoContext(ctx, "entry not recorded yet", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return
	}
	*g = saved
}
