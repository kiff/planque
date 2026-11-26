package planque

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

// A grant whose tools are still being registered in KIFF reads as connecting,
// never connected, to every reader outside planque: Setup's first read after
// the sign-in got KIFF's "Unknown tool" six times when it read the grant as
// connected while planque was still registering its tools. It's saved
// connected all the while, so the relay answers KIFF's probe, and reads
// connected once the last tool is in KIFF.
func TestConnectingUntilEveryToolIsInKIFF(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	read := func(path string) map[string]any {
		req, _ := http.NewRequest("GET", e.url+path, nil)
		req.Header.Set("Authorization", "Bearer k_setup")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return nil
		}
		defer resp.Body.Close()
		var out map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&out)
		return out
	}
	type seen struct{ tool, grant, words, agents, stored string }
	var during []seen
	e.kiff.OnConnect = func(name string) {
		g := read("/v1/accounts/" + tenant + "/grants/" + id)
		var listed string
		if gs, _ := read("/v1/accounts/" + tenant + "/agents/agt_taf/grants")["grants"].([]any); len(gs) == 1 {
			listed, _ = gs[0].(map[string]any)["status"].(string)
		}
		stored, err := e.s.loadGrant(t.Context(), tenant, id)
		if err != nil {
			t.Error(err)
			return
		}
		grant, _ := g["status"].(string)
		words, _ := g["status_words"].(string)
		during = append(during, seen{name, grant, words, listed, stored.Status})
	}
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if len(during) != 2 {
		t.Fatalf("KIFF registered %d tools, want 2: %+v", len(during), during)
	}
	for _, d := range during {
		// Stored connected, so the relay serves KIFF's probe; read as connecting by the grant, its words and the agent's list.
		if d.stored != StatusConnected || d.grant != StatusConnecting || d.agents != StatusConnecting || d.words != "Connecting · your tools are being set up in KIFF" {
			t.Errorf("while %s was being registered: %+v", d.tool, d)
		}
	}
	e.kiff.OnConnect = nil
	g := read("/v1/accounts/" + tenant + "/grants/" + id)
	if g["status"] != StatusConnected || g["status_words"] != "Connected" {
		t.Fatalf("after the last tool: %v %v", g["status"], g["status_words"])
	}
	stored, _ := e.s.loadGrant(t.Context(), tenant, id)
	if stored.Registering {
		t.Fatal("still registering once connected")
	}
}

// A registration KIFF refuses fails the grant, and it's never left reading as
// connecting.
func TestRegisteringEndsWithAFailure(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	e.kiff.RefuseConnectAfter = 1
	e.approveAndSignIn(id, digest)
	g, err := e.s.loadGrant(t.Context(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Status != StatusFailed || g.Registering || g.registering() {
		t.Fatalf("after a refused registration: %s, registering %v", g.Status, g.Registering)
	}
}

// A registration cut off mid-way (planque's run timed out or stopped between
// two tools) never leaves the grant reading as connecting for ever: past
// RegisterBy the sweep undoes it, what was registered leaves KIFF and the
// grant fails. Before that, the sweep leaves a registration under way alone.
func TestARegistrationCutOffIsUndoneByTheSweep(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	if n := len(e.kiff.Connections(tenant)); n != 2 {
		t.Fatalf("KIFF holds %d tools, want 2", n)
	}
	// As a run that stopped after registering: saved connected, still registering.
	ctx := context.Background()
	g, err := e.s.loadGrant(ctx, tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	g.Registering, g.RegisterBy = true, e.s.now().Add(registerTTL).UTC()
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	if g.deadline() != g.RegisterBy {
		t.Fatalf("the registration's deadline: %v", g.deadline())
	}

	e.at(registerTTL - time.Minute)
	if n, err := e.s.Sweep(ctx); err != nil || n != 0 {
		t.Fatalf("a sweep within the registration's time: %d %v", n, err)
	}
	if g, _ := e.s.loadGrant(ctx, tenant, id); !g.registering() || len(e.kiff.Connections(tenant)) != 2 {
		t.Fatalf("a registration under way was touched: %s registering %v", g.Status, g.Registering)
	}

	e.at(registerTTL + time.Minute)
	if n, err := e.s.Sweep(ctx); err != nil || n != 1 {
		t.Fatalf("the sweep past RegisterBy: %d %v", n, err)
	}
	g, _ = e.s.loadGrant(ctx, tenant, id)
	if g.Status != StatusFailed || g.Registering || !g.RegisterBy.IsZero() || !strings.Contains(g.Reason, "didn't finish") {
		t.Fatalf("after the sweep: %s, registering %v by %v, %q", g.Status, g.Registering, g.RegisterBy, g.Reason)
	}
	if n := len(e.kiff.Connections(tenant)); n != 0 {
		t.Fatalf("KIFF still holds %d of its tools", n)
	}
	if !g.deadline().IsZero() {
		t.Fatalf("a failed grant keeps a deadline: %v", g.deadline())
	}
}

// A registration that finishes leaves no deadline behind.
func TestAFinishedRegistrationHasNoDeadline(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, digest := e.request()
	e.approveAndSignIn(id, digest)
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Status != StatusConnected || g.Registering || !g.RegisterBy.IsZero() || !g.deadline().IsZero() {
		t.Fatalf("connected: %s, registering %v by %v, deadline %v", g.Status, g.Registering, g.RegisterBy, g.deadline())
	}
}
