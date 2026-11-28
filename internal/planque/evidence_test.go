package planque

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// The evidence pictures, kept sealed: a write's picture is kept under the
// grant with the operation, read back as a PNG with the account's key, and
// gone after 30 days and with the grant; a picture that isn't a PNG isn't
// kept.
func TestEvidencePicturesAreKeptSealed(t *testing.T) {
	e := newBrowserEnv(t)
	e.s.Evidence = e.s.SealedEvidence()
	var pic bytes.Buffer
	_ = png.Encode(&pic, image.NewRGBA(image.Rect(0, 0, 4, 4)))
	e.b.ScreenshotPNG = pic.Bytes()
	g := e.connect()
	out, sc := e.act(g)
	ev, _ := sc["evidence"].(map[string]any)
	ref, _ := ev["screenshot"].(string)
	if code(sc) != "" || !strings.HasPrefix(ref, "ev_") {
		t.Fatalf("act: %s", out)
	}
	items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "BFILEPART#"+g.ID+"#")
	for _, it := range items {
		if bytes.Contains(it.Doc, []byte("PNG")) || it.TTL == 0 {
			t.Fatal("a picture's part in the clear, or with no TTL")
		}
	}
	get := func(path string) (int, []byte, http.Header) {
		req, _ := http.NewRequest("GET", e.url+"/v1/accounts/"+tenant+"/grants/"+g.ID+path, nil)
		req.Header.Set("Authorization", "Bearer k_setup")
		res, err := e.client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, b, res.Header
	}
	st, body, hdr := get("/evidence/" + ref)
	if st != 200 || !bytes.Equal(body, pic.Bytes()) || hdr.Get("Content-Type") != browser.TypePNG || hdr.Get("Cache-Control") != "no-store" {
		t.Fatalf("read back: %d %v", st, hdr)
	}
	// A picture isn't a download: the files route doesn't serve it.
	if st, _, _ := get("/files/" + ref); st != 404 {
		t.Fatalf("a picture through the files route: %d", st)
	}
	later := time.Now().Add(browserFileKept + time.Hour)
	e.s.Now = func() time.Time { return later }
	if st, _, _ := get("/evidence/" + ref); st != 404 {
		t.Fatalf("after 30 days: %d", st)
	}
	e.s.Now = nil
	e.s.dropBrowser(context.Background(), g)
	if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "BFILE"); len(items) != 0 {
		t.Fatalf("left after the grant: %d", len(items))
	}
	if _, err := e.s.SealedEvidence().Put(context.Background(), tenant, g.ID, "op_x", []byte("\x89PNG not one")); err == nil {
		t.Fatal("a picture that isn't a PNG was kept")
	}
}
