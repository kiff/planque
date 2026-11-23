package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// pushService stands in for Chrome's push service: it records each push and
// answers status. Nothing leaves the test.
type pushService struct {
	mu     sync.Mutex
	status int
	got    []*http.Request
	onPush func()
}

func (p *pushService) RoundTrip(r *http.Request) (*http.Response, error) {
	p.mu.Lock()
	p.got = append(p.got, r)
	status, on := p.status, p.onPush
	p.mu.Unlock()
	if on != nil {
		on()
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader("")), Header: http.Header{}, Request: r}, nil
}

func (p *pushService) pushes() []*http.Request {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]*http.Request(nil), p.got...)
}

func newWaker(t *testing.T, svc *pushService) (*ChromeWaker, *ecdsa.PublicKey) {
	t.Helper()
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//nolint:staticcheck // the uncompressed point, as the extension's applicationServerKey
	public := elliptic.Marshal(elliptic.P256(), k.PublicKey.X, k.PublicKey.Y)
	return &ChromeWaker{Signer: stepkey.Local{Key: k}, Public: public, Subject: "https://planque.example", HTTP: &http.Client{Transport: svc}}, &k.PublicKey
}

// A paired Chrome with a subscription and no live connection.
func (r *relayRig) withPush(t *testing.T, endpoint string) {
	t.Helper()
	ctx := context.Background()
	var c Chrome
	it, err := getDoc(ctx, r.e.s.KV, tenantPK(tenant), skChrome(r.chrome), &c)
	if err != nil {
		t.Fatal(err)
	}
	sub := ChromePush{Endpoint: endpoint}
	raw := browserSubscription(t, endpoint)["keys"].(map[string]any)
	sub.Keys.P256DH, sub.Keys.Auth = raw["p256dh"].(string), raw["auth"].(string)
	c.Push = &sub
	if err := putDoc(ctx, r.e.s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1}, c, store.Cond{Ver: &it.Ver}); err != nil {
		t.Fatal(err)
	}
}

// The push is VAPID only: one POST to the endpoint, no body, a short TTL, and
// a JWT signed with planque's key for the push service's origin.
func TestAWakeIsOneEmptyPushWithAVAPIDTokenForThePushService(t *testing.T) {
	svc := &pushService{status: http.StatusCreated}
	w, pub := newWaker(t, svc)
	w.Now = func() time.Time { return time.Unix(1_800_000_000, 0) }
	var sub ChromePush
	sub.Endpoint = fcm
	raw := browserSubscription(t, fcm)["keys"].(map[string]any)
	sub.Keys.P256DH, sub.Keys.Auth = raw["p256dh"].(string), raw["auth"].(string)
	if err := w.Wake(context.Background(), sub); err != nil {
		t.Fatal(err)
	}
	got := svc.pushes()
	if len(got) != 1 || got[0].Method != http.MethodPost || got[0].URL.String() != fcm || got[0].ContentLength > 0 || got[0].Header.Get("TTL") != "30" {
		t.Fatalf("the push: %+v", got)
	}
	auth := got[0].Header.Get("Authorization")
	parts := strings.SplitN(strings.TrimPrefix(auth, "vapid t="), ", k=", 2)
	if len(parts) != 2 || parts[1] != base64.RawURLEncoding.EncodeToString(w.Public) {
		t.Fatalf("Authorization: %q", auth)
	}
	seg := strings.Split(parts[0], ".")
	if len(seg) != 3 {
		t.Fatalf("the JWT: %q", parts[0])
	}
	var claims struct {
		Aud string `json:"aud"`
		Exp int64  `json:"exp"`
		Sub string `json:"sub"`
	}
	body, _ := base64.RawURLEncoding.DecodeString(seg[1])
	if json.Unmarshal(body, &claims) != nil || claims.Aud != "https://fcm.googleapis.com" || claims.Exp != 1_800_003_600 || claims.Sub != "https://planque.example" {
		t.Fatalf("claims: %s", body)
	}
	sig, _ := base64.RawURLEncoding.DecodeString(seg[2])
	sum := sha256.Sum256([]byte(seg[0] + "." + seg[1]))
	if len(sig) != 64 || !ecdsa.Verify(pub, sum[:], new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:])) {
		t.Fatal("the JWT isn't signed with planque's VAPID key")
	}
	// Never anywhere but Chrome's push service, even from a subscription kept before.
	sub.Endpoint = "https://attacker.example/x"
	if err := w.Wake(context.Background(), sub); err == nil || len(svc.pushes()) != 1 {
		t.Fatalf("posted elsewhere: %v %d", err, len(svc.pushes()))
	}
}

// A step for a Chrome that isn't connected wakes it once; when its extension
// connects, the step goes out on that connection.
func TestAStepForAChromeThatIsntConnectedWakesItAndGoesOutOnceItConnects(t *testing.T) {
	r := newRelayRig(t)
	r.withPush(t, fcm)
	svc := &pushService{status: http.StatusCreated}
	r.e.s.ChromeWaker, _ = newWaker(t, svc)
	// The push service delivers; the extension connects a moment later.
	svc.onPush = func() { go func() { time.Sleep(300 * time.Millisecond); r.connect(t, "woken") }() }
	err := r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_1", ChromeID: r.chrome, Op: "read"})
	if err != nil {
		t.Fatalf("the step: %v", err)
	}
	if n := len(svc.pushes()); n != 1 {
		t.Fatalf("pushes: %d", n)
	}
	if m := r.poster.last("woken"); m["type"] != "step" {
		t.Fatalf("the step didn't go out on the woken connection: %v", m)
	}
}

// No connection within the wait (or no subscription, or no waker) is offline,
// as before; a close never wakes.
func TestAWakeThatDoesntConnectInTimeIsOfflineAndACloseNeverWakes(t *testing.T) {
	r := newRelayRig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()
	// No waker configured.
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_0", ChromeID: r.chrome, Op: "read"}); !errors.Is(err, errChromeOffline) {
		t.Fatalf("no waker: %v", err)
	}
	svc := &pushService{status: http.StatusCreated}
	r.e.s.ChromeWaker, _ = newWaker(t, svc)
	// No subscription kept.
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_1", ChromeID: r.chrome, Op: "read"}); !errors.Is(err, errChromeOffline) || len(svc.pushes()) != 0 {
		t.Fatalf("no subscription: %v %d", err, len(svc.pushes()))
	}
	r.withPush(t, fcm)
	if err := r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_2", ChromeID: r.chrome, Op: "close"}); !errors.Is(err, errChromeOffline) || len(svc.pushes()) != 0 {
		t.Fatalf("a close woke it: %v %d", err, len(svc.pushes()))
	}
	// Pushed, never connected: offline once the call's deadline comes, well before the wake's own wait.
	start := time.Now()
	if err := r.e.s.sendChromeStep(ctx, tenant, chromeStep{ID: "st_3", ChromeID: r.chrome, Op: "read"}); !errors.Is(err, errChromeOffline) {
		t.Fatalf("never connected: %v", err)
	}
	if took := time.Since(start); took > 2*time.Second || len(svc.pushes()) != 1 {
		t.Fatalf("the wait ignored the call's deadline: %v, pushes %d", took, len(svc.pushes()))
	}
}

// The push service saying the subscription is gone drops it, so the next step
// doesn't push to it again.
func TestASubscriptionThePushServiceSaysIsGoneIsDropped(t *testing.T) {
	r := newRelayRig(t)
	r.withPush(t, fcm)
	svc := &pushService{status: http.StatusGone}
	r.e.s.ChromeWaker, _ = newWaker(t, svc)
	if err := r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_1", ChromeID: r.chrome, Op: "read"}); !errors.Is(err, errChromeOffline) {
		t.Fatalf("gone: %v", err)
	}
	if p := r.kept(t); p != nil {
		t.Fatalf("a gone subscription was kept: %+v", p)
	}
	_ = r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_2", ChromeID: r.chrome, Op: "read"})
	if n := len(svc.pushes()); n != 1 {
		t.Fatalf("pushed to a gone subscription again: %d", n)
	}
	// Any other refusal keeps it: a busy push service says nothing about the subscription.
	r.withPush(t, fcm)
	svc.status = http.StatusTooManyRequests
	_ = r.e.s.sendChromeStep(context.Background(), tenant, chromeStep{ID: "st_3", ChromeID: r.chrome, Op: "read"})
	if r.kept(t) == nil {
		t.Fatal("a 429 dropped the subscription")
	}
}

// A redirect from the push service is refused, never followed: nothing is
// posted to where it points, and the Chrome isn't taken as woken.
func TestAPushServiceRedirectIsRefusedAndNeverFollowed(t *testing.T) {
	var mu sync.Mutex
	var hosts []string
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		hosts = append(hosts, r.URL.Host)
		mu.Unlock()
		if r.URL.Host == "fcm.googleapis.com" {
			h := http.Header{}
			h.Set("Location", "https://evil.example/x")
			return &http.Response{StatusCode: http.StatusTemporaryRedirect, Body: io.NopCloser(strings.NewReader("")), Header: h, Request: r}, nil
		}
		return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader("")), Header: http.Header{}, Request: r}, nil
	})
	w, _ := newWaker(t, &pushService{})
	w.HTTP = &http.Client{Transport: rt}
	var sub ChromePush
	sub.Endpoint = fcm
	raw := browserSubscription(t, fcm)["keys"].(map[string]any)
	sub.Keys.P256DH, sub.Keys.Auth = raw["p256dh"].(string), raw["auth"].(string)
	if err := w.Wake(context.Background(), sub); err == nil || !strings.Contains(err.Error(), "307") {
		t.Fatalf("a 307 was taken as woken: %v", err)
	}
	if len(hosts) != 1 || hosts[0] != "fcm.googleapis.com" {
		t.Fatalf("followed the redirect: %v", hosts)
	}
	// The waker's own client is never changed by it.
	if w.HTTP.CheckRedirect != nil {
		t.Fatal("the waker's client was changed")
	}
	// No contact, no push.
	w.Subject = ""
	if err := w.Wake(context.Background(), sub); err == nil || len(hosts) != 1 {
		t.Fatalf("pushed with no subject: %v %v", err, hosts)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
