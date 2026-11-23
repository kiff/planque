package planque

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The extension's manifest_hash and planque's allowlist are one contract:
// SHA-256 of the packed manifest.json, as 64 lower-case hex characters, in the
// pairing's key post and the relay's hello. On 22 Nov 2026 the extension sent
// base64url while deploy/chrome-builds.txt keeps hex, so no build could pair:
// both sides' tests used a placeholder hash. This test takes each of the
// extension's real manifests (testdata/chrome-manifest-<version>.json, kiff/tafs
// chrome/manifest.json byte for byte: 0.1.0 at 78c86cb4, 0.2.0 at e08e5d4),
// hashes it as the extension's manifestHashOf does (chrome/src/crypto.js),
// finds it in the shipped allowlist, and sends it through the real route. A
// release that changes the manifest adds its file here and its line to the
// allowlist in the same PR.
func TestTheExtensionsRealManifestHashIsInTheShippedAllowlist(t *testing.T) {
	files, err := filepath.Glob("testdata/chrome-manifest-*.json")
	if err != nil || len(files) < 2 {
		t.Fatalf("the extension's manifests: %v %v", files, err)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) { manifestPairs(t, file) })
	}
}

func manifestPairs(t *testing.T, file string) {
	manifest, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(manifest, &m); err != nil || m.Version == "" || filepath.Base(file) != "chrome-manifest-"+m.Version+".json" {
		t.Fatalf("the manifest's version: %q %v (in %s)", m.Version, err, file)
	}
	f, err := os.Open("../../deploy/chrome-builds.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	shipped, err := ChromeBuildsFrom(f)
	if err != nil {
		t.Fatalf("the shipped allowlist doesn't parse: %v", err)
	}

	sum := sha256.Sum256(manifest)
	sent := hex.EncodeToString(sum[:]) // as the extension sends it: 64 lower-case hex
	if b, ok := shipped[sent]; !ok || b.Version != m.Version {
		t.Fatalf("the extension's %s (%s) isn't in deploy/chrome-builds.txt: %v", sent, m.Version, shipped)
	}

	// Through the route, with the shipped allowlist: the hex hash passes the
	// build check, and the same bytes in another encoding never do.
	e := chromeEnv(t)
	e.s.ChromeBuilds = shipped
	// Each pairing is opened with its own signature: a second apart.
	base, n := time.Now(), 0
	e.phone.Now = func() time.Time { n++; return base.Add(time.Duration(n) * time.Second) }
	post := func(hash string) (int, map[string]any) {
		st, out := e.do("POST", "/v1/accounts/"+tenant+"/chromes/pairings", map[string]any{}, e.phone, "")
		if st != 201 {
			t.Fatalf("open: %d %v", st, out)
		}
		key, _ := p256Key(t)
		return e.do("POST", "/v1/extension/pairings/"+out["lookup"].(string), map[string]any{"key": key, "version": m.Version,
			"manifest_hash": hash, "browser": "Chrome 131", "os": "macOS 15", "proof": b64(32)}, nil, "")
	}
	if st, out := post(sent); st != 202 {
		t.Fatalf("the real build, as the extension sends its hash: %d %v", st, out)
	}
	for name, other := range map[string]string{
		"base64url":       base64.RawURLEncoding.EncodeToString(sum[:]),
		"upper-case hex":  strings.ToUpper(sent),
		"base64 (padded)": base64.StdEncoding.EncodeToString(sum[:]),
	} {
		if st, out := post(other); st != 403 {
			t.Errorf("the same build's hash as %s: %d %v, want 403", name, st, out)
		}
	}
}
