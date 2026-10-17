package fakes

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/device"
)

// Phone signs as tafs' ApprovalSigner does: a P-256 key, base64 SPKI DER
// public half, base64 DER ECDSA signatures over SHA-256.
type Phone struct{ Key *ecdsa.PrivateKey }

// NewPhone makes a phone with a fresh key.
func NewPhone() *Phone {
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return &Phone{Key: k}
}

// PublicKey is the key as tafs sends it.
func (p *Phone) PublicKey() string {
	der, _ := x509.MarshalPKIXPublicKey(&p.Key.PublicKey)
	return base64.StdEncoding.EncodeToString(der)
}

// SignBytes signs msg, as signConfirmed does.
func (p *Phone) SignBytes(msg []byte) string {
	sum := sha256.Sum256(msg)
	sig, _ := ecdsa.SignASN1(rand.Reader, p.Key, sum[:])
	return base64.StdEncoding.EncodeToString(sig)
}

// Sign adds Planque-Device-Signature to req, whose body is body, for the
// planque req is sent to.
func (p *Phone) Sign(req *http.Request, tenant string, body []byte) {
	at := time.Now().UTC().Format(time.RFC3339)
	msg := device.Message(req.URL.Scheme+"://"+req.URL.Host, tenant, req.Method, req.URL.RequestURI(), body, at)
	req.Header.Set(device.Header, "t="+at+",v1="+p.SignBytes(msg))
}
