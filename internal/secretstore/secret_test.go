package secretstore

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeStore struct {
	raw   string
	calls int
}

func (f *fakeStore) Do(_ context.Context, operation string, in, out any) error {
	f.calls++
	if operation != "GetSecretValue" {
		panic("unexpected write")
	}
	raw, _ := json.Marshal(map[string]string{"SecretString": f.raw})
	return json.Unmarshal(raw, out)
}
func TestReaderKeepsProjectAndHandoffFieldsSeparate(t *testing.T) {
	for _, field := range []string{"api_key", "secret"} {
		store := &fakeStore{raw: `{"` + field + `":"synthetic-value"}`}
		read := Reader(store, "configured-reference", field)
		for range 2 {
			got, err := read(context.Background())
			if err != nil || got != "synthetic-value" {
				t.Fatal("configured field unavailable")
			}
		}
		if store.calls != 1 {
			t.Fatal("secret store cache was bypassed")
		}
		if _, err := Reader(store, "configured-reference", "another_field")(context.Background()); err == nil {
			t.Fatal("a different secret field was accepted")
		}
	}
}
