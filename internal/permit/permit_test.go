package permit

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestJWKSTrust checks RFC 046 D4: an unknown kid refetches at most once
// a minute, a key gone from the latest set is not trusted, and the last
// good set is trusted for at most an hour when refetching fails.
func TestJWKSTrust(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	var mu sync.Mutex
	kids, fetches, down := []string{"k1"}, 0, false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		fetches++
		if down {
			http.Error(w, "down", 503)
			return
		}
		var keys []map[string]string
		for _, k := range kids {
			keys = append(keys, map[string]string{"kty": "OKP", "crv": "Ed25519", "kid": k, "x": base64.RawURLEncoding.EncodeToString(pub)})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": keys})
	}))
	defer srv.Close()
	now := time.Unix(1_800_000_000, 0)
	j := &JWKS{URL: srv.URL, Now: func() time.Time { return now }}
	ctx := context.Background()
	set := func(f func()) { mu.Lock(); f(); mu.Unlock() }
	count := func() int { mu.Lock(); defer mu.Unlock(); return fetches }

	if _, err := j.Key(ctx, "k1"); err != nil || count() != 1 {
		t.Fatalf("first fetch: %v, fetches %d", err, count())
	}
	// Unknown kids: one refetch, then none within the minute.
	for range 5 {
		if _, err := j.Key(ctx, "k2"); err == nil {
			t.Fatal("unknown kid trusted")
		}
	}
	if count() != 2 {
		t.Fatalf("unknown kid refetches: %d, want 1 more", count()-1)
	}
	// A new key is picked up after the minute.
	set(func() { kids = []string{"k1", "k2"} })
	now = now.Add(61 * time.Second)
	if _, err := j.Key(ctx, "k2"); err != nil {
		t.Fatalf("new key after a minute: %v", err)
	}
	// A removed key is not trusted after the next refresh.
	set(func() { kids = []string{"k2"} })
	now = now.Add(5 * time.Minute)
	if _, err := j.Key(ctx, "k1"); err == nil {
		t.Fatal("a key removed from the set is still trusted")
	}
	// KIFF unreachable: the last good set holds for an hour, then nothing.
	set(func() { down = true })
	now = now.Add(30 * time.Minute)
	if _, err := j.Key(ctx, "k2"); err != nil {
		t.Fatalf("within the hour: %v", err)
	}
	now = now.Add(31 * time.Minute)
	if _, err := j.Key(ctx, "k2"); err == nil {
		t.Fatal("trusted a key set older than an hour")
	}
}
