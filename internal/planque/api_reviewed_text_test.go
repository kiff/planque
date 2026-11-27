package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/words"
)

// A key grant from a reviewed catalog entry says its tools were reviewed; one
// asked by its own description's URL says they're what that published
// description lists. Neither shows the description's address or digest.
func TestTheKeyCardSaysReviewedOnlyForAReviewedEntry(t *testing.T) {
	q := Request{API: &APIRequest{URL: "https://raw.example/openapi.yaml", Digest: strings.Repeat("a", 64), Base: "https://api.vendor.example/v1",
		Auth: openapi.Auth{In: "header", Name: "Authorization"}}}
	for _, c := range []struct{ kind, want, not string }{
		{"reviewed", "Its tools are the ones reviewed for api.vendor.example.", "published description"},
		{"", "Its tools are the ones its published description lists for api.vendor.example.", "reviewed"},
	} {
		q.API.SourceKind = c.kind
		got := apiKeyText(words.Approve("en"), q, false, "")
		if !strings.Contains(got, c.want) || strings.Contains(got, c.not) {
			t.Errorf("source %q: %s", c.kind, got)
		}
		if strings.Contains(got, "raw.example") || strings.Contains(got, "SHA-256") || strings.Contains(got, strings.Repeat("a", 12)) {
			t.Errorf("source %q: the card shows the description's address or digest: %s", c.kind, got)
		}
	}
}

// A key grant is marked reviewed only when its description was checked against
// a catalog entry's digest; one asked by its own URL is not.
func TestOnlyACheckedDescriptionIsMarkedReviewed(t *testing.T) {
	fromEntry := &APIRequest{URL: "https://raw.example/openapi.yaml"}
	markReviewed(fromEntry, strings.Repeat("b", 64))
	if fromEntry.SourceKind != "reviewed" {
		t.Errorf("from a catalog entry: %q", fromEntry.SourceKind)
	}
	byURL := &APIRequest{URL: "https://raw.example/openapi.yaml"}
	markReviewed(byURL, "")
	if byURL.SourceKind != "" {
		t.Errorf("asked by its own URL: %q", byURL.SourceKind)
	}
	app := &APIRequest{SourceKind: "connector"}
	markReviewed(app, strings.Repeat("b", 64))
	if app.SourceKind != "connector" {
		t.Errorf("an existing kind is kept: %q", app.SourceKind)
	}
	markReviewed(nil, "x") // no API: nothing to mark, no panic
}
