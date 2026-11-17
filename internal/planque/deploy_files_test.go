package planque

import (
	"os"
	"strings"
	"testing"
)

// The function reads three reviewed files at start. CI's code-only deploy and
// deploy/aws.sh must both put every one of them in the zip: a file left out of
// one of them is a file the live function can't read after that deploy.
func TestCIAndAWSShipTheSameReviewedFiles(t *testing.T) {
	for _, path := range []string{"../../.github/workflows/ci.yml", "../../deploy/aws.sh"} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var zips []string
		for _, line := range strings.Split(string(raw), "\n") {
			if strings.Contains(line, "zip -q function.zip bootstrap server-apps.json") {
				zips = append(zips, line)
			}
		}
		if len(zips) == 0 {
			t.Fatalf("%s has no zip line that ships the reviewed files", path)
		}
		for _, line := range zips {
			// Whole names: "server-apps.json" contains the text "apps.json".
			shipped := map[string]bool{}
			for _, name := range strings.Fields(line) {
				shipped[strings.Trim(name, "();")] = true
			}
			for _, file := range []string{"server-apps.json", "catalog.json", "apps.json"} {
				if !shipped[file] {
					t.Fatalf("%s: the zip leaves out %s: %s", path, file, strings.TrimSpace(line))
				}
			}
		}
	}
}
