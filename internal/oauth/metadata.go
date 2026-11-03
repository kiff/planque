package oauth

import (
	"errors"
	"net/url"
	"strings"
)

// ClientDocument is the public identity of this installation. It contains
// no grant, account or credential data. Its URL is also its OAuth client id.
// https://modelcontextprotocol.io/specification/2025-11-25/basic/authorization
type ClientDocument struct {
	ClientID      string   `json:"client_id"`
	Name          string   `json:"client_name"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	AuthMethod    string   `json:"token_endpoint_auth_method"`
}

// NewClientDocument fixes the identifier and redirect to the installation's
// configured URLs, never the inbound Host header or a requesting agent.
func NewClientDocument(id, name, redirect string) (ClientDocument, error) {
	if !publicHTTPS(id, true) || !publicHTTPS(redirect, false) || strings.TrimSpace(name) == "" {
		return ClientDocument{}, errors.New("client metadata requires an HTTPS document URL with a path, an HTTPS redirect, and a client name")
	}
	return ClientDocument{ClientID: id, Name: name, RedirectURIs: []string{redirect},
		GrantTypes:    []string{"authorization_code", "refresh_token"},
		ResponseTypes: []string{"code"}, AuthMethod: "none"}, nil
}

func publicHTTPS(raw string, path bool) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil &&
		u.RawQuery == "" && u.Fragment == "" && (!path || u.Path != "" && u.Path != "/")
}
