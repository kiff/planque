package planque

import (
	"bufio"
	"context"
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
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/stepkey"
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

// chromeKiffOfferMax is the longest KIFF device offer an extension may pass
// to the phone: KIFF's own pairing offer is at most 4 KiB.
const chromeKiffOfferMax = 4096

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
	// KiffOffer is the extension's view-only KIFF device offer (kiff.pair.v1),
	// as its text: the phone pairs that device in the same step, and the
	// extension's proof covers it, so planque keeps it as given and checks
	// nothing in it (kiff/tafs docs/chrome.md).
	KiffOffer string `json:"kiff_offer,omitempty"`
	// ProfileID is the id the extension keeps for its Chrome profile, a hint
	// the phone uses to pre-select the Chrome a re-pairing replaces; it proves
	// nothing (kiff/tafs docs/chrome.md, "Your Chrome" 1). OneTouch says the
	// extension checks a phone_sig made with KIFF's pairing under one
	// fingerprint (5): the phone sends that only then. Both kept as given.
	ProfileID string `json:"profile_id,omitempty"`
	OneTouch  bool   `json:"one_touch,omitempty"`
	// OneTouchProof is the extension's HMAC over the flag (its own label,
	// under the code's secret, which planque never sees): the phone checks it,
	// so planque can't add the flag to an offer.
	OneTouchProof string `json:"one_touch_proof,omitempty"`
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
	// One fingerprint for the pairing and KIFF's view-only device (kiff/tafs
	// docs/chrome.md, "Your Chrome" 5): phone_sig is then the approval key's
	// kiff.requests.v1 signature, and the extension rebuilds it from these,
	// relayed as the phone sent them. All three or none.
	Tenant            string `json:"tenant,omitempty"`
	SignedAt          string `json:"signed_at,omitempty"`
	KiffRequestSHA256 string `json:"kiff_request_sha256,omitempty"`
	// Replaces is the Chrome this one replaces ("Your Chrome" 1), and Carried
	// each connection the replace moves with its sites, inside what the phone
	// signed for the extension: it keeps them as those grants' signed sites.
	Replaces string          `json:"replaces,omitempty"`
	Carried  []ChromeCarried `json:"carried,omitempty"`
}

// ChromeCarried is a connection a re-pairing moves, with its sites as the
// owner approved them.
type ChromeCarried struct {
	Grant string   `json:"grant"`
	Sites []string `json:"sites"`
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
	// ProfileID is the extension's id for its Chrome profile, as it paired
	// (a hint for the phone's re-pairing pick, "Your Chrome" 1).
	ProfileID string `json:"profile_id,omitempty"`
	// PairedWith is the re-pairing that carried connections here, as the
	// phone signed it: what vouches for those connections' sites.
	PairedWith *ChromePairedWith `json:"paired_with,omitempty"`
	// Push is where planque wakes this Chrome's extension for a step while it
	// isn't connected (chrome_push.go). Never listed.
	Push *ChromePush `json:"push,omitempty"`
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

// ChromeOffer is which accounts the owner's own Chrome is offered to: every
// account, or only Tenants. The zero value offers it to none.
type ChromeOffer struct {
	All     bool
	Tenants map[string]bool
}

// ChromeOfferFrom reads PLANQUE_CHROME_OFFER and PLANQUE_CHROME_TENANTS, as
// KIFF reads KIFF_APP_SIGNIN and its tenants: "on" offers the owner's Chrome
// to every account, "test" only to the accounts in tenants
// ("tenants:<id>[,<id>]", as PLANQUE_MANAGED_TEST_TENANTS), and "" or "off"
// to none. So a listed build alone turns nothing on, and a setting lost on
// a rerun turns it off, never on for everyone. Tenants are read only with
// "test"; "test" with none is a mistake.
func ChromeOfferFrom(offer, tenants string) (ChromeOffer, error) {
	switch strings.TrimSpace(offer) {
	case "", "off":
		return ChromeOffer{}, nil
	case "on":
		return ChromeOffer{All: true}, nil
	case "test":
		only, err := ManagedTestTenantsFrom(tenants)
		if err != nil {
			return ChromeOffer{}, fmt.Errorf("PLANQUE_CHROME_TENANTS %w", err)
		}
		if len(only) == 0 {
			return ChromeOffer{}, errors.New("PLANQUE_CHROME_OFFER=test needs PLANQUE_CHROME_TENANTS")
		}
		return ChromeOffer{Tenants: only}, nil
	}
	return ChromeOffer{}, errors.New(`PLANQUE_CHROME_OFFER is "on", "test" or "off"`)
}

// chromeOn says the owner's Chrome is offered to some account: a released
// build is listed, and the offer names every account or some.
func (s *Server) chromeOn() bool {
	return len(s.ChromeBuilds) > 0 && (s.ChromeOffer.All || len(s.ChromeOffer.Tenants) > 0)
}

// chromeOnFor says it's offered to tenant's account.
func (s *Server) chromeOnFor(tenant string) bool {
	return len(s.ChromeBuilds) > 0 && (s.ChromeOffer.All || s.ChromeOffer.Tenants[tenant])
}

func (s *Server) chromeOff(w http.ResponseWriter) bool {
	if s.chromeOn() {
		return false
	}
	chromeNotOffered(w)
	return true
}

// chromeOffFor is chromeOff for one account: another account is answered
// as if the Chrome weren't offered at all.
func (s *Server) chromeOffFor(w http.ResponseWriter, tenant string) bool {
	if s.chromeOnFor(tenant) {
		return false
	}
	chromeNotOffered(w)
	return true
}

func chromeNotOffered(w http.ResponseWriter) {
	writeErr(w, http.StatusNotFound, "a browser in your own Chrome isn't offered here")
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
	if s.chromeOffFor(w, sr.tenant) {
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
		Key           string `json:"key"`
		Version       string `json:"version"`
		ManifestHash  string `json:"manifest_hash"`
		Browser       string `json:"browser"`
		OS            string `json:"os"`
		Proof         string `json:"proof"`
		KiffOffer     string `json:"kiff_offer"`
		ProfileID     string `json:"profile_id"`
		OneTouch      bool   `json:"one_touch"`
		OneTouchProof string `json:"one_touch_proof"`
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
	if req.OneTouch || req.OneTouchProof != "" {
		if proof, err := base64.RawURLEncoding.DecodeString(req.OneTouchProof); !req.OneTouch || err != nil || len(proof) != sha256.Size {
			writeErr(w, http.StatusBadRequest, "one_touch comes with one_touch_proof, an HMAC-SHA256 in base64url")
			return
		}
	}
	if req.ProfileID != "" && !profileIDShape(req.ProfileID) {
		writeErr(w, http.StatusBadRequest, "profile_id is 32 lower-case hex digits")
		return
	}
	if len(req.KiffOffer) > chromeKiffOfferMax {
		writeErr(w, http.StatusBadRequest, "kiff_offer is at most 4 KiB, as KIFF's pairing offer is")
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
	// A pairing opened for an account the Chrome is no longer offered to
	// goes no further, and answers as an unknown code, so a sweep can't
	// tell it from no pairing.
	if !s.chromeOnFor(p.Tenant) {
		writeErr(w, http.StatusNotFound, "no pairing has that code, or it expired: start again on your phone")
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
	p.Key, p.Version, p.ManifestHash, p.Proof, p.KiffOffer = req.Key, req.Version, req.ManifestHash, req.Proof, req.KiffOffer
	p.ProfileID, p.OneTouch, p.OneTouchProof = req.ProfileID, req.OneTouch, req.OneTouchProof
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
	if s.chromeOffFor(w, sr.tenant) {
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
		if p.ProfileID != "" {
			out["profile_id"] = p.ProfileID
		}
		if p.OneTouch {
			out["one_touch"], out["one_touch_proof"] = true, p.OneTouchProof
		}
		if p.KiffOffer != "" {
			out["kiff_offer"] = p.KiffOffer
		}
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
	if s.chromeOffFor(w, sr.tenant) {
		return
	}
	var req struct {
		ChromeKeySHA256   string          `json:"chrome_key_sha256"`
		Name              string          `json:"name"`
		PhoneKey          string          `json:"phone_key"`
		PhoneProof        string          `json:"phone_proof"`
		PhoneSig          string          `json:"phone_sig"`
		Tenant            string          `json:"tenant"`
		SignedAt          string          `json:"signed_at"`
		KiffRequestSHA256 string          `json:"kiff_request_sha256"`
		Replaces          string          `json:"replaces"`
		Carried           []ChromeCarried `json:"carried"`
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
	// One fingerprint (Your Chrome 5): all three or none, the account's own,
	// each a single line, so the extension rebuilds the bytes the phone signed.
	if together := req.Tenant + req.SignedAt + req.KiffRequestSHA256; together != "" {
		if req.Tenant != sr.tenant || !oneTimeLine(req.SignedAt) || !hex64(req.KiffRequestSHA256) {
			writeErr(w, http.StatusBadRequest, "signed with KIFF's pairing: send tenant (this account), signed_at (RFC 3339) and kiff_request_sha256 (64 lower-case hex), all three")
			return
		}
	}
	// Re-pairing replaces (Your Chrome 1): the Chrome it replaces is the
	// account's, and only the connections it carries move, each on that
	// Chrome, with exactly the sites the owner approved for it.
	var replaced *Chrome
	var carried []*Grant
	if req.Replaces != "" {
		var old Chrome
		if !chromeIDRe.MatchString(req.Replaces) {
			writeErr(w, http.StatusBadRequest, "replaces must be one of this account's Chromes (chr_…)")
			return
		} else if _, err := getDoc(r.Context(), s.KV, tenantPK(sr.tenant), skChrome(req.Replaces), &old); errors.Is(err, errNotFound) {
			writeErr(w, http.StatusConflict, "the Chrome this one replaces isn't paired to this account any more: pair it as a new Chrome")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		replaced = &old
		if len(req.Carried) > chromeCarriedMax {
			writeErr(w, http.StatusBadRequest, "carried moves at most 50 connections at once")
			return
		}
		seen := map[string]bool{}
		for _, m := range req.Carried {
			g, err := s.loadGrant(r.Context(), sr.tenant, m.Grant)
			// Only a connection the owner approved for the old Chrome moves: its
			// sites signed for it (or carried to it by an earlier replace), and
			// live, never one only asked for, failed or ended (qa-kiff).
			if errors.Is(err, errNotFound) || err == nil && (seen[m.Grant] || g.Status == StatusDisconnected || g.Status == StatusPending ||
				g.Status == StatusFailed || g.Request.Browser == nil || g.Request.Browser.Driver != DriverExtension ||
				g.Request.Browser.Chrome != old.ID || g.ChromeSitesSig == "" && g.ChromeSitesCarried != old.ID ||
				!slices.Equal(g.Request.Browser.Sites, m.Sites)) {
				writeErr(w, http.StatusConflict, "carried names "+m.Grant+", which isn't a connection on the Chrome replaced with those sites: read the Chrome's connections again")
				return
			} else if err != nil {
				s.fail(w, r, err)
				return
			}
			// One an earlier re-pairing carried here has no signature of its own:
			// it moves again only if that pairing, as this phone signed it, carried
			// it with these sites. Otherwise it ends with the old Chrome, as before.
			if g.ChromeSitesSig == "" && !old.PairedWith.vouches(&old, req.PhoneKey, g) {
				writeErr(w, http.StatusConflict, "carried names "+m.Grant+", which an earlier re-pairing carried, and that pairing doesn't vouch for its sites as this phone signed them: leave it out, and it ends with the old Chrome")
				return
			}
			seen[m.Grant] = true
			carried = append(carried, g)
		}
	} else if len(req.Carried) > 0 {
		writeErr(w, http.StatusBadRequest, "carried moves connections only with replaces, the Chrome they're on")
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
		Version: p.Version, ManifestHash: p.ManifestHash, Browser: p.Browser, OS: p.OS, PairedAt: s.now().UTC(), ProfileID: p.ProfileID}
	// What vouches for the connections it carries, kept as the phone signed it
	// for the extension, so a later re-pairing can carry them again. One whose
	// signature doesn't hold is never kept: the extension refuses that pairing.
	if len(req.Carried) > 0 {
		pw := &ChromePairedWith{Body: chromePairBody{Pairing: r.PathValue("pairing"), ChromeKeySHA256: sum, Name: name, PhoneKey: req.PhoneKey,
			PhoneProof: req.PhoneProof, Carried: req.Carried}, PhoneSig: req.PhoneSig,
			Tenant: req.Tenant, SignedAt: req.SignedAt, KiffRequestSHA256: req.KiffRequestSHA256}
		if pw.holds() {
			c.PairedWith = pw
		}
	}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(c.Tenant), SK: skChrome(c.ID)}, c, store.Cond{NotExists: true}); err != nil {
		s.fail(w, r, err)
		return
	}
	// The relay finds the Chrome's account from its id alone.
	if err := putDoc(r.Context(), s.KV, store.Item{PK: pkChromeIndex(c.ID), SK: "-"}, chromeIndex{Tenant: c.Tenant}, store.Cond{NotExists: true}); err != nil {
		_, _, _ = s.KV.Delete(r.Context(), tenantPK(c.Tenant), skChrome(c.ID))
		s.fail(w, r, err)
		return
	}
	p.Confirmed = &ChromeConfirmation{ChromeID: c.ID, Name: name, PhoneKey: req.PhoneKey, PhoneProof: req.PhoneProof, PhoneSig: req.PhoneSig,
		Tenant: req.Tenant, SignedAt: req.SignedAt, KiffRequestSHA256: req.KiffRequestSHA256, Replaces: req.Replaces, Carried: req.Carried}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: it.PK, SK: it.SK, TTL: p.Expires.Unix()}, p, store.Cond{Ver: &it.Ver}); err != nil {
		// The Chrome was kept but the extension can't learn it: it goes.
		_, _, _ = s.KV.Delete(r.Context(), tenantPK(c.Tenant), skChrome(c.ID))
		_, _, _ = s.KV.Delete(r.Context(), pkChromeIndex(c.ID), "-")
		s.fail(w, r, err)
		return
	}
	// The carried connections move to the new Chrome, their sites now the
	// pairing's to vouch for; then the old Chrome goes, ending the rest (3).
	moved := []string{}
	if replaced != nil {
		for _, g := range carried {
			g.Request.Browser.Chrome, g.Request.Browser.ChromeName = c.ID, label(c.Name)
			g.ChromeSitesSig, g.ChromeSitesCarried = "", c.ID
			g.event(s.now(), "moved to the Chrome \""+c.Name+"\" that replaced \""+replaced.Name+"\", with the same sites")
			// One that can't be saved stays on the old Chrome and ends with it
			// below: fail closed, and `moved` leaves it out.
			if err := s.saveGrant(r.Context(), g); err != nil {
				s.log().ErrorContext(r.Context(), "a carried connection didn't move", slog.String("tenant", sr.tenant), slog.String("grant", g.ID), slog.String("err", err.Error()))
				continue
			}
			moved = append(moved, g.ID)
		}
		if _, err := s.removeChrome(r.Context(), sr.tenant, replaced.ID); err != nil {
			s.log().ErrorContext(r.Context(), "the replaced Chrome wasn't removed", slog.String("tenant", sr.tenant), slog.String("chrome", replaced.ID), slog.String("err", err.Error()))
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"chrome_id": c.ID, "name": name, "moved": moved})
}

// chromeCarriedMax is the most connections one re-pairing moves, as the
// extension keeps them (kiff/tafs docs/chrome.md, "Your Chrome" 1).
const chromeCarriedMax = 50

// oneTimeLine: an RFC 3339 time on one line.
func oneTimeLine(v string) bool {
	_, err := time.Parse(time.RFC3339, v)
	return err == nil && !strings.ContainsAny(v, "\r\n")
}

// hex64: a SHA-256 as 64 lower-case hex digits.
func hex64(v string) bool {
	if len(v) != 64 {
		return false
	}
	for _, c := range v {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
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
	// The same for its collect: as a pairing that isn't this key's.
	if !s.chromeOnFor(p.Tenant) {
		writeErr(w, http.StatusNotFound, "no pairing for this Chrome")
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
	// The pairing's id goes with the phone's confirmation: the phone signed
	// it, and the extension, which only knows the lookup, needs it to check
	// that signature (kiff/tafs docs/chrome.md, "Pairing").
	writeJSON(w, http.StatusOK, struct {
		Pairing string `json:"pairing"`
		*ChromeConfirmation
	}{p.ID, p.Confirmed})
}

// listChromes is the account's paired Chromes, for the phone.
func (s *Server) listChromes(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if s.chromeOffFor(w, sr.tenant) {
		return
	}
	items, err := s.KV.Query(r.Context(), tenantPK(sr.tenant), "CHROME#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	// Each Chrome's live connections ("Your Chrome" 3 and 4), decided here once
	// so no app writes the rule again: the ones unpairing it would end that
	// still matter (a failed one ends too, with nothing to lose).
	on, err := s.chromesConnections(r.Context(), sr.tenant)
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
		row := map[string]any{"chrome_id": c.ID, "name": c.Name, "browser": c.Browser, "os": c.OS, "version": c.Version, "paired_at": c.PairedAt}
		if c.ProfileID != "" {
			row["profile_id"] = c.ProfileID
		}
		if c.PairedWith != nil {
			row["paired_with"] = c.PairedWith
		}
		row["connections"] = on[c.ID]
		if on[c.ID] == nil {
			row["connections"] = []chromeConnection{}
		}
		out = append(out, row)
	}
	writeJSON(w, http.StatusOK, map[string]any{"chromes": out})
}

// unpairChrome removes a paired Chrome: its key goes at once, so nothing
// it signs is taken again.
func (s *Server) unpairChrome(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	// Any account may unpair while the Chrome is offered to some: an account
	// taken off the offer can still remove the Chrome it paired.
	if s.chromeOff(w) {
		return
	}
	id := r.PathValue("chrome")
	if !chromeIDRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "no such Chrome")
		return
	}
	found, err := s.removeChrome(r.Context(), sr.tenant, id)
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

// removeChrome removes a paired Chrome's key, so nothing it signs is taken
// again: a connect from it is told it was unpaired (4410) for a while, and a
// live connection ends now, so the extension wipes itself.
func (s *Server) removeChrome(ctx context.Context, tenant, id string) (bool, error) {
	_, found, err := s.KV.Delete(ctx, tenantPK(tenant), skChrome(id))
	if err != nil || !found {
		return false, err
	}
	_ = putDoc(ctx, s.KV, store.Item{PK: pkChromeIndex(id), SK: "-", TTL: s.now().Add(chromeUnpairedKept).Unix()},
		chromeIndex{Tenant: tenant, Unpaired: true}, store.Cond{})
	s.closeChrome(ctx, tenant, id, closeUnpaired, "this Chrome was unpaired")
	// Its acts' images go with it (chrome_screenshot.go).
	s.dropChromeShots(ctx, tenant, func(shot chromeShot) bool { return shot.Chrome == id })
	// Its connections end with it, so none is left on a Chrome that's gone
	// (kiff/tafs docs/chrome.md, "Your Chrome" 3). They come back only by
	// connecting again.
	s.endChromesGrants(ctx, tenant, id)
	return true, nil
}

// endChromesGrants disconnects every connection of the account that runs in
// the Chrome id, as the owner's disconnect does. One that can't be read or
// saved is logged and the rest still end.
func (s *Server) endChromesGrants(ctx context.Context, tenant, id string) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		s.log().ErrorContext(ctx, "an unpaired Chrome's connections couldn't be listed", slog.String("tenant", tenant), slog.String("chrome", id), slog.String("err", err.Error()))
		return
	}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.ID == "" || skGrant(g.ID) != it.SK || g.Status == StatusDisconnected {
			continue
		}
		if b := g.Request.Browser; b == nil || b.Driver != DriverExtension || b.Chrome != id {
			continue
		}
		g.ver = it.Ver
		if err := s.endGrant(ctx, &g, "disconnected: its Chrome was unpaired", "its Chrome was unpaired"); err != nil {
			s.log().ErrorContext(ctx, "an unpaired Chrome's connection didn't end", slog.String("tenant", tenant), slog.String("grant", g.ID), slog.String("err", err.Error()))
		}
	}
}

// labelUnpair is what the extension signs to unpair itself (kiff/tafs
// chrome/src/crypto.js LABEL.unpair).
const labelUnpair = "tafs-chrome/unpair/v1"

// chromeUnpairSkew is how far an unpair's signed time (at, in milliseconds)
// may be from planque's clock: an older one is refused, so a captured
// request can't unpair the Chrome paired after it under the same id.
const chromeUnpairSkew = 5 * time.Minute

// extensionUnpair is the extension unpairing itself from its side panel
// (POST /v1/extension/unpair {chrome_id, at, sig}): the same removal as the
// phone's, proved by the Chrome's own key over {chrome_id, at}
// (tafs-chrome/unpair/v1), at most chromeUnpairSkew from now. A Chrome that
// is unknown, already unpaired, or whose signature isn't its key's, gets the
// same 404, so nothing tells a stranger which Chromes exist. Done once, it
// is done: the key is gone, so a replay finds nothing.
func (s *Server) extensionUnpair(w http.ResponseWriter, r *http.Request) {
	if s.chromeOff(w) {
		return
	}
	if !s.chromeCounted(r, "unpair#"+chromeAddress(r), chromePostWindow, chromePostsPerAddress) {
		s.chromeTooMany(w)
		return
	}
	raw, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var m map[string]json.RawMessage
	var in struct {
		ChromeID string `json:"chrome_id"`
		At       int64  `json:"at"`
		Sig      string `json:"sig"`
	}
	if json.Unmarshal(raw, &m) != nil || len(m) != 3 || json.Unmarshal(raw, &in) != nil || !chromeIDRe.MatchString(in.ChromeID) || in.Sig == "" {
		writeErr(w, http.StatusBadRequest, "send {chrome_id, at, sig}: the Chrome's id, the time in milliseconds, and its key's signature of both")
		return
	}
	at := time.UnixMilli(in.At)
	if d := s.now().Sub(at); d > chromeUnpairSkew || d < -chromeUnpairSkew {
		writeErr(w, http.StatusBadRequest, "at is more than five minutes from planque's clock")
		return
	}
	gone := func() { writeErr(w, http.StatusNotFound, "no such Chrome") }
	var ix chromeIndex
	// The index names the Chrome's account; its record there is the authority (an unpaired one has none).
	if _, err := getDoc(r.Context(), s.KV, pkChromeIndex(in.ChromeID), "-", &ix); err != nil || ix.Tenant == "" {
		gone()
		return
	}
	var c Chrome
	if _, err := getDoc(r.Context(), s.KV, tenantPK(ix.Tenant), skChrome(in.ChromeID), &c); err != nil {
		gone()
		return
	}
	if !stepkey.Verify(c.Key, labelUnpair, map[string]json.RawMessage{"chrome_id": m["chrome_id"], "at": m["at"]}, in.Sig) {
		gone()
		return
	}
	found, err := s.removeChrome(r.Context(), ix.Tenant, in.ChromeID)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !found {
		gone()
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chrome_id": in.ChromeID, "status": "unpaired"})
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

// profileIDShape: the extension's profile id, 16 bytes as 32 lower-case hex.
func profileIDShape(v string) bool {
	if len(v) != 32 {
		return false
	}
	for _, c := range v {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// chromeConnection is one browser connection on a paired Chrome, as the
// Chromes list says it: its grant, the taf's name and the sites it reaches.
type chromeConnection struct {
	Grant string   `json:"grant"`
	Name  string   `json:"name"`
	Sites []string `json:"sites"`
}

// chromesConnections is each Chrome's live connections (asked for or
// approved, not ended or failed) driven by the owner's Chrome, by the
// Chrome's id, in the grants' own order.
func (s *Server) chromesConnections(ctx context.Context, tenant string) (map[string][]chromeConnection, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	on := map[string][]chromeConnection{}
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.ID == "" || skGrant(g.ID) != it.SK || g.Status == StatusDisconnected || g.Status == StatusFailed {
			continue
		}
		if b := g.Request.Browser; b != nil && b.Driver == DriverExtension && b.Chrome != "" {
			name := g.Request.AgentName
			if name == "" {
				name = g.Request.Agent
			}
			on[b.Chrome] = append(on[b.Chrome], chromeConnection{Grant: g.ID, Name: name, Sites: b.Sites})
		}
	}
	return on, nil
}
