package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

// widePlatform is a provider with hundreds of apps, answering a search with far
// more than anyone should read: only what the two reads use is implemented.
type widePlatform struct {
	connector.Provider
	apps, toolsPerApp int
}

func (w widePlatform) Search(context.Context, string, []string, string) (connector.Search, error) {
	out := connector.Search{Apps: map[string]connector.App{}}
	add := func(app string, n int) {
		out.Apps[app] = connector.App{Slug: app, Name: "App " + app, Setup: connector.SetupReady}
		for i := 0; i < n; i++ {
			out.Tools = append(out.Tools, connector.Tool{Slug: strings.ToUpper(fmt.Sprintf("%s_LIST_THING_%d", app, i)), Toolkit: app, Name: fmt.Sprintf("List thing %d", i),
				Description: strings.Repeat("a long description of what this tool does and returns ", 20), ReadOnly: true})
		}
	}
	add("bigapp", 25)
	for i := 0; i < w.apps; i++ {
		add(fmt.Sprintf("app%03d", i), w.toolsPerApp)
	}
	return out, nil
}

func (w widePlatform) Catalog(context.Context, string, string) (connector.Page[connector.App], error) {
	var page connector.Page[connector.App]
	for i := 0; i < w.apps; i++ {
		page.Items = append(page.Items, connector.App{Slug: fmt.Sprintf("app%03d", i), Name: fmt.Sprintf("App %03d", i), Setup: connector.SetupReady})
	}
	return page, nil
}

// With 200 apps (and one with 25 tools) answering a search, planque's answer
// stays the size it is for 14: at most 60 tools, 10 of any one app, each
// description cut, with "more" saying what was left out; and a catalog read
// of 400 apps lists 300. A search of a few apps is untouched.
func TestAWideCatalogKeepsEveryAnswerBounded(t *testing.T) {
	e, _ := newConnectorFixture(t)
	e.s.Connectors = map[string]connector.Provider{"composio": widePlatform{apps: 200, toolsPerApp: 6}}

	st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?provider=composio&q=things", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("search: %d %v", st, out)
	}
	raw, _ := json.Marshal(out)
	tools := out["tools"].([]any)
	if len(tools) != maxFoundTools {
		t.Fatalf("%d tools in the answer, want %d", len(tools), maxFoundTools)
	}
	perApp := map[string]int{}
	for _, it := range tools {
		m := it.(map[string]any)
		perApp[m["entry"].(string)]++
		if n := len([]rune(m["description"].(string))); n > maxFoundDescription+1 {
			t.Fatalf("a description of %d characters", n)
		}
	}
	for app, n := range perApp {
		if n > maxFoundPerApp {
			t.Fatalf("%s has %d tools", app, n)
		}
	}
	if perApp["connector:composio:bigapp"] != maxFoundPerApp {
		t.Errorf("the app with 25 tools got %d", perApp["connector:composio:bigapp"])
	}
	if more := int(out["more"].(float64)); more != 25+200*6-len(tools) {
		t.Errorf("more %d, want %d", more, 25+200*6-len(tools))
	}
	if len(raw) > 40_000 {
		t.Errorf("the answer is %d bytes", len(raw))
	}

	// A search of a few apps is what it was.
	small := widePlatform{apps: 2, toolsPerApp: 3}
	e.s.Connectors["composio"] = small
	_, out = e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?provider=composio&q=things", nil, nil, "k_setup")
	if n := len(out["tools"].([]any)); n != maxFoundPerApp+2*3 || int(out["more"].(float64)) != 25-maxFoundPerApp {
		t.Errorf("a small search: %d tools, more %v", n, out["more"])
	}

	// The catalog read lists at most 300 entries of the 400 apps and the
	// account's reviewed ones, and says how many it left out.
	e.s.Connectors["composio"] = widePlatform{apps: 400}
	_, out = e.do("GET", "/v1/accounts/"+tenant+"/catalog?provider=composio", nil, nil, "k_setup")
	entries := len(out["entries"].([]any))
	if entries != maxCatalogEntries || int(out["more"].(float64)) < 100 {
		t.Errorf("the catalog lists %d entries, more %v", entries, out["more"])
	}
}
