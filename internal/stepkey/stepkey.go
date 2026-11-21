// Package stepkey signs what planque sends the owner's own Chrome: the
// welcome that proves planque on each connect, and each step. The key is
// ECDSA P-256; a signature is over label ‖ 0x00 ‖ JCS(body) and travels as
// the raw 64-byte r‖s, base64url without padding, as the extension checks it
// with WebCrypto. In production the private key never leaves KMS.
package stepkey

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/jcs"
)

// Signer signs a message with planque's step key.
type Signer interface {
	// Sign returns the raw r‖s signature (64 bytes) of SHA-256(msg).
	Sign(ctx context.Context, msg []byte) ([]byte, error)
}

// Bytes is what is signed for body under label: label ‖ 0x00 ‖ JCS(body).
func Bytes(label string, body any) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	canon, err := jcs.Canonical(raw)
	if err != nil {
		return nil, err
	}
	out := append([]byte(label), 0)
	return append(out, canon...), nil
}

// SignBody signs body under label and returns the signature as it travels.
func SignBody(ctx context.Context, s Signer, label string, body any) (string, error) {
	msg, err := Bytes(label, body)
	if err != nil {
		return "", err
	}
	sig, err := s.Sign(ctx, msg)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sig), nil
}

// Verify reports whether sig (raw r‖s, base64url) is spki's signature of
// body under label. spki is the key as SPKI DER, base64url. Anything
// malformed is false.
func Verify(spki, label string, body any, sig string) bool {
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
	msg, err := Bytes(label, body)
	if err != nil {
		return false
	}
	sum := sha256.Sum256(msg)
	return ecdsa.Verify(k, sum[:], new(big.Int).SetBytes(raw[:32]), new(big.Int).SetBytes(raw[32:]))
}

// RawFromDER turns an X9.62 DER ECDSA P-256 signature, as KMS returns it,
// into the raw 64-byte r‖s.
func RawFromDER(der []byte) ([]byte, error) {
	var sig struct{ R, S *big.Int }
	rest, err := asn1.Unmarshal(der, &sig)
	if err != nil || len(rest) != 0 || sig.R == nil || sig.S == nil || sig.R.Sign() <= 0 || sig.S.Sign() <= 0 {
		return nil, errors.New("stepkey: not a DER ECDSA signature")
	}
	if sig.R.BitLen() > 256 || sig.S.BitLen() > 256 {
		return nil, errors.New("stepkey: not a P-256 signature")
	}
	out := make([]byte, 64)
	sig.R.FillBytes(out[:32])
	sig.S.FillBytes(out[32:])
	return out, nil
}

// KMS signs with an asymmetric ECC_NIST_P256 key in AWS KMS.
type KMS struct {
	Client *awsapi.Client
	KeyID  string // the key's id, ARN or alias
}

// Sign asks KMS to sign msg (ECDSA_SHA_256 over the raw message).
func (k KMS) Sign(ctx context.Context, msg []byte) ([]byte, error) {
	var out struct {
		Signature []byte `json:"Signature"`
	}
	in := map[string]any{"KeyId": k.KeyID, "Message": msg, "MessageType": "RAW", "SigningAlgorithm": "ECDSA_SHA_256"}
	if err := k.Client.Do(ctx, "Sign", in, &out); err != nil {
		return nil, fmt.Errorf("kms sign: %w", err)
	}
	return RawFromDER(out.Signature)
}

// PublicKey returns the key's public half as SPKI DER, base64url, for the
// extension's build to pin.
func (k KMS) PublicKey(ctx context.Context) (string, error) {
	var out struct {
		PublicKey []byte `json:"PublicKey"`
	}
	if err := k.Client.Do(ctx, "GetPublicKey", map[string]any{"KeyId": k.KeyID}, &out); err != nil {
		return "", fmt.Errorf("kms public key: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(out.PublicKey), nil
}

// Local signs with an in-process key: tests and a local run only.
type Local struct{ Key *ecdsa.PrivateKey }

// Sign signs SHA-256(msg) and returns r‖s.
func (l Local) Sign(_ context.Context, msg []byte) ([]byte, error) {
	sum := sha256.Sum256(msg)
	der, err := ecdsa.SignASN1(rand.Reader, l.Key, sum[:])
	if err != nil {
		return nil, err
	}
	return RawFromDER(der)
}

// PublicKey is the key's public half as SPKI DER, base64url.
func (l Local) PublicKey() string {
	der, _ := x509.MarshalPKIXPublicKey(&l.Key.PublicKey)
	return base64.RawURLEncoding.EncodeToString(der)
}
