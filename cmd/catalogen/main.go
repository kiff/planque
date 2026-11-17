// Command catalogen writes a reviewed key-catalog entry from a vendor's own
// API description: `catalogen -plan deploy/gen/<app>.json`.
//
// The plan names the vendor, where the key goes, the description file and
// which of its operations the entry offers. Every read the plan lists is
// offered as it is; every write is offered held for the owner's approval
// (effect "writes"), unless the plan pins an argument to an undoable value.
// What the entry says about reads and writes comes from the description,
// never from the plan: a write listed as a read, or the other way round, is
// refused.
//
// The description is copied to -specs and its SHA-256 pinned, so a changed
// description is refused before a grant is created. The entry is merged into
// -catalog under its id; an id that exists is never replaced (a change is a
// new revision). The result is decoded with the catalog's own checks before
// anything is written.
package main

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/planque"
)

// Plan is what a person or a worker writes for one app.
type Plan struct {
	ID     string       `json:"id"`     // "pipedrive-key@1"
	Vendor string       `json:"vendor"` // "Pipedrive"
	Title  string       `json:"title"`  // "Pipedrive: read deals and people, add a note"
	Spec   string       `json:"spec"`   // the vendor's description, a path from the plan's folder
	Base   string       `json:"base"`   // the API's fixed HTTPS base, or an account template (https://{account}.vendor.com/...)
	Bases  []string     `json:"bases"`  // instead of Base: the vendor's regions, a closed list the owner picks one of
	Auth   openapi.Auth `json:"auth"`   // where the customer's key goes
	Prefix string       `json:"prefix"` // tool names are <prefix>_<what it does>
	Pin    string       `json:"pin"`    // the base path planque pins ("/v1"); "" for none
	Reads  []Op         `json:"reads"`
	Writes []Op         `json:"writes"`
}

// Op is one operation the entry offers.
type Op struct {
	Operation string `json:"operation"`          // its name in the description
	Name      string `json:"name"`               // the tool's short name: "list_deals"
	Title     string `json:"title"`              // what it does, in the owner's words
	Noun      string `json:"noun,omitempty"`     // a write's resource noun ("deal")
	Resource  string `json:"resource,omitempty"` // a write's resource argument
	// Effect is what a write does when approved: writes (the default), sends,
	// records_money or moves_money. A write that reaches a customer says sends.
	Effect string `json:"effect,omitempty"`
	// Fixed pins an argument (a dotted path) to a value, for a write that
	// can be made final by it ("status": "draft").
	Fixed map[string]json.RawMessage `json:"fixed,omitempty"`
	// Require names arguments the call must give, for a pinned tool whose
	// description doesn't mark them required.
	Require []string `json:"require,omitempty"`
}

var urlRaw = "https://raw.githubusercontent.com/kiff/planque/main/deploy/specs/"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "catalogen:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("catalogen", flag.ContinueOnError)
	planFile := fs.String("plan", "", "the app's plan (JSON)")
	catalogFile := fs.String("catalog", "deploy/catalog.json", "the reviewed catalog entries")
	appsFile := fs.String("apps", "deploy/apps.json", "the registered apps")
	specs := fs.String("specs", "deploy/specs", "where the descriptions are kept")
	list := fs.Bool("list", false, "only list the description's operations (name, method, path, read or write) to choose from")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *planFile == "" {
		return errors.New("-plan is required")
	}
	var plan Plan
	if err := decodeStrict(*planFile, &plan); err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(*planFile), plan.Spec))
	if err != nil {
		return err
	}
	// The description is kept, and pinned, as JSON, so that every check on
	// shipped descriptions reads it, whatever form the vendor published it in.
	if raw, err = openapi.AsJSON(raw); err != nil {
		return fmt.Errorf("%s: %w", plan.Spec, err)
	}
	// Examples are annotations the pinned copy doesn't need: they carry sample
	// data (keys, addresses, names) that reads as real, and they are most of a
	// big description's size. Everything an operation's behaviour depends on stays.
	if raw, err = withoutExamples(raw); err != nil {
		return fmt.Errorf("%s: %w", plan.Spec, err)
	}
	if *list {
		snap, err := openapi.ParseToList(raw, plan.Auth)
		if err != nil {
			return err
		}
		for _, o := range snap.Operations {
			kind := "write"
			if o.ReadOnly {
				kind = "read"
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", o.Name, o.Method, o.Path, kind, o.Summary)
		}
		return nil
	}
	entry, specName, err := Generate(plan, raw)
	if err != nil {
		return err
	}
	return merge(*catalogFile, *appsFile, *specs, specName, raw, plan.ID, entry)
}

func decodeStrict(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

var (
	idRe     = regexp.MustCompile(`^([a-z][a-z0-9_-]{0,63})@([1-9][0-9]{0,8})$`)
	shortRe  = regexp.MustCompile(`^[a-z][a-z0-9_]{0,40}$`)
	prefixRe = regexp.MustCompile(`^[a-z][a-z0-9]{0,20}$`)
)

// Generate builds the entry for plan from raw, the vendor's description, and
// the name its copy is kept under.
func Generate(plan Plan, raw []byte) (map[string]any, string, error) {
	m := idRe.FindStringSubmatch(plan.ID)
	if m == nil || !prefixRe.MatchString(plan.Prefix) {
		return nil, "", errors.New("a plan needs an id like name@1 and a short lowercase prefix")
	}
	if len(plan.Reads)+len(plan.Writes) == 0 {
		return nil, "", errors.New("a plan offers at least one operation")
	}
	snap, err := openapi.ParseToList(raw, plan.Auth)
	if err != nil {
		return nil, "", err
	}
	byName := map[string]openapi.Operation{}
	for _, o := range snap.Operations {
		byName[o.Name] = o
	}
	var tools []planque.ToolRequest
	var ops []map[string]any
	var selected []string
	seen := map[string]bool{}
	for _, list := range []struct {
		ops   []Op
		write bool
	}{{plan.Reads, false}, {plan.Writes, true}} {
		for _, p := range list.ops {
			o, ok := byName[p.Operation]
			if !ok {
				return nil, "", fmt.Errorf("%s: not an operation of the description", p.Operation)
			}
			if o.ReadOnly == list.write {
				return nil, "", fmt.Errorf("%s is a %s %s in the description, and the plan lists it as a %s", p.Operation, o.Method, o.Path, map[bool]string{false: "read", true: "write"}[list.write])
			}
			if !shortRe.MatchString(p.Name) || p.Title == "" || strings.TrimSpace(p.Title) != p.Title {
				return nil, "", fmt.Errorf("%s: a short lowercase name and a title are needed", p.Operation)
			}
			name := plan.Prefix + "_" + p.Name
			if !openapi.ToolName.MatchString(name) || seen[name] || seen[p.Operation] {
				return nil, "", fmt.Errorf("%s: its name %q is taken or invalid", p.Operation, name)
			}
			seen[name], seen[p.Operation] = true, true
			selected = append(selected, p.Operation)
			t := planque.ToolRequest{Tool: p.Operation, Name: name, Title: p.Title, ReadOnly: !list.write}
			po := map[string]any{"tool": name, "title": p.Title}
			if !list.write {
				po["kind"] = "read"
			} else {
				required := requiredArgs(o)
				// A write that can copy an outside address (a Cc, a Bcc, a list of
				// recipients) does more than its title tells the owner: the plan
				// pins that argument, or the description drops it.
				for _, path := range addressProps(o.InputSchema) {
					if _, pinned := p.Fixed[path]; !pinned {
						return nil, "", fmt.Errorf("%s takes %q, an address list: pin it in the plan's fixed, or take it out of the description", p.Operation, path)
					}
				}
				if p.Resource == "" {
					p.Resource = firstPathParam(o)
				}
				if p.Resource != "" && !slices.Contains(required, p.Resource) {
					return nil, "", fmt.Errorf("%s: resource %q isn't one of its required arguments", p.Operation, p.Resource)
				}
				if p.Noun == "" && p.Resource != "" {
					return nil, "", fmt.Errorf("%s: a write on %q needs a noun for it", p.Operation, p.Resource)
				}
				if len(p.Fixed) > 0 {
					// A tool with values pinned is narrowed: its arguments are listed here,
					// what is pinned left out, and a call can't send anything else.
					args, err := narrowedArguments(o, p.Fixed, p.Require)
					if err != nil {
						return nil, "", fmt.Errorf("%s: %w", p.Operation, err)
					}
					t.Fixed, t.Arguments = p.Fixed, args
				}
				po["kind"] = "write"
				effect := cmp.Or(p.Effect, "writes")
				if !slices.Contains([]string{"writes", "sends", "records_money", "moves_money"}, effect) {
					return nil, "", fmt.Errorf("%s: effect %q isn't writes, sends, records_money or moves_money", p.Operation, effect)
				}
				po["effect"] = effect
				po["shape"] = "single"
				if p.Resource != "" {
					t.ResourceArg = p.Resource
					po["resource"] = map[string]any{"arg": p.Resource, "noun": p.Noun}
				}
				po["required_arguments"] = required
			}
			tools = append(tools, t)
			ops = append(ops, po)
		}
	}
	sum := sha256.Sum256(raw)
	specName := m[1] + "-openapi.json"
	api := map[string]any{"url": urlRaw + specName, "auth": plan.Auth, "operations": selected}
	if len(plan.Bases) > 0 {
		api["bases"] = plan.Bases
	} else {
		api["base"] = plan.Base
	}
	pin := map[string]any{"kind": "none"}
	if plan.Pin != "" {
		pin = map[string]any{"kind": "base_path", "value": plan.Pin}
	}
	profile := map[string]any{
		"v": "kiff.profile.v1", "vendor": plan.Vendor, "transport": "rest", "pin": pin,
		"source": "catalog", "catalog_entry": plan.ID, "test_proven": false, "operations": ops,
	}
	entry := map[string]any{
		"vendor": plan.Vendor, "title": plan.Title, "kind": "api_key",
		"api":                api,
		"description_sha256": hex.EncodeToString(sum[:]),
		"tools":              tools,
		"profile":            profile,
	}
	return entry, specName, nil
}

// requiredArgs are the arguments a call must give: the description's own
// required list, in its order.
func requiredArgs(o openapi.Operation) []string {
	var s struct {
		Required []string `json:"required"`
	}
	_ = json.Unmarshal(o.InputSchema, &s)
	return s.Required
}

// firstPathParam is the first path parameter of o that is a string, the
// thing it acts on: KIFF's Card can scope on a string argument only, so a
// write on a numeric id has no resource.
func firstPathParam(o openapi.Operation) string {
	var schema struct {
		Properties map[string]struct {
			Type any `json:"type"`
		} `json:"properties"`
	}
	_ = json.Unmarshal(o.InputSchema, &schema)
	for _, p := range o.Params {
		if p.In == "path" && schema.Properties[p.Name].Type == "string" {
			return p.Name
		}
	}
	return ""
}

// merge writes the description beside the others and adds the entry to the
// catalog, after the catalog's own decoding accepts the result.
func merge(catalogFile, appsFile, specs, specName string, raw []byte, id string, entry map[string]any) error {
	cb, err := os.ReadFile(catalogFile)
	if err != nil {
		return err
	}
	var cat map[string]json.RawMessage
	if err := json.Unmarshal(cb, &cat); err != nil {
		return err
	}
	if _, ok := cat[id]; ok {
		return fmt.Errorf("%s exists: a change is a new revision", id)
	}
	eb, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	cat[id] = eb
	out, err := json.MarshalIndent(cat, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	apps := map[string]appconfig.Application{}
	if ab, err := os.ReadFile(appsFile); err == nil {
		if apps, err = appconfig.Decode(ab); err != nil {
			return err
		}
	}
	if _, err := planque.DecodeCatalog(out, apps); err != nil {
		return fmt.Errorf("the catalog refuses %s: %w", id, err)
	}
	if err := os.WriteFile(filepath.Join(specs, specName), raw, 0o644); err != nil {
		return err
	}
	return os.WriteFile(catalogFile, out, 0o644)
}

// withoutExamples is a JSON description without its "example" and "examples"
// annotations. A schema property that is itself named so (under "properties")
// is kept: it is data, not an annotation.
func withoutExamples(raw []byte) ([]byte, error) {
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	var strip func(v any, underProperties bool) any
	strip = func(v any, underProperties bool) any {
		switch x := v.(type) {
		case map[string]any:
			for k, c := range x {
				if !underProperties && (k == "example" || k == "examples") {
					delete(x, k)
					continue
				}
				x[k] = strip(c, k == "properties")
			}
			return x
		case []any:
			for i, c := range x {
				x[i] = strip(c, false)
			}
			return x
		}
		return v
	}
	return json.Marshal(strip(doc, false))
}

// addressName is a property that names who else a message goes to.
var addressName = regexp.MustCompile(`(?i)^(cc|bcc|to|recipients?)$|^(cc|bcc|to)[_-]|[_-](cc|bcc|recipients?)$|^(cc|bcc)[_-]?(emails?|addresses)$`)

// addressProps are the dotted paths, in a tool's input schema, of properties
// that name addresses a message is copied to.
func addressProps(schema json.RawMessage) []string {
	var doc any
	if json.Unmarshal(schema, &doc) != nil {
		return nil
	}
	var out []string
	var walk func(v any, path string)
	walk = func(v any, path string) {
		switch x := v.(type) {
		case map[string]any:
			if props, ok := x["properties"].(map[string]any); ok {
				for name, c := range props {
					p := name
					if path != "" {
						p = path + "." + name
					}
					if addressName.MatchString(name) {
						out = append(out, p)
					}
					walk(c, p)
				}
			}
			for _, k := range []string{"items", "oneOf", "anyOf", "allOf"} {
				if c, ok := x[k]; ok {
					walk(c, path)
				}
			}
		case []any:
			for _, c := range x {
				walk(c, path)
			}
		}
	}
	walk(doc, "")
	slices.Sort(out)
	return out
}

// narrowedArguments are the arguments of a tool with values pinned: each
// property of the operation's input but those pinned, a plain scalar.
func narrowedArguments(o openapi.Operation, fixed map[string]json.RawMessage, require []string) (map[string]planque.Argument, error) {
	var schema struct {
		Properties map[string]struct {
			Type any `json:"type"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(o.InputSchema, &schema); err != nil {
		return nil, err
	}
	args := map[string]planque.Argument{}
	for name, prop := range schema.Properties {
		if _, pinned := fixed[name]; pinned {
			continue
		}
		typ, _ := prop.Type.(string)
		if !slices.Contains([]string{"string", "integer", "number", "boolean"}, typ) {
			return nil, fmt.Errorf("argument %q is a %v: a tool with pinned values takes plain arguments only, so take it out of the description", name, prop.Type)
		}
		args[name] = planque.Argument{Path: name, Type: typ, Required: slices.Contains(schema.Required, name) || slices.Contains(require, name)}
	}
	for _, name := range require {
		if _, ok := args[name]; !ok {
			return nil, fmt.Errorf("%q is required but isn't an argument of the tool", name)
		}
	}
	for name := range fixed {
		if strings.Contains(name, ".") {
			continue
		}
		if _, ok := schema.Properties[name]; !ok {
			return nil, fmt.Errorf("%q is pinned but isn't an argument of the operation", name)
		}
	}
	return args, nil
}
