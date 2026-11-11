package oauth

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDiscoverClientMetadataWithoutRegistration(t *testing.T) {
	requests := 0
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		status, body := 200, ""
		switch r.URL.String() {
		case "https://tools.example/mcp":
			status = 401
		case "https://tools.example/.well-known/oauth-protected-resource/mcp":
			body = `{"resource":"https://tools.example/mcp","authorization_servers":["https://auth.example"]}`
		case "https://auth.example/.well-known/oauth-authorization-server":
			body = `{"issuer":"https://auth.example","authorization_endpoint":"https://auth.example/authorize","token_endpoint":"https://auth.example/token","code_challenge_methods_supported":["S256"],"client_id_metadata_document_supported":true}`
		default:
			t.Fatalf("unexpected discovery or registration request: %s", r.URL)
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	m, err := c.Discover(context.Background(), "https://tools.example/mcp")
	if err != nil || !m.ClientMetadata || m.Registration != "" || requests != 3 {
		t.Fatalf("metadata=%+v requests=%d err=%v", m, requests, err)
	}
	q := AuthURL(m, "https://planque.example/oauth/client-metadata.json", "https://planque.example/oauth/callback", "state", "challenge")
	if !strings.Contains(q, "client_id=https%3A%2F%2Fplanque.example%2Foauth%2Fclient-metadata.json") {
		t.Fatal("authorization did not use the metadata URL as its client id")
	}
}

func TestClientDocumentUsesOnlyFixedHTTPSURLs(t *testing.T) {
	doc, err := NewClientDocument("https://planque.example/oauth/client-metadata.json", "tafs", "https://planque.example/oauth/callback")
	if err != nil || doc.ClientID != "https://planque.example/oauth/client-metadata.json" || len(doc.RedirectURIs) != 1 || doc.AuthMethod != "none" {
		t.Fatalf("%+v %v", doc, err)
	}
	for _, id := range []string{"http://planque.example/client.json", "https://planque.example", "https://user@planque.example/client.json", "https://planque.example/client.json?tenant=x", "https://planque.example/client.json#x"} {
		if _, err := NewClientDocument(id, "tafs", "https://planque.example/oauth/callback"); err == nil {
			t.Fatalf("accepted invalid client document URL %q", id)
		}
	}
}

// A server's metadata may name its origin for a server at a path, as Qonto
// does for https://mcp.qonto.com/mcp: the token is bound to what it names.
// Another host, or a sibling path, is someone else's resource.
func TestDiscoverAcceptsTheServersOriginAndRefusesAnotherResource(t *testing.T) {
	discover := func(resource string) (Metadata, error) {
		c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			status, body, h := 200, "", make(http.Header)
			switch r.URL.String() {
			case "https://bank.example/mcp":
				status = 401
				h.Set("WWW-Authenticate", `Bearer resource_metadata="https://bank.example/.well-known/oauth-protected-resource"`)
			case "https://bank.example/.well-known/oauth-protected-resource":
				body = `{"resource":"` + resource + `","authorization_servers":["https://bank.example"],"scopes_supported":["client.read","request_transfers.write"]}`
			case "https://bank.example/.well-known/oauth-authorization-server":
				body = `{"issuer":"https://bank.example","authorization_endpoint":"https://bank.example/authorize","token_endpoint":"https://bank.example/token","registration_endpoint":"https://bank.example/register","code_challenge_methods_supported":["S256"],"token_endpoint_auth_methods_supported":["none"]}`
			default:
				t.Fatalf("unexpected request: %s", r.URL)
			}
			return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(body))}, nil
		})}}
		return c.Discover(context.Background(), "https://bank.example/mcp")
	}
	for _, resource := range []string{"https://bank.example", "https://bank.example/", "https://bank.example/mcp", "https://BANK.example"} {
		m, err := discover(resource)
		if err != nil || m.Resource != resource {
			t.Fatalf("%s: %+v %v", resource, m, err)
		}
		if q := AuthURL(m, "client", "https://planque.example/oauth/callback", "state", "challenge"); !strings.Contains(q, "resource="+url.QueryEscape(resource)) {
			t.Fatalf("%s: the sign-in isn't bound to the resource named: %s", resource, q)
		}
	}
	for _, resource := range []string{"https://other.example", "https://bank.example.evil", "https://bank.example/mcpx", "https://bank.example/mcp/tools",
		"http://bank.example", "https://bank.example:8443", "https://bank.example/?x=1", "https://user@bank.example", ""} {
		if m, err := discover(resource); err == nil {
			t.Fatalf("%q accepted as the resource for /mcp: %+v", resource, m)
		}
	}
}
