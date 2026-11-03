// Package connector is the shared boundary to an existing connection
// platform. Providers supply authentication and concrete tools; planque
// retains the owner's approval, exact-call permit and execution record.
package connector

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrPermit = errors.New("the permit expired before the call could be sent")

var ErrUnavailable = errors.New("connector service is unavailable")

// CallError is ErrUnavailable with what the log needs to say which call the
// provider refused: the call (method and path, no query) and the status it
// answered, 0 for no answer. It never carries the body, the key or a query.
type CallError struct {
	Provider string
	Call     string
	Status   int
}

func (e *CallError) Error() string {
	status := "no answer"
	if e.Status != 0 {
		status = strconv.Itoa(e.Status)
	}
	return ErrUnavailable.Error() + ": " + e.Provider + ": " + e.Call + ": " + status
}

func (e *CallError) Unwrap() error { return ErrUnavailable }

var ErrAccount = errors.New("the selected connection is not active for this account and app")

// ErrNotSetUp: the provider can't sign anyone in to this app here (it has no
// enabled sign-in set up for it), so it isn't offered.
var ErrNotSetUp = errors.New("this app can't be connected here yet")
var ErrDefinition = errors.New("the connector tool definition changed; approve its current definition")

// Identity is allocated by planque, never selected in a tool call.
type Identity struct {
	User    string
	Toolkit string
	// AuthConfig pins the sign-in setup whose words the owner approved.
	// Empty for grants made before this binding was recorded.
	AuthConfig string
}

type App struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	// Provider metadata for a new grant, not catalog display fields.
	AuthConfig string `json:"-"`
	AuthScheme string `json:"-"`
}

type Page[T any] struct {
	Items   []T    `json:"items"`
	Next    string `json:"next_cursor,omitempty"`
	Skipped int    `json:"skipped,omitempty"`
}

// Tool is one concrete, versioned upstream action. Platform meta-executors
// cannot become granted business tools.
type Tool struct {
	Slug        string          `json:"slug"`
	Toolkit     string          `json:"toolkit"`
	Version     string          `json:"version"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema"`
	ReadOnly    bool            `json:"read_only"`
	Scopes      []string        `json:"scopes,omitempty"`
}

var slug = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_-]{0,255}$`)
var version = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$`)

func (t Tool) Validate() error {
	if !slug.MatchString(t.Slug) || !slug.MatchString(t.Toolkit) || !version.MatchString(t.Version) || strings.EqualFold(t.Version, "latest") || len(t.Name) == 0 || !json.Valid(t.InputSchema) {
		return ErrDefinition
	}
	if strings.HasPrefix(strings.ToUpper(t.Slug), "COMPOSIO_") || strings.EqualFold(t.Toolkit, "composio") {
		return ErrDefinition
	}
	var schema struct {
		Type       string                     `json:"type"`
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if json.Unmarshal(t.InputSchema, &schema) != nil || schema.Type != "object" {
		return ErrDefinition
	}
	return nil
}

// Alias fits a gateway tool name regardless of a provider's naming limits.
func (t Tool) Alias() string {
	sum := sha256.Sum256([]byte(t.Toolkit + "\n" + t.Slug))
	return "connector_" + hex.EncodeToString(sum[:16])
}

type Account struct {
	ID      string `json:"id"`
	User    string `json:"user_id"`
	Toolkit string `json:"toolkit"`
	Status  string `json:"status"`
}

func (a Account) Matches(id Identity) bool {
	return a.ID != "" && a.User == id.User && a.Toolkit == id.Toolkit && a.Status == "ACTIVE"
}

type Authorization struct {
	Session string
	Account string
	URL     string
}

// Revocation distinguishes disabling access from completing upstream
// revocation, which connection platforms may finish asynchronously.
type Revocation struct {
	Disabled bool   `json:"disabled"`
	Pending  bool   `json:"pending"`
	Job      string `json:"job,omitempty"`
}

type Result struct {
	Data json.RawMessage
	// Refused is a provider error. It never establishes that no effect occurred.
	Refused bool
	Unknown bool
}

// Platform describes the operator-configured service, never a caller URL.
type Platform struct {
	Name string
	URL  string
}

// Provider exposes concrete actions. Execute must check its final guard after
// preflight, immediately before sending, and never retry an uncertain write.
type Provider interface {
	Platform() Platform
	Catalog(context.Context, string, string) (Page[App], error)
	App(context.Context, string) (App, error)
	Authorize(context.Context, Identity, string, string) (Authorization, error)
	Account(context.Context, string) (Account, error)
	Tools(context.Context, string, string) (Page[Tool], error)
	Tool(context.Context, string, string) (Tool, error)
	Execute(context.Context, Identity, string, Tool, json.RawMessage, func() bool) (Result, error)
	Revoke(context.Context, Identity, string) (Revocation, error)
}
