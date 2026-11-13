package planque

import (
	"context"
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

// Negative 5: a bare server's catalog is unknown before its sign-in, so its
// first request can't name setup reads: it connects first, then asks for
// its reads with its tools, two steps as before. A request with tools at
// once can't either.
func TestABareServerKeepsItsTwoSteps(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	grants := "/v1/accounts/" + tenant + "/grants"
	for name, body := range map[string]map[string]any{
		"a bare server":      {"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}, "setup_reads": []string{"list_charges"}},
		"with tools at once": {"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{map[string]any{"tool": "list_charges", "title": "List charges", "read_only": true}}, "setup_reads": []string{"list_charges"}},
	} {
		st, out := e.do("POST", grants, body, nil, "k_setup")
		if st != 400 || !strings.Contains(out["error"].(string), "setup_reads goes on a catalog entry's first request") {
			t.Fatalf("%s: %d %v", name, st, out)
		}
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("in KIFF: %v", c)
	}
	// Connected first, it's listed: its reads come with its tools, a second approval.
	id := e.listedGrant()
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Status != StatusListed || len(g.Tools) != 0 {
		t.Fatalf("the bare server's first step: %s %v", g.Status, g.Tools)
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
