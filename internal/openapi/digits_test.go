package openapi

import (
	"encoding/json"
	"os"
	"testing"
)

// Moneybird's ids, as its description types them: a string or an integer of
// digits alone. Its ids are 18 digits, past 2^53 - 1, so an agent that sent
// one as a JSON number had the call refused (KIFF can't sign a permit over
// such a number). The agent is offered them as text only, in the path and in
// a body alike, and a "null" the type allows stays.
const moneybirdShape = `{
 "openapi": "3.1.0",
 "info": {"title": "Books"},
 "servers": [{"url": "https://books.example/api/v2"}],
 "security": [{"key": []}],
 "components": {
  "securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "X-API-Key"}},
  "schemas": {
   "identifier": {"description": "A unique record identifier", "pattern": "^\\d+$", "type": ["string", "integer"]},
   "administration_id": {"allOf": [{"$ref": "#/components/schemas/identifier", "description": "An administration"}]},
   "maybe_contact": {"pattern": "^[0-9]+$", "type": ["string", "integer", "null"]}
  },
  "parameters": {"administration_id": {"name": "administration_id", "in": "path", "required": true, "schema": {"$ref": "#/components/schemas/administration_id"}}}
 },
 "paths": {
  "/{administration_id}/sales_invoices": {
   "parameters": [{"$ref": "#/components/parameters/administration_id"}],
   "get": {"operationId": "list_invoices", "parameters": [{"name": "page", "in": "query", "schema": {"type": "integer"}}]},
   "post": {"operationId": "create_invoice", "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object",
     "properties": {"contact_id": {"$ref": "#/components/schemas/maybe_contact"}, "total": {"type": ["string", "integer"]}, "lines": {"type": "array", "items": {"type": "object", "properties": {"ledger_id": {"$ref": "#/components/schemas/identifier"}}}}}}}}}}
  }
 }
}`

func TestADigitsOnlyIdIsOfferedAsTextOnly(t *testing.T) {
	s := parse(t, moneybirdShape, key)
	var list struct {
		Properties map[string]map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(op(t, s, "list_invoices").InputSchema, &list); err != nil {
		t.Fatal(err)
	}
	admin, _ := json.Marshal(list.Properties["administration_id"])
	if got := string(admin); !contains(got, `"type":"string"`) || contains(got, `"integer"`) {
		t.Errorf("administration_id is %s, want a string only", got)
	}
	// An integer that isn't an id is left alone.
	if page := list.Properties["page"]["type"]; page != "integer" {
		t.Errorf("page is %v, want integer", page)
	}
	body, _ := json.Marshal(op(t, s, "create_invoice").InputSchema)
	var create struct {
		Properties map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(op(t, s, "create_invoice").InputSchema, &create); err != nil {
		t.Fatal(err)
	}
	contact, _ := json.Marshal(create.Properties["contact_id"])
	if got := string(contact); !contains(got, `"type":["string","null"]`) {
		t.Errorf("contact_id is %s, want string or null (%s)", got, body)
	}
	// A string-or-integer with no digits pattern isn't an id: left as it is.
	total, _ := json.Marshal(create.Properties["total"])
	if got := string(total); !contains(got, `"integer"`) {
		t.Errorf("total is %s, want it left as string or integer", got)
	}
	lines, _ := json.Marshal(create.Properties["lines"])
	if got := string(lines); contains(got, `"integer"`) {
		t.Errorf("a line's ledger_id is %s, want a string only", got)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// planque's own Moneybird description (deploy/specs): reading the invoices takes the administration's id as text.
func TestMoneybirdsAdministrationIdIsText(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/specs/moneybird-key-openapi.json")
	if err != nil {
		t.Skip("no Moneybird description here:", err)
	}
	s, err := Parse(raw, Auth{In: "header", Name: "Authorization", Prefix: "Bearer "})
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, o := range s.Operations {
		var in struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		if json.Unmarshal(o.InputSchema, &in) != nil {
			continue
		}
		if admin, ok := in.Properties["administration_id"]; ok {
			checked++
			if contains(string(admin), `"integer"`) {
				t.Errorf("%s: administration_id is %s, want text only", o.Name, admin)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no operation takes administration_id")
	}
}
