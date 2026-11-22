package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/lambda"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// fakePoster is API Gateway's management endpoint: what was posted to each
// connection, and which were dropped.
type fakePoster struct {
	mu      sync.Mutex
	posts   map[string][]map[string]any
	dropped map[string]bool
}

func (f *fakePoster) Post(_ context.Context, conn string, msg []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var m map[string]any
	_ = json.Unmarshal(msg, &m)
	if f.posts == nil {
		f.posts = map[string][]map[string]any{}
	}
	f.posts[conn] = append(f.posts[conn], m)
	return nil
}

func (f *fakePoster) Drop(_ context.Context, conn string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.dropped == nil {
		f.dropped = map[string]bool{}
	}
	f.dropped[conn] = true
	return nil
}

func (f *fakePoster) last(conn string) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.posts[conn]) == 0 {
		return nil
	}
	return f.posts[conn][len(f.posts[conn])-1]
}

// relayRig is a planque with the relay on and one paired Chrome, whose key
// the test holds as the extension does.
type relayRig struct {
	e       *env
	poster  *fakePoster
	planque stepkey.Local
	ext     stepkey.Local
	chrome  string
}

func newRelayRig(t *testing.T) *relayRig {
	t.Helper()
	return newRelayRigOn(t, chromeEnv(t))
}

// newRelayRigOn is newRelayRig on e, whose Chrome offer the test has set.
func newRelayRigOn(t *testing.T, e *env) *relayRig {
	t.Helper()
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ek, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	r := &relayRig{e: e, poster: &fakePoster{}, planque: stepkey.Local{Key: pk}, ext: stepkey.Local{Key: ek}, chrome: "chr_workmac0001"}
	e.s.ChromeSigner, e.s.ChromePoster = r.planque, r.poster
	der, _ := x509.MarshalPKIXPublicKey(&ek.PublicKey)
	c := Chrome{ID: r.chrome, Tenant: tenant, Name: "Work Mac", Key: base64.RawURLEncoding.EncodeToString(der)}
	ctx := context.Background()
	if err := putDoc(ctx, e.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(c.ID)}, c, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	if err := putDoc(ctx, e.s.KV, store.Item{PK: pkChromeIndex(c.ID), SK: "-"}, chromeIndex{Tenant: tenant}, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	return r
}

func (r *relayRig) event(route, conn string, msg any) int {
	var body []byte
	if msg != nil {
		body, _ = json.Marshal(msg)
	}
	ev := lambda.WSEvent{Route: route, Connection: conn, Body: body}
	if route == "$connect" {
		ev.Query = map[string]string{"chrome": msg.(string)}
		ev.Body = nil
	}
	return r.e.s.WebSocket(context.Background(), ev)
}

// hello is the extension's proof for nonceP, signed with key.
// Each edit changes the body before it's signed.
func (r *relayRig) hello(t *testing.T, key stepkey.Local, nonceP, manifest string, edits ...func(map[string]any)) map[string]any {
	t.Helper()
	body := map[string]any{"chrome_id": r.chrome, "nonce_p": nonceP, "nonce_e": "ne-123", "version": "1.0.0", "manifest_hash": manifest}
	for _, edit := range edits {
		edit(body)
	}
	sig, err := stepkey.SignBody(context.Background(), key, labelHello, body)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]any{"type": "hello", "sig": sig}
	for k, v := range body {
		m[k] = v
	}
	return m
}

// connect runs the whole connect: open, challenge, hello, welcome.
func (r *relayRig) connect(t *testing.T, conn string) {
	t.Helper()
	if st := r.event("$connect", conn, r.chrome); st != http.StatusOK {
		t.Fatalf("$connect: %d", st)
	}
	r.event("$default", conn, map[string]any{"type": "open"})
	ch := r.poster.last(conn)
	if ch["type"] != "challenge" || ch["nonce_p"] == "" {
		t.Fatalf("challenge: %v", ch)
	}
	r.event("$default", conn, r.hello(t, r.ext, ch["nonce_p"].(string), releasedManifest))
	w := r.poster.last(conn)
	if w["type"] != "welcome" || !stepkey.Verify(r.planque.PublicKey(), labelWelcome, map[string]string{"chrome_id": r.chrome, "nonce_p": ch["nonce_p"].(string), "nonce_e": "ne-123"}, w["sig"].(string)) {
		t.Fatalf("welcome: %v", w)
	}
	// planque's clock, for an extension whose computer's clock is off.
	if at, ok := w["now"].(float64); !ok || time.Since(time.UnixMilli(int64(at))).Abs() > 5*time.Second {
		t.Fatalf("the welcome's now: %v", w["now"])
	}
}

func closedWith(t *testing.T, r *relayRig, conn string, code float64) {
	t.Helper()
	m := r.poster.last(conn)
	if m["type"] != "close" || m["code"] != code || !r.poster.dropped[conn] {
		t.Fatalf("%s: want close %v, got %v (dropped %v)", conn, code, m, r.poster.dropped[conn])
	}
}

// The extension speaks first; planque challenges, the extension proves its
// paired key, planque proves its step key. Then steps go out signed, and
// only results the Chrome signed are kept, the first one for a step.
func TestTheRelayProvesBothEndsThenSendsSignedStepsAndKeepsSignedResults(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	r.connect(t, "c1")
	r.event("$default", "c1", map[string]any{"type": "ping"})
	if r.poster.last("c1")["type"] != "pong" {
		t.Fatalf("ping: %v", r.poster.last("c1"))
	}
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_1", ChromeID: r.chrome, Run: "run_1", Grant: "grt_1", Sites: []string{"shop.example"}, Op: "read", URL: "https://shop.example/"}); err != nil {
		t.Fatal(err)
	}
	m := r.poster.last("c1")
	raw, _ := json.Marshal(m["step"])
	var step map[string]any
	_ = json.Unmarshal(raw, &step)
	if m["type"] != "step" || step["nonce"] == "" || !stepkey.Verify(r.planque.PublicKey(), labelStep, step, m["sig"].(string)) {
		t.Fatalf("step: %v", m)
	}
	if exp := time.UnixMilli(int64(step["exp"].(float64))); exp.After(time.Now().Add(31*time.Second)) || exp.Before(time.Now()) {
		t.Fatalf("a step lapses within 30 s of being sent, not at %v", exp)
	}
	nonce := step["nonce"].(string)
	// A result counts only when the Chrome signed it, over the step's own
	// nonce and its own id, for a step planque sent it: the first such one.
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	result := func(key stepkey.Local, id, nonce, chromeID, page string) map[string]any {
		body := map[string]any{"id": id, "nonce": nonce, "chrome_id": chromeID, "ok": true, "page": map[string]any{"title": page}, "error": nil}
		sig, _ := stepkey.SignBody(ctx, key, labelResult, body)
		msg := map[string]any{"type": "result", "sig": sig}
		for k, v := range body {
			msg[k] = v
		}
		return msg
	}
	for name, m := range map[string]map[string]any{
		"another key":          result(stepkey.Local{Key: other}, "st_1", nonce, r.chrome, "forged"),
		"another step's nonce": result(r.ext, "st_1", "not-the-nonce", r.chrome, "forged"),
		"another Chrome's id":  result(r.ext, "st_1", nonce, "chr_someoneelse01", "forged"),
	} {
		r.event("$default", "c1", m)
		if _, found, _ := r.e.s.chromeResult(ctx, r.chrome, "st_1"); found {
			t.Fatalf("%s: a result was kept", name)
		}
	}
	// A step not sent yet: its result is never filed, so it can't stand in
	// for the real one later.
	r.event("$default", "c1", result(r.ext, "st_9", nonce, r.chrome, "early"))
	if _, found, _ := r.e.s.chromeResult(ctx, r.chrome, "st_9"); found {
		t.Fatal("a result for a step never sent was kept")
	}
	r.event("$default", "c1", result(r.ext, "st_1", nonce, r.chrome, "Shop"))
	r.event("$default", "c1", result(r.ext, "st_1", nonce, r.chrome, "Second"))
	got, found, err := r.e.s.chromeResult(ctx, r.chrome, "st_1")
	if err != nil || !found || string(got["page"]) != `{"title":"Shop"}` {
		t.Fatalf("result: %v %v %s", found, err, got["page"])
	}
	// A step id is sent once.
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_1", ChromeID: r.chrome}); err == nil {
		t.Fatal("a step id was sent twice")
	}
	// A result after its step lapsed (exp plus the grace) isn't filed.
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_3", ChromeID: r.chrome}); err != nil {
		t.Fatal(err)
	}
	raw, _ = json.Marshal(r.poster.last("c1")["step"])
	var late map[string]any
	_ = json.Unmarshal(raw, &late)
	r.e.s.Now = func() time.Time { return time.Now().Add(chromeStepTTL + chromeResultGrace + time.Minute) }
	r.event("$default", "c1", result(r.ext, "st_3", late["nonce"].(string), r.chrome, "late"))
	r.e.s.Now = nil
	if _, found, _ := r.e.s.chromeResult(ctx, r.chrome, "st_3"); found {
		t.Fatal("a result after its step lapsed was kept")
	}
	// Disconnected: a step finds no Chrome.
	r.event("$disconnect", "c1", nil)
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_2", ChromeID: r.chrome}); err != errChromeOffline {
		t.Fatalf("after disconnect: %v", err)
	}
}

// Each refusal closes with its code, as a message, then drops: an unknown
// Chrome (4401), an unreleased build (4403), a hello for another challenge
// or by another key (4401), anything before the proofs (4401), and an
// unpaired Chrome (4410), whether it connects after or is live at the time.
func TestTheRelayClosesWhatItCantTrustWithItsCode(t *testing.T) {
	t.Run("unknown Chrome", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", "chr_nobodyknows1")
		r.event("$default", "c1", map[string]any{"type": "open"})
		closedWith(t, r, "c1", 4401)
	})
	t.Run("unreleased build", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		nonce := r.poster.last("c1")["nonce_p"].(string)
		r.event("$default", "c1", r.hello(t, r.ext, nonce, "bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11bb11"))
		closedWith(t, r, "c1", 4403)
	})
	t.Run("another challenge", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		r.event("$default", "c1", r.hello(t, r.ext, "not-the-challenge", releasedManifest))
		closedWith(t, r, "c1", 4401)
	})
	t.Run("no nonce of its own", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		r.event("$default", "c1", r.hello(t, r.ext, r.poster.last("c1")["nonce_p"].(string), releasedManifest, func(b map[string]any) { b["nonce_e"] = "" }))
		closedWith(t, r, "c1", 4401)
	})
	t.Run("another Chrome's id, signed with this one's key", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		r.event("$default", "c1", r.hello(t, r.ext, r.poster.last("c1")["nonce_p"].(string), releasedManifest, func(b map[string]any) { b["chrome_id"] = "chr_someoneelse01" }))
		closedWith(t, r, "c1", 4401)
	})
	t.Run("another key", func(t *testing.T) {
		r := newRelayRig(t)
		other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		r.event("$default", "c1", r.hello(t, stepkey.Local{Key: other}, r.poster.last("c1")["nonce_p"].(string), releasedManifest))
		closedWith(t, r, "c1", 4401)
	})
	t.Run("a second open", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		r.event("$default", "c1", map[string]any{"type": "open"})
		closedWith(t, r, "c1", 4401)
	})
	t.Run("a step's result before the proofs", func(t *testing.T) {
		r := newRelayRig(t)
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "result", "id": "x"})
		closedWith(t, r, "c1", 4401)
	})
	t.Run("unpaired, connecting after", func(t *testing.T) {
		r := newRelayRig(t)
		if st, _ := r.e.do("DELETE", "/v1/accounts/"+tenant+"/chromes/"+r.chrome, nil, r.e.phone, ""); st != 200 {
			t.Fatalf("unpair: %d", st)
		}
		r.event("$connect", "c1", r.chrome)
		r.event("$default", "c1", map[string]any{"type": "open"})
		closedWith(t, r, "c1", 4410)
	})
	t.Run("unpaired while live", func(t *testing.T) {
		r := newRelayRig(t)
		r.connect(t, "c1")
		if st, _ := r.e.do("DELETE", "/v1/accounts/"+tenant+"/chromes/"+r.chrome, nil, r.e.phone, ""); st != 200 {
			t.Fatalf("unpair: %d", st)
		}
		closedWith(t, r, "c1", 4410)
		if err := r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_1", ChromeID: r.chrome}); err != errChromeOffline {
			t.Fatalf("a step after unpairing: %v", err)
		}
	})
}

// A paired Chrome whose account isn't offered it is refused as an unknown
// one, and no step is sent to it.
func TestTheRelayRefusesAChromeWhoseAccountIsntOffered(t *testing.T) {
	r := newRelayRig(t)
	r.e.s.ChromeOffer = ChromeOffer{Tenants: map[string]bool{"tnt_other": true}}
	if st := r.event("$connect", "c1", r.chrome); st != http.StatusOK {
		t.Fatalf("connect: %d", st)
	}
	r.event("$default", "c1", map[string]any{"type": "open"})
	closedWith(t, r, "c1", 4401)
}

// A Chrome connected while its account was offered gets no step once the
// account is taken off, though its connection is still open.
func TestNoStepGoesToAConnectedChromeOnceItsAccountIsntOffered(t *testing.T) {
	r := newRelayRig(t)
	r.e.s.ChromeOffer = ChromeOffer{Tenants: map[string]bool{tenant: true}}
	r.connect(t, "c1")
	sent := len(r.poster.posts["c1"])
	r.e.s.ChromeOffer = ChromeOffer{Tenants: map[string]bool{"tnt_other": true}}
	err := r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_1", ChromeID: r.chrome, Op: "read", URL: "https://shop.example/"})
	if !errors.Is(err, errChromeOffline) || len(r.poster.posts["c1"]) != sent {
		t.Fatalf("a step for an account no longer offered: %v, %d posts after %d", err, len(r.poster.posts["c1"]), sent)
	}
}

// With no step key or no WebSocket API, or the Chrome not offered, every
// connect is refused.
func TestTheRelayIsOffWithoutItsKeyAndAPI(t *testing.T) {
	r := newRelayRig(t)
	r.e.s.ChromeSigner = nil
	if st := r.event("$connect", "c1", r.chrome); st != http.StatusForbidden {
		t.Fatalf("no step key: %d", st)
	}
	r.e.s.ChromeSigner, r.e.s.ChromeBuilds = r.planque, nil
	if st := r.event("$connect", "c1", r.chrome); st != http.StatusForbidden {
		t.Fatalf("not offered: %d", st)
	}
}
