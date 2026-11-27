package appcredentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// A secret header's value is kept sealed with the client secret and hidden like it.
func TestASecretHeaderValueIsSealedAndHidden(t *testing.T) {
	kv := &store.Memory{}
	s := &Sealed{KV: kv, Sealer: &seal.Local{Key: make([]byte, 32)}}
	ctx := context.Background()
	want := Credentials{ClientID: "demo-client", Secret: "demo-secret", HeaderValue: "gate-value"}
	if err := s.Put(ctx, "gated", want); err != nil {
		t.Fatal(err)
	}
	it, _, _ := kv.Get(ctx, "OAUTHAPP#gated", "CREDENTIALS")
	if strings.Contains(string(it.Doc), "gate-value") {
		t.Fatal("the header's value was stored in plaintext")
	}
	got, err := s.Read(ctx, "gated")
	if err != nil || got != want {
		t.Fatalf("it did not round-trip: %v", err)
	}
	for _, shown := range []string{got.String(), fmt.Sprintf("%v %+v %#v", got, got, got)} {
		if strings.Contains(shown, "gate-value") || strings.Contains(shown, "demo-secret") {
			t.Fatalf("a credential printed its secret: %s", shown)
		}
	}
	if b, _ := json.Marshal(got); strings.Contains(string(b), "gate-value") {
		t.Fatalf("the header's value is in the JSON: %s", b)
	}
	if _, err := decode([]byte(`{"client_id":"c","header_secret":"a\nb"}`)); err == nil {
		t.Fatal("a header value with a line break was accepted")
	}
}

// namedAPI answers for any one secret it's told the name of, and records
// which names were asked.
type namedAPI struct{ asked []string }

func (f *namedAPI) Do(_ context.Context, action string, in, out any) error {
	id := in.(map[string]string)["SecretId"]
	f.asked = append(f.asked, id)
	if action != "GetSecretValue" || id != "planque/apps/googlesheets" {
		return errors.New("no such secret")
	}
	return json.Unmarshal([]byte(`{"SecretString":"{\"client_id\":\"sheets-client\",\"client_secret\":\"sheets-secret\"}"}`), out)
}

// With a prefix, an app's client is the secret named prefix + its id: a new
// app needs no configuration change. The id still has to look like an app
// id, so no grant can reach another name, and an app's own reference wins.
func TestAnAppsClientIsTheSecretNamedForIt(t *testing.T) {
	api := &namedAPI{}
	s := &SecretsManager{API: api, Prefix: "planque/apps/", References: map[string]string{"workspace": "configured-app-secret"}}
	c, err := s.Read(context.Background(), "googlesheets")
	if err != nil || c.ClientID != "sheets-client" || c.Secret != "sheets-secret" {
		t.Fatalf("by name: %v %v", c, err)
	}
	for _, id := range []string{"../workspace", "googlesheets/../x", "", "Google", "arn:aws:secretsmanager:example", "a/b"} {
		if _, err := s.Read(context.Background(), id); !errors.Is(err, ErrUnavailable) {
			t.Fatalf("%q was looked up", id)
		}
	}
	if _, err := s.Read(context.Background(), "nosuchapp"); !errors.Is(err, ErrUnavailable) {
		t.Fatal("a missing secret isn't unavailable")
	}
	if _, err := s.Read(context.Background(), "workspace"); !errors.Is(err, ErrUnavailable) {
		t.Fatal("the configured reference wasn't the one asked")
	}
	want := []string{"planque/apps/googlesheets", "planque/apps/nosuchapp", "configured-app-secret"}
	if strings.Join(api.asked, ",") != strings.Join(want, ",") {
		t.Fatalf("asked %v, want %v", api.asked, want)
	}
	// Without a prefix, only configured references, as before.
	plain := &SecretsManager{API: &namedAPI{}, References: map[string]string{}}
	if _, err := plain.Read(context.Background(), "googlesheets"); !errors.Is(err, ErrUnavailable) {
		t.Fatal("no prefix, yet an unconfigured app was looked up")
	}
}

// No way of printing Credentials prints the secret.
func TestCredentialsNeverPrintTheirSecret(t *testing.T) {
	c := Credentials{ClientID: "sheets-client", Secret: "sheets-secret"}
	for _, f := range []string{"%v", "%+v", "%#v", "%s", "%q"} {
		if got := fmt.Sprintf(f, c); strings.Contains(got, "sheets-secret") || !strings.Contains(got, "sheets-client") {
			t.Errorf("%s: %s", f, got)
		}
	}
	for _, f := range []string{"%v", "%+v", "%#v"} {
		if got := fmt.Sprintf(f, &c); strings.Contains(got, "sheets-secret") {
			t.Errorf("pointer %s: %s", f, got)
		}
		if got := fmt.Sprintf(f, struct{ C Credentials }{c}); strings.Contains(got, "sheets-secret") {
			t.Errorf("nested %s: %s", f, got)
		}
	}
}

type appBundle map[string]string

func (a appBundle) App(_ context.Context, id string) ([]byte, bool) {
	v, ok := a[id]
	return []byte(v), ok
}

type nextReader struct{ calls []string }

func (n *nextReader) Read(_ context.Context, id string) (Credentials, error) {
	n.calls = append(n.calls, id)
	return Credentials{ClientID: "own-" + id}, nil
}

// An app's credentials come from the one secret when it holds them, else
// from the reader before; a bad id reaches neither, and credentials in the
// bundle that don't decode are refused, not taken from elsewhere.
func TestBundledReadsTheOneSecretFirst(t *testing.T) {
	ctx := context.Background()
	next := &nextReader{}
	b := &Bundled{Bundle: appBundle{"googlesheets": `{"client_id":"c","client_secret":"s"}`, "broken": `{"client_id":""}`}, Next: next}
	if c, err := b.Read(ctx, "googlesheets"); err != nil || c.ClientID != "c" || c.Secret != "s" {
		t.Fatalf("from the bundle: %v %v", c, err)
	}
	if c, err := b.Read(ctx, "slack"); err != nil || c.ClientID != "own-slack" {
		t.Fatalf("from its own: %v %v", c, err)
	}
	if _, err := b.Read(ctx, "broken"); err == nil {
		t.Fatal("broken credentials in the bundle were accepted")
	}
	if _, err := b.Read(ctx, "../x"); err == nil {
		t.Fatal("a bad id was read")
	}
	if len(next.calls) != 1 || next.calls[0] != "slack" {
		t.Fatalf("reached the next reader: %v", next.calls)
	}
}
