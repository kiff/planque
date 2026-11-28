package planque

import (
	"context"
	"encoding/json"
)

// sharedSignIn is another live grant on a connection's sign-in, as the
// connection's row shows it: its id and its taf, never a token reference.
type sharedSignIn struct {
	Grant string `json:"grant"`
	// Agent and Taf are the taf the other grant is for: its id and its name
	// for people, as the grant's request names them.
	Agent string `json:"agent,omitempty"`
	Taf   string `json:"taf,omitempty"`
}

// sharesSignIn: o is another live grant on g's sign-in, as removal reads
// it (signInUsers): one that reuses the same token, or for a provider's app
// one that may be on the same vendor account. A disconnected or failed grant
// uses nothing; a pending one counts.
func sharesSignIn(g, o *Grant) bool {
	if o.ID == g.ID || o.Status == StatusDisconnected || o.Status == StatusFailed {
		return false
	}
	return sameSignIn(g, o) || mayShareVendorAccount(g, o)
}

// sharedSignInsOf are the grants among all that share g's sign-in, as its
// row shows them: an empty list, never nil, when none does.
func sharedSignInsOf(g *Grant, all []*Grant) []sharedSignIn {
	out := []sharedSignIn{}
	for _, o := range all {
		if sharesSignIn(g, o) {
			out = append(out, sharedSignIn{Grant: o.ID, Agent: o.Request.Agent, Taf: label(o.Request.AgentName)})
		}
	}
	return out
}

// accountGrants reads every grant of the account, as signInUsers does.
func (s *Server) accountGrants(ctx context.Context, tenant string) ([]*Grant, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "GRANT#")
	if err != nil {
		return nil, err
	}
	out := make([]*Grant, 0, len(items))
	for _, it := range items {
		var o Grant
		if json.Unmarshal(it.Doc, &o) != nil {
			continue
		}
		out = append(out, &o)
	}
	return out, nil
}
