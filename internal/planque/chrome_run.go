package planque

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/kiff/planque/internal/store"
)

// One run per taf turn in the owner's Chrome. A grant's calls that follow one
// another share one run: its tab, its tab group, its bar and its Stop, so the
// person sees one taf at work and not a new tab per call. The extension keys
// all of that by the step's run id (tafs chrome/src/runner.js), so the same
// id across calls is the same tab, navigated in place.
//
//   - The run is the grant's while its calls keep coming: a call within
//     chromeRunIdle of the last one takes it, under a lease, so two calls at
//     once never share a tab (the second runs alone, as each call did before).
//   - A call on the page the run's tab already shows reads it where it is,
//     never reloading it. Any other address opens in the same tab.
//   - The run ends when the person stops it (the extension answers stopped:
//     the taf is told, and the next call has a new run), when its Chrome
//     drops or doesn't answer, when the grant's signed sites change, or once
//     it's idle: the next call closes it then, and the extension's own idle
//     stop (5 minutes) closes one no call follows.
//
// A run's sites are fixed at its first step in the extension, so a grant
// whose signed sites changed starts a new run.

// chromeRunIdle is how long a run waits for the grant's next call: well under
// the extension's own idle stop (RUN_IDLE_MS, 5 minutes), so the extension
// never stops a run planque still counts on.
const chromeRunIdle = 2 * time.Minute

// chromeRunLease is how long a call holds its run: its own budget and a
// margin, so a call that died lets go on its own.
const chromeRunLease = chromeCallBudget + 10*time.Second

// skChromeRun is a grant's run in its Chrome.
func skChromeRun(g *Grant) string { return "CHROMERUN#" + g.ID }

// chromeRunDoc is a grant's run: its id and step count, the extension's read
// id and the address its tab shows, the sites' signature and the Chrome it
// began under, the kiff/agents run whose calls it serves, and when its last
// call ended. Run is empty once it ended.
type chromeRunDoc struct {
	Run    string `json:"run,omitempty"`
	N      int    `json:"n,omitempty"`
	ReadID string `json:"read_id,omitempty"`
	URL    string `json:"url,omitempty"`
	Sites  string `json:"sites"`
	Chrome string `json:"chrome"`
	// AgentRun is the kiff/agents run whose calls this run serves (the
	// permit's run, chromeStep.AgentRun): only that run's next call takes it
	// over, since a run may hold the person's own tab, which joins only the
	// run they started from the panel. "" for calls whose permit named none.
	AgentRun string    `json:"agent_run,omitempty"`
	At       time.Time `json:"at"`
}

// heldRun is the run a call holds: the version it wrote, and what the run's
// tab showed when the call took it.
type heldRun struct {
	ver    int64
	url    string
	reused bool
}

// sitesKey is what a run's sites were signed as: the phone's signature, or the
// Chrome a carried signature came with.
func (c *chromeSession) sitesKey() string {
	if c.signed.Sig != "" {
		return c.signed.Sig
	}
	return "carried:" + c.chrome
}

// claimRun makes the grant's run this call's: the one under way when it's
// still the grant's (same Chrome, same signed sites, the same kiff/agents run,
// idle less than chromeRunIdle, no other call holding it), else a new one, under the call's
// own operation id, after closing the one before. When the record can't be
// read or written, or another call holds the run, the call runs alone, as
// each call did before: nothing is shared and its run closes with it.
func (s *Server) claimRun(ctx context.Context, g *Grant, c *chromeSession) {
	now := s.now()
	pk, sk := tenantPK(g.Tenant), skChromeRun(g)
	var doc chromeRunDoc
	it, err := getDoc(ctx, s.KV, pk, sk, &doc)
	found := err == nil
	if err != nil && !errors.Is(err, errNotFound) {
		return
	}
	if found && it.Lease >= now.UnixMilli() {
		return
	}
	next := store.Item{PK: pk, SK: sk, Ver: it.Ver + 1, Lease: now.Add(chromeRunLease).UnixMilli(), TTL: now.Add(chromeRunPausedFor + chromeRunLease).Unix()}
	cond := store.Cond{NotExists: true}
	if found {
		cond = store.Cond{Ver: &it.Ver, LeaseFreeAt: now.UnixMilli()}
	}
	// A run the person paused stays the grant's past its idle time
	// (chromeRunPausedFor), so the step sent again on Resume goes on in its tab.
	// Never another kiff/agents run's: a schedule's, the phone's or the watch's
	// call right after the person's panel run gets a run, and a tab, of its own.
	if found && doc.Run != "" && doc.Chrome == c.chrome && doc.Sites == c.sitesKey() && doc.AgentRun == c.agentRun &&
		(now.Sub(doc.At) < chromeRunIdle || s.runPaused(ctx, g.Tenant, doc.Chrome, doc.Run, now)) {
		if putDoc(ctx, s.KV, next, doc, cond) != nil {
			return
		}
		c.run, c.n, c.readID, c.sent = doc.Run, doc.N, doc.ReadID, true
		c.held = &heldRun{ver: next.Ver, url: doc.URL, reused: true}
		return
	}
	if putDoc(ctx, s.KV, next, chromeRunDoc{Run: c.run, Sites: c.sitesKey(), Chrome: c.chrome, AgentRun: c.agentRun, At: now}, cond) != nil {
		return
	}
	c.held = &heldRun{ver: next.Ver}
	s.noteRun(ctx, c)
	// The run before, idle or under other sites, ends: its group greys and
	// its tab is the person's again.
	if found && doc.Run != "" && doc.Run != c.run && doc.Chrome == c.chrome {
		old := &chromeSession{s: s, tenant: c.tenant, chrome: c.chrome, grant: c.grant, run: doc.Run, n: doc.N, taf: c.taf,
			sites: c.sites, signed: c.signed, sent: true}
		_, _ = old.Close(ctx)
	}
}

// restart moves the call onto a new run, under id: the run it took has no tab
// any more in the Chrome (its service worker started again, or the person
// closed the tab).
func (c *chromeSession) restart(id string) {
	c.run, c.n, c.readID, c.sent = id, 0, "", false
	if c.held != nil {
		c.held.reused = false
		c.held.url = ""
	}
}

// releaseRun hands the run back at the call's end, for the grant's next call:
// its step count, read id and address. A run that ended (stopped, dropped, or
// with no answer) is closed and left for the next call to start anew.
func (s *Server) releaseRun(ctx context.Context, c *chromeSession) {
	now := s.now()
	doc := chromeRunDoc{Run: c.run, N: c.n, ReadID: c.readID, URL: c.url, Sites: c.sitesKey(), Chrome: c.chrome, AgentRun: c.agentRun, At: now}
	if c.ended {
		doc = chromeRunDoc{Sites: doc.Sites, Chrome: c.chrome, AgentRun: c.agentRun, At: now}
	}
	it := store.Item{PK: tenantPK(c.tenant), SK: skChromeRun(&Grant{ID: c.grant}), Ver: c.held.ver + 1, TTL: now.Add(chromeRunPausedFor + chromeRunLease).Unix()}
	ver := c.held.ver
	if err := putDoc(ctx, s.KV, it, doc, store.Cond{Ver: &ver}); err != nil {
		s.log().WarnContext(ctx, "chrome run not kept", slog.String("grant", c.grant), slog.String("err", err.Error()))
	}
	if c.ended && c.sent {
		c.n++
		_ = s.sendChromeStep(ctx, c.tenant, c.base(chromeStep{Op: "close"}))
	}
}

// firstPage is a call's first step: the page at u. A run that took over a tab
// already on u reads it where it is; otherwise u opens, in the run's tab. A
// tab the person took elsewhere meanwhile (the read comes back on another
// address) opens u again, in the same tab. A run whose tab is gone in the
// Chrome (no_page) goes on as a new run, under id, and u opens there.
func (c *chromeSession) firstPage(ctx context.Context, u, id string) (chromePage, error) {
	if c.held != nil && c.held.reused && c.held.url != "" && c.held.url == u {
		p, err := c.read(ctx, 0)
		if err == nil && p.URL == u {
			return p, nil
		}
		var r *chromeRefusal
		if err != nil && (!errors.As(err, &r) || r.Code != "no_page") {
			return p, err
		}
		if err != nil {
			c.restart(id)
			c.s.noteRun(ctx, c)
		}
	}
	return c.stepPage(ctx, chromeStep{Op: "open", URL: u})
}
