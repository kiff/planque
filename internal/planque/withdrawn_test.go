package planque

import (
	"os"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/appconfig"
)

// decodeOffered is the shipped catalog with every withdrawn entry offered again: for the tests of a
// key entry's own checks (its address, its tools, its IBAN read), which a grant made before the
// withdrawal still runs.
func decodeOffered(raw []byte, apps map[string]appconfig.Application) (map[string]CatalogEntry, error) {
	entries, err := DecodeCatalog(raw, apps)
	for id, e := range entries {
		if e.Withdrawn {
			e.Withdrawn = false
			entries[id] = e
		}
	}
	return entries, err
}

// Qonto connects only through its own sign-in, so its key entries are withdrawn: kept in the
// catalog, but offered to no account, the test accounts included. Not listed, not found by the tool
// search, and a new grant for either is refused as an entry that doesn't exist. A grant made
// before stays as it was.
func TestQontoByKeyIsWithdrawn(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	shipped, err := DecodeCatalog(raw, shippedApps(t))
	if err != nil {
		t.Fatal(err)
	}
	keys := []string{"qonto-key@1", "qonto-key@2"}
	for _, id := range keys {
		if e, ok := shipped[id]; !ok {
			t.Fatalf("%s was deleted, not withdrawn", id)
		} else if !e.Withdrawn {
			t.Fatalf("%s is still offered", id)
		}
	}

	// A grant made while the entry was offered.
	e := qontoKeyEnv(t)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@2/grants",
		map[string]any{"agent": "agt_taf", "tools": []string{"list_client_invoices"}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("a grant before the withdrawal: %d %v", st, out)
	}
	before, err := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil {
		t.Fatal(err)
	}

	e.s.CatalogEntries = shipped
	listed := func() bool {
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("catalog: %d %v", st, out)
		}
		return strings.Contains(mustJSONString(out), `"qonto-key@`)
	}
	found := func() bool {
		byApp, byTool := e.s.reviewedTools(tenant, "qonto invoices", map[string]catalogEntryView{})
		return strings.Contains(mustJSONString(byApp)+mustJSONString(byTool), `"qonto-key@`)
	}
	for _, testAccount := range []bool{false, true} {
		e.s.ManagedTestTenants = map[string]bool{tenant: testAccount}
		if listed() {
			t.Errorf("test account %v: Qonto's key entries are listed", testAccount)
		}
		if found() {
			t.Errorf("test account %v: the tool search finds Qonto's key entries", testAccount)
		}
		for _, id := range keys {
			if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/"+id+"/grants", map[string]any{"agent": "agt_other"}, nil, "k_setup"); st != 404 {
				t.Errorf("test account %v: a new grant for %s: %d %v", testAccount, id, st, out)
			}
		}
	}

	after, err := e.s.loadGrant(t.Context(), tenant, before.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != before.Status || mustJSONString(after.Request) != mustJSONString(before.Request) {
		t.Fatalf("the grant made before changed: %+v, was %+v", after, before)
	}
}
