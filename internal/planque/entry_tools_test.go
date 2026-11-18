package planque

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// qontoKeyEnv serves the shipped catalog, with qonto-key@2's description.
func qontoKeyEnv(t *testing.T) *env {
	t.Helper()
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi-2.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	e := newInProcessEnv(t)
	e.enrol()
	e.s.CatalogEntries = entries
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		return profileResponse(200, string(description)), nil
	})}
	return e
}

// A reviewed entry is asked for with only the tools the taf will use, named
// as the tools search names them: the request, its profile and what KIFF
// holds before the key are those alone. Coach's uses reach planque as
// tools, so this is how a taf asks for least privilege.
func TestAReviewedEntryIsAskedWithOnlyTheToolsTheTafUses(t *testing.T) {
	e := qontoKeyEnv(t)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@2/grants",
		map[string]any{"agent": "agt_taf", "agent_name": "Overdue", "tools": []string{"list_client_invoices", "get_organization"}}, nil, "k_setup")
	if st != 201 || out["needs_key"] != true {
		t.Fatalf("ask: %d %v", st, out)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, tool := range g.Request.Tools {
		got = append(got, tool.Tool)
	}
	// In the entry's order, whatever the ask's.
	if strings.Join(got, ",") != "get_organization,list_client_invoices" {
		t.Fatalf("asked for: %v", got)
	}
	p, err := profile.Parse(g.Request.Profile)
	if err != nil || len(p.Operations) != 2 {
		t.Fatalf("the profile: %+v %v", p.Operations, err)
	}
	if c := e.kiff.Connections(tenant); len(c) != 2 {
		t.Fatalf("declared in KIFF: %+v", c)
	}
	draft := e.s.CatalogEntries["qonto-key@2"].Tools[7]
	if draft.Tool != "create_client_invoice" || strings.Contains(g.Text, draft.Title) || !strings.Contains(g.Text, g.Request.Tools[1].Title) {
		t.Fatalf("what the owner reads names a tool it doesn't ask for: %s", g.Text)
	}

	// Left out, every tool of the entry, as before.
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@2/grants", map[string]any{"agent": "agt_two"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("no tools: %d %v", st, out)
	}
	all, _ := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
	if len(all.Request.Tools) != len(e.s.CatalogEntries["qonto-key@2"].Tools) {
		t.Fatalf("no tools named: %d tools", len(all.Request.Tools))
	}
}

// A tool the entry doesn't have, one named twice, or none is refused, saying
// which tools it has; nothing is asked for.
func TestAReviewedEntryRefusesToolsItDoesntHave(t *testing.T) {
	e := qontoKeyEnv(t)
	for name, tools := range map[string][]string{
		"not the entry's":    {"list_client_invoices", "send_transfer"},
		"named twice":        {"list_clients", "list_clients"},
		"by its KIFF name":   {"qonto_list_clients"},
		"an empty selection": {},
	} {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@2/grants", map[string]any{"agent": "agt_taf", "tools": tools}, nil, "k_setup")
		if st != 400 {
			t.Fatalf("%s: %d %v", name, st, out)
		}
		if name != "an empty selection" && !strings.Contains(out["error"].(string), "list_client_invoices") {
			t.Fatalf("%s: the refusal doesn't say what the entry has: %v", name, out)
		}
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("declared after a refusal: %+v", c)
	}
}

// A tool whose profile reads another of the entry's tools to check its calls
// can't be asked for without it: the refusal names both.
func TestANarrowedEntryKeepsTheReadsItsChecksNeed(t *testing.T) {
	entry := &CatalogEntry{
		Tools: []ToolRequest{{Tool: "refund", Name: "x_refund"}, {Tool: "get_charge", Name: "x_get_charge", ReadOnly: true}},
		Profile: []byte(`{"v":"kiff.profile.v1","vendor":"X","transport":"rest","pin":{"kind":"none"},"source":"catalog","catalog_entry":"x@1","operations":[` +
			`{"tool":"x_refund","title":"Reply","kind":"write","effect":"sends","shape":"single","verify_read":{"tool":"x_get_charge","done_when":{"path":"status","equals":"sent"}}},` +
			`{"tool":"x_get_charge","title":"Read a charge","kind":"read","shape":"single"}]}`),
	}
	if _, _, err := entryTools(entry, []string{"refund"}); err == nil || !strings.Contains(err.Error(), "refund reads get_charge") {
		t.Fatalf("a write without the read it checks with: %v", err)
	}
	tools, prof, err := entryTools(entry, []string{"get_charge", "refund"})
	if err != nil || len(tools) != 2 || tools[0].Tool != "refund" {
		t.Fatalf("both: %+v %v", tools, err)
	}
	if p, err := profile.Parse(prof); err != nil || len(p.Operations) != 2 {
		t.Fatalf("both's profile: %v", err)
	}
}
