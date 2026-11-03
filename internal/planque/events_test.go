package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

type eventFixture struct {
	mu                                   sync.Mutex
	e                                    *env
	now                                  time.Time
	user                                 string
	creates, hookCreates, deletes, sends int
	failCreate                           int            // the create that fails upstream, from 1
	failDeletes                          int            // upstream deletes that fail before one succeeds
	reuse                                map[int]string // the id the provider answers for a create, from 1
	configs                              []map[string]any
	deleted                              []string
	listCalls                            int
	listReadOnly                         bool
	status                               int
	message                              door.Message
	definition                           connector.TriggerType
}

func newEventFixture(t *testing.T) (*env, *eventFixture, *Grant) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	f := &eventFixture{e: e, now: time.Now().UTC().Truncate(time.Second), status: 202, listReadOnly: true, definition: connector.TriggerType{Slug: "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", Toolkit: "slackbot", Version: "20300101_00", Name: "Channel message", Config: json.RawMessage(`{"type":"object"}`), Payload: json.RawMessage(`{"type":"object"}`)}}
	e.s.Now = func() time.Time { return f.now }
	g := &Grant{ID: "grt_" + oauth.Random(12), Tenant: tenant, Status: StatusListed, SignedIn: true, Account: "account_events", HeldUntil: f.now.Add(holdTTL), Request: Request{Version: "planque.grant.v1", Tenant: tenant, Server: composio.API, Agent: "agt_taf", AskedBy: "agt_setup", Connector: &ConnectorRequest{Provider: "composio", Toolkit: "slackbot"}, API: &APIRequest{Base: composio.API}}}
	g.Request.Grant = g.ID
	identity := e.s.connectorIdentity(g)
	f.user = identity.User
	g.Connector = &ConnectorState{User: identity.User, Account: "account_events"}
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if r.URL.Host == "receiver.example" {
			raw, _ := io.ReadAll(r.Body)
			if connector.VerifyWebhook("synthetic-handoff", r.Header.Get("Planque-Event-Id"), r.Header.Get("Planque-Event-Timestamp"), r.Header.Get("Planque-Event-Signature"), raw, f.now) != nil {
				t.Error("handoff signature failed")
			}
			if err := json.Unmarshal(raw, &f.message); err != nil {
				t.Error(err)
			}
			f.sends++
			return profileResponse(f.status, `{}`), nil
		}
		if r.URL.Host != "backend.composio.dev" || r.Header.Get("x-api-key") != "synthetic-project-key" {
			t.Fatal("unexpected provider endpoint/credential")
		}
		switch path := r.URL.Path; {
		case path == "/api/v3.1/connected_accounts/account_events":
			if r.Method == "DELETE" {
				return profileResponse(200, `{"success":true,"revoke_job_id":"job_events"}`), nil
			}
			raw, _ := json.Marshal(map[string]any{"id": "account_events", "user_id": f.user, "status": "ACTIVE", "toolkit": map[string]string{"slug": "slackbot"}})
			return profileResponse(200, string(raw)), nil
		case strings.HasPrefix(path, "/api/v3.1/triggers_types/SLACKBOT_"):
			d := f.definition
			d.Slug = strings.TrimPrefix(path, "/api/v3.1/triggers_types/")
			raw, _ := json.Marshal(map[string]any{"slug": d.Slug, "name": d.Name, "version": d.Version, "toolkit": map[string]string{"slug": "slackbot"}, "config": d.Config, "payload": d.Payload})
			return profileResponse(200, string(raw)), nil
		case path == "/api/v3.1/webhook_subscriptions":
			if r.Method == "GET" {
				return profileResponse(200, `{"items":[]}`), nil
			}
			f.hookCreates++
			return profileResponse(201, `{"id":"hook_events","webhook_url":"https://planque.example/connectors/composio/webhook","version":"V3","enabled_events":["composio.trigger.message"],"secret":"synthetic-webhook-secret"}`), nil
		case strings.HasPrefix(path, "/api/v3.1/trigger_instances/SLACKBOT_") && strings.HasSuffix(path, "/upsert"):
			f.creates++
			if f.creates == f.failCreate {
				return profileResponse(400, `{"error":"refused"}`), nil
			}
			var body struct {
				Config map[string]any `json:"trigger_config"`
			}
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			f.configs = append(f.configs, body.Config)
			id := "trigger_events"
			if f.creates > 1 {
				id += "_" + strconv.Itoa(f.creates)
			}
			if reused, ok := f.reuse[f.creates]; ok {
				id = reused
			}
			return profileResponse(201, `{"trigger_id":"`+id+`"}`), nil
		case strings.HasPrefix(path, "/api/v3.1/trigger_instances/manage/"):
			if f.failDeletes > 0 {
				f.failDeletes--
				return profileResponse(500, `{"error":"unavailable"}`), nil
			}
			f.deletes++
			id := strings.TrimPrefix(path, "/api/v3.1/trigger_instances/manage/")
			f.deleted = append(f.deleted, id)
			return profileResponse(200, `{"trigger_id":"`+id+`"}`), nil
		case path == "/api/v3.1/tools/SLACKBOT_LIST_ALL_CHANNELS":
			tags := []string{}
			if f.listReadOnly {
				tags = append(tags, "readOnlyHint")
			}
			raw, _ := json.Marshal(map[string]any{"slug": "SLACKBOT_LIST_ALL_CHANNELS", "name": "List all channels", "description": "Lists channels.", "toolkit": map[string]string{"slug": "slackbot"}, "version": "20261101_00", "input_parameters": map[string]any{"type": "object", "properties": map[string]any{"limit": map[string]string{"type": "integer"}, "types": map[string]string{"type": "string"}, "cursor": map[string]string{"type": "string"}, "exclude_archived": map[string]string{"type": "boolean"}}}, "tags": tags})
			return profileResponse(200, string(raw)), nil
		case path == "/api/v3.1/tools/execute/SLACKBOT_LIST_ALL_CHANNELS":
			f.listCalls++
			var body struct {
				Arguments map[string]any `json:"arguments"`
			}
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			if body.Arguments["exclude_archived"] != true || body.Arguments["types"] != "public_channel,private_channel" {
				t.Errorf("channel list asked for archived or other conversations: %v", body.Arguments)
			}
			if body.Arguments["cursor"] == nil {
				return profileResponse(200, `{"successful":true,"data":{"channels":[{"id":"CSALES01","name":"sales","is_private":false,"topic":{"value":"not for the phone"}},{"id":"not-an-id","name":"bad"}],"response_metadata":{"next_cursor":"page2"}}}`), nil
			}
			return profileResponse(200, `{"successful":true,"data":{"channels":[{"id":"GSECRET1","name":"leads","is_private":true}],"response_metadata":{"next_cursor":""}}}`), nil
		default:
			t.Fatalf("unexpected provider API: %s", r.URL)
			return nil, nil
		}
	})}
	e.s.Connectors = map[string]connector.Provider{"composio": &composio.Client{HTTP: client, Secret: func(context.Context) (string, error) { return "synthetic-project-key", nil }}}
	e.s.DoorSink = &door.HTTP{URL: "https://receiver.example/v1/doors/planque", Secret: func(context.Context) (string, error) { return "synthetic-handoff", nil }, Client: client, Now: func() time.Time { return f.now }}
	return e, f, g
}
func (f *eventFixture) enable(t *testing.T, g *Grant) {
	t.Helper()
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CDEMO123"}}
	if st, _ := f.e.do("POST", path, ask, nil, ""); st != 401 || f.creates != 0 || f.hookCreates != 0 {
		t.Fatal("trigger/project webhook created without device approval")
	}
	ask["account"] = "another_account"
	if st, _ := f.e.do("POST", path, ask, f.e.phone, ""); st != 409 || f.creates != 0 || f.hookCreates != 0 {
		t.Fatal("a signed request changed the connection account")
	}
	ask["account"] = g.Account
	if st, out := f.e.do("POST", path, ask, f.e.phone, ""); st != 201 || out["id"] != "trigger_events" || out["definition"].(map[string]any)["version"] != f.definition.Version {
		t.Fatalf("enable: %d %v", st, out)
	}
	if len(f.configs) != 1 || f.configs[0]["channel_id"] != "CDEMO123" || f.configs[0]["is_bot_message"] != false {
		t.Fatalf("the upstream config is not the one approved channel: %v", f.configs)
	}
	saved, _ := f.e.s.loadGrant(context.Background(), tenant, g.ID)
	if len(saved.TriggerIDs) != 1 || !saved.HeldUntil.IsZero() {
		t.Fatal("active door was treated as an unused sign-in")
	}
	var hook webhookRecord
	_, err := getDoc(context.Background(), f.e.kv, webhookPK("composio"), "-", &hook)
	if err != nil || bytes.Contains(hook.Sealed, []byte("synthetic-webhook-secret")) {
		t.Fatal("project signing secret was not sealed")
	}
}
func (f *eventFixture) event(t *testing.T, change func(map[string]any), tamper bool) int {
	t.Helper()
	body := map[string]any{"id": "upstream_event", "type": "composio.trigger.message", "metadata": map[string]any{"trigger_id": "trigger_events", "trigger_slug": f.definition.Slug, "connected_account_id": "account_events", "user_id": f.user}, "data": map[string]any{"channel": "CDEMO123", "user": "UDEMO123", "text": "Please check this.", "ts": "1800000000.001", "thread_ts": "1799999999.001"}}
	if change != nil {
		change(body)
	}
	raw, _ := json.Marshal(body)
	timestamp := fmtUnix(f.now)
	sig := connector.SignWebhook("synthetic-webhook-secret", "delivery_one", timestamp, raw)
	if tamper {
		raw = append(raw, ' ')
	}
	req, _ := http.NewRequest("POST", f.e.url+"/connectors/composio/webhook", bytes.NewReader(raw))
	req.Header.Set("webhook-id", "delivery_one")
	req.Header.Set("webhook-timestamp", timestamp)
	req.Header.Set("webhook-signature", sig)
	res, err := f.e.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	return res.StatusCode
}
func fmtUnix(t time.Time) string { return strconv.FormatInt(t.Unix(), 10) }
func TestSignedEventsRouteExactAccountChannelAndDeliverOnce(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g)
	if st := f.event(t, nil, true); st != 401 || f.sends != 0 {
		t.Fatal("tampered event was forwarded")
	}
	for _, field := range []string{"connected_account_id", "user_id", "trigger_slug"} {
		st := f.event(t, func(v map[string]any) { v["metadata"].(map[string]any)[field] = "another_value" }, false)
		if st != 403 || f.sends != 0 {
			t.Fatalf("mismatched %s forwarded: %d", field, st)
		}
	}
	for _, field := range []string{"channel", "bot_id"} {
		st := f.event(t, func(v map[string]any) { v["data"].(map[string]any)[field] = "OUTSIDE" }, false)
		if st != 202 || f.sends != 0 {
			t.Fatalf("unapproved %s forwarded: %d", field, st)
		}
	}
	if st := f.event(t, nil, false); st != 202 || f.sends != 1 {
		t.Fatalf("event: status%d sends%d", st, f.sends)
	}
	m := f.message
	if m.Tenant != tenant || m.Agent != "agt_taf" || m.Grant != g.ID || m.Account != "account_events" || m.Channel != "CDEMO123" || m.Sender != "UDEMO123" || m.Thread != "1799999999.001" || m.TS != "1800000000.001" || m.Timestamp != "2027-01-15T08:00:00.001Z" {
		t.Fatalf("wrong server-bound event: %+v", m)
	}
	if st := f.event(t, nil, false); st != 202 || f.sends != 1 {
		t.Fatal("delivery replay was sent twice")
	}
	// Even a provider retry with another envelope id is the same Slack message.
	if st := f.event(t, func(v map[string]any) { v["id"] = "another_envelope" }, false); st != 202 || f.sends != 1 {
		t.Fatal("provider retry escaped stable message deduplication")
	}
	var record eventDelivery
	if _, err := getDoc(context.Background(), e.kv, deliveryPK, m.ID, &record); err != nil || record.State != "delivered" || len(record.Sealed) != 0 {
		t.Fatal("completed delivery retained message content or outcome was not saved")
	}
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers/trigger_events"
	if st, out := e.do("DELETE", path, map[string]string{"account": g.Account}, e.phone, ""); st != 200 || f.deletes != 1 {
		t.Fatalf("disable: %d %v", st, out)
	}
	if st := f.event(t, func(v map[string]any) { v["data"].(map[string]any)["ts"] = "1800000001.001" }, false); st != 202 || f.sends != 1 {
		t.Fatal("disabled trigger forwarded")
	}
}
func TestDoorTransientRetryAndAccountErase(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g)
	f.status = 503
	if st := f.event(t, nil, false); st != 503 || f.sends != 1 {
		t.Fatal("transient failure was not pending")
	}
	id := f.message.ID
	if st := f.event(t, nil, false); st != 503 || f.sends != 1 {
		t.Fatal("delivery lease did not suppress a concurrent retry")
	}
	f.now = f.now.Add(61 * time.Second)
	f.status = 202
	if n, err := e.s.deliverEvents(context.Background()); err != nil || n != 1 || f.sends != 2 || f.message.ID != id {
		t.Fatal("durable retry changed message identity or did not complete")
	}
	var enrollment Enrolment
	_, _ = getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &enrollment)
	if _, _, err := e.s.eraseAll(context.Background(), tenant, enrollment, "test-device"); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := e.kv.Get(context.Background(), deliveryPK, id); err != nil || ok {
		t.Fatal("account erasure retained a global delivery")
	}
	if _, ok, err := e.kv.Get(context.Background(), triggerPK("composio", "trigger_events"), "-"); err != nil || ok {
		t.Fatal("account erasure retained global trigger routing")
	}
}

func TestConcurrentSignedEventRetriesForwardOneMessage(t *testing.T) {
	_, f, g := newEventFixture(t)
	f.enable(t, g)
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status := f.event(t, nil, false)
			if status != 202 && status != 503 {
				t.Errorf("concurrent event status %d", status)
			}
		}()
	}
	wg.Wait()
	if f.sends != 1 {
		t.Fatalf("concurrent replay forwarded %d messages", f.sends)
	}
}
func TestDisablingTriggerCancelsPendingDelivery(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g)
	f.status = 503
	if status := f.event(t, nil, false); status != 503 {
		t.Fatal("delivery was not queued")
	}
	id := f.message.ID
	// The body is optional: the phone removes a door with a bare signed DELETE.
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers/trigger_events"
	if st, out := e.do("DELETE", path, nil, e.phone, ""); st != 200 || f.deletes != 1 {
		t.Fatalf("disable: %d %v", st, out)
	}
	f.now = f.now.Add(61 * time.Second)
	f.status = 202
	if _, err := e.s.deliverEvents(context.Background()); err != nil {
		t.Fatal(err)
	}
	if f.sends != 1 {
		t.Fatal("a disabled route's queued message was sent")
	}
	var record eventDelivery
	if _, err := getDoc(context.Background(), e.kv, deliveryPK, id, &record); err != nil || record.State != "dropped" || len(record.Sealed) != 0 {
		t.Fatal("cancelled delivery retained its content")
	}
}

func TestEventApprovalsExportWithoutProjectOrHandoffSecrets(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g)
	status, out := e.do("GET", "/v1/accounts/"+tenant+"/export", nil, e.phone, "")
	if status != 200 || len(out["triggers"].([]any)) != 1 {
		t.Fatalf("oversight export: %d %v", status, out)
	}
	raw, _ := json.Marshal(out)
	for _, secret := range []string{"synthetic-project-key", "synthetic-webhook-secret", "synthetic-handoff"} {
		if bytes.Contains(raw, []byte(secret)) {
			t.Fatal("a secret escaped in the oversight export")
		}
	}
	record := out["triggers"].([]any)[0].(map[string]any)
	if record["approval"] == nil || record["definition"] == nil || record["account"] != "account_events" {
		t.Fatal("export omitted the signed event authority")
	}
}

func TestDoorSubscriptionsAreOnePerChannelAndNeverWider(t *testing.T) {
	e, f, g := newEventFixture(t)
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	for _, ask := range []map[string]any{
		{"account": g.Account, "slug": "SLACKBOT_REACTION_ADDED", "channels": []string{"CDEMO123"}},
		{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{}},
		{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"DDIRECT1"}},
		{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CDEMO123", "CDEMO123"}},
		{"account": g.Account, "slug": "SLACKBOT_DIRECT_MESSAGE_RECEIVED", "channels": []string{"CDEMO123"}},
		{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CDEMO123"}, "config": map[string]any{}},
	} {
		if st, _ := e.do("POST", path, ask, e.phone, ""); st != 400 || f.creates != 0 {
			t.Fatalf("%v: %d, %d upstream creates", ask, st, f.creates)
		}
	}
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CSALES01", "GLEADS01"}}
	st, out := e.do("POST", path, ask, e.phone, "")
	if st != 201 || out["id"] != "trigger_events" || len(out["triggers"].([]any)) != 2 {
		t.Fatalf("two channels: %d %v", st, out)
	}
	if len(f.configs) != 2 || f.configs[0]["channel_id"] != "CSALES01" || f.configs[1]["channel_id"] != "GLEADS01" || f.configs[0]["is_bot_message"] != false || f.configs[1]["is_bot_message"] != false {
		t.Fatalf("each channel needs its own subscription: %v", f.configs)
	}
	for i, id := range []string{"trigger_events", "trigger_events_2"} {
		var record triggerRecord
		if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", id), "-", &record); err != nil || !record.Active || len(record.Channels) != 1 || record.Channels[0] != ask["channels"].([]string)[i] {
			t.Fatalf("%s routes %v", id, record.Channels)
		}
	}
	// A DM door: no channel filter upstream, its DM ids route delivery.
	ask = map[string]any{"account": g.Account, "slug": "SLACKBOT_DIRECT_MESSAGE_RECEIVED", "channels": []string{"DOWNER01", "DOWNER02"}}
	if st, out := e.do("POST", path, ask, e.phone, ""); st != 201 || len(out["triggers"].([]any)) != 1 {
		t.Fatalf("direct messages: %d %v", st, out)
	}
	if len(f.configs) != 3 || len(f.configs[2]) != 0 {
		t.Fatalf("a DM subscription takes no channel filter: %v", f.configs[2])
	}
	saved, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if len(saved.TriggerIDs) != 3 {
		t.Fatalf("grant names %v", saved.TriggerIDs)
	}
}

func TestDoorSubscriptionFailureTakesBackWhatItMade(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.failCreate = 3
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CONE0001", "CTWO0002", "CTHREE03"}}
	if st, _ := e.do("POST", path, ask, e.phone, ""); st < 400 {
		t.Fatalf("a refused third subscription answered %d", st)
	}
	if !slices.Equal(f.deleted, []string{"trigger_events", "trigger_events_2"}) {
		t.Fatalf("upstream subscriptions left listening: deleted %v", f.deleted)
	}
	for _, id := range []string{"trigger_events", "trigger_events_2"} {
		var record triggerRecord
		if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", id), "-", &record); err != nil || record.Active {
			t.Fatalf("%s still routes deliveries", id)
		}
	}
	saved, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if len(saved.TriggerIDs) != 0 || saved.HeldUntil.IsZero() {
		t.Fatal("a failed approval changed the connection")
	}
	// The channel event of a taken-back subscription goes nowhere.
	if st := f.event(t, func(v map[string]any) { v["data"].(map[string]any)["channel"] = "CONE0001" }, false); st != 202 || f.sends != 0 {
		t.Fatal("a taken-back subscription forwarded")
	}
}

func TestDoorEventReadsAreTheConnectionsOwnAgents(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g)
	e.kiff.Keys["k_stranger"] = kiff.Me{APIKey: true, TenantID: tenant, AgentID: "agt_stranger"}
	e.kiff.Keys["k_account"] = kiff.Me{APIKey: true, TenantID: tenant}
	base := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	for _, path := range []string{base, base + "/types?slug=SLACKBOT_CHANNEL_MESSAGE_RECEIVED"} {
		for _, key := range []string{"k_stranger", "k_account"} {
			if st, _ := e.do("GET", path, nil, nil, key); st != 403 {
				t.Fatalf("%s read %s: %d", key, path, st)
			}
		}
		for _, key := range []string{"k_taf", "k_setup"} {
			if st, out := e.do("GET", path, nil, nil, key); st != 200 {
				t.Fatalf("%s read %s: %d %v", key, path, st, out)
			}
		}
	}
	if st, _ := e.do("GET", "/v1/accounts/"+tenant+"/grants/grt_missing/triggers", nil, nil, "k_taf"); st != 404 {
		t.Fatal("a missing grant was not 404")
	}
}

func TestChannelsListThroughTheConnectionsReadOnlyAction(t *testing.T) {
	e, f, g := newEventFixture(t)
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/channels"
	if st, _ := e.do("GET", path, nil, nil, "k_taf"); st != 401 || f.listCalls != 0 {
		t.Fatal("channels listed without the device's signature")
	}
	st, out := e.do("GET", path, nil, e.phone, "")
	if st != 200 || f.listCalls != 2 {
		t.Fatalf("channels: %d %v (%d pages)", st, out, f.listCalls)
	}
	raw, _ := json.Marshal(out["channels"])
	if string(raw) != `[{"id":"CSALES01","name":"sales","private":false},{"id":"GSECRET1","name":"leads","private":true}]` {
		t.Fatalf("channels: %s", raw)
	}
	f.listReadOnly = false
	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	if st, out := e.do("GET", path, nil, e.phone, ""); st != 503 || f.listCalls != 2 {
		t.Fatalf("a list action that is no longer read-only was run: %d %v", st, out)
	}
}

func TestGrantViewNamesItsEntryAndProvider(t *testing.T) {
	e, _, g := newEventFixture(t)
	v := e.s.view(g)
	if v.Entry != "connector:composio:slackbot" || v.Provider != "composio" {
		t.Fatalf("view: entry %q provider %q", v.Entry, v.Provider)
	}
	g.Request.CatalogEntry = "connector:composio:slackbot"
	if v := e.s.view(g); v.Entry != "connector:composio:slackbot" {
		t.Fatalf("view entry %q", v.Entry)
	}
}

func TestUndoneSubscriptionsTheProviderKeptAreDeletedBySweep(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.failCreate, f.failDeletes = 2, 1
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CONE0001", "CTWO0002"}}
	if st, _ := e.do("POST", path, ask, e.phone, ""); st < 400 || f.deletes != 0 {
		t.Fatalf("a refused second subscription answered %d", st)
	}
	notes, err := e.kv.Query(context.Background(), triggerCleanupPK, "")
	if err != nil || len(notes) != 1 || bytes.Contains(notes[0].Doc, []byte("CONE0001")) {
		t.Fatalf("the kept subscription was not noted for the sweep, or the note holds more: %v %v", notes, err)
	}
	var record triggerRecord
	if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", "trigger_events"), "-", &record); err != nil || record.Active {
		t.Fatal("the kept subscription still routes")
	}
	if n, err := e.s.cleanTriggers(context.Background()); err != nil || n != 1 || !slices.Equal(f.deleted, []string{"trigger_events"}) {
		t.Fatalf("sweep: %d %v, deleted %v", n, err, f.deleted)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 0 {
		t.Fatal("a deleted subscription's note stayed")
	}
}

func TestSweepLeavesAPendingIndexBeingEnqueued(t *testing.T) {
	e, f, _ := newEventFixture(t)
	stamp, _ := json.Marshal(pendingStamp{At: f.now.UnixMilli()})
	if err := e.kv.Put(context.Background(), store.Item{PK: pendingEventsPK, SK: "evt_enqueueing", Doc: stamp}, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	if _, err := e.s.deliverEvents(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := e.kv.Get(context.Background(), pendingEventsPK, "evt_enqueueing"); !ok {
		t.Fatal("the sweep removed the index of a delivery still being written")
	}
	f.now = f.now.Add(pendingGrace)
	if _, err := e.s.deliverEvents(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := e.kv.Get(context.Background(), pendingEventsPK, "evt_enqueueing"); ok {
		t.Fatal("a stale index with no delivery stayed")
	}
}

func TestAFailedApprovalLeavesSubscriptionsThatWereAlreadyOn(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.enable(t, g) // CDEMO123 is on as trigger_events
	before, _, _ := e.kv.Get(context.Background(), triggerPK("composio", "trigger_events"), "-")
	f.reuse, f.failCreate = map[int]string{2: "trigger_events"}, 3
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CDEMO123", "CNEW0002"}}
	if st, _ := e.do("POST", path, ask, e.phone, ""); st < 400 {
		t.Fatalf("a refused second subscription answered %d", st)
	}
	after, _, _ := e.kv.Get(context.Background(), triggerPK("composio", "trigger_events"), "-")
	if !bytes.Equal(before.Doc, after.Doc) || before.Ver != after.Ver {
		t.Fatal("the record of a door that was on was rewritten by a failed request")
	}
	var record triggerRecord
	if _, err := getDoc(context.Background(), e.kv, triggerPK("composio", "trigger_events"), "-", &record); err != nil || !record.Active || slices.Contains(f.deleted, "trigger_events") {
		t.Fatalf("a door that was on was taken back: active %v, deleted %v", record.Active, f.deleted)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 0 {
		t.Fatal("a door that was on was noted for deletion")
	}
	if st := f.event(t, nil, false); st != 202 || f.sends != 1 {
		t.Fatal("the door that was on stopped delivering")
	}
}

func TestReapprovalRetiresTheSweepsNoteToDelete(t *testing.T) {
	e, f, g := newEventFixture(t)
	f.failCreate, f.failDeletes = 2, 1
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers"
	ask := map[string]any{"account": g.Account, "slug": "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "channels": []string{"CDEMO123", "CTWO0002"}}
	if st, _ := e.do("POST", path, ask, e.phone, ""); st < 400 {
		t.Fatalf("a refused second subscription answered %d", st)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 1 {
		t.Fatal("the kept subscription was not noted")
	}
	// Approved again: the provider answers the subscription it kept.
	f.reuse = map[int]string{3: "trigger_events"}
	time.Sleep(1100 * time.Millisecond) // a signed request is told apart by its second
	ask["channels"] = []string{"CDEMO123"}
	if st, out := e.do("POST", path, ask, e.phone, ""); st != 201 || out["id"] != "trigger_events" {
		t.Fatalf("reapproval: %d %v", st, out)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 0 {
		t.Fatal("the reapproved subscription is still noted for deletion")
	}
	// Even a note that outlived its reapproval never deletes a routing subscription.
	_ = putDoc(context.Background(), e.kv, store.Item{PK: triggerCleanupPK, SK: "composio#trigger_events"}, triggerCleanup{Provider: "composio", Trigger: "trigger_events", Tenant: tenant, Grant: g.ID}, store.Cond{})
	if n, err := e.s.cleanTriggers(context.Background()); err != nil || n != 0 || len(f.deleted) != 0 {
		t.Fatalf("sweep deleted a reapproved subscription: %d %v %v", n, err, f.deleted)
	}
	if notes, _ := e.kv.Query(context.Background(), triggerCleanupPK, ""); len(notes) != 0 {
		t.Fatal("an obsolete note stayed")
	}
	if st := f.event(t, nil, false); st != 202 || f.sends != 1 {
		t.Fatal("the reapproved door does not deliver")
	}
}

func TestARootSlackMessageCarriesItsOwnTS(t *testing.T) {
	_, f, g := newEventFixture(t)
	f.enable(t, g)
	if st := f.event(t, func(v map[string]any) { delete(v["data"].(map[string]any), "thread_ts") }, false); st != 202 || f.sends != 1 {
		t.Fatalf("root message: %d, %d sent", st, f.sends)
	}
	raw, _ := json.Marshal(f.message)
	var handoff map[string]any
	_ = json.Unmarshal(raw, &handoff)
	if handoff["ts"] != "1800000000.001" || handoff["thread"] != nil {
		t.Fatalf("a root message's handoff cannot be replied under: %s", raw)
	}
}
