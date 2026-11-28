package browserworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// downloads are the files one session's pages gave: saved by the browser
// into the session's own profile, under the id it gave each, never opened;
// read, typed and handed over by Downloads, then deleted. A file past
// browser.MaxDownload is stopped as it arrives.
type downloads struct {
	mu    sync.Mutex
	dir   string
	items map[string]*download
	order []string
}

type download struct {
	name     string // the page's suggested name
	state    string // inProgress, completed, canceled
	received int64
	tooLarge bool
}

// item is the download guid, made when first heard of. Under mu.
func (d *downloads) item(guid string) *download {
	it := d.items[guid]
	if it == nil {
		it = &download{state: "inProgress"}
		d.items[guid] = it
		d.order = append(d.order, guid)
	}
	return it
}

// downloadWait bounds how long Downloads waits for a file still arriving.
var downloadWait = 5 * time.Second

var guidRe = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)

// startDownloads lets the session's pages give files: saved into the
// profile's downloads folder, with the browser's events, so each is
// followed and a large one stopped.
func (s *session) startDownloads(ctx context.Context) error {
	dir := filepath.Join(s.c.dir, "downloads")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	s.dl = &downloads{dir: dir, items: map[string]*download{}}
	s.c.handle(s.followDownload)
	return s.c.call(ctx, "", "Browser.setDownloadBehavior", map[string]any{"behavior": "allowAndName", "downloadPath": dir, "eventsEnabled": true}, nil)
}

// followDownload keeps each download's progress, and stops one that grows
// past the cap.
func (s *session) followDownload(m cdpMessage) {
	switch m.Method {
	case "Browser.downloadWillBegin":
		var p struct {
			GUID              string `json:"guid"`
			SuggestedFilename string `json:"suggestedFilename"`
		}
		if json.Unmarshal(m.Params, &p) != nil || !guidRe.MatchString(p.GUID) {
			return
		}
		s.dl.mu.Lock()
		s.dl.item(p.GUID).name = p.SuggestedFilename
		s.dl.mu.Unlock()
	case "Browser.downloadProgress":
		var p struct {
			GUID          string  `json:"guid"`
			ReceivedBytes float64 `json:"receivedBytes"`
			TotalBytes    float64 `json:"totalBytes"`
			State         string  `json:"state"`
		}
		if json.Unmarshal(m.Params, &p) != nil || !guidRe.MatchString(p.GUID) {
			return
		}
		s.dl.mu.Lock()
		// Events are handled each in a goroutine of its own, so they can
		// come in any order: an item is made by whichever comes first, and
		// an ended one never goes back to arriving.
		d := s.dl.item(p.GUID)
		if d.state == "inProgress" {
			d.state = p.State
		}
		d.received = max(d.received, int64(p.ReceivedBytes))
		stop := p.State == "inProgress" && !d.tooLarge && (p.ReceivedBytes > browser.MaxDownload || p.TotalBytes > browser.MaxDownload)
		if stop {
			d.tooLarge = true
		}
		s.dl.mu.Unlock()
		if stop {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.c.call(ctx, "", "Browser.cancelDownload", map[string]any{"guid": p.GUID}, nil)
		}
	}
}

// Downloads hands over the files given since the last ask, each typed and
// taken, or refused with why, and forgets them; a file still arriving is
// waited for a moment, then refused.
func (s *session) Downloads(ctx context.Context) ([]browser.Download, error) {
	if s.closed {
		return nil, errors.New("the session is closed")
	}
	if s.dl == nil {
		return nil, nil
	}
	deadline := time.Now().Add(downloadWait)
	for {
		s.dl.mu.Lock()
		arriving := false
		for _, d := range s.dl.items {
			arriving = arriving || d.state == "inProgress" && !d.tooLarge
		}
		s.dl.mu.Unlock()
		if !arriving || time.Now().After(deadline) {
			break
		}
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	s.dl.mu.Lock()
	defer s.dl.mu.Unlock()
	out := []browser.Download{}
	for _, guid := range s.dl.order {
		d := s.dl.items[guid]
		path := filepath.Join(s.dl.dir, guid)
		out = append(out, takeDownload(d, path))
		_ = os.Remove(path)
		delete(s.dl.items, guid)
	}
	s.dl.order = nil
	return out, nil
}

// takeDownload reads one saved file and types it, or says why it isn't
// taken. The file is read only as bytes: nothing opens it.
func takeDownload(d *download, path string) browser.Download {
	name := browser.DownloadName(d.name, "")
	switch {
	case d.tooLarge:
		return browser.Download{Name: name, Size: int(d.received), Refused: "larger than 5 MB"}
	case d.state != "completed":
		return browser.Download{Name: name, Refused: "it didn't finish arriving"}
	}
	st, err := os.Stat(path)
	if err != nil || !st.Mode().IsRegular() {
		return browser.Download{Name: name, Refused: "it wasn't saved"}
	}
	if st.Size() > browser.MaxDownload {
		return browser.Download{Name: name, Size: int(st.Size()), Refused: "larger than 5 MB"}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return browser.Download{Name: name, Refused: "it couldn't be read"}
	}
	typ, err := browser.TypeOf(b)
	if err != nil {
		return browser.Download{Name: name, Size: len(b), Refused: err.Error()}
	}
	sum := sha256.Sum256(b)
	return browser.Download{Name: browser.DownloadName(d.name, typ), Type: typ, Size: len(b), SHA256: hex.EncodeToString(sum[:]), Data: b,
		Formulas: typ == browser.TypeCSV && browser.CSVFormulas(b)}
}
