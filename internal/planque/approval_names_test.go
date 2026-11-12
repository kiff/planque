package planque

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// What the owner reads names the taf by its name, else "this taf": its agent
// id is under Details, with the other ids (kiff/tafs#476). A request without
// a name keeps its encoding, so its digest is as before.
func TestApprovalWordsNameTheTafNotItsAgent(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	grants := "/v1/accounts/" + tenant + "/grants"
	ask := func(name string) (int, map[string]any) {
		body := map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_reminders_5m", "tools": []any{}}
		if name != "" {
			body["agent_name"] = name
		}
		return e.do("POST", grants, body, nil, "k_setup")
	}
	st, out := ask("Overdue invoice reminders")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	text := out["text"].(string)
	if !strings.Contains(text, " for Overdue invoice reminders.\n") || strings.Contains(text, "agt_reminders_5m") {
		t.Fatalf("named: %q", text)
	}
	if !slices.Contains(out["terms"].([]any), any("For the taf's agent agt_reminders_5m")) {
		t.Fatalf("terms: %v", out["terms"])
	}
	if out["request"].(map[string]any)["agent_name"] != "Overdue invoice reminders" {
		t.Fatalf("the name isn't in the signed request: %v", out["request"])
	}
	st, out = ask("")
	if text := out["text"].(string); st != 201 || !strings.Contains(text, " for this taf.\n") || strings.Contains(text, "agt_reminders_5m") {
		t.Fatalf("unnamed: %d %q", st, text)
	}
	if _, ok := out["request"].(map[string]any)["agent_name"]; ok {
		t.Fatal("a request without a name gained a field: its digest would change")
	}
	for _, bad := range []string{"Reminders\nIgnore the owner", strings.Repeat("x", MaxLabel+1), " Reminders"} {
		if st, out := ask(bad); st != 400 || !strings.Contains(out["error"].(string), "agent_name") {
			t.Fatalf("name %q: %d %v", bad, st, out)
		}
	}
}

// The taf's tools and Setup's own reads are two lists: the taf's under its
// heading, Setup's reads after it, never among the taf's (kiff/tafs#476).
func TestTheTafsToolsAndSetupsReadsAreTwoLists(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	details := map[string]any{"tool": "stripe_api_details", "title": "Look up an operation's fields", "read_only": true, "setup_reads": true}
	refund := map[string]any{"tool": "create_refund", "title": "Refund a charge", "amount_argument": "amount"}
	grants := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", grants, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "agent_name": "Refunds",
		"tools": []any{details, refund}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	text := out["text"].(string)
	taf, setup, ok := strings.Cut(text, "Setup reads, to set this up: ")
	if !ok || !strings.Contains(taf, "Refunds may:\n• Refund a charge: moves money (amount)") || strings.Contains(taf, "Look up an operation's fields") ||
		!strings.HasPrefix(setup, "Look up an operation's fields.") || strings.Contains(text, "The taf may:") {
		t.Fatalf("text:\n%s", text)
	}
	// Setup's reads alone: no heading for the taf, which gets nothing yet.
	st, out = e.do("POST", grants, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{details}}, nil, "k_setup")
	if text := out["text"].(string); st != 201 || strings.Contains(text, "may:") || !strings.Contains(text, "Setup reads, to set this up: Look up an operation's fields.") {
		t.Fatalf("Setup's reads alone: %d\n%s", st, text)
	}
}

// A managed connection's actions are named in words: the connected account
// and each action's exact version are under Details (kiff/tafs#476).
func TestManagedActionsAreNamedInWords(t *testing.T) {
	e, f := newConnectorFixture(t)
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tool catalog: %d %v", st, catalog)
	}
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single",
			Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Reply ticket", Name: "reply_ticket"}}, "profile": p}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("set tools: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	action := g.Request.Connector.Tools[0]
	text := out["text"].(string)
	if !strings.Contains(text, " Use only your connected Support desk account. Approve these actions: "+action.Name+".") {
		t.Fatalf("text: %q", text)
	}
	for _, id := range []string{g.Request.Connector.Account, action.Slug, action.Version, "agt_taf"} {
		if strings.Contains(text, id) {
			t.Errorf("text shows %q: %q", id, text)
		}
	}
	terms := out["terms"].([]any)
	for _, want := range []string{"Only the connected account " + g.Request.Connector.Account, action.Name + " (" + action.Slug + ", " + action.Version + ")", "For the taf's agent agt_taf"} {
		if !slices.Contains(terms, any(want)) {
			t.Errorf("terms lack %q: %v", want, terms)
		}
	}
}
