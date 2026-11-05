package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/store"
)

// completesProfile is the connector fixture's profile with reply_ticket
// completing later, on the desk's "new ticket" event.
func completesProfile() profile.Profile {
	return profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "reply_ticket", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single",
			Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"},
			Completes: &profile.Completes{Event: "DESK_DEMO_NEW_TICKET", Key: "id", EventKey: "ticket.reply", TimeoutS: 900,
				Words: "The desk answers when the customer replies; the taf waits for that."}}}}
}

// A call that completes later is refused before sending while the taf
// doesn't listen for its event; once subscribed, the answer says pending
// with the handle and the deadline, a replay answers the same, and the
// event that completes it is delivered with that handle as its key.
func TestACallThatCompletesLaterIsPendingAndItsEventCarriesTheKey(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnectorWith(t, e, f, completesProfile())
	ctx := context.Background()
	call := func(op string) (string, map[string]any) {
		args := mustJSON(map[string]any{"ticket": "ticket_demo", "text": "hello", IdempotencyArg: op})
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: g.Request.Agent, Op: op, Args: args})
		res := e.s.call(ctx, g, "reply_ticket", args, permit)
		_, sc := textOf(res)
		return allText(res), sc
	}
	// Not listening: refused before anything is sent.
	out, sc := call("op_unheard")
	if code(sc) != CodeNotListening || sc["sent"] != false || f.executions != 0 || !strings.Contains(out, "DESK_DEMO_NEW_TICKET") {
		t.Fatalf("not listening: %s (executions %d)", out, f.executions)
	}
	// The owner's subscription to the event, as subscribe keeps it (the
	// subscribing itself is the app-event tests').
	rec := triggerRecord{Tenant: tenant, Grant: g.ID, Provider: "composio", Account: g.Account, ID: "trigger_tickets", Definition: f.definitionOf("DESK_DEMO_NEW_TICKET"), Active: true, AppEvent: true}
	if err := putDoc(ctx, e.kv, store.Item{PK: triggerPK(rec.Provider, rec.ID), SK: "-", Ver: 1}, rec, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	if err := putDoc(ctx, e.kv, store.Item{PK: tenantPK(tenant), SK: "TRIGGER#" + rec.Provider + "#" + rec.ID}, map[string]string{"provider": rec.Provider, "trigger": rec.ID}, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	// Another tenant's subscription to the same event listens for nobody here.
	other := rec
	other.Tenant, other.Grant, other.ID = "tnt_other", "grt_other", "trigger_other"
	if err := putDoc(ctx, e.kv, store.Item{PK: triggerPK(other.Provider, other.ID), SK: "-", Ver: 1}, other, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	out, sc = call("op_heard")
	if code(sc) != "" || f.executions != 1 || sc["pending"] != true || sc["sent"] != true || sc["timeout_s"] != float64(900) {
		t.Fatalf("a call that completes later: %s (executions %d)", out, f.executions)
	}
	wait, _ := sc["wait_for"].(map[string]any)
	if wait["event"] != "DESK_DEMO_NEW_TICKET" || wait["key"] != "reply_demo" || wait["grant"] != g.ID || wait["words"] == "" || !strings.Contains(out, "the run waits up to 15 min") || !strings.Contains(out, "The desk answers when the customer replies") {
		t.Fatalf("the pending answer: %s", out)
	}
	// A replay answers the same, sending nothing more.
	if again, _ := call("op_heard"); again != out || f.executions != 1 {
		t.Fatalf("replay: %s (executions %d)", again, f.executions)
	}
	// The event that completes it carries the handle as its key.
	ev, ok := appEvent(g, triggerRecord{AppEvent: true, Provider: "composio", ID: "trigger_1", Definition: f.definitionOf("DESK_DEMO_NEW_TICKET")},
		eventOf("DESK_DEMO_NEW_TICKET", map[string]any{"ticket": map[string]any{"id": "ticket_9", "reply": "reply_demo"}}))
	if !ok || ev.Key != "reply_demo" || ev.Event != "DESK_DEMO_NEW_TICKET" {
		t.Fatalf("the event's key: %+v", ev)
	}
	// Another event of the app, or one with no handle, carries none.
	if ev, ok := appEvent(g, triggerRecord{AppEvent: true, Provider: "composio", ID: "trigger_1", Definition: f.definitionOf("DESK_DEMO_OTHER")},
		eventOf("DESK_DEMO_OTHER", map[string]any{"ticket": map[string]any{"reply": "reply_demo"}})); !ok || ev.Key != "" {
		t.Fatalf("another event's key: %+v", ev)
	}
	if ev, ok := appEvent(g, triggerRecord{AppEvent: true, Provider: "composio", ID: "trigger_1", Definition: f.definitionOf("DESK_DEMO_NEW_TICKET")},
		eventOf("DESK_DEMO_NEW_TICKET", map[string]any{"ticket": map[string]any{"id": "ticket_9"}})); !ok || ev.Key != "" {
		t.Fatalf("an event with no handle: %+v", ev)
	}
}

// A profile's completes is checked like every other term.
func TestProfileCompletesIsChecked(t *testing.T) {
	ok := completesProfile()
	if err := ok.Validate(); err != nil {
		t.Fatal(err)
	}
	for name, change := range map[string]func(c *profile.Completes){
		"no event":        func(c *profile.Completes) { c.Event = "" },
		"a bad key path":  func(c *profile.Completes) { c.Key = "data[id" },
		"no event key":    func(c *profile.Completes) { c.EventKey = "" },
		"too short":       func(c *profile.Completes) { c.TimeoutS = 30 },
		"over three days": func(c *profile.Completes) { c.TimeoutS = profile.MaxTimeout + 1 },
		"no words":        func(c *profile.Completes) { c.Words = " " },
		"two lines":       func(c *profile.Completes) { c.Words = "a\nb" },
	} {
		p := completesProfile()
		c := *p.Operations[0].Completes
		change(&c)
		p.Operations[0].Completes = &c
		if err := p.Validate(); err == nil {
			t.Errorf("%s passed", name)
		}
	}
	// A read never completes later.
	p := completesProfile()
	p.Operations[0].Kind, p.Operations[0].Effect, p.Operations[0].Resource, p.Operations[0].Shape = "read", "", nil, ""
	if err := p.Validate(); err == nil || !strings.Contains(err.Error(), "completes") {
		t.Errorf("a read with completes: %v", err)
	}
	// The pending answer: a result without the handle is left as it was.
	res := pending(toolResult(`{"other":1}`, false), ok.Operations[0].Completes, "grt_x")
	_, sc := textOf(res)
	if text := allText(res); sc["pending"] != nil || !strings.Contains(text, "carries no id") || !strings.HasPrefix(text, `{"other":1}`) {
		t.Errorf("a result with no handle: %s %v", text, sc)
	}
	// An error result is never pending.
	res = pending(toolResult("the desk refused it", true), ok.Operations[0].Completes, "grt_x")
	if _, sc := textOf(res); sc["pending"] != nil {
		t.Errorf("an error result: %v", sc)
	}
	// A path into a list and a number both read as the handle.
	if got := atPath(map[string]any{"items": []any{map[string]any{"id": json.Number("42")}}}, "items[0].id"); got != "42" {
		t.Errorf("atPath: %q", got)
	}
}

// A 64-bit numeric handle keeps its digits through the pending answer
// and the event, so the two match; the answer's other numbers reach the
// agent as the vendor wrote them (structured content and text alike).
func TestANumericHandleKeepsItsDigits(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnectorWith(t, e, f, completesProfile())
	c := completesProfile().Operations[0].Completes
	const big = "9007199254740993"
	for name, res := range map[string]any{
		"structured": json.RawMessage(`{"content":[{"type":"text","text":"ok"}],"structuredContent":{"id":` + big + `,"amount":12345678901234567,"rate":1.5}}`),
		"text":       toolResult(`{"id":`+big+`,"amount":12345678901234567,"rate":1.5}`, false),
	} {
		got := pending(res, c, g.ID)
		_, sc := textOf(got)
		out := allText(got)
		wait, _ := sc["wait_for"].(map[string]any)
		if sc["pending"] != true || wait["key"] != big {
			t.Errorf("%s: the handle: %v", name, sc)
		}
		raw, _ := json.Marshal(got)
		if !strings.Contains(string(raw), big) || !strings.Contains(string(raw), "12345678901234567") || !strings.Contains(string(raw), "1.5") || strings.Contains(string(raw), "9007199254740992") {
			t.Errorf("%s: the answer's numbers: %s (%s)", name, raw, out)
		}
		// The vendor's own text stays first in the answer, before the pending line.
		if !strings.HasPrefix(out, "{") && !strings.HasPrefix(out, "ok") {
			t.Errorf("%s: the vendor's text was dropped: %s", name, out)
		}
	}
	ev, ok := appEvent(g, triggerRecord{AppEvent: true, Provider: "composio", ID: "trigger_1", Definition: f.definitionOf("DESK_DEMO_NEW_TICKET")},
		eventOf("DESK_DEMO_NEW_TICKET", map[string]any{"ticket": map[string]any{"id": "ticket_9", "reply": json.Number(big)}}))
	if !ok || ev.Key != big {
		t.Fatalf("the event's key: %+v", ev)
	}
}

// allText is every text content of a result, in order, joined by a line.
func allText(res any) string {
	var m struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	_ = json.Unmarshal(mustJSON(res), &m)
	var out []string
	for _, c := range m.Content {
		out = append(out, c.Text)
	}
	return strings.Join(out, "\n")
}

// definitionOf is a trigger type of the fixture's app with that slug.
func (f *connectorFixture) definitionOf(slug string) connector.TriggerType {
	return connector.TriggerType{Slug: slug, Toolkit: "desk_demo", Version: "20261101_00", Name: slug, Config: json.RawMessage(`{"type":"object"}`), Payload: json.RawMessage(`{"type":"object"}`)}
}

// eventOf is a provider event of that slug with that data.
func eventOf(slug string, data map[string]any) connector.Event {
	raw, _ := json.Marshal(data)
	e := connector.Event{ID: "upstream_" + slug, Data: raw, Timestamp: "2026-11-05T10:00:00Z"}
	e.Metadata.Slug = slug
	return e
}
