package planque

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	// Each read fetches a description of up to 16 MB: counted per account a day.
	if err := s.countAPIRead(r.Context(), tenant); err != nil {
		if errors.Is(err, errAPIReads) {
			writeErr(w, http.StatusTooManyRequests, "This account has read "+strconv.Itoa(apiReadsADay)+" API descriptions today. Try again tomorrow.")
			return
		}
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
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "the body is {\"api\": {\"url\", \"auth\"}}")
		return
	}
	a := req.API
	a.URL = strings.TrimSpace(a.URL)
	if a.URL == "" {
		writeErr(w, http.StatusBadRequest, "api.url is required: the vendor's OpenAPI description, in JSON or YAML")
		return
	}
	raw, err := s.REST.Fetch(r.Context(), a.URL)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "Couldn't read the API description at "+rest.ShownURL(a.URL)+": "+err.Error())
		return
	}
	snap, err := openapi.Parse(raw, a.Auth)
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
	ops := make([]operation, 0, len(snap.Operations))
	for _, o := range snap.Operations {
		ops = append(ops, operation{Name: o.Name, Title: label(o.Summary), ReadOnly: o.ReadOnly})
	}
	writeJSON(w, http.StatusOK, map[string]any{"title": label(snap.Title), "servers": snap.Servers, "accounts": snap.Accounts,
		"operations": ops, "skipped": snap.Skipped})
}
