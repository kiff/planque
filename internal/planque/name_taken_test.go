package planque

import (
	"strings"
	"testing"
)

// kiff/tafs#138: a tool name another of the account's connections holds is
// refused when the agent sets the tools, not after the owner approved, when
// KIFF would refuse to register it. The grant's own name stays its own, and
// a tool left unnamed gets planque's per-grant name, which no other holds.
func TestAToolNameAnotherConnectionHoldsIsRefusedFirst(t *testing.T) {
	e := newEnv(t)
	first := e.listedGrant()
	// A second connection on the same, already enrolled, phone.
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	second := out["id"].(string)
	if st, out := e.approveAndSignIn(second, out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		t.Fatalf("the second connection: %d %v", st, out)
	}
	at := func(id string) string { return "/v1/accounts/" + tenant + "/grants/" + id + "/tools" }
	read := map[string]any{"tools": []any{map[string]any{"tool": "retrieve_charge", "name": "read_payment", "read_only": true}}}
	st, out := e.do("POST", at(first), read, nil, "k_setup")
	if st != 200 {
		t.Fatalf("the first connection's read: %d %v", st, out)
	}
	if st, out := e.do("POST", at(first), read, nil, "k_setup"); st != 200 {
		t.Fatalf("its own name again: %d %v", st, out)
	}
	st, out = e.do("POST", at(second), read, nil, "k_setup")
	if st != 409 || !strings.Contains(out["error"].(string), "read_payment is already connected") || !strings.Contains(out["error"].(string), "another name") {
		t.Fatalf("the same name on another connection: %d %v", st, out)
	}
	unnamed := map[string]any{"tools": []any{map[string]any{"tool": "retrieve_charge", "read_only": true}}}
	if st, out := e.do("POST", at(second), unnamed, nil, "k_setup"); st != 200 {
		t.Fatalf("left unnamed: %d %v", st, out)
	}
}
