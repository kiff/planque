package secretstore

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"
)

// Bundle is planque's one secret (PLANQUE_SECRET): a JSON object of named
// values, each a non-empty one-line string, and "apps", each registered
// application's credentials by its id:
//
//	{"composio_api_key": "…", "door_secret": "…", "kiff_event_secret": "…",
//	 "apps": {"googlesheets": {"client_id": "…", "client_secret": "…"}}}
//
// It's read at most once a minute. While a part isn't in it, or it can't be
// read, each part is read from its own secret as before (Or, App), so moving
// the values into it never takes one away; the old secrets go only once it's
// proven live. A nil Bundle holds nothing.
type Bundle struct {
	api       JSONAPI
	reference string

	mu    sync.Mutex
	doc   bundleDoc
	until time.Time
}

type bundleDoc struct {
	values map[string]string
	apps   map[string]json.RawMessage
}

// NewBundle reads the secret at reference; an empty reference is no bundle.
func NewBundle(api JSONAPI, reference string) *Bundle {
	if api == nil || strings.TrimSpace(reference) == "" {
		return nil
	}
	return &Bundle{api: api, reference: reference}
}

func (b *Bundle) read(ctx context.Context) (bundleDoc, bool) {
	if b == nil {
		return bundleDoc{}, false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if time.Now().Before(b.until) {
		return b.doc, true
	}
	var out struct {
		Text   string `json:"SecretString"`
		Binary []byte `json:"SecretBinary"`
	}
	if err := b.api.Do(ctx, "GetSecretValue", map[string]string{"SecretId": b.reference, "VersionStage": "AWSCURRENT"}, &out); err != nil {
		// Not cached: the next read tries again, and this one falls back.
		return bundleDoc{}, false
	}
	raw := out.Binary
	if out.Text != "" {
		raw = []byte(out.Text)
	}
	defer clear(raw)
	var fields map[string]json.RawMessage
	d := json.NewDecoder(bytes.NewReader(raw))
	if len(raw) > 64<<10 || d.Decode(&fields) != nil || d.Decode(new(any)) != io.EOF {
		return bundleDoc{}, false
	}
	doc := bundleDoc{values: map[string]string{}}
	for name, v := range fields {
		if name == "apps" {
			if json.Unmarshal(v, &doc.apps) != nil {
				return bundleDoc{}, false
			}
			continue
		}
		var s string
		// A value that isn't a one-line string is left out: its part is
		// read from its own secret, as if the bundle lacked it.
		if json.Unmarshal(v, &s) == nil && strings.TrimSpace(s) != "" && !strings.ContainsAny(s, "\r\n\x00") {
			doc.values[name] = s
		}
	}
	b.doc, b.until = doc, time.Now().Add(time.Minute)
	return doc, true
}

// Or reads the value name from the bundle, else through fallback (its own
// secret, as before); with no fallback, ErrUnavailable.
func (b *Bundle) Or(name string, fallback func(context.Context) (string, error)) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		if doc, ok := b.read(ctx); ok {
			if v, ok := doc.values[name]; ok {
				return v, nil
			}
		}
		if fallback == nil {
			return "", ErrUnavailable
		}
		return fallback(ctx)
	}
}

// App is the application id's credentials as the bundle holds them, still
// encoded; false when it holds none for id (or can't be read).
func (b *Bundle) App(ctx context.Context, id string) ([]byte, bool) {
	doc, ok := b.read(ctx)
	if !ok {
		return nil, false
	}
	raw, ok := doc.apps[id]
	if !ok {
		return nil, false
	}
	return bytes.Clone(raw), true
}
