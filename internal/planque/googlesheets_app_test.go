package planque

import (
	"slices"
	"testing"
)

// The shipped Google Sheets app asks for the one scope its two tools need: a
// change to the scope or to the operations is a change to what an owner
// approves, so it fails here until it is reviewed.
func TestTheShippedGoogleSheetsAppAsksForOnlyWhatItsToolsNeed(t *testing.T) {
	apps := shippedApps(t)
	app, ok := apps["googlesheets"]
	if !ok {
		t.Fatal("deploy/apps.json has no googlesheets app")
	}
	if want := []string{"https://www.googleapis.com/auth/spreadsheets"}; !slices.Equal(app.OAuth.Scopes, want) {
		t.Fatalf("scopes %v, want %v", app.OAuth.Scopes, want)
	}
	reads, writes := map[string]bool{}, map[string]bool{}
	for _, op := range app.Snapshot.Operations {
		if op.ReadOnly {
			reads[op.Name] = true
		} else {
			writes[op.Name] = true
		}
	}
	if len(reads) != 1 || !reads["sheets_get_values"] || len(writes) != 1 || !writes["sheets_append_values"] {
		t.Fatalf("operations: reads %v, writes %v", reads, writes)
	}
}
