package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// While paying is off in the hosted browser (PLANQUE_BROWSER_PAY not on, the
// default), a browser grant asked with every tool gets no browse_pay, one
// asked with browse_pay is refused, and a grant that holds it from before
// pays nothing: refused before any browser starts, whatever the approval.
// The hosted browser still opens, reads and reads a checkout.
func TestTheHostedBrowserPaysNothingWhilePayingIsOff(t *testing.T) {
	e := newBrowserEnv(t)
	// A grant connected while paying was on, with browse_pay.
	g := e.connect()
	e.s.BrowserPay = false

	st, out := e.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}}})
	if st != 201 {
		t.Fatalf("every tool, paying off: %d %v", st, out)
	}
	if fresh, err := e.s.loadGrant(t.Context(), tenant, out["id"].(string)); err != nil {
		t.Fatal(err)
	} else {
		for _, tool := range fresh.Request.Tools {
			if tool.Tool == browsePay {
				t.Fatalf("browse_pay was offered with paying off: %+v", fresh.Request.Tools)
			}
		}
		if len(fresh.Request.Tools) == 0 {
			t.Fatal("no tool at all was offered")
		}
	}
	st, out = e.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}, "currency": "EUR"},
		"tools": []map[string]any{{"tool": browseCheckout}, {"tool": browsePay}}})
	if st != 400 || !strings.Contains(mustJSONString(out), "paying in the hosted browser is off") {
		t.Fatalf("browse_pay asked for, paying off: %d %v", st, out)
	}

	if _, sc := e.call(g, browseCheckout, map[string]any{"url": fakes.ShopCheckout}); code(sc) != "" {
		t.Fatalf("a checkout read, paying off: %v", sc)
	}
	sessions := e.b.Sessions
	out2, sc := e.call(g, browsePay, map[string]any{"url": fakes.ShopCheckout, "ref": "place"})
	if code(sc) != CodePayOff || sc["sent"] != false || len(e.b.Orders) != 0 || e.b.Sessions != sessions || !strings.Contains(out2, "nothing was paid") {
		t.Fatalf("a pay, paying off: %s %v", out2, sc)
	}
	// Turned on, the same grant pays, once.
	e.s.BrowserPay = true
	if _, sc := e.call(g, browsePay, map[string]any{"url": fakes.ShopCheckout, "ref": "place"}); code(sc) != "" || len(e.b.Orders) != 1 {
		t.Fatalf("a pay, paying on: %v", sc)
	}
}
