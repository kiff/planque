// Package fakes are test doubles that speak the same standards and
// contracts as the real services: a vendor whose MCP server requires
// OAuth (as Stripe's does), and KIFF (apps/api, the gateway, permits).
package fakes

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Vendor is an MCP server behind its own authorization server.
type Vendor struct {
	*httptest.Server
	TokenTTL     time.Duration // access token lifetime (default 1h)
	RefreshDelay time.Duration // widens refresh races in tests
	CallDelay    time.Duration // a slow tools/call
	AdvertiseIss bool          // authorization_response_iss_parameter_supported
	Account      string        // named in the token response

	mu        sync.Mutex
	clients   map[string]string // client_id → redirect
	codes     map[string]code
	access    map[string]time.Time // token → expiry
	refresh   map[string]bool      // valid refresh tokens
	n         int
	Calls     []VendorCall
	Refreshes int
	Revoked   []string
}

type code struct{ client, redirect, challenge, resource string }

// VendorCall is one tools/call the vendor ran.
type VendorCall struct {
	Tool string
	Args map[string]any
}

// NewVendor starts a vendor.
func NewVendor() *Vendor {
	v := &Vendor{Account: "acct_fake123", clients: map[string]string{}, codes: map[string]code{},
		access: map[string]time.Time{}, refresh: map[string]bool{}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/oauth-protected-resource/mcp", v.prm)
	mux.HandleFunc("GET /.well-known/oauth-authorization-server/as", v.asMeta)
	mux.HandleFunc("POST /as/register", v.register)
	mux.HandleFunc("GET /as/authorize", v.authorize)
	mux.HandleFunc("POST /as/token", v.token)
	mux.HandleFunc("POST /as/revoke", v.revoke)
	mux.HandleFunc("POST /mcp", v.mcp)
	v.Server = httptest.NewServer(mux)
	return v
}

// MCPURL is the vendor's MCP server.
func (v *Vendor) MCPURL() string { return v.URL + "/mcp" }

// Issuer is the vendor's authorization server.
func (v *Vendor) Issuer() string { return v.URL + "/as" }

func (v *Vendor) next(prefix string) string {
	v.n++
	return fmt.Sprintf("%s_%d", prefix, v.n)
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
	if len(req.RedirectURIs) != 1 {
		writeJSON(w, 400, map[string]string{"error": "invalid_redirect_uri"})
		return
	}
	v.mu.Lock()
	id := v.next("client")
	v.clients[id] = req.RedirectURIs[0]
	v.mu.Unlock()
	writeJSON(w, 201, map[string]any{"client_id": id, "redirect_uris": req.RedirectURIs, "token_endpoint_auth_method": "none"})
}

// authorize stands in for the owner signing in and consenting: it
// redirects straight back with a code.
func (v *Vendor) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	v.mu.Lock()
	defer v.mu.Unlock()
	redirect, ok := v.clients[q.Get("client_id")]
	if !ok || q.Get("redirect_uri") != redirect || q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" {
		http.Error(w, "bad authorization request", 400)
		return
	}
	c := v.next("code")
	v.codes[c] = code{q.Get("client_id"), redirect, q.Get("code_challenge"), q.Get("resource")}
	u, _ := url.Parse(redirect)
	rq := u.Query()
	rq.Set("code", c)
	rq.Set("state", q.Get("state"))
	rq.Set("iss", v.Issuer())
	u.RawQuery = rq.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (v *Vendor) issue(w http.ResponseWriter) {
	ttl := v.TokenTTL
	if ttl == 0 {
		ttl = time.Hour
	}
	a, rt := v.next("access"), v.next("refresh")
	v.access[a] = time.Now().Add(ttl)
	v.refresh[rt] = true
	writeJSON(w, 200, map[string]any{"access_token": a, "token_type": "bearer", "expires_in": int(ttl.Seconds()),
		"refresh_token": rt, "stripe_user_id": v.Account})
}

func (v *Vendor) token(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	switch r.PostForm.Get("grant_type") {
	case "authorization_code":
		v.mu.Lock()
		defer v.mu.Unlock()
		c, ok := v.codes[r.PostForm.Get("code")]
		delete(v.codes, r.PostForm.Get("code"))
		sum := sha256.Sum256([]byte(r.PostForm.Get("code_verifier")))
		if !ok || c.client != r.PostForm.Get("client_id") || c.redirect != r.PostForm.Get("redirect_uri") ||
			c.challenge != base64.RawURLEncoding.EncodeToString(sum[:]) || c.resource != v.MCPURL() {
			writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
			return
		}
		v.issue(w)
	case "refresh_token":
		time.Sleep(v.RefreshDelay)
		v.mu.Lock()
		defer v.mu.Unlock()
		rt := r.PostForm.Get("refresh_token")
		if !v.refresh[rt] {
			writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
			return
		}
		delete(v.refresh, rt) // rotated: the old one stops working
		v.Refreshes++
		v.issue(w)
	default:
		writeJSON(w, 400, map[string]string{"error": "unsupported_grant_type"})
	}
}

func (v *Vendor) revoke(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	t := r.PostForm.Get("token")
	v.mu.Lock()
	delete(v.access, t)
	delete(v.refresh, t)
	v.Revoked = append(v.Revoked, t)
	v.mu.Unlock()
	w.WriteHeader(200)
}

// Expire makes every access token expire now, as a vendor may.
func (v *Vendor) Expire() {
	v.mu.Lock()
	defer v.mu.Unlock()
	for t := range v.access {
		v.access[t] = time.Now().Add(-time.Second)
	}
}

// CallCount is how many tools/call the vendor ran.
func (v *Vendor) CallCount() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.Calls)
}

// Live reports how many access and refresh tokens still work.
func (v *Vendor) Live() (access, refresh int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, exp := range v.access {
		if time.Now().Before(exp) {
			access++
		}
	}
	return access, len(v.refresh)
}

var vendorTools = []map[string]any{
	{"name": "create_refund", "description": "Refund a charge.", "inputSchema": map[string]any{"type": "object",
		"properties": map[string]any{"charge": map[string]any{"type": "string"}, "amount": map[string]any{"type": "integer"}},
		"required":   []string{"charge"}}},
	{"name": "list_charges", "description": "List charges.", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
	// Shaped like Stripe's: one tool for every write.
	{"name": "stripe_api_write", "description": "Write via any API operation.", "inputSchema": map[string]any{"type": "object",
		"properties": map[string]any{"stripe_api_operation_id": map[string]any{"type": "string"},
			"parameters": map[string]any{"type": "object"}, "human_confirmation": map[string]any{"type": "object"}},
		"required": []string{"stripe_api_operation_id"}}},
}

func (v *Vendor) mcp(w http.ResponseWriter, r *http.Request) {
	tok, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	v.mu.Lock()
	exp, ok := v.access[tok]
	v.mu.Unlock()
	if !ok || time.Now().After(exp) {
		w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+v.URL+`/.well-known/oauth-protected-resource/mcp"`)
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
		time.Sleep(v.CallDelay)
		v.mu.Lock()
		v.Calls = append(v.Calls, VendorCall{Tool: req.Params.Name, Args: req.Params.Arguments})
		n := len(v.Calls)
		v.mu.Unlock()
		b, _ := json.Marshal(map[string]any{"id": fmt.Sprintf("re_%d", n), "object": "refund", "amount": req.Params.Arguments["amount"]})
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
