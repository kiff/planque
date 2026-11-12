package planque

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
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
		if err != nil || derr != nil || e.ID != row.ID || digest != row.Digest || !s.catalogSignedBy(row.Published, "publish", row.ID, digest) {
			s.log().Warn("catalog entry not served: not as signed", "entry", row.ID)
			continue
		}
		published[row.ID] = e
	}
	c.published, c.retired, c.at, c.read = published, retired, s.now(), true
}

// catalogSignedBy says whether sig is an enrolled operator's device signature,
// for action, over this revision with this digest. Its age isn't checked:
// it was within its window when planque took it.
func (s *Server) catalogSignedBy(sig *catalogSigned, action, id, digest string) bool {
	if sig == nil || sig.Digest != digest {
		return false
	}
	op, ok := s.CatalogOperators[sig.Operator]
	if !ok || op.key == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(op.key)
	k, isEC := pub.(*ecdsa.PublicKey)
	der, derr := base64.StdEncoding.DecodeString(sig.Signature)
	if err != nil || !isEC || derr != nil {
		return false
	}
	sum := sha256.Sum256(CatalogSignMessage(s.BaseURL, action, id, digest, sig.Nonce, sig.SignedAt))
	return ecdsa.VerifyASN1(k, sum[:], der)
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
		// The row's signed block is the record that counts: published and
		// retired are each set once and never edited. This log of them is
		// kept beside it, for reading in order.
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
	view := catalogRowView(row)
	if e, err := decodeEntryDocument(row.Entry); err == nil {
		view["described"] = s.describeEntry(r.Context(), e)
	}
	writeJSON(w, http.StatusOK, view)
}

// describedTool is one of an entry's tools as planque reads its pinned
// description, never as the entry labels it.
type describedTool struct {
	Tool     string           `json:"tool"`
	Method   string           `json:"method,omitempty"`
	Path     string           `json:"path,omitempty"`
	ReadOnly bool             `json:"read_only"`
	Final    []describedFinal `json:"final_arguments,omitempty"`
	Why      string           `json:"why,omitempty"` // why it can't be read
}

type describedFinal struct {
	Argument string `json:"argument"`
	Undo     any    `json:"undo"`
}

// describeEntry is what a device signing e draws for read or write: each
// tool as planque reads the description e pins, with the arguments that can
// make it final. It's outside the signed digest: the device checks the
// entry, and this tells it what the entry's tools do. A key entry's
// description is read and must match its pin; an MCP entry's tools are
// listed by its server once connected, so there's nothing to read here.
func (s *Server) describeEntry(ctx context.Context, e CatalogEntry) map[string]any {
	if e.Kind != KindAPIKey || e.API == nil {
		return map[string]any{"from": "server", "why": "an MCP server lists its tools once connected: nothing is read before"}
	}
	snap, why := s.pinnedDescription(ctx, e)
	if why != "" {
		return map[string]any{"from": "description", "why": why}
	}
	tools := make([]describedTool, 0, len(e.Tools))
	for _, t := range e.Tools {
		i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == t.Tool })
		if i < 0 {
			tools = append(tools, describedTool{Tool: t.Tool, Why: "not in its description"})
			continue
		}
		op := snap.Operations[i]
		d := describedTool{Tool: t.Tool, Method: op.Method, Path: op.Path, ReadOnly: op.ReadOnly}
		finals, err := FinalArguments(op)
		if err != nil {
			d.Why = err.Error()
		}
		for _, f := range finals {
			d.Final = append(d.Final, describedFinal{Argument: f.Argument, Undo: f.Undo})
		}
		tools = append(tools, d)
	}
	return map[string]any{"from": "description", "description_sha256": e.DescriptionSHA256, "tools": tools}
}

// describedCache is the parse of each pinned description, by its SHA-256 and
// the auth it was parsed with. Only a description that matched its pin is
// kept: it can't go stale, since a different file has a different hash. A
// fetch that failed or a description that changed is read again next time.
type describedCache struct {
	mu    sync.Mutex
	snaps map[string]*openapi.Snapshot
}

// maxDescribed bounds the parses kept; past it the cache starts again.
const maxDescribed = 64

// pinnedDescription is e's API description, parsed, once it matched its pin:
// from the cache when it was read before, else fetched, checked and kept.
// why says what stopped it.
func (s *Server) pinnedDescription(ctx context.Context, e CatalogEntry) (snap *openapi.Snapshot, why string) {
	auth, _ := json.Marshal(e.API.Auth)
	key := e.DescriptionSHA256 + "\x00" + string(auth)
	c := &s.described
	c.mu.Lock()
	snap = c.snaps[key]
	c.mu.Unlock()
	if snap != nil {
		return snap, ""
	}
	raw, err := s.REST.Fetch(ctx, e.API.URL)
	if err != nil {
		return nil, "its description couldn't be read: " + err.Error()
	}
	if sum := sha256.Sum256(raw); hex.EncodeToString(sum[:]) != e.DescriptionSHA256 {
		return nil, "its description changed since it was pinned"
	}
	parsed, err := openapi.Parse(raw, e.API.Auth)
	if err != nil {
		return nil, err.Error()
	}
	c.mu.Lock()
	if c.snaps == nil || len(c.snaps) >= maxDescribed {
		c.snaps = map[string]*openapi.Snapshot{}
	}
	c.snaps[key] = parsed
	c.mu.Unlock()
	return parsed, ""
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

// undoableValues are the values that keep a write undoable, in the order
// a pin prefers them. An enum that holds one is flagged, whatever its other
// values are: pinning to it is the safe direction, and a person can still
// leave the tool out.
var undoableValues = []string{"draft", "pending", "unsent", "saved", "proposal", "quote", "preview"}

// finalWords name a boolean that makes a write final when true; undoWords
// name one that keeps it undoable when true ("draft", "is_draft").
var (
	finalWords = []string{"send", "final", "issue", "confirm", "publish", "submit", "post", "commit", "approve", "pay"}
	undoWords  = []string{"draft", "preview"}
)

// FinalArgument is a write's argument whose value makes the write final,
// by its dotted path, and the value that keeps it undoable.
type FinalArgument struct {
	Tool     string
	Argument string
	Undo     any
}

// errArgumentsUnread: a write's arguments couldn't be read to check them.
var errArgumentsUnread = errors.New("its arguments couldn't be read to check them")

// FinalArguments are the arguments of op that can make it final, nested
// objects included, each by its dotted path (as fixed_arguments name them,
// up to maxPathParts parts): an enum that holds an undoable value and
// another, a boolean named for a final step (pinned false) or for the
// undoable state (pinned true). A read has none. A write whose arguments
// can't be read is an error: the check fails closed.
func FinalArguments(op openapi.Operation) ([]FinalArgument, error) {
	if op.ReadOnly || len(op.InputSchema) == 0 {
		return nil, nil
	}
	var root map[string]json.RawMessage
	if json.Unmarshal(op.InputSchema, &root) != nil {
		return nil, errArgumentsUnread
	}
	var out []FinalArgument
	if err := finalIn(op.Name, root, nil, &out); err != nil {
		return nil, err
	}
	slices.SortFunc(out, func(a, b FinalArgument) int { return strings.Compare(a.Argument, b.Argument) })
	return out, nil
}

// finalIn walks one object schema's properties, at path.
func finalIn(tool string, schema map[string]json.RawMessage, path []string, out *[]FinalArgument) error {
	raw, ok := schema["properties"]
	if !ok {
		return nil
	}
	var props map[string]map[string]json.RawMessage
	if json.Unmarshal(raw, &props) != nil {
		return errArgumentsUnread
	}
	for name, p := range props {
		at := append(slices.Clone(path), name)
		types, err := schemaTypes(p["type"])
		if err != nil {
			return err
		}
		switch {
		case slices.Contains(types, "object") || p["properties"] != nil:
			if len(at) < maxPathParts {
				if err := finalIn(tool, p, at, out); err != nil {
					return err
				}
			}
		case slices.Contains(types, "boolean"):
			n := strings.ToLower(name)
			if containsAny(n, undoWords) {
				*out = append(*out, FinalArgument{Tool: tool, Argument: strings.Join(at, "."), Undo: true})
			} else if containsAny(n, finalWords) {
				*out = append(*out, FinalArgument{Tool: tool, Argument: strings.Join(at, "."), Undo: false})
			}
		case p["enum"] != nil:
			var values []any
			if json.Unmarshal(p["enum"], &values) != nil {
				return errArgumentsUnread
			}
			var undo any
			others := 0
			for _, v := range values {
				s, isString := v.(string)
				if isString && undo == nil && slices.Contains(undoableValues, strings.ToLower(s)) {
					undo = v
					continue
				}
				if v != nil {
					others++
				}
			}
			if undo != nil && others > 0 {
				*out = append(*out, FinalArgument{Tool: tool, Argument: strings.Join(at, "."), Undo: undo})
			}
		}
	}
	return nil
}

// schemaTypes reads a JSON Schema type: a string, a list (OpenAPI 3.1's
// ["string", "null"]), or absent.
func schemaTypes(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var one string
	if json.Unmarshal(raw, &one) == nil {
		return []string{one}, nil
	}
	var many []string
	if json.Unmarshal(raw, &many) == nil {
		return many, nil
	}
	return nil, errArgumentsUnread
}

func containsAny(s string, words []string) bool {
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
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
		// Read or write is what the description says, never the entry's own label.
		i := slices.IndexFunc(snap.Operations, func(o openapi.Operation) bool { return o.Name == t.Tool })
		if i < 0 {
			return e.ID + ": " + t.Tool + " isn't in its description", nil
		}
		finals, err := FinalArguments(snap.Operations[i])
		if err != nil {
			return e.ID + ": " + t.Tool + ": " + err.Error(), nil
		}
		for _, f := range finals {
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

// newestRevisions is entries with only the newest revision of each name,
// for what's offered: a search or a list never offers an older one. A grant
// still names the revision it was approved with.
func newestRevisions(entries map[string]CatalogEntry) map[string]CatalogEntry {
	newest := map[string]string{}
	rev := func(id string) int {
		_, r, _ := strings.Cut(id, "@")
		n := 0
		_, _ = fmt.Sscanf(r, "%d", &n)
		return n
	}
	for id := range entries {
		name, _, _ := strings.Cut(id, "@")
		if have, ok := newest[name]; !ok || rev(id) > rev(have) {
			newest[name] = id
		}
	}
	out := make(map[string]CatalogEntry, len(newest))
	for _, id := range newest {
		out[id] = entries[id]
	}
	return out
}
