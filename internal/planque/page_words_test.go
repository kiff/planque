package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/words"
)

// doIn is do, with the app's Accept-Language.
func (e *env) doIn(lang, method, path string, body any, phone *fakes.Phone) (int, map[string]any) {
	e.t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, e.url+path, bytes.NewReader(b))
	phone.Sign(req, tenant, b)
	req.Header.Set("Accept-Language", lang)
	client := e.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	raw, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(raw, &out)
	return resp.StatusCode, out
}

// browse opens link in a browser that asks for lang ("" asks for none).
func (e *env) browse(link, lang string) (int, string) {
	e.t.Helper()
	req, _ := http.NewRequest("GET", link, nil)
	if lang != "" {
		req.Header.Set("Accept-Language", lang)
	}
	client := e.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(raw)
}

func says(page, lang, line string) bool {
	return strings.Contains(page, `<html lang="`+lang+`"`) && strings.Contains(page, template.HTMLEscapeString(line))
}

// A sign-in's result page speaks the language the app approved in, whatever
// the browser's; with no grant to go by, the browser's; else English, as
// before.
func TestSignInPagesSpeakTheAppsLanguageThenTheBrowsers(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	st, out := e.doIn("de-DE, en;q=0.5", "POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone)
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Lang != "de" {
		t.Fatalf("the app's language wasn't kept: %q", g.Lang)
	}
	link, _ := url.Parse(out["sign_in_url"].(string))
	state := url.QueryEscape(link.Query().Get("state"))
	de := words.For("de")
	st, page := e.browse(e.url+"/oauth/callback?state="+state+"&error=access_denied", "fr-FR,fr;q=0.9")
	if st != 200 || !says(page, "de", de.NotConnected) || !says(page, "de", de.NotCompleted(hostOf(e.vendor.URL))) {
		t.Fatalf("the vendor's refusal, in German: %d %s", st, page)
	}
	// The link is used: no grant to go by, so the browser's language.
	st, page = e.browse(e.url+"/oauth/callback?state="+state+"&code=x", "fr-FR,fr;q=0.9")
	if fr := words.For("fr"); st != 400 || !says(page, "fr", fr.LinkUnknown) {
		t.Fatalf("a used link, in French: %d %s", st, page)
	}
	if st, page = e.browse(e.url+"/oauth/callback", "pl"); st != 400 || !says(page, "pl", words.For("pl").MissingState) {
		t.Fatalf("no state, in Polish: %d %s", st, page)
	}
	st, page = e.browse(e.url+"/oauth/callback?state=nope", "")
	if st != 400 || !says(page, "en", "This sign-in link is unknown, expired or already used. Approve again on your phone.") || !strings.Contains(page, "<h1>Not connected</h1>") {
		t.Fatalf("English, unchanged: %d %s", st, page)
	}
}

// The sign-in view speaks the language the app opened it in: its page, its
// buttons, its script's words and its answers; an unknown view, the browser's.
func TestTheSignInViewSpeaksTheAppsLanguage(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	st, out := e.doIn("nl-NL", "POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/sign-in", map[string]any{"url": fakes.ShopLogin}, e.phone)
	if st != 200 {
		t.Fatalf("start: %d %v", st, out)
	}
	viewURL := out["view"].(string)
	nl := words.For("nl")
	st, page := e.browse(viewURL, "it")
	if st != 200 || !says(page, "nl", nl.ViewNote("shop.example", int(signInViewTTL.Minutes()))) || !strings.Contains(page, "<title>"+nl.ViewTitle("shop.example")+"</title>") ||
		!strings.Contains(page, ">"+nl.Done+"</button>") || !strings.Contains(page, `"`+nl.DidntWork+`"`) || strings.Contains(page, "Sign in to shop.example yourself") {
		t.Fatalf("the view, in Dutch: %d %s", st, page)
	}
	if st, _, raw := e.view("POST", viewURL+"/done", nil); st != 200 || !strings.Contains(string(raw), nl.DoneWords) {
		t.Fatalf("done, in Dutch: %d %s", st, raw)
	}
	// Done forgets the view's token, and with it its grant: the browser's
	// language.
	es := words.For("es")
	if st, page := e.browse(viewURL, "es-ES"); st != 404 || !says(page, "es", es.ViewGone) || !says(page, "es", es.SignInView) {
		t.Fatalf("the finished view, in the browser's Spanish: %d %s", st, page)
	}
	if st, page := e.browse(e.url+"/browser/view/"+strings.Repeat("x", 129), "es-ES"); st != 404 || !says(page, "es", es.ViewUnknown) {
		t.Fatalf("a token that can't be one, in Spanish: %d %s", st, page)
	}
}
