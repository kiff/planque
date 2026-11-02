package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/jcs"
)

// These are KIFF's synthetic kiff.profile.v1 / kiff.words.v1 vectors.
func TestContractVectors(t *testing.T) {
	files, err := filepath.Glob("testdata/*.expected.json")
	if err != nil || len(files) != 5 {
		t.Fatal("missing contract vectors")
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			raw, err := os.ReadFile(strings.Replace(file, ".expected.json", ".json", 1))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Parse(raw); err != nil {
				t.Fatal(err)
			}
			expected, _ := os.ReadFile(file)
			var want struct {
				ProfileSHA256 string          `json:"profile_sha256"`
				Words         json.RawMessage `json:"words"`
				WordsSHA256   string          `json:"words_sha256"`
			}
			if err := json.Unmarshal(expected, &want); err != nil {
				t.Fatal(err)
			}
			digest, err := Digest(raw)
			if err != nil || digest != want.ProfileSHA256 {
				t.Fatalf("profile digest: %s %v", digest, err)
			}
			words, err := jcs.SHA256(want.Words)
			if err != nil || words != want.WordsSHA256 {
				t.Fatalf("words digest: %s %v", words, err)
			}
		})
	}
}

func TestProfileRefusesUnknownAndTrailingFields(t *testing.T) {
	raw, _ := os.ReadFile("testdata/holded-invoices.json")
	for _, changed := range []string{strings.Replace(string(raw), `"v":`, `"unknown":true,"v":`, 1), string(raw) + " null", string(raw) + " {}"} {
		if _, err := Parse([]byte(changed)); err == nil {
			t.Fatal("accepted an ambiguous or unknown contract value")
		}
	}
}
