package main

import (
	"io"
	"log/slog"
	"testing"
)

// The server starts on the shipped catalog with no apps file: an entry that
// names an app it doesn't have is left out, not a reason to refuse to start.
// Production had no PLANQUE_APPS_FILE when googlesheets@1 merged: the new code
// exited at start, answered 502, and the deploy's health check rolled it back.
func TestTheServerStartsOnTheShippedCatalogWithNoAppsFile(t *testing.T) {
	t.Setenv("PLANQUE_BASE_URL", "https://planque.test")
	t.Setenv("PLANQUE_APP_RETURN", "https://app.test/return")
	t.Setenv("PLANQUE_LOCAL_KEY", "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=")
	t.Setenv("KIFF_API", "https://kiff.test")
	t.Setenv("KIFF_ISSUER", "https://kiff.test")
	t.Setenv("PLANQUE_CATALOG_FILE", "../../deploy/catalog.json")
	t.Setenv("PLANQUE_APPS_FILE", "")
	s, err := build(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("the server doesn't start: %v", err)
	}
	if _, ok := s.CatalogEntries["googlesheets@1"]; ok {
		t.Error("an entry whose app isn't registered was loaded")
	}
	if _, ok := s.CatalogEntries["zendesk-key@1"]; !ok {
		t.Error("the other entries weren't loaded")
	}
}
