package planque

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Proof says an app connects today: a check connected it (a phone proof, an
// Eval), the way How says ("key", with the app's own key, or "sign-in"), and
// the app has been listed as working since At (YYYY-MM-DD), where Link shows
// it. It's per app, as kiff.dev's apps list has it, not per catalog entry: an
// entry that reaches the app another way carries the same proof, and How says
// which way was run. It's reviewed data (proved.json), never code: proving an
// app is one data change, and every reader (tafs cloud's gallery, kiff.dev's
// apps list) reads the same flag on the catalog.
type Proof struct {
	At   string `json:"at"`
	How  string `json:"how"`
	Link string `json:"link,omitempty"`
}

//go:embed proved.json
var provedJSON []byte

// Proved is proved.json, keyed by the app's folded name (provedKey). A file
// that doesn't read stops the build's tests, never a request.
var Proved = mustDecodeProved(provedJSON)

var provedName = regexp.MustCompile(`^[a-z0-9]{1,64}$`)

// DecodeProved reads a proofs file: each app's folded name, the way it was
// checked, the day it's listed as working since, and a public link to where
// that's recorded.
func DecodeProved(raw []byte) (map[string]Proof, error) {
	var rows map[string]Proof
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&rows) != nil || dec.Decode(new(any)) != io.EOF || rows == nil {
		return nil, fmt.Errorf("the proofs must be an object of apps, each with at and link")
	}
	for name, p := range rows {
		if !provedName.MatchString(name) || provedKey(name) != name {
			return nil, fmt.Errorf("proofs: %q isn't an app's folded name (lower case letters and digits)", name)
		}
		if p.How != "key" && p.How != "sign-in" {
			return nil, fmt.Errorf("proofs: %s's how %q is neither key nor sign-in", name, p.How)
		}
		if _, err := time.Parse(time.DateOnly, p.At); err != nil {
			return nil, fmt.Errorf("proofs: %s's at %q isn't a day (YYYY-MM-DD)", name, p.At)
		}
		if p.Link != "" {
			if u, err := url.Parse(p.Link); err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
				return nil, fmt.Errorf("proofs: %s's link isn't an https URL", name)
			}
		}
	}
	return rows, nil
}

func mustDecodeProved(raw []byte) map[string]Proof {
	p, err := DecodeProved(raw)
	if err != nil {
		panic(err)
	}
	return p
}

// provedKey folds an app's name as the gallery's recipes name it: lower
// case, without a leading underscore, with dots, spaces, dashes and
// underscores dropped (Google Calendar, google-calendar and googlecalendar
// are one app).
func provedKey(name string) string {
	name = strings.TrimLeft(strings.ToLower(name), "_")
	return strings.NewReplacer(" ", "", ".", "", "-", "", "_", "", "\t", "").Replace(name)
}

// proofOf is the proof of the app an entry names, by its slug: after
// connector:<provider>:, or before @ in a reviewed entry. Never by its vendor,
// which is the provider's words: a name planque didn't choose proves nothing.
func proofOf(proved map[string]Proof, entry string) *Proof {
	slug := entry
	if _, s, ok := connectorEntry(entry); ok {
		slug = s
	} else if at := strings.LastIndex(entry, "@"); at > 0 {
		slug = entry[:at]
	}
	if p, ok := proved[provedKey(slug)]; ok && slug != "" {
		return &p
	}
	return nil
}
