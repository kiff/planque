package browserworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// planque's client, over the worker's API, is a Worker like the local
// one: the same shop, the same facts, the state round trip, and the key
// checked before anything else; an idle session ends on its own.
func TestWorkerAPIServesPlanquesClient(t *testing.T) {
	w, s := worker(t), shop(t)
	sum := sha256.Sum256([]byte("worker-key"))
	srv := &Server{Local: w, KeySHA256: hex.EncodeToString(sum[:]), Idle: 400 * time.Millisecond}
	api := httptest.NewServer(srv.Handler())
	defer api.Close()
	defer srv.Shutdown(context.Background())
	ctx := context.Background()

	wrong := &browser.Client{URL: api.URL, Key: "not-the-key"}
	if _, err := wrong.Start(ctx, browser.Start{Grant: "grt_x"}); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatalf("a wrong key: %v", err)
	}
	none := &browser.Client{URL: api.URL}
	if _, err := none.Start(ctx, browser.Start{Grant: "grt_x"}); err == nil {
		t.Fatal("no key started a session")
	}

	c := &browser.Client{URL: api.URL, Key: "worker-key"}
	sess, err := c.Start(ctx, browser.Start{Grant: "grt_x"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := sess.Open(ctx, s.URL+"/checkout")
	if err != nil {
		t.Fatal(err)
	}
	if p.Checkout == nil || p.Checkout.Total != "2.50" || find(p.Fields, "Card number") == nil || !find(p.Fields, "Card number").Masked {
		t.Fatalf("over the API: %+v", p)
	}
	p, err = sess.Do(ctx, stepFill(find(p.Fields, "Coupon").Ref, "SPRING"))
	if err != nil || find(p.Fields, "Coupon").Value != "SPRING" {
		t.Fatalf("a step over the API: %v %+v", err, p.Fields)
	}
	if png, err := sess.Screenshot(ctx); err != nil || !strings.HasPrefix(string(png), "\x89PNG") {
		t.Errorf("a screenshot over the API: %v %d", err, len(png))
	}
	st, err := sess.Close(ctx)
	if err != nil || !strings.Contains(string(st), "cookies") {
		t.Fatalf("close over the API: %v %s", err, st)
	}
	// Closed is closed.
	if _, err := sess.Open(ctx, s.URL+"/cart"); err == nil || !strings.Contains(err.Error(), "no such session") {
		t.Errorf("a closed session: %v", err)
	}
	// An idle session is ended by the worker, so a caller that died
	// leaves no browser behind.
	idle, err := c.Start(ctx, browser.Start{Grant: "grt_y"})
	if err != nil {
		t.Fatal(err)
	}
	// Each call refreshes the idle time, so the test waits without one.
	time.Sleep(1200 * time.Millisecond)
	if _, err := idle.Open(ctx, s.URL+"/cart"); err == nil || !strings.Contains(err.Error(), "no such session") {
		t.Errorf("the idle session is still there: %v", err)
	}
	// The health check needs no key.
	res, err := http.Get(api.URL + "/healthz")
	if err != nil || res.StatusCode != 200 {
		t.Errorf("healthz: %v %v", err, res)
	}
}
