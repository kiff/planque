package planque

import (
	"crypto/sha256"
	"strings"
	"testing"
)

// The nightly connect matrix reads whether the managed sign-in can sign
// people in to an app: {entry, state} only, with its own key and nothing
// else's, and never a sign-in made, a set-up started or the provider's ids.
func TestTheSignInStateIsReadWithItsOwnKeyOnly(t *testing.T) {
	e, f := newConnectorFixture(t)
	path := func(entry string) string { return "/v1/catalog/" + entry + "/sign-in" }
	// Without a reader key set, the route isn't there.
	if st, _ := e.do("GET", path("connector:composio:desk_demo"), nil, nil, "reader-key"); st != 404 {
		t.Fatalf("no reader set: %d", st)
	}
	sum := sha256.Sum256([]byte("reader-key"))
	e.s.SignInReader = sum[:]
	for _, key := range []string{"", "wrong", "k_setup"} {
		if st, _ := e.do("GET", path("connector:composio:desk_demo"), nil, nil, key); st != 401 {
			t.Fatalf("key %q: %d", key, st)
		}
	}
	for entry, want := range map[string]string{
		"connector:composio:desk_demo":  "ready",
		"connector:composio:own_demo":   "needs_setup",
		"connector:composio:ondemand_1": "on_the_owners_pick",
	} {
		st, out := e.do("GET", path(entry), nil, nil, "reader-key")
		if st != 200 || out["entry"] != entry || out["state"] != want || len(out) != 2 {
			t.Fatalf("%s: %d %v", entry, st, out)
		}
		if b := mustJSONString(out); strings.Contains(b, "ac_") || strings.Contains(b, "auth_config") {
			t.Fatalf("a provider id went out: %s", b)
		}
	}
	// Reading an app on the owner's pick sets nothing up.
	if len(f.setUps) != 0 {
		t.Fatalf("a read set an app up: %v", f.setUps)
	}
	// Not the managed sign-in's: a reviewed entry, another provider, no slug.
	for _, entry := range []string{"qonto-key@2", "connector:other:desk_demo", "connector:composio:", "connector:composio:a%2Fb"} {
		if st, _ := e.do("GET", path(entry), nil, nil, "reader-key"); st != 404 {
			t.Fatalf("%s: %d", entry, st)
		}
	}
	// The reader key reads nothing else.
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "reader-key"); st != 401 && st != 403 {
		t.Fatalf("the reader key read the catalog: %d", st)
	}
}
