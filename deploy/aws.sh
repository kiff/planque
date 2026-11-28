#!/usr/bin/env bash
# planque on AWS: one Lambda behind a function URL, one DynamoDB table,
# one KMS key. Idempotent. Never prints account ids or ARNs. The table keeps
# point-in-time recovery on. The key is rotated yearly, and its policy lets
# only planque's role seal and open with it (keypolicy).
#
#   deploy/aws.sh up      create or update; prints the function URL
#   deploy/aws.sh down    remove everything (the KMS key is scheduled
#                         for deletion after 7 days, AWS's minimum)
#
# Environment: AWS_REGION (default eu-west-1), PLANQUE_NAME (default
# planque), and for the function: PLANQUE_APP_RETURN (required: the owner's
# app, which confirms each sign-in), KIFF_API, KIFF_GATEWAY, KIFF_ISSUER.
#
# PLANQUE_PUBLIC_URL is planque's public name when it has one in front of
# the function URL (e.g. https://planque.kiff.dev through CloudFront); the
# signatures, the relay audience and the vendor redirect all use it.
# PLANQUE_COMPOSIO_SECRET optionally enables managed connections. It names an
# existing Secrets Manager secret; deployment resolves its exact ARN for the
# runtime policy and never reads or prints its value. PLANQUE_DOOR_SECRET
# similarly names the {secret} handoff credential; PLANQUE_DOOR_URL is the
# fixed HTTPS event receiver, and PLANQUE_EVENT_URL the fixed HTTPS receiver
# of app events, signed with the same handoff credential. Each configured
# secret gets its own exact resource. PLANQUE_KIFF_EVENT_URL is KIFF's fixed
# HTTPS receiver of the refund events a profile's case block names, and
# PLANQUE_KIFF_EVENT_SECRET names the Secrets Manager secret they're signed
# with: both or neither, and the function may read that one secret only.
# Without them, no event goes to KIFF.
# PLANQUE_SECRET names planque's one secret, a JSON of composio_api_key,
# door_secret, kiff_event_secret, proof_state_token and apps ({id: {client_id,
# client_secret}}),
# read first; a value it lacks is still read from its own secret above, so
# those go only once it's proven live. It gets its own exact resource too, is
# kept on a rerun like the rest (kept_settings, resolved once for the checks,
# the grants and the function), and stands in for PLANQUE_KIFF_EVENT_SECRET
# beside PLANQUE_KIFF_EVENT_URL.
#
# deploy/server-apps.json, the reviewed list of which app a vendor's own
# server is, ships with the function and is read at start
# (PLANQUE_SERVER_APPS_FILE).
#
# deploy/apps.json, the reviewed definitions of the apps planque signs in to
# with KIFF's own OAuth client (PLANQUE_APPS_FILE), ships the same way. Each
# app's client is the Secrets Manager secret planque/apps/<its id>
# (PLANQUE_APP_SECRET_PREFIX), and the function may read secrets under that
# prefix alone: a new app needs its definition row and its secret, nothing
# here.
#
# deploy/catalog.json, the reviewed vendor catalog (connections an agent asks
# for by entry, such as Qonto by key), ships with it too and is read at start
# (PLANQUE_CATALOG_FILE). Each key entry pins its description's SHA-256; the
# descriptions it names live in deploy/specs.
#
# PLANQUE_CHROME_RELAY=on sets up the relay to the owner's own Chrome: an
# asymmetric KMS key that signs each step (alias/$PLANQUE_NAME-chrome-steps,
# ECC P-256, sign and verify only; planque's role may sign with it and read
# its public half, nothing else), and a WebSocket API ($PLANQUE_NAME-chrome,
# stage v1) whose $connect, $disconnect and $default routes invoke the
# function, which may post to and drop its connections only. The function
# gets PLANQUE_CHROME_STEP_KEY and PLANQUE_CHROME_WS_ENDPOINT, and up prints
# the wss:// address the extension connects to. PLANQUE_CHROME_RELAY=off
# drops both settings (the relay refuses every connect) and keeps the key and
# the API; unset, a rerun keeps what the function has.
#
# Which accounts the owner's Chrome is offered to is its own setting:
# PLANQUE_CHROME_OFFER=on (every account), test (only the accounts in
# PLANQUE_CHROME_TENANTS, tenants:<id>[,<id>]) or off. Unset, it's offered to
# none, whatever deploy/chrome-builds.txt lists. A rerun keeps both as the
# function has them.
#
# Which accounts planque's hosted browser is offered to is the same:
# PLANQUE_BROWSER_OFFER=on, test (only PLANQUE_BROWSER_TENANTS) or off. Unset,
# it's offered to none, whatever PLANQUE_BROWSER_URL says. A rerun keeps both.
#
# FAKEKIFF=1 deploys the KIFF stand-in used by the live check instead,
# named $PLANQUE_NAME-fakekiff, configured from FAKEKIFF_ENV (JSON).
# FAKEVENDOR=1 deploys the fake vendor (an MCP server behind its own OAuth,
# which consents on its own) for the live check and kiff/agents' nightly
# test, named $PLANQUE_NAME-fakevendor; FAKEVENDOR_TOKEN guards /fake/calls.
# FAKEVENDOR_DESK_KEY, the desk's and billing's key, is kept from the
# function when not given.
# Its state is in its own table of that name, so any instance answers.
# FAKEVENDOR_REDIRECTS (default https://planque.kiff.dev/oauth/callback) is
# the only redirect a client may register: it consents on its own.
set -euo pipefail

REGION=${AWS_REGION:-eu-west-1}
NAME=${PLANQUE_NAME:-planque}
ROOT=$(cd "$(dirname "$0")/.." && pwd)
aws() ( command aws --region "$REGION" --output text "$@" )

ACCOUNT=$(aws sts get-caller-identity --query Account)
TABLE=$NAME
ALIAS=alias/$NAME
STEP_ALIAS=alias/$NAME-chrome-steps
WSAPI=$NAME-chrome
WS_STAGE=v1
RELAY=${PLANQUE_CHROME_RELAY:-}
[[ -z "$RELAY" || "$RELAY" == on || "$RELAY" == off ]] || { echo "PLANQUE_CHROME_RELAY is on or off" >&2; exit 2; }
ROLE=$NAME-lambda
FN=$NAME
PKG=./cmd/planque
FAKE=""
if [[ "${FAKEKIFF:-}" == 1 ]]; then FAKE=kiff; fi
if [[ "${FAKEVENDOR:-}" == 1 ]]; then FAKE=vendor; fi
if [[ "$FAKE" == vendor ]]; then
  FN=$NAME-fakevendor ROLE=$NAME-fakevendor-lambda PKG=./internal/fakes/cmd/fakevendor TABLE=$NAME-fakevendor
fi
if [[ "$FAKE" == kiff ]]; then
  FN=$NAME-fakekiff ROLE=$NAME-fakekiff-lambda PKG=./internal/fakes/cmd/fakekiff
fi

build() {
  local out; out=$(mktemp -d)
  (cd "$ROOT" && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags=-s -o "$out/bootstrap" "$PKG")
  if [[ -z "$FAKE" ]]; then
    # The reviewed list of which app a vendor's own server is, read at start.
    cp "$ROOT/deploy/server-apps.json" "$out/server-apps.json"
    # The reviewed vendor catalog, read at start too.
    cp "$ROOT/deploy/catalog.json" "$out/catalog.json"
    cp "$ROOT/deploy/chrome-builds.txt" "$out/chrome-builds.txt"
    # The reviewed apps planque signs in to with KIFF's own client.
    cp "$ROOT/deploy/apps.json" "$out/apps.json"
    (cd "$out" && zip -q function.zip bootstrap server-apps.json catalog.json apps.json chrome-builds.txt)
  else
    (cd "$out" && zip -q function.zip bootstrap)
  fi
  echo "$out/function.zip"
}

role() {
  if ! aws iam get-role --role-name "$ROLE" --query Role.RoleName >/dev/null 2>&1; then
    aws iam create-role --role-name "$ROLE" --query Role.RoleName --assume-role-policy-document \
      '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' >/dev/null
    sleep 10 # IAM is eventually consistent
  fi
  local logs="{\"Effect\":\"Allow\",\"Action\":[\"logs:CreateLogGroup\",\"logs:CreateLogStream\",\"logs:PutLogEvents\"],\"Resource\":\"arn:aws:logs:$REGION:$ACCOUNT:log-group:/aws/lambda/$FN:*\"}"
  local stmts="$logs"
  if [[ "$FAKE" == vendor ]]; then
    stmts="$logs,{\"Effect\":\"Allow\",\"Action\":[\"dynamodb:GetItem\",\"dynamodb:PutItem\",\"dynamodb:DeleteItem\",\"dynamodb:Query\"],\"Resource\":\"arn:aws:dynamodb:$REGION:$ACCOUNT:table/$TABLE\"}"
  fi
  if [[ -z "$FAKE" ]]; then
    local key; key=$(aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.Arn)
    stmts="$logs,{\"Effect\":\"Allow\",\"Action\":[\"dynamodb:GetItem\",\"dynamodb:PutItem\",\"dynamodb:DeleteItem\",\"dynamodb:Query\",\"dynamodb:Scan\"],\"Resource\":\"arn:aws:dynamodb:$REGION:$ACCOUNT:table/$TABLE\"},{\"Effect\":\"Allow\",\"Action\":[\"kms:GenerateDataKey\",\"kms:Decrypt\"],\"Resource\":\"$key\"}"
    # The live index: read only (store.QueryLive), the table's own index alone.
    stmts="$stmts,{\"Effect\":\"Allow\",\"Action\":[\"dynamodb:Query\"],\"Resource\":\"arn:aws:dynamodb:$REGION:$ACCOUNT:table/$TABLE/index/live\"}"
    # KIFF's own OAuth clients, one secret per app, under planque/apps/ alone.
    stmts="$stmts,{\"Effect\":\"Allow\",\"Action\":[\"secretsmanager:GetSecretValue\"],\"Resource\":\"arn:aws:secretsmanager:$REGION:$ACCOUNT:secret:planque/apps/*\"}"
    # The Chrome relay: sign steps with its own key, and post to or drop the
    # WebSocket API's connections, nothing else on either.
    if [[ "${RELAY:-}" == on ]]; then
      local step_key api; step_key=$(aws kms describe-key --key-id "$STEP_ALIAS" --query KeyMetadata.Arn); api=$(wsapi_id)
      stmts="$stmts,{\"Effect\":\"Allow\",\"Action\":[\"kms:Sign\",\"kms:GetPublicKey\"],\"Resource\":\"$step_key\"}"
      stmts="$stmts,{\"Effect\":\"Allow\",\"Action\":[\"execute-api:ManageConnections\"],\"Resource\":\"arn:aws:execute-api:$REGION:$ACCOUNT:$api/$WS_STAGE/*/@connections/*\"}"
    fi
    # KIFF's receiver and its signing secret go together, as planque checks
    # at start: half of the pair is refused here, before anything changes.
    if [[ -n "${PLANQUE_KIFF_EVENT_URL:-}" && -z "${PLANQUE_KIFF_EVENT_SECRET:-}" && -z "${PLANQUE_SECRET:-}" ]] || [[ -z "${PLANQUE_KIFF_EVENT_URL:-}" && -n "${PLANQUE_KIFF_EVENT_SECRET:-}" ]]; then
      echo "PLANQUE_KIFF_EVENT_URL and PLANQUE_KIFF_EVENT_SECRET go together: set both, or neither" >&2
      return 1
    fi
    local secret_reference
    for secret_reference in "${PLANQUE_COMPOSIO_SECRET:-}" "${PLANQUE_DOOR_SECRET:-}" "${PLANQUE_KIFF_EVENT_SECRET:-}" "${PLANQUE_SECRET:-}"; do
      [[ -z "$secret_reference" ]] && continue
      local connector_secret
      connector_secret=$(aws secretsmanager describe-secret --secret-id "$secret_reference" --query ARN)
      [[ "$connector_secret" == arn:aws:secretsmanager:"$REGION":"$ACCOUNT":secret:* ]] || { echo "the configured connector secret could not be resolved" >&2; return 1; }
      stmts="$stmts,{\"Effect\":\"Allow\",\"Action\":[\"secretsmanager:GetSecretValue\"],\"Resource\":\"$connector_secret\"}"
    done
  fi
  aws iam put-role-policy --role-name "$ROLE" --policy-name "$ROLE" \
    --policy-document "{\"Version\":\"2012-10-17\",\"Statement\":[$stmts]}"
}

table() {
  if ! aws dynamodb describe-table --table-name "$TABLE" --query Table.TableName >/dev/null 2>&1; then
    aws dynamodb create-table --table-name "$TABLE" --billing-mode PAY_PER_REQUEST \
      --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S \
      --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE \
      --sse-specification Enabled=true --query TableDescription.TableName >/dev/null
    aws dynamodb wait table-exists --table-name "$TABLE"
    aws dynamodb update-time-to-live --table-name "$TABLE" \
      --time-to-live-specification Enabled=true,AttributeName=ttl --query TimeToLiveSpecification.Enabled >/dev/null
  fi
  # Deletion protection, checked on every run, so a table rebuilt or changed
  # by hand gets it back: planque's records can't be dropped by one command.
  # The fake vendor's test table goes without.
  if [[ -z "$FAKE" && "$(aws dynamodb describe-table --table-name "$TABLE" --query Table.DeletionProtectionEnabled)" != "True" ]]; then
    aws dynamodb update-table --table-name "$TABLE" --deletion-protection-enabled --query TableDescription.TableName >/dev/null
  fi
  # The live index (store.LiveIndex): a grant record carries "live" only while
  # it counts, so the index holds live grants alone, kept by DynamoDB in the
  # same write. Added once; reads use it only with PLANQUE_LIVE_GSI=on, set
  # after it's active.
  if [[ -z "$FAKE" ]] && [[ "$(aws dynamodb describe-table --table-name "$TABLE" --query "length(Table.GlobalSecondaryIndexes[?IndexName=='live'] || \`[]\`)")" == 0 ]]; then
    aws dynamodb update-table --table-name "$TABLE" --attribute-definitions AttributeName=live,AttributeType=S \
      --global-secondary-index-updates '[{"Create":{"IndexName":"live","KeySchema":[{"AttributeName":"live","KeyType":"HASH"}],"Projection":{"ProjectionType":"ALL"}}}]' \
      --query TableDescription.TableName >/dev/null
  fi
  # Point-in-time recovery, asked for on every run, so a table rebuilt or
  # changed by hand gets it back: without it, a bad write can't be undone.
  # A table created seconds ago may refuse it for a short while
  # (ContinuousBackupsUnavailableException), so it's asked again.
  local try
  for try in 1 2 3 4 5 6; do
    if aws dynamodb update-continuous-backups --table-name "$TABLE" \
      --point-in-time-recovery-specification PointInTimeRecoveryEnabled=true \
      --query ContinuousBackupsDescription.ContinuousBackupsStatus >/dev/null; then
      return 0
    fi
    [[ $try == 6 ]] && return 1
    sleep "${PITR_RETRY_SECONDS:-10}"
  done
}

kmskey() {
  if ! aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.KeyId >/dev/null 2>&1; then
    local id; id=$(aws kms create-key --description "planque envelope encryption" --query KeyMetadata.KeyId)
    aws kms create-alias --alias-name "$ALIAS" --target-key-id "$id"
  fi
  # Rotated yearly by KMS: old versions stay for opening what they sealed, so
  # nothing is sealed again.
  aws kms enable-key-rotation --key-id "$(aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.KeyId)"
}

# keypolicy: only planque's role seals and opens with the key. The account's
# administrators administer it (describe, rotate, disable, schedule deletion,
# change this policy) but can't use it, and no grant can be made on it, so no
# IAM policy, however broad, opens a stored secret. Changing this policy is
# the break-glass, and CloudTrail records it. Run after role(): the policy
# names the role, and KMS refuses a principal that doesn't exist yet.
keypolicy() {
  local id role_arn
  id=$(aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.KeyId)
  role_arn=$(aws iam get-role --role-name "$ROLE" --query Role.Arn)
  local admin='"kms:DescribeKey","kms:GetKeyPolicy","kms:PutKeyPolicy","kms:ListKeyPolicies","kms:GetKeyRotationStatus","kms:EnableKeyRotation","kms:DisableKeyRotation","kms:ListKeyRotations","kms:RotateKeyOnDemand","kms:EnableKey","kms:DisableKey","kms:ListAliases","kms:CreateAlias","kms:UpdateAlias","kms:DeleteAlias","kms:UpdateKeyDescription","kms:ListGrants","kms:RevokeGrant","kms:RetireGrant","kms:ListResourceTags","kms:TagResource","kms:UntagResource","kms:ScheduleKeyDeletion","kms:CancelKeyDeletion"'
  local policy="{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"AdministerNeverUse\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"arn:aws:iam::$ACCOUNT:root\"},\"Action\":[$admin],\"Resource\":\"*\"},{\"Sid\":\"PlanqueSealsAndOpens\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"$role_arn\"},\"Action\":[\"kms:GenerateDataKey\",\"kms:Decrypt\"],\"Resource\":\"*\"}]}"
  aws kms put-key-policy --key-id "$id" --policy-name default --policy "$policy"
}

# stepkey makes the relay's signing key: ECC P-256, sign and verify only. An
# asymmetric key can't rotate on its own; a new key is a new pairing of trust
# (the extension reads planque's public half at each welcome).
stepkey() {
  if ! aws kms describe-key --key-id "$STEP_ALIAS" --query KeyMetadata.KeyId >/dev/null 2>&1; then
    local id; id=$(aws kms create-key --description "planque Chrome step signing" --key-spec ECC_NIST_P256 \
      --key-usage SIGN_VERIFY --query KeyMetadata.KeyId)
    aws kms create-alias --alias-name "$STEP_ALIAS" --target-key-id "$id"
  fi
}

# stepkeypolicy: as keypolicy, the account's administrators administer the
# step key but can't sign with it; only planque's role signs and reads its
# public half. Run after role().
stepkeypolicy() {
  local id role_arn
  id=$(aws kms describe-key --key-id "$STEP_ALIAS" --query KeyMetadata.KeyId)
  role_arn=$(aws iam get-role --role-name "$ROLE" --query Role.Arn)
  local admin='"kms:DescribeKey","kms:GetKeyPolicy","kms:PutKeyPolicy","kms:ListKeyPolicies","kms:EnableKey","kms:DisableKey","kms:ListAliases","kms:CreateAlias","kms:UpdateAlias","kms:DeleteAlias","kms:UpdateKeyDescription","kms:ListGrants","kms:RevokeGrant","kms:RetireGrant","kms:ListResourceTags","kms:TagResource","kms:UntagResource","kms:ScheduleKeyDeletion","kms:CancelKeyDeletion"'
  local policy="{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"AdministerNeverSign\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"arn:aws:iam::$ACCOUNT:root\"},\"Action\":[$admin],\"Resource\":\"*\"},{\"Sid\":\"PlanqueSignsSteps\",\"Effect\":\"Allow\",\"Principal\":{\"AWS\":\"$role_arn\"},\"Action\":[\"kms:Sign\",\"kms:GetPublicKey\"],\"Resource\":\"*\"}]}"
  aws kms put-key-policy --key-id "$id" --policy-name default --policy "$policy"
}

# wsapi_id is the relay's WebSocket API's id, empty when there's none.
wsapi_id() {
  local id; id=$(aws apigatewayv2 get-apis --query "Items[?Name=='$WSAPI'].ApiId | [0]")
  [[ "$id" == None ]] && id=""
  echo "$id"
}

# wsapi makes the relay's WebSocket API. Every message goes to $default: the
# function reads the message's own type.
wsapi() {
  if [[ -z "$(wsapi_id)" ]]; then
    aws apigatewayv2 create-api --name "$WSAPI" --protocol-type WEBSOCKET \
      --route-selection-expression '$request.body.action' --query ApiId >/dev/null
  fi
}

# wsapi_wire points the API's three routes at the function, lets the API (and
# only it) invoke the function, and deploys stage v1. Run once the function
# exists.
wsapi_wire() {
  local id fn_arn integ
  id=$(wsapi_id)
  fn_arn=$(aws lambda get-function --function-name "$FN" --query Configuration.FunctionArn)
  integ=$(aws apigatewayv2 get-integrations --api-id "$id" --query 'Items[0].IntegrationId')
  if [[ -z "$integ" || "$integ" == None ]]; then
    integ=$(aws apigatewayv2 create-integration --api-id "$id" --integration-type AWS_PROXY --integration-method POST \
      --integration-uri "arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/$fn_arn/invocations" --query IntegrationId)
  fi
  local rk
  for rk in '$connect' '$disconnect' '$default'; do
    if [[ "$(aws apigatewayv2 get-routes --api-id "$id" --query "Items[?RouteKey=='$rk'].RouteId | [0]")" == None ]]; then
      aws apigatewayv2 create-route --api-id "$id" --route-key "$rk" --target "integrations/$integ" --query RouteId >/dev/null
    fi
  done
  aws lambda add-permission --function-name "$FN" --statement-id chrome-ws --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com --source-arn "arn:aws:execute-api:$REGION:$ACCOUNT:$id/*" >/dev/null 2>&1 || true
  # A few Chromes per account, each a message every few seconds at most.
  local throttle="ThrottlingBurstLimit=100,ThrottlingRateLimit=50"
  if ! aws apigatewayv2 get-stage --api-id "$id" --stage-name "$WS_STAGE" --query StageName >/dev/null 2>&1; then
    aws apigatewayv2 create-stage --api-id "$id" --stage-name "$WS_STAGE" --default-route-settings "$throttle" --query StageName >/dev/null
  else
    aws apigatewayv2 update-stage --api-id "$id" --stage-name "$WS_STAGE" --default-route-settings "$throttle" --query StageName >/dev/null
  fi
  aws apigatewayv2 create-deployment --api-id "$id" --stage-name "$WS_STAGE" --query DeploymentId >/dev/null
}

function_env() {
  if [[ "$FAKE" == vendor ]]; then
    # The desk's and billing's key (FAKEVENDOR_DESK_KEY) is set on the
    # function by hand and kept as it is, unless given here.
    { aws lambda get-function-configuration --function-name "$FN" --query Environment.Variables --output json 2>/dev/null || echo '{}'; } |
    python3 -c 'import json,os,sys; v={"FAKEVENDOR_BASE": sys.argv[1], "FAKEVENDOR_TOKEN": os.environ["FAKEVENDOR_TOKEN"], "FAKEVENDOR_TABLE": sys.argv[2], "FAKEVENDOR_REDIRECTS": os.environ.get("FAKEVENDOR_REDIRECTS") or "https://planque.kiff.dev/oauth/callback"}
now=json.load(sys.stdin) or {}
key=os.environ.get("FAKEVENDOR_DESK_KEY") or now.get("FAKEVENDOR_DESK_KEY")
if key: v["FAKEVENDOR_DESK_KEY"]=key
print(json.dumps({"Variables": v}))' "$1" "$TABLE"
  elif [[ "$FAKE" == kiff ]]; then
    python3 -c 'import json,os,sys; print(json.dumps({"Variables": json.loads(os.environ["FAKEKIFF_ENV"]) | {"FAKEKIFF_ISSUER": sys.argv[1]}}))' "$1"
  else
    # The browser worker's address and key are set by
    # deploy/browserworker/aws.sh and kept here as they are; so is which
    # accounts are offered Stripe refunds (PLANQUE_STRIPE_REFUNDS), which
    # test accounts may make a managed connection to any app
    # (PLANQUE_MANAGED_TEST_TENANTS, from the deployer, never this repo), which
    # accounts are shown an app's own entry in place of the managed one
    # (PLANQUE_OWN_CLIENT_TENANTS, the same), which accounts a robot test
    # user runs, with no live Stripe key and no browse_pay
    # (PLANQUE_PROOF_TEST_TENANTS, the same), which accounts the owner's
    # Chrome is offered to (PLANQUE_CHROME_OFFER: on, test or off, and
    # PLANQUE_CHROME_TENANTS for test; lost, it's off for every account),
    # which accounts the hosted browser is offered to (PLANQUE_BROWSER_OFFER
    # and PLANQUE_BROWSER_TENANTS, the same way), the
    # grants whose sign-ins are never used (PLANQUE_NEVER_USE, from the
    # deployer, never this repo), the
    # catalog's operators, unless PLANQUE_CATALOG_OPERATORS is given, and
    # the servers that get Stripe's describing read with their sign-in
    # besides Stripe's own (PLANQUE_DESCRIBER_SERVERS): live planque sets
    # it for the Eval's vendor, a host we run, so the Eval measures the
    # path a Stripe customer gets; a rerun keeps it.
    # planque's switches (PLANQUE_DIGEST=v2, PLANQUE_WINDOW_REQUIRED=true,
    # PLANQUE_COMPUTER_WINDOW_REQUIRED=true)
    # are taken from the environment when given, else kept from the
    # function, so a rerun never turns one off; set one to another value
    # ("v1", "false") to turn it off. cmd/planque checks every PLANQUE_
    # setting the code reads is named here.
    { aws lambda get-function-configuration --function-name "$FN" --query Environment.Variables --output json 2>/dev/null || echo '{}'; } |
    python3 -c 'import json,os,sys; v={"PLANQUE_BASE_URL": os.environ.get("PLANQUE_PUBLIC_URL") or sys.argv[1], "PLANQUE_TABLE": sys.argv[2], "PLANQUE_KMS_KEY": sys.argv[3], "PLANQUE_SERVER_APPS_FILE": "/var/task/server-apps.json", "PLANQUE_CATALOG_FILE": "/var/task/catalog.json", "PLANQUE_APPS_FILE": "/var/task/apps.json", "PLANQUE_CHROME_BUILDS_FILE": "/var/task/chrome-builds.txt", "PLANQUE_APP_SECRET_PREFIX": "planque/apps/"}
for k in ("KIFF_API","KIFF_GATEWAY","KIFF_ISSUER","KIFF_JWKS_URL","PLANQUE_APP_RETURN","PLANQUE_COMPOSIO_SECRET","PLANQUE_DOOR_URL","PLANQUE_DOOR_SECRET","PLANQUE_EVENT_URL","PLANQUE_KIFF_EVENT_URL","PLANQUE_KIFF_EVENT_SECRET","PLANQUE_SECRET","PLANQUE_SIGNIN_READER_SHA256","PLANQUE_CATALOG_OPERATORS","PLANQUE_DESCRIBER_SERVERS","PLANQUE_MANAGED_TEST_TENANTS","PLANQUE_PROOF_TEST_TENANTS","PLANQUE_OWN_CLIENT_TENANTS","PLANQUE_PROOF_STATE_URL","PLANQUE_PROOF_STATE_TOKEN","PLANQUE_CHROME_STEP_KEY","PLANQUE_CHROME_WS_ENDPOINT","PLANQUE_CHROME_OFFER","PLANQUE_CHROME_TENANTS","PLANQUE_BROWSER_OFFER","PLANQUE_BROWSER_TENANTS","PLANQUE_CHROME_SCREENSHOTS","PLANQUE_CHROME_HAND","PLANQUE_HOSTED_HAND","PLANQUE_SEALED_FILL","PLANQUE_NEVER_USE"):
    if os.environ.get(k): v[k]=os.environ[k]
now=json.load(sys.stdin) or {}
for k in ("PLANQUE_BROWSER_URL","PLANQUE_BROWSER_KEY","PLANQUE_BROWSER_PAY","PLANQUE_BROWSER_DOWNLOADS","PLANQUE_EVIDENCE","PLANQUE_FULL_SCOPES","PLANQUE_ENFORCE_PINS","PLANQUE_BROWSER_SESSIONS","PLANQUE_STRIPE_PROFILE","PLANQUE_STRIPE_REFUNDS","PLANQUE_CATALOG_OPERATORS","PLANQUE_DESCRIBER_SERVERS","PLANQUE_MANAGED_TEST_TENANTS","PLANQUE_PROOF_TEST_TENANTS","PLANQUE_OWN_CLIENT_TENANTS","PLANQUE_PROOF_STATE_URL","PLANQUE_PROOF_STATE_TOKEN","PLANQUE_CHROME_STEP_KEY","PLANQUE_CHROME_WS_ENDPOINT","PLANQUE_CHROME_OFFER","PLANQUE_CHROME_TENANTS","PLANQUE_BROWSER_OFFER","PLANQUE_BROWSER_TENANTS","PLANQUE_CHROME_SCREENSHOTS","PLANQUE_CHROME_HAND","PLANQUE_HOSTED_HAND","PLANQUE_SEALED_FILL","PLANQUE_NEVER_USE"):
    if now.get(k) and not v.get(k): v[k]=now[k]
if os.environ.get("PLANQUE_CHROME_RELAY") == "off":
    v.pop("PLANQUE_CHROME_STEP_KEY", None); v.pop("PLANQUE_CHROME_WS_ENDPOINT", None)
for k in ("PLANQUE_DIGEST","PLANQUE_WINDOW_REQUIRED","PLANQUE_COMPUTER_WINDOW_REQUIRED","PLANQUE_LIVE_GSI","PLANQUE_LIVE_LIST"):
    if os.environ.get(k) or now.get(k): v[k]=os.environ.get(k) or now[k]
print(json.dumps({"Variables": v}))' "$1" "$TABLE" "$ALIAS"
  fi
}

# schedule runs the sweep every hour (EventBridge Scheduler): a sign-in
# nobody confirmed, or whose tools nobody approved, is revoked on time even if
# nobody reads it again. The event is a function URL request, so the handler
# needs nothing new.
schedule() {
  local srole=$NAME-scheduler
  if ! aws iam get-role --role-name "$srole" --query Role.RoleName >/dev/null 2>&1; then
    aws iam create-role --role-name "$srole" --query Role.RoleName --assume-role-policy-document \
      '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"scheduler.amazonaws.com"},"Action":"sts:AssumeRole"}]}' >/dev/null
    sleep 10 # IAM is eventually consistent
  fi
  local fn_arn; fn_arn=$(aws lambda get-function --function-name "$FN" --query Configuration.FunctionArn)
  aws iam put-role-policy --role-name "$srole" --policy-name "$srole" --policy-document \
    "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"lambda:InvokeFunction\",\"Resource\":\"$fn_arn\"}]}"
  local target; target=$(python3 -c 'import json,sys; print(json.dumps({"Arn": sys.argv[1], "RoleArn": sys.argv[2], "Input": json.dumps({"rawPath": "/v1/sweep", "headers": {}, "body": "", "requestContext": {"domainName": "schedule.internal", "http": {"method": "POST", "sourceIp": "127.0.0.1"}}})}))' \
    "$fn_arn" "$(aws iam get-role --role-name "$srole" --query Role.Arn)")
  if aws scheduler get-schedule --name "$NAME-sweep" --query Name >/dev/null 2>&1; then
    aws scheduler update-schedule --name "$NAME-sweep" --schedule-expression "rate(1 hour)" \
      --flexible-time-window Mode=OFF --target "$target" --query ScheduleArn >/dev/null
  else
    aws scheduler create-schedule --name "$NAME-sweep" --schedule-expression "rate(1 hour)" \
      --flexible-time-window Mode=OFF --target "$target" --query ScheduleArn >/dev/null
  fi
}

# A rerun keeps planque's settings as the function has them. Each named here
# is taken from the environment when given, else from the function, so a
# rerun without one never turns a connector, a receiver or its secret off,
# nor moves planque's public address (PLANQUE_PUBLIC_URL, kept from
# PLANQUE_BASE_URL). Set one to "off" to turn it off. The role's grants
# follow, since they're read from the same variables.
kept_settings() {
  # A function that isn't there yet keeps nothing. One that is, but whose
  # settings can't be read, stops the run before anything is written: read
  # as none, the rerun would drop every one of them.
  local now='{}' probe
  if probe=$(aws lambda get-function --function-name "$FN" --query Configuration.FunctionName 2>&1); then
    now=$(aws lambda get-function-configuration --function-name "$FN" --query Environment.Variables --output json) ||
      { echo "the function's settings couldn't be read; nothing was changed" >&2; exit 1; }
  elif [[ "$probe" != *ResourceNotFoundException* ]]; then
    # Only a function that's confirmed missing keeps nothing; any other
    # failure (throttled, denied, offline) stops the run.
    echo "whether the function exists couldn't be read; nothing was changed" >&2; exit 1
  fi
  # The settings go to the parser on stdin, never as an argument: they hold
  # values (a token among them) that a process list would show.
  local exports
  exports=$(printf '%s' "$now" | python3 -c '
import json,os,shlex,sys
now=json.loads(sys.stdin.read() or "{}") or {}
kept={n: now.get(n) for n in ("PLANQUE_APP_RETURN","PLANQUE_COMPOSIO_SECRET","PLANQUE_DOOR_URL","PLANQUE_DOOR_SECRET","PLANQUE_EVENT_URL","PLANQUE_KIFF_EVENT_URL","PLANQUE_KIFF_EVENT_SECRET","PLANQUE_SIGNIN_READER_SHA256","PLANQUE_SECRET")}
kept["PLANQUE_PUBLIC_URL"]=now.get("PLANQUE_BASE_URL")
for n,v in kept.items():
    given=os.environ.get(n)
    if given=="off": print("unset %s" % n)
    elif not given and v: print("export %s=%s" % (n, shlex.quote(v)))
') || { echo "the function's settings couldn't be parsed; nothing was changed" >&2; exit 1; }
  eval "$exports"
}

up() {
  if [[ -n "$FAKE" && -n "$RELAY" ]]; then echo "PLANQUE_CHROME_RELAY is for planque itself, not a fake" >&2; exit 2; fi
  if [[ -z "$FAKE" ]]; then kept_settings; fi
  if [[ -z "$FAKE" ]]; then table; kmskey; fi
  if [[ "$RELAY" == on ]]; then stepkey; wsapi; fi
  if [[ "$FAKE" == vendor ]]; then table; fi
  role
  if [[ -z "$FAKE" ]]; then keypolicy; fi
  if [[ "$RELAY" == on ]]; then stepkeypolicy; fi
  local zip; zip=$(build)
  local role_arn; role_arn=$(aws iam get-role --role-name "$ROLE" --query Role.Arn)
  if aws lambda get-function --function-name "$FN" --query Configuration.FunctionName >/dev/null 2>&1; then
    aws lambda update-function-code --function-name "$FN" --zip-file "fileb://$zip" --query FunctionName >/dev/null
    aws lambda wait function-updated-v2 --function-name "$FN"
  else
    for i in 1 2 3 4 5; do # a new role takes a moment to be assumable
      if aws lambda create-function --function-name "$FN" --runtime provided.al2023 --architectures arm64 \
        --handler bootstrap --role "$role_arn" --zip-file "fileb://$zip" --timeout 30 --memory-size 256 \
        --query FunctionName >/dev/null 2>&1; then break; fi
      sleep 5
    done
    aws lambda wait function-active-v2 --function-name "$FN"
  fi
  if [[ "$FAKE" == kiff ]]; then
    # Its state is in memory: one instance.
    aws lambda put-function-concurrency --function-name "$FN" --reserved-concurrent-executions 1 >/dev/null
  elif [[ "$FAKE" == vendor ]]; then
    aws lambda delete-function-concurrency --function-name "$FN" >/dev/null
  fi
  if ! aws lambda get-function-url-config --function-name "$FN" --query FunctionUrl >/dev/null 2>&1; then
    aws lambda create-function-url-config --function-name "$FN" --auth-type NONE --query FunctionUrl >/dev/null
    aws lambda add-permission --function-name "$FN" --statement-id url --action lambda:InvokeFunctionUrl \
      --principal '*' --function-url-auth-type NONE >/dev/null
    aws lambda add-permission --function-name "$FN" --statement-id url-invoke --action lambda:InvokeFunction \
      --principal '*' --invoked-via-function-url >/dev/null 2>&1 || true
  fi
  local url; url=$(aws lambda get-function-url-config --function-name "$FN" --query FunctionUrl)
  url=${url%/}
  local ws=""
  if [[ "$RELAY" == on ]]; then
    wsapi_wire
    ws="$(wsapi_id).execute-api.$REGION.amazonaws.com/$WS_STAGE"
    export PLANQUE_CHROME_STEP_KEY="$STEP_ALIAS" PLANQUE_CHROME_WS_ENDPOINT="https://$ws"
  fi
  aws lambda update-function-configuration --function-name "$FN" --environment "$(function_env "$url")" \
    --query FunctionName >/dev/null
  aws lambda wait function-updated-v2 --function-name "$FN"
  # Its logs are kept 30 days, as KIFF's privacy page and DPA (Annex 4) say
  # a log group has no limit unless it is given one.
  aws logs create-log-group --log-group-name "/aws/lambda/$FN" >/dev/null 2>&1 || true
  aws logs put-retention-policy --log-group-name "/aws/lambda/$FN" --retention-in-days 30
  if [[ -z "$FAKE" ]]; then schedule; fi
  echo "$url"
  if [[ -n "$ws" ]]; then echo "wss://$ws"; fi
}

down() {
  if [[ -z "$FAKE" ]]; then
    aws scheduler delete-schedule --name "$NAME-sweep" 2>/dev/null || true
    aws iam delete-role-policy --role-name "$NAME-scheduler" --policy-name "$NAME-scheduler" 2>/dev/null || true
    aws iam delete-role --role-name "$NAME-scheduler" 2>/dev/null || true
  fi
  aws lambda delete-function --function-name "$FN" >/dev/null 2>&1 || true
  aws iam delete-role-policy --role-name "$ROLE" --policy-name "$ROLE" 2>/dev/null || true
  aws iam delete-role --role-name "$ROLE" 2>/dev/null || true
  command aws --region "$REGION" logs delete-log-group --log-group-name "/aws/lambda/$FN" 2>/dev/null || true
  if [[ -z "$FAKE" || "$FAKE" == vendor ]]; then
    # planque's own table has deletion protection: it stays, unless a person
    # turns that off first (update-table --no-deletion-protection-enabled).
    aws dynamodb delete-table --table-name "$TABLE" --query TableDescription.TableName >/dev/null 2>&1 ||
      echo "kept table $TABLE (deletion protection is on, or it's gone already)" >&2
  fi
  if [[ -z "$FAKE" ]]; then
    local api; api=$(wsapi_id 2>/dev/null || true)
    if [[ -n "$api" ]]; then aws apigatewayv2 delete-api --api-id "$api"; fi
    local step; step=$(aws kms describe-key --key-id "$STEP_ALIAS" --query KeyMetadata.KeyId 2>/dev/null || true)
    if [[ -n "$step" ]]; then
      aws kms delete-alias --alias-name "$STEP_ALIAS"
      aws kms schedule-key-deletion --key-id "$step" --pending-window-in-days 7 --query KeyState >/dev/null
    fi
    local id; id=$(aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.KeyId 2>/dev/null || true)
    if [[ -n "$id" ]]; then
      aws kms delete-alias --alias-name "$ALIAS"
      aws kms schedule-key-deletion --key-id "$id" --pending-window-in-days 7 --query KeyState >/dev/null
    fi
  fi
  echo "removed $FN"
}

case "${1:-}" in
  up) up ;;
  down) down ;;
  *) echo "usage: $0 up|down" >&2; exit 2 ;;
esac
