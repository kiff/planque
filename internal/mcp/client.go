// Package mcp is the slice of the Model Context Protocol planque speaks
// over Streamable HTTP: a client for the vendor's MCP server, and the
// server side of planque's relay, which KIFF's gateway calls.
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Version is the protocol revision planque asks for and serves.
const Version = "2025-06-18"

// Tool is a tool as a server lists it.
type Tool struct {
	Name        string          `json:"name"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"inputSchema"`
	Annotations json.RawMessage `json:"annotations,omitempty"`
}

// Server is how a server names itself at initialize (serverInfo).
type Server struct {
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
}

// Label is the server's name for people: its title, else its name.
func (s Server) Label() string {
	if t := strings.TrimSpace(s.Title); t != "" {
		return t
	}
	return strings.TrimSpace(s.Name)
}

// Label is the tool's name for people: its title, else its annotations'
// title; "" when the server gives neither.
func (t Tool) Label() string {
	if s := strings.TrimSpace(t.Title); s != "" {
		return s
	}
	var a struct {
		Title string `json:"title"`
	}
	_ = json.Unmarshal(t.Annotations, &a)
	return strings.TrimSpace(a.Title)
}

// RPCError is a JSON-RPC error from the server.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("mcp error %d: %s", e.Code, e.Message) }

// ErrUnauthorized: the server refused the token (HTTP 401).
var ErrUnauthorized = errors.New("mcp: the server refused the token")

// CallError is a tools/call that failed in transport. Sent says whether
// the request may have reached the server: if it did, the call's effect
// is unknown and it is never sent again.
type CallError struct {
	Sent bool
	Err  error
}

func (e *CallError) Error() string { return e.Err.Error() }
func (e *CallError) Unwrap() error { return e.Err }

// Client calls a vendor's MCP server.
type Client struct {
	HTTP *http.Client
	Name string // clientInfo.name
}

// ListTools lists the server's tools, following pagination, with how the
// server names itself.
func (c *Client) ListTools(ctx context.Context, url, token string) ([]Tool, Server, error) {
	s, err := c.open(ctx, url, token)
	if err != nil {
		return nil, Server{}, err
	}
	var out []Tool
	cursor := ""
	for page := 0; page < 20; page++ {
		params := map[string]any{}
		if cursor != "" {
			params["cursor"] = cursor
		}
		raw, err := s.rpc(ctx, 2+page, "tools/list", params)
		if err != nil {
			return nil, Server{}, err
		}
		var r struct {
			Tools      []Tool `json:"tools"`
			NextCursor string `json:"nextCursor"`
		}
		if err := json.Unmarshal(raw, &r); err != nil {
			return nil, Server{}, err
		}
		out = append(out, r.Tools...)
		if r.NextCursor == "" {
			break
		}
		cursor = r.NextCursor
	}
	return out, s.server, nil
}

// CallTool calls one tool and returns its result. An error from the
// handshake is never Sent; one after tools/call left is.
func (c *Client) CallTool(ctx context.Context, url, token, name string, args json.RawMessage) (json.RawMessage, error) {
	s, err := c.open(ctx, url, token)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return nil, err
		}
		return nil, &CallError{Sent: false, Err: err}
	}
	raw, err := s.rpc(ctx, 2, "tools/call", map[string]any{"name": name, "arguments": args})
	var rpc *RPCError
	switch {
	case err == nil:
		return raw, nil
	case errors.Is(err, ErrUnauthorized), errors.As(err, &rpc):
		// Answered: refused before it ran, or the tool's own error.
		return nil, err
	default:
		return nil, &CallError{Sent: true, Err: err}
	}
}

type session struct {
	c       *Client
	url     string
	token   string
	id      string
	version string // negotiated at initialize; sent on every later request
	server  Server // how the server named itself at initialize
}

func (c *Client) open(ctx context.Context, url, token string) (*session, error) {
	s := &session{c: c, url: url, token: token}
	name := c.Name
	if name == "" {
		name = "planque"
	}
	raw, err := s.rpc(ctx, 1, "initialize", map[string]any{
		"protocolVersion": Version,
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]string{"name": name, "version": "1"},
	})
	if err != nil {
		return nil, fmt.Errorf("initialize: %w", err)
	}
	var init struct {
		ProtocolVersion string `json:"protocolVersion"`
		ServerInfo      Server `json:"serverInfo"`
	}
	_ = json.Unmarshal(raw, &init)
	s.version, s.server = init.ProtocolVersion, init.ServerInfo
	if err := s.post(ctx, map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"}, nil); err != nil {
		return nil, fmt.Errorf("initialized: %w", err)
	}
	return s, nil
}

func (s *session) rpc(ctx context.Context, id int, method string, params any) (json.RawMessage, error) {
	var msg message
	err := s.post(ctx, map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params}, &msg)
	if err != nil {
		return nil, err
	}
	if msg.Error != nil {
		return nil, msg.Error
	}
	return msg.Result, nil
}

type message struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error"`
}

func (s *session) post(ctx context.Context, body map[string]any, into *message) error {
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if s.version != "" {
		req.Header.Set("MCP-Protocol-Version", s.version)
	}
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}
	if s.id != "" {
		req.Header.Set("Mcp-Session-Id", s.id)
	}
	hc := s.c.HTTP
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if sid := resp.Header.Get("Mcp-Session-Id"); sid != "" {
		s.id = sid
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if into == nil {
		if resp.StatusCode >= 300 {
			return fmt.Errorf("status %d", resp.StatusCode)
		}
		return nil
	}
	want, _ := json.Marshal(body["id"])
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Refused at the transport, before any tool ran: an answer, and
		// what it says is kept.
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		var m message
		if json.Unmarshal(b, &m) == nil && m.Error != nil {
			*into = m
			return nil
		}
		msg := strings.TrimSpace(string(b))
		if len(msg) > 300 {
			msg = msg[:300] + "…"
		}
		return &RPCError{Code: -resp.StatusCode, Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, msg)}
	}
	m, err := readMessage(resp, string(want))
	if err != nil {
		return fmt.Errorf("status %d: %w", resp.StatusCode, err)
	}
	*into = m
	return nil
}

// readMessage reads the response with id want from a JSON or SSE body.
func readMessage(resp *http.Response, want string) (message, error) {
	var msg message
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/event-stream") {
		err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&msg)
		return msg, err
	}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64<<10), 8<<20)
	for sc.Scan() {
		line, ok := strings.CutPrefix(sc.Text(), "data:")
		if !ok {
			continue
		}
		msg = message{}
		if json.Unmarshal([]byte(strings.TrimSpace(line)), &msg) == nil && string(msg.ID) == want {
			return msg, nil
		}
	}
	if err := sc.Err(); err != nil {
		return message{}, err
	}
	return message{}, errors.New("stream ended without a response")
}
