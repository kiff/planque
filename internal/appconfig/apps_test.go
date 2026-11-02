package appconfig

import (
	"encoding/json"
	"testing"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
)

func validApplication() Application {
	return Application{Name: "Workspace", Kind: REST, Server: "https://api.example/v1", Source: "https://docs.example/operations",
		AuthMethod: "client_secret_post", Permissions: []string{"Read tickets"},
		OAuth: oauth.Metadata{Issuer: "https://auth.example", Authorization: "https://auth.example/authorize", Token: "https://auth.example/token",
			Scopes: []string{"tickets:read"}, NoResourceParameter: true},
		Snapshot: &openapi.Snapshot{Operations: []openapi.Operation{{Name: "listTickets", Method: "GET", Path: "/tickets", ReadOnly: true,
			InputSchema: json.RawMessage(`{"type":"object"}`)}}}}
}

func TestApplicationConfigurationIsCredentialFreeAndValidated(t *testing.T) {
	app := validApplication()
	b, _ := json.Marshal(map[string]Application{"workspace": app})
	if _, err := Decode(b); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{`null`, `{"workspace":{"client_secret":"secret"}}`, `{} {}`, `{"../workspace":{}}`} {
		if _, err := Decode([]byte(raw)); err == nil {
			t.Fatalf("accepted invalid application configuration: %s", raw)
		}
	}
	for _, change := range []func(*Application){
		func(a *Application) { a.OAuth.Token = "http://auth.example/token" },
		func(a *Application) {
			a.OAuth.AuthorizationParams = map[string]string{"redirect_uri": "https://elsewhere.example"}
		},
		func(a *Application) { a.OAuth.NoPKCE = true; a.AuthMethod = "none" },
		func(a *Application) { a.Kind = MCP },
		func(a *Application) { a.Permissions = nil },
		func(a *Application) { a.Snapshot.Operations[0].ReadOnly = false },
	} {
		app := validApplication()
		change(&app)
		if err := app.Validate(); err == nil {
			t.Fatal("accepted an unsafe or inconsistent application definition")
		}
	}
}
