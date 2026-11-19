package planque

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/words"
)

// Each time Coach asks again for an app, the earlier request from the same
// agent that still waits for the owner ends, so one row waits at most: the
// phone read seven "Waiting for your approval" rows after a night of retries.
// A request from another agent, for another app, or one that holds a sign-in
// is left as it was; an ended or failed grant never reads as waiting.
func TestANewRequestEndsTheSameAgentsEarlierPendingOne(t *testing.T) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	e.s.Applications = shippedApps(t)
	credentials := &appcredentials.Sealed{KV: e.kv, Sealer: e.s.Sealer}
	for _, app := range []string{"googlesheets"} {
		if err := credentials.Put(ctx, app, appcredentials.Credentials{ClientID: "demo-client", Secret: "demo-secret"}); err != nil {
			t.Fatal(err)
		}
	}
	e.s.ApplicationCredentials = credentials
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))
	raw, err := os.ReadFile("../../deploy/catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	if e.s.CatalogEntries, err = DecodeCatalog(raw, e.s.Applications); err != nil {
		t.Fatal(err)
	}
	ask := func(entry, agent string) string {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/catalog/"+entry+"/grants", map[string]any{"agent": agent}, nil, "k_setup")
		if st != 201 {
			t.Fatalf("%s for %s: %d %v", entry, agent, st, out)
		}
		return out["grant"].(map[string]any)["id"].(string)
	}
	load := func(id string) *Grant {
		g, err := e.s.loadGrant(ctx, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
		return g
	}

	first := ask("googlesheets@1", "agt_orders")
	other := ask("googlesheets@1", "agt_other")
	if load(first).Status != StatusPending || load(other).Status != StatusPending {
		t.Fatal("the first requests aren't pending")
	}
	second := ask("googlesheets@1", "agt_orders")

	if g := load(first); g.Status != StatusDisconnected || g.Reason != "" || !strings.Contains(g.Events[len(g.Events)-1].What, "replaced by a newer request") || !strings.Contains(g.Events[len(g.Events)-1].What, second) {
		t.Errorf("the earlier request: %s %q %v", g.Status, g.Reason, g.Events)
	}
	if !strings.Contains(logs.String(), `"msg":"grant ended"`) || !strings.Contains(logs.String(), `"from":"pending"`) || !strings.Contains(logs.String(), `"by":"a newer request, `+second+`"`) {
		t.Errorf("the replacement isn't logged:\n%s", logs.String())
	}
	if g := load(other); g.Status != StatusPending {
		t.Errorf("another agent's request ended: %s", g.Status)
	}
	if g := load(second); g.Status != StatusPending {
		t.Errorf("the new request is %s", g.Status)
	}

	// One that holds a sign-in isn't ended by a newer ask.
	signed := load(second)
	signed.SignedIn = true
	if err := e.s.saveGrant(ctx, signed); err != nil {
		t.Fatal(err)
	}
	third := ask("googlesheets@1", "agt_orders")
	if g := load(second); g.Status != StatusPending || !g.SignedIn {
		t.Errorf("a request holding a sign-in was ended: %s", g.Status)
	}
	if load(third).Status != StatusPending {
		t.Error("the newest isn't pending")
	}

	// Another app from the same agent isn't ended.
	elsewhere := ask("googlesheets@1", "agt_two_apps")
	moved := load(elsewhere)
	moved.Request.CatalogEntry = "connector:composio:other"
	if err := e.s.saveGrant(ctx, moved); err != nil {
		t.Fatal(err)
	}
	ask("googlesheets@1", "agt_two_apps")
	if g := load(elsewhere); g.Status != StatusPending {
		t.Errorf("a request for another app was ended: %s", g.Status)
	}

	// A grant that ended or failed never reads as waiting, in any language.
	for _, lang := range []string{"en", "fr", "es"} {
		w := words.Approve(lang)
		ended := load(first)
		if got := statusIn(w, ended); got == w.WaitingForYou || got != w.Disconnected {
			t.Errorf("%s: an ended request reads %q", lang, got)
		}
		failed := *ended
		failed.Status = StatusFailed
		if got := statusIn(w, &failed); got == w.WaitingForYou || got != w.NotConnected {
			t.Errorf("%s: a failed request reads %q", lang, got)
		}
	}
}
