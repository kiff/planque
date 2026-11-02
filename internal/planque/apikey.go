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
	URL  string       `json:"url"`
	Base string       `json:"base"`
	Auth openapi.Auth `json:"auth"`
}

// readAPI reads the description an API-key grant names, checks the tools
// asked for against it and keeps it as the grant's snapshot. msg is the
// refusal in words, for the agent; err is planque's own failure.
func (s *Server) readAPI(ctx context.Context, tenant, grant string, a apiAsk, tools []ToolRequest) (*APIRequest, *openapi.Snapshot, string, error) {
	a.URL, a.Base = strings.TrimSpace(a.URL), strings.TrimSpace(a.Base)
	a.Auth.In, a.Auth.Name = strings.TrimSpace(a.Auth.In), strings.TrimSpace(a.Auth.Name)
	if a.URL == "" {
		return nil, nil, "api.url is required: the vendor's OpenAPI description, in JSON", nil
	}
	raw, err := s.REST.Fetch(ctx, a.URL)
	if err != nil {
		return nil, nil, "Couldn't read the API description at " + a.URL + ": " + err.Error(), nil
	}
	snap, err := openapi.Parse(raw, a.Auth)
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
	return &APIRequest{URL: a.URL, Digest: snap.Digest, Base: strings.TrimRight(a.Base, "/"), Auth: a.Auth}, snap, "", nil
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

// vendorTools are the tools g's vendor offers, as an MCP server lists them:
// listed by its MCP server, or read from an API-key grant's snapshot.
func (s *Server) vendorTools(ctx context.Context, g *Grant, access string) ([]mcpTool, mcp.Server, error) {
	if g.Request.Kind != KindAPIKey {
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
// taking the key for (needs_reconnect) takes a new one the same way.
func (s *Server) giveKey(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		Digest    string `json:"digest"`
		KeySHA256 string `json:"key_sha256"`
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
	switch {
	case g.Request.Kind != KindAPIKey:
		writeErr(w, http.StatusConflict, "this grant signs in at the vendor; it takes no key")
		return
	case g.TokenOf != "":
		writeErr(w, http.StatusConflict, "this grant uses another connection's key: approve it with no key")
		return
	case req.Digest != g.Digest:
		writeErr(w, http.StatusConflict, "the signed digest is not this grant's: nothing was granted")
		return
	case g.Status != StatusPending && g.Status != StatusNeedsReconnect, g.SignedIn && g.Status == StatusPending:
		writeErr(w, http.StatusConflict, "this grant is "+g.Status+"; it takes no key now")
		return
	}
	a := sr.approval
	reconnect := g.Status == StatusNeedsReconnect
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
		writeJSON(w, http.StatusOK, s.view(g))
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
		writeJSON(w, http.StatusOK, s.view(g))
		return
	}
	if err := s.connect(r.Context(), g, key); err != nil {
		s.log().WarnContext(r.Context(), "connect failed", slog.String("grant", g.ID), slog.String("err", err.Error()))
		writeErr(w, http.StatusBadGateway, g.Reason)
		return
	}
	writeJSON(w, http.StatusOK, s.view(g))
}

// callAPI is one call of an API-key grant's tool: the request its operation
// makes, with the key where the grant says. A 2xx is the result; a 4xx the
// vendor's refusal, done; a 5xx or a lost answer unknown, never sent again;
// a call that never reached the vendor, not sent.
func (s *Server) callAPI(ctx context.Context, g *Grant, tool GrantedTool, key string, args map[string]any) (state string, res any) {
	snap, err := s.loadSnapshot(ctx, g.Tenant, g.ID)
	if err != nil {
		return OpNotSent, toolResult("Not sent: planque couldn't read this connection's API description.", true)
	}
	i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == tool.Tool })
	if i < 0 {
		return OpNotSent, toolResult("Not sent: the API description has no operation "+tool.Tool+".", true)
	}
	method, target, body, err := snap.Operations[i].Request(g.Request.API.Base, args)
	if err != nil {
		return OpNotSent, toolResult("Refused by planque: "+err.Error()+". Nothing was sent.", true)
	}
	host := hostOf(g.Request.API.Base)
	if hostOf(target) != host {
		return OpNotSent, toolResult("Refused by planque: the call would leave "+host+". Nothing was sent.", true)
	}
	ans, err := s.REST.Call(ctx, method, target, body, g.Request.API.Auth, key)
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
