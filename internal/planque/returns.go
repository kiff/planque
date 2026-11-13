package planque

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
)

// maxReturns bounds the ids and links a catalog tool names.
const maxReturns = 20

// plainPath is a path in an answer as returns names it.
var plainPath = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// returnsOf are the ids and links a tool's answer gives back, by their
// path in it, read from the answer's JSON schema (properties, an array's
// items, three levels down): a property named id, or ending in Id, _id,
// url, link, or Link. Nil when the schema says nothing of them.
func returnsOf(schema json.RawMessage) []string {
	if len(schema) == 0 {
		return nil
	}
	var root map[string]any
	if json.Unmarshal(schema, &root) != nil {
		return nil
	}
	var out []string
	var walk func(node map[string]any, path string, depth int)
	walk = func(node map[string]any, path string, depth int) {
		if depth > 3 || len(out) >= maxReturns {
			return
		}
		if items, ok := node["items"].(map[string]any); ok {
			walk(items, path, depth)
		}
		props, _ := node["properties"].(map[string]any)
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			child, _ := props[k].(map[string]any)
			p := k
			if path != "" {
				p = path + "." + k
			}
			// The vendor's names, read by Coach: plain ones only, bounded.
			if len(p) > 128 || !plainPath.MatchString(p) {
				continue
			}
			if identifying(k) && len(out) < maxReturns {
				out = append(out, p)
			}
			if child != nil {
				walk(child, p, depth+1)
			}
		}
	}
	walk(root, "", 0)
	return out
}

// identifying says whether a property's name is an id or a link.
func identifying(name string) bool {
	l := strings.ToLower(name)
	return l == "id" || strings.HasSuffix(name, "Id") || strings.HasSuffix(name, "ID") || strings.HasSuffix(l, "_id") ||
		strings.HasSuffix(l, "url") || strings.HasSuffix(l, "link")
}
