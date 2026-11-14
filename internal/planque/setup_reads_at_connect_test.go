package planque

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/words"
)

// askWithSetupReads is a connect-first request to the fake vendor whose
// setup reads are named at the request, as a catalog entry's are: the
// request is stored as the entry route makes it, and its digest is what
// the owner signs.
func (e *env) askWithSetupReads(reads ...string) (id, sum string) {
	e.t.Helper()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		e.t.Fatalf("request: %d %v", st, out)
	}
	ctx := context.Background()
	g, _ := e.s.loadGrant(ctx, tenant, out["id"].(string))
	g.Request.SetupReads = reads
	g.Text, g.Texts = texts(func(w words.Approval) string { return requestText(w, g.Request, false, "") })
	var err error
	if g.Digest, err = digest(g.Request); err != nil {
		e.t.Fatal(err)
	}
	if err := e.s.saveGrant(ctx, g); err != nil {
		e.t.Fatal(err)
	}
	return g.ID, g.Digest
}

// A catalog entry's first request names its setup reads: the owner approves
// them with the sign-in, and once signed in they're connected, Setup's own
// Card asked of KIFF, with no second approval. A narrowed write whose
// fields they read can then be added at once, with the taf's Card.
func TestSetupReadsNamedAtTheRequestConnectWithTheSignIn(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	if st, out := e.enrolAs(phone); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	e.enrolled, e.phone = true, phone
	id, digest := e.askWithSetupReads("list_charges")
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if !strings.Contains(g.Text, "Setup reads, to set this up: list_charges.") {
		t.Fatalf("the owner's text doesn't name the read:\n%s", g.Text)
	}
	st, out := e.approveAndSignIn(id, digest)
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	if len(g.Tools) != 1 || !g.Tools[0].SetupReads || !g.Tools[0].ReadOnly || g.Tools[0].Tool != "list_charges" {
		t.Fatalf("the grant's tools: %+v", g.Tools)
	}
	if _, ok := e.kiff.Connection(tenant, g.Tools[0].Name); !ok {
		t.Fatal("the read isn't in KIFF")
	}
	if n := len(e.kiff.SetupReads); n != 1 || e.kiff.SetupReads[0].Grant != id {
		t.Fatalf("Setup's Card wasn't asked of KIFF: %+v", e.kiff.SetupReads)
	}
	// What it has now counts as approved: the next tools are an addition.
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	list := map[string]any{"tool": "create_refund", "title": "Refund a charge", "name": "create_refund", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2}
	if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup"); st != 200 || out["declared"] != true ||
		out["status_words"] != "Connected · new tools waiting for your approval" {
		t.Fatalf("the taf's tools after the reads: %d %v", st, out)
	}
}

// A read the vendor doesn't list, or doesn't mark read-only, connects
// nothing: the grant is listed as any connect-first one, and the reads are
// asked again with the tools. The owner approved a read, never a write.
func TestSetupReadsThatArentReadsConnectNothing(t *testing.T) {
	for _, name := range []string{"create_refund", "no_such_tool"} {
		t.Run(name, func(t *testing.T) {
			e := newEnv(t)
			e.enrol()
			id, digest := e.askWithSetupReads(name)
			st, out := e.approveAndSignIn(id, digest)
			if st != 200 || out["status"] != StatusListed {
				t.Fatalf("confirm: %d %v", st, out)
			}
			if c := e.kiff.Connections(tenant); len(c) != 0 {
				t.Fatalf("in KIFF: %v", c)
			}
			g, _ := e.s.loadGrant(context.Background(), tenant, id)
			if len(g.Tools) != 0 || len(g.Request.Tools) != 0 {
				t.Fatalf("the grant: tools %v, asked %v", g.Tools, g.Request.Tools)
			}
		})
	}
}

// A bare server's first request names its setup reads too (the Eval's setup
// connects by server): the owner approves them with the sign-in, and once
// signed in they're checked against what the server lists and connected,
// with no second approval. A name it doesn't list connects none, and the
// grant's events say why. A request with its tools at once, a key or a
// reused sign-in can't name them.
func TestABareServerTakesItsSetupReadsAtConnect(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	grants := "/v1/accounts/" + tenant + "/grants"
	for name, body := range map[string]map[string]any{
		"with tools at once": {"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{map[string]any{"tool": "list_charges", "title": "List charges", "read_only": true}}, "setup_reads": []string{"list_charges"}},
		"a reused sign-in":   {"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}, "reuse": "grt_x", "setup_reads": []string{"list_charges"}},
	} {
		st, out := e.do("POST", grants, body, nil, "k_setup")
		if st != 400 || !strings.Contains(fmt.Sprint(out["error"]), "setup_reads goes on a server's first request") {
			t.Fatalf("%s: %d %v", name, st, out)
		}
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("in KIFF: %v", c)
	}
	// Named at connect: approved with the sign-in, connected once signed in.
	st, out := e.do("POST", grants, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}, "setup_reads": []string{"stripe_api_details", "list_charges"}}, nil, "k_setup")
	if st != 201 || !strings.Contains(fmt.Sprint(out["text"]), "stripe_api_details") {
		t.Fatalf("request: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, out := e.approveAndSignIn(id, out["digest"].(string)); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if len(g.Tools) != 2 || !g.Tools[0].SetupReads || !g.Tools[1].SetupReads {
		t.Fatalf("the grant's tools: %+v", g.Tools)
	}
	// A name the server doesn't list: nothing connects, and the events say why.
	st, out = e.do("POST", grants, map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf2", "tools": []any{}, "setup_reads": []string{"list_everything"}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("request: %d %v", st, out)
	}
	id = out["id"].(string)
	if st, out := e.approveAndSignIn(id, out["digest"].(string)); st != 200 || out["status"] != StatusListed {
		t.Fatalf("an unlisted read: %d %v", st, out)
	}
	g, _ = e.s.loadGrant(context.Background(), tenant, id)
	said := false
	for _, ev := range g.Events {
		said = said || strings.Contains(ev.What, "lists no tool named list_everything")
	}
	if len(g.Tools) != 0 || !said {
		t.Fatalf("an unlisted read: tools %v, events %v", g.Tools, g.Events)
	}
}

// A catalog entry's request takes setup_reads, each named once, at most ten.
func TestAnEntrysSetupReadsAreNamedOnce(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	entry := "/v1/accounts/" + tenant + "/catalog/" + StripeEntry + "/grants"
	for name, reads := range map[string][]string{
		"none":   {},
		"twice":  {"stripe_api_details", "stripe_api_details"},
		"blank":  {" "},
		"eleven": {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
	} {
		if st, out := e.do("POST", entry, map[string]any{"agent": "agt_taf", "setup_reads": reads}, nil, "k_setup"); st != 400 {
			t.Fatalf("%s: %d %v", name, st, out)
		}
	}
	st, out := e.do("POST", entry, map[string]any{"agent": "agt_taf", "setup_reads": []string{"stripe_api_details"}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("Stripe with its setup read: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
	if len(g.Request.SetupReads) != 1 || !strings.Contains(g.Text, "Setup reads, to set this up: stripe_api_details.") {
		t.Fatalf("the request: %v\n%s", g.Request.SetupReads, g.Text)
	}
}

// Stripe's own server, connected first with no reads named, gets its
// describing tool as a setup read with the sign-in: it can't be named before
// the catalog is listed, and a narrowed write needs it first (the Eval's red 1).
// Bound to Stripe's server: another server that lists a tool by that name gets
// nothing, and a request that names its reads keeps them.
func TestStripesServerGetsItsDescriberWithTheSignIn(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	entry := "/v1/accounts/" + tenant + "/catalog/" + StripeEntry + "/grants"
	st, out := e.do("POST", entry, map[string]any{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("Stripe, no reads named: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, out["grant"].(map[string]any)["id"].(string))
	if !slices.Equal(g.Request.SetupReads, []string{"stripe_api_details"}) || !strings.Contains(g.Text, "Setup reads, to set this up: stripe_api_details.") {
		t.Fatalf("the request: %v\n%s", g.Request.SetupReads, g.Text)
	}
	if serverDescribers("https://mcp.stripe.com.example.net") != nil || serverDescribers("https://evil.example/mcp") != nil {
		t.Fatal("a server that isn't Stripe's got a describer")
	}
	// Another server that lists stripe_api_details: nothing added.
	e2 := newEnv(t)
	e2.enrol()
	st, out = e2.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e2.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("another server: %d %v", st, out)
	}
	g2, _ := e2.s.loadGrant(context.Background(), tenant, out["id"].(string))
	if len(g2.Request.SetupReads) != 0 || strings.Contains(g2.Text, "stripe_api_details") {
		t.Fatalf("another server got a read: %v\n%s", g2.Request.SetupReads, g2.Text)
	}
}
