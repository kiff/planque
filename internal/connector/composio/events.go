package composio

import (
	"context"
	"encoding/json"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/jcs"
)

var _ connector.EventProvider = (*Client)(nil)

func (c *Client) TriggerType(ctx context.Context, toolkit, slug, version string) (connector.TriggerType, error) {
	if !identifier.MatchString(toolkit) || !identifier.MatchString(slug) {
		return connector.TriggerType{}, connector.ErrDefinition
	}
	q := url.Values{"toolkit_versions": {"latest"}}
	if version != "" {
		if !pinnedVersion.MatchString(version) {
			return connector.TriggerType{}, connector.ErrDefinition
		}
		q.Del("toolkit_versions")
		q.Set("toolkit_versions["+toolkit+"]", version)
	}
	var wire struct {
		Slug    string `json:"slug"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Toolkit struct {
			Slug string `json:"slug"`
		} `json:"toolkit"`
		Config  json.RawMessage `json:"config"`
		Payload json.RawMessage `json:"payload"`
	}
	if _, err := c.request(ctx, "GET", "/triggers_types/"+slug, q, nil, &wire); err != nil {
		return connector.TriggerType{}, err
	}
	result := connector.TriggerType{Slug: wire.Slug, Toolkit: wire.Toolkit.Slug, Version: wire.Version, Name: wire.Name, Config: wire.Config, Payload: wire.Payload}
	if result.Validate() != nil || result.Slug != slug || result.Toolkit != toolkit || !pinnedVersion.MatchString(result.Version) || version != "" && result.Version != version {
		return connector.TriggerType{}, connector.ErrDefinition
	}
	return result, nil
}
func (c *Client) CreateTrigger(ctx context.Context, id connector.Identity, account string, approved connector.TriggerType, config json.RawMessage) (string, error) {
	if approved.Validate() != nil || approved.Toolkit != id.Toolkit {
		return "", connector.ErrDefinition
	}
	var args map[string]json.RawMessage
	if json.Unmarshal(config, &args) != nil || args == nil {
		return "", connector.ErrDefinition
	}
	a, err := c.Account(ctx, account)
	if err != nil {
		return "", err
	}
	if !a.Matches(id) {
		return "", connector.ErrAccount
	}
	current, err := c.TriggerType(ctx, id.Toolkit, approved.Slug, approved.Version)
	if err != nil {
		return "", err
	}
	// Wording does not reinterpret an approved schema. A changed config or
	// payload definition requires approval before registering the listener.
	definition := func(t connector.TriggerType) (string, error) {
		b, err := json.Marshal([]any{t.Slug, t.Toolkit, t.Version, t.Config, t.Payload})
		if err != nil {
			return "", err
		}
		return jcs.SHA256(b)
	}
	want, e1 := definition(approved)
	got, e2 := definition(current)
	if e1 != nil || e2 != nil || want != got {
		return "", connector.ErrDefinition
	}
	var result struct {
		ID string `json:"trigger_id"`
	}
	status, err := c.request(ctx, "POST", "/trigger_instances/"+approved.Slug+"/upsert", nil, map[string]any{"user_id": id.User, "connected_account_id": account, "trigger_config": config, "toolkit_versions": map[string]string{id.Toolkit: approved.Version}}, &result)
	if status == 403 {
		return "", connector.ErrPermission
	}
	if err != nil {
		return "", err
	}
	if !identifier.MatchString(result.ID) {
		return "", connector.ErrDefinition
	}
	return result.ID, nil
}
func (c *Client) DeleteTrigger(ctx context.Context, id string) error {
	if !identifier.MatchString(id) {
		return connector.ErrDefinition
	}
	status, err := c.request(ctx, "DELETE", "/trigger_instances/manage/"+id, nil, nil, nil)
	if status == 404 || status == 410 {
		return nil
	}
	if status == 403 {
		return connector.ErrPermission
	}
	return err
}

type wireWebhook struct {
	ID      string   `json:"id"`
	URL     string   `json:"webhook_url"`
	Version string   `json:"version"`
	Events  []string `json:"enabled_events"`
	Secret  string   `json:"secret"`
}

func (w wireWebhook) webhook(callback string) (connector.Webhook, error) {
	if !identifier.MatchString(w.ID) || w.URL != callback || w.Version != "V3" || !slices.Contains(w.Events, "composio.trigger.message") || w.Secret == "" || len(w.Secret) > 4096 {
		return connector.Webhook{}, connector.ErrWebhook
	}
	return connector.Webhook{ID: w.ID, URL: w.URL, Version: w.Version, Events: w.Events, Secret: w.Secret}, nil
}
func (c *Client) EnsureWebhook(ctx context.Context, callback string) (connector.Webhook, error) {
	u, err := url.Parse(callback)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return connector.Webhook{}, connector.ErrWebhook
	}
	var page struct {
		Items []wireWebhook `json:"items"`
		Next  string        `json:"next_cursor"`
	}
	status, err := c.request(ctx, "GET", "/webhook_subscriptions", url.Values{"limit": {"50"}}, nil, &page)
	if status == 403 {
		return connector.Webhook{}, connector.ErrPermission
	}
	if err != nil {
		return connector.Webhook{}, err
	}
	// The project currently allows one subscription. Never replace another
	// application's destination or discard its event filters.
	if len(page.Items) > 1 || page.Next != "" {
		return connector.Webhook{}, connector.ErrWebhook
	}
	if len(page.Items) == 1 {
		return page.Items[0].webhook(callback)
	}
	var result wireWebhook
	status, err = c.request(ctx, "POST", "/webhook_subscriptions", nil, map[string]any{"webhook_url": callback, "enabled_events": []string{"composio.trigger.message"}, "version": "V3"}, &result)
	if status == 403 {
		return connector.Webhook{}, connector.ErrPermission
	}
	if status == 409 {
		// A concurrent caller may have created the one project subscription.
		// Reconcile with one GET; never resend an uncertain create request.
		var existing struct {
			Items []wireWebhook `json:"items"`
		}
		if _, e := c.request(ctx, "GET", "/webhook_subscriptions", nil, nil, &existing); e != nil || len(existing.Items) != 1 {
			return connector.Webhook{}, connector.ErrWebhook
		}
		return existing.Items[0].webhook(callback)
	}
	if err != nil {
		return connector.Webhook{}, err
	}
	return result.webhook(callback)
}

// triggerTypesLimit is the most trigger types the provider lists a page.
const triggerTypesLimit = 50

// TriggerTypes lists an app's events, a page at a time (kiff/planque#87):
// the definitions as TriggerType reads them one by one, pinned to a version,
// so Coach can pick an event as it picks an action.
func (c *Client) TriggerTypes(ctx context.Context, toolkit, cursor string) (connector.Page[connector.TriggerType], error) {
	if !identifier.MatchString(toolkit) || len(cursor) > 4096 {
		return connector.Page[connector.TriggerType]{}, connector.ErrDefinition
	}
	q := url.Values{"toolkit_slugs": {toolkit}, "toolkit_versions": {"latest"}, "limit": {strconv.Itoa(triggerTypesLimit)}}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var page struct {
		Items []struct {
			Slug        string `json:"slug"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Version     string `json:"version"`
			Toolkit     struct {
				Slug string `json:"slug"`
			} `json:"toolkit"`
			Config  json.RawMessage `json:"config"`
			Payload json.RawMessage `json:"payload"`
		} `json:"items"`
		Next string `json:"next_cursor"`
	}
	if _, err := c.request(ctx, "GET", "/triggers_types", q, nil, &page); err != nil {
		return connector.Page[connector.TriggerType]{}, err
	}
	result := connector.Page[connector.TriggerType]{Items: []connector.TriggerType{}, Next: page.Next}
	if len(page.Items) > triggerTypesLimit || len(page.Next) > 4096 || page.Next != "" && page.Next == cursor {
		return connector.Page[connector.TriggerType]{}, connector.ErrDefinition
	}
	for _, w := range page.Items {
		t := connector.TriggerType{Slug: w.Slug, Toolkit: w.Toolkit.Slug, Version: w.Version, Name: w.Name, Config: w.Config, Payload: w.Payload}
		if len(w.Description) <= 20000 {
			t.Description = w.Description
		}
		if t.Validate() != nil || !strings.EqualFold(t.Toolkit, toolkit) || !pinnedVersion.MatchString(t.Version) {
			result.Skipped++
			continue
		}
		result.Items = append(result.Items, t)
	}
	return result, nil
}
