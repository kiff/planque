package planque

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"slices"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/stepkey"
)

// labelPair is what the phone signs to confirm a pairing (kiff/tafs
// chrome/src/crypto.js LABEL.pair).
const labelPair = "tafs-chrome/pair/v1"

// ChromePairedWith is the re-pairing that carried connections to a Chrome,
// as the phone signed it for the extension ("Your Chrome" 1): the pairing's
// body, which names each connection carried with its sites, and the phone's
// signature over it. A connection carried here keeps no signature of its
// own, so this is what vouches for its sites: a phone re-pairing this Chrome
// checks it with its own key before carrying the connection again, and
// planque takes it again only when this holds.
type ChromePairedWith struct {
	Body     chromePairBody `json:"body"`
	PhoneSig string         `json:"phone_sig"`
	// The one-fingerprint form (Your Chrome 5): phone_sig is then the approval
	// key's kiff.requests.v1 signature over this body's request and KIFF's.
	Tenant            string `json:"tenant,omitempty"`
	SignedAt          string `json:"signed_at,omitempty"`
	KiffRequestSHA256 string `json:"kiff_request_sha256,omitempty"`
}

// chromePairBody is the body the phone signs and the extension checks
// (kiff/tafs chrome/src/pairing.js checkConfirmation), field for field.
type chromePairBody struct {
	Pairing         string          `json:"pairing"`
	ChromeKeySHA256 string          `json:"chrome_key_sha256"`
	Name            string          `json:"name"`
	PhoneKey        string          `json:"phone_key"`
	PhoneProof      string          `json:"phone_proof"`
	Carried         []ChromeCarried `json:"carried"`
}

// holds reports whether PhoneSig is the body's phone key's signature over it,
// in the form it was given: alone (tafs-chrome/pair/v1), or beside KIFF's
// request under one fingerprint, as the extension checks it.
func (p *ChromePairedWith) holds() bool {
	if p.Tenant+p.SignedAt+p.KiffRequestSHA256 == "" {
		return stepkey.Verify(p.Body.PhoneKey, labelPair, p.Body, p.PhoneSig)
	}
	if p.Tenant == "" || !oneTimeLine(p.SignedAt) || !hex64(p.KiffRequestSHA256) {
		return false
	}
	own, err := stepkey.Bytes(labelPair, p.Body)
	if err != nil {
		return false
	}
	sum := sha256.Sum256(own)
	ownHex := hex.EncodeToString(sum[:])
	if ownHex == p.KiffRequestSHA256 {
		return false
	}
	return verifyRaw(p.Body.PhoneKey, device.RequestsMessage(p.Tenant, p.SignedAt, ownHex, p.KiffRequestSHA256), p.PhoneSig)
}

// vouches reports whether this pairing, the one that made c, carried g to c
// with exactly the sites g has, signed by phoneKey: the key re-pairing c now.
// Anything else (another Chrome's key, another phone, other sites, a
// signature that doesn't hold) vouches for nothing.
func (p *ChromePairedWith) vouches(c *Chrome, phoneKey string, g *Grant) bool {
	if p == nil || g.Request.Browser == nil || p.Body.ChromeKeySHA256 != c.KeySHA256 || p.Body.PhoneKey != phoneKey {
		return false
	}
	named := slices.ContainsFunc(p.Body.Carried, func(m ChromeCarried) bool {
		return m.Grant == g.ID && slices.Equal(m.Sites, g.Request.Browser.Sites)
	})
	return named && p.holds()
}

// verifyRaw reports whether sig (raw r‖s, base64url) is spki's P-256
// signature of msg; spki is SPKI DER in base64url. Anything malformed is false.
func verifyRaw(spki string, msg []byte, sig string) bool {
	der, err := base64.RawURLEncoding.DecodeString(spki)
	if err != nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return false
	}
	k, ok := pub.(*ecdsa.PublicKey)
	if !ok || k.Curve.Params().Name != "P-256" {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(sig)
	if err != nil || len(raw) != 64 {
		return false
	}
	sum := sha256.Sum256(msg)
	return ecdsa.Verify(k, sum[:], new(big.Int).SetBytes(raw[:32]), new(big.Int).SetBytes(raw[32:]))
}
