package planque

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/kiff/planque/internal/stepkey"
)

// The relay's contract, shared with the extension: kiff/tafs
// vectors/chrome-relay.json, vendored byte for byte. The extension's
// chrome/test/vectors.test.js reads the same file, so a field either side
// signs, sends or checks differently fails one of the two suites.
const (
	relayVectorPath = "testdata/chrome-relay.json"
	// relayVectorSHA256 is the vendored file's hash, from kiff/tafs at
	// relayVectorCommit: a copy edited here, not taken again from there, fails.
	relayVectorSHA256 = "915c4fbac1bc0a1b4639042eaca8f83724a32869679281daf8796b33f1e4548c"
	relayVectorCommit = "c54c620f0d14345593a41a50cd346c4cfeff1300"
)

type relayVectorEntry struct {
	Name         string                     `json:"name"`
	Label        string                     `json:"label"`
	Body         json.RawMessage            `json:"body"`
	SigningBytes string                     `json:"signing_bytes"`
	Sig          string                     `json:"sig"`
	Wire         map[string]json.RawMessage `json:"wire"`
}

type relayVector struct {
	PlanqueKey   string             `json:"planque_key"`
	ExtensionKey string             `json:"extension_key"`
	ChromeID     string             `json:"chrome_id"`
	Steps        []relayVectorEntry `json:"steps"`
	Results      []relayVectorEntry `json:"results"`
	HandResults  []relayVectorEntry `json:"hand_results"`
}

func readRelayVector(t *testing.T) relayVector {
	t.Helper()
	raw, err := os.ReadFile(relayVectorPath)
	if err != nil {
		t.Fatal(err)
	}
	if sum := sha256.Sum256(raw); hex.EncodeToString(sum[:]) != relayVectorSHA256 {
		t.Fatalf("%s isn't kiff/tafs' vectors/chrome-relay.json at %s: take it again from there and pin its hash", relayVectorPath, relayVectorCommit[:7])
	}
	var v relayVector
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	if len(v.Steps) == 0 || len(v.Results) == 0 || len(v.HandResults) == 0 {
		t.Fatal("the vector has no steps, results or hand results")
	}
	return v
}

// signedBytes checks that label ‖ 0 ‖ JCS(body) is the vector's signing_bytes.
func signedBytes(t *testing.T, e relayVectorEntry, body any) {
	t.Helper()
	got, err := stepkey.Bytes(e.Label, body)
	if err != nil {
		t.Fatal(err)
	}
	want, err := base64.RawURLEncoding.DecodeString(e.SigningBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s: planque signs\n%q\nthe extension\n%q", e.Name, got, want)
	}
}

func TestEveryStepInTheRelayVectorIsOnePlanqueSignsAndSends(t *testing.T) {
	v := readRelayVector(t)
	for _, e := range v.Steps {
		if e.Label != labelStep {
			t.Fatalf("%s: label %q", e.Name, e.Label)
		}
		// The step as it travels decodes as planque's chromeStep, and planque
		// signs exactly the extension's bytes for it: no field is lost or added.
		var wireType string
		_ = json.Unmarshal(e.Wire["type"], &wireType)
		if wireType != "step" {
			t.Fatalf("%s: a step travels as {type: \"step\", step, sig}, got type %q", e.Name, wireType)
		}
		var st chromeStep
		dec := json.NewDecoder(bytes.NewReader(e.Wire["step"]))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&st); err != nil {
			t.Fatalf("%s: the step doesn't decode as planque's: %v", e.Name, err)
		}
		signedBytes(t, e, st)
		var sig string
		_ = json.Unmarshal(e.Wire["sig"], &sig)
		if sig != e.Sig || !stepkey.Verify(v.PlanqueKey, labelStep, st, sig) {
			t.Fatalf("%s: the step's signature doesn't verify under planque_key", e.Name)
		}
		if st.ChromeID != v.ChromeID {
			t.Fatalf("%s: chrome_id %q", e.Name, st.ChromeID)
		}
	}
}

func TestEveryResultInTheRelayVectorIsOnePlanqueVerifiesAndReads(t *testing.T) {
	v := readRelayVector(t)
	for _, e := range v.Results {
		if e.Label != labelResult {
			t.Fatalf("%s: label %q", e.Name, e.Label)
		}
		// What relayResult verifies: the body built from the wire message.
		body, ok := resultBody(e.Wire)
		if !ok {
			t.Fatalf("%s: a field relayResult needs is missing", e.Name)
		}
		signedBytes(t, e, body)
		var sig string
		_ = json.Unmarshal(e.Wire["sig"], &sig)
		if !stepkey.Verify(v.ExtensionKey, labelResult, body, sig) {
			t.Fatalf("%s: the result doesn't verify under extension_key", e.Name)
		}
		// A changed field breaks it: the signature covers the whole body.
		tampered := map[string]json.RawMessage{}
		for k, x := range body {
			tampered[k] = x
		}
		tampered["chrome_id"] = json.RawMessage(`"chr_other"`)
		if stepkey.Verify(v.ExtensionKey, labelResult, tampered, sig) {
			t.Fatalf("%s: verified with another chrome_id", e.Name)
		}

		// Then read as the step it answers.
		var okField bool
		_ = json.Unmarshal(body["ok"], &okField)
		_, handed := body["handed"]
		switch {
		case handed:
			if err := readHandedResult(body); err != nil {
				t.Fatalf("%s: a handed answer reads as %v", e.Name, err)
			}
			delete(body, "handed")
			if err := readHandedResult(body); err == nil {
				t.Fatalf("%s: an answer without handed was taken as handed", e.Name)
			}
		case okField:
			if _, err := readChromeResult(body); err != nil {
				t.Fatalf("%s: a page reads as %v", e.Name, err)
			}
		default:
			_, err := readChromeResult(body)
			var r *chromeRefusal
			if !errors.As(err, &r) || r.Code == "" {
				t.Fatalf("%s: a refusal reads as %v", e.Name, err)
			}
		}
	}
}

func TestEveryHandResultInTheRelayVectorIsOnePlanqueVerifies(t *testing.T) {
	v := readRelayVector(t)
	for _, e := range v.HandResults {
		if e.Label != labelHand {
			t.Fatalf("%s: label %q", e.Name, e.Label)
		}
		signed, sig, ok := handSigned(e.Wire)
		if !ok {
			t.Fatalf("%s: relayHandResult reads no hand-over or signature", e.Name)
		}
		signedBytes(t, e, signed)
		if !stepkey.Verify(v.ExtensionKey, labelHand, signed, sig) {
			t.Fatalf("%s: the hand result doesn't verify under extension_key", e.Name)
		}
		signed.Outcome += "x"
		if stepkey.Verify(v.ExtensionKey, labelHand, signed, sig) {
			t.Fatalf("%s: verified with another outcome", e.Name)
		}
	}
}
