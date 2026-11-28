package planque

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// healthEnv is the browser env with shop.example's fixed page reviewed, and
// a connection on it.
func healthEnv(t *testing.T) (*browserEnv, *Grant) {
	t.Helper()
	e := newBrowserEnv(t)
	e.s.HealthPages = map[string]HealthPage{"shop.example": {URL: "https://shop.example/account", Shows: "Signed in as"}}
	return e, e.connect()
}

// signIn seals a sign-in for shop.example that ends at ends.
func (e *browserEnv) signIn(g *Grant, ends time.Time) {
	e.t.Helper()
	state := fmt.Sprintf(`{"cookies":[{"name":"sid","value":"signed-in","domain":".shop.example","path":"/","expires":%d},{"name":"pref","value":"x","domain":"shop.example","path":"/","expires":-1}]}`, ends.Unix())
	e.s.keepBrowserState(context.Background(), g, []byte(state), nil)
}

func (e *browserEnv) health(g *Grant, method string) (int, map[string]any) {
	e.t.Helper()
	return e.do(method, "/v1/accounts/"+tenant+"/grants/"+g.ID+"/health", nil, nil, "k_setup")
}

func siteOf(t *testing.T, out map[string]any) map[string]any {
	t.Helper()
	sites, _ := out["sites"].([]any)
	if len(sites) != 1 {
		t.Fatalf("sites: %v", out)
	}
	return sites[0].(map[string]any)
}

// A connection's sign-in is checked on demand with the reviewed page: signed
// in, with when it was checked and when the site set the sign-in to end (the
// latest expiry of its cookies, never a name or value); the cookies the
// visit rotated are sealed as a session's are; and it runs again only after
// a day. GET reads what's kept and never runs anything.
func TestASignInIsCheckedAtMostOnceADay(t *testing.T) {
	e, g := healthEnv(t)
	ends := time.Now().Add(20 * 24 * time.Hour).Truncate(time.Second).UTC()
	e.signIn(g, ends)
	st, out := e.health(g, "POST")
	site := siteOf(t, out)
	if st != 200 || site["site"] != "shop.example" || site["verdict"] != browser.HealthSignedIn || site["checked_at"] == nil || site["expires_at"] != ends.Format(time.RFC3339) {
		t.Fatalf("check: %d %v", st, out)
	}
	if len(e.b.Checks) != 1 || e.b.Checks[0].URL != "https://shop.example/account" || e.b.Checks[0].Shows != "Signed in as" || e.b.Checks[0].Site != "shop.example" {
		t.Fatalf("the worker was asked: %+v", e.b.Checks)
	}
	if state, _, _ := e.s.openBrowserState(context.Background(), g); !strings.Contains(string(state), "signed-in-rotated") {
		t.Fatalf("the rotated sign-in wasn't sealed: %s", state)
	}
	for _, v := range []string{"signed-in", "sid", "pref"} {
		if strings.Contains(fmt.Sprint(out), v) {
			t.Fatalf("the answer carries %q: %v", v, out)
		}
	}
	// Within the day: the kept verdict, nothing run. GET never runs.
	if _, out := e.health(g, "POST"); siteOf(t, out)["verdict"] != browser.HealthSignedIn || len(e.b.Checks) != 1 {
		t.Fatalf("ran again within the day: %d checks", len(e.b.Checks))
	}
	later := time.Now().Add(25 * time.Hour)
	e.s.Now = func() time.Time { return later }
	defer func() { e.s.Now = nil }()
	if _, out := e.health(g, "GET"); siteOf(t, out)["verdict"] != browser.HealthSignedIn || len(e.b.Checks) != 1 {
		t.Fatalf("GET ran a check: %d", len(e.b.Checks))
	}
	if e.health(g, "POST"); len(e.b.Checks) != 2 {
		t.Fatalf("after a day: %d checks", len(e.b.Checks))
	}
}

// Without a saved sign-in there's nothing to open: signed out, and the
// worker isn't asked. A site with no reviewed page has no check. A check
// that can't run keeps "unknown" and waits six hours.
func TestASignInCheckWithoutASignInOrAPageOrAWorker(t *testing.T) {
	e, g := healthEnv(t)
	if _, out := e.health(g, "POST"); siteOf(t, out)["verdict"] != browser.HealthSignedOut || len(e.b.Checks) != 0 {
		t.Fatalf("no sign-in: %v, %d checks", out, len(e.b.Checks))
	}

	e2, g2 := healthEnv(t)
	e2.s.HealthPages = nil
	e2.signIn(g2, time.Now().Add(time.Hour))
	if _, out := e2.health(g2, "POST"); siteOf(t, out)["verdict"] != browser.HealthUnknown || siteOf(t, out)["why"] == nil || len(e2.b.Checks) != 0 {
		t.Fatalf("no reviewed page: %v", out)
	}

	e3, g3 := healthEnv(t)
	e3.signIn(g3, time.Now().Add(time.Hour))
	e3.b.HealthDown = true
	if _, out := e3.health(g3, "POST"); siteOf(t, out)["verdict"] != browser.HealthUnknown {
		t.Fatalf("worker down: %v", out)
	}
	e3.b.HealthDown = false
	if e3.health(g3, "POST"); len(e3.b.Checks) != 0 {
		t.Fatal("ran again before the backoff")
	}
	later := time.Now().Add(7 * time.Hour)
	e3.s.Now = func() time.Time { return later }
	defer func() { e3.s.Now = nil }()
	if _, out := e3.health(g3, "POST"); siteOf(t, out)["verdict"] != browser.HealthSignedIn || len(e3.b.Checks) != 1 {
		t.Fatalf("after the backoff: %v", out)
	}
}

// Only with the account's key, only for a hosted-browser connection, and
// only on an account the hosted browser is offered to.
func TestASignInCheckIsTheAccountsOwnOnAnOfferedBrowser(t *testing.T) {
	e, g := healthEnv(t)
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/health", nil, nil, ""); st == 200 {
		t.Fatal("checked with no key")
	}
	e.s.BrowserOffer.All, e.s.BrowserOffer.Tenants = false, nil
	if st, _ := e.health(g, "POST"); st != 403 || len(e.b.Checks) != 0 {
		t.Fatalf("not offered: %d", st)
	}
}

// The sign-in's end is the latest expiry of the site's own cookies: a
// parent domain's counts, another site's and session cookies don't.
func TestASignInEndsWhenItsSiteSetIt(t *testing.T) {
	state := `{"cookies":[{"domain":".shop.example","expires":2000},{"domain":"shop.example","expires":3000},{"domain":"other.example","expires":9000},{"domain":"shop.example","expires":-1}]}`
	if end, ok := cookieEnd([]byte(state), "shop.example"); !ok || end.Unix() != 3000 {
		t.Fatalf("end: %v %v", end, ok)
	}
	if end, ok := cookieEnd([]byte(state), "www.shop.example"); !ok || end.Unix() != 3000 {
		t.Fatalf("a subdomain: %v %v", end, ok)
	}
	if _, ok := cookieEnd([]byte(`{"cookies":[{"domain":"shop.example","expires":-1}]}`), "shop.example"); ok {
		t.Fatal("session cookies have no end")
	}
}

// The reviewed list reads, and its review rule is in its header.
func TestTheReviewedHealthPagesRead(t *testing.T) {
	_ = ReviewedHealthPages()
	if !strings.Contains(string(healthPagesJSON), "opening the URL must change nothing") {
		t.Fatal("the review rule isn't in health_pages.json")
	}
	for site, page := range ReviewedHealthPages() {
		if !strings.HasPrefix(page.URL, "https://"+site+"/") || strings.TrimSpace(page.Shows) == "" {
			t.Fatalf("%s's page %+v isn't an https page on it with its signed-in words", site, page)
		}
	}
}

// A connection's sign-in checks go with it.
func TestASignInChecksVerdictGoesWithItsConnection(t *testing.T) {
	e, g := healthEnv(t)
	e.signIn(g, time.Now().Add(time.Hour))
	e.health(g, "POST")
	if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "HEALTH#"+g.ID+"#"); len(items) != 1 {
		t.Fatalf("kept: %d", len(items))
	}
	e.s.dropBrowser(context.Background(), g)
	if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "HEALTH#"+g.ID+"#"); len(items) != 0 {
		t.Fatalf("left after the connection went: %d", len(items))
	}
}
