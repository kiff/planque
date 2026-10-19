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
)

// device adds a device of tenant's owner to KIFF and returns it, signing
// with its KIFF id (kiff/planque#10).
func (e *env) device(id string) *fakes.Phone {
	p := fakes.NewPhone()
	p.DeviceID = id
	e.kiff.Devices[id] = fakes.Device{Tenant: tenant, Subject: "owner", Platform: "ios", Phone: p, Member: true,
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
	var en Enrolment
	if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &en); err != nil || len(en.DeviceKey) != 0 {
		t.Fatalf("a KIFF device left a key in the enrolment: %v", err)
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
func TestAPhoneEnrolledWithItsKeyAndAKIFFDevice(t *testing.T) {
	e := newEnv(t)
	e.enrol() // e.phone, with its key, before kiff/planque#10
	e.kiff.Keys["k_connector2"] = e.kiff.Keys["k_connector"]
	ipad := e.device("dev_ipad")
	if st, out := e.do("POST", "/v1/accounts/"+tenant+"/enrolment", map[string]string{"connector_key": "k_connector2"}, ipad, ""); st != 201 {
		t.Fatalf("the iPad, next to the enrolled phone: %d %v", st, out)
	}
	if k, _ := e.s.connectorKey(context.Background(), tenant); k != "k_connector2" {
		t.Fatalf("connector key %q, want the latest", k)
	}
	id, digest := e.request()
	if st, out := e.approveAndSignIn(id, digest); st != 200 || out["status"] != StatusConnected {
		t.Fatalf("the phone enrolled with its key: %d %v", st, out)
	}
	// An old app on yet another phone still can't enrol its own key.
	other := fakes.NewPhone()
	if st, _ := e.do("POST", "/v1/accounts/"+tenant+"/enrolment",
		map[string]string{"device_key": other.PublicKey(), "connector_key": "k_connector"}, other, ""); st != 409 {
		t.Fatalf("another phone's key: %d, want 409", st)
	}
	// Nor can the phone's key, if a KIFF device set the account up.
	e2 := newEnv(t)
	e2.enrolAs(e2.device("dev_iphone"))
	if st, _ := e2.do("POST", "/v1/accounts/"+tenant+"/enrolment",
		map[string]string{"device_key": e2.phone.PublicKey(), "connector_key": "k_connector"}, e2.phone, ""); st != 409 {
		t.Fatalf("a key after a KIFF device: %d, want 409", st)
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
