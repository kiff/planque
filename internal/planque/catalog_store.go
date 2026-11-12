package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/jcs"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/store"
)

// The catalog as data: the verified tier. An entry is drafted and reviewed
// by an operator's API credential, and served only once an operator's
// enrolled device signed it, never the one who drafted or last changed it.
// Each revision (id@n) is immutable once published: a change is a new
// revision. planque serves an entry only while its stored digest is the one
// its publish signature covers. Each publish and retire is kept as a record
// that is never edited. Entries shipped in PLANQUE_CATALOG_FILE are served
// as before, until they are published here.

// catalogPK holds the catalog's rows: ENTRY#<id>, RECORD#<at>#<id> and
// NONCE#<nonce>.
const catalogPK = "CATALOG"

// The states of a catalog revision.
const (
	CatalogDraft     = "draft"
	CatalogReviewed  = "reviewed"
	CatalogPublished = "published"
	CatalogRetired   = "retired"
)

// catalogSignWindow is how old a publish or retire signature may be.
const catalogSignWindow = 5 * time.Minute

// catalogCacheFor is how long served entries are kept before the table is
// read again.
const catalogCacheFor = time.Minute

// CatalogOperator is a person allowed to work on the catalog: their API
// credential (drafting, changing, reviewing), by its SHA-256, and their
// enrolled device (publishing, retiring), catalog-only and apart from any
// tenant's phone. Either may be absent: an agent's credential has no device.
type CatalogOperator struct {
	ID          string `json:"id"`
	TokenSHA256 string `json:"token_sha256,omitempty"`
	DeviceKey   string `json:"device_key,omitempty"`
	key         []byte
}

var catalogOperatorID = regexp.MustCompile(`^op_[a-z0-9_-]{1,40}$`)

// DecodeCatalogOperators reads PLANQUE_CATALOG_OPERATORS: a JSON list.
func DecodeCatalogOperators(raw []byte) (map[string]CatalogOperator, error) {
	var list []CatalogOperator
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&list) != nil || dec.Decode(new(any)) != io.EOF || len(list) > 50 {
		return nil, errors.New("catalog operators must be a list of at most 50 {id, token_sha256, device_key}")
	}
	out := map[string]CatalogOperator{}
	for _, op := range list {
		if !catalogOperatorID.MatchString(op.ID) || out[op.ID].ID != "" {
			return nil, errors.New("a catalog operator needs a unique id op_…")
		}
		if op.TokenSHA256 == "" && op.DeviceKey == "" {
			return nil, fmt.Errorf("catalog operator %s has neither a credential nor a device", op.ID)
		}
		if op.TokenSHA256 != "" {
			if sum, err := hex.DecodeString(op.TokenSHA256); err != nil || len(sum) != 32 || op.TokenSHA256 != strings.ToLower(op.TokenSHA256) {
				return nil, fmt.Errorf("catalog operator %s: token_sha256 must be 64 lowercase hex", op.ID)
			}
		}
		if op.DeviceKey != "" {
			key, err := device.ParseKey(op.DeviceKey)
			if err != nil {
				return nil, fmt.Errorf("catalog operator %s: %w", op.ID, err)
			}
			op.key = key
		}
		out[op.ID] = op
	}
	return out, nil
}

// catalogRow is one revision as kept.
type catalogRow struct {
	ID        string          `json:"id"`
	Entry     json.RawMessage `json:"entry"`  // the entry as drafted, with its id
	Digest    string          `json:"digest"` // the JCS SHA-256 of Entry
	State     string          `json:"state"`
	DraftedBy string          `json:"drafted_by"`
	ChangedBy string          `json:"changed_by"`
	Changed   string          `json:"changed"`
	Published *catalogSigned  `json:"published,omitempty"`
	Retired   *catalogSigned  `json:"retired,omitempty"`
}

// catalogSigned is a publish or retire as signed.
type catalogSigned struct {
	Operator  string `json:"operator"`
	Device    string `json:"device"` // the key's fingerprint
	Digest    string `json:"digest"`
	Nonce     string `json:"nonce"`
	SignedAt  string `json:"signed_at"`
	Signature string `json:"signature"`
}

// CatalogSignMessage is what an operator's device signs to publish or
// retire a revision of the catalog at the planque at base.
func CatalogSignMessage(base, action, id, digest, nonce, signedAt string) []byte {
	return []byte(strings.Join([]string{"planque.catalog.v1", strings.TrimRight(base, "/"), action, id, digest, nonce, signedAt}, "\n"))
}

var catalogNonce = regexp.MustCompile(`^[A-Za-z0-9_-]{16,64}$`)

// entryDocument is the entry as signed: its fields, with its id.
func entryDocument(e CatalogEntry) (json.RawMessage, string, error) {
	raw, err := json.Marshal(e)
	if err != nil {
		return nil, "", err
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, "", err
	}
	doc["id"], _ = json.Marshal(e.ID)
	raw, _ = json.Marshal(doc)
	digest, err := jcs.SHA256(raw)
	return raw, digest, err
}

// decodeEntryDocument reads an entry as sent or kept: its fields and its id.
func decodeEntryDocument(raw []byte) (CatalogEntry, error) {
	var doc map[string]json.RawMessage
	if json.Unmarshal(raw, &doc) != nil {
		return CatalogEntry{}, errors.New("the entry must be a JSON object")
	}
	var id string
	if json.Unmarshal(doc["id"], &id) != nil {
		return CatalogEntry{}, errors.New("the entry needs its id, name@revision")
	}
	delete(doc, "id")
	rest, _ := json.Marshal(doc)
	var e CatalogEntry
	dec := json.NewDecoder(bytes.NewReader(rest))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&e); err != nil {
		return CatalogEntry{}, errors.New("the entry has a field a catalog entry doesn't take")
	}
	e.ID = id
	return e, nil
}

// catalogCache is the table's part of what's served: its published
// revisions as signed, and the ids it retired. It is read at most once a
// minute.
type catalogCache struct {
	mu        sync.Mutex
	at        time.Time
	read      bool
	published map[string]CatalogEntry
	retired   []string
}

// catalogEntries is every entry served: the shipped ones, then each
// published revision whose stored digest is the one its signature covers,
// less the ids retired here. A table that can't be read serves what was
// read last.
func (s *Server) catalogEntries() map[string]CatalogEntry {
	if s.KV == nil {
		return s.CatalogEntries
	}
	c := &s.catalogCache
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.read || s.now().Sub(c.at) >= catalogCacheFor {
		s.readCatalogTable(c)
	}
	out := make(map[string]CatalogEntry, len(s.CatalogEntries)+len(c.published))
	for id, e := range s.CatalogEntries {
		out[id] = e
	}
	for _, id := range c.retired {
		delete(out, id) // retired here, whatever was shipped
	}
	for id, e := range c.published {
		out[id] = e
	}
	return out
}

func (s *Server) readCatalogTable(c *catalogCache) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := s.catalogRows(ctx)
	if err != nil {
		s.log().Warn("catalog not read", "err", err.Error())
		return
	}
	published, retired := map[string]CatalogEntry{}, []string{}
	for _, row := range rows {
		switch row.State {
		case CatalogRetired:
			retired = append(retired, row.ID)
			continue
		case CatalogPublished:
		default:
			continue
		}
		e, err := decodeEntryDocument(row.Entry)
		_, digest, derr := entryDocument(e)
		if err != nil || derr != nil || e.ID != row.ID || digest != row.Digest || row.Published == nil || row.Published.Digest != digest {
			s.log().Warn("catalog entry not served: not as signed", "entry", row.ID)
			continue
		}
		published[row.ID] = e
	}
	c.published, c.retired, c.at, c.read = published, retired, s.now(), true
}

// forgetCatalog makes the next read see the table.
func (s *Server) forgetCatalog() {
	s.catalogCache.mu.Lock()
	s.catalogCache.read = false
	s.catalogCache.mu.Unlock()
}

func (s *Server) catalogRows(ctx context.Context) ([]catalogRow, error) {
	items, err := s.KV.Query(ctx, catalogPK, "ENTRY#")
	if err != nil {
		return nil, err
	}
	rows := make([]catalogRow, 0, len(items))
	for _, it := range items {
		var row catalogRow
		if json.Unmarshal(it.Doc, &row) == nil && "ENTRY#"+row.ID == it.SK {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (s *Server) catalogRow(ctx context.Context, id string) (*catalogRow, int64, error) {
	it, ok, err := s.KV.Get(ctx, catalogPK, "ENTRY#"+id)
	if err != nil || !ok {
		return nil, 0, err
	}
	var row catalogRow
	if err := json.Unmarshal(it.Doc, &row); err != nil {
		return nil, 0, err
	}
	return &row, it.Ver, nil
}

func (s *Server) putCatalogRow(ctx context.Context, row *catalogRow, ver *int64) error {
	doc, _ := json.Marshal(row)
	cond := store.Cond{NotExists: ver == nil, Ver: ver}
	return s.KV.Put(ctx, store.Item{PK: catalogPK, SK: "ENTRY#" + row.ID, Doc: doc}, cond)
}

// catalogOperator is the operator whose credential the request carries.
func (s *Server) catalogOperator(r *http.Request) (CatalogOperator, bool) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		return CatalogOperator{}, false
	}
	sum := sha256.Sum256([]byte(token))
	got := hex.EncodeToString(sum[:])
	for _, op := range s.CatalogOperators {
		if op.TokenSHA256 != "" && subtle.ConstantTimeCompare([]byte(op.TokenSHA256), []byte(got)) == 1 {
			return op, true
		}
	}
	return CatalogOperator{}, false
}

// catalogRowView is a revision as an operator reads it: the entry, its
// digest to sign, its state and who did what.
func catalogRowView(row *catalogRow) map[string]any {
	return map[string]any{"id": row.ID, "entry": row.Entry, "digest": row.Digest, "state": row.State,
		"drafted_by": row.DraftedBy, "changed_by": row.ChangedBy, "changed": row.Changed,
		"published": row.Published, "retired": row.Retired}
}

// draftCatalogEntry: POST /v1/catalog/entries {entry}. A new revision, as
// a draft. A revision that exists isn't replaced: a change after publish is
// a new revision.
func (s *Server) draftCatalogEntry(w http.ResponseWriter, r *http.Request) {
	op, ok := s.catalogOperator(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "a catalog operator's credential is required")
		return
	}
	e, msg := s.readCatalogEntry(r)
	if msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	doc, digest, err := entryDocument(e)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if _, shipped := s.CatalogEntries[e.ID]; shipped && !s.shippedMatches(e.ID, digest) {
		writeErr(w, http.StatusConflict, e.ID+" is shipped with different content: draft a new revision")
		return
	}
	row := &catalogRow{ID: e.ID, Entry: doc, Digest: digest, State: CatalogDraft, DraftedBy: op.ID, ChangedBy: op.ID, Changed: s.now().UTC().Format(time.RFC3339)}
	if err := s.putCatalogRow(r.Context(), row, nil); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, e.ID+" exists: a revision is never replaced, draft "+nextRevision(e.ID))
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, catalogRowView(row))
}

// changeCatalogEntry: PUT /v1/catalog/entries/{id} {entry}. A draft or a
// reviewed revision changes, and is a draft again, changed by this operator.
func (s *Server) changeCatalogEntry(w http.ResponseWriter, r *http.Request) {
	op, ok := s.catalogOperator(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "a catalog operator's credential is required")
		return
	}
	id := r.PathValue("id")
	row, ver, err := s.catalogRow(r.Context(), id)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if row == nil {
		writeErr(w, http.StatusNotFound, "no catalog revision "+id)
		return
	}
	if row.State != CatalogDraft && row.State != CatalogReviewed {
		writeErr(w, http.StatusConflict, id+" is "+row.State+": it never changes, draft "+nextRevision(id))
		return
	}
	e, msg := s.readCatalogEntry(r)
	if msg == "" && e.ID != id {
		msg = "the entry's id must be " + id
	}
	if msg != "" {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	doc, digest, err := entryDocument(e)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	row.Entry, row.Digest, row.State, row.ChangedBy, row.Changed = doc, digest, CatalogDraft, op.ID, s.now().UTC().Format(time.RFC3339)
	if err := s.putCatalogRow(r.Context(), row, &ver); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, id+" changed meanwhile: read it again")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, catalogRowView(row))
}

// reviewCatalogEntry: POST /v1/catalog/entries/{id}/review {digest}. The
// draft as read is reviewed: ready for a device to sign.
func (s *Server) reviewCatalogEntry(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.catalogOperator(r); !ok {
		writeErr(w, http.StatusUnauthorized, "a catalog operator's credential is required")
		return
	}
	var in struct {
		Digest string `json:"digest"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&in) != nil {
		writeErr(w, http.StatusBadRequest, "send {digest}: the revision's digest as read")
		return
	}
	id := r.PathValue("id")
	row, ver, err := s.catalogRow(r.Context(), id)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	switch {
	case row == nil:
		writeErr(w, http.StatusNotFound, "no catalog revision "+id)
		return
	case row.State != CatalogDraft:
		writeErr(w, http.StatusConflict, id+" is "+row.State+", not a draft")
		return
	case in.Digest != row.Digest:
		writeErr(w, http.StatusConflict, id+" changed since it was read: read it again")
		return
	}
	row.State = CatalogReviewed
	if err := s.putCatalogRow(r.Context(), row, &ver); errors.Is(err, store.ErrCondition) {
		writeErr(w, http.StatusConflict, id+" changed meanwhile: read it again")
		return
	} else if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, catalogRowView(row))
}

// signCatalogEntry: POST /v1/catalog/entries/{id}/publish or …/retire,
// {operator, digest, nonce, signed_at, signature}. The operator's enrolled
// device signed CatalogSignMessage. Publishing takes a reviewed revision
// that no write could make final unpinned, from an operator who neither
// drafted nor last changed it; retiring takes a published one.
func (s *Server) signCatalogEntry(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Operator  string `json:"operator"`
			Digest    string `json:"digest"`
			Nonce     string `json:"nonce"`
			SignedAt  string `json:"signed_at"`
			Signature string `json:"signature"`
		}
		if json.NewDecoder(io.LimitReader(r.Body, 8192)).Decode(&in) != nil || !catalogNonce.MatchString(in.Nonce) {
			writeErr(w, http.StatusBadRequest, "send {operator, digest, nonce, signed_at, signature}: a nonce of 16 to 64 letters, digits, _ or -")
			return
		}
		op, ok := s.CatalogOperators[in.Operator]
		if !ok || op.key == nil {
			writeErr(w, http.StatusForbidden, "no enrolled catalog device for "+in.Operator)
			return
		}
		id := r.PathValue("id")
		msg := CatalogSignMessage(s.BaseURL, action, id, in.Digest, in.Nonce, in.SignedAt)
		if err := device.Verify(op.key, msg, in.SignedAt, in.Signature, s.now()); err != nil {
			writeErr(w, http.StatusForbidden, err.Error())
			return
		}
		row, ver, err := s.catalogRow(r.Context(), id)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		if row == nil {
			writeErr(w, http.StatusNotFound, "no catalog revision "+id)
			return
		}
		if in.Digest != row.Digest {
			writeErr(w, http.StatusConflict, id+" isn't what was signed: read it again")
			return
		}
		switch action {
		case "publish":
			if row.State != CatalogReviewed {
				writeErr(w, http.StatusConflict, id+" is "+row.State+": only a reviewed revision is published")
				return
			}
			if op.ID == row.DraftedBy || op.ID == row.ChangedBy {
				writeErr(w, http.StatusForbidden, "who drafted or last changed "+id+" can't publish it")
				return
			}
			e, err := decodeEntryDocument(row.Entry)
			if err != nil {
				s.fail(w, r, err)
				return
			}
			if why, err := s.unpinnedFinal(r.Context(), e); err != nil {
				s.fail(w, r, err)
				return
			} else if why != "" {
				writeErr(w, http.StatusConflict, why)
				return
			}
		case "retire":
			if row.State != CatalogPublished {
				writeErr(w, http.StatusConflict, id+" is "+row.State+": only a published revision is retired")
				return
			}
		}
		// The nonce is spent first: a signature counts once, even if what follows fails.
		nonce := store.Item{PK: catalogPK, SK: "NONCE#" + in.Nonce, Doc: []byte(`{}`), TTL: s.now().Add(2 * catalogSignWindow).Unix()}
		if err := s.KV.Put(r.Context(), nonce, store.Cond{NotExists: true}); errors.Is(err, store.ErrCondition) {
			writeErr(w, http.StatusConflict, "this signature was used already")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		signed := &catalogSigned{Operator: op.ID, Device: device.Fingerprint(op.key), Digest: in.Digest, Nonce: in.Nonce, SignedAt: in.SignedAt, Signature: in.Signature}
		if action == "publish" {
			row.State, row.Published = CatalogPublished, signed
		} else {
			row.State, row.Retired = CatalogRetired, signed
		}
		if err := s.putCatalogRow(r.Context(), row, &ver); errors.Is(err, store.ErrCondition) {
			writeErr(w, http.StatusConflict, id+" changed meanwhile: read it again")
			return
		} else if err != nil {
			s.fail(w, r, err)
			return
		}
		at := s.now().UTC().Format(time.RFC3339Nano)
		record, _ := json.Marshal(map[string]any{"action": action, "id": id, "at": at, "signed": signed})
		if err := s.KV.Put(r.Context(), store.Item{PK: catalogPK, SK: "RECORD#" + at + "#" + id, Doc: record}, store.Cond{NotExists: true}); err != nil {
			s.log().ErrorContext(r.Context(), "catalog record not kept", "entry", id, "action", action, "err", err.Error())
		}
		s.log().InfoContext(r.Context(), "catalog "+action, "entry", id, "operator", op.ID, "digest", in.Digest)
		s.forgetCatalog()
		writeJSON(w, http.StatusOK, catalogRowView(row))
	}
}

// listCatalogEntries: GET /v1/catalog/entries, every revision and its state.
func (s *Server) listCatalogEntries(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.catalogOperator(r); !ok {
		writeErr(w, http.StatusUnauthorized, "a catalog operator's credential is required")
		return
	}
	rows, err := s.catalogRows(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	slices.SortFunc(rows, func(a, b catalogRow) int { return strings.Compare(a.ID, b.ID) })
	out := make([]map[string]any, 0, len(rows))
	for i := range rows {
		out = append(out, catalogRowView(&rows[i]))
	}
	writeJSON(w, http.StatusOK, map[string]any{"entries": out})
}

// getCatalogEntry: GET /v1/catalog/entries/{id}, one revision, with the
// digest a device signs. It takes no credential: the device that signs it
// reads exactly what it signs, and an entry holds no secret.
func (s *Server) getCatalogEntry(w http.ResponseWriter, r *http.Request) {
	row, _, err := s.catalogRow(r.Context(), r.PathValue("id"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if row == nil {
		writeErr(w, http.StatusNotFound, "no catalog revision "+r.PathValue("id"))
		return
	}
	writeJSON(w, http.StatusOK, catalogRowView(row))
}

// readCatalogEntry reads {entry} and checks it as the shipped catalog is.
func (s *Server) readCatalogEntry(r *http.Request) (CatalogEntry, string) {
	var in struct {
		Entry json.RawMessage `json:"entry"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&in) != nil || len(in.Entry) == 0 {
		return CatalogEntry{}, "send {entry}: the entry, with its id"
	}
	e, err := decodeEntryDocument(in.Entry)
	if err != nil {
		return CatalogEntry{}, err.Error()
	}
	if err := e.validate(s.Applications); err != nil {
		return CatalogEntry{}, e.ID + ": " + err.Error()
	}
	return e, ""
}

// shippedMatches says whether a shipped entry is the one with digest, so
// it can be published here as it is.
func (s *Server) shippedMatches(id, digest string) bool {
	_, d, err := entryDocument(s.CatalogEntries[id])
	return err == nil && d == digest
}

// nextRevision is id with its revision one up.
func nextRevision(id string) string {
	name, rev, _ := strings.Cut(id, "@")
	n := 0
	_, _ = fmt.Sscanf(rev, "%d", &n)
	return fmt.Sprintf("%s@%d", name, n+1)
}

// finalValues are the words that make a write final: done, and can't be
// taken back.
var finalValues = []string{"final", "finalized", "finalised", "issued", "issue", "sent", "send", "posted", "post", "confirmed", "confirm", "published", "publish", "submitted", "submit", "approved", "paid", "booked", "closed", "completed", "unpaid", "validated", "committed"}

// undoableValues are the values a write can be pinned to instead, in order.
var undoableValues = []string{"draft", "pending", "unsent", "saved", "proposal", "quote"}

// FinalArgument is a write's argument whose value makes the write final,
// and the value that keeps it undoable.
type FinalArgument struct {
	Tool     string
	Argument string
	Undo     any
}

// FinalArguments are the arguments of op that make it final: an enum with
// a final value and an undoable one, or a boolean named for a final step.
// A read has none.
func FinalArguments(op openapi.Operation) []FinalArgument {
	if op.ReadOnly || len(op.InputSchema) == 0 {
		return nil
	}
	var schema struct {
		Properties map[string]struct {
			Type        string `json:"type"`
			Enum        []any  `json:"enum"`
			Description string `json:"description"`
		} `json:"properties"`
	}
	if json.Unmarshal(op.InputSchema, &schema) != nil {
		return nil
	}
	var out []FinalArgument
	for name, p := range schema.Properties {
		if p.Type == "boolean" {
			n := strings.ToLower(name)
			for _, w := range []string{"send", "final", "issue", "confirm", "publish", "submit", "post"} {
				if strings.Contains(n, w) {
					out = append(out, FinalArgument{Tool: op.Name, Argument: name, Undo: false})
					break
				}
			}
			continue
		}
		var final bool
		var undo any
		for _, v := range p.Enum {
			s, _ := v.(string)
			s = strings.ToLower(s)
			if slices.Contains(finalValues, s) {
				final = true
			}
			for _, u := range undoableValues {
				if s == u && undo == nil {
					undo = v
				}
			}
		}
		if final && undo != nil {
			out = append(out, FinalArgument{Tool: op.Name, Argument: name, Undo: undo})
		}
	}
	slices.SortFunc(out, func(a, b FinalArgument) int { return strings.Compare(a.Argument, b.Argument) })
	return out
}

// unpinnedFinal says why e can't be published while one of its writes can
// be made final by an argument it doesn't pin to the undoable value: "" when
// it can. A key entry's description is read, and must be the one it pins.
// An MCP entry's tools are listed by its server once connected: nothing is
// read here.
func (s *Server) unpinnedFinal(ctx context.Context, e CatalogEntry) (string, error) {
	if e.Kind != KindAPIKey || e.API == nil {
		return "", nil
	}
	raw, err := s.REST.Fetch(ctx, e.API.URL)
	if err != nil {
		return e.ID + ": its description couldn't be read to check its writes: " + err.Error(), nil
	}
	if sum := sha256.Sum256(raw); hex.EncodeToString(sum[:]) != e.DescriptionSHA256 {
		return e.ID + ": its description changed since it was pinned", nil
	}
	snap, err := openapi.Parse(raw, e.API.Auth)
	if err != nil {
		return e.ID + ": " + err.Error(), nil
	}
	for _, t := range e.Tools {
		i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == t.Tool })
		if i < 0 || t.ReadOnly {
			continue
		}
		for _, f := range FinalArguments(snap.Operations[i]) {
			want, _ := json.Marshal(f.Undo)
			if got, pinned := t.Fixed[f.Argument]; !pinned || !sameJSON(got, want) {
				return fmt.Sprintf("%s: %s's %s can make it final; pin it to %s in fixed_arguments, or leave the tool out", e.ID, t.Tool, f.Argument, want), nil
			}
		}
	}
	return "", nil
}

func sameJSON(a, b []byte) bool {
	ca, err1 := jcs.Canonical(a)
	cb, err2 := jcs.Canonical(b)
	return err1 == nil && err2 == nil && bytes.Equal(ca, cb)
}
