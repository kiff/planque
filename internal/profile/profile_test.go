package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/jcs"
)

// These are KIFF's synthetic kiff.profile.v1 / kiff.words.v1 vectors.
func TestContractVectors(t *testing.T) {
	files, err := filepath.Glob("testdata/*.expected.json")
	if err != nil || len(files) != 5 {
		t.Fatal("missing contract vectors")
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			raw, err := os.ReadFile(strings.Replace(file, ".expected.json", ".json", 1))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Parse(raw); err != nil {
				t.Fatal(err)
			}
			expected, _ := os.ReadFile(file)
			var want struct {
				ProfileSHA256 string          `json:"profile_sha256"`
				Words         json.RawMessage `json:"words"`
				WordsSHA256   string          `json:"words_sha256"`
			}
			if err := json.Unmarshal(expected, &want); err != nil {
				t.Fatal(err)
			}
			digest, err := Digest(raw)
			if err != nil || digest != want.ProfileSHA256 {
				t.Fatalf("profile digest: %s %v", digest, err)
			}
			words, err := jcs.SHA256(want.Words)
			if err != nil || words != want.WordsSHA256 {
				t.Fatalf("words digest: %s %v", words, err)
			}
		})
	}
}

func TestProfileRefusesUnknownAndTrailingFields(t *testing.T) {
	raw, _ := os.ReadFile("testdata/holded-invoices.json")
	for _, changed := range []string{strings.Replace(string(raw), `"v":`, `"unknown":true,"v":`, 1), string(raw) + " null", string(raw) + " {}"} {
		if _, err := Parse([]byte(changed)); err == nil {
			t.Fatal("accepted an ambiguous or unknown contract value")
		}
	}
}

// An operation may say where the owner looks at the vendor's record of its
// resource (record_url, reviewed catalog only) and which read gives the
// resource's own name (label_read); the id is path-escaped into the address.
func TestRecordURLAndLabelRead(t *testing.T) {
	raw, _ := os.ReadFile("testdata/stripe-refunds.json")
	with := func(changes map[string]any, source string) ([]byte, error) {
		var p map[string]any
		if err := json.Unmarshal(raw, &p); err != nil {
			t.Fatal(err)
		}
		if source != "" {
			p["source"] = source
			delete(p, "catalog_entry")
		}
		ops := p["operations"].([]any)
		for _, o := range ops {
			op := o.(map[string]any)
			if op["tool"] == "create_refund" {
				for k, v := range changes {
					op[k] = v
				}
			}
		}
		out, _ := json.Marshal(p)
		_, err := Parse(out)
		return out, err
	}
	good := map[string]any{"record_url": "https://dashboard.stripe.com/payments/{resource}", "label_read": map[string]any{"tool": "get_payment_intent", "args": map[string]string{"payment_intent": "$args.payment_intent"}, "path": "description"}}
	if _, err := with(good, ""); err != nil {
		t.Fatalf("a record_url and a label_read on a catalog profile: %v", err)
	}
	for name, c := range map[string]struct {
		changes map[string]any
		source  string
		want    string
	}{
		"a drafted profile's record_url":  {map[string]any{"record_url": "https://dashboard.stripe.com/payments/{resource}"}, "drafted", "reviewed catalog only"},
		"http":                            {map[string]any{"record_url": "http://dashboard.stripe.com/payments/{resource}"}, "", "https"},
		"credentials in it":               {map[string]any{"record_url": "https://user:pw@dashboard.stripe.com/{resource}"}, "", "no credentials"},
		"two placeholders":                {map[string]any{"record_url": "https://dashboard.stripe.com/{resource}/{resource}"}, "", "at most once"},
		"another brace":                   {map[string]any{"record_url": "https://dashboard.stripe.com/{account}/payments"}, "", "only placeholder"},
		"the placeholder as the host":     {map[string]any{"record_url": "https://{resource}/payments"}, "", "path"},
		"the placeholder in the query":    {map[string]any{"record_url": "https://dash.example.com/resources/view?id={resource}"}, "", "path"},
		"the placeholder in the fragment": {map[string]any{"record_url": "https://dash.example.com/resource#{resource}"}, "", "path"},
		"a label_read that is a write":    {map[string]any{"label_read": map[string]any{"tool": "create_refund", "path": "id"}}, "", "is a write"},
		"a label_read's bad path":         {map[string]any{"label_read": map[string]any{"tool": "get_payment_intent", "path": "../x"}}, "", "dotted path"},
	} {
		if _, err := with(c.changes, c.source); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v, want %q", name, err, c.want)
		}
	}
	// A label_read needs a resource: list_contacts in another vector has none.
	raw2, _ := os.ReadFile("testdata/holded-invoices.json")
	var p map[string]any
	_ = json.Unmarshal(raw2, &p)
	for _, o := range p["operations"].([]any) {
		if op := o.(map[string]any); op["tool"] == "list_contacts" {
			op["label_read"] = map[string]any{"tool": "list_contacts", "path": "name"}
		}
	}
	out, _ := json.Marshal(p)
	if _, err := Parse(out); err == nil || !strings.Contains(err.Error(), "needs a resource") {
		t.Errorf("a label_read on an operation without a resource: %v", err)
	}
	// The id is path-escaped into the address: it can't leave the vendor's host or path.
	for id, want := range map[string]string{
		"pi_3Abc":          "https://dashboard.stripe.com/payments/pi_3Abc",
		"../../admin":      "https://dashboard.stripe.com/payments/..%2F..%2Fadmin",
		"x?y=1#z":          "https://dashboard.stripe.com/payments/x%3Fy=1%23z",
		"a@evil.example/b": "https://dashboard.stripe.com/payments/a@evil.example%2Fb",
	} {
		if got := RecordURL("https://dashboard.stripe.com/payments/{resource}", id); got != want {
			t.Errorf("%q: %s, want %s", id, got, want)
		}
	}
	if got := RecordURL("https://dashboard.stripe.com/payments", "pi_1"); got != "https://dashboard.stripe.com/payments" {
		t.Errorf("no placeholder: %s", got)
	}
	// An id a browser would read as a step up the path, or none, gives no link.
	for _, id := range []string{"", ".", ".."} {
		if got := RecordURL("https://dashboard.stripe.com/payments/{resource}", id); got != "" {
			t.Errorf("%q: %s, want no link", id, got)
		}
	}
}

// A refund's case block names reads of the same profile whose only
// reference is the call's charge, finds KIFF's own refunds by the
// operation's marker, says where a dispute shows, and may name the app
// event that carries a refund. Its reads are renamed with the operation's.
func TestACaseBlock(t *testing.T) {
	raw, err := os.ReadFile("testdata/stripe-refunds.json")
	if err != nil {
		t.Fatal(err)
	}
	good := func() map[string]any {
		return map[string]any{
			"kind":         "refund",
			"charge_arg":   "payment_intent",
			"account_read": map[string]any{"tool": "get_account", "path": "id"},
			"charge_read": map[string]any{"tool": "get_payment_intent", "args": map[string]any{"payment_intent": "$args.payment_intent"},
				"currency": "currency", "captured": "amount_received", "disputed": "latest_charge.disputed"},
			"refunds_read": map[string]any{"tool": "list_refunds", "args": map[string]any{"payment_intent": "$args.payment_intent"},
				"list": "data", "id": "id", "amount": "amount", "status": "status", "marker": "metadata.kiff_op"},
			"refund_event": map[string]any{"event": "STRIPE_REFUND_UPDATED_TRIGGER", "object": "object", "charge": "payment_intent"},
		}
	}
	with := func(edit func(c, op map[string]any)) error {
		var p map[string]any
		if err := json.Unmarshal(raw, &p); err != nil {
			t.Fatal(err)
		}
		ops := p["operations"].([]any)
		p["operations"] = append(ops, map[string]any{"tool": "get_account", "title": "Look up the account", "kind": "read"})
		for _, o := range ops {
			if op := o.(map[string]any); op["tool"] == "create_refund" {
				c := good()
				edit(c, op)
				op["case"] = c
			}
		}
		b, _ := json.Marshal(p)
		_, err := Parse(b)
		return err
	}
	if err := with(func(map[string]any, map[string]any) {}); err != nil {
		t.Fatalf("a good case block: %v", err)
	}
	if err := with(func(c, _ map[string]any) { delete(c, "refund_event") }); err != nil {
		t.Fatalf("without a refund event: %v", err)
	}
	for name, edit := range map[string]func(c, op map[string]any){
		"another kind":      func(c, _ map[string]any) { c["kind"] = "coupon" },
		"no marker":         func(_, op map[string]any) { delete(op, "marker") },
		"another marker":    func(c, _ map[string]any) { c["refunds_read"].(map[string]any)["marker"] = "metadata.other" },
		"a write as a read": func(c, _ map[string]any) { c["account_read"].(map[string]any)["tool"] = "create_refund" },
		"an unknown read":   func(c, _ map[string]any) { c["charge_read"].(map[string]any)["tool"] = "get_charge" },
		"another reference": func(c, _ map[string]any) {
			c["refunds_read"].(map[string]any)["args"] = map[string]any{"x": "$args.amount"}
		},
		"a result reference": func(c, _ map[string]any) {
			c["charge_read"].(map[string]any)["args"] = map[string]any{"x": "$result.id"}
		},
		"no dispute":       func(c, _ map[string]any) { delete(c["charge_read"].(map[string]any), "disputed") },
		"a bad event slug": func(c, _ map[string]any) { c["refund_event"].(map[string]any)["event"] = "two words" },
		"a bad event path": func(c, _ map[string]any) { c["refund_event"].(map[string]any)["object"] = "data[x]" },
		"not moving money": func(_, op map[string]any) { op["effect"] = "records_money" },
		"an unknown field": func(c, _ map[string]any) { c["claim"] = true },
	} {
		if err := with(edit); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}
