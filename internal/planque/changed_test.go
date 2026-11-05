package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// changedCall makes a permitted create_refund call on g, as KIFF's gateway
// would after allowing it, and answers planque's result.
func changedCall(t *testing.T, e *env, g *Grant, op string) string {
	t.Helper()
	args := map[string]any{"charge": "ch_9", "amount": 100, IdempotencyArg: op}
	raw, _ := json.Marshal(args)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "create_refund", Sub: g.Request.Agent, Op: op, Args: raw})
	return mustJSONString(e.s.call(context.Background(), g, "create_refund", raw, permit))
}

// changeRefund has the fake vendor answer tools/list with create_refund
// taking another argument than the one the owner approved.
func changeRefund(tools []map[string]any) []map[string]any {
	for _, t := range tools {
		if t["name"] == "create_refund" {
			t["inputSchema"] = map[string]any{"type": "object", "properties": map[string]any{
				"charge": map[string]any{"type": "string"}, "amount": map[string]any{"type": "integer"},
				"refund_to_card": map[string]any{"type": "string"}}, "required": []string{"charge", "amount", "refund_to_card"}}
		}
	}
	return tools
}

// A direct MCP tool the owner approved that changes at the vendor pauses
// its grant at the next call: that call isn't sent, and nor is any later
// one, without reading the vendor again. The grant says why, in its words
// and its events, and keeps the pause through a new sign-in.
func TestAChangedToolPausesItsGrant(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	if out := changedCall(t, e, g, "op_before"); strings.Contains(out, "Refused") || e.vendor.CallCount() != 1 {
		t.Fatalf("a call before any change: %s, vendor calls %d", out, e.vendor.CallCount())
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	if g.DefinitionsReadAt.IsZero() || g.Tools[0].VendorSchemaSHA256 == "" {
		t.Fatalf("the definitions weren't read, or the approved hash kept: %+v", g)
	}

	e.vendor.ChangeTools = changeRefund
	g.DefinitionsReadAt = time.Time{} // as if 10 minutes had passed
	out := changedCall(t, e, g, "op_changed")
	if !strings.Contains(out, "Refused by planque") || !strings.Contains(out, "changed create_refund since you approved it") ||
		!strings.Contains(out, `"code":"`+CodeOperationChanged+`"`) || e.vendor.CallCount() != 1 {
		t.Fatalf("a call after the change: %s, vendor calls %d", out, e.vendor.CallCount())
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	if g.Paused == nil || g.Paused.Tool != "create_refund" || !strings.Contains(statusWords(g), "Paused · ") ||
		!strings.Contains(statusWords(g), "connect again") || !strings.Contains(mustJSONString(g.Events), "paused: ") {
		t.Fatalf("the grant after the change: paused=%+v words=%q", g.Paused, statusWords(g))
	}

	// Paused: refused at once, without reading the vendor's tools again.
	lists := e.vendor.Lists
	if out := changedCall(t, e, g, "op_after"); !strings.Contains(out, "Refused by planque") || e.vendor.CallCount() != 1 || e.vendor.Lists != lists {
		t.Fatalf("a call on the paused grant: %s, vendor calls %d, lists %d (was %d)", out, e.vendor.CallCount(), e.vendor.Lists, lists)
	}
	// The vendor changing it back doesn't lift the pause: the owner approves.
	e.vendor.ChangeTools = nil
	g.DefinitionsReadAt = time.Time{}
	if out := changedCall(t, e, g, "op_back"); !strings.Contains(out, "Refused by planque") || e.vendor.CallCount() != 1 {
		t.Fatalf("a call after the vendor changed it back: %s, vendor calls %d", out, e.vendor.CallCount())
	}
}

// A tool the vendor no longer lists pauses the grant as a changed one does.
func TestARemovedToolPausesItsGrant(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.vendor.ChangeTools = func(tools []map[string]any) []map[string]any {
		out := tools[:0]
		for _, t := range tools {
			if t["name"] != "create_refund" {
				out = append(out, t)
			}
		}
		return out
	}
	if out := changedCall(t, e, g, "op_gone"); !strings.Contains(out, "no longer lists create_refund") || e.vendor.CallCount() != 0 {
		t.Fatalf("a call after the tool was removed: %s, vendor calls %d", out, e.vendor.CallCount())
	}
	if g, _ = e.s.loadGrant(context.Background(), tenant, g.ID); g.Paused == nil {
		t.Fatal("the grant isn't paused")
	}
}

// Unchanged, a grant's tools are read at most every definitionsReadEvery:
// calls within it read nothing more, and one after it reads again.
func TestDefinitionsAreReadAtMostEveryTenMinutes(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	readAgo := func(d time.Duration) {
		t.Helper()
		g, _ = e.s.loadGrant(ctx, tenant, g.ID)
		g.DefinitionsReadAt = time.Now().Add(-d)
		if err := e.s.saveGrant(ctx, g); err != nil {
			t.Fatal(err)
		}
	}
	readAgo(definitionsReadEvery - time.Minute)
	lists := e.vendor.Lists
	if out := changedCall(t, e, g, "op_1"); strings.Contains(out, "Refused") || e.vendor.Lists != lists {
		t.Fatalf("read %v ago: %s, lists %d (was %d)", definitionsReadEvery-time.Minute, out, e.vendor.Lists, lists)
	}
	readAgo(definitionsReadEvery + time.Minute)
	if out := changedCall(t, e, g, "op_2"); strings.Contains(out, "Refused") || e.vendor.Lists != lists+1 {
		t.Fatalf("read %v ago: %s, lists %d (want %d)", definitionsReadEvery+time.Minute, out, e.vendor.Lists, lists+1)
	}
	if g, _ = e.s.loadGrant(ctx, tenant, g.ID); time.Since(g.DefinitionsReadAt) > time.Minute {
		t.Fatalf("the read wasn't recorded: %v", g.DefinitionsReadAt)
	}
	if e.vendor.CallCount() != 2 {
		t.Fatalf("vendor calls %d, want 2", e.vendor.CallCount())
	}
}

// A list the vendor can't give (it signed the grant out) doesn't pause the
// grant: the call follows its usual path, which notices the sign-out.
func TestAnUnreadableListDoesntPause(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.vendor.SignOut()
	changedCall(t, e, g, "op_signed_out")
	if g, _ = e.s.loadGrant(context.Background(), tenant, g.ID); g.Paused != nil || g.Status != StatusNeedsReconnect {
		t.Fatalf("after a sign-out: paused=%+v status=%s, want needs_reconnect and no pause", g.Paused, g.Status)
	}
}

// A grant paused for a change doesn't touch another grant of the account.
func TestAPauseStaysOnItsGrant(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	if code, why := e.s.pauseGrant(context.Background(), g, "create_refund", "a test pause"); code != CodeOperationChanged || why == "" {
		t.Fatalf("pause: %q %q", code, why)
	}
	other := *g
	other.ID, other.Paused, other.Events = g.ID+"_other", nil, nil
	if code, why := pausedRefusal(&other); code != "" || why != "" {
		t.Fatalf("another grant refused: %q %q", code, why)
	}
	if again, _ := e.s.loadGrant(context.Background(), tenant, g.ID); again.Paused == nil || again.Paused.Why != "a test pause" {
		t.Fatalf("the pause wasn't kept: %+v", again.Paused)
	}
}

// A managed connector tool whose definition the provider reports changed
// pauses its grant: the call isn't sent, nor is any later one, and the
// grant says why.
func TestAChangedConnectorToolPausesItsGrant(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	call := func(op string) string {
		raw := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
		p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: g.Request.Agent, Op: op, Args: raw})
		return mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, p))
	}
	f.changed = true
	if out := call("op_changed"); !strings.Contains(out, "Not sent") || !strings.Contains(out, "since you approved it") || f.executions != 0 {
		t.Fatalf("a call after the provider changed the tool: %s, executions %d", out, f.executions)
	}
	g, _ = e.s.loadGrant(ctx, tenant, g.ID)
	if g.Paused == nil || g.Paused.Tool != "reply_ticket" || !strings.Contains(statusWords(g), "Paused · ") {
		t.Fatalf("the grant: paused=%+v words=%q", g.Paused, statusWords(g))
	}
	f.changed = false
	if out := call("op_after"); !strings.Contains(out, "Refused by planque") || !strings.Contains(out, `"code":"`+CodeOperationChanged+`"`) || f.executions != 0 {
		t.Fatalf("a call on the paused grant: %s, executions %d", out, f.executions)
	}
}

// A connector tool definition the provider answers but planque can't read
// or check (an unparseable body, another tool or version) pauses nothing:
// the call isn't sent, says so, and the grant stays as it was. Only a
// definition read and found different from the approved one pauses.
func TestAnUnreadableConnectorDefinitionDoesntPause(t *testing.T) {
	for name, body := range map[string]string{
		"unparseable":     `{"slug":"DESK_DEMO_REPLY","toolkit":`,
		"another version": `{"slug":"DESK_DEMO_REPLY","toolkit":{"slug":"desk_demo"},"version":"20261102_00","name":"Reply to a ticket","input_parameters":{"type":"object"}}`,
		"another tool":    `{"slug":"DESK_DEMO_OTHER","toolkit":{"slug":"desk_demo"},"version":"20261101_00","name":"Other","input_parameters":{"type":"object"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			e, f := newConnectorFixture(t)
			g := connectConnector(t, e, f)
			ctx := context.Background()
			call := func(op string) string {
				raw := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"` + op + `"}`)
				p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: g.Request.Agent, Op: op, Args: raw})
				return mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, p))
			}
			f.toolBody = body
			if out := call("op_unread"); !strings.Contains(out, "Not sent") || strings.Contains(out, "since you approved it") || f.executions != 0 {
				t.Fatalf("a call whose definition can't be read: %s, executions %d", out, f.executions)
			}
			g, _ = e.s.loadGrant(ctx, tenant, g.ID)
			if g.Paused != nil || strings.Contains(statusWords(g), "Paused") {
				t.Fatalf("the grant was paused: %+v %q", g.Paused, statusWords(g))
			}
			// Once the definition reads as approved, the next call is sent.
			f.toolBody = ""
			if out := call("op_next"); strings.Contains(out, "Not sent") || strings.Contains(out, "Refused") || f.executions != 1 {
				t.Fatalf("the next call: %s, executions %d", out, f.executions)
			}
		})
	}
}
