package planque

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"

	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/profile"
)

// A refund KIFF governs is followed after it's sent: its money can stay
// pending, then succeed, fail or be cancelled days later. The connected
// app's own event says so. When the grant's profile names that event in a
// refund operation's case block (case.refund_event), planque hands the event
// to KIFF as well as to the taf: KIFF reads the refund's state from it and
// never acts on anything else in it. Receiving one grants nothing.

// kindKIFFEvent is a delivery of an app event to KIFF.
const kindKIFFEvent = "kiff_event"

// refundEventOf reports whether slug is the refund event a case block of
// g's profile names.
func refundEventOf(g *Grant, slug string) bool {
	p, err := profile.Parse(activeProfile(g))
	if err != nil {
		return false
	}
	return slices.ContainsFunc(p.Operations, func(o profile.Operation) bool {
		return o.Case != nil && o.Case.RefundEvent != nil && o.Case.RefundEvent.Event == slug
	})
}

// kiffEvent is the app event e, as KIFF takes it: the same event, with an id
// of its own, the same for the provider's redelivery.
func kiffEvent(g *Grant, route triggerRecord, e connector.Event) (door.AppEvent, bool) {
	ev, ok := appEvent(g, route, e)
	if !ok || !refundEventOf(g, route.Definition.Slug) {
		return door.AppEvent{}, false
	}
	sum := sha256.Sum256([]byte("kiff_event\n" + route.Provider + "\n" + g.Tenant + "\n" + g.ID + "\n" + route.ID + "\n" + e.ID))
	ev.ID, ev.Key = "evk_"+hex.EncodeToString(sum[:]), ""
	return ev, true
}
