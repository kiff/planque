package fakes

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
)

// Mail: the fake vendor's mailbox, reached by an API key, as the desk is. It
// is what a job that writes email is proved on: a message is written as a
// draft, kept for review, or sent, and the two are different writes, so a
// grant can allow drafts and hold every send. Sending a draft is a write of
// its own. Nothing here is written for one job: a mailbox has folders, a
// search, drafts and sends, as any mail API does.
//
// Like the desk, it keeps no state between calls: reads answer the seeded
// messages, and each write is recorded. A run is the same every time.

// MailKeyHeader is where the mailbox takes its API key: the vendor's one key
// (DeskKey), in its own header.
const MailKeyHeader = "X-Mail-Key"

// MailBox is the mailbox Setup connects; MailEvalBox is the one an Eval case
// works, apart from it so a setup test reads the same messages as before.
const (
	MailBox     = "mb_demo"
	MailEvalBox = "mb_eval"
)

// kvMail keeps the mailbox's writes: SK when (ns) and an id, Doc a VendorCall.
const kvMail = "fakevendor#mail"

// mailFolders are a mailbox's folders, as its list takes them.
var mailFolders = []string{"inbox", "sent", "drafts"}

// mailMessages are MailBox's messages.
var mailMessages = []map[string]any{
	{"id": "ms_1001", "folder": "inbox", "from": "ana@example.com", "to": "owner@shop.example", "subject": "Table for six on Friday?",
		"date": "2026-01-12T09:14:00Z", "body": "Hello, do you have a table for six this Friday at 20:00?"},
	{"id": "ms_1002", "folder": "inbox", "from": "billing@supplier.example", "to": "owner@shop.example", "subject": "Your receipt for January",
		"date": "2026-01-13T07:02:00Z", "body": "Thank you for your payment. Your receipt is attached."},
	{"id": "ms_1003", "folder": "sent", "from": "owner@shop.example", "to": "ben@example.com", "subject": "Your order has shipped",
		"date": "2026-01-10T16:40:00Z", "body": "Hi Ben, your order left our warehouse today."},
}

// mailEvalMessages are MailEvalBox's messages, by what each one tests in a
// job that writes first emails from a list (SheetEval, sheet.go):
//
//	M-4001  sent to fay@fernco.example: she was written to already, so no second message
//	M-4002  inbox, her answer: a conversation that exists, not a new company to write to
//	M-4003  inbox, a newsletter: nothing to do with the list
//	M-4004  sent to a customer: an address that isn't on the list
var mailEvalMessages = []map[string]any{
	{"id": "M-4001", "folder": "sent", "from": "owner@shop.example", "to": "fay@fernco.example", "subject": "An idea for Fern & Co's returns",
		"date": "2026-01-08T10:05:00Z", "body": "Hello Fay, I noticed returns at Fern & Co go through a contact form. Would you be open to a short call?"},
	{"id": "M-4002", "folder": "inbox", "from": "fay@fernco.example", "to": "owner@shop.example", "subject": "Re: An idea for Fern & Co's returns",
		"date": "2026-01-09T08:31:00Z", "body": "Thanks for writing. Could you tell me more next week?"},
	{"id": "M-4003", "folder": "inbox", "from": "news@retailweekly.example", "to": "owner@shop.example", "subject": "This week in retail",
		"date": "2026-01-12T06:00:00Z", "body": "Five stories from the week."},
	{"id": "M-4004", "folder": "sent", "from": "owner@shop.example", "to": "lea@customer.example", "subject": "Your invoice",
		"date": "2026-01-11T14:20:00Z", "body": "Hello Lea, your invoice is attached."},
}

// mailBoxes are the mailboxes and their messages.
var mailBoxes = map[string][]map[string]any{MailBox: mailMessages, MailEvalBox: mailEvalMessages}

// MailURL is the mailbox's API, and MailDescription its OpenAPI description.
func (v *Vendor) MailURL() string         { return v.root() + "/mail/v1" }
func (v *Vendor) MailDescription() string { return v.root() + "/mail/openapi.json" }

// mailRoutes adds the mailbox to mux.
func (v *Vendor) mailRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /mail/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(strings.ReplaceAll(mailOpenAPI, "{{base}}", v.MailURL())))
	})
	mux.HandleFunc("GET /mail/v1/mailboxes/{mailbox}/messages", v.mail(func(w http.ResponseWriter, r *http.Request) {
		folder := r.URL.Query().Get("folder")
		// A folder the mailbox doesn't have is the caller's mistake, never an
		// empty list that reads as "no mail".
		if folder != "" && !slices.Contains(mailFolders, folder) {
			writeJSON(w, 400, map[string]string{"error": "folder is inbox, sent or drafts"})
			return
		}
		// with: messages from or to this address, whatever its case.
		with := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("with")))
		out := []map[string]any{}
		for _, m := range mailBoxes[r.PathValue("mailbox")] {
			if folder != "" && m["folder"] != folder {
				continue
			}
			if with != "" && strings.ToLower(m["from"].(string)) != with && strings.ToLower(m["to"].(string)) != with {
				continue
			}
			out = append(out, m)
		}
		writeJSON(w, 200, map[string]any{"data": out})
	}))
	mux.HandleFunc("GET /mail/v1/mailboxes/{mailbox}/messages/{message}", v.mail(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range mailBoxes[r.PathValue("mailbox")] {
			if m["id"] == r.PathValue("message") {
				writeJSON(w, 200, m)
				return
			}
		}
		writeJSON(w, 404, map[string]string{"error": "no such message"})
	}))
	// A draft: written, kept in the mailbox, sent to no one.
	mux.HandleFunc("POST /mail/v1/mailboxes/{mailbox}/drafts", v.mail(func(w http.ResponseWriter, r *http.Request) {
		v.mailWrite(w, r, "createDraft", "dr")
	}))
	// A message sent at once.
	mux.HandleFunc("POST /mail/v1/mailboxes/{mailbox}/messages", v.mail(func(w http.ResponseWriter, r *http.Request) {
		v.mailWrite(w, r, "sendMessage", "ms")
	}))
	// A draft sent: a write of its own, apart from writing it.
	mux.HandleFunc("POST /mail/v1/mailboxes/{mailbox}/drafts/{draft}/send", v.mail(func(w http.ResponseWriter, r *http.Request) {
		draft := r.PathValue("draft")
		if !strings.HasPrefix(draft, "dr_") {
			writeJSON(w, 404, map[string]string{"error": "no such draft"})
			return
		}
		call := VendorCall{Tool: "sendDraft", Args: map[string]any{"mailbox": r.PathValue("mailbox"), "draft": draft}}
		if !v.record(w, r, kvMail, &v.MailCalls, call, next("sd")) {
			return
		}
		writeJSON(w, 200, map[string]any{"id": draft, "sent": true})
	}))
}

// mail checks the key and the mailbox before h.
func (v *Vendor) mail(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v.DeskKey == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get(MailKeyHeader)), []byte(v.DeskKey)) != 1 {
			writeJSON(w, 401, map[string]string{"error": "a valid " + MailKeyHeader + " is required"})
			return
		}
		if _, ok := mailBoxes[r.PathValue("mailbox")]; !ok {
			writeJSON(w, 404, map[string]string{"error": "no such mailbox"})
			return
		}
		h(w, r)
	}
}

// mailWrite records a draft or a sent message: to, subject and body, each
// required, and to one address.
func (v *Vendor) mailWrite(w http.ResponseWriter, r *http.Request, op, prefix string) {
	var body map[string]any
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		writeJSON(w, 400, map[string]string{"error": "the body must be a JSON object"})
		return
	}
	args := map[string]any{"mailbox": r.PathValue("mailbox")}
	for _, field := range []string{"to", "subject", "body"} {
		text, _ := body[field].(string)
		if strings.TrimSpace(text) == "" {
			writeJSON(w, 422, map[string]string{"error": field + " is required"})
			return
		}
		args[field] = text
	}
	// One recipient a message: a list of addresses isn't one, so a write to
	// many is many writes, each counted.
	if to := args["to"].(string); strings.Count(to, "@") != 1 || strings.ContainsAny(to, ",; ") {
		writeJSON(w, 422, map[string]string{"error": "to is one email address"})
		return
	}
	id := next(prefix)
	if !v.record(w, r, kvMail, &v.MailCalls, VendorCall{Tool: op, Args: args}, id) {
		return
	}
	writeJSON(w, 201, map[string]any{"id": id, "to": args["to"]})
}

// record keeps one write under pk and in calls; false when it couldn't,
// answered already.
func (v *Vendor) record(w http.ResponseWriter, r *http.Request, pk string, calls *[]VendorCall, call VendorCall, id string) bool {
	doc, _ := json.Marshal(call)
	if err := v.keep(r.Context(), pk, fmt.Sprintf("%020d#%s", time.Now().UnixNano(), id), doc, 7*24*time.Hour); err != nil {
		writeJSON(w, 500, map[string]string{"error": "not recorded"})
		return false
	}
	v.mu.Lock()
	*calls = append(*calls, call)
	v.mu.Unlock()
	return true
}

// mailOpenAPI is the mailbox's description, as planque reads one: JSON,
// OpenAPI 3.0, one server, the key in a header. Each write says whether it
// reaches anyone, as a vendor's own reference would.
const mailOpenAPI = `{
 "openapi": "3.0.3",
 "info": {"title": "Fake Mail API", "version": "1"},
 "servers": [{"url": "{{base}}"}],
 "security": [{"mailKey": []}],
 "components": {"securitySchemes": {"mailKey": {"type": "apiKey", "in": "header", "name": "X-Mail-Key"}},
  "schemas": {"Message": {"type": "object", "required": ["to", "subject", "body"],
   "properties": {
    "to": {"type": "string", "description": "One email address"},
    "subject": {"type": "string"},
    "body": {"type": "string", "description": "The message's text"}}}}},
 "paths": {
  "/mailboxes/{mailbox}/messages": {
   "get": {"operationId": "listMessages", "summary": "List a mailbox's messages",
    "parameters": [{"name": "mailbox", "in": "path", "required": true, "schema": {"type": "string"}},
                   {"name": "folder", "in": "query", "schema": {"type": "string", "enum": ["inbox", "sent", "drafts"]}},
                   {"name": "with", "in": "query", "description": "Only messages from or to this address", "schema": {"type": "string"}}]},
   "post": {"operationId": "sendMessage", "summary": "Send a message now, to its recipient",
    "parameters": [{"name": "mailbox", "in": "path", "required": true, "schema": {"type": "string"}}],
    "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Message"}}}}}},
  "/mailboxes/{mailbox}/messages/{message}": {"get": {"operationId": "getMessage", "summary": "Fetch one message",
   "parameters": [{"name": "mailbox", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "message", "in": "path", "required": true, "schema": {"type": "string"}}]}},
  "/mailboxes/{mailbox}/drafts": {"post": {"operationId": "createDraft", "summary": "Write a draft: kept in the mailbox, sent to no one",
   "parameters": [{"name": "mailbox", "in": "path", "required": true, "schema": {"type": "string"}}],
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Message"}}}}}},
  "/mailboxes/{mailbox}/drafts/{draft}/send": {"post": {"operationId": "sendDraft", "summary": "Send a draft to its recipient",
   "parameters": [{"name": "mailbox", "in": "path", "required": true, "schema": {"type": "string"}},
                  {"name": "draft", "in": "path", "required": true, "schema": {"type": "string"}}]}}
 }
}`
