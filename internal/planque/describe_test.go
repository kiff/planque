package planque

import (
	"context"
	"encoding/json"
	"maps"
	"slices"
	"strings"
	"testing"
)

// detailsGrant is a grant, connected first, whose one tool is the vendor's
// describing tool, for Setup's own reads: what a broad write is narrowed
// on, once the owner approved it.
func (e *env) detailsGrant() string {
	e.t.Helper()
	id := e.listedGrant()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	details := map[string]any{"tool": "stripe_api_details", "name": "fields_ab12", "read_only": true, "setup_reads": true}
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{details}}, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("set tools: %d %v", st, out)
	}
	if st, out := e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve: %d %v", st, out)
	}
	return id
}

// addToDetails adds tool to a detailsGrant, as Setup does: the grant's id
// and planque's answer, the addition waiting for the owner.
func (e *env) addToDetails(tool map[string]any) (string, map[string]any) {
	e.t.Helper()
	id := e.detailsGrant()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []any{tool}}, nil, "k_setup")
	if st != 200 {
		e.t.Fatalf("add: %d %v", st, out)
	}
	return id, out
}

func (e *env) approveAddition(id string, out map[string]any) {
	e.t.Helper()
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		e.t.Fatalf("approve the addition: %d %v", st, out)
	}
}

// A broad write narrowed to one operation takes only that operation's
// fields, as the vendor's describing tool lists them, read once the owner
// approved it: a made-up field is refused before the owner is asked, not by
// the vendor at the first call after it.
func TestANarrowedWriteTakesOnlyTheOperationsFields(t *testing.T) {
	refund := func(args map[string]any) map[string]any {
		return map[string]any{"tool": "stripe_api_write", "name": "refund_cd34", "title": "Refund a payment", "amount_argument": "amount",
			"fixed_arguments": map[string]any{"stripe_api_operation_id": "PostRefunds"}, "arguments": args}
	}
	arg := func(path, typ string) map[string]any { return map[string]any{"path": path, "type": typ} }
	good := map[string]any{"charge": arg("parameters.charge", "string"), "amount": arg("parameters.amount", "integer"),
		"note": arg("parameters.metadata.note", "string"), "approval_token": arg("human_confirmation.approval_token", "string")}

	// Asked for at once, nothing is approved yet: the describing tool comes first.
	e := newEnv(t)
	e.enrol()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf",
		"tools": []any{refund(good)}}, nil, "k_setup")
	if st != 400 || !strings.Contains(out["error"].(string), "connect stripe_api_details first (with setup_reads), and add refund_cd34 once the owner approved it") {
		t.Fatalf("at once: %d %v", st, out)
	}
	// So on a listed grant, before the owner approved the describing tool.
	id := e.listedGrant()
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []any{refund(good)}}, nil, "k_setup"); st != 400 ||
		!strings.Contains(out["error"].(string), "connect stripe_api_details first") {
		t.Fatalf("listed: %d %v", st, out)
	}

	for name, c := range map[string]struct {
		args map[string]any
		want string
	}{
		"a made-up field": {map[string]any{"amount": arg("parameters.amount", "integer"), "amount_to_counter": arg("parameters.amount_to_counter", "integer")},
			"PostRefunds takes no field amount_to_counter: stripe_api_details lists amount, charge, metadata, payment_intent"},
		"inside a plain field": {map[string]any{"amount": arg("parameters.amount", "integer"), "x": arg("parameters.charge.id", "string")},
			"PostRefunds takes no field charge.id"},
		"all of parameters": {map[string]any{"amount": arg("human_confirmation.amount", "integer"), "all": arg("parameters", "string")},
			"parameters takes one field of PostRefunds, never all of parameters"},
	} {
		e := newEnv(t)
		id := e.detailsGrant()
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []any{refund(c.args)}}, nil, "k_setup")
		if st != 400 || !strings.Contains(out["error"].(string), c.want) {
			t.Errorf("%s: %d %v", name, st, out)
		}
	}

	// The operation's own fields, one inside an object it doesn't detail, and
	// a path outside its fields, are taken; the read is on the grant's record.
	e = newEnv(t)
	id, out = e.addToDetails(refund(good))
	if out["status"] != StatusPending {
		t.Fatalf("add: %v", out)
	}
	_, g := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+id, nil, nil, "k_setup")
	if b, _ := json.Marshal(g["events"]); !strings.Contains(string(b), "read PostRefunds's fields with stripe_api_details at ") {
		t.Fatalf("events: %s", b)
	}
	e.approveAddition(id, out)

	// An operation the vendor doesn't describe is refused with its words.
	e = newEnv(t)
	id = e.detailsGrant()
	payout := refund(map[string]any{"amount": arg("parameters.amount", "integer")})
	payout["fixed_arguments"] = map[string]any{"stripe_api_operation_id": "PostPayouts"}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []any{payout}}, nil, "k_setup"); st != 400 ||
		!strings.Contains(out["error"].(string), "doesn't describe PostPayouts: Operation 'PostPayouts' is not available") {
		t.Fatalf("an operation it doesn't describe: %d %v", st, out)
	}
}

// The describing tool reads the operation in the account the narrowed tool
// pins: Stripe describes nothing without stripe_context, so a check without
// it refused every narrowed tool on a real account (the 1 Nov phone run).
func TestTheFieldsAreReadInTheNarrowedToolsAccount(t *testing.T) {
	search := func(acct string, live bool) map[string]any {
		return map[string]any{"tool": "stripe_api_read", "name": "search_charges", "read_only": true,
			"fixed_arguments": map[string]any{"stripe_api_operation_id": "GetChargesCharge", "stripe_context": acct, "livemode": live},
			"arguments":       map[string]any{"id": map[string]any{"path": "parameters.id", "type": "string"}}}
	}
	e := newEnv(t)
	id := e.detailsGrant()
	at := "/v1/accounts/" + tenant + "/grants/" + id + "/tools"
	if st, out := e.do("POST", at, map[string]any{"tools": []any{search("acct_other", false)}}, nil, "k_setup"); st != 400 ||
		!strings.Contains(out["error"].(string), "No such account: 'acct_other'") {
		t.Fatalf("another account: %d %v", st, out)
	}
	// And in its mode: the test sign-in describes nothing in live mode.
	if st, out := e.do("POST", at, map[string]any{"tools": []any{search(e.vendor.Account, true)}}, nil, "k_setup"); st != 400 ||
		!strings.Contains(out["error"].(string), "no access to live mode") {
		t.Fatalf("live mode: %d %v", st, out)
	}
	if st, out := e.do("POST", at, map[string]any{"tools": []any{search(e.vendor.Account, false)}}, nil, "k_setup"); st != 200 {
		t.Fatalf("the sign-in's account: %d %v", st, out)
	}
}

// A narrowed read isn't held back: it needs no describing tool, as before,
// and the vendor refuses a field it doesn't take when it's called.
func TestANarrowedReadNeedsNoDescribingTool(t *testing.T) {
	e := newEnv(t)
	id := e.listedGrant()
	read := map[string]any{"tool": "stripe_api_read", "name": "get_charge", "read_only": true,
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "GetChargesCharge"},
		"arguments":       map[string]any{"id": map[string]any{"path": "parameters.id", "type": "string"}}}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []any{read}}, nil, "k_setup"); st != 200 {
		t.Fatalf("read: %d %v", st, out)
	}
}

func TestParseFields(t *testing.T) {
	text := "# PostDisputesDispute\n**POST /v1/disputes/{dispute}**\n\nUpdate a dispute\n\n## Path Parameters\n- **dispute** [string] (required)\n\n" +
		"## Body Parameters\n- **evidence** [object] - Evidence.\n  - **product_description** [string] - What was sold.\n  - **receipt** [string]\n" +
		"- **metadata** [object]\n- **submit** [boolean] - Whether to submit.\n\nSome text\n- **not_a_field** [string]\n"
	f := parseFields(text)
	if got := listFields(f); got != "dispute, evidence, metadata, not_a_field, submit" {
		// A line after the last section's list is still in that section:
		// only a heading ends it.
		t.Fatalf("fields: %s", got)
	}
	if got := strings.Join(slices.Sorted(maps.Keys(f["evidence"].kids)), ", "); got != "product_description, receipt" {
		t.Fatalf("evidence: %s", got)
	}
	if !f["metadata"].open || f["submit"].open {
		t.Fatalf("open: metadata %v, submit %v", f["metadata"].open, f["submit"].open)
	}
}

// A narrowed write defined before planque checked its fields is checked at
// its first call: an action naming a field its operation doesn't take (the
// Chargebacks answer's amount_to_counter) is refused definition_invalid,
// with the reason, instead of going to the vendor to fail on every call. One
// that holds is marked, and isn't read again.
func TestANarrowedWriteFromBeforeTheCheckIsCheckedAtItsFirstCall(t *testing.T) {
	arg := func(path, typ string) map[string]any { return map[string]any{"path": path, "type": typ} }
	refund := map[string]any{"tool": "stripe_api_write", "name": "refund_cd34", "title": "Refund a payment", "amount_argument": "amount",
		"fixed_arguments": map[string]any{"stripe_api_operation_id": "PostRefunds"},
		"arguments":       map[string]any{"charge": arg("parameters.charge", "string"), "amount": arg("parameters.amount", "integer")}}
	e := newEnv(t)
	id, out := e.addToDetails(refund)
	e.approveAddition(id, out)
	ctx := context.Background()
	load := func() (*Grant, int) {
		g, err := e.s.loadGrant(ctx, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
		return g, slices.IndexFunc(g.Tools, func(gt GrantedTool) bool { return gt.Name == "refund_cd34" })
	}

	// As before the check: a made-up field, never checked.
	g, i := load()
	g.Tools[i].Arguments["amount_to_counter"] = Argument{Path: "parameters.amount_to_counter", Type: "integer"}
	g.Tools[i].FieldsChecked = false
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	g, i = load()
	if code, why := e.s.fieldsHold(ctx, g, i); code != CodeDefinitionInvalid || !strings.Contains(why, "refund_cd34 needs changing: PostRefunds takes no field amount_to_counter") {
		t.Fatalf("a made-up field: %q %q", code, why)
	}
	if g, i = load(); g.Tools[i].FieldsChecked {
		t.Fatal("a tool that doesn't hold was marked checked")
	}

	// Its own fields: it holds, is marked, and the next call reads nothing.
	g, i = load()
	delete(g.Tools[i].Arguments, "amount_to_counter")
	g.Tools[i].FieldsChecked = false
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	g, i = load()
	if code, why := e.s.fieldsHold(ctx, g, i); code != "" || why != "" {
		t.Fatalf("its own fields: %q %q", code, why)
	}
	g, i = load()
	if !g.Tools[i].FieldsChecked {
		t.Fatal("a tool that holds wasn't marked")
	}
	reads := strings.Count(string(mustJSON(g.Events)), "read PostRefunds's fields")
	if code, _ := e.s.fieldsHold(ctx, g, i); code != "" {
		t.Fatalf("checked: %q", code)
	}
	if g, _ = load(); strings.Count(string(mustJSON(g.Events)), "read PostRefunds's fields") != reads {
		t.Fatal("a checked tool was read again")
	}
}

// Setup's listing of a taf's grants says which tools are Setup's own reads,
// so a check that the taf's Card covers its tools can leave them out.
func TestSetupsListingMarksItsOwnReads(t *testing.T) {
	e := newEnv(t)
	id, out := e.addToDetails(map[string]any{"tool": "create_refund", "name": "refund_cd34", "amount_argument": "amount", "amount_unit": "EUR", "amount_scale": 2})
	e.approveAddition(id, out)
	st, list := e.do("GET", "/v1/accounts/"+tenant+"/agents/agt_taf/grants", nil, nil, "k_setup")
	gs, _ := list["grants"].([]any)
	if st != 200 || len(gs) != 1 {
		t.Fatalf("Setup's list: %d %v", st, list)
	}
	marked := map[string]any{}
	for _, raw := range gs[0].(map[string]any)["tools"].([]any) {
		tool := raw.(map[string]any)
		marked[tool["name"].(string)] = tool["setup_reads"]
	}
	if marked["fields_ab12"] != true || marked["refund_cd34"] != nil || len(marked) != 2 {
		t.Fatalf("the tools as Setup reads them: %v", marked)
	}
}
