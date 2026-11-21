package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
	for _, scenario := range []string{"enabled", "disabled", "fake", "wrong_region", "handoff", "browser", "kiff_events", "kiff_events_half", "switches_kept", "switches_given", "relay", "relay_kept", "relay_off"} {
		t.Run(scenario, func(t *testing.T) {
			script := `set -euo pipefail
REGION=eu-west-1 ACCOUNT=000000000000 ROLE=demo-lambda FN=demo TABLE=demo ALIAS=alias/demo FAKE=""
unset PLANQUE_DOOR_SECRET PLANQUE_DOOR_URL PLANQUE_KIFF_EVENT_URL PLANQUE_KIFF_EVENT_SECRET PLANQUE_DIGEST PLANQUE_WINDOW_REQUIRED
export PLANQUE_APP_RETURN=https://app.example/confirm
export PLANQUE_COMPOSIO_SECRET=demo/connectors/composio
connector_arn=arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/connectors/composio-AbCdEf
current_env=null
aws() {
 case "$1 $2" in
 "iam get-role") return 0;;
 "lambda get-function-configuration") echo "$current_env";;
 "apigatewayv2 get-apis") echo abc123;;
 "kms describe-key") echo arn:aws:kms:eu-west-1:000000000000:key/demo;;
 "secretsmanager describe-secret")
  while [[ $# -gt 0 ]]; do
   if [[ "$1" == --secret-id ]]; then
    if [[ "$2" == demo/doors/handoff ]]; then echo arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/doors/handoff-ZyXwVu
    elif [[ "$2" == demo/kiff/events ]]; then echo arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/kiff/events-QrStUv
    else echo "$connector_arn"; fi
    return
   fi
   shift
  done
  return 1;;
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
			case "handoff":
				script += "export PLANQUE_DOOR_SECRET=demo/doors/handoff\nexport PLANQUE_DOOR_URL=https://receiver.example/v1/doors/planque\n"
			case "browser":
				script += `current_env='{"PLANQUE_BROWSER_URL":"http://worker.demo.internal:8787","PLANQUE_BROWSER_KEY":"k","PLANQUE_STALE":"x"}'` + "\n"
			case "kiff_events":
				script += "export PLANQUE_KIFF_EVENT_URL=https://api.example/v1/planque/refund-events PLANQUE_KIFF_EVENT_SECRET=demo/kiff/events\n"
			case "kiff_events_half":
				script += "export PLANQUE_KIFF_EVENT_URL=https://api.example/v1/planque/refund-events\n"
			case "switches_kept":
				script += `current_env='{"PLANQUE_DIGEST":"v2","PLANQUE_WINDOW_REQUIRED":"true"}'` + "\n"
			case "switches_given":
				script += `current_env='{"PLANQUE_DIGEST":"v2","PLANQUE_WINDOW_REQUIRED":"true"}'` + "\nexport PLANQUE_DIGEST=v1\n"
			case "relay":
				// As up() exports them once the key and the API are there.
				script += "RELAY=on\nexport PLANQUE_CHROME_RELAY=on PLANQUE_CHROME_STEP_KEY=alias/demo-chrome-steps PLANQUE_CHROME_WS_ENDPOINT=https://abc123.execute-api.eu-west-1.amazonaws.com/v1\nWS_STAGE=v1 STEP_ALIAS=alias/demo-chrome-steps WSAPI=demo-chrome\n" + extract("wsapi_id") + "\n"
			case "relay_kept", "relay_off":
				script += `current_env='{"PLANQUE_CHROME_STEP_KEY":"alias/demo-chrome-steps","PLANQUE_CHROME_WS_ENDPOINT":"https://abc123.execute-api.eu-west-1.amazonaws.com/v1"}'` + "\n"
				if scenario == "relay_off" {
					script += "RELAY=off\nexport PLANQUE_CHROME_RELAY=off\n"
				}
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
			// KIFF's receiver without its secret, or the other way round, is
			// refused before any policy is written.
			if scenario == "kiff_events_half" {
				if err == nil || strings.Contains(string(out), `"Statement"`) || !strings.Contains(string(out), "go together") {
					t.Fatalf("half of KIFF's event settings was deployed: %s", out)
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
			granted, apps, signs, posts := 0, 0, 0, 0
			for _, stmt := range policy.Statement {
				for _, action := range stmt.Action {
					// The Chrome relay: sign with its own key, post to and drop
					// the API's connections in stage v1, nothing wider.
					switch {
					case action == "kms:Sign" || action == "kms:GetPublicKey":
						if stmt.Resource != "arn:aws:kms:eu-west-1:000000000000:key/demo" || len(stmt.Action) != 2 {
							t.Fatalf("the step key's grant is wider than the key: %v", stmt)
						}
						signs++
					case strings.HasPrefix(action, "execute-api:"):
						if action != "execute-api:ManageConnections" || stmt.Resource != "arn:aws:execute-api:eu-west-1:000000000000:abc123/v1/*/@connections/*" {
							t.Fatalf("the WebSocket API's grant is wider than its connections: %v", stmt)
						}
						posts++
					}
					if strings.HasPrefix(action, "secretsmanager:") {
						// KIFF's own OAuth clients: one secret per app, under
						// planque/apps/ alone, never wider; counted apart.
						if action == "secretsmanager:GetSecretValue" && stmt.Resource == "arn:aws:secretsmanager:eu-west-1:000000000000:secret:planque/apps/*" {
							apps++
							continue
						}
						granted++
						if action != "secretsmanager:GetSecretValue" || (stmt.Resource != "arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/connectors/composio-AbCdEf" &&
							!(scenario == "handoff" && stmt.Resource == "arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/doors/handoff-ZyXwVu") &&
							!(scenario == "kiff_events" && stmt.Resource == "arn:aws:secretsmanager:eu-west-1:000000000000:secret:demo/kiff/events-QrStUv")) {
							t.Fatal("secret permission was broader than the exact configured resource")
						}
					}
				}
			}
			var env struct{ Variables map[string]string }
			if json.Unmarshal([]byte(lines[1]), &env) != nil {
				t.Fatal("invalid Lambda environment")
			}
			// The apps' clients are readable, under planque/apps/ alone, exactly
			// when the apps' secret prefix is set; a real deployment sets it.
			if prefix := env.Variables["PLANQUE_APP_SECRET_PREFIX"]; (prefix != "") != (apps == 1) || apps > 1 ||
				prefix != "" && (prefix != "planque/apps/" || env.Variables["PLANQUE_APPS_FILE"] != "/var/task/apps.json") {
				t.Fatalf("the apps' secret prefix %q and %d statements reading it don't go together", prefix, apps)
			}
			if scenario == "enabled" && apps != 1 {
				t.Fatal("a real deployment can't read its apps' clients")
			}
			if scenario == "browser" {
				if env.Variables["PLANQUE_BROWSER_URL"] != "http://worker.demo.internal:8787" || env.Variables["PLANQUE_BROWSER_KEY"] != "k" || env.Variables["PLANQUE_STALE"] != "" {
					t.Fatal("the browser worker's settings weren't kept as they were, or something else was")
				}
			} else if env.Variables["PLANQUE_BROWSER_URL"] != "" || env.Variables["PLANQUE_BROWSER_KEY"] != "" {
				t.Fatal("browser settings appeared from nowhere")
			}
			// A rerun keeps planque's switches as they are; one given wins.
			switch want := [2]string{env.Variables["PLANQUE_DIGEST"], env.Variables["PLANQUE_WINDOW_REQUIRED"]}; scenario {
			case "switches_kept":
				if want != [2]string{"v2", "true"} {
					t.Fatalf("a rerun turned a switch off: %v", want)
				}
			case "switches_given":
				if want != [2]string{"v1", "true"} {
					t.Fatalf("a switch given wasn't taken, or the other wasn't kept: %v", want)
				}
			default:
				if want != [2]string{} {
					t.Fatalf("a switch appeared from nowhere: %v", want)
				}
			}
			// The relay's grants and settings are there exactly when it's on, and a
			// rerun without PLANQUE_CHROME_RELAY keeps its settings.
			relayEnv := [2]string{env.Variables["PLANQUE_CHROME_STEP_KEY"], env.Variables["PLANQUE_CHROME_WS_ENDPOINT"]}
			wantEnv := [2]string{"alias/demo-chrome-steps", "https://abc123.execute-api.eu-west-1.amazonaws.com/v1"}
			switch scenario {
			case "relay":
				if signs != 2 || posts != 1 || relayEnv != wantEnv {
					t.Fatalf("the relay wasn't granted or configured: %d %d %v", signs, posts, relayEnv)
				}
			case "relay_kept":
				if signs != 0 || posts != 0 || relayEnv != wantEnv {
					t.Fatalf("a rerun dropped the relay's settings or granted it: %v", relayEnv)
				}
			default:
				if signs != 0 || posts != 0 || relayEnv != [2]string{} {
					t.Fatalf("the relay appeared from nowhere, or off kept it: %d %d %v", signs, posts, relayEnv)
				}
			}
			if scenario == "kiff_events" {
				// The connector's secret and KIFF's: two exact reads, nothing wider.
				if granted != 2 || env.Variables["PLANQUE_KIFF_EVENT_URL"] != "https://api.example/v1/planque/refund-events" || env.Variables["PLANQUE_KIFF_EVENT_SECRET"] != "demo/kiff/events" {
					t.Fatal("KIFF's event receiver or its one secret was not configured")
				}
			} else if env.Variables["PLANQUE_KIFF_EVENT_URL"] != "" || env.Variables["PLANQUE_KIFF_EVENT_SECRET"] != "" {
				t.Fatal("KIFF's event settings appeared from nowhere")
			}
			if scenario == "kiff_events" {
				// Checked above, with the connector's secret beside KIFF's.
			} else if scenario == "handoff" {
				if granted != 2 || env.Variables["PLANQUE_DOOR_SECRET"] != "demo/doors/handoff" || env.Variables["PLANQUE_DOOR_URL"] != "https://receiver.example/v1/doors/planque" {
					t.Fatal("handoff secret or fixed receiver was not configured")
				}
			} else if scenario == "enabled" || scenario == "browser" || strings.HasPrefix(scenario, "switches_") || strings.HasPrefix(scenario, "relay") {
				if granted != 1 || env.Variables["PLANQUE_COMPOSIO_SECRET"] != "demo/connectors/composio" {
					t.Fatal("managed connector was not configured")
				}
			} else if granted != 0 || env.Variables["PLANQUE_COMPOSIO_SECRET"] != "" {
				t.Fatal("secret access leaked into an unconfigured or fake deployment")
			}
		})
	}
}

// The table keeps point-in-time recovery on, whether this run creates it or
// finds it: a table rebuilt or changed by hand gets it back.
func TestDeploymentKeepsPointInTimeRecovery(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is unavailable")
	}
	raw, err := os.ReadFile("../../deploy/aws.sh")
	if err != nil {
		t.Fatal(err)
	}
	table := regexp.MustCompile(`(?ms)^table\(\) \{.*?^\}`).FindString(string(raw))
	if table == "" {
		t.Fatal("deployment function missing: table")
	}
	// unavailable is how often a new table refuses recovery before it takes
	// it (ContinuousBackupsUnavailableException): a few times, then never.
	runAs := func(exists bool, unavailable int, fake, protected string) (string, error) {
		found := "1"
		if exists {
			found = "0"
		}
		dir := t.TempDir()
		log := dir + "/calls"
		script := `set -euo pipefail
TABLE=demo PITR_RETRY_SECONDS=0 FAKE=` + fake + `
refusals=` + strconv.Itoa(unavailable) + `
aws() {
 echo "$*" >>` + log + `
 case "$1 $2" in
 "dynamodb describe-table")
  if [[ "$*" == *DeletionProtectionEnabled* ]]; then echo ` + protected + `; return 0; fi
  return ` + found + `;;
 "dynamodb create-table"|"dynamodb wait"|"dynamodb update-time-to-live"|"dynamodb update-table") return 0;;
 "dynamodb update-continuous-backups")
  n=$(cat ` + dir + `/refused 2>/dev/null || echo 0)
  if (( n < refusals )); then echo $((n+1)) >` + dir + `/refused; return 254; fi
  return 0;;
 *) echo "unexpected deployment operation" >&2; return 1;;
 esac
}
` + table + "\ntable\n"
		out, err := exec.Command(bash, "-c", script).CombinedOutput()
		logged, _ := os.ReadFile(log)
		if err != nil {
			return string(logged), fmt.Errorf("%v: %s", err, out)
		}
		return string(logged), nil
	}
	run := func(exists bool, unavailable int) (string, error) { return runAs(exists, unavailable, "", "False") }
	for _, exists := range []bool{false, true} {
		calls, err := run(exists, 0)
		if err != nil {
			t.Fatalf("table (exists=%v) failed: %v", exists, err)
		}
		if !strings.Contains(calls, "dynamodb update-continuous-backups --table-name demo --point-in-time-recovery-specification PointInTimeRecoveryEnabled=true") {
			t.Fatalf("table (exists=%v) didn't ask for point-in-time recovery:\n%s", exists, calls)
		}
		if strings.Contains(calls, "create-table") == exists {
			t.Fatalf("table (exists=%v) created the table wrongly:\n%s", exists, calls)
		}
	}
	// Deletion protection: turned on when it's off, left alone when it's on,
	// and never asked of the fake vendor's test table.
	const protect = "dynamodb update-table --table-name demo --deletion-protection-enabled"
	for _, c := range []struct {
		fake, protected string
		want            bool
	}{{"", "False", true}, {"", "True", false}, {"vendor", "False", false}} {
		calls, err := runAs(true, 0, c.fake, c.protected)
		if err != nil {
			t.Fatalf("table (fake=%q protected=%s) failed: %v", c.fake, c.protected, err)
		}
		if strings.Contains(calls, protect) != c.want {
			t.Fatalf("table (fake=%q protected=%s) asked for deletion protection=%v, want %v:\n%s", c.fake, c.protected, !c.want, c.want, calls)
		}
	}
	// A new table that refuses recovery for a moment still gets it, and the run goes on.
	calls, err := run(false, 2)
	if err != nil {
		t.Fatalf("a new table that refused recovery twice stopped the run: %v", err)
	}
	if n := strings.Count(calls, "update-continuous-backups"); n != 3 {
		t.Fatalf("recovery asked %d times, want 3 (two refusals, then taken):\n%s", n, calls)
	}
	// One that never takes it fails the run, so the deploy says so.
	if _, err := run(false, 100); err == nil {
		t.Fatal("a table that never took recovery let the run succeed")
	}
}

// Every PLANQUE_ setting the server reads is named by the deploy script, so
// a rerun neither drops nor forgets one. The ones left out run planque
// locally, or (PLANQUE_APP_SECRETS) need each secret granted, which the
// script doesn't do: it isn't set on the function.
func TestDeploymentNamesEverySettingTheServerReads(t *testing.T) {
	script, err := os.ReadFile("../../deploy/aws.sh")
	if err != nil {
		t.Fatal(err)
	}
	notByTheScript := map[string]bool{"PLANQUE_ADDR": true, "PLANQUE_LOCAL_KEY": true, "PLANQUE_APPS_FILE": true, "PLANQUE_APP_SECRETS": true}
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range regexp.MustCompile(`Getenv\("(PLANQUE_[A-Z_]+)"\)`).FindAllStringSubmatch(string(src), -1) {
		if !notByTheScript[m[1]] && !strings.Contains(string(script), `"`+m[1]+`"`) {
			t.Errorf("deploy/aws.sh doesn't name %s, so a rerun drops it", m[1])
		}
	}
}
