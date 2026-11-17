package planque

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// Every key entry in the shipped catalog (deploy/catalog.json), the Qonto ones
// included: it pins the description kept in deploy/specs as it is, what it
// calls a read is a read in that description and what it calls a write is held
// for the owner's approval, and a grant from it asks for the key and offers
// exactly its tools. Entries are written by cmd/catalogen from a vendor's own
// description, and this is what keeps a batch of them honest.
func TestEveryShippedKeyEntryPinsItsDescriptionAndHoldsItsWrites(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	for id, entry := range entries {
		if entry.Kind != KindAPIKey {
			continue
		}
		t.Run(id, func(t *testing.T) {
			description, err := os.ReadFile("../../deploy/specs/" + path.Base(entry.API.URL))
			if err != nil {
				t.Fatalf("the description %s names isn't kept in deploy/specs: %v", entry.API.URL, err)
			}
			if !strings.HasPrefix(entry.API.URL, "https://raw.githubusercontent.com/kiff/planque/main/deploy/specs/") {
				t.Fatalf("a description is read from planque's own deploy/specs: %s", entry.API.URL)
			}
			if sum := sha256.Sum256(description); entry.DescriptionSHA256 != hex.EncodeToString(sum[:]) {
				t.Fatal("the entry doesn't pin its description as it is: regenerate the entry with the description")
			}
			snap, err := openapi.ParseToList(description, entry.API.Auth)
			if err != nil {
				t.Fatal(err)
			}
			described := map[string]openapi.Operation{}
			for _, o := range snap.Operations {
				described[o.Name] = o
			}
			var profile struct {
				Operations []struct {
					Tool   string `json:"tool"`
					Kind   string `json:"kind"`
					Effect string `json:"effect"`
				} `json:"operations"`
			}
			if err := json.Unmarshal(entry.Profile, &profile); err != nil {
				t.Fatal(err)
			}
			kinds := map[string]string{}
			effects := map[string]string{}
			for _, o := range profile.Operations {
				kinds[o.Tool], effects[o.Tool] = o.Kind, o.Effect
			}
			for _, tool := range entry.Tools {
				op, ok := described[tool.Tool]
				if !ok {
					t.Fatalf("%s isn't an operation of its description", tool.Tool)
				}
				if op.ReadOnly != tool.ReadOnly {
					t.Fatalf("%s: the entry says read_only=%v, the description says %v (%s %s)", tool.Name, tool.ReadOnly, op.ReadOnly, op.Method, op.Path)
				}
				if tool.ReadOnly != (kinds[tool.Name] == "read") {
					t.Fatalf("%s: the profile's kind %q doesn't match its tool", tool.Name, kinds[tool.Name])
				}
				if !tool.ReadOnly && effects[tool.Name] == "" {
					t.Fatalf("%s: a write that no effect holds for the owner", tool.Name)
				}
			}

			// A grant from the entry asks for the key and offers exactly its tools.
			e := newInProcessEnv(t)
			e.enrol()
			e.s.CatalogEntries = entries
			e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != entry.API.URL {
					t.Errorf("fetched %s", r.URL)
				}
				return profileResponse(200, string(description)), nil
			})}
			ask := map[string]string{"agent": "agt_taf"}
			switch {
			case len(entry.API.Bases) > 0:
				ask["base"] = entry.API.Bases[len(entry.API.Bases)-1] // the last region, as a customer may pick
			case accountTemplate(entry.API.Base):
				ask["account"] = "acme"
			}
			st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/"+id+"/grants", ask, nil, "k_setup")
			if st != 201 || out["needs_key"] != true {
				t.Fatalf("key grant: %d %v", st, out)
			}
			g, err := e.s.loadGrant(t.Context(), tenant, out["grant"].(map[string]any)["id"].(string))
			if err != nil {
				t.Fatal(err)
			}
			if len(g.Request.Tools) != len(entry.Tools) {
				t.Fatalf("tools: %+v", g.Request.Tools)
			}
		})
	}
}

// addressName is a property that names who else a message goes to.
var addressName = regexp.MustCompile(`(?i)^(cc|bcc|to|recipients?)$|^(cc|bcc|to)[_-]|[_-](cc|bcc|recipients?)$|^(cc|bcc)[_-]?(emails?|addresses)$`)

// toolArguments are the arguments a shipped tool takes as a taf sees them:
// the entry's own list when it pins values, else the description's input.
func toolArguments(t *testing.T, entry CatalogEntry, tool ToolRequest) map[string]bool {
	t.Helper()
	got := map[string]bool{}
	if len(tool.Arguments) > 0 || len(tool.Fixed) > 0 {
		for name := range tool.Arguments {
			got[name] = true
		}
		return got
	}
	description, err := os.ReadFile("../../deploy/specs/" + path.Base(entry.API.URL))
	if err != nil {
		t.Fatal(err)
	}
	snap, err := openapi.ParseToList(description, entry.API.Auth)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range snap.Operations {
		if o.Name != tool.Tool {
			continue
		}
		var schema struct {
			Properties map[string]struct{} `json:"properties"`
		}
		if err := json.Unmarshal(o.InputSchema, &schema); err != nil {
			t.Fatal(err)
		}
		for name := range schema.Properties {
			got[name] = true
		}
	}
	return got
}

// No shipped write can copy an outside address, and no generated entry says
// it is proven: a held reply does what its title says, and a generator never
// marks its own work tested.
func TestNoShippedWriteCopiesAnAddressAndNoEntryClaimsProof(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	for id, entry := range entries {
		var profile struct {
			TestProven bool `json:"test_proven"`
		}
		if err := json.Unmarshal(entry.Profile, &profile); err != nil || profile.TestProven {
			t.Errorf("%s: test_proven must be false in a shipped entry", id)
		}
		if entry.Kind != KindAPIKey {
			continue
		}
		for _, tool := range entry.Tools {
			if tool.ReadOnly {
				continue
			}
			for name := range toolArguments(t, entry, tool) {
				if addressName.MatchString(name) {
					t.Errorf("%s: write %s takes %q, an address list a held write could copy", id, tool.Name, name)
				}
			}
		}
	}
}

// The three writes kiff-qa found doing more than their titles: Intercom's
// reply as the customer, Freshdesk's reply copying an address, Freshdesk's
// note made public. Each comes back as a failure here.
func TestTheReplyAndNoteWritesDoWhatTheirTitlesSay(t *testing.T) {
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, nil)
	if err != nil {
		t.Fatal(err)
	}
	find := func(id, name string) (CatalogEntry, ToolRequest) {
		t.Helper()
		for _, tool := range entries[id].Tools {
			if tool.Name == name {
				return entries[id], tool
			}
		}
		t.Fatalf("%s has no tool %s", id, name)
		return CatalogEntry{}, ToolRequest{}
	}
	// Intercom: the admin's own comment, never the customer's, with no attachment.
	entry, reply := find("intercom-key@1", "intercom_reply")
	var typ, kind string
	_ = json.Unmarshal(reply.Fixed["type"], &typ)
	_ = json.Unmarshal(reply.Fixed["message_type"], &kind)
	if typ != "admin" || kind != "comment" {
		t.Fatalf("intercom_reply must pin type admin and message_type comment: %v", reply.Fixed)
	}
	args := toolArguments(t, entry, reply)
	for name := range args {
		if !map[string]bool{"conversation_id": true, "admin_id": true, "body": true}[name] {
			t.Errorf("intercom_reply takes %q: only the conversation, the admin and the words", name)
		}
	}
	if !reply.Arguments["admin_id"].Required || !reply.Arguments["body"].Required {
		t.Errorf("intercom_reply must require admin_id and body: %v", reply.Arguments)
	}
	// Freshdesk: a reply goes to the ticket's requester alone; a note stays private.
	entry, freshReply := find("freshdesk-key@1", "freshdesk_reply")
	for name := range toolArguments(t, entry, freshReply) {
		if name != "id" && name != "body" {
			t.Errorf("freshdesk_reply takes %q: only the ticket and the words", name)
		}
	}
	entry, note := find("freshdesk-key@1", "freshdesk_add_note")
	for name := range toolArguments(t, entry, note) {
		if name != "id" && name != "body" {
			t.Errorf("freshdesk_add_note takes %q: a note is private by default and has no way to be made public", name)
		}
	}
}

// A look-alike account template or a foreign host among an entry's regions
// is not on the vendor's hosts.
func TestALookAlikeTemplateOrRegionIsNotOnTheVendorsHosts(t *testing.T) {
	hosts := vendorHosts["Freshdesk"]
	if !onVendorHost(sampleBase("https://{account}.freshdesk.com/api/v2"), hosts) {
		t.Fatal("the real template must be on the vendor's hosts")
	}
	for _, template := range []string{"https://{account}.freshdesk.com.evil.example/api/v2", "https://{account}.freshdesk.co/api/v2", "https://{account}.evil.example/api/v2"} {
		if onVendorHost(sampleBase(template), hosts) {
			t.Errorf("%s is not Freshdesk's", template)
		}
	}
	regions := vendorHosts["Intercom"]
	for _, base := range []string{"https://api.eu.intercom.io", "https://api.au.intercom.io"} {
		if !onVendorHost(base, regions) {
			t.Errorf("%s must be one of Intercom's regions", base)
		}
	}
	for _, base := range []string{"https://api.eu.intercom.io.evil.example", "https://api.intercom.io:444", "http://api.intercom.io", "https://evil.example"} {
		if onVendorHost(base, regions) {
			t.Errorf("%s is not one of Intercom's regions", base)
		}
	}
}
