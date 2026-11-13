// Package guides holds the adapted app guides beside it, built into
// planque so the server serves the ones in this tree (GET
// /v1/accounts/{tenant}/catalog/{entry}/guide).
package guides

import (
	"embed"
	"regexp"
	"strings"
)

//go:embed *.md
var files embed.FS

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// Adapted is an app's adapted guide, by its slug (gmail, googlesheets,
// stripe), and the day it was checked; ok is false when the app has none
// yet. This directory's own notes (README, apps) are never a guide.
func Adapted(slug string) (text, checked string, ok bool) {
	if !slugPattern.MatchString(slug) || slug == "readme" || slug == "apps" {
		return "", "", false
	}
	b, err := files.ReadFile(slug + ".md")
	if err != nil {
		return "", "", false
	}
	text = string(b)
	for _, line := range strings.SplitN(text, "\n", 8) {
		if rest, found := strings.CutPrefix(strings.TrimSpace(line), "Checked: "); found {
			checked, _, _ = strings.Cut(rest, " ")
			break
		}
	}
	return text, checked, true
}
