package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// A server like Stripe's: it negotiates an older revision, wants it on
// every later request, and refuses a bad call with a plain HTTP 422.
func TestStripeLikeServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Method != "initialize" && r.Header.Get("MCP-Protocol-Version") != "2025-03-26" {
			http.Error(w, `{"error":"unsupported protocol version"}`, http.StatusBadRequest)
			return
		}
		switch req.Method {
		case "initialize":
			w.Header().Set("Mcp-Session-Id", "s")
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID,
				"result": map[string]any{"protocolVersion": "2025-03-26"}})
		case "notifications/initialized":
			w.WriteHeader(http.StatusAccepted)
		case "tools/call":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error":"Missing required parameter stripe_context"}`))
		}
	}))
	defer srv.Close()
	_, err := (&Client{}).CallTool(t.Context(), srv.URL, "tok", "stripe_api_read", json.RawMessage(`{}`))
	var rpc *RPCError
	if !errors.As(err, &rpc) || !strings.Contains(rpc.Message, "stripe_context") {
		t.Fatalf("want the server's refusal as an answer, got %T %v", err, err)
	}
}

// sessionServer counts the sessions opened, knows only those it opened and
// hasn't forgotten, and answers each call.
type sessionServer struct {
	mu     sync.Mutex
	opened int
	live   map[string]bool
	calls  []string // session of each call
	// rpc404 answers an unknown session with HTTP 404 and a JSON-RPC error body.
	rpc404 bool
}

func (v *sessionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	v.mu.Lock()
	defer v.mu.Unlock()
	switch req.Method {
	case "initialize":
		v.opened++
		sid := fmt.Sprintf("s%d", v.opened)
		if v.live == nil {
			v.live = map[string]bool{}
		}
		v.live[sid] = true
		w.Header().Set("Mcp-Session-Id", sid)
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": map[string]any{"protocolVersion": Version}})
	case "notifications/initialized":
		w.WriteHeader(http.StatusAccepted)
	case "tools/call":
		sid := r.Header.Get("Mcp-Session-Id")
		if !v.live[sid] {
			if v.rpc404 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "error": map[string]any{"code": -32001, "message": "Session not found"}})
				return
			}
			http.Error(w, "unknown session", http.StatusNotFound)
			return
		}
		v.calls = append(v.calls, sid)
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": map[string]any{"content": []any{}}})
	}
}

// Calls to one server with one sign-in reuse its session (kiff/planque#47);
// another sign-in gets its own.
func TestCallsReuseTheSession(t *testing.T) {
	v := &sessionServer{}
	srv := httptest.NewServer(v)
	defer srv.Close()
	c := &Client{}
	for i := 0; i < 3; i++ {
		if _, err := c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`)); err != nil {
			t.Fatal(err)
		}
	}
	if v.opened != 1 || strings.Join(v.calls, ",") != "s1,s1,s1" {
		t.Fatalf("want one session for three calls, opened %d, calls %v", v.opened, v.calls)
	}
	if _, err := c.CallTool(t.Context(), srv.URL, "other", "get_charge", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	if v.opened != 2 || v.calls[3] != "s2" {
		t.Fatalf("another sign-in shares no session: opened %d, calls %v", v.opened, v.calls)
	}
}

// A server that forgot the session refused the call before it ran: the call
// goes once more on a new session, and only once.
func TestAForgottenSessionIsOpenedAgain(t *testing.T) {
	v := &sessionServer{}
	srv := httptest.NewServer(v)
	defer srv.Close()
	c := &Client{}
	if _, err := c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.live = map[string]bool{}
	v.mu.Unlock()
	if _, err := c.CallTool(t.Context(), srv.URL, "tok", "create_refund", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	if v.opened != 2 || strings.Join(v.calls, ",") != "s1,s2" {
		t.Fatalf("want the call once, on a new session: opened %d, calls %v", v.opened, v.calls)
	}
}

// A server that answers an unknown session with HTTP 404 and a JSON-RPC
// error body forgot it all the same: the call goes once more on a new
// session, and the dead one is never kept.
func TestAForgottenSessionWithAnRPCBodyIsOpenedAgain(t *testing.T) {
	v := &sessionServer{rpc404: true}
	srv := httptest.NewServer(v)
	defer srv.Close()
	c := &Client{}
	if _, err := c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.live = map[string]bool{}
	v.mu.Unlock()
	if _, err := c.CallTool(t.Context(), srv.URL, "tok", "create_refund", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`)); err != nil {
		t.Fatal(err)
	}
	if v.opened != 2 || strings.Join(v.calls, ",") != "s1,s2,s2" {
		t.Fatalf("want the call once on a new session, then that one kept: opened %d, calls %v", v.opened, v.calls)
	}
}

// A session unused for longer than Idle isn't reused.
func TestAnIdleSessionIsNotReused(t *testing.T) {
	v := &sessionServer{}
	srv := httptest.NewServer(v)
	defer srv.Close()
	now := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)
	c := &Client{Idle: time.Minute, now: func() time.Time { return now }}
	_, _ = c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`))
	now = now.Add(2 * time.Minute)
	_, _ = c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`))
	if v.opened != 2 {
		t.Fatalf("want a new session after the idle time, opened %d", v.opened)
	}
}

// Calls may go at once on one session, each with its own request id.
func TestCallsAtOnceShareTheSession(t *testing.T) {
	v := &sessionServer{}
	srv := httptest.NewServer(v)
	defer srv.Close()
	c := &Client{}
	_, _ = c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := c.CallTool(t.Context(), srv.URL, "tok", "get_charge", json.RawMessage(`{}`)); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if v.opened != 1 || len(v.calls) != 9 {
		t.Fatalf("want 9 calls on one session, opened %d, calls %d", v.opened, len(v.calls))
	}
}
