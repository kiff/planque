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

func classTool(toolkit, slug string, read bool) string {
	tags := `[]`
	if read {
		tags = `["readOnlyHint"]`
	}
	return `{"slug":"` + slug + `","name":"` + slug + `","toolkit":{"slug":"` + toolkit + `"},"version":"20261101_00","input_parameters":{"type":"object","properties":{"x":{"type":"string"}}},"tags":` + tags + `}`
}

// A tool in a held class (delete, share, send-to-many, money) or one planque
// can't place never comes back from a search or an app's tools, can't be read
// one at a time and is never executed, even on a connection that holds it;
// reads and ordinary writes are served as before. The skipped ones are counted.
func TestAHeldClassToolNeverComesBackNorRuns(t *testing.T) {
	old := connector.HoldClasses
	defer func() { connector.HoldClasses = old }()
	connector.HoldClasses = true
	sends := 0
	items := `{"items":[` + strings.Join([]string{
		classTool("notion", "NOTION_DELETE_BLOCK", false), classTool("notion", "NOTION_SHARE_PAGE", false), classTool("notion", "NOTION_SEND_CAMPAIGN", false),
		classTool("notion", "NOTION_REFUND_X", false), classTool("notion", "NOTION_FROBNICATE_X", false),
		classTool("notion", "NOTION_FETCH_PAGE", true), classTool("notion", "NOTION_ADD_PAGE_CONTENT", false),
	}, ",") + `]}`
	c := testClient(func(r *http.Request) (*http.Response, error) {
		switch path := strings.TrimPrefix(r.URL.Path, "/api/v3.1"); {
		case path == "/auth_configs":
			return response(200, `{"items":[{"id":"ac_notion","auth_scheme":"OAUTH2","status":"ENABLED","toolkit":{"slug":"notion"}}]}`), nil
		case strings.HasPrefix(path, "/toolkits/"):
			slug := strings.TrimPrefix(path, "/toolkits/")
			return response(200, `{"slug":"`+slug+`","name":"`+slug+`"}`), nil
		case path == "/tools":
			return response(200, items), nil
		case path == "/tools/NOTION_DELETE_BLOCK":
			return response(200, classTool("notion", "NOTION_DELETE_BLOCK", false)), nil
		case path == "/connected_accounts/account_one":
			return response(200, `{"id":"account_one","user_id":"owner_one","status":"ACTIVE","toolkit":{"slug":"notion"}}`), nil
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
	const served = "NOTION_FETCH_PAGE NOTION_ADD_PAGE_CONTENT"
	found, err := c.Search(ctx, "page", nil, "")
	if err != nil || slugs(found.Tools) != served || found.Skipped != 5 {
		t.Fatalf("search: %s %+v %v", slugs(found.Tools), found, err)
	}
	page, err := c.Tools(ctx, "notion", "")
	if err != nil || slugs(page.Items) != served || page.Skipped != 5 {
		t.Fatalf("the app's tools: %s %+v %v", slugs(page.Items), page, err)
	}
	if _, err := c.Tool(ctx, "NOTION_DELETE_BLOCK", "20261101_00"); !errors.Is(err, connector.ErrDefinition) {
		t.Fatalf("read one at a time: %v", err)
	}
	held := connector.Tool{Slug: "NOTION_DELETE_BLOCK", Toolkit: "notion", Version: "20261101_00", Name: "Delete a block", InputSchema: json.RawMessage(`{"type":"object","properties":{"x":{"type":"string"}}}`)}
	if _, err := c.Execute(ctx, connector.Identity{User: "owner_one", Toolkit: "notion"}, "account_one", held, json.RawMessage(`{"x":"b"}`), func() bool { return true }); !errors.Is(err, connector.ErrDefinition) || sends != 0 {
		t.Fatalf("executed: sends=%d err=%v", sends, err)
	}
}
