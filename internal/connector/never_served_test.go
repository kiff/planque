package connector

import (
	"encoding/json"
	"testing"
)

// Stripe's tools that void, delete, finalize, pay or write off a bill are
// never served for the managed app, in any case; its other tools are, and a
// tool of the same name in another app is (the list is per app).
func TestAnAppsNeverServedToolsAreNotDefinitions(t *testing.T) {
	tool := func(toolkit, slug string) Tool {
		return Tool{Slug: slug, Toolkit: toolkit, Version: "20261101_00", Name: slug, InputSchema: json.RawMessage(`{"type":"object","properties":{}}`)}
	}
	for _, slug := range NeverServedTools["stripe"] {
		if tool("stripe", slug).Validate() != ErrDefinition || !NeverServed("Stripe", slug) {
			t.Errorf("%s is served", slug)
		}
		if _, ok := NeverServedAlias("stripe", tool("stripe", slug).Alias()); !ok {
			t.Errorf("%s isn't found by its gateway name", slug)
		}
	}
	if !NeverServed("stripe", "stripe_pay_invoice") {
		t.Error("a slug in lower case is served")
	}
	// Negatives: Stripe's ordinary tools, and another app's tool of the same name.
	for _, ok := range []Tool{tool("stripe", "STRIPE_LIST_INVOICES"), tool("stripe", "STRIPE_SEND_INVOICE"), tool("stripe", "STRIPE_CREATE_REFUND"), tool("other_app", "STRIPE_VOID_INVOICE")} {
		if ok.Validate() != nil || NeverServed(ok.Toolkit, ok.Slug) {
			t.Errorf("%s of %s is refused", ok.Slug, ok.Toolkit)
		}
		if _, refused := NeverServedAlias(ok.Toolkit, ok.Alias()); refused {
			t.Errorf("%s of %s is refused by its gateway name", ok.Slug, ok.Toolkit)
		}
	}
}
