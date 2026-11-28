package planque

import (
	"testing"

	"github.com/kiff/planque/internal/store"
)

// A sign-in planque holds that no taf's tools are in KIFF for (listed, waiting
// for its confirm, or reused while its tools wait) marks its app signed_in, so
// the owner's list puts it under their own apps with its Disconnect. A
// connection that's connected or must sign in again says that instead, and a
// request with nothing held, or a finished one, marks nothing.
func TestAHeldSignInNothingUsesMarksItsAppSignedIn(t *testing.T) {
	s := &Server{KV: &store.Memory{}}
	grant := func(id, entry, status string, signedIn bool) *Grant {
		return &Grant{ID: id, Tenant: tenant, Status: status, SignedIn: signedIn, Request: Request{CatalogEntry: entry, Tenant: tenant}}
	}
	for _, g := range []*Grant{
		grant("grt_a_listed", "qonto@2", StatusListed, true),
		grant("grt_b_confirm", "stripe@3", StatusSignedIn, true),
		grant("grt_c_reused", "typeform-key@1", StatusPending, true),
		grant("grt_d_asked", "fifth-key@1", StatusPending, false),
		grant("grt_e_approved", "sixth-key@1", StatusApproved, false),
		grant("grt_f_ended", "seventh-key@1", StatusDisconnected, false),
		// Two connections to one app: the stronger state stands, whichever came first.
		grant("grt_g_listed", "gmail-key@1", StatusListed, true),
		grant("grt_h_connected", "gmail-key@1", StatusConnected, true),
		grant("grt_i_reconnect", "pipedrive-key@1", StatusNeedsReconnect, true),
		grant("grt_j_listed", "pipedrive-key@1", StatusListed, true),
		// Two held sign-ins to one app: the older grant.
		grant("grt_l_listed", "notion-key@1", StatusListed, true),
		grant("grt_k_listed", "notion-key@1", StatusListed, true),
	} {
		if err := s.saveGrant(t.Context(), g); err != nil {
			t.Fatal(err)
		}
	}
	states, err := s.entryStates(t.Context(), tenant)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]catalogEntryView{
		"qonto@2":         {State: "signed_in", Grant: "grt_a_listed"},
		"stripe@3":        {State: "signed_in", Grant: "grt_b_confirm"},
		"typeform-key@1":  {State: "signed_in", Grant: "grt_c_reused"},
		"gmail-key@1":     {State: "connected", Grant: "grt_h_connected"},
		"pipedrive-key@1": {State: "needs_reconnect", Grant: "grt_i_reconnect"},
		"notion-key@1":    {State: "signed_in", Grant: "grt_k_listed"},
	}
	if len(states) != len(want) {
		t.Errorf("states: got %d entries, want %d: %+v", len(states), len(want), states)
	}
	for entry, w := range want {
		if got := states[entry]; got.State != w.State || got.Grant != w.Grant {
			t.Errorf("%s: got %s by %s, want %s by %s", entry, got.State, got.Grant, w.State, w.Grant)
		}
	}
}
