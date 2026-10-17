// Command fakekiff runs the KIFF stand-in (internal/fakes) on Lambda for
// the live check, until KIFF ships the connector key role. It is a test
// tool, never part of planque's deployment.
//
//	FAKEKIFF_ISSUER  its own public URL
//	FAKEKIFF_SEED    base64url Ed25519 seed for permits
//	FAKEKIFF_KEYS    JSON: KIFF key → /v1/me answer
//	FAKEKIFF_CARDS   JSON: agent key → action → most per call
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/lambda"
)

func main() {
	k := fakes.NewKIFF()
	k.Issuer = strings.TrimRight(os.Getenv("FAKEKIFF_ISSUER"), "/")
	seed, err := base64.RawURLEncoding.DecodeString(os.Getenv("FAKEKIFF_SEED"))
	if err != nil || len(seed) != ed25519.SeedSize {
		log.Fatal("FAKEKIFF_SEED must be a base64url 32-byte seed")
	}
	k.Key, k.KID = ed25519.NewKeyFromSeed(seed), "live-1"
	var keys map[string]kiff.Me
	if err := json.Unmarshal([]byte(os.Getenv("FAKEKIFF_KEYS")), &keys); err != nil {
		log.Fatal("FAKEKIFF_KEYS: ", err)
	}
	k.Keys = keys
	if err := json.Unmarshal([]byte(os.Getenv("FAKEKIFF_CARDS")), &k.Cards); err != nil {
		log.Fatal("FAKEKIFF_CARDS: ", err)
	}
	if lambda.Running() {
		log.Fatal(lambda.Start(k.Handler()))
	}
	log.Fatal(http.ListenAndServe(":8081", k.Handler()))
}
