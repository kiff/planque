package planque

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// A page whose Search control is labelled by n: its read changes with n.
func pageN(n int) map[string]any {
	return shopPage(fmt.Sprintf("rd_%d", n), map[string]any{"controls": []any{map[string]any{"ref": "c1", "label": fmt.Sprintf("Search %d", n)}}})
}

func latestOf(t *testing.T, c *chromeRig, g *Grant) chromeRead {
	t.Helper()
	var r chromeRead
	if _, err := getDoc(context.Background(), c.s.KV, tenantPK(tenant), skChromeRead(g), &r); err != nil {
		t.Fatal(err)
	}
	return r
}

// The refusal's page: its read on its own line after planque's words, the
// page's words only inside one fence, and its read alone in the structured
// content.
func checkRefusedWithPage(t *testing.T, out string, sc map[string]any, words string) string {
	t.Helper()
	read, _ := sc["read"].(string)
	if read == "" || !strings.Contains(out, "\n"+readLine(read)+"\n") {
		t.Fatalf("the refusal doesn't give the page's read on its own line: %v / %s", sc, out)
	}
	open, close := strings.Index(out, "<<<page-words:"), strings.LastIndex(out, "<<</page-words:")
	if open < 0 || close < open || strings.Count(out, "<<<page-words:") != 1 {
		t.Fatalf("the page isn't in one fence: %s", out)
	}
	if i := strings.Index(out, words); i < open || i > close {
		t.Fatalf("the page's words %q are outside the fence: %s", words, out)
	}
	if strings.Contains(out[close:], words) || strings.Contains(fmt.Sprint(sc), words) || sc["page"] != nil {
		t.Fatalf("the page's words reach past the fence: %v", sc)
	}
	if fencedJSON(out[strings.Index(out, "<<<page-words:")-len(fenceNote)-1:])["read"] != read {
		t.Fatalf("the fenced page isn't the read named: %s", out)
	}
	return read
}

// Invariant: a step left without its read runs only on the taf's latest read of that very page
func TestAnActWithItsReadLeftOutRunsOnTheTafsLatestReadOfThatPageOnly(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	n := 1
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, pageN(n), "" }
	c.online(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); sc["read"] == nil {
		t.Fatalf("open: %v", sc)
	}
	// Left out, the read is the latest, and the step runs.
	if out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "steps": []any{map[string]any{"click": "c1"}}}); code(sc) != "" || sc["steps_run"] != float64(1) {
		t.Fatalf("an act with its read left out, on the page read last: %s", out)
	}
	// The page changed since: the latest read no longer holds, so nothing runs,
	// and the refusal carries the page as it reads now.
	n = 2
	out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "steps": []any{map[string]any{"click": "c1"}}})
	if code(sc) != browser.CodePageChanged || sc["sent"] != false {
		t.Fatalf("an act left without its read on a page that changed: %s", out)
	}
	read := checkRefusedWithPage(t, out, sc, "Search 2")
	if latestOf(t, c, g).Read != read {
		t.Fatal("the page the refusal gave isn't the latest read")
	}
	// Built from it, the next one runs.
	if out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "steps": []any{map[string]any{"click": "c1"}}}); code(sc) != "" {
		t.Fatalf("an act on the page the refusal gave: %s", out)
	}
	// Another address than the latest read's has no read to default to.
	if out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/items/lamp", "steps": []any{map[string]any{"click": "c1"}}}); code(sc) != browser.CodePageChanged {
		t.Fatalf("an act left without its read on another page: %s", out)
	}
}

// A read named and wrong is still wrong: no default replaces it.
func TestANamedReadIsNeverReplacedByTheLatest(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, pageN(1), "" }
	c.online(t)
	c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": "0123456789abcdef01234567", "steps": []any{map[string]any{"click": "c1"}}})
	if code(sc) != browser.CodePageChanged {
		t.Fatalf("an act naming a read that isn't one: %s", out)
	}
	checkRefusedWithPage(t, out, sc, "Search 1")
}

// Invariant: a step the owner's Chrome refused before doing anything comes back with the page as it reads now, read again, never acted on
func TestAStepTheChromeRefusedComesBackWithThePageReadAgain(t *testing.T) {
	for _, refused := range []string{"ref_unknown", "not_visible", "option_unknown", "disabled", browser.CodePageChanged} {
		t.Run(refused, func(t *testing.T) {
			c := newChromeRig(t)
			g := c.grant(t, browseOpen, browseAct)
			var ops []string
			c.answer = func(step map[string]any) (bool, map[string]any, string) {
				ops = append(ops, fmt.Sprint(step["op"]))
				if step["op"] == "act" {
					return false, nil, refused
				}
				return true, pageN(1), ""
			}
			c.online(t)
			_, opened := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
			ops = nil
			state, res := c.callState(g, browseAct, map[string]any{"url": "https://shop.example/", "read": opened["read"], "steps": []any{map[string]any{"click": "c1"}}})
			out, sc := textOf(res)
			if state != OpNotSent || code(sc) != refused || sc["sent"] != false || !strings.Contains(out, "Nothing was sent.") {
				t.Fatalf("%s: %s %s", refused, state, out)
			}
			if strings.Contains(out, "try another site") || !strings.Contains(out, chromeWhy[refused]) {
				t.Fatalf("the refusal isn't in planque's words: %s", out)
			}
			checkRefusedWithPage(t, out, sc, "Search 1")
			if strings.Join(ops, ",") != "open,act,read" {
				t.Fatalf("the steps sent: %v", ops)
			}
		})
	}
}

// A refusal of the owner's, or one a re-read can't get past, carries no page.
func TestARefusalNoReadGetsPastCarriesNoPage(t *testing.T) {
	for _, refused := range []string{"secret_field", "payment_step", "paused_by_owner", "stopped", "dialog_open"} {
		c := newChromeRig(t)
		g := c.grant(t, browseOpen, browseAct)
		c.answer = func(step map[string]any) (bool, map[string]any, string) {
			if step["op"] == "act" {
				return false, nil, refused
			}
			return true, pageN(1), ""
		}
		c.online(t)
		_, opened := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
		out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": opened["read"], "steps": []any{map[string]any{"click": "c1"}}})
		if code(sc) != refused || strings.Contains(out, "<<<page-words:") || sc["read"] != nil {
			t.Fatalf("%s: %s", refused, out)
		}
	}
}

// Invariant: planque has its own words for every code the owner's Chrome sends, and no word names a tool, which each taf names its own way
func TestEveryCodeTheChromeSendsHasPlanquesWordsNamingNoTool(t *testing.T) {
	// The codes chrome/src (guard.js, page.js, runner.js, relay.js) answers a
	// step with, as the extension sends them.
	sent := []string{
		"bad_step", "bad_url", "address_private", "site_not_listed", "site_blocked", "needs_site_permission",
		"page_changed", "ref_unknown", "not_visible", "secret_field", "payment_step", "dialog_open",
		"paused_by_owner", "paused_failures", "paused_no_progress", "stopped", "page_unresponsive",
		"browser_unavailable", "page_too_large", "option_unknown", "disabled", "no_page", "no_hand",
		"sealed_unknown", "release_missing",
	}
	tool := regexp.MustCompile(`\bbrowse_[a-z_]+`)
	for _, c := range sent {
		why, ok := chromeWhy[c]
		if !ok {
			t.Errorf("no words for %s", c)
		}
		if tool.MatchString(why) {
			t.Errorf("%s's words name a tool: %q", c, why)
		}
	}
	for c := range recoverable {
		if !chromeNotRun[c] && c != browser.CodePageChanged {
			t.Errorf("%s is read again after it, but isn't a step the Chrome refused before doing anything", c)
		}
	}
}

// A refusal that names how to read again names the open tool the taf has: its
// own connection's, or the one its reads are shared with.
func TestARefusalNamesTheOpenToolOfTheTafsOwnToolSet(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeHand = true
	opens := c.grantAs(t, "agt_wiki", browseOpen, browseAct)
	hands := c.grantAs(t, "agt_wiki", browseHand, browseHandResult)
	names := map[string]bool{}
	for _, g := range []*Grant{opens, hands} {
		for _, gt := range g.Tools {
			names[gt.Name] = true
		}
	}
	n := 1
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, pageN(n), "" }
	c.online(t)
	_, opened := c.call(opens, browseOpen, map[string]any{"url": "https://shop.example/"})
	n = 2
	act, _ := c.call(opens, browseAct, map[string]any{"url": "https://shop.example/", "read": opened["read"], "steps": []any{map[string]any{"click": "c1"}}})
	n = 3
	hand, _ := c.call(hands, browseHand, map[string]any{"url": "https://shop.example/", "read": opened["read"], "words": "Type your search", "ref": "c1"})
	for _, out := range []string{act, hand} {
		words := out[:strings.Index(out, "<<<page-words:")]
		named := regexp.MustCompile(`\bbrowse_[a-z_]+\b`).FindAllString(words, -1)
		if len(named) == 0 {
			t.Fatalf("the refusal names no open tool: %s", words)
		}
		for _, name := range named {
			if !names[name] {
				t.Fatalf("the refusal names %s, a tool this taf doesn't have (%v): %s", name, names, words)
			}
		}
	}
}

// Invariant: a hand-over tells the taf to end its turn and wait for the person, naming its own result tool, and a result read while the person hasn't answered says the same instead of inviting another read
func TestAHandOverAndItsWaitingAnswerTellTheTafToEndItsTurn(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeHand = true
	opens := c.grantAs(t, "agt_wiki", browseOpen, browseAct)
	hands := c.grantAs(t, "agt_wiki", browseHand, browseHandResult)
	resultName := grantToolName(hands, browseHandResult)
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, pageN(1), "" }
	c.online(t)
	c.call(opens, browseOpen, map[string]any{"url": "https://shop.example/"})
	out, sc := c.call(hands, browseHand, map[string]any{"url": "https://shop.example/", "words": "Type your search, then tap Done", "ref": "f2"})
	next, _ := sc["next"].(string)
	if sc["state"] != "handed" || !strings.Contains(next, "End your turn now") || !strings.Contains(next, resultName+" and this hand") {
		t.Fatalf("the hand-over's answer: %s", out)
	}
	out, sc = c.call(hands, browseHandResult, map[string]any{"hand": sc["hand"]})
	if sc["state"] != handWaiting || !strings.Contains(fmt.Sprint(sc["next"]), "don't ask again until they write that they're done") {
		t.Fatalf("a result read before the person answered: %s", out)
	}
}
