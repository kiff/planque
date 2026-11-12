package fakes

import (
	"net/http"
	"strings"
)

// Billing's docs: the same billing API (billing.go) described the way some
// vendors publish theirs, with no single OpenAPI file to fetch. An llms.txt
// index links one Markdown page per operation, and each page holds a
// fragment of the description in a fenced JSON block: its path and method,
// with the components it uses. A shared page holds the key's scheme and the
// server. The index also links guides that describe no operation.
//
// It is what a catalog builder that drafts from docs is checked against: the
// fragments, merged, describe the three operations billing's own file does,
// create_invoice's final-by-default status included, so a draft built from
// the docs can be compared with one built from the file, and called.

// BillingDocsIndex is billing's llms.txt.
func (v *Vendor) BillingDocsIndex() string { return v.root() + "/billing/llms.txt" }

// billingDocs are the docs' pages, by name: the index links each as
// /billing/docs/<name>.md. {{base}} is billing's API, {{docs}} the pages'.
var billingDocs = map[string]string{
	"getting-started": "# Getting started\n\n" +
		"Create an API key in your settings, then call the API with it. Every request needs the key.\n",

	"authentication": "# Authentication\n\n" +
		"Send your API key in the `X-Billing-Key` header. Requests go to the server below.\n\n" +
		"```json\n" + `{
  "openapi": "3.0.3",
  "info": {"title": "Fake Billing API", "version": "1"},
  "servers": [{"url": "{{base}}"}],
  "security": [{"billingKey": []}],
  "components": {"securitySchemes": {"billingKey": {"type": "apiKey", "in": "header", "name": "X-Billing-Key"}}}
}` + "\n```\n",

	"list-orders": "# List orders\n\n" +
		"Lists your orders, newest first. Filter by status to find the ones to invoice.\n\n" +
		"```json\n" + `{
  "paths": {
    "/orders": {"get": {"operationId": "list_orders", "summary": "List orders",
      "parameters": [{"name": "status", "in": "query", "schema": {"$ref": "#/components/schemas/OrderStatus"}}]}}
  },
  "components": {"schemas": {"OrderStatus": {"type": "string", "enum": ["open", "fulfilled"]}}}
}` + "\n```\n",

	"create-invoice": "# Create an invoice\n\n" +
		"Creates an invoice for an order. **An invoice is finalized unless you send `status: draft`**: a finalized " +
		"invoice is numbered and can't be changed or deleted.\n\n" +
		"```json\n" + `{
  "paths": {
    "/invoices": {"post": {"operationId": "create_invoice", "summary": "Create an invoice for an order",
      "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/NewInvoice"}}}}}}
  },
  "components": {"schemas": {"NewInvoice": {"type": "object", "required": ["order"],
    "properties": {
      "order": {"type": "string", "description": "The order to invoice"},
      "status": {"type": "string", "enum": ["draft", "finalized"], "default": "finalized",
        "description": "draft: kept for review, can be changed or deleted. finalized: numbered and final, can't be changed or deleted. Left out, the invoice is finalized."}}}}}
}` + "\n```\n",

	"send-invoice": "# Send an invoice\n\n" +
		"Emails an invoice to the customer.\n\n" +
		"```json\n" + `{
  "paths": {
    "/invoices/{invoice}/send": {"post": {"operationId": "send_invoice", "summary": "Email an invoice to the customer",
      "parameters": [{"name": "invoice", "in": "path", "required": true, "schema": {"type": "string"}}]}}
  }
}` + "\n```\n",
}

// billingDocsIndex is the llms.txt, as the convention has it: a title, a
// summary, then sections of links with a line each.
const billingDocsIndex = `# Fake Billing

> Invoicing for your orders. The API reference is one page per operation, each with its OpenAPI fragment.

## Guides

- [Getting started]({{docs}}/getting-started.md): your first API call

## API reference

- [Authentication]({{docs}}/authentication.md): the key, and the server
- [List orders]({{docs}}/list-orders.md): find the orders to invoice
- [Create an invoice]({{docs}}/create-invoice.md): finalized unless status is draft
- [Send an invoice]({{docs}}/send-invoice.md): email it to the customer
`

// billingDocsRoutes adds the docs to mux. They are public, as docs are: no key.
func (v *Vendor) billingDocsRoutes(mux *http.ServeMux) {
	// Filled per request: a test server's address is known only once it runs.
	fill := func(s string) string {
		return strings.NewReplacer("{{base}}", v.BillingURL(), "{{docs}}", v.root()+"/billing/docs").Replace(s)
	}
	mux.HandleFunc("GET /billing/llms.txt", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(fill(billingDocsIndex)))
	})
	mux.HandleFunc("GET /billing/docs/{page}", func(w http.ResponseWriter, r *http.Request) {
		name, ok := strings.CutSuffix(r.PathValue("page"), ".md")
		page, found := billingDocs[name]
		if !ok || !found {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		_, _ = w.Write([]byte(fill(page)))
	})
}
