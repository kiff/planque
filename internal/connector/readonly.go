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

// mayChangeName is a tool name that says it may run whatever it's given: a
// GraphQL document (GRAPH_QL, GQL), a command, anything raw or custom. When
// unsure, it's a write: a vendor's "custom properties" or "raw content"
// reads come out writes too, which only costs Setup a read it can't make.
var mayChangeName = regexp.MustCompile(`(?i)(^|[_\-.])(graph_?ql|gql|execute|proxy|raw|custom|api_call|http_request)([_\-.]|$)`)

// mayChangeWords is a description that says the tool can change things, or
// runs any request or document it's given.
var mayChangeWords = regexp.MustCompile(`(?i)including mutations|may perform mutations|performs? mutations|(executes?|runs?) (a |any )?(custom|raw|arbitrary) |arbitrary (request|query|call|endpoint|operation)|any (rest |http |api )?(endpoint|request|call)\b|custom (graphql|request|document|api call)|raw (http |api )?request`)

// mayChangeArgs are arguments that may take a free-form document, program,
// body or raw request: whatever is in them is what the vendor runs. When
// unsure, it's a write: a code argument is often a discount code to look up,
// and is counted all the same.
var mayChangeArgs = map[string]bool{
	"document": true, "mutation": true, "graphql": true, "graphql_query": true, "graphql_document": true,
	"body": true, "raw": true, "raw_request": true, "request_body": true, "payload": true, "http_method": true,
	"sql": true, "statement": true, "script": true, "command": true, "code": true,
}

// requestTarget is an argument that names where a raw request goes: beside
// a method, the tool sends whatever request it's given.
var requestTarget = map[string]bool{"url": true, "path": true, "endpoint": true}

// graphQLWords is a query argument's description that says the query is a
// GraphQL document, which can hold a mutation, not a search's words.
var graphQLWords = regexp.MustCompile(`(?i)graphql|mutation`)

// MayChange reports whether anything about a tool says it can change things,
// whatever the vendor's read-only mark: its name (GRAPH_QL, EXECUTE, RAW,
// CUSTOM), its description ("including mutations", any request), or
// its arguments, nested ones included: a free-form document, program, body
// or raw request, a method beside a url, path or endpoint, a query that
// takes GraphQL or doesn't say what it takes, or one given by reference. A search's query ("from:ana
// is:unread") is a read's argument, and stays one.
func MayChange(name, description string, inputSchema json.RawMessage) bool {
	if mayChangeName.MatchString(name) || mayChangeWords.MatchString(description) {
		return true
	}
	if len(inputSchema) == 0 {
		return false
	}
	var root schemaNode
	if json.Unmarshal(inputSchema, &root) != nil {
		return true // a schema that can't be read says nothing it only reads
	}
	seen := map[string]bool{}
	if root.mayChange(seen, 0) {
		return true
	}
	if seen["method"] {
		for target := range requestTarget {
			if seen[target] {
				return true
			}
		}
	}
	return false
}

// schemaNode is the part of a JSON Schema MayChange reads: an object's
// properties, an array's items, and the alternatives of anyOf, oneOf and
// allOf, at any depth, and a reference ($ref), which it doesn't follow.
type schemaNode struct {
	Ref         string                `json:"$ref"`
	Description string                `json:"description"`
	Properties  map[string]schemaNode `json:"properties"`
	Items       *schemaNode           `json:"items"`
	AnyOf       []schemaNode          `json:"anyOf"`
	OneOf       []schemaNode          `json:"oneOf"`
	AllOf       []schemaNode          `json:"allOf"`
}

// maxSchemaDepth bounds how deep mayChange reads a schema; past it, the
// tool says nothing it only reads.
const maxSchemaDepth = 32

// mayChange reports whether an argument in n says the tool can change
// things, noting each argument's name in seen for the raw request's shape.
func (n schemaNode) mayChange(seen map[string]bool, depth int) bool {
	// A node given by reference, or too deep, says nothing it only reads.
	if depth > maxSchemaDepth || n.Ref != "" {
		return true
	}
	for arg, p := range n.Properties {
		a := strings.ToLower(arg)
		seen[a] = true
		if mayChangeArgs[a] {
			return true
		}
		// A query that takes GraphQL, or doesn't say what it takes, is a
		// document to run: only one that says it's a search's words is a read's.
		if (a == "query" || a == "q") && (graphQLWords.MatchString(p.Description) || strings.TrimSpace(p.Description) == "") {
			return true
		}
		if p.mayChange(seen, depth+1) {
			return true
		}
	}
	if n.Items != nil && n.Items.mayChange(seen, depth+1) {
		return true
	}
	for _, alts := range [][]schemaNode{n.AnyOf, n.OneOf, n.AllOf} {
		for _, alt := range alts {
			if alt.mayChange(seen, depth+1) {
				return true
			}
		}
	}
	return false
}
