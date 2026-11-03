package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
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
	Provider string           `json:"provider"`
	Toolkit  string           `json:"toolkit"`
	Account  string           `json:"account,omitempty"`
	Tools    []connector.Tool `json:"tools,omitempty"`
}

// ConnectorState is private server state, never an agent-selected identity.
// Vendor OAuth credentials remain at the configured connection provider.
type ConnectorState struct {
	User       string                `json:"user"`
	Session    string                `json:"session"`
	Account    string                `json:"account"`
	Tools      []connector.Tool      `json:"tools,omitempty"`
	Revocation *connector.Revocation `json:"revocation,omitempty"`
}

func connectorEntry(id string) (string, string, bool) {
	p := strings.Split(id, ":")
	if len(p) != 3 || p[0] != "connector" || p[1] == "" || p[2] == "" {
		return "", "", false
	}
	return p[1], p[2], true
}
func (s *Server) connectorIdentity(g *Grant) connector.Identity {
	sum := sha256.Sum256([]byte(s.BaseURL + "\n" + g.Tenant))
	return connector.Identity{User: "planque_" + hex.EncodeToString(sum[:]), Toolkit: g.Request.Connector.Toolkit}
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
	}
	if err := decodeStrict(body, &ask); err != nil {
		s.fail(w, r, err)
		return
	}
	if msg := checkTools(ask.Agent, nil); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	app, err := p.App(r.Context(), toolkit)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	platform := p.Platform()
	if app.Slug != toolkit || !catalogHTTPS(platform.URL) || label(platform.Name) == "" {
		writeErr(w, http.StatusServiceUnavailable, "the catalog provider is not configured")
		return
	}
	id := "grt_" + oauth.Random(12)
	g := &Grant{ID: id, Tenant: tenant, Status: StatusPending, VendorName: label(app.Name),
		Request: Request{Version: "planque.grant.v1", Grant: id, Tenant: tenant, Server: platform.URL, SignInAt: platform.URL,
			Agent: ask.Agent, Tools: []ToolRequest{}, AskedBy: me.AgentID, CreatedAt: s.now().UTC().Format(time.RFC3339), Kind: "connector",
			CatalogEntry: r.PathValue("entry"), Connector: &ConnectorRequest{Provider: providerID, Toolkit: toolkit},
			API: &APIRequest{URL: platform.URL, Base: platform.URL, SourceKind: "connector"}},
		Audience: s.BaseURL + "/mcp/" + url.PathEscape(tenant) + "/" + id}
	if g.Request.AskedBy == "" {
		g.Request.AskedBy = "an account key"
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
	auth, err := p.Authorize(r.Context(), id, g.ID, s.BaseURL+"/connectors/callback?state="+url.QueryEscape(state))
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
		return OpNotSent, toolResult("Not sent: the approved connector definition or account is unavailable.", true)
	}
	result, err := p.Execute(ctx, s.connectorIdentity(g), g.Request.Connector.Account, *t.Connector, args, live)
	if err != nil {
		return OpNotSent, toolResult("Not sent: "+err.Error(), true)
	}
	if result.Unknown {
		return OpUnknown, toolResult("The connector did not return a confirmed result. This operation will not be sent again.", true)
	}
	if result.Refused {
		return OpDone, toolResult("The connector refused the operation. No claim is made about whether an effect occurred.", true)
	}
	return OpDone, toolResult(string(result.Data), false)
}
func (s *Server) revokeConnector(ctx context.Context, g *Grant) error {
	if g.Connector == nil {
		return nil
	}
	if g.Connector.Revocation != nil && g.Connector.Revocation.Disabled {
		return nil
	}
	p, err := s.provider(g)
	if err != nil {
		return err
	}
	if g.Connector.User != s.connectorIdentity(g).User {
		return connector.ErrAccount
	}
	receipt, err := p.Revoke(ctx, s.connectorIdentity(g), g.Connector.Account)
	if err != nil {
		return err
	}
	g.Connector.Revocation = &receipt
	if !receipt.Disabled {
		return errors.New("the connector did not confirm that access was disabled")
	}
	return nil
}
func connectorText(q Request) string {
	text := "Connect " + q.Connector.Toolkit + " through " + q.Connector.Provider + " for " + q.Agent + ". This connection provider holds the vendor credentials."
	if len(q.Tools) == 0 {
		return text + " Sign in and confirm on this device, then choose the actions to approve. No actions are granted yet."
	}
	text += " Use only connected account " + q.Connector.Account + ". Approve these exact versioned actions:"
	for _, t := range q.Connector.Tools {
		text += " " + t.Name + " (" + t.Slug + ", " + t.Version + ")."
	}
	return text
}
