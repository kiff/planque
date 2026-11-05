package planque

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/kiff/planque/internal/kiff"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
)

// The real adapter and real planque/KIFF handlers run in process. Only the
// upstream service is synthetic; assertions inspect the outgoing wire request.
type connectorFixture struct {
	mu                                sync.Mutex
	t                                 *testing.T
	user, callback, accountStatus     string
	authScheme, authConfig            string
	links                             int
	foreign, changed, lost, strayLink bool
	// toolBody, when set, is the provider's answer for the pinned tool.
	toolBody                                string
	authorizations, executions, revocations int
	sent                                    map[string]json.RawMessage
	aliases                                 map[string]bool // names already used, refused again as upstream does
}

func (f *connectorFixture) RoundTrip(r *http.Request) (*http.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.URL.Host != "backend.composio.dev" || r.Header.Get("x-api-key") != "synthetic-project-key" {
		f.t.Fatal("wrong provider endpoint or project credential")
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v3.1")
	body := `{}`
	switch {
	case path == "/toolkits":
		if r.URL.Query().Get("cursor") == "next" {
			body = `{"items":[{"slug":"calendar_demo","name":"Calendar demo"}]}`
		} else {
			body = `{"items":[{"slug":"desk_demo","name":"Support desk"},{"slug":"unset_demo","name":"Not set up"}],"next_cursor":"next"}`
		}
	case path == "/auth_configs":
		scheme := f.authScheme
		if scheme == "" {
			scheme = "OAUTH2"
		}
		body = `{"items":[{"id":"ac_desk_demo","auth_scheme":"` + scheme + `","status":"ENABLED","toolkit":{"slug":"desk_demo"}},{"id":"ac_calendar_demo","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"calendar_demo"}},{"id":"ac_unset_demo","auth_scheme":"API_KEY","status":"DISABLED","toolkit":{"slug":"unset_demo"}}]}`
	case path == "/toolkits/desk_demo":
		body = `{"slug":"desk_demo","name":"Support desk"}`
	case path == "/tool_router/session":
		var in struct {
			User        string            `json:"user_id"`
			AuthConfigs map[string]string `json:"auth_configs"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		f.user = in.User
		f.authConfig = in.AuthConfigs["desk_demo"]
		f.authorizations++
		b, _ := json.Marshal(map[string]any{"session_id": "session_demo", "config": map[string]any{"user_id": f.user, "auth_configs": in.AuthConfigs}})
		body = string(b)
	case path == "/tool_router/session/session_demo/link":
		f.links++
		var in map[string]string
		_ = json.NewDecoder(r.Body).Decode(&in)
		if r.Method != "POST" || in["toolkit"] != "desk_demo" {
			f.t.Fatal("wrong hosted connection request")
		}
		f.callback = in["callback_url"]
		if f.aliases == nil {
			f.aliases = map[string]bool{}
		}
		if in["alias"] == "" || f.aliases[in["alias"]] {
			return profileResponse(400, `{"error":{"message":"alias already in use"}}`), nil
		}
		f.aliases[in["alias"]] = true
		body = `{"connected_account_id":"account_demo","redirect_url":"https://connect.composio.dev/link/demo"}`
		if f.strayLink {
			body = `{"connected_account_id":"account_demo","redirect_url":"https://elsewhere.example/link/demo"}`
		}
	case path == "/connected_accounts/account_demo" && r.Method == "DELETE":
		if r.URL.Query().Get("revoke_on_delete") != "true" {
			f.t.Error("no upstream revocation requested")
		}
		f.revocations++
		body = `{"success":true,"revoke_job_id":"job_demo"}`
	case path == "/connected_accounts/account_demo":
		user := f.user
		if f.foreign {
			user = "another_owner"
		}
		b, _ := json.Marshal(map[string]any{"id": "account_demo", "user_id": user, "status": f.accountStatus, "toolkit": map[string]string{"slug": "desk_demo"}})
		body = string(b)
	case path == "/tools" || path == "/tools/DESK_DEMO_REPLY":
		schema := `{"type":"object","properties":{"ticket":{"type":"string"},"text":{"type":"string"}},"required":["ticket","text"],"additionalProperties":false}`
		if f.changed {
			schema = `{"type":"object","properties":{"other":{"type":"number"}}}`
		}
		tool := `{"slug":"DESK_DEMO_REPLY","toolkit":{"slug":"desk_demo"},"version":"20261101_00","name":"Reply to a ticket","input_parameters":` + schema + `}`
		if path == "/tools" {
			body = `{"items":[` + tool + `]}`
		} else {
			if r.URL.Query().Get("version") != "20261101_00" {
				f.t.Error("version was not pinned")
			}
			body = tool
			if f.toolBody != "" {
				body = f.toolBody
			}
		}
	case path == "/tools/execute/DESK_DEMO_REPLY":
		if r.Method != "POST" {
			f.t.Fatal("unexpected execute method")
		}
		f.executions++
		_ = json.NewDecoder(r.Body).Decode(&f.sent)
		if len(f.sent) != 4 || string(f.sent["connected_account_id"]) != `"account_demo"` || string(f.sent["version"]) != `"20261101_00"` {
			f.t.Error("execution did not bind the approved account and version")
		}
		if f.lost {
			return nil, errors.New("synthetic lost answer")
		}
		body = `{"data":{"id":"reply_demo"},"successful":true}`
	default:
		f.t.Fatalf("unexpected provider request %s %s", r.Method, r.URL)
	}
	return profileResponse(200, body), nil
}
func newConnectorFixture(t *testing.T) (*env, *connectorFixture) {
	e := newInProcessEnv(t)
	e.enrol()
	f := &connectorFixture{t: t, accountStatus: "ACTIVE"}
	base := e.client.Transport
	e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
			return base.RoundTrip(r)
		}
		var in struct {
			Profile  json.RawMessage      `json:"profile"`
			Approval kiff.ProfileApproval `json:"approval"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			t.Fatal(err)
		}
		message, _ := base64.StdEncoding.DecodeString(in.Approval.Message)
		if !e.phone.Verifies(message, in.Approval.Signature) || in.Approval.DeviceID != e.phone.DeviceID || r.Header.Get("Authorization") != "Bearer k_connector" {
			t.Fatal("profile application lacks signed device evidence")
		}
		hash, _ := profile.Digest(in.Profile)
		return profileResponse(200, `{"revision":1,"profile_sha256":"`+hash+`","words_sha256":"`+strings.Repeat("a", 64)+`"}`), nil
	})
	e.s.Connectors = map[string]connector.Provider{"composio": &composio.Client{HTTP: &http.Client{Transport: f}, Secret: func(context.Context) (string, error) { return "synthetic-project-key", nil }}}
	return e, f
}
func requestConnector(t *testing.T, e *env) map[string]any {
	t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 || out["needs_key"] != false {
		t.Fatalf("request: %d %v", st, out)
	}
	return out["grant"].(map[string]any)
}
func listConnector(t *testing.T, e *env, f *connectorFixture) string {
	t.Helper()
	v := requestConnector(t, e)
	id := v["id"].(string)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve sign-in: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
	if st, _ := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.device("dev_other"), ""); st != 403 {
		t.Fatalf("another device confirmed: %d", st)
	}
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	return id
}
func connectConnector(t *testing.T, e *env, f *connectorFixture) *Grant {
	t.Helper()
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tool catalog: %d %v", st, catalog)
	}
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Name: "reply_ticket"}}, "profile": p}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("set tools: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if g.Request.Connector.Account != "account_demo" || len(g.Request.Connector.Tools) != 1 {
		t.Fatal("account and definition absent from signed request")
	}
	changed := g.Request
	binding := *changed.Connector
	changed.Connector = &binding
	binding.Account = "different_account"
	sum, _ := digest(changed)
	if sum == g.Digest {
		t.Fatal("account is not covered by approval digest")
	}
	ph, _ := profile.Digest(g.Request.Profile)
	approved := ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}
	if st, _ := e.do("POST", path+"/approve", map[string]any{"digest": g.Digest}, e.phone, ""); st != 409 {
		t.Fatalf("profile evidence not required: %d", st)
	}
	st, out = e.do("POST", path+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": approved}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve tools: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if !profileReady(g) || g.Tools[0].Connector == nil {
		t.Fatal("connection lacks applied profile or immutable tool snapshot")
	}
	return g
}
func TestConnectorCatalogAndApprovalBoundary(t *testing.T) {
	e, f := newConnectorFixture(t)
	path := "/v1/accounts/" + tenant + "/catalog"
	// One read walks the provider's pages until it has enough apps that can
	// be connected, or reaches the end.
	st, out := e.do("GET", path, nil, nil, "k_setup")
	if st != 200 || out["next_cursor"] != "" || out["provider"] != "composio" || len(out["entries"].([]any)) != 2 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	st, out = e.do("GET", path+"?provider=composio&cursor=next", nil, nil, "k_setup")
	if st != 200 || out["entries"].([]any)[0].(map[string]any)["vendor"] != "calendar_demo" {
		t.Fatalf("from a cursor: %d %v", st, out)
	}
	v := requestConnector(t, e)
	id := v["id"].(string)
	path = "/v1/accounts/" + tenant + "/grants/" + id
	if v["sign_in_url"] != nil || f.authorizations != 0 || f.executions != 0 {
		t.Fatal("sign-in or action occurred before approval")
	}
	if st, _ := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, nil, ""); st != 401 {
		t.Fatalf("unsigned: %d", st)
	}
	if st, _ := e.do("POST", path+"/approve", map[string]any{"digest": strings.Repeat("0", 64)}, e.phone, ""); st != 409 {
		t.Fatalf("wrong digest: %d", st)
	}
	if f.authorizations != 0 {
		t.Fatal("invalid approval contacted auth API")
	}
	if st, _ := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d", st)
	}
	f.foreign = true
	u, _ := url.Parse(f.callback)
	st, _ = e.do("GET", u.RequestURI()+"&status=ACTIVE", nil, nil, "")
	if st == 302 || st == 200 {
		t.Fatal("callback accepted another user's account")
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusApproved || g.SignedIn {
		t.Fatal("foreign callback authenticated")
	}
}
func TestConnectorSignedGrantPermitReplayAndDisconnect(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"operation_demo"}`)
	if res := e.s.call(ctx, g, "reply_ticket", args, ""); !strings.Contains(mustJSONString(res), "Refused") || f.executions != 0 {
		t.Fatal("missing permit executed")
	}
	mint := func(op string, raw json.RawMessage) string {
		return e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: op, Args: raw})
	}
	permit := mint("operation_demo", args)
	first := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if f.executions != 1 || !strings.Contains(first, "reply_demo") {
		t.Fatalf("call: %s executions=%d", first, f.executions)
	}
	var sent map[string]any
	_ = json.Unmarshal(f.sent["arguments"], &sent)
	if len(sent) != 2 || sent["ticket"] != "ticket_demo" || sent["text"] != "hello" {
		t.Fatalf("wrong exact args: %v", sent)
	}
	if second := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); second != first || f.executions != 1 {
		t.Fatal("replay resent or changed result")
	}
	// A changed definition pauses the grant: TestAChangedConnectorToolPausesItsGrant.
	for _, failure := range []string{"foreign", "expired"} {
		f.foreign = failure == "foreign"
		f.changed = failure == "changed"
		f.accountStatus = "ACTIVE"
		if failure == "expired" {
			f.accountStatus = "EXPIRED"
		}
		raw := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + failure + `"}`)
		result := mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, mint(failure, raw)))
		if f.executions != 1 || !strings.Contains(result, "Not sent") {
			t.Fatalf("%s: %s sends=%d", failure, result, f.executions)
		}
	}
	f.foreign, f.changed = false, false
	f.accountStatus = "ACTIVE"
	f.lost = true
	raw := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"lost"}`)
	p := mint("lost", raw)
	result := mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, p))
	again := mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, p))
	if f.executions != 2 || result != again || !strings.Contains(result, "will not be sent again") {
		t.Fatal("lost answer was retried or not retained")
	}
	st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, "")
	if st != 200 || out["status"] != StatusDisconnected || f.revocations != 1 {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(ctx, tenant, g.ID)
	if g.Connector.Revocation == nil || !g.Connector.Revocation.Disabled || !g.Connector.Revocation.Pending || !strings.Contains(mustJSONString(out), "pending") {
		t.Fatal("revocation falsely completed")
	}
}

func TestConnectorConcurrentReplayAndAbandonedSignIn(t *testing.T) {
	t.Run("concurrent replay", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"concurrent"}`)
		p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "concurrent", Args: args})
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); e.s.call(context.Background(), g, "reply_ticket", args, p) }()
		}
		wg.Wait()
		if f.executions != 1 {
			t.Fatalf("concurrent replays sent %d times", f.executions)
		}
	})
	t.Run("abandoned sign-in", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		v := requestConnector(t, e)
		id := v["id"].(string)
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
		if st != 200 {
			t.Fatalf("approval: %d %v", st, out)
		}
		g, _ := e.s.loadGrant(context.Background(), tenant, id)
		e.s.Now = func() time.Time { return g.HeldUntil.Add(time.Second) }
		if _, err := e.s.Sweep(context.Background()); err != nil {
			t.Fatal(err)
		}
		expired, _ := e.s.loadGrant(context.Background(), tenant, id)
		if expired.Status != StatusFailed || f.revocations != 1 || expired.Connector.Revocation == nil || !expired.Connector.Revocation.Pending {
			t.Fatal("abandoned connection was not disabled with accurate revocation status")
		}
	})
}

// An approval that failed after the provider made its link can be approved
// again: each attempt names its connection afresh, since the provider keeps
// a name it has seen. A link to anywhere but its sign-in page is refused.
func TestConnectorApprovalRetriesAfterALinkWasMade(t *testing.T) {
	e, f := newConnectorFixture(t)
	v := requestConnector(t, e)
	path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string) + "/approve"
	f.strayLink = true
	if st, _ := e.do("POST", path, map[string]any{"digest": v["digest"]}, e.phone, ""); st == 200 {
		t.Fatal("a link off the provider's sign-in page was accepted")
	}
	f.strayLink = false
	// The phone signs the retry anew; its signature carries the second.
	time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
	st, out := e.do("POST", path, map[string]any{"digest": v["digest"]}, e.phone, "")
	if st != 200 || out["sign_in_url"] != "https://connect.composio.dev/link/demo" || f.links != 2 || len(f.aliases) != 2 {
		t.Fatalf("retry: %d %v, %d links", st, out, f.links)
	}
}

// A connection an account key asked (Connect apps) is finished by the taf it
// is for or the account's Setup: either may list it, read its catalog and
// propose its tools, and the owner still approves them on the phone. No
// other agent can.
func TestConnectorAnAccountKeysGrantIsFinishedByItsTafOrSetup(t *testing.T) {
	e, f := newConnectorFixture(t)
	// The account's own key (tafs cloud's Connect apps), the account's Setup
	// (Coach), the taf, and an agent that is neither.
	e.kiff.Keys["k_account"] = kiff.Me{APIKey: true, TenantID: tenant}
	e.kiff.Keys["k_coach"] = kiff.Me{APIKey: true, TenantID: tenant, AgentID: "agt_coach", Roles: []string{"setup_reader"}}
	e.kiff.Keys["k_taf"] = kiff.Me{APIKey: true, TenantID: tenant, AgentID: "agt_taf"}
	e.kiff.Keys["k_third"] = kiff.Me{APIKey: true, TenantID: tenant, AgentID: "agt_third"}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_account")
	if st != 201 {
		t.Fatalf("request with the account key: %d %v", st, out)
	}
	v := out["grant"].(map[string]any)
	id := v["id"].(string)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve sign-in: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	lists := func(key string) bool {
		_, out := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, key)
		for _, g := range out["grants"].([]any) {
			if g.(map[string]any)["id"] == id {
				return true
			}
		}
		return false
	}
	for _, key := range []string{"k_third"} {
		if st, _ := e.do("GET", path+"/catalog", nil, nil, key); st != 403 || lists(key) {
			t.Fatalf("%s, neither the taf nor Setup, read the grant: %d", key, st)
		}
	}
	for _, key := range []string{"k_taf", "k_coach"} {
		if st, _ := e.do("GET", path+"/catalog", nil, nil, key); st != 200 || !lists(key) {
			t.Fatalf("%s can't read the grant: %d", key, st)
		}
	}
	_, catalog := e.do("GET", path+"/catalog", nil, nil, "k_coach")
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	tools := map[string]any{"tools": []ToolRequest{{Tool: alias, Name: "reply_ticket"}}, "profile": p}
	if st, _ := e.do("POST", path+"/tools", tools, nil, "k_third"); st != 403 {
		t.Fatalf("a third agent proposed tools: %d", st)
	}
	if st, out := e.do("POST", path+"/tools", tools, nil, "k_coach"); st != 200 {
		t.Fatalf("Setup proposes tools: %d %v", st, out)
	}
	// Proposing granted nothing: the owner's approval on the phone is next.
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status == StatusConnected || g.Proposer != "agt_coach" || f.executions != 0 {
		t.Fatalf("proposing tools connected them: %s, proposer %q", g.Status, g.Proposer)
	}
}

// A managed connection's words name the app by its title and the taf by its
// name, and no connection provider: what the owner reads and signs. Without a
// name, the taf is named by its agent; a name that isn't one line is refused.
func TestConnectorGrantWordsNameTheAppAndTheTaf(t *testing.T) {
	e, _ := newConnectorFixture(t)
	path := "/v1/accounts/" + tenant + "/catalog/connector:composio:desk_demo/grants"
	st, out := e.do("POST", path, map[string]string{"agent": "agt_taf", "agent_name": "Refunds"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	g := out["grant"].(map[string]any)
	text, vendor := g["text"].(string), g["vendor"].(string)
	if vendor == "" || !strings.HasPrefix(text, "Connect "+vendor+" for Refunds. ") || strings.Contains(strings.ToLower(text), "composio") || strings.Contains(text, "desk_demo") {
		t.Fatalf("words: %q (vendor %q)", text, vendor)
	}
	st, out = e.do("POST", path, map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 || !strings.Contains(out["grant"].(map[string]any)["text"].(string), " for the taf agt_taf. ") {
		t.Fatalf("without a name: %d %v", st, out)
	}
	if st, out := e.do("POST", path, map[string]string{"agent": "agt_taf", "agent_name": "Refunds\nand more"}, nil, "k_setup"); st != 400 {
		t.Fatalf("a name on two lines: %d %v", st, out)
	}
}

// Only an app the provider project can sign someone in to is offered: one
// listed upstream without an enabled sign-in is neither in the catalog nor
// started.
func TestConnectorCatalogOffersOnlyAppsWithASignIn(t *testing.T) {
	e, _ := newConnectorFixture(t)
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog?provider=composio", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	var listed []string
	for _, v := range out["entries"].([]any) {
		listed = append(listed, v.(map[string]any)["entry"].(string))
	}
	if !slices.Contains(listed, "connector:composio:desk_demo") || slices.Contains(listed, "connector:composio:unset_demo") {
		t.Fatalf("listed: %v", listed)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:unset_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 404 {
		t.Fatalf("an app without a sign-in was started: %d %v", st, out)
	}
}
