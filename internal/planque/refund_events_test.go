package planque

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/store"
)

// refundCaseProfile is a desk-and-payments profile whose refund operation
// has a case block naming the app's refund event.
const refundCaseProfile = `{"v":"kiff.profile.v1","vendor":"Payments demo","transport":"rest","pin":{"kind":"none"},"source":"drafted",
 "operations":[
  {"tool":"get_account","title":"Look up the account","kind":"read"},
  {"tool":"get_charge","title":"Look up a charge","kind":"read","resource":{"arg":"charge","noun":"charge"}},
  {"tool":"list_refunds","title":"List a charge's refunds","kind":"read","resource":{"arg":"charge","noun":"charge"}},
  {"tool":"refund","title":"Refund a charge","kind":"write","effect":"moves_money","shape":"single",
   "resource":{"arg":"charge","noun":"charge"},
   "amount":{"from":"argument","arg":"amount","format":"integer_minor","scale":2,"unit":"EUR"},
   "marker":{"field":"metadata.kiff_op"},
   "case":{"kind":"refund","charge_arg":"charge",
    "account_read":{"tool":"get_account","path":"id"},
    "charge_read":{"tool":"get_charge","args":{"charge":"$args.charge"},"currency":"currency","captured":"amount_captured","disputed":"disputed"},
    "refunds_read":{"tool":"list_refunds","args":{"charge":"$args.charge"},"list":"data","id":"id","amount":"amount","status":"status","marker":"metadata.kiff_op"},
    "refund_event":{"event":"PAYMENTS_DEMO_REFUND_UPDATED","object":"refund","charge":"charge"}}}]}`

type recordingSink struct {
	mu  sync.Mutex
	got []door.AppEvent
}

func (r *recordingSink) DeliverEvent(_ context.Context, e door.AppEvent) (door.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.got = append(r.got, e)
	return door.Result{Taken: true, Status: 202}, nil
}

func (r *recordingSink) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.got)
}

// The refund event a profile's case block names goes to KIFF as well as to
// the taf, once, under an id of its own; another event of the app goes to
// the taf only, and a grant whose profile names no refund event hands KIFF
// nothing.
func TestARefundEventGoesToKIFFOnce(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnectorWith(t, e, f, completesProfile())
	ctx := context.Background()
	g.Request.Profile, g.Connected, g.ConnectedProfile = json.RawMessage(refundCaseProfile), nil, nil
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	taf, kiff := &recordingSink{}, &recordingSink{}
	e.s.EventSink, e.s.KIFFEventSink = taf, kiff
	route := triggerRecord{Tenant: tenant, Grant: g.ID, Provider: "composio", Account: g.Account, ID: "trigger_refunds",
		Definition: f.definitionOf("PAYMENTS_DEMO_REFUND_UPDATED"), Active: true, AppEvent: true}
	if err := putDoc(ctx, e.kv, store.Item{PK: triggerPK(route.Provider, route.ID), SK: "-", Ver: 1}, route, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	refund := eventOf("PAYMENTS_DEMO_REFUND_UPDATED", map[string]any{"refund": map[string]any{"id": "re_1", "amount": 2000, "status": "failed",
		"charge": "ch_1", "metadata": map[string]any{"kiff_op": "tc-1"}}})
	app, ok := appEvent(g, route, refund)
	if !ok {
		t.Fatal("not an app event")
	}
	ev, ok := kiffEvent(g, route, refund)
	if !ok || ev.ID == app.ID || !strings.HasPrefix(ev.ID, "evk_") || ev.Event != app.Event || string(ev.Payload) != string(app.Payload) || ev.Key != "" {
		t.Fatalf("KIFF's event: %+v (the taf's id %s)", ev, app.ID)
	}
	// The provider's redelivery is the same event.
	if again, _ := kiffEvent(g, route, refund); again.ID != ev.ID {
		t.Fatal("a redelivery got another id")
	}
	for i := 0; i < 2; i++ {
		b, _ := json.Marshal(app)
		if _, err := e.s.enqueueDoor(ctx, g, route.Provider, route.ID, app.ID, "app_event", b); err != nil {
			t.Fatal(err)
		}
		b, _ = json.Marshal(ev)
		if _, err := e.s.enqueueDoor(ctx, g, route.Provider, route.ID, ev.ID, kindKIFFEvent, b); err != nil {
			t.Fatal(err)
		}
	}
	if taf.count() != 1 || kiff.count() != 1 || kiff.got[0].ID != ev.ID || kiff.got[0].Tenant != tenant || kiff.got[0].Grant != g.ID {
		t.Fatalf("delivered: taf %d, KIFF %+v", taf.count(), kiff.got)
	}
	// Another event of the app is the taf's only.
	other := route
	other.Definition = f.definitionOf("PAYMENTS_DEMO_CHARGE_SUCCEEDED")
	if _, ok := kiffEvent(g, other, eventOf("PAYMENTS_DEMO_CHARGE_SUCCEEDED", map[string]any{"id": "ch_2"})); ok {
		t.Fatal("an event no case block names went to KIFF")
	}
	// A profile with no case block names none.
	g.Request.Profile = mustJSON(completesProfile())
	if _, ok := kiffEvent(g, route, refund); ok {
		t.Fatal("a grant with no case block handed KIFF an event")
	}
	// And a KIFF event kept before the profile changed is dropped, not sent.
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	stale := ev
	stale.ID = "evk_stale"
	b, _ := json.Marshal(stale)
	if _, err := e.s.enqueueDoor(ctx, g, route.Provider, route.ID, stale.ID, kindKIFFEvent, b); err != nil {
		t.Fatal(err)
	}
	if kiff.count() != 1 {
		t.Fatal("a refund event no case block names now was handed to KIFF")
	}
}

// The case block's reads are renamed to KIFF's tool names with the rest.
func TestTheCaseReadsAreRenamedForKIFF(t *testing.T) {
	before := []string{"get_account", "get_charge", "list_refunds", "refund"}
	tools := []ToolRequest{{Name: "acct_1a2b"}, {Name: "charge_3c4d"}, {Name: "refunds_5e6f"}, {Name: "refund_7a8b"}}
	raw, err := renameProfile(json.RawMessage(refundCaseProfile), before, tools)
	if err != nil {
		t.Fatal(err)
	}
	var p struct {
		Operations []struct {
			Tool string `json:"tool"`
			Case *struct {
				AccountRead struct{ Tool string } `json:"account_read"`
				ChargeRead  struct{ Tool string } `json:"charge_read"`
				RefundsRead struct{ Tool string } `json:"refunds_read"`
			} `json:"case"`
		} `json:"operations"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	c := p.Operations[3].Case
	if p.Operations[3].Tool != "refund_7a8b" || c == nil || c.AccountRead.Tool != "acct_1a2b" || c.ChargeRead.Tool != "charge_3c4d" || c.RefundsRead.Tool != "refunds_5e6f" {
		t.Fatalf("renamed: %s", raw)
	}
}
