package words

import (
	"strings"
	"testing"
)

// The approval card names the API's host only: the description's address and
// its digest are in the signed request (Request.API), not on the card the owner
// reads. "Reviewed" and "published" are two lines, so neither is said of the
// other kind of key grant.
func TestAPIAsReadNamesOnlyTheHost(t *testing.T) {
	for lang, a := range approvals {
		for name, f := range map[string]string{"APIAsReadF": a.APIAsReadF, "APIPublishedF": a.APIPublishedF} {
			if n := strings.Count(f, "%s"); n != 1 || strings.Contains(f, "SHA-256") {
				t.Errorf("%s %s: must take the host alone, with no digest: %q", lang, name, f)
			}
		}
		if a.APIAsReadF == a.APIPublishedF {
			t.Errorf("%s: the reviewed and published lines are the same: %q", lang, a.APIAsReadF)
		}
	}
}
