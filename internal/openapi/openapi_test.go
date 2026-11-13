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
		"swagger 2":   {`"openapi": "3.1.0"`, `"openapi": "2.0"`, "isn't OpenAPI 3.0 or 3.1"},
		"external":    {`"#/components/schemas/Refund"`, `"other.json#/Refund"`, "references another document"},
		"op servers":  {`"operationId": "putNotes",`, `"operationId": "putNotes", "servers": [],`, "sets its own servers"},
		"callbacks":   {`"operationId": "putNotes",`, `"operationId": "putNotes", "callbacks": {},`, "has callbacks"},
		"links":       {`"200": {"content"`, `"200": {"links": {}, "content"`, "has links"},
		"other auth":  {`"operationId": "putNotes",`, `"operationId": "putNotes", "security": [{"oauth": []}],`, "secured otherwise"},
		"bad pointer": {`"#/components/schemas/Refund"`, `"#/components/schemas/Nope"`, "names nothing"},
	} {
		raw := strings.Replace(desc, c.from, c.to, 1)
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

// HTTP Basic is a place for a key when the description names a basic
// scheme (Zendesk, Freshdesk, Gorgias); its username or password beside the
// key is checked as any argument.
func TestABasicKeyMatchesABasicScheme(t *testing.T) {
	basic := strings.Replace(desc, `"bearer": {"type": "http", "scheme": "bearer"}`, `"basic": {"type": "http", "scheme": "basic"}`, 1)
	basic = strings.Replace(basic, `"security": [{"key": []}]`, `"security": [{"basic": []}]`, 1)
	for _, a := range []Auth{{In: "basic", User: "me@example.com/token"}, {In: "basic", Password: "X"}, {In: "basic"}} {
		if _, err := Parse([]byte(basic), a); err != nil {
			t.Fatalf("%+v: %v", a, err)
		}
	}
	// A basic key where the description only takes its header key.
	if _, err := Parse([]byte(desc), Auth{In: "basic", User: "me"}); err == nil || !strings.Contains(err.Error(), "none of its security schemes") {
		t.Fatalf("basic where only a header key is taken: %v", err)
	}
	for name, a := range map[string]Auth{
		"a name":            {In: "basic", Name: "Authorization"},
		"a prefix":          {In: "basic", Prefix: "Basic "},
		"a colon":           {In: "basic", User: "me:you"},
		"both":              {In: "basic", User: "me", Password: "X"},
		"a newline":         {In: "basic", User: "me\nHost: elsewhere"},
		"too long":          {In: "basic", User: strings.Repeat("a", MaxBasicPart+1)},
		"a header's user":   {In: "header", Name: "X-API-Key", User: "me"},
		"a header, no name": {In: "header"},
		"somewhere else":    {In: "cookie", Name: "k"},
	} {
		var r *Refused
		if _, err := Parse([]byte(basic), a); !errors.As(err, &r) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// A description that names no security scheme at all (GitHub's) takes the
// key where the request says, from the vendor's docs; an operation that
// names one is still refused, as there's none it could be.
func TestADescriptionNamingNoSchemeTakesTheKeyAsAsked(t *testing.T) {
	none := strings.Replace(desc, `"securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "X-API-Key"}, "bearer": {"type": "http", "scheme": "bearer"}},`, ``, 1)
	none = strings.Replace(none, ` "security": [{"key": []}],`+"\n", "", 1)
	if none == desc || strings.Contains(none, "securitySchemes") || strings.Contains(none, `"security"`) {
		t.Fatalf("setup: %s", none)
	}
	s, err := Parse([]byte(none), Auth{In: "header", Name: "Authorization", Prefix: "token "})
	if err != nil || len(s.Operations) == 0 {
		t.Fatalf("%+v %v", s, err)
	}
	named := strings.Replace(none, `"operationId": "putNotes",`, `"operationId": "putNotes", "security": [{"oauth": []}],`, 1)
	if _, err := Parse([]byte(named), Auth{In: "header", Name: "Authorization"}); err == nil || !strings.Contains(err.Error(), "secured otherwise") {
		t.Fatalf("an operation naming a scheme that isn't there: %v", err)
	}
	// The key's place must still be one.
	if _, err := Parse([]byte(none), Auth{In: "header"}); err == nil {
		t.Fatal("no header name")
	}
}

// A body that's a form (Stripe's kind) is the tool's fields, written as a
// form: strings, numbers and lists of them, never nested; JSON is taken
// when the operation takes both.
func TestAFormBodyIsWrittenAsAForm(t *testing.T) {
	form := strings.Replace(desc, `"post": {"operationId": "create-refund", "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Refund"}}}}}`,
		`"post": {"operationId": "create-refund", "requestBody": {"required": true, "content": {"application/x-www-form-urlencoded": {"schema": {"type": "object", "required": ["amount"], "properties": {"amount": {"type": "integer"}, "reason": {"type": "string"}, "expand": {"type": "array", "items": {"type": "string"}}, "metadata": {"type": "object"}}}}}}}`, 1)
	if form == desc {
		t.Fatal("setup: the form body isn't in")
	}
	refund := op(t, parse(t, form, key), "create_refund")
	if refund.Body != "form" || refund.ContentType() != "application/x-www-form-urlencoded" {
		t.Fatalf("%+v", refund)
	}
	method, target, body, err := refund.Request("https://eu.shop.example/v1", map[string]any{"order": "o 1", "amount": json.Number("500"), "reason": "a & b", "expand": []any{"charge", "customer"}})
	if err != nil || method != "POST" || !strings.HasSuffix(target, "/orders/o%201/refunds") || string(body) != "amount=500&expand=charge&expand=customer&reason=a+%26+b" {
		t.Fatalf("%s %s %q %v", method, target, body, err)
	}
	if _, _, _, err := refund.Request("https://eu.shop.example/v1", map[string]any{"order": "o1", "amount": 5, "metadata": map[string]any{"k": "v"}}); !errors.Is(err, ErrArguments) || !strings.Contains(err.Error(), "no nested fields") {
		t.Fatalf("a nested field: %v", err)
	}
	// A JSON operation is still written as JSON.
	if j := op(t, parse(t, desc, key), "create_refund"); j.Body != "fields" || j.ContentType() != "application/json" {
		t.Fatalf("%+v", j)
	}
	// Both offered: JSON is taken.
	both := strings.Replace(form, `"application/x-www-form-urlencoded": {"schema"`, `"application/json": {"schema": {"$ref": "#/components/schemas/Refund"}}, "application/x-www-form-urlencoded": {"schema"`, 1)
	if b := op(t, parse(t, both, key), "create_refund"); b.Body != "fields" {
		t.Fatalf("both: %+v", b)
	}
	// A form that isn't fields of its own is left out, with why.
	whole := strings.Replace(form, `{"type": "object", "required": ["amount"], "properties": {"amount": {"type": "integer"}, "reason": {"type": "string"}, "expand": {"type": "array", "items": {"type": "string"}}, "metadata": {"type": "object"}}}`, `{"type": "string"}`, 1)
	s := parse(t, whole, key)
	for _, sk := range s.Skipped {
		if sk.Name == "create_refund" && strings.Contains(sk.Why, "form body isn't fields") {
			return
		}
	}
	t.Fatalf("a form that's one string: %+v", s.Skipped)
}

// A server at an account's own address (a variable with no default) no
// longer refuses the description: it's kept as an account server, for the
// request to name the account's as its base.
func TestAServerAtAnAccountsAddressIsKeptForTheBase(t *testing.T) {
	acct := strings.Replace(desc, `"servers": [{"url": "https://{region}.shop.example/v1", "variables": {"region": {"default": "eu"}}}]`,
		`"servers": [{"url": "https://{subdomain}.shop.example/v1", "variables": {"subdomain": {"description": "your account"}}}]`, 1)
	if acct == desc {
		t.Fatal("setup")
	}
	s := parse(t, acct, key)
	if len(s.Servers) != 0 || len(s.Accounts) != 1 || s.Accounts[0] != "https://{subdomain}.shop.example/v1" || len(s.Operations) == 0 {
		t.Fatalf("%+v", s)
	}
	// A variable with a default is still filled in.
	if s := parse(t, desc, key); len(s.Servers) != 1 || len(s.Accounts) != 0 {
		t.Fatalf("%+v", s)
	}
}

// One fixed header beside the key (Harvest-Account-Id): an operation that
// needs it as a parameter has it, never as an argument; it can't be a
// credential's or the request's own header.
func TestAFixedHeaderIsSentForTheOperationsThatNeedIt(t *testing.T) {
	needs := strings.Replace(desc, `"/signed": {
   "get": {"operationId": "signed", "parameters": [{"name": "X-Signature", "in": "header", "required": true}]}`, `"/signed": {
   "get": {"operationId": "signed", "parameters": [{"name": "Account-Id", "in": "header", "required": true}]}`, 1)
	if needs == desc {
		t.Fatal("setup")
	}
	without := parse(t, needs, key)
	for _, o := range without.Operations {
		if o.Name == "signed" {
			t.Fatalf("an operation that needs a header no one sends: %+v", o)
		}
	}
	fixed := key
	fixed.Fixed = &Fixed{Name: "account-id", Value: "123456"}
	signed := op(t, parse(t, needs, fixed), "signed")
	if strings.Contains(string(signed.InputSchema), "Account-Id") || len(signed.Params) != 0 {
		t.Fatalf("the fixed header became an argument: %+v %s", signed, signed.InputSchema)
	}
	for name, f := range map[string]Fixed{
		"Authorization": {Name: "Authorization", Value: "x"},
		"the key's own": {Name: "X-API-Key", Value: "x"},
		"Host":          {Name: "Host", Value: "elsewhere.example"},
		"Cookie":        {Name: "Cookie", Value: "x"},
		"a forwarded":   {Name: "X-Forwarded-For", Value: "10.0.0.1"},
		"a proxy's":     {Name: "Proxy-Authorization", Value: "x"},
		"a newline":     {Name: "Account-Id", Value: "1\r\nHost: elsewhere"},
		"no value":      {Name: "Account-Id"},
		"not a name":    {Name: "Account Id", Value: "1"},
		"too long":      {Name: "Account-Id", Value: strings.Repeat("1", MaxBasicPart+1)},
	} {
		a := key
		a.Fixed = &f
		var r *Refused
		if _, err := Parse([]byte(needs), a); !errors.As(err, &r) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// An account's address stays on the vendor's own server: each variable one
// label, https only (the PM on planque#160).
func TestAnAccountsAddressStaysOnTheVendorsServer(t *testing.T) {
	accounts := []string{"https://{subdomain}.zendesk.com/api/v2"}
	for _, ok := range []string{"https://acme.zendesk.com/api/v2", "https://acme-eu.zendesk.com/api/v2/", "https://a1.zendesk.com/api/v2"} {
		if !OnAccountServer(accounts, ok) {
			t.Errorf("refused: %s", ok)
		}
	}
	for _, bad := range []string{
		"https://evil.example/api/v2",
		"https://acme.zendesk.com.evil.example/api/v2",
		"https://evil.example/x.zendesk.com/api/v2",
		"https://a.b.zendesk.com/api/v2",
		"http://acme.zendesk.com/api/v2",
		"https://acme.zendesk.com/api/v2/../../elsewhere",
		"https://acme.zendesk.com/api/v3",
		"https://-acme.zendesk.com/api/v2",
		"https://acme.zendesk.com@evil.example/api/v2",
	} {
		if OnAccountServer(accounts, bad) {
			t.Errorf("taken: %s", bad)
		}
	}
}

// An account server whose host is all variables would take any domain: no
// base is ever on it (qa-kiff on planque#160).
func TestAHostOfVariablesTakesNoBase(t *testing.T) {
	for _, a := range []string{"https://{domain}.{tld}/api", "https://{host}/api", "https://{sub}.{domain}.com/api", "https://acme.{domain}/api"} {
		for _, base := range []string{"https://evil.example/api", "https://acme.evil.com/api", "https://acme.evil/api"} {
			if OnAccountServer([]string{a}, base) {
				t.Errorf("%s took %s", a, base)
			}
		}
	}
	// The vendor's domain written out, the variable in its host or its path.
	if !OnAccountServer([]string{"https://{sub}.vendor.example.com/api"}, "https://acme.vendor.example.com/api") ||
		!OnAccountServer([]string{"https://api.vendor.com/{account}/v1"}, "https://api.vendor.com/acme/v1") {
		t.Fatal("a vendor's own domain refused")
	}
}
