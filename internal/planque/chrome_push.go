package planque

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// labelPush is what the extension signs its push subscription under (kiff/tafs
// chrome/src/crypto.js LABEL.push).
const labelPush = "tafs-chrome/push/v1"

// pushHost is the only push service a subscription may name: Chrome's. planque
// posts to the endpoint to wake the extension, so an endpoint is never a URL a
// client picked.
const pushHost = "fcm.googleapis.com"

// ChromePush is a Chrome's Web Push subscription, as its extension sent and
// signed it after the welcome: where planque sends the push that wakes it for a
// step while it isn't connected. One per Chrome; a new one replaces it. Never
// in the Chromes list: the endpoint is the extension's own address.
type ChromePush struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// Refusal codes for a push message, answered as {type: "push_refused", code}.
const (
	codePushUnsigned = "push_unsigned" // the signature isn't this Chrome's key's over it
	codePushEndpoint = "push_endpoint" // not https on Chrome's push service, or too long
	codePushKeys     = "push_keys"     // the keys aren't a P-256 point and a 16-byte secret
	codePushNotKept  = "push_not_kept" // planque couldn't keep it now; send it again
)

// pushRefusal is why sub can't be kept, or "".
func (p ChromePush) pushRefusal() string {
	if len(p.Endpoint) > 1024 {
		return codePushEndpoint
	}
	u, err := url.Parse(p.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host != pushHost || u.User != nil || u.Path == "" || u.Path == "/" {
		return codePushEndpoint
	}
	point, err := b64any(p.Keys.P256DH)
	if err != nil || len(point) != 65 || point[0] != 4 {
		return codePushKeys
	}
	if auth, err := b64any(p.Keys.Auth); err != nil || len(auth) != 16 {
		return codePushKeys
	}
	return ""
}

// b64any decodes base64url, padded or not, as a browser's toJSON may give it.
func b64any(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(strings.TrimRight(s, "="))
}

// relayPush keeps the subscription a connected Chrome sent, when its own key
// signed it and it names Chrome's push service. A refusal answers why and keeps
// the connection: a run can still reach this Chrome while it's connected.
func (s *Server) relayPush(ctx context.Context, connID string, conn chromeConn, m map[string]json.RawMessage) {
	refuse := func(code string) { s.relayPost(ctx, connID, map[string]any{"type": "push_refused", "code": code}) }
	raw, ok := m["subscription"]
	var sig string
	if !ok || json.Unmarshal(m["sig"], &sig) != nil {
		refuse(codePushUnsigned)
		return
	}
	var sub ChromePush
	if json.Unmarshal(raw, &sub) != nil {
		refuse(codePushEndpoint)
		return
	}
	var c Chrome
	it, err := getDoc(ctx, s.KV, tenantPK(conn.Tenant), skChrome(conn.Chrome), &c)
	if err != nil {
		s.relayClose(ctx, connID, closeUnpaired, "this Chrome was unpaired")
		return
	}
	// Over the subscription exactly as sent (toJSON's expirationTime and all).
	if !stepkey.Verify(c.Key, labelPush, raw, sig) {
		refuse(codePushUnsigned)
		return
	}
	if code := sub.pushRefusal(); code != "" {
		refuse(code)
		return
	}
	if c.Push == nil || !samePush(*c.Push, sub) {
		if beforePushKept != nil {
			beforePushKept()
		}
		c.Push = &sub
		// On the version read: an unpairing in between wins, and nothing re-creates the Chrome.
		if putDoc(ctx, s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}, c, store.Cond{Ver: &it.Ver}) != nil {
			refuse(codePushNotKept)
			return
		}
	}
	s.relayPost(ctx, connID, map[string]any{"type": "push_kept"})
}

// beforePushKept runs, in tests only, between reading the Chrome and keeping
// its subscription.
var beforePushKept func()

func samePush(a, b ChromePush) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}
