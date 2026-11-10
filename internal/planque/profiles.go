package planque

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/rest"
)

// ProfileDigests is the phone's checked pair, under its existing request
// signature. planque checks the profile binding; KIFF independently checks
// both digests and the signed evidence when it applies the revision.
type ProfileDigests struct {
	ProfileSHA256 string `json:"profile_sha256"`
	WordsSHA256   string `json:"words_sha256"`
}

func profileApproval(raw json.RawMessage, approved *ProfileDigests) error {
	if len(raw) == 0 {
		if approved != nil {
			return apiError{http.StatusConflict, "this grant has no profile to approve"}
		}
		return nil
	}
	if approved == nil {
		return apiError{http.StatusConflict, "approve the profile and KIFF's words on the phone first"}
	}
	digest, err := profile.Digest(raw)
	words, hexErr := hex.DecodeString(approved.WordsSHA256)
	if err != nil || approved.ProfileSHA256 != digest || hexErr != nil || len(words) != 32 || approved.WordsSHA256 != strings.ToLower(approved.WordsSHA256) {
		return apiError{http.StatusConflict, "the signed profile or words digest does not match this grant"}
	}
	return nil
}

// profileFields derives connection measurement settings from the profile
// before the phone sees its grant. A proposal cannot declare itself reviewed
// catalog data or claim successful test calls.
func profileFields(raw json.RawMessage, tools []ToolRequest) (json.RawMessage, error) {
	return deriveProfileFields(raw, tools, false)
}

func deriveProfileFields(raw json.RawMessage, tools []ToolRequest, reviewed bool) (json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	p, err := profile.Parse(raw)
	if err != nil {
		return nil, err
	}
	if !reviewed && (p.Source != "drafted" || p.TestProven) {
		return nil, errors.New("only a reviewed catalog entry or a recorded test can supply catalog or test-proven evidence")
	}
	if len(p.Operations) != len(tools) {
		return nil, errors.New("the profile must name every requested tool exactly once")
	}
	byName := map[string]*ToolRequest{}
	for i := range tools {
		name := tools[i].Name
		if name == "" {
			name = tools[i].Tool
		}
		if _, exists := byName[name]; exists {
			return nil, errors.New("give distinct names to the profile's tools")
		}
		byName[name] = &tools[i]
	}
	for _, op := range p.Operations {
		for _, name := range op.RequiredArguments {
			if name == IdempotencyArg {
				return nil, errors.New("a profile cannot require planque's operation argument")
			}
		}
		t := byName[op.Tool]
		if t == nil {
			return nil, fmt.Errorf("profile operation %s is not one of the requested tools", op.Tool)
		}
		t.ReadOnly = op.Kind == "read"
		t.AmountArg, t.AmountPath, t.AmountFormat, t.AmountUnit, t.AmountScale, t.AmountLookup, t.UnitLookup = "", "", "", "", nil, nil, nil
		t.ResourceArg, t.ResourcePath = "", ""
		if op.Resource != nil {
			t.ResourceArg, t.ResourcePath = op.Resource.Arg, op.Resource.Path
		}
		if a := op.Amount; a != nil {
			t.AmountArg, t.AmountPath, t.AmountFormat, t.AmountUnit = a.Arg, a.Path, a.Format, a.Unit
			scale := a.Scale
			t.AmountScale = &scale
			if a.Read != nil {
				t.AmountArg = "profile_amount"
				t.AmountLookup, err = profileLookup(*a.Read, t.ResourceArg)
				if err != nil {
					return nil, err
				}
			}
			if a.UnitRead != nil {
				t.UnitLookup, err = profileLookup(*a.UnitRead, t.ResourceArg)
				if err != nil {
					return nil, err
				}
				// A dynamic unit still needs the operation's nominal currency
				// for KIFF's connection settings. Never invent one.
				return nil, errors.New("a unit_read profile also needs a supported nominal-unit connection contract")
			}
		}
	}
	return json.Marshal(p)
}

func profileTransport(raw json.RawMessage, api *APIRequest, server string) error {
	return bindProfileTransport(raw, api, server, "")
}

func bindProfileTransport(raw json.RawMessage, api *APIRequest, server, header string) error {
	if len(raw) == 0 {
		return nil
	}
	p, err := profile.Parse(raw)
	if err != nil {
		return err
	}
	want := "mcp"
	if api != nil {
		want = "rest"
	}
	if p.Transport != want {
		return errors.New("the profile transport does not match this connection")
	}
	if p.Pin.Kind == "header" {
		if api == nil || header == "" {
			return errors.New("this vendor needs a reviewed version-header definition before its header pin can be used")
		}
		pin := &rest.VersionPin{Name: header, Value: p.Pin.Value}
		if err := pin.Validate(); err != nil {
			return err
		}
		if api.Auth.In == "header" && strings.EqualFold(api.Auth.Name, header) {
			return errors.New("a version pin cannot replace the vendor credential")
		}
		api.VersionPin = pin
	} else if api != nil && api.VersionPin != nil {
		return errors.New("the profile must preserve this connection's approved version header")
	}
	if p.Pin.Kind == "base_path" {
		u, err := url.Parse(server)
		prefix := strings.TrimRight(p.Pin.Value, "/")
		if err != nil || !strings.HasPrefix(prefix, "/") || u.Path != prefix && !strings.HasPrefix(u.Path, prefix+"/") {
			return errors.New("the connection's base does not use the profile's pinned version path")
		}
	}
	return nil
}

func profileLookup(read profile.Read, resource string) (*kiff.UnitLookup, error) {
	if resource == "" || len(read.Args) != 1 {
		return nil, errors.New("a measurement read must take the call's resource in one argument")
	}
	for arg, value := range read.Args {
		if value != "$args."+resource {
			return nil, errors.New("a measurement read must take the call's resource")
		}
		return &kiff.UnitLookup{Tool: read.Tool, Argument: arg, Path: read.Path}, nil
	}
	return nil, errors.New("measurement read has no resource argument")
}

func toolNames(tools []ToolRequest) []string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
		if names[i] == "" {
			names[i] = t.Tool
		}
	}
	return names
}

// renameProfile follows names allocated by planque, before approval. Read
// references move with their tool; caller-provided vendor names never select
// another account's connection.
func renameProfile(raw json.RawMessage, before []string, tools []ToolRequest) (json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	p, err := profile.Parse(raw)
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	for i, name := range before {
		names[name] = tools[i].Name
	}
	for i := range p.Operations {
		o := &p.Operations[i]
		if name, ok := names[o.Tool]; ok {
			o.Tool = name
		}
		if o.Amount != nil {
			for _, read := range []*profile.Read{o.Amount.Read, o.Amount.UnitRead} {
				if read != nil {
					if name, ok := names[read.Tool]; ok {
						read.Tool = name
					}
				}
			}
		}
		if o.VerifyRead != nil {
			if name, ok := names[o.VerifyRead.Tool]; ok {
				o.VerifyRead.Tool = name
			}
		}
		if o.StatusRead != nil {
			if name, ok := names[o.StatusRead.Tool]; ok {
				o.StatusRead.Tool = name
			}
		}
		if c := o.Case; c != nil {
			// KIFF makes the case reads as the agent, by its own names.
			for _, tool := range []*string{&c.AccountRead.Tool, &c.ChargeRead.Tool, &c.RefundsRead.Tool} {
				if name, ok := names[*tool]; ok {
					*tool = name
				}
			}
		}
	}
	return json.Marshal(p)
}

func profileSchema(g *Grant, t GrantedTool) error {
	if len(g.Request.Profile) == 0 {
		return nil
	}
	path := t.AmountArg
	if t.AmountPath != "" {
		path = t.AmountPath
	}
	if path != "" && t.AmountLookup == nil {
		want := "integer"
		if t.AmountFormat == "decimal_string" {
			want = "string"
		}
		if schemaType(t.InputSchema, path) != want {
			return fmt.Errorf("%s's amount must be a %s in its schema", t.Tool, want)
		}
	}
	resource := t.ResourceArg
	if t.ResourcePath != "" {
		resource = t.ResourcePath
	}
	if resource != "" && schemaType(t.InputSchema, resource) != "string" {
		return fmt.Errorf("%s's resource must be a string in its schema", t.Tool)
	}
	return nil
}

func schemaType(raw json.RawMessage, path string) string {
	for _, name := range strings.Split(path, ".") {
		var s struct {
			Properties map[string]json.RawMessage `json:"properties"`
			Additional json.RawMessage            `json:"additionalProperties"`
		}
		if json.Unmarshal(raw, &s) != nil {
			return ""
		}
		raw = s.Properties[name]
		if len(raw) == 0 && len(s.Additional) > 0 && s.Additional[0] == '{' {
			raw = s.Additional
		}
	}
	var leaf struct {
		Type string `json:"type"`
	}
	if json.Unmarshal(raw, &leaf) != nil {
		return ""
	}
	return leaf.Type
}

func profileOperation(raw json.RawMessage, name string) (profile.Operation, bool) {
	if len(raw) == 0 {
		return profile.Operation{}, false
	}
	p, err := profile.Parse(raw)
	if err != nil {
		return profile.Operation{}, false
	}
	for _, op := range p.Operations {
		if op.Tool == name {
			return op, true
		}
	}
	return profile.Operation{}, false
}

func profileVendorSchema(g *Grant, name string, schema json.RawMessage) error {
	op, ok := profileOperation(g.Request.Profile, name)
	if !ok || op.Marker == nil {
		return nil
	}
	if strings.ContainsAny(op.Marker.Field, "[]") || schemaType(schema, op.Marker.Field) != "string" {
		return errors.New("the marker must name a string field in the vendor's input schema")
	}
	return nil
}

func profileArguments(g *Grant, name string, args map[string]any) error {
	op, ok := profileOperation(activeProfile(g), name)
	if !ok {
		return nil
	}
	for _, required := range op.RequiredArguments {
		if value, exists := args[required]; !exists || value == nil {
			return fmt.Errorf("the profile requires argument %s", required)
		}
	}
	return nil
}

func profileMarker(g *Grant, name, operation string, raw json.RawMessage) (json.RawMessage, error) {
	op, ok := profileOperation(activeProfile(g), name)
	if !ok || op.Marker == nil {
		return raw, nil
	}
	var args map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if dec.Decode(&args) != nil || args == nil {
		return nil, errors.New("the vendor arguments cannot hold the profile marker")
	}
	parts := strings.Split(op.Marker.Field, ".")
	parent := args
	for _, part := range parts[:len(parts)-1] {
		value, exists := parent[part]
		if !exists {
			value = map[string]any{}
			parent[part] = value
		}
		next, ok := value.(map[string]any)
		if !ok {
			return nil, errors.New("the profile marker overlaps a non-object vendor argument")
		}
		parent = next
	}
	field := parts[len(parts)-1]
	if _, supplied := parent[field]; supplied {
		return nil, errors.New("the profile marker is set by planque, not a caller")
	}
	parent[field] = operation
	return json.Marshal(args)
}

func activeProfile(g *Grant) json.RawMessage {
	if g.Connected != nil {
		return g.ConnectedProfile
	}
	return g.Request.Profile
}

func profileReady(g *Grant) bool {
	raw := activeProfile(g)
	if len(raw) == 0 {
		return true
	}
	digest, err := profile.Digest(raw)
	return err == nil && g.ProfileRevision != nil && g.ProfileRevision.ProfileSHA256 == digest
}

func (s *Server) applyProfile(ctx context.Context, key string, g *Grant) error {
	if len(g.Request.Profile) == 0 {
		return nil
	}
	if g.Approval == nil {
		return errors.New("the profile has no signed approval")
	}
	var signed struct {
		Profile *ProfileDigests `json:"kiff_profile"`
	}
	if json.Unmarshal([]byte(g.Approval.Body), &signed) != nil {
		return errors.New("the profile approval is unreadable")
	}
	if err := profileApproval(g.Request.Profile, signed.Profile); err != nil {
		return err
	}
	a := g.Approval
	revision, err := s.KIFF.ApplyProfile(ctx, key, g.Request.Profile, kiff.ProfileApproval{
		Message: base64.StdEncoding.EncodeToString([]byte(a.Message)), Body: base64.StdEncoding.EncodeToString([]byte(a.Body)), Signature: a.Signature, DeviceID: a.DeviceID})
	if err != nil {
		return err
	}
	if revision.ProfileSHA256 != signed.Profile.ProfileSHA256 || revision.WordsSHA256 != signed.Profile.WordsSHA256 {
		return errors.New("KIFF confirmed another profile or words digest")
	}
	g.ProfileRevision = &revision
	return nil
}

// withApprovedOperations is the profile for a grant that already has
// approved, plus the tools being added: when next names only the added
// tools' operations, approved with next's operations appended, so every
// approved operation and the approved profile's own fields stay byte for
// byte as the owner approved them (unchangedProfile then holds). Anything
// else (no approved profile, next empty or unreadable, or naming a tool the
// grant already has) is next as given, checked as before.
func withApprovedOperations(approved, next json.RawMessage, added []ToolRequest) json.RawMessage {
	if len(approved) == 0 || len(next) == 0 {
		return next
	}
	var base, more map[string]json.RawMessage
	if json.Unmarshal(approved, &base) != nil || json.Unmarshal(next, &more) != nil {
		return next
	}
	var have, adds []json.RawMessage
	if json.Unmarshal(base["operations"], &have) != nil || json.Unmarshal(more["operations"], &adds) != nil || len(adds) != len(added) {
		return next
	}
	names := map[string]bool{}
	for _, t := range added {
		names[t.Name] = true
		names[t.Tool] = true
	}
	for _, op := range adds {
		var o struct {
			Tool string `json:"tool"`
		}
		if json.Unmarshal(op, &o) != nil || o.Tool == "" || !names[o.Tool] {
			return next
		}
	}
	ops, err := json.Marshal(append(have, adds...))
	if err != nil {
		return next
	}
	base["operations"] = ops
	merged, err := json.Marshal(base)
	if err != nil {
		return next
	}
	return merged
}

// Adding tools never silently changes measurements or behavior of an
// existing operation. A changed operation needs a separately approved repair.
func unchangedProfile(old, next json.RawMessage) bool {
	if len(old) == 0 {
		return len(next) == 0
	}
	a, err := profile.Parse(old)
	if err != nil {
		return false
	}
	b, err := profile.Parse(next)
	if err != nil {
		return false
	}
	byName := map[string]profile.Operation{}
	for _, op := range b.Operations {
		byName[op.Tool] = op
	}
	for _, op := range a.Operations {
		if !reflect.DeepEqual(op, byName[op.Tool]) {
			return false
		}
	}
	a.Operations, b.Operations = nil, nil
	return reflect.DeepEqual(a, b)
}

func sameToolSettings(a, b ToolRequest) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}
