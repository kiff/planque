package planque

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// qonto-key@2 adds the quotes and supplier invoices as reads, and one write:
// an invoice drafted from an accepted quote, always a draft (status fixed),
// whose lines can't carry a discount. Nothing on it can finalise, send, pay
// or delete. @1 stays as it was, for the grants already given.
func TestShippedQontoKey2ReadsAndDraftsOnly(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi-2.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := entries["qonto-key@2"]
	if !ok {
		t.Fatal("qonto-key@2 isn't in the shipped catalog")
	}
	if _, ok := entries["qonto-key@1"]; !ok {
		t.Fatal("qonto-key@1 left the catalog: its grants still name it")
	}
	sum := sha256.Sum256(description)
	if entry.DescriptionSHA256 != hex.EncodeToString(sum[:]) || !strings.HasSuffix(entry.API.URL, "/deploy/specs/qonto-invoicing-openapi-2.json") {
		t.Fatal("deploy/catalog.json doesn't pin deploy/specs/qonto-invoicing-openapi-2.json as it is: update the digest with the description")
	}
	if entry.API.Base != "https://thirdparty.qonto.com/v2" || entry.API.Auth.In != "header" || entry.API.Auth.Name != "Authorization" || entry.API.Auth.Prefix != "" {
		t.Fatalf("Qonto's key goes raw in Authorization, on its v2 base: %+v", entry.API)
	}

	var writes []ToolRequest
	for _, tool := range entry.Tools {
		if !tool.ReadOnly {
			writes = append(writes, tool)
		}
	}
	if len(entry.Tools) != 8 || len(writes) != 1 {
		t.Fatalf("want 7 reads and 1 write, got %d tools with writes %+v", len(entry.Tools), writes)
	}
	draft := writes[0]
	if draft.Tool != "create_client_invoice" || string(draft.Fixed["status"]) != `"draft"` || len(draft.Fixed) != 1 {
		t.Fatalf("the one write must be create_client_invoice with status fixed to draft: %+v", draft)
	}
	if draft.ResourceArg != "quote_id" || !draft.Arguments["quote_id"].Required {
		t.Fatal("every draft names the quote it comes from, for the Card and the record")
	}
	for name, a := range draft.Arguments {
		if name == "status" || name == "discount" || strings.Contains(a.Path, "discount") || a.Path == "status" {
			t.Fatalf("the taf can set %s (%s)", name, a.Path)
		}
	}
	if strings.Contains(string(draft.Arguments["items"].Schema), "discount") {
		t.Fatal("an invoice line can carry a discount")
	}

	// Every operation the entry offers is one Qonto's reviewed description has, and
	// only the draft isn't a read.
	snap, err := openapi.ParseOperations(description, openapi.Auth{In: "header", Name: "Authorization"}, entry.API.Operations)
	if err != nil {
		t.Fatal(err)
	}
	if msg := apiToolsCheck(snap, entry.Tools); msg != "" {
		t.Fatal(msg)
	}
	for _, op := range snap.Operations {
		if !op.ReadOnly && op.Name != "create_client_invoice" {
			t.Fatalf("an operation that isn't a read: %s %s", op.Method, op.Path)
		}
	}
	for _, want := range []string{"list_quotes", "get_quote", "list_supplier_invoices"} {
		if !slices.ContainsFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == want && o.ReadOnly }) {
			t.Fatalf("%s isn't a read in the description", want)
		}
	}
	// A draft the taf could send as it builds it: Qonto's own body, with the status fixed.
	line := map[string]any{"title": "Pose", "quantity": "2", "vat_rate": "0.2", "unit_price": map[string]any{"value": "120.00", "currency": "EUR"}}
	body, err := draft.vendorArgs(map[string]any{"quote_id": "q", "client_id": "c", "issue_date": "2026-11-13", "due_date": "2026-12-13",
		"currency": "EUR", "iban": "FR76", "items": []any{line}})
	var sent map[string]any
	if err != nil || json.Unmarshal(body, &sent) != nil || sent["status"] != "draft" || sent["payment_methods"].(map[string]any)["iban"] != "FR76" {
		t.Fatalf("draft body %s, %v", body, err)
	}
}
