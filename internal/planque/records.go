package planque

import (
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

// Approval is the owner's signature on a grant, kept as evidence.
type Approval struct {
	DeviceKey   []byte    `json:"device_key"`
	Fingerprint string    `json:"fingerprint"`
	Message     string    `json:"message"` // the exact text signed
	Body        string    `json:"body"`    // the request body its hash covers
	Signature   string    `json:"signature"`
	SignedAt    string    `json:"signed_at"`
	ApprovedAt  time.Time `json:"approved_at"`
}

// GrantedTool is a vendor tool as planque relays it.
type GrantedTool struct {
	ToolRequest
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

// Event is one line of a grant's record.
type Event struct {
	At   time.Time `json:"at"`
	What string    `json:"what"`
}

// Grant statuses.
const (
	StatusPending      = "pending"      // waiting for the owner's signature
	StatusApproved     = "approved"     // signed; waiting for the vendor sign-in
	StatusSignedIn     = "signed_in"    // token held, unused; waiting for the approving phone to confirm
	StatusConnected    = "connected"    // token held, tools registered in KIFF
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
	Approval *Approval      `json:"approval,omitempty"`
	Tools    []GrantedTool  `json:"tools,omitempty"`
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
func pkClient(issuer, redirect string) string {
	return "CLIENT#" + issuer + "#" + redirect
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
	}
	return err
}

func (g *Grant) event(now time.Time, what string) {
	g.Events = append(g.Events, Event{At: now.UTC(), What: what})
}
