package planque

import (
	"regexp"
	"slices"
	"strings"

	"github.com/kiff/planque/internal/openapi"
)

// A key entry's address is one of three things (kiff/planque, the account
// template): fixed ("https://api.example.com/v1"); an account template on the
// vendor's own domain ("https://{account}.example.com/api/v2"), whose one
// variable the owner fills with one label; or a closed list of regions
// (Bases), the owner picking one. Whatever it is, the grant records the one
// address it was approved with, and a tool call can never change it.

// accountVar is the one variable an account template has.
const accountVar = "{account}"

// loginVar is the one variable a Basic username can have: the owner's own
// login (an email) at the vendor, as Zendesk's "{login}/token" or Gorgias's
// "{login}". The key stays the password; the login is no secret, and the
// approval shows the username it was approved with.
const loginVar = "{login}"

// loginTemplate reports whether an entry's Basic username names the login.
func loginTemplate(auth openapi.Auth) bool {
	return auth.In == "basic" && strings.Contains(auth.User, loginVar)
}

// loginUserValid reports whether a catalog entry's Basic username is one of
// the two forms the catalog takes: the login alone, or the login and "/token".
func loginUserValid(user string) bool { return user == loginVar || user == loginVar+"/token" }

// loginAddress is an email as one line of letters, digits and the few marks an
// address has: no space, colon, slash, quote or bracket, so it fills a Basic
// username as given and never ends it early.
var loginAddress = regexp.MustCompile(`^[A-Za-z0-9._%+-]{1,64}@[A-Za-z0-9-]{1,63}(\.[A-Za-z0-9-]{1,63})+$`)

// accountTemplate reports whether base names the account variable.
func accountTemplate(base string) bool { return strings.Contains(base, "{") }

// sampleBase is a template with a stand-in label, for the checks that need a
// real URL.
func sampleBase(base string) string { return strings.Replace(base, accountVar, "x", 1) }

// templateValid reports whether base is an account template the catalog
// takes: https, the variable once and alone as the host's first label, two
// fixed labels or more after it (the vendor's own domain), no port, no user
// info and no variable anywhere else.
func templateValid(base string) bool {
	if strings.Count(base, "{") != 1 || strings.Count(base, "}") != 1 || !strings.HasPrefix(base, "https://"+accountVar+".") {
		return false
	}
	return catalogHTTPS(sampleBase(base)) && openapi.OnAccountServer([]string{base}, sampleBase(base))
}

// resolveEntryAPI is the API a grant from a key entry is asked with: the
// entry's own, its address made one. account fills a template; base picks one
// of a list of regions. Anything else the ask names is refused, and why is
// the refusal in words ("" when it is fine).
func resolveEntryAPI(a *apiAsk, account, base, login string) (*apiAsk, string) {
	// Taken exactly as given: a stray space or newline is refused, never quietly trimmed.
	account = strings.ToLower(account)
	if a == nil {
		if account != "" || base != "" || login != "" {
			return nil, "this entry takes no account, region or login"
		}
		return nil, ""
	}
	out := *a
	out.Bases = nil
	switch {
	case len(a.Bases) > 0:
		if account != "" {
			return nil, "this entry is at one of its regions, not at an account's address: give base, not account"
		}
		if base == "" {
			return nil, "base is required: the region, one of " + strings.Join(a.Bases, ", ")
		}
		if !slices.Contains(a.Bases, base) {
			return nil, base + " isn't one of this entry's regions: " + strings.Join(a.Bases, ", ")
		}
		out.Base = base
	case accountTemplate(a.Base):
		if base != "" {
			return nil, "this entry is at an account's address: give account, not base"
		}
		if account == "" {
			return nil, "account is required: the owner's own name at the vendor, as " + a.Base
		}
		if !openapi.AccountLabel(account) {
			return nil, "account is one word of letters, digits and dashes, at most 63: no dot, @, colon or other character"
		}
		resolved := strings.Replace(a.Base, accountVar, account, 1)
		if !openapi.OnAccountServer([]string{a.Base}, resolved) {
			return nil, "account " + account + " doesn't give an address on the vendor's own domain"
		}
		out.Base = resolved
	default:
		if account != "" || base != "" {
			return nil, "this entry is at one fixed address: it takes no account or region"
		}
	}
	if loginTemplate(out.Auth) {
		if login == "" {
			return nil, "login is required: the owner's own login at the vendor, an email address"
		}
		if len(login) > 254 || !loginAddress.MatchString(login) {
			return nil, "login is one email address of letters, digits and . _ % + - : no space, colon, slash or other character"
		}
		out.Auth.User = strings.Replace(out.Auth.User, loginVar, login, 1)
	} else if login != "" {
		return nil, "this entry takes no login: its key is the whole credential"
	}
	return &out, ""
}
