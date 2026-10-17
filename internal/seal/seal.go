// Package seal encrypts what planque keeps secret: vendor tokens, the
// KIFF connector key, PKCE verifiers. Hosted, each secret has its own data
// key from KMS (envelope encryption); self-hosted or in tests, a local
// AES-256 key does the same job.
//
// Every secret is sealed with a context naming the record it belongs to
// (tenant, grant, kind). The context is authenticated, so a ciphertext
// moved to another record does not open.
package seal

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/kiff/planque/internal/awsapi"
)

// Sealer seals and opens secrets bound to a context.
type Sealer interface {
	Seal(ctx context.Context, plaintext []byte, binding map[string]string) ([]byte, error)
	Open(ctx context.Context, sealed []byte, binding map[string]string) ([]byte, error)
}

// ErrOpen: the secret does not open with this context.
var ErrOpen = errors.New("seal: the secret does not open for this record")

type envelope struct {
	V     int    `json:"v"`
	Key   []byte `json:"k,omitempty"` // the KMS-encrypted data key
	Nonce []byte `json:"n"`
	Data  []byte `json:"c"`
}

func aad(binding map[string]string) []byte {
	keys := make([]string, 0, len(binding))
	for k := range binding {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k + "=" + binding[k] + "\n")
	}
	return []byte(b.String())
}

func gcm(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func encrypt(key, plaintext []byte, binding map[string]string) (nonce, data []byte, err error) {
	g, err := gcm(key)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, g.NonceSize())
	_, _ = rand.Read(nonce)
	return nonce, g.Seal(nil, nonce, plaintext, aad(binding)), nil
}

func decrypt(key []byte, e envelope, binding map[string]string) ([]byte, error) {
	g, err := gcm(key)
	if err != nil {
		return nil, err
	}
	out, err := g.Open(nil, e.Nonce, e.Data, aad(binding))
	if err != nil {
		return nil, ErrOpen
	}
	return out, nil
}

// Local seals with one AES-256 key.
type Local struct{ Key []byte }

// NewLocal reads a base64 32-byte key.
func NewLocal(b64 string) (*Local, error) {
	k, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
	if err != nil || len(k) != 32 {
		return nil, errors.New("seal: the local key must be 32 bytes, base64")
	}
	return &Local{Key: k}, nil
}

func (l *Local) Seal(_ context.Context, p []byte, b map[string]string) ([]byte, error) {
	n, d, err := encrypt(l.Key, p, b)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{V: 1, Nonce: n, Data: d})
}

func (l *Local) Open(_ context.Context, s []byte, b map[string]string) ([]byte, error) {
	var e envelope
	if json.Unmarshal(s, &e) != nil || e.V != 1 {
		return nil, ErrOpen
	}
	return decrypt(l.Key, e, b)
}

// KMS seals each secret with its own data key from AWS KMS. The binding
// is also the KMS encryption context, so KMS logs which record each
// decryption was for and refuses one asked for another record.
type KMS struct {
	KeyID string
	API   *awsapi.Client
}

func (k *KMS) Seal(ctx context.Context, p []byte, b map[string]string) ([]byte, error) {
	var out struct {
		Plaintext      []byte `json:"Plaintext"`
		CiphertextBlob []byte `json:"CiphertextBlob"`
	}
	err := k.API.Do(ctx, "GenerateDataKey", map[string]any{
		"KeyId": k.KeyID, "KeySpec": "AES_256", "EncryptionContext": b,
	}, &out)
	if err != nil {
		return nil, err
	}
	n, d, err := encrypt(out.Plaintext, p, b)
	clear(out.Plaintext)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{V: 1, Key: out.CiphertextBlob, Nonce: n, Data: d})
}

func (k *KMS) Open(ctx context.Context, s []byte, b map[string]string) ([]byte, error) {
	var e envelope
	if json.Unmarshal(s, &e) != nil || e.V != 1 || len(e.Key) == 0 {
		return nil, ErrOpen
	}
	var out struct {
		Plaintext []byte `json:"Plaintext"`
	}
	err := k.API.Do(ctx, "Decrypt", map[string]any{
		"KeyId": k.KeyID, "CiphertextBlob": e.Key, "EncryptionContext": b,
	}, &out)
	if awsapi.Is(err, "InvalidCiphertextException") {
		return nil, ErrOpen
	}
	if err != nil {
		return nil, err
	}
	defer clear(out.Plaintext)
	return decrypt(out.Plaintext, e, b)
}
