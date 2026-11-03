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
}

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
	message, ok := slackDoorMessage(g, route, event)
	if !ok {
		w.WriteHeader(202)
		return
	}
	key := store.Item{PK: deliveryPK, SK: message.ID, Ver: 1, TTL: s.now().Add(7 * 24 * time.Hour).Unix()}
	body, _ := json.Marshal(message)
	sealed, err := s.Sealer.Seal(r.Context(), body, map[string]string{"kind": "door_event", "id": message.ID})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	record := eventDelivery{Tenant: g.Tenant, Grant: g.ID, Provider: provider, Trigger: route.ID, Sealed: sealed, State: "pending"}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(g.Tenant), SK: "DOOR_EVENT#" + message.ID, TTL: key.TTL}, map[string]string{"event": message.ID}, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	if err := s.KV.Put(r.Context(), store.Item{PK: pendingEventsPK, SK: message.ID, Doc: []byte("{}"), TTL: key.TTL}, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	if err := putDoc(r.Context(), s.KV, key, record, store.Cond{NotExists: true}); err != nil && !errors.Is(err, store.ErrCondition) {
		s.fail(w, r, err)
		return
	}
	retry, err := s.deliverEvent(r.Context(), key.SK)
	if err != nil || retry {
		writeErr(w, 503, "the door event is pending delivery")
		return
	}
	w.WriteHeader(202)
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
	return door.Message{V: "planque.door.v1", ID: "evt_" + hex.EncodeToString(sum[:]), Tenant: g.Tenant, Agent: g.Request.Agent, Grant: g.ID, Provider: route.Provider, Account: route.Account, Trigger: route.ID, Slug: route.Definition.Slug, Channel: p.Channel, Sender: p.User, Text: p.Text, Thread: p.Thread, Timestamp: timestamp}, true
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
	if _, err := getDoc(ctx, s.KV, triggerPK(delivery.Provider, delivery.Trigger), "-", &route); err != nil {
		return true, err
	}
	if !route.Active || route.Tenant != g.Tenant || route.Grant != g.ID || g.Status == StatusDisconnected || g.Status == StatusFailed || g.Status == StatusNeedsReconnect || !g.SignedIn || s.overdue(g) {
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
	cleartext, err := s.Sealer.Open(ctx, delivery.Sealed, map[string]string{"kind": "door_event", "id": id})
	if err != nil {
		return true, err
	}
	defer clear(cleartext)
	var message door.Message
	if json.Unmarshal(cleartext, &message) != nil || message.ID != id || message.Tenant != g.Tenant || message.Agent != g.Request.Agent || message.Account != route.Account || message.Trigger != route.ID || !slices.Contains(route.Channels, message.Channel) {
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
	if s.DoorSink != nil {
		result, err = s.DoorSink.Deliver(ctx, message)
		if err != nil {
			result.Retry = true
		}
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
			_, _, _ = s.KV.Delete(ctx, pendingEventsPK, it.SK)
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
