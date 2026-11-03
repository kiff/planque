package secretstore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
	"time"
)

var ErrUnavailable = errors.New("the configured secret is unavailable")

type JSONAPI interface {
	Do(context.Context, string, any, any) error
}

// Reader resolves one operator-selected reference and JSON field. A tool cannot
// select a secret, and neither service errors nor reference names escape.
// The short cache avoids a secret-store round trip for each preflight GET.
func Reader(api JSONAPI, reference, field string) func(context.Context) (string, error) {
	var mu sync.Mutex
	var key string
	var until time.Time
	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		if time.Now().Before(until) {
			return key, nil
		}
		if api == nil || strings.TrimSpace(reference) == "" || field == "" {
			return "", ErrUnavailable
		}
		var out struct {
			Text   string `json:"SecretString"`
			Binary []byte `json:"SecretBinary"`
		}
		if err := api.Do(ctx, "GetSecretValue", map[string]string{"SecretId": reference, "VersionStage": "AWSCURRENT"}, &out); err != nil {
			return "", ErrUnavailable
		}
		raw := out.Binary
		if out.Text != "" {
			raw = []byte(out.Text)
		}
		defer clear(raw)
		var v map[string]string
		d := json.NewDecoder(bytes.NewReader(raw))
		d.DisallowUnknownFields()
		if len(raw) > 64<<10 || d.Decode(&v) != nil || len(v) != 1 || d.Decode(new(any)) != io.EOF || strings.TrimSpace(v[field]) == "" || strings.ContainsAny(v[field], "\r\n\x00") {
			return "", ErrUnavailable
		}
		key, until = v[field], time.Now().Add(time.Minute)
		return key, nil
	}
}
