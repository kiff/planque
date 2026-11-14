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

// A schema generated from typed models, as the provider's output_parameters
// are: the answer's data is a reference into $defs (the docs show it with no
// type), an optional field is anyOf with null, and a model may contain
// itself. Each is followed, so Sheets' create names spreadsheetId.
func TestReturnsFollowAModelSchemasReferences(t *testing.T) {
	schema := []byte(`{"type":"object","title":"CreateGoogleSheet1Response","required":["data","successful"],
	 "properties":{
	  "data":{"$ref":"#/$defs/CreateGoogleSheet1ResponseData","description":"Data from the action execution"},
	  "error":{"anyOf":[{"type":"string"},{"type":"null"}],"default":null},
	  "successful":{"type":"boolean"}},
	 "$defs":{
	  "CreateGoogleSheet1ResponseData":{"type":"object","properties":{
	    "spreadsheetId":{"type":"string"},
	    "spreadsheetUrl":{"anyOf":[{"type":"string"},{"type":"null"}]},
	    "properties":{"$ref":"#/$defs/SpreadsheetProperties"},
	    "sheets":{"type":"array","items":{"$ref":"#/$defs/Sheet"}}}},
	  "SpreadsheetProperties":{"type":"object","properties":{"title":{"type":"string"},"locale":{"type":"string"}}},
	  "Sheet":{"type":"object","properties":{"properties":{"anyOf":[{"$ref":"#/$defs/SheetProperties"},{"type":"null"}]}}},
	  "SheetProperties":{"type":"object","properties":{"sheetId":{"type":"integer"},"title":{"type":"string"},"parent":{"$ref":"#/$defs/SheetProperties"}}}}}`)
	got := returnsOf(schema)
	for _, want := range []string{"data.spreadsheetId", "data.spreadsheetUrl", "data.sheets.properties.sheetId"} {
		if !slices.Contains(got, want) {
			t.Fatalf("returns: %v, missing %s", got, want)
		}
	}
	// Under "definitions" too; and a reference out of the schema is never read.
	old := []byte(`{"properties":{"data":{"allOf":[{"$ref":"#/definitions/D"}]},"other":{"$ref":"https://example.com/x.json"}},"definitions":{"D":{"properties":{"file_id":{"type":"string"}}}}}`)
	if got := returnsOf(old); !slices.Equal(got, []string{"data.file_id"}) {
		t.Fatalf("definitions: %v", got)
	}
}
