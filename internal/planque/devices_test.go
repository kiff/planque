package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
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
		Capabilities: []string{"holds.answer", "tools.connect", "account.export", "account.delete"}}
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
	if st, out := e.enrolAs(ipad); st != 200 || out["kept"] != true {
		t.Fatalf("enrol from the iPad, a second device: %d %v, want 200 kept", st, out)
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

// A device KIFF refuses is logged with why, so the refusal can be told apart
// from the 401 the phone reads: here another account's device, after the
// phone switched accounts and kept its old id.
func TestARefusedDeviceIsLoggedWithItsReason(t *testing.T) {
	e := newEnv(t)
	var logs bytes.Buffer
	e.s.Logger = slog.New(slog.NewJSONHandler(&logs, nil))
	other := e.device("dev_other")
	d := e.kiff.Devices["dev_other"]
	d.Tenant = "tnt_other"
	e.kiff.Devices["dev_other"] = d
	if st, _ := e.enrolAs(other); st != 401 {
		t.Fatalf("enrol from another account's device: %d, want 401", st)
	}
	var line map[string]any
	if err := json.Unmarshal(logs.Bytes(), &line); err != nil {
		t.Fatalf("one JSON log line: %q %v", logs.String(), err)
	}
	for k, want := range map[string]string{"msg": "device refused", "tenant": tenant, "device": "dev_other", "reason": "unknown_device", "capability": "tools.connect"} {
		if line[k] != want {
			t.Errorf("log %s = %v, want %q", k, line[k], want)
		}
	}
	if strings.Contains(logs.String(), "k_connector") || strings.Contains(logs.String(), "v1=") {
		t.Fatalf("the log carries a key or a signature: %s", logs.String())
	}
}

// A second device's setup hands planque a new connector key, but the tools
// the enrolled key connected are that key's in KIFF: planque keeps it while
// KIFF takes it, and answers kept. Once KIFF no longer takes it (revoked),
// a new key replaces it.
func TestTheEnrolledConnectorKeyIsKeptWhileItWorks(t *testing.T) {
	e := newEnv(t)
	iphone, ipad := e.device("dev_iphone"), e.device("dev_ipad")
	e.kiff.Keys["k_connector2"] = e.kiff.Keys["k_connector"]
	enrol := func(p *fakes.Phone, key string) (int, map[string]any) {
		return e.do("POST", "/v1/accounts/"+tenant+"/enrolment", map[string]string{"connector_key": key}, p, "")
	}
	stored := func() string {
		var en Enrolment
		if _, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &en); err != nil {
			t.Fatal(err)
		}
		k, err := e.s.openConnector(context.Background(), en)
		if err != nil {
			t.Fatal(err)
		}
		return k
	}
	if st, out := enrol(iphone, "k_connector"); st != 201 {
		t.Fatalf("first enrolment: %d %v", st, out)
	}
	if st, out := enrol(ipad, "k_connector2"); st != 200 || out["kept"] != true || stored() != "k_connector" {
		t.Fatalf("a second device's key: %d %v, stored %s", st, out, stored())
	}
	// A new key must still be a good one: kept or not, a bad one is refused.
	if st, _ := enrol(ipad, "k_admin"); st != 400 {
		t.Fatalf("a bad key: %d, want 400", st)
	}
	// The enrolled key revoked in KIFF: the new one replaces it.
	delete(e.kiff.Keys, "k_connector")
	if st, out := enrol(iphone, "k_connector2"); st != 201 || out["kept"] != nil || stored() != "k_connector2" {
		t.Fatalf("after the enrolled key was revoked: %d %v, stored %s", st, out, stored())
	}
}

// kiff/tafs#494: a routine call (approving a sign-in) signed by the device
// key inside an approval window KIFF says is open goes through, and the
// record names the key and the window; one naming a window that isn't open
// is refused with window_required; the approval key signs it with no
// window; a safe call (disconnecting) takes the device key whatever the
// window; and with the window required, a device-key signature naming no
// window is refused.
func TestARoutineCallIsSignedInsideAnApprovalWindow(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	approval := fakes.NewPhone()
	approval.DeviceID = "dev_phone"
	d := e.kiff.Devices["dev_phone"]
	d.Approval = approval
	e.kiff.Devices["dev_phone"] = d
	if st, _ := e.enrolAs(phone); st != 201 {
		t.Fatal("enrol")
	}
	window := "win_" + strings.Repeat("a", 26)
	e.kiff.Windows[window] = "dev_phone"
	id, digest := e.request()
	approve := "/v1/accounts/" + tenant + "/grants/" + id + "/approve"
	// A window that isn't open for this device: refused, with the code the app acts on.
	phone.Window = "win_" + strings.Repeat("b", 26)
	if st, out := e.do("POST", approve, map[string]string{"digest": digest}, phone, ""); st != 403 || out["code"] != "window_required" {
		t.Fatalf("a window that isn't open: %d %v", st, out)
	}
	// The open window: a tap.
	phone.Window = window
	st, out := e.do("POST", approve, map[string]string{"digest": digest}, phone, "")
	if st != 200 {
		t.Fatalf("approve inside the window: %d %v", st, out)
	}
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Approval == nil || g.Approval.Key != "device" || g.Approval.Window != window || !strings.HasPrefix(g.Approval.Message, "planque.request.v2\n") || !strings.HasSuffix(g.Approval.Message, "\n"+window) {
		t.Fatalf("the approval doesn't name the key and the window: %+v", g.Approval)
	}
	// The approval key, no window: critical, and the record says so.
	id2, digest2 := e.request()
	approve2 := "/v1/accounts/" + tenant + "/grants/" + id2 + "/approve"
	if st, out := e.do("POST", approve2, map[string]string{"digest": digest2}, approval, ""); st != 200 {
		t.Fatalf("approve with the approval key: %d %v", st, out)
	}
	if g2, _ := e.s.loadGrant(context.Background(), tenant, id2); g2.Approval == nil || g2.Approval.Key != "critical" || g2.Approval.Window != "" {
		t.Fatalf("the approval key's record: %+v", g2.Approval)
	}
	// A safe call takes the device key, window open or not, and never asks KIFF about the window.
	phone.Window = "win_" + strings.Repeat("c", 26)
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id2, nil, phone, ""); st != 200 {
		t.Fatalf("disconnect with a closed window named: %d", st)
	}
	if asked := e.kiff.WindowsAsked; len(asked) == 0 || asked[len(asked)-1] != "" {
		t.Fatalf("a safe route sent the window to KIFF: %v", asked)
	}
	// With the window required, a device-key signature naming none is refused; the approval key still goes.
	e.s.WindowRequired = true
	phone.Window = ""
	id3, digest3 := e.request()
	approve3 := "/v1/accounts/" + tenant + "/grants/" + id3 + "/approve"
	if st, out := e.do("POST", approve3, map[string]string{"digest": digest3}, phone, ""); st != 403 || out["code"] != "window_required" {
		t.Fatalf("the device key alone, window required: %d %v", st, out)
	}
	if st, _ := e.do("POST", approve3, map[string]string{"digest": digest3}, approval, ""); st != 200 {
		t.Fatal("the approval key, window required")
	}
	// Disconnecting stays the device key's.
	if st, _ := e.do("DELETE", "/v1/accounts/"+tenant+"/grants/"+id3, nil, phone, ""); st != 200 {
		t.Fatalf("disconnect, window required: %d", st)
	}
}

// An app enrols first with no connector key: planque answers kept while
// KIFF still takes the enrolled key, and needs_key otherwise, writing
// nothing either way. The app mints a key only on needs_key, so no key is
// made that nothing uses (KIFF lets no device revoke one). A request with
// no key never creates an enrolment, and it's verified as the signed
// routes are, against the enrolled connector.
func TestAnEnrolmentWithNoKeyAsksWhetherOneIsNeeded(t *testing.T) {
	e := newEnv(t)
	iphone, ipad := e.device("dev_iphone"), e.device("dev_ipad")
	enrol := func(p *fakes.Phone, body map[string]string) (int, map[string]any) {
		return e.do("POST", "/v1/accounts/"+tenant+"/enrolment", body, p, "")
	}
	enrolment := func() (Enrolment, bool) {
		var en Enrolment
		_, err := getDoc(context.Background(), e.kv, tenantPK(tenant), skEnrolment, &en)
		return en, err == nil
	}
	// Nothing enrolled: needs_key, and nothing is written.
	if st, out := enrol(iphone, map[string]string{}); st != 409 || out["code"] != "needs_key" {
		t.Fatalf("no enrolment, no key: %d %v", st, out)
	}
	if _, ok := enrolment(); ok {
		t.Fatal("a request with no key created an enrolment")
	}
	if st, out := enrol(iphone, map[string]string{"connector_key": "k_connector"}); st != 201 {
		t.Fatalf("enrolling with a key: %d %v", st, out)
	}
	before, _ := enrolment()
	// Enrolled, and KIFF takes its key: kept, nothing written, any device of the account.
	if st, out := enrol(ipad, map[string]string{"connector_key": ""}); st != 200 || out["kept"] != true {
		t.Fatalf("enrolled, no key: %d %v", st, out)
	}
	if after, _ := enrolment(); !after.EnrolledAt.Equal(before.EnrolledAt) || string(after.Connector) != string(before.Connector) {
		t.Fatal("a kept enrolment with no key was rewritten")
	}
	// Signed as the signed routes are: a device KIFF doesn't know gets no answer about the account.
	if st, out := enrol(fakes.NewPhone(), map[string]string{}); st == 200 || out["kept"] == true {
		t.Fatalf("an unknown device: %d %v", st, out)
	}
	// The enrolled key revoked in KIFF: needs_key, and the enrolment is left as it was.
	delete(e.kiff.Keys, "k_connector")
	if st, out := enrol(ipad, map[string]string{}); st != 409 || out["code"] != "needs_key" {
		t.Fatalf("the enrolled key revoked, no key: %d %v", st, out)
	}
	if after, _ := enrolment(); string(after.Connector) != string(before.Connector) {
		t.Fatal("needs_key changed the enrolment")
	}
}

// An app pairs each listed grant with its raw entry by position, so a
// repeated id would let one grant's digest stand for another's. planque
// never lists one: a record whose id isn't its key's, or repeats one read,
// lists nothing, with its own code, and the app signs nothing.
func TestTheGrantsAreNeverListedWithARepeatedID(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	id, _ := e.request()
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants", nil, e.phone, ""); st != 200 || len(out["grants"].([]any)) != 1 {
		t.Fatalf("one grant: %d %v", st, out)
	}
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	// The same record under a second key: two entries for one id.
	if err := putDoc(context.Background(), e.kv, store.Item{PK: tenantPK(tenant), SK: skGrant("grt_zzzzsecond")}, g, store.Cond{}); err != nil {
		t.Fatal(err)
	}
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants?include=disconnected", nil, e.phone, "")
	if st != 500 || out["code"] != "grant_repeated" || out["grants"] != nil {
		t.Fatalf("a repeated id was listed: %d %v", st, out)
	}
}

// Each 409 an app may meet when enrolling has its own code, so needs_key
// is never mistaken for another: a request already made, or a race.
func TestAReplayedSignatureHasItsOwnCode(t *testing.T) {
	e := newEnv(t)
	phone := e.device("dev_phone")
	body, _ := json.Marshal(map[string]string{"connector_key": "k_connector"})
	req, _ := http.NewRequest("POST", e.url+"/v1/accounts/"+tenant+"/enrolment", bytes.NewReader(body))
	phone.Sign(req, tenant, body)
	send := func() (int, map[string]any) {
		again, _ := http.NewRequest("POST", req.URL.String(), bytes.NewReader(body))
		again.Header = req.Header.Clone()
		resp, err := http.DefaultClient.Do(again)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var out map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&out)
		return resp.StatusCode, out
	}
	if st, out := send(); st != 201 {
		t.Fatalf("enrol: %d %v", st, out)
	}
	if st, out := send(); st != 409 || out["code"] != "signature_replayed" {
		t.Fatalf("the same signed request again: %d %v", st, out)
	}
}
