// Package kiff is the part of KIFF planque calls: who a key speaks for
// (apps/api /v1/me), whether a person's device signed an approval
// (/v1/devices/verify), and the gateway's connections, where planque
// registers a vendor's tools in relay mode (kiff-cloud RFC 046 D1, D6).
package kiff

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// Client calls KIFF with the caller's own key.
type Client struct {
	API     string // https://api.kiff.dev
	Gateway string // https://mcp.kiff.dev
	HTTP    *http.Client
}

// Me is who a key speaks for, as apps/api answers it.
type Me struct {
	APIKey        bool     `json:"api_key"`
	TenantID      string   `json:"tenant_id"`
	AgentID       string   `json:"agent_id"`
	Roles         []string `json:"roles"`
	Capabilities  []string `json:"capabilities"`
	DeviceSession bool     `json:"device_session"`
}

// ErrUnauthorized: KIFF refused the key.
var ErrUnauthorized = errors.New("KIFF did not accept the key")

// connectorAllowed are the only capabilities a connector key may hold:
// planque keeps it, so it reads, connects tools and writes the gateway's
// domain, nothing more. A capability KIFF adds later is refused until it
// is added here.
var connectorAllowed = []string{"read", "personal", "tools.connect", "domains.manage"}

// ConnectorProblem says why a key is not a connector key for tenant, or "".
func (m Me) ConnectorProblem(tenant string) string {
	switch {
	case !m.APIKey || m.DeviceSession:
		return "the connector credential must be a KIFF API key"
	case m.TenantID != tenant:
		return "the connector key is for another KIFF account"
	case m.AgentID != "":
		return "the connector key must not be an agent's key"
	case !slices.Contains(m.Capabilities, "tools.connect"):
		return "the connector key cannot connect tools"
	}
	for _, c := range m.Capabilities {
		if !slices.Contains(connectorAllowed, c) {
			return "the connector key can do more than connect tools (" + c + "); planque keeps only a connector key"
		}
	}
	return ""
}

// Me asks apps/api who key speaks for.
func (c *Client) Me(ctx context.Context, key string) (Me, error) {
	var me Me
	status, err := c.do(ctx, http.MethodGet, c.API+"/v1/me", key, nil, &me)
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return Me{}, ErrUnauthorized
	}
	if err != nil {
		return Me{}, err
	}
	return me, nil
}

// Verified is KIFF's answer on a device's signature (kiff-cloud #1267).
type Verified struct {
	Valid bool `json:"valid"`
	// Why not: unknown_device, signature_invalid, not_a_member or not_allowed.
	Reason   string `json:"reason"`
	DeviceID string `json:"device_id"`
	Subject  string `json:"subject"` // the person whose device it is
	Platform string `json:"platform"`
	Attested bool   `json:"attested"`
}

// VerifyDevice asks KIFF, with the account's connector key, whether its
// device deviceID signed msg (sig: base64 DER ECDSA over its SHA-256), and
// whether that device's person may do what capability names. KIFF keeps
// the devices a person approves with, so planque keeps none
// (kiff/planque#10).
func (c *Client) VerifyDevice(ctx context.Context, connectorKey, deviceID string, msg []byte, sig, capability string) (Verified, error) {
	var v Verified
	status, err := c.do(ctx, http.MethodPost, c.API+"/v1/devices/verify", connectorKey, map[string]string{
		"device_id": deviceID, "message": base64.StdEncoding.EncodeToString(msg), "signature": sig, "capability": capability}, &v)
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return Verified{}, ErrUnauthorized
	}
	if err != nil {
		return Verified{}, err
	}
	return v, nil
}

// Connection is a tool registered in KIFF's gateway, as planque sends it.
type Connection struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Tool           string `json:"tool"`
	Description    string `json:"description,omitempty"`
	AmountArg      string `json:"amount_argument,omitempty"`
	AmountUnit     string `json:"amount_unit,omitempty"`
	AmountScale    *int   `json:"amount_scale,omitempty"`
	IdempotencyArg string `json:"idempotency_argument"`
	ReadOnly       bool   `json:"read_only"`
	Mode           string `json:"mode"`
	Audience       string `json:"audience"`
	Credential     string `json:"credential"` // the relay's transport credential
}

// Connect registers a tool. KIFF lists planque's relay before it accepts.
func (c *Client) Connect(ctx context.Context, key string, conn Connection) error {
	_, err := c.do(ctx, http.MethodPost, c.Gateway+"/v1/connections", key, conn, nil)
	return err
}

// Disconnect removes a tool. A tool already gone is not an error.
func (c *Client) Disconnect(ctx context.Context, key, name string) error {
	status, err := c.do(ctx, http.MethodDelete, c.Gateway+"/v1/connections/"+url.PathEscape(name), key, nil, nil)
	if status == http.StatusNotFound {
		return nil
	}
	return err
}

func (c *Client) do(ctx context.Context, method, u, key string, in, out any) (int, error) {
	var body io.Reader
	if in != nil {
		b, _ := json.Marshal(in)
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "planque")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(b, &e)
		msg := strings.TrimSpace(e.Error)
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return resp.StatusCode, fmt.Errorf("KIFF answered %d: %s", resp.StatusCode, msg)
	}
	if out != nil {
		if err := json.Unmarshal(b, out); err != nil {
			return resp.StatusCode, err
		}
	}
	return resp.StatusCode, nil
}
