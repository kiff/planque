package planque

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/store"
	"net/http"
	"regexp"
	"slices"
	"time"
)

const deliveryPK = "CONNECTOR#DELIVERIES"
const pendingEventsPK = "CONNECTOR#PENDING_EVENTS"

var channelID = regexp.MustCompile(`^[CGD][A-Z0-9]{2,40}$`)
var senderID = regexp.MustCompile(`^[UW][A-Z0-9]{2,40}$`)
var slackTime = regexp.MustCompile(`^[0-9]{1,12}\.[0-9]{1,9}$`)

type triggerRecord struct {
	Tenant     string                `json:"tenant"`
	Grant      string                `json:"grant"`
	Provider   string                `json:"provider"`
	Account    string                `json:"account"`
	User       string                `json:"user"`
	ID         string                `json:"id"`
	Definition connector.TriggerType `json:"definition"`
	Config     json.RawMessage       `json:"config"`
	Channels   []string              `json:"channels"`
	Approval   Approval              `json:"approval"`
	Active     bool                  `json:"active"`
}
type webhookRecord struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Sealed []byte `json:"sealed"`
}

func triggerPK(provider, id string) string { return "CONNECTOR#" + provider + "#TRIGGER#" + id }
func webhookPK(provider string) string     { return "CONNECTOR#" + provider + "#WEBHOOK" }
func supportedTrigger(t connector.TriggerType) bool {
	return t.Toolkit == "slackbot" && (t.Slug == "SLACKBOT_CHANNEL_MESSAGE_RECEIVED" || t.Slug == "SLACKBOT_DIRECT_MESSAGE_RECEIVED")
}
func (s *Server) eventProvider(g *Grant) (connector.EventProvider, error) {
	p, err := s.provider(g)
	if err != nil {
		return nil, err
	}
	events, ok := p.(connector.EventProvider)
	if !ok {
		return nil, apiError{503, "this provider does not support managed events"}
	}
	return events, nil
}
func (s *Server) ensureProjectWebhook(ctx context.Context, provider string, p connector.EventProvider) error {
	callback := s.BaseURL + "/connectors/" + provider + "/webhook"
	hook, err := p.EnsureWebhook(ctx, callback)
	if err != nil {
		return err
	}
	sealed, err := s.Sealer.Seal(ctx, []byte(hook.Secret), map[string]string{"kind": "connector_webhook", "provider": provider})
	if err != nil {
		return err
	}
	return putDoc(ctx, s.KV, store.Item{PK: webhookPK(provider), SK: "-"}, webhookRecord{ID: hook.ID, URL: callback, Sealed: sealed}, store.Cond{})
}
func (s *Server) triggerType(w http.ResponseWriter, r *http.Request) {
	if _, err := s.agentKey(r, r.PathValue("tenant")); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), r.PathValue("tenant"), r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Request.Connector == nil || g.Request.Connector.Toolkit != "slackbot" {
		writeErr(w, 400, "this connection has no supported door events")
		return
	}
	p, err := s.eventProvider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	definition, err := p.TriggerType(r.Context(), g.Request.Connector.Toolkit, r.URL.Query().Get("slug"), "")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !supportedTrigger(definition) {
		writeErr(w, 400, "this event type is not supported for a door")
		return
	}
	writeJSON(w, 200, map[string]any{"definition": definition, "account": g.Account})
}
func (s *Server) createTrigger(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var ask struct {
		Account    string                `json:"account"`
		Definition connector.TriggerType `json:"definition"`
		Config     json.RawMessage       `json:"config"`
		Channels   []string              `json:"channels"`
	}
	if err := decodeStrict(sr.body, &ask); err != nil {
		s.fail(w, r, err)
		return
	}
	if len(sr.body) > 64<<10 {
		writeErr(w, 400, "the event approval exceeds 64 KB")
		return
	}
	if s.DoorSink == nil {
		writeErr(w, 503, "the door receiver is not configured")
		return
	}
	if ask.Definition.Validate() != nil || !supportedTrigger(ask.Definition) || len(ask.Channels) == 0 || len(ask.Channels) > 100 || len(ask.Config) > 16<<10 {
		writeErr(w, 400, "approve a supported event definition and 1 to 100 channel ids")
		return
	}
	for _, channel := range ask.Channels {
		if !channelID.MatchString(channel) {
			writeErr(w, 400, "each channel must be an exact Slack channel id")
			return
		}
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !g.SignedIn || g.Request.Connector == nil || g.Connector == nil || g.Status == StatusFailed || g.Status == StatusDisconnected || g.Status == StatusNeedsReconnect || ask.Account != g.Connector.Account || ask.Definition.Toolkit != g.Request.Connector.Toolkit {
		writeErr(w, 409, "confirm the exact connected account before enabling its events")
		return
	}
	if _, err := s.expire(r.Context(), g); err != nil || s.overdue(g) || g.Status == StatusFailed {
		writeErr(w, 409, "the connection deadline passed")
		return
	}
	if err := s.checkConnector(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	p, err := s.eventProvider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if err := s.ensureProjectWebhook(r.Context(), g.Request.Connector.Provider, p); err != nil {
		s.fail(w, r, err)
		return
	}
	id, err := p.CreateTrigger(r.Context(), s.connectorIdentity(g), ask.Account, ask.Definition, ask.Config)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	key := store.Item{PK: triggerPK(g.Request.Connector.Provider, id), SK: "-", Ver: 1}
	var prior triggerRecord
	old, loadErr := getDoc(r.Context(), s.KV, key.PK, key.SK, &prior)
	cond := store.Cond{NotExists: true}
	if loadErr == nil {
		if prior.Tenant != g.Tenant || prior.Grant != g.ID || prior.Account != ask.Account {
			writeErr(w, 409, "that upstream event belongs to another approved connection")
			return
		}
		cond = store.Cond{Ver: &old.Ver}
		key.Ver = old.Ver + 1
	} else if !errors.Is(loadErr, errNotFound) {
		s.fail(w, r, loadErr)
		return
	}
	record := triggerRecord{Tenant: g.Tenant, Grant: g.ID, Provider: g.Request.Connector.Provider, Account: ask.Account, User: s.connectorIdentity(g).User, ID: id, Definition: ask.Definition, Config: ask.Config, Channels: ask.Channels, Approval: sr.approval, Active: true}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(g.Tenant), SK: "TRIGGER#" + record.Provider + "#" + record.ID}, map[string]string{"provider": record.Provider, "trigger": record.ID}, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	if err := putDoc(r.Context(), s.KV, key, record, cond); err != nil {
		s.fail(w, r, err)
		return
	}
	if !slices.Contains(g.TriggerIDs, id) {
		g.TriggerIDs = append(g.TriggerIDs, id)
	}
	g.HeldUntil = time.Time{}
	g.event(s.now(), "door events approved "+sr.approval.Who())
	if err := s.saveGrant(r.Context(), g); err != nil {
		record.Active = false
		_ = putDoc(r.Context(), s.KV, key, record, store.Cond{})
		_ = p.DeleteTrigger(r.Context(), id)
		s.fail(w, r, err)
		return
	}
	writeJSON(w, 201, map[string]any{"trigger": id, "account": ask.Account, "channels": ask.Channels})
}
func (s *Server) listTriggers(w http.ResponseWriter, r *http.Request) {
	if _, err := s.agentKey(r, r.PathValue("tenant")); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), r.PathValue("tenant"), r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	records := []triggerRecord{}
	if g.Request.Connector != nil {
		for _, id := range g.TriggerIDs {
			var record triggerRecord
			if _, err := getDoc(r.Context(), s.KV, triggerPK(g.Request.Connector.Provider, id), "-", &record); err != nil {
				s.fail(w, r, err)
				return
			}
			if record.Tenant != g.Tenant || record.Grant != g.ID {
				writeErr(w, 409, "the event binding does not match")
				return
			}
			records = append(records, record)
		}
	}
	writeJSON(w, 200, map[string]any{"triggers": records})
}
func (s *Server) stopTrigger(ctx context.Context, g *Grant, id string) error {
	if g.Request.Connector == nil {
		return connector.ErrDefinition
	}
	var record triggerRecord
	it, err := getDoc(ctx, s.KV, triggerPK(g.Request.Connector.Provider, id), "-", &record)
	if errors.Is(err, errNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if record.Tenant != g.Tenant || record.Grant != g.ID || record.Account != g.Connector.Account {
		return connector.ErrAccount
	}
	// Disable local delivery before contacting the provider. A revoked route
	// cannot be revived by a delayed webhook or an upstream deletion failure.
	record.Active = false
	it.Ver++
	if err := putDoc(ctx, s.KV, it, record, store.Cond{}); err != nil {
		return err
	}
	p, err := s.eventProvider(g)
	if err != nil {
		return err
	}
	return p.DeleteTrigger(ctx, id)
}
func (s *Server) stopTriggers(ctx context.Context, g *Grant) error {
	for _, id := range g.TriggerIDs {
		if err := s.stopTrigger(ctx, g, id); err != nil {
			return err
		}
	}
	return nil
}
func (s *Server) deleteTrigger(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var ask struct {
		Account string `json:"account"`
	}
	if err := decodeStrict(sr.body, &ask); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	id := r.PathValue("trigger")
	if g.Connector == nil || ask.Account != g.Connector.Account || !slices.Contains(g.TriggerIDs, id) {
		writeErr(w, 404, "no such approved trigger on this connection")
		return
	}
	if err := s.stopTrigger(r.Context(), g, id); err != nil {
		s.fail(w, r, err)
		return
	}
	g.TriggerIDs = slices.DeleteFunc(g.TriggerIDs, func(t string) bool { return t == id })
	if len(g.TriggerIDs) == 0 && (g.Status == StatusListed || g.Status == StatusPending && len(g.Tools) == 0) {
		g.HeldUntil = s.now().Add(holdTTL).UTC()
	}
	g.event(s.now(), "door events disabled "+sr.approval.Who())
	if err := s.saveGrant(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, 200, map[string]any{"trigger": id, "disabled": true})
}

func (s *Server) accountTriggers(ctx context.Context, tenant string) ([]triggerRecord, error) {
	indexes, err := s.KV.Query(ctx, tenantPK(tenant), "TRIGGER#")
	if err != nil {
		return nil, err
	}
	records := []triggerRecord{}
	for _, index := range indexes {
		var binding struct {
			Provider string `json:"provider"`
			Trigger  string `json:"trigger"`
		}
		if json.Unmarshal(index.Doc, &binding) != nil || binding.Provider == "" || binding.Trigger == "" {
			return nil, connector.ErrDefinition
		}
		var record triggerRecord
		if _, err := getDoc(ctx, s.KV, triggerPK(binding.Provider, binding.Trigger), "-", &record); errors.Is(err, errNotFound) {
			continue
		} else if err != nil {
			return nil, err
		}
		if record.Tenant != tenant {
			return nil, connector.ErrAccount
		}
		records = append(records, record)
	}
	return records, nil
}
