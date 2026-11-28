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
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/spt"
	"github.com/kiff/planque/internal/store"
)

// A purchase's shared payment token (Thune v1, batch B1, path a): for a case
// whose path is token, planque turns the case's single-use Issuing card into
// a card PaymentMethod (Stripe's preview; no number or code is read) and
// issues one shared payment token from it, scoped to exactly the case: the
// seller's Stripe profile KIFF signed, the card's cap, the case's currency,
// and ending at the decision's valid_until. As for the card, planque runs
// only the call KIFF signed. The token is a bearer credential for that
// seller: it's kept sealed for its life, never logged, never answered to the
// taf, and revoked when the case ends, used or not.

// toolIssueToken is the call KIFF signs for a case's token: its arguments
// are the case, as issue_card's are.
const toolIssueToken = "issue_token"

// errTokensOff refuses a token when no token connection is configured.
var errTokensOff = errors.New("purchase tokens are off: no Stripe connection for shared payment tokens is configured")

// stripeProfile is a Stripe business network profile's id.
var stripeProfile = regexp.MustCompile(`^profile_[A-Za-z0-9_]{1,200}$`)

func skCaseToken(tenant, caseID string, version int) string {
	return fmt.Sprintf("TOKEN#%s#%s#%d", tenant, caseID, version)
}

// issuedToken is planque's record of a token: its id sealed, never in the
// clear, and when its case ends.
type issuedToken struct {
	Sealed  []byte    `json:"sealed"`
	Tenant  string    `json:"tenant"`
	Case    string    `json:"case"`
	Version int       `json:"case_version"`
	Card    string    `json:"card"`
	Ends    time.Time `json:"ends"`
	Made    time.Time `json:"made"`
}

func purchaseTokenBinding(tenant, caseID string, version int) map[string]string {
	return map[string]string{"tenant": tenant, "kind": "purchase_token", "case": caseID, "version": fmt.Sprint(version)}
}

// tokensAudience is the audience a token permit names.
func (s *Server) tokensAudience() string { return s.BaseURL + "/thune/tokens" }

// tokenAnswer is what a token's issue says: never its id.
type tokenAnswer struct {
	State     string    `json:"state"`
	Seller    string    `json:"seller"`
	MaxAmount int64     `json:"max_amount"`
	Currency  string    `json:"currency"`
	ExpiresAt time.Time `json:"expires_at"`
}

// issueToken issues the shared payment token for the purchase case in
// caseJSON (issue_token's arguments, as issue_card's) of tenant, under
// token, KIFF's permit for exactly that call. Nothing goes to Stripe unless
// the permit verifies and names these exact arguments; the case checks as
// a card's does, its path is token, it names the seller's Stripe profile,
// and its card already exists, unexpired. A token that can't be kept is
// revoked at once.
func (s *Server) issueToken(ctx context.Context, tenant, token string, caseJSON []byte) (tokenAnswer, error) {
	if s.Issuing == nil || s.Tokens == nil {
		return tokenAnswer{}, errTokensOff
	}
	p, err := s.Permits.Verify(ctx, token, permit.Want{Audience: s.tokensAudience(), Tenant: tenant, Tool: toolIssueToken, Args: caseJSON})
	if err != nil {
		return tokenAnswer{}, err
	}
	now := s.now()
	if !p.Live(now) {
		return tokenAnswer{}, errors.New("the permit for this token has expired: ask KIFF for the case again")
	}
	c, err := cardCase(caseJSON)
	if err != nil {
		return tokenAnswer{}, err
	}
	if err := c.Check(now); err != nil {
		return tokenAnswer{}, err
	}
	var call cardCall
	var pc purchaseCase
	_ = json.Unmarshal(caseJSON, &call)
	_ = json.Unmarshal(call.Case, &pc)
	switch {
	case pc.Path != "token":
		return tokenAnswer{}, fmt.Errorf("%w: its path is %q, not token: a card fill never gets a token", issuing.ErrNotReady, pc.Path)
	case !stripeProfile.MatchString(pc.Merchant.StripeProfile):
		return tokenAnswer{}, fmt.Errorf("%w: it names no seller's Stripe profile, so no token can be scoped to the seller", issuing.ErrNotReady)
	}
	var card issuedCard
	if _, err := getDoc(ctx, s.KV, pkIssuing, skCaseCard(tenant, c.ID, c.Version), &card); errors.Is(err, errNotFound) {
		return tokenAnswer{}, fmt.Errorf("%w: the case has no card yet: issue_card comes first", issuing.ErrNotReady)
	} else if err != nil {
		return tokenAnswer{}, err
	}
	if !now.Before(card.Ends) {
		return tokenAnswer{}, fmt.Errorf("%w: the case's card has ended", issuing.ErrNotReady)
	}
	pm, err := s.Issuing.PaymentMethod(ctx, card.Card)
	if err != nil {
		return tokenAnswer{}, err
	}
	tok, err := s.Tokens.Issue(ctx, spt.Issue{PaymentMethod: pm, SellerProfile: pc.Merchant.StripeProfile, Currency: c.Currency,
		MaxAmount: c.MaxAmount, ExpiresAt: c.Expires, IdempotencyKey: fmt.Sprintf("thune-tok-%s-v%d", c.ID, c.Version), Purchase: c.ID})
	if err != nil {
		return tokenAnswer{}, err
	}
	sealed, err := s.Sealer.Seal(ctx, []byte(tok.ID), purchaseTokenBinding(tenant, c.ID, c.Version))
	if err == nil {
		rec := issuedToken{Sealed: sealed, Tenant: tenant, Case: c.ID, Version: c.Version, Card: card.Card, Ends: c.Expires.UTC(), Made: now.UTC()}
		err = putDoc(ctx, s.KV, store.Item{PK: pkIssuing, SK: skCaseToken(tenant, c.ID, c.Version), TTL: c.Expires.Add(cardKept).Unix()}, rec, store.Cond{})
	}
	if err != nil {
		if rerr := s.Tokens.Revoke(ctx, tok.ID); rerr != nil {
			s.log().ErrorContext(ctx, "purchase token: one without a record could not be revoked", slog.String("case", c.ID), slog.String("err", rerr.Error()))
		}
		return tokenAnswer{}, fmt.Errorf("the token's record couldn't be kept, so the token was revoked: %w", err)
	}
	s.log().InfoContext(ctx, "purchase token issued", slog.String("tenant", tenant), slog.String("case", c.ID), slog.Int("version", c.Version), slog.String("card", card.Card))
	return tokenAnswer{State: "issued", Seller: pc.Merchant.StripeProfile, MaxAmount: c.MaxAmount, Currency: strings.ToUpper(c.Currency), ExpiresAt: c.Expires.UTC()}, nil
}

// sweepTokens revokes every token whose case has ended, and forgets it. It
// returns how many it revoked.
func (s *Server) sweepTokens(ctx context.Context) (int, error) {
	if s.Tokens == nil {
		return 0, nil
	}
	items, err := s.KV.Query(ctx, pkIssuing, "TOKEN#")
	if err != nil {
		return 0, err
	}
	now, revoked := s.now(), 0
	var first error
	for _, it := range items {
		var rec issuedToken
		if json.Unmarshal(it.Doc, &rec) != nil || now.Before(rec.Ends) {
			continue
		}
		id, err := s.Sealer.Open(ctx, rec.Sealed, purchaseTokenBinding(rec.Tenant, rec.Case, rec.Version))
		if err != nil {
			first = cmpErr(first, err)
			continue
		}
		if err := s.Tokens.Revoke(ctx, string(id)); err != nil {
			// A token already used up, expired or revoked can't be revoked
			// again: read it, and forget it once it's no longer active.
			if t, gerr := s.Tokens.Get(ctx, string(id)); gerr != nil || t.Status == "active" {
				first = cmpErr(first, err)
				continue
			}
		}
		clear(id)
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
		revoked++
	}
	return revoked, first
}
