package planque

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/words"
)

// A paused grant's status words come from the kind of pause, in the person's
// language, and never from its Why: Why is English, for the log and for Coach,
// and names the tool ("Google Sheets changed GOOGLESHEETS_BATCH_UPDATE since
// you approved it"). A German person read it verbatim.
func TestAPausedGrantReadsInThePersonsWordsNeverItsEnglishWhy(t *testing.T) {
	for _, c := range []struct {
		kind string
		pick func(words.Approval) string
	}{
		{pauseChanged, func(w words.Approval) string { return w.PausedChanged }},
		{pauseGone, func(w words.Approval) string { return w.PausedGone }},
		{"", func(w words.Approval) string { return w.PausedChanged }}, // paused before kinds were kept
	} {
		for _, lang := range words.Languages {
			w := words.Approve(lang)
			g := &Grant{Status: StatusPending, Paused: &Pause{Tool: "sheets_batch_update", Kind: c.kind, Why: "Google Sheets changed GOOGLESHEETS_BATCH_UPDATE since you approved it"}}
			got := statusIn(w, g)
			if want := fmt.Sprintf(w.PausedF, c.pick(w)); got != want {
				t.Errorf("%s kind %q: %q, want %q", lang, c.kind, got, want)
			}
			for _, leak := range []string{"GOOGLESHEETS", "BATCH_UPDATE", "batch_update", "Google Sheets changed"} {
				if strings.Contains(got, leak) {
					t.Errorf("%s kind %q: %q carries %q", lang, c.kind, got, leak)
				}
			}
		}
	}
}
