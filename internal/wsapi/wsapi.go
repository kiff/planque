// Package wsapi posts to and drops the connections of an API Gateway
// WebSocket API, through its management endpoint
// (https://{api}.execute-api.{region}.amazonaws.com/{stage}/@connections/{id}),
// signed with SigV4 and the standard library only.
package wsapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kiff/planque/internal/awsapi"
)

// ErrGone: the connection is closed; nothing reaches it any more.
var ErrGone = errors.New("wsapi: the connection is gone")

// Poster sends a message to one connection, or drops it.
type Poster interface {
	Post(ctx context.Context, conn string, msg []byte) error
	Drop(ctx context.Context, conn string) error
}

// API is one WebSocket API stage's management endpoint.
type API struct {
	Endpoint string // https://{api}.execute-api.{region}.amazonaws.com/{stage}
	Region   string
	Creds    func() (awsapi.Credentials, error)
	HTTP     *http.Client
	Now      func() time.Time
}

func (a API) do(ctx context.Context, method, conn string, body []byte) error {
	creds, err := a.Creds()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, a.Endpoint+"/@connections/"+url.PathEscape(conn), bytes.NewReader(body))
	if err != nil {
		return err
	}
	now := time.Now
	if a.Now != nil {
		now = a.Now
	}
	awsapi.Sign(req, body, creds, "execute-api", a.Region, now())
	hc := a.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 5 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusGone:
		return ErrGone
	case resp.StatusCode >= 300:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("wsapi %s: %d %s", method, resp.StatusCode, b)
	}
	return nil
}

// Post sends msg to conn.
func (a API) Post(ctx context.Context, conn string, msg []byte) error {
	return a.do(ctx, http.MethodPost, conn, msg)
}

// Drop closes conn. One already gone is dropped.
func (a API) Drop(ctx context.Context, conn string) error {
	if err := a.do(ctx, http.MethodDelete, conn, nil); err != nil && !errors.Is(err, ErrGone) {
		return err
	}
	return nil
}
