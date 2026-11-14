package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
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

// connections is the account's connections as tafs cloud reads them with a
// person's session.
func (e *env) connections(key string) (int, []map[string]any) {
	e.t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/connections", nil, nil, key)
	var cs []map[string]any
	if list, ok := out["connections"].([]any); ok {
		for _, c := range list {
			cs = append(cs, c.(map[string]any))
		}
	}
	return st, cs
}

// Any of the account's keys or sessions reads how its tafs reach their
// vendors, whoever asked for each connection, as it reads the catalog: each
// one's agent, vendor, kind and what asks for the same again, with no grant
// id, approval, text, tool or key. Another account's key reads nothing. Only
// a connection that holds a sign-in is listed: not one asked for and never
// signed in to, and not a disconnected one.
func TestTheAccountsConnectionsSayEachKind(t *testing.T) {
	e := newEnv(t)
	e.kiff.Keys["k_person"] = kiff.Me{TenantID: tenant}
	g := e.connect()
	v := newRESTVendor(t)
	e.connectAPI(v)
	for _, key := range []string{"k_person", "k_taf"} {
		st, cs := e.connections(key)
		if st != 200 || len(cs) != 2 {
			t.Fatalf("%s: %d %v", key, st, cs)
		}
		kinds := map[string]map[string]any{}
		for _, c := range cs {
			for _, k := range []string{"id", "text", "approval", "tools", "digest"} {
				if _, ok := c[k]; ok {
					t.Fatalf("%s: a connection carries %s: %v", key, k, c)
				}
			}
			kinds[c["kind"].(string)] = c
		}
		if c := kinds[KindVendorSignIn]; c == nil || c["agent"] != g.Request.Agent || c["server"] != e.vendor.MCPURL() || c["status"] != StatusConnected {
			t.Fatalf("%s: the vendor's own sign-in: %v", key, kinds)
		}
		if c := kinds[KindKey]; c == nil || c["api"].(map[string]any)["url"] != v.URL+"/openapi.json" || c["server"] != nil {
			t.Fatalf("%s: the key: %v", key, kinds)
		}
	}
	// A connection asked for and never signed in to isn't one the account has: pending with no sign-in,
	// then approved on the phone but not signed in at the vendor. Neither is listed.
	st, asked := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_new", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, asked)
	}
	if _, cs := e.connections("k_person"); len(cs) != 2 {
		t.Fatalf("a pending grant with no sign-in is listed: %v", cs)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+asked["id"].(string)+"/approve", map[string]any{"digest": asked["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	if _, cs := e.connections("k_person"); len(cs) != 2 {
		t.Fatalf("an approved grant with no sign-in is listed: %v", cs)
	}
	_, raw := e.do("GET", "/v1/accounts/"+tenant+"/connections", nil, nil, "k_person")
	if b, _ := json.Marshal(raw); strings.Contains(string(b), v.key) {
		t.Fatal("the connections carry the key")
	}
	if st, _ := e.connections("k_other"); st != 403 {
		t.Fatalf("another account's key: %d", st)
	}
	if st, _ := e.connections(""); st != 401 {
		t.Fatalf("no key: %d", st)
	}
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, ""); st != 200 {
		t.Fatalf("disconnect: %d", st)
	}
	if _, cs := e.connections("k_person"); len(cs) != 1 || cs[0]["kind"] != KindKey {
		t.Fatalf("after a disconnect: %v", cs)
	}
}

// Each connection names the taf it was asked for, by the request's
// agent_name, so tafs cloud can tell a removed taf's connection from another
// agent's whose id only looks like it.
func TestTheAccountsConnectionsNameTheirTaf(t *testing.T) {
	e, f := newConnectorFixture(t)
	e.kiff.Keys["k_person"] = kiff.Me{TenantID: tenant}
	e.tafName = "Weekly news digest"
	connectConnector(t, e, f)
	st, cs := e.connections("k_person")
	if st != 200 || len(cs) != 1 || cs[0]["agent_name"] != "Weekly news digest" {
		t.Fatalf("connections: %d %v", st, cs)
	}
}
