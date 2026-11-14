package words

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var verb = regexp.MustCompile(`%[sd]`)

// Every language has every line, with English's values in English's order,
// and none is left in English (but for a few words the languages share).
func TestEveryLanguageHasEveryLine(t *testing.T) {
	en := reflect.ValueOf(catalog["en"])
	shared := map[string]bool{"Tab": true, "Enter": true}
	for _, lang := range Languages {
		l, ok := catalog[lang]
		if !ok {
			t.Fatalf("%s: no lines", lang)
		}
		v := reflect.ValueOf(l)
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Field(i).Name
			got, want := v.Field(i).String(), en.Field(i).String()
			if strings.TrimSpace(got) == "" {
				t.Errorf("%s %s: empty", lang, name)
			}
			if g, w := verb.FindAllString(got, -1), verb.FindAllString(want, -1); !reflect.DeepEqual(g, w) {
				t.Errorf("%s %s: values %v, English has %v", lang, name, g, w)
			}
			if lang != "en" && got == want && !shared[want] && !strings.HasPrefix(want, "%") {
				t.Errorf("%s %s: left in English", lang, name)
			}
		}
	}
	if len(catalog) != len(Languages) {
		t.Errorf("%d catalogs, %d languages", len(catalog), len(Languages))
	}
}

// The app's language wins; then the browser's, by its q (earlier on a tie);
// then English.
func TestPickTheAppsLanguageThenTheBrowsersThenEnglish(t *testing.T) {
	for _, c := range []struct{ stored, accept, want string }{
		{"de", "fr-FR,fr;q=0.9", "de"},
		{"", "fr-FR,fr;q=0.9,en;q=0.8", "fr"},
		{"sv", "sv-SE, pt-PT;q=0.8", "pt"},
		{"", "sv-SE,da;q=0.9", "en"},
		{"", "", "en"},
		{"", "de;q=0, nl;q=0.5", "nl"},
		{"PL", "", "pl"},
		{"", "pt_BR", "pt"},
		{"", "fr;q=0.5, de;q=0.8", "de"},
		{"", "nl;q=0.7, it;q=0.7", "nl"},
		{"", "de;q=2, fr;q=0.1", "fr"},
	} {
		if got := Pick(c.stored, c.accept); got != c.want {
			t.Errorf("Pick(%q, %q) = %q, want %q", c.stored, c.accept, got, c.want)
		}
	}
	if For("sv").NotConnected != "Not connected" {
		t.Error("a language planque doesn't speak isn't English")
	}
}

// Each language uses the glossary's word for a connection, where a line says it.
func TestTheGlossarysConnection(t *testing.T) {
	word := map[string]string{"de": "Verbindung", "fr": "connexion", "es": "conexión", "it": "collegamento", "pt": "ligação", "nl": "koppeling", "pl": "połączenie"}
	for lang, w := range word {
		if !strings.Contains(strings.ToLower(catalog[lang].NotCompletedRenewF), strings.ToLower(w)) {
			t.Errorf("%s: the connection isn't %q: %s", lang, w, catalog[lang].NotCompletedRenewF)
		}
	}
}

// Every language has every connection word, and a format verb where the
// English has one, so no row reads blank or with a stray %s.
func TestEveryLanguageHasEveryConnectionWord(t *testing.T) {
	en := connections["en"]
	for lang, c := range connections {
		for name, pair := range map[string][2]string{
			"SignInAgain": {c.SignInAgain, en.SignInAgain}, "Remove": {c.Remove, en.Remove}, "RemoveUnfinished": {c.RemoveUnfinished, en.RemoveUnfinished},
			"StopsF": {c.StopsF, en.StopsF}, "NothingStops": {c.NothingStops, en.NothingStops}, "WorkingLastF": {c.WorkingLastF, en.WorkingLastF},
			"WorkingNotYet": {c.WorkingNotYet, en.WorkingNotYet}, "Working": {c.Working, en.Working}, "ExpiredOnF": {c.ExpiredOnF, en.ExpiredOnF},
			"ExpiredF": {c.ExpiredF, en.ExpiredF}, "UnfinishedOnF": {c.UnfinishedOnF, en.UnfinishedOnF}, "Unfinished": {c.Unfinished, en.Unfinished},
			"DayMonthF": {c.DayMonthF, en.DayMonthF},
		} {
			if pair[0] == "" || strings.Count(pair[0], "%") != strings.Count(pair[1], "%") {
				t.Errorf("%s %s: %q", lang, name, pair[0])
			}
		}
		for i, m := range c.Months {
			if m == "" {
				t.Errorf("%s month %d", lang, i+1)
			}
		}
	}
	if len(connections) != len(approvals) {
		t.Errorf("connection words in %d languages, approval words in %d", len(connections), len(approvals))
	}
}
