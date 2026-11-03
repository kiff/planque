package planque

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
)

// maxEventConfig bounds the filters an owner approves for one app event.
const maxEventConfig = 8 << 10

// eventConfig checks the owner's filters for an app event against its
// trigger type's config schema, and answers them compacted, exactly as
// signed: an object whose every field the schema declares, with each
// required field present, each value of its declared type and, where the
// schema lists values, one of them. Anything the schema doesn't describe is
// refused, so no subscription listens wider than what the owner read.
func eventConfig(schema, config json.RawMessage) (json.RawMessage, error) {
	refuse := func(msg string) (json.RawMessage, error) {
		return nil, apiError{http.StatusBadRequest, msg}
	}
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	if len(config) > maxEventConfig {
		return refuse("the event's filters exceed 8 KB")
	}
	var s struct {
		Type       string `json:"type"`
		Properties map[string]struct {
			Type string            `json:"type"`
			Enum []json.RawMessage `json:"enum"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if json.Unmarshal(schema, &s) != nil || (s.Type != "" && s.Type != "object") {
		return refuse("this event's filters can't be checked")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(config, &fields); err != nil || fields == nil {
		return refuse("the event's filters are a JSON object")
	}
	for _, name := range s.Required {
		if v, ok := fields[name]; !ok || string(v) == "null" {
			return refuse("the event needs its filter " + name)
		}
	}
	for name, v := range fields {
		p, ok := s.Properties[name]
		if !ok {
			return refuse("the event has no filter " + name)
		}
		if p.Type != "" && jsonType(v) != p.Type && !(p.Type == "number" && jsonType(v) == "integer") {
			return refuse("the filter " + name + " is a " + p.Type)
		}
		if len(p.Enum) > 0 && !slices.ContainsFunc(p.Enum, func(e json.RawMessage) bool { return jsonEqual(e, v) }) {
			return refuse("the filter " + name + " isn't one of the values this event takes")
		}
	}
	var out bytes.Buffer
	if err := json.Compact(&out, config); err != nil {
		return refuse("the event's filters are a JSON object")
	}
	return out.Bytes(), nil
}

// jsonType is a JSON value's schema type; a whole number is an integer.
func jsonType(v json.RawMessage) string {
	var x any
	if json.Unmarshal(v, &x) != nil {
		return ""
	}
	switch n := x.(type) {
	case string:
		return "string"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case float64:
		if n == float64(int64(n)) {
			return "integer"
		}
		return "number"
	case nil:
		return "null"
	}
	return ""
}

func jsonEqual(a, b json.RawMessage) bool {
	var x, y bytes.Buffer
	return json.Compact(&x, a) == nil && json.Compact(&y, b) == nil && bytes.Equal(x.Bytes(), y.Bytes())
}
