package fakes

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

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
	for _, name := range []string{"list_accounts", "list_charges", "stripe_api_read", "stripe_api_details"} {
		if !readOnly[name] {
			t.Errorf("%s isn't marked read-only", name)
		}
		if _, ok := v.read(context.Background(), name, map[string]any{"stripe_api_operation_id": "GetChargesCharge", "parameters": map[string]any{"id": "ch_1"}}); !ok {
			t.Errorf("%s isn't answered as a read", name)
		}
	}
	for _, name := range []string{"create_refund", "payout_balance", "stripe_api_write"} {
		if readOnly[name] {
			t.Errorf("%s is marked read-only", name)
		}
		if _, ok := v.read(context.Background(), name, nil); ok {
			t.Errorf("%s was answered as a read, not recorded as a call", name)
		}
	}
}

// One charge is read as Stripe reads it (kiff/tafs#174): the operation, and
// the charge at parameters.id, which stripe_api_details names. A read that
// names the field after the thing is refused, as Stripe refused it.
func TestVendorReadsAChargeByTheOperationsOwnField(t *testing.T) {
	v, _ := NewVendorAt("https://vendor.test")
	read := func(params map[string]any) any {
		out, _ := v.read(context.Background(), "stripe_api_read", map[string]any{"stripe_api_operation_id": "GetChargesCharge", "parameters": params})
		return out
	}
	if c, ok := read(map[string]any{"id": "ch_e2e_1"}).(map[string]any); !ok || c["id"] != "ch_e2e_1" || c["currency"] != "eur" || c["amount_refunded"] != 0 {
		t.Fatalf("by id: %v", c)
	}
	if e, ok := read(map[string]any{"charge": "ch_e2e_1"}).(vendorError); !ok || !strings.Contains(string(e), "Missing required path parameters: id") {
		t.Fatalf("by charge: %v", e)
	}
	d, _ := v.read(context.Background(), "stripe_api_details", map[string]any{"stripe_api_operation_id": "GetChargesCharge"})
	b, _ := json.Marshal(d)
	if !strings.Contains(string(b), `"name":"id","required":true`) {
		t.Fatalf("details = %s", b)
	}
	if _, ok := v.read(context.Background(), "stripe_api_read", map[string]any{"stripe_api_operation_id": "PostRefunds", "parameters": map[string]any{"id": "x"}}); !ok {
		t.Fatal("a write operation through the read tool isn't refused as a read")
	} else if out, _ := v.read(context.Background(), "stripe_api_read", map[string]any{"stripe_api_operation_id": "PostRefunds"}); func() bool { _, isErr := out.(vendorError); return !isErr }() {
		t.Fatalf("a write operation was read: %v", out)
	}
}

// GetRefunds lists the refunds the vendor made, newest first, as Stripe does:
// by either refund tool, of one charge, and within a created range.
func TestVendorListsItsRefunds(t *testing.T) {
	v, _ := NewVendorAt("https://vendor.test")
	ctx := context.Background()
	day := time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)
	for _, c := range []struct {
		at   time.Time
		call VendorCall
	}{
		{day.Add(-time.Hour), VendorCall{Tool: "create_refund", Args: map[string]any{"charge": "ch_yesterday", "amount": 100.0}}},
		{day.Add(9 * time.Hour), VendorCall{Tool: "create_refund", Args: map[string]any{"charge": "ch_a", "amount": 300.0}}},
		{day.Add(10 * time.Hour), VendorCall{Tool: "stripe_api_write", Args: map[string]any{"stripe_api_operation_id": "PostRefunds",
			"parameters": map[string]any{"charge": "ch_b", "amount": 800.0}}}},
		{day.Add(11 * time.Hour), VendorCall{Tool: "payout_balance", Args: map[string]any{"amount": 5000.0}}},
	} {
		doc, _ := json.Marshal(c.call)
		if err := v.keep(ctx, kvCall, fmt.Sprintf("%020d#re_%d", c.at.UnixNano(), c.at.Unix()), doc, time.Hour); err != nil {
			t.Fatal(err)
		}
	}
	list := func(params map[string]any) []string {
		t.Helper()
		out, ok := v.read(ctx, "stripe_api_read", map[string]any{"stripe_api_operation_id": "GetRefunds", "parameters": params})
		if !ok {
			t.Fatal("GetRefunds isn't a read")
		}
		if e, isErr := out.(vendorError); isErr {
			t.Fatalf("GetRefunds refused: %s", e)
		}
		var got []string
		for _, r := range out.(map[string]any)["data"].([]any) {
			m := r.(map[string]any)
			got = append(got, fmt.Sprintf("%s %v", m["charge"], m["amount"]))
		}
		return got
	}
	if got := list(nil); !slices.Equal(got, []string{"ch_b 800", "ch_a 300", "ch_yesterday 100"}) {
		t.Errorf("every refund, newest first, and no payout: %v", got)
	}
	today := map[string]any{"gte": float64(day.Unix()), "lt": float64(day.Add(24 * time.Hour).Unix())}
	if got := list(map[string]any{"created": today}); !slices.Equal(got, []string{"ch_b 800", "ch_a 300"}) {
		t.Errorf("today's refunds: %v", got)
	}
	if got := list(map[string]any{"charge": "ch_a"}); !slices.Equal(got, []string{"ch_a 300"}) {
		t.Errorf("one charge's refunds: %v", got)
	}
	if d, _ := v.read(ctx, "stripe_api_details", map[string]any{"stripe_api_operation_id": "GetRefunds"}); d == nil {
		t.Error("stripe_api_details doesn't know GetRefunds")
	}
}
