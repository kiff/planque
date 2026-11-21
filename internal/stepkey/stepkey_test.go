package stepkey

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"math/big"
	"testing"
)

// A signature is over label ‖ 0x00 ‖ JCS(body), as raw r‖s: it verifies for
// that label and body only, whatever the key order, and never for another
// key, label or body.
func TestASignatureHoldsForItsLabelAndBodyOnly(t *testing.T) {
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	l := Local{Key: k}
	ctx := context.Background()
	sig, err := SignBody(ctx, l, "tafs-chrome/step/v1", map[string]any{"id": "st_1", "op": "read"})
	if err != nil {
		t.Fatal(err)
	}
	if !Verify(l.PublicKey(), "tafs-chrome/step/v1", map[string]any{"op": "read", "id": "st_1"}, sig) {
		t.Fatal("its own body, keys in another order, didn't verify")
	}
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	for name, ok := range map[string]bool{
		"another label": Verify(l.PublicKey(), "tafs-chrome/welcome/v1", map[string]any{"id": "st_1", "op": "read"}, sig),
		"another body":  Verify(l.PublicKey(), "tafs-chrome/step/v1", map[string]any{"id": "st_1", "op": "act"}, sig),
		"another key":   Verify(Local{Key: other}.PublicKey(), "tafs-chrome/step/v1", map[string]any{"id": "st_1", "op": "read"}, sig),
		"not a sig":     Verify(l.PublicKey(), "tafs-chrome/step/v1", map[string]any{"id": "st_1", "op": "read"}, "AAAA"),
	} {
		if ok {
			t.Errorf("%s verified", name)
		}
	}
}

// KMS answers DER; the extension reads r‖s, 32 bytes each, left-padded: a
// short r or s (one in 256 signatures) keeps its place.
func TestRawFromDERPadsAShortRAndS(t *testing.T) {
	r, s := big.NewInt(5), new(big.Int).Lsh(big.NewInt(1), 250)
	der, _ := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	raw, err := RawFromDER(der)
	if err != nil || len(raw) != 64 || raw[31] != 5 || raw[0] != 0 || new(big.Int).SetBytes(raw[32:]).Cmp(s) != 0 {
		t.Fatalf("%x %v", raw, err)
	}
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	for range 200 {
		msg := make([]byte, 16)
		_, _ = rand.Read(msg)
		sum := sha256.Sum256(msg)
		der, _ := ecdsa.SignASN1(rand.Reader, k, sum[:])
		raw, err := RawFromDER(der)
		if err != nil || !ecdsa.Verify(&k.PublicKey, sum[:], new(big.Int).SetBytes(raw[:32]), new(big.Int).SetBytes(raw[32:])) {
			t.Fatalf("a DER signature didn't survive: %v", err)
		}
	}
	if _, err := RawFromDER(append(der, 0)); err == nil {
		t.Fatal("trailing bytes were taken")
	}
}

// Only a P-256 key counts: a P-224 key's r and s, padded to 32 bytes each,
// make a well-formed 64-byte signature that would otherwise verify.
func TestAKeyOnAnotherCurveNeverVerifies(t *testing.T) {
	k, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	body := map[string]any{"id": "st_1"}
	msg, err := Bytes("tafs-chrome/result/v1", body)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, k, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	raw := make([]byte, 64)
	r.FillBytes(raw[:32])
	s.FillBytes(raw[32:])
	der, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ecdsa.Verify(&k.PublicKey, sum[:], new(big.Int).SetBytes(raw[:32]), new(big.Int).SetBytes(raw[32:])) {
		t.Fatal("the padded P-224 signature isn't valid on its own curve")
	}
	if Verify(base64.RawURLEncoding.EncodeToString(der), "tafs-chrome/result/v1", body, base64.RawURLEncoding.EncodeToString(raw)) {
		t.Fatal("a P-224 key verified")
	}
}
