package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kiff/planque/internal/store"
)

// Deleting the account from the app (kiff/tafs#185, App Store Review
// Guideline 5.1.1(v)): tafs cloud forwards the owner's signed request,
// before KIFF's own deletion, since planque checks the signature with the
// account's connector key, which KIFF's deletion revokes.

// eraseAccount deletes everything planque keeps for the account. Asked
// again once it's done, there is no enrolment left to check a signature
// with and nothing kept, so it answers 204; anything else needs the
// signature of someone KIFF says may delete the account (account.delete,
// admins only), as KIFF's own deletion does.
func (s *Server) eraseAccount(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	items, err := s.KV.Query(r.Context(), tenantPK(tenant), "")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.signedFor(capDelete, s.erase)(w, r)
}

// erase disconnects each grant (its tools leave KIFF's gateway, and its
// vendor tokens are revoked where the vendor allows, then deleted) and
// deletes every record under the account. The enrolment, with the KIFF
// connector key, goes last, so a deletion cut short is finished by asking
// again with a new signature. A vendor or KIFF that doesn't confirm is
// named in problems; planque's own copies are deleted either way.
func (s *Server) erase(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	grants, problems, err := s.eraseAll(r.Context(), sr.tenant, sr.enrolment, sr.approval.Fingerprint)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant": sr.tenant, "erased": true, "grants": grants, "problems": problems})
}

// eraseAll disconnects each grant and deletes every record of the account,
// the enrolment last; by names who asked, for the log.
func (s *Server) eraseAll(ctx context.Context, tenant string, e Enrolment, by string) (int, []string, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return 0, nil, err
	}
	problems := []string{}
	key, keyErr := s.openConnector(ctx, e)
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil {
			continue
		}
		problems = append(problems, s.eraseGrant(ctx, &g, key, keyErr)...)
	}
	all, err := s.KV.Query(ctx, tenantPK(tenant), "")
	if err != nil {
		return 0, nil, err
	}
	for _, it := range all {
		if it.SK == skEnrolment {
			continue
		}
		if strings.HasPrefix(it.SK, "DOOR_EVENT#") {
			if _, _, err := s.KV.Delete(ctx, pendingEventsPK, strings.TrimPrefix(it.SK, "DOOR_EVENT#")); err != nil {
				return 0, nil, err
			}
			if _, _, err := s.KV.Delete(ctx, deliveryPK, strings.TrimPrefix(it.SK, "DOOR_EVENT#")); err != nil {
				return 0, nil, err
			}
		}
		if strings.HasPrefix(it.SK, "TRIGGER#") {
			var binding struct {
				Provider string `json:"provider"`
				Trigger  string `json:"trigger"`
			}
			if json.Unmarshal(it.Doc, &binding) == nil && binding.Provider != "" && binding.Trigger != "" {
				if _, _, err := s.KV.Delete(ctx, triggerPK(binding.Provider, binding.Trigger), "-"); err != nil {
					return 0, nil, err
				}
			}
		}

		if _, _, err := s.KV.Delete(ctx, it.PK, it.SK); err != nil {
			return 0, nil, err
		}
	}
	if _, _, err := s.KV.Delete(ctx, tenantPK(tenant), skEnrolment); err != nil {
		return 0, nil, err
	}
	s.log().InfoContext(ctx, "account erased", slog.String("tenant", tenant), slog.Int("grants", len(items)),
		slog.Int("problems", len(problems)), slog.String("by", by))
	return len(items), problems, nil
}

// eraseGrant takes g's tools out of KIFF's gateway, unless it was already
// disconnected or failed, and revokes its vendor tokens.
func (s *Server) eraseGrant(ctx context.Context, g *Grant, key string, keyErr error) []string {
	var problems []string
	if g.Status != StatusDisconnected && g.Status != StatusFailed && len(g.inKIFF()) > 0 {
		if keyErr != nil {
			problems = append(problems, "the KIFF connector key could not be opened")
		} else {
			for _, t := range g.inKIFF() {
				if err := s.KIFF.Disconnect(ctx, key, t.Name); err != nil {
					problems = append(problems, "KIFF kept "+t.Name+": "+err.Error())
				}
			}
		}
	}
	return append(problems, s.revoke(ctx, g)...)
}

// finishedPerCall is the most finished grants one call deletes: the rest
// answer more: true and the next call carries on (the Eval's account held 960
// grants and 8,004 operations, too much for one function's time).
var finishedPerCall = 200

// eraseFinished deletes the account's finished grants (disconnected or failed,
// and nothing else) with each one's own records: its operations, snapshots,
// live row and due-index row. It never touches a live grant, the enrolment, or
// a finished grant that still holds a token or a provider access not confirmed
// disabled: that one is skipped and counted, so a record never goes while a
// secret of it remains. A grant's own record goes last, so a call cut short is
// finished by the next one. One account event says who signed and the counts.
func (s *Server) eraseFinished(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	ctx := r.Context()
	items, err := s.KV.Query(ctx, tenantPK(sr.tenant), "GRANT#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var drop []Grant
	skipped := 0
	for _, it := range items {
		var g Grant
		if json.Unmarshal(it.Doc, &g) != nil || g.ID == "" || skGrant(g.ID) != it.SK || g.live() {
			continue
		}
		if ok, err := s.holdsNothing(ctx, sr.tenant, &g); err != nil {
			s.fail(w, r, err)
			return
		} else if !ok {
			skipped++
			continue
		}
		// A browser grant's screenshots are kept outside the table: they go
		// first, and the grant is skipped if they can't be dropped.
		if g.Request.Kind == KindBrowser && s.Evidence != nil {
			if err := s.Evidence.Drop(ctx, sr.tenant, g.ID); err != nil {
				skipped++
				continue
			}
		}
		drop = append(drop, g)
	}
	more := len(drop) > finishedPerCall
	if more {
		drop = drop[:finishedPerCall]
	}
	gone := map[string]bool{}
	for _, g := range drop {
		gone[g.ID] = true
	}
	// The due index is dropped first for these grants: the sweeper never
	// meets a grant that's gone.
	if len(gone) > 0 {
		due, err := s.KV.Query(ctx, pkDue, "")
		if err != nil {
			s.fail(w, r, err)
			return
		}
		for _, it := range due {
			if _, tenant, id, ok := parseDue(it.SK); ok && tenant == sr.tenant && gone[id] {
				if _, _, err := s.KV.Delete(ctx, pkDue, it.SK); err != nil {
					s.fail(w, r, err)
					return
				}
			}
		}
	}
	ops := 0
	for _, g := range drop {
		for _, prefix := range []string{"OP#" + g.ID + "#", "SNAPSHOT#" + g.ID + "#", skVendorTools(g.ID)} {
			rows, err := s.KV.Query(ctx, tenantPK(sr.tenant), prefix)
			if err != nil {
				s.fail(w, r, err)
				return
			}
			for _, it := range rows {
				if _, _, err := s.KV.Delete(ctx, it.PK, it.SK); err != nil {
					s.fail(w, r, err)
					return
				}
				if strings.HasPrefix(prefix, "OP#") {
					ops++
				}
			}
		}
		if _, _, err := s.KV.Delete(ctx, tenantPK(sr.tenant), skLive(g.ID)); err != nil {
			s.fail(w, r, err)
			return
		}
		if _, _, err := s.KV.Delete(ctx, tenantPK(sr.tenant), skGrant(g.ID)); err != nil {
			s.fail(w, r, err)
			return
		}
	}
	at := s.now().UTC()
	event := map[string]any{"what": "finished grants deleted", "at": at.Format(time.RFC3339), "by": sr.approval.Fingerprint,
		"grants": len(drop), "operations": ops, "skipped": skipped, "more": more}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(sr.tenant), SK: "ACCOUNT_EVENT#" + at.Format(time.RFC3339Nano), Ver: 1}, event, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(ctx, "finished grants deleted", slog.String("tenant", sr.tenant), slog.Int("grants", len(drop)),
		slog.Int("operations", ops), slog.Int("skipped", skipped), slog.Bool("more", more), slog.String("by", sr.approval.Fingerprint))
	writeJSON(w, http.StatusOK, map[string]any{"grants": len(drop), "operations": ops, "skipped": skipped, "more": more,
		"grants_before": len(items), "grants_after": len(items) - len(drop)})
}

// holdsNothing says a finished grant of the signed account has no token,
// renewal, browser session, held browser sign-in or checkout row left, and no
// provider access that isn't confirmed disabled: only then may its record go.
// A record that names another account, or none, is never the signed one's to
// delete: it's skipped, whatever its rows say.
func (s *Server) holdsNothing(ctx context.Context, tenant string, g *Grant) (bool, error) {
	if g.Tenant != tenant {
		return false, nil
	}
	for _, sk := range []string{skToken(g.ID), skRenewal(g.ID), skBrowserState(g.ID), skBrowserSignIn(g.ID)} {
		_, ok, err := s.KV.Get(ctx, tenantPK(tenant), sk)
		if err != nil || ok {
			return false, err
		}
	}
	if rows, err := s.KV.Query(ctx, tenantPK(tenant), "CHECKOUT#"+g.ID+"#"); err != nil || len(rows) > 0 {
		return false, err
	}
	if g.Connector != nil && g.Connector.Account != "" && (g.Connector.Revocation == nil || !g.Connector.Revocation.Disabled) {
		return false, nil
	}
	return true, nil
}
