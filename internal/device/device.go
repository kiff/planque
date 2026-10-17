// Package device checks what an owner's phone signs. The key is the one
// tafs keeps in the Secure Enclave (ApprovalSigner): P-256, its public half
// sent as base64 SPKI DER, each signature a base64 DER ECDSA signature over
// the message's SHA-256.
//
// A signed call carries Planque-Device-Signature: t=<signed_at>,v1=<sig>
// over the same shape as KIFF's kiff.request.v1, with planque's own prefix
// so a signature for one service is never valid at the other, and the
// planque it is for, so it is never valid at another planque either:
//
//	planque.request.v1
//	<planque base URL>      e.g. https://planque.example.com, no trailing slash
//	<tenant>
//	<METHOD>
//	<path and query>
//	<sha256(body) hex>
//	<signed_at>             RFC 3339
package device

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// Header carries a request's signature.
const Header = "Planque-Device-Signature"

// MaxSkew is how far signed_at may be from planque's clock.
const MaxSkew = 5 * time.Minute

var (
	ErrKey       = errors.New("the device key must be a P-256 public key, base64 SPKI DER")
	ErrSignature = errors.New("the device signature does not verify for this request")
	ErrSignedAt  = errors.New("signed_at must be RFC 3339, within 5 minutes of now")
	ErrHeader    = errors.New(Header + " must be t=<signed_at>,v1=<signature>")
)

// ParseKey decodes a base64 SPKI DER key and accepts only P-256 ECDSA.
func ParseKey(b64 string) ([]byte, error) {
	der, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
	if err != nil || len(der) == 0 || len(der) > 512 {
		return nil, ErrKey
	}
	if _, err := publicKey(der); err != nil {
		return nil, err
	}
	return der, nil
}

func publicKey(der []byte) (*ecdsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, ErrKey
	}
	k, ok := pub.(*ecdsa.PublicKey)
	if !ok || k.Curve != elliptic.P256() {
		return nil, ErrKey
	}
	return k, nil
}

// Fingerprint is a short, readable name for a key, as KIFF shows it.
func Fingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	h := hex.EncodeToString(sum[:8])
	return h[0:4] + " " + h[4:8] + " " + h[8:12] + " " + h[12:16]
}

// Message is the exact text a phone signs for a request to the planque at
// base.
func Message(base, tenant, method, uri string, body []byte, signedAt string) []byte {
	sum := sha256.Sum256(body)
	return []byte(strings.Join([]string{"planque.request.v1", strings.TrimRight(base, "/"), tenant, method, uri,
		hex.EncodeToString(sum[:]), signedAt}, "\n"))
}

// ParseHeader splits t=<signed_at>,v1=<signature>.
func ParseHeader(h string) (signedAt, sig string, err error) {
	for part := range strings.SplitSeq(h, ",") {
		k, v, _ := strings.Cut(strings.TrimSpace(part), "=")
		switch k {
		case "t":
			signedAt = v
		case "v1":
			sig = v // Cut splits at the first '=', so base64 padding stays
		}
	}
	if signedAt == "" || sig == "" {
		return "", "", ErrHeader
	}
	return signedAt, sig, nil
}

// Verify checks sigB64 over msg with the key, and that signedAt is near now.
func Verify(keyDER, msg []byte, signedAt, sigB64 string, now time.Time) error {
	at, err := time.Parse(time.RFC3339, signedAt)
	if err != nil {
		return ErrSignedAt
	}
	if d := now.Sub(at); d > MaxSkew || d < -MaxSkew {
		return ErrSignedAt
	}
	k, err := publicKey(keyDER)
	if err != nil {
		return err
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(sigB64))
	if err != nil || len(sig) == 0 {
		return ErrSignature
	}
	digest := sha256.Sum256(msg)
	if !ecdsa.VerifyASN1(k, digest[:], sig) {
		return ErrSignature
	}
	return nil
}
