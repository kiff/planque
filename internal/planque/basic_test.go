package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/words"
)

// The owner reads where a Basic key goes, with the username given beside it,
// in each language; the password given beside a key-as-username isn't said.
func TestTheApprovalSaysWhereABasicKeyGoes(t *testing.T) {
	q := Request{API: &APIRequest{URL: "https://vendor.example/openapi.yaml", Digest: strings.Repeat("a", 64), Base: "https://vendor.example/api"}}
	q.API.Auth = openapi.Auth{In: "basic", User: "me@example.com/token"}
	if got := apiKeyText(words.Approve("en"), q, false, ""); !strings.Contains(got, "as the password of an HTTP Basic sign-in, with the username me@example.com/token") {
		t.Fatalf("the key as the password: %s", got)
	}
	q.API.Auth = openapi.Auth{In: "basic", Password: "X"}
	if got := apiKeyText(words.Approve("en"), q, false, ""); !strings.Contains(got, "as the username of an HTTP Basic sign-in") {
		t.Fatalf("the key as the username: %s", got)
	}
	for _, lang := range []string{"de", "fr", "es", "it", "pt", "nl", "pl"} {
		w := words.Approve(lang)
		if w.KeyAsBasicUser == "" || !strings.Contains(w.KeyAsBasicPasswordF, "%s") || w.KeyAsBasicUser == words.Approve("en").KeyAsBasicUser {
			t.Errorf("%s: %q %q", lang, w.KeyAsBasicPasswordF, w.KeyAsBasicUser)
		}
	}
}
