package planque

import (
	"context"
	"strings"
	"testing"
	"time"
)

// A key grant's tools are declared to KIFF when the grant is asked for, so
// a Card can be drafted over them and signed with the key, but nothing is
// usable until the key is in: a call before it is refused before anything
// reaches the vendor, saying it waits for its key. A request cancelled, or
// ended by its 7-day deadline, takes its declarations with it.

// askedKeyGrant asks for an API-key grant to v with its two tools, as Setup
// does before the owner gives the key.
func (e *env) askedKeyGrant(v *restVendor) (id string) {
	e.t.Helper()
	e.enrol()
	st, out := e.askAPIGrant(v, apiTools, nil)
	if st != 201 || out["status"] != StatusPending {
		e.t.Fatalf("ask: %d %v", st, out)
	}
	return out["id"].(string)
}

func (e *env) declared(names ...string) []string {
	e.t.Helper()
	var have []string
	for _, n := range names {
		if _, ok := e.kiff.Connection(tenant, n); ok {
			have = append(have, n)
		}
	}
	return have
}

func TestAKeyGrantsToolCalledBeforeItsKeyIsRefusedAndSendsNothing(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.askedKeyGrant(v)
	if got := e.declared("create_refund", "list_charges"); len(got) != 2 {
		t.Fatalf("declared to KIFF while the grant waits for its key: %v, want both tools", got)
	}
	ctx := context.Background()
	for tool, args := range map[string]map[string]any{
		"list_charges":  {"limit": 2},
		"create_refund": {"charge": "ch_1", "amount": 300},
	} {
		res, err := e.kiff.Call(ctx, "k_taf", tool, args)
		if err == nil && res.Decision == "allowed" {
			text, isErr := resultText(t, res.Result)
			if !isErr {
				t.Errorf("%s ran before its key was given: %s", tool, text)
			} else if !strings.Contains(strings.ToLower(text), "key") {
				t.Errorf("%s refused without saying it waits for its key: %s", tool, text)
			}
		}
	}
	if calls := v.seen(); len(calls) != 0 {
		t.Fatalf("the vendor was called before the key was given: %+v", calls)
	}
}

// Cancelled by the owner, ended by its deadline, or ended with its taf: none of its tools stays
// declared, so no Card can be drafted over a request that's gone.
func TestAKeyRequestThatEndsLeavesNothingDeclared(t *testing.T) {
	for _, how := range []string{"cancelled", "expired", "its taf removed"} {
		t.Run(how, func(t *testing.T) {
			e := newEnv(t)
			v := newRESTVendor(t)
			id := e.askedKeyGrant(v)
			if got := e.declared("create_refund", "list_charges"); len(got) != 2 {
				t.Fatalf("declared before it ended: %v, want both tools", got)
			}
			switch how {
			case "cancelled":
				if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 || live[out["status"].(string)] {
					t.Fatalf("cancel: %d %v", st, out)
				}
			case "its taf removed":
				e.people()
				if st, out := e.endAgent(tenant, "agt_taf", "s_owner", ""); st != 200 {
					t.Fatalf("end the taf's agent: %d %v", st, out)
				}
			case "expired":
				e.later(pendingDeadline + time.Minute)
				if _, err := e.s.Sweep(t.Context()); err != nil {
					t.Fatal(err)
				}
				if s := e.status(id); live[s] {
					t.Fatalf("past its deadline: %s", s)
				}
			}
			if got := e.declared("create_refund", "list_charges"); len(got) != 0 {
				t.Fatalf("still declared to KIFF after the request %s: %v", how, got)
			}
			if calls := v.seen(); len(calls) != 0 {
				t.Fatalf("the vendor was called: %+v", calls)
			}
		})
	}
}
