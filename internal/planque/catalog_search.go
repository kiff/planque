package planque

import (
	"net/http"
	"sort"
	"strings"

	"github.com/kiff/planque/internal/connector"
)

// A tool found by words, with its app as it stands (kiff/planque#87).
type foundTool struct {
	Entry       string `json:"entry"`
	App         string `json:"app"`
	AppTitle    string `json:"app_title"`
	AppState    string `json:"app_state"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ReadOnly    bool   `json:"read_only"`
}

// appState is a catalog entry's state for an app as the provider lists it:
// what the owner can do with it before any grant exists.
func appState(app connector.App) string {
	if app.Setup == connector.SetupNeedsSetup {
		return "needs_setup"
	}
	return "available"
}

// searchTools finds tools across the provider's apps by words, for Coach to
// propose apps before connecting: GET /v1/accounts/{tenant}/catalog/tools
// ?q=…&apps=slug,slug&cursor=…. Only the provider's catalog is searched;
// configured reviewed entries are found by name through the catalog read.
// Platform meta tools never come back.
func (s *Server) searchTools(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" || len(q) > 200 {
		writeErr(w, http.StatusBadRequest, "q is the words to search tools by, at most 200 characters")
		return
	}
	providerID := r.URL.Query().Get("provider")
	if providerID == "" {
		providerID = s.defaultProvider()
	}
	provider := s.Connectors[providerID]
	if provider == nil {
		writeErr(w, http.StatusBadRequest, "no such catalog provider")
		return
	}
	var apps []string
	for _, a := range strings.Split(r.URL.Query().Get("apps"), ",") {
		if a = strings.TrimSpace(a); a != "" {
			apps = append(apps, a)
		}
	}
	found, err := provider.Search(r.Context(), q, apps, r.URL.Query().Get("cursor"))
	if err != nil {
		if err == connector.ErrDefinition {
			writeErr(w, http.StatusBadRequest, "apps names at most 20 apps, as the catalog names them")
			return
		}
		writeErr(w, http.StatusServiceUnavailable, "the connection catalog is unavailable")
		return
	}
	tools := []foundTool{}
	for _, t := range found.Tools {
		app, ok := found.Apps[t.Toolkit]
		// The managed connector's Stripe and Qonto aren't offered: each is
		// planque's direct entry.
		if !ok || directEntryFor(t.Toolkit) != "" {
			continue
		}
		tools = append(tools, foundTool{Entry: "connector:" + providerID + ":" + t.Toolkit, App: t.Toolkit, AppTitle: label(app.Name), AppState: appState(app),
			Slug: t.Slug, Name: label(t.Name), Description: t.Description, ReadOnly: t.ReadOnly})
	}
	writeJSON(w, http.StatusOK, map[string]any{"tools": tools, "provider": providerID, "next_cursor": found.Next, "skipped": found.Skipped})
}

// An app's event, for a listing (kiff/planque#87): what Coach reads to pick
// one as it picks an action. The config is the filters' schema the owner's
// approval names exactly.
type eventView struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
	Config      any    `json:"config"`
	// Door: the event is one planque serves as a door (Slack messages),
	// approved by its channels; any other wakes a taf as an app event.
	Door bool `json:"door"`
}

func eventViews(page connector.Page[connector.TriggerType]) []eventView {
	views := []eventView{}
	for _, t := range page.Items {
		views = append(views, eventView{Slug: t.Slug, Name: label(t.Name), Description: t.Description, Version: t.Version, Config: t.Config, Door: supportedTrigger(t)})
	}
	return views
}

// catalogEvents lists an app's events before any connection exists:
// GET /v1/accounts/{tenant}/catalog/{entry}/events?cursor=….
func (s *Server) catalogEvents(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	providerID, toolkit, ok := connectorEntry(r.PathValue("entry"))
	if !ok {
		writeErr(w, http.StatusNotFound, "only a connection provider's app has events to list")
		return
	}
	provider := s.Connectors[providerID]
	if provider == nil {
		writeErr(w, http.StatusNotFound, "no such catalog provider")
		return
	}
	events, ok := provider.(connector.EventProvider)
	if !ok {
		writeErr(w, http.StatusServiceUnavailable, "this provider does not support managed events")
		return
	}
	s.writeEvents(w, r, events, toolkit)
}

// grantEvents lists the events of a connection's app, for the agent that
// asked for it: GET /v1/accounts/{tenant}/grants/{id}/events?cursor=….
func (s *Server) grantEvents(w http.ResponseWriter, r *http.Request) {
	g, ok := s.agentsGrant(w, r)
	if !ok {
		return
	}
	if g.Request.Connector == nil {
		writeErr(w, http.StatusBadRequest, "this connection has no app events")
		return
	}
	p, err := s.eventProvider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.writeEvents(w, r, p, g.Request.Connector.Toolkit)
}

func (s *Server) writeEvents(w http.ResponseWriter, r *http.Request, p connector.EventProvider, toolkit string) {
	// The platform's own toolkit has no app behind it: never an app's events.
	if strings.EqualFold(toolkit, "composio") {
		writeErr(w, http.StatusBadRequest, "that isn't an app the catalog names")
		return
	}
	page, err := p.TriggerTypes(r.Context(), toolkit, r.URL.Query().Get("cursor"))
	if err != nil {
		if err == connector.ErrDefinition {
			writeErr(w, http.StatusBadRequest, "that isn't an app the catalog names")
			return
		}
		writeErr(w, http.StatusServiceUnavailable, "the connection catalog is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": eventViews(page), "next_cursor": page.Next, "skipped": page.Skipped})
}

// defaultProvider is the catalog provider a read names when it names none:
// the first by id, as the catalog read picks it.
func (s *Server) defaultProvider() string {
	ids := make([]string, 0, len(s.Connectors))
	for id := range s.Connectors {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}
