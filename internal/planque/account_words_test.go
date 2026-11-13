package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// A managed app's connection is named by its account as the provider shows
// it, never by the provider's id for the connection (ca_…): in the grant's
// view, its listing for Setup and its catalog. With no name shown, it's left
// out, and the app names the row by the app.
func TestAManagedConnectionIsNamedByItsAccountNeverItsID(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	if g.Account == "" || !strings.HasPrefix(g.Account, "account_demo") {
		t.Fatalf("the fixture's provider id: %q", g.Account)
	}
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID
	for _, at := range []string{path, path + "/catalog"} {
		st, out := e.do("GET", at, nil, nil, "k_setup")
		if st != 200 || out["account"] != "Demo Shop" {
			t.Fatalf("%s: %d account %v", at, st, out["account"])
		}
	}
	// Its name for people, the same in every app: the app and its account, and the day it was asked for.
	_, view := e.do("GET", path, nil, nil, "k_setup")
	if view["title"] != "Support desk · Demo Shop" || len(view["since"].(string)) != 10 {
		t.Fatalf("title %v, since %v", view["title"], view["since"])
	}
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/asked-grants", nil, nil, "k_setup")
	if st != 200 || strings.Contains(stringOf(out), "account_demo") {
		t.Fatalf("Setup's listing names the provider's id: %d %v", st, out)
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
