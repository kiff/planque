package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

// liveProvider answers each read with what refuse says, and records what it was asked.
type liveProvider struct {
	connector.Provider
	t      *testing.T
	tools  map[string]connector.Tool
	status int // 0 answers; else refused with this status
	calls  []string
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
// connection has (dev-android on #330).
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
