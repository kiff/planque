package connector

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrPermission = errors.New("the connector project key does not allow this operation")
var ErrWebhook = errors.New("the project's webhook does not match this installation")

// TriggerType is a concrete upstream definition. Its version, configuration
// schema and payload schema are included in the device's subscription approval.
type TriggerType struct {
	Slug    string          `json:"slug"`
	Toolkit string          `json:"toolkit"`
	Version string          `json:"version"`
	Name    string          `json:"name"`
	Config  json.RawMessage `json:"config"`
	Payload json.RawMessage `json:"payload"`
	// Description is the provider's, for a listing (kiff/planque#87); it
	// isn't part of what the device approves.
	Description string `json:"description,omitempty"`
}

func (t TriggerType) Validate() error {
	if !slug.MatchString(t.Slug) || !slug.MatchString(t.Toolkit) || !version.MatchString(t.Version) || t.Version == "latest" || t.Name == "" || !json.Valid(t.Config) || !json.Valid(t.Payload) {
		return ErrDefinition
	}
	return nil
}

type Webhook struct {
	ID      string   `json:"id"`
	URL     string   `json:"webhook_url"`
	Version string   `json:"version"`
	Events  []string `json:"enabled_events"`
	// Secret is returned only to the server, sealed before storage, never JSON.
	Secret string `json:"-"`
}

// EventProvider is optional: direct MCP and other action providers need not
// implement event delivery. A trigger always names an exact connected account.
type EventProvider interface {
	TriggerType(context.Context, string, string, string) (TriggerType, error)
	// TriggerTypes lists an app's events, a page at a time (kiff/planque#87).
	TriggerTypes(context.Context, string, string) (Page[TriggerType], error)
	CreateTrigger(context.Context, Identity, string, TriggerType, json.RawMessage) (string, error)
	DeleteTrigger(context.Context, string) error
	EnsureWebhook(context.Context, string) (Webhook, error)
}
