package openapi

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const desc = `{
 "openapi": "3.1.0",
 "info": {"title": "Shop API"},
 "servers": [{"url": "https://{region}.shop.example/v1", "variables": {"region": {"default": "eu"}}}],
 "security": [{"key": []}],
 "components": {
  "securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "X-API-Key"}, "bearer": {"type": "http", "scheme": "bearer"}},
  "schemas": {
   "Order": {"type": "object", "properties": {"id": {"type": "string"}, "lines": {"type": "array", "items": {"$ref": "#/components/schemas/Order"}}}},
   "Refund": {"type": "object", "required": ["amount"], "properties": {"amount": {"type": "integer"}, "reason": {"type": "string"}}}
  },
  "parameters": {"Order": {"name": "order", "in": "path", "required": true, "schema": {"type": "string"}}}
 },
 "paths": {
  "/orders/{order}": {
   "parameters": [{"$ref": "#/components/parameters/Order"}],
   "get": {"operationId": "getOrder", "summary": "Read an order", "parameters": [{"name": "expand", "in": "query", "schema": {"type": "string"}}],
           "responses": {"200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Order"}}}}}},
   "delete": {"summary": "Cancel an order"}
  },
  "/orders/{order}/refunds": {
   "parameters": [{"$ref": "#/components/parameters/Order"}],
   "post": {"operationId": "create-refund", "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Refund"}}}}}
  },
  "/notes": {
   "put": {"operationId": "putNotes", "requestBody": {"content": {"application/json": {"schema": {"type": "array", "items": {"type": "string"}}}}}}
  },
  "/upload": {
   "post": {"operationId": "upload", "requestBody": {"content": {"multipart/form-data": {"schema": {"type": "object"}}}}}
  },
  "/signed": {
   "get": {"operationId": "signed", "parameters": [{"name": "X-Signature", "in": "header", "required": true}]}
  }
 }
}`

var key = Auth{In: "header", Name: "x-api-key"}

func parse(t *testing.T, raw string, a Auth) *Snapshot {
	t.Helper()
	s, err := Parse([]byte(raw), a)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func op(t *testing.T, s *Snapshot, name string) Operation {
	t.Helper()
	for _, o := range s.Operations {
		if o.Name == name {
			return o
		}
	}
	t.Fatalf("no operation %s in %+v", name, s.Operations)
	return Operation{}
}

func TestADescriptionReadsIntoTools(t *testing.T) {
	s := parse(t, desc, key)
	if s.Title != "Shop API" || len(s.Servers) != 1 || s.Servers[0] != "https://eu.shop.example/v1" || len(s.Digest) != 64 {
		t.Fatalf("%+v", s)
	}
	get := op(t, s, "getOrder")
	if get.Method != "GET" || !get.ReadOnly || get.Path != "/orders/{order}" || get.Summary != "Read an order" || get.Body != "" {
		t.Fatalf("%+v", get)
	}
	var schema struct {
		Properties map[string]any `json:"properties"`
		Required   []string       `json:"required"`
	}
	_ = json.Unmarshal(get.InputSchema, &schema)
	if schema.Properties["order"] == nil || schema.Properties["expand"] == nil || len(schema.Required) != 1 || schema.Required[0] != "order" {
		t.Fatalf("%s", get.InputSchema)
	}
	// An operation with no id is named by its method and path; a write isn't read-only.
	del := op(t, s, "delete_orders_order")
	if del.ReadOnly || del.Method != "DELETE" {
		t.Fatalf("%+v", del)
	}
	// An id that isn't a tool name is made one; an object body's fields are arguments.
	refund := op(t, s, "create_refund")
	_ = json.Unmarshal(refund.InputSchema, &schema)
	if refund.Body != "fields" || schema.Properties["amount"] == nil || strings.Join(schema.Required, ",") != "amount,order" {
		t.Fatalf("%+v %s", refund, refund.InputSchema)
	}
	// A body that isn't an object is the body argument.
	if notes := op(t, s, "putNotes"); notes.Body != "body" {
		t.Fatalf("%+v", notes)
	}
	// Left out, with why: a body that isn't JSON, a required header.
	why := map[string]string{}
	for _, sk := range s.Skipped {
		why[sk.Name] = sk.Why
	}
	if !strings.Contains(why["upload"], "isn't JSON") || !strings.Contains(why["signed"], "header parameter") {
		t.Fatalf("%+v", s.Skipped)
	}
	// A recursive schema resolves without looping.
	if !strings.Contains(string(get.InputSchema), `"order"`) {
		t.Fatal(string(get.InputSchema))
	}
}

func TestTheKeyMustBeOneOfItsSchemes(t *testing.T) {
	bearer := Auth{In: "header", Name: "Authorization", Prefix: "Bearer "}
	// Every operation here takes the key scheme only: a bearer key is secured otherwise.
	if _, err := Parse([]byte(desc), bearer); err == nil || !strings.Contains(err.Error(), "secured otherwise") {
		t.Fatalf("bearer where only the key is taken: %v", err)
	}
	either := strings.Replace(desc, `"security": [{"key": []}]`, `"security": [{"key": []}, {"bearer": []}]`, 1)
	if _, err := Parse([]byte(either), bearer); err != nil {
		t.Fatalf("a bearer scheme the operations take: %v", err)
	}
	for _, a := range []Auth{{In: "query", Name: "x-api-key"}, {In: "header", Name: "X-Other"}, {In: "cookie", Name: "k"}, {In: "header"}} {
		var r *Refused
		if _, err := Parse([]byte(desc), a); !errors.As(err, &r) {
			t.Fatalf("%+v: %v", a, err)
		}
	}
}

func TestWhatIsntReadIsRefusedWithTheReason(t *testing.T) {
	for name, c := range map[string]struct{ from, to, why string }{
		"yaml":        {`{`, "openapi: 3.0.0\n", "JSON descriptions only"},
		"swagger 2":   {`"openapi": "3.1.0"`, `"openapi": "2.0"`, "isn't OpenAPI 3.0 or 3.1"},
		"external":    {`"#/components/schemas/Refund"`, `"other.json#/Refund"`, "references another document"},
		"op servers":  {`"operationId": "putNotes",`, `"operationId": "putNotes", "servers": [],`, "sets its own servers"},
		"callbacks":   {`"operationId": "putNotes",`, `"operationId": "putNotes", "callbacks": {},`, "has callbacks"},
		"links":       {`"200": {"content"`, `"200": {"links": {}, "content"`, "has links"},
		"other auth":  {`"operationId": "putNotes",`, `"operationId": "putNotes", "security": [{"oauth": []}],`, "secured otherwise"},
		"bad pointer": {`"#/components/schemas/Refund"`, `"#/components/schemas/Nope"`, "names nothing"},
	} {
		raw := strings.Replace(desc, c.from, c.to, 1)
		if name == "yaml" {
			raw = "openapi: 3.0.0\npaths: {}\n"
		}
		_, err := Parse([]byte(raw), key)
		var r *Refused
		if !errors.As(err, &r) || !strings.Contains(err.Error(), c.why) {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := Parse(make([]byte, MaxBytes+1), key); err == nil || !strings.Contains(err.Error(), "16 MB") {
		t.Fatal(err)
	}
}

func TestACallFillsItsRequest(t *testing.T) {
	s := parse(t, desc, key)
	method, target, body, err := op(t, s, "create_refund").Request("https://eu.shop.example/v1/", map[string]any{
		"order": "o/1 x", "amount": json.Number("4000"), "reason": "damaged"})
	if err != nil || method != "POST" || target != "https://eu.shop.example/v1/orders/o%2F1%20x/refunds" {
		t.Fatalf("%s %s %v", method, target, err)
	}
	var got map[string]any
	if json.Unmarshal(body, &got) != nil || got["amount"] != float64(4000) || got["reason"] != "damaged" || got["order"] != nil {
		t.Fatalf("%s", body)
	}
	_, target, body, err = op(t, s, "getOrder").Request("https://eu.shop.example/v1", map[string]any{"order": "o1", "expand": "lines"})
	if err != nil || target != "https://eu.shop.example/v1/orders/o1?expand=lines" || body != nil {
		t.Fatalf("%s %s %v", target, body, err)
	}
	// Missing, extra and misplaced arguments are refused before anything is sent.
	for _, args := range []map[string]any{{}, {"order": "o1", "nope": "x"}, {"order": map[string]any{"a": 1}}} {
		if _, _, _, err := op(t, s, "getOrder").Request("https://eu.shop.example/v1", args); !errors.Is(err, ErrArguments) {
			t.Errorf("%v: %v", args, err)
		}
	}
	// A path parameter of . or .. would make another endpoint once the path is
	// normalized: refused, as is an empty one.
	for _, order := range []string{"..", ".", ""} {
		if _, _, _, err := op(t, s, "getOrder").Request("https://eu.shop.example/v1", map[string]any{"order": order}); !errors.Is(err, ErrArguments) {
			t.Errorf("%q: %v", order, err)
		}
	}
	if _, _, _, err := op(t, s, "getOrder").Request("ftp://eu.shop.example/v1", map[string]any{"order": "o1"}); !errors.Is(err, ErrArguments) {
		t.Fatal(err)
	}
	_, _, body, err = op(t, s, "putNotes").Request("https://eu.shop.example", map[string]any{"body": []any{"a", "b"}})
	if err != nil || string(body) != `["a","b"]` {
		t.Fatalf("%s %v", body, err)
	}
}
