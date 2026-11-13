package planque

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

// The key and the Card in one approval, on a real key app's own description (qa-kiff). planque declares a key
// request's tools to KIFF before the key, so the taf's Card is drafted over them and the owner signs the key and the
// Card together; this reads Shippo's published description as it is, the key as "ShippoToken <key>" in
// Authorization, and asks for a read and a write by planque's tool names.
//
// It needs the description, which isn't kept in this repo:
//
//	curl -sSfo /tmp/shippo.yaml https://docs.goshippo.com/spec/shippoapi/public-api.yaml
//	PLANQUE_SHIPPO_DESCRIPTION=/tmp/shippo.yaml go test ./internal/planque -run TestTheKeyAndTheCardInOneApproval -v
//
// Shippo itself is never called: its server is a local stand-in, and the key is a made-up one.
func TestTheKeyAndTheCardInOneApprovalOnShipposDescription(t *testing.T) {
	path := os.Getenv("PLANQUE_SHIPPO_DESCRIPTION")
	if path == "" {
		t.Skip("PLANQUE_SHIPPO_DESCRIPTION isn't set: see this test's comment")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	const key = "shippo_test_made_up_0001"
	var mu sync.Mutex
	var heard []string
	vendor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/public-api.yaml" {
			_, _ = w.Write(raw)
			return
		}
		mu.Lock()
		heard = append(heard, r.Method+" "+r.URL.Path+" "+r.Header.Get("Authorization"))
		mu.Unlock()
		if r.Header.Get("Authorization") != "ShippoToken "+key {
			w.WriteHeader(401)
			_, _ = io.WriteString(w, `{"detail":"Invalid token."}`)
			return
		}
		_, _ = io.WriteString(w, `{"results":[],"next":null}`)
	}))
	defer vendor.Close()

	e := newEnv(t)
	e.enrol()
	// As connect-an-app@13 asks: the job's operations as tools, never none, so the request is declared.
	tools := []map[string]any{
		{"tool": "ListAddresses", "name": "list_addresses", "title": "List addresses", "read_only": true},
		{"tool": "CreateAddress", "name": "create_address", "title": "Create an address"},
	}
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": tools,
		"api": map[string]any{"url": vendor.URL + "/public-api.yaml", "base": vendor.URL,
			"auth": map[string]any{"in": "header", "name": "Authorization", "prefix": "ShippoToken "}}}, nil, "k_setup")
	if st != 201 || out["status"] != StatusPending || out["declared"] != true {
		t.Fatalf("the key request wasn't declared: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)

	// Before the key: both tools are in KIFF, so a Card can be drafted over them now, and a call reaches no one.
	declared := map[string]string{}
	for _, name := range []string{"list_addresses", "create_address"} {
		c, ok := e.kiff.Connection(tenant, name)
		if !ok {
			t.Fatalf("%s isn't declared in KIFF before the key: no Card can be drafted over it", name)
		}
		declared[name] = c.Credential
	}
	ctx := context.Background()
	e.kiff.Cards["k_taf"]["list_addresses"] = 0
	res, err := e.kiff.Call(ctx, "k_taf", "list_addresses", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "waits for its key") {
		t.Fatalf("a call before the key: %s", text)
	}
	if len(heard) != 0 {
		t.Fatalf("the vendor heard of it before the key: %v", heard)
	}

	// One approval: the owner signs the request's digest with the key, on the phone. It connects with the
	// registration the Card was drafted over, not a new one.
	if st, out := e.giveKey(id, digest, key); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("the one approval: %d %v", st, out)
	}
	for name, cred := range declared {
		if c, ok := e.kiff.Connection(tenant, name); !ok || c.Credential != cred {
			t.Fatalf("%s was registered again at connect: a Card drafted before the key would be over another tool", name)
		}
	}
	res, _ = e.kiff.Call(ctx, "k_taf", "list_addresses", map[string]any{})
	if text, isErr := resultText(t, res.Result); isErr {
		t.Fatalf("a call with the key: %s", text)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(heard) != 1 || heard[0] != "GET /addresses ShippoToken "+key {
		t.Fatalf("the vendor heard %v, want one GET /addresses with the key as ShippoToken", heard)
	}
}
