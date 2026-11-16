package planque

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/store"
)

// keepUsedLink marks state as used once its sign-in came back, until the
// link's own deadline. A browser that lands on the link again (a "click
// here" after the redirect, a reopened tab) is then told what happened.
// Only the grant's place is kept: no code, no token. Not keeping it is
// said in the log and changes nothing else.
func (s *Server) keepUsedLink(ctx context.Context, state string, p pendingSignIn) {
	until := p.Expires
	if min := s.now().Add(time.Minute); until.Before(min) {
		until = min
	}
	if err := putDoc(ctx, s.KV, store.Item{PK: pkSignInUsed(state), SK: "-", TTL: until.Unix()}, usedSignIn{Tenant: p.Tenant, Grant: p.Grant}, store.Cond{}); err != nil {
		s.log().WarnContext(ctx, "used sign-in link not kept", slog.String("grant", p.Grant), slog.String("err", err.Error()))
	}
}

// spentLink answers a sign-in link planque no longer waits on. One that
// came back before says so: the sign-in already went through when its
// grant is signed in or connected, so the owner goes back to tafs, where
// it finishes; else the link was used, and the sign-in starts again from
// tafs. Any other link is unknown or expired, as before. The page names
// the app, never a grant, an account or a code.
func (s *Server) spentLink(w http.ResponseWriter, r *http.Request, state string) {
	var u usedSignIn
	if _, err := getDoc(r.Context(), s.KV, pkSignInUsed(state), "-", &u); err == nil {
		if g, err := s.loadGrant(r.Context(), u.Tenant, u.Grant); err == nil {
			lang, l := speech(r, g)
			if g.SignedIn || g.Status == StatusSignedIn || g.Status == StatusConnected {
				s.page(w, lang, http.StatusOK, l.SignedIn, l.SignInWentThrough(s.appName(g)))
				return
			}
			s.page(w, lang, http.StatusConflict, l.NotConnected, l.LinkUsedStartAgain)
			return
		}
	}
	lang, l := speech(r, nil)
	s.page(w, lang, http.StatusBadRequest, l.NotConnected, l.LinkUnknown)
}
