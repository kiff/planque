package fakes

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// The desk: the fake vendor's second API, a support desk reached by an API
// key, not MCP (kiff/agents#293). It isn't payments: every call is scoped by
// its workspace, it has a count-limited write (a reply) and an amount-free
// one (a label), and its reads are named apart from its writes. planque
// reaches it as an API-key grant, through its OpenAPI description.

// DeskKeyHeader is where the desk takes its API key.
const DeskKeyHeader = "X-Desk-Key"

// DeskWorkspace is the one workspace the desk has.
const DeskWorkspace = "ws_demo"

// kvDesk keeps the desk's writes: SK when (ns) and an id, Doc a VendorCall.
const kvDesk = "fakevendor#desk"

// deskTickets are the workspace's tickets.
var deskTickets = []map[string]any{
	{"id": "tk_1001", "subject": "Order arrived damaged", "status": "open", "requester": "ana@example.com"},
	{"id": "tk_1002", "subject": "How do I change my plan?", "status": "open", "requester": "ben@example.com"},
	{"id": "tk_1003", "subject": "Invoice address is wrong", "status": "pending", "requester": "chloe@example.com"},
}

// DeskURL is the desk's API, and DeskDescription its OpenAPI description.
func (v *Vendor) DeskURL() string         { return v.root() + "/desk/v1" }
func (v *Vendor) DeskDescription() string { return v.root() + "/desk/openapi.json" }

// deskRoutes adds the desk to mux.
func (v *Vendor) deskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /desk/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(strings.ReplaceAll(deskOpenAPI, "{{base}}", v.DeskURL())))
	})
	mux.HandleFunc("GET /desk/v1/workspaces/{workspace}/tickets", v.desk(func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		out := []map[string]any{}
		for _, t := range deskTickets {
			if status == "" || t["status"] == status {
				out = append(out, t)
			}
		}
		writeJSON(w, 200, map[string]any{"data": out})
	}))
	mux.HandleFunc("GET /desk/v1/workspaces/{workspace}/tickets/{ticket}", v.desk(func(w http.ResponseWriter, r *http.Request) {
		t, ok := deskTicket(r.PathValue("ticket"))
		if !ok {
			writeJSON(w, 404, map[string]string{"error": "no such ticket"})
			return
		}
		writeJSON(w, 200, t)
	}))
	mux.HandleFunc("POST /desk/v1/workspaces/{workspace}/tickets/{ticket}/replies", v.desk(func(w http.ResponseWriter, r *http.Request) {
		v.deskWrite(w, r, "postReply", "body", "rp")
	}))
	mux.HandleFunc("POST /desk/v1/workspaces/{workspace}/tickets/{ticket}/labels", v.desk(func(w http.ResponseWriter, r *http.Request) {
		v.deskWrite(w, r, "applyLabel", "label", "lb")
	}))
}

// desk checks the key and the workspace before h.
func (v *Vendor) desk(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v.DeskKey == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get(DeskKeyHeader)), []byte(v.DeskKey)) != 1 {
			writeJSON(w, 401, map[string]string{"error": "a valid " + DeskKeyHeader + " is required"})
			return
		}
		if r.PathValue("workspace") != DeskWorkspace {
			writeJSON(w, 404, map[string]string{"error": "no such workspace"})
			return
		}
		h(w, r)
	}
}

// deskWrite records one write on a ticket: field is its one required text.
func (v *Vendor) deskWrite(w http.ResponseWriter, r *http.Request, op, field, prefix string) {
	ticket := r.PathValue("ticket")
	if _, ok := deskTicket(ticket); !ok {
		writeJSON(w, 404, map[string]string{"error": "no such ticket"})
		return
	}
	var body map[string]any
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		writeJSON(w, 400, map[string]string{"error": "the body must be a JSON object"})
		return
	}
	text, _ := body[field].(string)
	if strings.TrimSpace(text) == "" {
		writeJSON(w, 422, map[string]string{"error": field + " is required"})
		return
	}
	call := VendorCall{Tool: op, Args: map[string]any{"workspace": r.PathValue("workspace"), "ticket": ticket, field: text}}
	id := next(prefix)
	doc, _ := json.Marshal(call)
	if err := v.keep(r.Context(), kvDesk, fmt.Sprintf("%020d#%s", time.Now().UnixNano(), id), doc, 7*24*time.Hour); err != nil {
		writeJSON(w, 500, map[string]string{"error": "not recorded"})
		return
	}
	v.mu.Lock()
	v.DeskCalls = append(v.DeskCalls, call)
	v.mu.Unlock()
	writeJSON(w, 201, map[string]any{"id": id, "ticket": ticket})
}

func deskTicket(id string) (map[string]any, bool) {
	for _, t := range deskTickets {
		if t["id"] == id {
			return t, true
		}
	}
	return nil, false
}

// DeskCallsSoFar is every write the desk recorded, across its instances.
func (v *Vendor) DeskCallsSoFar(r *http.Request) ([]VendorCall, error) {
	its, err := v.KV.Query(r.Context(), kvDesk, "")
	if err != nil {
		return nil, err
	}
	out := make([]VendorCall, 0, len(its))
	for _, it := range its {
		var c VendorCall
		if json.Unmarshal(it.Doc, &c) == nil {
			out = append(out, c)
		}
	}
	return out, nil
}

// deskOpenAPI is the desk's description, as planque reads one: JSON,
// OpenAPI 3.0, one server, the key in a header.
const deskOpenAPI = `{
 "openapi": "3.0.3",
 "info": {"title": "Fake Desk API", "version": "1"},
 "servers": [{"url": "{{base}}"}],
 "security": [{"deskKey": []}],
 "components": {"securitySchemes": {"deskKey": {"type": "apiKey", "in": "header", "name": "X-Desk-Key"}}},
 "paths": {
  "/workspaces/{workspace}/tickets": {"get": {"operationId": "searchTickets", "summary": "Search a workspace's tickets",
   "parameters": [{"name": "workspace", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "status", "in": "query", "schema": {"type": "string", "enum": ["open", "pending", "solved"]}}]}},
  "/workspaces/{workspace}/tickets/{ticket}": {"get": {"operationId": "fetchTicket", "summary": "Fetch one ticket",
   "parameters": [{"name": "workspace", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "ticket", "in": "path", "required": true, "schema": {"type": "string"}}]}},
  "/workspaces/{workspace}/tickets/{ticket}/replies": {"post": {"operationId": "postReply", "summary": "Reply to a ticket's requester",
   "parameters": [{"name": "workspace", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "ticket", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["body"],
    "properties": {"body": {"type": "string", "description": "The reply's text, sent to the requester"}}}}}}}},
  "/workspaces/{workspace}/tickets/{ticket}/labels": {"post": {"operationId": "applyLabel", "summary": "Label a ticket",
   "parameters": [{"name": "workspace", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "ticket", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["label"],
    "properties": {"label": {"type": "string"}}}}}}}}
 }
}`
