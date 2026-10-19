package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// While it looks, the asking agent reads what the owner's sign-in covers
// (kiff/tafs#89): it may call a vendor tool that is read-only and takes no
// arguments, and nothing else. planque records each read, and the vendor
// counts no call. The owner then approves tools whose pinned arguments come
// with the vendor's own words for them.
func TestTheAgentReadsWhatTheSignInCovers(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	st, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	describe := func(key, tool string) (int, map[string]any) {
		return e.do("POST", acct+"/"+id+"/describe", map[string]string{"tool": tool}, nil, key)
	}
	if st, _ := describe("k_setup", "list_accounts"); st != 409 {
		t.Fatalf("before the sign-in: %d", st)
	}
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusListed {
		t.Fatalf("confirm: %d %v", st, out)
	}

	st, out = describe("k_setup", "list_accounts")
	if st != 200 || out["tool"] != "list_accounts" {
		t.Fatalf("describe: %d %v", st, out)
	}
	if b, _ := json.Marshal(out["result"]); !strings.Contains(string(b), "acct_fake123") || !strings.Contains(string(b), "test") {
		t.Fatalf("the sign-in's accounts: %s", b)
	}
	if e.vendor.CallCount() != 0 {
		t.Fatal("reading what the sign-in covers counted as a vendor call")
	}
	for name, c := range map[string]struct {
		key, tool string
		want      int
	}{
		"another agent":                    {"k_taf", "list_accounts", 403},
		"a tool that changes things":       {"k_setup", "stripe_api_write", 403},
		"a tool that takes arguments":      {"k_setup", "create_refund", 403},
		"a no-argument tool not read-only": {"k_setup", "list_charges", 403},
		"a tool the vendor doesn't list":   {"k_setup", "delete_account", 404},
	} {
		if st, _ := describe(c.key, c.tool); st != c.want {
			t.Errorf("%s: %d, want %d", name, st, c.want)
		}
	}
	if e.vendor.CallCount() != 0 {
		t.Fatal("a refused describe reached the vendor")
	}
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	reads := 0
	for _, ev := range g.Events {
		if strings.HasPrefix(ev.What, "read list_accounts") {
			reads++
		}
	}
	if reads != 1 {
		t.Fatalf("reads recorded: %d", reads)
	}

	// The tools approval explains each pinned argument in the vendor's words.
	refund := map[string]any{"tool": "create_refund", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2,
		"fixed_arguments": map[string]any{"mode": "test"},
		"arguments":       map[string]any{"charge": map[string]any{"path": "charge", "type": "string", "required": true}, "amount": map[string]any{"path": "amount", "type": "integer", "required": true}}}
	st, out = e.do("POST", acct+"/"+id+"/tools", map[string]any{"tools": []any{refund}}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("set tools: %d %v", st, out)
	}
	if text := out["text"].(string); !strings.Contains(text, `mode = "test" (the vendor says: "Whether to act in live mode (real money) or test mode.")`) {
		t.Fatalf("the approval doesn't give the vendor's words for the pinned mode:\n%s", text)
	}
}

// Enough reads: a grant reads what its sign-in covers a bounded number of times.
func TestDescribeIsBounded(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	acct := "/v1/accounts/" + tenant + "/grants"
	_, out := e.do("POST", acct, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	id := out["id"].(string)
	if st, out := e.approveAndSignIn(id, out["digest"].(string)); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	for i := 0; i < MaxDescribes; i++ {
		if st, out := e.do("POST", acct+"/"+id+"/describe", map[string]string{"tool": "list_accounts"}, nil, "k_setup"); st != 200 {
			t.Fatalf("read %d: %d %v", i+1, st, out)
		}
	}
	if st, _ := e.do("POST", acct+"/"+id+"/describe", map[string]string{"tool": "list_accounts"}, nil, "k_setup"); st != 429 {
		t.Fatalf("one read too many: %d", st)
	}
}

func TestTakesNoArguments(t *testing.T) {
	for schema, want := range map[string]bool{
		``:                                  true,
		`{"type":"object","properties":{}}`: true,
		`{"type":"object"}`:                 true,
		`{"type":"object","properties":{"id":{"type":"string"}}}`: false,
		`{"type":"object","required":["id"]}`:                     false,
		`not json`:                                                false,
	} {
		if got := takesNoArguments(json.RawMessage(schema)); got != want {
			t.Errorf("takesNoArguments(%s) = %v", schema, got)
		}
	}
}
