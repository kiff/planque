package planque

import (
	"context"
	"strings"
	"testing"
	"time"
)

// A key grant asked with its tools has them declared in KIFF at once, so a
// Card can be drafted over them before the key: the view says so, and a
// call to either is refused before anything reaches the vendor, saying it
// waits for its key. Given the key, the same registration serves the
// grant, its credential unchanged, and a call goes through.
func TestAKeyGrantsToolsAreDeclaredAtTheRequestAndUsableOnlyWithTheKey(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	st, out := e.askAPIGrant(v, apiTools, nil)
	if st != 201 || out["status"] != StatusPending || out["declared"] != true {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	declared := map[string]string{}
	for _, name := range []string{"create_refund", "list_charges"} {
		c, ok := e.kiff.Connection(tenant, name)
		if !ok {
			t.Fatalf("%s isn't declared in KIFF", name)
		}
		declared[name] = c.Credential
	}

	ctx := context.Background()
	e.kiff.Cards["k_taf"]["list_charges"] = 0
	res, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{"limit": 2})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("%+v %v", res, err)
	}
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "waits for its key") || !strings.Contains(text, "Nothing was sent") {
		t.Fatalf("a call before the key: %s", text)
	}
	if n := len(v.seen()); n != 0 {
		t.Fatalf("the vendor heard of it: %+v", v.seen())
	}

	if st, out := e.giveKey(id, out["digest"].(string), v.key); st != 200 || out["status"] != StatusConnected || out["declared"] == true {
		t.Fatalf("key: %d %v", st, out)
	}
	for name, cred := range declared {
		if c, ok := e.kiff.Connection(tenant, name); !ok || c.Credential != cred {
			t.Fatalf("%s was registered again, or went, at connect", name)
		}
	}
	res, _ = e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{"limit": 2})
	if text, isErr := resultText(t, res.Result); isErr {
		t.Fatalf("a call with the key: %s", text)
	}
	if n := len(v.seen()); n != 1 {
		t.Fatalf("the vendor heard %d calls, want 1", n)
	}
}

// A request that ends before its key takes its declarations with it:
// cancelled by the owner, ended by its 7-day deadline, or ended with its
// taf. No Card can then be drafted over a request that's gone.
func TestADeclaredRequestThatEndsLeavesNothingInKIFF(t *testing.T) {
	for _, how := range []string{"cancelled", "past its deadline", "ended with its taf"} {
		t.Run(how, func(t *testing.T) {
			e := newEnv(t)
			v := newRESTVendor(t)
			e.enrol()
			e.manager()
			st, out := e.askAPIGrant(v, apiTools, nil)
			if st != 201 {
				t.Fatalf("ask: %d %v", st, out)
			}
			id := out["id"].(string)
			if len(e.kiff.Connections(tenant)) != 2 {
				t.Fatalf("declared: %v", e.kiff.Connections(tenant))
			}
			switch how {
			case "cancelled":
				if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 {
					t.Fatalf("cancel: %d %v", st, out)
				}
			case "past its deadline":
				e.at(pendingTTL + time.Minute)
				if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
					t.Fatalf("sweep: %d %v", n, err)
				}
			case "ended with its taf":
				if st, out := e.end("agt_taf", "s_manager", ""); st != 200 {
					t.Fatalf("end: %d %v", st, out)
				}
			}
			if c := e.kiff.Connections(tenant); len(c) != 0 {
				t.Fatalf("still in KIFF: %v", c)
			}
			if g, _ := e.s.loadGrant(t.Context(), tenant, id); g.Declared || len(g.Tools) != 0 && g.Status == StatusPending {
				t.Fatalf("the grant: %s, declared %v", g.Status, g.Declared)
			}
		})
	}
}

// Only a key grant asked with its tools, waiting for its key, is declared:
// one asked without tools picks them after the key, and a sign-in's tools
// come from the vendor only once signed in.
func TestOnlyAKeyGrantAskedWithItsToolsIsDeclared(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	if st, out := e.askAPIGrant(v, []map[string]any{}, nil); st != 201 || out["declared"] == true {
		t.Fatalf("no tools: %d %v", st, out)
	}
	if _, _ = e.request(); len(e.kiff.Connections(tenant)) != 0 {
		t.Fatalf("declared: %v", e.kiff.Connections(tenant))
	}
}

// A key grant declined before its key withdraws its declarations at once,
// rather than after its deadline: nothing is left in KIFF, the request ends
// failed, and a second withdrawal finds nothing waiting.
func TestADeclaredKeyGrantDeclinedLeavesNothingInKIFF(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	st, out := e.askAPIGrant(v, apiTools, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("ask: %d %v", st, out)
	}
	acct := "/v1/accounts/" + tenant + "/grants/" + out["id"].(string)
	if st, out := e.do("DELETE", acct+"/tools", nil, nil, "k_other"); st == 200 {
		t.Fatalf("another account's agent withdrew it: %v", out)
	}
	st, out = e.do("DELETE", acct+"/tools", nil, nil, "k_setup")
	if st != 200 || out["status"] != StatusFailed || out["declared"] == true {
		t.Fatalf("decline: %d %v", st, out)
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("still in KIFF: %v", c)
	}
	if st, _ := e.do("DELETE", acct+"/tools", nil, nil, "k_setup"); st != 409 {
		t.Fatalf("declined again: %d", st)
	}
	if st, _ := e.giveKey(out["id"].(string), out["digest"].(string), v.key); st == 200 {
		t.Fatal("a key was taken for a declined request")
	}
}
