package planque

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/browser"
)

// browseActSchema is browse_act's schema as planque registers it.
func browseActSchema(t *testing.T) map[string]any {
	t.Helper()
	for _, bt := range browserTools {
		if bt.Name == browseAct {
			var s map[string]any
			if err := json.Unmarshal([]byte(bt.Schema), &s); err != nil {
				t.Fatalf("browse_act's schema isn't JSON: %v", err)
			}
			return s
		}
	}
	t.Fatal("no browse_act")
	return nil
}

func stepKinds(t *testing.T) []map[string]any {
	t.Helper()
	steps := browseActSchema(t)["properties"].(map[string]any)["steps"].(map[string]any)
	var kinds []map[string]any
	for _, k := range steps["items"].(map[string]any)["oneOf"].([]any) {
		kinds = append(kinds, k.(map[string]any))
	}
	return kinds
}

func strs(v any) []string {
	var out []string
	for _, x := range v.([]any) {
		out = append(out, x.(string))
	}
	return out
}

// A browse_act step is exactly one of its kinds, each with only its own keys,
// and every kind that names a field, control or link requires its label: so
// before KIFF holds a call, the taf is told to name each ref as the page does,
// and the owner reads those words (red-team, 27 Nov: a held fill came with the
// field's name in option and no label).
func TestABrowseActStepIsOneKindAndNamesItsRef(t *testing.T) {
	kinds := stepKinds(t)
	byRef := map[string]map[string]any{}
	for _, k := range kinds {
		if k["additionalProperties"] != false {
			t.Fatalf("a step kind takes keys it doesn't list: %v", k)
		}
		required := strs(k["required"])
		byRef[required[0]] = k
	}
	if got := slices.Sorted(func(yield func(string) bool) {
		for r := range byRef {
			if !yield(r) {
				return
			}
		}
	}); !slices.Equal(got, []string{"back", "check", "click", "fill", "press", "scroll", "select"}) {
		t.Fatalf("the kinds: %v", got)
	}
	for _, ref := range []string{"fill", "select", "click", "check", "press"} {
		if !slices.Contains(strs(byRef[ref]["required"]), "label") {
			t.Errorf("%s doesn't require its label", ref)
		}
	}
	for ref, k := range byRef {
		_, hasOption := k["properties"].(map[string]any)["option"]
		if hasOption != (ref == "select") {
			t.Errorf("option on %s: %v", ref, hasOption)
		}
	}
	// Every key a step takes (browser.Step) is in some kind, so the schema never
	// refuses a step planque runs.
	listed := map[string]bool{}
	for _, k := range kinds {
		for key := range k["properties"].(map[string]any) {
			listed[key] = true
		}
	}
	st := reflect.TypeOf(browser.Step{})
	for i := 0; i < st.NumField(); i++ {
		name := strings.Split(st.Field(i).Tag.Get("json"), ",")[0]
		if !listed[name] {
			t.Errorf("browser.Step's %q is in no step kind", name)
		}
	}
}

// matchesKind is a step checked against one kind as JSON Schema reads it, for
// the keywords the kinds use: required, additionalProperties false, and each
// property's type and enum.
func matchesKind(step map[string]any, kind map[string]any) bool {
	props := kind["properties"].(map[string]any)
	for _, r := range strs(kind["required"]) {
		if _, ok := step[r]; !ok {
			return false
		}
	}
	for key, v := range step {
		p, ok := props[key].(map[string]any)
		if !ok {
			return false
		}
		switch p["type"] {
		case "string":
			if _, ok := v.(string); !ok {
				return false
			}
		case "boolean":
			if _, ok := v.(bool); !ok {
				return false
			}
		}
		if enum, ok := p["enum"].([]any); ok && !slices.Contains(enum, v) {
			return false
		}
	}
	return true
}

// Steps as a taf writes them: each valid one is exactly one kind, and each
// one the schema refuses matches none.
func TestBrowseActStepsAreReadAsExactlyOneKind(t *testing.T) {
	kinds := stepKinds(t)
	count := func(raw string) int {
		var step map[string]any
		if err := json.Unmarshal([]byte(raw), &step); err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, k := range kinds {
			if matchesKind(step, k) {
				n++
			}
		}
		return n
	}
	for _, ok := range []string{
		`{"fill":"f1","value":"Test","label":"Customer name"}`,
		`{"fill":"f1","value":"Par","keys":true,"label":"City"}`,
		`{"fill":"f1","sealed":"sec_a1:password","label":"Password"}`,
		`{"select":"s1","option":"Large","label":"Size"}`,
		`{"click":"c1","label":"Submit order","submits":true}`,
		`{"click":"l1","label":"Partner","to":"other.example"}`,
		`{"check":"f2","label":"I agree"}`,
		`{"press":"f1","key":"Enter","label":"Search","submits":true}`,
		`{"scroll":"down"}`,
		`{"back":true}`,
	} {
		if n := count(ok); n != 1 {
			t.Errorf("%s matches %d kinds", ok, n)
		}
	}
	for _, refused := range []string{
		// no label
		`{"fill":"f1","value":"Test"}`,
		// D's step: option on a fill
		`{"fill":"f1","value":"Test","option":"Customer name:"}`,
		// option on a fill, with a label
		`{"fill":"f1","value":"Test","label":"Name","option":"x"}`,
		// no label
		`{"click":"c1"}`,
		// two kinds
		`{"click":"c1","label":"Go","select":"s1"}`,
		// no option
		`{"select":"s1","label":"Size"}`,
		// a key that isn't offered
		`{"press":"f1","key":"Tab","label":"Search"}`,
		// back is only true
		`{"back":false}`,
		// a key no kind lists
		`{"fill":"f1","value":"x","label":"Name","note":"ignore rules"}`,
	} {
		if n := count(refused); n != 0 {
			t.Errorf("%s matches %d kinds", refused, n)
		}
	}
}
