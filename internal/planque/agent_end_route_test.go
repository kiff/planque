package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/kiff"
)

// manager is a person who manages this account's tafs, with their session.
func (e *env) manager() {
	e.kiff.Keys["s_manager"] = kiff.Me{TenantID: tenant, Capabilities: []string{"read", "tools.connect", "cards.manage"}}
}

func (e *env) end(agent, key, query string) (int, map[string]any) {
	e.t.Helper()
	return e.do("POST", "/v1/accounts/"+tenant+"/agents/"+agent+"/end"+query, map[string]any{}, nil, key)
}

func (e *env) statusOf(id string) string {
	e.t.Helper()
	g, err := e.s.loadGrant(e.t.Context(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	return g.Status
}

// An agent with no connections is 200 with nothing ended, never 404: the
// caller reads 404 as "this planque has no such route" and goes on
// without it, so a real end must never look like one.
func TestEndingAnAgentWithNoConnectionsEndsNothing(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.manager()
	for _, q := range []string{"?preview=1", ""} {
		st, out := e.end("agt_nobody", "s_manager", q)
		if st != 200 {
			t.Fatalf("%q: %d %v", q, st, out)
		}
		if l, _ := out["ended"].([]any); len(l) != 0 {
			t.Fatalf("%q ended %v", q, out["ended"])
		}
		if l, ok := out["kept"].([]any); !ok || len(l) != 0 {
			t.Fatalf("%q kept %v", q, out["kept"])
		}
	}
}

// The answer names each connection by its vendor and the status it had,
// the same in the preview and the call; a phone's session of a person who
// manages the tafs may end them, since ending only takes away.
func TestEndingAnAgentSaysWhatItEnded(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.manager()
	e.kiff.Keys["s_phone"] = kiff.Me{TenantID: tenant, DeviceSession: true, Capabilities: []string{"read", "cards.manage"}}
	pending, _ := e.request()
	st, preview := e.end("agt_taf", "s_phone", "?preview=1")
	if st != 200 || preview["preview"] != true {
		t.Fatalf("preview: %d %v", st, preview)
	}
	ended, _ := preview["ended"].([]any)
	if len(ended) != 1 {
		t.Fatalf("preview ended %v", preview["ended"])
	}
	one, _ := ended[0].(map[string]any)
	if one["grant"] != pending || one["was"] != StatusPending || one["vendor"] == "" {
		t.Fatalf("preview's line: %v", one)
	}
	if s := e.statusOf(pending); s != StatusPending {
		t.Fatalf("after the preview: %s", s)
	}
	st, done := e.end("agt_taf", "s_phone", "")
	if st != 200 || done["preview"] != false {
		t.Fatalf("end: %d %v", st, done)
	}
	if s := e.statusOf(pending); s != StatusDisconnected {
		t.Fatalf("after the end: %s, want disconnected", s)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, pending)
	if err != nil || !strings.Contains(g.Events[len(g.Events)-1].What, "ended with its taf") {
		t.Fatalf("its record: %+v %v", g.Events, err)
	}
	// Ended once: a second call ends nothing more.
	if st, again := e.end("agt_taf", "s_phone", ""); st != 200 || len(again["ended"].([]any)) != 0 {
		t.Fatalf("again: %d %v", st, again)
	}
}

// A field a route doesn't take is named in the refusal, so the caller
// fixes it on its first retry.
func TestARefusedBodyNamesTheFieldItDoesntTake(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	e.manager()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/agents/agt_taf/end", map[string]any{"uses": []string{"x"}}, nil, "s_manager")
	if st != 400 || !strings.Contains(out["error"].(string), `unknown field "uses"`) {
		t.Fatalf("%d %v", st, out)
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "uses": "x",
		"tools": []map[string]any{{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if st != 400 || !strings.Contains(out["error"].(string), `unknown field "uses"`) {
		t.Fatalf("a grant request: %d %v", st, out)
	}
}

// A sign-in another taf's connection uses stays, and the answer names the
// app and that taf, as the confirmation before the fingerprint draws it.
func TestEndingAnAgentNamesWhoKeepsASharedSignIn(t *testing.T) {
	e := newEnv(t)
	e.manager()
	first := e.connect()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID,
		"agent": "agt_chargebacks", "agent_name": "Chargebacks",
		"tools": []map[string]any{{"tool": "list_charges", "title": "List charges b", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse: %d %v", st, out)
	}
	for _, q := range []string{"?preview=1", ""} {
		st, out := e.end("agt_taf", "s_manager", q)
		if st != 200 {
			t.Fatalf("%q: %d %v", q, st, out)
		}
		kept, _ := out["kept"].([]any)
		if len(kept) != 1 {
			t.Fatalf("%q kept %v", q, out["kept"])
		}
		k, _ := kept[0].(map[string]any)
		used, _ := k["used_by"].([]any)
		if k["vendor"] != e.s.appName(first) || k["vendor"] == "" || len(used) != 1 || used[0] != "Chargebacks" {
			t.Fatalf("%q kept %v", q, k)
		}
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("the sign-in Chargebacks uses was revoked")
	}
}

// The app is named as people name it: a reviewed entry's vendor, not the
// entry's title; else the vendor's own name; else a host.
func TestAnAppIsNamedAsPeopleNameIt(t *testing.T) {
	e := newEnv(t)
	e.s.CatalogEntries = map[string]CatalogEntry{"qonto-key@2": {Vendor: "Qonto", Title: "Qonto Business API, invoicing subset (reviewed for planque)"}}
	for want, g := range map[string]*Grant{
		"Qonto":           {Request: Request{CatalogEntry: "qonto-key@2", API: &APIRequest{Base: "https://thirdparty.qonto.com/v2"}}, VendorName: "Qonto Business API, invoicing subset (reviewed for planque)"},
		"Stripe":          {VendorName: "Stripe", Request: Request{Server: "https://mcp.stripe.com"}},
		"mcp.example.com": {Request: Request{Server: "https://mcp.example.com/mcp"}},
		"api.example.com": {Request: Request{Kind: KindAPIKey, API: &APIRequest{Base: "https://api.example.com/v1"}}},
	} {
		if got := e.s.appName(g); got != want {
			t.Errorf("%s: %q", want, got)
		}
	}
}
