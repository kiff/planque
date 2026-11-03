package composio

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/connector"
)

type JSONAPI interface {
	Do(context.Context, string, any, any) error
}

// SecretFromStore resolves one operator-selected reference. A tool cannot
// select a secret, and neither service errors nor reference names escape.
// The short cache avoids a secret-store round trip for each preflight GET.
func SecretFromStore(api JSONAPI, reference string) Secret {
	var mu sync.Mutex
	var key string
	var until time.Time
	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		if time.Now().Before(until) {
			return key, nil
		}
		if api == nil || strings.TrimSpace(reference) == "" {
			return "", connector.ErrUnavailable
		}
		var out struct {
			Text   string `json:"SecretString"`
			Binary []byte `json:"SecretBinary"`
		}
		if err := api.Do(ctx, "GetSecretValue", map[string]string{"SecretId": reference, "VersionStage": "AWSCURRENT"}, &out); err != nil {
			return "", connector.ErrUnavailable
		}
		raw := out.Binary
		if out.Text != "" {
			raw = []byte(out.Text)
		}
		defer clear(raw)
		var v struct {
			Key string `json:"api_key"`
		}
		d := json.NewDecoder(bytes.NewReader(raw))
		d.DisallowUnknownFields()
		if len(raw) > 64<<10 || d.Decode(&v) != nil || d.Decode(new(any)) != io.EOF || strings.TrimSpace(v.Key) == "" || strings.ContainsAny(v.Key, "\r\n\x00") {
			return "", connector.ErrUnavailable
		}
		key, until = v.Key, time.Now().Add(time.Minute)
		return key, nil
	}
}
