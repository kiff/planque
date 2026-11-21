package planque

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const releasedManifest = "aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11aa11"

func p256Key(t *testing.T) (string, string) {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(der)
	return base64.RawURLEncoding.EncodeToString(der), hex.EncodeToString(sum[:])
}

func b64(n int) string { return base64.RawURLEncoding.EncodeToString(make([]byte, n)) }

func chromeEnv(t *testing.T) *env {
	t.Helper()
	e := newEnv(t)
	e.enrol()
	e.s.ChromeBuilds = map[string]ChromeBuild{releasedManifest: {ManifestSHA256: releasedManifest, Version: "1.0.0"}}
	return e
}

func extensionPosts(e *env, lookup, key string, manifest string) (int, map[string]any) {
	return e.do("POST", "/v1/extension/pairings/"+lookup, map[string]any{"key": key, "version": "1.0.0", "manifest_hash": manifest,
		"browser": "Chrome 131", "os": "macOS 15", "proof": b64(32)}, nil, "")
}

// The whole pairing: the phone opens it, the extension posts its key, the
// phone reads and confirms it, the extension collects the phone's key and
// proofs once; the Chrome is listed, and unpairing removes it.
func TestAChromeIsPairedByThePhoneAndPlanqueRelaysBothKeys(t *testing.T) {
	e := chromeEnv(t)
	acct := "/v1/accounts/" + tenant
	st, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, "")
	if st != 201 || !lookupPattern.MatchString(out["lookup"].(string)) {
		t.Fatalf("open: %d %v", st, out)
	}
	pairing, lookup := out["pairing"].(string), out["lookup"].(string)

	extKey, extSum := p256Key(t)
	if st, out := extensionPosts(e, lookup, extKey, releasedManifest); st != 202 {
		t.Fatalf("the extension's key: %d %v", st, out)
	}
	// Before the phone confirms, the extension learns nothing.
	if st, out := e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+extKey, nil, nil, ""); st != 202 || out["status"] != "waiting_for_phone" {
		t.Fatalf("collect early: %d %v", st, out)
	}
	st, out = e.do("GET", acct+"/chromes/pairings/"+pairing, nil, e.phone, "")
	if st != 200 || out["key"] != extKey || out["key_sha256"] != extSum || out["proof"] == "" || out["browser"] != "Chrome 131" {
		t.Fatalf("the phone's read: %d %v", st, out)
	}
	phoneKey, _ := p256Key(t)
	st, out = e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", map[string]any{"chrome_key_sha256": extSum, "name": "Work Mac",
		"phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)}, e.phone, "")
	if st != 200 || !chromeIDRe.MatchString(out["chrome_id"].(string)) {
		t.Fatalf("confirm: %d %v", st, out)
	}
	chrome := out["chrome_id"].(string)

	// Another key collects nothing; the posting key collects the phone's
	// key and proofs, once.
	if st, _ := e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+phoneKey, nil, nil, ""); st != 404 {
		t.Fatalf("another key collected: %d", st)
	}
	st, out = e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+extKey, nil, nil, "")
	if st != 200 || out["chrome_id"] != chrome || out["phone_key"] != phoneKey || out["name"] != "Work Mac" || out["phone_sig"] != b64(64) {
		t.Fatalf("collect: %d %v", st, out)
	}
	if st, _ := e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+extKey, nil, nil, ""); st != 404 {
		t.Fatalf("collected twice: %d", st)
	}

	st, out = e.do("GET", acct+"/chromes", nil, e.phone, "")
	if list := out["chromes"].([]any); st != 200 || len(list) != 1 || list[0].(map[string]any)["name"] != "Work Mac" {
		t.Fatalf("list: %d %v", st, out)
	}
	if st, _ := e.do("DELETE", acct+"/chromes/"+chrome, nil, e.phone, ""); st != 200 {
		t.Fatalf("unpair: %d", st)
	}
	// The phone signs each request anew: the same request in the same second
	// is a replay.
	anew := func() { time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second))) }
	anew()
	if _, out := e.do("GET", acct+"/chromes", nil, e.phone, ""); len(out["chromes"].([]any)) != 0 {
		t.Fatalf("still listed after unpairing: %v", out)
	}
	anew()
	if st, _ := e.do("DELETE", acct+"/chromes/"+chrome, nil, e.phone, ""); st != 404 {
		t.Fatalf("unpaired twice: %d", st)
	}
}

// Each refusal: off with no released build; an unreleased build; a second
// Chrome on the same code; a confirmation of another key than the one the
// phone checked; an expired code; the phone's own routes need its
// signature.
func TestAChromePairingRefusesWhatItCantTrust(t *testing.T) {
	acct := "/v1/accounts/" + tenant
	t.Run("off", func(t *testing.T) {
		e := newEnv(t)
		e.enrol()
		if st, _ := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, ""); st != 404 {
			t.Fatalf("offered with no released build: %d", st)
		}
	})
	open := func(t *testing.T, e *env) (string, string) {
		_, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, "")
		return out["pairing"].(string), out["lookup"].(string)
	}
	t.Run("unreleased build", func(t *testing.T) {
		e := chromeEnv(t)
		_, lookup := open(t, e)
		key, _ := p256Key(t)
		if st, out := extensionPosts(e, lookup, key, strings.Repeat("b", 64)); st != 403 || !strings.Contains(out["error"].(string), "isn't a released one") {
			t.Fatalf("an unreleased build: %d %v", st, out)
		}
	})
	t.Run("a second Chrome", func(t *testing.T) {
		e := chromeEnv(t)
		pairing, lookup := open(t, e)
		first, firstSum := p256Key(t)
		second, _ := p256Key(t)
		if st, _ := extensionPosts(e, lookup, first, releasedManifest); st != 202 {
			t.Fatal("first")
		}
		// The second post answers as an unknown code does, and spends it.
		if st, _ := extensionPosts(e, lookup, second, releasedManifest); st != 404 {
			t.Fatalf("a second Chrome took the code: %d", st)
		}
		if _, out := e.do("GET", acct+"/chromes/pairings/"+pairing, nil, e.phone, ""); out["status"] != "contested" || out["key"] != nil {
			t.Fatalf("the phone wasn't told: %v", out)
		}
		if st, _ := e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+first, nil, nil, ""); st != 404 {
			t.Fatalf("the first Chrome still collects a contested code: %d", st)
		}
		phoneKey, _ := p256Key(t)
		st, _ := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", map[string]any{"chrome_key_sha256": firstSum, "name": "Work Mac",
			"phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)}, e.phone, "")
		if st != 409 {
			t.Fatalf("a contested code was confirmed: %d", st)
		}
	})
	t.Run("another key confirmed", func(t *testing.T) {
		e := chromeEnv(t)
		pairing, lookup := open(t, e)
		key, _ := p256Key(t)
		extensionPosts(e, lookup, key, releasedManifest)
		phoneKey, _ := p256Key(t)
		_, otherSum := p256Key(t)
		if st, _ := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", map[string]any{"chrome_key_sha256": otherSum, "name": "Work Mac",
			"phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)}, e.phone, ""); st != 409 {
			t.Fatalf("confirmed a key the phone didn't check: %d", st)
		}
		if _, out := e.do("GET", acct+"/chromes", nil, e.phone, ""); len(out["chromes"].([]any)) != 0 {
			t.Fatal("a Chrome was kept")
		}
	})
	t.Run("expired", func(t *testing.T) {
		e := chromeEnv(t)
		_, lookup := open(t, e)
		e.s.Now = func() time.Time { return time.Now().Add(chromePairingTTL + time.Minute) }
		key, _ := p256Key(t)
		if st, _ := extensionPosts(e, lookup, key, releasedManifest); st != 404 {
			t.Fatalf("an expired code: %d", st)
		}
	})
	t.Run("the phone's routes are signed", func(t *testing.T) {
		e := chromeEnv(t)
		if st, _ := e.do("POST", acct+"/chromes/pairings", map[string]any{}, nil, ""); st != 401 {
			t.Fatalf("opened without the phone: %d", st)
		}
		if st, _ := e.do("GET", acct+"/chromes", nil, nil, ""); st != 401 {
			t.Fatalf("listed without the phone: %d", st)
		}
	})
}

func TestChromeBuildsFrom(t *testing.T) {
	got, err := ChromeBuildsFrom(strings.NewReader("# released\n" + releasedManifest + " 1.0.0\n\n"))
	if err != nil || got[releasedManifest].Version != "1.0.0" || len(got) != 1 {
		t.Fatalf("%v %v", got, err)
	}
	if _, err := ChromeBuildsFrom(strings.NewReader("ABC 1.0.0\n")); err == nil {
		t.Fatal("a bad line was taken")
	}
}

// Codes are drawn only from the 30 characters the extension and the phones
// read, each about as often as another.
func TestALookupIsInTheContractsAlphabet(t *testing.T) {
	const contract = "23456789ABCDEFGHJKMNPQRSTVWXYZ"
	seen := map[rune]int{}
	for range 3000 {
		l := newLookup()
		if !lookupPattern.MatchString(l) {
			t.Fatalf("%q isn't a lookup", l)
		}
		for _, c := range l {
			if !strings.ContainsRune(contract, c) {
				t.Fatalf("%q has %q, which the extension can't read", l, c)
			}
			seen[c]++
		}
	}
	if len(seen) != len(contract) {
		t.Fatalf("drew %d of %d characters", len(seen), len(contract))
	}
	if lookupPattern.MatchString("UUUU") {
		t.Fatal("the pattern takes U")
	}
}

// A byte at or above 240 (8·30) would favour the first characters; it's
// drawn again, never folded in.
func TestALookupDrawsAgainOverTheLastWholeRound(t *testing.T) {
	bytes := []byte{240, 255, 0, 29, 30, 241, 59}
	i := 0
	if got := lookupFrom(func() byte { i++; return bytes[i-1] }); got != "2Z2Z" {
		t.Fatalf("%q", got)
	}
}

// The extension's routes carry no signature, so a sweep is limited: per
// address on posts and on collects, and across every address on posts.
func TestTheExtensionsRoutesAreLimited(t *testing.T) {
	t.Run("posts from one address", func(t *testing.T) {
		e := chromeEnv(t)
		key, _ := p256Key(t)
		for i := range chromePostsPerAddress {
			if st, _ := extensionPosts(e, "2222", key, releasedManifest); st != 404 {
				t.Fatalf("post %d: %d", i, st)
			}
		}
		if st, _ := extensionPosts(e, "2223", key, releasedManifest); st != 429 {
			t.Fatalf("post over the limit: %d", st)
		}
	})
	t.Run("collects from one address", func(t *testing.T) {
		e := chromeEnv(t)
		for i := range chromeReadsPerAddress {
			if st, _ := e.do("GET", "/v1/extension/pairings/2222?key=x", nil, nil, ""); st != 404 {
				t.Fatalf("collect %d: %d", i, st)
			}
		}
		if st, _ := e.do("GET", "/v1/extension/pairings/2222?key=x", nil, nil, ""); st != 429 {
			t.Fatalf("collect over the limit: %d", st)
		}
	})
	t.Run("posts from many addresses", func(t *testing.T) {
		e := chromeEnv(t)
		h := e.s.Handler()
		post := func(addr string) int {
			r := httptest.NewRequest("POST", "/v1/extension/pairings/2222", strings.NewReader(`{}`))
			r.RemoteAddr = addr
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			return w.Code
		}
		for i := range chromePostsInAll {
			if st := post(fmt.Sprintf("198.51.%d.%d:443", i/250, i%250)); st == 429 {
				t.Fatalf("post %d refused under the limit", i)
			}
		}
		if st := post("203.0.113.9:443"); st != 429 {
			t.Fatalf("post over the limit in all: %d", st)
		}
	})
}

// The extension's view-only KIFF device offer rides its key post and reaches
// the phone as given; the collect answer names the pairing the phone signed,
// which the extension needs to check that signature (kiff/tafs docs/chrome.md).
func TestTheKiffOfferReachesThePhoneAndTheCollectNamesThePairing(t *testing.T) {
	e := chromeEnv(t)
	acct := "/v1/accounts/" + tenant
	_, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, "")
	pairing, lookup := out["pairing"].(string), out["lookup"].(string)
	extKey, extSum := p256Key(t)
	offer := `{"exp":"2026-11-21T12:04:30Z","key_store":"software","kind":"extension","name":"tafs for Chrome","platform":"chrome","v":"kiff.pair.v1"}`

	big := strings.Repeat("x", chromeKiffOfferMax+1)
	if st, _ := e.do("POST", "/v1/extension/pairings/"+lookup, map[string]any{"key": extKey, "version": "1.0.0", "manifest_hash": releasedManifest,
		"browser": "Chrome 131", "os": "macOS 15", "proof": b64(32), "kiff_offer": big}, nil, ""); st != 400 {
		t.Fatalf("an offer over 4 KiB: %d", st)
	}
	if st, out := e.do("POST", "/v1/extension/pairings/"+lookup, map[string]any{"key": extKey, "version": "1.0.0", "manifest_hash": releasedManifest,
		"browser": "Chrome 131", "os": "macOS 15", "proof": b64(32), "kiff_offer": offer}, nil, ""); st != 202 {
		t.Fatalf("the extension's key with its KIFF offer: %d %v", st, out)
	}
	if st, out := e.do("GET", acct+"/chromes/pairings/"+pairing, nil, e.phone, ""); st != 200 || out["kiff_offer"] != offer {
		t.Fatalf("the phone reads the offer as given: %d %v", st, out)
	}
	phoneKey, _ := p256Key(t)
	if st, out := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", map[string]any{"chrome_key_sha256": extSum, "name": "Work Mac",
		"phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64)}, e.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	st, out := e.do("GET", "/v1/extension/pairings/"+lookup+"?key="+extKey, nil, nil, "")
	if st != 200 || out["pairing"] != pairing || out["phone_key"] != phoneKey {
		t.Fatalf("the collect names the pairing: %d %v", st, out)
	}
}
