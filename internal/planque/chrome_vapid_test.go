package planque

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// planque makes its VAPID key once, keeps the private key sealed and the public
// key beside it, and every later need reads the same key.
func TestPlanquesVAPIDKeyIsMadeOnceAndKeptSealed(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	k1, pub1, err := r.e.s.vapidKey(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var d vapidDoc
	if _, err := getDoc(ctx, r.e.s.KV, pkVAPID, "-", &d); err != nil {
		t.Fatal(err)
	}
	if d.Public != base64.RawURLEncoding.EncodeToString(pub1) || len(pub1) != 65 || pub1[0] != 4 {
		t.Fatalf("the public key kept: %q", d.Public)
	}
	// Sealed: the private scalar isn't in the item.
	if bytes.Contains(d.Sealed, k1.D.FillBytes(make([]byte, 32))) {
		t.Fatal("the private key is kept in the clear")
	}
	k2, pub2, err := r.e.s.vapidKey(ctx)
	if err != nil || !k2.Equal(k1) || !bytes.Equal(pub2, pub1) {
		t.Fatalf("a second read gave another key: %v", err)
	}
}

// Instances making the key at once keep one: the first write wins and the
// others read it.
func TestInstancesMakingTheVAPIDKeyAtOnceKeepOne(t *testing.T) {
	r := newRelayRig(t)
	var wg sync.WaitGroup
	pubs := make([][]byte, 8)
	for i := range pubs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, pub, err := r.e.s.vapidKey(context.Background())
			if err != nil {
				t.Error(err)
			}
			pubs[i] = pub
		}(i)
	}
	wg.Wait()
	for _, p := range pubs[1:] {
		if !bytes.Equal(p, pubs[0]) {
			t.Fatal("two instances kept two keys")
		}
	}
}

// A public key kept beside a private key that isn't its own is refused: the
// extension would have pinned a key planque can't sign for.
func TestAVAPIDPublicKeyThatIsntThePrivateKeysIsRefused(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	if _, _, err := r.e.s.vapidKey(ctx); err != nil {
		t.Fatal(err)
	}
	var d vapidDoc
	it, _ := getDoc(ctx, r.e.s.KV, pkVAPID, "-", &d)
	other := browserSubscription(t, fcm)["keys"].(map[string]any)["p256dh"].(string)
	d.Public = other
	if err := putDoc(ctx, r.e.s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1}, d, store.Cond{Ver: &it.Ver}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r.e.s.vapidKey(ctx); err == nil {
		t.Fatal("a mismatched public key was taken")
	}
}

// The waker planque builds signs with the key it keeps, and the route gives
// the extension's build that key's public half; off without a subject.
func TestTheWakerAndThePushKeyRouteUsePlanquesOwnKey(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	if w := r.e.s.chromeWaker(ctx); w != nil {
		t.Fatal("a waker with no VAPID subject")
	}
	st, _ := r.e.do("GET", "/v1/extension/push-key", nil, nil, "")
	if st != http.StatusNotFound {
		t.Fatalf("the push key with no subject: %d", st)
	}
	r.e.s.VAPIDSubject = "https://planque.example"
	w := r.e.s.chromeWaker(ctx)
	if w == nil || w.Subject != "https://planque.example" {
		t.Fatalf("the waker: %+v", w)
	}
	k, pub, _ := r.e.s.vapidKey(ctx)
	if !bytes.Equal(w.Public, pub) {
		t.Fatal("the waker's public key isn't the kept one")
	}
	sig, err := w.Signer.Sign(ctx, []byte("header.payload"))
	sum := sha256.Sum256([]byte("header.payload"))
	if err != nil || len(sig) != 64 || !ecdsa.Verify(&k.PublicKey, sum[:], new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:])) {
		t.Fatal("the waker doesn't sign with the kept key")
	}
	st, out := r.e.do("GET", "/v1/extension/push-key", nil, nil, "")
	b, _ := json.Marshal(out)
	if st != http.StatusOK || out["key"] != base64.RawURLEncoding.EncodeToString(pub) {
		t.Fatalf("the push key: %d %s", st, b)
	}
	if w2 := r.e.s.chromeWaker(ctx); w2 != w {
		t.Fatal("the waker was built again")
	}
}
