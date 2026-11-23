package planque

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// A sealed fill (tafs docs/chrome.md, "Sealed values"): a taf's browse_act
// whose one step fills a sign-in saved in the owner's Chrome, named by its
// handle. KIFF holds every one for the owner, and the owner's phone releases
// it with its own signature over the exact field the owner saw: the Chrome,
// the step, the handle, the origin, the field's ref and the read it is on.
//
// The extension names each read anew, so the call the owner allows must act
// on the read the owner was shown. So the gateway, before KIFF holds the
// call, asks planque to describe it: a plain browse_open of the page that
// KIFF allowed on its own (op "<operation>:open"), carrying the gateway's
// SealedDescribeHeader. planque opens the page as the sealed call's own
// run's step 1, keeps that run open, and answers with the block the owner's
// card draws, signed with its step key. The allowed call (the operation
// itself, with the phone's release in SealedReleaseHeader) then sends step 2
// on that read, and nothing else. Every other case fails closed: the Chrome
// fills nothing without a release that fits.

// The gateway's headers for a sealed fill. A taf can't set either: the
// gateway sends only its own.
const (
	// SealedDescribeHeader asks a browse_open to describe a sealed fill:
	// base64url JSON of sealedDescribe.
	SealedDescribeHeader = "KIFF-Sealed-Describe"
	// SealedReleaseHeader carries the owner's phone's release on the allowed
	// sealed call: base64url JSON, as the phone made it.
	SealedReleaseHeader = "KIFF-Sealed-Release"
)

// labelDescribe is what planque signs a sealed fill's block under, with its
// step key: KIFF pins the public half.
const labelDescribe = "planque/describe/v1"

// The codes of a sealed fill.
const (
	// CodeSealedOff: sealed fills aren't turned on here.
	CodeSealedOff = "sealed_fill_off"
	// CodeSealedBusy: this Chrome already has as many sealed fills waiting
	// for their owner as it may.
	CodeSealedBusy = "sealed_busy"
	// CodeReleaseMissing: no release, or no page kept open for it: nothing
	// was filled.
	CodeReleaseMissing = "release_missing"
	// CodeReleaseStale: the release, or the call, isn't for the field the
	// owner was shown: nothing was filled.
	CodeReleaseStale = "release_stale"
)

// sealedKept is how long a described sealed fill's page stays open for its
// owner: as long as its hold, within the extension's 5-minute idle limit.
// maxSealedKept is how many may wait on one Chrome at a time.
const (
	sealedKept    = 3 * time.Minute
	maxSealedKept = 2
)

// sealedFillStep is the act's step in its run: the describe's open is 1.
const sealedFillStep = 2

// sealedDescribe is what the gateway asks a describe for: the sealed call's
// operation, and the read, field and handle its one step names.
type sealedDescribe struct {
	Op     string `json:"op"`
	Read   string `json:"read"`
	Fill   string `json:"fill"`
	Sealed string `json:"sealed"`
}

type sealedHeadersKey struct{}

type sealedHeaders struct {
	describe *sealedDescribe
	release  json.RawMessage
	bad      string
}

// withSealedHeaders keeps a call's sealed-fill headers for callChrome. A
// header that doesn't read is kept as bad: the call it came with is refused.
func withSealedHeaders(ctx context.Context, h http.Header) context.Context {
	var sh sealedHeaders
	if v := strings.TrimSpace(h.Get(SealedDescribeHeader)); v != "" {
		var d sealedDescribe
		if raw, err := base64.RawURLEncoding.DecodeString(v); err != nil || json.Unmarshal(raw, &d) != nil || d.Op == "" {
			sh.bad = "the describe header doesn't read"
		} else {
			sh.describe = &d
		}
	}
	if v := strings.TrimSpace(h.Get(SealedReleaseHeader)); v != "" {
		raw, err := base64.RawURLEncoding.DecodeString(v)
		if err != nil || !json.Valid(raw) {
			sh.bad = "the release header doesn't read"
		} else {
			sh.release = raw
		}
	}
	if sh.describe == nil && sh.release == nil && sh.bad == "" {
		return ctx
	}
	return context.WithValue(ctx, sealedHeadersKey{}, sh)
}

func sealedHeadersFrom(ctx context.Context) sealedHeaders {
	sh, _ := ctx.Value(sealedHeadersKey{}).(sealedHeaders)
	return sh
}

// sealedRun is a described sealed fill's run, kept open in the owner's
// Chrome for the one act its owner may allow.
type sealedRun struct {
	Tenant string `json:"tenant"`
	Grant  string `json:"grant"`
	Agent  string `json:"agent"`
	Op     string `json:"op"`
	Chrome string `json:"chrome"`
	// URL is the address the call opened, as the allowed call names it.
	URL      string    `json:"url"`
	Origin   string    `json:"origin"`
	Read     string    `json:"read"`
	ReadID   string    `json:"read_id"`
	FieldRef string    `json:"field_ref"`
	Sealed   string    `json:"sealed"`
	Opened   time.Time `json:"opened"`
}

func skSealedPrefix(chrome string) string { return "SEALEDRUN#" + chrome + "#" }

// skSealedRun keys a kept run by the Chrome, the connection and the sealed
// call's operation: only that connection's allowed call finds it.
func skSealedRun(chrome, grant, op string) string { return skSealedPrefix(chrome) + grant + "#" + op }

// sealedBlock is the hold's browser block: what the owner's card draws, and
// the release's fields (KIFF's handlers/browser_block.go reads exactly these).
type sealedBlock struct {
	V       string `json:"v"`
	Tenant  string `json:"tenant"`
	Op      string `json:"op"`
	Site    string `json:"site"`
	Origin  string `json:"origin"`
	Action  string `json:"action_words"`
	Context struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"context"`
	Chrome struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"chrome"`
	Release struct {
		ChromeID string `json:"chrome_id"`
		Op       string `json:"op"`
		Sealed   string `json:"sealed"`
		Origin   string `json:"origin"`
		FieldRef string `json:"field_ref"`
		ReadID   string `json:"read_id"`
	} `json:"release_fields"`
}

// The block's words are cut short: the owner reads the field and where it
// is, never the page.
const (
	maxBlockTitle = 120
	maxBlockLabel = 80
	maxBlockText  = 280
)

// sealedRefusal is why a call with sealed steps can't go as one, or "": only
// a browse_act in the owner's Chrome, while sealed fills are on, whose one
// step is the sealed fill.
func (s *Server) sealedRefusal(g *Grant, tool string, steps []browser.Step) (string, string) {
	sealed := slices.ContainsFunc(steps, func(st browser.Step) bool { return st.Sealed != "" })
	switch {
	case !sealed:
		return "", ""
	case !s.SealedFill:
		return CodeSealedOff, "filling a saved sign-in isn't turned on"
	case !chromeGrant(g) || tool != browseAct:
		return CodeSealedOff, "a saved sign-in is filled only by browse_act, in the owner's Chrome"
	case len(steps) != 1:
		return CodeBadArguments, "a call that fills a saved sign-in has that one step and nothing else"
	}
	return "", ""
}

// keptSealedRuns are the Chrome's kept runs still within their time. Any
// past it is closed and dropped now: planque's own sweep is hourly, so the
// next sealed call on that Chrome ends the ones its owner never answered,
// and the extension ends an idle run after 5 minutes in any case.
func (s *Server) keptSealedRuns(ctx context.Context, g *Grant, chrome string) ([]sealedRun, error) {
	items, err := s.KV.Query(ctx, tenantPK(g.Tenant), skSealedPrefix(chrome))
	if err != nil {
		return nil, err
	}
	var live []sealedRun
	for _, it := range items {
		var r sealedRun
		if json.Unmarshal(it.Doc, &r) != nil {
			continue
		}
		if s.now().Sub(r.Opened) < sealedKept {
			live = append(live, r)
			continue
		}
		if _, took, err := s.KV.Delete(ctx, it.PK, it.SK); err == nil && took {
			s.closeSealedRun(ctx, r)
		}
	}
	return live, nil
}

// closeSealedRun ends a kept run in the owner's Chrome, without waiting:
// closing only takes power away.
func (s *Server) closeSealedRun(ctx context.Context, r sealedRun) {
	// The run's steps were 1 (open) and maybe 2 (the fill): close is 3.
	_ = s.sendChromeStep(ctx, r.Tenant, chromeStep{ID: chromeStepID(r.Op, sealedFillStep+1), ChromeID: r.Chrome, Run: r.Op, Grant: r.Grant, Op: "close"})
}

// describeSealed answers the gateway's describe of a sealed fill: a plain
// browse_open that KIFF allowed as "<operation>:open", run as the sealed
// call's own step 1 and kept open, answered with the signed block. The taf
// is told nothing: the answer goes to the gateway.
func (s *Server) describeSealed(ctx context.Context, g *Grant, p permit.Permit, a browserArgs, d *sealedDescribe) (string, any) {
	part := browser.SealedPart(d.Sealed)
	switch {
	case !s.SealedFill:
		return OpNotSent, browserRefusal(CodeSealedOff, "filling a saved sign-in isn't turned on", false)
	case p.Op != d.Op+":open":
		// KIFF's permit is for this describe's own open, and no other.
		return OpNotSent, browserRefusal(CodePermitInvalid, "the describe's permit isn't for this sealed call's open", false)
	case part == "" || d.Fill == "" || d.Read == "":
		return OpNotSent, browserRefusal(CodeBadArguments, "a describe names the sealed call's read, field and handle", false)
	}
	chrome := g.Request.Browser.Chrome
	kept, err := s.keptSealedRuns(ctx, g, chrome)
	if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read the Chrome's waiting fills", false)
	}
	if slices.ContainsFunc(kept, func(r sealedRun) bool { return r.Grant == g.ID && r.Op == d.Op }) {
		// One describe per hold: the page is already open for it.
		return OpNotSent, browserRefusal(CodeSealedBusy, "this sealed fill was described already, and its page is open for the owner", false)
	}
	if len(kept) >= maxSealedKept {
		return OpNotSent, browserRefusal(CodeSealedBusy, "this Chrome already has as many saved sign-ins waiting for their owner as it may; try again once they're answered", false)
	}
	taf := cmpStr(g.Request.AgentName, g.Request.Agent)
	sess, code, why := s.startChrome(ctx, g, d.Op, taf)
	if code != "" {
		return OpNotSent, browserRefusal(code, why, false)
	}
	defer func() { _, _ = sess.Close(ctx) }()
	if !p.Live(s.now()) {
		return OpNotSent, toolResult("Not sent: the permit expired before the call could be sent.", true)
	}
	cp, err := sess.stepPage(ctx, chromeStep{Op: "open", URL: a.URL})
	if err != nil {
		return OpNotSent, chromeFailed(err, 0)
	}
	page := cp.Page
	if code, why := browser.Allowed(page.URL, g.Request.Browser.Sites); code != "" {
		return OpNotSent, browserRefusal(code, "the page redirected: "+why, false)
	}
	// The field the owner is shown is the one the taf named, on the read it
	// named: the same read gives each ref the same element.
	if browser.Read(page) != d.Read {
		return OpNotSent, browserRefusal(browser.CodePageChanged, "the page isn't the read the sealed step was built from: read it again first", false)
	}
	i := slices.IndexFunc(page.Fields, func(f browser.Field) bool { return f.Ref == d.Fill })
	if i < 0 {
		return OpNotSent, browserRefusal(browser.CodeRefUnknown, "the page has no field "+d.Fill, false)
	}
	f := page.Fields[i]
	if part == "password" && f.Kind != "password" || part == "username" && f.Kind != "text" && f.Kind != "email" {
		return OpNotSent, browserRefusal(browser.CodeSecretField, "a saved "+part+" goes only into a "+map[string]string{"password": "password", "username": "username or email"}[part]+" field", false)
	}
	u, err := url.Parse(page.URL)
	if err != nil || u.Scheme != "https" {
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "the page isn't on an https origin", false)
	}
	origin := u.Scheme + "://" + u.Host
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skChrome(chrome), &c); err != nil {
		return OpNotSent, browserRefusal(CodeChromeUnpaired, "the Chrome this connection names isn't paired to this account any more", false)
	}
	var b sealedBlock
	b.V, b.Tenant, b.Op, b.Site, b.Origin = labelDescribe, g.Tenant, chromeStepID(d.Op, sealedFillStep), browser.Site(page.URL), origin
	b.Action = map[string]string{"password": "Fill your saved password into ", "username": "Fill your saved username or email into "}[part] +
		"“" + blockWords(f.Label, maxBlockLabel) + "”"
	b.Context.Title = blockWords(page.Title, maxBlockTitle)
	b.Context.Text = blockWords(browser.Mask(page, maxBlockText*4).Text, maxBlockText)
	b.Chrome.ID, b.Chrome.Name = chrome, label(c.Name)
	b.Release.ChromeID, b.Release.Op, b.Release.Sealed, b.Release.Origin, b.Release.FieldRef, b.Release.ReadID =
		chrome, b.Op, d.Sealed, origin, d.Fill, sess.readID
	sig, err := stepkey.SignBody(ctx, s.ChromeSigner, labelDescribe, b)
	if err != nil {
		return OpNotSent, browserRefusal(browser.CodeBrowserDown, "planque could not sign the fill's description", false)
	}
	run := sealedRun{Tenant: g.Tenant, Grant: g.ID, Agent: g.Request.Agent, Op: d.Op, Chrome: chrome, URL: a.URL, Origin: origin,
		Read: d.Read, ReadID: sess.readID, FieldRef: d.Fill, Sealed: d.Sealed, Opened: s.now().UTC()}
	it := store.Item{PK: tenantPK(g.Tenant), SK: skSealedRun(chrome, g.ID, d.Op), TTL: s.now().Add(sealedKept + time.Minute).Unix()}
	if err := putDoc(ctx, s.KV, it, run, store.Cond{NotExists: true}); err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not keep the page open for the owner", false)
	}
	// Kept: the run stays open in the owner's Chrome for the allowed call.
	sess.kept = true
	s.log().InfoContext(ctx, "sealed fill described", slog.String("grant", g.ID), slog.String("op", d.Op), slog.String("site", b.Site))
	block, _ := json.Marshal(b)
	res := toolResult("Described for the owner's card.", false)
	res["structuredContent"] = map[string]any{"browser_block": map[string]any{"block": json.RawMessage(block), "sig": sig}}
	return OpDone, res
}

// sealedRelease is the phone's release, as the extension takes it.
type sealedRelease struct {
	ChromeID string `json:"chrome_id"`
	Op       string `json:"op"`
	Origin   string `json:"origin"`
	Nonce    string `json:"nonce"`
	Exp      int64  `json:"exp"`
	Sig      string `json:"sig"`
}

// actSealed sends an allowed sealed fill: step 2 of the run its describe kept,
// on the read the owner was shown, with the phone's release. Taken once:
// the kept run is gone as it's read, so a second call finds none.
func (s *Server) actSealed(ctx context.Context, g *Grant, p permit.Permit, a browserArgs, release json.RawMessage) (string, any) {
	st := a.Steps[0]
	chrome := g.Request.Browser.Chrome
	if _, err := s.keptSealedRuns(ctx, g, chrome); err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read the Chrome's waiting fills", false)
	}
	it, took, err := s.KV.Delete(ctx, tenantPK(g.Tenant), skSealedRun(chrome, g.ID, p.Op))
	if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read the page kept for the owner", false)
	}
	var r sealedRun
	if !took || json.Unmarshal(it.Doc, &r) != nil || s.now().Sub(r.Opened) >= sealedKept {
		if took {
			s.closeSealedRun(ctx, r)
		}
		return OpNotSent, browserRefusal(CodeReleaseMissing, "not filled: the page the owner was shown is no longer open for it; ask again", false)
	}
	// From here, whatever happens, the kept run ends with this call.
	stale := func(why string) (string, any) {
		s.closeSealedRun(ctx, r)
		return OpNotSent, browserRefusal(CodeReleaseStale, "not filled: "+why+"; ask again", false)
	}
	if len(release) == 0 {
		s.closeSealedRun(ctx, r)
		return OpNotSent, browserRefusal(CodeReleaseMissing, "not filled: the call came without the owner's release", false)
	}
	var rel sealedRelease
	if json.Unmarshal(release, &rel) != nil {
		return stale("the release doesn't read")
	}
	switch {
	case r.Agent != g.Request.Agent:
		return stale("the page was kept for another taf")
	case a.URL != r.URL || a.Read != r.Read || st.Fill != r.FieldRef || st.Sealed != r.Sealed:
		return stale("the call isn't the fill the owner was shown")
	case rel.ChromeID != r.Chrome || rel.Op != chromeStepID(r.Op, sealedFillStep) || rel.Origin != r.Origin:
		return stale("the release is for another Chrome, step or site")
	case rel.Exp <= s.now().UnixMilli():
		return stale("the release expired")
	}
	taf := cmpStr(g.Request.AgentName, g.Request.Agent)
	sess, code, why := s.startChrome(ctx, g, r.Op, taf)
	if code != "" {
		s.closeSealedRun(ctx, r)
		return OpNotSent, browserRefusal(code, why, false)
	}
	// The run continues where the describe left it: its step 1 was the open.
	sess.n, sess.readID, sess.sent = sealedFillStep-1, r.ReadID, true
	defer func() { _, _ = sess.Close(ctx) }()
	if !p.Live(s.now()) {
		return OpNotSent, toolResult("Not sent: the permit expired before the call could be sent.", true)
	}
	act, _ := json.Marshal(chromeAct{Kind: "fill", Ref: r.FieldRef, Sealed: r.Sealed})
	cp, err := sess.stepPage(ctx, chromeStep{Op: "act", ReadID: r.ReadID, Act: act, Release: release})
	if cr := (*chromeRefusal)(nil); errors.Is(err, errChromeOffline) || errors.As(err, &cr) && cr.NotRun() {
		return OpNotSent, chromeFailed(err, 0)
	}
	if err != nil {
		var cr *chromeRefusal
		if errors.As(err, &cr) && (cr.Code == CodeReleaseMissing || cr.Code == "sealed_unknown") {
			// The extension refused before filling anything.
			return OpNotSent, chromeFailed(err, 0)
		}
		s.log().WarnContext(ctx, "sealed fill failed", slog.String("grant", g.ID), slog.String("op", p.Op), slog.String("err", cut(err.Error(), 200)))
		return OpUnknown, toolResult("The fill failed in the owner's Chrome ("+chromeUnknownWhy(err)+"): what the page did isn't known. planque will not fill it again; read the page before a new call.", true)
	}
	page := cp.Page
	s.keepChromeRead(ctx, g, page, sess.readID, "")
	return OpDone, fencedAnswer(map[string]any{"site": browser.Site(page.URL), "steps_run": 1, "read": browser.Read(page), "page": browser.Mask(page, maxPageText)}, fenceTag())
}

// blockWords is s on one line, with no control or direction characters, cut
// to n bytes: a page's words, drawn as words on the owner's card.
func blockWords(s string, n int) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) {
			return ' '
		}
		return r
	}, s)
	return cut(strings.Join(strings.Fields(s), " "), n)
}

func cmpStr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
