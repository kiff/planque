package store

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
		if _, ok := in.ExpressionAttributeValues[":p"]; ok != (in.KeyConditionExpression != "pk = :pk") {
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
	if len(conds) != 2 || conds[0] != "pk = :pk" || conds[1] != "pk = :pk AND begins_with(sk, :p)" {
		t.Fatalf("conditions: %q", conds)
	}
}
