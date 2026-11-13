package planque

import (
	"testing"
	"time"
)

// A request nobody signs ends on its own:
// a pending grant gets a deadline, 7 days from when it was asked, and the
// sweep then ends it as not approved in time. Before, a plain pending
// grant had no deadline, so it waited forever, outliving its taf and its
// agent, ready to be approved by mistake.
const pendingDeadline = 7 * 24 * time.Hour

// later moves the server's clock and the phone's together, so the phone
// can still sign at that time.
func (e *env) later(d time.Duration) {
	e.at(d)
	e.phone.Now = e.s.Now
}

func (e *env) status(id string) string {
	e.t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("read %s: %d %v", id, st, out)
	}
	s, _ := out["status"].(string)
	return s
}

func TestAPendingGrantNobodySignsEndsAfterSevenDays(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()

	// The day before: the sweep leaves it, and it still waits.
	e.later(pendingDeadline - time.Hour)
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 0 {
		t.Fatalf("a sweep before the deadline: %d, %v", n, err)
	}
	if s := e.status(id); s != StatusPending {
		t.Fatalf("before the deadline: %s, want pending", s)
	}

	// Past it, the sweep alone ends it, though nobody read it.
	e.later(pendingDeadline + time.Minute)
	if n, err := e.s.Sweep(t.Context()); err != nil || n != 1 {
		t.Fatalf("the sweep past the deadline: %d, %v; want it to end the request", n, err)
	}
	if s := e.status(id); s != StatusFailed {
		t.Fatalf("past the deadline: %s, want failed (not approved in time)", s)
	}
	// The record stays, as the account's history, and can't be approved now.
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, ""); st == 200 {
		t.Fatalf("a request past its deadline was approved: %v", out)
	}
	if a, r := e.vendor.Live(); a != 0 || r != 0 {
		t.Fatalf("tokens live at the vendor: %d, %d", a, r)
	}
}

// Reading a pending request doesn't move its deadline, and a request the
// sweep hasn't reached yet is ended by the next read, as every expiry is.
func TestReadingAPendingGrantNeitherExtendsNorOutlivesItsDeadline(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, _ := e.request()
	e.later(6 * 24 * time.Hour)
	if s := e.status(id); s != StatusPending {
		t.Fatalf("on day 6: %s", s)
	}
	e.later(pendingDeadline + time.Minute)
	if s := e.status(id); s != StatusFailed {
		t.Fatalf("read past the deadline, before any sweep: %s, want failed", s)
	}
}

// Only a request waiting for the owner's signature has this deadline: a
// connected grant stays connected however long it's been, and a request
// approved in time connects as before.
func TestThePendingDeadlineLeavesAConnectedGrant(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	e.later(30 * 24 * time.Hour)
	if _, err := e.s.Sweep(t.Context()); err != nil {
		t.Fatal(err)
	}
	if s := e.status(g.ID); s != StatusConnected {
		t.Fatalf("a connected grant after 30 days: %s, want connected", s)
	}

	e2 := newEnv(t)
	e2.enrol()
	id, digest := e2.request()
	e2.later(pendingDeadline - time.Hour)
	if st, out := e2.approveAndSignIn(id, digest); st != 200 || out["status"] == StatusFailed {
		t.Fatalf("approved the day before the deadline: %d %v", st, out)
	}
}
