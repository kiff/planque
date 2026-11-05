package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// asked is Setup's listing of every connection it asked for, by grant id.
func (e *env) asked() map[string]map[string]any {
	e.t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("Setup's list: %d %v", st, out)
	}
	byID := map[string]map[string]any{}
	for _, g := range out["grants"].([]any) {
		v := g.(map[string]any)
		byID[v["id"].(string)] = v
	}
	return byID
}

// Setup's listing says how each connection reaches its vendor, so a new taf
// is offered the kind the account already uses: the vendor's own sign-in
// (its server, no mode planque knows), a provider's app (its entry and its
// mode, as the phone's view says it) or a REST key (its API's public facts,
// never the key).
func TestSetupsListingSaysEachConnectionsKind(t *testing.T) {
	t.Run("the vendor's own sign-in", func(t *testing.T) {
		e := newEnv(t)
		id := e.connectedRefunds()
		g := e.asked()[id]
		if g["kind"] != KindVendorSignIn || g["server"] != e.vendor.MCPURL() {
			t.Fatalf("kind and server: %v", g)
		}
		for _, k := range []string{"entry", "api", "mode", "mode_words"} {
			if _, ok := g[k]; ok {
				t.Fatalf("a vendor's own sign-in has no %s: %v", k, g)
			}
		}
	})

	t.Run("a REST key", func(t *testing.T) {
		e := newEnv(t)
		v := newRESTVendor(t)
		g := e.connectAPI(v)
		_, raw := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_setup")
		b, _ := json.Marshal(raw)
		if strings.Contains(string(b), v.key) {
			t.Fatal("the listing carries the key")
		}
		got := e.asked()[g.ID]
		api, _ := got["api"].(map[string]any)
		auth, _ := api["auth"].(map[string]any)
		if got["kind"] != KindKey || api["url"] != v.URL+"/openapi.json" || api["base"] != v.URL+"/v1" || auth["in"] != "header" || auth["name"] != "X-API-Key" {
			t.Fatalf("kind and api: %v", got)
		}
		if _, ok := got["mode"]; ok {
			t.Fatalf("a key grant has no mode: %v", got)
		}
	})

	t.Run("a provider's app", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		got := e.asked()[g.ID]
		if got["kind"] != KindProvider || got["entry"] != "connector:composio:desk_demo" || got["mode"] != ModeUnknown || got["mode_words"] != "Mode not known yet: treated as live" {
			t.Fatalf("before any call: %v", got)
		}
		if _, ok := got["api"]; ok {
			t.Fatalf("a provider's app has no api: %v", got)
		}
		// Once a permitted call showed a test object, the listing says test, as the phone's view does.
		f.executeBody = `{"data":{"id":"in_1","object":"invoice","livemode":false},"successful":true}`
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_test"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_test", Args: args})
		e.s.call(context.Background(), g, "reply_ticket", args, permit)
		if got := e.asked()[g.ID]; got["mode"] != "test" || got["mode_words"] != "Test" {
			t.Fatalf("after a test object: %v", got)
		}
	})
}
