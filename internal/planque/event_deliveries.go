package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/store"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type eventDelivery struct {
	Tenant   string `json:"tenant"`
	Grant    string `json:"grant"`
	Provider string `json:"provider"`
	Trigger  string `json:"trigger"`
	Sealed   []byte `json:"sealed,omitempty"`
	State    string `json:"state"`
	Attempts int    `json:"attempts"`
	Status   int    `json:"status,omitempty"`
	// Kind is "app_event" for planque.event.v1; empty for a door message.
	Kind string `json:"kind,omitempty"`
}

// maxEventPayload bounds an app event's data. A larger one is dropped, and
// the grant notes it: a cut payload would not be the app's event.
const maxEventPayload = 64 << 10

func (s *Server) connectorWebhook(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if s.Connectors[provider] == nil {
		writeErr(w, 404, "no such event provider")
		return
	}
	raw, err := io.ReadAll(io.LimitReader(r.Body, (1<<20)+1))
	if err != nil || len(raw) > 1<<20 {
		writeErr(w, 400, "the event body exceeds its limit")
		return
	}
	var hook webhookRecord
	if _, err := getDoc(r.Context(), s.KV, webhookPK(provider), "-", &hook); err != nil {
		writeErr(w, 503, "the project webhook is not configured")
		return
	}
	secret, err := s.Sealer.Open(r.Context(), hook.Sealed, map[string]string{"kind": "connector_webhook", "provider": provider})
	if err != nil {
		writeErr(w, 503, "the project webhook cannot be verified")
		return
	}
	defer clear(secret)
	if err := connector.VerifyWebhook(string(secret), r.Header.Get("webhook-id"), r.Header.Get("webhook-timestamp"), r.Header.Get("webhook-signature"), raw, s.now()); err != nil {
		writeErr(w, 401, "the event signature is invalid or expired")
		return
	}
	event, err := connector.ParseEvent(raw)
	if err != nil {
		writeErr(w, 400, "the event envelope is unsupported")
		return
	}
	var route triggerRecord
	if _, err := getDoc(r.Context(), s.KV, triggerPK(provider, event.Metadata.Trigger), "-", &route); err != nil {
		writeErr(w, 503, "the approved event route is not available")
		return
	}
	if !route.Active {
		w.WriteHeader(202)
		return
	}
	if route.Provider != provider || route.ID != event.Metadata.Trigger || route.Account != event.Metadata.Account || route.User != event.Metadata.User || route.Definition.Slug != event.Metadata.Slug {
		writeErr(w, 403, "the event does not match its approved route")
		return
	}
	g, err := s.loadGrant(r.Context(), route.Tenant, route.Grant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Status == StatusDisconnected || g.Status == StatusFailed || g.Status == StatusNeedsReconnect || !g.SignedIn || s.overdue(g) {
		w.WriteHeader(202)
		return
	}
	if !slices.Contains(g.TriggerIDs, route.ID) {
		writeErr(w, 503, "the approved event route is still being registered")
		return
	}
	var id, kind string
	var body []byte
	if route.AppEvent {
		e, ok := appEvent(g, route, event)
		if !ok {
			if len(event.Data) > maxEventPayload {
				g.event(s.now(), "an app event over 64 KB was dropped ("+route.Definition.Slug+")")
				if err := s.saveGrant(r.Context(), g); err != nil {
					s.fail(w, r, err)
					return
				}
			}
			w.WriteHeader(202)
			return
		}
		id, kind = e.ID, "app_event"
		body, _ = json.Marshal(e)
	} else {
		message, ok := slackDoorMessage(g, route, event)
		if !ok {
			w.WriteHeader(202)
			return
		}
		id = message.ID
		body, _ = json.Marshal(message)
	}
	retry, err := s.enqueueDoor(r.Context(), g, provider, route.ID, id, kind, body)
	clear(body)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if retry {
		writeErr(w, 503, "the event is pending delivery")
		return
	}
	w.WriteHeader(202)
}

// sealContext binds a sealed delivery to its kind and id.
func sealContext(kind, id string) map[string]string {
	if kind == "app_event" || kind == "telegram" {
		return map[string]string{"kind": kind, "id": id}
	}
	return map[string]string{"kind": "door_event", "id": id}
}

// appEvent is a connected app's event for the taf whose agent the grant
// names (planque.event.v1). Its id is the same for the provider's
// redelivery of the same event, so tafs starts nothing twice.
func appEvent(g *Grant, route triggerRecord, e connector.Event) (door.AppEvent, bool) {
	if !route.AppEvent || len(e.Data) > maxEventPayload || !json.Valid(e.Data) || route.Definition.Slug != e.Metadata.Slug {
		return door.AppEvent{}, false
	}
	sum := sha256.Sum256([]byte("app_event\n" + route.Provider + "\n" + g.Tenant + "\n" + g.ID + "\n" + route.ID + "\n" + e.ID))
	at := e.Timestamp
	if _, err := time.Parse(time.RFC3339Nano, at); err != nil {
		at = ""
	}
	return door.AppEvent{V: "planque.event.v1", ID: "evt_" + hex.EncodeToString(sum[:]), Tenant: g.Tenant, Agent: g.Request.Agent, Grant: g.ID,
		Provider: route.Provider, Account: route.Account, Subscription: route.ID, Event: route.Definition.Slug, Payload: e.Data, Timestamp: at,
		Key: eventKey(g, route.Definition.Slug, e.Data)}, true
}

func slackDoorMessage(g *Grant, route triggerRecord, e connector.Event) (door.Message, bool) {
	if !supportedTrigger(route.Definition) {
		return door.Message{}, false
	}
	var p struct {
		Channel string `json:"channel"`
		User    string `json:"user"`
		Text    string `json:"text"`
		TS      string `json:"ts"`
		Thread  string `json:"thread_ts"`
		Bot     string `json:"bot_id"`
		Subtype string `json:"subtype"`
	}
	if json.Unmarshal(e.Data, &p) != nil || !channelID.MatchString(p.Channel) || !senderID.MatchString(p.User) || !slackTime.MatchString(p.TS) || p.Thread != "" && !slackTime.MatchString(p.Thread) || p.Bot != "" || p.Subtype != "" && p.Subtype != "file_share" && p.Subtype != "thread_broadcast" || !slices.Contains(route.Channels, p.Channel) || len(p.Text) > 128<<10 {
		return door.Message{}, false
	}
	if strings.HasPrefix(p.Channel, "D") != (route.Definition.Slug == "SLACKBOT_DIRECT_MESSAGE_RECEIVED") {
		return door.Message{}, false
	}
	sum := sha256.Sum256([]byte(route.Provider + "\n" + g.Tenant + "\n" + g.ID + "\n" + route.Account + "\n" + p.Channel + "\n" + p.TS))
	seconds, fraction, _ := strings.Cut(p.TS, ".")
	sec, _ := strconv.ParseInt(seconds, 10, 64)
	nsec, _ := strconv.ParseInt(fraction+strings.Repeat("0", 9-len(fraction)), 10, 64)
	timestamp := time.Unix(sec, nsec).UTC().Format(time.RFC3339Nano)
	return door.Message{V: "planque.door.v1", ID: "evt_" + hex.EncodeToString(sum[:]), Tenant: g.Tenant, Agent: g.Request.Agent, Grant: g.ID, Provider: route.Provider, Account: route.Account, Trigger: route.ID, Slug: route.Definition.Slug, Channel: p.Channel, Sender: p.User, Text: p.Text, TS: p.TS, Thread: p.Thread, Timestamp: timestamp}, true
}
func (s *Server) deliverEvent(ctx context.Context, id string) (bool, error) {
	var delivery eventDelivery
	it, err := getDoc(ctx, s.KV, deliveryPK, id, &delivery)
	if err != nil {
		return true, err
	}
	if delivery.State != "pending" {
		_, _, err := s.KV.Delete(ctx, pendingEventsPK, id)
		return false, err
	}
	if it.Lease > s.now().UnixMilli() {
		return true, nil
	}
	g, err := s.loadGrant(ctx, delivery.Tenant, delivery.Grant)
	if errors.Is(err, errNotFound) {
		for _, key := range [][2]string{{deliveryPK, id}, {pendingEventsPK, id}, {tenantPK(delivery.Tenant), "DOOR_EVENT#" + id}} {
			if _, _, err := s.KV.Delete(ctx, key[0], key[1]); err != nil {
				return true, err
			}
		}
		return false, nil
	}
	if err != nil {
		return true, err
	}
	var route triggerRecord
	closed := false
	if delivery.Kind == "telegram" {
		// A Telegram door's route is its grant's current hook.
		closed = s.telegramClosed(g, delivery.Trigger)
	} else if _, err := getDoc(ctx, s.KV, triggerPK(delivery.Provider, delivery.Trigger), "-", &route); err != nil {
		return true, err
	} else {
		closed = !route.Active || route.Tenant != g.Tenant || route.Grant != g.ID || g.Status == StatusDisconnected || g.Status == StatusFailed || g.Status == StatusNeedsReconnect || !g.SignedIn || s.overdue(g)
	}
	if closed {
		delivery.State = "dropped"
		delivery.Sealed = nil
		it.Ver++
		it.Lease = 0
		if err := putDoc(ctx, s.KV, it, delivery, store.Cond{}); err != nil {
			return true, err
		}
		_, _, err := s.KV.Delete(ctx, pendingEventsPK, id)
		return false, err
	}
	previous := it.Ver
	it.Ver++
	it.Lease = s.now().Add(30 * time.Second).UnixMilli()
	if err := putDoc(ctx, s.KV, it, delivery, store.Cond{Ver: &previous}); errors.Is(err, store.ErrCondition) {
		return true, nil
	} else if err != nil {
		return true, err
	}
	cleartext, err := s.Sealer.Open(ctx, delivery.Sealed, sealContext(delivery.Kind, id))
	if err != nil {
		return true, err
	}
	defer clear(cleartext)
	var message door.Message
	var app door.AppEvent
	var tg door.TelegramMessage
	var valid bool
	if delivery.Kind == "telegram" {
		valid = json.Unmarshal(cleartext, &tg) == nil && tg.ID == id && tg.Tenant == g.Tenant && tg.Agent == g.Request.Agent && tg.Grant == g.ID
	} else if delivery.Kind == "app_event" {
		valid = json.Unmarshal(cleartext, &app) == nil && app.ID == id && app.Tenant == g.Tenant && app.Agent == g.Request.Agent && app.Account == route.Account && app.Subscription == route.ID && route.AppEvent
	} else {
		valid = json.Unmarshal(cleartext, &message) == nil && message.ID == id && message.Tenant == g.Tenant && message.Agent == g.Request.Agent && message.Account == route.Account && message.Trigger == route.ID && slices.Contains(route.Channels, message.Channel)
	}
	if !valid {
		delivery.State = "dropped"
		delivery.Sealed = nil
		previous = it.Ver
		it.Ver++
		it.Lease = 0
		if err := putDoc(ctx, s.KV, it, delivery, store.Cond{Ver: &previous}); err != nil {
			return true, err
		}
		_, _, err := s.KV.Delete(ctx, pendingEventsPK, id)
		return false, err
	}
	result := door.Result{Retry: true}
	switch tgSink, _ := s.DoorSink.(door.TelegramSink); {
	case delivery.Kind == "telegram" && tgSink != nil:
		result, err = tgSink.DeliverTelegram(ctx, tg)
	case delivery.Kind == "app_event" && s.EventSink != nil:
		result, err = s.EventSink.DeliverEvent(ctx, app)
	case delivery.Kind == "" && s.DoorSink != nil:
		result, err = s.DoorSink.Deliver(ctx, message)
	}
	if err != nil {
		result.Retry = true
	}
	delivery.Attempts++
	delivery.Status = result.Status
	retry := result.Retry
	switch {
	case result.Taken:
		delivery.State = "delivered"
		delivery.Sealed = nil
	case !retry:
		delivery.State = "dropped"
		delivery.Sealed = nil
	}
	previous = it.Ver
	it.Ver++
	it.Lease = 0
	if retry {
		it.Lease = s.now().Add(time.Minute).UnixMilli()
	}
	if err := putDoc(ctx, s.KV, it, delivery, store.Cond{Ver: &previous}); err != nil {
		return true, err
	}
	if !retry {
		_, _, err := s.KV.Delete(ctx, pendingEventsPK, id)
		if err != nil {
			return true, err
		}
	}
	return retry, nil
}

// pendingStamp is when a pending index was written; an index from before it
// had one reads as old.
type pendingStamp struct {
	At int64 `json:"at,omitempty"`
}

// pendingGrace is how long the sweep leaves a pending index whose delivery
// is not written yet: far longer than a webhook takes to write both.
const pendingGrace = 5 * time.Minute

func (s *Server) deliverEvents(ctx context.Context) (int, error) {
	items, err := s.KV.Query(ctx, pendingEventsPK, "")
	if err != nil {
		return 0, err
	}
	done, attempted := 0, 0
	for _, it := range items {
		if attempted >= 100 {
			break
		}
		var d eventDelivery
		row, err := getDoc(ctx, s.KV, deliveryPK, it.SK, &d)
		if errors.Is(err, errNotFound) {
			// The webhook writes this index before the delivery: a young index
			// with no delivery yet is one being enqueued, not a stale one.
			var stamp pendingStamp
			_ = json.Unmarshal(it.Doc, &stamp)
			if s.now().Sub(time.UnixMilli(stamp.At)) >= pendingGrace {
				_, _, _ = s.KV.Delete(ctx, pendingEventsPK, it.SK)
			}
			continue
		}
		if err != nil {
			return done, err
		}
		if row.Lease > s.now().UnixMilli() {
			continue
		}
		attempted++
		retry, err := s.deliverEvent(ctx, it.SK)
		if err != nil {
			return done, err
		}
		if !retry {
			done++
		}
	}
	return done, nil
}
