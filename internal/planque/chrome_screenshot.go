package planque

import (
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
)

// A taf's acting steps in the owner's own Chrome are kept as images, for the
// record (the owner's decision; tafs docs/chrome.md). The
// extension takes one just before each act, with password, card and code
// fields and anything it filled from Saved sign-ins masked, and never on a
// blocked or sensitive site or a payment frame. It sends it here, signed with
// its own key, and planque keeps it sealed under the account for 30 days.
//
//   - Only an act planque sent to that Chrome may have an image: sending one
//     records an allowance, and an upload is taken only against it, once.
//   - It's a WebP image, at most chromeShotMax bytes, chromeShotWidth
//     pixels wide and chromeShotHeight high; anything else is refused and nothing is kept.
//   - It's sealed (planque's Sealer, bound to the account and the step) and
//     kept with a TTL of 30 days, under the account, so an account's
//     erasure deletes it with the rest. It's never sent to a model.
//   - The owner reads it with their device's signature, by the step's id.
//
// Off unless PLANQUE_CHROME_SCREENSHOTS=on: both routes answer 404 and no
// allowance is recorded. It's one switch for every account (no per-account
// choice), and it goes on in production only on the owner's yes.

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
	// chromeShotHeight is the tallest image taken: a page is captured as the
	// tab shows it, and a phone never decodes a canvas taller than this.
	chromeShotHeight = 8000
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

// webpSize is a WebP image's width and height in pixels, read from its
// first chunk (VP8, VP8L or VP8X), or false when the bytes aren't a WebP
// image.
func webpSize(b []byte) (int, int, bool) {
	if len(b) < 30 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WEBP" {
		return 0, 0, false
	}
	if int(binary.LittleEndian.Uint32(b[4:8]))+8 != len(b) {
		return 0, 0, false
	}
	switch string(b[12:16]) {
	case "VP8 ":
		// A key frame's start code, then 14 bits of width and of height.
		if b[23] != 0x9d || b[24] != 0x01 || b[25] != 0x2a {
			return 0, 0, false
		}
		return int(binary.LittleEndian.Uint16(b[26:28]) & 0x3fff), int(binary.LittleEndian.Uint16(b[28:30]) & 0x3fff), true
	case "VP8L":
		// The signature, then 14 bits of width-1 and 14 of height-1.
		if b[20] != 0x2f {
			return 0, 0, false
		}
		v := binary.LittleEndian.Uint32(b[21:25])
		return int(v&0x3fff) + 1, int(v>>14&0x3fff) + 1, true
	case "VP8X":
		// The canvas: 24 bits of width-1, then 24 of height-1.
		return int(b[24]) | int(b[25])<<8 | int(b[26])<<16 + 1, int(b[27]) | int(b[28])<<8 | int(b[29])<<16 + 1, true
	}
	return 0, 0, false
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
	if width, height, ok := webpSize(img); !ok || width < 1 || width > chromeShotWidth || height < 1 || height > chromeShotHeight {
		writeErr(w, http.StatusUnsupportedMediaType, "the image is a WebP image at most 1280 pixels wide and 8000 high: nothing was kept")
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
	if ok, err := s.shotReadable(r.Context(), sr.tenant, shot); err != nil {
		s.fail(w, r, err)
		return
	} else if !ok {
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

// shotsOf names whose images a cleanup deletes: a Chrome's (unpaired, or
// replaced by a re-pairing, its carried connections' images with it), or a
// connection's (disconnected).
type shotsOf struct {
	Chrome string `json:"chrome,omitempty"`
	Grant  string `json:"grant,omitempty"`
}

func (o shotsOf) match(shot chromeShot) bool {
	return o.Chrome != "" && shot.Chrome == o.Chrome || o.Grant != "" && shot.Grant == o.Grant
}

// pkChromeShotClean holds the cleanups that didn't finish, one item each, for
// the hourly sweep to finish (sweepChromeShots).
const pkChromeShotClean = "CHROMESHOTCLEAN"

// chromeShotClean is a cleanup left for the sweep: the account and whose images.
type chromeShotClean struct {
	Tenant string    `json:"tenant"`
	Of     shotsOf   `json:"of"`
	At     time.Time `json:"at"`
}

func skChromeShotClean(tenant string, of shotsOf) string {
	return tenant + "#" + of.Chrome + "#" + of.Grant
}

// errShotsLeft says a cleanup left images it couldn't delete: none can be read
// meanwhile (shotReadable), and the sweep tries again.
var errShotsLeft = errors.New("some images could not be deleted yet")

// dropChromeShots deletes the account's kept images that are of what ended,
// sooner than their 30 days, as the pages a write kept go on a disconnect.
// A read of the images, or a delete, that fails leaves the cleanup for the
// hourly sweep, and says so (errShotsLeft). They can't be read meanwhile:
// their Chrome or connection is gone, which shotReadable checks.
func (s *Server) dropChromeShots(ctx context.Context, tenant string, of shotsOf) error {
	err := s.deleteChromeShots(ctx, tenant, of)
	if err == nil {
		return nil
	}
	s.log().WarnContext(ctx, "chrome screenshots not all deleted", slog.String("tenant", tenant), slog.String("err", err.Error()))
	now := s.now()
	it := store.Item{PK: pkChromeShotClean, SK: skChromeShotClean(tenant, of), TTL: now.Add(chromeShotKept + 24*time.Hour).Unix()}
	if perr := putDoc(ctx, s.KV, it, chromeShotClean{Tenant: tenant, Of: of, At: now}, store.Cond{}); perr != nil {
		s.log().ErrorContext(ctx, "chrome screenshots cleanup not recorded for the sweep", slog.String("tenant", tenant), slog.String("err", perr.Error()))
	}
	return errShotsLeft
}

// deleteChromeShots deletes them once, and says whether any read or delete failed.
func (s *Server) deleteChromeShots(ctx context.Context, tenant string, of shotsOf) error {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "CHROMESHOT#")
	if err != nil {
		return err
	}
	var first error
	for _, it := range items {
		var shot chromeShot
		if json.Unmarshal(it.Doc, &shot) != nil || !of.match(shot) {
			continue
		}
		if _, _, err := s.KV.Delete(ctx, it.PK, it.SK); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// sweepChromeShots finishes the cleanups that didn't: each is tried again, and
// goes once its images are deleted. It returns how many finished.
func (s *Server) sweepChromeShots(ctx context.Context) (int, error) {
	items, err := s.KV.Query(ctx, pkChromeShotClean, "")
	if err != nil {
		return 0, err
	}
	done := 0
	var first error
	for _, it := range items {
		var c chromeShotClean
		if json.Unmarshal(it.Doc, &c) != nil || c.Tenant == "" {
			_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
			continue
		}
		if err := s.deleteChromeShots(ctx, c.Tenant, c.Of); err != nil {
			first = cmpErr(first, err)
			continue
		}
		if _, _, err := s.KV.Delete(ctx, it.PK, it.SK); err != nil {
			first = cmpErr(first, err)
			continue
		}
		done++
	}
	return done, first
}

// shotReadable says whether a kept image may still be read: its Chrome is
// still paired to the account, and its connection is still there and not
// disconnected. So an unpair, a re-pairing that replaced the Chrome, or a
// disconnect cuts access at once, whether or not its images are deleted yet,
// and an image a restore from backup brings back stays unread.
func (s *Server) shotReadable(ctx context.Context, tenant string, shot chromeShot) (bool, error) {
	var c Chrome
	if _, err := getDoc(ctx, s.KV, tenantPK(tenant), skChrome(shot.Chrome), &c); errors.Is(err, errNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	g, err := s.loadGrant(ctx, tenant, shot.Grant)
	if errors.Is(err, errNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return g.Status != StatusDisconnected, nil
}

// chromeShotsListed is the most images one operation's list names: the
// newest, as an operation may act many times.
const chromeShotsListed = 20

// listChromeShots is GET /v1/accounts/{tenant}/chromes/screenshots?op=…,
// the owner's signed list of one operation's images, for a taf's Activity
// row by its proposal_id (KIFF's operation, which planque runs a Chrome
// call as): {shots: [{step, taken_at, sha256}]}, oldest first, the newest
// chromeShotsListed at most, never the bytes (each comes from
// readChromeShot, with its checks). Only that operation's own images, kept
// and within their 30 days; none, or screenshots off, is an empty list, so
// a row with none just shows no link.
func (s *Server) listChromeShots(w http.ResponseWriter, r *http.Request, sr signedRequest) {
	ops := r.URL.Query()["op"]
	if len(ops) != 1 || ops[0] == "" || len(ops[0]) > 512 {
		writeErr(w, http.StatusBadRequest, "send op, the Activity row's proposal_id, once")
		return
	}
	op := ops[0]
	type listed struct {
		Step    string `json:"step"`
		TakenAt string `json:"taken_at"`
		SHA256  string `json:"sha256"`
	}
	out := []listed{}
	if !s.ChromeScreenshots {
		writeJSON(w, http.StatusOK, map[string]any{"shots": out})
		return
	}
	// A step's id is the operation and its number, or a hash of the
	// operation when that is too long (chromeStepID): both prefixes are
	// read, and each image is matched by the operation it was kept for.
	h := sha256.Sum256([]byte(op))
	var shots []chromeShot
	for _, prefix := range []string{skChromeShot(op + "."), skChromeShot("op_" + hex.EncodeToString(h[:])[:40] + ".")} {
		items, err := s.KV.Query(r.Context(), tenantPK(sr.tenant), prefix)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		for _, it := range items {
			var shot chromeShot
			if json.Unmarshal(it.Doc, &shot) != nil || shot.Run != op || s.now().Sub(shot.At) > chromeShotKept ||
				slices.ContainsFunc(shots, func(k chromeShot) bool { return k.Step == shot.Step }) {
				continue
			}
			if ok, err := s.shotReadable(r.Context(), sr.tenant, shot); err != nil {
				s.fail(w, r, err)
				return
			} else if !ok {
				continue
			}
			shots = append(shots, shot)
		}
	}
	slices.SortFunc(shots, func(a, b chromeShot) int { return cmp.Or(a.At.Compare(b.At), cmp.Compare(a.Step, b.Step)) })
	if len(shots) > chromeShotsListed {
		shots = shots[len(shots)-chromeShotsListed:]
	}
	for _, shot := range shots {
		out = append(out, listed{Step: shot.Step, TakenAt: shot.At.Format(time.RFC3339), SHA256: shot.SHA256})
	}
	writeJSON(w, http.StatusOK, map[string]any{"shots": out})
}
