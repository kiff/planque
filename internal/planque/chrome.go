package planque

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

// The owner's own Chrome as a browser for their tafs: a tafs extension,
// paired by the owner's phone, that planque sends checked steps to. This
// file is the pairing. Pairing starts on the phone, which keeps a secret it
// never sends; the extension and the phone each prove the other's key with
// it, so planque relays both keys and can substitute neither.
//
//  1. The phone opens a pairing (signed): planque answers its lookup code.
//  2. The extension, given lookup and secret by the owner, posts its key
//     with a proof made with the secret. planque can't check the proof.
//  3. The phone reads the pairing, checks the proof, and confirms it
//     (signed), naming the Chrome and adding its own key and proofs.
//  4. The extension reads the confirmation with its key, and checks the
//     phone's proof and signature itself.
//
// Off unless ChromeBuilds lists at least one released build: an
// unreleased build can't pair or connect.

// chromePairingTTL is how long a pairing code works, from the phone.
const chromePairingTTL = 10 * time.Minute

// lookupAlphabet has no look-alike characters (no 0/O, 1/I/L, U/V): the
// extension and the phones read codes in exactly these 30.
const lookupAlphabet = "23456789ABCDEFGHJKMNPQRSTVWXYZ"

var (
	lookupPattern  = regexp.MustCompile(`^[2-9A-HJKMNP-TV-Z]{4}$`)
	manifestHashRe = regexp.MustCompile(`^[0-9a-f]{64}$`)
	chromeIDRe     = regexp.MustCompile(`^chr_[A-Za-z0-9_-]{8,40}$`)
)

// ChromeBuild is one released build of the extension: the SHA-256 of its
// packed manifest, which the extension sends at pairing and on each
// connect, and its version.
type ChromeBuild struct {
	ManifestSHA256 string
	Version        string
}

// ChromeBuildsFrom reads the release allowlist: one build per line,
// "<manifest sha256> <version>", # comments and blank lines skipped.
func ChromeBuildsFrom(r io.Reader) (map[string]ChromeBuild, error) {
	out := map[string]ChromeBuild{}
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if i := strings.Index(line, "#"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 2 || !manifestHashRe.MatchString(f[0]) {
			return nil, fmt.Errorf("chrome builds line %d: want \"<manifest sha256, lower-case hex> <version>\"", n)
		}
		out[f[0]] = ChromeBuild{ManifestSHA256: f[0], Version: f[1]}
	}
	return out, sc.Err()
}

// ChromePairing is a pairing while it's made, kept until it's confirmed
// and read by the extension, or for chromePairingTTL.
type ChromePairing struct {
	ID      string    `json:"id"`
	Tenant  string    `json:"tenant"`
	Lookup  string    `json:"lookup"`
	Expires time.Time `json:"expires"`
	// What the extension posted (step 2).
	Key          string `json:"key,omitempty"` // SPKI DER, base64url
	Version      string `json:"version,omitempty"`
	ManifestHash string `json:"manifest_hash,omitempty"`
	Browser      string `json:"browser,omitempty"`
	OS           string `json:"os,omitempty"`
	Proof        string `json:"proof,omitempty"`
	// Contested: a second Chrome posted on this code. The code is spent:
	// the phone is told and opens a new one, and no Chrome collects it.
	Contested bool `json:"contested,omitempty"`
	// What the phone confirmed (step 3).
	Confirmed *ChromeConfirmation `json:"confirmed,omitempty"`
}

// ChromeConfirmation is the phone's confirmation, relayed to the extension
// as the phone wrote it.
type ChromeConfirmation struct {
	ChromeID   string `json:"chrome_id"`
	Name       string `json:"name"`
	PhoneKey   string `json:"phone_key"`
	PhoneProof string `json:"phone_proof"`
	PhoneSig   string `json:"phone_sig"`
}

// Chrome is a paired Chrome: its key, and what it said it was.
type Chrome struct {
	ID           string    `json:"id"`
	Tenant       string    `json:"tenant"`
	Name         string    `json:"name"`
	Key          string    `json:"key"`
	KeySHA256    string    `json:"key_sha256"`
	Version      string    `json:"version"`
	ManifestHash string    `json:"manifest_hash"`
	Browser      string    `json:"browser"`
	OS           string    `json:"os"`
	PairedAt     time.Time `json:"paired_at"`
}

func skChromePairing(id string) string { return "CHROMEPAIRING#" + id }
func skChrome(id string) string        { return "CHROME#" + id }
func pkChromeLookup(lookup string) string {
	return "CHROMELOOKUP#" + lookup
}

// chromeLookupRef is the lookup code's index: which pairing it opens.
type chromeLookupRef struct {
	Tenant  string `json:"tenant"`
	Pairing string `json:"pairing"`
}

// chromeOn says the owner's Chrome is offered: a released build is listed.
func (s *Server) chromeOn() bool { return len(s.ChromeBuilds) > 0 }

func (s *Server) chromeOff(w http.ResponseWriter) bool {
	if s.chromeOn() {
		return false
	}
	writeErr(w, http.StatusNotFound, "a browser in your own Chrome isn't offered here")
	return true
}

// The extension's two routes carry no signature, so each is limited: per
// address, so one script can't sweep the codes, and in all, so many
// addresses can't either. A code is 4 of 30 characters (810,000) and lives
// ten minutes.
const (
	chromePostsPerAddress = 10  // key posts per address in chromePostWindow
	chromePostsInAll      = 300 // key posts across every address per minute
	chromeReadsPerAddress = 60  // collects per address per minute
	chromePostWindow      = 10 * time.Minute
)

// chromeCounted counts one request under name in the window it falls in,
// and says whether it is within max. A count that can't be kept refuses.
func (s *Server) chromeCounted(r *http.Request, name string, window time.Duration, most int) bool {
	now := s.now().UTC()
	start := now.Truncate(window)
	key := store.Item{PK: "CHROMELIMIT#" + name, SK: start.Format(time.RFC3339), TTL: start.Add(2 * window).Unix()}
	for range 5 {
		it, ok, err := s.KV.Get(r.Context(), key.PK, key.SK)
		if err != nil {
			return false
		}
		var n struct {
			Count int `json:"count"`
		}
		if ok && json.Unmarshal(it.Doc, &n) != nil {
			n.Count = 0
		}
		if n.Count >= most {
			return false
		}
		n.Count++
		cond := store.Cond{NotExists: true}
		if ok {
			ver := it.Ver
			cond, key.Ver = store.Cond{Ver: &ver}, it.Ver+1
		}
		if err := putDoc(r.Context(), s.KV, key, n, cond); err == nil {
			return true
		} else if !errors.Is(err, store.ErrCondition) {
			return false
		}
	}
	return false
}

// chromeAddress is the address a request came from, without its port.
func chromeAddress(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

func (s *Server) chromeTooMany(w http.ResponseWriter) {
	writeErr(w, http.StatusTooManyRequests, "too many tries from this network: wait a few minutes, then start again on your phone")
}

// spkiP256 says key is a P-256 public key, SPKI DER in base64url, and gives
// its SHA-256 in hex.
func spkiP256(key string) (string, error) {
	der, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil || len(der) > 200 {
		return "", errors.New("key must be an SPKI public key in base64url")
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return "", errors.New("key must be an SPKI public key in base64url")
	}
	if k, ok := pub.(*ecdsa.PublicKey); !ok || k.Curve != elliptic.P256() {
		return "", errors.New("key must be an ECDSA P-256 key")
	}
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:]), nil
}

// newLookup draws each character evenly: a byte at or above the largest
// multiple of 30 is drawn again, so no character is likelier than another.
func newLookup() string {
	return lookupFrom(func() byte {
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic(err)
		}
		return b[0]
	})
}

func lookupFrom(next func() byte) string {
	top := 256 - 256%len(lookupAlphabet)
	b := make([]byte, 0, 4)
	for len(b) < 4 {
		if c := int(next()); c < top {
			b = append(b, lookupAlphabet[c%len(lookupAlphabet)])
		}
	}
	return string(b)
}

// openChromePairing is step 1: the owner's phone opens a pairing.
func (s *Server) openChromePairing(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOff(w) {
		return
	}
	now := s.now()
	p := ChromePairing{ID: "pr_" + oauth.Random(12), Tenant: sr.tenant, Expires: now.Add(chromePairingTTL).UTC()}
	// A lookup code is unique among live pairings: the index is written
	// first, conditionally, and a taken code is drawn again.
	for try := 0; ; try++ {
		p.Lookup = newLookup()
		err := putDoc(r.Context(), s.KV, store.Item{PK: pkChromeLookup(p.Lookup), SK: "-", TTL: p.Expires.Unix()},
			chromeLookupRef{Tenant: p.Tenant, Pairing: p.ID}, store.Cond{NotExists: true})
		if err == nil {
			break
		}
		if !errors.Is(err, store.ErrCondition) || try == 4 {
			s.fail(w, r, err)
			return
		}
	}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(p.Tenant), SK: skChromePairing(p.ID), TTL: p.Expires.Unix()}, p, store.Cond{NotExists: true}); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"pairing": p.ID, "lookup": p.Lookup, "expires_at": p.Expires})
}

// loadChromeLookup finds the live pairing a lookup code opens.
func (s *Server) loadChromeLookup(r *http.Request, lookup string) (ChromePairing, store.Item, error) {
	var p ChromePairing
	if !lookupPattern.MatchString(lookup) {
		return p, store.Item{}, errNotFound
	}
	var ref chromeLookupRef
	if _, err := getDoc(r.Context(), s.KV, pkChromeLookup(lookup), "-", &ref); err != nil {
		return p, store.Item{}, err
	}
	it, err := getDoc(r.Context(), s.KV, tenantPK(ref.Tenant), skChromePairing(ref.Pairing), &p)
	if err != nil {
		return p, it, err
	}
	if !s.now().Before(p.Expires) || p.Lookup != lookup {
		return p, it, errNotFound
	}
	return p, it, nil
}

// postChromeKey is step 2: the extension posts its key, once. The first
// extension to use a code holds it; a second is refused, and the phone,
// checking the proof, refuses any key not made with its secret.
func (s *Server) postChromeKey(w http.ResponseWriter, r *http.Request) {
	if s.chromeOff(w) {
		return
	}
	if !s.chromeCounted(r, "post#"+chromeAddress(r), chromePostWindow, chromePostsPerAddress) ||
		!s.chromeCounted(r, "post", time.Minute, chromePostsInAll) {
		s.chromeTooMany(w)
		return
	}
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		Key          string `json:"key"`
		Version      string `json:"version"`
		ManifestHash string `json:"manifest_hash"`
		Browser      string `json:"browser"`
		OS           string `json:"os"`
		Proof        string `json:"proof"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	if _, err := spkiP256(req.Key); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if b, ok := s.ChromeBuilds[req.ManifestHash]; !ok || b.Version != req.Version {
		writeErr(w, http.StatusForbidden, "this build of the extension isn't a released one: install it from the Chrome Web Store")
		return
	}
	if proof, err := base64.RawURLEncoding.DecodeString(req.Proof); err != nil || len(proof) != sha256.Size {
		writeErr(w, http.StatusBadRequest, "proof must be an HMAC-SHA256 in base64url")
		return
	}
	p, it, err := s.loadChromeLookup(r, r.PathValue("lookup"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no pairing has that code, or it expired: start again on your phone")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	// A code a Chrome already used is spent by a second post: the phone is
	// told, and no Chrome collects it. The answer is the unknown code's,
	// so a sweep can't tell a used code from no code.
	if p.Key != "" || p.Contested {
		if !p.Contested {
			p.Contested = true
			if err := putDoc(r.Context(), s.KV, store.Item{PK: it.PK, SK: it.SK, TTL: p.Expires.Unix()}, p, store.Cond{Ver: &it.Ver}); err != nil && !errors.Is(err, store.ErrCondition) {
				s.fail(w, r, err)
				return
			}
		}
		writeErr(w, http.StatusNotFound, "no pairing has that code, or it expired: start again on your phone")
		return
	}
	p.Key, p.Version, p.ManifestHash, p.Proof = req.Key, req.Version, req.ManifestHash, req.Proof
	p.Browser, p.OS = oneLine(req.Browser, 60), oneLine(req.OS, 60)
	if err := putDoc(r.Context(), s.KV, store.Item{PK: it.PK, SK: it.SK, TTL: p.Expires.Unix()}, p, store.Cond{Ver: &it.Ver}); errors.Is(err, store.ErrCondition) {
		// Another post won the race: this one is the second.
		writeErr(w, http.StatusNotFound, "no pairing has that code, or it expired: start again on your phone")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"status": "waiting_for_phone"})
}

// readChromePairing is the phone reading a pairing, to check the proof.
func (s *Server) readChromePairing(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOff(w) {
		return
	}
	var p ChromePairing
	if _, err := getDoc(r.Context(), s.KV, tenantPK(sr.tenant), skChromePairing(r.PathValue("pairing")), &p); errors.Is(err, errNotFound) || err == nil && !s.now().Before(p.Expires) {
		writeErr(w, http.StatusNotFound, "no such pairing, or it expired")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	out := map[string]any{"pairing": p.ID, "expires_at": p.Expires, "status": "waiting_for_chrome"}
	if p.Key != "" {
		sum, _ := spkiP256(p.Key)
		out["status"], out["key"], out["key_sha256"], out["proof"] = "waiting_for_phone", p.Key, sum, p.Proof
		out["version"], out["browser"], out["os"] = p.Version, p.Browser, p.OS
	}
	if p.Confirmed != nil {
		out["status"], out["chrome_id"] = "confirmed", p.Confirmed.ChromeID
	}
	if p.Contested && p.Confirmed == nil {
		out = map[string]any{"pairing": p.ID, "expires_at": p.Expires, "status": "contested"}
	}
	writeJSON(w, http.StatusOK, out)
}

// confirmChromePairing is step 3: the phone confirms the Chrome whose key
// it checked, and planque keeps the Chrome for the account.
func (s *Server) confirmChromePairing(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOff(w) {
		return
	}
	var req struct {
		ChromeKeySHA256 string `json:"chrome_key_sha256"`
		Name            string `json:"name"`
		PhoneKey        string `json:"phone_key"`
		PhoneProof      string `json:"phone_proof"`
		PhoneSig        string `json:"phone_sig"`
	}
	if err := decodeStrict(sr.body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	name := oneLine(req.Name, 40)
	if name == "" {
		writeErr(w, http.StatusBadRequest, "name the Chrome (\"Work Mac\")")
		return
	}
	if _, err := spkiP256(req.PhoneKey); err != nil {
		writeErr(w, http.StatusBadRequest, "phone_key: "+err.Error())
		return
	}
	if proof, err := base64.RawURLEncoding.DecodeString(req.PhoneProof); err != nil || len(proof) != sha256.Size {
		writeErr(w, http.StatusBadRequest, "phone_proof must be an HMAC-SHA256 in base64url")
		return
	}
	if sig, err := base64.RawURLEncoding.DecodeString(req.PhoneSig); err != nil || len(sig) != 64 {
		writeErr(w, http.StatusBadRequest, "phone_sig must be a P-256 signature, r‖s in base64url")
		return
	}
	var p ChromePairing
	it, err := getDoc(r.Context(), s.KV, tenantPK(sr.tenant), skChromePairing(r.PathValue("pairing")), &p)
	if errors.Is(err, errNotFound) || err == nil && !s.now().Before(p.Expires) {
		writeErr(w, http.StatusNotFound, "no such pairing, or it expired")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	if p.Key == "" {
		writeErr(w, http.StatusConflict, "no Chrome has used this code yet")
		return
	}
	if p.Confirmed != nil {
		writeErr(w, http.StatusConflict, "this pairing is already confirmed")
		return
	}
	if p.Contested {
		writeErr(w, http.StatusConflict, "another Chrome also used this code, so it can't be trusted: start again on your phone")
		return
	}
	// The phone names the key it checked the proof against: another key
	// posted in between is never confirmed.
	sum, _ := spkiP256(p.Key)
	if req.ChromeKeySHA256 != sum {
		writeErr(w, http.StatusConflict, "the Chrome's key isn't the one you checked: start again on your phone")
		return
	}
	c := Chrome{ID: "chr_" + oauth.Random(12), Tenant: sr.tenant, Name: name, Key: p.Key, KeySHA256: sum,
		Version: p.Version, ManifestHash: p.ManifestHash, Browser: p.Browser, OS: p.OS, PairedAt: s.now().UTC()}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(c.Tenant), SK: skChrome(c.ID)}, c, store.Cond{NotExists: true}); err != nil {
		s.fail(w, r, err)
		return
	}
	p.Confirmed = &ChromeConfirmation{ChromeID: c.ID, Name: name, PhoneKey: req.PhoneKey, PhoneProof: req.PhoneProof, PhoneSig: req.PhoneSig}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: it.PK, SK: it.SK, TTL: p.Expires.Unix()}, p, store.Cond{Ver: &it.Ver}); err != nil {
		// The Chrome was kept but the extension can't learn it: it goes.
		_, _, _ = s.KV.Delete(r.Context(), tenantPK(c.Tenant), skChrome(c.ID))
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chrome_id": c.ID, "name": name})
}

// collectChromePairing is step 4: the extension, with the key it posted,
// reads the phone's confirmation once, and the pairing ends.
func (s *Server) collectChromePairing(w http.ResponseWriter, r *http.Request) {
	if s.chromeOff(w) {
		return
	}
	if !s.chromeCounted(r, "read#"+chromeAddress(r), time.Minute, chromeReadsPerAddress) {
		s.chromeTooMany(w)
		return
	}
	p, it, err := s.loadChromeLookup(r, r.PathValue("lookup"))
	// A pairing that isn't this key's, or that a second Chrome contested,
	// says nothing more than an unknown one.
	if errors.Is(err, errNotFound) || err == nil && (p.Key == "" || p.Contested || r.URL.Query().Get("key") != p.Key) {
		writeErr(w, http.StatusNotFound, "no pairing for this Chrome")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	if p.Confirmed == nil {
		writeJSON(w, http.StatusAccepted, map[string]any{"status": "waiting_for_phone"})
		return
	}
	if _, _, err := s.KV.Delete(r.Context(), it.PK, it.SK); err != nil {
		s.fail(w, r, err)
		return
	}
	_, _, _ = s.KV.Delete(r.Context(), pkChromeLookup(p.Lookup), "-")
	writeJSON(w, http.StatusOK, p.Confirmed)
}

// listChromes is the account's paired Chromes, for the phone.
func (s *Server) listChromes(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOff(w) {
		return
	}
	items, err := s.KV.Query(r.Context(), tenantPK(sr.tenant), "CHROME#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	out := []map[string]any{}
	for _, it := range items {
		var c Chrome
		if json.Unmarshal(it.Doc, &c) != nil {
			continue
		}
		out = append(out, map[string]any{"chrome_id": c.ID, "name": c.Name, "browser": c.Browser, "os": c.OS, "version": c.Version, "paired_at": c.PairedAt})
	}
	writeJSON(w, http.StatusOK, map[string]any{"chromes": out})
}

// unpairChrome removes a paired Chrome: its key goes at once, so nothing
// it signs is taken again.
func (s *Server) unpairChrome(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOff(w) {
		return
	}
	id := r.PathValue("chrome")
	if !chromeIDRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "no such Chrome")
		return
	}
	_, found, err := s.KV.Delete(r.Context(), tenantPK(sr.tenant), skChrome(id))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !found {
		writeErr(w, http.StatusNotFound, "no such Chrome")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chrome_id": id, "status": "unpaired"})
}

// oneLine is s on one line, control characters removed, at most max runes.
func oneLine(s string, max int) string {
	t := strings.Join(strings.Fields(strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s)), " ")
	if r := []rune(t); len(r) > max {
		return string(r[:max])
	}
	return t
}
