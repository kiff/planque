package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// signInReader reads whether the managed sign-in works for an app from planque's own route
// (GET /v1/catalog/{entry}/sign-in), with the one key that reads it and nothing else: it holds no provider key,
// no AWS credential, and reads no provider directly. The run gives it CONNECT_MATRIX_PLANQUE_URL (planque's
// address) and CONNECT_MATRIX_SIGNIN_KEY (the reader key); without both, the column says it wasn't checked.
type signInReader struct {
	base string
	key  string
	http *http.Client
}

func managedReader() *signInReader {
	base := strings.TrimRight(os.Getenv("CONNECT_MATRIX_PLANQUE_URL"), "/")
	key := os.Getenv("CONNECT_MATRIX_SIGNIN_KEY")
	if base == "" || key == "" || !strings.HasPrefix(base, "https://") {
		return nil
	}
	return &signInReader{base: base, key: key, http: &http.Client{Timeout: 30 * time.Second,
		// Never follow a redirect: the key goes to planque's address and nowhere else.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}

// managedProvider is the catalog entry's provider part for the managed sign-in's apps.
const managedProvider = "composio"

// managed says whether the managed sign-in can sign someone in to the app, as planque answers it: its state only,
// never a session, an account or a call.
func managed(p *signInReader, slug string) string {
	if slug == "" {
		return "-"
	}
	if p == nil {
		return "not checked (no sign-in reader key in this run)"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	entry := "connector:" + managedProvider + ":" + slug
	req, err := http.NewRequestWithContext(ctx, "GET", p.base+"/v1/catalog/"+url.PathEscape(entry)+"/sign-in", nil)
	if err != nil {
		return "? couldn't be read"
	}
	req.Header.Set("Authorization", "Bearer "+p.key)
	res, err := p.http.Do(req)
	if err != nil {
		return "? couldn't be read"
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusUnauthorized:
		return "? the sign-in reader key was refused"
	case http.StatusNotFound:
		// The route is off (no reader key set in planque), or the entry isn't a managed sign-in app: either way, unread.
		return "? not answered (route off, or not a managed sign-in app)"
	}
	var out struct {
		State string `json:"state"`
	}
	if res.StatusCode != http.StatusOK || json.NewDecoder(io.LimitReader(res.Body, 4096)).Decode(&out) != nil {
		return "? couldn't be read"
	}
	switch out.State {
	case "ready":
		return "✅ ready"
	case "on_the_owners_pick":
		return "✅ on the owner's pick"
	case "needs_setup":
		return "❌ needs an operator's own credentials"
	case "not_at_provider":
		return "❌ not at the provider"
	}
	return "? couldn't be read"
}
