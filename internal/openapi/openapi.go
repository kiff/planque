// Package openapi reads a vendor's published OpenAPI description (3.0 or
// 3.1, JSON) into the tools an API-key grant offers, and builds the HTTP
// request for each call (kiff/planque#64).
//
// What the owner approves is a snapshot: the description is read once, its
// internal references resolved, and every call uses what was read. The
// subset is small on purpose. Anything outside it is refused with the
// reason, before the owner is asked: external references, servers set per
// path or operation, an operation secured otherwise than the grant's key,
// callbacks and links, and a key sent any other way than in a header or the
// query.
package openapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode"

	"github.com/kiff/planque/internal/jcs"
)

// Limits on a description, as the owner's approval and planque's memory
// allow.
const (
	MaxBytes      = 16 << 20
	MaxOperations = 500
	// maxNodes bounds a description once its references are resolved: a
	// small document can reference itself into a huge one.
	maxNodes = 400_000
)

// Auth is where a grant's key goes on each call: a header (Authorization
// with a "Bearer " prefix, X-API-Key…), a query parameter, or HTTP Basic.
type Auth struct {
	In     string `json:"in"` // "header", "query" or "basic"
	Name   string `json:"name"`
	Prefix string `json:"prefix,omitempty"`
	// HTTP Basic: the key is the password, with User as the username
	// (an email, "me@example.com/token"), or, with no User, the username,
	// with Password as the password ("X", or none). Neither is a secret:
	// the owner reads them in the approval.
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	// Fixed is one more header every call carries, as the vendor asks
	// beside the key (Harvest-Account-Id): not a secret, the owner reads it.
	Fixed *Fixed `json:"fixed,omitempty"`
}

// Fixed is a header and its value, the same on every call.
type Fixed struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// headerToken is an HTTP header name.
var headerToken = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+.^_|~-]{1,64}$`)

// notFixed are the headers a fixed one can never be: the request's own, and
// ones that would move or carry a credential.
var notFixed = map[string]bool{"authorization": true, "proxy-authorization": true, "host": true, "cookie": true, "content-type": true,
	"content-length": true, "transfer-encoding": true, "connection": true, "upgrade": true, "te": true, "trailer": true, "expect": true,
	"accept": true, "accept-encoding": true, "forwarded": true, "via": true, "keep-alive": true,
	// A method override turns a call the Card reads as a GET into a DELETE or a POST, at a vendor that honours it.
	"x-http-method-override": true, "x-http-method": true, "x-method-override": true}

func (f *Fixed) check(a Auth) error {
	name := strings.ToLower(f.Name)
	if !headerToken.MatchString(f.Name) || notFixed[name] || strings.HasPrefix(name, "x-forwarded-") || strings.HasPrefix(name, "proxy-") ||
		a.In == "header" && strings.EqualFold(f.Name, a.Name) {
		return refuse("%q can't be a fixed header", f.Name)
	}
	if f.Value == "" || len(f.Value) > MaxBasicPart || strings.ContainsFunc(f.Value, unicode.IsControl) {
		return refuse("a fixed header's value is one line of at most %d characters", MaxBasicPart)
	}
	return nil
}

// MaxBasicPart is the most a Basic username or password given beside the
// key may be.
const MaxBasicPart = 200

// Check is why a is no place for a key, or nil.
func (a Auth) Check() error {
	switch a.In {
	case "header", "query":
		if strings.TrimSpace(a.Name) == "" {
			return refuse("the key's header or query parameter needs a name")
		}
		if a.User != "" || a.Password != "" {
			return refuse("a username or password goes with HTTP Basic only")
		}
	case "basic":
		if a.Name != "" || a.Prefix != "" {
			return refuse("an HTTP Basic key has no header name or prefix")
		}
		for _, part := range []string{a.User, a.Password} {
			if len(part) > MaxBasicPart || strings.ContainsFunc(part, unicode.IsControl) {
				return refuse("a Basic username or password is one line of at most %d characters", MaxBasicPart)
			}
		}
		if strings.Contains(a.User, ":") {
			return refuse("a Basic username can't hold a colon")
		}
		if a.User != "" && a.Password != "" {
			return refuse("with a Basic username given, the key is the password: give no password")
		}
	default:
		return refuse(`the key goes in a header, the query or HTTP Basic ("in": "header", "query" or "basic")`)
	}
	if a.Fixed != nil {
		return a.Fixed.check(a)
	}
	return nil
}

// Param is a path or query parameter of an operation.
type Param struct {
	Name     string `json:"name"`
	In       string `json:"in"` // "path" or "query"
	Required bool   `json:"required,omitempty"`
}

// Operation is one tool an API-key grant offers: the vendor's operation and
// how a call's arguments fill its request.
type Operation struct {
	Name        string  `json:"name"`
	Method      string  `json:"method"`
	Path        string  `json:"path"` // the template, "/charges/{id}"
	Summary     string  `json:"summary,omitempty"`
	Description string  `json:"description,omitempty"`
	ReadOnly    bool    `json:"read_only,omitempty"`
	Params      []Param `json:"params,omitempty"`
	// Body is how the body is filled: "" for none, "fields" when the JSON
	// body is an object whose properties are the tool's other arguments,
	// "body" when the whole JSON body is the tool's "body" argument, "form"
	// when the body is a form of the tool's other arguments, each a string,
	// a number or a list of them.
	Body        string          `json:"body,omitempty"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// Skipped is an operation left out of the tools, and why.
type Skipped struct {
	Name string `json:"name"`
	Why  string `json:"why"`
}

// Snapshot is a description as read: what the owner approves, and what
// every call uses.
type Snapshot struct {
	Version string `json:"version"`
	// Digest is the SHA-256 of the description in planque's normal form: its
	// JSON with keys in order and numbers as written, so a description and
	// its YAML twin, or the same JSON indented otherwise, share it.
	Digest  string   `json:"digest"`
	Title   string   `json:"title,omitempty"`
	Servers []string `json:"servers,omitempty"`
	// Accounts are servers at an account's own address, a variable with no
	// default ("https://{subdomain}.freshdesk.com/api/v2"): the request
	// names the account's, as its base, and the owner reads it.
	Accounts   []string    `json:"accounts,omitempty"`
	Operations []Operation `json:"operations"`
	Skipped    []Skipped   `json:"skipped,omitempty"`
}

// Refused is a description outside what planque reads, with the reason in
// words.
type Refused struct{ Why string }

func (e *Refused) Error() string { return "the API description can't be used: " + e.Why }

func refuse(format string, a ...any) error { return &Refused{Why: fmt.Sprintf(format, a...)} }

var methods = []string{"get", "put", "post", "delete", "patch", "head"}

// ToolName is the rule planque's tool names follow.
var ToolName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,63}$`)

// Parse reads raw, a description, for a grant whose key goes as auth says.
func Parse(raw []byte, auth Auth) (*Snapshot, error) {
	return parseSelected(raw, auth, nil, MaxOperations)
}

// MaxListed is the most operations a description is read for when it's only
// listed (ParseToList): GitHub's has about 1,300. A snapshot is still bound
// to MaxOperations.
const MaxListed = 5000

// ParseToList reads raw as Parse does, past MaxOperations: for listing a big
// API's operations to pick from, never for a grant's snapshot.
func ParseToList(raw []byte, auth Auth) (*Snapshot, error) {
	return parseSelected(raw, auth, nil, MaxListed)
}

// ParseOperations retains only the named operations while keeping the
// digest of the complete source document. The operation limit applies to
// this selection, not to every operation the vendor publishes.
func ParseOperations(raw []byte, auth Auth, operations []string) (*Snapshot, error) {
	if len(operations) == 0 || len(operations) > MaxOperations {
		return nil, refuse("select between 1 and %d operations", MaxOperations)
	}
	selected := make(map[string]bool, len(operations))
	for _, name := range operations {
		if !ToolName.MatchString(name) || selected[name] {
			return nil, refuse("selected operations need distinct tool names")
		}
		selected[name] = true
	}
	return parseSelected(raw, auth, selected, MaxOperations)
}

func parseSelected(raw []byte, auth Auth, selected map[string]bool, most int) (*Snapshot, error) {
	if len(raw) > MaxBytes {
		return nil, refuse("it is over 16 MB")
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, refuse("it is empty")
	}
	// Many vendors publish their description in YAML only: it's read as the
	// same JSON, then checked, digested and limited as any.
	if trimmed[0] != '{' {
		asJSON, err := yamlAsJSON(trimmed)
		if err != nil {
			return nil, err
		}
		trimmed = asJSON
	}
	var doc map[string]any
	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return nil, refuse("it isn't valid JSON")
	}
	normal, err := json.Marshal(doc)
	if err != nil {
		return nil, refuse("it can't be read as JSON")
	}
	version, _ := doc["openapi"].(string)
	if !strings.HasPrefix(version, "3.0.") && !strings.HasPrefix(version, "3.1.") {
		return nil, refuse("it isn't OpenAPI 3.0 or 3.1")
	}
	r := &resolver{doc: doc}
	if auth.Fixed != nil {
		r.fixed = auth.Fixed.Name
	}
	if err := r.checkRefs(doc); err != nil {
		return nil, err
	}
	scheme, err := matchAuth(doc, auth)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(normal)
	s := &Snapshot{Version: version, Digest: hex.EncodeToString(sum[:])}
	if info, ok := doc["info"].(map[string]any); ok {
		s.Title, _ = info["title"].(string)
	}
	for _, sv := range list(doc["servers"]) {
		u, account, err := serverURL(sv)
		if err != nil {
			return nil, err
		}
		if account {
			s.Accounts = append(s.Accounts, u)
			continue
		}
		s.Servers = append(s.Servers, u)
	}
	paths, _ := doc["paths"].(map[string]any)
	keys := make([]string, 0, len(paths))
	for p := range paths {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	global := list(doc["security"])
	names := map[string]bool{}
	// One operation planque can't call is left out with its reason, never the
	// whole description: the others still read. Only what spoils the whole
	// (its size once resolved) refuses it.
	leftOut := func(err error) (string, error) {
		var rf *Refused
		if !errors.As(err, &rf) || r.nodes > maxNodes {
			return "", err
		}
		return rf.Why, nil
	}
	for _, p := range keys {
		item, err := r.object(paths[p])
		pathWhy := ""
		if err != nil {
			if pathWhy, err = leftOut(err); err != nil {
				return nil, err
			}
		} else if _, ok := item["servers"]; ok {
			pathWhy = "its path sets its own servers"
		}
		var shared []rawParam
		if pathWhy == "" {
			if shared, err = r.params(item["parameters"]); err != nil {
				if pathWhy, err = leftOut(err); err != nil {
					return nil, err
				}
			}
		}
		for _, m := range methods {
			raw, ok := item[m]
			if !ok {
				continue
			}
			o, err := r.object(raw)
			if err != nil {
				why, err := leftOut(err)
				if err != nil {
					return nil, err
				}
				s.Skipped = append(s.Skipped, Skipped{Name: opName(m, p, nil), Why: why})
				continue
			}
			if selected != nil && !selected[opName(m, p, o)] {
				continue
			}
			if pathWhy != "" {
				s.Skipped = append(s.Skipped, Skipped{Name: opName(m, p, o), Why: pathWhy})
				continue
			}
			op, why, err := r.operation(p, m, o, shared, global, scheme)
			if err != nil {
				if why, err = leftOut(err); err != nil {
					return nil, err
				}
				if op.Name == "" {
					op.Name = opName(m, p, o)
				}
			}
			if why != "" {
				s.Skipped = append(s.Skipped, Skipped{Name: op.Name, Why: why})
				continue
			}
			if names[op.Name] {
				s.Skipped = append(s.Skipped, Skipped{Name: op.Name, Why: "another operation already has its name"})
				continue
			}
			names[op.Name] = true
			s.Operations = append(s.Operations, op)
			if len(s.Operations) > most {
				return nil, refuse("it has more than %d operations", most)
			}
		}
	}
	if r.nodes > maxNodes {
		return nil, refuse("it is too large once its references are resolved")
	}
	if len(s.Operations) == 0 {
		if len(s.Skipped) > 0 {
			return nil, refuse("it has no operation planque can call (%s: %s)", s.Skipped[0].Name, s.Skipped[0].Why)
		}
		return nil, refuse("it has no operation planque can call")
	}
	if selected != nil && len(names) != len(selected) {
		var missing []string
		for name := range selected {
			if !names[name] {
				missing = append(missing, name)
			}
		}
		sort.Strings(missing)
		return nil, refuse("selected operations are missing or unsupported: %s", strings.Join(missing, ", "))
	}
	return s, nil
}

// Digest identifies what executing an operation depends on. Vendor wording
// can change without pausing the operation; request shape and constraints
// cannot. Schema property names, including a property named description,
// remain part of the identity.
func (op Operation) Digest() (string, error) {
	var schema any
	if err := json.Unmarshal(op.InputSchema, &schema); err != nil {
		return "", err
	}
	semantic := struct {
		Name, Method, Path, Body string
		ReadOnly                 bool
		Params                   []Param
		Schema                   any
	}{op.Name, op.Method, op.Path, op.Body, op.ReadOnly, op.Params, semanticSchema(schema)}
	b, err := json.Marshal(semantic)
	if err != nil {
		return "", err
	}
	return jcs.SHA256(b)
}

func semanticSchema(v any) any {
	switch v := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, value := range v {
			switch k {
			case "description", "title", "$comment", "examples", "example", "externalDocs":
				continue
			case "properties", "patternProperties", "$defs", "definitions", "dependentSchemas":
				if named, ok := value.(map[string]any); ok {
					children := make(map[string]any, len(named))
					for name, child := range named {
						children[name] = semanticSchema(child)
					}
					out[k] = children
					continue
				}
			case "const", "enum", "default":
				// These values are data, not schemas: their object keys are
				// meaningful even when they happen to be annotation names.
				out[k] = value
				continue
			}
			out[k] = semanticSchema(value)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, child := range v {
			out[i] = semanticSchema(child)
		}
		return out
	default:
		return v
	}
}

// matchAuth is the name of the description's security scheme the grant's
// key is, or why there's none.
func matchAuth(doc map[string]any, a Auth) (string, error) {
	if err := a.Check(); err != nil {
		return "", err
	}
	comps, _ := doc["components"].(map[string]any)
	schemes, _ := comps["securitySchemes"].(map[string]any)
	var offered []string
	names := make([]string, 0, len(schemes))
	for n := range schemes {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		sc, _ := schemes[n].(map[string]any)
		typ, _ := sc["type"].(string)
		switch typ {
		case "apiKey":
			in, _ := sc["in"].(string)
			name, _ := sc["name"].(string)
			offered = append(offered, fmt.Sprintf("%s (%s %s)", n, in, name))
			if in == a.In && (in == "header" && strings.EqualFold(name, a.Name) || in == "query" && name == a.Name) {
				return n, nil
			}
		case "http":
			scheme, _ := sc["scheme"].(string)
			offered = append(offered, fmt.Sprintf("%s (http %s)", n, scheme))
			if strings.EqualFold(scheme, "bearer") && a.In == "header" && strings.EqualFold(a.Name, "Authorization") &&
				strings.EqualFold(strings.TrimSpace(a.Prefix), "Bearer") {
				return n, nil
			}
			if strings.EqualFold(scheme, "basic") && a.In == "basic" {
				return n, nil
			}
		default:
			offered = append(offered, fmt.Sprintf("%s (%s, not read)", n, typ))
		}
	}
	if len(offered) == 0 {
		// A description that names no scheme at all (GitHub's): the key goes
		// as the request says, from the vendor's docs, and the owner reads
		// where in the approval. An operation that names a scheme is still
		// refused, since there's none it could be.
		return "", nil
	}
	return "", refuse("none of its security schemes is a key sent as asked; it offers %s", strings.Join(offered, ", "))
}

// serverURL is a server's URL with its variables' defaults filled in.
// serverURL is a server's URL with its variables at their defaults, or, when
// one has none, the URL as written, at an account's own address.
func serverURL(v any) (u string, account bool, err error) {
	sv, _ := v.(map[string]any)
	u, _ = sv["url"].(string)
	vars, _ := sv["variables"].(map[string]any)
	for name, def := range vars {
		d, _ := def.(map[string]any)
		if val, _ := d["default"].(string); val != "" {
			u = strings.ReplaceAll(u, "{"+name+"}", val)
		}
	}
	if strings.ContainsAny(u, "{}") {
		if len(u) > 300 || strings.ContainsFunc(u, unicode.IsControl) {
			return "", false, refuse("a server URL can't be read")
		}
		return u, true, nil
	}
	return u, false, nil
}

// accountPart is what a server variable with no default may be: one label,
// so an address can't leave the vendor's own domain or path.
const accountPart = `[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?`

// OnAccountServer reports whether base is one of the account servers, each
// variable filled with one label, https only.
func OnAccountServer(accounts []string, base string) bool {
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(base, "https://") {
		return false
	}
	for _, a := range accounts {
		if !vendorDomain(a) {
			continue
		}
		var re strings.Builder
		re.WriteString("^")
		rest := strings.TrimRight(a, "/")
		for {
			i := strings.Index(rest, "{")
			j := strings.Index(rest, "}")
			if i < 0 || j < i {
				re.WriteString(regexp.QuoteMeta(rest))
				break
			}
			re.WriteString(regexp.QuoteMeta(rest[:i]))
			re.WriteString(accountPart)
			rest = rest[j+1:]
		}
		re.WriteString("$")
		if m, err := regexp.MatchString(re.String(), base); err == nil && m {
			return true
		}
	}
	return false
}

// vendorDomain reports whether an account server's host ends in the vendor's
// own domain, written out: after its last variable, two labels or more
// ("{subdomain}.zendesk.com"). A host that is all variables
// ("{domain}.{tld}") would take any domain, so no base is ever on it.
func vendorDomain(account string) bool {
	rest, ok := strings.CutPrefix(account, "https://")
	if !ok {
		return false
	}
	host, _, _ := strings.Cut(rest, "/")
	i := strings.LastIndex(host, "}")
	if i < 0 {
		// The host is written out: a variable is in the path only.
		return host != ""
	}
	return fixedDomain.MatchString(host[i+1:])
}

// fixedDomain is a host's end after its last variable: two labels or more.
var fixedDomain = regexp.MustCompile(`^(\.[A-Za-z0-9-]+){2,}$`)

// resolver resolves internal references against doc.
type resolver struct {
	doc   map[string]any
	nodes int
	// fixed is the header every call carries (Auth.Fixed): an operation
	// that needs it as a parameter has it.
	fixed string
}

// checkRefs refuses any reference outside the description.
func (r *resolver) checkRefs(v any) error {
	switch t := v.(type) {
	case map[string]any:
		if ref, ok := t["$ref"].(string); ok && !strings.HasPrefix(ref, "#/") {
			return refuse("it references another document (%s)", ref)
		}
		for _, c := range t {
			if err := r.checkRefs(c); err != nil {
				return err
			}
		}
	case []any:
		for _, c := range t {
			if err := r.checkRefs(c); err != nil {
				return err
			}
		}
	}
	return nil
}

// pointer is the value an internal reference names.
func (r *resolver) pointer(ref string) (any, error) {
	var cur any = r.doc
	for _, part := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		if p, err := url.PathUnescape(part); err == nil {
			part = p
		}
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, refuse("reference %s names nothing", ref)
		}
		if cur, ok = m[part]; !ok {
			return nil, refuse("reference %s names nothing", ref)
		}
	}
	return cur, nil
}

// object is v with a top-level reference followed, as an object.
func (r *resolver) object(v any) (map[string]any, error) {
	for range 16 {
		m, ok := v.(map[string]any)
		if !ok {
			return nil, refuse("an entry isn't an object")
		}
		ref, isRef := m["$ref"].(string)
		if !isRef {
			return m, nil
		}
		var err error
		if v, err = r.pointer(ref); err != nil {
			return nil, err
		}
	}
	return nil, refuse("a reference chain is too long")
}

// inline is v with every reference resolved, for a schema the agent reads.
// A reference to a schema already being resolved (a recursive one) becomes
// a plain object, so nothing loops.
func (r *resolver) inline(v any, seen []string) (any, error) {
	r.nodes++
	if r.nodes > maxNodes {
		return nil, refuse("it is too large once its references are resolved")
	}
	switch t := v.(type) {
	case map[string]any:
		if ref, ok := t["$ref"].(string); ok {
			if slices.Contains(seen, ref) {
				return map[string]any{"type": "object"}, nil
			}
			target, err := r.pointer(ref)
			if err != nil {
				return nil, err
			}
			return r.inline(target, append(seen, ref))
		}
		out := make(map[string]any, len(t))
		for k, c := range t {
			rc, err := r.inline(c, seen)
			if err != nil {
				return nil, err
			}
			out[k] = rc
		}
		return out, nil
	case []any:
		out := make([]any, len(t))
		for i, c := range t {
			rc, err := r.inline(c, seen)
			if err != nil {
				return nil, err
			}
			out[i] = rc
		}
		return out, nil
	}
	return v, nil
}

type rawParam struct {
	Param
	Schema      any
	Description string
}

// params are a parameters list, resolved.
func (r *resolver) params(v any) ([]rawParam, error) {
	var out []rawParam
	for _, e := range list(v) {
		m, err := r.object(e)
		if err != nil {
			return nil, err
		}
		p := rawParam{}
		p.Name, _ = m["name"].(string)
		p.In, _ = m["in"].(string)
		p.Required, _ = m["required"].(bool)
		p.Description, _ = m["description"].(string)
		if p.In == "path" {
			p.Required = true
		}
		schema := m["schema"]
		if schema == nil {
			schema = map[string]any{"type": "string"}
		}
		if p.Schema, err = r.inline(schema, nil); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

var notWord = regexp.MustCompile(`[^A-Za-z0-9]+`)

// opName is an operation's tool name: its operationId where that is a tool
// name, else its method and path in words ("get_charges_id").
func opName(method, path string, o map[string]any) string {
	if id, _ := o["operationId"].(string); id != "" {
		if n := strings.Trim(notWord.ReplaceAllString(id, "_"), "_"); ToolName.MatchString(n) {
			return n
		}
	}
	n := strings.Trim(notWord.ReplaceAllString(method+"_"+path, "_"), "_")
	if len(n) > 64 {
		n = n[:64]
	}
	if n == "" || !ToolName.MatchString(n) {
		n = "op_" + n
	}
	return strings.ToLower(n)
}

// operation reads one operation. why is set when it's left out of the
// tools; err when the whole description is refused.
func (r *resolver) operation(path, method string, o map[string]any, shared []rawParam, global []any, scheme string) (Operation, string, error) {
	op := Operation{Name: opName(method, path, o), Method: strings.ToUpper(method), Path: path, ReadOnly: method == "get" || method == "head"}
	op.Summary, _ = o["summary"].(string)
	op.Description, _ = o["description"].(string)
	if _, ok := o["servers"]; ok {
		return op, "", refuse("operation %s sets its own servers", op.Name)
	}
	if _, ok := o["callbacks"]; ok {
		return op, "", refuse("operation %s has callbacks", op.Name)
	}
	if responses, ok := o["responses"].(map[string]any); ok {
		for _, rv := range responses {
			if resp, err := r.object(rv); err == nil {
				if _, ok := resp["links"]; ok {
					return op, "", refuse("operation %s has links", op.Name)
				}
			}
		}
	}
	security := global
	if sec, ok := o["security"]; ok {
		security = list(sec)
	}
	if len(security) > 0 && !slices.ContainsFunc(security, func(req any) bool {
		m, _ := req.(map[string]any)
		_, ok := m[scheme]
		return ok
	}) {
		return op, "", refuse("operation %s is secured otherwise than the grant's key", op.Name)
	}
	own, err := r.params(o["parameters"])
	if err != nil {
		return op, "", err
	}
	byKey := map[string]rawParam{}
	var order []string
	for _, p := range append(shared, own...) {
		k := p.In + ":" + p.Name
		if _, ok := byKey[k]; !ok {
			order = append(order, k)
		}
		byKey[k] = p
	}
	props := map[string]any{}
	var required []string
	for _, k := range order {
		p := byKey[k]
		switch p.In {
		case "path", "query":
		case "header":
			if r.fixed != "" && strings.EqualFold(p.Name, r.fixed) {
				// Sent on every call, from the grant, never an argument.
				continue
			}
			if p.Required {
				return op, "it needs a " + p.In + " parameter (" + p.Name + "), which planque doesn't send", nil
			}
			continue
		default:
			if p.Required {
				return op, "it needs a " + p.In + " parameter (" + p.Name + "), which planque doesn't send", nil
			}
			continue
		}
		if !ToolName.MatchString(p.Name) && !validArgName(p.Name) {
			return op, "its parameter " + p.Name + " can't be an argument name", nil
		}
		schema, _ := p.Schema.(map[string]any)
		if schema == nil {
			schema = map[string]any{}
		}
		if p.Description != "" {
			if _, has := schema["description"]; !has {
				schema["description"] = p.Description
			}
		}
		props[p.Name] = schema
		op.Params = append(op.Params, p.Param)
		if p.Required {
			required = append(required, p.Name)
		}
	}
	if rb, ok := o["requestBody"]; ok {
		body, err := r.object(rb)
		if err != nil {
			return op, "", err
		}
		content, _ := body["content"].(map[string]any)
		var schema any
		found, form := false, false
		for ct, media := range content {
			if base := strings.TrimSpace(strings.Split(ct, ";")[0]); base == "application/json" || strings.HasSuffix(base, "+json") {
				m, _ := media.(map[string]any)
				schema, found = m["schema"], true
				break
			}
		}
		if !found {
			// A form-encoded body (Stripe's kind), when that's all it takes.
			if media, ok := content["application/x-www-form-urlencoded"].(map[string]any); ok {
				schema, found, form = media["schema"], true, true
			}
		}
		if !found {
			return op, "its body isn't JSON or a form", nil
		}
		resolved, err := r.inline(schema, nil)
		if err != nil {
			return op, "", err
		}
		bodyRequired, _ := body["required"].(bool)
		bs, _ := resolved.(map[string]any)
		bprops, _ := bs["properties"].(map[string]any)
		clash := false
		for k := range bprops {
			if _, ok := props[k]; ok {
				clash = true
			}
		}
		if typ, _ := bs["type"].(string); form && ((typ != "object" && !(typ == "" && bprops != nil)) || bprops == nil || clash) {
			return op, "its form body isn't fields of their own", nil
		}
		if typ, _ := bs["type"].(string); (typ == "object" || typ == "" && bprops != nil) && bprops != nil && !clash {
			op.Body = "fields"
			if form {
				op.Body = "form"
			}
			for k, v := range bprops {
				props[k] = v
			}
			for _, req := range list(bs["required"]) {
				if s, ok := req.(string); ok && bodyRequired {
					required = append(required, s)
				}
			}
		} else {
			if _, ok := props["body"]; ok {
				return op, "a parameter is named body, which its JSON body would take", nil
			}
			op.Body = "body"
			if resolved == nil {
				resolved = map[string]any{}
			}
			props["body"] = resolved
			if bodyRequired {
				required = append(required, "body")
			}
		}
	}
	schema := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		sort.Strings(required)
		schema["required"] = slices.Compact(required)
	}
	op.InputSchema, err = json.Marshal(schema)
	if err != nil {
		return op, "", refuse("operation %s's arguments can't be written as JSON", op.Name)
	}
	return op, "", nil
}

var argName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.\-\[\]]{0,63}$`)

func validArgName(n string) bool { return argName.MatchString(n) }

func list(v any) []any {
	l, _ := v.([]any)
	return l
}

// ErrArguments is a call whose arguments don't fill its operation.
var ErrArguments = errors.New("the arguments don't fit the operation")

// Request is the HTTP request a call makes: base is the grant's server URL,
// args the call's arguments. The key isn't in it: the relay adds it.
func (op Operation) Request(base string, args map[string]any) (method, target string, body []byte, err error) {
	u, err := url.Parse(strings.TrimRight(base, "/"))
	// Which schemes may be reached is the caller's to check (planque's rest client: https only).
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return "", "", nil, fmt.Errorf("%w: the grant's server isn't a URL", ErrArguments)
	}
	rest := map[string]any{}
	for k, v := range args {
		rest[k] = v
	}
	path := op.Path
	q := u.Query()
	for _, p := range op.Params {
		v, ok := rest[p.Name]
		delete(rest, p.Name)
		if !ok || v == nil {
			if p.Required {
				return "", "", nil, fmt.Errorf("%w: %s is required", ErrArguments, p.Name)
			}
			continue
		}
		switch p.In {
		case "path":
			s, ok := scalar(v)
			if !ok {
				return "", "", nil, fmt.Errorf("%w: %s must be a string or a number", ErrArguments, p.Name)
			}
			// "." and ".." aren't escaped, and a server that normalizes the
			// path would answer another operation than the one approved.
			if s == "" || s == "." || s == ".." {
				return "", "", nil, fmt.Errorf("%w: %s can't be empty, . or ..", ErrArguments, p.Name)
			}
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(s))
		case "query":
			if l, isList := v.([]any); isList {
				for _, e := range l {
					s, ok := scalar(e)
					if !ok {
						return "", "", nil, fmt.Errorf("%w: %s must hold strings or numbers", ErrArguments, p.Name)
					}
					q.Add(p.Name, s)
				}
				continue
			}
			s, ok := scalar(v)
			if !ok {
				return "", "", nil, fmt.Errorf("%w: %s must be a string, a number or a list of them", ErrArguments, p.Name)
			}
			q.Set(p.Name, s)
		}
	}
	if strings.ContainsAny(path, "{}") {
		return "", "", nil, fmt.Errorf("%w: a path parameter is missing", ErrArguments)
	}
	switch op.Body {
	case "fields":
		if len(rest) > 0 {
			body, err = json.Marshal(rest)
		}
	case "form":
		form := url.Values{}
		for _, k := range sortedKeys(rest) {
			if l, isList := rest[k].([]any); isList {
				for _, e := range l {
					s, ok := scalar(e)
					if !ok {
						return "", "", nil, fmt.Errorf("%w: %s must hold strings or numbers", ErrArguments, k)
					}
					form.Add(k, s)
				}
				continue
			}
			if rest[k] == nil {
				continue
			}
			s, ok := scalar(rest[k])
			if !ok {
				return "", "", nil, fmt.Errorf("%w: %s must be a string, a number or a list of them: a form has no nested fields", ErrArguments, k)
			}
			form.Set(k, s)
		}
		if len(form) > 0 {
			body = []byte(form.Encode())
		}
	case "body":
		if b, ok := rest["body"]; ok {
			body, err = json.Marshal(b)
			delete(rest, "body")
		}
		if len(rest) > 0 {
			return "", "", nil, fmt.Errorf("%w: %s isn't one of its arguments", ErrArguments, firstKey(rest))
		}
	default:
		if len(rest) > 0 {
			return "", "", nil, fmt.Errorf("%w: %s isn't one of its arguments", ErrArguments, firstKey(rest))
		}
	}
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: the body can't be written as JSON", ErrArguments)
	}
	// path holds its parameters escaped: the URL keeps them as they are.
	escaped := strings.TrimRight(u.EscapedPath(), "/") + path
	unescaped, err := url.PathUnescape(escaped)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: the path can't be written", ErrArguments)
	}
	u.Path, u.RawPath = unescaped, escaped
	u.RawQuery = q.Encode()
	return op.Method, u.String(), body, nil
}

// ContentType is how op's body is written: a form, or JSON.
func (op Operation) ContentType() string {
	if op.Body == "form" {
		return "application/x-www-form-urlencoded"
	}
	return "application/json"
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func firstKey(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys[0]
}

// scalar is v as it goes in a path or query: a string as it is, a number
// as written.
func scalar(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case json.Number:
		return t.String(), true
	case float64:
		return fmt.Sprint(t), true
	case bool:
		return fmt.Sprint(t), true
	}
	return "", false
}
