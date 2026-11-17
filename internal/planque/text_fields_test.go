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
	if r := markTextFields(snap, tools); r != nil {
		t.Fatalf("refused: %+v", r)
	}
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
		// On #267: a body in parts, a list of strings, and a value only under a text-named parent.
		`{"properties": {"content": {"type": "array", "items": {"type": "object", "properties": {"type": {"type": "string"}, "value": {"type": "string"}}}},
		  "meta": {"properties": {"value": {"type": "string"}}}}}`: {"content[].value"},
		`{"properties": {"notes": {"type": "array", "items": {"type": "string"}}}}`: {"notes[]"},
		// Words of a name: camelCase, "_" and "-", and the words of a post, an issue or a ticket.
		`{"properties": {"htmlBody": {"type": "string"}, "body_html": {"type": "string"}, "reply-text": {"type": "string"}, "internal_note": {"type": "string"}}}`: {"body_html", "htmlBody", "internal_note", "reply-text"},
		`{"properties": {"commentary": {"type": "string"}, "visibility": {"type": "string", "enum": ["PUBLIC"]}}}`:                                                 {"commentary"},
		`{"properties": {"title": {"type": "string"}, "description": {"type": "string"}, "labels": {"type": "array", "items": {"type": "string"}}}}`:               {"description", "title"},
		`{"properties": {"to": {"type": "string"}, "Body": {"type": "string"}}}`:                                                                                   {"Body"},
		// A write that reaches someone with no text-named argument: its free strings are what it says, never an id.
		`{"properties": {"to": {"type": "string"}, "payload": {"type": "string"}, "template_id": {"type": "string"}, "locale": {"type": "string"}}}`: {"payload"},
		// One that reaches no one and says nothing: a status and an id.
		`{"properties": {"ticket_id": {"type": "string"}, "status": {"type": "string", "enum": ["open", "closed"]}}}`: nil,
	} {
		if got, ok := textFields([]byte(schema)); !ok || !slices.Equal(got, want) {
			t.Errorf("%s: %v (%v), want %v", schema, got, ok, want)
		}
	}
	// More free strings than disclose takes, on a write that reaches someone: it can't be told.
	if got, ok := textFields([]byte(`{"properties": {"phone": {"type": "string"}, "a": {"type": "string"}, "b": {"type": "string"}, "c": {"type": "string"}, "d": {"type": "string"}, "e": {"type": "string"}}}`)); ok {
		t.Errorf("five free strings to a phone: %v", got)
	}
	// At most four, the longest first.
	got, ok := textFields([]byte(`{"properties": {"email": {"type": "string", "format": "email"}, "a": {"type": "string", "maxLength": 10}, "b": {"type": "string", "maxLength": 5000}, "c": {"type": "string"}}}`))
	if !ok || !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("three free strings to an email: %v %v", got, ok)
	}
}
