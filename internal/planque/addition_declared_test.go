package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/fakes"
)

// An addition declared in KIFF while it waits leaves nothing there once it
// ends unapproved: declined, past its day, or replaced by another before the
// owner approved it. The grant keeps what it had, on its first credential, and no
// Card can be signed over tools nobody approved.
func TestADeclaredAdditionThatEndsLeavesNothingInKIFF(t *testing.T) {
	for _, how := range []string{"declined", "past its day", "revised"} {
		t.Run(how, func(t *testing.T) {
			e := newEnv(t)
			id := e.connectedRefunds()
			acct := "/v1/accounts/" + tenant + "/grants/" + id
			refund, _ := e.kiff.Connection(tenant, "create_refund")
			list := map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}
			st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup")
			if st != 200 || out["declared"] != true {
				t.Fatalf("add: %d %v", st, out)
			}
			if _, ok := e.kiff.Connection(tenant, "list_charges"); !ok {
				t.Fatal("the addition isn't declared")
			}
			switch how {
			case "declined":
				if st, _ := e.do("DELETE", acct+"/tools", nil, nil, "k_other"); st == 200 {
					t.Fatal("another account's agent withdrew the addition")
				}
				if st, out := e.do("DELETE", acct+"/tools", nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected || out["declared"] == true {
					t.Fatalf("decline: %d %v", st, out)
				}
				if st, out := e.do("DELETE", acct+"/tools", nil, nil, "k_setup"); st != 409 {
					t.Fatalf("declining again, with nothing waiting: %d %v", st, out)
				}
			case "past its day":
				e.at(holdTTL + time.Minute)
				if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusConnected || out["declared"] == true {
					t.Fatalf("after a day: %d %v", st, out)
				}
			case "revised":
				read := map[string]any{"tool": "stripe_api_read", "title": "Read Stripe", "name": "stripe_api_read", "read_only": true}
				st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{read}}, nil, "k_setup")
				if st != 200 || out["status"] != StatusPending || out["declared"] != true {
					t.Fatalf("revise: %d %v", st, out)
				}
				if _, ok := e.kiff.Connection(tenant, "stripe_api_read"); !ok {
					t.Fatal("the revised addition isn't declared")
				}
			}
			if _, ok := e.kiff.Connection(tenant, "list_charges"); ok {
				t.Fatal("the ended addition is still in KIFF")
			}
			if c, ok := e.kiff.Connection(tenant, "create_refund"); !ok || c.Credential != refund.Credential {
				t.Fatal("the grant's own tool went, or changed, with the addition")
			}
			if d := e.refund(); d != "allowed" {
				t.Fatalf("a refund after the addition ended: %s", d)
			}
			g, _ := e.s.loadGrant(context.Background(), tenant, id)
			switch how {
			case "declined", "past its day":
				if len(g.Pending) != 0 || g.PendingTransport != "" || len(g.Tools) != 1 {
					t.Fatalf("the grant after it ended: pending %v, transport %q, tools %d", g.Pending, g.PendingTransport, len(g.Tools))
				}
			case "revised":
				if len(g.Pending) != 1 || g.Pending[0].Name != "stripe_api_read" {
					t.Fatalf("declared after the revision: %v", g.Pending)
				}
			}
		})
	}
}

// One approval, one fingerprint: the owner approves an addition and the
// Card drafted over it together (kiff.requests.v1, tools_and_cards), signed
// with the approval key beside KIFF's confirm of the Card. planque checks the
// pairing as on any request.
// Category: authority
// Attack: the pairing over a digest the addition moved past, from the device key, paired with itself, used twice
// Refused: each, with nothing approved; the approval key's over the current digest connects the declared tools
func TestAnAdditionApprovedTogetherWithItsCard(t *testing.T) {
	e := newEnv(t)
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
	id := e.connectedRefunds()
	acct := "/v1/accounts/" + tenant + "/grants/" + id
	approve := acct + "/approve"
	card := sha256.Sum256([]byte("KIFF's confirm of the Card's draft"))
	with := hex.EncodeToString(card[:])

	list := map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}
	_, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup")
	moved := out["digest"].(string)
	read := map[string]any{"tool": "stripe_api_read", "title": "Read Stripe", "name": "stripe_api_read", "read_only": true}
	st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{read}}, nil, "k_setup")
	if st != 200 || out["declared"] != true {
		t.Fatalf("revise: %d %v", st, out)
	}
	digest := out["digest"].(string)
	declared, _ := e.kiff.Connection(tenant, "stripe_api_read")
	unapproved := func(what string) {
		t.Helper()
		g, _ := e.s.loadGrant(context.Background(), tenant, id)
		if g.Status != StatusPending || len(g.Tools) != 1 || len(g.Pending) != 1 {
			t.Fatalf("%s: %s, tools %d, pending %d", what, g.Status, len(g.Tools), len(g.Pending))
		}
	}

	// Negative 2: the pairing over the digest the addition moved past.
	approval.With = with
	if st, out := e.do("POST", approve, map[string]string{"digest": moved}, approval, ""); st == http.StatusOK {
		t.Fatalf("the pairing over a moved digest: %d %v", st, out)
	}
	unapproved("after a moved digest")

	// The device key's pairing: refused.
	phone.With = with
	if st, out := e.do("POST", approve, map[string]string{"digest": digest}, phone, ""); st != http.StatusForbidden {
		t.Fatalf("the device key's pairing: %d %v", st, out)
	}
	phone.With = ""
	unapproved("after the device key's pairing")

	// Negative 3: paired with itself.
	approval.With, approval.WithItself = "", true
	if st, out := e.do("POST", approve, map[string]string{"digest": digest}, approval, ""); st != http.StatusUnauthorized {
		t.Fatalf("paired with itself: %d %v", st, out)
	}
	approval.WithItself = false
	unapproved("after a pairing with itself")

	// The approval key's pairing over the current digest: connected, the
	// declared registration kept, and the record holds the paired message.
	approval.With = with
	b, _ := json.Marshal(map[string]string{"digest": digest})
	send := func(header string) (int, string) {
		t.Helper()
		req, _ := http.NewRequest("POST", e.url+approve, bytes.NewReader(b))
		if header == "" {
			approval.Sign(req, tenant, b)
		} else {
			req.Header.Set(device.Header, header)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode, req.Header.Get(device.Header)
	}
	st, signed := send("")
	if st != http.StatusOK {
		t.Fatalf("the approval key's pairing: %d", st)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Status != StatusConnected || len(g.Tools) != 2 || len(g.Pending) != 0 {
		t.Fatalf("after the approval: %s, tools %d, pending %d", g.Status, len(g.Tools), len(g.Pending))
	}
	if g.Approval == nil || g.Approval.Key != "critical" || !bytes.Contains([]byte(g.Approval.Message), []byte(with)) {
		t.Fatalf("the record: %+v", g.Approval)
	}
	if c, ok := e.kiff.Connection(tenant, "stripe_api_read"); !ok || c.Credential != declared.Credential {
		t.Fatal("the declared tool was registered again, or went, at the approval")
	}

	// Negative 3: the same signed approval again.
	if st, _ := send(signed); st == http.StatusOK {
		t.Fatalf("the pairing used twice: %d", st)
	}
}

// A grant signed in with nothing connected yet has its first tools declared
// in KIFF once they're set, so the taf's Card can be drafted over them and
// approved with them: one approval for the tools and the Card. A call before
// the approval is refused and reaches nobody; approved, the same
// registration serves it. Declined, they leave KIFF and the grant is listed
// again, its sign-in held; revised, the earlier tools leave KIFF; not
// approved within a day, they leave with the sign-in.
func TestASignedInGrantsFirstToolsAreDeclared(t *testing.T) {
	for _, how := range []string{"approved", "declined", "revised", "past its day"} {
		t.Run(how, func(t *testing.T) {
			e := newEnv(t)
			id := e.listedGrant()
			acct := "/v1/accounts/" + tenant + "/grants/" + id
			list := map[string]any{"tool": "list_charges", "title": "List charges", "name": "list_charges", "read_only": true}
			st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{list}}, nil, "k_setup")
			if st != 200 || out["status"] != StatusPending || out["declared"] != true {
				t.Fatalf("set tools: %d %v", st, out)
			}
			declared, ok := e.kiff.Connection(tenant, "list_charges")
			if !ok {
				t.Fatal("list_charges isn't declared while it waits")
			}
			res, err := e.kiff.Call(context.Background(), "k_taf", "list_charges", map[string]any{})
			if err != nil || res.Decision != "allowed" {
				t.Fatalf("a call before the approval: %+v %v", res, err)
			}
			if text, isErr := resultText(t, res.Result); !isErr || !strings.Contains(text, "waits for the owner's approval") || !strings.Contains(text, "Nothing was sent") {
				t.Fatalf("a call before the approval went through: %s", text)
			}
			switch how {
			case "approved":
				signIns, _ := e.vendor.Live()
				st, out := e.do("POST", acct+"/approve", map[string]string{"digest": out["digest"].(string)}, e.phone, "")
				if st != 200 || out["status"] != StatusConnected || out["declared"] == true || out["sign_in_url"] != nil {
					t.Fatalf("approve: %d %v", st, out)
				}
				if a, _ := e.vendor.Live(); a != signIns {
					t.Fatalf("approving signed in again: %d live tokens, were %d", a, signIns)
				}
				if c, ok := e.kiff.Connection(tenant, "list_charges"); !ok || c.Credential != declared.Credential {
					t.Fatal("list_charges was registered again, or went, at the approval")
				}
				res, _ := e.kiff.Call(context.Background(), "k_taf", "list_charges", map[string]any{})
				if text, isErr := resultText(t, res.Result); isErr {
					t.Fatalf("a call after the approval: %s", text)
				}
			case "declined":
				signIns, _ := e.vendor.Live()
				st, out := e.do("DELETE", acct+"/tools", nil, nil, "k_setup")
				if st != 200 || out["status"] != StatusListed || out["declared"] == true || out["signed_in"] != true {
					t.Fatalf("decline: %d %v", st, out)
				}
				if c := e.kiff.Connections(tenant); len(c) != 0 {
					t.Fatalf("still in KIFF after the decline: %v", c)
				}
				if a, _ := e.vendor.Live(); a != signIns {
					t.Fatalf("the decline revoked the sign-in: %d live tokens, were %d", a, signIns)
				}
				// Other tools on the same sign-in, declared again.
				read := map[string]any{"tool": "stripe_api_read", "title": "Read Stripe", "name": "stripe_api_read", "read_only": true}
				if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{read}}, nil, "k_setup"); st != 200 || out["declared"] != true {
					t.Fatalf("tools after the decline: %d %v", st, out)
				}
			case "revised":
				read := map[string]any{"tool": "stripe_api_read", "title": "Read Stripe", "name": "stripe_api_read", "read_only": true}
				if st, out := e.do("POST", acct+"/tools", map[string]any{"tools": []any{read}}, nil, "k_setup"); st != 200 || out["declared"] != true {
					t.Fatalf("revise: %d %v", st, out)
				}
				if _, ok := e.kiff.Connection(tenant, "list_charges"); ok {
					t.Fatal("the revised-away tool is still in KIFF")
				}
				if _, ok := e.kiff.Connection(tenant, "stripe_api_read"); !ok {
					t.Fatal("the revised tool isn't declared")
				}
			case "past its day":
				e.at(holdTTL + time.Minute)
				if st, out := e.do("GET", acct, nil, nil, "k_setup"); st != 200 || out["status"] != StatusFailed || out["declared"] == true {
					t.Fatalf("after a day: %d %v", st, out)
				}
				if c := e.kiff.Connections(tenant); len(c) != 0 {
					t.Fatalf("still in KIFF after a day: %v", c)
				}
			}
		})
	}
}
