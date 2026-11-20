package planque

import (
	"context"
	"encoding/json"
	"slices"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/store"
)

// 20 Nov (planque#272): an app approves an app event and planque makes the
// subscription, then the app asks tafs cloud to record it. If that never
// lands (a restart, the dial moving first) the subscription stays live with
// nothing that knows it. A caller that asks for confirmation has a deadline: a
// subscription it never confirms is removed, one it confirms is kept, and one
// made without the option is never touched.

func approveAppEvent(t *testing.T, e *env, f *eventFixture, g *Grant, confirm bool, channel string) string {
	t.Helper()
	f.definition.Slug = "SLACKBOT_REACTION_ADDED"
	f.definition.Config = json.RawMessage(`{"type":"object","properties":{"channel_id":{"type":"string"}},"required":["channel_id"]}`)
	body := map[string]any{"account": g.Account, "slug": "SLACKBOT_REACTION_ADDED", "config": map[string]any{"channel_id": channel}}
	if confirm {
		body["confirm"] = true
	}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/triggers", body, e.phone, "")
	if st != 201 {
		t.Fatalf("approve: %d %v", st, out)
	}
	return out["trigger"].(string)
}

func confirmPath(g *Grant, id string) string {
	return "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers/" + id + "/confirm"
}

func unconfirmedNotes(t *testing.T, e *env) int {
	t.Helper()
	notes, err := e.kv.Query(context.Background(), triggerUnconfirmedPK, "")
	if err != nil {
		t.Fatal(err)
	}
	return len(notes)
}

// Category: authority
// Attack: a live subscription nothing records, kept for as long as the connection lasts
// Refused: removed upstream once, off the connection, after its deadline and not before
func TestAnUnconfirmedAppEventSubscriptionIsRemovedAfterItsDeadline(t *testing.T) {
	e, f, g := newEventFixture(t)
	id := approveAppEvent(t, e, f, g, true, "CDEMO123")
	if unconfirmedNotes(t, e) != 1 {
		t.Fatal("the subscription's deadline wasn't noted")
	}
	f.now = f.now.Add(confirmWithin - time.Second)
	if n, err := e.s.sweepUnconfirmed(context.Background()); err != nil || n != 0 || len(f.deleted) != 0 {
		t.Fatalf("swept before its deadline: %d %v deleted %v", n, err, f.deleted)
	}
	f.now = f.now.Add(2 * time.Second)
	if n, err := e.s.sweepUnconfirmed(context.Background()); err != nil || n != 1 || !slices.Equal(f.deleted, []string{id}) {
		t.Fatalf("sweep after the deadline: %d %v deleted %v", n, err, f.deleted)
	}
	var record triggerRecord
	if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", id), "-", &record); err != nil || record.Active || record.Removing {
		t.Fatalf("the removed subscription still routes, or is stuck removing: %+v %v", record, err)
	}
	saved, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if slices.Contains(saved.TriggerIDs, id) {
		t.Fatal("the removed subscription stayed on its connection")
	}
	if n, err := e.s.sweepUnconfirmed(context.Background()); err != nil || n != 0 || len(f.deleted) != 1 {
		t.Fatalf("a second sweep deleted again: %d %v deleted %v", n, err, f.deleted)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 0 || unconfirmedNotes(t, e) != 0 {
		t.Fatal("a settled removal left a note")
	}
}

// Category: authority
// Attack: a good subscription removed because its caller was slow, or one made without the option swept
// Refused: confirmed in time, or confirmed late but before the sweep, it stays; with no option it is never swept
func TestAConfirmedOrUnaskedSubscriptionIsNeverSwept(t *testing.T) {
	e, f, g := newEventFixture(t)
	confirmed := approveAppEvent(t, e, f, g, true, "CONFIRMED1")
	late := approveAppEvent(t, e, f, g, true, "LATE0001")
	unasked := approveAppEvent(t, e, f, g, false, "UNASKED1")
	if unconfirmedNotes(t, e) != 2 {
		t.Fatalf("notes: %d, want the two that asked for confirmation", unconfirmedNotes(t, e))
	}
	if st, out := e.do("POST", confirmPath(g, confirmed), nil, nil, "k_taf"); st != 200 || out["confirmed"] != true {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if unconfirmedNotes(t, e) != 1 {
		t.Fatal("a confirmed subscription's deadline note stayed")
	}
	f.now = f.now.Add(confirmWithin + time.Minute) // past the deadline, before any sweep
	if st, out := e.do("POST", confirmPath(g, late), nil, nil, "k_taf"); st != 200 {
		t.Fatalf("a late confirm before the sweep: %d %v", st, out)
	}
	if n, err := e.s.sweepUnconfirmed(context.Background()); err != nil || n != 0 || len(f.deleted) != 0 {
		t.Fatalf("the sweep removed a confirmed or unasked subscription: %d %v deleted %v", n, err, f.deleted)
	}
	saved, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	for _, id := range []string{confirmed, late, unasked} {
		var record triggerRecord
		if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", id), "-", &record); err != nil || !record.Active || !slices.Contains(saved.TriggerIDs, id) {
			t.Fatalf("%s no longer routes: %+v %v", id, record, err)
		}
	}
}

// A confirmation after the removal is refused, so the caller approves again;
// another agent can't confirm; a door has nothing to confirm.
// Category: authority
// Attack: confirming a subscription that is gone, or someone else's
// Refused: 409, 403, 400
func TestConfirmingARemovedOrForeignSubscriptionIsRefused(t *testing.T) {
	e, f, g := newEventFixture(t)
	e.kiff.Keys["k_stranger"] = kiff.Me{APIKey: true, TenantID: tenant, AgentID: "agt_stranger"}
	id := approveAppEvent(t, e, f, g, true, "CDEMO123")
	if st, _ := e.do("POST", confirmPath(g, id), nil, nil, "k_stranger"); st != 403 {
		t.Fatalf("another agent confirmed: %d", st)
	}
	if st, _ := e.do("POST", confirmPath(g, "trigger_other"), nil, nil, "k_taf"); st != 404 {
		t.Fatalf("a trigger the connection doesn't hold: %d", st)
	}
	door := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CONE0001"}, "confirm": true}
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/triggers", door, e.phone, ""); st != 400 {
		t.Fatalf("a door asked for confirmation: %d", st)
	}
	f.now = f.now.Add(confirmWithin + time.Second)
	if n, err := e.s.sweepUnconfirmed(context.Background()); err != nil || n != 1 {
		t.Fatalf("sweep: %d %v", n, err)
	}
	if st, _ := e.do("POST", confirmPath(g, id), nil, nil, "k_taf"); st != 404 && st != 409 {
		t.Fatalf("confirmed a removed subscription: %d", st)
	}
}

// The sweep has read an unconfirmed note and noted the removal, and the caller
// confirms before the removal claims the record: the confirmation moved the
// record to a newer version, so the removal drops its note and deletes nothing.
// Category: authority
// Attack: a confirmation racing the sweep, so a good subscription is deleted upstream
// Refused: nothing deleted, the subscription stays on
func TestAConfirmationBeatsASweepThatAlreadyReadItsNote(t *testing.T) {
	e, f, g := newEventFixture(t)
	id := approveAppEvent(t, e, f, g, true, "CDEMO123")
	notes, err := e.kv.Query(context.Background(), triggerUnconfirmedPK, "")
	if err != nil || len(notes) != 1 {
		t.Fatalf("notes: %v %v", notes, err)
	}
	var note triggerUnconfirmed
	if err := json.Unmarshal(notes[0].Doc, &note); err != nil {
		t.Fatal(err)
	}
	cleanup := triggerCleanup{Provider: note.Provider, Trigger: note.Trigger, Tenant: note.Tenant, Grant: note.Grant, Fence: note.Fence}
	if err := putDoc(context.Background(), e.kv, store.Item{PK: triggerCleanupPK, SK: cleanup.Provider + "#" + cleanup.Trigger}, cleanup, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	// The caller confirms now: after the sweep read the note, before it claims the record.
	if st, out := e.do("POST", confirmPath(g, id), nil, nil, "k_taf"); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	p := e.s.Connectors["composio"].(connector.EventProvider)
	if done, err := e.s.removeTrigger(context.Background(), p, cleanup); err != nil || !done || len(f.deleted) != 0 {
		t.Fatalf("the sweep's removal after the confirmation: done %v err %v deleted %v", done, err, f.deleted)
	}
	var record triggerRecord
	if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", id), "-", &record); err != nil || !record.Active || record.Removing {
		t.Fatalf("the confirmed subscription was taken: %+v %v", record, err)
	}
}
