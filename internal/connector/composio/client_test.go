package composio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
)

type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
func testClient(f transport) *Client {
	return &Client{Secret: func(context.Context) (string, error) { return "synthetic-project-key", nil }, HTTP: &http.Client{Transport: f}}
}

const toolJSON = `{"slug":"DESK_REPLY","name":"Reply to a ticket","toolkit":{"slug":"desk"},"version":"20300101_00","input_parameters":{"type":"object","properties":{"ticket":{"type":"string"},"text":{"type":"string"}},"required":["ticket","text"]},"tags":[],"scopes":["tickets:reply"]}`

func approvedTool(t *testing.T) connector.Tool {
	t.Helper()
	var w wireTool
	if err := json.Unmarshal([]byte(toolJSON), &w); err != nil {
		t.Fatal(err)
	}
	tl, err := w.tool()
	if err != nil {
		t.Fatal(err)
	}
	return tl
}

const accountJSON = `{"id":"account_one","user_id":"owner_one","status":"ACTIVE","toolkit":{"slug":"desk"}}`

func TestConcreteExecutionPinsAccountVersionAndArguments(t *testing.T) {
	calls := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("x-api-key") != "synthetic-project-key" {
			t.Error("project credential missing")
		}
		switch r.URL.Path {
		case "/api/v3.1/connected_accounts/account_one":
			return response(200, accountJSON), nil
		case "/api/v3.1/tools/DESK_REPLY":
			if r.URL.Query().Get("version") != "20300101_00" {
				t.Error("tool version not pinned")
			}
			return response(200, toolJSON), nil
		case "/api/v3.1/tools/execute/DESK_REPLY":
			calls++
			var body map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if len(body) != 4 || string(body["user_id"]) != `"owner_one"` || string(body["connected_account_id"]) != `"account_one"` || string(body["version"]) != `"20300101_00"` {
				t.Fatalf("execution binding: %s", mustJSON(body))
			}
			var args map[string]any
			_ = json.Unmarshal(body["arguments"], &args)
			if args["ticket"] != "ticket_one" || args["text"] != "hello" {
				t.Fatal("arguments changed")
			}
			return response(200, `{"data":{"id":"reply_one"},"successful":true,"error":null}`), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
			return nil, nil
		}
	})
	got, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one", approvedTool(t), json.RawMessage(`{"ticket":"ticket_one","text":"hello"}`), func() bool { return true })
	if err != nil || got.Refused || got.Unknown || !strings.Contains(string(got.Data), "reply_one") || calls != 1 {
		t.Fatalf("result %+v err=%v sends=%d", got, err, calls)
	}
}
func mustJSON(v any) string { b, _ := json.Marshal(v); return string(b) }

func TestUnapprovedIdentityAndChangedSchemaNeverExecute(t *testing.T) {
	for _, scenario := range []string{"other_user", "other_toolkit", "expired", "changed_schema", "meta_tool", "instant_account", "invalid_arguments"} {
		t.Run(scenario, func(t *testing.T) {
			sends := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				if strings.Contains(r.URL.Path, "/execute/") {
					sends++
					return response(200, `{}`), nil
				}
				if strings.Contains(r.URL.Path, "/connected_accounts/") {
					body := accountJSON
					switch scenario {
					case "other_user":
						body = strings.Replace(body, "owner_one", "owner_two", 1)
					case "other_toolkit":
						body = strings.Replace(body, `"desk"`, `"mail"`, 1)
					case "expired":
						body = strings.Replace(body, "ACTIVE", "EXPIRED", 1)
					}
					return response(200, body), nil
				}
				body := toolJSON
				if scenario == "changed_schema" {
					body = strings.Replace(body, `"type":"string"`, `"type":"integer"`, 1)
				}
				return response(200, body), nil
			})
			tl := approvedTool(t)
			account := "account_one"
			args := json.RawMessage(`{}`)
			if scenario == "meta_tool" {
				tl.Slug = "COMPOSIO_MULTI_EXECUTE_TOOL"
			}
			if scenario == "instant_account" {
				account = "instant_account"
			}
			if scenario == "invalid_arguments" {
				args = json.RawMessage(`[]`)
			}
			_, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, account, tl, args, func() bool { return true })
			if err == nil || sends != 0 {
				t.Fatalf("%s executed: sends=%d err=%v", scenario, sends, err)
			}
		})
	}
}

func TestExecutionFailuresAreNeverRetried(t *testing.T) {
	for _, status := range []int{0, 302, 408, 500, 429, 400, 200} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			sends := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				if strings.Contains(r.URL.Path, "/connected_accounts/") {
					return response(200, accountJSON), nil
				}
				if !strings.Contains(r.URL.Path, "/execute/") {
					return response(200, toolJSON), nil
				}
				sends++
				if status == 0 {
					return nil, errors.New("answer lost")
				}
				res := response(status, `{"successful":false,"error":"refused"}`)
				if status == 302 {
					res.Header.Set("Location", "https://another.example/leak")
				}
				return res, nil
			})
			got, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one", approvedTool(t), json.RawMessage(`{}`), func() bool { return true })
			unknown := status == 0 || status == 302 || status == 408 || status == 500
			if err != nil || sends != 1 || got.Unknown != unknown || !unknown && !got.Refused {
				t.Fatalf("status%d: %+v %v sends=%d", status, got, err, sends)
			}
		})
	}
}

func TestCatalogKeepsPaginationAndUsesUpstreamApps(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/api/v3.1/auth_configs" {
			return response(200, `{"items":[{"id":"ac_app_outside_starter_list","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"app_outside_starter_list"}}]}`), nil
		}
		if r.URL.Query().Get("cursor") == "page-three" {
			return response(200, `{"items":[]}`), nil
		}
		if r.URL.Query().Get("search") != "project" || r.URL.Query().Get("cursor") != "page-two" || r.URL.Path != "/api/v3.1/toolkits" {
			t.Fatal("catalog query changed")
		}
		return response(200, `{"items":[{"slug":"app_outside_starter_list","name":"Project app","meta":{"logo":"https://assets.example/icon.png","description":"Project management"}}],"next_cursor":"page-three"}`), nil
	})
	page, err := c.Catalog(context.Background(), "project", "page-two")
	if err != nil || len(page.Items) != 1 || page.Items[0].Slug != "app_outside_starter_list" || page.Next != "" {
		t.Fatalf("catalog: %+v %v", page, err)
	}
}

func TestFlatSchemasKeepNamesAndRequiredArguments(t *testing.T) {
	raw := json.RawMessage(`{"description":{"type":"string","required":true,"enum":["one","two"]},"type":{"type":"integer","minimum":1,"required":false}}`)
	got, err := inputSchema(raw)
	if err != nil {
		t.Fatal(err)
	}
	var s struct {
		Properties map[string]map[string]any `json:"properties"`
		Required   []string                  `json:"required"`
	}
	if json.Unmarshal(got, &s) != nil || s.Properties["description"]["enum"] == nil || s.Properties["type"]["minimum"] != float64(1) || len(s.Required) != 1 || s.Required[0] != "description" {
		t.Fatalf("schema changed: %s", got)
	}
}

func TestRevocationIsAccountScopedAndReportsThePendingJob(t *testing.T) {
	deletes := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.Method == "GET" {
			return response(200, accountJSON), nil
		}
		deletes++
		if r.Method != "DELETE" || r.URL.Path != "/api/v3.1/connected_accounts/account_one" || r.URL.Query().Get("revoke_on_delete") != "true" {
			t.Fatal("wrong revocation target")
		}
		return response(200, `{"success":true,"revoke_job_id":"revocation_one"}`), nil
	})
	if _, err := c.Revoke(context.Background(), connector.Identity{User: "owner_two", Toolkit: "desk"}, "account_one"); err == nil || deletes != 0 {
		t.Fatal("another account was revoked")
	}
	receipt, err := c.Revoke(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one")
	if err != nil || !receipt.Disabled || !receipt.Pending || receipt.Job != "revocation_one" || deletes != 1 {
		t.Fatalf("receipt: %+v %v", receipt, err)
	}
}

func TestAuthorizationRestrictsTheSessionAndChecksItsUser(t *testing.T) {
	for _, other := range []bool{false, true} {
		t.Run(fmt.Sprint(other), func(t *testing.T) {
			links := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/api/v3.1/tool_router/session" {
					var body struct {
						User     string `json:"user_id"`
						Toolkits struct {
							Enable []string `json:"enable"`
						} `json:"toolkits"`
						Execute struct {
							Multi bool `json:"enable_multi_execute"`
						} `json:"execute"`
						Workbench map[string]any `json:"workbench"`
					}
					if json.NewDecoder(r.Body).Decode(&body) != nil || body.User != "owner_one" || len(body.Toolkits.Enable) != 1 || body.Toolkits.Enable[0] != "desk" || body.Execute.Multi || len(body.Workbench) != 1 || body.Workbench["enable"] != false {
						t.Fatal("session was not restricted")
					}
					user := "owner_one"
					if other {
						user = "owner_two"
					}
					return response(201, `{"session_id":"session_one","config":{"user_id":"`+user+`"}}`), nil
				}
				links++
				if r.URL.Path != "/api/v3.1/tool_router/session/session_one/link" {
					t.Fatal("wrong session")
				}
				var body map[string]string
				_ = json.NewDecoder(r.Body).Decode(&body)
				if body["toolkit"] != "desk" || body["callback_url"] != "https://planque.example/connected" {
					t.Fatal("link target changed")
				}
				return response(201, `{"connected_account_id":"account_one","redirect_url":"https://app.composio.dev/link/link_one"}`), nil
			})
			got, err := c.Authorize(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "connection_one", "https://planque.example/connected")
			if other {
				if err == nil || links != 0 {
					t.Fatal("another user's session created a sign-in link")
				}
				return
			}
			if err != nil || got.Account != "account_one" || links != 1 {
				t.Fatalf("authorization %+v %v", got, err)
			}
		})
	}
}

func TestMalformedExecutionAnswerRemainsUnknown(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if strings.Contains(r.URL.Path, "/connected_accounts/") {
			return response(200, accountJSON), nil
		}
		if strings.Contains(r.URL.Path, "/execute/") {
			return response(200, `not JSON`), nil
		}
		return response(200, toolJSON), nil
	})
	got, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one", approvedTool(t), json.RawMessage(`{}`), func() bool { return true })
	if err != nil || !got.Unknown {
		t.Fatalf("malformed answer: %+v %v", got, err)
	}
}

func TestPermitExpiryDuringPreflightNeverExecutes(t *testing.T) {
	sends, secrets := 0, 0
	live := true
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/connected_accounts/account_one":
			return response(200, accountJSON), nil
		case "/api/v3.1/tools/DESK_REPLY":
			return response(200, toolJSON), nil
		default:
			sends++
			return response(200, `{"data":{}}`), nil
		}
	})
	c.Secret = func(context.Context) (string, error) {
		secrets++
		if secrets == 3 {
			live = false
		}
		return "synthetic-project-key", nil
	}
	_, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one", approvedTool(t), json.RawMessage(`{}`), func() bool { return live })
	if !errors.Is(err, connector.ErrPermit) || sends != 0 {
		t.Fatalf("expired during secret resolution: err=%v sends=%d", err, sends)
	}
}

func TestDiscoveryRetainsUnderscoreSlugsAndSkipsOddEntries(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/auth_configs":
			return response(200, `{"items":[{"id":"ac__1password","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"_1password"}},{"id":"ac__21risk","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"_21risk"}},{"id":"ac__2chat","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"_2chat"}},{"id":"ac_desk","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"desk"}}]}`), nil
		case "/api/v3.1/toolkits":
			if r.URL.Query().Get("cursor") == "next_page" {
				return response(200, `{"items":[]}`), nil
			}
			return response(200, `{"items":[{"slug":"_1password","name":"1Password"},{"slug":"_21risk","name":"21Risk"},{"slug":"_2chat","name":"2Chat"},{"slug":"../bad","name":"Bad"},{"slug":42,"name":"Wrong type"},{"slug":"empty_name","name":""},{"slug":"desk","name":"Desk"}],"next_cursor":"next_page"}`), nil
		case "/api/v3.1/toolkits/_1password":
			return response(200, `{"slug":"_1password","name":"1Password"}`), nil
		case "/api/v3.1/tools":
			if r.URL.Query().Get("toolkit_slug") != "_1password" {
				t.Fatal("underscore toolkit query changed")
			}
			valid := strings.ReplaceAll(strings.ReplaceAll(toolJSON, "DESK_REPLY", "_1PASSWORD_GET_ITEM"), `"desk"`, `"_1password"`)
			wrong := strings.ReplaceAll(valid, `"_1password"`, `"another_app"`)
			meta := strings.ReplaceAll(valid, "_1PASSWORD_GET_ITEM", "COMPOSIO_EXECUTE_TOOL")
			return response(200, `{"items":[null,{"slug":42},`+wrong+`,`+meta+`,`+valid+`],"next_cursor":"tool_page"}`), nil
		case "/api/v3.1/tools/_1PASSWORD_GET_ITEM":
			return response(200, strings.ReplaceAll(strings.ReplaceAll(toolJSON, "DESK_REPLY", "_1PASSWORD_GET_ITEM"), `"desk"`, `"_1password"`)), nil
		default:
			t.Fatalf("unexpected request: %s", r.URL)
			return nil, nil
		}
	})
	page, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(page.Items) != 4 || page.Skipped != 3 || page.Next != "" {
		t.Fatalf("discovery page: %+v %v", page, err)
	}
	for i, slug := range []string{"_1password", "_21risk", "_2chat", "desk"} {
		if page.Items[i].Slug != slug {
			t.Fatal("valid app was discarded")
		}
	}
	app, err := c.App(context.Background(), "_1password")
	if err != nil || app.Slug != "_1password" {
		t.Fatal("direct underscore app lookup failed")
	}
	tools, err := c.Tools(context.Background(), app.Slug, "")
	if err != nil || len(tools.Items) != 1 || tools.Skipped != 4 || tools.Next != "tool_page" || tools.Items[0].Validate() != nil {
		t.Fatalf("tool page: %+v %v", tools, err)
	}
	tool, err := c.Tool(context.Background(), tools.Items[0].Slug, tools.Items[0].Version)
	if err != nil || tool.Toolkit != "_1password" || tool.Slug != "_1PASSWORD_GET_ITEM" {
		t.Fatal("direct underscore tool lookup failed")
	}
	for _, bad := range []string{"../escape", "%2Fescape", "with space", "", strings.Repeat("_", 257)} {
		if _, err := c.App(context.Background(), bad); err == nil {
			t.Fatal("unsafe app identifier accepted")
		}
	}
}

// A refused or lost call stays ErrUnavailable, and says which call and what
// the provider answered, for the log: never the body, the key or a query.
func TestARefusedCallSaysWhichCallAndItsStatus(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/link") {
			return response(400, `{"error":{"message":"synthetic refusal detail"}}`), nil
		}
		return response(201, `{"session_id":"session_one","config":{"user_id":"owner_one"}}`), nil
	})
	_, err := c.Authorize(context.Background(), connector.Identity{User: "owner_one", Toolkit: "gmail"}, "connection_one", "https://planque.example/connected")
	var call *connector.CallError
	if !errors.Is(err, connector.ErrUnavailable) || !errors.As(err, &call) || call.Status != 400 || call.Call != "POST /tool_router/session/session_one/link" {
		t.Fatalf("err = %v", err)
	}
	if msg := err.Error(); !strings.Contains(msg, "composio: POST /tool_router/session/session_one/link: 400") || strings.Contains(msg, "synthetic") {
		t.Fatalf("message = %q", msg)
	}
	lost := testClient(func(*http.Request) (*http.Response, error) { return nil, errors.New("synthetic network failure") })
	_, err = lost.Authorize(context.Background(), connector.Identity{User: "owner_one", Toolkit: "gmail"}, "connection_one", "https://planque.example/connected")
	if !errors.Is(err, connector.ErrUnavailable) || !strings.HasSuffix(err.Error(), "POST /tool_router/session: no answer") {
		t.Fatalf("lost: %v", err)
	}
}

// The catalog offers only apps the project can sign someone in to: an
// enabled auth config, read across pages and cached briefly. A disabled or
// missing one hides the app and refuses to start it; a failed read offers
// nothing.
func TestOnlyAppsWithAnEnabledSignInAreOffered(t *testing.T) {
	reads, fail := 0, false
	now := time.Unix(1800000000, 0)
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/auth_configs":
			reads++
			if fail {
				return response(500, `{}`), nil
			}
			if r.URL.Query().Get("cursor") == "" {
				return response(200, `{"items":[{"id":"ac_gmail","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"gmail"}},{"id":"ac_notion","auth_scheme":"OAUTH2","status":"DISABLED","toolkit":{"slug":"notion"}}],"next_cursor":"more"}`), nil
			}
			return response(200, `{"items":[{"id":"ac_slack","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"slack"}}]}`), nil
		case "/api/v3.1/toolkits":
			return response(200, `{"items":[{"slug":"gmail","name":"Gmail"},{"slug":"notion","name":"Notion"},{"slug":"slack","name":"Slack"},{"slug":"hubspot","name":"HubSpot"}]}`), nil
		case "/api/v3.1/toolkits/gmail":
			return response(200, `{"slug":"gmail","name":"Gmail"}`), nil
		case "/api/v3.1/toolkits/notion", "/api/v3.1/toolkits/hubspot":
			return response(200, `{"slug":"`+strings.TrimPrefix(r.URL.Path, "/api/v3.1/toolkits/")+`","name":"An app","auth_schemes":["OAUTH2"]}`), nil
		}
		t.Fatalf("unexpected request: %s", r.URL)
		return nil, nil
	})
	c.Now = func() time.Time { return now }
	page, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(page.Items) != 2 || page.Items[0].Slug != "gmail" || page.Items[1].Slug != "slack" {
		t.Fatalf("catalog: %+v %v", page, err)
	}
	if app, err := c.App(context.Background(), "gmail"); err != nil || app.Setup != connector.SetupReady || app.AuthConfig != "ac_gmail" {
		t.Fatalf("gmail: %+v %v", app, err)
	}
	// An app with no enabled sign-in and no managed one is known, and needs setting up.
	for _, id := range []string{"notion", "hubspot"} {
		if app, err := c.App(context.Background(), id); err != nil || app.Setup != connector.SetupNeedsSetup || app.AuthConfig != "" {
			t.Fatalf("%s: %+v %v", id, app, err)
		}
	}
	if reads != 2 {
		t.Fatalf("the sign-ins were read %d times, want one read of two pages", reads)
	}
	now = now.Add(configuredFor + time.Second)
	fail = true
	if _, err := c.Catalog(context.Background(), "", ""); !errors.Is(err, connector.ErrUnavailable) {
		t.Fatalf("a failed read still offered apps: %v", err)
	}
}

// Connect apps opens on one upstream read of the whole catalog, kept with the
// sign-ins it was filtered by: a second open reads nothing; a search, or new
// sign-ins after the cache, read again.
func TestTheFirstPageIsOneReadAndKept(t *testing.T) {
	toolkits, configs := 0, 0
	now := time.Unix(1800000000, 0)
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/auth_configs":
			configs++
			return response(200, `{"items":[{"id":"ac_gmail","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"gmail"}}]}`), nil
		case "/api/v3.1/toolkits":
			toolkits++
			if r.URL.Query().Get("limit") != "1000" {
				t.Errorf("limit = %s", r.URL.Query().Get("limit"))
			}
			return response(200, `{"items":[{"slug":"gmail","name":"Gmail"},{"slug":"hubspot","name":"HubSpot"}]}`), nil
		}
		t.Fatalf("unexpected request: %s", r.URL)
		return nil, nil
	})
	c.Now = func() time.Time { return now }
	for range 2 {
		if page, err := c.Catalog(context.Background(), "", ""); err != nil || len(page.Items) != 1 || page.Items[0].Slug != "gmail" {
			t.Fatalf("first page: %+v %v", page, err)
		}
	}
	if toolkits != 1 || configs != 1 {
		t.Fatalf("two opens read the catalog %d and the sign-ins %d times", toolkits, configs)
	}
	if _, err := c.Catalog(context.Background(), "gm", ""); err != nil || toolkits != 2 {
		t.Fatalf("a search wasn't read: %v (%d)", err, toolkits)
	}
	now = now.Add(configuredFor + time.Second)
	if _, err := c.Catalog(context.Background(), "", ""); err != nil || toolkits != 3 || configs != 2 {
		t.Fatalf("after the sign-ins expired: %v (%d, %d)", err, toolkits, configs)
	}
}

// A search lists every app the provider has, each saying how its sign-in
// stands; the unsearched first page keeps to the apps set up already. A read
// that stops at its bound with pages left says so (kiff/planque#87).
func TestCatalogSaysHowEachAppIsSetUpAndWhenARegionWasNotRead(t *testing.T) {
	pages := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/auth_configs":
			return response(200, `{"items":[{"id":"ac_desk","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"desk"}}]}`), nil
		case "/api/v3.1/toolkits":
			pages++
			items := `{"slug":"desk","name":"Desk","auth_schemes":["OAUTH2"]},` +
				`{"slug":"mail","name":"Mail","composio_managed_auth_schemes":["OAUTH2"],"auth_schemes":["OAUTH2","API_KEY"]},` +
				`{"slug":"keyed","name":"Keyed","auth_schemes":["API_KEY"]},` +
				`{"slug":"open","name":"Open","no_auth":true},` +
				`{"slug":"own","name":"Own app","auth_schemes":["OAUTH2"]}`
			next := ""
			if r.URL.Query().Get("search") == "endless" {
				next = fmt.Sprintf(`,"next_cursor":"p%d"`, pages)
			}
			return response(200, `{"items":[`+items+`]`+next+`}`), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
			return nil, nil
		}
	})
	first, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(first.Items) != 1 || first.Items[0].Slug != "desk" || first.Items[0].Setup != connector.SetupReady || first.Partial {
		t.Fatalf("first page: %+v err=%v", first, err)
	}
	all, err := c.Catalog(context.Background(), "app", "")
	if err != nil || len(all.Items) != 5 || all.Partial {
		t.Fatalf("search: %+v err=%v", all, err)
	}
	want := map[string][2]string{"desk": {connector.SetupReady, "OAUTH2"}, "mail": {connector.SetupOnDemand, "OAUTH2"}, "keyed": {connector.SetupOnDemand, "API_KEY"}, "open": {connector.SetupOnDemand, "NO_AUTH"}, "own": {connector.SetupNeedsSetup, "OAUTH2"}}
	for _, a := range all.Items {
		if w := want[a.Slug]; a.Setup != w[0] || a.AuthScheme != w[1] || (a.Setup == connector.SetupReady) != (a.AuthConfig != "") {
			t.Errorf("%s: setup %q scheme %q config %q", a.Slug, a.Setup, a.AuthScheme, a.AuthConfig)
		}
	}
	pages = 0
	endless, err := c.Catalog(context.Background(), "endless", "")
	if err != nil || !endless.Partial || endless.Next == "" || pages != catalogPages {
		t.Fatalf("an endless catalog: partial=%v next=%q pages=%d err=%v", endless.Partial, endless.Next, pages, err)
	}
}

// Setting an app up asks the provider for one managed sign-in, takes it only
// when the provider says it manages it, and offers the app from then on.
func TestSetUpMakesOneManagedSignInAndOffersTheApp(t *testing.T) {
	made := 0
	managed := true
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.URL.Path == "/api/v3.1/auth_configs" && r.Method == "GET":
			// The provider lists a sign-in once made, as the fresh read before a set-up sees it.
			if made > 0 {
				return response(200, `{"items":[{"id":"ac_mail_new","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"mail"}}]}`), nil
			}
			return response(200, `{"items":[]}`), nil
		case r.URL.Path == "/api/v3.1/auth_configs" && r.Method == "POST":
			made++
			var body map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if len(body) != 2 || string(body["toolkit"]) != `{"slug":"mail"}` || string(body["auth_config"]) != `{"type":"use_composio_managed_auth"}` {
				t.Fatalf("set-up request: %s", mustJSON(body))
			}
			return response(201, fmt.Sprintf(`{"toolkit":{"slug":"mail"},"auth_config":{"id":"ac_mail_new","auth_scheme":"oauth2","is_composio_managed":%v}}`, managed)), nil
		case r.URL.Path == "/api/v3.1/toolkits/mail":
			return response(200, `{"slug":"mail","name":"Mail","composio_managed_auth_schemes":["OAUTH2"]}`), nil
		case r.URL.Path == "/api/v3.1/toolkits/own":
			return response(200, `{"slug":"own","name":"Own app","auth_schemes":["OAUTH2"]}`), nil
		case r.URL.Path == "/api/v3.1/toolkits/gone":
			return response(404, `{"error":{"message":"no such toolkit"}}`), nil
		case r.URL.Path == "/api/v3.1/toolkits":
			return response(200, `{"items":[{"slug":"mail","name":"Mail","composio_managed_auth_schemes":["OAUTH2"]}]}`), nil
		default:
			t.Fatalf("unexpected path %s %s", r.Method, r.URL.Path)
			return nil, nil
		}
	})
	if _, err := c.SetUp(context.Background(), "own"); !errors.Is(err, connector.ErrNotSetUp) || made != 0 {
		t.Fatalf("an app that needs setup was set up: %v (%d)", err, made)
	}
	if _, err := c.App(context.Background(), "gone"); !errors.Is(err, connector.ErrNotSetUp) {
		t.Fatalf("an app the provider doesn't have: %v", err)
	}
	app, err := c.SetUp(context.Background(), "mail")
	if err != nil || made != 1 || app.Setup != connector.SetupReady || app.AuthConfig != "ac_mail_new" || app.AuthScheme != "OAUTH2" {
		t.Fatalf("set up: %+v err=%v made=%d", app, err, made)
	}
	// From now on the app is offered without a search, and set up again asks nothing.
	first, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(first.Items) != 1 || first.Items[0].AuthConfig != "ac_mail_new" {
		t.Fatalf("offered after set-up: %+v err=%v", first, err)
	}
	if again, err := c.SetUp(context.Background(), "mail"); err != nil || made != 1 || again.AuthConfig != "ac_mail_new" {
		t.Fatalf("set up twice: %+v err=%v made=%d", again, err, made)
	}
	// An answer the provider doesn't manage is refused, and nothing is pinned.
	managed, made = false, 0
	c2 := testClient(c.HTTP.Transport.(transport))
	if _, err := c2.SetUp(context.Background(), "mail"); !errors.Is(err, connector.ErrDefinition) {
		t.Fatalf("an unmanaged sign-in was taken: %v", err)
	}
}

// A search by words asks the provider's text search, narrows to one app
// upstream and to several here, never returns a platform meta tool, and
// brings each tool's app as it stands, read once per app.
func TestSearchFindsToolsAcrossAppsAndSaysHowEachAppStands(t *testing.T) {
	appReads := map[string]int{}
	var seen url.Values
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.URL.Path == "/api/v3.1/auth_configs":
			return response(200, `{"items":[{"id":"ac_desk","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"desk"}}]}`), nil
		case r.URL.Path == "/api/v3.1/tools":
			seen = r.URL.Query()
			return response(200, `{"items":[`+toolJSON+`,`+
				`{"slug":"MAIL_SEND","name":"Send an email","toolkit":{"slug":"mail"},"version":"20300101_00","input_parameters":{"type":"object","properties":{"to":{"type":"string"}}},"tags":[]},`+
				`{"slug":"COMPOSIO_SEARCH_TOOLS","name":"Search tools","toolkit":{"slug":"composio"},"version":"20300101_00","input_parameters":{"type":"object","properties":{}},"tags":[]},`+
				`{"slug":"OWN_READ","name":"Read","toolkit":{"slug":"own"},"version":"20300101_00","input_parameters":{"type":"object","properties":{}},"tags":["readOnlyHint"]}],"next_cursor":"p2"}`), nil
		case strings.HasPrefix(r.URL.Path, "/api/v3.1/toolkits/"):
			slug := strings.TrimPrefix(r.URL.Path, "/api/v3.1/toolkits/")
			appReads[slug]++
			managed := `,"composio_managed_auth_schemes":["OAUTH2"]`
			if slug == "own" {
				managed = `,"auth_schemes":["OAUTH2"]`
			}
			return response(200, `{"slug":"`+slug+`","name":"The `+slug+` app"`+managed+`}`), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
			return nil, nil
		}
	})
	found, err := c.Search(context.Background(), " reply ticket ", nil, "")
	if err != nil || found.Next != "p2" || found.Skipped != 1 {
		t.Fatalf("search: %+v err=%v", found, err)
	}
	if seen.Get("query") != "reply ticket" || seen.Get("limit") != "50" || seen.Has("toolkit_slug") {
		t.Fatalf("asked with %v", seen)
	}
	slugs := []string{}
	for _, tl := range found.Tools {
		slugs = append(slugs, tl.Slug)
	}
	if strings.Join(slugs, " ") != "DESK_REPLY MAIL_SEND OWN_READ" || !found.Tools[2].ReadOnly {
		t.Fatalf("tools: %v", slugs)
	}
	if found.Apps["desk"].Setup != connector.SetupReady || found.Apps["mail"].Setup != connector.SetupOnDemand || found.Apps["own"].Setup != connector.SetupNeedsSetup || found.Apps["own"].Name != "The own app" {
		t.Fatalf("apps: %+v", found.Apps)
	}
	if appReads["desk"] != 1 || appReads["mail"] != 1 || appReads["own"] != 1 || len(appReads) != 3 {
		t.Fatalf("app reads: %v", appReads)
	}
	// One app narrows upstream; several narrow here.
	one, err := c.Search(context.Background(), "reply", []string{"desk"}, "")
	if err != nil || seen.Get("toolkit_slug") != "desk" || len(one.Tools) != 1 || one.Tools[0].Slug != "DESK_REPLY" {
		t.Fatalf("one app: %+v %v err=%v", one.Tools, seen, err)
	}
	two, err := c.Search(context.Background(), "reply", []string{"desk", "own"}, "")
	if err != nil || seen.Has("toolkit_slug") || len(two.Tools) != 2 {
		t.Fatalf("two apps: %+v err=%v", two.Tools, err)
	}
	for _, bad := range [][]string{{"../x"}, make([]string, 21)} {
		if _, err := c.Search(context.Background(), "reply", bad, ""); !errors.Is(err, connector.ErrDefinition) {
			t.Fatalf("bad apps %v: %v", bad, err)
		}
	}
	if _, err := c.Search(context.Background(), "  ", nil, ""); !errors.Is(err, connector.ErrDefinition) {
		t.Fatalf("no words: %v", err)
	}
}

// An app's events are listed a page at a time, each pinned to a version;
// one that isn't a definition is left out and counted.
func TestTriggerTypesListsAnAppsEvents(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v3.1/triggers_types" || r.URL.Query().Get("toolkit_slugs") != "desk" || r.URL.Query().Get("limit") != "50" {
			t.Fatalf("unexpected request %s", r.URL)
		}
		return response(200, `{"items":[`+
			`{"slug":"DESK_NEW_TICKET","name":"New ticket","description":"A ticket was opened.","version":"20300101_00","toolkit":{"slug":"desk"},"config":{"type":"object","properties":{"queue":{"type":"string"}}},"payload":{"type":"object"}},`+
			`{"slug":"DESK_LATEST","name":"Latest","version":"latest","toolkit":{"slug":"desk"},"config":{},"payload":{}},`+
			`{"slug":"OTHER_EVENT","name":"Elsewhere","version":"20300101_00","toolkit":{"slug":"other"},"config":{},"payload":{}}],"next_cursor":"p2"}`), nil
	})
	page, err := c.TriggerTypes(context.Background(), "desk", "")
	if err != nil || len(page.Items) != 1 || page.Skipped != 2 || page.Next != "p2" {
		t.Fatalf("events: %+v err=%v", page, err)
	}
	if e := page.Items[0]; e.Slug != "DESK_NEW_TICKET" || e.Description != "A ticket was opened." || e.Version != "20300101_00" || !strings.Contains(string(e.Config), "queue") {
		t.Fatalf("event: %+v", e)
	}
	if _, err := c.TriggerTypes(context.Background(), "../x", ""); !errors.Is(err, connector.ErrDefinition) {
		t.Fatalf("bad app: %v", err)
	}
}

// A refused call keeps who refused it and why: the vendor's own answer passed
// on by the provider, or the provider's own refusal with its status, each
// redacted (the project key never comes back in a reason).
func TestARefusalKeepsWhoRefusedAndWhy(t *testing.T) {
	run := func(status int, body string) connector.Result {
		t.Helper()
		c := testClient(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case "/api/v3.1/connected_accounts/account_one":
				return response(200, accountJSON), nil
			case "/api/v3.1/tools/DESK_REPLY":
				return response(200, toolJSON), nil
			case "/api/v3.1/tools/execute/DESK_REPLY":
				return response(status, body), nil
			}
			t.Fatalf("unexpected path %s", r.URL.Path)
			return nil, nil
		})
		got, err := c.Execute(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk"}, "account_one", approvedTool(t), json.RawMessage(`{"ticket":"ticket_one","text":"hello"}`), func() bool { return true })
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	if got := run(200, `{"data":{},"successful":false,"error":"No such ticket: 'ticket_one'"}`); !got.Refused || !got.ByVendor || got.Reason != "No such ticket: 'ticket_one'" || got.Status != 0 {
		t.Fatalf("the vendor's refusal: %+v", got)
	}
	// The vendor's status, when the provider passes it on with the refusal (a refused key).
	if got := run(200, `{"data":{"http_error":"Forbidden","status_code":403},"successful":false,"error":"Forbidden"}`); !got.Refused || !got.ByVendor || got.Status != 403 || got.Reason != "Forbidden" {
		t.Fatalf("the vendor's refusal with its status: %+v", got)
	}
	for _, data := range []string{`{"status_code":42}`, `{"status_code":403.5}`, `[]`} {
		if got := run(200, `{"data":`+data+`,"successful":false,"error":"no"}`); got.Status != 0 {
			t.Fatalf("a status from %s: %+v", data, got)
		}
	}
	if got := run(403, `{"error":{"message":"key synthetic-project-key may not use this account"}}`); !got.Refused || got.ByVendor || got.Status != 403 || got.Reason != "key [redacted] may not use this account" {
		t.Fatalf("the provider's refusal: %+v", got)
	}
	if got := run(502, `{"message":"upstream timed out"}`); !got.Unknown || got.Refused || got.Status != 502 || got.Reason != "upstream timed out" {
		t.Fatalf("no confirmed result: %+v", got)
	}
}
