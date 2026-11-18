package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

// KIFF is the one record of whether an app works (its proof monitor): an
// entry shows as proved here only when KIFF says it is green, and every
// other answer, an error, a timeout, a refusal or no setting at all, shows no
// proof. planque never works out an app's freshness itself.
const (
	proofStateTTL     = 60 * time.Second // KIFF caches the same for as long
	proofStateFailTTL = 15 * time.Second // a failure is not asked again at once
	proofStateTimeout = 2 * time.Second
	proofStateWidth   = 8  // asks in flight at once
	proofStateMost    = 64 // entries asked for in one read
	proofStateBody    = 4 << 10
)

// ProofState reads KIFF's state for one catalog entry revision
// (GET <URL>?entry=<name@rev>, with the read-only state token as its bearer).
// Nil means KIFF isn't asked: the shipped proved.json stands, as before.
type ProofState struct {
	url   string
	token string
	http  *http.Client
	now   func() time.Time

	mu    sync.Mutex
	cache map[string]proofRead
}

type proofRead struct {
	answer proofAnswer
	ok     bool
	at     time.Time
}

// proofAnswer is KIFF's answer for an entry: green or unknown, never more.
type proofAnswer struct {
	Entry  string `json:"entry"`
	App    string `json:"app"`
	State  string `json:"state"`
	Since  string `json:"since"`
	Reason string `json:"reason"`
}

// NewProofState is the reader for KIFF's state at rawURL with token. It
// returns nil, so nothing is asked, when either is empty or the address isn't
// an HTTPS one (HTTP only to a local address, for tests), has user info, a
// query or a fragment. The token is only ever sent to that one address.
func NewProofState(rawURL, token string) *ProofState {
	rawURL, token = strings.TrimSpace(rawURL), strings.TrimSpace(token)
	u, err := url.Parse(rawURL)
	if err != nil || rawURL == "" || token == "" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil
	}
	local := u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "::1"
	if u.Scheme != "https" && !(u.Scheme == "http" && local) {
		return nil
	}
	return &ProofState{
		url: rawURL, token: token, now: time.Now, cache: map[string]proofRead{},
		http: &http.Client{
			Timeout: proofStateTimeout,
			// The token goes to the address set and nowhere else.
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

// validProofEntry is KIFF's own shape for an entry name (name@rev); a
// connector's id, with its colons, is not one.
func validProofEntry(entry string) bool {
	name, rev, ok := strings.Cut(entry, "@")
	if !ok || name == "" || rev == "" || len(name) > 64 || len(rev) > 64 || !(name[0] >= 'a' && name[0] <= 'z' || name[0] >= '0' && name[0] <= '9') {
		return false
	}
	for _, r := range name {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' || r == '.' || r == '-') {
			return false
		}
	}
	for _, r := range rev {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '.' || r == '-') {
			return false
		}
	}
	return true
}

// Greens is, for each of the entries KIFF says is green, its answer. An entry
// KIFF doesn't call green, or that couldn't be read, isn't in it.
func (p *ProofState) Greens(ctx context.Context, entries []string) map[string]proofAnswer {
	out := map[string]proofAnswer{}
	if p == nil {
		return out
	}
	var ask []string
	for _, e := range entries {
		if validProofEntry(e) && !slices.Contains(ask, e) {
			ask = append(ask, e)
		}
	}
	if len(ask) > proofStateMost {
		ask = ask[:proofStateMost]
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	gate := make(chan struct{}, proofStateWidth)
	for _, e := range ask {
		wg.Add(1)
		gate <- struct{}{}
		go func() {
			defer func() { <-gate; wg.Done() }()
			if a, ok := p.read(ctx, e); ok && a.State == "green" {
				mu.Lock()
				out[e] = a
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return out
}

func (p *ProofState) read(ctx context.Context, entry string) (proofAnswer, bool) {
	now := p.now()
	p.mu.Lock()
	c, hit := p.cache[entry]
	p.mu.Unlock()
	if hit {
		ttl := proofStateTTL
		if !c.ok {
			ttl = proofStateFailTTL
		}
		if now.Sub(c.at) < ttl {
			return c.answer, c.ok
		}
	}
	a, ok := p.fetch(ctx, entry)
	p.mu.Lock()
	p.cache[entry] = proofRead{answer: a, ok: ok, at: now}
	p.mu.Unlock()
	return a, ok
}

func (p *ProofState) fetch(ctx context.Context, entry string) (proofAnswer, bool) {
	ctx, cancel := context.WithTimeout(ctx, proofStateTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url+"?entry="+url.QueryEscape(entry), nil)
	if err != nil {
		return proofAnswer{}, false
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Accept", "application/json")
	res, err := p.http.Do(req)
	if err != nil {
		return proofAnswer{}, false // the error can carry the address: never logged here
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return proofAnswer{}, false
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, proofStateBody+1))
	if err != nil || len(body) > proofStateBody {
		return proofAnswer{}, false
	}
	var a proofAnswer
	if err := json.Unmarshal(body, &a); err != nil || a.Entry != entry || (a.State != "green" && a.State != "unknown") {
		return proofAnswer{}, false
	}
	return a, true
}

// provedMap is the proofs a reader sees. When KIFF is asked, they are KIFF's
// green entries alone (the newest revision of each reviewed entry the catalog
// holds), keyed by the entry's folded name and reaching that entry, and the
// shipped proved.json is not shown at all: an entry KIFF has no green run for
// is not proved. Otherwise it is s.Proved, or proved.json, as before.
func (s *Server) provedMap(ctx context.Context) map[string]Proof {
	if s.ProofState == nil {
		if s.Proved != nil {
			return s.Proved
		}
		return Proved
	}
	entries := newestRevisions(s.catalogEntries())
	ids := make([]string, 0, len(entries))
	for id := range entries {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	out := map[string]Proof{}
	for id, a := range s.ProofState.Greens(ctx, ids) {
		name, _, _ := strings.Cut(id, "@")
		if !provedName.MatchString(provedKey(name)) {
			continue
		}
		how := "oauth"
		if entries[id].Kind == KindAPIKey {
			how = "key"
		}
		at := a.Since
		if len(at) >= 10 {
			at = at[:10]
		}
		if _, err := time.Parse(time.DateOnly, at); err != nil {
			continue
		}
		// By the reviewed entry's own vendor (planque's words, never the
		// provider's), reaching each of its green entries.
		key := provedKey(entries[id].Vendor)
		if !provedName.MatchString(key) {
			key = provedKey(name)
		}
		have, there := out[key]
		if there && have.How != how {
			continue
		}
		if !there || at < have.At {
			have.At = at
		}
		have.How = how
		if !slices.Contains(have.Entries, name) {
			have.Entries = append(have.Entries, name)
		}
		out[key] = have
	}
	return out
}
