package planque

import (
	"errors"
	"net/http"
)

// switchChrome moves one browser connection to another of the account's
// paired Chromes (kiff/tafs docs/chrome.md, "Your Chrome" 2):
// PUT /v1/accounts/{t}/grants/{id}/chrome {chrome_id, chrome_sites_sig}. Its
// sites and Card stay; the phone signs its sites for the new Chrome
// (tafs-chrome/sites/v1 over {chrome_id, grant, sites}), since the
// extension runs nothing whose sites aren't signed for it. Routine, as every
// signed planque route that gives: a tap inside an open approval window, else
// the fingerprint.
func (s *Server) switchChrome(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var req struct {
		ChromeID       string `json:"chrome_id"`
		ChromeSitesSig string `json:"chrome_sites_sig"`
	}
	if err := decodeStrict(sr.body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	if !chromeIDRe.MatchString(req.ChromeID) || !chromeSigRe.MatchString(req.ChromeSitesSig) {
		writeErr(w, http.StatusBadRequest, "send {chrome_id, chrome_sites_sig}: the Chrome to use, and the phone's signature over the connection's sites for it")
		return
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	b := g.Request.Browser
	if b == nil || b.Driver != DriverExtension {
		writeErr(w, http.StatusConflict, "this connection doesn't use the owner's own Chrome")
		return
	}
	// Only a connection the owner approved moves (qa-kiff): never one only
	// asked for, failed or ended.
	switch g.Status {
	case StatusApproved, StatusSignedIn, StatusConnected, StatusListed, StatusNeedsReconnect:
	default:
		writeErr(w, http.StatusConflict, "only a connection the owner approved, and that hasn't ended, can use another Chrome")
		return
	}
	if b.Chrome == req.ChromeID {
		writeErr(w, http.StatusConflict, "this connection already uses that Chrome")
		return
	}
	if !s.relayOn() || !s.chromeOnFor(g.Tenant) {
		writeErr(w, http.StatusConflict, "the owner's own Chrome isn't offered to this account now")
		return
	}
	var c Chrome
	if _, err := getDoc(r.Context(), s.KV, tenantPK(sr.tenant), skChrome(req.ChromeID), &c); errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no paired Chrome of this account has that id")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	from := b.ChromeName
	b.Chrome, b.ChromeName = c.ID, label(c.Name)
	g.ChromeSitesSig = req.ChromeSitesSig
	g.event(s.now(), "switched to the Chrome \""+c.Name+"\" from \""+from+"\" "+sr.approval.Who()+", with the same sites")
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, s.view(r, g))
}
