package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
)

// ApplicationApproval freezes the public app configuration and client id
// in the signed request. Later configuration changes cannot reinterpret an
// existing grant. REST operations are stored separately by their digest.
type ApplicationApproval struct {
	ID         string                `json:"id"`
	ClientID   string                `json:"client_id"`
	Definition appconfig.Application `json:"definition"`
}

func (s *Server) applicationRequest(ctx context.Context, tenant, grant, id string, tools []ToolRequest) (*ApplicationApproval, *APIRequest, error) {
	app, ok := s.Applications[id]
	if !ok || app.Validate() != nil || s.ApplicationCredentials == nil {
		return nil, nil, apiError{http.StatusBadRequest, "that registered vendor app is not configured"}
	}
	credentials, err := s.ApplicationCredentials.Read(ctx, id)
	if err != nil || credentials.ClientID == "" || app.AuthMethod != "none" && credentials.Secret == "" {
		return nil, nil, appcredentials.ErrUnavailable
	}
	// Definitions include maps and schemas: freeze a copy, not the server's
	// live configuration object, before it becomes approval evidence.
	raw, _ := json.Marshal(app)
	var frozen appconfig.Application
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, nil, errors.New("registered app configuration could not be copied")
	}
	app = frozen
	var api *APIRequest
	if app.Kind == appconfig.REST {
		snapshot := app.Snapshot
		snapshot.Digest = ""
		raw, err := json.Marshal(snapshot)
		if err != nil {
			return nil, nil, err
		}
		snapshot.Digest, err = jcs.SHA256(raw)
		if err != nil {
			return nil, nil, err
		}
		if msg := apiToolsCheck(snapshot, tools); msg != "" {
			return nil, nil, apiError{http.StatusBadRequest, msg}
		}
		if err := s.saveSnapshot(ctx, tenant, grant, snapshot); err != nil {
			return nil, nil, err
		}
		var names []string
		for _, operation := range snapshot.Operations {
			names = append(names, operation.Name)
		}
		api = &APIRequest{URL: app.Source, Digest: snapshot.Digest, Base: app.Server, SourceKind: "reviewed",
			Auth: openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer "}, Operations: names}
		app.Snapshot = nil
	}
	return &ApplicationApproval{ID: id, ClientID: credentials.ClientID, Definition: app}, api, nil
}

// registration supplies the current secret for the approved client id,
// pinned to the approved endpoints. Secret rotation works without changing
// the grant; replacing the client id requires a newly approved request.
func (s *Server) registration(ctx context.Context, g *Grant) ([]oauth.Registration, error) {
	approved := g.Request.Application
	if approved == nil {
		return nil, nil
	}
	if s.ApplicationCredentials == nil || approved.Definition.Kind != g.Request.Kind || approved.Definition.Server != g.Request.Server {
		return nil, appcredentials.ErrUnavailable
	}
	meta, _ := json.Marshal(g.Vendor)
	pinned, _ := json.Marshal(approved.Definition.OAuth)
	if string(meta) != string(pinned) {
		return nil, errors.New("the grant's OAuth configuration does not match its approval")
	}
	credentials, err := s.ApplicationCredentials.Read(ctx, approved.ID)
	if err != nil || credentials.ClientID != approved.ClientID || g.ClientID != "" && credentials.ClientID != g.ClientID {
		return nil, appcredentials.ErrUnavailable
	}
	return []oauth.Registration{{ClientID: credentials.ClientID, Secret: credentials.Secret, AuthMethod: approved.Definition.AuthMethod,
		TokenURL: g.Vendor.Token, RevocationURL: g.Vendor.Revocation}}, nil
}

func (s *Server) exchangeToken(ctx context.Context, g *Grant, code, verifier string) (oauth.Token, error) {
	registration, err := s.registration(ctx, g)
	if err != nil {
		return oauth.Token{}, err
	}
	return s.OAuth.Exchange(ctx, g.Vendor, g.ClientID, s.redirect(), code, verifier, registration...)
}

func (s *Server) refreshToken(ctx context.Context, g *Grant, refresh string) (oauth.Token, error) {
	registration, err := s.registration(ctx, g)
	if err != nil {
		return oauth.Token{}, err
	}
	return s.OAuth.Refresh(ctx, g.Vendor, g.ClientID, refresh, registration...)
}

func (s *Server) revokeToken(ctx context.Context, g *Grant, token, hint string) error {
	registration, err := s.registration(ctx, g)
	if err != nil {
		return err
	}
	return s.OAuth.Revoke(ctx, g.Vendor, g.ClientID, token, hint, registration...)
}

func appWords(request Request) string {
	if request.Application == nil {
		return ""
	}
	app := request.Application.Definition
	var words strings.Builder
	words.WriteString("\n\nAccess at " + app.Name + ":")
	for _, permission := range app.Permissions {
		words.WriteString("\n• " + permission)
	}
	if request.API != nil {
		words.WriteString("\nThe approved operations are the reviewed definitions from " + app.Source + ".")
	}
	return words.String()
}
