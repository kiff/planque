package composio

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type secretAPI struct {
	calls int
	raw   string
	fail  bool
}

func (s *secretAPI) Do(_ context.Context, operation string, in, out any) error {
	s.calls++
	if operation != "GetSecretValue" || in.(map[string]string)["SecretId"] != "configured-reference" || in.(map[string]string)["VersionStage"] != "AWSCURRENT" {
		return errors.New("unexpected secret lookup")
	}
	if s.fail {
		return errors.New("sensitive backing-service response")
	}
	raw, _ := json.Marshal(map[string]string{"SecretString": s.raw})
	return json.Unmarshal(raw, out)
}
func TestSecretIsConfiguredCachedAndNeverExposesBackingErrors(t *testing.T) {
	api := &secretAPI{raw: `{"api_key":"synthetic-project-key"}`}
	read := SecretFromStore(api, "configured-reference")
	for i := 0; i < 2; i++ {
		key, err := read(context.Background())
		if err != nil || key != "synthetic-project-key" {
			t.Fatal("configured key was not read")
		}
	}
	if api.calls != 1 {
		t.Fatalf("secret was read %d times", api.calls)
	}
	api.fail = true
	_, err := SecretFromStore(api, "configured-reference")(context.Background())
	if err == nil || strings.Contains(err.Error(), "sensitive") {
		t.Fatal("backing service error escaped")
	}
	for _, raw := range []string{`{"api_key":"key","extra":true}`, `{"api_key":""}`, `{"api_key":"key"} {}`, `{"api_key":"bad\nkey"}`} {
		api.fail = false
		api.raw = raw
		if _, err := SecretFromStore(api, "configured-reference")(context.Background()); err == nil {
			t.Fatalf("malformed secret accepted: %s", raw)
		}
	}
}
