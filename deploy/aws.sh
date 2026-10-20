#!/usr/bin/env bash
# planque on AWS: one Lambda behind a function URL, one DynamoDB table,
# one KMS key. Idempotent. Never prints account ids or ARNs.
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
#
# FAKEKIFF=1 deploys the KIFF stand-in used by the live check instead,
# named $PLANQUE_NAME-fakekiff, configured from FAKEKIFF_ENV (JSON).
# FAKEVENDOR=1 deploys the fake vendor (an MCP server behind its own OAuth,
# which consents on its own) for the live check and kiff/agents' nightly
# test, named $PLANQUE_NAME-fakevendor; FAKEVENDOR_TOKEN guards /fake/calls.
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
  (cd "$out" && zip -q function.zip bootstrap)
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
    stmts="$logs,{\"Effect\":\"Allow\",\"Action\":[\"dynamodb:GetItem\",\"dynamodb:PutItem\",\"dynamodb:DeleteItem\",\"dynamodb:Query\"],\"Resource\":\"arn:aws:dynamodb:$REGION:$ACCOUNT:table/$TABLE\"},{\"Effect\":\"Allow\",\"Action\":[\"kms:GenerateDataKey\",\"kms:Decrypt\"],\"Resource\":\"$key\"}"
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
}

kmskey() {
  if ! aws kms describe-key --key-id "$ALIAS" --query KeyMetadata.KeyId >/dev/null 2>&1; then
    local id; id=$(aws kms create-key --description "planque envelope encryption" --query KeyMetadata.KeyId)
    aws kms create-alias --alias-name "$ALIAS" --target-key-id "$id"
  fi
}

function_env() {
  if [[ "$FAKE" == vendor ]]; then
    python3 -c 'import json,os,sys; print(json.dumps({"Variables": {"FAKEVENDOR_BASE": sys.argv[1], "FAKEVENDOR_TOKEN": os.environ["FAKEVENDOR_TOKEN"], "FAKEVENDOR_TABLE": sys.argv[2], "FAKEVENDOR_REDIRECTS": os.environ.get("FAKEVENDOR_REDIRECTS") or "https://planque.kiff.dev/oauth/callback"}}))' "$1" "$TABLE"
  elif [[ "$FAKE" == kiff ]]; then
    python3 -c 'import json,os,sys; print(json.dumps({"Variables": json.loads(os.environ["FAKEKIFF_ENV"]) | {"FAKEKIFF_ISSUER": sys.argv[1]}}))' "$1"
  else
    python3 -c 'import json,os,sys; v={"PLANQUE_BASE_URL": os.environ.get("PLANQUE_PUBLIC_URL") or sys.argv[1], "PLANQUE_TABLE": sys.argv[2], "PLANQUE_KMS_KEY": sys.argv[3]}
for k in ("KIFF_API","KIFF_GATEWAY","KIFF_ISSUER","KIFF_JWKS_URL","PLANQUE_APP_RETURN"):
    if os.environ.get(k): v[k]=os.environ[k]
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

up() {
  if [[ -z "$FAKE" ]]; then table; kmskey; fi
  if [[ "$FAKE" == vendor ]]; then table; fi
  role
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
  aws lambda update-function-configuration --function-name "$FN" --environment "$(function_env "$url")" \
    --query FunctionName >/dev/null
  aws lambda wait function-updated-v2 --function-name "$FN"
  if [[ -z "$FAKE" ]]; then schedule; fi
  echo "$url"
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
    aws dynamodb delete-table --table-name "$TABLE" --query TableDescription.TableName >/dev/null 2>&1 || true
  fi
  if [[ -z "$FAKE" ]]; then
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
