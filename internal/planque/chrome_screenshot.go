package planque

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// A taf's acting steps in the owner's own Chrome are kept as images, for the
// record (the owner's decision, 22 Nov 2026; tafs docs/chrome.md). The
// extension takes one just before each act, with password, card and code
// fields and anything it filled from Saved sign-ins masked, and never on a
// blocked or sensitive site or a payment frame. It sends it here, signed with
// its own key, and planque keeps it sealed under the account for 30 days.
//
//   - Only an act planque sent to that Chrome may have an image: sending one
//     records an allowance, and an upload is taken only against it, once.
//   - It's a WebP image, at most chromeShotMax bytes and chromeShotWidth
//     pixels wide; anything else is refused and nothing is kept.
//   - It's sealed (planque's Sealer, bound to the account and the step) and
//     kept with a TTL of 30 days, under the account, so an account's
//     erasure deletes it with the rest. It's never sent to a model.
//   - The owner reads it with their device's signature, by the step's id.
//
// Off unless PLANQUE_CHROME_SCREENSHOTS=on: both routes answer 404 and no
// allowance is recorded.

const (
	labelScreenshot = "tafs-chrome/screenshot/v1"
	// chromeShotKept is how long an image is kept.
	chromeShotKept = 30 * 24 * time.Hour
	// chromeShotMax is the largest image taken, after the extension's own
	// re-encoding (it aims at about 50 KB).
	chromeShotMax = 120 << 10
	// chromeShotBody is the largest upload body: the image in base64url and
	// the rest of the JSON.
	chromeShotBody = 200 << 10
	// chromeShotWidth is the widest image taken, in pixels.
	chromeShotWidth = 1280
	// chromeShotAfter is how long after a step's own expiry its image may
	// still come: the image is taken before the act, and sent after it.
	chromeShotAfter = 2 * time.Minute
	// chromeShotsPerAddress bounds uploads per address in chromePostWindow:
	// one per act, far more than a run makes.
	chromeShotsPerAddress = 300
)

func skChromeShot(step string) string { return "CHROMESHOT#" + step }

// pkChromeShotAllowed is the allowance an act's sending records: the step's
// account, grant and run, until the step lapses plus chromeShotAfter.
func pkChromeShotAllowed(chrome, step string) string {
	return "CHROMESHOTOK#" + chrome + "#" + step
}

type chromeShotAllowed struct {
	Tenant string `json:"tenant"`
	Grant  string `json:"grant"`
	Run    string `json:"run"`
}

// chromeShot is a kept image: what it belongs to, and its sealed bytes.
type chromeShot struct {
	Step   string    `json:"step"`
	Chrome string    `json:"chrome"`
	Grant  string    `json:"grant"`
	Run    string    `json:"run"`
	SHA256 string    `json:"sha256"`
	Bytes  int       `json:"bytes"`
	At     time.Time `json:"at"`
	Sealed []byte    `json:"sealed"`
}

func chromeShotBinding(tenant, step string) map[string]string {
	return map[string]string{"tenant": tenant, "kind": "chrome_screenshot", "step": step}
}

// allowChromeShot records that step, an act planque is sending, may have an
// image. Recorded only with screenshots on; a failure to record it leaves
// the step without an image, never without its act.
func (s *Server) allowChromeShot(ctx context.Context, tenant string, step chromeStep) {
	if !s.ChromeScreenshots || step.Op != "act" {
		return
	}
	ttl := time.UnixMilli(step.Exp).Add(chromeShotAfter).Unix() + 1
	_ = putDoc(ctx, s.KV, store.Item{PK: pkChromeShotAllowed(step.ChromeID, step.ID), SK: "-", TTL: ttl},
		chromeShotAllowed{Tenant: tenant, Grant: step.Grant, Run: step.Run}, store.Cond{NotExists: true})
}

// webpWidth is a WebP image's width in pixels, read from its first chunk
// (VP8, VP8L or VP8X), or false when the bytes aren't a WebP image.
func webpWidth(b []byte) (int, bool) {
	if len(b) < 30 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WEBP" {
		return 0, false
	}
	if int(binary.LittleEndian.Uint32(b[4:8]))+8 != len(b) {
		return 0, false
	}
	switch string(b[12:16]) {
	case "VP8 ":
		// A key frame's start code, then 14 bits of width.
		if b[23] != 0x9d || b[24] != 0x01 || b[25] != 0x2a {
			return 0, false
		}
		return int(binary.LittleEndian.Uint16(b[26:28]) & 0x3fff), true
	case "VP8L":
		if b[20] != 0x2f {
			return 0, false
		}
		return int(binary.LittleEndian.Uint16(b[21:23])&0x3fff) + 1, true
	case "VP8X":
		return int(b[24]) | int(b[25])<<8 | int(b[26])<<16 + 1, true
	}
	return 0, false
}

// postChromeShot is POST /v1/extension/steps/{step}/screenshot from the
// extension: {chrome_id, image (base64url of the WebP bytes), sig}, sig its
// key's signature over {chrome_id, step, sha256} of the image. Taken once,
// only against the allowance its act recorded; anything else is 404, 413 or
// 415, and nothing is kept.
func (s *Server) postChromeShot(w http.ResponseWriter, r *http.Request) {
	if !s.ChromeScreenshots || s.chromeOff(w) {
		if !s.ChromeScreenshots {
			writeErr(w, http.StatusNotFound, "screenshots aren't kept here")
		}
		return
	}
	if !s.chromeCounted(r, "shot#"+chromeAddress(r), chromePostWindow, chromeShotsPerAddress) {
		s.chromeTooMany(w)
		return
	}
	step := r.PathValue("step")
	// The image travels in base64url, a third over its bytes: its own cap,
	// not readBody's.
	raw, err := io.ReadAll(io.LimitReader(r.Body, chromeShotBody+1))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "the body is unreadable")
		return
	}
	if len(raw) > chromeShotBody {
		writeErr(w, http.StatusRequestEntityTooLarge, "the image is over 120 KB: nothing was kept")
		return
	}
	var in struct {
		ChromeID string `json:"chrome_id"`
		Image    string `json:"image"`
		Sig      string `json:"sig"`
	}
	if step == "" || json.Unmarshal(raw, &in) != nil || !chromeIDRe.MatchString(in.ChromeID) || in.Image == "" || in.Sig == "" {
		writeErr(w, http.StatusBadRequest, "send {chrome_id, image, sig}: the Chrome's id, the WebP image in base64url, and its key's signature over {chrome_id, step, sha256}")
		return
	}
	img, err := base64.RawURLEncoding.DecodeString(in.Image)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "image is the WebP bytes in base64url, unpadded")
		return
	}
	if len(img) > chromeShotMax {
		writeErr(w, http.StatusRequestEntityTooLarge, "the image is over 120 KB: nothing was kept")
		return
	}
	if width, ok := webpWidth(img); !ok || width < 1 || width > chromeShotWidth {
		writeErr(w, http.StatusUnsupportedMediaType, "the image is a WebP image at most 1280 pixels wide: nothing was kept")
		return
	}
	sum := sha256.Sum256(img)
	digest := hex.EncodeToString(sum[:])
	gone := func() { writeErr(w, http.StatusNotFound, "no step of this Chrome takes an image") }
	var ok chromeShotAllowed
	if _, err := getDoc(r.Context(), s.KV, pkChromeShotAllowed(in.ChromeID, step), "-", &ok); err != nil || ok.Tenant == "" {
		gone()
		return
	}
	var c Chrome
	if _, err := getDoc(r.Context(), s.KV, tenantPK(ok.Tenant), skChrome(in.ChromeID), &c); err != nil {
		gone()
		return
	}
	if !stepkey.Verify(c.Key, labelScreenshot, map[string]string{"chrome_id": in.ChromeID, "step": step, "sha256": digest}, in.Sig) {
		gone()
		return
	}
	// Once: the allowance goes as the image is taken.
	if _, found, err := s.KV.Delete(r.Context(), pkChromeShotAllowed(in.ChromeID, step), "-"); err != nil {
		s.fail(w, r, err)
		return
	} else if !found {
		gone()
		return
	}
	sealed, err := s.Sealer.Seal(r.Context(), img, chromeShotBinding(ok.Tenant, step))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	now := s.now().UTC()
	shot := chromeShot{Step: step, Chrome: in.ChromeID, Grant: ok.Grant, Run: ok.Run, SHA256: digest, Bytes: len(img), At: now, Sealed: sealed}
	item := store.Item{PK: tenantPK(ok.Tenant), SK: skChromeShot(step), TTL: now.Add(chromeShotKept).Unix()}
	if err := putDoc(r.Context(), s.KV, item, shot, store.Cond{NotExists: true}); err != nil {
		if errors.Is(err, store.ErrCondition) {
			gone()
			return
		}
		s.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// readChromeShot is GET /v1/accounts/{tenant}/chromes/screenshots/{step}, the
// owner's signed read of one step's image for the taf's Activity:
// {step, run, grant, taken_at, sha256, image} with the WebP bytes in
// base64url. One past its 30 days is gone, whether or not the store has
// deleted it yet.
func (s *Server) readChromeShot(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	if !s.ChromeScreenshots {
		writeErr(w, http.StatusNotFound, "screenshots aren't kept here")
		return
	}
	step := r.PathValue("step")
	var shot chromeShot
	if _, err := getDoc(r.Context(), s.KV, tenantPK(sr.tenant), skChromeShot(step), &shot); err != nil || s.now().Sub(shot.At) > chromeShotKept {
		if err != nil && !errors.Is(err, errNotFound) {
			s.fail(w, r, err)
			return
		}
		writeErr(w, http.StatusNotFound, "no image for that step")
		return
	}
	img, err := s.Sealer.Open(r.Context(), shot.Sealed, chromeShotBinding(sr.tenant, step))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"step": shot.Step, "run": shot.Run, "grant": shot.Grant,
		"taken_at": shot.At.Format(time.RFC3339), "sha256": shot.SHA256, "image": base64.RawURLEncoding.EncodeToString(img)})
}

// dropChromeShots deletes the account's kept images that match: a
// connection's when it ends, a Chrome's when it's unpaired. The images go
// sooner than their 30 days, as the pages a write kept do on a disconnect.
// A failure leaves the rest to the TTL.
func (s *Server) dropChromeShots(ctx context.Context, tenant string, match func(chromeShot) bool) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "CHROMESHOT#")
	if err != nil {
		return
	}
	for _, it := range items {
		var shot chromeShot
		if json.Unmarshal(it.Doc, &shot) != nil || !match(shot) {
			continue
		}
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
	}
}
