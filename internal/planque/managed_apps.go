package planque

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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

// managed_test_tenants.json lists the accounts KIFF tests with (Evals, phone
// proofs) by tenant id: they may make a managed connection to any app, for
// trying one before it's published. It widens only this gate.
//
//go:embed managed_test_tenants.json
var managedTestTenantsJSON []byte

// ManagedApps is managed_apps.json by folded name (provedKey), and
// ManagedTestTenants managed_test_tenants.json. A file that doesn't read
// stops the build's tests, never a request.
var (
	ManagedApps        = mustDecodeNames(managedAppsJSON, "managed_apps.json", provedKey)
	ManagedTestTenants = mustDecodeNames(managedTestTenantsJSON, "managed_test_tenants.json", tenantKey)
)

var tenantName = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// tenantKey is a tenant id as it is, or "" (not a name) when it isn't one.
func tenantKey(id string) string {
	if !tenantName.MatchString(id) {
		return ""
	}
	return id
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
	apps, tenants := s.ManagedApps, s.ManagedTestTenants
	if apps == nil {
		apps = ManagedApps
	}
	if tenants == nil {
		tenants = ManagedTestTenants
	}
	return tenants[tenant] || apps[provedKey(toolkit)]
}

// errNotDisclosed is what a person reads when a managed connection to an
// app isn't offered yet.
const errNotDisclosed = "This app can't be connected through tafs yet."
