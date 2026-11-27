package planque

import (
	"encoding/json"
	"os"
	"testing"
)

// Moneybird and QuickBooks are parked until they're tested: no account but
// the test accounts (PLANQUE_MANAGED_TEST_TENANTS) is offered either, by
// Moneybird's key entry or by the managed sign-in, and the entry stays in the
// catalog rather than being deleted. The test accounts still reach both, which
// is how they get tested.
func TestMoneybirdAndQuickBooksAreParked(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries map[string]CatalogEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatal(err)
	}
	key, ok := entries["moneybird-key@1"]
	if !ok {
		t.Fatal("Moneybird's key entry was deleted, not parked")
	}
	if !key.TestOnly {
		t.Fatal("Moneybird's key entry is offered to every account")
	}

	e, _ := gated(t, nil, []string{"tnt_test"})
	s := e.s
	s.ManagedApps, s.ManagedExcept = nil, nil // the shipped file
	if s.entryOffered("tnt_any", key) {
		t.Error("Moneybird's key entry is offered to an account that isn't a test account")
	}
	if !s.entryOffered("tnt_test", key) {
		t.Error("a test account can't reach Moneybird's key entry to test it")
	}
	for _, toolkit := range []string{"moneybird", "quickbooks"} {
		if s.managedListed("tnt_any", toolkit) {
			t.Errorf("%s's managed sign-in is offered to an account that isn't a test account", toolkit)
		}
		if !s.managedListed("tnt_test", toolkit) {
			t.Errorf("a test account can't reach %s's managed sign-in to test it", toolkit)
		}
	}
	// Every other app stays as it was.
	for _, toolkit := range []string{"googlesheets", "pipedrive", "notion"} {
		if !s.managedListed("tnt_any", toolkit) {
			t.Errorf("%s isn't allowed", toolkit)
		}
	}
}
