package planque

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// The connection platform's sign-ins (the managed connections) reach a
// processor that KIFF's subprocessors page names, and that page says which
// apps it is used for. managed_apps.json is a copy of that list, the apps by
// the names people know them by ("Google Sheets"): planque starts a managed
// connection only for an app on it, so a connection the page doesn't
// disclose can't be made. An empty list refuses every one. An app is added
// to the published list first, then to this copy; the two are kept equal.
//
//go:embed managed_apps.json
var managedAppsJSON []byte

// ManagedApps is managed_apps.json by folded name (provedKey). A file that
// doesn't read stops the build's tests, never a request.
var ManagedApps = mustDecodeNames(managedAppsJSON, "managed_apps.json", provedKey)

var tenantName = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// ManagedTestTenantsFrom reads tenants:<id>[,<id>] (PLANQUE_MANAGED_TEST_TENANTS,
// set at deploy and never kept in this repository): the accounts KIFF tests
// with (Evals, phone proofs), which may make a managed connection to any app
// to try one before it's published. It widens only this gate. Anything else,
// or an id that isn't one, is an error: the setting is refused, never half read.
func ManagedTestTenantsFrom(v string) (map[string]bool, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return map[string]bool{}, nil
	}
	list, ok := strings.CutPrefix(v, "tenants:")
	if !ok {
		return nil, fmt.Errorf("must be tenants:<id>[,<id>]")
	}
	out := map[string]bool{}
	for _, t := range strings.Split(list, ",") {
		t = strings.TrimSpace(t)
		if !tenantName.MatchString(t) {
			return nil, fmt.Errorf("%q isn't a tenant id", t)
		}
		out[t] = true
	}
	return out, nil
}

// DecodeNames reads a JSON array of distinct, non-empty names, keyed by
// fold(name). Two names that fold the same are a repeat.
func DecodeNames(raw []byte, fold func(string) string) (map[string]bool, error) {
	var names []string
	dec := json.NewDecoder(bytes.NewReader(raw))
	if dec.Decode(&names) != nil || dec.Decode(new(any)) != io.EOF || names == nil {
		return nil, fmt.Errorf("must be a JSON array of names")
	}
	out := make(map[string]bool, len(names))
	for _, n := range names {
		k := fold(n)
		if k == "" || len(n) > 128 {
			return nil, fmt.Errorf("%q isn't a name", n)
		}
		if out[k] {
			return nil, fmt.Errorf("%q is listed twice", n)
		}
		out[k] = true
	}
	return out, nil
}

func mustDecodeNames(raw []byte, file string, fold func(string) string) map[string]bool {
	m, err := DecodeNames(raw, fold)
	if err != nil {
		panic(file + ": " + err.Error())
	}
	return m
}

// managedAllowed: tenant may start a managed connection to toolkit. The
// toolkit's slug is matched by its folded form to a folded published name
// ("googlesheets" to "Google Sheets"); an app whose slug isn't its name
// folded is refused until the two are made to match, never let through.
func (s *Server) managedAllowed(tenant, toolkit string) bool {
	if s.ownClientFor(tenant, toolkit) {
		return false
	}
	if s.ManagedTestTenants[tenant] {
		return true
	}
	// A mark, set by an operator's device, decides first: marked on, or marked
	// off (withdrawn), whatever the seed says. A vendor with no mark follows the seed.
	if on, marked := s.managedByMark(toolkit); marked {
		return on
	}
	apps := s.ManagedApps
	if apps == nil {
		apps = ManagedApps
	}
	return apps[provedKey(toolkit)]
}

// ownClientFor: tenant is one of PLANQUE_OWN_CLIENT_TENANTS (the same
// tenants:<id>[,<id>] form as the test accounts, set at deploy and never kept
// in this repository) and a reviewed entry reaches toolkit through an app's
// own client, so the managed connector's toolkit isn't offered to it: the
// account is shown, and can connect, only the entry that records tafs's own
// client (the way directEntryFor does for Stripe and Qonto, for every
// account). It decides before the test accounts' widening. Connections made
// before keep working: only a new managed connection, and the listings, ask.
func (s *Server) ownClientFor(tenant, toolkit string) bool {
	if !s.OwnClientTenants[tenant] {
		return false
	}
	key := provedKey(toolkit)
	for _, e := range newestRevisions(s.catalogEntries()) {
		if e.Application != "" && provedKey(e.Application) == key {
			return true
		}
	}
	return false
}

// managedClosed says g is a managed connection to an app that isn't on the
// list now (nor one of the test accounts' to try): the app was taken off it.
// A grant made while it was on is closed with it: its calls, the tools added
// to it and a sign-in resumed on it are refused, so taking an app off the list
// stops what was made, not only what is asked next.
func (s *Server) managedClosed(g *Grant) bool {
	return g.Request.Connector != nil && !s.managedAllowed(g.Tenant, g.Request.Connector.Toolkit)
}

// errAppClosed is what a person reads when an app was taken off the list and
// their connection to it is closed.
const errAppClosed = "This app isn't offered through tafs any more, so this connection is closed. Nothing was sent."

// errNotDisclosed is what a person reads when a managed connection to an
// app isn't offered yet.
const errNotDisclosed = "This app can't be connected through tafs yet."
