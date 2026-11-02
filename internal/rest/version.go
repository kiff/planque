package rest

import (
	"errors"
	"regexp"
	"strings"
)

// VersionPin is a reviewed vendor version header. It is public approval
// data, held separately from authentication and transport headers.
type VersionPin struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var versionHeaderName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]{0,63}$`)

func (p VersionPin) Validate() error {
	name := strings.ToLower(p.Name)
	// A pin names a version; it cannot replace authentication, HTTP framing,
	// content negotiation or a routing header.
	if !versionHeaderName.MatchString(p.Name) || !strings.Contains(name, "version") || len(p.Value) > 1024 || strings.TrimSpace(p.Value) == "" || strings.ContainsFunc(p.Value, func(r rune) bool { return r < 0x20 || r > 0x7e }) {
		return errors.New("a version pin needs a version header name and a single-line value")
	}
	return nil
}
