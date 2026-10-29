package planque

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kiff/planque/internal/store"
)

// Deleting an account from one signed intent, driven by tafs cloud. The
// owner's phone signs, with its approval key:
//
//	tafs.delete-account.v2
//	tenant=<tenant id>
//	device=<KIFF device id>
//	issued=<RFC 3339, UTC>
//	expires=<issued + at most 10 minutes>
//	nonce=<32 hex>
//
// one field a line, LF between, no newline at the end. Accept checks it,
// with KIFF (the device signed it, and its person may delete the account)
// and planque's own replay record on its nonce, and keeps it, deleting
// nothing. Finish erases the account as DELETE /v1/accounts/{tenant} does,
// for the accepted intent only, with no deadline; it is the same when asked
// again. A later intent takes the place of an accepted one (KIFF may have
// refused that one, so it never finished). Status says whether this intent
// was accepted, deleting nothing.

const (
	deleteIntentHead = "tafs.delete-account.v2"
	deleteIntentMax  = 10 * time.Minute
	skDeletion       = "DELETION"
)

var deleteNonce = regexp.MustCompile(`^[0-9a-f]{32}$`)

type deleteIntent struct {
	Tenant, Device, Nonce string
	Issued, Expires       time.Time
}

// parseDeleteIntent reads an intent in exactly its signed form.
func parseDeleteIntent(text string) (deleteIntent, bool) {
	lines := strings.Split(text, "\n")
	if len(lines) != 6 || lines[0] != deleteIntentHead {
		return deleteIntent{}, false
	}
	f := map[string]string{}
	for i, k := range []string{"tenant", "device", "issued", "expires", "nonce"} {
		v, ok := strings.CutPrefix(lines[i+1], k+"=")
		if !ok || v == "" {
			return deleteIntent{}, false
		}
		f[k] = v
	}
	issued, err1 := time.Parse(time.RFC3339, f["issued"])
	expires, err2 := time.Parse(time.RFC3339, f["expires"])
	if err1 != nil || err2 != nil || !expires.After(issued) || expires.Sub(issued) > deleteIntentMax || !deleteNonce.MatchString(f["nonce"]) {
		return deleteIntent{}, false
	}
	return deleteIntent{Tenant: f["tenant"], Device: f["device"], Nonce: f["nonce"], Issued: issued, Expires: expires}, true
}

// deletion is an accepted intent, kept on the account until it's erased.
type deletion struct {
	Intent     string    `json:"intent"`
	Signature  string    `json:"signature"`
	AcceptedAt time.Time `json:"accepted_at"`
}

type deletionBody struct {
	Intent    string `json:"intent"`
	Signature string `json:"signature"`
}

func readDeletion(r *http.Request) (deletionBody, deleteIntent, error) {
	body, err := readBody(r)
	if err != nil {
		return deletionBody{}, deleteIntent{}, err
	}
	var req deletionBody
	if err := decodeStrict(body, &req); err != nil {
		return deletionBody{}, deleteIntent{}, err
	}
	in, ok := parseDeleteIntent(req.Intent)
	if !ok || req.Signature == "" {
		return deletionBody{}, deleteIntent{}, apiError{http.StatusBadRequest, "send the signed intent: tafs.delete-account.v2, six lines, expiring at most 10 minutes after it's issued, and its signature"}
	}
	return req, in, nil
}

// acceptDeletion serves POST /v1/account-deletions.
func (s *Server) acceptDeletion(w http.ResponseWriter, r *http.Request) {
	req, in, err := readDeletion(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	ctx := r.Context()
	now := s.now().UTC()
	if !now.Before(in.Expires) || now.Before(in.Issued.Add(-time.Minute)) {
		writeErr(w, http.StatusBadRequest, "the intent expired: sign the deletion again")
		return
	}
	var e Enrolment
	if _, err := getDoc(ctx, s.KV, tenantPK(in.Tenant), skEnrolment, &e); errors.Is(err, errNotFound) {
		// Nothing kept for the account: there's nothing to delete here.
		writeJSON(w, http.StatusAccepted, map[string]any{"accepted_at": now, "nothing_kept": true})
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	key, err := s.openConnector(ctx, e)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	v, err := s.KIFF.VerifyDevice(ctx, key, in.Device, []byte(req.Intent), req.Signature, capDelete)
	if err != nil {
		s.log().ErrorContext(ctx, "device check failed", slog.String("tenant", in.Tenant), slog.String("err", err.Error()))
		writeErr(w, http.StatusBadGateway, "KIFF could not check the device that signed")
		return
	}
	if !v.Valid {
		s.fail(w, r, deviceRefused(v.Reason, capDelete))
		return
	}
	// The same intent again answers as the first time did; its nonce is
	// spent only once.
	var kept deletion
	if _, err := getDoc(ctx, s.KV, tenantPK(in.Tenant), skDeletion, &kept); err == nil && kept.Intent == req.Intent {
		writeJSON(w, http.StatusAccepted, map[string]any{"accepted_at": kept.AcceptedAt})
		return
	}
	sum := sha256.Sum256([]byte("delete-intent\n" + in.Tenant + "\n" + in.Nonce))
	err = s.KV.Put(ctx, store.Item{PK: pkSigned(hex.EncodeToString(sum[:])), SK: "-", Doc: []byte("{}"), TTL: now.Add(2 * deleteIntentMax).Unix()},
		store.Cond{NotExists: true})
	if errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, "this intent's nonce was already used; sign a new one")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	// A later intent takes an earlier one's place: KIFF may have refused that one.
	d := deletion{Intent: req.Intent, Signature: req.Signature, AcceptedAt: now}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(in.Tenant), SK: skDeletion}, d, store.Cond{}); err != nil {
		s.fail(w, r, err)
		return
	}
	s.log().InfoContext(ctx, "account deletion accepted", slog.String("tenant", in.Tenant), slog.String("device", v.DeviceID))
	writeJSON(w, http.StatusAccepted, map[string]any{"accepted_at": now})
}

// accepted is the account's accepted deletion when it is exactly this
// intent and signature. erased: nothing is kept for the account at all.
func (s *Server) accepted(ctx context.Context, req deletionBody, in deleteIntent) (ok, erased bool, err error) {
	items, err := s.KV.Query(ctx, tenantPK(in.Tenant), "")
	if err != nil {
		return false, false, err
	}
	if len(items) == 0 {
		return false, true, nil
	}
	var d deletion
	if _, err := getDoc(ctx, s.KV, tenantPK(in.Tenant), skDeletion, &d); errors.Is(err, errNotFound) {
		return false, false, nil
	} else if err != nil {
		return false, false, err
	}
	same := subtle.ConstantTimeCompare([]byte(d.Intent), []byte(req.Intent)) == 1 &&
		subtle.ConstantTimeCompare([]byte(d.Signature), []byte(req.Signature)) == 1
	return same, false, nil
}

// finishDeletion serves POST /v1/account-deletions/finish: the accepted
// intent erases the account, with no deadline; asked again, it's done.
func (s *Server) finishDeletion(w http.ResponseWriter, r *http.Request) {
	req, in, err := readDeletion(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	ctx := r.Context()
	ok, erased, err := s.accepted(ctx, req, in)
	switch {
	case err != nil:
		s.fail(w, r, err)
		return
	case erased:
		writeJSON(w, http.StatusOK, map[string]any{"tenant": in.Tenant, "erased": true})
		return
	case !ok:
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not_accepted", "detail": "planque finishes only a deletion it accepted, with its signature"})
		return
	}
	var e Enrolment
	if _, err := getDoc(ctx, s.KV, tenantPK(in.Tenant), skEnrolment, &e); err != nil && !errors.Is(err, errNotFound) {
		s.fail(w, r, err)
		return
	}
	grants, problems, err := s.eraseAll(ctx, in.Tenant, e, "deletion intent "+in.Device)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenant": in.Tenant, "erased": true, "grants": grants, "problems": problems})
}

// deletionStatus serves POST /v1/account-deletions/status: whether this
// intent was accepted, deleting nothing.
func (s *Server) deletionStatus(w http.ResponseWriter, r *http.Request) {
	req, in, err := readDeletion(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	ok, erased, err := s.accepted(r.Context(), req, in)
	switch {
	case err != nil:
		s.fail(w, r, err)
	case erased:
		writeJSON(w, http.StatusOK, map[string]any{"status": "erased"})
	case ok:
		writeJSON(w, http.StatusOK, map[string]any{"status": "accepted"})
	default:
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not_accepted"})
	}
}
