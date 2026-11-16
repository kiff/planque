package guides

import (
	"slices"
	"strings"
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

// Every guide ends with what a failed call returns, whatever the app, with
// the app's own name in it; the shared part is never a guide of its own.
func TestEveryGuideEndsWithWhatAFailedCallReturns(t *testing.T) {
	text, checked, ok := Adapted("airtable")
	if !ok || checked == "" {
		t.Fatalf("airtable: ok %v, checked %q", ok, checked)
	}
	part := text[strings.Index(text, "## What the taf gets back"):]
	if !strings.Contains(part, "`not_granted`") || !strings.Contains(part, "Airtable refused it (it answered 403)") {
		t.Errorf("airtable's guide doesn't end with the shared part, with its name:\n%s", part)
	}
	if strings.Contains(text, "{app}") {
		t.Error("a {app} is left in the served guide")
	}
	if strings.Count(text, "## What the taf gets back") != 1 {
		t.Error("the shared part isn't there exactly once")
	}
	if _, _, ok := Adapted("_refusals"); ok {
		t.Error("the shared part is served as a guide")
	}
	if got := Reads("airtable"); slices.Contains(got, "not_granted") {
		t.Errorf("the shared table's codes count as reads: %v", got)
	}
}
