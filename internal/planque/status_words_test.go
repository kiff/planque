package planque

import "testing"

// Each status has one line in the owner's words, the same in every app.
func TestStatusWords(t *testing.T) {
	vendor := Request{Server: "https://mcp.stripe.com/"}
	for _, c := range []struct {
		g    Grant
		want string
	}{
		{Grant{Status: StatusConnected}, "Connected"},
		{Grant{Status: StatusPending}, "Waiting for your approval"},
		{Grant{Status: StatusPending, SignedIn: true}, "Signed in · tools waiting for your approval"},
		{Grant{Status: StatusPending, SignedIn: true, Connected: []ToolRequest{{Tool: "list_charges"}}}, "Connected · new tools waiting for your approval"},
		{Grant{Status: StatusListed, SignedIn: true}, "Signed in · the tools are being chosen"},
		{Grant{Status: StatusApproved, Request: vendor}, "Approved · waiting for the sign-in at mcp.stripe.com"},
		{Grant{Status: StatusApproved, Request: vendor, VendorName: "Stripe"}, "Approved · waiting for the sign-in at Stripe"},
		{Grant{Status: StatusSignedIn}, "Signed in · waiting for your phone to confirm"},
		{Grant{Status: StatusFailed, Reason: "the token could not be stored"}, "Not connected"},
		{Grant{Status: StatusDisconnected}, "Disconnected"},
	} {
		if got := statusWords(&c.g); got != c.want {
			t.Errorf("%s: %q, want %q", c.g.Status, got, c.want)
		}
	}
}
