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
		// Probes on the first version (kiff/planque#118 review).
		{"RAW in the name", "X_RAW", "Reads.", `{}`},
		{"RAW and a query in the name", "SHOPIFY_RAW_QUERY", "Reads.", `{}`},
		{"CUSTOM in the name", "NOTION_CUSTOM_ACTION", "Reads.", `{}`},
		{"GQL in the name", "X_GQL", "Reads.", `{}`},
		{"a body one level down", "X_FETCH", "Fetches.", `{"type":"object","properties":{"request":{"type":"object","properties":{"body":{"type":"object"},"method":{"type":"string"}}}}}`},
		{"a body in an array's items", "X_BATCH", "Reads.", `{"type":"object","properties":{"calls":{"type":"array","items":{"type":"object","properties":{"payload":{"type":"object"}}}}}}`},
		{"a body in one of the alternatives", "X_READ", "Reads.", `{"type":"object","properties":{"input":{"anyOf":[{"type":"string"},{"type":"object","properties":{"document":{"type":"string"}}}]}}}`},
		{"a method beside a path", "X_SEND_REQUEST", "Sends.", `{"type":"object","properties":{"method":{"type":"string"},"path":{"type":"string"},"data":{"type":"object"}}}`},
		{"a method beside a url", "X_LIST", "Lists.", `{"type":"object","properties":{"method":{"type":"string"},"url":{"type":"string"}}}`},
		{"a sql argument", "X_RUN_SQL", "Runs.", `{"type":"object","properties":{"sql":{"type":"string"}}}`},
		{"a statement argument", "X_SQL_QUERY", "Execute SQL against the warehouse.", `{"type":"object","properties":{"statement":{"type":"string"}}}`},
		{"a script argument", "X_RUN_SCRIPT", "Runs.", `{"type":"object","properties":{"script":{"type":"string"}}}`},
		{"a code argument", "X_EVAL", "Evaluates.", `{"type":"object","properties":{"code":{"type":"string"}}}`},
		{"a query that doesn't say what it takes", "SHOPIFY_ADMIN_QUERY", "Reads the store.", `{"type":"object","properties":{"query":{"type":"string"}}}`},
		{"an argument given by reference", "X_FETCH", "Fetches.", `{"type":"object","properties":{"request":{"$ref":"#/$defs/Req"}},"$defs":{"Req":{"type":"object","properties":{"body":{"type":"object"}}}}}`},
		{"a schema given by reference", "X_GET", "Gets.", `{"$ref":"#/$defs/Args","$defs":{"Args":{"type":"object","properties":{"id":{"type":"string"}}}}}`},
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
		// A vendor's custom fields, types, roles and emoji are things it stores:
		// listing them is a read (19 Nov, the wave 1a apps' getters).
		{"ASANA_GET_CUSTOM_FIELD", "Returns the complete record for a single custom field.", `{"type":"object","properties":{"custom_field_gid":{"type":"string"}}}`},
		{"ASANA_GET_CUSTOM_TYPES", "Returns a list of custom types in a workspace.", `{"type":"object","properties":{"workspace":{"type":"string"}}}`},
		{"JIRA_LIST_CUSTOM_FIELD_CONTEXTS", "Returns a paginated list of contexts for a custom field.", `{"type":"object","properties":{"fieldId":{"type":"string"}}}`},
		{"TRELLO_GET_MEMBER_CUSTOM_EMOJI", "Gets a member's custom emoji.", `{"type":"object","properties":{"id":{"type":"string"}}}`},
		// "arbitrary requested fields" is a search that returns ids, not a request to run.
		{"JIRA_SEARCH_FOR_ISSUES_USING_JQL_GET", "Searches for issues using JQL and returns issue ids, rather than arbitrary requested fields.", `{"type":"object","properties":{"jql":{"type":"string","description":"A JQL search."}}}`},
	} {
		if MayChange(c.name, c.desc, json.RawMessage(c.schema)) {
			t.Errorf("%s: not a read", c.name)
		}
	}
}

// "custom" still holds a tool back where it names what the tool runs, and an
// arbitrary request in the description still does.
func TestCustomBesideWhatARunsStillMayChange(t *testing.T) {
	for _, c := range []struct{ name, desc string }{
		{"ASANA_CUSTOM_REQUEST", "Sends a request."},
		{"X_CUSTOM_API_CALL", "Calls the API."},
		{"X_RUN_CUSTOM_QUERY", "Runs a query."},
		{"X_CUSTOM_GRAPH_QL", "Runs GraphQL."},
		{"X_EXECUTE_CUSTOM_FIELD", "Reads a field."},
		{"X_LIST_THINGS", "Runs an arbitrary request against the API."},
		{"X_LIST_THINGS", "Executes a custom document."},
	} {
		if !MayChange(c.name, c.desc, nil) {
			t.Errorf("%s (%s): read as a read", c.name, c.desc)
		}
	}
}
