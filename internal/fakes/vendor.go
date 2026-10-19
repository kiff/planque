// Package fakes are test doubles that speak the same standards and
// contracts as the real services: a vendor whose MCP server requires
// OAuth (as Stripe's does), and KIFF (apps/api, the gateway, permits).
package fakes

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/store"
)

// Vendor is an MCP server behind its own authorization server.
type Vendor struct {
	*httptest.Server
	TokenTTL     time.Duration // access token lifetime (default 1h)
	RefreshDelay time.Duration // widens refresh races in tests
	CallDelay    time.Duration // a slow tools/call
	AdvertiseIss bool          // authorization_response_iss_parameter_supported
	Account      string        // named in the token response
	// Base is the vendor's public URL when it isn't served by the embedded
	// test server (a deployed fake, cmd/fakevendor).
	Base string
	// KV keeps its clients, codes, tokens and calls: in memory by default,
	// a table for a deployed fake, so any of its instances answers.
	KV store.KV
	// Redirects, when set, are the only redirect URIs a client may
	// register: a deployed fake consents on its own, so any other would
	// make it an open redirector.
	Redirects []string

	mu        sync.Mutex
	Calls     []VendorCall // this instance's
	Refreshes int
	Revoked   []string
}

type code struct {
	Client, Redirect, Challenge, Resource string
}

// VendorCall is one tools/call the vendor ran.
type VendorCall struct {
	Tool string
	Args map[string]any
}

// The vendor's records in KV, by kind.
const (
	kvClient  = "fakevendor#client"  // SK client_id, Doc its redirect
	kvCode    = "fakevendor#code"    // SK the code, Doc a code
	kvAccess  = "fakevendor#access"  // SK the token, Doc its expiry (RFC 3339)
	kvRefresh = "fakevendor#refresh" // SK the token
	kvCall    = "fakevendor#call"    // SK when (ns) and an id, Doc a VendorCall
)

// keep writes one record, kept for ttl.
func (v *Vendor) keep(ctx context.Context, pk, sk string, doc []byte, ttl time.Duration) error {
	return v.KV.Put(ctx, store.Item{PK: pk, SK: sk, Doc: doc, TTL: time.Now().Add(ttl).Unix()}, store.Cond{})
}

// NewVendor starts a vendor on a local test server.
func NewVendor() *Vendor {
	v, h := NewVendorAt("")
	v.Server = httptest.NewServer(h)
	return v
}

// NewVendorAt is a vendor served at base by the handler it returns: for a
// deployed fake the live check and the nightly test sign in at.
func NewVendorAt(base string) (*Vendor, http.Handler) {
	v := &Vendor{Account: "acct_fake123", Base: strings.TrimRight(base, "/"), KV: &store.Memory{}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/oauth-protected-resource/mcp", v.prm)
	mux.HandleFunc("GET /.well-known/oauth-authorization-server/as", v.asMeta)
	mux.HandleFunc("POST /as/register", v.register)
	mux.HandleFunc("GET /as/authorize", v.authorize)
	mux.HandleFunc("POST /as/token", v.token)
	mux.HandleFunc("POST /as/revoke", v.revoke)
	mux.HandleFunc("POST /mcp", v.mcp)
	return v, mux
}

// root is where the vendor is reached.
func (v *Vendor) root() string {
	if v.Base != "" {
		return v.Base
	}
	return v.URL
}

// MCPURL is the vendor's MCP server.
func (v *Vendor) MCPURL() string { return v.root() + "/mcp" }

// Issuer is the vendor's authorization server.
func (v *Vendor) Issuer() string { return v.root() + "/as" }

// next is a new id, unique across the vendor's instances.
func next(prefix string) string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}

func (v *Vendor) prm(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"resource": v.MCPURL(), "authorization_servers": []string{v.Issuer()}, "scopes_supported": []string{"mcp"}})
}

func (v *Vendor) asMeta(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{
		"issuer": v.Issuer(), "authorization_endpoint": v.Issuer() + "/authorize", "token_endpoint": v.Issuer() + "/token",
		"registration_endpoint": v.Issuer() + "/register", "revocation_endpoint": v.Issuer() + "/revoke",
		"code_challenge_methods_supported": []string{"S256"}, "token_endpoint_auth_methods_supported": []string{"none"},
		"grant_types_supported": []string{"authorization_code", "refresh_token"}, "response_types_supported": []string{"code"},
		"authorization_response_iss_parameter_supported": v.AdvertiseIss,
	})
}

func (v *Vendor) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RedirectURIs []string `json:"redirect_uris"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if len(req.RedirectURIs) != 1 || len(v.Redirects) > 0 && !slices.Contains(v.Redirects, req.RedirectURIs[0]) {
		writeJSON(w, 400, map[string]string{"error": "invalid_redirect_uri"})
		return
	}
	id := next("client")
	if err := v.keep(r.Context(), kvClient, id, []byte(req.RedirectURIs[0]), 30*24*time.Hour); err != nil {
		writeJSON(w, 500, map[string]string{"error": "server_error"})
		return
	}
	writeJSON(w, 201, map[string]any{"client_id": id, "redirect_uris": req.RedirectURIs, "token_endpoint_auth_method": "none"})
}

// authorize stands in for the owner signing in and consenting: it
// redirects straight back with a code.
func (v *Vendor) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	it, ok, err := v.KV.Get(r.Context(), kvClient, q.Get("client_id"))
	redirect := string(it.Doc)
	if err != nil || !ok || q.Get("redirect_uri") != redirect || q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" {
		http.Error(w, "bad authorization request", 400)
		return
	}
	c := next("code")
	doc, _ := json.Marshal(code{q.Get("client_id"), redirect, q.Get("code_challenge"), q.Get("resource")})
	if err := v.keep(r.Context(), kvCode, c, doc, 10*time.Minute); err != nil {
		http.Error(w, "server error", 500)
		return
	}
	u, _ := url.Parse(redirect)
	rq := u.Query()
	rq.Set("code", c)
	rq.Set("state", q.Get("state"))
	rq.Set("iss", v.Issuer())
	u.RawQuery = rq.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (v *Vendor) issue(w http.ResponseWriter, r *http.Request) {
	ttl := v.TokenTTL
	if ttl == 0 {
		ttl = time.Hour
	}
	a, rt := next("access"), next("refresh")
	exp := time.Now().Add(ttl)
	if v.keep(r.Context(), kvAccess, a, []byte(exp.Format(time.RFC3339Nano)), ttl+time.Hour) != nil ||
		v.keep(r.Context(), kvRefresh, rt, []byte("1"), 30*24*time.Hour) != nil {
		writeJSON(w, 500, map[string]string{"error": "server_error"})
		return
	}
	writeJSON(w, 200, map[string]any{"access_token": a, "token_type": "bearer", "expires_in": int(ttl.Seconds()),
		"refresh_token": rt, "account": v.Account})
}

func (v *Vendor) token(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	switch r.PostForm.Get("grant_type") {
	case "authorization_code":
		// Deleting it is using it: a code works once.
		it, ok, err := v.KV.Delete(r.Context(), kvCode, r.PostForm.Get("code"))
		var c code
		if ok {
			_ = json.Unmarshal(it.Doc, &c)
		}
		sum := sha256.Sum256([]byte(r.PostForm.Get("code_verifier")))
		if err != nil || !ok || c.Client != r.PostForm.Get("client_id") || c.Redirect != r.PostForm.Get("redirect_uri") ||
			c.Challenge != base64.RawURLEncoding.EncodeToString(sum[:]) || c.Resource != v.MCPURL() {
			writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
			return
		}
		v.issue(w, r)
	case "refresh_token":
		time.Sleep(v.RefreshDelay)
		// Rotated: the old one stops working.
		if _, ok, err := v.KV.Delete(r.Context(), kvRefresh, r.PostForm.Get("refresh_token")); err != nil || !ok {
			writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
			return
		}
		v.mu.Lock()
		v.Refreshes++
		v.mu.Unlock()
		v.issue(w, r)
	default:
		writeJSON(w, 400, map[string]string{"error": "unsupported_grant_type"})
	}
}

func (v *Vendor) revoke(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	t := r.PostForm.Get("token")
	_, _, _ = v.KV.Delete(r.Context(), kvAccess, t)
	_, _, _ = v.KV.Delete(r.Context(), kvRefresh, t)
	v.mu.Lock()
	v.Revoked = append(v.Revoked, t)
	v.mu.Unlock()
	w.WriteHeader(200)
}

// Expire makes every access token expire now, as a vendor may.
func (v *Vendor) Expire() {
	ctx := context.Background()
	its, _ := v.KV.Query(ctx, kvAccess, "")
	for _, it := range its {
		it.Doc = []byte(time.Now().Add(-time.Second).Format(time.RFC3339Nano))
		_ = v.KV.Put(ctx, it, store.Cond{})
	}
}

// CallsSoFar is every tools/call the vendor ran, from all its instances,
// in order.
func (v *Vendor) CallsSoFar(ctx context.Context) ([]VendorCall, error) {
	its, err := v.KV.Query(ctx, kvCall, "")
	if err != nil {
		return nil, err
	}
	out := make([]VendorCall, 0, len(its))
	for _, it := range its {
		var c VendorCall
		if json.Unmarshal(it.Doc, &c) == nil {
			out = append(out, c)
		}
	}
	return out, nil
}

// CallCount is how many tools/call the vendor ran.
func (v *Vendor) CallCount() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.Calls)
}

// Live reports how many access and refresh tokens still work.
func (v *Vendor) Live() (access, refresh int) {
	ctx := context.Background()
	as, _ := v.KV.Query(ctx, kvAccess, "")
	for _, it := range as {
		if exp, err := time.Parse(time.RFC3339Nano, string(it.Doc)); err == nil && time.Now().Before(exp) {
			access++
		}
	}
	rs, _ := v.KV.Query(ctx, kvRefresh, "")
	return access, len(rs)
}

var vendorTools = []map[string]any{
	{"name": "create_refund", "description": "Refund a charge.", "inputSchema": map[string]any{"type": "object",
		"properties": map[string]any{"charge": map[string]any{"type": "string"}, "amount": map[string]any{"type": "integer"},
			"mode": map[string]any{"type": "string", "description": "Whether to act in live mode (real money) or test mode."}},
		"required": []string{"charge"}}},
	{"name": "list_charges", "description": "List charges.", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
	// What the sign-in covers, as Stripe's list_available_accounts_or_orgs
	// says it: read-only, and it takes nothing (kiff/tafs#89).
	{"name": "list_accounts", "description": "Lists the accounts in this session, with their mode.",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}, "annotations": map[string]any{"readOnlyHint": true}},
	// Shaped like Stripe's: one tool for every write.
	{"name": "stripe_api_write", "description": "Write via any API operation.", "inputSchema": map[string]any{"type": "object",
		"properties": map[string]any{"stripe_api_operation_id": map[string]any{"type": "string"},
			"parameters": map[string]any{"type": "object"}, "human_confirmation": map[string]any{"type": "object"}},
		"required": []string{"stripe_api_operation_id"}}},
}

func (v *Vendor) mcp(w http.ResponseWriter, r *http.Request) {
	tok, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	it, ok, err := v.KV.Get(r.Context(), kvAccess, tok)
	exp, perr := time.Parse(time.RFC3339Nano, string(it.Doc))
	if err != nil || !ok || perr != nil || time.Now().After(exp) {
		w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+v.root()+`/.well-known/oauth-protected-resource/mcp"`)
		writeJSON(w, 401, map[string]string{"error": "invalid_token"})
		return
	}
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
		Params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	result := func(res any) {
		writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": res})
	}
	switch req.Method {
	case "initialize":
		w.Header().Set("Mcp-Session-Id", "s1")
		result(map[string]any{"protocolVersion": "2025-06-18", "capabilities": map[string]any{"tools": map[string]any{}},
			"serverInfo": map[string]string{"name": "fake-vendor", "version": "1"}})
	case "notifications/initialized":
		w.WriteHeader(202)
	case "tools/list":
		result(map[string]any{"tools": vendorTools})
	case "tools/call":
		if req.Params.Name == "list_accounts" {
			// It changes nothing, so it isn't one of the vendor's calls.
			b, _ := json.Marshal(map[string]any{"accounts": []map[string]string{{"account": v.Account, "mode": "test"}}})
			result(map[string]any{"content": []map[string]string{{"type": "text", "text": string(b)}}})
			return
		}
		time.Sleep(v.CallDelay)
		call := VendorCall{Tool: req.Params.Name, Args: req.Params.Arguments}
		id := next("re")
		doc, _ := json.Marshal(call)
		if err := v.keep(r.Context(), kvCall, fmt.Sprintf("%020d#%s", time.Now().UnixNano(), id), doc, 7*24*time.Hour); err != nil {
			writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "error": map[string]any{"code": -32603, "message": "not recorded"}})
			return
		}
		v.mu.Lock()
		v.Calls = append(v.Calls, call)
		v.mu.Unlock()
		b, _ := json.Marshal(map[string]any{"id": id, "object": "refund", "amount": req.Params.Arguments["amount"]})
		result(map[string]any{"content": []map[string]string{{"type": "text", "text": string(b)}}})
	default:
		writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": req.ID, "error": map[string]any{"code": -32601, "message": "method not found"}})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
