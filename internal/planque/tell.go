package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/store"
)

// pkTell holds one item per sign-out KIFF hasn't taken yet, keyed by when
// the vendor signed the grant out (tellSK), so the scheduled sweep sends it
// again. KIFF dedupes a resend of one sign-out, and alerts once an hour at
// most per grant.
const pkTell = "TELL"

// tellFor is how long a sign-out is told again before planque gives up: by
// then the owner has opened the taf, or will see it there.
const tellFor = 7 * 24 * time.Hour

// tellWait bounds one notice, so a call the vendor just refused isn't held
// up long by KIFF.
const tellWait = 3 * time.Second

func tellSK(at time.Time, tenant, id string) string {
	return at.UTC().Format(time.RFC3339) + "#" + tenant + "#" + id
}

type tellDoc struct {
	Reason string `json:"reason"`
}

// markTell notes that g's sign-out is to be told to KIFF, before g is saved
// signed out: a notice is never lost to a failure after the save.
func (s *Server) markTell(ctx context.Context, g *Grant, reason string) error {
	doc, _ := json.Marshal(tellDoc{Reason: reason})
	it := store.Item{PK: pkTell, SK: tellSK(g.SignedOutAt, g.Tenant, g.ID), Doc: doc, TTL: g.SignedOutAt.Add(tellFor).Unix()}
	return s.KV.Put(ctx, it, store.Cond{})
}

// tellSignedOut tells KIFF that g's vendor signed it out, and, once KIFF
// takes it, forgets the note. Best effort: a failure waits for the sweep.
func (s *Server) tellSignedOut(ctx context.Context, g *Grant, reason string) bool {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), tellWait)
	defer cancel()
	key, err := s.connectorKey(ctx, g.Tenant)
	if err == nil {
		err = s.KIFF.ConnectionSignedOut(ctx, key, kiff.SignedOut{Grant: g.ID, Agent: g.Request.Agent, Vendor: vendorName(g),
			Reason: reason, SignedOutAt: g.SignedOutAt.UTC().Format(time.RFC3339)})
	}
	if err != nil {
		s.log().WarnContext(ctx, "sign-out not told to KIFF yet", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return false
	}
	if _, _, err := s.KV.Delete(ctx, pkTell, tellSK(g.SignedOutAt, g.Tenant, g.ID)); err != nil {
		s.log().WarnContext(ctx, "told sign-out still marked", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	return true
}

// tellAgain sends the sign-outs KIFF hasn't taken. A note whose grant is
// gone, signed in again, or signed out at another time is dropped.
func (s *Server) tellAgain(ctx context.Context) (int, error) {
	items, err := s.KV.QueryTo(ctx, pkTell, "~")
	if err != nil {
		return 0, err
	}
	told := 0
	for _, it := range items {
		parts := strings.SplitN(it.SK, "#", 3)
		if len(parts) != 3 {
			_, _, _ = s.KV.Delete(ctx, pkTell, it.SK)
			continue
		}
		at, perr := time.Parse(time.RFC3339, parts[0])
		g, err := s.loadGrant(ctx, parts[1], parts[2])
		if perr != nil || errors.Is(err, errNotFound) || err == nil && (!g.signedOut() || !g.SignedOutAt.Equal(at)) {
			_, _, _ = s.KV.Delete(ctx, pkTell, it.SK)
			continue
		}
		if err != nil {
			continue
		}
		var doc tellDoc
		_ = json.Unmarshal(it.Doc, &doc)
		if s.tellSignedOut(ctx, g, doc.Reason) {
			told++
		}
	}
	return told, nil
}
