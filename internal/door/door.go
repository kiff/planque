// Package door delivers authenticated event messages to one configured tafs
// receiver. Receiving an event never grants a vendor action or a reply.
package door

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kiff/planque/internal/connector"
)

type Message struct {
	V        string `json:"v"`
	ID       string `json:"id"`
	Tenant   string `json:"tenant"`
	Agent    string `json:"agent"`
	Grant    string `json:"grant"`
	Provider string `json:"provider"`
	Account  string `json:"account"`
	Trigger  string `json:"trigger"`
	Slug     string `json:"trigger_slug"`
	Channel  string `json:"channel"`
	Sender   string `json:"sender"`
	Text     string `json:"text"`
	// TS is the message's own Slack ts, verbatim: a reply to a root message
	// goes under it, as Thread is the parent of a message already in a thread.
	TS        string `json:"ts"`
	Thread    string `json:"thread,omitempty"`
	Timestamp string `json:"timestamp"`
}

// AppEvent is any connected app's event, for a taf that wakes on it
// (planque.event.v1): the subscription the owner approved, and the event's
// data as the provider sent it. Receiving one never grants an action.
type AppEvent struct {
	V            string          `json:"v"`
	ID           string          `json:"id"`
	Tenant       string          `json:"tenant"`
	Agent        string          `json:"agent"`
	Grant        string          `json:"grant"`
	Provider     string          `json:"provider"`
	Account      string          `json:"account"`
	Subscription string          `json:"subscription"`
	Event        string          `json:"event"`
	Payload      json.RawMessage `json:"payload"`
	Timestamp    string          `json:"timestamp"`
}

// TelegramMessage is a message to a Telegram bot connected here, handed to
// tafs as a door message (planque.door.v1, "app": "telegram"). Ids are
// Telegram's numeric ids as strings; a @username is never a key.
type TelegramMessage struct {
	V        string `json:"v"`
	App      string `json:"app"`
	ID       string `json:"id"`
	Tenant   string `json:"tenant"`
	Agent    string `json:"agent"`
	Grant    string `json:"grant"`
	Provider string `json:"provider"`
	Chat     string `json:"chat"`
	Sender   string `json:"sender"`
	// Thread is the forum topic's message_thread_id, only in a forum group.
	Thread    string `json:"thread,omitempty"`
	MessageID string `json:"message_id"`
	Name      string `json:"name,omitempty"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

// TelegramSink takes Telegram messages, as Sink takes Slack's.
type TelegramSink interface {
	DeliverTelegram(context.Context, TelegramMessage) (Result, error)
}

// EventSink takes app events, as Sink takes door messages.
type EventSink interface {
	DeliverEvent(context.Context, AppEvent) (Result, error)
}

type Result struct {
	Taken  bool
	Retry  bool
	Status int
}
type Sink interface {
	Deliver(context.Context, Message) (Result, error)
}
type HTTP struct {
	URL    string
	Secret func(context.Context) (string, error)
	Client *http.Client
	Now    func() time.Time
}

var ErrUnavailable = errors.New("the door receiver is unavailable")

func (c *HTTP) Deliver(ctx context.Context, m Message) (Result, error) {
	if m.V != "planque.door.v1" || m.ID == "" {
		return Result{}, ErrUnavailable
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return Result{}, ErrUnavailable
	}
	return c.post(ctx, m.ID, raw)
}

// DeliverTelegram sends a Telegram message as Deliver sends Slack's: to the
// same receiver, signed with the same secret, retried on the same answers.
func (c *HTTP) DeliverTelegram(ctx context.Context, m TelegramMessage) (Result, error) {
	if m.V != "planque.door.v1" || m.App != "telegram" || m.ID == "" {
		return Result{}, ErrUnavailable
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return Result{}, ErrUnavailable
	}
	return c.post(ctx, m.ID, raw)
}

// DeliverEvent sends an app event as Deliver sends a door message: signed
// with the same secret, retried on the same answers.
func (c *HTTP) DeliverEvent(ctx context.Context, e AppEvent) (Result, error) {
	if e.V != "planque.event.v1" || e.ID == "" || !json.Valid(e.Payload) {
		return Result{}, ErrUnavailable
	}
	raw, err := json.Marshal(e)
	if err != nil {
		return Result{}, ErrUnavailable
	}
	return c.post(ctx, e.ID, raw)
}

func (c *HTTP) post(ctx context.Context, id string, raw []byte) (Result, error) {
	u, err := url.Parse(c.URL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || c.Secret == nil {
		return Result{}, ErrUnavailable
	}
	if len(raw) > 256<<10 {
		return Result{}, ErrUnavailable
	}
	secret, err := c.Secret(ctx)
	if err != nil || secret == "" {
		return Result{}, ErrUnavailable
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewReader(raw))
	if err != nil {
		return Result{}, ErrUnavailable
	}
	now := time.Now()
	if c.Now != nil {
		now = c.Now()
	}
	ts := strconv.FormatInt(now.Unix(), 10)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Planque-Event-Id", id)
	req.Header.Set("Planque-Event-Timestamp", ts)
	req.Header.Set("Planque-Event-Signature", connector.SignWebhook(secret, id, ts, raw))
	hc := c.Client
	if hc == nil {
		hc = &http.Client{Timeout: 10 * time.Second}
	}
	client := *hc
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := client.Do(req)
	if err != nil {
		return Result{Retry: true}, nil
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
	// The receiver deduplicates the id. Only 5xx/no-answer deliveries may be
	// retried; authorization failures, closed doors and redirects are terminal.
	return Result{Taken: response.StatusCode == http.StatusAccepted, Retry: response.StatusCode >= 500, Status: response.StatusCode}, nil
}
