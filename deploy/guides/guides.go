// Package guides holds the adapted app guides beside it, built into
// planque so the server serves the ones in this tree (GET
// /v1/accounts/{tenant}/catalog/{entry}/guide).
package guides

import (
	"embed"
	"regexp"
	"slices"
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

// Reads are the tools an app's adapted guide marks "read" in its tools
// table's "Read or write" column, by their slugs as the guide writes them:
// a person checked each one against the vendor's own reference. A tool the
// guide doesn't list, or marks otherwise ("write", "write, outside"), is not
// one: unknown stays a write.
func Reads(slug string) []string {
	text, _, ok := Adapted(slug)
	if !ok {
		return nil
	}
	var out []string
	col := -1
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			col = -1
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		if col < 0 {
			if len(cells) > 0 && cells[0] == "Tool" {
				col = slices.Index(cells, "Read or write")
			}
			continue
		}
		if col >= len(cells) || cells[col] != "read" {
			continue
		}
		for _, m := range toolSlug.FindAllStringSubmatch(cells[0], -1) {
			out = append(out, m[1])
		}
	}
	return out
}

var toolSlug = regexp.MustCompile("`([A-Z][A-Z0-9_]+)`")
