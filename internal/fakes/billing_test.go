package fakes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// Billing answers only with its key in its own header, makes a finalized
// invoice when status is left out, records each write once, and its
// description is one planque reads: a read and two writes, create_invoice
// with status in its body, where a grant can pin it.
func TestBillingIsAnAPIKeyVendorWhoseStatusDecidesFinal(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	v.DeskKey = "vendor-key-0123456789"
	do := func(method, path, header, key, body string) int {
		req, _ := http.NewRequest(method, v.BillingURL()+path, bytes.NewReader([]byte(body)))
		if key != "" {
			req.Header.Set(header, key)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = io.Copy(io.Discard, res.Body)
		res.Body.Close()
		return res.StatusCode
	}
	key := v.DeskKey
	for _, c := range []struct {
		method, path, header, key, body string
		want                            int
	}{
		{"GET", "/orders", BillingKeyHeader, "", "", 401},
		{"GET", "/orders", BillingKeyHeader, "wrong", "", 401},
		{"GET", "/orders", DeskKeyHeader, key, "", 401}, // the desk's header isn't billing's
		{"GET", "/orders?status=fulfilled", BillingKeyHeader, key, "", 200},
		{"POST", "/invoices", BillingKeyHeader, key, `{"order":"or_9"}`, 404},
		{"POST", "/invoices", BillingKeyHeader, key, `{"order":"or_2001","status":"sent"}`, 422},
		{"POST", "/invoices", BillingKeyHeader, key, `{"order":"or_2001","status":"draft"}`, 201},
		{"POST", "/invoices", BillingKeyHeader, key, `{"order":"or_2002"}`, 201},
		{"POST", "/invoices/x_1/send", BillingKeyHeader, key, "", 404},
		{"POST", "/invoices/in_1/send", BillingKeyHeader, key, "", 200},
	} {
		if got := do(c.method, c.path, c.header, c.key, c.body); got != c.want {
			t.Errorf("%s %s (%s): %d, want %d", c.method, c.path, c.body, got, c.want)
		}
	}
	calls := v.BillingCalls
	if len(calls) != 3 || calls[0].Args["status"] != InvoiceDraft || calls[1].Args["status"] != InvoiceFinalized || calls[2].Tool != "send_invoice" {
		t.Fatalf("recorded %+v; want a draft, then a finalized invoice for status left out, then a send", calls)
	}

	res, err := http.Get(v.BillingDescription())
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	spec, err := openapi.Parse(raw, openapi.Auth{In: "header", Name: BillingKeyHeader})
	if err != nil {
		t.Fatal(err)
	}
	ops := map[string]openapi.Operation{}
	for _, op := range spec.Operations {
		ops[op.Name] = op
	}
	for name, read := range map[string]bool{"list_orders": true, "create_invoice": false, "send_invoice": false} {
		op, ok := ops[name]
		if !ok || op.ReadOnly != read {
			t.Errorf("%s: found %t, read-only %t; want read-only %t", name, ok, op.ReadOnly, read)
		}
	}
	create := ops["create_invoice"]
	var schema struct {
		Properties map[string]struct {
			Enum        []string `json:"enum"`
			Description string   `json:"description"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if err := json.Unmarshal(create.InputSchema, &schema); err != nil {
		t.Fatal(err)
	}
	status := schema.Properties["status"]
	if len(status.Enum) != 2 || !bytes.Contains([]byte(status.Description), []byte("Left out, the invoice is finalized")) {
		t.Errorf("create_invoice's status is %+v; want draft or finalized, saying what leaving it out does", status)
	}

	// A grant that pins status to draft: the body the call makes carries it,
	// and the vendor makes a draft.
	method, target, body, err := create.Request(v.BillingURL(), map[string]any{"order": "or_2001", "status": InvoiceDraft})
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set(BillingKeyHeader, key)
	pinned, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	pinned.Body.Close()
	if last := v.BillingCalls[len(v.BillingCalls)-1]; pinned.StatusCode != 201 || last.Args["status"] != InvoiceDraft {
		t.Fatalf("a call with status pinned to draft: %d, recorded %+v", pinned.StatusCode, last)
	}
}
