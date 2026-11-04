package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"unicode/utf16"

	"github.com/kiff/planque/internal/door"
)

// telegramFake is the Bot API for one bot, and the tafs receiver its door
// messages go to.
type telegramFake struct {
	*httptest.Server
	token   string
	mu      sync.Mutex
	hookURL string
	secret  string
	deletes int
	sends   []map[string]any
	// failSend makes the nth sendMessage (1-based) fail with Telegram's 400.
	failSend  int
	delivered []map[string]any
}

const telegramToken = "123456:synthetic-bot-token"

func newTelegramFake(t *testing.T, e *env) *telegramFake {
	f := &telegramFake{token: telegramToken}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		token, method, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/bot"), "/")
		if token != f.token {
			w.WriteHeader(401)
			_, _ = io.WriteString(w, `{"ok":false,"error_code":401,"description":"Unauthorized"}`)
			return
		}
		var body map[string]any
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		switch method {
		case "getMe":
			_, _ = io.WriteString(w, `{"ok":true,"result":{"id":777,"is_bot":true,"first_name":"Shop","username":"shop_help_bot"}}`)
		case "setWebhook":
			f.hookURL, _ = body["url"].(string)
			f.secret, _ = body["secret_token"].(string)
			if a, _ := body["allowed_updates"].([]any); len(a) != 1 || a[0] != "message" {
				t.Errorf("allowed_updates: %v", body["allowed_updates"])
			}
			_, _ = io.WriteString(w, `{"ok":true,"result":true}`)
		case "deleteWebhook":
			f.deletes++
			_, _ = io.WriteString(w, `{"ok":true,"result":true}`)
		case "sendMessage":
			f.sends = append(f.sends, body)
			if len(f.sends) == f.failSend {
				w.WriteHeader(400)
				_, _ = io.WriteString(w, `{"ok":false,"error_code":400,"description":"Bad Request: chat not found `+f.token+`"}`)
				return
			}
			_, _ = io.WriteString(w, `{"ok":true,"result":{"message_id":`+string(rune('0'+len(f.sends)))+`}}`)
		default:
			t.Errorf("unexpected Bot API method %s", method)
		}
	}))
	t.Cleanup(f.Server.Close)
	e.s.TelegramBase = f.URL
	e.s.DoorSink = &door.HTTP{URL: "https://receiver.example/v1/doors/planque", Secret: func(context.Context) (string, error) { return "synthetic-handoff", nil },
		Client: &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
			var m map[string]any
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &m)
			f.mu.Lock()
			f.delivered = append(f.delivered, m)
			f.mu.Unlock()
			return &http.Response{StatusCode: 202, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		})}}
	e.kiff.Cards["k_taf"]["send_message"] = 0
	return f
}

// deliver is Telegram posting update to the webhook it was given, with
// secret in its header.
func (f *telegramFake) deliver(t *testing.T, url, secret string, update map[string]any) int {
	t.Helper()
	raw, _ := json.Marshal(update)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(raw))
	if secret != "" {
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

func (f *telegramFake) seen() ([]map[string]any, []map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]map[string]any(nil), f.sends...), append([]map[string]any(nil), f.delivered...)
}

func telegramMessage(update int64, chat map[string]any, from map[string]any, text string, extra map[string]any) map[string]any {
	m := map[string]any{"message_id": 42, "date": 1800000000, "chat": chat, "from": from, "text": text}
	for k, v := range extra {
		m[k] = v
	}
	return map[string]any{"update_id": update, "message": m}
}

var (
	ana     = map[string]any{"id": 5001, "is_bot": false, "first_name": "Ana", "last_name": "Ruiz", "username": "ana_r"}
	anaChat = map[string]any{"id": 5001, "type": "private"}
	group   = map[string]any{"id": -100200300, "type": "supergroup"}
)

// The whole Telegram connection: the token checked and its webhook set with
// a secret; only a person's text to the bot reaches tafs, once per update;
// the taf answers with send_message through KIFF; disconnecting deletes the
// webhook.
func TestATelegramBotEndToEnd(t *testing.T) {
	e := newEnv(t)
	f := newTelegramFake(t, e)
	e.enrol()
	ask := map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []map[string]any{}, "api": map[string]any{"url": TelegramDescription}}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", ask, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	// The owner reads what giving the token does, never planque's placement of it.
	if txt := out["text"].(string); !strings.Contains(txt, "points the bot's webhook at planque") || strings.Contains(txt, telegramTokenHeader) || strings.Contains(txt, "nothing is changed") {
		t.Fatalf("text: %s", txt)
	}

	// A token Telegram refuses keeps nothing and sets no webhook.
	if st, out := e.giveKey(id, digest, "999:not-a-token"); st != 400 || !strings.Contains(out["error"].(string), "Telegram didn't take that token") {
		t.Fatalf("bad token: %d %v", st, out)
	}
	if hooks, _ := e.kv.Query(context.Background(), telegramHookPK, ""); len(hooks) != 0 || f.hookURL != "" {
		t.Fatal("a refused token left a webhook")
	}
	if st, out := e.giveKey(id, digest, f.token); st != 200 || out["status"] != StatusListed {
		t.Fatalf("token: %d %v", st, out)
	}
	if !strings.HasPrefix(f.hookURL, e.url+"/telegram/") || len(f.secret) < 20 {
		t.Fatalf("webhook: %q", f.hookURL)
	}
	// Neither the token nor the webhook's secret is kept in the clear.
	for _, pk := range []string{tenantPK(tenant), telegramHookPK} {
		items, _ := e.kv.Query(context.Background(), pk, "")
		for _, it := range items {
			if bytes.Contains(it.Doc, []byte(f.token)) || bytes.Contains(it.Doc, []byte(f.secret)) {
				t.Fatalf("a secret in the clear in %s", it.SK)
			}
		}
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []map[string]any{{"tool": "send_message", "name": "send_message", "title": "Answer in Telegram"}}}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tools: %d %v", st, out)
	}
	if txt := out["text"].(string); !strings.Contains(txt, "Answer in Telegram") || !strings.Contains(txt, "deletes the webhook") || strings.Contains(txt, telegramTokenHeader) {
		t.Fatalf("tools text: %s", txt)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}

	// Deliveries without this webhook's secret are refused before the body is read.
	hello := telegramMessage(9001, anaChat, ana, "Hello, is my order shipped?", nil)
	for _, c := range []struct{ url, secret string }{{f.hookURL, ""}, {f.hookURL, "wrong-secret"}, {e.url + "/telegram/" + "unknownhook", f.secret}} {
		if st := f.deliver(t, c.url, c.secret, hello); st != 401 {
			t.Fatalf("%s %q: %d", c.url, c.secret, st)
		}
	}
	if _, d := f.seen(); len(d) != 0 {
		t.Fatalf("delivered without the secret: %v", d)
	}
	// A private message reaches tafs, in the agreed shape; Telegram's
	// redelivery of the same update doesn't reach it twice.
	for range 2 {
		if st := f.deliver(t, f.hookURL, f.secret, hello); st != 200 {
			t.Fatalf("private: %d", st)
		}
	}
	_, d := f.seen()
	if len(d) != 1 {
		t.Fatalf("deliveries: %v", d)
	}
	want := map[string]any{"v": "planque.door.v1", "app": "telegram", "tenant": tenant, "agent": "agt_taf", "grant": id, "provider": "telegram",
		"chat": "5001", "sender": "5001", "message_id": "42", "name": "Ana Ruiz", "text": "Hello, is my order shipped?", "timestamp": float64(1800000000)}
	for k, v := range want {
		if d[0][k] != v {
			t.Fatalf("%s: %v, want %v (%v)", k, d[0][k], v, d[0])
		}
	}
	if _, ok := d[0]["thread"]; ok || !strings.HasPrefix(d[0]["id"].(string), "evt_") {
		t.Fatalf("%v", d[0])
	}
	// Not for the taf: a bot's message, a group message that doesn't name the
	// bot, a channel post, a message with no text, a private chat not the
	// sender's. Each is answered 200 so Telegram doesn't retry it.
	drops := []map[string]any{
		telegramMessage(9002, anaChat, map[string]any{"id": 888, "is_bot": true, "first_name": "Other"}, "beep", nil),
		telegramMessage(9003, group, ana, "anyone around?", nil),
		telegramMessage(9004, group, ana, "@other_bot hi", map[string]any{"entities": []map[string]any{{"type": "mention", "offset": 0, "length": 10}}}),
		telegramMessage(9005, map[string]any{"id": -100999, "type": "channel"}, ana, "@shop_help_bot hi", map[string]any{"entities": []map[string]any{{"type": "mention", "offset": 0, "length": 14}}}),
		telegramMessage(9006, anaChat, ana, "", nil),
		telegramMessage(9007, map[string]any{"id": 6002, "type": "private"}, ana, "hi", nil),
	}
	for _, u := range drops {
		if st := f.deliver(t, f.hookURL, f.secret, u); st != 200 {
			t.Fatalf("drop: %d", st)
		}
	}
	if _, d := f.seen(); len(d) != 1 {
		t.Fatalf("delivered what isn't for the taf: %v", d[1:])
	}
	// A group message naming the bot (after an emoji: offsets are UTF-16), in a forum topic.
	text := "👋 @Shop_Help_Bot where is order 12?"
	mention := telegramMessage(9008, group, ana, text, map[string]any{"message_thread_id": 77, "is_topic_message": true,
		"entities": []map[string]any{{"type": "mention", "offset": 3, "length": 14}}})
	if st := f.deliver(t, f.hookURL, f.secret, mention); st != 200 {
		t.Fatalf("mention: %d", st)
	}
	if _, d := f.seen(); len(d) != 2 || d[1]["chat"] != "-100200300" || d[1]["thread"] != "77" || d[1]["text"] != text {
		t.Fatalf("group mention: %v", d)
	}

	// The taf answers with send_message, under its Card, in the topic, replying.
	ctx := context.Background()
	res, err := e.kiff.Call(ctx, "k_taf", "send_message", map[string]any{"chat_id": "-100200300", "text": "It left today.", "message_thread_id": 77, "reply_to_message_id": 42})
	if err != nil || res.Decision != "allowed" {
		t.Fatalf("%+v %v", res, err)
	}
	sends, _ := f.seen()
	if txt, isErr := resultText(t, res.Result); isErr || len(sends) != 1 || sends[0]["chat_id"] != "-100200300" || sends[0]["text"] != "It left today." ||
		sends[0]["message_thread_id"] != float64(77) || sends[0]["reply_parameters"].(map[string]any)["message_id"] != float64(42) {
		t.Fatalf("%s %v", txt, sends)
	}
	// A long answer goes as several messages, the reply on the first only.
	long := strings.Repeat("a", 4000) + "\n" + strings.Repeat("b", 4000) + "\n" + strings.Repeat("c", 100)
	res, _ = e.kiff.Call(ctx, "k_taf", "send_message", map[string]any{"chat_id": "5001", "text": long, "reply_to_message_id": 42})
	sends, _ = f.seen()
	if txt, isErr := resultText(t, res.Result); isErr || len(sends) != 4 || !strings.Contains(txt, `"parts":3`) ||
		sends[1]["reply_parameters"] == nil || sends[2]["reply_parameters"] != nil || sends[3]["reply_parameters"] != nil ||
		sends[1]["text"].(string)+sends[2]["text"].(string)+sends[3]["text"].(string) != long {
		t.Fatalf("%s %d sends", txt, len(sends))
	}
	// A refusal after the first part says what went out; the token is never in it.
	f.mu.Lock()
	f.failSend = 6
	f.mu.Unlock()
	res, _ = e.kiff.Call(ctx, "k_taf", "send_message", map[string]any{"chat_id": "5001", "text": strings.Repeat("x", 5000)})
	if txt, isErr := resultText(t, res.Result); !isErr || !strings.Contains(txt, "Sent 1 of 2 parts") || strings.Contains(txt, f.token) {
		t.Fatalf("part failed: %s", txt)
	}
	// A chat id that isn't one: refused by planque, nothing sent.
	before, _ := f.seen()
	res, _ = e.kiff.Call(ctx, "k_taf", "send_message", map[string]any{"chat_id": "@ana_r", "text": "hi"})
	if txt, isErr := resultText(t, res.Result); !isErr || !strings.Contains(txt, "Nothing was sent") {
		t.Fatalf("username as chat: %s", txt)
	}
	if after, _ := f.seen(); len(after) != len(before) {
		t.Fatal("sent to a username")
	}

	// Disconnecting deletes the webhook at Telegram; its deliveries stop here.
	if st, out := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if f.deletes != 1 {
		t.Fatalf("deleteWebhook calls: %d", f.deletes)
	}
	if st := f.deliver(t, f.hookURL, f.secret, telegramMessage(9009, anaChat, ana, "still there?", nil)); st != 401 {
		t.Fatalf("after disconnect: %d", st)
	}
}

// A delivery for a grant that isn't connected now is taken (Telegram stops
// retrying) and not handed on, and the queued one is dropped.
func TestATelegramMessageForAGrantNotConnectedIsDropped(t *testing.T) {
	e := newEnv(t)
	f := newTelegramFake(t, e)
	e.enrol()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []map[string]any{}, "api": map[string]any{"url": TelegramDescription}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, _ := e.giveKey(id, out["digest"].(string), f.token); st != 200 {
		t.Fatal(st)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	g.Status = StatusNeedsReconnect
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	if st := f.deliver(t, f.hookURL, f.secret, telegramMessage(1, anaChat, ana, "hi", nil)); st != 200 {
		t.Fatal(st)
	}
	if _, d := f.seen(); len(d) != 0 {
		t.Fatalf("delivered for a grant that isn't connected: %v", d)
	}
}

func TestSplitTelegramKeepsEachPartUnderTheLimitInUTF16(t *testing.T) {
	cases := []string{
		strings.Repeat("😀", 3000),                                    // 2 units each
		strings.Repeat("word ", 2000),                                // spaces
		strings.Repeat("a", 3000) + "\n" + strings.Repeat("b", 3000), // a line break
		strings.Repeat("a", 4096),
		"short",
	}
	for _, text := range cases {
		parts := splitTelegram(text)
		if strings.Join(parts, "") != text {
			t.Fatal("parts don't make the text")
		}
		for _, p := range parts {
			if n := len(utf16.Encode([]rune(p))); n > telegramMax || n == 0 {
				t.Fatalf("a part of %d units", n)
			}
		}
	}
	if parts := splitTelegram(strings.Repeat("a", 3000) + "\n" + strings.Repeat("b", 3000)); len(parts) != 2 || !strings.HasSuffix(parts[0], "\n") {
		t.Fatalf("not cut at the line break: %d parts", len(parts))
	}
}
