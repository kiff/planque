package fakes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// The desk reads only with its key and workspace, records each write once,
// and its description is one planque reads: four operations, the key in a
// header, reads named apart from writes.
func TestTheDeskIsAnAPIKeyVendor(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	v.DeskKey = "desk-key-0123456789"
	do := func(method, path, key, body string) int {
		req, _ := http.NewRequest(method, v.DeskURL()+path, bytes.NewReader([]byte(body)))
		if key != "" {
			req.Header.Set(DeskKeyHeader, key)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = io.Copy(io.Discard, res.Body)
		res.Body.Close()
		return res.StatusCode
	}
	for _, c := range []struct {
		method, path, key, body string
		want                    int
	}{
		{"GET", "/workspaces/ws_demo/tickets", "", "", 401},
		{"GET", "/workspaces/ws_demo/tickets", "wrong", "", 401},
		{"GET", "/workspaces/ws_other/tickets", v.DeskKey, "", 404},
		{"GET", "/workspaces/ws_demo/tickets?status=open", v.DeskKey, "", 200},
		{"GET", "/workspaces/ws_demo/tickets/tk_1001", v.DeskKey, "", 200},
		{"GET", "/workspaces/ws_demo/tickets/tk_9", v.DeskKey, "", 404},
		{"POST", "/workspaces/ws_demo/tickets/tk_1001/replies", v.DeskKey, `{"body":""}`, 422},
		{"POST", "/workspaces/ws_demo/tickets/tk_1001/replies", v.DeskKey, `{"body":"Sorry, a new one is on its way."}`, 201},
		{"POST", "/workspaces/ws_demo/tickets/tk_1001/labels", v.DeskKey, `{"label":"damaged"}`, 201},
		// The Eval's workspace is its own: its tickets aren't the setup test's, and theirs aren't in it.
		{"GET", "/workspaces/ws_eval/tickets/T-3006", v.DeskKey, "", 200},
		{"GET", "/workspaces/ws_eval/tickets/tk_1001", v.DeskKey, "", 404},
		{"GET", "/workspaces/ws_demo/tickets/T-3001", v.DeskKey, "", 404},
	} {
		if got := do(c.method, c.path, c.key, c.body); got != c.want {
			t.Errorf("%s %s: %d, want %d", c.method, c.path, got, c.want)
		}
	}
	if len(v.DeskCalls) != 2 || v.DeskCalls[0].Tool != "postReply" || v.DeskCalls[1].Args["label"] != "damaged" {
		t.Fatalf("recorded %+v", v.DeskCalls)
	}

	res, err := http.Get(v.DeskDescription())
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if !json.Valid(raw) {
		t.Fatal("the description isn't JSON")
	}
	spec, err := openapi.Parse(raw, openapi.Auth{In: "header", Name: DeskKeyHeader})
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, op := range spec.Operations {
		names[op.Name] = true
	}
	for _, want := range []string{"searchTickets", "fetchTicket", "postReply", "applyLabel"} {
		if !names[want] {
			t.Errorf("no operation %s in %v", want, names)
		}
	}
}
