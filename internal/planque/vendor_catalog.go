package planque

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/profile"
)

// CatalogEntry is an operator-installed, reviewed connection definition.
// It holds public protocol facts, never credentials. Revisions are immutable:
// changing operations or access requires a new id and owner approval.
type CatalogEntry struct {
	ID                string          `json:"-"`
	Vendor            string          `json:"vendor"`
	Title             string          `json:"title"`
	Icon              string          `json:"icon,omitempty"`
	Kind              string          `json:"kind,omitempty"`
	Server            string          `json:"server,omitempty"`
	Application       string          `json:"application,omitempty"`
	API               *apiAsk         `json:"api,omitempty"`
	VersionHeader     string          `json:"version_header,omitempty"`
	DescriptionSHA256 string          `json:"description_sha256,omitempty"`
	Tools             []ToolRequest   `json:"tools"`
	Profile           json.RawMessage `json:"profile"`
}

var catalogID = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}@[1-9][0-9]{0,8}$`)

// DecodeCatalog validates the entire configuration before serving it. No
// client request can install a catalog entry or replace its profile.
func DecodeCatalog(raw []byte, apps map[string]appconfig.Application) (map[string]CatalogEntry, error) {
	if len(raw) > openapi.MaxBytes {
		return nil, errors.New("vendor catalog exceeds 16 MB")
	}
	var entries map[string]CatalogEntry
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&entries) != nil || dec.Decode(new(any)) != io.EOF || entries == nil || len(entries) > 500 {
		return nil, errors.New("vendor catalog must be an object of at most 500 documented entries")
	}
	for id, entry := range entries {
		entry.ID = id
		if err := entry.validate(apps); err != nil {
			return nil, fmt.Errorf("vendor catalog %s: %w", id, err)
		}
		entries[id] = entry
	}
	return entries, nil
}

func catalogHTTPS(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && u.RawQuery == "" && u.Fragment == ""
}

func (e CatalogEntry) validate(apps map[string]appconfig.Application) error {
	if !catalogID.MatchString(e.ID) || label(e.Vendor) != e.Vendor || e.Vendor == "" || label(e.Title) != e.Title || e.Title == "" {
		return errors.New("a catalog entry needs a versioned id, vendor and title")
	}
	if e.Icon != "" && !catalogHTTPS(e.Icon) {
		return errors.New("a catalog icon must be a fixed HTTPS URL")
	}
	p, err := profile.Parse(e.Profile)
	if err != nil {
		return err
	}
	if p.Source != "catalog" || p.CatalogEntry != e.ID || p.Vendor != e.Vendor || p.TestProven {
		return errors.New("the reviewed profile must name this catalog entry and vendor; test evidence is recorded separately")
	}
	if len(e.Tools) == 0 {
		return errors.New("a catalog entry must supply its operations before approval")
	}
	raw, _ := json.Marshal(e.Tools)
	var tools []ToolRequest
	_ = json.Unmarshal(raw, &tools)
	if _, err := deriveProfileFields(e.Profile, tools, true); err != nil {
		return err
	}
	if msg := checkTools("catalog", tools); msg != "" {
		return errors.New(msg)
	}
	if msg := oneOfWithoutKey(e.Kind, tools); msg != "" {
		return errors.New(msg)
	}
	server, header := e.Server, e.VersionHeader
	var api *APIRequest
	switch {
	case e.Application != "":
		app, ok := apps[e.Application]
		if !ok || app.Validate() != nil || e.Kind != "" && e.Kind != app.Kind || e.Server != "" && e.Server != app.Server || e.API != nil || e.DescriptionSHA256 != "" {
			return errors.New("a catalog application must use its configured endpoints and operations")
		}
		if e.VersionHeader != "" {
			return errors.New("an application supplies its own reviewed version header")
		}
		server, header = app.Server, app.VersionHeader
		if app.Kind == appconfig.REST {
			api = &APIRequest{}
			if msg := apiToolsCheck(app.Snapshot, tools); msg != "" {
				return errors.New(msg)
			}
		}
	case e.Kind == KindAPIKey:
		if e.API == nil || !catalogHTTPS(e.API.URL) || !catalogHTTPS(e.API.Base) || e.Server != "" && e.Server != e.API.Base {
			return errors.New("a key catalog entry needs fixed API source and base URLs")
		}
		sum, err := hex.DecodeString(e.DescriptionSHA256)
		if err != nil || len(sum) != 32 || e.DescriptionSHA256 != strings.ToLower(e.DescriptionSHA256) {
			return errors.New("a key catalog entry needs the reviewed description's SHA-256")
		}
		if e.API.Auth.In != "header" && e.API.Auth.In != "query" || strings.TrimSpace(e.API.Auth.Name) == "" {
			return errors.New("a key catalog entry must say where the key goes")
		}
		server, api = e.API.Base, &APIRequest{}
	case e.Kind == "":
		if !catalogHTTPS(server) || e.API != nil || e.DescriptionSHA256 != "" {
			return errors.New("an MCP catalog entry needs its fixed server")
		}
	default:
		return errors.New("unsupported catalog connection kind")
	}
	return bindProfileTransport(e.Profile, api, server, header)
}

type catalogEntryView struct {
	Entry  string `json:"entry"`
	Vendor string `json:"vendor"`
	Title  string `json:"title"`
	Icon   string `json:"icon,omitempty"`
	How    string `json:"how"`
	State  string `json:"state"`
	Grant  string `json:"grant,omitempty"`
}

// vendorCatalog combines configured entries with paginated provider discovery.
// Connection state is read only for the authenticated account.
func (s *Server) vendorCatalog(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	states, err := s.entryStates(r.Context(), tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	entries := []catalogEntryView{}
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	for id, e := range s.CatalogEntries {
		if q != "" && !strings.Contains(strings.ToLower(e.Vendor+" "+e.Title), q) {
			continue
		}
		v := catalogEntryView{Entry: id, Vendor: e.Vendor, Title: e.Title, Icon: e.Icon, How: "oauth", State: "available"}
		if e.Kind == KindAPIKey {
			v.How = "key"
		}
		if state, ok := states[id]; ok {
			v.State, v.Grant = state.State, state.Grant
		}
		entries = append(entries, v)
	}
	// The browser (kiff/planque#86): planque's own entry, asked for with
	// POST …/catalog/browser/grants, or POST …/grants with kind browser.
	if s.Browser != nil && (q == "" || strings.Contains(strings.ToLower(browserVendor+" use websites"), q)) {
		v := catalogEntryView{Entry: BrowserEntry, Vendor: browserVendor, Title: "Use websites", How: KindBrowser, State: "available"}
		if state, ok := states[BrowserEntry]; ok {
			v.State, v.Grant = state.State, state.Grant
		}
		entries = append(entries, v)
	}
	// Stripe: planque's own entry, a direct sign-in to Stripe's MCP server
	// with test mode (stripe_refunds.go), listed on the first page.
	if r.URL.Query().Get("cursor") == "" && (q == "" || strings.Contains(strings.ToLower(stripeVendor), q)) {
		v := catalogEntryView{Entry: StripeEntry, Vendor: stripeVendor, Title: stripeVendor, Icon: stripeIcon, How: "oauth", State: "available"}
		if state, ok := states[StripeEntry]; ok {
			v.State, v.Grant = state.State, state.Grant
		}
		entries = append(entries, v)
	}
	// Qonto: planque's own entry, a direct sign-in to Qonto's MCP server
	// with reviewed scopes (qonto.go), listed on the first page.
	if r.URL.Query().Get("cursor") == "" && (q == "" || strings.Contains(strings.ToLower(qontoVendor), q)) {
		v := catalogEntryView{Entry: QontoEntry, Vendor: qontoVendor, Title: qontoVendor, Icon: qontoIcon, How: "oauth", State: "available"}
		if state, ok := states[QontoEntry]; ok {
			v.State, v.Grant = state.State, state.Grant
		}
		entries = append(entries, v)
	}
	providerID := r.URL.Query().Get("provider")
	if providerID == "" {
		ids := make([]string, 0, len(s.Connectors))
		for id := range s.Connectors {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		if len(ids) > 0 {
			providerID = ids[0]
		}
	}
	next := ""
	skipped := 0
	partial := false
	if providerID != "" {
		provider := s.Connectors[providerID]
		if provider == nil {
			writeErr(w, http.StatusBadRequest, "no such catalog provider")
			return
		}
		page, err := provider.Catalog(r.Context(), r.URL.Query().Get("q"), r.URL.Query().Get("cursor"))
		if err != nil {
			writeErr(w, http.StatusServiceUnavailable, "the connection catalog is unavailable")
			return
		}
		next = page.Next
		skipped = page.Skipped
		partial = page.Partial
		for _, app := range page.Items {
			id := "connector:" + providerID + ":" + app.Slug
			// The managed connector's Stripe and Qonto aren't offered: each is
			// a direct entry above, Stripe's the one with test mode, Qonto's
			// the one with reviewed scopes. A connection made to one before
			// stays listed, so it can be disconnected.
			if _, made := states[id]; directEntryFor(app.Slug) != "" && !made {
				continue
			}
			v := catalogEntryView{Entry: id, Vendor: app.Slug, Title: label(app.Name), Icon: app.Icon, How: "oauth", State: "available"}
			// An app that needs an operator's own credentials first is listed,
			// said as such, and not started (kiff/planque#87).
			if app.Setup == connector.SetupNeedsSetup {
				v.State = "needs_setup"
			}
			if state, ok := states[id]; ok {
				v.State, v.Grant = state.State, state.Grant
			}
			entries = append(entries, v)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Entry < entries[j].Entry })
	writeJSON(w, http.StatusOK, map[string]any{"entries": entries, "by_key": true, "provider": providerID, "next_cursor": next, "skipped": skipped, "partial": partial})
}

// entryStates is each catalog entry the tenant has a connection to, with its
// state (connected, or needs_reconnect) and the grant that says so: the
// catalog and the tool search mark an app the same way.
func (s *Server) entryStates(ctx context.Context, tenant string) (map[string]catalogEntryView, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	states := map[string]catalogEntryView{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.Tenant != tenant {
			continue
		}
		// Any connection to the app counts, a Coach-started one included.
		entry := s.appEntry(&g)
		if entry == "" {
			continue
		}
		state := ""
		switch {
		case len(g.TriggerIDs) > 0 && g.SignedIn && g.Status != StatusDisconnected && g.Status != StatusFailed && g.Status != StatusNeedsReconnect, g.Status == StatusConnected || g.Status == StatusPending && len(g.Connected) > 0:
			state = "connected"
		case g.Status == StatusNeedsReconnect:
			state = "needs_reconnect"
		}
		if state == "" {
			continue
		}
		previous := states[entry]
		if previous.State == "connected" || previous.State == state && previous.Grant < g.ID {
			continue
		}
		states[entry] = catalogEntryView{State: state, Grant: g.ID}
	}
	return states, nil
}

func (s *Server) requestCatalogGrant(w http.ResponseWriter, r *http.Request) {
	// Authenticate before revealing whether the entry exists.
	if _, err := s.agentKey(r, r.PathValue("tenant")); err != nil {
		s.fail(w, r, err)
		return
	}
	id := r.PathValue("entry")
	if id == StripeEntry {
		// A direct sign-in to Stripe's MCP server, connected first: its
		// tools are chosen from what the server lists.
		s.requestGrantFrom(w, r, &CatalogEntry{ID: StripeEntry, Vendor: stripeVendor, Title: stripeVendor, Server: stripeServer, Tools: []ToolRequest{}})
		return
	}
	if id == QontoEntry {
		// A direct sign-in to Qonto's MCP server, connected first: its tools
		// are chosen from what the server lists.
		s.requestGrantFrom(w, r, &CatalogEntry{ID: QontoEntry, Vendor: qontoVendor, Title: qontoVendor, Server: qontoServer, Tools: []ToolRequest{}})
		return
	}
	if id == BrowserEntry {
		if s.Browser == nil {
			writeErr(w, http.StatusServiceUnavailable, "this installation has no browser")
			return
		}
		s.requestGrantFrom(w, r, &CatalogEntry{ID: BrowserEntry, Vendor: browserVendor, Title: "Use websites", Kind: KindBrowser})
		return
	}
	entry, ok := s.CatalogEntries[id]
	if !ok {
		if provider, toolkit, match := connectorEntry(id); match {
			switch directEntryFor(toolkit) {
			case StripeEntry:
				writeErr(w, http.StatusConflict, "Stripe connects directly, with test mode: ask for catalog entry "+StripeEntry)
				return
			case QontoEntry:
				writeErr(w, http.StatusConflict, "Qonto connects directly, with reviewed scopes: ask for catalog entry "+QontoEntry)
				return
			}
			s.requestConnectorGrant(w, r, provider, toolkit)
			return
		}
		writeErr(w, http.StatusNotFound, "no such catalog entry")
		return
	}
	entry.ID = id
	if err := entry.validate(s.Applications); err != nil {
		writeErr(w, http.StatusServiceUnavailable, "that catalog entry is not configured for this installation")
		return
	}
	// Grant naming and profile derivation mutate their inputs. Freeze a new
	// copy so requests and concurrent accounts never mutate installed facts.
	raw, err := json.Marshal(entry)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var frozen CatalogEntry
	if err := json.Unmarshal(raw, &frozen); err != nil {
		s.fail(w, r, err)
		return
	}
	frozen.ID = id
	s.requestGrantFrom(w, r, &frozen)
}
