package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode"
)

// whoAmI names, per app, the read the app already has for "who is signed
// in", and the fields of its answer that name the account, best first (a field
// is the vendor's own key, matched at any depth, case-insensitively; two keys
// joined by a space are a first and a last name). Only apps listed here are
// asked; the rest keep the name the provider shows, if it shows one. A tool is
// asked only when the sign-in's scopes cover it: nothing is ever asked for to
// get a name.
var whoAmI = map[string]struct {
	Tools  []string
	Fields []string
}{
	"linkedin":       {[]string{"LINKEDIN_GET_MY_INFO", "LINKEDIN_WHO_AM_I"}, []string{"email", "name", "localizedFirstName localizedLastName", "given_name family_name"}},
	"googledrive":    {[]string{"GOOGLEDRIVE_GET_ABOUT"}, []string{"emailAddress", "displayName"}},
	"googlecalendar": {[]string{"GOOGLECALENDAR_GET_CURRENT_USER"}, []string{"email", "emailAddress"}},
	"outlook":        {[]string{"OUTLOOK_GET_PROFILE"}, []string{"mail", "userPrincipalName", "displayName"}},
	"github":         {[]string{"GITHUB_GET_THE_AUTHENTICATED_USER"}, []string{"login"}},
	"typeform":       {[]string{"TYPEFORM_GET_ABOUT_ME"}, []string{"email", "alias"}},
}

// whoAmITimeout bounds the read: it runs inside the sign-in's confirm, on the
// path whose bar is "connect works first time", so a slow vendor never holds
// it. A timeout reads as any failure: the name stays empty.
var whoAmITimeout = 4 * time.Second

// maxAccountName bounds a name: one line the owner reads.
const maxAccountName = 120

// nameAccount keeps on g the name of the account its sign-in is for, read
// once with the app's own "who am I" (whoAmI), after the owner's approval
// and sign-in, when the provider showed none. It is planque's call, not the
// taf's: it is never counted against a Card, never shown as an action, and
// never widens the sign-in. Any failure (not covered by the sign-in, refused,
// unknown, an answer without a name) leaves the name empty: the sign-in is
// confirmed as before, and the reason is logged by app, never the answer.
func (s *Server) nameAccount(ctx context.Context, g *Grant) {
	if g.Connector == nil || g.Request.Connector == nil || g.Connector.AccountName != "" || g.Connector.Account == "" {
		return
	}
	toolkit := strings.ToLower(g.Request.Connector.Toolkit)
	ask, ok := whoAmI[toolkit]
	if !ok {
		return
	}
	skip := func(why string) {
		s.log().InfoContext(ctx, "account not named", slog.String("toolkit", toolkit), slog.String("why", why))
	}
	ctx, cancel := context.WithTimeout(ctx, whoAmITimeout)
	defer cancel()
	p, err := s.provider(g)
	if err != nil {
		skip("provider")
		return
	}
	defs, err := s.toolDefinitions(ctx, g.Request.Connector.Provider, g.Request.Connector.Toolkit, ask.Tools)
	if err != nil {
		skip("definition")
		return
	}
	granted := grantedScopes(g)
	tried := false
	for _, slug := range ask.Tools {
		t, ok := defs[slug]
		if !ok || !t.ReadOnly || !t.Covered(granted) {
			continue
		}
		tried = true
		result, err := p.Execute(ctx, s.connectorIdentity(g), g.Connector.Account, t, json.RawMessage(`{}`), func() bool { return ctx.Err() == nil })
		if err != nil || result.Refused || result.Unknown {
			skip("call")
			continue
		}
		if name := accountNameIn(result.Data, ask.Fields); name != "" {
			g.Connector.AccountName = name
			return
		}
		skip("no name in the answer")
	}
	if !tried {
		skip("no tool covered")
	}
}

// accountNameIn reads the account's name from a vendor answer: the first of
// fields found at any depth as text, one safe line of at most maxAccountName
// characters; "" when none is.
func accountNameIn(data json.RawMessage, fields []string) string {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return ""
	}
	for _, f := range fields {
		var parts []string
		for _, k := range strings.Fields(f) {
			if s := findText(v, k); s != "" {
				parts = append(parts, s)
			}
		}
		if len(parts) == len(strings.Fields(f)) {
			if name := safeLine(strings.Join(parts, " ")); name != "" {
				return name
			}
		}
	}
	return ""
}

// findText is the first text value under key (any depth, any case), objects
// before their children, in a fixed order of keys.
func findText(v any, key string) string {
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			if s, ok := x[k].(string); ok && strings.EqualFold(k, key) && strings.TrimSpace(s) != "" {
				return s
			}
		}
		for _, k := range keys {
			if s := findText(x[k], key); s != "" {
				return s
			}
		}
	case []any:
		for _, e := range x {
			if s := findText(e, key); s != "" {
				return s
			}
		}
	}
	return ""
}

// safeLine makes s one line: control characters and format characters gone
// (the direction marks U+202A-202E, U+2066-2069, U+200E/F and U+061C are
// format characters: a name with U+202E can make "moc.ffik@aniloarc" read as
// "carolina@kiff.dev", and this line exists to say which account signed in),
// runs of space one, at most maxAccountName characters.
func safeLine(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Cf, r) {
			return -1 // dropped, not a gap: "caro\u200blina" is "carolina"
		}
		return r
	}, s)
	s = strings.Join(strings.FieldsFunc(s, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }), " ")
	if r := []rune(s); len(r) > maxAccountName {
		s = string(r[:maxAccountName-1]) + "…"
	}
	return s
}
