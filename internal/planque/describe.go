package planque

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/mcp"
)

// A vendor's broad tool reaches any operation (Stripe's stripe_api_write
// writes to any endpoint), and a grant narrows it to one by pinning the
// operation. The fields the taf may set are then the operation's own, as the
// vendor's describing tool lists them: an argument at a path the operation
// doesn't take is refused by the vendor at the first call, after the owner
// approved it and wrote a Card on it. planque checks each path against the
// vendor's list before the owner is asked.

// describer is a vendor's describing tool and the broad tools it describes.
type describer struct {
	Tool   string   // the vendor's tool that lists an operation's fields
	OpArg  string   // the argument that pins the operation, on both tools
	Params string   // where the operation's fields go in a broad tool's arguments
	Writes []string // broad tools that write
	Reads  []string // broad tools that only read
	// Context are the arguments that say which account and mode to act
	// in: the describing tool reads the operation in the narrowed tool's
	// account and mode, so they go with the read as the narrowed tool pins
	// them. Stripe refuses stripe_api_details without stripe_context or
	// without livemode.
	Context []string
}

var describers = []describer{{Tool: "stripe_api_details", OpArg: "stripe_api_operation_id", Params: narrowedParams,
	Writes: []string{"stripe_api_write"}, Reads: []string{"stripe_api_read"}, Context: []string{"stripe_context", "livemode"}}}

// narrowedTo is the describer of t and the operation t pins, when t is a
// broad tool narrowed to one operation; write says t writes.
func narrowedTo(t ToolRequest) (d describer, op string, write, ok bool) {
	for _, d := range describers {
		write := slices.Contains(d.Writes, t.Tool)
		if !write && !slices.Contains(d.Reads, t.Tool) {
			continue
		}
		if json.Unmarshal(t.Fixed[d.OpArg], &op) != nil || op == "" {
			return describer{}, "", false, false
		}
		return d, op, write, true
	}
	return describer{}, "", false, false
}

// field is one of an operation's fields, with the fields inside it. open
// is an object or array whose inner fields the vendor doesn't list.
type field struct {
	kids map[string]*field
	open bool
}

var fieldLine = regexp.MustCompile(`^( *)- \*\*([^*\s]+)\*\* \[([a-z]+)\]`)

// parseFields reads an operation's fields from its description, as Stripe's
// stripe_api_details writes it: each "## … Parameters" section lists its
// fields as "- **name** [type] …", a field inside another two spaces deeper.
func parseFields(text string) map[string]*field {
	top := map[string]*field{}
	var stack []*field // the field at each depth
	in := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			in = strings.HasPrefix(line, "## ") && strings.HasSuffix(strings.TrimSpace(line), "Parameters")
			stack = stack[:0]
			continue
		}
		m := fieldLine.FindStringSubmatch(line)
		if !in || m == nil {
			continue
		}
		depth := len(m[1]) / 2
		if depth > len(stack) {
			continue // deeper than any field above it: not one of its fields
		}
		stack = stack[:depth]
		f := &field{kids: map[string]*field{}, open: m[3] == "object" || m[3] == "array"}
		parent := top
		if depth > 0 {
			parent = stack[depth-1].kids
		}
		parent[m[2]] = f
		stack = append(stack, f)
	}
	return top
}

// fieldsCheck says which of t's paths into d's fields op doesn't take, with
// the fields it does take, or "".
// marker is the path planque writes the profile's marker at, or "".
func fieldsCheck(t ToolRequest, d describer, op string, fields map[string]*field, marker string) string {
	var paths []string
	for _, name := range slices.Sorted(maps.Keys(t.Arguments)) {
		paths = append(paths, t.Arguments[name].Path)
	}
	paths = append(paths, slices.Sorted(maps.Keys(t.Fixed))...)
	if marker != "" {
		paths = append(paths, marker)
	}
	for _, p := range paths {
		if p == d.Params {
			return fmt.Sprintf("%s: %s takes one field of %s, never all of %s", t.Tool, p, op, d.Params)
		}
		rest, ok := strings.CutPrefix(p, d.Params+".")
		if !ok {
			continue
		}
		level, parts := fields, strings.Split(rest, ".")
		for i, part := range parts {
			f := level[part]
			if f == nil {
				return fmt.Sprintf("%s takes no field %s: %s lists %s", op, rest, d.Tool, listFields(fields))
			}
			if len(f.kids) == 0 {
				// A field inside one whose inner fields aren't listed is the
				// vendor's to check; inside a plain value there is none.
				if !f.open && i < len(parts)-1 {
					return fmt.Sprintf("%s takes no field %s: %s lists %s", op, rest, d.Tool, listFields(fields))
				}
				break
			}
			level = f.kids
		}
	}
	return ""
}

func listFields(fields map[string]*field) string {
	if len(fields) == 0 {
		return "no fields for it"
	}
	return strings.Join(slices.Sorted(maps.Keys(fields)), ", ")
}

// operationsCheck checks each of tools narrowed to one operation against
// the operation's fields, read with the vendor's describing tool when the
// owner already approved it on g (approved): planque reads nothing through
// a vendor's tools before the owner approves them (kiff/planque#23). A
// narrowed write without it is refused; a narrowed read isn't checked. g
// is nil for a grant not made yet. It returns what to refuse, or "".
func (s *Server) operationsCheck(ctx context.Context, g *Grant, approved, tools []ToolRequest) (string, error) {
	read := map[string]map[string]*field{}
	for _, t := range tools {
		d, op, write, ok := narrowedTo(t)
		if !ok {
			continue
		}
		name := cmp.Or(t.Name, t.Tool)
		if g == nil || !slices.ContainsFunc(approved, func(a ToolRequest) bool { return a.Tool == d.Tool }) {
			if write {
				return fmt.Sprintf("%s narrows %s to %s, whose fields come from %s: connect %s first (with setup_reads), and add %s once the owner approved it",
					name, t.Tool, op, d.Tool, d.Tool, name), nil
			}
			continue
		}
		// Read once per operation and account.
		key := op
		for _, k := range d.Context {
			key += "\n" + string(t.Fixed[k])
		}
		fields, ok := read[key]
		if !ok {
			var why string
			var err error
			if fields, why, err = s.describe(ctx, g, d, op, t.Fixed); err != nil || why != "" {
				return why, err
			}
			read[key] = fields
		}
		// The marker planque writes into a write's fields is one of them too
		// (profileVendorSchema leaves it to this check).
		marker := ""
		if po, ok := profileOperation(activeProfile(g), name); ok && write && po.Marker != nil {
			marker = po.Marker.Field
		}
		if msg := fieldsCheck(t, d, op, fields, marker); msg != "" {
			return msg, nil
		}
	}
	return "", nil
}

// describe reads op's fields with d's tool, with g's token, in the account
// the narrowed tool pins (fixed, d.Context), and puts the read on g's
// record. why is the vendor's refusal.
func (s *Server) describe(ctx context.Context, g *Grant, d describer, op string, fixed map[string]json.RawMessage) (fields map[string]*field, why string, err error) {
	in := map[string]json.RawMessage{}
	for _, k := range d.Context {
		if v, ok := fixed[k]; ok {
			in[k] = v
		}
	}
	in[d.OpArg], _ = json.Marshal(op)
	args, _ := json.Marshal(in)
	access, ver, err := s.accessToken(ctx, g, 0)
	if err != nil {
		return nil, "", err
	}
	var res json.RawMessage
	for attempt := 0; ; attempt++ {
		res, err = s.MCP.CallTool(ctx, g.Request.Server, access, d.Tool, args)
		if !errors.Is(err, mcp.ErrUnauthorized) || attempt > 0 {
			break
		}
		if access, ver, err = s.accessToken(ctx, g, ver); err != nil {
			return nil, "", err
		}
	}
	if err != nil {
		return nil, "", fmt.Errorf("couldn't read %s's fields at %s: %w", op, hostOf(g.Request.Server), err)
	}
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if err := json.Unmarshal(res, &out); err != nil {
		return nil, "", fmt.Errorf("%s's description of %s isn't an MCP result: %w", hostOf(g.Request.Server), op, err)
	}
	var text strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			text.WriteString(c.Text)
		}
	}
	g.event(s.now(), fmt.Sprintf("read %s's fields with %s at %s, to check the tools asked for", op, d.Tool, hostOf(g.Request.Server)))
	if out.IsError {
		return nil, fmt.Sprintf("%s doesn't describe %s: %s", hostOf(g.Request.Server), op, strings.TrimSpace(text.String())), nil
	}
	return parseFields(text.String()), "", nil
}

// fieldsHold checks g's i-th tool, once, when it is a narrowed write not yet
// checked: against its operation's fields, read with the describing tool
// the owner approved on g, as adding it is checked now. It returns the
// refusal's code and reason, or "". A tool that holds is marked so; that
// mark failing to save only means it is checked again next time.
func (s *Server) fieldsHold(ctx context.Context, g *Grant, i int) (code, why string) {
	t := g.Tools[i]
	if _, _, write, ok := narrowedTo(t.ToolRequest); !ok || !write || t.FieldsChecked {
		return "", ""
	}
	approved := make([]ToolRequest, len(g.Tools))
	for j, gt := range g.Tools {
		approved[j] = gt.ToolRequest
	}
	msg, err := s.operationsCheck(ctx, g, approved, []ToolRequest{t.ToolRequest})
	if err != nil {
		s.log().WarnContext(ctx, "fields unread; refusing", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("err", err.Error()))
		return CodeDefinitionUnread, "planque couldn't read the operation's fields to check " + t.Name
	}
	if msg != "" {
		return CodeDefinitionInvalid, t.Name + " needs changing: " + msg
	}
	g.Tools[i].FieldsChecked = true
	if err := s.saveGrant(ctx, g); err != nil {
		s.log().WarnContext(ctx, "fields checked; mark not saved", slog.String("grant", g.ID), slog.String("tool", t.Name), slog.String("err", err.Error()))
	}
	return "", ""
}
