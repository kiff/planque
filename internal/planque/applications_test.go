package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

func deskApplication() appconfig.Application {
	return appconfig.Application{Name: "Support desk", Kind: appconfig.REST, Server: "https://api.example/v1", Source: "https://docs.example/tickets",
		AuthMethod: "client_secret_post", Permissions: []string{"Reply to support tickets"},
		OAuth: oauth.Metadata{Issuer: "https://auth.example", Authorization: "https://auth.example/authorize", Token: "https://auth.example/token",
			Revocation: "https://auth.example/revoke", Scopes: []string{"tickets:reply"}, NoResourceParameter: true},
		Snapshot: &openapi.Snapshot{Operations: []openapi.Operation{{Name: "replyToTicket", Method: "POST", Path: "/tickets/{ticket}/replies", Body: "fields",
			Params:      []openapi.Param{{Name: "ticket", In: "path", Required: true}},
			InputSchema: json.RawMessage(`{"type":"object","properties":{"ticket":{"type":"string"},"text":{"type":"string"}},"required":["ticket","text"]}`)}}}}
}

func TestRegisteredApplicationApprovalKeepsItsOriginalDefinition(t *testing.T) {
	ctx := context.Background()
	kv := &store.Memory{}
	sealer := &seal.Local{Key: make([]byte, 32)}
	credentials := &appcredentials.Sealed{KV: kv, Sealer: sealer}
	if err := credentials.Put(ctx, "desk", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
		t.Fatal(err)
	}
	s := &Server{KV: kv, Sealer: sealer, Applications: map[string]appconfig.Application{"desk": deskApplication()}, ApplicationCredentials: credentials}
	approval, api, err := s.applicationRequest(ctx, "tenant", "grant", "desk", nil)
	if err != nil || api == nil || len(api.Digest) != 64 {
		t.Fatalf("registered application request failed: %v", err)
	}
	g := &Grant{Request: Request{Kind: appconfig.REST, Server: approval.Definition.Server, API: api, Application: approval}, Vendor: approval.Definition.OAuth, ClientID: approval.ClientID}
	before, _ := digest(g.Request)
	changed := s.Applications["desk"]
	changed.OAuth.Token = "https://another.example/token"
	changed.OAuth.Scopes[0] = "different-access"
	s.Applications["desk"] = changed
	registration, err := s.registration(ctx, g)
	if err != nil || registration[0].TokenURL != "https://auth.example/token" || approval.Definition.OAuth.Scopes[0] != "tickets:reply" {
		t.Fatal("a later app configuration reinterpreted an approved grant")
	}
	if err := credentials.Put(ctx, "desk", appcredentials.Credentials{ClientID: "demo-client", Secret: "rotated-secret"}); err != nil {
		t.Fatal(err)
	}
	registration, err = s.registration(ctx, g)
	after, _ := digest(g.Request)
	if err != nil || registration[0].Secret != "rotated-secret" || before != after {
		t.Fatal("secret rotation changed the approved request or failed to reach the transport")
	}
	b, _ := json.Marshal(g)
	if strings.Contains(string(b), "rotated-secret") || strings.Contains(string(b), "demo-secret") {
		t.Fatal("app credentials escaped into a grant record")
	}
	if err := credentials.Put(ctx, "desk", appcredentials.Credentials{ClientID: "replacement-client", Secret: "new-secret"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.registration(ctx, g); err == nil {
		t.Fatal("a changed client identity reused the old approval")
	}
}

type applicationTransport func(*http.Request) (*http.Response, error)

func (f applicationTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRegisteredRESTApplicationNeedsThePhoneThenCallsWithItsBearer(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "desk", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
		t.Fatal(err)
	}
	e.s.Applications = map[string]appconfig.Application{"desk": deskApplication()}
	e.s.ApplicationCredentials = credentials
	var mu sync.Mutex
	tokens, vendorCalls, revocations := 0, 0, 0
	transport := applicationTransport(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		body := `{}`
		switch r.URL.String() {
		case "https://auth.example/token":
			if err := r.ParseForm(); err != nil || r.PostForm.Get("client_id") != "demo-client" || r.PostForm.Get("client_secret") != "demo-secret" || r.PostForm.Get("resource") != "" {
				t.Error("token exchange did not use the configured registered client")
			}
			tokens++
			body = `{"access_token":"desk-access","refresh_token":"desk-refresh","token_type":"Bearer","expires_in":3600}`
		case "https://api.example/v1/tickets/ticket_1/replies":
			if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer desk-access" {
				t.Error("the operation was not sent with its approved method and bearer token")
			}
			vendorCalls++
			body = `{"id":"reply_1"}`
		case "https://auth.example/revoke":
			if err := r.ParseForm(); err != nil || r.PostForm.Get("client_secret") != "demo-secret" {
				t.Error("revocation did not use the registered client")
			}
			revocations++
		default:
			t.Errorf("unexpected registered-app request: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	e.s.OAuth.HTTP = &http.Client{Transport: transport}
	e.s.REST.HTTP = &http.Client{Transport: transport}
	counts := func() (int, int, int) {
		mu.Lock()
		defer mu.Unlock()
		return tokens, vendorCalls, revocations
	}
	path := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", path, map[string]any{"application": "desk", "agent": "agt_taf", "tools": []map[string]any{{"tool": "replyToTicket", "name": "reply_ticket", "title": "Reply to a ticket", "resource_argument": "ticket"}}}, nil, "k_setup")
	if st != 201 || strings.Contains(mustJSONString(out), "demo-secret") || !strings.Contains(out["text"].(string), "Reply to support tickets") {
		t.Fatalf("registered connection request: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	if st, _ := e.do("POST", path+"/"+id+"/approve", map[string]string{"digest": digest}, nil, ""); st != 401 {
		t.Fatal("registered app connected without the phone's signature")
	}
	if st, _ := e.do("POST", path+"/"+id+"/approve", map[string]string{"digest": strings.Repeat("0", 64)}, e.phone, ""); st != 409 {
		t.Fatalf("registered app accepted a different approval digest: %d", st)
	}
	if st, _ := e.do("GET", "/v1/accounts/tnt_other/grants/"+id, nil, nil, "k_other"); st != 404 {
		t.Fatalf("another account could read the registered grant: %d", st)
	}
	tokenCount, vendorCount, _ := counts()
	if tokenCount != 0 || vendorCount != 0 {
		t.Fatal("a request used the owner's account before approval")
	}
	st, approved := e.do("POST", path+"/"+id+"/approve", map[string]string{"digest": digest}, e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, approved)
	}
	authURL, _ := url.Parse(approved["sign_in_url"].(string))
	if authURL.Query().Get("client_id") != "demo-client" || authURL.Query().Get("resource") != "" || authURL.Query().Get("code_challenge") == "" {
		t.Fatal("sign-in did not preserve registered identity and PKCE")
	}
	code := e.signIn(e.url + "/oauth/callback?state=" + url.QueryEscape(authURL.Query().Get("state")) + "&code=demo-code")
	if _, ok := e.kiff.Connection(tenant, "reply_ticket"); ok {
		t.Fatal("tools were registered before the approving phone confirmed the sign-in")
	}
	st, out = e.do("POST", path+"/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	e.kiff.Cards["k_taf"]["reply_ticket"] = 0
	res, err := e.kiff.Call(ctx, "k_taf", "reply_ticket", map[string]any{"ticket": "ticket_1", "text": "We received your request."})
	_, vendorCount, _ = counts()
	if err != nil || res.Decision != "allowed" || vendorCount != 1 {
		t.Fatalf("registered REST call: %v %s, %d vendor calls", err, res.Decision, vendorCount)
	}
	// Expiry uses the existing refresh lease and the same registered client.
	g, _ := e.s.loadGrant(ctx, tenant, id)
	token, item, err := e.s.openToken(ctx, g)
	if err != nil {
		t.Fatal(err)
	}
	token.Expiry = time.Now().Add(-time.Minute)
	if err := e.s.storeToken(ctx, g, token, &item.Ver); err != nil {
		t.Fatal(err)
	}
	if _, _, err := e.s.accessToken(ctx, g, 0); err != nil {
		t.Fatal(err)
	}
	tokenCount, _, _ = counts()
	if tokenCount != 2 {
		t.Fatal("expired app token was not refreshed")
	}
	st, out = e.do("DELETE", path+"/"+id, nil, e.phone, "")
	_, _, revokeCount := counts()
	if st != 200 || out["status"] != StatusDisconnected || revokeCount != 2 {
		t.Fatalf("registered app disconnect: %d %v, revocations=%d", st, out, revokeCount)
	}
}

func TestRegisteredApplicationRequestCannotReplaceTheReviewedDefinition(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	e.s.Applications = map[string]appconfig.Application{"desk": deskApplication()}
	for _, fields := range []map[string]any{
		{"application": "unknown"},
		{"application": "desk", "server": "https://different.example"},
		{"application": "desk", "kind": appconfig.MCP},
		{"application": "desk", "api": map[string]any{"url": "https://different.example/api.json"}},
		{"kind": appconfig.REST},
	} {
		fields["agent"], fields["tools"] = "agt_taf", []any{}
		st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", fields, nil, "k_setup")
		if st != 400 {
			t.Fatalf("request could replace or bypass registered application facts: %d", st)
		}
	}
}
