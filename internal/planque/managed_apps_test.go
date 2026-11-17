package planque

import (
	"context"
	"net/http"
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
	return e, c
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

// The shipped lists read, with no empty or repeated names.
func TestTheShippedListsRead(t *testing.T) {
	if _, err := DecodeNames(managedAppsJSON, provedKey); err != nil {
		t.Fatal("managed_apps.json:", err)
	}
	if _, err := DecodeNames(managedTestTenantsJSON, tenantKey); err != nil {
		t.Fatal("managed_test_tenants.json:", err)
	}
	for _, bad := range []string{`{}`, `null`, `[""]`, `["Gmail","gmail"]`, `["Google Sheets","GoogleSheets"]`, `[] []`} {
		if _, err := DecodeNames([]byte(bad), provedKey); err == nil {
			t.Errorf("%s read as a list", bad)
		}
	}
	for _, bad := range []string{`["tnt 1"]`, `["../x"]`} {
		if _, err := DecodeNames([]byte(bad), tenantKey); err == nil {
			t.Errorf("%s read as tenants", bad)
		}
	}
}
