//go:build live

// Package live checks planque end to end against real services: planque
// on AWS Lambda with DynamoDB and KMS, Stripe's MCP server in a sandbox,
// and the KIFF stand-in (internal/fakes/cmd/fakekiff) on Lambda until
// KIFF ships the connector key role. It is run step by step, because a
// person signs in at Stripe in the middle:
//
//	LIVE_STATE=… go test -tags live -run TestLiveSetup ./live     # keys for fakekiff
//	deploy/aws.sh up (twice: FAKEKIFF=1 first)                     # see docs/live.md
//	LIVE_STATE=… go test -tags live -run TestLiveConnect ./live   # prints the sign-in link, waits, confirms
//	LIVE_STATE=… go test -tags live -run TestLiveCalls ./live     # refunds, refusals, refresh, disconnect, export
//
// LIVE_STATE holds the phone's key and the stand-in's keys: keep it out of
// the repository. STRIPE_API_KEY must be a test-mode key (sk_test_).
package live

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

type state struct {
	Tenant string `json:"tenant"`
	Phone  []byte `json:"phone"` // PKCS#8
	// Device is the phone's id in KIFF: it signs in to KIFF as a device
	// first, and every signature names it (d=), kiff/planque#12.
	Device    string `json:"device"`
	Seed      []byte `json:"seed"`
	Connector string `json:"connector"`
	Setup     string `json:"setup"`
	Taf       string `json:"taf"`
	Grant     string `json:"grant,omitempty"`
	SignIn    string `json:"sign_in,omitempty"`
}

func path(t *testing.T) string {
	p := os.Getenv("LIVE_STATE")
	if p == "" {
		t.Skip("LIVE_STATE is not set")
	}
	return p
}

func load(t *testing.T) *state {
	b, err := os.ReadFile(path(t))
	if err != nil {
		t.Fatal(err)
	}
	var s state
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatal(err)
	}
	return &s
}

func (s *state) save(t *testing.T) {
	b, _ := json.MarshalIndent(s, "", " ")
	if err := os.WriteFile(path(t), b, 0o600); err != nil {
		t.Fatal(err)
	}
}

func (s *state) phone(t *testing.T) *fakes.Phone {
	k, err := x509.ParsePKCS8PrivateKey(s.Phone)
	if err != nil {
		t.Fatal(err)
	}
	if s.Device == "" {
		t.Fatal("LIVE_STATE has no device: sign the phone in to KIFF as a device and save its id (kiff/planque#12)")
	}
	return &fakes.Phone{Key: k.(*ecdsa.PrivateKey), DeviceID: s.Device}
}

func env(t *testing.T, k string) string {
	v := strings.TrimRight(os.Getenv(k), "/")
	if v == "" {
		t.Fatalf("%s is not set", k)
	}
	return v
}

func random(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// TestLiveSetup makes the phone, the stand-in's keys and its signing key,
// and writes the stand-in's configuration next to the state.
func TestLiveSetup(t *testing.T) {
	p := path(t)
	phone := fakes.NewPhone()
	der, _ := x509.MarshalPKCS8PrivateKey(phone.Key)
	seed := make([]byte, ed25519.SeedSize)
	_, _ = rand.Read(seed)
	s := &state{Tenant: "tnt_live_" + random(4), Phone: der, Seed: seed,
		Connector: "kc_" + random(24), Setup: "ka_setup_" + random(24), Taf: "ka_taf_" + random(24)}
	s.save(t)
	keys := map[string]kiff.Me{
		s.Connector: {APIKey: true, TenantID: s.Tenant, Roles: []string{"connector"}, Capabilities: []string{"read", "tools.connect"}},
		s.Setup:     {APIKey: true, TenantID: s.Tenant, AgentID: "agt_live_setup"},
		s.Taf:       {APIKey: true, TenantID: s.Tenant, AgentID: "agt_live_taf"},
	}
	// The taf's Card: refunds up to 5.00 a call; reading the account moves nothing.
	cards := map[string]map[string]int64{s.Taf: {"stripe_refund": 500, "stripe_payment_intent": 0}}
	kb, _ := json.Marshal(keys)
	cb, _ := json.Marshal(cards)
	cfg, _ := json.Marshal(map[string]string{"FAKEKIFF_SEED": base64.RawURLEncoding.EncodeToString(seed),
		"FAKEKIFF_KEYS": string(kb), "FAKEKIFF_CARDS": string(cb)})
	if err := os.WriteFile(p+".fakekiff.json", cfg, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Logf("state for %s written", s.Tenant)
}

func call(t *testing.T, method, u string, body any, sign func(*http.Request, []byte), bearer string) (int, map[string]any) {
	t.Helper()
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, u, bytes.NewReader(b))
	if sign != nil {
		sign(req, b)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if b != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := (&http.Client{Timeout: 90 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out map[string]any
	if json.Unmarshal(raw, &out) != nil {
		out = map[string]any{"raw": string(raw)}
	}
	return resp.StatusCode, out
}

// Stripe's OAuth tools name the account and the mode on every call; the
// grant pins both, so the taf can reach only this account, in test mode.
func refundTool(account string) map[string]any {
	return map[string]any{"tool": "stripe_api_write", "name": "stripe_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2,
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "PostRefunds", "stripe_context": account, "livemode": false},
		"arguments": map[string]any{
			"payment_intent": map[string]any{"path": "parameters.payment_intent", "type": "string", "required": true},
			"amount":         map[string]any{"path": "parameters.amount", "type": "integer", "required": true},
			"approval_token": map[string]any{"path": "human_confirmation.approval_token", "type": "string",
				"description": "Stripe's approval token, after the owner confirmed the refund at Stripe."},
		}}
}

func readTool(account string) map[string]any {
	return map[string]any{"tool": "stripe_api_read", "name": "stripe_payment_intent", "read_only": true,
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "GetPaymentIntentsIntent", "stripe_context": account, "livemode": false},
		"arguments":       map[string]any{"intent": map[string]any{"path": "parameters.id", "type": "string", "required": true}}}
}

// TestLiveConnect enrols the phone, has the setup agent ask for Stripe,
// approves on the phone, and waits for the owner to sign in at Stripe.
func TestLiveConnect(t *testing.T) {
	s := load(t)
	pq := env(t, "PLANQUE_URL")
	phone := s.phone(t)
	sign := func(r *http.Request, b []byte) { phone.Sign(r, s.Tenant, b) }
	acct := pq + "/v1/accounts/" + s.Tenant
	account := stripe(t, "GET", "/v1/account", nil)["id"].(string) // the owner's Stripe account
	if st, out := call(t, "POST", acct+"/enrolment", map[string]string{"connector_key": s.Connector}, sign, ""); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	st, out := call(t, "POST", acct+"/grants", map[string]any{"server": "https://mcp.stripe.com", "agent": "agt_live_taf",
		"tools": []any{refundTool(account), readTool(account)}}, nil, s.Setup)
	if st != 201 {
		t.Fatalf("grant request: %d %v", st, out)
	}
	s.Grant = out["id"].(string)
	t.Logf("the owner reads:\n%s", out["text"])
	// The phone checks the digest against the request it shows, then signs.
	st, out = call(t, "POST", acct+"/grants/"+s.Grant+"/approve", map[string]string{"digest": out["digest"].(string)}, sign, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	s.SignIn = out["sign_in_url"].(string)
	s.save(t)
	_ = os.WriteFile(path(t)+".signin", []byte(s.SignIn+"\n"), 0o600)
	t.Logf("sign in at Stripe (a sandbox): %s", s.SignIn)
	// The browser that signs in lands in the owner's app (here the
	// stand-in's /fake/app) with a code; the phone that approved confirms
	// with it.
	fk := env(t, "FAKEKIFF_URL")
	deadline := time.Now().Add(30 * time.Minute)
	for time.Now().Before(deadline) {
		_, g := call(t, "GET", acct+"/grants/"+s.Grant, nil, nil, s.Setup)
		switch g["status"] {
		case "signed_in":
			st, last := call(t, "GET", fk+"/fake/app/last?grant="+s.Grant, nil, nil, "")
			if st != 200 {
				t.Fatalf("signed in, but the app got no code: %d %v", st, last)
			}
			st, out := call(t, "POST", acct+"/grants/"+s.Grant+"/confirm", map[string]string{"confirm": last["confirm"].(string)}, sign, "")
			if st != 200 || out["status"] != "connected" {
				t.Fatalf("confirm: %d %v", st, out)
			}
			t.Logf("signed in, confirmed on the phone, connected: tools %v", out["tools"])
			return
		case "failed":
			t.Fatalf("failed: %v", g["reason"])
		}
		time.Sleep(3 * time.Second)
	}
	t.Fatal("nobody signed in within 30 minutes")
}

// stripe calls Stripe's REST API with the test key, to check what the
// MCP calls did. It refuses anything but a test key.
func stripe(t *testing.T, method, p string, form url.Values) map[string]any {
	t.Helper()
	key := env(t, "STRIPE_API_KEY")
	if !strings.HasPrefix(key, "sk_test_") {
		t.Fatal("STRIPE_API_KEY must be a test-mode key")
	}
	req, _ := http.NewRequest(method, "https://api.stripe.com"+p, strings.NewReader(form.Encode()))
	req.SetBasicAuth(key, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if resp.StatusCode != 200 {
		t.Fatalf("stripe %s %s: %d %v", method, p, resp.StatusCode, out["error"])
	}
	if lm, _ := out["livemode"].(bool); lm {
		t.Fatal("a live-mode object: stop")
	}
	return out
}

// gateway is the taf calling a tool through the KIFF stand-in.
func gateway(t *testing.T, s *state, tool string, args map[string]any, parallel int) []map[string]any {
	t.Helper()
	st, out := call(t, "POST", env(t, "FAKEKIFF_URL")+"/fake/call", map[string]any{"tool": tool, "arguments": args, "parallel": parallel}, nil, s.Taf)
	if st != 200 {
		t.Fatalf("gateway: %d %v", st, out)
	}
	var res []map[string]any
	b, _ := json.Marshal(out["results"])
	_ = json.Unmarshal(b, &res)
	for _, e := range out["errors"].([]any) {
		if e != "" {
			t.Fatalf("gateway forward: %v", e)
		}
	}
	return res
}

func text(r map[string]any) string {
	b, _ := json.Marshal(r["result"])
	return string(b)
}

func aws(service, target, typ string) *awsapi.Client {
	return &awsapi.Client{Service: service, Region: "eu-west-1", Target: target, JSONType: typ, Creds: awsapi.FromEnv}
}

// TestLiveCalls runs the calls the issue lists, then disconnects.
func TestLiveCalls(t *testing.T) {
	s := load(t)
	pq, fk := env(t, "PLANQUE_URL"), env(t, "FAKEKIFF_URL")
	ctx := context.Background()
	acct := pq + "/v1/accounts/" + s.Tenant
	_, g := call(t, "GET", acct+"/grants/"+s.Grant, nil, nil, s.Setup)
	if g["status"] != "connected" {
		t.Fatalf("grant is %v", g["status"])
	}
	audience := g["audience"].(string)

	// A payment to refund, made with the test key.
	pi := stripe(t, "POST", "/v1/payment_intents", url.Values{"amount": {"2000"}, "currency": {"eur"},
		"payment_method": {"pm_card_visa"}, "confirm": {"true"}, "automatic_payment_methods[enabled]": {"true"},
		"automatic_payment_methods[allow_redirects]": {"never"}, "description": {"planque live check"}})
	piID := pi["id"].(string)

	// The sign-in must be the same Stripe environment as the test key:
	// the payment it just made is readable through planque.
	read := map[string]any{"intent": piID}
	res := gateway(t, s, "stripe_payment_intent", read, 1)[0]
	if !strings.Contains(text(res), piID) || !strings.Contains(text(res), "succeeded") {
		t.Fatalf("the Stripe sign-in is another environment than the test key's: %.400s", text(res))
	}
	t.Logf("read %s through KIFF → planque → Stripe MCP: same test environment as the test key", piID)
	refunds := func() int {
		return len(stripe(t, "GET", "/v1/refunds?payment_intent="+piID, nil)["data"].([]any))
	}

	// Outside the Card (8.00 > 5.00): KIFF mints no permit, nothing is sent.
	if r := gateway(t, s, "stripe_refund", map[string]any{"payment_intent": piID, "amount": 800}, 1)[0]; r["decision"] != "refused" {
		t.Fatalf("outside the Card: %v", r)
	}
	if refunds() != 0 {
		t.Fatal("a refund outside the Card reached Stripe")
	}
	t.Log("refund of 8.00 outside the Card: refused by KIFF, nothing sent")

	// Inside the Card (3.00).
	r := gateway(t, s, "stripe_refund", map[string]any{"payment_intent": piID, "amount": 300}, 1)[0]
	t.Logf("refund of 3.00 inside the Card: %v, Stripe answered %.600s", r["decision"], text(r))
	writeFile(t, ".refund", text(r))

	// Permits a gateway gone wrong might send: none reaches Stripe.
	seed := ed25519.NewKeyFromSeed(s.Seed)
	k := &fakes.KIFF{Issuer: fk, Key: seed, KID: "live-1"}
	mint := func(op string, args map[string]any, p fakes.Permit) string {
		args["planque_op"] = op
		b, _ := json.Marshal(args)
		p.Op, p.Args, p.Sub = op, b, "agt_live_taf"
		if p.Audience == "" {
			p.Audience = audience
		}
		if p.Tenant == "" {
			p.Tenant = s.Tenant
		}
		if p.Tool == "" {
			p.Tool = "stripe_refund"
		}
		return k.Mint(p)
	}
	forward := func(args map[string]any, permit string) string {
		b, _ := json.Marshal(args)
		st, out := call(t, "POST", fk+"/fake/forward", map[string]any{"name": "stripe_refund", "arguments": json.RawMessage(b), "permit": permit}, nil, s.Taf)
		if st != 200 {
			t.Fatalf("forward: %d %v", st, out)
		}
		rb, _ := json.Marshal(out["result"])
		return string(rb)
	}
	before := refunds()
	args := func(op string, amount int) map[string]any {
		return map[string]any{"payment_intent": piID, "amount": amount, "planque_op": op}
	}
	for _, c := range []struct {
		name, want string
		args       map[string]any
		permit     string
	}{
		{"no permit", "no execution permit", args("op_live_1", 100), ""},
		{"other arguments", "other arguments", args("op_live_2", 1900), mint("op_live_2", map[string]any{"payment_intent": piID, "amount": 100}, fakes.Permit{})},
		{"other account", "another KIFF account", args("op_live_3", 100), mint("op_live_3", map[string]any{"payment_intent": piID, "amount": 100}, fakes.Permit{Tenant: "tnt_other"})},
		{"other audience", "another audience", args("op_live_4", 100), mint("op_live_4", map[string]any{"payment_intent": piID, "amount": 100}, fakes.Permit{Audience: pq + "/mcp/x/y"})},
		{"expired", "expired", args("op_live_5", 100), mint("op_live_5", map[string]any{"payment_intent": piID, "amount": 100}, fakes.Permit{IssuedAt: time.Now().Add(-3 * time.Minute)})},
	} {
		if out := forward(c.args, c.permit); !strings.Contains(out, "Refused by planque") || !strings.Contains(out, c.want) {
			t.Errorf("%s: %s", c.name, out)
		} else {
			t.Logf("%s: refused, nothing sent", c.name)
		}
	}
	if refunds() != before {
		t.Fatal("a refused permit reached Stripe")
	}
	// Replay: one permit presented twice runs once.
	p := mint("op_live_replay", map[string]any{"intent": piID}, fakes.Permit{Tool: "stripe_payment_intent"})
	first := forwardTool(t, fk, s, "stripe_payment_intent", map[string]any{"intent": piID, "planque_op": "op_live_replay"}, p)
	again := forwardTool(t, fk, s, "stripe_payment_intent", map[string]any{"intent": piID, "planque_op": "op_live_replay"}, p)
	if first != again {
		t.Fatal("a replayed permit got a different answer")
	}
	t.Log("replayed permit: the first answer, not a second call")

	// Refresh: the stored token is marked expired; parallel calls on
	// several Lambda instances refresh it once.
	kv := &store.Dynamo{Table: "planque", API: aws("dynamodb", "DynamoDB_20120810", "application/x-amz-json-1.0")}
	sealer := &seal.KMS{KeyID: "alias/planque", API: aws("kms", "TrentService", "application/x-amz-json-1.1")}
	pk, sk := "TENANT#"+s.Tenant, "TOKEN#"+s.Grant
	it, ok, err := kv.Get(ctx, pk, sk)
	if err != nil || !ok {
		t.Fatalf("token record: %v %v", ok, err)
	}
	var rec struct {
		Sealed []byte    `json:"sealed"`
		Expiry time.Time `json:"expiry"`
	}
	_ = json.Unmarshal(it.Doc, &rec)
	binding := map[string]string{"tenant": s.Tenant, "grant": s.Grant, "kind": "token"}
	plain, err := sealer.Open(ctx, rec.Sealed, binding)
	if err != nil {
		t.Fatal(err)
	}
	var tok struct {
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}
	_ = json.Unmarshal(plain, &tok)
	t.Logf("token sealed with KMS (%d bytes), version %d, expires %s, refresh token held: %v",
		len(rec.Sealed), it.Ver, rec.Expiry.Format(time.RFC3339), tok.Refresh != "")
	if tok.Access == "" {
		t.Fatal("the sealed token did not open to an access token")
	}
	oldAccess, oldRefresh := tok.Access, sha256.Sum256([]byte(tok.Refresh))
	rec.Expiry = time.Now().Add(-time.Minute)
	it.Doc, _ = json.Marshal(rec)
	ver := it.Ver
	it.Ver++
	if err := kv.Put(ctx, it, store.Cond{Ver: &ver}); err != nil {
		t.Fatal(err)
	}
	for _, r := range gateway(t, s, "stripe_payment_intent", read, 6) {
		if !strings.Contains(text(r), piID) {
			t.Fatalf("after refresh: %.300s", text(r))
		}
	}
	it2, _, _ := kv.Get(ctx, pk, sk)
	_ = json.Unmarshal(it2.Doc, &rec)
	plain, _ = sealer.Open(ctx, rec.Sealed, binding)
	_ = json.Unmarshal(plain, &tok)
	if it2.Ver != it.Ver+1 || tok.Access == oldAccess {
		t.Fatalf("refreshes: version %d → %d, want exactly one refresh", it.Ver, it2.Ver)
	}
	t.Logf("6 parallel calls, one refresh (version %d → %d); refresh token rotated: %v",
		it.Ver, it2.Ver, sha256.Sum256([]byte(tok.Refresh)) != oldRefresh)
	access := tok.Access

	// Disconnect on the phone: tools leave KIFF, the token is revoked at
	// Stripe and deleted here.
	phone := s.phone(t)
	sign := func(r *http.Request, b []byte) { phone.Sign(r, s.Tenant, b) }
	st, out := call(t, "DELETE", acct+"/grants/"+s.Grant, nil, sign, "")
	if st != 200 || out["status"] != "disconnected" {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	evs, _ := json.Marshal(out["events"])
	t.Logf("disconnect record: %s", evs)
	if _, ok, _ := kv.Get(ctx, pk, sk); ok {
		t.Fatal("the token is still in DynamoDB")
	}
	_, conns := call(t, "GET", fk+"/v1/connections", nil, nil, s.Connector)
	if conns["connections"] != nil {
		t.Fatalf("KIFF still lists %v", conns["connections"])
	}
	req, _ := http.NewRequest("POST", "https://mcp.stripe.com", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"planque-live","version":"1"}}}`))
	req.Header.Set("Authorization", "Bearer "+access)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	t.Logf("Stripe answers the revoked access token with %d", resp.StatusCode)
	if resp.StatusCode != 401 {
		t.Errorf("Stripe still accepts the access token after disconnect: %d", resp.StatusCode)
	}

	// The oversight record, checked offline.
	st, out = call(t, "GET", acct+"/export", nil, sign, "")
	if st != 200 {
		t.Fatalf("export: %d", st)
	}
	b, _ := json.MarshalIndent(out, "", " ")
	writeFile(t, ".export.json", string(b))
	grants := out["grants"].([]any)
	a := grants[0].(map[string]any)["approval"].(map[string]any)
	if a["device_id"] != s.Device || !phone.Verifies([]byte(a["message"].(string)), a["signature"].(string)) {
		t.Fatal("the exported approval doesn't name the phone, or doesn't verify with its key")
	}
	t.Logf("export: %d grant(s); the approval names the phone's KIFF device and verifies with its key", len(grants))
}

func forwardTool(t *testing.T, fk string, s *state, name string, args map[string]any, permit string) string {
	b, _ := json.Marshal(args)
	st, out := call(t, "POST", fk+"/fake/forward", map[string]any{"name": name, "arguments": json.RawMessage(b), "permit": permit}, nil, s.Taf)
	if st != 200 {
		t.Fatalf("forward: %d %v", st, out)
	}
	rb, _ := json.Marshal(out["result"])
	return string(rb)
}

func writeFile(t *testing.T, suffix, s string) {
	_ = os.WriteFile(path(t)+suffix, []byte(s), 0o600)
}

var _ = fmt.Sprint

// TestLiveDisconnect disconnects the current grant on the phone.
func TestLiveDisconnect(t *testing.T) {
	s := load(t)
	phone := s.phone(t)
	st, out := call(t, "DELETE", env(t, "PLANQUE_URL")+"/v1/accounts/"+s.Tenant+"/grants/"+s.Grant, nil,
		func(r *http.Request, b []byte) { phone.Sign(r, s.Tenant, b) }, "")
	if st != 200 || out["status"] != "disconnected" {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	evs := out["events"].([]any)
	t.Logf("%v", evs[len(evs)-1])
}
