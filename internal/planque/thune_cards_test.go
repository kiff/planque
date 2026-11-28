package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/issuing"
	"github.com/kiff/planque/internal/jcs"
)

// cardStripe is a fake Stripe Issuing: it makes, lists and cancels cards and
// counts what it was asked.
type cardStripe struct {
	mu      sync.Mutex
	creates int
	cards   map[string]map[string]any
	n       int
}

func (f *cardStripe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	form, _ := url.ParseQuery(string(b))
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/v1/issuing/cards":
		f.creates++
		f.n++
		id := "ic_" + strings.Repeat("x", f.n)
		f.cards[id] = map[string]any{"id": id, "status": "active", "last4": "4242", "currency": form.Get("currency"),
			"metadata": map[string]string{"case": form.Get("metadata[case]"), "case_version": form.Get("metadata[case_version]"), "permit": form.Get("metadata[permit]")}}
		_ = json.NewEncoder(w).Encode(f.cards[id])
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/v1/issuing/cards/"):
		c := f.cards[strings.TrimPrefix(r.URL.Path, "/v1/issuing/cards/")]
		c["status"] = form.Get("status")
		_ = json.NewEncoder(w).Encode(c)
	case r.Method == http.MethodGet:
		var data []map[string]any
		for _, c := range f.cards {
			if c["status"] == "active" {
				data = append(data, c)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
	}
}

func (f *cardStripe) status(id string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c := f.cards[id]; c != nil {
		return c["status"].(string)
	}
	return ""
}

func withCards(t *testing.T, e *env) *cardStripe {
	t.Helper()
	f := &cardStripe{cards: map[string]map[string]any{}}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	cl, err := issuing.New("rk_test_cards", "ich_test")
	if err != nil {
		t.Fatal(err)
	}
	cl.Base, cl.HTTP = srv.URL, srv.Client()
	e.s.Issuing = cl
	return f
}

// cardArgs is issue_card's arguments for a kiff.purchase.v1 case, as
// contracts/purchase.json writes it, with change applied to the case first.
func cardArgs(state string, ends time.Time, change func(map[string]any)) []byte {
	c := map[string]any{"schema": "kiff.purchase.v1", "case_id": "case_7", "version": 1, "path": "token", "currency": "EUR",
		"price": map[string]any{"total": 1899, "max": 1999},
		"grant": map[string]any{"kind": "single_use_card", "cap": 1999, "categories": []string{"book_stores"}, "uses": 1}}
	if change != nil {
		change(c)
	}
	raw, _ := json.Marshal(c)
	sum, _ := jcs.SHA256(raw)
	b, _ := json.Marshal(map[string]any{"case": json.RawMessage(raw), "digest": sum, "state": state, "valid_until": ends.UTC().Format(time.RFC3339)})
	return b
}

func purchaseCaseArgs(state string, ends time.Time) []byte { return cardArgs(state, ends, nil) }

// A card is made only under KIFF's permit for exactly that case: the right
// audience, account, tool and arguments, still live. Every other permit
// sends nothing to Stripe, and a made card is recorded with its case and
// permit.
func TestACardIsMadeOnlyUnderKIFFsPermitForThatCase(t *testing.T) {
	e := newEnv(t)
	f := withCards(t, e)
	ctx := context.Background()
	c := purchaseCaseArgs("approved", time.Now().Add(10*time.Minute))
	mint := func(change func(*fakes.Permit)) string {
		p := fakes.Permit{Audience: e.s.cardsAudience(), Tenant: tenant, Tool: toolIssueCard, Sub: "agt_taf", Op: "op_card", Args: c}
		if change != nil {
			change(&p)
		}
		return e.kiff.Mint(p)
	}
	for name, token := range map[string]string{
		"no permit":        "",
		"another account":  mint(func(p *fakes.Permit) { p.Tenant = "t_other" }),
		"another tool":     mint(func(p *fakes.Permit) { p.Tool = "browse_pay" }),
		"another case":     mint(func(p *fakes.Permit) { p.Args = purchaseCaseArgs("approved", time.Now().Add(5*time.Minute)) }),
		"another audience": mint(func(p *fakes.Permit) { p.Audience = e.s.BaseURL + "/mcp/x" }),
		"expired":          mint(func(p *fakes.Permit) { p.IssuedAt = time.Now().Add(-5 * time.Minute) }),
	} {
		if _, err := e.s.issueCard(ctx, tenant, token, c); err == nil {
			t.Errorf("%s: a card was made", name)
		}
	}
	// A permit for a case that isn't approved: verified, and still no card.
	pending := purchaseCaseArgs("pending", time.Now().Add(10*time.Minute))
	if _, err := e.s.issueCard(ctx, tenant, e.kiff.Mint(fakes.Permit{Audience: e.s.cardsAudience(), Tenant: tenant, Tool: toolIssueCard, Sub: "agt_taf", Op: "op_p", Args: pending}), pending); err == nil {
		t.Error("a card for a pending case")
	}
	// A permit for arguments that don't hold: a digest that isn't the
	// case's, a merchant code where Stripe's category name goes, a cap above
	// the case's maximum, no single-use card. Each is signed, and still no card.
	ends := time.Now().Add(10 * time.Minute)
	badDigest := func() []byte {
		var m map[string]any
		_ = json.Unmarshal(cardArgs("approved", ends, nil), &m)
		m["digest"] = strings.Repeat("0", 64)
		b, _ := json.Marshal(m)
		return b
	}()
	for name, args := range map[string][]byte{
		"a digest that isn't the case's": badDigest,
		"a merchant code":                cardArgs("approved", ends, func(c map[string]any) { c["grant"].(map[string]any)["categories"] = []string{"5943"} }),
		"a cap above the maximum":        cardArgs("approved", ends, func(c map[string]any) { c["grant"].(map[string]any)["cap"] = 2500 }),
		"two uses":                       cardArgs("approved", ends, func(c map[string]any) { c["grant"].(map[string]any)["uses"] = 2 }),
		"no payment path":                cardArgs("approved", ends, func(c map[string]any) { c["path"] = "none" }),
		"another schema":                 cardArgs("approved", ends, func(c map[string]any) { c["schema"] = "kiff.purchase.v0" }),
	} {
		tok := e.kiff.Mint(fakes.Permit{Audience: e.s.cardsAudience(), Tenant: tenant, Tool: toolIssueCard, Sub: "agt_taf", Op: "op_x", Args: args})
		if _, err := e.s.issueCard(ctx, tenant, tok, args); err == nil {
			t.Errorf("%s: a card was made", name)
		}
	}
	if f.creates != 0 {
		t.Fatalf("%d cards were asked of Stripe", f.creates)
	}

	card, err := e.s.issueCard(ctx, tenant, mint(nil), c)
	if err != nil {
		t.Fatal(err)
	}
	var rec issuedCard
	if _, err := getDoc(ctx, e.kv, pkIssuing, skIssuedCard(card.ID), &rec); err != nil {
		t.Fatalf("no record: %v", err)
	}
	if rec.Tenant != tenant || rec.Case != "case_7" || rec.Version != 1 || rec.Permit == "" || rec.Permit != card.Permit {
		t.Fatalf("record: %+v, card %+v", rec, card)
	}
}

// The sweep cancels a card whose case has ended and any active card planque
// has no record of; a card whose case still runs stays.
func TestTheSweepCancelsCardsPastTheirCaseAndCardsWithNoRecord(t *testing.T) {
	e := newEnv(t)
	f := withCards(t, e)
	ctx := context.Background()
	c := purchaseCaseArgs("approved", time.Now().Add(10*time.Minute))
	running, err := e.s.issueCard(ctx, tenant, e.kiff.Mint(fakes.Permit{Audience: e.s.cardsAudience(), Tenant: tenant, Tool: toolIssueCard, Sub: "agt_taf", Op: "op_1", Args: c}), c)
	if err != nil {
		t.Fatal(err)
	}
	// One made outside planque (by hand, or with a stolen key).
	f.mu.Lock()
	f.cards["ic_stray"] = map[string]any{"id": "ic_stray", "status": "active", "last4": "0005", "metadata": map[string]string{}}
	f.mu.Unlock()

	n, err := e.s.sweepCards(ctx)
	if err != nil || n != 1 {
		t.Fatalf("swept %d, %v", n, err)
	}
	if f.status("ic_stray") != "canceled" || f.status(running.ID) != "active" {
		t.Fatalf("stray %s, running %s", f.status("ic_stray"), f.status(running.ID))
	}
	// After its case ends, the running card goes too.
	e.s.Now = func() time.Time { return time.Now().Add(11 * time.Minute) }
	if n, err := e.s.sweepCards(ctx); err != nil || n != 1 || f.status(running.ID) != "canceled" {
		t.Fatalf("after the case: swept %d, %v, status %s", n, err, f.status(running.ID))
	}
}

// With no Issuing connection, no card and no sweep.
func TestNoCardsWithoutAnIssuingConnection(t *testing.T) {
	e := newEnv(t)
	if _, err := e.s.issueCard(context.Background(), tenant, "x", nil); err != errCardsOff {
		t.Fatalf("got %v", err)
	}
	if n, err := e.s.sweepCards(context.Background()); n != 0 || err != nil {
		t.Fatalf("swept %d, %v", n, err)
	}
}
