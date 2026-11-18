package planque

import (
	"os"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

// The older tests use fake tools with any name (a void, a delete); the class
// hold has its own tests, which turn it on.
func TestMain(m *testing.M) {
	connector.HoldClasses = false
	os.Exit(m.Run())
}
