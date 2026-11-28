package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/fakes"
)

// act runs the shop's add-to-cart steps, the call that gives the files.
func (e *browserEnv) act(g *Grant) (string, map[string]any) {
	e.t.Helper()
	return e.call(g, "browse_act", map[string]any{"url": fakes.ShopItem, "read": e.b.ReadOf(fakes.ShopItem),
		"steps": []map[string]any{{"fill": "qty", "value": "2", "label": "Quantity"}, {"click": "add", "label": "Add to cart", "submits": true}}})
}

func pdfOf(n int) []byte { return append([]byte("%PDF-1.4\n"), []byte(strings.Repeat("0", n))...) }

// A file a step's page gave is kept sealed under the grant, in parts, and
// listed in the answer as data; refused ones are listed with why; the
// account's key reads it back whole, to save and never to run; it's gone
// after 30 days, and with the grant.
func TestDownloadsAreKeptInTheGrantsStore(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	e.s.BrowserDownloads = map[string]bool{tenant: true}
	big := pdfOf(600 << 10) // three parts
	sum := sha256.Sum256(big)
	e.b.Files = []browser.Download{
		{Name: "invoice-1001.pdf", Type: browser.TypePDF, Size: len(big), SHA256: hex.EncodeToString(sum[:]), Data: big},
		{Name: "tool.exe", Refused: "only a PDF, CSV, PNG or JPEG file is taken"},
		// The worker said PDF, but the bytes aren't one: checked again here.
		{Name: "lie.pdf", Type: browser.TypePDF, Data: []byte("MZ\x90\x00")},
	}
	out, sc := e.act(g)
	dl, _ := sc["downloads"].([]any)
	if code(sc) != "" || len(dl) != 3 || e.b.StartedWithDownloads != 1 {
		t.Fatalf("act: %s", out)
	}
	kept := dl[0].(map[string]any)
	id, _ := kept["file"].(string)
	if !strings.HasPrefix(id, "bf_") || kept["type"] != browser.TypePDF || kept["sha256"] != hex.EncodeToString(sum[:]) || kept["name"] != "invoice-1001.pdf" {
		t.Fatalf("kept: %v", kept)
	}
	if dl[1].(map[string]any)["refused"] == nil || dl[2].(map[string]any)["refused"] == nil || dl[2].(map[string]any)["file"] != nil {
		t.Fatalf("refused: %v", dl[1:])
	}
	// Sealed: no part holds the file's bytes in the clear.
	items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "BFILEPART#"+g.ID+"#")
	if len(items) != 3 {
		t.Fatalf("parts: %d", len(items))
	}
	for _, it := range items {
		if strings.Contains(string(it.Doc), "%PDF") || it.TTL == 0 {
			t.Fatalf("a part in the clear, or with no TTL: %d", it.TTL)
		}
	}
	// Read back with the account's key, whole, as a file to save.
	get := func(key string) (*http.Response, []byte) {
		req, _ := http.NewRequest("GET", e.url+"/v1/accounts/"+tenant+"/grants/"+g.ID+"/files/"+id, nil)
		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
		res, err := e.client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res, b
	}
	if res, _ := get(""); res.StatusCode == 200 {
		t.Fatal("read with no key")
	}
	res, body := get("k_setup")
	if res.StatusCode != 200 || string(body) != string(big) || res.Header.Get("Content-Type") != browser.TypePDF ||
		!strings.HasPrefix(res.Header.Get("Content-Disposition"), "attachment") || res.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("read back: %d %v (%d bytes)", res.StatusCode, res.Header, len(body))
	}
	// Gone after 30 days.
	later := time.Now().Add(browserFileKept + time.Hour)
	e.s.Now = func() time.Time { return later }
	if res, _ := get("k_setup"); res.StatusCode != 404 {
		t.Fatalf("after 30 days: %d", res.StatusCode)
	}
	e.s.Now = nil
	// Gone with the grant.
	e.s.dropBrowser(context.Background(), g)
	for _, prefix := range []string{"BFILE#" + g.ID + "#", "BFILEPART#" + g.ID + "#"} {
		if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), prefix); len(items) != 0 {
			t.Fatalf("left after the grant: %d under %s", len(items), prefix)
		}
	}
}

// Off for the account, a session isn't started with downloads, and no file
// is kept or listed.
func TestDownloadsAreOffUnlessTheAccountIsListed(t *testing.T) {
	e := newBrowserEnv(t)
	g := e.connect()
	e.b.Files = []browser.Download{{Name: "invoice.pdf", Type: browser.TypePDF, Data: pdfOf(10)}}
	out, sc := e.act(g)
	if code(sc) != "" || sc["downloads"] != nil || e.b.StartedWithDownloads != 0 {
		t.Fatalf("off: %s", out)
	}
	if items, _ := e.s.KV.Query(context.Background(), tenantPK(tenant), "BFILE"); len(items) != 0 {
		t.Fatalf("kept while off: %d", len(items))
	}
}
