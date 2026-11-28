package planque

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/store"
)

// Downloads into the grant's store (browser batch 1): a file a page gave in
// the hosted browser, on the grant's sites, typed by its bytes and checked
// by the worker (browser.TypeOf), is kept sealed under the grant for
// browserFileKept, in chunks under the item size, and deleted sooner with
// the grant or the account. The taf's answer lists what was taken, as data;
// the owner's apps read a file with the account's key. Only for the
// accounts BrowserDownloads lists (PLANQUE_BROWSER_DOWNLOADS).

// browserFileKept is how long a downloaded file is kept.
const browserFileKept = 30 * 24 * time.Hour

// browserFileChunk is the most of a file one item holds, sealed: under the
// store's item limit with room for the seal and the keys.
const browserFileChunk = 256 << 10

// browserFile is a downloaded file's record: what it is, not its bytes.
type browserFile struct {
	ID     string    `json:"id"`
	Grant  string    `json:"grant"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Size   int       `json:"size"`
	SHA256 string    `json:"sha256"`
	Chunks int       `json:"chunks"`
	Site   string    `json:"site,omitempty"`
	At     time.Time `json:"at"`
}

type browserFileChunkRecord struct {
	Sealed []byte `json:"sealed"`
}

func skBrowserFile(grant, id string) string { return "BFILE#" + grant + "#" + id }
func skBrowserFileChunk(grant, id string, n int) string {
	return fmt.Sprintf("BFILEPART#%s#%s#%04d", grant, id, n)
}

func browserFileBinding(g *Grant, id string, n int) map[string]string {
	return map[string]string{"tenant": g.Tenant, "grant": g.ID, "kind": "browser_file", "file": id, "chunk": strconv.Itoa(n)}
}

// downloadsOn says the account's hosted browser takes downloads.
func (s *Server) downloadsOn(tenant string) bool { return s.BrowserDownloads[tenant] }

// keepDownloads keeps the files a session's pages gave, and answers each as
// the taf reads it: the taken ones with their id, type, size and digest, the
// refused ones with why. A file is checked again here (its type from its
// bytes, the cap, the digest), whatever the worker said.
func (s *Server) keepDownloads(ctx context.Context, g *Grant, site string, list []browser.Download) []map[string]any {
	out := []map[string]any{}
	for _, d := range list {
		name := browser.DownloadName(d.Name, "")
		if d.Refused != "" {
			out = append(out, map[string]any{"name": name, "refused": d.Refused})
			continue
		}
		typ, err := browser.TypeOf(d.Data)
		sum := sha256.Sum256(d.Data)
		switch {
		case len(d.Data) > browser.MaxDownload:
			out = append(out, map[string]any{"name": name, "refused": "larger than 5 MB"})
			continue
		case err != nil:
			out = append(out, map[string]any{"name": name, "refused": err.Error()})
			continue
		case d.SHA256 != "" && d.SHA256 != hex.EncodeToString(sum[:]):
			out = append(out, map[string]any{"name": name, "refused": "its digest isn't the worker's"})
			continue
		}
		f := browserFile{ID: "bf_" + rand.Text(), Grant: g.ID, Name: browser.DownloadName(d.Name, typ), Type: typ, Size: len(d.Data),
			SHA256: hex.EncodeToString(sum[:]), Site: site, At: s.now().UTC()}
		if err := s.putBrowserFile(ctx, g, &f, d.Data); err != nil {
			s.log().WarnContext(ctx, "a download wasn't kept", slog.String("grant", g.ID), slog.String("err", err.Error()))
			out = append(out, map[string]any{"name": f.Name, "refused": "it couldn't be kept"})
			continue
		}
		g.event(s.now(), fmt.Sprintf("a file was downloaded at %s and kept for 30 days: %s (%s, %d bytes)", site, f.Name, f.Type, f.Size))
		kept := map[string]any{"file": f.ID, "name": f.Name, "type": f.Type, "size": f.Size, "sha256": f.SHA256}
		// A cell a spreadsheet would run: said, for whoever opens it.
		if typ == browser.TypeCSV && browser.CSVFormulas(d.Data) {
			kept["formulas"] = true
		}
		out = append(out, kept)
	}
	return out
}

// putBrowserFile seals data in chunks under the grant, then its record;
// nothing is left of it if any part fails.
func (s *Server) putBrowserFile(ctx context.Context, g *Grant, f *browserFile, data []byte) error {
	ttl := f.At.Add(browserFileKept).Unix()
	f.Chunks = (len(data) + browserFileChunk - 1) / browserFileChunk
	for n := 0; n < f.Chunks; n++ {
		part := data[n*browserFileChunk : min((n+1)*browserFileChunk, len(data))]
		sealed, err := s.Sealer.Seal(ctx, part, browserFileBinding(g, f.ID, n))
		if err == nil {
			err = putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skBrowserFileChunk(g.ID, f.ID, n), Ver: 1, TTL: ttl}, browserFileChunkRecord{Sealed: sealed}, store.Cond{NotExists: true})
		}
		if err != nil {
			s.dropBrowserFile(ctx, g.Tenant, g.ID, f.ID, n)
			return err
		}
	}
	if err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(g.Tenant), SK: skBrowserFile(g.ID, f.ID), Ver: 1, TTL: ttl}, f, store.Cond{NotExists: true}); err != nil {
		s.dropBrowserFile(ctx, g.Tenant, g.ID, f.ID, f.Chunks)
		return err
	}
	return nil
}

// dropBrowserFile deletes a file's first chunks parts and its record.
func (s *Server) dropBrowserFile(ctx context.Context, tenant, grant, id string, chunks int) {
	for n := 0; n < chunks; n++ {
		_, _, _ = s.KV.Delete(ctx, tenantPK(tenant), skBrowserFileChunk(grant, id, n))
	}
	_, _, _ = s.KV.Delete(ctx, tenantPK(tenant), skBrowserFile(grant, id))
}

// dropBrowserFiles deletes every file kept under g: its records and parts.
func (s *Server) dropBrowserFiles(ctx context.Context, g *Grant) error {
	var first error
	for _, prefix := range []string{"BFILE#" + g.ID + "#", "BFILEPART#" + g.ID + "#"} {
		items, err := s.KV.Query(ctx, tenantPK(g.Tenant), prefix)
		if err != nil {
			first = cmpErr(first, err)
			continue
		}
		for _, it := range items {
			if _, _, err := s.KV.Delete(ctx, it.PK, it.SK); err != nil {
				first = cmpErr(first, err)
			}
		}
	}
	return first
}

// errFileGone: no such file under the grant, or past its time.
var errFileGone = errors.New("no such file")

// openBrowserFile is a kept file and its bytes, checked against its digest.
func (s *Server) openBrowserFile(ctx context.Context, g *Grant, id string) (browserFile, []byte, error) {
	var f browserFile
	if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skBrowserFile(g.ID, id), &f); errors.Is(err, errNotFound) {
		return f, nil, errFileGone
	} else if err != nil {
		return f, nil, err
	}
	if !s.now().Before(f.At.Add(browserFileKept)) {
		return f, nil, errFileGone
	}
	data := make([]byte, 0, f.Size)
	for n := 0; n < f.Chunks; n++ {
		var part browserFileChunkRecord
		if _, err := getDoc(ctx, s.KV, tenantPK(g.Tenant), skBrowserFileChunk(g.ID, id, n), &part); err != nil {
			return f, nil, errFileGone
		}
		b, err := s.Sealer.Open(ctx, part.Sealed, browserFileBinding(g, id, n))
		if err != nil {
			return f, nil, err
		}
		data = append(data, b...)
	}
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != f.SHA256 || len(data) != f.Size {
		return f, nil, errors.New("the file doesn't match its digest")
	}
	return f, data, nil
}

// grantFile answers GET /v1/accounts/{tenant}/grants/{id}/files/{file},
// with the account's key: the file's bytes, as its type, to save and never
// to run (attachment, no sniffing, no-store).
func (s *Server) grantFile(w http.ResponseWriter, r *http.Request) {
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
	f, data, err := s.openBrowserFile(r.Context(), g, r.PathValue("file"))
	if errors.Is(err, errFileGone) {
		writeErr(w, http.StatusNotFound, "no such file, or it was deleted after 30 days")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	w.Header().Set("Content-Type", f.Type)
	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(f.Name, `"`, "")+`"`)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	w.Header().Set("X-Planque-SHA256", f.SHA256)
	_, _ = w.Write(data)
}
