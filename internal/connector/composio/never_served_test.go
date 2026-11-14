package composio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/connector"
)

func stripeTool(toolkit, slug string) string {
	return `{"slug":"` + slug + `","name":"` + slug + `","toolkit":{"slug":"` + toolkit + `"},"version":"20261101_00","input_parameters":{"type":"object","properties":{"invoice":{"type":"string"}}},"tags":[]}`
}

// Stripe's void, delete, finalize, pay and write-off tools never come back
// from a search or the app's tools, can't be read one at a time, and are
// never executed; Stripe's other tools, and another app's tool of the same
// name, are served as before.
func TestStripesNeverServedToolsNeverComeBackNorRun(t *testing.T) {
	sends := 0
	items := `{"items":[` + strings.Join([]string{
		stripeTool("stripe", "STRIPE_VOID_INVOICE"), stripeTool("stripe", "STRIPE_PAY_INVOICE"),
		stripeTool("stripe", "STRIPE_LIST_INVOICES"), stripeTool("stripe", "STRIPE_SEND_INVOICE"), stripeTool("stripe", "STRIPE_CREATE_REFUND"),
	}, ",") + `]}`
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch path := strings.TrimPrefix(r.URL.Path, "/api/v3.1"); {
		case path == "/auth_configs":
			return response(200, `{"items":[{"id":"ac_stripe","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"stripe"}}]}`), nil
		case strings.HasPrefix(path, "/toolkits/"):
			slug := strings.TrimPrefix(path, "/toolkits/")
			return response(200, `{"slug":"`+slug+`","name":"`+slug+`"}`), nil
		case path == "/tools" && r.URL.Query().Get("toolkit_slug") == "other_app":
			return response(200, `{"items":[`+stripeTool("other_app", "STRIPE_VOID_INVOICE")+`]}`), nil
		case path == "/tools":
			return response(200, items), nil
		case path == "/tools/STRIPE_PAY_INVOICE":
			return response(200, stripeTool("stripe", "STRIPE_PAY_INVOICE")), nil
		case path == "/connected_accounts/account_one":
			return response(200, `{"id":"account_one","user_id":"owner_one","status":"ACTIVE","toolkit":{"slug":"stripe"}}`), nil
		case strings.HasPrefix(path, "/tools/execute/"):
			sends++
			return response(200, `{"data":{},"successful":true}`), nil
		}
		t.Fatalf("unexpected request %s", r.URL)
		return nil, nil
	})
	ctx := context.Background()
	slugs := func(tools []connector.Tool) string {
		var out []string
		for _, tl := range tools {
			out = append(out, tl.Slug)
		}
		return strings.Join(out, " ")
	}
	const served = "STRIPE_LIST_INVOICES STRIPE_SEND_INVOICE STRIPE_CREATE_REFUND"
	found, err := c.Search(ctx, "invoice", nil, "")
	if err != nil || slugs(found.Tools) != served || found.Skipped != 2 {
		t.Fatalf("search: %s %+v %v", slugs(found.Tools), found, err)
	}
	page, err := c.Tools(ctx, "stripe", "")
	if err != nil || slugs(page.Items) != served || page.Skipped != 2 {
		t.Fatalf("the app's tools: %s %v", slugs(page.Items), err)
	}
	// Per app: another app's tool of the same name is served.
	other, err := c.Tools(ctx, "other_app", "")
	if err != nil || slugs(other.Items) != "STRIPE_VOID_INVOICE" {
		t.Fatalf("another app's tool: %s %v", slugs(other.Items), err)
	}
	if _, err := c.Tool(ctx, "STRIPE_PAY_INVOICE", "20261101_00"); !errors.Is(err, connector.ErrDefinition) {
		t.Fatalf("read one at a time: %v", err)
	}
	// One a connection holds from before is never executed.
	held := connector.Tool{Slug: "STRIPE_PAY_INVOICE", Toolkit: "stripe", Version: "20261101_00", Name: "Pay an invoice", InputSchema: json.RawMessage(`{"type":"object","properties":{"invoice":{"type":"string"}}}`)}
	if _, err := c.Execute(ctx, connector.Identity{User: "owner_one", Toolkit: "stripe"}, "account_one", held, json.RawMessage(`{"invoice":"in_1"}`), func() bool { return true }); !errors.Is(err, connector.ErrDefinition) || sends != 0 {
		t.Fatalf("executed: sends=%d err=%v", sends, err)
	}
}
