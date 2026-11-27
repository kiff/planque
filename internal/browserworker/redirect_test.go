package browserworker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// A listed site that redirects off the session's sites reaches nothing: the
// real browser follows the 302, and the worker's proxy refuses the next hop,
// whether it's another site, a listed name that resolves to the task's
// metadata address, or that address written out. The red-team's items 2
// and 3 (27 Nov) couldn't make a granted site redirect off-site live, so
// this closes them with real Chromium.
func TestWorkerRefusesARedirectOffTheSites(t *testing.T) {
	if _, err := chromePath(""); err != nil {
		t.Skip(err)
	}
	// What a refused hop would have reached: a server that counts, mapped
	// for other.example so that only the proxy's rule stands between them.
	var trapped atomic.Int32
	trap := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trapped.Add(1)
		_, _ = w.Write([]byte(`{"secret":"metadata"}`))
	}))
	defer trap.Close()

	site := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, _ = w.Write([]byte(`<html><head><title>Shop</title></head><body><h1>The shop</h1></body></html>`))
		case "/off":
			http.Redirect(w, r, "https://other.example/secret", http.StatusFound)
		case "/meta-name":
			http.Redirect(w, r, "https://meta.shop.example/v2/metadata", http.StatusFound)
		case "/meta-ip":
			http.Redirect(w, r, "https://169.254.170.2/v2/metadata", http.StatusFound)
		case "/meta-http":
			http.Redirect(w, r, "http://169.254.170.2/v2/metadata", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer site.Close()

	var mu sync.Mutex
	var asked []string
	w := &Local{Options: Options{
		AllowLocal:   true,
		NavigateWait: 20 * time.Second,
		Map: map[string]string{
			"shop.example:443":  site.Listener.Addr().String(),
			"other.example:443": trap.Listener.Addr().String(),
		},
		// Every name answers with the task's metadata address: a listed
		// name is refused for it, and shop.example only reaches the site
		// through the tests' map.
		Resolve: func(_ context.Context, host string) ([]netip.Addr, error) {
			mu.Lock()
			defer mu.Unlock()
			asked = append(asked, host)
			return []netip.Addr{netip.MustParseAddr("169.254.170.2")}, nil
		},
	}}
	ctx := context.Background()
	sess, err := w.Start(ctx, startFor(site, "shop.example"))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)

	// The listed site itself opens.
	p, err := sess.Open(ctx, "https://shop.example/")
	if err != nil || !strings.Contains(p.Text, "The shop") {
		t.Fatalf("the listed site: %+v, %v", p, err)
	}
	for _, path := range []string{"/off", "/meta-name", "/meta-ip", "/meta-http"} {
		before := sess.(*session).proxy.Refused()
		p, err := sess.Open(ctx, "https://shop.example"+path)
		// Refused as the egress refuses, never a page: over plain http too,
		// where the browser would show the proxy's answer at the address.
		if err == nil || !strings.Contains(err.Error(), "egress") {
			t.Errorf("%s: want the egress's refusal, got %+v, %v", path, p, err)
		}
		if sess.(*session).proxy.Refused() <= before {
			t.Errorf("%s: the proxy refused nothing", path)
		}
	}
	// The listed site still opens after them.
	if p, err := sess.Open(ctx, "https://shop.example/"); err != nil || !strings.Contains(p.Text, "The shop") {
		t.Errorf("the listed site after the refusals: %+v, %v", p, err)
	}
	if n := trapped.Load(); n != 0 {
		t.Errorf("the unlisted site was reached %d times", n)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, h := range asked {
		if h == "other.example" {
			t.Error("the unlisted site was resolved: it must be refused by name first")
		}
	}
}
