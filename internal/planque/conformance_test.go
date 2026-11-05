package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// A conformance kind is one way planque reaches a vendor, set up to the
// point where a call needs only KIFF's permit: a grant connected and its
// tools approved on the owner's device.
type conformanceKind struct {
	name string
	// connect returns the server, the connected grant, the tool to call,
	// the call's arguments for an operation id, and how many calls the
	// vendor has received so far.
	connect func(t *testing.T) (e *env, g *Grant, tool string, args func(op string) map[string]any, sent func() int)
}

// conformanceKinds are every kind of connection planque makes. A new kind
// is added here, and so passes the same refusals as the others.
var conformanceKinds = []conformanceKind{
	{"mcp", func(t *testing.T) (*env, *Grant, string, func(string) map[string]any, func() int) {
		e := newEnv(t)
		g := e.connect()
		return e, g, "create_refund", func(op string) map[string]any {
			return map[string]any{"charge": "ch_2", "amount": 100, IdempotencyArg: op}
		}, e.vendor.CallCount
	}},
	{"rest_api_key", func(t *testing.T) (*env, *Grant, string, func(string) map[string]any, func() int) {
		e := newEnv(t)
		v := newRESTVendor(t)
		g := connectRESTForConformance(t, e, v)
		return e, g, "create_refund", func(op string) map[string]any {
			return map[string]any{"charge": "ch_2", "amount": 100, IdempotencyArg: op}
		}, func() int { return len(v.seen()) }
	}},
	{"browser", func(t *testing.T) (*env, *Grant, string, func(string) map[string]any, func() int) {
		e := newBrowserEnv(t)
		g := e.connect()
		return e.env, g, "browse_act", func(op string) map[string]any {
			return map[string]any{"url": fakes.ShopItem, "steps": []map[string]any{{"fill": "qty", "value": "2"}}, IdempotencyArg: op}
		}, e.b.Sent
	}},
	{"managed_connector", func(t *testing.T) (*env, *Grant, string, func(string) map[string]any, func() int) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		return e, g, "reply_ticket", func(op string) map[string]any {
			return map[string]any{"ticket": "ticket_demo", "text": "hello", IdempotencyArg: op}
		}, func() int { return f.executions }
	}},
}

// connectRESTForConformance connects an API-key grant to v, as
// TestAnAPIKeyGrantEndToEnd does step by step: asked for with no tools,
// the key given on the phone, the tools approved.
func connectRESTForConformance(t *testing.T, e *env, v *restVendor) *Grant {
	t.Helper()
	e.enrol()
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	if st, out := e.giveKey(id, digest, v.key); st != 200 || out["status"] != StatusListed {
		t.Fatalf("key: %d %v", st, out)
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": apiTools}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// Every kind of connection refuses a call that KIFF's permit doesn't
// authorize exactly, and sends the vendor nothing for it: no permit, a
// permit for other arguments, another tool, another grant, another KIFF
// account, one expired or issued for too long, one not signed by KIFF,
// one naming another operation. A permitted call is sent once; the same
// call again (a replay) sends nothing more and answers as the first did.
// planque#71: one set of refusals for every connection kind, so a kind
// added later can't skip one.
func TestConformanceEveryKindRefusesWhatItsPermitDoesntAuthorize(t *testing.T) {
	for _, kind := range conformanceKinds {
		t.Run(kind.name, func(t *testing.T) {
			e, g, tool, args, sent := kind.connect(t)
			ctx := context.Background()
			call := func(a map[string]any, permit string) string {
				t.Helper()
				raw, _ := json.Marshal(a)
				return mustJSONString(e.s.call(ctx, g, tool, raw, permit))
			}
			mint := func(p fakes.Permit) string {
				if p.Audience == "" {
					p.Audience = g.Audience
				}
				if p.Tenant == "" {
					p.Tenant = tenant
				}
				if p.Tool == "" {
					p.Tool = tool
				}
				if p.Sub == "" {
					p.Sub = g.Request.Agent
				}
				if p.Args == nil {
					p.Args = mustJSON(args(p.Op))
				}
				return e.kiff.Mint(p)
			}
			otherKey := fakes.NewKIFF()
			otherKey.Issuer, otherKey.KID = e.kiff.Issuer, e.kiff.KID
			otherArgs := args("op_a")
			for k, v := range otherArgs {
				if k != IdempotencyArg {
					if _, isNumber := v.(int); isNumber {
						otherArgs[k] = 999
					} else {
						otherArgs[k] = "other_value"
					}
					break
				}
			}

			before := sent()
			for _, c := range []struct{ name, callOp, permit, want string }{
				{"no permit", "op_m", "", "no execution permit"},
				{"other arguments", "op_a", mint(fakes.Permit{Op: "op_a", Args: mustJSON(otherArgs)}), "other arguments"},
				{"other tool", "op_b", mint(fakes.Permit{Op: "op_b", Tool: "another_tool"}), "another tool"},
				{"other grant", "op_c", mint(fakes.Permit{Op: "op_c", Audience: g.Audience + "_other"}), "another audience"},
				{"other account", "op_d", mint(fakes.Permit{Op: "op_d", Tenant: "tnt_other"}), "another KIFF account"},
				{"expired", "op_e", mint(fakes.Permit{Op: "op_e", IssuedAt: time.Now().Add(-2 * time.Minute)}), "expired"},
				{"too long-lived", "op_f", mint(fakes.Permit{Op: "op_f", Lifetime: time.Hour}), "lifetime"},
				{"not KIFF's key", "op_g", otherKey.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: tool, Sub: g.Request.Agent, Op: "op_g", Args: mustJSON(args("op_g"))}), "signature"},
				{"another operation", "op_other", mint(fakes.Permit{Op: "op_h", Args: mustJSON(args("op_other"))}), "operation"},
			} {
				out := call(args(c.callOp), c.permit)
				if !strings.Contains(out, "Refused by planque") || !strings.Contains(out, c.want) {
					t.Errorf("%s: %s, want a refusal naming %q", c.name, out, c.want)
				}
				if sent() != before {
					t.Fatalf("%s: the vendor received a call (%d, was %d)", c.name, sent(), before)
				}
			}

			// A permitted call is sent once; its replay sends nothing more.
			permit := mint(fakes.Permit{Op: "op_ok"})
			first := call(args("op_ok"), permit)
			if strings.Contains(first, "Refused") || sent() != before+1 {
				t.Fatalf("a permitted call: %s, vendor calls %d, want %d", first, sent(), before+1)
			}
			if again := call(args("op_ok"), permit); again != first || sent() != before+1 {
				t.Fatalf("a replay: %s (first %s), vendor calls %d, want %d", again, first, sent(), before+1)
			}
		})
	}
}

// Every kind of connection keeps its grant to its own account: another
// account's key reaches none of the grant's routes and learns nothing of
// it, the grant stays as it was, and a tool the grant doesn't include is
// refused. The vendor receives nothing. (Before the owner approves the
// tools, nothing is called either; each kind's own test covers that:
// TestNothingIsReadBeforeTheToolsAreApproved, TestConnectorCatalogAndApprovalBoundary,
// TestAnAPIKeyGrantEndToEnd.)
func TestConformanceEveryKindKeepsItsGrantToItsAccount(t *testing.T) {
	for _, kind := range conformanceKinds {
		t.Run(kind.name, func(t *testing.T) {
			e, g, tool, args, sent := kind.connect(t)
			before := sent()
			at := "/v1/accounts/" + tenant + "/grants/" + g.ID
			// k_other is a valid key for another account. A route a key may
			// use refuses it for this grant (403/404); a route that takes
			// only the owner's device signature answers 401 to any key. A
			// 401 elsewhere would mean the key wasn't read at all, and the
			// test would pass on a broken fixture.
			for _, c := range []struct {
				method, path string
				deviceOnly   bool
			}{
				{"GET", at, false}, {"GET", at + "/catalog", false}, {"POST", at + "/tools", false}, {"DELETE", at, true},
				{"GET", "/v1/accounts/" + tenant + "/grants", true},
				{"GET", "/v1/accounts/tnt_other/grants/" + g.ID, false},
			} {
				st, out := e.do(c.method, c.path, map[string]any{"tools": []any{map[string]any{"tool": tool}}}, nil, "k_other")
				if c.deviceOnly && st != 401 || !c.deviceOnly && st != 403 && st != 404 {
					t.Errorf("another account's %s %s: %d %v", c.method, c.path, st, out)
				}
				if b, _ := json.Marshal(out); strings.Contains(string(b), g.ID) || strings.Contains(string(b), tool) {
					t.Errorf("another account's %s %s learned about the grant: %s", c.method, c.path, b)
				}
			}
			after, err := e.s.loadGrant(context.Background(), tenant, g.ID)
			if err != nil || after.Status != StatusConnected || len(after.Tools) != len(g.Tools) {
				t.Fatalf("the grant after another account's attempts: %v %+v", err, after)
			}
			// A tool the grant doesn't include, even with a permit for it.
			raw := mustJSON(args("op_x"))
			permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "not_granted", Sub: g.Request.Agent, Op: "op_x", Args: raw})
			if out := mustJSONString(e.s.call(context.Background(), g, "not_granted", raw, permit)); !strings.Contains(out, "Refused by planque") || !strings.Contains(out, "does not include") {
				t.Errorf("a tool the grant doesn't include: %s", out)
			}
			if sent() != before {
				t.Fatalf("the vendor received a call (%d, was %d)", sent(), before)
			}
		})
	}
}
