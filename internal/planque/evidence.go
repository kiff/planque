package planque

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/kiff/planque/internal/browser"
)

// The evidence pictures (browser batch 1): a picture of the page after a
// hosted browser's write, kept with the operation, sealed under the grant
// as a download is (browser_files.go), for 30 days, deleted sooner with the
// grant or the account. Never on a page that pays, and never while a card,
// password or code field holds anything (browser.Pictured, the threat
// model's T3.2). On with PLANQUE_EVIDENCE=on (cmd/planque).

// sealedEvidence is browser.Evidence on planque's own table.
type sealedEvidence struct{ s *Server }

// SealedEvidence is the evidence store that keeps pictures sealed under
// the grant, in planque's table.
func (s *Server) SealedEvidence() browser.Evidence { return sealedEvidence{s} }

// Put keeps one PNG for op, and answers its reference. Anything that isn't
// a PNG, or is larger than 5 MB, is refused.
func (e sealedEvidence) Put(ctx context.Context, tenant, grant, op string, png []byte) (string, error) {
	if len(png) > browser.MaxDownload {
		return "", errors.New("the picture is larger than 5 MB")
	}
	if typ, err := browser.TypeOf(png); err != nil || typ != browser.TypePNG {
		return "", errors.New("the picture isn't a PNG")
	}
	sum := sha256.Sum256(png)
	f := browserFile{ID: "ev_" + rand.Text(), Grant: grant, Kind: "evidence", Name: op + ".png", Type: browser.TypePNG,
		Size: len(png), SHA256: hex.EncodeToString(sum[:]), At: e.s.now().UTC()}
	if err := e.s.putBlob(ctx, tenant, &f, png); err != nil {
		return "", err
	}
	return f.ID, nil
}

// Drop deletes every picture kept under the grant: dropBrowserFiles deletes
// the grant's blobs of every kind.
func (e sealedEvidence) Drop(ctx context.Context, tenant, grant string) error {
	return e.s.dropBrowserFiles(ctx, &Grant{Tenant: tenant, ID: grant})
}

// grantEvidence answers GET /v1/accounts/{tenant}/grants/{id}/evidence/{ref},
// with the account's key: the picture, as a PNG, for the taf's Record.
func (s *Server) grantEvidence(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	g, err := s.loadGrant(r.Context(), tenant, r.PathValue("id"))
	if errors.Is(err, errNotFound) {
		writeErr(w, http.StatusNotFound, "no such grant")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	_, data, err := s.openBlob(r.Context(), g, r.PathValue("ref"), "evidence")
	if errors.Is(err, errFileGone) {
		writeErr(w, http.StatusNotFound, "no such picture, or it was deleted after 30 days")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	w.Header().Set("Content-Type", browser.TypePNG)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	_, _ = w.Write(data)
}
