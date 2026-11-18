package planque

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"testing"
)

// The reviewed catalog that ships with the function (deploy/catalog.json)
// decodes, pins the description it names (deploy/specs), and gives Qonto by
// key as reads only: nothing on it can create, finalise or send.
func TestShippedQontoCatalogIsReadOnlyAndPinsItsDescription(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := entries["qonto-key@1"]
	if !ok {
		t.Fatal("qonto-key@1 isn't in the shipped catalog")
	}
	sum := sha256.Sum256(description)
	if entry.DescriptionSHA256 != hex.EncodeToString(sum[:]) {
		t.Fatal("deploy/catalog.json doesn't pin deploy/specs/qonto-invoicing-openapi.json as it is: update the digest with the description")
	}
	if entry.API.Base != "https://thirdparty.qonto.com/v2" || entry.API.Auth.In != "header" || entry.API.Auth.Name != "Authorization" || entry.API.Auth.Prefix != "" {
		t.Fatalf("Qonto's key goes raw in Authorization, on its v2 base: %+v", entry.API)
	}
	for _, tool := range entry.Tools {
		if !tool.ReadOnly || tool.Tool == "create_client_invoice" {
			t.Fatalf("a Qonto tool that isn't a read: %+v", tool)
		}
	}

	// A grant from the entry: the key is asked for, and only reads are offered.
	e := newInProcessEnv(t)
	e.enrol()
	e.s.CatalogEntries = entries
	e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != entry.API.URL {
			t.Errorf("fetched %s", r.URL)
		}
		return profileResponse(200, string(description)), nil
	})}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 || out["needs_key"] != true {
		t.Fatalf("Qonto key grant: %d %v", st, out)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Request.Tools) != 4 {
		t.Fatalf("tools: %+v", g.Request.Tools)
	}
	for _, tool := range g.Request.Tools {
		if !tool.ReadOnly {
			t.Fatalf("granted a write: %+v", tool)
		}
	}
}
