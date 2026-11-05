// Package browser is the browser as a connection (kiff/planque#86): what a
// page is to planque, the worker that drives the real browser, and the
// guard planque applies to every step before the worker runs it. The
// worker reports facts about a page (its fields, their kinds, its
// controls); planque decides. Nothing here is site-specific.
package browser

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"regexp"
	"strings"

	"github.com/kiff/planque/internal/egress"
)

// Page is a page as the worker reads it: the facts planque and the agent
// work from. Values of secret fields are masked by Mask before an agent
// sees them, whatever the worker sent.
type Page struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	// Text is the page's readable text, cut by the worker at a size.
	Text     string    `json:"text,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
	Controls []Control `json:"controls,omitempty"`
	Links    []Link    `json:"links,omitempty"`
	// PaymentFrame: the page embeds a payment provider's frame.
	PaymentFrame bool `json:"payment_frame,omitempty"`
	// Checkout is what the worker reads on a checkout page: the lines, the
	// total and its currency. Nil when the page isn't one it can read.
	Checkout *Checkout `json:"checkout,omitempty"`
}

// Field is a form field, with a stable ref the agent names it by.
type Field struct {
	Ref   string `json:"ref"`
	Label string `json:"label"`
	// Kind is the field's kind as the worker reads it from its type and
	// autocomplete: text, email, tel, number, search, select, checkbox,
	// textarea, password, card, code (a one-time code), or other.
	Kind    string   `json:"kind"`
	Value   string   `json:"value,omitempty"`
	Options []string `json:"options,omitempty"`
	Masked  bool     `json:"masked,omitempty"`
}

// Control is a button or a submit, with a stable ref.
type Control struct {
	Ref   string `json:"ref"`
	Label string `json:"label"`
	// Submits: the control submits a form.
	Submits bool `json:"submits,omitempty"`
	// PlacesOrder: the worker read, from the page's and the control's
	// semantics (a checkout form's submit, structured data, a payment
	// provider's button), that the control places an order or pays. A
	// fact in no language; the words are the fallback.
	PlacesOrder bool `json:"places_order,omitempty"`
}

// Link is a link on the page.
type Link struct {
	Ref   string `json:"ref"`
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Checkout is a checkout page as read: the lines, the total and the
// currency, as the page shows them. Total is a decimal string ("12.50").
type Checkout struct {
	Lines    []Line `json:"lines"`
	Total    string `json:"total"`
	Currency string `json:"currency"`
}

// Line is one item in a checkout.
type Line struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

// Same says two checkouts show the same lines, total and currency.
func (c Checkout) Same(o Checkout) bool {
	if c.Total != o.Total || c.Currency != o.Currency || len(c.Lines) != len(o.Lines) {
		return false
	}
	for i := range c.Lines {
		if c.Lines[i] != o.Lines[i] {
			return false
		}
	}
	return true
}

// Step is one thing the agent asks the browser to do on a page: fill a
// field, select an option, or click a control. Exactly one is set.
type Step struct {
	Fill   string `json:"fill,omitempty"`
	Value  string `json:"value,omitempty"`
	Select string `json:"select,omitempty"`
	Option string `json:"option,omitempty"`
	Click  string `json:"click,omitempty"`
}

// Ref is the field or control the step names.
func (s Step) Ref() string {
	switch {
	case s.Fill != "":
		return s.Fill
	case s.Select != "":
		return s.Select
	}
	return s.Click
}

// Check says the step is well formed: one of fill, select or click.
func (s Step) Check() error {
	n := 0
	for _, r := range []string{s.Fill, s.Select, s.Click} {
		if r != "" {
			n++
		}
	}
	switch {
	case n != 1:
		return errors.New("a step is one of {fill, value}, {select, option} or {click}")
	case s.Fill == "" && s.Value != "", s.Select == "" && s.Option != "":
		return errors.New("value goes with fill, option with select")
	case s.Select != "" && s.Option == "":
		return errors.New("select names the option to pick")
	case len(s.Value) > 4096:
		return errors.New("a value is at most 4,096 characters")
	}
	return nil
}

// Start is what a session begins with: the grant's kept storage state
// (opened by planque, sealed at rest) and the sites it may reach.
type Start struct {
	Grant string
	State []byte
	Sites []string
}

// Worker drives a real browser, in a session per call, and keeps nothing
// between sessions.
type Worker interface {
	Start(ctx context.Context, start Start) (Session, error)
}

// Session is one call's browser. Close returns the storage state after,
// for planque to seal, and ends the browser.
type Session interface {
	Open(ctx context.Context, url string) (Page, error)
	Do(ctx context.Context, step Step) (Page, error)
	Screenshot(ctx context.Context) ([]byte, error)
	Close(ctx context.Context) ([]byte, error)
}

// Evidence keeps a screenshot with an operation's record, encrypted at
// rest, for as long as the call records: deleted with the grant and the
// account. Put returns a reference the record keeps.
type Evidence interface {
	Put(ctx context.Context, tenant, grant, op string, png []byte) (string, error)
	Drop(ctx context.Context, tenant, grant string) error
}

// The codes a browser call is refused with, beyond the relay's own.
const (
	CodeSiteNotListed   = "site_not_listed"     // the URL's site isn't one the grant lists
	CodeAddressPrivate  = "address_private"     // a literal or local address, never a site
	CodeBadURL          = "bad_url"             // not an https URL planque can open
	CodePaymentStep     = "payment_step"        // a step that would pay or place an order: only browse_pay may
	CodeSecretField     = "secret_field"        // a step typing into a password, card or code field
	CodeRefUnknown      = "ref_unknown"         // the step names no field or control on the page
	CodeTotalUnread     = "total_unread"        // no total and currency read on the checkout
	CodeCheckoutChanged = "checkout_changed"    // the checkout isn't what was read before
	CodeCurrencyOther   = "currency_other"      // the checkout isn't in the connection's currency
	CodeBrowserDown     = "browser_unavailable" // the worker couldn't run the browser
)

// Secret says the field's value is never shown and never typed by an
// agent: a password, a card number or security code, a one-time code.
func Secret(f Field) bool {
	switch f.Kind {
	case "password", "card", "code":
		return true
	}
	return false
}

// Mask is p with every secret field's value removed and said masked, a
// card number shown anywhere in its text or a field's value replaced by
// its last four digits, and its text cut at a size. The worker masks too;
// this holds whatever it sent.
func Mask(p Page, maxText int) Page {
	out := p
	out.Title = MaskText(p.Title)
	out.Text = MaskText(p.Text)
	out.Fields = make([]Field, len(p.Fields))
	for i, f := range p.Fields {
		if Secret(f) {
			f.Value, f.Masked = "", true
		} else {
			f.Value = MaskText(f.Value)
		}
		out.Fields[i] = f
	}
	if maxText > 0 && len(out.Text) > maxText {
		out.Text = out.Text[:maxText] + "\n(cut)"
	}
	return out
}

// cardRun is a run of 13 to 19 digits, spaces or dashes between groups
// allowed, as a card number is written.
var cardRun = regexp.MustCompile(`\b(?:\d[ -]?){12,18}\d\b`)

// MaskText is text with every card number in it masked to its last four
// digits: a run of 13 to 19 digits that passes the Luhn check, with or
// without spaces or dashes. An order number or a phone number doesn't
// pass, as a rule; one that does is masked all the same, which costs
// nothing. A one-time code can't be told from any other number, so the
// page's own masking of its field is what keeps it out.
func MaskText(s string) string {
	if s == "" {
		return s
	}
	return cardRun.ReplaceAllStringFunc(s, func(run string) string {
		digits := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, run)
		if len(digits) < 13 || len(digits) > 19 || !luhn(digits) {
			return run
		}
		return "•••• " + digits[len(digits)-4:]
	})
}

// luhn says digits pass the Luhn check, as every card number does.
func luhn(digits string) bool {
	sum, double := 0, false
	for i := len(digits) - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

// payWords are the words of a control that pays or places an order, in
// the languages of the sites tafs use most (en, fr, de, es, it, nl, pt).
// The words are the fallback behind the worker's PlacesOrder fact, and
// the detection errs to refusing: a page that only looks like it pays is
// one browse_act won't submit.
var payWords = regexp.MustCompile(`(?i)\b(` +
	// en
	`pay|pay now|place (your |the |my )?order|buy( it)?( now)?|complete (the |your |my )?(purchase|order|payment)|confirm (and pay|purchase|order|payment)|purchase|subscribe|checkout now|order now|` +
	// fr
	`commander|acheter|payer|payer maintenant|(confirmer|valider|passer|finaliser) (la |ma )?commande|r[ée]gler|` +
	// de
	`kaufen|jetzt kaufen|bestellen|bezahlen|jetzt bezahlen|(zahlungspflichtig|kostenpflichtig) bestellen|bestellung (abschicken|absenden|aufgeben)|kauf abschlie[sß]+en|` +
	// es
	`comprar|comprar ahora|pagar|pagar ahora|realizar (el )?pedido|confirmar (el |la )?(pedido|compra|pago)|finalizar (la )?compra|tramitar pedido|` +
	// it
	`acquista( ora)?|compra( ora)?|paga( ora)?|ordina( ora)?|(conferma|completa) (l'?)?(ordine|acquisto|pagamento)|procedi al pagamento|` +
	// nl
	`kopen|nu kopen|koop nu|betalen|nu betalen|afrekenen|bestelling (plaatsen|bevestigen)|` +
	// pt
	`comprar agora|pagar agora|(finalizar|concluir) (a )?compra|confirmar (o |a )?(pedido|compra|pagamento)|efetuar pagamento|fazer pedido` +
	`)\b`)

// Pays says the control places an order or pays: the worker read it from
// the page, or its words say so.
func Pays(c Control) bool { return c.PlacesOrder || payWords.MatchString(c.Label) }

// Payment says p is a payment page: it has a card or security-code field,
// embeds a payment frame, or has a control whose words place an order.
func Payment(p Page) bool {
	if p.PaymentFrame {
		return true
	}
	for _, f := range p.Fields {
		if f.Kind == "card" {
			return true
		}
	}
	for _, c := range p.Controls {
		if Pays(c) {
			return true
		}
	}
	return false
}

// Refuse says why step may not run on p: its code and words, or "" when
// it may. A fill or select of a secret field, a click that pays, any
// click on a control of a checkout or payment page (labelled or not,
// submitting or not: a button's effect can't be read from outside), and a
// ref the page doesn't have are refused. Links are followed.
func Refuse(p Page, step Step) (code, why string) {
	ref := step.Ref()
	if step.Click != "" {
		for _, c := range p.Controls {
			if c.Ref != ref {
				continue
			}
			switch {
			case Pays(c):
				return CodePaymentStep, "clicking " + quote(c.Label) + " would pay or place an order; only browse_pay may, held for the owner"
			case p.Checkout != nil || Payment(p):
				return CodePaymentStep, "on a checkout or payment page the taf clicks no control (" + quote(c.Label) + "); only browse_pay may, held for the owner"
			}
			return "", ""
		}
		for _, l := range p.Links {
			if l.Ref == ref {
				return "", ""
			}
		}
		return CodeRefUnknown, "the page has no control or link " + ref
	}
	for _, f := range p.Fields {
		if f.Ref != ref {
			continue
		}
		switch {
		case f.Kind == "card":
			return CodePaymentStep, "typing into " + quote(f.Label) + " would enter card details; the taf never types those"
		case Secret(f):
			return CodeSecretField, "typing into " + quote(f.Label) + " would enter a password or a code; the taf never types those"
		}
		return "", ""
	}
	return CodeRefUnknown, "the page has no field " + ref
}

func quote(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "that control"
	}
	return fmt.Sprintf("%q", s)
}

var siteName = regexp.MustCompile(`^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)

// CheckSite says name is a site a grant may list: a host name, lower
// case, no scheme, port or path, and not a local one.
func CheckSite(name string) error {
	if !siteName.MatchString(name) {
		return fmt.Errorf("%q isn't a site: a host name such as shop.example, no scheme, port or path", name)
	}
	if local(name) {
		return fmt.Errorf("%q isn't a public site", name)
	}
	return nil
}

// local says host is one a browser never goes to on a taf's behalf: an
// address written out, or a name that is local by convention. A name that
// resolves to a private address is refused by the worker's egress, which
// sees the address (egress.Public).
func local(host string) bool {
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if ip, err := netip.ParseAddr(strings.Trim(host, "[]")); err == nil {
		return !egress.Public(ip)
	}
	if host == "localhost" || !strings.Contains(host, ".") {
		return true
	}
	for _, suffix := range []string{".localhost", ".local", ".internal", ".home", ".lan", ".corp", ".intranet", ".arpa"} {
		if strings.HasSuffix(host, suffix) {
			return true
		}
	}
	return false
}

// Allowed says why raw isn't a URL the grant may open: its code and words,
// or "" when it is. sites, when given, are the only hosts (and their
// subdomains) the grant reaches.
func Allowed(raw string, sites []string) (code, why string) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return CodeBadURL, "the URL must be https, with a host and no user info"
	}
	host := strings.TrimSuffix(strings.ToLower(u.Hostname()), ".")
	if _, err := netip.ParseAddr(strings.Trim(host, "[]")); err == nil {
		return CodeAddressPrivate, "the browser opens sites by name, never an address"
	}
	if local(host) {
		return CodeAddressPrivate, host + " is a local name, not a site"
	}
	if len(sites) == 0 {
		return "", ""
	}
	for _, s := range sites {
		if host == s || strings.HasSuffix(host, "."+s) {
			return "", ""
		}
	}
	return CodeSiteNotListed, host + " isn't one of the sites this connection may use (" + strings.Join(sites, ", ") + ")"
}

// Site is the host a URL names, lower case; "" when it isn't a URL.
func Site(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}
