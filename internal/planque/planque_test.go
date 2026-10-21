package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/egress"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

const tenant = "tnt_test"

type env struct {
	t      *testing.T
	app    *app
	vendor *fakes.Vendor
	kiff   *fakes.KIFF
	phone  *fakes.Phone
	s      *Server
	url    string
	kv     *store.Memory
}

func newEnv(t *testing.T) *env {
	t.Helper()
	e := &env{t: t, vendor: fakes.NewVendor(), kiff: fakes.NewKIFF(), phone: fakes.NewPhone(), kv: &store.Memory{}, app: newApp(t)}
	t.Cleanup(e.vendor.Close)
	ks := httptest.NewServer(e.kiff.Handler())
	t.Cleanup(ks.Close)
	e.kiff.Issuer = ks.URL
	e.kiff.Keys = map[string]kiff.Me{
		"k_connector": {APIKey: true, TenantID: tenant, Roles: []string{"connector"}, Capabilities: []string{"read", "tools.connect"}},
		"k_admin":     {APIKey: true, TenantID: tenant, Roles: []string{"admin"}, Capabilities: []string{"read", "tools.connect", "keys.manage", "cards.manage"}},
		// A capability KIFF might add later, which planque has not vetted.
		"k_future": {APIKey: true, TenantID: tenant, Roles: []string{"connector"}, Capabilities: []string{"read", "tools.connect", "webhooks.manage"}},
		"k_setup":  {APIKey: true, TenantID: tenant, AgentID: "agt_setup"},
		"k_taf":    {APIKey: true, TenantID: tenant, AgentID: "agt_taf"},
		"k_other":  {APIKey: true, TenantID: "tnt_other", AgentID: "agt_x"},
	}
	// The taf's Card: refunds up to 500 a call; listing moves nothing.
	e.kiff.Cards["k_taf"] = map[string]int64{"create_refund": 500, "list_charges": 0}
	var h http.Handler
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }))
	t.Cleanup(ps.Close)
	e.url = ps.URL
	e.s = e.server()
	h = e.s.Handler()
	// The owner's phone is a device KIFF knows, signing with its id (d=), as
	// the tafs app does: planque takes no other signature (kiff/planque#12).
	e.phone = e.device("dev_phone")
	return e
}

// server is one planque process on the shared table, as on Lambda.
func (e *env) server() *Server {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	return &Server{BaseURL: e.url, AppReturn: e.app.URL + "/signed-in", KV: e.kv, Sealer: &seal.Local{Key: key},
		KIFF:    &kiff.Client{API: e.kiff.Issuer, Gateway: e.kiff.Issuer},
		Permits: &permit.Verifier{Issuer: e.kiff.Issuer, Keys: &permit.JWKS{URL: e.kiff.Issuer + "/.well-known/kiff-permit-keys.json"}},
		OAuth:   &oauth.Client{AllowHTTP: true}, MCP: &mcp.Client{}}
}

// do sends a request; phone signs it when not nil, key is a KIFF bearer.
func (e *env) do(method, path string, body any, phone *fakes.Phone, key string) (int, map[string]any) {
	e.t.Helper()
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, e.url+path, bytes.NewReader(b))
	if phone != nil {
		phone.Sign(req, tenant, b)
	}
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	raw, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(raw, &out)
	return resp.StatusCode, out
}

func (e *env) enrol() {
	e.t.Helper()
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/enrolment",
		map[string]string{"device_key": e.phone.PublicKey(), "connector_key": "k_connector"}, e.phone, ""); st != 201 {
		e.t.Fatalf("enrol: %d %v", st, out)
	}
}

func (e *env) request() (id, digest string) {
	e.t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "create_refund", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2},
			{"tool": "list_charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		e.t.Fatalf("request: %d %v", st, out)
	}
	return out["id"].(string), out["digest"].(string)
}

// app stands in for the owner's app (tafs): where the browser that
// finished a vendor sign-in is sent, with the code that confirms it.
type app struct {
	*httptest.Server
	mu   sync.Mutex
	last map[string]string
}

func newApp(t *testing.T) *app {
	a := &app{}
	a.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.mu.Lock()
		a.last = map[string]string{"grant": r.URL.Query().Get("grant"), "confirm": r.URL.Query().Get("confirm"), "tenant": r.URL.Query().Get("tenant")}
		a.mu.Unlock()
		w.Write([]byte("app opened"))
	}))
	t.Cleanup(a.Close)
	return a
}

func (a *app) code() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.last["confirm"]
}

// signIn opens a sign-in link in a browser whose app is e.app, and
// returns the confirmation code the app received.
func (e *env) signIn(link string) string {
	e.t.Helper()
	resp, err := http.Get(link)
	if err != nil {
		e.t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || string(page) != "app opened" {
		e.t.Fatalf("sign-in did not reach the app: %d %s", resp.StatusCode, page)
	}
	return e.app.code()
}

// approveAndSignIn approves grant id on e.phone, signs in at the vendor
// and confirms on the phone. It returns the confirm answer.
func (e *env) approveAndSignIn(id, digest string) (int, map[string]any) {
	e.t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 200 {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	code := e.signIn(out["sign_in_url"].(string))
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, "")
}

// connect runs the whole connection and returns the grant.
func (e *env) connect() *Grant {
	e.t.Helper()
	e.enrol()
	id, digest := e.request()
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("confirm: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	return g
}

// relayCall forwards a call to planque's relay with permit, as the gateway.
func (e *env) relayCall(g *Grant, tool string, args map[string]any, permit string) string {
	e.t.Helper()
	c, ok := e.kiff.Connection(tenant, "create_refund")
	if !ok {
		e.t.Fatal("create_refund is not connected")
	}
	c.Tool = tool
	b, _ := json.Marshal(args)
	res, err := e.kiff.Forward(context.Background(), c, b, permit)
	if err != nil {
		e.t.Fatalf("forward: %v", err)
	}
	return string(res)
}

// 1. Enrolled once, from tafs.
func TestEnrol(t *testing.T) {
	e := newEnv(t)
	path := "/v1/accounts/" + tenant + "/enrolment"
	body := func(key string) map[string]string {
		return map[string]string{"device_key": e.phone.PublicKey(), "connector_key": key}
	}
	if st, _ := e.do("POST", path, body("k_connector"), nil, ""); st != 401 {
		t.Fatalf("unsigned enrolment: %d, want 401", st)
	}
	if st, out := e.do("POST", path, body("k_admin"), e.phone, ""); st != 400 || !strings.Contains(out["error"].(string), "more than connect tools") {
		t.Fatalf("admin key: %d %v, want 400 more than connect tools", st, out)
	}
	if st, out := e.do("POST", path, body("k_future"), e.phone, ""); st != 400 || !strings.Contains(out["error"].(string), "webhooks.manage") {
		t.Fatalf("unvetted capability: %d %v, want 400", st, out)
	}
	if st, _ := e.do("POST", path, body("k_setup"), e.phone, ""); st != 400 {
		t.Fatalf("agent key as connector: %d, want 400", st)
	}
	if st, _ := e.do("POST", path, body("k_other"), e.phone, ""); st != 400 {
		t.Fatalf("other account's key: %d, want 400", st)
	}
	// Signed by a key other than the one enrolled: no proof of possession.
	other := fakes.NewPhone()
	b, _ := json.Marshal(body("k_connector"))
	req, _ := http.NewRequest("POST", e.url+path, bytes.NewReader(b))
	other.Sign(req, tenant, b)
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 401 {
		t.Fatalf("enrolment signed by another key: %d, want 401", resp.StatusCode)
	}
	e.enrol()
	// The connector key is kept sealed, never in the clear.
	var en Enrolment
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &en); err != nil || len(en.Connector) == 0 ||
		bytes.Contains(en.Connector, []byte("k_connector")) {
		t.Fatal("the connector key is stored in the clear")
	}
	// A phone signing with its key alone can't enrol, even sending it (kiff/planque#12).
	if st, _ := e.do("POST", path, map[string]string{"device_key": other.PublicKey(), "connector_key": "k_connector"}, other, ""); st != 401 {
		t.Fatalf("a key-only phone: %d, want 401", st)
	}
}

// 2 and 3. The agent asks; only the phone's signature on that exact
// request approves it.
func TestRequestAndApprove(t *testing.T) {
	e := newEnv(t)
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "a",
		"tools": []map[string]any{{"tool": "list_charges"}}}, nil, "k_setup"); st != 409 {
		t.Fatalf("grant before enrolment: %d, want 409", st)
	}
	e.enrol()
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "a",
		"tools": []map[string]any{{"tool": "list_charges"}}}, nil, "k_other"); st != 403 {
		t.Fatalf("another account's agent: %d, want 403", st)
	}
	id, digest := e.request()
	_, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
	text := out["text"].(string)
	for _, want := range []string{"taf agt_taf", "create_refund: moves money (amount), only within the taf's Card", "list_charges: reads only", "sign in at 127.0.0.1"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q:\n%s", want, text)
		}
	}
	var q Request
	b, _ := json.Marshal(out["request"])
	_ = json.Unmarshal(b, &q)
	if d, _ := digest2(q); d != digest {
		t.Fatalf("digest is not the JCS SHA-256 of the request")
	}
	approve := "/v1/accounts/" + tenant + "/grants/" + id + "/approve"
	if st, _ := e.do("POST", approve, map[string]string{"digest": digest}, nil, ""); st != 401 {
		t.Fatalf("unsigned approval: %d, want 401", st)
	}
	if st, _ := e.do("POST", approve, map[string]string{"digest": digest}, fakes.NewPhone(), ""); st != 401 {
		t.Fatalf("approval by another phone: %d, want 401", st)
	}
	if st, _ := e.do("POST", approve, map[string]string{"digest": strings.Repeat("0", 64)}, e.phone, ""); st != 409 {
		t.Fatalf("approval of another digest: %d, want 409", st)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusPending || g.Approval != nil {
		t.Fatalf("after refused approvals the grant is %s", g.Status)
	}
	if len(e.kiff.Connections(tenant)) != 0 || e.vendor.CallCount() != 0 {
		t.Fatal("something was connected without an approval")
	}
	// A signed request replayed byte for byte is refused.
	b, _ = json.Marshal(map[string]string{"digest": digest})
	req, _ := http.NewRequest("POST", e.url+approve, bytes.NewReader(b))
	e.phone.Sign(req, tenant, b)
	sig := req.Header.Get(device.Header)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Fatalf("approval: %d", resp.StatusCode)
	}
	req, _ = http.NewRequest("POST", e.url+approve, bytes.NewReader(b))
	req.Header.Set(device.Header, sig)
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 409 {
		t.Fatalf("replayed approval: %d, want 409", resp.StatusCode)
	}
}

func digest2(q Request) (string, error) { return digest(q) }

// 4 and 5. Standard sign-in, the token sealed, the tools registered.
func TestSignInAndRegister(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	if g.Status != StatusConnected || g.Account != "acct_fake123" {
		t.Fatalf("grant %s, account %q", g.Status, g.Account)
	}
	// The token is sealed and bound to its grant.
	var rec tokenRecord
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skToken(g.ID), &rec); err != nil ||
		bytes.Contains(rec.Sealed, []byte("access_")) || bytes.Contains(rec.Sealed, []byte("refresh_")) {
		t.Fatalf("the token is missing or stored in the clear")
	}
	other := *g
	other.ID = "grt_other"
	if _, err := e.s.Sealer.Open(context.Background(), rec.Sealed, tokenBinding(&other)); err == nil {
		t.Fatal("a sealed token opens for another grant")
	}
	conns := e.kiff.Connections(tenant)
	if len(conns) != 2 {
		t.Fatalf("connections: %+v", conns)
	}
	for _, c := range conns {
		if c.Mode != "relay" || c.Audience != g.Audience || c.IdempotencyArg != IdempotencyArg || c.URL != g.Audience {
			t.Errorf("connection %+v", c.Connection)
		}
		if strings.Contains(c.Credential, "access_") {
			t.Error("KIFF received the vendor token")
		}
	}
	if c, _ := e.kiff.Connection(tenant, "create_refund"); c.AmountArg != "amount" || c.AmountUnit != "EUR" || *c.AmountScale != 2 {
		t.Errorf("create_refund amount: %+v", c.Connection)
	}
}

func TestCallbackChecks(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	signIn := out["sign_in_url"].(string)
	noFollow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, _ := noFollow.Get(signIn)
	back := resp.Header.Get("Location")
	// An unknown state grants nothing.
	r2, _ := http.Get(strings.Replace(back, "state=", "state=x", 1))
	if r2.StatusCode != 400 {
		t.Fatalf("unknown state: %d", r2.StatusCode)
	}
	// An answer naming another issuer grants nothing (RFC 9207).
	r3, _ := http.Get(strings.Replace(back, "iss=", "iss=https%3A%2F%2Fevil.example&x=", 1))
	if r3.StatusCode != 400 {
		t.Fatalf("other issuer: %d", r3.StatusCode)
	}
	// The state was taken by that attempt: the link works once.
	r4, _ := http.Get(back)
	if r4.StatusCode != 400 {
		t.Fatalf("reused state: %d", r4.StatusCode)
	}
	if len(e.kiff.Connections(tenant)) != 0 {
		t.Fatal("connected without a valid sign-in")
	}
}

// 6. Permit-gated calls.
func TestPermitGatedCalls(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()

	// Inside the Card it goes through, without planque's argument.
	res, err := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 300})
	if err != nil || res.Decision != "allowed" || !strings.Contains(string(res.Result), `\"object\":\"refund\"`) {
		t.Fatalf("refund inside the Card: %+v %v", res, err)
	}
	if e.vendor.CallCount() != 1 {
		t.Fatalf("vendor calls: %d", e.vendor.CallCount())
	}
	if _, ok := e.vendor.Calls[0].Args[IdempotencyArg]; ok {
		t.Fatal("the vendor received planque's operation argument")
	}
	// Outside it, KIFF mints no permit and nothing is sent.
	res, _ = e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 800})
	if res.Decision != "refused" || e.vendor.CallCount() != 1 {
		t.Fatalf("refund outside the Card: %+v, vendor calls %d", res, e.vendor.CallCount())
	}

	args := func(op string) map[string]any {
		return map[string]any{"charge": "ch_2", "amount": 100, IdempotencyArg: op}
	}
	mint := func(p fakes.Permit) string {
		if p.Audience == "" {
			p.Audience = g.Audience
		}
		if p.Tenant == "" {
			p.Tenant = tenant
		}
		if p.Tool == "" {
			p.Tool = "create_refund"
		}
		if p.Args == nil {
			p.Args, _ = json.Marshal(args(p.Op))
		}
		return e.kiff.Mint(p)
	}
	otherKey := fakes.NewKIFF()
	otherKey.Issuer, otherKey.KID = e.kiff.Issuer, e.kiff.KID
	type refusal struct{ name, op, callOp, permit, want string }
	cases := []refusal{
		{"missing", "op_m", "op_m", "", "no execution permit"},
		{"other arguments", "op_a", "op_a", mint(fakes.Permit{Op: "op_a", Args: json.RawMessage(`{"charge":"ch_2","amount":999,"planque_op":"op_a"}`)}), "other arguments"},
		{"other tool", "op_b", "op_b", mint(fakes.Permit{Op: "op_b", Tool: "list_charges"}), "another tool"},
		{"other audience", "op_c", "op_c", mint(fakes.Permit{Op: "op_c", Audience: e.url + "/mcp/" + tenant + "/grt_other"}), "another audience"},
		{"other account", "op_d", "op_d", mint(fakes.Permit{Op: "op_d", Tenant: "tnt_other"}), "another KIFF account"},
		{"expired", "op_e", "op_e", mint(fakes.Permit{Op: "op_e", IssuedAt: time.Now().Add(-2 * time.Minute)}), "expired"},
		{"too long-lived", "op_f", "op_f", mint(fakes.Permit{Op: "op_f", Lifetime: time.Hour}), "lifetime"},
		{"not KIFF's key", "op_g", "op_g", otherKey.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Op: "op_g", Args: mustJSON(args("op_g"))}), "signature"},
		// Signed over the call's arguments, but naming another op.
		{"op not the call's", "op_h", "op_other", mint(fakes.Permit{Op: "op_h", Args: mustJSON(args("op_other"))}), "operation"},
	}
	for _, c := range cases {
		out := e.relayCall(g, "create_refund", args(c.callOp), c.permit)
		if !strings.Contains(out, "Refused by planque") || !strings.Contains(out, c.want) {
			t.Errorf("%s: %s", c.name, out)
		}
	}
	// A tampered permit: one claim changed after signing.
	good := mint(fakes.Permit{Op: "op_t"})
	parts := strings.Split(good, ".")
	claims, _ := base64.RawURLEncoding.DecodeString(parts[1])
	parts[1] = base64.RawURLEncoding.EncodeToString(bytes.Replace(claims, []byte(tenant), []byte("tnt_evil"), 1))
	if out := e.relayCall(g, "create_refund", args("op_t"), strings.Join(parts, ".")); !strings.Contains(out, "signature") {
		t.Errorf("tampered: %s", out)
	}
	if n := e.vendor.CallCount(); n != 1 {
		t.Fatalf("refused permits reached the vendor: %d calls", n)
	}

	// A valid permit runs once; replayed, it returns the first answer.
	p := mint(fakes.Permit{Op: "op_once"})
	first := e.relayCall(g, "create_refund", args("op_once"), p)
	again := e.relayCall(g, "create_refund", args("op_once"), p)
	if first != again || e.vendor.CallCount() != 2 {
		t.Fatalf("replay: %s / %s, vendor calls %d", first, again, e.vendor.CallCount())
	}
	// Ten concurrent presentations of one new permit: one execution, even
	// while the vendor is slow to answer.
	e.vendor.CallDelay = 200 * time.Millisecond
	p = mint(fakes.Permit{Op: "op_race"})
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() { defer wg.Done(); e.relayCall(g, "create_refund", args("op_race"), p) }()
	}
	wg.Wait()
	e.vendor.CallDelay = 0
	if e.vendor.CallCount() != 3 {
		t.Fatalf("concurrent presentations: vendor calls %d, want 3", e.vendor.CallCount())
	}
	// An expired permit for a recorded op still learns how it ended.
	late := mint(fakes.Permit{Op: "op_once", IssuedAt: time.Now().Add(-5 * time.Minute)})
	if out := e.relayCall(g, "create_refund", args("op_once"), late); out != first {
		t.Fatalf("expired permit for a recorded op: %s", out)
	}
	// Without the transport credential, or with another, the relay
	// answers nothing.
	for _, cred := range []string{"", "Bearer not-the-credential"} {
		req, _ := http.NewRequest("POST", g.Audience, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`))
		if cred != "" {
			req.Header.Set("Authorization", cred)
		}
		if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 401 {
			t.Fatalf("relay with credential %q: %d", cred, resp.StatusCode)
		}
	}
}

// barrierKV holds the first two reads that find no record under prefix
// until both happened, so two requests race to write the same record.
type barrierKV struct {
	store.KV
	wg     *sync.WaitGroup
	prefix string
}

func (b barrierKV) Get(ctx context.Context, pk, sk string) (store.Item, bool, error) {
	it, ok, err := b.KV.Get(ctx, pk, sk)
	if strings.HasPrefix(sk, b.prefix) && !ok {
		b.wg.Done()
		b.wg.Wait()
	}
	return it, ok, err
}

// Two presentations that both found no record: only one claims the op.
func TestClaimIsExclusive(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	var wg sync.WaitGroup
	wg.Add(2)
	e.s.KV = barrierKV{KV: e.kv, wg: &wg, prefix: "OP#"}
	args := map[string]any{"charge": "ch_9", "amount": 100, IdempotencyArg: "op_claim"}
	p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Op: "op_claim", Args: mustJSON(args)})
	var done sync.WaitGroup
	for range 2 {
		done.Add(1)
		go func() { defer done.Done(); e.relayCall(g, "create_refund", args, p) }()
	}
	done.Wait()
	if n := e.vendor.CallCount(); n != 1 {
		t.Fatalf("vendor calls: %d, want 1", n)
	}
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

// 6. The token refreshes when it expires; concurrent calls in two
// processes refresh once, and the rotated refresh token is kept.
func TestRefreshOnce(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	e.vendor.RefreshDelay = 300 * time.Millisecond
	e.vendor.Expire()
	// Mark the stored token as expired too, as time passing would.
	tok, it, _ := e.s.openToken(ctx, g)
	tok.Expiry = time.Now().Add(-time.Second)
	_ = e.s.storeToken(ctx, g, tok, &it.Ver)

	second := e.server() // another Lambda instance on the same table
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			if i%2 == 0 {
				var res fakes.CallResult
				res, err = e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": fmt.Sprintf("ch_%d", i), "amount": 10})
				if err == nil && !strings.Contains(string(res.Result), `re_`) {
					err = fmt.Errorf("call %d: %s", i, res.Result)
				}
			} else {
				_, _, err = second.accessToken(ctx, g, 0)
			}
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if e.vendor.Refreshes != 1 {
		t.Fatalf("refreshes: %d, want 1", e.vendor.Refreshes)
	}
	if e.vendor.CallCount() != 10 {
		t.Fatalf("vendor calls: %d, want 10", e.vendor.CallCount())
	}
	// The rotated refresh token was kept: the next refresh works.
	e.vendor.Expire()
	res, err := e.kiff.Call(ctx, "k_taf", "create_refund", map[string]any{"charge": "ch_x", "amount": 10})
	if err != nil || !strings.Contains(string(res.Result), "re_") || e.vendor.Refreshes != 2 {
		t.Fatalf("after rotation: %s %v, refreshes %d", res.Result, err, e.vendor.Refreshes)
	}
}

// 7 and 8. Disconnecting, and the oversight record.
func TestDisconnectAndExport(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	if _, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{}); err != nil {
		t.Fatal(err)
	}
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup"); st != 401 {
		t.Fatalf("an agent disconnecting: %d, want 401", st)
	}
	st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, "")
	if st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if n := len(e.kiff.Connections(tenant)); n != 0 {
		t.Fatalf("KIFF still lists %d tools", n)
	}
	if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), skToken(g.ID)); ok {
		t.Fatal("the token is still stored")
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the vendor still honours %d access and %d refresh tokens", a, r)
	}
	if out := e.relayCallRaw(g); out != 401 {
		t.Fatalf("relay after disconnect: %d", out)
	}

	// The export carries what is needed to check each approval offline.
	st, out = e.do("GET", "/v1/accounts/"+tenant+"/export", nil, e.phone, "")
	if st != 200 {
		t.Fatalf("export: %d %v", st, out)
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/export", nil, nil, "k_setup"); st != 401 {
		t.Fatalf("export without the phone: %d", st)
	}
	var ex struct {
		Grants []grantView `json:"grants"`
	}
	b, _ := json.Marshal(out)
	_ = json.Unmarshal(b, &ex)
	if len(ex.Grants) != 1 || ex.Grants[0].Approval == nil {
		t.Fatalf("export: %s", b)
	}
	gv := ex.Grants[0]
	a := gv.Approval
	// Signed by the device it names, whose key KIFF keeps.
	if a.DeviceID != "dev_phone" || !e.phone.Verifies([]byte(a.Message), a.Signature) {
		t.Fatal("the exported signature does not verify with the named device's key")
	}
	bodySum := sha256.Sum256([]byte(a.Body))
	if lines := strings.Split(a.Message, "\n"); lines[1] != e.url || lines[5] != fmt.Sprintf("%x", bodySum) {
		t.Fatal("the signed message does not cover the exported body")
	}
	rq, _ := json.Marshal(gv.Request)
	d, _ := jcs.SHA256(rq)
	if !strings.Contains(a.Body, d) || d != gv.Digest {
		t.Fatal("the signed body does not name the exported request's digest")
	}
	var whats []string
	for _, ev := range gv.Events {
		whats = append(whats, ev.What)
	}
	for _, want := range []string{"requested by agt_setup", "approved on tafs on phone", "signed in at", "registered in KIFF", "disconnected on tafs on phone"} {
		if !strings.Contains(strings.Join(whats, "\n"), want) {
			t.Errorf("record lacks %q: %v", want, whats)
		}
	}
	// Listing works with the phone's signature too.
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants", nil, e.phone, ""); st != 200 || len(out["grants"].([]any)) != 1 {
		t.Fatalf("list: %d %v", st, out)
	}
}

func (e *env) relayCallRaw(g *Grant) int {
	req, _ := http.NewRequest("POST", g.Audience, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`))
	req.Header.Set("Authorization", "Bearer anything")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	return resp.StatusCode
}

// A connection whose registration KIFF refuses leaves nothing behind.
func TestRegistrationFailureUndoes(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	// Another tool already holds the name list_charges in KIFF, so the
	// second registration is refused after the first one succeeded.
	e.kiff.Add(fakes.Connection{Connection: kiff.Connection{Name: "list_charges", Mode: "sealed"}, Tenant: tenant})
	id, digest := e.request()
	if st, out := e.approveAndSignIn(id, digest); st != 502 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "list_charges") {
		t.Fatalf("grant %s: %s", g.Status, g.Reason)
	}
	if _, ok := e.kiff.Connection(tenant, "create_refund"); ok {
		t.Fatal("create_refund stayed registered after the connection failed")
	}
	if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skToken(id)); ok {
		t.Fatal("the token was kept after a failed connection")
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the vendor still honours %d access and %d refresh tokens", a, r)
	}
}

// A broad vendor tool narrowed to one operation: KIFF governs the lifted
// amount, the vendor receives the pinned operation, and the agent can set
// nothing else.
func TestNarrowedTool(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.kiff.Cards["k_taf"] = map[string]int64{"stripe_refund": 500}
	refund := map[string]any{"tool": "stripe_api_write", "name": "stripe_refund", "amount_argument": "amount",
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "PostRefunds"},
		"arguments": map[string]any{
			"payment_intent": map[string]any{"path": "parameters.payment_intent", "type": "string", "required": true},
			"amount":         map[string]any{"path": "parameters.amount", "type": "integer", "required": true},
			"approval_token": map[string]any{"path": "human_confirmation.approval_token", "type": "string"},
		}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{refund}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	if text := out["text"].(string); !strings.Contains(text, `stripe_refund (stripe_api_write with stripe_api_operation_id = "PostRefunds"; the taf sets only amount, approval_token, payment_intent): moves money (amount)`) {
		t.Fatalf("text:\n%s", text)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	if st, out := e.approveAndSignIn(id, digest); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if c, ok := e.kiff.Connection(tenant, "stripe_refund"); !ok || c.Tool != "stripe_refund" || c.AmountArg != "amount" {
		t.Fatalf("registered: %+v %v", c.Connection, ok)
	}
	ctx := context.Background()
	res, err := e.kiff.Call(ctx, "k_taf", "stripe_refund", map[string]any{"payment_intent": "pi_1", "amount": 300})
	if err != nil || res.Decision != "allowed" || e.vendor.CallCount() != 1 {
		t.Fatalf("refund: %+v %v", res, err)
	}
	got, _ := json.Marshal(e.vendor.Calls[0])
	if want := `{"Tool":"stripe_api_write","Args":{"parameters":{"amount":300,"payment_intent":"pi_1"},"stripe_api_operation_id":"PostRefunds"}}`; string(got) != want {
		t.Fatalf("vendor received %s\nwant %s", got, want)
	}
	if res, _ := e.kiff.Call(ctx, "k_taf", "stripe_refund", map[string]any{"payment_intent": "pi_1", "amount": 900}); res.Decision != "refused" {
		t.Fatalf("over the Card: %+v", res)
	}
	for _, bad := range []map[string]any{
		{"payment_intent": "pi_1", "amount": 10, "stripe_api_operation_id": "PostPayouts"},
		{"payment_intent": "pi_1", "amount": "10"},
		{"amount": 10},
	} {
		res, _ := e.kiff.Call(ctx, "k_taf", "stripe_refund", bad)
		if !strings.Contains(string(res.Result), "Refused by planque") {
			t.Errorf("%v: %s", bad, res.Result)
		}
	}
	if e.vendor.CallCount() != 1 {
		t.Fatalf("vendor calls: %d, want 1", e.vendor.CallCount())
	}
	// A narrowing that would let the agent reach a pinned path is refused.
	refund["arguments"] = map[string]any{"op": map[string]any{"path": "stripe_api_operation_id", "type": "string"},
		"amount": map[string]any{"path": "parameters.amount", "type": "integer"}}
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{refund}}, nil, "k_setup"); st != 400 {
		t.Fatalf("overlapping narrowing: %d, want 400", st)
	}
}

// Two devices enrolling at once: one is enrolled, the other is told to try
// again, never silently replaced.
func TestEnrolmentRace(t *testing.T) {
	e := newEnv(t)
	var wg sync.WaitGroup
	wg.Add(2)
	e.s.KV = barrierKV{KV: e.kv, wg: &wg, prefix: skEnrolment}
	phones := []*fakes.Phone{e.device("dev_a"), e.device("dev_b")}
	statuses := make([]int, 2)
	var done sync.WaitGroup
	for i, p := range phones {
		done.Add(1)
		go func() {
			defer done.Done()
			statuses[i], _ = e.do("POST", "/v1/accounts/"+tenant+"/enrolment",
				map[string]string{"connector_key": "k_connector"}, p, "")
		}()
	}
	done.Wait()
	if statuses[0]+statuses[1] != 201+409 {
		t.Fatalf("statuses %v, want one 201 and one 409", statuses)
	}
	var en Enrolment
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &en); err != nil || len(en.Connector) == 0 {
		t.Fatalf("no enrolment after the race: %v", err)
	}
}

// An agent cannot make planque reach its own network: with the client
// planque runs with, a server on a local address is refused when the
// connection is made, and nothing it would answer comes back.
func TestGrantRefusesLocalServers(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.s.OAuth = &oauth.Client{HTTP: egress.Client(5 * time.Second), AllowHTTP: true}
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte("internal-secret"))
	}))
	defer internal.Close()
	for _, server := range []string{e.vendor.MCPURL(), internal.URL, "http://localhost:" + strings.TrimPrefix(internal.URL, "http://127.0.0.1:")} {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": server, "agent": "agt_taf",
			"tools": []map[string]any{{"tool": "list_charges"}}}, nil, "k_setup")
		msg, _ := out["error"].(string)
		if st != 502 || !strings.Contains(msg, "public addresses") || strings.Contains(msg, "internal-secret") {
			t.Errorf("%s: %d %s", server, st, msg)
		}
	}
}

// A sign-in link someone else opens connects nothing: the code goes to
// their browser and app, and only the phone that approved may confirm.
// It is the attack of sending "connect your Stripe" to a victim.
func TestSignInBoundToApprover(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	link := out["sign_in_url"].(string)
	// The victim opens the link: the vendor signs them in, and the code
	// lands in their browser's app.
	victimCode := e.signIn(link)
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusSignedIn || len(e.kiff.Connections(tenant)) != 0 {
		t.Fatalf("after the sign-in, before any confirmation: %s, %d tools in KIFF", g.Status, len(e.kiff.Connections(tenant)))
	}
	confirm := "/v1/accounts/" + tenant + "/grants/" + id + "/confirm"
	// The approver, without the code, cannot confirm.
	for _, guess := range []string{"", "guess", strings.Repeat("A", len(victimCode))} {
		if st, _ := e.do("POST", confirm, map[string]string{"confirm": guess}, e.phone, ""); st != 403 && st != 400 {
			t.Fatalf("confirm with %q: %d", guess, st)
		}
	}
	// The victim's own phone cannot confirm someone else's grant.
	if st, _ := e.do("POST", confirm, map[string]string{"confirm": victimCode}, fakes.NewPhone(), ""); st != 401 {
		t.Fatalf("another phone confirming: %d", st)
	}
	if len(e.kiff.Connections(tenant)) != 0 || e.vendor.CallCount() != 0 {
		t.Fatal("something was connected without the approving phone's confirmation")
	}
	// Approving again revokes the unconfirmed sign-in's token. (A second
	// later: the same signed request twice is refused as a replay.)
	time.Sleep(1100 * time.Millisecond)
	_, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if out["sign_in_url"] == nil {
		t.Fatalf("re-approve: %v", out)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the unconfirmed token still works: %d access, %d refresh", a, r)
	}
}

// A sign-in not confirmed in time is revoked, not connected.
func TestUnconfirmedSignInExpires(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	code := e.signIn(out["sign_in_url"].(string))
	// Time has passed since the sign-in.
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	g.ConfirmBy = time.Now().Add(-time.Minute)
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 409 {
		t.Fatalf("late confirm: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || len(e.kiff.Connections(tenant)) != 0 {
		t.Fatalf("a late confirmation: %s, %d tools", g.Status, len(e.kiff.Connections(tenant)))
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the late sign-in's token still works: %d access, %d refresh", a, r)
	}
}

// RFC 9207: a server that says it names itself must have named itself.
func TestCallbackRequiresIssWhenAdvertised(t *testing.T) {
	e := newEnv(t)
	e.vendor.AdvertiseIss = true
	e.enrol()
	id, digest := e.request()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	noFollow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, _ := noFollow.Get(out["sign_in_url"].(string))
	back, _ := url.Parse(resp.Header.Get("Location"))
	q := back.Query()
	q.Del("iss")
	back.RawQuery = q.Encode()
	if r, _ := http.Get(back.String()); r.StatusCode != 400 {
		t.Fatalf("answer without iss: %d", r.StatusCode)
	}
}

// A signature made for one planque is not valid at another.
func TestSignatureNamesThePlanque(t *testing.T) {
	e := newEnv(t)
	b, _ := json.Marshal(map[string]string{"device_key": e.phone.PublicKey(), "connector_key": "k_connector"})
	req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tenant+"/enrolment", bytes.NewReader(b))
	// Signed for https://other-planque.example, sent here.
	other, _ := http.NewRequest("POST", "https://other-planque.example/v1/accounts/"+tenant+"/enrolment", nil)
	e.phone.Sign(other, tenant, b)
	req.Header.Set(device.Header, other.Header.Get(device.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("a signature for another planque: %d, want 401", resp.StatusCode)
	}
}

func TestUnboundWriteToolTextSaysSo(t *testing.T) {
	text := grantText(Request{Server: "https://mcp.example.com", SignInAt: "https://as.example.com", Agent: "a",
		Tools: []ToolRequest{{Tool: "api_write", Name: "api_write"}}})
	if strings.Contains(text, "only within") || !strings.Contains(text, "can make any change this tool allows") {
		t.Fatalf("text:\n%s", text)
	}
}

// A sign-in nobody ever confirms does not keep a live token at the
// vendor: the next read of the grant revokes it, whoever reads.
func TestNobodyConfirms(t *testing.T) {
	for _, read := range []string{"status", "list", "export"} {
		t.Run(read, func(t *testing.T) {
			e := newEnv(t)
			e.enrol()
			id, digest := e.request()
			_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
			e.signIn(out["sign_in_url"].(string))
			if a, _ := e.vendor.Live(); a != 1 {
				t.Fatalf("live tokens after the sign-in: %d", a)
			}
			g, _ := e.s.loadGrant(context.Background(), tenant, id)
			g.ConfirmBy = time.Now().Add(-time.Minute)
			if err := e.s.saveGrant(context.Background(), g); err != nil {
				t.Fatal(err)
			}
			switch read {
			case "status":
				e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
			case "list":
				e.do("GET", "/v1/accounts/"+tenant+"/grants", nil, e.phone, "")
			case "export":
				e.do("GET", "/v1/accounts/"+tenant+"/export", nil, e.phone, "")
			}
			if a, r := e.vendor.Live(); a != 0 || r != 0 {
				t.Fatalf("after a %s read the token still works: %d access, %d refresh", read, a, r)
			}
			g, _ = e.s.loadGrant(context.Background(), tenant, id)
			if g.Status != StatusFailed || g.ConfirmHash != "" {
				t.Fatalf("grant %s", g.Status)
			}
			if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skToken(id)); ok {
				t.Fatal("the token is still stored")
			}
		})
	}
}

// Connect first (contract 5 on kiff/agents#91): a grant asked for with no
// tools is approved and signed in, and planque lists the vendor's tools
// without registering any. The asking agent reads the catalog and sets the
// tools, and the owner approves them with no second sign-in.
func TestConnectFirst(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st != 201 || !strings.Contains(out["text"].(string), "Only what it offers is read: its list of tools. None of them is called") {
		t.Fatalf("request: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	if strings.Contains(out["text"].(string), "agt_setup") {
		t.Fatalf("the owner's text names the asking agent's id: %s", out["text"])
	}
	if st, out := e.do("GET", acct+"/"+id+"/catalog", nil, nil, "k_setup"); st != 409 {
		t.Fatalf("a catalog before the sign-in: %d %v", st, out)
	}
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if n := len(e.kiff.Connections(tenant)); n != 0 {
		t.Fatalf("a listed grant registered %d tools in KIFF", n)
	}

	if st, _ := e.do("GET", acct+"/"+id+"/catalog", nil, nil, "k_taf"); st != 403 {
		t.Fatalf("another agent reading the catalog: %d", st)
	}
	st, out = e.do("GET", acct+"/"+id+"/catalog", nil, nil, "k_setup")
	if st != 200 || out["account"] != "acct_fake123" {
		t.Fatalf("catalog: %d %v", st, out)
	}
	names := []string{}
	for _, c := range out["tools"].([]any) {
		names = append(names, c.(map[string]any)["name"].(string))
		if c.(map[string]any)["input_schema"] == nil {
			t.Fatalf("a catalog tool without its schema: %v", c)
		}
	}
	if !slices.Contains(names, "create_refund") || !slices.Contains(names, "stripe_api_write") {
		t.Fatalf("catalog tools: %v", names)
	}
	if e.vendor.CallCount() != 0 {
		t.Fatal("listing called a tool")
	}

	refund := map[string]any{"tool": "create_refund", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	// Only the agent that asked sets the tools, only tools the vendor lists, and the
	// owner can't approve a listed grant before it has tools.
	if st, _ := e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_taf"); st != 403 {
		t.Fatalf("another agent setting the tools: %d", st)
	}
	if st, _ := e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "delete_account"}}}, nil, "k_setup"); st != 400 {
		t.Fatalf("a tool the vendor doesn't list: %d", st)
	}
	if st, _ := e.do("POST", acct+"/"+id+"/approve", map[string]string{"digest": digest}, e.phone, ""); st != 409 {
		t.Fatalf("approving a listed grant: %d", st)
	}
	st, out = e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending || out["digest"] == digest || !strings.Contains(out["text"].(string), "create_refund") ||
		!strings.Contains(out["text"].(string), "already signed in") || strings.Contains(out["text"].(string), "You will sign in") ||
		!strings.Contains(out["text"].(string), "acct_fake123") {
		t.Fatalf("set tools: %d %v", st, out)
	}
	toolsDigest := out["digest"].(string)
	if out["signed_in"] != true {
		t.Fatalf("a pending grant whose sign-in is held doesn't say so: %v", out)
	}
	if st, _ := e.do("POST", acct+"/"+id+"/approve", map[string]string{"digest": digest}, e.phone, ""); st != 409 {
		t.Fatalf("approving the old digest: %d", st)
	}

	// Approving the tools connects at once: no sign-in link, the tools are in KIFF.
	st, out = e.do("POST", acct+"/"+id+"/approve", map[string]string{"digest": toolsDigest}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] != nil {
		t.Fatalf("approve the tools: %d %v", st, out)
	}
	if _, ok := e.kiff.Connection(tenant, "create_refund"); !ok {
		t.Fatal("create_refund isn't in KIFF")
	}
	res, err := e.kiff.Call(context.Background(), "k_taf", "create_refund", map[string]any{"charge": "ch_1", "amount": 300})
	if err != nil || res.Decision != "allowed" || e.vendor.CallCount() != 1 {
		t.Fatalf("a refund inside the Card: %+v %v, vendor calls %d", res, err, e.vendor.CallCount())
	}
	if st, _ := e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup"); st != 409 {
		t.Fatalf("setting tools on a connected grant: %d", st)
	}
}

// A listed grant the owner disconnects revokes its held sign-in.
func TestDisconnectingAListedGrantRevokesItsSignIn(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	_, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	id := out["id"].(string)
	if st, out := e.approveAndSignIn(id, out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if a, _ := e.vendor.Live(); a != 1 {
		t.Fatalf("live access tokens after the sign-in: %d", a)
	}
	if st, out := e.do("DELETE", acct+"/"+id, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live at the vendor: %d access, %d refresh", a, r)
	}
}

// A request without a tools field is refused: connect first is an explicit [].
func TestMissingToolsIsNotConnectFirst(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf"}, nil, "k_setup"); st != 400 {
		t.Fatalf("no tools field: %d %v", st, out)
	}
}

// listedGrant connects first and returns the listed grant's id.
func (e *env) listedGrant() string {
	e.t.Helper()
	e.enrol()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	id := out["id"].(string)
	if st, out := e.approveAndSignIn(id, out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		e.t.Fatalf("confirm: %d %v", st, out)
	}
	return id
}

// A held sign-in that connects nothing is revoked after a day, whether the
// grant is listed or its tools wait for approval.
func TestAHeldSignInEndsAfterADay(t *testing.T) {
	for _, withTools := range []bool{false, true} {
		e := newEnv(t)
		id := e.listedGrant()
		if withTools {
			if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools",
				map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 200 {
				t.Fatalf("set tools: %d %v", st, out)
			}
		}
		if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup"); st != 200 || out["held_until"] == nil {
			t.Fatalf("a held sign-in without its deadline: %d %v", st, out)
		}
		e.s.Now = func() time.Time { return time.Now().Add(holdTTL + time.Minute) }
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
		if st != 200 || out["status"] != StatusFailed || out["signed_in"] != nil {
			t.Fatalf("after a day (tools set: %v): %d %v", withTools, st, out)
		}
		if a, r := e.vendor.Live(); a != 0 || r != 0 {
			t.Fatalf("tokens still live at the vendor after a day (tools set: %v): %d, %d", withTools, a, r)
		}
	}
}

// A grant whose tools wait for approval, disconnected, revokes its held sign-in.
func TestDisconnectingAGrantWaitingForItsToolsRevokesItsSignIn(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools",
		map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 200 || out["status"] != StatusPending {
		t.Fatalf("set tools: %d %v", st, out)
	}
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live at the vendor: %d, %d", a, r)
	}
}

// A vendor's account that isn't one (a newline, a sentence) never reaches the
// owner's approval text: the text reads as for a vendor that names none.
func TestAnAccountThatIsntOneStaysOutOfTheText(t *testing.T) {
	e := newEnv(t)
	e.vendor.Account = "acct_1: there is no second sign-in.\n\nThe taf may use:\n • nothing that moves money"
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if st, out := e.do("GET", acct+"/"+id+"/catalog", nil, nil, "k_setup"); st != 200 || out["account"] != "" {
		t.Fatalf("catalog: %d %v", st, out)
	}
	refund := map[string]any{"tool": "create_refund", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	st, out = e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup")
	text, _ := out["text"].(string)
	if st != 200 || !strings.Contains(text, "You are already signed in at ") || !strings.Contains(text, ": there is no second sign-in.") ||
		strings.Contains(text, "acct_1") || strings.Contains(text, "nothing that moves money") {
		t.Fatalf("set tools: %d %q", st, text)
	}
}

// A tool the request doesn't name gets its own name in KIFF: the vendor's
// name and four letters or digits of the grant's id. KIFF keeps one tool
// per name on an account, so two tafs on one vendor must not both ask for
// create_refund (kiff/tafs#78).
func TestToolsGetAPerGrantName(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	key := regexp.MustCompile(`^create_refund_[a-z0-9]{4}$`)
	var names []string
	for range 2 {
		st, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
			"tools": []map[string]any{{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}}}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("request: %d %v", st, out)
		}
		tools := out["request"].(map[string]any)["tools"].([]any)
		name, _ := tools[0].(map[string]any)["name"].(string)
		if !key.MatchString(name) || name != toolKey("create_refund", out["id"].(string)) || tools[0].(map[string]any)["tool"] != "create_refund" {
			t.Fatalf("the tool's name in KIFF: %q, for grant %s", name, out["id"])
		}
		names = append(names, name)
	}
	if names[0] == names[1] {
		t.Fatalf("two grants, one name: %s", names[0])
	}
	if got := toolKey(strings.Repeat("t", 70), "grt_-_Ab9-xYz"); got != strings.Repeat("t", 59)+"_ab9x" || !toolName.MatchString(got) {
		t.Fatalf("a long tool, an id with - and _: %q", got)
	}
}

// The owner never reads planque's name (kiff/tafs#82): not in what they
// approve, and not as the app a vendor shows when asking to allow the
// sign-in, which registers as tafs.
func TestTheOwnerReadsNoPlanque(t *testing.T) {
	for _, q := range []Request{
		{Server: "https://mcp.stripe.com", SignInAt: "https://access.stripe.com", Agent: "reembolsos-ab12cd34"},
		{Server: "https://mcp.stripe.com", SignInAt: "https://access.stripe.com", Agent: "reembolsos-ab12cd34",
			Tools: []ToolRequest{{Tool: "create_refund", Name: "create_refund_ab12", AmountArg: "amount"}}},
	} {
		for _, text := range []string{grantText(q), requestText(q, true, "acct_1")} {
			if strings.Contains(strings.ToLower(text), "planque") || !strings.Contains(text, "sealed vault") {
				t.Fatalf("the owner reads: %s", text)
			}
		}
	}
	if ClientName != "tafs" || pkClient("https://as", "https://cb", "tafs") == pkClient("https://as", "https://cb", "planque") {
		t.Fatal("the vendor shows tafs, and a new name registers again")
	}
}

// The hourly sweep revokes a held sign-in once its day is up, though nobody
// reads the grant again (kiff/kiff-cloud#1244 promises it within a day). A
// grant that moved on before its deadline is left as it is.
func TestTheSweepEndsAHeldSignInNobodyReads(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 0 {
		t.Fatalf("a sweep before the deadline ended %d, %v", n, err)
	}
	if st, out := e.do("POST", "/v1/sweep", nil, nil, ""); st != 200 || out["ended"] != float64(0) {
		t.Fatalf("POST /v1/sweep before the deadline: %d %v", st, out)
	}
	e.s.Now = func() time.Time { return time.Now().Add(holdTTL + time.Minute) }
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("no live token before the sweep")
	}
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("the sweep ended %d, %v; want 1", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live at the vendor after the sweep: %d, %d", a, r)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil || g.Status != StatusFailed {
		t.Fatalf("the grant after the sweep: %+v %v", g, err)
	}
	if n, _ := e.s.Sweep(t.Context()); n != 0 {
		t.Fatalf("a second sweep ended %d", n)
	}
	if items, _ := e.s.KV.Query(t.Context(), pkDue, ""); len(items) != 0 {
		t.Fatalf("the index kept %d entries", len(items))
	}
}

// A money tool names what it acts on and the read that says its currency
// (kiff/kiff-cloud#1299, kiff/tafs#167): both reach KIFF, the lookup by
// the read's name in KIFF, so KIFF counts only calls in the Card's currency.
func TestResourceAndUnitLookupReachKIFF(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	refund := map[string]any{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2,
		"resource_argument": "charge", "unit_lookup": map[string]any{"tool": "retrieve_charge", "argument": "charge", "path": "currency"}}
	read := map[string]any{"tool": "retrieve_charge", "read_only": true}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{refund, read}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	if st, out := e.approveAndSignIn(out["id"].(string), out["digest"].(string)); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	var refundName, readName string
	for _, c := range e.kiff.Connections(tenant) {
		switch c.Tool[:strings.LastIndex(c.Tool, "_")] {
		case "create_refund":
			refundName = c.Name
		case "retrieve_charge":
			readName = c.Name
		}
	}
	c, ok := e.kiff.Connection(tenant, refundName)
	if !ok || c.ResourceArg != "charge" || c.UnitLookup == nil || *c.UnitLookup != (kiff.UnitLookup{Tool: readName, Argument: "charge", Path: "currency"}) || readName == "" {
		t.Fatalf("registered: %+v (read %q)", c.Connection, readName)
	}

	for name, tools := range map[string][]any{
		"no unit": {map[string]any{"tool": "create_refund", "amount_argument": "amount", "resource_argument": "charge",
			"unit_lookup": map[string]any{"tool": "retrieve_charge", "argument": "charge", "path": "currency"}}, read},
		"not a read": {map[string]any{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "resource_argument": "charge",
			"unit_lookup": map[string]any{"tool": "list_charges", "argument": "charge", "path": "currency"}}, map[string]any{"tool": "list_charges"}},
		"not in the request": {map[string]any{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "resource_argument": "charge",
			"unit_lookup": map[string]any{"tool": "retrieve_charge", "argument": "charge", "path": "currency"}}},
		"a bad resource": {map[string]any{"tool": "create_refund", "amount_argument": "amount", "resource_argument": "amount"}},
		// Two reads of one vendor tool: the lookup must name its read by
		// name, as Stripe's narrowed reads share stripe_api_read.
		"a vendor tool asked twice": {map[string]any{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "resource_argument": "charge",
			"unit_lookup": map[string]any{"tool": "retrieve_charge", "argument": "charge", "path": "currency"}},
			map[string]any{"tool": "retrieve_charge", "name": "get_charge", "read_only": true},
			map[string]any{"tool": "retrieve_charge", "name": "get_charge_again", "read_only": true}},
	} {
		if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
			"tools": tools}, nil, "k_setup"); st != 400 {
			t.Errorf("%s: %d %v, want 400", name, st, out)
		}
	}
}

// The lookup names its read by name when two reads share a vendor tool.
func TestUnitLookupByTheReadsName(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	refund := map[string]any{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2,
		"resource_argument": "charge", "unit_lookup": map[string]any{"tool": "get_charge", "argument": "charge", "path": "currency"}}
	reads := []any{map[string]any{"tool": "retrieve_charge", "name": "get_charge", "read_only": true},
		map[string]any{"tool": "retrieve_charge", "name": "get_charge_again", "read_only": true}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": append([]any{refund}, reads...)}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	if st, out := e.approveAndSignIn(out["id"].(string), out["digest"].(string)); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	lookups := 0
	for _, c := range e.kiff.Connections(tenant) {
		if c.UnitLookup != nil {
			lookups++
			if c.UnitLookup.Tool != "get_charge" {
				t.Fatalf("lookup: %+v", c.UnitLookup)
			}
		}
	}
	if lookups != 1 {
		t.Fatalf("lookups = %d, want 1", lookups)
	}
}
