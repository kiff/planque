package permit

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

// HaltConfirmHeader carries KIFF's ask that planque confirm the owner's halt
// (kiff-cloud#1061): a JWT signed with the permit key, of a typ no permit
// has, so it is never mistaken for one.
const HaltConfirmHeader = "KIFF-Halt-Confirm"

const (
	haltAskTyp = "kiff-halt-confirm+jwt"
	// haltAskLifetime is the longest exp-iat an ask may have.
	haltAskLifetime = 60 * time.Second
)

// HaltAsk is the halt KIFF asks planque to confirm: the account, the agent
// ("" for the whole account) and when it was halted.
type HaltAsk struct {
	Tenant, Agent string
	HaltedAt      time.Time
}

// ErrHaltAsk refuses an ask that isn't KIFF's, isn't live, or names another
// halt.
var ErrHaltAsk = errors.New("this is not KIFF's signed ask to confirm that halt")

// VerifyHaltAsk checks that token is KIFF's live ask to confirm want.
func (v *Verifier) VerifyHaltAsk(ctx context.Context, token string, want HaltAsk) error {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return ErrHaltAsk
	}
	var h header
	if !decodePart(parts[0], &h) || h.Typ != haltAskTyp || h.Alg != alg || h.Kid == "" {
		return ErrHaltAsk
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return ErrHaltAsk
	}
	key, err := v.Keys.Key(ctx, h.Kid)
	if err != nil || !ed25519.Verify(key, []byte(parts[0]+"."+parts[1]), sig) {
		return ErrHaltAsk
	}
	var c struct {
		Iss      string `json:"iss"`
		Iat      *int64 `json:"iat"`
		Exp      *int64 `json:"exp"`
		Tenant   string `json:"tenant"`
		Agent    string `json:"agent"`
		HaltedAt string `json:"halted_at"`
	}
	if !decodePart(parts[1], &c) || c.Iat == nil || c.Exp == nil || c.Iss != v.Issuer {
		return ErrHaltAsk
	}
	at, err := time.Parse(time.RFC3339Nano, c.HaltedAt)
	if err != nil || c.Tenant != want.Tenant || c.Agent != want.Agent || !at.Equal(want.HaltedAt) {
		return ErrHaltAsk
	}
	iat, exp, now := time.Unix(*c.Iat, 0), time.Unix(*c.Exp, 0), v.now()
	if exp.Sub(iat) > haltAskLifetime || exp.Before(iat) || iat.After(now.Add(Tolerance)) || !now.Before(exp.Add(Tolerance)) {
		return ErrHaltAsk
	}
	return nil
}
