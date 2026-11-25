package planque

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// A tab of the person's the run joined (Chrome 0.2.15): what planque gives
// the taf when the person moves in it, and what it never gives.

// The person's page, as a Chrome would send it after they moved: another
// address on the run's sites, with words and a field value of theirs.
func personsPage() map[string]any {
	return shopPage("rd_person", map[string]any{"url": "https://shop.example/account", "title": "Your account",
		"text":   "Ana Ruiz, 12 rue des Lilas, card ending 4242.",
		"fields": []any{map[string]any{"ref": "f9", "label": "Name", "kind": "text", "value": "Ana Ruiz"}}})
}

// Invariant: a step's refused answer gives the taf nothing of the page it carried, and a re-read after it gives the page only while the tab is still on the call's address
func TestAPageThePersonMovedToNeverReachesTheTaf(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	var ops []string
	moved := false
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		ops = append(ops, fmt.Sprint(step["op"]))
		switch {
		case step["op"] == "act":
			// The person moved: the Chrome refuses the act, and a careless one sends their page along.
			moved = true
			return false, personsPage(), browser.CodePageChanged
		case moved:
			// A Chrome that let the re-read through: the person's page, on the run's sites.
			return true, personsPage(), ""
		}
		return true, pageN(1), ""
	}
	c.online(t)
	_, opened := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	before := latestOf(t, c, g)
	ops = nil
	out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": opened["read"], "steps": []any{map[string]any{"click": "c1"}}})
	if code(sc) != browser.CodePageChanged || sc["sent"] != false {
		t.Fatalf("the refused act: %s", out)
	}
	for _, theirs := range []string{"Ana Ruiz", "rue des Lilas", "4242", "Your account", "/account"} {
		if strings.Contains(out, theirs) || strings.Contains(fmt.Sprint(sc), theirs) {
			t.Fatalf("the person's page reached the taf (%q): %s %v", theirs, out, sc)
		}
	}
	if strings.Contains(out, "<<<page-words:") || sc["read"] != nil {
		t.Fatalf("a page came with the refusal: %s", out)
	}
	// Read where it is, the act, then read again (a read changes nothing) and dropped: never kept as the grant's latest read.
	if strings.Join(ops, ",") != "read,act,read" {
		t.Fatalf("the steps sent: %v", ops)
	}
	if latestOf(t, c, g).Read != before.Read {
		t.Fatal("the person's page was kept as the grant's latest read")
	}
}

// Invariant: a tab given back to the person ends the run's work there: planque's plain words for the person, nothing read again, no other way tried
func TestATabGivenBackEndsTheRunsWorkThereInPlainWords(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	var ops []string
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		ops = append(ops, fmt.Sprint(step["op"]))
		if step["op"] == "act" {
			return false, personsPage(), "given_back"
		}
		return true, pageN(1), ""
	}
	c.online(t)
	_, opened := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	ops = nil
	state, res := c.callState(g, browseAct, map[string]any{"url": "https://shop.example/", "read": opened["read"], "steps": []any{map[string]any{"click": "c1"}}})
	out, sc := textOf(res)
	if state != OpNotSent || code(sc) != "given_back" || sc["sent"] != false {
		t.Fatalf("given back: %s %s", state, out)
	}
	if !strings.Contains(out, chromeWhy["given_back"]) || !strings.Contains(out, notAnotherWay) {
		t.Fatalf("not planque's plain words: %s", out)
	}
	if strings.Contains(out, "Ana Ruiz") || strings.Contains(out, "<<<page-words:") || sc["read"] != nil {
		t.Fatalf("a page came with it: %s", out)
	}
	if strings.Join(ops, ",") != "read,act" {
		t.Fatalf("read again after the tab was given back: %v", ops)
	}
}
