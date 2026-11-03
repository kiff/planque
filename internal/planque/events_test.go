package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/oauth"
)

type eventFixture struct {
	mu                                   sync.Mutex
	e                                    *env
	now                                  time.Time
	user                                 string
	creates, hookCreates, deletes, sends int
	status                               int
	message                              door.Message
	definition                           connector.TriggerType
}

func newEventFixture(t *testing.T) (*env, *eventFixture, *Grant) {
	e := newInProcessEnv(t)
	e.enrol()
	ctx := context.Background()
	f := &eventFixture{e: e, now: time.Now().UTC().Truncate(time.Second), status: 202, definition: connector.TriggerType{Slug: "SLACKBOT_CHANNEL_MESSAGE_RECEIVED", Toolkit: "slackbot", Version: "20300101_00", Name: "Channel message", Config: json.RawMessage(`{"type":"object"}`), Payload: json.RawMessage(`{"type":"object"}`)}}
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
		switch r.URL.Path {
		case "/api/v3.1/connected_accounts/account_events":
			if r.Method == "DELETE" {
				return profileResponse(200, `{"success":true,"revoke_job_id":"job_events"}`), nil
			}
			raw, _ := json.Marshal(map[string]any{"id": "account_events", "user_id": f.user, "status": "ACTIVE", "toolkit": map[string]string{"slug": "slackbot"}})
			return profileResponse(200, string(raw)), nil
		case "/api/v3.1/triggers_types/SLACKBOT_CHANNEL_MESSAGE_RECEIVED":
			raw, _ := json.Marshal(map[string]any{"slug": f.definition.Slug, "name": f.definition.Name, "version": f.definition.Version, "toolkit": map[string]string{"slug": "slackbot"}, "config": f.definition.Config, "payload": f.definition.Payload})
			return profileResponse(200, string(raw)), nil
		case "/api/v3.1/webhook_subscriptions":
			if r.Method == "GET" {
				return profileResponse(200, `{"items":[]}`), nil
			}
			f.hookCreates++
			return profileResponse(201, `{"id":"hook_events","webhook_url":"https://planque.example/connectors/composio/webhook","version":"V3","enabled_events":["composio.trigger.message"],"secret":"synthetic-webhook-secret"}`), nil
		case "/api/v3.1/trigger_instances/SLACKBOT_CHANNEL_MESSAGE_RECEIVED/upsert":
			f.creates++
			return profileResponse(201, `{"trigger_id":"trigger_events"}`), nil
		case "/api/v3.1/trigger_instances/manage/trigger_events":
			f.deletes++
			return profileResponse(200, `{"trigger_id":"trigger_events"}`), nil
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
	ask := map[string]any{"account": g.Account, "definition": f.definition, "config": map[string]bool{"is_bot_message": false}, "channels": []string{"CDEMO123"}}
	if st, _ := f.e.do("POST", path, ask, nil, ""); st != 401 || f.creates != 0 || f.hookCreates != 0 {
		t.Fatal("trigger/project webhook created without device approval")
	}
	ask["account"] = "another_account"
	if st, _ := f.e.do("POST", path, ask, f.e.phone, ""); st != 409 || f.creates != 0 || f.hookCreates != 0 {
		t.Fatal("a signed request changed the connection account")
	}
	ask["account"] = g.Account
	if st, out := f.e.do("POST", path, ask, f.e.phone, ""); st != 201 {
		t.Fatalf("enable: %d %v", st, out)
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
	if m.Tenant != tenant || m.Agent != "agt_taf" || m.Grant != g.ID || m.Account != "account_events" || m.Channel != "CDEMO123" || m.Sender != "UDEMO123" || m.Thread != "1799999999.001" || m.Timestamp != "2027-01-15T08:00:00.001Z" {
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
	path := "/v1/accounts/" + tenant + "/grants/" + g.ID + "/triggers/trigger_events"
	if st, out := e.do("DELETE", path, map[string]string{"account": g.Account}, e.phone, ""); st != 200 {
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
