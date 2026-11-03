package planque

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

// newInProcessEnv exercises the real HTTP handlers, signatures and permits
// through an injected transport. It needs no listening port or live service.
func newInProcessEnv(t *testing.T) *env {
	t.Helper()
	e := &env{t: t, kiff: fakes.NewKIFF(), kv: &store.Memory{}, app: &app{}, url: "https://planque.example"}
	e.kiff.Issuer = "https://kiff.example"
	e.kiff.Keys = map[string]kiff.Me{
		"k_connector": {APIKey: true, TenantID: tenant, Roles: []string{"connector"}, Capabilities: []string{"read", "tools.connect"}},
		"k_setup":     {APIKey: true, TenantID: tenant, AgentID: "agt_setup"},
		"k_taf":       {APIKey: true, TenantID: tenant, AgentID: "agt_taf"},
		"k_other":     {APIKey: true, TenantID: "tnt_other", AgentID: "agt_other"},
	}
	e.kiff.Cards["k_taf"] = map[string]int64{}
	var handler http.Handler
	e.client = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		switch r.URL.Host {
		case "planque.example":
			handler.ServeHTTP(w, r)
		case "kiff.example":
			e.kiff.Handler().ServeHTTP(w, r)
		case "app.example":
			e.app.mu.Lock()
			e.app.last = map[string]string{"confirm": r.URL.Query().Get("confirm")}
			e.app.mu.Unlock()
			_, _ = w.WriteString("app opened")
		default:
			t.Fatalf("unexpected in-process request host %s", r.URL.Host)
		}
		return w.Result(), nil
	})}
	e.kiff.HTTP = e.client
	e.s = &Server{BaseURL: e.url, AppReturn: "https://app.example/signed-in", KV: e.kv, Sealer: &seal.Local{Key: make([]byte, 32)},
		KIFF:    &kiff.Client{API: e.kiff.Issuer, Gateway: e.kiff.Issuer, HTTP: e.client},
		Permits: &permit.Verifier{Issuer: e.kiff.Issuer, Keys: &permit.JWKS{URL: e.kiff.Issuer + "/.well-known/kiff-permit-keys.json", HTTP: e.client}},
		OAuth:   &oauth.Client{}, MCP: &mcp.Client{}, REST: &rest.Client{}}
	handler = e.s.Handler()
	e.phone = e.device("dev_phone")
	return e
}
