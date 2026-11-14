package guides

import (
	"slices"
	"testing"
)

// A guide's reads are the rows its tools table marks "read", each slug in the
// cell (two on one row too); never a write, a write outside, or an app with no guide.
func TestReadsAreTheRowsTheGuideMarksRead(t *testing.T) {
	reads := Reads("gmail")
	for _, want := range []string{"GMAIL_FETCH_EMAILS", "GMAIL_FETCH_MESSAGE_BY_THREAD_ID", "GMAIL_LIST_LABELS", "GMAIL_LIST_DRAFTS", "GMAIL_GET_DRAFT"} {
		if !slices.Contains(reads, want) {
			t.Errorf("%s isn't among gmail's reads %v", want, reads)
		}
	}
	for _, not := range []string{"GMAIL_CREATE_EMAIL_DRAFT", "GMAIL_SEND_EMAIL", "GMAIL_REPLY_TO_THREAD", "GMAIL_DELETE_MESSAGE"} {
		if slices.Contains(reads, not) {
			t.Errorf("%s is a read: %v", not, reads)
		}
	}
	if got := Reads("no-such-app"); got != nil {
		t.Errorf("an app with no guide has reads: %v", got)
	}
	if got := Reads("readme"); got != nil {
		t.Errorf("the README is no guide: %v", got)
	}
}
