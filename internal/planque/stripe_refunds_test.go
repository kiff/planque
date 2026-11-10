package planque

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/profile"
)

const stripeTestAccount = "acct_test1"

// refundAsk is the one tool a caller asks for refunds, in a mode.
func refundAsk(live bool) []ToolRequest {
	mode, _ := json.Marshal(live)
	op, _ := json.Marshal(stripeRefund)
	return []ToolRequest{{Tool: stripeWrite, Name: "stripe_refund", Fixed: map[string]json.RawMessage{stripeOperation: op, stripeLivemode: mode}}}
}

// The reviewed request is the five tools, the refund last, each narrowed
// tool on the grant's account and the asked mode, and planque's catalog
// profile for them: the refund moves money in euros, its marker inside the
// operation's parameters, and every case read is one of the three granted.
func TestStripeRefundsProfileIsReviewedAndComplete(t *testing.T) {
	tools, raw, err := stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, nil, "grt_abcd", "Stripe")
	if err != nil {
		t.Fatal(err)
	}
	want := []struct{ tool, op string }{{stripeDetails, ""}, {stripeRead, stripeGetAccount}, {stripeRead, stripeGetCharge}, {stripeRead, stripeListRefunds}, {stripeWrite, stripeRefund}}
	if len(tools) != len(want) {
		t.Fatalf("tools: %+v", tools)
	}
	for i, w := range want {
		if tools[i].Tool != w.tool || stripeOperationOf(tools[i]) != w.op {
			t.Fatalf("tool %d is %s %s, want %s %s", i, tools[i].Tool, stripeOperationOf(tools[i]), w.tool, w.op)
		}
		if w.op == "" {
			if !tools[i].SetupReads || len(tools[i].Fixed) > 0 {
				t.Fatalf("the describing read: %+v", tools[i])
			}
			continue
		}
		// Every narrowed tool is on the grant's account, in the asked mode.
		if string(tools[i].Fixed[stripeContext]) != `"`+stripeTestAccount+`"` || string(tools[i].Fixed[stripeLivemode]) != "false" {
			t.Fatalf("tool %d's account and mode: %s", i, tools[i].Fixed)
		}
		if msg := checkNarrowing(tools[i]); msg != "" {
			t.Fatalf("tool %d: %s", i, msg)
		}
	}
	if string(tools[3].Fixed["parameters.limit"]) != "100" || string(tools[1].Fixed["parameters.id"]) != `"`+stripeTestAccount+`"` {
		t.Fatalf("the reads' fixed fields: %s %s", tools[3].Fixed, tools[1].Fixed)
	}
	if tools[4].Name != "stripe_refund" || tools[4].Arguments["charge"].Path != "parameters.charge" || tools[4].Arguments["amount"].Path != "parameters.amount" {
		t.Fatalf("the refund: %+v", tools[4])
	}
	p, err := profile.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "catalog" || p.CatalogEntry != StripeRefundsEntry || p.Transport != "mcp" || p.TestProven {
		t.Fatalf("profile header: %+v", p)
	}
	refund := p.Operations[4]
	c := refund.Case
	if refund.Effect != profile.EffectMovesMoney || refund.Amount.Unit != "EUR" || refund.Marker == nil || refund.Marker.Field != "parameters.metadata.kiff_op" || c == nil {
		t.Fatalf("refund operation: %+v", refund)
	}
	names := toolNames(tools)
	if c.AccountRead.Tool != names[1] || c.AccountRead.Path != "id" || c.ChargeRead.Tool != names[2] || c.RefundsRead.Tool != names[3] ||
		c.ChargeRead.Args["charge"] != "$args.charge" || c.RefundsRead.Args["charge"] != "$args.charge" ||
		c.RefundsRead.Marker != "metadata.kiff_op" || c.RefundsRead.List != "data" || c.RefundEvent != nil {
		t.Fatalf("case block: %+v", c)
	}
	// It passes as a reviewed profile, and never as a drafted one.
	if _, err := deriveProfileFields(raw, tools, true); err != nil {
		t.Fatalf("reviewed: %v", err)
	}
	if _, err := deriveProfileFields(raw, tools, false); err == nil {
		t.Fatal("planque's reviewed profile passed as a caller's draft")
	}
	// Names planque allocates follow into every case read.
	named := make([]ToolRequest, len(tools))
	for i := range tools {
		named[i] = tools[i]
		named[i].Name = toolKey(tools[i].Tool, "grt_abcd") + string(rune('a'+i))
	}
	renamed, err := renameProfile(raw, names, named)
	if err != nil {
		t.Fatal(err)
	}
	q, _ := profile.Parse(renamed)
	if rc := q.Operations[4].Case; rc.AccountRead.Tool != named[1].Name || rc.ChargeRead.Tool != named[2].Name || rc.RefundsRead.Tool != named[3].Name {
		t.Fatalf("case reads after renaming: %+v", rc)
	}
	// Live mode isn't offered yet: Stripe may ask for its own confirmation
	// of each refund there, which hasn't been checked.
	if _, _, err := stripeRefundsAsk(refundAsk(true), nil, stripeTestAccount, nil, "grt_abcd", "Stripe"); err == nil || !strings.Contains(err.Error(), "test mode only") {
		t.Fatalf("live: %v", err)
	}
}

// One tool exactly, with only its operation and mode, no profile of the
// caller's, and an account from Stripe's sign-in.
func TestStripeRefundsAskTakesOneToolAndNoProfile(t *testing.T) {
	withFixed := func(k, v string) []ToolRequest {
		ts := refundAsk(false)
		ts[0].Fixed[k] = json.RawMessage(v)
		return ts
	}
	withArgs := refundAsk(false)
	withArgs[0].Arguments = map[string]Argument{"charge": {Path: "parameters.charge", Type: "string"}}
	two := append(refundAsk(false), ToolRequest{Tool: stripeRead})
	noMode := refundAsk(false)
	delete(noMode[0].Fixed, stripeLivemode)
	noMode[0].Fixed["parameters.amount"] = json.RawMessage(`1`)
	for name, tools := range map[string][]ToolRequest{
		"another tool":             two,
		"the caller's account":     withFixed(stripeContext, `"acct_other"`),
		"the caller's arguments":   withArgs,
		"a mode that isn't a bool": withFixed(stripeLivemode, `"false"`)[:1],
		"no mode":                  noMode,
	} {
		if name == "a mode that isn't a bool" {
			tools[0].Fixed = map[string]json.RawMessage{stripeOperation: tools[0].Fixed[stripeOperation], stripeLivemode: json.RawMessage(`"false"`)}
		}
		if _, _, err := stripeRefundsAsk(tools, nil, stripeTestAccount, nil, "grt_abcd", "Stripe"); err == nil || !strings.Contains(err.Error(), "asks one tool") {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, _, err := stripeRefundsAsk(refundAsk(false), json.RawMessage(`{"v":"kiff.profile.v1"}`), stripeTestAccount, nil, "grt_abcd", "Stripe"); err == nil || !strings.Contains(err.Error(), "without a profile") {
		t.Errorf("a caller's profile: %v", err)
	}
	if _, _, err := stripeRefundsAsk(refundAsk(false), nil, "", nil, "grt_abcd", "Stripe"); err == nil || !strings.Contains(err.Error(), "sign in to Stripe again") {
		t.Errorf("no account from the sign-in: %v", err)
	}
}

// A drafted profile never carries a case block: a refund's money rules
// are reviewed code.
func TestADraftedProfileWithACaseBlockIsRefused(t *testing.T) {
	tools := []ToolRequest{{Tool: "get_account"}, {Tool: "get_charge"}, {Tool: "list_refunds"}, {Tool: "refund"}}
	_, err := deriveProfileFields(json.RawMessage(refundCaseProfile), tools, false)
	if err == nil || !strings.Contains(err.Error(), "case block") {
		t.Fatalf("drafted case: %v", err)
	}
}

// The case's marker may sit inside the operation's parameters only on a
// tool narrowed to one write, and only under exactly "parameters.": any
// other prefix is refused by the profile's own rules, and a tool that isn't
// narrowed is refused by planque.
func TestTheCaseMarkerInsideParametersIsOnlyForANarrowedWrite(t *testing.T) {
	tools, raw, err := stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, nil, "grt_abcd", "Stripe")
	if err != nil {
		t.Fatal(err)
	}
	// The same profile on the refund unnarrowed: refused.
	plain := append([]ToolRequest(nil), tools...)
	plain[4] = ToolRequest{Tool: stripeWrite, Name: tools[4].Name}
	if _, err := deriveProfileFields(raw, plain, true); err == nil || !strings.Contains(err.Error(), "narrowed to one write") {
		t.Fatalf("an unnarrowed tool marker inside parameters: %v", err)
	}
	// Any other prefix: the profile itself is refused.
	other := strings.Replace(string(raw), `"field":"parameters.metadata.kiff_op"`, `"field":"params.metadata.kiff_op"`, 1)
	if other == string(raw) {
		t.Fatal("the marker field wasn't found in the profile")
	}
	if _, err := profile.Parse(json.RawMessage(other)); err == nil || !strings.Contains(err.Error(), "refunds_read.marker") {
		t.Fatalf("another prefix: %v", err)
	}
	inner := strings.Replace(string(raw), `"field":"parameters.metadata.kiff_op"`, `"field":"parameters.parameters.metadata.kiff_op"`, 1)
	if _, err := profile.Parse(json.RawMessage(inner)); err == nil {
		t.Fatal("a doubled prefix passed")
	}
}

// The marker planque writes into a narrowed write's fields is checked
// against the operation's own description, like its arguments: an
// operation that lists no metadata refuses it.
func TestANarrowedWritesMarkerIsOneOfItsOperationsFields(t *testing.T) {
	tools, _, err := stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, nil, "grt_abcd", "Stripe")
	if err != nil {
		t.Fatal(err)
	}
	d, op, write, ok := narrowedTo(tools[4])
	if !ok || !write {
		t.Fatal("the refund isn't a narrowed write")
	}
	with := parseFields("## Body Parameters\n- **amount** [integer] The amount.\n- **charge** [string] The charge.\n- **metadata** [object] Set of key-value pairs.\n")
	without := parseFields("## Body Parameters\n- **amount** [integer] The amount.\n- **charge** [string] The charge.\n")
	marker := "parameters.metadata.kiff_op"
	if msg := fieldsCheck(tools[4], d, op, with, marker); msg != "" {
		t.Fatalf("with metadata: %s", msg)
	}
	if msg := fieldsCheck(tools[4], d, op, without, marker); !strings.Contains(msg, "metadata.kiff_op") {
		t.Fatalf("without metadata: %q", msg)
	}
	if msg := fieldsCheck(tools[4], d, op, without, ""); msg != "" {
		t.Fatalf("no marker: %s", msg)
	}
	// The broad tool's schema doesn't list the operation's fields: its
	// marker is left to that check, for a narrowed write only.
	schema := json.RawMessage(`{"type":"object","properties":{"stripe_api_operation_id":{"type":"string"},"parameters":{"type":"object"}}}`)
	g := &Grant{Request: Request{Tools: tools}}
	_, g.Request.Profile, _ = stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, nil, "grt_abcd", "Stripe")
	if err := profileVendorSchema(g, tools[4].Name, schema); err != nil {
		t.Fatalf("a narrowed write's marker against the broad schema: %v", err)
	}
	g.Request.Tools[4] = ToolRequest{Tool: stripeWrite, Name: tools[4].Name}
	if err := profileVendorSchema(g, tools[4].Name, schema); err == nil {
		t.Fatal("an unnarrowed tool's marker passed the schema it isn't in")
	}
}

// stripeDirectGrant is a listed direct Stripe connection, signed in, whose
// sign-in named its account, with Stripe's three broad tools listed.
func stripeDirectGrant(t *testing.T, e *env, account string) string {
	t.Helper()
	g := &Grant{ID: "grt_stripe1", Tenant: tenant, Status: StatusListed, SignedIn: true, VendorName: "Stripe", Account: account,
		Request: Request{Version: "planque.grant.v1", Grant: "grt_stripe1", Tenant: tenant, Server: stripeServer, Agent: "agt_taf",
			AskedBy: "agt_setup", Tools: []ToolRequest{}, CatalogEntry: StripeEntry},
		Catalog: []CatalogTool{
			{Name: stripeDetails, ReadOnly: true, InputSchema: json.RawMessage(`{"type":"object"}`)},
			{Name: stripeRead, ReadOnly: true, InputSchema: json.RawMessage(`{"type":"object"}`)},
			{Name: stripeWrite, InputSchema: json.RawMessage(`{"type":"object"}`)},
		}}
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	return "/v1/accounts/" + tenant + "/grants/" + g.ID
}

// On a direct Stripe connection, asking for the refund takes planque's
// whole request: one tool asked, five granted, planque's reviewed profile.
// Off, the refund is any narrowed write, which needs the describing tool
// approved first.
func TestAStripeRefundConnectionTakesPlanquesReviewedRequest(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	path := stripeDirectGrant(t, e, stripeTestAccount)
	ask := func(body map[string]any) (int, map[string]any) {
		return e.do("POST", path+"/tools", body, nil, "k_setup")
	}

	// Off: a narrowed write as any other, which names its arguments and
	// needs its describing tool approved first: this one is refused.
	if st, out := ask(map[string]any{"tools": refundAsk(false)}); st != 400 || out["status"] != nil {
		t.Fatalf("off: %d %v", st, out)
	}

	e.s.StripeRefunds = true
	if st, out := ask(map[string]any{"tools": append(refundAsk(false), ToolRequest{Tool: stripeRead})}); st != 400 || !strings.Contains(out["error"].(string), "asks one tool") {
		t.Fatalf("another tool: %d %v", st, out)
	}
	drafted := json.RawMessage(`{"v":"kiff.profile.v1","vendor":"Stripe","transport":"mcp","pin":{"kind":"none"},"source":"drafted","operations":[]}`)
	if st, out := ask(map[string]any{"tools": refundAsk(false), "profile": drafted}); st != 400 || !strings.Contains(out["error"].(string), "without a profile") {
		t.Fatalf("a caller's profile: %d %v", st, out)
	}
	st, out := ask(map[string]any{"tools": refundAsk(false)})
	if st != 200 || out["status"] != StatusPending {
		t.Fatalf("the refund: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, "grt_stripe1")
	if err != nil {
		t.Fatal(err)
	}
	p, err := profile.Parse(g.Request.Profile)
	if err != nil || p.Source != "catalog" || p.CatalogEntry != StripeRefundsEntry || len(p.Operations) != 5 || p.Operations[4].Case == nil {
		t.Fatalf("the grant's profile: %v %+v", err, p)
	}
	// Every case read names a tool the grant holds, under the name planque
	// gave it, and every narrowed tool is on the account Stripe named.
	held := map[string]ToolRequest{}
	for _, tr := range g.Request.Tools {
		held[tr.Name] = tr
	}
	c := p.Operations[4].Case
	for _, n := range []string{c.AccountRead.Tool, c.ChargeRead.Tool, c.RefundsRead.Tool, p.Operations[4].Tool} {
		tr, ok := held[n]
		if !ok || string(tr.Fixed[stripeContext]) != `"`+stripeTestAccount+`"` {
			t.Fatalf("%s isn't among the grant's tools on the account: %v", n, held)
		}
	}
	if len(g.Request.Tools) != 5 || !g.Request.Tools[0].SetupReads {
		t.Fatalf("the grant's tools: %+v", g.Request.Tools)
	}
	// Once it has tools, refunds aren't added to it.
	if st, out := ask(map[string]any{"tools": refundAsk(false)}); st != 409 {
		t.Fatalf("asked again: %d %v", st, out)
	}

	// A connection whose sign-in didn't name the account: sign in again.
	e2 := newInProcessEnv(t)
	e2.enrol()
	e2.s.StripeRefunds = true
	path2 := stripeDirectGrant(t, e2, "")
	if st, out := e2.do("POST", path2+"/tools", map[string]any{"tools": refundAsk(false)}, nil, "k_setup"); st != 409 || !strings.Contains(out["error"].(string), "sign in to Stripe again") {
		t.Fatalf("no account: %d %v", st, out)
	}
}

// fakeCatalogProvider lists the managed connector's Stripe and one other
// app, in its catalog and in a tool search.
type fakeCatalogProvider struct{ connector.Provider }

func (fakeCatalogProvider) Platform() connector.Platform {
	return connector.Platform{Name: "Provider", URL: "https://connect.example"}
}
func (fakeCatalogProvider) Catalog(context.Context, string, string) (connector.Page[connector.App], error) {
	return connector.Page[connector.App]{Items: []connector.App{
		{Slug: stripeToolkit, Name: "Stripe", Setup: connector.SetupReady},
		{Slug: "desk_demo", Name: "Desk", Setup: connector.SetupReady}}}, nil
}
func (fakeCatalogProvider) Search(context.Context, string, []string, string) (connector.Search, error) {
	return connector.Search{
		Tools: []connector.Tool{{Slug: stripeRefund, Toolkit: stripeToolkit, Name: "Refund"}, {Slug: "DESK_LIST", Toolkit: "desk_demo", Name: "List"}},
		Apps: map[string]connector.App{stripeToolkit: {Slug: stripeToolkit, Name: "Stripe", Setup: connector.SetupReady},
			"desk_demo": {Slug: "desk_demo", Name: "Desk", Setup: connector.SetupReady}}}, nil
}

// In the app, Stripe is one entry: planque's direct sign-in to Stripe's MCP
// server, the one with test mode. The managed connector's Stripe is neither
// listed, nor found in a tool search, nor connected.
func TestTheAppShowsOneStripe(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	e.s.Connectors = map[string]connector.Provider{"fake": fakeCatalogProvider{}}
	base := "/v1/accounts/" + tenant + "/catalog"

	st, out := e.do("GET", base, nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	stripes := []string{}
	for _, v := range out["entries"].([]any) {
		if v := v.(map[string]any); v["title"] == "Stripe" {
			stripes = append(stripes, v["entry"].(string))
		}
	}
	if len(stripes) != 1 || stripes[0] != StripeEntry {
		t.Fatalf("Stripe entries: %v", stripes)
	}

	st, out = e.do("GET", base+"/tools?q=refund", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("search: %d %v", st, out)
	}
	for _, v := range out["tools"].([]any) {
		if v.(map[string]any)["app"] == stripeToolkit {
			t.Fatalf("the managed Stripe found: %v", v)
		}
	}

	if st, out := e.do("POST", base+"/connector:fake:"+stripeToolkit+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 409 {
		t.Fatalf("the managed Stripe asked: %d %v", st, out)
	}
	st, out = e.do("POST", base+"/"+StripeEntry+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("Stripe asked: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
	if err != nil || g.Request.Server != stripeServer || g.Request.Connector != nil || g.Request.CatalogEntry != StripeEntry ||
		g.Vendor.AccountField != stripeAccountField {
		t.Fatalf("the grant: %v %+v", err, g.Request)
	}

	// A managed Stripe connection made before stays listed, connected, so
	// the owner can disconnect it.
	old := &Grant{ID: "grt_oldstripe", Tenant: tenant, Status: StatusConnected, VendorName: "Stripe",
		Request: Request{Version: "planque.grant.v1", Grant: "grt_oldstripe", Tenant: tenant, Agent: "agt_taf", Kind: "connector", Tools: []ToolRequest{},
			Connector: &ConnectorRequest{Provider: "fake", Toolkit: stripeToolkit, Title: "Stripe"}}}
	if err := e.s.saveGrant(context.Background(), old); err != nil {
		t.Fatal(err)
	}
	st, out = e.do("GET", base, nil, nil, "k_setup")
	listed := false
	for _, v := range out["entries"].([]any) {
		if v := v.(map[string]any); v["entry"] == "connector:fake:"+stripeToolkit {
			listed = v["state"] == "connected" && v["grant"] == old.ID
		}
	}
	if st != 200 || !listed {
		t.Fatalf("the old managed Stripe connection: %d %v", st, out["entries"])
	}
}

// Stripe's accounts tool is read from its structuredContent, or else from
// the same JSON as text; anything else names no account.
func TestTheSessionsAccountsAreReadAsStripeAnswers(t *testing.T) {
	one := `{"accounts":[{"stripe_context":"acct_1","livemode":false,"name":"Shop"}]}`
	for name, c := range map[string]struct {
		res  string
		want []sessionAccount
		ok   bool
	}{
		"structured":             {`{"content":[{"type":"text","text":"x"}],"structuredContent":` + one + `}`, []sessionAccount{{"acct_1", false}}, true},
		"text only":              {`{"content":[{"type":"text","text":` + strconvQuote(one) + `}]}`, []sessionAccount{{"acct_1", false}}, true},
		"two":                    {`{"structuredContent":{"accounts":[{"stripe_context":"acct_1","livemode":false},{"stripe_context":"acct_2","livemode":true}]}}`, []sessionAccount{{"acct_1", false}, {"acct_2", true}}, true},
		"none":                   {`{"structuredContent":{"accounts":[]}}`, []sessionAccount{}, true},
		"an error":               {`{"isError":true,"structuredContent":` + one + `}`, nil, false},
		"another shape":          {`{"structuredContent":{"data":[]}}`, nil, false},
		"not an account":         {`{"structuredContent":{"accounts":[{"stripe_context":"org_1","livemode":false}]}}`, nil, false},
		"a sentence":             {`{"structuredContent":{"accounts":[{"stripe_context":"acct_1 or another","livemode":false}]}}`, nil, false},
		"no mode":                {`{"structuredContent":{"accounts":[{"stripe_context":"acct_1"}]}}`, nil, false},
		"not even an MCP result": {`[1]`, nil, false},
	} {
		got, ok := sessionAccounts(json.RawMessage(c.res))
		if ok != c.ok || (ok && !slices.Equal(got, c.want)) {
			t.Errorf("%s: %v %v", name, got, ok)
		}
	}
}

func strconvQuote(s string) string { b, _ := json.Marshal(s); return string(b) }

// A Stripe sign-in whose token answer names no account: right after it,
// planque reads the account with Stripe's own tool, itself, and keeps it
// with its mode, as the grant's account. A refund connection then pins it;
// one whose account is in live mode is refused.
func TestASignInWithNoAccountKeepsTheOneItsAccountsToolLists(t *testing.T) {
	e := newEnv(t)
	e.vendor.UnnamedAccount = true
	e.enrol()
	id, digest := e.request()
	ctx := context.Background()
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	g.Vendor.AccountsTool = stripeAccountsTool
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	if st, out := e.approveAndSignIn(id, digest); st != 200 {
		t.Fatalf("sign in: %d %v", st, out)
	}
	if g, err = e.s.loadGrant(ctx, tenant, id); err != nil {
		t.Fatal(err)
	}
	if g.Account != "acct_fake123" || g.AccountLive == nil || *g.AccountLive {
		t.Fatalf("the grant's account: %q %v", g.Account, g.AccountLive)
	}
	if !strings.Contains(string(mustJSON(g.Events)), "read which account the sign-in covers: one, kept") {
		t.Fatalf("the read isn't on the record: %s", mustJSON(g.Events))
	}
	// A sign-in that covers two accounts keeps neither: the owner signs in
	// to one.
	e3 := newEnv(t)
	e3.vendor.UnnamedAccount, e3.vendor.MoreAccounts = true, []string{"acct_other456"}
	e3.enrol()
	id3, digest3 := e3.request()
	g3, _ := e3.s.loadGrant(ctx, tenant, id3)
	g3.Vendor.AccountsTool = stripeAccountsTool
	if err := e3.s.saveGrant(ctx, g3); err != nil {
		t.Fatal(err)
	}
	if st, out := e3.approveAndSignIn(id3, digest3); st != 200 {
		t.Fatalf("sign in: %d %v", st, out)
	}
	if g3, _ = e3.s.loadGrant(ctx, tenant, id3); g3.Account != "" || !strings.Contains(string(mustJSON(g3.Events)), "2 accounts, so none is kept") {
		t.Fatalf("two accounts: %q %s", g3.Account, mustJSON(g3.Events))
	}
	// Without the tool set, nothing is read and no account is kept.
	e2 := newEnv(t)
	e2.vendor.UnnamedAccount = true
	e2.enrol()
	id2, digest2 := e2.request()
	if st, out := e2.approveAndSignIn(id2, digest2); st != 200 {
		t.Fatalf("sign in: %d %v", st, out)
	}
	if g2, _ := e2.s.loadGrant(ctx, tenant, id2); g2.Account != "" || g2.AccountLive != nil {
		t.Fatalf("an account kept with no tool to read it: %q", g2.Account)
	}
	// The account's mode must be the one asked.
	live := true
	if _, _, err := stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, &live, "grt_abcd", "Stripe"); err == nil || !strings.Contains(err.Error(), "account is in live mode") {
		t.Fatalf("a live account for a test-mode refund: %v", err)
	}
	test := false
	if _, _, err := stripeRefundsAsk(refundAsk(false), nil, stripeTestAccount, &test, "grt_abcd", "Stripe"); err != nil {
		t.Fatalf("a test account: %v", err)
	}
}
