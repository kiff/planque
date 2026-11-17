package planque

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// Every key entry in the shipped catalog (deploy/catalog.json), the Qonto ones
// included: it pins the description kept in deploy/specs as it is, what it
// calls a read is a read in that description and what it calls a write is held
// for the owner's approval, and a grant from it asks for the key and offers
// exactly its tools. Entries are written by cmd/catalogen from a vendor's own
// description, and this is what keeps a batch of them honest.
func TestEveryShippedKeyEntryPinsItsDescriptionAndHoldsItsWrites(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	for id, entry := range entries {
		if entry.Kind != KindAPIKey {
			continue
		}
		t.Run(id, func(t *testing.T) {
			description, err := os.ReadFile("../../deploy/specs/" + path.Base(entry.API.URL))
			if err != nil {
				t.Fatalf("the description %s names isn't kept in deploy/specs: %v", entry.API.URL, err)
			}
			if !strings.HasPrefix(entry.API.URL, "https://raw.githubusercontent.com/kiff/planque/main/deploy/specs/") {
				t.Fatalf("a description is read from planque's own deploy/specs: %s", entry.API.URL)
			}
			if sum := sha256.Sum256(description); entry.DescriptionSHA256 != hex.EncodeToString(sum[:]) {
				t.Fatal("the entry doesn't pin its description as it is: regenerate the entry with the description")
			}
			snap, err := openapi.ParseToList(description, entry.API.Auth)
			if err != nil {
				t.Fatal(err)
			}
			described := map[string]openapi.Operation{}
			for _, o := range snap.Operations {
				described[o.Name] = o
			}
			var profile struct {
				Operations []struct {
					Tool   string `json:"tool"`
					Kind   string `json:"kind"`
					Effect string `json:"effect"`
				} `json:"operations"`
			}
			if err := json.Unmarshal(entry.Profile, &profile); err != nil {
				t.Fatal(err)
			}
			kinds := map[string]string{}
			effects := map[string]string{}
			for _, o := range profile.Operations {
				kinds[o.Tool], effects[o.Tool] = o.Kind, o.Effect
			}
			for _, tool := range entry.Tools {
				op, ok := described[tool.Tool]
				if !ok {
					t.Fatalf("%s isn't an operation of its description", tool.Tool)
				}
				if op.ReadOnly != tool.ReadOnly {
					t.Fatalf("%s: the entry says read_only=%v, the description says %v (%s %s)", tool.Name, tool.ReadOnly, op.ReadOnly, op.Method, op.Path)
				}
				if tool.ReadOnly != (kinds[tool.Name] == "read") {
					t.Fatalf("%s: the profile's kind %q doesn't match its tool", tool.Name, kinds[tool.Name])
				}
				if !tool.ReadOnly && effects[tool.Name] == "" {
					t.Fatalf("%s: a write that no effect holds for the owner", tool.Name)
				}
			}

			// A grant from the entry asks for the key and offers exactly its tools.
			e := newInProcessEnv(t)
			e.enrol()
			e.s.CatalogEntries = entries
			e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != entry.API.URL {
					t.Errorf("fetched %s", r.URL)
				}
				return profileResponse(200, string(description)), nil
			})}
			ask := map[string]string{"agent": "agt_taf"}
			switch {
			case len(entry.API.Bases) > 0:
				ask["base"] = entry.API.Bases[len(entry.API.Bases)-1] // the last region, as a customer may pick
			case accountTemplate(entry.API.Base):
				ask["account"] = "acme"
			}
			st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/"+id+"/grants", ask, nil, "k_setup")
			if st != 201 || out["needs_key"] != true {
				t.Fatalf("key grant: %d %v", st, out)
			}
			g, err := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
			if err != nil {
				t.Fatal(err)
			}
			if len(g.Request.Tools) != len(entry.Tools) {
				t.Fatalf("tools: %+v", g.Request.Tools)
			}
		})
	}
}
