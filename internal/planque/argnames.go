package planque

import (
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
)

// An agent's model takes an input schema's property names only within
// modelArg. A vendor's own names can be outside it, such as Qonto's query
// parameters "filter[status]", and a tool listing one is refused whole, so
// the run fails before anything is sent. The relay gives the agent a name
// within it, and puts the vendor's own name back on each call.
var modelArg = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,64}$`)

// argNames maps each of props' names to the name the agent is given: the
// same name when the model takes it, else a spelling it takes
// ("filter[status]" is "filter_status", "filter[]" is "filter"), never one
// another of props already has. The same props always give the same names.
func argNames(props map[string]any) map[string]string {
	names := slices.Sorted(maps.Keys(props))
	taken := map[string]bool{}
	for _, n := range names {
		if modelArg.MatchString(n) {
			taken[n] = true
		}
	}
	out := make(map[string]string, len(names))
	for _, n := range names {
		if modelArg.MatchString(n) {
			out[n] = n
			continue
		}
		base := safeArg(n)
		s := base
		for i := 2; taken[s]; i++ {
			s = fmt.Sprintf("%s_%d", base, i)
		}
		taken[s] = true
		out[n] = s
	}
	return out
}

// safeArg spells n within modelArg: "]" dropped, every other character
// outside it an underscore, runs of them one, none at either end.
func safeArg(n string) string {
	var b strings.Builder
	for _, r := range n {
		switch {
		case r == ']':
		case r < 128 && (r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-'):
			b.WriteRune(r)
		default:
			if !strings.HasSuffix(b.String(), "_") {
				b.WriteByte('_')
			}
		}
	}
	s := strings.Trim(b.String(), "_")
	if s == "" {
		s = "arg"
	}
	if len(s) > 56 {
		s = s[:56]
	}
	return s
}

// forAgent is schema with its property names (and required) as the agent
// is given them.
func forAgent(schema map[string]any) {
	props, _ := schema["properties"].(map[string]any)
	if len(props) == 0 {
		return
	}
	names := argNames(props)
	renamed := make(map[string]any, len(props))
	for n, v := range props {
		renamed[names[n]] = v
	}
	schema["properties"] = renamed
	if req, ok := schema["required"].([]any); ok {
		out := make([]any, len(req))
		for i, r := range req {
			if s, ok := r.(string); ok && names[s] != "" {
				out[i] = names[s]
			} else {
				out[i] = r
			}
		}
		schema["required"] = out
	}
}

// fromAgent puts the vendor's own names back on a call's arguments, by the
// granted tool's schema: each argument named as forAgent gave it is renamed
// to the property it stands for. Any other name is left as sent.
func fromAgent(inputSchema json.RawMessage, args map[string]any) {
	var schema struct {
		Properties map[string]any `json:"properties"`
	}
	if json.Unmarshal(inputSchema, &schema) != nil {
		return
	}
	for vendor, agent := range argNames(schema.Properties) {
		if vendor == agent {
			continue
		}
		if v, ok := args[agent]; ok {
			delete(args, agent)
			args[vendor] = v
		}
	}
}
