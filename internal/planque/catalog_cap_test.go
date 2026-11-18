package planque

import (
	"fmt"
	"testing"
)

// Past the cap, the apps the account has a connection to stay in the list
// wherever they sort: an owner never loses sight of what they connected to
// because the catalog grew. The rest keep the order given.
func TestTheCatalogCapKeepsTheAccountsOwnConnections(t *testing.T) {
	var entries []catalogEntryView
	for i := 0; i < 350; i++ {
		entries = append(entries, catalogEntryView{Entry: fmt.Sprintf("connector:composio:app%03d", i), State: "available"})
	}
	entries[349].State = "connected"
	entries[348].State = "needs_reconnect"
	kept, more := capCatalog(entries)
	if len(kept) != maxCatalogEntries || more != 50 {
		t.Fatalf("%d kept, %d more", len(kept), more)
	}
	has := func(id string) bool {
		for _, e := range kept {
			if e.Entry == id {
				return true
			}
		}
		return false
	}
	if !has("connector:composio:app349") || !has("connector:composio:app348") {
		t.Error("an app the account is connected to dropped out of the list")
	}
	if !has("connector:composio:app297") || has("connector:composio:app298") {
		t.Error("the rest didn't keep the order given")
	}
	if few, none := capCatalog(entries[:10]); len(few) != 10 || none != 0 {
		t.Errorf("a short list was cut: %d %d", len(few), none)
	}
}

// Exactly the cap is no cut; one more is cut by one.
func TestTheCatalogCapCutsOnlyPastTheLimit(t *testing.T) {
	for n, wantMore := range map[int]int{maxCatalogEntries: 0, maxCatalogEntries + 1: 1} {
		var entries []catalogEntryView
		for i := 0; i < n; i++ {
			entries = append(entries, catalogEntryView{Entry: fmt.Sprintf("connector:composio:app%03d", i), State: "available"})
		}
		kept, more := capCatalog(entries)
		if len(kept) != maxCatalogEntries || more != wantMore {
			t.Errorf("%d entries: %d kept, %d more, want %d kept, %d more", n, len(kept), more, maxCatalogEntries, wantMore)
		}
	}
}
