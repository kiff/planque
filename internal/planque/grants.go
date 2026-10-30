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
	"strconv"
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

// freeNames gives each of a grant's tools a name nothing else holds, then
// resolves their unit lookups. A name taken already, by another of the
// account's connections or by any tool in KIFF, would be refused by KIFF
// only at connect, after the owner approved and signed in (kiff/tafs#138,
// the owner's iPad, 21 Oct). A tool's name is KIFF's, not the owner's (they
// read its title), so planque gives it a free one instead, for any vendor;
// the asking agent reads the names back from the grant. own are the grant's
// tools already in KIFF, whose names stay its own. Both ways of naming tools
// come here: a request with its tools, and tools set after the listing.
func (s *Server) freeNames(ctx context.Context, tenant, grant string, own []GrantedTool, tools []ToolRequest) error {
	taken, err := s.namesInUse(ctx, tenant, grant)
	if err != nil {
		return err
	}
	for _, name := range s.namesInKIFF(ctx, tenant) {
		if !slices.ContainsFunc(own, func(t GrantedTool) bool { return t.Name == name }) {
			taken[name] = "KIFF"
		}
	}
	renameTaken(tools, grant, taken)
	resolveLookups(tools)
	return nil
}

// renameTaken gives each tool whose name is taken a free one, made from it
// and the grant (toolKey), and points the unit lookups that named it there.
func renameTaken(tools []ToolRequest, grant string, taken map[string]string) {
	renamed := map[string]string{}
	inRequest := func(name string) bool {
		return slices.ContainsFunc(tools, func(t ToolRequest) bool { return t.Name == name })
	}
	for i := range tools {
		if _, ok := taken[tools[i].Name]; !ok {
			continue
		}
		fresh := toolKey(tools[i].Name, grant)
		for n := 2; ; n++ {
			if _, ok := taken[fresh]; !ok && !inRequest(fresh) {
				break
			}
			fresh = numberedKey(tools[i].Name, grant, n)
		}
		renamed[tools[i].Name], tools[i].Name = fresh, fresh
		taken[fresh] = ""
	}
	for i := range tools {
		if l := tools[i].UnitLookup; l != nil {
			if to, ok := renamed[l.Tool]; ok {
				c := *l
				c.Tool = to
				tools[i].UnitLookup = &c
			}
		}
	}
}

// numberedKey is toolKey with n after it, its stem cut so the whole stays
// within the 64 characters a name may have, however many digits n takes.
func numberedKey(name, grant string, n int) string {
	num := strconv.Itoa(n)
	if max := 59 - len(num); len(name) > max {
		name = name[:max]
	}
	return toolKey(name, grant) + num
}

// namesInKIFF are the names of the tools connected in KIFF for the account,
// whoever connected them. Best effort: without them, only planque's own
// connections are known, and KIFF still refuses a taken name at connect.
func (s *Server) namesInKIFF(ctx context.Context, tenant string) []string {
	key, err := s.connectorKey(ctx, tenant)
	if err != nil {
		return nil
	}
	names, err := s.KIFF.Connections(ctx, key)
	if err != nil {
		s.log().WarnContext(ctx, "KIFF's connected tools couldn't be read", "tenant", tenant, "err", err.Error())
		return nil
	}
	return names
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
	if err := s.freeNames(r.Context(), tenant, id, nil, tools); err != nil {
		s.fail(w, r, err)
		return
	}
	g := &Grant{ID: id, Tenant: tenant, Status: StatusPending, Vendor: meta,
		Request: Request{Version: "planque.grant.v1", Grant: id, Tenant: tenant, Server: server, SignInAt: meta.Issuer,
			Agent: req.Agent, Tools: tools, AskedBy: asked, CreatedAt: s.now().UTC().Format(time.RFC3339)},
		Audience: s.BaseURL + "/mcp/" + url.PathEscape(tenant) + "/" + id}
	g.Text = grantText(g.Request)
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
	// A title names its tool on the owner's approval and Card: two alike
	// would read as one (kiff/tafs#138).
	titled := map[string]string{}
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
		case t.Title != "" && label(t.Title) != t.Title:
			return "a tool's title is what it does in the owner's words, on one line, at most 100 characters: " + name
		case t.Title != "" && titled[strings.ToLower(t.Title)] != "":
			return "tool " + name + " has the same title as " + titled[strings.ToLower(t.Title)] + ": give each tool its own"
		case t.SetupReads && !t.ReadOnly:
			return "setup_reads is only for a read-only tool: " + name
		}
		if t.Title != "" {
			titled[strings.ToLower(t.Title)] = name
		}
		if msg := checkNarrowing(t); msg != "" {
			return msg
		}
		seen[name] = true
	}
	return checkLookups(tools)
}

// grantText is what the owner reads before approving.
func grantText(q Request) string { return requestText(q, false, "") }

// requestText is what the owner reads. held is for tools asked for on a
// grant whose sign-in planque already holds: there is no second sign-in,
// and the text names the account it acts in, when the vendor names it.
func requestText(q Request, held bool, account string, catalog ...CatalogTool) string {
	if len(q.Tools) == 0 {
		return fmt.Sprintf("Connect your account at %s to see what it offers for the taf %s.\n"+
			"You will sign in at %s.\n\nOnly what it offers is read: its list of tools. None of them is called: "+
			"nothing is read from your account, nothing is changed, and nothing is given to the taf. "+
			"The tools the taf may use come next, and you approve them again. If they aren't approved within a day, the sign-in is revoked.\n\n"+
			"Your sign-in is kept in a sealed vault; the taf never sees it. You can disconnect at any time.",
			hostOf(q.Server), q.Agent, hostOf(q.SignInAt))
	}
	var b strings.Builder
	switch {
	case held && account != "":
		fmt.Fprintf(&b, "You're already signed in, to the account %s: no second sign-in.\n\nThe taf may:\n", account)
	case held:
		b.WriteString("You're already signed in: no second sign-in.\n\nThe taf may:\n")
	default:
		fmt.Fprintf(&b, "You will sign in at %s.\n\nThe taf may:\n", hostOf(q.SignInAt))
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
		// In words: the exact terms are in Details (requestTerms), signed all the same.
		fmt.Fprintf(&b, "• %s: %s\n", toolTitle(t, catalog...), how)
	}
	// Setup's own reads go through this connection, not KIFF: said in so
	// many words, once for all of them, and never naming planque (#82).
	var reads []string
	for _, t := range q.Tools {
		if t.SetupReads {
			reads = append(reads, toolTitle(t, catalog...))
		}
	}
	switch len(reads) {
	case 0:
	case 1:
		fmt.Fprintf(&b, "\nSetup reads, to set this up: %s. It only reads it, through KIFF with its own Card, at most %d reads a day; each read has KIFF's permit and is kept in KIFF's record.\n", reads[0], SetupReadsPerDay)
	default:
		fmt.Fprintf(&b, "\nSetup reads, to set this up: %s and %s. It only reads them, through KIFF with its own Card, at most %d reads a day; each read has KIFF's permit and is kept in KIFF's record.\n",
			strings.Join(reads[:len(reads)-1], ", "), reads[len(reads)-1], SetupReadsPerDay)
	}
	b.WriteString("\nYour sign-in is kept in a sealed vault; the taf never sees it. Each call is sent only with KIFF's permit for that exact call. You can disconnect at any time.")
	return b.String()
}

// requestTerms are the exact terms of each tool asked for, one line each,
// for the owner who opens Details: its name, the vendor's tool, what is
// pinned (with the vendor's words for it) and what the taf sets. What the
// phone signs is the request itself, these terms included.
func requestTerms(q Request, catalog ...CatalogTool) []string {
	out := make([]string, 0, len(q.Tools))
	for _, t := range q.Tools {
		line := t.Name
		if t.narrowed() {
			line += narrowingText(t, catalog...)
		} else {
			line += " (" + t.Tool + ")"
		}
		if t.SetupReads {
			line += "; Setup reads it too"
		}
		out = append(out, line)
	}
	return out
}

// additionText is what the owner reads first when tools are added to a
// connection the taf already has: what's new, and that what it does today
// keeps working while they decide.
func additionText(added, has []ToolRequest, catalog []CatalogTool) string {
	names := make([]string, 0, len(added))
	for _, t := range added {
		names = append(names, toolTitle(t, catalog...))
	}
	// A connection that holds only the asking agent's own reads gives the taf nothing
	// yet: these are its first tools, never an addition to what it has (kiff/tafs#236).
	if !slices.ContainsFunc(has, func(t ToolRequest) bool { return !t.SetupReads }) {
		return "These are the taf's first tools on this connection: " + strings.Join(names, ", ") +
			". Setup's own reads keep working while you decide.\n\n"
	}
	return "This adds to a connection the taf already has: " + strings.Join(names, ", ") +
		". What it does today keeps working while you decide.\n\n"
}

// toolTitle names a tool for the owner: the title the asking agent gave it,
// else the vendor's own label, else its vendor tool's name, never the name
// planque or the agent gave it (a grant prefix the owner can't read).
func toolTitle(t ToolRequest, catalog ...CatalogTool) string {
	if t.Title != "" {
		return t.Title
	}
	if i := slices.IndexFunc(catalog, func(c CatalogTool) bool { return c.Name == t.Tool }); i >= 0 && catalog[i].Title != "" {
		return catalog[i].Title
	}
	return t.Tool
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
	Vendor string `json:"vendor"`
	Status string `json:"status"`
	// StatusWords is Status in the owner's words, so every app shows the
	// same line for a connection and none keeps its own.
	StatusWords string `json:"status_words"`
	Text        string `json:"text"`
	// Terms are each tool's exact terms, for the owner who opens Details:
	// the text names the tools in words (kiff/tafs#138, item 52).
	Terms    []string      `json:"terms,omitempty"`
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
	v := grantView{ID: g.ID, Vendor: vendorName(g), Status: g.Status, StatusWords: statusWords(g), Text: g.Text, Terms: requestTerms(g.Request, g.Catalog...), Digest: g.Digest, Request: g.Request,
		Account: g.Account, Audience: g.Audience, Reason: g.Reason, Events: ownerEvents(g),
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

// statusWords is a grant's status as the owner reads it, one line for any app.
func statusWords(g *Grant) string {
	switch g.Status {
	case StatusConnected:
		return "Connected"
	case StatusPending:
		switch {
		case g.Connected != nil:
			// Tools added to a connected grant: its tools keep working.
			return "Connected · new tools waiting for your approval"
		case g.SignedIn:
			return "Signed in · tools waiting for your approval"
		}
		return "Waiting for your approval"
	case StatusListed:
		return "Signed in · the tools are being chosen"
	case StatusNeedsReconnect:
		return "Signed out by " + vendorName(g) + " · sign in again"
	case StatusApproved:
		return "Approved · waiting for the sign-in at " + vendorName(g)
	case StatusSignedIn:
		return "Signed in · waiting for you to confirm"
	case StatusFailed:
		return "Not connected"
	case StatusDisconnected:
		return "Disconnected"
	}
	return g.Status
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
	_, _ = s.expire(r.Context(), g)
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
	_, _ = s.expire(r.Context(), g)
	if s.overdue(g) {
		writeErr(w, http.StatusConflict, "this grant's deadline passed: its sign-in is being revoked")
		return
	}
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
	_, _ = s.expire(r.Context(), g)
	if s.overdue(g) {
		writeErr(w, http.StatusConflict, "this grant's deadline passed: its sign-in is being revoked")
		return
	}
	// Listed, or tools set and not approved yet (the agent may revise them),
	// or connected: then the tools are added to those it has, so a taf
	// keeps one connection per vendor, and the owner approves the addition
	// once, with no second sign-in (kiff/tafs#138, the owner's retest).
	adding := g.Status == StatusConnected || (g.Status == StatusPending && g.Connected != nil)
	if !g.SignedIn || (g.Status != StatusListed && g.Status != StatusPending && g.Status != StatusConnected) {
		writeErr(w, http.StatusConflict, "only a listed or connected grant takes tools; this one is "+g.Status)
		return
	}
	if len(g.Catalog) == 0 {
		writeErr(w, http.StatusConflict, "this grant was asked for with its tools at once, so it has no catalog to choose from; connect first (request with no tools)")
		return
	}
	// What it has: on a connected grant, its tools; while an addition waits,
	// the tools it had before it (the addition is revised).
	var has []ToolRequest
	switch {
	case g.Status == StatusConnected:
		has = g.Request.Tools
	case adding:
		has = g.Connected
	}
	if len(req.Tools) == 0 || len(has)+len(req.Tools) > 20 {
		writeErr(w, http.StatusBadRequest, fmt.Sprintf("tools must name 1 to %d of the catalog's tools", 20-len(has)))
		return
	}

	if msg := checkTools(g.Request.Agent, req.Tools); msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	for i := range req.Tools {
		j := slices.IndexFunc(g.Catalog, func(c CatalogTool) bool { return c.Name == req.Tools[i].Tool })
		if j < 0 {
			writeErr(w, http.StatusBadRequest, hostOf(g.Request.Server)+" lists no tool named "+req.Tools[i].Tool)
			return
		}
		// Read-only is the vendor's word, never the asking agent's: the owner
		// reads "reads only", and KIFF treats it as a read (kiff/agents#139).
		if req.Tools[i].ReadOnly && !g.Catalog[j].ReadOnly {
			writeErr(w, http.StatusBadRequest, notReadOnly(g, req.Tools[i].Tool))
			return
		}
		if req.Tools[i].Name == "" {
			req.Tools[i].Name = toolKey(req.Tools[i].Tool, g.ID)
		}
		if n := req.Tools[i].Name; slices.ContainsFunc(has, func(h ToolRequest) bool { return h.Name == n }) {
			writeErr(w, http.StatusConflict, "this grant already has "+n+": send only the tools to add; what it has stays as it is")
			return
		}
	}
	if err := s.freeNames(r.Context(), tenant, g.ID, g.Tools, req.Tools); err != nil {
		s.fail(w, r, err)
		return
	}
	if !adding {
		g.Request.Tools = req.Tools
		g.Text = requestText(g.Request, true, g.Account, g.Catalog...)
		g.Status, g.Approval = StatusPending, nil
		g.event(s.now(), "tools set by "+me.AgentID+"; waiting for the owner's approval")
	} else {
		if g.Status == StatusConnected {
			g.Connected = slices.Clone(g.Request.Tools)
		}
		g.Request.Tools = append(slices.Clone(g.Connected), req.Tools...)
		// The owner reads what's added first, then the whole connection, as
		// it will be: the approval signs all of it. The earlier approval
		// stays on the grant until this one replaces it.
		g.Text = additionText(req.Tools, g.Connected, g.Catalog) + requestText(g.Request, true, g.Account, g.Catalog...)
		g.Status, g.HeldUntil = StatusPending, s.now().Add(holdTTL).UTC()
		g.event(s.now(), fmt.Sprintf("%d tools added by %s; waiting for the owner's approval; its %d tools keep working meanwhile", len(req.Tools), me.AgentID, len(g.Connected)))
	}
	if g.Digest, err = digest(g.Request); err != nil {
		s.fail(w, r, err)
		return
	}
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
// expire ends g if its deadline passed: its tokens revoked at the vendor
// and dropped, then g saved as failed. ended says it did. An error leaves g
// as stored, its tokens kept until the vendor confirms revoking them, for
// the next sweep or read to try again (kiff/planque#23).
func (s *Server) expire(ctx context.Context, g *Grant) (ended bool, err error) {
	if g.Status == StatusPending && g.Connected != nil && !g.HeldUntil.IsZero() && s.now().After(g.HeldUntil) {
		// Tools added to a connected grant and never approved: they go, and
		// the grant stays connected with what it had. Its sign-in is in use.
		was := *g
		g.Request.Tools, g.Connected = g.Connected, nil
		g.Text = requestText(g.Request, true, g.Account, g.Catalog...)
		if g.Digest, err = digest(g.Request); err != nil {
			*g = was
			return false, err
		}
		g.Status, g.HeldUntil = StatusConnected, time.Time{}
		g.event(s.now(), "the added tools were not approved within a day; what it had stays connected")
		if err := s.saveGrant(ctx, g); err != nil {
			*g = was
			return false, err
		}
		return false, nil
	}
	var reason string
	held := g.SignedIn && (g.Status == StatusListed || g.Status == StatusPending)
	switch {
	case held && !g.HeldUntil.IsZero() && s.now().After(g.HeldUntil):
		reason = "its tools were not approved within a day; the sign-in was revoked"
	case g.Status == StatusSignedIn && g.Reconnect && s.now().After(g.ConfirmBy):
		// A new sign-in for a grant its vendor signed out, never confirmed:
		// that token goes, and the grant waits for a sign-in again.
		if err := s.revokeNow(ctx, g); err != nil {
			return false, err
		}
		was := *g
		g.ConfirmHash, g.ConfirmBy = "", time.Time{}
		g.Status, g.Reconnect = StatusNeedsReconnect, false
		g.event(s.now(), "the new sign-in was not confirmed in time; still waiting for the owner to sign in again")
		if err := s.saveGrant(ctx, g); err != nil {
			*g = was
			return false, err
		}
		return false, nil
	case g.Status == StatusSignedIn && s.now().After(g.ConfirmBy):
		reason = "the sign-in was not confirmed in time; its token was revoked"
	default:
		return false, nil
	}
	if err := s.revokeNow(ctx, g); err != nil {
		return false, err
	}
	was := *g
	g.HeldUntil, g.ConfirmHash, g.ConfirmBy = time.Time{}, "", time.Time{}
	g.Status, g.Reason = StatusFailed, reason
	g.event(s.now(), "failed: "+reason)
	if err := s.saveGrant(ctx, g); err != nil {
		*g = was
		return false, err
	}
	s.log().InfoContext(ctx, "expired sign-in revoked", slog.String("grant", g.ID), slog.String("reason", reason))
	return true, nil
}

// revokeNow revokes g's tokens at the vendor, then drops them. A token
// planque can't open, or one the vendor didn't confirm revoking, is kept,
// so revoking can be tried again: dropping it would leave it live at the
// vendor with no way to revoke it.
func (s *Server) revokeNow(ctx context.Context, g *Grant) error {
	tok, _, err := s.openToken(ctx, g)
	if errors.Is(err, errNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open the sign-in's token to revoke it: %w", err)
	}
	for _, t := range []struct{ tok, hint string }{{tok.Refresh, "refresh_token"}, {tok.Access, "access_token"}} {
		if t.tok == "" {
			continue
		}
		if err := s.OAuth.Revoke(ctx, g.Vendor, g.ClientID, t.tok, t.hint); err != nil {
			return fmt.Errorf("the vendor did not confirm revoking the %s: %w", t.hint, err)
		}
	}
	if _, _, err := s.KV.Delete(ctx, tenantPK(g.Tenant), skToken(g.ID)); err != nil {
		return fmt.Errorf("drop the revoked token: %w", err)
	}
	return nil
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
	if g.Status != StatusPending && g.Status != StatusApproved && g.Status != StatusSignedIn && g.Status != StatusNeedsReconnect {
		writeErr(w, http.StatusConflict, "this grant is "+g.Status)
		return
	}
	_, _ = s.expire(r.Context(), g)
	if g.Status == StatusFailed {
		writeErr(w, http.StatusConflict, "this grant is failed: "+g.Reason)
		return
	}
	if s.overdue(g) {
		writeErr(w, http.StatusConflict, "this grant's deadline passed: its sign-in is being revoked")
		return
	}

	if g.SignedIn && g.Status == StatusPending {
		// Tools set on a grant whose sign-in is held: approving them
		// connects at once, with no second sign-in.
		a := sr.approval
		tok, _, err := s.openToken(r.Context(), g)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if g.Connected != nil {
			// The approval of the addition replaces the grant's only once it
			// connects: a failed one leaves the approval of what stays.
			if err := s.addTools(r.Context(), g, tok.Access, a); err != nil {
				s.log().WarnContext(r.Context(), "adding tools failed", slog.String("grant", g.ID), slog.String("err", err.Error()))
				writeErr(w, http.StatusBadGateway, err.Error())
				return
			}
			writeJSON(w, http.StatusOK, s.view(g))
			return
		}
		g.Approval = &a
		g.event(s.now(), "tools approved "+a.Who())
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
	switch {
	case g.Status == StatusPending:
		g.Approval = &a
		g.event(s.now(), "approved "+a.Who())
	case g.Status == StatusNeedsReconnect:
		// The same tools, signed again: the sign-in that follows connects
		// the grant again, with no change in KIFF. This device confirms it.
		g.Approval, g.Reconnect = &a, true
		g.event(s.now(), "a new sign-in was approved "+a.Who()+" to reconnect")
	default:
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
	if g.Reconnect {
		// Reconnecting acts in the same vendor account, or not at all.
		if g.Account != "" && tok.Account != "" && tok.Account != g.Account {
			for _, t := range []struct{ tok, hint string }{{tok.Refresh, "refresh_token"}, {tok.Access, "access_token"}} {
				if t.tok != "" {
					_ = s.OAuth.Revoke(r.Context(), g.Vendor, g.ClientID, t.tok, t.hint)
				}
			}
			g.Status, g.Reconnect = StatusNeedsReconnect, false
			g.event(s.now(), "a sign-in to another account at "+vendorName(g)+" was refused; still waiting for the owner to sign in again")
			_ = s.saveGrant(r.Context(), g)
			s.page(w, http.StatusConflict, "Not connected", "That sign-in is for another account at "+vendorName(g)+". Sign in to the account this taf uses.")
			return
		}
		// The token the vendor ended goes; the new one takes its place.
		s.dropToken(r.Context(), g)
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
	g.event(s.now(), "Signed in to "+vendorName(g)+accountNote(g.Account)+"; waiting for the phone to confirm")
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
		_, _ = s.expire(r.Context(), g)
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
	if g.Reconnect {
		// Signed in again: the same tools, still in KIFF, now with this token.
		g.Reconnect, g.Reason = false, ""
		g.Status, g.HeldUntil = StatusConnected, time.Time{}
		g.event(s.now(), "connected again with the new sign-in; its tools are as they were")
		if err := s.saveGrant(r.Context(), g); err != nil {
			s.fail(w, r, err)
			return
		}
		s.log().InfoContext(r.Context(), "grant reconnected", slog.String("grant", g.ID))
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
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

// readOnlyHint is whether the vendor marks t read-only (MCP's readOnlyHint).
func readOnlyHint(t mcpTool) bool {
	var hints struct {
		ReadOnly bool `json:"readOnlyHint"`
	}
	_ = json.Unmarshal(t.Annotations, &hints)
	return hints.ReadOnly
}

func notReadOnly(g *Grant, tool string) string {
	return tool + " isn't read-only at " + hostOf(g.Request.Server) + ": the vendor doesn't mark it so, so it can't be granted as a read"
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
		g.Catalog = append(g.Catalog, CatalogTool{Name: t.Name, Title: label(t.Label()), Description: t.Description,
			InputSchema: t.InputSchema, ReadOnly: readOnlyHint(t)})
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
		gt, why := granted(g, t, listed)
		if why != "" {
			s.revoke(ctx, g)
			return s.failGrant(ctx, g, why)
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
	var done, titles []string
	for _, t := range g.Tools {
		err := s.KIFF.Connect(ctx, key, kiff.Connection{Name: t.Name, Title: t.Title, URL: g.Audience, Tool: t.Name,
			Description: t.Description, AmountArg: t.AmountArg, AmountUnit: t.AmountUnit, AmountScale: t.AmountScale,
			IdempotencyArg: IdempotencyArg, ReadOnly: t.ReadOnly, Mode: "relay", Audience: g.Audience, Credential: transport,
			ResourceArg: t.ResourceArg, UnitLookup: t.UnitLookup, VendorTool: t.Tool})
		if err != nil {
			return s.undo(ctx, g, done, "KIFF did not register "+t.Name+": "+err.Error())
		}
		done = append(done, t.Name)
		titles = append(titles, toolTitle(t.ToolRequest, CatalogTool{Name: t.Tool, Title: t.Title}))
	}
	g.event(s.now(), "Connected in KIFF: "+strings.Join(titles, ", "))
	s.setupReads(ctx, key, g)
	return s.saveGrant(ctx, g)
}

// granted is t as the vendor lists it, to register in KIFF; why says why it
// can't be.
func granted(g *Grant, t ToolRequest, listed []mcpTool) (GrantedTool, string) {
	i := slices.IndexFunc(listed, func(l mcpTool) bool { return l.Name == t.Tool })
	if i < 0 {
		names := make([]string, 0, len(listed))
		for _, l := range listed {
			names = append(names, l.Name)
		}
		return GrantedTool{}, hostOf(g.Request.Server) + " has no tool named " + t.Tool + "; it lists: " + strings.Join(names, ", ")
	}
	// A tool asked for as read-only connects only if the vendor marks it
	// so (kiff/agents#139); the owner approved "reads only".
	if t.ReadOnly && !readOnlyHint(listed[i]) {
		return GrantedTool{}, notReadOnly(g, t.Tool)
	}
	// KIFF's Card and record name it by the asking agent's title, else the vendor's label.
	title := t.Title
	if title == "" {
		title = label(listed[i].Label())
	}
	gt := GrantedTool{ToolRequest: t, Title: title, Description: listed[i].Description, InputSchema: listed[i].InputSchema}
	if t.narrowed() {
		gt.InputSchema = t.schema()
		gt.Description = "Through planque: " + t.Tool + narrowingText(t) + "."
	}
	return gt, ""
}

// addTools connects the tools added to a connected grant, once the owner
// approved them: only those are registered in KIFF, with a transport
// credential of their own, and the tools it had stay as they were. If
// anything fails, the grant goes back to connected with those alone, and
// its sign-in stays: an addition never takes a taf's tools away.
func (s *Server) addTools(ctx context.Context, g *Grant, access string, a Approval) error {
	had, hadTools, hadTransports, hadApproval := g.Connected, g.Tools, g.Transports, g.Approval
	g.Approval = &a
	g.event(s.now(), "added tools approved "+a.Who())
	back := func(registered []string, reason string) error {
		if key, err := s.connectorKey(ctx, g.Tenant); err == nil {
			for _, name := range registered {
				_ = s.KIFF.Disconnect(ctx, key, name)
			}
		}
		g.Request.Tools, g.Connected, g.Tools, g.Transports, g.Approval = had, nil, hadTools, hadTransports, hadApproval
		g.Text = requestText(g.Request, true, g.Account, g.Catalog...)
		if d, err := digest(g.Request); err == nil {
			g.Digest = d
		}
		g.Status, g.HeldUntil = StatusConnected, time.Time{}
		g.event(s.now(), "the added tools weren't connected: "+reason+"; what it had stays connected")
		if err := s.saveGrant(ctx, g); err != nil {
			return err
		}
		return errors.New(reason)
	}
	listed, _, err := s.MCP.ListTools(ctx, g.Request.Server, access)
	if err != nil {
		return back(nil, "Couldn't list the tools at "+hostOf(g.Request.Server)+": "+err.Error())
	}
	var added []GrantedTool
	for _, t := range g.Request.Tools[len(had):] {
		gt, why := granted(g, t, listed)
		if why != "" {
			return back(nil, why)
		}
		added = append(added, gt)
	}
	transport := oauth.Random(32)
	sum := sha256.Sum256([]byte(transport))
	// As in connect: saved as connected with the new tools and credential
	// first, so the relay answers KIFF's probe while they're registered.
	g.Tools = append(slices.Clone(hadTools), added...)
	g.Transports = append(slices.Clone(hadTransports), hex.EncodeToString(sum[:]))
	g.Connected, g.Status, g.HeldUntil = nil, StatusConnected, time.Time{}
	if err := s.saveGrant(ctx, g); err != nil {
		g.Connected, g.Tools, g.Transports, g.Approval = had, hadTools, hadTransports, hadApproval
		return err
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		return back(nil, "the KIFF connector key could not be opened")
	}
	var done, titles []string
	for _, t := range added {
		err := s.KIFF.Connect(ctx, key, kiff.Connection{Name: t.Name, Title: t.Title, URL: g.Audience, Tool: t.Name,
			Description: t.Description, AmountArg: t.AmountArg, AmountUnit: t.AmountUnit, AmountScale: t.AmountScale,
			IdempotencyArg: IdempotencyArg, ReadOnly: t.ReadOnly, Mode: "relay", Audience: g.Audience, Credential: transport,
			ResourceArg: t.ResourceArg, UnitLookup: t.UnitLookup, VendorTool: t.Tool})
		if err != nil {
			return back(done, "KIFF did not register "+t.Name+": "+err.Error())
		}
		done = append(done, t.Name)
		titles = append(titles, toolTitle(t.ToolRequest, CatalogTool{Name: t.Tool, Title: t.Title}))
	}
	g.event(s.now(), "Added in KIFF: "+strings.Join(titles, ", "))
	s.setupReads(ctx, key, g)
	return s.saveGrant(ctx, g)
}

// setupReads asks KIFF for the Card the asking agent (Setup) reads g's
// setup_reads tools on, with the owner's signed approval of g, which KIFF
// checks itself (kiff/tafs#138). Its reads then go through KIFF's gateway
// like any call. Without that Card each read is refused, so a failure is
// on the grant's record; the connection itself stands.
func (s *Server) setupReads(ctx context.Context, key string, g *Grant) {
	var tools []string
	for _, t := range g.Tools {
		if t.SetupReads && t.ReadOnly {
			tools = append(tools, t.Name)
		}
	}
	if len(tools) == 0 {
		return
	}
	if g.Approval == nil || g.Approval.DeviceID == "" {
		g.event(s.now(), "Setup's reads weren't set up in KIFF: no approval from a device KIFF knows")
		return
	}
	req := kiff.SetupReads{Agent: g.Request.AskedBy, Grant: g.ID, Tools: tools}
	req.Approval.DeviceID, req.Approval.Message, req.Approval.Body, req.Approval.Signature = g.Approval.DeviceID, g.Approval.Message, g.Approval.Body, g.Approval.Signature
	if err := s.KIFF.IssueSetupReads(ctx, key, req); err != nil {
		s.log().WarnContext(ctx, "setup reads not issued", slog.String("grant", g.ID), slog.String("err", err.Error()))
		g.event(s.now(), "Setup's reads weren't set up in KIFF: "+err.Error())
		return
	}
	g.event(s.now(), fmt.Sprintf("Setup reads through KIFF, on its own Card: at most %d reads a day", SetupReadsPerDay))
}

// SetupReadsPerDay is what KIFF's Card for Setup's reads allows
// (kiff-cloud SetupReadsPerDay), for the owner's text.
const SetupReadsPerDay = 50

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
	g.Transport, g.Transports, g.Connected = "", nil, nil
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
	g.Status, g.Transport, g.Transports, g.Connected = StatusDisconnected, "", nil, nil
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
		_, _ = s.expire(ctx, &g)
		out = append(out, s.view(&g))
	}
	return out, nil
}

// namesInUse are the tool names the account's other grants hold, each with
// its vendor's host: a connected one's tools as registered in KIFF, and
// before that, at any step on the way (listed with tools set, pending,
// approved, signed in), the tools it asks for, which it registers next. A
// failed or disconnected grant holds none.
func (s *Server) namesInUse(ctx context.Context, tenant, except string) (map[string]string, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.ID == except {
			continue
		}
		switch g.Status {
		case StatusFailed, StatusDisconnected:
		case StatusConnected:
			for _, t := range g.Tools {
				out[t.Name] = hostOf(g.Request.Server)
			}
		default:
			for _, t := range g.Request.Tools {
				out[t.Name] = hostOf(g.Request.Server)
			}
		}
	}
	return out, nil
}

// agentGrants are the grants the asking agent asked for a taf's agent, by
// its key: so Setup, improving a taf, finds the connection it already has
// to a vendor and adds to it, instead of asking for another (kiff/tafs#138).
// Only what Setup needs: no approval, no text.
func (s *Server) agentGrants(w http.ResponseWriter, r *http.Request) {
	s.askedGrants(w, r, r.PathValue("agent"))
}

// everyAskedGrant is every grant the asking agent asked for, for any taf's
// agent, each naming its agent: so Setup that lost a taf's agent id finds its
// own connection again instead of asking the owner to connect once more
// (kiff/planque#49).
func (s *Server) everyAskedGrant(w http.ResponseWriter, r *http.Request) {
	s.askedGrants(w, r, "")
}

// askedGrants lists the grants the asking agent asked for: for agent, or for
// any agent when it is empty.
func (s *Server) askedGrants(w http.ResponseWriter, r *http.Request, agent string) {
	tenant := r.PathValue("tenant")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	items, err := s.KV.Query(r.Context(), tenantPK(tenant), "GRANT#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	type tool struct {
		Name     string `json:"name"`
		Tool     string `json:"tool"`
		Title    string `json:"title,omitempty"`
		ReadOnly bool   `json:"read_only,omitempty"`
	}
	type grant struct {
		ID      string `json:"id"`
		Agent   string `json:"agent"`
		Vendor  string `json:"vendor"`
		Server  string `json:"server"`
		Status  string `json:"status"`
		Account string `json:"account,omitempty"`
		Created string `json:"created,omitempty"`
		Tools   []tool `json:"tools"`
	}
	out := []grant{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || (agent != "" && g.Request.Agent != agent) || g.Request.AskedBy != me.AgentID {
			continue
		}
		g.ver = it.Ver
		_, _ = s.expire(r.Context(), &g)
		if g.Status == StatusFailed || g.Status == StatusDisconnected {
			continue
		}
		v := grant{ID: g.ID, Agent: g.Request.Agent, Vendor: vendorName(&g), Server: g.Request.Server, Status: g.Status, Account: g.Account,
			Created: g.Request.CreatedAt, Tools: []tool{}}
		for _, t := range g.Request.Tools {
			v.Tools = append(v.Tools, tool{Name: t.Name, Tool: t.Tool, Title: t.Title, ReadOnly: t.ReadOnly})
		}
		out = append(out, v)
	}
	writeJSON(w, http.StatusOK, map[string]any{"grants": out})
}

// listGrants is the account's grants and their signed approvals, for the
// owner: those still in use, so the list doesn't grow with every
// connection ever made. A disconnected grant's record stays until the
// account ends (its vendor tokens are gone at disconnect), and
// ?include=disconnected lists it too, as the export does.
func (s *Server) listGrants(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	gs, err := s.grants(r.Context(), sr.tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if r.URL.Query().Get("include") != "disconnected" {
		gs = slices.DeleteFunc(gs, func(g grantView) bool { return g.Status == StatusDisconnected })
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
			"over SHA-256(message), by the key of the device approval.device_id names, which KIFF keeps: KIFF checked it " +
			"when the approval was made, and checks it again on request (POST /v1/devices/verify). An approval made before " +
			"kiff/planque#12 may instead carry approval.device_key (SPKI DER), the key itself.",
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

// Events written before #16 named who signed by ids: "by user_… on device
// dev_…" (#11) or "on the phone ab12 cd34 ef56 7890". These match only
// those old words; new events never carry them.
var (
	legacyDevice = regexp.MustCompile(` by \S+ on device (dev_[a-z0-9]+)`)
	legacyPhone  = regexp.MustCompile(` on the phone [0-9a-f]{4}( [0-9a-f]{4}){3}`)
)

// ownerEvents is g's record as the owner reads it: an event written before
// #16 says who signed in today's words (the device's name from the grant's
// approval when it is that device, else the same fallbacks as Who), with the
// stored line kept beside it as Recorded. Storage is never rewritten.
func ownerEvents(g *Grant) []Event {
	out := make([]Event, len(g.Events))
	for i, e := range g.Events {
		what := legacyDevice.ReplaceAllStringFunc(e.What, func(m string) string {
			id := legacyDevice.FindStringSubmatch(m)[1]
			if a := g.Approval; a != nil && a.DeviceID == id && a.DeviceName != "" {
				return " on " + a.DeviceName
			}
			return " on a device signed in to KIFF"
		})
		what = legacyPhone.ReplaceAllString(what, " on the enrolled phone")
		out[i] = e
		if what != e.What {
			out[i].What, out[i].Recorded = what, e.What
		}
	}
	return out
}

// setupRead was the asking agent reading a setup_reads tool through
// planque itself (kiff/tafs#174), with no permit. Those reads go through
// KIFF's gateway now, on Setup's own read Card (kiff/tafs#138), so this
// answers 410.
func (s *Server) setupRead(w http.ResponseWriter, r *http.Request) {
	// Setup's reads go through KIFF's gateway now, decided on its own read
	// Card, and the relay sends them only with KIFF's permit (kiff/tafs#138).
	writeErr(w, http.StatusGone, "Setup reads through KIFF's gateway now: call the tool there, by the name the grant lists")
}
