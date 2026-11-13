// Package device checks what an owner's phone signs. The key is the one
// tafs keeps in the Secure Enclave (ApprovalSigner): P-256, its public half
// sent as base64 SPKI DER, each signature a base64 DER ECDSA signature over
// the message's SHA-256.
//
// A signed call carries Planque-Device-Signature: t=<signed_at>,v1=<sig>
// (and d=<device id> from a device KIFF knows, kiff/planque#10) over the
// same shape as KIFF's kiff.request.v1, with planque's own prefix
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
//
// Inside an approval window (kiff/tafs#494: one fingerprint opened it at
// KIFF, then taps), the device key signs t=<signed_at>,v2=<sig>,w=<window
// id>, over the same lines under planque.request.v2 with the window's id
// last, so the signature can't be moved to another window. KIFF says
// whether that window is open for the device.
package device

import (
	"crypto/sha256"
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
	ErrSignature = errors.New("the device signature does not verify for this request")
	ErrSignedAt  = errors.New("signed_at must be RFC 3339, within 5 minutes of now")
	ErrHeader    = errors.New(Header + " must be t=<signed_at>,v1=<signature>, and d=<device id> for a device KIFF knows; inside an approval window, v2=<signature>,w=<window id>; signed together with another request, v3=<signature>,with=sha256:<hex>")
)

// Message is the exact text a phone signs for a request to the planque at
// base.
func Message(base, tenant, method, uri string, body []byte, signedAt string) []byte {
	sum := sha256.Sum256(body)
	return []byte(strings.Join([]string{"planque.request.v1", strings.TrimRight(base, "/"), tenant, method, uri,
		hex.EncodeToString(sum[:]), signedAt}, "\n"))
}

// MessageWindow is what a phone signs inside an approval window
// (kiff/tafs#494): Message's lines under planque.request.v2, with the
// window's id last. With no window, it is Message.
func MessageWindow(base, tenant, method, uri string, body []byte, signedAt, window string) []byte {
	if window == "" {
		return Message(base, tenant, method, uri, body, signedAt)
	}
	sum := sha256.Sum256(body)
	return []byte(strings.Join([]string{"planque.request.v2", strings.TrimRight(base, "/"), tenant, method, uri,
		hex.EncodeToString(sum[:]), signedAt, window}, "\n"))
}

// Signed is a request's Planque-Device-Signature.
type Signed struct {
	At        string // signed_at, RFC 3339
	Signature string // base64 DER ECDSA
	// Device is KIFF's id for the device that signed (dev_…), which KIFF
	// verifies; empty from a phone enrolled with its key before
	// kiff/planque#10, which planque verifies itself.
	Device string
	// Window is the approval window the device key signed in (win_…), a v2
	// signature (kiff/tafs#494); empty for a v1.
	Window string
	// With is a v3's other request (kiff.requests.v1): the lowercase hex
	// SHA-256 of the request message it was signed together with.
	With string
}

// ParseHeader splits t=<signed_at>,v1=<signature>[,d=<device id>], or
// t=<signed_at>,v2=<signature>,w=<window id>[,d=<device id>].
func ParseHeader(h string) (Signed, error) {
	var s Signed
	version := ""
	for part := range strings.SplitSeq(h, ",") {
		k, v, _ := strings.Cut(strings.TrimSpace(part), "=")
		switch k {
		case "t":
			s.At = v
		case "v1", "v2", "v3":
			if version != "" {
				return Signed{}, ErrHeader
			}
			version, s.Signature = k, v // Cut splits at the first '=', so base64 padding stays
		case "with":
			if s.With != "" {
				return Signed{}, ErrHeader
			}
			sum, ok := strings.CutPrefix(v, "sha256:")
			if _, err := hex.DecodeString(sum); !ok || err != nil || len(sum) != 64 || strings.ToLower(sum) != sum {
				return Signed{}, ErrHeader
			}
			s.With = sum
		case "w":
			s.Window = v
		case "d":
			s.Device = v
		}
	}
	// A v2 names its window, a v1 none.
	if s.At == "" || s.Signature == "" || (s.Device != "" && !validDeviceID(s.Device)) ||
		(version == "v2") != (s.Window != "") || (s.Window != "" && !validWindowID(s.Window)) ||
		(version == "v3") != (s.With != "") {
		return Signed{}, ErrHeader
	}
	return s, nil
}

// RequestsMessage is what a phone signs, once, with its approval key, to send
// two requests under one fingerprint (kiff.requests.v1, kiff/tafs
// docs/secure-connection.md): the tenant, signed_at, and the two request
// messages' SHA-256 in lowercase hex, sorted. Each service rebuilds its own
// request message (Message here, kiff.request.v1 at KIFF) and takes the
// other's hash as given. KIFF builds the same lines.
func RequestsMessage(tenant, signedAt, oneSHA256, otherSHA256 string) []byte {
	a, b := oneSHA256, otherSHA256
	if b < a {
		a, b = b, a
	}
	return []byte(strings.Join([]string{"kiff.requests.v1", tenant, signedAt, a, b}, "\n"))
}

// validWindowID: KIFF makes them win_ and 26 lowercase base32 characters.
func validWindowID(id string) bool {
	rest, ok := strings.CutPrefix(id, "win_")
	if !ok || len(rest) != 26 {
		return false
	}
	for _, c := range rest {
		if (c < 'a' || c > 'z') && (c < '2' || c > '7') {
			return false
		}
	}
	return true
}

// validDeviceID: KIFF makes them dev_ and lowercase base32.
func validDeviceID(id string) bool {
	rest, ok := strings.CutPrefix(id, "dev_")
	if !ok || rest == "" || len(id) > 64 {
		return false
	}
	for _, c := range rest {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// CheckTime says whether signedAt is RFC 3339 and within MaxSkew of now.
func CheckTime(signedAt string, now time.Time) error {
	at, err := time.Parse(time.RFC3339, signedAt)
	if err != nil {
		return ErrSignedAt
	}
	if d := now.Sub(at); d > MaxSkew || d < -MaxSkew {
		return ErrSignedAt
	}
	return nil
}
