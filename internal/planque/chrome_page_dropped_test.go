package planque

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// A step the Chrome signed as ran, with no page and page_dropped, is done with
// no page; page_dropped never turns any other answer into one.
func TestADroppedPageIsAStepThatRanAndNothingElse(t *testing.T) {
	read := func(fields map[string]string) error {
		body := map[string]json.RawMessage{}
		for k, v := range fields {
			body[k] = json.RawMessage(v)
		}
		_, err := readChromeResult(body)
		return err
	}
	if err := read(map[string]string{"ok": "true", "page": "null", "page_dropped": `"page_too_large"`}); !errors.Is(err, errPageDropped) {
		t.Errorf("ran, page dropped: %v", err)
	}
	for name, fields := range map[string]map[string]string{
		"another reason":     {"ok": "true", "page": "null", "page_dropped": `"other"`},
		"a page with no url": {"ok": "true", "page": `{"title":"t"}`, "page_dropped": `"page_too_large"`},
		"no page_dropped":    {"ok": "true", "page": "null"},
	} {
		var r *chromeRefusal
		if err := read(fields); errors.Is(err, errPageDropped) || !errors.As(err, &r) || r.Code != browser.CodeBrowserDown {
			t.Errorf("%s: %v", name, err)
		}
	}
	var r *chromeRefusal
	if err := read(map[string]string{"ok": "false", "page": "null", "error": `{"code":"not_visible"}`, "page_dropped": `"page_too_large"`}); !errors.As(err, &r) || r.Code != "not_visible" {
		t.Errorf("a refusal: %v", err)
	}
	if err := read(map[string]string{"ok": "true", "page": `{"url":"https://shop.example/","read_id":"rd_1"}`, "page_dropped": `"page_too_large"`}); err != nil {
		t.Errorf("a page: %v", err)
	}
	// A hand-over that ran keeps its answer: handed.
	if err := readHandedResult(map[string]json.RawMessage{"ok": json.RawMessage("true"), "page": json.RawMessage("null"), "handed": json.RawMessage("true"), "page_dropped": json.RawMessage(`"page_too_large"`)}); err != nil {
		t.Errorf("hand: %v", err)
	}
	// page_dropped is part of what the Chrome signed.
	body, ok := resultBody(map[string]json.RawMessage{"id": nil, "nonce": nil, "chrome_id": nil, "ok": nil, "page": nil, "error": nil, "page_dropped": json.RawMessage(`"page_too_large"`)})
	if !ok || string(body["page_dropped"]) != `"page_too_large"` {
		t.Errorf("signed body: %v", body)
	}
}

// An act that ran in the owner's Chrome, whose page was too large to send, is
// done: the taf is told it ran and to read again, the steps after it aren't
// sent, nothing is read for it, and its read isn't the latest one.
func TestAnActThatRanWithItsPageDroppedIsDoneNeverAFailure(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t)
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if step["op"] == "act" {
			return true, nil, "page_dropped"
		}
		return true, plainPage("rd_1"), ""
	}
	c.online(t)
	_, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	read := sc["read"].(string)
	before := len(c.sent())
	state, res := c.callState(g, browseAct, map[string]any{"url": "https://shop.example/", "read": read,
		"steps": []any{map[string]any{"click": "c1", "label": "Search"}, map[string]any{"fill": "f2", "value": "again", "label": "Search"}}})
	out, sc := textOf(res)
	if state != OpDone || sc["steps_run"] != float64(1) || sc["steps_not_sent"] != float64(1) || sc["page_dropped"] != "page_too_large" {
		t.Fatalf("dropped: %s %v", state, sc)
	}
	if isErr, _ := res.(map[string]any)["isError"].(bool); isErr {
		t.Errorf("a step that ran is an error: %s", out)
	}
	for _, want := range []string{"Step 1 ran", "Don't repeat it", "browse_open", "(1)"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in %q", want, out)
		}
	}
	if strings.Contains(out, "failed") || strings.Contains(out, "page-words") {
		t.Errorf("words: %q", out)
	}
	// The call's read of its page, one act, then nothing: no second act, and
	// no read or open for the dropped page.
	var ops []string
	for _, s := range c.sent()[before:] {
		ops = append(ops, s["op"].(string))
	}
	acts := 0
	for _, op := range ops {
		switch {
		case op == "act":
			acts++
		case acts > 0 && (op == "read" || op == "open"):
			t.Errorf("a %s after the dropped page: %v", op, ops)
		}
	}
	if acts != 1 {
		t.Errorf("acts sent: %v", ops)
	}
	// The read the act was planned from is still the latest one kept: no page
	// after the act was read or kept.
	latest, err := c.s.latestChromeRead(t.Context(), g)
	if err != nil || latest == nil || latest.Read != read {
		t.Errorf("latest read: %+v %v", latest, err)
	}
}

// A Chrome from before page_dropped answers an act that ran, whose page fit no
// cut, page_too_large: the step may have run, so it's unknown, never "failed".
func TestAnActAnsweredPageTooLargeMayHaveRun(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t)
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if step["op"] == "act" {
			return false, nil, "page_too_large"
		}
		return true, plainPage("rd_1"), ""
	}
	c.online(t)
	_, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	read := sc["read"].(string)
	state, res := c.callState(g, browseAct, map[string]any{"url": "https://shop.example/", "read": read, "steps": []any{map[string]any{"click": "c1", "label": "Search"}}})
	out, _ := textOf(res)
	if state != OpUnknown || !strings.Contains(out, "Step 1 may have run") || !strings.Contains(out, "Don't repeat it") || strings.Contains(out, "failed") {
		t.Fatalf("old build: %s %q", state, out)
	}
	if strings.Contains(out, "ignore planque") {
		t.Errorf("the Chrome's detail reached the taf: %q", out)
	}
}
