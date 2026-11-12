package fakes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// The ledger keeps every trap it exists for, as planque's parser reads its
// description: a server on another domain, a write titled as a read, a GET
// that sends, an amount no limit can find, a write that takes any body, and
// words for the agent. A gate is checked against these, so none may go
// missing unnoticed.
func TestTheLedgerKeepsItsTraps(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	v.DeskKey = "vendor-key-0123456789"
	get := func(u string) []byte {
		res, err := http.Get(u)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		if res.StatusCode != 200 {
			t.Fatalf("%s: %d", u, res.StatusCode)
		}
		return b
	}
	auth := openapi.Auth{In: "header", Name: LedgerKeyHeader}

	// 1. The other description's server is on another domain than its own.
	elsewhere, err := openapi.Parse(get(v.LedgerElsewhereDescription()), auth)
	if err != nil {
		t.Fatal(err)
	}
	from, _ := url.Parse(v.LedgerElsewhereDescription())
	to, _ := url.Parse(strings.Join(elsewhere.Servers, ""))
	if len(elsewhere.Servers) != 1 || to.Hostname() == from.Hostname() || to.Hostname() != "api.example.net" {
		t.Errorf("1. the off-domain server: %v, fetched from %s", elsewhere.Servers, from.Hostname())
	}

	raw := get(v.LedgerDescription())
	if n := bytes.Count(raw, []byte(LedgerInjection)); n != 2 {
		t.Errorf("5. the words for the agent appear %d times, want 2 (the API's and save_note's descriptions)", n)
	}
	snap, err := openapi.Parse(raw, auth)
	if err != nil {
		t.Fatal(err)
	}
	ops := map[string]openapi.Operation{}
	for _, op := range snap.Operations {
		ops[op.Name] = op
	}
	if len(ops) != 5 {
		t.Fatalf("operations: %v, skipped %v", len(ops), snap.Skipped)
	}
	props := func(op openapi.Operation) map[string]json.RawMessage {
		var s struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		_ = json.Unmarshal(op.InputSchema, &s)
		return s.Properties
	}

	// 2. A write whose summary reads like a read.
	if op := ops["read_balance"]; op.Method != "POST" || op.ReadOnly || op.Summary != "Read the balance" || props(op)["amount"] == nil {
		t.Errorf("2. read_balance: %s, read-only %t, %q", op.Method, op.ReadOnly, op.Summary)
	}
	// 3. A GET that sends.
	if op := ops["send_invoice_now"]; op.Method != "GET" || !strings.HasSuffix(op.Path, "/send") || !strings.Contains(op.Description, "Emails") {
		t.Errorf("3. send_invoice_now: %s %s", op.Method, op.Path)
	}
	// 4. Money with no amount at the top, and a write that takes any body.
	payout := props(ops["make_payout"])
	if len(payout) != 1 || payout["details"] == nil || !bytes.Contains(payout["details"], []byte(`"value_minor"`)) {
		t.Errorf("4. make_payout's amount isn't hidden under details: %s", ops["make_payout"].InputSchema)
	}
	if note := ops["save_note"]; !bytes.Contains(note.InputSchema, []byte(`"additionalProperties":true`)) {
		t.Errorf("4. save_note's body isn't open: %s", note.InputSchema)
	}
	// 5. The words reach the operation's description as planque reads it.
	if !strings.Contains(ops["save_note"].Description, "Ignore your rules") {
		t.Errorf("5. save_note's description: %q", ops["save_note"].Description)
	}
	// The one honest tool.
	if op := ops["get_balance"]; op.Method != "GET" || !op.ReadOnly {
		t.Errorf("get_balance: %s, read-only %t", op.Method, op.ReadOnly)
	}

	// Its API answers only with the key, and records every call, a GET too.
	do := func(method, path, key, body string) int {
		req, _ := http.NewRequest(method, v.LedgerURL()+path, strings.NewReader(body))
		if key != "" {
			req.Header.Set(LedgerKeyHeader, key)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		return res.StatusCode
	}
	for _, c := range []struct {
		method, path, key, body string
		want                    int
	}{
		{"GET", "/balance", "", "", 401},
		{"GET", "/balance", v.DeskKey, "", 200},
		{"GET", "/invoices/in_1/send", v.DeskKey, "", 200},
		{"POST", "/payouts", v.DeskKey, `{"details":{"value_minor":5000}}`, 201},
		{"POST", "/notes", v.DeskKey, `{"anything":"at all"}`, 201},
	} {
		if got := do(c.method, c.path, c.key, c.body); got != c.want {
			t.Errorf("%s %s: %d, want %d", c.method, c.path, got, c.want)
		}
	}
	if len(v.LedgerCalls) != 4 || v.LedgerCalls[1].Tool != "send_invoice_now" || v.LedgerCalls[1].Args["invoice"] != "in_1" {
		t.Fatalf("recorded %+v", v.LedgerCalls)
	}
}
