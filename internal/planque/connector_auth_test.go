package planque

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

func TestHostedKeyAppsReturnConnectedThroughTheSignedFlow(t *testing.T) {
	for _, scheme := range []string{"API_KEY", "BEARER_TOKEN"} {
		t.Run(scheme, func(t *testing.T) {
			e, f := newConnectorFixture(t)
			f.authScheme = scheme
			st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog?provider=composio", nil, nil, "k_setup")
			if st != 200 {
				t.Fatalf("catalog: %d %v", st, out)
			}
			for _, entry := range out["entries"].([]any) {
				if entry.(map[string]any)["how"] != "oauth" {
					t.Fatal("hosted sign-in was replaced by direct key collection")
				}
			}
			v := requestConnector(t, e)
			credential := "key"
			if scheme == "BEARER_TOKEN" {
				credential = "token"
			}
			want := "Connect Support desk for this taf. You'll enter your Support desk " + credential + " on the connection service's page; it keeps the " + credential + ", and the taf never sees it."
			if text := v["text"].(string); !strings.HasPrefix(text, want) || strings.Contains(strings.ToLower(text), "composio") || strings.Contains(text, "agt_taf") {
				t.Fatalf("words: %q", text)
			}
			path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string)
			if st, _ := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, nil, ""); st != 401 || f.authorizations != 0 || f.links != 0 {
				t.Fatalf("unsigned request started sign-in: %d", st)
			}
			// This helper exercises device approval, the real adapter's session
			// and /link, callback account checks, device confirmation, action
			// selection, a second signed approval and profile application.
			g := connectConnector(t, e, f)
			if g.Status != StatusConnected || f.authorizations != 1 || f.links != 1 || f.authConfig != "ac_desk_demo" || f.executions != 0 {
				t.Fatalf("connection: %+v sessions=%d links=%d config=%q", g, f.authorizations, f.links, f.authConfig)
			}
			if g.Request.Connector.AuthScheme != scheme || !strings.HasPrefix(grantText(g.Request), want) {
				t.Fatal("action selection changed the credential-entry words")
			}
			for _, field := range []string{"config", "scheme"} {
				changed := g.Request
				binding := *changed.Connector
				changed.Connector = &binding
				if field == "config" {
					binding.AuthConfig = "ac_different"
				} else {
					binding.AuthScheme = "OAUTH2"
				}
				sum, err := digest(changed)
				if err != nil || sum == g.Digest {
					t.Fatalf("%s is not bound to the approval digest: %v", field, err)
				}
			}
			// An operator's later scheme edit cannot rewrite a saved approval.
			f.authScheme = "OAUTH2"
			saved, err := e.s.loadGrant(context.Background(), tenant, g.ID)
			if err != nil || saved.Request.Connector.AuthScheme != scheme || !strings.HasPrefix(grantText(saved.Request), want) {
				t.Fatalf("saved approval changed: %v", err)
			}
		})
	}
}

func TestHostedKeyCallbackRequiresAnActiveOwnedAccount(t *testing.T) {
	for _, state := range []string{"PENDING", "EXPIRED", "foreign"} {
		t.Run(state, func(t *testing.T) {
			e, f := newConnectorFixture(t)
			f.authScheme = "API_KEY"
			v := requestConnector(t, e)
			path := "/v1/accounts/" + tenant + "/grants/" + v["id"].(string)
			if st, out := e.do("POST", path+"/approve", map[string]any{"digest": v["digest"]}, e.phone, ""); st != 200 {
				t.Fatalf("approve: %d %v", st, out)
			}
			if state == "foreign" {
				f.foreign = true
			} else {
				f.accountStatus = state
			}
			u, _ := url.Parse(f.callback)
			st, _ := e.do("GET", u.RequestURI()+"&status=ACTIVE", nil, nil, "")
			g, err := e.s.loadGrant(context.Background(), tenant, v["id"].(string))
			if st == 200 || st == 302 || err != nil || g.SignedIn || g.Status == StatusConnected {
				t.Fatalf("invalid account completed callback: %d %+v %v", st, g, err)
			}
		})
	}
}

func TestOlderConnectorGrantsKeepTheirEncodingAndWords(t *testing.T) {
	for _, title := range []string{"", "Support desk"} {
		old := `{"provider":"composio","toolkit":"desk_demo"`
		// Made before titles: its words stay as it was signed, agent id and all.
		prefix := "Connect desk_demo through composio for agt_taf. This connection provider holds the vendor credentials."
		if title != "" {
			old += `,"title":"Support desk"`
			// With a title and no name, it's "this taf": the agent id is under
			// Details (kiff/tafs#476). A grant's stored words never change.
			prefix = "Connect Support desk for this taf. Your sign-in to Support desk stays with the connection service; the taf never sees it."
		}
		old += `}`
		var binding ConnectorRequest
		if err := json.Unmarshal([]byte(old), &binding); err != nil {
			t.Fatal(err)
		}
		raw, err := json.Marshal(binding)
		if err != nil || string(raw) != old {
			t.Fatalf("older signed fields changed: %s %v", raw, err)
		}
		q := Request{Agent: "agt_taf", Connector: &binding}
		if got := connectorText(q); got != prefix+" Sign in and confirm on this device, then choose the actions to approve. No actions are granted yet." {
			t.Fatalf("older approval words changed: %q", got)
		}
	}
	// Its actions too: the account and each action's version, as it was signed.
	q := Request{Agent: "agt_taf", Tools: []ToolRequest{{Tool: "connector_x", Name: "reply_ticket"}}, Connector: &ConnectorRequest{Provider: "composio", Toolkit: "desk_demo", Account: "ca_old",
		Tools: []connector.Tool{{Slug: "DESK_REPLY", Toolkit: "desk_demo", Version: "20251001_00", Name: "Reply to a ticket"}}}}
	if got, want := connectorText(q), "Connect desk_demo through composio for agt_taf. This connection provider holds the vendor credentials. Use only connected account ca_old. Approve these exact versioned actions: Reply to a ticket (DESK_REPLY, 20251001_00)."; got != want {
		t.Fatalf("older actions' words changed:\n%q\nwant\n%q", got, want)
	}
}
