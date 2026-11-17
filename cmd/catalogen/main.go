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
	Base   string       `json:"base"`   // the API's fixed HTTPS base
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
					t.Fixed = p.Fixed
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
	pin := map[string]any{"kind": "none"}
	if plan.Pin != "" {
		pin = map[string]any{"kind": "base_path", "value": plan.Pin}
	}
	profile := map[string]any{
		"v": "kiff.profile.v1", "vendor": plan.Vendor, "transport": "rest", "pin": pin,
		"source": "catalog", "catalog_entry": plan.ID, "test_proven": false, "operations": ops,
	}
	specName := m[1] + "-openapi.json"
	entry := map[string]any{
		"vendor": plan.Vendor, "title": plan.Title, "kind": "api_key",
		"api":                map[string]any{"url": urlRaw + specName, "base": plan.Base, "auth": plan.Auth, "operations": selected},
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

// firstPathParam is the first path parameter of o, the thing it acts on.
func firstPathParam(o openapi.Operation) string {
	for _, p := range o.Params {
		if p.In == "path" {
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
