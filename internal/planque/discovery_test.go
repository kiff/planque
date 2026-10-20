package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// kiff/planque#23 (1): before the owner approves a grant's tools, nothing is
// read from their account. Connect-first lists the vendor's tools and calls
// none: describe is gone, so a no-argument read of business data (the fake's
// charges, with a customer's email) can't be called through it.
func TestNothingIsReadBeforeTheToolsAreApproved(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	for _, tool := range []string{"list_charges", "list_accounts"} {
		if st, out := e.do("POST", acct+"/describe", map[string]string{"tool": tool}, nil, "k_setup"); st != 410 {
			t.Fatalf("describe %s: %d %v", tool, st, out)
		}
	}
	st, out := e.do("GET", acct+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	b, _ := json.Marshal(out)
	if strings.Contains(string(b), "ch_customer_1") || strings.Contains(string(b), "customer@example.com") {
		t.Fatalf("the catalog carries the account's data: %s", b)
	}
	if calls, err := e.vendor.CallsSoFar(t.Context()); err != nil || len(calls) != 0 {
		t.Fatalf("the vendor ran %v before any tool was approved", calls)
	}
}
