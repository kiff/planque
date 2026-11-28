package planque

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/store"
)

// An account's hosted browser runs Server.BrowserSessions sessions at a
// time (one unless set, at most maxBrowserSessions): a taf's call, a pay, or
// the owner's live sign-in view. Two sessions of one grant at once are
// safe: each seals its sign-ins only over the version it opened with
// (keepBrowserState), so an older copy never overwrites a refreshed one. The worker's own cap is shared by every account, so without
// this one account's calls could hold every session. Each session takes a
// lease, a record under the account's own key that only planque writes,
// from the grant's tenant and never from a call's input; it lapses by itself
// (browserLeaseSpare past the call's deadline) if the call never frees it.
//
// A call that finds the account's browser busy is refused before anything
// is sent (browser_busy), with the seconds to wait. A call the owner
// approved waits for it instead, up to browserLeaseWait and never past its
// permit: asking the owner twice for one call is worse than a short wait.
// The wait is short enough that it and the call fit inside KIFF's forward
// timeout (kiffForwardTimeout): past that, KIFF records a call as maybe
// sent and never sends it again, which is worse than either.

// maxBrowserSessions is the most sessions one account's hosted browser may
// run at once, whatever is set: the worker's own cap is shared by every
// account.
const maxBrowserSessions = 3

// browserSessions is how many sessions one account's hosted browser runs
// at once: BrowserSessions, 1 when unset, at most maxBrowserSessions.
func (s *Server) browserSessions() int {
	return min(max(s.BrowserSessions, 1), maxBrowserSessions)
}

// kiffForwardTimeout is how long KIFF's gateway waits for a tool's answer
// (kiff/cloud's ForwardTimeout); browserCallRoom, what an approved call
// keeps of it for its own run after waiting for the account's browser.
const (
	kiffForwardTimeout = 30 * time.Second
	browserCallRoom    = 10 * time.Second
)

// browserRunNeeds is what an approved call keeps of its request's own
// deadline for the run, when the request has one: the wait ends that much
// before it.
var browserRunNeeds = 15 * time.Second

var (
	// browserLeaseWait is the longest an approved call waits for the
	// account's browser, leaving browserCallRoom of kiffForwardTimeout for
	// the call itself; browserLeasePoll, how often it looks again.
	browserLeaseWait = 20 * time.Second
	browserLeasePoll = time.Second
	// browserLeaseSpare is how long past its call's deadline a lease
	// outlives a call that never freed it; browserLeaseLongest, the lease
	// of a call with no deadline.
	browserLeaseSpare   = 30 * time.Second
	browserLeaseLongest = 5 * time.Minute
)

// errBrowserBusy: every session of the account's browser is taken.
var errBrowserBusy = errors.New("this account's browser is busy with another call")

// browserLease is one session's hold on the account's browser.
type browserLease struct {
	PK, SK string
	Ver    int64
}

func skBrowserLease(n int) string { return "BROWSERLEASE#" + strconv.Itoa(n) }

// takeBrowserLease holds one of tenant's sessions until until, or answers
// errBrowserBusy with how long until the first one frees.
func (s *Server) takeBrowserLease(ctx context.Context, tenant string, until time.Time) (*browserLease, time.Duration, error) {
	now := s.now()
	pk := tenantPK(tenant)
	soonest := time.Duration(-1)
	for n := range s.browserSessions() {
		sk := skBrowserLease(n)
		it, found, err := s.KV.Get(ctx, pk, sk)
		if err != nil {
			return nil, 0, err
		}
		if found && it.Lease > now.UnixMilli() {
			if left := time.UnixMilli(it.Lease).Sub(now); soonest < 0 || left < soonest {
				soonest = left
			}
			continue
		}
		next := store.Item{PK: pk, SK: sk, Doc: []byte("{}"), Ver: it.Ver + 1, Lease: until.UnixMilli(), TTL: until.Add(time.Hour).Unix()}
		cond := store.Cond{NotExists: true}
		if found {
			cond = store.Cond{Ver: &it.Ver, LeaseFreeAt: now.UnixMilli()}
		}
		if err := s.KV.Put(ctx, next, cond); errors.Is(err, store.ErrCondition) {
			// Another call took it between the read and the write.
			if soonest < 0 {
				soonest = browserLeasePoll
			}
			continue
		} else if err != nil {
			return nil, 0, err
		}
		return &browserLease{PK: pk, SK: sk, Ver: next.Ver}, 0, nil
	}
	return nil, max(soonest, time.Second), errBrowserBusy
}

// dropBrowserLease frees l, unless it lapsed and another session took it.
func (s *Server) dropBrowserLease(ctx context.Context, l *browserLease) {
	if l == nil {
		return
	}
	ver := l.Ver
	_ = s.KV.Put(ctx, store.Item{PK: l.PK, SK: l.SK, Doc: []byte("{}"), Ver: l.Ver + 1, TTL: s.now().Add(time.Hour).Unix()}, store.Cond{Ver: &ver})
}

// callLeaseUntil is how long a call's session holds the account's browser:
// past the call's own deadline by browserLeaseSpare, or browserLeaseLongest
// when it has none.
func (s *Server) callLeaseUntil(ctx context.Context) time.Time {
	if d, ok := ctx.Deadline(); ok {
		return d.Add(browserLeaseSpare)
	}
	return s.now().Add(browserLeaseLongest)
}

// leaseForCall takes the account's browser for a call: at once, or, for a
// call the owner approved (its permit names the approval), waiting up to
// browserLeaseWait and never past the permit. It answers the wait left when
// the browser stays busy.
func (s *Server) leaseForCall(ctx context.Context, tenant string, p permit.Permit) (*browserLease, time.Duration, error) {
	l, left, err := s.takeBrowserLease(ctx, tenant, s.callLeaseUntil(ctx))
	if !errors.Is(err, errBrowserBusy) || p.Approval == "" {
		return l, left, err
	}
	end := s.now().Add(browserLeaseWait)
	if exp := p.Expires; !exp.IsZero() && exp.Before(end) {
		end = exp
	}
	// And never so late that the run itself can't finish before the request's
	// own deadline, whatever its caller's timeout is.
	if d, ok := ctx.Deadline(); ok && d.Add(-browserRunNeeds).Before(end) {
		end = d.Add(-browserRunNeeds)
	}
	for s.now().Before(end) {
		select {
		case <-ctx.Done():
			return nil, left, errBrowserBusy
		case <-time.After(min(browserLeasePoll, max(end.Sub(s.now()), time.Millisecond))):
		}
		if l, left, err = s.takeBrowserLease(ctx, tenant, s.callLeaseUntil(ctx)); !errors.Is(err, errBrowserBusy) {
			return l, left, err
		}
	}
	return nil, left, errBrowserBusy
}
