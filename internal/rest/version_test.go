package rest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

type versionTransport func(*http.Request) (*http.Response, error)

func (f versionTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestReviewedVersionHeaderDoesNotReplaceAuthentication(t *testing.T) {
	calls := 0
	c := &Client{HTTP: &http.Client{Transport: versionTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Pay-Version") != "2030-01-01" || r.Header.Get("Authorization") != "Bearer synthetic-key" {
			t.Errorf("incorrect version or authentication header")
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}}
	auth := openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer "}
	pin := &VersionPin{Name: "Pay-Version", Value: "2030-01-01"}
	if _, err := c.Call(context.Background(), "POST", "https://vendor.example/v1", []byte(`{}`), auth, "synthetic-key", pin); err != nil {
		t.Fatal(err)
	}
	for _, pin := range []*VersionPin{{Name: "Authorization", Value: "replace"}, {Name: "Pay-Version", Value: "v1\r\nAuthorization: replace"}, {Name: "Pay-Version", Value: ""}} {
		if _, err := c.Call(context.Background(), "POST", "https://vendor.example/v1", nil, auth, "synthetic-key", pin); err == nil {
			t.Fatal("invalid pin reached the transport")
		}
	}
	auth.Name = "Pay-Version"
	if _, err := c.Call(context.Background(), "POST", "https://vendor.example/v1", nil, auth, "synthetic-key", pin); err == nil {
		t.Fatal("the version replaced authentication")
	}
	if calls != 1 {
		t.Fatalf("invalid pin made %d transport calls", calls-1)
	}
}

// A secret header goes to the host it was made for, never replaces the key's header, and a redirect (Qonto's
// gate answers one) is the answer: the header never follows it to another host.
func TestASecretHeaderStaysWithItsHost(t *testing.T) {
	var hosts []string
	c := &Client{HTTP: &http.Client{Transport: versionTransport(func(r *http.Request) (*http.Response, error) {
		hosts = append(hosts, r.URL.Host+" "+r.Header.Get("X-Gate-Token"))
		h := make(http.Header)
		h.Set("Location", "https://portal.example/login")
		return &http.Response{StatusCode: 302, Header: h, Body: io.NopCloser(strings.NewReader(``))}, nil
	})}}
	auth := openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer "}
	ctx := WithSecretHeader(context.Background(), SecretHeader{Name: "X-Gate-Token", Value: "gate", Host: "api.example"})
	if _, err := c.Call(ctx, "GET", "https://api.example/v2/x", nil, auth, "key"); err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 || hosts[0] != "api.example gate" {
		t.Fatalf("the call or a redirect went wrong: %v", hosts)
	}
	hosts = nil
	if _, err := c.Call(ctx, "GET", "https://other.example/v2/x", nil, auth, "key"); err == nil || len(hosts) != 0 {
		t.Fatalf("the secret header was sent to another host: %v %v", err, hosts)
	}
	bad := WithSecretHeader(context.Background(), SecretHeader{Name: "Authorization", Value: "gate", Host: "api.example"})
	if _, err := c.Call(bad, "GET", "https://api.example/v2/x", nil, auth, "key"); err == nil || len(hosts) != 0 {
		t.Fatalf("the secret header replaced authentication: %v %v", err, hosts)
	}
	none := WithSecretHeader(context.Background(), SecretHeader{Name: "X-Gate-Token", Host: "api.example"})
	if _, err := c.Call(none, "GET", "https://api.example/v2/x", nil, auth, "key"); err == nil || len(hosts) != 0 {
		t.Fatalf("a header without a value was sent: %v %v", err, hosts)
	}
}
