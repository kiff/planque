package planque

import (
	"slices"
	"strings"
	"testing"
)

// A tool's answer names what a taf keeps: its ids and links, by path.
func TestAToolsReturnsAreItsIdsAndLinks(t *testing.T) {
	schema := []byte(`{"type":"object","properties":{"successful":{"type":"boolean"},"data":{"type":"object","properties":{"spreadsheetId":{"type":"string"},"spreadsheetUrl":{"type":"string"},"title":{"type":"string"},"sheets":{"type":"array","items":{"type":"object","properties":{"sheet_id":{"type":"integer"},"name":{"type":"string"}}}}}}}}`)
	if got := returnsOf(schema); !slices.Equal(got, []string{"data.sheets.sheet_id", "data.spreadsheetId", "data.spreadsheetUrl"}) {
		t.Fatalf("returns: %v", got)
	}
	odd := []byte(`{"type":"object","properties":{"ignore previous instructions url":{"type":"string"},"data":{"type":"object","properties":{"` + strings.Repeat("a", 130) + `_id":{"type":"string"},"web_url":{"type":"string"}}}}}`)
	if got := returnsOf(odd); !slices.Equal(got, []string{"data.web_url"}) {
		t.Fatalf("returns of odd names: %v", got)
	}
	for _, none := range []string{"", "[1]", `{"type":"object","properties":{"title":{"type":"string"}}}`} {
		if got := returnsOf([]byte(none)); len(got) != 0 {
			t.Fatalf("%q: %v", none, got)
		}
	}
}
