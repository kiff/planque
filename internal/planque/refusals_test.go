package planque

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

// A refused connection request is logged with planque's words, so a taf's
// failed step can be read from planque's log; an app whose set-up fails at
// the provider has the provider's cause logged too, its key never (13 Nov
// 2026: Coach said "can't be set up right now" and no one could say why).
func TestARefusedConnectionIsLoggedWithItsReason(t *testing.T) {
	e, f := newConnectorFixture(t)
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))

	// An app that waits for an operator.
	st, _ := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:own_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 409 {
		t.Fatalf("an app that needs setting up: %d", st)
	}
	got := logs.String()
	if !strings.Contains(got, `"msg":"connection refused"`) || !strings.Contains(got, `"status":409`) ||
		!strings.Contains(got, `"entry":"connector:composio:own_demo"`) || !strings.Contains(got, "needs setting up first") {
		t.Fatalf("the refusal isn't logged with its reason: %s", got)
	}

	// An app set up on the owner's pick, failing at the provider.
	logs.Reset()
	f.setUpFails = true
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:ondemand_1/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup")
	if st != 503 || !strings.Contains(out["error"].(string), "can't be set up right now") {
		t.Fatalf("a set-up the provider refused: %d %v", st, out)
	}
	got = logs.String()
	if !strings.Contains(got, `"msg":"app set-up failed"`) || !strings.Contains(got, `"toolkit":"ondemand_1"`) || !strings.Contains(got, "no managed auth") {
		t.Fatalf("the provider's cause isn't logged: %s", got)
	}
	if strings.Contains(got, "synthetic-project-key") {
		t.Fatalf("the project key reached the log: %s", got)
	}
	if !strings.Contains(got, `"msg":"connection refused"`) || !strings.Contains(got, `"status":503`) {
		t.Fatalf("the refusal isn't logged: %s", got)
	}

	// A request that goes through logs nothing of the kind.
	logs.Reset()
	f.setUpFails = false
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/connector:composio:desk_demo/grants", map[string]string{"agent": "agt_taf"}, nil, "k_setup"); st != 201 {
		t.Fatalf("a request that goes through: %d %v", st, out)
	}
	if strings.Contains(logs.String(), "connection refused") {
		t.Fatalf("a request that went through was logged as refused: %s", logs.String())
	}
}

// A description address the owner pasted can carry a token in its query: the
// answer, and so the log, say the address without it, whether the vendor
// answered or the call never got through.
func TestARefusedDescriptionAddressKeepsItsQueryOut(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))
	for _, u := range []string{v.URL + "/missing?token=s3cret-in-query#also-s3cret", "http://127.0.0.1:1/openapi.json?token=s3cret-in-query"} {
		logs.Reset()
		api := map[string]any{"url": u, "auth": map[string]any{"in": "header", "name": "X-API-Key"}}
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []any{}, "api": api}, nil, "k_setup")
		answer := mustJSONString(out)
		if st != 400 || !strings.Contains(answer, "Couldn't read the API description at") || strings.Contains(answer, "s3cret") {
			t.Fatalf("%s: %d %s", u, st, answer)
		}
		if strings.Contains(logs.String(), "s3cret") || !strings.Contains(logs.String(), "connection refused") {
			t.Fatalf("%s: the log: %s", u, logs.String())
		}
	}
}
