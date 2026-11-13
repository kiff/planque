package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/store"
)

// apiReadsADay is the most API descriptions an account reads a day.
const apiReadsADay = 100

var errAPIReads = errors.New("the account's API description reads for the day are used")

// countAPIRead takes one of the day's description reads for the account, or
// says they're used (errAPIReads), as the on-demand set-up allowance counts.
func (s *Server) countAPIRead(ctx context.Context, tenant string) error {
	now := s.now().UTC()
	key := store.Item{PK: tenantPK(tenant), SK: "APIREADS#" + now.Format("2006-01-02"), TTL: now.Add(48 * time.Hour).Unix()}
	for range 5 {
		it, ok, err := s.KV.Get(ctx, key.PK, key.SK)
		if err != nil {
			return err
		}
		var a allowance
		if ok && json.Unmarshal(it.Doc, &a) != nil {
			a = allowance{}
		}
		if a.Count >= apiReadsADay {
			return errAPIReads
		}
		a.Count++
		cond := store.Cond{NotExists: true}
		if ok {
			ver := it.Ver
			cond = store.Cond{Ver: &ver}
			key.Ver = it.Ver + 1
		}
		if err = putDoc(ctx, s.KV, key, a, cond); !errors.Is(err, store.ErrCondition) {
			return err
		}
	}
	return errAPIReads
}

// apiOperations reads a vendor's public API description before any grant,
// POST /v1/accounts/{tenant}/api/operations {api: {url, auth}}: its
// operations as planque
// names them, so a key request can name the job's tools and be declared, and
// the key and the taf's Card come in one approval. It reads the description
// document only, as anyone could: no key, no call to the vendor's API, none of
// the account's data. Same parser, limits and refusals as a request; it keeps
// nothing: no grant, no key, no snapshot.
func (s *Server) apiOperations(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	if _, err := s.agentKey(r, tenant); err != nil {
		s.fail(w, r, err)
		return
	}
	body, err := readBody(r)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var req struct {
		API apiAsk `json:"api"`
		// Q is words from the job, to find a big API's few operations by.
		Q string `json:"q"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "the body is {\"api\": {\"url\", \"auth\"}}")
		return
	}
	if len(req.Q) > 200 {
		writeErr(w, http.StatusBadRequest, "q is a few words from the job, at most 200 characters")
		return
	}
	a := req.API
	a.URL = strings.TrimSpace(a.URL)
	if a.URL == "" {
		writeErr(w, http.StatusBadRequest, "api.url is required: the vendor's OpenAPI description, in JSON or YAML")
		return
	}
	// Each read fetches a description of up to 16 MB: counted per account a day, once the ask is one.
	if err := s.countAPIRead(r.Context(), tenant); err != nil {
		if errors.Is(err, errAPIReads) {
			writeErr(w, http.StatusTooManyRequests, "This account has read "+strconv.Itoa(apiReadsADay)+" API descriptions today. Try again tomorrow.")
			return
		}
		s.fail(w, r, err)
		return
	}
	raw, err := s.REST.Fetch(r.Context(), a.URL)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "Couldn't read the API description at "+rest.ShownURL(a.URL)+": "+err.Error())
		return
	}
	// Listing keeps nothing, so it reads past a snapshot's 500: a request still takes at most 500.
	snap, err := openapi.ParseToList(raw, a.Auth)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	// Names, titles and read or write: what Coach picks the job's tools by.
	type operation struct {
		Name     string `json:"name"`
		Title    string `json:"title,omitempty"`
		ReadOnly bool   `json:"read_only"`
	}
	answer := map[string]any{"title": label(snap.Title), "servers": snap.Servers, "accounts": snap.Accounts, "total": len(snap.Operations)}
	found := snap.Operations
	if q := strings.TrimSpace(req.Q); q != "" {
		found = matching(snap.Operations, q)
		answer["matched"] = len(found)
		if len(found) > maxFound {
			found = found[:maxFound]
		}
	} else if len(found) > maxListed {
		// A big API, with no words: its count, never a list Coach would spend a turn reading.
		answer["operations"] = []operation{}
		answer["more"] = strconv.Itoa(len(found)) + " operations: describe again with q, words from the job (\"refund charge\"), for the few it needs"
		writeJSON(w, http.StatusOK, answer)
		return
	}
	ops := make([]operation, 0, len(found))
	for _, o := range found {
		ops = append(ops, operation{Name: o.Name, Title: label(o.Summary), ReadOnly: o.ReadOnly})
	}
	answer["operations"] = ops
	if len(snap.Skipped) <= maxListed {
		answer["skipped"] = snap.Skipped
	}
	writeJSON(w, http.StatusOK, answer)
}

// maxListed is the most operations listed with no words; maxFound the most
// that words find.
const (
	maxListed = 200
	maxFound  = 50
)

// matching are the operations the words find, best first: by how many words
// its name holds, then how many its title does, so no summary, however many
// of the job's words it's stuffed with, ranks above an operation whose name
// holds more of them. Words under three letters, and a few with no meaning
// in an operation (the, and, for…), are ignored; an operation no word finds
// isn't listed.
func matching(ops []openapi.Operation, q string) []openapi.Operation {
	var words []string
	for _, w := range strings.FieldsFunc(strings.ToLower(q), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		if len(w) >= 3 && !stopWords[w] && !slices.Contains(words, w) {
			words = append(words, w)
		}
	}
	type scored struct {
		op          openapi.Operation
		name, title int
	}
	var hits []scored
	for _, o := range ops {
		name, title := strings.ToLower(o.Name), strings.ToLower(o.Summary)
		if len(title) > 200 {
			title = title[:200]
		}
		h := scored{op: o}
		for _, w := range words {
			if strings.Contains(name, w) {
				h.name++
			}
			if strings.Contains(title, w) {
				h.title++
			}
		}
		if h.name+h.title > 0 {
			hits = append(hits, h)
		}
	}
	slices.SortStableFunc(hits, func(a, b scored) int {
		if a.name != b.name {
			return b.name - a.name
		}
		return b.title - a.title
	})
	out := make([]openapi.Operation, len(hits))
	for i, h := range hits {
		out[i] = h.op
	}
	return out
}

// stopWords are words of three letters or more that name no operation.
var stopWords = map[string]bool{"the": true, "and": true, "for": true, "with": true, "from": true, "that": true, "this": true,
	"into": true, "its": true, "our": true, "your": true, "all": true, "any": true, "are": true, "was": true, "has": true}
