package planque

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/kiff/planque/internal/connector"
)

// appSlug is an app's slug at the provider.
var appSlug = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_-]{0,127}$`)

// signInState answers whether the managed sign-in can sign people in to an
// app, GET /v1/catalog/{entry}/sign-in, for the nightly connect matrix:
// {entry, state}, state one of ready, on_the_owners_pick, needs_setup and
// not_at_provider. It reads the provider's app metadata only: no sign-in is
// made, no session, no allowance taken, and nothing of the provider's own
// (its config ids) is in the answer. The same for every account, so no
// tenant. Only the sign-in reader key reads it; without one, it isn't there.
func (s *Server) signInState(w http.ResponseWriter, r *http.Request) {
	if len(s.SignInReader) != sha256.Size {
		http.NotFound(w, r)
		return
	}
	key, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	sum := sha256.Sum256([]byte(key))
	if !ok || key == "" || subtle.ConstantTimeCompare(sum[:], s.SignInReader) != 1 {
		writeErr(w, http.StatusUnauthorized, "the sign-in reader key is required")
		return
	}
	entry := r.PathValue("entry")
	providerID, toolkit, ok := connectorEntry(entry)
	p := s.Connectors[providerID]
	if !ok || p == nil || !appSlug.MatchString(toolkit) {
		writeErr(w, http.StatusNotFound, "that isn't an app of the managed sign-in's catalog")
		return
	}
	state := ""
	app, err := p.App(r.Context(), toolkit)
	switch {
	case errors.Is(err, connector.ErrNotSetUp):
		state = "not_at_provider"
	case err != nil:
		s.log().WarnContext(r.Context(), "sign-in state unread", "entry", entry, "err", err.Error())
		writeErr(w, http.StatusBadGateway, "the provider couldn't be read")
		return
	case app.Setup == connector.SetupReady:
		state = "ready"
	case app.Setup == connector.SetupNeedsSetup:
		state = "needs_setup"
	case app.Setup == connector.SetupOnDemand:
		state = "on_the_owners_pick"
		// One already set up on someone's pick is ready for the next.
		if held, kept, err := s.setUpOf(r.Context(), providerID, toolkit); err == nil && kept && held.AuthConfig != "" {
			state = "ready"
		}
	default:
		state = "needs_setup"
	}
	writeJSON(w, http.StatusOK, map[string]string{"entry": entry, "state": state})
}
