package planque

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
)

// whoProvider answers the apps' "who am I" reads, and records what it was
// asked: only the tools in answers exist, and a call to one not in scopes
// below fails the test, since planque must never ask for more than the
// sign-in has.
type whoProvider struct {
	connector.Provider
	t        *testing.T
	tools    map[string]connector.Tool
	answers  map[string]string
	refuse   bool
	failWith error
	calls    []string
}

func (p *whoProvider) Platform() connector.Platform {
	return connector.Platform{Name: "Test", URL: "https://provider.example"}
}

func (p *whoProvider) Tools(context.Context, string, string) (connector.Page[connector.Tool], error) {
	var items []connector.Tool
	for _, t := range p.tools {
		items = append(items, t)
	}
	return connector.Page[connector.Tool]{Items: items}, nil
}

func (p *whoProvider) Execute(_ context.Context, _ connector.Identity, _ string, t connector.Tool, args json.RawMessage, live func() bool) (connector.Result, error) {
	p.calls = append(p.calls, t.Slug)
	if p.failWith != nil {
		return connector.Result{}, p.failWith
	}
	if string(args) != `{}` || !live() {
		p.t.Fatalf("who-am-I asked with %s", args)
	}
	return connector.Result{Data: json.RawMessage(p.answers[t.Slug]), Refused: p.refuse}, nil
}

func whoGrant(e *env, toolkit string, scopes, granted []string) *Grant {
	g := &Grant{Tenant: tenant,
		Request:   Request{Server: "https://provider.example", Connector: &ConnectorRequest{Provider: "composio", Toolkit: toolkit, Scopes: scopes}},
		Connector: &ConnectorState{Account: "account_1", Granted: granted}}
	g.Connector.User = e.s.connectorIdentity(g).User
	return g
}

// The app's own "who am I" names the account the owner signed in with, for
// each app in the table, on the narrowest sign-in planque asks for as well as
// the wide one; a sign-in that doesn't cover it is never asked, and a refusal,
// an error or an answer without a name leaves the name empty (kiff/planque
// account name after sign-in).
func TestTheAccountIsNamedByTheAppsOwnWhoAmI(t *testing.T) {
	const g = "https://www.googleapis.com/auth/"
	tool := func(slug, toolkit string, scopes ...string) connector.Tool {
		return connector.Tool{Slug: slug, Toolkit: toolkit, Version: "20300101_00", Name: slug, InputSchema: json.RawMessage(`{"type":"object"}`), ReadOnly: true, Scopes: scopes}
	}
	for _, c := range []struct {
		name, toolkit string
		tool          connector.Tool
		answer        string
		asked, got    []string // the sign-in's scopes: asked, and granted when the provider says
		want          string
		called        bool
	}{
		{"LinkedIn, the usual sign-in", "linkedin", tool("LINKEDIN_GET_MY_INFO", "linkedin", "profile", "r_liteprofile"),
			`{"data":{"response_data":{"localizedFirstName":"Carolina","localizedLastName":"Arias","id":"x1"}}}`, []string{"openid", "profile", "w_member_social"}, nil, "Carolina Arias", true},
		{"LinkedIn, an email wins", "linkedin", tool("LINKEDIN_GET_MY_INFO", "linkedin", "profile"),
			`{"data":{"email":"carolina@kiff.dev","name":"Carolina Arias"}}`, []string{"profile"}, nil, "carolina@kiff.dev", true},
		{"Drive, the narrowest it asks (file-only) covers about", "googledrive", tool("GOOGLEDRIVE_GET_ABOUT", "googledrive", g+"drive", g+"drive.file", g+"drive.metadata.readonly"),
			`{"data":{"user":{"displayName":"Carolina Arias","emailAddress":"carolina@kiff.dev"}}}`, []string{g + "drive.file"}, nil, "carolina@kiff.dev", true},
		{"Drive, a sign-in granted less than asked is not asked", "googledrive", tool("GOOGLEDRIVE_GET_ABOUT", "googledrive", g+"drive", g+"drive.readonly"),
			`{"data":{"user":{"emailAddress":"carolina@kiff.dev"}}}`, []string{g + "drive.file"}, []string{g + "docs"}, "", false},
		{"Calendar, events alone doesn't cover the current user", "googlecalendar", tool("GOOGLECALENDAR_GET_CURRENT_USER", "googlecalendar", g+"calendar", g+"calendar.readonly"),
			`{"data":{"email":"carolina@kiff.dev"}}`, []string{g + "calendar.events"}, nil, "", false},
		{"Calendar, the whole scope covers it", "googlecalendar", tool("GOOGLECALENDAR_GET_CURRENT_USER", "googlecalendar", g+"calendar", g+"calendar.readonly"),
			`{"data":{"email":"carolina@kiff.dev"}}`, []string{g + "calendar"}, nil, "carolina@kiff.dev", true},
		{"Calendar, an opaque id is never taken for the account", "googlecalendar", tool("GOOGLECALENDAR_GET_CURRENT_USER", "googlecalendar", g+"calendar", g+"calendar.readonly"),
			`{"data":{"id":"x9f2opaque"}}`, []string{g + "calendar"}, nil, "", true},
		{"Outlook, User.Read", "outlook", tool("OUTLOOK_GET_PROFILE", "outlook", "User.Read", "User.ReadBasic.All"),
			`{"data":{"mail":null,"userPrincipalName":"carolona@kiff.dev","displayName":"Carolina Arias Moreno"}}`, []string{"Mail.Read", "User.Read"}, nil, "carolona@kiff.dev", true},
		{"Outlook, mail only doesn't cover the profile", "outlook", tool("OUTLOOK_GET_PROFILE", "outlook", "User.Read", "User.ReadBasic.All"),
			`{"data":{"mail":"x@y.z"}}`, []string{"Mail.Read"}, nil, "", false},
		{"GitHub", "github", tool("GITHUB_GET_THE_AUTHENTICATED_USER", "github", "user"),
			`{"data":{"login":"carolina-test","id":7}}`, []string{"user"}, nil, "carolina-test", true},
		{"Typeform, accounts:read", "typeform", tool("TYPEFORM_GET_ABOUT_ME", "typeform", "accounts:read"),
			`{"data":{"email":"carolina@kiff.dev","alias":"carolina"}}`, []string{"accounts:read", "forms:read"}, nil, "carolina@kiff.dev", true},
		{"Typeform, without accounts:read", "typeform", tool("TYPEFORM_GET_ABOUT_ME", "typeform", "accounts:read"),
			`{"data":{"email":"carolina@kiff.dev"}}`, []string{"forms:read"}, nil, "", false},
		{"an answer without a name", "github", tool("GITHUB_GET_THE_AUTHENTICATED_USER", "github", "user"),
			`{"data":{"id":7}}`, []string{"user"}, nil, "", true},
		{"a name that is not one line is made one", "github", tool("GITHUB_GET_THE_AUTHENTICATED_USER", "github", "user"),
			`{"data":{"login":"carol\n  ina\u0000"}}`, []string{"user"}, nil, "carol ina", true},
		{"an app not in the table", "notion", tool("NOTION_GET_ABOUT_ME", "notion", "read"),
			`{"data":{"name":"Someone"}}`, nil, nil, "", false},
	} {
		e := newInProcessEnv(t)
		p := &whoProvider{t: t, tools: map[string]connector.Tool{c.tool.Slug: c.tool}, answers: map[string]string{c.tool.Slug: c.answer}}
		e.s.Connectors = map[string]connector.Provider{"composio": p}
		gr := whoGrant(e, c.toolkit, c.asked, c.got)
		e.s.nameAccount(context.Background(), gr)
		if gr.Connector.AccountName != c.want || (len(p.calls) > 0) != c.called {
			t.Errorf("%s: named %q (want %q), asked %v (want asked %v)", c.name, gr.Connector.AccountName, c.want, p.calls, c.called)
		}
	}

	// A refused answer, a failed call, and a name the provider already showed.
	tl := connector.Tool{Slug: "GITHUB_GET_THE_AUTHENTICATED_USER", Toolkit: "github", Version: "20300101_00", Name: "me", InputSchema: json.RawMessage(`{"type":"object"}`), ReadOnly: true, Scopes: []string{"user"}}
	for name, p := range map[string]*whoProvider{
		"refused":     {refuse: true, answers: map[string]string{tl.Slug: `{"data":{"login":"leaked"},"error":"nope"}`}},
		"errored":     {failWith: errors.New("down"), answers: map[string]string{tl.Slug: `{"data":{"login":"leaked"}}`}},
		"not a read":  {answers: map[string]string{tl.Slug: `{"data":{"login":"leaked"}}`}},
		"not defined": {answers: map[string]string{tl.Slug: `{"data":{"login":"leaked"}}`}},
	} {
		e := newInProcessEnv(t)
		p.t = t
		p.tools = map[string]connector.Tool{tl.Slug: tl}
		switch name {
		case "not a read":
			w := tl
			w.ReadOnly = false
			p.tools = map[string]connector.Tool{tl.Slug: w}
		case "not defined":
			p.tools = nil
		}
		e.s.Connectors = map[string]connector.Provider{"composio": p}
		gr := whoGrant(e, "github", []string{"user"}, nil)
		e.s.nameAccount(context.Background(), gr)
		if gr.Connector.AccountName != "" || (name == "not a read" || name == "not defined") && len(p.calls) != 0 {
			t.Errorf("%s: named %q, asked %v", name, gr.Connector.AccountName, p.calls)
		}
	}
	e := newInProcessEnv(t)
	p := &whoProvider{t: t, tools: map[string]connector.Tool{tl.Slug: tl}, answers: map[string]string{tl.Slug: `{"data":{"login":"other"}}`}}
	e.s.Connectors = map[string]connector.Provider{"composio": p}
	gr := whoGrant(e, "github", []string{"user"}, nil)
	gr.Connector.AccountName = "Demo Shop"
	e.s.nameAccount(context.Background(), gr)
	if gr.Connector.AccountName != "Demo Shop" || len(p.calls) != 0 {
		t.Errorf("a name the provider showed was asked again or replaced: %q, %v", gr.Connector.AccountName, p.calls)
	}
}

// The name is read after the owner's fingerprint and sign-in, by planque and
// not the taf: the sign-in's confirm names the account when the provider shows
// none, nothing is asked before it, and no tool of the taf runs.
func TestTheSignInIsNamedAtConfirmAndNotBefore(t *testing.T) {
	whoAmI["desk_demo"] = struct {
		Tools  []string
		Fields []string
	}{[]string{"DESK_DEMO_WHOAMI"}, []string{"login"}}
	t.Cleanup(func() { delete(whoAmI, "desk_demo") })
	e, f := newConnectorFixture(t)
	f.noName, f.who = true, true
	v := requestConnector(t, e)
	id := v["id"].(string)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
		t.Fatalf("approve sign-in: %d %v", st, out)
	}
	code := e.signIn(f.callback + "&connected_account_id=untrusted&status=ACTIVE")
	if f.whoCalls != 0 {
		t.Fatalf("the account was read before the sign-in was confirmed: %d", f.whoCalls)
	}
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	gr, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if gr.Connector.AccountName != "carolina-test" || f.whoCalls != 1 || f.executions != 0 {
		t.Fatalf("named %q after %d who-am-I reads and %d of the taf's calls", gr.Connector.AccountName, f.whoCalls, f.executions)
	}
}

// A who-am-I that never answers is cut at its own short deadline and reads as
// any failure: no name, and the sign-in's confirm is not held for it
// (mac-local's review of the account-name change).
func TestAWhoAmIThatNeverAnswersIsCutShort(t *testing.T) {
	old := whoAmITimeout
	whoAmITimeout = 50 * time.Millisecond
	t.Cleanup(func() { whoAmITimeout = old })
	tl := connector.Tool{Slug: "GITHUB_GET_THE_AUTHENTICATED_USER", Toolkit: "github", Version: "20300101_00", Name: "me", InputSchema: json.RawMessage(`{"type":"object"}`), ReadOnly: true, Scopes: []string{"user"}}
	e := newInProcessEnv(t)
	p := &slowWho{whoProvider: whoProvider{t: t, tools: map[string]connector.Tool{tl.Slug: tl}}}
	e.s.Connectors = map[string]connector.Provider{"composio": p}
	gr := whoGrant(e, "github", []string{"user"}, nil)
	start := time.Now()
	e.s.nameAccount(context.Background(), gr)
	if gr.Connector.AccountName != "" || time.Since(start) > 2*time.Second || !p.cut {
		t.Fatalf("named %q after %v, cut %v", gr.Connector.AccountName, time.Since(start), p.cut)
	}
}

// slowWho answers its read only when its context ends.
type slowWho struct {
	whoProvider
	cut bool
}

func (p *slowWho) Execute(ctx context.Context, _ connector.Identity, _ string, _ connector.Tool, _ json.RawMessage, _ func() bool) (connector.Result, error) {
	<-ctx.Done()
	p.cut = true
	return connector.Result{}, ctx.Err()
}

// Names come from vendor answers now, so a sign-in again to the same account
// is the same account whatever the case it comes back in.
// Category: authority
// Attack: a sign-in again refused for a difference of case only, so the right account is turned away (and the owner told it is another one)
// Refused: nothing refused: it reconnects
func TestASignInAgainToTheSameAccountInAnotherCaseReconnects(t *testing.T) {
	e, f := newConnectorFixture(t)
	ctx := context.Background()
	f.distinct = true
	first := connectConnector(t, e, f)
	signOut401(t, e, f, first)
	g, _ := e.s.loadGrant(ctx, tenant, first.ID)
	path := "/v1/accounts/" + tenant + "/grants/" + first.ID
	e.at(time.Minute)
	e.phone.Now = func() time.Time { return time.Now().Add(time.Minute) }
	if st, out := e.do("POST", path+"/approve", reapprove(g), e.phone, ""); st != 200 {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	f.accountName = "DEMO SHOP"
	code := e.signIn(f.callback + "&status=ACTIVE")
	if st, out := e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm the same account in another case: %d %v", st, out)
	}
}
