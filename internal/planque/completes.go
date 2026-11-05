package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kiff/planque/internal/profile"
	"github.com/kiff/planque/internal/store"
)

// A slow tool (kiff/planque#86, the slow-tools design): a write whose
// outcome comes later as one of the app's events. Its profile says so
// (Operation.Completes); the relay's answer then says "pending" with the
// handle the answer carries and the deadline, and the event that
// completes it is delivered with the same handle as its key, so the run
// that made the call is matched to it without reading the vendor's
// payload. A call whose completing event the taf doesn't listen for is
// refused before sending: a silent timeout would be worse.

// CodeNotListening: the taf has no subscription to the event that would
// complete the call.
const CodeNotListening = "not_listening"

// completesOf is the operation's completes, if its profile has one.
func completesOf(g *Grant, name string) *profile.Completes {
	op, ok := profileOperation(activeProfile(g), name)
	if !ok {
		return nil
	}
	return op.Completes
}

// listening says the taf has an active subscription to event on g.
func (s *Server) listening(ctx context.Context, g *Grant, event string) (bool, error) {
	triggers, err := s.accountTriggers(ctx, g.Tenant)
	if err != nil {
		return false, err
	}
	for _, t := range triggers {
		if t.Grant == g.ID && t.AppEvent && t.Active && !t.Removing && t.Definition.Slug == event {
			return true, nil
		}
	}
	return false, nil
}

// resultJSON is the JSON a tool result carries: its structured content,
// else its first text content parsed as JSON, else nothing.
func resultJSON(res any) map[string]any {
	raw, err := json.Marshal(res)
	if err != nil {
		return nil
	}
	var r struct {
		Structured map[string]any `json:"structuredContent"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if decodeNumbers(raw, &r) != nil {
		return nil
	}
	if r.Structured != nil {
		return r.Structured
	}
	for _, c := range r.Content {
		if c.Type != "text" {
			continue
		}
		var m map[string]any
		if decodeNumbers([]byte(c.Text), &m) == nil && m != nil {
			return m
		}
		return nil
	}
	return nil
}

// decodeNumbers decodes raw into v keeping every number's digits as the
// vendor wrote them (json.Number): a 64-bit numeric handle rounded to a
// float64 would never match the event's key, and the answer re-sent to
// the agent would carry changed numbers.
func decodeNumbers(raw []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(v)
}

// atPath is the value at a dotted path (with [n] for an index) in v, as
// a string; "" when there is none or it isn't a string or a number.
func atPath(v any, path string) string {
	for _, part := range strings.Split(strings.ReplaceAll(path, "[", ".["), ".") {
		if part == "" {
			continue
		}
		if strings.HasPrefix(part, "[") {
			i, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(part, "["), "]"))
			list, ok := v.([]any)
			if err != nil || !ok || i < 0 || i >= len(list) {
				return ""
			}
			v = list[i]
			continue
		}
		m, ok := v.(map[string]any)
		if !ok {
			return ""
		}
		v = m[part]
	}
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	return ""
}

// pending is the result of a call that completes later, as the agent
// reads it: the vendor's answer as it was, said pending, with the handle
// and the deadline in its structured content. When the answer carries no
// handle, it is left as it was: the run goes on, and the words say so.
func pending(res any, c *profile.Completes, grant string) any {
	// The answer as JSON, whatever Go shape it came in (the vendor's raw
	// MCP result, or a toolResult): its own copy, with its numbers kept.
	var m map[string]any
	raw, err := json.Marshal(res)
	if err != nil || decodeNumbers(raw, &m) != nil || m == nil {
		return res
	}
	if isErr, _ := m["isError"].(bool); isErr {
		return res
	}
	key := atPath(any(resultJSON(res)), c.Key)
	var text string
	if key == "" {
		text = "Sent. The answer carries no " + c.Key + ", so the outcome can't be waited for: " + c.Words
	} else {
		text = fmt.Sprintf("Sent. The outcome comes later as the event %s; the run waits up to %s. %s", c.Event, waitWords(c.TimeoutS), c.Words)
	}
	content, _ := m["content"].([]any)
	content = append(content, map[string]any{"type": "text", "text": text})
	m["content"] = content
	if key != "" {
		sc, _ := m["structuredContent"].(map[string]any)
		if sc == nil {
			sc = map[string]any{}
		}
		sc["pending"], sc["sent"], sc["timeout_s"] = true, true, c.TimeoutS
		sc["wait_for"] = map[string]any{"grant": grant, "event": c.Event, "key": key, "words": c.Words}
		m["structuredContent"] = sc
	}
	return m
}

func waitWords(seconds int) string {
	switch {
	case seconds < 3600:
		return fmt.Sprintf("%d min", seconds/60)
	case seconds < 86400:
		return fmt.Sprintf("%d h", seconds/3600)
	}
	return fmt.Sprintf("%d days", seconds/86400)
}

// eventKey is the handle an event carries for a call that completes on
// it, by the grant's profile: the value at completes.event_key in the
// payload, for the operation that completes on this event slug.
func eventKey(g *Grant, slug string, payload json.RawMessage) string {
	p, err := profile.Parse(activeProfile(g))
	if err != nil {
		return ""
	}
	for _, op := range p.Operations {
		if op.Completes == nil || op.Completes.Event != slug {
			continue
		}
		var data any
		if decodeNumbers(payload, &data) != nil {
			return ""
		}
		return atPath(data, op.Completes.EventKey)
	}
	return ""
}

// completed is finish for a call that ended OpDone: a call that completes
// later is answered as pending; any other as it was.
func (s *Server) completed(ctx context.Context, g *Grant, name string, claim store.Item, op operation, state string, res any) any {
	if state == OpDone {
		if c := completesOf(g, name); c != nil {
			res = pending(res, c, g.ID)
		}
	}
	return s.finish(ctx, claim, op, state, res)
}
