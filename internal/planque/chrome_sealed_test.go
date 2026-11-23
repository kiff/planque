package planque

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// callSealed is call for a sealed fill: op chosen by the test, and the
// gateway's headers as the relay reads them.
func (c *chromeRig) callSealed(g *Grant, tool, op string, args map[string]any, h http.Header) (string, map[string]any) {
	c.t.Helper()
	for _, gt := range g.Tools {
		if gt.Tool == tool {
			tool = gt.Name
		}
	}
	args[IdempotencyArg] = op
	raw := mustJSON(args)
	ctx := withSealedHeaders(context.Background(), h)
	return textOf(c.s.call(ctx, g, tool, raw, c.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: tool, Sub: g.Request.Agent, Op: op, Args: raw})))
}

func describeHeader(op, read, fill, sealed string) http.Header {
	raw, _ := json.Marshal(sealedDescribe{Op: op, Read: read, Fill: fill, Sealed: sealed})
	h := http.Header{}
	h.Set(SealedDescribeHeader, base64.RawURLEncoding.EncodeToString(raw))
	return h
}

func releaseHeader(rel map[string]any) http.Header {
	raw, _ := json.Marshal(rel)
	h := http.Header{}
	h.Set(SealedReleaseHeader, base64.RawURLEncoding.EncodeToString(raw))
	return h
}

// sealedRig is a Chrome rig with sealed fills on, whose taf has read the
// shop's sign-in page: its read, for the steps to name.
func sealedRig(t *testing.T) (*chromeRig, *Grant, string) {
	t.Helper()
	c := newChromeRig(t)
	c.s.SealedFill = true
	g := c.grant(t)
	n := 0
	c.answer = func(step map[string]any) (bool, map[string]any, string) {
		n++
		// Every read gets a new id, as the extension's do.
		return true, shopPage("rd_"+string(rune('a'+n)), nil), ""
	}
	c.online(t)
	_, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"})
	read, _ := sc["read"].(string)
	if read == "" {
		t.Fatalf("the taf's read: %v", sc)
	}
	// Its run closes after the call: counted before a test counts its own.
	waitFor(t, func() bool { return strings.Join(c.ops(), ",") == "open,close" })
	return c, g, read
}

// describe describes op's sealed fill of f1 with the password, as the
// gateway does, and returns the block and its signature.
func (c *chromeRig) describe(t *testing.T, g *Grant, op, read string) (map[string]any, string) {
	t.Helper()
	_, sc := c.callSealed(g, browseOpen, op+":open", map[string]any{"url": "https://shop.example/"}, describeHeader(op, read, "f1", "sec_a1:password"))
	bb, _ := sc["browser_block"].(map[string]any)
	block, _ := bb["block"].(map[string]any)
	sig, _ := bb["sig"].(string)
	if block == nil || sig == "" {
		t.Fatalf("describe: %v", sc)
	}
	return block, sig
}

func (c *chromeRig) release(op string, edit func(map[string]any)) map[string]any {
	r := map[string]any{"chrome_id": c.r.chrome, "op": op + ".2", "origin": "https://shop.example",
		"nonce": "c2l4dGVlbi1ieXRlLW5vbmNl", "exp": time.Now().Add(time.Minute).UnixMilli(), "sig": "SIGNED-BY-THE-PHONE"}
	if edit != nil {
		edit(r)
	}
	return r
}

func sealedAct(read string) map[string]any {
	return map[string]any{"url": "https://shop.example/", "read": read, "steps": []any{map[string]any{"fill": "f1", "sealed": "sec_a1:password"}}}
}

// Off unless turned on, and only a browse_act in the owner's Chrome whose one
// step is the fill: anything else is refused before the Chrome is asked.
func TestASealedFillIsRefusedUnlessOnAndItsCallsOnlyStep(t *testing.T) {
	c, g, read := sealedRig(t)
	before := len(c.sent())
	c.s.SealedFill = false
	if _, sc := c.callSealed(g, browseAct, "op_off", sealedAct(read), nil); code(sc) != CodeSealedOff {
		t.Fatalf("off: %v", sc)
	}
	if _, sc := c.callSealed(g, browseOpen, "op_off:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_off", read, "f1", "sec_a1:password")); code(sc) != CodeSealedOff {
		t.Fatalf("a describe while off: %v", sc)
	}
	c.s.SealedFill = true
	two := sealedAct(read)
	two["steps"] = []any{map[string]any{"fill": "f2", "value": "lamp"}, map[string]any{"fill": "f1", "sealed": "sec_a1:password"}}
	if _, sc := c.callSealed(g, browseAct, "op_two", two, nil); code(sc) != CodeBadArguments {
		t.Fatalf("two steps: %v", sc)
	}
	for _, step := range []map[string]any{
		{"fill": "f1", "sealed": "sec_a1:password", "value": "hunter2"},
		{"fill": "f1", "sealed": "sec_a1:password:username"},
		{"fill": "f1", "sealed": "sec_a1:pin"},
		{"click": "c1", "sealed": "sec_a1:password"},
	} {
		args := sealedAct(read)
		args["steps"] = []any{step}
		if _, sc := c.callSealed(g, browseAct, "op_bad_"+randomSuffix(), args, nil); code(sc) != CodeBadArguments {
			t.Errorf("%v: %v", step, sc)
		}
	}
	if len(c.sent()) != before {
		t.Fatalf("a refused sealed call reached the Chrome: %v", c.ops()[before:])
	}
}

// Invariant: the owner's card shows planque's signed block for the exact field the owner's phone then releases, and the run that read stays open for that one fill
// Category: authority
// Attack: the gateway (or anyone with the connection) describes one page and lets another be filled, or the block's fields don't name the read the extension will check
// Refused: the block is signed with planque's step key and names step 2 of the call's own run, the read the describe made, the field, the handle, the origin and the Chrome; the run isn't closed; the open is the call's step 1
func TestADescribeOpensTheCallsRunKeepsItAndSignsTheBlock(t *testing.T) {
	c, g, read := sealedRig(t)
	before := len(c.sent())
	block, sig := c.describe(t, g, "op_fill", read)
	if !stepkey.Verify(c.r.planque.PublicKey(), labelDescribe, block, sig) {
		t.Fatal("the block isn't signed with planque's step key")
	}
	steps := c.sent()[before:]
	if len(steps) != 1 || steps[0]["op"] != "open" || steps[0]["run"] != "op_fill" || steps[0]["id"] != "op_fill.1" {
		t.Fatalf("the describe's steps: %v", steps)
	}
	rf, _ := block["release_fields"].(map[string]any)
	want := map[string]any{"chrome_id": c.r.chrome, "op": "op_fill.2", "sealed": "sec_a1:password", "origin": "https://shop.example", "field_ref": "f1"}
	for k, v := range want {
		if rf[k] != v {
			t.Errorf("release_fields.%s = %v, want %v", k, rf[k], v)
		}
	}
	if rf["read_id"] == "" || rf["read_id"] == nil {
		t.Errorf("no read_id: %v", rf)
	}
	if block["v"] != labelDescribe || block["tenant"] != tenant || block["op"] != "op_fill.2" || block["site"] != "shop.example" ||
		!strings.Contains(block["action_words"].(string), "“Password”") {
		t.Fatalf("block: %v", block)
	}
	// Nothing of the page's secrets is in it.
	if raw, _ := json.Marshal(block); strings.Contains(string(raw), "hunter2") {
		t.Fatalf("a field's value is in the block: %s", raw)
	}
	time.Sleep(50 * time.Millisecond)
	if ops := c.ops()[before:]; len(ops) != 1 {
		t.Fatalf("the kept run was closed: %v", ops)
	}
	// One describe per hold.
	if _, sc := c.callSealed(g, browseOpen, "op_fill:open2", map[string]any{"url": "https://shop.example/"}, describeHeader("op_fill", read, "f1", "sec_a1:password")); code(sc) == "" {
		t.Fatalf("a second describe of the same fill: %v", sc)
	}
}

// A describe's open is only its own: KIFF's permit for "<operation>:open",
// on the read the taf named, of a field the handle's part fits.
func TestADescribeIsRefusedForAnotherOpReadOrField(t *testing.T) {
	c, g, read := sealedRig(t)
	if _, sc := c.callSealed(g, browseOpen, "op_x:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_y", read, "f1", "sec_a1:password")); code(sc) != CodePermitInvalid {
		t.Errorf("another op: %v", sc)
	}
	if _, sc := c.callSealed(g, browseOpen, "op_r:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_r", "rd-other", "f1", "sec_a1:password")); code(sc) != "page_changed" {
		t.Errorf("another read: %v", sc)
	}
	if _, sc := c.callSealed(g, browseOpen, "op_k:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_k", read, "f2", "sec_a1:password")); code(sc) != "secret_field" {
		t.Errorf("a password into a search field: %v", sc)
	}
	if _, sc := c.callSealed(g, browseAct, "op_a:open", sealedAct(read), describeHeader("op_a", read, "f1", "sec_a1:password")); code(sc) != CodeBadArguments {
		t.Errorf("a describe on browse_act: %v", sc)
	}
	// The open follows the plain rules: only a start page or a read link.
	if _, sc := c.callSealed(g, browseOpen, "op_u:open", map[string]any{"url": "https://shop.example/account?next=x"}, describeHeader("op_u", read, "f1", "sec_a1:password")); code(sc) != CodeURLNotRead {
		t.Errorf("an address the taf wrote: %v", sc)
	}
	if n, _ := c.s.KV.Query(context.Background(), tenantPK(tenant), skSealedPrefix(c.r.chrome)); len(n) != 0 {
		t.Fatalf("a refused describe kept a run: %d", len(n))
	}
}

// Invariant: an allowed sealed fill sends step 2 of the kept run, on the read the owner saw, with the phone's release as it came, once; then the run closes
// Category: authority
// Attack: the taf re-sends its allowed sealed call to fill the sign-in a second time, or names another field, handle or read on the allowed call
// Refused: a second call finds no kept run and sends nothing; a call or release that doesn't fit the kept run is release_stale, sends nothing, and ends the run
func TestTheAllowedFillActsOnceOnTheKeptRead(t *testing.T) {
	c, g, read := sealedRig(t)
	block, _ := c.describe(t, g, "op_fill", read)
	rf := block["release_fields"].(map[string]any)
	before := len(c.sent())
	rel := c.release("op_fill", nil)
	out, sc := c.callSealed(g, browseAct, "op_fill", sealedAct(read), releaseHeader(rel))
	if code(sc) != "" || sc["steps_run"] != float64(1) {
		t.Fatalf("the fill: %s %v", out, sc)
	}
	waitFor(t, func() bool { return len(c.sent()) >= before+2 })
	act, closed := c.sent()[before], c.sent()[before+1]
	if act["id"] != "op_fill.2" || act["run"] != "op_fill" || act["op"] != "act" || act["read_id"] != rf["read_id"] {
		t.Fatalf("the act: %v", act)
	}
	a := act["act"].(map[string]any)
	if a["kind"] != "fill" || a["ref"] != "f1" || a["sealed"] != "sec_a1:password" || a["value"] != nil {
		t.Fatalf("the act's fill: %v", a)
	}
	if mustJSONString(act["release"]) != mustJSONString(rel) {
		t.Fatalf("the release wasn't relayed as the phone made it: %v", act["release"])
	}
	if closed["op"] != "close" || closed["id"] != "op_fill.3" {
		t.Fatalf("the run wasn't closed after: %v", closed)
	}
	after := len(c.sent())
	if _, sc := c.callSealed(g, browseAct, "op_fill_again", sealedAct(read), releaseHeader(c.release("op_fill_again", nil))); code(sc) != CodeReleaseMissing {
		t.Fatalf("a second fill: %v", sc)
	}
	if len(c.sent()) != after {
		t.Fatalf("a second fill reached the Chrome: %v", c.ops()[after:])
	}
}

// Every way the allowed call or its release doesn't fit the kept run fails
// closed: nothing is filled, and the run ends.
func TestASealedFillThatDoesntFitFailsClosed(t *testing.T) {
	cases := map[string]func(c *chromeRig, op, read string) (map[string]any, http.Header){
		"no release": func(c *chromeRig, op, read string) (map[string]any, http.Header) { return sealedAct(read), nil },
		"another Chrome": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			return sealedAct(read), releaseHeader(c.release(op, func(r map[string]any) { r["chrome_id"] = "chr_other" }))
		},
		"another step": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			return sealedAct(read), releaseHeader(c.release(op, func(r map[string]any) { r["op"] = op + ".1" }))
		},
		"another origin": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			return sealedAct(read), releaseHeader(c.release(op, func(r map[string]any) { r["origin"] = "https://evil.example" }))
		},
		"expired": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			return sealedAct(read), releaseHeader(c.release(op, func(r map[string]any) { r["exp"] = time.Now().Add(-time.Second).UnixMilli() }))
		},
		"another field": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			a := sealedAct(read)
			a["steps"] = []any{map[string]any{"fill": "f2", "sealed": "sec_a1:password"}}
			return a, releaseHeader(c.release(op, nil))
		},
		"another handle": func(c *chromeRig, op, read string) (map[string]any, http.Header) {
			a := sealedAct(read)
			a["steps"] = []any{map[string]any{"fill": "f1", "sealed": "sec_b2:password"}}
			return a, releaseHeader(c.release(op, nil))
		},
	}
	for name, mk := range cases {
		t.Run(name, func(t *testing.T) {
			c, g, read := sealedRig(t)
			c.describe(t, g, "op_fill", read)
			before := len(c.sent())
			args, h := mk(c, "op_fill", read)
			if _, sc := c.callSealed(g, browseAct, "op_fill", args, h); code(sc) != CodeReleaseStale && code(sc) != CodeReleaseMissing {
				t.Fatalf("%v", sc)
			}
			waitFor(t, func() bool { return len(c.sent()) > before })
			if ops := c.ops()[before:]; len(ops) != 1 || ops[0] != "close" {
				t.Fatalf("not filled and closed: %v", ops)
			}
			if items, _ := c.s.KV.Query(context.Background(), tenantPK(tenant), skSealedPrefix(c.r.chrome)); len(items) != 0 {
				t.Fatal("the kept run outlived a refused fill")
			}
		})
	}
}

// Another connection naming the same operation finds nothing: the kept run is
// the describing connection's.
func TestAKeptRunIsOnlyItsConnections(t *testing.T) {
	c, g, read := sealedRig(t)
	other := c.grantAs(t, "agt_other")
	c.describe(t, g, "op_fill", read)
	before := len(c.sent())
	if _, sc := c.callSealed(other, browseAct, "op_fill", sealedAct(read), releaseHeader(c.release("op_fill", nil))); code(sc) != CodeReleaseMissing {
		t.Fatalf("another connection: %v", sc)
	}
	if len(c.sent()) != before {
		t.Fatalf("another connection's call reached the Chrome: %v", c.ops()[before:])
	}
}

// At most two sealed fills wait on one Chrome; one past its time is closed
// by the next sealed call, and its allowed call finds nothing.
func TestKeptRunsAreCappedAndEndWithTheirTime(t *testing.T) {
	c, g, read := sealedRig(t)
	c.describe(t, g, "op_1", read)
	c.describe(t, g, "op_2", read)
	if _, sc := c.callSealed(g, browseOpen, "op_3:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_3", read, "f1", "sec_a1:password")); code(sc) != CodeSealedBusy {
		t.Fatalf("a third: %v", sc)
	}
	// Both past their time, as their holds are.
	ctx := context.Background()
	items, _ := c.s.KV.Query(ctx, tenantPK(tenant), skSealedPrefix(c.r.chrome))
	for _, it := range items {
		var r sealedRun
		_ = json.Unmarshal(it.Doc, &r)
		r.Opened = r.Opened.Add(-sealedKept - time.Second)
		it.Doc, _ = json.Marshal(r)
		if err := c.s.KV.Put(ctx, it, store.Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	before := len(c.sent())
	if _, sc := c.callSealed(g, browseAct, "op_1", sealedAct(read), releaseHeader(c.release("op_1", nil))); code(sc) != CodeReleaseMissing {
		t.Fatalf("a fill past its time: %v", sc)
	}
	waitFor(t, func() bool { return len(c.sent()) >= before+2 })
	for _, s := range c.sent()[before:] {
		if s["op"] != "close" {
			t.Fatalf("only closes: %v", c.ops()[before:])
		}
	}
	// The Chrome has room again.
	c.describe(t, g, "op_4", read)
}

// A describe is claimed before anything reaches the Chrome: one racing it for
// the same fill is refused there, and an allowed fill that finds only that
// unfinished claim fills nothing.
func TestARacingDescribeIsRefusedBeforeTheChromeIsAsked(t *testing.T) {
	c, g, read := sealedRig(t)
	ctx := context.Background()
	// Another describe of op_fill has claimed it, and hasn't opened the page yet.
	claim := sealedRun{Tenant: tenant, Grant: g.ID, Agent: g.Request.Agent, Op: "op_fill", Chrome: c.r.chrome, Opened: time.Now().UTC()}
	it := store.Item{PK: tenantPK(tenant), SK: skSealedRun(c.r.chrome, g.ID, "op_fill"), Ver: 1, TTL: time.Now().Add(time.Hour).Unix()}
	if err := putDoc(ctx, c.s.KV, it, claim, store.Cond{NotExists: true}); err != nil {
		t.Fatal(err)
	}
	before := len(c.sent())
	if _, sc := c.callSealed(g, browseOpen, "op_fill:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_fill", read, "f1", "sec_a1:password")); code(sc) != CodeSealedBusy {
		t.Fatalf("a racing describe: %v", sc)
	}
	if _, sc := c.callSealed(g, browseAct, "op_fill", sealedAct(read), releaseHeader(c.release("op_fill", nil))); code(sc) != CodeReleaseMissing {
		t.Fatalf("a fill on an unfinished claim: %v", sc)
	}
	time.Sleep(50 * time.Millisecond)
	for _, s := range c.sent()[before:] {
		if s["op"] != "close" {
			t.Fatalf("a step other than close reached the Chrome: %v", c.ops()[before:])
		}
	}
}

// A refused describe gives its claim back: the Chrome's room isn't held by a
// fill that was never described.
func TestARefusedDescribeKeepsNothing(t *testing.T) {
	c, g, read := sealedRig(t)
	if _, sc := c.callSealed(g, browseOpen, "op_k:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_k", read, "f2", "sec_a1:password")); code(sc) != "secret_field" {
		t.Fatalf("%v", sc)
	}
	if items, _ := c.s.KV.Query(context.Background(), tenantPK(tenant), skSealedPrefix(c.r.chrome)); len(items) != 0 {
		t.Fatalf("a refused describe kept its claim: %d", len(items))
	}
	// Its room is free for the next (an operation is answered once, so a new one).
	c.describe(t, g, "op_k2", read)
}

// claimRaceKV puts a competing claim for the same fill just before the first
// claim of a kept run: two describes that both found the fill unclaimed.
type claimRaceKV struct {
	store.KV
	raced bool
}

func (k *claimRaceKV) Put(ctx context.Context, it store.Item, c store.Cond) error {
	if !k.raced && strings.HasPrefix(it.SK, "SEALEDRUN#") && c.NotExists {
		k.raced = true
		if err := k.KV.Put(ctx, it, store.Cond{NotExists: true}); err != nil {
			return err
		}
	}
	return k.KV.Put(ctx, it, c)
}

// Two describes of one fill, both past the "described already" read: only the
// first claim holds, and the second never asks the Chrome.
func TestTwoDescribesRacingForOneFillClaimItOnce(t *testing.T) {
	c, g, read := sealedRig(t)
	c.s.KV = &claimRaceKV{KV: c.s.KV}
	before := len(c.sent())
	if _, sc := c.callSealed(g, browseOpen, "op_fill:open", map[string]any{"url": "https://shop.example/"}, describeHeader("op_fill", read, "f1", "sec_a1:password")); code(sc) != CodeSealedBusy {
		t.Fatalf("the describe that lost the race: %v", sc)
	}
	time.Sleep(50 * time.Millisecond)
	if len(c.sent()) != before {
		t.Fatalf("the losing describe reached the Chrome: %v", c.ops()[before:])
	}
}
