package planque

import (
	"net/url"
	"regexp"
	"strings"
)

// A Slack message that mentions @channel, @here, @everyone or a user group
// notifies every member of that channel, workspace or group. One call to a
// "single" send then reaches everyone, which a Card's count of calls can't
// see. So planque refuses a Slack write whose arguments carry such a
// mention, on every route to Slack (a key to its Web API, or the managed
// app), before anything is sent. A taf that needs to reach everyone asks the
// person, who posts it.

// CodeBroadcast is the refusal of a Slack write that would notify everyone.
const CodeBroadcast = "broadcast_mention"

// slackBroadcastMention matches Slack's own mention syntax (<!channel>,
// <!here>, <!everyone>, <!subteam^ID>, with or without a label after "|"),
// and the plain words Slack turns into one when a message is sent with
// link_names or parse=full (@channel, @here, @everyone). Case doesn't count.
var slackBroadcastMention = regexp.MustCompile(`(?i)<!(channel|here|everyone|subteam\^)|(^|[^\w@.])@(channel|here|everyone)\b`)

// isSlack says t, on g, is a call to Slack: the managed app's tool (its
// toolkit slack or slackbot), or a key grant to Slack's Web API (slack.com).
func isSlack(g *Grant, t GrantedTool) bool {
	if t.Connector != nil {
		k := strings.ToLower(t.Connector.Toolkit)
		return k == "slack" || k == "slackbot"
	}
	if g != nil && g.Request.API != nil {
		u, err := url.Parse(g.Request.API.Base)
		return err == nil && strings.EqualFold(u.Hostname(), "slack.com")
	}
	return false
}

// slackRead says t is a read, so its words can't notify anyone. Only a key
// grant's tool can be: its reviewed entry's mark says so. A managed app's tool
// is never taken as a read on its provider's mark, which can be wrong (a send
// tagged read-only, or a name like CHAT_POST_MESSAGE that no verb list
// places), so every managed Slack call is checked. A read whose own words
// carry a broadcast mention is refused too; it can be asked without one.
func slackRead(t GrantedTool) bool {
	return t.Connector == nil && t.ReadOnly
}

// slackBroadcast is why a call to t with args is refused, or "": a Slack tool
// that isn't a read whose arguments, at any depth, mention everyone.
func slackBroadcast(g *Grant, t GrantedTool, args map[string]any) string {
	if !isSlack(g, t) || slackRead(t) {
		return ""
	}
	if m := broadcastIn(args); m != "" {
		return "a Slack message that mentions " + m + " notifies everyone there, so the taf doesn't send it; the person can post it themselves"
	}
	return ""
}

// broadcastIn is the first broadcast mention in any string within v, or "".
func broadcastIn(v any) string {
	switch x := v.(type) {
	case string:
		if loc := slackBroadcastMention.FindStringIndex(x); loc != nil {
			return strings.TrimLeft(x[loc[0]:loc[1]], " \t\n\r([{,;:!?\"'")
		}
	case map[string]any:
		for _, e := range x {
			if m := broadcastIn(e); m != "" {
				return m
			}
		}
	case []any:
		for _, e := range x {
			if m := broadcastIn(e); m != "" {
				return m
			}
		}
	}
	return ""
}
