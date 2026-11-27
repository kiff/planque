package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRegisteredAppCredentialsAreBoundToClientAndEndpoint(t *testing.T) {
	for _, method := range []string{"client_secret_post", "client_secret_basic", "none"} {
		t.Run(method, func(t *testing.T) {
			app := Registration{ClientID: "client id", Secret: "test-secret+value", TokenURL: "https://auth.example/token", AuthMethod: method}
			if method == "none" {
				app.Secret = ""
			}
			calls := 0
			c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				switch method {
				case "client_secret_post":
					if r.PostForm.Get("client_secret") != app.Secret || r.PostForm.Get("client_id") != app.ClientID || r.Header.Get("Authorization") != "" {
						t.Fatal("registered app did not use its configured body authentication")
					}
				case "client_secret_basic":
					id, secret, ok := r.BasicAuth()
					if !ok || id != url.QueryEscape(app.ClientID) || secret != url.QueryEscape(app.Secret) || r.PostForm.Get("client_secret") != "" || r.PostForm.Get("client_id") != "" {
						t.Fatal("registered app used two authentication methods or incorrect Basic encoding")
					}
				case "none":
					if r.PostForm.Get("client_id") != app.ClientID || r.PostForm.Get("client_secret") != "" || r.Header.Get("Authorization") != "" {
						t.Fatal("public client leaked a credential")
					}
				}
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"access_token":"new-access","refresh_token":"rotated-refresh","token_type":"Bearer"}`))}, nil
			})}}
			m := Metadata{Token: app.TokenURL, Resource: "https://api.example"}
			if _, err := c.Exchange(context.Background(), m, app.ClientID, "https://planque.example/oauth/callback", "code", "verifier", app); err != nil {
				t.Fatal(err)
			}
			tok, err := c.Refresh(context.Background(), m, app.ClientID, "old-refresh", app)
			if err != nil || tok.Refresh != "rotated-refresh" || calls != 2 {
				t.Fatalf("refresh token was not rotated: %v", err)
			}
			for _, bad := range []Metadata{{Token: "https://elsewhere.example/token"}, {Token: app.TokenURL + "?redirect=x"}} {
				if _, err := c.Refresh(context.Background(), bad, app.ClientID, "old-refresh", app); err == nil {
					t.Fatal("an altered endpoint could receive app credentials")
				}
			}
			if _, err := c.Refresh(context.Background(), m, "another-client", "old-refresh", app); err == nil || calls != 2 {
				t.Fatal("another client could use this registration or a refused request was sent")
			}
			b, _ := json.Marshal(app)
			if strings.Contains(string(b), "test-secret") || strings.Contains(string(b), "client_secret\"") {
				t.Fatal("serializing a registration exposed the app secret")
			}
		})
	}
}

func TestTokenRequestsDoNotReplayCredentialsToRedirects(t *testing.T) {
	calls := 0
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Host != "auth.example" {
			t.Fatal("OAuth credentials reached a redirect target")
		}
		return &http.Response{StatusCode: 307, Header: http.Header{"Location": {"https://elsewhere.example/token"}}, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}}
	app := Registration{ClientID: "client", Secret: "secret", TokenURL: "https://auth.example/token", AuthMethod: "client_secret_post"}
	if _, err := c.Refresh(context.Background(), Metadata{Token: app.TokenURL}, app.ClientID, "refresh", app); err == nil || calls != 1 {
		t.Fatal("a redirected token request was accepted or replayed")
	}
}

func TestMalformedTokenResponsesAreNotAccepted(t *testing.T) {
	for _, body := range []string{`{"access_token":"partial",`, `{"access_token":"access","error":"invalid_client"}`} {
		c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
		})}}
		if _, err := c.Refresh(context.Background(), Metadata{Token: "https://auth.example/token"}, "client", "refresh"); err == nil {
			t.Fatal("an incomplete or error-bearing token response was accepted")
		}
	}
}

// A server whose token answer names the account its own way (Stripe's
// stripe_user_id) gives it only where planque set Metadata.AccountField;
// "account" is read either way, and a value that isn't an account is left out.
func TestTheAccountFieldIsReadOnlyWhereSet(t *testing.T) {
	answer := ""
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(answer))}, nil
	})}}
	exchange := func(field string) string {
		t.Helper()
		tok, err := c.Exchange(context.Background(), Metadata{Token: "https://auth.example/token", Resource: "https://mcp.example", AccountField: field}, "client", "https://planque.example/oauth/callback", "code", "verifier")
		if err != nil {
			t.Fatal(err)
		}
		return tok.Account
	}
	answer = `{"access_token":"a","token_type":"bearer","stripe_user_id":"acct_1"}`
	if got := exchange("stripe_user_id"); got != "acct_1" {
		t.Fatalf("with the field set: %q", got)
	}
	if got := exchange(""); got != "" {
		t.Fatalf("a vendor's own field read without being set: %q", got)
	}
	answer = `{"access_token":"a","token_type":"bearer","account":"acct_2","stripe_user_id":"acct_1"}`
	if got := exchange("stripe_user_id"); got != "acct_2" {
		t.Fatalf("account comes first: %q", got)
	}
	for _, bad := range []string{`"two words"`, `7`, `{"id":"acct_1"}`} {
		answer = `{"access_token":"a","token_type":"bearer","stripe_user_id":` + bad + `}`
		if got := exchange("stripe_user_id"); got != "" {
			t.Fatalf("%s read as %q", bad, got)
		}
	}
}

// A vendor's gate that wants a header on every call to its token endpoint (Qonto's sandbox) gets the registered
// app's header on the exchange, the refresh and the revocation, and a registration without one sends none.
func TestARegisteredSecretHeaderGoesToTheRegisteredEndpointsOnly(t *testing.T) {
	app := Registration{ClientID: "c", Secret: "s", TokenURL: "https://auth.example/token", RevocationURL: "https://auth.example/revoke",
		AuthMethod: "client_secret_post", Header: "X-Gate-Token", HeaderValue: "gate-value"}
	seen := map[string]string{}
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen[r.URL.Path] = r.Header.Get("X-Gate-Token")
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"access_token":"a","refresh_token":"r","token_type":"Bearer"}`))}, nil
	})}}
	m := Metadata{Token: app.TokenURL, Revocation: app.RevocationURL}
	ctx := context.Background()
	if _, err := c.Exchange(ctx, m, "c", "https://planque.example/oauth/callback", "code", "verifier", app); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Refresh(ctx, m, "c", "old", app); err != nil {
		t.Fatal(err)
	}
	if err := c.Revoke(ctx, m, "c", "old", "refresh_token", app); err != nil {
		t.Fatal(err)
	}
	if seen["/token"] != "gate-value" || seen["/revoke"] != "gate-value" {
		t.Fatalf("the registered header didn't reach the registered endpoints: %v", seen)
	}
	// A name with no value is a configuration fault: nothing is sent.
	app.HeaderValue = ""
	before := len(seen)
	seen = map[string]string{}
	if _, err := c.Refresh(ctx, m, "c", "old", app); err == nil || len(seen) != 0 || before == 0 {
		t.Fatalf("a request left without its header's value: %v %v", err, seen)
	}
	// No header configured: none is sent.
	app.Header = ""
	if _, err := c.Refresh(ctx, m, "c", "old", app); err != nil || seen["/token"] != "" {
		t.Fatalf("a header was sent that the registration doesn't name: %v %v", err, seen)
	}
	if b, _ := json.Marshal(Registration{Header: "X-Gate-Token", HeaderValue: "gate-value"}); strings.Contains(string(b), "gate-value") {
		t.Fatalf("the header's value is in the registration's JSON: %s", b)
	}
}
