package planque

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/browser"
)

func TestAHandOversWordsArePlainTextOfOneTo280Characters(t *testing.T) {
	for _, ok := range []string{"Fill in the card number, then tap Done", strings.Repeat("é", 280), "Zahlen Sie hier — danke"} {
		if why := handWordsRefusal(ok); why != "" {
			t.Errorf("%q: %s", ok, why)
		}
	}
	for _, bad := range []string{"", strings.Repeat("a", 281), "Pay\nnow", "tab\there", "Done ‮enoD", "a⁦b⁩", "‏text", "x y", string([]byte{0xff})} {
		if handWordsRefusal(bad) == "" {
			t.Errorf("%q was taken", bad)
		}
	}
}

func TestAfterAHandOverNoFieldValueIsKeptWhateverTheExtensionSent(t *testing.T) {
	page := browser.Page{URL: "https://shop.example/pay", Text: "Your card ending 4242", Fields: []browser.Field{
		{Ref: "f1", Kind: "card", Value: "4242 4242 4242 4242"},
		{Ref: "f2", Kind: "checkbox", Value: "on"},
		{Ref: "f3", Kind: "select", Value: "Visa", Options: []string{"Visa", "Mastercard"}},
	}}
	got := blankFields(page)
	for _, f := range got.Fields {
		if f.Value != "" {
			t.Errorf("%s kept %q", f.Ref, f.Value)
		}
	}
	if got.Fields[2].Options[1] != "Mastercard" || got.Text != page.Text {
		t.Fatal("the page's own list or text changed")
	}
	if page.Fields[0].Value == "" {
		t.Fatal("blanking changed the page it was given")
	}
}

// One hand-over's life: waiting, then the first answer kept (a resend changes
// nothing), read once with no values, then gone; only by its own grant.
func TestAHandOverKeepsTheFirstAnswerAndIsReadOnceByItsGrant(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	now := time.Now()
	r.e.s.Now = func() time.Time { return now }
	if err := r.e.s.keepHand(ctx, r.chrome, "st_h1", chromeHand{Tenant: tenant, Grant: "g1", Run: "run1", Words: "Fill in the card", Opened: now}); err != nil {
		t.Fatal(err)
	}
	if h, err := r.e.s.readHand(ctx, r.chrome, "st_h1", "g1"); err != nil || h.Outcome != handWaiting {
		t.Fatalf("before the person: %+v %v", h, err)
	}
	if _, err := r.e.s.readHand(ctx, r.chrome, "st_h1", "g2"); !errors.Is(err, errHandUnknown) {
		t.Fatalf("another grant read it: %v", err)
	}
	page := &browser.Page{URL: "https://shop.example/pay", Fields: []browser.Field{{Ref: "f1", Kind: "card", Value: "4242"}}}
	if err := r.e.s.answerHand(ctx, r.chrome, "st_h1", handDone, page); err != nil {
		t.Fatal(err)
	}
	if err := r.e.s.answerHand(ctx, r.chrome, "st_h1", handDeclined, nil); !errors.Is(err, errHandAnswered) {
		t.Fatalf("a second answer: %v", err)
	}
	if err := r.e.s.answerHand(ctx, "chr_other", "st_h1", handDone, nil); !errors.Is(err, errHandUnknown) {
		t.Fatalf("another Chrome answered it: %v", err)
	}
	h, err := r.e.s.readHand(ctx, r.chrome, "st_h1", "g1")
	if err != nil || h.Outcome != handDone || h.Page == nil || h.Page.Fields[0].Value != "" {
		t.Fatalf("the answer: %+v %v", h, err)
	}
	if _, err := r.e.s.readHand(ctx, r.chrome, "st_h1", "g1"); !errors.Is(err, errHandUnknown) {
		t.Fatalf("read twice: %v", err)
	}
}

// Past its half hour a hand-over takes no answer and reads as expired, once.
func TestAHandOverPastItsHalfHourIsExpired(t *testing.T) {
	r := newRelayRig(t)
	ctx := context.Background()
	now := time.Now()
	r.e.s.Now = func() time.Time { return now }
	_ = r.e.s.keepHand(ctx, r.chrome, "st_h2", chromeHand{Tenant: tenant, Grant: "g1", Opened: now})
	now = now.Add(handOpenFor)
	if err := r.e.s.answerHand(ctx, r.chrome, "st_h2", handDone, nil); !errors.Is(err, errHandExpired) {
		t.Fatalf("an answer after the half hour: %v", err)
	}
	if h, err := r.e.s.readHand(ctx, r.chrome, "st_h2", "g1"); err != nil || h.Outcome != handExpired {
		t.Fatalf("read after: %+v %v", h, err)
	}
	if _, err := r.e.s.readHand(ctx, r.chrome, "st_h2", "g1"); !errors.Is(err, errHandUnknown) {
		t.Fatalf("expired read twice: %v", err)
	}
	if err := r.e.s.answerHand(ctx, r.chrome, "st_h2", "maybe", nil); !errors.Is(err, errHandUnknown) {
		t.Fatalf("an unknown outcome: %v", err)
	}
}
