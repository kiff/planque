package planque

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/kiff"
)

const maxPathParts = 4

var argTypes = []string{"string", "integer", "number", "boolean", "array"}

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

// checkPin says what is wrong with one of a tool's fixed arguments, or nil.
// The words are for Coach, which fixes the request: they name the tool and
// the pin, say which part is wrong, and show the form planque wants.
func checkPin(t ToolRequest, p string, v json.RawMessage) *toolRefusal {
	tool := refTool(t)
	if _, ok := pathParts(p); !ok {
		return &toolRefusal{Code: refPinShape, Tool: tool, Field: p, Want: "dotted path → JSON value",
			Msg: fmt.Sprintf("%s: fixed argument %q isn't a dotted path: name the vendor's argument (up to %d names joined by dots, no spaces), and give its value as JSON, like {\"author\": \"urn:li:person:…\"} or {\"content.visibility\": \"PUBLIC\"}", tool, p, maxPathParts)}
	}
	if !json.Valid(v) {
		return &toolRefusal{Code: refPinShape, Tool: tool, Field: p, Want: "dotted path → JSON value",
			Msg: fmt.Sprintf("%s: fixed argument %q's value isn't JSON: write the value itself as JSON, a string in quotes, like {%q: \"urn:li:person:…\"}", tool, p, p)}
	}
	return nil
}

// refTool is the tool a refusal names: the request's name for it, else its
// vendor tool.
func refTool(t ToolRequest) string { return cmp.Or(t.Name, t.Tool) }

// checkNarrowing says what is wrong with a tool's fixed and lifted
// arguments, or nil.
func checkNarrowing(t ToolRequest) *toolRefusal {
	tool := refTool(t)
	arg := func(code, name, want, msg string) *toolRefusal {
		return &toolRefusal{Code: code, Tool: tool, Field: name, Want: want, Msg: msg}
	}
	// A broad write narrowed to an operation that moves money or closes a bill for good (voiding,
	// paying, writing off, deleting or finalizing a bill or a credit):
	// refused before anything is asked, as a managed app's never-served tool.
	if d, op, write, ok := narrowedTo(t); ok && write && neverSent(t.Tool, op) {
		return &toolRefusal{Code: refToolNotServed, Tool: tool, Field: d.OpArg, Want: "another operation",
			Msg: fmt.Sprintf("%s: planque never narrows %s to %s: it moves money or closes a bill for good (voiding, paying, writing off, deleting or finalizing a bill or a credit), so it is refused; leave it out", tool, t.Tool, strings.TrimSpace(op))}
	}
	paths := map[string]bool{}
	// In a set order, so a request names the same pin each time.
	for _, p := range slices.Sorted(maps.Keys(t.Fixed)) {
		if r := checkPin(t, p, t.Fixed[p]); r != nil {
			return r
		}
		paths[p] = true
	}
	for name, a := range t.Arguments {
		switch {
		case !toolName.MatchString(name) || name == IdempotencyArg:
			return arg(refArgumentName, name, "letters, digits and underscores, starting with a letter", "argument "+name+" is not a valid name")
		case !slices.Contains(argTypes, a.Type):
			return arg(refArgumentType, name, "string, integer, number, boolean or array", "argument "+name+" must be a string, integer, number, boolean or array")
		case a.Type == "array" && checkShape(a.Schema, true) != "":
			return arg(refArgumentSchema, name, "an array of plain values, or of objects naming their fields", "argument "+name+"'s schema: "+checkShape(a.Schema, true))
		case a.Type != "array" && len(a.Schema) > 0:
			return arg(refArgumentSchema, name, "no schema", "argument "+name+" is a "+a.Type+": only an array argument has a schema")
		case a.OneOf != nil && (a.Type != "string" || !toolName.MatchString(a.OneOf.Tool) || !pathExpr(a.OneOf.Path)):
			return arg(refArgumentOneOf, name, "a string argument's {tool, path}", "argument "+name+"'s one_of is a string argument's read (a tool of the same request) and the dotted path of its values")
		}
		if _, ok := pathParts(a.Path); !ok {
			return arg(refArgumentPath, name, "dotted path", "argument "+name+" needs a dotted path in the vendor's arguments")
		}
		for p := range paths {
			if p == a.Path || strings.HasPrefix(p, a.Path+".") || strings.HasPrefix(a.Path, p+".") {
				return arg(refArgumentOverlap, name, "", "argument "+name+" overlaps another argument at "+p)
			}
		}
		paths[a.Path] = true
	}
	if t.AmountLookup != nil {
		// KIFF reads the amount from the resource: the agent never sets it.
		if _, set := t.Arguments[t.AmountArg]; set {
			return arg(refAmountArgument, "amount_argument", "not one of the tool's arguments", "amount_argument "+t.AmountArg+" is read by amount_lookup, so it can't also be one of the tool's arguments")
		}
	} else if t.AmountArg != "" && t.narrowed() {
		want := "integer"
		if t.AmountFormat == "decimal_string" {
			want = "string"
		}
		if t.Arguments[t.AmountArg].Type != want {
			return arg(refAmountArgument, "amount_argument", "one of the tool's "+want+" arguments", "amount_argument "+t.AmountArg+" must be one of the tool's "+want+" arguments")
		}
	}
	if r := t.ResourceArg; r != "" {
		switch {
		case !toolName.MatchString(r) || r == t.AmountArg || r == IdempotencyArg:
			return arg(refResourceArgument, "resource_argument", "", "resource_argument "+r+" is not a valid name")
		case t.narrowed() && t.Arguments[r].Type != "string":
			return arg(refResourceArgument, "resource_argument", "one of the tool's string arguments", "resource_argument "+r+" must be one of the tool's string arguments")
		}
	}
	return nil
}

// checkLookups says what is wrong with the tools' unit and amount lookups,
// or nil: each names a read of the same request, by its name or its vendor
// tool, and the read's string argument that takes the resource.
func checkLookups(tools []ToolRequest) *toolRefusal {
	for k, t := range tools {
		if l := t.UnitLookup; l != nil {
			if t.AmountArg == "" || t.AmountUnit == "" || t.ResourceArg == "" {
				return &toolRefusal{Code: refLookup, Tool: refTool(t), Field: "unit_lookup", Want: "amount_argument, amount_unit and resource_argument",
					Msg: t.Tool + "'s unit_lookup needs its amount_argument, amount_unit and resource_argument"}
			}
			if r := checkLookup(tools, k, *l, "unit_lookup", "currency"); r != nil {
				return r
			}
		}
		if l := t.AmountLookup; l != nil {
			if t.AmountArg == "" || t.ResourceArg == "" {
				return &toolRefusal{Code: refLookup, Tool: refTool(t), Field: "amount_lookup", Want: "amount_argument and resource_argument",
					Msg: t.Tool + "'s amount_lookup needs its amount_argument (the name its Card counts) and resource_argument"}
			}
			if r := checkLookup(tools, k, *l, "amount_lookup", "amount"); r != nil {
				return r
			}
		}
	}
	return nil
}

// checkLookup checks one lookup of tools[k]: a read-only tool of the same
// request, its string argument, and a dotted path in its answer.
func checkLookup(tools []ToolRequest, k int, l kiff.UnitLookup, field, example string) *toolRefusal {
	t := tools[k]
	ref := func(sub, want, msg string) *toolRefusal {
		return &toolRefusal{Code: refLookup, Tool: refTool(t), Field: field + sub, Want: want, Msg: msg}
	}
	i := lookupIndex(tools, l.Tool)
	if i < 0 || i == k || !tools[i].ReadOnly {
		return ref(".tool", "a read-only tool of the same request", t.Tool+"'s "+field+" must name a read-only tool of the same request")
	}
	read := tools[i]
	if !toolName.MatchString(l.Argument) || (len(read.Arguments) > 0 && read.Arguments[l.Argument].Type != "string") {
		return ref(".argument", "one of "+read.Tool+"'s string arguments", t.Tool+"'s "+field+".argument must be one of "+read.Tool+"'s string arguments")
	}
	if _, ok := pathParts(l.Path); !ok {
		return ref(".path", "dotted path, such as "+example, t.Tool+"'s "+field+".path must be a dotted path in the read's answer, such as "+example)
	}
	return nil
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

// resolveLookups points each unit and amount lookup at its read's name in
// KIFF, once every tool has one.
func resolveLookups(tools []ToolRequest) {
	for i := range tools {
		for _, l := range []**kiff.UnitLookup{&tools[i].UnitLookup, &tools[i].AmountLookup} {
			if *l == nil {
				continue
			}
			if j := lookupIndex(tools, (*l).Tool); j >= 0 {
				*l = &kiff.UnitLookup{Tool: tools[j].Name, Argument: (*l).Argument, Path: (*l).Path}
			}
		}
	}
}

// narrowed reports whether the agent sees planque's own argument list: the
// arguments it names, or none when every argument the tool takes is pinned
// (a read of fulfilled orders, its status fixed), so the taf sets nothing.
func (t ToolRequest) narrowed() bool { return len(t.Arguments) > 0 || len(t.Fixed) > 0 }

// schema is the input schema of a narrowed tool: its arguments, nothing
// else.
func (t ToolRequest) schema() json.RawMessage { return t.schemaOver(nil) }

// schemaOver is schema with each scalar argument's allowed values, as the
// vendor's input schema lists them at its path (its enum): the agent sees
// "open, pending or solved", not a bare string it has to guess, and a call
// with another value is refused before the vendor gets it (outsideEnum).
func (t ToolRequest) schemaOver(vendor json.RawMessage) json.RawMessage {
	props := map[string]any{}
	required := []string{}
	for _, name := range slices.Sorted(maps.Keys(t.Arguments)) {
		a := t.Arguments[name]
		p := map[string]any{"type": a.Type}
		if a.Type == "array" {
			p = map[string]any{}
			_ = json.Unmarshal(a.Schema, &p)
		} else if e := vendorEnum(vendor, a.Path, a.Type); e != nil {
			p["enum"] = e
		}
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
			if a.Type == "array" {
				if why := fits(v, a.Schema, name); why != "" {
					return nil, fmt.Errorf("argument %s", why)
				}
			} else if !hasType(v, a.Type) {
				return nil, fmt.Errorf("argument %s must be a %s", name, a.Type)
			}
			parts, _ := pathParts(a.Path)
			set(out, parts, v)
		}
		for name, a := range t.Arguments {
			if _, ok := args[name]; a.Required && !ok {
				return nil, fmt.Errorf("argument %s is required", name)
			}
			// A read sends the object an argument goes in even when the
			// agent set none of its arguments: called with {} (a schedule's
			// check), Stripe's read still gets its parameters object, which
			// its tool requires. Never a write's: an empty object there could
			// mean something the owner didn't approve.
			if parts, _ := pathParts(a.Path); t.ReadOnly && len(parts) > 1 {
				parents(out, parts[:len(parts)-1])
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

// parents makes sure the nested objects at path exist in m, empty when new.
// It only adds what's absent: a value already at a key, of any kind, is never
// replaced, so nothing the owner approved is dropped.
func parents(m map[string]any, path []string) {
	for _, k := range path {
		v, ok := m[k]
		if !ok {
			next := map[string]any{}
			m[k] = next
			m = next
			continue
		}
		next, isObject := v.(map[string]any)
		if !isObject {
			return
		}
		m = next
	}
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
	if len(t.Arguments) == 0 {
		return s + "; the taf sets nothing)"
	}
	return s + "; the taf sets only " + strings.Join(slices.Sorted(maps.Keys(t.Arguments)), ", ") + ")"
}

// argDescription is the vendor's description of the argument at a dotted
// path in a tool's input schema, on one line and at most 200 characters; ""
// when the schema doesn't describe it.
func argDescription(schema json.RawMessage, path string) string {
	var leaf struct {
		Description string `json:"description"`
	}
	if json.Unmarshal(vendorLeaf(schema, path), &leaf) != nil {
		return ""
	}
	d := strings.Join(strings.Fields(leaf.Description), " ")
	d = strings.ReplaceAll(d, `"`, "'")
	if r := []rune(d); len(r) > 200 {
		d = strings.TrimSpace(string(r[:199])) + "…"
	}
	return d
}

// vendorLeaf is the schema at a dotted path in a tool's input schema, through
// its properties; nil when the schema doesn't have it.
func vendorLeaf(schema json.RawMessage, path string) json.RawMessage {
	for _, name := range strings.Split(path, ".") {
		var s struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		if json.Unmarshal(schema, &s) != nil {
			return nil
		}
		next, ok := s.Properties[name]
		if !ok {
			return nil
		}
		schema = next
	}
	return schema
}

// maxEnum bounds the allowed values an argument carries from the vendor's
// schema, and the length of each: a longer list, or a value of another type
// than the argument's, isn't carried, and the argument stays as it was.
const (
	maxEnum      = 50
	maxEnumValue = 128
)

// vendorEnum is the vendor's allowed values for the argument at path, each of
// the argument's type, or nil.
func vendorEnum(schema json.RawMessage, path, typ string) []any {
	leaf := vendorLeaf(schema, path)
	if leaf == nil {
		return nil
	}
	var s struct {
		Enum []json.RawMessage `json:"enum"`
	}
	if json.Unmarshal(leaf, &s) != nil || len(s.Enum) == 0 || len(s.Enum) > maxEnum {
		return nil
	}
	out := make([]any, 0, len(s.Enum))
	for _, raw := range s.Enum {
		var v any
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if dec.Decode(&v) != nil || !hasType(v, typ) {
			return nil
		}
		if str, ok := v.(string); ok && len([]rune(str)) > maxEnumValue {
			return nil
		}
		out = append(out, v)
	}
	return out
}

// outsideEnum says which argument of a call has a value its tool's schema
// doesn't allow, naming the values it does, or nil. The schema is the one
// the agent was given (GrantedTool.InputSchema), whose enums schemaOver
// carried from the vendor's.
func outsideEnum(schema json.RawMessage, args map[string]any) error {
	var s struct {
		Properties map[string]struct {
			Enum []json.RawMessage `json:"enum"`
		} `json:"properties"`
	}
	if json.Unmarshal(schema, &s) != nil {
		return nil
	}
	for _, name := range slices.Sorted(maps.Keys(args)) {
		allowed := s.Properties[name].Enum
		if len(allowed) == 0 {
			continue
		}
		// Compared as canonical JSON, so the same number written two ways
		// (1 and 1.0, 10 and 1e1) is the same value, and so is a string
		// escaped two ways.
		raw, err := json.Marshal(args[name])
		if err != nil {
			return fmt.Errorf("argument %s can't be read", name)
		}
		got, err := jcs.Canonical(raw)
		if err != nil {
			return fmt.Errorf("argument %s can't be read", name)
		}
		if slices.ContainsFunc(allowed, func(a json.RawMessage) bool {
			c, err := jcs.Canonical(a)
			return err == nil && bytes.Equal(c, got)
		}) {
			continue
		}
		words := make([]string, len(allowed))
		for i, a := range allowed {
			var v any
			dec := json.NewDecoder(bytes.NewReader(a))
			dec.UseNumber()
			if dec.Decode(&v) == nil {
				words[i] = fmt.Sprint(v)
			} else {
				words[i] = string(a)
			}
		}
		return fmt.Errorf("argument %s must be one of %s", name, strings.Join(words, ", "))
	}
	return nil
}

// shape is the part of JSON Schema an array argument's reviewed schema may
// use: a type, an object's named fields and which are required, an array's
// items, a string's allowed values and length.
type shape struct {
	Type        string                     `json:"type"`
	Properties  map[string]json.RawMessage `json:"properties,omitempty"`
	Required    []string                   `json:"required,omitempty"`
	Items       json.RawMessage            `json:"items,omitempty"`
	Enum        []string                   `json:"enum,omitempty"`
	MaxLength   *int                       `json:"maxLength,omitempty"`
	MaxItems    *int                       `json:"maxItems,omitempty"`
	Description string                     `json:"description,omitempty"`
	// AdditionalProperties is always false, written or not: a field the
	// schema doesn't name is refused.
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}

const maxShapeDepth = 4

// itemTypes are what an array argument's items may be: objects naming their
// fields, or plain values (a draft's cc addresses), which have no fields to
// hide and are checked as a plain argument is.
var itemTypes = []string{"object", "string", "integer", "number", "boolean"}

// annotations are JSON Schema keywords that describe a value without
// restricting what planque lets through: a provider's generated schema
// carries them (an email draft's cc has a title, a default and examples), and
// an agent copies them into an argument's schema. format is one: planque has
// never checked it (fits doesn't), and the vendor still checks its own.
var annotations = []string{"title", "default", "examples", "example", "$comment", "deprecated", "readOnly", "writeOnly", "nullable", "format"}

// normaliseArguments puts each array argument's schema in the form planque
// checks and keeps (normaliseSchema), in place: what is stored, declared in
// KIFF, shown to the owner and signed is the normalised schema.
func normaliseArguments(t ToolRequest) {
	for name, a := range t.Arguments {
		if a.Type == "array" && len(a.Schema) > 0 {
			a.Schema = normaliseSchema(a.Schema)
			t.Arguments[name] = a
		}
	}
}

// normaliseSchema drops a schema's annotations, at every level, and reads an
// optional value written as anyOf (or oneOf) of it and {"type": "null"} as
// the value alone: the argument is optional already when not required, and
// planque refuses a null. Nothing that restricts a value is dropped; any other
// anyOf, oneOf or keyword stays, for checkShape to refuse. The schema is
// rewritten only when something was dropped, else returned as it came.
func normaliseSchema(raw json.RawMessage) json.RawMessage {
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&m) != nil || m == nil || !normaliseNode(m, 0) {
		return raw
	}
	b, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return b
}

func normaliseNode(m map[string]any, depth int) bool {
	if depth > 2*maxShapeDepth {
		return false // refused as too deep anyway
	}
	changed := false
	for _, k := range []string{"anyOf", "oneOf"} {
		if x := optionalOf(m[k]); x != nil && !overlaps(m, x, k) {
			delete(m, k)
			for kx, v := range x {
				if _, ok := m[kx]; !ok || kx != "description" {
					m[kx] = v // the outer description, when both have one
				}
			}
			changed = true
		}
	}
	for _, k := range annotations {
		if _, ok := m[k]; ok {
			delete(m, k)
			changed = true
		}
	}
	if props, ok := m["properties"].(map[string]any); ok {
		for _, p := range props {
			if pm, ok := p.(map[string]any); ok && normaliseNode(pm, depth+1) {
				changed = true
			}
		}
	}
	if items, ok := m["items"].(map[string]any); ok && normaliseNode(items, depth+1) {
		changed = true
	}
	return changed
}

// optionalOf is X when alts is exactly [X, {"type": "null"}] (in either
// order, annotations aside), else nil.
func optionalOf(alts any) map[string]any {
	list, ok := alts.([]any)
	if !ok || len(list) != 2 {
		return nil
	}
	for i, alt := range list {
		if isNull(alt) {
			x, _ := list[1-i].(map[string]any)
			if isNull(x) {
				return nil
			}
			return x
		}
	}
	return nil
}

// isNull reports whether v is the schema {"type": "null"}, annotations aside.
func isNull(v any) bool {
	m, ok := v.(map[string]any)
	if !ok || m["type"] != "null" {
		return false
	}
	for k := range m {
		if k != "type" && !slices.Contains(annotations, k) {
			return false
		}
	}
	return true
}

// overlaps reports whether x sets a keyword m sets too, its words
// (description, annotations) and the anyOf being read aside: then neither
// can be taken for the other, and the anyOf stays, to be refused.
func overlaps(m, x map[string]any, of string) bool {
	for k := range x {
		if _, ok := m[k]; ok && k != of && k != "description" && !slices.Contains(annotations, k) {
			return true
		}
	}
	return false
}

// checkShape says what is wrong with a reviewed schema, or "": only the
// keywords shape has, an array at the top (top) with object items, every
// object naming its fields, and at most maxShapeDepth levels.
func checkShape(raw json.RawMessage, top bool) string {
	return checkShapeAt(raw, top, 0)
}

func checkShapeAt(raw json.RawMessage, top bool, depth int) string {
	var s shape
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if len(raw) == 0 || dec.Decode(&s) != nil {
		return "it must be a JSON Schema of type, properties, required, items, enum, maxLength, maxItems and description only"
	}
	if depth > maxShapeDepth {
		return "it nests too deep"
	}
	if s.AdditionalProperties != nil && *s.AdditionalProperties {
		return "an object takes only the fields it names"
	}
	if top && s.Type != "array" {
		return "an array argument's schema is an array"
	}
	switch s.Type {
	case "array":
		if top {
			var it shape
			if json.Unmarshal(s.Items, &it) != nil || !slices.Contains(itemTypes, it.Type) {
				return "its items are objects or plain values (string, integer, number or boolean)"
			}
		}
		return checkShapeAt(s.Items, false, depth+1)
	case "object":
		if len(s.Properties) == 0 {
			return "an object names its fields"
		}
		for _, r := range s.Required {
			if _, ok := s.Properties[r]; !ok {
				return "required field " + r + " isn't one of its fields"
			}
		}
		for _, name := range slices.Sorted(maps.Keys(s.Properties)) {
			if why := checkShapeAt(s.Properties[name], false, depth+1); why != "" {
				return name + ": " + why
			}
		}
	case "string", "integer", "number", "boolean":
	default:
		return "type " + s.Type + " isn't one an argument can have"
	}
	return ""
}

// fits says how v doesn't fit a reviewed schema, naming where (at), or ""
// when it does. A field the schema doesn't name doesn't fit.
func fits(v any, raw json.RawMessage, at string) string {
	var s shape
	_ = json.Unmarshal(raw, &s)
	switch s.Type {
	case "array":
		list, ok := v.([]any)
		if !ok {
			return at + " must be an array"
		}
		if s.MaxItems != nil && len(list) > *s.MaxItems {
			return fmt.Sprintf("%s takes at most %d items", at, *s.MaxItems)
		}
		for i, e := range list {
			if why := fits(e, s.Items, fmt.Sprintf("%s[%d]", at, i)); why != "" {
				return why
			}
		}
	case "object":
		m, ok := v.(map[string]any)
		if !ok {
			return at + " must be an object"
		}
		for _, k := range slices.Sorted(maps.Keys(m)) {
			p, named := s.Properties[k]
			if !named {
				return at + "." + k + " is not part of this grant"
			}
			if why := fits(m[k], p, at+"."+k); why != "" {
				return why
			}
		}
		for _, r := range s.Required {
			if _, ok := m[r]; !ok {
				return at + "." + r + " is required"
			}
		}
	case "string":
		str, ok := v.(string)
		if !ok {
			return at + " must be a string"
		}
		if len(s.Enum) > 0 && !slices.Contains(s.Enum, str) {
			return at + " must be one of " + strings.Join(s.Enum, ", ")
		}
		if s.MaxLength != nil && len([]rune(str)) > *s.MaxLength {
			return fmt.Sprintf("%s takes at most %d characters", at, *s.MaxLength)
		}
	default:
		if !hasType(v, s.Type) {
			return at + " must be a " + s.Type
		}
	}
	return ""
}

// pathExpr says whether p is a dotted path of plain field names.
func pathExpr(p string) bool {
	parts, ok := pathParts(p)
	if !ok {
		return false
	}
	for _, k := range parts {
		if !toolName.MatchString(k) {
			return false
		}
	}
	return true
}

// checkOneOf says what is wrong with the tools' one_of reads, or nil: each
// names a read-only tool of the same request, by its vendor tool.
func checkOneOf(tools []ToolRequest) *toolRefusal {
	for _, t := range tools {
		for _, name := range slices.Sorted(maps.Keys(t.Arguments)) {
			o := t.Arguments[name].OneOf
			if o == nil {
				continue
			}
			i := slices.IndexFunc(tools, func(r ToolRequest) bool { return r.Tool == o.Tool })
			if i < 0 || !tools[i].ReadOnly {
				return &toolRefusal{Code: refOneOf, Tool: refTool(t), Field: name, Want: "a read-only tool of the same request",
					Msg: t.Tool + "'s argument " + name + ": one_of must name a read-only tool of the same request"}
			}
		}
	}
	return nil
}

// valuesAt is every string at a dotted path in v, going through lists.
func valuesAt(v any, parts []string) []string {
	if len(parts) == 0 {
		switch x := v.(type) {
		case string:
			return []string{x}
		case []any:
			var out []string
			for _, e := range x {
				out = append(out, valuesAt(e, nil)...)
			}
			return out
		}
		return nil
	}
	switch x := v.(type) {
	case map[string]any:
		return valuesAt(x[parts[0]], parts[1:])
	case []any:
		var out []string
		for _, e := range x {
			out = append(out, valuesAt(e, parts)...)
		}
		return out
	}
	return nil
}

// sameValue compares an argument with an allowed value as people write
// them: spaces dropped, letters in any case ("FR76 3000…" is "fr763000…").
func sameValue(a, b string) bool {
	norm := func(s string) string { return strings.ToUpper(strings.Join(strings.Fields(s), "")) }
	return norm(a) != "" && norm(a) == norm(b)
}

// oneOfWithoutKey says why tools can't carry a one_of on a grant of this
// kind, or "": the check runs only where planque calls the vendor's API with
// the grant's key, so anywhere else it would be signed and never run.
func oneOfWithoutKey(kind string, tools []ToolRequest) string {
	if kind == KindAPIKey {
		return ""
	}
	for _, t := range tools {
		for _, name := range slices.Sorted(maps.Keys(t.Arguments)) {
			if t.Arguments[name].OneOf != nil {
				return t.Tool + "'s argument " + name + ": one_of is checked only on an API-key grant, so it can't be set here"
			}
		}
	}
	return ""
}
