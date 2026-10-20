package fakes

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/store"
)

// A deployed fake runs as several instances over one table: a client
// registered at one, signed in at another and called at a third is one
// vendor, and every instance reads every call.
func TestVendorInstancesShareTheirState(t *testing.T) {
	kv := &store.Memory{}
	var hs [3]http.Handler
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	vs := make([]*Vendor, 3)
	for i := range hs {
		vs[i], hs[i] = NewVendorAt(srv.URL)
		vs[i].KV = kv
	}
	n := 0
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { hs[n%3].ServeHTTP(w, r); n++ })
	noFollow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

	resp, err := http.Post(srv.URL+"/as/register", "application/json", strings.NewReader(`{"redirect_uris":["https://planque.test/callback"]}`))
	if err != nil {
		t.Fatal(err)
	}
	var reg struct {
		ClientID string `json:"client_id"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&reg)
	resp.Body.Close()

	verifier := "a-verifier-long-enough-for-pkce-0123456789abcdef"
	sum := sha256.Sum256([]byte(verifier))
	q := url.Values{"client_id": {reg.ClientID}, "redirect_uri": {"https://planque.test/callback"}, "code_challenge_method": {"S256"},
		"code_challenge": {base64.RawURLEncoding.EncodeToString(sum[:])}, "resource": {srv.URL + "/mcp"}}
	resp, err = noFollow.Get(srv.URL + "/as/authorize?" + q.Encode())
	if err != nil || resp.StatusCode != http.StatusFound {
		t.Fatalf("authorize: %v %v", resp, err)
	}
	loc, _ := url.Parse(resp.Header.Get("Location"))
	resp.Body.Close()

	exchange := url.Values{"grant_type": {"authorization_code"}, "code": {loc.Query().Get("code")}, "client_id": {reg.ClientID},
		"redirect_uri": {"https://planque.test/callback"}, "code_verifier": {verifier}}
	resp, err = http.PostForm(srv.URL+"/as/token", exchange)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("token: %v %v", resp, err)
	}
	var tok struct {
		Access string `json:"access_token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tok)
	resp.Body.Close()
	if resp, _ := http.PostForm(srv.URL+"/as/token", exchange); resp.StatusCode != 400 {
		t.Fatalf("a code worked twice, at another instance: %d", resp.StatusCode)
	}

	req, _ := http.NewRequest("POST", srv.URL+"/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_refund","arguments":{"charge":"ch_1","amount":300}}}`))
	req.Header.Set("Authorization", "Bearer "+tok.Access)
	if resp, err := http.DefaultClient.Do(req); err != nil || resp.StatusCode != 200 {
		t.Fatalf("tools/call: %v %v", resp, err)
	}
	for i, v := range vs {
		calls, err := v.CallsSoFar(context.Background())
		if err != nil || len(calls) != 1 || calls[0].Tool != "create_refund" {
			t.Fatalf("instance %d reads %v %v", i, calls, err)
		}
	}
}

// A vendor that consents on its own registers only the redirects it is
// given, or it would send anyone anywhere.
func TestVendorRegistersOnlyItsRedirects(t *testing.T) {
	v, h := NewVendorAt("https://vendor.test")
	v.Redirects = []string{"https://planque.test/oauth/callback"}
	for uri, want := range map[string]int{"https://planque.test/oauth/callback": 201, "https://elsewhere.test/": 400} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("POST", "/as/register", strings.NewReader(`{"redirect_uris":["`+uri+`"]}`)))
		if rec.Code != want {
			t.Errorf("registering %s: %d, want %d", uri, rec.Code, want)
		}
	}
}

// A taf finds and checks its work with the vendor's read-only tools before
// it acts (kiff/tafs#107): they answer without being one of the vendor's
// calls, and only the writes are recorded.
func TestVendorReadsAreNotCalls(t *testing.T) {
	v, _ := NewVendorAt("https://vendor.test")
	readOnly := map[string]bool{}
	for _, tool := range vendorTools {
		a, _ := tool["annotations"].(map[string]any)
		readOnly[tool["name"].(string)] = a["readOnlyHint"] == true
	}
	for _, name := range []string{"list_accounts", "list_charges", "retrieve_charge"} {
		if !readOnly[name] {
			t.Errorf("%s isn't marked read-only", name)
		}
		if _, ok := v.read(name, map[string]any{"charge": "ch_1"}); !ok {
			t.Errorf("%s isn't answered as a read", name)
		}
	}
	for _, name := range []string{"create_refund", "payout_balance", "stripe_api_write"} {
		if readOnly[name] {
			t.Errorf("%s is marked read-only", name)
		}
		if _, ok := v.read(name, nil); ok {
			t.Errorf("%s was answered as a read, not recorded as a call", name)
		}
	}
	out, _ := v.read("retrieve_charge", map[string]any{"charge": "ch_e2e_1"})
	if c := out.(map[string]any); c["id"] != "ch_e2e_1" || c["amount_refunded"] != 0 {
		t.Fatalf("retrieve_charge = %v", out)
	}
}
