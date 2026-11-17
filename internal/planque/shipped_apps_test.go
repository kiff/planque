package planque

import (
	"os"
	"testing"

	"github.com/kiff/planque/internal/appconfig"
)

// shippedApps reads the registered apps that ship with the function
// (deploy/apps.json), which the shipped catalog's entries may name.
func shippedApps(t testing.TB) map[string]appconfig.Application {
	t.Helper()
	raw, err := os.ReadFile("../../deploy/apps.json")
	if err != nil {
		t.Fatal(err)
	}
	apps, err := appconfig.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	return apps
}
