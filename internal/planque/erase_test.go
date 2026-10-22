package planque

import (
	"context"
	"strings"
	"testing"
)

// Deleting the account from the app (kiff/tafs#185): with the owner's
// signature, everything planque keeps for the account goes. Its tools leave
// KIFF, the vendor's tokens are revoked, and its records and the
// enrolment are deleted.
func TestDeletingTheAccountErasesEverything(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	acct := "/v1/accounts/" + tenant
	// A second grant, signed in with no tools yet.
	_, out := e.do("POST", acct+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st, out := e.approveAndSignIn(out["id"].(string), out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		t.Fatalf("the listed grant: %d %v", st, out)
	}

	// Not without the owner's signature: an agent's key is refused, and
	// nothing goes.
	if st, _ := e.do("DELETE", acct, nil, nil, "k_setup"); st != 401 {
		t.Fatalf("an agent deleting the account: %d, want 401", st)
	}
	// Nor with an editor's: they connect tools, but only an admin may delete
	// the account, in KIFF as here (review on #45).
	editor := e.device("dev_editor")
	ed := e.kiff.Devices["dev_editor"]
	ed.Capabilities = []string{"holds.answer", "tools.connect"}
	e.kiff.Devices["dev_editor"] = ed
	if st, out := e.do("DELETE", acct, nil, editor, ""); st != 403 || !strings.Contains(out["error"].(string), "can't delete the account") {
		t.Fatalf("an editor deleting the account: %d %v", st, out)
	}
	// Or export it (kiff/planque#46).
	if st, out := e.do("GET", acct+"/export", nil, editor, ""); st != 403 || !strings.Contains(out["error"].(string), "can't export the account") {
		t.Fatalf("an editor exporting the account: %d %v", st, out)
	}
	if n := len(e.kiff.Connections(tenant)); n != 2 {
		t.Fatalf("KIFF lists %d tools after a refused delete, want 2", n)
	}
	if a, _ := e.vendor.Live(); a == 0 {
		t.Fatal("a refused delete revoked the vendor's tokens")
	}

	st, out := e.do("DELETE", acct, nil, e.phone, "")
	if st != 200 || out["erased"] != true || out["grants"] != float64(2) {
		t.Fatalf("delete: %d %v", st, out)
	}
	if ps, _ := out["problems"].([]any); len(ps) != 0 {
		t.Fatalf("problems: %v", ps)
	}
	if n := len(e.kiff.Connections(tenant)); n != 0 {
		t.Fatalf("KIFF still lists %d tools", n)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("the vendor still honours %d access and %d refresh tokens", a, r)
	}
	if items, _ := e.kv.Query(ctx, tenantPK(tenant), ""); len(items) != 0 {
		t.Fatalf("%d records kept for the account", len(items))
	}
	if code := e.relayCallRaw(g); code != 401 {
		t.Fatalf("relay after the delete: %d", code)
	}

	// Asked again, there is nothing kept, and nothing to check a signature
	// with: done.
	if st, _ := e.do("DELETE", acct, nil, e.phone, ""); st != 204 {
		t.Fatalf("delete again: %d, want 204", st)
	}
	// An account planque keeps nothing for has nothing to delete.
	if st, _ := e.do("DELETE", "/v1/accounts/tnt_other", nil, nil, "k_other"); st != 204 {
		t.Fatalf("an account with nothing kept: %d, want 204", st)
	}
}
