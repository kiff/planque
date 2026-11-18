package connector

import (
	"os"
	"testing"
)

// The older tests use fake tools with any name (a void, a delete); the class
// hold has its own tests, which turn it on.
func TestMain(m *testing.M) {
	HoldClasses = false
	os.Exit(m.Run())
}
