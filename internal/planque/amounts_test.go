package planque

import (
	"encoding/json"
	"math/big"
	"testing"
)

// Each of Stripe's decimals: a currency's smallest unit in its own unit, and
// no figure for a special case or a currency whose decimals aren't known.
func TestAMinorAmountIsSaidInItsCurrencysOwnUnit(t *testing.T) {
	for _, c := range []struct {
		minor    int64
		currency string
		want     string
	}{
		{700, "eur", "7.00 EUR"},
		{100000, "EUR", "1000.00 EUR"},
		{5, "usd", "0.05 USD"},
		{0, "usd", "0.00 USD"},
		{-1250, "gbp", "-12.50 GBP"},
		{700, "jpy", "700 JPY"},
		{700, "mga", "700 MGA"},
		{700, "kwd", "0.700 KWD"},
		{1234560, "bhd", "1234.560 BHD"},
	} {
		places, ok := stripeDecimals(c.currency)
		if !ok {
			t.Errorf("%s: no decimals, want %q", c.currency, c.want)
			continue
		}
		if got := MajorAmount(big.NewInt(c.minor), c.currency, places); got != c.want {
			t.Errorf("%d %s: %q, want %q", c.minor, c.currency, got, c.want)
		}
	}
	for _, none := range []string{"isk", "huf", "twd", "ugx", "xyz", "", "eu", "euro", "BTC", "CLF"} {
		if places, ok := stripeDecimals(none); ok {
			t.Errorf("%q: %d decimals, want none: a special case or an unknown currency is never guessed", none, places)
		}
	}
}

// The Eval's refund and its charge: each amount beside its currency gets its
// major form, at any depth; other numbers, and a figure planque wrote, don't.
func TestAnAnswersAmountsGetTheirMajorForm(t *testing.T) {
	in := `{"id":"re_1","object":"refund","amount":700,"currency":"eur","charge":{"amount":100000,"amount_refunded":700,"currency":"eur","created":1763370000},
"lines":[{"unit_amount":1500,"quantity":3,"currency":"jpy"}],"count":2}`
	var got map[string]any
	if err := json.Unmarshal([]byte(WithMajorAmounts(in, stripeDecimals)), &got); err != nil {
		t.Fatal(err)
	}
	charge := got["charge"].(map[string]any)
	line := got["lines"].([]any)[0].(map[string]any)
	for field, c := range map[string]struct {
		got  any
		want string
	}{
		"amount":                 {got["amount_major"], "7.00 EUR"},
		"charge.amount":          {charge["amount_major"], "1000.00 EUR"},
		"charge.amount_refunded": {charge["amount_refunded_major"], "7.00 EUR"},
		"lines.unit_amount":      {line["unit_amount_major"], "1500 JPY"},
	} {
		if c.got != c.want {
			t.Errorf("%s_major: %v, want %q", field, c.got, c.want)
		}
	}
	for _, none := range []string{"created_major", "quantity_major"} {
		if _, ok := charge[none]; ok {
			t.Errorf("%s: a number that isn't an amount got a major form", none)
		}
		if _, ok := line[none]; ok {
			t.Errorf("%s: a number that isn't an amount got a major form", none)
		}
	}
	if _, ok := got["count_major"]; ok {
		t.Error("count: no currency beside it, so no major form")
	}
}

// An answer with nothing to add is passed on byte for byte: text, an amount
// without a currency, a currency whose decimals aren't known, an amount
// already major (a string, a decimal), a field already there, or what isn't
// one JSON value.
func TestAnAnswerWithNoMinorAmountIsUnchanged(t *testing.T) {
	for _, in := range []string{
		"Refunded.",
		`{"amount": 700}`,
		`{"amount": 700, "currency": "xyz"}`,
		`{"amount": 700, "currency": "isk"}`,
		`{"amount": "7.00", "currency": "eur"}`,
		`{"amount": 7.5, "currency": "eur"}`,
		`{"amount": 700, "currency": "eur", "amount_major": "theirs"}`,
		`[1, 2, 3]`,
		`{"amount": 700, "currency": "eur"} trailing`,
		"",
	} {
		if got := WithMajorAmounts(in, stripeDecimals); got != in {
			t.Errorf("%q came back as %q", in, got)
		}
	}
}

// Only a vendor whose own API says its amounts are minor gets them: Stripe's
// grant does, in its text and its structured content; any other grant's
// answer, whatever it holds, is returned as it was.
func TestOnlyAVendorThatSaysSoGetsMajorAmounts(t *testing.T) {
	answer := json.RawMessage(`{"content":[{"type":"text","text":"{\"amount\":700,\"currency\":\"eur\"}"}],"structuredContent":{"amount":700,"currency":"eur"}}`)
	stripe := &Grant{Request: Request{CatalogEntry: StripeEntry}}
	out, _ := json.Marshal(withMajorAmounts(stripe, answer))
	var got struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Structured map[string]any `json:"structuredContent"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Content[0].Text != `{"amount":700,"amount_major":"7.00 EUR","currency":"eur"}` {
		t.Errorf("text: %s", got.Content[0].Text)
	}
	if got.Structured["amount_major"] != "7.00 EUR" {
		t.Errorf("structured content: %v", got.Structured)
	}
	for _, entry := range []string{"", "qonto", "connector:composio:stripe", "some-api"} {
		other := &Grant{Request: Request{CatalogEntry: entry}}
		if res := withMajorAmounts(other, answer); string(res.(json.RawMessage)) != string(answer) {
			t.Errorf("entry %q: the answer changed: %v", entry, res)
		}
	}
}
