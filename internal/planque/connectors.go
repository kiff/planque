package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

// ConnectorRequest is part of the device-signed grant. Once tools are chosen,
// the exact account and concrete versioned definitions are covered by its digest.
type ConnectorRequest struct {
	Provider string `json:"provider"`
	Toolkit  string `json:"toolkit"`
	// Title is the app's name for people ("Gmail"), as the provider lists it:
	// the words the owner reads and signs. Empty on grants made before it.
	Title string `json:"title,omitempty"`
	// AuthConfig and AuthScheme freeze the hosted sign-in setup and wording.
	// Omitted on older grants so their existing digests and words stay valid.
	AuthConfig string           `json:"auth_config,omitempty"`
	AuthScheme string           `json:"auth_scheme,omitempty"`
	Account    string           `json:"account,omitempty"`
	Tools      []connector.Tool `json:"tools,omitempty"`
}

// ConnectorState is private server state, never an agent-selected identity.
// Vendor credentials remain at the configured connection provider.
type ConnectorState struct {
	User       string                `json:"user"`
	Session    string                `json:"session"`
	Account    string                `json:"account"`
	Tools      []connector.Tool      `json:"tools,omitempty"`
	Revocation *connector.Revocation `json:"revocation,omitempty"`
	// AccountName is the vendor account's name as the provider shows it.
	AccountName string `json:"account_name,omitempty"`
	// Mode is the vendor account's mode, "live" or "test", as the objects
	// permitted calls returned said it (their "livemode"), and when: live
	// wins, so once any live object is seen it stays live. Empty is
	// unknown, which every reader treats as live: planque never calls the
	// vendor outside a permit to find out.
	Mode       string    `json:"mode,omitempty"`
	ModeSeenAt time.Time `json:"mode_seen_at,omitzero"`
}

// ModeUnknown is what a grant says before a permitted call has shown its
// account's mode; it is treated as live.
const ModeUnknown = "unknown"

// modeWords say a grant's mode to the owner, the same in every app.
var modeWords = map[string]string{
	"live":      "Live",
	"test":      "Test",
	ModeUnknown: "Mode not known yet: treated as live",
}

// modeOf reads the "livemode" a vendor's answer carries, at its top and in
// the objects it lists, at most three levels down: "live" if any object
// read is live, "test" only if every one read is not, "" when none says (an
// empty list says nothing). Live wins, so a test object never hides a live
// one.
func modeOf(data json.RawMessage) string {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return ""
	}
	live, test := false, false
	var walk func(any, int)
	walk = func(v any, depth int) {
		if depth > 3 || live {
			return
		}
		switch x := v.(type) {
		case map[string]any:
			if b, ok := x["livemode"].(bool); ok {
				if b {
					live = true
					return
				}
				test = true
			}
			for _, k := range []string{"data", "response_data", "object", "result"} {
				walk(x[k], depth+1)
			}
		case []any:
			for _, e := range x {
				walk(e, depth+1)
			}
		}
	}
	walk(v, 0)
	switch {
	case live:
		return "live"
	case test:
		return "test"
	}
	return ""
}

// noteAccount keeps on the grant what a call showed of its account: the
// provider's name for it, and its mode, where live wins: a live answer sets
// live whatever was there, a test answer counts only while the mode is
// unknown, so a test reading never hides a live one. Best effort: a grant
// changed meanwhile keeps its version, and the next call notes it again.
func (s *Server) noteAccount(ctx context.Context, g *Grant, r connector.Result) {
	if g.Connector == nil {
		return
	}
	mode := ""
	if !r.Refused && !r.Unknown {
		switch seen := modeOf(r.Data); {
		case seen == "live" && g.Connector.Mode != "live":
			mode = "live"
		case seen == "test" && g.Connector.Mode == "":
			mode = "test"
		}
	}
	name := r.AccountName
	if mode == "" && (name == "" || name == g.Connector.AccountName) {
		return
	}
	if mode != "" {
		g.Connector.Mode, g.Connector.ModeSeenAt = mode, s.now()
	}
	if name != "" {
		g.Connector.AccountName = name
	}
	_ = s.saveGrant(ctx, g)
}

func connectorEntry(id string) (string, string, bool) {
	p := strings.Split(id, ":")
	if len(p) != 3 || p[0] != "connector" || p[1] == "" || p[2] == "" {
		return "", "", false
	}
	return p[1], p[2], true
}

// onDemandADay is the most apps an account sets up on demand in one day
// (kiff/planque#87): each is a sign-in the provider sets up for everyone,
// counted against its allowance, so a runaway stops here.
const onDemandADay = 10

var errSettingUp = errors.New("the app's sign-in is being set up")
var errAllowance = errors.New("the account's on-demand set-ups for the day are used")

// setUpRecord is the sign-in set up on demand for an app, one per provider
// and app: the lock two owners picking the app at once take once, what it
// made, and which account's pick made it. An unfinished one, older than a
// minute, is taken over.
type setUpRecord struct {
	Toolkit    string `json:"toolkit"`
	Tenant     string `json:"tenant"`
	AuthConfig string `json:"auth_config,omitempty"`
	AuthScheme string `json:"auth_scheme,omitempty"`
	StartedAt  string `json:"started_at"`
	MadeAt     string `json:"made_at,omitempty"`
}

func setUpPK(provider string) string { return "CONNECTOR#" + provider + "#SETUP" }

// setUpOf is the sign-in set up on demand for the app, if one was.
func (s *Server) setUpOf(ctx context.Context, providerID, toolkit string) (setUpRecord, bool, error) {
	it, ok, err := s.KV.Get(ctx, setUpPK(providerID), toolkit)
	if err != nil || !ok {
		return setUpRecord{}, false, err
	}
	var held setUpRecord
	if json.Unmarshal(it.Doc, &held) != nil {
		return setUpRecord{}, false, nil
	}
	return held, true, nil
}

// allowance is an account's on-demand set-ups so far one day, by the day
// in UTC: one record, bumped under its version before a set-up starts, so
// picks arriving together never pass the day's allowance between them.
type allowance struct {
	Count int `json:"count"`
}

func allowanceSK(day string) string { return "ONDEMAND#" + day }

// takeAllowance takes one of the day's set-ups for the account, or says the
// allowance is used (errAllowance). Contended five times, it gives up as a
// set-up in progress.
func (s *Server) takeAllowance(ctx context.Context, tenant string) error {
	return s.adjustAllowance(ctx, tenant, 1)
}

// giveBackAllowance returns a set-up that didn't happen.
func (s *Server) giveBackAllowance(ctx context.Context, tenant string) {
	_ = s.adjustAllowance(ctx, tenant, -1)
}

func (s *Server) adjustAllowance(ctx context.Context, tenant string, by int) error {
	now := s.now().UTC()
	key := store.Item{PK: tenantPK(tenant), SK: allowanceSK(now.Format("2006-01-02")), TTL: now.Add(48 * time.Hour).Unix()}
	for range 5 {
		it, ok, err := s.KV.Get(ctx, key.PK, key.SK)
		if err != nil {
			return err
		}
		var a allowance
		if ok && json.Unmarshal(it.Doc, &a) != nil {
			a = allowance{}
		}
		if by > 0 && a.Count >= onDemandADay {
			return errAllowance
		}
		a.Count = max(a.Count+by, 0)
		cond := store.Cond{NotExists: true}
		if ok {
			ver := it.Ver
			cond = store.Cond{Ver: &ver}
			key.Ver = it.Ver + 1
		}
		err = putDoc(ctx, s.KV, key, a, cond)
		if errors.Is(err, store.ErrCondition) {
			continue
		}
		return err
	}
	return errSettingUp
}

// setUpApp sets up the sign-in for an on-demand app once across every
// instance. A record that holds the sign-in made answers from it, so the
// provider is asked once and the day's allowance isn't touched; otherwise
// one of the allowance is taken first, then the record under its key is the
// lock: one still being made says so, unless it is stale, which is taken
// over. A set-up that fails gives the lock and the allowance back.
func (s *Server) setUpApp(ctx context.Context, p connector.Provider, providerID, tenant string, app connector.App) (connector.App, error) {
	if held, ok, err := s.setUpOf(ctx, providerID, app.Slug); err != nil {
		return app, err
	} else if ok && held.AuthConfig != "" {
		app.AuthConfig, app.AuthScheme, app.Setup = held.AuthConfig, held.AuthScheme, connector.SetupReady
		return app, nil
	}
	if err := s.takeAllowance(ctx, tenant); err != nil {
		return app, err
	}
	key := store.Item{PK: setUpPK(providerID), SK: app.Slug}
	now := s.now().UTC()
	record := setUpRecord{Toolkit: app.Slug, Tenant: tenant, StartedAt: now.Format(time.RFC3339)}
	err := putDoc(ctx, s.KV, key, record, store.Cond{NotExists: true})
	if errors.Is(err, store.ErrCondition) {
		it, ok, getErr := s.KV.Get(ctx, key.PK, key.SK)
		if getErr != nil {
			s.giveBackAllowance(ctx, tenant)
			return app, getErr
		}
		var held setUpRecord
		if ok && json.Unmarshal(it.Doc, &held) == nil && held.AuthConfig != "" {
			// Made meanwhile by another pick: taken as it is, and no set-up of ours.
			s.giveBackAllowance(ctx, tenant)
			app.AuthConfig, app.AuthScheme, app.Setup = held.AuthConfig, held.AuthScheme, connector.SetupReady
			return app, nil
		}
		started, _ := time.Parse(time.RFC3339, held.StartedAt)
		if ok && now.Sub(started) < time.Minute {
			s.giveBackAllowance(ctx, tenant)
			return app, errSettingUp
		}
		// Stale: a set-up that never finished. Take it over at its version.
		ver := it.Ver
		key.Ver = it.Ver + 1
		if err := putDoc(ctx, s.KV, key, record, store.Cond{Ver: &ver}); err != nil {
			s.giveBackAllowance(ctx, tenant)
			return app, errSettingUp
		}
	} else if err != nil {
		s.giveBackAllowance(ctx, tenant)
		return app, err
	}
	made, err := p.SetUp(ctx, app.Slug)
	if err != nil {
		_, _, _ = s.KV.Delete(ctx, key.PK, key.SK)
		s.giveBackAllowance(ctx, tenant)
		return app, err
	}
	record.AuthConfig, record.AuthScheme, record.MadeAt = made.AuthConfig, made.AuthScheme, s.now().UTC().Format(time.RFC3339)
	if err := putDoc(ctx, s.KV, key, record, store.Cond{}); err != nil {
		return app, err
	}
	return made, nil
}

func (s *Server) connectorIdentity(g *Grant) connector.Identity {
	sum := sha256.Sum256([]byte(s.BaseURL + "\n" + g.Tenant))
	return connector.Identity{User: "planque_" + hex.EncodeToString(sum[:]), Toolkit: g.Request.Connector.Toolkit, AuthConfig: g.Request.Connector.AuthConfig}
}
func (s *Server) provider(g *Grant) (connector.Provider, error) {
	if g.Request.Connector == nil {
		return nil, connector.ErrDefinition
	}
	p := s.Connectors[g.Request.Connector.Provider]
	if p == nil || p.Platform().URL != g.Request.Server {
		return nil, connector.ErrUnavailable
	}
	return p, nil
}
func (s *Server) checkConnector(ctx context.Context, g *Grant) error {
	p, err := s.provider(g)
	if err != nil {
		return err
	}
	if g.Connector == nil || g.Connector.Revocation != nil && g.Connector.Revocation.Disabled {
		return connector.ErrAccount
	}
	id := s.connectorIdentity(g)
	if id.User != g.Connector.User || g.Connector.Account == "" || g.Request.Connector.Account != "" && g.Request.Connector.Account != g.Connector.Account {
		return connector.ErrAccount
	}
	a, err := p.Account(ctx, g.Connector.Account)
	if err != nil {
		return err
	}
	if a.ID != g.Connector.Account || !a.Matches(id) {
		return connector.ErrAccount
	}
	return nil
}

func (s *Server) requestConnectorGrant(w http.ResponseWriter, r *http.Request, providerID, toolkit string) {
	tenant := r.PathValue("tenant")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if _, ok, err := s.KV.Get(r.Context(), tenantPK(tenant), skEnrolment); err != nil {
		s.fail(w, r, err)
		return
	} else if !ok {
		writeErr(w, http.StatusConflict, "the owner's phone is not enrolled with planque for this account yet")
		return
	}
	p := s.Connectors[providerID]
	if p == nil {
		writeErr(w, http.StatusNotFound, "no such catalog provider")
		return
	}
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var ask struct {
		Agent string `json:"agent"`
		// AgentName is the taf's name for people ("Refunds"), optional: the
		// grant's words name the taf by it instead of the agent id.
		AgentName string `json:"agent_name"`
	}
	if err := decodeStrict(body, &ask); err != nil {
		s.fail(w, r, err)
		return
	}
	if msg := checkTools(ask.Agent, nil); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	if ask.AgentName != "" && label(ask.AgentName) != ask.AgentName {
		writeErr(w, http.StatusBadRequest, "agent_name is the taf's name, on one line, at most "+strconv.Itoa(MaxLabel)+" characters")
		return
	}
	app, err := p.App(r.Context(), toolkit)
	if errors.Is(err, connector.ErrNotSetUp) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	platform := p.Platform()
	if app.Slug != toolkit || !catalogHTTPS(platform.URL) || label(platform.Name) == "" {
		writeErr(w, http.StatusServiceUnavailable, "the catalog provider is not configured")
		return
	}
	// An app the provider can sign people in to by itself is set up on the
	// owner's pick, once; one that needs our own developer credentials waits
	// for an operator (kiff/planque#87).
	onDemand := false
	switch app.Setup {
	case connector.SetupNeedsSetup:
		writeErr(w, http.StatusConflict, "This app needs setting up first. Ask tafs to enable it.")
		return
	case connector.SetupOnDemand:
		if app, err = s.setUpApp(r.Context(), p, providerID, tenant, app); err != nil {
			switch {
			case errors.Is(err, errAllowance):
				writeErr(w, http.StatusTooManyRequests, "You've set up "+strconv.Itoa(onDemandADay)+" new apps today. Try again tomorrow.")
			case errors.Is(err, errSettingUp):
				writeErr(w, http.StatusServiceUnavailable, "This app is being set up. Try again in a moment.")
			default:
				writeErr(w, http.StatusServiceUnavailable, "This app can't be set up right now. Try again later.")
			}
			return
		}
		onDemand = true
	case connector.SetupReady:
		// A sign-in set up on an owner's pick, whenever: every grant on it is on demand.
		held, ok, err := s.setUpOf(r.Context(), providerID, toolkit)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		onDemand = ok && held.AuthConfig != "" && held.AuthConfig == app.AuthConfig
	}
	id := "grt_" + oauth.Random(12)
	g := &Grant{ID: id, Tenant: tenant, Status: StatusPending, VendorName: label(app.Name), OnDemand: onDemand,
		Request: Request{Version: "planque.grant.v1", Grant: id, Tenant: tenant, Server: platform.URL, SignInAt: platform.URL,
			Agent: ask.Agent, AgentName: ask.AgentName, Tools: []ToolRequest{}, AskedBy: me.AgentID, CreatedAt: s.now().UTC().Format(time.RFC3339), Kind: "connector",
			CatalogEntry: r.PathValue("entry"), Connector: &ConnectorRequest{Provider: providerID, Toolkit: toolkit, Title: label(app.Name), AuthConfig: app.AuthConfig, AuthScheme: app.AuthScheme},
			API: &APIRequest{URL: platform.URL, Base: platform.URL, SourceKind: "connector"}},
		Audience: s.BaseURL + "/mcp/" + url.PathEscape(tenant) + "/" + id}
	if g.Request.AskedBy == "" {
		g.Request.AskedBy = askedByAccountKey
	}
	g.Text = grantText(g.Request)
	g.Digest, err = digest(g.Request)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g.event(s.now(), "connection requested")
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"grant": s.view(g), "needs_key": false})
}

func (s *Server) approveConnector(w http.ResponseWriter, r *http.Request, sr signedRequest, g *Grant) {
	p, err := s.provider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.SignedIn && g.Status == StatusPending {
		if err := s.checkConnector(r.Context(), g); err != nil {
			s.fail(w, r, err)
			return
		}
		if g.Request.Connector.Account != g.Connector.Account || len(g.Request.Connector.Tools) == 0 {
			writeErr(w, http.StatusConflict, "choose the connection's concrete tools before approving")
			return
		}
		if g.Connected != nil {
			err = s.addTools(r.Context(), g, "", sr.approval)
		} else {
			a := sr.approval
			g.Approval = &a
			g.event(s.now(), "tools approved "+a.Who())
			err = s.connect(r.Context(), g, "")
		}
		if err != nil {
			s.fail(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
	// A replacement account always needs a fresh grant and tool approval.
	// Never swap accounts underneath an already approved profile.
	if g.Connector != nil {
		writeErr(w, http.StatusConflict, "this connection already has a sign-in; finish it or disconnect and request a new grant")
		return
	}
	if s.AppReturn == "" {
		writeErr(w, http.StatusServiceUnavailable, "sign-in confirmation is not configured")
		return
	}
	state := oauth.Random(24)
	exp := s.now().Add(signInTTL)
	if err := putDoc(r.Context(), s.KV, store.Item{PK: pkSignIn(state), SK: "-", TTL: exp.Unix()}, pendingSignIn{Tenant: g.Tenant, Grant: g.ID, Expires: exp}, store.Cond{NotExists: true}); err != nil {
		s.fail(w, r, err)
		return
	}
	id := s.connectorIdentity(g)
	// Each attempt names its connection afresh: the provider keeps a name
	// once a link is made, so an attempt that failed after that point must
	// not block the next one.
	auth, err := p.Authorize(r.Context(), id, g.ID+"_"+oauth.Random(6), s.BaseURL+"/connectors/callback?state="+url.QueryEscape(state))
	if err != nil {
		_, _, _ = s.KV.Delete(r.Context(), pkSignIn(state), "-")
		s.fail(w, r, err)
		return
	}
	a := sr.approval
	g.Connector = &ConnectorState{User: id.User, Session: auth.Session, Account: auth.Account}
	g.Approval, g.Status = &a, StatusApproved
	g.HeldUntil = exp.UTC()
	g.event(s.now(), "sign-in approved "+a.Who())
	if err := s.saveGrant(r.Context(), g); err != nil {
		_, _ = p.Revoke(r.Context(), id, auth.Account)
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": g.Status, "sign_in_url": auth.URL, "expires_at": exp.UTC()})
}

func (s *Server) connectorCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		writeErr(w, http.StatusBadRequest, "missing sign-in state")
		return
	}
	it, ok, err := s.KV.Delete(r.Context(), pkSignIn(state), "-")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var pending pendingSignIn
	if !ok || json.Unmarshal(it.Doc, &pending) != nil || !s.now().Before(pending.Expires) {
		writeErr(w, http.StatusBadRequest, "unknown, expired or used sign-in")
		return
	}
	g, err := s.loadGrant(r.Context(), pending.Tenant, pending.Grant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Request.Connector == nil || g.Status != StatusApproved {
		writeErr(w, http.StatusConflict, "the connection is not waiting for sign-in")
		return
	}
	// Callback parameters never choose an account or establish authentication.
	if err := s.checkConnector(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	u, err := url.Parse(s.AppReturn)
	if err != nil || u.Scheme == "" {
		writeErr(w, http.StatusServiceUnavailable, "sign-in confirmation is not configured")
		return
	}
	code := oauth.Random(24)
	sum := sha256.Sum256([]byte(code))
	g.Account, g.ConfirmHash, g.ConfirmBy = g.Connector.Account, hex.EncodeToString(sum[:]), s.now().Add(signInTTL).UTC()
	g.Status = StatusSignedIn
	g.HeldUntil = time.Time{}
	g.event(s.now(), "signed in; waiting for the approving device to confirm")
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	q := u.Query()
	q.Set("tenant", g.Tenant)
	q.Set("grant", g.ID)
	q.Set("confirm", code)
	u.RawQuery = q.Encode()
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (s *Server) connectorTools(ctx context.Context, g *Grant) ([]mcpTool, mcp.Server, error) {
	_, err := s.provider(g)
	if err != nil {
		return nil, mcp.Server{}, err
	}
	if err := s.checkConnector(ctx, g); err != nil {
		return nil, mcp.Server{}, err
	}
	definitions := g.Request.Connector.Tools
	if len(g.Request.Tools) == 0 {
		return nil, mcp.Server{Name: g.VendorName}, nil
	}
	listed := make([]mcpTool, 0, len(definitions))
	aliases := map[string]bool{}
	for _, t := range definitions {
		if t.Validate() != nil || t.Toolkit != g.Request.Connector.Toolkit || aliases[t.Alias()] {
			return nil, mcp.Server{}, connector.ErrDefinition
		}
		aliases[t.Alias()] = true
		annotations, _ := json.Marshal(map[string]bool{"readOnlyHint": t.ReadOnly})
		listed = append(listed, mcpTool{Name: t.Alias(), Title: label(t.Name), Description: t.Description, InputSchema: t.InputSchema, Annotations: annotations})
	}
	return listed, mcp.Server{Name: g.VendorName}, nil
}

// Tool discovery pages are separate records. A catalog with thousands of
// actions never makes a grant exceed the store's per-item size limit.
func connectorToolKey(g *Grant, alias string) string { return "CONNECTOR#" + g.ID + "#" + alias }
func (s *Server) connectorCatalog(w http.ResponseWriter, r *http.Request, g *Grant) {
	if !g.SignedIn || g.Status == StatusFailed || g.Status == StatusDisconnected {
		writeErr(w, http.StatusConflict, "confirm this connection before reading its tool catalog")
		return
	}
	if err := s.checkConnector(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	p, err := s.provider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	page, err := p.Tools(r.Context(), g.Request.Connector.Toolkit, r.URL.Query().Get("cursor"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	tools := make([]CatalogTool, 0, len(page.Items))
	for _, t := range page.Items {
		if t.Validate() != nil || t.Toolkit != g.Request.Connector.Toolkit {
			s.fail(w, r, connector.ErrDefinition)
			return
		}
		if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(g.Tenant), SK: connectorToolKey(g, t.Alias()), TTL: s.now().Add(holdTTL).Unix()}, t, store.Cond{}); err != nil {
			s.fail(w, r, err)
			return
		}
		tools = append(tools, CatalogTool{Name: t.Alias(), Title: label(t.Name), Description: t.Description, InputSchema: t.InputSchema, ReadOnly: t.ReadOnly})
	}
	writeJSON(w, http.StatusOK, map[string]any{"grant": g.ID, "server": g.Request.Server, "account": g.Account, "tools": tools, "next_cursor": page.Next, "skipped": page.Skipped})
}
func (s *Server) selectConnectorTools(ctx context.Context, g *Grant, tools []ToolRequest) error {
	if err := s.checkConnector(ctx, g); err != nil {
		return err
	}
	g.Catalog = nil
	g.Connector.Tools = nil
	for _, ask := range tools {
		var definition connector.Tool
		it, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), connectorToolKey(g, ask.Tool), &definition)
		if err != nil || it.TTL <= s.now().Unix() || definition.Alias() != ask.Tool || definition.Toolkit != g.Request.Connector.Toolkit {
			return apiError{http.StatusBadRequest, "read the app's tool catalog before selecting this action"}
		}
		g.Connector.Tools = append(g.Connector.Tools, definition)
		g.Catalog = append(g.Catalog, CatalogTool{Name: ask.Tool, Title: label(definition.Name), Description: definition.Description, InputSchema: definition.InputSchema, ReadOnly: definition.ReadOnly})
	}
	return nil
}
func bindConnectorTools(g *Grant, requested []ToolRequest) error {
	if g.Connector == nil {
		return connector.ErrAccount
	}
	definitions := make([]connector.Tool, 0, len(requested))
	seen := map[string]bool{}
	for _, ask := range requested {
		if seen[ask.Tool] {
			continue
		}
		seen[ask.Tool] = true
		var found *connector.Tool
		for _, existing := range g.Tools {
			if existing.Tool == ask.Tool && existing.Connector != nil {
				copy := *existing.Connector
				found = &copy
				break
			}
		}
		if found == nil {
			for _, candidate := range g.Connector.Tools {
				if candidate.Alias() == ask.Tool {
					copy := candidate
					found = &copy
					break
				}
			}
		}
		if found == nil || found.Validate() != nil || found.Toolkit != g.Request.Connector.Toolkit {
			return connector.ErrDefinition
		}
		definitions = append(definitions, *found)
	}
	raw, err := json.Marshal(definitions)
	if err != nil {
		return err
	}
	if len(raw) > 64<<10 {
		return apiError{http.StatusBadRequest, "the selected tool definitions exceed 64 KB; select fewer actions"}
	}
	g.Connector.Tools = nil
	g.Request.Connector.Tools = definitions
	g.Request.Connector.Account = g.Connector.Account
	b, err := json.Marshal(definitions)
	if err != nil {
		return err
	}
	g.Request.API.Digest, err = jcs.SHA256(b)
	return err
}

func (s *Server) callConnector(ctx context.Context, g *Grant, t GrantedTool, args json.RawMessage, live func() bool) (string, map[string]any) {
	p, err := s.provider(g)
	if err != nil || t.Connector == nil || g.Connector == nil || g.Request.Connector.Account != g.Connector.Account || g.Connector.User != s.connectorIdentity(g).User || g.Connector.Revocation != nil && g.Connector.Revocation.Disabled {
		return OpNotSent, refusal(CodeConnectionDown, "Not sent: the approved connector definition or account is unavailable.")
	}
	result, err := p.Execute(ctx, s.connectorIdentity(g), g.Request.Connector.Account, *t.Connector, args, live)
	if errors.Is(err, connector.ErrChanged) {
		// The provider's definition of an approved tool was read and changed:
		// the grant is paused until the owner approves it again
		// (kiff/planque#71). A definition that couldn't be read or checked
		// (ErrDefinition alone) pauses nothing: it's "Not sent" below.
		_, why := s.pauseGrant(ctx, g, t.Name, vendorName(g)+" changed "+t.Name+" since you approved it")
		return OpNotSent, refusal(CodeOperationChanged, "Not sent: "+why+".")
	}
	if err != nil {
		// Every error here is before the provider ran the action: a refusal
		// with its code, so KIFF never reads it as the vendor's.
		code := CodeConnectionDown
		switch {
		case errors.Is(err, connector.ErrPermit):
			code = CodePermitExpired
		case errors.Is(err, connector.ErrDefinition):
			code = CodeDefinitionUnread
		}
		return OpNotSent, refusal(code, "Not sent: "+err.Error())
	}
	s.noteAccount(ctx, g, result)
	if result.Unknown {
		return OpUnknown, toolResult("The connector did not return a confirmed result"+providerSaid(result)+". This operation will not be sent again.", true)
	}
	if result.Refused {
		return OpDone, toolResult(refusedWords(vendorName(g), result)+" No claim is made about whether an effect occurred.", true)
	}
	return OpDone, toolResult(string(result.Data), false)
}

// refusedWords says who refused a connector call and why, so the taf can
// tell the owner: the vendor's own answer ("Stripe refused it: No such
// invoice: 'in_1'.") or the connection provider's ("The connection provider
// refused it (it answered 403): …"). planque's own refusals are "Not sent: …" and
// KIFF's never reach here, so the three read apart.
func refusedWords(vendor string, r connector.Result) string {
	if r.ByVendor {
		if r.Reason == "" {
			return vendor + " refused it and gave no reason."
		}
		return vendor + " refused it: " + sentence(r.Reason)
	}
	return "The connection provider refused it" + providerSaid(r) + "."
}

// providerSaid is the provider's status and reason, when it gave them, in
// the form KIFF's gateway reads a status from ("answered 403").
func providerSaid(r connector.Result) string {
	out := ""
	if r.Status != 0 {
		out = fmt.Sprintf(" (it answered %d)", r.Status)
	}
	if r.Reason != "" {
		out += ": " + strings.TrimRight(r.Reason, ".")
	}
	return out
}

func sentence(s string) string {
	if strings.HasSuffix(s, ".") || strings.HasSuffix(s, "…") {
		return s
	}
	return s + "."
}

func (s *Server) revokeConnector(ctx context.Context, g *Grant) error {
	triggerErr := s.stopTriggers(ctx, g)
	if g.Connector == nil {
		return triggerErr
	}
	if g.Connector.Revocation != nil && g.Connector.Revocation.Disabled {
		return triggerErr
	}
	p, err := s.provider(g)
	if err != nil {
		return errors.Join(triggerErr, err)
	}
	if g.Connector.User != s.connectorIdentity(g).User {
		return errors.Join(triggerErr, connector.ErrAccount)
	}
	receipt, err := p.Revoke(ctx, s.connectorIdentity(g), g.Connector.Account)
	if err != nil {
		return errors.Join(triggerErr, err)
	}
	g.Connector.Revocation = &receipt
	if !receipt.Disabled {
		return errors.New("the connector did not confirm that access was disabled")
	}
	return triggerErr
}

// connectorText is what the owner reads and signs for a managed connection.
// A grant made with the app's title names the app and the taf in a person's
// words, and no connection provider; one made before keeps the words it was
// signed with.
func connectorText(q Request) string {
	text := "Connect " + q.Connector.Toolkit + " through " + q.Connector.Provider + " for " + q.Agent + ". This connection provider holds the vendor credentials."
	if q.Connector.Title != "" {
		taf := "the taf " + q.Agent
		if q.AgentName != "" {
			taf = q.AgentName
		}
		text = "Connect " + q.Connector.Title + " for " + taf + ". Your sign-in to " + q.Connector.Title + " stays with the connection service; the taf never sees it."
		switch q.Connector.AuthScheme {
		case "API_KEY", "BEARER_TOKEN":
			credential := "key"
			if q.Connector.AuthScheme == "BEARER_TOKEN" {
				credential = "token"
			}
			text = "Connect " + q.Connector.Title + " for " + taf + ". You'll enter your " + q.Connector.Title + " " + credential + " on the connection service's page; it keeps the " + credential + ", and the taf never sees it."
		}
	}
	if len(q.Tools) == 0 {
		return text + " Sign in and confirm on this device, then choose the actions to approve. No actions are granted yet."
	}
	text += " Use only connected account " + q.Connector.Account + ". Approve these exact versioned actions:"
	for _, t := range q.Connector.Tools {
		text += " " + t.Name + " (" + t.Slug + ", " + t.Version + ")."
	}
	return text
}
