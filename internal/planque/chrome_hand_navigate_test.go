package planque

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// navRig is a Chrome rig with the hand-over on whose shop has two pages: its
// start page, and the lamp's page a click on the start page's link (l1)
// navigates to. The lamp's page renders a link differently on every load, as a
// news site's front page does, so no two of its reads are the same.
func navRig(t *testing.T) (*chromeRig, *Grant) {
	t.Helper()
	c := newChromeRig(t)
	c.s.ChromeHand = true
	g := c.grant(t)
	// No control places an order: a page with one is a checkout, where a taf follows no link.
	searchOnly := []any{map[string]any{"ref": "c1", "label": "Search"}}
	loads := 0
	lamp := func() map[string]any {
		loads++
		return shopPage(fmt.Sprintf("rd_lamp_%d", loads), map[string]any{
			"url": "https://shop.example/items/lamp", "title": "Lamp", "controls": searchOnly,
			"links": []any{map[string]any{"ref": "l1", "label": fmt.Sprintf("Today %d", loads), "url": "https://shop.example/"}},
		})
	}
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		switch {
		case step["op"] == "act":
			return true, lamp(), ""
		case step["url"] == "https://shop.example/items/lamp":
			return true, lamp(), ""
		case step["op"] == "hand":
			return true, lamp(), ""
		}
		return true, shopPage("rd_start", map[string]any{"controls": searchOnly}), ""
	}
	c.online(t)
	return c, g
}

// Invariant: after a click navigates, the read the act's answer gives is the one a hand-over on the new page names, and it lands
// Category: correctness
// Attack: none (a run that failed for real: an open, a click on a link that navigated, then five hand-overs refused for their read, and the tool blocked for the run)
// Refused: a hand-over naming the page before the click is refused with the read to name now and the open tool the taf has; the act's own read, and the read the refusal gives, are each handed
func TestAHandOverAfterAClickThatNavigatesLandsOnTheReadTheClickGave(t *testing.T) {
	c, g := navRig(t)
	out, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	start, _ := sc["read"].(string)
	if start == "" || !strings.HasPrefix(out, readLine(start)+"\n") {
		t.Fatalf("the open doesn't say its read first: %q", out)
	}

	// The click on the start page's link navigates to the lamp's page.
	out, sc = c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": start, "steps": []any{map[string]any{"click": "l1"}}})
	if code(sc) != "" || sc["steps_run"] != float64(1) {
		t.Fatalf("the click: %s", out)
	}
	after, _ := sc["read"].(string)
	if after == "" || after == start || !strings.HasPrefix(out, readLine(after)+"\n") {
		t.Fatalf("the act's answer doesn't say the new page's read first (%s, was %s): %q", after, start, out)
	}
	page, _ := sc["page"].(map[string]any)
	if page["url"] != "https://shop.example/items/lamp" {
		t.Fatalf("the act's answer isn't the page the click landed on: %v", page)
	}

	// The taf hands over on the page the click landed on, naming the read the
	// click gave: it lands, though the page loads differently every time.
	hand := map[string]any{"url": "https://shop.example/items/lamp", "read": after, "words": "Type your search, then tap Done", "ref": "f2"}
	if _, sc := c.call(g, browseHand, hand); code(sc) != "" || sc["state"] != "handed" {
		t.Fatalf("a hand-over naming the click's read: %v", sc)
	}

	// Naming the page from before the click is the mistake the taf made: it's
	// refused with the read to name now and the open tool this taf has, so the
	// next hand-over lands without another call.
	stale := map[string]any{"url": "https://shop.example/items/lamp", "read": start, "words": "Type your search, then tap Done", "ref": "f2"}
	words, sc := c.call(g, browseHand, stale)
	if code(sc) != browser.CodePageChanged {
		t.Fatalf("a hand-over naming the start page's read: %v", sc)
	}
	latest, err := c.s.latestChromeRead(t.Context(), g)
	if err != nil || latest == nil {
		t.Fatalf("the latest read: %v %v", latest, err)
	}
	if !strings.Contains(words, "naming read "+latest.Read) || !strings.Contains(words, openNameOf(g)) || strings.Contains(words, "with  ") || strings.Contains(words, "with first") {
		t.Fatalf("the refusal doesn't give the read to name (%s) and the open tool by name: %s", latest.Read, words)
	}
	hand["read"] = latest.Read
	if _, sc := c.call(g, browseHand, hand); code(sc) != "" || sc["state"] != "handed" {
		t.Fatalf("a hand-over naming the read it was told: %v", sc)
	}
}

// openNameOf is the name g's browse_open is called by.
func openNameOf(g *Grant) string {
	for _, gt := range g.Tools {
		if gt.Tool == browseOpen {
			return gt.Name
		}
	}
	return ""
}
