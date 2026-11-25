package planque

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeContract checks msg, as planque posts it, against testdata/<name>: the
// file another side's tests read (tafs copies it byte for byte into its
// contracts/). What changes from run to run is made stable first: the ids in
// ids by their names, planque's own signature over a welcome, and its clock.
// UPDATE_CONTRACTS=1 writes the file again.
func writeContract(t *testing.T, name string, msg map[string]any, ids map[string]string) {
	t.Helper()
	m := map[string]any{}
	for k, v := range msg {
		m[k] = v
	}
	if _, ok := m["sig"]; ok {
		m["sig"] = "<planque's signature>"
	}
	if _, ok := m["now"]; ok {
		m["now"] = 0
	}
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw) + "\n"
	for id, as := range ids {
		got = strings.ReplaceAll(got, id, as)
	}
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_CONTRACTS") == "1" {
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil || string(want) != got {
		t.Fatalf("planque's %s changed: write it again (UPDATE_CONTRACTS=1), and tafs' contracts/ copy with it:\n%s", name, got)
	}
}
