package planque

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/store"
)

// An API-key grant (kiff/planque#64) reaches a vendor that has a REST API
// and no MCP server. The grant names the vendor's published OpenAPI
// description; planque reads it once into a snapshot, which the owner
// approves by its digest and every call uses. The owner gives the key on
// their phone, straight to planque, with the approval: it is sealed like a
// sign-in's token, held as one (so adding tools needs no second key), never
// shown again, and deleted when the grant is disconnected.

// KindAPIKey is an API-key grant's kind.
const KindAPIKey = "api_key"

// KeyHeader carries the vendor's key to POST …/key. It's outside the signed
// body, which keeps only its SHA-256, so the approval kept as evidence never
// holds the key.
const KeyHeader = "Planque-Vendor-Key"

// snapshotChunk is how much of a compressed snapshot one record holds,
// under the table's item limit.
const snapshotChunk = 300 << 10

// saveSnapshot keeps a grant's snapshot, compressed, in as many records as
// it takes.
func (s *Server) saveSnapshot(ctx context.Context, tenant, grant string, snap *openapi.Snapshot) error {
	raw, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(raw); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	b := buf.Bytes()
	for n := 0; len(b) > 0 || n == 0; n++ {
		part := b[:min(len(b), snapshotChunk)]
		b = b[len(part):]
		doc, _ := json.Marshal(map[string]any{"part": part, "last": len(b) == 0})
		if err := s.KV.Put(ctx, store.Item{PK: tenantPK(tenant), SK: skSnapshot(grant, n), Doc: doc, Ver: 1}, store.Cond{}); err != nil {
			return err
		}
	}
	return nil
}

// loadSnapshot is the snapshot a grant was asked for with.
func (s *Server) loadSnapshot(ctx context.Context, tenant, grant string) (*openapi.Snapshot, error) {
	var buf bytes.Buffer
	for n := 0; ; n++ {
		it, ok, err := s.KV.Get(ctx, tenantPK(tenant), skSnapshot(grant, n))
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("the grant's API description isn't kept")
		}
		var rec struct {
			Part []byte `json:"part"`
			Last bool   `json:"last"`
		}
		if err := json.Unmarshal(it.Doc, &rec); err != nil {
			return nil, err
		}
		buf.Write(rec.Part)
		if rec.Last {
			break
		}
	}
	zr, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(io.LimitReader(zr, 64<<20))
	if err != nil {
		return nil, err
	}
	var snap openapi.Snapshot
	return &snap, json.Unmarshal(raw, &snap)
}

// dropSnapshot deletes a grant's snapshot records.
func (s *Server) dropSnapshot(ctx context.Context, tenant, grant string) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "SNAPSHOT#"+grant+"#")
	if err != nil {
		return
	}
	for _, it := range items {
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
	}
}

// apiAsk is what a grant request says about the API it calls.
type apiAsk struct {
	URL        string       `json:"url"`
	Base       string       `json:"base"`
	Auth       openapi.Auth `json:"auth"`
	Operations []string     `json:"operations,omitempty"`
}

// readAPI reads the description an API-key grant names, checks the tools
// asked for against it and keeps it as the grant's snapshot. msg is the
// refusal in words, for the agent; err is planque's own failure.
func (s *Server) readAPI(ctx context.Context, tenant, grant string, a apiAsk, tools []ToolRequest, expectedDigest ...string) (*APIRequest, *openapi.Snapshot, string, error) {
	a.URL, a.Base = strings.TrimSpace(a.URL), strings.TrimSpace(a.Base)
	a.Auth.In, a.Auth.Name = strings.TrimSpace(a.Auth.In), strings.TrimSpace(a.Auth.Name)
	if a.URL == "" {
		return nil, nil, "api.url is required: the vendor's OpenAPI description, in JSON or YAML", nil
	}
	var raw []byte
	var err error
	if a.URL == TelegramDescription {
		// A Telegram bot (kiff/planque#83): planque's own description of
		// the Bot API, its server and where the token goes.
		raw, a.Base, a.Auth, a.Operations = []byte(telegramAPI), s.telegramAPIBase(), telegramAuth, nil
	} else if raw, err = s.REST.Fetch(ctx, a.URL); err != nil {
		return nil, nil, "Couldn't read the API description at " + a.URL + ": " + err.Error(), nil
	}
	if len(expectedDigest) > 0 && expectedDigest[0] != "" {
		sum := sha256.Sum256(raw)
		if hex.EncodeToString(sum[:]) != expectedDigest[0] {
			return nil, nil, "the vendor description changed since this catalog entry was reviewed", nil
		}
	}
	var snap *openapi.Snapshot
	if len(a.Operations) > 0 {
		snap, err = openapi.ParseOperations(raw, a.Auth, a.Operations)
	} else {
		snap, err = openapi.Parse(raw, a.Auth)
	}
	if err != nil {
		return nil, nil, err.Error(), nil
	}
	if a.Base == "" && len(snap.Servers) > 0 {
		a.Base = snap.Servers[0]
	}
	if a.Base == "" {
		return nil, nil, "api.base is required: the description names no server", nil
	}
	if err := s.REST.CheckURL(a.Base); err != nil {
		return nil, nil, "api.base: " + err.Error(), nil
	}
	if msg := apiToolsCheck(snap, tools); msg != "" {
		return nil, nil, msg, nil
	}
	if err := s.saveSnapshot(ctx, tenant, grant, snap); err != nil {
		return nil, nil, "", err
	}
	return &APIRequest{URL: a.URL, Digest: snap.Digest, Base: strings.TrimRight(a.Base, "/"), Auth: a.Auth,
		Operations: slices.Clone(a.Operations)}, snap, "", nil
}

// apiToolsCheck says why tools can't be asked for from snap, "" when they can.
func apiToolsCheck(snap *openapi.Snapshot, tools []ToolRequest) string {
	for _, t := range tools {
		i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == t.Tool })
		if i < 0 {
			return "the API has no operation named " + t.Tool + "; its catalog lists them once the key is given"
		}
		if t.ReadOnly && !snap.Operations[i].ReadOnly {
			return t.Tool + " isn't a read: only GET and HEAD operations are read-only"
		}
	}
	return ""
}

// reusable is the connected grant a new one may reuse the sign-in of
// (kiff/planque#36): this account's, connected, with its sign-in held, and to
// the vendor the request names, if it names one. msg says why not.
func (s *Server) reusable(ctx context.Context, tenant, id, kind, server string) (*Grant, string, error) {
	g, err := s.loadGrant(ctx, tenant, id)
	if errors.Is(err, errNotFound) {
		return nil, "this account has no connection " + id + " to reuse", nil
	}
	if err != nil {
		return nil, "", err
	}
	switch {
	case g.Request.Connector != nil:
		return nil, "request a new connector grant or add tools to this connection", nil
	case g.Status != StatusConnected || !g.SignedIn:
		return nil, "connection " + id + " isn't connected now (" + g.Status + "), so its sign-in can't be reused: ask for a new one", nil
	case kind != "" && kind != g.Request.Kind:
		return nil, "connection " + id + " isn't an " + kind + " connection", nil
	case server != "" && g.Request.Kind != KindAPIKey && server != g.Request.Server:
		return nil, "connection " + id + " is to " + hostOf(g.Request.Server) + ", not " + hostOf(server), nil
	}
	if _, _, err := s.openToken(ctx, g); err != nil {
		return nil, "connection " + id + "'s sign-in isn't held any more: ask for a new one", nil
	}
	return g, "", nil
}

// errSignInRefused: the vendor no longer takes the sign-in a request asked to
// reuse, so reusing it would only fail once the owner approved.
var errSignInRefused = errors.New("planque: the vendor refused the sign-in to reuse")

// signInWorks checks, before a request is offered with "no second sign-in",
// that g's sign-in still works at its vendor: a live token (refreshed when
// due) and one read of its server's tools, the token refreshed once if the
// server refuses it. A refusal is errSignInRefused, and g is signed out as
// any refused refresh signs it out. A vendor that can't be reached isn't a
// refusal: reuse goes on, as it did. Only an OAuth sign-in to an MCP server
// is read; a key or a managed app's sign-in is checked where it's used.
func (s *Server) signInWorks(ctx context.Context, g *Grant) error {
	if g.Request.Connector != nil || g.Request.API != nil || g.Request.Kind == KindAPIKey || g.Request.Kind == KindBrowser {
		return nil
	}
	refused := func(err error) error {
		var so signedOutError
		if errors.As(err, &so) || errors.Is(err, errNoToken) {
			return errSignInRefused
		}
		return err
	}
	access, ver, err := s.accessToken(ctx, g, 0)
	if err != nil {
		return refused(err)
	}
	for attempt := 0; ; attempt++ {
		_, _, err = s.MCP.ListTools(ctx, g.Request.Server, access)
		if !errors.Is(err, mcp.ErrUnauthorized) {
			return nil
		}
		if attempt > 0 {
			return errSignInRefused
		}
		if access, ver, err = s.accessToken(ctx, g, ver); err != nil {
			return refused(err)
		}
	}
}

// vendorTools are the tools g's vendor offers, as an MCP server lists them:
// listed by its MCP server, or read from an API-key grant's snapshot.
func (s *Server) vendorTools(ctx context.Context, g *Grant, access string) ([]mcpTool, mcp.Server, error) {
	if g.Request.Connector != nil {
		return s.connectorTools(ctx, g)
	}
	if g.Request.Kind == KindBrowser {
		return browserMCPTools(), mcp.Server{Name: browserVendorName(g.Request.Browser)}, nil
	}
	if g.Request.API == nil {
		return s.MCP.ListTools(ctx, g.Request.Server, access)
	}
	snap, err := s.loadSnapshot(ctx, g.Tenant, g.ID)
	if err != nil {
		return nil, mcp.Server{}, err
	}
	out := make([]mcpTool, 0, len(snap.Operations))
	for _, o := range snap.Operations {
		ann, _ := json.Marshal(map[string]any{"readOnlyHint": o.ReadOnly})
		desc := o.Description
		if desc == "" {
			desc = o.Summary
		}
		out = append(out, mcpTool{Name: o.Name, Title: label(o.Summary), Description: o.Method + " " + o.Path + ". " + desc,
			InputSchema: o.InputSchema, Annotations: ann})
	}
	return out, mcp.Server{Name: hostOf(g.Request.API.Base), Title: snap.Title}, nil
}

// operations is a listed API-key grant's catalog with each tool's method and
// path, for what the owner reads.
func withOperations(g *Grant, snap *openapi.Snapshot) {
	for i, c := range g.Catalog {
		if j := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == c.Name }); j >= 0 {
			g.Catalog[i].Method, g.Catalog[i].Path = snap.Operations[j].Method, snap.Operations[j].Path
		}
	}
}

// giveKey is the owner's approval of an API-key grant, with the key: the
// phone signs {"digest", "key_sha256"} and sends the key in KeyHeader. The
// key is sealed and held like a sign-in's token. A grant with no tools is
// then listed; one with tools connects at once. A grant the vendor stopped
// taking the key for (needs_reconnect) takes a new one the same way, and a
// connected one takes a replacement once one read with it works (replaceKey).
func (s *Server) giveKey(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		Digest    string          `json:"digest"`
		KeySHA256 string          `json:"key_sha256"`
		Profile   *ProfileDigests `json:"kiff_profile,omitempty"`
	}
	if err := decodeStrict(sr.body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	key := strings.TrimSpace(r.Header.Get(KeyHeader))
	sum := sha256.Sum256([]byte(key))
	want, _ := hex.DecodeString(req.KeySHA256)
	if key == "" || len(want) != len(sum) || subtle.ConstantTimeCompare(want, sum[:]) != 1 {
		writeErr(w, http.StatusBadRequest, "the key in "+KeyHeader+" isn't the one whose SHA-256 was signed: nothing was kept")
		return
	}
	if len(key) > 4096 || strings.ContainsAny(key, "\r\n") {
		writeErr(w, http.StatusBadRequest, "that isn't a key planque can send")
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
	lang, text, signedOK := g.signed(req.Digest)
	sr.approval.Lang, sr.approval.Text = lang, text
	switch {
	case g.Request.Kind != KindAPIKey:
		writeErr(w, http.StatusConflict, "this grant signs in at the vendor; it takes no key")
		return
	case g.TokenOf != "" && g.Status == StatusConnected:
		writeErr(w, http.StatusConflict, "this connection uses another connection's key: replace it on that one, and every connection that uses it gets the new key")
		return
	case g.TokenOf != "":
		writeErr(w, http.StatusConflict, "this grant uses another connection's key: approve it with no key")
		return
	case !signedOK:
		writeErr(w, http.StatusConflict, "the signed digest is not this grant's: nothing was granted")
		return
	case g.Status != StatusPending && g.Status != StatusNeedsReconnect && g.Status != StatusConnected, g.SignedIn && g.Status == StatusPending:
		writeErr(w, http.StatusConflict, "this grant is "+g.Status+"; it takes no key now")
		return
	}
	a := sr.approval
	if err := profileApproval(g.Request.Profile, req.Profile); err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Status == StatusConnected {
		s.replaceKey(w, r, a, g, key)
		return
	}
	reconnect := g.Status == StatusNeedsReconnect
	if isTelegram(g) {
		bot, msg, err := s.setTelegramHook(r.Context(), g, key, nil)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if msg != "" {
			writeErr(w, http.StatusBadRequest, msg)
			return
		}
		if g.Telegram != nil && g.Telegram.Hook != bot.Hook {
			_, _, _ = s.KV.Delete(r.Context(), telegramHookPK, g.Telegram.Hook)
		}
		g.Telegram = bot
	}
	var prev *int64
	if reconnect {
		if _, it, err := s.openToken(r.Context(), g); err == nil {
			prev = &it.Ver
		}
	}
	if err := s.storeToken(r.Context(), g, oauth.Token{Access: key}, prev); err != nil {
		s.fail(w, r, err)
		return
	}
	g.Approval, g.SignedIn = &a, true
	if reconnect {
		g.Status, g.Reason = StatusConnected, ""
		g.event(s.now(), "a new key was given "+a.Who()+": connected again")
		if err := s.saveGrant(r.Context(), g); err != nil {
			s.fail(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, s.view(r, g))
		return
	}
	g.event(s.now(), "approved, with the key, "+a.Who())
	if len(g.Request.Tools) == 0 {
		if err := s.list(r.Context(), g, key); err != nil {
			writeErr(w, http.StatusBadGateway, g.Reason)
			return
		}
		if snap, err := s.loadSnapshot(r.Context(), g.Tenant, g.ID); err == nil {
			withOperations(g, snap)
			_ = s.saveGrant(r.Context(), g)
		}
		writeJSON(w, http.StatusOK, s.view(r, g))
		return
	}
	if err := s.connect(r.Context(), g, key); err != nil {
		s.log().WarnContext(r.Context(), "connect failed", slog.String("grant", g.ID), slog.String("err", err.Error()))
		writeErr(w, http.StatusBadGateway, g.Reason)
		return
	}
	writeJSON(w, http.StatusOK, s.view(r, g))
}

// replaceKey is a new key for a connected API-key grant (gap 13 of
// kiff/tafs#476): the owner replaces one the vendor refuses, from the
// connection's sheet in tafs, signed as the first was. planque reads once
// with the new key before keeping it, so a key the vendor refuses, or one it
// can't check, changes nothing and the old key stays. The grant, its tools
// and its Card stay as they are, and a connection that reuses this key gets
// the new one too.
func (s *Server) replaceKey(w http.ResponseWriter, r *http.Request, a Approval, g *Grant, key string) {
	ctx := r.Context()
	var checked string
	oldHook := ""
	undo := func() {}
	if isTelegram(g) {
		if g.Telegram == nil {
			writeErr(w, http.StatusConflict, "this Telegram connection names no bot: disconnect it and connect it again")
			return
		}
		bot, msg, err := s.setTelegramHook(ctx, g, key, g.Telegram)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if msg != "" {
			writeErr(w, http.StatusBadRequest, msg)
			return
		}
		oldHook, g.Telegram = g.Telegram.Hook, bot
		checked = "Telegram took it for @" + bot.Username
		// The bot's webhook is the new hook now: until the grant names it, a
		// failure points it back, so the taf's messages keep coming.
		undo = func() { s.repointTelegramHook(ctx, key, oldHook, bot.Hook) }
	} else {
		read, status, msg, err := s.checkKey(ctx, g, key)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if msg != "" {
			writeErr(w, status, msg)
			return
		}
		checked = "one read with it worked (" + read + ")"
	}
	var prev *int64
	if _, it, err := s.openToken(ctx, g); err == nil {
		prev = &it.Ver
	} else if !errors.Is(err, errNotFound) {
		undo()
		s.fail(w, r, err)
		return
	}
	if err := s.storeToken(ctx, g, oauth.Token{Access: key}, prev); err != nil {
		undo()
		s.fail(w, r, err)
		return
	}
	g.Reason = ""
	g.event(s.now(), "key replaced "+a.Who()+": "+checked+"; its tools are as they were")
	if err := s.saveGrant(ctx, g); err != nil {
		undo()
		s.fail(w, r, err)
		return
	}
	if oldHook != "" && oldHook != g.Telegram.Hook {
		_, _, _ = s.KV.Delete(ctx, telegramHookPK, oldHook)
	}
	s.log().InfoContext(ctx, "key replaced", slog.String("grant", g.ID))
	writeJSON(w, http.StatusOK, s.view(r, g))
}

// checkKey reads once with a new key for g, before it replaces the one
// held: read is the request it made ("GET /charges"). msg, with the status
// to answer, says why the key isn't kept: the vendor refused it, didn't
// answer, or the API has no read to check it with. Nothing changes then.
func (s *Server) checkKey(ctx context.Context, g *Grant, key string) (read string, status int, msg string, err error) {
	snap, err := s.loadSnapshot(ctx, g.Tenant, g.ID)
	if err != nil {
		return "", 0, "", err
	}
	op, method, target, body, ok := keyRead(g, snap)
	if !ok {
		return "", http.StatusConflict, "planque can't check a new key on this connection: its API has no read that takes no arguments. Nothing changed: to use another key, disconnect it and connect it again", nil
	}
	host := hostOf(g.Request.API.Base)
	read = method + " " + op.Path
	ans, err := s.REST.Call(ctx, method, target, body, g.Request.API.Auth, key, g.Request.API.VersionPin)
	switch {
	case err != nil:
		return read, http.StatusBadGateway, "planque couldn't check the new key: " + host + " didn't answer " + read + ". Nothing changed: try again", nil
	case ans.Status == http.StatusUnauthorized || ans.Status == http.StatusForbidden:
		return read, http.StatusBadRequest, fmt.Sprintf("%s refused the new key: it answered %d to %s. Nothing changed: the old key stays", host, ans.Status, read), nil
	case ans.Status < 200 || ans.Status > 299:
		return read, http.StatusBadGateway, fmt.Sprintf("planque couldn't check the new key: %s answered %d to %s. Nothing changed: try again", host, ans.Status, read), nil
	}
	return read, 0, "", nil
}

// keyRead is the read a new key is checked with, and its request: a
// read-only operation of g's API that needs no argument and sends no body,
// one of g's own tools first. Its path must be one of the API's own, under
// its base, on its host: the key goes nowhere else. ok is false when the API
// has none.
func keyRead(g *Grant, snap *openapi.Snapshot) (op openapi.Operation, method, target string, body []byte, ok bool) {
	own := map[string]bool{}
	for _, t := range g.Tools {
		own[t.Tool] = true
	}
	host := hostOf(g.Request.API.Base)
	for _, mine := range []bool{true, false} {
		for _, o := range snap.Operations {
			if own[o.Name] != mine || !o.ReadOnly || o.Body != "" {
				continue
			}
			m, t, b, err := o.Request(g.Request.API.Base, nil)
			if err != nil || !strings.HasPrefix(o.Path, "/") || strings.Contains(o.Path, "://") || strings.HasPrefix(o.Path, "//") || hostOf(t) != host {
				continue
			}
			return o, m, t, b, true
		}
	}
	return openapi.Operation{}, "", "", nil, false
}

// callAPI is one call of an API-key grant's tool: the request its operation
// makes, with the key where the grant says. A 2xx is the result; a 4xx the
// vendor's refusal, done; a 5xx or a lost answer unknown, never sent again;
// a call that never reached the vendor, not sent.
func (s *Server) callAPI(ctx context.Context, g *Grant, tool GrantedTool, key string, args map[string]any) (state string, res any) {
	if isTelegram(g) {
		return s.callTelegram(ctx, g, tool, key, args)
	}
	snap, err := s.loadSnapshot(ctx, g.Tenant, g.ID)
	if err != nil {
		return OpNotSent, toolResult("Not sent: planque couldn't read this connection's API description.", true)
	}
	i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == tool.Tool })
	if i < 0 {
		return OpNotSent, toolResult("Not sent: the API description has no operation "+tool.Tool+".", true)
	}
	if why := s.checkOneOf(ctx, g, snap, tool, key, args); why != "" {
		return OpNotSent, toolResult("Refused by planque: "+why+". Nothing was sent.", true)
	}
	method, target, body, err := snap.Operations[i].Request(g.Request.API.Base, args)
	if err != nil {
		return OpNotSent, toolResult("Refused by planque: "+err.Error()+". Nothing was sent.", true)
	}
	host := hostOf(g.Request.API.Base)
	if hostOf(target) != host {
		return OpNotSent, toolResult("Refused by planque: the call would leave "+host+". Nothing was sent.", true)
	}
	ans, err := s.REST.Call(ctx, method, target, body, g.Request.API.Auth, key, g.Request.API.VersionPin)
	var ce *rest.CallError
	switch {
	case errors.As(err, &ce) && !ce.Sent:
		return OpNotSent, toolResult("Not sent: "+host+" could not be reached.", true)
	case err != nil:
		s.log().WarnContext(ctx, "vendor answer lost", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return OpUnknown, unknownResult(g)
	case ans.Status >= 500:
		return OpUnknown, toolResult(fmt.Sprintf("Sent to %s, which answered %d: it may or may not have done it. planque will not send it again; check the account before making a new call.", host, ans.Status), true)
	case ans.Status >= 300 && ans.Status < 400:
		return OpUnknown, toolResult(fmt.Sprintf("Sent to %s, which answered %d, a redirect planque doesn't follow: it may or may not have done it. Check the account before making a new call.", host, ans.Status), true)
	}
	text := rest.HideKey(answerText(ans), key)
	if ans.Status >= 400 {
		// A status alone never says nothing happened (a batch can answer 400
		// after part of it ran): the refusal carries no evidence, and the
		// call stays counted in KIFF until something established says so.
		return OpDone, toolResult(fmt.Sprintf("%s refused it (%d): %s", host, ans.Status, text), true)
	}
	return OpDone, toolResult(text, false)
}

// answerText is a vendor's answer as the tool's text: its JSON as it is,
// else its text, said cut when it was.
func answerText(a rest.Answer) string {
	t := strings.TrimSpace(string(a.Body))
	if t == "" {
		t = fmt.Sprintf(`{"status": %d}`, a.Status)
	}
	if a.Cut {
		t += "\n(cut at 256 KB)"
	}
	return t
}

// checkOneOf says why a call's argument isn't one its one_of read allows,
// or "": for each such argument, planque calls the read (with no arguments,
// the grant's key, on the grant's host) and looks for the value at its path.
// A read that fails, or answers anything but a 2xx JSON object, refuses the
// call: planque never sends a value it couldn't check.
func (s *Server) checkOneOf(ctx context.Context, g *Grant, snap *openapi.Snapshot, tool GrantedTool, key string, args map[string]any) string {
	for _, name := range slices.Sorted(maps.Keys(tool.Arguments)) {
		a := tool.Arguments[name]
		if a.OneOf == nil {
			continue
		}
		parts, _ := pathParts(a.Path)
		v := valuesAt(args, parts)
		if len(v) != 1 {
			continue // not set: the vendor's own required check answers
		}
		i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == a.OneOf.Tool })
		if i < 0 || !snap.Operations[i].ReadOnly {
			return "argument " + name + " can't be checked: the API description has no read " + a.OneOf.Tool
		}
		method, target, body, err := snap.Operations[i].Request(g.Request.API.Base, nil)
		if err != nil || hostOf(target) != hostOf(g.Request.API.Base) {
			return "argument " + name + " can't be checked against " + a.OneOf.Tool
		}
		ans, err := s.REST.Call(ctx, method, target, body, g.Request.API.Auth, key, g.Request.API.VersionPin)
		var read any
		if err != nil || ans.Status < 200 || ans.Status >= 300 || ans.Cut || json.Unmarshal(ans.Body, &read) != nil {
			return "argument " + name + " can't be checked: " + a.OneOf.Tool + " couldn't be read"
		}
		want, _ := pathParts(a.OneOf.Path)
		if !slices.ContainsFunc(valuesAt(read, want), func(x string) bool { return sameValue(v[0], x) }) {
			return "argument " + name + " isn't one the account's " + a.OneOf.Tool + " lists"
		}
	}
	return ""
}
