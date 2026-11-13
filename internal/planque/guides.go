package planque

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/connector"
)

// A guide is what Coach reads about an app before it offers a job with it.
// Its generated part is written from the catalog as planque reads it: how
// the app connects, each tool's access, inputs and returns. Nothing in it
// is hand-written, so a run over the whole catalog gives every app its
// part at once; the rest of a guide (limits, plans, what to tell the owner
// on an error) is adapted by hand beside it, in deploy/guides.

// GuideTool is one tool as a guide lists it.
type GuideTool struct {
	Name     string
	Title    string
	ReadOnly bool
	// Access are the scopes any one of which lets it run, as the vendor's
	// API lists them; empty when the catalog doesn't say.
	Access []string
	// Inputs are its required arguments, then its optional ones.
	Inputs, Optional []string
	// Returns are the ids and links its answer gives back, by path.
	Returns []string
}

// GuideApp is one app of the catalog, as a guide's generated part has it.
type GuideApp struct {
	Slug, Name string
	// Connects says how: "sign-in" (a managed app), "key" (a reviewed
	// entry with an API key), or "server" (a vendor's own server, whose
	// tools are listed only once signed in).
	Connects string
	// Scheme is the managed app's sign-in scheme (OAUTH2, API_KEY…).
	Scheme string
	Tools  []GuideTool
	// Skipped counts tools the catalog listed that couldn't be read.
	Skipped int
}

// guideTools bounds the tools one app's guide reads: past it, the guide
// says how many were left out.
const guideTools = 1000

// ManagedGuideApp is a managed app's generated part: each tool the
// provider lists, its access, inputs and returns.
func ManagedGuideApp(app connector.App, tools []connector.Tool) GuideApp {
	g := GuideApp{Slug: app.Slug, Name: cmp.Or(app.Name, app.Slug), Connects: "sign-in", Scheme: app.AuthScheme}
	for _, t := range tools {
		required, optional := schemaInputs(t.InputSchema)
		g.Tools = append(g.Tools, GuideTool{Name: t.Slug, Title: label(cmp.Or(t.Name, t.Slug)), ReadOnly: t.ReadOnly,
			Access: t.Scopes, Inputs: required, Optional: optional, Returns: slices.DeleteFunc(returnsOf(t.Output), serviceNamed)})
	}
	slices.SortFunc(g.Tools, func(a, b GuideTool) int { return strings.Compare(a.Name, b.Name) })
	return g
}

// EntryGuideApp is a reviewed catalog entry's generated part: the tools
// planque installed for it, as reviewed.
func EntryGuideApp(id string, e CatalogEntry) GuideApp {
	g := GuideApp{Slug: id, Name: cmp.Or(e.Vendor, id), Connects: "key"}
	if e.Server != "" && e.API == nil {
		g.Connects = "server"
	}
	for _, t := range e.Tools {
		var fixed []string
		for k := range t.Fixed {
			fixed = append(fixed, k)
		}
		slices.Sort(fixed)
		var args []string
		for k := range t.Arguments {
			args = append(args, k)
		}
		slices.Sort(args)
		g.Tools = append(g.Tools, GuideTool{Name: cmp.Or(t.Name, t.Tool), Title: t.Title, ReadOnly: t.ReadOnly, Inputs: args})
	}
	return g
}

// serviceNamed says whether a field is named after the connection service
// rather than the vendor: its own plumbing, never part of a public guide.
func serviceNamed(name string) bool {
	return strings.Contains(strings.ToLower(name), "composio")
}

// schemaInputs reads a tool's arguments from its input schema: the
// required ones, then the rest, each sorted.
func schemaInputs(schema json.RawMessage) (required, optional []string) {
	var s struct {
		Properties map[string]json.RawMessage `json:"properties"`
		Required   []string                   `json:"required"`
	}
	if json.Unmarshal(schema, &s) != nil {
		return nil, nil
	}
	for name := range s.Properties {
		if serviceNamed(name) {
			continue
		}
		if slices.Contains(s.Required, name) {
			required = append(required, name)
		} else {
			optional = append(optional, name)
		}
	}
	slices.Sort(required)
	slices.Sort(optional)
	return required, optional
}

// CollectGuideApps reads every app the provider can sign someone in to
// now (the list Connect apps opens on), and each one's tools.
func CollectGuideApps(ctx context.Context, p connector.Provider) ([]GuideApp, error) {
	var apps []connector.App
	cursor := ""
	for range catalogPagesRead {
		page, err := p.Catalog(ctx, "", cursor)
		if err != nil {
			return nil, err
		}
		apps = append(apps, page.Items...)
		if page.Next == "" || page.Next == cursor {
			break
		}
		cursor = page.Next
	}
	out := make([]GuideApp, 0, len(apps))
	for _, a := range apps {
		var tools []connector.Tool
		skipped, cursor := 0, ""
		for len(tools) < guideTools {
			page, err := p.Tools(ctx, a.Slug, cursor)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", a.Slug, err)
			}
			tools, skipped = append(tools, page.Items...), skipped+page.Skipped
			if page.Next == "" || page.Next == cursor {
				break
			}
			cursor = page.Next
		}
		g := ManagedGuideApp(a, tools)
		g.Skipped = skipped
		out = append(out, g)
	}
	return out, nil
}

const catalogPagesRead = 10

// GeneratedGuide is an app's generated part, in Markdown, dated.
func GeneratedGuide(a GuideApp, checked string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: generated\n\nChecked: %s · Generated from planque's catalog: don't edit by hand, run `planque guides` again.\n\n", a.Name, checked)
	b.WriteString("## How it connects\n\n")
	switch a.Connects {
	case "key":
		b.WriteString("- **By the owner's API key**, reviewed in planque's catalog: the tools below are the ones installed, as reviewed. The key and the taf's Card come in one approval.\n")
	case "server":
		b.WriteString("- **A sign-in to the vendor's own server.** Its tools are listed only once signed in: read the grant's catalog then, and offer only its `tools`.\n")
	default:
		fmt.Fprintf(&b, "- **A sign-in, held by the connection service** (%s). The taf never sees a password or token.\n", schemeWords(a.Scheme))
		b.WriteString("- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.\n")
		b.WriteString("- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.\n")
	}
	reads, writes := 0, 0
	for _, t := range a.Tools {
		if t.ReadOnly {
			reads++
		} else {
			writes++
		}
	}
	fmt.Fprintf(&b, "\n## Tools\n\n%d tools: %d read, %d write (a tool not marked read-only counts as a write).", len(a.Tools), reads, writes)
	if a.Skipped > 0 {
		fmt.Fprintf(&b, " %d more couldn't be read and are left out.", a.Skipped)
	}
	b.WriteString("\n\n| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |\n|---|---|---|---|---|---|\n")
	for _, t := range a.Tools {
		kind := "write"
		if t.ReadOnly {
			kind = "read"
		}
		inputs := codeList(t.Inputs)
		if len(t.Optional) > 0 {
			opt := "*" + strings.Join(t.Optional, ", ") + "*"
			if len(t.Optional) > 8 {
				opt = fmt.Sprintf("*%s and %d more*", strings.Join(t.Optional[:8], ", "), len(t.Optional)-8)
			}
			inputs = strings.TrimPrefix(inputs+"; "+opt, "; ")
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %s | %s |\n", t.Name, cell(t.Title), kind, inputs, codeList(t.Returns), codeList(shortScopes(t.Access)))
	}
	return b.String()
}

func schemeWords(scheme string) string {
	switch strings.ToUpper(scheme) {
	case "OAUTH2", "OAUTH1":
		return "the owner signs in on the vendor's own page"
	case "API_KEY", "BEARER_TOKEN", "BASIC", "BASIC_WITH_JWT":
		return "the owner pastes a key from the vendor's settings"
	case "":
		return "how is the vendor's choice"
	}
	return "by " + strings.ToLower(scheme)
}

// shortScopes drops a scope URL's common prefix, so a table stays legible.
func shortScopes(scopes []string) []string {
	out := make([]string, 0, len(scopes))
	for _, s := range scopes {
		s = strings.TrimPrefix(s, "https://www.googleapis.com/auth/")
		out = append(out, s)
	}
	return out
}

func codeList(xs []string) string {
	if len(xs) == 0 {
		return ""
	}
	return "`" + strings.Join(xs, "`, `") + "`"
}

// cell keeps a table cell on one line.
func cell(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "|", "\\|")
}

// GuideBatches lists every app in batches of size, first the ones named in
// first (what Explore shows and people connect first), in that order, then
// the rest alphabetically: each batch is one worker's share of the
// adapted parts.
func GuideBatches(apps []GuideApp, first []string, size int, checked string) string {
	ordered := slices.Clone(apps)
	rank := func(a GuideApp) int {
		if i := slices.Index(first, a.Slug); i >= 0 {
			return i
		}
		return len(first)
	}
	slices.SortStableFunc(ordered, func(a, b GuideApp) int {
		return cmp.Or(cmp.Compare(rank(a), rank(b)), strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)))
	})
	var b strings.Builder
	fmt.Fprintf(&b, "# Apps and their guides, in batches\n\nChecked: %s · Generated: %d apps, %d batches of at most %d. Each app's generated part is in `generated/`; its adapted guide goes beside this file, named by its slug.\n", checked, len(ordered), (len(ordered)+size-1)/max(size, 1), size)
	for i := 0; i < len(ordered); i += size {
		fmt.Fprintf(&b, "\n## Batch %d\n\n| App | Slug | Connects | Tools |\n|---|---|---|---|\n", i/size+1)
		for _, a := range ordered[i:min(i+size, len(ordered))] {
			fmt.Fprintf(&b, "| %s | `%s` | %s | %d |\n", cell(a.Name), a.Slug, a.Connects, len(a.Tools))
		}
	}
	return b.String()
}
