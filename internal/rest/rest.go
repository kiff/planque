// Package rest is how planque reaches a vendor's REST API for an API-key
// grant (kiff/planque#64): it fetches the vendor's OpenAPI description once,
// and makes each call the relay builds. Redirects are never followed: a call
// goes to the grant's server and nowhere else.
package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kiff/planque/internal/openapi"
)

// MaxAnswer is how much of a vendor's answer a call keeps.
const MaxAnswer = 256 << 10

// Client reaches vendors' REST APIs.
type Client struct {
	// HTTP reaches vendors: in production, only at public addresses.
	HTTP *http.Client
	// AllowHTTP lets tests use a plain-http fake vendor.
	AllowHTTP bool
}

func (c *Client) http() *http.Client {
	base := c.HTTP
	if base == nil {
		base = &http.Client{Timeout: 30 * time.Second}
	}
	cp := *base
	// A redirect is the vendor's answer, never followed: the key goes to the
	// grant's server only.
	cp.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &cp
}

// CheckURL refuses a URL planque won't reach: https only (http in tests),
// with a host, and no user info.
func (c *Client) CheckURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" || u.User != nil || u.Fragment != "" {
		return fmt.Errorf("%s isn't a URL planque can reach", raw)
	}
	if u.Scheme != "https" && !(c.AllowHTTP && u.Scheme == "http") {
		return fmt.Errorf("%s isn't an https URL", raw)
	}
	return nil
}

// Fetch reads a description, at most openapi.MaxBytes.
func (c *Client) Fetch(ctx context.Context, raw string) ([]byte, error) {
	if err := c.CheckURL(raw); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.http().Do(req)
	if err != nil {
		return nil, fmt.Errorf("it couldn't be read: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("it answered %d", res.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, openapi.MaxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("it couldn't be read: %w", err)
	}
	if len(b) > openapi.MaxBytes {
		return nil, &openapi.Refused{Why: "it is over 2 MB"}
	}
	return b, nil
}

// Answer is a vendor's answer to a call.
type Answer struct {
	Status int
	Body   []byte // at most MaxAnswer
	Cut    bool   // the answer was longer and was cut
}

// CallError is a call that has no answer. Sent says whether the request
// may have reached the vendor: a call that may have is never sent again.
type CallError struct {
	Sent bool
	Err  error
}

func (e *CallError) Error() string { return e.Err.Error() }
func (e *CallError) Unwrap() error { return e.Err }

// Call makes one request, with the key where auth says.
func (c *Client) Call(ctx context.Context, method, target string, body []byte, auth openapi.Auth, key string) (Answer, error) {
	if err := c.CheckURL(target); err != nil {
		return Answer{}, &CallError{Err: err}
	}
	if auth.In == "query" {
		u, _ := url.Parse(target)
		q := u.Query()
		q.Set(auth.Name, key)
		u.RawQuery = q.Encode()
		target = u.String()
	}
	var rd io.Reader
	if body != nil {
		rd = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, rd)
	if err != nil {
		return Answer{}, &CallError{Err: err}
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth.In == "header" {
		req.Header.Set(auth.Name, auth.Prefix+key)
	}
	res, err := c.http().Do(req)
	if err != nil {
		// A request that failed before a connection was made never reached
		// the vendor; any later failure may have.
		return Answer{}, &CallError{Sent: !notConnected(err), Err: err}
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, MaxAnswer+1))
	if err != nil {
		return Answer{}, &CallError{Sent: true, Err: err}
	}
	a := Answer{Status: res.StatusCode, Body: b}
	if len(b) > MaxAnswer {
		a.Body, a.Cut = b[:MaxAnswer], true
	}
	return a, nil
}

// notConnected: err happened before planque connected to the vendor (the
// name didn't resolve, the address was refused, the dial failed).
func notConnected(err error) bool {
	var ue *url.Error
	if errors.As(err, &ue) && ue.Timeout() {
		return false // a timeout may come after the request was written
	}
	msg := err.Error()
	for _, s := range []string{"no such host", "connection refused", "not a public", "blocked", "dial tcp"} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}
