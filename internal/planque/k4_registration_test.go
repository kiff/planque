package planque

import (
	"encoding/json"
	"strings"
	"testing"
)

// KIFF classes each call by the tool's app and its own name (K4), so a tool's
// registration names its app as a lowercase slug, the one KIFF's class tables
// use: Qonto's and Stripe's own servers by their name, a managed app by its
// toolkit. A name planque can't make a slug of is left out, never guessed.
func TestEachRegistrationNamesItsAppAsKIFFsTablesDo(t *testing.T) {
	e := newEnv(t)
	cases := []struct {
		name string
		g    *Grant
		want string
	}{
		{"Qonto's own server", &Grant{Request: Request{Server: qontoServer}}, "qonto"},
		{"Stripe's own server", &Grant{Request: Request{Server: stripeServer}}, "stripe"},
		{"Qonto's entry", &Grant{Request: Request{CatalogEntry: QontoEntry, Server: qontoServer}}, "qonto"},
		{"a managed app", &Grant{Request: Request{Connector: &ConnectorRequest{Provider: "p", Toolkit: "GMAIL"}}}, "gmail"},
		{"another server", &Grant{Request: Request{Server: "https://mcp.vendor.example/mcp"}}, ""},
	}
	for _, c := range cases {
		if got := e.s.appSlug(c.g); got != c.want {
			t.Errorf("%s: app %q, want %q", c.name, got, c.want)
		}
	}
	for in, want := range map[string]string{"Qonto": "qonto", "Google Sheets": "google-sheets", "google_drive": "google-drive",
		"": "", "Café": "", "-x": "", strings.Repeat("a", 65): ""} {
		if got := slugOf(in); got != want {
			t.Errorf("slugOf(%q) = %q, want %q", in, got, want)
		}
	}
}

// The vendor's destructiveHint is carried from its listing to KIFF, as the
// vendor says it and nothing else: KIFF only raises a class on it.
func TestTheVendorsDestructiveHintReachesKIFFAsListed(t *testing.T) {
	listed := []mcpTool{
		{Name: "delete_client", Annotations: json.RawMessage(`{"readOnlyHint":false,"destructiveHint":true}`)},
		{Name: "list_clients", Annotations: json.RawMessage(`{"readOnlyHint":true,"destructiveHint":false}`)},
		{Name: "update_client"},
	}
	g := &Grant{Request: Request{Server: qontoServer}, Audience: "https://relay.example/g1"}
	want := map[string]bool{"delete_client": true, "list_clients": false, "update_client": false}
	for _, l := range listed {
		gt, why := granted(g, ToolRequest{Tool: l.Name, Name: l.Name}, listed)
		if why != "" {
			t.Fatalf("%s: %s", l.Name, why)
		}
		if gt.DestructiveHint != want[l.Name] {
			t.Errorf("%s: destructive %v, want %v", l.Name, gt.DestructiveHint, want[l.Name])
		}
		conn := kiffConnection(g, gt, "cred", "qonto")
		if conn.App != "qonto" || conn.DestructiveHint != want[l.Name] {
			t.Errorf("%s: registered as app %q, destructive %v", l.Name, conn.App, conn.DestructiveHint)
		}
	}
}

// With no app and no hint, the registration is byte for byte what it was
// before K4: neither field is sent.
func TestARegistrationWithNoAppOrHintSendsNeither(t *testing.T) {
	g := &Grant{Request: Request{Server: "https://mcp.vendor.example/mcp"}, Audience: "https://relay.example/g1"}
	b, err := json.Marshal(kiffConnection(g, GrantedTool{ToolRequest: ToolRequest{Tool: "do_it", Name: "do_it"}}, "cred", ""))
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); strings.Contains(s, `"app"`) || strings.Contains(s, `"destructive_hint"`) {
		t.Fatalf("sent a field KIFF's gateway may not know yet: %s", s)
	}
}
