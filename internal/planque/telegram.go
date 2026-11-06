package planque

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/store"
)

// A Telegram bot is an API-key grant (kiff/planque#83): the owner gives the
// bot's token, made at BotFather, like any key. planque checks it with
// Telegram, sets the bot's webhook with a secret of its own, hands each
// message people send the bot to tafs as a door message, and sends the
// taf's answers with send_message, a governed call like any other. The
// Bot API has no OpenAPI description, so planque carries the one operation
// it offers; nothing is fetched.

// TelegramDescription is the api.url of a Telegram bot grant: named, never
// fetched (planque carries telegramAPI).
const TelegramDescription = "https://core.telegram.org/bots/api"

// telegramBase is the Bot API, unless the server names another (tests).
const telegramBase = "https://api.telegram.org"

// telegramTokenHeader is the description's key placement. A Telegram
// token goes in the request path, which callTelegram builds; the header is
// never sent.
const telegramTokenHeader = "Planque-Telegram-Token"

// telegramMax is the most characters (UTF-16 code units) one Telegram
// message takes.
const telegramMax = 4096

const telegramAPI = `{
 "openapi": "3.0.3",
 "info": {"title": "Telegram Bot API"},
 "security": [{"botToken": []}],
 "components": {"securitySchemes": {"botToken": {"type": "apiKey", "in": "header", "name": "Planque-Telegram-Token"}}},
 "paths": {
  "/sendMessage": {"post": {"operationId": "send_message", "summary": "Send a message in a chat",
   "description": "Posts text in a chat the bot is in, as the bot. Text over 4,096 characters is sent as several messages, the reply on the first.",
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["chat_id", "text"], "properties": {
    "chat_id": {"type": "string", "description": "The chat's id, as the door message gives it (a group's is negative)"},
    "text": {"type": "string", "description": "The message"},
    "message_thread_id": {"type": "integer", "description": "The forum topic, when the message came in one"},
    "reply_to_message_id": {"type": "integer", "description": "The message this answers"}
   }}}}}}}
 }
}`

// telegramAuth is where the description says the token goes.
var telegramAuth = openapi.Auth{In: "header", Name: telegramTokenHeader}

// isTelegram: g is a Telegram bot grant.
func isTelegram(g *Grant) bool {
	return g.Request.API != nil && g.Request.API.URL == TelegramDescription
}

func (s *Server) telegramAPIBase() string {
	if s.TelegramBase != "" {
		return s.TelegramBase
	}
	return telegramBase
}

// telegramText is what the owner reads for a Telegram bot grant: the token
// they give, the webhook it sets, what reaches the taf, and the tools.
func telegramText(q Request, held bool, account string, catalog ...CatalogTool) string {
	hook := "Giving it points the bot's webhook at planque, replacing any the bot had: messages people send the bot in a private chat, " +
		"and group messages that mention it, go to " + tafWords(q) + " as they arrive. Other messages are dropped."
	vault := "The token is kept in a sealed vault and is never shown again; the taf never sees it. It's sent only to Telegram."
	if len(q.Tools) == 0 {
		return "Connect a Telegram bot with its token from BotFather, for " + tafWords(q) + ". Telegram checks the token first.\n\n" +
			hook + " The taf can't answer until you approve its tools next. If they aren't approved within a day, the token is deleted.\n\n" +
			vault + " You can disconnect at any time, which deletes the webhook and the token."
	}
	var b strings.Builder
	if !held {
		b.WriteString("You will give a Telegram bot's token from BotFather.\n" + hook + "\n\n")
	}
	b.WriteString(toolLines(q, held, account, catalog...))
	b.WriteString("\n" + vault + " Each message is sent only with KIFF's permit for that exact call. You can disconnect at any time, which deletes the webhook and the token.")
	return b.String()
}

// TelegramBot is the bot a Telegram grant's token is, and the webhook
// planque set for it.
type TelegramBot struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	// Hook is the opaque id in the webhook's URL. Telegram keeps one
	// webhook per bot, so a newer grant for the same bot takes it over.
	Hook string `json:"hook"`
}

// telegramBotView is the bot as a grant's view shows it.
type telegramBotView struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// telegramHook is one webhook's record: whose it is, and its secret, sealed.
type telegramHook struct {
	Tenant string `json:"tenant"`
	Grant  string `json:"grant"`
	Bot    int64  `json:"bot"`
	Sealed []byte `json:"sealed"`
}

const telegramHookPK = "TELEGRAM#HOOKS"

func telegramHookContext(hook string) map[string]string {
	return map[string]string{"kind": "telegram_hook", "hook": hook}
}

// telegramTokenForm is a BotFather token's form; nothing else is put in a URL.
var telegramTokenForm = regexp.MustCompile(`^[0-9]{1,20}:[A-Za-z0-9_-]{1,200}$`)

// hideToken is text without token, as it is or URL-escaped.
func hideToken(text, token string) string {
	text = rest.HideKey(text, token)
	return rest.HideKey(rest.HideKey(text, url.PathEscape(token)), url.QueryEscape(token))
}

// telegramCall is one Bot API method with the token, and what Telegram
// answered: its result, or its description of the refusal.
func (s *Server) telegramCall(ctx context.Context, token, method string, body any) (json.RawMessage, int, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	ans, err := s.REST.Call(ctx, "POST", s.telegramAPIBase()+"/bot"+token+"/"+method, raw, openapi.Auth{}, token)
	if err != nil {
		// The token is in the URL, and some errors carry the URL: none
		// leaves here with it, whatever its form.
		var ce *rest.CallError
		if errors.As(err, &ce) {
			return nil, 0, &rest.CallError{Sent: ce.Sent, Err: errors.New(hideToken(ce.Err.Error(), token))}
		}
		return nil, 0, errors.New(hideToken(err.Error(), token))
	}
	var r struct {
		OK          bool            `json:"ok"`
		Result      json.RawMessage `json:"result"`
		Description string          `json:"description"`
	}
	if json.Unmarshal(ans.Body, &r) != nil || !r.OK {
		why := hideToken(r.Description, token)
		if why == "" {
			why = fmt.Sprintf("answered %d", ans.Status)
		}
		return nil, ans.Status, errors.New(why)
	}
	return r.Result, ans.Status, nil
}

// setTelegramHook checks token with Telegram and points the bot's webhook
// at a new hook of g's, with a secret only Telegram and planque hold. On
// any failure nothing is kept and the bot's webhook is as it was. With same,
// the token must be that bot's (a token replaced, gap 13 of kiff/tafs#476).
func (s *Server) setTelegramHook(ctx context.Context, g *Grant, token string, same *TelegramBot) (*TelegramBot, string, error) {
	if !telegramTokenForm.MatchString(token) {
		return nil, "that isn't a Telegram bot token's form (digits, a colon, then letters, digits, - or _): nothing was kept", nil
	}
	me, _, err := s.telegramCall(ctx, token, "getMe", map[string]any{})
	if err != nil {
		return nil, "Telegram didn't take that token (" + err.Error() + "): nothing was kept", nil
	}
	var bot struct {
		ID       int64  `json:"id"`
		IsBot    bool   `json:"is_bot"`
		Username string `json:"username"`
	}
	if json.Unmarshal(me, &bot) != nil || !bot.IsBot || bot.ID == 0 || bot.Username == "" {
		return nil, "that token isn't a Telegram bot's: nothing was kept", nil
	}
	if same != nil && bot.ID != same.ID {
		return nil, "that token is @" + bot.Username + "'s, not @" + same.Username + "'s: a new token must be the same bot's. Nothing was kept: the old one stays", nil
	}
	hook, secret := rand.Text(), rand.Text()
	sealed, err := s.Sealer.Seal(ctx, []byte(secret), telegramHookContext(hook))
	if err != nil {
		return nil, "", err
	}
	rec := telegramHook{Tenant: g.Tenant, Grant: g.ID, Bot: bot.ID, Sealed: sealed}
	if err := putDoc(ctx, s.KV, store.Item{PK: telegramHookPK, SK: hook, Ver: 1}, rec, store.Cond{NotExists: true}); err != nil {
		return nil, "", err
	}
	if _, _, err := s.telegramCall(ctx, token, "setWebhook", map[string]any{
		"url": s.BaseURL + "/telegram/" + hook, "secret_token": secret, "allowed_updates": []string{"message"}}); err != nil {
		_, _, _ = s.KV.Delete(ctx, telegramHookPK, hook)
		return nil, "Telegram didn't set the bot's webhook (" + err.Error() + "): nothing was kept", nil
	}
	return &TelegramBot{ID: bot.ID, Username: bot.Username, Hook: hook}, "", nil
}

// dropTelegramHook deletes g's webhook at Telegram, while its token is
// still held, and its record here.
func (s *Server) dropTelegramHook(ctx context.Context, g *Grant) []string {
	if g.Telegram == nil || g.TokenOf != "" {
		return nil
	}
	var problems []string
	if tok, _, err := s.openToken(ctx, g); err == nil {
		if _, _, err := s.telegramCall(ctx, tok.Access, "deleteWebhook", map[string]any{}); err != nil {
			problems = append(problems, "Telegram did not confirm deleting the bot's webhook: "+err.Error())
		}
	}
	_, _, _ = s.KV.Delete(ctx, telegramHookPK, g.Telegram.Hook)
	return problems
}

// telegramUpdate is the part of a Telegram update a door reads.
type telegramUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  *struct {
		MessageID       int64 `json:"message_id"`
		MessageThreadID int64 `json:"message_thread_id"`
		IsTopic         bool  `json:"is_topic_message"`
		Date            int64 `json:"date"`
		From            *struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"from"`
		Chat struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"chat"`
		Text     string `json:"text"`
		Entities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			User   *struct {
				ID int64 `json:"id"`
			} `json:"user"`
		} `json:"entities"`
	} `json:"message"`
}

// telegramWebhook takes a Telegram delivery for one bot. The secret header
// is checked before anything in the body is read; a message for the taf is
// kept and handed to tafs as the Slack door's are, once per update.
func (s *Server) telegramWebhook(w http.ResponseWriter, r *http.Request) {
	hook := r.PathValue("hook")
	var rec telegramHook
	it, err := getDoc(r.Context(), s.KV, telegramHookPK, hook, &rec)
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusUnauthorized, "the delivery's secret is not this webhook's")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	secret, err := s.Sealer.Open(r.Context(), rec.Sealed, telegramHookContext(it.SK))
	if err != nil {
		writeErr(w, http.StatusServiceUnavailable, "the webhook cannot be verified")
		return
	}
	defer clear(secret)
	got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if got == "" || subtle.ConstantTimeCompare([]byte(got), secret) != 1 {
		writeErr(w, http.StatusUnauthorized, "the delivery's secret is not this webhook's")
		return
	}
	raw, err := io.ReadAll(io.LimitReader(r.Body, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		writeErr(w, http.StatusBadRequest, "the update exceeds its limit")
		return
	}
	var u telegramUpdate
	if json.Unmarshal(raw, &u) != nil || u.UpdateID == 0 {
		writeErr(w, http.StatusBadRequest, "that isn't a Telegram update")
		return
	}
	g, err := s.loadGrant(r.Context(), rec.Tenant, rec.Grant)
	if errors.Is(err, errNotFound) {
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	// Telegram retries anything but a 2xx, so what isn't for the taf is
	// answered 200 and dropped: a superseded hook, a grant not connected,
	// an update that isn't a person's text to the bot.
	m, ok := telegramDoorMessage(g, it.SK, u)
	if !ok || s.telegramClosed(g, it.SK) {
		w.WriteHeader(http.StatusOK)
		return
	}
	body, _ := json.Marshal(m)
	retry, err := s.enqueueDoor(r.Context(), g, "telegram", it.SK, m.ID, "telegram", body)
	clear(body)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if retry {
		writeErr(w, http.StatusServiceUnavailable, "the message is pending delivery")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// telegramClosed: g takes no message through hook now.
func (s *Server) telegramClosed(g *Grant, hook string) bool {
	return !isTelegram(g) || g.Telegram == nil || g.Telegram.Hook != hook || g.Status == StatusDisconnected || g.Status == StatusFailed ||
		g.Status == StatusNeedsReconnect || !g.SignedIn || s.overdue(g)
}

// telegramDoorMessage is u as a door message for g's taf: a person's text
// in a private chat with the bot, or in a group where it names the bot.
// Its id is the same for Telegram's redelivery of the same update.
func telegramDoorMessage(g *Grant, hook string, u telegramUpdate) (door.TelegramMessage, bool) {
	m := u.Message
	if g.Telegram == nil || m == nil || m.From == nil || m.From.IsBot || m.From.ID <= 0 || strings.TrimSpace(m.Text) == "" || len(m.Text) > 128<<10 {
		return door.TelegramMessage{}, false
	}
	switch m.Chat.Type {
	case "private":
		if m.Chat.ID != m.From.ID {
			return door.TelegramMessage{}, false
		}
	case "group", "supergroup":
		if !mentions(m.Text, g.Telegram, m.Entities) {
			return door.TelegramMessage{}, false
		}
	default:
		return door.TelegramMessage{}, false
	}
	sum := sha256.Sum256([]byte("telegram\n" + g.Tenant + "\n" + g.ID + "\n" + hook + "\n" + strconv.FormatInt(u.UpdateID, 10)))
	out := door.TelegramMessage{V: "planque.door.v1", App: "telegram", ID: "evt_" + hex.EncodeToString(sum[:]), Tenant: g.Tenant, Agent: g.Request.Agent,
		Grant: g.ID, Provider: "telegram", Chat: strconv.FormatInt(m.Chat.ID, 10), Sender: strconv.FormatInt(m.From.ID, 10),
		MessageID: strconv.FormatInt(m.MessageID, 10), Name: strings.TrimSpace(m.From.FirstName + " " + m.From.LastName), Text: m.Text, Timestamp: m.Date}
	if m.IsTopic && m.MessageThreadID != 0 {
		out.Thread = strconv.FormatInt(m.MessageThreadID, 10)
	}
	return out, true
}

// mentions: text names bot, by @username or as a mention of its user.
// Telegram's offsets count UTF-16 code units.
func mentions(text string, bot *TelegramBot, entities []struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	User   *struct {
		ID int64 `json:"id"`
	} `json:"user"`
}) bool {
	units := utf16.Encode([]rune(text))
	for _, e := range entities {
		switch {
		case e.Type == "text_mention" && e.User != nil && e.User.ID == bot.ID:
			return true
		case e.Type == "mention" && e.Offset >= 0 && e.Length > 0 && e.Offset+e.Length <= len(units):
			if strings.EqualFold(string(utf16.Decode(units[e.Offset:e.Offset+e.Length])), "@"+bot.Username) {
				return true
			}
		}
	}
	return false
}

// splitTelegram cuts text into messages of at most telegramMax UTF-16
// units, at a line break where one is in the second half of a part, else
// at a space, else where it must.
func splitTelegram(text string) []string {
	var parts []string
	for {
		units := utf16.Encode([]rune(text))
		if len(units) <= telegramMax {
			return append(parts, text)
		}
		// The longest prefix of whole runes that fits.
		n, size := 0, 0
		for i, r := range text {
			w := utf16.RuneLen(r)
			if w < 0 {
				w = 1
			}
			if size+w > telegramMax {
				n = i
				break
			}
			size += w
		}
		cut := n
		if i := strings.LastIndex(text[:n], "\n"); i > n/2 {
			cut = i + 1
		} else if i := strings.LastIndex(text[:n], " "); i > n/2 {
			cut = i + 1
		}
		parts = append(parts, text[:cut])
		text = text[cut:]
	}
}

// callTelegram is one send_message call: one Telegram message, or several
// for a long text, the reply on the first. Nothing sent, it's not sent; a
// first part whose answer was lost, unknown; once a part went out, done,
// saying which parts did.
func (s *Server) callTelegram(ctx context.Context, g *Grant, tool GrantedTool, token string, args map[string]any) (string, any) {
	snap, err := s.loadSnapshot(ctx, g.Tenant, g.ID)
	if err != nil {
		return OpNotSent, toolResult("Not sent: planque couldn't read this connection's API description.", true)
	}
	var op *openapi.Operation
	for i := range snap.Operations {
		if snap.Operations[i].Name == tool.Tool {
			op = &snap.Operations[i]
		}
	}
	if op == nil || op.Name != "send_message" {
		return OpNotSent, toolResult("Not sent: the Telegram connection has no operation "+tool.Tool+".", true)
	}
	// The description's own check of the arguments, before anything goes.
	if _, _, _, err := op.Request(s.telegramAPIBase(), args); err != nil {
		return OpNotSent, toolResult("Refused by planque: "+err.Error()+". Nothing was sent.", true)
	}
	chat := fmt.Sprint(args["chat_id"])
	text, _ := args["text"].(string)
	if _, err := strconv.ParseInt(chat, 10, 64); err != nil || strings.TrimSpace(text) == "" {
		return OpNotSent, toolResult("Refused by planque: chat_id is a chat's numeric id and text isn't empty. Nothing was sent.", true)
	}
	parts := splitTelegram(text)
	var sent []string
	for i, part := range parts {
		body := map[string]any{"chat_id": chat, "text": part}
		if v, ok := args["message_thread_id"]; ok {
			body["message_thread_id"] = v
		}
		if v, ok := args["reply_to_message_id"]; ok && i == 0 {
			body["reply_parameters"] = map[string]any{"message_id": v, "allow_sending_without_reply": true}
		}
		res, status, err := s.telegramCall(ctx, token, "sendMessage", body)
		var ce *rest.CallError
		switch {
		case err == nil:
			var msg struct {
				MessageID int64 `json:"message_id"`
			}
			_ = json.Unmarshal(res, &msg)
			sent = append(sent, strconv.FormatInt(msg.MessageID, 10))
			continue
		case i == 0 && errors.As(err, &ce) && !ce.Sent:
			return OpNotSent, toolResult("Not sent: Telegram could not be reached.", true)
		case i == 0 && status >= 400 && status < 500:
			return OpDone, toolResult(fmt.Sprintf("Telegram refused it (%d): %s", status, err.Error()), true)
		case i == 0:
			s.log().WarnContext(ctx, "telegram answer lost", slog.String("grant", g.ID), slog.String("err", err.Error()))
			return OpUnknown, toolResult("Sent to Telegram, but its answer was lost: the message may or may not be in the chat. planque will not send it again.", true)
		}
		return OpDone, toolResult(fmt.Sprintf("Sent %d of %d parts (messages %s); part %d failed: %s. planque will not send the rest again.",
			len(sent), len(parts), strings.Join(sent, ", "), i+1, err.Error()), true)
	}
	out, _ := json.Marshal(map[string]any{"ok": true, "chat_id": chat, "message_ids": sent, "parts": len(parts)})
	return OpDone, toolResult(string(out), false)
}

// enqueueDoor keeps a door message, sealed, and tries to hand it over:
// once per id, whatever the provider redelivers. retry: it waits for the
// sweep.
func (s *Server) enqueueDoor(ctx context.Context, g *Grant, provider, trigger, id, kind string, body []byte) (bool, error) {
	key := store.Item{PK: deliveryPK, SK: id, Ver: 1, TTL: s.now().Add(7 * 24 * time.Hour).Unix()}
	sealed, err := s.Sealer.Seal(ctx, body, sealContext(kind, id))
	if err != nil {
		return false, err
	}
	record := eventDelivery{Tenant: g.Tenant, Grant: g.ID, Provider: provider, Trigger: trigger, Sealed: sealed, State: "pending", Kind: kind}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: "DOOR_EVENT#" + id, TTL: key.TTL}, map[string]string{"event": id}, store.Cond{}); err != nil {
		return false, err
	}
	stamp, _ := json.Marshal(pendingStamp{At: s.now().UnixMilli()})
	if err := s.KV.Put(ctx, store.Item{PK: pendingEventsPK, SK: id, Doc: stamp, TTL: key.TTL}, store.Cond{}); err != nil {
		return false, err
	}
	if err := putDoc(ctx, s.KV, key, record, store.Cond{NotExists: true}); err != nil && !errors.Is(err, store.ErrCondition) {
		return false, err
	}
	retry, err := s.deliverEvent(ctx, key.SK)
	return retry || err != nil, nil
}
