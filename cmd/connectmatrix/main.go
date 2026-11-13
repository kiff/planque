// Command connectmatrix reads each app's published API description as planque does, and says whether the app can
// be connected by its key: the description found and read (JSON or YAML), its key expressed, its tools listed. Each
// red is a generic cause to fix in planque, never one vendor's patch. It reads public pages only and, given
// planque's sign-in reader key (managed.go), planque's answer of whether the managed sign-in works for each app:
// never a session, an account or a call, and no provider key.
//
//	go run ./cmd/connectmatrix deploy/connect-matrix.json
//
// prints a Markdown table and exits 0: it reports, it doesn't gate. This repo holds the code and the app list only:
// no workflow runs it here, so its results and the reader key stay out of a public log.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/kiff/planque/internal/openapi"
	"go.yaml.in/yaml/v3"
)

type app struct {
	App     string `json:"app"`
	URL     string `json:"url"`
	None    string `json:"none"`
	Managed string `json:"managed"`
}

// row is one app's line in the table.
type row struct {
	App, Description, Route, Ops, Why, Managed string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: connectmatrix <apps.json>")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var list struct {
		Apps []app `json:"apps"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	client := &http.Client{Timeout: 90 * time.Second}
	provider := managedReader()
	var rows []row
	for _, a := range list.Apps {
		r := check(client, a)
		r.Managed = managed(provider, a.Managed)
		rows = append(rows, r)
	}
	fmt.Print(table(rows))
}

func table(rows []row) string {
	var b strings.Builder
	green := 0
	for _, r := range rows {
		if strings.HasPrefix(r.Route, "✅") {
			green++
		}
	}
	fmt.Fprintf(&b, "## Connect matrix: %d of %d apps can be connected by key\n\n", green, len(rows))
	b.WriteString("| App | Description | Key route | Tools | Why | Managed sign-in |\n|---|---|---|---|---|---|\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", r.App, r.Description, r.Route, r.Ops, strings.ReplaceAll(r.Why, "|", "/"), r.Managed)
	}
	return b.String()
}

func check(client *http.Client, a app) row {
	r := row{App: a.App, Route: "❌", Ops: "-"}
	if a.URL == "" {
		r.Description, r.Why = "none", a.None
		return r
	}
	body, err := fetch(client, a.URL)
	if err != nil {
		r.Description, r.Why = "not fetched", err.Error()
		return r
	}
	r.Description = describe(body)
	if !strings.HasPrefix(r.Description, "OpenAPI 3") {
		r.Why = map[string]string{"Google Discovery": "a Google Discovery document, not OpenAPI: a managed sign-in only",
			"Swagger 2": "Swagger 2.0, not OpenAPI 3", "HTML": "a web page, not a description"}[r.Description]
		return r
	}
	r.Route, r.Ops, r.Why = keyRoute(body)
	return r
}

// fetch reads one public description, at most a little over planque's own limit, so an oversized one is seen as such.
func fetch(client *http.Client, url string) ([]byte, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, errors.New("couldn't be fetched")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("answered %d", res.StatusCode)
	}
	return io.ReadAll(io.LimitReader(res.Body, openapi.MaxBytes+1))
}

func describe(body []byte) string {
	t := bytes.TrimSpace(body)
	size := fmt.Sprintf(" (%s)", sizeOf(len(body)))
	switch {
	case len(t) == 0:
		return "empty"
	case t[0] == '<':
		return "HTML"
	case t[0] == '{':
		var doc map[string]any
		if json.Unmarshal(t, &doc) != nil {
			if len(body) > openapi.MaxBytes {
				return "JSON over 16 MB"
			}
			return "not JSON"
		}
		switch {
		case doc["kind"] == "discovery#restDescription":
			return "Google Discovery"
		case doc["swagger"] != nil:
			return "Swagger 2"
		case doc["openapi"] != nil:
			return fmt.Sprintf("OpenAPI %v, JSON%s", doc["openapi"], size)
		}
		return "JSON"
	}
	if m := regexp.MustCompile(`(?m)^openapi:\s*['"]?([0-9.]+)`).FindSubmatch(t); m != nil {
		return fmt.Sprintf("OpenAPI %s, YAML%s", m[1], size)
	}
	if regexp.MustCompile(`(?m)^swagger:`).Match(t) {
		return "Swagger 2"
	}
	return "unknown"
}

func sizeOf(n int) string {
	if n >= 1<<20 {
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	}
	return fmt.Sprintf("%d KB", n>>10)
}

// keyRoute reads the description as planque would for a key grant: whole, then as Coach asks, with three reads
// selected (uses). It answers the route, the tools read, and why not.
func keyRoute(body []byte) (route, ops, why string) {
	if len(body) > openapi.MaxBytes {
		return "❌", "-", "over 16 MB"
	}
	doc, err := asDoc(body)
	if err != nil {
		return "❌", "-", err.Error()
	}
	auths, others := keyAuths(doc)
	if len(auths) == 0 {
		if len(others) == 0 {
			return "❌", "-", "names no way to send a key"
		}
		return "❌", "-", "its key goes as " + strings.Join(others, ", ") + ", which planque can't send"
	}
	var lastErr error
	for _, a := range auths {
		s, err := openapi.Parse(body, a)
		if err == nil {
			return "✅ " + authWords(a), fmt.Sprint(len(s.Operations)), skipped(s)
		}
		lastErr = err
	}
	sel := reads(doc, 3)
	for _, a := range auths {
		s, err := openapi.ParseOperations(body, a, sel)
		if err == nil {
			return "✅ " + authWords(a) + ", as selected", fmt.Sprint(len(s.Operations)), "whole: " + refusal(lastErr)
		}
		lastErr = err
	}
	return "❌", "-", refusal(lastErr)
}

func authWords(a openapi.Auth) string {
	if a.Prefix != "" {
		return fmt.Sprintf("%s %s %q", a.In, a.Name, a.Prefix)
	}
	return a.In + " " + a.Name
}

func refusal(err error) string {
	var r *openapi.Refused
	if errors.As(err, &r) {
		return r.Why
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

func skipped(s *openapi.Snapshot) string {
	if len(s.Skipped) == 0 {
		return ""
	}
	whys := map[string]int{}
	for _, k := range s.Skipped {
		whys[k.Why]++
	}
	var out []string
	for w, n := range whys {
		out = append(out, fmt.Sprintf("%d skipped: %s", n, w))
	}
	sort.Strings(out)
	if len(out) > 3 {
		out = append(out[:3], "…")
	}
	return strings.Join(out, "; ")
}

// asDoc reads the description's JSON (a YAML one decoded as YAML 1.2), to look at its schemes and paths. Whether
// planque reads it is openapi.Parse's answer, never this.
func asDoc(body []byte) (map[string]any, error) {
	t := bytes.TrimSpace(body)
	var v any
	if len(t) > 0 && t[0] == '{' {
		if err := json.Unmarshal(t, &v); err != nil {
			return nil, errors.New("it isn't valid JSON")
		}
	} else if err := yaml.Unmarshal(t, &v); err != nil {
		return nil, errors.New("it isn't valid YAML")
	}
	doc, ok := stringKeys(v).(map[string]any)
	if !ok {
		return nil, errors.New("it isn't an OpenAPI description")
	}
	return doc, nil
}

func stringKeys(v any) any {
	switch t := v.(type) {
	case map[string]any:
		for k, x := range t {
			t[k] = stringKeys(x)
		}
		return t
	case map[any]any:
		out := make(map[string]any, len(t))
		for k, x := range t {
			out[fmt.Sprint(k)] = stringKeys(x)
		}
		return out
	case []any:
		for i, x := range t {
			t[i] = stringKeys(x)
		}
	}
	return v
}

// keyAuths are the ways planque can send a key that the description's schemes name, and the others, in words.
func keyAuths(doc map[string]any) (auths []openapi.Auth, others []string) {
	comps, _ := doc["components"].(map[string]any)
	schemes, _ := comps["securitySchemes"].(map[string]any)
	names := make([]string, 0, len(schemes))
	for n := range schemes {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		s, _ := schemes[n].(map[string]any)
		typ, _ := s["type"].(string)
		switch typ {
		case "apiKey":
			in, _ := s["in"].(string)
			name, _ := s["name"].(string)
			if in == "header" || in == "query" {
				auths = append(auths, openapi.Auth{In: in, Name: name})
			} else {
				others = append(others, "an apiKey in "+in)
			}
		case "http":
			sch, _ := s["scheme"].(string)
			if strings.EqualFold(sch, "bearer") {
				auths = append(auths, openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer "})
			} else {
				others = append(others, "HTTP "+strings.ToLower(sch))
			}
		case "oauth2", "openIdConnect":
			others = append(others, "OAuth (a managed sign-in)")
		default:
			others = append(others, typ)
		}
	}
	return auths, others
}

var notWord = regexp.MustCompile(`[^A-Za-z0-9]+`)

// opName is planque's name for an operation (as openapi names it).
func opName(method, path string, o map[string]any) string {
	if id, _ := o["operationId"].(string); id != "" {
		if n := strings.Trim(notWord.ReplaceAllString(id, "_"), "_"); openapi.ToolName.MatchString(n) {
			return n
		}
	}
	n := strings.Trim(notWord.ReplaceAllString(method+"_"+path, "_"), "_")
	if len(n) > 64 {
		n = n[:64]
	}
	if n == "" || !openapi.ToolName.MatchString(n) {
		n = "op_" + n
	}
	return strings.ToLower(n)
}

// reads are up to n GET operations by planque's names, in path order: what a taf asks for to read.
func reads(doc map[string]any, n int) []string {
	paths, _ := doc["paths"].(map[string]any)
	keys := make([]string, 0, len(paths))
	for p := range paths {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	var out []string
	for _, p := range keys {
		item, _ := paths[p].(map[string]any)
		if o, ok := item["get"].(map[string]any); ok {
			out = append(out, opName("get", p, o))
			if len(out) == n {
				break
			}
		}
	}
	return out
}
