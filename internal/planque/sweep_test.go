package planque

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/store"
)

// faultyKV fails the reads and writes it is told to, as a store that is
// briefly unreachable.
type faultyKV struct {
	store.KV
	mu      sync.Mutex
	failGet func(pk, sk string) bool
	failPut func(it store.Item) bool
}

func (f *faultyKV) Get(ctx context.Context, pk, sk string) (store.Item, bool, error) {
	f.mu.Lock()
	fail := f.failGet != nil && f.failGet(pk, sk)
	f.mu.Unlock()
	if fail {
		return store.Item{}, false, errors.New("store unreachable")
	}
	return f.KV.Get(ctx, pk, sk)
}

func (f *faultyKV) Put(ctx context.Context, it store.Item, c store.Cond) error {
	f.mu.Lock()
	fail := f.failPut != nil && f.failPut(it)
	f.mu.Unlock()
	if fail {
		return errors.New("store unreachable")
	}
	return f.KV.Put(ctx, it, c)
}

func (f *faultyKV) set(get func(pk, sk string) bool, put func(store.Item) bool) {
	f.mu.Lock()
	f.failGet, f.failPut = get, put
	f.mu.Unlock()
}

func (e *env) faulty() *faultyKV {
	f := &faultyKV{KV: e.kv}
	e.s.KV = f
	return f
}

func (e *env) dueEntries() int {
	e.t.Helper()
	items, err := e.kv.Query(e.t.Context(), pkDue, "")
	if err != nil {
		e.t.Fatal(err)
	}
	return len(items)
}

func (e *env) at(d time.Duration) { e.s.Now = func() time.Time { return time.Now().Add(d) } }

// kiff/planque#23 A: a sweep that can't read a due grant keeps its entry,
// and the next sweep revokes the sign-in, though nobody read the grant.
func TestASweepThatCantReadAGrantTriesAgain(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	f := e.faulty()
	e.at(holdTTL + time.Minute)
	f.set(func(pk, sk string) bool { return sk == skGrant(id) }, nil)
	if n, err := e.s.Sweep(t.Context()); err == nil || n != 0 {
		t.Fatalf("a sweep that couldn't read the grant: %d, %v", n, err)
	}
	if e.dueEntries() == 0 {
		t.Fatal("the sweep dropped the work it couldn't do")
	}
	f.set(nil, nil)
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("the next sweep: %d, %v", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live after the next sweep: %d, %d", a, r)
	}
}

// kiff/planque#23 B: a grant is never saved with a deadline the sweep
// can't find. Here the phone's confirm can't index the held sign-in's day,
// so it fails, and the grant stays signed in under its confirm deadline,
// which the sweep then ends.
func TestAGrantIsNotSavedWithoutItsDeadlineIndexed(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	id := out["id"].(string)
	st, ap := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, ap)
	}
	code := e.signIn(ap["sign_in_url"].(string))
	f := e.faulty()
	f.set(nil, func(it store.Item) bool { return it.PK == pkDue })
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st == 200 && out["status"] == StatusListed {
		t.Fatalf("listed without its deadline indexed: %d %v", st, out)
	}
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil || g.Status != StatusSignedIn {
		t.Fatalf("the grant after a confirm that couldn't index: %+v %v", g, err)
	}
	f.set(nil, nil)
	e.at(signInTTL + time.Minute)
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("the sweep: %d, %v", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live after the sweep: %d, %d", a, r)
	}
}

// kiff/planque#23 C: an entry is never due before its deadline, so a sweep
// in the deadline's own second leaves it for one after.
func TestASweepInTheDeadlinesSecondKeepsItsEntry(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	// The deadline at second S + 900 ms, a sweep at S + 100 ms.
	deadline := g.HeldUntil.Truncate(time.Second).Add(900 * time.Millisecond)
	g.HeldUntil = deadline
	if err := e.s.saveGrant(t.Context(), g); err != nil {
		t.Fatal(err)
	}
	e.s.Now = func() time.Time { return deadline.Add(-800 * time.Millisecond) }
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 0 {
		t.Fatalf("a sweep before the deadline: %d, %v", n, err)
	}
	if !strings.Contains(dueSK(deadline, tenant, id), ceilSecond(deadline).Format(time.RFC3339)) {
		t.Fatal("the entry isn't rounded up")
	}
	e.s.Now = func() time.Time { return deadline.Add(2 * time.Second) }
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("a sweep after the deadline: %d, %v", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live: %d, %d", a, r)
	}
}

// A vendor that doesn't confirm revoking keeps the token and the work: the
// grant isn't marked failed, and the next sweep revokes. Meanwhile the grant
// acts on nothing: its deadline passed.
func TestAVendorThatDoesntConfirmRevokingIsAskedAgain(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	// Down for the sweep and the read after it.
	e.vendor.RevokeFails = 2
	e.at(holdTTL + time.Minute)
	if n, err := e.s.Sweep(t.Context()); err == nil || n != 0 {
		t.Fatalf("a sweep the vendor didn't confirm: %d, %v", n, err)
	}
	if g, _ := e.s.loadGrant(t.Context(), tenant, id); g.Status == StatusFailed {
		t.Fatal("marked failed with its tokens not revoked")
	}
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools",
		map[string]any{"tools": []any{map[string]any{"tool": "list_charges", "read_only": true}}}, nil, "k_setup"); st != 409 {
		t.Fatalf("an overdue grant took its tools: %d", st)
	}
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("the next sweep: %d, %v", n, err)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens still live: %d, %d", a, r)
	}
	if e.dueEntries() != 0 {
		t.Fatal("the entry outlived its work")
	}
}
