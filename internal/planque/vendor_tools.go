package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"time"
	"unicode/utf8"

	"github.com/kiff/planque/internal/store"
)

// The vendor's own tool definitions, as its MCP server last listed them to a
// grant: kept beside the grant, before planque filters anything, so a tool
// planque doesn't serve can be checked against what the vendor really says
// (its hints, its input and output schemas, its description) without anyone
// handling the owner's token outside planque. Step 0 of the policy move
// reads list_bank_accounts this way, and step 9 passes this list through.
//
// Definitions only: what tools/list answers, never a call's result, a header
// or a token. Nothing serves it to a taf; it's read from the store.

// maxVendorTools bounds the definitions kept for one grant, well under
// DynamoDB's 400 KB item. What fits first is every tool's name, title,
// annotations and input schema, which decide what a tool is and how it's
// asked: then descriptions, then output schemas, as long as there's room.
// Only a tool whose name didn't fit is counted as left.
const maxVendorTools = 256 << 10

// maxVendorDescription is the longest description kept for one tool: past
// it, the description is cut, and the record says so.
const maxVendorDescription = 2 << 10

// vendorToolsKept is how long a listing is kept after its last read.
const vendorToolsKept = 30 * 24 * time.Hour

func skVendorTools(grant string) string { return "VENDORTOOLS#" + grant + "#" }

// vendorToolsRecord is one listing of a grant's vendor tools.
type vendorToolsRecord struct {
	Grant  string    `json:"grant"`
	Server string    `json:"server"`
	ReadAt time.Time `json:"read_at"`
	Tools  []mcpTool `json:"tools"`
	// Cut names, by tool, the parts not kept whole: "description" (cut at
	// maxVendorDescription, or left out) and "output_schema" (left out).
	Cut map[string][]string `json:"cut,omitempty"`
	// Left is how many tools the listing had whose name didn't fit.
	Left int `json:"left,omitempty"`
}

// keepVendorTools keeps listed, as the vendor sent it, for g. Best effort:
// a listing that can't be kept is logged, and the grant goes on as before.
func (s *Server) keepVendorTools(ctx context.Context, g *Grant, listed []mcpTool) {
	if len(listed) == 0 {
		return
	}
	now := s.now().UTC()
	rec := vendorToolsRecord{Grant: g.ID, Server: g.Request.Server, ReadAt: now}
	rec.Tools, rec.Cut, rec.Left = boundVendorTools(listed)
	it := store.Item{PK: tenantPK(g.Tenant), SK: skVendorTools(g.ID), TTL: now.Add(vendorToolsKept).Unix()}
	if err := putDoc(ctx, s.KV, it, rec, store.Cond{}); err != nil {
		s.log().WarnContext(ctx, "vendor tools not kept", slog.String("grant", g.ID), slog.String("error", err.Error()))
	}
}

// boundVendorTools is what's kept of listed within maxVendorTools, in list
// order: every tool's name, title, annotations and input schema first; then
// each description, cut at maxVendorDescription; then each output schema.
// cut names what wasn't kept whole, and left counts the tools whose name
// didn't fit at all.
func boundVendorTools(listed []mcpTool) (kept []mcpTool, cut map[string][]string, left int) {
	size := func(t mcpTool) int {
		b, err := json.Marshal(t)
		if err != nil {
			return maxVendorTools + 1
		}
		return len(b)
	}
	used := 0
	for i, t := range listed {
		base := mcpTool{Name: t.Name, Title: t.Title, Annotations: t.Annotations, InputSchema: t.InputSchema}
		n := size(base)
		if used+n > maxVendorTools {
			left = len(listed) - i
			break
		}
		used += n
		kept = append(kept, base)
	}
	cut = map[string][]string{}
	mark := func(name, part string) {
		if !slices.Contains(cut[name], part) {
			cut[name] = append(cut[name], part)
		}
	}
	grow := func(i int, next mcpTool) bool {
		d := size(next) - size(kept[i])
		if used+d > maxVendorTools {
			return false
		}
		used += d
		kept[i] = next
		return true
	}
	for i := range kept {
		t := listed[i]
		if t.Description == "" {
			continue
		}
		next := kept[i]
		next.Description = t.Description
		if len(next.Description) > maxVendorDescription {
			next.Description = cutUTF8(next.Description, maxVendorDescription)
			mark(t.Name, "description")
		}
		if !grow(i, next) {
			mark(t.Name, "description")
		}
	}
	for i := range kept {
		t := listed[i]
		if len(t.OutputSchema) == 0 {
			continue
		}
		next := kept[i]
		next.OutputSchema = t.OutputSchema
		if !grow(i, next) {
			mark(t.Name, "output_schema")
		}
	}
	if len(cut) == 0 {
		cut = nil
	}
	return kept, cut, left
}

// cutUTF8 is s cut to at most n bytes, never inside a character.
func cutUTF8(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n]
}
