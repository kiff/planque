package connector

import "strings"

// NeverServedTools are the tools planque never serves for a managed app, by
// the app's slug (lower case) and then the tool's: each one moves money or
// closes a bill for good (it voids a bill, writes a debt off, deletes one,
// makes one owed or charges the customer). Like the platform's meta tools
// (Validate), they never come back from a search or an app's tools, can't
// be asked for or chosen, and a call to one is never sent, even on a
// connection that holds it from before. Keyed by app, so another app's tool of the same
// name is served as it is. Tests may swap it; nothing else changes it.
var NeverServedTools = map[string][]string{
	"stripe": {
		"STRIPE_POST_V1_INVOICES_SEND_VOID_AN_INVOICE",
		"STRIPE_VOID_INVOICE",
		"STRIPE_POST_INVOICES_INVOICE_VOID",
		"STRIPE_POST_V1_INVOICES_VOID_AN_INVOICE",
		"STRIPE_POST_V1_INVOICES_FINALIZE_VOID_AN_INVOICE",
		"STRIPE_DELETE_V1_INVOICES_VOID_AN_INVOICE",
		"STRIPE_VOID_CREDIT_NOTE",
		"STRIPE_POST_BILLING_CREDIT_GRANTS_VOID_CREDIT_GRANT",
		"STRIPE_DELETE_INVOICE",
		"STRIPE_FINALIZE_INVOICE",
		"STRIPE_PAY_INVOICE",
		"STRIPE_MARK_INVOICE_UNCOLLECTIBLE",
	},
}

// NeverServed says planque never serves the tool slug of the app toolkit.
func NeverServed(toolkit, slug string) bool {
	for _, s := range NeverServedTools[strings.ToLower(toolkit)] {
		if strings.EqualFold(s, slug) {
			return true
		}
	}
	return false
}

// NeverServedAlias is the slug of the app's never-served tool whose gateway
// name (Tool.Alias) is alias, when there is one.
func NeverServedAlias(toolkit, alias string) (string, bool) {
	for _, s := range NeverServedTools[strings.ToLower(toolkit)] {
		if (Tool{Toolkit: toolkit, Slug: s}).Alias() == alias {
			return s, true
		}
	}
	return "", false
}
