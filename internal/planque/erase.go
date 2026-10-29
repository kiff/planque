package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
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
	if g.Status != StatusDisconnected && g.Status != StatusFailed && len(g.Tools) > 0 {
		if keyErr != nil {
			problems = append(problems, "the KIFF connector key could not be opened")
		} else {
			for _, t := range g.Tools {
				if err := s.KIFF.Disconnect(ctx, key, t.Name); err != nil {
					problems = append(problems, "KIFF kept "+t.Name+": "+err.Error())
				}
			}
		}
	}
	return append(problems, s.revoke(ctx, g)...)
}
