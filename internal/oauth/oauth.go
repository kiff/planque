// Package oauth is MCP's authorization standard, from the client's side:
// protected resource metadata (RFC 9728), authorization server metadata
// (RFC 8414), dynamic client registration (RFC 7591), PKCE with S256
// (RFC 7636), resource indicators (RFC 8707) and revocation (RFC 7009).
// planque implements it once; any vendor speaking it works the same way.
package oauth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"
)

// Metadata is what planque needs to sign in to one vendor server.
type Metadata struct {
	Resource      string   `json:"resource"` // the MCP server, as its metadata names it
	Issuer        string   `json:"issuer"`
	Authorization string   `json:"authorization_endpoint"`
	Token         string   `json:"token_endpoint"`
	Registration  string   `json:"registration_endpoint"`
	Revocation    string   `json:"revocation_endpoint,omitempty"`
	Scopes        []string `json:"scopes,omitempty"`
	// IssParameter: the server always names itself in its answer (RFC 9207).
	IssParameter bool `json:"iss_parameter,omitempty"`
}

// Client speaks to authorization servers.
type Client struct {
	HTTP *http.Client
	// AllowHTTP admits http:// endpoints, for tests on localhost only.
	AllowHTTP bool
}

func (c *Client) hc() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 15 * time.Second}
}

// Discover finds how to sign in to the MCP server at server: the
// resource_metadata its 401 names (RFC 9728 §5.1), or the well-known
// location, then its first authorization server's RFC 8414 metadata.
func (c *Client) Discover(ctx context.Context, server string) (Metadata, error) {
	if err := c.checkURL(server); err != nil {
		return Metadata{}, err
	}
	prmURL, err := c.challenge(ctx, server)
	if err != nil {
		return Metadata{}, err
	}
	if prmURL == "" {
		prmURL = wellKnown(server, "oauth-protected-resource")
	}
	var prm struct {
		Resource             string   `json:"resource"`
		AuthorizationServers []string `json:"authorization_servers"`
		Scopes               []string `json:"scopes_supported"`
	}
	if err := c.getJSON(ctx, prmURL, &prm); err != nil {
		return Metadata{}, fmt.Errorf("protected resource metadata: %w", err)
	}
	// RFC 9728 §3.3: the metadata must be for the resource we called.
	if trim(prm.Resource) != trim(server) {
		return Metadata{}, fmt.Errorf("the server's metadata names resource %q, not %q", prm.Resource, server)
	}
	if len(prm.AuthorizationServers) == 0 {
		return Metadata{}, errors.New("the server's metadata names no authorization server")
	}
	issuer := prm.AuthorizationServers[0]
	if err := c.checkURL(issuer); err != nil {
		return Metadata{}, err
	}
	var as struct {
		Issuer        string   `json:"issuer"`
		Authorization string   `json:"authorization_endpoint"`
		Token         string   `json:"token_endpoint"`
		Registration  string   `json:"registration_endpoint"`
		Revocation    string   `json:"revocation_endpoint"`
		Challenge     []string `json:"code_challenge_methods_supported"`
		AuthMethods   []string `json:"token_endpoint_auth_methods_supported"`
		Grants        []string `json:"grant_types_supported"`
		IssParameter  bool     `json:"authorization_response_iss_parameter_supported"`
	}
	if err := c.getJSON(ctx, wellKnown(issuer, "oauth-authorization-server"), &as); err != nil {
		if err2 := c.getJSON(ctx, trim(issuer)+"/.well-known/openid-configuration", &as); err2 != nil {
			return Metadata{}, fmt.Errorf("authorization server metadata: %w", err)
		}
	}
	// RFC 8414 §3.3: the issuer must be the one we asked about.
	if as.Issuer != issuer {
		return Metadata{}, fmt.Errorf("the authorization server's metadata names issuer %q, not %q", as.Issuer, issuer)
	}
	switch {
	case !slices.Contains(as.Challenge, "S256"):
		return Metadata{}, errors.New("the authorization server does not offer PKCE with S256")
	case as.Registration == "":
		return Metadata{}, errors.New("the authorization server does not offer dynamic client registration")
	case len(as.AuthMethods) > 0 && !slices.Contains(as.AuthMethods, "none"):
		return Metadata{}, errors.New("the authorization server does not accept public clients")
	case len(as.Grants) > 0 && !slices.Contains(as.Grants, "authorization_code"):
		return Metadata{}, errors.New("the authorization server does not offer the authorization code grant")
	}
	for _, u := range []string{as.Authorization, as.Token, as.Registration} {
		if err := c.checkURL(u); err != nil {
			return Metadata{}, err
		}
	}
	if as.Revocation != "" {
		if err := c.checkURL(as.Revocation); err != nil {
			return Metadata{}, err
		}
	}
	return Metadata{Resource: prm.Resource, Issuer: as.Issuer, Authorization: as.Authorization, Token: as.Token,
		Registration: as.Registration, Revocation: as.Revocation, Scopes: prm.Scopes, IssParameter: as.IssParameter}, nil
}

// challenge makes an unauthenticated MCP request and reads the
// resource_metadata parameter of the 401's WWW-Authenticate.
func (c *Client) challenge(ctx context.Context, server string) (string, error) {
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"planque","version":"1"}}}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := c.hc().Do(req)
	if err != nil {
		return "", fmt.Errorf("the MCP server could not be reached: %w", err)
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("the MCP server answered %d without a token; planque connects servers that require sign-in", resp.StatusCode)
	}
	for _, h := range resp.Header.Values("WWW-Authenticate") {
		if v := authParam(h, "resource_metadata"); v != "" {
			if err := c.checkURL(v); err != nil {
				return "", err
			}
			return v, nil
		}
	}
	return "", nil
}

// authParam reads one auth-param from a WWW-Authenticate value.
func authParam(h, name string) string {
	i := strings.Index(h, name+"=")
	if i < 0 {
		return ""
	}
	v := h[i+len(name)+1:]
	if strings.HasPrefix(v, `"`) {
		v = v[1:]
		if j := strings.IndexByte(v, '"'); j >= 0 {
			return v[:j]
		}
		return ""
	}
	if j := strings.IndexAny(v, ", "); j >= 0 {
		v = v[:j]
	}
	return v
}

// wellKnown inserts /.well-known/<name> between host and path, as RFC
// 8414 §3.1 and RFC 9728 §3.1 say.
func wellKnown(id, name string) string {
	u, err := url.Parse(id)
	if err != nil {
		return id
	}
	p := strings.TrimSuffix(u.Path, "/")
	u.Path = "/.well-known/" + name + p
	u.RawQuery, u.Fragment = "", ""
	return u.String()
}

func trim(s string) string { return strings.TrimRight(s, "/") }

func (c *Client) checkURL(s string) error {
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return fmt.Errorf("%q is not a URL", s)
	}
	if u.Scheme == "https" || (c.AllowHTTP && u.Scheme == "http") {
		return nil
	}
	return fmt.Errorf("%q must use https", s)
}

func (c *Client) getJSON(ctx context.Context, u string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.hc().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: status %d", u, resp.StatusCode)
	}
	// What the server sent is never repeated back, only that it was not
	// metadata.
	if json.NewDecoder(io.LimitReader(resp.Body, 256<<10)).Decode(v) != nil {
		return fmt.Errorf("%s: not JSON metadata", u)
	}
	return nil
}

// Register registers planque as a public client with one exact redirect
// URI (RFC 7591) and returns its client_id.
func (c *Client) Register(ctx context.Context, m Metadata, name, redirect string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"client_name":                name,
		"redirect_uris":              []string{redirect},
		"grant_types":                []string{"authorization_code", "refresh_token"},
		"response_types":             []string{"code"},
		"token_endpoint_auth_method": "none",
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.Registration, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.hc().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r struct {
		ClientID     string   `json:"client_id"`
		RedirectURIs []string `json:"redirect_uris"`
		Error        string   `json:"error"`
		Description  string   `json:"error_description"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&r)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK || r.ClientID == "" {
		return "", fmt.Errorf("client registration failed (%d): %s", resp.StatusCode, short(r.Error))
	}
	// The server may not change where it sends the owner back.
	if len(r.RedirectURIs) > 0 && !slices.Equal(r.RedirectURIs, []string{redirect}) {
		return "", errors.New("the authorization server registered other redirect URIs than planque's")
	}
	return r.ClientID, nil
}

// PKCE is a code verifier and its S256 challenge.
type PKCE struct{ Verifier, Challenge string }

// NewPKCE makes a 43-character verifier (RFC 7636 §4.1).
func NewPKCE() PKCE {
	v := Random(32)
	sum := sha256.Sum256([]byte(v))
	return PKCE{Verifier: v, Challenge: base64.RawURLEncoding.EncodeToString(sum[:])}
}

// Random returns n random bytes, base64url.
func Random(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// AuthURL is the vendor's own sign-in page for this request.
func AuthURL(m Metadata, clientID, redirect, state, challenge string) string {
	q := url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"redirect_uri":          {redirect},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"resource":              {m.Resource},
	}
	if len(m.Scopes) > 0 {
		q.Set("scope", strings.Join(m.Scopes, " "))
	}
	sep := "?"
	if strings.Contains(m.Authorization, "?") {
		sep = "&"
	}
	return m.Authorization + sep + q.Encode()
}

// Token is a vendor's answer at its token endpoint.
type Token struct {
	Access  string    `json:"access_token"`
	Refresh string    `json:"refresh_token,omitempty"`
	Expiry  time.Time `json:"expiry,omitzero"`
	Scope   string    `json:"scope,omitempty"`
	// Account is the vendor account the token acts in, when the token
	// response names it (Stripe: stripe_user_id or account).
	Account string `json:"account,omitempty"`
}

// ErrInvalidGrant: the code or refresh token is no longer valid.
var ErrInvalidGrant = errors.New("oauth: invalid_grant")

// Exchange trades an authorization code for a token.
func (c *Client) Exchange(ctx context.Context, m Metadata, clientID, redirect, code, verifier string) (Token, error) {
	return c.token(ctx, m, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirect},
		"client_id":     {clientID},
		"code_verifier": {verifier},
		"resource":      {m.Resource},
	})
}

// accountPattern is what a vendor's account id may be: short, one word.
var accountPattern = regexp.MustCompile(`^[A-Za-z0-9_.:@+-]{1,64}$`)

// Refresh trades a refresh token for a new token.
func (c *Client) Refresh(ctx context.Context, m Metadata, clientID, refresh string) (Token, error) {
	return c.token(ctx, m, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh},
		"client_id":     {clientID},
		"resource":      {m.Resource},
	})
}

func (c *Client) token(ctx context.Context, m Metadata, form url.Values) (Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.Token, strings.NewReader(form.Encode()))
	if err != nil {
		return Token{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := c.hc().Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()
	var r struct {
		Access      string `json:"access_token"`
		Type        string `json:"token_type"`
		Refresh     string `json:"refresh_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Scope       string `json:"scope"`
		StripeUser  string `json:"stripe_user_id"`
		Account     string `json:"account"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&r)
	if r.Error == "invalid_grant" {
		return Token{}, ErrInvalidGrant
	}
	if resp.StatusCode != http.StatusOK || r.Access == "" {
		return Token{}, fmt.Errorf("token endpoint (%d): %s", resp.StatusCode, short(r.Error))
	}
	if r.Type != "" && !strings.EqualFold(r.Type, "bearer") {
		return Token{}, fmt.Errorf("token type %q is not bearer", r.Type)
	}
	t := Token{Access: r.Access, Refresh: r.Refresh, Scope: r.Scope, Account: r.StripeUser}
	if t.Account == "" {
		t.Account = r.Account
	}
	// The owner reads the account in an approval's text: one that doesn't look
	// like an account (a newline, a sentence) is left out, as when none is named.
	if !accountPattern.MatchString(t.Account) {
		t.Account = ""
	}
	if r.ExpiresIn > 0 {
		t.Expiry = time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)
	}
	return t, nil
}

// Revoke revokes token at the vendor (RFC 7009), when it offers that.
func (c *Client) Revoke(ctx context.Context, m Metadata, clientID, token, hint string) error {
	if m.Revocation == "" {
		return nil
	}
	form := url.Values{"token": {token}, "token_type_hint": {hint}, "client_id": {clientID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.Revocation, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.hc().Do(req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("revocation endpoint: status %d", resp.StatusCode)
	}
	return nil
}

// short keeps an OAuth error code readable and bounded; descriptions are
// not repeated.
func short(code string) string {
	if len(code) > 64 {
		code = code[:64]
	}
	return code
}
