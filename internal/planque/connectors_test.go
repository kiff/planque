package planque

import (
	"cmp"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/kiff"

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
	toolBody string
	// executeStatus and executeBody, when set, are the provider's answer to an execution.
	executeStatus                           int
	executeBody                             string
	authorizations, executions, revocations int
	// setUps are the apps whose sign-in planque asked the provider to set up, in order.
	setUps  []string
	sent    map[string]json.RawMessage
	aliases map[string]bool // names already used, refused again as upstream does
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
			// unset_demo: the provider signs people in by itself, so it can be set up on
			// demand; own_demo needs our own developer app first; the ondemand_N apps
			// are there to count against the day's allowance.
			items := `{"slug":"desk_demo","name":"Support desk"},{"slug":"unset_demo","name":"Not set up","composio_managed_auth_schemes":["OAUTH2"]},{"slug":"own_demo","name":"Own app","auth_schemes":["OAUTH2"]}`
			for i := 1; i <= 12; i++ {
				items += fmt.Sprintf(`,{"slug":"ondemand_%d","name":"On demand %d","composio_managed_auth_schemes":["OAUTH2"]}`, i, i)
			}
			body = `{"items":[` + items + `],"next_cursor":"next"}`
		}
	case path == "/toolkits/unset_demo":
		body = `{"slug":"unset_demo","name":"Not set up","composio_managed_auth_schemes":["OAUTH2"]}`
	case path == "/toolkits/own_demo":
		body = `{"slug":"own_demo","name":"Own app","auth_schemes":["OAUTH2"]}`
	case strings.HasPrefix(path, "/toolkits/ondemand_"):
		slug := strings.TrimPrefix(path, "/toolkits/")
		body = `{"slug":"` + slug + `","name":"On demand","composio_managed_auth_schemes":["OAUTH2"]}`
	case path == "/auth_configs" && r.Method == "POST":
		var in struct {
			Toolkit struct {
				Slug string `json:"slug"`
			} `json:"toolkit"`
			Config map[string]any `json:"auth_config"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		if in.Config["type"] != "use_composio_managed_auth" || len(in.Config) != 1 {
			f.t.Fatalf("a set-up with something of ours in it: %v", in.Config)
		}
		f.setUps = append(f.setUps, in.Toolkit.Slug)
		body = `{"toolkit":{"slug":"` + in.Toolkit.Slug + `"},"auth_config":{"id":"ac_` + in.Toolkit.Slug + `_new","auth_scheme":"OAUTH2","is_composio_managed":true}}`
		return profileResponse(201, body), nil
	case path == "/auth_configs":
		scheme := f.authScheme
		if scheme == "" {
			scheme = "OAUTH2"
		}
		body = `{"items":[{"id":"ac_desk_demo","auth_scheme":"` + scheme + `","status":"ENABLED","toolkit":{"slug":"desk_demo"}},{"id":"ac_calendar_demo","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"calendar_demo"}},{"id":"ac_unset_demo","auth_scheme":"API_KEY","status":"DISABLED","toolkit":{"slug":"unset_demo"}}]}`
	case path == "/tool_router/session":
		var in struct {
			User        string            `json:"user_id"`
			AuthConfigs map[string]string `json:"auth_configs"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		f.user = in.User
		f.authConfig = in.AuthConfigs["desk_demo"]
		if f.authConfig == "" {
			f.authConfig = in.AuthConfigs["unset_demo"]
		}
		f.authorizations++
		b, _ := json.Marshal(map[string]any{"session_id": "session_demo", "config": map[string]any{"user_id": f.user, "auth_configs": in.AuthConfigs}})
		body = string(b)
	case path == "/tool_router/session/session_demo/link":
		f.links++
		var in map[string]string
		_ = json.NewDecoder(r.Body).Decode(&in)
		if r.Method != "POST" || in["toolkit"] != "desk_demo" && in["toolkit"] != "unset_demo" {
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
		// As the provider shows a vendor account: its name readable, its id and mode redacted.
		b, _ := json.Marshal(map[string]any{"id": "account_demo", "user_id": user, "status": f.accountStatus, "toolkit": map[string]string{"slug": "desk_demo"},
			"data": map[string]any{"displayName": "Demo Shop", "livemode": "REDACTED", "stripe_user_id": "REDACTED"}})
		body = string(b)
	case path == "/tools" && r.URL.Query().Get("query") != "":
		// A search by words: a tool of a set-up app, one of an app that needs setup,
		// and a platform meta tool that must never come back.
		if r.URL.Query().Has("toolkit_slug") {
			f.t.Fatal("a search for several apps was narrowed upstream to one")
		}
		body = `{"items":[{"slug":"DESK_DEMO_REPLY","toolkit":{"slug":"desk_demo"},"version":"20261101_00","name":"Reply to a ticket","description":"Replies to a ticket.","input_parameters":{"type":"object","properties":{"ticket":{"type":"string"}}}},` +
			`{"slug":"OWN_DEMO_SEND","toolkit":{"slug":"own_demo"},"version":"20261101_00","name":"Send a reply","input_parameters":{"type":"object","properties":{}}},` +
			`{"slug":"COMPOSIO_MULTI_EXECUTE_TOOL","toolkit":{"slug":"composio"},"version":"20261101_00","name":"Run anything","input_parameters":{"type":"object","properties":{}}}]}`
	case path == "/triggers_types":
		if r.URL.Query().Get("toolkit_slugs") != "desk_demo" {
			body = `{"items":[]}`
			break
		}
		body = `{"items":[{"slug":"DESK_DEMO_NEW_TICKET","name":"New ticket","description":"A ticket was opened.","version":"20261101_00","toolkit":{"slug":"desk_demo"},"config":{"type":"object","properties":{"queue":{"type":"string"}}},"payload":{"type":"object"}}]}`
	case path == "/toolkits/desk_demo":
		body = `{"slug":"desk_demo","name":"Support desk"}`
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
		if f.executeBody != "" {
			return profileResponse(cmp.Or(f.executeStatus, 200), f.executeBody), nil
		}
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
	return connectConnectorWith(t, e, f, profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}})
}

// connectConnectorWith is connectConnector with the profile given.
func connectConnectorWith(t *testing.T, e *env, f *connectorFixture, p profile.Profile) *Grant {
	t.Helper()
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tool catalog: %d %v", st, catalog)
	}
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
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

// The unsearched catalog offers the apps with a sign-in set up; a search
// finds every app the provider has, each saying how it stands. An app the
// provider can sign people in to by itself is set up on the owner's pick,
// once, and its grant pins that sign-in; one that needs our own developer
// credentials is listed as such and not started (kiff/planque#87).
func TestConnectorAnyAppIsSetUpOnTheOwnersPick(t *testing.T) {
	e, f := newConnectorFixture(t)
	catalog := "/v1/accounts/" + tenant + "/catalog"
	states := func(q string) map[string]string {
		t.Helper()
		st, out := e.do("GET", catalog+"?provider=composio"+q, nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("catalog%s: %d %v", q, st, out)
		}
		m := map[string]string{}
		for _, v := range out["entries"].([]any) {
			e := v.(map[string]any)
			m[e["entry"].(string)] = e["state"].(string)
		}
		return m
	}
	short := states("")
	if short["connector:composio:desk_demo"] != "available" || short["connector:composio:unset_demo"] != "" || short["connector:composio:own_demo"] != "" {
		t.Fatalf("the short list: %v", short)
	}
	found := states("&q=demo")
	if found["connector:composio:desk_demo"] != "available" || found["connector:composio:unset_demo"] != "available" || found["connector:composio:own_demo"] != "needs_setup" {
		t.Fatalf("a search: %v", found)
	}
	// An app that needs our own credentials isn't started.
	if st, out := e.do("POST", catalog+"/connector:composio:own_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 409 || !strings.Contains(out["error"].(string), "needs setting up") || len(f.setUps) != 0 {
		t.Fatalf("an app that needs setup was started: %d %v %v", st, out, f.setUps)
	}
	// One the provider signs people in to is set up on the pick, once for every pick after.
	st, out := e.do("POST", catalog+"/connector:composio:unset_demo/grants", map[string]string{"agent": "agt_taf", "agent_name": "Invoices"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("on demand: %d %v", st, out)
	}
	g := out["grant"].(map[string]any)
	if g["on_demand"] != true || g["custody"] != "provider" || !strings.HasPrefix(g["text"].(string), "Connect Not set up for Invoices. Your sign-in to Not set up stays with the connection service") {
		t.Fatalf("the on-demand grant: %v", g)
	}
	if st, out := e.do("POST", catalog+"/connector:composio:unset_demo/grants", map[string]string{"agent": "agt_other"}, nil, "k_setup"); st != 201 || out["grant"].(map[string]any)["on_demand"] != true {
		t.Fatalf("a second pick: %d %v", st, out)
	}
	if !slices.Equal(f.setUps, []string{"unset_demo"}) {
		t.Fatalf("the sign-in was set up %v times", f.setUps)
	}
	// It is offered without a search from now on, and its approval pins the sign-in that was made.
	if now := states(""); now["connector:composio:unset_demo"] != "connected" && now["connector:composio:unset_demo"] != "available" {
		t.Fatalf("offered after set-up: %v", now)
	}
	id := g["id"].(string)
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]any{"digest": g["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	if f.authConfig != "ac_unset_demo_new" {
		t.Fatalf("the hosted sign-in pinned %q", f.authConfig)
	}
	// A grant planque signs in to itself says so.
	st, out = e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
	if st != 200 || out["custody"] != "provider" {
		t.Fatalf("custody on a read: %d %v", st, out)
	}
}

// An account sets up at most ten apps on demand a day: each is a sign-in
// the provider sets up for everyone, against its allowance.
func TestConnectorOnDemandSetUpsAreCappedADay(t *testing.T) {
	e, f := newConnectorFixture(t)
	for i := 1; i <= onDemandADay; i++ {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_"+strconv.Itoa(i)+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("app %d: %d %v", i, st, out)
		}
	}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_11/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 429 || !strings.Contains(out["error"].(string), "10 new apps today") || len(f.setUps) != onDemandADay {
		t.Fatalf("the eleventh: %d %v (%d set up)", st, out, len(f.setUps))
	}
	// An app set up already doesn't count: a pick of one is a grant as any other.
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
		t.Fatalf("a pre-set app under the cap: %d %v", st, out)
	}
	// The next day, the allowance is new.
	e.s.Now = func() time.Time { return time.Now().Add(24 * time.Hour) }
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_11/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
		t.Fatalf("the next day: %d %v", st, out)
	}
}

// Coach finds tools by words across apps before connecting any: each tool
// comes with its app as it stands, a platform meta tool never does, and an
// app's events are listed the same way, before and after a connection
// (kiff/planque#87).
func TestConnectorToolsAndEventsAreFoundBeforeConnecting(t *testing.T) {
	e, f := newConnectorFixture(t)
	base := "/v1/accounts/" + tenant + "/catalog"
	if st, out := e.do("GET", base+"/tools", nil, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "q is the words") {
		t.Fatalf("no words: %d %v", st, out)
	}
	if st, _ := e.do("GET", base+"/tools?q=reply", nil, nil, ""); st != 401 {
		t.Fatalf("unauthenticated search: %d", st)
	}
	st, out := e.do("GET", base+"/tools?q=reply", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("search: %d %v", st, out)
	}
	tools := out["tools"].([]any)
	if len(tools) != 2 || out["skipped"].(float64) != 1 {
		t.Fatalf("found: %v", out)
	}
	desk, own := tools[0].(map[string]any), tools[1].(map[string]any)
	if desk["entry"] != "connector:composio:desk_demo" || desk["app_state"] != "available" || desk["app_title"] != "Support desk" || desk["slug"] != "DESK_DEMO_REPLY" || desk["description"] != "Replies to a ticket." {
		t.Fatalf("a set-up app's tool: %v", desk)
	}
	if own["entry"] != "connector:composio:own_demo" || own["app_state"] != "needs_setup" || own["app_title"] != "Own app" {
		t.Fatalf("an app that needs setup: %v", own)
	}
	for _, tl := range tools {
		if strings.HasPrefix(tl.(map[string]any)["slug"].(string), "COMPOSIO_") {
			t.Fatalf("a platform meta tool came back: %v", tl)
		}
	}
	// Narrowed to the apps named.
	st, out = e.do("GET", base+"/tools?q=reply&apps=own_demo,desk_demo", nil, nil, "k_setup")
	if st != 200 || len(out["tools"].([]any)) != 2 {
		t.Fatalf("named apps: %d %v", st, out)
	}
	// An app's events, before any connection.
	st, out = e.do("GET", base+"/connector:composio:desk_demo/events", nil, nil, "k_setup")
	if st != 200 || len(out["events"].([]any)) != 1 {
		t.Fatalf("events: %d %v", st, out)
	}
	ev := out["events"].([]any)[0].(map[string]any)
	if ev["slug"] != "DESK_DEMO_NEW_TICKET" || ev["description"] != "A ticket was opened." || ev["door"] != false || ev["version"] != "20261101_00" || !strings.Contains(fmt.Sprint(ev["config"]), "queue") {
		t.Fatalf("event: %v", ev)
	}
	if st, _ := e.do("GET", base+"/support-desk@1/events", nil, nil, "k_setup"); st != 404 {
		t.Fatalf("a reviewed entry's events: %d", st)
	}
	// The platform's own toolkit is no app: its events are refused, and the provider isn't asked.
	if st, out := e.do("GET", base+"/connector:composio:composio/events", nil, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "isn't an app") {
		t.Fatalf("the platform's own events: %d %v", st, out)
	}
	// And after: the connection's own events, for its agent only.
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id + "/events"
	st, out = e.do("GET", path, nil, nil, "k_setup")
	if st != 200 || len(out["events"].([]any)) != 1 {
		t.Fatalf("a connection's events: %d %v", st, out)
	}
	if st, _ := e.do("GET", path, nil, nil, "k_other"); st != 403 && st != 401 {
		t.Fatalf("another agent read the events: %d", st)
	}
}

// The day's allowance holds when picks arrive together: at most ten
// sign-ins are set up, whatever the order, and the provider's cached
// sign-ins are never written while another read holds them.
func TestConnectorOnDemandCapHoldsUnderARace(t *testing.T) {
	e, f := newConnectorFixture(t)
	var wg sync.WaitGroup
	results := make([]int, 13)
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			st, _ := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_"+strconv.Itoa(i)+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
			results[i] = st
		}(i)
	}
	// A catalog read at the same time holds the sign-ins the set-ups replace.
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.do("GET", "/v1/accounts/"+tenant+"/catalog?provider=composio&q=demand", nil, nil, "k_setup")
		}()
	}
	wg.Wait()
	made, refused := 0, 0
	for _, st := range results[1:] {
		switch st {
		case 201:
			made++
		case 429:
			refused++
		}
	}
	f.mu.Lock()
	setUps := len(f.setUps)
	f.mu.Unlock()
	if setUps > onDemandADay || made > onDemandADay || made+refused != 12 {
		t.Fatalf("%d set up, %d made, %d refused, cap %d", setUps, made, refused, onDemandADay)
	}
}

// A refused call tells the taf who refused it and why, so it can tell the
// owner: the vendor's own reason, the provider's with its status, or a lost
// answer with the provider's words, each still saying no claim is made about
// an effect. Before, every refusal was one fixed line, the same for any id.
func TestARefusedCallSaysWhoRefusedAndWhy(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	call := func(op string) string {
		t.Helper()
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: op, Args: args})
		return mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	}
	f.executeBody = `{"data":{},"successful":false,"error":"No such ticket: 'ticket_demo'"}`
	if got := call("op_vendor"); !strings.Contains(got, "refused it: No such ticket: 'ticket_demo'. No claim is made about whether an effect occurred.") || !strings.Contains(got, `"isError":true`) {
		t.Fatalf("the vendor's refusal: %s", got)
	}
	f.executeStatus, f.executeBody = 403, `{"error":{"message":"this account may not run that action"}}`
	if got := call("op_provider"); !strings.Contains(got, "The connection provider refused it (it answered 403): this account may not run that action. No claim is made") {
		t.Fatalf("the provider's refusal: %s", got)
	}
	f.executeStatus, f.executeBody = 502, `{"message":"upstream timed out"}`
	if got := call("op_unknown"); !strings.Contains(got, "The connector did not return a confirmed result (it answered 502): upstream timed out. This operation will not be sent again.") {
		t.Fatalf("no confirmed result: %s", got)
	}
}

// A managed connection's mode is what a permitted call showed, never a call
// of planque's own: unknown (treated as live) until an answer carries
// "livemode", then kept. An empty list says nothing. The account's name is
// the provider's, and its redacted id and mode never pass for one.
func TestAConnectionsModeIsWhatAPermittedCallShowed(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	call := func(op string) {
		t.Helper()
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: op, Args: args})
		e.s.call(ctx, g, "reply_ticket", args, permit)
	}
	grant := func() map[string]any {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("read the grant: %d %v", st, out)
		}
		return out
	}
	if v := grant(); v["mode"] != ModeUnknown || v["mode_words"] != "Mode not known yet: treated as live" {
		t.Fatalf("before any call: %v, %v", v["mode"], v["mode_words"])
	}
	executions := f.executions
	f.executeBody = `{"data":{"object":"list","data":[],"has_more":false},"successful":true}`
	call("op_empty")
	if v := grant(); v["mode"] != ModeUnknown || v["account_name"] != "Demo Shop" {
		t.Fatalf("after an empty list: mode %v, account %v", v["mode"], v["account_name"])
	}
	f.executeBody = `{"data":{"object":"list","data":[{"id":"in_1","object":"invoice","livemode":true}]},"successful":true}`
	call("op_live")
	if v := grant(); v["mode"] != "live" || v["mode_words"] != "Live" {
		t.Fatalf("after a live object: %v, %v", v["mode"], v["mode_words"])
	}
	// Kept: a later answer doesn't change it.
	f.executeBody = `{"data":{"id":"in_2","livemode":false},"successful":true}`
	call("op_test")
	if v := grant(); v["mode"] != "live" {
		t.Fatalf("a later answer changed the mode: %v", v["mode"])
	}
	// Only the three permitted calls reached the vendor: planque made none of its own.
	if f.executions != executions+3 {
		t.Fatalf("executions: %d, want %d", f.executions, executions+3)
	}
}

func TestModeOf(t *testing.T) {
	for raw, want := range map[string]string{
		`{"livemode":false}`: "test",
		`{"data":{"object":"list","data":[{"livemode":true}]}}`: "live",
		`{"response_data":{"livemode":false}}`:                  "test",
		`{"data":{"object":"list","data":[]}}`:                  "",
		`{"livemode":"REDACTED"}`:                               "",
		`not json`:                                              "",
	} {
		if got := modeOf(json.RawMessage(raw)); got != want {
			t.Errorf("modeOf(%s) = %q, want %q", raw, got, want)
		}
	}
}

// A connector call planque stops before the provider runs it is a refusal
// with its code, so KIFF never takes it for the vendor's error: nothing was
// sent.
func TestAConnectorCallStoppedBeforeTheProviderIsARefusalWithItsCode(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	f.accountStatus = "EXPIRED"
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_down"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_down", Args: args})
	executions := f.executions
	got := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if !strings.Contains(got, `"refused_by":"connection"`) || !strings.Contains(got, `"code":"`+CodeConnectionDown+`"`) || !strings.Contains(got, `"sent":false`) || !strings.Contains(got, "Not sent: ") {
		t.Fatalf("a call stopped before the provider: %s", got)
	}
	if f.executions != executions {
		t.Fatal("the provider ran it")
	}
}
