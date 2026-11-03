package planque

import (
	"bytes"
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
	// Removing: a sweep or a failed approval has claimed this subscription to
	// delete it upstream. It routes nothing, and no approval reuses it until
	// the deletion is done.
	Removing bool `json:"removing,omitempty"`
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

// agentsGrant is the grant, for a KIFF key of its own agent or of the agent
// that asked for it: another agent of the account never reads its events.
func (s *Server) agentsGrant(w http.ResponseWriter, r *http.Request) (*Grant, bool) {
	me, err := s.agentKey(r, r.PathValue("tenant"))
	if err != nil {
		s.fail(w, r, err)
		return nil, false
	}
	g, err := s.loadGrant(r.Context(), r.PathValue("tenant"), r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return nil, false
	}
	if err != nil {
		s.fail(w, r, err)
		return nil, false
	}
	if me.AgentID == "" || me.AgentID != g.Request.Agent && me.AgentID != g.Request.AskedBy {
		writeErr(w, http.StatusForbidden, "only this connection's own agent may read its events")
		return nil, false
	}
	return g, true
}
func (s *Server) triggerType(w http.ResponseWriter, r *http.Request) {
	g, ok := s.agentsGrant(w, r)
	if !ok {
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

// createTrigger subscribes channels of the grant's connected account to a
// supported event, as the device approved: {account, slug, channels}. planque
// resolves the definition for slug and pins its version, and builds each
// upstream config from the approved channels only, so no subscription can
// cover more than the owner approved: a channel-message door gets one
// subscription per channel (a null channel_id would mean every channel), a
// direct-message door one subscription with no channel filter, its DM ids
// gating delivery. It answers the subscriptions' ids and the pinned definition.
func (s *Server) createTrigger(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	var ask struct {
		Account  string   `json:"account"`
		Slug     string   `json:"slug"`
		Channels []string `json:"channels"`
	}
	if len(sr.body) > 64<<10 {
		writeErr(w, 400, "the event approval exceeds 64 KB")
		return
	}
	if err := decodeStrict(sr.body, &ask); err != nil {
		s.fail(w, r, err)
		return
	}
	if s.DoorSink == nil {
		writeErr(w, 503, "the door receiver is not configured")
		return
	}
	subscriptions, err := doorSubscriptions(ask.Slug, ask.Channels)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !g.SignedIn || g.Request.Connector == nil || g.Connector == nil || g.Status == StatusFailed || g.Status == StatusDisconnected || g.Status == StatusNeedsReconnect || ask.Account != g.Connector.Account {
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
	definition, err := p.TriggerType(r.Context(), g.Request.Connector.Toolkit, ask.Slug, "")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if definition.Validate() != nil || !supportedTrigger(definition) || definition.Slug != ask.Slug || definition.Toolkit != g.Request.Connector.Toolkit {
		writeErr(w, 400, "this event type is not supported for a door")
		return
	}
	if err := s.ensureProjectWebhook(r.Context(), g.Request.Connector.Provider, p); err != nil {
		s.fail(w, r, err)
		return
	}
	// Any failure takes back every upstream subscription this request made:
	// nothing listens upstream that planque keeps no active record of.
	var made []madeTrigger
	undo := func(err error) {
		s.takeBack(r.Context(), p, made)
		s.fail(w, r, err)
	}
	ids := []string{}
	for _, sub := range subscriptions {
		record, ver, kept, err := s.subscribe(r.Context(), g, p, definition, sub, sr.approval)
		// A subscription that was already on before this request stays on.
		if record.ID != "" && !kept {
			made = append(made, madeTrigger{record, ver})
		}
		if err != nil {
			undo(err)
			return
		}
		ids = append(ids, record.ID)
	}
	for _, id := range ids {
		if !slices.Contains(g.TriggerIDs, id) {
			g.TriggerIDs = append(g.TriggerIDs, id)
		}
	}
	g.HeldUntil = time.Time{}
	g.event(s.now(), "door events approved "+sr.approval.Who())
	if err := s.saveGrant(r.Context(), g); err != nil {
		undo(err)
		return
	}
	writeJSON(w, 201, map[string]any{"id": ids[0], "trigger": ids[0], "triggers": ids, "account": ask.Account, "channels": ask.Channels, "definition": definition})
}

// triggerCleanupPK holds upstream subscriptions a failed approval made and
// could not delete at once: the sweep deletes them later. Each names only
// the provider, the subscription, whose approval made it, and Fence: the
// version of its record when it was noted (0 for none). A record written
// after that is a newer approval, which the note never undoes.
const triggerCleanupPK = "CONNECTOR#TRIGGER_CLEANUP"
const takeBackTimeout = 20 * time.Second

type triggerCleanup struct {
	Provider string `json:"provider"`
	Trigger  string `json:"trigger"`
	Tenant   string `json:"tenant"`
	Grant    string `json:"grant"`
	Fence    int64  `json:"fence"`
}

// madeTrigger is a subscription a request turned on, and the version of
// its record that request is answerable for (0 for none).
type madeTrigger struct {
	record triggerRecord
	ver    int64
}

// takeBack undoes the subscriptions a failed approval made, even when its
// request is gone: each is noted for the sweep, then removed as the sweep
// would remove it, and the note goes only once it is.
func (s *Server) takeBack(ctx context.Context, p connector.EventProvider, made []madeTrigger) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), takeBackTimeout)
	defer cancel()
	for _, m := range made {
		note := triggerCleanup{Provider: m.record.Provider, Trigger: m.record.ID, Tenant: m.record.Tenant, Grant: m.record.Grant, Fence: m.ver}
		if err := putDoc(ctx, s.KV, store.Item{PK: triggerCleanupPK, SK: note.Provider + "#" + note.Trigger}, note, store.Cond{}); err != nil {
			s.log().WarnContext(ctx, "an undone subscription was not noted for the sweep", "trigger", note.Trigger, "err", err.Error())
		}
		if _, err := s.removeTrigger(ctx, p, note); err != nil {
			s.log().WarnContext(ctx, "an undone subscription is left for the sweep", "trigger", note.Trigger, "err", err.Error())
		}
	}
}

// removeTrigger deletes a noted subscription upstream, once it holds the
// claim on its record: a write conditioned on the version the note was taken
// at, with a lease. A newer record is a newer approval: the note is dropped
// and nothing is deleted. The claim is exclusive through the provider call:
// another worker takes it over only once its lease has passed, and the holder
// gives up its own call well before that, so two deletions never overlap.
// While it is held, the record routes nothing and an approval that reuses
// the id is refused. It is settled only by a write conditioned on the
// version it claimed: a worker that lost its claim settles nothing. done
// says the note is settled.
func (s *Server) removeTrigger(ctx context.Context, p connector.EventProvider, note triggerCleanup) (done bool, err error) {
	key := store.Item{PK: triggerPK(note.Provider, note.Trigger), SK: "-"}
	noteKey := [2]string{triggerCleanupPK, note.Provider + "#" + note.Trigger}
	now := s.now()
	key.Lease = now.Add(removeLease).UnixMilli()
	var record triggerRecord
	it, err := getDoc(ctx, s.KV, key.PK, key.SK, &record)
	created := false
	switch {
	case errors.Is(err, errNotFound):
		// No record (never written, or erased): claim with one that routes nothing.
		record = triggerRecord{Tenant: note.Tenant, Grant: note.Grant, Provider: note.Provider, ID: note.Trigger, Removing: true}
		key.Ver = 1
		if err := putDoc(ctx, s.KV, key, record, store.Cond{NotExists: true}); err != nil {
			return false, ignoreCondition(err)
		}
		created = true
	case err != nil:
		return false, err
	case record.Removing:
		// Another worker's claim: taken over only once its lease has passed.
		if it.Lease >= now.UnixMilli() {
			return false, nil
		}
		key.Ver = it.Ver + 1
		if err := putDoc(ctx, s.KV, key, record, store.Cond{Ver: &it.Ver, LeaseFreeAt: now.UnixMilli()}); err != nil {
			return false, ignoreCondition(err)
		}
	case it.Ver > note.Fence:
		_, _, err := s.KV.Delete(ctx, noteKey[0], noteKey[1])
		return err == nil, err
	default:
		record.Active, record.Removing = false, true
		key.Ver = it.Ver + 1
		if err := putDoc(ctx, s.KV, key, record, store.Cond{Ver: &it.Ver}); err != nil {
			return false, ignoreCondition(err)
		}
	}
	call, cancel := context.WithTimeout(ctx, removeCall)
	err = p.DeleteTrigger(call, note.Trigger)
	cancel()
	if err != nil {
		return false, err
	}
	claimed := key.Ver
	record.Removing = false
	key.Ver, key.Lease = claimed+1, 0
	if err := putDoc(ctx, s.KV, key, record, store.Cond{Ver: &claimed}); err != nil {
		return false, ignoreCondition(err)
	}
	if created {
		if _, _, err := s.KV.Delete(ctx, key.PK, key.SK); err != nil {
			return false, err
		}
	}
	_, _, err = s.KV.Delete(ctx, noteKey[0], noteKey[1])
	return err == nil, err
}

// removeLease is how long a removal claim is exclusive; removeCall, how long
// its holder waits for the provider, far within it.
const removeLease, removeCall = 2 * time.Minute, 30 * time.Second

// ignoreCondition: a lost race is no error; the next sweep reads again.
func ignoreCondition(err error) error {
	if errors.Is(err, store.ErrCondition) {
		return nil
	}
	return err
}

// cleanTriggers removes the subscriptions takeBack could not.
func (s *Server) cleanTriggers(ctx context.Context) (int, error) {
	items, err := s.KV.Query(ctx, triggerCleanupPK, "")
	if err != nil {
		return 0, err
	}
	cleaned := 0
	var errs []error
	for _, it := range items {
		var note triggerCleanup
		if json.Unmarshal(it.Doc, &note) != nil {
			continue
		}
		p, ok := s.Connectors[note.Provider].(connector.EventProvider)
		if !ok {
			continue
		}
		done, err := s.removeTrigger(ctx, p, note)
		if err != nil {
			errs = append(errs, err)
		}
		if done {
			cleaned++
		}
	}
	return cleaned, errors.Join(errs...)
}

// supportedSlugs are the events a door may subscribe to, for now.
var supportedSlugs = []string{"SLACKBOT_CHANNEL_MESSAGE_RECEIVED", "SLACKBOT_DIRECT_MESSAGE_RECEIVED"}

var conversationID = regexp.MustCompile(`^[CG][A-Z0-9]{2,40}$`)
var directID = regexp.MustCompile(`^D[A-Z0-9]{2,40}$`)

// subscription is one upstream subscription a door asks for: its config and
// the channels its deliveries are routed to.
type subscription struct {
	config   json.RawMessage
	channels []string
}

// doorSubscriptions are the upstream subscriptions for an approved slug and
// channels. A channel message takes one channel per subscription, never a null
// channel_id, and never the bot's own posts; a direct message has no channel
// filter upstream, so one subscription covers its approved DM ids.
func doorSubscriptions(slug string, channels []string) ([]subscription, error) {
	if !slices.Contains(supportedSlugs, slug) {
		return nil, apiError{http.StatusBadRequest, "this event type is not supported for a door"}
	}
	if len(channels) == 0 || len(channels) > 100 {
		return nil, apiError{http.StatusBadRequest, "approve 1 to 100 channel ids"}
	}
	seen := map[string]bool{}
	for _, channel := range channels {
		if seen[channel] {
			return nil, apiError{http.StatusBadRequest, "each channel is approved once"}
		}
		seen[channel] = true
	}
	if slug == "SLACKBOT_DIRECT_MESSAGE_RECEIVED" {
		for _, channel := range channels {
			if !directID.MatchString(channel) {
				return nil, apiError{http.StatusBadRequest, "a direct-message door takes exact Slack DM ids"}
			}
		}
		return []subscription{{config: json.RawMessage(`{}`), channels: channels}}, nil
	}
	subs := make([]subscription, 0, len(channels))
	for _, channel := range channels {
		if !conversationID.MatchString(channel) {
			return nil, apiError{http.StatusBadRequest, "each channel must be an exact Slack channel id"}
		}
		config, err := json.Marshal(map[string]any{"channel_id": channel, "is_bot_message": false})
		if err != nil {
			return nil, err
		}
		subs = append(subs, subscription{config: config, channels: []string{channel}})
	}
	return subs, nil
}

// subscribe makes one upstream subscription and its record, then the tenant's
// index after the record: an index never names a record that isn't there. The
// provider may answer an id planque already has: for this connection it is
// reused, and kept says it was already on before this request: its record
// stays as it was, and a failure later in the request leaves it on. One being
// removed is refused until it is gone. The record it returns carries the
// upstream id whenever planque made one, even with an error, and ver the
// version of its record this request is answerable for, so the caller can
// take it back; one that belongs to another connection, or is being removed,
// is never returned to be taken back.
func (s *Server) subscribe(ctx context.Context, g *Grant, p connector.EventProvider, definition connector.TriggerType, sub subscription, approval Approval) (record triggerRecord, ver int64, kept bool, err error) {
	id, err := p.CreateTrigger(ctx, s.connectorIdentity(g), g.Connector.Account, definition, sub.config)
	if err != nil {
		return triggerRecord{}, 0, false, err
	}
	record = triggerRecord{Tenant: g.Tenant, Grant: g.ID, Provider: g.Request.Connector.Provider, Account: g.Connector.Account, User: s.connectorIdentity(g).User, ID: id, Definition: definition, Config: sub.config, Channels: sub.channels, Approval: approval}
	key := store.Item{PK: triggerPK(record.Provider, id), SK: "-", Ver: 1}
	var prior triggerRecord
	old, err := getDoc(ctx, s.KV, key.PK, key.SK, &prior)
	cond := store.Cond{NotExists: true}
	if err == nil {
		if prior.Tenant != g.Tenant || prior.Grant != g.ID || prior.Account != record.Account && !prior.Removing {
			return triggerRecord{}, 0, false, apiError{http.StatusConflict, "that upstream event belongs to another approved connection"}
		}
		if prior.Removing {
			return triggerRecord{}, 0, false, apiError{http.StatusConflict, "this channel's previous door is still being removed: try again in a few minutes"}
		}
		if prior.Active {
			// Already on: its record, and the approval that turned it on, stay
			// as they were, whatever happens to the rest of this request. The
			// provider reuses an instance only for the same config; any other
			// answer is not this channel's door, and the request fails.
			if !bytes.Equal(prior.Config, sub.config) {
				return prior, old.Ver, true, apiError{http.StatusBadGateway, "the provider answered another channel's door for this channel"}
			}
			return prior, old.Ver, true, nil
		}
		ver = old.Ver
		cond = store.Cond{Ver: &old.Ver}
		key.Ver = old.Ver + 1
	} else if !errors.Is(err, errNotFound) {
		return record, 0, false, err
	}
	record.Active = true
	if err := putDoc(ctx, s.KV, key, record, cond); err != nil {
		record.Active = false
		return record, ver, false, err
	}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: "TRIGGER#" + record.Provider + "#" + record.ID}, map[string]string{"provider": record.Provider, "trigger": record.ID}, store.Cond{}); err != nil {
		return record, key.Ver, false, err
	}
	return record, key.Ver, false, nil
}

func (s *Server) listTriggers(w http.ResponseWriter, r *http.Request) {
	g, ok := s.agentsGrant(w, r)
	if !ok {
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
	// The body is optional: with none, the subscription is the grant's own account's.
	if len(bytes.TrimSpace(sr.body)) > 0 {
		if err := decodeStrict(sr.body, &ask); err != nil {
			s.fail(w, r, err)
			return
		}
	}
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	id := r.PathValue("trigger")
	if g.Connector != nil && ask.Account == "" {
		ask.Account = g.Connector.Account
	}
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

// The connection's own action that lists its channels, pinned and read-only.
const listChannelsSlug, listChannelsVersion = "SLACKBOT_LIST_ALL_CHANNELS", "20261101_00"
const maxChannels = 200

// listChannels is the channels a door may be set on, for the owner choosing
// one on the device: ids, names and whether each is private, no content. It
// reads through the connection's own pinned list action, which must stay
// read-only, and pages up to 200 channels.
func (s *Server) listChannels(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	g, err := s.loadGrant(r.Context(), sr.tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if g.Request.Connector == nil || g.Request.Connector.Toolkit != "slackbot" {
		writeErr(w, 400, "this connection has no channels to choose from")
		return
	}
	if !g.SignedIn || g.Connector == nil || g.Status == StatusFailed || g.Status == StatusDisconnected || g.Status == StatusNeedsReconnect {
		writeErr(w, 409, "sign in to this connection before choosing its channels")
		return
	}
	if err := s.checkConnector(r.Context(), g); err != nil {
		s.fail(w, r, err)
		return
	}
	p, err := s.provider(g)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	tool, err := p.Tool(r.Context(), listChannelsSlug, listChannelsVersion)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if !tool.ReadOnly || tool.Toolkit != g.Request.Connector.Toolkit {
		writeErr(w, 503, "the channel list is not a read-only action of this connection")
		return
	}
	type channel struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Private bool   `json:"private"`
	}
	channels := []channel{}
	cursor := ""
	for len(channels) < maxChannels {
		args := map[string]any{"exclude_archived": true, "types": "public_channel,private_channel", "limit": maxChannels - len(channels)}
		if cursor != "" {
			args["cursor"] = cursor
		}
		raw, _ := json.Marshal(args)
		result, err := p.Execute(r.Context(), s.connectorIdentity(g), g.Connector.Account, tool, raw, func() bool { return r.Context().Err() == nil })
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if result.Refused || result.Unknown {
			writeErr(w, 502, "the connection did not list its channels")
			return
		}
		var page struct {
			Channels []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				IsPrivate bool   `json:"is_private"`
			} `json:"channels"`
			Meta struct {
				Next string `json:"next_cursor"`
			} `json:"response_metadata"`
		}
		if json.Unmarshal(result.Data, &page) != nil {
			writeErr(w, 502, "the connection listed its channels in an unknown form")
			return
		}
		for _, c := range page.Channels {
			if conversationID.MatchString(c.ID) && len(channels) < maxChannels {
				channels = append(channels, channel{ID: c.ID, Name: c.Name, Private: c.IsPrivate})
			}
		}
		if page.Meta.Next == "" || page.Meta.Next == cursor || len(page.Channels) == 0 {
			break
		}
		cursor = page.Meta.Next
	}
	writeJSON(w, 200, map[string]any{"channels": channels})
}
