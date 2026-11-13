package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"regexp"
	"strconv"

	"go.yaml.in/yaml/v3"
)

// jsonNumber is JSON's number grammar: a YAML number written so is kept
// as written, digit for digit.
var jsonNumber = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)

// yamlAsJSON reads a YAML description as the JSON it stands for, so it is
// checked, digested and limited as a JSON one is. Only plain data is read:
// one document of mappings, sequences and scalars. Anchors, aliases and merge
// keys are refused, so a small file can't stand for a huge one, and so are
// tags beyond YAML's own, which JSON has no form for. A key is a string, as
// in JSON: an unquoted 200 is "200".
func yamlAsJSON(raw []byte) ([]byte, error) {
	// The parser builds every node before any is counted, and a file of
	// short keys within 16 MB is millions of them: seconds of work. A cheap
	// count first: each "key: value" is at most two nodes, each list item
	// one, and a colon or comma inside a string only counts more. A real
	// description is far below (Zendesk's 2.5 MB: about 140,000).
	if 2*bytes.Count(raw, []byte(":"))+bytes.Count(raw, []byte("- "))+bytes.Count(raw, []byte(",")) > 2*maxNodes {
		return nil, refuse("it is too large once read")
	}
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	var doc yaml.Node
	if err := dec.Decode(&doc); err != nil {
		return nil, refuse("it isn't valid JSON or YAML")
	}
	var more yaml.Node
	if err := dec.Decode(&more); !errors.Is(err, io.EOF) {
		return nil, refuse("it has more than one YAML document")
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, refuse("it isn't an OpenAPI description")
	}
	nodes := 0
	v, err := yamlValue(doc.Content[0], &nodes)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, refuse("it can't be read as JSON")
	}
	if len(b) > MaxBytes {
		return nil, refuse("it is over 16 MB as JSON")
	}
	return b, nil
}

func yamlValue(n *yaml.Node, nodes *int) (any, error) {
	if *nodes++; *nodes > maxNodes {
		return nil, refuse("it is too large once read")
	}
	if n.Anchor != "" || n.Kind == yaml.AliasNode {
		return nil, refuse("it uses YAML anchors or aliases, which planque doesn't read")
	}
	switch n.Kind {
	case yaml.MappingNode:
		out := make(map[string]any, len(n.Content)/2)
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := n.Content[i]
			if k.Kind != yaml.ScalarNode || k.Anchor != "" {
				return nil, refuse("a YAML key isn't a plain string")
			}
			switch k.Tag {
			case "!!str", "!!int", "!!float", "!!bool", "!!null", "!!timestamp":
			case "!!merge":
				return nil, refuse("it uses YAML merge keys, which planque doesn't read")
			default:
				return nil, refuse("a YAML key has the tag %s, which planque doesn't read", k.Tag)
			}
			if _, dup := out[k.Value]; dup {
				return nil, refuse("the key %q is given twice", k.Value)
			}
			v, err := yamlValue(n.Content[i+1], nodes)
			if err != nil {
				return nil, err
			}
			out[k.Value] = v
		}
		return out, nil
	case yaml.SequenceNode:
		out := make([]any, 0, len(n.Content))
		for _, c := range n.Content {
			v, err := yamlValue(c, nodes)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case yaml.ScalarNode:
		return yamlScalar(n)
	}
	return nil, refuse("it isn't plain YAML data")
}

func yamlScalar(n *yaml.Node) (any, error) {
	switch n.Tag {
	case "!!str", "!!timestamp":
		// An unquoted date (version: 2018-02-08) is the text written, as JSON carries it.
		return n.Value, nil
	case "!!null":
		return nil, nil
	case "!!bool":
		var b bool
		if err := n.Decode(&b); err != nil {
			return nil, refuse("it isn't valid YAML")
		}
		return b, nil
	case "!!int", "!!float":
		if jsonNumber.MatchString(n.Value) {
			return json.Number(n.Value), nil
		}
		// Another YAML spelling (0x1F, 1_000, +1): its value, if JSON has one.
		var f float64
		if err := n.Decode(&f); err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return nil, refuse("the number %q has no JSON form", n.Value)
		}
		if n.Tag == "!!int" {
			var i int64
			if err := n.Decode(&i); err != nil {
				return nil, refuse("the number %q has no JSON form", n.Value)
			}
			return json.Number(strconv.FormatInt(i, 10)), nil
		}
		return json.Number(strconv.FormatFloat(f, 'g', -1, 64)), nil
	}
	return nil, refuse("a value has the tag %s, which planque doesn't read", n.Tag)
}
