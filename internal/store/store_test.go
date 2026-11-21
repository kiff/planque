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
	"time"

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

func testDynamo(url string) *Dynamo {
	return &Dynamo{Table: "t", API: &awsapi.Client{Service: "dynamodb", Region: "eu-west-1", Target: "DynamoDB_20120810",
		JSONType: "application/x-amz-json-1.0", Endpoint: url,
		Creds: func() (awsapi.Credentials, error) {
			return awsapi.Credentials{AccessKey: "AKID", SecretKey: "secret"}, nil
		}},
		Backoff: func(int) time.Duration { return 0 }}
}

// A throttled read is sent again, page by page, keeping its place; one that
// stays throttled fails after throttleTries. A throttled write is sent again
// too: DynamoDB wrote nothing (a write whose answer was lost never is).
func TestDynamoReadsAgainWhenThrottled(t *testing.T) {
	var mu sync.Mutex
	var bodies []map[string]any
	throttle := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in map[string]any
		_ = json.NewDecoder(r.Body).Decode(&in)
		mu.Lock()
		defer mu.Unlock()
		bodies = append(bodies, in)
		if throttle > 0 {
			throttle--
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"__type":"com.amazonaws.dynamodb.v20120810#ThrottlingException","message":"Throughput exceeds the current capacity"}`)
			return
		}
		if _, ok := in["ExclusiveStartKey"]; ok {
			_, _ = io.WriteString(w, `{"Items":[{"pk":{"S":"TENANT#a"},"sk":{"S":"OP#g#2"},"doc":{"B":"e30="}}]}`)
			return
		}
		_, _ = io.WriteString(w, `{"Items":[{"pk":{"S":"TENANT#a"},"sk":{"S":"OP#g#1"},"doc":{"B":"e30="}}],"LastEvaluatedKey":{"pk":{"S":"TENANT#a"},"sk":{"S":"OP#g#1"}}}`)
	}))
	defer srv.Close()
	d, ctx := testDynamo(srv.URL), context.Background()
	mu.Lock()
	throttle = 2
	mu.Unlock()
	items, next, err := d.ScanPage(ctx, "OP#", "", 100)
	if err != nil || len(items) != 1 || items[0].SK != "OP#g#1" || next == "" {
		t.Fatalf("first page after two throttles: %v %q %v", items, next, err)
	}
	mu.Lock()
	first := bodies[len(bodies)-1]
	if len(bodies) != 3 || first["Limit"] != float64(100) || first["ConsistentRead"] != nil {
		t.Fatalf("calls %d, last %v: want 3 calls, Limit 100, eventually consistent", len(bodies), first)
	}
	throttle = 1
	mu.Unlock()
	items, next, err = d.ScanPage(ctx, "OP#", next, 100)
	if err != nil || len(items) != 1 || items[0].SK != "OP#g#2" || next != "" {
		t.Fatalf("second page from the cursor: %v %q %v", items, next, err)
	}
	mu.Lock()
	if start, _ := json.Marshal(bodies[len(bodies)-1]["ExclusiveStartKey"]); !strings.Contains(string(start), "OP#g#1") {
		t.Fatalf("the cursor wasn't sent: %s", start)
	}
	throttle = 100
	mu.Unlock()
	if _, _, err := d.ScanPage(ctx, "OP#", "", 100); !awsapi.Is(err, "ThrottlingException") {
		t.Fatalf("a read throttled every time: %v", err)
	}
	mu.Lock()
	n := len(bodies)
	mu.Unlock()
	// A throttle is DynamoDB's answer that nothing was written: a write it
	// throttled is sent again, a conditional one included, up to throttleTries.
	if err := d.Put(ctx, Item{PK: "pk", SK: "sk", Doc: []byte("{}")}, Cond{}); !awsapi.Is(err, "ThrottlingException") {
		t.Fatalf("a put throttled every time: %v", err)
	}
	mu.Lock()
	if len(bodies) != n+throttleTries {
		t.Fatalf("a put throttled every time was sent %d times, want %d", len(bodies)-n, throttleTries)
	}
	throttle, n = 1, len(bodies)
	mu.Unlock()
	ver := int64(1)
	if err := d.Put(ctx, Item{PK: "pk", SK: "sk", Doc: []byte("{}"), Ver: 2}, Cond{Ver: &ver}); err != nil {
		t.Fatalf("a conditional put throttled once: %v", err)
	}
	mu.Lock()
	if len(bodies) != n+2 {
		t.Fatalf("a put throttled once was sent %d times, want 2", len(bodies)-n)
	}
	throttle, n = 1, len(bodies)
	mu.Unlock()
	if _, _, err := d.Delete(ctx, "pk", "sk"); err != nil {
		t.Fatalf("a delete throttled once: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != n+2 {
		t.Fatalf("a delete throttled once was sent %d times, want 2", len(bodies)-n)
	}
}

// Memory's pages cover every pk, in order, each item once.
func TestMemoryScanPages(t *testing.T) {
	m, ctx := &Memory{}, context.Background()
	for _, k := range [][2]string{{"TENANT#b", "OP#1"}, {"TENANT#a", "OP#2"}, {"TENANT#a", "OP#1"}, {"TENANT#a", "GRANT#1"}} {
		if err := m.Put(ctx, Item{PK: k[0], SK: k[1], Doc: []byte("{}")}, Cond{}); err != nil {
			t.Fatal(err)
		}
	}
	var got []string
	cursor := ""
	for range 10 {
		items, next, err := m.ScanPage(ctx, "OP#", cursor, 2)
		if err != nil {
			t.Fatal(err)
		}
		for _, it := range items {
			got = append(got, it.PK+" "+it.SK)
		}
		if cursor = next; cursor == "" {
			break
		}
	}
	if strings.Join(got, ",") != "TENANT#a OP#1,TENANT#a OP#2,TENANT#b OP#1" {
		t.Fatalf("pages: %v", got)
	}
}

// DeleteIf removes an item only at the version read: one written since
// stays, and none there is a failed condition.
func TestDeleteIfOnlyAtTheVersionRead(t *testing.T) {
	ctx := context.Background()
	m := &Memory{}
	if err := m.Put(ctx, Item{PK: "p", SK: "s", Doc: []byte("{}"), Ver: 2}, Cond{}); err != nil {
		t.Fatal(err)
	}
	if err := m.DeleteIf(ctx, "p", "s", 1); err != ErrCondition {
		t.Fatalf("an older version: %v", err)
	}
	if _, ok, _ := m.Get(ctx, "p", "s"); !ok {
		t.Fatal("deleted at another version")
	}
	if err := m.DeleteIf(ctx, "p", "s", 2); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := m.Get(ctx, "p", "s"); ok {
		t.Fatal("left at its version")
	}
	if err := m.DeleteIf(ctx, "p", "s", 2); err != ErrCondition {
		t.Fatalf("none there: %v", err)
	}
}

// Dynamo's DeleteIf sends the version as the delete's condition, and a
// failed one is ErrCondition.
func TestDynamoDeleteIfSendsTheVersion(t *testing.T) {
	var got struct {
		ConditionExpression       string
		ExpressionAttributeValues map[string]map[string]string
	}
	fail := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		if fail {
			w.Header().Set("Content-Type", "application/x-amz-json-1.0")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"__type":"com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException","message":"The conditional request failed"}`)
			return
		}
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()
	d := &Dynamo{Table: "t", API: &awsapi.Client{Service: "dynamodb", Region: "eu-west-1", Target: "DynamoDB_20120810",
		JSONType: "application/x-amz-json-1.0", Endpoint: srv.URL,
		Creds: func() (awsapi.Credentials, error) {
			return awsapi.Credentials{AccessKey: "AKID", SecretKey: "secret"}, nil
		}}}
	if err := d.DeleteIf(context.Background(), "p", "s", 7); err != nil {
		t.Fatal(err)
	}
	if got.ConditionExpression != "ver = :ver" || got.ExpressionAttributeValues[":ver"]["N"] != "7" {
		t.Fatalf("sent %q %v", got.ConditionExpression, got.ExpressionAttributeValues)
	}
	fail = true
	if err := d.DeleteIf(context.Background(), "p", "s", 7); err != ErrCondition {
		t.Fatalf("a failed condition: %v", err)
	}
}
