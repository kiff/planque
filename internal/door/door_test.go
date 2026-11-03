package door

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kiff/planque/internal/connector"
)

type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestDoorSignsExactMessageAndRetriesOnlyTransientAnswers(t *testing.T) {
	for _, status := range []int{0, 202, 200, 302, 400, 401, 404, 409, 429, 500, 503} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			sends := 0
			now := time.Unix(1800000000, 0)
			c := &HTTP{URL: "https://receiver.example/v1/doors/planque", Secret: func(context.Context) (string, error) { return "synthetic-handoff", nil }, Now: func() time.Time { return now }, Client: &http.Client{Transport: transport(func(r *http.Request) (*http.Response, error) {
				sends++
				raw, _ := io.ReadAll(r.Body)
				if r.URL.Host != "receiver.example" || r.Header.Get("Planque-Event-Id") != "event_one" || connector.VerifyWebhook("synthetic-handoff", r.Header.Get("Planque-Event-Id"), r.Header.Get("Planque-Event-Timestamp"), r.Header.Get("Planque-Event-Signature"), raw, now) != nil {
					t.Error("handoff signature or receiver changed")
				}
				if status == 0 {
					return nil, errors.New("lost answer")
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Location": []string{"https://other.example"}}, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
			})}}
			got, err := c.Deliver(context.Background(), Message{V: "planque.door.v1", ID: "event_one"})
			if err != nil || sends != 1 || got.Taken != (status == 202) || got.Retry != (status == 0 || status >= 500) {
				t.Fatalf("%d: %+v %v sends%d", status, got, err, sends)
			}
		})
	}
}
