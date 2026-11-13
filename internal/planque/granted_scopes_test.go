package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// A managed app's catalog offers only the tools its sign-in's granted
// scopes let run: a search the vendor would refuse (403, insufficient
// scopes) is listed apart, with why in words, and a Card can't grant it,
// since it can't be bound to the connection. Its scopes known neither way,
// a tool is offered as before, and the vendor decides at the call.
// Category: authority
// Attack: binding a tool the sign-in's scopes don't let run
// Refused: 409, with the reason in words; nothing is bound
func TestAManagedAppOffersOnlyWhatItsSignInLetsRun(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.search, f.granted = true, "https://www.googleapis.com/auth/spreadsheets openid"
	f.replyScopes = []string{"https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/drive"}
	f.searchScopes = []string{"https://www.googleapis.com/auth/drive.readonly", "https://www.googleapis.com/auth/drive.metadata.readonly"}
	id := listConnector(t, e, f)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if len(g.Connector.Granted) != 2 {
		t.Fatalf("the granted scopes weren't kept: %v", g.Connector.Granted)
	}
	st, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
	tools, unusable := catalog["tools"].([]any), catalog["unusable"].([]any)
	if len(tools) != 1 || len(unusable) != 1 {
		t.Fatalf("offered %v, apart %v", tools, unusable)
	}
	reply, search := tools[0].(map[string]any), unusable[0].(map[string]any)
	if !strings.Contains(strings.ToLower(reply["title"].(string)), "reply") {
		t.Fatalf("the offered tool: %v", reply)
	}
	// What it gives back, from the provider's description of its answer: what a taf keeps.
	if r, _ := reply["returns"].([]any); len(r) != 2 || r[0] != "data.ticket_id" || r[1] != "data.url" {
		t.Fatalf("the reply's returns: %v", reply["returns"])
	}
	if why, _ := search["why"].(string); !strings.Contains(why, "only partly signed in: Search tickets needs access its sign-in didn't give") {
		t.Fatalf("the search's why: %v", search)
	}
	// Bound anyway: refused, with the same words.
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "search_tickets", Title: "Search tickets", Kind: "read", Shape: "single"}}}
	st, out := e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: search["name"].(string), Title: "Search tickets", Name: "search_tickets"}}, "profile": p}, nil, "k_setup")
	if st == 200 || !strings.Contains(out["error"].(string), "can't be used on this connection") {
		t.Fatalf("binding the search: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); len(g.Request.Tools) != 0 {
		t.Fatalf("bound: %v", g.Request.Tools)
	}
}

// With no scopes granted known, the catalog is as before: every tool.
func TestAManagedAppWhoseScopesArentKnownOffersEveryTool(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.search = true
	f.searchScopes = []string{"https://www.googleapis.com/auth/drive.readonly"}
	id := listConnector(t, e, f)
	st, catalog := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id+"/catalog", nil, nil, "k_setup")
	if st != 200 || len(catalog["tools"].([]any)) != 2 || len(catalog["unusable"].([]any)) != 0 {
		t.Fatalf("catalog: %d %v", st, catalog)
	}
}

// A managed app's sign-in reused for another taf names its account, as the
// provider shows it, in the owner's text: an owner who holds two can then
// say it's the other one.
func TestAReusedSignInNamesItsAccount(t *testing.T) {
	e, f := newConnectorFixture(t)
	first := connectConnector(t, e, f)
	grants := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", grants, map[string]any{"reuse": first.ID, "agent": "agt_refunds", "agent_name": "Refund requests", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse: %d %v", st, out)
	}
	path := grants + "/" + out["id"].(string)
	_, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "refund_reply", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single",
			Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	st, out = e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Refund reply", Name: "refund_reply"}}, "profile": p}, nil, "k_setup")
	if st != 200 || !strings.HasPrefix(out["text"].(string), "You're already signed in, to the account Demo Shop: no second sign-in.") {
		t.Fatalf("tools: %d %v", st, out["text"])
	}
}
