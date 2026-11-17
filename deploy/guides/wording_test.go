package guides

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// testdata/coach-wording.json is a copy of kiff/tafs's
// contracts/coach-wording.json (at 3a50e11): the words Coach must never use
// about a connection, a sign-in or a step, which kiff/agents' show tool
// refuses on a card and the Eval's gate fails. Coach reads these guides and
// copies their words (Eval 43975057022, 17 Nov: "I'm connecting Google
// Sheets with read-only access to cells", from googlesheets.md), so no guide
// may use them either. Copy the file again whenever tafs changes it.
type wording struct {
	Patterns   []string `json:"patterns"`
	UnlessNext []string `json:"unless_next"`
	res        []*regexp.Regexp
}

func wordingRules(t *testing.T) []wording {
	t.Helper()
	raw, err := os.ReadFile("testdata/coach-wording.json")
	if err != nil {
		t.Fatal(err)
	}
	var rules []wording
	if err := json.Unmarshal(raw, &rules); err != nil {
		t.Fatal(err)
	}
	if len(rules) == 0 {
		t.Fatal("no wording rules")
	}
	for i := range rules {
		for _, p := range rules[i].Patterns {
			re, err := regexp.Compile(`(?i)(?:^|[^\pL\pN])(?:` + p + `)(?:$|[^\pL\pN])`)
			if err != nil {
				t.Fatalf("pattern %q: %v", p, err)
			}
			rules[i].res = append(rules[i].res, re)
		}
	}
	return rules
}

// claim is the first phrase in text one of the rules catches, the way
// kiff/agents' show.Claim reads them: each pattern on whole words without
// case, and a match whose captured next word is in unless_next ("read only
// the rows") let through.
func claim(text string, rules []wording) string {
	edge := func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }
	for _, w := range rules {
		for _, re := range w.res {
			for at := 0; at <= len(text); {
				loc := re.FindStringSubmatchIndex(text[at:])
				if loc == nil {
					break
				}
				group := ""
				if len(loc) > 3 && loc[2] >= 0 {
					group = text[at+loc[2] : at+loc[3]]
				}
				if next := strings.TrimSpace(group); next != "" && slicesHasFold(w.UnlessNext, next) {
					step := loc[1]
					if r, _ := utf8.DecodeRuneInString(group); edge(r) && loc[2] > 0 {
						step = loc[2]
					}
					at += max(step, 1)
					continue
				}
				phrase := text[at+loc[0] : at+loc[1]]
				if group != "" {
					phrase = strings.TrimSuffix(strings.TrimRightFunc(phrase, edge), group)
				}
				return strings.ToLower(strings.TrimFunc(phrase, edge))
			}
		}
	}
	return ""
}

func slicesHasFold(words []string, s string) bool {
	for _, w := range words {
		if strings.EqualFold(strings.TrimSpace(w), s) {
			return true
		}
	}
	return false
}

// No served guide calls a connection, a sign-in or a step read-only or
// harmless: each says what the sign-in asks for, and that the Card limits
// the taf.
func TestNoGuideCallsAConnectionReadOnly(t *testing.T) {
	rules := wordingRules(t)
	entries, err := files.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	served := 0
	for _, e := range entries {
		slug, ok := strings.CutSuffix(e.Name(), ".md")
		if !ok {
			continue
		}
		text, _, ok := Adapted(strings.ToLower(slug))
		if !ok {
			continue
		}
		served++
		for i, line := range strings.Split(text, "\n") {
			if got := claim(line, rules); got != "" {
				t.Errorf("%s.md:%d says %q: say what the sign-in asks for, and that the Card limits the taf", slug, i+1, got)
			}
		}
	}
	if served < 30 {
		t.Errorf("only %d guides were checked", served)
	}
}

// The copy catches what the gate fails, in English and French at least, and
// passes what the taf reads.
func TestTheWordingCopyCatchesTheGatesClaims(t *testing.T) {
	rules := wordingRules(t)
	for _, bad := range []string{"read-only Sheets", "a read only sign-in", "Read means it changes nothing in Asana.",
		"it can't change anything", "accès en lecture seule", "read only what the taf needs"} {
		if claim(bad, rules) == "" {
			t.Errorf("%q passes", bad)
		}
	}
	for _, good := range []string{"it will read only the rows added since its last run", "`SUPABASE_RUN_READ_ONLY_QUERY`",
		"read just what changed", "if you decline, nothing changes", "full Drive, even for a job that writes nothing"} {
		if got := claim(good, rules); got != "" {
			t.Errorf("%q fails on %q", good, got)
		}
	}
}
