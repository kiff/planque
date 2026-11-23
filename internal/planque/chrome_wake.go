package planque

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// chromeWakeWait bounds how long a step waits, after the push, for the woken
// extension to connect; never past the call's own deadline.
const chromeWakeWait = 8 * time.Second

// ChromeWaker sends the push that wakes the owner's Chrome's extension for a
// step while it isn't connected (Web Push, RFC 8030, with VAPID, RFC 8292). The
// push carries no payload: it only says "connect", and the step itself goes out
// signed over the relay once it does, so nothing is encrypted to the
// subscription's keys.
type ChromeWaker struct {
	// Signer is the VAPID key: raw r‖s over SHA-256, as ES256 wants.
	Signer stepkey.Signer
	// Public is its uncompressed P-256 point (65 bytes), the extension's
	// applicationServerKey.
	Public []byte
	// Subject is how the push service can reach planque's operator: an https
	// URL or a mailto.
	Subject string
	HTTP    *http.Client
	Now     func() time.Time
}

// errPushGone is the push service saying the subscription no longer exists.
var errPushGone = errors.New("the push subscription is gone")

// Wake posts one push to sub. A subscription the push service says is gone
// (404, 410) is errPushGone; anything else not accepted is an error.
func (w *ChromeWaker) Wake(ctx context.Context, sub ChromePush) error {
	// Checked again here: planque only ever posts to Chrome's push service.
	if code := sub.pushRefusal(); code != "" {
		return fmt.Errorf("chrome wake: %s", code)
	}
	u, _ := url.Parse(sub.Endpoint)
	jwt, err := w.token(ctx, u.Scheme+"://"+u.Host)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.Endpoint, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("TTL", "30")
	req.Header.Set("Urgency", "high")
	req.Header.Set("Authorization", "vapid t="+jwt+", k="+base64.RawURLEncoding.EncodeToString(w.Public))
	client := w.HTTP
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone:
		return errPushGone
	default:
		return fmt.Errorf("chrome wake: the push service answered %d", resp.StatusCode)
	}
}

// token is the VAPID JWT for aud, the push service's origin: ES256, valid an hour.
func (w *ChromeWaker) token(ctx context.Context, aud string) (string, error) {
	now := time.Now
	if w.Now != nil {
		now = w.Now
	}
	header, _ := json.Marshal(map[string]string{"typ": "JWT", "alg": "ES256"})
	claims, _ := json.Marshal(map[string]any{"aud": aud, "exp": now().Add(time.Hour).Unix(), "sub": w.Subject})
	input := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(claims)
	sig, err := w.Signer.Sign(ctx, []byte(input))
	if err != nil {
		return "", err
	}
	return input + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

// wakeChrome pushes to a Chrome that isn't connected and waits, up to
// chromeWakeWait and never past ctx's deadline, for its proven connection:
// the connection to post the step to. With no waker, no subscription, or no
// connection in time, the Chrome is offline, as before. A subscription the
// push service says is gone is dropped.
func (s *Server) wakeChrome(ctx context.Context, tenant, chrome string) (string, error) {
	if s.ChromeWaker == nil {
		return "", errChromeOffline
	}
	var c Chrome
	it, err := getDoc(ctx, s.KV, tenantPK(tenant), skChrome(chrome), &c)
	if err != nil || c.Push == nil {
		return "", errChromeOffline
	}
	if err := s.ChromeWaker.Wake(ctx, *c.Push); err != nil {
		if errors.Is(err, errPushGone) {
			// On the version read: a subscription sent meanwhile is kept, and an unpairing wins.
			c.Push = nil
			_ = putDoc(ctx, s.KV, store.Item{PK: it.PK, SK: it.SK, Ver: it.Ver + 1, TTL: it.TTL}, c, store.Cond{Ver: &it.Ver})
		}
		s.log().WarnContext(ctx, "chrome wake: not sent", slog.String("chrome", chrome), slog.String("err", err.Error()))
		return "", errChromeOffline
	}
	until := time.Now().Add(chromeWakeWait)
	if d, ok := ctx.Deadline(); ok && d.Before(until) {
		until = d
	}
	for {
		var live struct {
			Conn string `json:"conn"`
		}
		if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skChromeLive(chrome), &live); err == nil && live.Conn != "" {
			return live.Conn, nil
		}
		if !time.Now().Before(until) {
			return "", errChromeOffline
		}
		select {
		case <-ctx.Done():
			return "", errChromeOffline
		case <-time.After(chromePoll):
		}
	}
}
