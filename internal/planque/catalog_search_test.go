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
	if e.s.CatalogEntries, err = DecodeCatalog(conf, nil); err != nil {
		t.Fatal(err)
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
	// Both entries: @1's 4 reads, and @2's 7 reads and its one draft.
	if st != 200 || reviewed(tools) != 12 || reviewed(tools[:12]) != 12 {
		t.Fatalf("by the app's name: %d %v", st, out)
	}
	for _, x := range tools[:12] {
		draft := x["entry"] == "qonto-key@2" && x["slug"] == "create_client_invoice"
		if !strings.HasPrefix(x["app_title"].(string), "Qonto invoicing") || x["app_state"] != "available" || x["read_only"] != !draft || x["slug"] == "" || x["name"] == "" {
			t.Fatalf("a reviewed tool as found: %v", x)
		}
	}
	// The managed connector's Qonto is still left out.
	for _, x := range tools {
		if x["app"] == qontoToolkit {
			t.Fatalf("the managed Qonto was offered: %v", x)
		}
	}
	// By a tool's words alone: only the tools they name, and after the
	// provider's apps. A word as generic as "invoices" doesn't put a reviewed
	// entry ahead of the app the person may use.
	if st, tools, out := search("q=invoices"); st != 200 || reviewed(tools) != 3 || reviewed(tools[:1]) != 0 || reviewed(tools[len(tools)-3:]) != 3 {
		t.Fatalf("by the tools' words: %d %v", st, out)
	}
	// The app's name and a tool's words: the app's name puts it first.
	if st, tools, out := search("q=qonto%20invoices"); st != 200 || reviewed(tools) != 12 || reviewed(tools[:1]) != 1 {
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
	if st, tools, out := search("q=Qonto"); st != 200 || reviewed(tools) != 12 || out["partial"] != true {
		t.Fatalf("provider down, reviewed found: %d %v", st, out)
	}
	if st, _, out := search("q=refund"); st != 503 {
		t.Fatalf("provider down, nothing reviewed: %d %v", st, out)
	}
}
