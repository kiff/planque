package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// describePay is the gateway's describe of a pay: browse_pay with its
// header, under a permit for the op given.
func (e *browserEnv) describePay(g *Grant, header, op string, args map[string]any) (string, map[string]any) {
	e.t.Helper()
	h := http.Header{}
	h.Set(CheckoutDescribeHeader, header)
	args[IdempotencyArg] = op
	raw := mustJSON(args)
	ctx := withCheckoutHeader(context.Background(), h)
	return textOf(e.s.call(ctx, g, browsePay, raw, e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: browsePay, Sub: g.Request.Agent, Op: op, Args: raw})))
}

// The gateway's describe of a pay answers the checkout planque read and
// kept, signed under its own label, starting no browser and paying
// nothing; only under a permit for that describe, only after a read, and
// never for a checkout whose pay ended unknown.
func TestAPayIsDescribedFromItsKeptRead(t *testing.T) {
	e := newBrowserEnv(t)
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := stepkey.Local{Key: key}
	e.s.ChromeSigner = signer
	g := e.connect()
	e.b.Delivery = &browser.Address{Name: "Acme SAS", Line1: "12 rue des Lilas", PostalCode: "69003", City: "Lyon", Country: "FR"}
	pay := map[string]any{"url": fakes.ShopCheckout, "ref": "place"}

	// Before a read: no block.
	if _, sc := e.describePay(g, "op_pay0", "op_pay0:describe", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeTotalUnread {
		t.Fatalf("before a read: %v", sc)
	}
	if _, sc := e.call(g, "browse_checkout", map[string]any{"url": fakes.ShopCheckout}); code(sc) != "" {
		t.Fatalf("read: %v", sc)
	}
	sessions := e.b.Sessions
	// A permit for another op, a bad header, or steps: refused.
	for name, c := range map[string][3]string{
		"another op":    {"op_pay1", "op_pay1"},
		"a bad header":  {"op pay 1", "op pay 1:describe"},
		"the pay's own": {"op_pay1", "op_pay1:open"},
	} {
		if _, sc := e.describePay(g, c[0], c[1], map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) == "" {
			t.Errorf("%s: described", name)
		}
	}
	if _, sc := e.describePay(g, "op_pay5", "op_pay5:describe", map[string]any{"url": fakes.ShopCheckout}); sc["browser_block"] != nil {
		t.Errorf("no ref: described %v", sc)
	}

	_, sc := e.describePay(g, "op_pay1", "op_pay1:describe", pay)
	bb, _ := sc["browser_block"].(map[string]any)
	if code(sc) != "" || bb == nil {
		t.Fatalf("describe: %v", sc)
	}
	raw, _ := json.Marshal(bb["block"])
	var b checkoutBlock
	if err := json.Unmarshal(raw, &b); err != nil {
		t.Fatal(err)
	}
	if b.V != labelCheckout || b.Tenant != tenant || b.Op != "op_pay1" || b.Grant != g.ID || b.Site != "shop.example" || b.Ref != "place" ||
		b.Total != "2.50" || b.Currency != "EUR" || len(b.Lines) != 1 || b.Delivery == nil || b.Delivery.Line1 != "12 rue des Lilas" || b.ReadAt.IsZero() {
		t.Fatalf("block: %+v", b)
	}
	if !stepkey.Verify(signer.PublicKey(), labelCheckout, b, bb["sig"].(string)) {
		t.Fatal("the block's signature doesn't verify under the checkout label")
	}
	if stepkey.Verify(signer.PublicKey(), labelDescribe, b, bb["sig"].(string)) {
		t.Fatal("the block verifies as a sealed fill's")
	}
	// No browser was started, nothing was paid.
	if e.b.Sessions != sessions || len(e.b.Orders) != 0 {
		t.Fatalf("the describe started %d sessions, placed %d orders", e.b.Sessions-sessions, len(e.b.Orders))
	}
	// Too long after the read: no block.
	var rec checkoutRecord
	it, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skCheckout(g.ID, fakes.ShopCheckout), &rec)
	if err != nil {
		t.Fatal(err)
	}
	rec.ReadAt = rec.ReadAt.Add(-checkoutKept - time.Minute)
	if err := putDoc(context.Background(), e.s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}, rec, store.Cond{Ver: &it.Ver}); err != nil {
		t.Fatal(err)
	}
	if _, sc := e.describePay(g, "op_pay6", "op_pay6:describe", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeTotalUnread {
		t.Errorf("stale: %v", sc)
	}
	// A checkout whose earlier pay ended unknown: no block.
	it, _ = getDoc(context.Background(), e.s.KV, tenantPK(tenant), skCheckout(g.ID, fakes.ShopCheckout), &rec)
	rec.ReadAt, rec.UnknownAt = time.Now().UTC(), time.Now().UTC()
	if err := putDoc(context.Background(), e.s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}, rec, store.Cond{Ver: &it.Ver}); err != nil {
		t.Fatal(err)
	}
	if _, sc := e.describePay(g, "op_pay7", "op_pay7:describe", map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != browser.CodeOutcomeUnknown {
		t.Errorf("unknown: %v", sc)
	}
	// Only browse_pay describes.
	h := http.Header{}
	h.Set(CheckoutDescribeHeader, "op_x")
	args := map[string]any{"url": fakes.ShopCheckout, IdempotencyArg: "op_x:describe"}
	raw = mustJSON(args)
	if out, sc := textOf(e.s.call(withCheckoutHeader(context.Background(), h), g, "browse_checkout", raw, e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "browse_checkout", Sub: g.Request.Agent, Op: "op_x:describe", Args: raw}))); code(sc) != CodeBadArguments || !strings.Contains(out, "only browse_pay") {
		t.Errorf("another tool: %s", out)
	}
}
