// Package composio reuses Composio's catalog, authentication and tools via
// its versioned HTTP API. It does not embed another agent runtime or SDK.
package composio

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
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
	scoped     map[string]authConfig // the scoped sign-ins planque made (Scope), by name
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
	// Managed: the provider's own OAuth app, not an operator's.
	Managed bool
}

// scopedPrefix names the sign-ins Scope makes. One is never an app's usual
// sign-in: it asks only for some tools' scopes.
const scopedPrefix = "planque-scoped-"

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
	scoped := map[string]authConfig{}
	cursor := ""
	for range 50 {
		q := url.Values{"limit": {"100"}}
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		var page struct {
			Items []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Scheme  string `json:"auth_scheme"`
				Status  string `json:"status"`
				Managed bool   `json:"is_composio_managed"`
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
				config := authConfig{ID: a.ID, Scheme: strings.ToUpper(a.Scheme), Managed: a.Managed}
				// A scoped sign-in is some tools' only, never the app's usual one.
				if strings.HasPrefix(a.Name, scopedPrefix) {
					if a.Managed {
						scoped[a.Name] = config
					}
					continue
				}
				// Stable across upstream page ordering when an app has several
				// enabled configs. The grant and hosted session pin this one.
				if old, ok := m[slug]; !ok || a.ID < old.ID {
					m[slug] = config
				}
			}
		}
		if page.Next == "" || page.Next == cursor {
			c.mu.Lock()
			c.configured, c.scoped, c.readAt = m, scoped, now
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
		app.AuthConfig, app.AuthScheme, app.Setup, app.Managed = config.ID, config.Scheme, connector.SetupReady, config.Managed
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
	config := authConfig{ID: made.Config.ID, Scheme: scheme, Managed: true}
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
	app.AuthConfig, app.AuthScheme, app.Setup, app.Managed = config.ID, config.Scheme, connector.SetupReady, config.Managed
	return app, nil
}

// Scope finds the narrowest sign-in to app's provider-managed sign-in that
// covers tools (gap 5 of kiff/tafs#476). It reads the app's managed scope
// ceiling and the scopes always asked, asks the provider which of the
// ceiling's scopes the tools need at least, and pins a managed sign-in that
// asks only for those: made once per app and scope set, and found again by
// its name. A new one is paid from the allowance first. ErrUnscoped for
// anything it can't do, so the caller keeps the app's usual sign-in.
func (c *Client) Scope(ctx context.Context, app connector.App, tools []string, budget connector.Allowance) (connector.Scoped, error) {
	if !app.Managed || app.AuthScheme != "OAUTH2" || !identifier.MatchString(app.Slug) || len(tools) == 0 {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	for _, t := range tools {
		if !identifier.MatchString(t) {
			return connector.Scoped{}, connector.ErrDefinition
		}
	}
	var kit struct {
		Managed []struct {
			Mode   string `json:"mode"`
			Scopes struct {
				Available []string `json:"available"`
			} `json:"scopes"`
		} `json:"composio_managed_auth"`
		Details []struct {
			Mode     string   `json:"mode"`
			Required []string `json:"required_scopes"`
		} `json:"auth_config_details"`
	}
	if _, err := c.request(ctx, "GET", "/toolkits/"+app.Slug, nil, nil, &kit); err != nil {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	var ceiling, always []string
	for _, m := range kit.Managed {
		if strings.EqualFold(m.Mode, "OAUTH2") {
			ceiling = cleanScopes(m.Scopes.Available)
		}
	}
	for _, d := range kit.Details {
		if strings.EqualFold(d.Mode, "OAUTH2") {
			always = cleanScopes(d.Required)
		}
	}
	if len(ceiling) == 0 {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	var rec struct {
		Scopes struct {
			Least []string `json:"least_privilege"`
		} `json:"scopes"`
	}
	body := map[string]any{"tools": tools, "auth_scheme": "OAUTH2", "available_scopes": ceiling}
	// A tool its ceiling can't cover (scope_not_available), or the beta
	// failing: wide, the app's usual access less the personal data its
	// usual sign-in asks for (14 Nov: a drafts-only Gmail taf's sign-in
	// asked for a date of birth, addresses, phone numbers and contacts).
	if _, err := c.request(ctx, "POST", "/toolkits/"+app.Slug+"/scopes/recommended", nil, body, &rec); err != nil {
		return c.wide(ctx, app, ceiling, always, budget)
	}
	least := cleanScopes(rec.Scopes.Least)
	if len(least) == 0 || slices.ContainsFunc(least, func(s string) bool { return !slices.Contains(ceiling, s) }) {
		return c.wide(ctx, app, ceiling, always, budget)
	}
	id, err := c.pin(ctx, app, least, budget)
	if err != nil {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	return connector.Scoped{AuthConfig: id, Scopes: cleanScopes(append(slices.Clone(least), always...))}, nil
}

// wide is app's sign-in when it can't be narrowed to the tools: its managed
// ceiling less the personal scopes (personalScope), pinned as a sign-in of
// its own, and said to be wide. ErrUnscoped when that leaves out nothing
// (the usual sign-in is then as narrow) or nothing is left.
func (c *Client) wide(ctx context.Context, app connector.App, ceiling, always []string, budget connector.Allowance) (connector.Scoped, error) {
	plain := slices.DeleteFunc(slices.Clone(ceiling), personalScope)
	if len(plain) == 0 || len(plain) == len(ceiling) {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	id, err := c.pin(ctx, app, plain, budget)
	if err != nil {
		return connector.Scoped{}, connector.ErrUnscoped
	}
	return connector.Scoped{AuthConfig: id, Scopes: cleanScopes(append(slices.Clone(plain), slices.DeleteFunc(slices.Clone(always), personalScope)...)), Wide: true}, nil
}

// personal are the scopes that read a person's own data and no app's tool
// needs to sign in: Google's People scopes for a date of birth, addresses,
// phone numbers, other email addresses, gender, organisation, age range,
// language, contacts and the directory. The account's own address
// (userinfo.email, openid) is never among them.
var personal = []string{"user.birthday.read", "user.addresses.read", "user.phonenumbers.read", "user.emails.read",
	"user.gender.read", "user.organization.read", "profile.agerange.read", "profile.language.read",
	"contacts", "contacts.readonly", "contacts.other.readonly", "directory.readonly"}

// personalScope says s reads a person's own data (personal), in Google's URL form.
func personalScope(s string) bool {
	name, ok := strings.CutPrefix(s, "https://www.googleapis.com/auth/")
	return ok && slices.Contains(personal, name)
}

// pin is the managed sign-in to app asking for scopes, made once per app
// and scope set and found again by its name. A new one is paid from the
// allowance first.
func (c *Client) pin(ctx context.Context, app connector.App, scopes []string, budget connector.Allowance) (string, error) {
	sum := sha256.Sum256([]byte(strings.Join(scopes, " ")))
	name := scopedPrefix + strings.ToLower(app.Slug) + "-" + hex.EncodeToString(sum[:8])
	if _, err := c.signInSetUp(ctx); err != nil {
		return "", err
	}
	c.mu.Lock()
	made, ok := c.scoped[name]
	c.mu.Unlock()
	if ok {
		return made.ID, nil
	}
	var out struct {
		Toolkit struct {
			Slug string `json:"slug"`
		} `json:"toolkit"`
		Config struct {
			ID      string `json:"id"`
			Scheme  string `json:"auth_scheme"`
			Managed bool   `json:"is_composio_managed"`
		} `json:"auth_config"`
	}
	// A new sign-in at the provider costs the account as an on-demand set-up does.
	if budget.Take != nil {
		if err := budget.Take(ctx); err != nil {
			return "", err
		}
	}
	giveBack := func() {
		if budget.GiveBack != nil {
			budget.GiveBack(ctx)
		}
	}
	config := map[string]any{"type": "use_composio_managed_auth", "name": name, "credentials": map[string]any{"scopes": scopes}}
	if _, err := c.request(ctx, "POST", "/auth_configs", nil, map[string]any{"toolkit": map[string]string{"slug": app.Slug}, "auth_config": config}, &out); err != nil {
		giveBack()
		return "", err
	}
	if !strings.EqualFold(out.Toolkit.Slug, app.Slug) || !out.Config.Managed || !identifier.MatchString(out.Config.ID) || !strings.EqualFold(out.Config.Scheme, "OAUTH2") {
		giveBack()
		return "", connector.ErrUnscoped
	}
	c.mu.Lock()
	next := make(map[string]authConfig, len(c.scoped)+1)
	for k, v := range c.scoped {
		next[k] = v
	}
	next[name] = authConfig{ID: out.Config.ID, Scheme: "OAUTH2", Managed: true}
	c.scoped = next
	c.mu.Unlock()
	return out.Config.ID, nil
}

// cleanScopes is scopes trimmed, without empty or repeated ones, sorted.
func cleanScopes(scopes []string) []string {
	out := []string{}
	for _, s := range scopes {
		if s = strings.TrimSpace(s); s != "" && len(s) <= 512 && !strings.ContainsAny(s, " \t\r\n") && !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	slices.Sort(out)
	return out
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
		// The vendor account's name, where the provider keeps it readable;
		// its id and mode it redacts.
		Data  map[string]json.RawMessage `json:"data"`
		State struct {
			Val map[string]json.RawMessage `json:"val"`
		} `json:"state"`
	}
	if _, err := c.request(ctx, "GET", "/connected_accounts/"+id, nil, nil, &a); err != nil {
		return connector.Account{}, err
	}
	if a.ID != id || !identifier.MatchString(a.User) || !identifier.MatchString(a.Toolkit.Slug) {
		return connector.Account{}, connector.ErrAccount
	}
	name := ""
	for _, m := range []map[string]json.RawMessage{a.Data, a.State.Val} {
		for _, k := range []string{"displayName", "display_name"} {
			if v, ok := m[k]; ok && name == "" {
				if n := connector.Reason(v); n != "REDACTED" {
					name = n
				}
			}
		}
	}
	return connector.Account{ID: a.ID, User: a.User, Toolkit: a.Toolkit.Slug, Status: a.Status, Name: name, Scopes: grantedScopes(a.Data, a.State.Val)}, nil
}

// grantedScopes are the scopes the vendor's token answer granted, where the
// provider keeps them readable on the connected account ("scope", as OAuth
// gives it, space-separated; or a list under "scopes"); nil when it doesn't.
func grantedScopes(fields ...map[string]json.RawMessage) []string {
	for _, m := range fields {
		for _, k := range []string{"scope", "scopes"} {
			v, ok := m[k]
			if !ok {
				continue
			}
			var one string
			var list []string
			switch {
			case json.Unmarshal(v, &one) == nil && one != "REDACTED":
				list = strings.FieldsFunc(one, func(r rune) bool { return r == ' ' || r == ',' })
			case json.Unmarshal(v, &list) == nil:
			}
			if scopes := cleanScopes(list); len(scopes) > 0 && len(scopes) <= 100 {
				return scopes
			}
		}
	}
	return nil
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
