package planque

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/store"
)

// callIn is call, made in the kiff/agents run run, as KIFF signs it in the
// permit (claim "run").
func (c *chromeRig) callIn(g *Grant, run, tool string, args map[string]any) (string, map[string]any) {
	c.t.Helper()
	op := "op_" + strings.ToLower(tool) + "_" + randomSuffix()
	args[IdempotencyArg] = op
	raw := mustJSON(args)
	return textOf(c.s.call(context.Background(), g, tool, raw, c.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: tool, Sub: g.Request.Agent, Op: op, Args: raw, Run: run})))
}

// ownRuns is the panel's read of its Chrome's runs, with KIFF session key.
func ownRuns(t *testing.T, e *env, key string) (int, []any) {
	t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/chromes/self/runs", nil, nil, key)
	runs, _ := out["runs"].([]any)
	return st, runs
}

// The extension's KIFF offer's key is kept with the Chrome at the confirm,
// from the offer planque relayed; an offer that isn't the extension's keeps
// none.
func TestTheConfirmKeepsTheKeyTheExtensionOfferedToKIFF(t *testing.T) {
	if got := kiffOfferPub(`{"exp":"2026-11-21T12:04:30Z","key_store":"software","kind":"extension","name":"tafs for Chrome","nonce":"n","platform":"chrome","pub":"MFkw_ext","v":"kiff.pair.v1"}`); got != "MFkw_ext" {
		t.Fatalf("the extension's offer: %q", got)
	}
	for _, offer := range []string{"", "not json", `{"v":"kiff.pair.v1","platform":"ios","pub":"MFkw_phone"}`, `{"v":"kiff.pair.v2","platform":"chrome","pub":"MFkw"}`} {
		if got := kiffOfferPub(offer); got != "" {
			t.Errorf("%q keeps %q", offer, got)
		}
	}
	e := chromeEnv(t)
	acct := "/v1/accounts/" + tenant
	_, out := e.do("POST", acct+"/chromes/pairings", map[string]any{}, e.phone, "")
	pairing, lookup := out["pairing"].(string), out["lookup"].(string)
	extKey, extSum := p256Key(t)
	offer := `{"exp":"2026-11-21T12:04:30Z","key_store":"software","kind":"extension","name":"tafs for Chrome","nonce":"n","platform":"chrome","pub":"MFkw_ext","v":"kiff.pair.v1"}`
	if st, out := e.do("POST", "/v1/extension/pairings/"+lookup, map[string]any{"key": extKey, "version": "1.0.0", "manifest_hash": releasedManifest,
		"browser": "Chrome 131", "os": "macOS 15", "proof": b64(32), "kiff_offer": offer}, nil, ""); st != 202 {
		t.Fatalf("the extension's key: %d %v", st, out)
	}
	phoneKey, _ := p256Key(t)
	confirm := map[string]any{"chrome_key_sha256": extSum, "name": "Work Mac", "phone_key": phoneKey, "phone_proof": b64(32), "phone_sig": b64(64), "kiff_pub": "MFkw_claimed"}
	// A key the confirm names is nobody's word: refused, never kept.
	if st, _ := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", confirm, e.phone, ""); st != 400 {
		t.Fatalf("a confirm naming a key: %d, want 400", st)
	}
	delete(confirm, "kiff_pub")
	st, out := e.do("POST", acct+"/chromes/pairings/"+pairing+"/confirm", confirm, e.phone, "")
	if st != 200 {
		t.Fatalf("confirm: %d %v", st, out)
	}
	var c Chrome
	if _, err := getDoc(context.Background(), e.s.KV, tenantPK(tenant), skChrome(out["chrome_id"].(string)), &c); err != nil {
		t.Fatal(err)
	}
	if c.KiffPub != "MFkw_ext" || c.KiffDevice != "" {
		t.Fatalf("kept: pub %q device %q", c.KiffPub, c.KiffDevice)
	}
	// Never listed.
	if _, out := e.do("GET", acct+"/chromes", nil, e.phone, ""); strings.Contains(mustJSONString(out), "MFkw_ext") {
		t.Fatalf("the list shows the key: %v", out)
	}
}

// Each run that sends a step to a Chrome is listed for that Chrome's own
// KIFF session, newest first, by the run KIFF signed and the taf's agent;
// any other session, a Chrome that isn't this one, or a call with no run
// lists nothing, and unpairing drops them all.
func TestTheRunsThatWorkedInAChromeAreListedForItsOwnSessionOnly(t *testing.T) {
	c := newChromeRig(t)
	g := c.grant(t)
	c.answer = func(step map[string]any) (bool, map[string]any, string) { return true, plainPage("rd_1"), "" }
	c.online(t)
	ctx := context.Background()
	var ch Chrome
	it, err := getDoc(ctx, c.s.KV, tenantPK(tenant), skChrome(c.r.chrome), &ch)
	if err != nil {
		t.Fatal(err)
	}
	ch.KiffPub = "MFkw_ext"
	if err := putDoc(ctx, c.s.KV, store.Item{PK: it.PK, SK: it.SK}, ch, store.Cond{Ver: &it.Ver}); err != nil {
		t.Fatal(err)
	}
	c.kiff.ChromeDevices[tenant+" MFkw_ext"] = "dev_panel"
	c.kiff.ChromeDevices[tenant+" MFkw_other"] = "dev_other"
	c.kiff.Keys["s_panel"] = kiff.Me{TenantID: tenant, DeviceSession: true, DevicePlatform: "chrome", DeviceID: "dev_panel", Capabilities: []string{"read"}}
	c.kiff.Keys["s_other_chrome"] = kiff.Me{TenantID: tenant, DeviceSession: true, DevicePlatform: "chrome", DeviceID: "dev_other", Capabilities: []string{"read"}}
	c.kiff.Keys["s_phone"] = kiff.Me{TenantID: tenant, DeviceSession: true, DevicePlatform: "ios", DeviceID: "dev_panel", Capabilities: []string{"read"}}
	c.kiff.Keys["s_elsewhere"] = kiff.Me{TenantID: "t_other", DeviceSession: true, DevicePlatform: "chrome", DeviceID: "dev_panel", Capabilities: []string{"read"}}

	if st, runs := ownRuns(t, c.env, "s_panel"); st != 200 || len(runs) != 0 {
		t.Fatalf("before any run: %d %v", st, runs)
	}
	if _, sc := c.call(g, browseOpen, map[string]any{"url": "https://shop.example/"}); sc["site"] != "shop.example" {
		t.Fatalf("a call with no run: %v", sc)
	}
	if _, runs := ownRuns(t, c.env, "s_panel"); len(runs) != 0 {
		t.Fatalf("a call with no run is listed: %v", runs)
	}
	c.callIn(g, "run_first", browseOpen, map[string]any{"url": "https://shop.example/"})
	c.callIn(g, "run_first", browseOpen, map[string]any{"url": "https://shop.example/"})
	time.Sleep(5 * time.Millisecond)
	c.callIn(g, "run_second", browseOpen, map[string]any{"url": "https://shop.example/"})
	st, runs := ownRuns(t, c.env, "s_panel")
	if st != 200 || len(runs) != 2 {
		t.Fatalf("two runs: %d %v", st, runs)
	}
	newest, oldest := runs[0].(map[string]any), runs[1].(map[string]any)
	if newest["run"] != "run_second" || oldest["run"] != "run_first" || newest["taf"] != g.Request.Agent || len(newest) != 3 {
		t.Fatalf("newest first, by run and taf: %v", runs)
	}
	asked := c.kiff.ChromeDevicesAsked
	if asked != 1 {
		t.Fatalf("KIFF asked %d times, want once: the device is kept once named", asked)
	}
	for _, key := range []string{"s_other_chrome", "s_phone", "k_setup"} {
		if st, runs := ownRuns(t, c.env, key); st != 200 || len(runs) != 0 {
			t.Errorf("%s: %d %v, want none", key, st, runs)
		}
	}
	if st, _ := ownRuns(t, c.env, "s_elsewhere"); st != 403 {
		t.Errorf("another account's session: %d, want 403", st)
	}
	// A connection moved to another Chrome takes its runs out of this one's
	// panel; another connection's run in this Chrome stays.
	c.s.markChromeWorked(ctx, tenant, c.r.chrome, "run_other_connection", "agt_other", "grt_other")
	home := Chrome{ID: "chr_homelaptop01", Tenant: tenant, Name: "Home laptop", Key: "k"}
	if err := putDoc(ctx, c.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(home.ID)}, home, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	newSig := strings.Replace(testSitesSig, "s1gned", "s1GNED", 1)
	if st, out := c.do("PUT", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/chrome", map[string]any{"chrome_id": home.ID, "chrome_sites_sig": newSig}, c.phone, ""); st != 200 {
		t.Fatalf("switch: %d %v", st, out)
	}
	if _, runs := ownRuns(t, c.env, "s_panel"); len(runs) != 1 || runs[0].(map[string]any)["run"] != "run_other_connection" {
		t.Fatalf("after the move, the old Chrome lists: %v", runs)
	}
	if left, _ := c.s.KV.Query(ctx, tenantPK(tenant), skChromeWorked(home.ID, "")); len(left) != 0 {
		t.Fatalf("the move carried %d runs to the new Chrome", len(left))
	}
	// Unpairing drops the Chrome's marks with it.
	if _, err := c.s.removeChrome(ctx, tenant, c.r.chrome); err != nil {
		t.Fatal(err)
	}
	if left, _ := c.s.KV.Query(ctx, tenantPK(tenant), skChromeWorked(c.r.chrome, "")); len(left) != 0 {
		t.Fatalf("marks left after unpairing: %d", len(left))
	}
	if _, runs := ownRuns(t, c.env, "s_panel"); len(runs) != 0 {
		t.Fatalf("an unpaired Chrome lists: %v", runs)
	}
}

// KIFF makes the extension's device when the extension claims its pairing,
// after the phone's confirm: a Chrome whose device KIFF doesn't name yet
// lists nothing, is asked about again only after a minute, and keeps the
// device once named. A Chrome paired before its key was kept is never asked.
func TestAChromesDeviceIsAskedForUntilKIFFNamesIt(t *testing.T) {
	e := chromeEnv(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 25, 21, 0, 0, 0, time.UTC)
	e.s.Now = func() time.Time { return now }
	for _, c := range []Chrome{{ID: "chr_newpairing01", Tenant: tenant, Name: "Work Mac", Key: "k", KeySHA256: "s", KiffPub: "MFkw_new"},
		{ID: "chr_oldpairing01", Tenant: tenant, Name: "Old Mac", Key: "k2", KeySHA256: "s2"}} {
		if err := putDoc(ctx, e.s.KV, store.Item{PK: tenantPK(tenant), SK: skChrome(c.ID)}, c, store.Cond{NotExists: true}); err != nil {
			t.Fatal(err)
		}
	}
	e.s.markChromeWorked(ctx, tenant, "chr_newpairing01", "run_1", "agt_taf", "grt_1")
	e.s.markChromeWorked(ctx, tenant, "chr_oldpairing01", "run_2", "agt_taf", "grt_2")
	e.kiff.Keys["s_panel"] = kiff.Me{TenantID: tenant, DeviceSession: true, DevicePlatform: "chrome", DeviceID: "dev_panel", Capabilities: []string{"read"}}

	if _, runs := ownRuns(t, e, "s_panel"); len(runs) != 0 || e.kiff.ChromeDevicesAsked != 1 {
		t.Fatalf("not claimed yet: %v, asked %d (the old pairing is never asked)", runs, e.kiff.ChromeDevicesAsked)
	}
	e.kiff.ChromeDevices[tenant+" MFkw_new"] = "dev_panel"
	now = now.Add(30 * time.Second)
	if _, runs := ownRuns(t, e, "s_panel"); len(runs) != 0 || e.kiff.ChromeDevicesAsked != 1 {
		t.Fatalf("within the minute: %v, asked %d", runs, e.kiff.ChromeDevicesAsked)
	}
	now = now.Add(31 * time.Second)
	_, runs := ownRuns(t, e, "s_panel")
	if len(runs) != 1 || runs[0].(map[string]any)["run"] != "run_1" || e.kiff.ChromeDevicesAsked != 2 {
		t.Fatalf("named: %v, asked %d", runs, e.kiff.ChromeDevicesAsked)
	}
	var c Chrome
	if _, err := getDoc(ctx, e.s.KV, tenantPK(tenant), skChrome("chr_newpairing01"), &c); err != nil || c.KiffDevice != "dev_panel" || c.KiffAsked != nil {
		t.Fatalf("kept: %+v %v", c, err)
	}
	now = now.Add(time.Hour)
	ownRuns(t, e, "s_panel")
	if e.kiff.ChromeDevicesAsked != 2 {
		t.Fatalf("asked again once named: %d", e.kiff.ChromeDevicesAsked)
	}
}
