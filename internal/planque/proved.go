package planque

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"
)

// Proof says an app connects today: a check connected it (a phone proof, an
// Eval), the way How says, in the catalog's words ("key", with the app's own
// key, or "oauth", a sign-in), and the app has been listed as working since
// At (YYYY-MM-DD), where Link shows it. A catalog entry carries it only when
// it's the way that was run: its own How is the proof's, and it names the app
// (its slug, or one of Entries, the reviewed entries that reach it, such as
// qonto-key for Qonto). An app proved by key with no key entry is proved on
// the catalog's by_key way alone (Server.catalog's proved). It's reviewed
// data (proved.json), never code: proving an app is one data change, and
// every reader (tafs cloud's gallery, kiff.dev's apps list) reads the same.
type Proof struct {
	At      string   `json:"at"`
	How     string   `json:"how"`
	Entries []string `json:"entries,omitempty"`
	// Route narrows a sign-in proof to the way it was run: "managed" (the
	// provider's own sign-in, connector:<provider>:<app>) or "own" (a reviewed
	// entry through tafs's own client, such as googlesheets@1). Empty reaches
	// both, as every proof did before.
	Route string `json:"route,omitempty"`
	Link  string `json:"link,omitempty"`
}

//go:embed proved.json
var provedJSON []byte

// Proved is proved.json, keyed by the app's folded name (provedKey). A file
// that doesn't read stops the build's tests, never a request.
var Proved = mustDecodeProved(provedJSON)

var (
	provedName   = regexp.MustCompile(`^[a-z0-9]{1,64}$`)
	reviewedName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
)

// DecodeProved reads a proofs file: each app's folded name, the way it was
// checked, the reviewed entries that reach it, the day it's listed as working
// since, and a public link to where that's recorded.
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
		if p.How != "key" && p.How != "oauth" {
			return nil, fmt.Errorf("proofs: %s's how %q is neither key nor oauth", name, p.How)
		}
		if p.Route != "" && p.Route != "managed" && p.Route != "own" {
			return nil, fmt.Errorf("proofs: %s's route %q is neither managed nor own", name, p.Route)
		}
		if p.Route != "" && p.How != "oauth" {
			return nil, fmt.Errorf("proofs: %s's route only narrows a sign-in (oauth) proof", name)
		}
		for _, e := range p.Entries {
			if !reviewedName.MatchString(e) {
				return nil, fmt.Errorf("proofs: %s's entry %q isn't a reviewed entry's name (without @)", name, e)
			}
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

// proofOf is the proof an entry carries: the app it names, by its slug (after
// connector:<provider>:, or before @ in a reviewed entry) or as one of a
// proof's Entries, proved the entry's own way (how). Never by its vendor,
// which is the provider's words: a name planque didn't choose proves nothing.
func proofOf(proved map[string]Proof, entry, how string) *Proof {
	slug := entry
	managed := false
	if _, s, ok := connectorEntry(entry); ok {
		slug, managed = s, true
	} else if at := strings.LastIndex(entry, "@"); at > 0 {
		slug = entry[:at]
	}
	if slug == "" {
		return nil
	}
	runWay := func(p Proof) bool {
		return p.How == how && (p.Route == "" || (p.Route == "managed") == managed)
	}
	if p, ok := proved[provedKey(slug)]; ok && runWay(p) {
		return &p
	}
	for _, p := range proved {
		if runWay(p) && slices.Contains(p.Entries, slug) {
			return &p
		}
	}
	return nil
}

// publicProof is what a reader with no account sees of a proof: the day the
// app has been listed as working since, the way it was run, and the public
// link to where that's recorded. Never which reviewed entries reach it.
type publicProof struct {
	At   string `json:"at"`
	How  string `json:"how"`
	Link string `json:"link,omitempty"`
}

// servePublicProofs answers GET /v1/proved with no account: {"proved": {app:
// {at, how, link}}}, by folded name, the same map the catalog's top level
// carries. kiff.dev's apps list reads it, so the site and tafs cloud's gallery
// read one source. It holds nothing but proved.json, which is public already.
func (s *Server) servePublicProofs(w http.ResponseWriter, r *http.Request) {
	// The catalog's own source (vendor_catalog.go), so the two never disagree.
	proved := s.provedMap(r.Context())
	out := make(map[string]publicProof, len(proved))
	for name, p := range proved {
		out[name] = publicProof{At: p.At, How: p.How, Link: p.Link}
	}
	w.Header().Set("Content-Type", "application/json")
	// Public data, re-read by kiff.dev every 10 minutes: caches may keep it as long.
	w.Header().Set("Cache-Control", "public, max-age=600")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"proved": out})
}
