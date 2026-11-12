package fakes

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/device"
)

// Phone signs as tafs' ApprovalSigner does: a P-256 key, base64 SPKI DER
// public half, base64 DER ECDSA signatures over SHA-256.
type Phone struct {
	Key *ecdsa.PrivateKey
	// DeviceID is KIFF's id for this device, sent as d= (kiff/planque#10);
	// empty for a phone enrolled with its key before that.
	DeviceID string
	// Window, when set, is the approval window the phone signs in
	// (kiff/tafs#494): Sign makes a v2 signature naming it.
	Window string
	// With, when set, is the lowercase hex SHA-256 of another request's
	// message: the phone signs both under one signature (kiff.requests.v1).
	With string
	// WithItself, for a test, pairs the request with its own message.
	WithItself bool
}

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
	msg := device.MessageWindow(req.URL.Scheme+"://"+req.URL.Host, tenant, req.Method, req.URL.RequestURI(), body, at, p.Window)
	h := "t=" + at + ",v1=" + p.SignBytes(msg)
	if p.Window != "" {
		h = "t=" + at + ",v2=" + p.SignBytes(msg) + ",w=" + p.Window
	}
	if p.With != "" || p.WithItself {
		own := sha256.Sum256(msg)
		with := p.With
		if p.WithItself {
			with = hex.EncodeToString(own[:])
		}
		h = "t=" + at + ",v3=" + p.SignBytes(device.RequestsMessage(tenant, at, hex.EncodeToString(own[:]), with)) + ",with=sha256:" + with
	}
	if p.DeviceID != "" {
		h += ",d=" + p.DeviceID
	}
	req.Header.Set(device.Header, h)
}

// Verifies says whether sigB64 is this phone's signature over msg.
func (p *Phone) Verifies(msg []byte, sigB64 string) bool {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return false
	}
	sum := sha256.Sum256(msg)
	return ecdsa.VerifyASN1(&p.Key.PublicKey, sum[:], sig)
}
