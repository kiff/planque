package planque

import (
	"strings"
	"testing"
)

// kiff/agents#139: read-only is the vendor's word. The asking agent can't
// grant a vendor write as a read: not when it sets a listed grant's tools,
// and not when it asks for known tools, which fail closed at connect,
// before anything reaches KIFF.
func TestReadOnlyIsTheVendorsWord(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	at := "/v1/accounts/" + tenant + "/grants/" + id + "/tools"
	st, out := e.do("POST", at, map[string]any{"tools": []any{map[string]any{"tool": "payout_balance", "read_only": true}}}, nil, "k_setup")
	if st != 400 || !strings.Contains(out["error"].(string), "isn't read-only") {
		t.Fatalf("a vendor write asked for as a read: %d %v", st, out)
	}
	if st, out := e.do("POST", at, map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 200 ||
		!strings.Contains(out["text"].(string), "reads only") {
		t.Fatalf("a vendor read: %d %v", st, out)
	}

	// Known tools, no catalog yet: checked against the vendor's list at connect.
	e2 := newEnv(t)
	e2.enrol()
	st, out = e2.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e2.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{map[string]any{"tool": "payout_balance", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	if st, out := e2.approveAndSignIn(out["id"].(string), out["digest"].(string)); st == 200 && out["status"] == StatusConnected {
		t.Fatalf("a vendor write connected as a read: %d %v", st, out)
	}
	if _, ok := e2.kiff.Connection(tenant, "payout_balance"); ok {
		t.Fatal("KIFF has the write registered as a read")
	}
	if a, r := e2.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the sign-in wasn't revoked: %d, %d", a, r)
	}
}
