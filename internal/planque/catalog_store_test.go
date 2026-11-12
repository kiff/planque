package planque

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/store"
)

// catalogDevice is an operator's enrolled device, for a test.
type catalogDevice struct {
	key *ecdsa.PrivateKey
	pub string
}

func newCatalogDevice(t *testing.T) catalogDevice {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, _ := x509.MarshalPKIXPublicKey(&k.PublicKey)
	return catalogDevice{key: k, pub: base64.StdEncoding.EncodeToString(der)}
}

func (d catalogDevice) sign(t *testing.T, base, action, id, digest, nonce, at string) string {
	t.Helper()
	sum := sha256.Sum256(CatalogSignMessage(base, action, id, digest, nonce, at))
	sig, err := ecdsa.SignASN1(rand.Reader, d.key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(sig)
}

func tokenSHA(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// catalogEnv is planque with two operators, Ada (who drafts) and Ben (who
// publishes), each with a credential and a device, and an agent's
// credential with no device. Qonto by key's description is served.
func catalogEnv(t *testing.T) (*env, CatalogEntry, map[string]catalogDevice) {
	t.Helper()
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	shipped, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	e := newInProcessEnv(t)
	e.enrol()
	devices := map[string]catalogDevice{"op_ada": newCatalogDevice(t), "op_ben": newCatalogDevice(t)}
	raw, _ := json.Marshal([]map[string]string{
		{"id": "op_ada", "token_sha256": tokenSHA("t_ada"), "device_key": devices["op_ada"].pub},
		{"id": "op_ben", "token_sha256": tokenSHA("t_ben"), "device_key": devices["op_ben"].pub},
		{"id": "op_agent", "token_sha256": tokenSHA("t_agent")},
	})
	if e.s.CatalogOperators, err = DecodeCatalogOperators(raw); err != nil {
		t.Fatal(err)
	}
	e.s.CatalogEntries = shipped
	entry := shipped["qonto-key@1"]
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		return profileResponse(200, string(description)), nil
	})}
	return e, entry, devices
}

func entryBody(t *testing.T, e CatalogEntry) map[string]any {
	t.Helper()
	doc, _, err := entryDocument(e)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	_ = json.Unmarshal(doc, &m)
	return map[string]any{"entry": m}
}

// The verified tier's whole life: an operator drafts and reviews, another
// operator's device publishes, and planque serves it then, and only as
// signed; retired, it isn't served.
func TestTheCatalogKeptAsData(t *testing.T) {
	e, qonto, devices := catalogEnv(t)
	base := "/v1/catalog/entries"
	entry := qonto
	entry.ID = "qonto-read@1"
	entry.Title = "Qonto, read only"
	var p map[string]any
	_ = json.Unmarshal(entry.Profile, &p)
	p["catalog_entry"] = entry.ID
	entry.Profile, _ = json.Marshal(p)

	// No credential, or a tenant's key: nothing.
	if st, _ := e.do("POST", base, entryBody(t, entry), nil, ""); st != 401 {
		t.Fatalf("no credential: %d", st)
	}
	if st, _ := e.do("POST", base, entryBody(t, entry), nil, "k_setup"); st != 401 {
		t.Fatalf("a tenant's key: %d", st)
	}
	// An agent drafts (it has no device).
	st, out := e.do("POST", base, entryBody(t, entry), nil, "t_agent")
	if st != 201 || out["state"] != CatalogDraft || out["drafted_by"] != "op_agent" {
		t.Fatalf("draft: %d %v", st, out)
	}
	digest := out["digest"].(string)
	if st, out := e.do("POST", base, entryBody(t, entry), nil, "t_ada"); st != 409 || !strings.Contains(out["error"].(string), "qonto-read@2") {
		t.Fatalf("a revision is never replaced: %d %v", st, out)
	}
	// Not served while a draft.
	if _, served := e.s.catalogEntries()["qonto-read@1"]; served {
		t.Fatal("a draft was served")
	}
	// A device reads it with no credential, and planque's own reading of
	// each tool from the pinned description, not the entry's labels.
	st, one := e.do("GET", base+"/qonto-read@1", nil, nil, "")
	described, _ := one["described"].(map[string]any)
	tools, _ := described["tools"].([]any)
	if st != 200 || one["digest"] != digest || described["from"] != "description" || len(tools) != len(entry.Tools) {
		t.Fatalf("one revision: %d %v", st, one)
	}
	for _, x := range tools {
		tool := x.(map[string]any)
		if tool["method"] != "GET" || tool["read_only"] != true || tool["why"] != nil {
			t.Fatalf("a Qonto read as planque reads it: %v", tool)
		}
	}
	publish := func(op, action, id, digest, nonce string, at time.Time) (int, map[string]any) {
		t.Helper()
		ts := at.UTC().Format(time.RFC3339)
		d, enrolled := devices[op]
		if !enrolled {
			d = newCatalogDevice(t) // a key no operator enrolled
		}
		return e.do("POST", base+"/"+id+"/"+action, map[string]string{"operator": op, "digest": digest, "nonce": nonce,
			"signed_at": ts, "signature": d.sign(t, e.url, action, id, digest, nonce, ts)}, nil, "")
	}
	now := e.s.now()
	if st, out := publish("op_ben", "publish", "qonto-read@1", digest, "nonce-draft-000001", now); st != 409 || !strings.Contains(out["error"].(string), "reviewed") {
		t.Fatalf("a draft published: %d %v", st, out)
	}
	if st, out := e.do("POST", base+"/qonto-read@1/review", map[string]string{"digest": "0" + digest[1:]}, nil, "t_ada"); st != 409 {
		t.Fatalf("reviewed what wasn't read: %d %v", st, out)
	}
	if st, out := e.do("POST", base+"/qonto-read@1/review", map[string]string{"digest": digest}, nil, "t_ada"); st != 200 || out["state"] != CatalogReviewed {
		t.Fatalf("review: %d %v", st, out)
	}
	// Ada changes it: a draft again, changed by her, so only Ben may publish it later.
	entry.Icon = qonto.Icon
	entry.Title = "Qonto, read only, reviewed"
	st, out = e.do("PUT", base+"/qonto-read@1", entryBody(t, entry), nil, "t_ada")
	if st != 200 || out["state"] != CatalogDraft || out["changed_by"] != "op_ada" || out["digest"] == digest {
		t.Fatalf("change: %d %v", st, out)
	}
	digest = out["digest"].(string)
	if st, _ := e.do("POST", base+"/qonto-read@1/review", map[string]string{"digest": digest}, nil, "t_ben"); st != 200 {
		t.Fatalf("review again: %d", st)
	}
	// Who last changed it can't publish it; a signature over another digest,
	// another planque, a stale one or one from no device doesn't count.
	if st, out := publish("op_ada", "publish", "qonto-read@1", digest, "nonce-ada-00000001", now); st != 403 || !strings.Contains(out["error"].(string), "last changed") {
		t.Fatalf("the changer published: %d %v", st, out)
	}
	if st, _ := publish("op_agent", "publish", "qonto-read@1", digest, "nonce-agent-000001", now); st != 403 {
		t.Fatalf("a credential with no device published: %d", st)
	}
	if st, _ := publish("op_ben", "publish", "qonto-read@1", "0"+digest[1:], "nonce-other-000001", now); st != 409 {
		t.Fatalf("another digest: %d", st)
	}
	if st, _ := publish("op_ben", "publish", "qonto-read@1", digest, "nonce-stale-000001", now.Add(-6*time.Minute)); st != 403 {
		t.Fatalf("a stale signature: %d", st)
	}
	ts := now.UTC().Format(time.RFC3339)
	elsewhere := devices["op_ben"].sign(t, "https://other.example", "publish", "qonto-read@1", digest, "nonce-else-0000001", ts)
	if st, _ := e.do("POST", base+"/qonto-read@1/publish", map[string]string{"operator": "op_ben", "digest": digest, "nonce": "nonce-else-0000001", "signed_at": ts, "signature": elsewhere}, nil, ""); st != 403 {
		t.Fatalf("a signature for another planque: %d", st)
	}
	st, out = publish("op_ben", "publish", "qonto-read@1", digest, "nonce-ben-00000001", now)
	if st != 200 || out["state"] != CatalogPublished {
		t.Fatalf("publish: %d %v", st, out)
	}
	if got, served := e.s.catalogEntries()["qonto-read@1"]; !served || got.Title != "Qonto, read only, reviewed" {
		t.Fatalf("published, not served: %+v", got)
	}
	// The same signature twice doesn't count; a published revision never changes.
	if st, _ := publish("op_ben", "retire", "qonto-read@1", digest, "nonce-ben-00000001", now); st != 409 {
		t.Fatalf("a nonce used twice: %d", st)
	}
	if st, _ := e.do("PUT", base+"/qonto-read@1", entryBody(t, entry), nil, "t_ada"); st != 409 {
		t.Fatalf("a published revision changed: %d", st)
	}
	// Its record names the device and the digest.
	records, _ := e.kv.Query(t.Context(), catalogPK, "RECORD#")
	if len(records) != 1 || !strings.Contains(string(records[0].Doc), digest) || !strings.Contains(string(records[0].Doc), "op_ben") {
		t.Fatalf("records: %d", len(records))
	}
	// A row edited behind planque's back isn't served.
	it, _, _ := e.kv.Get(t.Context(), catalogPK, "ENTRY#qonto-read@1")
	tampered := strings.Replace(string(it.Doc), "Qonto, read only, reviewed", "Qonto, read and write", 1)
	it.Doc = []byte(tampered)
	_ = e.kv.Put(t.Context(), it, storeAny())
	e.s.forgetCatalog()
	if _, served := e.s.catalogEntries()["qonto-read@1"]; served {
		t.Fatal("an entry changed after its signature was served")
	}
	// So isn't a row marked published with a matching digest and a signature
	// no enrolled device made.
	var forged catalogRow
	_ = json.Unmarshal([]byte(strings.Replace(tampered, "Qonto, read and write", "Qonto, read only, reviewed", 1)), &forged)
	good := *forged.Published
	forged.Published.Signature = devices["op_ada"].sign(t, e.url, "publish", "qonto-read@1", digest, "nonce-forged-00001", ts)
	forged.Published.Operator = "op_ben"
	it.Doc, _ = json.Marshal(forged)
	_ = e.kv.Put(t.Context(), it, storeAny())
	e.s.forgetCatalog()
	if _, served := e.s.catalogEntries()["qonto-read@1"]; served {
		t.Fatal("an entry with a forged signature was served")
	}
	forged.Published = &good
	it.Doc, _ = json.Marshal(forged)
	_ = e.kv.Put(t.Context(), it, storeAny())
	e.s.forgetCatalog()
	if _, served := e.s.catalogEntries()["qonto-read@1"]; !served {
		t.Fatal("the entry as signed, not served")
	}
	// Retired: not served, nor is the shipped revision of the same id.
	if st, out := publish("op_ben", "retire", "qonto-read@1", digest, "nonce-ben-00000002", now); st != 200 || out["state"] != CatalogRetired {
		t.Fatalf("retire: %d %v", st, out)
	}
	if _, served := e.s.catalogEntries()["qonto-read@1"]; served {
		t.Fatal("a retired entry was served")
	}
	// The shipped entries are served all along.
	if _, served := e.s.catalogEntries()["qonto-key@1"]; !served {
		t.Fatal("the shipped entry went")
	}
}

// A write that a value of one of its arguments makes final isn't published
// until that argument is pinned to the value that keeps it undoable.
func TestAWriteThatCanBeFinalIsPinnedBeforePublishing(t *testing.T) {
	create := openapi.Operation{Name: "create_invoice", Method: "POST", InputSchema: json.RawMessage(`{"type":"object","properties":{
		"order":{"type":"string"},
		"status":{"type":"string","enum":["draft","finalized"],"default":"finalized"},
		"send_email":{"type":"boolean"},
		"currency":{"type":"string","enum":["eur","usd"]}}}`)}
	got, err := FinalArguments(create)
	if err != nil || len(got) != 2 || got[0].Argument != "send_email" || got[0].Undo != false || got[1].Argument != "status" || got[1].Undo != "draft" {
		t.Fatalf("final arguments: %+v %v", got, err)
	}
	if read := (openapi.Operation{Name: "list", ReadOnly: true, InputSchema: create.InputSchema}); func() bool { f, _ := FinalArguments(read); return f != nil }() {
		t.Fatal("a read has no final argument")
	}
	for name, c := range map[string]struct {
		schema string
		want   string // "path=undo", "" for none, "error" when unreadable
	}{
		// OpenAPI 3.1's nullable types don't turn the check off.
		"a nullable property beside it": {`{"properties":{"note":{"type":["string","null"]},"status":{"type":"string","enum":["draft","unpaid"]}}}`, `status="draft"`},
		"nested in an object":           {`{"properties":{"invoice":{"type":"object","properties":{"status":{"type":"string","enum":["draft","finalized"]}}}}}`, `invoice.status="draft"`},
		"under the whole body":          {`{"properties":{"body":{"type":["object","null"],"properties":{"status":{"enum":["draft","issued"]}}}}}`, `body.status="draft"`},
		"a boolean for the draft state": {`{"properties":{"is_draft":{"type":"boolean","default":false}}}`, `is_draft=true`},
		"a boolean for a final step":    {`{"properties":{"send_now":{"type":"boolean"}}}`, `send_now=false`},
		"a final value not in any list": {`{"properties":{"status":{"type":"string","enum":["draft","open"],"default":"open"}}}`, `status="draft"`},
		"nothing undoable to pin to":    {`{"properties":{"state":{"type":"string","enum":["open","paid"]}}}`, ``},
		"an undoable value alone":       {`{"properties":{"status":{"enum":["draft"]}}}`, ``},
		"a type that isn't a schema's":  {`{"properties":{"status":{"type":7,"enum":["draft","sent"]}}}`, `error`},
		"properties that aren't a map":  {`{"properties":["status"]}`, `error`},
	} {
		got, err := FinalArguments(openapi.Operation{Name: "create_invoice", Method: "POST", InputSchema: json.RawMessage(c.schema)})
		var gotS string
		switch {
		case err != nil:
			gotS = "error"
		case len(got) == 1:
			u, _ := json.Marshal(got[0].Undo)
			gotS = got[0].Argument + "=" + string(u)
		case len(got) > 1:
			gotS = "several"
		}
		if gotS != c.want {
			t.Errorf("%s: %q, want %q", name, gotS, c.want)
		}
	}

	// Through a key entry's description: refused unpinned, then published pinned.
	spec := `{"openapi":"3.0.3","info":{"title":"Billing","version":"1"},"servers":[{"url":"https://billing.example/v1"}],
 "security":[{"k":[]}],"components":{"securitySchemes":{"k":{"type":"apiKey","in":"header","name":"X-Key"}}},
 "paths":{"/invoices":{"post":{"operationId":"create_invoice","requestBody":{"required":true,"content":{"application/json":{"schema":
  {"type":"object","properties":{"order":{"type":"string"},"status":{"type":"string","enum":["draft","finalized"]}}}}}}}}}}`
	sum := sha256.Sum256([]byte(spec))
	e := newInProcessEnv(t)
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		return profileResponse(200, spec), nil
	})}
	entry := CatalogEntry{ID: "billing@1", Kind: KindAPIKey, DescriptionSHA256: hex.EncodeToString(sum[:]),
		API:   &apiAsk{URL: "https://billing.example/openapi.json", Base: "https://billing.example/v1", Auth: openapi.Auth{In: "header", Name: "X-Key"}},
		Tools: []ToolRequest{{Tool: "create_invoice"}}}
	why, err := e.s.unpinnedFinal(t.Context(), entry)
	if err != nil || !strings.Contains(why, "create_invoice's status can make it final") || !strings.Contains(why, `"draft"`) {
		t.Fatalf("unpinned: %q %v", why, err)
	}
	entry.Tools[0].Fixed = map[string]json.RawMessage{"status": json.RawMessage(`"finalized"`)}
	if why, _ := e.s.unpinnedFinal(t.Context(), entry); why == "" {
		t.Fatal("pinned to final, published")
	}
	entry.Tools[0].Fixed = map[string]json.RawMessage{"status": json.RawMessage(` "draft" `)}
	if why, _ := e.s.unpinnedFinal(t.Context(), entry); why != "" {
		t.Fatalf("pinned to draft: %q", why)
	}
	// The entry's own read-only label doesn't count: the description says
	// it's a write. A tool the description doesn't have is refused.
	entry.Tools[0] = ToolRequest{Tool: "create_invoice", ReadOnly: true}
	if why, _ := e.s.unpinnedFinal(t.Context(), entry); !strings.Contains(why, "can make it final") {
		t.Fatalf("a write labelled read-only: %q", why)
	}
	entry.Tools[0] = ToolRequest{Tool: "send_invoice"}
	if why, _ := e.s.unpinnedFinal(t.Context(), entry); !strings.Contains(why, "send_invoice isn't in its description") {
		t.Fatalf("a tool not in the description: %q", why)
	}
	entry.DescriptionSHA256 = strings.Repeat("0", 64)
	if why, _ := e.s.unpinnedFinal(t.Context(), entry); !strings.Contains(why, "changed since it was pinned") {
		t.Fatalf("changed description: %q", why)
	}
}

// The operators' setting refuses what it can't use.
func TestCatalogOperatorsSetting(t *testing.T) {
	for name, raw := range map[string]string{
		"not a list":     `{"id":"op_a"}`,
		"bad id":         `[{"id":"ada","token_sha256":"` + tokenSHA("x") + `"}]`,
		"twice":          `[{"id":"op_a","token_sha256":"` + tokenSHA("x") + `"},{"id":"op_a","token_sha256":"` + tokenSHA("y") + `"}]`,
		"neither":        `[{"id":"op_a"}]`,
		"short sha":      `[{"id":"op_a","token_sha256":"abc"}]`,
		"bad device key": `[{"id":"op_a","device_key":"bm90IGEga2V5"}]`,
		"unknown field":  `[{"id":"op_a","token":"x"}]`,
	} {
		if _, err := DecodeCatalogOperators([]byte(raw)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

// storeAny is a write with no condition.
func storeAny() store.Cond { return store.Cond{} }

// What's offered is the newest revision of each name.
func TestOnlyTheNewestRevisionIsOffered(t *testing.T) {
	got := newestRevisions(map[string]CatalogEntry{"desk@1": {}, "desk@10": {}, "desk@2": {}, "mail@3": {}})
	if len(got) != 2 {
		t.Fatalf("%v", got)
	}
	if _, ok := got["desk@10"]; !ok {
		t.Fatalf("not the newest: %v", got)
	}
}

// A revision's read takes no credential, so its description is downloaded
// and parsed once, not on every read: two reads of the same revision reach
// the description's host once. Only a description that matched its pin is
// kept; one that couldn't be read is read again next time.
func TestADescribedRevisionIsReadOnce(t *testing.T) {
	e, qonto, _ := catalogEnv(t)
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	fetched, down := 0, true
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		fetched++
		if down {
			return profileResponse(503, "down"), nil
		}
		return profileResponse(200, string(description)), nil
	})}
	entry := qonto
	entry.ID = "qonto-read@1"
	var p map[string]any
	_ = json.Unmarshal(entry.Profile, &p)
	p["catalog_entry"] = entry.ID
	entry.Profile, _ = json.Marshal(p)
	base := "/v1/catalog/entries"
	if st, out := e.do("POST", base, entryBody(t, entry), nil, "t_agent"); st != 201 {
		t.Fatalf("draft: %d %v", st, out)
	}
	read := func() map[string]any {
		t.Helper()
		st, one := e.do("GET", base+"/qonto-read@1", nil, nil, "")
		if st != 200 {
			t.Fatalf("read: %d %v", st, one)
		}
		d, _ := one["described"].(map[string]any)
		return d
	}
	// The host down: said, and not kept.
	if d := read(); !strings.Contains(fmt.Sprint(d["why"]), "couldn't be read") {
		t.Fatalf("a description that couldn't be read: %v", d)
	}
	mu.Lock()
	down, fetched = false, 0
	mu.Unlock()
	for i := 0; i < 3; i++ {
		if d := read(); d["from"] != "description" || d["why"] != nil || len(d["tools"].([]any)) != len(entry.Tools) {
			t.Fatalf("read %d: %v", i, d)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if fetched != 1 {
		t.Fatalf("the description was downloaded %d times for three reads, want 1", fetched)
	}
}
