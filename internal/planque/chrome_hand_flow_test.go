package planque

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/stepkey"
)

// The hand-over's tools are offered only while ChromeHand is on, and only with
// the owner's Chrome: asked for otherwise, the grant is refused in words.
func TestTheHandOverIsOfferedOnlyWithTheOwnersChromeOnceItsOn(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t)
	if slicesHas(g.Tools, browseHand) || slicesHas(g.Tools, browseHandResult) {
		t.Fatal("hand tools granted while off")
	}
	chrome := map[string]any{"driver": DriverExtension, "chrome": c.r.chrome, "sites": []string{"shop.example"}}
	st, out := c.ask(map[string]any{"browser": chrome, "tools": []any{map[string]any{"tool": browseOpen}, map[string]any{"tool": browseHand}}})
	if st != 400 || !strings.Contains(out["error"].(string), "isn't offered yet") {
		t.Fatalf("asked while off: %d %v", st, out)
	}
	c.s.ChromeHand = true
	g = c.grant(t)
	if !slicesHas(g.Tools, browseHand) || !slicesHas(g.Tools, browseHandResult) {
		t.Fatalf("not in a Chrome grant's defaults: %v", g.Tools)
	}
	hosted := c.connect()
	if slicesHas(hosted.Tools, browseHand) || slicesHas(hosted.Tools, browseHandResult) {
		t.Fatal("hand tools in a hosted grant")
	}
	st, out = c.ask(map[string]any{"browser": map[string]any{"sites": []string{"shop.example"}}, "tools": []any{map[string]any{"tool": browseHandResult}}})
	if st != 400 || !strings.Contains(out["error"].(string), "owner's Chrome's only") {
		t.Fatalf("asked for a hosted browser: %d %v", st, out)
	}
}

// handRig is a Chrome rig with the hand-over on, its grant, and the read of
// shop.example the taf hands a step on.
func handRig(t *testing.T) (*chromeRig, *Grant, string) {
	t.Helper()
	c := newChromeRig(t)
	c.s.ChromeHand = true
	g := c.grant(t)
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		return true, shopPage("rd_1", nil), ""
	}
	c.online(t)
	_, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	read, _ := sc["read"].(string)
	if read == "" {
		t.Fatalf("open: %v", sc)
	}
	return c, g, read
}

// answerHandAs is the extension's signed answer to a hand-over, sent with key.
func (c *chromeRig) answerHandAs(key stepkey.Local, hand, outcome string, page map[string]any) map[string]any {
	c.t.Helper()
	body := map[string]any{"hand": hand, "outcome": outcome}
	if page != nil {
		body["page"] = page
	}
	sig, err := stepkey.SignBody(context.Background(), key, labelHand, body)
	if err != nil {
		c.t.Fatal(err)
	}
	msg := map[string]any{"type": "hand_result", "sig": sig}
	for k, v := range body {
		msg[k] = v
	}
	c.r.event("$default", "c1", msg)
	return c.r.poster.last("c1")
}

// A hand-over's life end to end: a signed "hand" step on the taf's read, the
// run left open (no close), waiting until the person answers, the answer
// kept once from this Chrome's key, then read once with no field values.
func TestAHandOverGoesToThePersonAndTheirAnswerIsReadOnce(t *testing.T) {
	c, g, read := handRig(t)
	_, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": read, "words": "Type your password, then tap Done", "ref": "f1"})
	hand, _ := sc["hand"].(string)
	if sc["state"] != "handed" || hand == "" {
		t.Fatalf("hand: %v", sc)
	}
	var step map[string]any
	waitFor(t, func() bool {
		for _, s := range c.sent() {
			if s["op"] == "hand" {
				step = s
			}
		}
		return step != nil
	})
	h, _ := step["hand"].(map[string]any)
	if step["id"] != hand || h["words"] != "Type your password, then tap Done" || h["ref"] != "f1" || step["read_id"] != "rd_1" {
		t.Fatalf("the step: %v", step)
	}
	// The run stays open for the person: no close follows the hand step.
	time.Sleep(100 * time.Millisecond)
	if ops := c.ops(); slices.Contains(ops[slices.Index(ops, "hand"):], "close") {
		t.Fatalf("the run was closed after the hand-over: %v", ops)
	}
	if _, sc := c.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handWaiting {
		t.Fatalf("before the person: %v", sc)
	}
	ok, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	other := stepkey.Local{Key: ok}
	if m := c.answerHandAs(other, hand, handDone, nil); m["type"] != "hand_refused" || m["code"] != "unsigned" {
		t.Fatalf("another key's answer: %v", m)
	}
	if m := c.answerHandAs(c.r.ext, "st_nobody", handDone, nil); m["type"] != "hand_refused" || m["code"] != "unknown_hand" {
		t.Fatalf("an unknown hand-over: %v", m)
	}
	after := shopPage("rd_2", nil)
	if m := c.answerHandAs(c.r.ext, hand, handDone, after); m["type"] != "hand_kept" || m["hand"] != hand {
		t.Fatalf("the answer: %v", m)
	}
	if m := c.answerHandAs(c.r.ext, hand, handDeclined, nil); m["type"] != "hand_kept" {
		t.Fatalf("a resend: %v", m)
	}
	out, sc := c.call(g, browseHandResult, map[string]any{"hand": hand})
	if sc["state"] != handDone || sc["read"] == nil || strings.Contains(out, "hunter2") {
		t.Fatalf("the answer read: %s", out)
	}
	if _, sc := c.call(g, browseHandResult, map[string]any{"hand": hand}); code(sc) != CodeBadArguments {
		t.Fatalf("read twice: %v", sc)
	}
}

// A hand-over is refused before anything reaches the Chrome when its words,
// read or ref aren't right; its answer is read only by its own grant.
func TestABadHandOverIsNeverSent(t *testing.T) {
	c, g, read := handRig(t)
	for name, args := range map[string]map[string]any{
		"words": {"url": "https://shop.example/", "read": read, "words": "Pay\nnow"},
		"read":  {"url": "https://shop.example/", "read": "rd_another", "words": "Tap Done"},
		"ref":   {"url": "https://shop.example/", "read": read, "words": "Tap Done", "ref": "f9"},
	} {
		if _, sc := c.call(g, browseHand, args); code(sc) == "" {
			t.Errorf("%s: %v", name, sc)
		}
	}
	waitFor(t, func() bool { return len(c.sent()) > 0 })
	for _, op := range c.ops() {
		if op == "hand" {
			t.Fatalf("a refused hand-over was sent: %v", c.ops())
		}
	}
	_, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": read, "words": "Tap Done"})
	hand, _ := sc["hand"].(string)
	g2 := c.grant(t)
	name := ""
	for _, gt := range g2.Tools {
		if gt.Tool == browseHandResult {
			name = gt.Name // a second grant's tools are named apart
		}
	}
	if _, sc := c.call(g2, name, map[string]any{"hand": hand}); code(sc) != CodeBadArguments {
		t.Fatalf("another grant read it: %v", sc)
	}
	if _, sc := c.call(g, browseHandResult, map[string]any{"hand": hand, "url": "https://shop.example/"}); code(sc) != CodeBadArguments {
		t.Fatalf("a result with a url: %v", sc)
	}
	if _, sc := c.call(g, browseHandResult, map[string]any{"hand": hand}); sc["state"] != handWaiting {
		t.Fatalf("the grant's own read: %v", sc)
	}
}

// A page whose read changes between loads (a news item, a rotating link)
// still takes a hand-over on the latest read: the same address, and the
// element handed unchanged by ref, kind and label. A changed element, another
// address, or a read that isn't the latest is still refused.
func TestAHandOverLandsOnAPageWhoseReadChangedButNotOnAnotherElement(t *testing.T) {
	c, g, read := handRig(t)
	// The rig reads its answer under c.mu while the last call's close is still
	// being answered: swap it under the same lock.
	news := func(label string, fieldLabel string) {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.answer = func(step map[string]any) (bool, map[string]any, string) {
			p := shopPage("rd_"+step["id"].(string), map[string]any{
				"links": []any{map[string]any{"ref": "l1", "label": "Lamp", "url": "https://shop.example/items/lamp"},
					map[string]any{"ref": "l2", "label": label, "url": "https://shop.example/news/today"}}})
			if fieldLabel != "" {
				p["fields"] = []any{map[string]any{"ref": "f1", "label": fieldLabel, "kind": "password"}, map[string]any{"ref": "f2", "label": "Search", "kind": "search"}}
			}
			return true, p, ""
		}
	}
	news("Today's news", "")
	_, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": read, "words": "Type your password, then tap Done", "ref": "f1"})
	if sc["state"] != "handed" {
		t.Fatalf("a page whose links changed: %v", sc)
	}
	// The latest read is now the page just opened: the first read is no
	// longer it, so it's refused, and the one this call gave is taken.
	news("Tomorrow's news", "")
	if _, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": read, "words": "Tap Done", "ref": "f1"}); code(sc) != browser.CodePageChanged {
		t.Fatalf("a read that isn't the latest: %v", sc)
	}
	_, sc = c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	latest, _ := sc["read"].(string)
	// The element handed changed: refused, whatever else holds.
	news("Later news", "Card number")
	if _, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": latest, "words": "Type it, then tap Done", "ref": "f1"}); code(sc) != browser.CodePageChanged {
		t.Fatalf("another element under the same ref: %v", sc)
	}
	_, sc = c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	latest, _ = sc["read"].(string)
	// A control is never taken for a field of the same ref.
	news("Even later news", "")
	if _, sc := c.call(g, browseHand, map[string]any{"url": "https://shop.example/", "read": latest, "words": "Tap Done", "ref": "c9"}); code(sc) == "" {
		t.Fatalf("a ref on neither read: %v", sc)
	}
}

// Invariant: a taf whose hand-over is refused for its read is told the page's read as it is now, and the open tool it actually has, so its next hand-over can land
// Category: correctness
// Attack: none (a run that failed for real: the hand-over's tools on a connection of their own, a page whose links change on every load, and a taf that named the fence's tag as the read, five times)
// Refused: the refusal names the current read and the sibling connection's own open tool; a hand-over naming that read is handed; every Chrome answer says its read first
func TestAHandOverRefusedForItsReadSaysTheReadToNameAndTheOpenToolTheTafHas(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeHand = true
	// One taf, one Chrome: its open and act on one connection, its hand-over on another, sharing their reads.
	opens := c.grantAs(t, "agt_wiki", browseOpen, browseAct)
	hands := c.grantAs(t, "agt_wiki", browseHand, browseHandResult)
	openName := ""
	for _, gt := range opens.Tools {
		if gt.Tool == browseOpen {
			openName = gt.Name
		}
	}
	if openName == "" || slicesHas(hands.Tools, browseOpen) {
		t.Fatalf("the rig's connections: %v / %v", opens.Tools, hands.Tools)
	}
	// A page that renders a link differently on every load, as a news site's front page does.
	n := 0
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		n++
		return true, shopPage("rd_"+strings.Repeat("a", n), map[string]any{"links": []any{map[string]any{"ref": "l1", "label": fmt.Sprintf("Today %d", n), "url": "https://shop.example/items/lamp"}}}), ""
	}
	c.online(t)
	out, sc := c.call(opens, browseOpen, map[string]any{"url": "https://shop.example/"})
	read, _ := sc["read"].(string)
	if read == "" || !strings.HasPrefix(out, readLine(read)+"\n") {
		t.Fatalf("the open's answer doesn't say its read first: %q", out)
	}
	// The taf names the fence's tag, a 24-hex value, as it did.
	tag := out[strings.Index(out, "<<<page-words:")+len("<<<page-words:") : strings.Index(out, ">>>")]
	words, sc := c.call(hands, browseHand, map[string]any{"url": "https://shop.example/", "read": tag, "words": "Type your search, then tap Done", "ref": "f1"})
	if code(sc) != browser.CodePageChanged {
		t.Fatalf("a hand-over naming the fence's tag: %v", sc)
	}
	var latest chromeRead
	if _, err := getDoc(context.Background(), c.s.KV, tenantPK(tenant), skChromeRead(hands), &latest); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(words, "reads as "+latest.Read+" now") || !strings.Contains(words, "naming read "+latest.Read) || !strings.Contains(words, openName) {
		t.Fatalf("the refusal doesn't give the read to name (%s) and the open tool the taf has (%s): %s", latest.Read, openName, words)
	}
	// Named as told, it lands, though the page loaded differently again.
	if _, sc := c.call(hands, browseHand, map[string]any{"url": "https://shop.example/", "read": latest.Read, "words": "Type your search, then tap Done", "ref": "f1"}); code(sc) != "" || sc["state"] != "handed" {
		t.Fatalf("a hand-over naming the read it was told: %v", sc)
	}
}
