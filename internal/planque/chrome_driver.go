package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// The owner's own Chrome as a browser grant's worker (tafs docs/chrome.md,
// "Steps"): each session call is one step, signed and sent over the relay
// (chrome_relay.go), and its answer is the page the Chrome signed. planque's
// guard runs before a step leaves, as for the hosted worker, and again on
// the page that comes back: the extension's own checks are the second guard,
// never the only one.

// DriverExtension is a browser grant whose browser is the owner's paired
// Chrome, through the tafs extension. Without it, the grant's browser is the
// hosted worker.
const DriverExtension = "extension"

// The codes planque adds for the owner's Chrome. The extension's own codes
// (docs/chrome.md, "What the extension refuses") pass through as they are.
const (
	// CodeChromeOffline: no proven connection from the grant's Chrome when
	// the step was due, so nothing was sent.
	CodeChromeOffline = "chrome_offline"
	// CodeURLNotRead: the URL isn't a listed site's start page or a link from
	// this grant's latest read; a taf never writes an address of its own.
	CodeURLNotRead = "url_not_read"
	// CodeSitesUnsigned: the grant has no site list signed by the owner's
	// phone, which the extension needs before it runs any step.
	CodeSitesUnsigned = "sites_unsigned"
	// CodeChromeUnpaired: the grant's Chrome isn't paired to this account any
	// more, or the owner's Chrome isn't offered to it.
	CodeChromeUnpaired = "chrome_unpaired"
)

// chromeWait bounds how long a step waits for its signed result. planque's
// call has 30 s in all (its function's timeout), so a step that found the
// Chrome connected waits at most this long, and the call stops sending once
// its own budget, chromeCallBudget, is spent.
const (
	chromeWait       = 20 * time.Second
	chromeCallBudget = 25 * time.Second
	chromePoll       = 150 * time.Millisecond
)

// chromeWaitFor is chromeWait, shorter in tests.
var chromeWaitFor = chromeWait

// afterResultCheck runs, in tests only, between a step's look for its result
// and its look at the connection.
var afterResultCheck func()

// latestReadKept is how long a grant's latest read stays the one its links
// come from: as long as a run's window (docs/chrome.md, "Runs").
const latestReadKept = 2 * time.Hour

// chromeCodeRe is the shape of a code the extension may send back; anything
// else is read as browser_unavailable.
var chromeCodeRe = regexp.MustCompile(`^[a-z_]{1,40}$`)

// chromeNotRun are the extension's codes for a step it refused before doing
// anything: the page is as it was.
var chromeNotRun = map[string]bool{
	"bad_step": true, "bad_url": true, "address_private": true, "site_not_listed": true, "site_blocked": true,
	"needs_site_permission": true, "page_changed": true, "ref_unknown": true, "not_visible": true, "secret_field": true,
	"payment_step": true, "dialog_open": true, "paused_by_owner": true, "paused_failures": true, "paused_no_progress": true, "stopped": true,
}

// chromeRefusal is a step the Chrome answered with an error: its code, and
// the dialog that blocked the page, when one did.
type chromeRefusal struct {
	Code   string
	Detail string
	Dialog *chromeDialog
}

func (r *chromeRefusal) Error() string {
	if r.Detail != "" {
		return r.Code + ": " + r.Detail
	}
	return r.Code
}

// NotRun says the extension refused the step before doing anything.
func (r *chromeRefusal) NotRun() bool { return chromeNotRun[r.Code] }

type chromeDialog struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}

// errChromeNoAnswer: the step was sent, and no signed result came in time.
var errChromeNoAnswer = errors.New("the Chrome didn't answer in time")

// chromeSitesSigned is the phone's signature over the grant's sites
// (tafs-chrome/sites/v1 over {chrome_id, grant, sites}), relayed as the
// phone made it: the extension checks it with the phone key it pinned when
// it was paired.
type chromeSitesSigned struct {
	Sites []string `json:"sites"`
	Sig   string   `json:"sig"`
}

// chromeAct is one act of a step, as the extension runs it.
type chromeAct struct {
	Kind  string `json:"kind"`
	Ref   string `json:"ref"`
	Value string `json:"value,omitempty"`
}

// chromePage is the page a step's result carries: planque's browser.Page,
// plus the extension's read id and a dialog that blocked the page.
type chromePage struct {
	browser.Page
	ReadID string        `json:"read_id"`
	Dialog *chromeDialog `json:"dialog,omitempty"`
	// More and NextOffset page a long read's text (op "read" with offset).
	More       bool `json:"more,omitempty"`
	NextOffset int  `json:"next_offset,omitempty"`
}

// chromeSession is one call's run in the owner's Chrome: its own tab group,
// opened by the first step and closed with the call, as the hosted worker's
// session is one call's browser.
type chromeSession struct {
	s        *Server
	tenant   string
	chrome   string
	grant    string
	run      string
	taf      string
	sites    []string
	signed   chromeSitesSigned
	deadline time.Time
	n        int
	readID   string
	sent     bool
}

// startChrome is a grant's session in its paired Chrome, for the operation
// op: the pairing and the offer are read again here, on every call, so a
// Chrome unpaired or an account no longer offered stops the next step.
func (s *Server) startChrome(ctx context.Context, g *Grant, op, taf string) (*chromeSession, string, string) {
	b := g.Request.Browser
	if b == nil || b.Driver != DriverExtension || b.Chrome == "" {
		return nil, CodeChromeUnpaired, "this connection names no paired Chrome"
	}
	if !s.relayOn() || !s.chromeOnFor(g.Tenant) {
		return nil, CodeChromeUnpaired, "the owner's own Chrome isn't offered to this account now"
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skChrome(b.Chrome), &c); errors.Is(err, errNotFound) {
		return nil, CodeChromeUnpaired, "the Chrome this connection names isn't paired to this account any more"
	} else if err != nil {
		return nil, CodeRecordUnavailable, "planque could not read the Chrome's pairing"
	}
	if g.ChromeSitesSig == "" {
		return nil, CodeSitesUnsigned, "the owner's phone hasn't signed this connection's sites, so the Chrome would run nothing"
	}
	return &chromeSession{s: s, tenant: g.Tenant, chrome: b.Chrome, grant: g.ID, run: op, taf: taf, sites: b.Sites,
		signed: chromeSitesSigned{Sites: b.Sites, Sig: g.ChromeSitesSig}, deadline: time.Now().Add(chromeCallBudget)}, "", ""
}

func (c *chromeSession) ID() string { return c.run }

// Open sends an open step for u and returns the page the Chrome read.
func (c *chromeSession) Open(ctx context.Context, u string) (browser.Page, error) {
	return c.step(ctx, chromeStep{Op: "open", URL: u})
}

// Do sends one act, on the read it was planned from.
func (c *chromeSession) Do(ctx context.Context, st browser.Step) (browser.Page, error) {
	act := chromeAct{}
	switch {
	case st.Fill != "":
		act = chromeAct{Kind: "fill", Ref: st.Fill, Value: st.Value}
	case st.Select != "":
		act = chromeAct{Kind: "select", Ref: st.Select, Value: st.Option}
	case st.Click != "":
		act = chromeAct{Kind: "click", Ref: st.Click}
	default:
		return browser.Page{}, &chromeRefusal{Code: CodeBadArguments, Detail: "a step names nothing to do"}
	}
	raw, _ := json.Marshal(act)
	return c.step(ctx, chromeStep{Op: "act", ReadID: c.readID, Act: raw})
}

// Screenshot: the owner's Chrome sends none, in v1 (tafs Privacy: "it takes
// no screenshots").
func (c *chromeSession) Screenshot(context.Context) ([]byte, error) {
	return nil, errors.New("the owner's Chrome sends no screenshots")
}

// Input and Frame are the hosted sign-in view's: in the owner's Chrome, the
// owner signs in themselves.
func (c *chromeSession) Input(context.Context, browser.Input) error {
	return errors.New("the owner's Chrome has no live view: sign in there")
}

func (c *chromeSession) Frame(context.Context) ([]byte, error) {
	return nil, errors.New("the owner's Chrome has no live view")
}

// Close ends the run's tab group, without waiting: closing only takes power
// away. Nothing is kept for planque to seal: the sign-ins are the Chrome's.
func (c *chromeSession) Close(ctx context.Context) ([]byte, error) {
	if c.sent {
		c.n++
		_ = c.s.sendChromeStep(ctx, c.tenant, c.base(chromeStep{Op: "close"}))
	}
	return nil, nil
}

// base fills in what every step of the run carries.
func (c *chromeSession) base(st chromeStep) chromeStep {
	st.ID = chromeStepID(c.run, c.n)
	st.ChromeID, st.Run, st.Taf, st.Grant, st.Sites = c.chrome, c.run, c.taf, c.grant, c.sites
	signed := c.signed
	st.SitesSigned = &signed
	return st
}

// chromeStepID names a run's nth step, within the relay's 64 characters.
func chromeStepID(run string, n int) string {
	id := fmt.Sprintf("%s.%d", run, n)
	if len(id) <= 64 {
		return id
	}
	h := sha256.Sum256([]byte(run))
	return fmt.Sprintf("op_%s.%d", hex.EncodeToString(h[:])[:40], n)
}

// step sends st and waits for the Chrome's signed result: the page, or a
// chromeRefusal. A Chrome with no proven connection is errChromeOffline,
// with nothing sent; a sent step with no answer in time is errChromeNoAnswer.
func (c *chromeSession) step(ctx context.Context, st chromeStep) (browser.Page, error) {
	p, err := c.stepPage(ctx, st)
	return p.Page, err
}

// stepPage is step, with the page's paging.
func (c *chromeSession) stepPage(ctx context.Context, st chromeStep) (chromePage, error) {
	if time.Now().After(c.deadline) {
		return chromePage{}, errChromeNoAnswer
	}
	c.n++
	st = c.base(st)
	if err := c.s.sendChromeStep(ctx, c.tenant, st); err != nil {
		return chromePage{}, err
	}
	c.sent = true
	until := time.Now().Add(min(chromeWaitFor, time.Until(c.deadline)))
	for {
		// The result first: one the Chrome signed for this step and filed
		// before its connection ended still stands.
		if p, found, err := c.result(ctx, st.ID); found || err != nil {
			return p, err
		}
		if afterResultCheck != nil {
			afterResultCheck()
		}
		// Then the connection: one that ended won't answer, so the step
		// stops now rather than at its timeout. A result filed just before
		// it ended still wins.
		if live, err := c.s.chromeLive(ctx, c.tenant, c.chrome); err == nil && !live {
			if p, found, err := c.result(ctx, st.ID); found || err != nil {
				return p, err
			}
			return chromePage{}, errChromeDropped
		}
		if time.Now().After(until) {
			return chromePage{}, errChromeNoAnswer
		}
		select {
		case <-ctx.Done():
			return chromePage{}, errChromeNoAnswer
		case <-time.After(chromePoll):
		}
	}
}

// result is step id's signed result, once the Chrome filed one.
func (c *chromeSession) result(ctx context.Context, id string) (chromePage, bool, error) {
	body, found, err := c.s.chromeResult(ctx, c.chrome, id)
	if err != nil || !found {
		return chromePage{}, false, err
	}
	p, err := readChromeResult(body)
	if err == nil {
		c.readID = p.ReadID
	}
	return p, true, err
}

// readChromeResult is a signed result's page, or its refusal.
func readChromeResult(body map[string]json.RawMessage) (chromePage, error) {
	var ok bool
	_ = json.Unmarshal(body["ok"], &ok)
	if !ok {
		var e struct {
			Code   string `json:"code"`
			Detail string `json:"detail"`
		}
		_ = json.Unmarshal(body["error"], &e)
		if !chromeCodeRe.MatchString(e.Code) {
			e.Code = browser.CodeBrowserDown
		}
		r := &chromeRefusal{Code: e.Code, Detail: e.Detail}
		if e.Code == "dialog_open" {
			var p chromePage
			if json.Unmarshal(body["page"], &p) == nil && p.Dialog != nil {
				r.Dialog = p.Dialog
			}
		}
		return chromePage{}, r
	}
	var p chromePage
	if err := json.Unmarshal(body["page"], &p); err != nil || p.URL == "" {
		// An answer with no page isn't a page: what the step did isn't known.
		return chromePage{}, &chromeRefusal{Code: browser.CodeBrowserDown, Detail: "the Chrome's answer carried no page"}
	}
	if p.Dialog != nil {
		return chromePage{}, &chromeRefusal{Code: "dialog_open", Dialog: p.Dialog}
	}
	return p, nil
}

// chromeSigRe is the shape of a phone's signature: ECDSA P-256's raw r‖s,
// base64url, as the phones sign every tafs-chrome label.
var chromeSigRe = regexp.MustCompile(`^[A-Za-z0-9_-]{86}$`)

// chromeApprovable says why an approval of g can't carry sig, or "": only an
// extension grant takes the phone's signature over its sites, it needs one,
// and its Chrome must still be paired to the account and offered to it.
func (s *Server) chromeApprovable(ctx context.Context, g *Grant, sig string) string {
	b := g.Request.Browser
	if b == nil || b.Driver != DriverExtension {
		return "chrome_sites_sig is only for a grant that uses the owner's own Chrome"
	}
	if !chromeSigRe.MatchString(sig) {
		return "a grant that uses the owner's own Chrome is approved with chrome_sites_sig, the phone's signature over its sites: nothing was granted"
	}
	if !s.relayOn() || !s.chromeOnFor(g.Tenant) {
		return "the owner's own Chrome isn't offered to this account: nothing was granted"
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skChrome(b.Chrome), &c); err != nil {
		return "the Chrome this grant names isn't paired to this account any more: nothing was granted"
	}
	return ""
}

// pickChrome is the paired Chrome a browser grant asked with none uses: the
// only one paired, or with several the one named (as the owner named it when
// pairing). Never a guess: several and no name, or a name that matches none
// or more than one, is a message saying which names there are.
func (s *Server) pickChrome(ctx context.Context, tenant, name string) (string, string, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "CHROME#")
	if err != nil {
		return "", "", err
	}
	var all []Chrome
	for _, it := range items {
		var c Chrome
		if json.Unmarshal(it.Doc, &c) == nil && chromeIDRe.MatchString(c.ID) {
			all = append(all, c)
		}
	}
	if len(all) == 0 {
		return "", "no Chrome is paired to this account: the owner pairs one on the phone first", nil
	}
	names := make([]string, 0, len(all))
	for _, c := range all {
		names = append(names, label(c.Name))
	}
	name = strings.TrimSpace(name)
	if name == "" {
		if len(all) == 1 {
			return all[0].ID, "", nil
		}
		return "", "several Chromes are paired: name the one this taf uses in browser.chrome_name (" + strings.Join(names, ", ") + ")", nil
	}
	found := ""
	for _, c := range all {
		if strings.EqualFold(label(c.Name), name) {
			if found != "" {
				return "", "two paired Chromes have that name: the owner renames one on the phone first", nil
			}
			found = c.ID
		}
	}
	if found == "" {
		return "", "no paired Chrome has that name (" + strings.Join(names, ", ") + ")", nil
	}
	return found, "", nil
}
