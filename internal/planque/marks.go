package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/store"
)

// A managed mark says a vendor's apps may be connected through the
// connection platform: planque's gate for managed connections reads it, and
// the catalog read says which reviewed entries it covers. It is kept in
// planque's table, never in an entry (an entry can't change once published,
// a mark must be able to), and it is set only by an enrolled operator's
// device, signed, so no session, credential or generator can set one.
//
// A mark names the vendor and the exact hosts its apps may use (the same
// exact form as the review checks' vendorHosts: a host, or "*.domain" for a
// vendor that serves one per customer). An entry is served as managed only
// when its vendor is marked and every endpoint it has is on those hosts, so
// an entry that calls itself by a managed vendor's name but points elsewhere
// is not managed, and a revision that adds a host isn't until the owner
// signs the mark again with it.
//
// Withdrawing a mark stops new managed connections at once (within the
// seconds a planque instance keeps its read). It revokes nothing: what happens
// to a connection already made stays the owner's decision. managed_apps.json
// is the seed: a vendor with no mark follows it, and a withdrawn mark wins
// over it.

const (
	markPK     = "MARKS"
	markActive = "active"
	markGone   = "withdrawn"
)

// markHost is one exact host, or "*." and a domain of two labels or more.
var markHost = regexp.MustCompile(`^(\*\.)?([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)

// markFollows is the date of the published list's update a mark follows.
var markFollows = regexp.MustCompile(`^20[0-9]{2}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)

// managedMark is a mark as kept and as signed.
type managedMark struct {
	Vendor  string         `json:"vendor"`  // the folded name (provedKey)
	Name    string         `json:"name"`    // as people know it ("Google Sheets")
	Hosts   []string       `json:"hosts"`   // exact, sorted
	Follows string         `json:"follows"` // the published list's date
	State   string         `json:"state"`
	Digest  string         `json:"digest"` // what the device signed
	Signed  *catalogSigned `json:"signed"`
	Seq     int            `json:"seq"`
}

// markDigest is what an operator's device signs to set a mark: the SHA-256 of
// its fields in a fixed form.
func markDigest(name string, hosts []string, follows string) string {
	hosts = slices.Clone(hosts)
	sort.Strings(hosts)
	doc, _ := json.Marshal(struct {
		Name    string   `json:"name"`
		Hosts   []string `json:"hosts"`
		Follows string   `json:"follows"`
	}{name, hosts, follows})
	sum := sha256.Sum256(doc)
	return hex.EncodeToString(sum[:])
}

// markCache is the marks as one planque instance last read them.
type markCache struct {
	mu   sync.Mutex
	at   time.Time
	read map[string]managedMark
}

// markTTL is how long an instance keeps what it read: a mark set or withdrawn
// on another instance reaches it within this.
const markTTL = 5 * time.Second

func (s *Server) forgetMarks() {
	s.markCache.mu.Lock()
	s.markCache.read = nil
	s.markCache.mu.Unlock()
}

// managedMarks are the current marks by vendor.
func (s *Server) managedMarks() map[string]managedMark {
	c := &s.markCache
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.read != nil && s.now().Sub(c.at) < markTTL {
		return c.read
	}
	out := map[string]managedMark{}
	// Read for the gate, which has no request of its own to read it under.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items, err := s.KV.Query(ctx, markPK, "MANAGED#")
	if err != nil {
		// Unreadable: keep what was read, else follow the seed alone.
		if c.read != nil {
			return c.read
		}
		return out
	}
	for _, it := range items {
		var m managedMark
		if json.Unmarshal(it.Doc, &m) == nil && "MANAGED#"+m.Vendor == it.SK && s.markSigned(m) {
			out[m.Vendor] = m
		}
	}
	c.read, c.at = out, s.now()
	return out
}

// markSigned: the mark's signature holds, with the key of an operator still
// enrolled: a row changed after it was signed, or written by anything but an
// enrolled device, isn't a mark.
func (s *Server) markSigned(m managedMark) bool {
	action := "mark-managed"
	if m.State == markGone {
		action = "unmark-managed"
	}
	return m.Signed != nil && m.Signed.Digest == m.Digest && s.catalogSignedBy(m.Signed, action, m.Vendor, m.Digest) &&
		(m.State == markGone || markDigest(m.Name, m.Hosts, m.Follows) == m.Digest)
}

// managedByMark says whether vendor's managed connections are marked on (true),
// marked off (false), or not marked (ok false: the seed decides).
func (s *Server) managedByMark(vendor string) (on, ok bool) {
	m, found := s.managedMarks()[provedKey(vendor)]
	if !found {
		return false, false
	}
	return m.State == markActive, true
}

// hostOnMark reports whether raw is an https address on one of hosts: exact,
// no port, no user info; "*.domain" takes subdomains of it.
func hostOnMark(raw string, hosts []string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	for _, h := range hosts {
		if suffix, ok := strings.CutPrefix(h, "*."); ok {
			if strings.HasSuffix(host, "."+suffix) {
				return true
			}
		} else if host == h {
			return true
		}
	}
	return false
}

// entryEndpoints are every address a reviewed entry sends a call or a person
// to, with an account template read with a stand-in label.
func (s *Server) entryEndpoints(e CatalogEntry) []string {
	var out []string
	add := func(a string) {
		if a != "" {
			out = append(out, a)
		}
	}
	switch {
	case e.Application != "":
		app, ok := s.Applications[e.Application]
		if !ok {
			return nil
		}
		add(app.Server)
		for _, ep := range []string{app.OAuth.Issuer, app.OAuth.Authorization, app.OAuth.Token, app.OAuth.Registration, app.OAuth.Revocation} {
			add(ep)
		}
	case e.API != nil:
		if len(e.API.Bases) > 0 {
			out = append(out, e.API.Bases...)
		} else if accountTemplate(e.API.Base) {
			add(sampleBase(e.API.Base))
		} else {
			add(e.API.Base)
		}
	default:
		add(e.Server)
	}
	return out
}

// entryManaged says whether the reviewed entry is served as managed: its
// vendor has an active mark and every endpoint it has is on the mark's hosts.
func (s *Server) entryManaged(e CatalogEntry) bool {
	m, ok := s.managedMarks()[provedKey(e.Vendor)]
	if !ok || m.State != markActive {
		return false
	}
	endpoints := s.entryEndpoints(e)
	if len(endpoints) == 0 {
		return false
	}
	for _, ep := range endpoints {
		if !hostOnMark(ep, m.Hosts) {
			return false
		}
	}
	return true
}

// markSignature is what a request to set or withdraw a mark carries.
type markSignature struct {
	Operator  string `json:"operator"`
	Nonce     string `json:"nonce"`
	SignedAt  string `json:"signed_at"`
	Signature string `json:"signature"`
}

// verifyMark checks the operator's device signed action over (vendor, digest),
// recently, once, and returns the operator. It writes the refusal itself.
func (s *Server) verifyMark(w http.ResponseWriter, r *http.Request, sig markSignature, action, vendor, digest string) (CatalogOperator, bool) {
	if !catalogNonce.MatchString(sig.Nonce) {
		writeErr(w, http.StatusBadRequest, "a nonce of 16 to 64 letters, digits, _ or -")
		return CatalogOperator{}, false
	}
	op, ok := s.CatalogOperators[sig.Operator]
	if !ok || op.key == nil {
		writeErr(w, http.StatusForbidden, "no enrolled device for "+sig.Operator)
		return CatalogOperator{}, false
	}
	msg := CatalogSignMessage(s.BaseURL, action, vendor, digest, sig.Nonce, sig.SignedAt)
	if err := device.Verify(op.key, msg, sig.SignedAt, sig.Signature, s.now()); err != nil {
		writeErr(w, http.StatusForbidden, err.Error())
		return CatalogOperator{}, false
	}
	// The nonce is spent first: a signature counts once, even if what follows fails.
	nonce := store.Item{PK: markPK, SK: "NONCE#" + sig.Nonce, Doc: []byte(`{}`), TTL: s.now().Add(2 * catalogSignWindow).Unix()}
	if err := s.KV.Put(r.Context(), nonce, store.Cond{NotExists: true}); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, "this signature was used already")
		return CatalogOperator{}, false
	} else if err != nil {
		s.fail(w, r, err)
		return CatalogOperator{}, false
	}
	return op, true
}

// setManagedMark: POST /v1/marks/managed {name, hosts, follows, operator,
// nonce, signed_at, signature}. The owner's device signs the mark's digest;
// planque computes it from what was sent, so a signature covers exactly it.
func (s *Server) setManagedMark(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name    string   `json:"name"`
		Hosts   []string `json:"hosts"`
		Follows string   `json:"follows"`
		markSignature
	}
	if json.NewDecoder(io.LimitReader(r.Body, 8192)).Decode(&in) != nil {
		writeErr(w, http.StatusBadRequest, "send {name, hosts, follows, operator, nonce, signed_at, signature}")
		return
	}
	vendor := provedKey(in.Name)
	if vendor == "" || label(in.Name) != in.Name || len(in.Name) > 128 {
		writeErr(w, http.StatusBadRequest, "name is the vendor as people know it, on one line")
		return
	}
	if len(in.Hosts) == 0 || len(in.Hosts) > 16 || !markFollows.MatchString(in.Follows) {
		writeErr(w, http.StatusBadRequest, "hosts: 1 to 16 exact hosts; follows: the published list's date, YYYY-MM-DD")
		return
	}
	for i, h := range in.Hosts {
		if h != strings.ToLower(h) || !markHost.MatchString(h) || slices.Contains(in.Hosts[:i], h) {
			writeErr(w, http.StatusBadRequest, "each host is distinct, lowercase, exact (a host, or *.domain of two labels or more)")
			return
		}
	}
	hosts := slices.Clone(in.Hosts)
	sort.Strings(hosts)
	digest := markDigest(in.Name, hosts, in.Follows)
	op, ok := s.verifyMark(w, r, in.markSignature, "mark-managed", vendor, digest)
	if !ok {
		return
	}
	s.putMark(w, r, &managedMark{Vendor: vendor, Name: in.Name, Hosts: hosts, Follows: in.Follows, State: markActive, Digest: digest,
		Signed: &catalogSigned{Operator: op.ID, Device: device.Fingerprint(op.key), Digest: digest, Nonce: in.Nonce, SignedAt: in.SignedAt, Signature: in.Signature}})
}

// withdrawManagedMark: POST /v1/marks/managed/{vendor}/withdraw {operator,
// nonce, signed_at, signature}, over the mark's digest as it is.
func (s *Server) withdrawManagedMark(w http.ResponseWriter, r *http.Request) {
	var in markSignature
	if json.NewDecoder(io.LimitReader(r.Body, 8192)).Decode(&in) != nil {
		writeErr(w, http.StatusBadRequest, "send {operator, nonce, signed_at, signature}")
		return
	}
	vendor := provedKey(r.PathValue("vendor"))
	cur, ok := s.managedMarks()[vendor]
	if !ok || cur.State != markActive {
		writeErr(w, http.StatusConflict, "no active managed mark for "+r.PathValue("vendor"))
		return
	}
	op, ok := s.verifyMark(w, r, in, "unmark-managed", vendor, cur.Digest)
	if !ok {
		return
	}
	gone := cur
	gone.State = markGone
	gone.Signed = &catalogSigned{Operator: op.ID, Device: device.Fingerprint(op.key), Digest: cur.Digest, Nonce: in.Nonce, SignedAt: in.SignedAt, Signature: in.Signature}
	s.putMark(w, r, &gone)
}

// putMark keeps the mark as current and, beside it, as one more line of the
// history, which is only ever added to.
func (s *Server) putMark(w http.ResponseWriter, r *http.Request, m *managedMark) {
	it, found, err := s.KV.Get(r.Context(), markPK, "MANAGED#"+m.Vendor)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	cond := store.Cond{NotExists: true}
	if found {
		var prev managedMark
		_ = json.Unmarshal(it.Doc, &prev)
		m.Seq = prev.Seq + 1
		ver := it.Ver
		cond = store.Cond{Ver: &ver}
	} else {
		m.Seq = 1
	}
	doc, _ := json.Marshal(m)
	if err := s.KV.Put(r.Context(), store.Item{PK: markPK, SK: "MANAGED#" + m.Vendor, Doc: doc}, cond); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, "the mark changed meanwhile: read it again")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	at := s.now().UTC().Format(time.RFC3339Nano)
	if err := s.KV.Put(r.Context(), store.Item{PK: markPK, SK: "LOG#" + at + "#" + m.Vendor, Doc: doc}, store.Cond{NotExists: true}); err != nil {
		s.log().ErrorContext(r.Context(), "mark history not kept", "vendor", m.Vendor, "err", err.Error())
	}
	s.log().InfoContext(r.Context(), "managed mark "+m.State, "vendor", m.Vendor, "seq", m.Seq, "operator", m.Signed.Operator, "digest", m.Digest)
	s.forgetMarks()
	writeJSON(w, http.StatusOK, markView(*m))
}

// markView is a mark as read: no credential is in it.
func markView(m managedMark) map[string]any {
	by := ""
	if m.Signed != nil {
		by = m.Signed.Operator
	}
	return map[string]any{"vendor": m.Vendor, "name": m.Name, "hosts": m.Hosts, "follows": m.Follows, "state": m.State, "seq": m.Seq, "signed_by": by}
}

// listMarks: GET /v1/marks, the managed marks as planque holds them.
func (s *Server) listMarks(w http.ResponseWriter, r *http.Request) {
	marks := s.managedMarks()
	out := make([]map[string]any, 0, len(marks))
	for _, m := range marks {
		out = append(out, markView(m))
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["vendor"].(string) < out[j]["vendor"].(string) })
	writeJSON(w, http.StatusOK, map[string]any{"managed": out})
}
