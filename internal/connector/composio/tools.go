package composio

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/jcs"
)

var pinnedVersion = regexp.MustCompile(`^[0-9]{8}_[0-9]{2,6}$`)

type wireTool struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Toolkit     struct {
		Slug string `json:"slug"`
	} `json:"toolkit"`
	Version    string          `json:"version"`
	Input      json.RawMessage `json:"input_parameters"`
	Tags       []string        `json:"tags"`
	Scopes     []string        `json:"scopes"`
	Deprecated bool            `json:"is_deprecated"`
}

func (w wireTool) tool() (connector.Tool, error) {
	schema, err := inputSchema(w.Input)
	if err != nil {
		return connector.Tool{}, err
	}
	t := connector.Tool{Slug: w.Slug, Toolkit: w.Toolkit.Slug, Version: w.Version, Name: w.Name, Description: w.Description, InputSchema: schema, Scopes: w.Scopes}
	for _, tag := range w.Tags {
		if tag == "readOnlyHint" {
			t.ReadOnly = true
		}
	}
	if !pinnedVersion.MatchString(w.Version) || w.Deprecated || len(w.Description) > 20000 {
		return connector.Tool{}, connector.ErrDefinition
	}
	return t, t.Validate()
}

// Both the JSON Schema representation and the documented flat parameter
// representation retain their constraints. Property names are never treated
// as schema keywords while converting the latter.
func inputSchema(raw json.RawMessage) (json.RawMessage, error) {
	var root map[string]json.RawMessage
	if len(raw) > 64<<10 || json.Unmarshal(raw, &root) != nil || root == nil {
		return nil, connector.ErrDefinition
	}
	var typ string
	if json.Unmarshal(root["type"], &typ) == nil && typ == "object" {
		return raw, nil
	}
	props := map[string]json.RawMessage{}
	required := []string{}
	for name, property := range root {
		var p map[string]json.RawMessage
		if json.Unmarshal(property, &p) != nil || p == nil {
			return nil, connector.ErrDefinition
		}
		if flag, ok := p["required"]; ok {
			var yes bool
			if json.Unmarshal(flag, &yes) != nil {
				return nil, connector.ErrDefinition
			}
			if yes {
				required = append(required, name)
			}
			delete(p, "required")
		}
		b, err := json.Marshal(p)
		if err != nil {
			return nil, connector.ErrDefinition
		}
		props[name] = b
	}
	sort.Strings(required)
	return json.Marshal(map[string]any{"type": "object", "properties": props, "required": required, "additionalProperties": false})
}

func (c *Client) Tools(ctx context.Context, toolkit, cursor string) (connector.Page[connector.Tool], error) {
	if !identifier.MatchString(toolkit) || len(cursor) > 4096 {
		return connector.Page[connector.Tool]{}, connector.ErrDefinition
	}
	q := url.Values{"toolkit_slug": {toolkit}, "limit": {"100"}, "include_deprecated": {"false"}}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var page struct {
		Items []json.RawMessage `json:"items"`
		Next  string            `json:"next_cursor"`
	}
	if _, err := c.request(ctx, "GET", "/tools", q, nil, &page); err != nil {
		return connector.Page[connector.Tool]{}, err
	}
	result := connector.Page[connector.Tool]{Items: []connector.Tool{}, Next: page.Next}
	if len(page.Items) > 100 || len(page.Next) > 4096 || page.Next != "" && page.Next == cursor {
		return result, connector.ErrDefinition
	}
	for _, raw := range page.Items {
		var wire wireTool
		if json.Unmarshal(raw, &wire) != nil {
			result.Skipped++
			continue
		}
		t, err := wire.tool()
		if err != nil || t.Toolkit != toolkit {
			result.Skipped++
			continue
		}
		result.Items = append(result.Items, t)
	}
	return result, nil
}
func (c *Client) Tool(ctx context.Context, slug, version string) (connector.Tool, error) {
	if !identifier.MatchString(slug) || strings.HasPrefix(slug, "COMPOSIO_") {
		return connector.Tool{}, connector.ErrDefinition
	}
	var wire wireTool
	if _, err := c.request(ctx, "GET", "/tools/"+slug, url.Values{"version": {version}}, nil, &wire); err != nil {
		return connector.Tool{}, err
	}
	t, err := wire.tool()
	if err != nil || t.Slug != slug || t.Version != version {
		return connector.Tool{}, connector.ErrDefinition
	}
	return t, nil
}

func toolDigest(t connector.Tool) (string, error) {
	scopes := append([]string(nil), t.Scopes...)
	sort.Strings(scopes)
	b, err := json.Marshal(struct {
		Slug, Toolkit, Version string
		Schema                 json.RawMessage
		ReadOnly               bool
		Scopes                 []string
	}{t.Slug, t.Toolkit, t.Version, t.InputSchema, t.ReadOnly, scopes})
	if err != nil {
		return "", err
	}
	return jcs.SHA256(b)
}

// Execute is called only after planque verifies the permit and claims its
// once-only operation. It independently checks ownership and the exact
// pinned tool schema before making the single execution request.
func (c *Client) Execute(ctx context.Context, id connector.Identity, account string, approved connector.Tool, args json.RawMessage, live func() bool) (connector.Result, error) {
	if approved.Validate() != nil || approved.Toolkit != id.Toolkit {
		return connector.Result{}, connector.ErrDefinition
	}
	var object map[string]json.RawMessage
	if json.Unmarshal(args, &object) != nil || object == nil {
		return connector.Result{}, connector.ErrDefinition
	}
	a, err := c.Account(ctx, account)
	if err != nil {
		return connector.Result{}, err
	}
	if !a.Matches(id) {
		return connector.Result{}, connector.ErrAccount
	}
	current, err := c.Tool(ctx, approved.Slug, approved.Version)
	if err != nil {
		return connector.Result{}, err
	}
	want, e1 := toolDigest(approved)
	got, e2 := toolDigest(current)
	if e1 != nil || e2 != nil {
		return connector.Result{}, connector.ErrDefinition
	}
	if want != got {
		return connector.Result{}, connector.ErrChanged
	}
	if live == nil || !live() {
		return connector.Result{}, connector.ErrPermit
	}
	var result struct {
		Data       json.RawMessage `json:"data"`
		Error      json.RawMessage `json:"error"`
		Successful *bool           `json:"successful"`
	}
	// Only these four documented fields are sent. Callers cannot select
	// another account, use instant credentials, provide custom auth or ask
	// the platform to turn natural language into a different action.
	body := map[string]any{"user_id": id.User, "connected_account_id": account, "version": approved.Version, "arguments": json.RawMessage(args)}
	status, err := c.requestGuard(ctx, "POST", "/tools/execute/"+approved.Slug, nil, body, &result, live)
	if errors.Is(err, connector.ErrPermit) {
		return connector.Result{}, err
	}
	if err != nil {
		// The provider itself answered with an error: its reason, if it gave one.
		var call *connector.CallError
		why := ""
		if errors.As(err, &call) {
			why = call.Message
		}
		if status < 300 || status >= 500 || status == 408 || status >= 300 && status < 400 {
			return connector.Result{Unknown: true, Reason: why, Status: status, AccountName: a.Name}, nil
		}
		return connector.Result{Refused: true, Reason: why, Status: status, AccountName: a.Name}, nil
	}
	refused := result.Successful != nil && !*result.Successful
	if len(result.Error) > 0 && string(result.Error) != "null" && string(result.Error) != `""` {
		refused = true
	}
	if len(result.Data) == 0 {
		result.Data = json.RawMessage(`{}`)
	}
	if !refused {
		return connector.Result{Data: result.Data, AccountName: a.Name}, nil
	}
	// The provider ran the action and the vendor said no: the vendor's own words,
	// and its HTTP status when the provider passes it on (a refused key's 403).
	why := connector.Reason(result.Error)
	if why == "" {
		why = connector.Reason(result.Data)
	}
	return connector.Result{Data: result.Data, Refused: true, Reason: why, ByVendor: true, Status: vendorStatus(result.Data), AccountName: a.Name}, nil
}

// vendorStatus is the vendor's HTTP status in what the provider returned with
// its refusal ({"status_code": 403, "http_error": …}), or 0 when it gives
// none that reads as one.
func vendorStatus(data json.RawMessage) int {
	var d struct {
		StatusCode json.Number `json:"status_code"`
	}
	if json.Unmarshal(data, &d) != nil {
		return 0
	}
	n, err := d.StatusCode.Int64()
	if err != nil || n < 100 || n > 599 {
		return 0
	}
	return int(n)
}

// searchLimit is how many tools one search answers at most.
const searchLimit = 50

// Search is the provider's text search over tools' names, slugs and
// descriptions (kiff/planque#87), narrowed to one app when exactly one is
// named and to the named apps here otherwise. Platform meta tools never come
// back (Validate), and each tool's app comes with it as it stands, read once
// per app, so a caller can say which apps are set up before connecting.
func (c *Client) Search(ctx context.Context, query string, toolkits []string, cursor string) (connector.Search, error) {
	query = strings.TrimSpace(query)
	if query == "" || len(query) > 200 || len(cursor) > 4096 || len(toolkits) > 20 {
		return connector.Search{}, connector.ErrDefinition
	}
	want := map[string]bool{}
	for _, t := range toolkits {
		if !identifier.MatchString(t) {
			return connector.Search{}, connector.ErrDefinition
		}
		want[strings.ToLower(t)] = true
	}
	q := url.Values{"query": {query}, "limit": {strconv.Itoa(searchLimit)}, "include_deprecated": {"false"}}
	if len(toolkits) == 1 {
		q.Set("toolkit_slug", toolkits[0])
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var page struct {
		Items []json.RawMessage `json:"items"`
		Next  string            `json:"next_cursor"`
	}
	if _, err := c.request(ctx, "GET", "/tools", q, nil, &page); err != nil {
		return connector.Search{}, err
	}
	result := connector.Search{Tools: []connector.Tool{}, Apps: map[string]connector.App{}, Next: page.Next}
	if len(page.Items) > searchLimit || len(page.Next) > 4096 || page.Next != "" && page.Next == cursor {
		return connector.Search{}, connector.ErrDefinition
	}
	for _, raw := range page.Items {
		var wire wireTool
		if json.Unmarshal(raw, &wire) != nil {
			result.Skipped++
			continue
		}
		t, err := wire.tool()
		if err != nil {
			result.Skipped++
			continue
		}
		if len(want) > 0 && !want[strings.ToLower(t.Toolkit)] {
			continue
		}
		if _, seen := result.Apps[t.Toolkit]; !seen {
			app, err := c.App(ctx, t.Toolkit)
			if err != nil {
				result.Skipped++
				continue
			}
			result.Apps[t.Toolkit] = app
		}
		result.Tools = append(result.Tools, t)
	}
	return result, nil
}
