package planque

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// runs is the run of each step sent, in order, and whether a close was sent.
func (c *chromeRig) runs() (runs []string, closed []string) {
	for _, s := range c.sent() {
		run, _ := s["run"].(string)
		if s["op"] == "close" {
			closed = append(closed, run)
			continue
		}
		if len(runs) == 0 || runs[len(runs)-1] != run {
			runs = append(runs, run)
		}
	}
	return runs, closed
}

// tabRig is a Chrome whose run has one tab: a read reads it as it is, an open
// or an act loads a page there. Each load is a new read id, as the
// extension's are.
func tabRig(t *testing.T) (*chromeRig, *Grant) {
	t.Helper()
	c := newChromeRig(t)
	g := c.grant(t)
	var mu sync.Mutex
	shown := map[string]map[string]any{}
	loads := 0
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		mu.Lock()
		defer mu.Unlock()
		run, _ := step["run"].(string)
		if step["op"] == "read" {
			if p := shown[run]; p != nil {
				return true, p, ""
			}
			return false, nil, "no_page"
		}
		loads++
		u, _ := step["url"].(string)
		if u == "" {
			u = "https://shop.example/"
		}
		p := shopPage("rd_"+strings.Repeat("x", loads), map[string]any{"url": u, "controls": []any{map[string]any{"ref": "c1", "label": "Search"}},
			"links": []any{map[string]any{"ref": "l1", "label": "Lamp", "url": "https://shop.example/items/lamp"}, map[string]any{"ref": "l2", "label": "Desk", "url": "https://shop.example/items/desk"}}})
		shown[run] = p
		return true, p, ""
	}
	c.online(t)
	return c, g
}

// Invariant: a taf's calls that follow one another in the owner's Chrome are one run: one tab, one group, its bar and Stop for the whole turn
// Category: correctness
// Attack: none (Chrome 0.2.4's try: each call opened a new tab and group, so the bar, Stop and "Working in this Chrome" had nothing live to hold)
// Refused: nothing; three opens in one turn send one run id, no close between them, and a call on the page the tab shows reads it where it is
func TestATafsCallsInOneTurnAreOneRunInOneTab(t *testing.T) {
	c, g := tabRig(t)
	for _, u := range []string{"https://shop.example/", "https://shop.example/items/lamp", "https://shop.example/items/desk"} {
		if _, sc := c.call(g, browseOpen, map[string]any{"url": u}); code(sc) != "" {
			// The lamp's and the desk's links come from the start page's read.
			t.Fatalf("open %s: %v", u, sc)
		}
		// Each page links back to both, as a site's pages do.
	}
	runs, closed := c.runs()
	if len(runs) != 1 || len(closed) != 0 {
		t.Fatalf("three opens: runs %v, closed %v", runs, closed)
	}
	if ops := strings.Join(c.ops(), ","); ops != "open,open,open" {
		t.Fatalf("the steps: %s", ops)
	}
	// Step ids go on across the calls, never used twice.
	seen := map[string]bool{}
	for _, s := range c.sent() {
		id := s["id"].(string)
		if seen[id] {
			t.Fatalf("a step id used twice: %s", id)
		}
		seen[id] = true
	}
	// A read of the page the tab shows doesn't load it again.
	if _, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/items/desk"}); code(sc) != "" {
		t.Fatalf("read: %v", sc)
	}
	if ops := strings.Join(c.ops(), ","); ops != "open,open,open,read" {
		t.Fatalf("a read of the tab's page loaded it again: %s", ops)
	}
}

// Invariant: two calls at once never share a tab: the second runs alone, in its own tab, closed when it ends
func TestACallWhileTheRunIsBusyRunsAloneAndCloses(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	// Another call holds the grant's run now.
	it, ok, _ := c.s.KV.Get(context.Background(), tenantPK(g.Tenant), skChromeRun(g))
	if !ok {
		t.Fatal("no run kept for the grant")
	}
	it.Lease = time.Now().Add(time.Minute).UnixMilli()
	it.Ver++
	if err := c.s.KV.Put(context.Background(), it, storeAny()); err != nil {
		t.Fatal(err)
	}
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	waitFor(t, func() bool { _, closed := c.runs(); return len(closed) == 1 })
	runs, closed := c.runs()
	if len(runs) != 2 || runs[0] != first[0] || runs[1] == first[0] || closed[0] != runs[1] {
		t.Fatalf("the busy run's call: runs %v, closed %v", runs, closed)
	}
}

// Invariant: a run the person stopped is never used again: the taf is told, and its next call is a new run
func TestARunThePersonStoppedEndsAndTheNextCallIsANewRun(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	answer := c.answer
	c.mu.Lock()
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if step["run"] == first[0] {
			return false, nil, "stopped"
		}
		return answer(step)
	}
	c.mu.Unlock()
	out, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/items/lamp"})
	if code(sc) != "stopped" || !strings.Contains(out, chromeWhy["stopped"]) {
		t.Fatalf("a stopped run: %s", out)
	}
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("the next call: %v", sc)
	}
	runs, _ := c.runs()
	if len(runs) != 2 || runs[1] == first[0] {
		t.Fatalf("after the stop: runs %v", runs)
	}
}

// Invariant: an idle run ends: the grant's next call is a new run, and the old one is closed, so its group greys
func TestAnIdleRunIsClosedByTheNextCall(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	// Its last call ended longer ago than a run waits.
	var doc chromeRunDoc
	it, err := getDoc(context.Background(), c.s.KV, tenantPK(g.Tenant), skChromeRun(g), &doc)
	if err != nil {
		t.Fatal(err)
	}
	doc.At = time.Now().Add(-chromeRunIdle - time.Second)
	it.Ver++
	if err := putDoc(context.Background(), c.s.KV, it, doc, storeAny()); err != nil {
		t.Fatal(err)
	}
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	waitFor(t, func() bool { _, closed := c.runs(); return len(closed) == 1 })
	runs, closed := c.runs()
	if len(runs) != 2 || runs[1] == first[0] || closed[0] != first[0] {
		t.Fatalf("after idling: runs %v, closed %v", runs, closed)
	}
}

// Invariant: a run whose tab is gone in the Chrome goes on as a new run, and the call's page opens there
func TestARunWithNoTabGoesOnAsANewRun(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	// The extension's service worker started again: it holds no tab for the run.
	answer := c.answer
	c.mu.Lock()
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if step["run"] == first[0] && step["op"] == "read" {
			return false, nil, "no_page"
		}
		return answer(step)
	}
	c.mu.Unlock()
	if _, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	if ops := strings.Join(c.ops(), ","); ops != "open,read,open" {
		t.Fatalf("the steps: %s", ops)
	}
	runs, _ := c.runs()
	if len(runs) != 2 || runs[1] == first[0] {
		t.Fatalf("runs %v", runs)
	}
	// The new run is the grant's from now on.
	if _, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	if again, _ := c.runs(); len(again) != 2 {
		t.Fatalf("runs %v", again)
	}
}

// Invariant: a run's sites are fixed at its first step: a grant whose signed sites changed starts a new run
func TestNewSignedSitesStartANewRun(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	var doc chromeRunDoc
	it, err := getDoc(context.Background(), c.s.KV, tenantPK(g.Tenant), skChromeRun(g), &doc)
	if err != nil {
		t.Fatal(err)
	}
	doc.Sites = "an earlier signature"
	it.Ver++
	if err := putDoc(context.Background(), c.s.KV, it, doc, storeAny()); err != nil {
		t.Fatal(err)
	}
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	waitFor(t, func() bool { _, closed := c.runs(); return len(closed) == 1 })
	runs, closed := c.runs()
	if len(runs) != 2 || runs[1] == first[0] || closed[0] != first[0] {
		t.Fatalf("new sites: runs %v, closed %v", runs, closed)
	}
}

// Invariant: a page is read in place only while the run's tab still shows the address asked for: a tab the person took elsewhere opens that address again
func TestATabThePersonTookElsewhereOpensTheAskedPageAgain(t *testing.T) {
	c, g := tabRig(t)
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); code(sc) != "" {
		t.Fatalf("%v", sc)
	}
	first, _ := c.runs()
	// The person followed a link in the run's tab: a read now reads that page.
	answer := c.answer
	c.mu.Lock()
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		if step["op"] == "read" {
			return true, shopPage("rd_elsewhere", map[string]any{"url": "https://shop.example/items/lamp"}), ""
		}
		return answer(step)
	}
	c.mu.Unlock()
	_, sc := c.call(g, browseRead, map[string]any{"url": "https://shop.example/"})
	if code(sc) != "" || sc["url"] != "https://shop.example/" {
		t.Fatalf("the asked page: %v", sc)
	}
	if ops := strings.Join(c.ops(), ","); ops != "open,read,open" {
		t.Fatalf("the steps: %s", ops)
	}
	// Still the one run, in the same tab.
	if runs, closed := c.runs(); len(runs) != 1 || runs[0] != first[0] || len(closed) != 0 {
		t.Fatalf("runs %v, closed %v", runs, closed)
	}
}
