package planque

import (
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
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
	mu                            sync.Mutex
	t                             *testing.T
	user, callback, accountStatus string
	// accountDown makes the provider fail to answer a connected account's read (503).
	accountDown bool
	// accountName, when set, is the vendor account's name the provider shows (else "Demo Shop").
	accountName string
	// noName: the provider shows no name for the vendor account.
	noName                            bool
	authScheme, authConfig            string
	links                             int
	foreign, changed, lost, strayLink bool
	// distinct, when set, gives each link after the first an account of its
	// own (account_demo_2, …), as the provider does; revoked lists the
	// accounts revoked, in order.
	distinct bool
	revoked  []string
	// managed marks desk_demo's sign-in as the provider's own; ceiling, when
	// set, is its managed scope ceiling, always is what it always asks, and
	// setUpFails fails an app's set-up at the provider, its answer naming the project key.
	setUpFails bool
	// noRecommend fails the scope recommendation (gap 5 of kiff/tafs#476).
	managed     bool
	ceiling     []string
	always      []string
	noRecommend bool
	// scopedMade lists the scoped sign-ins made, as name: scopes.
	scopedMade []string
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
	// granted, when set, is the scope the vendor granted, as the provider
	// keeps it on the connected account; search adds a search tool to the
	// catalog, needing searchScopes, and replyScopes are the reply's.
	granted                   string
	search                    bool
	searchScopes, replyScopes []string
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
		if name, _ := in.Config["name"].(string); strings.HasPrefix(name, "planque-scoped-") {
			// A scoped sign-in: the managed kind, named, asking for the scopes given and nothing of ours.
			creds, _ := in.Config["credentials"].(map[string]any)
			if in.Config["type"] != "use_composio_managed_auth" || len(in.Config) != 3 || len(creds) != 1 || creds["scopes"] == nil {
				f.t.Fatalf("a scoped sign-in with something else in it: %v", in.Config)
			}
			b, _ := json.Marshal(creds["scopes"])
			f.scopedMade = append(f.scopedMade, name+": "+string(b))
			body = `{"toolkit":{"slug":"` + in.Toolkit.Slug + `"},"auth_config":{"id":"ac_scoped_` + strconv.Itoa(len(f.scopedMade)) + `","auth_scheme":"OAUTH2","is_composio_managed":true}}`
			return profileResponse(201, body), nil
		}
		if in.Config["type"] != "use_composio_managed_auth" || len(in.Config) != 1 {
			f.t.Fatalf("a set-up with something of ours in it: %v", in.Config)
		}
		if f.setUpFails {
			return profileResponse(400, `{"error":{"message":"Toolkit has no managed auth for key synthetic-project-key"}}`), nil
		}
		f.setUps = append(f.setUps, in.Toolkit.Slug)
		body = `{"toolkit":{"slug":"` + in.Toolkit.Slug + `"},"auth_config":{"id":"ac_` + in.Toolkit.Slug + `_new","auth_scheme":"OAUTH2","is_composio_managed":true}}`
		return profileResponse(201, body), nil
	case path == "/auth_configs":
		scheme := f.authScheme
		if scheme == "" {
			scheme = "OAUTH2"
		}
		managed := "false"
		if f.managed {
			managed = "true"
		}
		body = `{"items":[{"id":"ac_desk_demo","auth_scheme":"` + scheme + `","status":"ENABLED","is_composio_managed":` + managed + `,"toolkit":{"slug":"desk_demo"}},{"id":"ac_calendar_demo","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"calendar_demo"}},{"id":"ac_unset_demo","auth_scheme":"API_KEY","status":"DISABLED","toolkit":{"slug":"unset_demo"}}]}`
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
		account := "account_demo"
		if f.distinct && f.links > 1 {
			account += "_" + strconv.Itoa(f.links)
		}
		body = `{"connected_account_id":"` + account + `","redirect_url":"https://connect.composio.dev/link/demo"}`
		if f.strayLink {
			body = `{"connected_account_id":"account_demo","redirect_url":"https://elsewhere.example/link/demo"}`
		}
	case strings.HasPrefix(path, "/connected_accounts/account_demo") && r.Method == "DELETE":
		if r.URL.Query().Get("revoke_on_delete") != "true" {
			f.t.Error("no upstream revocation requested")
		}
		f.revocations++
		f.revoked = append(f.revoked, strings.TrimPrefix(path, "/connected_accounts/"))
		body = `{"success":true,"revoke_job_id":"job_demo"}`
	case strings.HasPrefix(path, "/connected_accounts/account_demo") && f.accountDown:
		return profileResponse(http.StatusServiceUnavailable, `{"error":{"message":"unavailable"}}`), nil
	case strings.HasPrefix(path, "/connected_accounts/account_demo"):
		user := f.user
		if f.foreign {
			user = "another_owner"
		}
		// As the provider shows a vendor account: its name readable, its id and mode redacted.
		data := map[string]any{"displayName": cmp.Or(f.accountName, "Demo Shop"), "livemode": "REDACTED", "stripe_user_id": "REDACTED"}
		if f.noName {
			delete(data, "displayName")
		}
		if f.granted != "" {
			data["scope"] = f.granted
		}
		b, _ := json.Marshal(map[string]any{"id": strings.TrimPrefix(path, "/connected_accounts/"), "user_id": user, "status": f.accountStatus, "toolkit": map[string]string{"slug": "desk_demo"},
			"data": data})
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
		kit := map[string]any{"slug": "desk_demo", "name": "Support desk"}
		if f.ceiling != nil {
			kit["composio_managed_auth"] = []any{map[string]any{"mode": "OAUTH2", "scopes": map[string]any{"available": f.ceiling}}}
			kit["auth_config_details"] = []any{map[string]any{"mode": "OAUTH2", "required_scopes": f.always}}
		}
		b, _ := json.Marshal(kit)
		body = string(b)
	case path == "/toolkits/desk_demo/scopes/recommended" && r.Method == "POST":
		var in struct {
			Tools     []string `json:"tools"`
			Scheme    string   `json:"auth_scheme"`
			Available []string `json:"available_scopes"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		if f.noRecommend {
			return profileResponse(400, `{"error":{"message":"scope_not_available"}}`), nil
		}
		if in.Scheme != "OAUTH2" || !slices.Equal(in.Available, f.ceiling) {
			f.t.Fatalf("a recommendation asked outside the managed ceiling: %+v", in)
		}
		// Each tool's least-privilege scope: a reply writes tickets, a list reads them.
		var least []string
		for _, tool := range in.Tools {
			scope := "tickets.read"
			if strings.HasSuffix(tool, "_REPLY") {
				scope = "tickets.write"
			}
			if !slices.Contains(least, scope) {
				least = append(least, scope)
			}
		}
		b, _ := json.Marshal(map[string]any{"auth_scheme": "OAUTH2", "scopes": map[string]any{"least_privilege": least, "fewest": least}})
		body = string(b)
	case path == "/tools" || path == "/tools/DESK_DEMO_REPLY":
		schema := `{"type":"object","properties":{"ticket":{"type":"string"},"text":{"type":"string"}},"required":["ticket","text"],"additionalProperties":false}`
		if f.changed {
			schema = `{"type":"object","properties":{"other":{"type":"number"}}}`
		}
		replyScopes, _ := json.Marshal(f.replyScopes)
		tool := `{"slug":"DESK_DEMO_REPLY","toolkit":{"slug":"desk_demo"},"version":"20261101_00","name":"Reply to a ticket","scopes":` + string(replyScopes) + `,"input_parameters":` + schema +
			`,"output_parameters":{"type":"object","properties":{"data":{"type":"object","properties":{"ticket_id":{"type":"string"},"url":{"type":"string"},"status":{"type":"string"}}}}}}`
		if path == "/tools" {
			items := tool
			if f.search {
				searchScopes, _ := json.Marshal(f.searchScopes)
				items += `,{"slug":"DESK_DEMO_SEARCH","toolkit":{"slug":"desk_demo"},"version":"20261101_00","name":"Search tickets","tags":["readOnlyHint"],"scopes":` + string(searchScopes) +
					`,"input_parameters":{"type":"object","properties":{"q":{"type":"string"}},"additionalProperties":false}}`
			}
			body = `{"items":[` + items + `]}`
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
		if len(f.sent) != 4 || !strings.HasPrefix(string(f.sent["connected_account_id"]), `"account_demo`) || string(f.sent["version"]) != `"20261101_00"` {
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

// requestConnector asks for a sign-in of its own, even when the account holds
// one to the app (new_sign_in): a test of a held sign-in reused asks without it.
func requestConnector(t *testing.T, e *env) map[string]any {
	t.Helper()
	ask := map[string]any{"agent": "agt_taf", "new_sign_in": true}
	if e.tafName != "" {
		ask["agent_name"] = e.tafName
	}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", ask, nil, "k_setup")
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
	st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Reply ticket", Name: "reply_ticket"}}, "profile": p}, nil, "k_setup")
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
	// be connected, or reaches the end, beside planque's own Stripe and Qonto.
	st, out := e.do("GET", path, nil, nil, "k_setup")
	if st != 200 || out["next_cursor"] != "" || out["provider"] != "composio" || len(out["entries"].([]any)) != 4 {
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

// A managed sign-in that never came back to the app (kiff/planque#87, on the
// phone): the approving device asks again, and planque asks the provider,
// never answering 409. Signed in there, the grant finishes as the callback
// would, and the device confirms with the code it got. Not signed in, a fresh
// link for the same grant, and the unfinished account is revoked. Only the
// device that approved may ask.
func TestAManagedSignInThatNeverCameBackFinishesOrStartsAgain(t *testing.T) {
	// The phone signs each request anew; its signature carries the second.
	anew := func() { time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second))) }
	approve := func(t *testing.T, e *env) (string, map[string]any) {
		t.Helper()
		v := requestConnector(t, e)
		path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string)
		if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 || out["sign_in_url"] != "https://connect.composio.dev/link/demo" {
			t.Fatalf("approve: %d %v", st, out)
		}
		return path, v
	}
	t.Run("signed in at the provider, its redirect lost", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		path, v := approve(t, e)
		// No callback: the owner signed in, and the browser came back to the app without its code.
		if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.device("dev_other"), ""); st != 403 {
			t.Fatalf("another device asked: %d %v", st, out)
		}
		anew()
		st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
		if st != 200 || out["status"] != StatusSignedIn || out["confirm"] == nil || out["sign_in_url"] != nil || f.links != 1 {
			t.Fatalf("asked again: %d %v, %d links", st, out, f.links)
		}
		if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": out["confirm"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusListed {
			t.Fatalf("confirm: %d %v", st, out)
		}
	})
	t.Run("signed in, the callback's code lost on the way to the app", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		path, v := approve(t, e)
		lost := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
		anew()
		st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
		if st != 200 || out["status"] != StatusSignedIn || out["confirm"] == nil || out["confirm"] == lost {
			t.Fatalf("asked again: %d %v", st, out)
		}
		// The lost code no longer confirms: the new one does.
		if st, _ := e.do("POST", path+"/confirm", map[string]string{"confirm": lost}, e.phone, ""); st != 403 {
			t.Fatalf("the lost code confirmed: %d", st)
		}
		if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": out["confirm"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusListed {
			t.Fatalf("confirm: %d %v", st, out)
		}
	})
	t.Run("not signed in at the provider", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		f.distinct, f.accountStatus = true, "INITIATED"
		path, v := approve(t, e)
		anew()
		st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
		if st != 200 || out["status"] != StatusApproved || out["sign_in_url"] != "https://connect.composio.dev/link/demo" || f.links != 2 {
			t.Fatalf("asked again: %d %v, %d links", st, out, f.links)
		}
		g, _ := e.s.loadGrant(context.Background(), tenant, v["id"].(string))
		if g.Connector.Account != "account_demo_2" || !slices.Equal(f.revoked, []string{"account_demo"}) {
			t.Fatalf("the fresh sign-in is %q, revoked %v", g.Connector.Account, f.revoked)
		}
		// The fresh link's sign-in finishes as any does.
		f.accountStatus = "ACTIVE"
		code := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
		if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusListed {
			t.Fatalf("confirm: %d %v", st, out)
		}
	})
}

// A connection asked for with the taf's tools signs in for only what they
// need (gap 5 of kiff/tafs#476): the provider's least-privilege scopes for
// them, within its own sign-in's ceiling, on a managed sign-in made once per
// app and scope set and pinned for the sign-in. The scopes are in what the
// owner signs, the words say it, and the sign-in serves only those tools.
func TestAConnectionAskedWithItsToolsSignsInForOnlyWhatTheyNeed(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.managed, f.ceiling, f.always = true, []string{"admin", "tickets.read", "tickets.write"}, []string{"offline_access"}
	ask := func(tools ...string) (map[string]any, *Grant) {
		t.Helper()
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "agent_name": "Weekly tickets", "tools": tools}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("request: %d %v", st, out)
		}
		v := out["grant"].(map[string]any)
		g, _ := e.s.loadGrant(context.Background(), tenant, v["id"].(string))
		return v, g
	}
	v, g := ask("DESK_DEMO_LIST_TICKETS")
	c := g.Request.Connector
	if !slices.Equal(c.AskedTools, []string{"DESK_DEMO_LIST_TICKETS"}) || !slices.Equal(c.Scopes, []string{"offline_access", "tickets.read"}) || c.UsualAccess || c.AuthConfig != "ac_scoped_1" {
		t.Fatalf("the request: %+v", c)
	}
	if !slices.Equal(f.scopedMade, []string{`planque-scoped-desk_demo-` + scopedSuffix([]string{"tickets.read"}) + `: ["tickets.read"]`}) {
		t.Fatalf("scoped sign-ins made: %v", f.scopedMade)
	}
	// No provider's tool name ("List tickets" from DESK_DEMO_LIST_TICKETS): the actions come next, by their titles.
	if words := v["text"].(string); !strings.Contains(words, " Support desk is asked only for what this taf's actions need. ") || strings.Contains(words, "List tickets") {
		t.Fatalf("words: %q", words)
	}
	// The scopes are what the owner signs.
	widened := g.Request
	conn := *widened.Connector
	conn.Scopes = []string{"admin"}
	widened.Connector = &conn
	if sum, _ := digest(widened); sum == g.Digest {
		t.Fatal("the scopes aren't in the signed request")
	}
	// The sign-in pins the scoped sign-in, never the app's usual one.
	path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string)
	if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 || f.authConfig != "ac_scoped_1" {
		t.Fatalf("sign-in: %d %v, pinned %q", st, out, f.authConfig)
	}
	code := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}
	// A tool it didn't ask for isn't served by this sign-in: a new connection is.
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket",
		Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Reply ticket", Name: "reply_ticket"}}, "profile": p}, nil, "k_setup"); st != 409 || !strings.Contains(out["error"].(string), "isn't among them") {
		t.Fatalf("a tool outside the sign-in: %d %v", st, out)
	}
	// The same tools again, even asked in lower case: the scoped sign-in made before is pinned again.
	if _, again := ask("desk_demo_list_tickets"); again.Request.Connector.AuthConfig != "ac_scoped_1" || len(f.scopedMade) != 1 ||
		!slices.Equal(again.Request.Connector.AskedTools, []string{"DESK_DEMO_LIST_TICKETS"}) {
		t.Fatalf("asked again: %+v, %d made", again.Request.Connector, len(f.scopedMade))
	}
	for _, bad := range [][]string{{"not a slug"}, {"DESK_DEMO_LIST_TICKETS", "DESK_DEMO_LIST_TICKETS"}} {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "tools": bad}, nil, "k_setup")
		if st != 400 || !strings.Contains(out["error"].(string), "tools") {
			t.Fatalf("tools %q: %d %v", bad, st, out)
		}
	}
}

// When the provider can't narrow the sign-in (no managed ceiling, its
// recommendation failing, an operator's own sign-in, or the account's
// set-ups for the day used), the connection asks for the app's usual access
// and says so: the Card still bounds what the taf calls. It's recorded.
func TestAConnectionThatCantBeNarrowedSaysSo(t *testing.T) {
	useTheDay := func(t *testing.T, e *env, n int) {
		t.Helper()
		for i := 1; i <= n; i++ {
			if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_"+strconv.Itoa(i)+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
				t.Fatalf("set-up %d: %d %v", i, st, out)
			}
		}
	}
	for _, tc := range []struct {
		name string
		set  func(*testing.T, *env, *connectorFixture)
	}{
		{"its recommendation failed", func(_ *testing.T, _ *env, f *connectorFixture) {
			f.managed, f.ceiling, f.noRecommend = true, []string{"tickets.write"}, true
		}},
		{"no managed ceiling", func(_ *testing.T, _ *env, f *connectorFixture) { f.managed = true }},
		{"an operator's own sign-in", func(_ *testing.T, _ *env, f *connectorFixture) { f.ceiling = []string{"tickets.write"} }},
		{"the day's set-ups used", func(t *testing.T, e *env, f *connectorFixture) {
			f.managed, f.ceiling = true, []string{"tickets.write"}
			useTheDay(t, e, onDemandADay)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, f := newConnectorFixture(t)
			tc.set(t, e, f)
			st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "agent_name": "Replies", "tools": []string{"DESK_DEMO_REPLY"}}, nil, "k_setup")
			if st != 201 {
				t.Fatalf("request: %d %v", st, out)
			}
			v := out["grant"].(map[string]any)
			g, _ := e.s.loadGrant(context.Background(), tenant, v["id"].(string))
			if c := g.Request.Connector; !c.UsualAccess || c.Scopes != nil || c.AuthConfig != "ac_desk_demo" || len(f.scopedMade) != 0 {
				t.Fatalf("the request: %+v, made %v", c, f.scopedMade)
			}
			// Why it asks for more, and that the taf gets only what's approved; no provider's tool name.
			if words := v["text"].(string); !strings.Contains(words, " Support desk will ask for more access than this taf uses, because its sign-in can't be narrowed to this taf's actions. Sign in, then approve its actions: the taf gets only those, nothing before.") ||
				strings.Contains(words, "Reply.") {
				t.Fatalf("words: %q", words)
			}
		})
	}
}

// A new scoped sign-in is a set-up at the provider, so it counts against the
// account's day as setting up an app does: after it, one set-up fewer is left.
func TestANewScopedSignInCountsAsASetUp(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.managed, f.ceiling = true, []string{"tickets.read", "tickets.write"}
	for i := 1; i < onDemandADay; i++ {
		if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_"+strconv.Itoa(i)+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
			t.Fatalf("set-up %d: %d %v", i, st, out)
		}
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "tools": []string{"DESK_DEMO_REPLY"}}, nil, "k_setup"); st != 201 || len(f.scopedMade) != 1 {
		t.Fatalf("the scoped sign-in: %d %v, %d made", st, out, len(f.scopedMade))
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_"+strconv.Itoa(onDemandADay)+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 429 {
		t.Fatalf("the day's last set-up after a scoped sign-in: %d %v", st, out)
	}
}

// scopedSuffix is the name's suffix Scope gives a scope set.
func scopedSuffix(scopes []string) string {
	sum := sha256.Sum256([]byte(strings.Join(scopes, " ")))
	return hex.EncodeToString(sum[:8])
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
	tools := map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Reply ticket", Name: "reply_ticket"}}, "profile": p}
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
// name, it's "this taf", and no card shows its agent id; a name that isn't one
// line is refused.
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
	// Without a name, "this taf": its agent id is in neither the words nor Details (kiff/tafs#476).
	st, out = e.do("POST", path, map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if g, _ := out["grant"].(map[string]any); st != 201 || !strings.Contains(g["text"].(string), " for this taf. ") || strings.Contains(g["text"].(string), "agt_taf") ||
		strings.Contains(mustJSONString(g["terms"]), "agt_taf") {
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
	if g["on_demand"] != true || g["custody"] != "provider" || !strings.HasPrefix(g["text"].(string), "Connect Not set up for Invoices. Your Not set up sign-in is kept by our sign-in partner; the taf never sees it.") {
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
	// With the vendor's status, in the words KIFF's gateway reads it from.
	f.executeBody = `{"data":{"http_error":"Forbidden","status_code":403},"successful":false,"error":"Forbidden"}`
	if got := call("op_vendor_status"); !strings.Contains(got, "refused it (it answered 403): Forbidden. No claim is made about whether an effect occurred.") {
		t.Fatalf("the vendor's refusal with its status: %s", got)
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

// A managed connection's mode is what permitted calls showed, never a call
// of planque's own: unknown (treated as live) until an answer carries
// "livemode". Live wins: a test answer never hides a live one, and once live
// it stays live. An empty list says nothing. The account's name is the
// provider's, and its redacted id and mode never pass for one.
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
	f.executeBody = `{"data":{"id":"in_1","object":"invoice","livemode":false},"successful":true}`
	call("op_test")
	if v := grant(); v["mode"] != "test" || v["mode_words"] != "Test" {
		t.Fatalf("after a test object: %v, %v", v["mode"], v["mode_words"])
	}
	// Live wins: a live object after a test one makes it live.
	f.executeBody = `{"data":{"object":"list","data":[{"id":"in_2","livemode":false},{"id":"in_3","livemode":true}]},"successful":true}`
	call("op_live")
	if v := grant(); v["mode"] != "live" || v["mode_words"] != "Live" {
		t.Fatalf("after a live object: %v, %v", v["mode"], v["mode_words"])
	}
	// And stays live: a later test answer changes nothing.
	f.executeBody = `{"data":{"id":"in_4","livemode":false},"successful":true}`
	call("op_test_again")
	if v := grant(); v["mode"] != "live" {
		t.Fatalf("a test answer after live changed the mode: %v", v["mode"])
	}
	// Only the four permitted calls reached the vendor: planque made none of its own.
	if f.executions != executions+4 {
		t.Fatalf("executions: %d, want %d", f.executions, executions+4)
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

// Two tafs' connections on one provider's connected account: disconnecting
// one ends only that one, and the account stays signed in at the provider
// for the other; disconnecting the last one revokes it.
func TestDisconnectingOneTafKeepsAProviderAccountAnotherUses(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	a := connectConnector(t, e, f)
	b := connectConnector(t, e, f)
	if b.Connector.Account != a.Connector.Account {
		// The provider gave each its own account: pin b on a's, as one
		// connected account reused by a second taf is.
		b.Connector.Account = a.Connector.Account
		if err := e.s.saveGrant(ctx, b); err != nil {
			t.Fatal(err)
		}
	}
	st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+a.ID, nil, e.phone, "")
	if st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect the first: %d %v", st, out)
	}
	if f.revocations != 0 {
		t.Fatalf("disconnecting one taf revoked the account another uses: %d revocations", f.revocations)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, b.ID); g.Status != StatusConnected {
		t.Fatalf("the other taf's connection: %s", g.Status)
	}
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+b.ID, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect the second: %d %v", st, out)
	}
	if f.revocations != 1 {
		t.Fatalf("the last connection on the account didn't revoke it: %d revocations", f.revocations)
	}
}

// A new taf reuses a managed app's sign-in another taf has (kiff/tafs#476,
// gap 2: each new taf asked for Agent Mail's key again): its connection is
// signed in and listed at once, on the same account, with no second sign-in
// or key; the owner approves its tools alone; and the account stays while
// either uses it.
func TestANewTafReusesAManagedAppsSignIn(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	first := connectConnector(t, e, f)
	authorizations, links := f.authorizations, f.links
	grants := "/v1/accounts/" + tenant + "/grants"
	for name, body := range map[string]map[string]any{
		"with its tools at once": {"reuse": first.ID, "agent": "agt_refunds", "tools": []any{map[string]any{"tool": "connector_x"}}},
		"with a server too":      {"reuse": first.ID, "agent": "agt_refunds", "tools": []any{}, "server": "https://mcp.example"},
	} {
		if st, out := e.do("POST", grants, body, nil, "k_setup"); st != 400 {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
	st, out := e.do("POST", grants, map[string]any{"reuse": first.ID, "agent": "agt_refunds", "agent_name": "Refund requests", "tools": []any{}}, nil, "k_setup")
	if st != 201 || out["status"] != StatusListed || out["signed_in"] != true ||
		!strings.HasPrefix(out["text"].(string), "Connect Support desk for Refund requests with the sign-in you already have there: no second sign-in. Our sign-in partner keeps it; the taf never sees it.") {
		t.Fatalf("reuse: %d %v", st, out)
	}
	id := out["id"].(string)
	g, _ := e.s.loadGrant(ctx, tenant, id)
	if g.Connector.Account != first.Connector.Account || g.Request.Reuse != first.ID || g.Request.Agent != "agt_refunds" || f.authorizations != authorizations || f.links != links {
		t.Fatalf("the reused sign-in: %+v %+v; sign-ins %d→%d", g.Connector, g.Request, authorizations, f.authorizations)
	}
	// Its tools, from the catalog, approved by the owner on their device: no sign-in.
	path := grants + "/" + id
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "refund_reply", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single",
			Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	st, out = e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Refund reply", Name: "refund_reply"}}, "profile": p}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending || !strings.Contains(out["text"].(string), "with the sign-in you already have there") ||
		!strings.Contains(out["text"].(string), "Approve these actions: ") {
		t.Fatalf("tools: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	ph, _ := profile.Digest(g.Request.Profile)
	st, out = e.do("POST", path+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || f.authorizations != authorizations || f.links != links {
		t.Fatalf("approve: %d %v; sign-ins %d→%d", st, out, authorizations, f.authorizations)
	}
	// The first taf leaves; the account stays for the second.
	if st, _ := e.do("DELETE", grants+"/"+first.ID, nil, e.phone, ""); st != 200 || f.revocations != 0 {
		t.Fatalf("disconnect the first: %d, %d revocations", st, f.revocations)
	}
	// A connection that isn't connected any more can't be reused.
	if st, out := e.do("POST", grants, map[string]any{"reuse": first.ID, "agent": "agt_third", "tools": []any{}}, nil, "k_setup"); st != 409 || !strings.Contains(out["error"].(string), "isn't connected now") {
		t.Fatalf("reuse a disconnected one: %d %v", st, out)
	}
}

// A person's browser comes back to the connector's callback: whatever went
// wrong, it reads a page in its language, never planque's JSON. A sign-in
// cancelled at the provider says nothing was granted, and the grant waits.
func TestConnectorCallbackAnswersThePersonWithAPage(t *testing.T) {
	e, f := newConnectorFixture(t)
	page := func(path, lang string) (int, string, string) {
		t.Helper()
		req, _ := http.NewRequest("GET", e.url+path, nil)
		if lang != "" {
			req.Header.Set("Accept-Language", lang)
		}
		client := e.client
		if client == nil {
			client = http.DefaultClient
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, resp.Header.Get("Content-Type"), html.UnescapeString(string(raw))
	}
	isPage := func(name string, st int, ct, body, want string) {
		t.Helper()
		if !strings.HasPrefix(ct, "text/html") || strings.Contains(body, `{"error"`) || !strings.Contains(body, want) {
			t.Fatalf("%s: %d %s %q, want a page saying %q", name, st, ct, body, want)
		}
	}
	st, ct, body := page("/connectors/callback", "fr")
	isPage("no state", st, ct, body, "Ce lien n'a pas son état de connexion.")
	st, ct, body = page("/connectors/callback?state=nope", "")
	isPage("unknown state", st, ct, body, "This sign-in link is unknown, expired or already used.")

	// Approved, then cancelled at the provider: no account is active for it.
	v := requestConnector(t, e)
	id := v["id"].(string)
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	f.accountStatus = "INITIATED"
	u, _ := url.Parse(f.callback)
	st, ct, body = page(u.RequestURI(), "de")
	isPage("cancelled", st, ct, body, "wurde nicht abgeschlossen. Es wurde nichts gewährt.")
	if st == 200 || st == 302 {
		t.Fatalf("a cancelled sign-in answered %d", st)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.SignedIn || g.Status != StatusApproved {
		t.Fatalf("a cancelled sign-in moved the grant: %s signed=%v", g.Status, g.SignedIn)
	}
	// The link worked once: the same callback again is a used link, as a page.
	st, ct, body = page(u.RequestURI(), "")
	isPage("used", st, ct, body, "This sign-in link is unknown, expired or already used.")
}

// The sign-in link comes with the app's name as people name it, so the app's
// card says "Sign in to <app>" and never the provider's host the link goes
// through.
func TestASignInLinkNamesTheApp(t *testing.T) {
	e, _ := newConnectorFixture(t)
	v := requestConnector(t, e)
	path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string)
	st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
	if st != 200 || out["sign_in_url"] == nil {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, v["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	link, _ := url.Parse(out["sign_in_url"].(string))
	name, _ := out["app_name"].(string)
	if name == "" || name != e.s.appName(g) || name != g.VendorName || strings.Contains(name, link.Host) {
		t.Fatalf("app_name %q: want the app's own name (%q), not the link's host %q", name, g.VendorName, link.Host)
	}
}

// A managed app's catalog names each tool's provider slug beside its name,
// and setting tools takes the slug as well as the name: the app's guide and
// a search give the slug, which the catalog lists under another name (14 Nov,
// Gmail: Coach looked up GMAIL_FETCH_EMAILS and the setup read as failed).
func TestAManagedToolIsTakenByItsSlugAsByItsName(t *testing.T) {
	e, f := newConnectorFixture(t)
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	first := catalog["tools"].([]any)[0].(map[string]any)
	name, slug := first["name"].(string), first["slug"].(string)
	if slug == "" || slug == name || !strings.HasPrefix(name, "connector_") {
		t.Fatalf("the catalog's tool: %v", first)
	}
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	for _, asked := range []string{slug, strings.ToLower(slug)} {
		st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: asked, Title: "Reply ticket", Name: "reply_ticket"}}, "profile": p}, nil, "k_setup")
		if st != 200 {
			t.Fatalf("set tools by %q: %d %v", asked, st, out)
		}
		g, _ := e.s.loadGrant(context.Background(), tenant, id)
		if len(g.Request.Tools) != 1 || g.Request.Tools[0].Tool != name || len(g.Request.Connector.Tools) != 1 || g.Request.Connector.Tools[0].Slug != slug {
			t.Fatalf("by %q, the request: %+v", asked, g.Request.Tools)
		}
	}
	// A slug the catalog doesn't hold is refused as before, in Coach's own words.
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: "DESK_DEMO_NOTHING", Title: "Nothing", Name: "nothing"}}, "profile": p}, nil, "k_setup"); st != 400 ||
		!strings.Contains(out["error"].(string), "read the app's tool catalog before selecting this action") {
		t.Fatalf("an unknown slug: %d %v", st, out)
	}
}
