package planque

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/issuing"
	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/store"
)

// A purchase's card (Thune v1, batch B): planque makes one single-use
// virtual card through Stripe Issuing for one purchase case KIFF approved,
// and runs exactly that call. KIFF decides whether to buy, and KIFF's live
// authorization decides each charge: a card in no approved case of KIFF's is
// declined there whatever planque did. planque's part is to make the card
// only under KIFF's permit for that exact case, to keep a record of every
// card it made, and to cancel any card that outlives its case or that it
// has no record of.

// toolIssueCard is the call KIFF signs for a card: its arguments are the
// case, so the permit's hash binds the card to exactly that case.
const toolIssueCard = "issue_card"

// errCardsOff refuses a card when no Issuing connection is configured.
var errCardsOff = errors.New("purchase cards are off: no Stripe Issuing connection is configured")

// pkIssuing holds one record per card planque made, by card id: what the
// sweep checks Stripe's active cards against.
const pkIssuing = "ISSUING"

func skIssuedCard(id string) string { return "CARD#" + id }

// skCaseCard is a case version's card record, by the case: one card per
// case version.
func skCaseCard(tenant, caseID string, version int) string {
	return fmt.Sprintf("CASECARD#%s#%s#%d", tenant, caseID, version)
}

// cardKept is how long a card's record stays after its case ends: long
// enough for reconciliation to tie Stripe's events to the case.
const cardKept = 30 * 24 * time.Hour

// issuedCard is planque's record of a card it made. Never its number.
type issuedCard struct {
	Card    string    `json:"card"`
	Tenant  string    `json:"tenant"`
	Case    string    `json:"case"`
	Version int       `json:"case_version"`
	Permit  string    `json:"permit"`
	Ends    time.Time `json:"ends"`
	Made    time.Time `json:"made"`
}

// cardCall is issue_card's arguments: the purchase case as
// contracts/purchase.json writes it (kiff.purchase.v1), its digest, and
// KIFF's decision on it: approved, and until when.
type cardCall struct {
	Case       json.RawMessage `json:"case"`
	Digest     string          `json:"digest"`
	State      string          `json:"state"`
	ValidUntil time.Time       `json:"valid_until"`
}

// purchaseCase is the part of a kiff.purchase.v1 case a card is made from.
type purchaseCase struct {
	Schema   string `json:"schema"`
	ID       string `json:"case_id"`
	Version  int    `json:"version"`
	Path     string `json:"path"`
	Currency string `json:"currency"`
	Price    struct {
		Max int64 `json:"max"`
	} `json:"price"`
	Grant struct {
		Kind       string   `json:"kind"`
		Cap        int64    `json:"cap"`
		Categories []string `json:"categories"`
		Uses       int      `json:"uses"`
	} `json:"grant"`
	// Merchant's stripe_profile is the seller's Stripe profile a shared
	// payment token is issued to (path token), as KIFF signed it.
	Merchant struct {
		Site          string `json:"site"`
		StripeProfile string `json:"stripe_profile"`
	} `json:"merchant"`
}

// merchantCode is a merchant category code: four digits. Stripe's controls
// take its category names, so a code is refused, never guessed at.
var merchantCode = regexp.MustCompile(`^[0-9]{4}$`)

// cardCase reads issue_card's arguments into the card's case, checking
// what planque can check itself: the digest is the case's own (JCS, as
// tafs computes it), it asks for one single-use card on a payment path,
// its cap is within the case's maximum, and its categories are Stripe's
// names. issuing.Case.Check then checks the state, the amount and the
// window.
func cardCase(raw []byte) (issuing.Case, error) {
	var call cardCall
	if err := json.Unmarshal(raw, &call); err != nil || len(call.Case) == 0 {
		return issuing.Case{}, fmt.Errorf("%w: the arguments aren't a purchase case and its decision", issuing.ErrNotReady)
	}
	sum, err := jcs.SHA256(call.Case)
	if err != nil || sum != call.Digest {
		return issuing.Case{}, fmt.Errorf("%w: the digest isn't the case's", issuing.ErrNotReady)
	}
	var pc purchaseCase
	if err := json.Unmarshal(call.Case, &pc); err != nil || pc.Schema != "kiff.purchase.v1" {
		return issuing.Case{}, fmt.Errorf("%w: the case isn't a kiff.purchase.v1 case", issuing.ErrNotReady)
	}
	switch {
	case pc.Grant.Kind != "single_use_card" || pc.Grant.Uses != 1:
		return issuing.Case{}, fmt.Errorf("%w: it asks for no single-use card", issuing.ErrNotReady)
	case pc.Path != "token" && pc.Path != "card_fill":
		return issuing.Case{}, fmt.Errorf("%w: it names no payment path a card serves", issuing.ErrNotReady)
	case pc.Price.Max <= 0 || pc.Grant.Cap > pc.Price.Max:
		return issuing.Case{}, fmt.Errorf("%w: the card's cap is above the case's maximum", issuing.ErrNotReady)
	}
	for _, cat := range pc.Grant.Categories {
		if merchantCode.MatchString(cat) {
			return issuing.Case{}, fmt.Errorf("%w: %s is a merchant category code; the card's controls take Stripe's category name", issuing.ErrNotReady, cat)
		}
	}
	return issuing.Case{ID: pc.ID, Version: pc.Version, State: call.State, MaxAmount: pc.Grant.Cap,
		Currency: strings.ToLower(pc.Currency), Categories: pc.Grant.Categories, Expires: call.ValidUntil}, nil
}

// cardsAudience is the audience a card permit names.
func (s *Server) cardsAudience() string { return s.BaseURL + "/thune/cards" }

// issueCard makes the card for the purchase case in caseJSON (issue_card's
// arguments, cardCall) of tenant, under
// token, KIFF's permit for exactly that call. Nothing goes to Stripe unless
// the permit verifies, is live, and names this tenant, this tool and these
// exact arguments, and the case it carries is approved and within its
// window. A card that can't be recorded is canceled at once: planque never
// leaves a card it has no record of.
func (s *Server) issueCard(ctx context.Context, tenant, token string, caseJSON []byte) (issuing.Card, error) {
	if s.Issuing == nil {
		return issuing.Card{}, errCardsOff
	}
	p, err := s.Permits.Verify(ctx, token, permit.Want{Audience: s.cardsAudience(), Tenant: tenant, Tool: toolIssueCard, Args: caseJSON})
	if err != nil {
		return issuing.Card{}, err
	}
	now := s.now()
	if !p.Live(now) {
		return issuing.Card{}, errors.New("the permit for this card has expired: ask KIFF for the case again")
	}
	c, err := cardCase(caseJSON)
	if err != nil {
		return issuing.Card{}, err
	}
	card, err := s.Issuing.Create(ctx, c, p.ID, now)
	if err != nil {
		return issuing.Card{}, err
	}
	rec := issuedCard{Card: card.ID, Tenant: tenant, Case: c.ID, Version: c.Version, Permit: p.ID, Ends: c.Expires.UTC(), Made: now.UTC()}
	it := store.Item{PK: pkIssuing, SK: skIssuedCard(card.ID), TTL: c.Expires.Add(cardKept).Unix()}
	err = putDoc(ctx, s.KV, it, rec, store.Cond{})
	if err == nil {
		// The case's card, by case and version: what its token is made from.
		err = putDoc(ctx, s.KV, store.Item{PK: pkIssuing, SK: skCaseCard(tenant, c.ID, c.Version), TTL: it.TTL}, rec, store.Cond{})
	}
	if err != nil {
		if cerr := s.Issuing.Cancel(ctx, card.ID); cerr != nil {
			s.log().ErrorContext(ctx, "issuing: a card without a record could not be canceled", slog.String("card", card.ID), slog.String("case", c.ID), slog.String("err", cerr.Error()))
		}
		return issuing.Card{}, fmt.Errorf("the card's record couldn't be kept, so the card was canceled: %w", err)
	}
	s.log().InfoContext(ctx, "issuing: card made", slog.String("card", card.ID), slog.String("tenant", tenant), slog.String("case", c.ID), slog.Int("version", c.Version), slog.String("permit", p.ID))
	return card, nil
}

// sweepCards cancels every active card that outlived its case, and every
// active card planque has no record of (one made outside planque, with the
// key or by hand): that one is also logged as an alert. It returns how many
// it canceled. KIFF's live authorization already declines both; this ends
// them at Stripe too.
func (s *Server) sweepCards(ctx context.Context) (int, error) {
	if s.Issuing == nil {
		return 0, nil
	}
	// Every page is read; on an error the cards read so far are still
	// swept, and the error is reported.
	active, err := s.Issuing.Active(ctx, 100)
	now, canceled := s.now(), 0
	first := err
	for _, card := range active {
		var rec issuedCard
		_, err := getDoc(ctx, s.KV, pkIssuing, skIssuedCard(card.ID), &rec)
		switch {
		case errors.Is(err, errNotFound):
			s.log().ErrorContext(ctx, "issuing alert: an active card planque has no record of, canceled", slog.String("card", card.ID), slog.String("case", card.Case))
		case err != nil:
			first = cmpErr(first, err)
			continue
		case now.Before(rec.Ends):
			continue // its case is still running
		}
		if err := s.Issuing.Cancel(ctx, card.ID); err != nil {
			first = cmpErr(first, err)
			continue
		}
		canceled++
	}
	return canceled, first
}
