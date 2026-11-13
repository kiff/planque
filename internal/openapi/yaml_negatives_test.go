package openapi

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// The negatives for reading a vendor's description in YAML: converted to JSON, then held to the same checks, digest
// and limits as a JSON description. What must never get through: a small file that blows up, a file that is over a
// limit once converted, one description with two digests, or a YAML 1.1 reading that changes what an operation takes.

// descYAML is desc, the JSON description the other tests read, written as a vendor writes YAML: response codes
// unquoted, so they're read as numbers.
const descYAML = `openapi: 3.1.0
info:
  title: Shop API
servers:
- url: https://{region}.shop.example/v1
  variables:
    region:
      default: eu
security:
- key: []
components:
  securitySchemes:
    key:
      type: apiKey
      in: header
      name: X-API-Key
    bearer:
      type: http
      scheme: bearer
  schemas:
    Order:
      type: object
      properties:
        id:
          type: string
        lines:
          type: array
          items:
            $ref: '#/components/schemas/Order'
    Refund:
      type: object
      required:
      - amount
      properties:
        amount:
          type: integer
        reason:
          type: string
  parameters:
    Order:
      name: order
      in: path
      required: true
      schema:
        type: string
paths:
  /orders/{order}:
    parameters:
    - $ref: '#/components/parameters/Order'
    get:
      operationId: getOrder
      summary: Read an order
      parameters:
      - name: expand
        in: query
        schema:
          type: string
      responses:
        200:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Order'
    delete:
      summary: Cancel an order
  /orders/{order}/refunds:
    parameters:
    - $ref: '#/components/parameters/Order'
    post:
      operationId: create-refund
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Refund'
  /notes:
    put:
      operationId: putNotes
      requestBody:
        content:
          application/json:
            schema:
              type: array
              items:
                type: string
  /upload:
    post:
      operationId: upload
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
  /signed:
    get:
      operationId: signed
      parameters:
      - name: X-Signature
        in: header
        required: true
`

func refused(t *testing.T, what string, raw []byte) {
	t.Helper()
	start := time.Now()
	_, err := Parse(raw, key)
	var r *Refused
	if !errors.As(err, &r) {
		t.Errorf("%s: not refused (err %v)", what, err)
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Errorf("%s: took %s to refuse", what, d)
	}
}

// The same description, in YAML or in JSON, is one description: one digest, the same operations, each with the
// same digest. Otherwise an owner approves one and a re-fetch in the other reads as a change, or worse, doesn't.
func TestOneDescriptionInYAMLOrJSONHasOneDigest(t *testing.T) {
	j := parse(t, desc, key)
	y := parse(t, descYAML, key)
	if j.Digest != y.Digest {
		t.Errorf("digest: JSON %s, YAML %s", j.Digest, y.Digest)
	}
	if len(j.Operations) != len(y.Operations) {
		t.Fatalf("operations: JSON %d, YAML %d", len(j.Operations), len(y.Operations))
	}
	for _, o := range j.Operations {
		jd, _ := o.Digest()
		yd, _ := op(t, y, o.Name).Digest()
		if jd != yd {
			t.Errorf("%s: operation digest JSON %s, YAML %s", o.Name, jd, yd)
		}
	}
	// Whitespace and key order are the file's, not the description's.
	if again := parse(t, strings.ReplaceAll(desc, "\n ", "\n   "), key); again.Digest != j.Digest {
		t.Errorf("the same JSON reindented: %s, want %s", again.Digest, j.Digest)
	}
}

// A small file of anchors and aliases that expands without bound is refused, in bounded time, before it's built.
func TestAYAMLAliasBombIsRefused(t *testing.T) {
	var b strings.Builder
	b.WriteString("a0: &a0 [\"lol\",\"lol\",\"lol\",\"lol\",\"lol\",\"lol\",\"lol\",\"lol\",\"lol\"]\n")
	for i := 1; i <= 9; i++ {
		fmt.Fprintf(&b, "a%d: &a%d [*a%d,*a%d,*a%d,*a%d,*a%d,*a%d,*a%d,*a%d,*a%d]\n", i, i, i-1, i-1, i-1, i-1, i-1, i-1, i-1, i-1, i-1)
	}
	bomb := strings.Replace(descYAML, "info:\n  title: Shop API\n", "info:\n  title: Shop API\n  x-bomb:\n"+indent(b.String(), "    "), 1)
	if len(bomb) > 4<<10+len(descYAML) {
		t.Fatalf("the bomb isn't small: %d bytes", len(bomb))
	}
	refused(t, "an alias bomb", []byte(bomb))
	// One alias, even a harmless one, and a merge key (an alias by another name).
	refused(t, "one alias", []byte(strings.Replace(descYAML, "info:\n  title: Shop API\n", "info:\n  title: &t Shop API\n  x-copy: *t\n", 1)))
	refused(t, "a merge key", []byte(strings.Replace(descYAML, "info:\n  title: Shop API\n", "info:\n  title: Shop API\nx-base: &b {a: 1}\nx-more:\n  <<: *b\n", 1)))
}

func indent(s, by string) string {
	return by + strings.ReplaceAll(strings.TrimRight(s, "\n"), "\n", "\n"+by) + "\n"
}

// A YAML file is held to the limits as converted: more operations than planque reads, or a file under 16 MB as
// YAML that is over it as JSON, is refused as a JSON one would be.
func TestAYAMLDescriptionOverALimitOnceConvertedIsRefused(t *testing.T) {
	var ops strings.Builder
	for i := 0; i <= MaxOperations; i++ {
		fmt.Fprintf(&ops, "  /p%d:\n    get:\n      operationId: op%d\n", i, i)
	}
	many := strings.Replace(descYAML, "paths:\n", "paths:\n"+ops.String(), 1)
	refused(t, fmt.Sprintf("%d operations in YAML", MaxOperations+1), []byte(many))

	// Short keys and values in block style: about 6 bytes of YAML make 10 of JSON.
	var big strings.Builder
	big.WriteString(strings.Replace(descYAML, "info:\n  title: Shop API\n", "info:\n  title: Shop API\n  x-pad:\n", 1))
	for i := 0; big.Len() < MaxBytes-(64<<10); i++ {
		fmt.Fprintf(&big, "    k%x: v\n", i)
	}
	if big.Len() > MaxBytes {
		t.Fatalf("the YAML is over the limit itself: %d", big.Len())
	}
	refused(t, "a YAML file over 16 MB once converted", []byte(big.String()))
}

// YAML 1.1 reads unquoted yes, no, on and off as booleans, and big numbers as floats. A description read that way
// takes different arguments than the vendor wrote: an enum of "no" becomes false, a maximum loses its last digits.
// Read as YAML 1.2 and with exact numbers, it's the JSON description, digest and all.
func TestYAMLReadsWhatTheVendorWroteNotYAML11(t *testing.T) {
	jsonVariant := strings.Replace(desc, `"reason": {"type": "string"}`, `"reason": {"type": "string", "enum": ["yes", "no", "on", "off"]}, "cents": {"type": "integer", "maximum": 9007199254740993}`, 1)
	yamlVariant := strings.Replace(descYAML, "        reason:\n          type: string\n", "        reason:\n          type: string\n          enum: [yes, no, on, off]\n        cents:\n          type: integer\n          maximum: 9007199254740993\n", 1)
	if yamlVariant == descYAML {
		t.Fatal("setup: the YAML variant didn't change")
	}
	j, y := parse(t, jsonVariant, key), parse(t, yamlVariant, key)
	jd, _ := op(t, j, "create_refund").Digest()
	yd, _ := op(t, y, "create_refund").Digest()
	if jd != yd || j.Digest != y.Digest {
		t.Errorf("YAML read differently from JSON: operation %s vs %s, description %s vs %s\n%s", yd, jd, y.Digest, j.Digest, op(t, y, "create_refund").InputSchema)
	}
}

// What has no JSON meaning is refused rather than guessed: a tag, a second document, a key given twice.
func TestYAMLWithNoJSONMeaningIsRefused(t *testing.T) {
	for what, yaml := range map[string]string{
		"a tag":         strings.Replace(descYAML, "title: Shop API", "title: !!binary U2hvcCBBUEk=", 1),
		"a custom tag":  strings.Replace(descYAML, "title: Shop API", "title: !vendor Shop API", 1),
		"two documents": descYAML + "---\nopenapi: 3.1.0\n",
		"a key twice":   strings.Replace(descYAML, "info:\n  title: Shop API\n", "info:\n  title: Shop API\n  title: Other API\n", 1),
	} {
		if yaml == descYAML {
			t.Fatalf("setup: %s didn't change the YAML", what)
		}
		refused(t, what, []byte(yaml))
	}
}
