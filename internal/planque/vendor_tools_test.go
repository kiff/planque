package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
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

// A listing past the bound keeps what fits and counts the rest.
func TestAVendorListingPastItsBoundCountsWhatItLeftOut(t *testing.T) {
	e := newEnv(t)
	g := &Grant{ID: "grt_big", Tenant: tenant, Request: Request{Server: "https://mcp.vendor.example/mcp"}}
	big := strings.Repeat("x", 64<<10)
	var listed []mcpTool
	for i := 0; i < 8; i++ {
		listed = append(listed, mcpTool{Name: "tool_" + string(rune('a'+i)), Description: big})
	}
	e.s.keepVendorTools(context.Background(), g, listed)
	var rec vendorToolsRecord
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skVendorTools(g.ID), &rec); err != nil {
		t.Fatal(err)
	}
	if len(rec.Tools) == 0 || len(rec.Tools)+rec.Left != len(listed) || rec.Left == 0 {
		t.Fatalf("kept %d, left %d of %d", len(rec.Tools), rec.Left, len(listed))
	}
}
