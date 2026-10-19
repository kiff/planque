package planque

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

// signInTTL is how long a sign-in link works.
const signInTTL = 15 * time.Minute

// holdTTL is how long planque keeps a confirmed sign-in that connects
// nothing (connect first): the tools must be set and approved by then.
const holdTTL = 24 * time.Hour

var toolName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,63}$`)

// toolKey is a tool's name in KIFF when the request names none: the
// vendor's name and four letters or digits of the grant's id. KIFF keeps
// one tool per name on an account, so two tafs on one vendor each get
// their own (kiff/tafs#78). A stopgap until the gateway makes the key
// (RFC 048, kiff/kiff-cloud#1266).
func toolKey(tool, grant string) string {
	var suffix []rune
	for _, r := range strings.ToLower(strings.TrimPrefix(grant, "grt_")) {
		if len(suffix) == 4 {
			break
		}
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			suffix = append(suffix, r)
		}
	}
	if len(tool) > 59 {
		tool = tool[:59]
	}
	return tool + "_" + string(suffix)
}

// requestGrant is an agent asking for a vendor account for a taf. It
// grants nothing: it returns the plain words the owner will see, and the
// digest their phone signs to approve exactly this.
func (s *Server) requestGrant(w http.ResponseWriter, r *http.Request) {
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
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		Server string         `json:"server"`
		Agent  string         `json:"agent"`
		Tools  *[]ToolRequest `json:"tools"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	if req.Tools == nil {
		writeErr(w, http.StatusBadRequest, "tools is required: the server's tools for the taf, or [] to connect first and see what it offers")
		return
	}
	if msg := checkTools(req.Agent, *req.Tools); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	server := strings.TrimSpace(req.Server)
	meta, err := s.OAuth.Discover(r.Context(), server)
	if err != nil {
		writeErr(w, http.StatusBadGateway, "Couldn't find how to sign in to "+server+": "+err.Error())
		return
	}
	id := "grt_" + oauth.Random(12)
	asked := me.AgentID
	if asked == "" {
		asked = "an account key"
	}
	tools := *req.Tools
	for i := range tools {
		if tools[i].Name == "" {
			tools[i].Name = toolKey(tools[i].Tool, id)
		}
	}
	g := &Grant{ID: id, Tenant: tenant, Status: StatusPending, Vendor: meta,
		Request: Request{Version: "planque.grant.v1", Grant: id, Tenant: tenant, Server: server, SignInAt: meta.Issuer,
			Agent: req.Agent, Tools: tools, AskedBy: asked, CreatedAt: s.now().UTC().Format(time.RFC3339)},
		Audience: s.BaseURL + "/mcp/" + url.PathEscape(tenant) + "/" + id}
	g.Text = grantText(g.Request)
	// Connect-first's text says the sign-in may be read (kiff/tafs#89).
	g.Describable = len(tools) == 0
	g.Digest, err = digest(g.Request)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g.event(s.now(), "requested by "+asked+" for agent "+req.Agent)
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(r.Context(), "grant requested", slog.String("grant", id), slog.String("server", hostOf(server)))
	writeJSON(w, http.StatusCreated, s.view(g))
}

func checkTools(agent string, tools []ToolRequest) string {
	if strings.TrimSpace(agent) == "" {
		return "agent is required: the KIFF agent of the taf that will use the account"
	}
	if len(tools) > 20 {
		return "tools may name at most 20 of the server's tools"
	}
	seen := map[string]bool{}
	for _, t := range tools {
		name := t.Name
		if name == "" {
			name = t.Tool
		}
		switch {
		case !toolName.MatchString(t.Tool) || !toolName.MatchString(name):
			return "tool names must start with a letter and use letters, digits and underscores"
		case seen[name]:
			return "tool " + name + " is asked for twice"
		case t.AmountArg != "" && t.ReadOnly:
			return "a tool that moves an amount is not read-only: " + t.Tool
		case t.AmountArg == IdempotencyArg:
			return IdempotencyArg + " is planque's own argument"
		}
		if msg := checkNarrowing(t); msg != "" {
			return msg
		}
		seen[name] = true
	}
	return ""
}

// grantText is what the owner reads before approving.
func grantText(q Request) string { return requestText(q, false, "") }

// requestText is what the owner reads. held is for tools asked for on a
// grant whose sign-in planque already holds: there is no second sign-in,
// and the text names the account it acts in, when the vendor names it.
func requestText(q Request, held bool, account string, catalog ...CatalogTool) string {
	if len(q.Tools) == 0 {
		return fmt.Sprintf("Connect your account at %s to see what it offers for the taf %s.\n"+
			"You will sign in at %s.\n\nOnly what it offers is read: its list of tools, and what your sign-in covers (such as which accounts, test or live), "+
			"through its tools that change nothing and take no input. Nothing is changed and nothing is given to the taf. "+
			"The tools the taf may use come next, and you approve them again. If they aren't approved within a day, the sign-in is revoked.\n\n"+
			"Your sign-in is kept in a sealed vault; the taf never sees it. You can disconnect at any time.",
			hostOf(q.Server), q.Agent, hostOf(q.SignInAt))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Connect your account at %s to the taf %s.\n", hostOf(q.Server), q.Agent)
	switch {
	case held && account != "":
		fmt.Fprintf(&b, "You are already signed in at %s, to the account %s: there is no second sign-in.\n\nThe taf may use:\n", hostOf(q.SignInAt), account)
	case held:
		fmt.Fprintf(&b, "You are already signed in at %s: there is no second sign-in.\n\nThe taf may use:\n", hostOf(q.SignInAt))
	default:
		fmt.Fprintf(&b, "You will sign in at %s.\n\nThe taf may use:\n", hostOf(q.SignInAt))
	}
	for _, t := range q.Tools {
		how := "makes changes, only within the taf's Card"
		switch {
		case t.ReadOnly:
			how = "reads only"
		case t.AmountArg != "":
			how = "moves money (" + t.AmountArg + "), only within the taf's Card"
		case !t.narrowed():
			// Nothing narrows it and no amount bounds it: the Card can
			// allow or refuse each call, not limit what it changes.
			how = "can make any change this tool allows; the taf's Card can allow or refuse each call, not limit what it changes"
		}
		fmt.Fprintf(&b, "• %s%s: %s\n", t.Name, narrowingText(t, catalog...), how)
	}
	b.WriteString("\nYour sign-in is kept in a sealed vault; the taf never sees it. Each call is sent only with KIFF's permit for that exact call. You can disconnect at any time.")
	return b.String()
}

// digest is the SHA-256 of the request's canonical JSON: what the phone
// signs to approve it.
func digest(q Request) (string, error) {
	b, err := json.Marshal(q)
	if err != nil {
		return "", err
	}
	return jcs.SHA256(b)
}

// grantView is a grant as planque shows it: never a secret.
type grantView struct {
	ID string `json:"id"`
	// Vendor is the vendor's name for people: what its server calls itself,
	// else its host. The app shows it as it is.
	Vendor   string        `json:"vendor"`
	Status   string        `json:"status"`
	Text     string        `json:"text"`
	Digest   string        `json:"digest"`
	Request  Request       `json:"request"`
	Account  string        `json:"account,omitempty"`
	Tools    []string      `json:"tools,omitempty"`
	Audience string        `json:"audience"`
	Reason   string        `json:"reason,omitempty"`
	Approval *approvalView `json:"approval,omitempty"`
	// SignedIn: planque holds a confirmed sign-in for this grant, so
	// approving its tools needs no second sign-in.
	SignedIn  bool      `json:"signed_in,omitempty"`
	HeldUntil time.Time `json:"held_until,omitzero"`
	Events    []Event   `json:"events"`
}

type approvalView struct {
	Fingerprint string    `json:"device"`
	DeviceKey   []byte    `json:"device_key,omitempty"` // SPKI DER, base64
	DeviceID    string    `json:"device_id,omitempty"`  // KIFF's, for a device KIFF knows
	Message     string    `json:"message"`
	Body        string    `json:"body"`
	Signature   string    `json:"signature"`
	SignedAt    string    `json:"signed_at"`
	ApprovedAt  time.Time `json:"approved_at"`
}

func (s *Server) view(g *Grant) grantView {
	v := grantView{ID: g.ID, Vendor: vendorName(g), Status: g.Status, Text: g.Text, Digest: g.Digest, Request: g.Request,
		Account: g.Account, Audience: g.Audience, Reason: g.Reason, Events: g.Events,
		SignedIn: g.SignedIn && g.Status != StatusDisconnected && g.Status != StatusFailed, HeldUntil: g.HeldUntil}
	for _, t := range g.Tools {
		v.Tools = append(v.Tools, t.Name)
	}
	if a := g.Approval; a != nil {
		v.Approval = &approvalView{Fingerprint: a.Fingerprint, DeviceKey: a.DeviceKey, DeviceID: a.DeviceID, Message: a.Message,
			Body: a.Body, Signature: a.Signature, SignedAt: a.SignedAt, ApprovedAt: a.ApprovedAt}
	}
	return v
}

// getGrant shows a grant to an agent of its account, so the setup agent
// can follow it.
func (s *Server) getGrant(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
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
	s.expire(r.Context(), g)
	writeJSON(w, http.StatusOK, s.view(g))
}

// catalog is what the vendor's server lists, for the agent choosing the
// tools of a grant asked for with no tools: each tool's own words and
// arguments, and the account the sign-in acts in when the vendor names it.
func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
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
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if me.AgentID == "" || me.AgentID != g.Request.AskedBy {
		writeErr(w, http.StatusForbidden, "only the agent that asked for this grant may read its catalog")
		return
	}
	s.expire(r.Context(), g)
	if len(g.Catalog) == 0 || g.Status == StatusFailed || g.Status == StatusDisconnected {
		writeErr(w, http.StatusConflict, "this grant has no catalog: it is "+g.Status+", and only a grant asked for with no tools lists them, once the owner signed in")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"grant": g.ID, "server": g.Request.Server, "account": g.Account, "tools": g.Catalog})
}

// setTools is the asking agent choosing, from the catalog, the tools a
// listed grant gives the taf. The request, its text and its digest change,
// and the owner approves them on the phone; the sign-in is already held.
func (s *Server) setTools(w http.ResponseWriter, r *http.Request) {
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
		Tools []ToolRequest `json:"tools"`
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
		writeErr(w, http.StatusForbidden, "only the agent that asked for this grant may set its tools")
		return
	}
	s.expire(r.Context(), g)
	// Listed, or tools set and not approved yet (the agent may revise them).
	if !g.SignedIn || len(g.Catalog) == 0 || (g.Status != StatusListed && g.Status != StatusPending) {
		writeErr(w, http.StatusConflict, "only a listed grant takes its tools; this one is "+g.Status)
		return
	}
	if len(req.Tools) == 0 {
		writeErr(w, http.StatusBadRequest, "tools must name 1 to 20 of the catalog's tools")
		return
	}
	if msg := checkTools(g.Request.Agent, req.Tools); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	for i := range req.Tools {
		if !slices.ContainsFunc(g.Catalog, func(c CatalogTool) bool { return c.Name == req.Tools[i].Tool }) {
			writeErr(w, http.StatusBadRequest, hostOf(g.Request.Server)+" lists no tool named "+req.Tools[i].Tool)
			return
		}
		if req.Tools[i].Name == "" {
			req.Tools[i].Name = toolKey(req.Tools[i].Tool, g.ID)
		}
	}
	g.Request.Tools = req.Tools
	g.Text = requestText(g.Request, true, g.Account, g.Catalog...)
	if g.Digest, err = digest(g.Request); err != nil {
		s.fail(w, r, err)
		return
	}
	g.Status, g.Approval = StatusPending, nil
	g.event(s.now(), "tools set by "+me.AgentID+"; waiting for the owner's approval")
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, s.view(g))
}

// expire ends a sign-in nobody confirmed in time: its token is revoked at
// the vendor and the grant fails. A sign-in that is never confirmed would
// otherwise keep a live token at the vendor, so every read of the grant
// (status, list, export) checks, as confirm does.
func (s *Server) expire(ctx context.Context, g *Grant) {
	held := g.SignedIn && (g.Status == StatusListed || g.Status == StatusPending)
	if held && !g.HeldUntil.IsZero() && s.now().After(g.HeldUntil) {
		s.revoke(ctx, g)
		g.HeldUntil = time.Time{}
		_ = s.failGrant(ctx, g, "its tools were not approved within a day; the sign-in was revoked")
		s.log().InfoContext(ctx, "held sign-in revoked", slog.String("grant", g.ID))
		return
	}
	if g.Status != StatusSignedIn || !s.now().After(g.ConfirmBy) {
		return
	}
	s.revoke(ctx, g)
	g.ConfirmHash, g.ConfirmBy = "", time.Time{}
	_ = s.failGrant(ctx, g, "the sign-in was not confirmed in time; its token was revoked")
	s.log().InfoContext(ctx, "unconfirmed sign-in revoked", slog.String("grant", g.ID))
}

// approve is the owner's Face ID on the phone: it signs this request,
// whose body names the grant's digest. Only then does planque register
// itself at the vendor and hand back the vendor's own sign-in link.
func (s *Server) approve(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		Digest string `json:"digest"`
	}
	if err := decodeStrict(sr.body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if req.Digest != g.Digest {
		writeErr(w, http.StatusConflict, "the signed digest is not this grant's: nothing was granted")
		return
	}
	if g.Status == StatusListed {
		writeErr(w, http.StatusConflict, "this grant has no tools to approve yet")
		return
	}
	if g.Status != StatusPending && g.Status != StatusApproved && g.Status != StatusSignedIn {
		writeErr(w, http.StatusConflict, "this grant is "+g.Status)
		return
	}
	s.expire(r.Context(), g)
	if g.Status == StatusFailed {
		writeErr(w, http.StatusConflict, "this grant is failed: "+g.Reason)
		return
	}
	if g.SignedIn && g.Status == StatusPending {
		// Tools set on a grant whose sign-in is held: approving them
		// connects at once, with no second sign-in.
		a := sr.approval
		g.Approval = &a
		g.event(s.now(), "tools approved "+a.Who())
		tok, _, err := s.openToken(r.Context(), g)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if err := s.connect(r.Context(), g, tok.Access); err != nil {
			s.log().WarnContext(r.Context(), "connect failed", slog.String("grant", g.ID), slog.String("err", err.Error()))
			writeErr(w, http.StatusBadGateway, g.Reason)
			return
		}
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
	if g.Status == StatusSignedIn {
		// A new link replaces a sign-in nobody confirmed: its token goes.
		s.revoke(r.Context(), g)
		g.ConfirmHash, g.ConfirmBy = "", time.Time{}
		g.event(s.now(), "an unconfirmed sign-in was revoked")
	}
	clientID, err := s.client(r.Context(), g.Vendor)
	if err != nil {
		writeErr(w, http.StatusBadGateway, "Couldn't register with "+hostOf(g.Vendor.Issuer)+": "+err.Error())
		return
	}
	pkce := oauth.NewPKCE()
	state := oauth.Random(24)
	sealed, err := s.Sealer.Seal(r.Context(), []byte(pkce.Verifier), map[string]string{"tenant": g.Tenant, "grant": g.ID, "kind": "pkce"})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	exp := s.now().Add(signInTTL)
	err = putDoc(r.Context(), s.KV, store.Item{PK: pkSignIn(state), SK: "-", TTL: exp.Unix()},
		pendingSignIn{Tenant: g.Tenant, Grant: g.ID, Verifier: sealed, Expires: exp}, store.Cond{NotExists: true})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	a := sr.approval
	if g.Status == StatusPending {
		g.Approval = &a
		g.event(s.now(), "approved "+a.Who())
	} else {
		g.event(s.now(), "a new sign-in link was signed for "+a.Who())
	}
	g.Status, g.ClientID = StatusApproved, clientID
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": g.Status,
		"sign_in_url": oauth.AuthURL(g.Vendor, clientID, s.redirect(), state, pkce.Challenge), "expires_at": exp.UTC()})
}

// ClientName is the name a vendor shows the owner when it asks them to
// allow the sign-in: the app they are connecting, tafs, not planque
// (kiff/tafs#82).
const ClientName = "tafs"

// client returns planque's client_id at the vendor's authorization
// server, registering once per server and name: a new name registers
// again, so the vendor shows it.
func (s *Server) client(ctx context.Context, m oauth.Metadata) (string, error) {
	var c struct {
		ClientID string `json:"client_id"`
	}
	pk := pkClient(m.Issuer, s.redirect(), ClientName)
	if _, err := getDoc(ctx, s.KV, pk, "-", &c); err == nil && c.ClientID != "" {
		return c.ClientID, nil
	}
	id, err := s.OAuth.Register(ctx, m, ClientName, s.redirect())
	if err != nil {
		return "", err
	}
	c.ClientID = id
	if err := putDoc(ctx, s.KV, store.Item{PK: pk, SK: "-"}, c, store.Cond{}); err != nil {
		return "", err
	}
	return id, nil
}

// callback is where the vendor sends the owner back. The state must be
// one planque issued and not yet used; the code is exchanged with the
// PKCE verifier and the exact redirect URI.
func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state == "" {
		s.page(w, http.StatusBadRequest, "Not connected", "This link is missing its sign-in state.")
		return
	}
	// Taking the state deletes it: a link works once.
	it, ok, err := s.KV.Delete(r.Context(), pkSignIn(state), "-")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var p pendingSignIn
	if !ok || json.Unmarshal(it.Doc, &p) != nil || s.now().After(p.Expires) {
		s.page(w, http.StatusBadRequest, "Not connected", "This sign-in link is unknown, expired or already used. Approve again on your phone.")
		return
	}
	g, err := s.loadGrant(r.Context(), p.Tenant, p.Grant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Status != StatusApproved {
		s.page(w, http.StatusConflict, "Not connected", "This grant is "+g.Status+".")
		return
	}
	if e := q.Get("error"); e != "" {
		g.event(s.now(), "the vendor sign-in ended: "+e)
		_ = s.saveGrant(r.Context(), g)
		s.page(w, http.StatusOK, "Not connected", "The sign-in at "+hostOf(g.Vendor.Issuer)+" was not completed. Nothing was granted.")
		return
	}
	// RFC 9207: when the server names itself, it must be the one we used,
	// and a server that says it always does must have.
	if iss := q.Get("iss"); (iss == "" && g.Vendor.IssParameter) || (iss != "" && iss != g.Vendor.Issuer) {
		s.page(w, http.StatusBadRequest, "Not connected", "This answer came from another sign-in server. Nothing was granted.")
		return
	}
	verifier, err := s.Sealer.Open(r.Context(), p.Verifier, map[string]string{"tenant": g.Tenant, "grant": g.ID, "kind": "pkce"})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	tok, err := s.OAuth.Exchange(r.Context(), g.Vendor, g.ClientID, s.redirect(), q.Get("code"), string(verifier))
	if err != nil {
		s.failGrant(r.Context(), g, "the sign-in code was not accepted: "+err.Error())
		s.page(w, http.StatusBadGateway, "Not connected", "The sign-in at "+hostOf(g.Vendor.Issuer)+" couldn't be finished.")
		return
	}
	if s.AppReturn == "" {
		s.failGrant(r.Context(), g, "planque has no app to confirm sign-ins with (PLANQUE_APP_RETURN)")
		s.page(w, http.StatusInternalServerError, "Not connected", "Sign-ins can't be confirmed right now.")
		return
	}
	// The token is kept but not used: nothing is connected until the phone
	// that approved confirms with the code only this browser receives.
	if err := s.storeToken(r.Context(), g, tok, nil); err != nil {
		s.failGrant(r.Context(), g, "the token could not be stored")
		s.page(w, http.StatusInternalServerError, "Not connected", "The sign-in couldn't be kept.")
		return
	}
	code := oauth.Random(24)
	sum := sha256.Sum256([]byte(code))
	g.Account, g.ConfirmHash, g.ConfirmBy = tok.Account, hex.EncodeToString(sum[:]), s.now().Add(signInTTL).UTC()
	g.Status = StatusSignedIn
	g.event(s.now(), "signed in at "+hostOf(g.Vendor.Issuer)+accountNote(g.Account)+"; waiting for the phone to confirm")
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.revoke(r.Context(), g)
		s.fail(w, r, err)
		return
	}
	u, err := url.Parse(s.AppReturn)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	rq := u.Query()
	rq.Set("tenant", g.Tenant)
	rq.Set("grant", g.ID)
	rq.Set("confirm", code)
	u.RawQuery = rq.Encode()
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// confirm is the approving phone finishing a sign-in: it signs the code
// the vendor's sign-in handed to the owner's app. Only then are the tools
// registered and the grant connected. A link someone else opened hands
// its code to their browser, never to the phone that approved.
func (s *Server) confirm(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		Confirm string `json:"confirm"`
	}
	if err := decodeStrict(sr.body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Status != StatusSignedIn {
		writeErr(w, http.StatusConflict, "this grant is "+g.Status)
		return
	}
	if g.Approval == nil || !g.Approval.SameSigner(sr.approval) {
		writeErr(w, http.StatusForbidden, "only the device that approved this grant may confirm it")
		return
	}
	if s.now().After(g.ConfirmBy) {
		s.expire(r.Context(), g)
		writeErr(w, http.StatusConflict, "the sign-in was not confirmed in time; it was revoked")
		return
	}
	sum := sha256.Sum256([]byte(req.Confirm))
	want, _ := hex.DecodeString(g.ConfirmHash)
	if req.Confirm == "" || subtle.ConstantTimeCompare(sum[:], want) != 1 {
		writeErr(w, http.StatusForbidden, "this is not the code of this grant's sign-in")
		return
	}
	g.ConfirmHash, g.ConfirmBy = "", time.Time{}
	g.event(s.now(), "sign-in confirmed "+sr.approval.Who())
	tok, _, err := s.openToken(r.Context(), g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g.SignedIn = true
	if len(g.Request.Tools) == 0 {
		// Connect first: the vendor's tools are listed for the asking agent
		// to choose from, and nothing is registered in KIFF.
		if err := s.list(r.Context(), g, tok.Access); err != nil {
			writeErr(w, http.StatusBadGateway, g.Reason)
			return
		}
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
	if err := s.connect(r.Context(), g, tok.Access); err != nil {
		s.log().WarnContext(r.Context(), "connect failed", slog.String("grant", g.ID), slog.String("err", err.Error()))
		writeErr(w, http.StatusBadGateway, g.Reason)
		return
	}
	writeJSON(w, http.StatusOK, s.view(g))
}

// list reads the vendor's tools into the grant's catalog, with the token
// of a confirmed sign-in, and keeps the sign-in: the grant is listed.
func (s *Server) list(ctx context.Context, g *Grant, access string) error {
	listed, server, err := s.MCP.ListTools(ctx, g.Request.Server, access)
	if err != nil {
		s.revoke(ctx, g)
		return s.failGrant(ctx, g, "Couldn't list the tools at "+hostOf(g.Request.Server)+": "+err.Error())
	}
	g.VendorName = label(server.Label())
	g.Catalog = nil
	for _, t := range listed {
		var hints struct {
			ReadOnly bool `json:"readOnlyHint"`
		}
		_ = json.Unmarshal(t.Annotations, &hints)
		g.Catalog = append(g.Catalog, CatalogTool{Name: t.Name, Title: label(t.Label()), Description: t.Description,
			InputSchema: t.InputSchema, ReadOnly: hints.ReadOnly})
	}
	g.Status, g.HeldUntil = StatusListed, s.now().Add(holdTTL).UTC()
	g.event(s.now(), fmt.Sprintf("listed %d tools at %s; waiting for the tools the taf needs", len(g.Catalog), hostOf(g.Request.Server)))
	return s.saveGrant(ctx, g)
}

// connect reads the granted tools from the vendor with the stored token,
// and registers them in KIFF's gateway with the connector key. If any
// step fails, nothing stays connected and the token is revoked.
func (s *Server) connect(ctx context.Context, g *Grant, access string) error {
	listed, server, err := s.MCP.ListTools(ctx, g.Request.Server, access)
	if err != nil {
		s.revoke(ctx, g)
		return s.failGrant(ctx, g, "Couldn't list the tools at "+hostOf(g.Request.Server)+": "+err.Error())
	}
	if n := label(server.Label()); n != "" {
		g.VendorName = n
	}
	g.Tools = nil
	for _, t := range g.Request.Tools {
		i := slices.IndexFunc(listed, func(l mcpTool) bool { return l.Name == t.Tool })
		if i < 0 {
			names := make([]string, 0, len(listed))
			for _, l := range listed {
				names = append(names, l.Name)
			}
			s.revoke(ctx, g)
			return s.failGrant(ctx, g, hostOf(g.Request.Server)+" has no tool named "+t.Tool+"; it lists: "+strings.Join(names, ", "))
		}
		gt := GrantedTool{ToolRequest: t, Title: label(listed[i].Label()), Description: listed[i].Description, InputSchema: listed[i].InputSchema}
		if t.narrowed() {
			gt.InputSchema = t.schema()
			gt.Description = "Through planque: " + t.Tool + narrowingText(t) + "."
		}
		g.Tools = append(g.Tools, gt)
	}
	transport := oauth.Random(32)
	sum := sha256.Sum256([]byte(transport))
	g.Transport = hex.EncodeToString(sum[:])
	// The relay must answer KIFF's probe before the registration returns,
	// so the grant is saved as connected first; a failure undoes it.
	g.Status, g.HeldUntil = StatusConnected, time.Time{}
	if err := s.saveGrant(ctx, g); err != nil {
		s.revoke(ctx, g)
		return err
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		return s.undo(ctx, g, nil, "the KIFF connector key could not be opened")
	}
	var done []string
	for _, t := range g.Tools {
		err := s.KIFF.Connect(ctx, key, kiff.Connection{Name: t.Name, Title: t.Title, URL: g.Audience, Tool: t.Name,
			Description: t.Description, AmountArg: t.AmountArg, AmountUnit: t.AmountUnit, AmountScale: t.AmountScale,
			IdempotencyArg: IdempotencyArg, ReadOnly: t.ReadOnly, Mode: "relay", Audience: g.Audience, Credential: transport})
		if err != nil {
			return s.undo(ctx, g, done, "KIFF did not register "+t.Name+": "+err.Error())
		}
		done = append(done, t.Name)
	}
	g.event(s.now(), "registered in KIFF in relay mode: "+strings.Join(done, ", "))
	return s.saveGrant(ctx, g)
}

func accountNote(a string) string {
	if a == "" {
		return ""
	}
	return " (account " + a + ")"
}

// revoke revokes g's tokens at the vendor (RFC 7009) and deletes them
// here, whatever the vendor answers. It returns what the vendor did not
// confirm.
func (s *Server) revoke(ctx context.Context, g *Grant) []string {
	var problems []string
	if tok, _, err := s.openToken(ctx, g); err == nil {
		for _, t := range []struct{ tok, hint string }{{tok.Refresh, "refresh_token"}, {tok.Access, "access_token"}} {
			if t.tok == "" {
				continue
			}
			if err := s.OAuth.Revoke(ctx, g.Vendor, g.ClientID, t.tok, t.hint); err != nil {
				problems = append(problems, "the vendor did not confirm revoking the "+t.hint+": "+err.Error())
			}
		}
	}
	s.dropToken(ctx, g)
	return problems
}

// undo removes what connect registered and revokes the token.
func (s *Server) undo(ctx context.Context, g *Grant, registered []string, reason string) error {
	if key, err := s.connectorKey(ctx, g.Tenant); err == nil {
		for _, name := range registered {
			_ = s.KIFF.Disconnect(ctx, key, name)
		}
	}
	s.revoke(ctx, g)
	g.Transport = ""
	return s.failGrant(ctx, g, reason)
}

func (s *Server) failGrant(ctx context.Context, g *Grant, reason string) error {
	g.Status, g.Reason = StatusFailed, reason
	g.event(s.now(), "failed: "+reason)
	if err := s.saveGrant(ctx, g); err != nil {
		return err
	}
	return errors.New(reason)
}

// disconnect is the owner's signed act: the tools leave KIFF, the token
// is revoked at the vendor and deleted here.
func (s *Server) disconnect(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Status == StatusDisconnected {
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
	var problems []string
	if key, err := s.connectorKey(r.Context(), g.Tenant); err != nil {
		problems = append(problems, "the KIFF connector key could not be opened")
	} else {
		for _, t := range g.Tools {
			if err := s.KIFF.Disconnect(r.Context(), key, t.Name); err != nil {
				problems = append(problems, "KIFF kept "+t.Name+": "+err.Error())
			}
		}
	}
	problems = append(problems, s.revoke(r.Context(), g)...)
	g.Status, g.Transport = StatusDisconnected, ""
	what := "disconnected " + sr.approval.Who() + ": token deleted"
	if len(problems) > 0 {
		what += "; " + strings.Join(problems, "; ")
	}
	g.event(s.now(), what)
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, s.view(g))
}

func (s *Server) grants(ctx context.Context, tenant string) ([]grantView, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	out := []grantView{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil {
			continue
		}
		g.ver = it.Ver
		s.expire(ctx, &g)
		out = append(out, s.view(&g))
	}
	return out, nil
}

// listGrants is every grant and its signed approval, for the owner.
func (s *Server) listGrants(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	gs, err := s.grants(r.Context(), sr.tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"grants": gs})
}

// export is the human-oversight record as a file: each grant, what the
// owner read, and the phone's signature with what is needed to check it
// without planque.
func (s *Server) export(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	gs, err := s.grants(r.Context(), sr.tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="planque-oversight-`+s.now().UTC().Format("2006-01-02")+`.json"`)
	writeJSON(w, http.StatusOK, map[string]any{
		"format":      "planque.oversight.v1",
		"tenant":      sr.tenant,
		"exported_at": s.now().UTC(),
		"how_to_verify": "Each approval.message is the exact text the phone signed: its sixth line is the SHA-256 of approval.body, " +
			"whose digest is the SHA-256 of the RFC 8785 form of request. approval.signature is a DER ECDSA P-256 signature " +
			"over SHA-256(message), by approval.device_key (SPKI DER).",
		"grants": gs,
	})
}

// MaxLabel bounds a name for people that a vendor gives: a name, not a text.
const MaxLabel = 100

// label is a vendor-given name for people, on one line and bounded; "" when
// it has a control character or is too long, so the caller's fallback shows.
func label(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || len([]rune(s)) > MaxLabel || strings.ContainsFunc(s, unicode.IsControl) {
		return ""
	}
	return s
}

// vendorName is the vendor as the owner reads it: its own name, else the
// host of its server.
func vendorName(g *Grant) string {
	if g.VendorName != "" {
		return g.VendorName
	}
	return hostOf(g.Request.Server)
}
