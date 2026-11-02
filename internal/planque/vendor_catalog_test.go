package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/profile"
)

func deskCatalogEntry() CatalogEntry {
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "base_path", Value: "/v1"}, Source: "catalog", CatalogEntry: "desk@1",
		Operations: []profile.Operation{{Tool: "reply", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single", Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	raw, _ := json.Marshal(p)
	return CatalogEntry{ID: "desk@1", Vendor: p.Vendor, Title: "Support replies", Application: "desk", Tools: []ToolRequest{{Tool: "replyToTicket", Name: "reply", Title: "Reply to a ticket"}}, Profile: raw}
}

func installDeskCatalog(t *testing.T, e *env) CatalogEntry {
	t.Helper()
	e.s.Applications = map[string]appconfig.Application{"desk": deskApplication()}
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(context.Background(), "desk", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials
	entry := deskCatalogEntry()
	raw, _ := json.Marshal(map[string]CatalogEntry{entry.ID: entry})
	var err error
	e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications)
	if err != nil {
		t.Fatal(err)
	}
	return entry
}

func TestCatalogGrantUsesReviewedOperationsAndAccountState(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	entry := installDeskCatalog(t, e)
	path := "/v1/accounts/" + tenant + "/catalog"
	if st, _ := e.do("GET", path, nil, nil, "k_other"); st != 403 {
		t.Fatalf("other account: %d", st)
	}
	for _, q := range []string{"?q=SUPPORT", "?q=replies"} {
		st, out := e.do("GET", path+q, nil, nil, "k_setup")
		if st != 200 || len(out["entries"].([]any)) != 1 || out["by_key"] != true {
			t.Fatalf("search: %d %v", st, out)
		}
	}
	if _, out := e.do("GET", path+"?q=unrelated", nil, nil, "k_setup"); len(out["entries"].([]any)) != 0 {
		t.Fatal("unrelated entry returned")
	}
	endpoint := path + "/" + entry.ID + "/grants"
	for _, body := range []any{map[string]any{"agent": "agt_taf", "profile": json.RawMessage(entry.Profile)}, map[string]any{"agent": "agt_taf", "tools": []any{}}, map[string]any{}} {
		if st, _ := e.do("POST", endpoint, body, nil, "k_setup"); st != 400 {
			t.Fatalf("catalog body accepted: %d %v", st, body)
		}
	}
	st, out := e.do("POST", endpoint, map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("catalog grant: %d %v", st, out)
	}
	view := out["grant"].(map[string]any)
	if view["status"] != StatusPending || view["sign_in_url"] != nil || out["needs_key"] != false {
		t.Fatalf("preapproval response: %v", out)
	}
	id := view["id"].(string)
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Approval != nil || g.SignedIn || len(g.Tools) != 0 || g.Request.CatalogEntry != entry.ID || g.Request.Agent != "agt_taf" {
		t.Fatalf("granted before approval: %+v", g)
	}
	p, err := profile.Parse(g.Request.Profile)
	if err != nil || p.Source != "catalog" || p.Operations[0].Tool != g.Request.Tools[0].Name || g.Request.Tools[0].ResourceArg != "ticket" {
		t.Fatalf("profile not bound: %+v %v", g.Request, err)
	}
	expected, _ := digest(g.Request)
	if expected != view["digest"] {
		t.Fatal("digest does not cover final request")
	}
	// The entry remains pristine when a second grant has to allocate another
	// name. Its profile must follow the allocated name before being signed.
	st, out = e.do("POST", endpoint, map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("second grant: %d %v", st, out)
	}
	second := out["grant"].(map[string]any)
	next, _ := e.s.loadGrant(context.Background(), tenant, second["id"].(string))
	pp, _ := profile.Parse(next.Request.Profile)
	if next.Request.Tools[0].Name == g.Request.Tools[0].Name || pp.Operations[0].Tool != next.Request.Tools[0].Name || string(e.s.CatalogEntries[entry.ID].Profile) != string(entry.Profile) {
		t.Fatal("entry or profile names were mutated")
	}
	// Account state cannot be supplied by another account's grant, even if
	// both accounts have the same public catalog entry installed.
	g.Status = StatusConnected
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	other := *g
	other.Tenant = "tnt_other"
	other.ID = "grt_elsewhere"
	other.ver = 0
	other.Status = StatusNeedsReconnect
	if err := e.s.saveGrant(context.Background(), &other); err != nil {
		t.Fatal(err)
	}
	_, out = e.do("GET", path, nil, nil, "k_setup")
	listed := out["entries"].([]any)[0].(map[string]any)
	if listed["state"] != "connected" || listed["grant"] != id || listed["how"] != "oauth" {
		t.Fatalf("state: %v", out)
	}
	g.Status = StatusNeedsReconnect
	_ = e.s.saveGrant(context.Background(), g)
	_, out = e.do("GET", path, nil, nil, "k_setup")
	if out["entries"].([]any)[0].(map[string]any)["state"] != "needs_reconnect" {
		t.Fatalf("reconnect: %v", out)
	}
	g.Status = StatusDisconnected
	_ = e.s.saveGrant(context.Background(), g)
	_, out = e.do("GET", path, nil, nil, "k_setup")
	listed = out["entries"].([]any)[0].(map[string]any)
	if listed["state"] != "available" || listed["grant"] != nil {
		t.Fatalf("disconnected: %v", out)
	}
}

func TestCatalogKeyGrantRejectsAChangedDescription(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	p := refundProfile()
	p.Source = "catalog"
	p.CatalogEntry = "pay@1"
	raw, _ := json.Marshal(p)
	description := strings.ReplaceAll(restDesc, `"amount": {"type": "integer"}`, `"amount": {"type": "string"}, "metadata": {"type":"object","additionalProperties":{"type":"string"}}`)
	sum := sha256.Sum256([]byte(description))
	entry := CatalogEntry{ID: "pay@1", Vendor: p.Vendor, Title: "Refund payments", Kind: KindAPIKey, Profile: raw, Tools: []ToolRequest{{Tool: "createRefund", Name: "refund", Title: "Refund a charge"}},
		API: &apiAsk{URL: "https://pay.example/openapi.json", Base: "https://pay.example/v1", Auth: openapi.Auth{In: "header", Name: "X-API-Key"}}, DescriptionSHA256: hex.EncodeToString(sum[:])}
	conf, _ := json.Marshal(map[string]CatalogEntry{entry.ID: entry})
	var err error
	e.s.CatalogEntries, err = DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) { return profileResponse(200, description), nil })}
	path := "/v1/accounts/" + tenant + "/catalog/" + entry.ID + "/grants"
	st, out := e.do("POST", path, map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 || out["needs_key"] != true {
		t.Fatalf("key grant: %d %v", st, out)
	}
	description += "\n"
	if st, out := e.do("POST", path, map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "changed") {
		t.Fatalf("drift: %d %v", st, out)
	}
	items, _ := e.kv.Query(context.Background(), tenantPK(tenant), "GRANT#")
	if len(items) != 1 {
		t.Fatal("a changed description created another grant")
	}
}

func TestCatalogConfigurationMustMatchItsReviewedProfile(t *testing.T) {
	for name, change := range map[string]func(*CatalogEntry){
		"revision": func(e *CatalogEntry) { e.ID = "desk@2" },
		"vendor":   func(e *CatalogEntry) { e.Vendor = "Another desk" },
		"endpoint": func(e *CatalogEntry) { e.Server = "https://another.example" },
		"tools":    func(e *CatalogEntry) { e.Tools = nil },
		"profile": func(e *CatalogEntry) {
			p, _ := profile.Parse(e.Profile)
			p.Source = "drafted"
			p.CatalogEntry = ""
			e.Profile, _ = json.Marshal(p)
		},
	} {
		t.Run(name, func(t *testing.T) {
			e := deskCatalogEntry()
			change(&e)
			raw, _ := json.Marshal(map[string]CatalogEntry{e.ID: e})
			if _, err := DecodeCatalog(raw, map[string]appconfig.Application{"desk": deskApplication()}); err == nil {
				t.Fatal("invalid entry accepted")
			}
		})
	}
	if _, err := DecodeCatalog([]byte(`{"desk@1":{"unknown":true}}`), nil); err == nil {
		t.Fatal("unknown config accepted")
	}
}
