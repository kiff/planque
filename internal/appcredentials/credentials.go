// Package appcredentials supplies registered OAuth client credentials to
// the token transport. References are operator configuration, never grant
// arguments, and credential values never enter approval or export JSON.
package appcredentials

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

// Credentials is used only while authenticating the registered client.
type Credentials struct {
	ClientID string `json:"client_id"`
	Secret   string `json:"-"`
}

// String names the client and hides its secret, so no log or %v/%+v of a
// Credentials can print it.
func (c Credentials) String() string {
	if c.Secret == "" {
		return "client " + c.ClientID + " (no secret)"
	}
	return "client " + c.ClientID + " (secret hidden)"
}

// GoString is String for %#v.
func (c Credentials) GoString() string { return c.String() }

// Reader resolves an operator-configured app id to its credentials.
type Reader interface {
	Read(context.Context, string) (Credentials, error)
}

type JSONAPI interface {
	Do(context.Context, string, any, any) error
}

var appID = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

var ErrUnavailable = errors.New("registered app credentials are unavailable")

// SecretsManager reads only the references configured by the operator: an
// app id's own reference, or, with Prefix, the secret named Prefix + id
// (planque/apps/googlesheets), which the process's role may read under
// that prefix alone. Either way the id must be a configured app's, by its
// pattern, so a grant cannot turn an id into an arbitrary secret-store
// lookup, and a new app needs only its definition and its secret.
type SecretsManager struct {
	API        JSONAPI
	References map[string]string
	Prefix     string
}

func (s *SecretsManager) Read(ctx context.Context, id string) (Credentials, error) {
	ref, ok := s.References[id]
	if !ok && s.Prefix != "" {
		ref, ok = s.Prefix+id, true
	}
	if !appID.MatchString(id) || !ok || strings.TrimSpace(ref) == "" || s.API == nil {
		return Credentials{}, ErrUnavailable
	}
	var out struct {
		SecretString string `json:"SecretString"`
		SecretBinary []byte `json:"SecretBinary"`
	}
	if err := s.API.Do(ctx, "GetSecretValue", map[string]string{"SecretId": ref, "VersionStage": "AWSCURRENT"}, &out); err != nil {
		// Do not return the backing service's text or configured reference
		// through a grant response: this error can reach an agent.
		return Credentials{}, ErrUnavailable
	}
	raw := out.SecretBinary
	if out.SecretString != "" {
		raw = []byte(out.SecretString)
	}
	defer clear(raw)
	return decode(raw)
}

// Sealed keeps application credentials using the same KV and envelope
// encryption interfaces as grants. It also works in a self-hosted process.
// Put is an operator API in Go, deliberately absent from the HTTP handler.
type Sealed struct {
	KV     store.KV
	Sealer seal.Sealer
}

func binding(id string) map[string]string {
	return map[string]string{"kind": "oauth_application", "application": id}
}

func (s *Sealed) Read(ctx context.Context, id string) (Credentials, error) {
	if !appID.MatchString(id) || s.KV == nil || s.Sealer == nil {
		return Credentials{}, ErrUnavailable
	}
	it, ok, err := s.KV.Get(ctx, "OAUTHAPP#"+id, "CREDENTIALS")
	if err != nil || !ok {
		return Credentials{}, ErrUnavailable
	}
	var rec struct {
		Sealed []byte `json:"sealed"`
	}
	if json.Unmarshal(it.Doc, &rec) != nil {
		return Credentials{}, ErrUnavailable
	}
	raw, err := s.Sealer.Open(ctx, rec.Sealed, binding(id))
	if err != nil {
		return Credentials{}, ErrUnavailable
	}
	defer clear(raw)
	return decode(raw)
}

func (s *Sealed) Put(ctx context.Context, id string, credentials Credentials) error {
	if !appID.MatchString(id) || s.KV == nil || s.Sealer == nil {
		return ErrUnavailable
	}
	raw, err := json.Marshal(struct {
		ClientID string `json:"client_id"`
		Secret   string `json:"client_secret"`
	}{credentials.ClientID, credentials.Secret})
	if err != nil {
		return ErrUnavailable
	}
	defer clear(raw)
	if _, err := decode(raw); err != nil {
		return err
	}
	cipher, err := s.Sealer.Seal(ctx, raw, binding(id))
	if err != nil {
		return ErrUnavailable
	}
	doc, _ := json.Marshal(map[string]any{"sealed": cipher})
	if err := s.KV.Put(ctx, store.Item{PK: "OAUTHAPP#" + id, SK: "CREDENTIALS", Doc: doc}, store.Cond{}); err != nil {
		return ErrUnavailable
	}
	return nil
}

func decode(raw []byte) (Credentials, error) {
	if len(raw) == 0 || len(raw) > 64<<10 {
		return Credentials{}, ErrUnavailable
	}
	var wire struct {
		ClientID string `json:"client_id"`
		Secret   string `json:"client_secret"`
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&wire) != nil || strings.TrimSpace(wire.ClientID) == "" || len(wire.ClientID) > 4096 || len(wire.Secret) > 32768 || strings.ContainsAny(wire.ClientID+wire.Secret, "\r\n\x00") {
		return Credentials{}, ErrUnavailable
	}
	if d.Decode(new(any)) != io.EOF {
		return Credentials{}, ErrUnavailable
	}
	return Credentials{ClientID: wire.ClientID, Secret: wire.Secret}, nil
}
