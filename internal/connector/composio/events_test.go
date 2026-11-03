package composio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

const triggerJSON = `{"slug":"SLACKBOT_CHANNEL_MESSAGE_RECEIVED","name":"Channel message","toolkit":{"slug":"slackbot"},"version":"20300101_00","config":{"type":"object"},"payload":{"type":"object"}}`

func TestEventsPinAccountVersionAndLeaveOtherProjectWebhookUntouched(t *testing.T) {
	for _, scenario := range []string{"accepted", "foreign", "changed", "forbidden", "other_webhook"} {
		t.Run(scenario, func(t *testing.T) {
			creates := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				switch r.URL.Path {
				case "/api/v3.1/connected_accounts/account_one":
					body := strings.ReplaceAll(accountJSON, `"desk"`, `"slackbot"`)
					if scenario == "foreign" {
						body = strings.ReplaceAll(body, "owner_one", "another_owner")
					}
					return response(200, body), nil
				case "/api/v3.1/triggers_types/SLACKBOT_CHANNEL_MESSAGE_RECEIVED":
					if r.URL.Query().Get("toolkit_versions[slackbot]") != "20300101_00" {
						t.Error("trigger version was not pinned")
					}
					body := triggerJSON
					if scenario == "changed" {
						body = strings.Replace(body, `"config":{"type":"object"}`, `"config":{"type":"array"}`, 1)
					}
					return response(200, body), nil
				case "/api/v3.1/trigger_instances/SLACKBOT_CHANNEL_MESSAGE_RECEIVED/upsert":
					creates++
					if scenario == "forbidden" {
						return response(403, `{"error":"private upstream body"}`), nil
					}
					var body map[string]json.RawMessage
					_ = json.NewDecoder(r.Body).Decode(&body)
					if len(body) != 4 || string(body["user_id"]) != `"owner_one"` || string(body["connected_account_id"]) != `"account_one"` || string(body["toolkit_versions"]) != `{"slackbot":"20300101_00"}` {
						t.Error("trigger identity/version binding changed")
					}
					return response(201, `{"trigger_id":"trigger_one"}`), nil
				case "/api/v3.1/webhook_subscriptions":
					if r.Method != "GET" {
						t.Error("another application's webhook was changed")
					}
					return response(200, `{"items":[{"id":"hook_one","webhook_url":"https://other.example/webhook","version":"V3","secret":"private-webhook-secret","enabled_events":["composio.trigger.message"]}]}`), nil
				default:
					t.Fatal("unexpected event API request")
					return nil, nil
				}
			})
			definition := connector.TriggerType{Slug: "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", Toolkit: "slackbot", Version: "20300101_00", Name: "Channel message", Config: json.RawMessage(`{"type":"object"}`), Payload: json.RawMessage(`{"type":"object"}`)}
			if scenario == "other_webhook" {
				if _, err := c.EnsureWebhook(context.Background(), "https://ours.example/webhook"); !errors.Is(err, connector.ErrWebhook) {
					t.Fatal("foreign webhook was accepted")
				}
				return
			}
			id, err := c.CreateTrigger(context.Background(), connector.Identity{User: "owner_one", Toolkit: "slackbot"}, "account_one", definition, json.RawMessage(`{"is_bot_message":false}`))
			switch scenario {
			case "accepted":
				if err != nil || id != "trigger_one" || creates != 1 {
					t.Fatalf("create: %s %v", id, err)
				}
			case "foreign", "changed":
				if err == nil || creates != 0 {
					t.Fatal("unapproved identity/definition registered")
				}
			case "forbidden":
				if !errors.Is(err, connector.ErrPermission) || creates != 1 {
					t.Fatal("permission refusal not preserved")
				}
			}
		})
	}
}
func TestWebhookSigningSecretNeverMarshalsAndCreateIsNotRetried(t *testing.T) {
	posts := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.Method == "GET" {
			return response(200, `{"items":[]}`), nil
		}
		posts++
		return nil, errors.New("lost create answer")
	})
	if _, err := c.EnsureWebhook(context.Background(), "https://ours.example/webhook"); err == nil || posts != 1 {
		t.Fatal("uncertain project webhook creation was retried")
	}
	raw, _ := json.Marshal(connector.Webhook{ID: "hook_one", Secret: "private-webhook-secret"})
	if strings.Contains(string(raw), "private-webhook-secret") {
		t.Fatal("signing secret escaped to JSON")
	}
}
