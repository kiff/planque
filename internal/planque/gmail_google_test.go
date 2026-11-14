package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/profile"
)

// shippedApps is deploy/apps.json, as planque reads it at start.
func shippedApps(t *testing.T) map[string]appconfig.Application {
	t.Helper()
	raw, err := os.ReadFile("../../deploy/apps.json")
	if err != nil {
		t.Fatal(err)
	}
	apps, err := appconfig.Decode(raw)
	if err != nil {
		t.Fatalf("deploy/apps.json would refuse planque's start: %v", err)
	}
	return apps
}

const gmailEntry = "gmail-google@1"

var gmailScopes = []string{"https://www.googleapis.com/auth/gmail.readonly", "https://www.googleapis.com/auth/gmail.compose"}

// Gmail by tafs's own Google client: Google's endpoints, the two Gmail
// scopes alone, a refresh token asked for, PKCE kept, and the reviewed
// operations under /gmail/v1/users/me, of which only the drafts write.
func TestShippedGmailAppAndItsCatalogEntry(t *testing.T) {
	apps := shippedApps(t)
	app, ok := apps["gmail"]
	if !ok || app.Kind != appconfig.REST || app.Server != "https://gmail.googleapis.com" || app.AuthMethod != "client_secret_post" {
		t.Fatalf("the Gmail app: %+v", app)
	}
	o := app.OAuth
	if o.Issuer != "https://accounts.google.com" || o.Authorization != "https://accounts.google.com/o/oauth2/v2/auth" || o.Token != "https://oauth2.googleapis.com/token" ||
		o.Revocation != "https://oauth2.googleapis.com/revoke" || !slices.Equal(o.Scopes, gmailScopes) || !o.NoResourceParameter || o.NoPKCE ||
		len(o.AuthorizationParams) != 2 || o.AuthorizationParams["access_type"] != "offline" || o.AuthorizationParams["prompt"] != "consent" {
		t.Fatalf("Gmail's OAuth: %+v", o)
	}
	var names []string
	for _, op := range app.Snapshot.Operations {
		names = append(names, op.Name)
		if !strings.HasPrefix(op.Path, "/gmail/v1/users/me/") || op.ReadOnly != (op.Method == "GET") {
			t.Fatalf("operation %s: %s %s", op.Name, op.Method, op.Path)
		}
	}
	if strings.Join(names, ",") != "listMessages,getMessage,getThread,listLabels,listDrafts,createDraft,updateDraft,sendDraft" {
		t.Fatalf("operations: %v", names)
	}
	conf, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := DecodeCatalog(conf, apps)
	if err != nil {
		t.Fatal(err)
	}
	entry := entries[gmailEntry]
	if entry.Application != "gmail" || entry.Kind != "" || entry.Server != "" || len(entry.Tools) != len(names) {
		t.Fatalf("the Gmail entry doesn't resolve to the app: %+v", entry)
	}
	if apiToolsCheck(app.Snapshot, entry.Tools) != nil {
		t.Fatal("the entry's tools aren't the app's operations")
	}
	// Without the app, the entry is refused rather than served half-made.
	if _, err := DecodeCatalog(conf, nil); err == nil {
		t.Fatal("the Gmail entry was accepted without its registered app")
	}
}

// The whole connection against a fake Google: approve, sign in, exchange
// the code with the registered client, list, draft, refresh, revoke. The
// client secret is never in a grant view, the export or a log line.
func TestGmailByTafsOwnGoogleClient(t *testing.T) {
	const secret = "gmail-client-secret-value"
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	// KIFF applies the reviewed profile the phone approved.
	base := e.client.Transport
	e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
			return base.RoundTrip(r)
		}
		var in struct{ Profile json.RawMessage }
		_ = json.NewDecoder(r.Body).Decode(&in)
		hash, _ := profile.Digest(in.Profile)
		return profileResponse(200, `{"revision":1,"profile_sha256":"`+hash+`","words_sha256":"`+strings.Repeat("a", 64)+`"}`), nil
	})
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	e.s.Applications = shippedApps(t)
	conf, _ := os.ReadFile("../../deploy/catalog.json")
	var err error
	if e.s.CatalogEntries, err = DecodeCatalog(conf, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "gmail", appcredentials.Credentials{ClientID: "tafs-client.apps.example", Secret: secret}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials

	var mu sync.Mutex
	var seen []string
	form := func(r *http.Request) url.Values {
		if err := r.ParseForm(); err != nil || r.PostForm.Get("client_id") != "tafs-client.apps.example" || r.PostForm.Get("client_secret") != secret || r.PostForm.Has("resource") {
			t.Errorf("%s wasn't sent by the registered client alone: %v", r.URL.Path, r.PostForm)
		}
		return r.PostForm
	}
	transport := applicationTransport(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		body := `{}`
		switch r.Method + " " + r.URL.Host + r.URL.Path {
		case "POST oauth2.googleapis.com/token":
			f := form(r)
			switch f.Get("grant_type") {
			case "authorization_code":
				if f.Get("code") != "google-code" || f.Get("code_verifier") == "" || f.Get("redirect_uri") != e.url+"/oauth/callback" {
					t.Errorf("code exchange: %v", f)
				}
				body = `{"access_token":"ya29.first","refresh_token":"1//refresh","token_type":"Bearer","expires_in":3599,"scope":"` + strings.Join(gmailScopes, " ") + `"}`
			case "refresh_token":
				if f.Get("refresh_token") != "1//refresh" {
					t.Errorf("refresh: %v", f)
				}
				body = `{"access_token":"ya29.second","token_type":"Bearer","expires_in":3599}`
			}
			seen = append(seen, "token:"+f.Get("grant_type"))
		case "POST oauth2.googleapis.com/revoke":
			form(r)
			seen = append(seen, "revoke")
		case "GET gmail.googleapis.com/gmail/v1/users/me/messages":
			if r.Header.Get("Authorization") != "Bearer ya29.first" || r.URL.Query().Get("q") != "from:ana@example.com" || r.URL.Query().Get("maxResults") != "5" {
				t.Errorf("list: %s %v", r.URL, r.Header.Get("Authorization"))
			}
			body = `{"messages":[{"id":"m1","threadId":"t1"}],"resultSizeEstimate":1}`
			seen = append(seen, "list")
		case "POST gmail.googleapis.com/gmail/v1/users/me/drafts":
			raw, _ := io.ReadAll(r.Body)
			if string(raw) != `{"message":{"raw":"VG86IGFuYUBleGFtcGxlLmNvbQ","threadId":"t1"}}` || r.Header.Get("Authorization") != "Bearer ya29.first" {
				t.Errorf("draft create: %s", raw)
			}
			body = `{"id":"r-1","message":{"id":"m2","threadId":"t1"}}`
			seen = append(seen, "draft")
		default:
			t.Errorf("unexpected request to Google: %s %s", r.Method, r.URL)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	e.s.OAuth.HTTP = &http.Client{Transport: transport}
	e.s.REST.HTTP = &http.Client{Transport: transport}

	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/"+gmailEntry+"/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 || out["needs_key"] != false {
		t.Fatalf("Gmail grant: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	g, _ := e.s.loadGrant(ctx, tenant, id)
	path := "/v1/accounts/" + tenant + "/grants/" + id
	st, out = e.do("POST", path+"/approve", reapprove(g), e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, out)
	}
	auth, _ := url.Parse(out["sign_in_url"].(string))
	q := auth.Query()
	if auth.Host != "accounts.google.com" || auth.Path != "/o/oauth2/v2/auth" || q.Get("client_id") != "tafs-client.apps.example" ||
		q.Get("access_type") != "offline" || q.Get("prompt") != "consent" || q.Get("code_challenge") == "" || q.Get("code_challenge_method") != "S256" ||
		q.Has("resource") || q.Get("scope") != strings.Join(gmailScopes, " ") || q.Get("redirect_uri") != e.url+"/oauth/callback" {
		t.Fatalf("Google's sign-in: %s", auth)
	}
	code := e.signIn(e.url + "/oauth/callback?state=" + url.QueryEscape(q.Get("state")) + "&code=google-code")
	if st, out = e.do("POST", path+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}

	for name := range map[string]bool{"gmail_list_messages": true, "gmail_create_draft": true} {
		e.kiff.Cards["k_taf"][name] = 0
	}
	if res, err := e.kiff.Call(ctx, "k_taf", "gmail_list_messages", map[string]any{"q": "from:ana@example.com", "maxResults": 5}); err != nil || res.Decision != "allowed" {
		t.Fatalf("list: %v %+v", err, res)
	}
	draft := map[string]any{"message": map[string]any{"raw": "VG86IGFuYUBleGFtcGxlLmNvbQ", "threadId": "t1"}}
	if res, err := e.kiff.Call(ctx, "k_taf", "gmail_create_draft", draft); err != nil || res.Decision != "allowed" {
		t.Fatalf("draft: %v %+v", err, res)
	}

	g, _ = e.s.loadGrant(ctx, tenant, id)
	token, item, err := e.s.openToken(ctx, g)
	if err != nil {
		t.Fatal(err)
	}
	token.Expiry = time.Now().Add(-time.Minute)
	if err := e.s.storeToken(ctx, g, token, &item.Ver); err != nil {
		t.Fatal(err)
	}
	if access, _, err := e.s.accessToken(ctx, g, 0); err != nil || access != "ya29.second" {
		t.Fatalf("refresh: %q %v", access, err)
	}

	_, view := e.do("GET", path, nil, nil, "k_setup")
	_, export := e.do("GET", "/v1/accounts/"+tenant+"/export", nil, e.phone, "")
	if st, out = e.do("DELETE", path, nil, e.phone, ""); st != 200 || out["status"] != StatusDisconnected {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	mu.Lock()
	defer mu.Unlock()
	if strings.Join(seen, ",") != "token:authorization_code,list,draft,token:refresh_token,revoke,revoke" {
		t.Fatalf("Google was asked: %v", seen)
	}
	raw, _ := json.Marshal(g)
	for what, text := range map[string]string{"grant view": mustJSONString(view), "export": mustJSONString(export), "grant record": string(raw), "logs": logs.String()} {
		if strings.Contains(text, secret) {
			t.Fatalf("the client secret is in the %s", what)
		}
	}
}

// gmailSearch is a provider whose search finds the managed Gmail's draft tool.
type gmailSearch struct{ fakeCatalogProvider }

func (gmailSearch) Search(context.Context, string, []string, string) (connector.Search, error) {
	return connector.Search{Tools: []connector.Tool{{Slug: "DESK_LIST", Toolkit: "desk_demo", Name: "List"}, {Slug: "GMAIL_CREATE_EMAIL_DRAFT", Toolkit: "gmail", Name: "Create email draft"}},
		Apps: map[string]connector.App{"desk_demo": {Slug: "desk_demo", Name: "Desk"}, "gmail": {Slug: "gmail", Name: "Gmail", Setup: connector.SetupReady}}}, nil
}

// When the reviewed Gmail and the managed Gmail are both found, the reviewed
// one comes first, even found by a tool's words alone; another reviewed
// entry found so still comes after the provider's apps.
func TestToolSearchListsReviewedGmailBeforeTheManagedOne(t *testing.T) {
	conf, _ := os.ReadFile("../../deploy/catalog.json")
	e := newInProcessEnv(t)
	e.enrol()
	var err error
	if e.s.CatalogEntries, err = DecodeCatalog(conf, shippedApps(t)); err != nil {
		t.Fatal(err)
	}
	e.s.Connectors = map[string]connector.Provider{"fake": gmailSearch{}}
	for _, q := range []string{"gmail", "labels%20clients"} {
		st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog/tools?q="+q, nil, nil, "k_setup")
		var order []string
		for _, x := range out["tools"].([]any) {
			order = append(order, x.(map[string]any)["entry"].(string))
		}
		reviewed, managed := slices.Index(order, gmailEntry), slices.Index(order, "connector:fake:gmail")
		if st != 200 || reviewed < 0 || managed < 0 || reviewed > managed || order[0] != gmailEntry {
			t.Fatalf("%s: %d %v", q, st, order)
		}
		if q != "gmail" && !strings.HasPrefix(order[len(order)-1], "qonto-key@") {
			t.Fatalf("Qonto's draft, found by its words, didn't stay after the provider's apps: %v", order)
		}
	}
}
