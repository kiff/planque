package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const keyed = `{"openapi": "3.1.0", "info": {"title": "Shop"}, "servers": [{"url": "https://api.shop.example/v1"}],
 "security": [{"key": []}],
 "components": {"securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}},
 "paths": {"/orders": {"get": {"operationId": "listOrders", "responses": {"200": {"description": "ok"}}}}}}`

// The matrix says what planque would do with each description, in words a fix can start from.
func TestTheMatrixReadsEachDescriptionAsPlanqueDoes(t *testing.T) {
	if route, ops, why := keyRoute([]byte(keyed)); !strings.HasPrefix(route, "✅ header X-API-Key") || ops != "1" || why != "" {
		t.Errorf("a key in a header: %q %q %q", route, ops, why)
	}
	yamlTwin := "openapi: 3.1.0\ninfo:\n  title: Shop\nservers:\n- url: https://api.shop.example/v1\nsecurity:\n- key: []\ncomponents:\n  securitySchemes:\n    key:\n      type: apiKey\n      in: header\n      name: X-API-Key\npaths:\n  /orders:\n    get:\n      operationId: listOrders\n      responses:\n        200:\n          description: ok\n"
	if route, _, why := keyRoute([]byte(yamlTwin)); !strings.HasPrefix(route, "✅") {
		t.Errorf("the same description in YAML: %q %q", route, why)
	}
	oauth := strings.Replace(keyed, `{"type": "apiKey", "in": "header", "name": "X-API-Key"}`, `{"type": "oauth2", "flows": {}}`, 1)
	if route, _, why := keyRoute([]byte(oauth)); route != "❌" || !strings.Contains(why, "managed sign-in") {
		t.Errorf("OAuth only: %q %q", route, why)
	}
	basic := strings.Replace(keyed, `{"type": "apiKey", "in": "header", "name": "X-API-Key"}`, `{"type": "http", "scheme": "basic"}`, 1)
	if route, _, why := keyRoute([]byte(basic)); route != "❌" || !strings.Contains(why, "HTTP basic") {
		t.Errorf("HTTP Basic: %q %q", route, why)
	}
	for body, want := range map[string]string{
		`{"kind": "discovery#restDescription"}`: "Google Discovery",
		`{"swagger": "2.0"}`:                    "Swagger 2",
		"<html><body>docs</body></html>":        "HTML",
		keyed:                                   "OpenAPI 3.1.0, JSON",
		yamlTwin:                                "OpenAPI 3.1.0, YAML",
	} {
		if got := describe([]byte(body)); !strings.HasPrefix(got, want) {
			t.Errorf("describe %.30q: %q, want %q", body, got, want)
		}
	}
}

// Without the sign-in reader key the managed column says so, and never guesses.
func TestTheManagedColumnSaysWhenItWasntChecked(t *testing.T) {
	if got := managed(nil, "gmail"); !strings.Contains(got, "not checked") {
		t.Errorf("no provider: %q", got)
	}
	if got := managed(nil, ""); got != "-" {
		t.Errorf("no slug: %q", got)
	}
	out := table([]row{{App: "a", Route: "✅ header X"}, {App: "b", Route: "❌", Why: "x | y"}})
	if !strings.Contains(out, "1 of 2 apps") || strings.Contains(out, "x | y") {
		t.Errorf("table: %s", out)
	}
}

// The managed column reads planque's sign-in state route with the reader key: each state in words, a refused key
// and a route that's off said as unread, and the key never sent anywhere a redirect points.
func TestTheManagedColumnReadsPlanquesSignInState(t *testing.T) {
	var elsewhere string
	other := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { elsewhere = r.Header.Get("Authorization") }))
	defer other.Close()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer reader" {
			w.WriteHeader(401)
			return
		}
		slug := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/catalog/connector:composio:"), "/sign-in")
		switch slug {
		case "moved":
			http.Redirect(w, r, other.URL+"/steal", http.StatusFound)
		case "gone":
			w.WriteHeader(404)
		default:
			_, _ = io.WriteString(w, `{"entry":"connector:composio:`+slug+`","state":"`+slug+`"}`)
		}
	}))
	defer srv.Close()
	p := &signInReader{base: srv.URL, key: "reader", http: srv.Client()}
	p.http.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	for slug, want := range map[string]string{
		"ready": "✅ ready", "on_the_owners_pick": "✅ on the owner's pick", "needs_setup": "❌ needs", "not_at_provider": "❌ not at the provider",
		"something_new": "? couldn't be read", "gone": "? not answered", "moved": "? couldn't be read",
	} {
		if got := managed(p, slug); !strings.HasPrefix(got, want) {
			t.Errorf("%s: %q, want %q", slug, got, want)
		}
	}
	if elsewhere != "" {
		t.Errorf("the reader key followed a redirect: %q", elsewhere)
	}
	if got := managed(&signInReader{base: srv.URL, key: "wrong", http: srv.Client()}, "ready"); !strings.Contains(got, "refused") {
		t.Errorf("a wrong key: %q", got)
	}
	t.Setenv("CONNECT_MATRIX_PLANQUE_URL", "http://planque.example")
	t.Setenv("CONNECT_MATRIX_SIGNIN_KEY", "reader")
	if managedReader() != nil {
		t.Error("the reader key would go over http")
	}
}
