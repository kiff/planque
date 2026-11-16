package planque

import (
	"strings"
	"testing"
)

func TestTheShippedProofsRead(t *testing.T) {
	p, err := DecodeProved(provedJSON)
	if err != nil || len(p) == 0 {
		t.Fatalf("proved.json: %v %v", p, err)
	}
	for name, proof := range p {
		if !strings.HasPrefix(proof.Link, "https://") {
			t.Errorf("%s has no public link: %+v", name, proof)
		}
	}
}

func TestProofsRefuseWhatTheyCantRead(t *testing.T) {
	for _, c := range []struct{ name, raw, want string }{
		{"not an object", `[]`, "object of apps"},
		{"two objects", `{} {}`, "object of apps"},
		{"null", `null`, "object of apps"},
		{"an unknown field", `{"gmail":{"at":"2026-11-13","by":"me"}}`, "object of apps"},
		{"an unfolded name", `{"Google Calendar":{"at":"2026-11-16"}}`, "folded name"},
		{"a dash in the name", `{"google-calendar":{"at":"2026-11-16"}}`, "folded name"},
		{"no day", `{"gmail":{}}`, "isn't a day"},
		{"a day with a time", `{"gmail":{"at":"2026-11-13T10:00:00Z"}}`, "isn't a day"},
		{"an http link", `{"gmail":{"at":"2026-11-13","link":"http://example.com/a"}}`, "https"},
		{"a link with a user", `{"gmail":{"at":"2026-11-13","link":"https://me@example.com/a"}}`, "https"},
		{"a relative link", `{"gmail":{"at":"2026-11-13","link":"/a"}}`, "https"},
	} {
		if p, err := DecodeProved([]byte(c.raw)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v %v, want %q", c.name, p, err, c.want)
		}
	}
	if p, err := DecodeProved([]byte(`{"gmail":{"at":"2026-11-13"}}`)); err != nil || p["gmail"].At != "2026-11-13" {
		t.Errorf("a proof with no link: %v %v", p, err)
	}
}

func TestAProofMatchesTheAppHoweverItsNamed(t *testing.T) {
	proved := map[string]Proof{"googlecalendar": {At: "2026-11-16"}, "shippo": {At: "2026-11-13"}}
	for _, c := range []struct{ entry, vendor, want string }{
		{"connector:composio:googlecalendar", "", "2026-11-16"},
		{"connector:composio:GOOGLE_CALENDAR", "", "2026-11-16"},
		{"shippo@1", "", "2026-11-13"},
		{"shipping", "Shippo", "2026-11-13"},
		{"connector:composio:calendar", "Google Calendar", "2026-11-16"},
		{"connector:composio:notion", "Notion", ""},
		{"", "", ""},
	} {
		got := proofOf(proved, c.entry, c.vendor)
		if (got == nil) != (c.want == "") || got != nil && got.At != c.want {
			t.Errorf("%s / %s: %+v, want %q", c.entry, c.vendor, got, c.want)
		}
	}
}

func TestTheCatalogSaysWhichAppsAreProved(t *testing.T) {
	e, _ := newConnectorFixture(t)
	e.s.Proved = map[string]Proof{"deskdemo": {At: "2026-11-16", Link: "https://example.com/proof"}}
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
			if !has || proof["at"] != "2026-11-16" || proof["link"] != "https://example.com/proof" {
				t.Errorf("desk_demo isn't proved: %v", row)
			}
		} else if _, there := row["proved"]; there {
			t.Errorf("%v is proved and shouldn't be", row["entry"])
		}
	}
	if seen != 1 {
		t.Fatalf("desk_demo isn't in the catalog: %v", out)
	}
}
