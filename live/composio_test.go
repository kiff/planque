//go:build live

package live

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/egress"
	"github.com/kiff/planque/internal/jcs"
)

// TestLiveComposioDiscovery reads public catalog/definition metadata only. It
// creates no session or connected account and never executes a vendor action.
func TestLiveComposioDiscovery(t *testing.T) {
	reference := os.Getenv("LIVE_COMPOSIO_SECRET")
	if reference == "" {
		t.Skip("LIVE_COMPOSIO_SECRET is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	credentials, err := awsapi.FromEnv()
	if err != nil {
		// A developer's named/SSO profile can be resolved by their AWS CLI. The
		// captured credentials remain in process and are never printed or exported.
		command := exec.CommandContext(ctx, "aws", "configure", "export-credentials", "--format", "process")
		raw, e := command.Output()
		if e != nil {
			t.Fatal("AWS credentials could not be resolved")
		}
		var exported struct{ AccessKeyId, SecretAccessKey, SessionToken string }
		e = json.Unmarshal(raw, &exported)
		clear(raw)
		if e != nil || exported.AccessKeyId == "" || exported.SecretAccessKey == "" {
			t.Fatal("AWS credentials could not be resolved")
		}
		credentials = awsapi.Credentials{AccessKey: exported.AccessKeyId, SecretKey: exported.SecretAccessKey, SessionToken: exported.SessionToken}
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}
	httpClient := egress.Client(30 * time.Second)
	store := &awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", HTTP: httpClient, Creds: func() (awsapi.Credentials, error) { return credentials, nil }}
	provider := &composio.Client{HTTP: httpClient, Secret: composio.SecretFromStore(store, reference)}
	apps := map[string]bool{}
	cursors := map[string]bool{}
	cursor := ""
	pages := 0
	skippedApps := 0
	for {
		page, err := provider.Catalog(ctx, "", cursor)
		if err != nil {
			t.Fatal("live catalog discovery failed; no credentials or response bodies are logged")
		}
		pages++
		skippedApps += page.Skipped
		for _, app := range page.Items {
			apps[app.Slug] = true
		}
		if page.Next == "" {
			break
		}
		if cursors[page.Next] || pages >= 100 {
			t.Fatal("live catalog pagination did not terminate")
		}
		cursors[page.Next] = true
		cursor = page.Next
	}
	if len(apps) == 0 {
		t.Fatal("the live catalog was empty")
	}
	toolkit := os.Getenv("LIVE_COMPOSIO_TOOLKIT")
	if toolkit == "" {
		toolkit = "github"
	}
	if !apps[toolkit] {
		t.Fatal("the selected public toolkit was absent from discovery")
	}
	app, err := provider.App(ctx, toolkit)
	if err != nil || app.Slug != toolkit {
		t.Fatal("live app lookup did not match discovery")
	}
	cursor = ""
	cursors = map[string]bool{}
	tools := map[string]bool{}
	var sample *connector.Tool
	toolPages := 0
	skippedTools := 0
	for {
		page, err := provider.Tools(ctx, toolkit, cursor)
		if err != nil {
			t.Fatal("live tool discovery failed conformance")
		}
		toolPages++
		skippedTools += page.Skipped
		for _, tool := range page.Items {
			if tool.Validate() != nil || tool.Toolkit != toolkit {
				t.Fatal("live tool definition was not concrete and versioned")
			}
			tools[tool.Slug] = true
			if sample == nil {
				copy := tool
				sample = &copy
			}
		}
		if page.Next == "" {
			break
		}
		if cursors[page.Next] || toolPages >= 100 {
			t.Fatal("live tool pagination did not terminate")
		}
		cursors[page.Next] = true
		cursor = page.Next
	}
	if sample == nil {
		t.Fatal("the selected public toolkit had no concrete tools")
	}
	pinned, err := provider.Tool(ctx, sample.Slug, sample.Version)
	if err != nil || pinned.Slug != sample.Slug || pinned.Toolkit != toolkit || pinned.Version != sample.Version {
		t.Fatal("live pinned tool lookup did not match discovery")
	}
	expected, e1 := jcs.SHA256(sample.InputSchema)
	actual, e2 := jcs.SHA256(pinned.InputSchema)
	if e1 != nil || e2 != nil || expected != actual {
		t.Fatal("live pinned input schema differed from discovery")
	}
	t.Logf("Read-only conformance: %d apps in %d pages (%d skipped); %d concrete tools in %d pages (%d skipped); pinned schema matched. No authentication session or vendor action was created.", len(apps), pages, skippedApps, len(tools), toolPages, skippedTools)
}
