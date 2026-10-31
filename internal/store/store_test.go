package store

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kiff/planque/internal/awsapi"
)

// An empty prefix queries the whole partition: DynamoDB refuses an empty
// begins_with value.
func TestDynamoQueryWithAnEmptyPrefix(t *testing.T) {
	var conds []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		var in struct {
			KeyConditionExpression    string
			ExpressionAttributeValues map[string]any
		}
		_ = json.Unmarshal(b, &in)
		conds = append(conds, in.KeyConditionExpression)
		if _, ok := in.ExpressionAttributeValues[":p"]; ok != strings.Contains(in.KeyConditionExpression, "begins_with") {
			t.Errorf("%q with values %v", in.KeyConditionExpression, in.ExpressionAttributeValues)
		}
		_, _ = io.WriteString(w, `{"Items":[]}`)
	}))
	defer srv.Close()
	d := &Dynamo{Table: "t", API: &awsapi.Client{Service: "dynamodb", Region: "eu-west-1", Target: "DynamoDB_20120810",
		JSONType: "application/x-amz-json-1.0", Endpoint: srv.URL,
		Creds: func() (awsapi.Credentials, error) {
			return awsapi.Credentials{AccessKey: "AKID", SecretKey: "secret"}, nil
		}}}
	for _, prefix := range []string{"", "GRANT#"} {
		if _, err := d.Query(context.Background(), "pk", prefix); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := d.QueryTo(context.Background(), "DUE", "2026-10-20T10:00:00Z#~#~"); err != nil {
		t.Fatal(err)
	}
	if len(conds) != 3 || conds[0] != "pk = :pk" || conds[1] != "pk = :pk AND begins_with(sk, :p)" || conds[2] != "pk = :pk AND sk <= :max" {
		t.Fatalf("conditions: %q", conds)
	}
}

// QueryTo reads a range up to its bound, in order, and nothing past it.
func TestMemoryQueryToStopsAtItsBound(t *testing.T) {
	m := &Memory{}
	for _, sk := range []string{"2026-10-20T09:00:00Z#a#1", "2026-10-20T11:00:00Z#a#2", "2026-10-20T10:00:00Z#b#3"} {
		if err := m.Put(context.Background(), Item{PK: "DUE", SK: sk, Doc: []byte("{}")}, Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := m.QueryTo(context.Background(), "DUE", "2026-10-20T10:00:00Z#~#~")
	if err != nil || len(got) != 2 || got[0].SK != "2026-10-20T09:00:00Z#a#1" || got[1].SK != "2026-10-20T10:00:00Z#b#3" {
		t.Fatalf("QueryTo = %+v %v", got, err)
	}
}

// A read is sent again when the connection drops before DynamoDB answers; a
// write isn't, since one whose answer was lost may already have landed.
func TestDynamoReadsAgainWhenTheConnectionDrops(t *testing.T) {
	// Written from the server's handler goroutines, read by the test.
	var mu sync.Mutex
	calls := map[string]int{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.Header.Get("X-Amz-Target")
		mu.Lock()
		calls[target]++
		first := calls[target] == 1
		mu.Unlock()
		if first {
			// The connection drops with no answer: the client reads EOF.
			conn, _, _ := w.(http.Hijacker).Hijack()
			_ = conn.Close()
			return
		}
		_, _ = io.WriteString(w, `{"Items":[],"Item":{"pk":{"S":"pk"},"sk":{"S":"sk"},"doc":{"S":"{}"}}}`)
	}))
	defer srv.Close()
	d := &Dynamo{Table: "t", API: &awsapi.Client{Service: "dynamodb", Region: "eu-west-1", Target: "DynamoDB_20120810",
		JSONType: "application/x-amz-json-1.0", Endpoint: srv.URL,
		Creds: func() (awsapi.Credentials, error) {
			return awsapi.Credentials{AccessKey: "AKID", SecretKey: "secret"}, nil
		}}}
	ctx := context.Background()
	if _, err := d.Query(ctx, "pk", "GRANT#"); err != nil {
		t.Fatalf("a query after one dropped connection: %v", err)
	}
	if _, _, err := d.Get(ctx, "pk", "sk"); err != nil {
		t.Fatalf("a get after one dropped connection: %v", err)
	}
	if err := d.Put(ctx, Item{PK: "pk", SK: "sk", Doc: []byte("{}")}, Cond{}); err == nil {
		t.Fatal("a put whose connection dropped succeeded: it must not be sent again")
	}
	mu.Lock()
	defer mu.Unlock()
	if calls["DynamoDB_20120810.Query"] != 2 || calls["DynamoDB_20120810.GetItem"] != 2 || calls["DynamoDB_20120810.PutItem"] != 1 {
		t.Fatalf("calls: %v", calls)
	}
}
