package jcs

import (
	"math"
	"testing"
)

func TestCanonical(t *testing.T) {
	cases := []struct{ in, want string }{
		// RFC 8785 §3.2.2 example, sorted and with ES number forms.
		{`{"numbers":[333333333.33333329,1E30,4.50,2e-3,0.000000000000000000000000001],"string":"\u20ac$\u000F\u000aA'\u0042\u0022\u005c\\\"\/","literals":[null,true,false]}`,
			`{"literals":[null,true,false],"numbers":[333333333.3333333,1e+30,4.5,0.002,1e-27],"string":"€$\u000f\nA'B\"\\\\\"/"}`},
		// Keys sort by UTF-16 code units, not bytes (RFC 8785 §3.2.3).
		{`{"\ud83d\ude00":1,"\ufb33":2,"a":3}`, `{"a":3,"😀":1,"דּ":2}`},
		{`{"b":{"z":1,"a":[1,2]},"a":"x"}`, `{"a":"x","b":{"a":[1,2],"z":1}}`},
	}
	for _, c := range cases {
		got, err := Canonical([]byte(c.in))
		if err != nil {
			t.Fatalf("%s: %v", c.in, err)
		}
		if string(got) != c.want {
			t.Errorf("Canonical(%s)\n got %s\nwant %s", c.in, got, c.want)
		}
	}
}

func TestNumber(t *testing.T) {
	cases := map[float64]string{
		0: "0", 1: "1", -1: "-1", 100: "100", 1e21: "1e+21", 1e20: "100000000000000000000",
		0.000001: "0.000001", 0.0000001: "1e-7", 123.456: "123.456", 9007199254740991: "9007199254740991",
		5e-324: "5e-324", math.MaxFloat64: "1.7976931348623157e+308",
	}
	for f, want := range cases {
		if got, _ := Number(f); got != want {
			t.Errorf("Number(%v) = %s, want %s", f, got, want)
		}
	}
}
