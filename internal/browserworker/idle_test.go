package browserworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// slowWorker is a worker whose calls take longer than the idle time: what
// a slow runner made of the real one (the CI failure on bb279b7).
type slowWorker struct {
	delay  time.Duration
	closed atomic.Int32
}

func (w *slowWorker) Start(context.Context, browser.Start) (browser.Session, error) {
	return &slowSession{w: w}, nil
}

type slowSession struct{ w *slowWorker }

func (s *slowSession) Open(ctx context.Context, u string) (browser.Page, error) {
	time.Sleep(s.w.delay)
	return browser.Page{URL: u, Title: "slow"}, nil
}
func (s *slowSession) Do(context.Context, browser.Step) (browser.Page, error) {
	time.Sleep(s.w.delay)
	return browser.Page{Title: "stepped"}, nil
}
func (s *slowSession) Screenshot(context.Context) ([]byte, error) { return []byte("png"), nil }
func (s *slowSession) Close(context.Context) ([]byte, error) {
	s.w.closed.Add(1)
	return []byte(`{"cookies":[]}`), nil
}

// A session in a call is never idle, however long the call takes: idle
// time counts from the end of the last call, and the next call finds the
// session. Left alone, it is ended.
func TestWorkerNeverReapsASessionInACall(t *testing.T) {
	w := &slowWorker{delay: 600 * time.Millisecond}
	sum := sha256.Sum256([]byte("k"))
	srv := &Server{Local: w, KeySHA256: hex.EncodeToString(sum[:]), Idle: 200 * time.Millisecond}
	api := httptest.NewServer(srv.Handler())
	defer api.Close()
	c := &browser.Client{URL: api.URL, Key: "k"}
	ctx := context.Background()
	sess, err := c.Start(ctx, browser.Start{Grant: "g"})
	if err != nil {
		t.Fatal(err)
	}
	if p, err := sess.Open(ctx, "https://shop.example/"); err != nil || p.Title != "slow" {
		t.Fatalf("a slow open: %v %+v", err, p)
	}
	if p, err := sess.Do(ctx, browser.Step{Click: "c1"}); err != nil || p.Title != "stepped" {
		t.Fatalf("the step after a slow open: %v %+v", err, p)
	}
	if n := w.closed.Load(); n != 0 {
		t.Fatalf("the session was ended during its calls (%d)", n)
	}
	time.Sleep(600 * time.Millisecond)
	if _, err := sess.Open(ctx, "https://shop.example/"); err == nil {
		t.Error("the idle session is still there")
	}
	if n := w.closed.Load(); n != 1 {
		t.Errorf("closed %d times, want 1", n)
	}
}
