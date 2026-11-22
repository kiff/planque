package words

import "testing"

// liveLanguages are the languages KIFF's legal pages are in, kiff/cloud's
// components.PageLanguages: the site's eight, then the rest. The owner reads
// planque's approvals in the language they read everything else in, so
// planque speaks each. A language added there is added here, with its file.
var liveLanguages = []string{"en", "de", "fr", "es", "it", "pt", "nl", "pl",
	"ko", "ja", "et", "lv", "lt", "el", "cs", "sk", "sl", "hu", "ro", "bg", "hr", "tr", "id", "vi", "th", "zh",
	"ar", "sr", "sq", "ka", "hy", "az", "bn", "he", "ca", "mk", "mn", "ne", "km", "lo", "si", "ta", "tg", "dv", "dz"}

// Every live language is one planque speaks, once, and picks: the tests that
// loop over Languages then hold every key in each.
func TestEveryLiveLanguageIsSpoken(t *testing.T) {
	have := map[string]int{}
	for _, l := range Languages {
		have[l]++
	}
	for l, n := range have {
		if n > 1 {
			t.Errorf("%s is in Languages %d times", l, n)
		}
	}
	for _, l := range liveLanguages {
		if have[l] == 0 {
			t.Errorf("%s is live but planque has no words for it", l)
		}
		if Pick("", l) != l {
			t.Errorf("%s: an app asking for it gets %s", l, Pick("", l))
		}
	}
	if len(have) != len(liveLanguages) {
		t.Errorf("planque speaks %d languages, %d are live", len(have), len(liveLanguages))
	}
}
