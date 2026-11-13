package planque

import (
	"strings"
	"testing"

	"github.com/kiff/planque/internal/profile"
)

// The negatives for narrowing a reviewed entry to the tools a taf uses:
// what a narrowed request must never let through.

// The shipped entry's write checks its IBAN against the business's own
// accounts, read with get_organization: asked for alone it's refused and
// nothing is declared; with that read it's asked for, and nothing else is.
func TestTheDraftIsNeverAskedForWithoutTheReadItsIBANCheckNeeds(t *testing.T) {
	e := qontoKeyEnv(t)
	path := "/v1/accounts/" + tenant + "/catalog/qonto-key@2/grants"
	if st, out := e.do("POST", path, map[string]any{"agent": "agt_taf", "tools": []string{"create_client_invoice"}}, nil, "k_setup"); st != 400 {
		t.Fatalf("the draft without the read its IBAN check needs: %d %v", st, out)
	}
	if c := e.kiff.Connections(tenant); len(c) != 0 {
		t.Fatalf("declared after a refusal: %+v", c)
	}
	st, out := e.do("POST", path, map[string]any{"agent": "agt_taf", "tools": []string{"get_organization", "create_client_invoice"}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("the draft with its read: %d %v", st, out)
	}
	var names []string
	for _, c := range e.kiff.Connections(tenant) {
		names = append(names, c.Name)
	}
	if strings.Join(names, ",") != "qonto_organization,qonto_draft_client_invoice" && strings.Join(names, ",") != "qonto_draft_client_invoice,qonto_organization" {
		t.Fatalf("declared in KIFF: %v; want the read and the draft only", names)
	}
}

// Reads only: no write is declared, so a Card can't be drafted over one and
// a call to it reaches nothing.
func TestAReadsOnlySelectionDeclaresNoWrite(t *testing.T) {
	e := qontoKeyEnv(t)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/qonto-key@2/grants",
		map[string]any{"agent": "agt_taf", "tools": []string{"list_client_invoices", "get_client_invoice"}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	if _, ok := e.kiff.Connection(tenant, "qonto_draft_client_invoice"); ok {
		t.Fatal("the draft is declared for a taf that asked only to read")
	}
	if len(e.kiff.Connections(tenant)) != 2 {
		t.Fatalf("declared: %+v, want the two reads", e.kiff.Connections(tenant))
	}
}

// Every read a tool's profile does to check its calls counts as one the
// narrowed request must name with it: its amount and unit reads, its verify
// read and its status read.
func TestEveryReadAToolsChecksDoIsOneItNeeds(t *testing.T) {
	o := profile.Operation{
		Tool:       "x_send",
		Amount:     &profile.Amount{From: "read", Read: &profile.Read{Tool: "x_amount"}, UnitRead: &profile.Read{Tool: "x_unit"}},
		VerifyRead: &profile.VerifyRead{Tool: "x_verify"},
		StatusRead: &profile.StatusRead{Tool: "x_status"},
	}
	got := map[string]bool{}
	for _, r := range profileReads(o) {
		got[r] = true
	}
	for _, want := range []string{"x_amount", "x_unit", "x_verify", "x_status"} {
		if !got[want] {
			t.Errorf("%s isn't counted as a read x_send's checks need: %v", want, profileReads(o))
		}
	}
}
