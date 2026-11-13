package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/words"
)

// The owner reads where a Basic key goes, with the username given beside it,
// in each language; the password given beside a key-as-username isn't said.
func TestTheApprovalSaysWhereABasicKeyGoes(t *testing.T) {
	q := Request{API: &APIRequest{URL: "https://vendor.example/openapi.yaml", Digest: strings.Repeat("a", 64), Base: "https://vendor.example/api"}}
	q.API.Auth = openapi.Auth{In: "basic", User: "me@example.com/token"}
	if got := apiKeyText(words.Approve("en"), q, false, ""); !strings.Contains(got, "as the password of an HTTP Basic sign-in, with the username me@example.com/token") {
		t.Fatalf("the key as the password: %s", got)
	}
	q.API.Auth = openapi.Auth{In: "basic", Password: "X"}
	if got := apiKeyText(words.Approve("en"), q, false, ""); !strings.Contains(got, "as the username of an HTTP Basic sign-in") {
		t.Fatalf("the key as the username: %s", got)
	}
	for _, lang := range []string{"de", "fr", "es", "it", "pt", "nl", "pl"} {
		w := words.Approve(lang)
		if w.KeyAsBasicUser == "" || !strings.Contains(w.KeyAsBasicPasswordF, "%s") || w.KeyAsBasicUser == words.Approve("en").KeyAsBasicUser {
			t.Errorf("%s: %q %q", lang, w.KeyAsBasicPasswordF, w.KeyAsBasicUser)
		}
	}
}

// The owner reads the fixed header and its value beside where the key goes.
func TestTheApprovalSaysTheFixedHeader(t *testing.T) {
	q := Request{API: &APIRequest{URL: "https://vendor.example/openapi.yaml", Digest: strings.Repeat("a", 64), Base: "https://acme.vendor.example/api"}}
	q.API.Auth = openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer ", Fixed: &openapi.Fixed{Name: "Harvest-Account-Id", Value: "123456"}}
	if got := apiKeyText(words.Approve("en"), q, false, ""); !strings.Contains(got, "in the Authorization header, with the Harvest-Account-Id header set to 123456 on each call") {
		t.Fatalf("%s", got)
	}
	for _, lang := range []string{"de", "fr", "es", "it", "pt", "nl", "pl"} {
		if w := words.Approve(lang); strings.Count(w.FixedHeaderF, "%s") != 2 || w.FixedHeaderF == words.Approve("en").FixedHeaderF {
			t.Errorf("%s: %q", lang, w.FixedHeaderF)
		}
	}
}

// A description at each account's own address, asked for with no base: the
// answer says the API is at the account's address, for Coach to ask the
// owner, rather than refusing the description.
func TestAnAccountsAddressIsAskedFor(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	v.mu.Lock()
	v.desc = strings.Replace(restDesc, `"info": {"title": "Pay API"},`, `"info": {"title": "Pay API"}, "servers": [{"url": "https://{subdomain}.pay.example/v1", "variables": {"subdomain": {"description": "yours"}}}],`, 1)
	v.mu.Unlock()
	api := map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []any{}, "api": api}, nil, "k_setup")
	if st != 400 || !strings.Contains(mustJSONString(out), "at each account's own address (https://{subdomain}.pay.example/v1)") {
		t.Fatalf("%d %v", st, out)
	}
}

// A base off the vendor's account server is refused; one on it goes on.
func TestAnAccountsBaseStaysOnTheVendorsServer(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	v.mu.Lock()
	v.desc = strings.Replace(restDesc, `"info": {"title": "Pay API"},`, `"info": {"title": "Pay API"}, "servers": [{"url": "https://{subdomain}.pay.example/v1", "variables": {"subdomain": {"description": "yours"}}}],`, 1)
	v.mu.Unlock()
	ask := func(base string) (int, string) {
		api := map[string]any{"url": v.URL + "/openapi.json", "base": base, "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []any{}, "api": api}, nil, "k_setup")
		return st, mustJSONString(out)
	}
	for _, bad := range []string{"https://evil.example/v1", "https://acme.pay.example.evil.example/v1", "https://a.b.pay.example/v1"} {
		if st, out := ask(bad); st != 400 || !strings.Contains(out, "an account's own address") {
			t.Errorf("%s: %d %s", bad, st, out)
		}
	}
	if _, out := ask("https://acme.pay.example/v1"); strings.Contains(out, "an account's own address") {
		t.Fatalf("the account's own address was refused: %s", out)
	}
}

// A key that's in the connection's fixed header is refused when the owner
// gives it: the fixed header is shown and kept with the request, never in
// the vault.
func TestAKeyInTheFixedHeaderIsRefusedWhenGiven(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	body := map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []map[string]any{{"tool": "listCharges", "title": "List charges"}},
		"api": map[string]any{"url": v.URL + "/openapi.json", "base": v.URL + "/v1",
			"auth": map[string]any{"in": "header", "name": "X-API-Key", "fixed": map[string]any{"name": "Account-Id", "value": "acct-" + v.key}}}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", body, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	st, out = e.giveKey(out["id"].(string), out["digest"].(string), v.key)
	if st != 400 || !strings.Contains(mustJSONString(out), "fixed header") || strings.Contains(mustJSONString(out), v.key) {
		t.Fatalf("give: %d %v", st, out)
	}
}
