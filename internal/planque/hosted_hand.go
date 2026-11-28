package planque

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/store"
)

// "Your turn" in the hosted browser: the hand-over the owner's Chrome has
// (chrome_hand.go), through the hosted browser's live view, mid-run, on the
// phone. The taf's browse_hand opens the page in the grant's browser and
// keeps that session open, held as the grant's sign-in view; the phone asks
// for the view with a signed request and the owner does the step in it.
// Done seals what the site kept as the grant's sign-in, over the version
// the session opened with, and the taf reads "done" with browse_hand_result;
// the app's Not now ends the session and keeps nothing. The view stays
// one-time, short-lived, unrecorded and on the grant's sites only, and a
// page that pays is never handed over. While it's open, the taf's other
// calls on the grant wait. Offered while HostedHand is on.

func pkHostedHand(grant, hand string) string { return "HOSTEDHAND#" + grant + "#" + hand }

// hostedHand hands the person the step a in the hosted browser: the words
// checked, the page opened on the grant's sites, refused if it pays, the
// ref (if any) one of its fields or controls; then the session is kept open
// for the phone's view, and the taf gets the hand-over's id.
func (s *Server) hostedHand(ctx context.Context, g *Grant, p permit.Permit, a browserArgs, sites []string) (string, any) {
	if !s.HostedHand {
		return OpNotSent, browserRefusal(CodeNotGranted, "the hand-over isn't offered in the hosted browser", false)
	}
	if why := handWordsRefusal(a.Words); why != "" {
		return OpNotSent, browserRefusal(CodeBadArguments, why, false)
	}
	if len(a.Steps) > 0 || a.Hand != "" {
		return OpNotSent, browserRefusal(CodeBadArguments, browseHand+" takes url, words, and read and ref when it names an element: no steps", false)
	}
	resumer, ok := s.Browser.(browser.Resumer)
	if s.Browser == nil || !ok {
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "this installation's browser can't hold a hand-over", false)
	}
	state, from, err := s.openBrowserState(ctx, g)
	if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "the browser's sign-ins couldn't be opened", false)
	}
	// The hand-over holds the account's browser while it waits for the
	// person, as the sign-in view does.
	lease, left, err := s.takeBrowserLease(ctx, g.Tenant, s.now().Add(signInViewTTL))
	if errors.Is(err, errBrowserBusy) {
		return OpNotSent, busyRefusal(left)
	} else if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read this account's browser sessions", false)
	}
	sess, err := resumer.(browser.Worker).Start(ctx, browser.Start{Grant: g.ID, State: state, Sites: sites})
	if err != nil {
		s.dropBrowserLease(context.WithoutCancel(ctx), lease)
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "the browser couldn't be started: "+err.Error(), false)
	}
	kept := false
	defer func() {
		if !kept {
			_, _ = sess.Close(context.WithoutCancel(ctx))
			s.dropBrowserLease(context.WithoutCancel(ctx), lease)
		}
	}()
	if !p.Live(s.now()) {
		return OpNotSent, toolResult("Not sent: the permit expired before the call could be sent.", true)
	}
	page, err := sess.Open(ctx, a.URL)
	if err != nil {
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "the page couldn't be opened: "+err.Error(), false)
	}
	if code, why := browser.Allowed(page.URL, sites); code != "" {
		return OpNotSent, browserRefusal(code, "the page redirected: "+why, false)
	}
	// Never on a page that pays: a card is typed only in a checkout's own
	// fill, never in a view (the threat model's T3.4).
	if browser.Payment(page) {
		return OpNotSent, browserRefusal(browser.CodePaymentStep, "this page takes a payment or places an order, and such a page is never handed over: paying is browse_pay's, held for the owner", false)
	}
	if a.Read != "" && browser.Read(page) != a.Read {
		return OpNotSent, browserRefusal(browser.CodePageChanged, "the page isn't the read named ("+a.Read+"); it reads as "+browser.Read(page)+" now: hand it over again naming that read", false)
	}
	if a.Ref != "" && (a.Read == "" || !pageHasRef(page, a.Ref)) {
		return OpNotSent, browserRefusal(CodeBadArguments, "ref "+a.Ref+" isn't a field or control of that read (a ref goes with the read it's from)", false)
	}
	id := "hh_" + rand.Text()
	now := s.now()
	if err := s.keepHandAt(ctx, pkHostedHand(g.ID, id), chromeHand{Tenant: g.Tenant, Grant: g.ID, Words: a.Words, Ref: a.Ref, Opened: now}); err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not keep the hand-over", false)
	}
	rec := signInRecord{Session: sess.ID(), Site: browser.Site(page.URL), Expires: now.Add(signInViewTTL), OpenedAt: now.UTC(), Lease: lease,
		Hand: id, Words: a.Words, From: from}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skBrowserSignIn(g.ID), Ver: 1, TTL: rec.Expires.Add(time.Minute).Unix()}, rec, store.Cond{}); err != nil {
		_, _, _ = s.KV.Delete(ctx, pkHostedHand(g.ID, id), "-")
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not keep the hand-over", false)
	}
	kept = true
	g.event(now, "the taf handed the person a step at "+rec.Site+": "+a.Words)
	_ = s.saveGrant(ctx, g)
	return OpDone, browserAnswer(map[string]any{"state": "handed", "hand": id, "site": rec.Site,
		"next": "The person does it on their phone, in this connection's browser, within " + signInViewTTL.String() + ". End your turn now and tell them it's their turn; once they write that they're done, read their answer with " + grantToolName(g, browseHandResult) + " and this hand."})
}

// readHostedHand is what the taf reads of a hosted hand-over: the person's
// answer (read once), waiting while its view is still open, else expired.
func (s *Server) readHostedHand(ctx context.Context, g *Grant, hand string) (chromeHand, error) {
	h, err := s.readHandAt(ctx, pkHostedHand(g.ID, hand), g.ID, handOpenFor)
	if err != nil || h.Outcome != handWaiting {
		return h, err
	}
	// Waiting only while its view is open: a view that lapsed, or that the
	// owner's own sign-in replaced, ended the hand-over with nothing kept.
	if open, err := s.openSignIn(ctx, g); err == nil && (open == nil || open.Hand != hand) {
		_, _, _ = s.KV.Delete(ctx, pkHostedHand(g.ID, hand), "-")
		h.Outcome = handExpired
	}
	return h, nil
}

// openHand is g's open hosted hand-over hand, or why not, as a status.
func (s *Server) openHand(r *http.Request, sr signedRequest) (*Grant, *signInRecord, int, string) {
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		return nil, nil, http.StatusNotFound, "no such grant"
	} else if err != nil {
		return nil, nil, http.StatusInternalServerError, "the grant couldn't be read"
	}
	open, err := s.openSignIn(r.Context(), g)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, "the hand-over couldn't be read"
	}
	if open == nil || open.Hand == "" || open.Hand != r.PathValue("hand") {
		return nil, nil, http.StatusGone, "this step is no longer waiting for you: it was done, put off, or it lapsed"
	}
	return g, open, 0, ""
}

// startHandView is the phone's signed ask for a hosted hand-over's view:
// POST …/grants/{id}/hands/{hand}/view. It answers a one-time URL for the
// view of the session the taf handed over, good until the hand-over lapses;
// a second ask makes the first URL useless.
func (s *Server) startHandView(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	g, open, status, why := s.openHand(r, sr)
	if why != "" {
		writeErr(w, status, why)
		return
	}
	token := rand.Text()
	var cur signInRecord
	it, err := getDoc(r.Context(), s.KV, tenantPK(g.Tenant), skBrowserSignIn(g.ID), &cur)
	if err != nil || cur.Hand != open.Hand {
		writeErr(w, http.StatusGone, "this step is no longer waiting for you")
		return
	}
	old := cur.TokenSHA256
	cur.TokenSHA256 = tokenHash(token)
	ttl := cur.Expires.Add(time.Minute).Unix()
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skBrowserSignIn(g.ID), Ver: it.Ver + 1, TTL: ttl}, cur, store.Cond{Ver: &it.Ver}); err != nil {
		s.fail(w, r, err)
		return
	}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: pkSignInToken(cur.TokenSHA256), SK: "-", Ver: 1, TTL: ttl}, signInToken{Tenant: g.Tenant, Grant: g.ID, Expires: cur.Expires}, store.Cond{NotExists: true}); err != nil {
		s.fail(w, r, err)
		return
	}
	if old != "" {
		_, _, _ = s.KV.Delete(r.Context(), pkSignInToken(old), "-")
	}
	g.event(s.now(), "the view of the step handed over was opened "+sr.approval.Who())
	keepLang(r, g)
	_ = s.saveGrant(r.Context(), g)
	writeJSON(w, http.StatusOK, map[string]any{"view": s.BaseURL + "/browser/view/" + token, "expires": cur.Expires.UTC().Format(time.RFC3339),
		"site": cur.Site, "words": cur.Words})
}

// declineHand is the phone's signed Not now on a hosted hand-over: POST
// …/grants/{id}/hands/{hand}/decline. The session ends, nothing is kept,
// and the taf reads "declined".
func (s *Server) declineHand(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	g, open, status, why := s.openHand(r, sr)
	if why != "" {
		writeErr(w, status, why)
		return
	}
	if err := s.answerHandAt(r.Context(), pkHostedHand(g.ID, open.Hand), handDeclined, nil); err != nil && !errors.Is(err, errHandAnswered) {
		s.log().WarnContext(r.Context(), "a hand-over's Not now wasn't kept", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	s.dropSignIn(r.Context(), g)
	g.event(s.now(), "the person put off the step handed to them at "+open.Site+"; nothing was kept")
	_ = s.saveGrant(r.Context(), g)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// handOnPaymentPage is why a hand-over's view ends now, or "": the page its
// session shows pays (or left the grant's sites), and a card is never typed
// in a view, however the person got there (the threat model's T3.4). It's
// checked before each of the person's inputs, on each frame and at Done; a
// page that can't be read stops the input without ending the view. ended
// says the hand-over was ended: put off, nothing kept, the taf told.
func (s *Server) handOnPaymentPage(ctx context.Context, g *Grant, sess browser.Session, rec *signInRecord) (why string, ended bool) {
	if rec.Hand == "" {
		return "", false
	}
	l, ok := sess.(browser.Looker)
	if !ok {
		return "this browser can't check the page, so nothing is sent to it", false
	}
	page, err := l.Look(ctx)
	if err != nil {
		return "the page couldn't be checked, so nothing was sent to it; try again in a moment", false
	}
	sites := []string(nil)
	if g.Request.Browser != nil {
		sites = g.Request.Browser.Sites
	}
	code, _ := browser.Allowed(page.URL, sites)
	if !browser.Payment(page) && code == "" {
		return "", false
	}
	_ = s.answerHandAt(ctx, pkHostedHand(g.ID, rec.Hand), handDeclined, nil)
	s.dropSignIn(ctx, g)
	g.event(s.now(), "a step handed to the person ended at a page that takes a payment (or off the connection's sites); nothing was typed there and nothing was kept")
	_ = s.saveGrant(ctx, g)
	return "This page takes a payment, and a step handed to you can't be done on it. The view has ended and nothing was kept; the taf is told you put it off.", true
}

// waitingHand is one hosted hand-over still waiting for the person, as
// tafs cloud lists it in Needs you: no URL and no token, since the view is
// the phone's signed ask, made fresh on each open.
type waitingHand struct {
	Grant   string `json:"grant"`
	Agent   string `json:"agent"`
	Hand    string `json:"hand"`
	Site    string `json:"site"`
	Words   string `json:"words"`
	Since   string `json:"since"`
	Expires string `json:"expires"`
}

// accountHands answers GET /v1/accounts/{tenant}/hands, with the account's
// key as /connections: the hosted hand-overs still waiting for the person,
// none done, put off or lapsed.
func (s *Server) accountHands(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	items, err := s.KV.Query(r.Context(), tenantPK(tenant), "BROWSER_SIGNIN#")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	out := []waitingHand{}
	now := s.now()
	for _, it := range items {
		var rec signInRecord
		if json.Unmarshal(it.Doc, &rec) != nil || rec.Hand == "" || !now.Before(rec.Expires) {
			continue
		}
		g, err := s.loadGrant(r.Context(), tenant, strings.TrimPrefix(it.SK, "BROWSER_SIGNIN#"))
		if err != nil || g.Status != StatusConnected {
			continue
		}
		out = append(out, waitingHand{Grant: g.ID, Agent: g.Request.Agent, Hand: rec.Hand, Site: rec.Site, Words: rec.Words,
			Since: rec.OpenedAt.UTC().Format(time.RFC3339), Expires: rec.Expires.UTC().Format(time.RFC3339)})
	}
	writeJSON(w, http.StatusOK, map[string]any{"hands": out})
}
