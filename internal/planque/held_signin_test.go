package planque

import (
	"context"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// A new taf asks for an app the account already holds a sign-in to: that
// sign-in is reused, with no browser and no new sign-in, and the new
// connection gets only the tools the owner approves for the new taf, never
// the other connection's. Its tools are declared while they wait, so the
// tools and the taf's Card take one approval.
// Category: authority
// Attack: the reused connection carrying the held one's tools, a sign-in narrowed to other tools reused
// Refused: each; the new connection holds its own tools only, and a sign-in that doesn't cover the ask isn't reused
func TestANewTafReusesTheHeldSignInOfAnAppByItsEntry(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	first := connectConnector(t, e, f)
	held, _ := e.kiff.Connection(tenant, "reply_ticket")
	authorizations, links := f.authorizations, f.links
	entry := "/v1/accounts/" + tenant + "/catalog/connector:composio:desk_demo/grants"

	st, out := e.do("POST", entry, map[string]any{"agent": "agt_refunds", "agent_name": "Refund requests"}, nil, "k_setup")
	if st != 201 || out["needs_key"] != false {
		t.Fatalf("request: %d %v", st, out)
	}
	v := out["grant"].(map[string]any)
	if v["status"] != StatusListed || v["signed_in"] != true || v["sign_in_url"] != nil || f.authorizations != authorizations || f.links != links {
		t.Fatalf("the held sign-in wasn't reused: %v; sign-ins %d→%d", v, authorizations, f.authorizations)
	}
	id := v["id"].(string)
	g, _ := e.s.loadGrant(ctx, tenant, id)
	if g.Request.Reuse != first.ID || g.Connector.Account != first.Connector.Account || len(g.Tools) != 0 || len(g.Request.Tools) != 0 {
		t.Fatalf("the reused connection: reuse %q, account %q, tools %v, asked %v", g.Request.Reuse, g.Connector.Account, g.Tools, g.Request.Tools)
	}

	// Its own tool, declared while it waits: one approval for it and the Card.
	path := "/v1/accounts/" + tenant + "/grants/" + id
	_, catalog := e.do("GET", path+"/catalog", nil, nil, "k_setup")
	alias := catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	p := profile.Profile{V: profile.Version, Vendor: "Support desk", Transport: "rest", Pin: profile.Pin{Kind: "none"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "refund_reply", Title: "Reply to a ticket", Kind: "write", Effect: profile.EffectSends, Shape: "single",
			Resource: &profile.Resource{Arg: "ticket", Noun: "ticket", Outside: true}, RequiredArguments: []string{"ticket", "text"}}}}
	st, out = e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Refund reply", Name: "refund_reply"}}, "profile": p}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending || out["declared"] != true {
		t.Fatalf("tools: %d %v", st, out)
	}
	if _, ok := e.kiff.Connection(tenant, "refund_reply"); !ok {
		t.Fatal("refund_reply isn't declared while it waits")
	}
	// Declined: it leaves KIFF, and the connection is listed again on the held sign-in.
	if st, out := e.do("DELETE", path+"/tools", nil, nil, "k_setup"); st != 200 || out["status"] != StatusListed || out["signed_in"] != true {
		t.Fatalf("decline: %d %v", st, out)
	}
	if _, ok := e.kiff.Connection(tenant, "refund_reply"); ok {
		t.Fatal("refund_reply is still in KIFF after the decline")
	}
	_, catalog = e.do("GET", path+"/catalog", nil, nil, "k_setup")
	alias = catalog["tools"].([]any)[0].(map[string]any)["name"].(string)
	st, out = e.do("POST", path+"/tools", map[string]any{"tools": []ToolRequest{{Tool: alias, Title: "Refund reply", Name: "refund_reply"}}, "profile": p}, nil, "k_setup")
	if st != 200 || out["status"] != StatusPending || out["declared"] != true {
		t.Fatalf("tools again: %d %v", st, out)
	}
	declared, ok := e.kiff.Connection(tenant, "refund_reply")
	if !ok {
		t.Fatal("refund_reply isn't declared again")
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	ph, _ := profile.Digest(g.Request.Profile)
	st, out = e.do("POST", path+"/approve", map[string]any{"digest": g.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || f.authorizations != authorizations || f.links != links {
		t.Fatalf("approve: %d %v; sign-ins %d→%d", st, out, authorizations, f.authorizations)
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	if len(g.Tools) != 1 || g.Tools[0].Name != "refund_reply" {
		t.Fatalf("the reused connection's tools: %v", toolNamesOf(g.Tools))
	}
	if c, ok := e.kiff.Connection(tenant, "refund_reply"); !ok || c.Credential != declared.Credential {
		t.Fatal("refund_reply was registered again, or went, at the approval")
	}
	if c, ok := e.kiff.Connection(tenant, "reply_ticket"); !ok || c.Credential != held.Credential {
		t.Fatal("the held connection's tool changed")
	}

	// Asked for a sign-in of its own: one is started, as before.
	st, out = e.do("POST", entry, map[string]any{"agent": "agt_third", "new_sign_in": true}, nil, "k_setup")
	if st != 201 || out["grant"].(map[string]any)["status"] != StatusPending || out["grant"].(map[string]any)["signed_in"] == true {
		t.Fatalf("a sign-in of its own: %d %v", st, out)
	}

	// A held sign-in narrowed to other tools doesn't cover the ask: not reused.
	for _, gid := range []string{first.ID, id} {
		h, _ := e.s.loadGrant(ctx, tenant, gid)
		h.Request.Connector.AskedTools, h.Request.Connector.UsualAccess = []string{"DESK_REPLY_ONLY"}, false
		if err := e.s.saveGrant(ctx, h); err != nil {
			t.Fatal(err)
		}
	}
	st, out = e.do("POST", entry, map[string]any{"agent": "agt_fourth"}, nil, "k_setup")
	if st != 201 || out["grant"].(map[string]any)["status"] != StatusPending || out["grant"].(map[string]any)["signed_in"] == true {
		t.Fatalf("a narrowed sign-in reused for an ask it may not cover: %d %v", st, out)
	}
}
