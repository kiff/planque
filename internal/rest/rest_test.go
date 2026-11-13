package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

func TestACallSendsTheKeyWhereItGoesAndFollowsNoRedirect(t *testing.T) {
	var got []string
	away := 0
	elsewhere := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { away++ }))
	defer elsewhere.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+r.Header.Get("Authorization"))
		if r.URL.Path == "/away" {
			http.Redirect(w, r, elsewhere.URL, http.StatusFound)
			return
		}
		_, _ = w.Write([]byte(strings.Repeat("x", MaxAnswer+10)))
	}))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	ctx := context.Background()
	a, err := c.Call(ctx, "GET", srv.URL+"/a?x=1", nil, openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer "}, "k1")
	if err != nil || a.Status != 200 || !a.Cut || len(a.Body) != MaxAnswer || got[0] != "GET /a?x=1 Bearer k1" {
		t.Fatalf("%+v %v %v", a.Status, err, got)
	}
	if _, err := c.Call(ctx, "GET", srv.URL+"/q", nil, openapi.Auth{In: "query", Name: "api_key"}, "k2"); err != nil || got[1] != "GET /q?api_key=k2 " {
		t.Fatalf("%v %v", err, got)
	}
	a, err = c.Call(ctx, "POST", srv.URL+"/away", []byte(`{}`), openapi.Auth{In: "header", Name: "X-Key"}, "k3")
	if err != nil || a.Status != http.StatusFound || away != 0 {
		t.Fatalf("a redirect: %d %v, followed %d", a.Status, err, away)
	}
	// Only https, unless a test allows http.
	if _, err := (&Client{}).Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "header", Name: "X"}, "k"); err == nil {
		t.Fatal("http without AllowHTTP")
	}
	var ce *CallError
	if _, err := c.Call(ctx, "GET", "http://127.0.0.1:1/x", nil, openapi.Auth{In: "header", Name: "X"}, "k"); !errors.As(err, &ce) || ce.Sent {
		t.Fatalf("a refused connection is not sent: %v", err)
	}
}

func TestADescriptionOver16MBIsRefused(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("{" + strings.Repeat(" ", openapi.MaxBytes) + "}"))
	}))
	defer srv.Close()
	var r *openapi.Refused
	if _, err := (&Client{AllowHTTP: true}).Fetch(context.Background(), srv.URL); !errors.As(err, &r) {
		t.Fatal(err)
	}
}

// A failed call's error is logged: it never holds the key, though Go's
// url.Error carries the URL a query key is in.
func TestAFailedCallsErrorHoldsNoKey(t *testing.T) {
	c := &Client{AllowHTTP: true}
	key := "sk_live/secret+1"
	for _, target := range []string{"http://127.0.0.1:1/x?y=1", "http://nohost.invalid/x"} {
		_, err := c.Call(context.Background(), "GET", target, nil, openapi.Auth{In: "query", Name: "api_key"}, key)
		if err == nil || strings.Contains(err.Error(), "secret") {
			t.Fatalf("%s: %v", target, err)
		}
	}
	if got := HideKey(`{"error":"invalid key sk_live/secret+1"}`, key); got != `{"error":"invalid key [key]"}` {
		t.Fatal(got)
	}
}

// HTTP Basic: the key as the password with the username given, or as the
// username with the password given; a key with a colon can't be a username,
// and a version pin can't replace the Authorization it fills.
func TestABasicKeyIsSentAsHTTPBasic(t *testing.T) {
	var user, pass string
	var ok bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok = r.BasicAuth()
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	ctx := context.Background()
	if _, err := c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "basic", User: "me@example.com/token"}, "k3y"); err != nil || !ok || user != "me@example.com/token" || pass != "k3y" {
		t.Fatalf("the key as the password: %v %v %q %q", err, ok, user, pass)
	}
	if _, err := c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "basic", Password: "X"}, "k3y"); err != nil || user != "k3y" || pass != "X" {
		t.Fatalf("the key as the username: %v %q %q", err, user, pass)
	}
	user = ""
	if _, err := c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "basic"}, "k:3y"); err == nil || user != "" || strings.Contains(err.Error(), "k:3y") {
		t.Fatalf("a key with a colon as the username: %v (sent as %q)", err, user)
	}
	if _, err := c.Call(ctx, "GET", srv.URL, nil, openapi.Auth{In: "basic", User: "me"}, "k3y", &VersionPin{Name: "Authorization", Value: "2024-01-01"}); err == nil {
		t.Fatal("a version pin replaced Basic's Authorization")
	}
}

// A form body goes as a form.
func TestAFormBodyGoesAsAForm(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Content-Type")
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	if _, err := c.CallAs(context.Background(), "POST", srv.URL, "application/x-www-form-urlencoded", []byte("a=1"), openapi.Auth{In: "header", Name: "X"}, "k"); err != nil || got != "application/x-www-form-urlencoded" {
		t.Fatalf("%v %q", err, got)
	}
	if _, err := c.Call(context.Background(), "POST", srv.URL, []byte("{}"), openapi.Auth{In: "header", Name: "X"}, "k"); err != nil || got != "application/json" {
		t.Fatalf("%v %q", err, got)
	}
}

// The fixed header goes on every call, and a version pin can't replace it.
func TestAFixedHeaderGoesOnEveryCall(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Account-Id")
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	auth := openapi.Auth{In: "header", Name: "Authorization", Prefix: "Bearer ", Fixed: &openapi.Fixed{Name: "Account-Id", Value: "123456"}}
	if _, err := c.Call(context.Background(), "GET", srv.URL, nil, auth, "k"); err != nil || got != "123456" {
		t.Fatalf("%v %q", err, got)
	}
	if _, err := c.Call(context.Background(), "GET", srv.URL, nil, auth, "k", &VersionPin{Name: "account-id", Value: "2024-01-01"}); err == nil {
		t.Fatal("a version pin replaced the fixed header")
	}
}

// The fixed header never carries the key: its approval says where the key goes.
func TestAFixedHeaderNeverCarriesTheKey(t *testing.T) {
	sent := false
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { sent = true }))
	defer srv.Close()
	c := &Client{AllowHTTP: true}
	auth := openapi.Auth{In: "header", Name: "X-API-Key", Fixed: &openapi.Fixed{Name: "Account-Id", Value: "acct-sk_live_1"}}
	if _, err := c.Call(context.Background(), "GET", srv.URL, nil, auth, "sk_live_1"); err == nil || sent || strings.Contains(err.Error(), "sk_live_1") {
		t.Fatalf("%v, sent %v", err, sent)
	}
}
