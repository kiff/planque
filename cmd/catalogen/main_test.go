package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// A small description of a desk: a read, a reply that can copy outside
// addresses, a note that can be made public.
const deskSpec = `{
 "openapi": "3.0.3", "info": {"title": "Desk", "version": "1"},
 "servers": [{"url": "https://api.desk.example/v1"}],
 "components": {"securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "Authorization"}}},
 "security": [{"key": []}],
 "paths": {
  "/tickets/{id}": {"get": {"operationId": "getTicket", "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "ok"}}}},
  "/tickets/{id}/reply": {"post": {"operationId": "reply", "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
    "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["body"],
      "properties": {"body": {"type": "string"}, "bcc_emails": {"type": "array", "items": {"type": "string"}}, "cc": {"type": "string"}}}}}},
    "responses": {"200": {"description": "ok"}}}},
  "/tickets/{id}/notes": {"post": {"operationId": "note", "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
    "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["body"],
      "properties": {"body": {"type": "string"}, "kind": {"type": "string"}}}}}},
    "responses": {"200": {"description": "ok"}}}}
 }}`

func deskPlan() Plan {
	return Plan{ID: "desk-key@1", Vendor: "Desk", Title: "Desk: read a ticket, reply", Spec: "desk.json", Base: "https://api.desk.example/v1", Prefix: "desk",
		Auth:  openapi.Auth{In: "header", Name: "Authorization"},
		Reads: []Op{{Operation: "getTicket", Name: "read_ticket", Title: "Read one ticket"}}}
}

// A write that can copy an address is refused until the plan pins it, so a
// held reply can't copy anyone the owner didn't see.
func TestAWriteThatCanCopyAnAddressIsRefused(t *testing.T) {
	plan := deskPlan()
	plan.Writes = []Op{{Operation: "reply", Name: "reply", Title: "Reply to a ticket", Noun: "ticket", Effect: "sends"}}
	_, _, err := Generate(plan, []byte(deskSpec))
	if err == nil || !strings.Contains(err.Error(), "address list") {
		t.Fatalf("a reply with bcc_emails and cc: %v", err)
	}
	// Pinned, it goes: those two are left out of its arguments.
	plan.Writes[0].Fixed = map[string]json.RawMessage{"bcc_emails": json.RawMessage(`[]`), "cc": json.RawMessage(`""`)}
	entry, _, err := Generate(plan, []byte(deskSpec))
	if err != nil {
		t.Fatalf("pinned: %v", err)
	}
	raw, _ := json.Marshal(entry)
	var got struct {
		Tools []struct {
			Name  string                    `json:"name"`
			Args  map[string]map[string]any `json:"arguments"`
			Fixed map[string]any            `json:"fixed_arguments"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	for _, tool := range got.Tools {
		if tool.Name != "desk_reply" {
			continue
		}
		if _, ok := tool.Args["bcc_emails"]; ok || tool.Args["body"] == nil || tool.Args["id"] == nil || len(tool.Args) != 2 {
			t.Fatalf("the pinned reply's arguments: %v", tool.Args)
		}
		if len(tool.Fixed) != 2 {
			t.Fatalf("its pinned values: %v", tool.Fixed)
		}
		return
	}
	t.Fatal("no reply tool")
}

// What the entry calls a read or a write comes from the description.
func TestAPlanCannotCallAWriteARead(t *testing.T) {
	plan := deskPlan()
	plan.Reads = append(plan.Reads, Op{Operation: "note", Name: "note", Title: "A note"})
	if _, _, err := Generate(plan, []byte(deskSpec)); err == nil || !strings.Contains(err.Error(), "plan lists it as a read") {
		t.Fatalf("a note as a read: %v", err)
	}
}

// A generator never marks its own entry proven.
func TestAGeneratedEntryIsNeverProven(t *testing.T) {
	entry, _, err := Generate(deskPlan(), []byte(deskSpec))
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(entry["profile"])
	var p struct {
		TestProven *bool `json:"test_proven"`
	}
	if err := json.Unmarshal(raw, &p); err != nil || p.TestProven == nil || *p.TestProven {
		t.Fatalf("test_proven: %s", raw)
	}
}
