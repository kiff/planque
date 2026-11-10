package planque

import (
	"encoding/json"
	"slices"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/profile"
)

// Stripe refunds through KIFF's governed path. A refund's money rules are
// reviewed code, never a model's draft: when a managed Stripe connection
// asks for the refund tool, planque supplies the whole profile itself, the
// way browserProfile supplies the browser's. That is the refund operation
// (its amount, its marker), its case block (the reads KIFF builds the
// refund's case from, and the refund event it follows), and the three case
// reads, granted with it. The caller asks the four tools, read from the
// app's catalog like any other, and no profile. Such a connection holds
// those four and no other: the taf's other Stripe tools go on another
// connection, with the profile Coach drafts.
//
// It is offered only when the installation turns it on (StripeRefunds):
// the paths below are the Stripe objects' own, and they're used in
// production only once a read of each case tool on a test-mode connection
// has shown the connector answers with them.

// stripeRefundUnit is the currency the Card counts a refund in. A connection
// takes no unit read (deriveProfileFields), so step one is euros: KIFF's
// case refuses a refund on a charge in any other currency, as a currency
// mismatch, before anything is counted or sent.
const stripeRefundUnit = "EUR"

// StripeRefundsEntry names planque's reviewed profile for it.
const StripeRefundsEntry = "stripe-refunds@1"

// The managed connector's Stripe tools it uses, by their provider slugs.
const (
	stripeToolkit       = "stripe"
	stripeRefund        = "STRIPE_CREATE_REFUND"
	stripeGetAccount    = "STRIPE_GET_ACCOUNT"
	stripeGetCharge     = "STRIPE_RETRIEVE_CHARGE"
	stripeListRefunds   = "STRIPE_LIST_REFUNDS"
	stripeRefundCreated = "STRIPE_REFUND_CREATED"
)

// stripeRefundTools are the four tools, the refund last.
var stripeRefundTools = []string{stripeGetAccount, stripeGetCharge, stripeListRefunds, stripeRefund}

// stripeAlias is the alias the connector's catalog lists slug under.
func stripeAlias(slug string) string {
	return connector.Tool{Toolkit: stripeToolkit, Slug: slug}.Alias()
}

// stripeRefundsGrant reports whether g is a managed Stripe connection.
func stripeRefundsGrant(g *Grant) bool {
	c := g.Request.Connector
	return c != nil && c.Toolkit == stripeToolkit
}

// asksStripeRefund reports whether tools ask for the refund tool.
func asksStripeRefund(tools []ToolRequest) bool {
	refund := stripeAlias(stripeRefund)
	return slices.ContainsFunc(tools, func(t ToolRequest) bool { return t.Tool == refund })
}

// errStripeRefundsSet: a refund connection holds the refund and its three
// case reads, all four asked, and no other tool.
var errStripeRefundsSet = apiError{400, "a Stripe connection with the refund tool asks exactly four tools from the app's catalog: " +
	stripeRefund + " and its three case reads, " + stripeGetAccount + ", " + stripeGetCharge + " and " + stripeListRefunds +
	". Ask the taf's other Stripe tools on another connection"}

// stripeRefundsAsk is the request for a refund connection: the four tools,
// in a fixed order, each under the name the caller gave it, and planque's
// reviewed profile for them. A caller's own profile is refused, since this
// one is planque's.
func stripeRefundsAsk(tools []ToolRequest, callerProfile json.RawMessage, vendor string) ([]ToolRequest, json.RawMessage, error) {
	if len(callerProfile) > 0 {
		return nil, nil, apiError{400, "the refund tool's profile is planque's reviewed one (" + StripeRefundsEntry + "): send the tools without a profile"}
	}
	if len(tools) != len(stripeRefundTools) {
		return nil, nil, errStripeRefundsSet
	}
	byAlias := map[string]ToolRequest{}
	for _, t := range tools {
		byAlias[t.Tool] = t
	}
	out := make([]ToolRequest, 0, len(stripeRefundTools))
	for _, slug := range stripeRefundTools {
		t, asked := byAlias[stripeAlias(slug)]
		if !asked {
			return nil, nil, errStripeRefundsSet
		}
		out = append(out, t)
	}
	raw, err := stripeRefundsProfile(out, vendor)
	return out, raw, err
}

// stripeRefundsProfile is the reviewed profile for the four tools, under
// their names before planque allocates its own (renameProfile follows).
func stripeRefundsProfile(tools []ToolRequest, vendor string) (json.RawMessage, error) {
	names := toolNames(tools)
	account, charge, refunds, refund := names[0], names[1], names[2], names[3]
	chargeArgs := map[string]string{"charge_id": "$args.charge"}
	refundsArgs := map[string]string{"charge": "$args.charge"}
	p := profile.Profile{V: profile.Version, Vendor: vendor, Transport: "rest", Pin: profile.Pin{Kind: "none"},
		Source: "catalog", CatalogEntry: StripeRefundsEntry,
		Operations: []profile.Operation{
			{Tool: account, Title: "Look up the Stripe account", Kind: "read"},
			{Tool: charge, Title: "Look up a charge", Kind: "read", Resource: &profile.Resource{Arg: "charge_id", Noun: "charge"}},
			{Tool: refunds, Title: "List a charge's refunds", Kind: "read", Resource: &profile.Resource{Arg: "charge", Noun: "charge"}},
			{Tool: refund, Title: "Refund a charge", Kind: "write", Effect: profile.EffectMovesMoney, Shape: "single",
				Resource:          &profile.Resource{Arg: "charge", Noun: "charge"},
				Amount:            &profile.Amount{From: "argument", Arg: "amount", Format: "integer_minor", Scale: 2, Unit: stripeRefundUnit},
				Marker:            &profile.Marker{Field: "metadata.kiff_op"},
				RequiredArguments: []string{"charge", "amount"},
				Case: &profile.Case{Kind: profile.CaseKindRefund, ChargeArg: "charge",
					AccountRead: profile.Read{Tool: account, Path: "id"},
					ChargeRead: profile.CaseChargeRead{Tool: charge, Args: chargeArgs,
						Currency: "currency", Captured: "amount_captured", Disputed: "disputed"},
					RefundsRead: profile.CaseRefundsRead{Tool: refunds, Args: refundsArgs,
						List: "data", ID: "id", Amount: "amount", Status: "status", Marker: "metadata.kiff_op"},
					RefundEvent: &profile.CaseRefundEvent{Event: stripeRefundCreated, Object: "data.object", Charge: "charge"}}},
		}}
	return json.Marshal(p)
}

// draftedCase reports whether a profile that isn't a reviewed one carries
// a case block: a refund's money rules are never a model's draft.
func draftedCase(p profile.Profile) bool {
	if p.Source == "catalog" {
		return false
	}
	return slices.ContainsFunc(p.Operations, func(o profile.Operation) bool { return o.Case != nil })
}
