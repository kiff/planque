#!/usr/bin/env bash
# The browser worker on AWS (kiff/planque#86), next to planque in the same
# region: one Fargate task on a private address, reachable only by planque,
# out to the web only on TCP 443. Idempotent. Never prints account ids, ARNs
# or the key.
#
#   deploy/browserworker/aws.sh            plan: what exists, what would be
#                                          created or changed, the monthly
#                                          cost; reads only, changes nothing
#   deploy/browserworker/aws.sh up         create or update (needs
#                                          BROWSER_WORKER_APPLY=1)
#   deploy/browserworker/aws.sh down       take planque off the worker and
#                                          remove it all (needs
#                                          BROWSER_WORKER_APPLY=1)
#
# The shape:
#   - a VPC with two private subnets (two zones) and one public subnet that
#     holds the NAT gateway; nothing comes in from the internet;
#   - the worker as an ECS Fargate service (arm64, 1 vCPU, 2 GB, one task),
#     no public IP, named worker.<namespace> in a private Cloud Map namespace;
#   - planque's Lambda joins the private subnets: it reaches the worker by
#     that name, and everything else it calls (KIFF, vendors, KMS, Secrets
#     Manager) through the NAT; DynamoDB through a gateway endpoint;
#   - security groups: the worker takes 8787 from planque's group only and
#     sends TCP 443 only. No UDP leaves it: the VPC resolver (its DNS) isn't
#     filtered by security groups, and nothing else is allowed. planque's
#     group sends TCP 443 anywhere and 8787 to the worker.
#   - the key: generated here once (or kept from planque's configuration),
#     held in a variable; planque gets PLANQUE_BROWSER_KEY, the worker only
#     its SHA-256 (PLANQUE_WORKER_KEY_SHA256). planque's other settings are
#     kept as they are.
#
# Environment: AWS_REGION (default eu-west-1), PLANQUE_NAME (default
# planque), BROWSER_VPC_CIDR (default 10.42.0.0/16), BROWSER_NAMESPACE
# (default planque.internal), BROWSER_CPU / BROWSER_MEMORY (default 1024 /
# 2048), BROWSER_IMAGE_TAG (default the commit). Building the image needs
# docker with buildx (linux/arm64).
set -euo pipefail

REGION=${AWS_REGION:-eu-west-1}
NAME=${PLANQUE_NAME:-planque}
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
CIDR=${BROWSER_VPC_CIDR:-10.42.0.0/16}
NS=${BROWSER_NAMESPACE:-planque.internal}
CPU=${BROWSER_CPU:-1024}
MEM=${BROWSER_MEMORY:-2048}
W=$NAME-browser
FN=$NAME
LAMBDA_ROLE=$NAME-lambda
EXEC_ROLE=$W-exec
LOGS=/ecs/$W
PORT=8787
MODE=${1:-plan}
aws() ( command aws --region "$REGION" --output text "$@" )

case "$MODE" in
  plan) ;;
  up|down)
    [[ "${BROWSER_WORKER_APPLY:-}" == 1 ]] || { echo "$MODE changes AWS: set BROWSER_WORKER_APPLY=1 once the owner has said go" >&2; exit 2; } ;;
  *) echo "usage: $0 [plan|up|down]" >&2; exit 2 ;;
esac

ACCOUNT=$(aws sts get-caller-identity --query Account)
BASE=${CIDR%/*}; BASE=${BASE%.*.*}
ZONES=("${REGION}a" "${REGION}b")

say() { echo "$*" >&2; }
# step names a change: plan prints it, up and down make it.
step() { if [[ "$MODE" == plan ]]; then say "  would: $1"; return 1; fi; say "  $1"; }
have() { say "  exists: $1"; }

tagged() { # the id of a resource of kind $1 tagged Name=$2, or nothing
  local q
  case "$1" in
    vpc) q=Vpcs[0].VpcId ;; subnet) q=Subnets[0].SubnetId ;;
    internet-gateway) q=InternetGateways[0].InternetGatewayId ;;
    route-table) q=RouteTables[0].RouteTableId ;; nat-gateway) q=NatGateways[0].NatGatewayId ;;
    security-group) q=SecurityGroups[0].GroupId ;; address) q=Addresses[0].AllocationId ;;
  esac
  local filters=(Name=tag:Name,Values="$2")
  [[ "$1" == nat-gateway ]] && filters+=(Name=state,Values=pending,available)
  local v; v=$(aws ec2 "describe-${1}s" --filters "${filters[@]}" --query "$q" 2>/dev/null || true)
  [[ "$v" == None ]] && v=""
  echo "$v"
}
nametag() { echo "ResourceType=$1,Tags=[{Key=Name,Value=$2},{Key=planque,Value=browser}]"; }

# --- the network ---------------------------------------------------------

VPC="" PUB="" PRIV_A="" PRIV_B="" IGW="" NAT="" RT_PUB="" RT_PRIV="" SG_W="" SG_P=""

network() {
  VPC=$(tagged vpc "$W")
  if [[ -n "$VPC" ]]; then have "VPC $W"; elif step "create VPC $W ($CIDR)"; then
    VPC=$(aws ec2 create-vpc --cidr-block "$CIDR" --tag-specifications "$(nametag vpc "$W")" --query Vpc.VpcId)
    aws ec2 wait vpc-available --vpc-ids "$VPC"
    aws ec2 modify-vpc-attribute --vpc-id "$VPC" --enable-dns-support '{"Value":true}'
    aws ec2 modify-vpc-attribute --vpc-id "$VPC" --enable-dns-hostnames '{"Value":true}'
  fi
  subnet PUB "$W-public" "$BASE.0.0/24" "${ZONES[0]}"
  subnet PRIV_A "$W-private-a" "$BASE.1.0/24" "${ZONES[0]}"
  subnet PRIV_B "$W-private-b" "$BASE.2.0/24" "${ZONES[1]}"

  IGW=$(tagged internet-gateway "$W")
  if [[ -n "$IGW" ]]; then have "internet gateway (for the NAT only)"; elif step "create the internet gateway, for the NAT only"; then
    IGW=$(aws ec2 create-internet-gateway --tag-specifications "$(nametag internet-gateway "$W")" --query InternetGateway.InternetGatewayId)
    aws ec2 attach-internet-gateway --internet-gateway-id "$IGW" --vpc-id "$VPC"
  fi
  RT_PUB=$(tagged route-table "$W-public")
  if [[ -n "$RT_PUB" ]]; then have "public route table"; elif step "route the public subnet to the internet gateway"; then
    RT_PUB=$(aws ec2 create-route-table --vpc-id "$VPC" --tag-specifications "$(nametag route-table "$W-public")" --query RouteTable.RouteTableId)
    aws ec2 create-route --route-table-id "$RT_PUB" --destination-cidr-block 0.0.0.0/0 --gateway-id "$IGW" >/dev/null
    aws ec2 associate-route-table --route-table-id "$RT_PUB" --subnet-id "$PUB" >/dev/null
  fi

  NAT=$(tagged nat-gateway "$W")
  if [[ -n "$NAT" ]]; then have "NAT gateway"; elif step "create one NAT gateway (and its Elastic IP) in the public subnet"; then
    local eip; eip=$(tagged address "$W")
    [[ -n "$eip" ]] || eip=$(aws ec2 allocate-address --domain vpc --tag-specifications "$(nametag elastic-ip "$W")" --query AllocationId)
    NAT=$(aws ec2 create-nat-gateway --subnet-id "$PUB" --allocation-id "$eip" --tag-specifications "$(nametag natgateway "$W")" --query NatGateway.NatGatewayId)
    aws ec2 wait nat-gateway-available --nat-gateway-ids "$NAT"
  fi
  RT_PRIV=$(tagged route-table "$W-private")
  if [[ -n "$RT_PRIV" ]]; then have "private route table"; elif step "route the private subnets out through the NAT, DynamoDB through a gateway endpoint"; then
    RT_PRIV=$(aws ec2 create-route-table --vpc-id "$VPC" --tag-specifications "$(nametag route-table "$W-private")" --query RouteTable.RouteTableId)
    aws ec2 create-route --route-table-id "$RT_PRIV" --destination-cidr-block 0.0.0.0/0 --nat-gateway-id "$NAT" >/dev/null
    aws ec2 associate-route-table --route-table-id "$RT_PRIV" --subnet-id "$PRIV_A" >/dev/null
    aws ec2 associate-route-table --route-table-id "$RT_PRIV" --subnet-id "$PRIV_B" >/dev/null
    aws ec2 create-vpc-endpoint --vpc-id "$VPC" --vpc-endpoint-type Gateway --service-name "com.amazonaws.$REGION.dynamodb" \
      --route-table-ids "$RT_PRIV" --tag-specifications "$(nametag vpc-endpoint "$W-dynamodb")" --query VpcEndpoint.State >/dev/null
  fi
  groups
}

subnet() { # subnet VAR name cidr zone
  local id; id=$(tagged subnet "$2")
  if [[ -n "$id" ]]; then have "subnet $2"; elif step "create subnet $2 ($3, $4)"; then
    id=$(aws ec2 create-subnet --vpc-id "$VPC" --cidr-block "$3" --availability-zone "$4" \
      --tag-specifications "$(nametag subnet "$2")" --query Subnet.SubnetId)
  fi
  printf -v "$1" %s "$id"
}

# groups makes both security groups and sets their rules exactly: each
# group's default allow-all egress is removed, so what's listed is all.
groups() {
  SG_W=$(tagged security-group "$W-worker")
  SG_P=$(tagged security-group "$W-planque")
  if [[ -n "$SG_W" && -n "$SG_P" ]]; then have "security groups (rules reset to the listed ones)"; fi
  if ! step "security groups: worker in 8787 from planque only, out TCP 443 only; planque out TCP 443 and 8787 to the worker"; then return 0; fi
  [[ -n "$SG_W" ]] || SG_W=$(aws ec2 create-security-group --vpc-id "$VPC" --group-name "$W-worker" \
    --description "planque browser worker" --tag-specifications "$(nametag security-group "$W-worker")" --query GroupId)
  [[ -n "$SG_P" ]] || SG_P=$(aws ec2 create-security-group --vpc-id "$VPC" --group-name "$W-planque" \
    --description "planque Lambda, browser side" --tag-specifications "$(nametag security-group "$W-planque")" --query GroupId)
  local g
  for g in "$SG_W" "$SG_P"; do
    local rules; rules=$(command aws --region "$REGION" ec2 describe-security-group-rules --filters Name=group-id,Values="$g" --output json)
    local out in
    out=$(python3 -c 'import json,sys; print(" ".join(r["SecurityGroupRuleId"] for r in json.load(sys.stdin)["SecurityGroupRules"] if r["IsEgress"]))' <<<"$rules")
    in=$(python3 -c 'import json,sys; print(" ".join(r["SecurityGroupRuleId"] for r in json.load(sys.stdin)["SecurityGroupRules"] if not r["IsEgress"]))' <<<"$rules")
    # shellcheck disable=SC2086
    [[ -z "$out" ]] || aws ec2 revoke-security-group-egress --group-id "$g" --security-group-rule-ids $out >/dev/null
    # shellcheck disable=SC2086
    [[ -z "$in" ]] || aws ec2 revoke-security-group-ingress --group-id "$g" --security-group-rule-ids $in >/dev/null
  done
  aws ec2 authorize-security-group-ingress --group-id "$SG_W" \
    --ip-permissions "IpProtocol=tcp,FromPort=$PORT,ToPort=$PORT,UserIdGroupPairs=[{GroupId=$SG_P}]" >/dev/null
  aws ec2 authorize-security-group-egress --group-id "$SG_W" \
    --ip-permissions "IpProtocol=tcp,FromPort=443,ToPort=443,IpRanges=[{CidrIp=0.0.0.0/0}]" >/dev/null
  aws ec2 authorize-security-group-egress --group-id "$SG_P" \
    --ip-permissions "IpProtocol=tcp,FromPort=443,ToPort=443,IpRanges=[{CidrIp=0.0.0.0/0}]" \
    "IpProtocol=tcp,FromPort=$PORT,ToPort=$PORT,UserIdGroupPairs=[{GroupId=$SG_W}]" >/dev/null
}

# --- the worker ----------------------------------------------------------

IMAGE=""

image() {
  local repo=$W tag=${BROWSER_IMAGE_TAG:-$(git -C "$ROOT" rev-parse --short=12 HEAD)}
  local host=$ACCOUNT.dkr.ecr.$REGION.amazonaws.com
  IMAGE=$host/$repo:$tag
  if aws ecr describe-repositories --repository-names "$repo" --query 'repositories[0].repositoryName' >/dev/null 2>&1; then
    have "image repository $repo"
  elif step "create image repository $repo (scan on push, tags immutable)"; then
    aws ecr create-repository --repository-name "$repo" --image-scanning-configuration scanOnPush=true \
      --image-tag-mutability IMMUTABLE --query repository.repositoryName >/dev/null
    aws ecr put-lifecycle-policy --repository-name "$repo" --lifecycle-policy-text \
      '{"rules":[{"rulePriority":1,"selection":{"tagStatus":"any","countType":"imageCountMoreThan","countNumber":5},"action":{"type":"expire"}}]}' >/dev/null
  fi
  if aws ecr describe-images --repository-name "$repo" --image-ids imageTag="$tag" --query 'imageDetails[0].imageTags[0]' >/dev/null 2>&1; then
    have "image $repo:$tag"
  elif step "build deploy/browserworker/Dockerfile for linux/arm64 and push it as $repo:$tag"; then
    aws ecr get-login-password | docker login --username AWS --password-stdin "$host" >/dev/null
    docker buildx build --platform linux/arm64 -f "$ROOT/deploy/browserworker/Dockerfile" -t "$IMAGE" --push "$ROOT" >&2
  fi
}

KEY="" KEY_SUM=""

# key keeps planque's key when it has one, otherwise makes one. Never printed.
key() {
  KEY=$(command aws --region "$REGION" lambda get-function-configuration --function-name "$FN" --output json \
    | python3 -c 'import json,sys; print((json.load(sys.stdin).get("Environment") or {}).get("Variables",{}).get("PLANQUE_BROWSER_KEY",""))')
  if [[ -n "$KEY" ]]; then have "planque's browser key (kept)"; elif step "make planque's browser key (32 random bytes)"; then
    KEY=$(openssl rand -hex 32)
  fi
  if [[ -n "$KEY" ]]; then KEY_SUM=$(printf %s "$KEY" | shasum -a 256 | cut -d' ' -f1); fi
}

NS_ID="" SVC_REG=""

discovery() {
  NS_ID=$(aws servicediscovery list-namespaces --filters Name=NAME,Values="$NS",Condition=EQ --query 'Namespaces[0].Id' 2>/dev/null || true)
  [[ "$NS_ID" == None ]] && NS_ID=""
  if [[ -n "$NS_ID" ]]; then have "private namespace $NS"; elif step "create the private DNS namespace $NS in the VPC"; then
    local op; op=$(aws servicediscovery create-private-dns-namespace --name "$NS" --vpc "$VPC" --query OperationId)
    until [[ "$(aws servicediscovery get-operation --operation-id "$op" --query Operation.Status)" == SUCCESS ]]; do sleep 5; done
    NS_ID=$(aws servicediscovery list-namespaces --filters Name=NAME,Values="$NS",Condition=EQ --query 'Namespaces[0].Id')
  fi
  SVC_REG=$(aws servicediscovery list-services --filters Name=NAMESPACE_ID,Values="${NS_ID:-none}",Condition=EQ \
    --query "Services[?Name=='worker']|[0].Arn" 2>/dev/null || true)
  [[ "$SVC_REG" == None ]] && SVC_REG=""
  if [[ -n "$SVC_REG" ]]; then have "name worker.$NS"; elif step "register the name worker.$NS (A record, 10 s)"; then
    SVC_REG=$(aws servicediscovery create-service --name worker --namespace-id "$NS_ID" \
      --dns-config 'RoutingPolicy=MULTIVALUE,DnsRecords=[{Type=A,TTL=10}]' \
      --health-check-custom-config FailureThreshold=1 --query Service.Arn)
  fi
}

# The image has no curl: bash asks /healthz itself.
HEALTH="exec 3<>/dev/tcp/127.0.0.1/$PORT && printf 'GET /healthz HTTP/1.0\\r\\n\\r\\n' >&3 && head -n1 <&3 | grep -q ' 200 '"

worker() {
  if aws iam get-role --role-name "$EXEC_ROLE" --query Role.RoleName >/dev/null 2>&1; then
    have "task execution role (pull the image, write its logs)"
  elif step "create the task execution role (pull the image, write its logs; nothing else)"; then
    aws iam create-role --role-name "$EXEC_ROLE" --query Role.RoleName --assume-role-policy-document \
      '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ecs-tasks.amazonaws.com"},"Action":"sts:AssumeRole"}]}' >/dev/null
    aws iam attach-role-policy --role-name "$EXEC_ROLE" --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
    sleep 10 # IAM is eventually consistent
  fi
  if step "keep the worker's logs 30 days in $LOGS"; then
    aws logs create-log-group --log-group-name "$LOGS" >/dev/null 2>&1 || true
    aws logs put-retention-policy --log-group-name "$LOGS" --retention-in-days 30
  fi
  if [[ "$(aws ecs describe-clusters --clusters "$W" --query 'clusters[0].status' 2>/dev/null || true)" == ACTIVE ]]; then
    have "ECS cluster $W"
  elif step "create the ECS cluster $W"; then
    aws ecs create-cluster --cluster-name "$W" --query cluster.status >/dev/null
  fi
  local td=""
  if step "register the task: arm64, $CPU CPU units, $MEM MB, no task role (it calls nothing in AWS), PLANQUE_WORKER_KEY_SHA256"; then
    td=$(python3 -c 'import json,sys
img,logs,region,cpu,mem,port,ks,role,fam,health=sys.argv[1:]
print(json.dumps({"family":fam,"requiresCompatibilities":["FARGATE"],"networkMode":"awsvpc","cpu":cpu,"memory":mem,
 "runtimePlatform":{"cpuArchitecture":"ARM64","operatingSystemFamily":"LINUX"},"executionRoleArn":role,
 "containerDefinitions":[{"name":"worker","image":img,"essential":True,"stopTimeout":30,
  "portMappings":[{"containerPort":int(port),"protocol":"tcp"}],
  "environment":[{"name":"PLANQUE_WORKER_KEY_SHA256","value":ks}],
  "healthCheck":{"command":["CMD","bash","-c",health],"interval":30,"timeout":5,"retries":3,"startPeriod":20},
  "logConfiguration":{"logDriver":"awslogs","options":{"awslogs-group":logs,"awslogs-region":region,"awslogs-stream-prefix":"worker"}}}]}))' \
      "$IMAGE" "$LOGS" "$REGION" "$CPU" "$MEM" "$PORT" "$KEY_SUM" \
      "$(aws iam get-role --role-name "$EXEC_ROLE" --query Role.Arn)" "$W" "$HEALTH")
    td=$(aws ecs register-task-definition --cli-input-json "$td" --query taskDefinition.taskDefinitionArn)
  fi
  local net="awsvpcConfiguration={subnets=[$PRIV_A,$PRIV_B],securityGroups=[$SG_W],assignPublicIp=DISABLED}"
  if [[ "$(aws ecs describe-services --cluster "$W" --services worker --query 'services[0].status' 2>/dev/null || true)" == ACTIVE ]]; then
    have "ECS service worker"
    if step "roll the service onto the new task"; then
      aws ecs update-service --cluster "$W" --service worker --task-definition "$td" --network-configuration "$net" --query service.status >/dev/null
    fi
  elif step "create the ECS service worker: one task, private subnets, no public IP, named worker.$NS"; then
    aws ecs create-service --cluster "$W" --service-name worker --task-definition "$td" --desired-count 1 \
      --launch-type FARGATE --platform-version LATEST --network-configuration "$net" \
      --service-registries "registryArn=$SVC_REG" --query service.status >/dev/null
  fi
  if [[ "$MODE" == up ]]; then aws ecs wait services-stable --cluster "$W" --services worker; fi
}

# --- planque -------------------------------------------------------------

planque() {
  if ! step "planque's Lambda: may manage its network interfaces, joins the private subnets with group $W-planque, gets PLANQUE_BROWSER_URL=http://worker.$NS:$PORT and PLANQUE_BROWSER_KEY (other settings kept)"; then return 0; fi
  aws iam attach-role-policy --role-name "$LAMBDA_ROLE" --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole
  local env; env=$(command aws --region "$REGION" lambda get-function-configuration --function-name "$FN" --output json \
    | PLANQUE_BROWSER_KEY="$KEY" python3 -c 'import json,os,sys
v=(json.load(sys.stdin).get("Environment") or {}).get("Variables",{})
v["PLANQUE_BROWSER_URL"]=sys.argv[1]; v["PLANQUE_BROWSER_KEY"]=os.environ["PLANQUE_BROWSER_KEY"]
print(json.dumps({"Variables":v}))' "http://worker.$NS:$PORT")
  aws lambda update-function-configuration --function-name "$FN" --environment "$env" \
    --vpc-config "SubnetIds=$PRIV_A,$PRIV_B,SecurityGroupIds=$SG_P" --query FunctionName >/dev/null
  aws lambda wait function-updated-v2 --function-name "$FN"
}

cost() {
  python3 - "$CPU" "$MEM" <<'EOF' >&2
import sys
cpu, mem = int(sys.argv[1]) / 1024, int(sys.argv[2]) / 1024
h = 730  # hours in a month
# eu-west-1 on-demand list prices; check them before quoting.
fargate = (cpu * 0.03238 + mem * 0.00356) * h        # Fargate, Linux/ARM
nat = 0.048 * h                                        # one NAT gateway
ip = 0.005 * h                                         # its public IPv4 address
dns = 0.50 + 1.00                                      # private hosted zone, Cloud Map
ecr = 0.10 * 1.0                                       # about 1 GB of images
print("\nmonthly, before traffic (USD, eu-west-1, on-demand):")
for k, v in (("Fargate worker, always on", fargate), ("NAT gateway", nat), ("NAT public IPv4", ip),
             ("private DNS, Cloud Map", dns), ("image storage", ecr)):
    print(f"  {k:<28} {v:7.2f}")
print(f"  {'total':<28} {fargate + nat + ip + dns + ecr:7.2f}")
print("plus per GB: NAT processing 0.048 (what planque and the browser fetch), logs 0.57 ingested")
print("DynamoDB goes through its free gateway endpoint, not the NAT.")
EOF
}

# --- down ----------------------------------------------------------------

down() {
  say "taking planque off the worker first"
  local env; env=$(command aws --region "$REGION" lambda get-function-configuration --function-name "$FN" --output json \
    | python3 -c 'import json,sys
v=(json.load(sys.stdin).get("Environment") or {}).get("Variables",{})
v.pop("PLANQUE_BROWSER_URL",None); v.pop("PLANQUE_BROWSER_KEY",None)
print(json.dumps({"Variables":v}))')
  aws lambda update-function-configuration --function-name "$FN" --environment "$env" \
    --vpc-config SubnetIds=[],SecurityGroupIds=[] --query FunctionName >/dev/null
  aws lambda wait function-updated-v2 --function-name "$FN"
  aws iam detach-role-policy --role-name "$LAMBDA_ROLE" --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole 2>/dev/null || true

  aws ecs update-service --cluster "$W" --service worker --desired-count 0 --query service.status >/dev/null 2>&1 || true
  aws ecs delete-service --cluster "$W" --service worker --force --query service.status >/dev/null 2>&1 || true
  aws ecs wait services-inactive --cluster "$W" --services worker 2>/dev/null || true
  aws ecs delete-cluster --cluster "$W" --query cluster.status >/dev/null 2>&1 || true
  local td
  for td in $(aws ecs list-task-definitions --family-prefix "$W" --query taskDefinitionArns 2>/dev/null || true); do
    aws ecs deregister-task-definition --task-definition "$td" --query taskDefinition.status >/dev/null || true
  done
  NS_ID=$(aws servicediscovery list-namespaces --filters Name=NAME,Values="$NS",Condition=EQ --query 'Namespaces[0].Id' 2>/dev/null || true)
  if [[ -n "$NS_ID" && "$NS_ID" != None ]]; then
    local s; for s in $(aws servicediscovery list-services --filters Name=NAMESPACE_ID,Values="$NS_ID",Condition=EQ --query 'Services[].Id'); do
      aws servicediscovery delete-service --id "$s" || true
    done
    aws servicediscovery delete-namespace --id "$NS_ID" --query OperationId >/dev/null || true
  fi
  aws ecr delete-repository --repository-name "$W" --force --query repository.repositoryName >/dev/null 2>&1 || true
  command aws --region "$REGION" logs delete-log-group --log-group-name "$LOGS" 2>/dev/null || true
  aws iam detach-role-policy --role-name "$EXEC_ROLE" --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy 2>/dev/null || true
  aws iam delete-role --role-name "$EXEC_ROLE" 2>/dev/null || true

  VPC=$(tagged vpc "$W")
  if [[ -n "$VPC" ]]; then
    local e
    for e in $(aws ec2 describe-vpc-endpoints --filters Name=vpc-id,Values="$VPC" --query 'VpcEndpoints[].VpcEndpointId'); do
      aws ec2 delete-vpc-endpoints --vpc-endpoint-ids "$e" >/dev/null || true
    done
    NAT=$(tagged nat-gateway "$W")
    if [[ -n "$NAT" ]]; then aws ec2 delete-nat-gateway --nat-gateway-id "$NAT" >/dev/null; aws ec2 wait nat-gateway-deleted --nat-gateway-ids "$NAT"; fi
    e=$(tagged address "$W"); [[ -z "$e" ]] || aws ec2 release-address --allocation-id "$e"
    # Lambda's interfaces linger a while after it leaves the VPC.
    local g; for g in "$W-worker" "$W-planque"; do
      local id; id=$(tagged security-group "$g")
      [[ -z "$id" ]] && continue
      for _ in $(seq 1 60); do aws ec2 delete-security-group --group-id "$id" 2>/dev/null && break; sleep 30; done
    done
    for g in "$W-public" "$W-private"; do
      local rt; rt=$(tagged route-table "$g"); [[ -z "$rt" ]] && continue
      local a; for a in $(aws ec2 describe-route-tables --route-table-ids "$rt" --query 'RouteTables[0].Associations[].RouteTableAssociationId'); do
        aws ec2 disassociate-route-table --association-id "$a"
      done
      aws ec2 delete-route-table --route-table-id "$rt"
    done
    IGW=$(tagged internet-gateway "$W")
    if [[ -n "$IGW" ]]; then aws ec2 detach-internet-gateway --internet-gateway-id "$IGW" --vpc-id "$VPC"; aws ec2 delete-internet-gateway --internet-gateway-id "$IGW"; fi
    for g in "$W-public" "$W-private-a" "$W-private-b"; do
      local sn; sn=$(tagged subnet "$g"); [[ -z "$sn" ]] || aws ec2 delete-subnet --subnet-id "$sn"
    done
    aws ec2 delete-vpc --vpc-id "$VPC"
  fi
  say "removed the browser worker; planque runs without a browser"
}

if [[ "$MODE" == down ]]; then down; exit 0; fi
say "browser worker, $REGION ($MODE):"
network
image
key
discovery
worker
planque
cost
if [[ "$MODE" == up ]]; then say "worker.$NS:$PORT is up; planque offers the browser"; fi
