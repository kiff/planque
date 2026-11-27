package planque

import (
	"testing"

	"github.com/kiff/planque/internal/connector"
)

// A Slack write that mentions everyone is refused on both routes to Slack;
// an ordinary message, a read, and another app's write with the same words
// go through.
func TestASlackMessageToEveryoneIsRefusedOnEveryRoute(t *testing.T) {
	key := &Grant{Request: Request{API: &APIRequest{Base: "https://slack.com/api"}}}
	managed := &Grant{}
	other := &Grant{Request: Request{API: &APIRequest{Base: "https://api.example.com"}}}
	post := GrantedTool{ToolRequest: ToolRequest{Tool: "chat_postMessage", Name: "slack_post_message"}}
	read := GrantedTool{ToolRequest: ToolRequest{Tool: "conversations_history", Name: "slack_read_channel", ReadOnly: true}}
	send := GrantedTool{ToolRequest: ToolRequest{Tool: "SLACK_SEND_MESSAGE", Name: "send"},
		Connector: &connector.Tool{Toolkit: "slack", Slug: "SLACK_SEND_MESSAGE"}}
	bot := GrantedTool{ToolRequest: ToolRequest{Tool: "SLACKBOT_SEND_MESSAGE", Name: "bot"},
		Connector: &connector.Tool{Toolkit: "SlackBot", Slug: "SLACKBOT_SEND_MESSAGE"}}
	gmail := GrantedTool{ToolRequest: ToolRequest{Tool: "GMAIL_SEND_EMAIL", Name: "mail"},
		Connector: &connector.Tool{Toolkit: "gmail", Slug: "GMAIL_SEND_EMAIL"}}

	everyone := []map[string]any{
		{"channel": "C1", "text": "<!channel> the release is out"},
		{"channel": "C1", "text": "Heads up <!here|here>"},
		{"channel": "C1", "text": "<!EVERYONE> please read"},
		{"channel": "C1", "text": "ping <!subteam^S123|@support>"},
		{"channel": "C1", "text": "@channel the release is out"},
		{"channel": "C1", "text": "Thanks, (@here) all"},
		{"channel": "C1", "blocks": []any{map[string]any{"type": "section", "text": map[string]any{"text": "<!channel>"}}}},
		{"channel": "C1", "text": "ok", "attachments": []any{map[string]any{"fallback": "@everyone"}}},
	}
	for _, args := range everyone {
		for _, c := range []struct {
			g    *Grant
			tool GrantedTool
		}{{key, post}, {managed, send}, {managed, bot}} {
			if why := slackBroadcast(c.g, c.tool, args); why == "" {
				t.Errorf("%s with %v: sent, want refused", c.tool.Name, args)
			}
		}
	}

	ordinary := []map[string]any{
		{"channel": "C1", "text": "The release is out. Thanks, <@U123>."},
		{"channel": "C1", "text": "Write to support@here.example or ops@channel.example"},
		{"channel": "C1", "text": "the #here channel and the word everyone"},
	}
	for _, args := range ordinary {
		if why := slackBroadcast(key, post, args); why != "" {
			t.Errorf("%v: refused (%s), want sent", args, why)
		}
		if why := slackBroadcast(managed, send, args); why != "" {
			t.Errorf("managed %v: refused (%s), want sent", args, why)
		}
	}

	// A read is never refused for its words, and another app's write isn't Slack's.
	loud := map[string]any{"channel": "C1", "text": "<!channel>"}
	if why := slackBroadcast(key, read, loud); why != "" {
		t.Errorf("a read: refused (%s)", why)
	}
	if why := slackBroadcast(managed, gmail, loud); why != "" {
		t.Errorf("Gmail: refused (%s)", why)
	}
	if why := slackBroadcast(other, post, loud); why != "" {
		t.Errorf("a key grant to another API: refused (%s)", why)
	}
}

// A managed Slack tool is checked whatever its provider marks it, so a send
// wrongly marked read-only can't notify everyone. A managed read without a
// broadcast mention goes through; a key grant's reviewed read stays a read.
func TestAManagedSlackToolIsCheckedWhateverItsMark(t *testing.T) {
	managed := &Grant{}
	loud := map[string]any{"channel": "C1", "text": "<!channel> the release is out"}
	for _, c := range []struct{ toolkit, slug string }{
		{"slack", "SLACK_SEND_MESSAGE"}, {"slack", "SLACK_CHAT_POST_MESSAGE"},
		{"slack", "SLACK_UPDATE_MESSAGE"}, {"slackbot", "SLACKBOT_SEND_MESSAGE"},
	} {
		marked := GrantedTool{ToolRequest: ToolRequest{Tool: c.slug, Name: "send", ReadOnly: true},
			Connector: &connector.Tool{Toolkit: c.toolkit, Slug: c.slug, ReadOnly: true}}
		if why := slackBroadcast(managed, marked, loud); why == "" {
			t.Errorf("%s marked read-only: sent, want refused", c.slug)
		}
	}
	read := GrantedTool{ToolRequest: ToolRequest{Tool: "SLACK_SEARCH_MESSAGES", Name: "read", ReadOnly: true},
		Connector: &connector.Tool{Toolkit: "slack", Slug: "SLACK_SEARCH_MESSAGES", ReadOnly: true}}
	if why := slackBroadcast(managed, read, map[string]any{"query": "release notes"}); why != "" {
		t.Errorf("a managed read: refused (%s)", why)
	}
	key := &Grant{Request: Request{API: &APIRequest{Base: "https://slack.com/api"}}}
	reviewed := GrantedTool{ToolRequest: ToolRequest{Tool: "search.messages", Name: "search", ReadOnly: true}}
	if why := slackBroadcast(key, reviewed, map[string]any{"query": "@here"}); why != "" {
		t.Errorf("a key grant's reviewed read: refused (%s)", why)
	}
}
