package planque

import (
	"strings"
	"testing"
)

// sharesOf reads a row's shares_sign_in as grant ids, and whether it was there.
func sharesOf(t *testing.T, row map[string]any) ([]string, []map[string]any, bool) {
	t.Helper()
	raw, ok := row["shares_sign_in"]
	if !ok {
		return nil, nil, false
	}
	list, ok := raw.([]any)
	if !ok {
		t.Fatalf("shares_sign_in isn't a list: %v", raw)
	}
	var ids []string
	var entries []map[string]any
	for _, x := range list {
		m := x.(map[string]any)
		ids = append(ids, m["grant"].(string))
		entries = append(entries, m)
	}
	return ids, entries, true
}

// A connection's row names the account's other live grants on its sign-in,
// pending ones included, as removal counts them, so an app can tell a
// connection no taf needs from one whose sign-in another grant still uses.
// A grant ended uses nothing, and a row names only grant ids and tafs.
func TestARowNamesTheOtherLiveGrantsOnItsSignIn(t *testing.T) {
	e := newEnv(t)
	first := e.connect()
	e.kiff.Cards["k_taf"]["list_charges_b"] = 0
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "list_charges", "title": "List charges b", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse: %d %v", st, out)
	}
	second := out["id"].(string)

	// The second waits for its approval: it counts, as removal counts it.
	ids, entries, ok := sharesOf(t, listRow(t, e, first.ID, ""))
	if !ok || len(ids) != 1 || ids[0] != second {
		t.Fatalf("the first's row with the second pending: %v (present %t)", ids, ok)
	}
	if entries[0]["agent"] != "agt_taf" {
		t.Fatalf("the other grant's taf: %v", entries[0])
	}
	for k := range entries[0] {
		if k != "grant" && k != "agent" && k != "taf" {
			t.Fatalf("a shared sign-in shows %q: %v", k, entries[0])
		}
	}
	if s := strings.ToLower(mustJSONString(entries)); strings.Contains(s, "tok") {
		t.Fatalf("a shared sign-in names a token: %s", s)
	}
	if ids, _, ok := sharesOf(t, listRow(t, e, second, "")); !ok || len(ids) != 1 || ids[0] != first.ID {
		t.Fatalf("the second's row: %v (present %t)", ids, ok)
	}

	// A single grant's answer, which a taf's own key can read, names no other
	// taf's grant: only the phone's signed list does.
	for _, key := range []string{"k_setup", "k_taf"} {
		st, one := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+first.ID, nil, nil, key)
		if st != 200 {
			t.Fatalf("get with %s: %d %v", key, st, one)
		}
		if _, ok := one["shares_sign_in"]; ok {
			t.Fatalf("a single answer to %s names the grants on its sign-in: %v", key, one["shares_sign_in"])
		}
		if s := mustJSONString(one); strings.Contains(s, second) {
			t.Fatalf("a single answer to %s names another grant: %s", key, s)
		}
	}

	// The second ended: the first's sign-in is its own again, an empty list.
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+second, nil, e.phone, ""); st != 200 {
		t.Fatal("remove the second")
	}
	if ids, _, ok := sharesOf(t, listRow(t, e, first.ID, "")); !ok || len(ids) != 0 {
		t.Fatalf("the first's row once the second ended: %v (present %t)", ids, ok)
	}
}

// Two connections of a provider's app on the same vendor account name each
// other, as removal keeps the account for the other; two named different
// accounts name nobody.
func TestARowNamesAConnectionThatMayBeOnTheSameVendorAccount(t *testing.T) {
	e, f := newConnectorFixture(t)
	a, b := twoOwnSignIns(t, e, f) // both on "Demo Shop"
	if ids, _, ok := sharesOf(t, listRow(t, e, a.ID, "")); !ok || len(ids) != 1 || ids[0] != b.ID {
		t.Fatalf("same vendor account: %v (present %t)", ids, ok)
	}
	b, err := e.s.loadGrant(t.Context(), tenant, b.ID) // as the list left it
	if err != nil {
		t.Fatal(err)
	}
	b.Connector.AccountName = "Another Shop"
	if err := e.s.saveGrant(t.Context(), b); err != nil {
		t.Fatal(err)
	}
	if ids, _, ok := sharesOf(t, listRow(t, e, a.ID, "")); !ok || len(ids) != 0 {
		t.Fatalf("different vendor accounts: %v (present %t)", ids, ok)
	}
}

// What a row says is what removal does: every grant a row names is one
// signInUsers counts, and the other way round.
func TestARowsSharedSignInsAreRemovalsUsers(t *testing.T) {
	e := newEnv(t)
	first := e.connect()
	e.kiff.Cards["k_taf"]["list_charges_b"] = 0
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "reuse": first.ID, "agent": "agt_taf",
		"tools": []map[string]any{{"tool": "list_charges", "title": "List charges b", "name": "list_charges_b", "read_only": true}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("reuse: %d %v", st, out)
	}
	ctx := t.Context()
	all, err := e.s.accountGrants(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range all {
		users, err := e.s.signInUsers(ctx, g)
		if err != nil {
			t.Fatal(err)
		}
		var shown []string
		for _, x := range sharedSignInsOf(g, all) {
			shown = append(shown, x.Grant)
		}
		if strings.Join(users, ",") != strings.Join(shown, ",") {
			t.Fatalf("%s: removal counts %v, the row shows %v", g.ID, users, shown)
		}
	}
}
