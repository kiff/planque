package planque

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/profile"
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

// The shipped catalog offers Google Sheets through our own registered client:
// a grant asked from googlesheets@1 records that client (never the managed
// connection provider's), the one scope its two tools need, a read and a held
// write, and nothing an owner didn't see. Without the entry, nobody could reach
// our client at all: the managed route is refused to ordinary accounts.
func TestAGrantFromTheShippedGoogleSheetsEntryRecordsOurOwnClient(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	e.s.Applications = shippedApps(t)
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "googlesheets", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	if e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	entry, ok := e.s.CatalogEntries["googlesheets@1"]
	if !ok || entry.Application != "googlesheets" || entry.Kind != "" || entry.API != nil {
		t.Fatalf("the shipped entry: %+v", entry)
	}

	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/googlesheets@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("grant: %d %v", st, out)
	}
	g, err := e.s.loadGrant(ctx, tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	r := g.Request
	if r.Connector != nil || r.Kind != appconfig.REST || r.CatalogEntry != "googlesheets@1" {
		t.Fatalf("not our own client's route: connector %v, kind %q, entry %q", r.Connector, r.Kind, r.CatalogEntry)
	}
	if r.Application == nil || r.Application.ID != "googlesheets" || r.Application.ClientID != "demo-client" || r.Application.Definition.Server != "https://sheets.googleapis.com" {
		t.Fatalf("the registered client wasn't recorded: %+v", r.Application)
	}
	if want := []string{"https://www.googleapis.com/auth/spreadsheets"}; !slices.Equal(r.Application.Definition.OAuth.Scopes, want) {
		t.Fatalf("scopes %v, want %v", r.Application.Definition.OAuth.Scopes, want)
	}
	if len(r.Tools) != 2 || r.Tools[0].Tool != "sheets_get_values" || r.Tools[0].Name != "google_sheets_read_range" ||
		r.Tools[1].Tool != "sheets_append_values" || r.Tools[1].Name != "google_sheets_add_row" ||
		!r.Tools[0].ReadOnly || r.Tools[1].ReadOnly {
		t.Fatalf("tools: %+v", r.Tools)
	}
	p, err := profile.Parse(r.Profile)
	if err != nil || p.Source != "catalog" || len(p.Operations) != 2 {
		t.Fatalf("profile: %+v %v", p, err)
	}
	reads, held := 0, 0
	for _, op := range p.Operations {
		switch {
		case op.Kind == "read":
			reads++
		case op.Kind == "write" && op.Effect == profile.EffectWrites:
			held++
		}
	}
	if reads != 1 || held != 1 {
		t.Fatalf("one read and one held write wanted, got %d and %d: %+v", reads, held, p.Operations)
	}
}
