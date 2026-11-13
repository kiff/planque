package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// A key request is declared only with operation names planque accepts, so
// Coach reads them first, before any grant: the names a request then takes,
// read from the public description, with nothing kept and no call to the
// vendor's API.
func TestAnAPIsOperationsAreReadBeforeAnyGrant(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	path := "/v1/accounts/" + tenant + "/api/operations"
	api := map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
	if st, _ := e.do("POST", path, map[string]any{"api": api}, nil, ""); st != 401 {
		t.Fatalf("no key: %d", st)
	}
	st, out := e.do("POST", path, map[string]any{"api": api}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("%d %v", st, out)
	}
	names := map[string]bool{}
	for _, o := range out["operations"].([]any) {
		m := o.(map[string]any)
		names[m["name"].(string)] = m["read_only"].(bool)
	}
	if ro, ok := names["listCharges"]; !ok || !ro || len(names) < 2 {
		t.Fatalf("operations: %v", out["operations"])
	}
	// Nothing was called at the vendor, and no grant was made.
	if calls := v.seen(); len(calls) != 0 {
		t.Fatalf("the vendor's API was called: %+v", calls)
	}
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/agents/agt_taf/grants", nil, nil, "k_setup"); st == 200 && strings.Contains(mustJSONString(out), "grt_") {
		t.Fatalf("a grant was made: %v", out)
	}
	// The names it gives are the ones a key request takes, and that request is declared.
	tool := ""
	for n, ro := range names {
		if ro {
			tool = n
			break
		}
	}
	st, out = e.askAPIGrant(v, []map[string]any{{"tool": tool, "title": "Read"}}, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("a request with a name it gave, not declared: %d %v", st, out)
	}
	// Refused as a request is: an unreadable description, a key place that isn't one.
	for _, bad := range []map[string]any{
		{"url": v.URL + "/missing?token=s3cret", "auth": map[string]any{"in": "header", "name": "X-API-Key"}},
		{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "cookie", "name": "k"}},
		{"url": "", "auth": map[string]any{"in": "header", "name": "X-API-Key"}},
	} {
		st, out := e.do("POST", path, map[string]any{"api": bad}, nil, "k_setup")
		if st != 400 || strings.Contains(mustJSONString(out), "s3cret") {
			t.Fatalf("%v: %d %v", bad, st, out)
		}
	}
}

// Each read fetches a description of up to 16 MB, so an account has 100 a
// day: past them, 429, and nothing is fetched.
func TestAnAccountsAPIReadsAreCountedADay(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	day := e.s.now().UTC().Format("2006-01-02")
	key := store.Item{PK: tenantPK(tenant), SK: "APIREADS#" + day}
	if err := putDoc(t.Context(), e.s.KV, key, allowance{Count: apiReadsADay - 1}, store.Cond{NotExists: true}); err != nil {
		t.Fatal(err)
	}
	api := map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
	path := "/v1/accounts/" + tenant + "/api/operations"
	if st, out := e.do("POST", path, map[string]any{"api": api}, nil, "k_setup"); st != 200 {
		t.Fatalf("the last of the day: %d %v", st, out)
	}
	st, out := e.do("POST", path, map[string]any{"api": api}, nil, "k_setup")
	if st != 429 || !strings.Contains(mustJSONString(out), "today") {
		t.Fatalf("past the day's reads: %d %v", st, out)
	}
}
