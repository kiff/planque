package planque

import (
	"net/http"
	"strings"

	"github.com/kiff/planque/deploy/guides"
)

// guideSlug is the app a catalog entry connects, as its guide is named: a
// provider's app by its toolkit (connector:composio:gmail → gmail), a
// reviewed entry by its name without its version or kind (qonto-key@2 →
// qonto, stripe → stripe).
func guideSlug(entry string) string {
	if _, toolkit, ok := connectorEntry(entry); ok {
		return strings.ToLower(toolkit)
	}
	name, _, _ := strings.Cut(entry, "@")
	return strings.TrimSuffix(strings.ToLower(name), "-key")
}

// catalogGuide answers an app's adapted guide, for an agent of the
// account, before any connection: GET /v1/accounts/{tenant}/catalog/{entry}/guide.
// Coach reads it before it offers or connects the app. It names no owner's
// data: it's the same for every account.
func (s *Server) catalogGuide(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	entry := r.PathValue("entry")
	slug := guideSlug(entry)
	text, checked, ok := guides.Adapted(slug)
	if !ok {
		writeErr(w, http.StatusNotFound, "no guide for this app yet: read its catalog for its tools, what they return and the access each needs")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"entry": entry, "app": slug, "checked": checked, "guide": text})
}
