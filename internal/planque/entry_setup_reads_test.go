package planque

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/appcredentials"
)

// A first request for an entry that lists its tools may name setup reads: they
// pick which of its read-only tools the asking agent reads itself once the
// owner approved them. Coach asks for a Sheets sheet's header this way, and
// planque refused it ("setup_reads goes on a server's first request…"), so a
// normal first connect showed "1 failed". A name that isn't a read of the
// entry is still refused, naming the reads, and the shipped entry itself is
// never marked.
func TestAnEntryThatListsItsToolsTakesSetupReadsOnItsReads(t *testing.T) {
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
	path := "/v1/accounts/" + tenant + "/catalog/googlesheets@1/grants"
	ask := func(reads ...string) (int, map[string]any) {
		return e.do("POST", path, map[string]any{"agent": "agt_taf", "setup_reads": reads}, nil, "k_setup")
	}

	for _, name := range []string{"sheets_get_values", "google_sheets_read_range"} {
		st, out := ask(name)
		if st != 201 {
			t.Fatalf("setup read %q: %d %v", name, st, out)
		}
		g, err := e.s.loadGrant(ctx, tenant, out["grant"].(map[string]any)["id"].(string))
		if err != nil || len(g.Request.Tools) != 2 {
			t.Fatalf("the grant: %v %v", g, err)
		}
		if !g.Request.Tools[0].SetupReads || g.Request.Tools[1].SetupReads {
			t.Errorf("setup read marks: read %v, write %v", g.Request.Tools[0].SetupReads, g.Request.Tools[1].SetupReads)
		}
	}
	for _, t0 := range e.s.CatalogEntries["googlesheets@1"].Tools {
		if t0.SetupReads {
			t.Errorf("the shipped entry's %s was marked", t0.Tool)
		}
	}

	// The write, a name it doesn't have, and a repeat are each refused, naming the reads.
	for name, reads := range map[string][]string{"the write": {"sheets_append_values"}, "unknown": {"nope"}, "twice": {"sheets_get_values", "sheets_get_values"}} {
		st, out := ask(reads...)
		if st != 400 || !(strings.Contains(fmt.Sprint(out["error"]), "read-only tools") || strings.Contains(fmt.Sprint(out["error"]), "each tool once")) {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
}
