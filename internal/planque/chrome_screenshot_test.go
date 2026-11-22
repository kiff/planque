package planque

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/kiff/planque/internal/stepkey"
)

// webp is a minimal WebP image header (VP8X) width pixels wide, padded to n
// bytes: enough for planque's checks, which read the size and the width only.
func webp(width, n int) []byte {
	b := make([]byte, n)
	copy(b, "RIFF")
	binary.LittleEndian.PutUint32(b[4:], uint32(n-8))
	copy(b[8:], "WEBPVP8X")
	binary.LittleEndian.PutUint32(b[16:], 10)
	w := width - 1
	b[24], b[25], b[26] = byte(w), byte(w>>8), byte(w>>16)
	b[27] = 0 // height 1
	return b
}

// The owner's decision (22 Nov 2026): an acting step in their own Chrome is
// kept as an image, for the record. Only for an act planque sent to that
// Chrome, signed by that Chrome's key, once; sealed under the account for
// 30 days; read by the owner's device; never with the switch off.
func TestAnActsImageIsKeptOnceAndReadByTheOwner(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeScreenshots = true
	c.r.connect(t, "c1")
	ctx := context.Background()
	now := time.Now()
	c.s.Now = func() time.Time { return now }
	// The phone signs on planque's clock, a second apart, so no read is a replay.
	tick := 0
	c.phone.Now = func() time.Time { tick++; return now.Add(time.Duration(tick) * time.Second) }
	send := func(id, op string) {
		t.Helper()
		if err := c.s.sendChromeStep(ctx, tenant, chromeStep{ID: id, ChromeID: c.r.chrome, Run: "run_" + id, Grant: "grt_x", Sites: []string{"shop.example"}, Op: op}); err != nil {
			t.Fatal(err)
		}
	}
	upload := func(step string, img []byte, key stepkey.Local, signedStep string) int {
		t.Helper()
		sum := sha256.Sum256(img)
		sig, err := stepkey.SignBody(ctx, key, labelScreenshot, map[string]string{"chrome_id": c.r.chrome, "step": signedStep, "sha256": hex.EncodeToString(sum[:])})
		if err != nil {
			t.Fatal(err)
		}
		st, _ := c.do("POST", "/v1/extension/steps/"+step+"/screenshot", map[string]any{"chrome_id": c.r.chrome,
			"image": base64.RawURLEncoding.EncodeToString(img), "sig": sig}, nil, "")
		return st
	}
	read := func(step string) (int, map[string]any) {
		return c.do("GET", "/v1/accounts/"+tenant+"/chromes/screenshots/"+step, nil, c.phone, "")
	}

	send("stp_act1", "act")
	img := webp(1280, 50<<10)
	if st := upload("stp_act1", img, c.r.ext, "stp_act1"); st != http.StatusNoContent {
		t.Fatalf("the act's image: %d", st)
	}
	if st := upload("stp_act1", img, c.r.ext, "stp_act1"); st != http.StatusNotFound {
		t.Errorf("the same step again: %d, want 404 (once)", st)
	}
	st, out := read("stp_act1")
	if st != 200 || out["image"] != base64.RawURLEncoding.EncodeToString(img) || out["run"] != "run_stp_act1" || out["step"] != "stp_act1" {
		t.Fatalf("the owner's read: %d %v", st, out)
	}
	// Kept sealed: the stored record never holds the image's bytes.
	var shot chromeShot
	if _, err := getDoc(ctx, c.s.KV, tenantPK(tenant), skChromeShot("stp_act1"), &shot); err != nil || len(shot.Sealed) == 0 || string(shot.Sealed[:4]) == "RIFF" {
		t.Errorf("the stored record: %v", err)
	}
	// Unsigned, the read is refused.
	if st, _ := c.do("GET", "/v1/accounts/"+tenant+"/chromes/screenshots/stp_act1", nil, nil, ""); st != 401 {
		t.Errorf("an unsigned read: %d", st)
	}

	// Refused, and nothing kept: a step that only reads or opens; a step
	// planque never sent; another key; a signature over another step; not
	// WebP, too wide or too big.
	send("stp_read1", "read")
	send("stp_act2", "act")
	send("stp_act3", "act")
	send("stp_act4", "act")
	other := stepkey.Local{Key: c.r.planque.Key}
	for name, got := range map[string]int{
		"a read":                  upload("stp_read1", img, c.r.ext, "stp_read1"),
		"never sent":              upload("stp_nope", img, c.r.ext, "stp_nope"),
		"another key":             upload("stp_act2", img, other, "stp_act2"),
		"signed for another step": upload("stp_act2", img, c.r.ext, "stp_act3"),
	} {
		if got != http.StatusNotFound {
			t.Errorf("%s: %d, want 404", name, got)
		}
	}
	if st := upload("stp_act2", []byte("not an image at all, not a WebP image, plainly text bytes"), c.r.ext, "stp_act2"); st != http.StatusUnsupportedMediaType {
		t.Errorf("not WebP: %d", st)
	}
	if st := upload("stp_act3", webp(1281, 4<<10), c.r.ext, "stp_act3"); st != http.StatusUnsupportedMediaType {
		t.Errorf("1281 px wide: %d", st)
	}
	if st := upload("stp_act4", webp(800, 121<<10), c.r.ext, "stp_act4"); st != http.StatusRequestEntityTooLarge {
		t.Errorf("over 120 KB: %d", st)
	}
	for _, step := range []string{"stp_read1", "stp_act2", "stp_act3", "stp_act4"} {
		if st, _ := read(step); st != 404 {
			t.Errorf("%s was kept: %d", step, st)
		}
	}
	// A refused upload leaves the act's allowance: the right image still goes.
	if st := upload("stp_act2", webp(640, 30<<10), c.r.ext, "stp_act2"); st != http.StatusNoContent {
		t.Errorf("the right image after a refused one: %d", st)
	}

	// Past its 30 days, it's gone, whether or not the store has deleted it.
	now = now.Add(chromeShotKept + time.Minute)
	if st, _ := read("stp_act1"); st != 404 {
		t.Errorf("after 30 days: %d", st)
	}

	// Off: no allowance is recorded, and both routes answer 404.
	now = time.Now()
	c.s.ChromeScreenshots = false
	send("stp_act5", "act")
	c.s.ChromeScreenshots = true
	if st := upload("stp_act5", img, c.r.ext, "stp_act5"); st != http.StatusNotFound {
		t.Errorf("an act sent while off: %d", st)
	}
	c.s.ChromeScreenshots = false
	send("stp_act6", "act")
	if st := upload("stp_act6", img, c.r.ext, "stp_act6"); st != http.StatusNotFound {
		t.Errorf("upload while off: %d", st)
	}
	if st, _ := read("stp_act2"); st != 404 {
		t.Errorf("read while off: %d", st)
	}
}

func TestWebPWidthReadsEachKindOfFirstChunk(t *testing.T) {
	if w, ok := webpWidth(webp(1280, 64)); !ok || w != 1280 {
		t.Errorf("VP8X: %d %v", w, ok)
	}
	vp8 := webp(1, 64)
	copy(vp8[12:], "VP8 ")
	vp8[23], vp8[24], vp8[25] = 0x9d, 0x01, 0x2a
	binary.LittleEndian.PutUint16(vp8[26:], 640)
	if w, ok := webpWidth(vp8); !ok || w != 640 {
		t.Errorf("VP8: %d %v", w, ok)
	}
	vp8l := webp(1, 64)
	copy(vp8l[12:], "VP8L")
	vp8l[20] = 0x2f
	binary.LittleEndian.PutUint16(vp8l[21:], 1023)
	if w, ok := webpWidth(vp8l); !ok || w != 1024 {
		t.Errorf("VP8L: %d %v", w, ok)
	}
	bad := webp(100, 64)
	binary.LittleEndian.PutUint32(bad[4:], 1000) // a size that isn't the file's
	if _, ok := webpWidth(bad); ok {
		t.Error("a RIFF size that isn't the file's was read")
	}
	if _, ok := webpWidth([]byte("GIF89a")); ok {
		t.Error("not WebP was read")
	}
}

// An image goes sooner than its 30 days when what it belongs to ends: a
// Chrome's when it's unpaired, a connection's when it's disconnected (as the
// pages a write kept do), and nothing else's.
func TestImagesGoWithTheirChromeAndTheirConnection(t *testing.T) {
	c := newChromeRig(t)
	c.s.ChromeScreenshots = true
	c.r.connect(t, "c1")
	ctx := context.Background()
	now := time.Now()
	c.s.Now = func() time.Time { return now }
	tick := 0
	c.phone.Now = func() time.Time { tick++; return now.Add(time.Duration(tick) * time.Second) }
	keep := func(id, grant string) {
		t.Helper()
		if err := c.s.sendChromeStep(ctx, tenant, chromeStep{ID: id, ChromeID: c.r.chrome, Run: "run_" + id, Grant: grant, Sites: []string{"shop.example"}, Op: "act"}); err != nil {
			t.Fatal(err)
		}
		img := webp(800, 20<<10)
		sum := sha256.Sum256(img)
		sig, _ := stepkey.SignBody(ctx, c.r.ext, labelScreenshot, map[string]string{"chrome_id": c.r.chrome, "step": id, "sha256": hex.EncodeToString(sum[:])})
		if st, out := c.do("POST", "/v1/extension/steps/"+id+"/screenshot", map[string]any{"chrome_id": c.r.chrome,
			"image": base64.RawURLEncoding.EncodeToString(img), "sig": sig}, nil, ""); st != http.StatusNoContent {
			t.Fatalf("keep %s: %d %v", id, st, out)
		}
	}
	kept := func(id string) bool {
		st, _ := c.do("GET", "/v1/accounts/"+tenant+"/chromes/screenshots/"+id, nil, c.phone, "")
		return st == 200
	}

	// A connection's images go when it's disconnected; another's stay.
	g := c.grant(t, browseOpen, browseAct)
	keep("stp_mine", g.ID)
	keep("stp_other", "grt_another")
	if st, out := c.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, c.phone, ""); st != 200 {
		t.Fatalf("disconnect: %d %v", st, out)
	}
	if kept("stp_mine") {
		t.Error("the disconnected connection's image is still kept")
	}
	if !kept("stp_other") {
		t.Error("another connection's image went with it")
	}

	// A Chrome's images go when it's unpaired.
	if st, out := c.do("DELETE", "/v1/accounts/"+tenant+"/chromes/"+c.r.chrome, nil, c.phone, ""); st != 200 && st != 204 {
		t.Fatalf("unpair: %d %v", st, out)
	}
	if kept("stp_other") {
		t.Error("the unpaired Chrome's image is still kept")
	}
}
