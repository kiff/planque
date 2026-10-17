// Package jcs is the JSON Canonicalization Scheme (RFC 8785): the exact
// bytes KIFF hashes into a permit's args_sha256 (RFC 046 D2).
package jcs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"
)

var maxSafeInteger = big.NewInt(1<<53 - 1)

// Canonical returns the RFC 8785 form of the JSON value in b.
func Canonical(b []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, errors.New("jcs: trailing data")
	}
	var out bytes.Buffer
	if err := write(&out, v); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// SHA256 is the lowercase hex SHA-256 of b's canonical form.
func SHA256(b []byte) (string, error) {
	c, err := Canonical(b)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(c)
	return hex.EncodeToString(sum[:]), nil
}

func write(out *bytes.Buffer, v any) error {
	switch v := v.(type) {
	case nil:
		out.WriteString("null")
	case bool:
		out.WriteString(strconv.FormatBool(v))
	case string:
		writeString(out, v)
	case json.Number:
		// Integers beyond 2^53 - 1 alias as doubles, so two different
		// values could hash the same. Refused, as KIFF refuses them
		// (kiff-cloud apps/api/internal/permit/jcs.go).
		if i, ok := new(big.Int).SetString(string(v), 10); ok && new(big.Int).Abs(i).Cmp(maxSafeInteger) > 0 {
			return fmt.Errorf("jcs: integer %s is beyond 2^53 - 1 and cannot be hashed without ambiguity", v)
		}
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return fmt.Errorf("jcs: number %s: %w", v, err)
		}
		s, err := Number(f)
		if err != nil {
			return err
		}
		out.WriteString(s)
	case []any:
		out.WriteByte('[')
		for i, e := range v {
			if i > 0 {
				out.WriteByte(',')
			}
			if err := write(out, e); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		// Properties sort by their UTF-16 code units (RFC 8785 §3.2.3).
		slices.SortFunc(keys, func(a, b string) int {
			return slices.Compare(utf16.Encode([]rune(a)), utf16.Encode([]rune(b)))
		})
		out.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				out.WriteByte(',')
			}
			writeString(out, k)
			out.WriteByte(':')
			if err := write(out, v[k]); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("jcs: unexpected %T", v)
	}
	return nil
}

// writeString escapes only what RFC 8785 §3.2.2.2 requires.
func writeString(out *bytes.Buffer, s string) {
	out.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			out.WriteString(`\"`)
		case '\\':
			out.WriteString(`\\`)
		case '\b':
			out.WriteString(`\b`)
		case '\f':
			out.WriteString(`\f`)
		case '\n':
			out.WriteString(`\n`)
		case '\r':
			out.WriteString(`\r`)
		case '\t':
			out.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(out, `\u%04x`, r)
			} else {
				out.WriteRune(r)
			}
		}
	}
	out.WriteByte('"')
}

// Number formats f as ECMAScript's Number.prototype.toString does, which
// RFC 8785 §3.2.2.3 requires.
func Number(f float64) (string, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "", errors.New("jcs: NaN and Infinity are not JSON")
	}
	if f == 0 {
		return "0", nil
	}
	sign := ""
	if f < 0 {
		sign, f = "-", -f
	}
	// Shortest round-trip digits and the decimal exponent: d.ddd e±n.
	e := strconv.FormatFloat(f, 'e', -1, 64)
	mant, exp, _ := strings.Cut(e, "e")
	digits := strings.Replace(mant, ".", "", 1)
	n, _ := strconv.Atoi(exp)
	n++ // the position of the decimal point after the first digit, as ES counts it
	k := len(digits)
	var s string
	switch {
	case k <= n && n <= 21:
		s = digits + strings.Repeat("0", n-k)
	case 0 < n && n <= 21:
		s = digits[:n] + "." + digits[n:]
	case -6 < n && n <= 0:
		s = "0." + strings.Repeat("0", -n) + digits
	default:
		d := digits[:1]
		if k > 1 {
			d += "." + digits[1:]
		}
		x := n - 1
		if x >= 0 {
			s = d + "e+" + strconv.Itoa(x)
		} else {
			s = d + "e-" + strconv.Itoa(-x)
		}
	}
	return sign + s, nil
}
