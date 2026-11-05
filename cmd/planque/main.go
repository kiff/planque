// Command planque is the one binary: on AWS Lambda behind a function
// URL, or as an HTTP server anywhere else.
//
// Configuration (environment):
//
//	PLANQUE_BASE_URL   planque's public URL, e.g. its function URL (required)
//	PLANQUE_APP_RETURN where a finished vendor sign-in sends the owner's
//	                   browser with its confirmation code: the owner's app,
//	                   e.g. tafs://planque/signed-in (required)
//	PLANQUE_TABLE      DynamoDB table; without it records stay in memory
//	PLANQUE_KMS_KEY    KMS key id or alias for envelope encryption, or
//	PLANQUE_LOCAL_KEY  a base64 32-byte AES key (self-hosted)
//	AWS_REGION         the region of the table and key (eu-west-1)
//	KIFF_API           default https://api.kiff.dev
//	KIFF_GATEWAY       default https://mcp.kiff.dev
//	KIFF_ISSUER        the permit issuer, default KIFF_API
//	KIFF_JWKS_URL      default KIFF_ISSUER/.well-known/kiff-permit-keys.json
//	PLANQUE_DOOR_URL fixed HTTPS receiver for approved incoming door events
//	PLANQUE_DOOR_SECRET Secrets Manager reference containing {secret}
//	PLANQUE_EVENT_URL fixed HTTPS receiver for approved app events, signed
//	                   with PLANQUE_DOOR_SECRET; unset, no app event can be approved
//	PLANQUE_COMPOSIO_SECRET optional Secrets Manager reference for {api_key}
//	PLANQUE_ADDR       listen address outside Lambda, default :8080
//	PLANQUE_APPS_FILE  optional credential-free registered app definitions
//	PLANQUE_APP_SECRETS optional JSON map of app ids to secret references
package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kiff/planque/internal/appconfig"
	"github.com/kiff/planque/internal/appcredentials"
	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/browser"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/door"
	"github.com/kiff/planque/internal/egress"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/lambda"
	"github.com/kiff/planque/internal/mcp"
	"github.com/kiff/planque/internal/oauth"
	"github.com/kiff/planque/internal/permit"
	"github.com/kiff/planque/internal/planque"
	"github.com/kiff/planque/internal/rest"
	"github.com/kiff/planque/internal/seal"
	"github.com/kiff/planque/internal/secretstore"
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
	appReturn := os.Getenv("PLANQUE_APP_RETURN")
	if u, err := url.Parse(appReturn); appReturn == "" || err != nil || u.Scheme == "" {
		return nil, errors.New("PLANQUE_APP_RETURN is required: the owner's app, which confirms each sign-in")
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
	applications := map[string]appconfig.Application{}
	if path := os.Getenv("PLANQUE_APPS_FILE"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New("the registered application definitions could not be read")
		}
		applications, err = appconfig.Decode(raw)
		if err != nil {
			return nil, err
		}
	}
	catalog := map[string]planque.CatalogEntry{}
	if path := os.Getenv("PLANQUE_CATALOG_FILE"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New("the vendor catalog could not be read")
		}
		catalog, err = planque.DecodeCatalog(raw, applications)
		if err != nil {
			return nil, err
		}
	}
	var applicationCredentials appcredentials.Reader = &appcredentials.Sealed{KV: kv, Sealer: sealer}
	if raw := os.Getenv("PLANQUE_APP_SECRETS"); raw != "" {
		var references map[string]string
		if json.Unmarshal([]byte(raw), &references) != nil || len(references) == 0 {
			return nil, errors.New("PLANQUE_APP_SECRETS must map application ids to secret references")
		}
		applicationCredentials = &appcredentials.SecretsManager{References: references, API: &awsapi.Client{
			Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: creds}}
	}
	api := cmp.Or(os.Getenv("KIFF_API"), "https://api.kiff.dev")
	issuer := cmp.Or(os.Getenv("KIFF_ISSUER"), api)
	// KIFF is configured here; vendors are named by agents and servers, so
	// they are reached only at public addresses.
	hc := &http.Client{Timeout: 30 * time.Second}
	vendors := egress.Client(30 * time.Second)
	connectors := map[string]connector.Provider{}
	if reference := strings.TrimSpace(os.Getenv("PLANQUE_COMPOSIO_SECRET")); reference != "" {
		connectors["composio"] = &composio.Client{HTTP: vendors, Secret: composio.SecretFromStore(&awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: creds}, reference)}
	}
	var doorSink door.Sink
	doorURL, doorSecret := strings.TrimSpace(os.Getenv("PLANQUE_DOOR_URL")), strings.TrimSpace(os.Getenv("PLANQUE_DOOR_SECRET"))
	if doorURL != "" || doorSecret != "" {
		u, err := url.Parse(doorURL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || doorSecret == "" {
			return nil, errors.New("PLANQUE_DOOR_URL and PLANQUE_DOOR_SECRET must name the fixed HTTPS receiver and its secret reference")
		}
		reader := secretstore.Reader(&awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: creds}, doorSecret, "secret")
		doorSink = &door.HTTP{URL: doorURL, Secret: reader, Client: egress.Client(10 * time.Second)}
	}
	var eventSink door.EventSink
	if eventURL := strings.TrimSpace(os.Getenv("PLANQUE_EVENT_URL")); eventURL != "" {
		u, err := url.Parse(eventURL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || doorSecret == "" {
			return nil, errors.New("PLANQUE_EVENT_URL must name the fixed HTTPS receiver, and PLANQUE_DOOR_SECRET the secret it's signed with")
		}
		reader := secretstore.Reader(&awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: creds}, doorSecret, "secret")
		eventSink = &door.HTTP{URL: eventURL, Secret: reader, Client: egress.Client(10 * time.Second)}
	}
	// The browser worker (kiff/planque#86): where it listens, on a private
	// address, and the key it checks. Without both, no browser is offered.
	var browserWorker browser.Worker
	workerURL, workerKey := strings.TrimSpace(os.Getenv("PLANQUE_BROWSER_URL")), strings.TrimSpace(os.Getenv("PLANQUE_BROWSER_KEY"))
	if workerURL != "" || workerKey != "" {
		u, err := url.Parse(workerURL)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || workerKey == "" {
			return nil, errors.New("PLANQUE_BROWSER_URL and PLANQUE_BROWSER_KEY must name the browser worker and the key it checks")
		}
		browserWorker = &browser.Client{URL: strings.TrimRight(workerURL, "/"), Key: workerKey, HTTP: &http.Client{Timeout: 90 * time.Second}}
	}
	return &planque.Server{
		Browser:    browserWorker,
		DoorSink:   doorSink,
		EventSink:  eventSink,
		Connectors: connectors,
		BaseURL:    base,
		AppReturn:  appReturn,
		KV:         kv,
		Sealer:     sealer,
		KIFF:       &kiff.Client{API: api, Gateway: cmp.Or(os.Getenv("KIFF_GATEWAY"), "https://mcp.kiff.dev"), HTTP: hc},
		Permits: &permit.Verifier{Issuer: issuer, Keys: &permit.JWKS{
			URL: cmp.Or(os.Getenv("KIFF_JWKS_URL"), issuer+"/.well-known/kiff-permit-keys.json"), HTTP: hc}},
		OAuth:                  &oauth.Client{HTTP: vendors},
		MCP:                    &mcp.Client{HTTP: vendors},
		REST:                   &rest.Client{HTTP: vendors},
		Logger:                 logger,
		WindowRequired:         strings.EqualFold(strings.TrimSpace(os.Getenv("PLANQUE_WINDOW_REQUIRED")), "true"),
		CatalogEntries:         catalog,
		Applications:           applications,
		ApplicationCredentials: applicationCredentials,
	}, nil
}
