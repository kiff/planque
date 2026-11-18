package planque

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// kiffState is a stand-in for KIFF's GET /v1/ops/proofs/state.
func kiffState(t *testing.T, answers map[string]string, calls *atomic.Int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get("Authorization") != "Bearer state-token" {
			http.Error(w, "no", http.StatusUnauthorized)
			return
		}
		entry := r.URL.Query().Get("entry")
		switch state := answers[entry]; state {
		case "":
			http.Error(w, "boom", http.StatusInternalServerError)
		case "wrong-entry":
			_ = json.NewEncoder(w).Encode(map[string]string{"entry": "other@1", "state": "green", "since": "2026-11-17T10:00:00Z"})
		case "slow":
			select {
			case <-time.After(3 * time.Second):
			case <-r.Context().Done():
			}
			http.Error(w, "late", http.StatusGatewayTimeout)
		case "garbage":
			_, _ = w.Write([]byte("<html>"))
		default:
			_ = json.NewEncoder(w).Encode(map[string]string{"entry": entry, "app": strings.SplitN(entry, "-", 2)[0], "state": state, "since": "2026-11-17T10:00:00Z"})
		}
	}))
}

func proofEntries() map[string]CatalogEntry {
	return map[string]CatalogEntry{
		"zendesk-key@1":   {Vendor: "Zendesk", Kind: KindAPIKey},
		"qonto-key@1":     {Vendor: "Qonto", Kind: KindAPIKey},
		"pipedrive-key@1": {Vendor: "Pipedrive", Kind: KindAPIKey},
	}
}

// KIFF's green answers are the only proofs: the shipped file is not shown, an
// unknown or unreadable or mismatched answer shows nothing.
func TestProvedComesFromKIFFAloneWhenItIsAsked(t *testing.T) {
	var calls atomic.Int32
	srv := kiffState(t, map[string]string{
		"zendesk-key@1":   "green",
		"qonto-key@1":     "unknown",
		"pipedrive-key@1": "wrong-entry",
	}, &calls)
	defer srv.Close()
	s := &Server{CatalogEntries: proofEntries(), ProofState: NewProofState(srv.URL+"/v1/ops/proofs/state", "state-token")}
	got := s.provedMap(t.Context())
	if len(got) != 1 || got["zendesk"].At != "2026-11-17" || got["zendesk"].How != "key" {
		t.Fatalf("proofs: %+v", got)
	}
	if _, shipped := got["gmail"]; shipped {
		t.Fatal("a shipped proof shows while KIFF is asked")
	}
	if p := proofOf(got, "zendesk-key@1", "key"); p == nil {
		t.Fatal("the zendesk entry doesn't carry its proof")
	}
	if p := proofOf(got, "qonto-key@1", "key"); p != nil {
		t.Fatalf("an unknown answer proves: %+v", p)
	}
	// The public route says the same.
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/v1/proved", nil))
	var body struct {
		Proved map[string]publicProof `json:"proved"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || len(body.Proved) != 1 || body.Proved["zendesk"].At != "2026-11-17" {
		t.Fatalf("%v %s", err, w.Body)
	}
}

// With no setting, nothing changes: the shipped proofs stand.
func TestProvedIsTheShippedFileWhenKIFFIsNotAsked(t *testing.T) {
	if got := (&Server{CatalogEntries: proofEntries()}).provedMap(t.Context()); len(got) != len(Proved) {
		t.Fatalf("%d, want %d", len(got), len(Proved))
	}
	for _, c := range [][2]string{{"", ""}, {"https://x.example/state", ""}, {"", "t"}, {"http://kiff.example/state", "t"}, {"https://u:p@kiff.example/state", "t"}, {"https://kiff.example/state?x=1", "t"}, {"not a url", "t"}} {
		if NewProofState(c[0], c[1]) != nil {
			t.Errorf("%q/%q makes a reader", c[0], c[1])
		}
	}
}

// KIFF down, slow or answering wrongly proves nothing, and never takes the
// catalog down; a failure is not asked again at once, a success is cached.
func TestAFailingKIFFShowsNoProofAndIsNotHammered(t *testing.T) {
	var calls atomic.Int32
	srv := kiffState(t, map[string]string{"zendesk-key@1": "green", "qonto-key@1": "garbage"}, &calls)
	p := NewProofState(srv.URL, "state-token")
	now := time.Now()
	p.now = func() time.Time { return now }
	ids := []string{"zendesk-key@1", "qonto-key@1", "pipedrive-key@1", "not an entry", "connector:composio:gmail"}
	if g := p.Greens(t.Context(), ids); len(g) != 1 || g["zendesk-key@1"].State != "green" {
		t.Fatalf("%+v", g)
	}
	if calls.Load() != 3 {
		t.Fatalf("%d calls: only the three real entries are asked", calls.Load())
	}
	p.Greens(t.Context(), ids)
	if calls.Load() != 3 {
		t.Fatalf("%d calls: answers are kept", calls.Load())
	}
	now = now.Add(proofStateFailTTL + time.Second)
	p.Greens(t.Context(), ids)
	if calls.Load() != 5 {
		t.Fatalf("%d calls: failures are asked again after their short wait, successes not", calls.Load())
	}
	now = now.Add(proofStateTTL)
	srv.Close() // now down
	if g := p.Greens(t.Context(), ids); len(g) != 0 {
		t.Fatalf("a dead KIFF still proves: %+v", g)
	}
	// Wrong token: refused, no proof.
	var c2 atomic.Int32
	srv2 := kiffState(t, map[string]string{"zendesk-key@1": "green"}, &c2)
	defer srv2.Close()
	if g := NewProofState(srv2.URL, "other").Greens(t.Context(), ids); len(g) != 0 {
		t.Fatalf("a refused token proves: %+v", g)
	}
}

// The token goes to the address set and nowhere else: a redirect is not followed.
func TestTheStateTokenIsNotSentOnARedirect(t *testing.T) {
	var elsewhere atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { elsewhere.Add(1) }))
	defer other.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusFound)
	}))
	defer srv.Close()
	if g := NewProofState(srv.URL, "state-token").Greens(t.Context(), []string{"zendesk-key@1"}); len(g) != 0 || elsewhere.Load() != 0 {
		t.Fatalf("%+v, %d calls elsewhere", g, elsewhere.Load())
	}
}

// A slow KIFF never slows the catalog: one deadline for the whole read, no
// proof for what didn't answer, and no new asks inside the failure wait.
func TestASlowKIFFIsWaitedOnOnceForTheWholeRead(t *testing.T) {
	var calls atomic.Int32
	answers := map[string]string{}
	entries := map[string]CatalogEntry{}
	for _, n := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"} {
		answers[n+"-key@1"] = "slow"
		entries[n+"-key@1"] = CatalogEntry{Vendor: strings.ToUpper(n), Kind: KindAPIKey}
	}
	srv := kiffState(t, answers, &calls)
	defer srv.Close()
	s := &Server{CatalogEntries: entries, ProofState: NewProofState(srv.URL, "state-token")}
	start := time.Now()
	if got := s.provedMap(t.Context()); len(got) != 0 {
		t.Fatalf("proofs from a KIFF that never answered: %+v", got)
	}
	if took := time.Since(start); took > 1500*time.Millisecond {
		t.Fatalf("the read took %s", took)
	}
	asked := calls.Load()
	s.provedMap(t.Context())
	if calls.Load() != asked {
		t.Fatalf("%d new asks inside the failure wait", calls.Load()-asked)
	}
}

// expired is a state KIFF now answers: it proves nothing, by name.
func TestAnExpiredEntryProvesNothing(t *testing.T) {
	var calls atomic.Int32
	srv := kiffState(t, map[string]string{"zendesk-key@1": "expired", "qonto-key@1": "green"}, &calls)
	defer srv.Close()
	s := &Server{CatalogEntries: map[string]CatalogEntry{
		"zendesk-key@1": {Vendor: "Zendesk", Kind: KindAPIKey},
		"qonto-key@1":   {Vendor: "Qonto", Kind: KindAPIKey},
	}, ProofState: NewProofState(srv.URL, "state-token")}
	got := s.provedMap(t.Context())
	if _, ok := got["zendesk"]; ok || len(got) != 1 || got["qonto"].How != "key" {
		t.Fatalf("%+v", got)
	}
}
