package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// testPhone is an owner's phone: its approval key, which signs a pairing for
// the extension (tafs-chrome/pair/v1, or kiff.requests.v1 beside KIFF's).
type testPhone struct {
	k    *ecdsa.PrivateKey
	spki string
}

func newTestPhone(t *testing.T) testPhone {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return testPhone{k, base64.RawURLEncoding.EncodeToString(der)}
}

// raw signs msg as the phone sends it: r‖s, base64url.
func (p testPhone) raw(t *testing.T, msg []byte) string {
	t.Helper()
	sum := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, p.k, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	out := make([]byte, 64)
	r.FillBytes(out[:32])
	s.FillBytes(out[32:])
	return base64.RawURLEncoding.EncodeToString(out)
}

func (p testPhone) signPair(t *testing.T, body chromePairBody) string {
	t.Helper()
	msg, err := stepkey.Bytes(labelPair, body)
	if err != nil {
		t.Fatal(err)
	}
	return p.raw(t, msg)
}

// A connection one re-pairing carried keeps no signature of its own; the
// pairing that carried it, as the phone signed it for the extension, is what
// vouches for its sites ("Your Chrome" 1). planque keeps that pairing on the
// new Chrome and lists it (paired_with), so the phone can check it with its
// own key and carry the connection again. planque takes it again only when
// that pairing vouches for it: made for that Chrome's key, signed by the phone
// re-pairing now, naming the connection with exactly its sites, its signature
// holding. Otherwise the re-pairing is refused, and a re-pairing that leaves
// the connection out ends it with the old Chrome, as before.
// Category: authority
// Attack: carrying again a connection whose earlier pairing names another Chrome's key, another phone, other sites, or has a signature that doesn't hold
// Refused: each, with nothing moved; left out, it ends with the old Chrome
func TestAConnectionCarriedOnceIsCarriedAgainOnlyAsItsPairingVouches(t *testing.T) {
	c := newChromeRig(t)
	c.s.Browser = nil
	ctx := context.Background()
	acct := "/v1/accounts/" + tenant
	st, out := c.do("POST", acct+"/catalog/browser@1/grants", map[string]any{"agent": "agt_wiki", "agent_name": "Wiki titles",
		"tools": []any{map[string]any{"tool": browseOpen}}, "browser": map[string]any{"sites": []string{"wikipedia.org"}}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	g, _ := c.s.loadGrant(ctx, tenant, id)
	g.ChromeSitesSig, g.Status = testSitesSig, StatusConnected
	if err := c.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	phone := newTestPhone(t)
	carried := []ChromeCarried{{Grant: id, Sites: []string{"wikipedia.org"}}}

	// rePair opens a pairing, lets a new extension post its key, and confirms
	// it replacing from, carrying what's given, signed by signer for the
	// extension (the body names phone's key, as the phone sends it).
	rePair := func(from string, carry []ChromeCarried, signer testPhone) (int, map[string]any) {
		t.Helper()
		// Each its own signed request: the same one within a second is a replay.
		time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
		st, out := c.do("POST", acct+"/chromes/pairings", map[string]any{}, c.phone, "")
		if st != 201 {
			t.Fatalf("open: %d %v", st, out)
		}
		pairing, lookup := out["pairing"].(string), out["lookup"].(string)
		extKey, extSum := p256Key(t)
		if st, out := extensionPosts(c.env, lookup, extKey, releasedManifest); st != 202 {
			t.Fatalf("post: %d %v", st, out)
		}
		body := chromePairBody{Pairing: pairing, ChromeKeySHA256: extSum, Name: "Work Mac", PhoneKey: signer.spki, PhoneProof: b64(32), Carried: carry}
		send := map[string]any{"chrome_key_sha256": extSum, "name": "Work Mac", "phone_key": signer.spki, "phone_proof": b64(32),
			"phone_sig": signer.signPair(t, body), "replaces": from}
		if len(carry) > 0 {
			send["carried"] = carry
		}
		return c.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", send, c.phone, "")
	}
	chromeOf := func(chromeID string) Chrome {
		t.Helper()
		var ch Chrome
		if _, err := getDoc(ctx, c.s.KV, tenantPK(tenant), skChrome(chromeID), &ch); err != nil {
			t.Fatal(err)
		}
		return ch
	}
	on := func() *Grant {
		t.Helper()
		g, err := c.s.loadGrant(ctx, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
		return g
	}

	// The first re-pairing carries the connection, whose sites the phone signed:
	// the new Chrome keeps that pairing, as signed, and the list gives it.
	st, out = rePair(c.r.chrome, carried, phone)
	if st != 200 {
		t.Fatalf("first re-pairing: %d %v", st, out)
	}
	b := out["chrome_id"].(string)
	if g := on(); g.Request.Browser.Chrome != b || g.ChromeSitesSig != "" || g.ChromeSitesCarried != b {
		t.Fatalf("after the first re-pairing: %+v sig=%q carried=%q", g.Request.Browser, g.ChromeSitesSig, g.ChromeSitesCarried)
	}
	kept := chromeOf(b)
	if kept.PairedWith == nil || !kept.PairedWith.holds() || kept.PairedWith.Body.PhoneKey != phone.spki ||
		kept.PairedWith.Body.ChromeKeySHA256 != kept.KeySHA256 || len(kept.PairedWith.Body.Carried) != 1 {
		t.Fatalf("the pairing kept on the new Chrome: %+v", kept.PairedWith)
	}
	st, out = c.do("GET", acct+"/chromes", nil, c.phone, "")
	listed := false
	for _, row := range out["chromes"].([]any) {
		if m := row.(map[string]any); m["chrome_id"] == b {
			pw, _ := m["paired_with"].(map[string]any)
			body, _ := pw["body"].(map[string]any)
			listed = pw["phone_sig"] == kept.PairedWith.PhoneSig && body["phone_key"] == phone.spki && body["chrome_key_sha256"] == kept.KeySHA256
		}
	}
	if st != 200 || !listed {
		t.Fatalf("the list doesn't give the pairing: %d %v", st, out)
	}

	// Carried again from that Chrome, it's taken only as its pairing vouches.
	// Each way the pairing doesn't is refused, and the connection stays put.
	for what, edit := range map[string]func(p *ChromePairedWith){
		"made for another Chrome's key": func(p *ChromePairedWith) { p.Body.ChromeKeySHA256 = strings.Repeat("0", 64) },
		"with other sites":              func(p *ChromePairedWith) { p.Body.Carried[0].Sites = []string{"wikipedia.org", "bank.example"} },
		"not naming the connection":     func(p *ChromePairedWith) { p.Body.Carried[0].Grant = "grt_someoneelse0" },
		"a signature that doesn't hold": func(p *ChromePairedWith) { p.PhoneSig = b64(64) },
		"none kept":                     nil,
	} {
		ch := chromeOf(b)
		was := ch.PairedWith
		if edit == nil {
			ch.PairedWith = nil
		} else {
			cp := *was
			cp.Body.Carried = []ChromeCarried{{Grant: was.Body.Carried[0].Grant, Sites: append([]string(nil), was.Body.Carried[0].Sites...)}}
			edit(&cp)
			// Signed as it reads by the same phone, so only the part named fails:
			// the signature itself is its own case.
			if what != "a signature that doesn't hold" {
				cp.PhoneSig = phone.signPair(t, cp.Body)
			}
			ch.PairedWith = &cp
		}
		if err := putDoc(ctx, c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(b)}, ch, store.Cond{}); err != nil {
			t.Fatal(err)
		}
		if st, out := rePair(b, carried, phone); st != 409 {
			t.Errorf("%s: %d %v, want 409", what, st, out)
		}
		ch.PairedWith = was
		if err := putDoc(ctx, c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(b)}, ch, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	// Another phone, its own key signing: the pairing that carried it here
	// wasn't this phone's, so it doesn't vouch.
	if st, out := rePair(b, carried, newTestPhone(t)); st != 409 {
		t.Errorf("another phone: %d %v, want 409", st, out)
	}
	if g := on(); g.Request.Browser.Chrome != b || g.Status != StatusConnected {
		t.Fatalf("a refused re-pairing moved or ended the connection: %+v %s", g.Request.Browser, g.Status)
	}

	// As its pairing vouches, by the same phone: it moves again, and the next
	// Chrome keeps this pairing in turn.
	st, out = rePair(b, carried, phone)
	if st != 200 {
		t.Fatalf("second re-pairing: %d %v", st, out)
	}
	cc := out["chrome_id"].(string)
	if g := on(); g.Request.Browser.Chrome != cc || g.ChromeSitesCarried != cc || g.Status != StatusConnected {
		t.Fatalf("after the second re-pairing: %+v %s", g.Request.Browser, g.Status)
	}
	if pw := chromeOf(cc).PairedWith; pw == nil || !pw.holds() {
		t.Fatalf("the second pairing wasn't kept: %+v", pw)
	}

	// Left out of a re-pairing, it ends with the old Chrome, as before.
	if st, out := rePair(cc, nil, phone); st != 200 {
		t.Fatalf("a re-pairing that carries nothing: %d %v", st, out)
	}
	if g := on(); g.Status != StatusDisconnected {
		t.Fatalf("left out, it didn't end: %s", g.Status)
	}
}

// The one-fingerprint form (Your Chrome 5): the phone's signature is its
// approval key's kiff.requests.v1 over this pairing's request and KIFF's, as
// the extension checks it; a pairing signed so holds, and the same signature
// read with any part changed, or paired with itself, doesn't.
func TestAPairingSignedBesideKIFFsHoldsOnlyAsSigned(t *testing.T) {
	phone := newTestPhone(t)
	body := chromePairBody{Pairing: "pair_1", ChromeKeySHA256: strings.Repeat("ab", 32), Name: "Work Mac", PhoneKey: phone.spki,
		PhoneProof: b64(32), Carried: []ChromeCarried{{Grant: "grt_wiki0000000001", Sites: []string{"wikipedia.org"}}}}
	own, err := stepkey.Bytes(labelPair, body)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(own)
	kiff := strings.Repeat("c", 64)
	at := "2026-11-23T01:00:00Z"
	pw := ChromePairedWith{Body: body, Tenant: tenant, SignedAt: at, KiffRequestSHA256: kiff,
		PhoneSig: phone.raw(t, device.RequestsMessage(tenant, at, hex.EncodeToString(sum[:]), kiff))}
	if !pw.holds() {
		t.Fatal("signed beside KIFF's request, it doesn't hold")
	}
	for what, edit := range map[string]func(p *ChromePairedWith){
		"other sites": func(p *ChromePairedWith) {
			p.Body.Carried = []ChromeCarried{{Grant: "grt_wiki0000000001", Sites: []string{"bank.example"}}}
		},
		"another account":         func(p *ChromePairedWith) { p.Tenant = "t_other" },
		"another time":            func(p *ChromePairedWith) { p.SignedAt = "2026-11-23T01:00:01Z" },
		"another KIFF request":    func(p *ChromePairedWith) { p.KiffRequestSHA256 = strings.Repeat("d", 64) },
		"read as a pairing alone": func(p *ChromePairedWith) { p.Tenant, p.SignedAt, p.KiffRequestSHA256 = "", "", "" },
		"half the form":           func(p *ChromePairedWith) { p.Tenant = "" },
	} {
		cp := pw
		edit(&cp)
		if cp.holds() {
			t.Errorf("%s: holds", what)
		}
	}
	// Paired with itself: the extension never takes it, nor does planque.
	self := pw
	self.KiffRequestSHA256 = hex.EncodeToString(sum[:])
	self.PhoneSig = phone.raw(t, device.RequestsMessage(tenant, at, self.KiffRequestSHA256, self.KiffRequestSHA256))
	if self.holds() {
		t.Error("a pairing signed beside itself holds")
	}
}
