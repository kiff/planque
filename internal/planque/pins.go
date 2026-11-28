package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/permit"
)

// toolPins are the values planque sets in a tool's vendor body after the
// permit, as KIFF's gateway registers them (vendor_pins) to sign into each
// permit: Qonto's draft pins (qontoDraftOnly), the account a Stripe tool acts
// in (callAccount), and a profile's marker (profileMarker), whose value is
// the permit's op. They're read from the same steps the relay runs, on empty
// arguments, so what's registered is exactly what the relay sets. None for
// any other tool.
func toolPins(g *Grant, t GrantedTool) []kiff.Pin {
	var pins []kiff.Pin
	if _, set, why := qontoDraftOnly(g.Request.Server, t, json.RawMessage(`{}`), g.Country); why == "" {
		paths := mapsKeys(set)
		slices.Sort(paths)
		for _, path := range paths {
			if pinValue(set[path]) {
				pins = append(pins, kiff.Pin{Path: path, Value: set[path]})
			}
		}
	}
	if held, why := callAccount(g, t.Tool, json.RawMessage(`{}`)); why == "" && !bytes.Equal(bytes.TrimSpace(held), []byte(`{}`)) {
		var set map[string]any
		if json.Unmarshal(held, &set) == nil {
			if acct, ok := set[stripeContext].(string); ok {
				pins = append(pins, kiff.Pin{Path: stripeContext, Value: acct})
			}
		}
	}
	if op, ok := profileOperation(activeProfile(g), t.Name); ok && op.Marker != nil {
		pins = append(pins, kiff.Pin{Path: op.Marker.Field, Marker: true})
	}
	return pins
}

// pinValue reports whether v can be pinned: a permit signs only strings and
// booleans.
func pinValue(v any) bool {
	switch v.(type) {
	case string, bool:
		return true
	}
	return false
}

func mapsKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// pinDiffs compares what planque pins for a tool with what KIFF signed in
// the permit and with the body about to be sent, both ways. Each difference
// is "<path>: <what differs>", never a value. Empty when they agree, or when
// neither side has a pin.
func pinDiffs(want []kiff.Pin, p permit.Permit, body json.RawMessage) []string {
	var sent map[string]any
	_ = json.Unmarshal(body, &sent)
	var out []string
	pinned := map[string]bool{}
	for _, w := range want {
		pinned[w.Path] = true
		expect := w.Value
		if w.Marker {
			expect = p.Op
		}
		signed, ok := p.Pins[w.Path]
		switch {
		case !ok:
			out = append(out, w.Path+": pinned here, not signed")
		case signed != expect:
			out = append(out, w.Path+": signed with another value")
		}
		if got, ok := valueAt(sent, w.Path); !ok || got != expect {
			out = append(out, w.Path+": the body sent doesn't carry the pinned value")
		}
	}
	for path := range p.Pins {
		if !pinned[path] {
			out = append(out, path+": signed, not pinned here")
		}
	}
	slices.Sort(out)
	return out
}

// valueAt reads a dotted path of object keys in a JSON object.
func valueAt(m map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var cur any = m
	for _, part := range parts {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = obj[part]; !ok {
			return nil, false
		}
	}
	return cur, true
}

// pinRefusal is why a call whose pins don't hold is refused, for an account
// that enforces pins: a tool with pins whose permit carries none signed, or
// any difference pinDiffs finds (pinned here and not signed, signed with
// another value, the body not carrying the pinned value, or signed and not
// pinned here). Paths only, never a value. "" when nothing is pinned on
// either side, or every pin holds.
func pinRefusal(want []kiff.Pin, p permit.Permit, body json.RawMessage) string {
	if len(want) == 0 && len(p.Pins) == 0 {
		return ""
	}
	if !p.PinsSigned {
		return "the tool's pinned values weren't signed by KIFF"
	}
	if d := pinDiffs(want, p, body); len(d) > 0 {
		return "the pinned values don't hold (" + strings.Join(d, "; ") + ")"
	}
	return ""
}

// enforcesPins reports whether tenant's calls are refused on any pin
// difference (EnforcePins; "*" is every account).
func (s *Server) enforcesPins(tenant string) bool {
	return s.EnforcePins[tenant] || s.EnforcePins["*"]
}

// logPinDiffs logs, for one call, each difference between what planque pins
// and what KIFF signed (pinDiffs), or that KIFF signed no pins for a tool
// that has some. Paths only, never a value. It decides nothing yet.
func (s *Server) logPinDiffs(ctx context.Context, g *Grant, tool GrantedTool, p permit.Permit, body json.RawMessage) {
	want := toolPins(g, tool)
	if len(want) == 0 && len(p.Pins) == 0 {
		return
	}
	if !p.PinsSigned {
		s.log().InfoContext(ctx, "pins unsigned", slog.String("grant", g.ID), slog.String("tool", tool.Name))
		return
	}
	for _, d := range pinDiffs(want, p, body) {
		s.log().WarnContext(ctx, "pin differs", slog.String("grant", g.ID), slog.String("tool", tool.Name), slog.String("diff", d))
	}
}
