package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/store"
)

// Enrolment is a KIFF account set up to connect tools: its connector key.
// Items written before kiff/planque#12 also hold a phone's key
// (device_key); it is ignored.
type Enrolment struct {
	Tenant     string    `json:"tenant"`
	Connector  []byte    `json:"connector"` // the KIFF connector key, sealed
	EnrolledAt time.Time `json:"enrolled_at"`
}

// ToolRequest is one vendor tool an agent asks for, possibly narrowed.
//
// A vendor may offer one broad tool (Stripe's stripe_api_write writes to
// any endpoint). A grant narrows it: Fixed are vendor arguments planque
// always sends and the agent cannot set ("stripe_api_operation_id":
// "PostRefunds"), and Arguments are the only arguments the agent may
// set, each placed at a path in the vendor's arguments. Arguments sit at
// the top level, where KIFF's Card can govern them.
type ToolRequest struct {
	Tool string `json:"tool"`           // the vendor's tool name
	Name string `json:"name,omitempty"` // its name in KIFF and at the relay: the Card's action
	// Title is what the tool does, in the owner's words, as the asking agent
	// puts it ("Refund a payment"): the approval, KIFF's Card and the record
	// name the tool by it, never by its name (kiff/tafs#138, items 52, 53).
	Title string `json:"title,omitempty"`
	// SetupReads lets the asking agent itself read this tool, through
	// planque, once the owner approved it: for a vendor's tool that lists an
	// operation's fields, so Setup names them from the vendor, never by
	// guess (kiff/tafs#174). Only a read-only tool; the owner's text says
	// so, and every read is in the grant's record.
	SetupReads   bool                       `json:"setup_reads,omitempty"`
	AmountArg    string                     `json:"amount_argument,omitempty"`
	AmountPath   string                     `json:"amount_path,omitempty"`
	AmountFormat string                     `json:"amount_format,omitempty"`
	AmountUnit   string                     `json:"amount_unit,omitempty"`
	AmountScale  *int                       `json:"amount_scale,omitempty"`
	ReadOnly     bool                       `json:"read_only,omitempty"`
	Fixed        map[string]json.RawMessage `json:"fixed_arguments,omitempty"` // path → value
	Arguments    map[string]Argument        `json:"arguments,omitempty"`       // name → where it goes
	// ResourceArg is the string argument saying what a call acts on (a
	// payment): KIFF's Card can scope on it, and the owner's record names
	// it (kiff/tafs#167).
	ResourceArg  string `json:"resource_argument,omitempty"`
	ResourcePath string `json:"resource_path,omitempty"`
	// UnitLookup is the read of this grant that says which currency a call
	// is in, so KIFF counts only calls in the Card's.
	UnitLookup *kiff.UnitLookup `json:"unit_lookup,omitempty"`
	// AmountLookup is the read of this grant that says a call's amount from
	// its resource (a dispute's), for a write the agent gives no amount to:
	// KIFF's Card counts the amount read.
	AmountLookup *kiff.UnitLookup `json:"amount_lookup,omitempty"`
}

// Argument is one argument of a narrowed tool.
type Argument struct {
	Path        string `json:"path"` // dotted path in the vendor's arguments
	Type        string `json:"type"` // string, integer, number, boolean or array
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	// Schema is an array argument's reviewed shape (an invoice's lines): its
	// items, each an object whose fields are all named. A call is refused if
	// it sends a field the schema doesn't name, so a field left out of it (a
	// line's discount) can't be set. Only an array argument has one.
	Schema json.RawMessage `json:"schema,omitempty"`
	// OneOf is a string argument whose value must be one the account's own
	// read gives (an invoice's IBAN among the business's bank accounts):
	// before each call planque makes that read, with the grant's key, and
	// refuses a value it doesn't list. API-key grants only.
	OneOf *OneOf `json:"one_of,omitempty"`
}

// OneOf names a read-only operation of the same grant, called with no
// arguments, and the dotted path of the values it allows in its answer
// (through any lists on the way: "organization.bank_accounts.iban").
type OneOf struct {
	Tool string `json:"tool"`
	Path string `json:"path"`
}

// Request is exactly what is asked, and what the owner signs.
type Request struct {
	Version  string `json:"v"`
	Grant    string `json:"grant"`
	Tenant   string `json:"tenant"`
	Server   string `json:"server"`
	SignInAt string `json:"sign_in_at"` // the vendor's authorization server
	Agent    string `json:"agent"`
	// AgentName is the taf's name for people, when the request gave one.
	// Left out when empty, so a grant without it keeps its digest.
	AgentName string        `json:"agent_name,omitempty"`
	Tools     []ToolRequest `json:"tools"`
	AskedBy   string        `json:"asked_by"` // the KIFF agent or key that asked
	CreatedAt string        `json:"created_at"`
	// Kind is "api_key" for a vendor reached by its REST API with a key the
	// owner gives (kiff/planque#64); empty for an MCP server with a sign-in.
	// Both, like API and Reuse, are left out when empty, so an MCP grant's
	// digest is what it always was.
	Kind string `json:"kind,omitempty"`
	// API is what an API-key grant calls: the description it was read from,
	// that description's SHA-256 (the snapshot every call uses), the server
	// and where the key goes. The owner signs all of it.
	API *APIRequest `json:"api,omitempty"`
	// Reuse names a connected grant whose sign-in (or key) this one uses:
	// approving it needs no second sign-in (kiff/planque#36).
	Reuse string `json:"reuse,omitempty"`
	// Application is the exact registered app configuration approved for
	// this grant. It contains public facts and a client id, never a secret.
	Application *ApplicationApproval `json:"application,omitempty"`
	// Profile is the exact public profile the device approves. Keeping the
	// original JSON preserves the same canonical digest the phone checks.
	Profile      json.RawMessage   `json:"profile,omitempty"`
	CatalogEntry string            `json:"catalog_entry,omitempty"`
	Connector    *ConnectorRequest `json:"connector,omitempty"`
	// Browser is what a browser grant (kiff/planque#86) is limited to: its
	// sites and the currency it pays in. Left out for any other kind.
	Browser *BrowserRequest `json:"browser,omitempty"`
}

// APIRequest is the API an API-key grant calls.
type APIRequest struct {
	URL    string       `json:"url"`
	Digest string       `json:"digest"`
	Base   string       `json:"base"`
	Auth   openapi.Auth `json:"auth"`
	// Operations is the selected subset. Its names are signed together
	// with the full document digest, so selection cannot change silently.
	Operations []string         `json:"operations,omitempty"`
	SourceKind string           `json:"source_kind,omitempty"`
	VersionPin *rest.VersionPin `json:"version_header,omitempty"`
}

// Approval is the owner's signature on a grant, kept as evidence. The device
// is named by KIFF's id and person (kiff/planque#10). Approvals made before
// kiff/planque#12 by a phone enrolled with its key name that key, and are
// kept as they were recorded.
type Approval struct {
	DeviceKey   []byte    `json:"device_key,omitempty"`
	DeviceID    string    `json:"device_id,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	DeviceName  string    `json:"device_name,omitempty"` // as KIFF names it, for the owner to read
	Fingerprint string    `json:"fingerprint"`           // the key's fingerprint, or KIFF's device id
	Message     string    `json:"message"`               // the exact text signed
	Body        string    `json:"body"`                  // the request body its hash covers
	Signature   string    `json:"signature"`
	SignedAt    string    `json:"signed_at"`
	ApprovedAt  time.Time `json:"approved_at"`
	// Key is which key signed, as KIFF says it (kiff/tafs#494): "device" or
	// "critical" (the approval key, behind the fingerprint); empty before
	// KIFF said. Window is the approval window a device-key signature was
	// made in, so an audit tells a tap from a fingerprint.
	Key    string `json:"key,omitempty"`
	Window string `json:"window,omitempty"`
	// Lang and Text are the words the owner signed with digest_v 2, and
	// their language: the receipt of what they read. Empty for a digest
	// that covers the request alone.
	Lang string `json:"lang,omitempty"`
	Text string `json:"text,omitempty"`
}

// Who names who signed, for the grant's events: the person and their
// device KIFF knows, or the phone enrolled with its key.
//
// The words are for the owner, who reads a grant's last event in their app
// (#15): the device's name as KIFF gives it, never an id or a fingerprint.
// Those stay in the approval's own fields, as the evidence.
func (a Approval) Who() string {
	switch {
	case a.DeviceName != "":
		return "on " + a.DeviceName
	case a.DeviceID != "":
		return "on a device signed in to KIFF"
	}
	return "on the enrolled phone"
}

// SameSigner says whether b was signed on the same device as a.
func (a Approval) SameSigner(b Approval) bool {
	if a.DeviceID != "" || b.DeviceID != "" {
		return a.DeviceID == b.DeviceID
	}
	return len(a.DeviceKey) > 0 && bytes.Equal(a.DeviceKey, b.DeviceKey)
}

// GrantedTool is a vendor tool as planque relays it.
type GrantedTool struct {
	ToolRequest
	// Title is the vendor's name for the tool, for people (its MCP title).
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
	Connector   *connector.Tool `json:"connector,omitempty"`
	// VendorSchemaSHA256 is the canonical hash of the tool's input schema as
	// the vendor listed it when the owner approved it (kiff/planque#71): a
	// different one later pauses the grant.
	VendorSchemaSHA256 string `json:"vendor_schema_sha256,omitempty"`
	// FieldsChecked says planque read this narrowed write's operation and
	// found every field it names there: set at its first call once checked,
	// so a tool defined before the check (kiff/planque#61) is checked too.
	FieldsChecked bool `json:"fields_checked,omitempty"`
}

// CatalogTool is one tool the vendor's server lists: its own words and
// arguments, never a value from the account.
type CatalogTool struct {
	Name        string          `json:"name"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
	ReadOnly    bool            `json:"read_only,omitempty"` // the server's readOnlyHint
	// Method and Path are an API-key grant's operation: what a call sends.
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

// Event is one line of a grant's record.
type Event struct {
	At   time.Time `json:"at"`
	What string    `json:"what"`
	// Recorded is the line as it was stored, when the view shows it in the
	// owner's words instead (an event written before #16 named who signed
	// by ids). The stored record itself is never rewritten.
	Recorded string `json:"recorded,omitempty"`
}

// Grant statuses.
const (
	StatusPending   = "pending"   // waiting for the owner's signature
	StatusApproved  = "approved"  // signed; waiting for the vendor sign-in
	StatusSignedIn  = "signed_in" // token held, unused; waiting for the approving phone to confirm
	StatusConnected = "connected" // token held, tools registered in KIFF
	// StatusListed: a grant asked for with no tools (connect first): the
	// sign-in is confirmed and held, and the vendor's tools are listed in
	// its catalog. Nothing is in KIFF until the asking agent sets the
	// tools and the owner approves them.
	StatusListed       = "listed"
	StatusFailed       = "failed"       // see Reason
	StatusDisconnected = "disconnected" // token deleted, tools removed
	// StatusNeedsReconnect: a connected grant whose sign-in the vendor
	// ended (it refused the refresh). Its tools stay in KIFF and its calls
	// are refused, saying so, until the owner signs in again (Reconnect).
	StatusNeedsReconnect = "needs_reconnect"
)

// Grant is one connection of a vendor account to one taf.
type Grant struct {
	ID      string  `json:"id"`
	Tenant  string  `json:"tenant"`
	Status  string  `json:"status"`
	Request Request `json:"request"`
	Text    string  `json:"text"`
	// Texts are Text in every language planque speaks, by language, written
	// with it: what the owner's phone shows, and signs with digest_v 2. A
	// grant made before has none, and is shown in English.
	Texts  map[string]string `json:"texts,omitempty"`
	Digest string            `json:"digest"`
	// Proposer is the agent that proposed the tools of a grant an account
	// key asked (its taf, or the account's Setup); empty otherwise.
	Proposer string         `json:"proposer,omitempty"`
	Vendor   oauth.Metadata `json:"vendor"`
	ClientID string         `json:"client_id,omitempty"`
	Account  string         `json:"account,omitempty"` // the vendor account the token is bound to
	// AccountLive is the account's mode, where the vendor's own list of
	// the sign-in's accounts said it (Stripe's livemode); nil when unsaid.
	AccountLive *bool `json:"account_live,omitempty"`
	// VendorName is the vendor's name for people, as its MCP server names
	// itself (serverInfo: title, else name), once listed or connected.
	VendorName string `json:"vendor_name,omitempty"`
	// Lang is the language the owner's app asked in (its Accept-Language)
	// when it approved or opened a sign-in: the language of the pages planque
	// then shows in the browser.
	Lang string `json:"lang,omitempty"`
	// Entry is the catalog entry for the grant's app, which planque recorded
	// on a grant made before its request named it (entry.go). It's kept
	// outside the request on purpose: the request is what the owner signed,
	// and stays byte for byte as signed. A grant made since names its entry
	// in Request.CatalogEntry, signed with the rest.
	Entry    string        `json:"entry,omitempty"`
	Approval *Approval     `json:"approval,omitempty"`
	Tools    []GrantedTool `json:"tools,omitempty"`
	// SignedIn: the owner's sign-in was confirmed on the phone and its
	// token is held, so approving new tools needs no second sign-in.
	SignedIn bool `json:"signed_in,omitempty"`
	// OnDemand: the app's sign-in was set up at the provider on this
	// owner's pick (kiff/planque#87). Coach drafts such an app's writes and
	// money to ask the owner by default.
	OnDemand bool `json:"on_demand,omitempty"`
	// Reconnect: the owner approved a new sign-in for a grant the vendor
	// signed out. Confirmed, it is connected again with the same tools, in
	// KIFF as they were; until then its calls are refused, saying why.
	Reconnect bool `json:"reconnect,omitempty"`
	// SignedOutAt is when the vendor last ended its sign-in: KIFF is told
	// of each sign-out once, by it.
	SignedOutAt time.Time `json:"signed_out_at,omitzero"`
	// HeldUntil is when a held sign-in that isn't connected (a listed
	// grant, or its tools waiting for approval) is revoked if nothing
	// connected it: planque keeps no live vendor token the taf doesn't use.
	HeldUntil time.Time `json:"held_until,omitzero"`
	// Catalog is what the vendor's MCP server listed, once signed in, for
	// a grant asked for with no tools: what the asking agent chooses from.
	Catalog []CatalogTool `json:"catalog,omitempty"`
	// Transport is the SHA-256 of the credential KIFF's gateway presents
	// to the relay. The credential itself is only at KIFF.
	Transport string `json:"transport,omitempty"`
	// Transports are the SHA-256s of the credentials of tools added to the
	// grant once connected: each addition registers its tools in KIFF with
	// a credential of its own, since KIFF keeps the first one's (kiff/tafs#138).
	Transports []string `json:"transports,omitempty"`
	// Connected are the tools a connected grant had when the asking agent
	// added more: set while the addition waits for the owner's approval,
	// during which those tools keep working. Approved, the grant is
	// connected with both; not approved within a day, with these alone.
	Connected        []ToolRequest   `json:"connected,omitempty"`
	ConnectedProfile json.RawMessage `json:"connected_profile,omitempty"`
	// ProfileRevision is KIFF's accepted revision, never the proposing key's
	// assertion that an approval applied.
	ProfileRevision *kiff.ProfileRevision `json:"profile_revision,omitempty"`
	Connector       *ConnectorState       `json:"connector,omitempty"`
	// ConfirmHash is the SHA-256 of the code the vendor sign-in handed to
	// the browser that completed it; the approving phone presents it to
	// confirm, by ConfirmBy.
	TriggerIDs  []string  `json:"trigger_ids,omitempty"`
	ConfirmHash string    `json:"confirm_hash,omitempty"`
	ConfirmBy   time.Time `json:"confirm_by,omitzero"`
	Audience    string    `json:"audience"`
	Reason      string    `json:"reason,omitempty"`
	// Paused: a tool the owner approved changed at the vendor; nothing is
	// sent on the grant until the owner approves the tools again.
	Paused *Pause `json:"paused,omitempty"`
	// DefinitionsReadAt is when a direct MCP server's tools were last read
	// to notice a change (checkDefinitions).
	DefinitionsReadAt time.Time `json:"definitions_read_at,omitzero"`
	Events            []Event   `json:"events"`
	// TokenOf is the grant whose sign-in (or key) this one uses, for a grant
	// asked for with reuse: its token is that grant's, never a copy, so a
	// refresh happens once for both. Cleared when that sign-in is gone.
	TokenOf string `json:"token_of,omitempty"`
	// Telegram is the bot a Telegram grant's token is, and its webhook.
	Telegram *TelegramBot `json:"telegram,omitempty"`
	// Declared: a key grant's tools are in KIFF before the key is given, so
	// a Card can be drafted over them and signed with the key, but nothing
	// is usable: the relay refuses every call until the grant is connected.
	// Its Tools and Transport are what was registered.
	Declared bool `json:"declared,omitempty"`
	// Pending are a connected grant's added tools, declared in KIFF while the
	// addition waits for the owner's approval, so the taf's Card can be
	// drafted over them and signed with that approval (tools_and_cards). The
	// relay lists them and refuses every call of them until approved; with
	// the approval they join Tools, their registration and PendingTransport
	// unchanged. The addition ending unapproved takes them out of KIFF.
	Pending          []GrantedTool `json:"pending,omitempty"`
	PendingTransport string        `json:"pending_transport,omitempty"`
	// Renewal is a new sign-in the owner approved for a connected grant
	// (gap 13 of kiff/tafs#476): its tools keep working with the sign-in in
	// use until the approving device confirms the new one, which then
	// replaces it. Not finished by its By, it's dropped.
	Renewal *Renewal `json:"renewal,omitempty"`
	ver     int64
	// dueAt is the deadline g's index entry was written for, when g was
	// loaded or last saved; zero when unknown or none.
	dueAt time.Time
	// keptFor is how many other live grants kept this one's sign-in when it
	// was disconnected; not stored.
	keptFor int
	// revokedShared: this grant was the last on a sign-in it reused, and
	// its disconnect or expiry revoked that sign-in; not stored.
	revokedShared bool
}

// Renewal is a connected grant's new sign-in under way: who approved it, the
// client it signs in with, the account it came back in, and by when it must
// come back, then be confirmed (By). Its token, once the vendor gave it, is
// held apart from the one in use (skRenewal) until it's confirmed.
type Renewal struct {
	Approval Approval `json:"approval"`
	ClientID string   `json:"client_id,omitempty"`
	Account  string   `json:"account,omitempty"`
	// ConfirmHash is the SHA-256 of the code the vendor's sign-in handed to
	// the browser that finished it: set once it came back.
	ConfirmHash string    `json:"confirm_hash,omitempty"`
	By          time.Time `json:"by"`
}

// pendingSignIn is one vendor sign-in in progress, keyed by its state.
type pendingSignIn struct {
	Tenant   string    `json:"tenant"`
	Grant    string    `json:"grant"`
	Verifier []byte    `json:"verifier"` // sealed
	Expires  time.Time `json:"expires"`
}

// tokenRecord is a vendor token, sealed. Expiry is kept in the clear so
// a call knows when to refresh without opening it.
type tokenRecord struct {
	Sealed []byte    `json:"sealed"`
	Expiry time.Time `json:"expiry,omitzero"`
}

// Operation states: each permit's operation runs at most once.
const (
	OpStarted = "started" // claimed; the lease holder may send
	OpDone    = "done"    // the vendor answered; its answer is kept
	OpNotSent = "not_sent"
	OpUnknown = "unknown" // sent, the answer was lost: never resent
)

// operation is one permit's execution, keyed by its op.
type operation struct {
	Tool            string                `json:"tool"`
	ArgsHash        string                `json:"args_sha256"`
	State           string                `json:"state"`
	StartedAt       time.Time             `json:"started_at"`
	Decision        string                `json:"decision,omitempty"`
	Result          json.RawMessage       `json:"result,omitempty"`
	Profile         json.RawMessage       `json:"profile,omitempty"`
	ProfileRevision *kiff.ProfileRevision `json:"profile_revision,omitempty"`
}

func tenantPK(t string) string { return "TENANT#" + t }

const skEnrolment = "ENROLMENT"

func skGrant(id string) string { return "GRANT#" + id }
func skToken(id string) string { return "TOKEN#" + id }

// skRenewal holds the token of a connected grant's new sign-in until it's
// confirmed (Renewal).
func skRenewal(id string) string { return "RENEWAL#" + id }
func skSnapshot(id string, n int) string {
	return fmt.Sprintf("SNAPSHOT#%s#%03d", id, n)
}
func skOp(id, op string) string    { return "OP#" + id + "#" + op }
func pkSignIn(state string) string { return "SIGNIN#" + state }
func pkClient(issuer, redirect, name string) string {
	return "CLIENT#" + issuer + "#" + redirect + "#" + name
}
func pkSigned(hash string) string { return "SIGNED#" + hash }

var errNotFound = errors.New("not found")

func getDoc(ctx context.Context, kv store.KV, pk, sk string, v any) (store.Item, error) {
	it, ok, err := kv.Get(ctx, pk, sk)
	if err != nil {
		return it, err
	}
	if !ok {
		return it, errNotFound
	}
	return it, json.Unmarshal(it.Doc, v)
}

func putDoc(ctx context.Context, kv store.KV, it store.Item, v any, c store.Cond) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	it.Doc = b
	return kv.Put(ctx, it, c)
}

func (s *Server) loadGrant(ctx context.Context, tenant, id string) (*Grant, error) {
	var g Grant
	it, err := getDoc(ctx, s.KV, tenantPK(tenant), skGrant(id), &g)
	if err != nil {
		return nil, err
	}
	g.ver = it.Ver
	g.dueAt = g.deadline()
	return &g, nil
}

// saveGrant writes g if nobody changed it since it was loaded. A deadline g
// had when loaded and no longer has loses its index entry once g is saved
// (a request's 7 days, once it's signed): an entry left behind would be
// harmless, since the sweep checks the grant itself, but it would wait there
// for days.
func (s *Server) saveGrant(ctx context.Context, g *Grant) error {
	c := store.Cond{Ver: &g.ver}
	if g.ver == 0 {
		c = store.Cond{NotExists: true}
	}
	if err := s.markDue(ctx, g); err != nil {
		return err
	}
	err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skGrant(g.ID), Ver: g.ver + 1}, g, c)
	if err != nil {
		return err
	}
	g.ver++
	now := g.deadline()
	if was := g.dueAt; !was.IsZero() && (now.IsZero() || dueSK(was, g.Tenant, g.ID) != dueSK(now, g.Tenant, g.ID)) {
		_, _, _ = s.KV.Delete(ctx, pkDue, dueSK(was, g.Tenant, g.ID))
	}
	g.dueAt = now
	return nil
}

func (g *Grant) event(now time.Time, what string) {
	g.Events = append(g.Events, Event{At: now.UTC(), What: what})
}
