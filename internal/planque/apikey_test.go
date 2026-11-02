package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// restVendor is a vendor with a REST API, an OpenAPI description and an API
// key, and no MCP server (kiff/planque#64).
type restVendor struct {
	*httptest.Server
	key string
	mu  sync.Mutex
	// desc is the description it serves now.
	desc  string
	calls []restCall
	// elsewhere counts requests to where a redirect points.
	elsewhere *httptest.Server
	away      int
}

type restCall struct {
	Method, Path, Query, Key string
	Body                     map[string]any
}

const restDesc = `{
 "openapi": "3.0.3",
 "info": {"title": "Pay API"},
 "security": [{"apiKey": []}],
 "components": {"securitySchemes": {"apiKey": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}},
 "paths": {
  "/charges": {"get": {"operationId": "listCharges", "summary": "List charges", "parameters": [{"name": "limit", "in": "query", "schema": {"type": "integer"}}]}},
  "/charges/{charge}/refunds": {"post": {"operationId": "createRefund", "summary": "Refund a charge",
   "parameters": [{"name": "charge", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["amount"], "properties": {"amount": {"type": "integer"}}}}}}}}
 }
}`

func newRESTVendor(t *testing.T) *restVendor {
	v := &restVendor{key: "sk_test_secret_123", desc: restDesc}
	v.elsewhere = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.mu.Lock()
		v.away++
		v.mu.Unlock()
	}))
	t.Cleanup(v.elsewhere.Close)
	v.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.mu.Lock()
		defer v.mu.Unlock()
		if r.URL.Path == "/openapi.json" {
			_, _ = io.WriteString(w, v.desc)
			return
		}
		if r.URL.Path == "/openapi.yaml" {
			_, _ = io.WriteString(w, "openapi: 3.0.3\n")
			return
		}
		c := restCall{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery, Key: r.Header.Get("X-API-Key")}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &c.Body)
		v.calls = append(v.calls, c)
		if c.Key != v.key {
			w.WriteHeader(401)
			_, _ = io.WriteString(w, `{"error":"bad key"}`)
			return
		}
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/charges/ch_boom/"):
			w.WriteHeader(503)
		case strings.HasPrefix(r.URL.Path, "/v1/charges/ch_away/"):
			http.Redirect(w, r, v.elsewhere.URL+"/steal", http.StatusFound)
		case c.Body["amount"] == float64(499):
			w.WriteHeader(402)
			_, _ = io.WriteString(w, `{"error":"insufficient funds"}`)
		case r.Method == "POST":
			_, _ = io.WriteString(w, `{"id":"re_1","status":"succeeded"}`)
		default:
			_, _ = io.WriteString(w, `{"data":[{"id":"ch_1"}]}`)
		}
	}))
	t.Cleanup(v.Server.Close)
	return v
}

func (v *restVendor) seen() []restCall {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]restCall(nil), v.calls...)
}

// askAPIGrant asks for an API-key grant to v, as Setup does.
func (e *env) askAPIGrant(v *restVendor, tools []map[string]any, extra map[string]any) (int, map[string]any) {
	e.t.Helper()
	body := map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": tools,
		"api": map[string]any{"url": v.URL + "/openapi.json", "base": v.URL + "/v1", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}}
	for k, x := range extra {
		body[k] = x
	}
	return e.do("POST", "/v1/accounts/"+tenant+"/grants", body, nil, "k_setup")
}

// giveKey is the owner's approval with the key, from the phone: the body
// signed names the key's SHA-256, and the key goes in its header.
func (e *env) giveKey(id, digest, key string) (int, map[string]any) {
	e.t.Helper()
	sum := sha256.Sum256([]byte(key))
	b, _ := json.Marshal(map[string]string{"digest": digest, "key_sha256": hex.EncodeToString(sum[:])})
	req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tenant+"/grants/"+id+"/key", bytes.NewReader(b))
	e.phone.Sign(req, tenant, b)
	req.Header.Set(KeyHeader, key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

var apiTools = []map[string]any{
	{"tool": "createRefund", "name": "create_refund", "title": "Refund a charge", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2},
	{"tool": "listCharges", "name": "list_charges", "read_only": true},
}

// connectAPI connects an API-key grant to v with its two tools: asked with
// none, the key given (listed), the tools set and approved.
func (e *env) connectAPI(v *restVendor) *Grant {
	e.t.Helper()
	e.enrol()
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 {
		e.t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, out := e.giveKey(id, out["digest"].(string), v.key); st != 200 || out["status"] != StatusListed {
		e.t.Fatalf("key: %d %v", st, out)
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": apiTools}, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	return g
}

func resultText(t *testing.T, raw json.RawMessage) (string, bool) {
	t.Helper()
	var r struct {
		Content []struct{ Text string } `json:"content"`
		IsError bool                    `json:"isError"`
	}
	if err := json.Unmarshal(raw, &r); err != nil || len(r.Content) == 0 {
		t.Fatalf("result: %s", raw)
	}
	return r.Content[0].Text, r.IsError
}

// The whole of an API-key grant: asked for with no tools, the key given on the
// phone, the catalog from the description, the tools approved with no second
// key, and calls that reach the vendor only with KIFF's permit.
func TestAnAPIKeyGrantEndToEnd(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 || out["status"] != StatusPending || !strings.Contains(out["text"].(string), "API key you give") ||
		!strings.Contains(out["text"].(string), "X-API-Key header") {
		t.Fatalf("ask: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	// Approve is the key's: a grant with no key yet isn't approved without it.
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, ""); st != 409 {
		t.Fatalf("approve without the key: %d", st)
	}
	// A key that isn't the one whose hash was signed is refused, and nothing kept.
	sum := sha256.Sum256([]byte("another key"))
	b, _ := json.Marshal(map[string]string{"digest": digest, "key_sha256": hex.EncodeToString(sum[:])})
	req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tenant+"/grants/"+id+"/key", bytes.NewReader(b))
	e.phone.Sign(req, tenant, b)
	req.Header.Set(KeyHeader, v.key)
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 400 {
		t.Fatalf("a key other than the signed one: %d", resp.StatusCode)
	}
	if st, out := e.giveKey(id, digest, v.key); st != 200 || out["status"] != StatusListed {
		t.Fatalf("key: %d %v", st, out)
	}
	// The key is sealed: no record holds it, the approval kept as evidence included.
	items, _ := e.kv.Query(context.Background(), tenantPK(tenant), "")
	for _, it := range items {
		if bytes.Contains(it.Doc, []byte(v.key)) {
			t.Fatalf("the key is in the clear in %s", it.SK)
		}
	}
	st, out = e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id+"/catalog", nil, nil, "k_setup")
	if st != 200 || !strings.Contains(mustJSONString(out), `"method":"POST"`) || !strings.Contains(mustJSONString(out), "createRefund") {
		t.Fatalf("catalog: %d %v", st, out)
	}
	// Only the description was read: nothing called at the vendor.
	if len(v.seen()) != 0 {
		t.Fatalf("calls before any tool was approved: %+v", v.seen())
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": apiTools}, nil, "k_setup")
	if st != 200 || !strings.Contains(out["text"].(string), "no second key") {
		t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	if c, ok := e.kiff.Connection(tenant, "create_refund"); !ok || c.AmountArg != "amount" || c.Mode != "relay" {
		t.Fatalf("not registered in KIFF: %+v", c)
	}

	ctx := context.Background()
	res, err := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 300})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("%+v %v", res, err)
	}
	text, isErr := resultText(t, res.Result)
	calls := v.seen()
	if isErr || !strings.Contains(text, "re_1") || len(calls) != 1 || calls[0].Method != "POST" || calls[0].Path != "/v1/charges/ch_1/refunds" ||
		calls[0].Key != v.key || calls[0].Body["amount"] != float64(300) || calls[0].Body["charge"] != nil || calls[0].Body[IdempotencyArg] != nil {
		t.Fatalf("%s %v %+v", text, isErr, calls)
	}
	// Outside the Card: refused by KIFF, nothing sent.
	if res, _ := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 800}); res.Decision != "refused" || len(v.seen()) != 1 {
		t.Fatalf("over the Card: %+v, %d calls", res, len(v.seen()))
	}
	// A read, with its query.
	res, _ = e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{"limit": 2})
	if text, isErr := resultText(t, res.Result); isErr || !strings.Contains(text, "ch_1") || v.seen()[1].Query != "limit=2" || v.seen()[1].Method != "GET" {
		t.Fatalf("%s %+v", text, v.seen())
	}
	// No permit, or another call's: nothing sent.
	g, _ := e.s.loadGrant(ctx, tenant, id)
	for _, p := range []string{"", e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Sub: "agt_taf", Op: "inv_x",
		Args: mustJSON(map[string]any{"charge": "ch_2", "amount": 1, IdempotencyArg: "inv_x"})})} {
		if out := e.relayCall(g, "create_refund", map[string]any{"charge": "ch_1", "amount": 1, IdempotencyArg: "inv_x"}, p); !strings.Contains(out, "Refused by planque") {
			t.Fatalf("permit %q: %s", p, out)
		}
	}
	if len(v.seen()) != 2 {
		t.Fatalf("sent without a valid permit: %+v", v.seen())
	}

	// Disconnecting deletes the key and the snapshot.
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	for _, it := range mustQuery(t, e) {
		if it == skToken(id) || strings.HasPrefix(it, "SNAPSHOT#"+id) {
			t.Fatalf("%s kept after disconnect", it)
		}
	}
	if _, ok := e.kiff.Connection(tenant, "create_refund"); ok {
		t.Fatal("KIFF still has the tool")
	}
}

func mustQuery(t *testing.T, e *env) []string {
	t.Helper()
	items, err := e.kv.Query(context.Background(), tenantPK(tenant), "")
	if err != nil {
		t.Fatal(err)
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.SK)
	}
	return out
}

func mustJSONString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// What the owner approved is the description as it was read: a changed one
// at the same URL changes nothing for the grant.
func TestAnAPIKeyGrantKeepsItsSnapshot(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	g := e.connectAPI(v)
	v.mu.Lock()
	v.desc = strings.ReplaceAll(restDesc, "/charges/{charge}/refunds", "/elsewhere/{charge}")
	v.mu.Unlock()
	if res, _ := e.kiff.Call(context.Background(), "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 10}); res.Decision != "allowed" {
		t.Fatalf("%+v", res)
	}
	if calls := v.seen(); len(calls) != 1 || calls[0].Path != "/v1/charges/ch_1/refunds" {
		t.Fatalf("%s: %+v", g.ID, calls)
	}
}

// The vendor's answers, as the gateway reads them: a refusal is done, with
// its words; a 5xx or a redirect is unknown, and a redirect isn't followed.
func TestAnAPIKeyCallsAnswers(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.connectAPI(v)
	ctx := context.Background()
	res, _ := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 499})
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "refused it (402)") || !strings.Contains(text, "insufficient funds") {
		t.Fatalf("%s", text)
	}
	res, _ = e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_boom", "amount": 5})
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "may or may not have done it") {
		t.Fatalf("%s", text)
	}
	res, _ = e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_away", "amount": 5})
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "redirect planque doesn't follow") {
		t.Fatalf("%s", text)
	}
	v.mu.Lock()
	away := v.away
	v.mu.Unlock()
	if away != 0 {
		t.Fatal("a redirect was followed, with the key")
	}
}

// A description outside what planque reads is refused when the grant is
// asked for, with the reason, before the owner is asked anything.
func TestAnAPIKeyGrantIsRefusedWithTheReason(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	for _, c := range []struct {
		api  map[string]any
		want string
	}{
		{map[string]any{"url": v.URL + "/openapi.yaml", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}, "JSON descriptions only"},
		{map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "query", "name": "key"}}, "none of its security schemes"},
		{map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}, "api.base is required"},
		{map[string]any{"url": v.URL + "/openapi.json", "base": "ftp://x.example", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}, "api.base"},
	} {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []any{}, "api": c.api}, nil, "k_setup")
		if st != 400 || !strings.Contains(mustJSONString(out), c.want) {
			t.Errorf("%v: %d %v", c.api, st, out)
		}
	}
	// A tool the description doesn't have.
	st, out := e.askAPIGrant(v, []map[string]any{{"tool": "deleteEverything"}}, nil)
	if st != 400 || !strings.Contains(mustJSONString(out), "no operation named deleteEverything") {
		t.Fatalf("%d %v", st, out)
	}
}

// A second grant to the same vendor and account reuses the first one's key
// or sign-in: one approval, no second key or sign-in. When the first is
// disconnected, the second must sign in (or give a key) again.
func TestASecondGrantReusesTheSignIn(t *testing.T) {
	ctx := context.Background()
	t.Run("api key", func(t *testing.T) {
		e := newEnv(t)
		v := newRESTVendor(t)
		first := e.connectAPI(v)
		e.kiff.Cards["k_taf"]["list_charges_b"] = 0
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"reuse": first.ID, "agent": "agt_taf",
			"tools": []map[string]any{{"tool": "listCharges", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
		if st != 201 || !strings.Contains(out["text"].(string), "no second key") {
			t.Fatalf("reuse: %d %v", st, out)
		}
		id := out["id"].(string)
		if st, out := e.giveKey(id, out["digest"].(string), v.key); st != 409 {
			t.Fatalf("a key for a reused one: %d %v", st, out)
		}
		_, out = e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
		if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
			t.Fatalf("approve: %d %v", st, out)
		}
		if res, _ := e.kiff.Call(ctx, "k_taf", "list_charges_b", map[string]any{}); res.Decision != "allowed" {
			t.Fatalf("%+v", res)
		} else if text, isErr := resultText(t, res.Result); isErr || v.seen()[0].Key != v.key {
			t.Fatalf("%s %+v", text, v.seen())
		}
		// One key, the first's: none of its own.
		if mustHave(mustQuery(t, e), skToken(id)) {
			t.Fatal("the reused key was copied")
		}
		// The first disconnected: the second can't call, and waits to be connected again.
		if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, e.phone, ""); st != 200 {
			t.Fatal("disconnect")
		}
		res, _ := e.kiff.Call(ctx, "k_taf", "list_charges_b", map[string]any{})
		if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "Not sent") {
			t.Fatalf("%s", text)
		}
		if g, _ := e.s.loadGrant(ctx, tenant, id); g.Status != StatusNeedsReconnect || g.TokenOf != "" {
			t.Fatalf("%s %q", g.Status, g.TokenOf)
		}
		// A reuse of a grant that isn't connected is refused when asked.
		if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"reuse": first.ID, "agent": "agt_taf", "tools": []any{}}, nil, "k_setup"); st != 409 {
			t.Fatalf("reuse of a disconnected grant: %d %v", st, out)
		}
	})
	t.Run("sign-in", func(t *testing.T) {
		e := newEnv(t)
		first := e.connect()
		e.kiff.Cards["k_taf"]["list_charges_b"] = 0
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
			"tools": []map[string]any{{"tool": "list_charges", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
		if st != 201 || !strings.Contains(out["text"].(string), "already signed in") {
			t.Fatalf("reuse: %d %v", st, out)
		}
		id := out["id"].(string)
		// One approval, no sign-in link: connected at once.
		st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
		if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] != nil {
			t.Fatalf("approve: %d %v", st, out)
		}
		if res, _ := e.kiff.Call(ctx, "k_taf", "list_charges_b", map[string]any{}); res.Decision != "allowed" {
			t.Fatalf("%+v", res)
		} else if text, isErr := resultText(t, res.Result); isErr {
			t.Fatalf("%s", text)
		}
		// Another vendor's server is refused.
		if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": "https://other.example/mcp", "reuse": first.ID, "agent": "agt_taf", "tools": []any{}}, nil, "k_setup"); st != 409 {
			t.Fatalf("reuse for another server: %d", st)
		}
		// The first disconnected (its sign-in revoked): the second signs in again.
		if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, e.phone, ""); st != 200 {
			t.Fatal("disconnect")
		}
		res, _ := e.kiff.Call(ctx, "k_taf", "list_charges_b", map[string]any{})
		if _, isErr := resultText(t, res.Result); !isErr {
			t.Fatal("a call through a revoked sign-in")
		}
		if g, _ := e.s.loadGrant(ctx, tenant, id); g.Status != StatusNeedsReconnect {
			t.Fatalf("%s", g.Status)
		}
	})
}

func mustHave(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
