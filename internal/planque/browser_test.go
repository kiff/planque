package planque

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/store"
)

// browserEnv is an in-process planque with the fake browser worker, the
// test shop at shop.example, and a KIFF that applies profiles.
type browserEnv struct {
	*env
	b  *fakes.Browser
	ev *fakes.Evidence
}

func newBrowserEnv(t *testing.T) *browserEnv {
	t.Helper()
	e := newInProcessEnv(t)
	e.enrol()
	base := e.client.Transport
	e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
			return base.RoundTrip(r)
		}
		var in struct {
			Profile json.RawMessage `json:"profile"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			t.Fatal(err)
		}
		hash, _ := profile.Digest(in.Profile)
		return profileResponse(200, `{"revision":1,"profile_sha256":"`+hash+`","words_sha256":"`+strings.Repeat("a", 64)+`"}`), nil
	})
	be := &browserEnv{env: e, b: fakes.NewBrowser(), ev: fakes.NewEvidence()}
	e.s.Browser, e.s.Evidence = be.b, be.ev
	// These tests pay; paying off is tested on its own (browser_pay_off_test.go).
	e.s.BrowserPay = true
	return be
}

// ask asks for a browser grant as Setup does.
func (e *browserEnv) ask(body map[string]any) (int, map[string]any) {
	e.t.Helper()
	req := map[string]any{"kind": KindBrowser, "agent": "agt_taf", "tools": []any{}}
	for k, v := range body {
		req[k] = v
	}
	return e.do("POST", "/v1/accounts/"+tenant+"/grants", req, nil, "k_setup")
}

// connect asks for a grant limited to shop.example paying in EUR, and
// approves it on the phone: connected at once, no sign-in.
func (e *browserEnv) connect() *Grant {
	e.t.Helper()
	st, out := e.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}})
	if st != 201 || out["status"] != StatusPending {
		e.t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	hash, _ := profile.Digest(g.Request.Profile)
	approved := ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": approved}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	g, err = e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		e.t.Fatal(err)
	}
	return g
}

// call is one tools/call at the relay with a permit KIFF minted for it.
func (e *browserEnv) call(g *Grant, tool string, args map[string]any) (string, map[string]any) {
	e.t.Helper()
	op := "op_" + strings.ToLower(tool) + "_" + randomSuffix()
	args[IdempotencyArg] = op
	raw := mustJSON(args)
	mint := fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: tool, Sub: g.Request.Agent, Op: op, Args: raw}
	if tool == browsePay {
		// KIFF holds every pay for its owner, so its permit names the approval.
		mint.Approval = "exc_" + op
	}
	return textOf(e.s.call(context.Background(), g, tool, raw, e.kiff.Mint(mint)))
}

// textOf is a tool result as the agent reads it (its text) and its
// structured content, both as JSON gives them.
func textOf(res any) (string, map[string]any) {
	var m struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Structured map[string]any `json:"structuredContent"`
	}
	_ = json.Unmarshal(mustJSON(res), &m)
	text := ""
	if len(m.Content) > 0 {
		text = m.Content[0].Text
	}
	if m.Structured["fenced"] == true {
		// The owner's Chrome: the answer is only inside the fence.
		return text, fencedJSON(text)
	}
	return text, m.Structured
}

var suffixes int

func randomSuffix() string {
	suffixes++
	return strings.Repeat("x", suffixes%5+1) + string(rune('a'+suffixes%26))
}

// fencedJSON is the answer inside a fence, as the taf reads it.
func fencedJSON(text string) map[string]any {
	start := strings.Index(text, ">>>\n")
	end := strings.LastIndex(text, "\n<<</page-words:")
	if start < 0 || end < start {
		return nil
	}
	var v map[string]any
	_ = json.Unmarshal([]byte(text[start+4:end]), &v)
	return v
}

func code(sc map[string]any) string {
	c, _ := sc["code"].(string)
	return c
}

// 1. Asking for the browser: the owner's words name the sites and what
// the taf may do; the profile is planque's, with browse_pay's amount read
// from browse_checkout; a bad ask is refused in words.
func TestBrowserGrantIsAskedForWithPlanquesOwnToolsAndProfile(t *testing.T) {
	e := newBrowserEnv(t)
	st, out := e.ask(map[string]any{"browser": map[string]any{"sites": []string{"Shop.Example", "docs.example.org"}, "currency": "EUR"}})
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	text := out["text"].(string)
	for _, want := range []string{"shop.example, docs.example.org", "Pay a checkout: pays a checkout in EUR only", "never a payment, a password, a card number or a code", "you sign in yourself"} {
		if !strings.Contains(text, want) {
			t.Errorf("the owner's text lacks %q:\n%s", want, text)
		}
	}
	if out["provider"] != "browser" || out["entry"] != "browser" || out["vendor"] != "shop.example, docs.example.org" {
		t.Errorf("view: provider %v entry %v vendor %v", out["provider"], out["entry"], out["vendor"])
	}
	req := out["request"].(map[string]any)
	b := req["browser"].(map[string]any)
	if sites, _ := b["sites"].([]any); len(sites) != 2 || sites[0] != "shop.example" {
		t.Errorf("sites: %v", b)
	}
	var p profile.Profile
	raw, _ := json.Marshal(req["profile"])
	if err := json.Unmarshal(raw, &p); err != nil || len(p.Operations) != 5 || p.Source != "catalog" || p.CatalogEntry != "browser@1" {
		t.Fatalf("profile: %v %+v", err, p)
	}
	for _, op := range p.Operations {
		switch op.Tool {
		case "browse_pay":
			// Every pay waits for its owner, whatever a Card's room.
			if op.Effect != profile.EffectMovesMoney || !op.AskEachCall || op.Amount == nil || op.Amount.From != "read" || op.Amount.Read.Tool != "browse_checkout" || op.Amount.Unit != "EUR" || op.Amount.Read.Args["url"] != "$args.url" {
				t.Errorf("browse_pay: %+v", op)
			}
		case "browse_act":
			if op.Effect != profile.EffectWrites || op.AskEachCall || op.Resource == nil || op.Resource.Arg != "url" {
				t.Errorf("browse_act: %+v", op)
			}
		default:
			if op.Kind != "read" {
				t.Errorf("%s isn't a read", op.Tool)
			}
		}
	}
	// The catalog lists the browser, and an ask through it is the same.
	st, cat := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
	if st != 200 || !strings.Contains(mustJSONString(cat), `"entry":"browser"`) {
		t.Errorf("catalog: %d %v", st, cat)
	}
	st, viaCatalog := e.do("POST", "/v1/accounts/"+tenant+"/catalog/browser/grants", map[string]any{"agent": "agt_taf", "browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}}, nil, "k_setup")
	if st != 201 || viaCatalog["needs_key"] != false || viaCatalog["grant"].(map[string]any)["request"].(map[string]any)["profile"].(map[string]any)["vendor"] != "shop.example" {
		t.Errorf("catalog ask: %d %v", st, viaCatalog)
	}

	for _, c := range []struct {
		name string
		body map[string]any
		want string
	}{
		{"a site with a scheme", map[string]any{"browser": map[string]any{"sites": []string{"https://shop.example"}}}, "isn't a site"},
		{"a local site", map[string]any{"browser": map[string]any{"sites": []string{"db.internal"}}}, "public site"},
		{"pay with no currency", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}}}, "currency is required"},
		{"no sites", map[string]any{"browser": map[string]any{"currency": "EUR"}}, "sites is required"},
		{"a shared host", map[string]any{"browser": map[string]any{"sites": []string{"shop.example", "github.io"}}}, "shared domain"},
		{"no browser at all", map[string]any{}, "sites is required"},
		{"pay with no checkout", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}, "tools": []map[string]any{{"tool": "browse_pay", "title": "Browse pay"}}}, "needs browse_checkout"},
		{"a tool it hasn't", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}, "tools": []map[string]any{{"tool": "browse_download", "title": "Browse download"}}}, "no tool browse_download"},
		{"a profile of the agent's", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}, "profile": map[string]any{"v": "x"}}, "profile is planque's"},
		{"a server", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"}, "server": "https://mcp.example"}, "profile is planque's"},
		{"a bad currency", map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "euros"}}, "ISO 4217"},
	} {
		st, out := e.ask(c.body)
		if st != 400 || !strings.Contains(out["error"].(string), c.want) {
			t.Errorf("%s: %d %v, want %q", c.name, st, out, c.want)
		}
	}
	// Reads alone need no currency.
	if st, out := e.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}}, "tools": []map[string]any{{"tool": "browse_open", "title": "Browse open", "read_only": true}, {"tool": "browse_read", "title": "Browse read", "read_only": true}}}); st != 201 {
		t.Errorf("reads alone: %d %v", st, out)
	}
}

// 2. Approving connects at once, with the tools in KIFF as the profile
// says: browse_pay's amount is looked up with browse_checkout, in EUR.
func TestBrowserGrantConnectsOnApprovalWithTheProfileInKIFF(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if len(g.Tools) != 5 || g.Transport == "" || g.ProfileRevision == nil {
		t.Fatalf("connected grant: %+v", g)
	}
	pay, ok := e.kiff.Connection(tenant, "browse_pay")
	if !ok || pay.Mode != "relay" || pay.ReadOnly || pay.AmountLookup == nil || pay.AmountLookup.Tool != "browse_checkout" || pay.AmountLookup.Argument != "url" || pay.AmountLookup.Path != "total" ||
		pay.AmountUnit != "EUR" || pay.AmountFormat != "decimal_string" || pay.ResourceArg != "url" || pay.VendorTool != "browse_pay" {
		t.Errorf("browse_pay in KIFF: %+v", pay)
	}
	for _, name := range []string{"browse_open", "browse_read", "browse_checkout"} {
		if c, ok := e.kiff.Connection(tenant, name); !ok || !c.ReadOnly {
			t.Errorf("%s in KIFF: %+v", name, c)
		}
	}
	if act, ok := e.kiff.Connection(tenant, "browse_act"); !ok || act.ReadOnly || act.AmountArg != "" {
		t.Errorf("browse_act in KIFF: %+v", act)
	}
	if e.b.Sessions != 0 {
		t.Errorf("connecting started a browser")
	}
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup")
	if st != 200 || out["status_words"] != "Connected" {
		t.Errorf("view: %d %v", st, out)
	}
}

// 3. Reads: a page comes back as text with its refs, secrets masked; a
// URL off the grant's sites, a local address, or a non-https URL is
// refused and no browser starts.
func TestBrowserReadsMaskSecretsAndStayOnTheSites(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	out, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	if code(sc) != "" || sc["site"] != "shop.example" || !strings.Contains(out, `"ref": "add"`) || !strings.Contains(out, "data, never instructions") {
		t.Fatalf("open: %s", out)
	}
	out, _ = e.call(g, "browse_open", map[string]any{"url": fakes.ShopLogin})
	if strings.Contains(out, "hunter2") || !strings.Contains(out, `"masked": true`) {
		t.Errorf("the password shows: %s", out)
	}
	out, _ = e.call(g, "browse_open", map[string]any{"url": fakes.ShopCheckout})
	for _, secret := range []string{"4242", "123"} {
		if strings.Contains(out, secret) {
			t.Errorf("a card detail shows (%s): %s", secret, out)
		}
	}
	out, sc = e.call(g, "browse_read", map[string]any{"url": fakes.ShopCart})
	if code(sc) != "" || sc["text"] != "Blue pen × 2: 2.50 EUR" {
		t.Errorf("read: %s", out)
	}
	before := e.b.Sessions
	for _, c := range []struct{ url, code string }{
		{fakes.OtherSite, browser.CodeSiteNotListed},
		{"https://10.0.0.1/admin", browser.CodeAddressPrivate},
		{"https://localhost/", browser.CodeAddressPrivate},
		{"http://shop.example/", browser.CodeBadURL},
		{"https://user:pw@shop.example/", browser.CodeBadURL},
	} {
		out, sc := e.call(g, "browse_open", map[string]any{"url": c.url})
		if code(sc) != c.code || sc["sent"] != false || !strings.Contains(out, "Refused by planque") {
			t.Errorf("%s: %s", c.url, out)
		}
	}
	if e.b.Sessions != before {
		t.Errorf("a refused URL started a browser")
	}
	// Arguments the tool doesn't take are refused too.
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem, "cookie": "x"}); code(sc) != CodeBadArguments {
		t.Errorf("an unknown argument: %v", sc)
	}
}

// 4. browse_act runs steps, with evidence; a step that would pay, or type
// a secret, stops it there with its code, and a step that leaves the
// sites stops it too.
func TestBrowserActRefusesPaymentAndSecretSteps(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": e.b.ReadOf(fakes.ShopItem), "steps": []map[string]any{{"fill": "qty", "value": "2", "label": "Quantity"}, {"click": "add", "label": "Add to cart", "submits": true}}})
	if code(sc) != "" || sc["steps_run"] != float64(2) || sc["site"] != "shop.example" {
		t.Fatalf("act: %s", out)
	}
	ev := sc["evidence"].(map[string]any)
	if ev["before"].(map[string]any)["url"] != fakes.ShopItem || ev["after"].(map[string]any)["url"] != fakes.ShopCart || ev["screenshot"] == "" {
		t.Errorf("evidence: %v", ev)
	}
	if refs := e.ev.Refs(tenant, g.ID); len(refs) != 1 || string(e.ev.Kept[refs[0]]) != "\x89PNG "+fakes.ShopCart {
		t.Errorf("screenshot kept: %v", refs)
	}

	sent := e.b.Sent()
	for _, c := range []struct {
		name  string
		url   string
		steps []map[string]any
		code  string
		ran   int
	}{
		{"placing the order after a fill", fakes.ShopCheckout, []map[string]any{{"fill": "email", "value": "a@example.com", "label": "Email"}, {"click": "place"}}, browser.CodePaymentStep, 1},
		{"a submit on the payment page", fakes.ShopCheckout, []map[string]any{{"click": "apply"}}, browser.CodePaymentStep, 0},
		{"typing a card number", fakes.ShopCheckout, []map[string]any{{"fill": "card", "value": "4000 0000 0000 0002"}}, browser.CodePaymentStep, 0},
		{"typing a security code", fakes.ShopCheckout, []map[string]any{{"fill": "cvc", "value": "999"}}, browser.CodeSecretField, 0},
		{"typing a password", fakes.ShopLogin, []map[string]any{{"fill": "pass", "value": "hunter2"}}, browser.CodeSecretField, 0},
		{"a ref that isn't there", fakes.ShopItem, []map[string]any{{"click": "buy_now"}}, browser.CodeRefUnknown, 0},
		{"Commander on a French checkout", fakes.ShopCheckoutFR, []map[string]any{{"click": "cmd"}}, browser.CodePaymentStep, 0},
		{"Jetzt kaufen on a checkout", fakes.ShopCheckoutFR, []map[string]any{{"click": "zb"}}, browser.CodePaymentStep, 0},
		{"an unlabelled button on a checkout", fakes.ShopCheckoutFR, []map[string]any{{"fill": "email", "value": "a@example.com", "label": "Courriel"}, {"click": "icon"}}, browser.CodePaymentStep, 1},
		{"a one-click Kaufen on a product page", fakes.ShopItemDE, []map[string]any{{"click": "kaufen"}}, browser.CodePaymentStep, 0},
		{"a redirect off the sites", fakes.Away, []map[string]any{{"click": "place"}}, browser.CodeSiteNotListed, 0},
		{"leaving the sites", fakes.Leaving, []map[string]any{{"click": "other", "label": "Our partner", "to": "other.example"}}, browser.CodeSiteNotListed, 1},
		{"a malformed step", fakes.ShopItem, []map[string]any{{"fill": "qty", "click": "add"}}, CodeBadArguments, 0},
	} {
		out, sc := e.call(g, "browse_act", map[string]any{"url": c.url, "read": e.b.ReadOf(c.url), "steps": c.steps})
		if code(sc) != c.code || !strings.Contains(out, "Refused by planque") {
			t.Errorf("%s: %s", c.name, out)
		}
		ran, _ := sc["steps_run"].(float64)
		if int(ran) != c.ran || sc["sent"] != (c.ran > 0) || e.b.Sent() != sent+c.ran {
			t.Errorf("%s: steps run %v, sent %v, worker saw %d (was %d)", c.name, sc["steps_run"], sc["sent"], e.b.Sent(), sent)
		}
		sent = e.b.Sent()
	}
	if len(e.b.Orders) != 0 {
		t.Fatalf("an order was placed: %v", e.b.Orders)
	}
	if strings.Contains(e.b.StepsJSON(), "hunter2") || strings.Contains(e.b.StepsJSON(), "4000") {
		t.Errorf("a secret was typed: %s", e.b.StepsJSON())
	}
}

// A listed URL that redirects off the sites: every tool refuses the page
// it landed on, nothing is read or recorded, and a pay clicks nothing.
func TestBrowserRefusesARedirectOffTheSites(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	for _, tool := range []string{"browse_open", "browse_read", "browse_checkout"} {
		out, sc := e.call(g, tool, map[string]any{"url": fakes.Away})
		if code(sc) != browser.CodeSiteNotListed || sc["sent"] != false || strings.Contains(out, "Another site") || sc["total"] != nil {
			t.Errorf("%s after a redirect: %s", tool, out)
		}
	}
	if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skCheckout(g.ID, fakes.Away)); ok {
		t.Error("the other site's checkout was recorded")
	}
	// Even with a checkout record under that URL (read before the site
	// started redirecting), a pay refuses the page it lands on.
	rec := checkoutRecord{Site: "shop.example", Checkout: browser.Checkout{Lines: []browser.Line{{Item: "Blue pen", Quantity: 2, Price: "1.25"}}, Total: "2.50", Currency: "EUR"}, ReadAt: e.s.now()}
	if err := putDoc(context.Background(), e.kv, store.Item{PK: tenantPK(tenant), SK: skCheckout(g.ID, fakes.Away), Ver: 1}, rec, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.Away, "ref": "place"})
	if code(sc) != browser.CodeSiteNotListed || len(e.b.Orders) != 0 || e.b.Sent() != 0 {
		t.Errorf("pay after a redirect: %s (orders %v, steps %d)", out, e.b.Orders, e.b.Sent())
	}
}

// 5. Paying: the checkout is read first (the total KIFF counts), then paid
// once with the permit; a changed checkout, another currency, an unread
// total, or the wrong control sends nothing; a replay answers as the first
// did.
func TestBrowserPaysOnceWhatWasRead(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	// Pay before any read: nothing starts.
	out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
	if code(sc) != browser.CodeTotalUnread || e.b.Sessions != 0 || !strings.Contains(out, "browse_checkout first") {
		t.Fatalf("pay before a read: %s", out)
	}
	out, sc = e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout})
	if code(sc) != "" || sc["total"] != "2.50" || sc["currency"] != "EUR" || len(sc["lines"].([]any)) != 1 {
		t.Fatalf("checkout: %s", out)
	}
	// The shop changes its total before the owner approves.
	e.b.Total = "3.00"
	out, sc = e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
	if code(sc) != browser.CodeCheckoutChanged || sc["sent"] != false || len(e.b.Orders) != 0 || !strings.Contains(out, "2.50 EUR") || !strings.Contains(out, "3.00 EUR") {
		t.Fatalf("pay after a change: %s", out)
	}
	// Read again, then another currency on the page.
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "3.00" {
		t.Fatalf("second read: %v", sc)
	}
	e.b.Currency = "USD"
	if out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeCurrencyOther || len(e.b.Orders) != 0 {
		t.Fatalf("pay in USD: %s", out)
	}
	// A checkout in another currency gives no total: KIFF reads none.
	if out, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); code(sc) != "" || sc["total"] != nil || !strings.Contains(out, "pays in EUR only") {
		t.Fatalf("checkout in USD: %s", out)
	}
	e.b.Currency = "EUR"
	// A page that isn't a checkout gives no total either.
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopItem}); sc["checkout"] != nil || sc["total"] != nil {
		t.Fatalf("an item page as a checkout: %v", sc)
	}
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "3.00" {
		t.Fatalf("third read: %v", sc)
	}
	// The wrong control: nothing placed.
	if out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "apply"}); code(sc) != browser.CodeRefUnknown || len(e.b.Orders) != 0 {
		t.Fatalf("pay with the coupon button: %s", out)
	}
	// Paid once, with the permit; the same permit again sends nothing.
	op := "op_pay_once"
	args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: op})
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args, Approval: "exc_" + op})
	first, _ := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
	if strings.Contains(first, "Refused") || len(e.b.Orders) != 1 || !strings.Contains(first, `"total": "3.00"`) || !strings.Contains(first, "Order 1001 placed") {
		t.Fatalf("pay: %s (orders %v)", first, e.b.Orders)
	}
	if again, _ := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit)); again != first || len(e.b.Orders) != 1 {
		t.Fatalf("replay: %s (orders %v)", again, e.b.Orders)
	}
	// The read was used: paying again on it needs a new read.
	if _, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeTotalUnread || len(e.b.Orders) != 1 {
		t.Fatalf("pay again without a read: %v", sc)
	}
	if refs := e.ev.Refs(tenant, g.ID); len(refs) != 1 {
		t.Errorf("the pay's screenshot: %v", refs)
	}
}

// A pay click that ended unknown (the order was placed, the answer lost) marks its checkout: any later pay on the
// same lines, total and currency is refused, whatever operation asks, and a new read doesn't clear it. A different
// cart is a new checkout. Once per operation was the only guard before: a new operation re-read the page and
// clicked again within 20 minutes.
// Category: money
// Attack: an agent's retry under a new operation after a click with a lost answer places a second order
// Refused: the second pay sends nothing (outcome_unknown), the same checkout read again gets no total, one order exists
func TestAnUnknownPayBlocksTheSameCheckoutUnderAnyOperation(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
		t.Fatalf("read: %v", sc)
	}
	var last map[string]any
	pay := func(op string) string {
		args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: op})
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args, Approval: "exc_" + op})
		out, sc := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
		last = sc
		return out
	}
	e.b.LoseAnswer = true
	first := pay("op_lost")
	if !strings.Contains(first, "isn't known") || len(e.b.Orders) != 1 {
		t.Fatalf("the click with a lost answer: %s (orders %v)", first, e.b.Orders)
	}
	e.b.LoseAnswer = false
	clicks := e.b.Sent()
	// A new operation on the same checkout sends nothing, however it asks.
	second := pay("op_retry")
	if code(last) != browser.CodeOutcomeUnknown || !strings.Contains(second, "without a known outcome") || len(e.b.Orders) != 1 || e.b.Sent() != clicks {
		t.Fatalf("a new operation: %s (orders %v, steps %d→%d)", second, e.b.Orders, clicks, e.b.Sent())
	}
	// Reading it again doesn't make it payable: no total, and the pay is still refused.
	if out, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != nil || !strings.Contains(out, "without a known outcome") {
		t.Fatalf("a read of the same checkout: %s", out)
	}
	if third := pay("op_third"); code(last) != browser.CodeOutcomeUnknown || len(e.b.Orders) != 1 {
		t.Fatalf("after a read: %s (orders %v)", third, e.b.Orders)
	}
	// A different cart (another total) is a new checkout: read, then paid once.
	e.b.Total = "9.00"
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "9.00" {
		t.Fatalf("a different cart's read: %v", sc)
	}
	if out := pay("op_other_cart"); strings.Contains(out, "Refused") || len(e.b.Orders) != 2 {
		t.Fatalf("a different cart: %s (orders %v)", out, e.b.Orders)
	}
}

// The mark is written before the click, not after its error: a pay at the same moment, and a click that never
// returns (the process dies), both find the checkout marked and send nothing.
// Category: money
// Attack: two pays on one checkout at once, or a pay whose process died mid-click and an agent's retry, place two orders
// Refused: exactly one click; the other pay is outcome_unknown with nothing sent
func TestThePayIsMarkedBeforeItsClick(t *testing.T) {
	read := func(e *browserEnv, g *Grant) {
		t.Helper()
		if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
			t.Fatalf("read: %v", sc)
		}
	}
	pay := func(e *browserEnv, g *Grant, op string) (string, map[string]any) {
		args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: op})
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args, Approval: "exc_" + op})
		return textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
	}
	t.Run("two pays at once", func(t *testing.T) {
		e := newBrowserEnv(t)
		g := e.connect()
		read(e, g)
		var other string
		var otherSC map[string]any
		e.b.BeforeClick = func() {
			e.b.BeforeClick = nil
			// The first pay is at its click; a second operation pays the same checkout now. The
			// account's browser lease (browser_lease.go) would make it wait; here the first's
			// lease has lapsed, so what's tested is the checkout's own mark.
			if _, _, err := e.s.KV.Delete(context.Background(), tenantPK(g.Tenant), skBrowserLease(0)); err != nil {
				t.Fatal(err)
			}
			other, otherSC = pay(e, g, "op_second")
		}
		first, _ := pay(e, g, "op_first")
		if strings.Contains(first, "Refused") || len(e.b.Orders) != 1 {
			t.Fatalf("the first pay: %s (orders %v)", first, e.b.Orders)
		}
		if code(otherSC) != browser.CodeOutcomeUnknown || !strings.Contains(other, "Nothing was sent") {
			t.Fatalf("the second pay: %s", other)
		}
	})
	t.Run("the process dies mid-click", func(t *testing.T) {
		e := newBrowserEnv(t)
		g := e.connect()
		read(e, g)
		e.b.BeforeClick = func() { e.b.BeforeClick = nil; panic("the process died") }
		func() {
			defer func() { _ = recover() }()
			pay(e, g, "op_died")
		}()
		if len(e.b.Orders) != 0 {
			t.Fatalf("the fake placed an order: %v", e.b.Orders)
		}
		// The agent's retry, a new operation, finds the checkout marked.
		if out, sc := pay(e, g, "op_retry"); code(sc) != browser.CodeOutcomeUnknown || len(e.b.Orders) != 0 {
			t.Fatalf("the retry after a death: %s (orders %v)", out, e.b.Orders)
		}
	})
}

// Two pays that both read the checkout's record before either marked it: only the first arms it. The second's write
// is conditional on the version it read, so it loses and clicks nothing (the table decides, not the order of reads).
func TestOnlyOnePayArmsACheckoutItBothRead(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
		t.Fatalf("read: %v", sc)
	}
	var rec checkoutRecord
	it, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skCheckout(g.ID, fakes.ShopCheckout), &rec)
	if err != nil {
		t.Fatal(err)
	}
	// Both pays hold the same record at the same version.
	if _, ok := e.s.armPay(context.Background(), g, fakes.ShopCheckout, rec, it.Ver); !ok {
		t.Fatal("the first pay couldn't arm the checkout")
	}
	refusal, ok := e.s.armPay(context.Background(), g, fakes.ShopCheckout, rec, it.Ver)
	if ok {
		t.Fatal("the second pay armed a checkout the first already armed")
	}
	if _, sc := textOf(refusal); code(sc) != browser.CodeOutcomeUnknown {
		t.Fatalf("the second pay's refusal: %v", refusal)
	}
}

// The mark lasts unknownKept, not checkoutKept: 21 minutes and 23 hours after the unknown click the same checkout
// is still refused; past 24 hours the record is gone, so a pay is refused for its own reason (no total was read).
func TestAnUnknownMarkLastsADayNotTwentyMinutes(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
		t.Fatalf("read: %v", sc)
	}
	var last map[string]any
	pay := func(op string, after time.Duration) string {
		// Only the table's clock moves: the mark's lifetime is the table's TTL, and a permit is signed in real time.
		e.kv.Now = func() time.Time { return time.Now().Add(after) }
		args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: op})
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args, Approval: "exc_" + op})
		out, sc := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
		last = sc
		return out
	}
	e.b.LoseAnswer = true
	pay("op_lost", 0)
	e.b.LoseAnswer = false
	if out := pay("op_21m", 21*time.Minute); code(last) != browser.CodeOutcomeUnknown || len(e.b.Orders) != 1 {
		t.Fatalf("21 minutes after: %s (orders %v)", out, e.b.Orders)
	}
	if out := pay("op_23h", 23*time.Hour); code(last) != browser.CodeOutcomeUnknown || len(e.b.Orders) != 1 {
		t.Fatalf("23 hours after: %s (orders %v)", out, e.b.Orders)
	}
	if out := pay("op_25h", 25*time.Hour); code(last) != browser.CodeTotalUnread || len(e.b.Orders) != 1 {
		t.Fatalf("25 hours after, with no read: %s (orders %v)", out, e.b.Orders)
	}
}

// A read checkout counts for a pay for 15 minutes: at 14:59 it pays, at 15:01 it is refused (no total was read), so a
// checkout is never payable after the 15-minute approval it was read for has lapsed. KIFF's hold lapses at 15 too.
// Category: money
// Attack: a pay on a checkout read 16 to 20 minutes ago, after the owner's approval has lapsed
// Refused: total_unread at 15:01, nothing sent; at 14:59 the same pay goes through once
func TestACheckoutReadCountsForFifteenMinutes(t *testing.T) {
	for _, c := range []struct {
		name   string
		after  time.Duration
		orders int
		code   string
	}{{"14:59", 14*time.Minute + 59*time.Second, 1, ""}, {"15:01", 15*time.Minute + time.Second, 0, browser.CodeTotalUnread}} {
		t.Run(c.name, func(t *testing.T) {
			e := newBrowserEnv(t)
			g := e.connect()
			// Planque's clock stands still on a whole second, so the table's expiry is exactly 15:00; only the table's
			// clock moves on, and a permit is signed in real time.
			base := time.Unix(time.Now().Unix(), 0)
			e.s.Now = func() time.Time { return base }
			e.kv.Now = func() time.Time { return base }
			if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
				t.Fatalf("read: %v", sc)
			}
			e.kv.Now = func() time.Time { return base.Add(c.after) }
			args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: "op_" + c.name})
			permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: "op_" + c.name, Args: args, IssuedAt: base, Approval: "exc_" + c.name})
			out, sc := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
			if code(sc) != c.code || len(e.b.Orders) != c.orders {
				t.Fatalf("a pay %s after the read: %s (orders %v)", c.name, out, e.b.Orders)
			}
		})
	}
}

// 6. The sign-in a site keeps is sealed between calls and goes with the
// grant: the next session starts with it; disconnecting deletes it, the
// screenshots, and the tools in KIFF.
func TestBrowserSignInIsSealedAndDeletedOnDisconnect(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopLogin, "read": e.b.ReadOf(fakes.ShopLogin), "steps": []map[string]any{{"fill": "user", "value": "ana@example.com", "label": "Email"}, {"click": "signin", "label": "Sign in", "submits": true}}}); code(sc) != "" {
		t.Fatalf("sign in: %v", sc)
	}
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != "" {
		t.Fatalf("open: %v", sc)
	}
	if n := len(e.b.States); n != 2 || e.b.States[0] != nil || !strings.Contains(string(e.b.States[1]), "signed-in") {
		t.Fatalf("the sessions' states: %q", e.b.States)
	}
	// At rest, the state is sealed: the record doesn't carry it in the clear.
	it, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skBrowserState(g.ID))
	if !ok || strings.Contains(string(it.Doc), "signed-in") {
		t.Fatalf("the state at rest: %v %s", ok, it.Doc)
	}
	if len(e.ev.Refs(tenant, g.ID)) != 1 {
		t.Fatalf("screenshots: %v", e.ev.Refs(tenant, g.ID))
	}
	st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, e.phone, "")
	if st != 200 || out["status"] != StatusDisconnected || !strings.Contains(mustJSONString(out["events"]), "sign-ins and screenshots deleted") {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if _, ok, _ := e.kv.Get(context.Background(), tenantPK(tenant), skBrowserState(g.ID)); ok {
		t.Error("the sign-in stayed")
	}
	if len(e.ev.Refs(tenant, g.ID)) != 0 {
		t.Error("the screenshots stayed")
	}
	if _, ok := e.kiff.Connection(tenant, "browse_pay"); ok {
		t.Error("browse_pay stayed in KIFF")
	}
}

// 7. A worker that can't be reached: refused, not sent, said so.
func TestBrowserDownIsRefusedNotSent(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	e.b.Down = true
	out, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	if code(sc) != browser.CodeBrowserDown || sc["sent"] != false || !strings.Contains(out, "Nothing was sent") {
		t.Errorf("down: %s", out)
	}
	e.s.Browser = nil
	if _, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem}); code(sc) != browser.CodeBrowserDown {
		t.Errorf("no worker: %v", sc)
	}
}

// The relay lists the browser's tools with planque's operation argument,
// as KIFF's gateway reads them.
func TestBrowserToolsAreListedAtTheRelay(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	tools := (&Server{}).relayTools(g)
	if len(tools) != 5 {
		t.Fatalf("listed: %d", len(tools))
	}
	for _, tl := range tools {
		if !strings.Contains(string(tl.InputSchema), IdempotencyArg) || tl.Title == "" || tl.Description == "" {
			t.Errorf("%s: %s", tl.Name, tl.InputSchema)
		}
	}
}

// An installation with no browser worker lists no browser and takes no
// browser grant: the kind is there only where it can run.
func TestBrowserIsOfferedOnlyWithAWorker(t *testing.T) {
	e := newBrowserEnv(t)
	e.s.Browser = nil
	st, cat := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
	if st != 200 || strings.Contains(mustJSONString(cat), `"entry":"browser"`) {
		t.Errorf("catalog without a worker: %d %v", st, cat)
	}
	if st, out := e.ask(map[string]any{"browser": map[string]any{"currency": "EUR"}}); st != 503 {
		t.Errorf("ask without a worker: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/browser/grants", map[string]any{"agent": "agt_taf"}, nil, "k_setup"); st != 503 {
		t.Errorf("catalog ask without a worker: %d %v", st, out)
	}
}

// A step runs only on the page the taf read: refs restart on each new
// document, so after a navigation a ref can name an element the taf never
// read. Here the second step names the cart's checkout link, which exists
// only on the page the first step's click led to; before, it ran there.
func TestBrowserActStepsRunOnlyOnThePageTheyWereReadFrom(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	before := e.b.Sent()
	out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": e.b.ReadOf(fakes.ShopItem), "steps": []map[string]any{{"click": "add", "label": "Add to cart", "submits": true}, {"click": "checkout", "label": "Checkout", "to": "shop.example"}}})
	if code(sc) != browser.CodePageChanged || sc["steps_run"] != float64(1) || !strings.Contains(out, "step 2") {
		t.Fatalf("a step on a page the taf didn't read: %s", out)
	}
	if e.b.Sent() != before+1 {
		t.Fatalf("steps sent: %d, want 1", e.b.Sent()-before)
	}
	// A read of another page refuses the first step, and nothing is sent.
	before = e.b.Sent()
	out, sc = e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": e.b.ReadOf(fakes.ShopCart), "steps": []map[string]any{{"fill": "qty", "value": "2", "label": "Quantity"}}})
	if code(sc) != browser.CodePageChanged || !strings.Contains(out, "step 1") || e.b.Sent() != before {
		t.Fatalf("a stale read: %s (sent %d)", out, e.b.Sent()-before)
	}
	// browse_open gives the read browse_act names, and the page's read after.
	_, sc = e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	if sc["read"] != e.b.ReadOf(fakes.ShopItem) {
		t.Fatalf("browse_open's read: %v", sc["read"])
	}
	_, sc = e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": sc["read"], "steps": []map[string]any{{"fill": "qty", "value": "2", "label": "Quantity"}}})
	if code(sc) != "" || sc["read"] != e.b.ReadOf(fakes.ShopItem) {
		t.Fatalf("a fill keeps the read (values aren't in it): %v", sc)
	}
}

// The same address re-rendered between the read and the step: a field
// replaced by another, or two fields whose refs swapped, makes another
// read, and the step is refused with nothing run.
func TestBrowserActRefusesAPageThatChangedAtTheSameAddress(t *testing.T) {
	form := func(fields ...browser.Field) browser.Page {
		return browser.Page{URL: fakes.ShopItem, Title: "Blue pen", Fields: fields,
			Controls: []browser.Control{{Ref: "add", Label: "Add to cart", Submits: true}}}
	}
	for _, c := range []struct {
		name         string
		read, acting browser.Page
	}{
		{"a field replaced", form(browser.Field{Ref: "qty", Label: "Quantity", Kind: "number"}), form(browser.Field{Ref: "qty", Label: "Gift note", Kind: "text"})},
		{"two refs swapped", form(browser.Field{Ref: "f1", Label: "Quantity", Kind: "number"}, browser.Field{Ref: "f2", Label: "Gift note", Kind: "text"}),
			form(browser.Field{Ref: "f2", Label: "Quantity", Kind: "number"}, browser.Field{Ref: "f1", Label: "Gift note", Kind: "text"})},
	} {
		t.Run(c.name, func(t *testing.T) {
			e := newBrowserEnv(t)
			g := e.connect()
			e.b.Pages = map[string]browser.Page{fakes.ShopItem: c.read}
			_, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
			read := sc["read"]
			e.b.Pages[fakes.ShopItem] = c.acting
			before := e.b.Sent()
			out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": read, "steps": []map[string]any{{"fill": "qty", "value": "2", "label": "Quantity"}}})
			if code(sc) != browser.CodePageChanged || !strings.Contains(out, "step 1") || e.b.Sent() != before {
				t.Fatalf("a re-rendered page: %s (sent %d)", out, e.b.Sent()-before)
			}
		})
	}
}

// A pay runs only on its owner's approval, whatever a Card's room: a
// permit KIFF minted from a Card's room names no approval, and pays
// nothing, before any browser starts. The same checkout pays once the
// owner's approval is named.
func TestAPayWithoutTheOwnersApprovalPaysNothing(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); sc["total"] != "2.50" {
		t.Fatalf("checkout: %v", sc)
	}
	sessions := e.b.Sessions
	op := "op_pay_room"
	args := mustJSON(map[string]any{"url": fakes.ShopCheckout, "ref": "place", IdempotencyArg: op})
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_pay", Sub: g.Request.Agent, Op: op, Args: args})
	out, sc := textOf(e.s.call(context.Background(), g, "browse_pay", args, permit))
	if code(sc) != browser.CodePayNotApproved || sc["sent"] != false || len(e.b.Orders) != 0 || e.b.Sessions != sessions {
		t.Fatalf("a pay a Card's room allowed: %s (orders %v, sessions %d → %d)", out, e.b.Orders, sessions, e.b.Sessions)
	}
	if out, sc := e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != "" || len(e.b.Orders) != 1 {
		t.Fatalf("the approved pay: %s (orders %v)", out, e.b.Orders)
	}
}

// A browse_act with no read is told plainly where the read comes from, and
// runs nothing.
func TestBrowserActWithoutAReadSaysWhereItComesFrom(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": "", "steps": []map[string]any{{"click": "add", "label": "Add to cart", "submits": true}}})
	if code(sc) != browser.CodePageChanged || !strings.Contains(out, "needs the read from browse_open") || len(e.b.Steps) != 0 {
		t.Fatalf("no read: %s", out)
	}
}

// The account's connections say a browser connection's signed sites, so tafs
// cloud gives a taf the page the person sends with a message only from one of
// them (a connection of any other kind says none:
// TestSetupsListingSaysEachConnectionsKind).
func TestTheConnectionsSayABrowserConnectionsSignedSites(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/connections", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("connections: %d %v", st, out)
	}
	byKind := map[string]map[string]any{}
	for _, c := range out["connections"].([]any) {
		v := c.(map[string]any)
		byKind[v["kind"].(string)] = v
	}
	sites, _ := byKind[KindBrowser]["sites"].([]any)
	if _, any := byKind[KindBrowser]["any_site"]; len(sites) != 1 || sites[0] != "shop.example" || any {
		t.Fatalf("the browser connection (%s): %v", g.ID, byKind[KindBrowser])
	}
}

// A browser connection signed for no sites works on any public site, and
// says so in words of its own, any_site: never as a missing list, which is
// also what every other kind (and an older planque) sends.
func TestABrowserConnectionForAnySiteSaysSo(t *testing.T) {
	for name, b := range map[string]*BrowserRequest{"no sites": {Currency: "EUR"}, "no browser limits": nil} {
		g := &Grant{Request: Request{Agent: "agt_taf", Kind: KindBrowser, Browser: b}, Status: StatusConnected}
		c := connectionOf(g)
		if !c.AnySite || c.Sites != nil {
			t.Fatalf("%s: %+v", name, c)
		}
	}
	key := connectionOf(&Grant{Request: Request{Agent: "agt_taf", API: &APIRequest{URL: "https://api.example/openapi.json"}}})
	if key.AnySite || key.Sites != nil {
		t.Fatalf("a key grant: %+v", key)
	}
}

// What a held browse_act shows the owner is each step's ref named by the
// page's own words, and those words are checked before the step runs: a step
// that calls the shop's "Add to cart" a "Save for later", or that leaves out
// that it submits a form, runs nothing; and a page that changes a button's
// words between the read and the act runs nothing either (red-team, 27 Nov).
func TestBrowserActRunsOnlyStepsNamedByThePagesWords(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	_, sc := e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	read := sc["read"]
	for _, c := range []struct {
		name string
		step map[string]any
	}{
		{"another control's words", map[string]any{"click": "add", "label": "Save for later", "submits": true}},
		{"a submit not said to submit", map[string]any{"click": "add", "label": "Add to cart"}},
		{"no words at all", map[string]any{"click": "add"}},
		{"a field under another name", map[string]any{"fill": "qty", "value": "2", "label": "Gift note"}},
	} {
		before := e.b.Sent()
		out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": read, "steps": []map[string]any{c.step}})
		if code(sc) != browser.CodeLabelMismatch || sc["sent"] != false || e.b.Sent() != before || !strings.Contains(out, "step 1") {
			t.Fatalf("%s: %s (sent %d)", c.name, out, e.b.Sent()-before)
		}
	}
	// The page's own words run.
	if out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": read, "steps": []map[string]any{{"click": "add", "label": "Add to cart", "submits": true}}}); code(sc) != "" || sc["steps_run"] != float64(1) {
		t.Fatalf("the page's words: %s", out)
	}
	// The button's words change between the read and the act: nothing runs.
	item := func(button string) browser.Page {
		return browser.Page{URL: fakes.ShopItem, Title: "Blue pen",
			Fields:   []browser.Field{{Ref: "qty", Label: "Quantity", Kind: "number", Value: "1"}},
			Controls: []browser.Control{{Ref: "add", Label: button, Submits: true}}}
	}
	e.b.Pages = map[string]browser.Page{fakes.ShopItem: item("Add to cart")}
	_, sc = e.call(g, "browse_open", map[string]any{"url": fakes.ShopItem})
	read = sc["read"]
	e.b.Pages[fakes.ShopItem] = item("Buy now")
	before := e.b.Sent()
	out, sc := e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": read, "steps": []map[string]any{{"click": "add", "label": "Add to cart", "submits": true}}})
	if code(sc) != browser.CodePageChanged || e.b.Sent() != before {
		t.Fatalf("a button's words changed after the read: %s (sent %d)", out, e.b.Sent()-before)
	}
}

// A purchase case is built from the checkout read: shipping, tax, the
// merchant and the delivery address are answered and kept with the lines
// and total, and a pay on a checkout whose delivery address changed since
// the read sends nothing.
func TestBrowserCheckoutCarriesTheDeliveryAddress(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	e.b.Delivery = &browser.Address{Name: "Acme SAS", Line1: "12 rue des Lilas", PostalCode: "69003", City: "Lyon", Country: "FR"}
	out, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout})
	d, _ := sc["delivery"].(map[string]any)
	if code(sc) != "" || d["line1"] != "12 rue des Lilas" || d["postal_code"] != "69003" || d["city"] != "Lyon" || d["country"] != "FR" {
		t.Fatalf("checkout: %s", out)
	}
	e.b.Delivery = &browser.Address{Line1: "1 avenue Foch", PostalCode: "75116", City: "Paris", Country: "FR"}
	out, sc = e.call(g, "browse_pay", map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
	if code(sc) != browser.CodeCheckoutChanged || sc["sent"] != false || len(e.b.Orders) != 0 || !strings.Contains(out, "delivery address") {
		t.Fatalf("pay after the address changed: %s", out)
	}
}
