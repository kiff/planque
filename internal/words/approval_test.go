package words

import (
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// approvalVerb is a format's verb, with its argument when it names one.
var approvalVerb = regexp.MustCompile(`%(?:\[(\d)\])?[sd]`)

// englishLeft is an English word of planque's own, left in another language.
var englishLeft = regexp.MustCompile(`(?i)\b(permit|vault|sign-in|approve)\b`)

// args are the arguments a format takes, in the order it takes them, each
// once: "%[2]s … %[1]s" takes 1 and 2.
func args(format string) []string {
	var out []string
	n := 0
	for _, m := range approvalVerb.FindAllStringSubmatch(format, -1) {
		arg := m[1]
		if arg == "" {
			n++
			arg = string(rune('0' + n))
		}
		if !slices.Contains(out, arg) {
			out = append(out, arg)
		}
	}
	slices.Sort(out)
	return out
}

// Every language has every approval word, taking English's values, and none
// is left in English.
func TestEveryLanguageHasEveryApprovalWord(t *testing.T) {
	en := reflect.ValueOf(approvals["en"])
	for _, lang := range Languages {
		a, ok := approvals[lang]
		if !ok {
			t.Fatalf("%s: no approval words", lang)
		}
		v := reflect.ValueOf(a)
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Field(i).Name
			got, want := v.Field(i).String(), en.Field(i).String()
			if strings.TrimSpace(got) == "" {
				t.Errorf("%s %s: empty", lang, name)
			}
			if g, w := args(got), args(want); !reflect.DeepEqual(g, w) {
				t.Errorf("%s %s: takes %v, English takes %v", lang, name, g, w)
			}
			if strings.HasSuffix(name, "F") != (len(args(want)) > 0) {
				t.Errorf("%s: a format is named …F, and only a format", name)
			}
			if lang != "en" && got == want {
				t.Errorf("%s %s: left in English", lang, name)
			}
			// The same breaks, so a text reads in the same paragraphs.
			if strings.Count(got, "\n") != strings.Count(want, "\n") {
				t.Errorf("%s %s: %d line breaks, English has %d", lang, name, strings.Count(got, "\n"), strings.Count(want, "\n"))
			}
			if lang != "en" && englishLeft.MatchString(got) {
				t.Errorf("%s %s: English left in: %q", lang, name, got)
			}
		}
	}
	if len(approvals) != len(Languages) {
		t.Errorf("%d languages of approval words, %d languages", len(approvals), len(Languages))
	}
}

// taf and Card are names (kiff/tafs docs/languages.md): no language
// translates them.
func TestTafAndCardStayNames(t *testing.T) {
	for _, lang := range Languages {
		a := approvals[lang]
		if !strings.Contains(a.Changes, "Card") || !strings.Contains(a.ThisTafFirst, "taf") {
			t.Errorf("%s: %q, %q", lang, a.Changes, a.ThisTafFirst)
		}
	}
	if Approve("sv").Connected != "Connected" {
		t.Error("a language planque doesn't speak isn't English")
	}
}

// A managed app's approval card, in every language, names no internal
// service ("the connection service"), says why when the app asks for more
// than the taf uses, and
// names no tools (they come next, by their titles): the format takes the app
// alone.
func TestTheConnectCardNamesNoServiceAndSaysWhy(t *testing.T) {
	service := regexp.MustCompile(`(?i)connection service|Verbindungsdienst|service de connexion|servicio de conexión|servizio di collegamento|serviço de ligação|koppelingsdienst|usług[ai]? połączeń|usłudze połączeń`)
	why := map[string]string{"en": "because", "de": "weil", "fr": "car", "es": "porque", "it": "perché", "pt": "porque", "nl": "omdat", "pl": "bo"}
	for _, lang := range Languages {
		a := approvals[lang]
		for _, s := range []string{a.ConnectTitledF, a.ConnectReusedF, a.ReusedNext, a.EnterKeyF, a.EnterTokenF, a.SignInThenChoose, a.MoreAccessF, a.OnlyWhatItUsesF} {
			if service.MatchString(s) {
				t.Errorf("%s names the service: %q", lang, s)
			}
		}
		if !strings.Contains(a.MoreAccessF, " "+why[lang]+" ") {
			t.Errorf("%s doesn't say why: %q", lang, a.MoreAccessF)
		}
		for _, f := range []string{a.MoreAccessF, a.OnlyWhatItUsesF} {
			if got := args(f); !slices.Equal(got, []string{"1"}) {
				t.Errorf("%s takes %v, not the app alone: %q", lang, got, f)
			}
		}
	}
}
