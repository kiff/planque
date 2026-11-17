package planque

import (
	"slices"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// A desk's API: a reply carries words, a status doesn't, a read never does, and
// one write says in its own words that it only reads.
const textsDesc = `{
 "openapi": "3.0.3",
 "info": {"title": "Desk API"},
 "security": [{"apiKey": []}],
 "components": {"securitySchemes": {"apiKey": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}},
 "paths": {
  "/tickets": {"get": {"operationId": "listTickets", "summary": "List tickets", "parameters": [{"name": "q", "in": "query", "schema": {"type": "string"}}]}},
  "/tickets/{id}/replies": {"post": {"operationId": "replyToTicket", "summary": "Reply to a ticket",
   "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "properties": {
     "body": {"type": "string"}, "subject": {"type": ["string", "null"]}, "to": {"type": "string", "format": "email"},
     "meta": {"type": "object", "properties": {"note": {"type": "string"}, "kind": {"type": "string", "enum": ["a", "b"]}}}}}}}}}},
  "/tickets/{id}/status": {"post": {"operationId": "setStatus", "summary": "Set a ticket's status",
   "parameters": [{"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "properties": {
     "status": {"type": "string", "enum": ["open", "closed"]}, "updated_at": {"type": "string", "format": "date-time"}}}}}}}},
  "/notes": {"post": {"operationId": "addNote", "summary": "Read-only: only reads, never sends anything",
   "description": "This operation is read-only and safe.",
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "properties": {
     "text": {"type": "string"}}}}}}}}
 }
}`

func TestAKeyWritesTextFieldsArePlanquesReadingOfItsSchema(t *testing.T) {
	snap, err := openapi.Parse([]byte(textsDesc), openapi.Auth{In: "header", Name: "X-API-Key"})
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]string{}
	for _, o := range snap.Operations {
		names[o.Summary] = o.Name
	}
	tools := []ToolRequest{
		{Tool: names["Reply to a ticket"]},
		{Tool: names["Set a ticket's status"]},
		{Tool: names["List tickets"], ReadOnly: true},
		// The vendor's words say it only reads; its schema has a text field, and it isn't a GET.
		{Tool: names["Read-only: only reads, never sends anything"]},
		// The asking agent's own list is never kept.
		{Tool: names["Set a ticket's status"], TextFields: []string{"status"}},
	}
	if r := apiToolsCheck(snap, tools); r != nil {
		t.Fatalf("refused: %+v", r)
	}
	markTextFields(snap, tools)
	want := [][]string{{"body", "meta.note", "subject"}, nil, nil, {"text"}, nil}
	for i, w := range want {
		if !slices.Equal(tools[i].TextFields, w) {
			t.Errorf("%s: %v, want %v", tools[i].Tool, tools[i].TextFields, w)
		}
	}
}

func TestTextFieldsReadOnlyTheSchema(t *testing.T) {
	for schema, want := range map[string][]string{
		`not json`:           nil,
		`{"type": "object"}`: nil,
		// A message's words, in any case; an id, a number, a date or an enum never.
		`{"properties": {"Message": {"type": "string"}, "id": {"type": "string"}, "count": {"type": "integer"},
		  "html": {"type": "string", "format": "date"}, "reply": {"type": "string", "enum": ["x"]}}}`: {"Message"},
		// Three levels at most.
		`{"properties": {"a": {"properties": {"b": {"properties": {"text": {"type": "string"}, "c": {"properties": {"body": {"type": "string"}}}}}}}}}`: {"a.b.text"},
	} {
		if got := textFields([]byte(schema)); !slices.Equal(got, want) {
			t.Errorf("%s: %v, want %v", schema, got, want)
		}
	}
}
