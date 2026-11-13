package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/oauth"
)

// kiffConnection is t as KIFF's gateway registers it: relayed to g, with
// transport as the gateway's credential.
func kiffConnection(g *Grant, t GrantedTool, transport string) kiff.Connection {
	return kiff.Connection{Name: t.Name, Title: t.Title, URL: g.Audience, Tool: t.Name,
		Description: t.Description, AmountArg: t.AmountArg, AmountUnit: t.AmountUnit, AmountScale: t.AmountScale,
		AmountPath: t.AmountPath, AmountFormat: t.AmountFormat, ResourcePath: t.ResourcePath,
		IdempotencyArg: IdempotencyArg, ReadOnly: t.ReadOnly, Mode: "relay", Audience: g.Audience, Credential: transport,
		ResourceArg: t.ResourceArg, UnitLookup: t.UnitLookup, AmountLookup: t.AmountLookup, VendorTool: t.Tool}
}

// declares says g's tools are declared to KIFF when it's asked for: a key
// grant asked with its tools, waiting for the owner's key. Its tools come
// from the API description pinned at the request, so they are known before
// the key; a sign-in's come from the vendor only once signed in.
func declares(g *Grant) bool {
	return g.Request.Kind == KindAPIKey && g.Request.API != nil && g.Status == StatusPending && !g.SignedIn &&
		g.TokenOf == "" && g.Connected == nil && len(g.Request.Tools) > 0 && !g.Declared
}

// declare registers a key grant's tools in KIFF as soon as it's asked for,
// so Coach can draft the taf's Card over them and the owner signs the key
// and the Card together. Declared isn't usable: the relay refuses every call
// until the key is in and the grant connects, and the same registration
// then serves it, its credential unchanged. A declaration that fails leaves
// the request as it was, with nothing in KIFF: it connects as before, once
// the key is given.
//
// refused is set when KIFF turned a tool down for what it is (a 4xx): the
// key would meet the same refusal after the owner gave it, so the request
// ends here instead, failed with KIFF's words, and nothing is left in KIFF.
func (s *Server) declare(ctx context.Context, g *Grant) (refused string) {
	if !declares(g) {
		return ""
	}
	note := func(why string) {
		g.event(s.now(), "its tools weren't declared to KIFF before the key: "+why+"; they're registered once the key is given")
		if err := s.saveGrant(ctx, g); err != nil {
			s.log().WarnContext(ctx, "declare: the grant's record wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
		}
	}
	listed, _, err := s.vendorTools(ctx, g, "")
	if err != nil {
		note("its API description couldn't be read: " + err.Error())
		return ""
	}
	var tools []GrantedTool
	for _, t := range g.Request.Tools {
		gt, why := granted(g, t, listed)
		if why != "" {
			note(why)
			return ""
		}
		tools = append(tools, gt)
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		note("the KIFF connector key could not be opened")
		return ""
	}
	transport := oauth.Random(32)
	sum := sha256.Sum256([]byte(transport))
	// Saved as declared first, so the relay answers KIFF's probe while the
	// tools are registered.
	g.Tools, g.Transport, g.Declared = tools, hex.EncodeToString(sum[:]), true
	if err := s.saveGrant(ctx, g); err != nil {
		g.Tools, g.Transport, g.Declared = nil, "", false
		s.log().WarnContext(ctx, "declare: the grant wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return ""
	}
	var done, titles []string
	for _, t := range tools {
		if err := s.KIFF.Connect(ctx, key, kiffConnection(g, t, transport)); err != nil {
			for _, name := range done {
				_ = s.KIFF.Disconnect(ctx, key, name)
			}
			g.Tools, g.Transport, g.Declared = nil, "", false
			if a := (*kiff.Answered)(nil); errors.As(err, &a) && a.Refused() {
				why := "KIFF won't take " + t.Name + ": " + a.Msg
				if err := s.failGrant(ctx, g, why); err != nil && err.Error() != why {
					s.log().WarnContext(ctx, "declare: the refused grant wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
				}
				return why
			}
			note("KIFF did not register " + t.Name + ": " + err.Error())
			return ""
		}
		done = append(done, t.Name)
		titles = append(titles, toolTitle(t.ToolRequest, CatalogTool{Name: t.Tool, Title: t.Title}))
	}
	g.event(s.now(), "Declared in KIFF, waiting for the key: "+strings.Join(titles, ", ")+"; nothing is usable until the key is given")
	if err := s.saveGrant(ctx, g); err != nil {
		s.log().WarnContext(ctx, "declare: the grant's record wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
	return ""
}

// undeclare takes a request's declared tools out of KIFF, when it ends
// before its key: nothing can be drafted over a request that's gone. An
// error leaves g as it was, for the next sweep or read to try again.
func (s *Server) undeclare(ctx context.Context, g *Grant) error {
	if !g.Declared || len(g.Tools) == 0 {
		g.Declared = false
		return nil
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		return fmt.Errorf("open the KIFF connector key to take the declared tools out: %w", err)
	}
	for _, t := range g.Tools {
		if err := s.KIFF.Disconnect(ctx, key, t.Name); err != nil {
			return fmt.Errorf("KIFF kept the declared %s: %w", t.Name, err)
		}
	}
	g.Tools, g.Transport, g.Declared = nil, "", false
	return nil
}

// connectDeclared connects a grant whose tools were declared at its
// request, once the owner gave the key: the tools are in KIFF already, with
// the credential the relay checks, so they only become usable. The
// approved profile and Setup's reads are applied as for any connection; if
// the profile fails, the declared tools go too, and the key is deleted.
func (s *Server) connectDeclared(ctx context.Context, g *Grant) error {
	g.Status, g.HeldUntil, g.Declared = StatusConnected, time.Time{}, false
	if err := s.saveGrant(ctx, g); err != nil {
		g.Status, g.Declared = StatusPending, true
		return err
	}
	var names, titles []string
	for _, t := range g.Tools {
		names = append(names, t.Name)
		titles = append(titles, toolTitle(t.ToolRequest, CatalogTool{Name: t.Tool, Title: t.Title}))
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		return s.undo(ctx, g, names, "the KIFF connector key could not be opened")
	}
	if err := s.applyProfile(ctx, key, g); err != nil {
		return s.undo(ctx, g, names, "KIFF did not apply the approved profile: "+err.Error())
	}
	g.event(s.now(), "Connected in KIFF: "+strings.Join(titles, ", ")+" (declared at the request)")
	s.setupReads(ctx, key, g)
	return s.saveGrant(ctx, g)
}

// declareAddition registers the tools added to a connected grant in KIFF
// while the addition waits for the owner's approval, so the taf's Card can be
// drafted over them and signed with that approval (tools_and_cards). They get
// a transport credential of their own, as any addition does, and stay
// unusable: the relay lists them and refuses every call until approved. A
// declaration that fails leaves the addition as it was, with nothing in KIFF:
// it's registered once approved, as before.
func (s *Server) declareAddition(ctx context.Context, g *Grant) {
	if g.Status != StatusPending || g.Connected == nil || len(g.Pending) > 0 || len(g.Request.Tools) <= len(g.Connected) {
		return
	}
	note := func(why string) {
		g.event(s.now(), "the added tools weren't declared to KIFF before the approval: "+why+"; they're registered once approved")
		if err := s.saveGrant(ctx, g); err != nil {
			s.log().WarnContext(ctx, "declare addition: the grant's record wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
		}
	}
	access := ""
	if g.Request.Connector == nil && g.Request.API == nil && g.Request.Kind != KindBrowser {
		tok, _, err := s.openToken(ctx, g)
		if err != nil {
			note("the sign-in couldn't be opened to list the tools")
			return
		}
		access = tok.Access
	}
	listed, _, err := s.vendorTools(ctx, g, access)
	if err != nil {
		note("the tools couldn't be listed: " + err.Error())
		return
	}
	var added []GrantedTool
	for _, t := range g.Request.Tools[len(g.Connected):] {
		gt, why := granted(g, t, listed)
		if why != "" {
			note(why)
			return
		}
		added = append(added, gt)
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		note("the KIFF connector key could not be opened")
		return
	}
	transport := oauth.Random(32)
	sum := sha256.Sum256([]byte(transport))
	// Saved first, so the relay answers KIFF's probe while they're registered.
	g.Pending, g.PendingTransport = added, hex.EncodeToString(sum[:])
	if err := s.saveGrant(ctx, g); err != nil {
		g.Pending, g.PendingTransport = nil, ""
		s.log().WarnContext(ctx, "declare addition: the grant wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
		return
	}
	var done, titles []string
	for _, t := range added {
		if err := s.KIFF.Connect(ctx, key, kiffConnection(g, t, transport)); err != nil {
			for _, name := range done {
				_ = s.KIFF.Disconnect(ctx, key, name)
			}
			g.Pending, g.PendingTransport = nil, ""
			note("KIFF did not register " + t.Name + ": " + err.Error())
			return
		}
		done = append(done, t.Name)
		titles = append(titles, toolTitle(t.ToolRequest, CatalogTool{Name: t.Tool, Title: t.Title}))
	}
	g.event(s.now(), "Declared in KIFF, waiting for the owner's approval: "+strings.Join(titles, ", ")+"; nothing is usable until approved")
	if err := s.saveGrant(ctx, g); err != nil {
		s.log().WarnContext(ctx, "declare addition: the grant's record wasn't saved", slog.String("grant", g.ID), slog.String("err", err.Error()))
	}
}

// undeclareAddition takes an addition's declared tools out of KIFF when it
// ends unapproved: past its day, or replaced by another. No Card can then be
// signed over tools nobody approved. A tool the grant has connected stays. An
// error leaves g as it was, for the next sweep or read to try again.
func (s *Server) undeclareAddition(ctx context.Context, g *Grant) error {
	if len(g.Pending) == 0 {
		g.PendingTransport = ""
		return nil
	}
	key, err := s.connectorKey(ctx, g.Tenant)
	if err != nil {
		return fmt.Errorf("open the KIFF connector key to take the declared tools out: %w", err)
	}
	for _, t := range g.Pending {
		if slices.ContainsFunc(g.Tools, func(h GrantedTool) bool { return h.Name == t.Name }) {
			continue
		}
		if err := s.KIFF.Disconnect(ctx, key, t.Name); err != nil {
			return fmt.Errorf("KIFF kept the declared %s: %w", t.Name, err)
		}
	}
	g.Pending, g.PendingTransport = nil, ""
	return nil
}

// sameNames says the declared tools are, by name and in order, the ones the
// addition asks for: what was declared is what the owner approves.
func sameNames(declared []GrantedTool, asked []ToolRequest) bool {
	if len(declared) != len(asked) {
		return false
	}
	for i := range declared {
		if declared[i].Name != asked[i].Name {
			return false
		}
	}
	return true
}

// toolNamesOf are the tools' names in KIFF.
func toolNamesOf(tools []GrantedTool) []string {
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		out = append(out, t.Name)
	}
	return out
}
