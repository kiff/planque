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
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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
	// Message is the provider's own reason for a refusal, as Reason makes it
	// safe to show: one line, bounded, secrets redacted. Never in Error().
	Message string
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

// ErrChanged is ErrDefinition when the provider's current definition of an
// approved tool was read and differs from what the owner approved. Only
// this one means the tool changed: a definition that can't be read or
// checked is ErrDefinition alone, and says nothing about a change.
var ErrChanged = fmt.Errorf("%w", ErrDefinition)

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
	// Setup says whether the provider can sign someone in to the app now
	// (SetupReady), can once the sign-in is set up on the owner's pick
	// (SetupOnDemand), or needs an operator's own developer credentials first
	// (SetupNeedsSetup).
	Setup string `json:"-"`
	// Managed says the pinned sign-in is the provider's own OAuth app, not
	// an operator's: its scopes are the provider's to narrow (Scope).
	Managed bool `json:"-"`
}

// Scoped is a sign-in that asks only for what some tools need (kiff-cloud
// #1487, gap 5): pinned instead of the app's usual one.
type Scoped struct {
	AuthConfig string
	// Scopes are what it asks for, sorted: the tools' least-privilege scopes
	// and the ones the provider always asks (a refresh, who you are).
	Scopes []string
}

// ErrUnscoped: the provider can't narrow the app's sign-in to the tools (its
// own sign-in has no scope ceiling, a tool needs a scope outside it, or the
// provider couldn't say). The caller keeps the app's usual sign-in, and says so.
var ErrUnscoped = errors.New("this app's sign-in can't be narrowed to these tools")

// How an app's sign-in stands at the provider (kiff/planque#87).
const (
	// SetupReady: the project has an enabled sign-in for the app.
	SetupReady = "ready"
	// SetupOnDemand: the provider signs people in with its own credentials,
	// so the sign-in is set up the moment an owner picks the app (SetUp).
	SetupOnDemand = "on_demand"
	// SetupNeedsSetup: the app needs an operator's own developer app or
	// credentials before anyone can sign in; it's listed, and not started.
	SetupNeedsSetup = "needs_setup"
)

type Page[T any] struct {
	Items   []T    `json:"items"`
	Next    string `json:"next_cursor,omitempty"`
	Skipped int    `json:"skipped,omitempty"`
	// Partial: the read stopped at its bound with pages left, so an app
	// that wasn't seen isn't absent.
	Partial bool `json:"partial,omitempty"`
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

// Search is what a search by words found: the tools, and each tool's app as
// it stands, so the caller can say which app is set up, on demand or needs
// setting up before anything is connected.
type Search struct {
	Tools   []Tool         `json:"tools"`
	Apps    map[string]App `json:"apps"`
	Next    string         `json:"next_cursor,omitempty"`
	Skipped int            `json:"skipped,omitempty"`
}

type Account struct {
	ID      string `json:"id"`
	User    string `json:"user_id"`
	Toolkit string `json:"toolkit"`
	Status  string `json:"status"`
	// Name is the vendor account's name as the provider shows it (a
	// business's name), one safe line; empty when it shows none.
	Name string `json:"name,omitempty"`
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
	// Why a call was refused or left unknown, as Reason makes it safe to show,
	// and who said it: the vendor's own answer passed on by the provider
	// (ByVendor), or the provider, with its HTTP status when it gave one.
	Reason   string
	ByVendor bool
	Status   int
	// AccountName is the account's name as the provider read it before the
	// call (Account.Name).
	AccountName string
}

// maxReason bounds a reason: one line the owner reads, never a payload.
const maxReason = 300

// secretShapes are what a reason must never carry, whatever the vendor
// echoed: vendor keys, bearer tokens, webhook secrets and JWTs.
var secretShapes = regexp.MustCompile(`(?i)\b(?:sk|rk|pk)_(?:live|test)_[A-Za-z0-9]+|\bwhsec_[A-Za-z0-9]+|\bbearer\s+[A-Za-z0-9._~+/=-]+|\beyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)

// Reason is a provider's or vendor's error made safe to show the owner and
// the taf: its message when the error is JSON ({"message"}, {"error":
// {"message"}}, {"error": "…"} or a string), one line, control characters
// gone, secrets in secretShapes and any of secrets redacted, at most
// maxReason characters. Empty when there's nothing to say.
func Reason(raw []byte, secrets ...string) string {
	text := strings.TrimSpace(string(raw))
	var v any
	if json.Unmarshal(raw, &v) == nil {
		text = messageOf(v)
	}
	for _, s := range secrets {
		if s = strings.TrimSpace(s); s != "" {
			text = strings.ReplaceAll(text, s, "[redacted]")
		}
	}
	text = secretShapes.ReplaceAllString(text, "[redacted]")
	text = strings.Join(strings.FieldsFunc(text, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }), " ")
	if r := []rune(text); len(r) > maxReason {
		text = string(r[:maxReason-1]) + "…"
	}
	return text
}

func messageOf(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case map[string]any:
		for _, k := range []string{"message", "error", "detail", "error_description"} {
			if m := messageOf(x[k]); m != "" {
				return m
			}
		}
	}
	return ""
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
	// Catalog with no search lists the apps with a sign-in set up; with one,
	// every app matching it, each saying how it can be set up (App.Setup).
	Catalog(context.Context, string, string) (Page[App], error)
	App(context.Context, string) (App, error)
	// SetUp sets up the sign-in for an app the provider can sign people in
	// to with its own credentials (SetupOnDemand), once, and returns the app
	// ready, with the sign-in to pin. It never takes a credential.
	SetUp(context.Context, string) (App, error)
	Authorize(context.Context, Identity, string, string) (Authorization, error)
	Account(context.Context, string) (Account, error)
	Tools(context.Context, string, string) (Page[Tool], error)
	// Search finds tools across apps by words (kiff/planque#87): the
	// provider's own text search over names and descriptions, narrowed to
	// the apps named when any are, with each tool's app as it stands.
	Search(context.Context, string, []string, string) (Search, error)
	Tool(context.Context, string, string) (Tool, error)
	Execute(context.Context, Identity, string, Tool, json.RawMessage, func() bool) (Result, error)
	Revoke(context.Context, Identity, string) (Revocation, error)
	// Scope finds the narrowest sign-in to app's provider-managed sign-in
	// that covers tools (the provider's tool slugs), made once per app and
	// scope set. ErrUnscoped when it can't.
	Scope(context.Context, App, []string) (Scoped, error)
}
