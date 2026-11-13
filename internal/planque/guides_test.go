package planque

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The generator reads every app the provider can sign someone in to, and
// each one's tools: its guide part names each tool's access, inputs and
// returns, and the app list puts the first apps first, in batches.
func TestEveryConnectableAppGetsItsGeneratedGuide(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.search, f.searchScopes = true, []string{"https://www.googleapis.com/auth/drive.readonly"}
	apps, err := CollectGuideApps(context.Background(), e.s.Connectors["composio"])
	if err != nil {
		t.Fatal(err)
	}
	var desk *GuideApp
	for i := range apps {
		if apps[i].Slug == "desk_demo" {
			desk = &apps[i]
		}
		if apps[i].Slug == "unset_demo" || apps[i].Slug == "own_demo" {
			t.Fatalf("an app that can't be signed in to now got a guide: %s", apps[i].Slug)
		}
	}
	if desk == nil || len(desk.Tools) != 2 {
		t.Fatalf("the desk app: %+v", apps)
	}
	md := GeneratedGuide(*desk, "2026-11-13")
	for _, want := range []string{"Checked: 2026-11-13", "`DESK_DEMO_REPLY`", "`data.ticket_id`, `data.url`", "| `DESK_DEMO_SEARCH` | Search tickets | write | *q* |", "`drive.readonly`", "2 tools: 0 read, 2 write"} {
		if !strings.Contains(md, want) {
			t.Fatalf("the guide lacks %q:\n%s", want, md)
		}
	}
	list := GuideBatches(append(apps, GuideApp{Slug: "zz_first", Name: "Zed"}), []string{"zz_first"}, 1, "2026-11-13")
	if !strings.Contains(list, "## Batch 1\n\n| App | Slug | Connects | Tools |\n|---|---|---|---|\n| Zed | `zz_first`") {
		t.Fatalf("the first app isn't first:\n%s", list)
	}
}

// A reviewed entry's guide part lists the tools installed, with the
// arguments planque maps for them.
func TestAReviewedEntryGetsItsGeneratedGuide(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(raw, nil)
	if err != nil {
		t.Fatal(err)
	}
	md := GeneratedGuide(EntryGuideApp("qonto-key@2", entries["qonto-key@2"]), "2026-11-13")
	if !strings.Contains(md, "By the owner's API key") || !strings.Contains(md, "| `qonto_draft_client_invoice` | Draft an invoice from an accepted quote | write | `client_id`") {
		t.Fatalf("the reviewed entry's guide:\n%s", md)
	}
}

// Every adapted guide is dated and short, as deploy/guides/README.md says.
func TestEveryAdaptedGuideIsDatedAndShort(t *testing.T) {
	files, err := filepath.Glob("../../deploy/guides/*.md")
	if err != nil || len(files) == 0 {
		t.Fatalf("no guides: %v", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "README.md") || strings.HasSuffix(f, "apps.md") {
			continue
		}
		raw, _ := os.ReadFile(f)
		lines := strings.SplitN(string(raw), "\n", 4)
		if len(lines) < 3 || !strings.HasPrefix(lines[2], "Checked: ") {
			t.Errorf("%s: the line under the title isn't its Checked date", f)
		}
		if n := len(strings.Fields(string(raw))); n >= 1200 {
			t.Errorf("%s: %d words, more than 1,200", f, n)
		}
	}
}
