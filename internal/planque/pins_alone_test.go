package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// A read whose every argument is pinned (fulfilled orders, its status
// fixed) is a tool of its own: the taf sets nothing, and planque sends the
// pins alone. Refusing it for naming no arguments failed a key request on
// its first try in the Eval of 13 Nov.
func TestAToolWhoseArgumentsAreAllPinnedIsDeclared(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	v.desc = strings.Replace(restDesc, `"parameters": [{"name": "limit"`, `"parameters": [{"name": "status", "in": "query", "required": true, "schema": {"type": "string"}}, {"name": "limit"`, 1)
	e.enrol()
	read := map[string]any{"tool": "listCharges", "name": "list_fulfilled", "title": "List fulfilled charges", "read_only": true,
		"fixed_arguments": map[string]any{"status": "fulfilled", "limit": 50}}
	st, out := e.askAPIGrant(v, []map[string]any{read}, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("a read with every argument pinned: %d %v", st, out)
	}
	if terms := mustJSONString(out["terms"]); !strings.Contains(terms, `the taf sets nothing`) || !strings.Contains(terms, `status = \"fulfilled\"`) {
		t.Fatalf("the owner's words: %s", terms)
	}

	// The taf's call sends the pins, and can't add an argument.
	tr := ToolRequest{Tool: "listCharges", Fixed: map[string]json.RawMessage{"status": json.RawMessage(`"fulfilled"`), "limit": json.RawMessage(`50`)}}
	got, err := tr.vendorArgs(map[string]any{})
	if err != nil || string(got) != `{"limit":50,"status":"fulfilled"}` {
		t.Fatalf("sent %s %v", got, err)
	}
	if _, err := tr.vendorArgs(map[string]any{"status": "open"}); err == nil || !strings.Contains(err.Error(), "not part of this grant") {
		t.Fatalf("a pinned argument set by the taf: %v", err)
	}
	if s := string(tr.schema()); !strings.Contains(s, `"properties":{}`) || !strings.Contains(s, `"additionalProperties":false`) {
		t.Fatalf("the taf's schema: %s", s)
	}

	// A narrowed tool still sends every argument its operation requires: a
	// pin that leaves one out is refused before anything is declared.
	for name, tool := range map[string]map[string]any{
		"a read missing its required status":  {"tool": "listCharges", "title": "List charges", "read_only": true, "fixed_arguments": map[string]any{"limit": 50}},
		"a write missing its required amount": {"tool": "createRefund", "title": "Refund a charge", "fixed_arguments": map[string]any{"charge": "ch_1"}},
		"a resource the taf can't set":        {"tool": "createRefund", "title": "Refund a charge", "fixed_arguments": map[string]any{"charge": "ch_1", "amount": 100}, "resource_argument": "charge"},
	} {
		st, out := e.askAPIGrant(v, []map[string]any{tool}, nil)
		if st != 400 {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}
	if st, out := e.askAPIGrant(v, []map[string]any{{"tool": "listCharges", "title": "List charges", "read_only": true, "fixed_arguments": map[string]any{"limit": 50}}}, nil); !strings.Contains(mustJSONString(out), "listCharges requires status: pin it in fixed_arguments, or name it in arguments for the taf to set") {
		t.Fatalf("the refusal's words: %d %v", st, out)
	}
}
