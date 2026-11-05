package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is planque's Worker over the browser worker's API (the worker
// is internal/browserworker, run next to planque). Every request carries
// the shared key; a session is one call's.
type Client struct {
	// URL is the worker, without a trailing slash.
	URL string
	// Key is the shared key the worker checks (its SHA-256 is configured there).
	Key string
	// HTTP reaches the worker: a private address, so not planque's egress
	// client. Nil, a default with a timeout.
	HTTP *http.Client
}

// KeyHeader carries the shared key: the same name the worker reads.
const KeyHeader = "Planque-Worker-Key"

func (c *Client) http() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 60 * time.Second}
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.URL, "/")+path, body)
	if err != nil {
		return err
	}
	req.Header.Set(KeyHeader, c.Key)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.http().Do(req)
	if err != nil {
		return fmt.Errorf("the browser worker couldn't be reached: %w", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 16<<20))
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		if e.Error == "" {
			e.Error = fmt.Sprintf("the browser worker answered %d", res.StatusCode)
		}
		return errors.New(e.Error)
	}
	if out != nil {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// Start begins a session on the worker.
func (c *Client) Start(ctx context.Context, start Start) (Session, error) {
	var out struct {
		Session string `json:"session"`
	}
	in := map[string]any{"grant": start.Grant, "sites": start.Sites}
	if len(start.State) > 0 {
		in["state"] = json.RawMessage(start.State)
	}
	if err := c.do(ctx, http.MethodPost, "/sessions", in, &out); err != nil {
		return nil, err
	}
	if out.Session == "" {
		return nil, errors.New("the browser worker gave no session")
	}
	return &remote{c: c, id: out.Session}, nil
}

type remote struct {
	c  *Client
	id string
}

func (r *remote) Open(ctx context.Context, url string) (Page, error) {
	var out struct {
		Page Page `json:"page"`
	}
	err := r.c.do(ctx, http.MethodPost, "/sessions/"+r.id+"/open", map[string]string{"url": url}, &out)
	return out.Page, err
}

func (r *remote) Do(ctx context.Context, step Step) (Page, error) {
	var out struct {
		Page Page `json:"page"`
	}
	err := r.c.do(ctx, http.MethodPost, "/sessions/"+r.id+"/step", map[string]any{"step": step}, &out)
	return out.Page, err
}

func (r *remote) Screenshot(ctx context.Context) ([]byte, error) {
	var out struct {
		PNG []byte `json:"png"`
	}
	err := r.c.do(ctx, http.MethodPost, "/sessions/"+r.id+"/screenshot", nil, &out)
	return out.PNG, err
}

func (r *remote) Close(ctx context.Context) ([]byte, error) {
	var out struct {
		State json.RawMessage `json:"state"`
	}
	err := r.c.do(ctx, http.MethodDelete, "/sessions/"+r.id, nil, &out)
	return out.State, err
}
