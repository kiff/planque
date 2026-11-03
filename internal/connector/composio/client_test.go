package composio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

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
		if r.URL.Query().Get("search") != "project" || r.URL.Query().Get("cursor") != "page-two" || r.URL.Path != "/api/v3.1/toolkits" {
			t.Fatal("catalog query changed")
		}
		return response(200, `{"items":[{"slug":"app_outside_starter_list","name":"Project app","meta":{"logo":"https://assets.example/icon.png","description":"Project management"}}],"next_cursor":"page-three"}`), nil
	})
	page, err := c.Catalog(context.Background(), "project", "page-two")
	if err != nil || len(page.Items) != 1 || page.Items[0].Slug != "app_outside_starter_list" || page.Next != "page-three" {
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
						Workbench struct {
							Enable bool `json:"enable"`
							Proxy  bool `json:"proxy_execution_enabled"`
						} `json:"workbench"`
					}
					if json.NewDecoder(r.Body).Decode(&body) != nil || body.User != "owner_one" || len(body.Toolkits.Enable) != 1 || body.Toolkits.Enable[0] != "desk" || body.Execute.Multi || body.Workbench.Enable || body.Workbench.Proxy {
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
		case "/api/v3.1/toolkits":
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
	if err != nil || len(page.Items) != 4 || page.Skipped != 3 || page.Next != "next_page" {
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
