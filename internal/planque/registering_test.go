package planque

import (
	"encoding/json"
	"net/http"
	"testing"
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
