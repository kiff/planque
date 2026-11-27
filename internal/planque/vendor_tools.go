package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

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
// DynamoDB's 400 KB item: past it, the rest are counted, not kept.
const maxVendorTools = 256 << 10

// vendorToolsKept is how long a listing is kept after its last read.
const vendorToolsKept = 30 * 24 * time.Hour

func skVendorTools(grant string) string { return "VENDORTOOLS#" + grant + "#" }

// vendorToolsRecord is one listing of a grant's vendor tools.
type vendorToolsRecord struct {
	Grant  string    `json:"grant"`
	Server string    `json:"server"`
	ReadAt time.Time `json:"read_at"`
	Tools  []mcpTool `json:"tools"`
	// Left is how many tools the listing had past maxVendorTools.
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
	size := 0
	for i, t := range listed {
		b, err := json.Marshal(t)
		if err != nil || size+len(b) > maxVendorTools {
			rec.Left = len(listed) - i
			break
		}
		size += len(b)
		rec.Tools = append(rec.Tools, t)
	}
	it := store.Item{PK: tenantPK(g.Tenant), SK: skVendorTools(g.ID), TTL: now.Add(vendorToolsKept).Unix()}
	if err := putDoc(ctx, s.KV, it, rec, store.Cond{}); err != nil {
		s.log().WarnContext(ctx, "vendor tools not kept", slog.String("grant", g.ID), slog.String("error", err.Error()))
	}
}
