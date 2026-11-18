package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/kiff/planque/internal/awsapi"
	"github.com/kiff/planque/internal/connector/composio"
	"github.com/kiff/planque/internal/egress"
	"github.com/kiff/planque/internal/planque"
)

// guidesFirst are the apps whose guides come first: what Explore shows and
// what people connect first. Every other app follows, alphabetically.
var guidesFirst = []string{"gmail", "googlesheets", "stripe", "github", "linkedin", "qonto",
	"googlecalendar", "googledrive", "googledocs", "outlook", "slack", "notion", "hubspot", "shopify", "calendly", "trello", "asana", "airtable"}

var guideSlug = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// guides writes every app's generated guide part, and the app list in
// batches: `planque guides -out deploy/guides`. It reads the managed apps
// with PLANQUE_COMPOSIO_SECRET (as the server does) and the reviewed
// entries from -catalog. Nothing of an owner's is read.
func guides(args []string) error {
	fs := flag.NewFlagSet("guides", flag.ContinueOnError)
	out := fs.String("out", "deploy/guides", "where the guides go")
	catalog := fs.String("catalog", "deploy/catalog.json", "the reviewed catalog entries")
	size := fs.Int("batch", 10, "apps per batch")
	date := fs.String("date", time.Now().UTC().Format(time.DateOnly), "the Checked date")
	if err := fs.Parse(args); err != nil {
		return err
	}
	var apps []planque.GuideApp
	if reference := strings.TrimSpace(os.Getenv("PLANQUE_COMPOSIO_SECRET")); reference != "" {
		region := cmp.Or(os.Getenv("AWS_REGION"), "eu-west-1")
		p := &composio.Client{HTTP: egress.Client(30 * time.Second), Secret: composio.SecretFromStore(&awsapi.Client{Service: "secretsmanager", Region: region,
			Target: "secretsmanager", JSONType: "application/x-amz-json-1.1", Creds: awsapi.FromEnv}, reference)}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		managed, err := planque.CollectGuideApps(ctx, p)
		if err != nil {
			return err
		}
		apps = append(apps, managed...)
	} else {
		fmt.Fprintln(os.Stderr, "PLANQUE_COMPOSIO_SECRET unset: only the reviewed entries")
	}
	if *catalog != "" {
		raw, err := os.ReadFile(*catalog)
		if err != nil {
			return err
		}
		entries, err := planque.DecodeCatalog(raw, nil)
		if err != nil {
			return err
		}
		ids := make([]string, 0, len(entries))
		for id := range entries {
			ids = append(ids, id)
		}
		slices.Sort(ids)
		for _, id := range ids {
			apps = append(apps, planque.EntryGuideApp(id, entries[id]))
		}
	}
	if len(apps) == 0 {
		return errors.New("no apps read")
	}
	gen := filepath.Join(*out, "generated")
	if err := os.MkdirAll(gen, 0o755); err != nil {
		return err
	}
	for _, a := range apps {
		name := strings.ReplaceAll(strings.ToLower(a.Slug), "@", "-")
		if !guideSlug.MatchString(name) {
			return fmt.Errorf("an app's slug can't name a file: %q", a.Slug)
		}
		if err := os.WriteFile(filepath.Join(gen, name+".md"), []byte(planque.GeneratedGuide(a, *date)), 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(*out, "apps.md"), []byte(planque.GuideBatches(apps, guidesFirst, *size, *date)), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%d apps written to %s\n", len(apps), *out)
	return nil
}
