package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const enumDesc = `{
 "openapi": "3.0.3",
 "info": {"title": "Desk API"},
 "security": [{"apiKey": []}],
 "components": {"securitySchemes": {"apiKey": {"type": "apiKey", "in": "header", "name": "X-API-Key"}}},
 "paths": {
  "/workspaces/{workspace}/tickets": {"get": {"operationId": "searchTickets", "summary": "Search the tickets",
   "parameters": [{"name": "workspace", "in": "path", "required": true, "schema": {"type": "string"}},
    {"name": "status", "in": "query", "schema": {"type": "string", "enum": ["open", "pending", "solved"]}}]}}
 }
}`

// A narrowed argument keeps the values the vendor's description allows (its
// enum): the agent is shown "open, pending or solved", and a value the vendor
// doesn't list ("new", which the vendor answers with an empty list, not an
// error) is refused before anything is sent, naming the values it takes. A
// value it lists goes through as before.
func TestANarrowedArgumentKeepsTheVendorsAllowedValues(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	v.Server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.mu.Lock()
		defer v.mu.Unlock()
		if r.URL.Path == "/openapi.json" {
			_, _ = io.WriteString(w, enumDesc)
			return
		}
		v.calls = append(v.calls, restCall{Method: r.Method, Path: r.URL.Path + "?" + r.URL.RawQuery, Key: r.Header.Get("X-API-Key")})
		if r.Header.Get("X-API-Key") != v.key {
			w.WriteHeader(401)
			return
		}
		_, _ = io.WriteString(w, `{"data":[{"id":"T-1","status":"open"}]}`)
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
	tools := []map[string]any{{"tool": "searchTickets", "name": "search_tickets", "title": "Search the tickets", "read_only": true,
		"fixed_arguments": map[string]any{"workspace": "ws_eval"},
		"arguments":       map[string]any{"status": map[string]any{"path": "status", "type": "string"}}}}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": tools}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}

	// What the agent is shown: the vendor's three values on status.
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties map[string]struct {
			Enum []string `json:"enum"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(g.Tools[0].InputSchema, &schema); err != nil || strings.Join(schema.Properties["status"].Enum, ",") != "open,pending,solved" {
		t.Fatalf("the agent's schema: %s", g.Tools[0].InputSchema)
	}

	e.kiff.Cards["k_taf"] = map[string]int64{"search_tickets": 0}
	ctx := context.Background()
	res, err := e.kiff.Call(ctx, "k_taf", "search_tickets", map[string]any{"status": "new"})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("%+v %v", res, err)
	}
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "must be one of open, pending, solved") || !strings.Contains(text, "Nothing was sent") {
		t.Fatalf("a value the vendor doesn't list: %s", text)
	}
	if n := len(v.seen()); n != 0 {
		t.Fatalf("the vendor was called: %+v", v.seen())
	}

	res, _ = e.kiff.Call(ctx, "k_taf", "search_tickets", map[string]any{"status": "open"})
	if text, isErr := resultText(t, res.Result); isErr || !strings.Contains(text, "T-1") {
		t.Fatalf("a value it lists: %s", text)
	}
	if c := v.seen(); len(c) != 1 || c[0].Path != "/v1/workspaces/ws_eval/tickets?status=open" {
		t.Fatalf("sent %+v", c)
	}
	// Left out, it's every ticket, as before.
	res, _ = e.kiff.Call(ctx, "k_taf", "search_tickets", map[string]any{})
	if _, isErr := resultText(t, res.Result); isErr || len(v.seen()) != 2 {
		t.Fatalf("no status: %+v", v.seen())
	}
}

// Only what the vendor lists, of the argument's type and within bounds, is
// carried; anything else leaves the argument a plain one, as it was.
func TestOnlyAUsableVendorEnumIsCarried(t *testing.T) {
	desc := func(enum string) json.RawMessage {
		return json.RawMessage(`{"type":"object","properties":{"filter":{"type":"object","properties":{"status":{"type":"string","enum":` + enum + `}}}}}`)
	}
	if got := vendorEnum(desc(`["a","b"]`), "filter.status", "string"); len(got) != 2 {
		t.Fatalf("a nested enum: %v", got)
	}
	long := `["` + strings.Repeat("x", maxEnumValue+1) + `"]`
	many := `[` + strings.Repeat(`"x",`, maxEnum) + `"y"]`
	for name, c := range map[string]struct {
		schema json.RawMessage
		path   string
		typ    string
	}{
		"another type":        {desc(`["a", 1]`), "filter.status", "string"},
		"no such path":        {desc(`["a"]`), "filter.state", "string"},
		"no enum":             {json.RawMessage(`{"properties":{"status":{"type":"string"}}}`), "status", "string"},
		"a value too long":    {desc(long), "filter.status", "string"},
		"too many values":     {desc(many), "filter.status", "string"},
		"no vendor schema":    {nil, "status", "string"},
		"an integer argument": {desc(`["a"]`), "filter.status", "integer"},
	} {
		if got := vendorEnum(c.schema, c.path, c.typ); got != nil {
			t.Errorf("%s: carried %v", name, got)
		}
	}
	// Every argument of a call is checked, not only the first that fits.
	schema := json.RawMessage(`{"properties":{"a":{"enum":["x"]},"b":{"enum":["y"]}}}`)
	if err := outsideEnum(schema, map[string]any{"a": "x", "b": "z"}); err == nil || !strings.Contains(err.Error(), "argument b must be one of y") {
		t.Fatalf("the second argument: %v", err)
	}
	if err := outsideEnum(schema, map[string]any{"a": "x", "b": "y"}); err != nil {
		t.Fatal(err)
	}
	if err := outsideEnum(json.RawMessage(`{"properties":{"n":{"enum":[1,2]}}}`), map[string]any{"n": json.Number("2")}); err != nil {
		t.Fatalf("a number it lists: %v", err)
	}
}
