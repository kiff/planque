// Package appconfig holds reviewed, credential-free definitions for
// registered vendor applications. Credential references are configured
// separately; agents cannot choose endpoints for app-secret requests.
package appconfig

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/rest"
)

const (
	MCP  = "mcp_oauth_app"
	REST = "rest_oauth_app"
)

// Application contains public protocol facts only. For REST applications,
// Snapshot consists of reviewed operations, cited by Source. Their data
// uses the same request builder as an OpenAPI snapshot.
type Application struct {
	Name          string            `json:"name"`
	Kind          string            `json:"kind"`
	Server        string            `json:"server"`
	Source        string            `json:"source"`
	OAuth         oauth.Metadata    `json:"oauth"`
	AuthMethod    string            `json:"token_endpoint_auth_method"`
	Permissions   []string          `json:"permission_words"`
	VersionHeader string            `json:"version_header,omitempty"`
	Snapshot      *openapi.Snapshot `json:"operations,omitempty"`
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

// Decode accepts only the documented configuration, with a finite size.
func Decode(raw []byte) (map[string]Application, error) {
	if len(raw) > openapi.MaxBytes {
		return nil, errors.New("registered app definitions exceed 16 MB")
	}
	apps := map[string]Application{}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&apps) != nil || d.Decode(new(any)) != io.EOF || apps == nil || len(apps) > 500 {
		return nil, errors.New("registered app definitions must be a JSON object of documented fields")
	}
	for id, app := range apps {
		if !idPattern.MatchString(id) {
			return nil, errors.New("registered app ids must use lowercase letters, numbers, underscores or hyphens")
		}
		if err := app.Validate(); err != nil {
			return nil, fmt.Errorf("registered app %s: %w", id, err)
		}
	}
	return apps, nil
}

func (a Application) Validate() error {
	if strings.TrimSpace(a.Name) == "" || a.Kind != MCP && a.Kind != REST {
		return errors.New("name and a registered MCP or REST connection kind are required")
	}
	for _, raw := range []string{a.Server, a.Source, a.OAuth.Issuer, a.OAuth.Authorization, a.OAuth.Token} {
		if !https(raw) {
			return errors.New("vendor and OAuth endpoints must be fixed HTTPS URLs")
		}
	}
	if a.OAuth.Revocation != "" && !https(a.OAuth.Revocation) {
		return errors.New("revocation must use a fixed HTTPS URL")
	}
	if a.AuthMethod != "none" && a.AuthMethod != "client_secret_post" && a.AuthMethod != "client_secret_basic" {
		return errors.New("unsupported client authentication method")
	}
	if a.Kind == MCP && (a.OAuth.NoPKCE || a.OAuth.NoResourceParameter || a.OAuth.Resource != a.Server) {
		return errors.New("MCP applications require PKCE and an exact resource binding")
	}
	if a.OAuth.NoPKCE && a.AuthMethod == "none" {
		return errors.New("a public client requires PKCE")
	}
	for _, scope := range a.OAuth.Scopes {
		if scope == "" || len(scope) > 1024 || strings.ContainsAny(scope, " \r\n\t\x00") {
			return errors.New("invalid scope")
		}
	}
	if len(a.OAuth.Scopes) > 0 && len(a.Permissions) == 0 {
		return errors.New("requested access needs permission words for the approval")
	}
	for _, words := range a.Permissions {
		if strings.TrimSpace(words) == "" || len(words) > 240 || strings.ContainsAny(words, "\r\n\x00") {
			return errors.New("permission words must be short, single-line descriptions")
		}
	}
	for name, value := range a.OAuth.AuthorizationParams {
		switch name {
		case "access_type", "prompt", "include_granted_scopes", "audience":
		default:
			return errors.New("unsupported authorization parameter")
		}
		if len(value) > 1024 || strings.ContainsAny(value, "\r\n\x00") {
			return errors.New("invalid authorization parameter")
		}
	}
	if a.VersionHeader != "" {
		if a.Kind != REST {
			return errors.New("only a REST application uses a version header")
		}
		if err := (rest.VersionPin{Name: a.VersionHeader, Value: "reviewed"}).Validate(); err != nil {
			return err
		}
	}
	if a.Kind == REST {
		if a.Snapshot == nil || len(a.Snapshot.Operations) == 0 || len(a.Snapshot.Operations) > openapi.MaxOperations {
			return errors.New("a REST application needs reviewed operations")
		}
		names := map[string]bool{}
		for _, operation := range a.Snapshot.Operations {
			if !openapi.ToolName.MatchString(operation.Name) || names[operation.Name] {
				return errors.New("reviewed operations need distinct tool names")
			}
			names[operation.Name] = true
			switch operation.Method {
			case "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE":
			default:
				return errors.New("unsupported reviewed operation method")
			}
			if operation.Method == "" || !strings.HasPrefix(operation.Path, "/") || strings.ContainsAny(operation.Path, "?#\r\n") || !json.Valid(operation.InputSchema) {
				return errors.New("invalid reviewed operation")
			}
			read := operation.Method == "GET" || operation.Method == "HEAD"
			if operation.ReadOnly != read {
				return errors.New("reviewed operation read-only state must match its method")
			}
		}
	}
	return nil
}

func https(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && u.RawQuery == "" && u.Fragment == ""
}
