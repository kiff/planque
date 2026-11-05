package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
)

// startView is the phone's signed ask for a sign-in view.
func (e *browserEnv) startView(id string, body map[string]any) (int, map[string]any) {
	e.t.Helper()
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/sign-in", body, e.phone, "")
}

// view calls one of the view's routes with no credential but the token.
func (e *browserEnv) view(method, path string, body any) (int, http.Header, []byte) {
	e.t.Helper()
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, path, bytes.NewReader(b))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := e.client.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header, raw
}

// 1. The view's life: the phone asks, signed; the page is served to its
// token alone; frames come from the browser and are kept nowhere; the
// owner's typing reaches the browser and never an agent; Done seals the
// cookies as the grant's sign-in, which the next call starts with.
func TestBrowserSignInViewSealsWhatTheOwnerSignedIn(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	// Unsigned: refused before anything starts.
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/sign-in", nil, nil, "k_setup"); st != 401 || e.b.Sessions != 0 {
		t.Fatalf("unsigned ask: %d, sessions %d", st, e.b.Sessions)
	}
	st, out := e.startView(g.ID, map[string]any{"url": fakes.ShopLogin})
	if st != 200 || e.b.Sessions != 1 {
		t.Fatalf("start: %d %v (sessions %d)", st, out, e.b.Sessions)
	}
	viewURL, _ := out["view"].(string)
	if !strings.HasPrefix(viewURL, e.url+"/browser/view/") || out["expires"] == "" || !strings.Contains(out["words"].(string), "Nothing you type is seen by the taf") {
		t.Fatalf("the view: %v", out)
	}
	// The page, its words, and no site-specific thing in it.
	st, hdr, page := e.view("GET", viewURL, nil)
	if st != 200 || !strings.Contains(string(page), "Sign in to shop.example yourself") || hdr.Get("Cache-Control") != "no-store" || !strings.Contains(hdr.Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatalf("the page: %d %s %v", st, page[:min(len(page), 200)], hdr)
	}
	// A frame, from the browser, as a JPEG.
	st, hdr, frame := e.view("POST", viewURL+"/frame", nil)
	if st != 200 || hdr.Get("Content-Type") != "image/jpeg" || !strings.Contains(string(frame), "JPEG "+fakes.ShopLogin) {
		t.Fatalf("frame: %d %s %q", st, hdr.Get("Content-Type"), frame)
	}
	// While the view is open, the taf's calls on the grant are refused.
	if out, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != CodeSigningIn || e.b.Sessions != 1 {
		t.Fatalf("a call during the sign-in: %s (sessions %d)", out, e.b.Sessions)
	}
	// The owner's acts: a tap, a URL off the sites refused, typing the
	// password, Enter. The fake signs the shop in on that.
	for _, in := range []map[string]any{{"kind": "tap", "x": 100, "y": 200}, {"kind": "type", "text": "hunter2"}, {"kind": "key", "key": "Enter"}} {
		if st, _, raw := e.view("POST", viewURL+"/input", in); st != 200 {
			t.Fatalf("input %v: %d %s", in, st, raw)
		}
	}
	if st, _, raw := e.view("POST", viewURL+"/input", map[string]any{"kind": "go", "url": fakes.OtherSite}); st != 400 || !strings.Contains(string(raw), "sites this connection may use") {
		t.Fatalf("go off the sites: %d %s", st, raw)
	}
	if st, _, raw := e.view("POST", viewURL+"/input", map[string]any{"kind": "key", "key": "F12"}); st != 400 {
		t.Fatalf("a key that isn't one: %d %s", st, raw)
	}
	if n := len(e.b.Inputs); n != 4 { // go (the start), tap, type, key; not the refused ones
		t.Fatalf("inputs reached the browser: %d %+v", n, e.b.Inputs)
	}
	// The password never appears in anything an agent reads: no page was
	// read for one. The grant's record names no secret either.
	if raw := mustJSONString(e.b.Inputs); !strings.Contains(raw, "hunter2") {
		t.Fatalf("the typing didn't reach the browser: %s", raw)
	}
	// Done: the cookies are sealed as the grant's sign-in.
	st, _, raw := e.view("POST", viewURL+"/done", nil)
	if st != 200 || !strings.Contains(string(raw), "The taf can use the site") {
		t.Fatalf("done: %d %s", st, raw)
	}
	it, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skBrowserState(g.ID))
	if !ok || strings.Contains(string(it.Doc), "signed-in") {
		t.Fatalf("the state at rest: %v %s", ok, it.Doc)
	}
	// The view is finished: its token opens nothing more.
	if st, _, _ := e.view("GET", viewURL, nil); st != 404 {
		t.Errorf("the finished view: %d", st)
	}
	if st, _, _ := e.view("POST", viewURL+"/frame", nil); st != 404 {
		t.Errorf("a frame after done: %d", st)
	}
	// The next call starts with the sign-in.
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("a call after the sign-in: %v", sc)
	}
	if last := e.b.States[len(e.b.States)-1]; !strings.Contains(string(last), "signed-in") {
		t.Errorf("the call didn't start with the sign-in: %s", last)
	}
	gv, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	events := mustJSONString(gv.Events)
	if !strings.Contains(events, "sign-in view was opened") || !strings.Contains(events, "signed in at shop.example on the phone") || strings.Contains(events, "hunter2") {
		t.Errorf("the grant's record: %s", events)
	}
}

// 2. What the view refuses: an unknown token, a lapsed view (nothing
// kept), a start on a grant that isn't a browser, a start URL off the
// sites, a worker that can't hold a view; and a disconnect ends an open
// view.
func TestBrowserSignInViewRefusals(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if st, _, _ := e.view("GET", e.url+"/browser/view/not-a-token", nil); st != 404 {
		t.Errorf("an unknown token: %d", st)
	}
	if st, _, raw := e.view("POST", e.url+"/browser/view/not-a-token/frame", nil); st != 404 || !strings.Contains(string(raw), "unknown") {
		t.Errorf("an unknown token's frame: %d %s", st, raw)
	}
	if st, out := e.startView(g.ID, map[string]any{"url": fakes.OtherSite}); st != 400 || !strings.Contains(out["error"].(string), "sites this connection may use") {
		t.Errorf("a start off the sites: %d %v", st, out)
	}
	// A lapsed view: nothing is kept, and the token opens nothing.
	st, out := e.startView(g.ID, nil)
	if st != 200 {
		t.Fatalf("start: %d %v", st, out)
	}
	viewURL := out["view"].(string)
	later := e.s.now().Add(signInViewTTL + time.Second)
	e.s.Now = func() time.Time { return later }
	if st, _, raw := e.view("GET", viewURL, nil); st != 404 || !strings.Contains(string(raw), "lapsed") {
		t.Errorf("a lapsed view: %d %s", st, raw)
	}
	if st, _, _ := e.view("POST", viewURL+"/done", nil); st != 404 {
		t.Errorf("done on a lapsed view: %d", st)
	}
	if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skBrowserState(g.ID)); ok {
		t.Error("a lapsed view kept a state")
	}
	// The lapsed record stays until its TTL: the taf's calls still see
	// no open view. (The clock goes back so permits and signatures hold.)
	e.s.Now = nil
	if _, _, err := e.kv.Delete(context.Background(), tenantPK(tenant), skBrowserSignIn(g.ID)); err != nil {
		t.Fatal(err)
	}
	// The taf's calls run again once the view lapsed.
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Errorf("a call after the lapse: %v", sc)
	}
	// A new view, then a disconnect: the view's session ends with the grant.
	st, out = e.startView(g.ID, map[string]any{"url": fakes.ShopLogin})
	if st != 200 {
		t.Fatalf("start again: %d %v", st, out)
	}
	viewURL = out["view"].(string)
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if st, _, _ := e.view("POST", viewURL+"/frame", nil); st != 404 {
		t.Errorf("a frame after the disconnect: %d", st)
	}
	// A worker that can't resume a session holds no view.
	e.s.Browser = noResume{e.b}
	h := e.connect()
	if st, out := e.startView(h.ID, nil); st != 503 || !strings.Contains(out["error"].(string), "can't hold a sign-in view") {
		t.Errorf("a worker with no resume: %d %v", st, out)
	}
	// A grant that isn't a browser.
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/grt_none/sign-in", nil, e.phone, ""); st != 404 {
		t.Errorf("no such grant: %d %v", st, out)
	}
}

// noResume is a worker without Resume.
type noResume struct{ w browser.Worker }

func (n noResume) Start(ctx context.Context, s browser.Start) (browser.Session, error) {
	return n.w.Start(ctx, s)
}

// 3. A view closed in the app before Done doesn't lock anyone out: the
// phone's next ask ends it, keeping nothing, and opens a fresh one.
func TestBrowserSignInViewAnotherAskReplacesTheOpenOne(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	st, first := e.startView(g.ID, map[string]any{"url": fakes.ShopLogin})
	if st != 200 {
		t.Fatalf("first view: %d %v", st, first)
	}
	firstURL := first["view"].(string)
	if st, _, raw := e.view("POST", firstURL+"/input", map[string]any{"kind": "type", "text": "half-typed"}); st != 200 {
		t.Fatalf("input: %d %s", st, raw)
	}
	st, second := e.startView(g.ID, nil)
	if st != 200 || e.b.Sessions != 2 {
		t.Fatalf("second view: %d %v (sessions %d)", st, second, e.b.Sessions)
	}
	secondURL := second["view"].(string)
	if secondURL == firstURL {
		t.Fatal("the second view reused the first one's token")
	}
	// The first view is over: its token opens nothing, its browser ended.
	if st, _, _ := e.view("GET", firstURL, nil); st != 404 {
		t.Errorf("the replaced view's page: %d", st)
	}
	if st, _, _ := e.view("POST", firstURL+"/done", nil); st != 404 {
		t.Errorf("the replaced view's Done: %d", st)
	}
	if _, err := e.b.Resume(context.Background(), "fake_1"); err == nil {
		t.Error("the replaced view's browser is still running")
	}
	// Nothing was kept from it.
	if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skBrowserState(g.ID)); ok {
		t.Error("the replaced view's state was kept")
	}
	// The new one works to the end.
	if st, _, raw := e.view("POST", secondURL+"/done", nil); st != 200 {
		t.Fatalf("done on the new view: %d %s", st, raw)
	}
	gv, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if events := mustJSONString(gv.Events); !strings.Contains(events, "ended for a new one") || strings.Contains(events, "half-typed") {
		t.Errorf("the grant's record: %s", events)
	}
}
