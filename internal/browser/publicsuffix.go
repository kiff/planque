package browser

import (
	_ "embed"
	"strings"
)

// The Public Suffix List (public_suffix_list.dat, MPL 2.0, refreshed with
// tools/publicsuffix.sh): the names under which anyone may register their
// own. That's ICANN's ("com", "co.uk") and its private section's shared
// hosts ("github.io", "vercel.app", "myshopify.com", "blogspot.com"). A
// browser grant's site covers its subdomains, so a site that is itself a
// public suffix would reach every other owner's site under it.
//
//go:embed public_suffix_list.dat
var publicSuffixRaw string

// publicSuffixes are the list's rules: "x" a suffix, "*.x" every name one
// label under x, "!x" an exception to a wildcard.
var publicSuffixes = func() map[string]bool {
	m := map[string]bool{}
	for _, line := range strings.Split(publicSuffixRaw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		m[strings.ToLower(line)] = true
	}
	return m
}()

// PublicSuffix reports that name (a host name, lower case) is itself a
// public suffix: a rule of the list, or one label under a wildcard rule that
// no exception names. A name under a suffix ("shop.github.io") isn't one.
func PublicSuffix(name string) bool {
	name = strings.TrimSuffix(strings.ToLower(name), ".")
	if name == "" {
		return false
	}
	if publicSuffixes[name] {
		return true
	}
	if publicSuffixes["!"+name] {
		return false
	}
	if _, parent, ok := strings.Cut(name, "."); ok && publicSuffixes["*."+parent] {
		return true
	}
	return false
}
