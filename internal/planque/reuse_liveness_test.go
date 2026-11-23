package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/store"
)

// liveProvider answers each read with what refuse says, and records what it was asked.
type liveProvider struct {
	connector.Provider
	t      *testing.T
	tools  map[string]connector.Tool
	user   string
	status int   // 0 answers; else refused with this status
	err    error // returned instead of an answer: a slow or unreachable vendor
	calls  []string
}

func (p *liveProvider) Account(_ context.Context, id string) (connector.Account, error) {
	return connector.Account{ID: id, User: p.user, Toolkit: "googlecalendar", Status: "ACTIVE"}, nil
}

func (p *liveProvider) Platform() connector.Platform {
	return connector.Platform{Name: "Test", URL: "https://provider.example"}
}

func (p *liveProvider) Tools(context.Context, string, string) (connector.Page[connector.Tool], error) {
	var items []connector.Tool
	for _, t := range p.tools {
		items = append(items, t)
	}
	return connector.Page[connector.Tool]{Items: items}, nil
}

func (p *liveProvider) Execute(_ context.Context, _ connector.Identity, _ string, t connector.Tool, args json.RawMessage, _ func() bool) (connector.Result, error) {
	p.calls = append(p.calls, t.Slug)
	if p.err != nil {
		return connector.Result{}, p.err
	}
	if string(args) != `{}` {
		p.t.Fatalf("the check asked %s with %s", t.Slug, args)
	}
	if p.status != 0 {
		return connector.Result{Refused: true, Status: p.status, Reason: "vendor's words"}, nil
	}
	return connector.Result{Data: json.RawMessage(`{"data":{"items":[]}}`)}, nil
}

// A managed sign-in to reuse is read once at its vendor first (21 Nov: Google Calendar's sign-in read "Connected", was
// reused, and was refused at its first call), with the app's own "who am I" only. A refused sign-in (401) isn't reused,
// and every connection on it reads signed out; a live one, or another refusal, goes on as before. A sign-in that doesn't
// cover "who am I" (Calendar on events alone) is never read at all: never the owner's content, even through a read its
// connection has (a review on #330).
func TestAReusedSignInIsReadAtItsVendorFirst(t *testing.T) {
	const events, calendar = "https://www.googleapis.com/auth/calendar.events", "https://www.googleapis.com/auth/calendar"
	list := connector.Tool{Slug: "GOOGLECALENDAR_EVENTS_LIST", Toolkit: "googlecalendar", Version: "20300101_00", Name: "List events", ReadOnly: true,
		InputSchema: json.RawMessage(`{"type":"object","properties":{"calendarId":{"type":"string"}}}`)}
	who := connector.Tool{Slug: "GOOGLECALENDAR_GET_CURRENT_USER", Toolkit: "googlecalendar", Version: "20300101_00", Name: "Who", ReadOnly: true,
		InputSchema: json.RawMessage(`{"type":"object"}`), Scopes: []string{calendar, "https://www.googleapis.com/auth/calendar.readonly"}}
	for _, c := range []struct {
		name      string
		status    int
		scope     string
		want      error
		signedOut bool
		asked     []string
	}{
		{"refused (401): not reused, signed out", http.StatusUnauthorized, calendar, errSignInRefused, true, []string{who.Slug}},
		{"live: reused", 0, calendar, nil, false, []string{who.Slug}},
		{"another refusal (403): reused as before", http.StatusForbidden, calendar, nil, false, []string{who.Slug}},
		{"events alone doesn't cover who am I: nothing read, not even its events", http.StatusUnauthorized, events, nil, false, nil},
	} {
		e := newInProcessEnv(t)
		p := &liveProvider{t: t, tools: map[string]connector.Tool{who.Slug: who}, status: c.status}
		e.s.Connectors = map[string]connector.Provider{"composio": p}
		g := whoGrant(e, "googlecalendar", []string{c.scope}, []string{c.scope})
		g.ID, g.Status, g.SignedIn = "grt_calendar", StatusConnected, true
		def := list
		g.Tools = append(g.Tools, GrantedTool{ToolRequest: ToolRequest{Tool: def.Alias(), Name: def.Slug, ReadOnly: true}, Connector: &def})
		if err := e.s.saveGrant(context.Background(), g); err != nil {
			t.Fatal(err)
		}
		if err := e.s.reuseSignInWorks(context.Background(), g); !errors.Is(err, c.want) || err != nil && c.want == nil {
			t.Errorf("%s: %v, want %v", c.name, err, c.want)
		}
		if len(p.calls) != len(c.asked) || len(c.asked) == 1 && p.calls[0] != c.asked[0] {
			t.Errorf("%s: asked %v, want %v", c.name, p.calls, c.asked)
		}
		after, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
		if (after.Status == StatusNeedsReconnect) != c.signedOut {
			t.Errorf("%s: the connection reads %s", c.name, after.Status)
		}
	}
}

// A review on #330: the check is wired into the reuse request, an error or a timeout signs no one out, and a 401 signs
// out every connection on that sign-in, the one that reused it too.
func TestAReuseRequestOfARefusedSignInIsRefusedAndSignsItOut(t *testing.T) {
	const calendar = "https://www.googleapis.com/auth/calendar"
	who := connector.Tool{Slug: "GOOGLECALENDAR_GET_CURRENT_USER", Toolkit: "googlecalendar", Version: "20300101_00", Name: "Who", ReadOnly: true,
		InputSchema: json.RawMessage(`{"type":"object"}`), Scopes: []string{calendar}}
	setup := func(p *liveProvider) (*env, *Grant, *Grant) {
		e := newInProcessEnv(t)
		e.s.Connectors = map[string]connector.Provider{"composio": p}
		e.s.ManagedTestTenants = map[string]bool{tenant: true}
		e.enrol()
		saved := func(id string) *Grant {
			g := whoGrant(e, "googlecalendar", []string{calendar}, []string{calendar})
			g.ID, g.Status, g.SignedIn = id, StatusConnected, true
			g.Request.API = &APIRequest{}
			if err := e.s.saveGrant(context.Background(), g); err != nil {
				t.Fatal(err)
			}
			return g
		}
		first, second := saved("grt_first"), saved("grt_second")
		p.user = first.Connector.User
		return e, first, second
	}
	reuse := func(e *env, from *Grant) (int, map[string]any) {
		return e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"reuse": from.ID, "agent": "agt_taf", "agent_name": "Today's agenda", "tools": []any{}}, nil, "k_setup")
	}
	status := func(e *env, id string) string {
		g, _ := e.s.loadGrant(context.Background(), tenant, id)
		return g.Status
	}

	// A 401: 409 sign_in_refused, nothing reused, and both connections on that sign-in read signed out.
	p := &liveProvider{t: t, tools: map[string]connector.Tool{who.Slug: who}, status: http.StatusUnauthorized}
	e, first, second := setup(p)
	before := len(e.kvGrants(t))
	st, out := reuse(e, first)
	if st != http.StatusConflict || out["code"] != "sign_in_refused" {
		t.Fatalf("a refused sign-in's reuse: %d %v", st, out)
	}
	if n := len(e.kvGrants(t)); n != before {
		t.Fatalf("a grant was made: %d → %d", before, n)
	}
	if status(e, first.ID) != StatusNeedsReconnect || status(e, second.ID) != StatusNeedsReconnect {
		t.Fatalf("the sign-in's connections: %s, %s", status(e, first.ID), status(e, second.ID))
	}

	// The vendor erring or timing out isn't a refusal: the reuse goes on, and no one is signed out.
	for _, err := range []error{errors.New("provider unreachable"), context.DeadlineExceeded} {
		p := &liveProvider{t: t, tools: map[string]connector.Tool{who.Slug: who}, err: err}
		e, first, second := setup(p)
		if st, out := reuse(e, first); st != http.StatusCreated {
			t.Fatalf("%v: the reuse: %d %v", err, st, out)
		}
		if status(e, first.ID) != StatusConnected || status(e, second.ID) != StatusConnected {
			t.Fatalf("%v: signed out: %s, %s", err, status(e, first.ID), status(e, second.ID))
		}
	}
}

// kvGrants are the account's grants as kept.
func (e *env) kvGrants(t *testing.T) []store.Item {
	t.Helper()
	items, err := e.s.KV.Query(context.Background(), tenantPK(tenant), "GRANT#")
	if err != nil {
		t.Fatal(err)
	}
	return items
}
