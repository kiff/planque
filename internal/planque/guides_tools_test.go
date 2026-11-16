package planque

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// A tool's name: letters and digits in two or more parts joined by
// underscores, as GMAIL_SEND_EMAIL or stripe_api_read.
var guideToolName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*(_[A-Za-z0-9]+)+$")

var guideBackticked = regexp.MustCompile("`([^`]+)`")

// A row of a guide's table whose first header cell is Tool.
type guideToolRow struct {
	cells, header []string
}

// The row's "Read or write" cell, if its table has one.
func (r guideToolRow) access() (string, bool) {
	i := slices.Index(r.header, "Read or write")
	if i < 0 || i >= len(r.cells) {
		return "", false
	}
	return r.cells[i], true
}

func guideToolRows(text string) []guideToolRow {
	var out []guideToolRow
	var header []string
	inTable := false
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			header, inTable = nil, false
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		if !inTable {
			inTable = true
			if cells[0] == "Tool" {
				header = cells
			}
			continue
		}
		if header == nil || strings.Trim(strings.Join(cells, ""), "-: ") == "" {
			continue
		}
		out = append(out, guideToolRow{cells, header})
	}
	return out
}

// The tools of an app's generated guide, each with its "Read or write"
// cell. An app whose tools come in several parts has one generated file
// per part, <slug>-<part>.md.
func generatedGuideTools(slug string) (map[string]string, bool) {
	dir := "../../deploy/guides/generated/"
	files := []string{dir + slug + ".md"}
	if _, err := os.Stat(files[0]); err != nil {
		files, _ = filepath.Glob(dir + slug + "-*.md")
	}
	if len(files) == 0 {
		return nil, false
	}
	tools := map[string]string{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, false
		}
		for _, r := range guideToolRows(string(raw)) {
			said, _ := r.access()
			for _, m := range guideBackticked.FindAllStringSubmatch(r.cells[0], -1) {
				tools[m[1]] = said
			}
		}
	}
	return tools, true
}

// Every tool an adapted guide names is a tool of that app's generated
// guide, or one of planque's own tools for an app with a describer; and a
// tools-table row marks it read exactly when the generated guide does. A
// guide that names a tool the catalog doesn't have, or calls a write a
// read, sends Coach to plan with it.
func TestEveryToolAGuideNamesIsInItsGeneratedGuide(t *testing.T) {
	files, err := filepath.Glob("../../deploy/guides/*.md")
	if err != nil || len(files) == 0 {
		t.Fatalf("no guides: %v", err)
	}
	// planque's own tools, and the argument names a guide may quote.
	served := map[string]string{}
	arguments := map[string]bool{}
	for _, d := range describers {
		served[d.Tool] = "read"
		for _, n := range d.Reads {
			served[n] = "read"
		}
		for _, n := range d.Writes {
			served[n] = "write"
		}
		arguments[d.OpArg] = true
		for _, n := range d.Context {
			arguments[n] = true
		}
	}
	checked := 0
	for _, f := range files {
		slug := strings.TrimSuffix(filepath.Base(f), ".md")
		if slug == "README" || slug == "apps" || slug == "_refusals" {
			continue
		}
		generated, ok := generatedGuideTools(slug)
		if !ok {
			t.Errorf("%s: no generated guide", slug)
			continue
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		text := string(raw)
		lookup := func(name string) (string, bool) {
			if said, ok := generated[name]; ok {
				return said, true
			}
			said, ok := served[name]
			return said, ok
		}
		checked++

		// The tools tables: every name, and its read or write.
		for _, r := range guideToolRows(text) {
			said, hasAccess := r.access()
			for _, m := range guideBackticked.FindAllStringSubmatch(r.cells[0], -1) {
				name := m[1]
				if !guideToolName.MatchString(name) {
					continue
				}
				gen, ok := lookup(name)
				if !ok {
					t.Errorf("%s: the tools table names `%s`, not a tool of the generated guide", slug, name)
					continue
				}
				if hasAccess && gen != "" && (said == "read") != (gen == "read") {
					t.Errorf("%s: `%s` is %q here, %q in the generated guide", slug, name, said, gen)
				}
			}
		}

		// Everywhere else: a name with one of the app's tool prefixes.
		prefixes := map[string]bool{}
		for name := range generated {
			if head, _, ok := strings.Cut(name, "_"); ok {
				prefixes[head+"_"] = true
			}
		}
		for _, m := range guideBackticked.FindAllStringSubmatch(text, -1) {
			name := m[1]
			if !guideToolName.MatchString(name) || arguments[name] {
				continue
			}
			head, _, _ := strings.Cut(name, "_")
			if !prefixes[head+"_"] {
				continue
			}
			if _, ok := lookup(name); !ok {
				t.Errorf("%s: names `%s`, not a tool of the generated guide", slug, name)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no adapted guides checked")
	}
}
