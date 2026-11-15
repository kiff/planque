package planque

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/words"
)

// A sign-in link that came back once, landed on again (the vendor's "click
// here" after its redirect, a reopened tab), says what happened: the
// sign-in already went through while the grant is signed in or connected,
// so the owner goes back to tafs. The link still works once: the second
// visit moves nothing. The owner on Android 104 saw "unknown, expired or
// already used. Approve again", and went round in circles.
func TestAUsedSignInLinkSaysTheSignInWentThrough(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	st, out := e.doIn("es-ES", "POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone)
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	link := out["sign_in_url"].(string)
	u, _ := url.Parse(link)
	again := e.url + "/oauth/callback?state=" + url.QueryEscape(u.Query().Get("state")) + "&code=demo-code"
	code := e.signIn(link)

	before, _ := e.s.loadGrant(context.Background(), tenant, id)
	if before.Status != StatusSignedIn {
		t.Fatalf("after the sign-in: %s", before.Status)
	}
	es := words.For("es")
	st, page := e.browse(again, "en")
	if st != 200 || !says(page, "es", es.SignedIn) || !says(page, "es", es.SignInWentThrough(e.s.appName(before))) {
		t.Fatalf("signed in, the link again, in the app's Spanish: %d %s", st, page)
	}
	after, _ := e.s.loadGrant(context.Background(), tenant, id)
	if after.Status != before.Status || len(after.Events) != len(before.Events) || after.ConfirmHash != before.ConfirmHash {
		t.Fatalf("the second visit moved the grant: %s, %d events", after.Status, len(after.Events))
	}
	for _, leak := range []string{id, tenant, code, u.Query().Get("state")} {
		if strings.Contains(page, leak) {
			t.Fatalf("the page shows %q: %s", leak, page)
		}
	}

	// Confirmed on the phone: connected, the same words.
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	connected, _ := e.s.loadGrant(context.Background(), tenant, id)
	if st, page = e.browse(again, "en"); st != 200 || !says(page, "es", es.SignInWentThrough(e.s.appName(connected))) {
		t.Fatalf("connected, the link again: %d %s", st, page)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Status != StatusConnected || len(g.Events) != len(connected.Events) {
		t.Fatalf("the visit after connecting moved the grant: %s", g.Status)
	}

	// A link planque never gave says what it said before.
	if st, page = e.browse(e.url+"/oauth/callback?state=never-given&code=x", ""); st != 400 || !says(page, "en", words.For("en").LinkUnknown) {
		t.Fatalf("an unknown link: %d %s", st, page)
	}
}

// The used link is kept by the state's hash, never the state itself, and
// only until the link's own deadline.
func TestAUsedSignInLinkIsKeptByItsHash(t *testing.T) {
	state := "a-state-planque-gave"
	if k := pkSignInUsed(state); strings.Contains(k, state) || !strings.HasPrefix(k, "SIGNINUSED#") || len(k) != len("SIGNINUSED#")+64 {
		t.Fatalf("the used link's key: %q", k)
	}
}
