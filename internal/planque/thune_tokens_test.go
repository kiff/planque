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
	"github.com/kiff/planque/internal/spt"
)

// tokenStripe is a fake Stripe for path (a): Issuing cards (cardStripe),
// card PaymentMethods from them, and shared payment tokens, keeping what it
// was asked.
type tokenStripe struct {
	*cardStripe
	mu      sync.Mutex
	pms     []url.Values
	pmVer   []string
	issues  []url.Values
	tokens  map[string]map[string]any
	revoked []string
}

func (f *tokenStripe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/v1/payment_methods":
		b, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(b))
		f.mu.Lock()
		f.pms = append(f.pms, form)
		f.pmVer = append(f.pmVer, r.Header.Get("Stripe-Version"))
		f.mu.Unlock()
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pm_from_" + form.Get("card[issuing_card]"), "type": "card"})
	case r.Method == http.MethodPost && r.URL.Path == "/v1/shared_payment/issued_tokens":
		b, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(b))
		f.mu.Lock()
		f.issues = append(f.issues, form)
		id := "spt_test_secret_" + form.Get("shared_metadata[purchase]")
		t := map[string]any{"id": id, "status": "active",
			"usage_limits":   map[string]any{"currency": form.Get("usage_limits[currency]"), "max_amount": json.Number(form.Get("usage_limits[max_amount]")), "expires_at": json.Number(form.Get("usage_limits[expires_at]"))},
			"seller_details": map[string]any{"network_business_profile": form.Get("seller_details[network_business_profile]")}}
		f.tokens[id] = t
		f.mu.Unlock()
		_ = json.NewEncoder(w).Encode(t)
	case strings.HasPrefix(r.URL.Path, "/v1/shared_payment/issued_tokens/"):
		rest := strings.TrimPrefix(r.URL.Path, "/v1/shared_payment/issued_tokens/")
		id, revoke := strings.CutSuffix(rest, "/revoke")
		f.mu.Lock()
		t := f.tokens[id]
		if revoke && t != nil {
			f.revoked = append(f.revoked, id)
			t["status"] = "deactivated"
		}
		f.mu.Unlock()
		_ = json.NewEncoder(w).Encode(t)
	default:
		f.cardStripe.ServeHTTP(w, r)
	}
}

func withTokens(t *testing.T, e *env) *tokenStripe {
	t.Helper()
	f := &tokenStripe{cardStripe: &cardStripe{cards: map[string]map[string]any{}}, tokens: map[string]map[string]any{}}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	cl, err := issuing.New("rk_test_cards", "ich_test")
	if err != nil {
		t.Fatal(err)
	}
	cl.Base, cl.HTTP = srv.URL, srv.Client()
	e.s.Issuing = cl
	e.s.Tokens = &spt.Client{HTTP: srv.Client(), Base: srv.URL, Version: issuing.PreviewVersion,
		Key: func(context.Context) (string, error) { return "rk_test_tokens", nil }}
	return f
}

const testProfile = "profile_test_61TU90nIeGjU7NNVXA6TU90m7ISQWsBxpcx9lASWWXTk"

func tokenCaseArgs(ends time.Time, change func(map[string]any)) []byte {
	return cardArgs("approved", ends, func(c map[string]any) {
		c["merchant"] = map[string]any{"site": "shop.example", "stripe_profile": testProfile}
		if change != nil {
			change(c)
		}
	})
}

// A case's token: only under KIFF's permit for exactly that case, only for
// a token case naming the seller's profile, only once its card exists; made
// from the card's PaymentMethod (the preview), scoped to the seller, the
// card's cap, the currency and the case's end; kept sealed and never
// answered; revoked when the case ends.
func TestAPurchaseTokenIsScopedToItsCaseAndNeverAnswered(t *testing.T) {
	e := newEnv(t)
	f := withTokens(t, e)
	ctx := context.Background()
	ends := time.Now().Add(10 * time.Minute)
	c := tokenCaseArgs(ends, nil)
	mintFor := func(tool, aud string, args []byte) string {
		return e.kiff.Mint(fakes.Permit{Audience: aud, Tenant: tenant, Tool: tool, Sub: "agt_taf", Op: "op_" + tool, Args: args})
	}
	// No card yet: nothing reaches Stripe's token routes.
	if _, err := e.s.issueToken(ctx, tenant, mintFor(toolIssueToken, e.s.tokensAudience(), c), c); err == nil || !strings.Contains(err.Error(), "issue_card comes first") {
		t.Fatalf("a token before the card: %v", err)
	}
	if _, err := e.s.issueCard(ctx, tenant, mintFor(toolIssueCard, e.s.cardsAudience(), c), c); err != nil {
		t.Fatal(err)
	}
	// Permits that aren't for exactly this call, and cases a token isn't for.
	for name, tok := range map[string]string{
		"no permit":         "",
		"the card's permit": mintFor(toolIssueCard, e.s.cardsAudience(), c),
		"another audience":  mintFor(toolIssueToken, e.s.cardsAudience(), c),
	} {
		if _, err := e.s.issueToken(ctx, tenant, tok, c); err == nil {
			t.Errorf("%s: a token was issued", name)
		}
	}
	for name, args := range map[string][]byte{
		"a card fill": tokenCaseArgs(ends, func(c map[string]any) { c["path"] = "card_fill" }),
		"no profile":  tokenCaseArgs(ends, func(c map[string]any) { c["merchant"] = map[string]any{"site": "shop.example"} }),
		"a profile that isn't one": tokenCaseArgs(ends, func(c map[string]any) {
			c["merchant"] = map[string]any{"site": "shop.example", "stripe_profile": "acct_123"}
		}),
	} {
		if _, err := e.s.issueToken(ctx, tenant, mintFor(toolIssueToken, e.s.tokensAudience(), args), args); err == nil {
			t.Errorf("%s: a token was issued", name)
		}
	}
	if len(f.issues) != 0 || len(f.pms) != 0 {
		t.Fatalf("Stripe was asked %d tokens, %d payment methods", len(f.issues), len(f.pms))
	}

	ans, err := e.s.issueToken(ctx, tenant, mintFor(toolIssueToken, e.s.tokensAudience(), c), c)
	if err != nil {
		t.Fatal(err)
	}
	if ans.State != "issued" || ans.Seller != testProfile || ans.MaxAmount != 1999 || ans.Currency != "EUR" || !ans.ExpiresAt.Equal(ends.UTC().Truncate(time.Second)) {
		t.Fatalf("answer: %+v", ans)
	}
	out, _ := json.Marshal(ans)
	if strings.Contains(string(out), "spt_") {
		t.Fatalf("the token's id was answered: %s", out)
	}
	// From the card's PaymentMethod, made with the preview version; scoped
	// exactly to the case.
	if len(f.pms) != 1 || !strings.HasPrefix(f.pms[0].Get("card[issuing_card]"), "ic_") || f.pms[0].Get("type") != "card" || f.pmVer[0] != issuing.PreviewVersion {
		t.Fatalf("payment method: %v %v", f.pms, f.pmVer)
	}
	is := f.issues[0]
	if is.Get("payment_method") != "pm_from_"+f.pms[0].Get("card[issuing_card]") || is.Get("seller_details[network_business_profile]") != testProfile ||
		is.Get("usage_limits[max_amount]") != "1999" || is.Get("usage_limits[currency]") != "eur" || is.Get("usage_limits[expires_at]") != itoa(ends.Unix()) {
		t.Fatalf("token request: %v", is)
	}
	// Kept sealed: the record never holds the id in the clear.
	items, _ := e.kv.Query(ctx, pkIssuing, "TOKEN#")
	if len(items) != 1 || strings.Contains(string(items[0].Doc), "spt_") {
		t.Fatalf("the token's record: %d items", len(items))
	}
	// While the case runs, the sweep leaves it; once it ends, it's revoked.
	if n, err := e.s.sweepTokens(ctx); err != nil || n != 0 || len(f.revoked) != 0 {
		t.Fatalf("sweep while running: %d %v %v", n, err, f.revoked)
	}
	later := ends.Add(time.Minute)
	e.s.Now = func() time.Time { return later }
	n, err := e.s.sweepTokens(ctx)
	e.s.Now = nil
	if err != nil || n != 1 || len(f.revoked) != 1 || !strings.HasPrefix(f.revoked[0], "spt_test_secret_") {
		t.Fatalf("sweep after the case: %d %v %v", n, err, f.revoked)
	}
	if items, _ := e.kv.Query(ctx, pkIssuing, "TOKEN#"); len(items) != 0 {
		t.Fatal("the token's record was kept after its revoke")
	}
}

func itoa(n int64) string { b, _ := json.Marshal(n); return string(b) }

// Off unless both connections are configured.
func TestNoTokensWithoutTheirConnection(t *testing.T) {
	e := newEnv(t)
	withCards(t, e)
	c := tokenCaseArgs(time.Now().Add(10*time.Minute), nil)
	if _, err := e.s.issueToken(context.Background(), tenant, "x", c); err == nil || !strings.Contains(err.Error(), "off") {
		t.Fatalf("no token connection: %v", err)
	}
}
