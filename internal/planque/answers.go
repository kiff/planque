package planque

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kiff/planque/internal/store"
)

// answerTTL is how long a vendor's answer stays in its operation's record.
// The outcome (done, not sent, unknown) stays for as long as the record, so
// an operation still never runs twice; the answer itself, the vendor's data
// (an email, a customer's record), is only kept for a replay soon after.
const answerTTL = 24 * time.Hour

// pkAnswerDue indexes each kept answer by when it expires, as pkDue does for
// grants, so the scheduled sweep finds them without reading every account.
const pkAnswerDue = "ANSWERDUE"

// answerGone stands in for an answer past answerTTL: the outcome, never
// the vendor's data.
func answerGone(isError bool) map[string]any {
	if isError {
		return toolResult("Answered once, with an error; planque no longer keeps the answer (it's kept 24 hours). Nothing new was sent; check the account before making a new call.", true)
	}
	return toolResult("Answered once; planque no longer keeps the answer (it's kept 24 hours). Nothing new was sent; read again with a new call.", false)
}

// keepAnswer sets when op's answer expires and indexes it, before the record
// is saved with it, so an answer is never stored where the sweep can't find
// it. Only a vendor's answer (OpDone) expires: a refusal or a lost answer is
// planque's own words.
func (s *Server) keepAnswer(ctx context.Context, claim store.Item, op *operation) error {
	if op.State != OpDone || len(op.Result) == 0 {
		return nil
	}
	op.ResultUntil = s.now().UTC().Add(answerTTL)
	tenant := strings.TrimPrefix(claim.PK, "TENANT#")
	it := store.Item{PK: pkAnswerDue, SK: dueSK(op.ResultUntil, tenant, claim.SK), Doc: []byte("{}"), TTL: op.ResultUntil.Add(7 * 24 * time.Hour).Unix()}
	if err := s.KV.Put(ctx, it, store.Cond{}); err != nil {
		return fmt.Errorf("index the answer's expiry: %w", err)
	}
	return nil
}

// answerExpired says op's answer is past its time: a replay gets the
// outcome alone, even before the sweep drops the answer.
func (s *Server) answerExpired(op operation) bool {
	return !op.ResultUntil.IsZero() && !s.now().Before(op.ResultUntil)
}

// dropAnswer replaces op's kept answer with the outcome alone.
func dropAnswer(op *operation) {
	var was struct {
		IsError bool `json:"isError"`
	}
	_ = json.Unmarshal(op.Result, &was)
	op.Result, _ = json.Marshal(answerGone(was.IsError))
	op.ResultUntil, op.AnswerDropped = time.Time{}, true
}

// SweepAnswers drops every kept answer past answerTTL from its record. It
// returns how many it dropped, and the first error. An entry goes only once
// its record no longer holds the answer; any failure keeps it for the next
// sweep.
func (s *Server) SweepAnswers(ctx context.Context) (int, error) {
	now, dropped := s.now(), 0
	items, err := s.KV.QueryTo(ctx, pkAnswerDue, now.UTC().Truncate(time.Second).Format(time.RFC3339)+"#~#~")
	if err != nil {
		return 0, err
	}
	var first error
	for _, it := range items {
		at, tenant, sk, ok := parseDue(it.SK)
		if ok && at.After(now) {
			break
		}
		if ok {
			n, err := s.dropDue(ctx, tenant, sk)
			if err != nil {
				s.log().WarnContext(ctx, "sweep: an answer waits for the next sweep", slog.String("op", sk), slog.String("err", err.Error()))
				first = cmpErr(first, err)
				continue
			}
			dropped += n
		}
		if _, _, err := s.KV.Delete(ctx, pkAnswerDue, it.SK); err != nil {
			first = cmpErr(first, err)
		}
	}
	return dropped, first
}

// dropDue drops the answer in the record at sk if it's past its time.
func (s *Server) dropDue(ctx context.Context, tenant, sk string) (int, error) {
	var op operation
	it, err := getDoc(ctx, s.KV, tenantPK(tenant), sk, &op)
	if errors.Is(err, errNotFound) {
		return 0, nil // the account was erased
	}
	if err != nil {
		return 0, err
	}
	if op.ResultUntil.IsZero() {
		return 0, nil // dropped already
	}
	if !s.answerExpired(op) {
		return 0, errors.New("the answer is not due yet")
	}
	dropAnswer(&op)
	next, ver := it, it.Ver
	next.Ver = ver + 1
	if err := putDoc(ctx, s.KV, next, op, store.Cond{Ver: &ver}); err != nil {
		return 0, err
	}
	return 1, nil
}

// onceAnswers marks the one-time clean-up of the answers kept before
// answerTTL existed as done.
var onceAnswers = store.Item{PK: "ONCE", SK: "answers-24h"}

// cleanProgress is where the one-time clean-up stands between sweeps: the
// next page's cursor, and whether this pass skipped a record (a write in
// between), so the pass is run again rather than marked done.
var cleanProgress = store.Item{PK: "ONCE", SK: "answers-24h-progress"}

type progress struct {
	Cursor  string `json:"cursor,omitempty"`
	Skipped bool   `json:"skipped,omitempty"`
}

// cleanPage and cleanBudget bound one sweep's share of the clean-up: pages
// of about cleanPage records, read for at most cleanBudget, so a sweep ends
// well inside the function's time limit and the next one carries on from
// the saved cursor.
var (
	cleanPage   = 200
	cleanBudget = 15 * time.Second
)

// CleanStoredAnswers is the one-time clean-up of the answers kept before
// answerTTL, which the owner approved: a record's answer older than a day is
// dropped as the sweep would drop it, and a younger one gets its expiry and
// its entry. Every record keeps its outcome, its arguments' hash and its
// state, all the replay guard reads. It reads the table page by page within
// cleanBudget and saves its cursor, so each sweep carries on where the last
// stopped. A pass that skipped a record (a write in between) starts again;
// one that skipped none marks it done, and it never runs again. It logs
// counts only.
func (s *Server) CleanStoredAnswers(ctx context.Context) (dropped, indexed int, err error) {
	if _, ok, err := s.KV.Get(ctx, onceAnswers.PK, onceAnswers.SK); err != nil || ok {
		return 0, 0, called("read whether it's done", err)
	}
	var at progress
	if _, err := getDoc(ctx, s.KV, cleanProgress.PK, cleanProgress.SK, &at); err != nil && !errors.Is(err, errNotFound) {
		return 0, 0, called("read its cursor", err)
	}
	begun := time.Now()
	for {
		items, next, err := s.KV.ScanPage(ctx, "OP#", at.Cursor, cleanPage)
		if err != nil {
			return dropped, indexed, called("scan a page (Scan)", err)
		}
		d, i, skipped, err := s.cleanItems(ctx, items)
		dropped, indexed = dropped+d, indexed+i
		if err != nil {
			return dropped, indexed, err
		}
		at.Cursor, at.Skipped = next, at.Skipped || skipped > 0
		if next == "" {
			break
		}
		if err := putDoc(ctx, s.KV, cleanProgress, at, store.Cond{}); err != nil {
			return dropped, indexed, called("save its cursor (PutItem)", err)
		}
		if time.Since(begun) > cleanBudget {
			s.log().InfoContext(ctx, "stored answers: the clean-up carries on at the next sweep", slog.Int("dropped", dropped), slog.Int("indexed", indexed))
			return dropped, indexed, nil
		}
	}
	if at.Skipped {
		// A pass that skipped a record runs again from the start.
		if err := putDoc(ctx, s.KV, cleanProgress, progress{}, store.Cond{}); err != nil {
			return dropped, indexed, called("start its pass again (PutItem)", err)
		}
		s.log().InfoContext(ctx, "stored answers: some records changed during the clean-up; it runs again", slog.Int("dropped", dropped), slog.Int("indexed", indexed))
		return dropped, indexed, nil
	}
	if err := putDoc(ctx, s.KV, onceAnswers, map[string]string{"at": s.now().UTC().Format(time.RFC3339)}, store.Cond{NotExists: true}); err != nil && !errors.Is(err, store.ErrCondition) {
		return dropped, indexed, called("mark it done (PutItem)", err)
	}
	_, _, _ = s.KV.Delete(ctx, cleanProgress.PK, cleanProgress.SK)
	s.log().InfoContext(ctx, "stored answers cleaned once", slog.Int("dropped", dropped), slog.Int("indexed", indexed))
	return dropped, indexed, nil
}

// cleanItems drops or indexes the kept answers among items.
func (s *Server) cleanItems(ctx context.Context, items []store.Item) (dropped, indexed, skipped int, err error) {
	for _, it := range items {
		if !strings.HasPrefix(it.PK, "TENANT#") {
			continue
		}
		var op operation
		if json.Unmarshal(it.Doc, &op) != nil || op.State != OpDone || len(op.Result) == 0 || !op.ResultUntil.IsZero() || op.AnswerDropped {
			continue
		}
		next, ver := it, it.Ver
		next.Ver = ver + 1
		young := false
		if until := op.StartedAt.Add(answerTTL); s.now().Before(until) {
			young = true
			op.ResultUntil = until.UTC()
			due := store.Item{PK: pkAnswerDue, SK: dueSK(until, strings.TrimPrefix(it.PK, "TENANT#"), it.SK), Doc: []byte("{}"), TTL: until.Add(7 * 24 * time.Hour).Unix()}
			if err := s.KV.Put(ctx, due, store.Cond{}); err != nil {
				return dropped, indexed, skipped, called("index an answer's expiry (PutItem)", err)
			}
		} else {
			dropAnswer(&op)
		}
		switch err := putDoc(ctx, s.KV, next, op, store.Cond{Ver: &ver}); {
		case errors.Is(err, store.ErrCondition):
			skipped++ // written in between: the pass runs again
		case err != nil:
			return dropped, indexed, skipped, called("save a record (PutItem)", err)
		case young:
			indexed++
		default:
			dropped++
		}
	}
	return dropped, indexed, skipped, nil
}

// called is err with the clean-up's step that met it, so the sweep's warning
// says which call stopped it (14 Nov: a throttle logged with no call named).
func called(step string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("the clean-up couldn't %s: %w", step, err)
}
