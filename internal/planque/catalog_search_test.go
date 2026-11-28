package planque

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

// searchDown is a provider whose tool search is down.
type searchDown struct{ fakeCatalogProvider }

func (searchDown) Search(context.Context, string, []string, string) (connector.Search, error) {
	return connector.Search{}, errors.New("down")
}

// Coach finds apps only by the tool search, so a reviewed entry is found
// there too, on the first page: first when a word names its app, after the
// provider's apps when a word names only one of its tools. Qonto by key (the
// shipped catalog) comes first for "Qonto", after the rest for "invoices".
func TestTheToolSearchFindsReviewedEntries(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	e := newInProcessEnv(t)
	e.enrol()
	if e.s.CatalogEntries, err = decodeOffered(conf, shippedApps(t)); err != nil {
		t.Fatal(err)
	}
	// This is about Qonto's entries: another vendor's tools that share a word
	// ("invoices") would move its counts, whatever else the catalog holds.
	for id := range e.s.CatalogEntries {
		if !strings.HasPrefix(id, "qonto-key@") {
			delete(e.s.CatalogEntries, id)
		}
	}
	e.s.Connectors = map[string]connector.Provider{"fake": fakeCatalogProvider{}}
	search := func(query string) (int, []map[string]any, map[string]any) {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?"+query, nil, nil, "k_setup")
		var tools []map[string]any
		if list, ok := out["tools"].([]any); ok {
			for _, x := range list {
				tools = append(tools, x.(map[string]any))
			}
		}
		return st, tools, out
	}
	reviewed := func(tools []map[string]any) (n int) {
		for _, x := range tools {
			if e, _ := x["entry"].(string); strings.HasPrefix(e, "qonto-key@") {
				n++
			}
		}
		return n
	}

	st, tools, out := search("q=Qonto")
	// Qonto's own sign-in first (ownSignInsFound), then the newest key
	// revision only (newestRevisions): @2's 7 reads, @1's four among them,
	// and its one draft. @1 isn't offered beside it.
	if st != 200 || tools[0]["entry"] != QontoEntry || reviewed(tools) != 8 || reviewed(tools[1:9]) != 8 {
		t.Fatalf("by the app's name: %d %v", st, out)
	}
	for _, x := range tools[1:9] {
		draft := x["entry"] == "qonto-key@2" && x["slug"] == "create_client_invoice"
		if !strings.HasPrefix(x["app_title"].(string), "Qonto invoicing") || x["app_state"] != "available" || x["read_only"] != !draft || x["slug"] == "" || x["name"] == "" {
			t.Fatalf("a reviewed tool as found: %v", x)
		}
	}
	// The managed connector's Qonto is still left out.
	for _, x := range tools {
		if strings.HasPrefix(x["entry"].(string), "connector:") && x["app"] == qontoToolkit {
			t.Fatalf("the managed Qonto was offered: %v", x)
		}
	}
	// By a tool's words alone: only the tools they name, and after the
	// provider's apps. A word as generic as "invoices" doesn't put a reviewed
	// entry ahead of the app the person may use.
	if st, tools, out := search("q=invoices"); st != 200 || reviewed(tools) != 2 || reviewed(tools[:1]) != 0 || reviewed(tools[len(tools)-2:]) != 2 {
		t.Fatalf("by the tools' words: %d %v", st, out)
	}
	// The app's name and a tool's words: the app's name puts it first.
	if st, tools, out := search("q=qonto%20invoices"); st != 200 || reviewed(tools) != 8 || reviewed(tools[1:2]) != 1 {
		t.Fatalf("by the app's name and the tools' words: %d %v", st, out)
	}
	// Words that name nothing reviewed, short words, a next page, or apps
	// named: the provider's alone.
	for _, query := range []string{"q=refund", "q=to%20a", "q=Qonto&cursor=c2", "q=Qonto&apps=desk_demo"} {
		if st, tools, out := search(query); st != 200 || reviewed(tools) != 0 {
			t.Fatalf("%s: %d %v", query, st, out)
		}
	}

	// The provider's search down: the reviewed entries found are said, as partial.
	e.s.Connectors = map[string]connector.Provider{"fake": searchDown{}}
	if st, tools, out := search("q=Qonto"); st != 200 || reviewed(tools) != 8 || out["partial"] != true {
		t.Fatalf("provider down, reviewed found: %d %v", st, out)
	}
	if st, _, out := search("q=refund"); st != 503 {
		t.Fatalf("provider down, nothing reviewed: %d %v", st, out)
	}
}

// A managed app no recorded run has proved is held by tafs cloud for every
// change. Coach asks a new taf's rules before any grant exists, so the search
// says it too, by the same rule as the grant (proofOf): on 27 Nov Coach asked
// "how many deals a day" for Pipedrive, and the Card was refused twice.
func TestTheToolSearchSaysAManagedAppIsUnproved(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	e.s.Connectors = map[string]connector.Provider{"fake": fakeCatalogProvider{}}
	desk := func() map[string]any {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?q=list", nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("search: %d %v", st, out)
		}
		for _, x := range out["tools"].([]any) {
			if m := x.(map[string]any); m["app"] == "desk_demo" {
				return m
			}
		}
		t.Fatalf("desk_demo not found: %v", out)
		return nil
	}
	// No proof for it: unproved, as its grant will say.
	e.s.Proved = map[string]Proof{}
	if got := desk(); got["unproved"] != true {
		t.Fatalf("no proof: %v", got)
	}
	// Proved on the managed route: nothing is said.
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-20", How: "oauth", Route: "managed"}}
	if got := desk(); got["unproved"] != nil {
		t.Fatalf("proved: %v", got)
	}
	// A proof through tafs's own client doesn't reach the managed route.
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-20", How: "oauth", Route: "own"}}
	if got := desk(); got["unproved"] != true {
		t.Fatalf("proved by its own client only: %v", got)
	}
}

// planque's own sign-in entries (qonto, stripe) have no tools until the owner
// signs in, so the search finds them by the vendor's name, first, as the app
// itself: in a Qonto re-test (27 Nov), "Qonto bank account balance" found only
// Qonto's key entries, a search with apps ["qonto"] found nothing, and Coach
// asked the owner for a key, which Qonto's guide forbids.
func TestTheToolSearchFindsAVendorsOwnSignInFirst(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	e := newInProcessEnv(t)
	e.enrol()
	if e.s.CatalogEntries, err = decodeOffered(conf, shippedApps(t)); err != nil {
		t.Fatal(err)
	}
	e.s.Connectors = map[string]connector.Provider{"fake": fakeCatalogProvider{}}
	search := func(query string) []map[string]any {
		t.Helper()
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?"+query, nil, nil, "k_setup")
		if st != 200 {
			t.Fatalf("%s: %d %v", query, st, out)
		}
		var tools []map[string]any
		for _, x := range out["tools"].([]any) {
			tools = append(tools, x.(map[string]any))
		}
		return tools
	}
	own := func(tools []map[string]any, entry string) (at int) {
		for i, x := range tools {
			if x["entry"] == entry {
				if x["sign_in"] != true || x["slug"] != "" || x["app_state"] != "available" || !strings.Contains(x["description"].(string), `request {entry: "`+entry+`"}`) {
					t.Fatalf("an own sign-in as found: %v", x)
				}
				return i
			}
		}
		return -1
	}
	// The re-test's words: Qonto's own sign-in, first, ahead of its key entries.
	tools := search("q=Qonto%20bank%20account%20balance")
	if own(tools, QontoEntry) != 0 || tools[0]["app_title"] != "Qonto" {
		t.Fatalf("by Qonto's name: %v", tools)
	}
	keys := 0
	for _, x := range tools[1:] {
		if strings.HasPrefix(x["entry"].(string), "qonto-") {
			keys++
		}
	}
	if keys == 0 {
		t.Fatalf("the key entries are still found, after it: %v", tools)
	}
	// The name with punctuation or a possessive is still the name.
	for _, query := range []string{"q=Qonto%27s%20balance", "q=qonto%3F", "q=%22Qonto%22%20invoices", "q=Qonto%E2%80%99s%20IBAN"} {
		if tools := search(query); own(tools, QontoEntry) != 0 {
			t.Fatalf("%s: %v", query, tools)
		}
	}
	// Asked as an app, with other words: found too, alone of the reviewed.
	if tools := search("q=balance&apps=qonto"); own(tools, QontoEntry) != 0 {
		t.Fatalf("with apps [qonto]: %v", tools)
	}
	// Stripe's the same.
	if tools := search("q=stripe%20refunds"); own(tools, StripeEntry) != 0 || own(tools, QontoEntry) != -1 {
		t.Fatalf("by Stripe's name: %v", tools)
	}
	// A proof test account connects Stripe only by a test key: its own sign-in isn't offered there; Qonto's still is.
	e.s.ProofTestTenants = map[string]bool{tenant: true}
	if tools := search("q=stripe%20qonto"); own(tools, StripeEntry) != -1 || own(tools, QontoEntry) != 0 {
		t.Fatalf("a proof test account: %v", tools)
	}
	e.s.ProofTestTenants = nil
	// Words that name neither, a word too short to name anything, or a next page: none.
	// A word only part of a vendor's name names it not: French "ont" in "Combien de factures ont été payées ?", or "strip".
	for _, query := range []string{"q=invoices", "q=qo", "q=Qonto&cursor=c2", "q=Combien%20de%20factures%20ont%20%C3%A9t%C3%A9%20pay%C3%A9es%20%3F",
		"q=strip%20the%20rows", "q=str%20tri%20rip", "q=qont"} {
		if tools := search(query); own(tools, QontoEntry) != -1 || own(tools, StripeEntry) != -1 {
			t.Fatalf("%s: %v", query, tools)
		}
	}
}
