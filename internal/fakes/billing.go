package fakes

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Billing: the fake vendor's third API, invoicing reached by an API key, as
// the desk is. One argument decides whether a write is final: an invoice is
// made as a draft, kept for review and changed freely, or finalized,
// numbered and unchangeable. Left out, it is finalized, as some vendors'
// invoicing APIs do. Sending an invoice is a write of its own. It is what
// kiff/agents' e2e checks a grant pins: a draft-only taf has create_invoice
// with status fixed to draft, and no send_invoice.

// BillingKeyHeader is where billing takes its API key: the vendor's one key,
// which the desk takes too (DeskKey), each in its own header.
const BillingKeyHeader = "X-Billing-Key"

// kvBilling keeps billing's writes: SK when (ns) and an id, Doc a VendorCall.
const kvBilling = "fakevendor#billing"

// Invoice statuses, as create_invoice takes them.
const (
	InvoiceDraft     = "draft"
	InvoiceFinalized = "finalized"
)

// billingOrders are the orders to invoice: two fulfilled, one not yet.
var billingOrders = []map[string]any{
	{"id": "or_2001", "customer": "Ana Shop", "amount": 12000, "currency": "eur", "status": "fulfilled"},
	{"id": "or_2002", "customer": "Ben Studio", "amount": 4500, "currency": "eur", "status": "fulfilled"},
	{"id": "or_2003", "customer": "Chloe Café", "amount": 8000, "currency": "eur", "status": "open"},
}

// BillingURL is billing's API, and BillingDescription its OpenAPI description.
func (v *Vendor) BillingURL() string         { return v.root() + "/billing/v1" }
func (v *Vendor) BillingDescription() string { return v.root() + "/billing/openapi.json" }

// billingRoutes adds billing to mux.
func (v *Vendor) billingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /billing/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(strings.ReplaceAll(billingOpenAPI, "{{base}}", v.BillingURL())))
	})
	mux.HandleFunc("GET /billing/v1/orders", v.billing(func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		out := []map[string]any{}
		for _, o := range billingOrders {
			if status == "" || o["status"] == status {
				out = append(out, o)
			}
		}
		writeJSON(w, 200, map[string]any{"data": out})
	}))
	mux.HandleFunc("POST /billing/v1/invoices", v.billing(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			writeJSON(w, 400, map[string]string{"error": "the body must be a JSON object"})
			return
		}
		order, _ := body["order"].(string)
		if _, ok := billingOrder(order); !ok {
			writeJSON(w, 404, map[string]string{"error": "no such order"})
			return
		}
		status := InvoiceFinalized
		if s, set := body["status"]; set {
			if s != InvoiceDraft && s != InvoiceFinalized {
				writeJSON(w, 422, map[string]string{"error": "status is draft or finalized"})
				return
			}
			status = s.(string)
		}
		id := next("in")
		if !v.billingWrite(w, r, VendorCall{Tool: "create_invoice", Args: map[string]any{"order": order, "status": status}}, id) {
			return
		}
		writeJSON(w, 201, map[string]any{"id": id, "order": order, "status": status})
	}))
	mux.HandleFunc("POST /billing/v1/invoices/{invoice}/send", v.billing(func(w http.ResponseWriter, r *http.Request) {
		invoice := r.PathValue("invoice")
		if !strings.HasPrefix(invoice, "in_") {
			writeJSON(w, 404, map[string]string{"error": "no such invoice"})
			return
		}
		if !v.billingWrite(w, r, VendorCall{Tool: "send_invoice", Args: map[string]any{"invoice": invoice}}, next("sd")) {
			return
		}
		writeJSON(w, 200, map[string]any{"id": invoice, "sent": true})
	}))
}

// billing checks the key before h.
func (v *Vendor) billing(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v.DeskKey == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get(BillingKeyHeader)), []byte(v.DeskKey)) != 1 {
			writeJSON(w, 401, map[string]string{"error": "a valid " + BillingKeyHeader + " is required"})
			return
		}
		h(w, r)
	}
}

// billingWrite records one write; false when it couldn't, answered already.
func (v *Vendor) billingWrite(w http.ResponseWriter, r *http.Request, call VendorCall, id string) bool {
	doc, _ := json.Marshal(call)
	if err := v.keep(r.Context(), kvBilling, fmt.Sprintf("%020d#%s", time.Now().UnixNano(), id), doc, 7*24*time.Hour); err != nil {
		writeJSON(w, 500, map[string]string{"error": "not recorded"})
		return false
	}
	v.mu.Lock()
	v.BillingCalls = append(v.BillingCalls, call)
	v.mu.Unlock()
	return true
}

func billingOrder(id string) (map[string]any, bool) {
	for _, o := range billingOrders {
		if o["id"] == id {
			return o, true
		}
	}
	return nil, false
}

// billingOpenAPI is billing's description, as planque reads one: JSON,
// OpenAPI 3.0, one server, the key in a header. create_invoice's status
// says what each value does and what leaving it out does, as a vendor's
// own reference would.
const billingOpenAPI = `{
 "openapi": "3.0.3",
 "info": {"title": "Fake Billing API", "version": "1"},
 "servers": [{"url": "{{base}}"}],
 "security": [{"billingKey": []}],
 "components": {"securitySchemes": {"billingKey": {"type": "apiKey", "in": "header", "name": "X-Billing-Key"}}},
 "paths": {
  "/orders": {"get": {"operationId": "list_orders", "summary": "List orders",
   "parameters": [{"name": "status", "in": "query", "schema": {"type": "string", "enum": ["open", "fulfilled"]}}]}},
  "/invoices": {"post": {"operationId": "create_invoice", "summary": "Create an invoice for an order",
   "requestBody": {"required": true, "content": {"application/json": {"schema": {"type": "object", "required": ["order"],
    "properties": {
     "order": {"type": "string", "description": "The order to invoice"},
     "status": {"type": "string", "enum": ["draft", "finalized"], "default": "finalized",
      "description": "draft: kept for review, can be changed or deleted. finalized: numbered and final, can't be changed or deleted. Left out, the invoice is finalized."}}}}}}}},
  "/invoices/{invoice}/send": {"post": {"operationId": "send_invoice", "summary": "Email an invoice to the customer",
   "parameters": [{"name": "invoice", "in": "path", "required": true, "schema": {"type": "string"}}]}}
 }
}`
