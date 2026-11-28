// Package planque is the service: an owner's phone enrols, an agent asks
// for a vendor account, the owner approves on the phone, signs in at the
// vendor, and planque relays each call KIFF permits to the vendor with a
// token only planque holds.
package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/issuing"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/spt"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
	"github.com/kiff/planque/internal/words"
	"github.com/kiff/planque/internal/wsapi"
)

// IdempotencyArg is the argument KIFF's gateway fills with the call's
// operation id; the permit binds it, and planque removes it before the
// vendor sees the call.
const IdempotencyArg = "planque_op"

// Server is planque.
type Server struct {
	// LiveIndex: the reads that want an account's live grants only (the
	// catalog, a held sign-in) read the table's live index, which DynamoDB
	// keeps with each grant record in the same write (store.Item.Live),
	// instead of the LIVE# rows. Off until the index exists
	// (PLANQUE_LIVE_GSI=on).
	LiveIndex bool
	// LiveListTest passes a vendor's whole live tool list through for the
	// test accounts' grants (PLANQUE_LIVE_LIST=test), for KIFF's Card to
	// decide each call's class: servesAll. Off, every account keeps the
	// tools planque serves by name.
	LiveListTest bool
	// BaseURL is where planque is reached publicly, without a trailing
	// slash: the vendor redirects the owner back to BaseURL/oauth/callback.
	BaseURL string
	// DigestV2: a grant's view carries digest_v 2, a digest over its words
	// in the app's language with the request (digestV2), once both apps
	// check it. Approval takes either digest all the same.
	DigestV2 bool
	KV       store.KV
	Sealer   seal.Sealer
	KIFF     *kiff.Client
	Permits  *permit.Verifier
	OAuth    *oauth.Client
	MCP      *mcp.Client
	// REST reaches an API-key grant's vendor (kiff/planque#64).
	REST *rest.Client
	// Issuing makes a purchase's single-use card through Stripe Issuing, in
	// test mode only (thune_cards.go); nil, no card is made.
	Issuing *issuing.Client
	// Tokens issues a purchase's shared payment token (path a, B1) from its
	// Issuing card, in test mode; nil, no token is made.
	Tokens *spt.Client
	// TelegramBase is the Telegram Bot API; empty, Telegram's own.
	TelegramBase string
	// Browser drives a browser grant's browser (kiff/planque#86); nil, no
	// browser call runs. Evidence keeps its screenshots; nil, none are.
	Browser browser.Worker
	// BrowserOffer is which accounts Browser is offered to
	// (PLANQUE_BROWSER_OFFER, BrowserOfferFrom): every account (All), or only
	// Tenants. Neither, it's offered to none, even with a worker configured.
	BrowserOffer BrowserOffer
	// ChromeBuilds are the extension's released builds, by the SHA-256 of
	// their packed manifest (deploy/chrome-builds.txt): only these pair and
	// connect. Empty, the owner's own Chrome isn't offered.
	ChromeBuilds map[string]ChromeBuild
	// ChromeOffer is which accounts the owner's own Chrome is offered to
	// once a build is listed (PLANQUE_CHROME_OFFER, ChromeOfferFrom): every
	// account (All), or only Tenants. Neither, it's offered to none.
	ChromeOffer ChromeOffer
	// ChromeScreenshots keeps an image of each acting step in the owner's
	// own Chrome, for the taf's Activity (PLANQUE_CHROME_SCREENSHOTS=on,
	// chrome_screenshot.go). Off, the routes answer 404 and nothing is kept.
	ChromeScreenshots bool
	// ChromeHand offers the hand-over's two tools (browse_hand,
	// browse_hand_result) on the owner's Chrome's new connections only: a
	// connection approved before keeps its tools. On in production unless
	// PLANQUE_CHROME_HAND=off (cmd/planque); false in a bare Server.
	ChromeHand bool
	// HostedHand offers the hand-over's two tools on the hosted browser's
	// new connections ("Your turn" through the live view, hosted_hand.go).
	// Off unless PLANQUE_HOSTED_HAND=on (cmd/planque).
	HostedHand bool
	// SealedFill takes a sealed fill in the owner's Chrome: the gateway's
	// describe of it, and its act with the phone's release (chrome_sealed.go).
	// Off unless PLANQUE_SEALED_FILL=on (cmd/planque); off, every sealed step
	// is refused before anything reaches the Chrome.
	SealedFill bool
	// ChromeSigner signs the relay's welcome and steps with planque's step
	// key (KMS in production), and ChromePoster reaches the extension's
	// WebSocket connections. Without both, the relay refuses every connect.
	ChromeSigner stepkey.Signer
	ChromePoster wsapi.Poster
	// ChromeWaker wakes a Chrome that isn't connected for a step, by Web
	// Push, when its extension gave planque where (chrome_wake.go). Nil: a
	// Chrome that isn't connected is offline, as before.
	ChromeWaker *ChromeWaker
	// VAPIDSubject is the contact in planque's VAPID tokens (an https URL or
	// a mailto). Set, with the relay on, planque wakes Chromes with its own
	// VAPID key (chrome_vapid.go); empty, it doesn't.
	VAPIDSubject string
	vapidMu      sync.Mutex
	vapidWaker   *ChromeWaker
	Evidence     browser.Evidence
	// Application definitions are reviewed public protocol facts; their
	// client credentials are resolved separately from an operator store.
	CatalogEntries map[string]CatalogEntry
	// CatalogOperators may draft, review, publish and retire the catalog
	// kept as data (catalog_store.go), by their credential and device.
	CatalogOperators map[string]CatalogOperator
	// markCache is the managed marks as this instance last read them.
	markCache    markCache
	catalogCache catalogCache
	// described keeps each pinned description's parse (describeEntry), so a
	// read of a revision never downloads and parses it twice.
	described describedCache
	// ServerApps is the reviewed list of which app a vendor's own server
	// is: its server, and the catalog entry it belongs to (entry.go).
	ServerApps map[string]string
	// Proved is which apps connect today (proved.go); nil is the reviewed
	// proved.json shipped in the binary.
	Proved map[string]Proof
	// ProofState, when set, makes KIFF's proof monitor the only source of
	// proved marks (proof_state.go); nil keeps Proved.
	ProofState *ProofState
	// DescriberServers are the servers, besides Stripe's own, whose first
	// request gets these describing tools as setup reads with its sign-in
	// (PLANQUE_DESCRIBER_SERVERS): the Eval's vendor, on a host we run.
	DescriberServers map[string][]string
	Connectors       map[string]connector.Provider
	DoorSink         door.Sink
	// EventSink takes app events for tafs (planque.event.v1); nil, no app
	// event can be approved.
	EventSink door.EventSink
	// KIFFEventSink takes the refund events a profile's case block names,
	// for KIFF (refund_events.go); nil hands them to no one but the taf.
	KIFFEventSink door.EventSink
	// StripeRefunds offers planque's reviewed Stripe refunds profile
	// (stripe_refunds.go) on a managed Stripe connection that asks for the
	// refund tool, to every account. Off, the refund tool takes a drafted
	// profile as any other tool does, with no case block.
	StripeRefunds bool
	// StripeRefundsTenants offers it to these accounts only, the others as
	// if it were off: how an installation tries it on one account.
	StripeRefundsTenants map[string]bool
	// ManagedApps replaces managed_apps.json when set (tests) with the list
	// mode: the apps a managed connection may be made to, by folded name.
	// ManagedExcept, when ManagedApps isn't set, replaces it with the all
	// mode: every app but these. ManagedTestTenants are the accounts that may
	// make one to any app (PLANQUE_MANAGED_TEST_TENANTS, ManagedTestTenantsFrom).
	ManagedApps        map[string]bool
	ManagedExcept      map[string]bool
	ManagedTestTenants map[string]bool
	// BrowserPay turns paying on in the hosted browser (PLANQUE_BROWSER_PAY
	// on): browse_pay is offered and runs, held for the owner each time. Off
	// (the default), the hosted browser reads, fills and clicks, but pays
	// nothing: browse_pay is neither offered nor run.
	BrowserPay bool
	// ProofTestTenants are the accounts a robot test user runs
	// (PLANQUE_PROOF_TEST_TENANTS, test_tenants.go): they take no live
	// Stripe key, and none of their tafs pays (browse_pay).
	ProofTestTenants map[string]bool
	// NeverUse are grants whose sign-ins are never reused or called, by any
	// grant on the same account at the app (PLANQUE_NEVER_USE, NeverUseFrom):
	// an account no taf is meant to reach.
	NeverUse map[string]bool
	// OwnClientTenants are the accounts shown an app's own entry instead of
	// the managed connector's toolkit (PLANQUE_OWN_CLIENT_TENANTS, ownClientFor).
	OwnClientTenants       map[string]bool
	Applications           map[string]appconfig.Application
	ApplicationCredentials appcredentials.Reader
	Logger                 *slog.Logger
	Now                    func() time.Time
	// WindowRequired makes a routine signed call (approving a sign-in or a
	// key, confirming, a trigger, enrolling) need either the phone's
	// approval key or its device key inside an open approval window
	// (kiff/tafs#494). Off, a device-key signature with no window is taken
	// as before, for apps from before the window; a signature that names a
	// window must name an open one either way.
	WindowRequired bool
	// ComputerWindowRequired is WindowRequired for computers alone: a routine
	// signed call by the device key of a device KIFF doesn't name as a phone's
	// or a tablet's (its platform ios, ipados or android) must name an open
	// approval window. A computer's device key (a TPM's, a Mac's, a software
	// key on disk) signs with nobody asked, so with no window it would approve
	// a connection nobody was present for. Phones are left to WindowRequired,
	// until their apps sign every routine call inside a window.
	ComputerWindowRequired bool
	// AppReturn is where the browser that completed a vendor sign-in is
	// sent, with the code that confirms it: the owner's app (tafs), for
	// example tafs://planque/signed-in. The operator sets it; no request
	// can, so a sign-in link someone else opens never hands its code back
	// to whoever approved the grant.
	AppReturn string
	// SignInReader is the SHA-256 of the one key that reads whether the
	// managed sign-in can sign people in to an app (GET /v1/catalog/{entry}/
	// sign-in): for the nightly connect matrix, and nothing else. Unset,
	// the route isn't there.
	SignInReader []byte
	// RefreshWait bounds how long a call waits for another one's refresh.
	RefreshWait time.Duration

	refreshing sync.Map // grant id → *sync.Mutex, one refresh per process
}

func (s *Server) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Server) log() *slog.Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

func (s *Server) redirect() string { return s.BaseURL + "/oauth/callback" }

// Handler serves planque's routes.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]bool{"ok": true}) })
	// Crawlers keep off planque.kiff.dev: its URLs are for programs, and a
	// crawler that follows one from a public page only gets errors.
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("User-agent: *\nDisallow: /\n"))
	})
	// The apps that connect today, the way a check ran them (proved.json), for
	// readers with no account: kiff.dev's apps list. Public, reviewed data only.
	mux.HandleFunc("GET /v1/proved", s.servePublicProofs)
	mux.HandleFunc("POST /v1/sweep", s.sweep)
	// The one-time clean-up of stored answers alone, for a schedule more
	// frequent than the sweep's while it runs (sweep.go).
	mux.HandleFunc("POST /v1/sweep/answers", s.sweepAnswersOnce)
	mux.HandleFunc("POST /v1/accounts/{tenant}/enrolment", s.enrol)
	// The owner's own Chrome (chrome.go): paired by the phone, the extension's
	// own routes unauthenticated, each step of pairing proving itself.
	mux.HandleFunc("POST /v1/accounts/{tenant}/chromes/pairings", s.signed(s.openChromePairing))
	mux.HandleFunc("GET /v1/accounts/{tenant}/chromes/pairings/{pairing}", s.signedSafe(s.readChromePairing))
	mux.HandleFunc("POST /v1/accounts/{tenant}/chromes/pairings/{pairing}/confirm", s.signed(s.confirmChromePairing))
	mux.HandleFunc("GET /v1/accounts/{tenant}/chromes", s.signedSafe(s.listChromes))
	mux.HandleFunc("GET /v1/accounts/{tenant}/chromes/self/runs", s.chromeOwnRuns)
	// An acting step's image in the owner's own Chrome: the extension sends
	// it, the owner reads it (chrome_screenshot.go; off unless switched on).
	mux.HandleFunc("POST /v1/extension/steps/{step}/screenshot", extensionCORS(s.postChromeShot))
	mux.HandleFunc("GET /v1/accounts/{tenant}/chromes/screenshots/{step}", s.signedSafe(s.readChromeShot))
	// One operation's images, listed for a taf's Activity row by its proposal_id.
	mux.HandleFunc("GET /v1/accounts/{tenant}/chromes/screenshots", s.signedSafe(s.listChromeShots))
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/chromes/{chrome}", s.signedSafe(s.unpairChrome))
	mux.HandleFunc("PUT /v1/accounts/{tenant}/grants/{id}/chrome", s.signed(s.switchChrome))
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/chrome/run", s.readChromeRun)
	mux.HandleFunc("POST /v1/extension/pairings/{lookup}", extensionCORS(s.postChromeKey))
	mux.HandleFunc("GET /v1/extension/pairings/{lookup}", extensionCORS(s.collectChromePairing))
	mux.HandleFunc("GET /v1/extension/push-key", extensionCORS(s.pushKey))
	mux.HandleFunc("POST /v1/extension/unpair", extensionCORS(s.extensionUnpair))
	// The extension has no host permission for planque, so Chrome asks first
	// (chrome_cors.go).
	mux.HandleFunc("OPTIONS /v1/extension/", extensionPreflight)
	mux.HandleFunc("POST /v1/accounts/{tenant}/halts/confirm", s.confirmHalt)
	mux.HandleFunc("GET /v1/catalog/{entry}/sign-in", s.signInState)
	// The catalog kept as data, for its operators (catalog_store.go).
	mux.HandleFunc("GET /v1/marks", s.listMarks)
	mux.HandleFunc("POST /v1/marks/managed", s.setManagedMark)
	mux.HandleFunc("POST /v1/marks/managed/{vendor}/withdraw", s.withdrawManagedMark)
	mux.HandleFunc("GET /v1/catalog-entries", s.listCatalogEntries)
	mux.HandleFunc("POST /v1/catalog-entries", s.draftCatalogEntry)
	mux.HandleFunc("GET /v1/catalog-entries/{id}", s.getCatalogEntry)
	mux.HandleFunc("PUT /v1/catalog-entries/{id}", s.changeCatalogEntry)
	mux.HandleFunc("POST /v1/catalog-entries/{id}/review", s.reviewCatalogEntry)
	mux.HandleFunc("POST /v1/catalog-entries/{id}/publish", s.signCatalogEntry("publish"))
	mux.HandleFunc("POST /v1/catalog-entries/{id}/retire", s.signCatalogEntry("retire"))
	mux.HandleFunc("GET /v1/accounts/{tenant}/catalog", s.vendorCatalog)
	mux.HandleFunc("GET /v1/accounts/{tenant}/catalog/tools", s.searchTools)
	mux.HandleFunc("GET /v1/accounts/{tenant}/catalog/{entry}/events", s.catalogEvents)
	mux.HandleFunc("GET /v1/accounts/{tenant}/catalog/{entry}/guide", s.catalogGuide)
	mux.HandleFunc("POST /v1/accounts/{tenant}/catalog/{entry}/grants", s.refusalsLogged(s.requestCatalogGrant))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants", s.refusalsLogged(s.requestGrant))
	mux.HandleFunc("POST /v1/accounts/{tenant}/api/operations", s.refusalsLogged(s.apiOperations))
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}", s.getGrant)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/catalog", s.catalog)
	// Reading what a sign-in covers, before its tools were approved, could
	// read business data no approval named (kiff/planque#23): it's gone.
	// Setup asks the owner which mode instead.
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/describe", func(w http.ResponseWriter, _ *http.Request) {
		writeErr(w, http.StatusGone, "reading what a sign-in covers was removed: nothing is read through a vendor's tools before the owner approves them. Ask the owner which account and mode it is")
	})
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/tools", s.setTools)
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/grants/{id}/tools", s.withdrawTools)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/read", s.setupRead)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/approve", s.signed(s.approve))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/key", s.signed(s.giveKey))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/connect-fresh", s.signed(s.connectFresh))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/confirm", s.signed(s.confirm))
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/grants/{id}", s.signedSafe(s.disconnect))
	// The live sign-in view of a browser grant (kiff/planque#86, step 4).
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/sign-in", s.signed(s.startBrowserSignIn))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/hands/{hand}/view", s.signed(s.startHandView))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/hands/{hand}/decline", s.signed(s.declineHand))
	mux.HandleFunc("GET /browser/view/{token}", s.browserView)
	mux.HandleFunc("POST /browser/view/{token}/frame", s.browserViewFrame)
	mux.HandleFunc("POST /browser/view/{token}/input", s.browserViewInput)
	mux.HandleFunc("POST /browser/view/{token}/done", s.browserViewDone)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants", s.signedSafe(s.listGrants))
	mux.HandleFunc("GET /v1/accounts/{tenant}/agents/{agent}/grants", s.agentGrants)
	mux.HandleFunc("POST /v1/accounts/{tenant}/agents/{agent}/end", s.endAgent)
	mux.HandleFunc("GET /v1/accounts/{tenant}/asked-grants", s.everyAskedGrant)
	mux.HandleFunc("GET /v1/accounts/{tenant}/connections", s.accountConnections)
	mux.HandleFunc("GET /v1/accounts/{tenant}/hands", s.accountHands)
	mux.HandleFunc("GET /v1/accounts/{tenant}/sessions/{session}/agents", s.sessionAgents)
	mux.HandleFunc("GET /v1/accounts/{tenant}/export", s.signedAs(kindSafe, capExport, s.export))
	// Deleting the account (kiff/tafs#185): everything planque keeps for it.
	mux.HandleFunc("DELETE /v1/accounts/{tenant}", s.eraseAccount)
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/finished-grants", s.signedFor(capDelete, s.eraseFinished))
	// Deleting it from one signed intent, which tafs cloud carries (deletions.go).
	mux.HandleFunc("POST /v1/account-deletions", s.acceptDeletion)
	mux.HandleFunc("POST /v1/account-deletions/finish", s.finishDeletion)
	mux.HandleFunc("POST /v1/account-deletions/status", s.deletionStatus)
	mux.HandleFunc("GET /oauth/callback", s.callback)
	mux.HandleFunc("GET /connectors/callback", s.connectorCallback)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/channels", s.signedSafe(s.listChannels))
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/triggers/types", s.triggerType)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/events", s.grantEvents)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/triggers", s.listTriggers)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/triggers", s.signed(s.createTrigger))
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/grants/{id}/triggers/{trigger}", s.signedSafe(s.deleteTrigger))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/triggers/{trigger}/confirm", s.confirmTrigger)
	mux.HandleFunc("POST /connectors/{provider}/webhook", s.connectorWebhook)
	mux.HandleFunc("POST /telegram/{hook}", s.telegramWebhook)
	mux.HandleFunc("GET /oauth/client-metadata.json", s.clientMetadata)
	mux.HandleFunc("POST /mcp/{tenant}/{id}", s.relay)
	return mux
}

// clientMetadata is the same document whose URL is sent as client_id when
// a vendor supports CIMD. No tenant or sign-in data belongs in this response.
func (s *Server) clientMetadata(w http.ResponseWriter, _ *http.Request) {
	doc, err := oauth.NewClientDocument(s.BaseURL+"/oauth/client-metadata.json", ClientName, s.redirect())
	if err != nil {
		writeErr(w, http.StatusServiceUnavailable, "this installation's public client metadata is not configured")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_ = json.NewEncoder(w).Encode(doc)
}

type apiError struct {
	status int
	msg    string
}

func (e apiError) Error() string { return e.msg }

// apiCoded is an apiError with a code the app acts on (kiff/tafs#494):
// critical_key_required, the approval key must sign; window_required, the
// device key signed in no open window, so the app opens one and signs again.
type apiCoded struct {
	apiError
	code string
	// words, when set, is the person's line, in their language: the apps
	// show it as sent, never a line of their own for the status.
	words string
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	var rf *toolRefusal
	if errors.As(err, &rf) {
		writeRefusal(w, rf)
		return
	}
	var ce apiCoded
	if errors.As(err, &ce) {
		body := map[string]string{"error": ce.msg, "code": ce.code}
		if ce.words != "" {
			body["words"] = ce.words
		}
		writeJSON(w, ce.status, body)
		return
	}
	var ae apiError
	if errors.As(err, &ae) {
		writeErr(w, ae.status, ae.msg)
		return
	}
	var cf connectFailed
	if errors.As(err, &cf) {
		s.log().WarnContext(r.Context(), "connect failed", slog.String("path", loggedPath(r.URL.Path)), slog.String("err", err.Error()))
		writeJSON(w, http.StatusConflict, map[string]string{"error": string(cf), "code": "connect_failed"})
		return
	}
	s.log().ErrorContext(r.Context(), "request failed", slog.String("path", loggedPath(r.URL.Path)), slog.String("err", err.Error()))
	writeErr(w, http.StatusInternalServerError, "planque could not complete this request")
}

// failPage is fail for a page a person's browser shows: the cause is
// logged as fail logs it, and the person reads that it couldn't be
// confirmed, in their language, never planque's JSON.
func (s *Server) failPage(w http.ResponseWriter, r *http.Request, g *Grant, err error) {
	s.log().ErrorContext(r.Context(), "request failed", slog.String("path", loggedPath(r.URL.Path)), slog.String("err", err.Error()))
	lang, l := speech(r, g)
	s.page(w, lang, http.StatusInternalServerError, l.NotConnected, l.CantConfirm)
}

// loggedPath is a request's path for the log, with a sign-in view's token
// (its only credential) left out.
func loggedPath(p string) string {
	if rest, ok := strings.CutPrefix(p, "/browser/view/"); ok {
		_, tail, _ := strings.Cut(rest, "/")
		if tail != "" {
			return "/browser/view/{token}/" + tail
		}
		return "/browser/view/{token}"
	}
	return p
}

func readBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r.Body, 64<<10+1))
	if err != nil || len(b) > 64<<10 {
		return nil, apiError{http.StatusBadRequest, "the body is too long or unreadable"}
	}
	return b, nil
}

func decodeStrict(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		// The field planque doesn't take is named, so a caller fixes it on
		// its first retry instead of guessing.
		if f, ok := strings.CutPrefix(err.Error(), "json: unknown field "); ok {
			return apiError{http.StatusBadRequest, "the body must be a JSON object with the documented fields: unknown field " + f}
		}
		return apiError{http.StatusBadRequest, "the body must be a JSON object with the documented fields"}
	}
	return nil
}

// signedRequest is a request the owner signed on their device.
type signedRequest struct {
	tenant    string
	body      []byte
	enrolment Enrolment
	approval  Approval // the signature, as evidence
}

// verifySigned checks Planque-Device-Signature over r, and that this exact
// signed request was never made before. The device is one KIFF knows (d=),
// verified by KIFF, asked with the account's connector key, which opens
// connector (kiff/planque#10). A signature without d=, from a phone that
// enrolled its own key before KIFF kept devices, is no longer taken
// (kiff/planque#12).
// What KIFF is asked a signer may do, as KIFF's own roles say: connect a
// vendor, or, only for those who may in KIFF too, export or delete the
// account (kiff/tafs#185, kiff/planque#46).
const (
	capConnect = "tools.connect"
	capExport  = "account.export"
	capDelete  = "account.delete"
)

// capWords says what a capability lets the signer do, for a refusal.
var capWords = map[string]string{capConnect: "connect tools", capExport: "export the account", capDelete: "delete the account"}

// The kinds of a signed route (kiff/tafs#494), decided here, never by the
// caller: safe takes authority away or reads, so the device key alone
// signs it; routine gives a little (a sign-in, a key, a confirm, a
// trigger), so inside an open approval window the device key signs it,
// else the approval key; critical widens what a taf may do, so only the
// approval key, behind the fingerprint, signs it. planque has no critical
// route of its own: the account's deletion is KIFF's to check.
const (
	kindSafe     = "safe"
	kindRoutine  = "routine"
	kindCritical = "critical"
)

func (s *Server) verifySigned(r *http.Request, tenant, capability string, body []byte, connector func() (string, error)) (Approval, error) {
	return s.verifySignedAs(r, kindRoutine, tenant, capability, body, connector)
}

// verifySignedAs is verifySigned for a route of this kind (kiff/tafs#494).
func (s *Server) verifySignedAs(r *http.Request, kind, tenant, capability string, body []byte, connector func() (string, error)) (Approval, error) {
	sg, err := device.ParseHeader(r.Header.Get(device.Header))
	if err != nil {
		return Approval{}, apiError{http.StatusUnauthorized, "this request needs the owner's signature from their device (" + device.Header + ")"}
	}
	msg := device.MessageWindow(s.BaseURL, tenant, r.Method, r.URL.RequestURI(), body, sg.At, sg.Window)
	// Signed together with another request (kiff.requests.v1): the phone
	// signed this request's message's hash beside the other's, once, with
	// its approval key. The use is this request's alone.
	var own, request string
	if sg.With != "" {
		request = string(msg)
		sum := sha256.Sum256(msg)
		own = hex.EncodeToString(sum[:])
		if own == sg.With {
			return Approval{}, apiError{http.StatusUnauthorized, "a request can't be signed together with itself"}
		}
		msg = device.RequestsMessage(tenant, sg.At, own, sg.With)
	}
	a := Approval{Message: string(msg), Request: request, Body: string(body), Signature: sg.Signature, SignedAt: sg.At, ApprovedAt: s.now().UTC()}
	// What the replay record names as the signer, with the message.
	if sg.Device == "" {
		return Approval{}, apiError{http.StatusUnauthorized, "this device isn't known to KIFF: sign in to KIFF on it again"}
	}
	var signer []byte
	{
		// planque keeps the window and the replay check: KIFF says only
		// that the device signed these bytes, not when.
		if err := device.CheckTime(sg.At, s.now()); err != nil {
			return Approval{}, apiError{http.StatusUnauthorized, err.Error()}
		}
		key, err := connector()
		if err != nil {
			return Approval{}, err
		}
		// KIFF is asked about the window only where it counts: a routine
		// route. A safe route takes the device key as it is, and a critical
		// one takes no window at all.
		window := ""
		if kind == kindRoutine {
			window = sg.Window
		}
		v, err := s.KIFF.VerifyDevice(r.Context(), key, sg.Device, msg, sg.Signature, capability, window)
		if err != nil {
			s.log().ErrorContext(r.Context(), "device check failed", slog.String("tenant", tenant), slog.String("err", err.Error()))
			return Approval{}, apiError{http.StatusBadGateway, "KIFF could not check the device that signed"}
		}
		if !v.Valid {
			// The phone reads only a 401 or 403; the reason is kept here, so a
			// refusal can be told apart (another account's device, a key KIFF
			// doesn't have) without reading KIFF's records. No key material.
			s.log().WarnContext(r.Context(), "device refused", slog.String("tenant", tenant), slog.String("device", sg.Device),
				slog.String("reason", v.Reason), slog.String("capability", capability))
			return Approval{}, deviceRefused(v.Reason, capability)
		}
		if sg.With != "" && v.Key != kiff.KeyCritical {
			s.log().WarnContext(r.Context(), "a signature over two requests not by the approval key", slog.String("tenant", tenant), slog.String("device", sg.Device), slog.String("key", v.Key))
			return Approval{}, apiError{http.StatusForbidden, "two requests signed together need the phone's approval key, behind its fingerprint"}
		}
		if err := keyFits(kind, v, sg.Window, s.WindowRequired || s.ComputerWindowRequired && !phonePlatforms[v.Platform]); err != nil {
			s.log().WarnContext(r.Context(), "device key short of the route", slog.String("tenant", tenant), slog.String("device", sg.Device),
				slog.String("kind", kind), slog.String("key", v.Key), slog.String("window", sg.Window))
			return Approval{}, err
		}
		a.DeviceID, a.Subject, a.DeviceName, a.Fingerprint, a.Key, a.Platform = v.DeviceID, v.Subject, strings.TrimSpace(v.Name), v.DeviceID, v.Key, v.Platform
		if v.Key == kiff.KeyDevice && window != "" {
			a.Window = window
		}
		signer = []byte("kiff-device\n" + v.DeviceID + "\n")
	}
	used := append(append([]byte{}, signer...), msg...)
	if own != "" {
		used = append(append(used, 0), own...)
	}
	sum := sha256.Sum256(used)
	ttl := s.now().Add(2 * device.MaxSkew).Unix()
	err = s.KV.Put(r.Context(), store.Item{PK: pkSigned(hex.EncodeToString(sum[:])), SK: "-", Doc: []byte("{}"), TTL: ttl},
		store.Cond{NotExists: true})
	if errors.Is(err, store.ErrCondition) {
		return Approval{}, apiCoded{apiError{http.StatusConflict, "this signed request was already made; sign a new one"}, "signature_replayed", ""}
	}
	if err != nil {
		return Approval{}, err
	}
	return a, nil
}

// phonePlatforms are the platforms KIFF names for a phone or a tablet, whose
// device key is the hardware key its app pairs (the Secure Enclave, Android's
// keystore). Any other, a computer's or one KIFF didn't name, is a computer's
// for ComputerWindowRequired.
var phonePlatforms = map[string]bool{"ios": true, "ipados": true, "android": true}

// keyFits says whether the key KIFF found, and the window the signature
// named, are enough for a route of this kind (kiff/tafs#494). The approval
// key is enough for any kind. A signature that names a window counts only
// while KIFF says that window is open. On a routine route, a device-key
// signature naming no window is enough only until windowRequired, for apps
// from before the window.
func keyFits(kind string, v kiff.Verified, window string, windowRequired bool) error {
	if v.Key == kiff.KeyCritical || kind == kindSafe {
		return nil
	}
	if kind == kindCritical {
		return apiCoded{apiError{http.StatusForbidden, "this widens what a taf may do: sign it with the phone's approval key, behind its fingerprint or face"}, "critical_key_required", ""}
	}
	if window != "" {
		if v.Window == nil || !v.Window.Open || v.Window.ID != window {
			return apiCoded{apiError{http.StatusForbidden, "the approval window this was signed in isn't open: open one with the phone's approval key, behind its fingerprint, and sign it again"}, "window_required", ""}
		}
		return nil
	}
	if windowRequired {
		return apiCoded{apiError{http.StatusForbidden, "this needs the phone's approval key, or its device key inside an open approval window: open one, behind the fingerprint, and sign it again"}, "window_required", ""}
	}
	return nil
}

// deviceRefused says why KIFF didn't vouch for a device's signature.
func deviceRefused(reason, capability string) apiError {
	switch reason {
	case "unknown_device":
		return apiError{http.StatusUnauthorized, "KIFF doesn't know the device that signed, or it was removed"}
	case "signature_invalid":
		return apiError{http.StatusUnauthorized, device.ErrSignature.Error()}
	case "not_a_member":
		return apiError{http.StatusForbidden, "the person whose device signed isn't a member of this account"}
	case "not_allowed":
		return apiError{http.StatusForbidden, "the person whose device signed can't " + capWords[capability] + " on this account"}
	}
	return apiError{http.StatusUnauthorized, "KIFF didn't vouch for the device that signed"}
}

// signed admits a request only with the owner's signature: from a device
// KIFF knows, of a person who may connect tools. The route is routine
// (kiff/tafs#494): the approval key, or the device key inside an open
// approval window.
func (s *Server) signed(next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return s.signedAs(kindRoutine, capConnect, next)
}

// signedSafe is signed for a route that takes authority away or reads:
// the device key alone, as today, window or not.
func (s *Server) signedSafe(next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return s.signedAs(kindSafe, capConnect, next)
}

// signedFor is signed, for a person KIFF says holds capability.
func (s *Server) signedFor(capability string, next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return s.signedAs(kindRoutine, capability, next)
}

// signedAs is signed for a route of this kind, for a person KIFF says holds
// capability.
func (s *Server) signedAs(kind, capability string, next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := r.PathValue("tenant")
		body, err := readBody(r)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		var e Enrolment
		if _, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &e); errors.Is(err, errNotFound) {
			writeErr(w, http.StatusForbidden, "this account isn't set up to connect tools yet")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		a, err := s.verifySignedAs(r, kind, tenant, capability, body, func() (string, error) { return s.openConnector(r.Context(), e) })
		if err != nil {
			s.fail(w, r, err)
			return
		}
		next(w, r, signedRequest{tenant: tenant, body: body, enrolment: e, approval: a})
	}
}

// enrol hands planque the account's KIFF connector key. A device KIFF
// knows (d=) signs it, and KIFF says whether its person may connect tools:
// any of the person's devices may (kiff/planque#10). The enrolled key is
// kept while KIFF still accepts it for the account: the tools it connected
// are that key's in KIFF, so a second device setting up changes nothing,
// and the answer says kept. A new key replaces it only once KIFF no longer
// takes it (revoked, or no longer the account's connector). Without d=, a phone enrols its own key and signs with
// it, as before: a different phone can't replace an enrolled one, and the
// enrolled phone may re-enrol to hand planque a new connector key.
func (s *Server) enrol(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		// DeviceKey is what a phone sent before KIFF kept devices. It is
		// still accepted, so apps that send it keep working, and ignored:
		// the device is KIFF's (d=), kiff/planque#12.
		DeviceKey    string `json:"device_key,omitempty"`
		ConnectorKey string `json:"connector_key"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	connector := strings.TrimSpace(req.ConnectorKey)
	var prev Enrolment
	prevItem, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &prev)
	// The write below is conditional on what was read here, so two
	// enrolments racing cannot both pass this check.
	cond := store.Cond{NotExists: true}
	if err == nil {
		cond = store.Cond{Ver: &prevItem.Ver}
	} else if !errors.Is(err, errNotFound) {
		s.fail(w, r, err)
		return
	}
	if connector == "" {
		// No key sent: an app asks first whether the account's enrolled key
		// still stands, and mints one only when it doesn't, so no key is made
		// that nothing uses (KIFF won't let a device revoke one).
		s.enrolWithoutKey(w, r, tenant, body, prev, prevItem.Ver > 0)
		return
	}
	// The connector key first: KIFF is asked with it whether the device signed.
	me, err := s.KIFF.Me(r.Context(), connector)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "KIFF did not accept the connector key")
		return
	}
	if p := me.ConnectorProblem(tenant); p != "" {
		writeErr(w, http.StatusBadRequest, p)
		return
	}
	a, err := s.verifySigned(r, tenant, capConnect, body, func() (string, error) { return connector, nil })
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if prevItem.Ver > 0 {
		kept, err := s.connectorStands(r.Context(), tenant, prev)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if kept {
			s.log().InfoContext(r.Context(), "enrolment kept", slog.String("tenant", tenant), slog.String("device", a.Fingerprint))
			writeJSON(w, http.StatusOK, map[string]any{"tenant": tenant, "device": a.Fingerprint, "kept": true})
			return
		}
	}
	sealed, err := s.Sealer.Seal(r.Context(), []byte(connector), map[string]string{"tenant": tenant, "kind": "connector"})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	e := Enrolment{Tenant: tenant, Connector: sealed, EnrolledAt: s.now().UTC()}
	err = putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(tenant), SK: skEnrolment, Ver: prevItem.Ver + 1}, e, cond)
	if errors.Is(err, store.ErrCondition) {
		s.fail(w, r, apiCoded{apiError{http.StatusConflict, "the enrolment changed while this one was checked; enrol again"}, "enrolment_changed", ""})
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(r.Context(), "device enrolled", slog.String("tenant", tenant), slog.String("device", a.Fingerprint))
	writeJSON(w, http.StatusCreated, map[string]string{"tenant": tenant, "device": a.Fingerprint})
}

// enrolWithoutKey answers an enrolment that sends no connector key. It
// never writes: with no enrolment, or one whose key KIFF no longer takes
// as the account's connector, it answers 409 needs_key, and the app mints
// a key and enrols again with it. Else the request is verified as the
// signed routes are, against the enrolled connector, and it answers kept.
func (s *Server) enrolWithoutKey(w http.ResponseWriter, r *http.Request, tenant string, body []byte, prev Enrolment, enrolled bool) {
	needsKey := apiCoded{apiError{http.StatusConflict, "no connector key is enrolled that KIFF still takes; mint one and enrol again with it"}, "needs_key", ""}
	if !enrolled {
		s.fail(w, r, needsKey)
		return
	}
	kept, err := s.connectorStands(r.Context(), tenant, prev)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !kept {
		s.fail(w, r, needsKey)
		return
	}
	a, err := s.verifySigned(r, tenant, capConnect, body, func() (string, error) { return s.openConnector(r.Context(), prev) })
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(r.Context(), "enrolment kept", slog.String("tenant", tenant), slog.String("device", a.Fingerprint))
	writeJSON(w, http.StatusOK, map[string]any{"tenant": tenant, "device": a.Fingerprint, "kept": true})
}

// agentKey checks the KIFF key an agent sends, for tenant.
func (s *Server) agentKey(r *http.Request, tenant string) (kiff.Me, error) {
	key, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || strings.TrimSpace(key) == "" {
		return kiff.Me{}, apiError{http.StatusUnauthorized, "a KIFF key is required"}
	}
	me, err := s.KIFF.Me(r.Context(), strings.TrimSpace(key))
	if errors.Is(err, kiff.ErrUnauthorized) {
		return kiff.Me{}, apiError{http.StatusUnauthorized, "KIFF did not accept the key"}
	}
	if err != nil {
		return kiff.Me{}, apiError{http.StatusBadGateway, "KIFF could not be asked who the key is"}
	}
	if me.TenantID != tenant {
		return kiff.Me{}, apiError{http.StatusForbidden, "the key is for another KIFF account"}
	}
	return me, nil
}

// speech is the language planque's pages speak to r, and its lines: the
// language g's app asked in, else the browser's, else English. g may be nil.
func speech(r *http.Request, g *Grant) (string, words.Lines) {
	stored := ""
	if g != nil {
		stored = g.Lang
	}
	lang := words.Pick(stored, r.Header.Get("Accept-Language"))
	return lang, words.For(lang)
}

// keepLang keeps on g the language the owner's app asks in, when it says.
func keepLang(r *http.Request, g *Grant) {
	if a := r.Header.Get("Accept-Language"); a != "" {
		g.Lang = words.Pick("", a)
	}
}

// page is a page in lang: a title and a line, both already in it.
func (s *Server) page(w http.ResponseWriter, lang string, status int, title, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.WriteHeader(status)
	_ = pageTmpl.Execute(w, map[string]string{"Lang": lang, "Title": title, "Message": msg})
}

// returnToApp hands a finished sign-in's return link u, which carries its
// one-use code, to the owner's app. A phone's or tablet's approval (a) gets a
// redirect, which its browser session (ASWebAuthenticationSession, Custom
// Tabs) catches. Any other device's, or one whose platform isn't known, gets
// a page with an Open tafs link: a computer's browser opens an app's link
// only on the person's click, never on a redirect with no click behind it
// (17 Nov, Windows: Chrome dropped it). Either way the link is in this one
// answer only (planque keeps the code's hash), never cached and never sent
// onward, and it's useless without the approving device's signed confirm.
func (s *Server) returnToApp(w http.ResponseWriter, r *http.Request, lang string, l words.Lines, a *Approval, app string, u *url.URL) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if a != nil && phonePlatforms[a.Platform] {
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.WriteHeader(http.StatusOK)
	// u is planque's own PLANQUE_APP_RETURN with the code added, never a
	// caller's: its scheme (tafs:) is let through as is.
	_ = openTmpl.Execute(w, map[string]any{"Lang": lang, "Title": l.SignedIn, "Message": l.SignedInOpen(app),
		"Open": l.OpenTafs, "Link": template.URL(u.String())})
}

var openTmpl = template.Must(template.New("o").Parse(`<!doctype html><html lang="{{.Lang}}"><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1"><title>{{.Title}} · planque</title>
<style>body{font:17px/1.5 -apple-system,system-ui,sans-serif;max-width:32rem;margin:15vh auto;padding:0 16px;color:#1d1d1f}
a{display:inline-block;margin-top:8px;padding:12px 22px;border-radius:10px;background:#1d1d1f;color:#fff;text-decoration:none;font-weight:600}
@media(prefers-color-scheme:dark){body{background:#111;color:#eee}a{background:#eee;color:#111}}</style>
<h1>{{.Title}}</h1><p>{{.Message}}</p><p><a href="{{.Link}}">{{.Open}}</a></p>`))

var pageTmpl = template.Must(template.New("p").Parse(`<!doctype html><html lang="{{.Lang}}"><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1"><title>{{.Title}} · planque</title>
<style>body{font:17px/1.5 -apple-system,system-ui,sans-serif;max-width:32rem;margin:15vh auto;padding:0 16px;color:#1d1d1f}
@media(prefers-color-scheme:dark){body{background:#111;color:#eee}}</style>
<h1>{{.Title}}</h1><p>{{.Message}}</p>`))

func hostOf(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return u
	}
	return p.Host
}

// connectorKey opens the account's KIFF connector key.
func (s *Server) connectorKey(ctx context.Context, tenant string) (string, error) {
	var e Enrolment
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skEnrolment, &e); err != nil {
		return "", err
	}
	return s.openConnector(ctx, e)
}

// connectorStands says whether KIFF still takes e's connector key as
// tenant's connector. KIFF refusing the key is false; KIFF not answering
// is an error, so an outage never swaps the account's key.
func (s *Server) connectorStands(ctx context.Context, tenant string, e Enrolment) (bool, error) {
	key, err := s.openConnector(ctx, e)
	if err != nil {
		return false, err
	}
	me, err := s.KIFF.Me(ctx, key)
	if errors.Is(err, kiff.ErrUnauthorized) {
		return false, nil
	}
	if err != nil {
		return false, apiError{http.StatusBadGateway, "KIFF could not be asked whether the enrolled connector key still works; enrol again"}
	}
	return me.ConnectorProblem(tenant) == "", nil
}

// openConnector opens an enrolment's KIFF connector key.
func (s *Server) openConnector(ctx context.Context, e Enrolment) (string, error) {
	k, err := s.Sealer.Open(ctx, e.Connector, map[string]string{"tenant": e.Tenant, "kind": "connector"})
	return string(k), err
}
