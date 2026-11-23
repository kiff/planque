// Command planque is the one binary: on AWS Lambda behind a function
// URL, or as an HTTP server anywhere else. `planque guides` writes every
// app's generated guide part instead (deploy/guides/README.md).
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
//	PLANQUE_APP_SECRET_PREFIX optional: an app's secret is the one named
//	                   this prefix plus its id (planque/apps/), unless
//	                   PLANQUE_APP_SECRETS names its own
//	PLANQUE_MANAGED_TEST_TENANTS optional tenants:<id>[,<id>]: the test
//	                   accounts that may make a managed connection to any app;
//	                   set at deploy, never kept in the repository
//	PLANQUE_CHROME_OFFER which accounts the owner's own Chrome is offered to,
//	                   once a build is listed: on (every one), test (only
//	                   PLANQUE_CHROME_TENANTS, tenants:<id>[,<id>]), or off;
//	                   unset, none
//	PLANQUE_CHROME_SCREENSHOTS on keeps an image of each acting step in the
//	                   owner's own Chrome, sealed, for 30 days; unset or
//	                   anything else, off
//	PLANQUE_CHROME_HAND off stops offering the hand-over's two tools
//	                   (browse_hand, browse_hand_result) to new connections
//	                   in the owner's Chrome; unset or anything else, on
//	PLANQUE_NEVER_USE  optional grants:<id>[,<id>]: grants whose sign-ins are
//	                   never reused or called, by any grant on the same
//	                   account at the app; set at deploy, never in the repository
//	PLANQUE_OWN_CLIENT_TENANTS optional tenants:<id>[,<id>]: accounts shown an
//	                   app's own entry instead of the managed connector's
//	                   toolkit for it; set at deploy, never kept in the repository
//	PLANQUE_SIGNIN_READER_SHA256 optional: the hex SHA-256 of the one key that
//	                   reads an app's managed sign-in state; unset, no one can
//	PLANQUE_DESCRIBER_SERVERS optional JSON map of HTTPS servers, besides
//	                   Stripe's own, to the describing tools their first
//	                   request gets with its sign-in: the Eval's vendor, which
//	                   serves Stripe's tools on a host we run; not secret
//	PLANQUE_CATALOG_OPERATORS optional JSON list of the catalog's operators:
//	                   {id, token_sha256, device_key}; unset, the catalog kept
//	                   as data can't be drafted or published
package main

import (
	"cmp"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
	"github.com/kiff/planque/internal/stepkey"
	"github.com/kiff/planque/internal/store"
	"github.com/kiff/planque/internal/wsapi"
)

// appSecretPrefix is a Secrets Manager name path: letters, digits and
// /_+=.@-, ending in a slash, so Prefix+id names one secret under it.
var appSecretPrefix = regexp.MustCompile(`^[A-Za-z0-9/_+=.@-]{1,200}/$`)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "guides" {
		if err := guides(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		return
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	s, err := build(logger)
	if err != nil {
		log.Fatal(err)
	}
	if lambda.Running() {
		log.Fatal(lambda.Start(s.Handler(), s))
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
	// Who may work on the catalog kept as data: each operator's credential
	// (by its SHA-256) and enrolled device, catalog-only.
	var operators map[string]planque.CatalogOperator
	if raw := os.Getenv("PLANQUE_CATALOG_OPERATORS"); raw != "" {
		var err error
		if operators, err = planque.DecodeCatalogOperators([]byte(raw)); err != nil {
			return nil, err
		}
	}
	// Which app a vendor's own server is: a reviewed list, read like the
	// catalog, shipped with the function by deploy/aws.sh.
	var serverApps map[string]string
	if path := os.Getenv("PLANQUE_SERVER_APPS_FILE"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New("the server apps could not be read")
		}
		if serverApps, err = planque.DecodeServerApps(raw); err != nil {
			return nil, err
		}
	}
	var describerServers map[string][]string
	if raw := os.Getenv("PLANQUE_DESCRIBER_SERVERS"); raw != "" {
		var err error
		if describerServers, err = planque.DecodeDescriberServers([]byte(raw)); err != nil {
			return nil, err
		}
	}
	// planque's one secret (PLANQUE_SECRET): the values below, read from it
	// first, each from its own secret while it lacks it or can't be read.
	secretsAPI := &awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: creds}
	bundle := secretstore.NewBundle(secretsAPI, strings.TrimSpace(os.Getenv("PLANQUE_SECRET")))
	// own reads one value from its own secret's field, as before; none when
	// it isn't named.
	own := func(reference, field string) func(context.Context) (string, error) {
		if reference == "" {
			return nil
		}
		return secretstore.Reader(secretsAPI, reference, field)
	}
	var applicationCredentials appcredentials.Reader = &appcredentials.Sealed{KV: kv, Sealer: sealer}
	raw, prefix := os.Getenv("PLANQUE_APP_SECRETS"), os.Getenv("PLANQUE_APP_SECRET_PREFIX")
	if raw != "" || prefix != "" {
		var references map[string]string
		if raw != "" && (json.Unmarshal([]byte(raw), &references) != nil || len(references) == 0) {
			return nil, errors.New("PLANQUE_APP_SECRETS must map application ids to secret references")
		}
		if prefix != "" && !appSecretPrefix.MatchString(prefix) {
			return nil, errors.New("PLANQUE_APP_SECRET_PREFIX must be a secret-name path ending in /")
		}
		applicationCredentials = &appcredentials.SecretsManager{References: references, Prefix: prefix, API: secretsAPI}
	}
	if bundle != nil {
		applicationCredentials = &appcredentials.Bundled{Bundle: bundle, Next: applicationCredentials}
	}
	api := cmp.Or(os.Getenv("KIFF_API"), "https://api.kiff.dev")
	issuer := cmp.Or(os.Getenv("KIFF_ISSUER"), api)
	// KIFF is configured here; vendors are named by agents and servers, so
	// they are reached only at public addresses.
	hc := &http.Client{Timeout: 30 * time.Second}
	vendors := egress.Client(30 * time.Second)
	connectors := map[string]connector.Provider{}
	// Composio is on with its own secret named, or with the one secret
	// (whose composio_api_key it reads first).
	if reference := strings.TrimSpace(os.Getenv("PLANQUE_COMPOSIO_SECRET")); reference != "" || bundle != nil {
		read := bundle.Or("composio_api_key", own(reference, "api_key"))
		connectors["composio"] = &composio.Client{HTTP: vendors, Secret: func(ctx context.Context) (string, error) {
			v, err := read(ctx)
			if err != nil {
				return "", connector.ErrUnavailable
			}
			return v, nil
		}}
	}
	var doorSink door.Sink
	doorURL, doorSecret := strings.TrimSpace(os.Getenv("PLANQUE_DOOR_URL")), strings.TrimSpace(os.Getenv("PLANQUE_DOOR_SECRET"))
	if doorURL != "" || doorSecret != "" {
		u, err := url.Parse(doorURL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (doorSecret == "" && bundle == nil) {
			return nil, errors.New("PLANQUE_DOOR_URL and PLANQUE_DOOR_SECRET (or PLANQUE_SECRET) must name the fixed HTTPS receiver and its secret")
		}
		reader := bundle.Or("door_secret", own(doorSecret, "secret"))
		doorSink = &door.HTTP{URL: doorURL, Secret: reader, Client: egress.Client(10 * time.Second)}
	}
	var eventSink door.EventSink
	if eventURL := strings.TrimSpace(os.Getenv("PLANQUE_EVENT_URL")); eventURL != "" {
		u, err := url.Parse(eventURL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (doorSecret == "" && bundle == nil) {
			return nil, errors.New("PLANQUE_EVENT_URL must name the fixed HTTPS receiver, and PLANQUE_DOOR_SECRET (or PLANQUE_SECRET) the secret it's signed with")
		}
		reader := bundle.Or("door_secret", own(doorSecret, "secret"))
		eventSink = &door.HTTP{URL: eventURL, Secret: reader, Client: egress.Client(10 * time.Second)}
	}
	// The refund events a profile's case block names, for KIFF: its own
	// fixed HTTPS receiver and its own signing secret.
	var kiffEventSink door.EventSink
	kiffURL, kiffSecret := strings.TrimSpace(os.Getenv("PLANQUE_KIFF_EVENT_URL")), strings.TrimSpace(os.Getenv("PLANQUE_KIFF_EVENT_SECRET"))
	if kiffURL != "" || kiffSecret != "" {
		u, err := url.Parse(kiffURL)
		if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (kiffSecret == "" && bundle == nil) {
			return nil, errors.New("PLANQUE_KIFF_EVENT_URL and PLANQUE_KIFF_EVENT_SECRET (or PLANQUE_SECRET) must name KIFF's fixed HTTPS receiver and its secret")
		}
		reader := bundle.Or("kiff_event_secret", own(kiffSecret, "secret"))
		kiffEventSink = &door.HTTP{URL: kiffURL, Secret: reader, Client: egress.Client(10 * time.Second)}
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
	// The tafs Chrome extension's released builds: only these pair and
	// connect, and with none the owner's own Chrome isn't offered.
	var chromeBuilds map[string]planque.ChromeBuild
	if path := os.Getenv("PLANQUE_CHROME_BUILDS_FILE"); path != "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.New("PLANQUE_CHROME_BUILDS_FILE: " + err.Error())
		}
		chromeBuilds, err = planque.ChromeBuildsFrom(f)
		f.Close()
		if err != nil {
			return nil, err
		}
	}
	// The relay to the owner's own Chrome: planque's step key in KMS, and the
	// WebSocket API's management endpoint. Half a setting is a mistake.
	var chromeSigner stepkey.Signer
	var chromePoster wsapi.Poster
	stepKey, wsEndpoint := strings.TrimSpace(os.Getenv("PLANQUE_CHROME_STEP_KEY")), strings.TrimSpace(os.Getenv("PLANQUE_CHROME_WS_ENDPOINT"))
	if stepKey != "" || wsEndpoint != "" {
		if stepKey == "" || !strings.HasPrefix(wsEndpoint, "https://") {
			return nil, errors.New("PLANQUE_CHROME_STEP_KEY and PLANQUE_CHROME_WS_ENDPOINT must name the step key and the WebSocket API's https endpoint")
		}
		chromeSigner = stepkey.KMS{KeyID: stepKey, Client: &awsapi.Client{Service: "kms", Region: region,
			Target: "TrentService", JSONType: "application/x-amz-json-1.1", Creds: creds}}
		chromePoster = wsapi.API{Endpoint: strings.TrimRight(wsEndpoint, "/"), Region: region, Creds: creds}
	}
	// Which accounts the owner's Chrome is offered to: every one, a few, or
	// none (the default), whatever builds are listed.
	chromeOffer, err := planque.ChromeOfferFrom(os.Getenv("PLANQUE_CHROME_OFFER"), os.Getenv("PLANQUE_CHROME_TENANTS"))
	if err != nil {
		return nil, err
	}
	neverUse, err := planque.NeverUseFrom(os.Getenv("PLANQUE_NEVER_USE"))
	if err != nil {
		return nil, err
	}
	testTenants, err := planque.ManagedTestTenantsFrom(os.Getenv("PLANQUE_MANAGED_TEST_TENANTS"))
	if err != nil {
		return nil, errors.New("PLANQUE_MANAGED_TEST_TENANTS " + err.Error())
	}
	ownTenants, err := planque.ManagedTestTenantsFrom(os.Getenv("PLANQUE_OWN_CLIENT_TENANTS"))
	if err != nil {
		return nil, errors.New("PLANQUE_OWN_CLIENT_TENANTS " + err.Error())
	}
	// KIFF's proof monitor, when set, is the only source of proved marks. Half
	// a setting is a mistake, never a silent fall back to the shipped file. Its
	// read-only token is read from planque's one secret (proof_state_token)
	// first, and from PLANQUE_PROOF_STATE_TOKEN while the secret lacks it, so
	// moving it there takes nothing away.
	proofState, err := proofStateFrom(os.Getenv("PLANQUE_PROOF_STATE_URL"), os.Getenv("PLANQUE_PROOF_STATE_TOKEN"), bundle)
	if err != nil {
		return nil, err
	}
	return &planque.Server{
		ProofState:         proofState,
		ManagedTestTenants: testTenants,
		NeverUse:           neverUse,
		OwnClientTenants:   ownTenants,
		DigestV2:           os.Getenv("PLANQUE_DIGEST") == "v2",
		LiveIndex:          os.Getenv("PLANQUE_LIVE_GSI") == "on",
		Browser:            browserWorker,
		ChromeBuilds:       chromeBuilds,
		ChromeOffer:        chromeOffer,
		ChromeScreenshots:  strings.TrimSpace(os.Getenv("PLANQUE_CHROME_SCREENSHOTS")) == "on",
		// The hand-over's tools, offered to new Chrome connections only: on
		// unless PLANQUE_CHROME_HAND=off (chromeHandOn).
		ChromeHand:   chromeHandOn(os.Getenv("PLANQUE_CHROME_HAND")),
		ChromeSigner: chromeSigner,
		ChromePoster: chromePoster,
		// The wake's VAPID contact: planque's own public address. Its key is planque's own (chrome_vapid.go).
		VAPIDSubject:  base,
		DoorSink:      doorSink,
		EventSink:     eventSink,
		KIFFEventSink: kiffEventSink,
		// The reviewed Stripe refunds profile, off unless turned on.
		StripeRefunds:        os.Getenv("PLANQUE_STRIPE_REFUNDS") == "on",
		StripeRefundsTenants: refundTenants(os.Getenv("PLANQUE_STRIPE_REFUNDS")),
		Connectors:           connectors,
		BaseURL:              base,
		AppReturn:            appReturn,
		SignInReader:         signInReader(os.Getenv("PLANQUE_SIGNIN_READER_SHA256")),
		KV:                   kv,
		Sealer:               sealer,
		KIFF:                 &kiff.Client{API: api, Gateway: cmp.Or(os.Getenv("KIFF_GATEWAY"), "https://mcp.kiff.dev"), HTTP: hc},
		Permits: &permit.Verifier{Issuer: issuer, Keys: &permit.JWKS{
			URL: cmp.Or(os.Getenv("KIFF_JWKS_URL"), issuer+"/.well-known/kiff-permit-keys.json"), HTTP: hc}},
		OAuth:                  &oauth.Client{HTTP: vendors},
		MCP:                    &mcp.Client{HTTP: vendors},
		REST:                   &rest.Client{HTTP: vendors},
		Logger:                 logger,
		WindowRequired:         strings.EqualFold(strings.TrimSpace(os.Getenv("PLANQUE_WINDOW_REQUIRED")), "true"),
		ComputerWindowRequired: strings.EqualFold(strings.TrimSpace(os.Getenv("PLANQUE_COMPUTER_WINDOW_REQUIRED")), "true"),
		CatalogEntries:         catalog,
		CatalogOperators:       operators,
		ServerApps:             serverApps,
		DescriberServers:       describerServers,
		Applications:           applications,
		ApplicationCredentials: applicationCredentials,
	}, nil
}

// refundTenants reads PLANQUE_STRIPE_REFUNDS=tenants:<id>[,<id>]: the
// accounts offered the Stripe refunds profile. Anything else lists none.
// chromeHandOn reads PLANQUE_CHROME_HAND: the hand-over is on by default,
// now that KIFF's Coach names its tools for the person's own step (kiff/agents
// browse_hand words, tafs a-browser-taf@3). Only "off" turns it off, so a
// typo never turns it off by accident and the env can still stop it at once.
func chromeHandOn(v string) bool {
	return strings.TrimSpace(v) != "off"
}

func refundTenants(v string) map[string]bool {
	list, ok := strings.CutPrefix(strings.TrimSpace(v), "tenants:")
	if !ok {
		return nil
	}
	out := map[string]bool{}
	for _, t := range strings.Split(list, ",") {
		if t = strings.TrimSpace(t); t != "" {
			out[t] = true
		}
	}
	return out
}

// signInReader is the sign-in reader key's SHA-256 from its hex, or none when
// it isn't 64 lowercase hex characters.
func signInReader(h string) []byte {
	b, err := hex.DecodeString(h)
	if err != nil || len(b) != 32 || strings.ToLower(h) != h {
		return nil
	}
	return b
}

// proofStateFrom is the reader for KIFF's proof state at rawURL: its token
// from bundle's proof_state_token, else token (PLANQUE_PROOF_STATE_TOKEN).
// No address and no token is no reader; an address with no token anywhere,
// a token with no address, or an address that isn't KIFF's is refused.
func proofStateFrom(rawURL, token string, bundle *secretstore.Bundle) (*planque.ProofState, error) {
	rawURL, token = strings.TrimSpace(rawURL), strings.TrimSpace(token)
	if rawURL == "" && token == "" {
		return nil, nil
	}
	var fallback func(context.Context) (string, error)
	if token != "" {
		fallback = func(context.Context) (string, error) { return token, nil }
	}
	var p *planque.ProofState
	if rawURL != "" && (bundle != nil || fallback != nil) {
		p = planque.NewProofStateFrom(rawURL, bundle.Or("proof_state_token", fallback))
	}
	if p == nil {
		return nil, errors.New("PLANQUE_PROOF_STATE_URL must name KIFF's proof state (https), with its read-only token in PLANQUE_SECRET (proof_state_token) or PLANQUE_PROOF_STATE_TOKEN")
	}
	return p, nil
}
