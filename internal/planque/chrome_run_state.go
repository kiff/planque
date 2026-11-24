package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// A run the person paused in their Chrome (tafs docs/chrome.md, the floating
// bar's Pause) waits for them instead of ending the taf's turn. The extension
// tells planque each time the person pauses, resumes or stops a run, in a
// message its Chrome key signs (run_state). planque keeps the newest state of
// the run, and the taf's runtime asks for it with its own key
// (GET …/grants/{grant}/chrome/run) while the run is paused, so it sends the
// step again on Resume rather than giving up. Nothing about the state goes
// through KIFF's gateway: the person's Record keeps the steps alone.

// labelRunState is what a Chrome signs a run's state under.
const labelRunState = "tafs-chrome/run_state/v1"

// The states a run's Chrome says it is in.
const (
	runWorking = "working"
	runPaused  = "paused"
	runStopped = "stopped"
	// runEnded is the read's answer for a grant with no run: never sent by a Chrome.
	runEnded = "ended"
	// runUnreported is the read's answer for a run whose Chrome never said its
	// state (an extension from before run_state): never sent by a Chrome. The
	// runtime doesn't wait on it, as it can't know when the person resumes.
	runUnreported = "unreported"
)

// chromeRunPausedFor is how long a paused run stays the grant's run with no
// call: past the runtime's own wait (10 minutes), and as long as the
// extension's sweep of a paused run (15 minutes), so neither side ends a run
// the other still counts on.
const chromeRunPausedFor = 15 * time.Minute

// runStateAhead is how far ahead of planque's clock a state's time may be.
const runStateAhead = 2 * time.Minute

// runStateKeptFor is how long a run's state item lives at all.
const runStateKeptFor = 6 * time.Hour

func skChromeRunState(chrome, run string) string { return "CHROMERUNSTATE#" + chrome + "#" + run }

// chromeRunState is one run of a Chrome: the grant it's for and when planque
// started it, and once its Chrome said one, its newest state, the time the
// Chrome gave it, and when planque received it.
type chromeRunState struct {
	Grant    string    `json:"grant"`
	Started  time.Time `json:"started"`
	State    string    `json:"state,omitempty"`
	At       time.Time `json:"at,omitempty"`
	Received time.Time `json:"received,omitempty"`
}

// noteRun records the call's run for its Chrome, once: which grant it's for,
// and when it started. A run_state names only a run noted here.
func (s *Server) noteRun(ctx context.Context, c *chromeSession) {
	if c.run == "" || c.chrome == "" {
		return
	}
	now := s.now()
	it := store.Item{PK: tenantPK(c.tenant), SK: skChromeRunState(c.chrome, c.run), TTL: now.Add(runStateKeptFor).Unix()}
	err := putDoc(ctx, s.KV, it, chromeRunState{Grant: c.grant, Started: now}, store.Cond{NotExists: true})
	if err != nil && !errors.Is(err, store.ErrCondition) {
		s.log().WarnContext(ctx, "chrome run not noted", slog.String("grant", c.grant), slog.String("err", err.Error()))
	}
}

// runPaused says whether the run's Chrome last said it is paused, within
// chromeRunPausedFor of planque receiving that.
func (s *Server) runPaused(ctx context.Context, tenant, chrome, run string, now time.Time) bool {
	var st chromeRunState
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skChromeRunState(chrome, run), &st); err != nil {
		return false
	}
	return st.State == runPaused && now.Sub(st.Received) < chromeRunPausedFor
}

// runStateBody is what a Chrome signs in a run_state: the run, its state, and when.
type runStateBody struct {
	Run   string `json:"run"`
	State string `json:"state"`
	At    string `json:"at"`
}

// relayRunState keeps a run's state from its Chrome, signed by that Chrome's
// key, for a run planque noted for that Chrome. Refused: an unknown run, a
// state after stopped (final), an unsigned or malformed one, or a time far
// ahead of planque's or before the run began. An older state than the one
// kept is answered as kept, and changes nothing.
func (s *Server) relayRunState(ctx context.Context, connID string, conn chromeConn, m map[string]json.RawMessage) {
	var signed runStateBody
	var sig string
	_ = json.Unmarshal(m["run"], &signed.Run)
	_ = json.Unmarshal(m["state"], &signed.State)
	_ = json.Unmarshal(m["at"], &signed.At)
	refuse := func(code string) {
		s.relayPost(ctx, connID, map[string]any{"type": "run_state_refused", "run": signed.Run, "code": code})
	}
	if json.Unmarshal(m["sig"], &sig) != nil || signed.Run == "" || len(signed.Run) > 200 {
		refuse("unsigned")
		return
	}
	switch signed.State {
	case runWorking, runPaused, runStopped:
	default:
		refuse("bad_state")
		return
	}
	at, err := time.Parse(time.RFC3339Nano, signed.At)
	if err != nil {
		refuse("bad_at")
		return
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChrome(conn.Chrome), &c); err != nil {
		s.relayClose(ctx, connID, closeUnpaired, "this Chrome was unpaired")
		return
	}
	if !stepkey.Verify(c.Key, labelRunState, signed, sig) {
		refuse("unsigned")
		return
	}
	pk, sk := tenantPK(conn.Tenant), skChromeRunState(conn.Chrome, signed.Run)
	// Two tries: a state from another connection of the same Chrome may land between the read and the write.
	for range 2 {
		var st chromeRunState
		it, err := getDoc(ctx, s.KV, pk, sk, &st)
		if errors.Is(err, errNotFound) {
			refuse("unknown_run")
			return
		}
		if err != nil {
			refuse("not_kept")
			return
		}
		now := s.now()
		switch {
		case at.After(now.Add(runStateAhead)) || at.Before(st.Started):
			refuse("bad_at")
			return
		case st.State == runStopped && signed.State == runStopped:
			s.relayPost(ctx, connID, map[string]any{"type": "run_state_kept", "run": signed.Run})
			return
		case st.State == runStopped:
			// Stopped is final: a late Resume never brings back a run the person ended.
			refuse("run_stopped")
			return
		case !st.At.IsZero() && !at.After(st.At):
			s.relayPost(ctx, connID, map[string]any{"type": "run_state_kept", "run": signed.Run})
			return
		}
		st.State, st.At, st.Received = signed.State, at, now
		ver := it.Ver
		next := store.Item{PK: pk, SK: sk, Ver: ver + 1, TTL: it.TTL}
		if err := putDoc(ctx, s.KV, next, st, store.Cond{Ver: &ver}); err == nil {
			s.relayPost(ctx, connID, map[string]any{"type": "run_state_kept", "run": signed.Run})
			return
		} else if !errors.Is(err, store.ErrCondition) {
			s.log().WarnContext(ctx, "chrome run state not kept", slog.String("err", err.Error()))
			refuse("not_kept")
			return
		}
	}
	refuse("not_kept")
}

// readChromeRun answers a taf's runtime, with the taf's own KIFF key, its
// grant's run in the owner's Chrome and that run's newest state: working,
// paused, stopped, unreported when its Chrome never said one, or ended when
// the grant has no run. Only the grant's own
// taf reads it, and nothing else of the run.
func (s *Server) readChromeRun(w http.ResponseWriter, r *http.Request) {
	tenant, id := r.PathValue("tenant"), r.PathValue("id")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), tenant, id)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if me.AgentID == "" || me.AgentID != g.Request.Agent {
		writeErr(w, http.StatusForbidden, "only the grant's own taf reads its run")
		return
	}
	var doc chromeRunDoc
	if _, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skChromeRun(g), &doc); errors.Is(err, errNotFound) || err == nil && doc.Run == "" {
		writeJSON(w, http.StatusOK, map[string]string{"run": "", "state": runEnded})
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	state := runUnreported
	var st chromeRunState
	if _, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skChromeRunState(doc.Chrome, doc.Run), &st); err == nil && st.State != "" {
		state = st.State
	} else if err != nil && !errors.Is(err, errNotFound) {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"run": doc.Run, "state": state})
}

// pauseCodes are the Chrome's answers for a run that's paused.
var pauseCodes = map[string]bool{"paused_by_owner": true, "paused_failures": true, "paused_no_progress": true}

// withPausedRun names the run and its grant in a refusal that says it's
// paused, so the taf's runtime knows which run to wait on and where to read
// it. Any other answer is unchanged.
func withPausedRun(res any, run, grant string) any {
	m, ok := res.(map[string]any)
	if !ok || run == "" {
		return res
	}
	sc, _ := m["structuredContent"].(map[string]any)
	if code, _ := sc["code"].(string); sc != nil && pauseCodes[code] {
		sc["run"], sc["grant"] = run, grant
	}
	return res
}
