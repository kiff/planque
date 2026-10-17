package planque

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/device"
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
	vendor *fakes.Vendor
	kiff   *fakes.KIFF
	phone  *fakes.Phone
	s      *Server
	url    string
	kv     *store.Memory
}

func newEnv(t *testing.T) *env {
	t.Helper()
	e := &env{t: t, vendor: fakes.NewVendor(), kiff: fakes.NewKIFF(), phone: fakes.NewPhone(), kv: &store.Memory{}}
	t.Cleanup(e.vendor.Close)
	ks := httptest.NewServer(e.kiff.Handler())
	t.Cleanup(ks.Close)
	e.kiff.Issuer = ks.URL
	e.kiff.Keys = map[string]kiff.Me{
		"k_connector": {APIKey: true, TenantID: tenant, Roles: []string{"connector"}, Capabilities: []string{"read", "tools.connect"}},
		"k_admin":     {APIKey: true, TenantID: tenant, Roles: []string{"admin"}, Capabilities: []string{"read", "tools.connect", "keys.manage", "cards.manage"}},
		"k_setup":     {APIKey: true, TenantID: tenant, AgentID: "agt_setup"},
		"k_taf":       {APIKey: true, TenantID: tenant, AgentID: "agt_taf"},
		"k_other":     {APIKey: true, TenantID: "tnt_other", AgentID: "agt_x"},
	}
	// The taf's Card: refunds up to 500 a call; listing moves nothing.
	e.kiff.Cards["k_taf"] = map[string]int64{"create_refund": 500, "list_charges": 0}
	var h http.Handler
	ps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }))
	t.Cleanup(ps.Close)
	e.url = ps.URL
	e.s = e.server()
	h = e.s.Handler()
	return e
}

// server is one planque process on the shared table, as on Lambda.
func (e *env) server() *Server {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	return &Server{BaseURL: e.url, KV: e.kv, Sealer: &seal.Local{Key: key},
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
		"tools": []map[string]any{{"tool": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2},
			{"tool": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		e.t.Fatalf("request: %d %v", st, out)
	}
	return out["id"].(string), out["digest"].(string)
}

// connect runs the whole connection and returns the grant.
func (e *env) connect() *Grant {
	e.t.Helper()
	e.enrol()
	id, digest := e.request()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 200 {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	resp, err := http.Get(out["sign_in_url"].(string)) // the owner signs in at the vendor
	if err != nil {
		e.t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || !strings.Contains(string(page), "Connected") {
		e.t.Fatalf("callback: %d %s", resp.StatusCode, page)
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
	// Another phone cannot replace the enrolled one.
	if st, _ := e.do("POST", path, map[string]string{"device_key": other.PublicKey(), "connector_key": "k_connector"}, other, ""); st != 409 {
		t.Fatalf("second phone: %d, want 409", st)
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
	if err != nil || res.Decision != "allowed" || !strings.Contains(string(res.Result), "re_1") {
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

// barrierKV holds the first two reads of an operation's record until
// both happened, so two presentations race to claim the same op.
type barrierKV struct {
	store.KV
	wg *sync.WaitGroup
}

func (b barrierKV) Get(ctx context.Context, pk, sk string) (store.Item, bool, error) {
	it, ok, err := b.KV.Get(ctx, pk, sk)
	if strings.HasPrefix(sk, "OP#") && !ok {
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
	e.s.KV = barrierKV{KV: e.kv, wg: &wg}
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
	pub, err := x509.ParsePKIXPublicKey(a.DeviceKey)
	if err != nil {
		t.Fatal(err)
	}
	sig, _ := base64.StdEncoding.DecodeString(a.Signature)
	sum := sha256.Sum256([]byte(a.Message))
	if !ecdsa.VerifyASN1(pub.(*ecdsa.PublicKey), sum[:], sig) {
		t.Fatal("the exported signature does not verify")
	}
	bodySum := sha256.Sum256([]byte(a.Body))
	if lines := strings.Split(a.Message, "\n"); lines[4] != fmt.Sprintf("%x", bodySum) {
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
	for _, want := range []string{"requested by agt_setup", "approved on the phone", "signed in at", "registered in KIFF", "disconnected on the phone"} {
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
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	resp, _ := http.Get(out["sign_in_url"].(string))
	if resp.StatusCode != 502 {
		t.Fatalf("callback: %d", resp.StatusCode)
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
	_, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if resp, _ := http.Get(out["sign_in_url"].(string)); resp.StatusCode != 200 {
		t.Fatalf("callback: %d", resp.StatusCode)
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
