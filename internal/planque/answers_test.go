package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/store"
)

func (e *env) answerEntries() int {
	e.t.Helper()
	items, err := e.kv.Query(e.t.Context(), pkAnswerDue, "")
	if err != nil {
		e.t.Fatal(err)
	}
	return len(items)
}

func (e *env) opRecord(g *Grant, op string) operation {
	e.t.Helper()
	var rec operation
	if _, err := getDoc(e.t.Context(), e.kv, tenantPK(tenant), skOp(g.ID, op), &rec); err != nil {
		e.t.Fatal(err)
	}
	return rec
}

// A vendor's answer is kept 24 hours for a replay, then dropped: the
// record keeps the outcome alone, and the call is still never sent again.
func TestAVendorsAnswerIsKeptADayThenOnlyItsOutcome(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_kept"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_kept", Args: args})
	first := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if f.executions != 1 || !strings.Contains(first, "reply_demo") {
		t.Fatalf("call: %s", first)
	}
	rec := e.opRecord(g, "op_kept")
	if rec.State != OpDone || rec.ResultUntil.IsZero() || rec.ResultUntil.Sub(time.Now()) > answerTTL+time.Minute || e.answerEntries() != 1 {
		t.Fatalf("answer kept without its expiry indexed: until=%v entries=%d", rec.ResultUntil, e.answerEntries())
	}
	// Within the day: the sweep leaves it, and a replay gets the answer.
	e.at(23 * time.Hour)
	if n, err := e.s.SweepAnswers(ctx); n != 0 || err != nil || e.answerEntries() != 1 {
		t.Fatalf("an answer dropped early: %d %v", n, err)
	}
	if again := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); again != first {
		t.Fatalf("replay within the day: %s", again)
	}
	// Past it, before any sweep: a replay gets the outcome, not the answer.
	e.at(25 * time.Hour)
	late := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if strings.Contains(late, "reply_demo") || !strings.Contains(late, "Answered once") || f.executions != 1 {
		t.Fatalf("an expired answer was served: %s", late)
	}
	// The sweep drops it from the record, and its entry.
	if n, err := e.s.SweepAnswers(ctx); n != 1 || err != nil {
		t.Fatalf("sweep: %d %v", n, err)
	}
	rec = e.opRecord(g, "op_kept")
	if strings.Contains(string(rec.Result), "reply_demo") || rec.State != OpDone || !rec.ResultUntil.IsZero() || e.answerEntries() != 0 {
		t.Fatalf("answer still kept: %s state=%s entries=%d", rec.Result, rec.State, e.answerEntries())
	}
	if again := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); f.executions != 1 || strings.Contains(again, "reply_demo") || !strings.Contains(again, "Answered once") {
		t.Fatalf("after the sweep, a replay sent again or got the answer: %s sends=%d", again, f.executions)
	}
}

// planque's own words (a refusal, a lost answer) carry no vendor data, so
// nothing about them expires.
func TestARefusalHasNoAnswerToExpire(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	f.accountStatus = "EXPIRED"
	raw := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_refused"}`)
	p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_refused", Args: raw})
	if res := mustJSONString(e.s.call(ctx, g, "reply_ticket", raw, p)); !strings.Contains(res, "Not sent") {
		t.Fatalf("not refused: %s", res)
	}
	if rec := e.opRecord(g, "op_refused"); rec.State != OpNotSent || !rec.ResultUntil.IsZero() || e.answerEntries() != 0 {
		t.Fatalf("a refusal was given an expiry: %+v", rec)
	}
}

// A sweep that can't save the record keeps the entry, and the next one
// drops the answer.
func TestASweepThatCantSaveAnAnswerTriesAgain(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_retry"}`)
	p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_retry", Args: args})
	e.s.call(ctx, g, "reply_ticket", args, p)
	kv := e.faulty()
	kv.set(nil, func(it store.Item) bool { return strings.HasPrefix(it.SK, "OP#") })
	e.at(25 * time.Hour)
	if n, err := e.s.SweepAnswers(ctx); n != 0 || err == nil || e.answerEntries() != 1 {
		t.Fatalf("entry dropped though the record wasn't saved: %d %v", n, err)
	}
	kv.set(nil, nil)
	if n, err := e.s.SweepAnswers(ctx); n != 1 || err != nil || e.answerEntries() != 0 {
		t.Fatalf("second sweep: %d %v", n, err)
	}
	if rec := e.opRecord(g, "op_retry"); strings.Contains(string(rec.Result), "reply_demo") {
		t.Fatal("answer still kept")
	}
}

// The one-time clean-up the owner approved: an answer kept before the
// 24 hours existed is dropped once it's a day old, and a younger one gets
// its expiry. Every record keeps its outcome and arguments' hash, so a
// replay is still never sent; it runs once only.
func TestTheOneTimeCleanUpDropsOldAnswersAndKeepsOutcomes(t *testing.T) {
	e, _ := newConnectorFixture(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seed := func(sk, state string, started time.Time, result string) {
		op := operation{Tool: "reply_ticket", ArgsHash: "hash_" + sk, State: state, StartedAt: started, Result: json.RawMessage(result)}
		if err := putDoc(ctx, e.kv, store.Item{PK: tenantPK(tenant), SK: sk, Ver: 2}, op, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	seed("OP#g#old", OpDone, now.Add(-72*time.Hour), `{"content":[{"type":"text","text":"customer_secret"}],"isError":false}`)
	seed("OP#g#young", OpDone, now.Add(-2*time.Hour), `{"content":[{"type":"text","text":"young_secret"}],"isError":false}`)
	seed("OP#g#refused", OpNotSent, now.Add(-72*time.Hour), `{"content":[{"type":"text","text":"Refused by planque"}],"isError":true}`)
	dropped, indexed, err := e.s.CleanStoredAnswers(ctx)
	if err != nil || dropped != 1 || indexed != 1 {
		t.Fatalf("clean-up: dropped=%d indexed=%d err=%v", dropped, indexed, err)
	}
	read := func(sk string) operation {
		var op operation
		if _, err := getDoc(ctx, e.kv, tenantPK(tenant), sk, &op); err != nil {
			t.Fatal(err)
		}
		return op
	}
	old := read("OP#g#old")
	if strings.Contains(string(old.Result), "customer_secret") || old.State != OpDone || old.ArgsHash != "hash_OP#g#old" || !strings.Contains(string(old.Result), "Answered once") {
		t.Fatalf("old record: %+v", old)
	}
	young := read("OP#g#young")
	if !strings.Contains(string(young.Result), "young_secret") || young.ResultUntil.Sub(now.Add(22*time.Hour)).Abs() > time.Second || e.answerEntries() != 1 {
		t.Fatalf("young record: until=%v entries=%d", young.ResultUntil, e.answerEntries())
	}
	if r := read("OP#g#refused"); r.State != OpNotSent || !strings.Contains(string(r.Result), "Refused by planque") {
		t.Fatalf("a refusal changed: %+v", r)
	}
	// Once only: a second run changes nothing, even with a new old answer.
	seed("OP#g#later", OpDone, now.Add(-72*time.Hour), `{"content":[{"type":"text","text":"later"}],"isError":false}`)
	if d, i, err := e.s.CleanStoredAnswers(ctx); d != 0 || i != 0 || err != nil {
		t.Fatalf("ran twice: %d %d %v", d, i, err)
	}
	// The young answer goes with the hourly sweep, a day after its call.
	e.at(23 * time.Hour)
	if n, err := e.s.SweepAnswers(ctx); n != 1 || err != nil || strings.Contains(string(read("OP#g#young").Result), "young_secret") {
		t.Fatalf("sweep of the young answer: %d %v", n, err)
	}
}

// racingKV bumps a record's version between the clean-up's Scan and its
// save, as a write landing in between would (kiff-qa's review of #219).
type racingKV struct {
	store.KV
	bump func()
}

func (r *racingKV) ScanPage(ctx context.Context, p, start string, limit int) ([]store.Item, string, error) {
	items, next, err := r.KV.ScanPage(ctx, p, start, limit)
	if r.bump != nil {
		r.bump()
		r.bump = nil
	}
	return items, next, err
}

// A record the clean-up couldn't save, old or young, isn't left with its
// answer: the marker waits, the next run reads it again, and the answer
// goes.
func TestTheCleanUpRunsAgainForARecordWrittenInBetween(t *testing.T) {
	for _, age := range []time.Duration{72 * time.Hour, 2 * time.Hour} {
		e, _ := newConnectorFixture(t)
		ctx := context.Background()
		op := operation{Tool: "reply_ticket", ArgsHash: "h", State: OpDone, StartedAt: time.Now().UTC().Add(-age),
			Result: json.RawMessage(`{"content":[{"type":"text","text":"customer_secret"}],"isError":false}`)}
		key := store.Item{PK: tenantPK(tenant), SK: "OP#g#raced"}
		put := func(ver int64) {
			it := key
			it.Ver = ver
			if err := putDoc(ctx, e.kv, it, op, store.Cond{}); err != nil {
				t.Fatal(err)
			}
		}
		put(2)
		e.s.KV = &racingKV{KV: e.kv, bump: func() { put(3) }}
		if _, _, err := e.s.CleanStoredAnswers(ctx); err != nil {
			t.Fatal(err)
		}
		if _, ok, _ := e.kv.Get(ctx, onceAnswers.PK, onceAnswers.SK); ok {
			t.Fatalf("age %v: marked done though a record was skipped", age)
		}
		e.at(48 * time.Hour)
		if _, _, err := e.s.CleanStoredAnswers(ctx); err != nil {
			t.Fatal(err)
		}
		if _, err := e.s.SweepAnswers(ctx); err != nil {
			t.Fatal(err)
		}
		var rec operation
		if _, err := getDoc(ctx, e.kv, key.PK, key.SK, &rec); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(rec.Result), "customer_secret") || rec.State != OpDone {
			t.Fatalf("age %v: a raced record keeps the vendor's answer: %s", age, rec.Result)
		}
		if _, ok, _ := e.kv.Get(ctx, onceAnswers.PK, onceAnswers.SK); !ok {
			t.Fatalf("age %v: the second run didn't finish", age)
		}
	}
}

// A table too big for one sweep: each sweep reads what its budget allows,
// saves its cursor, and the next carries on, until every old answer is gone
// and the clean-up is marked done.
func TestTheCleanUpCarriesOnAcrossSweeps(t *testing.T) {
	page, budget := cleanPage, cleanBudget
	cleanPage, cleanBudget = 1, 0
	t.Cleanup(func() { cleanPage, cleanBudget = page, budget })
	e, _ := newConnectorFixture(t)
	ctx := context.Background()
	for _, sk := range []string{"OP#g#a", "OP#g#b", "OP#g#c"} {
		op := operation{Tool: "reply_ticket", ArgsHash: "h", State: OpDone, StartedAt: time.Now().UTC().Add(-72 * time.Hour),
			Result: json.RawMessage(`{"content":[{"type":"text","text":"customer_secret"}],"isError":false}`)}
		if err := putDoc(ctx, e.kv, store.Item{PK: tenantPK(tenant), SK: sk, Ver: 1}, op, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	dropped := 0
	for run := 1; run <= 10; run++ {
		d, _, err := e.s.CleanStoredAnswers(ctx)
		if err != nil {
			t.Fatal(err)
		}
		dropped += d
		if _, ok, _ := e.kv.Get(ctx, onceAnswers.PK, onceAnswers.SK); ok {
			if run < 3 {
				t.Fatalf("marked done after %d runs of one record each", run)
			}
			break
		}
		if run == 10 {
			t.Fatal("never finished")
		}
	}
	if dropped != 3 {
		t.Fatalf("dropped %d answers across the runs, want 3", dropped)
	}
	for _, sk := range []string{"OP#g#a", "OP#g#b", "OP#g#c"} {
		var rec operation
		if _, err := getDoc(ctx, e.kv, tenantPK(tenant), sk, &rec); err != nil || strings.Contains(string(rec.Result), "customer_secret") {
			t.Fatalf("%s still holds its answer: %v", sk, err)
		}
	}
	if _, ok, _ := e.kv.Get(ctx, cleanProgress.PK, cleanProgress.SK); ok {
		t.Fatal("the cursor outlived the clean-up")
	}
}

// throttledKV answers its first ScanPage as DynamoDB does when a read
// stayed throttled past its retries.
type throttledKV struct {
	store.KV
	left int
}

func (k *throttledKV) ScanPage(ctx context.Context, p, start string, limit int) ([]store.Item, string, error) {
	if k.left > 0 {
		k.left--
		return nil, "", &awsapi.Error{Status: 400, Type: "ThrottlingException", Message: "Throughput exceeds the current capacity"}
	}
	return k.KV.ScanPage(ctx, p, start, limit)
}

// A sweep whose scan stayed throttled changes nothing and isn't marked
// done; the next one finishes and writes the marker.
func TestAThrottledCleanUpFinishesAtTheNextSweep(t *testing.T) {
	e, _ := newConnectorFixture(t)
	ctx := context.Background()
	op := operation{Tool: "reply_ticket", ArgsHash: "h", State: OpDone, StartedAt: time.Now().UTC().Add(-72 * time.Hour),
		Result: json.RawMessage(`{"content":[{"type":"text","text":"customer_secret"}],"isError":false}`)}
	if err := putDoc(ctx, e.kv, store.Item{PK: tenantPK(tenant), SK: "OP#g#t", Ver: 1}, op, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	e.s.KV = &throttledKV{KV: e.kv, left: 1}
	// The warning names the call that stopped it.
	if _, _, err := e.s.CleanStoredAnswers(ctx); !awsapi.Is(err, "ThrottlingException") || !strings.Contains(err.Error(), "couldn't scan a page (Scan)") {
		t.Fatalf("a throttled scan: %v", err)
	}
	if _, ok, _ := e.kv.Get(ctx, onceAnswers.PK, onceAnswers.SK); ok {
		t.Fatal("marked done though the scan was throttled")
	}
	if d, _, err := e.s.CleanStoredAnswers(ctx); err != nil || d != 1 {
		t.Fatalf("the next sweep: dropped=%d %v", d, err)
	}
	if _, ok, _ := e.kv.Get(ctx, onceAnswers.PK, onceAnswers.SK); !ok {
		t.Fatal("the next sweep didn't mark it done")
	}
}

// The clean-up paces its writes by their size, so one tenant's partition key
// stays under DynamoDB's per-key limit: a heavy record waits longer.
func TestTheCleanUpPacesItsWritesBySize(t *testing.T) {
	was := cleanPace
	cleanPace = 1000
	t.Cleanup(func() { cleanPace = was })
	begun := time.Now()
	pace(context.Background(), 100)
	if d := time.Since(begun); d < 90*time.Millisecond || d > time.Second {
		t.Fatalf("100 bytes at 1,000 a second took %v, want about 100 ms", d)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	begun = time.Now()
	pace(ctx, 1_000_000)
	if d := time.Since(begun); d > 100*time.Millisecond {
		t.Fatalf("a pace outlived its context: %v", d)
	}
}
