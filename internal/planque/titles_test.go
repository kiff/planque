package planque

import (
	"strings"
	"testing"
)

// Nothing reaches the owner's approval named by a tool's name: an agent's ask
// for a new connection, or for tools added to one, refuses any tool without a
// title, and names the ones it found, in the words kiff/agents asks with. Nothing
// is made: no grant, and no tools added.
func TestEveryToolAskedForNeedsATitle(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	grants := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", grants, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "create_refund", "title": "Refund a charge", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2},
			{"tool": "list_charges", "read_only": true}}}, nil, "k_setup")
	msg, _ := out["error"].(string)
	if st != 400 || !strings.Contains(msg, "each tool needs a title") || !strings.Contains(msg, "list_charges has none") || strings.Contains(msg, "create_refund") {
		t.Fatalf("a request with an untitled tool: %d %v", st, out)
	}
	if items, err := e.s.KV.Query(t.Context(), tenantPK(tenant), "GRANT#"); err != nil || len(items) != 0 {
		t.Fatalf("a refused request made a grant: %v %d", err, len(items))
	}

	id := e.connectedRefunds()
	before, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	st, out = e.do("POST", grants+"/"+id+"/tools", map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	if msg, _ := out["error"].(string); st != 400 || !strings.Contains(msg, "list_charges has none") {
		t.Fatalf("adding an untitled tool: %d %v", st, out)
	}
	got, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil || len(got.Request.Tools) != len(before.Request.Tools) {
		t.Fatalf("a refused addition changed the grant's tools: %v %v", err, got.Request.Tools)
	}
}
