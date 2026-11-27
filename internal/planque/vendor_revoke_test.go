package planque

import (
	"strings"
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

// What ending a connection did with its sign-in or key reads the same
// whichever door ended it: a key, a sign-in whose vendor offers no
// revocation, and a shared one at either kind of vendor each say what stays
// at the vendor.
func TestSignInEndedSaysWhatStaysAtTheVendor(t *testing.T) {
	noRevoke := oauth.Metadata{Token: "https://qonto.example/token"}
	revokes := oauth.Metadata{Token: "https://vendor.example/token", Revocation: "https://vendor.example/revoke"}
	cases := []struct {
		name string
		g    *Grant
		held bool
		want string
	}{
		{"kept for another taf", &Grant{keptFor: 1, Vendor: noRevoke}, true, "the sign-in stays, another taf's connection uses it"},
		{"a shared key, the last on it", &Grant{revokedShared: true, VendorName: "Moneybird", Request: Request{Kind: KindAPIKey}}, true,
			"the key it shared was deleted here, no other connection used it; it stays valid at Moneybird until the owner deletes it there"},
		{"a shared sign-in at a vendor with no revocation", &Grant{revokedShared: true, Vendor: noRevoke, VendorName: "Qonto"}, true,
			"the sign-in it shared was deleted here, no other connection used it; Qonto offers no revocation, so the sign-in ends there only when the owner removes tafs's access there"},
		{"a shared sign-in at a vendor that revokes", &Grant{revokedShared: true, Vendor: revokes}, true, "the sign-in it shared was revoked, no other connection used it"},
		{"a reused sign-in still in use", &Grant{TokenOf: "grt_first", Vendor: noRevoke}, true, "the sign-in it used stays with the connection it came from"},
		{"a key", &Grant{VendorName: "Moneybird", Request: Request{Kind: KindAPIKey}}, true, "key deleted here; it stays valid at Moneybird until the owner deletes it there"},
		{"a sign-in at a vendor with no revocation", &Grant{Vendor: noRevoke, VendorName: "Qonto"}, true,
			"token deleted here; Qonto offers no revocation, so the sign-in ends there only when the owner removes tafs's access there"},
		{"one that never signed in", &Grant{Vendor: noRevoke, VendorName: "Qonto"}, false, ""},
		{"a sign-in at a vendor that revokes", &Grant{Vendor: revokes}, true, ""},
	}
	for _, c := range cases {
		if got := signInEnded(c.g, c.held); got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// Both doors write it: the owner's Disconnect and its taf's end say the same
// of a sign-in whose vendor offers no revocation.
func TestBothDoorsSayASignInWithNoRevocationIsOnlyDeletedHere(t *testing.T) {
	const said = "offers no revocation, so the sign-in ends there only when the owner removes tafs's access there"
	noRevocation := func(e *env) *Grant {
		g := e.connect()
		stored, err := e.s.loadGrant(t.Context(), tenant, g.ID)
		if err != nil {
			t.Fatal(err)
		}
		stored.Vendor.Revocation = ""
		if err := e.s.saveGrant(t.Context(), stored); err != nil {
			t.Fatal(err)
		}
		return stored
	}
	last := func(e *env, id string) string {
		g, err := e.s.loadGrant(t.Context(), tenant, id)
		if err != nil {
			t.Fatal(err)
		}
		return g.Events[len(g.Events)-1].What
	}

	e := newEnv(t)
	g := noRevocation(e)
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, ""); st != 200 {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if what := last(e, g.ID); !strings.Contains(what, "token deleted here; ") || !strings.Contains(what, said) {
		t.Errorf("the owner's Disconnect: %q", what)
	}

	e = newEnv(t)
	e.enrol()
	e.people()
	g = noRevocation(e)
	if st, out := e.endAgent(tenant, g.Request.Agent, "s_owner", ""); st != 200 {
		t.Fatalf("end: %d %v", st, out)
	}
	if what := last(e, g.ID); !strings.Contains(what, "ended with its taf") || !strings.Contains(what, "token deleted here; ") || !strings.Contains(what, said) {
		t.Errorf("its taf's end: %q", what)
	}
}
