package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

var listReads int

// listRow is grant id as the owner's list shows it, read as an app reads
// it: signed by the phone, in e's language, in zone (none when "").
func listRow(t *testing.T, e *env, id, zone string) map[string]any {
	t.Helper()
	// Each read is its own signed request (a fresh URI), never a replay of the last.
	listReads++
	req, _ := http.NewRequest("GET", e.url+"/v1/accounts/"+tenant+"/grants?read="+strconv.Itoa(listReads), bytes.NewReader(nil))
	e.phone.Sign(req, tenant, nil)
	if e.lang != "" {
		req.Header.Set("Accept-Language", e.lang)
	}
	if zone != "" {
		req.Header.Set("X-Time-Zone", zone)
	}
	client := e.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out struct {
		Grants []map[string]any `json:"grants"`
	}
	if resp.StatusCode != 200 || json.NewDecoder(resp.Body).Decode(&out) != nil {
		t.Fatalf("list: %d", resp.StatusCode)
	}
	for _, g := range out.Grants {
		if g["id"] == id {
			return g
		}
	}
	t.Fatalf("%s isn't listed", id)
	return nil
}

// actionIDs are a row's actions' ids, in order.
func actionIDs(row map[string]any) []string {
	out := []string{}
	for _, a := range row["actions"].([]any) {
		out = append(out, a.(map[string]any)["id"].(string))
	}
	return out
}

func action(row map[string]any, id string) map[string]any {
	for _, a := range row["actions"].([]any) {
		if m := a.(map[string]any); m["id"] == id {
			return m
		}
	}
	return nil
}

func health(row map[string]any) map[string]any { return row["health"].(map[string]any) }

// seen is what a person reads of a row: title, used by, health and actions.
func seen(row map[string]any) string {
	return mustJSONString(map[string]any{"title": row["title"], "used_by": row["used_by"], "health": row["health"], "actions": row["actions"]})
}

// The owner's leftover: a connection signed in with no tools chosen and no
// taf behind it offers only to remove it, says it's unfinished, and never
// offers Sign in again, which would answer "this connection is listed".
func TestAnUnfinishedConnectionOffersOnlyToRemoveIt(t *testing.T) {
	e, f := newConnectorFixture(t)
	id := listConnector(t, e, f)
	row := listRow(t, e, id, "")
	if got := actionIDs(row); len(got) != 1 || got[0] != ActionRemoveUnfinished {
		t.Fatalf("actions: %v", got)
	}
	if a := action(row, ActionRemoveUnfinished); a["words"] != "Remove unfinished connection" || a["how"] != "delete" {
		t.Fatalf("remove unfinished: %v", a)
	}
	if h := health(row); h["state"] != HealthUnfinished || h["ok"] != false || !strings.HasPrefix(h["words"].(string), "Unfinished · signed in on ") || !strings.Contains(h["words"].(string), "no taf uses it") {
		t.Fatalf("health: %v", h)
	}
	if u := row["used_by"].([]any); len(u) != 0 {
		t.Fatalf("used by: %v", u)
	}
	// The provider's id for the connection is in no field a person reads.
	if s := seen(row); strings.Contains(s, "account_demo") || strings.Contains(s, "ca_") {
		t.Fatalf("a provider id is shown: %s", s)
	}
}

// A working connection names the taf on it, says it works and when it was
// last used (in the owner's zone), and offers only Remove access, whose
// confirm names what stops. Before any call it says it isn't used yet.
func TestAWorkingConnectionSaysWhoUsesItAndWhenAndOffersOnlyRemove(t *testing.T) {
	e, f := newConnectorFixture(t)
	e.tafName = "Unanswered customer emails"
	g := connectConnector(t, e, f)
	taf, vendor := g.Request.AgentName, vendorName(g)
	if taf == "" {
		t.Fatal("the fixture's taf has no name")
	}
	row := listRow(t, e, g.ID, "Europe/Paris")
	if u := row["used_by"].([]any); len(u) != 1 || u[0] != taf {
		t.Fatalf("used by: %v", u)
	}
	if got := actionIDs(row); len(got) != 1 || got[0] != ActionRemove {
		t.Fatalf("actions: %v", got)
	}
	if a := action(row, ActionRemove); a["words"] != "Remove access" || a["how"] != "delete" || a["stops"] != taf+" will stop using "+vendor+"." {
		t.Fatalf("remove: %v", a)
	}
	if h := health(row); h["state"] != HealthWorking || h["ok"] != true || h["words"] != "Working · not used yet" {
		t.Fatalf("health before a call: %v", h)
	}

	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_used"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_used", Args: args})
	if got := mustJSONString(e.s.call(context.Background(), g, "reply_ticket", args, permit)); strings.Contains(got, `"isError":true`) {
		t.Fatalf("the call: %s", got)
	}
	paris, _ := time.LoadLocation("Europe/Paris")
	want := "Working · last used " + e.s.now().In(paris).Format("15:04")
	if h := health(listRow(t, e, g.ID, "Europe/Paris")); h["words"] != want || h["at"] == nil {
		t.Fatalf("health after a call: %v, want %q", h, want)
	}
	// No zone named: the time is said in UTC, and says so.
	if h := health(listRow(t, e, g.ID, "")); h["words"] != "Working · last used "+e.s.now().UTC().Format("15:04")+" UTC" {
		t.Fatalf("health with no zone: %v", h)
	}
}

// A call that fails isn't a use: the last use stays as it was.
func TestAFailedCallIsntALastUse(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	f.executeBody = `{"data":{"status_code":500,"message":"down"},"successful":false,"error":"vendor down"}`
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_fail"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_fail", Args: args})
	_ = e.s.call(context.Background(), g, "reply_ticket", args, permit)
	if h := health(listRow(t, e, g.ID, "")); h["words"] != "Working · not used yet" {
		t.Fatalf("health after a failed call: %v", h)
	}
}

// A signed-out managed connection says when its sign-in expired and offers
// Sign in again on itself (planque's approve) before Remove access; in the
// owner's language too.
func TestASignedOutConnectionOffersSignInAgainAndSaysWhen(t *testing.T) {
	e, f := newConnectorFixture(t)
	e.tafName = "Unanswered customer emails"
	g := connectConnector(t, e, f)
	signOut401(t, e, f, g)
	vendor := vendorName(g)
	row := listRow(t, e, g.ID, "Europe/Paris")
	if got := actionIDs(row); len(got) != 2 || got[0] != ActionSignInAgain || got[1] != ActionRemove {
		t.Fatalf("actions: %v", got)
	}
	if a := action(row, ActionSignInAgain); a["words"] != "Sign in again" || a["how"] != "approve" {
		t.Fatalf("sign in again: %v", a)
	}
	h := health(row)
	if h["state"] != HealthSignedOut || h["ok"] != false || !strings.HasPrefix(h["words"].(string), vendor+"'s sign-in expired on ") || h["at"] == nil {
		t.Fatalf("health: %v", h)
	}
	if u := row["used_by"].([]any); len(u) != 1 || u[0] != g.Request.AgentName {
		t.Fatalf("used by: %v", u)
	}

	e.lang = "fr"
	row = listRow(t, e, g.ID, "Europe/Paris")
	if a := action(row, ActionSignInAgain); a["words"] != "Se reconnecter" {
		t.Fatalf("fr sign in again: %v", a)
	}
	if w := health(row)["words"].(string); !strings.HasPrefix(w, "La connexion à "+vendor+" a expiré le ") {
		t.Fatalf("fr health: %s", w)
	}
}

// A signed-out managed connection whose account's name was never kept
// can't sign in again on itself (no new sign-in could be shown to be the
// same account), so only Remove access is offered.
func TestASignedOutConnectionWithNoAccountNameOffersOnlyRemove(t *testing.T) {
	e, f := newConnectorFixture(t)
	f.noName = true
	g := connectConnector(t, e, f)
	signOut401(t, e, f, g)
	if got := actionIDs(listRow(t, e, g.ID, "")); len(got) != 1 || got[0] != ActionRemove {
		t.Fatalf("actions: %v", got)
	}
}
