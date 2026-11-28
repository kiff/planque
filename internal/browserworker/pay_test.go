package browserworker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/cardfill"
	"github.com/kiff/planque/internal/fakes"
)

// syncBuf is a log both the worker and its proxy write to.
type syncBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

var testCard = cardfill.Card{Number: "4242 4242 4242 4242", ExpMonth: 12, ExpYear: 2034, CVC: "123", Name: "Acme SAS"}

var testDelivery = cardfill.Address{Name: "Acme SAS", Line1: "12 rue des Lilas", PostalCode: "69003", City: "Lyon", Country: "FR"}

type payRig struct {
	shop  *fakes.Shop
	w     *PayWorker
	logs  *syncBuf
	armed int
}

// payRig is the payment worker on the test shop, reached as shop.example
// over TLS, in test mode, logging to a buffer the canary reads.
func newPayRig(t *testing.T) *payRig {
	t.Helper()
	if _, err := chromePath(""); err != nil {
		t.Skip(err)
	}
	shop := fakes.NewTLSShop()
	t.Cleanup(shop.Close)
	logs := &syncBuf{}
	r := &payRig{shop: shop, logs: logs}
	r.w = &PayWorker{
		Mode:       cardfill.Test,
		Root:       t.TempDir(),
		AllowLocal: true,
		Map:        map[string]string{"shop.example:443": shop.Listener.Addr().String()},
		// Nothing is resolved: the shop is reached through the tests' map
		// only, and any other host is refused.
		Resolve: func(context.Context, string) ([]netip.Addr, error) { return nil, errors.New("no resolver in tests") },
		Logger:  slog.New(slog.NewTextHandler(logs, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Wait:    20 * time.Second,
	}
	return r
}

// approved is the pay checkout as read before the payment: the lines,
// total and currency the owner approved.
func (r *payRig) approved(t *testing.T) browser.Checkout {
	t.Helper()
	w := &Local{Options: Options{AllowLocal: true, NavigateWait: 20 * time.Second}}
	ctx := context.Background()
	sess, err := w.Start(ctx, browser.Start{Grant: "grt_test"})
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close(ctx)
	p, err := sess.Open(ctx, r.shop.URL+"/pay")
	if err != nil || p.Checkout == nil {
		t.Fatalf("read the pay checkout: %+v, %v", p, err)
	}
	return *p.Checkout
}

func (r *payRig) job(t *testing.T, path string) PayJob {
	return PayJob{
		Checkout: "https://shop.example" + path,
		Merchant: "shop.example",
		Approved: r.approved(t),
		Delivery: testDelivery,
		Card:     testCard,
		Arm:      func(context.Context) error { r.armed++; return nil },
	}
}

// profiles is what the worker left under its root.
func profiles(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// canaryFiles runs the canary over every file under dir.
func canaryFiles(t *testing.T, dir string, cards ...string) []cardfill.Hit {
	t.Helper()
	var hits []cardfill.Hit
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		// UTF-16 too: some of the browser's stores write text that way.
		text := string(b) + "\n" + strings.ReplaceAll(string(b), "\x00", "")
		// The browser's own files hold its timestamps (17-digit microsecond
		// counts, about one in ten Luhn-valid), so here the canary looks for
		// the card filled: its last eight, which its whole number holds too,
		// in any spelling. The result and the logs get the full canary.
		for _, h := range cardfill.Canary(path, text, cards...) {
			if h.What != "card number" {
				hits = append(hits, h)
			}
		}
		return nil
	})
	return hits
}

// 1. The canary test, which gates B2: a test card is filled into the
// checkout and the order is placed once; no card number, nor its last
// eight, is in the result, the worker's log, the proxy's log, or any file
// the browser wrote while the card was in the fields; and the profile is
// gone after.
func TestPayFillsOnceAndTheCanaryFindsNothing(t *testing.T) {
	r := newPayRig(t)
	job := r.job(t, "/pay")
	var profileHits []cardfill.Hit
	var profileDir string
	r.w.afterFill = func(c *chrome) {
		profileDir = c.dir
		profileHits = canaryFiles(t, c.dir, testCard.Number)
	}
	res, err := r.w.Pay(context.Background(), job)
	if err != nil {
		t.Fatalf("pay: %+v, %v", res, err)
	}
	if res.Outcome != "ordered" || res.OrderRef != "A-1001" || !strings.Contains(res.Confirmation, "Order A-1001 placed") {
		t.Errorf("result: %+v", res)
	}
	if r.armed != 1 || r.shop.PayCalls != 1 || len(r.shop.Paid) != 1 || r.shop.Paid[0] != "4242" {
		t.Errorf("armed %d, the shop's orders %d, paid with %v: want one order, armed once, with the test card", r.armed, r.shop.PayCalls, r.shop.Paid)
	}
	if profileDir == "" {
		t.Fatal("the fill never ran")
	}
	if len(profileHits) > 0 {
		t.Errorf("the profile held the card while it was in the fields: %+v", profileHits)
	}
	out, _ := json.Marshal(res)
	for where, text := range map[string]string{"result": string(out), "logs": r.logs.String(), "error": errString(err)} {
		if hits := cardfill.Canary(where, text, testCard.Number); len(hits) > 0 {
			t.Errorf("the canary found %+v", hits)
		}
	}
	if left := profiles(t, r.w.Root); len(left) > 0 {
		t.Errorf("the payment left %v under its root", left)
	}
	if _, err := os.Stat(profileDir); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the profile is still there: %v", err)
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// 2. The delivery address is the case's, read from the address fields when
// the page has them; another address and nothing is filled or armed.
func TestPayChecksTheDeliveryAddress(t *testing.T) {
	r := newPayRig(t)
	if res, err := r.w.Pay(context.Background(), r.job(t, "/pay?fields=1")); err != nil || res.Outcome != "ordered" {
		t.Fatalf("the address fields: %+v, %v", res, err)
	}
	r.shop.Ship = [3]string{"1 avenue Foch", "75116", "Paris"}
	for _, path := range []string{"/pay", "/pay?fields=1"} {
		res, err := r.w.Pay(context.Background(), r.job(t, path))
		if !errors.Is(err, cardfill.ErrAddressNotShown) || res.Outcome != "not_ordered" {
			t.Errorf("%s with another address: %+v, %v", path, res, err)
		}
	}
	if r.armed != 1 || r.shop.PayCalls != 1 {
		t.Errorf("armed %d, orders %d: the other address must arm and order nothing", r.armed, r.shop.PayCalls)
	}
}

// 3. The checkout as approved: another total and nothing is filled.
func TestPayRefusesAChangedCheckout(t *testing.T) {
	r := newPayRig(t)
	job := r.job(t, "/pay")
	r.shop.Total = "25.00"
	res, err := r.w.Pay(context.Background(), job)
	if err == nil || res.Outcome != "not_ordered" || r.armed != 0 || r.shop.PayCalls != 0 {
		t.Errorf("a changed checkout: %+v, %v; armed %d, orders %d", res, err, r.armed, r.shop.PayCalls)
	}
}

// 4. One control that pays, never a guess between two.
func TestPayNeedsOneControlThatPays(t *testing.T) {
	r := newPayRig(t)
	res, err := r.w.Pay(context.Background(), r.job(t, "/pay?two=1"))
	if err == nil || !strings.Contains(err.Error(), "more than one") || res.Outcome != "not_ordered" || r.shop.PayCalls != 0 {
		t.Errorf("two controls: %+v, %v", res, err)
	}
}

// 5. The pay-once guard: armed before the click; a guard that refuses
// (another payment of the same case armed it first) clicks nothing.
func TestPayClicksNothingWhenTheGuardRefuses(t *testing.T) {
	r := newPayRig(t)
	job := r.job(t, "/pay")
	job.Arm = func(context.Context) error { return errors.New("armed already") }
	res, err := r.w.Pay(context.Background(), job)
	if err == nil || res.Outcome != "not_ordered" || r.shop.PayCalls != 0 {
		t.Errorf("a refused guard: %+v, %v; orders %d", res, err, r.shop.PayCalls)
	}
	job.Arm = nil
	if res, err := r.w.Pay(context.Background(), job); err == nil || res.Outcome != "not_ordered" {
		t.Errorf("no guard: %+v, %v", res, err)
	}
}

// 6. Test mode only, test cards only, and the merchant's site only: each
// refused before a browser starts.
func TestPayRefusesBeforeABrowserStarts(t *testing.T) {
	r := newPayRig(t)
	job := r.job(t, "/pay")
	cases := map[string]func() (*PayWorker, PayJob){
		"off": func() (*PayWorker, PayJob) {
			return &PayWorker{Mode: cardfill.Off, Root: r.w.Root, AllowLocal: true, Map: r.w.Map, Resolve: r.w.Resolve}, job
		},
		"a real card": func() (*PayWorker, PayJob) {
			j := job
			j.Card.Number = "4000 0012 3456 7899"
			return r.w, j
		},
		"another site": func() (*PayWorker, PayJob) {
			j := job
			j.Checkout = "https://other.example/pay"
			return r.w, j
		},
		"plain http": func() (*PayWorker, PayJob) {
			j := job
			j.Checkout = "http://shop.example/pay"
			return r.w, j
		},
	}
	for name, mk := range cases {
		w, j := mk()
		res, err := w.Pay(context.Background(), j)
		if err == nil || res.Outcome != "not_ordered" {
			t.Errorf("%s: %+v, %v", name, res, err)
		}
		if strings.Contains(errString(err), "4000") {
			t.Errorf("%s: the error carries the number: %v", name, err)
		}
	}
	if r.armed != 0 || r.shop.PayCalls != 0 || len(profiles(t, r.w.Root)) != 0 {
		t.Errorf("armed %d, orders %d, profiles %v", r.armed, r.shop.PayCalls, profiles(t, r.w.Root))
	}
}

// 7. A browser killed with the card in the fields leaves no file: the
// payment ends not ordered and its profile is deleted; and what a killed
// worker left is swept before the next payment.
func TestPayKilledLeavesNoFile(t *testing.T) {
	r := newPayRig(t)
	job := r.job(t, "/pay")
	r.w.afterFill = func(c *chrome) { _ = c.cmd.Process.Kill() }
	res, err := r.w.Pay(context.Background(), job)
	if err == nil || res.Outcome != "not_ordered" || r.armed != 0 || r.shop.PayCalls != 0 {
		t.Errorf("killed: %+v, %v; armed %d, orders %d", res, err, r.armed, r.shop.PayCalls)
	}
	if left := profiles(t, r.w.Root); len(left) > 0 {
		t.Errorf("the killed payment left %v", left)
	}

	// A worker killed itself leaves its profile; the next payment sweeps it.
	stale := filepath.Join(r.w.Root, payPrefix+"stale")
	if err := os.MkdirAll(filepath.Join(stale, "Default"), 0o700); err != nil {
		t.Fatal(err)
	}
	r.w.afterFill = nil
	if res, err := r.w.Pay(context.Background(), job); err != nil || res.Outcome != "ordered" {
		t.Fatalf("the next payment: %+v, %v", res, err)
	}
	if left := profiles(t, r.w.Root); len(left) > 0 {
		t.Errorf("after the next payment: %v", left)
	}
}

// 8. One payment at a time.
func TestPayOneAtATime(t *testing.T) {
	w := &PayWorker{Mode: cardfill.Test, Root: t.TempDir()}
	w.mu.Lock()
	defer w.mu.Unlock()
	if res, err := w.Pay(context.Background(), PayJob{}); !errors.Is(err, ErrPayBusy) || res.Outcome != "not_ordered" {
		t.Errorf("busy: %+v, %v", res, err)
	}
}

// "ordered" only when the page after says so: a failure page, even with a
// reference on it, and a page with no confirmation or reference are
// unknown, for reconciliation to settle.
func TestPayOrderedOnlyWhenThePageConfirmsIt(t *testing.T) {
	r := newPayRig(t)
	for path, want := range map[string]string{"/pay?end=failed": "unknown", "/pay?end=plain": "unknown", "/pay": "ordered"} {
		res, err := r.w.Pay(context.Background(), r.job(t, path))
		if res.Outcome != want || (want == "ordered") != (err == nil) {
			t.Errorf("%s: %+v, %v; want %s", path, res, err, want)
		}
	}
}

// An expiry year field that takes two digits gets two.
func TestPayTypesTheYearAsTheFieldTakesIt(t *testing.T) {
	r := newPayRig(t)
	if res, err := r.w.Pay(context.Background(), r.job(t, "/pay?split=1")); err != nil || res.Outcome != "ordered" {
		t.Fatalf("split expiry: %+v, %v", res, err)
	}
	if _, err := r.w.Pay(context.Background(), r.job(t, "/pay")); err != nil {
		t.Fatal(err)
	}
	if len(r.shop.Exp) != 2 || r.shop.Exp[0] != "12/34" || r.shop.Exp[1] != "12 / 34" {
		t.Errorf("expiries received: %q", r.shop.Exp)
	}
}

// Card fields inside a provider's frame aren't filled: refused, with why.
func TestPayRefusesCardFieldsInAProvidersFrame(t *testing.T) {
	r := newPayRig(t)
	res, err := r.w.Pay(context.Background(), r.job(t, "/pay?frame=1"))
	if res.Outcome != "not_ordered" || err == nil || !strings.Contains(err.Error(), "payment provider's frame") || r.armed != 0 || r.shop.PayCalls != 0 {
		t.Errorf("frame: %+v, %v", res, err)
	}
}
