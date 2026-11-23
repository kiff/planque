package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// A hand-over (tafs docs/chrome-panel.md, "A step handed to the person"): a
// taf hands the person one thing to do in the page (a card number, a check,
// a click), the extension shows it, and the step is answered at once as
// handed. The person's Done or Not now comes back later as its own signed
// message, kept here until the taf reads it.

// handOpenFor is how long a hand-over waits for the person; handKeptFor is
// how long its item lives at all, so a taf that never asks leaves nothing.
const (
	handOpenFor  = 30 * time.Minute
	handKeptFor  = time.Hour
	maxHandWords = 280
)

// Outcomes of a hand-over, as the taf reads them.
const (
	handWaiting  = "waiting"
	handDone     = "done"
	handDeclined = "declined"
	handExpired  = "expired"
)

func pkChromeHand(chrome, hand string) string { return "CHROMEHAND#" + chrome + "#" + hand }

// chromeHand is one hand-over: whose it is, what the taf asked, and once the
// person answered, how, and the page after (values blanked).
type chromeHand struct {
	Tenant  string        `json:"tenant"`
	Grant   string        `json:"grant"`
	Run     string        `json:"run"`
	Words   string        `json:"words"`
	Ref     string        `json:"ref,omitempty"`
	Opened  time.Time     `json:"opened"`
	Outcome string        `json:"outcome,omitempty"`
	Page    *browser.Page `json:"page,omitempty"`
}

var (
	errHandUnknown  = errors.New("unknown_hand")
	errHandExpired  = errors.New("expired")
	errHandAnswered = errors.New("already answered")
)

// handWordsRefusal is why a taf's words can't be shown as a hand-over, or "":
// plain text, 1 to 280 characters, with no control or direction characters,
// refused rather than stripped (the extension checks the same).
func handWordsRefusal(words string) string {
	if !utf8.ValidString(words) {
		return "the words aren't text"
	}
	n := utf8.RuneCountInString(words)
	if n < 1 || n > maxHandWords {
		return "the words are 1 to 280 characters"
	}
	for _, r := range words {
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) || r == ' ' || r == ' ' {
			return "the words hold a control or direction character"
		}
	}
	return ""
}

// blankFields is the page with nothing a person typed or chose in it: every
// field's value goes, and with it a checkbox's state and a select's choice
// (both read into Value). The page's own text stays, masked as any read: a page
// that echoes what was typed shows it there, as in any read.
func blankFields(p browser.Page) browser.Page {
	if len(p.Fields) == 0 {
		return p
	}
	fields := make([]browser.Field, len(p.Fields))
	for i, f := range p.Fields {
		f.Value = ""
		fields[i] = f
	}
	p.Fields = fields
	return p
}

// keepHand keeps a hand-over just handed to chrome as step hand.
func (s *Server) keepHand(ctx context.Context, chrome, hand string, h chromeHand) error {
	it := store.Item{PK: pkChromeHand(chrome, hand), SK: "-", TTL: h.Opened.Add(handKeptFor).Unix()}
	return putDoc(ctx, s.KV, it, h, store.Cond{NotExists: true})
}

// answerHand keeps the person's answer to hand, sent by chrome: the first one
// only (a resend is errHandAnswered, which the extension takes as kept), and
// only within the half hour. The page is kept with its values blanked,
// whatever the extension sent.
func (s *Server) answerHand(ctx context.Context, chrome, hand, outcome string, page *browser.Page) error {
	if outcome != handDone && outcome != handDeclined {
		return errHandUnknown
	}
	var h chromeHand
	it, err := getDoc(ctx, s.KV, pkChromeHand(chrome, hand), "-", &h)
	if errors.Is(err, errNotFound) {
		return errHandUnknown
	} else if err != nil {
		return err
	}
	if h.Outcome != "" {
		return errHandAnswered
	}
	if s.now().Sub(h.Opened) >= handOpenFor {
		return errHandExpired
	}
	h.Outcome = outcome
	if page != nil && outcome == handDone {
		blank := blankFields(*page)
		h.Page = &blank
	}
	return putDoc(ctx, s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}, h, store.Cond{Ver: &it.Ver})
}

// readHand is what the taf reads of hand, for its grant only: waiting,
// expired, or the answer, which is then gone (read once).
func (s *Server) readHand(ctx context.Context, chrome, hand, grant string) (chromeHand, error) {
	var h chromeHand
	it, err := getDoc(ctx, s.KV, pkChromeHand(chrome, hand), "-", &h)
	if errors.Is(err, errNotFound) || err == nil && h.Grant != grant {
		return chromeHand{}, errHandUnknown
	} else if err != nil {
		return chromeHand{}, err
	}
	switch {
	case h.Outcome != "":
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
		return h, nil
	case s.now().Sub(h.Opened) >= handOpenFor:
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
		h.Outcome = handExpired
		return h, nil
	default:
		h.Outcome = handWaiting
		return h, nil
	}
}

// browseHand hands the person one thing to do in the page the call opened:
// the words checked, the ref (if any) one of that read's fields or controls,
// then a signed "hand" step. The extension answers it at once as handed; the
// run stays open for the person, and the taf gets the hand-over's id.
func (s *Server) browseHand(ctx context.Context, g *Grant, sess *chromeSession, a browserArgs, page browser.Page) (string, any) {
	if why := handWordsRefusal(a.Words); why != "" {
		return OpNotSent, browserRefusal(CodeBadArguments, why, false)
	}
	if a.Read == "" || browser.Read(page) != a.Read {
		return OpNotSent, browserRefusal(browser.CodePageChanged, "the page isn't the read named ("+a.Read+"): read it again with browse_open first", false)
	}
	if a.Ref != "" && !pageHasRef(page, a.Ref) {
		return OpNotSent, browserRefusal(CodeBadArguments, "ref "+a.Ref+" isn't a field or control of that read", false)
	}
	st := chromeStep{Op: "hand", ReadID: sess.readID, Hand: &chromeHandStep{Words: a.Words, Ref: a.Ref}}
	// The id the step goes out with is the hand-over's: kept before it is sent,
	// so an answer can never come for a hand-over planque doesn't know.
	id := chromeStepID(sess.run, sess.n+1)
	if err := s.keepHand(ctx, sess.chrome, id, chromeHand{Tenant: g.Tenant, Grant: g.ID, Run: sess.run, Words: a.Words, Ref: a.Ref, Opened: s.now()}); err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not keep the hand-over", false)
	}
	if _, err := sess.stepPage(ctx, st); err != nil {
		_, _, _ = s.KV.Delete(ctx, pkChromeHand(sess.chrome, id), "-")
		return OpNotSent, chromeFailed(err, 0)
	}
	sess.kept = true
	return OpDone, map[string]any{"state": "handed", "hand": id,
		"next": "The person does it in their Chrome; read their answer with " + browseHandResult + " and this hand."}
}

// browseHandResult is the person's answer to one of this grant's hand-overs.
func (s *Server) browseHandResult(ctx context.Context, g *Grant, a browserArgs) (string, any) {
	if !chromeGrant(g) {
		return OpNotSent, browserRefusal(CodeBadArguments, browseHandResult+" is the owner's Chrome's only", false)
	}
	if a.Hand == "" || a.URL != "" || a.Read != "" || a.Words != "" || len(a.Steps) > 0 || a.Ref != "" || a.Offset != 0 {
		return OpNotSent, browserRefusal(CodeBadArguments, browseHandResult+" takes hand, the id browse_hand answered, and nothing else", false)
	}
	h, err := s.readHand(ctx, g.Request.Browser.Chrome, a.Hand, g.ID)
	if errors.Is(err, errHandUnknown) {
		return OpNotSent, browserRefusal(CodeBadArguments, "no hand-over "+a.Hand+" of this connection (or its answer was read already)", false)
	} else if err != nil {
		return OpNotSent, browserRefusal(CodeRecordUnavailable, "planque could not read the hand-over", false)
	}
	out := map[string]any{"hand": a.Hand, "state": h.Outcome}
	if h.Page != nil {
		tag := fenceTag()
		out["site"], out["read"], out["page"] = browser.Site(h.Page.URL), browser.Read(*h.Page), browser.Mask(*h.Page, maxPageText)
		return OpDone, fencedAnswer(out, tag)
	}
	return OpDone, out
}

// pageHasRef says ref is one of the page's fields or controls.
func pageHasRef(p browser.Page, ref string) bool {
	for _, f := range p.Fields {
		if f.Ref == ref {
			return true
		}
	}
	for _, c := range p.Controls {
		if c.Ref == ref {
			return true
		}
	}
	return false
}

// labelHand is what the extension signs a hand-over's answer under (kiff/tafs
// chrome/src/crypto.js LABEL.hand): never a step's result.
const labelHand = "tafs-chrome/hand/v1"

// maxHandResult bounds a hand_result message, page and all: the relay's own cap.
const maxHandResult = 120 << 10

// relayHandResult keeps the person's answer to a hand-over, from the Chrome it
// was handed to, signed by that Chrome's key; any ready connection of it, after
// any reconnect. The first answer is kept; a resend is answered as kept.
func (s *Server) relayHandResult(ctx context.Context, connID string, conn chromeConn, body []byte, m map[string]json.RawMessage) {
	var hand string
	_ = json.Unmarshal(m["hand"], &hand)
	refuse := func(code string) {
		s.relayPost(ctx, connID, map[string]any{"type": "hand_refused", "hand": hand, "code": code})
	}
	if len(body) > maxHandResult {
		refuse("too_large")
		return
	}
	var sig string
	var signed struct {
		Hand    string          `json:"hand"`
		Outcome string          `json:"outcome"`
		Page    json.RawMessage `json:"page,omitempty"`
	}
	signed.Hand = hand
	_ = json.Unmarshal(m["outcome"], &signed.Outcome)
	signed.Page = m["page"]
	if json.Unmarshal(m["sig"], &sig) != nil || hand == "" {
		refuse("unsigned")
		return
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChrome(conn.Chrome), &c); err != nil {
		s.relayClose(ctx, connID, closeUnpaired, "this Chrome was unpaired")
		return
	}
	if !stepkey.Verify(c.Key, labelHand, signed, sig) {
		refuse("unsigned")
		return
	}
	var page *browser.Page
	if len(signed.Page) > 0 && string(signed.Page) != "null" {
		var p browser.Page
		if json.Unmarshal(signed.Page, &p) != nil {
			refuse("unsigned")
			return
		}
		page = &p
	}
	switch err := s.answerHand(ctx, conn.Chrome, hand, signed.Outcome, page); {
	case err == nil, errors.Is(err, errHandAnswered):
		s.relayPost(ctx, connID, map[string]any{"type": "hand_kept", "hand": hand})
	case errors.Is(err, errHandExpired):
		refuse("expired")
	case errors.Is(err, errHandUnknown):
		refuse("unknown_hand")
	default:
		s.log().WarnContext(ctx, "chrome hand: not kept", slog.String("err", err.Error()))
		refuse("not_kept")
	}
}
