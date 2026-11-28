package planque

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

// A vendor's kept answer is sealed in its record, under the account and the
// operation: the record never holds it readable (with full scopes it can
// carry IBANs, balances and client records), and a replay opens it whole.
func TestAKeptAnswerIsSealedAndAReplayOpensItWhole(t *testing.T) {
	e, f := newConnectorFixture(t)
	g := connectConnector(t, e, f)
	ctx := context.Background()
	args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_sealed"}`)
	permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_sealed", Args: args})
	first := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
	if !strings.Contains(first, "reply_demo") {
		t.Fatalf("call: %s", first)
	}
	it, _, err := e.kv.Get(ctx, tenantPK(tenant), skOp(g.ID, "op_sealed"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(it.Doc), "reply_demo") {
		t.Fatalf("the answer is readable in its record: %s", it.Doc)
	}
	rec := e.opRecord(g, "op_sealed")
	if len(rec.Result) != 0 || len(rec.SealedResult) == 0 || rec.ResultUntil.IsZero() {
		t.Fatalf("not sealed: result %q, sealed %d bytes", rec.Result, len(rec.SealedResult))
	}
	if again := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); again != first || f.executions != 1 {
		t.Fatalf("replay: %s (sends %d)", again, f.executions)
	}
	// It opens for its own record only: under another operation's binding it doesn't.
	if _, err := e.s.Sealer.Open(ctx, rec.SealedResult, answerBinding(store.Item{PK: tenantPK(tenant), SK: skOp(g.ID, "op_other")})); err == nil {
		t.Fatal("the answer opened under another operation")
	}
	// Past its day, the sweep drops it to the outcome, seal and all.
	e.at(25 * time.Hour)
	if n, err := e.s.SweepAnswers(ctx); n != 1 || err != nil {
		t.Fatalf("sweep: %d %v", n, err)
	}
	rec = e.opRecord(g, "op_sealed")
	if rec.SealedResult != nil || !strings.Contains(string(rec.Result), "Answered once") || !rec.AnswerDropped {
		t.Fatalf("after the sweep: %+v", rec)
	}
}

// failSealer refuses to seal or open answers, as KMS throttled or down
// would; anything else goes through.
type failSealer struct {
	seal.Sealer
	seal, open bool
}

func (f failSealer) Seal(ctx context.Context, p []byte, b map[string]string) ([]byte, error) {
	if f.seal && b["kind"] == "answer" {
		return nil, errors.New("ThrottlingException")
	}
	return f.Sealer.Seal(ctx, p, b)
}

func (f failSealer) Open(ctx context.Context, s []byte, b map[string]string) ([]byte, error) {
	if f.open && b["kind"] == "answer" {
		return nil, errors.New("ThrottlingException: secret detail")
	}
	return f.Sealer.Open(ctx, s, b)
}

// A seal that fails keeps the outcome alone, never the answer readable, and
// the call's own answer is unchanged. An answer that doesn't open on a
// replay is a plain miss: the outcome, never an error that says why, and
// the call is still never sent again.
func TestAnAnswerThatCantBeSealedOrOpenedIsAPlainMiss(t *testing.T) {
	t.Run("seal fails", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		e.s.Sealer = failSealer{Sealer: e.s.Sealer, seal: true}
		ctx := context.Background()
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_noseal"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_noseal", Args: args})
		if res := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); !strings.Contains(res, "reply_demo") {
			t.Fatalf("the call's own answer: %s", res)
		}
		it, _, _ := e.kv.Get(ctx, tenantPK(tenant), skOp(g.ID, "op_noseal"))
		if strings.Contains(string(it.Doc), "reply_demo") {
			t.Fatalf("kept readable when the seal failed: %s", it.Doc)
		}
		if again := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit)); !strings.Contains(again, "Answered once") || f.executions != 1 {
			t.Fatalf("replay: %s (sends %d)", again, f.executions)
		}
	})
	t.Run("open fails", func(t *testing.T) {
		e, f := newConnectorFixture(t)
		g := connectConnector(t, e, f)
		ctx := context.Background()
		args := json.RawMessage(`{"ticket":"ticket_demo","text":"hello","planque_op":"op_noopen"}`)
		permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "reply_ticket", Sub: "agt_taf", Op: "op_noopen", Args: args})
		_ = e.s.call(ctx, g, "reply_ticket", args, permit)
		e.s.Sealer = failSealer{Sealer: e.s.Sealer, open: true}
		again := mustJSONString(e.s.call(ctx, g, "reply_ticket", args, permit))
		if !strings.Contains(again, "Answered once") || strings.Contains(again, "Throttling") || strings.Contains(again, "secret detail") || f.executions != 1 {
			t.Fatalf("replay when the answer doesn't open: %s (sends %d)", again, f.executions)
		}
	})
}
