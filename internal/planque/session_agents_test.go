package planque

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/kiff"
)

// A request names the chat that asked for it (session): planque keeps it
// beside the grant, never in its request, so no request served and no digest
// signed carries it; a chat's delete reads which agents it asked for; and an
// ended grant keeps no chat.
func TestASessionsAgentsAreKeptBesideTheGrantAndEndWithIt(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	// The person who deletes the chat, with their own session (agent_end.go).
	e.kiff.Keys["s_owner"] = kiff.Me{TenantID: tenant, Capabilities: []string{"read", "tools.connect", "cards.manage"}}
	ask := func(agent, session string) (int, map[string]any) {
		return e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": agent, "session": session,
			"tools": []map[string]any{{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
	}
	st, out := ask("agt_chat", "ses_chat01")
	if st != 201 {
		t.Fatalf("request with a session: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, _ := ask("agt_other", "ses_other02"); st != 201 {
		t.Fatalf("request for another chat: %d", st)
	}
	if st, _ := ask("agt_none", ""); st != 201 {
		t.Fatalf("request with no session: %d", st)
	}

	// Kept beside the request, never in it: the answer, the request and the
	// digest's input don't hold it.
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil || g.AskedIn != "ses_chat01" {
		t.Fatalf("asked in: %q %v", g.AskedIn, err)
	}
	served, _ := json.Marshal(out)
	request, _ := json.Marshal(g.Request)
	if strings.Contains(string(served), "ses_chat01") || strings.Contains(string(request), "ses_chat01") {
		t.Fatalf("the session reaches the request served or signed:\n%s\n%s", served, request)
	}

	// The chat's agents: only that chat's, with no grant id, token or tool.
	st, got := e.do("GET", "/v1/accounts/"+tenant+"/sessions/ses_chat01/agents", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("session agents: %d %v", st, got)
	}
	agents, _ := json.Marshal(got["agents"])
	if string(agents) != `[{"agent":"agt_chat","since":"`+g.Request.CreatedAt[:10]+`","status":"pending"}]` {
		t.Fatalf("session agents: %s", agents)
	}
	if strings.Contains(string(agents), id) {
		t.Fatalf("the answer names the grant: %s", agents)
	}
	if st, got := e.do("GET", "/v1/accounts/"+tenant+"/sessions/ses_nobody/agents", nil, nil, "k_setup"); st != 200 || len(got["agents"].([]any)) != 0 {
		t.Fatalf("a chat that asked nothing: %d %v", st, got)
	}

	// Ended, it keeps no chat, and the chat's delete finds nothing left.
	if st, out := e.endAgent(tenant, "agt_chat", "s_owner", ""); st != 200 {
		t.Fatalf("end: %d %v", st, out)
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Status != StatusDisconnected || g.AskedIn != "" {
		t.Fatalf("ended: %s, asked in %q", g.Status, g.AskedIn)
	}
	if _, got := e.do("GET", "/v1/accounts/"+tenant+"/sessions/ses_chat01/agents", nil, nil, "k_setup"); len(got["agents"].([]any)) != 0 {
		t.Fatalf("an ended grant still listed: %v", got)
	}
}

// A session is an id or nothing: anything else is refused before a grant is
// made, on the route too.
func TestASessionThatIsntAnIDIsRefused(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	for _, bad := range []string{"a b", "ses/../x", strings.Repeat("s", 129), "ses\n1", "é"} {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_x", "session": bad,
			"tools": []map[string]any{{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}}}, nil, "k_setup")
		if st != 400 || !strings.Contains(out["error"].(string), "session is the chat's id") {
			t.Errorf("%q: %d %v", bad, st, out)
		}
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/sessions/a.b/agents", nil, nil, "k_setup"); st != 400 {
		t.Errorf("route with a session that isn't an id: %d", st)
	}
	if !askedInOK("") || !askedInOK("ses_01-AB") || askedInOK("a.b") {
		t.Error("askedInOK")
	}
}
