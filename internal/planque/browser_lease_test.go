package planque

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/store"
)

// shortLeaseWaits makes an approved call's wait short for a test.
func shortLeaseWaits(t *testing.T, wait time.Duration) {
	t.Helper()
	w, p := browserLeaseWait, browserLeasePoll
	browserLeaseWait, browserLeasePoll = wait, 10*time.Millisecond
	t.Cleanup(func() { browserLeaseWait, browserLeasePoll = w, p })
}

// An account's hosted browser runs one session at a time: the second take
// is refused with the time left, another account's isn't, a freed lease is
// free again, and one that lapsed frees itself.
func TestAnAccountsBrowserRunsOneSessionAtATime(t *testing.T) {
	s := &Server{KV: &store.Memory{}}
	ctx := context.Background()
	now := time.Now()
	s.Now = func() time.Time { return now }
	a, _, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Minute))
	if err != nil || a == nil {
		t.Fatalf("the first: %v", err)
	}
	if _, left, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Minute)); !errors.Is(err, errBrowserBusy) || left < 59*time.Second || left > time.Minute {
		t.Fatalf("the second on the same account: %v, %s", err, left)
	}
	b, _, err := s.takeBrowserLease(ctx, "t_b", now.Add(time.Minute))
	if err != nil || b == nil {
		t.Fatalf("another account: %v", err)
	}
	// Another account freeing its own never frees this one's.
	s.dropBrowserLease(ctx, b)
	if _, _, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Minute)); !errors.Is(err, errBrowserBusy) {
		t.Fatalf("t_a freed by t_b's drop: %v", err)
	}
	s.dropBrowserLease(ctx, a)
	again, _, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Second))
	if err != nil {
		t.Fatalf("after the drop: %v", err)
	}
	// One that lapses frees itself, and its holder's late drop doesn't free
	// the session another call took since.
	now = now.Add(2 * time.Second)
	taken, _, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("after the lapse: %v", err)
	}
	s.dropBrowserLease(ctx, again)
	if _, _, err := s.takeBrowserLease(ctx, "t_a", now.Add(time.Minute)); !errors.Is(err, errBrowserBusy) {
		t.Fatalf("a lapsed holder's drop freed the new holder's session: %v", err)
	}
	s.dropBrowserLease(ctx, taken)
}

// Calls racing for the account's one session: exactly one wins.
func TestOnlyOneCallTakesTheAccountsBrowser(t *testing.T) {
	s := &Server{KV: &store.Memory{}}
	var wg sync.WaitGroup
	var mu sync.Mutex
	won := 0
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, _, err := s.takeBrowserLease(context.Background(), "t_a", time.Now().Add(time.Minute)); err == nil {
				mu.Lock()
				won++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if won != 1 {
		t.Fatalf("%d calls took the one session", won)
	}
}

// A call the owner approved waits for the account's browser, up to the
// bound and never past its permit; one nobody approved doesn't wait.
func TestAnApprovedCallWaitsForTheAccountsBrowser(t *testing.T) {
	shortLeaseWaits(t, 2*time.Second)
	s := &Server{KV: &store.Memory{}}
	ctx := context.Background()
	hold := func() *browserLease {
		l, _, err := s.takeBrowserLease(ctx, "t_a", time.Now().Add(time.Minute))
		if err != nil {
			t.Fatal(err)
		}
		return l
	}
	approved := permit.Permit{Approval: "exc_1", Expires: time.Now().Add(time.Minute)}

	// Freed while it waits: it runs.
	l := hold()
	go func() { time.Sleep(100 * time.Millisecond); s.dropBrowserLease(ctx, l) }()
	start := time.Now()
	got, _, err := s.leaseForCall(ctx, "t_a", approved)
	if err != nil || got == nil || time.Since(start) < 100*time.Millisecond {
		t.Fatalf("an approved call waiting for a freed browser: %v after %s", err, time.Since(start))
	}
	s.dropBrowserLease(ctx, got)

	// Not approved: refused at once.
	l = hold()
	start = time.Now()
	if _, left, err := s.leaseForCall(ctx, "t_a", permit.Permit{Expires: time.Now().Add(time.Minute)}); !errors.Is(err, errBrowserBusy) || left <= 0 || time.Since(start) > 500*time.Millisecond {
		t.Fatalf("an unapproved call: %v, %s, after %s", err, left, time.Since(start))
	}
	// Approved, never freed: refused when the bound runs out.
	start = time.Now()
	if _, _, err := s.leaseForCall(ctx, "t_a", approved); !errors.Is(err, errBrowserBusy) || time.Since(start) < 2*time.Second {
		t.Fatalf("an approved call past the bound: %v after %s", err, time.Since(start))
	}
	// Its permit ending first ends the wait first.
	start = time.Now()
	soon := permit.Permit{Approval: "exc_2", Expires: time.Now().Add(200 * time.Millisecond)}
	if _, _, err := s.leaseForCall(ctx, "t_a", soon); !errors.Is(err, errBrowserBusy) || time.Since(start) > time.Second {
		t.Fatalf("an approved call past its permit: %v after %s", err, time.Since(start))
	}
	s.dropBrowserLease(ctx, l)
}

// A taf's call on an account whose browser is busy is refused before any
// session starts, with when to try again; once the session is free, it runs.
// The lease is the account's own, from the grant: a call's input names none
// (planque reads browse_open's arguments strictly, so no key can carry one).
func TestABrowserCallOnABusyAccountIsRefusedBeforeAnythingStarts(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	ctx := context.Background()
	other, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	before := e.b.Sessions
	out, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	if code(sc) != browser.CodeBrowserBusy || sc["sent"] != false || e.b.Sessions != before || sc["retry_after"] == nil || !strings.Contains(out, "busy with another call") {
		t.Fatalf("a call on a busy account: %s %v (sessions %d → %d)", out, sc, before, e.b.Sessions)
	}
	e.s.dropBrowserLease(ctx, other)
	if out, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("once free: %s", out)
	}
	// The call freed its own lease as it ended.
	if l, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("the call kept the account's session: %v", err)
	} else {
		e.s.dropBrowserLease(ctx, l)
	}
}

// The owner's live sign-in view is the account's browser session while it's
// open: a call finds the account busy, the view can't open over a taf's call,
// and Done frees it.
func TestTheSignInViewHoldsTheAccountsBrowser(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	ctx := context.Background()
	// A taf's call holds the account's browser: the view waits for none.
	l, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if st, out := e.startView(g.ID, nil); st != 409 || out["code"] != browser.CodeBrowserBusy || out["retry_after"] == nil || e.b.Sessions != 0 {
		t.Fatalf("a view over a taf's call: %d %v (sessions %d)", st, out, e.b.Sessions)
	}
	e.s.dropBrowserLease(ctx, l)
	st, out := e.startView(g.ID, map[string]any{"url": fakes.ShopLogin})
	if st != 200 {
		t.Fatalf("the view: %d %v", st, out)
	}
	if _, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute)); !errors.Is(err, errBrowserBusy) {
		t.Fatalf("the open view doesn't hold the account's browser: %v", err)
	}
	if st, _, body := e.view("POST", out["view"].(string)+"/done", nil); st != 200 {
		t.Fatalf("done: %d %s", st, body)
	}
	if l, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute)); err != nil {
		t.Fatalf("done kept the account's browser: %v", err)
	} else {
		e.s.dropBrowserLease(ctx, l)
	}
}

// Through a whole call: one the owner approved, on an account whose browser
// another call holds, waits and then runs once that session ends.
func TestAnApprovedBrowserCallRunsOnceTheAccountsBrowserFrees(t *testing.T) {
	shortLeaseWaits(t, 5*time.Second)
	e := newBrowserEnv(t)
	g := e.connect()
	ctx := context.Background()
	l, _, err := e.s.takeBrowserLease(ctx, g.Tenant, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	go func() { time.Sleep(200 * time.Millisecond); e.s.dropBrowserLease(ctx, l) }()
	op := "op_open_approved"
	args := mustJSON(map[string]any{"url": fakes.ShopItem, IdempotencyArg: op})
	p := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_open", Sub: g.Request.Agent, Op: op, Args: args, Approval: "exc_" + op})
	before, start := e.b.Sessions, time.Now()
	out, sc := textOf(e.s.call(ctx, g, "browse_open", args, p))
	if code(sc) != "" || e.b.Sessions != before+1 || time.Since(start) < 200*time.Millisecond {
		t.Fatalf("an approved call on a busy account: %s (sessions %d → %d, after %s)", out, before, e.b.Sessions, time.Since(start))
	}
}
