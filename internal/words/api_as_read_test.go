package words

import (
	"strings"
	"testing"
)

// The approval card names the API's host only: the description's address and
// its digest are in the signed request (Request.API), not on the card the owner
// reads.
func TestAPIAsReadNamesOnlyTheHost(t *testing.T) {
	for lang, a := range approvals {
		if n := strings.Count(a.APIAsReadF, "%s"); n != 1 || strings.Contains(a.APIAsReadF, "SHA-256") {
			t.Errorf("%s: APIAsReadF must take the host alone, with no digest: %q", lang, a.APIAsReadF)
		}
	}
}
