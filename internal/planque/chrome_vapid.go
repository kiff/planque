package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// pkVAPID is the one item that keeps planque's VAPID key: the key that signs
// the push waking the owner's Chrome (chrome_wake.go). Made once, on first
// need, by planque itself; never typed, sent or shown by anyone. The private
// key is sealed (envelope, with planque's KMS key hosted) under vapidBinding;
// the public key beside it is what the extension pins as its
// applicationServerKey.
const pkVAPID = "VAPID"

var vapidBinding = map[string]string{"kind": "vapid"}

type vapidDoc struct {
	// Sealed is the private key, SEC 1 DER, sealed.
	Sealed []byte `json:"sealed"`
	// Public is the uncompressed P-256 point, base64url.
	Public string `json:"public"`
}

// vapidKey is planque's VAPID key, made and kept the first time it's needed:
// a conditional write, so two instances making one at once keep one.
func (s *Server) vapidKey(ctx context.Context) (*ecdsa.PrivateKey, []byte, error) {
	var d vapidDoc
	_, err := getDoc(ctx, s.KV, pkVAPID, "-", &d)
	if errors.Is(err, errNotFound) {
		k, gerr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if gerr != nil {
			return nil, nil, gerr
		}
		der, merr := x509.MarshalECPrivateKey(k)
		if merr != nil {
			return nil, nil, merr
		}
		sealed, serr := s.Sealer.Seal(ctx, der, vapidBinding)
		if serr != nil {
			return nil, nil, serr
		}
		pub, perr := k.PublicKey.ECDH()
		if perr != nil {
			return nil, nil, perr
		}
		d = vapidDoc{Sealed: sealed, Public: base64.RawURLEncoding.EncodeToString(pub.Bytes())}
		if putDoc(ctx, s.KV, store.Item{PK: pkVAPID, SK: "-"}, d, store.Cond{NotExists: true}) != nil {
			// Another instance kept one first: that one is the key.
			if _, err := getDoc(ctx, s.KV, pkVAPID, "-", &d); err != nil {
				return nil, nil, err
			}
		}
	} else if err != nil {
		return nil, nil, err
	}
	der, err := s.Sealer.Open(ctx, d.Sealed, vapidBinding)
	if err != nil {
		return nil, nil, err
	}
	k, err := x509.ParseECPrivateKey(der)
	if err != nil || k.Curve != elliptic.P256() {
		return nil, nil, errors.New("planque's VAPID key doesn't read")
	}
	pub, err := k.PublicKey.ECDH()
	if err != nil {
		return nil, nil, err
	}
	// The public key kept beside it must be this key's, or the extension pinned another.
	if base64.RawURLEncoding.EncodeToString(pub.Bytes()) != d.Public {
		return nil, nil, errors.New("planque's VAPID public key isn't its private key's")
	}
	return k, pub.Bytes(), nil
}

// chromeWaker is the waker wakeChrome uses: the one set on the server (tests),
// else planque's own, built once from its VAPID key when the relay runs and a
// VAPID subject is set. Nil: no wake, a Chrome that isn't connected is offline.
func (s *Server) chromeWaker(ctx context.Context) *ChromeWaker {
	if s.ChromeWaker != nil {
		return s.ChromeWaker
	}
	if s.VAPIDSubject == "" || !s.relayOn() {
		return nil
	}
	s.vapidMu.Lock()
	defer s.vapidMu.Unlock()
	if s.vapidWaker == nil {
		k, pub, err := s.vapidKey(ctx)
		if err != nil {
			s.log().ErrorContext(ctx, "chrome wake: no VAPID key", slog.String("err", err.Error()))
			return nil
		}
		s.vapidWaker = &ChromeWaker{Signer: stepkey.Local{Key: k}, Public: pub, Subject: s.VAPIDSubject}
	}
	return s.vapidWaker
}

// pushKey answers GET /v1/extension/push-key: planque's VAPID public key, for
// the extension's build to pin (config.js pushKey). Public by nature; made on
// the first ask. 404 while the wake isn't offered.
func (s *Server) pushKey(w http.ResponseWriter, r *http.Request) {
	wk := s.chromeWaker(r.Context())
	if wk == nil {
		writeErr(w, http.StatusNotFound, "planque doesn't wake the extension here")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"key": base64.RawURLEncoding.EncodeToString(wk.Public)})
}
