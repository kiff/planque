package appcredentials

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

type fakeAPI struct{ calls int }

func (f *fakeAPI) Do(_ context.Context, action string, in, out any) error {
	f.calls++
	if action != "GetSecretValue" || in.(map[string]string)["SecretId"] != "configured-app-secret" || in.(map[string]string)["VersionStage"] != "AWSCURRENT" {
		return errors.New("unexpected secret lookup")
	}
	return json.Unmarshal([]byte(`{"SecretString":"{\"client_id\":\"demo-client\",\"client_secret\":\"demo-secret\"}"}`), out)
}

func TestSecretReferencesCannotBeChosenByAGrant(t *testing.T) {
	api := &fakeAPI{}
	s := &SecretsManager{API: api, References: map[string]string{"workspace": "configured-app-secret"}}
	c, err := s.Read(context.Background(), "workspace")
	if err != nil || c.ClientID != "demo-client" || c.Secret != "demo-secret" {
		t.Fatalf("credential lookup failed: %v", err)
	}
	for _, id := range []string{"another-secret", "../workspace", "", "arn:aws:secretsmanager:example"} {
		if _, err := s.Read(context.Background(), id); !errors.Is(err, ErrUnavailable) {
			t.Fatal("unconfigured app lookup was accepted")
		}
	}
	if api.calls != 1 {
		t.Fatal("a request selected an arbitrary backing-store secret")
	}
	b, _ := json.Marshal(c)
	if strings.Contains(string(b), "demo-secret") {
		t.Fatal("a credential was included in JSON output")
	}
}

func TestApplicationCredentialsAreSealedAndBoundToTheirApp(t *testing.T) {
	kv := &store.Memory{}
	s := &Sealed{KV: kv, Sealer: &seal.Local{Key: make([]byte, 32)}}
	ctx := context.Background()
	want := Credentials{ClientID: "demo-client", Secret: "demo-secret"}
	if err := s.Put(ctx, "workspace", want); err != nil {
		t.Fatal(err)
	}
	it, _, err := kv.Get(ctx, "OAUTHAPP#workspace", "CREDENTIALS")
	if err != nil || strings.Contains(string(it.Doc), want.Secret) || strings.Contains(string(it.Doc), want.ClientID) {
		t.Fatal("application credentials were stored in plaintext")
	}
	got, err := s.Read(ctx, "workspace")
	if err != nil || got != want {
		t.Fatalf("sealed credentials did not round-trip: %v", err)
	}
	it.PK = "OAUTHAPP#another"
	if err := kv.Put(ctx, it, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Read(ctx, "another"); !errors.Is(err, ErrUnavailable) {
		t.Fatal("one app could open another app's credentials")
	}
}
