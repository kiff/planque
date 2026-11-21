package planque

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/lambda"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// The relay to the owner's own Chrome (tafs docs/chrome.md, "Connect"): the
// extension always opens the connection, a WebSocket through API Gateway,
// and nothing connects in. Each connect proves both ends: the extension
// signs planque's challenge with its paired key, and planque signs the
// extension's with its step key, which the extension's build pins. Only
// then does planque send steps, each signed, and take results, each signed
// by the extension.
//
// API Gateway can't post before $connect has answered, so the extension
// speaks first ({type: "open"}), and it can't close with a code, so planque
// says {type: "close", code} and then drops the connection.

// The signatures' labels, one per direction and purpose.
const (
	labelHello   = "tafs-chrome/hello/v1"
	labelWelcome = "tafs-chrome/welcome/v1"
	labelStep    = "tafs-chrome/step/v1"
	labelResult  = "tafs-chrome/result/v1"
)

// The codes planque closes a connection with.
const (
	closeUnknownChrome = 4401 // no such Chrome, or a proof that isn't its key's
	closeBuildRefused  = 4403 // the extension's build isn't a released one
	closeUnpaired      = 4410 // the Chrome was unpaired: the extension wipes itself
)

const (
	// chromeConnKept bounds a connection's record: API Gateway ends any
	// connection after two hours.
	chromeConnKept = 2 * time.Hour
	// chromeUnpairedKept is how long an unpaired Chrome is still told so
	// (4410) when it connects, rather than taken for an unknown one.
	chromeUnpairedKept = 30 * 24 * time.Hour
	// chromeResultKept is how long a step's signed result waits to be read.
	chromeResultKept = 10 * time.Minute
	// maxRelayMessage bounds what the relay reads: API Gateway's own limit
	// for one message is 128 KB.
	maxRelayMessage = 128 << 10
)

// chromeIndex finds a Chrome's account from its id alone, as the
// extension's connect names it. Unpaired keeps it for a while, so the
// extension is told it was unpaired.
type chromeIndex struct {
	Tenant   string `json:"tenant"`
	Unpaired bool   `json:"unpaired,omitempty"`
}

func pkChromeIndex(id string) string { return "CHROMEID#" + id }

// chromeConn is one connection: whose Chrome, and how far its proofs went.
type chromeConn struct {
	Chrome string `json:"chrome"`
	Tenant string `json:"tenant,omitempty"`
	// Close, when set, is the code the connection is closed with on its
	// first message: the Chrome is unknown, or was unpaired.
	Close  int    `json:"close,omitempty"`
	NonceP string `json:"nonce_p,omitempty"`
	Ready  bool   `json:"ready,omitempty"`
}

func pkChromeConn(conn string) string { return "CHROMECONN#" + conn }

// skChromeLive is a Chrome's live, proven connection.
func skChromeLive(chrome string) string { return "CHROMELIVE#" + chrome }

func pkChromeResult(chrome, step string) string { return "CHROMERESULT#" + chrome + "#" + step }

// relayOn says the relay runs: the owner's Chrome is offered, and planque
// has its step key and the WebSocket API to post through.
func (s *Server) relayOn() bool {
	return s.chromeOn() && s.ChromeSigner != nil && s.ChromePoster != nil
}

// WebSocket is the relay's side of API Gateway's WebSocket API.
func (s *Server) WebSocket(ctx context.Context, ev lambda.WSEvent) int {
	if !s.relayOn() {
		return http.StatusForbidden
	}
	switch ev.Route {
	case "$connect":
		return s.relayConnect(ctx, ev)
	case "$disconnect":
		s.relayDisconnect(ctx, ev.Connection)
		return http.StatusOK
	case "$default":
		if len(ev.Body) > maxRelayMessage {
			s.relayClose(ctx, ev.Connection, closeUnknownChrome, "message too large")
			return http.StatusOK
		}
		s.relayMessage(ctx, ev.Connection, ev.Body)
		return http.StatusOK
	}
	return http.StatusBadRequest
}

// relayConnect records the connection and whose Chrome it says it is. A
// Chrome planque doesn't know is still let in, so it can be told why on its
// first message: unknown (4401) or unpaired (4410).
func (s *Server) relayConnect(ctx context.Context, ev lambda.WSEvent) int {
	id := ev.Query["chrome"]
	conn := chromeConn{Chrome: id, Close: closeUnknownChrome}
	var ix chromeIndex
	if chromeIDRe.MatchString(id) {
		if _, err := getDoc(ctx, s.KV, pkChromeIndex(id), "-", &ix); err == nil {
			switch {
			case ix.Unpaired:
				conn.Close = closeUnpaired
			default:
				var c Chrome
				if _, err := getDoc(ctx, s.KV, tenantPK(ix.Tenant), skChrome(id), &c); err == nil {
					conn.Tenant, conn.Close = ix.Tenant, 0
				}
			}
		} else if !errors.Is(err, errNotFound) {
			return http.StatusServiceUnavailable
		}
	}
	it := store.Item{PK: pkChromeConn(ev.Connection), SK: "-", TTL: s.now().Add(chromeConnKept).Unix()}
	if err := putDoc(ctx, s.KV, it, conn, store.Cond{}); err != nil {
		return http.StatusServiceUnavailable
	}
	return http.StatusOK
}

func (s *Server) relayDisconnect(ctx context.Context, connID string) {
	var conn chromeConn
	if _, err := getDoc(ctx, s.KV, pkChromeConn(connID), "-", &conn); err != nil {
		return
	}
	_, _, _ = s.KV.Delete(ctx, pkChromeConn(connID), "-")
	if conn.Tenant == "" {
		return
	}
	var live struct {
		Conn string `json:"conn"`
	}
	if _, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChromeLive(conn.Chrome), &live); err == nil && live.Conn == connID {
		_, _, _ = s.KV.Delete(ctx, tenantPK(conn.Tenant), skChromeLive(conn.Chrome))
	}
}

// relayMessage takes one message: open, hello, ping, or a step's result.
func (s *Server) relayMessage(ctx context.Context, connID string, body []byte) {
	var m map[string]json.RawMessage
	if json.Unmarshal(body, &m) != nil {
		s.relayClose(ctx, connID, closeUnknownChrome, "not a message")
		return
	}
	var kind string
	_ = json.Unmarshal(m["type"], &kind)
	it, err := getDoc(ctx, s.KV, pkChromeConn(connID), "-", new(chromeConn))
	if err != nil {
		s.relayClose(ctx, connID, closeUnknownChrome, "no such connection")
		return
	}
	var conn chromeConn
	_ = json.Unmarshal(it.Doc, &conn)
	if conn.Close != 0 {
		s.relayClose(ctx, connID, conn.Close, map[int]string{closeUnknownChrome: "no such Chrome", closeUnpaired: "this Chrome was unpaired"}[conn.Close])
		return
	}
	switch {
	case kind == "open" && !conn.Ready && conn.NonceP == "":
		conn.NonceP = oauth.Random(24)
		cond := store.Cond{Ver: &it.Ver}
		next := store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}
		if putDoc(ctx, s.KV, next, conn, cond) != nil {
			s.relayClose(ctx, connID, closeUnknownChrome, "the connection couldn't be kept")
			return
		}
		s.relayPost(ctx, connID, map[string]any{"type": "challenge", "nonce_p": conn.NonceP})
	case kind == "hello" && !conn.Ready && conn.NonceP != "":
		s.relayHello(ctx, connID, it, conn, m)
	case kind == "ping" && conn.Ready:
		s.relayPost(ctx, connID, map[string]any{"type": "pong"})
	case kind == "pong" && conn.Ready:
	case kind == "result" && conn.Ready:
		s.relayResult(ctx, conn, m)
	default:
		// Anything else, or anything before the proofs, ends it.
		s.relayClose(ctx, connID, closeUnknownChrome, "unexpected message")
	}
}

// relayHello checks the extension's proof: its paired key signed planque's
// challenge, from a released build. Then planque proves itself.
func (s *Server) relayHello(ctx context.Context, connID string, it store.Item, conn chromeConn, m map[string]json.RawMessage) {
	body := map[string]json.RawMessage{}
	for _, k := range []string{"chrome_id", "nonce_p", "nonce_e", "version", "manifest_hash"} {
		v, ok := m[k]
		if !ok {
			s.relayClose(ctx, connID, closeUnknownChrome, "hello lacks "+k)
			return
		}
		body[k] = v
	}
	var hello struct {
		ChromeID     string `json:"chrome_id"`
		NonceP       string `json:"nonce_p"`
		NonceE       string `json:"nonce_e"`
		Version      string `json:"version"`
		ManifestHash string `json:"manifest_hash"`
	}
	raw, _ := json.Marshal(body)
	var sig string
	if json.Unmarshal(raw, &hello) != nil || json.Unmarshal(m["sig"], &sig) != nil {
		s.relayClose(ctx, connID, closeUnknownChrome, "malformed hello")
		return
	}
	if hello.ChromeID != conn.Chrome || hello.NonceP != conn.NonceP || hello.NonceE == "" || len(hello.NonceE) > 128 {
		s.relayClose(ctx, connID, closeUnknownChrome, "hello isn't for this connection")
		return
	}
	if b, ok := s.ChromeBuilds[hello.ManifestHash]; !ok || b.Version != hello.Version {
		s.relayClose(ctx, connID, closeBuildRefused, "this build of the extension isn't a released one")
		return
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChrome(conn.Chrome), &c); err != nil {
		s.relayClose(ctx, connID, closeUnpaired, "this Chrome was unpaired")
		return
	}
	if !stepkey.Verify(c.Key, labelHello, body, sig) {
		s.relayClose(ctx, connID, closeUnknownChrome, "hello isn't signed by this Chrome's key")
		return
	}
	conn.Ready = true
	next := store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}
	if putDoc(ctx, s.KV, next, conn, store.Cond{Ver: &it.Ver}) != nil {
		s.relayClose(ctx, connID, closeUnknownChrome, "the connection couldn't be kept")
		return
	}
	live := store.Item{PK: tenantPK(conn.Tenant), SK: skChromeLive(conn.Chrome), TTL: s.now().Add(chromeConnKept).Unix()}
	if putDoc(ctx, s.KV, live, map[string]string{"conn": connID}, store.Cond{}) != nil {
		s.relayClose(ctx, connID, closeUnknownChrome, "the connection couldn't be kept")
		return
	}
	welcome, err := stepkey.SignBody(ctx, s.ChromeSigner, labelWelcome, map[string]string{"chrome_id": conn.Chrome, "nonce_p": conn.NonceP, "nonce_e": hello.NonceE})
	if err != nil {
		s.log().ErrorContext(ctx, "chrome relay: welcome not signed", slog.String("err", err.Error()))
		s.relayClose(ctx, connID, closeUnknownChrome, "planque couldn't prove itself")
		return
	}
	s.relayPost(ctx, connID, map[string]any{"type": "welcome", "sig": welcome})
}

// relayResult keeps a step's result, once its signature is the Chrome's: a
// result is never taken on planque's word or the network's.
func (s *Server) relayResult(ctx context.Context, conn chromeConn, m map[string]json.RawMessage) {
	body := map[string]json.RawMessage{}
	for _, k := range []string{"id", "ok", "page", "error"} {
		v, ok := m[k]
		if !ok {
			return
		}
		body[k] = v
	}
	var id, sig string
	if json.Unmarshal(m["id"], &id) != nil || id == "" || len(id) > 64 || json.Unmarshal(m["sig"], &sig) != nil {
		return
	}
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChrome(conn.Chrome), &c); err != nil || !stepkey.Verify(c.Key, labelResult, body, sig) {
		s.log().WarnContext(ctx, "chrome relay: a result not signed by its Chrome", slog.String("chrome", conn.Chrome))
		return
	}
	it := store.Item{PK: pkChromeResult(conn.Chrome, id), SK: "-", TTL: s.now().Add(chromeResultKept).Unix()}
	// The first result for a step stands: a second can't replace it.
	_ = putDoc(ctx, s.KV, it, body, store.Cond{NotExists: true})
}

func (s *Server) relayPost(ctx context.Context, connID string, msg any) {
	b, _ := json.Marshal(msg)
	if err := s.ChromePoster.Post(ctx, connID, b); err != nil {
		s.log().WarnContext(ctx, "chrome relay: not posted", slog.String("err", err.Error()))
	}
}

// relayClose tells the extension why, then drops the connection.
func (s *Server) relayClose(ctx context.Context, connID string, code int, reason string) {
	s.relayPost(ctx, connID, map[string]any{"type": "close", "code": code, "reason": reason})
	if err := s.ChromePoster.Drop(ctx, connID); err != nil {
		s.log().WarnContext(ctx, "chrome relay: not dropped", slog.String("err", err.Error()))
	}
	_, _, _ = s.KV.Delete(ctx, pkChromeConn(connID), "-")
}

// chromeStep is one signed step, as the extension checks it.
type chromeStep struct {
	ID       string          `json:"id"`
	Nonce    string          `json:"nonce"`
	Exp      int64           `json:"exp"` // ms since the epoch, at most 30 s ahead
	ChromeID string          `json:"chrome_id"`
	Run      string          `json:"run"`
	Grant    string          `json:"grant"`
	Sites    []string        `json:"sites"`
	Op       string          `json:"op"`
	URL      string          `json:"url,omitempty"`
	ReadID   string          `json:"read_id,omitempty"`
	Act      json.RawMessage `json:"act,omitempty"`
}

// errChromeOffline: no proven connection from that Chrome now.
var errChromeOffline = errors.New("the Chrome isn't connected")

// sendChromeStep signs step and sends it to the Chrome's live connection.
func (s *Server) sendChromeStep(ctx context.Context, tenant string, step chromeStep) error {
	var live struct {
		Conn string `json:"conn"`
	}
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skChromeLive(step.ChromeID), &live); err != nil {
		if errors.Is(err, errNotFound) {
			return errChromeOffline
		}
		return err
	}
	step.Nonce = oauth.Random(16)
	step.Exp = s.now().Add(30 * time.Second).UnixMilli()
	sig, err := stepkey.SignBody(ctx, s.ChromeSigner, labelStep, step)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(map[string]any{"type": "step", "step": step, "sig": sig})
	return s.ChromePoster.Post(ctx, live.Conn, b)
}

// chromeResult is a step's result, once the Chrome sent one: found is
// false while it hasn't.
func (s *Server) chromeResult(ctx context.Context, chrome, step string) (map[string]json.RawMessage, bool, error) {
	var body map[string]json.RawMessage
	if _, err := getDoc(ctx, s.KV, pkChromeResult(chrome, step), "-", &body); err != nil {
		if errors.Is(err, errNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return body, true, nil
}

// closeChrome ends a Chrome's live connection with code, if it has one.
func (s *Server) closeChrome(ctx context.Context, tenant, chrome string, code int, reason string) {
	if s.ChromePoster == nil {
		return
	}
	var live struct {
		Conn string `json:"conn"`
	}
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skChromeLive(chrome), &live); err != nil {
		return
	}
	s.relayClose(ctx, live.Conn, code, reason)
	_, _, _ = s.KV.Delete(ctx, tenantPK(tenant), skChromeLive(chrome))
}
