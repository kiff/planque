package fakes

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// Billing's docs describe what billing's own file does: the index links a
// guide and four reference pages, each reference page holds one JSON
// fragment, and the fragments merged parse to the same three operations,
// digest for digest, create_invoice's final-by-default status included.
func TestBillingDocsDescribeWhatItsFileDoes(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	get := func(url string) (int, string) {
		res, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}

	status, index := get(v.BillingDocsIndex())
	if status != 200 || !strings.HasPrefix(index, "# Fake Billing\n") {
		t.Fatalf("llms.txt: %d %q", status, index)
	}
	if status, _ := get(v.root() + "/billing/docs/nothing.md"); status != 404 {
		t.Errorf("an unknown page: %d, want 404", status)
	}
	links := regexp.MustCompile(`\]\(([^)]+\.md)\)`).FindAllStringSubmatch(index, -1)
	if len(links) != 5 {
		t.Fatalf("the index links %d pages, want 5: %q", len(links), index)
	}
	fence := regexp.MustCompile("(?s)```json\n(.*?)\n```")
	merged := map[string]any{"paths": map[string]any{}, "components": map[string]any{}}
	fragments := 0
	for _, l := range links {
		status, page := get(l[1])
		if status != 200 {
			t.Fatalf("%s: %d", l[1], status)
		}
		blocks := fence.FindAllStringSubmatch(page, -1)
		if strings.HasSuffix(l[1], "/getting-started.md") {
			if len(blocks) != 0 {
				t.Errorf("the guide holds a fragment")
			}
			continue
		}
		if len(blocks) != 1 {
			t.Fatalf("%s holds %d fragments, want 1", l[1], len(blocks))
		}
		fragments++
		var f map[string]any
		if err := json.Unmarshal([]byte(blocks[0][1]), &f); err != nil {
			t.Fatalf("%s: %v", l[1], err)
		}
		for k, val := range f {
			switch k {
			case "paths", "components":
				into := merged[k].(map[string]any)
				for name, part := range val.(map[string]any) {
					if k == "components" {
						have, _ := into[name].(map[string]any)
						if have == nil {
							have = map[string]any{}
						}
						for n, c := range part.(map[string]any) {
							have[n] = c
						}
						into[name] = have
						continue
					}
					into[name] = part
				}
			default:
				merged[k] = val
			}
		}
	}
	if fragments != 4 {
		t.Fatalf("%d fragments, want 4", fragments)
	}
	raw, _ := json.Marshal(merged)
	auth := openapi.Auth{In: "header", Name: BillingKeyHeader}
	fromDocs, err := openapi.Parse(raw, auth)
	if err != nil {
		t.Fatalf("the fragments merged: %v", err)
	}
	_, file := get(v.BillingDescription())
	fromFile, err := openapi.Parse([]byte(file), auth)
	if err != nil {
		t.Fatal(err)
	}
	digests := func(s *openapi.Snapshot) map[string]string {
		out := map[string]string{}
		for _, op := range s.Operations {
			d, err := op.Digest()
			if err != nil {
				t.Fatal(err)
			}
			out[op.Name] = d
		}
		return out
	}
	docs, want := digests(fromDocs), digests(fromFile)
	if len(want) != 3 || len(docs) != len(want) {
		t.Fatalf("operations from the docs %v, from the file %v", docs, want)
	}
	for name, d := range want {
		if docs[name] != d {
			t.Errorf("%s from the docs isn't the file's", name)
		}
	}
	if strings.Join(fromDocs.Servers, " ") != strings.Join(fromFile.Servers, " ") || len(fromFile.Servers) != 1 {
		t.Errorf("the docs' servers %q, the file's %q", fromDocs.Servers, fromFile.Servers)
	}
}
