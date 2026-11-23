package planque

import (
	"context"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/stepkey"
)

// A browser's PushSubscription.toJSON(), as the extension sends it: expirationTime included.
func browserSubscription(t *testing.T, endpoint string) map[string]any {
	t.Helper()
	k, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	auth := make([]byte, 16)
	_, _ = rand.Read(auth)
	return map[string]any{"endpoint": endpoint, "expirationTime": nil, "keys": map[string]any{
		"p256dh": base64.RawURLEncoding.EncodeToString(k.PublicKey().Bytes()), "auth": base64.RawURLEncoding.EncodeToString(auth)}}
}

func (r *relayRig) push(t *testing.T, conn string, key stepkey.Local, sub map[string]any) map[string]any {
	t.Helper()
	sig, err := stepkey.SignBody(context.Background(), key, labelPush, sub)
	if err != nil {
		t.Fatal(err)
	}
	r.event("$default", conn, map[string]any{"type": "push", "subscription": sub, "sig": sig})
	return r.poster.last(conn)
}

func (r *relayRig) kept(t *testing.T) *ChromePush {
	t.Helper()
	var c Chrome
	if _, err := getDoc(context.Background(), r.e.s.KV, tenantPK(tenant), skChrome(r.chrome), &c); err != nil {
		t.Fatal(err)
	}
	return c.Push
}

const fcm = "https://fcm.googleapis.com/fcm/send/dXk3:APA91b-example"

// A connected Chrome's signed subscription is kept for its wake; a new one
// replaces it, and the Chromes list never shows it.
func TestAConnectedChromesSignedPushSubscriptionIsKeptAndReplaced(t *testing.T) {
	r := newRelayRig(t)
	r.connect(t, "c1")
	first := browserSubscription(t, fcm)
	if m := r.push(t, "c1", r.ext, first); m["type"] != "push_kept" {
		t.Fatalf("a signed subscription: %v", m)
	}
	if p := r.kept(t); p == nil || p.Endpoint != fcm || p.Keys.P256DH != first["keys"].(map[string]any)["p256dh"] {
		t.Fatalf("kept: %+v", p)
	}
	second := browserSubscription(t, fcm+"-rotated")
	if m := r.push(t, "c1", r.ext, second); m["type"] != "push_kept" || r.kept(t).Endpoint != fcm+"-rotated" {
		t.Fatalf("a rotated subscription: %v %+v", m, r.kept(t))
	}
	// The same again is kept without a write.
	before, _, _ := r.e.s.KV.Get(context.Background(), tenantPK(tenant), skChrome(r.chrome))
	if m := r.push(t, "c1", r.ext, second); m["type"] != "push_kept" {
		t.Fatalf("the same again: %v", m)
	}
	if after, _, _ := r.e.s.KV.Get(context.Background(), tenantPK(tenant), skChrome(r.chrome)); after.Ver != before.Ver {
		t.Fatalf("the same subscription wrote again: %d → %d", before.Ver, after.Ver)
	}
	// Never in the Chromes list.
	st, out := r.e.do("GET", "/v1/accounts/"+tenant+"/chromes", nil, r.e.phone, "")
	if b, _ := json.Marshal(out); st != http.StatusOK || strings.Contains(string(b), "fcm.googleapis") || strings.Contains(string(b), `"push"`) {
		t.Fatalf("the list shows the subscription: %d %s", st, b)
	}
}

// What isn't this Chrome's own, or names anywhere but Chrome's push service,
// is refused with why, and the connection stays: nothing is kept.
func TestAPushSubscriptionIsRefusedUnlessSignedAndOnChromesPushService(t *testing.T) {
	r := newRelayRig(t)
	r.connect(t, "c1")
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	for name, c := range map[string]struct {
		key  stepkey.Local
		sub  func() map[string]any
		code string
	}{
		"another key's signature": {stepkey.Local{Key: other}, func() map[string]any { return browserSubscription(t, fcm) }, codePushUnsigned},
		"http":                    {r.ext, func() map[string]any { return browserSubscription(t, "http://fcm.googleapis.com/fcm/send/x") }, codePushEndpoint},
		"another host":            {r.ext, func() map[string]any { return browserSubscription(t, "https://attacker.example/fcm/send/x") }, codePushEndpoint},
		"a host that ends alike":  {r.ext, func() map[string]any { return browserSubscription(t, "https://evilfcm.googleapis.com.example/x") }, codePushEndpoint},
		"a port":                  {r.ext, func() map[string]any { return browserSubscription(t, "https://fcm.googleapis.com:8443/fcm/send/x") }, codePushEndpoint},
		"user info":               {r.ext, func() map[string]any { return browserSubscription(t, "https://a@fcm.googleapis.com/fcm/send/x") }, codePushEndpoint},
		"no path":                 {r.ext, func() map[string]any { return browserSubscription(t, "https://fcm.googleapis.com/") }, codePushEndpoint},
		"too long":                {r.ext, func() map[string]any { return browserSubscription(t, fcm+strings.Repeat("a", 1024)) }, codePushEndpoint},
		"a short p256dh": {r.ext, func() map[string]any {
			s := browserSubscription(t, fcm)
			s["keys"].(map[string]any)["p256dh"] = base64.RawURLEncoding.EncodeToString(make([]byte, 33))
			return s
		}, codePushKeys},
		"a long auth": {r.ext, func() map[string]any {
			s := browserSubscription(t, fcm)
			s["keys"].(map[string]any)["auth"] = base64.RawURLEncoding.EncodeToString(make([]byte, 32))
			return s
		}, codePushKeys},
	} {
		if m := r.push(t, "c1", c.key, c.sub()); m["type"] != "push_refused" || m["code"] != c.code {
			t.Errorf("%s: %v, want push_refused %s", name, m, c.code)
		}
		if r.poster.dropped["c1"] {
			t.Fatalf("%s: the connection was closed", name)
		}
	}
	if p := r.kept(t); p != nil {
		t.Fatalf("a refused subscription was kept: %+v", p)
	}
	// Signed over what was sent: a field changed after signing doesn't hold.
	sub := browserSubscription(t, fcm)
	sig, _ := stepkey.SignBody(context.Background(), r.ext, labelPush, sub)
	sub["expirationTime"] = 1
	r.event("$default", "c1", map[string]any{"type": "push", "subscription": sub, "sig": sig})
	if m := r.poster.last("c1"); m["code"] != codePushUnsigned {
		t.Fatalf("edited after signing: %v", m)
	}
}

// Before the welcome, a push message is like any other unexpected one: it
// ends the connection, and nothing is kept.
func TestAPushBeforeTheWelcomeEndsTheConnection(t *testing.T) {
	r := newRelayRig(t)
	if st := r.event("$connect", "c1", r.chrome); st != http.StatusOK {
		t.Fatalf("$connect: %d", st)
	}
	r.push(t, "c1", r.ext, browserSubscription(t, fcm))
	closedWith(t, r, "c1", closeUnknownChrome)
	if p := r.kept(t); p != nil {
		t.Fatalf("kept before the welcome: %+v", p)
	}
}

// A Chrome unpaired while its subscription was being kept stays unpaired: the
// write is on the version read, and the answer says it wasn't kept.
func TestAPushSubscriptionNeverRecreatesAChromeUnpairedMeanwhile(t *testing.T) {
	r := newRelayRig(t)
	r.connect(t, "c1")
	beforePushKept = func() {
		_, _, _ = r.e.s.KV.Delete(context.Background(), tenantPK(tenant), skChrome(r.chrome))
	}
	defer func() { beforePushKept = nil }()
	if m := r.push(t, "c1", r.ext, browserSubscription(t, fcm)); m["type"] != "push_refused" || m["code"] != codePushNotKept {
		t.Fatalf("kept over an unpairing: %v", m)
	}
	if _, ok, _ := r.e.s.KV.Get(context.Background(), tenantPK(tenant), skChrome(r.chrome)); ok {
		t.Fatal("the unpaired Chrome came back")
	}
}
