package mcp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
