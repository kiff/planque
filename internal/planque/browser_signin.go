package planque

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/store"
)

// The live sign-in view (kiff/planque#86, step 4): the owner signs in to
// a site themselves, on their phone, in a view of the grant's own browser
// that planque serves. The phone asks for the view with a signed request;
// planque starts a session on the worker with the grant's kept sign-ins
// and answers a one-time URL the app opens. The page shows the browser
// as frames and sends the owner's taps, typing and keys to it; nothing
// of the page is read for an agent, and no frame is kept. When the owner
// taps Done, the session's cookies are sealed as the grant's sign-in. A
// view not finished in ten minutes lapses, and nothing is kept from it.
// While a view is open, the taf's calls on the grant are refused.

// signInViewTTL is how long a sign-in view is good for.
const signInViewTTL = 10 * time.Minute

// CodeSigningIn: a call refused while the owner's sign-in view is open.
const CodeSigningIn = "signing_in"

// signInRecord is an open sign-in view, on the grant.
type signInRecord struct {
	TokenSHA256 string    `json:"token_sha256"`
	Session     string    `json:"session"`
	Site        string    `json:"site,omitempty"`
	Expires     time.Time `json:"expires"`
	OpenedAt    time.Time `json:"opened_at"`
}

// signInToken is the view's token, hashed, pointing at its grant.
type signInToken struct {
	Tenant  string    `json:"tenant"`
	Grant   string    `json:"grant"`
	Expires time.Time `json:"expires"`
}

func skBrowserSignIn(id string) string { return "BROWSER_SIGNIN#" + id }
func pkSignInToken(hash string) string { return "BROWSERVIEW#" + hash }

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// openSignIn is the grant's open sign-in view, or nil.
func (s *Server) openSignIn(ctx context.Context, g *Grant) (*signInRecord, error) {
	var rec signInRecord
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skBrowserSignIn(g.ID), &rec); errors.Is(err, errNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if !s.now().Before(rec.Expires) {
		return nil, nil
	}
	return &rec, nil
}

// dropSignIn forgets an open view and ends its session, best effort.
func (s *Server) dropSignIn(ctx context.Context, g *Grant) {
	var rec signInRecord
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skBrowserSignIn(g.ID), &rec); err != nil {
		return
	}
	if r, ok := s.Browser.(browser.Resumer); ok && rec.Session != "" {
		if sess, err := r.Resume(ctx, rec.Session); err == nil {
			_, _ = sess.Close(ctx)
		}
	}
	_, _, _ = s.KV.Delete(ctx, tenantPK(g.Tenant), skBrowserSignIn(g.ID))
	_, _, _ = s.KV.Delete(ctx, pkSignInToken(rec.TokenSHA256), "-")
}

// startBrowserSignIn is the phone's signed ask for a sign-in view on a
// browser grant: POST …/grants/{id}/sign-in, body {url?}. It answers the
// view's URL, good for ten minutes, which the app opens.
func (s *Server) startBrowserSignIn(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		URL string `json:"url"`
	}
	if len(sr.body) > 0 {
		if err := decodeStrict(sr.body, &req); err != nil {
			s.fail(w, r, err)
			return
		}
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Request.Kind != KindBrowser {
		writeErr(w, http.StatusConflict, "this connection has no browser to sign in with")
		return
	}
	if g.Status != StatusConnected {
		writeErr(w, http.StatusConflict, "this connection is "+g.Status+"; it takes no sign-in now")
		return
	}
	resumer, ok := s.Browser.(browser.Resumer)
	if s.Browser == nil || !ok {
		writeErr(w, http.StatusServiceUnavailable, "this installation's browser can't hold a sign-in view")
		return
	}
	open, err := s.openSignIn(r.Context(), g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if open != nil {
		writeErr(w, http.StatusConflict, "a sign-in view is open for this connection: finish it on the phone, or wait for it to lapse")
		return
	}
	sites := []string(nil)
	if g.Request.Browser != nil {
		sites = g.Request.Browser.Sites
	}
	start := strings.TrimSpace(req.URL)
	if start == "" && len(sites) > 0 {
		start = "https://" + sites[0] + "/"
	}
	if start != "" {
		if code, why := browser.Allowed(start, sites); code != "" {
			writeErr(w, http.StatusBadRequest, "url: "+why)
			return
		}
	}
	state, err := s.openBrowserState(r.Context(), g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	sess, err := resumer.(browser.Worker).Start(r.Context(), browser.Start{Grant: g.ID, State: state, Sites: sites})
	if err != nil {
		s.log().WarnContext(r.Context(), "sign-in view not started", slog.String("grant", g.ID), slog.String("err", err.Error()))
		writeErr(w, http.StatusServiceUnavailable, "the browser couldn't be started: "+err.Error())
		return
	}
	if start != "" {
		if err := sess.Input(r.Context(), browser.Input{Kind: "go", URL: start}); err != nil {
			_, _ = sess.Close(r.Context())
			writeErr(w, http.StatusBadGateway, "the page couldn't be opened: "+err.Error())
			return
		}
	}
	token := rand.Text()
	now := s.now()
	rec := signInRecord{TokenSHA256: tokenHash(token), Session: sess.ID(), Site: browser.Site(start), Expires: now.Add(signInViewTTL), OpenedAt: now.UTC()}
	ttl := rec.Expires.Add(time.Minute).Unix()
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skBrowserSignIn(g.ID), Ver: 1, TTL: ttl}, rec, store.Cond{}); err != nil {
		_, _ = sess.Close(r.Context())
		s.fail(w, r, err)
		return
	}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: pkSignInToken(rec.TokenSHA256), SK: "-", Ver: 1, TTL: ttl}, signInToken{Tenant: g.Tenant, Grant: g.ID, Expires: rec.Expires}, store.Cond{NotExists: true}); err != nil {
		_, _ = sess.Close(r.Context())
		s.fail(w, r, err)
		return
	}
	g.event(now, "a sign-in view was opened "+sr.approval.Who())
	_ = s.saveGrant(r.Context(), g)
	writeJSON(w, http.StatusOK, map[string]any{"view": s.BaseURL + "/browser/view/" + token, "expires": rec.Expires.UTC().Format(time.RFC3339),
		"words": "Sign in to the site yourself in the view. Nothing you type is seen by the taf. Tap Done when you're signed in."})
}

// viewOf is the open view a token names: its grant and session. why is
// the refusal in words when there is none.
func (s *Server) viewOf(ctx context.Context, token string) (g *Grant, sess browser.Session, rec *signInRecord, why string) {
	if token == "" || len(token) > 128 {
		return nil, nil, nil, "This sign-in view is unknown."
	}
	var t signInToken
	if _, err := getDoc(ctx, s.KV, pkSignInToken(tokenHash(token)), "-", &t); err != nil {
		return nil, nil, nil, "This sign-in view is unknown, finished or lapsed. Open it again from your phone."
	}
	if !s.now().Before(t.Expires) {
		return nil, nil, nil, "This sign-in view has lapsed (ten minutes). Open it again from your phone; nothing was kept."
	}
	g, err := s.loadGrant(ctx, t.Tenant, t.Grant)
	if err != nil {
		return nil, nil, nil, "This sign-in view's connection is gone."
	}
	open, err := s.openSignIn(ctx, g)
	if err != nil || open == nil || open.TokenSHA256 != tokenHash(token) {
		return nil, nil, nil, "This sign-in view is finished. Open it again from your phone if you need to."
	}
	resumer, ok := s.Browser.(browser.Resumer)
	if !ok {
		return nil, nil, nil, "This installation's browser can't hold a sign-in view."
	}
	sess, err = resumer.Resume(ctx, open.Session)
	if err != nil {
		return nil, nil, nil, "The browser behind this view has ended. Open it again from your phone; nothing was kept."
	}
	return g, sess, open, ""
}

// viewHeaders keep the view and its frames out of caches and referrers.
func viewHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// browserView serves the view page.
func (s *Server) browserView(w http.ResponseWriter, r *http.Request) {
	viewHeaders(w)
	token := r.PathValue("token")
	_, _, rec, why := s.viewOf(r.Context(), token)
	if why != "" {
		s.page(w, http.StatusNotFound, "Sign-in view", why)
		return
	}
	w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src blob:; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; form-action 'none'; base-uri 'none'")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	site := rec.Site
	if site == "" {
		site = "the site"
	}
	_ = viewTmpl.Execute(w, map[string]any{"Token": token, "Site": site, "Minutes": int(signInViewTTL.Minutes())})
}

// browserViewFrame is the browser as shown now, JPEG. Never kept.
func (s *Server) browserViewFrame(w http.ResponseWriter, r *http.Request) {
	viewHeaders(w)
	_, sess, _, why := s.viewOf(r.Context(), r.PathValue("token"))
	if why != "" {
		writeErr(w, http.StatusNotFound, why)
		return
	}
	jpeg, err := sess.Frame(r.Context())
	if err != nil {
		writeErr(w, http.StatusBadGateway, "no frame: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	_, _ = w.Write(jpeg)
}

// browserViewInput is the owner's act: a tap, text, a key, a scroll, or
// going to a URL, which must be on the grant's sites.
func (s *Server) browserViewInput(w http.ResponseWriter, r *http.Request) {
	viewHeaders(w)
	g, sess, _, why := s.viewOf(r.Context(), r.PathValue("token"))
	if why != "" {
		writeErr(w, http.StatusNotFound, why)
		return
	}
	var in browser.Input
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if err := decodeStrict(body, &in); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := in.Check(); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.Kind == "go" {
		sites := []string(nil)
		if g.Request.Browser != nil {
			sites = g.Request.Browser.Sites
		}
		if code, why := browser.Allowed(in.URL, sites); code != "" {
			writeErr(w, http.StatusBadRequest, why)
			return
		}
	}
	if err := sess.Input(r.Context(), in); err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// browserViewDone ends the view: the session's cookies become the
// grant's sealed sign-in, the view is forgotten, and the grant's record
// says so.
func (s *Server) browserViewDone(w http.ResponseWriter, r *http.Request) {
	viewHeaders(w)
	g, sess, rec, why := s.viewOf(r.Context(), r.PathValue("token"))
	if why != "" {
		writeErr(w, http.StatusNotFound, why)
		return
	}
	state, err := sess.Close(r.Context())
	_, _, _ = s.KV.Delete(r.Context(), tenantPK(g.Tenant), skBrowserSignIn(g.ID))
	_, _, _ = s.KV.Delete(r.Context(), pkSignInToken(rec.TokenSHA256), "-")
	if err != nil {
		writeErr(w, http.StatusBadGateway, "the browser didn't hand its sign-in back: "+err.Error())
		return
	}
	s.keepBrowserState(r.Context(), g, state)
	what := "the sign-in view was closed; what the site kept is sealed for the taf's calls"
	if rec.Site != "" {
		what = "signed in at " + rec.Site + " on the phone; what the site kept is sealed for the taf's calls"
	}
	g.event(s.now(), what)
	_ = s.saveGrant(r.Context(), g)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "words": "Done. The taf can use the site as you signed in. You can close this view."})
}

// viewTmpl is the view: the browser as frames, the owner's input sent
// back. Its words are planque's; the app only opens the URL.
var viewTmpl = template.Must(template.New("view").Parse(`<!doctype html><html><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1"><title>Sign in to {{.Site}}</title>
<style>
body{margin:0;font:15px system-ui,sans-serif;background:#111;color:#eee}
.bar{display:flex;gap:6px;padding:8px;flex-wrap:wrap;align-items:center;background:#222}
.bar input{flex:1;min-width:120px;padding:8px;font:inherit;border:1px solid #555;border-radius:6px;background:#000;color:#eee}
button{padding:8px 12px;font:inherit;border:0;border-radius:6px;background:#3a6;color:#fff}
button.k{background:#555}
#done{background:#06c}
#frame{display:block;width:100%;height:auto;background:#000;touch-action:manipulation}
#note{padding:8px;color:#bbb;font-size:13px}
</style></head><body>
<div class="bar"><input id="url" placeholder="https://{{.Site}}/…" autocapitalize="off" autocorrect="off"><button id="go">Go</button></div>
<img id="frame" alt="The browser">
<div class="bar"><input id="text" placeholder="Type here, then Send" autocapitalize="off" autocorrect="off"><button id="send">Send</button>
<button class="k" data-key="Enter">Enter</button><button class="k" data-key="Backspace">⌫</button><button class="k" data-key="Tab">Tab</button>
<button class="k" id="up">▲</button><button class="k" id="down">▼</button><button id="done">Done</button></div>
<p id="note">Sign in to {{.Site}} yourself: tap the browser where you'd tap, type below and Send. Nothing you type is seen by the taf. This view lapses after {{.Minutes}} minutes; tap Done when you're signed in.</p>
<script>
(() => {
  const base = location.pathname.replace(/\/$/, '');
  const img = document.getElementById('frame'), note = document.getElementById('note');
  const W = 1280, H = 900;
  let busy = false, stopped = false, lastY = H / 2;
  const say = t => { note.textContent = t; };
  const post = (path, body) => fetch(base + path, {method: 'POST', headers: body ? {'Content-Type': 'application/json'} : {}, body: body ? JSON.stringify(body) : undefined});
  const input = async (body) => {
    const r = await post('/input', body);
    if (!r.ok) { const e = await r.json().catch(() => ({})); say(e.error || 'That didn\'t work.'); }
    setTimeout(frame, 250);
  };
  async function frame() {
    if (busy || stopped) return;
    busy = true;
    try {
      const r = await post('/frame');
      if (r.status === 404) { stopped = true; say((await r.json()).error); return; }
      if (r.ok) { const b = await r.blob(); const u = URL.createObjectURL(b); img.onload = () => URL.revokeObjectURL(u); img.src = u; }
    } catch (e) {} finally { busy = false; }
  }
  setInterval(frame, 900); frame();
  img.addEventListener('click', e => {
    const r = img.getBoundingClientRect();
    const x = Math.round((e.clientX - r.left) * W / r.width), y = Math.round((e.clientY - r.top) * H / r.height);
    lastY = y; input({kind: 'tap', x, y});
  });
  document.getElementById('go').onclick = () => { const u = document.getElementById('url').value.trim(); if (u) input({kind: 'go', url: u}); };
  document.getElementById('send').onclick = () => { const t = document.getElementById('text'); if (t.value) { input({kind: 'type', text: t.value}); t.value = ''; } };
  for (const b of document.querySelectorAll('button[data-key]')) b.onclick = () => input({kind: 'key', key: b.dataset.key});
  document.getElementById('up').onclick = () => input({kind: 'scroll', x: W / 2, y: lastY, delta_y: -400});
  document.getElementById('down').onclick = () => input({kind: 'scroll', x: W / 2, y: lastY, delta_y: 400});
  document.getElementById('done').onclick = async () => {
    stopped = true;
    const r = await post('/done'); const j = await r.json().catch(() => ({}));
    say(j.words || j.error || 'Done.');
    for (const b of document.querySelectorAll('button')) b.disabled = true;
  };
})();
</script></body></html>`))

// signingIn says whether the owner's sign-in view is open on g, for a
// call to be refused meanwhile.
func (s *Server) signingIn(ctx context.Context, g *Grant) bool {
	rec, err := s.openSignIn(ctx, g)
	return err == nil && rec != nil
}
