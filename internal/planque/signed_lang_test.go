package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/words"
)

// The vector both apps pin (kiff/tafs#709): the same request, language and
// words give the same digest_v 2, byte for byte.
func TestTheSharedDigestV2Vector(t *testing.T) {
	request := json.RawMessage(`{"agent":"books","asked_by":"setup-agent","created_at":"2026-11-11T16:00:00Z","grant":"g_02",` +
		`"server":"https://mcp.qonto.com/mcp","tenant":"t_01","tools":[{"read_only":true,"tool":"list_invoices"}],"v":"planque.grant.v1"}`)
	const text = "books pourra lire vos factures sur mcp.qonto.com"
	const want = "f1ca3ab371ecf97a82d134204e4931ef39ea474b7cd6a3e233b36f97499fbb32"
	if got, err := digestV2(request, "fr", text); err != nil || got != want {
		t.Fatalf("digest_v 2: %s %v, want %s", got, err, want)
	}
	// One word, or the language, changed: another digest.
	for _, c := range []struct{ lang, text string }{{"fr", strings.Replace(text, "lire", "modifier", 1)}, {"en", text}} {
		if got, _ := digestV2(request, c.lang, c.text); got == want {
			t.Errorf("%s %q has the vector's digest", c.lang, c.text)
		}
	}
}

// browserGrant asks for a browser grant, in e's language, and reads it as
// the owner's app lists it.
func browserGrant(t *testing.T, e *browserEnv) (*Grant, map[string]any) {
	t.Helper()
	st, out := e.ask(map[string]any{"agent_name": "Shop", "browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}})
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, out["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	st, list := e.do("GET", "/v1/accounts/"+tenant+"/grants", nil, e.phone, "")
	if st != 200 {
		t.Fatalf("list: %d %v", st, list)
	}
	for _, v := range list["grants"].([]any) {
		if v := v.(map[string]any); v["id"] == g.ID {
			return g, v
		}
	}
	t.Fatalf("grant %s isn't listed", g.ID)
	return nil, nil
}

// approveWith approves g on the phone with sum as its digest.
func approveWith(e *browserEnv, g *Grant, sum string) (int, map[string]any) {
	hash, _ := profile.Digest(g.Request.Profile)
	approved := ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}
	return e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]any{"digest": sum, "kiff_profile": approved}, e.phone, "")
}

// A French phone reads the grant in French, signs those very words, and
// planque keeps them, with their language, as the receipt.
func TestAGrantIsShownAndSignedInItsPersonsLanguage(t *testing.T) {
	e := newBrowserEnv(t)
	e.s.DigestV2 = true
	e.lang = "fr-FR,fr;q=0.9,en;q=0.8"
	g, v := browserGrant(t, e)

	fr := words.Approve("fr")
	if v["lang"] != "fr" || v["text"] != g.Texts["fr"] || !strings.Contains(g.Texts["fr"], strings.SplitN(fmt.Sprintf(fr.TheseSitesF, "Shop", "shop.example"), "\n", 2)[0]) {
		t.Fatalf("the view: lang %v, text %q", v["lang"], v["text"])
	}
	if v["status_words"] != fr.WaitingForYou {
		t.Errorf("status words %q, want %q", v["status_words"], fr.WaitingForYou)
	}
	want, _ := g.digestIn("fr")
	if v["digest_v"] != float64(2) || v["digest"] != want || want == g.Digest {
		t.Fatalf("digest_v %v digest %v, want 2 %s (v1 %s)", v["digest_v"], v["digest"], want, g.Digest)
	}
	// The digest covers the request the view sent, as the phone recomputes it.
	request, _ := json.Marshal(v["request"])
	if sum, _ := digestV2(request, "fr", v["text"].(string)); sum != want {
		t.Fatalf("the view's own request, words and language give %s, not %s", sum, want)
	}

	// Words the phone didn't show, or another language's: nothing is approved.
	request, _ = json.Marshal(g.Request)
	changed, _ := digestV2(request, "fr", strings.Replace(g.Texts["fr"], "Shop", "Shoq", 1))
	relabeled, _ := digestV2(request, "de", g.Texts["fr"])
	for _, sum := range []string{changed, relabeled} {
		if st, out := approveWith(e, g, sum); st != 409 {
			t.Fatalf("a digest over other words was taken: %d %v", st, out)
		}
	}

	if st, out := approveWith(e, g, want); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	if a := g.Approval; a == nil || a.Lang != "fr" || a.Text != g.Texts["fr"] {
		t.Fatalf("the receipt: %+v", g.Approval)
	}
}

// Without the switch, the view keeps digest v1 and no digest_v, as both
// apps read it today; the words are in the person's language all the same.
// And a language planque doesn't speak reads in English.
func TestDigestV1UntilTheSwitchAndEnglishForAnUnknownLanguage(t *testing.T) {
	e := newBrowserEnv(t)
	e.lang = "sv-SE"
	g, v := browserGrant(t, e)
	if v["lang"] != "en" || v["text"] != g.Text || v["digest"] != g.Digest || v["digest_v"] != nil {
		t.Fatalf("the view: lang %v digest %v digest_v %v", v["lang"], v["digest"], v["digest_v"])
	}
	if st, out := approveWith(e, g, g.Digest); st != 200 {
		t.Fatalf("approve with v1: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, g.ID)
	if g.Approval == nil || g.Approval.Lang != "" || g.Approval.Text != "" {
		t.Fatalf("a v1 signature covers no words: %+v", g.Approval)
	}
}

// A grant made before it wrote its words in each language reads, and is
// signed, in English: the only words it has.
func TestAnOlderGrantIsShownInEnglish(t *testing.T) {
	g := &Grant{Text: "Connect your account", Request: Request{Agent: "a"}}
	if lang, text := g.shown("fr"); lang != "en" || text != g.Text {
		t.Fatalf("shown: %s %q", lang, text)
	}
	en, _ := g.digestIn("en")
	if lang, text, ok := g.signed(en); !ok || lang != "en" || text != g.Text {
		t.Fatalf("signed: %s %q %v", lang, text, ok)
	}
	if _, _, ok := g.signed(""); ok {
		t.Fatal("an empty digest was taken")
	}
}

// Every builder writes every language: no grant's text is left empty or
// in English where another language is asked.
func TestEveryLanguagesTextIsWritten(t *testing.T) {
	e := newBrowserEnv(t)
	g, _ := browserGrant(t, e)
	for _, lang := range words.Languages {
		text := g.Texts[lang]
		if text == "" || (lang != "en" && text == g.Text) {
			t.Errorf("%s: %q", lang, text)
		}
	}
}
