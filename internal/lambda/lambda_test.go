package lambda

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

type wsRecorder struct{ got []WSEvent }

func (w *wsRecorder) WebSocket(_ context.Context, ev WSEvent) int {
	w.got = append(w.got, ev)
	return 200
}

// A WebSocket API's event goes to the WebSocket handler, never through the
// HTTP routes; a function URL request, whatever its path, goes to the
// routes. With no WebSocket handler, its events are refused.
func TestAWebSocketEventNeverReachesTheHTTPRoutes(t *testing.T) {
	routed := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { routed++ })
	ws := &wsRecorder{}
	event := []byte(`{"body":"{\"type\":\"open\"}","requestContext":{"routeKey":"$default","connectionId":"c1","identity":{"sourceIp":"192.0.2.1"}}}`)
	out, err := invoke(context.Background(), h, ws, event)
	var res map[string]int
	_ = json.Unmarshal(out, &res)
	if err != nil || res["statusCode"] != 200 || routed != 0 || len(ws.got) != 1 || ws.got[0].Connection != "c1" || string(ws.got[0].Body) != `{"type":"open"}` {
		t.Fatalf("%s %v routed %d %+v", out, err, routed, ws.got)
	}
	url := []byte(`{"rawPath":"/x","requestContext":{"routeKey":"$default","connectionId":"c1","http":{"method":"POST","sourceIp":"192.0.2.1"}}}`)
	if _, err := invoke(context.Background(), h, ws, url); err != nil || routed != 1 || len(ws.got) != 1 {
		t.Fatalf("a function URL request: %v routed %d", err, routed)
	}
	out, _ = invoke(context.Background(), h, nil, event)
	_ = json.Unmarshal(out, &res)
	if res["statusCode"] != 403 {
		t.Fatalf("no WebSocket handler: %s", out)
	}
}
