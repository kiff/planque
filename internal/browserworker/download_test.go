package browserworker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// dlSite serves files as attachments, and a page linking each; dlTLSSite
// the same over TLS, for a browser reaching it as a named site.
func dlSite(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(dlHandler())
	t.Cleanup(srv.Close)
	return srv
}

func dlTLSSite(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewTLSServer(dlHandler())
	t.Cleanup(srv.Close)
	return srv
}

func dlHandler() http.Handler {
	var pic bytes.Buffer
	_ = png.Encode(&pic, image.NewRGBA(image.Rect(0, 0, 2, 2)))
	files := map[string][]byte{
		"invoice.pdf": []byte("%PDF-1.4\n1 0 obj<< /Type /Catalog >>endobj\ntrailer<<>>\n%%EOF\n"),
		"active.pdf":  []byte("%PDF-1.4\n1 0 obj<< /OpenAction << /S /JavaScript /JS (app.alert(1)) >> >>endobj\n%%EOF\n"),
		"export.csv":  []byte("date,amount\n2026-11-01,12.50\n"),
		"tool.exe":    append([]byte("MZ\x90\x00\x03"), bytes.Repeat([]byte{0}, 64)...),
		"pic.png":     pic.Bytes(),
		"big.pdf":     append([]byte("%PDF-1.4\n"), bytes.Repeat([]byte("0"), browser.MaxDownload+1024)...),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		var links strings.Builder
		for name := range files {
			fmt.Fprintf(&links, `<p><a href="/dl/%s">%s</a></p>`, name, name)
		}
		fmt.Fprintf(&links, `<p><a href="%s">offsite</a></p>`, r.URL.Query().Get("off"))
		fmt.Fprint(w, "<!doctype html><html><head><title>Files</title></head><body>"+links.String()+"</body></html>")
	})
	mux.HandleFunc("/dl/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/dl/")
		b, ok := files[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
		_, _ = w.Write(b)
	})
	return mux
}

// clickLink clicks the link with that label and hands over the downloads.
func clickLink(t *testing.T, sess browser.Session, page browser.Page, label string) []browser.Download {
	t.Helper()
	ref := mustLink(t, page, label)
	_, _ = sess.Do(context.Background(), browser.Step{Click: ref})
	d, ok := sess.(browser.Downloader)
	if !ok {
		t.Fatal("the session takes no downloads")
	}
	list, err := d.Downloads(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return list
}

// A download is taken only as a PDF, CSV, PNG or JPEG, typed by its bytes,
// with nothing that runs or opens, under 5 MB; anything else is refused
// with why, and its file is gone after.
func TestDownloadsAreTypedCheckedAndCapped(t *testing.T) {
	w, srv := worker(t), dlSite(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, browser.Start{Grant: "grt_test", Downloads: true})
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	page, err := sess.Open(ctx, srv.URL+"/files")
	if err != nil {
		t.Fatal(err)
	}
	for label, want := range map[string]string{"invoice.pdf": browser.TypePDF, "export.csv": browser.TypeCSV, "pic.png": browser.TypePNG} {
		list := clickLink(t, sess, page, label)
		if len(list) != 1 || list[0].Type != want || list[0].Refused != "" || len(list[0].Data) == 0 || list[0].SHA256 == "" || list[0].Name != label {
			t.Errorf("%s: %+v", label, list)
		}
	}
	for label, why := range map[string]string{"active.pdf": "script", "tool.exe": "only a PDF", "big.pdf": "larger than 5 MB"} {
		list := clickLink(t, sess, page, label)
		if len(list) != 1 || !strings.Contains(list[0].Refused, why) || len(list[0].Data) != 0 {
			t.Errorf("%s: %+v", label, list)
		}
	}
	dir := filepath.Join(sess.(*session).c.dir, "downloads")
	if left, _ := os.ReadDir(dir); len(left) != 0 {
		t.Errorf("files left after they were handed over: %v", left)
	}
}

// Off, a download is denied: nothing is saved and none is handed over.
func TestDownloadsAreDeniedUnlessOn(t *testing.T) {
	w, srv := worker(t), dlSite(t)
	ctx := context.Background()
	sess, err := w.Start(ctx, browser.Start{Grant: "grt_test"})
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	page, err := sess.Open(ctx, srv.URL+"/files")
	if err != nil {
		t.Fatal(err)
	}
	if list := clickLink(t, sess, page, "invoice.pdf"); len(list) != 0 {
		t.Errorf("handed over with downloads off: %+v", list)
	}
	time.Sleep(300 * time.Millisecond)
	var found []string
	_ = filepath.Walk(sess.(*session).c.dir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.Contains(p, "invoice") {
			found = append(found, p)
		}
		return nil
	})
	if len(found) != 0 {
		t.Errorf("saved with downloads off: %v", found)
	}
}

// A file from a site the session doesn't reach never arrives: the egress
// refuses it, and nothing is taken; a file from the listed site, in the
// same session, is.
func TestADownloadOffTheSitesNeverArrives(t *testing.T) {
	if _, err := chromePath(""); err != nil {
		t.Skip(err)
	}
	site, other := dlTLSSite(t), dlTLSSite(t)
	w := &Local{Options: Options{AllowLocal: true, NavigateWait: 20 * time.Second,
		Map: map[string]string{"files.example:443": site.Listener.Addr().String(), "other.example:443": other.Listener.Addr().String()},
		Resolve: func(context.Context, string) ([]netip.Addr, error) { return nil, errors.New("no resolver in tests") }}}
	ctx := context.Background()
	sess, err := w.Start(ctx, browser.Start{Grant: "grt_test", Downloads: true, Sites: []string{"files.example"}})
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	page, err := sess.Open(ctx, "https://files.example/files?off=https://other.example/dl/invoice.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if list := clickLink(t, sess, page, "invoice.pdf"); len(list) != 1 || list[0].Type != browser.TypePDF {
		t.Errorf("the listed site's file in the same session: %+v", list)
	}
	before := sess.(*session).proxy.Refused()
	for _, d := range clickLink(t, sess, page, "offsite") {
		if d.Refused == "" || len(d.Data) != 0 {
			t.Errorf("an off-site file was taken: %+v", d)
		}
	}
	if sess.(*session).proxy.Refused() <= before {
		t.Error("the egress refused nothing: the off-site file wasn't asked for")
	}
}
