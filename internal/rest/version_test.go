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
