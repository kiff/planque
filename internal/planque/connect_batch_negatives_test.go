package planque

import (
	"strings"
	"testing"
)

// R6's guard (qa-kiff): an account's own address stays on the description's own server domain, https only, so the
// key never goes where a page or a request names. "https://{subdomain}.pay.example/v1": acme.pay.example, never
// pay.example.attacker.example, attacker.example or acme.pay.example.attacker.example.
func TestAnAccountsAddressStaysOnTheDescriptionsDomain(t *testing.T) {
	ask := func(base string) (int, string) {
		e := newEnv(t)
		v := newRESTVendor(t)
		e.enrol()
		v.mu.Lock()
		v.desc = strings.Replace(restDesc, `"info": {"title": "Pay API"},`, `"info": {"title": "Pay API"}, "servers": [{"url": "https://{subdomain}.pay.example/v1", "variables": {"subdomain": {"description": "yours"}}}],`, 1)
		v.mu.Unlock()
		api := map[string]any{"url": v.URL + "/openapi.json", "base": base, "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []any{}, "api": api}, nil, "k_setup")
		return st, mustJSONString(out)
	}
	for _, base := range []string{
		"https://attacker.example/v1",
		"https://pay.example.attacker.example/v1",
		"https://acme.pay.example.attacker.example/v1",
		"https://acme-pay.example/v1",
		"http://acme.pay.example/v1",
		"https://acme.pay.example:8443/v1",
		"https://a.b.pay.example/v1",
	} {
		if st, out := ask(base); st < 400 {
			t.Errorf("an account address off the description's domain was taken: %s (%d %s)", base, st, out)
		}
	}
}
