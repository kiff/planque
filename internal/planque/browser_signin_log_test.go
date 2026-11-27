package planque

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
)

// What the owner types in the live sign-in view passes through planque to
// the browser and is kept nowhere: no log line at any level, no stored
// record, and no answer of planque's carries it. That covers typing that
// works, typing into a view that has ended, a malformed act that names it,
// and the view's end.
func TestTheOwnersTypingInTheLiveViewIsNeverLoggedOrKept(t *testing.T) {
	const secret = "correct-horse-battery-staple-7731"
	e := newBrowserEnv(t)
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	g := e.connect()
	st, out := e.startView(g.ID, map[string]any{"url": fakes.ShopLogin})
	if st != 200 {
		t.Fatalf("start: %d %v", st, out)
	}
	viewURL := out["view"].(string)
	var answers []string
	for _, in := range []any{
		map[string]any{"kind": "tap", "x": 10, "y": 20},
		map[string]any{"kind": "type", "text": secret},
		map[string]any{"kind": "key", "key": "Enter"},
		// Malformed: an unknown field beside the text, and a kind that isn't one.
		map[string]any{"kind": "type", "text": secret, "extra": secret},
		map[string]any{"kind": "typed", "text": secret},
	} {
		_, _, raw := e.view("POST", viewURL+"/input", in)
		answers = append(answers, string(raw))
	}
	if !strings.Contains(mustJSONString(e.b.Inputs), secret) {
		t.Fatal("the typing didn't reach the browser")
	}
	// The view ends, then the owner types again into the finished view.
	_, _, raw := e.view("POST", viewURL+"/done", nil)
	answers = append(answers, string(raw))
	_, _, raw = e.view("POST", viewURL+"/input", map[string]any{"kind": "type", "text": secret})
	answers = append(answers, string(raw))

	if strings.Contains(logs.String(), secret) {
		t.Fatalf("a log line holds the owner's typing:\n%s", logs.String())
	}
	for i, a := range answers {
		if strings.Contains(a, secret) {
			t.Errorf("answer %d holds the owner's typing: %s", i, a)
		}
	}
	items, err := e.kv.Query(context.Background(), tenantPK(tenant), "")
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range items {
		if bytes.Contains(it.Doc, []byte(secret)) {
			t.Errorf("record %s holds the owner's typing", it.SK)
		}
	}
	if logs.Len() == 0 {
		t.Fatal("nothing was logged at all: the test reads no log")
	}
}
