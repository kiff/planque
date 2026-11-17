package planque

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

const zendeskTemplate = "https://{account}.zendesk.com/api/v2"

// An account template is the vendor's own domain with the owner's one label
// in front: nothing else is a template.
func TestOnlyAnAccountOnTheVendorsOwnDomainIsATemplate(t *testing.T) {
	for _, ok := range []string{zendeskTemplate, "https://{account}.eu.example.co.uk/v1", "https://{account}.example.com"} {
		if !templateValid(ok) {
			t.Errorf("%s should be a template", ok)
		}
	}
	for _, bad := range []string{
		"https://{domain}.{tld}/api",                 // every part variable: any domain
		"https://{account}.com/api",                  // one fixed label: a whole TLD
		"https://{account}.zendesk.com:444/api",      // a port
		"https://{account}.{region}.zendesk.com",     // a second variable
		"https://api.zendesk.com/{account}/v2",       // the variable in the path
		"https://x{account}.zendesk.com/v2",          // not alone in its label
		"https://zendesk.com.{account}/v2",           // a variable after the domain
		"http://{account}.zendesk.com/v2",            // not https
		"https://{account}@zendesk.com/v2",           // user info
		"https://{account}.zendesk.com/v2?{account}", // a second variable
		"https://{acount}.zendesk.com/v2",            // not the account variable
	} {
		if templateValid(bad) {
			t.Errorf("%s must not be a template", bad)
		}
	}
}

// The owner's account is one label, checked here with the one rule planque
// already applies to an account's address, and the result stays on the
// vendor's own domain.
func TestAnAccountIsOneLabelOnTheVendorsDomain(t *testing.T) {
	a := &apiAsk{URL: "https://example.com/spec.json", Base: zendeskTemplate}
	got, why := resolveEntryAPI(a, "Acme", "")
	if why != "" || got.Base != "https://acme.zendesk.com/api/v2" || a.Base != zendeskTemplate {
		t.Fatalf("acme: %q %+v (the entry must stay a template: %s)", why, got, a.Base)
	}
	for _, bad := range []string{
		"acme.evil.com", "acme.zendesk.com.evil.com", "acme@evil", "acme:444", "", " ", "-acme", "acme-",
		strings.Repeat("a", 64), "аcme", // a Cyrillic а
		"acme/..", "acme%2f", "acme\n", "acme evil", "a.b", "[::1]",
	} {
		if got, why := resolveEntryAPI(a, bad, ""); why == "" {
			t.Errorf("account %q was accepted as %s", bad, got.Base)
		}
	}
	// 63 characters is the most a label can be.
	if _, why := resolveEntryAPI(a, strings.Repeat("a", 63), ""); why != "" {
		t.Errorf("63 characters: %s", why)
	}
	// An entry at an account's address takes the account, never a base; the other kinds take neither.
	if _, why := resolveEntryAPI(a, "acme", "https://evil.example"); why == "" {
		t.Error("a base was taken beside an account")
	}
	if _, why := resolveEntryAPI(&apiAsk{URL: a.URL, Base: "https://api.example.com/v1"}, "acme", ""); why == "" {
		t.Error("a fixed entry took an account")
	}
	if _, why := resolveEntryAPI(&apiAsk{URL: a.URL, Base: "https://api.example.com/v1"}, "", ""); why != "" {
		t.Error("a fixed entry refused an ask with none")
	}
}

// Regions are a closed list per entry, the owner picking one; nothing else.
func TestRegionsAreAClosedList(t *testing.T) {
	a := &apiAsk{URL: "https://example.com/spec.json", Bases: []string{"https://api.example.com", "https://api.eu.example.com"}}
	got, why := resolveEntryAPI(a, "", "https://api.eu.example.com")
	if why != "" || got.Base != "https://api.eu.example.com" || len(got.Bases) != 0 {
		t.Fatalf("eu: %q %+v", why, got)
	}
	for _, bad := range []string{"", "https://api.evil.example", "https://api.eu.example.com/", "http://api.eu.example.com", "https://api.eu.example.com:444"} {
		if got, why := resolveEntryAPI(a, "", bad); why == "" {
			t.Errorf("region %q was accepted as %s", bad, got.Base)
		}
	}
	if _, why := resolveEntryAPI(a, "acme", ""); why == "" {
		t.Error("an account was taken on an entry with regions")
	}
}

// The catalog takes a template or regions only as they are described above.
func TestCatalogEntriesWithTemplatesAndRegionsAreChecked(t *testing.T) {
	entries, _ := shippedKeyEntry(t)
	base := entries["qonto-key@1"]
	with := func(mut func(*CatalogEntry)) error {
		e := base
		api := *base.API
		e.API = &api
		mut(&e)
		e.ID = "qonto-key@1"
		return e.validate(nil)
	}
	if err := with(func(e *CatalogEntry) { e.API.Base = "https://{account}.qonto.com/v2" }); err != nil {
		t.Fatalf("a template: %v", err)
	}
	if err := with(func(e *CatalogEntry) {
		e.API.Base, e.API.Bases = "", []string{"https://a.qonto.com/v2", "https://b.qonto.com/v2"}
	}); err != nil {
		t.Fatalf("regions: %v", err)
	}
	for name, mut := range map[string]func(*CatalogEntry){
		"every part variable":    func(e *CatalogEntry) { e.API.Base = "https://{domain}.{tld}/v2" },
		"a variable in the path": func(e *CatalogEntry) { e.API.Base = "https://thirdparty.qonto.com/{account}/v2" },
		"a template and regions": func(e *CatalogEntry) { e.API.Bases = []string{"https://a.qonto.com/v2", "https://b.qonto.com/v2"} },
		"one region":             func(e *CatalogEntry) { e.API.Base, e.API.Bases = "", []string{"https://a.qonto.com/v2"} },
		"a template as a region": func(e *CatalogEntry) {
			e.API.Base, e.API.Bases = "", []string{"https://{account}.qonto.com/v2", "https://b.qonto.com/v2"}
		},
		"a repeated region": func(e *CatalogEntry) {
			e.API.Base, e.API.Bases = "", []string{"https://a.qonto.com/v2", "https://a.qonto.com/v2"}
		},
		"a region off the pin": func(e *CatalogEntry) {
			e.API.Base, e.API.Bases = "", []string{"https://a.qonto.com/v2", "https://b.qonto.com/v3"}
		},
		"a plain-http region": func(e *CatalogEntry) {
			e.API.Base, e.API.Bases = "", []string{"https://a.qonto.com/v2", "http://b.qonto.com/v2"}
		},
	} {
		if err := with(mut); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func shippedKeyEntry(t *testing.T) (map[string]CatalogEntry, []byte) {
	t.Helper()
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	return entries, description
}

// A grant from an entry at an account's address is made at the address the
// owner's account names, which the approval shows and the grant keeps; a taf
// has no way to move it.
func TestAGrantFromAnAccountEntryKeepsItsAddress(t *testing.T) {
	entries, description := shippedKeyEntry(t)
	entry := entries["qonto-key@1"]
	api := *entry.API
	api.Base = "https://{account}.qonto.com/v2"
	entry.API = &api
	e := newInProcessEnv(t)
	e.enrol()
	e.s.CatalogEntries = map[string]CatalogEntry{"qonto-key@1": entry}
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		return profileResponse(200, string(description)), nil
	})}
	ask := func(body map[string]string) (int, map[string]any) {
		body["agent"] = "agt_taf"
		return e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@1/grants", body, nil, "k_setup")
	}
	for _, bad := range []map[string]string{{}, {"account": "acme.evil.com"}, {"account": "acme.qonto.com.evil.com"}, {"account": "acme@evil"},
		{"account": "acme:444"}, {"account": strings.Repeat("a", 64)}, {"base": "https://evil.example/v2"}} {
		if st, out := ask(bad); st != 400 {
			t.Errorf("%v: %d %v", bad, st, out)
		}
	}
	st, out := ask(map[string]string{"account": "Acme"})
	if st != 201 || out["needs_key"] != true {
		t.Fatalf("acme: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Request.API.Base != "https://acme.qonto.com/v2" || g.Request.CatalogEntry != "qonto-key@1" {
		t.Fatalf("the grant's address: %+v", g.Request.API)
	}
	if !strings.Contains(g.Text, "acme.qonto.com") {
		t.Fatalf("the approval doesn't show the address: %s", g.Text)
	}
	if strings.Contains(g.Text, "{account}") {
		t.Fatalf("the approval shows the template: %s", g.Text)
	}
}
