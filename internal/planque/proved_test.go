package planque

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestTheShippedProofsRead(t *testing.T) {
	p, err := DecodeProved(provedJSON)
	if err != nil || len(p) == 0 {
		t.Fatalf("proved.json: %v %v", p, err)
	}
	for name, proof := range p {
		if !strings.HasPrefix(proof.Link, "https://") || proof.How == "" {
			t.Errorf("%s has no public link: %+v", name, proof)
		}
	}
}

func TestProofsRefuseWhatTheyCantRead(t *testing.T) {
	for _, c := range []struct{ name, raw, want string }{
		{"not an object", `[]`, "object of apps"},
		{"two objects", `{} {}`, "object of apps"},
		{"null", `null`, "object of apps"},
		{"an unknown field", `{"gmail":{"at":"2026-11-13","how":"key","by":"me"}}`, "object of apps"},
		{"an unfolded name", `{"Google Calendar":{"at":"2026-11-16","how":"oauth"}}`, "folded name"},
		{"a dash in the name", `{"google-calendar":{"at":"2026-11-16","how":"oauth"}}`, "folded name"},
		{"another route", `{"gmail":{"at":"2026-11-13","how":"oauth","route":"mine"}}`, "neither managed nor own"},
		{"a route on a key proof", `{"qonto":{"at":"2026-11-13","how":"key","route":"own"}}`, "narrows a sign-in"},
		{"no how", `{"gmail":{"at":"2026-11-13"}}`, "neither key nor oauth"},
		{"another how", `{"gmail":{"at":"2026-11-13","how":"sign-in"}}`, "neither key nor oauth"},
		{"an entry with its revision", `{"qonto":{"at":"2026-11-13","how":"key","entries":["qonto-key@2"]}}`, "reviewed entry's name"},
		{"no day", `{"gmail":{"how":"oauth"}}`, "isn't a day"},
		{"a day with a time", `{"gmail":{"at":"2026-11-13T10:00:00Z","how":"oauth"}}`, "isn't a day"},
		{"an http link", `{"gmail":{"at":"2026-11-13","how":"oauth","link":"http://example.com/a"}}`, "https"},
		{"a link with a user", `{"gmail":{"at":"2026-11-13","how":"oauth","link":"https://me@example.com/a"}}`, "https"},
		{"a relative link", `{"gmail":{"at":"2026-11-13","how":"oauth","link":"/a"}}`, "https"},
	} {
		if p, err := DecodeProved([]byte(c.raw)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v %v, want %q", c.name, p, err, c.want)
		}
	}
	if p, err := DecodeProved([]byte(`{"gmail":{"at":"2026-11-13","how":"oauth"}}`)); err != nil || p["gmail"].At != "2026-11-13" || p["gmail"].How != "oauth" {
		t.Errorf("a proof with no link: %v %v", p, err)
	}
}

// A proof matches by the slug planque names an entry with, or one of its
// reviewed entries, and only the way it was proved; never by the vendor's
// words, which the provider supplies.
func TestAProofMatchesTheAppsSlugTheWayItWasRun(t *testing.T) {
	proved := map[string]Proof{
		"googlecalendar": {At: "2026-11-16", How: "oauth"},
		"qonto":          {At: "2026-11-13", How: "key", Entries: []string{"qonto-key"}},
		"googlesheets":   {At: "2026-11-21", How: "oauth", Route: "managed"},
		"notion":         {At: "2026-11-21", How: "oauth", Route: "own"},
	}
	for _, c := range []struct{ entry, how, want string }{
		{"connector:composio:googlecalendar", "oauth", "2026-11-16"},
		{"connector:composio:GOOGLE_CALENDAR", "oauth", "2026-11-16"},
		{"connector:composio:googlecalendar", "key", ""},
		{"qonto-key@2", "key", "2026-11-13"},
		{"qonto", "oauth", ""},
		{"connector:composio:qonto", "oauth", ""},
		{"connector:composio:calendar", "oauth", ""},
		{"connector:composio:notion", "oauth", ""},
		{"", "key", ""},
		{"connector:composio:googlesheets", "oauth", "2026-11-21"},
		{"googlesheets@1", "oauth", ""},
		{"notion@1", "oauth", "2026-11-21"},
		{"connector:composio:notion", "oauth", ""},
	} {
		got := proofOf(proved, c.entry, c.how)
		if (got == nil) != (c.want == "") || got != nil && got.At != c.want {
			t.Errorf("%s (%s): %+v, want %q", c.entry, c.how, got, c.want)
		}
	}
}

// The shipped proofs against the shipped catalog: each lands on exactly the
// entries its way in was run, and on no other. Qonto, proved by key, is on
// its reviewed key entry and not its sign-in; the apps proved by key with no
// key entry (Stripe's, Shippo's) are on none, only on by_key; each app
// proved by sign-in is on its provider entry.
func TestTheShippedProofsLandOnlyWhereTheyWereRun(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	shipped, err := DecodeCatalog(conf, shippedApps(t))
	if err != nil {
		t.Fatal(err)
	}
	type view struct{ entry, how string }
	views := []view{{StripeEntry, "oauth"}, {QontoEntry, "oauth"}}
	for id, e := range newestRevisions(shipped) {
		views = append(views, view{id, reviewedHow(e)})
	}
	for app := range Proved {
		views = append(views, view{"connector:composio:" + app, "oauth"})
	}
	landed := map[string][]string{}
	for app, p := range Proved {
		for _, v := range views {
			if proofOf(map[string]Proof{app: p}, v.entry, v.how) != nil {
				landed[app] = append(landed[app], v.entry)
			}
		}
	}
	for app := range landed {
		slices.Sort(landed[app])
	}
	want := map[string][]string{
		"qonto":          {"qonto-key@2"},
		"gmail":          {"connector:composio:gmail"},
		"typeform":       {"connector:composio:typeform"},
		"googlecalendar": {"connector:composio:googlecalendar"},
		"googledrive":    {"connector:composio:googledrive"},
		"googlesheets":   {"connector:composio:googlesheets"},
	}
	for app, p := range Proved {
		if !slices.Equal(landed[app], want[app]) {
			t.Errorf("%s (proved by %s) lands on %v, want %v", app, p.How, landed[app], want[app])
		}
	}
}

func TestTheCatalogSaysWhichAppsAreProved(t *testing.T) {
	e, _ := newConnectorFixture(t)
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-16", How: "oauth", Link: "https://example.com/proof"}}
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	seen := 0
	for _, v := range out["entries"].([]any) {
		row := v.(map[string]any)
		proof, has := row["proved"].(map[string]any)
		if row["entry"] == "connector:composio:desk_demo" {
			seen++
			if !has || proof["at"] != "2026-11-16" || proof["how"] != "oauth" || proof["link"] != "https://example.com/proof" {
				t.Errorf("desk_demo isn't proved: %v", row)
			}
		} else if _, there := row["proved"]; there {
			t.Errorf("%v is proved and shouldn't be", row["entry"])
		}
	}
	if seen != 1 {
		t.Fatalf("desk_demo isn't in the catalog: %v", out)
	}
	// Every proof, by app, beside by_key: an app proved by key with no key
	// entry is proved there.
	if all, _ := out["proved"].(map[string]any); out["by_key"] != true || len(all) != 1 || all["deskdemo"] == nil {
		t.Errorf("proved by app: %v", out["proved"])
	}
}

// kiff.dev's apps list reads the proofs with no account: the same map the
// catalog carries, with no reviewed entries, and nothing else.
func TestTheProofsAreReadWithNoAccount(t *testing.T) {
	w := httptest.NewRecorder()
	(&Server{}).Handler().ServeHTTP(w, httptest.NewRequest("GET", "/v1/proved", nil))
	if w.Code != 200 || w.Header().Get("Content-Type") != "application/json" || w.Header().Get("Cache-Control") != "public, max-age=600" {
		t.Fatalf("%d %q %q", w.Code, w.Header().Get("Content-Type"), w.Header().Get("Cache-Control"))
	}
	var got map[string]map[string]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil || len(got) != 1 || got["proved"] == nil {
		t.Fatalf("%v %s", err, w.Body)
	}
	if len(got["proved"]) != len(Proved) {
		t.Fatalf("%d apps, want %d", len(got["proved"]), len(Proved))
	}
	for name, p := range Proved {
		g := got["proved"][name]
		if g["at"] != p.At || g["how"] != p.How || (p.Link != "" && g["link"] != p.Link) {
			t.Errorf("%s: %v, want %+v", name, g, p)
		}
		for k := range g {
			if k != "at" && k != "how" && k != "link" {
				t.Errorf("%s carries %q to a reader with no account", name, k)
			}
		}
	}
	// It's a read: anything else is refused.
	w = httptest.NewRecorder()
	(&Server{}).Handler().ServeHTTP(w, httptest.NewRequest("POST", "/v1/proved", nil))
	if w.Code != 405 {
		t.Fatalf("POST: %d", w.Code)
	}
}

// The route reads the catalog's own source: a server's proofs when it has
// them, so the catalog and the public list never disagree.
func TestTheProofsReadWithNoAccountAreTheCatalogs(t *testing.T) {
	s := &Server{Proved: map[string]Proof{"gmail": {At: "2026-11-13", How: "oauth", Entries: []string{"gmail-sheet"}, Link: "https://example.com/proof"}}}
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/v1/proved", nil))
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != `{"proved":{"gmail":{"at":"2026-11-13","how":"oauth","link":"https://example.com/proof"}}}` {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}
