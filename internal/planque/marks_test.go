package planque

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiff/planque/internal/store"
)

var markNonce atomic.Int64

// markAsk sets a managed mark as an operator's device signs it.
func markAsk(t *testing.T, e *env, dev catalogDevice, operator, name string, hosts []string, follows string) (int, map[string]any) {
	t.Helper()
	digest := markDigest(name, hosts, follows)
	nonce := fmt.Sprintf("n%019d", markNonce.Add(1))
	at := e.s.now().UTC().Format(time.RFC3339)
	return e.do("POST", "/v1/marks/managed", map[string]any{"name": name, "hosts": hosts, "follows": follows, "operator": operator, "nonce": nonce, "signed_at": at,
		"signature": dev.sign(t, e.s.BaseURL, "mark-managed", provedKey(name), digest, nonce, at)}, nil, "")
}

func markWithdraw(t *testing.T, e *env, dev catalogDevice, operator, name, digest string) (int, map[string]any) {
	t.Helper()
	nonce := fmt.Sprintf("n%019d", markNonce.Add(1))
	at := e.s.now().UTC().Format(time.RFC3339)
	return e.do("POST", "/v1/marks/managed/"+provedKey(name)+"/withdraw", map[string]any{"operator": operator, "nonce": nonce, "signed_at": at,
		"signature": dev.sign(t, e.s.BaseURL, "unmark-managed", provedKey(name), digest, nonce, at)}, nil, "")
}

var intercomHosts = []string{"api.intercom.io", "api.eu.intercom.io", "api.au.intercom.io"}

// Only an enrolled operator's device, signing exactly this mark, once and
// recently, sets one.
func TestOnlyAnEnrolledDevicesSignatureSetsAManagedMark(t *testing.T) {
	e, _, devices := catalogEnv(t)
	ada, ben := devices["op_ada"], devices["op_ben"]
	good := func() map[string]any {
		digest := markDigest("Intercom", intercomHosts, "2026-11-17")
		at := e.s.now().UTC().Format(time.RFC3339)
		return map[string]any{"name": "Intercom", "hosts": intercomHosts, "follows": "2026-11-17", "operator": "op_ada", "nonce": "n0000000000000000001", "signed_at": at,
			"signature": ada.sign(t, e.s.BaseURL, "mark-managed", "intercom", digest, "n0000000000000000001", at)}
	}
	post := func(b map[string]any) int {
		st, _ := e.do("POST", "/v1/marks/managed", b, nil, "")
		return st
	}
	cases := map[string]func(map[string]any){
		"no signature":                func(b map[string]any) { b["signature"] = "" },
		"an operator with no device":  func(b map[string]any) { b["operator"] = "op_agent" },
		"another device's signature":  func(b map[string]any) { b["operator"] = "op_ben" },
		"hosts changed after signing": func(b map[string]any) { b["hosts"] = []string{"api.intercom.io", "api.evil.example"} },
		"another name":                func(b map[string]any) { b["name"] = "Intercom Pro" },
		"a stale signature": func(b map[string]any) {
			at := e.s.now().Add(-time.Hour).UTC().Format(time.RFC3339)
			b["signed_at"] = at
			b["signature"] = ada.sign(t, e.s.BaseURL, "mark-managed", "intercom", markDigest("Intercom", intercomHosts, "2026-11-17"), "n0000000000000000001", at)
		},
		"signed as a withdrawal": func(b map[string]any) {
			b["signature"] = ada.sign(t, e.s.BaseURL, "unmark-managed", "intercom", markDigest("Intercom", intercomHosts, "2026-11-17"), "n0000000000000000001", b["signed_at"].(string))
		},
	}
	for name, mutate := range cases {
		b := good()
		mutate(b)
		if st := post(b); st == http.StatusOK {
			t.Errorf("%s: accepted", name)
		}
	}
	_ = ben
	if st := post(good()); st != http.StatusOK {
		t.Fatalf("the good one: %d", st)
	}
	// A signature counts once.
	if st := post(good()); st != http.StatusConflict {
		t.Errorf("a replay: %d", st)
	}
	// What is asked is checked: no host list, a host with a port or an upper-case letter, a bad date.
	for name, b := range map[string]map[string]any{
		"no hosts":   {"name": "Intercom", "hosts": []string{}, "follows": "2026-11-17"},
		"a port":     {"name": "Intercom", "hosts": []string{"api.intercom.io:444"}, "follows": "2026-11-17"},
		"upper case": {"name": "Intercom", "hosts": []string{"API.intercom.io"}, "follows": "2026-11-17"},
		"a bare tld": {"name": "Intercom", "hosts": []string{"*.io"}, "follows": "2026-11-17"},
		"a bad date": {"name": "Intercom", "hosts": []string{"api.intercom.io"}, "follows": "yesterday"},
	} {
		if st := post(b); st != http.StatusBadRequest {
			t.Errorf("%s: %d", name, st)
		}
	}
}

// The gate follows a mark first, then the seed: a vendor with no mark follows
// the list; a withdrawn mark wins over it; a test tenant is let through whatever.
func TestTheGateFollowsTheMarkThenTheSeed(t *testing.T) {
	e, _, devices := catalogEnv(t)
	ada := devices["op_ada"]
	e.s.ManagedApps = map[string]bool{"googlesheets": true}
	e.s.ManagedTestTenants = map[string]bool{"tnt_test": true}
	if e.s.managedAllowed("tnt_x", "moneybird") || !e.s.managedAllowed("tnt_x", "googlesheets") || !e.s.managedAllowed("tnt_test", "moneybird") {
		t.Fatal("the seed and the test tenants, before any mark")
	}
	st, out := markAsk(t, e, ada, "op_ada", "Moneybird", []string{"moneybird.com"}, "2026-11-17")
	if st != http.StatusOK {
		t.Fatalf("set: %d %v", st, out)
	}
	if !e.s.managedAllowed("tnt_x", "Moneybird") || !e.s.managedAllowed("tnt_x", "moneybird") {
		t.Fatal("a marked vendor is allowed, by any spelling of its name")
	}
	st, out = markWithdraw(t, e, ada, "op_ada", "Moneybird", out["vendor"].(string)+"")
	if st == http.StatusOK {
		t.Fatalf("a withdrawal signed over the wrong digest was taken: %v", out)
	}
	cur := e.s.managedMarks()["moneybird"]
	if st, out := markWithdraw(t, e, ada, "op_ada", "Moneybird", cur.Digest); st != http.StatusOK || out["state"] != "withdrawn" {
		t.Fatalf("withdraw: %d %v", st, out)
	}
	if e.s.managedAllowed("tnt_x", "moneybird") {
		t.Fatal("a withdrawn mark stops new managed connections at once")
	}
	if !e.s.managedAllowed("tnt_test", "moneybird") {
		t.Fatal("a test tenant is never held by a mark")
	}
	// A seeded vendor, marked off, is off.
	if st, out := markAsk(t, e, ada, "op_ada", "Google Sheets", []string{"sheets.googleapis.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatalf("set sheets: %d %v", st, out)
	}
	cur = e.s.managedMarks()["googlesheets"]
	if st, _ := markWithdraw(t, e, ada, "op_ada", "Google Sheets", cur.Digest); st != http.StatusOK {
		t.Fatal("withdraw sheets")
	}
	if e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("a withdrawn mark wins over the seed")
	}
}

// An entry is managed only when its vendor is marked and every endpoint it
// has is on the mark's hosts: one host off the list, a look-alike vendor name
// pointing elsewhere, or a host added since the mark was signed is not.
func TestAnEntryIsManagedOnlyOnTheMarksHosts(t *testing.T) {
	e, _, devices := catalogEnv(t)
	ada := devices["op_ada"]
	intercom := CatalogEntry{Vendor: "Intercom", API: &apiAsk{Bases: intercomHosts2()}}
	if e.s.entryManaged(intercom) {
		t.Fatal("managed with no mark")
	}
	if st, out := markAsk(t, e, ada, "op_ada", "Intercom", []string{"api.intercom.io", "api.eu.intercom.io"}, "2026-11-17"); st != http.StatusOK {
		t.Fatalf("set: %d %v", st, out)
	}
	if e.s.entryManaged(intercom) {
		t.Fatal("an entry with a region the mark doesn't list was served as managed")
	}
	// The owner signs the mark again with the new host, and only then.
	if st, out := markAsk(t, e, ada, "op_ada", "Intercom", intercomHosts, "2026-11-17"); st != http.StatusOK {
		t.Fatalf("re-sign: %d %v", st, out)
	}
	if !e.s.entryManaged(intercom) {
		t.Fatal("not managed on its own marked hosts")
	}
	// A look-alike: Intercom's name, another host.
	if e.s.entryManaged(CatalogEntry{Vendor: "Intercom", API: &apiAsk{Base: "https://api.evil.example/v1"}}) ||
		e.s.entryManaged(CatalogEntry{Vendor: "Intercom", API: &apiAsk{Bases: []string{"https://api.intercom.io", "https://api.intercom.io.evil.example"}}}) ||
		e.s.entryManaged(CatalogEntry{Vendor: "Intercom", API: &apiAsk{Base: "https://api.intercom.io:444/v1"}}) ||
		e.s.entryManaged(CatalogEntry{Vendor: "Intercom", API: &apiAsk{Base: "http://api.intercom.io/v1"}}) {
		t.Fatal("an entry named for a managed vendor but pointing elsewhere was served as managed")
	}
	// An account template is read on the vendor's own domain.
	if st, out := markAsk(t, e, ada, "op_ada", "Freshdesk", []string{"*.freshdesk.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatalf("freshdesk: %d %v", st, out)
	}
	if !e.s.entryManaged(CatalogEntry{Vendor: "Freshdesk", API: &apiAsk{Base: "https://{account}.freshdesk.com/api/v2"}}) ||
		e.s.entryManaged(CatalogEntry{Vendor: "Freshdesk", API: &apiAsk{Base: "https://{account}.freshdesk.com.evil.example/api/v2"}}) {
		t.Fatal("a template is managed on the vendor's domain only")
	}
	// Withdrawn, no entry is.
	cur := e.s.managedMarks()["intercom"]
	if st, _ := markWithdraw(t, e, ada, "op_ada", "Intercom", cur.Digest); st != http.StatusOK {
		t.Fatal("withdraw")
	}
	if e.s.entryManaged(intercom) {
		t.Fatal("managed after the mark was withdrawn")
	}
}

func intercomHosts2() []string {
	out := make([]string, len(intercomHosts))
	for i, h := range intercomHosts {
		out[i] = "https://" + h
	}
	return out
}

// Every change of a mark is kept, in order and only added to; reading marks
// shows no signature; the catalog read says which entries are managed.
func TestTheMarkHistoryIsKeptAndTheCatalogSaysManaged(t *testing.T) {
	e, _, devices := catalogEnv(t)
	ada := devices["op_ada"]
	// Three changes within a fraction of a second, at times whose text sorts
	// wrongly when trailing zeros are dropped (.12, .123, .125): the history
	// must still read in the order they happened.
	base := time.Now().UTC().Truncate(time.Second)
	var offset time.Duration
	e.s.Now = func() time.Time { return base.Add(offset) }
	offset = 120 * time.Millisecond
	if st, _ := markAsk(t, e, ada, "op_ada", "Qonto", []string{"thirdparty.qonto.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatal("set")
	}
	_, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
	managed := 0
	for _, x := range out["entries"].([]any) {
		m := x.(map[string]any)
		if strings.HasPrefix(m["entry"].(string), "qonto-key@") && m["managed"] == true {
			managed++
		}
	}
	if managed == 0 {
		t.Fatalf("the catalog read doesn't say Qonto's entry is managed: %v", out)
	}
	cur := e.s.managedMarks()["qonto"]
	offset = 123 * time.Millisecond
	if st, _ := markWithdraw(t, e, ada, "op_ada", "Qonto", cur.Digest); st != http.StatusOK {
		t.Fatal("withdraw")
	}
	offset = 125 * time.Millisecond
	if st, _ := markAsk(t, e, ada, "op_ada", "Qonto", []string{"thirdparty.qonto.com"}, "2026-11-18"); st != http.StatusOK {
		t.Fatal("set again")
	}
	items, err := e.s.KV.Query(t.Context(), markPK, "LOG#")
	if err != nil || len(items) != 3 {
		t.Fatalf("history: %d %v", len(items), err)
	}
	var states []string
	for _, it := range items {
		var m managedMark
		_ = json.Unmarshal(it.Doc, &m)
		states = append(states, fmt.Sprintf("%s#%d", m.State, m.Seq))
	}
	if strings.Join(states, " ") != "active#1 withdrawn#2 active#3" {
		t.Fatalf("history in order: %v", states)
	}
	st, body := e.do("GET", "/v1/marks", nil, nil, "")
	raw, _ := json.Marshal(body)
	if st != http.StatusOK || strings.Contains(string(raw), "signature") || strings.Contains(string(raw), "nonce") {
		t.Fatalf("the marks read: %d %s", st, raw)
	}
}

// A row written into the table by anything but an enrolled device's signed
// request isn't a mark, and a mark whose fields changed after signing isn't.
func TestAForgedOrAlteredMarkRowIsNotAMark(t *testing.T) {
	e, _, devices := catalogEnv(t)
	forged, _ := json.Marshal(managedMark{Vendor: "evilcorp", Name: "Evilcorp", Hosts: []string{"evil.example"}, Follows: "2026-11-17", State: markActive,
		Digest: markDigest("Evilcorp", []string{"evil.example"}, "2026-11-17")})
	if err := e.s.KV.Put(t.Context(), storeItem(markPK, "MANAGED#evilcorp", forged), storeCond()); err != nil {
		t.Fatal(err)
	}
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "evilcorp") {
		t.Fatal("an unsigned row opened the gate")
	}
	// A real mark with its hosts widened in the table after signing.
	if st, _ := markAsk(t, e, devices["op_ada"], "op_ada", "Moneybird", []string{"moneybird.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatal("set")
	}
	it, _, _ := e.s.KV.Get(t.Context(), markPK, "MANAGED#moneybird")
	var m managedMark
	_ = json.Unmarshal(it.Doc, &m)
	m.Hosts = append(m.Hosts, "evil.example")
	doc, _ := json.Marshal(m)
	ver := it.Ver
	if err := e.s.KV.Put(t.Context(), storeItem(markPK, "MANAGED#moneybird", doc), storeCondVer(&ver)); err != nil {
		t.Fatal(err)
	}
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "moneybird") || e.s.entryManaged(CatalogEntry{Vendor: "Moneybird", API: &apiAsk{Base: "https://evil.example/v2"}}) {
		t.Fatal("a mark altered after signing was honoured")
	}
}

func storeItem(pk, sk string, doc []byte) store.Item { return store.Item{PK: pk, SK: sk, Doc: doc} }
func storeCond() store.Cond                          { return store.Cond{NotExists: true} }
func storeCondVer(v *int64) store.Cond               { return store.Cond{Ver: v} }

// A mark is read only while the operator whose device signed it is still
// enrolled: removing an operator (rotating their device out) stops what they
// signed from counting, as it does for a published catalog revision, and
// enrolling them again, or another operator signing the mark again, brings it
// back.
func TestMarksSignedByAnOperatorNoLongerEnrolledStopCounting(t *testing.T) {
	e, _, devices := catalogEnv(t)
	if st, _ := markAsk(t, e, devices["op_ben"], "op_ben", "Moneybird", []string{"moneybird.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatal("set")
	}
	entry := CatalogEntry{Vendor: "Moneybird", API: &apiAsk{Base: "https://moneybird.com/api/v2"}}
	if !e.s.managedAllowed("tnt_x", "moneybird") || !e.s.entryManaged(entry) {
		t.Fatal("a mark by an enrolled operator")
	}
	all := e.s.CatalogOperators
	rotated := map[string]CatalogOperator{}
	for id, op := range all {
		if id != "op_ben" {
			rotated[id] = op
		}
	}
	e.s.CatalogOperators = rotated
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "moneybird") || e.s.entryManaged(entry) {
		t.Fatal("a mark by an operator no longer enrolled still counts")
	}
	if got := e.s.managedMarks(); len(got) != 0 {
		t.Fatalf("marks read with no signer: %v", got)
	}
	// Another operator signs it again: it counts under that signature.
	if st, _ := markAsk(t, e, devices["op_ada"], "op_ada", "Moneybird", []string{"moneybird.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatal("signed again")
	}
	if !e.s.managedAllowed("tnt_x", "moneybird") {
		t.Fatal("a mark signed again by an enrolled operator")
	}
	// The first operator back: nothing changes the mark in the table.
	e.s.CatalogOperators = all
	e.s.forgetMarks()
	if !e.s.managedAllowed("tnt_x", "moneybird") {
		t.Fatal("still marked")
	}
}

// A reproduction: rotating an operator out also drops the withdrawals
// they signed, and a seeded vendor must not become managed again. A row that
// exists but no longer counts reads off, never the seed.
func TestRotatingOutAWithdrawingOperatorDoesNotBringTheSeedBack(t *testing.T) {
	e, _, devices := catalogEnv(t)
	e.s.ManagedApps = map[string]bool{"googlesheets": true}
	ben := devices["op_ben"]
	if st, _ := markAsk(t, e, ben, "op_ben", "Google Sheets", []string{"sheets.googleapis.com"}, "2026-11-17"); st != http.StatusOK {
		t.Fatal("set")
	}
	cur := e.s.managedMarks()["googlesheets"]
	if st, _ := markWithdraw(t, e, ben, "op_ben", "Google Sheets", cur.Digest); st != http.StatusOK {
		t.Fatal("withdraw")
	}
	if e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("withdrawn, and the seed let it through")
	}
	all := e.s.CatalogOperators
	rotated := map[string]CatalogOperator{}
	for id, op := range all {
		if id != "op_ben" {
			rotated[id] = op
		}
	}
	e.s.CatalogOperators = rotated
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("rotating the withdrawing operator out brought the seeded vendor back")
	}
	// Every other seeded vendor, with no row at all, still follows the seed.
	e.s.ManagedApps["moneybird"] = true
	if !e.s.managedAllowed("tnt_x", "moneybird") {
		t.Fatal("a vendor with no row follows the seed")
	}
	// A row altered after signing reads off too.
	it, _, _ := e.s.KV.Get(t.Context(), markPK, "MANAGED#googlesheets")
	var m managedMark
	_ = json.Unmarshal(it.Doc, &m)
	m.State = markActive
	doc, _ := json.Marshal(m)
	ver := it.Ver
	e.s.CatalogOperators = all
	if err := e.s.KV.Put(t.Context(), storeItem(markPK, "MANAGED#googlesheets", doc), storeCondVer(&ver)); err != nil {
		t.Fatal(err)
	}
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("a withdrawal turned back to active by editing the row was honoured")
	}
}

// failQuery is a table that can't be read.
type failQuery struct{ store.KV }

func (failQuery) Query(context.Context, string, string) ([]store.Item, error) {
	return nil, errors.New("down")
}

// While the table can't be read, and nothing was read before, the gate
// refuses managed connections: it never follows a seed it can't check
// against a withdrawal. A test tenant is still let through.
func TestAnUnreadableMarksTableRefusesManaged(t *testing.T) {
	e, _, _ := catalogEnv(t)
	e.s.ManagedApps = map[string]bool{"googlesheets": true}
	e.s.ManagedTestTenants = map[string]bool{"tnt_test": true}
	good := e.s.KV
	e.s.KV = failQuery{good}
	e.s.forgetMarks()
	if e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("a seeded vendor was allowed while the marks couldn't be read")
	}
	if !e.s.managedAllowed("tnt_test", "googlesheets") {
		t.Fatal("a test tenant is let through whatever")
	}
	e.s.KV = good
	e.s.forgetMarks()
	if !e.s.managedAllowed("tnt_x", "googlesheets") {
		t.Fatal("readable again: the seed")
	}
}
