package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// A managed app's connection is named for people by its app and its account
// as the provider shows it (`title`), never by the provider's id for the
// connection (ca_…). `account` stays that id: the apps send it back to name
// the connection a door or an event is for, and planque checks it there.
func TestAManagedConnectionIsNamedByItsAccountNeverItsID(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	if g.Account == "" || !strings.HasPrefix(g.Account, "account_demo") {
		t.Fatalf("the fixture's provider id: %q", g.Account)
	}
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID
	for _, at := range []string{path, path + "/catalog"} {
		st, out := e.do("GET", at, nil, nil, "k_setup")
		if st != 200 || out["account"] != g.Account {
			t.Fatalf("%s: %d account %v, want the provider's id kept", at, st, out["account"])
		}
	}
	// Its name for people, the same in every app: the app and its account, and the day it was asked for.
	_, view := e.do("GET", path, nil, nil, "k_setup")
	if view["title"] != "Support desk · Demo Shop" || len(view["since"].(string)) != 10 {
		t.Fatalf("title %v, since %v", view["title"], view["since"])
	}
	if strings.Contains(stringOf(view["title"]), "account_demo") {
		t.Fatalf("the title names the provider's id: %v", view["title"])
	}
	// The owner's approval text names the account by its name, never the id.
	if text, _ := view["text"].(string); strings.Contains(text, "account_demo") {
		t.Fatalf("the approval text names the provider's id: %s", text)
	}
}

// An unfinished managed connection, signed in with no tools approved (no
// taf behind it), is withdrawn after its day, as any held sign-in is: the
// sweep ends it, so the list never fills with leftovers.
func TestAnUnfinishedManagedConnectionEndsAfterADay(t *testing.T) {
	e, f := newConnectorFixture(t)
	id := listConnector(t, e, f)
	e.at(holdTTL + time.Minute)
	if n, err := e.s.Sweep(context.Background()); err != nil || n != 1 {
		t.Fatalf("sweep: %d %v", n, err)
	}
	if _, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup"); view == nil {
		t.Fatal("no view")
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusFailed || !strings.Contains(g.Reason, "deadline") && !strings.Contains(g.Reason, "within a day") {
		t.Fatalf("after a day: %s, %q", g.Status, g.Reason)
	}
}

func stringOf(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Coach reusing a connection for a new taf reads it from the account's
// connections, not the grant: so that view says unproved (and on_demand) as
// the grant view does, by the same rule. On 27 Nov Coach offered to reuse a
// Pipedrive connection and asked "how many deals a day" for its writes,
// which tafs cloud holds for every change of an unproved app.
func TestTheAccountsConnectionsSayUnprovedAsTheGrantDoes(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	listed := func() map[string]any {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/connections", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("connections: %d %v", st, out)
		}
		list, _ := out["connections"].([]any)
		if len(list) != 1 {
			t.Fatalf("want the one connection: %v", out)
		}
		return list[0].(map[string]any)
	}
	grant := func() map[string]any {
		t.Helper()
		_, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup")
		return view
	}
	slug := g.Request.Connector.Toolkit
	// No proof for the app: both say unproved.
	e.s.Proved = map[string]Proof{}
	if c, v := listed(), grant(); c["unproved"] != true || v["unproved"] != true {
		t.Fatalf("no proof: connection %v, grant %v", c["unproved"], v["unproved"])
	}
	// Proved on the managed route: neither says it.
	e.s.Proved = map[string]Proof{provedKey(slug): {At: "2026-11-20", How: "oauth", Route: "managed"}}
	if c, v := listed(), grant(); c["unproved"] != nil || v["unproved"] != nil {
		t.Fatalf("proved: connection %v, grant %v", c["unproved"], v["unproved"])
	}
	// The connection says on_demand exactly when its grant does.
	if c, v := listed(), grant(); c["on_demand"] != v["on_demand"] {
		t.Fatalf("on_demand: connection %v, grant %v", c["on_demand"], v["on_demand"])
	}
}
