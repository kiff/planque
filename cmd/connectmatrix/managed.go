package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/egress"
)

// managedProvider is the managed sign-in provider, read only, when the run has its key: the secret's reference in
// LIVE_COMPOSIO_SECRET and AWS credentials to read it, as the live tests do. Without them, the column says so.
func managedProvider() *composio.Client {
	reference := os.Getenv("LIVE_COMPOSIO_SECRET")
	if reference == "" {
		return nil
	}
	credentials, err := awsapi.FromEnv()
	if err != nil {
		return nil
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}
	httpClient := egress.Client(30 * time.Second)
	store := &awsapi.Client{Service: "secretsmanager", Region: region, Target: "secretsmanager", JSONType: "application/x-amz-json-1.1",
		HTTP: httpClient, Creds: func() (awsapi.Credentials, error) { return credentials, nil }}
	return &composio.Client{HTTP: httpClient, Secret: composio.SecretFromStore(store, reference)}
}

// managed says whether the provider can sign someone in to the app: its metadata only, never a session, an
// account or a call.
func managed(p *composio.Client, slug string) string {
	if slug == "" {
		return "-"
	}
	if p == nil {
		return "not checked (no provider key in this run)"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a, err := p.App(ctx, slug)
	switch {
	case errors.Is(err, connector.ErrNotSetUp):
		return "❌ not at the provider"
	case err != nil:
		return "? couldn't be read"
	}
	switch a.Setup {
	case connector.SetupReady:
		return "✅ ready"
	case connector.SetupOnDemand:
		return "✅ on the owner's pick"
	case connector.SetupNeedsSetup:
		return "❌ needs an operator's own credentials"
	}
	return "? " + a.Setup
}
