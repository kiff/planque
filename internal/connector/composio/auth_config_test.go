package composio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
)

func TestAuthConfigSelectionKeepsTheSameSchemeAcrossPages(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		first := `{"id":"ac_z","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"desk"}}`
		second := `{"id":"ac_a","auth_scheme":"API_KEY","status":"ENABLED","toolkit":{"slug":"desk"}}`
		if reverse {
			first, second = second, first
		}
		c := testClient(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case "/api/v3.1/auth_configs":
				if r.URL.Query().Get("cursor") == "more" {
					return response(200, `{"items":[`+second+`,{"id":"ac_0","auth_scheme":"BEARER_TOKEN","status":"DISABLED","toolkit":{"slug":"desk"}}]}`), nil
				}
				return response(200, `{"items":[`+first+`,{"id":"../bad","auth_scheme":"API_KEY","status":"ENABLED","toolkit":{"slug":"desk"}},{"auth_scheme":"API_KEY","status":"ENABLED","toolkit":{"slug":"missing_id"}},{"id":"ac_missing_scheme","status":"ENABLED","toolkit":{"slug":"missing_scheme"}}],"next_cursor":"more"}`), nil
			case "/api/v3.1/toolkits":
				return response(200, `{"items":[{"slug":"desk","name":"Support desk"},{"slug":"missing_id","name":"Missing ID"},{"slug":"missing_scheme","name":"Missing scheme"}]}`), nil
			case "/api/v3.1/toolkits/desk":
				return response(200, `{"slug":"desk","name":"Support desk"}`), nil
			}
			t.Fatalf("unexpected request: %s", r.URL)
			return nil, nil
		})
		page, err := c.Catalog(context.Background(), "", "")
		if err != nil || len(page.Items) != 1 {
			t.Fatalf("catalog: %+v %v", page, err)
		}
		app, err := c.App(context.Background(), "desk")
		if err != nil || app.AuthConfig != "ac_a" || app.AuthScheme != "API_KEY" || page.Items[0] != app {
			t.Fatalf("selection changed with page order %v: %+v %v", reverse, app, err)
		}
	}
}

func TestAuthSchemeRefreshReplacesCachedCatalogMetadata(t *testing.T) {
	now := time.Unix(1800000000, 0)
	scheme, config := "API_KEY", "ac_key"
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/api/v3.1/auth_configs" {
			return response(200, `{"items":[{"id":"`+config+`","auth_scheme":"`+scheme+`","status":"ENABLED","toolkit":{"slug":"desk"}}]}`), nil
		}
		return response(200, `{"items":[{"slug":"desk","name":"Support desk"}]}`), nil
	})
	c.Now = func() time.Time { return now }
	before, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(before.Items) != 1 || before.Items[0].AuthScheme != "API_KEY" {
		t.Fatalf("before: %+v %v", before, err)
	}
	scheme, config = "BEARER_TOKEN", "ac_token"
	now = now.Add(configuredFor + time.Second)
	after, err := c.Catalog(context.Background(), "", "")
	if err != nil || len(after.Items) != 1 || after.Items[0].AuthScheme != "BEARER_TOKEN" || after.Items[0].AuthConfig != "ac_token" || before.Items[0].AuthConfig != "ac_key" {
		t.Fatalf("refresh: before=%+v after=%+v %v", before, after, err)
	}
}

func TestAuthorizationUsesTheApprovedAuthConfig(t *testing.T) {
	for _, returned := range []string{"ac_approved", "ac_other", ""} {
		t.Run(returned, func(t *testing.T) {
			links := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/api/v3.1/tool_router/session" {
					var in struct {
						AuthConfigs map[string]string `json:"auth_configs"`
					}
					if err := json.NewDecoder(r.Body).Decode(&in); err != nil || len(in.AuthConfigs) != 1 || in.AuthConfigs["desk"] != "ac_approved" {
						t.Fatal("session didn't pin the approved auth config")
					}
					return response(201, `{"session_id":"session_one","config":{"user_id":"owner_one","auth_configs":{"desk":"`+returned+`"}}}`), nil
				}
				if r.URL.Path != "/api/v3.1/tool_router/session/session_one/link" || r.Method != "POST" {
					t.Fatal("wrong hosted link endpoint")
				}
				links++
				return response(201, `{"connected_account_id":"account_one","redirect_url":"https://app.composio.dev/link/demo"}`), nil
			})
			got, err := c.Authorize(context.Background(), connector.Identity{User: "owner_one", Toolkit: "desk", AuthConfig: "ac_approved"}, "connection_one", "https://planque.example/connected")
			if returned != "ac_approved" {
				if !errors.Is(err, connector.ErrDefinition) || links != 0 {
					t.Fatalf("different or missing config created a link: %+v %v links=%d", got, err, links)
				}
			} else if err != nil || got.Account != "account_one" || links != 1 {
				t.Fatalf("hosted link: %+v %v links=%d", got, err, links)
			}
		})
	}
}

// A scoped sign-in (gap 5 of kiff/tafs#476) is never an app's usual one, even
// with an id that sorts first, and the same scope set is found again by its
// name rather than made twice. An operator's own sign-in is never narrowed,
// and a recommendation outside the managed ceiling isn't taken.
func TestAScopedSignInIsNeverAnAppsUsualOne(t *testing.T) {
	sum := sha256.Sum256([]byte("tickets.read"))
	name := "planque-scoped-desk-" + hex.EncodeToString(sum[:8])
	recommend := `{"scopes":{"least_privilege":["tickets.read"]}}`
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v3.1/auth_configs":
			if r.Method == http.MethodPost {
				t.Fatal("a scoped sign-in was made again")
			}
			return response(200, `{"items":[{"id":"ac_0_scoped","name":"`+name+`","auth_scheme":"OAUTH2","status":"ENABLED","is_composio_managed":true,"toolkit":{"slug":"desk"}},`+
				`{"id":"ac_usual","auth_scheme":"OAUTH2","status":"ENABLED","is_composio_managed":true,"toolkit":{"slug":"desk"}},`+
				`{"id":"ac_own","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"own"}}]}`), nil
		case "/api/v3.1/toolkits/desk":
			return response(200, `{"slug":"desk","name":"Support desk","composio_managed_auth":[{"mode":"OAUTH2","scopes":{"available":["tickets.read","tickets.write"]}}],`+
				`"auth_config_details":[{"mode":"OAUTH2","required_scopes":["offline_access"]}]}`), nil
		case "/api/v3.1/toolkits/desk/scopes/recommended":
			return response(200, recommend), nil
		case "/api/v3.1/toolkits/own":
			return response(200, `{"slug":"own","name":"Own app"}`), nil
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL)
		return nil, nil
	})
	ctx := context.Background()
	app, err := c.App(ctx, "desk")
	if err != nil || app.AuthConfig != "ac_usual" || !app.Managed {
		t.Fatalf("the app's usual sign-in: %+v %v", app, err)
	}
	scoped, err := c.Scope(ctx, app, []string{"DESK_LIST_TICKETS"}, connector.Allowance{})
	if err != nil || scoped.AuthConfig != "ac_0_scoped" || strings.Join(scoped.Scopes, " ") != "offline_access tickets.read" {
		t.Fatalf("scoped: %+v %v", scoped, err)
	}
	own, err := c.App(ctx, "own")
	if err != nil || own.Managed {
		t.Fatalf("the operator's app: %+v %v", own, err)
	}
	if _, err := c.Scope(ctx, own, []string{"OWN_SEND"}, connector.Allowance{}); !errors.Is(err, connector.ErrUnscoped) {
		t.Fatalf("an operator's sign-in was narrowed: %v", err)
	}
	recommend = `{"scopes":{"least_privilege":["admin"]}}`
	if _, err := c.Scope(ctx, app, []string{"DESK_LIST_TICKETS"}, connector.Allowance{}); !errors.Is(err, connector.ErrUnscoped) {
		t.Fatalf("a scope outside the ceiling was taken: %v", err)
	}
}
