package planque

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/profile"
)

// Gap 13 of kiff/tafs#476: a connection whose key or sign-in the vendor
// refuses is fixed from the app, in place. The grant, its tools and KIFF stay
// as they are, and nothing changes until the new key or sign-in is proven.

func lastEvent(g *Grant) string {
	if len(g.Events) == 0 {
		return ""
	}
	return g.Events[len(g.Events)-1].What
}

// A connected key the vendor now refuses is replaced after one read with the
// new key. A key the vendor refuses changes nothing, and the old one stays.
func TestAConnectedKeyIsReplacedAfterOneReadWithIt(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	g := e.connectAPI(v)
	ctx := context.Background()
	old := v.key
	kiffTools := len(e.kiff.Connections(tenant))

	// The vendor rotates its key: the one held is refused, and the grant
	// says connected, as before.
	v.mu.Lock()
	v.key = "sk_test_rotated_456"
	v.mu.Unlock()
	res, _ := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{})
	if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "(401)") {
		t.Fatalf("a call with the rotated-out key: %s %v", text, isErr)
	}

	// A wrong key: the vendor refuses the read, and nothing changes.
	events := len(g.Events)
	before := len(v.seen())
	st, out := e.giveKey(g.ID, g.Digest, "sk_test_wrong_789")
	if st != 400 || !strings.Contains(mustJSONString(out), "refused the new key") || !strings.Contains(mustJSONString(out), "the old key stays") {
		t.Fatalf("a wrong key: %d %v", st, out)
	}
	if calls := v.seen()[before:]; len(calls) != 1 || calls[0].Method != "GET" || calls[0].Path != "/v1/charges" || calls[0].Key != "sk_test_wrong_789" {
		t.Fatalf("the check: %+v", calls)
	}
	if tok, _, _ := e.s.openToken(ctx, g); tok.Access != old {
		t.Fatal("a refused key replaced the one held")
	}
	if now, _ := e.s.loadGrant(ctx, tenant, g.ID); len(now.Events) != events || now.Status != StatusConnected {
		t.Fatalf("after a refused key: %s, events %d→%d", now.Status, events, len(now.Events))
	}

	// The vendor's new key: one read with it works, and it replaces the old.
	st, out = e.giveKey(g.ID, g.Digest, v.key)
	if st != 200 || out["status"] != StatusConnected || out["id"] != g.ID {
		t.Fatalf("the new key: %d %v", st, out)
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if ev := lastEvent(now); !strings.HasPrefix(ev, "key replaced on ") || !strings.Contains(ev, "(GET /charges)") || !strings.Contains(ev, "its tools are as they were") {
		t.Fatalf("event: %q", ev)
	}
	if len(now.Tools) != len(g.Tools) || now.Digest != g.Digest || len(e.kiff.Connections(tenant)) != kiffTools {
		t.Fatalf("tools %d→%d, digest changed %v, KIFF %d→%d", len(g.Tools), len(now.Tools), now.Digest != g.Digest, kiffTools, len(e.kiff.Connections(tenant)))
	}
	res, _ = e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{})
	if text, isErr := resultText(t, res.Result); isErr || !strings.Contains(text, "ch_1") || v.seen()[len(v.seen())-1].Key != v.key {
		t.Fatalf("a call with the new key: %s %v", text, isErr)
	}
	// Sealed, as the first key was.
	items, _ := e.kv.Query(ctx, tenantPK(tenant), "")
	for _, it := range items {
		if bytes.Contains(it.Doc, []byte(v.key)) {
			t.Fatalf("the new key is in the clear in %s", it.SK)
		}
	}
	// Signed for this grant, or nothing: another digest is refused.
	if st, _ := e.giveKey(g.ID, "sha256:"+strings.Repeat("0", 64), "sk_test_other_000"); st != 409 {
		t.Fatalf("a key signed for another digest: %d", st)
	}
	// Approve is for sign-ins: a key grant's replacement goes to /key.
	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": g.Digest}, e.phone, ""); st != 409 || !strings.Contains(mustJSONString(out), "POST …/key") {
		t.Fatalf("approve on a connected key grant: %d %v", st, out)
	}
}

// An API with no read that needs no argument can't check a new key, so the
// key isn't replaced: the owner is told how to use another one.
func TestAKeyThatCantBeCheckedIsNotReplaced(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	v.desc = strings.Replace(restDesc, `"parameters": [{"name": "limit", "in": "query", "schema": {"type": "integer"}}]`,
		`"parameters": [{"name": "customer", "in": "query", "required": true, "schema": {"type": "string"}}]`, 1)
	g := e.connectAPI(v)
	before := len(v.seen())
	st, out := e.giveKey(g.ID, g.Digest, "sk_test_rotated_456")
	if st != 409 || !strings.Contains(mustJSONString(out), "no read that takes no arguments") || !strings.Contains(mustJSONString(out), "disconnect it and connect it again") {
		t.Fatalf("a key that can't be checked: %d %v", st, out)
	}
	if len(v.seen()) != before {
		t.Fatalf("calls at the vendor: %+v", v.seen()[before:])
	}
	if tok, _, _ := e.s.openToken(context.Background(), g); tok.Access != v.key {
		t.Fatal("an unchecked key replaced the one held")
	}
}

// Signing in again on a connected grant: the sign-in in use goes on working
// until the approving phone confirms the new one, which then replaces it,
// with the same tools, and the old one is revoked at the vendor.
func TestAConnectedSignInIsReplacedOnlyOnceConfirmed(t *testing.T) {
	e := newEnv(t)
	e.vendor.Account = "acct_1"
	g := e.connect()
	ctx := context.Background()
	kiffTools := len(e.kiff.Connections(tenant))
	old, _, _ := e.s.openToken(ctx, g)

	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": g.Digest}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected || out["sign_in_url"] == nil {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if now.Status != StatusConnected || now.Renewal == nil || now.Reconnect {
		t.Fatalf("while signing in again: %s, renewal %v, reconnect %v", now.Status, now.Renewal, now.Reconnect)
	}
	// The connection works while the new sign-in is under way.
	if res, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{}); err != nil || strings.Contains(string(res.Result), "Not sent") || strings.Contains(string(res.Result), "signs in to") {
		t.Fatalf("a call while signing in again: %s %v", res.Result, err)
	}
	code := e.signIn(out["sign_in_url"].(string))
	if tok, _, _ := e.s.openToken(ctx, g); tok.Access != old.Access {
		t.Fatal("an unconfirmed sign-in replaced the one in use")
	}
	if res, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{}); err != nil || strings.Contains(string(res.Result), "Not sent") {
		t.Fatalf("a call before the confirmation: %s %v", res.Result, err)
	}
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/confirm", map[string]string{"confirm": "not-the-code"}, e.phone, ""); st != 403 {
		t.Fatalf("another code: %d", st)
	}

	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/confirm", map[string]string{"confirm": code}, e.phone, "")
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm: %d %v", st, out)
	}
	now, _ = e.s.loadGrant(ctx, tenant, g.ID)
	tok, _, _ := e.s.openToken(ctx, g)
	if now.Renewal != nil || tok.Access == old.Access || now.Account != "acct_1" || len(now.Tools) != len(g.Tools) || len(e.kiff.Connections(tenant)) != kiffTools {
		t.Fatalf("after the confirmation: renewal %v, token replaced %v, account %q, tools %d→%d", now.Renewal, tok.Access != old.Access, now.Account, len(g.Tools), len(now.Tools))
	}
	if ev := lastEvent(now); !strings.HasPrefix(ev, "signed in again on ") || !strings.Contains(ev, "its tools are as they were") {
		t.Fatalf("event: %q", ev)
	}
	if !slices.Contains(e.vendor.Revoked, old.Refresh) || !slices.Contains(e.vendor.Revoked, old.Access) {
		t.Fatalf("the replaced sign-in wasn't revoked: %v", e.vendor.Revoked)
	}
	if _, ok, _ := e.kv.Get(ctx, tenantPK(tenant), skRenewal(g.ID)); ok {
		t.Fatal("the new sign-in's token is still held apart")
	}
	if res, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{}); err != nil || strings.Contains(string(res.Result), "Not sent") {
		t.Fatalf("a call with the new sign-in: %s %v", res.Result, err)
	}
}

// A new sign-in in another vendor account changes nothing: it's revoked, and
// the connection works with the sign-in it has.
func TestANewSignInToAnotherAccountChangesNothing(t *testing.T) {
	e := newEnv(t)
	e.vendor.Account = "acct_1"
	g := e.connect()
	ctx := context.Background()
	old, _, _ := e.s.openToken(ctx, g)
	time.Sleep(1100 * time.Millisecond)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": g.Digest}, e.phone, "")
	if st != 200 {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	e.vendor.Account = "acct_2"
	revoked := len(e.vendor.Revoked)
	resp, err := http.Get(out["sign_in_url"].(string))
	if err != nil {
		t.Fatal(err)
	}
	page, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 409 || !strings.Contains(string(page), "another account") || !strings.Contains(string(page), "nothing changed") {
		t.Fatalf("a sign-in to another account: %d %s", resp.StatusCode, page)
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if tok, _, _ := e.s.openToken(ctx, g); tok.Access != old.Access || now.Status != StatusConnected || now.Renewal != nil || now.Account != "acct_1" {
		t.Fatalf("after another account: %s, renewal %v, account %q", now.Status, now.Renewal, now.Account)
	}
	if len(e.vendor.Revoked) != revoked+2 {
		t.Fatalf("the other account's sign-in wasn't revoked: %v", e.vendor.Revoked[revoked:])
	}
	if res, err := e.kiff.Call(ctx, "k_taf", "list_charges", map[string]any{}); err != nil || strings.Contains(string(res.Result), "Not sent") {
		t.Fatalf("a call after: %s %v", res.Result, err)
	}
}

// A new sign-in nobody confirms goes at its deadline, revoked at the vendor,
// and the connection stays as it was: it never acts on nothing meanwhile.
func TestAnUnconfirmedNewSignInGoesAndTheConnectionStays(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	ctx := context.Background()
	old, _, _ := e.s.openToken(ctx, g)
	time.Sleep(1100 * time.Millisecond)
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve", map[string]string{"digest": g.Digest}, e.phone, "")
	if st != 200 {
		t.Fatalf("sign in again: %d %v", st, out)
	}
	_ = e.signIn(out["sign_in_url"].(string)) // signed in, never confirmed
	held, err := e.s.loadGrant(ctx, tenant, g.ID)
	if err != nil || held.Renewal == nil || held.Renewal.ConfirmHash == "" {
		t.Fatalf("signed in: %v %+v", err, held.Renewal)
	}
	fresh, _ := e.s.openRenewal(ctx, held)
	e.s.Now = func() time.Time { return time.Now().Add(signInTTL + time.Minute) }
	if e.s.overdue(held) {
		t.Fatal("a connected grant with a new sign-in past its deadline acts on nothing")
	}
	if n, err := e.s.Sweep(ctx); err != nil || n != 0 {
		t.Fatalf("sweep: %d %v", n, err)
	}
	now, _ := e.s.loadGrant(ctx, tenant, g.ID)
	if now.Status != StatusConnected || now.Renewal != nil || len(now.Tools) != len(g.Tools) || !strings.Contains(lastEvent(now), "the connection stays as it was") {
		t.Fatalf("after the deadline: %s, renewal %v, %d tools, %q", now.Status, now.Renewal, len(now.Tools), lastEvent(now))
	}
	if !slices.Contains(e.vendor.Revoked, fresh.Refresh) || slices.Contains(e.vendor.Revoked, old.Refresh) {
		t.Fatalf("revoked: %v", e.vendor.Revoked)
	}
	if tok, _, _ := e.s.openToken(ctx, g); tok.Refresh != old.Refresh {
		t.Fatal("the sign-in in use changed")
	}
}

// Signing in again isn't offered where planque can't hold it to the same
// account (an app through a connection provider), or where the sign-in is
// another connection's: each says what to do instead.
func TestSigningInAgainSaysWhereItCant(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ph, _ := profile.Digest(g.Request.Profile)
	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+g.ID+"/approve",
		map[string]any{"digest": g.Digest, "kiff_profile": ProfileDigests{ProfileSHA256: ph, WordsSHA256: strings.Repeat("a", 64)}}, e.phone, "")
	if st != 409 || !strings.Contains(mustJSONString(out), "can't tell which account a new sign-in is") || !strings.Contains(mustJSONString(out), "disconnect it and connect it again") {
		t.Fatalf("a connected managed app: %d %v", st, out)
	}
	if now, _ := e.s.loadGrant(context.Background(), tenant, g.ID); now.Status != StatusConnected || now.Renewal != nil {
		t.Fatalf("after: %s %v", now.Status, now.Renewal)
	}
}

// A Telegram bot's token is replaced by the same bot's new one: Telegram takes
// it, and the bot's webhook moves to a new hook. Another bot's token changes
// nothing.
func TestATelegramTokenIsReplacedByTheSameBotsOnly(t *testing.T) {
	e := newEnv(t)
	f := newTelegramFake(t, e)
	e.enrol()
	ctx := context.Background()
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants", map[string]any{"kind": "api_key", "agent": "agt_taf", "tools": []map[string]any{}, "api": map[string]any{"url": TelegramDescription}}, nil, "k_setup")
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	id := out["id"].(string)
	if st, out := e.giveKey(id, out["digest"].(string), f.token); st != 200 {
		t.Fatalf("token: %d %v", st, out)
	}
	st, out = e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/tools", map[string]any{"tools": []map[string]any{{"tool": "send_message", "name": "send_message", "title": "Answer in Telegram"}}}, nil, "k_setup")
	if st != 200 {
		t.Fatalf("tools: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("approve: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(ctx, tenant, id)
	hook := f.hookURL

	// The bot's token is revoked at BotFather: a new one, for the same bot.
	f.mu.Lock()
	f.token = "123456:rotated-synthetic-token"
	f.others = map[string]string{"654321:another-synthetic-token": `{"id":888,"is_bot":true,"first_name":"Other","username":"other_bot"}`}
	f.mu.Unlock()
	st, out = e.giveKey(id, g.Digest, "654321:another-synthetic-token")
	if st != 400 || !strings.Contains(mustJSONString(out), "@other_bot's, not @shop_help_bot's") {
		t.Fatalf("another bot's token: %d %v", st, out)
	}
	if hooks, _ := e.kv.Query(ctx, telegramHookPK, ""); len(hooks) != 1 || f.hookURL != hook {
		t.Fatalf("another bot's token moved the webhook: %d hooks, %q", len(hooks), f.hookURL)
	}
	st, out = e.giveKey(id, g.Digest, f.token)
	if st != 200 || out["status"] != StatusConnected {
		t.Fatalf("the same bot's new token: %d %v", st, out)
	}
	now, _ := e.s.loadGrant(ctx, tenant, id)
	if ev := lastEvent(now); !strings.Contains(ev, "Telegram took it for @shop_help_bot") || now.Telegram == nil || now.Telegram.ID != 777 {
		t.Fatalf("after: %q %+v", ev, now.Telegram)
	}
	if hooks, _ := e.kv.Query(ctx, telegramHookPK, ""); len(hooks) != 1 || f.hookURL == hook || !strings.HasSuffix(f.hookURL, "/telegram/"+now.Telegram.Hook) {
		t.Fatalf("the webhook: %d hooks, %q", len(hooks), f.hookURL)
	}
	if tok, _, _ := e.s.openToken(ctx, now); tok.Access != f.token {
		t.Fatal("the new token isn't the one held")
	}
}
