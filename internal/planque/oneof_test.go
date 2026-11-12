package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const ibanDesc = `{
 "openapi": "3.0.3",
 "info": {"title": "Books API"},
 "security": [{"apiKey": []}],
 "components": {"securitySchemes": {"apiKey": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}},
 "paths": {
  "/organization": {"get": {"operationId": "get_organization", "summary": "Read the business"}},
  "/invoices": {"post": {"operationId": "create_invoice", "summary": "Create an invoice",
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["quote_id", "payment_methods"],
    "properties": {"quote_id": {"type": "string"}, "status": {"type": "string"}, "payment_methods": {"type": "object", "properties": {"iban": {"type": "string"}}}}}}}}}}
 }
}`

// An argument with one_of is checked against the account's own read before
// each call, with the grant's key: an IBAN the business doesn't have (one
// planted in a quote) is refused, and nothing reaches the vendor's write.
// Written as people write it, the business's own IBAN goes through. A read
// that fails refuses the call too: planque never sends what it couldn't check.
func TestAnArgumentIsCheckedAgainstTheAccountsOwnRead(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	orgDown := false
	v.Server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.mu.Lock()
		defer v.mu.Unlock()
		if r.URL.Path == "/openapi.json" {
			_, _ = io.WriteString(w, ibanDesc)
			return
		}
		c := restCall{Method: r.Method, Path: r.URL.Path, Key: r.Header.Get("X-API-Key")}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &c.Body)
		v.calls = append(v.calls, c)
		switch {
		case c.Key != v.key:
			w.WriteHeader(401)
		case r.URL.Path == "/v1/organization" && orgDown:
			w.WriteHeader(503)
		case r.URL.Path == "/v1/organization":
			_, _ = io.WriteString(w, `{"organization":{"bank_accounts":[{"iban":"FR7630001007941234567890185"},{"iban":"FR7612345000010009876543210"}]}}`)
		default:
			_, _ = io.WriteString(w, `{"client_invoice":{"id":"inv_1","status":"draft"}}`)
		}
	})
	e.enrol()
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, out := e.giveKey(id, out["digest"].(string), v.key); st != 200 {
		t.Fatalf("key: %d %v", st, out)
	}
	tools := []map[string]any{
		{"tool": "get_organization", "name": "organization", "title": "Read the business", "read_only": true},
		{"tool": "create_invoice", "name": "draft_invoice", "title": "Draft an invoice",
			"fixed_arguments": map[string]any{"status": "draft"},
			"arguments": map[string]any{
				"quote_id": map[string]any{"path": "quote_id", "type": "string", "required": true},
				"iban":     map[string]any{"path": "payment_methods.iban", "type": "string", "required": true, "one_of": map[string]any{"tool": "get_organization", "path": "organization.bank_accounts.iban"}},
			}},
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": tools}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	e.kiff.Cards["k_taf"] = map[string]int64{"draft_invoice": 0, "organization": 0}
	ctx := context.Background()
	writes := func() (n int) {
		for _, c := range v.seen() {
			if c.Path == "/v1/invoices" {
				n++
			}
		}
		return n
	}

	res, err := e.kiff.Call(ctx, "k_taf", "draft_invoice", map[string]any{"quote_id": "q_1", "iban": "DE89370400440532013000"})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("%+v %v", res, err)
	}
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "isn't one the account's get_organization lists") || !strings.Contains(text, "Nothing was sent") {
		t.Fatalf("a planted IBAN: %s", text)
	}
	if writes() != 0 {
		t.Fatalf("the write reached the vendor: %+v", v.seen())
	}
	if c := v.seen()[0]; c.Path != "/v1/organization" || c.Method != "GET" || c.Key != v.key {
		t.Fatalf("the check's read: %+v", c)
	}

	res, _ = e.kiff.Call(ctx, "k_taf", "draft_invoice", map[string]any{"quote_id": "q_1", "iban": "fr76 1234 5000 0100 0987 6543 210"})
	if text, isErr := resultText(t, res.Result); isErr || !strings.Contains(text, "inv_1") || writes() != 1 {
		t.Fatalf("the business's own IBAN: %s %+v", text, v.seen())
	}
	if last := v.seen()[len(v.seen())-1]; last.Body["status"] != "draft" || last.Body["payment_methods"].(map[string]any)["iban"] != "fr76 1234 5000 0100 0987 6543 210" {
		t.Fatalf("sent %+v", last.Body)
	}

	v.mu.Lock()
	orgDown = true
	v.mu.Unlock()
	res, _ = e.kiff.Call(ctx, "k_taf", "draft_invoice", map[string]any{"quote_id": "q_2", "iban": "FR7630001007941234567890185"})
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "couldn't be read") || writes() != 1 {
		t.Fatalf("the read down: %s", text)
	}
}

// one_of names a read-only tool of the same request, on a string argument.
func TestOneOfIsChecked(t *testing.T) {
	read := ToolRequest{Tool: "get_organization", Name: "organization", ReadOnly: true}
	write := func(a Argument) ToolRequest {
		return ToolRequest{Tool: "create_invoice", Name: "draft", Fixed: map[string]json.RawMessage{"status": json.RawMessage(`"draft"`)},
			Arguments: map[string]Argument{"iban": a}}
	}
	ok := Argument{Path: "payment_methods.iban", Type: "string", OneOf: &OneOf{Tool: "get_organization", Path: "organization.bank_accounts.iban"}}
	if msg := checkTools("agt", []ToolRequest{read, write(ok)}); msg != "" {
		t.Fatal(msg)
	}
	for name, c := range map[string][]ToolRequest{
		"no such read":      {write(ok)},
		"a write's read":    {{Tool: "get_organization", Name: "organization"}, write(ok)},
		"not a string":      {read, write(Argument{Path: "n", Type: "integer", OneOf: ok.OneOf})},
		"a path with a gap": {read, write(Argument{Path: "iban", Type: "string", OneOf: &OneOf{Tool: "get_organization", Path: "a..b"}})},
	} {
		if msg := checkTools("agt", c); msg == "" {
			t.Errorf("%s: accepted", name)
		}
	}
}
