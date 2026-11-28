package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

// The vendor's own tool definitions are kept beside the grant as the server
// listed them, before planque filters anything: annotations and output
// schemas included, and never a token.
func TestTheVendorsOwnToolDefinitionsAreKeptAsListed(t *testing.T) {
	e := newEnv(t)
	e.vendor.ChangeTools = func(tools []map[string]any) []map[string]any {
		return append(tools, map[string]any{
			"name":         "list_bank_accounts",
			"description":  "List the organization's bank accounts, with each one's balance.",
			"inputSchema":  map[string]any{"type": "object", "properties": map[string]any{"per_page": map[string]any{"type": "integer"}}},
			"annotations":  map[string]any{"readOnlyHint": true, "destructiveHint": false},
			"outputSchema": map[string]any{"type": "object", "properties": map[string]any{"bank_accounts": map[string]any{"type": "array"}}},
		})
	}
	id := e.listedGrant()

	var rec vendorToolsRecord
	it, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skVendorTools(id), &rec)
	if err != nil {
		t.Fatalf("no listing kept: %v", err)
	}
	if rec.Grant != id || rec.Server != e.vendor.MCPURL() || rec.ReadAt.IsZero() || rec.Left != 0 || it.TTL == 0 {
		t.Fatalf("record: %+v (ttl %d)", rec, it.TTL)
	}
	var kept *mcpTool
	for i := range rec.Tools {
		if rec.Tools[i].Name == "list_bank_accounts" {
			kept = &rec.Tools[i]
		}
	}
	if kept == nil || len(rec.Tools) < 2 {
		t.Fatalf("the vendor's list isn't kept whole: %+v", rec.Tools)
	}
	var hints map[string]any
	if json.Unmarshal(kept.Annotations, &hints) != nil || hints["readOnlyHint"] != true || hints["destructiveHint"] != false {
		t.Fatalf("annotations: %s", kept.Annotations)
	}
	if !strings.Contains(string(kept.OutputSchema), "bank_accounts") || !strings.Contains(string(kept.InputSchema), "per_page") ||
		!strings.Contains(kept.Description, "balance") {
		t.Fatalf("definition: %+v", kept)
	}
	// Definitions only: nothing of the sign-in.
	if doc := strings.ToLower(string(it.Doc)); strings.Contains(doc, "token") || strings.Contains(doc, "bearer") || strings.Contains(doc, "authorization") {
		t.Fatalf("the listing holds something of the sign-in: %s", it.Doc)
	}
}

// A listing whose names and input schemas alone are past the bound keeps the
// tools that fit, whole in what decides them, and counts the rest as left.
func TestAVendorListingPastItsBoundCountsWhatItLeftOut(t *testing.T) {
	e := newEnv(t)
	g := &Grant{ID: "grt_big", Tenant: tenant, Request: Request{Server: "https://mcp.vendor.example/mcp"}}
	schema := json.RawMessage(`{"type":"object","description":"` + strings.Repeat("x", 64<<10) + `"}`)
	var listed []mcpTool
	for i := 0; i < 8; i++ {
		listed = append(listed, mcpTool{Name: "tool_" + string(rune('a'+i)), InputSchema: schema})
	}
	e.s.keepVendorTools(context.Background(), g, listed)
	var rec vendorToolsRecord
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skVendorTools(g.ID), &rec); err != nil {
		t.Fatal(err)
	}
	if len(rec.Tools) == 0 || len(rec.Tools)+rec.Left != len(listed) || rec.Left == 0 {
		t.Fatalf("kept %d, left %d of %d", len(rec.Tools), rec.Left, len(listed))
	}
	for _, kept := range rec.Tools {
		if string(kept.InputSchema) != string(schema) {
			t.Fatalf("%s's input schema wasn't kept whole", kept.Name)
		}
	}
}

// Qonto's sign-in lists 61 tools whose definitions run 20 to 36 KB, mostly
// output schemas: re-test 6 kept the first 31 whole and left 30 out, so
// list_bank_accounts, the read step 0 checks, wasn't kept. Every tool's name,
// annotations and input schema is kept now; descriptions are cut at 2 KB and
// output schemas fill what's left, and the record names what wasn't kept whole.
func TestEveryToolOfAQontoSizedListingIsKeptByNameAnnotationsAndInputSchema(t *testing.T) {
	e := newEnv(t)
	g := &Grant{ID: "grt_qonto", Tenant: tenant, Request: Request{Server: qontoServer}}
	pad := func(n int) string { return strings.Repeat("y", n) }
	var listed []mcpTool
	for i := 0; i < 61; i++ {
		name := fmt.Sprintf("tool_%02d", i)
		if i == 45 {
			name = "list_bank_accounts"
		}
		listed = append(listed, mcpTool{
			Name:         name,
			Title:        "Title " + name,
			Description:  pad(3 << 10),
			InputSchema:  json.RawMessage(`{"type":"object","properties":{"q":{"type":"string","description":"` + pad(1<<10) + `"}}}`),
			Annotations:  json.RawMessage(`{"readOnlyHint":true,"destructiveHint":false}`),
			OutputSchema: json.RawMessage(`{"type":"object","description":"` + pad(24<<10) + `"}`),
		})
	}
	e.s.keepVendorTools(context.Background(), g, listed)
	var rec vendorToolsRecord
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skVendorTools(g.ID), &rec); err != nil {
		t.Fatal(err)
	}
	if rec.Left != 0 || len(rec.Tools) != len(listed) {
		t.Fatalf("kept %d, left %d of %d", len(rec.Tools), rec.Left, len(listed))
	}
	b, _ := json.Marshal(rec)
	if len(b) > maxVendorTools+(8<<10) {
		t.Fatalf("the record is %d bytes, past the bound", len(b))
	}
	outputs := 0
	for i, kept := range rec.Tools {
		want := listed[i]
		if kept.Name != want.Name || kept.Title != want.Title || string(kept.Annotations) != string(want.Annotations) || string(kept.InputSchema) != string(want.InputSchema) {
			t.Fatalf("%s isn't kept whole in its name, title, annotations and input schema", want.Name)
		}
		if len(kept.Description) == 0 || len(kept.Description) > maxVendorDescription || !slices.Contains(rec.Cut[want.Name], "description") {
			t.Fatalf("%s's description: %d bytes, cut %v", want.Name, len(kept.Description), rec.Cut[want.Name])
		}
		if len(kept.OutputSchema) > 0 {
			outputs++
			if string(kept.OutputSchema) != string(want.OutputSchema) || slices.Contains(rec.Cut[want.Name], "output_schema") {
				t.Fatalf("%s's output schema: kept %d bytes, cut %v", want.Name, len(kept.OutputSchema), rec.Cut[want.Name])
			}
		} else if !slices.Contains(rec.Cut[want.Name], "output_schema") {
			t.Fatalf("%s's output schema was left out and the record doesn't say so", want.Name)
		}
	}
	if outputs == 0 || outputs == len(listed) {
		t.Fatalf("output schemas kept: %d of %d, want what fits", outputs, len(listed))
	}
	if rec.Tools[45].Name != "list_bank_accounts" || string(rec.Tools[45].InputSchema) != string(listed[45].InputSchema) {
		t.Fatal("list_bank_accounts isn't kept")
	}
}

// A description is cut on a character's edge, never inside one.
func TestAVendorDescriptionIsCutOnACharactersEdge(t *testing.T) {
	s := strings.Repeat("é", maxVendorDescription) // two bytes each
	got := cutUTF8(s, maxVendorDescription+1)
	if !utf8.ValidString(got) || len(got) != maxVendorDescription {
		t.Fatalf("cut to %d bytes, valid %v", len(got), utf8.ValidString(got))
	}
}
