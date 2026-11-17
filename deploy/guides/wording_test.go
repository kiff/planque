package guides

import (
	"regexp"
	"strings"
	"testing"
)

// harmlessWords are the Eval's "never says read-only" patterns
// (kiff/agents cmd/e2e/planwords.go, readOnlyWords): Coach reads these
// guides and copies their words, so no guide may call a connection, a
// sign-in or a step read-only or harmless (Eval 43975057022, 17 Nov: "I'm
// connecting Google Sheets with read-only access to cells", from this
// directory's googlesheets.md). Keep them in step with the gate.
var harmlessWords = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bread-only\b|\bread only\b(\s+\pL+)?`),
	regexp.MustCompile(`(?i)\b(changes?|modif(y|ies)|touch(es)?) nothing\b`),
	regexp.MustCompile(`(?i)\b(won'?t|will not|doesn'?t|does not|can'?t|cannot|can not) (change|modify|touch) anything\b`),
	regexp.MustCompile(`(?i)\blecture seule\b`),
	regexp.MustCompile(`(?i)\bne (change|modifie|touche)(ra)? rien\b`),
	regexp.MustCompile(`(?i)\bne peut rien (changer|modifier)\b`),
	regexp.MustCompile(`(?i)\baucune modification\b`),
}

// jobWords are the words after "read only" that make it what the taf reads
// ("read only the rows added since its last run"), as the gate's readsOnly.
var jobWords = map[string]bool{"the": true, "a": true, "an": true, "its": true, "your": true, "my": true, "their": true,
	"this": true, "that": true, "these": true, "those": true, "new": true, "each": true, "every": true,
	"rows": true, "lines": true, "emails": true, "messages": true, "pages": true}

// harmlessClaim is the first phrase in text the Eval's gate would fail.
func harmlessClaim(text string) string {
	for _, re := range harmlessWords {
		for _, m := range re.FindAllStringSubmatch(text, -1) {
			if len(m) > 1 && m[1] != "" {
				if jobWords[strings.ToLower(strings.TrimSpace(m[1]))] {
					continue
				}
				return strings.TrimSuffix(m[0], m[1])
			}
			return m[0]
		}
	}
	return ""
}

// No served guide calls a connection, a sign-in or a step read-only or
// harmless: each says what the sign-in asks for, and that the Card limits
// the taf.
func TestNoGuideCallsAConnectionReadOnly(t *testing.T) {
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
			if got := harmlessClaim(line); got != "" {
				t.Errorf("%s.md:%d says %q: say what the sign-in asks for, and that the Card limits the taf", slug, i+1, got)
			}
		}
	}
	if served < 30 {
		t.Errorf("only %d guides were checked", served)
	}
}

// The check catches what the gate fails, and passes what the taf reads.
func TestHarmlessClaimMatchesTheGate(t *testing.T) {
	for _, bad := range []string{"read-only Sheets", "a read only sign-in", "Read means it changes nothing in Asana.",
		"it can't change anything", "accès en lecture seule", "ça ne changera rien", "read only what the taf needs"} {
		if harmlessClaim(bad) == "" {
			t.Errorf("%q passes", bad)
		}
	}
	for _, good := range []string{"it will read only the rows added since its last run", "`SUPABASE_RUN_READ_ONLY_QUERY`",
		"read just what changed", "if you decline, nothing changes", "full Drive, even for a job that writes nothing"} {
		if got := harmlessClaim(good); got != "" {
			t.Errorf("%q fails on %q", good, got)
		}
	}
}
