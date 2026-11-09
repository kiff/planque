package connector

import (
	"encoding/json"
	"testing"
)

// A tool a vendor tags read-only that can also change things is a write: the
// nine GraphQL tools of one managed toolkit, tagged read-only, each say they
// run a custom document "that may perform mutations".
func TestAToolThatCanChangeThingsIsNeverReadOnly(t *testing.T) {
	graphQL := map[string]string{
		"SHOPIFY_GRAPH_QL_COLLECTIONS":            "List, retrieve, or count Shopify collections using GraphQL presets, or execute a custom Admin GraphQL document, including mutations.",
		"SHOPIFY_GRAPH_QL_CUSTOMERS":              "List, search, retrieve, or count Shopify customers using GraphQL presets, or execute a custom Admin GraphQL document, including mutations.",
		"SHOPIFY_GRAPH_QL_DISCOUNTS":              "List, search, or retrieve Shopify discounts using GraphQL presets, or execute a custom Admin GraphQL document, including mutations.",
		"SHOPIFY_GRAPH_QL_DRAFT_ORDERS":           "Lists, retrieves, and counts Shopify draft orders using GraphQL, or executes a custom GraphQL document that may perform mutations.",
		"SHOPIFY_GRAPH_QL_INVENTORY":              "Lists and retrieves Shopify inventory items and their levels using GraphQL, or executes a custom GraphQL document that may perform mutations.",
		"SHOPIFY_GRAPH_QL_METAFIELDS_METAOBJECTS": "Reads Shopify metafields, metafield definitions, and metaobjects using GraphQL presets, or executes a custom GraphQL document that may perform mutations.",
		"SHOPIFY_GRAPH_QL_ORDERS":                 "Lists, retrieves, counts, and searches Shopify orders using GraphQL, or executes a custom GraphQL document that may perform mutations.",
		"SHOPIFY_GRAPH_QL_SHOP":                   "Reads Shopify shop and app installation metadata using GraphQL, or executes a custom GraphQL document that may perform mutations.",
		"SHOPIFY_GRAPH_QL_WEBHOOKS":               "Lists, retrieves, and counts Shopify webhook subscriptions using GraphQL, or executes a custom GraphQL document that may perform mutations.",
	}
	plain := json.RawMessage(`{"type":"object","properties":{"limit":{"type":"integer"}}}`)
	for name, desc := range graphQL {
		if !MayChange(name, desc, plain) {
			t.Errorf("%s: read-only", name)
		}
		// Its words alone say so, whatever it's called.
		if !MayChange("list_orders", desc, plain) {
			t.Errorf("%s's description alone: read-only", name)
		}
	}
	for _, c := range []struct{ why, name, desc, schema string }{
		{"a GraphQL query argument", "list_orders", "Lists orders.", `{"type":"object","properties":{"query":{"type":"string","description":"A GraphQL document to run; mutations are allowed."}}}`},
		{"a document argument", "read_things", "Reads things.", `{"type":"object","properties":{"document":{"type":"string"}}}`},
		{"a raw request body", "get_item", "Gets an item.", `{"type":"object","properties":{"body":{"type":"object"}}}`},
		{"a name that runs what it's given", "GRAPHQL_EXECUTE", "Reads data.", `{"type":"object","properties":{}}`},
		{"a raw request", "vendor_api", "Sends a raw HTTP request to the vendor's API.", `{"type":"object","properties":{}}`},
		{"a schema that can't be read", "list_orders", "Lists orders.", `{"properties":`},
	} {
		if !MayChange(c.name, c.desc, json.RawMessage(c.schema)) {
			t.Errorf("%s: read-only", c.why)
		}
	}
}

// An ordinary list or get stays a read, a search's query included: its
// words filter what it reads, they aren't run.
func TestAnOrdinaryReadStaysARead(t *testing.T) {
	for _, c := range []struct{ name, desc, schema string }{
		{"GMAIL_FETCH_EMAILS", "Lists Gmail messages with query and label filtering, pagination, and optional message details.",
			`{"type":"object","properties":{"query":{"type":"string","description":"Gmail advanced search query (e.g., 'from:user subject:meeting')."},"max_results":{"type":"integer"}}}`},
		{"SLACK_FETCH_CONVERSATION_HISTORY", "Fetches a chronological list of messages and events from a specified Slack conversation.", `{"type":"object","properties":{"channel":{"type":"string"},"limit":{"type":"integer"}}}`},
		{"STRIPE_LIST_REFUNDS", "Lists Stripe refunds, sorted by creation date descending (newest first).", `{"type":"object","properties":{"limit":{"type":"integer"}}}`},
		{"QUICKBOOKS_QUERY_INVOICES", "Query QuickBooks invoices with custom SQL-like filters or customer, transaction-date and payment-status filters, with pagination.",
			`{"type":"object","properties":{"query":{"type":"string","description":"SQL-like query string to filter invoices."}}}`},
		{"SLACK_SEARCH_MESSAGES", "Workspace-wide Slack message search with date ranges and filters.", `{"type":"object","properties":{"query":{"type":"string","description":"Search query supporting various modifiers."}}}`},
		{"GORGIAS_LIST_TICKETS", "Lists tickets in Gorgias with various filtering options.", `{"type":"object","properties":{"order_by":{"type":"string"}}}`},
		{"stripe_api_read", "Read data from any Stripe API GET operation.", `{"type":"object","properties":{"operation":{"type":"string"}}}`},
		{"list_charges", "Lists charges.", ``},
	} {
		if MayChange(c.name, c.desc, json.RawMessage(c.schema)) {
			t.Errorf("%s: not a read", c.name)
		}
	}
}
