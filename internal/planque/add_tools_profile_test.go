package planque

import (
	"encoding/json"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// Adding tools to a connection whose profile the owner approved: the caller
// (Coach) sends a profile of the new operations alone, never having the
// approved profile's bytes, and planque merges it into the approved one,
// which stays as it is. Before, every addition was refused (409 "adding
// tools cannot change or remove the existing profile", or 400 "the
// profile must name every requested tool exactly once").
func TestAnAdditionsProfileIsMergedIntoTheApprovedOne(t *testing.T) {
	op := func(tool, title, kind string) profile.Operation {
		o := profile.Operation{Tool: tool, Title: title, Kind: kind, Shape: "single"}
		if kind == "write" {
			o.Effect = profile.EffectSends
			o.Resource = &profile.Resource{Arg: "inbox_id", Noun: "inbox", Outside: true}
		}
		return o
	}
	doc := func(ops ...profile.Operation) json.RawMessage {
		b, err := json.Marshal(profile.Profile{V: profile.Version, Vendor: "Mail", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted", Operations: ops})
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	list := op("list_inboxes", "List inboxes", "read")
	approved := doc(list)
	has := []ToolRequest{{Tool: "connector_a", Title: "List inboxes", Name: "list_inboxes"}}
	added := []ToolRequest{{Tool: "connector_b", Title: "Send reminder", Name: "send_reminder"}, {Tool: "connector_c", Title: "List sent messages", Name: "list_sent_messages"}}
	next := doc(op("send_reminder", "Send a reminder", "write"), op("list_sent_messages", "List sent messages", "read"))

	merged := withApprovedOperations(approved, next, added)
	if !unchangedProfile(approved, merged) {
		t.Fatalf("the approved operations changed:\n%s", merged)
	}
	p, err := profile.Parse(merged)
	if err != nil || len(p.Operations) != 3 || p.Operations[0].Tool != "list_inboxes" {
		t.Fatalf("merged: %v %+v", err, p)
	}
	if _, err := profileFields(merged, append(append([]ToolRequest{}, has...), added...)); err != nil {
		t.Fatalf("the merged profile isn't the grant's tools: %v", err)
	}

	// A full profile (approved operations repeated) is taken as given, as before.
	full := doc(list, op("send_reminder", "Send a reminder", "write"), op("list_sent_messages", "List sent messages", "read"))
	if got := withApprovedOperations(approved, full, added); string(got) != string(full) {
		t.Fatalf("a full profile was rewritten:\n%s", got)
	}
	// An operation for a tool the grant has, or one that isn't being added: as given, so the old checks refuse it.
	again := doc(op("list_inboxes", "List every inbox", "read"), op("send_reminder", "Send a reminder", "write"))
	if got := withApprovedOperations(approved, again, added); string(got) != string(again) {
		t.Fatalf("an operation replacing an approved one was merged:\n%s", got)
	}
	if got := withApprovedOperations(approved, next, added[:1]); string(got) != string(next) {
		t.Fatalf("more operations than tools added was merged:\n%s", got)
	}
	// No approved profile yet, or no profile sent: as given.
	if got := withApprovedOperations(nil, next, added); string(got) != string(next) {
		t.Fatalf("with nothing approved: %s", got)
	}
	if got := withApprovedOperations(approved, nil, added); got != nil {
		t.Fatalf("with no profile sent: %s", got)
	}
}
