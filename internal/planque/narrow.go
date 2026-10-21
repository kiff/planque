package planque

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/kiff"
)

const maxPathParts = 4

var argTypes = []string{"string", "integer", "number", "boolean"}

func pathParts(p string) ([]string, bool) {
	parts := strings.Split(p, ".")
	if len(parts) == 0 || len(parts) > maxPathParts {
		return nil, false
	}
	for _, s := range parts {
		if s == "" || len(s) > 64 || strings.ContainsAny(s, " \t\r\n") {
			return nil, false
		}
	}
	return parts, true
}

// checkNarrowing says what is wrong with a tool's fixed and lifted
// arguments, or "".
func checkNarrowing(t ToolRequest) string {
	paths := map[string]bool{}
	for p, v := range t.Fixed {
		if _, ok := pathParts(p); !ok || !json.Valid(v) {
			return "fixed argument " + p + " must be a dotted path with a JSON value"
		}
		paths[p] = true
	}
	if len(t.Fixed) > 0 && len(t.Arguments) == 0 {
		return t.Tool + " has fixed arguments, so it must also name the arguments the agent may set"
	}
	for name, a := range t.Arguments {
		switch {
		case !toolName.MatchString(name) || name == IdempotencyArg:
			return "argument " + name + " is not a valid name"
		case !slices.Contains(argTypes, a.Type):
			return "argument " + name + " must be a string, integer, number or boolean"
		}
		if _, ok := pathParts(a.Path); !ok {
			return "argument " + name + " needs a dotted path in the vendor's arguments"
		}
		for p := range paths {
			if p == a.Path || strings.HasPrefix(p, a.Path+".") || strings.HasPrefix(a.Path, p+".") {
				return "argument " + name + " overlaps another argument at " + p
			}
		}
		paths[a.Path] = true
	}
	if t.AmountArg != "" && len(t.Arguments) > 0 && t.Arguments[t.AmountArg].Type != "integer" {
		return "amount_argument " + t.AmountArg + " must be one of the tool's integer arguments"
	}
	if r := t.ResourceArg; r != "" {
		switch {
		case !toolName.MatchString(r) || r == t.AmountArg || r == IdempotencyArg:
			return "resource_argument " + r + " is not a valid name"
		case len(t.Arguments) > 0 && t.Arguments[r].Type != "string":
			return "resource_argument " + r + " must be one of the tool's string arguments"
		}
	}
	return ""
}

// checkLookups says what is wrong with the tools' unit lookups, or "": each
// names a read of the same request, by its name or its vendor tool, and the
// read's string argument that takes the resource.
func checkLookups(tools []ToolRequest) string {
	for k, t := range tools {
		l := t.UnitLookup
		if l == nil {
			continue
		}
		if t.AmountArg == "" || t.AmountUnit == "" || t.ResourceArg == "" {
			return t.Tool + "'s unit_lookup needs its amount_argument, amount_unit and resource_argument"
		}
		i := lookupIndex(tools, l.Tool)
		if i < 0 || i == k || !tools[i].ReadOnly {
			return t.Tool + "'s unit_lookup must name a read-only tool of the same request"
		}
		read := tools[i]
		if !toolName.MatchString(l.Argument) || (len(read.Arguments) > 0 && read.Arguments[l.Argument].Type != "string") {
			return t.Tool + "'s unit_lookup.argument must be one of " + read.Tool + "'s string arguments"
		}
		if _, ok := pathParts(l.Path); !ok {
			return t.Tool + "'s unit_lookup.path must be a dotted path in the read's answer, such as currency"
		}
	}
	return ""
}

// lookupIndex is the tool a unit lookup names: by its name, else its
// vendor tool. -1 when none, or the vendor tool is asked for twice.
func lookupIndex(tools []ToolRequest, name string) int {
	at := -1
	for i, t := range tools {
		if t.Name == name && name != "" {
			return i
		}
		if t.Tool == name {
			if at >= 0 {
				return -1
			}
			at = i
		}
	}
	return at
}

// resolveLookups points each unit lookup at its read's name in KIFF, once
// every tool has one.
func resolveLookups(tools []ToolRequest) {
	for i := range tools {
		if l := tools[i].UnitLookup; l != nil {
			if j := lookupIndex(tools, l.Tool); j >= 0 {
				tools[i].UnitLookup = &kiff.UnitLookup{Tool: tools[j].Name, Argument: l.Argument, Path: l.Path}
			}
		}
	}
}

// narrowed reports whether the agent sees planque's own argument list.
func (t ToolRequest) narrowed() bool { return len(t.Arguments) > 0 }

// schema is the input schema of a narrowed tool: its arguments, nothing
// else.
func (t ToolRequest) schema() json.RawMessage {
	props := map[string]any{}
	required := []string{}
	for _, name := range slices.Sorted(maps.Keys(t.Arguments)) {
		a := t.Arguments[name]
		p := map[string]any{"type": a.Type}
		if a.Description != "" {
			p["description"] = a.Description
		}
		props[name] = p
		if a.Required {
			required = append(required, name)
		}
	}
	b, _ := json.Marshal(map[string]any{"type": "object", "properties": props, "required": required, "additionalProperties": false})
	return b
}

// vendorArgs builds what the vendor receives from the agent's arguments
// (planque's operation argument already removed). It refuses an argument
// the grant does not name, of the wrong type, or a missing required one.
func (t ToolRequest) vendorArgs(args map[string]any) (json.RawMessage, error) {
	out := map[string]any{}
	if !t.narrowed() {
		maps.Copy(out, args)
	} else {
		for name, v := range args {
			a, ok := t.Arguments[name]
			if !ok {
				return nil, fmt.Errorf("argument %s is not part of this grant", name)
			}
			if !hasType(v, a.Type) {
				return nil, fmt.Errorf("argument %s must be a %s", name, a.Type)
			}
			parts, _ := pathParts(a.Path)
			set(out, parts, v)
		}
		for name, a := range t.Arguments {
			if _, ok := args[name]; a.Required && !ok {
				return nil, fmt.Errorf("argument %s is required", name)
			}
		}
	}
	for p, raw := range t.Fixed {
		var v any
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		_ = dec.Decode(&v)
		parts, _ := pathParts(p)
		set(out, parts, v)
	}
	return json.Marshal(out)
}

func set(m map[string]any, parts []string, v any) {
	for _, k := range parts[:len(parts)-1] {
		next, ok := m[k].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[k] = next
		}
		m = next
	}
	m[parts[len(parts)-1]] = v
}

func hasType(v any, typ string) bool {
	switch typ {
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "integer":
		n, ok := v.(json.Number)
		if !ok {
			return false
		}
		_, err := n.Int64()
		return err == nil
	case "number":
		_, ok := v.(json.Number)
		return ok
	}
	return false
}

// narrowingText says, in the owner's words, how a tool is narrowed. A
// pinned argument the vendor's catalog describes comes with the vendor's own
// words for it (kiff/tafs#89), so the owner reads what the value means
// whatever the vendor calls it, and the words aren't the asking agent's.
func narrowingText(t ToolRequest, catalog ...CatalogTool) string {
	if !t.narrowed() {
		return ""
	}
	var schema json.RawMessage
	if i := slices.IndexFunc(catalog, func(c CatalogTool) bool { return c.Name == t.Tool }); i >= 0 {
		schema = catalog[i].InputSchema
	}
	var fixed []string
	for _, p := range slices.Sorted(maps.Keys(t.Fixed)) {
		f := p + " = " + string(t.Fixed[p])
		if d := argDescription(schema, p); d != "" {
			f += ` (the vendor says: "` + d + `")`
		}
		fixed = append(fixed, f)
	}
	s := " (" + t.Tool
	if len(fixed) > 0 {
		s += " with " + strings.Join(fixed, ", ")
	}
	return s + "; the taf sets only " + strings.Join(slices.Sorted(maps.Keys(t.Arguments)), ", ") + ")"
}

// argDescription is the vendor's description of the argument at a dotted
// path in a tool's input schema, on one line and at most 200 characters; ""
// when the schema doesn't describe it.
func argDescription(schema json.RawMessage, path string) string {
	for _, name := range strings.Split(path, ".") {
		var s struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		if json.Unmarshal(schema, &s) != nil {
			return ""
		}
		next, ok := s.Properties[name]
		if !ok {
			return ""
		}
		schema = next
	}
	var leaf struct {
		Description string `json:"description"`
	}
	if json.Unmarshal(schema, &leaf) != nil {
		return ""
	}
	d := strings.Join(strings.Fields(leaf.Description), " ")
	d = strings.ReplaceAll(d, `"`, "'")
	if r := []rune(d); len(r) > 200 {
		d = strings.TrimSpace(string(r[:199])) + "…"
	}
	return d
}
