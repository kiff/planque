package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

// Run the actual deployment functions against a synthetic AWS command. This
// verifies the generated policy/environment without touching an AWS account.
func TestDeploymentScopesManagedConnectorSecret(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is unavailable")
	}
	raw, err := os.ReadFile("../../deploy/aws.sh")
	if err != nil {
		t.Fatal(err)
	}
	extract := func(name string) string {
		r := regexp.MustCompile(`(?ms)^` + name + `\(\) \{.*?^\}`)
		v := r.FindString(string(raw))
		if v == "" {
			t.Fatal("deployment function missing: " + name)
		}
		return v
	}
	for _, scenario := range []string{"enabled", "disabled", "fake", "wrong_region"} {
		t.Run(scenario, func(t *testing.T) {
			script := `set -euo pipefail
REGION=eu-west-1 ACCOUNT=000000000000 ROLE=demo-lambda FN=demo TABLE=demo ALIAS=alias/demo FAKE=""
export PLANQUE_APP_RETURN=https://app.example/confirm
export PLANQUE_COMPOSIO_SECRET=demo/connectors/composio
connector_arn=arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/connectors/composio-AbCdEf
aws() {
 case "$1 $2" in
 "iam get-role") return 0;;
 "kms describe-key") echo arn:aws:kms:eu-west-1:000000000000:key/demo;;
 "secretsmanager describe-secret") echo "$connector_arn";;
 "iam put-role-policy")
  while [[ $# -gt 0 ]]; do
   if [[ "$1" == --policy-document ]]; then echo "$2";return;fi
   shift
  done
  return 1;;
 *) echo "unexpected deployment operation" >&2; return 1;;
 esac
}
`
			switch scenario {
			case "disabled":
				script += "unset PLANQUE_COMPOSIO_SECRET\n"
			case "fake":
				script += "FAKE=kiff\nexport FAKEKIFF_ENV='{}'\n"
			case "wrong_region":
				script += "connector_arn=arn:aws:secretsmanager:us-east-1:000000000000:secret:other\n"
			}
			script += extract("role") + "\n" + extract("function_env") + "\nrole\nfunction_env https://lambda.example\n"
			out, err := exec.Command(bash, "-c", script).CombinedOutput()
			if scenario == "wrong_region" {
				if err == nil || strings.Contains(string(out), `"Statement"`) {
					t.Fatal("an unresolved or foreign secret was allowed")
				}
				return
			}
			if err != nil {
				t.Fatalf("deployment configuration failed: %v", err)
			}
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(lines) != 2 {
				t.Fatal("unexpected deployment output")
			}
			var policy struct {
				Statement []struct {
					Action   []string
					Resource string
				}
			}
			if json.Unmarshal([]byte(lines[0]), &policy) != nil {
				t.Fatal("invalid IAM policy")
			}
			granted := 0
			for _, stmt := range policy.Statement {
				for _, action := range stmt.Action {
					if strings.HasPrefix(action, "secretsmanager:") {
						granted++
						if action != "secretsmanager:GetSecretValue" || stmt.Resource != "arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/connectors/composio-AbCdEf" {
							t.Fatal("secret permission was broader than the exact configured resource")
						}
					}
				}
			}
			var env struct{ Variables map[string]string }
			if json.Unmarshal([]byte(lines[1]), &env) != nil {
				t.Fatal("invalid Lambda environment")
			}
			if scenario == "enabled" {
				if granted != 1 || env.Variables["PLANQUE_COMPOSIO_SECRET"] != "demo/connectors/composio" {
					t.Fatal("managed connector was not configured")
				}
			} else if granted != 0 || env.Variables["PLANQUE_COMPOSIO_SECRET"] != "" {
				t.Fatal("secret access leaked into an unconfigured or fake deployment")
			}
		})
	}
}
