// Package composio reuses Composio's catalog, authentication and tools via
// its versioned HTTP API. It does not embed another agent runtime or SDK.
package composio

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/connector"
)

const API = "https://backend.composio.dev/api/v3.1"
const maxResponse = 16 << 20

// Secret resolves only the operator-configured project credential.
type Secret func(context.Context) (string, error)
type Client struct {
	HTTP   *http.Client
	Secret Secret
	// Now is the clock for the sign-in cache; nil is time.Now.
	Now func() time.Time

	mu         sync.Mutex
	configured map[string]authConfig // enabled sign-in selected for each toolkit
	readAt     time.Time
	firstPage  connector.Page[connector.App] // the unsearched first page, for readAt's sign-ins
	firstAt    time.Time
}

// configuredFor is how long the project's auth configs are trusted before
// they're read again.
const configuredFor = 5 * time.Minute

type authConfig struct {
	ID     string
	Scheme string
}

// signInSetUp reads which toolkits the project can sign someone in to: those
// with an enabled auth config. An app without one can't be connected, so it
// is neither listed nor started. The answer is cached briefly; a failed read
// offers nothing.
func (c *Client) signInSetUp(ctx context.Context) (map[string]authConfig, error) {
	now := time.Now()
	if c.Now != nil {
		now = c.Now()
	}
	c.mu.Lock()
	if c.configured != nil && now.Sub(c.readAt) < configuredFor {
		m := c.configured
		c.mu.Unlock()
		return m, nil
	}
	c.mu.Unlock()
	m := map[string]authConfig{}
	cursor := ""
	for range 50 {
		q := url.Values{"limit": {"100"}}
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		var page struct {
			Items []struct {
				ID      string `json:"id"`
				Scheme  string `json:"auth_scheme"`
				Status  string `json:"status"`
				Toolkit struct {
					Slug string `json:"slug"`
				} `json:"toolkit"`
			} `json:"items"`
			Next string `json:"next_cursor"`
		}
		if _, err := c.request(ctx, "GET", "/auth_configs", q, nil, &page); err != nil {
			return nil, err
		}
		for _, a := range page.Items {
			if strings.EqualFold(a.Status, "ENABLED") && identifier.MatchString(a.Toolkit.Slug) && identifier.MatchString(a.ID) && identifier.MatchString(a.Scheme) {
				slug := strings.ToLower(a.Toolkit.Slug)
				// Stable across upstream page ordering when an app has several
				// enabled configs. The grant and hosted session pin this one.
				if old, ok := m[slug]; !ok || a.ID < old.ID {
					m[slug] = authConfig{ID: a.ID, Scheme: strings.ToUpper(a.Scheme)}
				}
			}
		}
		if page.Next == "" || page.Next == cursor {
			c.mu.Lock()
			c.configured, c.readAt = m, now
			c.mu.Unlock()
			return m, nil
		}
		cursor = page.Next
	}
	return nil, connector.ErrDefinition
}

func (c *Client) Platform() connector.Platform { return connector.Platform{Name: "Composio", URL: API} }

var _ connector.Provider = (*Client)(nil)

var identifier = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_-]{0,255}$`)

func (c *Client) request(ctx context.Context, method, path string, q url.Values, body, out any) (int, error) {
	return c.requestGuard(ctx, method, path, q, body, out, nil)
}
func (c *Client) requestGuard(ctx context.Context, method, path string, q url.Values, body, out any, live func() bool) (int, error) {
	if c.Secret == nil {
		return 0, connector.ErrUnavailable
	}
	key, err := c.Secret(ctx)
	if err != nil || strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n\x00") {
		return 0, connector.ErrUnavailable
	}
	var raw []byte
	if body != nil {
		raw, err = json.Marshal(body)
		if err != nil {
			return 0, connector.ErrUnavailable
		}
	}
	target := API + path
	if len(q) > 0 {
		target += "?" + q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(raw))
	if err != nil {
		return 0, connector.ErrUnavailable
	}
	req.Header.Set("x-api-key", key)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	client := *hc
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	// No application retries: a POST with a lost answer must not be sent again.
	if live != nil && !live() {
		return 0, connector.ErrPermit
	}
	call := &connector.CallError{Provider: "composio", Call: method + " " + path}
	res, err := client.Do(req)
	if err != nil {
		return 0, call
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, maxResponse+1))
	if err != nil || len(b) > maxResponse {
		return res.StatusCode, connector.ErrUnavailable
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		call.Status = res.StatusCode
		call.Message = connector.Reason(b, key)
		return res.StatusCode, call
	}
	encodedKey, _ := json.Marshal(key)
	b = bytes.ReplaceAll(b, encodedKey[1:len(encodedKey)-1], []byte("[redacted]"))
	if out != nil && json.Unmarshal(b, out) != nil {
		return res.StatusCode, connector.ErrUnavailable
	}
	return res.StatusCode, nil
}

type wireApp struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Meta struct {
		Description string `json:"description"`
		Logo        string `json:"logo"`
	} `json:"meta"`
	// How the provider signs people in to it: the schemes it manages with
	// its own credentials, every scheme it takes, and whether none is needed.
	ManagedSchemes []string `json:"composio_managed_auth_schemes"`
	AuthSchemes    []string `json:"auth_schemes"`
	NoAuth         bool     `json:"no_auth"`
}

func (a wireApp) app() (connector.App, error) {
	if !identifier.MatchString(a.Slug) || strings.TrimSpace(a.Name) == "" || len(a.Name) > 100 {
		return connector.App{}, connector.ErrDefinition
	}
	icon := a.Meta.Logo
	if u, err := url.Parse(icon); err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil {
		icon = ""
	}
	return connector.App{Slug: a.Slug, Name: a.Name, Description: a.Meta.Description, Icon: icon}, nil
}

// keySchemes are sign-ins the provider's hosted page collects from the owner
// (a key or token), with no developer credentials of anyone's behind them.
var keySchemes = map[string]bool{"API_KEY": true, "BEARER_TOKEN": true, "BASIC": true}

// setUp says how the app's sign-in stands (kiff/planque#87): ready when the
// project has an enabled sign-in for it; on demand when the provider signs
// people in with its own credentials, takes only a key or token, or needs
// no sign-in; needs setup otherwise. The scheme is the one a sign-in set up
// on demand will have, for the words the owner reads before it exists.
func (a wireApp) setUp(app *connector.App, config authConfig, ok bool) {
	if ok {
		app.AuthConfig, app.AuthScheme, app.Setup = config.ID, config.Scheme, connector.SetupReady
		return
	}
	schemes := func(list []string) []string {
		out := []string{}
		for _, s := range list {
			if s = strings.ToUpper(strings.TrimSpace(s)); identifier.MatchString(s) {
				out = append(out, s)
			}
		}
		return out
	}
	if managed := schemes(a.ManagedSchemes); len(managed) > 0 {
		app.AuthScheme, app.Setup = managed[0], connector.SetupOnDemand
		for _, s := range managed {
			if s == "OAUTH2" {
				app.AuthScheme = s
			}
		}
		return
	}
	if a.NoAuth {
		app.AuthScheme, app.Setup = "NO_AUTH", connector.SetupOnDemand
		return
	}
	all := schemes(a.AuthSchemes)
	for _, s := range all {
		if keySchemes[s] {
			app.AuthScheme, app.Setup = s, connector.SetupOnDemand
			return
		}
	}
	app.Setup = connector.SetupNeedsSetup
	if len(all) > 0 {
		app.AuthScheme = all[0]
	}
}

// catalogPages bounds the upstream pages one catalog read walks, at
// catalogLimit apps each (the provider's maximum): the whole catalog is
// usually one page. A read that stops at the bound with pages left says so
// (Page.Partial), so an app that wasn't seen isn't taken as absent.
const catalogPages, catalogLimit = 10, 1000

func (c *Client) Catalog(ctx context.Context, search, cursor string) (connector.Page[connector.App], error) {
	if len(search) > 200 || len(cursor) > 4096 {
		return connector.Page[connector.App]{}, connector.ErrDefinition
	}
	setUp, err := c.signInSetUp(ctx)
	if err != nil {
		return connector.Page[connector.App]{}, err
	}
	// The first page with no search is what every Connect apps opens on:
	// kept as long as the sign-ins it was filtered with.
	first := search == "" && cursor == ""
	if first {
		c.mu.Lock()
		cached, ok := c.firstPage, !c.firstAt.IsZero() && c.firstAt.Equal(c.readAt)
		c.mu.Unlock()
		if ok {
			return cached, nil
		}
	}
	// With no search, only the apps with a sign-in set up: the short list
	// Connect apps opens on. A search finds every app, each saying how it
	// can be set up (kiff/planque#87).
	every := search != ""
	result := connector.Page[connector.App]{Items: []connector.App{}}
	for range catalogPages {
		q := url.Values{"limit": {strconv.Itoa(catalogLimit)}, "sort_by": {"alphabetically"}, "include_deprecated": {"false"}}
		if search != "" {
			q.Set("search", search)
		}
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		var page struct {
			Items []json.RawMessage `json:"items"`
			Next  string            `json:"next_cursor"`
		}
		if _, err := c.request(ctx, "GET", "/toolkits", q, nil, &page); err != nil {
			return connector.Page[connector.App]{}, err
		}
		if len(page.Items) > catalogLimit || len(page.Next) > 4096 || page.Next != "" && page.Next == cursor {
			return result, connector.ErrDefinition
		}
		for _, raw := range page.Items {
			var a wireApp
			if json.Unmarshal(raw, &a) != nil {
				result.Skipped++
				continue
			}
			app, err := a.app()
			if err != nil {
				result.Skipped++
				continue
			}
			config, ok := setUp[strings.ToLower(app.Slug)]
			a.setUp(&app, config, ok)
			if !every && app.Setup != connector.SetupReady {
				continue
			}
			result.Items = append(result.Items, app)
		}
		result.Next, cursor = page.Next, page.Next
		if page.Next == "" {
			break
		}
	}
	result.Partial = result.Next != ""
	if first {
		c.mu.Lock()
		c.firstPage, c.firstAt = result, c.readAt
		c.mu.Unlock()
	}
	return result, nil
}
func (c *Client) App(ctx context.Context, id string) (connector.App, error) {
	if !identifier.MatchString(id) {
		return connector.App{}, connector.ErrDefinition
	}
	setUp, err := c.signInSetUp(ctx)
	if err != nil {
		return connector.App{}, err
	}
	var a wireApp
	if status, err := c.request(ctx, "GET", "/toolkits/"+id, nil, nil, &a); err != nil {
		// An app the provider doesn't have can't be set up or connected here.
		if status == http.StatusNotFound {
			return connector.App{}, connector.ErrNotSetUp
		}
		return connector.App{}, err
	}
	if a.Slug != id {
		return connector.App{}, connector.ErrDefinition
	}
	app, err := a.app()
	config, ok := setUp[strings.ToLower(id)]
	a.setUp(&app, config, ok)
	return app, err
}

// SetUp sets up the provider-managed sign-in for an app it can sign people
// in to with its own credentials (kiff/planque#87): one auth config, of the
// managed kind, with nothing of ours in it. The answer is taken only when
// the provider says it manages it; the sign-in is then pinned by every
// grant as a pre-set app's is. The cached sign-ins take it at once, and the
// first catalog page is read again, so the app is offered from now on.
func (c *Client) SetUp(ctx context.Context, id string) (connector.App, error) {
	if !identifier.MatchString(id) {
		return connector.App{}, connector.ErrDefinition
	}
	// The sign-ins are read afresh first: one another instance set up a
	// moment ago is taken as it is, not made twice.
	c.mu.Lock()
	c.readAt = time.Time{}
	c.mu.Unlock()
	app, err := c.App(ctx, id)
	if err != nil {
		return connector.App{}, err
	}
	if app.Setup == connector.SetupReady {
		return app, nil
	}
	if app.Setup != connector.SetupOnDemand {
		return connector.App{}, connector.ErrNotSetUp
	}
	var made struct {
		Toolkit struct {
			Slug string `json:"slug"`
		} `json:"toolkit"`
		Config struct {
			ID      string `json:"id"`
			Scheme  string `json:"auth_scheme"`
			Managed bool   `json:"is_composio_managed"`
		} `json:"auth_config"`
	}
	body := map[string]any{"toolkit": map[string]string{"slug": id}, "auth_config": map[string]any{"type": "use_composio_managed_auth"}}
	if _, err := c.request(ctx, "POST", "/auth_configs", nil, body, &made); err != nil {
		return connector.App{}, err
	}
	scheme := strings.ToUpper(made.Config.Scheme)
	if !strings.EqualFold(made.Toolkit.Slug, id) || !made.Config.Managed || !identifier.MatchString(made.Config.ID) || !identifier.MatchString(scheme) {
		return connector.App{}, connector.ErrDefinition
	}
	config := authConfig{ID: made.Config.ID, Scheme: scheme}
	// The sign-ins handed out by signInSetUp are read without the lock, so
	// the map they hold is never written: a new map takes its place.
	c.mu.Lock()
	next := make(map[string]authConfig, len(c.configured)+1)
	for k, v := range c.configured {
		next[k] = v
	}
	if old, ok := next[strings.ToLower(id)]; !ok || config.ID < old.ID {
		next[strings.ToLower(id)] = config
	}
	config = next[strings.ToLower(id)]
	c.configured = next
	c.firstAt = time.Time{}
	c.mu.Unlock()
	app.AuthConfig, app.AuthScheme, app.Setup = config.ID, config.Scheme, connector.SetupReady
	return app, nil
}

// linkHosts serve the provider's hosted sign-in page; it moved from app. to
// connect. and both are accepted.
var linkHosts = map[string]bool{"app.composio.dev": true, "connect.composio.dev": true}

func (c *Client) Authorize(ctx context.Context, id connector.Identity, alias, callback string) (connector.Authorization, error) {
	if !identifier.MatchString(id.User) || !identifier.MatchString(id.Toolkit) || !identifier.MatchString(alias) || id.AuthConfig != "" && !identifier.MatchString(id.AuthConfig) {
		return connector.Authorization{}, connector.ErrDefinition
	}
	u, err := url.Parse(callback)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil {
		return connector.Authorization{}, connector.ErrDefinition
	}
	var s struct {
		ID     string `json:"session_id"`
		Config struct {
			User        string            `json:"user_id"`
			AuthConfigs map[string]string `json:"auth_configs"`
		} `json:"config"`
	}
	body := map[string]any{"user_id": id.User, "toolkits": map[string]any{"enable": []string{id.Toolkit}},
		"manage_connections": map[string]any{"enable": true, "enable_wait_for_connections": false, "enable_connection_removal": false},
		"multi_account":      map[string]any{"enable": true, "require_explicit_selection": true},
		"execute":            map[string]any{"enable_multi_execute": false}, "workbench": map[string]any{"enable": false}}
	if id.AuthConfig != "" {
		body["auth_configs"] = map[string]string{id.Toolkit: id.AuthConfig}
	}
	if _, err := c.request(ctx, "POST", "/tool_router/session", nil, body, &s); err != nil {
		return connector.Authorization{}, err
	}
	if !identifier.MatchString(s.ID) || s.Config.User != id.User {
		return connector.Authorization{}, connector.ErrAccount
	}
	if id.AuthConfig != "" && s.Config.AuthConfigs[id.Toolkit] != id.AuthConfig {
		return connector.Authorization{}, connector.ErrDefinition
	}
	var link struct {
		Account string `json:"connected_account_id"`
		URL     string `json:"redirect_url"`
	}
	if _, err := c.request(ctx, "POST", "/tool_router/session/"+s.ID+"/link", nil, map[string]string{"toolkit": id.Toolkit, "alias": alias, "callback_url": callback}, &link); err != nil {
		return connector.Authorization{}, err
	}
	target, err := url.Parse(link.URL)
	if err != nil || target.Scheme != "https" || target.User != nil || !linkHosts[target.Host] || !strings.HasPrefix(target.Path, "/link/") || !identifier.MatchString(link.Account) {
		return connector.Authorization{}, connector.ErrDefinition
	}
	return connector.Authorization{Session: s.ID, Account: link.Account, URL: link.URL}, nil
}

func (c *Client) Account(ctx context.Context, id string) (connector.Account, error) {
	if !identifier.MatchString(id) || id == "instant_account" {
		return connector.Account{}, connector.ErrAccount
	}
	var a struct {
		ID      string `json:"id"`
		User    string `json:"user_id"`
		Status  string `json:"status"`
		Toolkit struct {
			Slug string `json:"slug"`
		} `json:"toolkit"`
	}
	if _, err := c.request(ctx, "GET", "/connected_accounts/"+id, nil, nil, &a); err != nil {
		return connector.Account{}, err
	}
	if a.ID != id || !identifier.MatchString(a.User) || !identifier.MatchString(a.Toolkit.Slug) {
		return connector.Account{}, connector.ErrAccount
	}
	return connector.Account{ID: a.ID, User: a.User, Toolkit: a.Toolkit.Slug, Status: a.Status}, nil
}

func (c *Client) Revoke(ctx context.Context, id connector.Identity, account string) (connector.Revocation, error) {
	a, err := c.Account(ctx, account)
	if err != nil {
		return connector.Revocation{}, err
	}
	if a.User != id.User || a.Toolkit != id.Toolkit {
		return connector.Revocation{}, connector.ErrAccount
	}
	var receipt struct {
		Success bool   `json:"success"`
		Job     string `json:"revoke_job_id"`
	}
	_, err = c.request(ctx, "DELETE", "/connected_accounts/"+account, url.Values{"revoke_on_delete": {"true"}}, nil, &receipt)
	if err != nil || !receipt.Success {
		return connector.Revocation{}, connector.ErrUnavailable
	}
	// A successful delete disables this exact connection immediately. The
	// upstream revocation is a separate background job, never claimed done.
	return connector.Revocation{Disabled: true, Pending: true, Job: receipt.Job}, nil
}
