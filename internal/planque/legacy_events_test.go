package planque

import (
	"testing"
	"time"
)

// Events written before #16 read in the owner's words, with the stored
// line kept as Recorded; the grant itself is never rewritten.
func TestOldEventsReadInTheOwnersWords(t *testing.T) {
	at := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	g := &Grant{Approval: &Approval{DeviceID: "dev_ipad01", DeviceName: "tafs on iPad"}, Events: []Event{
		{At: at, What: "approved by user_2abc on device dev_ipad01"},
		{At: at, What: "disconnected by user_2abc on device dev_other9: token deleted"},
		{At: at, What: "sign-in confirmed on the phone ab12 cd34 ef56 7890"},
		{At: at, What: "signed in at mcp.stripe.com"},
		{At: at, What: "approved on tafs on iPad"},
	}}
	want := []struct{ what, recorded string }{
		{"approved on tafs on iPad", "approved by user_2abc on device dev_ipad01"},
		{"disconnected on a device signed in to KIFF: token deleted", "disconnected by user_2abc on device dev_other9: token deleted"},
		{"sign-in confirmed on the enrolled phone", "sign-in confirmed on the phone ab12 cd34 ef56 7890"},
		{"signed in at mcp.stripe.com", ""},
		{"approved on tafs on iPad", ""},
	}
	got := ownerEvents(g)
	for i, w := range want {
		if got[i].What != w.what || got[i].Recorded != w.recorded {
			t.Errorf("event %d: %q (recorded %q), want %q (recorded %q)", i, got[i].What, got[i].Recorded, w.what, w.recorded)
		}
	}
	if g.Events[0].What != "approved by user_2abc on device dev_ipad01" || g.Events[0].Recorded != "" {
		t.Fatal("the grant's own record was rewritten")
	}
}
