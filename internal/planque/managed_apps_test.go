package planque

import (
	"context"
	"net/http"
	"slices"
	"sync/atomic"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
)

// countingTransport counts the requests that reach the connection platform.
type countingTransport struct {
	next http.RoundTripper
	n    atomic.Int32
}

func (c *countingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.n.Add(1)
	return c.next.RoundTrip(r)
}

// gated is the connector fixture with the gate's lists set as asked (no
// test tenants unless given), and a count of what reaches the platform.
func gated(t *testing.T, apps, tenants []string) (*env, *countingTransport) {
	e, _, c := gatedFixture(t, apps, tenants)
	return e, c
}

func gatedFixture(t *testing.T, apps, tenants []string) (*env, *connectorFixture, *countingTransport) {
	t.Helper()
	e, f := newConnectorFixture(t)
	c := &countingTransport{next: f}
	e.s.Connectors = map[string]connector.Provider{"composio": &composio.Client{HTTP: &http.Client{Transport: c}, Secret: func(context.Context) (string, error) { return "synthetic-project-key", nil }}}
	e.s.ManagedApps, e.s.ManagedTestTenants = map[string]bool{}, map[string]bool{}
	for _, a := range apps {
		e.s.ManagedApps[provedKey(a)] = true
	}
	for _, id := range tenants {
		e.s.ManagedTestTenants[id] = true
	}
	return e, f, c
}

func askManaged(e *env) (int, map[string]any) {
	return e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf", "new_sign_in": true}, nil, "k_setup")
}

// A managed connection to an app the published list doesn't name is refused
// before anything reaches the platform: no app read, no set-up, no link.
func TestAManagedConnectionToAnUnlistedAppIsRefusedBeforeThePlatform(t *testing.T) {
	e, c := gated(t, []string{"Google Sheets"}, nil)
	st, out := askManaged(e)
	if st != http.StatusForbidden || out["error"] != errNotDisclosed {
		t.Fatalf("an unlisted app: %d %v", st, out)
	}
	if n := c.n.Load(); n != 0 {
		t.Fatalf("%d requests reached the platform", n)
	}
}

// Today's empty list refuses every managed connection, a held sign-in's
// reuse included.
func TestAnEmptyListRefusesEveryManagedConnection(t *testing.T) {
	e, c := gated(t, nil, nil)
	if st, _ := askManaged(e); st != http.StatusForbidden {
		t.Fatalf("an empty list let a connection through: %d", st)
	}
	st, _ := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]any{"agent": "agt_taf"}, nil, "k_setup")
	if st != http.StatusForbidden || c.n.Load() != 0 {
		t.Fatalf("a reuse: %d, %d requests reached the platform", st, c.n.Load())
	}
}

// A listed app connects, matched by its name as people know it.
func TestAListedAppConnects(t *testing.T) {
	e, c := gated(t, []string{"Desk Demo"}, nil)
	if st, out := askManaged(e); st != http.StatusCreated {
		t.Fatalf("a listed app: %d %v", st, out)
	}
	if c.n.Load() == 0 {
		t.Fatal("nothing reached the platform")
	}
}

// A test tenant tries any app; it widens only its own account.
func TestATestTenantTriesAnyAppAndOnlyItself(t *testing.T) {
	e, _ := gated(t, nil, []string{tenant})
	if st, out := askManaged(e); st != http.StatusCreated {
		t.Fatalf("a test tenant: %d %v", st, out)
	}
	e2, c := gated(t, nil, []string{"tnt_other"})
	if st, _ := askManaged(e2); st != http.StatusForbidden || c.n.Load() != 0 {
		t.Fatalf("another account's exception let this one through: %d", st)
	}
}

// A held sign-in reused through POST /grants (reuse) is checked as the
// catalog's reuse is: an app no longer allowed gets no new connection.
func TestAReuseThroughGrantsIsGated(t *testing.T) {
	e, f, c := gatedFixture(t, nil, []string{tenant})
	g := connectConnector(t, e, f)
	e.s.ManagedTestTenants = map[string]bool{}
	before := c.n.Load()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"reuse": g.ID, "agent": "agt_other", "tools": []any{}}, nil, "k_setup")
	if st != http.StatusForbidden || out["error"] != errNotDisclosed {
		t.Fatalf("a reuse through grants: %d %v", st, out)
	}
	if n := c.n.Load() - before; n != 0 {
		t.Fatalf("%d requests reached the platform", n)
	}
}

// A grant asked while its app was allowed, approved after the list no longer
// names it, starts no sign-in at the platform.
func TestApprovingAGrantAskedBeforeIsGated(t *testing.T) {
	e, c := gated(t, nil, []string{tenant})
	v := requestConnector(t, e)
	e.s.ManagedTestTenants = map[string]bool{}
	before := c.n.Load()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+v["id"].(string)+"/approve", map[string]any{"digest": v["digest"]}, e.phone, "")
	if st != http.StatusForbidden || out["error"] != errNotDisclosed {
		t.Fatalf("an approval after the app left the list: %d %v", st, out)
	}
	if n := c.n.Load() - before; n != 0 {
		t.Fatalf("%d requests reached the platform", n)
	}
}

// The test accounts come from the deploy's setting, in one form only.
func TestTestTenantsAreReadFromTheSetting(t *testing.T) {
	got, err := ManagedTestTenantsFrom(" tenants:tnt_a, tnt_b ")
	if err != nil || len(got) != 2 || !got["tnt_a"] || !got["tnt_b"] {
		t.Fatalf("tenants: %v %v", got, err)
	}
	if got, err := ManagedTestTenantsFrom(""); err != nil || len(got) != 0 {
		t.Fatalf("unset: %v %v", got, err)
	}
	for _, bad := range []string{"tnt_a", "tenants:", "tenants:tnt a", "tenants:tnt_a,,tnt_b", "tenants:../x", "all"} {
		if _, err := ManagedTestTenantsFrom(bad); err == nil {
			t.Errorf("%q read as test tenants", bad)
		}
	}
}

// The shipped list reads, with no empty or repeated names.
func TestTheShippedListsRead(t *testing.T) {
	if _, err := DecodeNames(managedAppsJSON, provedKey); err != nil {
		t.Fatal("managed_apps.json:", err)
	}
	for _, bad := range []string{`{}`, `null`, `[""]`, `["Gmail","gmail"]`, `["Google Sheets","GoogleSheets"]`, `[] []`} {
		if _, err := DecodeNames([]byte(bad), provedKey); err == nil {
			t.Errorf("%s read as a list", bad)
		}
	}
}

// The shipped list is the three Google apps the subprocessors page names
// ("Used only for: Gmail, Google Calendar, Google Sheets"), nothing more: an
// app that isn't on that page can't be connected through the platform.
func TestTheShippedListNamesTheGoogleAppsAndNothingElse(t *testing.T) {
	if len(ManagedApps) != 3 {
		t.Fatalf("the shipped list has %d apps: %v", len(ManagedApps), ManagedApps)
	}
	for _, name := range []string{"Gmail", "Google Calendar", "Google Sheets", "google sheets"} {
		if !ManagedApps[provedKey(name)] {
			t.Errorf("%s isn't on the shipped list", name)
		}
	}
	for _, name := range []string{"Moneybird", "Pipedrive", "Stripe", "Google Drive"} {
		if ManagedApps[provedKey(name)] {
			t.Errorf("%s is on the shipped list", name)
		}
	}
}

// An account named in PLANQUE_OWN_CLIENT_TENANTS isn't offered the managed
// connector's toolkit for an app that has an entry through tafs's own client:
// not in the catalog, not in the tool search, and a managed request for it is
// refused before the platform. Every other account, and every other app,
// keeps the managed route; and it decides ahead of the test accounts.
func TestAnOwnClientAccountIsNotOfferedTheManagedToolkit(t *testing.T) {
	e, c := gated(t, []string{"desk_demo", "own_demo"}, []string{tenant})
	e.s.CatalogEntries = map[string]CatalogEntry{"desk@1": {Application: "desk_demo", Vendor: "Support desk", Title: "Support desk",
		Tools: []ToolRequest{{Tool: "desk_read", Name: "desk_read", Title: "Read a desk", ReadOnly: true}}}}
	offered := func() (catalog, search []string) {
		_, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog?provider=composio", nil, nil, "k_setup")
		for _, it := range out["entries"].([]any) {
			catalog = append(catalog, it.(map[string]any)["entry"].(string))
		}
		_, out = e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?provider=composio&q=desk", nil, nil, "k_setup")
		for _, it := range out["tools"].([]any) {
			search = append(search, it.(map[string]any)["entry"].(string))
		}
		return
	}
	has := func(list []string, id string) bool { return slices.Contains(list, id) }
	managed := "connector:composio:desk_demo"

	// Not listed: the managed toolkit is offered, and connects.
	catalog, search := offered()
	if !has(catalog, managed) || !has(search, managed) {
		t.Fatalf("an ordinary account lost the managed toolkit: %v %v", catalog, search)
	}
	if st, out := askManaged(e); st != 201 {
		t.Fatalf("an ordinary account's managed request: %d %v", st, out)
	}

	e.s.OwnClientTenants = map[string]bool{tenant: true}
	catalog, search = offered()
	if has(catalog, managed) || has(search, managed) {
		t.Errorf("the managed toolkit is still offered: %v %v", catalog, search)
	}
	if !has(catalog, "connector:composio:calendar_demo") {
		t.Errorf("an app with no own entry lost its managed toolkit: %v", catalog)
	}
	if !has(catalog, "desk@1") || !has(search, "desk@1") {
		t.Errorf("the own entry isn't offered: %v %v", catalog, search)
	}
	before := c.n.Load()
	if st, out := askManaged(e); st != http.StatusForbidden || out["error"] != errNotDisclosed {
		t.Errorf("a managed request for the app: %d %v", st, out)
	}
	if n := c.n.Load() - before; n != 0 {
		t.Errorf("%d requests reached the platform", n)
	}

	// Another account is untouched.
	if !e.s.managedAllowed("other-tenant", "desk_demo") || e.s.ownClientFor("other-tenant", "desk_demo") {
		t.Error("an account that isn't listed was changed")
	}
}
