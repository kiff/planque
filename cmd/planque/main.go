// Command planque is the one binary: on AWS Lambda behind a function
// URL, or as an HTTP server anywhere else.
//
// Configuration (environment):
//
//	PLANQUE_BASE_URL   planque's public URL, e.g. its function URL (required)
//	PLANQUE_TABLE      DynamoDB table; without it records stay in memory
//	PLANQUE_KMS_KEY    KMS key id or alias for envelope encryption, or
//	PLANQUE_LOCAL_KEY  a base64 32-byte AES key (self-hosted)
//	AWS_REGION         the region of the table and key (eu-west-1)
//	KIFF_API           default https://api.kiff.dev
//	KIFF_GATEWAY       default https://mcp.kiff.dev
//	KIFF_ISSUER        the permit issuer, default KIFF_API
//	KIFF_JWKS_URL      default KIFF_ISSUER/.well-known/kiff-permit-keys.json
//	PLANQUE_ADDR       listen address outside Lambda, default :8080
package main

import (
	"cmp"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/lambda"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/planque"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	s, err := build(logger)
	if err != nil {
		log.Fatal(err)
	}
	if lambda.Running() {
		log.Fatal(lambda.Start(s.Handler()))
	}
	addr := cmp.Or(os.Getenv("PLANQUE_ADDR"), ":8080")
	logger.Info("planque listening", slog.String("addr", addr))
	srv := &http.Server{Addr: addr, Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func build(logger *slog.Logger) (*planque.Server, error) {
	base := strings.TrimRight(os.Getenv("PLANQUE_BASE_URL"), "/")
	if base == "" {
		return nil, errors.New("PLANQUE_BASE_URL is required")
	}
	region := cmp.Or(os.Getenv("AWS_REGION"), "eu-west-1")
	creds := awsapi.FromEnv
	var kv store.KV = &store.Memory{}
	if t := os.Getenv("PLANQUE_TABLE"); t != "" {
		kv = &store.Dynamo{Table: t, API: &awsapi.Client{Service: "dynamodb", Region: region,
			Target: "DynamoDB_20120810", JSONType: "application/x-amz-json-1.0", Creds: creds}}
	}
	var sealer seal.Sealer
	switch {
	case os.Getenv("PLANQUE_KMS_KEY") != "":
		sealer = &seal.KMS{KeyID: os.Getenv("PLANQUE_KMS_KEY"), API: &awsapi.Client{Service: "kms", Region: region,
			Target: "TrentService", JSONType: "application/x-amz-json-1.1", Creds: creds}}
	case os.Getenv("PLANQUE_LOCAL_KEY") != "":
		l, err := seal.NewLocal(os.Getenv("PLANQUE_LOCAL_KEY"))
		if err != nil {
			return nil, err
		}
		sealer = l
	default:
		return nil, errors.New("set PLANQUE_KMS_KEY or PLANQUE_LOCAL_KEY")
	}
	api := cmp.Or(os.Getenv("KIFF_API"), "https://api.kiff.dev")
	issuer := cmp.Or(os.Getenv("KIFF_ISSUER"), api)
	hc := &http.Client{Timeout: 30 * time.Second}
	return &planque.Server{
		BaseURL: base,
		KV:      kv,
		Sealer:  sealer,
		KIFF:    &kiff.Client{API: api, Gateway: cmp.Or(os.Getenv("KIFF_GATEWAY"), "https://mcp.kiff.dev"), HTTP: hc},
		Permits: &permit.Verifier{Issuer: issuer, Keys: &permit.JWKS{
			URL: cmp.Or(os.Getenv("KIFF_JWKS_URL"), issuer+"/.well-known/kiff-permit-keys.json"), HTTP: hc}},
		OAuth:  &oauth.Client{HTTP: hc},
		MCP:    &mcp.Client{HTTP: hc},
		Logger: logger,
	}, nil
}
