package planque

import (
	"fmt"
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

// A big API (Stripe's 608 operations, GitHub's 1,277) is described by words
// from the job: at most 50 operations, best first, with the total. With no
// words it answers its count and how to ask, never every name; a small one
// still lists whole. Listing never widens what a request takes: still 500.
func TestABigAPIIsDescribedByWords(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	var paths []string
	for i := range 600 {
		paths = append(paths, fmt.Sprintf(`"/things/%d": {"get": {"operationId": "getThing%d", "summary": "Read thing %d"}}`, i, i, i))
	}
	paths = append(paths,
		`"/charges/{c}/refunds": {"post": {"operationId": "createChargeRefund", "summary": "Refund a charge", "parameters": [{"name": "c", "in": "path", "required": true}]}}`,
		// Stuffed with every word: it may match, but never above the one a job names.
		`"/stuffed": {"get": {"operationId": "stuffed", "summary": "refund charge refund charge create refund for a charge"}}`)
	big := `{"openapi": "3.0.3", "info": {"title": "Big"}, "security": [{"k": []}], "components": {"securitySchemes": {"k": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}}, "paths": {` + strings.Join(paths, ",") + `}}`
	v.mu.Lock()
	v.desc = big
	v.mu.Unlock()
	path := "/v1/accounts/" + tenant + "/api/operations"
	api := map[string]any{"url": v.URL + "/openapi.json", "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
	// No words: the count and how to ask.
	st, out := e.do("POST", path, map[string]any{"api": api}, nil, "k_setup")
	if st != 200 || out["total"] != float64(602) || len(out["operations"].([]any)) != 0 || !strings.Contains(fmt.Sprint(out["more"]), "describe again with q") {
		t.Fatalf("no words: %d total %v ops %d more %v", st, out["total"], len(out["operations"].([]any)), out["more"])
	}
	// Words: the job's operation first, at most 50, with the total and how many matched.
	st, out = e.do("POST", path, map[string]any{"api": api, "q": "refund a charge"}, nil, "k_setup")
	ops := out["operations"].([]any)
	if st != 200 || len(ops) == 0 || ops[0].(map[string]any)["name"] != "createChargeRefund" || out["total"] != float64(602) {
		t.Fatalf("refund a charge: %d %v", st, out)
	}
	st, out = e.do("POST", path, map[string]any{"api": api, "q": "thing"}, nil, "k_setup")
	if st != 200 || len(out["operations"].([]any)) != 50 || out["matched"] != float64(600) {
		t.Fatalf("thing: %d matched %v listed %d", st, out["matched"], len(out["operations"].([]any)))
	}
	// Words that match nothing: none, said as zero, never a refusal.
	st, out = e.do("POST", path, map[string]any{"api": api, "q": "zebra quokka"}, nil, "k_setup")
	if st != 200 || len(out["operations"].([]any)) != 0 || out["matched"] != float64(0) {
		t.Fatalf("no match: %d %v", st, out)
	}
	if st, _ := e.do("POST", path, map[string]any{"api": api, "q": strings.Repeat("w", 201)}, nil, "k_setup"); st != 400 {
		t.Fatalf("a q too long: %d", st)
	}
	// A request with the job's operation, as describe named it: its snapshot is that operation alone, declared.
	st, out = e.askAPIGrant(v, []map[string]any{{"tool": "createChargeRefund", "title": "Refund a charge"}}, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("a big API's request with its operation: %d %v", st, out)
	}
	// One with no tools still reads the whole, and a whole is still at most 500: listing never widened it.
	st, out = e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 400 || !strings.Contains(mustJSONString(out), "more than 500 operations") {
		t.Fatalf("a whole request of a big API: %d %v", st, out)
	}
}
