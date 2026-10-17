// Command fakevendor runs the fake vendor (internal/fakes) on Lambda, for
// the live check and kiff/agents' nightly simulated-owner test: an MCP server
// behind its own OAuth (discovery, client registration, PKCE), as Stripe's
// is, that consents on its own. It is a test tool, never part of planque.
//
//	FAKEVENDOR_BASE   its own public URL
//	FAKEVENDOR_TOKEN  the bearer GET /fake/calls takes: what tools/call it ran
//	FAKEVENDOR_TABLE  the DynamoDB table its state is in, so any of its
//	                  instances answers (in memory without it)
package main

import (
	"cmp"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/lambda"
	"github.com/kiff/planque/internal/store"
)

func main() {
	base, token := os.Getenv("FAKEVENDOR_BASE"), os.Getenv("FAKEVENDOR_TOKEN")
	if base == "" || len(token) < 16 {
		log.Fatal("FAKEVENDOR_BASE and a FAKEVENDOR_TOKEN of 16 or more characters are required")
	}
	v, h := fakes.NewVendorAt(base)
	if t := os.Getenv("FAKEVENDOR_TABLE"); t != "" {
		v.KV = &store.Dynamo{Table: t, API: &awsapi.Client{Service: "dynamodb", Region: cmp.Or(os.Getenv("AWS_REGION"), "eu-west-1"),
			Target: "DynamoDB_20120810", JSONType: "application/x-amz-json-1.0", Creds: awsapi.FromEnv}}
	}
	mux := http.NewServeMux()
	mux.Handle("/", h)
	mux.HandleFunc("GET /fake/calls", func(w http.ResponseWriter, r *http.Request) {
		got, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		calls, err := v.CallsSoFar(r.Context())
		if err != nil {
			log.Printf("reading the calls: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"calls": calls})
	})
	if lambda.Running() {
		log.Fatal(lambda.Start(mux))
	}
	log.Fatal(http.ListenAndServe(":8082", mux))
}
