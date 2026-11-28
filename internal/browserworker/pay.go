package browserworker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/cardfill"
)

// PayWorker is the payment-only worker (Thune v1, path (b), test mode only
// until the card-data decision): it fills one card into one checkout and
// places that one order, in a browser started for that payment alone and
// ended after it. It is not the browser worker's session API: there is no
// screenshot, no live view, no kept state and no step from an agent. Each
// payment gets:
//   - a fresh profile, under Root, with autofill and saved cards and
//     passwords off, deleted when the payment ends;
//   - an off-the-record browser context inside it, so the browser writes
//     nothing of the payment to disk;
//   - egress only to the merchant's site and the published payment hosts
//     (cardfill.Egress), through a proxy of its own;
//   - no network or console events kept: only the page domain is on.
//
// The card is typed with the browser's text input into the focused field,
// so it is never part of a script, an error or a log. At most one payment
// runs at a time.
type PayWorker struct {
	// Mode is cardfill's mode: Off refuses every payment.
	Mode cardfill.Mode
	// Chrome is the browser's path; empty, chromePath looks.
	Chrome string
	// Root holds the payment profiles; empty, planque-pay in the system's
	// temporary directory. It is the worker's own: what a killed payment
	// left there is deleted before the next one starts.
	Root string
	// AllowLocal and Map, tests only, as Options'.
	AllowLocal bool
	Map        map[string]string
	Resolve    func(ctx context.Context, host string) ([]netip.Addr, error)
	Logger     *slog.Logger
	// Wait bounds each page load; 0 is 20 s.
	Wait time.Duration

	// afterFill, tests only, runs once the card is in the fields and
	// before the payment is armed, with the profile's directory.
	afterFill func(c *chrome)

	mu sync.Mutex
}

// PayJob is one payment: what KIFF's signed case says, and the card. No
// PayJob is built outside tests until the single-use card's number reaches
// the worker alone (the threat model's T3.7): nothing in planque's own code
// path ever holds a card number to put in one.
type PayJob struct {
	// Checkout is the checkout's https address, on Merchant.
	Checkout string
	// Merchant is the merchant's site, a host name, as the case names it.
	Merchant string
	// Approved is the checkout as read and approved: the page must show
	// the same lines, total and currency, or nothing is filled.
	Approved browser.Checkout
	// Delivery is the case's signed delivery address.
	Delivery cardfill.Address
	// Card is the one card filled.
	Card cardfill.Card
	// Arm marks the payment as sent before the click (planque's pay-once
	// guard: a conditional write, so a second payment of the same case
	// finds it marked and sends nothing). Its error means nothing is
	// clicked. Required.
	Arm func(ctx context.Context) error
}

// payPrefs are the profile's preferences: no autofill of addresses or
// cards, no saved passwords, no payment handlers.
var payPrefs = []byte(`{"autofill":{"enabled":false,"profile_enabled":false,"credit_card_enabled":false},` +
	`"credentials_enable_service":false,"profile":{"password_manager_enabled":false},` +
	`"payments":{"can_make_payment_enabled":false}}`)

const payPrefix = "planque-pay-"

// ErrPayBusy: a payment runs already.
var ErrPayBusy = errors.New("the payment worker is busy with another payment")

func (w *PayWorker) root() string {
	if w.Root != "" {
		return w.Root
	}
	return filepath.Join(os.TempDir(), "planque-pay")
}

func (w *PayWorker) log() *slog.Logger {
	if w.Logger != nil {
		return w.Logger
	}
	return slog.Default()
}

// Sweep deletes what an ended payment left under Root: a profile of a
// worker killed before it could delete its own. It runs before every
// payment, and may run at start.
func (w *PayWorker) Sweep() error {
	entries, err := os.ReadDir(w.root())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	var first error
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), payPrefix) {
			if err := os.RemoveAll(filepath.Join(w.root(), e.Name())); err != nil && first == nil {
				first = err
			}
		}
	}
	return first
}

// notOrdered and unknown are the two ways a payment ends without an order
// known: before the arm nothing was clicked; after it, the order may exist.
func notOrdered(err error) (cardfill.Result, error) {
	return cardfill.Result{Outcome: "not_ordered"}, err
}

func unknown(err error) (cardfill.Result, error) {
	return cardfill.Result{Outcome: "unknown"}, err
}

// Pay runs one payment and returns its result, cleaned (cardfill's
// Result.Clean: a result the canary finds a card in is dropped as
// unknown). The error says why a payment didn't end ordered; it never
// carries the card.
func (w *PayWorker) Pay(ctx context.Context, job PayJob) (cardfill.Result, error) {
	if !w.mu.TryLock() {
		return notOrdered(ErrPayBusy)
	}
	defer w.mu.Unlock()
	if err := cardfill.Fillable(w.Mode, job.Card.Number); err != nil {
		return notOrdered(err)
	}
	if job.Arm == nil {
		return notOrdered(errors.New("a payment needs its pay-once guard"))
	}
	if job.Approved.Total == "" || job.Approved.Currency == "" {
		return notOrdered(errors.New("a payment needs the checkout as approved"))
	}
	sites, err := cardfill.Egress(job.Merchant)
	if err != nil {
		return notOrdered(err)
	}
	if code, why := browser.Allowed(job.Checkout, sites[:1]); code != "" {
		return notOrdered(fmt.Errorf("the checkout isn't on the merchant's site: %s", why))
	}
	path, err := chromePath(w.Chrome)
	if err != nil {
		return notOrdered(err)
	}
	if err := os.MkdirAll(w.root(), 0o700); err != nil {
		return notOrdered(err)
	}
	if err := w.Sweep(); err != nil {
		return notOrdered(fmt.Errorf("an earlier payment's profile couldn't be deleted: %w", err))
	}

	p := &Proxy{Sites: sites, AllowLocal: w.AllowLocal, Map: w.Map, Logger: w.Logger, Resolve: w.Resolve}
	addr, err := p.Start()
	if err != nil {
		return notOrdered(err)
	}
	// No crash dump: the browser's crash reporter is off, and the process
	// may write no core file, so a crash with the card in memory leaves it
	// nowhere on disk.
	if err := noCoreDumps(); err != nil {
		p.Close()
		return notOrdered(fmt.Errorf("core dumps couldn't be turned off, so nothing was started: %w", err))
	}
	c, err := launchIn(context.Background(), path, addr, w.AllowLocal, w.root(), payPrefix, payPrefs,
		"--disable-features=AutofillServerCommunication", "--incognito", "--disable-crash-reporter", "--disable-breakpad")
	if err != nil {
		p.Close()
		return notOrdered(err)
	}
	s := &session{c: c, proxy: p, wait: w.Wait, logger: w.Logger}
	if s.wait == 0 {
		s.wait = 20 * time.Second
	}
	defer s.end()
	if err := w.attach(ctx, s); err != nil {
		return notOrdered(err)
	}

	page, err := s.Open(ctx, job.Checkout)
	if err != nil {
		return notOrdered(err)
	}
	if code, why := browser.Allowed(page.URL, sites[:1]); code != "" {
		return notOrdered(fmt.Errorf("the checkout redirected: %s", why))
	}
	if co := page.Checkout; co == nil || !co.Same(job.Approved) {
		return notOrdered(errors.New("the checkout isn't what was approved: nothing was filled"))
	}
	if err := w.addressShown(ctx, s, job.Delivery); err != nil {
		return notOrdered(err)
	}
	payRef, err := onePay(page)
	if err != nil {
		return notOrdered(err)
	}
	// Card fields inside a payment provider's cross-site frame (a Stripe
	// element) aren't reached yet: such a checkout is refused, never
	// half-filled.
	if page.PaymentFrame && !slices.ContainsFunc(page.Fields, func(f browser.Field) bool { return f.Kind == "card" }) {
		return notOrdered(errors.New("the checkout's card fields are inside the payment provider's frame, which this worker doesn't fill yet: nothing was filled or paid"))
	}

	if err := w.fill(ctx, s, job.Card); err != nil {
		return notOrdered(err)
	}
	if w.afterFill != nil {
		w.afterFill(c)
	}
	// Read again with the card in: the address and the checkout are still
	// the approved ones, and the control is still the one that pays.
	again, err := s.page(ctx)
	if err != nil {
		return notOrdered(err)
	}
	if co := again.Checkout; co == nil || !co.Same(job.Approved) {
		return notOrdered(errors.New("the checkout changed while the card was filled: nothing was paid"))
	}
	if err := w.addressShown(ctx, s, job.Delivery); err != nil {
		return notOrdered(err)
	}
	if ref, err := onePay(again); err != nil || ref != payRef {
		return notOrdered(errors.New("the control that pays changed while the card was filled: nothing was paid"))
	}

	// Armed before the click, never after its error.
	if err := job.Arm(ctx); err != nil {
		return notOrdered(fmt.Errorf("the payment couldn't be armed, so nothing was paid: %w", err))
	}
	after, err := s.Do(ctx, browser.Step{Click: payRef})
	if err != nil {
		w.log().Warn("pay click: outcome unknown", slog.String("merchant", job.Merchant), slog.String("err", err.Error()))
		return unknown(errors.New("the order was placed, but what the site did with it isn't known"))
	}
	// Ordered only when the page after says so: it left the checkout, and
	// it names an order reference or confirms the order in words. A
	// "payment failed" page, or a refused 3-D Secure, is unknown, for
	// reconciliation to settle: the owner never reads "bought" before the
	// money says so.
	r := cardfill.Result{Outcome: "unknown", Confirmation: confirmation(after)}
	if placed(after, page) && !failed(after) {
		if ref := orderRef(after.Text); ref != "" || confirmed(after) {
			r.Outcome, r.OrderRef = "ordered", ref
		}
	}
	clean, err := r.Clean(job.Card.Number)
	if err != nil {
		w.log().Warn("pay result dropped", slog.String("merchant", job.Merchant), slog.String("err", err.Error()))
		return clean, err
	}
	if clean.Outcome != "ordered" {
		return clean, errors.New("the site's page after the click doesn't show an order: reconciliation settles it")
	}
	return clean, nil
}

// attach opens the payment's page in an off-the-record context, downloads
// denied, WebRTC taken away from every document and frame, and only the
// page domain on: no network or console event is kept.
func (w *PayWorker) attach(ctx context.Context, s *session) error {
	var bc struct {
		ID string `json:"browserContextId"`
	}
	if err := s.c.call(ctx, "", "Target.createBrowserContext", map[string]any{"disposeOnDetach": true}, &bc); err != nil {
		return err
	}
	if err := s.c.call(ctx, "", "Browser.setDownloadBehavior", map[string]any{"behavior": "deny", "browserContextId": bc.ID}, nil); err != nil {
		return err
	}
	var created struct {
		TargetID string `json:"targetId"`
	}
	if err := s.c.call(ctx, "", "Target.createTarget", map[string]any{"url": "about:blank", "browserContextId": bc.ID}, &created); err != nil {
		return err
	}
	s.target = created.TargetID
	var attached struct {
		SessionID string `json:"sessionId"`
	}
	if err := s.c.call(ctx, "", "Target.attachToTarget", map[string]any{"targetId": s.target, "flatten": true}, &attached); err != nil {
		return err
	}
	s.sid = attached.SessionID
	if err := s.c.call(ctx, s.sid, "Page.enable", nil, nil); err != nil {
		return err
	}
	if err := s.c.call(ctx, s.sid, "Page.addScriptToEvaluateOnNewDocument", map[string]any{"source": noWebRTC}, nil); err != nil {
		return err
	}
	s.c.handle(s.armChild)
	return s.c.call(ctx, s.sid, "Target.setAutoAttach", map[string]any{"autoAttach": true, "waitForDebuggerOnStart": true, "flatten": true}, nil)
}

// addressJS reads the delivery address as the page shows it: its address
// fields by autocomplete token (billing ones left out), and its text. It
// reads no card field.
const addressJS = `(() => {
  const want = {'address-line1': 1, 'postal-code': 1, 'address-level2': 1, 'country': 1, 'country-name': 1};
  const fields = {};
  for (const el of document.querySelectorAll('input,select,textarea')) {
    const toks = (el.getAttribute('autocomplete') || '').toLowerCase().split(/\s+/);
    if (toks.includes('billing')) continue;
    for (const t of toks) {
      if (!want[t]) continue;
      const k = t === 'country-name' ? 'country' : t;
      fields[k] = el.tagName === 'SELECT' ? (el.selectedOptions[0] ? el.selectedOptions[0].value : '') : el.value;
    }
  }
  return JSON.stringify({fields, text: (document.body ? document.body.innerText : '').slice(0, 20000)});
})()`

func (w *PayWorker) addressShown(ctx context.Context, s *session, a cardfill.Address) error {
	var ev struct {
		Result struct {
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": addressJS, "returnByValue": true}, &ev); err != nil {
		return err
	}
	var shown struct {
		Fields map[string]string `json:"fields"`
		Text   string            `json:"text"`
	}
	if err := json.Unmarshal([]byte(ev.Result.Value), &shown); err != nil {
		return errors.New("the checkout's delivery address couldn't be read")
	}
	return cardfill.AddressShown(a, shown.Fields, shown.Text)
}

// focusJS focuses the one visible field with an autocomplete token and
// selects what it holds, so the text typed next replaces it. It carries no
// value: the card is typed, never put in a script.
const focusJS = `((tok) => {
  const els = [...document.querySelectorAll('input')].filter(e => !e.disabled && e.getClientRects().length > 0 &&
    (e.getAttribute('autocomplete') || '').toLowerCase().split(/\s+/).includes(tok));
  if (els.length !== 1) return els.length ? 'many' : 'none';
  els[0].focus();
  if (els[0].select) els[0].select();
  return document.activeElement === els[0] ? 'ok' : 'unfocused';
})`

// yearTwoJS says the visible expiry-year field takes two digits: a
// maxlength of 2, or a placeholder of two letters (YY, AA, JJ).
const yearTwoJS = `(() => {
  const el = [...document.querySelectorAll('input')].find(e => e.getClientRects().length > 0 &&
    (e.getAttribute('autocomplete') || '').toLowerCase().split(/\s+/).includes('cc-exp-year'));
  return !!el && (el.maxLength === 2 || /^[a-z]{2}$/i.test((el.placeholder || '').trim()));
})()`

// fill types the card into the checkout's card fields: the number, the
// expiry (one field, or month and year), the security code, and the
// name when the page has a field for it.
func (w *PayWorker) fill(ctx context.Context, s *session, card cardfill.Card) error {
	focus := func(tok string) (string, error) {
		var ev struct {
			Result struct {
				Value string `json:"value"`
			} `json:"result"`
		}
		arg, _ := json.Marshal(tok)
		err := s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": focusJS + "(" + string(arg) + ")", "returnByValue": true}, &ev)
		return ev.Result.Value, err
	}
	typeIn := func(tok, text string, required bool) error {
		got, err := focus(tok)
		if err != nil {
			return err
		}
		if got != "ok" {
			if !required && got == "none" {
				return nil
			}
			return fmt.Errorf("the checkout's %s field: %s", tok, got)
		}
		if err := s.c.call(ctx, s.sid, "Input.insertText", map[string]any{"text": text}, nil); err != nil {
			return fmt.Errorf("the checkout's %s field couldn't be typed into", tok)
		}
		return nil
	}
	if err := typeIn("cc-number", card.Number, true); err != nil {
		return err
	}
	exp := fmt.Sprintf("%02d / %02d", card.ExpMonth, card.ExpYear%100)
	if got, err := focus("cc-exp"); err != nil {
		return err
	} else if got == "ok" {
		if err := s.c.call(ctx, s.sid, "Input.insertText", map[string]any{"text": exp}, nil); err != nil {
			return errors.New("the checkout's expiry couldn't be typed")
		}
	} else {
		if err := typeIn("cc-exp-month", fmt.Sprintf("%02d", card.ExpMonth), true); err != nil {
			return err
		}
		// Two digits where the field takes two (its maxlength, or a YY
		// placeholder), four otherwise.
		year := fmt.Sprint(card.ExpYear)
		var two struct {
			Result struct {
				Value bool `json:"value"`
			} `json:"result"`
		}
		_ = s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": yearTwoJS, "returnByValue": true}, &two)
		if two.Result.Value {
			year = fmt.Sprintf("%02d", card.ExpYear%100)
		}
		if err := typeIn("cc-exp-year", year, true); err != nil {
			return err
		}
	}
	if err := typeIn("cc-csc", card.CVC, true); err != nil {
		return err
	}
	if card.Name != "" {
		if err := typeIn("cc-name", card.Name, false); err != nil {
			return err
		}
	}
	// Leave the fields: a page's checks run on blur.
	_ = s.c.call(ctx, s.sid, "Runtime.evaluate", map[string]any{"expression": "document.activeElement && document.activeElement.blur && document.activeElement.blur(), ''"}, nil)
	return nil
}

// onePay is the one control on the page that places the order. None, or
// more than one, and nothing is paid: a payment never guesses.
func onePay(p browser.Page) (string, error) {
	var refs []string
	for _, c := range p.Controls {
		if browser.Pays(c) {
			refs = append(refs, c.Ref)
		}
	}
	switch len(refs) {
	case 1:
		return refs[0], nil
	case 0:
		return "", errors.New("the checkout has no control that places the order")
	}
	return "", errors.New("the checkout has more than one control that places the order: nothing is paid")
}

// placed says the page after the click is no longer the checkout: another
// address, and no card field left. A page that still shows the card fields
// (a decline, a check that failed) isn't an order known.
func placed(after, before browser.Page) bool {
	if sameAddress(after.URL, before.URL) {
		return false
	}
	return !slices.ContainsFunc(after.Fields, func(f browser.Field) bool { return f.Kind == "card" })
}

// confirmWords are a confirmation page's words, in the languages tafs
// ships; failWords are a failure's, which win over them.
var (
	confirmWords = regexp.MustCompile(`(?i)\b(thank you for your (order|purchase)|order (is )?(placed|confirmed|received|complete)|your order has been (placed|received|confirmed)|merci pour votre (commande|achat)|commande (confirmée|validée|enregistrée)|vielen dank für ihre bestellung|bestellung (ist )?(eingegangen|bestätigt|abgeschlossen)|gracias por (tu|su) (pedido|compra)|pedido (confirmado|realizado|recibido)|grazie per (il tuo|l'|il vostro) ?(ordine|acquisto)|ordine (confermato|ricevuto|completato)|obrigado pela sua (encomenda|compra)|encomenda confirmada|bedankt voor je (bestelling|aankoop)|bestelling (is )?(geplaatst|ontvangen|bevestigd)|dziękujemy za (zamówienie|zakup)|zamówienie (zostało )?(przyjęte|złożone|potwierdzone))`)
	failWords    = regexp.MustCompile(`(?i)\b(declined|refused|failed|did not go through|didn't go through|unsuccessful|try again|error|refusée?|échoué|abgelehnt|fehlgeschlagen|rechazad[ao]|fallid[ao]|rifiutat[ao]|non riuscit[ao]|recusad[ao]|mislukt|geweigerd|odrzucon[ay]|nieudan[ay])\b`)
)

// confirmed says the page after the click confirms the order in words.
func confirmed(p browser.Page) bool { return confirmWords.MatchString(p.Title + "\n" + p.Text) }

// failed says the page after the click names a failure: never ordered,
// whatever reference or thanks it also shows.
func failed(p browser.Page) bool { return failWords.MatchString(p.Title + "\n" + p.Text) }

// sameAddress: equal once a fragment is dropped; empty is never the same.
func sameAddress(a, b string) bool {
	a, b = strings.SplitN(a, "#", 2)[0], strings.SplitN(b, "#", 2)[0]
	return a != "" && a == b
}

// maxConfirmation bounds the confirmation's text kept.
const maxConfirmation = 500

func confirmation(p browser.Page) string {
	t := strings.TrimSpace(browser.Mask(p, 0).Text)
	if r := []rune(t); len(r) > maxConfirmation {
		t = string(r[:maxConfirmation])
	}
	return t
}

var orderRefRE = regexp.MustCompile(`(?i)\b(?:order|commande|bestellung|pedido|ordine)\s*(?:number|no\.?|nr\.?|n°|nº|#)?\s*[:#]?\s*([A-Z0-9][A-Z0-9-]{2,31})\b`)

// orderRef is the order's reference as the confirmation page words it, or
// empty.
func orderRef(text string) string {
	if m := orderRefRE.FindStringSubmatch(text); m != nil {
		return m[1]
	}
	return ""
}
