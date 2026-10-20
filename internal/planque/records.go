package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/store"
)

// Enrolment is an owner's phone, enrolled for one KIFF account.
type Enrolment struct {
	Tenant     string    `json:"tenant"`
	DeviceKey  []byte    `json:"device_key"` // P-256 SPKI DER
	Connector  []byte    `json:"connector"`  // the KIFF connector key, sealed
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
	Tool        string                     `json:"tool"`           // the vendor's tool name
	Name        string                     `json:"name,omitempty"` // its name in KIFF and at the relay: the Card's action
	AmountArg   string                     `json:"amount_argument,omitempty"`
	AmountUnit  string                     `json:"amount_unit,omitempty"`
	AmountScale *int                       `json:"amount_scale,omitempty"`
	ReadOnly    bool                       `json:"read_only,omitempty"`
	Fixed       map[string]json.RawMessage `json:"fixed_arguments,omitempty"` // path → value
	Arguments   map[string]Argument        `json:"arguments,omitempty"`       // name → where it goes
}

// Argument is one argument of a narrowed tool.
type Argument struct {
	Path        string `json:"path"` // dotted path in the vendor's arguments
	Type        string `json:"type"` // string, integer, number or boolean
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

// Request is exactly what is asked, and what the owner signs.
type Request struct {
	Version   string        `json:"v"`
	Grant     string        `json:"grant"`
	Tenant    string        `json:"tenant"`
	Server    string        `json:"server"`
	SignInAt  string        `json:"sign_in_at"` // the vendor's authorization server
	Agent     string        `json:"agent"`
	Tools     []ToolRequest `json:"tools"`
	AskedBy   string        `json:"asked_by"` // the KIFF agent or key that asked
	CreatedAt string        `json:"created_at"`
}

// Approval is the owner's signature on a grant, kept as evidence. A device
// KIFF knows is named by KIFF's id and person (kiff/planque#10); a phone
// enrolled with its key, by that key.
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
}

// CatalogTool is one tool the vendor's server lists: its own words and
// arguments, never a value from the account.
type CatalogTool struct {
	Name        string          `json:"name"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
	ReadOnly    bool            `json:"read_only,omitempty"` // the server's readOnlyHint
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
)

// Grant is one connection of a vendor account to one taf.
type Grant struct {
	ID       string         `json:"id"`
	Tenant   string         `json:"tenant"`
	Status   string         `json:"status"`
	Request  Request        `json:"request"`
	Text     string         `json:"text"`
	Digest   string         `json:"digest"`
	Vendor   oauth.Metadata `json:"vendor"`
	ClientID string         `json:"client_id,omitempty"`
	Account  string         `json:"account,omitempty"` // the vendor account the token is bound to
	// VendorName is the vendor's name for people, as its MCP server names
	// itself (serverInfo: title, else name), once listed or connected.
	VendorName string        `json:"vendor_name,omitempty"`
	Approval   *Approval     `json:"approval,omitempty"`
	Tools      []GrantedTool `json:"tools,omitempty"`
	// SignedIn: the owner's sign-in was confirmed on the phone and its
	// token is held, so approving new tools needs no second sign-in.
	SignedIn bool `json:"signed_in,omitempty"`
	// HeldUntil is when a held sign-in that isn't connected (a listed
	// grant, or its tools waiting for approval) is revoked if nothing
	// connected it: planque keeps no live vendor token the taf doesn't use.
	HeldUntil time.Time `json:"held_until,omitzero"`
	// Catalog is what the vendor's MCP server listed, once signed in, for
	// a grant asked for with no tools: what the asking agent chooses from.
	Catalog []CatalogTool `json:"catalog,omitempty"`
	// Describable: the owner approved this grant on a text saying its
	// sign-in may be read, through the vendor's tools that change nothing
	// and take no input (kiff/tafs#89). A grant approved before that text,
	// or on a text naming tools, can't be described.
	Describable bool `json:"describable,omitempty"`
	// Transport is the SHA-256 of the credential KIFF's gateway presents
	// to the relay. The credential itself is only at KIFF.
	Transport string `json:"transport,omitempty"`
	// ConfirmHash is the SHA-256 of the code the vendor sign-in handed to
	// the browser that completed it; the approving phone presents it to
	// confirm, by ConfirmBy.
	ConfirmHash string    `json:"confirm_hash,omitempty"`
	ConfirmBy   time.Time `json:"confirm_by,omitzero"`
	Audience    string    `json:"audience"`
	Reason      string    `json:"reason,omitempty"`
	Events      []Event   `json:"events"`
	ver         int64
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

// Operation states (RFC 046 D5).
const (
	OpStarted = "started" // claimed; the lease holder may send
	OpDone    = "done"    // the vendor answered; its answer is kept
	OpNotSent = "not_sent"
	OpUnknown = "unknown" // sent, the answer was lost: never resent
)

// operation is one permit's execution, keyed by its op.
type operation struct {
	Tool      string          `json:"tool"`
	ArgsHash  string          `json:"args_sha256"`
	State     string          `json:"state"`
	StartedAt time.Time       `json:"started_at"`
	Decision  string          `json:"decision,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
}

func tenantPK(t string) string { return "TENANT#" + t }

const skEnrolment = "ENROLMENT"

func skGrant(id string) string     { return "GRANT#" + id }
func skToken(id string) string     { return "TOKEN#" + id }
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
	return &g, nil
}

// saveGrant writes g if nobody changed it since it was loaded.
func (s *Server) saveGrant(ctx context.Context, g *Grant) error {
	c := store.Cond{Ver: &g.ver}
	if g.ver == 0 {
		c = store.Cond{NotExists: true}
	}
	err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skGrant(g.ID), Ver: g.ver + 1}, g, c)
	if err == nil {
		g.ver++
		s.markDue(ctx, g)
	}
	return err
}

func (g *Grant) event(now time.Time, what string) {
	g.Events = append(g.Events, Event{At: now.UTC(), What: what})
}
