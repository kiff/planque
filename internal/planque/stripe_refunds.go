package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/profile"
)

// Stripe refunds through KIFF's governed path, on the direct Stripe grant
// (planque's own Stripe entry, with test mode). A refund's money rules are
// reviewed code, never a model's draft: when a direct Stripe connection asks
// for the refund, planque supplies the whole request itself, the way
// browserProfile supplies the browser's. The caller asks one tool,
// stripe_api_write narrowed to PostRefunds with the mode the owner picked
// (livemode), and no profile. planque adds:
//   - stripe_api_details, Setup's read of an operation's fields: once the
//     owner approved it, planque reads the narrowed tools' fields with it
//     (describe.go), and checks their paths, the marker's included, against
//     Stripe's own description before the refund is first sent;
//   - the case reads: stripe_api_read narrowed to GetAccountsAccount,
//     GetChargesCharge and GetRefunds;
//   - every narrowed tool's arguments, account and mode. The account is the
//     one Stripe's sign-in answer named (the grant's account, from
//     stripe_user_id), never the caller's;
//   - the profile, stripe-refunds@1: the refund in euros, with its marker
//     and its case block.
//
// Such a connection holds those five tools and no other: the taf's other
// Stripe tools go on another connection, with the profile Coach drafts.
//
// It is offered only when the installation turns it on (StripeRefunds).

// stripeRefundUnit is the currency the Card counts a refund in. A connection
// takes no unit read (deriveProfileFields), so step one is euros: KIFF's
// case refuses a refund on a charge in any other currency, as a currency
// mismatch, before anything is counted or sent. Stripe writes a currency in
// lower case (eur), and KIFF compares it without case.
const stripeRefundUnit = "EUR"

// StripeEntry is Stripe in the vendor catalog: planque's own entry for a
// direct sign-in to Stripe's MCP server, which offers live mode, sandboxes
// and test mode on one consent page, so a taf can be tried on a test
// account before going live. Where it exists, the managed connector's
// Stripe isn't offered: the app shows one Stripe.
const StripeEntry = "stripe"

// stripeServer is Stripe's own MCP server.
const stripeServer = "https://mcp.stripe.com"

// stripeAccountField is where Stripe's token answer names the account the
// sign-in is for (Stripe Connect's OAuth). planque keeps it as the grant's
// account (oauth.Metadata.AccountField).
const stripeAccountField = "stripe_user_id"

// stripeAccountsTool is Stripe's own tool that lists the accounts a sign-in
// covers: {"accounts":[{"stripe_context":"acct_…","livemode":false,
// "name":"…"}]}, in the result's structuredContent and, as text, in its
// content, as a test-mode sign-in showed.
const stripeAccountsTool = "list_available_accounts_or_orgs"

// signedInAccount is the account a sign-in covers, read with the vendor's
// own accounts tool (oauth.Metadata.AccountsTool), once, by planque itself,
// right after the owner signed in, with the new token: never a caller's
// word. It returns the account and its mode only when the list has exactly
// one account; with none, several, or a failed read it returns "", and a
// refund connection then asks the owner to sign in to one account. The read
// is on the grant's record.
func (s *Server) signedInAccount(ctx context.Context, g *Grant, access string) (string, *bool) {
	if s.MCP == nil {
		return "", nil
	}
	res, err := s.MCP.CallTool(ctx, g.Request.Server, access, g.Vendor.AccountsTool, json.RawMessage(`{}`))
	if err != nil {
		s.log().WarnContext(ctx, "the sign-in's accounts couldn't be read", slog.String("grant", g.ID), slog.String("err", err.Error()))
		g.event(s.now(), "couldn't read which account the sign-in covers")
		return "", nil
	}
	accounts, ok := sessionAccounts(res)
	switch {
	case !ok:
		g.event(s.now(), "read which account the sign-in covers; the answer named none")
		return "", nil
	case len(accounts) != 1:
		g.event(s.now(), fmt.Sprintf("read which account the sign-in covers: %d accounts, so none is kept", len(accounts)))
		return "", nil
	}
	live := accounts[0].Live
	g.event(s.now(), "read which account the sign-in covers: one, kept as the connection's account")
	return accounts[0].Context, &live
}

// sessionAccount is one account a sign-in covers, as Stripe lists it.
type sessionAccount struct {
	Context string `json:"stripe_context"`
	Live    bool   `json:"livemode"`
}

// sessionAccounts reads an accounts tool's result: its structuredContent,
// or else the same JSON as its text content. ok is false for an error, an
// answer of another shape, or an entry whose id isn't an account's.
func sessionAccounts(res json.RawMessage) ([]sessionAccount, bool) {
	type list struct {
		Accounts *[]struct {
			Context string `json:"stripe_context"`
			Live    *bool  `json:"livemode"`
		} `json:"accounts"`
	}
	var out struct {
		Structured *list `json:"structuredContent"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if json.Unmarshal(res, &out) != nil || out.IsError {
		return nil, false
	}
	l := out.Structured
	if l == nil || l.Accounts == nil {
		for _, c := range out.Content {
			var t list
			if c.Type == "text" && json.Unmarshal([]byte(c.Text), &t) == nil && t.Accounts != nil {
				l = &t
				break
			}
		}
	}
	if l == nil || l.Accounts == nil {
		return nil, false
	}
	accounts := make([]sessionAccount, 0, len(*l.Accounts))
	for _, a := range *l.Accounts {
		if !strings.HasPrefix(a.Context, "acct_") || !oauth.ValidAccount(a.Context) || a.Live == nil {
			return nil, false
		}
		accounts = append(accounts, sessionAccount{Context: a.Context, Live: *a.Live})
	}
	return accounts, true
}

// stripeVendor is the name people read.
const stripeVendor = "Stripe"

// stripeIcon is the Stripe entry's logo, the one the managed connector's
// Stripe lists: the apps load it through tafs cloud, never from here.
const stripeIcon = "https://logos.composio.dev/api/stripe"

// stripeToolkit is the managed connector's Stripe, which the direct entry
// replaces.
const stripeToolkit = "stripe"

// StripeRefundsEntry names planque's reviewed profile for it.
const StripeRefundsEntry = "stripe-refunds@1"

// Stripe's MCP tools, and the operations a refund connection narrows them
// to, as stripe_api_details names them.
const (
	stripeDetails       = "stripe_api_details"
	stripeRead          = "stripe_api_read"
	stripeWrite         = "stripe_api_write"
	stripeOperation     = "stripe_api_operation_id"
	stripeContext       = "stripe_context"
	stripeLivemode      = "livemode"
	stripeGetAccount    = "GetAccountsAccount"
	stripeGetCharge     = "GetChargesCharge"
	stripeListRefunds   = "GetRefunds"
	stripeRefund        = "PostRefunds"
	stripeRefundsLimit  = 100 // GetRefunds' most; its default is 10
	stripeRefundsMarker = "metadata.kiff_op"
)

// narrowedParams is where a narrowed broad tool takes its operation's
// fields (describer.Params), and where its marker goes.
const narrowedParams = profile.NarrowedParams

// stripeNames are the names a refund connection's narrowed tools are given,
// by operation, before the grant's suffix (toolKey).
var stripeNames = map[string]string{stripeGetAccount: "stripe_account", stripeGetCharge: "stripe_charge",
	stripeListRefunds: "stripe_refunds", stripeRefund: "stripe_refund"}

// stripeRefundsGrant reports whether g is a direct Stripe connection.
func stripeRefundsGrant(g *Grant) bool {
	return g.Request.Connector == nil && g.Request.Kind != KindBrowser && sameServer(g.Request.Server, stripeServer)
}

// stripeOperationOf is the operation t pins on one of Stripe's broad tools.
func stripeOperationOf(t ToolRequest) string {
	var op string
	_ = json.Unmarshal(t.Fixed[stripeOperation], &op)
	return op
}

// asksStripeRefund reports whether tools ask for the refund.
func asksStripeRefund(tools []ToolRequest) bool {
	return slices.ContainsFunc(tools, func(t ToolRequest) bool {
		return t.Tool == stripeWrite && stripeOperationOf(t) == stripeRefund
	})
}

// errStripeRefundsAsk: a refund connection is asked as one tool.
var errStripeRefundsAsk = apiError{400, "a Stripe connection with refunds asks one tool and no profile: " + stripeWrite +
	" with fixed_arguments " + stripeOperation + " " + stripeRefund + " and " + stripeLivemode +
	" (false: test mode, the only one offered for now), and nothing else. planque adds the case reads and the profile (" + StripeRefundsEntry +
	"). Ask the taf's other Stripe tools on another connection"}

// stripeRefundsAsk is the request for a refund connection: the five tools,
// in a fixed order with the refund last, each named for grant, and
// planque's reviewed profile for them. account is the grant's, as Stripe's
// accounts tool named it right after the sign-in, with its mode. A caller's own profile is refused, since this
// one is planque's.
func stripeRefundsAsk(tools []ToolRequest, callerProfile json.RawMessage, account string, accountLive *bool, grant, vendor string) ([]ToolRequest, json.RawMessage, error) {
	if len(callerProfile) > 0 {
		return nil, nil, apiError{400, "the refund's profile is planque's reviewed one (" + StripeRefundsEntry + "): send the tool without a profile"}
	}
	if len(tools) != 1 || len(tools[0].Arguments) > 0 || len(tools[0].Fixed) != 2 {
		return nil, nil, errStripeRefundsAsk
	}
	asked := tools[0]
	var live bool
	if json.Unmarshal(asked.Fixed[stripeLivemode], &live) != nil || stripeOperationOf(asked) != stripeRefund {
		return nil, nil, errStripeRefundsAsk
	}
	if live {
		// Stripe may ask the owner to confirm each refund in its Dashboard
		// when the sign-in is a person's own and the mode is live; only test
		// mode has been seen to refund at once.
		return nil, nil, apiError{409, "refunds are offered in Stripe's test mode only for now (livemode false): in live mode, Stripe may ask you to confirm each refund in its Dashboard, and that hasn't been checked yet"}
	}
	if account == "" {
		return nil, nil, apiError{409, "this Stripe connection doesn't know its Stripe account: sign in to Stripe again, choosing one account in test mode, then ask for refunds"}
	}
	if accountLive != nil && *accountLive != live {
		return nil, nil, apiError{409, "this Stripe connection's account is in live mode, and refunds are offered in test mode only for now: sign in to Stripe again in test mode"}
	}
	out := stripeRefundTools(account, live, grant, asked.Name, asked.Title)
	raw, err := stripeRefundsProfile(out, vendor)
	return out, raw, err
}

// stripeRefundTools are the five tools of a refund connection on account,
// in live or test mode: Setup's describing read, the three case reads and
// the refund. Each is named for grant, since three are the same broad tool;
// the refund keeps the name and title its caller gave it.
func stripeRefundTools(account string, live bool, grant, name, title string) []ToolRequest {
	acct, _ := json.Marshal(account)
	mode, _ := json.Marshal(live)
	narrowed := func(tool, op, title string, extra map[string]json.RawMessage, args map[string]Argument) ToolRequest {
		opv, _ := json.Marshal(op)
		fixed := map[string]json.RawMessage{stripeOperation: opv, stripeContext: acct, stripeLivemode: mode}
		for k, v := range extra {
			fixed[k] = v
		}
		return ToolRequest{Tool: tool, Name: toolKey(stripeNames[op], grant), Title: title, Fixed: fixed, Arguments: args}
	}
	charge := Argument{Path: narrowedParams + ".charge", Type: "string", Required: true, Description: "The charge's id (ch_…)."}
	limit, _ := json.Marshal(stripeRefundsLimit)
	details := ToolRequest{Tool: stripeDetails, Name: toolKey(stripeDetails, grant), Title: "Read a Stripe operation's fields", SetupReads: true, ReadOnly: true}
	// The account read takes the account from the grant, at the operation's
	// id (GetAccountsAccount, as stripe_api_details lists it); a narrowed
	// tool still names one argument (checkNarrowing): Stripe's expand, which
	// KIFF never sends.
	account_ := narrowed(stripeRead, stripeGetAccount, "Look up the Stripe account",
		map[string]json.RawMessage{narrowedParams + ".id": acct},
		map[string]Argument{"expand": {Path: narrowedParams + ".expand", Type: "string", Description: "Leave out."}})
	chargeRead := narrowed(stripeRead, stripeGetCharge, "Look up a charge", nil,
		map[string]Argument{"charge": {Path: narrowedParams + ".id", Type: "string", Required: true, Description: "The charge's id (ch_…)."}})
	refunds := narrowed(stripeRead, stripeListRefunds, "List a charge's refunds",
		map[string]json.RawMessage{narrowedParams + ".limit": limit}, map[string]Argument{"charge": charge})
	refund := narrowed(stripeWrite, stripeRefund, title, nil, map[string]Argument{
		"charge": charge,
		"amount": {Path: narrowedParams + ".amount", Type: "integer", Required: true, Description: "The amount to refund, in cents."},
	})
	if name != "" {
		refund.Name = name
	}
	if refund.Title == "" {
		refund.Title = "Refund a charge"
	}
	return []ToolRequest{details, account_, chargeRead, refunds, refund}
}

// stripeRefundsProfile is the reviewed profile for the five tools, under
// their names before planque allocates its own (renameProfile follows).
func stripeRefundsProfile(tools []ToolRequest, vendor string) (json.RawMessage, error) {
	names := toolNames(tools)
	details, account, charge, refunds, refund := names[0], names[1], names[2], names[3], names[4]
	byCharge := map[string]string{"charge": "$args.charge"}
	p := profile.Profile{V: profile.Version, Vendor: vendor, Transport: "mcp", Pin: profile.Pin{Kind: "none"},
		Source: "catalog", CatalogEntry: StripeRefundsEntry,
		Operations: []profile.Operation{
			{Tool: details, Title: "Read a Stripe operation's fields", Kind: "read"},
			{Tool: account, Title: "Look up the Stripe account", Kind: "read"},
			{Tool: charge, Title: "Look up a charge", Kind: "read", Resource: &profile.Resource{Arg: "charge", Noun: "charge"}},
			{Tool: refunds, Title: "List a charge's refunds", Kind: "read", Resource: &profile.Resource{Arg: "charge", Noun: "charge"}},
			{Tool: refund, Title: "Refund a charge", Kind: "write", Effect: profile.EffectMovesMoney, Shape: "single",
				Resource: &profile.Resource{Arg: "charge", Noun: "charge"},
				Amount:   &profile.Amount{From: "argument", Arg: "amount", Format: "integer_minor", Scale: 2, Unit: stripeRefundUnit},
				// planque writes the marker into Stripe's arguments, inside the
				// operation's parameters; the refund Stripe makes carries it
				// at metadata.kiff_op, where the refunds read finds it.
				Marker:            &profile.Marker{Field: narrowedParams + "." + stripeRefundsMarker},
				RequiredArguments: []string{"charge", "amount"},
				Case: &profile.Case{Kind: profile.CaseKindRefund, ChargeArg: "charge",
					AccountRead: profile.Read{Tool: account, Path: "id"},
					ChargeRead: profile.CaseChargeRead{Tool: charge, Args: byCharge,
						Currency: "currency", Captured: "amount_captured", Disputed: "disputed"},
					RefundsRead: profile.CaseRefundsRead{Tool: refunds, Args: byCharge,
						List: "data", ID: "id", Amount: "amount", Status: "status", Marker: stripeRefundsMarker}}},
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
