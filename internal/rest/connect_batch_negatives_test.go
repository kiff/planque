package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// The negatives for the connect batch at the client (qa-kiff): whatever reached it, a fixed header never replaces
// the key's Authorization or the Host, and never carries the key itself; a Basic key's failed call never says it.
func TestAFixedHeaderNeverReplacesTheKeyOrCarriesIt(t *testing.T) {
	var got http.Header
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, host = r.Header.Clone(), r.Host
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	ctx := context.Background()
	// Past Auth.Check (a grant stored before it, or a bug): the key's own header still wins, or nothing is sent.
	for _, a := range []openapi.Auth{
		{In: "header", Name: "Authorization", Prefix: "Bearer ", Fixed: &openapi.Fixed{Name: "Authorization", Value: "Bearer someone-else"}},
		{In: "basic", User: "me", Fixed: &openapi.Fixed{Name: "Authorization", Value: "Basic c29tZW9uZS1lbHNl"}},
	} {
		got = nil
		_, err := c.Call(ctx, "GET", srv.URL, nil, a, "k3y")
		if err == nil && got != nil && strings.Contains(got.Get("Authorization"), "someone-else") || got != nil && strings.Contains(got.Get("Authorization"), "c29tZW9uZS1lbHNl") {
			t.Errorf("a fixed header replaced the key's Authorization (%s): %q", a.In, got.Get("Authorization"))
		}
	}
	got = nil
	_, _ = c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "header", Name: "X-API-Key", Fixed: &openapi.Fixed{Name: "Host", Value: "evil.example"}}, "k3y")
	if got != nil && host == "evil.example" {
		t.Error("a fixed header moved the Host")
	}
	// The key never goes in the fixed header, which the owner reads in the approval and which isn't kept as a secret.
	got = nil
	_, err := c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "header", Name: "X-API-Key", Fixed: &openapi.Fixed{Name: "Harvest-Account-Id", Value: "k3y"}}, "k3y")
	if err == nil && got != nil && got.Get("Harvest-Account-Id") == "k3y" {
		t.Error("the fixed header repeats the key")
	}
}

func TestABasicKeysFailedCallNeverSaysIt(t *testing.T) {
	c := &Client{AllowHTTP: true}
	key := "sk_live/s3cret+1"
	for _, a := range []openapi.Auth{{In: "basic", User: "me@example.com/token"}, {In: "basic", Password: "X"}} {
		for _, target := range []string{"http://127.0.0.1:1/x", "http://nohost.invalid/x"} {
			_, err := c.Call(context.Background(), "GET", target, nil, a, key)
			if err == nil || strings.Contains(err.Error(), "s3cret") {
				t.Errorf("%s %+v: %v", target, a, err)
			}
		}
	}
}
