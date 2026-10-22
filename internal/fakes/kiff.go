package fakes

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/kiff"
)

// KIFF stands in for apps/api and the gateway, following their contracts:
// /v1/me, /v1/connections (kiff-cloud apps/gateway/internal/connection),
// and execution permits (RFC 046), minted only for a call the agent's
// Card allows.
type KIFF struct {
	Issuer string // the base URL, also the permit issuer
	Key    ed25519.PrivateKey
	KID    string
	Keys   map[string]kiff.Me // KIFF key → who it speaks for
	// Cards: agent key → action → the most one call may move. An action
	// missing from a Card is refused.
	Cards map[string]map[string]int64
	// Devices: the devices people approve with, by KIFF's id
	// (kiff-cloud #1267), checked at /v1/devices/verify.
	Devices map[string]Device

	// SetupReads are the read Cards asked for (kiff/tafs#138), in order;
	// SetupReadsRefused makes KIFF refuse them.
	SetupReads        []kiff.SetupReads
	SetupReadsRefused bool

	mu      sync.Mutex
	halts   map[string]halt       // tenant/agent ("" the account) → its halt
	down    bool                  // the halt routes don't answer
	conns   map[string]Connection // tenant/name → connection
	appLast map[string]string
	HTTP    *http.Client
}

// Device is a person's device, as KIFF keeps it.
type Device struct {
	Tenant, Subject, Platform, Name string
	Phone                           *Phone
	// What the person's role allows: "tools.connect" to connect tools;
	// "account.export" and "account.delete", an admin's, to export or
	// delete the account.
	Capabilities []string
	Member       bool // an active member of the account
	Revoked      bool
}

// Connection is a tool registered in the fake gateway.
type Connection struct {
	kiff.Connection
	Tenant string
}

// NewKIFF makes a fake KIFF with a fresh signing key.
func NewKIFF() *KIFF {
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	return &KIFF{Key: priv, KID: "test-1", Keys: map[string]kiff.Me{}, Cards: map[string]map[string]int64{}, Devices: map[string]Device{}}
}

// halt is an agent's halt, or its account's.
type halt struct {
	on bool
	at time.Time
}

// Halt halts agent of tenant ("" the account) at at, as the owner does.
func (k *KIFF) Halt(tenant, agent string, at time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.halts == nil {
		k.halts = map[string]halt{}
	}
	k.halts[tenant+"/"+agent] = halt{on: true, at: at}
}

// Resume resumes it, keeping the halt's time.
func (k *KIFF) Resume(tenant, agent string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	h := k.halts[tenant+"/"+agent]
	h.on = false
	k.halts[tenant+"/"+agent] = h
}

// HaltsDown makes the halt routes fail, as KIFF unreachable.
func (k *KIFF) HaltsDown(down bool) {
	k.mu.Lock()
	k.down = down
	k.mu.Unlock()
}

// setupReads answers as apps/api's POST /v1/me/setup-reads: a connector key
// asks for Setup's read Card with the owner's approval, kept here as asked.
func (k *KIFF) setupReads(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok || !slices.Contains(me.Roles, "connector") {
		writeJSON(w, 403, map[string]string{"error": "connector_key_required"})
		return
	}
	var req kiff.SetupReads
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if dec.Decode(&req) != nil {
		writeJSON(w, 400, map[string]string{"error": "bad body"})
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.SetupReadsRefused {
		writeJSON(w, 403, map[string]string{"error": "the owner's approval doesn't check"})
		return
	}
	k.SetupReads = append(k.SetupReads, req)
	writeJSON(w, 201, map[string]string{"id": "setup-reads-" + req.Grant})
}

// haltHTTP answers as apps/api's GET /v1/agents/{id}/halt and GET /v1/halt,
// to a key of the account that isn't an agent's: the agent's halt over its
// account's.
func (k *KIFF) haltHTTP(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok || me.AgentID != "" {
		writeJSON(w, 403, map[string]string{"error": "forbidden"})
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.down {
		writeJSON(w, 503, map[string]string{"error": "the halt could not be read"})
		return
	}
	// A trial (agent~trial) is its agent, as in apps/api.
	agent, _ := strings.CutSuffix(r.PathValue("id"), "~trial")
	own, all := k.halts[me.TenantID+"/"+agent], k.halts[me.TenantID+"/"]
	if r.PathValue("id") == "" {
		own = halt{}
	}
	at := own.at
	if all.at.After(at) {
		at = all.at
	}
	out := map[string]any{"state": "running", "halted": false}
	if own.on || all.on {
		out["state"], out["halted"] = "halted", true
	}
	if !at.IsZero() {
		out["halted_at"] = at
	}
	writeJSON(w, 200, out)
}

// Handler serves the fake's routes.
func (k *KIFF) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/me", k.me)
	mux.HandleFunc("POST /v1/devices/verify", k.verifyDevice)
	mux.HandleFunc("GET /.well-known/kiff-permit-keys.json", k.jwks)
	mux.HandleFunc("POST /v1/me/setup-reads", k.setupReads)
	mux.HandleFunc("GET /v1/agents/{id}/halt", k.haltHTTP)
	mux.HandleFunc("GET /v1/halt", k.haltHTTP)
	mux.HandleFunc("POST /v1/connections", k.connect)
	mux.HandleFunc("GET /v1/connections", k.list)
	mux.HandleFunc("DELETE /v1/connections/{name}", k.remove)
	// Test-only: an agent's tools/call through the gateway.
	mux.HandleFunc("POST /fake/call", k.callHTTP)
	mux.HandleFunc("POST /fake/forward", k.forwardHTTP)
	// Test-only: the owner's app, where a finished sign-in sends the
	// browser with its confirmation code.
	mux.HandleFunc("GET /fake/app", k.appOpened)
	mux.HandleFunc("GET /fake/app/last", k.appLastHTTP)
	return mux
}

func (k *KIFF) who(r *http.Request) (kiff.Me, bool) {
	key, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	me, ok := k.Keys[key]
	return me, ok
}

func (k *KIFF) me(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok {
		writeJSON(w, 401, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, 200, me)
}

// verifyDevice answers as apps/api's /v1/devices/verify (kiff-cloud
// #1267): a connector key asks about a device of its own account.
func (k *KIFF) verifyDevice(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok {
		writeJSON(w, 401, map[string]string{"error": "unauthorized"})
		return
	}
	if !me.APIKey || me.AgentID != "" || !slices.Contains(me.Roles, "connector") {
		writeJSON(w, 403, map[string]string{"error": "connector_key_required"})
		return
	}
	var req struct {
		Device     string `json:"device_id"`
		Message    string `json:"message"`
		Signature  string `json:"signature"`
		Capability string `json:"capability"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	msg, err := base64.StdEncoding.DecodeString(req.Message)
	if err != nil || len(msg) == 0 || strings.HasPrefix(string(msg), "kiff.") || !slices.Contains([]string{"tools.connect", "holds.answer", "account.export", "account.delete"}, req.Capability) {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	no := func(why string) { writeJSON(w, 200, map[string]any{"valid": false, "reason": why}) }
	d, ok := k.Devices[req.Device]
	switch {
	case !ok || d.Tenant != me.TenantID || d.Revoked:
		no("unknown_device")
	case !d.Phone.Verifies(msg, req.Signature):
		no("signature_invalid")
	case !d.Member:
		no("not_a_member")
	case !slices.Contains(d.Capabilities, req.Capability):
		no("not_allowed")
	default:
		writeJSON(w, 200, map[string]any{"valid": true, "device_id": req.Device, "subject": d.Subject, "platform": d.Platform, "name": d.Name, "attested": false})
	}
}

func (k *KIFF) jwks(w http.ResponseWriter, _ *http.Request) {
	pub := k.Key.Public().(ed25519.PublicKey)
	writeJSON(w, 200, map[string]any{"keys": []map[string]string{{"kty": "OKP", "crv": "Ed25519", "alg": "Ed25519",
		"kid": k.KID, "x": base64.RawURLEncoding.EncodeToString(pub)}}})
}

func (k *KIFF) canManage(me kiff.Me) bool {
	return me.AgentID == "" && !me.DeviceSession && slices.Contains(me.Capabilities, "tools.connect")
}

// connect checks a relay connection as the gateway does: the relay lists
// the tool with the transport credential, the idempotency argument is a
// string argument and the amount an integer one.
func (k *KIFF) connect(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok || !k.canManage(me) {
		writeJSON(w, 403, map[string]string{"error": "connecting tools requires tools.connect"})
		return
	}
	var c kiff.Connection
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		writeJSON(w, 400, map[string]string{"error": "bad body"})
		return
	}
	if c.Mode != "relay" || c.Audience == "" || c.IdempotencyArg == "" || c.Credential == "" {
		writeJSON(w, 400, map[string]string{"error": "a relay connection needs an audience, an idempotency argument and a credential"})
		return
	}
	tools, err := k.listTools(r.Context(), c.URL, c.Credential)
	if err != nil {
		writeJSON(w, 502, map[string]string{"error": "KIFF could not list the tools at that URL: " + err.Error()})
		return
	}
	var props map[string][]string
	for _, t := range tools {
		if t.Name == c.Tool {
			props = schemaTypes(t.InputSchema)
		}
	}
	switch {
	case props == nil:
		writeJSON(w, 502, map[string]string{"error": "the server has no tool named " + c.Tool})
		return
	case !slices.Contains(props[c.IdempotencyArg], "string"):
		writeJSON(w, 400, map[string]string{"error": "idempotency_argument must be a string argument of the tool"})
		return
	case c.AmountArg != "" && !slices.Contains(props[c.AmountArg], "integer"):
		writeJSON(w, 400, map[string]string{"error": "amount_argument must be an integer argument of the tool"})
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.conns == nil {
		k.conns = map[string]Connection{}
	}
	if _, dup := k.conns[me.TenantID+"/"+c.Name]; dup {
		writeJSON(w, 409, map[string]string{"error": "a tool named " + c.Name + " is already connected"})
		return
	}
	k.conns[me.TenantID+"/"+c.Name] = Connection{Connection: c, Tenant: me.TenantID}
	writeJSON(w, 201, map[string]string{"name": c.Name, "mode": c.Mode})
}

func (k *KIFF) list(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok {
		writeJSON(w, 401, map[string]string{"error": "unauthorized"})
		return
	}
	// As the gateway lists them: each connection with its name.
	var conns []map[string]string
	for _, c := range k.Connections(me.TenantID) {
		conns = append(conns, map[string]string{"name": c.Name})
	}
	writeJSON(w, 200, map[string]any{"connections": conns})
}

func (k *KIFF) remove(w http.ResponseWriter, r *http.Request) {
	me, ok := k.who(r)
	if !ok || !k.canManage(me) {
		writeJSON(w, 403, map[string]string{"error": "forbidden"})
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	key := me.TenantID + "/" + r.PathValue("name")
	if _, ok := k.conns[key]; !ok {
		writeJSON(w, 404, map[string]string{"error": "no such tool"})
		return
	}
	delete(k.conns, key)
	w.WriteHeader(204)
}

// Add registers a connection directly, as if made earlier.
func (k *KIFF) Add(c Connection) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.conns == nil {
		k.conns = map[string]Connection{}
	}
	k.conns[c.Tenant+"/"+c.Name] = c
}

func (k *KIFF) appOpened(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	k.mu.Lock()
	k.appLast = map[string]string{"tenant": q.Get("tenant"), "grant": q.Get("grant"), "confirm": q.Get("confirm")}
	k.mu.Unlock()
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Signed in. Confirm in the app."))
}

func (k *KIFF) appLastHTTP(w http.ResponseWriter, r *http.Request) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.appLast == nil || k.appLast["grant"] != r.URL.Query().Get("grant") {
		writeJSON(w, 404, map[string]string{"error": "not yet"})
		return
	}
	writeJSON(w, 200, k.appLast)
}

// Connections lists a tenant's connected tools.
func (k *KIFF) Connections(tenant string) []Connection {
	k.mu.Lock()
	defer k.mu.Unlock()
	var out []Connection
	for _, c := range k.conns {
		if c.Tenant == tenant {
			out = append(out, c)
		}
	}
	slices.SortFunc(out, func(a, b Connection) int { return strings.Compare(a.Name, b.Name) })
	return out
}

// Connection returns one connected tool.
func (k *KIFF) Connection(tenant, name string) (Connection, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	c, ok := k.conns[tenant+"/"+name]
	return c, ok
}

// Permit is what Mint signs.
type Permit struct {
	Audience, Tenant, Tool, Sub, Op string
	Args                            json.RawMessage
	IssuedAt                        time.Time
	Lifetime                        time.Duration
}

// AskHaltConfirm signs KIFF's ask that planque confirm agent's halt of
// tenant at at, as apps/api's halts.Planques does; life 0 is 60 seconds.
func (k *KIFF) AskHaltConfirm(tenant, agent string, at time.Time, life time.Duration) string {
	if life == 0 {
		life = 60 * time.Second
	}
	now := time.Now()
	claims := map[string]any{"iss": k.Issuer, "iat": now.Unix(), "exp": now.Add(life).Unix(), "jti": "pmt_" + randHex(8),
		"tenant": tenant, "agent": agent, "halted_at": at.UTC().Format(time.RFC3339Nano)}
	h, _ := json.Marshal(map[string]string{"alg": "Ed25519", "typ": "kiff-halt-confirm+jwt", "kid": k.KID})
	c, _ := json.Marshal(claims)
	in := b64(h) + "." + b64(c)
	return in + "." + b64(ed25519.Sign(k.Key, []byte(in)))
}

// Mint signs a permit as apps/api does (kiff-cloud apps/api/internal/permit).
func (k *KIFF) Mint(p Permit) string {
	if p.IssuedAt.IsZero() {
		p.IssuedAt = time.Now()
	}
	if p.Lifetime == 0 {
		p.Lifetime = 60 * time.Second
	}
	sum, err := jcs.SHA256(p.Args)
	if err != nil {
		panic(err)
	}
	claims := map[string]any{"iss": k.Issuer, "aud": p.Audience, "sub": p.Sub, "iat": p.IssuedAt.Unix(),
		"exp": p.IssuedAt.Add(p.Lifetime).Unix(), "jti": "pmt_" + randHex(8),
		"authorization_details": []map[string]string{{"type": "https://kiff.dev/permits/tool-call", "tenant": p.Tenant,
			"tool": p.Tool, "args_sha256": sum, "op": p.Op, "card": "crd_test", "decision": "dec_" + randHex(6)}}}
	h, _ := json.Marshal(map[string]string{"alg": "Ed25519", "typ": "kiff-permit+jwt", "kid": k.KID})
	c, _ := json.Marshal(claims)
	in := b64(h) + "." + b64(c)
	return in + "." + b64(ed25519.Sign(k.Key, []byte(in)))
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CallResult is what the gateway answers an agent.
type CallResult struct {
	Decision string          `json:"decision"` // allowed, refused
	Op       string          `json:"op,omitempty"`
	Result   json.RawMessage `json:"result,omitempty"`
}

// Call is an agent's tools/call through the gateway: the Card decides;
// only an allowed call gets a permit and is forwarded to the relay.
func (k *KIFF) Call(ctx context.Context, agentKey, name string, args map[string]any) (CallResult, error) {
	me, ok := k.Keys[agentKey]
	if !ok || me.AgentID == "" {
		return CallResult{}, errors.New("not an agent key")
	}
	c, ok := k.Connection(me.TenantID, name)
	if !ok {
		return CallResult{}, fmt.Errorf("no tool %s", name)
	}
	limit, ok := k.Cards[agentKey][name]
	amount := int64(0)
	if c.AmountArg != "" {
		switch v := args[c.AmountArg].(type) {
		case float64:
			amount = int64(v)
		case int:
			amount = int64(v)
		case int64:
			amount = v
		}
	}
	if !ok || amount > limit {
		return CallResult{Decision: "refused"}, nil
	}
	op := "inv_" + randHex(10)
	fwd := map[string]any{}
	for k, v := range args {
		fwd[k] = v
	}
	fwd[c.IdempotencyArg] = op
	body, _ := json.Marshal(fwd)
	permit := k.Mint(Permit{Audience: c.Audience, Tenant: me.TenantID, Tool: c.Tool, Sub: me.AgentID, Op: op, Args: body})
	res, err := k.Forward(ctx, c, body, permit)
	return CallResult{Decision: "allowed", Op: op, Result: res}, err
}

func (k *KIFF) callHTTP(w http.ResponseWriter, r *http.Request) {
	key, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	var req struct {
		Tool      string         `json:"tool"`
		Arguments map[string]any `json:"arguments"`
		Parallel  int            `json:"parallel"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSON(w, 400, map[string]string{"error": "bad body"})
		return
	}
	n := max(req.Parallel, 1)
	out := make([]CallResult, n)
	errs := make([]string, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := k.Call(r.Context(), key, req.Tool, req.Arguments)
			out[i] = res
			if err != nil {
				errs[i] = err.Error()
			}
		}()
	}
	wg.Wait()
	writeJSON(w, 200, map[string]any{"results": out, "errors": errs})
}

// forwardHTTP sends a given permit with the stored transport credential:
// a gateway gone wrong, for testing what the relay refuses.
func (k *KIFF) forwardHTTP(w http.ResponseWriter, r *http.Request) {
	key, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	me, ok := k.Keys[key]
	var req struct {
		Name      string          `json:"name"`
		Tool      string          `json:"tool"`
		Arguments json.RawMessage `json:"arguments"`
		Permit    string          `json:"permit"`
	}
	if !ok || json.NewDecoder(r.Body).Decode(&req) != nil {
		writeJSON(w, 400, map[string]string{"error": "bad request"})
		return
	}
	c, ok := k.Connection(me.TenantID, req.Name)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "no such tool"})
		return
	}
	if req.Tool != "" {
		c.Tool = req.Tool
	}
	res, err := k.Forward(r.Context(), c, req.Arguments, req.Permit)
	if err != nil {
		writeJSON(w, 502, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"result": res, "audience": c.Audience})
}

func (k *KIFF) hc() *http.Client {
	if k.HTTP != nil {
		return k.HTTP
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// Forward sends one tools/call to a relay the way the gateway's client
// does: server/discover, then a legacy session, the permit on tools/call
// only.
func (k *KIFF) Forward(ctx context.Context, c Connection, args json.RawMessage, permit string) (json.RawMessage, error) {
	var sid string
	post := func(id int, method string, params any, hdr map[string]string) (json.RawMessage, error) {
		m := map[string]any{"jsonrpc": "2.0", "method": method, "params": params}
		if id >= 0 {
			m["id"] = id
		}
		b, _ := json.Marshal(m)
		req, _ := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		req.Header.Set("Authorization", "Bearer "+c.Credential)
		if sid != "" {
			req.Header.Set("Mcp-Session-Id", sid)
		}
		for h, v := range hdr {
			req.Header.Set(h, v)
		}
		resp, err := k.hc().Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if s := resp.Header.Get("Mcp-Session-Id"); s != "" {
			sid = s
		}
		raw, _ := io.ReadAll(resp.Body)
		if id < 0 {
			return nil, nil
		}
		var msg struct {
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, raw)
		}
		if msg.Error != nil {
			return nil, fmt.Errorf("rpc %d: %s", msg.Error.Code, msg.Error.Message)
		}
		return msg.Result, nil
	}
	if _, err := post(0, "server/discover", map[string]any{}, nil); err == nil {
		return nil, errors.New("the relay claims the modern revision; this fake speaks the legacy one")
	}
	if _, err := post(1, "initialize", map[string]any{"protocolVersion": "2025-11-25", "capabilities": map[string]any{},
		"clientInfo": map[string]string{"name": "kiff-gateway", "version": "1"}}, nil); err != nil {
		return nil, err
	}
	_, _ = post(-1, "notifications/initialized", nil, nil)
	return post(2, "tools/call", map[string]any{"name": c.Tool, "arguments": args}, map[string]string{"KIFF-Permit": permit})
}

type listedTool struct {
	Name        string          `json:"name"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

func (k *KIFF) listTools(ctx context.Context, url, cred string) ([]listedTool, error) {
	c := Connection{Connection: kiff.Connection{URL: url, Credential: cred}}
	var out []listedTool
	b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": map[string]any{}})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cred)
	resp, err := k.hc().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var msg struct {
		Result struct {
			Tools []listedTool `json:"tools"`
		} `json:"result"`
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	out = msg.Result.Tools
	return out, nil
}

func schemaTypes(schema json.RawMessage) map[string][]string {
	var s struct {
		Properties map[string]struct {
			Type any `json:"type"`
		} `json:"properties"`
	}
	_ = json.Unmarshal(schema, &s)
	out := map[string][]string{}
	for name, p := range s.Properties {
		switch t := p.Type.(type) {
		case string:
			out[name] = []string{t}
		case []any:
			for _, v := range t {
				if v, ok := v.(string); ok {
					out[name] = append(out[name], v)
				}
			}
		}
	}
	return out
}
