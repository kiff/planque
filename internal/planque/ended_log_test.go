package planque

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// A grant that ends is logged: which, from which status to which, and who or
// what ended it, so a row that went from "waiting for your approval" to gone
// can be explained (the owner saw seven such rows and nothing said who ended
// them). Three ways to end, three lines; nothing in a line is a secret.
func TestAGrantThatEndsIsLoggedWithWhoEndedIt(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))
	ended := func(id string) map[string]any {
		t.Helper()
		for _, line := range strings.Split(logs.String(), "\n") {
			var m map[string]any
			if json.Unmarshal([]byte(line), &m) == nil && m["msg"] == "grant ended" && m["grant"] == id {
				return m
			}
		}
		t.Fatalf("no \"grant ended\" line for %s in:\n%s", id, logs.String())
		return nil
	}
	request := func() (string, string) {
		_, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"server": e.vendor.MCPURL(), "agent": "agt_taf", "tools": []any{}}, nil, "k_setup")
		return out["id"].(string), out["digest"].(string)
	}

	// The owner's signed disconnect, from a connected-or-waiting status.
	id, digest := request()
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, e.phone, ""); st != 200 {
		t.Fatalf("approve: %d", st)
	}
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id, nil, e.phone, ""); st != 200 {
		t.Fatalf("disconnect: %d", st)
	}
	if m := ended(id); m["from"] != StatusApproved || m["to"] != StatusDisconnected || m["by"] != "the owner's signed disconnect" || m["agent"] != "agt_taf" {
		t.Errorf("the disconnect line: %v", m)
	}

	// Nobody signed it in time: planque's sweep ends it.
	waiting, _ := request()
	e.at(pendingTTL + time.Minute)
	if n, err := e.s.Sweep(t.Context()); err != nil || n < 1 {
		t.Fatalf("the sweep: %d %v", n, err)
	}
	if m := ended(waiting); m["from"] != StatusPending || m["to"] != StatusFailed || m["by"] != "planque's sweep" {
		t.Errorf("the sweep line: %v", m)
	}
}
