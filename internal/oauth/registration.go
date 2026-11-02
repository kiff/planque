package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// Registration is configured by the operator for one registered app.
// Secret is deliberately excluded from JSON: grants, previews and exported
// approvals must never contain app credentials. The operator's sealed store
// supplies it only while making a token or revocation request.
type Registration struct {
	ClientID      string `json:"client_id"`
	AuthMethod    string `json:"token_endpoint_auth_method"`
	TokenURL      string `json:"token_endpoint"`
	RevocationURL string `json:"revocation_endpoint,omitempty"`
	Secret        string `json:"-"`
}

func (c *Client) formRequest(ctx context.Context, endpoint string, form url.Values, revoke bool, registrations []Registration) (*http.Request, error) {
	if len(registrations) > 1 {
		return nil, errors.New("one client registration is required")
	}
	var app Registration
	if len(registrations) == 1 {
		app = registrations[0]
		want := app.TokenURL
		if revoke {
			want = app.RevocationURL
		}
		// A discovery document or grant cannot move the registered app's
		// secret to a different token endpoint or client id.
		if app.ClientID == "" || app.ClientID != form.Get("client_id") || want == "" || endpoint != want || !publicHTTPS(want, false) {
			return nil, errors.New("the registered client does not match this endpoint and client id")
		}
		switch app.AuthMethod {
		case "none":
			if app.Secret != "" {
				return nil, errors.New("a public client must not supply an app secret")
			}
		case "client_secret_post":
			if app.Secret == "" {
				return nil, errors.New("the registered app secret is unavailable")
			}
			form.Set("client_secret", app.Secret)
		case "client_secret_basic":
			if app.Secret == "" {
				return nil, errors.New("the registered app secret is unavailable")
			}
			// RFC 6749 uses one client authentication method per request.
			form.Del("client_id")
		default:
			return nil, errors.New("the registered client authentication method is unsupported")
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if app.AuthMethod == "client_secret_basic" {
		req.SetBasicAuth(url.QueryEscape(app.ClientID), url.QueryEscape(app.Secret))
	}
	return req, nil
}

// credentialHTTP never forwards a code, refresh token or app secret to a
// redirect target, including a 307 that would otherwise replay the body.
func (c *Client) credentialHTTP() *http.Client {
	client := *c.hc()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &client
}
