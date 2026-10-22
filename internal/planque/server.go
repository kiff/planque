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
	// AppReturn is where the browser that completed a vendor sign-in is
	// sent, with the code that confirms it: the owner's app (tafs), for
	// example tafs://planque/signed-in. The operator sets it; no request
	// can, so a sign-in link someone else opens never hands its code back
	// to whoever approved the grant.
	AppReturn string
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
	mux.HandleFunc("POST /v1/sweep", s.sweep)
	mux.HandleFunc("POST /v1/accounts/{tenant}/enrolment", s.enrol)
	mux.HandleFunc("POST /v1/accounts/{tenant}/halts/confirm", s.confirmHalt)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants", s.requestGrant)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}", s.getGrant)
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants/{id}/catalog", s.catalog)
	// Reading what a sign-in covers, before its tools were approved, could
	// read business data no approval named (kiff/planque#23): it's gone.
	// Setup asks the owner which mode instead.
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/describe", func(w http.ResponseWriter, _ *http.Request) {
		writeErr(w, http.StatusGone, "reading what a sign-in covers was removed: nothing is read through a vendor's tools before the owner approves them. Ask the owner which account and mode it is")
	})
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/tools", s.setTools)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/read", s.setupRead)
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/approve", s.signed(s.approve))
	mux.HandleFunc("POST /v1/accounts/{tenant}/grants/{id}/confirm", s.signed(s.confirm))
	mux.HandleFunc("DELETE /v1/accounts/{tenant}/grants/{id}", s.signed(s.disconnect))
	mux.HandleFunc("GET /v1/accounts/{tenant}/grants", s.signed(s.listGrants))
	mux.HandleFunc("GET /v1/accounts/{tenant}/export", s.signedFor(capExport, s.export))
	// Deleting the account (kiff/tafs#185): everything planque keeps for it.
	mux.HandleFunc("DELETE /v1/accounts/{tenant}", s.eraseAccount)
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

// signedRequest is a request the owner signed on their device.
type signedRequest struct {
	tenant    string
	body      []byte
	enrolment Enrolment
	approval  Approval // the signature, as evidence
}

// verifySigned checks Planque-Device-Signature over r, and that this exact
// signed request was never made before. The device is one KIFF knows (d=),
// verified by KIFF, asked with the account's connector key, which opens
// connector (kiff/planque#10). A signature without d=, from a phone that
// enrolled its own key before KIFF kept devices, is no longer taken
// (kiff/planque#12).
// What KIFF is asked a signer may do, as KIFF's own roles say: connect a
// vendor, or, only for those who may in KIFF too, export or delete the
// account (kiff/tafs#185, kiff/planque#46).
const (
	capConnect = "tools.connect"
	capExport  = "account.export"
	capDelete  = "account.delete"
)

// capWords says what a capability lets the signer do, for a refusal.
var capWords = map[string]string{capConnect: "connect tools", capExport: "export the account", capDelete: "delete the account"}

func (s *Server) verifySigned(r *http.Request, tenant, capability string, body []byte, connector func() (string, error)) (Approval, error) {
	sg, err := device.ParseHeader(r.Header.Get(device.Header))
	if err != nil {
		return Approval{}, apiError{http.StatusUnauthorized, "this request needs the owner's signature from their device (" + device.Header + ")"}
	}
	msg := device.Message(s.BaseURL, tenant, r.Method, r.URL.RequestURI(), body, sg.At)
	a := Approval{Message: string(msg), Body: string(body), Signature: sg.Signature, SignedAt: sg.At, ApprovedAt: s.now().UTC()}
	// What the replay record names as the signer, with the message.
	if sg.Device == "" {
		return Approval{}, apiError{http.StatusUnauthorized, "this device isn't known to KIFF: sign in to KIFF on it again"}
	}
	var signer []byte
	{
		// planque keeps the window and the replay check: KIFF says only
		// that the device signed these bytes, not when.
		if err := device.CheckTime(sg.At, s.now()); err != nil {
			return Approval{}, apiError{http.StatusUnauthorized, err.Error()}
		}
		key, err := connector()
		if err != nil {
			return Approval{}, err
		}
		v, err := s.KIFF.VerifyDevice(r.Context(), key, sg.Device, msg, sg.Signature, capability)
		if err != nil {
			s.log().ErrorContext(r.Context(), "device check failed", slog.String("tenant", tenant), slog.String("err", err.Error()))
			return Approval{}, apiError{http.StatusBadGateway, "KIFF could not check the device that signed"}
		}
		if !v.Valid {
			return Approval{}, deviceRefused(v.Reason, capability)
		}
		a.DeviceID, a.Subject, a.DeviceName, a.Fingerprint = v.DeviceID, v.Subject, strings.TrimSpace(v.Name), v.DeviceID
		signer = []byte("kiff-device\n" + v.DeviceID + "\n")
	}
	sum := sha256.Sum256(append(append([]byte{}, signer...), msg...))
	ttl := s.now().Add(2 * device.MaxSkew).Unix()
	err = s.KV.Put(r.Context(), store.Item{PK: pkSigned(hex.EncodeToString(sum[:])), SK: "-", Doc: []byte("{}"), TTL: ttl},
		store.Cond{NotExists: true})
	if errors.Is(err, store.ErrCondition) {
		return Approval{}, apiError{http.StatusConflict, "this signed request was already made; sign a new one"}
	}
	if err != nil {
		return Approval{}, err
	}
	return a, nil
}

// deviceRefused says why KIFF didn't vouch for a device's signature.
func deviceRefused(reason, capability string) apiError {
	switch reason {
	case "unknown_device":
		return apiError{http.StatusUnauthorized, "KIFF doesn't know the device that signed, or it was removed"}
	case "signature_invalid":
		return apiError{http.StatusUnauthorized, device.ErrSignature.Error()}
	case "not_a_member":
		return apiError{http.StatusForbidden, "the person whose device signed isn't a member of this account"}
	case "not_allowed":
		return apiError{http.StatusForbidden, "the person whose device signed can't " + capWords[capability] + " on this account"}
	}
	return apiError{http.StatusUnauthorized, "KIFF didn't vouch for the device that signed"}
}

// signed admits a request only with the owner's signature: from a device
// KIFF knows, of a person who may connect tools.
func (s *Server) signed(next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return s.signedFor(capConnect, next)
}

// signedFor is signed, for a person KIFF says holds capability.
func (s *Server) signedFor(capability string, next func(http.ResponseWriter, *http.Request, signedRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := r.PathValue("tenant")
		body, err := readBody(r)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		var e Enrolment
		if _, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &e); errors.Is(err, errNotFound) {
			writeErr(w, http.StatusForbidden, "this account isn't set up to connect tools yet")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		a, err := s.verifySigned(r, tenant, capability, body, func() (string, error) { return s.openConnector(r.Context(), e) })
		if err != nil {
			s.fail(w, r, err)
			return
		}
		next(w, r, signedRequest{tenant: tenant, body: body, enrolment: e, approval: a})
	}
}

// enrol hands planque the account's KIFF connector key. A device KIFF
// knows (d=) signs it, and KIFF says whether its person may connect tools:
// any of the person's devices may, and the latest connector key is kept
// (kiff/planque#10). Without d=, a phone enrols its own key and signs with
// it, as before: a different phone can't replace an enrolled one, and the
// enrolled phone may re-enrol to hand planque a new connector key.
func (s *Server) enrol(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		// DeviceKey is what a phone sent before KIFF kept devices. It is
		// still accepted, so apps that send it keep working, and ignored:
		// the device is KIFF's (d=), kiff/planque#12.
		DeviceKey    string `json:"device_key,omitempty"`
		ConnectorKey string `json:"connector_key"`
	}
	if err := decodeStrict(body, &req); err != nil {
		s.fail(w, r, err)
		return
	}
	connector := strings.TrimSpace(req.ConnectorKey)
	var prev Enrolment
	prevItem, err := getDoc(r.Context(), s.KV, tenantPK(tenant), skEnrolment, &prev)
	// The write below is conditional on what was read here, so two
	// enrolments racing cannot both pass this check.
	cond := store.Cond{NotExists: true}
	if err == nil {
		cond = store.Cond{Ver: &prevItem.Ver}
	} else if !errors.Is(err, errNotFound) {
		s.fail(w, r, err)
		return
	}
	// The connector key first: KIFF is asked with it whether the device signed.
	me, err := s.KIFF.Me(r.Context(), connector)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "KIFF did not accept the connector key")
		return
	}
	if p := me.ConnectorProblem(tenant); p != "" {
		writeErr(w, http.StatusBadRequest, p)
		return
	}
	a, err := s.verifySigned(r, tenant, capConnect, body, func() (string, error) { return connector, nil })
	if err != nil {
		s.fail(w, r, err)
		return
	}
	sealed, err := s.Sealer.Seal(r.Context(), []byte(connector), map[string]string{"tenant": tenant, "kind": "connector"})
	if err != nil {
		s.fail(w, r, err)
		return
	}
	e := Enrolment{Tenant: tenant, Connector: sealed, EnrolledAt: s.now().UTC()}
	err = putDoc(r.Context(), s.KV, store.Item{PK: tenantPK(tenant), SK: skEnrolment, Ver: prevItem.Ver + 1}, e, cond)
	if errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, "the enrolment changed while this one was checked; enrol again")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(r.Context(), "device enrolled", slog.String("tenant", tenant), slog.String("device", a.Fingerprint))
	writeJSON(w, http.StatusCreated, map[string]string{"tenant": tenant, "device": a.Fingerprint})
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
	return s.openConnector(ctx, e)
}

// openConnector opens an enrolment's KIFF connector key.
func (s *Server) openConnector(ctx context.Context, e Enrolment) (string, error) {
	k, err := s.Sealer.Open(ctx, e.Connector, map[string]string{"tenant": e.Tenant, "kind": "connector"})
	return string(k), err
}
