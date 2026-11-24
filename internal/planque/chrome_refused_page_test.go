package planque

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
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

// Invariant: a step is built only on a read the taf was shown with its refs; browse_read and browse_checkout give a read and no refs, so neither is a step's default or a read a step may name
func TestAStepIsNeverBuiltOnAReadShownWithoutItsRefs(t *testing.T) {
	for _, reader := range []string{browseRead, browseCheckout} {
		for _, named := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/named=%v", reader, named), func(t *testing.T) {
				c := newChromeRig(t)
				g := c.grant(t, browseOpen, browseRead, browseCheckout, browseAct)
				n := 1
				var ops []string
				c.answer = func(step map[string]any) (bool, map[string]any, string) {
					ops = append(ops, fmt.Sprint(step["op"]))
					return true, pageN(n), ""
				}
				c.online(t)
				c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
				// The page re-renders: c1 is another control now, and only a read without refs sees it.
				n = 2
				c.call(g, reader, map[string]any{"url": "https://shop.example/"})
				latest := latestOf(t, c, g)
				if latest.Shown == latest.Read {
					t.Fatalf("%s made its read the one shown with refs", reader)
				}
				args := map[string]any{"url": "https://shop.example/", "steps": []any{map[string]any{"click": "c1"}}}
				if named {
					args["read"] = latest.Read
				}
				ops = nil
				out, sc := c.call(g, browseAct, args)
				if code(sc) != browser.CodePageChanged || slices.Contains(ops, "act") {
					t.Fatalf("%s then an act: %s (sent %v)", reader, out, ops)
				}
				checkRefusedWithPage(t, out, sc, "Search 2")
				// Shown with its refs by that refusal, the read now holds.
				ops = nil
				if out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "steps": []any{map[string]any{"click": "c1"}}}); code(sc) != "" || !slices.Contains(ops, "act") {
					t.Fatalf("an act on the page the refusal showed: %s", out)
				}
			})
		}
	}
}

// An older read shown without its refs is never one a step names, even when the page shows that state again.
func TestAnOlderReadShownWithoutItsRefsIsNeverNamedBack(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseRead, browseAct)
	n := 1
	var ops []string
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		ops = append(ops, fmt.Sprint(step["op"]))
		return true, pageN(n), ""
	}
	c.online(t)
	c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	n = 2
	_, read := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"})
	unshown, _ := read["read"].(string)
	n = 3
	c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	// The page goes back to the state only the text read saw.
	n = 2
	ops = nil
	out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": unshown, "steps": []any{map[string]any{"click": "c1"}}})
	if unshown == "" || code(sc) != browser.CodePageChanged || slices.Contains(ops, "act") {
		t.Fatalf("an act naming an older read shown without refs: %s (sent %v)", out, ops)
	}
	checkRefusedWithPage(t, out, sc, "Search 2")
}

// A page with a sign-in field. n changes nothing in it but its read.
func signInPage(n int) map[string]any {
	return shopPage(fmt.Sprintf("rd_s%d", n), map[string]any{"url": "https://shop.example/", "title": "Sign in",
		"fields":   []any{map[string]any{"ref": "f1", "label": "Email", "kind": "email"}, map[string]any{"ref": "f2", "label": "Password", "kind": "password"}},
		"controls": []any{map[string]any{"ref": "c1", "label": fmt.Sprintf("Sign in %d", n)}}})
}

// Invariant: at a page that asks for a password, a card number or a code, the taf is told in planque's own words to hand that step to the person, never to type it; and a step that would type one is refused with the same way out
func TestASignInPageIsHandedToThePersonNeverTyped(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeHand = true
	opens := c.grantAs(t, "agt_shop", browseOpen, browseAct)
	hands := c.grantAs(t, "agt_shop", browseHand, browseHandResult)
	handName := grantToolName(hands, browseHand)
	signIn := true
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if signIn {
			return true, signInPage(1), ""
		}
		return true, shopPage("rd_plain", map[string]any{"fields": []any{map[string]any{"ref": "f1", "label": "Search", "kind": "search"}}}), ""
	}
	c.online(t)
	out, sc := c.call(opens, browseOpen, map[string]any{"url": "https://shop.example/"})
	read, _ := sc["read"].(string)
	lines := strings.SplitN(out, "\n", 3)
	if len(lines) < 3 || lines[0] != readLine(read) || !strings.HasPrefix(lines[1], "Sign-in: ") || !strings.Contains(lines[1], "If your task needs that step, hand that step to the person with "+handName) {
		t.Fatalf("a sign-in page's answer doesn't say, after its read and outside the fence, to hand the sign-in over: %q", out)
	}
	// Typing the password is refused, and the refusal says the same way out.
	out, sc = c.call(opens, browseAct, map[string]any{"url": "https://shop.example/", "read": read, "steps": []any{map[string]any{"fill": "f2", "value": "hunter2"}}})
	if code(sc) != browser.CodeSecretField || !strings.Contains(out, "Hand that step to the person with "+handName) {
		t.Fatalf("typing a password: %s", out)
	}
	// A page with no such field says nothing of it.
	signIn = false
	if out, _ := c.call(opens, browseOpen, map[string]any{"url": "https://shop.example/"}); strings.Contains(out, "Sign-in: ") {
		t.Fatalf("a page with no sign-in field: %s", out)
	}
}

// With no hand-over tool, the way out is to ask the person in words.
func TestASignInPageWithNoHandOverToolAsksThePersonInWords(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, signInPage(1), "" }
	c.online(t)
	out, _ := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	if !strings.Contains(out, "Sign-in: ") || !strings.Contains(out, "If your task needs that step, ask the person to do that step themselves in this tab") || strings.Contains(out, "browse_hand") {
		t.Fatalf("no hand-over tool: %s", out)
	}
}

// Invariant: the owner's Chrome's other steps go to the extension as its acts, one each, and the hosted browser refuses them before anything runs
func TestTheChromesOtherStepsGoAsItsActsAndTheHostedBrowserRefusesThem(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t, browseOpen, browseAct)
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		return true, shopPage("rd_1", map[string]any{"fields": []any{map[string]any{"ref": "f1", "label": "City", "kind": "search"}},
			"controls": []any{map[string]any{"ref": "c1", "label": "Search"}}}), ""
	}
	c.online(t)
	_, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	read, _ := sc["read"].(string)
	for _, s := range []struct {
		step map[string]any
		act  string
	}{
		{map[string]any{"fill": "f1", "value": "Par", "keys": true}, `{"kind":"fill","ref":"f1","value":"Par","keys":true}`},
		{map[string]any{"check": "f1"}, `{"kind":"check","ref":"f1"}`},
		{map[string]any{"press": "f1", "key": "ArrowDown"}, `{"kind":"press","ref":"f1","value":"ArrowDown"}`},
		{map[string]any{"scroll": "c1"}, `{"kind":"scroll","ref":"c1"}`},
		{map[string]any{"scroll": "down"}, `{"kind":"scroll","value":"down"}`},
		{map[string]any{"back": true}, `{"kind":"back"}`},
	} {
		before := len(c.sent())
		if out, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": read, "steps": []any{s.step}}); code(sc) != "" {
			t.Fatalf("%v: %s", s.step, out)
		}
		var acts []string
		for _, st := range c.sent()[before:] {
			if st["op"] == "act" {
				raw, _ := json.Marshal(st["act"])
				acts = append(acts, string(raw))
			}
		}
		var want map[string]any
		_ = json.Unmarshal([]byte(s.act), &want)
		norm, _ := json.Marshal(want)
		if len(acts) != 1 || acts[0] != string(norm) {
			t.Fatalf("%v went as %v, want %s", s.step, acts, s.act)
		}
		_, sc = c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
		read, _ = sc["read"].(string)
	}
	// A key on a password field is refused by planque, before the Chrome is asked.
	before := len(c.sent())
	if _, sc := c.call(g, browseAct, map[string]any{"url": "https://shop.example/", "read": read, "steps": []any{map[string]any{"press": "f2", "key": "Enter"}}}); code(sc) != browser.CodeRefUnknown && code(sc) != browser.CodeSecretField {
		t.Fatalf("a key on a field that isn't a plain one: %v", sc)
	}
	for _, st := range c.sent()[before:] {
		if st["op"] == "act" {
			t.Fatalf("an act was sent: %v", st)
		}
	}

	e := newBrowserEnv(t)
	hosted := e.connect()
	for _, st := range []map[string]any{{"check": "qty"}, {"press": "qty", "key": "Enter"}, {"scroll": "down"}, {"back": true}, {"fill": "qty", "value": "2", "keys": true}} {
		sent := e.b.Sent()
		if out, sc := e.call(hosted, "browse_act", map[string]any{"url": fakes.ShopItem, "read": e.b.ReadOf(fakes.ShopItem), "steps": []map[string]any{st}}); code(sc) != CodeBadArguments || e.b.Sent() != sent {
			t.Fatalf("the hosted browser ran %v: %s", st, out)
		}
	}
}
