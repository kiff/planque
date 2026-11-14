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

// CleanStoredAnswers is the one-time clean-up of the answers kept before
// answerTTL, which the owner approved: a record's answer older than a day is
// dropped as the sweep would drop it, and a younger one gets its expiry and
// its entry. Every record keeps its outcome, its arguments' hash and its
// state, all the replay guard reads. It runs from the sweep until it has
// finished once, then never again; a record it couldn't save (a write in
// between) leaves the marker unwritten, so the next sweep runs it again. It
// logs counts only.
func (s *Server) CleanStoredAnswers(ctx context.Context) (dropped, indexed int, err error) {
	if _, ok, err := s.KV.Get(ctx, onceAnswers.PK, onceAnswers.SK); err != nil || ok {
		return 0, 0, err
	}
	items, err := s.KV.Scan(ctx, "OP#")
	if err != nil {
		return 0, 0, err
	}
	skipped := 0
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
				return dropped, indexed, err
			}
		} else {
			dropAnswer(&op)
		}
		switch err := putDoc(ctx, s.KV, next, op, store.Cond{Ver: &ver}); {
		case errors.Is(err, store.ErrCondition):
			skipped++ // written in between: the next run reads it again
		case err != nil:
			return dropped, indexed, err
		case young:
			indexed++
		default:
			dropped++
		}
	}
	if skipped > 0 {
		s.log().InfoContext(ctx, "stored answers: some records changed during the clean-up; it runs again", slog.Int("skipped", skipped))
		return dropped, indexed, nil
	}
	if err := putDoc(ctx, s.KV, onceAnswers, map[string]string{"at": s.now().UTC().Format(time.RFC3339)}, store.Cond{NotExists: true}); err != nil && !errors.Is(err, store.ErrCondition) {
		return dropped, indexed, err
	}
	s.log().InfoContext(ctx, "stored answers cleaned once", slog.Int("dropped", dropped), slog.Int("indexed", indexed))
	return dropped, indexed, nil
}
