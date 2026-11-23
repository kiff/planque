package planque

import (
	"context"
	"errors"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kiff/planque/internal/browser"
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
