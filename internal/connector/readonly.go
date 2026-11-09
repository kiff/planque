package connector

import (
	"encoding/json"
	"regexp"
	"strings"
)

// A vendor's read-only mark (MCP's readOnlyHint, a managed provider's tag)
// is the vendor's word, and some vendors give it to a tool that can also
// change things: a GraphQL tool tagged read-only "executes a custom GraphQL
// document that may perform mutations". KIFF, the gateway and the runtime
// treat a read-only tool as one that only reads: the owner is told "reads
// only", Setup may call it, a paused connection still lets it through, it is
// never read back and never counted. So a tool is read-only only when the
// vendor marks it so and nothing about it says it can change things: unknown
// means a write.

// mayChangeName is a tool name that says it runs whatever it's given: a
// GraphQL document, a command, a raw or custom request.
var mayChangeName = regexp.MustCompile(`(?i)(^|[_\-.])(graph_?ql|execute|proxy|raw_request|custom_(api|request|query)|api_call|http_request)([_\-.]|$)`)

// mayChangeWords is a description that says the tool can change things, or
// runs any request or document it's given.
var mayChangeWords = regexp.MustCompile(`(?i)including mutations|may perform mutations|performs? mutations|(executes?|runs?) (a |any )?(custom|raw|arbitrary) |arbitrary (request|query|call|endpoint|operation)|any (rest |http |api )?(endpoint|request|call)\b|custom (graphql|request|document|api call)|raw (http |api )?request`)

// mayChangeArgs are arguments that take a free-form document, body or raw
// request: whatever is in them is what the vendor runs.
var mayChangeArgs = map[string]bool{
	"document": true, "mutation": true, "graphql": true, "graphql_query": true, "graphql_document": true,
	"body": true, "raw": true, "raw_request": true, "request_body": true, "payload": true, "http_method": true,
}

// graphQLWords is a query argument's description that says the query is a
// GraphQL document, which can hold a mutation, not a search's words.
var graphQLWords = regexp.MustCompile(`(?i)graphql|mutation`)

// MayChange reports whether anything about a tool says it can change things,
// whatever the vendor's read-only mark: its name (GRAPH_QL, EXECUTE, a raw
// or custom request), its description ("including mutations", any request),
// or a free-form document, body or raw request among its arguments, a query
// argument that takes GraphQL included. A search's query ("from:ana
// is:unread") is a read's argument, and stays one.
func MayChange(name, description string, inputSchema json.RawMessage) bool {
	if mayChangeName.MatchString(name) || mayChangeWords.MatchString(description) {
		return true
	}
	var schema struct {
		Properties map[string]struct {
			Description string `json:"description"`
		} `json:"properties"`
	}
	if len(inputSchema) > 0 && json.Unmarshal(inputSchema, &schema) != nil {
		return true // a schema that can't be read says nothing it only reads
	}
	for arg, p := range schema.Properties {
		a := strings.ToLower(arg)
		if mayChangeArgs[a] {
			return true
		}
		if (a == "query" || a == "q") && graphQLWords.MatchString(p.Description) {
			return true
		}
	}
	return false
}
