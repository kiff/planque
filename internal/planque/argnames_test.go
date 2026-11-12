package planque

import (
	"encoding/json"
	"net/url"
	"os"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// A vendor's own argument names the model can't take ("filter[status]")
// are given to the agent in a spelling it takes, and put back on the call,
// so a whole tool list isn't refused for one name. Qonto's shipped
// description lists them, and its run failed on them.
func TestTheAgentIsGivenArgumentNamesItsModelTakes(t *testing.T) {
	names := argNames(map[string]any{"filter[status]": 1, "filter[]": 1, "filter": 1, "page": 1, "a b": 1, "[]": 1})
	want := map[string]string{"filter[status]": "filter_status", "filter[]": "filter_2", "filter": "filter", "page": "page", "a b": "a_b", "[]": "arg"}
	for k, v := range want {
		if names[k] != v {
			t.Fatalf("%q is given as %q, want %q (%v)", k, names[k], v, names)
		}
	}

	description, err := os.ReadFile("../../deploy/specs/qonto-invoicing-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	snap, err := openapi.Parse(description, openapi.Auth{In: "header", Name: "Authorization"})
	if err != nil {
		t.Fatal(err)
	}
	g := &Grant{}
	for _, o := range snap.Operations {
		g.Tools = append(g.Tools, GrantedTool{ToolRequest: ToolRequest{Name: o.Name, Tool: o.Name}, InputSchema: o.InputSchema})
	}
	bracketed := 0
	for _, tool := range relayTools(g) {
		var schema struct {
			Properties map[string]any `json:"properties"`
			Required   []string       `json:"required"`
		}
		if err := json.Unmarshal(tool.InputSchema, &schema); err != nil {
			t.Fatal(err)
		}
		for n := range schema.Properties {
			if !modelArg.MatchString(n) {
				t.Fatalf("%s gives the agent %q", tool.Name, n)
			}
			if n == "filter_status" {
				bracketed++
			}
		}
		for _, r := range schema.Required {
			if _, ok := schema.Properties[r]; !ok {
				t.Fatalf("%s requires %q, which it doesn't list", tool.Name, r)
			}
		}
	}
	if bracketed == 0 {
		t.Fatal("Qonto's filter[status] wasn't given as filter_status")
	}

	// The call: the agent's names, put back as Qonto's own in the request.
	for i, o := range snap.Operations {
		if o.Name != "list_client_invoices" {
			continue
		}
		args := map[string]any{"filter_status": "unpaid", "filter_due_date_to": "2026-11-12"}
		fromAgent(g.Tools[i].InputSchema, args)
		_, target, _, err := o.Request("https://thirdparty.qonto.com/v2", args)
		if err != nil {
			t.Fatal(err)
		}
		u, _ := url.Parse(target)
		if q := u.Query(); q.Get("filter[status]") != "unpaid" || q.Get("filter[due_date_to]") != "2026-11-12" {
			t.Fatalf("the request: %s", target)
		}
		return
	}
	t.Fatal("no list_client_invoices in Qonto's description")
}
