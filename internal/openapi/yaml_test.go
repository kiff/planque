package openapi

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.yaml.in/yaml/v3"
)

// asYAML is a JSON description written as YAML, as a vendor publishing YAML
// only would: response codes unquoted.
func asYAML(t *testing.T, raw string) []byte {
	t.Helper()
	var doc any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatal(err)
	}
	b, err := yaml.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	return []byte(strings.ReplaceAll(string(b), `"200":`, `200:`))
}

// A vendor that publishes its description in YAML only gives the same tools,
// and the same digest, as the same description in JSON.
func TestAYAMLDescriptionReadsAsItsJSON(t *testing.T) {
	raw := asYAML(t, desc)
	if !strings.Contains(string(raw), "200:") || strings.Contains(string(raw), `"200":`) {
		t.Fatalf("the response code isn't unquoted:\n%s", raw)
	}
	fromJSON, fromYAML := parse(t, desc, key), parse(t, string(raw), key)
	a, _ := json.Marshal(fromJSON)
	b, _ := json.Marshal(fromYAML)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("YAML read otherwise:\n%s\n%s", a, b)
	}
	// The same JSON indented otherwise has the same digest too.
	var doc any
	_ = json.Unmarshal([]byte(desc), &doc)
	flat, _ := json.Marshal(doc)
	if parse(t, string(flat), key).Digest != fromJSON.Digest {
		t.Fatal("reindented JSON has another digest")
	}
	// Selecting operations reads YAML too.
	if s, err := ParseOperations(raw, key, []string{"getOrder"}); err != nil || len(s.Operations) != 1 || s.Digest != fromJSON.Digest {
		t.Fatalf("%+v %v", s, err)
	}
}

// YAML's own spellings that JSON reads otherwise are read as JSON would:
// yes and no are words, and a large integer keeps its every digit.
func TestYAMLKeepsWordsAndNumbersExact(t *testing.T) {
	twin := `{"openapi": "3.0.3", "components": {"schemas": {"Answer": {"type": "string", "enum": ["yes", "no", "on", "off"]}, "Big": {"type": "integer", "maximum": 9007199254740993}}}, "paths": {}}`
	y := "openapi: 3.0.3\ncomponents:\n  schemas:\n    Answer:\n      type: string\n      enum: [yes, no, on, off]\n    Big:\n      type: integer\n      maximum: 9007199254740993\npaths: {}\n"
	j, err := yamlAsJSON([]byte(y))
	if err != nil {
		t.Fatal(err)
	}
	var got, want any
	dec := json.NewDecoder(strings.NewReader(string(j)))
	dec.UseNumber()
	_ = dec.Decode(&got)
	dec = json.NewDecoder(strings.NewReader(twin))
	dec.UseNumber()
	_ = dec.Decode(&want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("read otherwise:\n%s\n%s", j, twin)
	}
}

// An unquoted date, as many vendors write their API's version, is the text
// written, with the digest of its quoted twin.
func TestAnUnquotedDateIsItsText(t *testing.T) {
	quoted := strings.Replace(string(asYAML(t, desc)), "title: Shop API", "title: Shop API\n    version: '2018-02-08'", 1)
	bare := strings.Replace(quoted, "'2018-02-08'", "2018-02-08", 1)
	if bare == quoted || !strings.Contains(bare, "version: 2018-02-08\n") {
		t.Fatalf("the date isn't in:\n%s", bare)
	}
	j, err := yamlAsJSON([]byte(bare))
	if err != nil || !strings.Contains(string(j), `"version":"2018-02-08"`) {
		t.Fatalf("%s %v", j, err)
	}
	if parse(t, bare, key).Digest != parse(t, quoted, key).Digest {
		t.Fatal("the unquoted date has another digest")
	}
}

// What isn't plain data is refused with the reason, as JSON's faults are.
func TestYAMLThatIsntPlainDataIsRefused(t *testing.T) {
	head := "openapi: 3.0.0\npaths: {}\n"
	for name, c := range map[string]struct{ raw, why string }{
		"not yaml":      {"openapi: [3.0.0\n", "isn't valid JSON or YAML"},
		"a list":        {"- openapi: 3.0.0\n", "isn't an OpenAPI description"},
		"a scalar":      {"openapi\n", "isn't an OpenAPI description"},
		"swagger 2":     {"swagger: '2.0'\n", "isn't OpenAPI 3.0 or 3.1"},
		"blank":         {"  \n", "it is empty"},
		"an alias":      {head + "x: &a {k: 1}\ny: *a\n", "anchors or aliases"},
		"an anchor":     {head + "x: &a {k: 1}\n", "anchors or aliases"},
		"a merge key":   {head + "x:\n  <<: {k: 1}\n", "merge keys"},
		"binary":        {head + "x: !!binary aGk=\n", "tag !!binary"},
		"a custom tag":  {head + "x: !vendor 1\n", "tag !vendor"},
		"two documents": {head + "---\nopenapi: 3.0.0\n", "more than one YAML document"},
		"a key twice":   {head + "x: 1\nx: 2\n", "given twice"},
		"a list as key": {head + "? [a, b]\n: 1\n", "isn't a plain string"},
		"infinity":      {head + "x: .inf\n", "no JSON form"},
	} {
		_, err := Parse([]byte(c.raw), key)
		var r *Refused
		if !errors.As(err, &r) || !strings.Contains(err.Error(), c.why) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// A small file can't stand for a huge one: nested aliases are refused at the
// first, fast, and a YAML file within 16 MB that is over it as JSON is refused.
func TestYAMLCantGrowPastTheLimits(t *testing.T) {
	var bomb strings.Builder
	bomb.WriteString("openapi: 3.0.0\npaths: {}\na0: &a0 [x, x, x, x, x, x, x, x, x, x]\n")
	for i := 1; i < 10; i++ {
		prev := "*a" + string(rune('0'+i-1))
		bomb.WriteString("a" + string(rune('0'+i)) + ": &a" + string(rune('0'+i)) + " [" + strings.Repeat(prev+", ", 9) + prev + "]\n")
	}
	start := time.Now()
	if _, err := Parse([]byte(bomb.String()), key); err == nil || !strings.Contains(err.Error(), "anchors or aliases") || time.Since(start) > 2*time.Second {
		t.Fatalf("an alias bomb: %v in %s", err, time.Since(start))
	}
	// One string JSON writes longer: each "<" is six bytes as JSON (\u003c).
	big := "openapi: 3.0.0\npaths: {}\nx: '" + strings.Repeat("<", 4<<20) + "'\n"
	if _, err := Parse([]byte(big), key); err == nil || !strings.Contains(err.Error(), "over 16 MB as JSON") {
		t.Fatalf("over the limit as JSON: %v", err)
	}
}
