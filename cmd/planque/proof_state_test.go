package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiff/planque/internal/secretstore"
)

// oneSecret answers GetSecretValue with raw, as planque's one secret.
type oneSecret struct{ raw string }

func (o oneSecret) Do(_ context.Context, _ string, _ any, out any) error {
	b, _ := json.Marshal(map[string]string{"SecretString": o.raw})
	return json.Unmarshal(b, out)
}

// KIFF's proof state takes its token from planque's one secret first, then
// from PLANQUE_PROOF_STATE_TOKEN; half a setting is refused, so a token is
// never left unused nor an address asked without one.
func TestTheProofStateTokenComesFromTheOneSecretFirst(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Header.Get("Authorization"))
		_ = json.NewEncoder(w).Encode(map[string]string{"entry": r.URL.Query().Get("entry"), "state": "green", "since": "2026-11-17T10:00:00Z"})
	}))
	defer srv.Close()
	holding := secretstore.NewBundle(oneSecret{`{"proof_state_token":"from-secret"}`}, "planque")
	lacking := secretstore.NewBundle(oneSecret{`{"door_secret":"d"}`}, "planque")
	for name, c := range map[string]struct {
		token  string
		bundle *secretstore.Bundle
		sent   string
	}{
		"the secret holds it":            {"", holding, "Bearer from-secret"},
		"the secret over the env":        {"from-env", holding, "Bearer from-secret"},
		"the env while the secret lacks": {"from-env", lacking, "Bearer from-env"},
		"the env, no secret":             {"from-env", nil, "Bearer from-env"},
	} {
		got = nil
		p, err := proofStateFrom(srv.URL, c.token, c.bundle)
		if err != nil || p == nil {
			t.Fatalf("%s: %v %v", name, p, err)
		}
		if g := p.Greens(t.Context(), []string{"zendesk-key@1"}); len(g) != 1 || len(got) != 1 || got[0] != c.sent {
			t.Errorf("%s: sent %v, want %q", name, got, c.sent)
		}
	}
	// Lacking everywhere: the address is asked with no token, never.
	got = nil
	p, _ := proofStateFrom(srv.URL, "", lacking)
	if g := p.Greens(t.Context(), []string{"zendesk-key@1"}); len(g) != 0 || len(got) != 0 {
		t.Errorf("no token anywhere: %v sent, %d proved", got, len(g))
	}
	if p, err := proofStateFrom("", "", holding); p != nil || err != nil {
		t.Errorf("nothing set: %v %v", p, err)
	}
	for name, c := range map[string][2]string{
		"an address, no token, no secret": {srv.URL, ""},
		"a token, no address":             {"", "from-env"},
		"an address that isn't https":     {"http://kiff.example/state", "from-env"},
	} {
		if p, err := proofStateFrom(c[0], c[1], nil); p != nil || err == nil {
			t.Errorf("%s: taken", name)
		}
	}
}
