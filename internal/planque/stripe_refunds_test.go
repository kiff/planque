package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/store"
)

func stripeAsk(slugs ...string) []ToolRequest {
	out := make([]ToolRequest, 0, len(slugs))
	for _, s := range slugs {
		out = append(out, ToolRequest{Tool: stripeAlias(s)})
	}
	return out
}

// The reviewed profile is planque's catalog profile for exactly the four
// tools: the refund moves money with its marker and case block, and every
// case read is one of the three reads granted with it.
func TestStripeRefundsProfileIsReviewedAndComplete(t *testing.T) {
	// Asked in any order, they come back in the fixed one.
	tools, raw, err := stripeRefundsAsk(stripeAsk(stripeRefund, stripeListRefunds, stripeGetAccount, stripeGetCharge), nil, "Stripe")
	if err != nil {
		t.Fatal(err)
	}
	for i, slug := range stripeRefundTools {
		if tools[i].Tool != stripeAlias(slug) {
			t.Fatalf("tool %d is %s, want %s", i, tools[i].Tool, slug)
		}
	}
	p, err := profile.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "catalog" || p.CatalogEntry != StripeRefundsEntry || p.Transport != "rest" || p.TestProven {
		t.Fatalf("profile header: %+v", p)
	}
	refund := p.Operations[3]
	c := refund.Case
	if refund.Effect != profile.EffectMovesMoney || refund.Marker == nil || refund.Marker.Field != "metadata.kiff_op" || c == nil {
		t.Fatalf("refund operation: %+v", refund)
	}
	names := toolNames(tools)
	if c.AccountRead.Tool != names[0] || c.ChargeRead.Tool != names[1] || c.RefundsRead.Tool != names[2] ||
		c.ChargeRead.Args["charge_id"] != "$args.charge" || c.RefundsRead.Args["charge"] != "$args.charge" ||
		c.RefundsRead.Marker != refund.Marker.Field || c.RefundEvent == nil || c.RefundEvent.Event != stripeRefundCreated {
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
		named[i].Name = toolKey(stripeRefundTools[i], "grt_abcd")
	}
	renamed, err := renameProfile(raw, names, named)
	if err != nil {
		t.Fatal(err)
	}
	q, _ := profile.Parse(renamed)
	if rc := q.Operations[3].Case; rc.AccountRead.Tool != named[0].Name || rc.ChargeRead.Tool != named[1].Name || rc.RefundsRead.Tool != named[2].Name ||
		q.Operations[3].Amount.Unit != "EUR" {
		t.Fatalf("case reads after renaming: %+v", rc)
	}
}

// The four tools exactly, and no profile of the caller's.
func TestStripeRefundsAskTakesTheFourToolsAndNoProfile(t *testing.T) {
	for name, tools := range map[string][]ToolRequest{
		"a read missing": stripeAsk(stripeRefund, stripeGetAccount, stripeGetCharge),
		"another tool":   stripeAsk(stripeRefund, stripeGetAccount, stripeGetCharge, stripeListRefunds, "STRIPE_LIST_CUSTOMERS"),
		"one twice":      stripeAsk(stripeRefund, stripeGetAccount, stripeGetCharge, stripeGetCharge),
	} {
		if _, _, err := stripeRefundsAsk(tools, nil, "Stripe"); err == nil || !strings.Contains(err.Error(), "exactly four tools") {
			t.Errorf("%s: %v", name, err)
		}
	}
	all := stripeAsk(stripeRefund, stripeGetAccount, stripeGetCharge, stripeListRefunds)
	if _, _, err := stripeRefundsAsk(all, json.RawMessage(`{"v":"kiff.profile.v1"}`), "Stripe"); err == nil || !strings.Contains(err.Error(), "without a profile") {
		t.Errorf("a caller's profile: %v", err)
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

// fakeStripeProvider is a managed connector whose one account is a Stripe
// sign-in. Only what setting tools reads is answered.
type fakeStripeProvider struct {
	connector.Provider
	identity connector.Identity
}

func (f fakeStripeProvider) Platform() connector.Platform {
	return connector.Platform{Name: "Provider", URL: "https://connect.example"}
}
func (f fakeStripeProvider) Account(_ context.Context, id string) (connector.Account, error) {
	return connector.Account{ID: id, User: f.identity.User, Toolkit: stripeToolkit, Status: "ACTIVE"}, nil
}

// stripeGrant is a listed managed Stripe connection, signed in, whose
// catalog read left the four tools and one other.
func stripeGrant(t *testing.T, e *env) string {
	t.Helper()
	ctx := context.Background()
	g := &Grant{ID: "grt_stripe1", Tenant: tenant, Status: StatusListed, SignedIn: true, VendorName: "Stripe",
		Request: Request{Version: "planque.grant.v1", Grant: "grt_stripe1", Tenant: tenant, Server: "https://connect.example", Agent: "agt_taf",
			AskedBy: "agt_setup", Kind: "connector", Tools: []ToolRequest{},
			Connector: &ConnectorRequest{Provider: "fake", Toolkit: stripeToolkit, Title: "Stripe", AuthConfig: "ac_stripe"},
			API:       &APIRequest{URL: "https://connect.example", Base: "https://connect.example", SourceKind: "connector"}}}
	id := e.s.connectorIdentity(g)
	g.Connector = &ConnectorState{User: id.User, Account: "acct_stripe"}
	e.s.Connectors = map[string]connector.Provider{"fake": fakeStripeProvider{identity: id}}
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	for _, slug := range append(append([]string{}, stripeRefundTools...), "STRIPE_LIST_CUSTOMERS") {
		def := connector.Tool{Slug: slug, Toolkit: stripeToolkit, Version: "20261101_00", Name: slug, ReadOnly: slug != stripeRefund,
			InputSchema: json.RawMessage(`{"type":"object","properties":{}}`)}
		if err := putDoc(ctx, e.s.KV, store.Item{PK: tenantPK(tenant), SK: connectorToolKey(g, def.Alias()), TTL: time.Now().Add(time.Hour).Unix()}, def, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	return "/v1/accounts/" + tenant + "/grants/" + g.ID
}

// On a managed Stripe connection, the refund tool takes planque's reviewed
// profile: the four tools exactly, no profile of the caller's. Off, the
// refund tool is any tool, and takes a drafted profile with no case block.
func TestAStripeRefundConnectionTakesPlanquesReviewedProfile(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	path := stripeGrant(t, e)
	four := stripeAsk(stripeRefund, stripeGetAccount, stripeGetCharge, stripeListRefunds)

	// Off: the request needs the caller's profile, as for any connector tool.
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": four}, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "owner-approved profile") {
		t.Fatalf("off: %d %v", st, out)
	}

	e.s.StripeRefunds = true
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": four[:3]}, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "exactly four tools") {
		t.Fatalf("a read missing: %d %v", st, out)
	}
	drafted := json.RawMessage(`{"v":"kiff.profile.v1","vendor":"Stripe","transport":"rest","pin":{"kind":"none"},"source":"drafted","operations":[]}`)
	if st, out := e.do("POST", path+"/tools", map[string]any{"tools": four, "profile": drafted}, nil, "k_setup"); st != 400 || !strings.Contains(out["error"].(string), "without a profile") {
		t.Fatalf("a caller's profile: %d %v", st, out)
	}
	st, out := e.do("POST", path+"/tools", map[string]any{"tools": four}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending {
		t.Fatalf("the four tools: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, "grt_stripe1")
	if err != nil {
		t.Fatal(err)
	}
	p, err := profile.Parse(g.Request.Profile)
	if err != nil || p.Source != "catalog" || p.CatalogEntry != StripeRefundsEntry || len(p.Operations) != 4 || p.Operations[3].Case == nil {
		t.Fatalf("the grant's profile: %v %+v", err, p)
	}
	// Every case read names a tool the grant holds, under the name planque gave it.
	held := map[string]bool{}
	for _, tr := range g.Request.Tools {
		held[tr.Name] = true
	}
	c := p.Operations[3].Case
	for _, n := range []string{c.AccountRead.Tool, c.ChargeRead.Tool, c.RefundsRead.Tool, p.Operations[3].Tool} {
		if !held[n] {
			t.Fatalf("%s isn't among the grant's tools %v", n, held)
		}
	}
	if len(g.Request.Connector.Tools) != 4 {
		t.Fatalf("connector definitions bound: %d", len(g.Request.Connector.Tools))
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
	if err != nil || g.Request.Server != stripeServer || g.Request.Connector != nil || g.Request.CatalogEntry != StripeEntry {
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
