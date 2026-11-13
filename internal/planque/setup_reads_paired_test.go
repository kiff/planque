package planque

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/words"
)

// A key given together with its Card (key_and_cards: one fingerprint over
// the key's request and KIFF's confirm) still gets Setup's read Card: the
// approval KIFF checks names the key's request by its hash beside the
// Card's, so planque sends that request's own message with it. Without it,
// KIFF can't read what was approved, issues no Card, and Setup's first read
// is refused while planque counts it connected (Eval 43519074848).
func TestSetupReadsOfAKeyGivenTogetherWithItsCard(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	phone := e.device("dev_phone")
	approval := fakes.NewPhone()
	approval.DeviceID = "dev_phone"
	d := e.kiff.Devices["dev_phone"]
	d.Approval = approval
	e.kiff.Devices["dev_phone"] = d
	if st, out := e.enrolAs(phone); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	e.enrolled, e.phone = true, phone
	tools := []map[string]any{apiTools[0], {"tool": "listCharges", "title": "List charges", "name": "list_charges", "read_only": true, "setup_reads": true}}
	st, out := e.askAPIGrant(v, tools, nil)
	if st != 201 || out["declared"] != true {
		t.Fatalf("ask: %d %v", st, out)
	}
	id, digest := out["id"].(string), out["digest"].(string)
	card := sha256.Sum256([]byte("KIFF's confirm of the Card's draft"))
	approval.With = hex.EncodeToString(card[:])
	e.phone = approval
	if st, out := e.giveKey(id, digest, v.key); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("key with its Card: %d %v", st, out)
	}
	if n := len(e.kiff.SetupReads); n != 1 {
		g, _ := e.s.loadGrant(context.Background(), tenant, id)
		t.Fatalf("Setup's Card asked %d times; the grant's events: %v", n, g.Events)
	}
	a := e.kiff.SetupReads[0].Approval
	if !strings.HasPrefix(a.Message, "kiff.requests.v1\n") || !strings.HasPrefix(a.Request, "planque.request.") {
		t.Fatalf("the approval sent: message %q, request %q", a.Message, a.Request)
	}
	sum := sha256.Sum256([]byte(a.Request))
	if !strings.Contains(a.Message, hex.EncodeToString(sum[:])) || !strings.Contains(a.Request, "/grants/"+id+"/key") {
		t.Fatalf("the request isn't the one the paired message names:\n%s\n%s", a.Message, a.Request)
	}
}

// A request signed alone sends no request beside its message: the message
// is the request, as before.
func TestSetupReadsOfAKeyGivenAloneSendNoRequest(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	tools := []map[string]any{apiTools[0], {"tool": "listCharges", "title": "List charges", "name": "list_charges", "read_only": true, "setup_reads": true}}
	st, out := e.askAPIGrant(v, tools, nil)
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	if st, out := e.giveKey(out["id"].(string), out["digest"].(string), v.key); st != 200 {
		t.Fatalf("key: %d %v", st, out)
	}
	if len(e.kiff.SetupReads) != 1 || e.kiff.SetupReads[0].Approval.Request != "" || !strings.HasPrefix(e.kiff.SetupReads[0].Approval.Message, "planque.request.") {
		t.Fatalf("Setup's Card asked: %+v", e.kiff.SetupReads)
	}
}

// A key app's first request names its setup reads, as a catalog entry's
// does: given the key, they connect at once and Setup's Card is asked, as
// after a sign-in. Before, only a sign-in connected them, and a key left
// the grant listed with Setup's reads never in KIFF.
func TestSetupReadsNamedAtAKeyRequestConnectWithTheKey(t *testing.T) {
	e := newEnv(t)
	v := newRESTVendor(t)
	e.enrol()
	st, out := e.askAPIGrant(v, []map[string]any{}, nil)
	if st != 201 {
		t.Fatalf("ask: %d %v", st, out)
	}
	ctx := context.Background()
	g, _ := e.s.loadGrant(ctx, tenant, out["id"].(string))
	g.Request.SetupReads = []string{"listCharges"}
	g.Text, g.Texts = texts(func(w words.Approval) string { return requestText(w, g.Request, false, "") })
	var err error
	if g.Digest, err = digest(g.Request); err != nil {
		t.Fatal(err)
	}
	if err := e.s.saveGrant(ctx, g); err != nil {
		t.Fatal(err)
	}
	st, out = e.giveKey(g.ID, g.Digest, v.key)
	if st != 200 || out["status"] != StatusConnected {
		got, _ := e.s.loadGrant(ctx, tenant, g.ID)
		t.Fatalf("key: %d %v; events %v", st, out, got.Events)
	}
	g, _ = e.s.loadGrant(ctx, tenant, g.ID)
	if len(g.Tools) != 1 || !g.Tools[0].SetupReads || !g.Tools[0].ReadOnly {
		t.Fatalf("the grant's tools: %+v", g.Tools)
	}
	if len(e.kiff.SetupReads) != 1 || e.kiff.SetupReads[0].Grant != g.ID {
		t.Fatalf("Setup's Card wasn't asked of KIFF: %+v", e.kiff.SetupReads)
	}
}
