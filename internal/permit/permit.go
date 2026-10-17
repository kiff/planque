// Package permit verifies KIFF's execution permits (kiff-cloud RFC 046):
// a JWT signed with Ed25519 that authorizes at most one execution of one
// exact call. planque checks one before every call it makes to a vendor.
package permit

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/jcs"
)

// Header carries the permit on the gateway's tools/call request.
const Header = "KIFF-Permit"

const (
	typ        = "kiff-permit+jwt"
	alg        = "Ed25519"
	detailType = "https://kiff.dev/permits/tool-call"
	// Tolerance is the clock skew allowed on iat and exp (RFC 046 D5).
	Tolerance = 30 * time.Second
	// MaxLifetime is the most exp-iat any verifier accepts.
	MaxLifetime = 300 * time.Second
)

// Permit is what a verified permit authorizes.
type Permit struct {
	ID       string // jti
	Subject  string // sub: on whose authority (the agent)
	Tenant   string
	Tool     string
	ArgsHash string
	Op       string
	Card     string
	Decision string
	IssuedAt time.Time
	Expires  time.Time
}

// Live reports whether the permit may start a new execution at now:
// unexpired, with the clock tolerance (D5 step 6).
func (p Permit) Live(now time.Time) bool { return now.Before(p.Expires.Add(Tolerance)) }

// Want is the call a permit must name exactly.
type Want struct {
	Audience string
	Tenant   string
	Tool     string
	Args     json.RawMessage // the arguments as received, before any change
}

// Keys returns a trusted public key by kid.
type Keys interface {
	Key(ctx context.Context, kid string) (ed25519.PublicKey, error)
}

// Verifier checks a permit's authenticity (RFC 046 D5 steps 1–4).
type Verifier struct {
	Issuer string
	Keys   Keys
	// Lifetime is the longest exp-iat accepted: 60 seconds by default,
	// never above MaxLifetime.
	Lifetime time.Duration
	Now      func() time.Time
}

// Refusals. None is a reason to send anything.
var (
	ErrMissing   = errors.New("no execution permit: planque calls the vendor only with KIFF's permit for that exact call")
	ErrMalformed = errors.New("the permit is not a well-formed KIFF permit")
	ErrSignature = errors.New("the permit's signature does not verify with a key KIFF publishes")
	ErrIssuer    = errors.New("the permit is not from the configured KIFF issuer")
	ErrAudience  = errors.New("the permit is for another audience")
	ErrTenant    = errors.New("the permit is for another KIFF account")
	ErrTool      = errors.New("the permit is for another tool")
	ErrArgs      = errors.New("the permit is for other arguments")
	ErrLifetime  = errors.New("the permit's lifetime is longer than this relay accepts, or it was issued in the future")
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

type claims struct {
	Iss     string          `json:"iss"`
	Aud     json.RawMessage `json:"aud"`
	Sub     string          `json:"sub"`
	Iat     *int64          `json:"iat"`
	Exp     *int64          `json:"exp"`
	Jti     string          `json:"jti"`
	Details []struct {
		Type       string `json:"type"`
		Tenant     string `json:"tenant"`
		Tool       string `json:"tool"`
		ArgsSHA256 string `json:"args_sha256"`
		Op         string `json:"op"`
		Card       string `json:"card"`
		Decision   string `json:"decision"`
	} `json:"authorization_details"`
}

// Verify checks that token is a permit KIFF issued for exactly want. It
// does not check expiry: an expired authentic permit may still learn how
// its operation ended (D5 step 5). Use Permit.Live before sending.
func (v *Verifier) Verify(ctx context.Context, token string, want Want) (Permit, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return Permit{}, ErrMissing
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Permit{}, ErrMalformed
	}
	var h header
	if !decodePart(parts[0], &h) || h.Typ != typ || h.Alg != alg || h.Kid == "" {
		return Permit{}, ErrMalformed
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Permit{}, ErrMalformed
	}
	key, err := v.Keys.Key(ctx, h.Kid)
	if err != nil {
		return Permit{}, fmt.Errorf("%w: %v", ErrSignature, err)
	}
	if !ed25519.Verify(key, []byte(parts[0]+"."+parts[1]), sig) {
		return Permit{}, ErrSignature
	}
	var c claims
	if !decodePart(parts[1], &c) || c.Iat == nil || c.Exp == nil || len(c.Details) != 1 {
		return Permit{}, ErrMalformed
	}
	d := c.Details[0]
	if d.Type != detailType || d.Op == "" || d.ArgsSHA256 == "" {
		return Permit{}, ErrMalformed
	}
	if c.Iss != v.Issuer {
		return Permit{}, ErrIssuer
	}
	if !audienceIs(c.Aud, want.Audience) {
		return Permit{}, ErrAudience
	}
	if d.Tenant != want.Tenant {
		return Permit{}, ErrTenant
	}
	if d.Tool != want.Tool {
		return Permit{}, ErrTool
	}
	sum, err := jcs.SHA256(want.Args)
	if err != nil || sum != d.ArgsSHA256 {
		return Permit{}, ErrArgs
	}
	iat, exp := time.Unix(*c.Iat, 0), time.Unix(*c.Exp, 0)
	life := v.Lifetime
	if life <= 0 {
		life = 60 * time.Second
	}
	life = min(life, MaxLifetime)
	if exp.Sub(iat) > life || exp.Before(iat) || iat.After(v.now().Add(Tolerance)) {
		return Permit{}, ErrLifetime
	}
	return Permit{ID: c.Jti, Subject: c.Sub, Tenant: d.Tenant, Tool: d.Tool, ArgsHash: d.ArgsSHA256,
		Op: d.Op, Card: d.Card, Decision: d.Decision, IssuedAt: iat, Expires: exp}, nil
}

func (v *Verifier) now() time.Time {
	if v.Now != nil {
		return v.Now()
	}
	return time.Now()
}

func decodePart(s string, v any) bool {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return false
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	return dec.Decode(v) == nil
}

// audienceIs accepts aud as one string or an array holding exactly want.
func audienceIs(raw json.RawMessage, want string) bool {
	var one string
	if json.Unmarshal(raw, &one) == nil {
		return one == want && want != ""
	}
	var many []string
	if json.Unmarshal(raw, &many) == nil {
		return len(many) == 1 && many[0] == want && want != ""
	}
	return false
}

// JWKS trusts the keys KIFF publishes (RFC 046 D4): refetched at least
// every 5 minutes, and on an unknown kid at most once a minute. A kid
// missing from the latest fetch is not trusted. If refetching fails, the
// last good set is trusted for at most an hour from its fetch.
type JWKS struct {
	URL  string
	HTTP *http.Client
	Now  func() time.Time

	mu          sync.Mutex
	keys        map[string]ed25519.PublicKey
	fetchedAt   time.Time // last successful fetch
	lastUnknown time.Time // last refetch for an unknown kid
}

const (
	jwksRefresh  = 5 * time.Minute
	jwksUnknown  = time.Minute
	jwksStaleMax = time.Hour
)

// Key returns the trusted key named kid.
func (j *JWKS) Key(ctx context.Context, kid string) (ed25519.PublicKey, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	now := j.now()
	if j.keys == nil || now.Sub(j.fetchedAt) >= jwksRefresh {
		j.refetch(ctx, now)
	}
	if _, ok := j.keys[kid]; !ok && now.Sub(j.lastUnknown) >= jwksUnknown {
		j.lastUnknown = now
		j.refetch(ctx, now)
	}
	if j.keys == nil || now.Sub(j.fetchedAt) >= jwksStaleMax {
		return nil, errors.New("KIFF's permit keys could not be fetched recently")
	}
	k, ok := j.keys[kid]
	if !ok {
		return nil, errors.New("unknown key " + kid)
	}
	return k, nil
}

func (j *JWKS) now() time.Time {
	if j.Now != nil {
		return j.Now()
	}
	return time.Now()
}

// refetch replaces the key set on success and keeps it on failure.
func (j *JWKS) refetch(ctx context.Context, now time.Time) {
	keys, err := j.fetch(ctx)
	if err != nil {
		return
	}
	j.keys, j.fetchedAt = keys, now
}

func (j *JWKS) fetch(ctx context.Context) (map[string]ed25519.PublicKey, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, j.URL, nil)
	if err != nil {
		return nil, err
	}
	hc := j.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks: status %d", resp.StatusCode)
	}
	var set struct {
		Keys []struct {
			Kty, Crv, X, Kid, Alg string
		} `json:"keys"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&set); err != nil {
		return nil, err
	}
	keys := map[string]ed25519.PublicKey{}
	for _, k := range set.Keys {
		if k.Kty != "OKP" || k.Crv != "Ed25519" || k.Kid == "" || (k.Alg != "" && k.Alg != alg) {
			continue
		}
		x, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil || len(x) != ed25519.PublicKeySize {
			continue
		}
		keys[k.Kid] = ed25519.PublicKey(x)
	}
	return keys, nil
}
