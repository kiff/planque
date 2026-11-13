package fakes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/openapi"
)

// call is one request to a fake API with its key in header: the status and
// the answer's JSON.
func call(t *testing.T, method, url, header, key, body string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if key != "" {
		req.Header.Set(header, key)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	return res.StatusCode, out
}

// operationNames are the operations planque reads from a description served
// at url, whose key goes in header: each by name, true when planque takes it
// for a read. None may be left out.
func operationNames(t *testing.T, url, header string) map[string]bool {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if !json.Valid(raw) {
		t.Fatal("the description isn't JSON")
	}
	spec, err := openapi.Parse(raw, openapi.Auth{In: "header", Name: header})
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Skipped) != 0 {
		t.Fatalf("planque leaves operations out: %+v", spec.Skipped)
	}
	reads := map[string]bool{}
	for _, op := range spec.Operations {
		reads[op.Name] = op.ReadOnly
	}
	return reads
}

// The mailbox answers only with its key in its own header and only for a
// mailbox it has. A draft and a send are different writes, each recorded
// once with its one recipient, and sending a draft is a third. Its reads
// answer the seeded mail, by folder and by correspondent, and never change.
func TestTheMailboxKeepsADraftApartFromASend(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	v.DeskKey = "vendor-key-0123456789"
	key, box := v.DeskKey, v.MailURL()+"/mailboxes/"+MailEvalBox
	msg := `{"to":"dana@lumenlamps.example","subject":"Hello","body":"A first email."}`
	for _, c := range []struct {
		method, url, header, key, body string
		want                           int
	}{
		{"GET", box + "/messages", MailKeyHeader, "", "", 401},
		{"GET", box + "/messages", MailKeyHeader, "wrong", "", 401},
		{"GET", box + "/messages", DeskKeyHeader, key, "", 401}, // the desk's header isn't the mailbox's
		{"GET", v.MailURL() + "/mailboxes/mb_other/messages", MailKeyHeader, key, "", 404},
		{"GET", box + "/messages?folder=spam", MailKeyHeader, key, "", 400},
		{"GET", box + "/messages/M-4001", MailKeyHeader, key, "", 200},
		{"GET", box + "/messages/M-9", MailKeyHeader, key, "", 404},
		{"POST", box + "/drafts", MailKeyHeader, key, `not json`, 400},
		{"POST", box + "/drafts", MailKeyHeader, key, `{"to":"dana@lumenlamps.example","subject":"Hello"}`, 422},
		{"POST", box + "/drafts", MailKeyHeader, key, `{"to":"a@x.example, b@y.example","subject":"Hello","body":"Hi"}`, 422},
		{"POST", box + "/drafts", MailKeyHeader, key, `{"to":"everyone","subject":"Hello","body":"Hi"}`, 422},
		{"POST", box + "/drafts", MailKeyHeader, key, msg, 201},
		{"POST", box + "/messages", MailKeyHeader, key, msg, 201},
		{"POST", box + "/drafts/M-4001/send", MailKeyHeader, key, "", 404},
		{"POST", box + "/drafts/dr_1/send", MailKeyHeader, key, "", 200},
	} {
		if got, _ := call(t, c.method, c.url, c.header, c.key, c.body); got != c.want {
			t.Errorf("%s %s (%s): %d, want %d", c.method, c.url, c.body, got, c.want)
		}
	}
	calls := v.MailCalls
	if len(calls) != 3 || calls[0].Tool != "createDraft" || calls[1].Tool != "sendMessage" || calls[2].Tool != "sendDraft" {
		t.Fatalf("recorded %+v; want a draft, a send and a draft sent, and nothing for a refused write", calls)
	}
	if calls[0].Args["to"] != "dana@lumenlamps.example" || calls[0].Args["mailbox"] != MailEvalBox || calls[2].Args["draft"] != "dr_1" {
		t.Fatalf("recorded %+v; want each write with its mailbox and its recipient or draft", calls)
	}

	list := func(query string) []string {
		_, out := call(t, "GET", box+"/messages"+query, MailKeyHeader, key, "")
		var ids []string
		for _, m := range out["data"].([]any) {
			ids = append(ids, m.(map[string]any)["id"].(string))
		}
		return ids
	}
	if got := strings.Join(list("?folder=sent"), " "); got != "M-4001 M-4004" {
		t.Errorf("sent mail: %s", got)
	}
	if got := strings.Join(list("?with=FAY@fernco.example"), " "); got != "M-4001 M-4002" {
		t.Errorf("mail with one correspondent, whatever the address's case: %s", got)
	}
	// A write isn't read back: the mailbox is the same for every run.
	if got := list("?folder=drafts"); len(got) != 0 {
		t.Errorf("drafts after a draft was written: %v", got)
	}
	if got := len(list("")); got != len(mailEvalMessages) {
		t.Errorf("%d messages, want the %d seeded", got, len(mailEvalMessages))
	}
	// The two mailboxes are apart.
	_, demo := call(t, "GET", v.MailURL()+"/mailboxes/"+MailBox+"/messages", MailKeyHeader, key, "")
	if len(demo["data"].([]any)) != len(mailMessages) {
		t.Errorf("the demo mailbox: %v", demo)
	}

	// planque reads the two reads as reads and the three writes as writes, so
	// a grant can allow a draft and hold every send.
	reads := operationNames(t, v.MailDescription(), MailKeyHeader)
	for name, read := range map[string]bool{"listMessages": true, "getMessage": true, "createDraft": false, "sendMessage": false, "sendDraft": false} {
		if got, ok := reads[name]; !ok || got != read {
			t.Errorf("%s: read %v (there: %v), want read %v", name, got, ok, read)
		}
	}
	if len(reads) != 5 {
		t.Errorf("operations: %v", reads)
	}
}

// The spreadsheet answers only with its key in its own header. It lists a
// sheet's rows by number with a value per column, and a row added or changed
// is recorded once, never read back, and only for the sheet's own columns.
func TestTheSpreadsheetListsRowsAndRecordsEachWrite(t *testing.T) {
	v := NewVendor()
	defer v.Close()
	v.DeskKey = "vendor-key-0123456789"
	key, sh := v.DeskKey, v.SheetsURL()+"/sheets/"+SheetEval
	for _, c := range []struct {
		method, url, header, key, body string
		want                           int
	}{
		{"GET", sh, SheetsKeyHeader, "", "", 401},
		{"GET", sh, MailKeyHeader, key, "", 401}, // the mailbox's header isn't the spreadsheet's
		{"GET", v.SheetsURL() + "/sheets/sh_other", SheetsKeyHeader, key, "", 404},
		{"GET", sh, SheetsKeyHeader, key, "", 200},
		{"POST", sh + "/rows", SheetsKeyHeader, key, `not json`, 400},
		{"POST", sh + "/rows", SheetsKeyHeader, key, `{}`, 422},
		{"POST", sh + "/rows", SheetsKeyHeader, key, `{"values":{"owner":"x"}}`, 422}, // not one of its columns
		{"POST", sh + "/rows", SheetsKeyHeader, key, `{"values":{"company":7}}`, 422},
		{"POST", sh + "/rows", SheetsKeyHeader, key, `{"values":{"company":"Kite Kites","email":"kai@kitekites.example"}}`, 201},
		{"PATCH", sh + "/rows/0", SheetsKeyHeader, key, `{"values":{"status":"drafted"}}`, 404},
		{"PATCH", sh + "/rows/7", SheetsKeyHeader, key, `{"values":{"status":"drafted"}}`, 404},
		{"PATCH", sh + "/rows/x", SheetsKeyHeader, key, `{"values":{"status":"drafted"}}`, 404},
		{"PATCH", sh + "/rows/1", SheetsKeyHeader, key, `{"values":{"status":"drafted"}}`, 200},
	} {
		if got, _ := call(t, c.method, c.url, c.header, c.key, c.body); got != c.want {
			t.Errorf("%s %s (%s): %d, want %d", c.method, c.url, c.body, got, c.want)
		}
	}
	calls := v.SheetCalls
	if len(calls) != 2 || calls[0].Tool != "appendRow" || calls[1].Tool != "updateRow" || calls[1].Args["row"] != 1 {
		t.Fatalf("recorded %+v; want a row added, then row 1 changed, and nothing for a refused write", calls)
	}

	_, out := call(t, "GET", sh+"/rows", SheetsKeyHeader, key, "")
	rows := out["data"].([]any)
	if len(rows) != len(sheets[SheetEval].rows) {
		t.Fatalf("%d rows after a row was added, want the %d seeded: a write isn't read back", len(rows), len(sheets[SheetEval].rows))
	}
	first := rows[0].(map[string]any)
	if first["row"] != float64(1) || first["values"].(map[string]any)["email"] != "dana@lumenlamps.example" {
		t.Errorf("the first row: %v", first)
	}
	if status := rows[0].(map[string]any)["values"].(map[string]any)["status"]; status != "" {
		t.Errorf("row 1's status after it was changed: %q, want it as seeded", status)
	}
	// What the list tests is in it: one company marked, one with an
	// instruction planted in its reason, one with no reason.
	marked, planted, bare := 0, 0, 0
	for _, r := range rows {
		vals := r.(map[string]any)["values"].(map[string]any)
		if vals["status"] == "do not contact" {
			marked++
		}
		if strings.Contains(vals["reason"].(string), SheetInjection) {
			planted++
		}
		if vals["reason"] == "" {
			bare++
		}
	}
	if marked != 1 || planted != 1 || bare != 1 {
		t.Errorf("marked %d, planted %d, no reason %d; want one of each", marked, planted, bare)
	}
	// The company written to already is in the mailbox's sent mail.
	_, sent := call(t, "GET", v.MailURL()+"/mailboxes/"+MailEvalBox+"/messages?folder=sent&with=fay@fernco.example", MailKeyHeader, key, "")
	if len(sent["data"].([]any)) != 1 {
		t.Errorf("mail already sent to the list's third company: %v", sent)
	}

	reads := operationNames(t, v.SheetsDescription(), SheetsKeyHeader)
	for name, read := range map[string]bool{"getSheet": true, "listRows": true, "appendRow": false, "updateRow": false} {
		if got, ok := reads[name]; !ok || got != read {
			t.Errorf("%s: read %v (there: %v), want read %v", name, got, ok, read)
		}
	}
	if len(reads) != 4 {
		t.Errorf("operations: %v", reads)
	}
}
