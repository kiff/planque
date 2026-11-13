package planque

import (
	"strings"
	"testing"
)

// A refused connection request answers, beside its words, a code and the
// tool and field at fault, so Coach can say in plain words what to change.
func TestARefusalCarriesItsCodeAndField(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	refund := func(extra map[string]any) []map[string]any {
		tool := map[string]any{"tool": "createRefund", "name": "create_refund", "title": "Refund a charge"}
		for k, x := range extra {
			tool[k] = x
		}
		return []map[string]any{tool}
	}
	for name, c := range map[string]struct {
		tools []map[string]any
		extra map[string]any
		want  map[string]string
		words string
	}{
		"a pin that isn't a dotted path": {
			tools: refund(map[string]any{"fixed_arguments": map[string]any{"post author": "x"}}),
			want:  map[string]string{"code": "pin_shape", "tool": "create_refund", "field": "post author", "want": "dotted path → JSON value"},
			words: `create_refund: fixed argument "post author" isn't a dotted path`,
		},
		"a pin the operation doesn't take": {
			tools: refund(map[string]any{"fixed_arguments": map[string]any{"account": "a_1"}, "arguments": map[string]any{"charge": map[string]any{"path": "charge", "type": "string"}}}),
			want:  map[string]string{"code": "pin_unknown", "tool": "create_refund", "field": "account", "want": "amount, charge"},
			words: "fixed argument account isn't one of createRefund's arguments (amount, charge)",
		},
		"a narrowed tool without a required argument": {
			tools: refund(map[string]any{"arguments": map[string]any{"amount": map[string]any{"path": "amount", "type": "integer"}}}),
			want:  map[string]string{"code": "required_argument", "tool": "create_refund", "field": "charge"},
			words: "createRefund requires charge",
		},
		"no agent": {
			tools: refund(nil),
			extra: map[string]any{"agent": " "},
			want:  map[string]string{"code": "agent_required", "field": "agent"},
			words: "agent is required",
		},
	} {
		st, out := e.askAPIGrant(v, c.tools, c.extra)
		if st != 400 {
			t.Errorf("%s: %d %v", name, st, out)
			continue
		}
		for k, w := range c.want {
			if out[k] != w {
				t.Errorf("%s: %s is %v, want %q (%v)", name, k, out[k], w, out)
			}
		}
		if _, named := out["tool"]; c.want["tool"] == "" && named {
			t.Errorf("%s: names a tool: %v", name, out)
		}
		if words, _ := out["error"].(string); !strings.Contains(words, c.words) {
			t.Errorf("%s: words %q, want %q", name, words, c.words)
		}
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("a refused request reached KIFF: %+v", c)
	}
}
