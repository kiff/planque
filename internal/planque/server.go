// Package planque is the service: an owner's phone enrols, an agent asks
// for a vendor account, the owner approves on the phone, signs in at the
// vendor, and planque relays each call KIFF permits to the vendor with a
// token only planque holds.
package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

// IdempotencyArg is the argument KIFF's gateway fills with the call's
// operation id; the permit binds it, and planque removes it before the
// vendor sees the call.
const IdempotencyArg = "planque_op"

// Server is planque.
type Server struct {
	// BaseURL is where planque is reached publicly, without a trailing
	// slash: the vendor redirects the owner back to BaseURL/oauth/callback.
	BaseURL string
	KV      store.KV
	Sealer  seal.Sealer
	KIFF    *kiff.Client
	Permits *permit.Verifier
	OAuth   *oauth.Client
	MCP     *mcp.Client
	Logger  *slog.Logger
	Now     func() time.Time
	// RefreshWait bounds how long a call waits for another one's refresh.
	RefreshWait time.Duration

	refreshing sync.Map // grant id → *sync.Mutex, one refresh per process
}

func (s *Server) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Server) log() *slog.Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

func (s *Server) redirect() string { return s.BaseURL + "/oauth/callback" }

// Handler serves planque's routes.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]bool{"ok": true}) })
	mux.HandleFunc("POST /v1/accounts/{tenant}/enrolment", s.enrol)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants", s.requestGrant)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}", s.getGrant)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/approve", s.signed(s.approve))
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/grants/{id}", s.signed(s.disconnect))
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants", s.signed(s.listGrants))
	mux.HandleFunc("GET /v1/accounts/{tenant}/export", s.signed(s.export))
	mux.HandleFunc("GET /oauth/callback", s.callback)
	mux.HandleFunc("POST /mcp/{tenant}/{id}", s.relay)
	return mux
}

type apiError struct {
	status int
	msg    string
}

func (e apiError) Error() string { return e.msg }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	var ae apiError
	if errors.As(err, &ae) {
		writeErr(w, ae.status, ae.msg)
		return
	}
	s.log().ErrorContext(r.Context(), "request failed", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
	writeErr(w, http.StatusInternalServerError, "planque could not complete this request")
}

func readBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r.Body, 64<<10+1))
	if err != nil || len(b) > 64<<10 {
		return nil, apiError{http.StatusBadRequest, "the body is too long or unreadable"}
	}
	return b, nil
}

func decodeStrict(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return apiError{http.StatusBadRequest, "the body must be a JSON object with the documented fields"}
	}
	return nil
}

// signedRequest is a request the enrolled phone signed.
type signedRequest struct {
	tenant    string
	body      []byte
	enrolment Enrolment
	approval  Approval // the signature, as evidence
}

// verifySigned checks Planque-Device-Signature over r with keyDER, and
// that this exact signed request was never made before.
func (s *Server) verifySigned(r *http.Request, tenant string, body, keyDER []byte) (Approval, error) {
	signedAt, sig, err := device.ParseHeader(r.Header.Get(device.Header))
	if err != nil {
		return Approval{}, apiError{http.StatusUnauthorized, "this request needs the owner's signature from their phone (" + device.Header + ")"}
	}
	msg := device.Message(tenant, r.Method, r.URL.RequestURI(), body, signedAt)
	if err := device.Verify(keyDER, msg, signedAt, sig, s.now()); err != nil {
		return Approval{}, apiError{http.StatusUnauthorized, err.Error()}
	}
	sum := sha256.Sum256(append(append([]byte{}, keyDER...), msg...))
	ttl := s.now().Add(2 * device.MaxSkew).Unix()
	err = s.KV.Put(r.Context(), store.Item{PK: pkSigned(hex.EncodeToString(sum[:])), SK: "-", Doc: []byte("{}"), TTL: ttl},
		store.Cond{NotExists: true})
	if errors.Is(err, store.ErrCondition) {
		return Approval{}, apiError{http.StatusConflict, "this signed request was already made; sign a new one"}
	}
	if err != nil {
		return Approval{}, err
	}
	return Approval{DeviceKey: keyDER, Fingerprint: device.Fingerprint(keyDER), Message: string(msg),
		Body: string(body), Signature: sig, SignedAt: signedAt, ApprovedAt: s.now().UTC()}, nil
}

// signed admits a request only with the enrolled phone's signature.
func (s *Server) signed(next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := r.PathValue("tenant")
		body, err := readBody(r)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		var e Enrolment
		if _, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &e); errors.Is(err, errNotFound) {
			writeErr(w, http.StatusForbidden, "no phone is enrolled for this account")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		a, err := s.verifySigned(r, tenant, body, e.DeviceKey)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		next(w, r, signedRequest{tenant: tenant, body: body, enrolment: e, approval: a})
	}
}

// enrol binds an owner's phone and planque's KIFF connector key to an
// account. The phone signs the enrolment with the key it enrols. A
// different phone cannot replace an enrolled one; the enrolled phone may
// re-enrol to hand planque a new connector key.
func (s *Server) enrol(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		DeviceKey    string `json:"device_key"`
		ConnectorKey string `json:"connector_key"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	der, err := device.ParseKey(req.DeviceKey)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	var prev Enrolment
	_, err = getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &prev)
	switch {
	case err == nil && !bytes.Equal(prev.DeviceKey, der):
		writeErr(w, http.StatusConflict, "another phone is enrolled for this account")
		return
	case err != nil && !errors.Is(err, errNotFound):
		s.fail(w, r, err)
		return
	}
	if _, err := s.verifySigned(r, tenant, body, der); err != nil {
		s.fail(w, r, err)
		return
	}
	me, err := s.KIFF.Me(r.Context(), strings.TrimSpace(req.ConnectorKey))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "KIFF did not accept the connector key")
		return
	}
	if p := me.ConnectorProblem(tenant); p != "" {
		writeErr(w, http.StatusBadRequest, p)
		return
	}
	sealed, err := s.Sealer.Seal(r.Context(), []byte(strings.TrimSpace(req.ConnectorKey)), map[string]string{"tenant": tenant, "kind": "connector"})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	e := Enrolment{Tenant: tenant, DeviceKey: der, Connector: sealed, EnrolledAt: s.now().UTC()}
	if err := putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(tenant), SK: skEnrolment}, e, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(r.Context(), "phone enrolled", slog.String("tenant", tenant), slog.String("device", device.Fingerprint(der)))
	writeJSON(w, http.StatusCreated, map[string]string{"tenant": tenant, "device": device.Fingerprint(der)})
}

// agentKey checks the KIFF key an agent sends, for tenant.
func (s *Server) agentKey(r *http.Request, tenant string) (kiff.Me, error) {
	key, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || strings.TrimSpace(key) == "" {
		return kiff.Me{}, apiError{http.StatusUnauthorized, "a KIFF key is required"}
	}
	me, err := s.KIFF.Me(r.Context(), strings.TrimSpace(key))
	if errors.Is(err, kiff.ErrUnauthorized) {
		return kiff.Me{}, apiError{http.StatusUnauthorized, "KIFF did not accept the key"}
	}
	if err != nil {
		return kiff.Me{}, apiError{http.StatusBadGateway, "KIFF could not be asked who the key is"}
	}
	if me.TenantID != tenant {
		return kiff.Me{}, apiError{http.StatusForbidden, "the key is for another KIFF account"}
	}
	return me, nil
}

func (s *Server) page(w http.ResponseWriter, status int, title, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.WriteHeader(status)
	_ = pageTmpl.Execute(w, map[string]string{"Title": title, "Message": msg})
}

var pageTmpl = template.Must(template.New("p").Parse(`<!doctype html><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1"><title>{{.Title}} · planque</title>
<style>body{font:17px/1.5 -apple-system,system-ui,sans-serif;max-width:32rem;margin:15vh auto;padding:0 16px;color:#1d1d1f}
@media(prefers-color-scheme:dark){body{background:#111;color:#eee}}</style>
<h1>{{.Title}}</h1><p>{{.Message}}</p>`))

func hostOf(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return u
	}
	return p.Host
}

// connectorKey opens the account's KIFF connector key.
func (s *Server) connectorKey(ctx context.Context, tenant string) (string, error) {
	var e Enrolment
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skEnrolment, &e); err != nil {
		return "", err
	}
	k, err := s.Sealer.Open(ctx, e.Connector, map[string]string{"tenant": tenant, "kind": "connector"})
	return string(k), err
}
