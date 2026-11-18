package planque

import (
	"os"
	"slices"
	"testing"

	"github.com/kiff/planque/internal/appconfig"
)

// A server without an app registered starts with the shipped catalog minus the
// entries that name it: planque#291 merged an entry for an app production
// hadn't registered yet, and the whole service refused to start (a 502, rolled
// back by the deploy's health check). The entry comes back once its app is
// registered, and nothing else about the catalog changes.
func TestAServerWithoutAnAppStartsWithoutItsEntries(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	whole, err := DecodeCatalog(raw, shippedApps(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeCatalog(raw, map[string]appconfig.Application{}); err == nil {
		t.Fatal("the shipped catalog decodes with no apps registered, so this test proves nothing")
	}

	kept, left := LeaveOutUnregistered(raw, map[string]appconfig.Application{})
	if !slices.Contains(left, "googlesheets@1") {
		t.Fatalf("left out: %v", left)
	}
	got, err := DecodeCatalog(kept, map[string]appconfig.Application{})
	if err != nil {
		t.Fatalf("the catalog without its app's entries doesn't decode: %v", err)
	}
	if _, ok := got["googlesheets@1"]; ok {
		t.Error("an entry whose app isn't registered was kept")
	}
	for id, e := range whole {
		if _, gone := got[id]; gone {
			continue
		}
		if e.Application == "" {
			t.Errorf("%s names no application and was left out", id)
		}
	}
	if len(got)+len(left) != len(whole) {
		t.Errorf("%d kept + %d left out, but the catalog has %d", len(got), len(left), len(whole))
	}

	// With its app registered, nothing is left out and the bytes are the same.
	kept, left = LeaveOutUnregistered(raw, shippedApps(t))
	if len(left) != 0 || string(kept) != string(raw) {
		t.Errorf("with the apps registered, left out %v", left)
	}
	// What doesn't read at all is DecodeCatalog's to refuse.
	if kept, left = LeaveOutUnregistered([]byte("not json"), nil); string(kept) != "not json" || left != nil {
		t.Errorf("an unreadable catalog was changed: %q %v", kept, left)
	}
}
