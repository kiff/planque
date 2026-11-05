package fakes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/kiff/planque/internal/browser"
)

// Browser is a browser worker that drives no browser: a scripted site, the
// test shop at shop.example, with pages as a real worker would read them.
// The end to end on the browser connection runs on it and never touches a
// real merchant (kiff/planque#86). It keeps what each session did, so a
// test can say nothing was sent.
type Browser struct {
	mu sync.Mutex
	// Down: the worker can't be reached; Start fails.
	Down bool
	// Total is the shop's checkout total now; a test changes it between a
	// read and a pay.
	Total string
	// Currency is the shop's; a test sets another to pay in it.
	Currency string
	// Sessions is how many sessions were started.
	Sessions int
	// Steps are every step run, in order, with the page it ran on.
	Steps []BrowserStep
	// Orders are the orders placed: one per click on "Place order".
	Orders []string
	// States are the storage states sessions were started with.
	States [][]byte
}

// BrowserStep is one step a session ran.
type BrowserStep struct {
	Grant string
	URL   string
	Step  browser.Step
}

// NewBrowser is the shop, open, with a total of 2.50 EUR in its cart.
func NewBrowser() *Browser { return &Browser{Total: "2.50", Currency: "EUR"} }

// Sent is how many steps ran: what a vendor would have received.
func (b *Browser) Sent() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.Steps)
}

// The shop's pages.
const (
	ShopItem     = "https://shop.example/item/pen"
	ShopCart     = "https://shop.example/cart"
	ShopCheckout = "https://shop.example/checkout"
	ShopDone     = "https://shop.example/done"
	ShopLogin    = "https://shop.example/login"
	ShopSearch   = "https://shop.example/search"
	ShopResults  = "https://shop.example/results"
	OtherSite    = "https://other.example/"
	// Leaving is a page on the shop whose link goes to another site.
	Leaving = "https://shop.example/partners"
	// Away is a page on the shop that redirects to another site.
	Away = "https://shop.example/away"
	// ShopCheckoutFR is the checkout in French: no English word on it.
	ShopCheckoutFR = "https://shop.example/commande"
	// ShopItemDE is a product page with a one-click buy, in German.
	ShopItemDE = "https://shop.example/artikel/stift"
)

func (b *Browser) page(u string) (browser.Page, bool) {
	lines := []browser.Line{{Item: "Blue pen", Quantity: 2, Price: "1.25"}}
	switch u {
	case ShopItem:
		return browser.Page{URL: u, Title: "Blue pen", Text: "Blue pen. 1.25 EUR each. In stock.",
			Fields:   []browser.Field{{Ref: "qty", Label: "Quantity", Kind: "number", Value: "1"}},
			Controls: []browser.Control{{Ref: "add", Label: "Add to cart", Submits: true}},
			Links:    []browser.Link{{Ref: "cart", Label: "Cart", URL: ShopCart}, {Ref: "login", Label: "Sign in", URL: ShopLogin}}}, true
	case ShopCart:
		return browser.Page{URL: u, Title: "Your cart", Text: "Blue pen × 2: 2.50 EUR",
			Links: []browser.Link{{Ref: "checkout", Label: "Checkout", URL: ShopCheckout}}}, true
	case ShopCheckout:
		return browser.Page{URL: u, Title: "Checkout", Text: "Blue pen × 2. Total " + b.Total + " " + b.Currency,
			Fields: []browser.Field{
				{Ref: "email", Label: "Email", Kind: "email"},
				{Ref: "card", Label: "Card number", Kind: "card", Value: "4242 4242 4242 4242"},
				{Ref: "cvc", Label: "Security code", Kind: "code", Value: "123"},
				{Ref: "coupon", Label: "Coupon", Kind: "text"}},
			Controls: []browser.Control{{Ref: "apply", Label: "Apply coupon", Submits: true}, {Ref: "place", Label: "Place order", Submits: true, PlacesOrder: true}},
			Checkout: &browser.Checkout{Lines: lines, Total: b.Total, Currency: b.Currency}}, true
	case ShopDone:
		return browser.Page{URL: u, Title: "Thank you", Text: "Order 1001 placed. Total " + b.Total + " " + b.Currency + "."}, true
	case ShopLogin:
		return browser.Page{URL: u, Title: "Sign in",
			Fields:   []browser.Field{{Ref: "user", Label: "Email", Kind: "email"}, {Ref: "pass", Label: "Password", Kind: "password", Value: "hunter2"}},
			Controls: []browser.Control{{Ref: "signin", Label: "Sign in", Submits: true}}}, true
	case ShopSearch:
		return browser.Page{URL: u, Title: "Search",
			Fields:   []browser.Field{{Ref: "q", Label: "Search", Kind: "search"}},
			Controls: []browser.Control{{Ref: "go", Label: "Search", Submits: true}}}, true
	case ShopResults:
		return browser.Page{URL: u, Title: "Results", Text: "1 result: Blue pen", Links: []browser.Link{{Ref: "pen", Label: "Blue pen", URL: ShopItem}}}, true
	case Leaving:
		return browser.Page{URL: u, Title: "Partners", Links: []browser.Link{{Ref: "other", Label: "Our partner", URL: OtherSite}}}, true
	case Away:
		// The browser follows the redirect: the page is the other site's.
		return browser.Page{URL: OtherSite, Title: "Other", Text: "Another site.",
			Checkout: &browser.Checkout{Lines: lines, Total: b.Total, Currency: b.Currency},
			Controls: []browser.Control{{Ref: "place", Label: "Place order", Submits: true, PlacesOrder: true}}}, true
	case ShopCheckoutFR:
		return browser.Page{URL: u, Title: "Commande", Text: "Stylo bleu × 2. Total " + b.Total + " " + b.Currency,
			Fields:   []browser.Field{{Ref: "email", Label: "Courriel", Kind: "email"}},
			Controls: []browser.Control{{Ref: "cmd", Label: "Commander", Submits: true}, {Ref: "icon", Label: "", Submits: false}, {Ref: "zb", Label: "Jetzt kaufen"}},
			Checkout: &browser.Checkout{Lines: lines, Total: b.Total, Currency: b.Currency}}, true
	case ShopItemDE:
		return browser.Page{URL: u, Title: "Blauer Stift", Text: "Blauer Stift. 1,25 EUR.",
			Controls: []browser.Control{{Ref: "kaufen", Label: "Kaufen"}, {Ref: "merken", Label: "Merken"}}}, true
	case OtherSite:
		return browser.Page{URL: u, Title: "Other", Text: "Another site."}, true
	}
	return browser.Page{}, false
}

// Start opens a session on the shop.
func (b *Browser) Start(_ context.Context, start browser.Start) (browser.Session, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Down {
		return nil, errors.New("the browser worker is down")
	}
	b.Sessions++
	b.States = append(b.States, start.State)
	return &browserSession{b: b, grant: start.Grant, sites: start.Sites, state: start.State}, nil
}

type browserSession struct {
	b      *Browser
	grant  string
	sites  []string
	state  []byte
	page   browser.Page
	opened bool
	closed bool
}

func (s *browserSession) Open(_ context.Context, u string) (browser.Page, error) {
	s.b.mu.Lock()
	defer s.b.mu.Unlock()
	if s.closed {
		return browser.Page{}, errors.New("the session is closed")
	}
	p, ok := s.b.page(u)
	if !ok {
		return browser.Page{URL: u, Title: "Not found", Text: "404"}, nil
	}
	s.page, s.opened = p, true
	return p, nil
}

func (s *browserSession) Do(_ context.Context, step browser.Step) (browser.Page, error) {
	s.b.mu.Lock()
	defer s.b.mu.Unlock()
	if s.closed || !s.opened {
		return browser.Page{}, errors.New("no page is open")
	}
	s.b.Steps = append(s.b.Steps, BrowserStep{Grant: s.grant, URL: s.page.URL, Step: step})
	ref := step.Ref()
	if step.Fill != "" || step.Select != "" {
		for i := range s.page.Fields {
			if s.page.Fields[i].Ref == ref {
				if step.Fill != "" {
					s.page.Fields[i].Value = step.Value
				} else {
					s.page.Fields[i].Value = step.Option
				}
				return s.page, nil
			}
		}
		return browser.Page{}, fmt.Errorf("no field %s", ref)
	}
	for _, l := range s.page.Links {
		if l.Ref == ref {
			p, _ := s.b.page(l.URL)
			s.page = p
			return p, nil
		}
	}
	next := ""
	switch s.page.URL + "#" + ref {
	case ShopItem + "#add":
		next = ShopCart
	case ShopCheckout + "#place", ShopCheckoutFR + "#cmd", ShopCheckoutFR + "#icon", ShopCheckoutFR + "#zb", ShopItemDE + "#kaufen", OtherSite + "#place":
		s.b.Orders = append(s.b.Orders, "1001")
		next = ShopDone
	case ShopItemDE + "#merken":
		next = ShopItemDE
	case ShopCheckout + "#apply":
		next = ShopCheckout
	case ShopLogin + "#signin":
		s.state = []byte(`{"cookies":[{"name":"session","value":"signed-in"}]}`)
		next = ShopItem
	case ShopSearch + "#go":
		next = ShopResults
	default:
		return browser.Page{}, fmt.Errorf("no control %s", ref)
	}
	p, _ := s.b.page(next)
	s.page = p
	return p, nil
}

func (s *browserSession) Screenshot(context.Context) ([]byte, error) {
	return []byte("\x89PNG " + s.page.URL), nil
}

func (s *browserSession) Close(context.Context) ([]byte, error) {
	s.b.mu.Lock()
	defer s.b.mu.Unlock()
	s.closed = true
	if s.state == nil {
		return []byte(`{"cookies":[]}`), nil
	}
	return s.state, nil
}

// Evidence keeps screenshots in memory, as planque's bucket does at rest.
type Evidence struct {
	mu   sync.Mutex
	Kept map[string][]byte // tenant/grant/op → png
}

func NewEvidence() *Evidence { return &Evidence{Kept: map[string][]byte{}} }

func (e *Evidence) Put(_ context.Context, tenant, grant, op string, png []byte) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ref := tenant + "/" + grant + "/" + op
	e.Kept[ref] = png
	return ref, nil
}

func (e *Evidence) Drop(_ context.Context, tenant, grant string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for ref := range e.Kept {
		if strings.HasPrefix(ref, tenant+"/"+grant+"/") {
			delete(e.Kept, ref)
		}
	}
	return nil
}

// Refs lists what Evidence holds for a grant.
func (e *Evidence) Refs(tenant, grant string) []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	var out []string
	for ref := range e.Kept {
		if strings.HasPrefix(ref, tenant+"/"+grant+"/") {
			out = append(out, ref)
		}
	}
	return out
}

// StepsJSON is the steps run, as JSON, for a test's words.
func (b *Browser) StepsJSON() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	raw, _ := json.Marshal(b.Steps)
	return string(raw)
}
