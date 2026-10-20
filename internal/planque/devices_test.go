package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/store"
)

// device adds a device of tenant's owner to KIFF and returns it, signing
// with its KIFF id (kiff/planque#10).
func (e *env) device(id string) *fakes.Phone {
	p := fakes.NewPhone()
	p.DeviceID = id
	e.kiff.Devices[id] = fakes.Device{Tenant: tenant, Subject: "owner", Platform: "ios", Name: "tafs on " + strings.TrimPrefix(id, "dev_"), Phone: p, Member: true,
		Capabilities: []string{"holds.answer", "tools.connect"}}
	return p
}

// enrolAs hands planque the connector key, signed by a device KIFF knows.
func (e *env) enrolAs(p *fakes.Phone) (int, map[string]any) {
	e.t.Helper()
	return e.do("POST", "/v1/accounts/"+tenant+"/enrolment", map[string]string{"connector_key": "k_connector"}, p, "")
}

// The owner's iPhone and iPad, both KIFF devices: either sets the account
// up, either approves, and only the one that approved confirms.
func TestKIFFDevicesConnectFromEither(t *testing.T) {
	e := newEnv(t)
	iphone, ipad := e.device("dev_iphone"), e.device("dev_ipad")
	if st, out := e.enrolAs(iphone); st != 201 || out["device"] != "dev_iphone" {
		t.Fatalf("enrol from the iPhone: %d %v", st, out)
	}
	if st, out := e.enrolAs(ipad); st != 201 {
		t.Fatalf("enrol from the iPad, a second device: %d %v, want 201", st, out)
	}
	if it, ok, err := e.kv.Get(context.Background(), tenantPK(tenant), skEnrolment); err != nil || !ok || bytes.Contains(it.Doc, []byte("device_key")) {
		t.Fatalf("a KIFF device left a key in the enrolment: %s %v", it.Doc, err)
	}

	id, digest := e.request()
	approve := "/v1/accounts/" + tenant + "/grants/" + id + "/approve"
	st, out := e.do("POST", approve, map[string]string{"digest": digest}, ipad, "")
	if st != 200 {
		t.Fatalf("approve on the iPad: %d %v", st, out)
	}
	code := e.signIn(out["sign_in_url"].(string))
	confirm := "/v1/accounts/" + tenant + "/grants/" + id + "/confirm"
	if st, _ := e.do("POST", confirm, map[string]string{"confirm": code}, iphone, ""); st != 403 {
		t.Fatalf("the iPhone confirming the iPad's approval: %d, want 403", st)
	}
	if st, out := e.do("POST", confirm, map[string]string{"confirm": code}, ipad, ""); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("confirm on the iPad: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Approval == nil || g.Approval.DeviceID != "dev_ipad" || g.Approval.Subject != "owner" || len(g.Approval.DeviceKey) != 0 {
		t.Fatalf("the approval doesn't name the iPad: %+v", g.Approval)
	}
	if g.Approval.DeviceName != "tafs on ipad" {
		t.Fatalf("the approval keeps no device name: %+v", g.Approval)
	}
	// The record says who approved, and on which device.
	var said []string
	for _, ev := range g.Events {
		said = append(said, ev.What)
	}
	// No id or fingerprint in what the owner reads (#15): those are the approval's fields.
	for _, w := range said {
		if strings.Contains(w, "dev_") || strings.Contains(w, "owner") {
			t.Errorf("an event names an id: %q", w)
		}
	}
	for _, want := range []string{"approved on tafs on ipad", "sign-in confirmed on tafs on ipad"} {
		if !strings.Contains(strings.Join(said, "\n"), want) {
			t.Errorf("events lack %q: %v", want, said)
		}
	}
}

// KIFF decides which devices sign: one it doesn't know, a removed one,
// another account's, a forged signature, a person who left or who can't
// connect tools are all refused, and nothing is approved.
func TestKIFFRefusesTheDevice(t *testing.T) {
	e := newEnv(t)
	e.enrolAs(e.device("dev_iphone"))
	id, digest := e.request()
	approve := "/v1/accounts/" + tenant + "/grants/" + id + "/approve"

	forger := fakes.NewPhone()
	forger.DeviceID = "dev_iphone" // claims the iPhone, signs with its own key
	unknown := fakes.NewPhone()
	unknown.DeviceID = "dev_nobody"
	removed, other, left, viewer := e.device("dev_removed"), e.device("dev_other"), e.device("dev_left"), e.device("dev_viewer")
	e.kiff.Devices["dev_removed"] = withDevice(e.kiff.Devices["dev_removed"], func(d *fakes.Device) { d.Revoked = true })
	e.kiff.Devices["dev_other"] = withDevice(e.kiff.Devices["dev_other"], func(d *fakes.Device) { d.Tenant = "tnt_other" })
	e.kiff.Devices["dev_left"] = withDevice(e.kiff.Devices["dev_left"], func(d *fakes.Device) { d.Member = false })
	e.kiff.Devices["dev_viewer"] = withDevice(e.kiff.Devices["dev_viewer"], func(d *fakes.Device) { d.Capabilities = []string{"holds.answer"} })

	for _, c := range []struct {
		name  string
		phone *fakes.Phone
		want  int
	}{
		{"forged signature", forger, 401},
		{"unknown device", unknown, 401},
		{"removed device", removed, 401},
		{"another account's device", other, 401},
		{"a person who left", left, 403},
		{"a person who can't connect tools", viewer, 403},
	} {
		if st, out := e.do("POST", approve, map[string]string{"digest": digest}, c.phone, ""); st != c.want {
			t.Errorf("%s: %d %v, want %d", c.name, st, out, c.want)
		}
	}
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Status != StatusPending || g.Approval != nil {
		t.Fatalf("after refused devices the grant is %s", g.Status)
	}
}

func withDevice(d fakes.Device, f func(*fakes.Device)) fakes.Device {
	f(&d)
	return d
}

// planque keeps the time window and the replay check for a KIFF device:
// KIFF says only that the device signed the bytes.
func TestAKIFFDevicesSignatureIsFreshAndOnce(t *testing.T) {
	e := newEnv(t)
	ipad := e.device("dev_ipad")
	e.enrolAs(ipad)
	id, digest := e.request()
	approve := e.url + "/v1/accounts/" + tenant + "/grants/" + id + "/approve"
	b, _ := json.Marshal(map[string]string{"digest": digest})

	// Signed an hour ago, with a valid signature: refused before KIFF is asked.
	old := "2020-01-01T00:00:00Z"
	msg := device.Message(e.url, tenant, "POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", b, old)
	req, _ := http.NewRequest("POST", approve, bytes.NewReader(b))
	req.Header.Set(device.Header, "t="+old+",v1="+ipad.SignBytes(msg)+",d=dev_ipad")
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 401 {
		t.Fatalf("a stale signature: %d, want 401", resp.StatusCode)
	}

	req, _ = http.NewRequest("POST", approve, bytes.NewReader(b))
	ipad.Sign(req, tenant, b)
	sig := req.Header.Get(device.Header)
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 200 {
		t.Fatalf("approve: %d", resp.StatusCode)
	}
	req, _ = http.NewRequest("POST", approve, bytes.NewReader(b))
	req.Header.Set(device.Header, sig)
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != 409 {
		t.Fatalf("replayed: %d, want 409", resp.StatusCode)
	}
}

// The rollout: the phone enrolled with its key keeps signing as before
// while another device signs as a KIFF device, and the account's
// connector key is the latest one handed over.
// A phone that signs with its key alone, as apps did before KIFF kept
// devices, is refused, whatever it sends (kiff/planque#12). An enrolment
// stored then, with a phone's key in it, still serves KIFF devices.
func TestAKeyOnlyPhoneIsRefused(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	old := fakes.NewPhone() // no d=: its key alone
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/enrolment",
		map[string]string{"device_key": old.PublicKey(), "connector_key": "k_connector"}, old, ""); st != 401 {
		t.Fatalf("a key-only phone enrolling: %d, want 401", st)
	}
	id, digest := e.request()
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, old, ""); st != 401 {
		t.Fatalf("a key-only phone approving: %d, want 401", st)
	}
	// An enrolment written before #12 still has a phone's key: it is ignored.
	it, _, _ := e.kv.Get(context.Background(), tenantPK(tenant), skEnrolment)
	var doc map[string]any
	_ = json.Unmarshal(it.Doc, &doc)
	doc["device_key"] = old.PublicKey()
	it.Doc, _ = json.Marshal(doc)
	if err := e.kv.Put(context.Background(), it, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("the KIFF device on an enrolment from before #12: %d %v", st, out)
	}
}

// When KIFF can't be asked, the answer is KIFF's outage, not a refusal.
func TestKIFFDownIsNotARefusal(t *testing.T) {
	e := newEnv(t)
	ipad := e.device("dev_ipad")
	e.enrolAs(ipad)
	id, digest := e.request()
	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(503) }))
	t.Cleanup(down.Close)
	e.s.KIFF.API = down.URL
	st, out := e.do("POST", "/v1/accounts/"+tenant+"/grants/"+id+"/approve", map[string]string{"digest": digest}, ipad, "")
	if st != 502 || !strings.Contains(out["error"].(string), "KIFF could not check") {
		t.Fatalf("KIFF down: %d %v, want 502", st, out)
	}
}

// A malformed device id is a malformed header.
func TestADeviceIDIsKIFFs(t *testing.T) {
	for _, h := range []string{"t=x,v1=y,d=iphone", "t=x,v1=y,d=dev_", "t=x,v1=y,d=dev_ABC", "t=x,v1=y,d=dev_a/b"} {
		if _, err := device.ParseHeader(h); err == nil {
			t.Errorf("%q parsed", h)
		}
	}
	if s, err := device.ParseHeader("t=2026-10-19T10:00:00Z,v1=c2ln,d=dev_abc234"); err != nil || s.Device != "dev_abc234" {
		t.Fatalf("a good header: %+v %v", s, err)
	}
}

// A device KIFF gives no name for reads neutrally, never its id (#15): it may be another member's.
func TestAnUnnamedDeviceReadsWithoutItsID(t *testing.T) {
	a := Approval{DeviceID: "dev_x7", Subject: "user_1"}
	if got := a.Who(); got != "on a device signed in to KIFF" {
		t.Fatalf("Who() = %q", got)
	}
	if got := (Approval{Fingerprint: "ab12 cd34 ef56 7890"}).Who(); got != "on the enrolled phone" {
		t.Fatalf("a key-enrolled phone: %q", got)
	}
}
