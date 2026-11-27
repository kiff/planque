package planque

import (
	"testing"

	"github.com/kiff/planque/internal/oauth"
)

// What disconnecting does at the vendor, as the apps are told it: a sign-in
// is revoked only where the vendor's metadata names a revocation endpoint (so
// Qonto's isn't), a key the owner made stays at the vendor until the owner
// deletes it there, a managed sign-in is revoked through its platform, and a
// browser says nothing.
func TestVendorRevokeSaysWhatDisconnectingDoesAtTheVendor(t *testing.T) {
	signIn := func(revocation string) *Grant {
		return &Grant{Vendor: oauth.Metadata{Token: "https://vendor.example/token", Revocation: revocation}}
	}
	cases := []struct {
		name string
		g    *Grant
		want string
	}{
		{"a sign-in whose vendor revokes", signIn("https://vendor.example/revoke"), revokeAtVendor},
		{"a sign-in whose vendor offers no revocation", signIn(""), revokeNotOffered},
		{"a key", &Grant{Request: Request{Kind: KindAPIKey}}, revokeKeyStays},
		{"a managed sign-in", &Grant{Request: Request{Connector: &ConnectorRequest{Provider: "composio", Toolkit: "gmail"}}}, revokeAtVendor},
		{"a browser", &Grant{Request: Request{Kind: KindBrowser}}, ""},
		{"nothing held at a vendor yet", &Grant{}, ""},
	}
	for _, c := range cases {
		if got := vendorRevokeOf(c.g); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}
