package secretstore

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// failingStore answers each GetSecretValue with err, else as fakeStore.
type failingStore struct {
	fakeStore
	err error
}

func (f *failingStore) Do(ctx context.Context, operation string, in, out any) error {
	if f.err != nil {
		f.calls++
		return f.err
	}
	return f.fakeStore.Do(ctx, operation, in, out)
}

func old(v string) func(context.Context) (string, error) {
	return func(context.Context) (string, error) { return v, nil }
}

// Each value comes from the one secret, read once a minute; a value it lacks,
// or holds as anything but a one-line string, comes from its own secret;
// with neither, it's unavailable.
func TestTheOneSecretFirstThenEachOwn(t *testing.T) {
	ctx := context.Background()
	store := &failingStore{fakeStore: fakeStore{raw: `{"composio_api_key":"from-bundle","door_secret":"","kiff_event_secret":7,"multi":"a\nb","apps":{"googlesheets":{"client_id":"c","client_secret":"s"}}}`}}
	b := NewBundle(store, "planque")
	if got, err := b.Or("composio_api_key", old("from-own"))(ctx); err != nil || got != "from-bundle" {
		t.Fatalf("composio: %q %v", got, err)
	}
	for _, name := range []string{"door_secret", "kiff_event_secret", "multi", "missing"} {
		if got, err := b.Or(name, old("from-own"))(ctx); err != nil || got != "from-own" {
			t.Errorf("%s: %q %v, want its own secret's", name, got, err)
		}
	}
	if _, err := b.Or("missing", nil)(ctx); !errors.Is(err, ErrUnavailable) {
		t.Errorf("neither: %v", err)
	}
	if store.calls != 1 {
		t.Errorf("the one secret was read %d times, want once", store.calls)
	}
	raw, ok := b.App(ctx, "googlesheets")
	var creds map[string]string
	if !ok || json.Unmarshal(raw, &creds) != nil || creds["client_id"] != "c" {
		t.Errorf("app: %s %v", raw, ok)
	}
	if _, ok := b.App(ctx, "slack"); ok {
		t.Error("an app the bundle doesn't hold")
	}
}

// A one secret that can't be read, or isn't a JSON object, holds nothing:
// each value comes from its own secret, and a failed read isn't kept, so the
// next one tries again. No bundle named is no bundle.
func TestAnUnreadableOneSecretFallsBackAndIsReadAgain(t *testing.T) {
	ctx := context.Background()
	store := &failingStore{err: errors.New("throttled")}
	b := NewBundle(store, "planque")
	if got, err := b.Or("door_secret", old("from-own"))(ctx); err != nil || got != "from-own" {
		t.Fatalf("unreadable: %q %v", got, err)
	}
	store.err, store.raw = nil, `{"door_secret":"from-bundle"}`
	if got, _ := b.Or("door_secret", old("from-own"))(ctx); got != "from-bundle" {
		t.Fatalf("the failed read was kept: %q", got)
	}
	for _, raw := range []string{`not json`, `["door_secret"]`, `{"door_secret":"a"} {}`, `{"apps":"x","door_secret":"a"}`} {
		nb := NewBundle(&fakeStore{raw: raw}, "planque")
		if got, _ := nb.Or("door_secret", old("from-own"))(ctx); got != "from-own" {
			t.Errorf("%s: %q", raw, got)
		}
	}
	if NewBundle(store, " ") != nil || NewBundle(nil, "planque") != nil {
		t.Error("a bundle with no reference or store")
	}
	var none *Bundle
	if got, _ := none.Or("door_secret", old("from-own"))(ctx); got != "from-own" {
		t.Errorf("no bundle: %q", got)
	}
}
