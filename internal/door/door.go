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
	V         string `json:"v"`
	ID        string `json:"id"`
	Tenant    string `json:"tenant"`
	Agent     string `json:"agent"`
	Grant     string `json:"grant"`
	Provider  string `json:"provider"`
	Account   string `json:"account"`
	Trigger   string `json:"trigger"`
	Slug      string `json:"trigger_slug"`
	Channel   string `json:"channel"`
	Sender    string `json:"sender"`
	Text      string `json:"text"`
	Thread    string `json:"thread,omitempty"`
	Timestamp string `json:"timestamp"`
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
	u, err := url.Parse(c.URL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || c.Secret == nil {
		return Result{}, ErrUnavailable
	}
	secret, err := c.Secret(ctx)
	if err != nil || secret == "" {
		return Result{}, ErrUnavailable
	}
	raw, err := json.Marshal(m)
	if err != nil || len(raw) > 256<<10 || m.V != "planque.door.v1" || m.ID == "" {
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
	req.Header.Set("Planque-Event-Id", m.ID)
	req.Header.Set("Planque-Event-Timestamp", ts)
	req.Header.Set("Planque-Event-Signature", connector.SignWebhook(secret, m.ID, ts, raw))
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
	// The receiver deduplicates m.ID. Only 5xx/no-answer deliveries may be
	// retried; authorization failures, closed doors and redirects are terminal.
	return Result{Taken: response.StatusCode == http.StatusAccepted, Retry: response.StatusCode >= 500, Status: response.StatusCode}, nil
}
