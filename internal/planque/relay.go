package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/store"
)

// opLease is how long a claimed operation may stay started before it is
// treated as interrupted (RFC 046 D5).
const opLease = 30 * time.Second

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func rpcResult(w http.ResponseWriter, id json.RawMessage, result any) {
	writeJSON(w, http.StatusOK, map[string]any{"jsonrpc": "2.0", "id": id, "result": result})
}

func rpcError(w http.ResponseWriter, status int, id json.RawMessage, code int, msg string) {
	writeJSON(w, status, map[string]any{"jsonrpc": "2.0", "id": id, "error": map[string]any{"code": code, "message": msg}})
}

// toolResult is an MCP tool result with one text block.
func toolResult(text string, isError bool) map[string]any {
	return map[string]any{"content": []map[string]string{{"type": "text", "text": text}}, "isError": isError}
}

// relay is planque's MCP endpoint for one grant, which KIFF's gateway
// calls in relay mode. It is stateless and speaks the 2025 revisions;
// the gateway's server/discover gets "method not found", which tells it
// so.
func (s *Server) relay(w http.ResponseWriter, r *http.Request) {
	tenant, id := r.PathValue("tenant"), r.PathValue("id")
	g, err := s.loadGrant(r.Context(), tenant, id)
	if errors.Is(err, errNotFound) || (err == nil && g.Status != StatusConnected) {
		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
		writeErr(w, http.StatusUnauthorized, "no connected grant here")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	// The transport credential says the caller is KIFF's gateway for this
	// grant. Alone it executes nothing: each call also needs a permit.
	cred, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	sum := sha256.Sum256([]byte(strings.TrimSpace(cred)))
	want, _ := hex.DecodeString(g.Transport)
	if cred == "" || len(want) != len(sum) || subtle.ConstantTimeCompare(sum[:], want) != 1 {
		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
		writeErr(w, http.StatusUnauthorized, "the relay's transport credential is required")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "unreadable body")
		return
	}
	var req rpcRequest
	if json.Unmarshal(body, &req) != nil || req.JSONRPC != "2.0" {
		rpcError(w, http.StatusBadRequest, nil, -32700, "parse error")
		return
	}
	if len(req.ID) == 0 { // a notification
		w.WriteHeader(http.StatusAccepted)
		return
	}
	switch req.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(req.Params, &p)
		v := mcp.Version
		if slices.Contains([]string{"2025-03-26", "2025-06-18", "2025-11-25"}, p.ProtocolVersion) {
			v = p.ProtocolVersion
		}
		rpcResult(w, req.ID, map[string]any{"protocolVersion": v, "capabilities": map[string]any{"tools": map[string]any{}},
			"serverInfo": map[string]string{"name": "planque", "version": "1"}})
	case "ping":
		rpcResult(w, req.ID, map[string]any{})
	case "tools/list":
		rpcResult(w, req.ID, map[string]any{"tools": relayTools(g)})
	case "tools/call":
		var p struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if json.Unmarshal(req.Params, &p) != nil {
			rpcError(w, http.StatusOK, req.ID, -32602, "invalid params")
			return
		}
		res := s.call(r.Context(), g, p.Name, p.Arguments, r.Header.Get(permit.Header))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": res})
	default:
		rpcError(w, http.StatusOK, req.ID, -32601, "method not found")
	}
}

// relayTools lists the granted tools with planque's operation argument,
// which KIFF fills and binds into the permit.
func relayTools(g *Grant) []mcp.Tool {
	out := make([]mcp.Tool, 0, len(g.Tools))
	for _, t := range g.Tools {
		var schema map[string]any
		if json.Unmarshal(t.InputSchema, &schema) != nil || schema == nil {
			schema = map[string]any{"type": "object"}
		}
		props, _ := schema["properties"].(map[string]any)
		if props == nil {
			props = map[string]any{}
		}
		props[IdempotencyArg] = map[string]any{"type": "string", "description": "Set by KIFF: the call's operation id."}
		schema["properties"] = props
		req, _ := schema["required"].([]any)
		schema["required"] = append(req, IdempotencyArg)
		b, _ := json.Marshal(schema)
		out = append(out, mcp.Tool{Name: t.Name, Title: t.Title, Description: t.Description, InputSchema: b})
	}
	return out
}

// call runs one tools/call: the permit first, then the operation's
// record, then the vendor, at most once per op (RFC 046 D5).
func (s *Server) call(ctx context.Context, g *Grant, name string, raw json.RawMessage, token string) any {
	refused := func(why string) any {
		s.log().InfoContext(ctx, "call refused", slog.String("grant", g.ID), slog.String("tool", name), slog.String("why", why))
		return toolResult("Refused by planque: "+why+". Nothing was sent to "+hostOf(g.Request.Server)+".", true)
	}
	i := slices.IndexFunc(g.Tools, func(t GrantedTool) bool { return t.Name == name })
	if i < 0 {
		return refused("this grant does not include " + name)
	}
	tool := g.Tools[i]
	if len(raw) == 0 {
		raw = json.RawMessage("{}")
	}
	p, err := s.Permits.Verify(ctx, token, permit.Want{Audience: g.Audience, Tenant: g.Tenant, Tool: name, Args: raw})
	if err != nil {
		return refused(err.Error())
	}
	var args map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&args) != nil {
		return refused("the arguments are not a JSON object")
	}
	if op, _ := args[IdempotencyArg].(string); op != p.Op {
		return refused("the permit's operation is not this call's " + IdempotencyArg)
	}
	key := store.Item{PK: tenantPK(g.Tenant), SK: skOp(g.ID, p.Op)}
	var prev operation
	if it, err := getDoc(ctx, s.KV, key.PK, key.SK, &prev); err == nil {
		return s.recorded(ctx, g, it, prev, p)
	} else if !errors.Is(err, errNotFound) {
		return refused("planque could not read the operation's record")
	}
	now := s.now()
	if !p.Live(now) {
		return refused("the permit has expired")
	}
	// Claim the op: only the claimer may send.
	op := operation{Tool: name, ArgsHash: p.ArgsHash, State: OpStarted, StartedAt: now.UTC(), Decision: p.Decision}
	claim := key
	claim.Ver, claim.Lease = 1, now.Add(opLease).UnixMilli()
	if err := putDoc(ctx, s.KV, claim, op, store.Cond{NotExists: true}); errors.Is(err, store.ErrCondition) {
		it, err := getDoc(ctx, s.KV, key.PK, key.SK, &prev)
		if err != nil {
			return refused("planque could not read the operation's record")
		}
		return s.recorded(ctx, g, it, prev, p)
	} else if err != nil {
		return refused("planque could not record the operation")
	}
	delete(args, IdempotencyArg)
	vendorArgs, err := tool.vendorArgs(args)
	if err != nil {
		return s.finish(ctx, claim, op, OpNotSent, toolResult("Refused by planque: "+err.Error()+". Nothing was sent.", true))
	}

	access, ver, err := s.accessToken(ctx, g, 0)
	if err != nil {
		return s.finish(ctx, claim, op, OpNotSent, toolResult("Not sent: "+err.Error(), true))
	}
	var res json.RawMessage
	for attempt := 0; ; attempt++ {
		// The permit must still be live when the call leaves (D5).
		if !p.Live(s.now()) {
			return s.finish(ctx, claim, op, OpNotSent, toolResult("Not sent: the permit expired before the call could be sent.", true))
		}
		res, err = s.MCP.CallTool(ctx, g.Request.Server, access, tool.Tool, vendorArgs)
		if !errors.Is(err, mcp.ErrUnauthorized) || attempt > 0 {
			break
		}
		// The vendor refused the token before running anything: refresh
		// it once and try again.
		if access, ver, err = s.accessToken(ctx, g, ver); err != nil {
			return s.finish(ctx, claim, op, OpNotSent, toolResult("Not sent: "+err.Error(), true))
		}
	}
	var ce *mcp.CallError
	var rpc *mcp.RPCError
	switch {
	case err == nil:
		return s.finish(ctx, claim, op, OpDone, res)
	case errors.As(err, &rpc):
		return s.finish(ctx, claim, op, OpDone, toolResult(hostOf(g.Request.Server)+" returned an error: "+rpc.Message, true))
	case errors.Is(err, mcp.ErrUnauthorized), errors.As(err, &ce) && !ce.Sent:
		return s.finish(ctx, claim, op, OpNotSent, toolResult("Not sent: "+hostOf(g.Request.Server)+" could not be reached or refused planque's token.", true))
	default:
		s.log().WarnContext(ctx, "vendor answer lost", slog.String("grant", g.ID), slog.String("op", p.Op), slog.String("err", err.Error()))
		return s.finish(ctx, claim, op, OpUnknown, unknownResult(g))
	}
}

func unknownResult(g *Grant) map[string]any {
	return toolResult("Sent to "+hostOf(g.Request.Server)+", but its answer was lost. planque will not send it again; check the account before making a new call.", true)
}

// finish records how an operation ended and returns the result.
func (s *Server) finish(ctx context.Context, claim store.Item, op operation, state string, res any) any {
	op.State = state
	op.Result, _ = json.Marshal(res)
	it := claim
	it.Ver, it.Lease = claim.Ver+1, 0
	if err := putDoc(ctx, s.KV, it, op, store.Cond{Ver: &claim.Ver}); err != nil {
		s.log().ErrorContext(ctx, "operation not recorded", slog.String("op", claim.SK), slog.String("err", err.Error()))
	}
	return res
}

// recorded answers for an op already claimed: never by sending again.
func (s *Server) recorded(ctx context.Context, g *Grant, it store.Item, op operation, p permit.Permit) any {
	if op.ArgsHash != p.ArgsHash {
		return toolResult("Refused by planque: this operation was used for other arguments. Nothing was sent.", true)
	}
	switch op.State {
	case OpDone, OpNotSent, OpUnknown:
		return op.Result
	}
	if it.Lease > s.now().UnixMilli() {
		return toolResult("This call is in progress. Nothing new was sent; ask again shortly.", true)
	}
	// Claimed, and its holder stopped before recording an answer: it may
	// or may not have sent. A generic MCP tool has no lookup, so this goes
	// to a person.
	res := unknownResult(g)
	ver := it.Ver
	op.State = OpUnknown
	op.Result, _ = json.Marshal(res)
	next := it
	next.Ver, next.Lease = ver+1, 0
	_ = putDoc(ctx, s.KV, next, op, store.Cond{Ver: &ver})
	return res
}
