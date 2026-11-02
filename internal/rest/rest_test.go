package rest

import (
	"context"
	"errors"
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

func TestADescriptionOver2MBIsRefused(t *testing.T) {
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
