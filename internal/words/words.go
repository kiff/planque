// Package words holds the lines planque's own browser pages show a person:
// the result of a vendor's sign-in and the sign-in view, in the eight
// languages tafs ships. A page speaks the language of the app that started
// the connection, as the app said it, else the browser's. The terms follow
// tafs' glossary (kiff/tafs docs/languages.md): connection, Approve, Done.
package words

import (
	"fmt"
	"strconv"
	"strings"
)

// Lines is one language's lines. A field ending in a format verb takes
// its values through the methods below, never through the caller's
// English.
type Lines struct {
	// Titles.
	NotConnected     string
	NotSignedInAgain string
	SignInView       string

	// The result of a vendor's sign-in.
	MissingState          string
	LinkUnknown           string
	LinkUsed              string
	NotCompletedF         string // the vendor's host
	NotCompletedRenewF    string // the vendor's host
	OtherServer           string
	NotFinishedF          string // the vendor's host
	NotFinishedRenewF     string // the vendor's host
	CantConfirm           string
	OtherAccountUnknownF  string // the vendor's name
	OtherAccountF         string // the vendor's name
	SignInTafAccount      string
	SignInTafAccountRenew string
	NotKept               string
	NothingChanged        string

	// The sign-in view.
	ViewTitleF         string // the site
	Go                 string
	TypeHere           string
	Send               string
	Enter              string
	Tab                string
	Done               string
	FrameAlt           string
	ViewNoteF          string // the site, then its minutes
	ViewTitleAny       string // the title when no site is named
	ViewNoteAnyF       string // its minutes, when no site is named
	DidntWork          string
	DoneWords          string
	NotHandedBack      string
	ViewUnknown        string
	ViewGone           string
	ViewLapsedF        string // its minutes
	ViewConnectionGone string
	ViewFinished       string
	ViewNoBrowser      string
	ViewBrowserEnded   string
}

func (l Lines) NotCompleted(host string) string      { return fmt.Sprintf(l.NotCompletedF, host) }
func (l Lines) NotCompletedRenew(host string) string { return fmt.Sprintf(l.NotCompletedRenewF, host) }
func (l Lines) NotFinished(host string) string       { return fmt.Sprintf(l.NotFinishedF, host) }
func (l Lines) NotFinishedRenew(host string) string  { return fmt.Sprintf(l.NotFinishedRenewF, host) }
func (l Lines) OtherAccountUnknown(vendor string) string {
	return fmt.Sprintf(l.OtherAccountUnknownF, vendor)
}
func (l Lines) OtherAccount(vendor string) string { return fmt.Sprintf(l.OtherAccountF, vendor) }
func (l Lines) ViewTitle(site string) string      { return fmt.Sprintf(l.ViewTitleF, site) }
func (l Lines) ViewNote(site string, minutes int) string {
	return fmt.Sprintf(l.ViewNoteF, site, minutes)
}
func (l Lines) ViewNoteAny(minutes int) string { return fmt.Sprintf(l.ViewNoteAnyF, minutes) }
func (l Lines) ViewLapsed(minutes int) string  { return fmt.Sprintf(l.ViewLapsedF, minutes) }

// Languages are the languages planque speaks, English first: the ones
// tafs ships.
var Languages = []string{"en", "de", "fr", "es", "it", "pt", "nl", "pl"}

// For is lang's lines, English's for a language planque doesn't speak.
func For(lang string) Lines {
	if l, ok := catalog[lang]; ok {
		return l
	}
	return catalog["en"]
}

// Pick is the language a page speaks: stored (the language the app that
// started the connection asked in) when planque speaks it, else the first
// of acceptLanguage (an Accept-Language header, in its order) that planque
// speaks, else English. A tag counts by its language ("pt-PT" is "pt").
func Pick(stored, acceptLanguage string) string {
	if code := Code(stored); code != "" {
		return code
	}
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		if refused(params) {
			continue
		}
		if code := Code(tag); code != "" {
			return code
		}
	}
	return "en"
}

// refused says whether a language's parameters give it q=0: not acceptable.
func refused(params string) bool {
	for _, p := range strings.Split(params, ";") {
		k, v, ok := strings.Cut(strings.TrimSpace(p), "=")
		if !ok || strings.TrimSpace(k) != "q" {
			continue
		}
		q, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return err == nil && q == 0
	}
	return false
}

// Code is tag's language when planque speaks it ("de-CH" is "de"), or "".
func Code(tag string) string {
	code, _, _ := strings.Cut(strings.ToLower(strings.TrimSpace(tag)), "-")
	code, _, _ = strings.Cut(code, "_")
	if _, ok := catalog[code]; ok {
		return code
	}
	return ""
}
