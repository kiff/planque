package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/profile"
)

// Qonto's sandbox wants a staging token on every request, the token endpoint included. The value is the
// owner's, kept with the client secret: it reaches Qonto's own hosts, only there, and shows nowhere else.
func TestQontoSandboxSendsItsStagingTokenToItsHostsAlone(t *testing.T) {
	const staging = "stg-token-ZZ9"
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	e.s.Applications = shippedApps(t)
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "qonto_sandbox", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret", HeaderValue: staging}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials
	raw, _ := os.ReadFile("../../deploy/catalog.json")
	var err error
	if e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	got := map[string][]string{} // host -> the staging header seen on each request
	var elsewhere []string       // requests that reached a host Qonto's definition doesn't name
	redirect := false
	transport := applicationTransport(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		got[r.URL.Host] = append(got[r.URL.Host], r.Header.Get("X-Qonto-Staging-Token"))
		body, status := `{}`, 200
		switch r.URL.Host {
		case "oauth-sandbox.staging.qonto.co":
			b, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(b), "client_secret=demo-secret") {
				t.Errorf("the token request lacks the client secret: %s", b)
			}
			body = `{"access_token":"a","refresh_token":"r","token_type":"Bearer","expires_in":3600}`
		case "thirdparty-sandbox.staging.qonto.co":
			body = `{"organization":{"slug":"demo"}}`
			if redirect {
				h := make(http.Header)
				h.Set("Location", "https://portal.example/login")
				return &http.Response{StatusCode: 302, Header: h, Body: io.NopCloser(strings.NewReader(""))}, nil
			}
		default:
			elsewhere = append(elsewhere, r.URL.Host)
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	base := e.client.Transport
	e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
			return base.RoundTrip(r)
		}
		var in struct {
			Profile json.RawMessage `json:"profile"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		hash, _ := profile.Digest(in.Profile)
		return profileResponse(200, `{"revision":1,"profile_sha256":"`+hash+`","words_sha256":"`+strings.Repeat("a", 64)+`"}`), nil
	})
	e.s.OAuth.HTTP = &http.Client{Transport: transport}
	e.s.REST.HTTP = &http.Client{Transport: transport}

	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-sandbox@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("grant: %d %v", st, out)
	}
	seen := mustJSONString(out)
	id := out["grant"].(map[string]any)["id"].(string)
	gr, _ := e.s.loadGrant(ctx, tenant, id)
	hash, _ := profile.Digest(gr.Request.Profile)
	path := "/v1/accounts/" + tenant + "/grants"
	st, approved := e.do("POST", path+"/"+id+"/approve", map[string]any{"digest": gr.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, approved)
	}
	seen += mustJSONString(approved)
	authURL, _ := url.Parse(approved["sign_in_url"].(string))
	if strings.Contains(authURL.String(), staging) {
		t.Fatalf("the staging token is in the sign-in URL: %s", authURL)
	}
	code := e.signIn(e.url + "/oauth/callback?state=" + url.QueryEscape(authURL.Query().Get("state")) + "&code=demo-code")
	if st, out = e.do("POST", path+"/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	seen += mustJSONString(out)

	e.kiff.Cards["k_taf"]["qonto_organization"] = 0
	res, err := e.kiff.Call(ctx, "k_taf", "qonto_organization", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("first read: %s", mustJSONString(res))
	if s := mustJSONString(res); strings.Contains(s, "Refused") || strings.Contains(s, "Not sent") {
		t.Fatalf("the read didn't go: %s", s)
	}
	// Qonto's own gate answers a call with a 302 to its portal: it is the answer, never followed.
	mu.Lock()
	redirect = true
	mu.Unlock()
	res, err = e.kiff.Call(ctx, "k_taf", "qonto_organization", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	seen += mustJSONString(res)

	mu.Lock()
	defer mu.Unlock()
	for _, host := range []string{"oauth-sandbox.staging.qonto.co", "thirdparty-sandbox.staging.qonto.co"} {
		if len(got[host]) == 0 {
			t.Fatalf("no request reached %s: %v", host, got)
		}
		for _, h := range got[host] {
			if h != staging {
				t.Errorf("a request to %s carried %q, want the staging token", host, h)
			}
		}
	}
	if len(elsewhere) != 0 || len(got["portal.example"]) != 0 {
		t.Errorf("a request left Qonto's hosts: %v %v", elsewhere, got["portal.example"])
	}
	// The value is in no answer the account saw, and in no stored grant.
	if strings.Contains(seen, staging) || strings.Contains(seen, "demo-secret") {
		t.Errorf("a secret shows in what the account was shown: %s", seen)
	}
	stored, _ := e.s.loadGrant(ctx, tenant, id)
	if b, _ := json.Marshal(stored); strings.Contains(string(b), staging) || strings.Contains(string(b), "demo-secret") {
		t.Errorf("a secret is stored in the grant: %s", b)
	}
	if b, _ := json.Marshal(stored.Request.Application.Definition); !strings.Contains(string(b), `"secret_header":"X-Qonto-Staging-Token"`) {
		t.Errorf("the approved definition lacks the header's name: %s", b)
	}
}

// The sandbox entry holds test data, so it is for the test accounts alone: another account neither sees it
// in the catalog nor in a search, and a grant for it is refused as an entry that doesn't exist.
func TestTheQontoSandboxIsOfferedToTestAccountsOnly(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	e.s.Applications = shippedApps(t)
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "qonto_sandbox", appcredentials.Credentials{ClientID: "c", Secret: "s", HeaderValue: "h"}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials
	raw, _ := os.ReadFile("../../deploy/catalog.json")
	var err error
	if e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	listed := func() bool {
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("catalog: %d %v", st, out)
		}
		return strings.Contains(mustJSONString(out), "qonto-sandbox@1")
	}
	found := func() bool {
		byApp, byTool := e.s.reviewedTools(tenant, "qonto", map[string]catalogEntryView{})
		return strings.Contains(mustJSONString(byApp), `"qonto-sandbox@1"`) || strings.Contains(mustJSONString(byTool), `"qonto-sandbox@1"`)
	}
	e.s.ManagedTestTenants = map[string]bool{}
	if listed() || found() {
		t.Fatal("an account that isn't a test account is offered the Qonto sandbox")
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-sandbox@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st == 201 {
		t.Fatalf("an account that isn't a test account asked for the Qonto sandbox: %d %v", st, out)
	}
	// The application is the one that is test-only: a direct request for it is refused the same way.
	direct := func() int {
		st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"application": "qonto_sandbox", "agent": "agt_taf",
			"tools": []map[string]any{{"tool": "qonto_get_organization", "name": "qonto_organization", "title": "Read the organisation and its bank accounts"}}}, nil, "k_setup")
		return st
	}
	if st := direct(); st != 400 {
		t.Fatalf("an account that isn't a test account asked for the Qonto sandbox application directly: %d", st)
	}
	e.s.ManagedTestTenants = map[string]bool{tenant: true}
	if st := direct(); st != 201 {
		t.Fatalf("a test account can't ask for the Qonto sandbox application directly: %d", st)
	}
	if !listed() || !found() {
		t.Fatal("a test account isn't offered the Qonto sandbox")
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-sandbox@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
		t.Fatalf("a test account can't ask for the Qonto sandbox: %d %v", st, out)
	}
}
