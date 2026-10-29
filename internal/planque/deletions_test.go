package planque

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/fakes"
)

// An intent for the account, signed by phone, made at issued and expiring
// after life, with nonce.
func intent(phone *fakes.Phone, device string, issued time.Time, life time.Duration, nonce string) map[string]string {
	text := strings.Join([]string{"tafs.delete-account.v2", "tenant=" + tenant, "device=" + device,
		"issued=" + issued.UTC().Format(time.RFC3339), "expires=" + issued.Add(life).UTC().Format(time.RFC3339), "nonce=" + nonce}, "\n")
	return map[string]string{"intent": text, "signature": phone.SignBytes([]byte(text))}
}

// Deleting the account from one signed intent, which tafs cloud carries:
// accepted (deleting nothing) only from the account's device of a person who
// may delete it, once per nonce; finished only for the intent accepted, with
// no deadline, and the same when asked again; status says which, deleting
// nothing.
func TestDeletingTheAccountFromOneSignedIntent(t *testing.T) {
	e := newEnv(t)
	e.connect()
	ctx := context.Background()
	tools := len(e.kiff.Connections(tenant))
	now := time.Now()
	good := intent(e.phone, "dev_phone", now.Add(-time.Minute), 10*time.Minute, "0123456789abcdef0123456789abcdef")

	// KIFF accepts first, from the phone's session, checking the approval key.
	accept := func(b map[string]string) {
		if e.kiff.Deletions == nil {
			e.kiff.Deletions = map[string]string{}
		}
		e.kiff.Deletions[b["intent"]] = b["signature"]
	}

	// Refused at accept, deleting nothing: expired, too long-lived, malformed,
	// and any intent KIFF didn't accept (another key's signature, a person who
	// can't delete the account: KIFF refused those).
	other := e.device("dev_other")
	for _, c := range []struct {
		what string
		body map[string]string
		want int
	}{
		{"expired", intent(e.phone, "dev_phone", now.Add(-20*time.Minute), 10*time.Minute, "11111111111111111111111111111111"), 400},
		{"too long-lived", intent(e.phone, "dev_phone", now, 30*time.Minute, "22222222222222222222222222222222"), 400},
		{"malformed", map[string]string{"intent": "delete everything", "signature": "x"}, 400},
		{"not accepted by KIFF", good, 403},
		{"signed by another key", map[string]string{"intent": good["intent"], "signature": other.SignBytes([]byte(good["intent"]))}, 403},
	} {
		if st, out := e.do("POST", "/v1/account-deletions", c.body, nil, ""); st != c.want {
			t.Errorf("%s: %d %v, want %d", c.what, st, out, c.want)
		}
	}
	if n := len(e.kiff.Connections(tenant)); n != tools {
		t.Fatalf("a refused accept changed KIFF's tools: %d, was %d", n, tools)
	}
	accept(good)

	// Not accepted by planque yet: finish and status refuse, and nothing goes.
	if st, out := e.do("POST", "/v1/account-deletions/finish", good, nil, ""); st != 403 || out["error"] != "not_accepted" {
		t.Fatalf("finish before accept: %d %v", st, out)
	}
	if st, _ := e.do("POST", "/v1/account-deletions/status", good, nil, ""); st != 403 {
		t.Fatalf("status before accept: %d", st)
	}

	// Accepted: nothing deleted; the same intent again answers the same.
	if st, out := e.do("POST", "/v1/account-deletions", good, nil, ""); st != 202 {
		t.Fatalf("accept: %d %v", st, out)
	}
	if st, _ := e.do("POST", "/v1/account-deletions", good, nil, ""); st != 202 {
		t.Fatalf("the same intent again: %d", st)
	}
	if n := len(e.kiff.Connections(tenant)); n != tools {
		t.Fatalf("accept deleted tools: %d left of %d", n, tools)
	}
	if st, out := e.do("POST", "/v1/account-deletions/status", good, nil, ""); st != 200 || out["status"] != "accepted" {
		t.Fatalf("status after accept: %d %v", st, out)
	}
	// A different intent with a spent nonce is refused.
	spent := intent(e.phone, "dev_phone", now, 9*time.Minute, "0123456789abcdef0123456789abcdef")
	accept(spent)
	if st, _ := e.do("POST", "/v1/account-deletions", spent, nil, ""); st != 409 {
		t.Fatalf("a spent nonce: %d, want 409", st)
	}
	// A later intent with a new nonce takes the earlier one's place (KIFF may
	// have refused that one): the earlier one no longer finishes.
	later := intent(e.phone, "dev_phone", now, 10*time.Minute, "fedcba9876543210fedcba9876543210")
	accept(later)
	if st, _ := e.do("POST", "/v1/account-deletions", later, nil, ""); st != 202 {
		t.Fatalf("a later intent: %d", st)
	}
	if st, _ := e.do("POST", "/v1/account-deletions/finish", good, nil, ""); st != 403 {
		t.Fatalf("finishing with the replaced intent: %d, want 403", st)
	}

	// planque never erases on an intent KIFF doesn't vouch for.
	delete(e.kiff.Deletions, later["intent"])
	if st, out := e.do("POST", "/v1/account-deletions/finish", later, nil, ""); st != 403 || out["error"] != "not_accepted" {
		t.Fatalf("finish KIFF doesn't vouch for: %d %v", st, out)
	}
	accept(later)

	// Finished with the accepted one, past its expiry: everything goes.
	e.s.Now = func() time.Time { return now.Add(3 * time.Hour) }
	st, out := e.do("POST", "/v1/account-deletions/finish", later, nil, "")
	if st != 200 || out["erased"] != true {
		t.Fatalf("finish: %d %v", st, out)
	}
	if n := len(e.kiff.Connections(tenant)); n != 0 {
		t.Fatalf("KIFF still lists %d tools", n)
	}
	if items, _ := e.kv.Query(ctx, tenantPK(tenant), ""); len(items) != 0 {
		t.Fatalf("%d records kept", len(items))
	}
	// Asked again, it's done; status says erased.
	if st, out := e.do("POST", "/v1/account-deletions/finish", later, nil, ""); st != 200 || out["erased"] != true {
		t.Fatalf("finish again: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/account-deletions/status", later, nil, ""); st != 200 || out["status"] != "erased" {
		t.Fatalf("status after: %d %v", st, out)
	}
}

// An account planque keeps nothing for accepts a deletion with nothing to do.
func TestADeletionWithNothingKept(t *testing.T) {
	e := newEnv(t)
	body := intent(e.phone, "dev_phone", time.Now(), 10*time.Minute, "0123456789abcdef0123456789abcdef")
	if st, out := e.do("POST", "/v1/account-deletions", body, nil, ""); st != 202 || out["nothing_kept"] != true {
		t.Fatalf("accept with nothing kept: %d %v", st, out)
	}
	if st, out := e.do("POST", "/v1/account-deletions/finish", body, nil, ""); st != 200 || out["erased"] != true {
		t.Fatalf("finish with nothing kept: %d %v", st, out)
	}
}
